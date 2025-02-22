// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sync/atomic"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/thriftrw/protocol"
	"go.uber.org/thriftrw/wire"

	"github.com/uber/cadence/.gen/go/indexer"
	"github.com/uber/cadence/.gen/go/replicator"
	"github.com/uber/cadence/client/admin"
	"github.com/uber/cadence/common"
	"github.com/uber/cadence/common/config"
	"github.com/uber/cadence/common/persistence"
	"github.com/uber/cadence/common/types"
	"github.com/uber/cadence/common/types/mapper/thrift"
)

type (
	filterFn              func(*types.ReplicationTask) bool
	filterFnForVisibility func(*indexer.Message) bool

	kafkaMessageType int

	historyV2Task struct {
		Task         *types.ReplicationTask
		Events       []*types.HistoryEvent
		NewRunEvents []*types.HistoryEvent
	}
)

const (
	kafkaMessageTypeReplicationTask kafkaMessageType = iota
	kafkaMessageTypeVisibilityMsg
)

const (
	bufferSize                       = 8192
	preambleVersion0            byte = 0x59
	malformedMessage                 = "Input was malformed"
	chanBufferSize                   = 10000
	maxRereplicateEventID            = 999999
	defaultResendContextTimeout      = 30 * time.Second
)

var (
	r = regexp.MustCompile(`Partition: .*?, Offset: .*?, Key: .*?`)
)

type writerChannel struct {
	Type                   kafkaMessageType
	ReplicationTaskChannel chan *types.ReplicationTask
	VisibilityMsgChannel   chan *indexer.Message
}

func newWriterChannel(messageType kafkaMessageType) *writerChannel {
	ch := &writerChannel{
		Type: messageType,
	}
	switch messageType {
	case kafkaMessageTypeReplicationTask:
		ch.ReplicationTaskChannel = make(chan *types.ReplicationTask, chanBufferSize)
	case kafkaMessageTypeVisibilityMsg:
		ch.VisibilityMsgChannel = make(chan *indexer.Message, chanBufferSize)
	}
	return ch
}

func (ch *writerChannel) Close() {
	if ch.ReplicationTaskChannel != nil {
		close(ch.ReplicationTaskChannel)
	}
	if ch.VisibilityMsgChannel != nil {
		close(ch.VisibilityMsgChannel)
	}
}

// AdminKafkaParse parses the output of k8read and outputs replication tasks
func AdminKafkaParse(c *cli.Context) {
	inputFile := getInputFile(c.String(FlagInputFile))
	outputFile := getOutputFile(c.String(FlagOutputFilename))

	defer inputFile.Close()
	defer outputFile.Close()

	readerCh := make(chan []byte, chanBufferSize)
	writerCh := newWriterChannel(kafkaMessageType(c.Int(FlagMessageType)))
	doneCh := make(chan struct{})
	serializer := persistence.NewPayloadSerializer()

	var skippedCount int32
	skipErrMode := c.Bool(FlagSkipErrorMode)

	go startReader(inputFile, readerCh)
	go startParser(readerCh, writerCh, skipErrMode, &skippedCount)
	go startWriter(outputFile, writerCh, doneCh, &skippedCount, serializer, c)

	<-doneCh

	if skipErrMode {
		fmt.Printf("%v messages were skipped due to errors in parsing", atomic.LoadInt32(&skippedCount))
	}
}

func buildFilterFn(workflowID, runID string) filterFn {
	return func(task *types.ReplicationTask) bool {
		if len(workflowID) != 0 || len(runID) != 0 {
			if task.GetHistoryTaskV2Attributes() == nil {
				return false
			}
		}
		if len(workflowID) != 0 && task.GetHistoryTaskV2Attributes().WorkflowID != workflowID {
			return false
		}
		if len(runID) != 0 && task.GetHistoryTaskV2Attributes().RunID != runID {
			return false
		}
		return true
	}
}

func buildFilterFnForVisibility(workflowID, runID string) filterFnForVisibility {
	return func(msg *indexer.Message) bool {
		if len(workflowID) != 0 && msg.GetWorkflowID() != workflowID {
			return false
		}
		if len(runID) != 0 && msg.GetRunID() != runID {
			return false
		}
		return true
	}
}

func getOutputFile(outputFile string) *os.File {
	if len(outputFile) == 0 {
		return os.Stdout
	}
	f, err := os.Create(outputFile)
	if err != nil {
		ErrorAndExit("failed to create output file", err)
	}

	return f
}

func startReader(file *os.File, readerCh chan<- []byte) {
	defer close(readerCh)
	reader := bufio.NewReader(file)

	for {
		buf := make([]byte, bufferSize)
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				ErrorAndExit("Failed to read from reader", err)
			} else {
				break
			}

		}
		buf = buf[:n]
		readerCh <- buf
	}
}

func startParser(readerCh <-chan []byte, writerCh *writerChannel, skipErrors bool, skippedCount *int32) {
	defer writerCh.Close()

	var buffer []byte
Loop:
	for {
		select {
		case data, ok := <-readerCh:
			if !ok {
				break Loop
			}
			buffer = append(buffer, data...)
			data, nextBuffer := splitBuffer(buffer)
			buffer = nextBuffer
			parse(data, skipErrors, skippedCount, writerCh)
		}
	}
	parse(buffer, skipErrors, skippedCount, writerCh)
}

func startWriter(
	outputFile *os.File,
	writerCh *writerChannel,
	doneCh chan struct{},
	skippedCount *int32,
	serializer persistence.PayloadSerializer,
	c *cli.Context,
) {

	defer close(doneCh)

	skipErrMode := c.Bool(FlagSkipErrorMode)
	headerMode := c.Bool(FlagHeadersMode)

	switch writerCh.Type {
	case kafkaMessageTypeReplicationTask:
		writeReplicationTask(outputFile, writerCh, skippedCount, skipErrMode, headerMode, serializer, c)
	case kafkaMessageTypeVisibilityMsg:
		writeVisibilityMessage(outputFile, writerCh, skippedCount, skipErrMode, headerMode, c)
	}
}

func writeReplicationTask(
	outputFile *os.File,
	writerCh *writerChannel,
	skippedCount *int32,
	skipErrMode bool,
	headerMode bool,
	serializer persistence.PayloadSerializer,
	c *cli.Context,
) {
	filter := buildFilterFn(c.String(FlagWorkflowID), c.String(FlagRunID))
Loop:
	for {
		select {
		case task, ok := <-writerCh.ReplicationTaskChannel:
			if !ok {
				break Loop
			}
			if filter(task) {
				jsonStr, err := decodeReplicationTask(task, serializer)
				if err != nil {
					if !skipErrMode {
						ErrorAndExit(malformedMessage, fmt.Errorf("failed to encode into json, err: %v", err))
					} else {
						atomic.AddInt32(skippedCount, 1)
						continue Loop
					}
				}

				var outStr string
				if !headerMode {
					outStr = string(jsonStr)
				} else {
					outStr = fmt.Sprintf(
						"%v, %v, %v",
						task.GetHistoryTaskV2Attributes().DomainID,
						task.GetHistoryTaskV2Attributes().WorkflowID,
						task.GetHistoryTaskV2Attributes().RunID,
					)
				}
				_, err = outputFile.WriteString(fmt.Sprintf("%v\n", outStr))
				if err != nil {
					ErrorAndExit("Failed to write to file", fmt.Errorf("err: %v", err))
				}
			}
		}
	}
}

func writeVisibilityMessage(
	outputFile *os.File,
	writerCh *writerChannel,
	skippedCount *int32,
	skipErrMode bool,
	headerMode bool,
	c *cli.Context,
) {
	filter := buildFilterFnForVisibility(c.String(FlagWorkflowID), c.String(FlagRunID))
Loop:
	for {
		select {
		case msg, ok := <-writerCh.VisibilityMsgChannel:
			if !ok {
				break Loop
			}
			if filter(msg) {
				jsonStr, err := json.Marshal(msg)
				if err != nil {
					if !skipErrMode {
						ErrorAndExit(malformedMessage, fmt.Errorf("failed to encode into json, err: %v", err))
					} else {
						atomic.AddInt32(skippedCount, 1)
						continue Loop
					}
				}

				var outStr string
				if !headerMode {
					outStr = string(jsonStr)
				} else {
					outStr = fmt.Sprintf(
						"%v, %v, %v, %v, %v",
						msg.GetDomainID(),
						msg.GetWorkflowID(),
						msg.GetRunID(),
						msg.GetMessageType().String(),
						msg.GetVersion(),
					)
				}
				_, err = outputFile.WriteString(fmt.Sprintf("%v\n", outStr))
				if err != nil {
					ErrorAndExit("Failed to write to file", fmt.Errorf("err: %v", err))
				}
			}
		}
	}
}

func splitBuffer(buffer []byte) ([]byte, []byte) {
	matches := r.FindAllIndex(buffer, -1)
	if len(matches) == 0 {
		ErrorAndExit(malformedMessage, errors.New("header not found, did you generate dump with -v"))
	}
	splitIndex := matches[len(matches)-1][0]
	return buffer[:splitIndex], buffer[splitIndex:]
}

func parse(bytes []byte, skipErrors bool, skippedCount *int32, writerCh *writerChannel) {
	messages, skippedGetMsgCount := getMessages(bytes, skipErrors)
	switch writerCh.Type {
	case kafkaMessageTypeReplicationTask:
		msgs, skippedDeserializeCount := deserializeMessages(messages, skipErrors)
		atomic.AddInt32(skippedCount, skippedGetMsgCount+skippedDeserializeCount)
		for _, msg := range msgs {
			writerCh.ReplicationTaskChannel <- msg
		}
	case kafkaMessageTypeVisibilityMsg:
		msgs, skippedDeserializeCount := deserializeVisibilityMessages(messages, skipErrors)
		atomic.AddInt32(skippedCount, skippedGetMsgCount+skippedDeserializeCount)
		for _, msg := range msgs {
			writerCh.VisibilityMsgChannel <- msg
		}
	}
}

func getMessages(data []byte, skipErrors bool) ([][]byte, int32) {
	str := string(data)
	messagesWithHeaders := r.Split(str, -1)
	if len(messagesWithHeaders[0]) != 0 {
		ErrorAndExit(malformedMessage, errors.New("got data chunk to handle that does not start with valid header"))
	}
	messagesWithHeaders = messagesWithHeaders[1:]
	var rawMessages [][]byte
	var skipped int32
	for _, m := range messagesWithHeaders {
		if len(m) == 0 {
			ErrorAndExit(malformedMessage, errors.New("got empty message between valid headers"))
		}
		curr := []byte(m)
		messageStart := bytes.Index(curr, []byte{preambleVersion0})
		if messageStart == -1 {
			if !skipErrors {
				ErrorAndExit(malformedMessage, errors.New("failed to find message preamble"))
			} else {
				skipped++
				continue
			}
		}
		rawMessages = append(rawMessages, curr[messageStart:])
	}
	return rawMessages, skipped
}

func deserializeMessages(messages [][]byte, skipErrors bool) ([]*types.ReplicationTask, int32) {
	var replicationTasks []*types.ReplicationTask
	var skipped int32
	for _, m := range messages {
		var task replicator.ReplicationTask
		err := decode(m, &task)
		if err != nil {
			if !skipErrors {
				ErrorAndExit(malformedMessage, err)
			} else {
				skipped++
				continue
			}
		}
		replicationTasks = append(replicationTasks, thrift.ToReplicationTask(&task))
	}
	return replicationTasks, skipped
}

func decode(message []byte, val *replicator.ReplicationTask) error {
	reader := bytes.NewReader(message[1:])
	wireVal, err := protocol.Binary.Decode(reader, wire.TStruct)
	if err != nil {
		return err
	}
	return val.FromWire(wireVal)
}

func deserializeVisibilityMessages(messages [][]byte, skipErrors bool) ([]*indexer.Message, int32) {
	var visibilityMessages []*indexer.Message
	var skipped int32
	for _, m := range messages {
		var msg indexer.Message
		err := decodeVisibility(m, &msg)
		if err != nil {
			if !skipErrors {
				ErrorAndExit(malformedMessage, err)
			} else {
				skipped++
				continue
			}
		}
		visibilityMessages = append(visibilityMessages, &msg)
	}
	return visibilityMessages, skipped
}

func decodeVisibility(message []byte, val *indexer.Message) error {
	reader := bytes.NewReader(message[1:])
	wireVal, err := protocol.Binary.Decode(reader, wire.TStruct)
	if err != nil {
		return err
	}
	return val.FromWire(wireVal)
}

// ClustersConfig describes the kafka clusters
type ClustersConfig struct {
	Clusters map[string]config.ClusterConfig
	TLS      config.TLS
}

func doRereplicate(
	ctx context.Context,
	domainID string,
	wid string,
	rid string,
	endEventID *int64,
	endEventVersion *int64,
	sourceCluster string,
	adminClient admin.Client,
) {
	fmt.Printf("Start rereplication for wid: %v, rid:%v \n", wid, rid)
	if err := adminClient.ResendReplicationTasks(
		ctx,
		&types.ResendReplicationTasksRequest{
			DomainID:      domainID,
			WorkflowID:    wid,
			RunID:         rid,
			RemoteCluster: sourceCluster,
			EndEventID:    endEventID,
			EndVersion:    endEventVersion,
		},
	); err != nil {
		ErrorAndExit("Failed to resend ndc workflow", err)
	}
	fmt.Printf("Done rereplication for wid: %v, rid:%v \n", wid, rid)
}

// AdminRereplicate parses will re-publish replication tasks to topic
func AdminRereplicate(c *cli.Context) {
	numberOfShards := c.Int(FlagNumberOfShards)
	if numberOfShards <= 0 {
		ErrorAndExit("numberOfShards is must be > 0", nil)
		return
	}
	sourceCluster := getRequiredOption(c, FlagSourceCluster)

	adminClient := cFactory.ServerAdminClient(c)
	var endEventID, endVersion *int64
	if c.IsSet(FlagMaxEventID) {
		endEventID = common.Int64Ptr(c.Int64(FlagMaxEventID) + 1)
	}
	if c.IsSet(FlagEndEventVersion) {
		endVersion = common.Int64Ptr(c.Int64(FlagEndEventVersion))
	}
	domainID := getRequiredOption(c, FlagDomainID)
	wid := getRequiredOption(c, FlagWorkflowID)
	rid := getRequiredOption(c, FlagRunID)
	contextTimeout := defaultResendContextTimeout

	if c.GlobalIsSet(FlagContextTimeout) {
		contextTimeout = time.Duration(c.GlobalInt(FlagContextTimeout)) * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	doRereplicate(
		ctx,
		domainID,
		wid,
		rid,
		endEventID,
		endVersion,
		sourceCluster,
		adminClient,
	)
}

func decodeReplicationTask(
	task *types.ReplicationTask,
	serializer persistence.PayloadSerializer,
) ([]byte, error) {

	switch task.GetTaskType() {
	case types.ReplicationTaskTypeHistoryV2:
		historyV2 := task.GetHistoryTaskV2Attributes()
		events, err := serializer.DeserializeBatchEvents(
			persistence.NewDataBlobFromInternal(historyV2.Events),
		)
		if err != nil {
			return nil, err
		}
		var newRunEvents []*types.HistoryEvent
		if historyV2.NewRunEvents != nil {
			newRunEvents, err = serializer.DeserializeBatchEvents(
				persistence.NewDataBlobFromInternal(historyV2.NewRunEvents),
			)
			if err != nil {
				return nil, err
			}
		}
		historyV2.Events = nil
		historyV2.NewRunEvents = nil
		historyV2Attributes := &historyV2Task{
			Task:         task,
			Events:       events,
			NewRunEvents: newRunEvents,
		}
		return json.Marshal(historyV2Attributes)
	default:
		return json.Marshal(task)
	}
}
