// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package retry

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want Retry
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Retry) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Retry) error {
		wantr, gotr := w.want.(*retry), got.(*retry)

		if reflect.ValueOf(wantr.errorFn).Pointer() != reflect.ValueOf(gotr.errorFn).Pointer() {
			return errors.Errorf("errorFn: got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotr, wantr)
		}

		if reflect.ValueOf(wantr.warnFn).Pointer() != reflect.ValueOf(gotr.warnFn).Pointer() {
			return errors.Errorf("warnFn: got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotr, wantr)
		}

		return nil
	}
	tests := []test{
		{
			name: "returns l when opts is nil",
			want: want{
				want: &retry{
					errorFn: nopFunc,
					warnFn:  nopFunc,
				},
			},
		},

		func() test {
			fn := func(...interface{}) {}
			return test{
				name: "returns l when opts is not nil",
				args: args{
					opts: []Option{
						WithError(fn),
					},
				},
				want: want{
					want: &retry{
						errorFn: fn,
						warnFn:  nopFunc,
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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

			got := New(test.args.opts...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_retry_Out(t *testing.T) {
	type args struct {
		fn   func(vals ...interface{}) error
		vals []interface{}
	}
	type fields struct {
		warnFn  func(vals ...interface{})
		errorFn func(vals ...interface{})
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(args, *testing.T)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		func() test {
			var (
				wantCnt = 1
				gotCnt  = 0
			)

			fn := func(vals ...interface{}) error {
				gotCnt++
				return nil
			}

			return test{
				name: "called success when fn returns nil",
				args: args{
					fn: fn,
				},
				checkFunc: func() error {
					if gotCnt != wantCnt {
						return errors.Errorf("count: got: %d, want: %d", gotCnt, wantCnt)
					}
					return nil
				},
			}
		}(),

		func() test {
			err := errors.New("error")
			fn := func(vals ...interface{}) error {
				return err
			}

			var (
				gotWarnFnErr  error
				gotErrorFnErr error
			)

			warnFn := func(vals ...interface{}) {
				if len(vals) == 1 {
					gotWarnFnErr = vals[0].(error)
				}
			}

			errorFn := func(vals ...interface{}) {
				if len(vals) == 1 {
					gotErrorFnErr = vals[0].(error)
				}
			}

			return test{
				name: "panic occurs when fn returns error",
				args: args{
					fn: fn,
				},
				fields: fields{
					warnFn:  warnFn,
					errorFn: errorFn,
				},
				checkFunc: func() error {
					if !errors.Is(gotErrorFnErr, err) {
						return errors.Errorf("errorFn argument: got: %v, want: %v", gotErrorFnErr, err)
					}

					if !errors.Is(gotWarnFnErr, err) {
						return errors.Errorf("warnFn argument: got: %v, want: %v", gotWarnFnErr, err)
					}

					return nil
				},
				afterFunc: func(args args, t *testing.T) {
					t.Helper()
					if e := recover(); e == nil {
						t.Error("panic dose not occur")
					}
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args, tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			r := &retry{
				warnFn:  test.fields.warnFn,
				errorFn: test.fields.errorFn,
			}

			r.Out(test.args.fn, test.args.vals...)
			if err := checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_retry_Outf(t *testing.T) {
	type args struct {
		fn     func(format string, vals ...interface{}) error
		format string
		vals   []interface{}
	}
	type fields struct {
		warnFn  func(vals ...interface{})
		errorFn func(vals ...interface{})
	}
	type test struct {
		name       string
		args       args
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(args, *testing.T)
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		func() test {
			var (
				wantCnt    = 1
				wantFormat = "foramt"
				wantVals   = []interface{}{
					"vald",
				}
			)

			var (
				gotCnt    int
				gotFormat string
				gotVals   []interface{}
			)

			fn := func(format string, vals ...interface{}) error {
				gotCnt++
				gotFormat = format
				gotVals = vals
				return nil
			}

			return test{
				name: "called success when fn returns nil",
				args: args{
					fn:     fn,
					format: wantFormat,
					vals:   wantVals,
				},
				checkFunc: func() error {
					if gotCnt != wantCnt {
						return errors.Errorf("count: got: %d, want: %d", gotCnt, wantCnt)
					}

					if gotFormat != wantFormat {
						return errors.Errorf("format: got: %d, want: %d", gotFormat, wantFormat)
					}

					if gotCnt != wantCnt {
						return errors.Errorf("vals: got: %d, want: %d", gotVals, wantVals)
					}

					return nil
				},
			}
		}(),

		func() test {
			err := errors.New("error")
			fn := func(format string, vals ...interface{}) error {
				return err
			}

			var (
				gotWarnFnErr  error
				gotErrorFnErr error
			)

			warnFn := func(vals ...interface{}) {
				if len(vals) == 1 {
					gotWarnFnErr = vals[0].(error)
				}
			}

			errorFn := func(vals ...interface{}) {
				if len(vals) == 1 {
					gotErrorFnErr = vals[0].(error)
				}
			}

			return test{
				name: "panic occurs when fn returns error",
				args: args{
					fn: fn,
				},
				fields: fields{
					warnFn:  warnFn,
					errorFn: errorFn,
				},
				checkFunc: func() error {
					if !errors.Is(gotErrorFnErr, err) {
						return errors.Errorf("errorFn argument: got: %v, want: %v", gotErrorFnErr, err)
					}

					if !errors.Is(gotWarnFnErr, err) {
						return errors.Errorf("warnFn argument: got: %v, want: %v", gotWarnFnErr, err)
					}

					return nil
				},
				afterFunc: func(args args, t *testing.T) {
					t.Helper()
					if e := recover(); e == nil {
						t.Error("panic dose not occur")
					}
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args, tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			r := &retry{
				warnFn:  test.fields.warnFn,
				errorFn: test.fields.errorFn,
			}

			r.Outf(test.args.fn, test.args.format, test.args.vals...)
			if err := checkFunc(); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
