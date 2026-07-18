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
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/k8s/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BenchmarkJobSpec struct {
	ObjectConfig            *config.ObjectConfig        `json:"object_config,omitempty"              yaml:"object_config"`
	ClientConfig            *config.GRPCClient          `json:"client_config,omitempty"              yaml:"client_config"`
	Target                  *BenchmarkTarget            `json:"target,omitempty"                     yaml:"target"`
	Dataset                 *BenchmarkDataset           `json:"dataset,omitempty"                    yaml:"dataset"`
	UpdateConfig            *config.UpdateConfig        `json:"update_config,omitempty"              yaml:"update_config"`
	*config.GlobalConfig    `json:",omitempty" yaml:""` //nolint:tagliatelle,tagalign // wire format fixed by embedded config.GlobalConfig's own contract, not renameable; alignment can't include a trailing nolint comment
	RemoveConfig            *config.RemoveConfig        `json:"remove_config,omitempty"              yaml:"remove_config"`
	InsertConfig            *config.InsertConfig        `json:"insert_config,omitempty"              yaml:"insert_config"`
	ServerConfig            *config.Servers             `json:"server_config,omitempty"              yaml:"server_config"`
	SearchConfig            *config.SearchConfig        `json:"search_config,omitempty"              yaml:"search_config"`
	UpsertConfig            *config.UpsertConfig        `json:"upsert_config,omitempty"              yaml:"upsert_config"`
	JobType                 string                      `json:"job_type,omitempty"                   yaml:"job_type"`
	Rules                   []*config.BenchmarkJobRule  `json:"rules,omitempty"                      yaml:"rules"`
	Repetition              int                         `json:"repetition,omitempty"                 yaml:"repetition"`
	Replica                 int                         `json:"replica,omitempty"                    yaml:"replica"`
	RPS                     int                         `json:"rps,omitempty"                        yaml:"rps"`
	ConcurrencyLimit        int                         `json:"concurrency_limit,omitempty"          yaml:"concurrency_limit"`
	TTLSecondsAfterFinished int                         `json:"ttl_seconds_after_finished,omitempty" yaml:"ttl_seconds_after_finished"`
}

type BenchmarkJobStatus string

const (
	BenchmarkJobNotReady  = BenchmarkJobStatus("NotReady")
	BenchmarkJobCompleted = BenchmarkJobStatus("Completed")
	BenchmarkJobAvailable = BenchmarkJobStatus("Available")
	BenchmarkJobHealthy   = BenchmarkJobStatus("Healthy")
)

// BenchmarkTarget defines the desired state of BenchmarkTarget.
type BenchmarkTarget config.BenchmarkTarget

// BenchmarkDataset defines the desired state of BenchmarkDateset.
type BenchmarkDataset config.BenchmarkDataset

// BenchmarkDatasetRange defines the desired state of BenchmarkDatesetRange.
type BenchmarkDatasetRange config.BenchmarkDatasetRange

// BenchmarkJobRule defines the desired state of BenchmarkJobRule.
type BenchmarkJobRule config.BenchmarkJobRule

type ValdBenchmarkJob struct { //nolint:tagliatelle // generic embed field name confuses the linter's tag-name check
	resource.Base[ValdBenchmarkJob, *ValdBenchmarkJob] `json:"-"`

	metav1.TypeMeta   `                   json:",inline"`  //nolint:tagalign // aligned with the ObjectMeta tag below, whose trailing nolint comment tagalign can't account for
	Status            BenchmarkJobStatus                   `json:"status,omitempty"`
	metav1.ObjectMeta `                   json:"metadata"` //nolint:tagliatelle,tagalign // fixed by the Kubernetes object API wire format, not renameable; alignment can't include a trailing nolint comment
	Spec              BenchmarkJobSpec                     `json:"spec"`
}

// ValdBenchmarkJobList contains a list of ValdBenchmarkJob. The whole list
// kind is derived generically from the item type.
type ValdBenchmarkJobList = resource.List[ValdBenchmarkJob, *ValdBenchmarkJob]

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *BenchmarkDataset) DeepCopyInto(out *BenchmarkDataset) {
	*out = *in
	out.Range = resource.CopyPtr(in.Range)
}

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *BenchmarkJobSpec) DeepCopyInto(out *BenchmarkJobSpec) {
	*out = *in
	out.Target = resource.CopyPtr(in.Target)
	out.Dataset = resource.CopyPtrInto(in.Dataset)
	out.Rules = resource.CopySliceFunc(in.Rules, resource.CopyPtr[config.BenchmarkJobRule])
}

// DeepCopy returns a deep copy of the receiver.
func (in *BenchmarkJobSpec) DeepCopy() *BenchmarkJobSpec {
	return resource.CopyPtrInto(in)
}

// DeepCopyInto copies the receiver into out. in must be non-nil.
func (in *ValdBenchmarkJob) DeepCopyInto(out *ValdBenchmarkJob) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy returns a deep copy of the receiver. It disambiguates the promoted
// Base.DeepCopy from the equally promoted metav1.ObjectMeta.DeepCopy.
func (in *ValdBenchmarkJob) DeepCopy() *ValdBenchmarkJob {
	return in.Base.DeepCopy()
}
