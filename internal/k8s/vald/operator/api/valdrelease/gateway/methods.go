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
package gateway

import (
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/agent"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/common"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (l *Lb) getMaxReplica(ar int) int {
	if r := ar * 2; r >= 1 {
		return r
	}
	return 1
}

func (l *Lb) getMinReplica(ar int) int {
	if r := ar / 2; r >= 1 {
		return r
	}
	return 1
}

func (l *Lb) SetReplica(a agent.Agent) {
	ar := a.MaxReplicas
	l.MinReplicas = l.getMinReplica(ar)
	l.MaxReplicas = l.getMaxReplica(ar)
}

func (l *Lb) SetResources() {
	l.Resources = &v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("200m"),
			v1.ResourceMemory: resource.MustParse("150Mi"),
		},
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("2000m"),
			v1.ResourceMemory: resource.MustParse("700Mi"),
		},
	}
}

func (l *Lb) SetTopologySpreadConstraints() {
	l.TopologySpreadConstraints = []v1.TopologySpreadConstraint{
		common.BuildTopologySpreadConstraint(componentLabelGatewayLb),
	}
}
