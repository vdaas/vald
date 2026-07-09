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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newValdMirrorTargetFixture() *ValdMirrorTarget {
	return &ValdMirrorTarget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "target-fixture",
			Namespace: "default",
			Labels:    map[string]string{"app": "vald-mirror"},
		},
		Spec: MirrorTargetSpec{
			Colocation: "dc1",
			Target:     MirrorTarget{Host: "vald-mirror-gateway", Port: 8081},
		},
		Status: MirrorTargetStatus{Phase: MirrorTargetConnected},
	}
}

func TestValdMirrorTarget_DeepCopy(t *testing.T) {
	t.Parallel()

	orig := newValdMirrorTargetFixture()
	cp := orig.DeepCopy()
	if cp == nil {
		t.Fatal("DeepCopy() = nil, want copy")
	}

	cp.Labels["app"] = "mutated"
	cp.Spec.Target.Host = "mutated"
	cp.Status.Phase = MirrorTargetDisconnected

	if orig.Labels["app"] != "vald-mirror" {
		t.Errorf("original labels mutated: %v", orig.Labels)
	}
	if orig.Spec.Target.Host != "vald-mirror-gateway" {
		t.Errorf("original spec mutated: %v", orig.Spec)
	}
	if orig.Status.Phase != MirrorTargetConnected {
		t.Errorf("original status mutated: %v", orig.Status)
	}
}

func TestValdMirrorTarget_DeepCopyObject(t *testing.T) {
	t.Parallel()

	orig := newValdMirrorTargetFixture()
	obj := orig.DeepCopyObject()
	cp, ok := obj.(*ValdMirrorTarget)
	if !ok {
		t.Fatalf("DeepCopyObject() = %T, want *ValdMirrorTarget", obj)
	}
	if cp.GetName() != orig.GetName() {
		t.Errorf("copied name = %q, want %q", cp.GetName(), orig.GetName())
	}

	list := &ValdMirrorTargetList{Items: []ValdMirrorTarget{*orig}}
	lobj := list.DeepCopyObject()
	lcp, ok := lobj.(*ValdMirrorTargetList)
	if !ok {
		t.Fatalf("DeepCopyObject() = %T, want *ValdMirrorTargetList", lobj)
	}
	lcp.Items[0].Labels["app"] = "mutated"
	if orig.Labels["app"] != "vald-mirror" {
		t.Errorf("original mutated through list copy: %v", orig.Labels)
	}
}
