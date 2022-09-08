//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

package reader

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/go-cmp/cmp"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/db/storage/blob/s3/sdk/s3/s3iface"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestWithErrGroup(t *testing.T) {
	type T = reader
	type args struct {
		eg errgroup.Group
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when eg is not nil",
			args: args{
				eg: errgroup.Get(),
			},
			want: want{
				obj: &T{
					eg: errgroup.Get(),
				},
			},
		},
		{
			name: "set success when eg is nil",
			want: want{
				obj: new(T),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			got := WithErrGroup(test.args.eg)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithService(t *testing.T) {
	type T = reader
	type args struct {
		s s3iface.S3API
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when service is not nil",
			args: args{
				s: new(s3.S3),
			},
			want: want{
				obj: &T{
					service: new(s3.S3),
				},
			},
		},
		{
			name: "set success when service is nil",
			want: want{
				obj: new(T),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			got := WithService(test.args.s)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithBucket(t *testing.T) {
	type T = reader
	type args struct {
		bucket string
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
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when bucket is not nil",
			args: args{
				bucket: "vald",
			},
			want: want{
				obj: &T{
					bucket: "vald",
				},
			},
		},
		{
			name: "set success when bucket is nil",
			want: want{
				obj: new(T),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			got := WithBucket(test.args.bucket)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithMaxChunkSize(t *testing.T) {
	type T = reader
	type args struct {
		size int64
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
	defaultCheckFunc := func(w want, got *T) error {
		if !reflect.DeepEqual(got, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when size is positive number",
			args: args{
				size: 10,
			},
			want: want{
				obj: &T{
					maxChunkSize: 10,
				},
			},
		},
		{
			name: "set success when size is negative number",
			args: args{
				size: -10,
			},
			want: want{
				obj: &T{
					maxChunkSize: -10,
				},
			},
		},
		{
			name: "set success when size is zero",
			want: want{
				obj: new(T),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := WithMaxChunkSize(test.args.size)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithBackoff(t *testing.T) {
	type T = reader
	type args struct {
		enabled bool
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
	defaultCheckFunc := func(w want, got *T) error {
		if !reflect.DeepEqual(got, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when backoff enabled is true",
			args: args{
				enabled: true,
			},
			want: want{
				obj: &T{
					backoffEnabled: true,
				},
			},
		},
		{
			name: "set success when backoff enabled is false",
			want: want{
				obj: new(T),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := WithBackoff(test.args.enabled)
			obj := new(T)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithBackoffOpts(t *testing.T) {
	type T = reader
	type args struct {
		opts           []backoff.Option
		defaultOptions []backoff.Option
	}
	type want struct {
		obj *T
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *T) error
		beforeFunc func(args, *T)
		afterFunc  func(args, *T)
	}
	defaultCheckFunc := func(w want, got *T) error {
		opts := []cmp.Option{
			cmp.AllowUnexported(*got),
			cmp.AllowUnexported(*w.obj),
			cmp.Comparer(func(want, got []backoff.Option) bool {
				return len(got) == len(want)
			}),
			cmp.Comparer(func(want, got backoff.Option) bool {
				return reflect.ValueOf(got).Pointer() == reflect.ValueOf(want).Pointer()
			}),
		}
		if diff := cmp.Diff(w.obj, got, opts...); diff != "" {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.obj)
		}
		return nil
	}
	tests := []test{
		func() test {
			opts := []backoff.Option{
				backoff.WithRetryCount(1),
			}
			return test{
				name: "set success when opts is not nil and backoffOpts is nil",
				args: args{
					opts: opts,
				},
				want: want{
					obj: &T{
						backoffOpts: opts,
					},
				},
			}
		}(),
		func() test {
			defaultOptions := []backoff.Option{}
			opts := []backoff.Option{
				backoff.WithRetryCount(1),
			}
			return test{
				name: "set success when opts is not nil and backoffOpts is not nil",
				args: args{
					opts:           opts,
					defaultOptions: defaultOptions,
				},
				want: want{
					obj: &T{
						backoffOpts: append(defaultOptions, opts...),
					},
				},
				beforeFunc: func(args args, r *T) {
					r.backoffOpts = args.defaultOptions
				},
			}
		}(),
		func() test {
			return test{
				name: "set success when opts is nil",
				want: want{
					obj: new(T),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			obj := new(T)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args, obj)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args, obj)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := WithBackoffOpts(test.args.opts...)
			got(obj)
			if err := checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
