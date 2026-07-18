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

package valdrelease

import (
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	// AgentResourceRatio is the fraction of a node's resources allocated to agent
	// pods scheduled on it.
	AgentResourceRatio = 0.6
	// DefaultIndexPath is the on-disk index path used for PV-backed agents.
	DefaultIndexPath = "/var/ngt/index"
	// DefaultAgentMaxSurge / DefaultAgentMaxUnavailable are the rolling-update
	// defaults applied to the agent StatefulSet when the CR omits them.
	DefaultAgentMaxSurge       = "1"
	DefaultAgentMaxUnavailable = "1"

	componentLabelAgent        = "agent"
	componentLabelGatewayLb    = "gateway-lb"
	componentLabelDiscoverer   = "discoverer"
	componentLabelManagerIndex = "manager-index"
)

// ResourceParams carries the configuration-derived knobs needed to populate
// resource-related fields. The builder fills it from *config.Config and passes
// it to SetScaledResources, keeping this package free of config concerns.
type ResourceParams struct {
	DiscovererDSMaxSurge       string
	DiscovererDSMaxUnavailable string
	AgentPodsPerNode           int
}

// setTopologySpreadConstraints returns a single-element TSC spreading pods
// across nodes by hostname, scoped to the given app.kubernetes.io/component
// label value.
func setTopologySpreadConstraints(componentLabel string) *TopologySpreadConstraints {
	tsc := TopologySpreadConstraints{common.BuildTopologySpreadConstraint(componentLabel)}
	return &tsc
}

// fixedResources builds a Resources value from literal request/limit
// quantities. It is used by the components whose compute resources are
// fixed constants rather than derived from node/agent sizing (Agent's
// SetResources is computed dynamically and does not use this helper).
func fixedResources(reqCPU, reqMemory, limCPU, limMemory string) *Resources {
	return &Resources{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(reqCPU),
			corev1.ResourceMemory: resource.MustParse(reqMemory),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(limCPU),
			corev1.ResourceMemory: resource.MustParse(limMemory),
		},
	}
}

// --- Agent -----------------------------------------------------------------

// SetReplica sets min/max replicas from node count * pods-per-node.
func (a *Agent) SetReplica(nr, podsPerNode int) {
	rep := nr * podsPerNode
	a.MinReplicas = new(rep)
	a.MaxReplicas = new(rep)
}

// SetResources derives per-pod CPU/memory from the node machine resources.
// Memory limit is intentionally omitted (the NGT index grows after startup).
func (a *Agent) SetResources(mc corev1.ResourceList, podsPerNode int) {
	div := float64(podsPerNode)
	cpu := mc.Cpu().Value()
	memory := mc.Memory().Value()
	a.Resources = &Resources{
		Requests: common.NormalizeResourceList(corev1.ResourceList{
			corev1.ResourceCPU:    common.CalcResource(cpu, AgentResourceRatio, div),
			corev1.ResourceMemory: common.CalcResource(memory, AgentResourceRatio, div),
		}),
		Limits: common.NormalizeResourceList(corev1.ResourceList{
			corev1.ResourceCPU: common.CalcResource(cpu, AgentResourceRatio),
		}),
	}
}

// SetTopologySpreadConstraints spreads agent pods across nodes.
func (a *Agent) SetTopologySpreadConstraints() {
	a.TopologySpreadConstraints = setTopologySpreadConstraints(componentLabelAgent)
}

// SetPvEnable configures the agent for PV-backed operation.
func (a *Agent) SetPvEnable(sc, am, size string) {
	if a.Ngt == nil {
		a.Ngt = &AgentNgt{}
	}
	a.Ngt.EnableCopyOnWrite = new(true)
	a.Ngt.EnableInMemoryMode = new(false)
	a.Ngt.IndexPath = new(DefaultIndexPath)
	a.PersistentVolume = &AgentPersistentVolume{
		Enabled:      new(true),
		StorageClass: new(sc),
		Size:         new(size),
		AccessMode:   new(am),
	}
}

// --- Gateway (LB) ----------------------------------------------------------

// gatewayLbReplicaScale relates LB gateway replicas to the agent replica
// count: max = agents × scale, min = agents ÷ scale (both floored at 1).
const gatewayLbReplicaScale = 2

func (l *GatewayLb) getMaxReplica(ar int) int {
	if r := ar * gatewayLbReplicaScale; r >= 1 {
		return r
	}
	return 1
}

func (l *GatewayLb) getMinReplica(ar int) int {
	if r := ar / gatewayLbReplicaScale; r >= 1 {
		return r
	}
	return 1
}

// SetReplica derives LB replicas from the agent's replica count.
func (l *GatewayLb) SetReplica(a *Agent) {
	ar := 0
	if a != nil && a.MaxReplicas != nil {
		ar = *a.MaxReplicas
	}
	l.MinReplicas = new(l.getMinReplica(ar))
	l.MaxReplicas = new(l.getMaxReplica(ar))
}

// SetResources sets fixed compute resources for the LB gateway.
func (l *GatewayLb) SetResources() {
	l.Resources = fixedResources("200m", "150Mi", "2000m", "700Mi")
}

// SetTopologySpreadConstraints spreads LB gateway pods across nodes.
func (l *GatewayLb) SetTopologySpreadConstraints() {
	l.TopologySpreadConstraints = setTopologySpreadConstraints(componentLabelGatewayLb)
}

// --- Discoverer ------------------------------------------------------------

// ApplyDefaultsByKind fills in mode-specific defaults. The caller supplies the
// DaemonSet rolling-update knobs explicitly so this package stays free of
// configuration-source dependencies.
func (d *Discoverer) ApplyDefaultsByKind(daemonSetMaxSurge, daemonSetMaxUnavailable string) {
	kind := common.KindTypeDeployment
	if d.Kind != nil {
		kind = common.KindType(*d.Kind)
	}
	switch kind {
	case common.KindTypeDaemonSet:
		if d.ServiceType == nil {
			d.ServiceType = new(DiscovererServiceTypeNodePort)
		}
		if d.ExternalTrafficPolicy == nil {
			d.ExternalTrafficPolicy = new(string(corev1.ServiceExternalTrafficPolicyTypeLocal))
		}
		d.RollingUpdate = &RollingUpdate{
			MaxSurge:       new(daemonSetMaxSurge),
			MaxUnavailable: new(daemonSetMaxUnavailable),
		}
	case common.KindTypeDeployment:
		d.ServiceType = new(DiscovererServiceTypeClusterIP)
	case common.KindTypeStatefulSet:
		// The discoverer is never deployed as a StatefulSet; listed explicitly
		// so this switch stays exhaustive over common.KindType.
	}
}

// SetResources sets fixed compute resources for the discoverer.
func (d *Discoverer) SetResources() {
	d.Resources = fixedResources("200m", "65Mi", "600m", "200Mi")
}

// SetTopologySpreadConstraints spreads discoverer pods across nodes.
func (d *Discoverer) SetTopologySpreadConstraints() {
	d.TopologySpreadConstraints = setTopologySpreadConstraints(componentLabelDiscoverer)
}

// --- Manager (Index) -------------------------------------------------------

// SetResources sets fixed compute resources for the index manager.
func (i *ManagerIndex) SetResources() {
	i.Resources = fixedResources("200m", "80Mi", "1000m", "500Mi")
}

// SetTopologySpreadConstraints spreads index-manager pods across nodes.
func (i *ManagerIndex) SetTopologySpreadConstraints() {
	i.TopologySpreadConstraints = setTopologySpreadConstraints(componentLabelManagerIndex)
}

// --- Orchestrator ----------------------------------------------------------

// SetScaledResources populates the resource-related fields across every
// component that depends on the agent replica count or the node machine
// resources. Nil sub-specs are skipped, so callers only pay for the components
// they actually declared.
func (v *ValdRelease) SetScaledResources(ar int, am corev1.ResourceList, p ResourceParams) {
	if a := v.Spec.Agent; a != nil {
		a.SetReplica(ar, p.AgentPodsPerNode)
		a.SetResources(am, p.AgentPodsPerNode)
		a.SetTopologySpreadConstraints()
	}
	if g := v.Spec.Gateway; g != nil && g.Lb != nil {
		g.Lb.SetReplica(v.Spec.Agent)
		g.Lb.SetResources()
		g.Lb.SetTopologySpreadConstraints()
	}
	if d := v.Spec.Discoverer; d != nil {
		d.ApplyDefaultsByKind(p.DiscovererDSMaxSurge, p.DiscovererDSMaxUnavailable)
		d.SetResources()
		d.SetTopologySpreadConstraints()
	}
	if m := v.Spec.Manager; m != nil && m.Index != nil {
		if m.Index.Enabled != nil && *m.Index.Enabled {
			m.Index.SetResources()
			m.Index.SetTopologySpreadConstraints()
		}
	}
}
