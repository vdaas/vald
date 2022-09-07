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

// Package grpc provides grpc server logic
package grpc

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func TestWithIP(t *testing.T) {
	t.Parallel()
	type T = server
	type args struct {
		ip string
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant_error: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		{
			name: "set success when ip is not empty",
			args: args{
				ip: "192.168.1.1",
			},
			want: want{
				obj: &T{
					ip: "192.168.1.1",
				},
			},
		},
		{
			name: "set fail when ip is empty",
			args: args{
				ip: "",
			},
			want: want{
				obj: new(T),
				err: errors.New("invalid option, name: ip, val: "),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got := WithIP(test.args.ip)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithName(t *testing.T) {
	t.Parallel()
	type T = server
	type args struct {
		name string
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant_error: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when name is not empty",
			args: args{
				name: "agent handler",
			},
			want: want{
				obj: &T{
					name: "agent handler",
				},
			},
		},
		{
			name: "set fail when name is empty",
			args: args{
				name: "",
			},
			want: want{
				obj: new(T),
				err: errors.New("invalid option, name: name, val: "),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got := WithName(test.args.name)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithNGT(t *testing.T) {
	t.Parallel()
	type T = server
	type args struct {
		n service.NGT
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
		beforeFunc func(*testing.T, *args, *want)
		afterFunc  func(*testing.T, args, want)
	}

	defaultCheckFunc := func(w want, obj *T, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%s\",\n\t\t\t\twant_error: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when ngt is not nil",
			beforeFunc: func(t *testing.T, args *args, w *want) {
				t.Helper()
				n, err := service.New(&config.NGT{
					Dimension:    1024,
					DistanceType: "cos",
					ObjectType:   "uint8",
					VQueue: &config.VQueue{
						InsertBufferPoolSize: 50,
						DeleteBufferPoolSize: 50,
					},
					KVSDB: &config.KVSDB{
						Concurrency: 0,
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				args.n = n
				w.obj = &T{
					ngt: n,
				}
			},
		},
		{
			name: "set fail when ngt is nil",
			args: args{
				n: nil,
			},
			want: want{
				obj: new(T),
				err: errors.New("invalid option, name: ngt, val: <nil>"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(tt, &test.args, &test.want)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args, test.want)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := WithNGT(test.args.n)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithStreamConcurrency(t *testing.T) {
	t.Parallel()
	type T = server
	type args struct {
		c int
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant_error: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		{
			name: "set success when streamConcurrency > 0",
			args: args{
				c: 100,
			},
			want: want{
				obj: &T{
					streamConcurrency: 100,
				},
			},
		},
		{
			name: "set fail when streamConcurrency < 0",
			args: args{
				c: -500,
			},
			want: want{
				obj: &T{
					streamConcurrency: 0,
				},
				err: errors.New("invalid option, name: streamConcurrency, val: -500"),
			},
		},
		{
			name: "set fail when streamConcurrency  0",
			args: args{
				c: 0,
			},
			want: want{
				obj: &T{
					streamConcurrency: 0,
				},
				err: errors.New("invalid option, name: streamConcurrency, val: 0"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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

			got := WithStreamConcurrency(test.args.c)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithErrGroup(t *testing.T) {
	t.Parallel()
	type T = server
	type args struct {
		eg errgroup.Group
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant_error: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}
	tests := []test{
		func() test {
			eg := errgroup.Get()
			return test{
				name: "set success when eg is not nil",
				args: args{
					eg: eg,
				},
				want: want{
					obj: &T{
						eg: eg,
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "set fail when eg is nil",
				args: args{
					eg: nil,
				},
				want: want{
					obj: &T{
						eg: nil,
					},
					err: errors.New("invalid option, name: errGroup, val: <nil>"),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
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
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
