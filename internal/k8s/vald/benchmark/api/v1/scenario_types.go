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

import (
	"github.com/vdaas/vald/internal/k8s/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ValdBenchmarkScenarioSpec struct {
	Target  *BenchmarkTarget    `json:"target,omitempty"`
	Dataset *BenchmarkDataset   `json:"dataset,omitempty"`
	Jobs    []*BenchmarkJobSpec `json:"jobs,omitempty"`
}

type ValdBenchmarkScenarioStatus string

const (
	BenchmarkScenarioNotReady  ValdBenchmarkScenarioStatus = "NotReady"
	BenchmarkScenarioCompleted ValdBenchmarkScenarioStatus = "Completed"
	BenchmarkScenarioAvailable ValdBenchmarkScenarioStatus = "Available"
	BenchmarkScenarioHealthy   ValdBenchmarkScenarioStatus = "Healthy"
)

type ValdBenchmarkScenario struct {
	resource.Base[ValdBenchmarkScenario, *ValdBenchmarkScenario] `json:"-"`

	metav1.TypeMeta   `                            json:",inline"`
	Status            ValdBenchmarkScenarioStatus `json:"status,omitempty"`
	metav1.ObjectMeta `                            json:"metadata"`
	Spec              ValdBenchmarkScenarioSpec `json:"spec"`
}

// ValdBenchmarkScenarioList contains a list of ValdBenchmarkScenario. The
// whole list kind is derived generically from the item type.
type ValdBenchmarkScenarioList = resource.List[ValdBenchmarkScenario, *ValdBenchmarkScenario]

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *ValdBenchmarkScenario) DeepCopyInto(out *ValdBenchmarkScenario) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy returns a deep copy of the receiver. It disambiguates the promoted
// Base.DeepCopy from the equally promoted metav1.ObjectMeta.DeepCopy.
func (in *ValdBenchmarkScenario) DeepCopy() *ValdBenchmarkScenario {
	return in.Base.DeepCopy()
}

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *ValdBenchmarkScenarioSpec) DeepCopyInto(out *ValdBenchmarkScenarioSpec) {
	*out = *in
	out.Target = resource.CopyPtr(in.Target)
	out.Dataset = resource.CopyPtrInto(in.Dataset)
	out.Jobs = resource.CopySliceFunc(in.Jobs, resource.CopyPtrInto[BenchmarkJobSpec])
}
