// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

// IndexCreation represents the configurations for index creation.
type IndexCreation struct {
	// Discoverer represents the discoverer client configuration.
	Discoverer *DiscovererClient `json:"discoverer" yaml:"discoverer"`
	// AgentName represents the agent name.
	AgentName string `json:"agent_name" yaml:"agent_name"`
	// AgentNamespace represents the agent namespace.
	AgentNamespace string `json:"agent_namespace" yaml:"agent_namespace"`
	// AgentDNS represents the agent DNS.
	AgentDNS string `json:"agent_dns" yaml:"agent_dns"`
	// NodeName represents the node name.
	NodeName string `json:"node_name" yaml:"node_name"`
	// TargetAddrs represents the target addresses.
	TargetAddrs []string `json:"target_addrs" yaml:"target_addrs"`
	// AgentPort represents the agent port.
	AgentPort int `json:"agent_port" yaml:"agent_port"`
	// Concurrency represents the concurrency.
	Concurrency int `json:"concurrency" yaml:"concurrency"`
	// CreationPoolSize represents the creation pool size.
	CreationPoolSize uint32 `json:"creation_pool_size" yaml:"creation_pool_size"`
}

func (ic *IndexCreation) Bind() *IndexCreation {
	ic.AgentName = GetActualValue(ic.AgentName)
	ic.AgentNamespace = GetActualValue(ic.AgentNamespace)
	ic.AgentDNS = GetActualValue(ic.AgentDNS)
	ic.NodeName = GetActualValue(ic.NodeName)
	ic.TargetAddrs = GetActualValues(ic.TargetAddrs)

	if ic.Discoverer != nil {
		ic.Discoverer.Bind()
	}
	return ic
}
