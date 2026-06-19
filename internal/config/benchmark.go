//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package config

// BenchmarkJob represents the configuration for the internal benchmark search job.
type BenchmarkJob struct {
	// RemoveConfig represents the remove configuration.
	RemoveConfig *RemoveConfig `json:"remove_config,omitempty" yaml:"remove_config"`
	// Dataset represents the dataset configuration.
	Dataset *BenchmarkDataset `json:"dataset,omitempty" yaml:"dataset"`
	// ClientConfig represents the client configuration.
	ClientConfig *GRPCClient `json:"client_config,omitempty" yaml:"client_config"`
	// ObjectConfig represents the object configuration.
	ObjectConfig *ObjectConfig `json:"object_config,omitempty" yaml:"object_config"`
	// Target represents the target configuration.
	Target *BenchmarkTarget `json:"target,omitempty" yaml:"target"`
	// InsertConfig represents the insert configuration.
	InsertConfig *InsertConfig `json:"insert_config,omitempty" yaml:"insert_config"`
	// UpdateConfig represents the update configuration.
	UpdateConfig *UpdateConfig `json:"update_config,omitempty" yaml:"update_config"`
	// UpsertConfig represents the upsert configuration.
	UpsertConfig *UpsertConfig `json:"upsert_config,omitempty" yaml:"upsert_config"`
	// SearchConfig represents the search configuration.
	SearchConfig *SearchConfig `json:"search_config,omitempty" yaml:"search_config"`
	// JobType represents the job type.
	JobType string `json:"job_type,omitempty" yaml:"job_type"`
	// BeforeJobName represents the before job name.
	BeforeJobName string `json:"before_job_name,omitempty" yaml:"before_job_name"`
	// BeforeJobNamespace represents the before job namespace.
	BeforeJobNamespace string `json:"before_job_namespace,omitempty" yaml:"before_job_namespace"`
	// Rules represents the list of rules.
	Rules []*BenchmarkJobRule `json:"rules,omitempty" yaml:"rules"`
	// Repetition represents the repetition count.
	Repetition int `json:"repetition,omitempty" yaml:"repetition"`
	// Replica represents the replica count.
	Replica int `json:"replica,omitempty" yaml:"replica"`
	// RPS represents the requests per second.
	RPS int `json:"rps,omitempty" yaml:"rps"`
	// ConcurrencyLimit represents the concurrency limit.
	ConcurrencyLimit int `json:"concurrency_limit,omitempty" yaml:"concurrency_limit"`
}

// BenchmarkScenario represents the configuration for the internal benchmark scenario.
type BenchmarkScenario struct {
	// Target represents the target configuration.
	Target *BenchmarkTarget `json:"target,omitempty" yaml:"target"`
	// Dataset represents the dataset configuration.
	Dataset *BenchmarkDataset `json:"dataset,omitempty" yaml:"dataset"`
	// Jobs represents the list of benchmark jobs.
	Jobs []*BenchmarkJob `json:"jobs,omitempty" yaml:"jobs"`
}

// BenchmarkTarget defines the desired state of BenchmarkTarget.
type BenchmarkTarget struct {
	// Meta represents the metadata.
	Meta map[string]string `json:"meta,omitempty"`
	// Host represents the host.
	Host string `json:"host,omitempty"`
	// Port represents the port.
	Port int `json:"port,omitempty"`
}

func (t *BenchmarkTarget) Bind() *BenchmarkTarget {
	t.Host = GetActualValue(t.Host)
	return t
}

// BenchmarkDataset defines the desired state of BenchmarkDateset.
type BenchmarkDataset struct {
	// Range represents the range configuration.
	Range *BenchmarkDatasetRange `json:"range,omitempty"`
	// Name represents the dataset name.
	Name string `json:"name,omitempty"`
	// Group represents the dataset group.
	Group string `json:"group,omitempty"`
	// URL represents the dataset URL.
	URL string `json:"url,omitempty"`
	// Indexes represents the number of indexes.
	Indexes int `json:"indexes,omitempty"`
}

func (d *BenchmarkDataset) Bind() *BenchmarkDataset {
	d.Name = GetActualValue(d.Name)
	d.Group = GetActualValue(d.Group)
	d.URL = GetActualValue(d.URL)
	return d
}

// BenchmarkDatasetRange defines the desired state of BenchmarkDatesetRange.
type BenchmarkDatasetRange struct {
	// Start represents the start index.
	Start int `json:"start,omitempty"`
	// End represents the end index.
	End int `json:"end,omitempty"`
}

// BenchmarkJobRule defines the desired state of BenchmarkJobRule.
type BenchmarkJobRule struct {
	// Name represents the rule name.
	Name string `json:"name,omitempty"`
	// Type represents the rule type.
	Type string `json:"type,omitempty"`
}

func (r *BenchmarkJobRule) Bind() *BenchmarkJobRule {
	r.Name = GetActualValue(r.Name)
	r.Type = GetActualValue(r.Type)
	return r
}

// InsertConfig defines the desired state of insert config.
type InsertConfig struct {
	// Timestamp represents the timestamp.
	Timestamp string `json:"timestamp,omitempty"`
	// SkipStrictExistCheck enables skipping strict existence check.
	SkipStrictExistCheck bool `json:"skip_strict_exist_check,omitempty"`
}

func (cfg *InsertConfig) Bind() *InsertConfig {
	cfg.Timestamp = GetActualValue(cfg.Timestamp)
	return cfg
}

// UpdateConfig defines the desired state of update config.
type UpdateConfig struct {
	// Timestamp represents the timestamp.
	Timestamp string `json:"timestamp,omitempty"`
	// SkipStrictExistCheck enables skipping strict existence check.
	SkipStrictExistCheck bool `json:"skip_strict_exist_check,omitempty"`
	// DisableBalancedUpdate disables balanced update.
	DisableBalancedUpdate bool `json:"disable_balanced_update,omitempty"`
}

func (cfg *UpdateConfig) Bind() *UpdateConfig {
	cfg.Timestamp = GetActualValue(cfg.Timestamp)
	return cfg
}

// UpsertConfig defines the desired state of upsert config.
type UpsertConfig struct {
	// Timestamp represents the timestamp.
	Timestamp string `json:"timestamp,omitempty"`
	// SkipStrictExistCheck enables skipping strict existence check.
	SkipStrictExistCheck bool `json:"skip_strict_exist_check,omitempty"`
	// DisableBalancedUpdate disables balanced update.
	DisableBalancedUpdate bool `json:"disable_balanced_update,omitempty"`
}

func (cfg *UpsertConfig) Bind() *UpsertConfig {
	cfg.Timestamp = GetActualValue(cfg.Timestamp)
	return cfg
}

// SearchConfig defines the desired state of search config.
type SearchConfig struct {
	// Timeout represents the timeout.
	Timeout string `json:"timeout,omitempty"`
	// AggregationAlgorithm represents the aggregation algorithm.
	AggregationAlgorithm string `json:"aggregation_algorithm,omitempty"`
	// Epsilon represents the epsilon.
	Epsilon float32 `json:"epsilon,omitempty"`
	// Radius represents the radius.
	Radius float32 `json:"radius,omitempty"`
	// Num represents the number of results.
	Num int32 `json:"num,omitempty"`
	// MinNum represents the minimum number of results.
	MinNum int32 `json:"min_num,omitempty"`
	// EnableLinearSearch enables linear search.
	EnableLinearSearch bool `json:"enable_linear_search,omitempty"`
}

func (cfg *SearchConfig) Bind() *SearchConfig {
	cfg.Timeout = GetActualValue(cfg.Timeout)
	cfg.AggregationAlgorithm = GetActualValue(cfg.AggregationAlgorithm)
	return cfg
}

// RemoveConfig defines the desired state of remove config.
type RemoveConfig struct {
	// Timestamp represents the timestamp.
	Timestamp string `json:"timestamp,omitempty"`
	// SkipStrictExistCheck enables skipping strict existence check.
	SkipStrictExistCheck bool `json:"skip_strict_exist_check,omitempty"`
}

func (cfg *RemoveConfig) Bind() *RemoveConfig {
	cfg.Timestamp = GetActualValue(cfg.Timestamp)
	return cfg
}

// ObjectConfig defines the desired state of object config.
type ObjectConfig struct {
	// FilterConfig represents the filter configuration.
	FilterConfig FilterConfig `json:"filter_config" yaml:"filter_config"`
}

func (cfg *ObjectConfig) Bind() *ObjectConfig {
	cfg.FilterConfig.Bind()
	return cfg
}

// FilterTarget defines the desired state of filter target.
type FilterTarget struct {
	// Host represents the host.
	Host string `json:"host,omitempty" yaml:"host"`
	// Port represents the port.
	Port int32 `json:"port,omitempty" yaml:"port"`
}

func (cfg *FilterTarget) Bind() *FilterTarget {
	cfg.Host = GetActualValue(cfg.Host)
	return cfg
}

// FilterConfig defines the desired state of filter config.
type FilterConfig struct {
	// Targets represents the list of filter targets.
	Targets []*FilterTarget `json:"target,omitempty" yaml:"target"`
}

func (cfg *FilterConfig) Bind() *FilterConfig {
	for i := 0; i < len(cfg.Targets); i++ {
		if cfg.Targets[i] != nil {
			cfg.Targets[i].Bind()
		}
	}
	return cfg
}

// Bind binds the actual data from the Job receiver fields.
func (b *BenchmarkJob) Bind() *BenchmarkJob {
	b.JobType = GetActualValue(b.JobType)
	b.BeforeJobName = GetActualValue(b.BeforeJobName)
	b.BeforeJobNamespace = GetActualValue(b.BeforeJobNamespace)

	if b.Target != nil {
		b.Target.Bind()
	}
	if b.Dataset != nil {
		b.Dataset.Bind()
	}
	if b.InsertConfig != nil {
		b.InsertConfig.Bind()
	}
	if b.UpdateConfig != nil {
		b.UpdateConfig.Bind()
	}
	if b.UpsertConfig != nil {
		b.UpsertConfig.Bind()
	}
	if b.SearchConfig != nil {
		b.SearchConfig.Bind()
	}
	if b.RemoveConfig != nil {
		b.RemoveConfig.Bind()
	}
	if b.ObjectConfig != nil {
		b.ObjectConfig.Bind()
	}
	if b.ClientConfig != nil {
		b.ClientConfig.Bind()
	}
	if len(b.Rules) > 0 {
		for i := 0; i < len(b.Rules); i++ {
			if b.Rules[i] != nil {
				b.Rules[i].Bind()
			}
		}
	}
	return b
}

// Bind binds the actual data from the BenchmarkScenario receiver fields.
func (b *BenchmarkScenario) Bind() *BenchmarkScenario {
	if b.Target != nil {
		b.Target.Bind()
	}
	if b.Dataset != nil {
		b.Dataset.Bind()
	}
	if len(b.Jobs) > 0 {
		for i := range b.Jobs {
			if b.Jobs[i] != nil {
				b.Jobs[i].Bind()
			}
		}
	}
	return b
}

// BenchmarkJobImageInfo represents the docker image information for benchmark job.
type BenchmarkJobImageInfo struct {
	// Repository represents the image repository.
	Repository string `info:"repository" json:"repository,omitempty" yaml:"repository"`
	// Tag represents the image tag.
	Tag string `info:"tag" json:"tag,omitempty" yaml:"tag"`
	// PullPolicy represents the image pull policy.
	PullPolicy string `info:"pull_policy" json:"pull_policy,omitempty" yaml:"pull_policy"`
}

// Bind binds the actual data from the BenchmarkJobImageInfo receiver fields.
func (b *BenchmarkJobImageInfo) Bind() *BenchmarkJobImageInfo {
	b.Repository = GetActualValue(b.Repository)
	b.Tag = GetActualValue(b.Tag)
	b.PullPolicy = GetActualValue(b.PullPolicy)
	return b
}

// OperatorJobConfig represents the general job configuration for operator.
type OperatorJobConfig struct {
	// Image represents the image information.
	Image *BenchmarkJobImageInfo `info:"image" json:"image,omitempty" yaml:"image"`
	// BenchmarkJob represents the benchmark job configuration.
	*BenchmarkJob
}
