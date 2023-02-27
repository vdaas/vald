// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package glg

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/log/mock"
	"github.com/vdaas/vald/internal/log/retry"
)

func TestNew(t *testing.T) {
	type test struct {
		name string
		opts []Option
		want *logger
	}

	tests := []test{
		func() test {
			glg := glg.New()
			retry := retry.New()

			return test{
				name: "returns logger object when option and defaultOptions is set",
				opts: []Option{
					WithGlg(glg),
					WithRetry(retry),
				},
				want: &logger{
					glg:   glg,
					level: level.DEBUG,
					retry: retry,
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.opts...)
			if !reflect.DeepEqual(tt.want, l) {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, l)
			}
		})
	}
}

func TestSetLevelMode(t *testing.T) {
	type args struct {
		lv level.Level
	}

	type field struct {
		glg *glg.Glg
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(got *logger) error
	}

	tests := []test{
		{
			name: "returns logger object updated the glg object when lv is DEBUG",
			args: args{
				lv: level.DEBUG,
			},
			field: field{
				glg: glg.New(),
			},
			checkFunc: func(got *logger) error {
				g := got.glg

				if g.GetCurrentMode(glg.DEBG) != glg.STD {
					return errors.New("debug level is wrong")
				}

				if g.GetCurrentMode(glg.INFO) != glg.STD {
					return errors.New("info level is wrong")
				}

				if g.GetCurrentMode(glg.WARN) != glg.STD {
					return errors.New("warn level is wrong")
				}

				if g.GetCurrentMode(glg.ERR) != glg.STD {
					return errors.New("error level is wrong")
				}

				if g.GetCurrentMode(glg.FATAL) != glg.STD {
					return errors.New("fatal level is wrong")
				}

				return nil
			},
		},

		{
			name: "returns logger object updated the glg object when lv is INFO",
			args: args{
				lv: level.INFO,
			},
			field: field{
				glg: glg.New(),
			},
			checkFunc: func(got *logger) error {
				g := got.glg

				if g.GetCurrentMode(glg.DEBG) != glg.NONE {
					return errors.New("debug level is wrong")
				}

				if g.GetCurrentMode(glg.INFO) != glg.STD {
					return errors.New("info level is wrong")
				}

				if g.GetCurrentMode(glg.WARN) != glg.STD {
					return errors.New("warn level is wrong")
				}

				if g.GetCurrentMode(glg.ERR) != glg.STD {
					return errors.New("error level is wrong")
				}

				if g.GetCurrentMode(glg.FATAL) != glg.STD {
					return errors.New("fatal level is wrong")
				}

				return nil
			},
		},

		{
			name: "returns logger object updated the glg object when lv is WARN",
			args: args{
				lv: level.WARN,
			},
			field: field{
				glg: glg.New(),
			},
			checkFunc: func(got *logger) error {
				g := got.glg

				if g.GetCurrentMode(glg.DEBG) != glg.NONE {
					return errors.New("debug level is wrong")
				}

				if g.GetCurrentMode(glg.INFO) != glg.NONE {
					return errors.New("info level is wrong")
				}

				if g.GetCurrentMode(glg.WARN) != glg.STD {
					return errors.New("warn level is wrong")
				}

				if g.GetCurrentMode(glg.ERR) != glg.STD {
					return errors.New("error level is wrong")
				}

				if g.GetCurrentMode(glg.FATAL) != glg.STD {
					return errors.New("fatal level is wrong")
				}

				return nil
			},
		},

		{
			name: "returns logger object updated the glg object when lv is ERROR",
			args: args{
				lv: level.ERROR,
			},
			field: field{
				glg: glg.New(),
			},
			checkFunc: func(got *logger) error {
				g := got.glg

				if g.GetCurrentMode(glg.DEBG) != glg.NONE {
					return errors.New("debug level is wrong")
				}

				if g.GetCurrentMode(glg.INFO) != glg.NONE {
					return errors.New("info level is wrong")
				}

				if g.GetCurrentMode(glg.WARN) != glg.NONE {
					return errors.New("warn level is wrong")
				}

				if g.GetCurrentMode(glg.ERR) != glg.STD {
					return errors.New("error level is wrong")
				}

				if g.GetCurrentMode(glg.FATAL) != glg.STD {
					return errors.New("fatal level is wrong")
				}

				return nil
			},
		},

		{
			name: "returns logger object updated the glg object when lv is FATAL",
			args: args{
				lv: level.FATAL,
			},
			field: field{
				glg: glg.New(),
			},
			checkFunc: func(got *logger) error {
				g := got.glg

				if g.GetCurrentMode(glg.DEBG) != glg.NONE {
					return errors.New("debug level is wrong")
				}

				if g.GetCurrentMode(glg.INFO) != glg.NONE {
					return errors.New("info level is wrong")
				}

				if g.GetCurrentMode(glg.WARN) != glg.NONE {
					return errors.New("warn level is wrong")
				}

				if g.GetCurrentMode(glg.ERR) != glg.NONE {
					return errors.New("error level is wrong")
				}

				if g.GetCurrentMode(glg.FATAL) != glg.STD {
					return errors.New("fatal level is wrong")
				}

				return nil
			},
		},

		{
			name: "returns logger object updated the glg object when lv is Unknown",
			args: args{
				lv: level.Unknown,
			},
			field: field{
				glg: glg.New(),
			},
			checkFunc: func(got *logger) error {
				g := got.glg

				if g.GetCurrentMode(glg.DEBG) != glg.NONE {
					return errors.New("debug level is wrong")
				}

				if g.GetCurrentMode(glg.INFO) != glg.NONE {
					return errors.New("info level is wrong")
				}

				if g.GetCurrentMode(glg.WARN) != glg.NONE {
					return errors.New("warn level is wrong")
				}

				if g.GetCurrentMode(glg.ERR) != glg.NONE {
					return errors.New("error level is wrong")
				}

				if g.GetCurrentMode(glg.FATAL) != glg.NONE {
					return errors.New("fatal level is wrong")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := (&logger{
				glg: tt.field.glg,
			}).setLevelMode(tt.args.lv)
			if err := tt.checkFunc(l); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSetLogFormat(t *testing.T) {
	type args struct {
		fmt format.Format
	}

	type field struct {
		glg *glg.Glg
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(got *logger) error
	}

	tests := []test{
		{
			name: "returns logger object updated the glg object when format is JSON",
			args: args{
				fmt: format.JSON,
			},
			field: field{
				glg: glg.New().SetMode(glg.BOTH),
			},
			checkFunc: func(got *logger) error {
				buf := new(bytes.Buffer)
				got.glg.SetLevelWriter(glg.INFO, buf)
				got.glg.Info("vald")

				var obj map[string]interface{}
				if err := json.NewDecoder(buf).Decode(&obj); err != nil {
					return errors.New("not in JSON output logger")
				}
				return nil
			},
		},

		{
			name: "returns logger object without updating the glg object when format is invalid",
			args: args{
				fmt: format.Unknown,
			},
			field: field{
				glg: glg.New().SetMode(glg.BOTH),
			},
			checkFunc: func(got *logger) error {
				buf := new(bytes.Buffer)
				got.glg.AddLevelWriter(glg.INFO, buf)
				got.glg.Info("vald")

				var obj map[string]interface{}
				if err := json.NewDecoder(buf).Decode(&obj); err == nil {
					return errors.New("not in RAW output logger")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := (&logger{
				glg: tt.field.glg,
			}).setLogFormat(tt.args.fmt)

			if err := tt.checkFunc(l); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		vals interface{}
	}

	type field struct {
		glg   *glg.Glg
		retry retry.Retry
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var gotVals string
			retry := &mock.Retry{
				OutFunc: func(fn func(vals ...interface{}) error, vals ...interface{}) {
					gotVals = vals[0].(string)
				},
			}

			wantVals := "vals"

			return test{
				name: "output success",
				args: args{
					vals: wantVals,
				},
				field: field{
					retry: retry,
					glg:   glg.Get(),
				},
				checkFunc: func() error {
					if gotVals != wantVals {
						return errors.Errorf("vals not equals. want: %v, but got: %v", wantVals, gotVals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{
				glg:   tt.field.glg,
				retry: tt.field.retry,
			}
			l.Info(tt.args.vals)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestInfof(t *testing.T) {
	type args struct {
		vals   interface{}
		format string
	}

	type field struct {
		glg   *glg.Glg
		retry retry.Retry
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var (
				gotVals   string
				gotFormat string
			)
			retry := &mock.Retry{
				OutfFunc: func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{}) {
					gotFormat = format
					gotVals = vals[0].(string)
				},
			}

			wantVals := "vals"
			wantFormat := "format"

			return test{
				name: "output success",
				args: args{
					vals:   wantVals,
					format: wantFormat,
				},
				field: field{
					retry: retry,
					glg:   glg.Get(),
				},
				checkFunc: func() error {
					if gotFormat != wantFormat {
						return errors.Errorf("format not equals. want: %v, but got: %v", wantFormat, gotFormat)
					}

					if gotVals != wantVals {
						return errors.Errorf("vals not equals. want: %v, but got: %v", wantVals, gotVals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{
				glg:   tt.field.glg,
				retry: tt.field.retry,
			}
			l.Infof(tt.args.format, tt.args.vals)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDebug(t *testing.T) {
	type args struct {
		vals interface{}
	}

	type field struct {
		glg   *glg.Glg
		retry retry.Retry
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var gotVals string
			retry := &mock.Retry{
				OutFunc: func(fn func(vals ...interface{}) error, vals ...interface{}) {
					gotVals = vals[0].(string)
				},
			}

			wantVals := "vals"

			return test{
				name: "output success",
				args: args{
					vals: wantVals,
				},
				field: field{
					retry: retry,
					glg:   glg.Get(),
				},
				checkFunc: func() error {
					if gotVals != wantVals {
						return errors.Errorf("vals not equals. want: %v, but got: %v", wantVals, gotVals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{
				glg:   tt.field.glg,
				retry: tt.field.retry,
			}
			l.Debug(tt.args.vals)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDebugf(t *testing.T) {
	type args struct {
		vals   interface{}
		format string
	}

	type field struct {
		glg   *glg.Glg
		retry retry.Retry
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var (
				gotVals   string
				gotFormat string
			)
			retry := &mock.Retry{
				OutfFunc: func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{}) {
					gotFormat = format
					gotVals = vals[0].(string)
				},
			}

			wantVals := "vals"
			wantFormat := "format"

			return test{
				name: "output success",
				args: args{
					vals:   wantVals,
					format: wantFormat,
				},
				field: field{
					retry: retry,
					glg:   glg.Get(),
				},
				checkFunc: func() error {
					if gotFormat != wantFormat {
						return errors.Errorf("format not equals. want: %v, but got: %v", wantFormat, gotFormat)
					}

					if gotVals != wantVals {
						return errors.Errorf("vals not equals. want: %v, but got: %v", wantVals, gotVals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{
				glg:   tt.field.glg,
				retry: tt.field.retry,
			}
			l.Debugf(tt.args.format, tt.args.vals)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		vals interface{}
	}

	type field struct {
		glg   *glg.Glg
		retry retry.Retry
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var gotVals string
			retry := &mock.Retry{
				OutFunc: func(fn func(vals ...interface{}) error, vals ...interface{}) {
					gotVals = vals[0].(string)
				},
			}

			wantVals := "vals"

			return test{
				name: "output success",
				args: args{
					vals: wantVals,
				},
				field: field{
					retry: retry,
					glg:   glg.Get(),
				},
				checkFunc: func() error {
					if gotVals != wantVals {
						return errors.Errorf("vals not equals. want: %v, but got: %v", wantVals, gotVals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{
				glg:   tt.field.glg,
				retry: tt.field.retry,
			}
			l.Warn(tt.args.vals)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWarnf(t *testing.T) {
	type args struct {
		vals   interface{}
		format string
	}

	type field struct {
		glg   *glg.Glg
		retry retry.Retry
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var (
				gotVals   string
				gotFormat string
			)
			retry := &mock.Retry{
				OutfFunc: func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{}) {
					gotFormat = format
					gotVals = vals[0].(string)
				},
			}

			wantVals := "vals"
			wantFormat := "format"

			return test{
				name: "output success",
				args: args{
					vals:   wantVals,
					format: wantFormat,
				},
				field: field{
					retry: retry,
					glg:   glg.Get(),
				},
				checkFunc: func() error {
					if gotFormat != wantFormat {
						return errors.Errorf("format not equals. want: %v, but got: %v", wantFormat, gotFormat)
					}

					if gotVals != wantVals {
						return errors.Errorf("vals not equals. want: %v, but got: %v", wantVals, gotVals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{
				glg:   tt.field.glg,
				retry: tt.field.retry,
			}
			l.Warnf(tt.args.format, tt.args.vals)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		vals interface{}
	}

	type field struct {
		glg   *glg.Glg
		retry retry.Retry
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var gotVals string
			retry := &mock.Retry{
				OutFunc: func(fn func(vals ...interface{}) error, vals ...interface{}) {
					gotVals = vals[0].(string)
				},
			}

			wantVals := "vals"

			return test{
				name: "output success",
				args: args{
					vals: wantVals,
				},
				field: field{
					retry: retry,
					glg:   glg.Get(),
				},
				checkFunc: func() error {
					if gotVals != wantVals {
						return errors.Errorf("vals not equals. want: %v, but got: %v", wantVals, gotVals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{
				glg:   tt.field.glg,
				retry: tt.field.retry,
			}
			l.Error(tt.args.vals)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		vals   interface{}
		format string
	}

	type field struct {
		glg   *glg.Glg
		retry retry.Retry
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func() error
	}

	tests := []test{
		func() test {
			var (
				gotVals   string
				gotFormat string
			)
			retry := &mock.Retry{
				OutfFunc: func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{}) {
					gotFormat = format
					gotVals = vals[0].(string)
				},
			}

			wantVals := "vals"
			wantFormat := "format"

			return test{
				name: "output success",
				args: args{
					vals:   wantVals,
					format: wantFormat,
				},
				field: field{
					retry: retry,
					glg:   glg.Get(),
				},
				checkFunc: func() error {
					if gotFormat != wantFormat {
						return errors.Errorf("format not equals. want: %v, but got: %v", wantFormat, gotFormat)
					}

					if gotVals != wantVals {
						return errors.Errorf("vals not equals. want: %v, but got: %v", wantVals, gotVals)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logger{
				glg:   tt.field.glg,
				retry: tt.field.retry,
			}
			l.Errorf(tt.args.format, tt.args.vals)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}
