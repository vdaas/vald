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
package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

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
			assert.Equal(t, tt.want, a.MinReplicas)
			assert.Equal(t, tt.want, a.MaxReplicas)
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
			mc := v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse(tt.nodeCPU),
				v1.ResourceMemory: resource.MustParse(tt.nodeRAM),
			}
			a := &Agent{}
			a.SetResources(mc, tt.podsPerNode)

			assert.NotNil(t, a.Resources)

			gotCPUReq := a.Resources.Requests[v1.ResourceCPU]
			gotCPULim := a.Resources.Limits[v1.ResourceCPU]
			gotRAMReq := a.Resources.Requests[v1.ResourceMemory]

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
			_, hasMemLim := a.Resources.Limits[v1.ResourceMemory]
			assert.False(t, hasMemLim, "memory limit must not be set (NGT index grows post-startup)")
		})
	}
}

// TestAgent_SetPvEnable verifies that calling SetPvEnable populates NGT
// settings and PersistentVolume from the caller-provided values. PV size
// resolution lives in the domain layer (see Domain.AgentPvSize tests).
func TestAgent_SetPvEnable(t *testing.T) {
	a := &Agent{}
	a.SetPvEnable("fast-ssd", "ReadWriteOnce", "6Gi")

	assert.True(t, a.NGT.EnableCopyOnWrite)
	assert.False(t, a.NGT.EnableInMemoryMode)
	assert.Equal(t, DefaultIndexPath, a.NGT.IndexPath)

	assert.NotNil(t, a.PersistentVolume)
	assert.True(t, a.PersistentVolume.Enabled)
	assert.Equal(t, "fast-ssd", a.PersistentVolume.StorageClass)
	assert.Equal(t, "ReadWriteOnce", a.PersistentVolume.AccessMode)
	assert.Equal(t, "6Gi", a.PersistentVolume.Size)
}

func TestAgent_SetTopologySpreadConstraints(t *testing.T) {
	a := &Agent{}
	a.SetTopologySpreadConstraints()

	assert.Len(t, a.TopologySpreadConstraints, 1)
	tsc := a.TopologySpreadConstraints[0]
	assert.Equal(t, int32(1), tsc.MaxSkew)
	assert.Equal(t, "kubernetes.io/hostname", tsc.TopologyKey)
	assert.Equal(t, v1.DoNotSchedule, tsc.WhenUnsatisfiable)
	assert.NotNil(t, tsc.LabelSelector)
	assert.Equal(t, "agent", tsc.LabelSelector.MatchLabels["app.kubernetes.io/component"])
}
