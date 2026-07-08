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

package resource

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// List is a fully generic Kubernetes list object for the item type T. Because
// the structure of every list kind is identical (TypeMeta + ListMeta +
// Items), the whole type including its deepcopy behavior can be derived
// generically; API packages only declare an alias:
//
//	type ValdOperatorReleaseList = resource.List[ValdOperatorRelease, *ValdOperatorRelease]
//
// NOTE: generic instantiations carry mangled reflect type names, so list
// aliases MUST be registered on the scheme with an explicit kind via
// AddKnownTypeWithName (see AddListToScheme) instead of AddKnownTypes.
type List[T any, PT DeepCopyIntoer[T]] struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []T `json:"items"`
}

// DeepCopyInto copies the receiver into out.
func (in *List[T, PT]) DeepCopyInto(out *List[T, PT]) {
	if in == nil || out == nil {
		return
	}
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	out.Items = CopySliceInto[T, PT](in.Items)
}

// DeepCopy returns a deep copy of the receiver.
func (in *List[T, PT]) DeepCopy() *List[T, PT] {
	if in == nil {
		return nil
	}
	out := new(List[T, PT])
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject returns a deep copy of the receiver as a runtime.Object.
func (in *List[T, PT]) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	return in.DeepCopy()
}

// AddListToScheme registers a generic list instantiation on the scheme with
// an explicit kind name. AddKnownTypes cannot be used for generic list
// aliases because it derives the kind from the reflect type name, which is
// mangled for generic instantiations (e.g. "List[...]").
func AddListToScheme[T any, PT DeepCopyIntoer[T]](
	s *runtime.Scheme, gv schema.GroupVersion, kind string,
) {
	s.AddKnownTypeWithName(gv.WithKind(kind), &List[T, PT]{})
}
