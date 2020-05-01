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
package mock

import (
	"testing"

	"go.uber.org/goleak"
)

func TestRetry_Out(t *testing.T) {
	type args struct {
		fn   func(vals ...interface{}) error
		vals []interface{}
	}
	type fields struct {
		OutFunc  func(fn func(vals ...interface{}) error, vals ...interface{})
		OutfFunc func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{})
	}
	type want struct {
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           fn: nil,
		           vals: nil,
		       },
		       fields: fields {
		           OutFunc: nil,
		           OutfFunc: nil,
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
		           vals: nil,
		           },
		           fields: fields {
		           OutFunc: nil,
		           OutfFunc: nil,
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
			r := &Retry{
				OutFunc:  test.fields.OutFunc,
				OutfFunc: test.fields.OutfFunc,
			}

			r.Out(test.args.fn, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
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
		OutFunc  func(fn func(vals ...interface{}) error, vals ...interface{})
		OutfFunc func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{})
	}
	type want struct {
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           fn: nil,
		           format: "",
		           vals: nil,
		       },
		       fields: fields {
		           OutFunc: nil,
		           OutfFunc: nil,
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
		           format: "",
		           vals: nil,
		           },
		           fields: fields {
		           OutFunc: nil,
		           OutfFunc: nil,
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
			r := &Retry{
				OutFunc:  test.fields.OutFunc,
				OutfFunc: test.fields.OutfFunc,
			}

			r.Outf(test.args.fn, test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
