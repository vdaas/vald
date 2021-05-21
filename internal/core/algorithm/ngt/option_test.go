//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package ngt provides implementation of Go API for https://github.com/yahoojapan/NGT
package ngt

import (
	"math"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/comparator"
	"go.uber.org/goleak"
)

var (
	ngtComparator = []comparator.Option{
		comparator.AllowUnexported(ngt{}),
		// ignore C dependencies
		comparator.IgnoreFields(ngt{},
			"dimension", "prop", "ebuf", "index", "ospace"),
	}
)

func TestWithInMemoryMode(t *testing.T) {
	type T = ngt
	type args struct {
		flg bool
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, ngtComparator...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when flg is true",
			args: args{
				flg: true,
			},
			want: want{
				obj: &T{
					inMemory: true,
				},
			},
		},
		{
			name: "set success when flg is false",
			args: args{
				flg: false,
			},
			want: want{
				obj: &T{
					inMemory: false,
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithInMemoryMode(test.args.flg)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithIndexPath(t *testing.T) {
	type T = ngt
	type args struct {
		path string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, ngtComparator...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when index path is not empty string",
			args: args{
				path: "/tmp/ngt/index",
			},
			want: want{
				obj: &T{
					idxPath: "/tmp/ngt/index",
				},
			},
		},
		{
			name: "set success when the index path is empty string",
			args: args{
				path: "",
			},
			want: want{
				err: errors.New("invalid option, name: indexPath, val: "),
				obj: &T{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithIndexPath(test.args.path)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithBulkInsertChunkSize(t *testing.T) {
	type T = ngt
	type args struct {
		size int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, ngtComparator...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when the size is 100",
			args: args{
				size: 100,
			},
			want: want{
				obj: &T{
					bulkInsertChunkSize: 100,
				},
			},
		},
		{
			name: "set success when the size is 0",
			args: args{
				size: 0,
			},
			want: want{
				obj: &T{
					bulkInsertChunkSize: 0,
				},
			},
		},
		{
			name: "set success when the size is -100",
			args: args{
				size: -100,
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid option, name: BulkInsertChunkSize, val: -100"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithBulkInsertChunkSize(test.args.size)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDimension(t *testing.T) {
	type T = ngt
	type args struct {
		size int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%v\",\n\t\t\t\twant: \"%v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, ngtComparator...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when the size is 100",
			args: args{
				size: 100,
			},
			want: want{
				obj: func() *T {
					t := &T{
						dimension: 100,
					}
					t.setup()
					return t
				}(),
			},
		},
		{
			name: "set success when the size is 0",
			args: args{
				size: 0,
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid critical option, name: dimension, val: 0: dimension size 0 is invalid, the supporting dimension size must be between 2 ~ 65536"),
			},
		},
		{
			name: "set success when the size is 2",
			args: args{
				size: 2,
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when the size is 65536",
			args: args{
				size: 65536,
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "return error when the size is 1",
			args: args{
				size: 1,
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid critical option, name: dimension, val: 1: dimension size 1 is invalid, the supporting dimension size must be between 2 ~ 65536"),
			},
		},
		{
			name: "return error when the size is -100",
			args: args{
				size: -100,
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid critical option, name: dimension, val: -100: dimension size -100 is invalid, the supporting dimension size must be between 2 ~ 65536"),
			},
		},
		{
			name: "return error when the size is 65537",
			args: args{
				size: 65537,
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid critical option, name: dimension, val: 65537: dimension size 65537 is invalid, the supporting dimension size must be between 2 ~ 65536"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithDimension(test.args.size)
			obj := new(T)
			obj.setup()
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDistanceTypeByString(t *testing.T) {
	type T = ngt
	type args struct {
		dt string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, ngtComparator...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when distance type is l1",
			args: args{
				dt: "l1",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is L1",
			args: args{
				dt: "L1",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is l2",
			args: args{
				dt: "l2",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is angle",
			args: args{
				dt: "angle",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is ang",
			args: args{
				dt: "ang",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is hamming",
			args: args{
				dt: "hamming",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is ham",
			args: args{
				dt: "ham",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is cosine",
			args: args{
				dt: "cosine",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is cos",
			args: args{
				dt: "cos",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is normalizedangle",
			args: args{
				dt: "normalizedangle",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is nang",
			args: args{
				dt: "nang",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is normalizedcosine",
			args: args{
				dt: "normalizedcosine",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is ncos",
			args: args{
				dt: "ncos",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is jaccard",
			args: args{
				dt: "jaccard",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is jac",
			args: args{
				dt: "jac",
			},
			want: want{
				obj: &T{},
			},
		},

		{
			name: "set success when distance type includes _ character",
			args: args{
				dt: "normalized_angle",
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "return error when distance type is invalid",
			args: args{
				dt: "invalid type",
			},
			want: want{
				err: errors.ErrUnsupportedDistanceType,
				obj: &T{},
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithDistanceTypeByString(test.args.dt)
			obj := new(T)
			obj.setup()
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDistanceType(t *testing.T) {
	type T = ngt
	type args struct {
		t distanceType
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, ngtComparator...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when distance type is L1",
			args: args{
				t: L1,
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is L2",
			args: args{
				t: L2,
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is Angle",
			args: args{
				t: Angle,
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is Hamming",
			args: args{
				t: Hamming,
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is Cosine",
			args: args{
				t: Cosine,
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is NormalizedAngle",
			args: args{
				t: NormalizedAngle,
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is NormalizedCosine",
			args: args{
				t: NormalizedCosine,
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "set success when distance type is Jaccard",
			args: args{
				t: Jaccard,
			},
			want: want{
				obj: &T{},
			},
		},
		{
			name: "return error when distance type is -1",
			args: args{
				t: -1,
			},
			want: want{
				err: errors.ErrUnsupportedDistanceType,
				obj: &T{},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithDistanceType(test.args.t)
			obj := new(T)
			obj.setup()
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithObjectTypeByString(t *testing.T) {
	type T = ngt
	type args struct {
		ot string
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%v\",\n\t\t\t\twant: \"%v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, ngtComparator...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when object type is uint8",
			args: args{
				ot: "uint8",
			},
			want: want{
				obj: &T{
					objectType: Uint8,
				},
			},
		},
		{
			name: "set success when object type is float",
			args: args{
				ot: "float",
			},
			want: want{
				obj: &T{
					objectType: Float,
				},
			},
		},
		{
			name: "set success when object type is UINT8",
			args: args{
				ot: "UINT8",
			},
			want: want{
				obj: &T{
					objectType: Uint8,
				},
			},
		},
		{
			name: "set success when object type is FLOAT",
			args: args{
				ot: "FLOAT",
			},
			want: want{
				obj: &T{
					objectType: Float,
				},
			},
		},
		{
			name: "set success when object type is double",
			args: args{
				ot: "double",
			},
			want: want{
				obj: &T{
					objectType: Float,
				},
			},
		},
		{
			name: "set success when object type is DOUBLE",
			args: args{
				ot: "DOUBLE",
			},
			want: want{
				obj: &T{
					objectType: Float,
				},
			},
		},
		{
			name: "set success when object type is uint8-",
			args: args{
				ot: "uint8",
			},
			want: want{
				obj: &T{
					objectType: Uint8,
				},
			},
		},
		{
			name: "return error when object type is invalid",
			args: args{
				ot: "invalid",
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid critical option, name: objectType, val: 0: unsupported ObjectType"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithObjectTypeByString(test.args.ot)
			obj := new(T)
			obj.setup()
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithObjectType(t *testing.T) {
	type T = ngt
	type args struct {
		t objectType
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%v\",\n\t\t\t\twant: \"%v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, ngtComparator...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when object type is Uint8",
			args: args{
				t: Uint8,
			},
			want: want{
				obj: &T{
					objectType: Uint8,
				},
			},
		},
		{
			name: "set success when object type is Float",
			args: args{
				t: Float,
			},
			want: want{
				obj: &T{
					objectType: Float,
				},
			},
		},
		{
			name: "return error when object type is -1",
			args: args{
				t: -1,
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid critical option, name: objectType, val: -1: unsupported ObjectType"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithObjectType(test.args.t)
			obj := new(T)
			obj.setup()
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithCreationEdgeSize(t *testing.T) {
	type T = ngt
	type args struct {
		size int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, ngtComparator...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when size is 0",
			args: args{
				size: 0,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set success when size is 1",
			args: args{
				size: 1,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set success when size is -1",
			args: args{
				size: -1,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set success when size is MinInt64",
			args: args{
				size: math.MinInt64,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set success when size is MaxInt64",
			args: args{
				size: math.MaxInt64,
			},
			want: want{
				obj: new(T),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithCreationEdgeSize(test.args.size)
			obj := new(T)
			obj.setup()
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithSearchEdgeSize(t *testing.T) {
	type T = ngt
	type args struct {
		size int
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(obj, w.obj, ngtComparator...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when size is 0",
			args: args{
				size: 0,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set success when size is 1",
			args: args{
				size: 1,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set success when size is -1",
			args: args{
				size: -1,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set success when size is MinInt64",
			args: args{
				size: math.MinInt64,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set success when size is MaxInt64",
			args: args{
				size: math.MaxInt64,
			},
			want: want{
				obj: new(T),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithSearchEdgeSize(test.args.size)
			obj := new(T)
			obj.setup()
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDefaultPoolSize(t *testing.T) {
	type T = ngt
	type args struct {
		poolSize uint32
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%v\",\n\t\t\t\twant: \"%v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{

		{
			name: "set success when size is 1",
			args: args{
				poolSize: 1,
			},
			want: want{
				obj: &T{
					poolSize: 1,
				},
			},
		},
		{
			name: "set success when size is MaxUint32",
			args: args{
				poolSize: math.MaxUint32,
			},
			want: want{
				obj: &T{
					poolSize: math.MaxUint32,
				},
			},
		},
		{
			name: "return error when size is 0",
			args: args{
				poolSize: 0,
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid option, name: defaultPoolSize, val: 0"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithDefaultPoolSize(test.args.poolSize)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDefaultRadius(t *testing.T) {
	type T = ngt
	type args struct {
		radius float32
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when radius is 1",
			args: args{
				radius: 1,
			},
			want: want{
				obj: &T{
					radius: 1,
				},
			},
		},
		{
			name: "set success when radius is MaxFloat32",
			args: args{
				radius: math.MaxFloat32,
			},
			want: want{
				obj: &T{
					radius: math.MaxFloat32,
				},
			},
		},
		{
			name: "set success when radius is SmallestNonzeroFloat32",
			args: args{
				radius: math.SmallestNonzeroFloat32,
			},
			want: want{
				obj: &T{
					radius: math.SmallestNonzeroFloat32,
				},
			},
		},
		{
			name: "return error when radius is 0",
			args: args{
				radius: 0,
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid option, name: defaultRadius, val: 0"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithDefaultRadius(test.args.radius)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDefaultEpsilon(t *testing.T) {
	type T = ngt
	type args struct {
		epsilon float32
	}
	type want struct {
		obj *T
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when epsilon is 1",
			args: args{
				epsilon: 1,
			},
			want: want{
				obj: &T{
					epsilon: 1,
				},
			},
		},
		{
			name: "set success when epsilon is MaxFloat32",
			args: args{
				epsilon: math.MaxFloat32,
			},
			want: want{
				obj: &T{
					epsilon: math.MaxFloat32,
				},
			},
		},
		{
			name: "set success when epsilon is SmallestNonzeroFloat32",
			args: args{
				epsilon: math.SmallestNonzeroFloat32,
			},
			want: want{
				obj: &T{
					epsilon: math.SmallestNonzeroFloat32,
				},
			},
		},
		{
			name: "return error when epsilon is 0",
			args: args{
				epsilon: 0,
			},
			want: want{
				obj: &T{},
				err: errors.New("invalid option, name: defaultEpsilon, val: 0"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				tt.Cleanup(func() { test.afterFunc(test.args) })
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithDefaultEpsilon(test.args.epsilon)
			obj := new(T)
			if err := test.checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
