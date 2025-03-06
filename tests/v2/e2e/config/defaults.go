//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package config provides configuration types and logic for loading and binding configuration values.
// This file includes detailed Bind methods for all configuration types with extensive comments.
package config

import (
	"os"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/timeutil"
)

// Constant definitions for default host and port values.
const (
	localhost = "localhost"
	localPort = Port("8081")

	defaultNum                  uint64                  = 10000
	defaultOffset               uint64                  = 0
	defaultTimestamp            int64                   = 0
	defaultSkipStrictExistCheck bool                    = false
	defaultConcurrency          uint64                  = 10
	defaultTimeout              timeutil.DurationString = "3s"
	defaultWaitAfterInsert      timeutil.DurationString = "2m"
)

// Default holds the default configuration values.
// It is used to provide fallback values and defaults for the Bind process.
var Default = &Data{
	Target: &config.GRPCClient{
		Addrs: []string{net.JoinHostPort(localhost, localPort.Port())},
	},
	Strategies: []*Strategy{
		{
			Name:        "check Index Property",
			Concurrency: 1,
			Operations: []*Operation{
				{
					Name: "IndexProperty",
					Executions: []*Execution{
						{
							Name: "IndexProperty",
							Type: OpIndexProperty,
							Mode: OperationUnary,
						},
					},
				},
			},
		},
		{
			Name:        "Initial Insert and Wait",
			Concurrency: 1,
			Operations: []*Operation{
				{
					Name: "Insert -> IndexInfo",
					Executions: []*Execution{
						{
							TimeConfig: TimeConfig{
								Timeout: "",
								Delay:   "",
								Wait:    defaultWaitAfterInsert,
							},
							Name: "Insert",
							Type: OpInsert,
							Mode: OperationUnary,
							BaseConfig: &BaseConfig{
								Num:         defaultNum,
								Offset:      defaultOffset,
								Concurrency: defaultConcurrency,
							},
							ModificationConfig: &ModificationConfig{
								Timestamp:            defaultTimestamp,
								SkipStrictExistCheck: defaultSkipStrictExistCheck,
							},
						},
						{
							Name: "IndexInfo",
							Type: OpIndexInfo,
							Mode: OperationUnary,
						},
					},
				},
			},
		},
		{
			Concurrency: 4,
			Operations: []*Operation{
				{
					Executions: []*Execution{
						{
							Type: OpSearch,
							Mode: OperationUnary,
							BaseConfig: &BaseConfig{
								Num:         defaultNum,
								Offset:      defaultOffset,
								Concurrency: defaultConcurrency,
							},
							SearchConfig: []*SearchQuery{
								{
									Timeout: defaultTimeout,
								},
							},
						},
					},
				},
				{
					Executions: []*Execution{
						{
							Type: OpSearchByID,
							Mode: OperationUnary,
							BaseConfig: &BaseConfig{
								Num:         defaultNum,
								Offset:      defaultOffset,
								Concurrency: defaultConcurrency,
							},
							SearchConfig: []*SearchQuery{
								{
									Timeout: defaultTimeout,
								},
							},
						},
					},
				},
				{
					Executions: []*Execution{
						{
							Type: OpLinearSearch,
							Mode: OperationUnary,
							BaseConfig: &BaseConfig{
								Num:         defaultNum,
								Offset:      defaultOffset,
								Concurrency: defaultConcurrency,
							},
							SearchConfig: []*SearchQuery{
								{
									Timeout: defaultTimeout,
								},
							},
						},
					},
				},
				{
					Executions: []*Execution{
						{
							Type: OpLinearSearchByID,
							Mode: OperationUnary,
							BaseConfig: &BaseConfig{
								Num:         defaultNum,
								Offset:      defaultOffset,
								Concurrency: defaultConcurrency,
							},
							SearchConfig: []*SearchQuery{
								{
									Timeout: defaultTimeout,
								},
							},
						},
					},
				},
			},
		},
		{
			TimeConfig: TimeConfig{
				Timeout: "",
				Delay:   "",
				Wait:    defaultWaitAfterInsert,
			},
			Concurrency: 3,
			Operations: []*Operation{
				{
					Executions: []*Execution{
						{
							Type: OpObject,
							Mode: OperationUnary,
						},
					},
				},
				{
					Executions: []*Execution{
						{
							Type: OpExists,
							Mode: OperationUnary,
						},
					},
				},
				{
					Executions: []*Execution{
						{
							Type: OpTimestamp,
							Mode: OperationUnary,
						},
					},
				},
			},
		},
		{
			TimeConfig: TimeConfig{
				Timeout: "",
				Delay:   "",
				Wait:    defaultWaitAfterInsert,
			},
			Concurrency: 1,
			Operations: []*Operation{
				{
					Executions: []*Execution{
						{
							TimeConfig: TimeConfig{
								Timeout: "",
								Delay:   "",
								Wait:    defaultWaitAfterInsert,
							},
							Type: OpUpdate,
							Mode: OperationUnary,
							BaseConfig: &BaseConfig{
								Num:         defaultNum,
								Offset:      defaultOffset,
								Concurrency: defaultConcurrency,
							},
							ModificationConfig: &ModificationConfig{
								Timestamp:            defaultTimestamp,
								SkipStrictExistCheck: defaultSkipStrictExistCheck,
							},
						},
						{
							Type: OpIndexDetail,
							Mode: OperationUnary,
						},
					},
				},
			},
		},
		{
			TimeConfig: TimeConfig{
				Timeout: "",
				Delay:   "",
				Wait:    defaultWaitAfterInsert,
			},
			Concurrency: 2,
			Operations: []*Operation{
				{
					Executions: []*Execution{
						{
							Type: OpRemove,
							Mode: OperationUnary,
							BaseConfig: &BaseConfig{
								Num:         defaultNum,
								Offset:      defaultOffset,
								Concurrency: defaultConcurrency,
							},
							ModificationConfig: &ModificationConfig{
								Timestamp:            defaultTimestamp,
								SkipStrictExistCheck: defaultSkipStrictExistCheck,
							},
						},
						{
							Type: OpIndexStatistics,
							Mode: OperationUnary,
						},
					},
				},
				{
					Executions: []*Execution{
						{
							Type: OpUpsert,
							Mode: OperationUnary,
							BaseConfig: &BaseConfig{
								Num:         defaultNum,
								Offset:      defaultOffset,
								Concurrency: defaultConcurrency,
							},
							ModificationConfig: &ModificationConfig{
								Timestamp:            defaultTimestamp,
								SkipStrictExistCheck: defaultSkipStrictExistCheck,
							},
						},
						{
							Type: OpIndexDetail,
							Mode: OperationUnary,
						},
					},
				},
			},
		},
		{
			TimeConfig: TimeConfig{
				Timeout: "",
				Delay:   "",
				Wait:    defaultWaitAfterInsert,
			},
			Concurrency: 1,
			Operations: []*Operation{
				{
					Executions: []*Execution{
						{
							TimeConfig: TimeConfig{
								Timeout: "",
								Delay:   "",
								Wait:    defaultWaitAfterInsert,
							},
							Type: OpRemoveByTimestamp,
							Mode: OperationUnary,
							ModificationConfig: &ModificationConfig{
								Timestamp: defaultTimestamp,
							},
						},
						{
							Type: OpIndexDetail,
							Mode: OperationUnary,
						},
						{
							TimeConfig: TimeConfig{
								Timeout: "",
								Delay:   "",
								Wait:    defaultWaitAfterInsert,
							},
							Type: OpUpsert,
							Mode: OperationUnary,
							BaseConfig: &BaseConfig{
								Num:         defaultNum,
								Offset:      defaultOffset,
								Concurrency: defaultConcurrency,
							},
							ModificationConfig: &ModificationConfig{
								Timestamp:            defaultTimestamp,
								SkipStrictExistCheck: defaultSkipStrictExistCheck,
							},
						},
						{
							Type: OpIndexDetail,
							Mode: OperationUnary,
						},
					},
				},
			},
		},
		{
			Concurrency: 1,
			Operations: []*Operation{
				{
					Executions: []*Execution{
						{
							Type: OpIndexStatisticsDetail,
							Mode: OperationUnary,
						},
						{
							TimeConfig: TimeConfig{
								Timeout: "",
								Delay:   "",
								Wait:    defaultWaitAfterInsert,
							},
							Type: OpFlush,
							Mode: OperationUnary,
						},
						{
							Type: OpIndexInfo,
							Mode: OperationUnary,
						},
					},
				},
			},
		},
	},
	Dataset: &Dataset{
		Name: "fashion-mnist-784-euclidean.hdf5",
	},
	Kubernetes: &Kubernetes{
		KubeConfig: file.Join(os.Getenv("HOME"), ".kube", "config"),
		PortForward: &PortForward{
			Enabled:    false,
			PodName:    "vald-gateway-0",
			TargetPort: localPort,
			LocalPort:  localPort,
			Namespace:  "default",
		},
	},
	Metadata: map[string]string{
		"sample metadata key1": "sample metadata value1",
		"sample metadata key2": "sample metadata value2",
		"sample metadata key3": "sample metadata value3",
	},
}
