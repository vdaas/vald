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

package gateway

import (
	"maps"

	"github.com/vdaas/vald/internal/k8s/resource"
)

func (in *Gateway) DeepCopyInto(out *Gateway) {
	*out = *in
	in.Lb.DeepCopyInto(&out.Lb)
}

func (in *Ingress) DeepCopyInto(out *Ingress) {
	*out = *in
	out.Annotations = maps.Clone(in.Annotations)
}

func (in *Lb) DeepCopyInto(out *Lb) {
	*out = *in
	out.Logging = resource.CopyPtr(in.Logging)
	out.Affinity = resource.CopyPtrInto(in.Affinity)
	out.NodeSelector = maps.Clone(in.NodeSelector)
	out.Tolerations = resource.CopyPtrSliceInto(in.Tolerations)
	in.Ingress.DeepCopyInto(&out.Ingress)
	out.Hpa = resource.CopyPtr(in.Hpa)
	out.Resources = resource.CopyPtrInto(in.Resources)
	out.TopologySpreadConstraints = resource.CopySliceInto(in.TopologySpreadConstraints)
}
