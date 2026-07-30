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
	"testing"

	json "github.com/vdaas/vald/internal/encoding/json"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type nestedList = List[nested, *nested]

func newNestedList() *nestedList {
	return &nestedList{
		Items: []nested{
			{Names: []string{"a", "b"}},
			{Names: []string{"c"}},
		},
	}
}

func TestList_DeepCopy(t *testing.T) {
	t.Parallel()

	type test struct {
		mutate func(cp *nestedList)
		check  func(t *testing.T, orig *nestedList)
		name   string
	}

	tests := []test{
		{
			name:   "mutating copied item references does not affect the original",
			mutate: func(cp *nestedList) { cp.Items[0].Names[0] = testMutatedValue },
			check: func(t *testing.T, orig *nestedList) {
				t.Helper()
				if orig.Items[0].Names[0] != "a" {
					t.Errorf("original items mutated: %v", orig.Items)
				}
			},
		},
		{
			name:   "appending to the copied items does not affect the original",
			mutate: func(cp *nestedList) { cp.Items = append(cp.Items, nested{Names: []string{"x"}}) },
			check: func(t *testing.T, orig *nestedList) {
				t.Helper()
				if len(orig.Items) != 2 {
					t.Errorf("original items length changed: %d", len(orig.Items))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orig := newNestedList()
			cp := orig.DeepCopy()
			if cp == nil {
				t.Fatal("DeepCopy() = nil, want copy")
			}
			tc.mutate(cp)
			tc.check(t, orig)
		})
	}

	t.Run("nil receiver returns nil", func(t *testing.T) {
		t.Parallel()
		var l *nestedList
		if got := l.DeepCopy(); got != nil {
			t.Errorf("nil DeepCopy() = %v, want nil", got)
		}
		if got := l.DeepCopyObject(); got != nil {
			t.Errorf("nil DeepCopyObject() = %v, want nil", got)
		}
	})

	t.Run("DeepCopyObject returns a runtime.Object copy", func(t *testing.T) {
		t.Parallel()
		orig := newNestedList()
		obj := orig.DeepCopyObject()
		cp, ok := obj.(*nestedList)
		if !ok {
			t.Fatalf("DeepCopyObject() = %T, want *nestedList", obj)
		}
		cp.Items[0].Names[0] = testMutatedValue
		if orig.Items[0].Names[0] != "a" {
			t.Errorf("original mutated through DeepCopyObject: %v", orig.Items)
		}
	})
}

func TestList_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	orig := newNestedList()
	orig.Kind = "NestedList"
	orig.APIVersion = "vald.vdaas.org/v1"

	raw, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := new(nestedList)
	if err := json.Unmarshal(raw, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Kind != orig.Kind || got.APIVersion != orig.APIVersion {
		t.Errorf("TypeMeta round-trip mismatch: %+v", got.TypeMeta)
	}
	if len(got.Items) != len(orig.Items) || got.Items[0].Names[0] != "a" {
		t.Errorf("items round-trip mismatch: %+v", got.Items)
	}
}

func TestAddListToScheme(t *testing.T) {
	t.Parallel()

	gv := schema.GroupVersion{Group: "vald.vdaas.org", Version: "v1"}
	s := runtime.NewScheme()
	AddListToScheme[nested](s, gv, "NestedList")

	gvks, unversioned, err := s.ObjectKinds(new(nestedList))
	if err != nil {
		t.Fatalf("ObjectKinds: %v", err)
	}
	if unversioned {
		t.Error("list registered as unversioned")
	}
	want := gv.WithKind("NestedList")
	found := false
	for _, gvk := range gvks {
		if gvk == want {
			found = true
		}
	}
	if !found {
		t.Errorf("ObjectKinds = %v, want to contain %v", gvks, want)
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestList_DeepCopyInto(t *testing.T) {
// 	type args struct {
// 		out *List[T, PT]
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		in         *List[T, PT]
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
// 		           out:List[T, PT]{},
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
// 		           out:List[T, PT]{},
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
// 			in := &List[T, PT]{}
//
// 			in.DeepCopyInto(test.args.out)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestList_DeepCopyObject(t *testing.T) {
// 	type want struct {
// 		want runtime.Object
// 	}
// 	type test struct {
// 		name       string
// 		in         *List[T, PT]
// 		want       want
// 		checkFunc  func(want, runtime.Object) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got runtime.Object) error {
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
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
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
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
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
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			in := &List[T, PT]{}
//
// 			got := in.DeepCopyObject()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
