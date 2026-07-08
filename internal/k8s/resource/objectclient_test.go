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
	"testing"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newPodObjectClient(
	t *testing.T, objs ...Object,
) *ObjectClient[corev1.ConfigMap, *corev1.ConfigMap] {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add corev1 scheme: %v", err)
	}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
	return NewObjectClient[corev1.ConfigMap](c)
}

func newConfigMap(name, namespace string, data map[string]string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Data:       data,
	}
}

func TestObjectClient_Get(t *testing.T) {
	t.Parallel()

	type test struct {
		name      string
		objs      []Object
		getName   string
		wantErr   bool
		wantFound bool
	}

	tests := []test{
		{
			name:      "returns the object when it exists",
			objs:      []Object{newConfigMap("exists", "default", map[string]string{"k": "v"})},
			getName:   "exists",
			wantFound: true,
		},
		{
			name:    "returns NotFound when the object is missing",
			getName: "missing",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			oc := newPodObjectClient(t, tc.objs...)
			obj, err := oc.Get(context.Background(), tc.getName, "default")
			if tc.wantErr {
				if err == nil {
					t.Fatal("Get() error = nil, want error")
				}
				if !apierrors.IsNotFound(err) {
					t.Errorf("Get() error = %v, want NotFound", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Get() error = %v, want nil", err)
			}
			if tc.wantFound && obj.GetName() != tc.getName {
				t.Errorf("Get() name = %q, want %q", obj.GetName(), tc.getName)
			}
		})
	}
}

func TestObjectClient_CreateOrUpdate(t *testing.T) {
	t.Parallel()

	oc := newPodObjectClient(t)
	ctx := context.Background()

	obj := newConfigMap("target", "default", map[string]string{"k": "v1"})
	ope, err := oc.CreateOrUpdate(ctx, obj, func() error { return nil })
	if err != nil {
		t.Fatalf("CreateOrUpdate(create) error = %v", err)
	}
	if ope != OperationResultCreated {
		t.Errorf("CreateOrUpdate(create) = %v, want %v", ope, OperationResultCreated)
	}

	// mutate on the second pass: the object exists, so it must be updated
	ope, err = oc.CreateOrUpdate(ctx, obj, func() error {
		obj.Data["k"] = "v2"
		return nil
	})
	if err != nil {
		t.Fatalf("CreateOrUpdate(update) error = %v", err)
	}
	if ope != OperationResultUpdated {
		t.Errorf("CreateOrUpdate(update) = %v, want %v", ope, OperationResultUpdated)
	}

	got, err := oc.Get(ctx, "target", "default")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Data["k"] != "v2" {
		t.Errorf("data = %q, want %q", got.Data["k"], "v2")
	}
}

func TestObjectClient_Wait(t *testing.T) {
	t.Parallel()

	t.Run("returns the context error when canceled", func(t *testing.T) {
		t.Parallel()

		oc := newPodObjectClient(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		done, err := oc.Wait(ctx, "missing", "default", func(*corev1.ConfigMap) (bool, error) {
			return true, nil
		})
		if done || err == nil {
			t.Errorf("Wait() = (%v, %v), want (false, context error)", done, err)
		}
	})

	t.Run("returns true once eval reports done", func(t *testing.T) {
		t.Parallel()

		oc := newPodObjectClient(t, newConfigMap("ready", "default", nil))
		done, err := oc.Wait(
			context.Background(), "ready", "default",
			func(cm *corev1.ConfigMap) (bool, error) {
				return cm.GetName() == "ready", nil
			},
		)
		if err != nil {
			t.Fatalf("Wait() error = %v", err)
		}
		if !done {
			t.Error("Wait() done = false, want true")
		}
	})
}

func TestWatchObjects(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add corev1 scheme: %v", err)
	}
	c := fake.NewClientBuilder().WithScheme(scheme).Build()

	w, err := WatchObjects(context.Background(), c, new(corev1.ConfigMapList))
	if err != nil {
		t.Fatalf("WatchObjects() error = %v", err)
	}
	defer w.Stop()

	if err := c.Create(context.Background(), newConfigMap("watched", "default", nil)); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	ev := <-w.ResultChan()
	if ev.Type != watch.Added {
		t.Errorf("event type = %v, want %v", ev.Type, watch.Added)
	}
	cm, ok := ev.Object.(*corev1.ConfigMap)
	if !ok || cm.GetName() != "watched" {
		t.Errorf("event object = %v, want ConfigMap watched", ev.Object)
	}
}

func TestLabelSelector(t *testing.T) {
	t.Parallel()

	type test struct {
		name    string
		key     string
		op      selection.Operator
		vals    []string
		want    string
		wantErr bool
	}

	tests := []test{
		{
			name: "builds an equality selector",
			key:  "app",
			op:   selection.Equals,
			vals: []string{"vald-agent"},
			want: "app=vald-agent",
		},
		{
			name: "builds an exists selector",
			key:  "app",
			op:   selection.Exists,
			want: "app",
		},
		{
			name:    "returns error for an invalid key",
			key:     "in valid key!",
			op:      selection.Equals,
			vals:    []string{"v"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			sel, err := LabelSelector(tc.key, tc.op, tc.vals)
			if tc.wantErr {
				if err == nil {
					t.Error("LabelSelector() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("LabelSelector() error = %v", err)
			}
			if sel.String() != tc.want {
				t.Errorf("LabelSelector() = %q, want %q", sel.String(), tc.want)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
