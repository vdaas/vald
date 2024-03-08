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

// IndexOperator represents the configurations for index k8s operator.
type IndexOperator struct {
	// AgentName represent agents meta_name for service discovery
	AgentName string `json:"agent_name" yaml:"agent_name"`

	// AgentNamespace represent agent namespace location
	AgentNamespace string `json:"agent_namespace" yaml:"agent_namespace"`

	// Concurrency represents indexing concurrency.
	Concurrency int `json:"concurrency" yaml:"concurrency"`

	// ReadReplicaEnabled represents whether read replica is enabled or not.
	ReadReplicaEnabled bool `json:"read_replica_enabled" yaml:"read_replica_enabled"`

	// ReadReplicaLabelKey represents the label key for read replica.
	ReadReplicaLabelKey string `json:"read_replica_label_key" yaml:"read_replica_label_key"`
}

func (ic *IndexOperator) Bind() *IndexOperator {
	ic.AgentName = GetActualValue(ic.AgentName)
	ic.AgentNamespace = GetActualValue(ic.AgentNamespace)
	return ic
}
