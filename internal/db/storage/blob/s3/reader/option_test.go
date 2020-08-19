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

package reader

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/go-cmp/cmp"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

var (
	// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
	}
)

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
			return errors.Errorf("got = %v, want %v", obj, w.obj)
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithErrGroup(test.args.eg)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithService(t *testing.T) {
	type T = reader
	type args struct {
		s *s3.S3
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithService(test.args.s)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
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
			return errors.Errorf("got = %v, want %v", obj, w.obj)
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

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithBucket(test.args.bucket)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithKey(t *testing.T) {
	type T = reader
	type args struct {
		key string
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
			name: "set success when key is not empty string",
			args: args{
				key: "vdaas",
			},
			want: want{
				obj: &T{
					key: "vdaas",
				},
			},
		},
		{
			name: "set success when key is empty string",
			want: want{
				obj: new(T),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			got := WithKey(test.args.key)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
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
			return errors.Errorf("got = %v, want %v", got, w.obj)
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
			checkFunc: defaultCheckFunc,
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
			checkFunc: defaultCheckFunc,
		},
		{
			name: "set success when size is zero",
			want: want{
				obj: new(T),
			},
			checkFunc: defaultCheckFunc,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithMaxChunkSize(test.args.size)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
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
			return errors.Errorf("got = %v, want %v", got, w.obj)
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
			checkFunc: defaultCheckFunc,
		},
		{
			name: "set success when backoff enabled is false",
			want: want{
				obj: new(T),
			},
			checkFunc: defaultCheckFunc,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithBackoff(test.args.enabled)
			obj := new(T)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func TestWithBackoffOpts(t *testing.T) {
	type T = reader
	type args struct {
		opts        []backoff.Option
		defaultOpts []backoff.Option
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
				return reflect.DeepEqual(want, got) || (got != nil && want != nil)
			}),
		}
		if diff := cmp.Diff(w.obj, got, opts...); diff != "" {
			fmt.Println(diff)
			return errors.Errorf("got = %v, want %v", got, w.obj)
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
				checkFunc: defaultCheckFunc,
			}
		}(),
		func() test {
			defaultOpts := []backoff.Option{}
			opts := []backoff.Option{
				backoff.WithRetryCount(1),
			}
			return test{
				name: "set success when opts is not nil and backoffOpts is not nil",
				args: args{
					opts:        opts,
					defaultOpts: defaultOpts,
				},
				want: want{
					obj: &T{
						backoffOpts: append(defaultOpts, opts...),
					},
				},
				beforeFunc: func(args args, r *T) {
					r.backoffOpts = args.defaultOpts
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
		{
			name: "set success when opts is nil",
			want: want{
				obj: new(T),
			},
			checkFunc: defaultCheckFunc,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			obj := new(T)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args, obj)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args, obj)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := WithBackoffOpts(test.args.opts...)
			got(obj)
			if err := test.checkFunc(test.want, obj); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
