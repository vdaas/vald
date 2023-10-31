// Package config providers configuration type and load configuration logic
package config

// IndexCreator represents the index creator configurations.
type IndexCreator struct {
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

	// Concurrency represents indexing concurrency.
	Concurrency int `json:"concurrency" yaml:"concurrency"`

	// CreationPoolSize represents batch pool size for indexing.
	CreationPoolSize uint32 `yaml:"creation_pool_size" json:"creation_pool_size"`

	// TargetAddrs represents indexing target addresses.
	TargetAddrs []string `yaml:"target_addrs" json:"target_addrs"`

	// Discoverer represents agent discoverer service configuration.
	Discoverer *DiscovererClient `json:"discoverer" yaml:"discoverer"`
}

func (ic *IndexCreator) Bind() *IndexCreator {
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
