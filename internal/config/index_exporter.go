// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package config

// IndexExporter represents the configurations for index exportation.
type IndexExporter struct {
	// Concurrency represents indexing concurrency.
	Concurrency int `json:"concurrency" yaml:"concurrency"`

	// KVSBackgroundSyncInterval represents interval for checked id list kvs sync duration
	KVSBackgroundSyncInterval string `json:"kvs_background_sync_interval" yaml:"kvs_background_sync_interval"`

	// KVSBackgroundCompactionInterval represents interval for checked id list kvs compaction duration
	KVSBackgroundCompactionInterval string `json:"kvs_background_compaction_interval" yaml:"kvs_background_compaction_interval"`

	// IndexPath represents the export index file path
	IndexPath string `json:"index_path,omitempty" yaml:"index_path"`

	// Gateway represent gateway service configuration
	Gateway *GRPCClient `json:"gateway" yaml:"gateway"`
}

func (e *IndexExporter) Bind() *IndexExporter {
	e.KVSBackgroundCompactionInterval = GetActualValue(e.KVSBackgroundCompactionInterval)
	e.KVSBackgroundSyncInterval = GetActualValue(e.KVSBackgroundSyncInterval)
	e.IndexPath = GetActualValue(e.IndexPath)

	if e.Gateway != nil {
		e.Gateway = e.Gateway.Bind()
	}
	return e
}
