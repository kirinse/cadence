// Copyright (c) 2017-2021 Uber Technologies Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package basic

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/cadence"
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"

	"github.com/uber/cadence/bench/lib"
	"github.com/uber/cadence/bench/load/common"
	c "github.com/uber/cadence/common"
)

const (
	// TestName is the test name for basic test
	TestName = "basic"

	// LauncherWorkflowName is the workflow name for launching basic load test
	LauncherWorkflowName = "basic-load-test-workflow"
	// LauncherLaunchActivityName is the activity name for launching stress load test
	LauncherLaunchActivityName = "basic-load-test-launch-activity"
	// LauncherVerifyActivityName is the verification activity name
	LauncherVerifyActivityName = "basic-load-test-verify-activity"

	defaultStressWorkflowStartToCloseTimeout = 5 * time.Minute  // default 5m may not be not enough for long running workflow
	workflowWaitTimeBuffer                   = 10 * time.Second // time buffer of waiting for stressWorkflow execution. Hence the actual waiting time is  stressWorkflowStartToCloseTimeout + workflowWaitTimeBuffer
	maxLauncherActivityRetryCount            = 5                // number of retry for launcher activity
)

type (
	launcherActivityParams struct {
		RoutineID int                 // the ID of the launchActivity
		Count     int                 // stressWorkflows to start per launchActivity
		Config    lib.BasicTestConfig // config of this load test
	}

	verifyActivityParams struct {
		WorkflowStartTime    int64
		ListWorkflowPageSize int32
	}

	verifyActivityResult struct {
		OpenCount    int
		TimeoutCount int
		FailedCount  int
	}
)

// RegisterLauncher registers workflows and activities for basic load launching
func RegisterLauncher(w worker.Worker) {
	w.RegisterWorkflowWithOptions(launcherWorkflow, workflow.RegisterOptions{Name: LauncherWorkflowName})
	w.RegisterActivityWithOptions(launcherActivity, activity.RegisterOptions{Name: LauncherLaunchActivityName})
	w.RegisterActivityWithOptions(verifyResultActivity, activity.RegisterOptions{Name: LauncherVerifyActivityName})
}

func launcherWorkflow(ctx workflow.Context, config lib.BasicTestConfig) (string, error) {
	logger := workflow.GetLogger(ctx).With(zap.String("Test", TestName))
	workflowPerActivity := config.TotalLaunchCount / config.RoutineCount
	workflowTimeout := workflow.GetInfo(ctx).ExecutionStartToCloseTimeoutSeconds
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Duration(workflowTimeout) * time.Second,
		HeartbeatTimeout:       20 * time.Second,
		RetryPolicy: &workflow.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 1,
			MaximumInterval:    time.Second,
			ExpirationInterval: time.Second * time.Duration(workflowTimeout+1) * time.Second * maxLauncherActivityRetryCount,
			MaximumAttempts:    maxLauncherActivityRetryCount,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	startTime := workflow.Now(ctx)

	for i := 0; i < config.RoutineCount; i++ {
		actInput := launcherActivityParams{
			RoutineID: i,
			Count:     workflowPerActivity,
			Config:    config,
		}
		f := workflow.ExecuteActivity(
			ctx,
			LauncherLaunchActivityName,
			actInput,
		)
		err := f.Get(ctx, nil)
		if err != nil {
			return "", err
		}
	}

	workflowWaitTime := defaultStressWorkflowStartToCloseTimeout
	if config.ExecutionStartToCloseTimeoutInSeconds > 0 {
		workflowWaitTime = time.Duration(config.ExecutionStartToCloseTimeoutInSeconds)*time.Second + workflowWaitTimeBuffer
	}
	logger.Info(fmt.Sprintf("%v stressWorkflows are launched, now waiting for %v ...", config.TotalLaunchCount, workflowWaitTime))
	if err := workflow.Sleep(ctx, workflowWaitTime); err != nil {
		return "", fmt.Errorf("launcher workflow sleep failed: %v", err)
	}

	maxTolerantFailure := int32(float64(config.TotalLaunchCount) * config.FailureThreshold)
	actInput := verifyActivityParams{
		WorkflowStartTime:    startTime.UnixNano(),
		ListWorkflowPageSize: maxTolerantFailure + 10,
	}
	actOptions := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		RetryPolicy: &cadence.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2,
			MaximumAttempts:    5,
		},
	}
	var result verifyActivityResult
	ctx = workflow.WithActivityOptions(ctx, actOptions)
	if err := workflow.ExecuteActivity(
		ctx,
		LauncherVerifyActivityName,
		actInput,
	).Get(ctx, &result); err != nil {
		return "", err
	}
	passed := (result.TimeoutCount + result.OpenCount + result.FailedCount) <= int(maxTolerantFailure)
	finalResult := fmt.Sprintf("TEST PASSED: %v; \nDetails report: timeoutCount: %v, failedCount: %v, openCount:%v, launchCount: %v, maxThreshold:%v",
		passed, result.TimeoutCount, result.FailedCount, result.OpenCount, config.TotalLaunchCount, maxTolerantFailure)
	return finalResult, nil
}

func launcherActivity(ctx context.Context, params launcherActivityParams) error {
	logger := activity.GetLogger(ctx).With(zap.String("Test", TestName))

	var lastStartedID int
	if activity.HasHeartbeatDetails(ctx) {
		err := activity.GetHeartbeatDetails(ctx, &lastStartedID)
		if err != nil {
			logger.Error("failed to resume from last  checkpoint...start from beginning...")
		}
		logger.Info("resume from last checkpoint", zap.Int("checkpoint", lastStartedID))
	}

	logger.Info("Start activity routineID", zap.Int("RoutineID", params.RoutineID))

	cc := ctx.Value(lib.CtxKeyCadenceClient).(lib.CadenceClient)
	rc := ctx.Value(lib.CtxKeyRuntimeContext).(*lib.RuntimeContext)
	numTaskList := rc.Bench.NumTaskLists
	basicTestConfig := params.Config

	stressWorkflowTimeout := defaultStressWorkflowStartToCloseTimeout
	if basicTestConfig.ExecutionStartToCloseTimeoutInSeconds > 0 {
		stressWorkflowTimeout = time.Duration(basicTestConfig.ExecutionStartToCloseTimeoutInSeconds) * time.Second
	}
	workflowOptions := client.StartWorkflowOptions{
		ExecutionStartToCloseTimeout: stressWorkflowTimeout,
	}

	var sleepTime time.Duration
	if basicTestConfig.MinCadenceSleepInSeconds > 0 && basicTestConfig.MaxCadenceSleepInSeconds > 0 {
		minSleep := time.Duration(basicTestConfig.MinCadenceSleepInSeconds) * time.Second
		maxSleep := time.Duration(basicTestConfig.MaxCadenceSleepInSeconds) * time.Second
		jitter := rand.Float64() * float64(maxSleep-minSleep)
		sleepTime = minSleep + time.Duration(int64(jitter))
	}
	stressWorkflowInput := WorkflowParams{
		ChainSequence:    basicTestConfig.ChainSequence,
		ConcurrentCount:  basicTestConfig.ConcurrentCount,
		PayloadSizeBytes: basicTestConfig.PayloadSizeBytes,
		PanicWorkflow:    basicTestConfig.PanicStressWorkflow,
		CadenceSleep:     sleepTime,
	}

	startWFCtxTimeout := common.DefaultContextTimeout
	if basicTestConfig.ContextTimeoutInSeconds > 0 {
		startWFCtxTimeout = time.Duration(basicTestConfig.ContextTimeoutInSeconds) * time.Second
	}

	for startedID := lastStartedID; startedID < params.Count; startedID++ {
		stressWorkflowInput.TaskListNumber = rand.Intn(numTaskList)

		workflowOptions.ID = fmt.Sprintf("%d-%d", params.RoutineID, startedID)
		workflowOptions.TaskList = common.GetTaskListName(stressWorkflowInput.TaskListNumber)

		startWorkflowContext, cancelF := context.WithTimeout(context.Background(), startWFCtxTimeout)
		we, err := cc.StartWorkflow(startWorkflowContext, workflowOptions, stressWorkflowExecute, stressWorkflowInput)
		cancelF()
		if err == nil {
			logger.Debug("Created Workflow successfully", zap.String("WorkflowID", we.ID), zap.String("RunID", we.RunID))
		} else {
			if _, ok := err.(*shared.WorkflowExecutionAlreadyStartedError); ok {
				logger.Debug("Workflow already started in previous activity attempt")
			}
			logger.Error("Failed to start workflow execution", zap.Error(err))
			return err
		}
		activity.RecordHeartbeat(ctx, startedID)
		jitter := time.Duration(75 + rand.Intn(25))
		time.Sleep(jitter * time.Millisecond)
	}

	logger.Info("finish running launcher activity", zap.Int("StartedCount", params.Count))
	return nil
}

func verifyResultActivity(
	ctx context.Context,
	params verifyActivityParams,
) (verifyActivityResult, error) {
	cc := ctx.Value(lib.CtxKeyCadenceClient).(lib.CadenceClient)
	info := activity.GetInfo(ctx)
	var opens, timeouts, faileds int

	// verify if any open workflow
	listOpenWorkflowRequest := &shared.ListOpenWorkflowExecutionsRequest{
		Domain:          c.StringPtr(info.WorkflowDomain),
		MaximumPageSize: &params.ListWorkflowPageSize,
		StartTimeFilter: &shared.StartTimeFilter{
			EarliestTime: c.Int64Ptr(params.WorkflowStartTime),
			LatestTime:   c.Int64Ptr(time.Now().UnixNano()),
		},
		TypeFilter: &shared.WorkflowTypeFilter{
			Name: c.StringPtr(stressWorkflowName),
		},
	}
	openWfs, err := cc.ListOpenWorkflow(ctx, listOpenWorkflowRequest)
	if err != nil {
		return verifyActivityResult{}, err
	}

	if len(openWfs.Executions) > 0 {
		opens = len(openWfs.Executions)
	}

	// verify if any failed workflow
	closeStatus := shared.WorkflowExecutionCloseStatusFailed
	listWorkflowRequest := &shared.ListClosedWorkflowExecutionsRequest{
		Domain:          c.StringPtr(info.WorkflowDomain),
		MaximumPageSize: &params.ListWorkflowPageSize,
		StartTimeFilter: &shared.StartTimeFilter{
			EarliestTime: c.Int64Ptr(params.WorkflowStartTime),
			LatestTime:   c.Int64Ptr(time.Now().UnixNano()),
		},
		TypeFilter: &shared.WorkflowTypeFilter{
			Name: c.StringPtr(stressWorkflowName),
		},
		StatusFilter: &closeStatus,
	}
	wfs, err := cc.ListClosedWorkflow(ctx, listWorkflowRequest)
	if err != nil {
		return verifyActivityResult{}, err
	}

	if len(wfs.Executions) > 0 {
		faileds = len(wfs.Executions)
	}

	// verify if any timeouted workflow
	closeStatus = shared.WorkflowExecutionCloseStatusTimedOut
	listWorkflowRequest.StatusFilter = &closeStatus
	wfs, err = cc.ListClosedWorkflow(ctx, listWorkflowRequest)
	if err != nil {
		return verifyActivityResult{}, err
	}

	if len(wfs.Executions) > 0 {
		timeouts = len(wfs.Executions)
	}

	return verifyActivityResult{
		OpenCount:    opens,
		FailedCount:  faileds,
		TimeoutCount: timeouts,
	}, nil
}
