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
package safety

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"go.uber.org/goleak"
)

func TestRecoverFunc(t *testing.T) {
	type test struct {
		name       string
		fn         func() error
		runtimeErr bool
		want       error
	}

	tests := []test{
		{
			name: "returns error when system paniced caused by runtime error",
			fn: func() error {
				_ = []string{}[10]
				return nil
			},
			runtimeErr: true,
			want:       errors.New("system paniced caused by runtime error: runtime error: index out of range [10] with length 0"),
		},

		{
			name: "returns error when system paniced caused by panic with string value",
			fn: func() error {
				panic("panic")
			},
			want: errors.New("panic recovered: panic"),
		},

		{
			name: "returns error when system paniced caused by panic with error",
			fn: func() error {
				panic(errors.Errorf("error"))
			},
			want: errors.New("error"),
		},

		{
			name: "returns error when system paniced caused by panic with int value",
			fn: func() error {
				panic(10)
			},
			want: errors.New("panic recovered: 10"),
		},
	}

	log.Init()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if ok := tt.runtimeErr; ok {
					if want, got := tt.want, recover().(error); !errors.Is(got, want) {
						t.Errorf("not equals. want: %v, got: %v", want, got)
					}
				}
			}()

			got := RecoverFunc(tt.fn)()
			if !errors.Is(got, tt.want) {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
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
			return errors.Errorf("got = %v, want %v", got, w.want)
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

			got := RecoverWithoutPanicFunc(test.args.fn)
			if err := test.checkFunc(test.want, got); err != nil {
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
			return errors.Errorf("got = %v, want %v", got, w.want)
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

			got := recoverFunc(test.args.fn, test.args.withPanic)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
