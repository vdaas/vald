// +k8s:deepcopy-gen=package
// +groupName=vald.vdaas.org
package valdrelease

import (
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/agent"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/defaults"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/discoverer"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/gateway"
	"github.com/vdaas/vald/hack/vald-operator/internal/pkg/api/valdrelease/manager"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "vald.vdaas.org", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

var GVK = schema.GroupVersionKind{
	Group:   GroupVersion.Group,
	Version: GroupVersion.Version,
	Kind:    "ValdRelease",
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ValdRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Spec      `json:"spec,omitempty"`
	Status VrsStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ValdReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ValdRelease `json:"items"`
}

type Spec struct {
	Defaults   defaults.Defaults     `json:"defaults"`
	Gateway    gateway.Gateway       `json:"gateway"`
	Agent      agent.Agent           `json:"agent"`
	Manager    manager.Manager       `json:"manager"`
	Discoverer discoverer.Discoverer `json:"discoverer"`
}

type Status string

type VrsStatus struct {
	Status    Status           `json:"status,omitempty"`
	Condition metav1.Condition `json:"condition,omitempty"`
}

func init() {
	SchemeBuilder.Register(&ValdRelease{}, &ValdReleaseList{})
}
