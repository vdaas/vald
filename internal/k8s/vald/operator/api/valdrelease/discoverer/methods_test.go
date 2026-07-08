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
package discoverer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/common"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestDiscoverer_SetResources(t *testing.T) {
	d := &Discoverer{}
	d.SetResources()

	assert.NotNil(t, d.Resources)
	assert.Equal(t, resource.MustParse("200m"), d.Resources.Requests[v1.ResourceCPU])
	assert.Equal(t, resource.MustParse("65Mi"), d.Resources.Requests[v1.ResourceMemory])
	assert.Equal(t, resource.MustParse("600m"), d.Resources.Limits[v1.ResourceCPU])
	assert.Equal(t, resource.MustParse("200Mi"), d.Resources.Limits[v1.ResourceMemory])
}

func TestDiscoverer_SetTopologySpreadConstraints(t *testing.T) {
	d := &Discoverer{}
	d.SetTopologySpreadConstraints()

	assert.Len(t, d.TopologySpreadConstraints, 1)
	tsc := d.TopologySpreadConstraints[0]
	assert.Equal(t, int32(1), tsc.MaxSkew)
	assert.Equal(t, "kubernetes.io/hostname", tsc.TopologyKey)
	assert.Equal(t, v1.DoNotSchedule, tsc.WhenUnsatisfiable)
	assert.NotNil(t, tsc.LabelSelector)
	assert.Equal(t, "discoverer", tsc.LabelSelector.MatchLabels["app.kubernetes.io/component"])
}

func TestDiscoverer_ApplyDefaultsByKind(t *testing.T) {
	const (
		maxSurge       = "30%"
		maxUnavailable = "0%"
	)

	tests := []struct {
		name                  string
		kind                  common.KindType
		existingServiceType   v1.ServiceType
		existingTrafficPolicy v1.ServiceExternalTrafficPolicyType
		wantServiceType       v1.ServiceType
		wantTrafficPolicy     v1.ServiceExternalTrafficPolicyType
		wantRollingUpdate     bool
	}{
		{
			name:              "DaemonSet sets NodePort and Local defaults",
			kind:              common.KindTypeDaemonSet,
			wantServiceType:   v1.ServiceTypeNodePort,
			wantTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeLocal,
			wantRollingUpdate: true,
		},
		{
			name:                "DaemonSet does not override existing ServiceType",
			kind:                common.KindTypeDaemonSet,
			existingServiceType: v1.ServiceTypeClusterIP,
			wantServiceType:     v1.ServiceTypeClusterIP,
			wantTrafficPolicy:   v1.ServiceExternalTrafficPolicyTypeLocal,
			wantRollingUpdate:   true,
		},
		{
			name:                  "DaemonSet does not override existing ExternalTrafficPolicy",
			kind:                  common.KindTypeDaemonSet,
			existingTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
			wantServiceType:       v1.ServiceTypeNodePort,
			wantTrafficPolicy:     v1.ServiceExternalTrafficPolicyTypeCluster,
			wantRollingUpdate:     true,
		},
		{
			name:              "Deployment sets ClusterIP",
			kind:              common.KindTypeDeployment,
			wantServiceType:   v1.ServiceTypeClusterIP,
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
			d := &Discoverer{
				Kind:                  tt.kind,
				ServiceType:           tt.existingServiceType,
				ExternalTrafficPolicy: tt.existingTrafficPolicy,
			}
			d.ApplyDefaultsByKind(maxSurge, maxUnavailable)

			assert.Equal(t, tt.wantServiceType, d.ServiceType)
			assert.Equal(t, tt.wantTrafficPolicy, d.ExternalTrafficPolicy)

			if tt.wantRollingUpdate {
				assert.NotNil(t, d.RollingUpdate)
				assert.Equal(t, maxSurge, d.RollingUpdate.MaxSurge.String())
				assert.Equal(t, maxUnavailable, d.RollingUpdate.MaxUnavailable.String())
			} else {
				assert.Nil(t, d.RollingUpdate)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
