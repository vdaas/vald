// +k8s:deepcopy-gen=package
// +groupName=vald.vdaas.org
package discoverer

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/common"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/defaults"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	componentLabelDiscoverer = "discoverer"
)

// Note: This struct is same as `appsv1.RollingUpdateDaemonSet` in k8s.io/api/apps/v1.
// Note: However, this is a special structure for vrs, so we define it here just in case.
type RollingUpdateValdreelase struct {
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty" protobuf:"bytes,1,opt,name=maxUnavailable"`
	MaxSurge       *intstr.IntOrString `json:"maxSurge,omitempty" protobuf:"bytes,2,opt,name=maxSurge"`
}

// +schemagen:begin

type Discoverer struct {
	Logging                   *defaults.Logging                   `json:"logging,omitempty"`
	Affinity                  *v1.Affinity                        `json:"affinity,omitempty"`
	NodeSelector              map[string]string                   `json:"nodeSelector,omitempty"`
	Tolerations               *[]v1.Toleration                    `json:"tolerations,omitempty"`
	ClusterRole               ClusterRole                         `json:"clusterRole,omitempty"`
	ClusterRoleBinding        ClusterRoleBinding                  `json:"clusterRoleBinding,omitempty"`
	Kind                      common.KindType                     `json:"kind,omitempty"`
	ServiceType               v1.ServiceType                      `json:"serviceType,omitempty"`
	ExternalTrafficPolicy     v1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
	Resources                 *v1.ResourceRequirements            `json:"resources,omitempty"`
	TopologySpreadConstraints []v1.TopologySpreadConstraint       `json:"topologySpreadConstraints,omitempty"`
	RollingUpdate             *RollingUpdateValdreelase           `json:"rollingUpdate,omitempty"`
}

type ClusterRole struct {
	Enabled bool   `json:"enabled,omitempty"`
	Name    string `json:"name,omitempty"`
}

type ClusterRoleBinding struct {
	Enabled bool   `json:"enabled,omitempty"`
	Name    string `json:"name,omitempty"`
}
