// +k8s:deepcopy-gen=package
// +groupName=vald.vdaas.org
package gateway

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/defaults"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

const (
	componentLabelGatewayLb = "gateway-lb"
)

// +schemagen:begin

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
	Enabled        bool              `json:"enabled"`
	Annotations    map[string]string `json:"annotations,omitempty"`
	Host           string            `json:"host"`
	DefaultBackend DefaultBackend    `json:"defaultBackend,omitempty"`
	PathType       v1beta1.PathType  `json:"pathType,omitempty"`
	ServicePort    string            `json:"servicePort"`
}

type DefaultBackend struct {
	Enabled bool `json:"enabled"`
}
type GatewayConfig struct {
	IndexReplica int `json:"index_replica"`
}
