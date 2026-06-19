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

// NGT represent the ngt core configuration for server.
type NGT struct {
	// KVSDB represents the KVS DB configuration.
	KVSDB *KVSDB `json:"kvsdb,omitempty" yaml:"kvsdb"`
	// VQueue represents the vector queue configuration.
	VQueue *VQueue `json:"vqueue,omitempty" yaml:"vqueue"`
	// ObjectType represents the object type.
	ObjectType string `json:"object_type,omitempty" yaml:"object_type" info:"object_type"`
	// ExportIndexInfoDuration represents the export index info duration.
	ExportIndexInfoDuration string `json:"export_index_info_duration,omitempty" yaml:"export_index_info_duration"`
	// IndexPath represents the index path.
	IndexPath string `json:"index_path,omitempty" yaml:"index_path"`
	// DistanceType represents the distance type.
	DistanceType string `json:"distance_type,omitempty" yaml:"distance_type" info:"distance_type"`
	// MinLoadIndexTimeout represents the minimum load index timeout.
	MinLoadIndexTimeout string `json:"min_load_index_timeout,omitempty" yaml:"min_load_index_timeout"`
	// PodNamespace represents the pod namespace.
	PodNamespace string `json:"namespace,omitempty" yaml:"namespace"`
	// PodName represents the pod name.
	PodName string `json:"pod_name,omitempty" yaml:"pod_name"`
	// AutoIndexDurationLimit represents the auto index duration limit.
	AutoIndexDurationLimit string `json:"auto_index_duration_limit,omitempty" yaml:"auto_index_duration_limit"`
	// AutoIndexCheckDuration represents the auto index check duration.
	AutoIndexCheckDuration string `json:"auto_index_check_duration,omitempty" yaml:"auto_index_check_duration"`
	// AutoSaveIndexDuration represents the auto save index duration.
	AutoSaveIndexDuration string `json:"auto_save_index_duration,omitempty" yaml:"auto_save_index_duration"`
	// LoadIndexTimeoutFactor represents the load index timeout factor.
	LoadIndexTimeoutFactor string `json:"load_index_timeout_factor,omitempty" yaml:"load_index_timeout_factor"`
	// InitialDelayMaxDuration represents the initial delay maximum duration.
	InitialDelayMaxDuration string `json:"initial_delay_max_duration,omitempty" yaml:"initial_delay_max_duration"`
	// MaxLoadIndexTimeout represents the maximum load index timeout.
	MaxLoadIndexTimeout string `json:"max_load_index_timeout,omitempty" yaml:"max_load_index_timeout"`
	// SearchEdgeSize represents the search edge size.
	SearchEdgeSize int `json:"search_edge_size,omitempty" yaml:"search_edge_size"`
	// BrokenIndexHistoryLimit represents the broken index history limit.
	BrokenIndexHistoryLimit int `json:"broken_index_history_limit,omitempty" yaml:"broken_index_history_limit"`
	// AutoIndexLength represents the auto index length.
	AutoIndexLength int `json:"auto_index_length,omitempty" yaml:"auto_index_length"`
	// Dimension represents the dimension of the vector.
	Dimension int `json:"dimension,omitempty" yaml:"dimension" info:"dimension"`
	// ErrorBufferLimit represents the error buffer limit.
	ErrorBufferLimit uint64 `json:"error_buffer_limit,omitempty" yaml:"error_buffer_limit"`
	// CreationEdgeSize represents the creation edge size.
	CreationEdgeSize int `json:"creation_edge_size,omitempty" yaml:"creation_edge_size"`
	// BulkInsertChunkSize represents the bulk insert chunk size.
	BulkInsertChunkSize int `json:"bulk_insert_chunk_size,omitempty" yaml:"bulk_insert_chunk_size"`
	// DefaultEpsilon represents the default epsilon.
	DefaultEpsilon float32 `json:"default_epsilon,omitempty" yaml:"default_epsilon"`
	// EpsilonForCreation represents the epsilon for creation.
	EpsilonForCreation float32 `json:"epsilon_for_creation,omitempty" yaml:"epsilon_for_creation"`
	// DefaultPoolSize represents the default pool size.
	DefaultPoolSize uint32 `json:"default_pool_size,omitempty" yaml:"default_pool_size"`
	// DefaultRadius represents the default radius.
	DefaultRadius float32 `json:"default_radius,omitempty" yaml:"default_radius"`
	// EnableInMemoryMode enables in-memory mode.
	EnableInMemoryMode bool `json:"enable_in_memory_mode,omitempty" yaml:"enable_in_memory_mode"`
	// EnableCopyOnWrite enables copy on write.
	EnableCopyOnWrite bool `json:"enable_copy_on_write,omitempty" yaml:"enable_copy_on_write"`
	// IsReadReplica enables read replica.
	IsReadReplica bool `json:"is_readreplica" yaml:"is_readreplica"`
	// EnableExportIndexInfoToK8s enables exporting index info to k8s.
	EnableExportIndexInfoToK8s bool `json:"enable_export_index_info_to_k8s" yaml:"enable_export_index_info_to_k8s"`
	// EnableProactiveGC enables proactive GC.
	EnableProactiveGC bool `json:"enable_proactive_gc,omitempty" yaml:"enable_proactive_gc"`
	// EnableStatistics enables statistics.
	EnableStatistics bool `json:"enable_statistics" yaml:"enable_statistics"`
}

// KVSDB represent the ngt vector bidirectional kv store configuration.
type KVSDB struct {
	// Concurrency represents kvsdb range loop processing concurrency
	Concurrency int `json:"concurrency,omitempty" yaml:"concurrency,omitempty"`
}

// Bind binds the actual data from the KVSDB receiver fields.
func (k *KVSDB) Bind() *KVSDB {
	// No string fields or nested structs to bind in KVSDB itself
	return k
}

// VQueue represent the ngt vector queue buffer size.
type VQueue struct {
	// InsertBufferPoolSize represents insert time ordered slice buffer size
	InsertBufferPoolSize int `json:"insert_buffer_pool_size,omitempty" yaml:"insert_buffer_pool_size"`

	// DeleteBufferPoolSize represents delete time ordered slice buffer size
	DeleteBufferPoolSize int `json:"delete_buffer_pool_size,omitempty" yaml:"delete_buffer_pool_size"`
}

// Bind binds the actual data from the VQueue receiver fields.
func (vq *VQueue) Bind() *VQueue {
	// No string fields or nested structs to bind in VQueue itself
	return vq
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
	n.ExportIndexInfoDuration = GetActualValue(n.ExportIndexInfoDuration)

	if n.VQueue == nil {
		n.VQueue = new(VQueue)
	}
	n.VQueue.Bind()

	if n.KVSDB == nil {
		n.KVSDB = new(KVSDB)
	}
	n.KVSDB.Bind()

	return n
}
