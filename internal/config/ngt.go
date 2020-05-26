//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	// IndexPath represent the ngt index file path
	IndexPath string `yaml:"index_path" json:"index_path"`

	// Dimension represent the ngt index dimension
	Dimension int `yaml:"dimension" json:"dimension"`

	// BulkInsertChunkSize represent the bulk insert chunk size
	BulkInsertChunkSize int `yaml:"bulk_insert_chunk_size" json:"bulk_insert_chunk_size"`

	// DistanceType represent the ngt index distance type
	DistanceType string `yaml:"distance_type" json:"distance_type"`

	// ObjectType represent the ngt index object type float or int
	ObjectType string `yaml:"object_type" json:"object_type"`

	// CreationEdgeSize represent the index edge count
	CreationEdgeSize int `yaml:"creation_edge_size" json:"creation_edge_size"`

	// SearchEdgeSize represent the search edge size
	SearchEdgeSize int `yaml:"search_edge_size" json:"search_edge_size"`

	// AutoIndexDurationLimit represents auto indexing duration limit
	AutoIndexDurationLimit string `yaml:"auto_index_duration_limit" json:"auto_index_duration_limit"`

	// AutoIndexCheckDuration represent checking loop duration about auto indexing execution
	AutoIndexCheckDuration string `yaml:"auto_index_check_duration" json:"auto_index_check_duration"`

	// AutoSaveIndexDuration represent checking loop duration about auto save index execution
	AutoSaveIndexDuration string `yaml:"auto_save_index_duration" json:"auto_save_index_duration"`

	// AutoIndexLength represent auto index length limit
	AutoIndexLength int `yaml:"auto_index_length" json:"auto_index_length"`

	// AutoCreateIndexPoolSize represent the pool size for create index operation
	AutoCreateIndexPoolSize uint32 `yaml:"auto_create_index_pool_size" json:"auto_create_index_pool_size"`

	// InitialDelayMaxDuration represent maximum duration for initial delay
	InitialDelayMaxDuration string `yaml:"initial_delay_max_duration" json:"initial_delay_max_duration"`

	// EnableInMemoryMode enables on memory ngt indexing mode
	EnableInMemoryMode bool `yaml:"enable_in_memory_mode" json:"enable_in_memory_mode"`
}

func (n *NGT) Bind() *NGT {
	n.IndexPath = GetActualValue(n.IndexPath)
	n.DistanceType = GetActualValue(n.DistanceType)
	n.ObjectType = GetActualValue(n.ObjectType)
	n.AutoIndexCheckDuration = GetActualValue(n.AutoIndexCheckDuration)
	n.AutoIndexDurationLimit = GetActualValue(n.AutoIndexDurationLimit)
	n.AutoSaveIndexDuration = GetActualValue(n.AutoSaveIndexDuration)
	n.InitialDelayMaxDuration = GetActualValue(n.InitialDelayMaxDuration)
	return n
}
