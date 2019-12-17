//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

type Gateway struct {
	// AgentPort represent agent port number
	AgentPort int `json:"agent_port" yaml:"agent_port"`

	// AgentName represent agents meta_name for service discovery
	AgentName string `json:"agent_name" yaml:"agent_name"`

	// AgentDNS represent agents dns A record for service discovery
	AgentDNS string `json:"agent_dns" yaml:"agent_dns"`

	// IndexReplica represents index replication count
	IndexReplica int `json:"index_replica" yaml:"index_replica"`

	// Discoverer represent agent discoverer service configuration
	Discoverer *Discoverer `json:"discoverer" yaml:"discoverer"`

	// Meta represent meta data service configuration
	Meta *Meta `json:"meta" yaml:"meta"`

	// BackupManager represent backup manager configuration
	BackupManager *BackupManager `json:"backup" yaml:"backup"`

	// EgressFilter represents egress filter configuration
	EgressFilter *EgressFilter `json:"egress_filter" yaml:"egress_filter"`
}

func (g *Gateway) Bind() *Gateway {
	g.AgentName = GetActualValue(g.AgentName)

	g.AgentDNS = GetActualValue(g.AgentDNS)

	if g.Discoverer != nil {
		g.Discoverer = g.Discoverer.Bind()
	}
	if g.Meta != nil {
		g.Meta = g.Meta.Bind()
	}
	if g.BackupManager != nil {
		g.BackupManager = g.BackupManager.Bind()
	}

	if g.EgressFilter != nil {
		g.EgressFilter = g.EgressFilter.Bind()
	}
	return g
}
