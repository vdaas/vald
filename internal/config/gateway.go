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

type Gateway struct {
	// AgentPort represent agent port number
	AgentPort int `json:"agent_port" yaml:"agent_port"`

	// AgentName represent agents meta_name for service discovery
	AgentName string `json:"agent_name" yaml:"agent_name"`

	// BackoffEnabled enables backoff algorithms for each external request
	BackoffEnabled bool `json:"backoff_enabled" yaml:"backoff_enabled"`

	// Discoverer represent agent discoverer service configuration
	Discoverer *Discoverer `json:"discoverer" yaml:"discoverer"`

	// MetaProxy represent metadata gateway configuration
	Meta *MetaProxy `json:"meta" yaml:"meta"`

	// IndexManager represent index manager configuration
	IndexManager *IndexManager `json:"meta" yaml:"meta"`
}

func (g *Gateway) Bind() *Gateway {
	g.AgentName = GetActualValue(g.AgentName)

	if g.Discoverer != nil {
		g.Discoverer = g.Discoverer.Bind()
	}
	if g.Meta != nil {
		g.Meta = g.Meta.Bind()
	}
	if g.IndexManager.Bind() != nil {
		g.IndexManager = g.IndexManager.Bind()
	}
	return g
}
