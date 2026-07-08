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

package v1

import "github.com/vdaas/vald/internal/k8s/resource"

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *Agent) DeepCopyInto(out *Agent) {
	*out = *in
	out.PersistentVolume = resource.CopyPtr(in.PersistentVolume)
}

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *Gateway) DeepCopyInto(out *Gateway) {
	*out = *in
	out.Ingress = resource.CopyPtr(in.Ingress)
}

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *ValdOperatorRelease) DeepCopyInto(out *ValdOperatorRelease) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy returns a deep copy of the receiver. It disambiguates the promoted
// Base.DeepCopy from the equally promoted metav1.ObjectMeta.DeepCopy.
func (in *ValdOperatorRelease) DeepCopy() *ValdOperatorRelease {
	return in.Base.DeepCopy()
}

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *ValdOperatorReleaseInfra) DeepCopyInto(out *ValdOperatorReleaseInfra) {
	*out = *in
	out.Clusters = resource.CopySlice(in.Clusters)
	out.NodePools = resource.CopyMap(in.NodePools)
}

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *ValdOperatorReleaseSpec) DeepCopyInto(out *ValdOperatorReleaseSpec) {
	*out = *in
	out.Infrastructure = resource.CopySliceInto(in.Infrastructure)
	in.VectorEngine.DeepCopyInto(&out.VectorEngine)
}

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *ValdOperatorReleaseStatus) DeepCopyInto(out *ValdOperatorReleaseStatus) {
	*out = *in
	out.Conditions = resource.CopySliceInto(in.Conditions)
}

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *Vald) DeepCopyInto(out *Vald) {
	*out = *in
	in.Agent.DeepCopyInto(&out.Agent)
	in.Gateway.DeepCopyInto(&out.Gateway)
	in.Overlay.DeepCopyInto(&out.Overlay)
}

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *VectorEngine) DeepCopyInto(out *VectorEngine) {
	*out = *in
	in.Vald.DeepCopyInto(&out.Vald)
}
