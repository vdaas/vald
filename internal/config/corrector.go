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

// Corrector represents the index correction configurations.
type Corrector struct {
	// Discoverer represents the discoverer client configuration.
	Discoverer *DiscovererClient `json:"discoverer" yaml:"discoverer"`
	// Gateway represents the gateway client configuration.
	Gateway *GRPCClient `json:"gateway" yaml:"gateway"`
	// AgentName represents the agent name.
	AgentName string `json:"agent_name" yaml:"agent_name"`
	// AgentNamespace represents the agent namespace.
	AgentNamespace string `json:"agent_namespace" yaml:"agent_namespace"`
	// AgentDNS represents the agent DNS.
	AgentDNS string `json:"agent_dns" yaml:"agent_dns"`
	// NodeName represents the node name.
	NodeName string `json:"node_name" yaml:"node_name"`
	// KVSBackgroundSyncInterval represents the interval for KVS background synchronization.
	KVSBackgroundSyncInterval string `json:"kvs_background_sync_interval" yaml:"kvs_background_sync_interval"`
	// KVSBackgroundCompactionInterval represents the interval for KVS background compaction.
	KVSBackgroundCompactionInterval string `json:"kvs_background_compaction_interval" yaml:"kvs_background_compaction_interval"`
	// AgentPort represents the agent port.
	AgentPort int `json:"agent_port" yaml:"agent_port"`
	// StreamListConcurrency represents the stream list concurrency.
	StreamListConcurrency int `json:"stream_list_concurrency" yaml:"stream_list_concurrency"`
	// IndexReplica represents the index replica count.
	IndexReplica int `json:"index_replica" yaml:"index_replica"`
}

// Bind binds the actual data from the Indexer receiver field.
func (c *Corrector) Bind() *Corrector {
	c.AgentName = GetActualValue(c.AgentName)
	c.AgentNamespace = GetActualValue(c.AgentNamespace)
	c.AgentDNS = GetActualValue(c.AgentDNS)
	c.NodeName = GetActualValue(c.NodeName)
	c.KVSBackgroundCompactionInterval = GetActualValue(c.KVSBackgroundCompactionInterval)
	c.KVSBackgroundSyncInterval = GetActualValue(c.KVSBackgroundSyncInterval)

	if c.Discoverer != nil {
		c.Discoverer = c.Discoverer.Bind()
	}
	if c.Gateway != nil {
		c.Gateway = c.Gateway.Bind()
	}
	return c
}
