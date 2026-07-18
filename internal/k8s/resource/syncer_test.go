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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/sync"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	testNamespaceDefault = "default"
	testNameAlpha        = "alpha"
	testNameOrphan       = "orphan"
	testNameKeep         = "keep"
	testFreshValue       = "fresh"
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
			Namespace:  testNamespaceDefault,
			UID:        types.UID("owner-uid"),
			Generation: 7,
		},
	}
}

// listConfigMaps returns every ConfigMap in testNamespaceDefault as the prune
// candidate set, mirroring how callers list their owned resource kind.
func listConfigMaps(c k8s.Client) func(ctx context.Context) ([]Object, error) {
	return func(ctx context.Context) ([]Object, error) {
		list := &corev1.ConfigMapList{}
		if err := c.List(ctx, list, kclient.InNamespace(testNamespaceDefault)); err != nil {
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

	res, err := s.Sync(context.Background(), newSyncOwner(), nil, listConfigMaps(c))
	assert.NoError(t, err)
	assert.Contains(t, res, "no_resources")
}

func TestSyncer_Sync_ExistingError(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	s := NewSyncer(c, scheme, "")

	cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: testNameAlpha, Namespace: testNamespaceDefault}}
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

	cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: testNameAlpha, Namespace: testNamespaceDefault}}
	cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	res, err := s.Sync(context.Background(), owner, []Object{cm}, listConfigMaps(c))
	assert.NoError(t, err)
	assert.Contains(t, res, "/v1/ConfigMap/default/alpha")

	got := &corev1.ConfigMap{}
	err = c.Get(context.Background(), kclient.ObjectKey{Name: testNameAlpha, Namespace: testNamespaceDefault}, got)
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
			Name:      testNameOrphan,
			Namespace: testNamespaceDefault,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         appsv1.SchemeGroupVersion.String(),
					Kind:               testKindDeployment,
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

	keep := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: testNameKeep, Namespace: testNamespaceDefault}}
	keep.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	res, err := s.Sync(context.Background(), owner, []Object{keep}, listConfigMaps(c))
	assert.NoError(t, err)
	assert.Contains(t, res, "/v1/ConfigMap/default/keep")
	assert.Equal(t, SyncResults{
		"/v1/ConfigMap/default/keep":   res["/v1/ConfigMap/default/keep"],
		"/v1/ConfigMap/default/orphan": prunedResult,
	}, res)

	err = c.Get(context.Background(), kclient.ObjectKey{Name: testNameOrphan, Namespace: testNamespaceDefault}, &corev1.ConfigMap{})
	assert.Error(t, err, "orphan should have been deleted")
}

func TestSyncer_Sync_NilExistingSkipsPrune(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	owner := newSyncOwner()

	orphan := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testNameOrphan,
			Namespace: testNamespaceDefault,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         appsv1.SchemeGroupVersion.String(),
					Kind:               testKindDeployment,
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

	keep := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: testNameKeep, Namespace: testNamespaceDefault}}
	keep.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	_, err := s.Sync(context.Background(), owner, []Object{keep}, nil)
	assert.NoError(t, err)

	err = c.Get(context.Background(), kclient.ObjectKey{Name: testNameOrphan, Namespace: testNamespaceDefault}, &corev1.ConfigMap{})
	assert.NoError(t, err, "orphan must survive when pruning is skipped")
}

func TestSyncer_Sync_UpdatesExistingObject(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	owner := newSyncOwner()

	existing := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testNameAlpha,
			Namespace: testNamespaceDefault,
			Labels:    map[string]string{"managed-generation": "3"},
		},
		Data: map[string]string{testDataKey: "stale"},
	}
	existing.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()
	s := NewSyncer(c, scheme, "")

	desired := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: testNameAlpha, Namespace: testNamespaceDefault},
		Data:       map[string]string{testDataKey: testFreshValue},
	}
	desired.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

	res, err := s.Sync(context.Background(), owner, []Object{desired}, listConfigMaps(c))
	assert.NoError(t, err)
	assert.Equal(t, OperationResultUpdated, res["/v1/ConfigMap/default/alpha"],
		"an existing object whose state differs from desired must be updated")

	got := &corev1.ConfigMap{}
	assert.NoError(t, c.Get(context.Background(), kclient.ObjectKey{Name: testNameAlpha, Namespace: testNamespaceDefault}, got))
	assert.Equal(t, testFreshValue, got.Data[testDataKey], "desired data must converge onto the live object")
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
				"name":      testNameAlpha,
				"namespace": testNamespaceDefault,
			},
			"data": data,
		}}
		return u
	}

	existing := newUnstructuredCM(map[string]any{testDataKey: "stale", "obsolete": "gone"})
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()
	s := NewSyncer(c, scheme, "")

	desired := newUnstructuredCM(map[string]any{testDataKey: testFreshValue})

	res, err := s.Sync(context.Background(), owner, []Object{desired}, nil)
	assert.NoError(t, err)
	assert.Equal(t, OperationResultUpdated, res["/v1/ConfigMap/default/alpha"])

	got := newUnstructuredCM(nil)
	got.Object = map[string]any{"apiVersion": "v1", "kind": "ConfigMap"}
	got.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
	assert.NoError(t, c.Get(context.Background(), kclient.ObjectKey{Name: testNameAlpha, Namespace: testNamespaceDefault}, got))
	data, _, err := unstructured.NestedStringMap(got.Object, "data")
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{testDataKey: testFreshValue}, data,
		"desired payload must replace the live payload including stale keys")
	assert.Equal(t, "7", got.GetLabels()["managed-generation"])
	assert.Len(t, got.GetOwnerReferences(), 1)
}

// TestPayloadFieldIndexes pins the payload-field cache contract: repeated
// lookups return the identical cached plan, and the plan contains exactly the
// exported non-TypeMeta/ObjectMeta/Status fields.
func TestPayloadFieldIndexes(t *testing.T) {
	t.Parallel()

	typ := reflect.TypeFor[corev1.ConfigMap]()
	idx := payloadFieldIndexes(typ)
	assert.Equal(t, idx, payloadFieldIndexes(typ),
		"repeated lookups must return the same cached plan")

	want := make([]int, 0, typ.NumField())
	names := make([]string, 0, typ.NumField())
	for i := range typ.NumField() {
		f := typ.Field(i)
		if f.IsExported() && f.Name != "TypeMeta" && f.Name != "ObjectMeta" && f.Name != "Status" {
			want = append(want, i)
			names = append(names, f.Name)
		}
	}
	assert.Equal(t, want, idx)
	assert.Contains(t, names, "Data", "ConfigMap payload must include Data")
}

// TestApplyDesiredState_ConcurrentTypes hammers the typed reflect path from
// many goroutines across two distinct struct types, so the payload-field
// cache's Load/Store paths run concurrently (validated under -race) while
// every worker's copy result stays correct.
func TestApplyDesiredState_ConcurrentTypes(t *testing.T) {
	t.Parallel()

	const (
		workers    = 8
		iterations = 100
	)

	var wg sync.WaitGroup
	for range workers {
		wg.Add(2)
		go func() {
			defer wg.Done()
			desired := newConfigMap(testNameAlpha, testNamespaceDefault,
				map[string]string{testDataKey: testFreshValue})
			obj := newConfigMap(testNameAlpha, testNamespaceDefault,
				map[string]string{testDataKey: "outdated"})
			for range iterations {
				if !assert.NoError(t, applyDesiredState(obj, desired)) {
					return
				}
			}
			assert.Equal(t, testFreshValue, obj.Data[testDataKey])
		}()
		go func() {
			defer wg.Done()
			desired := &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{Replicas: new(int32(3))},
			}
			obj := &appsv1.Deployment{}
			for range iterations {
				if !assert.NoError(t, applyDesiredState(obj, desired)) {
					return
				}
			}
			if assert.NotNil(t, obj.Spec.Replicas) {
				assert.Equal(t, int32(3), *obj.Spec.Replicas)
			}
		}()
	}
	wg.Wait()
}

// TestSyncer_Prune_EmptyGVKFailSafe verifies that prune never deletes an
// owned object whose GVK is empty (unknown identity, e.g. the caller's GVK
// restoration failed): its syncKey degenerates to "///ns/name", can never
// match a desired entry, and without the fail-safe the object would be
// wrongly deleted. The warn-and-skip branch is the only code path that keeps
// such an object alive, so the survival assertions below exercise exactly
// that (log.Warnf) path. The populated-GVK case pins the normal prune
// behavior unchanged.
func TestSyncer_Prune_EmptyGVKFailSafe(t *testing.T) {
	tests := []struct {
		name       string
		stampGVK   bool
		wantPruned bool
	}{
		{
			name:       "empty GVK: owned orphan is skipped instead of pruned",
			stampGVK:   false,
			wantPruned: false,
		},
		{
			name:       "populated GVK: owned orphan is pruned as before",
			stampGVK:   true,
			wantPruned: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := newSchemeForSyncerTests(t)
			owner := newSyncOwner()

			orphan := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      testNameOrphan,
					Namespace: testNamespaceDefault,
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion:         appsv1.SchemeGroupVersion.String(),
							Kind:               testKindDeployment,
							Name:               owner.Name,
							UID:                owner.UID,
							Controller:         new(true),
							BlockOwnerDeletion: new(true),
						},
					},
				},
			}

			c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(orphan).Build()
			s := NewSyncer(c, scheme, "")

			keep := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: testNameKeep, Namespace: testNamespaceDefault}}
			keep.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))

			// existing mirrors a caller whose GVK restoration either failed
			// (empty TypeMeta on every item, including the just-applied keep —
			// the production failure that would prune the entire desired set)
			// or succeeded (stamped TypeMeta).
			existing := func(ctx context.Context) ([]Object, error) {
				list := &corev1.ConfigMapList{}
				if err := c.List(ctx, list, kclient.InNamespace(testNamespaceDefault)); err != nil {
					return nil, err
				}
				out := make([]Object, 0, len(list.Items))
				for i := range list.Items {
					if tt.stampGVK {
						list.Items[i].SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("ConfigMap"))
					} else {
						list.Items[i].TypeMeta = metav1.TypeMeta{}
					}
					out = append(out, &list.Items[i])
				}
				return out, nil
			}

			res, err := s.Sync(context.Background(), owner, []Object{keep}, existing)
			assert.NoError(t, err)

			err = c.Get(context.Background(),
				kclient.ObjectKey{Name: testNameOrphan, Namespace: testNamespaceDefault}, &corev1.ConfigMap{})
			if tt.wantPruned {
				assert.Equal(t, prunedResult, res["/v1/ConfigMap/default/orphan"])
				assert.Error(t, err, "orphan with a known GVK must be pruned")
			} else {
				for key, ope := range res {
					assert.NotEqual(t, prunedResult, ope, "nothing may be pruned when GVKs are unknown (got %s)", key)
				}
				assert.NoError(t, err, "orphan with an empty GVK must never be pruned")
			}

			err = c.Get(context.Background(),
				kclient.ObjectKey{Name: testNameKeep, Namespace: testNamespaceDefault}, &corev1.ConfigMap{})
			assert.NoError(t, err, "the applied object must survive pruning in every case")
		})
	}
}

func TestSyncer_Sync_EmptyDesiredPrunesOwned(t *testing.T) {
	scheme := newSchemeForSyncerTests(t)
	owner := newSyncOwner()

	orphan := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testNameOrphan,
			Namespace: testNamespaceDefault,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         appsv1.SchemeGroupVersion.String(),
					Kind:               testKindDeployment,
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

	res, err := s.Sync(context.Background(), owner, nil, listConfigMaps(c))
	assert.NoError(t, err)
	assert.Contains(t, res, "no_resources")
	assert.Equal(t, prunedResult, res["/v1/ConfigMap/default/orphan"],
		"an empty desired set must still prune every owned object")

	err = c.Get(context.Background(), kclient.ObjectKey{Name: testNameOrphan, Namespace: testNamespaceDefault}, &corev1.ConfigMap{})
	assert.Error(t, err, "owned orphan must be deleted when the desired set is empty")
}

// NOT IMPLEMENTED BELOW
//
// func TestNewSyncer(t *testing.T) {
// 	type args struct {
// 		api             k8s.Client
// 		scheme          *runtime.Scheme
// 		generationLabel string
// 	}
// 	type want struct {
// 		want *Syncer
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *Syncer) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *Syncer) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           api:nil,
// 		           scheme:nil,
// 		           generationLabel:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           api:nil,
// 		           scheme:nil,
// 		           generationLabel:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := NewSyncer(test.args.api, test.args.scheme, test.args.generationLabel)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestSyncer_Sync(t *testing.T) {
// 	type args struct {
// 		ctx      context.Context
// 		owner    Object
// 		desired  []Object
// 		existing func(ctx context.Context) ([]Object, error)
// 	}
// 	type fields struct {
// 		api             k8s.Client
// 		scheme          *runtime.Scheme
// 		generationLabel string
// 	}
// 	type want struct {
// 		want SyncResults
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, SyncResults, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got SyncResults, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           owner:nil,
// 		           desired:nil,
// 		           existing:nil,
// 		       },
// 		       fields: fields {
// 		           api:nil,
// 		           scheme:nil,
// 		           generationLabel:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           owner:nil,
// 		           desired:nil,
// 		           existing:nil,
// 		           },
// 		           fields: fields {
// 		           api:nil,
// 		           scheme:nil,
// 		           generationLabel:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &Syncer{
// 				api:             test.fields.api,
// 				scheme:          test.fields.scheme,
// 				generationLabel: test.fields.generationLabel,
// 			}
//
// 			got, err := s.Sync(test.args.ctx, test.args.owner, test.args.desired, test.args.existing)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_payloadFieldIndexes(t *testing.T) {
// 	type args struct {
// 		t reflect.Type
// 	}
// 	type want struct {
// 		want []int
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []int) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           t:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           t:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := payloadFieldIndexes(test.args.t)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_applyDesiredState(t *testing.T) {
// 	type args struct {
// 		obj     Object
// 		desired Object
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           obj:nil,
// 		           desired:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           obj:nil,
// 		           desired:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			err := applyDesiredState(test.args.obj, test.args.desired)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestSyncer_prune(t *testing.T) {
// 	type args struct {
// 		ctx      context.Context
// 		owner    Object
// 		applied  SyncResults
// 		existing func(ctx context.Context) ([]Object, error)
// 	}
// 	type fields struct {
// 		api             k8s.Client
// 		scheme          *runtime.Scheme
// 		generationLabel string
// 	}
// 	type want struct {
// 		want SyncResults
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, SyncResults, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got SyncResults, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           owner:nil,
// 		           applied:nil,
// 		           existing:nil,
// 		       },
// 		       fields: fields {
// 		           api:nil,
// 		           scheme:nil,
// 		           generationLabel:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           owner:nil,
// 		           applied:nil,
// 		           existing:nil,
// 		           },
// 		           fields: fields {
// 		           api:nil,
// 		           scheme:nil,
// 		           generationLabel:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &Syncer{
// 				api:             test.fields.api,
// 				scheme:          test.fields.scheme,
// 				generationLabel: test.fields.generationLabel,
// 			}
//
// 			got, err := s.prune(test.args.ctx, test.args.owner, test.args.applied, test.args.existing)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_syncKey(t *testing.T) {
// 	type args struct {
// 		obj Object
// 	}
// 	type want struct {
// 		want string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got string) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           obj:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           obj:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := syncKey(test.args.obj)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
