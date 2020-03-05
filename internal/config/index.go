//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

type Indexer struct {
	// AgentPort represent agent port number
	AgentPort int `json:"agent_port" yaml:"agent_port"`

	// AgentName represent agents meta_name for service discovery
	AgentName string `json:"agent_name" yaml:"agent_name"`

	// AgentNamespace represent agent namespace location
	AgentNamespace string `json:"agent_namespace" yaml:"agent_namespace"`

	// AgentDNS represent agents dns A record for service discovery
	AgentDNS string `json:"agent_dns" yaml:"agent_dns"`

	// Concurrency represents indexing concurrency
	Concurrency int `json:"concurrency" yaml:"concurrency"`

	// AutoIndexLimit auto indexing duration limit
	AutoIndexDurationLimit string `yaml:"auto_index_duration_limit" json:"auto_index_duration_limit"`

	// AutoIndexCheckDuration represent checking loop duration about auto indexing execution
	AutoIndexCheckDuration string `yaml:"auto_index_check_duration" json:"auto_index_check_duration"`

	// AutoIndexLength represent minimum auto index length
	AutoIndexLength uint32 `yaml:"auto_index_length" json:"auto_index_length"`

	// NodeName represents node name
	NodeName string `json:"node_name" yaml:"node_name"`

	// Discoverer represent agent discoverer service configuration
	Discoverer *DiscovererClient `json:"discoverer" yaml:"discoverer"`
}

func (im *Indexer) Bind() *Indexer {
	im.AgentName = GetActualValue(im.AgentName)
	im.AgentNamespace = GetActualValue(im.AgentNamespace)
	im.AgentDNS = GetActualValue(im.AgentDNS)
	im.AutoIndexDurationLimit = GetActualValue(im.AutoIndexDurationLimit)
	im.AutoIndexCheckDuration = GetActualValue(im.AutoIndexCheckDuration)
	im.NodeName = GetActualValue(im.NodeName)

	if im.Discoverer != nil {
		im.Discoverer = im.Discoverer.Bind()
	}
	return im
}
