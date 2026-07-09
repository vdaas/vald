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

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/internal/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newSchemeForSyncerTests(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	assert.NoError(t, corev1.AddToScheme(scheme))
	assert.NoError(t, appsv1.AddToScheme(scheme))
	return scheme
}

func newSyncOwner() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "owner",
			Namespace:  "default",
			UID:        types.UID("owner-uid"),
			Generation: 7,
		},
	}
}

// listConfigMaps returns every ConfigMap in the namespace as the prune
// candidate set, mirroring how callers list their owned resource kind.
func listConfigMaps(c ObjectAPI, ns string) func(ctx context.Context) ([]Object, error) {
	return func(ctx context.Context) ([]Object, error) {
		list := &corev1.ConfigMapList{}
		if err := c.List(ctx, list, kclient.InNamespace(ns)); err != nil {
			return nil, err
		}
		out := make([]Object, len(list.Items))
		for i := range list.Items {
			list.Items[i].SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
			out[i] = &list.Items[i]
		}
		return out, nil
	}
}

func TestSyncKey_IncludesGVK(t *testing.T) {
	cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "n", Namespace: "ns"}}
	cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	got := syncKey(cm)
	assert.Equal(t, "/v1/ConfigMap/ns/n", got)
}

func TestSyncer_Sync_EmptyDesired(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	s := NewSyncer(c, scheme, "")

	res, err := s.Sync(context.Background(), newSyncOwner(), nil, listConfigMaps(c, "default"))
	assert.NoError(t, err)
	assert.Contains(t, res, "no_resources")
}

func TestSyncer_Sync_ExistingError(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	s := NewSyncer(c, scheme, "")

	cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "alpha", Namespace: "default"}}
	cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	_, err := s.Sync(context.Background(), newSyncOwner(), []Object{cm},
		func(context.Context) ([]Object, error) { return nil, errors.New("boom") })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

func TestSyncer_Sync_CreatesObjectsWithOwnerAndLabels(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	s := NewSyncer(c, scheme, "")
	owner := newSyncOwner()

	cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "alpha", Namespace: "default"}}
	cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	res, err := s.Sync(context.Background(), owner, []Object{cm}, listConfigMaps(c, "default"))
	assert.NoError(t, err)
	assert.Contains(t, res, "/v1/ConfigMap/default/alpha")

	got := &corev1.ConfigMap{}
	err = c.Get(context.Background(), kclient.ObjectKey{Name: "alpha", Namespace: "default"}, got)
	assert.NoError(t, err)
	assert.Equal(t, "7", got.GetLabels()["managed-generation"])
	assert.Len(t, got.GetOwnerReferences(), 1)
	assert.Equal(t, types.UID("owner-uid"), got.GetOwnerReferences()[0].UID)
}

func TestSyncer_Sync_PrunesOrphans(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	owner := newSyncOwner()

	orphan := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orphan",
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         appsv1.SchemeGroupVersion.String(),
					Kind:               "Deployment",
					Name:               owner.Name,
					UID:                owner.UID,
					Controller:         new(true),
					BlockOwnerDeletion: new(true),
				},
			},
		},
	}
	orphan.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(orphan).Build()
	s := NewSyncer(c, scheme, "")

	keep := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "keep", Namespace: "default"}}
	keep.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	res, err := s.Sync(context.Background(), owner, []Object{keep}, listConfigMaps(c, "default"))
	assert.NoError(t, err)
	assert.Contains(t, res, "/v1/ConfigMap/default/keep")
	assert.Equal(t, SyncResults{
		"/v1/ConfigMap/default/keep":   res["/v1/ConfigMap/default/keep"],
		"/v1/ConfigMap/default/orphan": prunedResult,
	}, res)

	err = c.Get(context.Background(), kclient.ObjectKey{Name: "orphan", Namespace: "default"}, &corev1.ConfigMap{})
	assert.Error(t, err, "orphan should have been deleted")
}

func TestSyncer_Sync_NilExistingSkipsPrune(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	owner := newSyncOwner()

	orphan := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orphan",
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         appsv1.SchemeGroupVersion.String(),
					Kind:               "Deployment",
					Name:               owner.Name,
					UID:                owner.UID,
					Controller:         new(true),
					BlockOwnerDeletion: new(true),
				},
			},
		},
	}
	orphan.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(orphan).Build()
	s := NewSyncer(c, scheme, "")

	keep := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "keep", Namespace: "default"}}
	keep.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	_, err := s.Sync(context.Background(), owner, []Object{keep}, nil)
	assert.NoError(t, err)

	err = c.Get(context.Background(), kclient.ObjectKey{Name: "orphan", Namespace: "default"}, &corev1.ConfigMap{})
	assert.NoError(t, err, "orphan must survive when pruning is skipped")
}

func TestSyncer_Sync_UpdatesExistingObject(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	owner := newSyncOwner()

	existing := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "alpha",
			Namespace: "default",
			Labels:    map[string]string{"managed-generation": "3"},
		},
		Data: map[string]string{"key": "stale"},
	}
	existing.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()
	s := NewSyncer(c, scheme, "")

	desired := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "alpha", Namespace: "default"},
		Data:       map[string]string{"key": "fresh"},
	}
	desired.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	res, err := s.Sync(context.Background(), owner, []Object{desired}, listConfigMaps(c, "default"))
	assert.NoError(t, err)
	assert.Equal(t, OperationResultUpdated, res["/v1/ConfigMap/default/alpha"],
		"an existing object whose state differs from desired must be updated")

	got := &corev1.ConfigMap{}
	assert.NoError(t, c.Get(context.Background(), kclient.ObjectKey{Name: "alpha", Namespace: "default"}, got))
	assert.Equal(t, "fresh", got.Data["key"], "desired data must converge onto the live object")
	assert.Equal(t, "7", got.GetLabels()["managed-generation"], "generation label must follow the owner generation")
	assert.Len(t, got.GetOwnerReferences(), 1)
	assert.Equal(t, types.UID("owner-uid"), got.GetOwnerReferences()[0].UID)
}

func TestSyncer_Sync_UpdatesExistingUnstructured(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	owner := newSyncOwner()

	newUnstructuredCM := func(data map[string]any) *unstructured.Unstructured {
		u := &unstructured.Unstructured{Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]any{
				"name":      "alpha",
				"namespace": "default",
			},
			"data": data,
		}}
		return u
	}

	existing := newUnstructuredCM(map[string]any{"key": "stale", "obsolete": "gone"})
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()
	s := NewSyncer(c, scheme, "")

	desired := newUnstructuredCM(map[string]any{"key": "fresh"})

	res, err := s.Sync(context.Background(), owner, []Object{desired}, nil)
	assert.NoError(t, err)
	assert.Equal(t, OperationResultUpdated, res["/v1/ConfigMap/default/alpha"])

	got := newUnstructuredCM(nil)
	got.Object = map[string]any{"apiVersion": "v1", "kind": "ConfigMap"}
	got.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
	assert.NoError(t, c.Get(context.Background(), kclient.ObjectKey{Name: "alpha", Namespace: "default"}, got))
	data, _, err := unstructured.NestedStringMap(got.Object, "data")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"key": "fresh"}, data,
		"desired payload must replace the live payload including stale keys")
	assert.Equal(t, "7", got.GetLabels()["managed-generation"])
	assert.Len(t, got.GetOwnerReferences(), 1)
}

func TestSyncer_Sync_EmptyDesiredPrunesOwned(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	owner := newSyncOwner()

	orphan := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orphan",
			Namespace: "default",
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         appsv1.SchemeGroupVersion.String(),
					Kind:               "Deployment",
					Name:               owner.Name,
					UID:                owner.UID,
					Controller:         new(true),
					BlockOwnerDeletion: new(true),
				},
			},
		},
	}
	orphan.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(orphan).Build()
	s := NewSyncer(c, scheme, "")

	res, err := s.Sync(context.Background(), owner, nil, listConfigMaps(c, "default"))
	assert.NoError(t, err)
	assert.Contains(t, res, "no_resources")
	assert.Equal(t, prunedResult, res["/v1/ConfigMap/default/orphan"],
		"an empty desired set must still prune every owned object")

	err = c.Get(context.Background(), kclient.ObjectKey{Name: "orphan", Namespace: "default"}, &corev1.ConfigMap{})
	assert.Error(t, err, "owned orphan must be deleted when the desired set is empty")
}
