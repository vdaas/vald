//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// LB represents the configuration for load balancer.
type LB struct {
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

	// IndexReplica represents index replication count
	IndexReplica int `json:"index_replica" yaml:"index_replica"`

	// ReadReplicaReplicas represents replica count of read replica Deployment
	ReadReplicaReplicas uint64 `json:"read_replica_replicas" yaml:"read_replica_replicas"`

	// ReadReplicaClient represents read replica client configuration
	ReadReplicaClient ReadReplicaClient `json:"read_replica_client" yaml:"read_replica_client"`

	// Discoverer represent agent discoverer service configuration
	Discoverer *DiscovererClient `json:"discoverer" yaml:"discoverer"`

	// MultiOperationConcurrency
	MultiOperationConcurrency int `json:"multi_operation_concurrency" yaml:"multi_operation_concurrency"`
}

// Bind binds the actual data from the LB receiver fields.
func (g *LB) Bind() *LB {
	g.AgentName = GetActualValue(g.AgentName)
	g.AgentNamespace = GetActualValue(g.AgentNamespace)
	g.AgentDNS = GetActualValue(g.AgentDNS)
	g.NodeName = GetActualValue(g.NodeName)

	if g.Discoverer != nil {
		g.Discoverer = g.Discoverer.Bind()
	}
	return g
}

// ReadReplicaClient represents a configuration of grpc client for read replica.
type ReadReplicaClient struct {
	Duration           string      `json:"duration"             yaml:"duration"`
	Client             *GRPCClient `json:"client"               yaml:"client"`
	AgentClientOptions *GRPCClient `json:"agent_client_options" yaml:"agent_client_options"`
}

// Bind binds the actual data from the ReadReplicaClient receiver field.
func (d *ReadReplicaClient) Bind() *ReadReplicaClient {
	d.Duration = GetActualValue(d.Duration)
	if d.Client != nil {
		d.Client.Bind()
	} else {
		d.Client = newGRPCClientConfig()
	}
	if d.AgentClientOptions != nil {
		d.AgentClientOptions.Bind()
	} else {
		d.AgentClientOptions = newGRPCClientConfig()
	}
	return d
}
