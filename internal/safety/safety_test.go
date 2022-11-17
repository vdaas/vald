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
package safety

import (
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	info.Init("")
	os.Exit(m.Run())
}

func TestRecoverFunc(t *testing.T) {
	type args struct {
		fn func() error
	}
	type want struct {
		want      func() error
		wantPanic func() error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, func() error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got func() error) error {
		gotErr := got()

		// if wantPanic is not nil then the panic should be recovered and this line should not be executed
		if w.wantPanic != nil {
			return errors.Errorf("wantPanic is not nil, but got return error: %v", gotErr)
		}

		if (w.want == nil && got != nil) || (w.want != nil && got == nil) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w)
		}
		wantErr := w.want()
		if !errors.Is(gotErr, wantErr) {
			return errors.Errorf("got error= %v, want error= %v", gotErr, wantErr)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns error when system panicked caused by runtime error",
			args: args{
				fn: func() error {
					_ = []string{}[10]
					return nil
				},
			},
			want: want{
				wantPanic: func() error {
					return errors.New("system panicked caused by runtime error: runtime error: index out of range [10] with length 0")
				},
			},
		},
		{
			name: "returns error when system panicked caused by panic with string value",
			args: args{
				fn: func() error {
					panic("panic")
				},
			},
			want: want{
				want: func() error {
					return errors.New("panic recovered: panic")
				},
			},
		},
		{
			name: "returns error when system panicked caused by panic with error",
			args: args{
				fn: func() error {
					panic(errors.Errorf("error"))
				},
			},
			want: want{
				want: func() error {
					return errors.Errorf("error")
				},
			},
		},
		{
			name: "returns error when system panicked caused by panic with int value",
			args: args{
				fn: func() error {
					panic(10)
				},
			},
			want: want{
				want: func() error {
					return errors.New("panic recovered: 10")
				},
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)

			defer func(w want, tt *testing.T) {
				gotPanic := recover()
				if w.wantPanic == nil && gotPanic == nil {
					return
				}
				panicErr, ok := gotPanic.(error)
				if !ok {
					tt.Errorf("cannot cast panic to error, panic: %v", gotPanic)
					return
				}
				if want := w.wantPanic(); !errors.Is(want, panicErr) {
					tt.Errorf("want: %v, got: %v", want, panicErr)
				}
			}(test.want, tt)

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

			got := RecoverFunc(test.args.fn)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestRecoverWithoutPanicFunc(t *testing.T) {
	type args struct {
		fn func() error
	}
	type want struct {
		want func() error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, func() error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got func() error) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           fn: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           fn: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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

			got := RecoverWithoutPanicFunc(test.args.fn)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_recoverFunc(t *testing.T) {
	type args struct {
		fn        func() error
		withPanic bool
	}
	type want struct {
		want func() error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, func() error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got func() error) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           fn: nil,
		           withPanic: false,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           fn: nil,
		           withPanic: false,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
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

			got := recoverFunc(test.args.fn, test.args.withPanic)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
