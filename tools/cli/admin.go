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

import "github.com/urfave/cli/v2"

func newAdminWorkflowCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "show",
			Aliases: []string{"show"},
			Usage:   "show workflow history from database",
			Flags: []cli.Flag{
				// v2 history events
				&cli.StringFlag{
					Name:  FlagTreeID,
					Usage: "TreeID",
				},
				&cli.StringFlag{
					Name:  FlagBranchID,
					Usage: "BranchID",
				},
				&cli.StringFlag{
					Name:    FlagOutputFilename,
					Aliases: FlagOutputFilenameAlias,
					Usage:   "output file",
				},

				// for persistence connection
				// TODO need to support other database: https://github.com/uber/cadence/issues/2777
				&cli.StringFlag{
					Name:  FlagDBAddress,
					Usage: "persistence address(right now only cassandra is supported)",
				},
				&cli.IntFlag{
					Name:  FlagDBPort,
					Value: 9042,
					Usage: "persistence port",
				},
				&cli.StringFlag{
					Name:  FlagUsername,
					Usage: "cassandra username",
				},
				&cli.StringFlag{
					Name:  FlagPassword,
					Usage: "cassandra password",
				},
				&cli.StringFlag{
					Name:  FlagKeyspace,
					Usage: "cassandra keyspace",
				},
				&cli.BoolFlag{
					Name:  FlagEnableTLS,
					Usage: "enable TLS over cassandra connection",
				},
				&cli.StringFlag{
					Name:  FlagTLSCertPath,
					Usage: "cassandra tls client cert path (tls must be enabled)",
				},
				&cli.StringFlag{
					Name:  FlagTLSKeyPath,
					Usage: "cassandra tls client key path (tls must be enabled)",
				},
				&cli.StringFlag{
					Name:  FlagTLSCaPath,
					Usage: "cassandra tls client ca path (tls must be enabled)",
				},
				&cli.BoolFlag{
					Name:  FlagTLSEnableHostVerification,
					Usage: "cassandra tls verify hostname and server cert (tls must be enabled)",
				},

				// support mysql query
				&cli.IntFlag{
					Name:    FlagShardID,
					Aliases: FlagShardIDAlias,
					Usage:   "ShardID",
				},
			},
			Action: func(c *cli.Context) error {
				AdminShowWorkflow(c)
				return nil
			},
		},
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "Describe internal information of workflow execution",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    FlagWorkflowID,
					Aliases: FlagWorkflowIDAlias,
					Usage:   "WorkflowID",
				},
				&cli.StringFlag{
					Name:    FlagRunID,
					Aliases: FlagRunIDAlias,
					Usage:   "RunID",
				},
			},
			Action: func(c *cli.Context) error {
				AdminDescribeWorkflow(c)
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"del"},
			Usage:   "Delete current workflow execution and the mutableState record",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    FlagWorkflowID,
					Aliases: FlagWorkflowIDAlias,
					Usage:   "WorkflowID",
				},
				&cli.StringFlag{
					Name:    FlagRunID,
					Aliases: FlagRunIDAlias,
					Usage:   "RunID",
				},
				&cli.BoolFlag{
					Name:    FlagSkipErrorMode,
					Aliases: FlagSkipErrorModeAlias,
					Usage:   "skip errors when deleting history",
				},

				// for persistence connection
				// TODO need to support other database: https://github.com/uber/cadence/issues/2777
				&cli.StringFlag{
					Name:  FlagDBAddress,
					Usage: "persistence address(right now only cassandra is supported)",
				},
				&cli.IntFlag{
					Name:  FlagDBPort,
					Value: 9042,
					Usage: "persistence port",
				},
				&cli.StringFlag{
					Name:  FlagUsername,
					Usage: "cassandra username",
				},
				&cli.StringFlag{
					Name:  FlagPassword,
					Usage: "cassandra password",
				},
				&cli.StringFlag{
					Name:  FlagKeyspace,
					Usage: "cassandra keyspace",
				},
				&cli.BoolFlag{
					Name:  FlagEnableTLS,
					Usage: "use TLS over cassandra connection",
				},
				&cli.StringFlag{
					Name:  FlagTLSCertPath,
					Usage: "cassandra tls client cert path (tls must be enabled)",
				},
				&cli.StringFlag{
					Name:  FlagTLSKeyPath,
					Usage: "cassandra tls client key path (tls must be enabled)",
				},
				&cli.StringFlag{
					Name:  FlagTLSCaPath,
					Usage: "cassandra tls client ca path (tls must be enabled)",
				},
				&cli.BoolFlag{
					Name:  FlagTLSEnableHostVerification,
					Usage: "cassandra tls verify hostname and server cert (tls must be enabled)",
				},
			},
			Action: func(c *cli.Context) error {
				AdminDeleteWorkflow(c)
				return nil
			},
		},
	}
}

func newAdminShardManagementCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "closeShard",
			Aliases: []string{"clsh"},
			Usage:   "close a shard given a shard id",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  FlagShardID,
					Usage: "ShardID for the cadence cluster to manage",
				},
			},
			Action: func(c *cli.Context) error {
				AdminShardManagement(c)
				return nil
			},
		},
		{
			Name:    "removeTask",
			Aliases: []string{"rmtk"},
			Usage:   "remove a task based on shardID, typeID and taskID",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  FlagShardID,
					Usage: "ShardID for the cadence cluster to manage",
				},
				&cli.Int64Flag{
					Name:  FlagRemoveTaskID,
					Usage: "task id which user want to specify",
				},
				&cli.IntFlag{
					Name:  FlagRemoveTypeID,
					Usage: "type id which user want to specify: 2 (transfer task), 3 (timer task), 4 (replication task)",
				},
			},
			Action: func(c *cli.Context) error {
				AdminRemoveTask(c)
				return nil
			},
		},
	}
}

func newAdminHistoryHostCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "Describe internal information of history host",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    FlagWorkflowID,
					Aliases: FlagWorkflowIDAlias,
					Usage:   "WorkflowID",
				},
				&cli.StringFlag{
					Name:    FlagHistoryAddress,
					Aliases: FlagHistoryAddressAlias,
					Usage:   "History Host address(IP:PORT)",
				},
				&cli.IntFlag{
					Name:    FlagShardID,
					Aliases: FlagShardIDAlias,
					Usage:   "ShardID",
				},
				&cli.BoolFlag{
					Name:    FlagPrintFullyDetail,
					Aliases: FlagPrintFullyDetailAlias,
					Usage:   "Print fully detail",
				},
			},
			Action: func(c *cli.Context) error {
				AdminDescribeHistoryHost(c)
				return nil
			},
		},
		{
			Name:    "getshard",
			Aliases: []string{"gsh"},
			Usage:   "Get shardID for a workflowID",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    FlagWorkflowID,
					Aliases: FlagWorkflowIDAlias,
					Usage:   "WorkflowID",
				},
				&cli.IntFlag{
					Name:  FlagNumberOfShards,
					Usage: "NumberOfShards for the cadence cluster(see config for numHistoryShards)",
				},
			},
			Action: func(c *cli.Context) error {
				AdminGetShardID(c)
				return nil
			},
		},
	}
}

func newAdminDomainCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "register",
			Aliases: []string{"re"},
			Usage:   "Register workflow domain",
			Flags:   adminRegisterDomainFlags,
			Action: func(c *cli.Context) error {
				newDomainCLI(c, true).RegisterDomain(c)
				return nil
			},
		},
		{
			Name:    "update",
			Aliases: []string{"up", "u"},
			Usage:   "Update existing workflow domain",
			Flags:   adminUpdateDomainFlags,
			Action: func(c *cli.Context) error {
				newDomainCLI(c, true).UpdateDomain(c)
				return nil
			},
		},
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "Describe existing workflow domain",
			Flags:   adminDescribeDomainFlags,
			Action: func(c *cli.Context) error {
				newDomainCLI(c, true).DescribeDomain(c)
				return nil
			},
		},
		{
			Name:    "getdomainidorname",
			Aliases: []string{"getdn"},
			Usage:   "Get domainID or domainName",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    FlagDomain,
					Aliases: FlagDomainAlias,
					Usage:   "DomainName",
				},
				&cli.StringFlag{
					Name:  FlagDomainID,
					Usage: "Domain ID(uuid)",
				},

				// for persistence connection
				// TODO need to support other database: https://github.com/uber/cadence/issues/2777
				&cli.StringFlag{
					Name:  FlagDBAddress,
					Usage: "persistence address(right now only cassandra is supported)",
				},
				&cli.IntFlag{
					Name:  FlagDBPort,
					Value: 9042,
					Usage: "persistence port",
				},
				&cli.StringFlag{
					Name:  FlagUsername,
					Usage: "cassandra username",
				},
				&cli.StringFlag{
					Name:  FlagPassword,
					Usage: "cassandra password",
				},
				&cli.StringFlag{
					Name:  FlagKeyspace,
					Usage: "cassandra keyspace",
				},
				&cli.BoolFlag{
					Name:  FlagEnableTLS,
					Usage: "use TLS over cassandra connection",
				},
				&cli.StringFlag{
					Name:  FlagTLSCertPath,
					Usage: "cassandra tls client cert path (tls must be enabled)",
				},
				&cli.StringFlag{
					Name:  FlagTLSKeyPath,
					Usage: "cassandra tls client key path (tls must be enabled)",
				},
				&cli.StringFlag{
					Name:  FlagTLSCaPath,
					Usage: "cassandra tls client ca path (tls must be enabled)",
				},
				&cli.BoolFlag{
					Name:  FlagTLSEnableHostVerification,
					Usage: "cassandra tls verify hostname and server cert (tls must be enabled)",
				},
			},
			Action: func(c *cli.Context) error {
				AdminGetDomainIDOrName(c)
				return nil
			},
		},
	}
}

func newAdminKafkaCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "parse",
			Aliases: []string{"par"},
			Usage:   "Parse replication tasks from kafka messages",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    FlagInputFile,
					Aliases: FlagInputFileAlias,
					Usage:   "Input file to use, if not present assumes piping",
				},
				&cli.StringFlag{
					Name:    FlagWorkflowID,
					Aliases: FlagWorkflowIDAlias,
					Usage:   "WorkflowID, if not provided then no filters by WorkflowID are applied",
				},
				&cli.StringFlag{
					Name:    FlagRunID,
					Aliases: FlagRunIDAlias,
					Usage:   "RunID, if not provided then no filters by RunID are applied",
				},
				&cli.StringFlag{
					Name:    FlagOutputFilename,
					Aliases: FlagOutputFilenameAlias,
					Usage:   "Output file to write to, if not provided output is written to stdout",
				},
				&cli.BoolFlag{
					Name:    FlagSkipErrorMode,
					Aliases: FlagSkipErrorModeAlias,
					Usage:   "Skip errors in parsing messages",
				},
				&cli.BoolFlag{
					Name:    FlagHeadersMode,
					Aliases: FlagHeadersModeAlias,
					Usage:   "Output headers of messages in format: DomainID, WorkflowID, RunID, FirstEventID, NextEventID",
				},
				&cli.IntFlag{
					Name:    FlagMessageType,
					Aliases: FlagMessageTypeAlias,
					Usage:   "Kafka message type (0: replicationTasks; 1: visibility)",
					Value:   0,
				},
			},
			Action: func(c *cli.Context) error {
				AdminKafkaParse(c)
				return nil
			},
		},
		{
			Name:    "purgeTopic",
			Aliases: []string{"purge"},
			Usage:   "purge Kafka topic by consumer group",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  FlagCluster,
					Usage: "Name of the Kafka cluster to publish replicationTasks",
				},
				&cli.StringFlag{
					Name:  FlagTopic,
					Usage: "Topic to publish replication task",
				},
				&cli.StringFlag{
					Name:  FlagGroup,
					Usage: "Group to read DLQ",
				},
				&cli.StringFlag{
					Name: FlagHostFile,
					Usage: "Kafka host config file in format of: " + `
tls:
    enabled: false
    certFile: ""
    keyFile: ""
    caFile: ""
clusters:
	localKafka:
		brokers:
		- 127.0.0.1
		- 127.0.0.2`,
				},
			},
			Action: func(c *cli.Context) error {
				AdminPurgeTopic(c)
				return nil
			},
		},
		{
			Name:    "mergeDLQ",
			Aliases: []string{"mgdlq"},
			Usage:   "Merge replication tasks to target topic(from input file or DLQ topic)",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    FlagInputFile,
					Aliases: FlagInputFileAlias,
					Usage:   "Input file to use to read as JSON of ReplicationTask, separated by line",
				},
				&cli.StringFlag{
					Name:    FlagInputTopic,
					Aliases: FlagInputTopicAlias,
					Usage:   "Input topic to read ReplicationTask",
				},
				&cli.StringFlag{
					Name:  FlagInputCluster,
					Usage: "Name of the Kafka cluster for reading DLQ topic for ReplicationTask",
				},
				&cli.Int64Flag{
					Name:  FlagStartOffset,
					Usage: "Starting offset for reading DLQ topic for ReplicationTask",
				},
				&cli.StringFlag{
					Name:  FlagCluster,
					Usage: "Name of the Kafka cluster to publish replicationTasks",
				},
				&cli.StringFlag{
					Name:  FlagTopic,
					Usage: "Topic to publish replication task",
				},
				&cli.StringFlag{
					Name:  FlagGroup,
					Usage: "Group to read DLQ",
				},
				&cli.StringFlag{
					Name: FlagHostFile,
					Usage: "Kafka host config file in format of: " + `
tls:
    enabled: false
    certFile: ""
    keyFile: ""
    caFile: ""
clusters:
	localKafka:
		brokers:
		- 127.0.0.1
		- 127.0.0.2`,
				},
			},
			Action: func(c *cli.Context) error {
				AdminMergeDLQ(c)
				return nil
			},
		},
		{
			Name:    "rereplicate",
			Aliases: []string{"rrp"},
			Usage:   "Rereplicate replication tasks to target topic from history tables",
			Flags: []cli.Flag{

				&cli.StringFlag{
					Name:  FlagTargetCluster,
					Usage: "Name of targetCluster to receive the replication task",
				},
				&cli.IntFlag{
					Name:  FlagNumberOfShards,
					Usage: "NumberOfShards is required to calculate shardID. (see server config for numHistoryShards)",
				},

				// for multiple workflow
				&cli.StringFlag{
					Name:    FlagInputFile,
					Aliases: FlagInputFileAlias,
					Usage:   "Input file to read multiple workflow line by line. For each line: domainID,workflowID,runID,minEventID,maxEventID (minEventID/maxEventID are optional.)",
				},

				// for one workflow
				&cli.Int64Flag{
					Name:  FlagMinEventID,
					Usage: "MinEventID. Optional, default to all events",
				},
				&cli.Int64Flag{
					Name:  FlagMaxEventID,
					Usage: "MaxEventID Optional, default to all events",
				},
				&cli.StringFlag{
					Name:    FlagWorkflowID,
					Aliases: FlagWorkflowIDAlias,
					Usage:   "WorkflowID",
				},
				&cli.StringFlag{
					Name:    FlagRunID,
					Aliases: FlagRunIDAlias,
					Usage:   "RunID",
				},
				&cli.StringFlag{
					Name:  FlagDomainID,
					Usage: "DomainID",
				},

				// for persistence connection
				// TODO need to support other database: https://github.com/uber/cadence/issues/2777
				&cli.StringFlag{
					Name:  FlagDBAddress,
					Usage: "persistence address(right now only cassandra is supported)",
				},
				&cli.IntFlag{
					Name:  FlagDBPort,
					Value: 9042,
					Usage: "persistence port",
				},
				&cli.StringFlag{
					Name:  FlagUsername,
					Usage: "cassandra username",
				},
				&cli.StringFlag{
					Name:  FlagPassword,
					Usage: "cassandra password",
				},
				&cli.StringFlag{
					Name:  FlagKeyspace,
					Usage: "cassandra keyspace",
				},
				&cli.BoolFlag{
					Name:  FlagEnableTLS,
					Usage: "use TLS over cassandra connection",
				},
				&cli.StringFlag{
					Name:  FlagTLSCertPath,
					Usage: "cassandra tls client cert path (tls must be enabled)",
				},
				&cli.StringFlag{
					Name:  FlagTLSKeyPath,
					Usage: "cassandra tls client key path (tls must be enabled)",
				},
				&cli.StringFlag{
					Name:  FlagTLSCaPath,
					Usage: "cassandra tls client ca path (tls must be enabled)",
				},
				&cli.BoolFlag{
					Name:  FlagTLSEnableHostVerification,
					Usage: "cassandra tls verify hostname and server cert (tls must be enabled)",
				},

				// kafka
				&cli.StringFlag{
					Name:  FlagCluster,
					Usage: "Name of the Kafka cluster to publish replicationTasks",
				},
				&cli.StringFlag{
					Name:  FlagTopic,
					Usage: "Topic to publish replication task",
				},
				&cli.StringFlag{
					Name: FlagHostFile,
					Usage: "Kafka host config file in format of: " + `
tls:
    enabled: false
    certFile: ""
    keyFile: ""
    caFile: ""
clusters:
	localKafka:
		brokers:
		- 127.0.0.1
		- 127.0.0.2`,
				},
			},
			Action: func(c *cli.Context) error {
				AdminRereplicate(c)
				return nil
			},
		},
	}
}

func newAdminElasticSearchCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "catIndex",
			Aliases: []string{"cind"},
			Usage:   "Cat Indices on ElasticSearch",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  FlagURL,
					Usage: "URL of ElasticSearch cluster",
				},
				&cli.StringFlag{
					Name:    FlagMuttleyDestination,
					Aliases: FlagMuttleyDestinationAlias,
					Usage:   "Optional muttely destination to ElasticSearch cluster",
				},
			},
			Action: func(c *cli.Context) error {
				AdminCatIndices(c)
				return nil
			},
		},
		{
			Name:    "index",
			Aliases: []string{"ind"},
			Usage:   "Index docs on ElasticSearch",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  FlagURL,
					Usage: "URL of ElasticSearch cluster",
				},
				&cli.StringFlag{
					Name:    FlagMuttleyDestination,
					Aliases: FlagMuttleyDestinationAlias,
					Usage:   "Optional muttely destination to ElasticSearch cluster",
				},
				&cli.StringFlag{
					Name:  FlagIndex,
					Usage: "ElasticSearch target index",
				},
				&cli.StringFlag{
					Name:    FlagInputFile,
					Aliases: FlagInputFileAlias,
					Usage:   "Input file of indexer.Message in json format, separated by newline",
				},
				&cli.IntFlag{
					Name:    FlagBatchSize,
					Aliases: FlagBatchSizeAlias,
					Usage:   "Optional batch size of actions for bulk operations",
					Value:   1000,
				},
			},
			Action: func(c *cli.Context) error {
				AdminIndex(c)
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"del"},
			Usage:   "Delete docs on ElasticSearch",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  FlagURL,
					Usage: "URL of ElasticSearch cluster",
				},
				&cli.StringFlag{
					Name:    FlagMuttleyDestination,
					Aliases: FlagMuttleyDestinationAlias,
					Usage:   "Optional muttely destination to ElasticSearch cluster",
				},
				&cli.StringFlag{
					Name:  FlagIndex,
					Usage: "ElasticSearch target index",
				},
				&cli.StringFlag{
					Name:    FlagInputFile,
					Aliases: FlagInputFileAlias,
					Usage: "Input file name. Redirect cadence wf list result (with tale format) to a file and use as delete input. " +
						"First line should be table header like WORKFLOW TYPE | WORKFLOW ID | RUN ID | ...",
				},
				&cli.IntFlag{
					Name:    FlagBatchSize,
					Aliases: FlagBatchSizeAlias,
					Usage:   "Optional batch size of actions for bulk operations",
					Value:   1000,
				},
				&cli.IntFlag{
					Name:  FlagRPS,
					Usage: "Optional batch request rate per second",
					Value: 30,
				},
			},
			Action: func(c *cli.Context) error {
				AdminDelete(c)
				return nil
			},
		},
		{
			Name:    "report",
			Aliases: []string{"rep"},
			Usage:   "Generate Report by Aggregation functions on ElasticSearch",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  FlagURL,
					Usage: "URL of ElasticSearch cluster",
				},
				&cli.StringFlag{
					Name:  FlagIndex,
					Usage: "ElasticSearch target index",
				},
				&cli.StringFlag{
					Name:  FlagListQuery,
					Usage: "SQL query of the report",
				},
				&cli.StringFlag{
					Name:  FlagOutputFormat,
					Usage: "Additional output format (html or csv)",
				},
				&cli.StringFlag{
					Name:  FlagOutputFilename,
					Usage: "Additional output filename with path",
				},
			},
			Action: func(c *cli.Context) error {
				GenerateReport(c)
				return nil
			},
		},
	}
}

func newAdminTaskListCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "describe",
			Aliases: []string{"desc"},
			Usage:   "Describe pollers and status information of tasklist",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    FlagTaskList,
					Aliases: FlagTaskListAlias,
					Usage:   "TaskList description",
				},
				&cli.StringFlag{
					Name:    FlagTaskListType,
					Aliases: FlagTaskListTypeAlias,
					Value:   "decision",
					Usage:   "Optional TaskList type [decision|activity]",
				},
			},
			Action: func(c *cli.Context) error {
				AdminDescribeTaskList(c)
				return nil
			},
		},
	}
}

func newAdminClusterCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "add-search-attr",
			Aliases: []string{"asa"},
			Usage:   "whitelist search attribute",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  FlagSearchAttributesKey,
					Usage: "Search Attribute key to be whitelisted",
				},
				&cli.IntFlag{
					Name:  FlagSearchAttributesType,
					Value: -1,
					Usage: "Search Attribute value type. [0:String, 1:Keyword, 2:Int, 3:Double, 4:Bool, 5:Datetime]",
				},
				&cli.StringFlag{
					Name:    FlagSecurityToken,
					Aliases: FlagSecurityTokenAlias,
					Usage:   "Optional token for security check",
				},
			},
			Action: func(c *cli.Context) error {
				AdminAddSearchAttribute(c)
				return nil
			},
		},
	}
}
