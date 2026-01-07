//go:build e2e

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

// Package config provides configuration types and logic for loading and binding configuration values.
// This file includes refactored Bind methods (always returning error) and non-Bind functions,
// with named return values and proper ordering of sections.
package config

import (
	"io/fs"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/timeutil"
	"github.com/vdaas/vald/internal/timeutil/rate"
	"github.com/vdaas/vald/tests/v2/e2e/kubernetes"
	"github.com/vdaas/vald/tests/v2/e2e/metrics"
	"sigs.k8s.io/yaml"
)

////////////////////////////////////////////////////////////////////////////////
// Struct Section
////////////////////////////////////////////////////////////////////////////////

// Data represents the complete configuration for the application.
type Data struct {
	config.GlobalConfig `json:",inline" yaml:",inline"`
	TimeConfig          `json:",inline" yaml:",inline"`
	Target              *config.GRPCClient `json:"target,omitempty"          yaml:"target,omitempty"`
	Strategies          []*Strategy        `json:"strategies,omitempty"      yaml:"strategies,omitempty"`
	Dataset             *Dataset           `json:"dataset,omitempty"         yaml:"dataset,omitempty"`
	Kubernetes          *Kubernetes        `json:"kubernetes,omitempty"      yaml:"kubernetes,omitempty"`
	Metrics             *Metrics           `json:"metrics,omitempty"         yaml:"metrics,omitempty"`
	Metadata            map[string]string  `json:"metadata,omitempty"        yaml:"metadata,omitempty"`
	MetaString          string             `json:"metadata_string,omitempty" yaml:"metadata_string,omitempty"`
	FilePath            string             `json:"-"                         yaml:"-"`
	Collector           metrics.Collector  `json:"-"                         yaml:"-"`
}

// Metrics represents the configuration for the metrics collector.
type Metrics struct {
	Enabled               bool          `json:"enabled"                        yaml:"enabled"`
	Histogram             *Histogram    `json:"histogram,omitempty"            yaml:"histogram,omitempty"`
	LatencyHistogram      *Histogram    `json:"latency_histogram,omitempty"    yaml:"latency_histogram,omitempty"`
	QueueWaitHistogram    *Histogram    `json:"queue_wait_histogram,omitempty" yaml:"queue_wait_histogram,omitempty"`
	TDigest               *TDigest      `json:"tdigest,omitempty"              yaml:"tdigest,omitempty"`
	LatencyTDigest        *TDigest      `json:"latency_tdigest,omitempty"      yaml:"latency_tdigest,omitempty"`
	QueueWaitTDigest      *TDigest      `json:"queue_wait_tdigest,omitempty"   yaml:"queue_wait_tdigest,omitempty"`
	Exemplar              *Exemplar     `json:"exemplar,omitempty"             yaml:"exemplar,omitempty"`
	RangeScales           []*RangeScale `json:"range_scales,omitempty"         yaml:"range_scales,omitempty"`
	TimeScales            []*TimeScale  `json:"time_scales,omitempty"          yaml:"time_scales,omitempty"`
	CustomCounters        []string      `json:"custom_counters,omitempty"      yaml:"custom_counters,omitempty"`
	DetailedErrorTracking bool          `json:"detailed_error_tracking"        yaml:"detailed_error_tracking"`
}

// RangeScale represents the configuration for a range scale.
type RangeScale struct {
	Name     string `json:"name,omitempty"     yaml:"name,omitempty"`
	Width    uint64 `json:"width,omitempty"    yaml:"width,omitempty"`
	Capacity uint64 `json:"capacity,omitempty" yaml:"capacity,omitempty"`
}

// TimeScale represents the configuration for a time scale.
type TimeScale struct {
	Name     string `json:"name,omitempty"     yaml:"name,omitempty"`
	Width    uint64 `json:"width,omitempty"    yaml:"width,omitempty"`
	Capacity uint64 `json:"capacity,omitempty" yaml:"capacity,omitempty"`
}

// Histogram represents the configuration for a histogram.
type Histogram struct {
	NumShards int `json:"num_shards,omitempty" yaml:"num_shards,omitempty"`
}

// TDigest represents the configuration for a TDigest.
type TDigest struct {
	Compression              float64   `json:"compression,omitempty"                yaml:"compression,omitempty"`
	CompressionTriggerFactor float64   `json:"compression_trigger_factor,omitempty" yaml:"compression_trigger_factor,omitempty"`
	Quantiles                []float64 `json:"quantiles,omitempty"                  yaml:"quantiles,omitempty"`
	NumShards                int       `json:"num_shards,omitempty"                 yaml:"num_shards,omitempty"`
}

// Exemplar represents the configuration for an exemplar.
type Exemplar struct {
	Capacity     int `json:"capacity,omitempty"      yaml:"capacity,omitempty"`
	NumShards    int `json:"num_shards,omitempty"    yaml:"num_shards,omitempty"`
	SamplingRate int `json:"sampling_rate,omitempty" yaml:"sampling_rate,omitempty"`
}

// Strategy represents a test strategy.
type Strategy struct {
	TimeConfig  `yaml:",inline" json:",inline"`
	Name        string            `yaml:"name"                 json:"name,omitempty"`
	Repeats     *Repeats          `yaml:"repeats"              json:"repeats,omitempty"`
	Concurrency uint64            `yaml:"concurrency"          json:"concurrency,omitempty"`
	Operations  []*Operation      `yaml:"operations,omitempty" json:"operations,omitempty"`
	Metrics     *Metrics          `yaml:"metrics,omitempty"    json:"metrics,omitempty"`
	Collector   metrics.Collector `yaml:"-"                    json:"-"`
}

// Operation represents an individual operation configuration.
type Operation struct {
	TimeConfig `yaml:",inline" json:",inline"`
	Name       string            `yaml:"name,omitempty"       json:"name,omitempty"`
	Repeats    *Repeats          `yaml:"repeats"              json:"repeats,omitempty"`
	Executions []*Execution      `yaml:"executions,omitempty" json:"executions,omitempty"`
	Metrics    *Metrics          `yaml:"metrics,omitempty"    json:"metrics,omitempty"`
	Collector  metrics.Collector `yaml:"-"                    json:"-"`
}

// Execution represents the execution details for a given operation.
type Execution struct {
	*BaseConfig  `yaml:",inline,omitempty" json:",inline,omitempty"`
	TimeConfig   `yaml:",inline" json:",inline"`
	Name         string              `yaml:"name"                   json:"name,omitempty"`
	Strategy     string              `yaml:"-"                      json:"-"`
	Operation    string              `yaml:"-"                      json:"-"`
	Repeats      *Repeats            `yaml:"repeats"                json:"repeats,omitempty"`
	Type         OperationType       `yaml:"type"                   json:"type,omitempty"`
	Mode         OperationMode       `yaml:"mode"                   json:"mode,omitempty"`
	Search       *SearchQuery        `yaml:"search,omitempty"       json:"search,omitempty"`
	Agent        *AgentConfig        `yaml:"agent,omitempty"        json:"agent,omitempty"`
	Kubernetes   *KubernetesConfig   `yaml:"kubernetes,omitempty"   json:"kubernetes,omitempty"`
	Modification *ModificationConfig `yaml:"modification,omitempty" json:"modification,omitempty"`
	Expect       []Expect            `yaml:"expect,omitempty"       json:"expect,omitempty"`
	Collector    metrics.Collector   `yaml:"-"                      json:"-"`
	Metrics      *Metrics            `yaml:"metrics,omitempty"      json:"metrics,omitempty"`
}

// TimeConfig holds time-related configuration values.
type TimeConfig struct {
	Delay   timeutil.DurationString `yaml:"delay"   json:"delay,omitempty"`
	Wait    timeutil.DurationString `yaml:"wait"    json:"wait,omitempty"`
	Timeout timeutil.DurationString `yaml:"timeout" json:"timeout,omitempty"`
}

// BaseConfig represents basic operational configuration parameters.
type BaseConfig struct {
	Num         uint64       `yaml:"num,omitempty"         json:"num,omitempty"`
	Offset      uint64       `yaml:"offset,omitempty"      json:"offset,omitempty"`
	BulkSize    uint64       `yaml:"bulk_size,omitempty"   json:"bulk_size,omitempty"`
	Parallelism uint64       `yaml:"parallelism,omitempty" json:"parallelism,omitempty"`
	QPS         uint64       `yaml:"qps,omitempty"         json:"qps,omitempty"`
	Limiter     rate.Limiter `yaml:"-"                     json:"-"`
}

// SearchQuery represents the parameters for a search query.
type SearchQuery struct {
	K               uint32                              `yaml:"k,omitempty"         json:"k,omitempty"`
	Radius          float32                             `yaml:"radius,omitempty"    json:"radius,omitempty"`
	Epsilon         float32                             `yaml:"epsilon,omitempty"   json:"epsilon,omitempty"`
	AlgorithmString string                              `yaml:"algorithm,omitempty" json:"algorithm_string,omitempty"`
	MinNum          uint32                              `yaml:"min_num,omitempty"   json:"min_num,omitempty"`
	Ratio           float32                             `yaml:"ratio,omitempty"     json:"ratio,omitempty"`
	Nprobe          uint32                              `yaml:"nprobe,omitempty"    json:"nprobe,omitempty"`
	Timeout         timeutil.DurationString             `yaml:"timeout,omitempty"   json:"timeout,omitempty"`
	Algorithm       payload.Search_AggregationAlgorithm `yaml:"-"                   json:"-"`
}

// ModificationConfig represents settings for modifications like insert or update.
type ModificationConfig struct {
	SkipStrictExistCheck bool  `yaml:"skip_strict_exist_check,omitempty" json:"skip_strict_exist_check,omitempty"`
	Timestamp            int64 `yaml:"timestamp,omitempty"               json:"timestamp,omitempty"`
}

// AgentConfig represents settings for agent for createting index
type AgentConfig struct {
	PoolSize uint32 `yaml:"pool_size,omitempty" json:"pool_size,omitempty"`
}

// KubernetesConfig holds Kubernetes-specific settings.
type KubernetesConfig struct {
	Kind          KubernetesKind   `yaml:"kind"           json:"kind,omitempty"`
	Namespace     string           `yaml:"namespace"      json:"namespace,omitempty"`
	Name          string           `yaml:"name"           json:"name,omitempty"`
	LabelSelector string           `yaml:"label_selector" json:"label_selector,omitempty"`
	Action        KubernetesAction `yaml:"action"         json:"action,omitempty"`
	Status        KubernetesStatus `yaml:"status"         json:"status,omitempty"`
}

// Kubernetes holds configuration for Kubernetes environments.
type Kubernetes struct {
	KubeConfig  string       `yaml:"kubeconfig"            json:"kube_config,omitempty"`
	PortForward *PortForward `yaml:"portforward,omitempty" json:"port_forward,omitempty"`
}

// PortForward holds configuration for port forwarding.
type PortForward struct {
	Enabled     bool   `yaml:"enabled"      json:"enabled,omitempty"`
	TargetPort  Port   `yaml:"target_port"  json:"target_port,omitempty"`
	LocalPort   Port   `yaml:"local_port"   json:"local_port,omitempty"`
	Namespace   string `yaml:"namespace"    json:"namespace,omitempty"`
	ServiceName string `yaml:"service_name" json:"service_name,omitempty"`
}

// Port represents a port as a string.
type Port string

// Dataset holds dataset-related configuration.
type Dataset struct {
	Name string `yaml:"name" json:"name,omitempty"`
}

// Expect holds expected results for executions.
type Expect struct {
	StatusCode StatusCode `yaml:"status_code,omitempty" json:"status_code,omitempty"`
	Path       string     `yaml:"path,omitempty"        json:"path,omitempty"`
	Op         Operator   `yaml:"op,omitempty"          json:"op,omitempty"`
	Value      any        `yaml:"value,omitempty"       json:"value,omitempty"`
}

// Repeats holds the repeat configuration for operations.
type Repeats struct {
	Enabled       bool                    `yaml:"enabled,omitempty"        json:"enabled,omitempty"`
	ExitCondition ExitCondition           `yaml:"exit_condition,omitempty" json:"exit_condition,omitempty"`
	Count         uint64                  `yaml:"count,omitempty"          json:"count,omitempty"`
	Interval      timeutil.DurationString `yaml:"interval,omitempty"       json:"interval,omitempty"`
}

////////////////////////////////////////////////////////////////////////////////
// Bind Section
////////////////////////////////////////////////////////////////////////////////

// Bind binds and validates the Data configuration.
// It processes nested configurations and metadata.
func (d *Data) Bind() (bound *Data, err error) {
	if d == nil ||
		d.Strategies == nil || len(d.Strategies) == 0 ||
		d.Dataset == nil ||
		d.Target == nil {
		return nil, errors.Wrap(errors.ErrInvalidConfig, "missing required fields on Data")
	}
	d.GlobalConfig.Bind()

	// Bind gRPC Target configuration if provided.
	if d.Target != nil {
		d.Target.Bind()
	}
	// Bind Metrics.
	if d.Metrics != nil {
		if m, err := d.Metrics.Bind(); err != nil {
			return nil, errors.Wrap(err, "failed to bind Metrics configuration")
		} else if m != nil {
			d.Metrics = m
		}
		if d.Metrics.Enabled {
			if d.Collector, err = metrics.NewCollector(d.Metrics.Opts()...); err != nil {
				return nil, errors.Wrap(err, "failed to create metrics collector")
			}
		}
	}
	// Bind Dataset.
	if d.Dataset != nil {
		if ds, err := d.Dataset.Bind(); err != nil {
			return nil, errors.Wrapf(err, "failed to bind dataset configuration for %s", d.Dataset.Name)
		} else if ds != nil {
			d.Dataset = ds
		}
	}
	// Bind Kubernetes.
	if d.Kubernetes != nil {
		if k, err := d.Kubernetes.Bind(); err != nil {
			return nil, errors.Wrapf(err, "failed to bind Kubernetes configuration for %s", d.Kubernetes.KubeConfig)
		} else if k != nil {
			d.Kubernetes = k
		}
	}
	// Process metadata.
	if d.Metadata == nil {
		d.Metadata = make(map[string]string)
	}
	for _, meta := range strings.Split(config.GetActualValue(d.MetaString), ",") {
		key, val, ok := strings.Cut(meta, "=")
		if ok && key != "" && val != "" {
			d.Metadata[config.GetActualValue(key)] = config.GetActualValue(val)
		}
	}

	// Bind each Strategy.
	var cnt int
	for _, strategy := range d.Strategies {
		if strategy != nil {
			var bs *Strategy
			if bs, err = strategy.Bind(d.Metrics); err != nil {
				return nil, errors.Wrapf(err, "failed to bind strategy: %s", strategy.Name)
			} else if bs != nil {
				d.Strategies[cnt] = bs
				cnt++
			}
		}
	}
	d.Strategies = d.Strategies[:cnt]
	return d, nil
}

// Bind binds and validates the Strategy configuration.
func (s *Strategy) Bind(parentMetrics *Metrics) (bound *Strategy, err error) {
	if s == nil || s.Operations == nil || len(s.Operations) == 0 {
		return nil, errors.Wrapf(errors.ErrInvalidConfig, "missing required fields on Strategy %s", s.Name)
	}
	s.Name = config.GetActualValue(s.Name)
	s.TimeConfig.Bind()

	// Bind Metrics.
	if s.Metrics != nil {
		if m, err := parentMetrics.Merge(s.Metrics).Bind(); err != nil {
			return nil, errors.Wrap(err, "failed to bind Metrics configuration")
		} else if m != nil {
			s.Metrics = m
		}
	} else {
		s.Metrics = parentMetrics
	}
	if s.Metrics != nil && s.Metrics.Enabled {
		if s.Collector, err = metrics.NewCollector(s.Metrics.Opts()...); err != nil {
			return nil, errors.Wrap(err, "failed to create metrics collector")
		}
	}
	var cnt int
	for _, op := range s.Operations {
		if op != nil {
			var bo *Operation
			if bo, err = op.Bind(s.Name, s.Metrics); err != nil {
				return nil, errors.Wrapf(err, "failed to bind operation: %s", op.Name)
			} else if bo != nil {
				s.Operations[cnt] = bo
				cnt++
			}
		}
	}
	return s, nil
}

// Bind binds and validates the Operation configuration.
func (o *Operation) Bind(strategy string, parentMetrics *Metrics) (bound *Operation, err error) {
	if o == nil || o.Executions == nil || len(o.Executions) == 0 {
		return nil, errors.Wrapf(errors.ErrInvalidConfig, "missing required fields on Operation %s", o.Name)
	}
	o.Name = config.GetActualValue(o.Name)
	o.TimeConfig.Bind()
	// Bind Metrics.
	if o.Metrics != nil {
		if m, err := parentMetrics.Merge(o.Metrics).Bind(); err != nil {
			return nil, errors.Wrap(err, "failed to bind Metrics configuration")
		} else if m != nil {
			o.Metrics = m
		}
	} else {
		o.Metrics = parentMetrics
	}
	if o.Metrics != nil && o.Metrics.Enabled {
		if o.Collector, err = metrics.NewCollector(o.Metrics.Opts()...); err != nil {
			return nil, errors.Wrap(err, "failed to create metrics collector")
		}
	}
	var cnt int
	for _, exec := range o.Executions {
		exec.Strategy = strategy
		exec.Operation = o.Name
		var be *Execution
		if be, err = exec.Bind(o.Metrics); err != nil {
			return nil, errors.Wrapf(err, "failed to bind execution: %s", exec.Name)
		} else if be != nil {
			o.Executions[cnt] = be
			cnt++
		}
	}
	return o, nil
}

// Bind binds and validates the Execution configuration.
func (e *Execution) Bind(parentMetrics *Metrics) (bound *Execution, err error) {
	if e == nil {
		return nil, errors.Wrap(errors.ErrInvalidConfig, "missing required fields on Execution")
	}

	// Bind Metrics.
	if e.Metrics != nil {
		if m, err := parentMetrics.Merge(e.Metrics).Bind(); err != nil {
			return nil, errors.Wrap(err, "failed to bind Metrics configuration")
		} else if m != nil {
			e.Metrics = m
		}
	} else {
		e.Metrics = parentMetrics
	}
	if e.Metrics != nil && e.Metrics.Enabled {
		if e.Collector, err = metrics.NewCollector(e.Metrics.Opts()...); err != nil {
			return nil, errors.Wrap(err, "failed to create metrics collector")
		}
	}

	// Bind OperationType and OperationMode.
	if e.Type, err = e.Type.Bind(); err != nil {
		return nil, errors.Wrapf(err, "failed to bind OperationType: %s on Execution %s", e.Type, e.Name)
	}
	if e.Mode, err = e.Mode.Bind(); err != nil {
		return nil, errors.Wrapf(err, "failed to bind OperationMode: %s on Execution %s", e.Mode, e.Name)
	}
	e.Name = config.GetActualValue(e.Name)
	e.TimeConfig.Bind()
	if e.Expect != nil {
		for i, ex := range e.Expect {
			if ex.StatusCode, err = ex.StatusCode.Bind(); err != nil {
				return nil, errors.Wrapf(err, "failed to bind StatusCodes for Execution %s of type %s", e.Name, e.Type)
			}
			if e.Mode != OperationUnary && ex.Value != nil {
				return nil, errors.Wrapf(errors.ErrInvalidConfig, "Expect.Value is only supported for unary operations in Execution %s of type %s", e.Name, e.Type)
			}
			if ex.Op, err = ex.Op.Bind(); err != nil {
				return nil, errors.Wrapf(err, "failed to bind Expect.Op for Execution %s of type %s", e.Name, e.Type)
			}
			e.Expect[i] = ex
		}
	}
	switch e.Type {
	case OpSearch,
		OpSearchByID,
		OpLinearSearch,
		OpLinearSearchByID,
		OpInsert,
		OpUpdate,
		OpUpsert,
		OpRemove,
		OpRemoveByTimestamp,
		OpObject,
		OpListObject,
		OpTimestamp,
		OpExists:
		if e.BaseConfig == nil || e.BaseConfig.Num == 0 {
			return nil, errors.Wrapf(errors.ErrInvalidConfig, "BaseConfig and its Num are required for Execution %s of type %s", e.Name, e.Type)
		}
		if e.BaseConfig.QPS > 0 {
			e.Limiter = rate.NewLimiter(int(e.BaseConfig.QPS))
		}
		if e.Mode == OperationMultiple && e.BaseConfig.BulkSize == 0 {
			return nil, errors.Errorf("bulk_size must be greater than 0 for multiple operations for Execution %s of type %s of mode %s", e.Name, e.Type, e.Mode)
		}
		switch e.Type {
		case OpSearch,
			OpSearchByID,
			OpLinearSearch,
			OpLinearSearchByID:
			if e.Search == nil {
				return nil, errors.Errorf("missing required fields on SearchQuery for Execution %s of type %s", e.Name, e.Type)
			}
			if sq, err := e.Search.Bind(); err != nil {
				return nil, errors.Wrapf(err, "failed to bind SearchQuery for Execution %s of type %s", e.Name, e.Type)
			} else if sq != nil {
				e.Search = sq
			}
		case OpInsert,
			OpUpdate,
			OpUpsert,
			OpRemove,
			OpRemoveByTimestamp:
			if e.Modification != nil {
				if m, err := e.Modification.Bind(); err != nil {
					return nil, errors.Wrapf(err, "failed to bind ModificationConfig for Execution %s of type %s", e.Name, e.Type)
				} else if m != nil {
					e.Modification = m
				}
			}
		}
	case OpIndexInfo,
		OpIndexDetail,
		OpIndexStatistics,
		OpIndexStatisticsDetail,
		OpIndexProperty,
		OpFlush:
	case OpCreateIndex,
		OpSaveIndex,
		OpCreateAndSaveIndex:
		if e.Agent == nil {
			e.Agent = new(AgentConfig)
		}
		if e.Agent.PoolSize == 0 {
			e.Agent.PoolSize = uint32(runtime.GOMAXPROCS(-1))
		}
	case OpKubernetes:
		if e.Kubernetes != nil {
			if ek, err := e.Kubernetes.Bind(); err != nil {
				return nil, errors.Wrapf(err, "failed to bind Kubernetes configuration for Execution: %s, detail %v", e.Name, e.Kubernetes)
			} else if ek != nil {
				e.Kubernetes = ek
			}
		}
	case OpClient:
		// do nothing
	case OpWait:
		// do nothing
	default:
		return nil, errors.Wrapf(errors.ErrInvalidConfig, "unsupported operation type %s for Execution %s", e.Type, e.Name)
	}
	bound = e
	return bound, err
}

// Bind binds the TimeConfig by expanding environment variables.
func (t *TimeConfig) Bind() (bound *TimeConfig) {
	if t == nil {
		return nil
	}
	t.Delay = config.GetActualValue(t.Delay)
	t.Wait = config.GetActualValue(t.Wait)
	t.Timeout = config.GetActualValue(t.Timeout)
	return t
}

// Bind binds and validates the SearchQuery configuration.
func (sq *SearchQuery) Bind() (bound *SearchQuery, err error) {
	if sq == nil {
		return nil, errors.Wrap(errors.ErrInvalidConfig, "missing required fields on SearchQuery")
	}
	sq.Timeout = config.GetActualValue(sq.Timeout)
	dur, err := sq.Timeout.Duration()
	if err != nil || dur < 0 {
		sq.Timeout = defaultTimeout
	}
	sq.AlgorithmString = config.GetActualValue(sq.AlgorithmString)
	if sq.K == 0 {
		sq.K = defaultTopK
	}
	if sq.Radius == 0 {
		sq.Radius = -1
	}
	switch strings.TrimForCompare(sq.AlgorithmString) {
	case "concurrentqueue", "queue", "cqueue", "cq":
		sq.Algorithm = payload.Search_ConcurrentQueue
	case "sortslice", "slice", "sslice", "ss":
		sq.Algorithm = payload.Search_SortSlice
	case "sortpoolslice", "poolslice", "spslice", "pslice", "sps", "ps":
		sq.Algorithm = payload.Search_SortPoolSlice
	case "pairingheap", "pairheap", "pheap", "heap", "ph":
		sq.Algorithm = payload.Search_PairingHeap
	default:
		sq.Algorithm = payload.Search_ConcurrentQueue
	}
	return sq, nil
}

// Bind binds and validates the ModificationConfig.
func (m *ModificationConfig) Bind() (bound *ModificationConfig, err error) {
	if m.Timestamp < 0 {
		m.Timestamp = 0
	}
	return m, nil
}

// Bind expands environment variables for OperationType.
func (ot OperationType) Bind() (bound OperationType, err error) {
	if ot == "" {
		return "", errors.Wrap(errors.ErrInvalidConfig, "missing required fields on OperationType")
	}
	switch strings.TrimForCompare(config.GetActualValue(ot)) {
	case "search", "ser", "s":
		return OpSearch, nil
	case "searchbyid", "serid", "sid", "sbyid":
		return OpSearchByID, nil
	case "linearsearch", "lsearch", "lser", "ls":
		return OpLinearSearch, nil
	case "linearsearchbyid", "lsearchbyid", "lserid", "lsbyid":
		return OpLinearSearchByID, nil
	case "insert", "ins", "i":
		return OpInsert, nil
	case "update", "upd", "u":
		return OpUpdate, nil
	case "upsert", "usert", "upst", "us":
		return OpUpsert, nil
	case "remove", "rem", "r", "delete", "del", "d":
		return OpRemove, nil
	case "removebytimestamp", "removets", "remts", "rmts", "dts":
		return OpRemoveByTimestamp, nil
	case "object", "obj", "o":
		return OpObject, nil
	case "listobject", "listobj", "lobj", "lo":
		return OpListObject, nil
	case "timestamp", "ts", "t":
		return OpTimestamp, nil
	case "exists", "exist", "ex", "e":
		return OpExists, nil
	case "indexinfo", "index", "info", "ii":
		return OpIndexInfo, nil
	case "indexdetail", "detail", "id":
		return OpIndexDetail, nil
	case "indexstatistics", "statistics", "stat", "is":
		return OpIndexStatistics, nil
	case "indexstatisticsdetail", "statisticsdetail", "statdetail", "isd":
		return OpIndexStatisticsDetail, nil
	case "indexproperty", "property", "prop", "ip":
		return OpIndexProperty, nil
	case "flush", "fl", "f":
		return OpFlush, nil
	case "kubernetes", "kube", "k8s":
		return OpKubernetes, nil
	case "client", "cli", "c", "grpc":
		return OpClient, nil
	case "wait":
		return OpWait, nil
	}
	return bound, nil
}

// Bind expands environment variables for OperationMode.
func (om OperationMode) Bind() (bound OperationMode, err error) {
	switch strings.TrimForCompare(config.GetActualValue(om)) {
	case "unary", "un", "u":
		return OperationUnary, nil
	case "stream", "str", "s":
		return OperationStream, nil
	case "multiple", "multi", "m":
		return OperationMultiple, nil
	default:
		return OperationOther, nil
	}
}

// Bind expands environment variables for StatusCode.
func (op Operator) Bind() (bound Operator, err error) {
	switch strings.TrimForCompare(config.GetActualValue(op)) {
	case Le, Lt, Ge, Gt, Eq, Ne:
		return op, nil
	case "":
		return Eq, nil
	}
	return op, errors.New("Unsupported operator: " + string(op))
}

// Bind expands environment variables for StatusCode.
func (sc StatusCode) Bind() (bound StatusCode, err error) {
	switch strings.TrimForCompare(config.GetActualValue(sc)) {
	case StatusCodeOK:
		return StatusCodeOK, nil
	case StatusCodeCanceled:
		return StatusCodeCanceled, nil
	case StatusCodeUnknown:
		return StatusCodeUnknown, nil
	case StatusCodeInvalidArgument:
		return StatusCodeInvalidArgument, nil
	case StatusCodeDeadlineExceeded:
		return StatusCodeDeadlineExceeded, nil
	case StatusCodeNotFound:
		return StatusCodeNotFound, nil
	case StatusCodeAlreadyExists:
		return StatusCodeAlreadyExists, nil
	case StatusCodePermissionDenied:
		return StatusCodePermissionDenied, nil
	case StatusCodeResourceExhausted:
		return StatusCodeResourceExhausted, nil
	case StatusCodeFailedPrecondition:
		return StatusCodeFailedPrecondition, nil
	case StatusCodeAborted:
		return StatusCodeAborted, nil
	case StatusCodeOutOfRange:
		return StatusCodeOutOfRange, nil
	case StatusCodeUnimplemented:
		return StatusCodeUnimplemented, nil
	case StatusCodeInternal:
		return StatusCodeInternal, nil
	case StatusCodeUnavailable:
		return StatusCodeUnavailable, nil
	case StatusCodeDataLoss:
		return StatusCodeDataLoss, nil
	case StatusCodeUnauthenticated:
		return StatusCodeUnauthenticated, nil
	}
	return bound, nil
}

// Bind expands environment variables for KubernetesKind.
func (kk KubernetesKind) Bind() (bound KubernetesKind, err error) {
	switch strings.TrimForCompare(config.GetActualValue(kk)) {
	case "configmap", "config", "cm":
		return ConfigMap, nil
	case "cronjob", "cron", "cj":
		return CronJob, nil
	case "daemonset", "daemon", "ds":
		return DaemonSet, nil
	case "deployment", "deploy", "dep":
		return Deployment, nil
	case "job", "jb":
		return Job, nil
	case "pod", "pd":
		return Pod, nil
	case "secret", "sec":
		return Secret, nil
	case "service", "svc":
		return Service, nil
	case "statefulset", "stateful", "sts":
		return StatefulSet, nil
	}
	return bound, nil
}

// Bind expands environment variables for KubernetesAction.
func (ka KubernetesAction) Bind() (bound KubernetesAction, err error) {
	switch strings.TrimForCompare(config.GetActualValue(ka)) {
	case "rollout", "roll", "ro":
		return KubernetesActionRollout, nil
	case "delete", "del", "d":
		return KubernetesActionDelete, nil
	case "get", "g":
		return KubernetesActionGet, nil
	case "exec", "e":
		return KubernetesActionExec, nil
	case "apply", "a":
		return KubernetesActionApply, nil
	case "create", "c":
		return KubernetesActionCreate, nil
	case "patch", "p":
		return KubernetesActionPatch, nil
	case "scale", "s":
		return KubernetesActionScale, nil
	case "wait", "wa", "w":
		return KubernetesActionWait, nil
	}
	return bound, nil
}

// Bind expands environment variables for KubernetesAction.
func (ks KubernetesStatus) Bind() (bound KubernetesStatus, err error) {
	switch strings.TrimForCompare(config.GetActualValue(ks)) {
	case "unknown", "u":
		return KubernetesStatusUnknown, nil
	case "pending", "pen", "p":
		return KubernetesStatusPending, nil
	case "updating", "update":
		return KubernetesStatusUpdating, nil
	case "available", "a":
		return KubernetesStatusAvailable, nil
	case "degraded", "degrade", "d":
		return KubernetesStatusDegraded, nil
	case "failed", "fail", "f":
		return KubernetesStatusFailed, nil
	case "completed", "complete", "c":
		return KubernetesStatusCompleted, nil
	case "scheduled", "schedule", "sc":
		return KubernetesStatusScheduled, nil
	case "scaling", "scale", "s":
		return KubernetesStatusScaling, nil
	case "paused", "pause":
		return KubernetesStatusPaused, nil
	case "terminating", "terminate", "t":
		return KubernetesStatusTerminating, nil
	case "notready", "r":
		return KubernetesStatusNotReady, nil
	case "bound", "b":
		return KubernetesStatusBound, nil
	case "loadbalancing", "locabalance", "l":
		return KubernetesStatusLoadBalancing, nil
	}
	return bound, nil
}

// Bind binds and validates the KubernetesConfig.
func (k *KubernetesConfig) Bind() (bound *KubernetesConfig, err error) {
	if k == nil {
		return nil, errors.Wrap(errors.ErrInvalidConfig, "missing required fields on KubernetesConfig")
	}
	k.Namespace = config.GetActualValue(k.Namespace)
	k.Name = config.GetActualValue(k.Name)
	k.LabelSelector = config.GetActualValue(k.LabelSelector)
	if k.Action, err = k.Action.Bind(); err != nil {
		return nil, err
	}
	if k.Kind, err = k.Kind.Bind(); err != nil {
		return nil, err
	}
	if k.Status, err = k.Status.Bind(); err != nil {
		return nil, err
	}
	if k.Namespace == "" || (k.Name == "" && k.LabelSelector == "") || k.Action == "" || k.Kind == "" {
		return nil, errors.Errorf("kubernetes config: namespace: %s, name: %s or label_selector: %s, action: %s, and kind: %s must be provided",
			k.Namespace, k.Name, k.LabelSelector, k.Action, k.Kind)
	}
	if k.LabelSelector != "" {
		if k.Action != KubernetesActionWait {
			return nil, errors.Errorf("kubernetes config: label_selector is currently only supported for wait action")
		}
	}
	return k, nil
}

// Bind binds and validates the Kubernetes configuration.
func (k *Kubernetes) Bind() (bound *Kubernetes, err error) {
	if k == nil {
		return nil, errors.Wrap(errors.ErrInvalidConfig, "missing required fields on Kubernetes")
	}
	k.KubeConfig = config.GetActualValue(k.KubeConfig)
	if k.KubeConfig == "" {
		log.Warn("Kubernetes.KubeConfig is empty; please check your configuration")
	} else if !file.Exists(k.KubeConfig) {
		log.Warn("Kubernetes: kubeconfig file does not exist: ", k.KubeConfig)
	}
	if k.PortForward != nil {
		if k.PortForward, err = k.PortForward.Bind(); err != nil {
			return nil, err
		}
	}
	return k, nil
}

// Bind binds and validates the PortForward configuration.
func (pf *PortForward) Bind() (bound *PortForward, err error) {
	if pf == nil {
		return nil, errors.Wrap(errors.ErrInvalidConfig, "missing required fields on PortForward")
	}
	if !pf.Enabled {
		return pf, nil
	}
	pf.ServiceName = config.GetActualValue(pf.ServiceName)
	pf.Namespace = config.GetActualValue(pf.Namespace)
	if pf.ServiceName == "" {
		return nil, errors.New("portforward: service name cannot be empty")
	}
	if pf.Namespace == "" {
		return nil, errors.New("portforward: namespace cannot be empty")
	}
	if _, err = pf.TargetPort.Bind(); err != nil {
		return nil, err
	}
	if _, err = pf.LocalPort.Bind(); err != nil {
		return nil, err
	}
	if pf.TargetPort.Port() == 0 {
		pf.TargetPort = localPort
	}
	if pf.LocalPort.Port() == 0 {
		pf.LocalPort = localPort
	}
	return pf, nil
}

// Bind expands environment variables for Port.
func (p *Port) Bind() (bound *Port, err error) {
	port := config.GetActualValue(*p)
	return &port, nil
}

// Bind binds and validates the Dataset configuration.
func (d *Dataset) Bind() (bound *Dataset, err error) {
	if d == nil {
		return nil, errors.Wrap(errors.ErrInvalidConfig, "missing required fields on Dataset")
	}
	d.Name = config.GetActualValue(d.Name)
	if d.Name == "" || !file.Exists(d.Name) {
		return nil, errors.Errorf("dataset name: %s cannot be empty", d.Name)
	}
	return d, nil
}

// Bind binds and validates the Metrics configuration.
func (m *Metrics) Bind() (bound *Metrics, err error) {
	if m == nil || !m.Enabled {
		return nil, nil
	}
	if m.Histogram != nil {
		if h, err := m.Histogram.Bind(); err != nil {
			return nil, errors.Wrap(err, "failed to bind Histogram configuration")
		} else if h != nil {
			m.Histogram = h
		}
	}
	if m.LatencyHistogram != nil {
		if h, err := m.LatencyHistogram.Bind(); err != nil {
			return nil, errors.Wrap(err, "failed to bind LatencyHistogram configuration")
		} else if h != nil {
			m.LatencyHistogram = h
		}
	} else if m.Histogram != nil {
		m.LatencyHistogram = m.Histogram
	}
	if m.QueueWaitHistogram != nil {
		if h, err := m.QueueWaitHistogram.Bind(); err != nil {
			return nil, errors.Wrap(err, "failed to bind QueueWaitHistogram configuration")
		} else if h != nil {
			m.QueueWaitHistogram = h
		}
	} else if m.Histogram != nil {
		m.QueueWaitHistogram = m.Histogram
	}
	if m.TDigest != nil {
		if t, err := m.TDigest.Bind(); err != nil {
			return nil, errors.Wrap(err, "failed to bind TDigest configuration")
		} else if t != nil {
			m.TDigest = t
		}
	}
	if m.LatencyTDigest != nil {
		if t, err := m.LatencyTDigest.Bind(); err != nil {
			return nil, errors.Wrap(err, "failed to bind LatencyTDigest configuration")
		} else if t != nil {
			m.LatencyTDigest = t
		}
	} else if m.TDigest != nil {
		m.LatencyTDigest = m.TDigest
	}
	if m.QueueWaitTDigest != nil {
		if t, err := m.QueueWaitTDigest.Bind(); err != nil {
			return nil, errors.Wrap(err, "failed to bind QueueWaitTDigest configuration")
		} else if t != nil {
			m.QueueWaitTDigest = t
		}
	} else if m.TDigest != nil {
		m.QueueWaitTDigest = m.TDigest
	}
	if m.Exemplar != nil {
		if e, err := m.Exemplar.Bind(); err != nil {
			return nil, errors.Wrap(err, "failed to bind Exemplar configuration")
		} else if e != nil {
			m.Exemplar = e
		}
	}
	if m.RangeScales != nil {
		for i, rs := range m.RangeScales {
			if r, err := rs.Bind(); err != nil {
				return nil, errors.Wrapf(err, "failed to bind RangeScale configuration for %s", rs.Name)
			} else if r != nil {
				m.RangeScales[i] = r
			}
		}
	}
	if m.TimeScales != nil {
		for i, ts := range m.TimeScales {
			if t, err := ts.Bind(); err != nil {
				return nil, errors.Wrapf(err, "failed to bind TimeScale configuration for %s", ts.Name)
			} else if t != nil {
				m.TimeScales[i] = t
			}
		}
	}
	return m, nil
}

func (m *Metrics) Opts() (opts []metrics.Option) {
	if m == nil {
		return nil
	}
	opts = make([]metrics.Option, 0, len(m.RangeScales)+len(m.TimeScales)+6)
	if m.CustomCounters != nil {
		opts = append(opts, metrics.WithCustomCounters(m.CustomCounters...))
	}
	if m.LatencyHistogram != nil {
		hopts := make([]metrics.HistogramOption, 0, 4)
		if m.LatencyHistogram.NumShards > 0 {
			hopts = append(hopts, metrics.WithHistogramNumShards(m.LatencyHistogram.NumShards))
		}
		opts = append(opts, metrics.WithLatencyHistogram(hopts...))
	}
	if m.QueueWaitHistogram != nil {
		hopts := make([]metrics.HistogramOption, 0, 4)
		if m.QueueWaitHistogram.NumShards > 0 {
			hopts = append(hopts, metrics.WithHistogramNumShards(m.QueueWaitHistogram.NumShards))
		}
		opts = append(opts, metrics.WithQueueWaitHistogram(hopts...))
	}
	if m.LatencyTDigest != nil {
		opts = append(opts, metrics.WithLatencyTDigest(
			metrics.WithTDigestCompression(m.LatencyTDigest.Compression),
			metrics.WithTDigestCompressionTriggerFactor(m.LatencyTDigest.CompressionTriggerFactor),
			metrics.WithTDigestQuantiles(m.LatencyTDigest.Quantiles...),
			metrics.WithTDigestNumShards(m.LatencyTDigest.NumShards),
		))
	}
	if m.QueueWaitTDigest != nil {
		opts = append(opts, metrics.WithQueueWaitTDigest(
			metrics.WithTDigestCompression(m.QueueWaitTDigest.Compression),
			metrics.WithTDigestCompressionTriggerFactor(m.QueueWaitTDigest.CompressionTriggerFactor),
			metrics.WithTDigestQuantiles(m.QueueWaitTDigest.Quantiles...),
			metrics.WithTDigestNumShards(m.QueueWaitTDigest.NumShards),
		))
	}
	if m.Exemplar != nil {
		opts = append(opts, metrics.WithExemplar(
			metrics.WithExemplarCapacity(m.Exemplar.Capacity),
			metrics.WithExemplarNumShards(m.Exemplar.NumShards),
			metrics.WithExemplarSamplingRate(m.Exemplar.SamplingRate),
		))
	}
	opts = append(opts, metrics.WithDetailedErrorTracking(m.DetailedErrorTracking))
	if m.RangeScales != nil {
		for _, rs := range m.RangeScales {
			opts = append(opts, metrics.WithRangeScale(rs.Name, rs.Width, rs.Capacity))
		}
	}
	if m.TimeScales != nil {
		for _, ts := range m.TimeScales {
			opts = append(opts, metrics.WithTimeScale(ts.Name, time.Duration(ts.Width), ts.Capacity))
		}
	}
	return opts
}

// Bind binds and validates the Histogram configuration.
func (h *Histogram) Bind() (bound *Histogram, err error) {
	if h == nil {
		return nil, nil
	}
	return h, nil
}

// Bind binds and validates the TDigest configuration.
func (t *TDigest) Bind() (bound *TDigest, err error) {
	if t == nil {
		return nil, nil
	}
	return t, nil
}

// Bind binds and validates the Exemplar configuration.
func (e *Exemplar) Bind() (bound *Exemplar, err error) {
	if e == nil {
		return nil, nil
	}
	return e, nil
}

// Bind binds and validates the RangeScale configuration.
func (rs *RangeScale) Bind() (bound *RangeScale, err error) {
	if rs == nil {
		return nil, nil
	}
	rs.Name = config.GetActualValue(rs.Name)
	return rs, nil
}

// Bind binds and validates the TimeScale configuration.
func (ts *TimeScale) Bind() (bound *TimeScale, err error) {
	if ts == nil {
		return nil, nil
	}
	ts.Name = config.GetActualValue(ts.Name)
	return ts, nil
}

////////////////////////////////////////////////////////////////////////////////
// Func Section
////////////////////////////////////////////////////////////////////////////////

// Merge merges two Metrics configurations.
// The fields of `child` will override the fields of `parent` if they are set.
func (parent *Metrics) Merge(child *Metrics) (merged *Metrics) {
	if parent == nil && child == nil {
		return nil
	}
	if parent == nil {
		return child
	}
	if child == nil {
		return parent
	}

	merged = new(Metrics)

	if parent.Enabled ||
		(!parent.Enabled && child.Enabled) {
		merged.Enabled = true
	}

	if child.Histogram != nil {
		merged.Histogram = child.Histogram
	} else {
		merged.Histogram = parent.Histogram
	}
	if child.LatencyHistogram != nil {
		merged.LatencyHistogram = child.LatencyHistogram
	} else {
		merged.LatencyHistogram = parent.LatencyHistogram
	}
	if child.QueueWaitHistogram != nil {
		merged.QueueWaitHistogram = child.QueueWaitHistogram
	} else {
		merged.QueueWaitHistogram = parent.QueueWaitHistogram
	}
	if child.TDigest != nil {
		merged.TDigest = child.TDigest
	} else {
		merged.TDigest = parent.TDigest
	}
	if child.LatencyTDigest != nil {
		merged.LatencyTDigest = child.LatencyTDigest
	} else {
		merged.LatencyTDigest = parent.LatencyTDigest
	}
	if child.QueueWaitTDigest != nil {
		merged.QueueWaitTDigest = child.QueueWaitTDigest
	} else {
		merged.QueueWaitTDigest = parent.QueueWaitTDigest
	}
	if child.Exemplar != nil {
		merged.Exemplar = child.Exemplar
	} else {
		merged.Exemplar = parent.Exemplar
	}
	if child.RangeScales != nil {
		merged.RangeScales = child.RangeScales
	} else {
		merged.RangeScales = parent.RangeScales
	}
	if child.TimeScales != nil {
		merged.TimeScales = child.TimeScales
	} else {
		merged.TimeScales = parent.TimeScales
	}
	if child.CustomCounters != nil {
		merged.CustomCounters = child.CustomCounters
	} else {
		merged.CustomCounters = parent.CustomCounters
	}

	if parent.DetailedErrorTracking ||
		(!parent.DetailedErrorTracking && child.DetailedErrorTracking) {
		merged.DetailedErrorTracking = true
	}

	return merged
}

// Timing interface provides access to time configuration values.
type Timing interface {
	GetDelay() timeutil.DurationString
	GetWait() timeutil.DurationString
	GetTimeout() timeutil.DurationString
}

// GetDelay returns the Delay value from TimeConfig.
func (t *TimeConfig) GetDelay() timeutil.DurationString {
	if t == nil {
		return ""
	}
	return t.Delay
}

// GetWait returns the Wait value from TimeConfig.
func (t *TimeConfig) GetWait() timeutil.DurationString {
	if t == nil {
		return ""
	}
	return t.Wait
}

// GetTimeout returns the Timeout value from TimeConfig.
func (t *TimeConfig) GetTimeout() timeutil.DurationString {
	if t == nil {
		return ""
	}
	return t.Timeout
}

type Repeater interface {
	GetRepeats() *Repeats
}

func (d Data) GetRepeats() *Repeats {
	return &Repeats{} // Data level repetition is not supported. Use Strategy, Operation, or Execution level repetition instead, as these levels are designed to handle repeated operations.
}

func (s Strategy) GetRepeats() *Repeats {
	return s.Repeats
}

func (o Operation) GetRepeats() *Repeats {
	return o.Repeats
}

func (e Execution) GetRepeats() *Repeats {
	return e.Repeats
}

// Equals compares StatusCode with a given string ignoring case.
func (sc StatusCode) Equals(c string) bool {
	bound, _ := sc.Bind() // error ignored as Bind never errors
	return strings.EqualFold(bound.String(), c)
}

func (sc StatusCode) Status() codes.Code {
	switch strings.TrimForCompare(sc) {
	case StatusCodeOK:
		return codes.OK
	case StatusCodeCanceled:
		return codes.Canceled
	case StatusCodeUnknown:
		return codes.Unknown
	case StatusCodeInvalidArgument:
		return codes.InvalidArgument
	case StatusCodeDeadlineExceeded:
		return codes.DeadlineExceeded
	case StatusCodeNotFound:
		return codes.NotFound
	case StatusCodeAlreadyExists:
		return codes.AlreadyExists
	case StatusCodePermissionDenied:
		return codes.PermissionDenied
	case StatusCodeResourceExhausted:
		return codes.ResourceExhausted
	case StatusCodeFailedPrecondition:
		return codes.FailedPrecondition
	case StatusCodeAborted:
		return codes.Aborted
	case StatusCodeOutOfRange:
		return codes.OutOfRange
	case StatusCodeUnimplemented:
		return codes.Unimplemented
	case StatusCodeInternal:
		return codes.Internal
	case StatusCodeUnavailable:
		return codes.Unavailable
	case StatusCodeDataLoss:
		return codes.DataLoss
	case StatusCodeUnauthenticated:
		return codes.Unauthenticated
	}
	return codes.Unknown
}

// String returns the string representation of StatusCode.
func (sc StatusCode) String() string {
	return string(sc)
}

// Port returns the numeric value of the Port.
func (p Port) Port() uint16 {
	bp, _ := p.Bind() // error ignored as Bind never errors
	port, err := strconv.ParseUint(string(*bp), 10, 16)
	if err != nil {
		return 0
	}
	return uint16(port)
}

func (ks KubernetesStatus) Status() kubernetes.ResourceStatus {
	switch strings.TrimForCompare(ks) {
	case KubernetesStatusUnknown:
		return kubernetes.StatusUnknown
	case KubernetesStatusPending:
		return kubernetes.StatusPending
	case KubernetesStatusUpdating:
		return kubernetes.StatusUpdating
	case KubernetesStatusAvailable:
		return kubernetes.StatusAvailable
	case KubernetesStatusDegraded:
		return kubernetes.StatusDegraded
	case KubernetesStatusFailed:
		return kubernetes.StatusFailed
	case KubernetesStatusCompleted:
		return kubernetes.StatusCompleted
	case KubernetesStatusScheduled:
		return kubernetes.StatusScheduled
	case KubernetesStatusScaling:
		return kubernetes.StatusScaling
	case KubernetesStatusPaused:
		return kubernetes.StatusPaused
	case KubernetesStatusTerminating:
		return kubernetes.StatusTerminating
	case KubernetesStatusNotReady:
		return kubernetes.StatusNotReady
	case KubernetesStatusBound:
		return kubernetes.StatusBound
	case KubernetesStatusLoadBalancing:
		return kubernetes.StatusLoadBalancing
	}
	return kubernetes.StatusUnknown
}

// //////////////////////////////////////////////////////////////////////////////
// Const Section
// //////////////////////////////////////////////////////////////////////////////
const (
	localPort      Port = "8081"
	defaultTopK         = uint32(10)
	defaultTimeout      = timeutil.DurationString("3s")
)

func replaceEnvInValues(v any) any {
	switch val := v.(type) {
	case string:
		str := config.GetActualValue(val)
		// Convert env-injected scalars to native types so yaml.Unmarshal can bind to typed fields.
		// Preference order: big uints (for positive numbers), signed ints (for negatives), then floats; handle true/false explicitly (not 0/1).
		if len(str) > 0 && str[0] != '-' {
			if u, err := strconv.ParseUint(str, 10, 64); err == nil {
				return u
			}
		}
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			return i
		}
		if f, err := strconv.ParseFloat(str, 64); err == nil {
			return f
		}
		switch strings.ToLower(str) {
		case "true":
			return true
		case "false":
			return false
		}
		return str
	case []any:
		for i, e := range val {
			val[i] = replaceEnvInValues(e)
		}
		return val
	case map[string]any:
		for k, e := range val {
			val[k] = replaceEnvInValues(e)
		}
		return val
	default:
		return val
	}
}

// TODO: This function is copied from the internal/config package and modified to mitigate the risk. We need to merge this back to the internal/config package
// when possible.
// read returns config struct or error when decoding the configuration file to actually *Config struct.
func read[T any](path string, cfg T) (err error) {
	f, err := file.Open(path, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		if f != nil {
			if err != nil {
				err = errors.Join(f.Close(), err)
				return
			}
			err = f.Close()
		}
	}()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	var raw map[string]any
	switch ext := filepath.Ext(path); ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return err
		}
	case ".json":
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}
	default:
		return errors.ErrUnsupportedConfigFileType(ext)
	}
	replaced := replaceEnvInValues(raw)
	intermediate, err := yaml.Marshal(replaced)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(intermediate, cfg)
}

////////////////////////////////////////////////////////////////////////////////
// Load Function
////////////////////////////////////////////////////////////////////////////////

// Load reads the configuration from the specified file path,
// binds and validates the configuration, and returns the complete Data configuration.
func Load(path string) (cfg *Data, err error) {
	log.Debugf("loading test client configuration from %s", path)
	cfg = new(Data)
	if err = read(path, cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to read configuration from %s", path)
	}
	if cfg == nil || len(cfg.Strategies) == 0 || cfg.Dataset == nil {
		return nil, errors.Errorf("failed to load configuration from %s", path)
	}
	if cfg, err = cfg.Bind(); err != nil {
		return nil, errors.Wrapf(err, "failed to bind configuration from %s", path)
	}
	log.Debug(config.ToRawYaml(cfg))
	return cfg, nil
}
