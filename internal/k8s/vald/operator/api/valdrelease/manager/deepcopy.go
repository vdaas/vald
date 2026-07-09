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

package manager

import (
	"maps"

	"github.com/vdaas/vald/internal/k8s/resource"
)

func (in *Creator) DeepCopyInto(out *Creator) {
	*out = *in
	out.Concurrency = resource.CopyPtr(in.Concurrency)
	out.Affinity = resource.CopyPtrInto(in.Affinity)
	out.NodeSelector = maps.Clone(in.NodeSelector)
	out.Tolerations = resource.CopyPtrSliceInto(in.Tolerations)
	out.Image = resource.CopyPtr(in.Image)
}

func (in *Index) DeepCopyInto(out *Index) {
	*out = *in
	out.Logging = resource.CopyPtr(in.Logging)
	in.Indexer.DeepCopyInto(&out.Indexer)
	out.Saver = resource.CopyPtrInto(in.Saver)
	out.Creator = resource.CopyPtrInto(in.Creator)
	out.Resources = resource.CopyPtrInto(in.Resources)
	out.TopologySpreadConstraints = resource.CopySliceInto(in.TopologySpreadConstraints)
	out.Affinity = resource.CopyPtrInto(in.Affinity)
	out.NodeSelector = maps.Clone(in.NodeSelector)
	out.Tolerations = resource.CopyPtrSliceInto(in.Tolerations)
}

func (in *Indexer) DeepCopyInto(out *Indexer) {
	*out = *in
	out.AutoIndexLength = resource.CopyPtr(in.AutoIndexLength)
	out.Concurrency = resource.CopyPtr(in.Concurrency)
}

func (in *Manager) DeepCopyInto(out *Manager) {
	*out = *in
	in.Index.DeepCopyInto(&out.Index)
}

func (in *Saver) DeepCopyInto(out *Saver) {
	*out = *in
	out.Affinity = resource.CopyPtrInto(in.Affinity)
	out.NodeSelector = maps.Clone(in.NodeSelector)
	out.Tolerations = resource.CopyPtrSliceInto(in.Tolerations)
	out.Image = resource.CopyPtr(in.Image)
}
