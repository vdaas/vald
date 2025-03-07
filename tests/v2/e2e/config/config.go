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
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/timeutil"
)

// Data represents the complete configuration for the application.
// It encapsulates all configuration sections, including gRPC target, search settings, operation settings,
// Kubernetes settings, dataset details, and additional metadata.
type Data struct {
	config.GlobalConfig `json:",inline,omitempty" yaml:",inline,omitempty"`
	Target              *config.GRPCClient `json:"target,omitempty"          yaml:"target,omitempty"`          // gRPC target configuration.
	Strategies          []*Strategy        `json:"strategies,omitempty"      yaml:"strategies,omitempty"`      // test strategies
	Dataset             *Dataset           `json:"dataset,omitempty"         yaml:"dataset,omitempty"`         // Dataset configuration.
	Kubernetes          *Kubernetes        `json:"kubernetes,omitempty"      yaml:"kubernetes,omitempty"`      // Kubernetes-related configuration.
	Metadata            map[string]string  `json:"metadata,omitempty"        yaml:"metadata,omitempty"`        // Additional metadata provided as key-value pairs.
	MetaString          string             `json:"metadata_string,omitempty" yaml:"metadata_string,omitempty"` // Raw metadata string (e.g., "KEY1=VAL1,KEY2=VAL2") to be parsed.
}

// Strateguy represents a test strategy that includes a slice of operations to be executed
// the operations are executed in concurrent goroutines with the specified delay between them.
type Strategy struct {
	TimeConfig  `             yaml:",inline,omitempty"    json:",inline,omitempty"`
	Name        string       `yaml:"name"                 json:"name,omitempty"` // Name of the strategy.
	Concurrency uint64       `yaml:"concurrency"          json:"concurrency,omitempty"`
	Operations  []*Operation `yaml:"operations,omitempty" json:"operations,omitempty"`
}

type Operation struct {
	TimeConfig `             yaml:",inline,omitempty"    json:",inline,omitempty"`
	Name       string       `yaml:"name,omitempty"       json:"name,omitempty"`
	Executions []*Execution `yaml:"executions,omitempty" json:"executions,omitempty"`
}

type Execution struct {
	*BaseConfig         `               yaml:",inline,omitempty"               json:",inline,omitempty"`
	*ModificationConfig `               yaml:",inline,omitempty"               json:",inline,omitempty"`
	*KubernetesConfig   `               yaml:",inline,omitempty"               json:",inline,omitempty"`
	TimeConfig          `               yaml:",inline,omitempty"               json:",inline,omitempty"`
	Name                string         `yaml:"name"                            json:"name,omitempty"` // Name of the execution.
	Type                OperationType  `yaml:"type"                            json:"type,omitempty"`
	Mode                OperationMode  `yaml:"mode"                            json:"mode,omitempty"`
	ExpectedStatusCodes StatusCodes    `yaml:"expected_status_codes,omitempty" json:"expected_status_codes,omitempty"`
	SearchConfig        []*SearchQuery `yaml:"search_config,omitempty"         json:"search_config,omitempty"`
}

type TimeConfig struct {
	Delay   timeutil.DurationString `yaml:"delay"   json:"delay,omitempty"`
	Wait    timeutil.DurationString `yaml:"wait"    json:"wait,omitempty"`
	Timeout timeutil.DurationString `yaml:"timeout" json:"timeout,omitempty"`
}
type Timing interface {
	GetDelay() timeutil.DurationString
	GetWait() timeutil.DurationString
	GetTimeout() timeutil.DurationString
}

func (t *TimeConfig) GetDelay() timeutil.DurationString {
	if t == nil {
		return ""
	}
	return t.Delay
}

func (t *TimeConfig) GetWait() timeutil.DurationString {
	if t == nil {
		return ""
	}
	return t.Wait
}

func (t *TimeConfig) GetTimeout() timeutil.DurationString {
	if t == nil {
		return ""
	}
	return t.Timeout
}

type BaseConfig struct {
	Num         uint64 `yaml:"num,omitempty"         json:"num,omitempty"`         // Number of items to process.
	Offset      uint64 `yaml:"offset,omitempty"      json:"offset,omitempty"`      // Starting offset for the operation.
	BulkSize    uint64 `yaml:"bulk_size,omitempty"   json:"bulk_size,omitempty"`   // Bulk size for multi-xxx operations.
	Parallelism uint64 `yaml:"parallelism,omitempty" json:"parallelism,omitempty"` // Parallelism for operations.
}

// SearchQuery represents the detailed parameters for a single search query.
type SearchQuery struct {
	K               uint32                              `yaml:"k,omitempty"         json:"k,omitempty"`                // Number of top results to return.
	Radius          float32                             `yaml:"radius,omitempty"    json:"radius,omitempty"`           // Radius for search (if applicable).
	Epsilon         float32                             `yaml:"epsilon,omitempty"   json:"epsilon,omitempty"`          // Epsilon for approximate search algorithms.
	AlgorithmString string                              `yaml:"algorithm,omitempty" json:"algorithm_string,omitempty"` // Algorithm identifier as a string; will be normalized and mapped.
	MinNum          uint32                              `yaml:"min_num,omitempty"   json:"min_num,omitempty"`          // Minimum number of items required for the operation.
	Ratio           float32                             `yaml:"ratio,omitempty"     json:"ratio,omitempty"`            // Ratio parameter for search (algorithm dependent).
	Nprobe          uint32                              `yaml:"nprobe,omitempty"    json:"nprobe,omitempty"`           // Number of probes for the search algorithm.
	Timeout         timeutil.DurationString             `yaml:"timeout,omitempty"   json:"timeout,omitempty"`          // Timeout value as a time.Duration.
	Algorithm       payload.Search_AggregationAlgorithm `yaml:"-"                   json:"-"`                          // Mapped algorithm constant based on AlgorithmString.
}

// Setting represents basic operation settings used across multiple operations (e.g., insert, update).
// It includes numeric values such as the number of items to process, an offset, and a timestamp.
type ModificationConfig struct {
	SkipStrictExistCheck bool  `yaml:"skip_strict_exist_check,omitempty" json:"skip_strict_exist_check,omitempty"` // Flag to indicate if strict existence checks should be skipped.
	Timestamp            int64 `yaml:"timestamp,omitempty"               json:"timestamp,omitempty"`               // Timestamp value for the operation; used for versioning.
}

type (
	StatusCode  string
	StatusCodes []StatusCode
)

func (sc StatusCode) Bind() StatusCode {
	return config.GetActualValue(sc)
}

func (sc StatusCode) Equals(c string) bool {
	return strings.EqualFold(sc.String(), StatusCode(c).Bind().String())
}

func (sc StatusCode) String() string {
	return string(sc)
}

func (sc StatusCodes) Bind() StatusCodes {
	for i, c := range sc {
		sc[i] = c.Bind()
	}
	return sc
}

func (sc StatusCodes) Equals(c string) bool {
	for _, s := range sc {
		if s.Equals(c) {
			return true
		}
	}
	return false
}

type KubernetesConfig struct {
	Action    KubernetesAction       `yaml:"action"    json:"action,omitempty"`
	Namespace string                 `yaml:"namespace" json:"namespace,omitempty"`
	Name      string                 `yaml:"name"      json:"name,omitempty"`
	Args      map[any]any            `yaml:"args"      json:"args,omitempty"`
	Resource  KubernetesResourceType `yaml:"resource"  json:"resource,omitempty"`
}

type Port string

func (p *Port) Bind() *Port {
	bp := config.GetActualValue(*p)
	p = &bp
	return p
}

func (p Port) Port() uint16 {
	port, err := strconv.ParseUint(string(*p.Bind()), 10, 16)
	if err != nil {
		return 0
	}
	return uint16(port)
}

// Kubernetes holds configuration settings specific to Kubernetes environments.
type Kubernetes struct {
	KubeConfig  string       `yaml:"kubeconfig"            json:"kube_config,omitempty"`  // File path to the kubeconfig.
	PortForward *PortForward `yaml:"portforward,omitempty" json:"port_forward,omitempty"` // Port forwarding settings.
}

// PortForward holds configuration for port forwarding when running in a Kubernetes environment.
type PortForward struct {
	Enabled     bool   `yaml:"enabled"      json:"enabled,omitempty"`      // Flag to enable or disable port forwarding.
	TargetPort  Port   `yaml:"target_port"  json:"target_port,omitempty"`  // The port forward target port number.
	LocalPort   Port   `yaml:"local_port"   json:"local_port,omitempty"`   // The local port number; if not set, it defaults to TargetPort.
	Namespace   string `yaml:"namespace"    json:"namespace,omitempty"`    // The Kubernetes namespace of the pod.
	ServiceName string `yaml:"service_name" json:"service_name,omitempty"` // The Kubernetes service name of the pod.
}

// Dataset holds information about the dataset to be used, such as the filename.
type Dataset struct {
	Name string `yaml:"name" json:"name,omitempty"` // Name (or path) of the dataset file.
}

// Bind processes and validates the Data configuration.
// It calls Bind on all nested configurations and processes the metadata fields.
// This method ensures that no field is omitted during the binding process.
func (d *Data) Bind() *Data {
	// Process the gRPC Target configuration if provided.
	if d.Target != nil {
		d.Target.Bind()
	}

	if d.Strategies != nil {
		for i, strategy := range d.Strategies {
			d.Strategies[i] = strategy.Bind()
		}
	}

	// Process the dataset configuration.
	if d.Dataset != nil {
		d.Dataset.Bind()
	}

	// Process Kubernetes configuration.
	if d.Kubernetes != nil {
		d.Kubernetes.Bind()
	}

	// Process metadata.
	// Expand any environment variables or placeholders in the raw metadata string.
	d.MetaString = config.GetActualValue(d.MetaString)
	// Initialize the Metadata map if it is nil.
	if d.Metadata == nil {
		d.Metadata = make(map[string]string)
	}
	// Split the MetaString into individual key=value pairs and add them to the Metadata map.
	for _, meta := range strings.Split(d.MetaString, ",") {
		key, val, ok := strings.Cut(meta, "=")
		if ok && key != "" && val != "" {
			d.Metadata[key] = val
		}
	}
	// Re-evaluate each key and value in the Metadata map to ensure all dynamic values are updated.
	for key, val := range d.Metadata {
		newKey := config.GetActualValue(key)
		val = config.GetActualValue(val)
		if key != newKey {
			delete(d.Metadata, key)
		}
		d.Metadata[newKey] = val
	}

	return d
}

func (s *Strategy) Bind() *Strategy {
	if s == nil {
		return nil
	}

	s.Delay = config.GetActualValue(s.Delay)
	s.Timeout = config.GetActualValue(s.Timeout)
	s.Wait = config.GetActualValue(s.Wait)

	if s.Operations != nil {
		for i, operation := range s.Operations {
			s.Operations[i] = operation.Bind()
		}
	}

	return s
}

func (o *Operation) Bind() *Operation {
	if o == nil {
		return nil
	}

	o.Delay = config.GetActualValue(o.Delay)
	o.Timeout = config.GetActualValue(o.Timeout)
	o.Wait = config.GetActualValue(o.Wait)

	if o.Executions != nil {
		for i := range o.Executions {
			o.Executions[i].Bind()
		}
	}

	return o
}

func (e *Execution) Bind() *Execution {
	if e == nil {
		return nil
	}

	e.Delay = config.GetActualValue(e.Delay)
	e.Wait = config.GetActualValue(e.Wait)
	e.Timeout = config.GetActualValue(e.Timeout)

	if e.BaseConfig != nil {
		e.BaseConfig.Bind()
	}
	if e.ModificationConfig != nil {
		e.ModificationConfig.Bind()
	}
	if e.KubernetesConfig != nil {
		e.KubernetesConfig.Bind()
	}
	if e.SearchConfig != nil {
		for i := range e.SearchConfig {
			e.SearchConfig[i].Bind()
		}
	}

	return e
}

func (b *BaseConfig) Bind() *BaseConfig {
	if b == nil {
		return nil
	}

	return b
}

// Bind validates and processes the SearchQuery parameters.
// It parses the timeout string into a time.Duration and maps the algorithm string to a constant.
func (sq *SearchQuery) Bind() *SearchQuery {
	if sq == nil {
		return nil
	}
	// Use a default timeout if parsing fails.
	sq.Timeout = config.GetActualValue(sq.Timeout)

	// Expand and normalize the algorithm string.
	sq.AlgorithmString = config.GetActualValue(sq.AlgorithmString)
	// Normalize by converting to lowercase and trimming common separators.
	switch strings.TrimFunc(strings.ToLower(sq.AlgorithmString), func(r rune) bool {
		switch r {
		case ' ', '-', '_', ':', ',':
			return true
		}
		return false
	}) {
	case "concurrentqueue", "queue", "cqueue", "cq":
		sq.Algorithm = payload.Search_ConcurrentQueue
	case "sortslice", "slice", "sslice", "ss":
		sq.Algorithm = payload.Search_SortSlice
	case "sortpoolslice", "poolslice", "spslice", "pslice", "sps", "ps":
		sq.Algorithm = payload.Search_SortPoolSlice
	case "pairingheap", "pairheap", "pheap", "heap":
		sq.Algorithm = payload.Search_PairingHeap
	default:
		// Default to ConcurrentQueue if the algorithm string is unrecognized.
		sq.Algorithm = payload.Search_ConcurrentQueue
	}
	return sq
}

// Bind processes and validates the Dataset configuration by expanding environment variables.
func (d *Dataset) Bind() *Dataset {
	if d == nil {
		return nil
	}
	// Expand any dynamic values in the dataset name.
	d.Name = config.GetActualValue(d.Name)
	return d
}

// Bind validates and processes the modification configuration for operations such as insert, update, etc.
// It logs warnings if certain numeric values are zero or invalid.
func (s *ModificationConfig) Bind() *ModificationConfig {
	// If the timestamp is negative, reset it to zero and log a warning.
	if s.Timestamp < 0 {
		log.Warn("ModificationConfig.Timestamp is negative, resetting to 0")
		s.Timestamp = 0
	}
	return s
}

func (k *KubernetesConfig) Bind() *KubernetesConfig {
	if k == nil {
		return nil
	}

	k.Namespace = config.GetActualValue(k.Namespace)
	k.Name = config.GetActualValue(k.Name)
	k.Action = config.GetActualValue(k.Action)
	k.Resource = config.GetActualValue(k.Resource)

	return k
}

// Bind processes Kubernetes configuration by validating the kubeconfig file path,
// expanding environment variables, and binding nested port forwarding settings.
func (k *Kubernetes) Bind() *Kubernetes {
	if k == nil {
		return nil
	}
	// Expand the kubeconfig path to replace any environment variables.
	k.KubeConfig = config.GetActualValue(k.KubeConfig)
	if k.KubeConfig == "" {
		log.Warn("Kubernetes.KubeConfig is empty; please check your configuration")
	} else if !file.Exists(k.KubeConfig) {
		// Warn if the specified kubeconfig file does not exist.
		log.Warn("Kubernetes: kubeconfig file does not exist: ", k.KubeConfig)
	}
	// Bind the PortForward configuration if it is provided.
	if k.PortForward != nil {
		k.PortForward.Bind()
	}
	return k
}

// Bind validates and processes the PortForward configuration.
// It expands environment variables for PodName and Namespace and sets default port values if necessary.
func (pf *PortForward) Bind() *PortForward {
	if pf == nil {
		return nil
	}
	// Expand dynamic values in the pod name and namespace.
	pf.ServiceName = config.GetActualValue(pf.ServiceName)
	pf.Namespace = config.GetActualValue(pf.Namespace)

	pf.TargetPort.Bind()
	pf.LocalPort.Bind()

	// If TargetPort is not set, default it to the localPort constant.
	if pf.TargetPort.Port() == 0 {
		pf.TargetPort = localPort
	}
	// If LocalPort is not set, default it to the same value as TargetPort.
	if pf.LocalPort.Port() == 0 {
		pf.LocalPort = localPort
	}
	return pf
}

// Load reads the configuration from the specified file path.
// If reading fails, it merges the read configuration with the default configuration.
// Finally, it calls Bind to perform all necessary post-processing on the configuration.
func Load(path string) (cfg *Data, err error) {
	log.Debugf("loading test client configuration from %s", path)
	cfg = new(Data)

	// Attempt to read the configuration from the file.
	err = config.Read(path, &cfg)
	if err != nil {
		log.Error(err)
		// If reading fails, merge the configuration with default values.
		cfg, err = config.Merge(cfg, Default)
		if err != nil {
			return nil, err
		}
	}
	log.Info(config.ToRawYaml(cfg))

	if cfg != nil {
		// Process and validate all configuration values.
		cfg.Bind()
	} else {
		return nil, errors.ErrInvalidConfig
	}

	// Log the entire configuration as raw YAML for debugging purposes.
	log.Debug(config.ToRawYaml(cfg))
	return cfg, nil
}
