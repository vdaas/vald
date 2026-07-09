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

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "github.com/vdaas/vald/internal/k8s/vald/operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestResolveAgentNodePool(t *testing.T) {
	makeInfra := func(generalReplicas int, agentReplicas int, withAgent bool) v1.ValdOperatorReleaseInfra {
		pools := v1.NodePools{
			v1.NodePoolTypeGeneral: v1.NodePool{
				Name:            "general",
				Replicas:        generalReplicas,
				MachineResource: v1.MachineResource{Cpu: "4", Memory: "8Gi"},
			},
		}
		if withAgent {
			pools[v1.NodePoolTypeValdAgent] = v1.NodePool{
				Name:            "agent",
				Replicas:        agentReplicas,
				MachineResource: v1.MachineResource{Cpu: "16", Memory: "32Gi"},
			}
		}
		return v1.ValdOperatorReleaseInfra{NodePools: pools}
	}

	t.Run("agent pool present with replicas: use agent", func(t *testing.T) {
		got := resolveAgentNodePool(makeInfra(3, 2, true))
		assert.Equal(t, 2, got.NodeCount)
		assert.Equal(t, resource.MustParse("16"), got.MachineResource[corev1.ResourceCPU])
	})

	t.Run("agent pool present but replicas == 0: fall back to general", func(t *testing.T) {
		got := resolveAgentNodePool(makeInfra(3, 0, true))
		assert.Equal(t, 3, got.NodeCount)
		assert.Equal(t, resource.MustParse("4"), got.MachineResource[corev1.ResourceCPU])
	})

	t.Run("agent pool absent: fall back to general", func(t *testing.T) {
		got := resolveAgentNodePool(makeInfra(3, 0, false))
		assert.Equal(t, 3, got.NodeCount)
		assert.Equal(t, resource.MustParse("4"), got.MachineResource[corev1.ResourceCPU])
	})

	t.Run("no pools at all: empty spec without panicking", func(t *testing.T) {
		got := resolveAgentNodePool(v1.ValdOperatorReleaseInfra{})
		assert.Zero(t, got.NodeCount)
		assert.Nil(t, got.MachineResource)
	})

	t.Run("agent pool with zero replicas and general pool absent: empty spec", func(t *testing.T) {
		got := resolveAgentNodePool(v1.ValdOperatorReleaseInfra{
			NodePools: v1.NodePools{
				v1.NodePoolTypeValdAgent: v1.NodePool{
					Name:            "agent",
					Replicas:        0,
					MachineResource: v1.MachineResource{Cpu: "16", Memory: "32Gi"},
				},
			},
		})
		assert.Zero(t, got.NodeCount)
		assert.Nil(t, got.MachineResource)
	})
}

func TestAgentPvSize(t *testing.T) {
	const oneGi = int64(1) << 30
	const oneMi = int64(1) << 20

	tests := []struct {
		name         string
		memoryBytes  int64
		bufferRatio  float64
		minSizeBytes int64
		wantBytes    int64
	}{
		{"memory * ratio above min: ratio applies", 4 * oneGi, 1.5, oneGi, 6 * oneGi},
		{"memory * ratio below min: min applies", 100 * oneMi, 1.5, oneGi, oneGi},
		{"configurable min raises the floor", 100 * oneMi, 1.5, 2 * oneGi, 2 * oneGi},
		{"configurable ratio scales the calculation", 4 * oneGi, 2.0, oneGi, 8 * oneGi},
		{"rounds up to whole Gi", 1500 * oneMi, 1.0, oneGi, 2 * oneGi},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := agentPvSize(tt.memoryBytes, tt.bufferRatio, tt.minSizeBytes)
			gotQty := resource.MustParse(got)
			assert.Equal(t, tt.wantBytes, gotQty.Value(), "PV size in bytes")
		})
	}
}
