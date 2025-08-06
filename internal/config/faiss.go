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

package config

// Faiss represent the faiss core configuration for server.
type Faiss struct {
	// VQueue represents the faiss vector queue buffer size
	VQueue                  *VQueue `json:"vqueue,omitempty"                     yaml:"vqueue"`
	// KVSDB represents the faiss bidirectional kv store configuration
	KVSDB                   *KVSDB  `json:"kvsdb,omitempty"                      yaml:"kvsdb"`
	// AutoIndexDurationLimit represents auto indexing duration limit
	AutoIndexDurationLimit  string  `json:"auto_index_duration_limit,omitempty"  yaml:"auto_index_duration_limit"`
	// InitialDelayMaxDuration represents maximum duration for initial delay
	InitialDelayMaxDuration string  `json:"initial_delay_max_duration,omitempty" yaml:"initial_delay_max_duration"`
	// LoadIndexTimeoutFactor represents a factor of load index timeout
	LoadIndexTimeoutFactor  string  `json:"load_index_timeout_factor,omitempty"  yaml:"load_index_timeout_factor"`
	// MethodType represents the method type
	MethodType              string  `json:"method_type,omitempty"                yaml:"method_type"                info:"method_type"`
	// MetricType represents the metric type
	MetricType              string  `json:"metric_type,omitempty"                yaml:"metric_type"                info:"metric_type"`
	// MaxLoadIndexTimeout represents maximum duration of load index timeout
	MaxLoadIndexTimeout     string  `json:"max_load_index_timeout,omitempty"     yaml:"max_load_index_timeout"`
	// AutoIndexCheckDuration represents checking loop duration about auto indexing execution
	AutoIndexCheckDuration  string  `json:"auto_index_check_duration,omitempty"  yaml:"auto_index_check_duration"`
	// AutoSaveIndexDuration represents checking loop duration about auto save index execution
	AutoSaveIndexDuration   string  `json:"auto_save_index_duration,omitempty"   yaml:"auto_save_index_duration"`
	// IndexPath represents the faiss index file path
	IndexPath               string  `json:"index_path,omitempty"                 yaml:"index_path"`
	// MinLoadIndexTimeout represents minimum duration of load index timeout
	MinLoadIndexTimeout     string  `json:"min_load_index_timeout,omitempty"     yaml:"min_load_index_timeout"`
	// AutoIndexLength represents auto index length limit
	AutoIndexLength         int     `json:"auto_index_length,omitempty"          yaml:"auto_index_length"`
	// M represents the number of subquantizers
	// ref: https://github.com/facebookresearch/faiss/wiki/Faiss-indexes-(composite)#cell-probe-method-with-a-pq-index-as-coarse-quantizer
	M                       int     `json:"m,omitempty"                          yaml:"m"                          info:"m"`
	// NbitsPerIdx represents the number of bit per subvector index
	// ref: https://github.com/facebookresearch/faiss/wiki/FAQ#can-i-ignore-warning-clustering-xxx-points-to-yyy-centroids
	NbitsPerIdx             int     `json:"nbits_per_idx,omitempty"              yaml:"nbits_per_idx"              info:"nbits_per_idx"`
	// Nlist represents the number of Voronoi cells
	// ref: https://github.com/facebookresearch/faiss/wiki/Faster-search
	Nlist                   int     `json:"nlist,omitempty"                      yaml:"nlist"                      info:"nlist"`
	// Dimension represents the faiss index dimension
	Dimension               int     `json:"dimension,omitempty"                  yaml:"dimension"                  info:"dimension"`
	// EnableInMemoryMode enables on memory faiss indexing mode
	EnableInMemoryMode      bool    `json:"enable_in_memory_mode,omitempty"      yaml:"enable_in_memory_mode"`
	// EnableProactiveGC enables more proactive GC call for reducing heap memory allocation
	EnableProactiveGC       bool    `json:"enable_proactive_gc,omitempty"        yaml:"enable_proactive_gc"`
	// EnableCopyOnWrite enables copy on write saving
	EnableCopyOnWrite       bool    `json:"enable_copy_on_write,omitempty"       yaml:"enable_copy_on_write"`
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
