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

// Package discoverer defines the discoverer component model of the generated ValdRelease.
package discoverer

import (
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/common"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/defaults"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	componentLabelDiscoverer = "discoverer"
)

type Discoverer struct {
	Logging                   *defaults.Logging                   `json:"logging,omitempty"`
	Affinity                  *v1.Affinity                        `json:"affinity,omitempty"`
	NodeSelector              map[string]string                   `json:"nodeSelector,omitempty"`
	Tolerations               *[]v1.Toleration                    `json:"tolerations,omitempty"`
	ClusterRole               ClusterRole                         `json:"clusterRole"`
	ClusterRoleBinding        ClusterRoleBinding                  `json:"clusterRoleBinding"`
	Kind                      common.KindType                     `json:"kind,omitempty"`
	ServiceType               v1.ServiceType                      `json:"serviceType,omitempty"`
	ExternalTrafficPolicy     v1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
	Resources                 *v1.ResourceRequirements            `json:"resources,omitempty"`
	TopologySpreadConstraints []v1.TopologySpreadConstraint       `json:"topologySpreadConstraints,omitempty"`
	RollingUpdate             *RollingUpdateValdRelease           `json:"rollingUpdate,omitempty"`
}

type ClusterRole struct {
	Enabled bool   `json:"enabled,omitempty"`
	Name    string `json:"name,omitempty"`
}

type ClusterRoleBinding struct {
	Enabled bool   `json:"enabled,omitempty"`
	Name    string `json:"name,omitempty"`
}

// Note: This struct is same as `appsv1.RollingUpdateDaemonSet` in k8s.io/api/apps/v1.
// Note: However, this is a special structure for vrs, so we define it here just in case.
type RollingUpdateValdRelease struct {
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty" protobuf:"bytes,1,opt,name=maxUnavailable"`
	MaxSurge       *intstr.IntOrString `json:"maxSurge,omitempty"       protobuf:"bytes,2,opt,name=maxSurge"`
}
