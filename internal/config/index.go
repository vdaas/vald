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

// Indexer represents the Indexer configurations.
type Indexer struct {
	// Discoverer represents the discoverer client configuration.
	Discoverer *DiscovererClient `json:"discoverer" yaml:"discoverer"`
	// AutoSaveIndexDurationLimit represents the auto save index duration limit.
	AutoSaveIndexDurationLimit string `json:"auto_save_index_duration_limit" yaml:"auto_save_index_duration_limit"`
	// AgentNamespace represents the agent namespace.
	AgentNamespace string `json:"agent_namespace" yaml:"agent_namespace"`
	// AgentDNS represents the agent DNS.
	AgentDNS string `json:"agent_dns" yaml:"agent_dns"`
	// AutoIndexDurationLimit represents the auto index duration limit.
	AutoIndexDurationLimit string `json:"auto_index_duration_limit" yaml:"auto_index_duration_limit"`
	// AutoSaveIndexWaitDuration represents the auto save index wait duration.
	AutoSaveIndexWaitDuration string `json:"auto_save_index_wait_duration" yaml:"auto_save_index_wait_duration"`
	// AutoIndexCheckDuration represents the auto index check duration.
	AutoIndexCheckDuration string `json:"auto_index_check_duration" yaml:"auto_index_check_duration"`
	// NodeName represents the node name.
	NodeName string `json:"node_name" yaml:"node_name"`
	// AgentName represents the agent name.
	AgentName string `json:"agent_name" yaml:"agent_name"`
	// Concurrency represents the concurrency.
	Concurrency int `json:"concurrency" yaml:"concurrency"`
	// AgentPort represents the agent port.
	AgentPort int `json:"agent_port" yaml:"agent_port"`
	// AutoIndexLength represents the auto index length.
	AutoIndexLength uint32 `json:"auto_index_length" yaml:"auto_index_length"`
	// CreationPoolSize represents the creation pool size.
	CreationPoolSize uint32 `json:"creation_pool_size" yaml:"creation_pool_size"`
}

// Bind binds the actual data from the Indexer receiver field.
func (im *Indexer) Bind() *Indexer {
	im.AgentName = GetActualValue(im.AgentName)
	im.AgentNamespace = GetActualValue(im.AgentNamespace)
	im.AgentDNS = GetActualValue(im.AgentDNS)
	im.AutoIndexDurationLimit = GetActualValue(im.AutoIndexDurationLimit)
	im.AutoIndexCheckDuration = GetActualValue(im.AutoIndexCheckDuration)
	im.AutoSaveIndexDurationLimit = GetActualValue(im.AutoSaveIndexDurationLimit)
	im.AutoSaveIndexWaitDuration = GetActualValue(im.AutoSaveIndexWaitDuration)
	im.NodeName = GetActualValue(im.NodeName)

	if im.Discoverer == nil {
		im.Discoverer = new(DiscovererClient)
	}
	// Assuming DiscovererClient.Bind() is compliant and im.Discoverer is now non-nil
	im.Discoverer.Bind()

	return im
}
