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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vdaas/vald/internal/k8s"
	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	testDataKey   = "key"
	testDataValue = "value"
)

func testConfigMapGVK() schema.GroupVersionKind {
	return corev1.SchemeGroupVersion.WithKind(testKindConfigMap)
}

// schemeOverrideClient wraps a working client but reports a different scheme,
// simulating a client whose scheme does not know the fetched type (GVK
// restoration must then skip without failing the fetch).
type schemeOverrideClient struct {
	k8s.Client
	scheme *runtime.Scheme
}

func (c *schemeOverrideClient) Scheme() *runtime.Scheme { return c.scheme }

func newGVKRestoreClient(t *testing.T, emptyScheme bool, objs ...k8s.Object) k8s.Client {
	t.Helper()
	scheme := runtime.NewScheme()
	assert.NoError(t, corev1.AddToScheme(scheme))
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
	if emptyScheme {
		return &schemeOverrideClient{Client: c, scheme: runtime.NewScheme()}
	}
	return c
}

func TestGetObject_GVKRestore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		wantGVK     schema.GroupVersionKind
		name        string
		emptyScheme bool
	}{
		{
			name:    "restores GVK for a scheme-registered type",
			wantGVK: testConfigMapGVK(),
		},
		{
			name:        "skips restore without error for a type the scheme cannot resolve",
			emptyScheme: true,
			wantGVK:     schema.GroupVersionKind{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := newGVKRestoreClient(t, tc.emptyScheme,
				newConfigMap("cfg", "ns", map[string]string{testDataKey: testDataValue}))
			obj, err := GetObject(context.Background(), c, "cfg", "ns", &corev1.ConfigMap{})
			assert.NoError(t, err)
			assert.Equal(t, tc.wantGVK, obj.GetObjectKind().GroupVersionKind())
		})
	}
}

func TestRefreshObject_GVKRestore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		wantGVK     schema.GroupVersionKind
		name        string
		emptyScheme bool
	}{
		{
			name:    "restores GVK for a scheme-registered type",
			wantGVK: testConfigMapGVK(),
		},
		{
			name:        "skips restore without error for a type the scheme cannot resolve",
			emptyScheme: true,
			wantGVK:     schema.GroupVersionKind{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := newGVKRestoreClient(t, tc.emptyScheme,
				newConfigMap("cfg", "ns", map[string]string{testDataKey: testDataValue}))
			obj, err := RefreshObject(context.Background(), c,
				newConfigMap("cfg", "ns", nil))
			assert.NoError(t, err)
			assert.Equal(t, tc.wantGVK, obj.GetObjectKind().GroupVersionKind())
		})
	}
}

func TestListObjects_GVKRestore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		wantList    schema.GroupVersionKind
		wantItem    schema.GroupVersionKind
		name        string
		emptyScheme bool
	}{
		{
			name:     "restores the list GVK and derives the item GVK",
			wantList: corev1.SchemeGroupVersion.WithKind(testKindConfigMap + "List"),
			wantItem: testConfigMapGVK(),
		},
		{
			name:        "skips restore without error for a list the scheme cannot resolve",
			emptyScheme: true,
			wantList:    schema.GroupVersionKind{},
			wantItem:    schema.GroupVersionKind{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := newGVKRestoreClient(t, tc.emptyScheme,
				newConfigMap("a", "ns", nil), newConfigMap("b", "ns", nil))
			list, err := ListObjects(context.Background(), c, &corev1.ConfigMapList{},
				k8s.InNamespace("ns"))
			assert.NoError(t, err)
			assert.Equal(t, tc.wantList, list.GetObjectKind().GroupVersionKind())
			assert.Len(t, list.Items, 2)
			for i := range list.Items {
				assert.Equal(t, tc.wantItem, list.Items[i].GroupVersionKind(),
					"item %d must carry the derived GVK", i)
			}
		})
	}
}

func TestRestoreObjectGVK(t *testing.T) {
	t.Parallel()

	customGVK := schema.GroupVersionKind{Group: "example.com", Version: "v2", Kind: "Custom"}

	t.Run("nil scheme leaves the object untouched", func(t *testing.T) {
		t.Parallel()

		cm := &corev1.ConfigMap{}
		restoreObjectGVK(nil, cm)
		assert.True(t, cm.GroupVersionKind().Empty())
	})

	t.Run("populated TypeMeta is never overwritten", func(t *testing.T) {
		t.Parallel()

		scheme := runtime.NewScheme()
		assert.NoError(t, corev1.AddToScheme(scheme))
		cm := &corev1.ConfigMap{}
		cm.SetGroupVersionKind(customGVK)
		restoreObjectGVK(scheme, cm)
		assert.Equal(t, customGVK, cm.GroupVersionKind(),
			"restore must not replace an already-populated TypeMeta with the scheme GVK")
	})
}

func TestRestoreListGVK(t *testing.T) {
	t.Parallel()

	itemGVKs := func(t *testing.T, list ObjectListType) []schema.GroupVersionKind {
		t.Helper()
		items, err := apimeta.ExtractList(list)
		assert.NoError(t, err)
		out := make([]schema.GroupVersionKind, 0, len(items))
		for _, item := range items {
			out = append(out, item.GetObjectKind().GroupVersionKind())
		}
		return out
	}

	customGVK := schema.GroupVersionKind{Group: "example.com", Version: "v2", Kind: "Custom"}
	testGV := schema.GroupVersion{Group: "test.vdaas.org", Version: "v1"}

	t.Run("fills empty item TypeMeta and preserves populated ones", func(t *testing.T) {
		t.Parallel()

		scheme := runtime.NewScheme()
		assert.NoError(t, corev1.AddToScheme(scheme))
		list := &corev1.ConfigMapList{Items: []corev1.ConfigMap{{}, {}}}
		list.Items[0].SetGroupVersionKind(customGVK)

		restoreListGVK(scheme, list)
		assert.Equal(t,
			corev1.SchemeGroupVersion.WithKind(testKindConfigMap+listKindSuffix),
			list.GroupVersionKind())
		assert.Equal(t,
			[]schema.GroupVersionKind{customGVK, testConfigMapGVK()},
			itemGVKs(t, list))
	})

	t.Run("derives the item kind for a generic List registered via AddListToScheme", func(t *testing.T) {
		t.Parallel()

		scheme := runtime.NewScheme()
		AddListToScheme[corev1.ConfigMap](scheme, testGV, testKindConfigMap+listKindSuffix)
		list := &List[corev1.ConfigMap, *corev1.ConfigMap]{Items: []corev1.ConfigMap{{}}}

		restoreListGVK(scheme, list)
		assert.Equal(t, testGV.WithKind(testKindConfigMap+listKindSuffix), list.GroupVersionKind())
		assert.Equal(t,
			[]schema.GroupVersionKind{testGV.WithKind(testKindConfigMap)},
			itemGVKs(t, list))
	})

	t.Run("leaves items untouched when the list kind lacks the List suffix", func(t *testing.T) {
		t.Parallel()

		list := &corev1.ConfigMapList{Items: []corev1.ConfigMap{{}}}
		list.SetGroupVersionKind(testGV.WithKind("Weird"))

		restoreListGVK(runtime.NewScheme(), list)
		assert.Equal(t, testGV.WithKind("Weird"), list.GroupVersionKind())
		assert.Equal(t, []schema.GroupVersionKind{{}}, itemGVKs(t, list))
	})

	t.Run("leaves items untouched when the list kind is exactly List", func(t *testing.T) {
		t.Parallel()

		list := &corev1.ConfigMapList{Items: []corev1.ConfigMap{{}}}
		list.SetGroupVersionKind(testGV.WithKind(listKindSuffix))

		restoreListGVK(runtime.NewScheme(), list)
		assert.Equal(t, []schema.GroupVersionKind{{}}, itemGVKs(t, list))
	})
}

func TestToUnstructured(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "cfg", Namespace: "ns"},
		Data:       map[string]string{testDataKey: "value"},
	}
	cm.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind(testKindConfigMap))

	us, err := ToUnstructured(cm)
	assert.NoError(t, err)
	assert.Equal(t, testKindConfigMap, us.GetKind())
	assert.Equal(t, "cfg", us.GetName())
	assert.Equal(t, "ns", us.GetNamespace())

	data, ok := us.Object["data"].(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "value", data[testDataKey])
}

func TestObjectsOf(t *testing.T) {
	items := []corev1.ConfigMap{
		{ObjectMeta: metav1.ObjectMeta{Name: "a", Namespace: "ns"}},
		{ObjectMeta: metav1.ObjectMeta{Name: "b", Namespace: "ns"}},
	}

	out := ObjectsOf(items)
	assert.Len(t, out, 2)
	assert.Equal(t, "a", out[0].GetName())
	assert.Equal(t, "b", out[1].GetName())

	// The returned objects must alias the input slice elements, not copies.
	out[0].SetName("mutated")
	assert.Equal(t, "mutated", items[0].GetName())
}

// NOT IMPLEMENTED BELOW
//
// func TestIgnoreNotFound(t *testing.T) {
// 	type args struct {
// 		err error
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
// 		           err:nil,
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
// 		           err:nil,
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
// 			err := IgnoreNotFound(test.args.err)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestSemanticDeepEqual(t *testing.T) {
// 	type args struct {
// 		a1 any
// 		a2 any
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
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
// 		           a1:nil,
// 		           a2:nil,
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
// 		           a1:nil,
// 		           a2:nil,
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
// 			got := SemanticDeepEqual(test.args.a1, test.args.a2)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_restoreObjectGVK(t *testing.T) {
// 	type args struct {
// 		scheme *runtime.Scheme
// 		obj    runtime.Object
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           scheme:nil,
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
// 		           scheme:nil,
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
// 			restoreObjectGVK(test.args.scheme, test.args.obj)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_restoreListGVK(t *testing.T) {
// 	type args struct {
// 		scheme *runtime.Scheme
// 		list   ObjectListType
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           scheme:nil,
// 		           list:nil,
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
// 		           scheme:nil,
// 		           list:nil,
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
// 			restoreListGVK(test.args.scheme, test.args.list)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestGetObject(t *testing.T) {
// 	type args struct {
// 		ctx       context.Context
// 		c         k8s.Client
// 		name      string
// 		namespace string
// 		obj       T
// 	}
// 	type want struct {
// 		want T
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, T, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got T, err error) error {
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
// 		           c:nil,
// 		           name:"",
// 		           namespace:"",
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
// 		           ctx:nil,
// 		           c:nil,
// 		           name:"",
// 		           namespace:"",
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
// 			got, err := GetObject(test.args.ctx, test.args.c, test.args.name, test.args.namespace, test.args.obj)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestRefreshObject(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		c   k8s.Client
// 		obj T
// 	}
// 	type want struct {
// 		want T
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, T, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got T, err error) error {
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
// 		           c:nil,
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
// 		           ctx:nil,
// 		           c:nil,
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
// 			got, err := RefreshObject(test.args.ctx, test.args.c, test.args.obj)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestListObjects(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		c    k8s.Client
// 		list L
// 		opts []ListOption
// 	}
// 	type want struct {
// 		want L
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, L, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got L, err error) error {
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
// 		           c:nil,
// 		           list:nil,
// 		           opts:nil,
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
// 		           c:nil,
// 		           list:nil,
// 		           opts:nil,
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
// 			got, err := ListObjects(test.args.ctx, test.args.c, test.args.list, test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
