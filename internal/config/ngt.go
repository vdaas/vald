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

// NGT represent the ngt core configuration for server.
type NGT struct {
	// PodName represent the ngt pod name
	PodName string `yaml:"pod_name" json:"pod_name,omitempty"`

	// PodNamespace represent the ngt pod namespace
	PodNamespace string `yaml:"namespace" json:"namespace,omitempty"`

	// IndexPath represent the ngt index file path
	IndexPath string `json:"index_path,omitempty" yaml:"index_path"`

	// Dimension represent the ngt index dimension
	Dimension int `info:"dimension" json:"dimension,omitempty" yaml:"dimension"`

	// BulkInsertChunkSize represent the bulk insert chunk size
	BulkInsertChunkSize int `json:"bulk_insert_chunk_size,omitempty" yaml:"bulk_insert_chunk_size"`

	// DistanceType represent the ngt index distance type
	// it should be `l1`, `l2`, `angle`, `hamming`, `cosine`,`poincare`, `lorentz`, `jaccard`, `sparsejaccard`, `normalizedangle` or `normalizedcosine`. for further details about NGT libraries supported distance is https://github.com/yahoojapan/NGT/wiki/Command-Quick-Reference and vald agent's supported NGT distance type is https://pkg.go.dev/github.com/vdaas/vald/internal/core/algorithm/ngt#pkg-constants
	DistanceType string `info:"distance_type" json:"distance_type,omitempty" yaml:"distance_type"`

	// ObjectType represent the ngt index object type float or int
	ObjectType string `info:"object_type" json:"object_type,omitempty" yaml:"object_type"`

	// CreationEdgeSize represent the index edge count
	CreationEdgeSize int `json:"creation_edge_size,omitempty" yaml:"creation_edge_size"`

	// SearchEdgeSize represent the search edge size
	SearchEdgeSize int `json:"search_edge_size,omitempty" yaml:"search_edge_size"`

	// AutoIndexDurationLimit represents auto indexing duration limit
	AutoIndexDurationLimit string `json:"auto_index_duration_limit,omitempty" yaml:"auto_index_duration_limit"`

	// AutoIndexCheckDuration represent checking loop duration about auto indexing execution
	AutoIndexCheckDuration string `json:"auto_index_check_duration,omitempty" yaml:"auto_index_check_duration"`

	// AutoSaveIndexDuration represent checking loop duration about auto save index execution
	AutoSaveIndexDuration string `json:"auto_save_index_duration,omitempty" yaml:"auto_save_index_duration"`

	// AutoIndexLength represent auto index length limit
	AutoIndexLength int `json:"auto_index_length,omitempty" yaml:"auto_index_length"`

	// InitialDelayMaxDuration represent maximum duration for initial delay
	InitialDelayMaxDuration string `json:"initial_delay_max_duration,omitempty" yaml:"initial_delay_max_duration"`

	// EnableInMemoryMode enables on memory ngt indexing mode
	EnableInMemoryMode bool `json:"enable_in_memory_mode,omitempty" yaml:"enable_in_memory_mode"`

	// DefaultPoolSize represent default create index batch pool size
	DefaultPoolSize uint32 `json:"default_pool_size,omitempty" yaml:"default_pool_size"`

	// DefaultRadius represent default radius used for search
	DefaultRadius float32 `json:"default_radius,omitempty" yaml:"default_radius"`

	// DefaultEpsilon represent default epsilon used for search
	DefaultEpsilon float32 `json:"default_epsilon,omitempty" yaml:"default_epsilon"`

	// MinLoadIndexTimeout represents minimum duration of load index timeout
	MinLoadIndexTimeout string `json:"min_load_index_timeout,omitempty" yaml:"min_load_index_timeout"`

	// MaxLoadIndexTimeout represents maximum duration of load index timeout
	MaxLoadIndexTimeout string `json:"max_load_index_timeout,omitempty" yaml:"max_load_index_timeout"`

	// LoadIndexTimeoutFactor represents a factor of load index timeout
	LoadIndexTimeoutFactor string `json:"load_index_timeout_factor,omitempty" yaml:"load_index_timeout_factor"`

	// EnableProactiveGC enables more proactive GC call for reducing heap memory allocation
	EnableProactiveGC bool `json:"enable_proactive_gc,omitempty" yaml:"enable_proactive_gc"`

	// EnableCopyOnWrite enables copy on write saving
	EnableCopyOnWrite bool `json:"enable_copy_on_write,omitempty" yaml:"enable_copy_on_write"`

	// VQueue represent the ngt vector queue buffer size
	VQueue *VQueue `json:"vqueue,omitempty" yaml:"vqueue"`

	// KVSDB represent the ngt bidirectional kv store configuration
	KVSDB *KVSDB `json:"kvsdb,omitempty" yaml:"kvsdb"`

	// BrokenIndexHistoryLimit represents the maximum number of broken index generations that will be backed up
	BrokenIndexHistoryLimit int `json:"broken_index_history_limit,omitempty" yaml:"broken_index_history_limit"`

	// ErrorBufferLimit represents the maximum number of core ngt error buffer pool size limit
	ErrorBufferLimit uint64 `json:"error_buffer_limit,omitempty" yaml:"error_buffer_limit"`

	// IsReadReplica represents whether the ngt is read replica or not
	IsReadReplica bool `json:"is_readreplica" yaml:"is_readreplica"`
}

// KVSDB represent the ngt vector bidirectional kv store configuration.
type KVSDB struct {
	// Concurrency represents kvsdb range loop processing concurrency
	Concurrency int `json:"concurrency,omitempty" yaml:"concurrency,omitempty"`
}

// VQueue represent the ngt vector queue buffer size.
type VQueue struct {
	// InsertBufferPoolSize represents insert time ordered slice buffer size
	InsertBufferPoolSize int `json:"insert_buffer_pool_size,omitempty" yaml:"insert_buffer_pool_size"`

	// DeleteBufferPoolSize represents delete time ordered slice buffer size
	DeleteBufferPoolSize int `json:"delete_buffer_pool_size,omitempty" yaml:"delete_buffer_pool_size"`
}

// Bind returns NGT object whose some string value is filed value or environment value.
func (n *NGT) Bind() *NGT {
	n.PodName = GetActualValue(n.PodName)
	n.PodNamespace = GetActualValue(n.PodNamespace)
	n.IndexPath = GetActualValue(n.IndexPath)
	n.DistanceType = GetActualValue(n.DistanceType)
	n.ObjectType = GetActualValue(n.ObjectType)
	n.AutoIndexCheckDuration = GetActualValue(n.AutoIndexCheckDuration)
	n.AutoIndexDurationLimit = GetActualValue(n.AutoIndexDurationLimit)
	n.AutoSaveIndexDuration = GetActualValue(n.AutoSaveIndexDuration)
	n.InitialDelayMaxDuration = GetActualValue(n.InitialDelayMaxDuration)
	n.MinLoadIndexTimeout = GetActualValue(n.MinLoadIndexTimeout)
	n.MaxLoadIndexTimeout = GetActualValue(n.MaxLoadIndexTimeout)
	n.LoadIndexTimeoutFactor = GetActualValue(n.LoadIndexTimeoutFactor)

	if n.VQueue == nil {
		n.VQueue = new(VQueue)
	}
	if n.KVSDB == nil {
		n.KVSDB = new(KVSDB)
	}

	return n
}
