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

package resource

import (
	"context"
	"strings"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	applyconfigurationscorev1 "k8s.io/client-go/applyconfigurations/core/v1"
)

// unextractableList is a runtime.Object without an Items field, so
// apimeta.ExtractList fails on it.
type unextractableList struct{ metav1.TypeMeta }

func (l *unextractableList) DeepCopyObject() runtime.Object {
	out := *l
	return &out
}

// unextractableListClient returns a list object without an Items field so
// that item extraction fails on every List call.
type unextractableListClient struct{}

func (unextractableListClient) Get(context.Context, string, metav1.GetOptions) (*corev1.Pod, error) {
	return nil, nil
}

func (unextractableListClient) List(
	context.Context, metav1.ListOptions,
) (*unextractableList, error) {
	return &unextractableList{}, nil
}

func (unextractableListClient) Watch(context.Context, metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}

func (unextractableListClient) Create(
	context.Context, *corev1.Pod, metav1.CreateOptions,
) (*corev1.Pod, error) {
	return nil, nil
}

func (unextractableListClient) Delete(context.Context, string, metav1.DeleteOptions) error {
	return nil
}

func (unextractableListClient) Update(
	context.Context, *corev1.Pod, metav1.UpdateOptions,
) (*corev1.Pod, error) {
	return nil, nil
}

func (unextractableListClient) Apply(
	context.Context, *applyconfigurationscorev1.PodApplyConfiguration, metav1.ApplyOptions,
) (*corev1.Pod, error) {
	return nil, nil
}

func (unextractableListClient) Patch(
	context.Context, string, types.PatchType, []byte, metav1.PatchOptions, ...string,
) (*corev1.Pod, error) {
	return nil, nil
}

// TestWaitForStatus_LabelSelector_ExtractError guards the error ordering in
// the labelSelector path: an item extraction failure must surface immediately
// instead of being swallowed by the len==0 polling continue (which previously
// spun until the 5 minute timeout).
func TestWaitForStatus_LabelSelector_ExtractError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	matched, err := WaitForStatus(ctx, unextractableListClient{}, "", "app=vald", StatusAvailable)
	if matched {
		t.Error("WaitForStatus() matched = true, want false")
	}
	if err == nil {
		t.Fatal("WaitForStatus() error = nil, want extract failure")
	}
	if !strings.Contains(err.Error(), "failed to extract items") {
		t.Errorf("WaitForStatus() error = %v, want extract failure", err)
	}
}

func TestCheckResourceState_Details(t *testing.T) {
	t.Parallel()

	type test struct {
		name       string
		obj        Object
		wantStatus ResourceStatus
		wantInfo   []string
	}

	tests := []test{
		{
			name: "paused deployment reports name and reason",
			obj: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: "vald-lb-gateway"},
				Spec:       appsv1.DeploymentSpec{Paused: true},
			},
			wantStatus: StatusPaused,
			wantInfo:   []string{"vald-lb-gateway", "Deployment is paused."},
		},
		{
			name: "available deployment reports success",
			obj: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: "vald-agent"},
				Status: appsv1.DeploymentStatus{
					Replicas:          1,
					UpdatedReplicas:   1,
					AvailableReplicas: 1,
				},
			},
			wantStatus: StatusAvailable,
			wantInfo:   []string{"Deployment is fully operational."},
		},
		{
			name: "degraded deployment reports replica shortage",
			obj: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: "vald-discoverer"},
				Status: appsv1.DeploymentStatus{
					Replicas:        1,
					UpdatedReplicas: 1,
				},
			},
			wantStatus: StatusDegraded,
			wantInfo:   []string{"vald-discoverer", "Only 0 out of 1 replicas available."},
		},
		{
			name: "pending pod reports phase",
			obj: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "vald-agent-0"},
				Status:     corev1.PodStatus{Phase: corev1.PodPending},
			},
			wantStatus: StatusPending,
			wantInfo:   []string{"vald-agent-0", "Pod is pending scheduling or initialization."},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			status, info, err := checkResourceState(tc.obj)
			if err != nil {
				t.Fatalf("CheckResourceState() error = %v", err)
			}
			if status != tc.wantStatus {
				t.Errorf("CheckResourceState() status = %v, want %v", status, tc.wantStatus)
			}
			for _, want := range tc.wantInfo {
				if !strings.Contains(info, want) {
					t.Errorf("CheckResourceState() info = %q, want it to contain %q", info, want)
				}
			}
		})
	}
}
