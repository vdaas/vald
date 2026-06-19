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

// Faiss represent the faiss core configuration for server.
type Faiss struct {
	// VQueue represents the vector queue configuration.
	VQueue *VQueue `json:"vqueue,omitempty" yaml:"vqueue"`
	// KVSDB represents the KVS DB configuration.
	KVSDB *KVSDB `json:"kvsdb,omitempty" yaml:"kvsdb"`
	// AutoIndexDurationLimit represents the auto index duration limit.
	AutoIndexDurationLimit string `json:"auto_index_duration_limit,omitempty" yaml:"auto_index_duration_limit"`
	// InitialDelayMaxDuration represents the initial delay maximum duration.
	InitialDelayMaxDuration string `json:"initial_delay_max_duration,omitempty" yaml:"initial_delay_max_duration"`
	// LoadIndexTimeoutFactor represents the load index timeout factor.
	LoadIndexTimeoutFactor string `json:"load_index_timeout_factor,omitempty" yaml:"load_index_timeout_factor"`
	// MethodType represents the method type.
	MethodType string `json:"method_type,omitempty" yaml:"method_type" info:"method_type"`
	// MetricType represents the metric type.
	MetricType string `json:"metric_type,omitempty" yaml:"metric_type" info:"metric_type"`
	// MaxLoadIndexTimeout represents the maximum load index timeout.
	MaxLoadIndexTimeout string `json:"max_load_index_timeout,omitempty" yaml:"max_load_index_timeout"`
	// AutoIndexCheckDuration represents the auto index check duration.
	AutoIndexCheckDuration string `json:"auto_index_check_duration,omitempty" yaml:"auto_index_check_duration"`
	// AutoSaveIndexDuration represents the auto save index duration.
	AutoSaveIndexDuration string `json:"auto_save_index_duration,omitempty" yaml:"auto_save_index_duration"`
	// IndexPath represents the index path.
	IndexPath string `json:"index_path,omitempty" yaml:"index_path"`
	// MinLoadIndexTimeout represents the minimum load index timeout.
	MinLoadIndexTimeout string `json:"min_load_index_timeout,omitempty" yaml:"min_load_index_timeout"`
	// AutoIndexLength represents the auto index length.
	AutoIndexLength int `json:"auto_index_length,omitempty" yaml:"auto_index_length"`
	// M represents the number of neighbors for graph construction.
	M int `json:"m,omitempty" yaml:"m" info:"m"`
	// NbitsPerIdx represents the number of bits per index.
	NbitsPerIdx int `json:"nbits_per_idx,omitempty" yaml:"nbits_per_idx" info:"nbits_per_idx"`
	// Nlist represents the number of cells.
	Nlist int `json:"nlist,omitempty" yaml:"nlist" info:"nlist"`
	// Dimension represents the dimension of the vector.
	Dimension int `json:"dimension,omitempty" yaml:"dimension" info:"dimension"`
	// EnableInMemoryMode enables in-memory mode.
	EnableInMemoryMode bool `json:"enable_in_memory_mode,omitempty" yaml:"enable_in_memory_mode"`
	// EnableProactiveGC enables proactive GC.
	EnableProactiveGC bool `json:"enable_proactive_gc,omitempty" yaml:"enable_proactive_gc"`
	// EnableCopyOnWrite enables copy on write.
	EnableCopyOnWrite bool `json:"enable_copy_on_write,omitempty" yaml:"enable_copy_on_write"`
}

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
	f.VQueue.Bind() // Call Bind on VQueue

	if f.KVSDB == nil {
		f.KVSDB = new(KVSDB)
	}
	f.KVSDB.Bind() // Call Bind on KVSDB

	return f
}
