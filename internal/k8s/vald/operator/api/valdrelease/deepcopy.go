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

func (in *Spec) DeepCopyInto(out *Spec) {
	*out = *in
	in.Defaults.DeepCopyInto(&out.Defaults)
	in.Gateway.DeepCopyInto(&out.Gateway)
	in.Agent.DeepCopyInto(&out.Agent)
	in.Manager.DeepCopyInto(&out.Manager)
	in.Discoverer.DeepCopyInto(&out.Discoverer)
}

func (in *ValdRelease) DeepCopyInto(out *ValdRelease) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy returns a deep copy of the receiver. It disambiguates the promoted
// Base.DeepCopy from the equally promoted metav1.ObjectMeta.DeepCopy.
func (in *ValdRelease) DeepCopy() *ValdRelease {
	return in.Base.DeepCopy()
}

func (in *VrsStatus) DeepCopyInto(out *VrsStatus) {
	*out = *in
	in.Condition.DeepCopyInto(&out.Condition)
}
