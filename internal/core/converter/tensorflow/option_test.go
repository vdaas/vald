//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package tensorflow provides implementation of Go API for extract data to vector
package tensorflow

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestWithSessionOptions(t *testing.T) {
	type T = tensorflow
	type args struct {
		opts *SessionOptions
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set default",
			args: args{
				opts: nil,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set value",
			args: args{
				opts: &SessionOptions{},
			},
			want: want{
				obj: &T{
					options: &SessionOptions{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithSessionOptions(test.args.opts)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithSessionTarget(t *testing.T) {
	type T = tensorflow
	type args struct {
		tgt string
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set default",
			args: args{
				tgt: "",
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set value",
			args: args{
				tgt: "test",
			},
			want: want{
				obj: &T{
					sessionTarget: "test",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithSessionTarget(test.args.tgt)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithSessionConfig(t *testing.T) {
	type T = tensorflow
	type args struct {
		cfg []byte
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set default",
			args: args{
				cfg: nil,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set value",
			args: args{
				cfg: []byte{0},
			},
			want: want{
				obj: &T{
					sessionConfig: []byte{0},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithSessionConfig(test.args.cfg)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithOperations(t *testing.T) {
	type T = tensorflow
	type args struct {
		opes []*Operation
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		obj        *T
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set default",
			args: args{
				opes: nil,
			},
			want: want{
				obj: new(T),
			},
			obj: new(T),
		},
		{
			name: "set value: tensorflow.operations != nil",
			args: args{
				opes: []*Operation{},
			},
			want: want{
				obj: &T{
					operations: []*Operation{},
				},
			},
			obj: &T{
				operations: []*Operation{},
			},
		},
		{
			name: "set value: tensorflow.operations == nil",
			args: args{
				opes: []*Operation{},
			},
			want: want{
				obj: &T{
					operations: []*Operation{},
				},
			},
			obj: new(T),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithOperations(test.args.opes...)
			got(test.obj)
			if err := test.checkFunc(test.want, test.obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithExportPath(t *testing.T) {
	type T = tensorflow
	type args struct {
		path string
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set default",
			args: args{
				path: "",
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set value",
			args: args{
				path: "test",
			},
			want: want{
				obj: &T{
					exportDir: "test",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithExportPath(test.args.path)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithTags(t *testing.T) {
	type T = tensorflow
	type args struct {
		tags []string
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		obj        *T
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set default",
			args: args{
				tags: nil,
			},
			want: want{
				obj: new(T),
			},
			obj: new(T),
		},
		{
			name: "set value: tensorflow.tags != nil",
			args: args{
				tags: []string{"test"},
			},
			want: want{
				obj: &T{
					tags: []string{
						"test",
						"test",
					},
				},
			},
			obj: &T{
				tags: []string{
					"test",
				},
			},
		},
		{
			name: "set value: tensorflow.tags == nil",
			args: args{
				tags: []string{"test"},
			},
			want: want{
				obj: &T{
					tags: []string{
						"test",
					},
				},
			},
			obj: new(T),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithTags(test.args.tags...)
			got(test.obj)
			if err := test.checkFunc(test.want, test.obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithFeed(t *testing.T) {
	type T = tensorflow
	type args struct {
		operationName string
		outputIndex   int
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set value",
			args: args{
				operationName: "test",
				outputIndex:   0,
			},
			want: want{
				obj: &T{
					feeds: []OutputSpec{
						OutputSpec{
							operationName: "test",
							outputIndex:   0,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithFeed(test.args.operationName, test.args.outputIndex)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithFeeds(t *testing.T) {
	type T = tensorflow
	type args struct {
		operationNames []string
		outputIndexes  []int
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set value",
			args: args{
				operationNames: []string{"test"},
				outputIndexes:  []int{0},
			},
			want: want{
				obj: &T{
					feeds: []OutputSpec{
						OutputSpec{
							operationName: "test",
							outputIndex:   0,
						},
					},
				},
			},
		},
		{
			name: "operationNames == nil",
			args: args{
				operationNames: nil,
				outputIndexes:  []int{0},
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "outputIndexes == nil",
			args: args{
				operationNames: []string{"test"},
				outputIndexes:  nil,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "operationName length != outputIndexes length",
			args: args{
				operationNames: []string{},
				outputIndexes:  []int{0},
			},
			want: want{
				obj: new(T),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithFeeds(test.args.operationNames, test.args.outputIndexes)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithFetch(t *testing.T) {
	type T = tensorflow
	type args struct {
		operationName string
		outputIndex   int
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set value",
			args: args{
				operationName: "test",
				outputIndex:   0,
			},
			want: want{
				obj: &T{
					fetches: []OutputSpec{
						OutputSpec{
							operationName: "test",
							outputIndex:   0,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithFetch(test.args.operationName, test.args.outputIndex)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithFetches(t *testing.T) {
	type T = tensorflow
	type args struct {
		operationNames []string
		outputIndexes  []int
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set value",
			args: args{
				operationNames: []string{"test"},
				outputIndexes:  []int{0},
			},
			want: want{
				obj: &T{
					fetches: []OutputSpec{
						OutputSpec{
							operationName: "test",
							outputIndex:   0,
						},
					},
				},
			},
		},
		{
			name: "operationName == nil",
			args: args{
				operationNames: nil,
				outputIndexes:  []int{0},
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "outputIndexs == nil",
			args: args{
				operationNames: []string{"test"},
				outputIndexes:  nil,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "operationNames length != outputIndexs length",
			args: args{
				operationNames: []string{},
				outputIndexes:  []int{0},
			},
			want: want{
				obj: new(T),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithFetches(test.args.operationNames, test.args.outputIndexes)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithNdim(t *testing.T) {
	type T = tensorflow
	type args struct {
		ndim uint8
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args)
		afterFunc  func(args)
	}

	defaultCheckFunc := func(w want, obj *T) error {
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got = %v, want %v", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set defalut",
			args: args{
				ndim: 0,
			},
			want: want{
				obj: new(T),
			},
		},
		{
			name: "set value",
			args: args{
				ndim: 1,
			},
			want: want{
				obj: &T{
					ndim: 1,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithNdim(test.args.ndim)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
