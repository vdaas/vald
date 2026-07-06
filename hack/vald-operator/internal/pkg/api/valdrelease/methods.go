package valdrelease

import (
	v1 "k8s.io/api/core/v1"
)

// ResourceParams carries the configuration-derived knobs needed to populate
// resource-related fields on a ValdRelease. The caller (typically the builder
// driven by *config.Config) is responsible for filling this in; keeping the
// values as a plain struct lets the api package stay independent of the
// configuration source.
type ResourceParams struct {
	AgentPodsPerNode           int
	DiscovererDSMaxSurge       string
	DiscovererDSMaxUnavailable string
}

func (v *ValdRelease) SetRelationalResources(ar int, am v1.ResourceList, p ResourceParams) {
	v.Spec.Agent.SetReplica(ar, p.AgentPodsPerNode)
	v.Spec.Agent.SetResources(am, p.AgentPodsPerNode)
	v.Spec.Agent.SetTopologySpreadConstraints()
	v.Spec.Gateway.Lb.SetReplica(v.Spec.Agent)
	v.Spec.Gateway.Lb.SetResources()
	v.Spec.Gateway.Lb.SetTopologySpreadConstraints()
	v.Spec.Discoverer.ApplyDefaultsByKind(p.DiscovererDSMaxSurge, p.DiscovererDSMaxUnavailable)
	v.Spec.Discoverer.SetResources()
	v.Spec.Discoverer.SetTopologySpreadConstraints()
	if v.Spec.Manager.Index.Enabled {
		v.Spec.Manager.Index.SetResources()
		v.Spec.Manager.Index.SetTopologySpreadConstraints()
	}
}
