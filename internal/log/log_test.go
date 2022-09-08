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

package log

import (
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/glg"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/log/mock"
	"github.com/vdaas/vald/internal/log/nop"
	"github.com/vdaas/vald/internal/test/goleak"
)

// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestInit(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		l logger.Logger
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, logger.Logger) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got logger.Logger) error {
		if !reflect.DeepEqual(got, l) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.l)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "initialize success when option is nil",
				want: want{
					glg.New(
						glg.WithLevel(level.DEBUG.String()),
					),
				},
				beforeFunc: func(args) {
					once = sync.Once{}
				},
			}
		}(),

		func() test {
			return test{
				name: "initialize success when option is not nil",
				args: args{
					opts: []Option{
						WithLevel(level.FATAL.String()),
					},
				},
				want: want{
					glg.New(
						glg.WithLevel(level.FATAL.String()),
					),
				},
				beforeFunc: func(args) {
					once = sync.Once{}
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			Init(test.args.opts...)
			if err := checkFunc(test.want, l); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_getLogger(t *testing.T) {
	type args struct {
		o *option
	}
	type want struct {
		want logger.Logger
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, logger.Logger) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got logger.Logger) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "returns zap object when *option.logType is ZAP",
				args: args{
					o: &option{
						logType: logger.ZAP,
					},
				},
				checkFunc: func(w want, got logger.Logger) error {
					if got == nil {
						return errors.New("got object is empty")
					}

					return nil
				},
			}
		}(),

		{
			name: "returns nop logger when *option.logType is NOP",
			args: args{
				o: &option{
					logType: logger.NOP,
				},
			},
			want: want{
				want: nop.New(),
			},
		},

		{
			name: "returns glg object when *option.logType is GLG",
			args: args{
				o: &option{
					logType: logger.GLG,
				},
			},
			want: want{
				want: glg.New(
					glg.WithLevel(level.Unknown.String()),
				),
			},
			checkFunc: func(w want, got logger.Logger) error {
				if got == nil {
					return errors.New("got object is empty")
				}
				return nil
			},
		},

		{
			name: "returns glg object when *option is empty",
			args: args{
				o: new(option),
			},
			want: want{
				want: glg.New(
					glg.WithLevel(level.Unknown.String()),
				),
			},
			checkFunc: func(w want, got logger.Logger) error {
				if got == nil {
					return errors.New("got object is empty")
				}
				return nil
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := getLogger(test.args.o)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestBold(t *testing.T) {
	type args struct {
		str string
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns concat string with bash sequence when str is `Vald`",
			args: args{
				str: "Vald",
			},
			want: want{
				want: "\033[1mVald\033[22m",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := Bold(test.args.str)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type want struct {
		vals []interface{}
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	tests := []test{
		func() test {
			var got []interface{}

			ml := &mock.Logger{
				DebugFunc: func(vals ...interface{}) {
					got = vals
				},
			}

			w := want{
				vals: []interface{}{
					"vald",
				},
			}

			return test{
				name: "output success",
				args: args{
					vals: w.vals,
				},
				want: w,
				beforeFunc: func(args) {
					l = ml
				},
				afterFunc: func(args) {
					l = nil
				},
				checkFunc: func(w want) error {
					if !reflect.DeepEqual(got, w.vals) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			Debug(test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDebugf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type want struct {
		format string
		vals   []interface{}
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	tests := []test{
		func() test {
			var (
				gotFormat string
				gotVals   []interface{}
			)

			ml := &mock.Logger{
				DebugfFunc: func(format string, vals ...interface{}) {
					gotFormat, gotVals = format, vals
				},
			}

			w := want{
				format: "format",
				vals: []interface{}{
					"vald",
				},
			}

			return test{
				name: "output success",
				args: args{
					format: w.format,
					vals:   w.vals,
				},
				want: w,
				beforeFunc: func(args) {
					l = ml
				},
				afterFunc: func(args) {
					l = nil
				},
				checkFunc: func(w want) error {
					if !reflect.DeepEqual(gotFormat, w.format) {
						return errors.Errorf("format got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotFormat, w.format)
					}
					if !reflect.DeepEqual(gotVals, w.vals) {
						return errors.Errorf("format got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVals, w.vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			Debugf(test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type want struct {
		vals []interface{}
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	tests := []test{
		func() test {
			var got []interface{}

			ml := &mock.Logger{
				InfoFunc: func(vals ...interface{}) {
					got = vals
				},
			}

			w := want{
				vals: []interface{}{
					"vald",
				},
			}

			return test{
				name: "output success",
				args: args{
					vals: w.vals,
				},
				want: w,
				beforeFunc: func(args) {
					l = ml
				},
				afterFunc: func(args) {
					l = nil
				},
				checkFunc: func(want) error {
					if !reflect.DeepEqual(got, w.vals) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			Info(test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestInfof(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type want struct {
		format string
		vals   []interface{}
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	tests := []test{
		func() test {
			var (
				gotFormat string
				gotVals   []interface{}
			)

			ml := &mock.Logger{
				InfofFunc: func(format string, vals ...interface{}) {
					gotFormat, gotVals = format, vals
				},
			}

			w := want{
				format: "format",
				vals: []interface{}{
					"vald",
				},
			}

			return test{
				name: "output success",
				args: args{
					format: w.format,
					vals:   w.vals,
				},
				want: w,
				beforeFunc: func(args) {
					l = ml
				},
				afterFunc: func(args) {
					l = nil
				},
				checkFunc: func(w want) error {
					if !reflect.DeepEqual(gotFormat, w.format) {
						return errors.Errorf("format got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotFormat, w.format)
					}
					if !reflect.DeepEqual(gotVals, w.vals) {
						return errors.Errorf("format got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVals, w.vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			Infof(test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type want struct {
		vals []interface{}
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	tests := []test{
		func() test {
			var got []interface{}

			ml := &mock.Logger{
				WarnFunc: func(vals ...interface{}) {
					got = vals
				},
			}

			w := want{
				vals: []interface{}{
					"vald",
				},
			}

			return test{
				name: "output success",
				args: args{
					vals: w.vals,
				},
				want: w,
				beforeFunc: func(args) {
					l = ml
				},
				afterFunc: func(args) {
					l = nil
				},
				checkFunc: func(want) error {
					if !reflect.DeepEqual(got, w.vals) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			Warn(test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWarnf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type want struct {
		format string
		vals   []interface{}
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	tests := []test{
		func() test {
			var (
				gotFormat string
				gotVals   []interface{}
			)

			ml := &mock.Logger{
				WarnfFunc: func(format string, vals ...interface{}) {
					gotFormat, gotVals = format, vals
				},
			}

			w := want{
				format: "format",
				vals: []interface{}{
					"vald",
				},
			}

			return test{
				name: "output success",
				args: args{
					format: w.format,
					vals:   w.vals,
				},
				want: w,
				beforeFunc: func(args) {
					l = ml
				},
				afterFunc: func(args) {
					l = nil
				},
				checkFunc: func(w want) error {
					if !reflect.DeepEqual(gotFormat, w.format) {
						return errors.Errorf("format got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotFormat, w.format)
					}
					if !reflect.DeepEqual(gotVals, w.vals) {
						return errors.Errorf("format got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVals, w.vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			Warnf(test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type want struct {
		vals []interface{}
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	tests := []test{
		func() test {
			var got []interface{}

			ml := &mock.Logger{
				ErrorFunc: func(vals ...interface{}) {
					got = vals
				},
			}

			w := want{
				vals: []interface{}{
					"vald",
				},
			}

			return test{
				name: "output success",
				args: args{
					vals: w.vals,
				},
				want: w,
				beforeFunc: func(args) {
					l = ml
				},
				afterFunc: func(args) {
					l = nil
				},
				checkFunc: func(w want) error {
					if !reflect.DeepEqual(got, w.vals) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			Error(test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type want struct {
		format string
		vals   []interface{}
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	tests := []test{
		func() test {
			var (
				gotFormat string
				gotVals   []interface{}
			)

			ml := &mock.Logger{
				ErrorfFunc: func(format string, vals ...interface{}) {
					gotFormat, gotVals = format, vals
				},
			}

			w := want{
				format: "format",
				vals: []interface{}{
					"vald",
				},
			}

			return test{
				name: "output success",
				args: args{
					format: w.format,
					vals:   w.vals,
				},
				want: w,
				beforeFunc: func(args) {
					l = ml
				},
				afterFunc: func(args) {
					l = nil
				},
				checkFunc: func(w want) error {
					if !reflect.DeepEqual(gotFormat, w.format) {
						return errors.Errorf("format got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotFormat, w.format)
					}
					if !reflect.DeepEqual(gotVals, w.vals) {
						return errors.Errorf("format got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVals, w.vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			Errorf(test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestFatal(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type want struct {
		vals []interface{}
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	tests := []test{
		func() test {
			var got []interface{}

			ml := &mock.Logger{
				FatalFunc: func(vals ...interface{}) {
					got = vals
				},
			}

			w := want{
				vals: []interface{}{
					"vald",
				},
			}

			return test{
				name: "output success",
				args: args{
					vals: w.vals,
				},
				want: w,
				beforeFunc: func(args) {
					l = ml
				},
				afterFunc: func(args) {
					l = nil
				},
				checkFunc: func(w want) error {
					if !reflect.DeepEqual(got, w.vals) {
						return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			Fatal(test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestFatalf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type want struct {
		format string
		vals   []interface{}
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	tests := []test{
		func() test {
			var (
				gotFormat string
				gotVals   []interface{}
			)

			ml := &mock.Logger{
				FatalfFunc: func(format string, vals ...interface{}) {
					gotFormat, gotVals = format, vals
				},
			}

			w := want{
				format: "format",
				vals: []interface{}{
					"vald",
				},
			}

			return test{
				name: "output success",
				args: args{
					format: w.format,
					vals:   w.vals,
				},
				want: w,
				beforeFunc: func(args) {
					l = ml
				},
				afterFunc: func(args) {
					l = nil
				},
				checkFunc: func(w want) error {
					if !reflect.DeepEqual(gotFormat, w.format) {
						return errors.Errorf("format got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotFormat, w.format)
					}
					if !reflect.DeepEqual(gotVals, w.vals) {
						return errors.Errorf("format got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVals, w.vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}

			Fatalf(test.args.format, test.args.vals...)
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestDebugd(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
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
		           msg: "",
		           details: nil,
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
		           msg: "",
		           details: nil,
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

			Debugd(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestInfod(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
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
		           msg: "",
		           details: nil,
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
		           msg: "",
		           details: nil,
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

			Infod(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestWarnd(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
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
		           msg: "",
		           details: nil,
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
		           msg: "",
		           details: nil,
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

			Warnd(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestErrord(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
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
		           msg: "",
		           details: nil,
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
		           msg: "",
		           details: nil,
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

			Errord(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestFatald(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type want struct{}
	type test struct {
		name       string
		args       args
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
		           msg: "",
		           details: nil,
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
		           msg: "",
		           details: nil,
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

			Fatald(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestClose(t *testing.T) {
	t.Parallel()
	type want struct {
		err error
	}
	type test struct {
		name       string
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			err := Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
