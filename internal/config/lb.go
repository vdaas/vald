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

// LB represents the configuration for load balancer.
type LB struct {
	// ReadReplicaClient represents the read replica client configuration.
	ReadReplicaClient ReadReplicaClient `json:"read_replica_client" yaml:"read_replica_client"`
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
	// AgentPort represents the agent port.
	AgentPort int `json:"agent_port" yaml:"agent_port"`
	// IndexReplica represents the index replica count.
	IndexReplica int `json:"index_replica" yaml:"index_replica"`
	// ReadReplicaReplicas represents the read replica replicas count.
	ReadReplicaReplicas uint64 `json:"read_replica_replicas" yaml:"read_replica_replicas"`
	// MultiOperationConcurrency represents the multi operation concurrency.
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
	// Client represents the gRPC client configuration.
	Client *GRPCClient `json:"client" yaml:"client"`
	// AgentClientOptions represents the agent client options.
	AgentClientOptions *GRPCClient `json:"agent_client_options" yaml:"agent_client_options"`
	// Duration represents the duration.
	Duration string `json:"duration" yaml:"duration"`
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
