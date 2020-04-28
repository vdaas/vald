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

func TestLogger_Debug(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type fields struct {
		DebugFunc  func(vals ...interface{})
		DebugfFunc func(format string, vals ...interface{})
		InfoFunc   func(vals ...interface{})
		InfofFunc  func(format string, vals ...interface{})
		WarnFunc   func(vals ...interface{})
		WarnfFunc  func(format string, vals ...interface{})
		ErrorFunc  func(vals ...interface{})
		ErrorfFunc func(format string, vals ...interface{})
		FatalFunc  func(vals ...interface{})
		FatalfFunc func(format string, vals ...interface{})
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
		           vals: nil,
		       },
		       fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
		           vals: nil,
		           },
		           fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
			l := &Logger{
				DebugFunc:  test.fields.DebugFunc,
				DebugfFunc: test.fields.DebugfFunc,
				InfoFunc:   test.fields.InfoFunc,
				InfofFunc:  test.fields.InfofFunc,
				WarnFunc:   test.fields.WarnFunc,
				WarnfFunc:  test.fields.WarnfFunc,
				ErrorFunc:  test.fields.ErrorFunc,
				ErrorfFunc: test.fields.ErrorfFunc,
				FatalFunc:  test.fields.FatalFunc,
				FatalfFunc: test.fields.FatalfFunc,
			}

			l.Debug(test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLogger_Debugf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type fields struct {
		DebugFunc  func(vals ...interface{})
		DebugfFunc func(format string, vals ...interface{})
		InfoFunc   func(vals ...interface{})
		InfofFunc  func(format string, vals ...interface{})
		WarnFunc   func(vals ...interface{})
		WarnfFunc  func(format string, vals ...interface{})
		ErrorFunc  func(vals ...interface{})
		ErrorfFunc func(format string, vals ...interface{})
		FatalFunc  func(vals ...interface{})
		FatalfFunc func(format string, vals ...interface{})
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
		           format: "",
		           vals: nil,
		       },
		       fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
		           format: "",
		           vals: nil,
		           },
		           fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
			l := &Logger{
				DebugFunc:  test.fields.DebugFunc,
				DebugfFunc: test.fields.DebugfFunc,
				InfoFunc:   test.fields.InfoFunc,
				InfofFunc:  test.fields.InfofFunc,
				WarnFunc:   test.fields.WarnFunc,
				WarnfFunc:  test.fields.WarnfFunc,
				ErrorFunc:  test.fields.ErrorFunc,
				ErrorfFunc: test.fields.ErrorfFunc,
				FatalFunc:  test.fields.FatalFunc,
				FatalfFunc: test.fields.FatalfFunc,
			}

			l.Debugf(test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLogger_Info(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type fields struct {
		DebugFunc  func(vals ...interface{})
		DebugfFunc func(format string, vals ...interface{})
		InfoFunc   func(vals ...interface{})
		InfofFunc  func(format string, vals ...interface{})
		WarnFunc   func(vals ...interface{})
		WarnfFunc  func(format string, vals ...interface{})
		ErrorFunc  func(vals ...interface{})
		ErrorfFunc func(format string, vals ...interface{})
		FatalFunc  func(vals ...interface{})
		FatalfFunc func(format string, vals ...interface{})
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
		           vals: nil,
		       },
		       fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
		           vals: nil,
		           },
		           fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
			l := &Logger{
				DebugFunc:  test.fields.DebugFunc,
				DebugfFunc: test.fields.DebugfFunc,
				InfoFunc:   test.fields.InfoFunc,
				InfofFunc:  test.fields.InfofFunc,
				WarnFunc:   test.fields.WarnFunc,
				WarnfFunc:  test.fields.WarnfFunc,
				ErrorFunc:  test.fields.ErrorFunc,
				ErrorfFunc: test.fields.ErrorfFunc,
				FatalFunc:  test.fields.FatalFunc,
				FatalfFunc: test.fields.FatalfFunc,
			}

			l.Info(test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLogger_Infof(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type fields struct {
		DebugFunc  func(vals ...interface{})
		DebugfFunc func(format string, vals ...interface{})
		InfoFunc   func(vals ...interface{})
		InfofFunc  func(format string, vals ...interface{})
		WarnFunc   func(vals ...interface{})
		WarnfFunc  func(format string, vals ...interface{})
		ErrorFunc  func(vals ...interface{})
		ErrorfFunc func(format string, vals ...interface{})
		FatalFunc  func(vals ...interface{})
		FatalfFunc func(format string, vals ...interface{})
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
		           format: "",
		           vals: nil,
		       },
		       fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
		           format: "",
		           vals: nil,
		           },
		           fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
			l := &Logger{
				DebugFunc:  test.fields.DebugFunc,
				DebugfFunc: test.fields.DebugfFunc,
				InfoFunc:   test.fields.InfoFunc,
				InfofFunc:  test.fields.InfofFunc,
				WarnFunc:   test.fields.WarnFunc,
				WarnfFunc:  test.fields.WarnfFunc,
				ErrorFunc:  test.fields.ErrorFunc,
				ErrorfFunc: test.fields.ErrorfFunc,
				FatalFunc:  test.fields.FatalFunc,
				FatalfFunc: test.fields.FatalfFunc,
			}

			l.Infof(test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLogger_Warn(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type fields struct {
		DebugFunc  func(vals ...interface{})
		DebugfFunc func(format string, vals ...interface{})
		InfoFunc   func(vals ...interface{})
		InfofFunc  func(format string, vals ...interface{})
		WarnFunc   func(vals ...interface{})
		WarnfFunc  func(format string, vals ...interface{})
		ErrorFunc  func(vals ...interface{})
		ErrorfFunc func(format string, vals ...interface{})
		FatalFunc  func(vals ...interface{})
		FatalfFunc func(format string, vals ...interface{})
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
		           vals: nil,
		       },
		       fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
		           vals: nil,
		           },
		           fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
			l := &Logger{
				DebugFunc:  test.fields.DebugFunc,
				DebugfFunc: test.fields.DebugfFunc,
				InfoFunc:   test.fields.InfoFunc,
				InfofFunc:  test.fields.InfofFunc,
				WarnFunc:   test.fields.WarnFunc,
				WarnfFunc:  test.fields.WarnfFunc,
				ErrorFunc:  test.fields.ErrorFunc,
				ErrorfFunc: test.fields.ErrorfFunc,
				FatalFunc:  test.fields.FatalFunc,
				FatalfFunc: test.fields.FatalfFunc,
			}

			l.Warn(test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLogger_Warnf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type fields struct {
		DebugFunc  func(vals ...interface{})
		DebugfFunc func(format string, vals ...interface{})
		InfoFunc   func(vals ...interface{})
		InfofFunc  func(format string, vals ...interface{})
		WarnFunc   func(vals ...interface{})
		WarnfFunc  func(format string, vals ...interface{})
		ErrorFunc  func(vals ...interface{})
		ErrorfFunc func(format string, vals ...interface{})
		FatalFunc  func(vals ...interface{})
		FatalfFunc func(format string, vals ...interface{})
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
		           format: "",
		           vals: nil,
		       },
		       fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
		           format: "",
		           vals: nil,
		           },
		           fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
			l := &Logger{
				DebugFunc:  test.fields.DebugFunc,
				DebugfFunc: test.fields.DebugfFunc,
				InfoFunc:   test.fields.InfoFunc,
				InfofFunc:  test.fields.InfofFunc,
				WarnFunc:   test.fields.WarnFunc,
				WarnfFunc:  test.fields.WarnfFunc,
				ErrorFunc:  test.fields.ErrorFunc,
				ErrorfFunc: test.fields.ErrorfFunc,
				FatalFunc:  test.fields.FatalFunc,
				FatalfFunc: test.fields.FatalfFunc,
			}

			l.Warnf(test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLogger_Error(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type fields struct {
		DebugFunc  func(vals ...interface{})
		DebugfFunc func(format string, vals ...interface{})
		InfoFunc   func(vals ...interface{})
		InfofFunc  func(format string, vals ...interface{})
		WarnFunc   func(vals ...interface{})
		WarnfFunc  func(format string, vals ...interface{})
		ErrorFunc  func(vals ...interface{})
		ErrorfFunc func(format string, vals ...interface{})
		FatalFunc  func(vals ...interface{})
		FatalfFunc func(format string, vals ...interface{})
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
		           vals: nil,
		       },
		       fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
		           vals: nil,
		           },
		           fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
			l := &Logger{
				DebugFunc:  test.fields.DebugFunc,
				DebugfFunc: test.fields.DebugfFunc,
				InfoFunc:   test.fields.InfoFunc,
				InfofFunc:  test.fields.InfofFunc,
				WarnFunc:   test.fields.WarnFunc,
				WarnfFunc:  test.fields.WarnfFunc,
				ErrorFunc:  test.fields.ErrorFunc,
				ErrorfFunc: test.fields.ErrorfFunc,
				FatalFunc:  test.fields.FatalFunc,
				FatalfFunc: test.fields.FatalfFunc,
			}

			l.Error(test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLogger_Errorf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type fields struct {
		DebugFunc  func(vals ...interface{})
		DebugfFunc func(format string, vals ...interface{})
		InfoFunc   func(vals ...interface{})
		InfofFunc  func(format string, vals ...interface{})
		WarnFunc   func(vals ...interface{})
		WarnfFunc  func(format string, vals ...interface{})
		ErrorFunc  func(vals ...interface{})
		ErrorfFunc func(format string, vals ...interface{})
		FatalFunc  func(vals ...interface{})
		FatalfFunc func(format string, vals ...interface{})
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
		           format: "",
		           vals: nil,
		       },
		       fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
		           format: "",
		           vals: nil,
		           },
		           fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
			l := &Logger{
				DebugFunc:  test.fields.DebugFunc,
				DebugfFunc: test.fields.DebugfFunc,
				InfoFunc:   test.fields.InfoFunc,
				InfofFunc:  test.fields.InfofFunc,
				WarnFunc:   test.fields.WarnFunc,
				WarnfFunc:  test.fields.WarnfFunc,
				ErrorFunc:  test.fields.ErrorFunc,
				ErrorfFunc: test.fields.ErrorfFunc,
				FatalFunc:  test.fields.FatalFunc,
				FatalfFunc: test.fields.FatalfFunc,
			}

			l.Errorf(test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLogger_Fatal(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type fields struct {
		DebugFunc  func(vals ...interface{})
		DebugfFunc func(format string, vals ...interface{})
		InfoFunc   func(vals ...interface{})
		InfofFunc  func(format string, vals ...interface{})
		WarnFunc   func(vals ...interface{})
		WarnfFunc  func(format string, vals ...interface{})
		ErrorFunc  func(vals ...interface{})
		ErrorfFunc func(format string, vals ...interface{})
		FatalFunc  func(vals ...interface{})
		FatalfFunc func(format string, vals ...interface{})
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
		           vals: nil,
		       },
		       fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
		           vals: nil,
		           },
		           fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
			l := &Logger{
				DebugFunc:  test.fields.DebugFunc,
				DebugfFunc: test.fields.DebugfFunc,
				InfoFunc:   test.fields.InfoFunc,
				InfofFunc:  test.fields.InfofFunc,
				WarnFunc:   test.fields.WarnFunc,
				WarnfFunc:  test.fields.WarnfFunc,
				ErrorFunc:  test.fields.ErrorFunc,
				ErrorfFunc: test.fields.ErrorfFunc,
				FatalFunc:  test.fields.FatalFunc,
				FatalfFunc: test.fields.FatalfFunc,
			}

			l.Fatal(test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLogger_Fatalf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type fields struct {
		DebugFunc  func(vals ...interface{})
		DebugfFunc func(format string, vals ...interface{})
		InfoFunc   func(vals ...interface{})
		InfofFunc  func(format string, vals ...interface{})
		WarnFunc   func(vals ...interface{})
		WarnfFunc  func(format string, vals ...interface{})
		ErrorFunc  func(vals ...interface{})
		ErrorfFunc func(format string, vals ...interface{})
		FatalFunc  func(vals ...interface{})
		FatalfFunc func(format string, vals ...interface{})
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
		           format: "",
		           vals: nil,
		       },
		       fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
		           format: "",
		           vals: nil,
		           },
		           fields: fields {
		           DebugFunc: nil,
		           DebugfFunc: nil,
		           InfoFunc: nil,
		           InfofFunc: nil,
		           WarnFunc: nil,
		           WarnfFunc: nil,
		           ErrorFunc: nil,
		           ErrorfFunc: nil,
		           FatalFunc: nil,
		           FatalfFunc: nil,
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
			l := &Logger{
				DebugFunc:  test.fields.DebugFunc,
				DebugfFunc: test.fields.DebugfFunc,
				InfoFunc:   test.fields.InfoFunc,
				InfofFunc:  test.fields.InfofFunc,
				WarnFunc:   test.fields.WarnFunc,
				WarnfFunc:  test.fields.WarnfFunc,
				ErrorFunc:  test.fields.ErrorFunc,
				ErrorfFunc: test.fields.ErrorfFunc,
				FatalFunc:  test.fields.FatalFunc,
				FatalfFunc: test.fields.FatalfFunc,
			}

			l.Fatalf(test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
