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
		want       want
		checkFunc  func(want, Retry) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
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
			fn := func(...any) {}
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
		fn   func(vals ...any) error
		vals []any
	}
	type fields struct {
		warnFn  func(vals ...any)
		errorFn func(vals ...any)
	}
	type test struct {
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
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

			fn := func(vals ...any) error {
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
			var (
				gotWarnFnErr  error
				gotErrorFnErr error
				fnCnt         int
				errorFnCnt    int
			)

			fn := func(vals ...any) error {
				fnCnt++
				return err
			}

			warnFn := func(vals ...any) {
				if len(vals) == 1 {
					gotWarnFnErr = vals[0].(error)
				}
			}

			errorFn := func(vals ...any) {
				errorFnCnt++
				if len(vals) == 1 {
					gotErrorFnErr = vals[0].(error)
				}
			}

			return test{
				name: "does not panic and reports the final error when fn keeps returning error",
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

					// fn is tried 3 times; errorFn is notified after the 2nd and the 3rd failure.
					if fnCnt != 3 {
						return errors.Errorf("fn count: got: %d, want: %d", fnCnt, 3)
					}

					if errorFnCnt != 2 {
						return errors.Errorf("errorFn count: got: %d, want: %d", errorFnCnt, 2)
					}

					return nil
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					if e := recover(); e != nil {
						t.Errorf("unexpected panic: %v", e)
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
				defer test.afterFunc(tt, test.args)
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
		fn     func(format string, vals ...any) error
		format string
		vals   []any
	}
	type fields struct {
		warnFn  func(vals ...any)
		errorFn func(vals ...any)
	}
	type test struct {
		fields     fields
		checkFunc  func() error
		beforeFunc func(args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
	}
	defaultCheckFunc := func() error {
		return nil
	}
	tests := []test{
		func() test {
			var (
				wantCnt    = 1
				wantFormat = "foramt"
				wantVals   = []any{
					"vald",
				}
			)

			var (
				gotCnt    int
				gotFormat string
				gotVals   []any
			)

			fn := func(format string, vals ...any) error {
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
			var (
				gotWarnFnErr  error
				gotErrorFnErr error
				fnCnt         int
				errorFnCnt    int
			)

			fn := func(format string, vals ...any) error {
				fnCnt++
				return err
			}

			warnFn := func(vals ...any) {
				if len(vals) == 1 {
					gotWarnFnErr = vals[0].(error)
				}
			}

			errorFn := func(vals ...any) {
				errorFnCnt++
				if len(vals) == 1 {
					gotErrorFnErr = vals[0].(error)
				}
			}

			return test{
				name: "does not panic and reports the final error when fn keeps returning error",
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

					// fn is tried 3 times; errorFn is notified after the 2nd and the 3rd failure.
					if fnCnt != 3 {
						return errors.Errorf("fn count: got: %d, want: %d", fnCnt, 3)
					}

					if errorFnCnt != 2 {
						return errors.Errorf("errorFn count: got: %d, want: %d", errorFnCnt, 2)
					}

					return nil
				},
				afterFunc: func(t *testing.T, args args) {
					t.Helper()
					if e := recover(); e != nil {
						t.Errorf("unexpected panic: %v", e)
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
				defer test.afterFunc(tt, test.args)
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

// NOT IMPLEMENTED BELOW

// Test_retry_reentrantReport covers the production wiring where warnFn/errorFn
// log through the same failing logger and so call back into Out with the report.
func Test_retry_reentrantReport(t *testing.T) {
	const maxDepth = 10
	tests := []struct {
		call func(r *retry, fn func(...any) error)
		name string
	}{
		{
			name: "Out terminates when warnFn/errorFn re-enter Out with a failing fn",
			call: func(r *retry, fn func(...any) error) {
				r.Out(fn, "msg")
			},
		},
		{
			name: "Outf terminates when warnFn/errorFn re-enter Out with a failing fn",
			call: func(r *retry, fn func(...any) error) {
				r.Outf(func(string, ...any) error { return fn() }, "%s", "msg")
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			fail := errors.New("write failed")
			var writes int
			fn := func(...any) error {
				writes++
				return fail
			}
			r := new(retry)
			var depth, maxSeen, reports int
			reenter := func(vals ...any) {
				reports++
				depth++
				defer func() { depth-- }()
				if depth > maxSeen {
					maxSeen = depth
				}
				if depth > maxDepth {
					return // bound the old unguarded recursion so the test fails instead of overflowing the stack
				}
				r.Out(fn, vals...) // what log.Warn/log.Error do through the glg logger
			}
			r.warnFn = reenter
			r.errorFn = reenter

			test.call(r, fn)

			if maxSeen != 1 {
				tt.Errorf("report nesting depth: got: %d, want: 1", maxSeen)
			}
			// 3 attempts for the message, each failure reported once and each report written once
			if reports != 3 || writes != 6 {
				tt.Errorf("reports/writes: got: %d/%d, want: 3/6", reports, writes)
			}
		})
	}
}
