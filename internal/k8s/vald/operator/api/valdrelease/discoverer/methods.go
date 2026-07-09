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
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/common"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// ApplyDefaultsByKind fills in mode-specific defaults. The caller supplies the
// DaemonSet rolling-update knobs explicitly so this package stays free of
// configuration-source dependencies.
func (d *Discoverer) ApplyDefaultsByKind(daemonSetMaxSurge, daemonSetMaxUnavailable string) {
	switch d.Kind {
	case common.KindTypeDaemonSet:
		if d.ServiceType == "" {
			d.ServiceType = v1.ServiceTypeNodePort
		}
		if d.ExternalTrafficPolicy == "" {
			d.ExternalTrafficPolicy = v1.ServiceExternalTrafficPolicyTypeLocal
		}
		ms := intstr.FromString(daemonSetMaxSurge)
		mu := intstr.FromString(daemonSetMaxUnavailable)
		d.RollingUpdate = &RollingUpdateValdRelease{
			MaxSurge:       &ms,
			MaxUnavailable: &mu,
		}
	case common.KindTypeDeployment:
		d.ServiceType = v1.ServiceTypeClusterIP
	}
}

func (d *Discoverer) SetResources() {
	d.Resources = &v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("200m"),
			v1.ResourceMemory: resource.MustParse("65Mi"),
		},
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("600m"),
			v1.ResourceMemory: resource.MustParse("200Mi"),
		},
	}
}

func (d *Discoverer) SetTopologySpreadConstraints() {
	d.TopologySpreadConstraints = []v1.TopologySpreadConstraint{
		common.BuildTopologySpreadConstraint(componentLabelDiscoverer),
	}
}
