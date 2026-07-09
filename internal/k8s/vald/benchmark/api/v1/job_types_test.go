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
	"testing"

	"github.com/vdaas/vald/internal/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newValdBenchmarkJobFixture() *ValdBenchmarkJob {
	return &ValdBenchmarkJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job-fixture",
			Namespace: "default",
			Labels:    map[string]string{"app": "vald-benchmark"},
		},
		Spec: BenchmarkJobSpec{
			JobType: "search",
			Target:  &BenchmarkTarget{Host: "vald-lb-gateway", Port: 8081},
			Dataset: &BenchmarkDataset{
				Name:  "fashion-mnist",
				Range: &config.BenchmarkDatasetRange{Start: 1, End: 1000},
			},
			Rules: []*config.BenchmarkJobRule{{Name: "rule-a"}},
		},
		Status: BenchmarkJobAvailable,
	}
}

func TestValdBenchmarkJob_DeepCopy(t *testing.T) {
	t.Parallel()

	orig := newValdBenchmarkJobFixture()
	cp := orig.DeepCopy()
	if cp == nil {
		t.Fatal("DeepCopy() = nil, want copy")
	}

	cp.Labels["app"] = "mutated"
	cp.Spec.Target.Host = "mutated"
	cp.Spec.Dataset.Range.End = 9
	cp.Spec.Rules[0].Name = "mutated"

	if orig.Labels["app"] != "vald-benchmark" {
		t.Errorf("original labels mutated: %v", orig.Labels)
	}
	if orig.Spec.Target.Host != "vald-lb-gateway" {
		t.Errorf("original target mutated: %v", orig.Spec.Target)
	}
	if orig.Spec.Dataset.Range.End != 1000 {
		t.Errorf("original dataset range mutated: %v", orig.Spec.Dataset.Range)
	}
	if orig.Spec.Rules[0].Name != "rule-a" {
		t.Errorf("original rules mutated: %v", orig.Spec.Rules[0])
	}
}

func TestValdBenchmarkJob_DeepCopyObject(t *testing.T) {
	t.Parallel()

	orig := newValdBenchmarkJobFixture()
	obj := orig.DeepCopyObject()
	cp, ok := obj.(*ValdBenchmarkJob)
	if !ok {
		t.Fatalf("DeepCopyObject() = %T, want *ValdBenchmarkJob", obj)
	}
	if cp.GetName() != orig.GetName() {
		t.Errorf("copied name = %q, want %q", cp.GetName(), orig.GetName())
	}

	list := &ValdBenchmarkJobList{Items: []ValdBenchmarkJob{*orig}}
	lobj := list.DeepCopyObject()
	lcp, ok := lobj.(*ValdBenchmarkJobList)
	if !ok {
		t.Fatalf("DeepCopyObject() = %T, want *ValdBenchmarkJobList", lobj)
	}
	lcp.Items[0].Labels["app"] = "mutated"
	if orig.Labels["app"] != "vald-benchmark" {
		t.Errorf("original mutated through list copy: %v", orig.Labels)
	}
}
