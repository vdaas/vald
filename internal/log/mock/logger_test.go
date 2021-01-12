//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestLogger_Debug(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type fields struct {
		DebugFunc func(vals ...interface{})
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
			wantVals := []interface{}{
				"Vald",
			}
			var cnt int
			return test{
				name: "Call DebugFunc",
				args: args{
					vals: wantVals,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						DebugFunc: func(vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) {
								t.Errorf("got = %v, want = %v", vals, wantVals)
							}
							cnt++
						},
					}
				},
				checkFunc: func(want) error {
					if cnt != 1 {
						return errors.Errorf("got cnt = %d, want = %d", cnt, 1)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			fields := test.fieldsFunc(tt)
			l := &Logger{
				DebugFunc: fields.DebugFunc,
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
		DebugfFunc func(format string, vals ...interface{})
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
			wantVals := []interface{}{
				"Vald",
			}
			wantFormat := "json"
			var cnt int
			return test{
				name: "Call DebugfFunc",
				args: args{
					vals:   wantVals,
					format: wantFormat,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						DebugfFunc: func(format string, vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) || !reflect.DeepEqual(format, wantFormat) {
								t.Errorf("got = %v, want = %v", vals, wantVals)
							}
							cnt++
						},
					}
				},
				checkFunc: func(want) error {
					if cnt != 1 {
						return errors.Errorf("got cnt = %d, want = %d", cnt, 1)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			fields := test.fieldsFunc(tt)
			l := &Logger{
				DebugfFunc: fields.DebugfFunc,
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
		InfoFunc func(vals ...interface{})
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
			wantVals := []interface{}{
				"Vald",
			}
			var cnt int
			return test{
				name: "Call InfoFunc",
				args: args{
					vals: wantVals,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						InfoFunc: func(vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) {
								t.Errorf("got = %v, want = %v", vals, wantVals)
							}
							cnt++
						},
					}
				},
				checkFunc: func(want) error {
					if cnt != 1 {
						return errors.Errorf("got cnt = %d, want = %d", cnt, 1)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			fields := test.fieldsFunc(tt)
			l := &Logger{
				InfoFunc: fields.InfoFunc,
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
		InfofFunc func(format string, vals ...interface{})
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
			wantVals := []interface{}{
				"Vald",
			}
			wantFormat := "json"
			var cnt int
			return test{
				name: "Call InfofFunc",
				args: args{
					vals:   wantVals,
					format: wantFormat,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						InfofFunc: func(format string, vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) || !reflect.DeepEqual(format, wantFormat) {
								t.Errorf("got = %v, want = %v", vals, wantVals)
							}
							cnt++
						},
					}
				},
				checkFunc: func(want) error {
					if cnt != 1 {
						return errors.Errorf("got cnt = %d, want = %d", cnt, 1)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			fields := test.fieldsFunc(t)
			l := &Logger{
				InfofFunc: fields.InfofFunc,
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
		WarnFunc func(vals ...interface{})
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
			wantVals := []interface{}{
				"Vald",
			}
			var cnt int
			return test{
				name: "Call WarnFunc",
				args: args{
					vals: wantVals,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						WarnFunc: func(vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) {
								t.Errorf("got = %v, want = %v", vals, wantVals)
							}
							cnt++
						},
					}
				},
				checkFunc: func(want) error {
					if cnt != 1 {
						return errors.Errorf("got cnt = %d, want = %d", cnt, 1)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			fields := test.fieldsFunc(tt)
			l := &Logger{
				WarnFunc: fields.WarnFunc,
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
		WarnfFunc func(format string, vals ...interface{})
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
			wantVals := []interface{}{
				"Vald",
			}
			wantFormat := "json"
			var cnt int
			return test{
				name: "Call WarnfFunc",
				args: args{
					vals:   wantVals,
					format: wantFormat,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						WarnfFunc: func(format string, vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) || !reflect.DeepEqual(format, wantFormat) {
								t.Errorf("got = %v, want = %v", vals, wantVals)
							}
							cnt++
						},
					}
				},
				checkFunc: func(want) error {
					if cnt != 1 {
						return errors.Errorf("got cnt = %d, want = %d", cnt, 1)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			fields := test.fieldsFunc(tt)
			l := &Logger{
				WarnfFunc: fields.WarnfFunc,
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
		ErrorFunc func(vals ...interface{})
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
			wantVals := []interface{}{
				"Vald",
			}
			var cnt int
			return test{
				name: "Call ErrorFunc",
				args: args{
					vals: wantVals,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						ErrorFunc: func(vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) {
								t.Errorf("got = %v, want = %v", vals, wantVals)
							}
							cnt++
						},
					}
				},
				checkFunc: func(want) error {
					if cnt != 1 {
						return errors.Errorf("got cnt = %d, want = %d", cnt, 1)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			fields := test.fieldsFunc(tt)
			l := &Logger{
				ErrorFunc: fields.ErrorFunc,
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
		ErrorfFunc func(format string, vals ...interface{})
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
			wantVals := []interface{}{
				"Vald",
			}
			wantFormat := "json"
			var cnt int
			return test{
				name: "Call ErrorfFunc",
				args: args{
					vals:   wantVals,
					format: wantFormat,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						ErrorfFunc: func(format string, vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) || !reflect.DeepEqual(format, wantFormat) {
								t.Errorf("got = %v, want = %v", vals, wantVals)
							}
							cnt++
						},
					}
				},
				checkFunc: func(want) error {
					if cnt != 1 {
						return errors.Errorf("got cnt = %d, want = %d", cnt, 1)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			fields := test.fieldsFunc(tt)
			l := &Logger{
				ErrorfFunc: fields.ErrorfFunc,
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
		FatalFunc func(vals ...interface{})
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
			wantVals := []interface{}{
				"Vald",
			}
			var cnt int
			return test{
				name: "Call FatalFunc",
				args: args{
					vals: wantVals,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						FatalFunc: func(vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) {
								t.Errorf("got = %v, want = %v", vals, wantVals)
							}
							cnt++
						},
					}
				},
				checkFunc: func(want) error {
					if cnt != 1 {
						return errors.Errorf("got cnt = %d, want = %d", cnt, 1)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			fields := test.fieldsFunc(tt)
			l := &Logger{
				FatalFunc: fields.FatalFunc,
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
		FatalfFunc func(format string, vals ...interface{})
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
			wantVals := []interface{}{
				"Vald",
			}
			wantFormat := "json"
			var cnt int
			return test{
				name: "Call FatalfFunc",
				args: args{
					vals:   wantVals,
					format: wantFormat,
				},
				fieldsFunc: func(t *testing.T) fields {
					t.Helper()
					return fields{
						FatalfFunc: func(format string, vals ...interface{}) {
							if !reflect.DeepEqual(vals, wantVals) || !reflect.DeepEqual(format, wantFormat) {
								t.Errorf("got = %v, want = %v", vals, wantVals)
							}
							cnt++
						},
					}
				},
				checkFunc: func(want) error {
					if cnt != 1 {
						return errors.Errorf("got cnt = %d, want = %d", cnt, 1)
					}
					return nil
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			fields := test.fieldsFunc(tt)
			l := &Logger{
				FatalfFunc: fields.FatalfFunc,
			}

			l.Fatalf(test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
