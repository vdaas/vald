//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

package operator

import (
	job "github.com/vdaas/vald/internal/k8s/benchmark/job"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type BenchmarkOperatorSpec struct {
	Target  *job.BenchmarkTarget
	Dataset *job.BenchmarkDataset
	Jobs    []*job.BenchmarkJobSpec
}

type BenchmarkOperatorStatus string

const (
	BenchmarkOperatorNotReady  = BenchmarkOperatorStatus("NotReady")
	BenchmarkOperatorAvailable = BenchmarkOperatorStatus("Available")
	BenchmarkOperatorHealthy   = BenchmarkOperatorStatus("Healthy")
)

type BenchmarkOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BenchmarkOperatorSpec   `json:"spec,omitempty"`
	Status BenchmarkOperatorStatus `json:"status,omitempty"`
}

type BenchmarkOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BenchmarkOperator `json:"items"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkOperator) DeepCopyInto(out *BenchmarkOperator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkOperator.
func (in *BenchmarkOperator) DeepCopy() *BenchmarkOperator {
	if in == nil {
		return nil
	}
	out := new(BenchmarkOperator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BenchmarkOperator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkOperatorList) DeepCopyInto(out *BenchmarkOperatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BenchmarkOperator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkOperatorList.
func (in *BenchmarkOperatorList) DeepCopy() *BenchmarkOperatorList {
	if in == nil {
		return nil
	}
	out := new(BenchmarkOperatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BenchmarkOperatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BenchmarkOperatorSpec) DeepCopyInto(out *BenchmarkOperatorSpec) {
	*out = *in
	if in.Target != nil {
		in, out := &in.Target, &out.Target
		*out = new(job.BenchmarkTarget)
		**out = **in
	}
	if in.Dataset != nil {
		in, out := &in.Dataset, &out.Dataset
		*out = new(job.BenchmarkDataset)
		(*in).DeepCopyInto(*out)
	}
	if in.Jobs != nil {
		in, out := &in.Jobs, &out.Jobs
		*out = make([]*job.BenchmarkJobSpec, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(job.BenchmarkJobSpec)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BenchmarkOperatorSpec.
func (in *BenchmarkOperatorSpec) DeepCopy() *BenchmarkOperatorSpec {
	if in == nil {
		return nil
	}
	out := new(BenchmarkOperatorSpec)
	in.DeepCopyInto(out)
	return out
}
