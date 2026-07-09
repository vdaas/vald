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
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/defaults"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

const (
	componentLabelGatewayLb = "gateway-lb"
)

type Gateway struct {
	Lb Lb `json:"lb"`
}

type Lb struct {
	Logging                   *defaults.Logging             `json:"logging,omitempty"`
	Affinity                  *v1.Affinity                  `json:"affinity,omitempty"`
	NodeSelector              map[string]string             `json:"nodeSelector,omitempty"`
	Tolerations               *[]v1.Toleration              `json:"tolerations,omitempty"`
	Ingress                   Ingress                       `json:"ingress"`
	MinReplicas               int                           `json:"minReplicas"`
	MaxReplicas               int                           `json:"maxReplicas"`
	Hpa                       *Hpa                          `json:"hpa,omitempty"`
	GatewayConfig             GatewayConfig                 `json:"gateway_config"`
	ServiceType               v1.ServiceType                `json:"serviceType,omitempty"`
	Resources                 *v1.ResourceRequirements      `json:"resources,omitempty"`
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

type Hpa struct {
	TargetCPUUtilizationPercentage int `json:"targetCPUUtilizationPercentage,omitempty"`
}

type Ingress struct {
	Enabled        bool                  `json:"enabled"`
	Annotations    map[string]string     `json:"annotations,omitempty"`
	Host           string                `json:"host"`
	DefaultBackend DefaultBackend        `json:"defaultBackend"`
	PathType       networkingv1.PathType `json:"pathType,omitempty"`
	ServicePort    string                `json:"servicePort"`
}

type DefaultBackend struct {
	Enabled bool `json:"enabled"`
}
type GatewayConfig struct {
	IndexReplica int `json:"index_replica"`
}
