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
package v1

import (
	"github.com/vdaas/vald/internal/k8s/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// GroupVersion is group version used to register these objects.
	GroupVersion = schema.GroupVersion{Group: "vald.vdaas.org", Version: "v1"}

	// SchemeBuilder registers the item type via AddKnownTypes and the generic
	// list alias via AddListToScheme: generic instantiations carry mangled
	// reflect type names, so the list kind must be registered explicitly.
	SchemeBuilder = runtime.NewSchemeBuilder(func(s *runtime.Scheme) error {
		s.AddKnownTypes(GroupVersion, &ValdMirrorTarget{})
		resource.AddListToScheme[ValdMirrorTarget](s, GroupVersion, "ValdMirrorTargetList")
		metav1.AddToGroupVersion(s, GroupVersion)
		return nil
	})

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// ValdMirrorTarget is a mirror information.
type ValdMirrorTarget struct {
	resource.Base[ValdMirrorTarget, *ValdMirrorTarget] `json:"-"`

	Status            MirrorTargetStatus `json:"status"`
	metav1.TypeMeta   `                   json:",inline"`
	metav1.ObjectMeta `                   json:"metadata"`
	Spec              MirrorTargetSpec `json:"spec"`
}

// ValdMirrorTargetList is the whole list of all ValdMirrorTarget which have
// been registered with master. The whole list kind is derived generically
// from the item type.
type ValdMirrorTargetList = resource.List[ValdMirrorTarget, *ValdMirrorTarget]

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *ValdMirrorTarget) DeepCopyInto(out *ValdMirrorTarget) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
}

// DeepCopy returns a deep copy of the receiver. It disambiguates the promoted
// Base.DeepCopy from the equally promoted metav1.ObjectMeta.DeepCopy.
func (in *ValdMirrorTarget) DeepCopy() *ValdMirrorTarget {
	return in.Base.DeepCopy()
}

// MirrorTargetSpec is a description of a ValdMirrorTarget.
type MirrorTargetSpec struct {
	Colocation string       `json:"colocation,omitempty"`
	Target     MirrorTarget `json:"target"`
}

// MirrorTarget is a target information.
type MirrorTarget struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

// MirrorTargetStatus is a status of ValdMirrorTarget.
type MirrorTargetStatus struct {
	LastTransitionTime metav1.Time       `json:"lastTransitionTime"`
	Phase              MirrorTargetPhase `json:"phase,omitempty"`
}

// MirrorTargetPhase is a label for the condition of a ValdMirrorTarget at the current time.
type MirrorTargetPhase string

const (
	// MirrorTargetConnected means that the ValdMirrorTarget has been accepted by the system.
	MirrorTargetPending = MirrorTargetPhase("Pending")

	// MirrorTargetConnected means that the target was connected.
	MirrorTargetConnected = MirrorTargetPhase("Connected")

	// MirrorTargetDisconnected means that the target was disconnected.
	MirrorTargetDisconnected = MirrorTargetPhase("Disconnected")

	// MirrorTargetUnknown means that for some reason the state of the ValdMirrorTarget could not be obtained.
	MirrorTargetUnknown = MirrorTargetPhase("Unknown")
)
