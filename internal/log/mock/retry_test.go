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
package mock

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestRetry_Out(t *testing.T) {
	type args struct {
		fn   func(vals ...interface{}) error
		vals []interface{}
	}
	type fields struct {
		OutFunc func(fn func(vals ...interface{}) error, vals ...interface{})
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fieldsFunc func(*testing.T) fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		func() test {
			var (
				wantFn = func(vals ...interface{}) error {
					return nil
				}
				wantVals = []interface{}{
					"Vald",
				}
				cnt int
			)
			return test{
				name: "Call Out",
				args: args{
					fn:   wantFn,
					vals: wantVals,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						OutFunc: func(fn func(vals ...interface{}) error, vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) {
								t.Errorf("vals got = %v, want = %v", vals, wantVals)
							}
							if reflect.ValueOf(fn).Pointer() != reflect.ValueOf(wantFn).Pointer() {
								t.Errorf("fn got = %p, want = %p", fn, wantFn)
							}
							cnt++
						},
					}
				},
				checkFunc: func(w want) error {
					if cnt != 1 {
						return errors.Errorf("cnt got = %d, want = %d", cnt, 1)
					}
					return nil
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
			fields := test.fieldsFunc(tt)
			r := &Retry{
				OutFunc: fields.OutFunc,
			}

			r.Out(test.args.fn, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestRetry_Outf(t *testing.T) {
	type args struct {
		fn     func(format string, vals ...interface{}) error
		format string
		vals   []interface{}
	}
	type fields struct {
		OutfFunc func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{})
	}
	type want struct{}
	type test struct {
		name       string
		args       args
		fieldsFunc func(*testing.T) fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		func() test {
			var (
				wantFn = func(format string, vals ...interface{}) error {
					return nil
				}
				wantVals = []interface{}{
					"Vald",
				}
				wantFormat = "json"
				cnt        int
			)
			return test{
				name: "Call Outf",
				args: args{
					fn:     wantFn,
					format: wantFormat,
					vals:   wantVals,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						OutfFunc: func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) {
								t.Errorf("vals got = %v, want = %v", vals, wantVals)
							}
							if format != wantFormat {
								t.Errorf("format got = %s, want = %s", format, wantFormat)
							}
							if reflect.ValueOf(fn).Pointer() != reflect.ValueOf(wantFn).Pointer() {
								t.Errorf("fn got = %p, want = %p", fn, wantFn)
							}
							cnt++
						},
					}
				},
				checkFunc: func(w want) error {
					if cnt != 1 {
						return errors.Errorf("cnt got = %d, want = %d", cnt, 1)
					}
					return nil
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
			fields := test.fieldsFunc(tt)
			r := &Retry{
				OutfFunc: fields.OutfFunc,
			}

			r.Outf(test.args.fn, test.args.format, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
