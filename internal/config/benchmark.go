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

// BenchmarkJob represents the configuration for the internal benchmark search job.
type BenchmarkJob struct {
	Target        *BenchmarkTarget    `json:"target" yaml:"target"`
	JobType       string              `json:"job_type"       yaml:"job_type"`
	Dataset       *BenchmarkDataset   `json:"dataset"        yaml:"dataset"`
	Replica       int                 `json:"replica" yaml:"replica"`
	Repetition    int                 `json:"repetition" yaml:"repetition"`
	Dimension     int                 `json:"dimension"      yaml:"dimension"`
	Iter          int                 `json:"iter"           yaml:"iter"`
	Num           uint32              `json:"num"            yaml:"num"`
	MinNum        uint32              `json:"min_num"        yaml:"min_num"`
	Radius        float64             `json:"radius"         yaml:"radius"`
	Epsilon       float64             `json:"epsilon"        yaml:"epsilon"`
	Timeout       string              `json:"timeout"        yaml:"timeout"`
	Rules         []*BenchmarkJobRule `json:"rules,omitempty" yaml:"rules,omitempty"`
	GatewayClient *GRPCClient         `json:"gateway_client" yaml:"gateway_client"`
}

// BenchmarkScenario represents the configuration for the internal benchmark scenario.
type BenchmarkScenario struct {
	Target  *BenchmarkTarget  `json:"target" yaml:"target"`
	Dataset *BenchmarkDataset `jon:"dataset" yaml:"dataset"`
	Jobs    []*BenchmarkJob   `job:"jobs" yaml:jobs`
}

// BenchmarkTarget defines the desired state of BenchmarkTarget.
type BenchmarkTarget struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// BenchmarkDataset defines the desired state of BenchmarkDateset.
type BenchmarkDataset struct {
	Name    string                 `json:"name" yaml:"name"`
	Group   string                 `json:"group" yaml:"group"`
	Indexes int                    `json:"indexes" yaml:"indexes"`
	Range   *BenchmarkDatasetRange `json:"range" yaml:"range"`
}

// BenchmarkDatasetRange defines the desired state of BenchmarkDatesetRange.
type BenchmarkDatasetRange struct {
	Start int `json:"start" yaml:"start"`
	End   int `json:"end" yaml:"end"`
}

// BenchmarkJobRule defines the desired state of BenchmarkJobRule.
type BenchmarkJobRule struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
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
