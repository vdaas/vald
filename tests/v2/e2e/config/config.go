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
// This file includes refactored Bind methods (always returning error) and non-Bind functions,
// with named return values and proper ordering of sections.
package config

import (
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/timeutil"
	"github.com/vdaas/vald/internal/timeutil/rate"
	"github.com/vdaas/vald/tests/v2/e2e/kubernetes"
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
	Metadata            map[string]string  `json:"metadata,omitempty"        yaml:"metadata,omitempty"`
	MetaString          string             `json:"metadata_string,omitempty" yaml:"metadata_string,omitempty"`
	FilePath            string             `json:"-"                         yaml:"-"`
}

// Strategy represents a test strategy.
type Strategy struct {
	TimeConfig  `             yaml:",inline"              json:",inline"`
	Name        string       `yaml:"name"                 json:"name,omitempty"`
	Concurrency uint64       `yaml:"concurrency"          json:"concurrency,omitempty"`
	Operations  []*Operation `yaml:"operations,omitempty" json:"operations,omitempty"`
}

// Operation represents an individual operation configuration.
type Operation struct {
	TimeConfig `             yaml:",inline"              json:",inline"`
	Name       string       `yaml:"name,omitempty"       json:"name,omitempty"`
	Executions []*Execution `yaml:"executions,omitempty" json:"executions,omitempty"`
}

// Execution represents the execution details for a given operation.
type Execution struct {
	*BaseConfig  `                    yaml:",inline,omitempty"      json:",inline,omitempty"`
	TimeConfig   `                    yaml:",inline"                json:",inline"`
	Name         string              `yaml:"name"                   json:"name,omitempty"`
	Type         OperationType       `yaml:"type"                   json:"type,omitempty"`
	Mode         OperationMode       `yaml:"mode"                   json:"mode,omitempty"`
	Search       *SearchQuery        `yaml:"search,omitempty"       json:"search,omitempty"`
	Kubernetes   *KubernetesConfig   `yaml:"kubernetes,omitempty"   json:"kubernetes,omitempty"`
	Modification *ModificationConfig `yaml:"modification,omitempty" json:"modification,omitempty"`
	Expect       []Expect            `yaml:"expect,omitempty"       json:"expect,omitempty"`
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
	// Bind each Strategy.
	var cnt int
	for _, strategy := range d.Strategies {
		if strategy != nil {
			var bs *Strategy
			if bs, err = strategy.Bind(); err != nil {
				return nil, errors.Wrapf(err, "failed to bind strategy: %s", strategy.Name)
			} else if bs != nil {
				d.Strategies[cnt] = bs
				cnt++
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
	return d, nil
}

// Bind binds and validates the Strategy configuration.
func (s *Strategy) Bind() (bound *Strategy, err error) {
	if s == nil || s.Operations == nil || len(s.Operations) == 0 {
		return nil, errors.Wrapf(errors.ErrInvalidConfig, "missing required fields on Strategy %s", s.Name)
	}
	s.Name = config.GetActualValue(s.Name)
	s.TimeConfig.Bind()
	var cnt int
	for _, op := range s.Operations {
		if op != nil {
			var bo *Operation
			if bo, err = op.Bind(); err != nil {
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
func (o *Operation) Bind() (bound *Operation, err error) {
	if o == nil || o.Executions == nil || len(o.Executions) == 0 {
		return nil, errors.Wrapf(errors.ErrInvalidConfig, "missing required fields on Operation %s", o.Name)
	}
	o.Name = config.GetActualValue(o.Name)
	o.TimeConfig.Bind()
	var cnt int
	for _, exec := range o.Executions {
		var be *Execution
		if be, err = exec.Bind(); err != nil {
			return nil, errors.Wrapf(err, "failed to bind execution: %s", exec.Name)
		} else if be != nil {
			o.Executions[cnt] = be
			cnt++
		}
	}
	return o, nil
}

// Bind binds and validates the Execution configuration.
func (e *Execution) Bind() (bound *Execution, err error) {
	if e == nil {
		return nil, errors.Wrap(errors.ErrInvalidConfig, "missing required fields on Execution")
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
	return
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
	switch trimStringForCompare(sq.AlgorithmString) {
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
	switch trimStringForCompare(config.GetActualValue(ot)) {
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
	switch trimStringForCompare(config.GetActualValue(om)) {
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
	switch trimStringForCompare(config.GetActualValue(op)) {
	case Le, Lt, Ge, Gt, Eq, Ne:
		return op, nil
	case "":
		return Eq, nil
	}
	return op, errors.New("Unsupported operator: " + string(op))
}

// Bind expands environment variables for StatusCode.
func (sc StatusCode) Bind() (bound StatusCode, err error) {
	switch trimStringForCompare(config.GetActualValue(sc)) {
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
	switch trimStringForCompare(config.GetActualValue(kk)) {
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
	switch trimStringForCompare(config.GetActualValue(ka)) {
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
	switch trimStringForCompare(config.GetActualValue(ks)) {
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

////////////////////////////////////////////////////////////////////////////////
// Func Section
////////////////////////////////////////////////////////////////////////////////

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

// Equals compares StatusCode with a given string ignoring case.
func (sc StatusCode) Equals(c string) bool {
	bound, _ := sc.Bind() // error ignored as Bind never errors
	return strings.EqualFold(bound.String(), c)
}

func (sc StatusCode) Status() codes.Code {
	switch trimStringForCompare(sc) {
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
	switch trimStringForCompare(ks) {
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

////////////////////////////////////////////////////////////////////////////////
// Load Function
////////////////////////////////////////////////////////////////////////////////

// Load reads the configuration from the specified file path,
// binds and validates the configuration, and returns the complete Data configuration.
func Load(path string) (cfg *Data, err error) {
	log.Debugf("loading test client configuration from %s", path)
	cfg = new(Data)
	if err = config.Read(path, &cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to read configuration from %s", path)
	}
	if cfg == nil {
		return nil, errors.Errorf("failed to load configuration from %s", path)
	}
	if cfg, err = cfg.Bind(); err != nil {
		return nil, errors.Wrapf(err, "failed to bind configuration from %s", path)
	}
	log.Debug(config.ToRawYaml(cfg))
	return cfg, nil
}

var reps = strings.NewReplacer(" ", "", "-", "", "_", "", ":", "", ";", "", ",", "", ".", "")

func trimStringForCompare[S ~string](str S) S {
	return S(reps.Replace(string(str)))
}
