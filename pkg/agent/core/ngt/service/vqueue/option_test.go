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

// Package vqueue manages the vector cache layer for reducing FFI overhead for fast Agent processing.
package vqueue

import (
	"context"
	"math"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestWithErrGroup(t *testing.T) {
	type T = vqueue
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			eg, _ := errgroup.New(context.Background())
			return test{
				name: "set success when the eg is not nil",
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
				name: "set fails when the eg is nil",
				args: args{},
				want: want{
					err: errors.NewErrInvalidOption("errgroup", nil),
					obj: new(T),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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

func TestWithInsertBufferPoolSize(t *testing.T) {
	type T = vqueue
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
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			size := 100
			return test{
				name: "set success when size is 100",
				args: args{
					size: size,
				},
				want: want{
					obj: &T{
						iBufSize: size,
					},
				},
			}
		}(),
		func() test {
			size := 1
			return test{
				name: "set success when size is 1",
				args: args{
					size: size,
				},
				want: want{
					obj: &T{
						iBufSize: size,
					},
				},
			}
		}(),
		func() test {
			size := math.MaxInt64
			return test{
				name: "set success when size is maximum value of int",
				args: args{
					size: size,
				},
				want: want{
					obj: &T{
						iBufSize: size,
					},
				},
			}
		}(),
		func() test {
			size := 0
			return test{
				name: "set fails when size is 0",
				args: args{
					size: size,
				},
				want: want{
					err: errors.NewErrInvalidOption("insertBufferPoolSize", size),
					obj: new(T),
				},
			}
		}(),
		func() test {
			size := -1
			return test{
				name: "set fails when size is -1",
				args: args{
					size: size,
				},
				want: want{
					err: errors.NewErrInvalidOption("insertBufferPoolSize", size),
					obj: new(T),
				},
			}
		}(),
		func() test {
			size := -100
			return test{
				name: "set fails when size is -100",
				args: args{
					size: size,
				},
				want: want{
					err: errors.NewErrInvalidOption("insertBufferPoolSize", size),
					obj: new(T),
				},
			}
		}(),
		func() test {
			size := math.MinInt64
			return test{
				name: "set success when size is minimum value of int",
				args: args{
					size: size,
				},
				want: want{
					err: errors.NewErrInvalidOption("insertBufferPoolSize", size),
					obj: new(T),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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

			got := WithInsertBufferPoolSize(test.args.size)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWithDeleteBufferPoolSize(t *testing.T) {
	type T = vqueue
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
		if !reflect.DeepEqual(obj, w.obj) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", obj, w.obj)
		}
		return nil
	}

	tests := []test{
		func() test {
			size := 100
			return test{
				name: "set success when size is 100",
				args: args{
					size: size,
				},
				want: want{
					obj: &T{
						dBufSize: size,
					},
				},
			}
		}(),
		func() test {
			size := 1
			return test{
				name: "set success when size is 1",
				args: args{
					size: size,
				},
				want: want{
					obj: &T{
						dBufSize: size,
					},
				},
			}
		}(),
		func() test {
			size := math.MaxInt64
			return test{
				name: "set success when size is maximum value of int",
				args: args{
					size: size,
				},
				want: want{
					obj: &T{
						dBufSize: size,
					},
				},
			}
		}(),
		func() test {
			size := 0
			return test{
				name: "set fails when size is 0",
				args: args{
					size: size,
				},
				want: want{
					err: errors.NewErrInvalidOption("deleteBufferPoolSize", size),
					obj: new(T),
				},
			}
		}(),
		func() test {
			size := -1
			return test{
				name: "set fails when size is -1",
				args: args{
					size: size,
				},
				want: want{
					err: errors.NewErrInvalidOption("deleteBufferPoolSize", size),
					obj: new(T),
				},
			}
		}(),
		func() test {
			size := -100
			return test{
				name: "set fails when size is -100",
				args: args{
					size: size,
				},
				want: want{
					err: errors.NewErrInvalidOption("deleteBufferPoolSize", size),
					obj: new(T),
				},
			}
		}(),
		func() test {
			size := math.MinInt64
			return test{
				name: "set success when size is minimum value of int",
				args: args{
					size: size,
				},
				want: want{
					err: errors.NewErrInvalidOption("deleteBufferPoolSize", size),
					obj: new(T),
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
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

			got := WithDeleteBufferPoolSize(test.args.size)
			obj := new(T)
			if err := checkFunc(test.want, obj, got(obj)); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
