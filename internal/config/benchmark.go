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

import v1 "github.com/vdaas/vald/internal/k8s/vald/benchmark/api/v1"

// BenchmarkJob represents the configuration for the internal benchmark search job.
type BenchmarkJob struct {
	Target       *v1.BenchmarkTarget    `json:"target,omitempty" yaml:"target"`
	Dataset      *v1.BenchmarkDataset   `json:"dataset,omitempty" yaml:"dataset"`
	Dimension    int                    `json:"dimension,omitempty" yaml:"dimension"`
	Replica      int                    `json:"replica,omitempty" yaml:"replica"`
	Repetition   int                    `json:"repetition,omitempty" yaml:"repetition"`
	JobType      string                 `json:"job_type,omitempty" yaml:"job_type"`
	InsertConfig *v1.InsertConfig       `json:"insert_config,omitempty" yaml:"insert_config"`
	UpdateConfig *v1.UpdateConfig       `json:"update_config,omitempty" yaml:"update_config"`
	UpsertConfig *v1.UpsertConfig       `json:"upsert_config,omitempty" yaml:"upsert_config"`
	SearchConfig *v1.SearchConfig       `json:"search_config,omitempty" yaml:"search_config"`
	RemoveConfig *v1.RemoveConfig       `json:"remove_config,omitempty" yaml:"remove_config"`
	ClientConfig *GRPCClient            `json:"client_config,omitempty" yaml:"client_config"`
	Rules        []*v1.BenchmarkJobRule `json:"rules,omitempty" yaml:"rules"`
}

// BenchmarkScenario represents the configuration for the internal benchmark scenario.
type BenchmarkScenario struct {
	Target  *v1.BenchmarkTarget  `json:"target" yaml:"target"`
	Dataset *v1.BenchmarkDataset `jon:"dataset" yaml:"dataset"`
	Jobs    []*BenchmarkJob      `job:"jobs" yaml:jobs`
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
