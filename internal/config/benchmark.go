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
	Target        *v1.BenchmarkTarget    `json:"target" yaml:"target"`
	Dataset       *v1.BenchmarkDataset   `json:"dataset"        yaml:"dataset"`
	Replica       int                    `json:"replica" yaml:"replica"`
	Repetition    int                    `json:"repetition" yaml:"repetition"`
	JobType       string                 `json:"job_type"       yaml:"job_type"`
	Dimension     int                    `json:"dimension"      yaml:"dimension"`
	Epsilon       float64                `json:"epsilon"        yaml:"epsilon"`
	Radius        float64                `json:"radius"         yaml:"radius"`
	Iter          int                    `json:"iter"           yaml:"iter"`
	Num           uint32                 `json:"num"            yaml:"num"`
	MinNum        uint32                 `json:"min_num"        yaml:"min_num"`
	Timeout       string                 `json:"timeout"        yaml:"timeout"`
	Rules         []*v1.BenchmarkJobRule `json:"rules,omitempty" yaml:"rules,omitempty"`
	GatewayClient *GRPCClient            `json:"gateway_client" yaml:"gateway_client"`
}

// BenchmarkScenario represents the configuration for the internal benchmark scenario.
type BenchmarkScenario struct {
	Target  *v1.BenchmarkTarget  `json:"target" yaml:"target"`
	Dataset *v1.BenchmarkDataset `jon:"dataset" yaml:"dataset"`
	Jobs    []*BenchmarkJob      `job:"jobs" yaml:jobs`
}

// Bind binds the actual data from the Job receiver fields.
func (b *BenchmarkJob) Bind() *BenchmarkJob {
	b.Timeout = GetActualValue(b.Timeout)
	b.JobType = GetActualValue(b.JobType)

	if b.GatewayClient != nil {
		b.GatewayClient = b.GatewayClient.Bind()
	}
	return b
}

// Bind binds the actual data from the BenchmarkScenario receiver fields.
func (b *BenchmarkScenario) Bind() *BenchmarkScenario {
	return b
}
