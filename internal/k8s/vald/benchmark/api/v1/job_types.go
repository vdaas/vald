//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type BenchmarkJobSpec struct {
	Target                  *BenchmarkTarget           `json:"target,omitempty"                     yaml:"target"`
	Dataset                 *BenchmarkDataset          `json:"dataset,omitempty"                    yaml:"dataset"`
	Dimension               int                        `json:"dimension,omitempty"                  yaml:"dimension"`
	Replica                 int                        `json:"replica,omitempty"                    yaml:"replica"`
	Repetition              int                        `json:"repetition,omitempty"                 yaml:"repetition"`
	JobType                 string                     `json:"job_type,omitempty"                   yaml:"job_type"`
	InsertConfig            *config.InsertConfig       `json:"insert_config,omitempty"              yaml:"insert_config"`
	UpdateConfig            *config.UpdateConfig       `json:"update_config,omitempty"              yaml:"update_config"`
	UpsertConfig            *config.UpsertConfig       `json:"upsert_config,omitempty"              yaml:"upsert_config"`
	SearchConfig            *config.SearchConfig       `json:"search_config,omitempty"              yaml:"search_config"`
	RemoveConfig            *config.RemoveConfig       `json:"remove_config,omitempty"              yaml:"remove_config"`
	ObjectConfig            *config.ObjectConfig       `json:"object_config,omitempty"              yaml:"object_config"`
	ClientConfig            *config.GRPCClient         `json:"client_config,omitempty"              yaml:"client_config"`
	Rules                   []*config.BenchmarkJobRule `json:"rules,omitempty"                      yaml:"rules"`
	RPS                     int                        `json:"rps,omitempty"                        yaml:"rps"`
	ConcurrencyLimit        int                        `json:"concurrency_limit,omitempty"          yaml:"concurrency_limit"`
	TTLSecondsAfterFinished int                        `json:"ttl_seconds_after_finished,omitempty" yaml:"ttl_seconds_after_finished"`
	GlobalConfig            *config.GlobalConfig       `json:"global_config,omitempty"              yaml:"global_config"`
	ServerConfig            *config.Servers            `json:"server_config,omitempty"              yaml:"server_config"`
}

type BenchmarkJobStatus string

const (
	BenchmarkJobNotReady  = BenchmarkJobStatus("NotReady")
	BenchmarkJobCompleted = BenchmarkJobStatus("Completed")
	BenchmarkJobAvailable = BenchmarkJobStatus("Available")
	BenchmarkJobHealthy   = BenchmarkJobStatus("Healthy")
)

// BenchmarkTarget defines the desired state of BenchmarkTarget
type BenchmarkTarget config.BenchmarkTarget

// BenchmarkDataset defines the desired state of BenchmarkDateset
type BenchmarkDataset config.BenchmarkDataset

// BenchmarkDatasetRange defines the desired state of BenchmarkDatesetRange
type BenchmarkDatasetRange config.BenchmarkDatasetRange

// BenchmarkJobRule defines the desired state of BenchmarkJobRule
type BenchmarkJobRule config.BenchmarkJobRule

type ValdBenchmarkJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   BenchmarkJobSpec   `json:"spec,omitempty"`
	Status BenchmarkJobStatus `json:"status,omitempty"`
}

type ValdBenchmarkJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []ValdBenchmarkJob `json:"items,omitempty"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkDataset) DeepCopyInto(out *BenchmarkDataset) {
	*out = *in
	if in.Range != nil {
		in, out := &in.Range, &out.Range
		*out = new(config.BenchmarkDatasetRange)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkDataset.
func (in *BenchmarkDataset) DeepCopy() *BenchmarkDataset {
	if in == nil {
		return nil
	}
	out := new(BenchmarkDataset)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkDatasetRange) DeepCopyInto(out *BenchmarkDatasetRange) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkDatasetRange.
func (in *BenchmarkDatasetRange) DeepCopy() *BenchmarkDatasetRange {
	if in == nil {
		return nil
	}
	out := new(BenchmarkDatasetRange)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkJobRule) DeepCopyInto(out *BenchmarkJobRule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkJobRule.
func (in *BenchmarkJobRule) DeepCopy() *BenchmarkJobRule {
	if in == nil {
		return nil
	}
	out := new(BenchmarkJobRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkJobSpec) DeepCopyInto(out *BenchmarkJobSpec) {
	*out = *in
	if in.Target != nil {
		in, out := &in.Target, &out.Target
		*out = new(BenchmarkTarget)
		**out = **in
	}
	if in.Dataset != nil {
		in, out := &in.Dataset, &out.Dataset
		*out = new(BenchmarkDataset)
		(*in).DeepCopyInto(*out)
	}
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]*config.BenchmarkJobRule, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(config.BenchmarkJobRule)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkJobSpec.
func (in *BenchmarkJobSpec) DeepCopy() *BenchmarkJobSpec {
	if in == nil {
		return nil
	}
	out := new(BenchmarkJobSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkTarget) DeepCopyInto(out *BenchmarkTarget) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkTarget.
func (in *BenchmarkTarget) DeepCopy() *BenchmarkTarget {
	if in == nil {
		return nil
	}
	out := new(BenchmarkTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValdBenchmarkJob) DeepCopyInto(out *ValdBenchmarkJob) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkOperator.
func (in *ValdBenchmarkJob) DeepCopy() *ValdBenchmarkJob {
	if in == nil {
		return nil
	}
	out := new(ValdBenchmarkJob)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ValdBenchmarkJob) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValdBenchmarkJobList) DeepCopyInto(out *ValdBenchmarkJobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ValdBenchmarkJob, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkOperatorList.
func (in *ValdBenchmarkJobList) DeepCopy() *ValdBenchmarkJobList {
	if in == nil {
		return nil
	}
	out := new(ValdBenchmarkJobList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ValdBenchmarkJobList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
