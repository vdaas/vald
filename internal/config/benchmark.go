//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package config providers configuration type and load configuration logic
package config

// BenchmarkJob represents the configuration for the internal benchmark search job.
type BenchmarkJob struct {
	Target             *BenchmarkTarget    `json:"target,omitempty"               yaml:"target"`
	Dataset            *BenchmarkDataset   `json:"dataset,omitempty"              yaml:"dataset"`
	Replica            int                 `json:"replica,omitempty"              yaml:"replica"`
	Repetition         int                 `json:"repetition,omitempty"           yaml:"repetition"`
	JobType            string              `json:"job_type,omitempty"             yaml:"job_type"`
	InsertConfig       *InsertConfig       `json:"insert_config,omitempty"        yaml:"insert_config"`
	UpdateConfig       *UpdateConfig       `json:"update_config,omitempty"        yaml:"update_config"`
	UpsertConfig       *UpsertConfig       `json:"upsert_config,omitempty"        yaml:"upsert_config"`
	SearchConfig       *SearchConfig       `json:"search_config,omitempty"        yaml:"search_config"`
	RemoveConfig       *RemoveConfig       `json:"remove_config,omitempty"        yaml:"remove_config"`
	ObjectConfig       *ObjectConfig       `json:"object_config,omitempty"        yaml:"object_config"`
	ClientConfig       *GRPCClient         `json:"client_config,omitempty"        yaml:"client_config"`
	Rules              []*BenchmarkJobRule `json:"rules,omitempty"                yaml:"rules"`
	BeforeJobName      string              `json:"before_job_name,omitempty"      yaml:"before_job_name"`
	BeforeJobNamespace string              `json:"before_job_namespace,omitempty" yaml:"before_job_namespace"`
	RPS                int                 `json:"rps,omitempty"                  yaml:"rps"`
	ConcurrencyLimit   int                 `json:"concurrency_limit,omitempty"    yaml:"concurrency_limit"`
}

// BenchmarkScenario represents the configuration for the internal benchmark scenario.
type BenchmarkScenario struct {
	Target  *BenchmarkTarget  `json:"target,omitempty"  yaml:"target"`
	Dataset *BenchmarkDataset `json:"dataset,omitempty" yaml:"dataset"`
	Jobs    []*BenchmarkJob   `json:"jobs,omitempty"    yaml:"jobs"`
}

// BenchmarkTarget defines the desired state of BenchmarkTarget.
type BenchmarkTarget struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

func (t *BenchmarkTarget) Bind() *BenchmarkTarget {
	t.Host = GetActualValue(t.Host)
	return t
}

// BenchmarkDataset defines the desired state of BenchmarkDateset.
type BenchmarkDataset struct {
	Name    string                 `json:"name,omitempty"`
	Group   string                 `json:"group,omitempty"`
	Indexes int                    `json:"indexes,omitempty"`
	Range   *BenchmarkDatasetRange `json:"range,omitempty"`
	URL     string                 `json:"url,omitempty"`
}

func (d *BenchmarkDataset) Bind() *BenchmarkDataset {
	d.Name = GetActualValue(d.Name)
	d.Group = GetActualValue(d.Group)
	d.URL = GetActualValue(d.URL)
	return d
}

// BenchmarkDatasetRange defines the desired state of BenchmarkDatesetRange.
type BenchmarkDatasetRange struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

// BenchmarkJobRule defines the desired state of BenchmarkJobRule.
type BenchmarkJobRule struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

func (r *BenchmarkJobRule) Bind() *BenchmarkJobRule {
	r.Name = GetActualValue(r.Name)
	r.Type = GetActualValue(r.Type)
	return r
}

// InsertConfig defines the desired state of insert config.
type InsertConfig struct {
	SkipStrictExistCheck bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp            string `json:"timestamp,omitempty"`
}

func (cfg *InsertConfig) Bind() *InsertConfig {
	cfg.Timestamp = GetActualValue(cfg.Timestamp)
	return cfg
}

// UpdateConfig defines the desired state of update config.
type UpdateConfig struct {
	SkipStrictExistCheck  bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp             string `json:"timestamp,omitempty"`
	DisableBalancedUpdate bool   `json:"disable_balanced_update,omitempty"`
}

func (cfg *UpdateConfig) Bind() *UpdateConfig {
	cfg.Timestamp = GetActualValue(cfg.Timestamp)
	return cfg
}

// UpsertConfig defines the desired state of upsert config.
type UpsertConfig struct {
	SkipStrictExistCheck  bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp             string `json:"timestamp,omitempty"`
	DisableBalancedUpdate bool   `json:"disable_balanced_update,omitempty"`
}

func (cfg *UpsertConfig) Bind() *UpsertConfig {
	cfg.Timestamp = GetActualValue(cfg.Timestamp)
	return cfg
}

// SearchConfig defines the desired state of search config.
type SearchConfig struct {
	Epsilon              float32 `json:"epsilon,omitempty"`
	Radius               float32 `json:"radius,omitempty"`
	Num                  int32   `json:"num,omitempty"`
	MinNum               int32   `json:"min_num,omitempty"`
	Timeout              string  `json:"timeout,omitempty"`
	EnableLinearSearch   bool    `json:"enable_linear_search,omitempty"`
	AggregationAlgorithm string  `json:"aggregation_algorithm,omitempty"`
}

func (cfg *SearchConfig) Bind() *SearchConfig {
	cfg.Timeout = GetActualValue(cfg.Timeout)
	cfg.AggregationAlgorithm = GetActualValue(cfg.AggregationAlgorithm)
	return cfg
}

// RemoveConfig defines the desired state of remove config.
type RemoveConfig struct {
	SkipStrictExistCheck bool   `json:"skip_strict_exist_check,omitempty"`
	Timestamp            string `json:"timestamp,omitempty"`
}

func (cfg *RemoveConfig) Bind() *RemoveConfig {
	cfg.Timestamp = GetActualValue(cfg.Timestamp)
	return cfg
}

// ObjectConfig defines the desired state of object config.
type ObjectConfig struct {
	FilterConfig FilterConfig `json:"filter_config,omitempty" yaml:"filter_config"`
}

func (cfg *ObjectConfig) Bind() *ObjectConfig {
	cfg.FilterConfig = *cfg.FilterConfig.Bind()
	return cfg
}

// FilterTarget defines the desired state of filter target.
type FilterTarget struct {
	Host string `json:"host,omitempty" yaml:"host"`
	Port int32  `json:"port,omitempty" yaml:"port"`
}

func (cfg *FilterTarget) Bind() *FilterTarget {
	cfg.Host = GetActualValue(cfg.Host)
	return cfg
}

// FilterConfig defines the desired state of filter config.
type FilterConfig struct {
	Targets []*FilterTarget `json:"target,omitempty" yaml:"target"`
}

func (cfg *FilterConfig) Bind() *FilterConfig {
	for i := 0; i < len(cfg.Targets); i++ {
		cfg.Targets[i] = cfg.Targets[i].Bind()
	}
	return cfg
}

// Bind binds the actual data from the Job receiver fields.
func (b *BenchmarkJob) Bind() *BenchmarkJob {
	b.JobType = GetActualValue(b.JobType)
	b.BeforeJobName = GetActualValue(b.BeforeJobName)
	b.BeforeJobNamespace = GetActualValue(b.BeforeJobNamespace)

	if b.Target != nil {
		b.Target = b.Target.Bind()
	}
	if b.Dataset != nil {
		b.Dataset = b.Dataset.Bind()
	}
	if b.InsertConfig != nil {
		b.InsertConfig = b.InsertConfig.Bind()
	}
	if b.UpdateConfig != nil {
		b.UpdateConfig = b.UpdateConfig.Bind()
	}
	if b.UpsertConfig != nil {
		b.UpsertConfig = b.UpsertConfig.Bind()
	}
	if b.SearchConfig != nil {
		b.SearchConfig = b.SearchConfig.Bind()
	}
	if b.RemoveConfig != nil {
		b.RemoveConfig = b.RemoveConfig.Bind()
	}
	if b.ObjectConfig != nil {
		b.ObjectConfig = b.ObjectConfig.Bind()
	}
	if b.ClientConfig != nil {
		b.ClientConfig = b.ClientConfig.Bind()
	}
	if len(b.Rules) > 0 {
		for i := 0; i < len(b.Rules); i++ {
			b.Rules[i] = b.Rules[i].Bind()
		}
	}
	return b
}

// Bind binds the actual data from the BenchmarkScenario receiver fields.
func (b *BenchmarkScenario) Bind() *BenchmarkScenario {
	return b
}

// BenchmarkJobImageInfo represents the docker image information for benchmark job.
type BenchmarkJobImageInfo struct {
	Image      string `info:"image"       json:"image,omitempty"       yaml:"image"`
	PullPolicy string `info:"pull_policy" json:"pull_policy,omitempty" yaml:"pull_policy"`
}

// Bind binds the actual data from the BenchmarkJobImageInfo receiver fields.
func (b *BenchmarkJobImageInfo) Bind() *BenchmarkJobImageInfo {
	b.Image = GetActualValue(b.Image)
	b.PullPolicy = GetActualValue(b.PullPolicy)
	return b
}
