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

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/agent"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/discoverer"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/gateway"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/manager"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func newTestValdRelease() *ValdRelease {
	return &ValdRelease{
		Spec: Spec{
			Agent:      agent.Agent{},
			Gateway:    gateway.Gateway{},
			Discoverer: discoverer.Discoverer{Kind: "Deployment"},
			Manager: manager.Manager{
				Index: manager.Index{
					Enabled: true,
					Creator: &manager.Creator{Enabled: true},
					Saver:   &manager.Saver{Enabled: true},
				},
			},
		},
	}
}

func TestSetRelationalResources(t *testing.T) {
	mc := v1.ResourceList{
		v1.ResourceCPU:    resource.MustParse("8"),
		v1.ResourceMemory: resource.MustParse("16Gi"),
	}

	tests := []struct {
		name         string
		indexEnabled bool
	}{
		{"index enabled: manager index resources are set", true},
		{"index disabled: manager index resources stay nil", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr := newTestValdRelease()
			vr.Spec.Manager.Index.Enabled = tt.indexEnabled

			vr.SetRelationalResources(3, mc, ResourceParams{
				AgentPodsPerNode:           2,
				DiscovererDSMaxSurge:       "30%",
				DiscovererDSMaxUnavailable: "0%",
			})

			assert.NotZero(t, vr.Spec.Agent.MinReplicas)
			assert.NotZero(t, vr.Spec.Agent.MaxReplicas)
			assert.NotNil(t, vr.Spec.Agent.Resources)
			assert.NotZero(t, vr.Spec.Gateway.Lb.MinReplicas)
			assert.NotZero(t, vr.Spec.Gateway.Lb.MaxReplicas)
			assert.NotNil(t, vr.Spec.Gateway.Lb.Resources)
			assert.NotNil(t, vr.Spec.Discoverer.Resources)
			if tt.indexEnabled {
				assert.NotNil(t, vr.Spec.Manager.Index.Resources)
				assert.NotEmpty(t, vr.Spec.Manager.Index.TopologySpreadConstraints)
			} else {
				assert.Nil(t, vr.Spec.Manager.Index.Resources)
				assert.Empty(t, vr.Spec.Manager.Index.TopologySpreadConstraints)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
