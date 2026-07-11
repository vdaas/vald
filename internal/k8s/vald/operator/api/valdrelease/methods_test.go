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
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// --- Agent -----------------------------------------------------------------

func TestAgent_SetReplica(t *testing.T) {
	const podsPerNode = 2
	tests := []struct {
		name         string
		nodeReplicas int
		want         int
	}{
		{"2 nodes × 2 pods", 2, 4},
		{"3 nodes × 2 pods", 3, 6},
		{"1 node × 2 pods", 1, 2},
		{"0 nodes", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{}
			a.SetReplica(tt.nodeReplicas, podsPerNode)
			assert.Equal(t, tt.want, *a.MinReplicas)
			assert.Equal(t, tt.want, *a.MaxReplicas)
		})
	}
}

// TestAgent_SetResources_Values verifies exact resource values for representative node specs.
//
// Formula:
//
//	CPU request  = nodeCPU_cores × ResourceRatio / podsPerNode  [millicores]
//	CPU limit    = nodeCPU_cores × ResourceRatio                [millicores]
//	RAM request  = nodeRAM_bytes × ResourceRatio / podsPerNode  [rounded up to M]
//	RAM limit    = NOT SET (NGT index grows post-startup; hard limit causes OOM kills)
const testNodeRAM = "10000M"

func TestAgent_SetResources_Values(t *testing.T) {
	tests := []struct {
		name            string
		podsPerNode     int
		nodeCPU         string
		nodeRAM         string
		wantCPUReqMilli int64
		wantCPULimMilli int64
		wantRAMReqBytes int64
	}{
		{
			// 16 × 0.6 / 2 = 4800m req, 9600m lim; 10000M × 0.6 / 2 = 3000M req
			name:        "16CPU 10000M node, 2 pods/node",
			podsPerNode: 2, nodeCPU: "16", nodeRAM: testNodeRAM,
			wantCPUReqMilli: 4800, wantCPULimMilli: 9600,
			wantRAMReqBytes: 3_000_000_000,
		},
		{
			// 8 × 0.6 / 2 = 2400m req, 4800m lim
			name:        "8CPU 10000M node, 2 pods/node",
			podsPerNode: 2, nodeCPU: "8", nodeRAM: testNodeRAM,
			wantCPUReqMilli: 2400, wantCPULimMilli: 4800,
			wantRAMReqBytes: 3_000_000_000,
		},
		{
			// podsPerNode=1: req == lim
			name:        "16CPU node, 1 pod/node",
			podsPerNode: 1, nodeCPU: "16", nodeRAM: testNodeRAM,
			wantCPUReqMilli: 9600, wantCPULimMilli: 9600,
			wantRAMReqBytes: 6_000_000_000,
		},
		{
			// podsPerNode=4: req = lim / 4
			name:        "16CPU node, 4 pods/node",
			podsPerNode: 4, nodeCPU: "16", nodeRAM: testNodeRAM,
			wantCPUReqMilli: 2400, wantCPULimMilli: 9600,
			wantRAMReqBytes: 1_500_000_000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(tt.nodeCPU),
				corev1.ResourceMemory: resource.MustParse(tt.nodeRAM),
			}
			a := &Agent{}
			a.SetResources(mc, tt.podsPerNode)

			assert.NotNil(t, a.Resources)

			gotCPUReq := a.Resources.Requests[corev1.ResourceCPU]
			gotCPULim := a.Resources.Limits[corev1.ResourceCPU]
			gotRAMReq := a.Resources.Requests[corev1.ResourceMemory]

			assert.Equal(t, tt.wantCPUReqMilli, gotCPUReq.MilliValue(), "CPU request (milli)")
			assert.Equal(t, tt.wantCPULimMilli, gotCPULim.MilliValue(), "CPU limit (milli)")
			assert.Equal(t, tt.wantRAMReqBytes, gotRAMReq.Value(), "RAM request (bytes)")

			// Structural invariant: limit covers all pods on the node, request is per-pod.
			assert.Equal(
				t,
				gotCPUReq.MilliValue()*int64(tt.podsPerNode),
				gotCPULim.MilliValue(),
				"CPU limit must equal request × podsPerNode",
			)

			// Memory limit must NOT be set.
			_, hasMemLim := a.Resources.Limits[corev1.ResourceMemory]
			assert.False(t, hasMemLim, "memory limit must not be set (NGT index grows post-startup)")
		})
	}
}

// TestAgent_SetPvEnable verifies that calling SetPvEnable populates NGT
// settings and PersistentVolume from the caller-provided values.
func TestAgent_SetPvEnable(t *testing.T) {
	a := &Agent{}
	a.SetPvEnable("fast-ssd", "ReadWriteOnce", "6Gi")

	assert.True(t, *a.Ngt.EnableCopyOnWrite)
	assert.False(t, *a.Ngt.EnableInMemoryMode)
	assert.Equal(t, DefaultIndexPath, *a.Ngt.IndexPath)

	assert.NotNil(t, a.PersistentVolume)
	assert.True(t, *a.PersistentVolume.Enabled)
	assert.Equal(t, "fast-ssd", *a.PersistentVolume.StorageClass)
	assert.Equal(t, "ReadWriteOnce", *a.PersistentVolume.AccessMode)
	assert.Equal(t, "6Gi", *a.PersistentVolume.Size)
}

func TestAgent_SetTopologySpreadConstraints(t *testing.T) {
	a := &Agent{}
	a.SetTopologySpreadConstraints()

	assert.NotNil(t, a.TopologySpreadConstraints)
	assert.Len(t, *a.TopologySpreadConstraints, 1)
	tsc := (*a.TopologySpreadConstraints)[0]
	assert.Equal(t, int32(1), tsc.MaxSkew)
	assert.Equal(t, "kubernetes.io/hostname", tsc.TopologyKey)
	assert.Equal(t, corev1.DoNotSchedule, tsc.WhenUnsatisfiable)
	assert.NotNil(t, tsc.LabelSelector)
	assert.Equal(t, "agent", tsc.LabelSelector.MatchLabels["app.kubernetes.io/component"])
}

// --- Gateway (LB) ----------------------------------------------------------

func TestGatewayLb_SetReplica(t *testing.T) {
	tests := []struct {
		name     string
		agentMax int
		wantMin  int
		wantMax  int
	}{
		{"6 agent replicas", 6, 3, 12},
		{"1 agent replica → min stays 1", 1, 1, 2},
		{"2 agent replicas", 2, 1, 4},
		// ar=0: both getMaxReplica(0)=1, getMinReplica(0)=1
		{"0 agent replicas → floor to 1", 0, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := &GatewayLb{}
			a := &Agent{MaxReplicas: new(tt.agentMax)}
			lb.SetReplica(a)
			assert.Equal(t, tt.wantMin, *lb.MinReplicas)
			assert.Equal(t, tt.wantMax, *lb.MaxReplicas)
		})
	}
}

func TestGatewayLb_SetResources(t *testing.T) {
	lb := &GatewayLb{}
	lb.SetResources()

	assert.NotNil(t, lb.Resources)
	assert.NotEmpty(t, lb.Resources.Requests[corev1.ResourceCPU])
	assert.NotEmpty(t, lb.Resources.Requests[corev1.ResourceMemory])
	assert.NotEmpty(t, lb.Resources.Limits[corev1.ResourceCPU])
	assert.NotEmpty(t, lb.Resources.Limits[corev1.ResourceMemory])
}

func TestGatewayLb_SetTopologySpreadConstraints(t *testing.T) {
	lb := &GatewayLb{}
	lb.SetTopologySpreadConstraints()

	assert.NotNil(t, lb.TopologySpreadConstraints)
	assert.Len(t, *lb.TopologySpreadConstraints, 1)
	assert.Equal(t, "gateway-lb", (*lb.TopologySpreadConstraints)[0].LabelSelector.MatchLabels["app.kubernetes.io/component"])
}

// --- Discoverer ------------------------------------------------------------

func TestDiscoverer_SetResources(t *testing.T) {
	d := &Discoverer{}
	d.SetResources()

	assert.NotNil(t, d.Resources)
	assert.Equal(t, resource.MustParse("200m"), d.Resources.Requests[corev1.ResourceCPU])
	assert.Equal(t, resource.MustParse("65Mi"), d.Resources.Requests[corev1.ResourceMemory])
	assert.Equal(t, resource.MustParse("600m"), d.Resources.Limits[corev1.ResourceCPU])
	assert.Equal(t, resource.MustParse("200Mi"), d.Resources.Limits[corev1.ResourceMemory])
}

func TestDiscoverer_SetTopologySpreadConstraints(t *testing.T) {
	d := &Discoverer{}
	d.SetTopologySpreadConstraints()

	assert.NotNil(t, d.TopologySpreadConstraints)
	assert.Len(t, *d.TopologySpreadConstraints, 1)
	tsc := (*d.TopologySpreadConstraints)[0]
	assert.Equal(t, int32(1), tsc.MaxSkew)
	assert.Equal(t, "kubernetes.io/hostname", tsc.TopologyKey)
	assert.Equal(t, corev1.DoNotSchedule, tsc.WhenUnsatisfiable)
	assert.NotNil(t, tsc.LabelSelector)
	assert.Equal(t, "discoverer", tsc.LabelSelector.MatchLabels["app.kubernetes.io/component"])
}

func TestDiscoverer_ApplyDefaultsByKind(t *testing.T) {
	const (
		maxSurge       = "30%"
		maxUnavailable = "0%"
	)

	// Empty string in the existing/want columns represents a nil pointer on the
	// generated type (the field is omitted).
	tests := []struct {
		name                  string
		kind                  string
		existingServiceType   string
		existingTrafficPolicy string
		wantServiceType       string
		wantTrafficPolicy     string
		wantRollingUpdate     bool
	}{
		{
			name:              "DaemonSet sets NodePort and Local defaults",
			kind:              string(DiscovererKindDaemonSet),
			wantServiceType:   string(corev1.ServiceTypeNodePort),
			wantTrafficPolicy: string(corev1.ServiceExternalTrafficPolicyTypeLocal),
			wantRollingUpdate: true,
		},
		{
			name:                "DaemonSet does not override existing ServiceType",
			kind:                string(DiscovererKindDaemonSet),
			existingServiceType: string(corev1.ServiceTypeClusterIP),
			wantServiceType:     string(corev1.ServiceTypeClusterIP),
			wantTrafficPolicy:   string(corev1.ServiceExternalTrafficPolicyTypeLocal),
			wantRollingUpdate:   true,
		},
		{
			name:                  "DaemonSet does not override existing ExternalTrafficPolicy",
			kind:                  string(DiscovererKindDaemonSet),
			existingTrafficPolicy: string(corev1.ServiceExternalTrafficPolicyTypeCluster),
			wantServiceType:       string(corev1.ServiceTypeNodePort),
			wantTrafficPolicy:     string(corev1.ServiceExternalTrafficPolicyTypeCluster),
			wantRollingUpdate:     true,
		},
		{
			name:              "Deployment sets ClusterIP",
			kind:              string(DiscovererKindDeployment),
			wantServiceType:   string(corev1.ServiceTypeClusterIP),
			wantRollingUpdate: false,
		},
		{
			name:              "unknown kind makes no change",
			kind:              "Unknown",
			wantServiceType:   "",
			wantTrafficPolicy: "",
			wantRollingUpdate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discoverer{Kind: new(DiscovererKind(tt.kind))}
			if tt.existingServiceType != "" {
				d.ServiceType = new(DiscovererServiceType(tt.existingServiceType))
			}
			if tt.existingTrafficPolicy != "" {
				d.ExternalTrafficPolicy = new(tt.existingTrafficPolicy)
			}
			d.ApplyDefaultsByKind(maxSurge, maxUnavailable)

			if tt.wantServiceType == "" {
				assert.Nil(t, d.ServiceType)
			} else {
				assert.NotNil(t, d.ServiceType)
				assert.Equal(t, tt.wantServiceType, string(*d.ServiceType))
			}

			if tt.wantTrafficPolicy == "" {
				assert.Nil(t, d.ExternalTrafficPolicy)
			} else {
				assert.NotNil(t, d.ExternalTrafficPolicy)
				assert.Equal(t, tt.wantTrafficPolicy, *d.ExternalTrafficPolicy)
			}

			if tt.wantRollingUpdate {
				assert.NotNil(t, d.RollingUpdate)
				assert.Equal(t, maxSurge, *d.RollingUpdate.MaxSurge)
				assert.Equal(t, maxUnavailable, *d.RollingUpdate.MaxUnavailable)
			} else {
				assert.Nil(t, d.RollingUpdate)
			}
		})
	}
}

// --- Manager (Index) -------------------------------------------------------

func TestManagerIndex_SetResources(t *testing.T) {
	idx := &ManagerIndex{}
	idx.SetResources()

	assert.NotNil(t, idx.Resources)
	assert.NotEmpty(t, idx.Resources.Requests[corev1.ResourceCPU])
	assert.NotEmpty(t, idx.Resources.Requests[corev1.ResourceMemory])
	assert.NotEmpty(t, idx.Resources.Limits[corev1.ResourceCPU])
	assert.NotEmpty(t, idx.Resources.Limits[corev1.ResourceMemory])
}

func TestManagerIndex_SetTopologySpreadConstraints(t *testing.T) {
	idx := &ManagerIndex{}
	idx.SetTopologySpreadConstraints()

	assert.NotNil(t, idx.TopologySpreadConstraints)
	assert.Len(t, *idx.TopologySpreadConstraints, 1)
	assert.Equal(t, "manager-index", (*idx.TopologySpreadConstraints)[0].LabelSelector.MatchLabels["app.kubernetes.io/component"])
}
