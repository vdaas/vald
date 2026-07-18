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

package valdrelease

import (
	"github.com/vdaas/vald/internal/k8s/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupVersion / GVK identify the ValdRelease custom resource.
//
// // standard controller-runtime scheme-registration pattern; GroupVersion/GVK/AddToScheme are the public API consumed by every reconciler in this group
var (
	GroupVersion = schema.GroupVersion{Group: "vald.vdaas.org", Version: "v1"}
	GVK          = schema.GroupVersionKind{Group: GroupVersion.Group, Version: GroupVersion.Version, Kind: "ValdRelease"}

	// SchemeBuilder registers the item type via AddKnownTypes and the generic
	// list alias via AddListToScheme: generic instantiations carry mangled
	// reflect type names, so the list kind must be registered explicitly.
	SchemeBuilder = runtime.NewSchemeBuilder(func(s *runtime.Scheme) error {
		s.AddKnownTypes(GroupVersion, &ValdRelease{})
		resource.AddListToScheme[ValdRelease](s, GroupVersion, "ValdReleaseList")
		metav1.AddToGroupVersion(s, GroupVersion)
		return nil
	})

	// AddToScheme registers the ValdRelease types with a runtime.Scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// Status is the ValdRelease lifecycle status string.
type Status string

// VrsStatus is the observed status of a ValdRelease.
type VrsStatus struct {
	Status    Status           `json:"status,omitempty"`
	Condition metav1.Condition `json:"condition"`
}

// DeepCopyInto copies the receiver into out.
func (in *VrsStatus) DeepCopyInto(out *VrsStatus) {
	*out = *in
	in.Condition.DeepCopyInto(&out.Condition)
}

// ValdRelease is the CRD object. Its Spec is the generated chart-values type
// (Values), derived from charts/vald/values.schema.json instead of a
// hand-maintained mirror. The embedded resource.Base promotes DeepCopy /
// DeepCopyObject generically; DeepCopyInto is provided by the generator.
//
// // fieldalignment: resource.Base must stay the zero-size first field (contract documented on Base) and the K8s API convention fixes the
// TypeMeta/ObjectMeta/Spec/Status order
type ValdRelease struct { //nolint:tagliatelle // generic embed field name confuses the linter's tag-name check
	resource.Base[ValdRelease, *ValdRelease] `json:"-"`

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"` //nolint:tagliatelle // fixed by the Kubernetes object API wire format, not renameable

	Spec   Values    `json:"spec"`
	Status VrsStatus `json:"status"`
}

// DeepCopyInto copies the receiver into out.
func (in *ValdRelease) DeepCopyInto(out *ValdRelease) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// ValdReleaseList is a list of ValdRelease, derived generically from the item
// type.
type ValdReleaseList = resource.List[ValdRelease, *ValdRelease]

// Compile-time proof that ValdRelease is a typed runtime.Object.
var _ runtime.Object = &ValdRelease{}
