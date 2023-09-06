//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

package errors

import (
	"math"
	"runtime"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/test/goleak"
)

func TestErrPanicRecovered(t *testing.T) {
	type args struct {
		err error
		rec interface{}
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	defaultErr := New("runtime error")
	tests := []test{
		func() test {
			r := math.MaxFloat64
			return test{
				name: "return an error when err is not empty and rec is int value",
				args: args{
					err: defaultErr,
					rec: r,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("panic recovered: %v", r).Error()),
				},
			}
		}(),
		func() test {
			r := math.MaxFloat64
			return test{
				name: "return an error when err is not empty and rec is float value",
				args: args{
					err: defaultErr,
					rec: r,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("panic recovered: %v", r).Error()),
				},
			}
		}(),
		func() test {
			tString := "10h"
			var r time.Duration
			r, _ = time.ParseDuration(tString)
			return test{
				name: "return an error when err is not empty and rec is time.Duration value",
				args: args{
					err: defaultErr,
					rec: r,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("panic recovered: %v", r).Error()),
				},
			}
		}(),
		func() test {
			r := "10h"
			return test{
				name: "return an error when err is not empty and rec is string value",
				args: args{
					err: defaultErr,
					rec: r,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("panic recovered: %v", r).Error()),
				},
			}
		}(),
		func() test {
			r := []byte{0x00, 0x01}
			return test{
				name: "return an error when err is not empty and rec is byte slice value",
				args: args{
					err: defaultErr,
					rec: r,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("panic recovered: %v", r).Error()),
				},
			}
		}(),
		func() test {
			r := struct {
				connectionType string
				code           int
			}{
				connectionType: "grpc",
				code:           404,
			}
			return test{
				name: "return an error when err is not empty and rec is struct value",
				args: args{
					err: defaultErr,
					rec: r,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("panic recovered: %v", r).Error()),
				},
			}
		}(),
		func() test {
			r := map[string]int{
				"uuid": 12345678,
			}
			return test{
				name: "return an error when err is not empty and rec is map value",
				args: args{
					err: defaultErr,
					rec: r,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("panic recovered: %v", r).Error()),
				},
			}
		}(),
		func() test {
			r := true
			return test{
				name: "return an error when err is not empty and rec is boolean",
				args: args{
					err: defaultErr,
					rec: r,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("panic recovered: %v", r).Error()),
				},
			}
		}(),
		func() test {
			r := &defaultErr
			return test{
				name: "return an error when err is not empty and rec is pointer",
				args: args{
					err: defaultErr,
					rec: r,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("panic recovered: %v", r).Error()),
				},
			}
		}(),
		func() test {
			r := "10h"
			return test{
				name: "return an error when err is empty and rec has value",
				args: args{
					rec: r,
				},
				want: want{
					want: Wrap(nil, Errorf("panic recovered: %v", r).Error()),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an error when err is not empty and rec is nil",
				args: args{
					err: defaultErr,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("panic recovered: %v", nil).Error()),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an error when err and rec are empty or nil",
				args: args{},
				want: want{
					want: Wrap(nil, Errorf("panic recovered: %v", nil).Error()),
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

			got := ErrPanicRecovered(test.args.err, test.args.rec)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrPanicString(t *testing.T) {
	type args struct {
		err error
		msg string
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	defaultErr := New("runtime panic string error")
	defaultMsg := "success"
	tests := []test{
		func() test {
			return test{
				name: "return an error when err is not nil and msg is not empty",
				args: args{
					err: defaultErr,
					msg: defaultMsg,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("panic recovered: %v", defaultMsg).Error()),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an error when err is not nil and msg is empty",
				args: args{
					err: defaultErr,
				},
				want: want{
					want: Wrap(defaultErr, New("panic recovered: ").Error()),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an error when err is nil and msg is not empty",
				args: args{
					msg: defaultMsg,
				},
				want: want{
					want: Wrap(nil, Errorf("panic recovered: %v", defaultMsg).Error()),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an error when err is nil and msg is empty",
				args: args{},
				want: want{
					want: Wrap(nil, New("panic recovered: ").Error()),
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

			got := ErrPanicString(test.args.err, test.args.msg)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

type runtimeErr struct {
	err error
}

func (runtimeErr) RuntimeError() {}
func (e runtimeErr) Error() string {
	return e.err.Error()
}

func TestErrRuntimeError(t *testing.T) {
	type args struct {
		err error
		r   runtime.Error
	}
	type want struct {
		want error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got error) error {
		if !Is(got, w.want) {
			return Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	defaultErr := New("runtime panic string error")
	defaultRuntimeErr := &runtimeErr{
		err: New("runtime error is occurred"),
	}
	tests := []test{
		func() test {
			return test{
				name: "return an error when err is not nil and r is not empty",
				args: args{
					err: defaultErr,
					r:   defaultRuntimeErr,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("system panicked caused by runtime error: %v", defaultRuntimeErr).Error()),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an error when err is not nil and msg is nil",
				args: args{
					err: defaultErr,
				},
				want: want{
					want: Wrap(defaultErr, Errorf("system panicked caused by runtime error: %v", nil).Error()),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an error when err is nil and msg is not empty",
				args: args{
					r: defaultRuntimeErr,
				},
				want: want{
					want: Wrap(nil, Errorf("system panicked caused by runtime error: %v", defaultRuntimeErr).Error()),
				},
			}
		}(),
		func() test {
			return test{
				name: "return an error when err is nil and msg is nil",
				args: args{},
				want: want{
					want: Wrap(nil, Errorf("system panicked caused by runtime error: %v", nil).Error()),
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

			got := ErrRuntimeError(test.args.err, test.args.r)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
