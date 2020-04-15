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

package log

import (
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/glg"
	"github.com/vdaas/vald/internal/log/level"
	logger "github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/log/mock"
)

func TestInit(t *testing.T) {
	type test struct {
		name      string
		opts      []Option
		checkFunc func(Logger) error
	}

	tests := []test{
		func() test {
			l := glg.New()
			return test{
				name: "set l object when option is not empty",
				opts: []Option{
					WithLogger(l),
				},
				checkFunc: func(got Logger) error {
					if !reflect.DeepEqual(got, l) {
						return errors.Errorf("not equals. want: %v, but got: %v", got, l)
					}
					return nil
				},
			}
		}(),

		func() test {
			return test{
				name: "set l object when option is empty",
				opts: []Option{},
				checkFunc: func(got Logger) error {
					if got == nil {
						return errors.New("l is nil")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				l = nil
				once = sync.Once{}
			}()

			Init(tt.opts...)

			if err := tt.checkFunc(l); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGetLogger(t *testing.T) {
	type test struct {
		name string
		o    *option
		want Logger
	}

	tests := []test{
		{
			name: "returns glg object when l type is GLG",
			o: &option{
				logType: logger.GLG,
				level:   level.DEBUG,
				format:  format.JSON,
			},
			want: glg.New(
				glg.WithLevel(level.DEBUG.String()),
				glg.WithFormat(format.JSON.String()),
			),
		},

		func() test {
			l := glg.New()

			return test{
				name: "returns l when l type is Unknown",
				o: &option{
					logType: logger.Unknown,
					logger:  l,
				},
				want: l,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getLogger(tt.o)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}

func TestBold(t *testing.T) {
	type test struct {
		name string
		str  string
		want string
	}

	tests := []test{
		{
			name: "returns concat string with bash sequence",
			str:  "Vald",
			want: "\033[1mVald\033[22m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Bold(tt.str)
			if tt.want != got {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	type args struct {
		vals []interface{}
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var want []interface{}

			l := &mock.Logger{
				DebugFunc: func(vals ...interface{}) {
					want = vals
				},
			}

			vals := []interface{}{
				"vald",
			}

			return test{
				name: "output success",
				args: args{
					vals: vals,
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if !reflect.DeepEqual(vals, want) {
						return errors.Errorf("not equals. want: %v, got: %v", want, vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l = tt.global.l

			Debug(tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDebugf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var (
				wantFormat string
				wantVals   []interface{}
			)

			l := &mock.Logger{
				DebugfFunc: func(format string, vals ...interface{}) {
					wantFormat = format
					wantVals = vals
				},
			}

			format := "%v"
			vals := []interface{}{
				"vald",
			}

			return test{
				name: "output success",
				args: args{
					format: format,
					vals:   vals,
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if !reflect.DeepEqual(format, wantFormat) {
						return errors.Errorf("vals is not equals. want: %v, got: %v", wantFormat, format)
					}

					if !reflect.DeepEqual(vals, wantVals) {
						return errors.Errorf("vals is not equals. want: %v, got: %v", wantVals, vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l = tt.global.l

			Debugf(tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		vals []interface{}
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var want []interface{}

			l := &mock.Logger{
				InfoFunc: func(vals ...interface{}) {
					want = vals
				},
			}

			vals := []interface{}{
				"vald",
			}

			return test{
				name: "output success",
				args: args{
					vals: vals,
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if !reflect.DeepEqual(vals, want) {
						return errors.Errorf("not equals. want: %v, got: %v", want, vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l = tt.global.l

			Info(tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestInfof(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var (
				wantFormat string
				wantVals   []interface{}
			)

			l := &mock.Logger{
				InfofFunc: func(format string, vals ...interface{}) {
					wantFormat = format
					wantVals = vals
				},
			}

			format := "%v"
			vals := []interface{}{
				"vald",
			}

			return test{
				name: "output success",
				args: args{
					format: format,
					vals:   vals,
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if !reflect.DeepEqual(format, wantFormat) {
						return errors.Errorf("vals is not equals. want: %v, got: %v", wantFormat, format)
					}

					if !reflect.DeepEqual(vals, wantVals) {
						return errors.Errorf("vals is not equals. want: %v, got: %v", wantVals, vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l = tt.global.l

			Infof(tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		vals []interface{}
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var want []interface{}

			l := &mock.Logger{
				WarnFunc: func(vals ...interface{}) {
					want = vals
				},
			}

			vals := []interface{}{
				"vald",
			}

			return test{
				name: "output success",
				args: args{
					vals: vals,
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if !reflect.DeepEqual(vals, want) {
						return errors.Errorf("not equals. want: %v, got: %v", want, vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l = tt.global.l

			Warn(tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWarnf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var (
				wantFormat string
				wantVals   []interface{}
			)

			l := &mock.Logger{
				WarnfFunc: func(format string, vals ...interface{}) {
					wantFormat = format
					wantVals = vals
				},
			}

			format := "%v"
			vals := []interface{}{
				"vald",
			}

			return test{
				name: "output success",
				args: args{
					format: format,
					vals:   vals,
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if !reflect.DeepEqual(format, wantFormat) {
						return errors.Errorf("vals is not equals. want: %v, got: %v", wantFormat, format)
					}

					if !reflect.DeepEqual(vals, wantVals) {
						return errors.Errorf("vals is not equals. want: %v, got: %v", wantVals, vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l = tt.global.l

			Warnf(tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		vals []interface{}
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var want []interface{}

			l := &mock.Logger{
				ErrorFunc: func(vals ...interface{}) {
					want = vals
				},
			}

			vals := []interface{}{
				"vald",
			}

			return test{
				name: "output success",
				args: args{
					vals: vals,
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if !reflect.DeepEqual(vals, want) {
						return errors.Errorf("not equals. want: %v, got: %v", want, vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l = tt.global.l

			Error(tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var (
				wantFormat string
				wantVals   []interface{}
			)

			l := &mock.Logger{
				ErrorfFunc: func(format string, vals ...interface{}) {
					wantFormat = format
					wantVals = vals
				},
			}

			format := "fmt"
			vals := []interface{}{
				"vald",
			}

			return test{
				name: "output success",
				args: args{
					format: format,
					vals:   vals,
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if !reflect.DeepEqual(format, wantFormat) {
						return errors.Errorf("vals is not equals. want: %v, got: %v", wantFormat, format)
					}

					if !reflect.DeepEqual(vals, wantVals) {
						return errors.Errorf("vals is not equals. want: %v, got: %v", wantVals, vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l = tt.global.l

			Errorf(tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestFatal(t *testing.T) {
	type args struct {
		vals []interface{}
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var want []interface{}

			l := &mock.Logger{
				FatalFunc: func(vals ...interface{}) {
					want = vals
				},
			}

			vals := []interface{}{
				"vald",
			}

			return test{
				name: "output success",
				args: args{
					vals: vals,
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if !reflect.DeepEqual(vals, want) {
						return errors.Errorf("not equals. want: %v, got: %v", want, vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l = tt.global.l

			Fatal(tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestFatalf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var (
				wantFormat string
				wantVals   []interface{}
			)

			l := &mock.Logger{
				FatalfFunc: func(format string, vals ...interface{}) {
					wantFormat = format
					wantVals = vals
				},
			}

			format := "%v"
			vals := []interface{}{
				"vald",
			}

			return test{
				name: "output success",
				args: args{
					format: format,
					vals:   vals,
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if !reflect.DeepEqual(format, wantFormat) {
						return errors.Errorf("vals is not equals. want: %v, got: %v", wantFormat, format)
					}

					if !reflect.DeepEqual(vals, wantVals) {
						return errors.Errorf("vals is not equals. want: %v, got: %v", wantVals, vals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l = tt.global.l

			Fatalf(tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}
