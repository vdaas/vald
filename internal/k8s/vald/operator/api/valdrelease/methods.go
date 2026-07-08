// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
package valdrelease

import v1 "k8s.io/api/core/v1"

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
