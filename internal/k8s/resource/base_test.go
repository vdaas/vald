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
	"bytes"
	"maps"
	"slices"
	"testing"

	json "github.com/vdaas/vald/internal/encoding/json"
)

const (
	testAppLabelValue = "vald"
	testMutatedValue  = "mutated"
)

// nested owns a reference field so that the synthetic type exercises the
// DeepCopyIntoer-based helpers.
type nested struct {
	Names []string `json:"names,omitempty"`
}

func (in *nested) DeepCopyInto(out *nested) {
	*out = *in
	out.Names = slices.Clone(in.Names)
}

// synthetic exercises every helper: slice, map, ptr, nested references.
type synthetic struct { //nolint:tagliatelle // generic embed field name confuses the linter's tag-name check
	Base[synthetic, *synthetic] `json:"-"`

	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
	Count  *int              `json:"count,omitempty"`
	Nested *nested           `json:"nested,omitempty"`
	Items  []nested          `json:"items,omitempty"`
}

func (in *synthetic) DeepCopyInto(out *synthetic) {
	*out = *in
	out.Labels = maps.Clone(in.Labels)
	out.Count = CopyPtr(in.Count)
	out.Nested = CopyPtrInto(in.Nested)
	out.Items = CopySliceInto(in.Items)
}

// syntheticNoBase mirrors synthetic without the Base embed for the JSON
// invariance check.
type syntheticNoBase struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
	Count  *int              `json:"count,omitempty"`
	Nested *nested           `json:"nested,omitempty"`
	Items  []nested          `json:"items,omitempty"`
}

func newSynthetic() *synthetic {
	count := 7
	return &synthetic{
		Name:   "origin",
		Labels: map[string]string{"app": testAppLabelValue},
		Count:  &count,
		Nested: &nested{Names: []string{"a", "b"}},
		Items:  []nested{{Names: []string{"c"}}, {Names: []string{"d"}}},
	}
}

func TestBase_DeepCopy(t *testing.T) {
	t.Parallel()

	type test struct {
		mutate func(cp *synthetic)
		check  func(t *testing.T, orig *synthetic)
		name   string
	}

	tests := []test{
		{
			name:   "mutating the copied map does not affect the original",
			mutate: func(cp *synthetic) { cp.Labels["app"] = testMutatedValue },
			check: func(t *testing.T, orig *synthetic) {
				t.Helper()
				if orig.Labels["app"] != testAppLabelValue {
					t.Errorf("original map mutated: %v", orig.Labels)
				}
			},
		},
		{
			name:   "mutating the copied pointer does not affect the original",
			mutate: func(cp *synthetic) { *cp.Count = 99 },
			check: func(t *testing.T, orig *synthetic) {
				t.Helper()
				if *orig.Count != 7 {
					t.Errorf("original pointer mutated: %d", *orig.Count)
				}
			},
		},
		{
			name:   "mutating the copied nested reference does not affect the original",
			mutate: func(cp *synthetic) { cp.Nested.Names[0] = testMutatedValue },
			check: func(t *testing.T, orig *synthetic) {
				t.Helper()
				if orig.Nested.Names[0] != "a" {
					t.Errorf("original nested slice mutated: %v", orig.Nested.Names)
				}
			},
		},
		{
			name:   "mutating the copied slice elements does not affect the original",
			mutate: func(cp *synthetic) { cp.Items[1].Names[0] = testMutatedValue },
			check: func(t *testing.T, orig *synthetic) {
				t.Helper()
				if orig.Items[1].Names[0] != "d" {
					t.Errorf("original items mutated: %v", orig.Items)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orig := newSynthetic()
			cp := orig.DeepCopy()
			if cp == nil {
				t.Fatal("DeepCopy() = nil, want copy")
			}
			tc.mutate(cp)
			tc.check(t, orig)
		})
	}
}

func TestBase_NilReceiver(t *testing.T) {
	t.Parallel()

	var b *Base[synthetic, *synthetic]
	if got := b.DeepCopy(); got != nil {
		t.Errorf("nil receiver DeepCopy() = %v, want nil", got)
	}
	if got := b.DeepCopyObject(); got != nil {
		t.Errorf("nil receiver DeepCopyObject() = %v, want nil", got)
	}
}

func TestBase_DeepCopyObject(t *testing.T) {
	t.Parallel()

	// synthetic does not implement runtime.Object (no ObjectKind), so
	// DeepCopyObject must return nil. The positive case is covered by the
	// scheme-registered API types (e.g. ValdOperatorRelease) in their own package.
	if got := newSynthetic().DeepCopyObject(); got != nil {
		t.Errorf("DeepCopyObject() = %v, want nil for non API type", got)
	}
}

func TestBase_JSONMarshalInvariance(t *testing.T) {
	t.Parallel()

	count := 7
	with := newSynthetic()
	without := &syntheticNoBase{
		Name:   "origin",
		Labels: map[string]string{"app": testAppLabelValue},
		Count:  &count,
		Nested: &nested{Names: []string{"a", "b"}},
		Items:  []nested{{Names: []string{"c"}}, {Names: []string{"d"}}},
	}

	got, err := json.Marshal(with)
	if err != nil {
		t.Fatalf("marshal with Base: %v", err)
	}
	want, err := json.Marshal(without)
	if err != nil {
		t.Fatalf("marshal without Base: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("json output changed by Base embed:\nwith:    %s\nwithout: %s", got, want)
	}
}

func TestCopyHelpers_Nil(t *testing.T) {
	t.Parallel()

	if CopyPtr[int](nil) != nil {
		t.Error("CopyPtr(nil) != nil")
	}
	if CopyPtrInto[nested](nil) != nil {
		t.Error("CopyPtrInto(nil) != nil")
	}
	if CopySliceInto[nested](nil) != nil {
		t.Error("CopySliceInto(nil) != nil")
	}
	if CopyPtrSliceInto[nested](nil) != nil {
		t.Error("CopyPtrSliceInto(nil) != nil")
	}
}

// misplacedBase violates the Base contract: the embed is not the first field.
// Field order is the point of this type (VerifyBase must detect the offset
// violation), so it cannot be reordered for memory layout.
type misplacedBase struct { //nolint:govet,tagliatelle // field order deliberately violates Base's first-field contract; generic embed confuses the tag-name check
	Name string

	Base[misplacedBase, *misplacedBase] `json:"-"`
}

func (in *misplacedBase) DeepCopyInto(out *misplacedBase) { *out = *in }

// missingBase violates the Base contract: no Base embed at all.
type missingBase struct {
	Name string
}

func (in *missingBase) DeepCopyInto(out *missingBase) { *out = *in }

func TestVerifyBase(t *testing.T) {
	t.Parallel()

	type test struct {
		verify  func() error
		name    string
		wantErr bool
	}

	tests := []test{
		{
			name:   "valid first-field embed passes",
			verify: VerifyBase[synthetic, *synthetic],
		},
		{
			name:    "embed not at offset 0 fails",
			verify:  VerifyBase[misplacedBase, *misplacedBase],
			wantErr: true,
		},
		{
			name:    "missing embed fails",
			verify:  VerifyBase[missingBase, *missingBase],
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if err := tc.verify(); tc.wantErr != (err != nil) {
				t.Errorf("VerifyBase() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNewQuantity(t *testing.T) {
// 	type args struct {
// 		value  int64
// 		format k8sresource.Format
// 	}
// 	type want struct {
// 		want *Quantity
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *Quantity) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *Quantity) error {
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
// 		           value:0,
// 		           format:nil,
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
// 		           value:0,
// 		           format:nil,
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
// 			got := NewQuantity(test.args.value, test.args.format)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestMustParse(t *testing.T) {
// 	type args struct {
// 		str string
// 	}
// 	type want struct {
// 		want Quantity
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Quantity) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Quantity) error {
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
// 		           str:"",
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
// 		           str:"",
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
// 			got := MustParse(test.args.str)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestBase_self(t *testing.T) {
// 	type want struct {
// 		want PT
// 	}
// 	type test struct {
// 		name       string
// 		b          *Base[T, PT]
// 		want       want
// 		checkFunc  func(want, PT) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got PT) error {
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
// 			b := &Base[T, PT]{}
//
// 			got := b.self()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestCopyPtr(t *testing.T) {
// 	type args struct {
// 		in *T
// 	}
// 	type want struct {
// 		want *T
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *T) error {
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
// 		           in:nil,
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
// 		           in:nil,
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
// 			got := CopyPtr(test.args.in)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestCopyPtrInto(t *testing.T) {
// 	type args struct {
// 		in PT
// 	}
// 	type want struct {
// 		want PT
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, PT) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got PT) error {
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
// 		           in:nil,
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
// 		           in:nil,
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
// 			got := CopyPtrInto(test.args.in)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestCopySliceInto(t *testing.T) {
// 	type args struct {
// 		in []T
// 	}
// 	type want struct {
// 		want []T
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []T) error {
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
// 		           in:nil,
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
// 		           in:nil,
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
// 			got := CopySliceInto(test.args.in)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestCopyPtrSliceInto(t *testing.T) {
// 	type args struct {
// 		in *[]T
// 	}
// 	type want struct {
// 		want *[]T
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *[]T) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *[]T) error {
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
// 		           in:nil,
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
// 		           in:nil,
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
// 			got := CopyPtrSliceInto(test.args.in)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestCopySliceFunc(t *testing.T) {
// 	type args struct {
// 		in S
// 		cp func(E) E
// 	}
// 	type want struct {
// 		want S
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, S) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got S) error {
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
// 		           in:nil,
// 		           cp:nil,
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
// 		           in:nil,
// 		           cp:nil,
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
// 			got := CopySliceFunc(test.args.in, test.args.cp)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
