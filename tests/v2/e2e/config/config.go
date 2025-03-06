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
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/timeutil"
)

// Data represents the complete configuration for the application.
// It encapsulates all configuration sections, including gRPC target, search settings, operation settings,
// Kubernetes settings, dataset details, and additional metadata.
type Data struct {
	Target           *config.GRPCClient `yaml:"target" json:"target,omitempty"`                           // gRPC target configuration.
	Search           *SearchConfig      `yaml:"search" json:"search,omitempty"`                           // Configuration for search operations.
	SearchByID       *SearchConfig      `yaml:"search_by_id" json:"search_by_id,omitempty"`               // Configuration for search-by-id operations.
	LinearSearch     *SearchConfig      `yaml:"linear_search" json:"linear_search,omitempty"`             // Configuration for linear search operations.
	LinearSearchByID *SearchConfig      `yaml:"linear_search_by_id" json:"linear_search_by_id,omitempty"` // Configuration for linear search-by-id operations.
	Insert           *Setting           `yaml:"insert" json:"insert,omitempty"`                           // Configuration for insert operations.
	Update           *Setting           `yaml:"update" json:"update,omitempty"`                           // Configuration for update operations.
	Upsert           *Setting           `yaml:"upsert" json:"upsert,omitempty"`                           // Configuration for upsert operations.
	Remove           *Setting           `yaml:"remove" json:"remove,omitempty"`                           // Configuration for remove operations.
	Object           *Setting           `yaml:"object" json:"object,omitempty"`                           // Configuration for object retrieval.
	Index            *WaitAfterInsert   `yaml:"index" json:"index,omitempty"`                             // Configuration for waiting period after insert.
	Dataset          *Dataset           `yaml:"dataset" json:"dataset,omitempty"`                         // Dataset configuration.
	Kubernetes       *Kubernetes        `yaml:"kubernetes" json:"kubernetes,omitempty"`                   // Kubernetes-related configuration.
	Metadata         map[string]string  `yaml:"metadata" json:"metadata,omitempty"`                       // Additional metadata provided as key-value pairs.
	MetaString       string             `yaml:"metadata_string" json:"meta_string,omitempty"`             // Raw metadata string (e.g., "KEY1=VAL1,KEY2=VAL2") to be parsed.
}

// SearchConfig holds configuration parameters specific to search operations.
// It defines the total number of items, an offset for pagination, and a slice of detailed search queries.
type SearchConfig struct {
	Num         uint64         `yaml:"num" json:"num,omitempty"`                 // Total number of items to be used for search.
	Offset      uint64         `yaml:"offset" json:"offset,omitempty"`           // Starting offset for the search operation.
	BulkSize    uint64         `yaml:"bulk_size" json:"bulk_size,omitempty"`     // Bulk size for multi-search operations.
	Concurrency uint64         `yaml:"concurrency" json:"concurrency,omitempty"` // Concurrency for search operations.
	Queries     []*SearchQuery `yaml:"queries" json:"queries,omitempty"`         // Slice of detailed search query configurations.
}

// SearchQuery represents the detailed parameters for a single search query.
type SearchQuery struct {
	K               uint32                              `yaml:"k" json:"k,omitempty"`                               // Number of top results to return.
	Radius          float32                             `yaml:"radius" json:"radius,omitempty"`                     // Radius for search (if applicable).
	Epsilon         float32                             `yaml:"epsilon" json:"epsilon,omitempty"`                   // Epsilon for approximate search algorithms.
	TimeoutString   string                              `yaml:"timeout_string" json:"timeout_string,omitempty"`     // Timeout value as a string; will be parsed to time.Duration.
	AlgorithmString string                              `yaml:"algorithm_string" json:"algorithm_string,omitempty"` // Algorithm identifier as a string; will be normalized and mapped.
	MinNum          uint32                              `yaml:"min_num" json:"min_num,omitempty"`                   // Minimum number of items required for the operation.
	Ratio           float32                             `yaml:"ratio" json:"ratio,omitempty"`                       // Ratio parameter for search (algorithm dependent).
	Nprobe          uint32                              `yaml:"nprobe" json:"nprobe,omitempty"`                     // Number of probes for the search algorithm.
	Timeout         time.Duration                       // Parsed timeout value.
	Algorithm       payload.Search_AggregationAlgorithm // Mapped algorithm constant based on AlgorithmString.
}

// Setting represents basic operation settings used across multiple operations (e.g., insert, update).
// It includes numeric values such as the number of items to process, an offset, and a timestamp.
type Setting struct {
	Num                  uint64 `yaml:"num" json:"num,omitempty"`                                         // Number of items to process.
	Offset               uint64 `yaml:"offset" json:"offset,omitempty"`                                   // Starting offset for the operation.
	BulkSize             uint64 `yaml:"bulk_size" json:"bulk_size,omitempty"`                             // Bulk size for multi-xxx operations.
	Concurrency          uint64 `yaml:"concurrency" json:"concurrency,omitempty"`                         // Concurrency for operations.
	SkipStrictExistCheck bool   `yaml:"skip_strict_exist_check" json:"skip_strict_exist_check,omitempty"` // Flag to indicate if strict existence checks should be skipped.
	Timestamp            int64  `yaml:"timestamp" json:"timestamp,omitempty"`                             // Timestamp value for the operation; used for versioning.
}

// WaitAfterInsert holds configuration regarding the waiting period after an insert operation,
// typically used before starting the indexing process.
type WaitAfterInsert struct {
	String          string        `yaml:"wait_after_insert" json:"string,omitempty"` // Wait duration as a string (e.g., "3m").
	WaitAfterInsert time.Duration // Parsed wait duration value.
}

// Kubernetes holds configuration settings specific to Kubernetes environments.
type Kubernetes struct {
	KubeConfig  string       `yaml:"kubeconfig" json:"kubeconfig,omitempty"`   // File path to the kubeconfig.
	PortForward *PortForward `yaml:"portforward" json:"portforward,omitempty"` // Port forwarding settings.
}

// PortForward holds configuration for port forwarding when running in a Kubernetes environment.
type PortForward struct {
	Enabled    bool   `yaml:"enabled" json:"enabled,omitempty"`         // Flag to enable or disable port forwarding.
	PodName    string `yaml:"pod_name" json:"pod_name,omitempty"`       // The name of the pod to forward from.
	TargetPort uint16 `yaml:"target_port" json:"target_port,omitempty"` // The port forward target port number.
	LocalPort  uint16 `yaml:"local_port" json:"local_port,omitempty"`   // The local port number; if not set, it defaults to TargetPort.
	Namespace  string `yaml:"namespace" json:"namespace,omitempty"`     // The Kubernetes namespace of the pod.
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

	// Process all search-related configurations.
	if d.Search != nil {
		d.Search.Bind()
	}
	if d.SearchByID != nil {
		d.SearchByID.Bind()
	}
	if d.LinearSearch != nil {
		d.LinearSearch.Bind()
	}
	if d.LinearSearchByID != nil {
		d.LinearSearchByID.Bind()
	}

	// Process all operation settings configurations.
	if d.Insert != nil {
		d.Insert.Bind()
	}
	if d.Update != nil {
		d.Update.Bind()
	}
	if d.Upsert != nil {
		d.Upsert.Bind()
	}
	if d.Remove != nil {
		d.Remove.Bind()
	}
	if d.Object != nil {
		d.Object.Bind()
	}

	// Process the wait duration for index creation.
	if d.Index != nil {
		d.Index.Bind()
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

// Bind processes each SearchConfig by iterating through its Queries and binding them.
func (sc *SearchConfig) Bind() *SearchConfig {
	if sc == nil {
		return nil
	}
	if sc.Num == 0 {
		sc.Num = defaultNum
	}

	if sc.Concurrency == 0 {
		sc.Concurrency = defaultConcurrency
	}

	// Iterate over all search queries and bind each one.
	for i, query := range sc.Queries {
		if query != nil {
			sc.Queries[i] = query.Bind()
		}
	}
	return sc
}

// Bind validates and processes the SearchQuery parameters.
// It parses the timeout string into a time.Duration and maps the algorithm string to a constant.
func (sq *SearchQuery) Bind() *SearchQuery {
	if sq == nil {
		return nil
	}
	// Expand and parse the timeout value.
	sq.TimeoutString = config.GetActualValue(sq.TimeoutString)
	// Use a default timeout if parsing fails.
	sq.Timeout = timeutil.ParseWithDefault(sq.TimeoutString, Default.Search.Queries[0].Timeout)

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

// Bind converts the wait duration from string to time.Duration for WaitAfterInsert.
// It uses a default duration (3 minutes) if parsing fails.
func (w *WaitAfterInsert) Bind() *WaitAfterInsert {
	if w == nil {
		return nil
	}
	// Expand the wait duration string.
	w.String = config.GetActualValue(w.String)
	// Parse the duration string, defaulting to 3 minutes on failure.
	w.WaitAfterInsert = timeutil.ParseWithDefault(w.String, Default.Index.WaitAfterInsert)
	return w
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

// Bind validates and processes the Setting configuration for operations such as insert, update, etc.
// It logs warnings if certain numeric values are zero or invalid.
func (s *Setting) Bind() *Setting {
	// Warn if the number of items is zero; this might disable the operation.
	if s.Num == 0 {
		s.Num = defaultNum
	}

	if s.Concurrency == 0 {
		s.Concurrency = defaultConcurrency
	}
	// If the timestamp is negative, reset it to zero and log a warning.
	if s.Timestamp < 0 {
		log.Warn("Setting.Timestamp is negative, resetting to 0")
		s.Timestamp = 0
	}
	return s
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
	pf.PodName = config.GetActualValue(pf.PodName)
	pf.Namespace = config.GetActualValue(pf.Namespace)

	// If TargetPort is not set, default it to the localPort constant.
	if pf.TargetPort == 0 {
		pf.TargetPort = localPort
	}
	// If LocalPort is not set, default it to the same value as TargetPort.
	if pf.LocalPort == 0 {
		pf.LocalPort = localPort
	}
	return pf
}

// Constant definitions for default host and port values.
const (
	localhost        = "localhost"
	localPort uint16 = 8081

	defaultNum                  uint64 = 10000
	defaultOffset               uint64 = 0
	defaultTimestamp            int64  = 0
	defaultSkipStrictExistCheck        = false
	defaultConcurrency          uint64 = 10
)

// Default holds the default configuration values.
// It is used to provide fallback values and defaults for the Bind process.
var Default = &Data{
	Target: &config.GRPCClient{
		Addrs: []string{net.JoinHostPort(localhost, localPort)},
	},
	Search: &SearchConfig{
		Num:         defaultNum,
		Offset:      defaultOffset,
		Concurrency: defaultConcurrency,
		Queries: []*SearchQuery{
			{
				Timeout: time.Second * 3,
			},
		},
	},
	SearchByID: &SearchConfig{
		Num:         defaultNum,
		Offset:      defaultOffset,
		Concurrency: defaultConcurrency,
	},
	LinearSearch: &SearchConfig{
		Num:         defaultNum,
		Offset:      defaultOffset,
		Concurrency: defaultConcurrency,
	},
	LinearSearchByID: &SearchConfig{
		Num:         defaultNum,
		Offset:      defaultOffset,
		Concurrency: defaultConcurrency,
	},
	Insert: &Setting{
		Num:                  defaultNum,
		Offset:               defaultOffset,
		Concurrency:          defaultConcurrency,
		SkipStrictExistCheck: defaultSkipStrictExistCheck,
		Timestamp:            defaultTimestamp,
	},
	Update: &Setting{
		Num:                  defaultNum,
		Offset:               defaultOffset,
		Concurrency:          defaultConcurrency,
		SkipStrictExistCheck: defaultSkipStrictExistCheck,
		Timestamp:            defaultTimestamp,
	},
	Upsert: &Setting{
		Num:                  defaultNum,
		Offset:               defaultOffset,
		Concurrency:          defaultConcurrency,
		SkipStrictExistCheck: defaultSkipStrictExistCheck,
		Timestamp:            defaultTimestamp,
	},
	Remove: &Setting{
		Num:                  defaultNum,
		Offset:               defaultOffset,
		Concurrency:          defaultConcurrency,
		SkipStrictExistCheck: defaultSkipStrictExistCheck,
		Timestamp:            defaultTimestamp,
	},
	Object: &Setting{
		Num:                  defaultNum,
		Offset:               defaultOffset,
		Concurrency:          defaultConcurrency,
		SkipStrictExistCheck: defaultSkipStrictExistCheck,
		Timestamp:            defaultTimestamp,
	},
	Index: &WaitAfterInsert{
		String:          "3m",
		WaitAfterInsert: time.Minute * 3,
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
}

// newData returns initial Data struct.
func newData() *Data {
	return &Data{
		Target:           &config.GRPCClient{},
		Search:           &SearchConfig{},
		SearchByID:       &SearchConfig{},
		LinearSearch:     &SearchConfig{},
		LinearSearchByID: &SearchConfig{},
		Insert:           &Setting{},
		Update:           &Setting{},
		Upsert:           &Setting{},
		Remove:           &Setting{},
		Object:           &Setting{},
		Index:            &WaitAfterInsert{},
		Dataset:          &Dataset{},
		Kubernetes:       &Kubernetes{},
	}
}

// Load reads the configuration from the specified file path.
// If reading fails, it merges the read configuration with the default configuration.
// Finally, it calls Bind to perform all necessary post-processing on the configuration.
func Load(path string) (cfg *Data, err error) {
	log.Debugf("loading test client configuration from %s", path)
	cfg = newData()

	// Attempt to read the configuration from the file.
	err = config.Read(path, cfg)
	if err != nil {
		// If reading fails, merge the configuration with default values.
		cfg, err = config.Merge(cfg, Default)
		if err != nil {
			return nil, err
		}
	}

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
