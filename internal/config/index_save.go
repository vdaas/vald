// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// IndexSave represents the configurations for index save.
type IndexSave struct {
	// AgentPort represent agent port number.
	AgentPort int `json:"agent_port" yaml:"agent_port"`

	// AgentName represent agents meta_name for service discovery.
	AgentName string `json:"agent_name" yaml:"agent_name"`

	// AgentNamespace represent agent namespace location.
	AgentNamespace string `json:"agent_namespace" yaml:"agent_namespace"`

	// AgentDNS represent agents dns A record for service discovery.
	AgentDNS string `json:"agent_dns" yaml:"agent_dns"`

	// NodeName represents node name.
	NodeName string `json:"node_name" yaml:"node_name"`

	// Concurrency represents indexing concurrency.
	Concurrency int `json:"concurrency" yaml:"concurrency"`

	// TargetAddrs represents indexing target addresses.
	TargetAddrs []string `json:"target_addrs" yaml:"target_addrs"`

	// Discoverer represents agent discoverer service configuration.
	Discoverer *DiscovererClient `json:"discoverer" yaml:"discoverer"`
}

func (is *IndexSave) Bind() *IndexSave {
	is.AgentName = GetActualValue(is.AgentName)
	is.AgentNamespace = GetActualValue(is.AgentNamespace)
	is.AgentDNS = GetActualValue(is.AgentDNS)
	is.NodeName = GetActualValue(is.NodeName)
	is.TargetAddrs = GetActualValues(is.TargetAddrs)

	if is.Discoverer != nil {
		is.Discoverer.Bind()
	}
	return is
}
