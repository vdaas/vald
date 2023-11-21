//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Corrector represents the index correction configurations.
type Corrector struct {
	// AgentPort represent agent port number
	AgentPort int `json:"agent_port" yaml:"agent_port"`

	// AgentName represent agents meta_name for service discovery
	AgentName string `json:"agent_name" yaml:"agent_name"`

	// AgentNamespace represent agent namespace location
	AgentNamespace string `json:"agent_namespace" yaml:"agent_namespace"`

	// AgentDNS represent agents dns A record for service discovery
	AgentDNS string `json:"agent_dns" yaml:"agent_dns"`

	// NodeName represents node name
	NodeName string `json:"node_name" yaml:"node_name"`

	// StreamConcurrency represent stream concurrency for StreamListObject rpc client
	// this directly affects the memory usage of this job
	StreamListConcurrency int `json:"stream_list_concurrency" yaml:"stream_list_concurrency"`

	// KvsAsyncWriteConcurrency represent concurrency for kvs async write
	KvsAsyncWriteConcurrency int `json:"kvs_async_write_concurrency" yaml:"kvs_async_write_concurrency"`

	// IndexReplica represent index replica count. This should be equal to the lb setting
	IndexReplica int `json:"index_replica" yaml:"index_replica"`

	// Discoverer represent agent discoverer service configuration
	Discoverer *DiscovererClient `json:"discoverer" yaml:"discoverer"`
}

// Bind binds the actual data from the Indexer receiver field.
func (c *Corrector) Bind() *Corrector {
	c.AgentName = GetActualValue(c.AgentName)
	c.AgentNamespace = GetActualValue(c.AgentNamespace)
	c.AgentDNS = GetActualValue(c.AgentDNS)
	c.NodeName = GetActualValue(c.NodeName)

	if c.Discoverer != nil {
		c.Discoverer = c.Discoverer.Bind()
	}
	return c
}
