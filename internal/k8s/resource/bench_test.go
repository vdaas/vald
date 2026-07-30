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

package resource

import (
	"context"
	"strconv"
	"testing"

	"github.com/vdaas/vald/internal/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newBenchScheme(b *testing.B) *runtime.Scheme {
	b.Helper()
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		b.Fatal(err)
	}
	return scheme
}

func newBenchClient(b *testing.B, objs ...k8s.Object) k8s.Client {
	b.Helper()
	return fake.NewClientBuilder().WithScheme(newBenchScheme(b)).WithObjects(objs...).Build()
}

func BenchmarkGetObject(b *testing.B) {
	c := newBenchClient(b, newConfigMap("cfg", "ns", map[string]string{testDataKey: testDataValue}))
	ctx := context.Background()
	b.ReportAllocs()
	for b.Loop() {
		obj, err := GetObject(ctx, c, "cfg", "ns", &corev1.ConfigMap{})
		if err != nil {
			b.Fatal(err)
		}
		if obj.GroupVersionKind().Empty() {
			b.Fatal("GVK not restored")
		}
	}
}

func BenchmarkListObjects(b *testing.B) {
	for _, n := range []int{10, 100} {
		b.Run("items="+strconv.Itoa(n), func(b *testing.B) {
			objs := make([]k8s.Object, n)
			for i := range n {
				objs[i] = newConfigMap("cfg-"+strconv.Itoa(i), "ns", nil)
			}
			c := newBenchClient(b, objs...)
			ctx := context.Background()
			b.ReportAllocs()
			for b.Loop() {
				list, err := ListObjects(ctx, c, &corev1.ConfigMapList{}, k8s.InNamespace("ns"))
				if err != nil {
					b.Fatal(err)
				}
				if len(list.Items) != n {
					b.Fatalf("unexpected item count %d", len(list.Items))
				}
			}
		})
	}
}

// BenchmarkRestoreObjectGVK isolates the per-fetch GVK restoration cost from
// the fake-client Get/List overhead measured above.
func BenchmarkRestoreObjectGVK(b *testing.B) {
	scheme := newBenchScheme(b)
	cm := &corev1.ConfigMap{}
	b.ReportAllocs()
	for b.Loop() {
		cm.TypeMeta = metav1.TypeMeta{}
		restoreObjectGVK(scheme, cm)
	}
	if cm.GroupVersionKind().Empty() {
		b.Fatal("GVK not restored")
	}
}

func BenchmarkRestoreObjectGVKParallel(b *testing.B) {
	scheme := newBenchScheme(b)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		cm := &corev1.ConfigMap{}
		ran := false
		for pb.Next() {
			cm.TypeMeta = metav1.TypeMeta{}
			restoreObjectGVK(scheme, cm)
			ran = true
		}
		// goroutines that received no iterations never restore anything.
		if ran && cm.GroupVersionKind().Empty() {
			b.Error("GVK not restored")
		}
	})
}

func BenchmarkRestoreListGVK(b *testing.B) {
	for _, n := range []int{10, 100} {
		b.Run("items="+strconv.Itoa(n), func(b *testing.B) {
			scheme := newBenchScheme(b)
			list := &corev1.ConfigMapList{Items: make([]corev1.ConfigMap, n)}
			b.ReportAllocs()
			for b.Loop() {
				list.TypeMeta = metav1.TypeMeta{}
				for i := range list.Items {
					list.Items[i].TypeMeta = metav1.TypeMeta{}
				}
				restoreListGVK(scheme, list)
			}
			if list.GroupVersionKind().Empty() {
				b.Fatal("list GVK not restored")
			}
		})
	}
}

// BenchmarkApplyDesiredState exercises the typed (reflect) path that Syncer
// runs once per desired object inside ctrl.CreateOrUpdate's mutate callback.
func BenchmarkApplyDesiredState(b *testing.B) {
	desired := newConfigMap(testNameAlpha, testNamespaceDefault,
		map[string]string{testDataKey: testDataValue})
	desired.SetLabels(map[string]string{defaultGenerationLabel: "7"})
	obj := desired.DeepCopy()
	b.ReportAllocs()
	for b.Loop() {
		if err := applyDesiredState(obj, desired); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkApplyDesiredStateParallel(b *testing.B) {
	desired := newConfigMap(testNameAlpha, testNamespaceDefault,
		map[string]string{testDataKey: testDataValue})
	desired.SetLabels(map[string]string{defaultGenerationLabel: "7"})
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		obj := desired.DeepCopy()
		for pb.Next() {
			if err := applyDesiredState(obj, desired); err != nil {
				b.Error(err)
				return
			}
		}
	})
}

// BenchmarkSyncKey documents the identity-key generation cost per object per
// reconcile (desired once + existing once during prune).
func BenchmarkSyncKey(b *testing.B) {
	cm := newConfigMap(testNameAlpha, testNamespaceDefault, nil)
	cm.SetGroupVersionKind(testConfigMapGVK())
	b.ReportAllocs()
	var key string
	for b.Loop() {
		key = syncKey(cm)
	}
	if key == "" {
		b.Fatal("empty key")
	}
}
