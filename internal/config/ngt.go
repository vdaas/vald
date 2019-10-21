//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	IndexPath string `yaml:"index_path"`

	// Dimension represent the ngt index dimension
	Dimension int `yaml:"dimension"`

	// BulkInsertChunkSize represent the bulk insert chunk size
	BulkInsertChunkSize int `yaml:"bulk_insert_chunk_size"`

	// DistanceType represent the ngt index distance type
	DistanceType string `yaml:"distance_type"`

	// ObjectType represent the ngt index object type float or int
	ObjectType string `yaml:"object_type"`

	// CreationEdgeSize represent the index edge count
	CreationEdgeSize int `yaml:"creation_edge_size"`

	// SearchEdgeSize represent the search edge size
	SearchEdgeSize int `yaml:"search_edge_size"`
}

func (n *NGT) Bind() *NGT {
	n.IndexPath = GetActualValue(n.IndexPath)
	n.DistanceType = GetActualValue(n.DistanceType)
	n.ObjectType = GetActualValue(n.ObjectType)
	return n
}
