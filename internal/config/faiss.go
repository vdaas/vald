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

// Faiss represent the faiss core configuration for server.
type Faiss struct {
	// IndexPath represents the faiss index file path
	IndexPath string `yaml:"index_path" json:"index_path,omitempty"`

	// Dimension represents the faiss index dimension
	Dimension int `yaml:"dimension" json:"dimension,omitempty" info:"dimension"`

	// Nlist represents the number of Voronoi cells
	// ref: https://github.com/facebookresearch/faiss/wiki/Faster-search
	Nlist int `yaml:"nlist" json:"nlist,omitempty" info:"nlist"`

	// M represents the number of subquantizers
	// ref: https://github.com/facebookresearch/faiss/wiki/Faiss-indexes-(composite)#cell-probe-method-with-a-pq-index-as-coarse-quantizer
	M int `yaml:"m" json:"m,omitempty" info:"m"`

	// NbitsPerIdx represents the number of bit per subvector index
	// ref: https://github.com/facebookresearch/faiss/wiki/FAQ#can-i-ignore-warning-clustering-xxx-points-to-yyy-centroids
	NbitsPerIdx int `yaml:"nbits_per_idx" json:"nbits_per_idx,omitempty" info:"nbits_per_idx"`

	// MetricType represents the metric type
	MetricType string `yaml:"metric_type" json:"metric_type,omitempty" info:"metric_type"`

	// EnableInMemoryMode enables on memory faiss indexing mode
	EnableInMemoryMode bool `yaml:"enable_in_memory_mode" json:"enable_in_memory_mode,omitempty"`

	// AutoIndexCheckDuration represents checking loop duration about auto indexing execution
	AutoIndexCheckDuration string `yaml:"auto_index_check_duration" json:"auto_index_check_duration,omitempty"`

	// AutoSaveIndexDuration represents checking loop duration about auto save index execution
	AutoSaveIndexDuration string `yaml:"auto_save_index_duration" json:"auto_save_index_duration,omitempty"`

	// AutoIndexDurationLimit represents auto indexing duration limit
	AutoIndexDurationLimit string `yaml:"auto_index_duration_limit" json:"auto_index_duration_limit,omitempty"`

	// AutoIndexLength represents auto index length limit
	AutoIndexLength int `yaml:"auto_index_length" json:"auto_index_length,omitempty"`

	// InitialDelayMaxDuration represents maximum duration for initial delay
	InitialDelayMaxDuration string `yaml:"initial_delay_max_duration" json:"initial_delay_max_duration,omitempty"`

	// MinLoadIndexTimeout represents minimum duration of load index timeout
	MinLoadIndexTimeout string `yaml:"min_load_index_timeout" json:"min_load_index_timeout,omitempty"`

	// MaxLoadIndexTimeout represents maximum duration of load index timeout
	MaxLoadIndexTimeout string `yaml:"max_load_index_timeout" json:"max_load_index_timeout,omitempty"`

	// LoadIndexTimeoutFactor represents a factor of load index timeout
	LoadIndexTimeoutFactor string `yaml:"load_index_timeout_factor" json:"load_index_timeout_factor,omitempty"`

	// EnableProactiveGC enables more proactive GC call for reducing heap memory allocation
	EnableProactiveGC bool `yaml:"enable_proactive_gc" json:"enable_proactive_gc,omitempty"`

	// EnableCopyOnWrite enables copy on write saving
	EnableCopyOnWrite bool `yaml:"enable_copy_on_write" json:"enable_copy_on_write,omitempty"`

	// VQueue represents the faiss vector queue buffer size
	VQueue *VQueue `json:"vqueue,omitempty" yaml:"vqueue"`

	// KVSDB represents the faiss bidirectional kv store configuration
	KVSDB *KVSDB `json:"kvsdb,omitempty" yaml:"kvsdb"`
}

//// KVSDB represent the faiss vector bidirectional kv store configuration
// type KVSDB struct {
// 	// Concurrency represents kvsdb range loop processing concurrency
// 	Concurrency int `json:"concurrency,omitempty" yaml:"concurrency,omitempty"`
// }
//// VQueue represent the faiss vector queue buffer size
// type VQueue struct {
// 	// InsertBufferPoolSize represents insert time ordered slice buffer size
// 	InsertBufferPoolSize int `json:"insert_buffer_pool_size,omitempty" yaml:"insert_buffer_pool_size"`
//
// 	// DeleteBufferPoolSize represents delete time ordered slice buffer size
// 	DeleteBufferPoolSize int `json:"delete_buffer_pool_size,omitempty" yaml:"delete_buffer_pool_size"`
// }

// Bind returns Faiss object whose some string value is filed value or environment value.
func (f *Faiss) Bind() *Faiss {
	f.IndexPath = GetActualValue(f.IndexPath)
	f.MetricType = GetActualValue(f.MetricType)
	f.AutoIndexCheckDuration = GetActualValue(f.AutoIndexCheckDuration)
	f.AutoIndexDurationLimit = GetActualValue(f.AutoIndexDurationLimit)
	f.AutoSaveIndexDuration = GetActualValue(f.AutoSaveIndexDuration)
	f.InitialDelayMaxDuration = GetActualValue(f.InitialDelayMaxDuration)
	f.MinLoadIndexTimeout = GetActualValue(f.MinLoadIndexTimeout)
	f.MaxLoadIndexTimeout = GetActualValue(f.MaxLoadIndexTimeout)
	f.LoadIndexTimeoutFactor = GetActualValue(f.LoadIndexTimeoutFactor)

	if f.VQueue == nil {
		f.VQueue = new(VQueue)
	}
	if f.KVSDB == nil {
		f.KVSDB = new(KVSDB)
	}

	return f
}
