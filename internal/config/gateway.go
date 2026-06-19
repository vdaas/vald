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

// Gateway represents the list of configurations for gateway.
type Gateway struct {
	// Discoverer represents the discoverer client configuration.
	Discoverer *DiscovererClient `json:"discoverer" yaml:"discoverer"`
	// Meta represents the metadata configuration.
	Meta *Meta `json:"meta" yaml:"meta"`
	// BackupManager represents the backup manager configuration.
	BackupManager *BackupManager `json:"backup" yaml:"backup"`
	// EgressFilter represents the egress filter configuration.
	EgressFilter *EgressFilter `json:"egress_filter" yaml:"egress_filter"`
	// AgentName represents the agent name.
	AgentName string `json:"agent_name" yaml:"agent_name"`
	// AgentNamespace represents the agent namespace.
	AgentNamespace string `json:"agent_namespace" yaml:"agent_namespace"`
	// AgentDNS represents the agent DNS.
	AgentDNS string `json:"agent_dns" yaml:"agent_dns"`
	// NodeName represents the node name.
	NodeName string `json:"node_name" yaml:"node_name"`
	// AgentPort represents the agent port.
	AgentPort int `json:"agent_port" yaml:"agent_port"`
	// IndexReplica represents the index replica count.
	IndexReplica int `json:"index_replica" yaml:"index_replica"`
}

// Bind binds the actual data from the Gateway receiver field on the Gateway.
func (g *Gateway) Bind() *Gateway {
	g.AgentName = GetActualValue(g.AgentName)
	g.AgentNamespace = GetActualValue(g.AgentNamespace)

	g.AgentDNS = GetActualValue(g.AgentDNS)

	g.NodeName = GetActualValue(g.NodeName)

	if g.Discoverer != nil {
		g.Discoverer.Bind()
	}

	if g.Meta != nil {
		g.Meta.Bind()
	}

	if g.BackupManager != nil {
		g.BackupManager.Bind()
	}

	if g.EgressFilter != nil {
		g.EgressFilter.Bind()
	}

	return g
}
