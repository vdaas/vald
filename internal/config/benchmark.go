//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package config providers configuration type and load configuration logic
package config

import "github.com/vdaas/vald/internal/k8s/client"

// BenchmarkJob represents the configuration for the internal benchmark search job.
type BenchmarkJob struct {
	Target             *BenchmarkTarget    `json:"target,omitempty" yaml:"target"`
	Dataset            *BenchmarkDataset   `json:"dataset,omitempty" yaml:"dataset"`
	Dimension          int                 `json:"dimension,omitempty" yaml:"dimension"`
	Replica            int                 `json:"replica,omitempty" yaml:"replica"`
	Repetition         int                 `json:"repetition,omitempty" yaml:"repetition"`
	JobType            string              `json:"job_type,omitempty" yaml:"job_type"`
	InsertConfig       *InsertConfig       `json:"insert_config,omitempty" yaml:"insert_config"`
	UpdateConfig       *UpdateConfig       `json:"update_config,omitempty" yaml:"update_config"`
	UpsertConfig       *UpsertConfig       `json:"upsert_config,omitempty" yaml:"upsert_config"`
	SearchConfig       *SearchConfig       `json:"search_config,omitempty" yaml:"search_config"`
	RemoveConfig       *RemoveConfig       `json:"remove_config,omitempty" yaml:"remove_config"`
	ClientConfig       *GRPCClient         `json:"client_config,omitempty" yaml:"client_config"`
	Rules              []*BenchmarkJobRule `json:"rules,omitempty" yaml:"rules"`
	BeforeJobName      string              `json:"before_job_name,omitempty" yaml:"before_job_name"`
	BeforeJobNamespace string              `json:"before_job_namespace,omitempty" yaml:"before_job_namespace"`
	Client             client.Client       `json:"client,omitempty" yaml:"client"`
}

// BenchmarkScenario represents the configuration for the internal benchmark scenario.
type BenchmarkScenario struct {
	Target  *BenchmarkTarget  `json:"target" yaml:"target"`
	Dataset *BenchmarkDataset `jon:"dataset" yaml:"dataset"`
	Jobs    []*BenchmarkJob   `job:"jobs" yaml:jobs`
}

// BenchmarkTarget defines the desired state of BenchmarkTarget
type BenchmarkTarget struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

// BenchmarkDataset defines the desired state of BenchmarkDateset
type BenchmarkDataset struct {
	Name    string                 `json:"name,omitempty"`
	Group   string                 `json:"group,omitempty"`
	Indexes int                    `json:"indexes,omitempty"`
	Range   *BenchmarkDatasetRange `json:"range,omitempty"`
}

// BenchmarkDatasetRange defines the desired state of BenchmarkDatesetRange
type BenchmarkDatasetRange struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

// BenchmarkJobRule defines the desired state of BenchmarkJobRule
type BenchmarkJobRule struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// InsertConfig defines the desired state of insert config
type InsertConfig struct {
	SkipStrictExistCheck bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp            string `json:"timestamp,omitempty"`
}

// UpdateConfig defines the desired state of update config
type UpdateConfig struct {
	SkipStrictExistCheck bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp            string `json:"timestamp,omitempty"`
}

// UpsertConfig defines the desired state of upsert config
type UpsertConfig struct {
	SkipStrictExistCheck bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp            string `json:"timestamp,omitempty"`
}

// SearchConfig defines the desired state of search config
type SearchConfig struct {
	Epsilon float32 `json:"epsilon,omitempty"`
	Radius  float32 `json:"radius,omitempty"`
	Num     int32   `json:"num,omitempty"`
	MinNum  int32   `json:"min_num,omitempty"`
	Timeout string  `json:"timeout,omitempty"`
}

// RemoveConfig defines the desired state of remove config
type RemoveConfig struct {
	SkipStrictExistCheck bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp            string `json:"timestamp,omitempty"`
}

// Bind binds the actual data from the Job receiver fields.
func (b *BenchmarkJob) Bind() *BenchmarkJob {
	b.JobType = GetActualValue(b.JobType)

	if b.ClientConfig != nil {
		b.ClientConfig = b.ClientConfig.Bind()
	}
	return b
}

// Bind binds the actual data from the BenchmarkScenario receiver fields.
func (b *BenchmarkScenario) Bind() *BenchmarkScenario {
	return b
}
