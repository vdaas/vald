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
package zap

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/test/goleak"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want *logger
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, *logger, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got *logger, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			called := false

			return test{
				name: "return new logger instance correctly",
				args: args{
					opts: []Option{
						func(l *logger) {
							called = true
						},
					},
				},
				want: want{
					want: &logger{},
					err:  nil,
				},
				checkFunc: func(w want, got *logger, err error) error {
					if !called {
						return errors.New("Option function is not applied")
					}

					if got == nil {
						return errors.New("logger is not returned")
					}

					if !errors.Is(err, w.err) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
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

			got, err := New(test.args.opts...)
			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_initialize(t *testing.T) {
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type args struct {
		sinkPath    string
		errSinkPath string
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		args       args
		want       want
		checkFunc  func(want, *logger, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, l *logger, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			_, _, err := zap.Open("")

			return test{
				name: "returns error when sinkPath is invalid",
				fields: fields{
					format: format.RAW,
					level:  level.DEBUG,
				},
				args: args{
					sinkPath:    "",
					errSinkPath: "stderr",
				},
				want: want{
					err: err,
				},
				checkFunc: func(w want, l *logger, err error) error {
					if l.logger != nil {
						return errors.New("l.logger is not empty")
					}

					if l.sugar != nil {
						return errors.New("l.sugar is not empty")
					}

					return defaultCheckFunc(w, l, err)
				},
			}
		}(),
		func() test {
			_, _, err := zap.Open("xxx:///")

			return test{
				name: "returns error when errSinkPath is invalid",
				fields: fields{
					format: format.RAW,
					level:  level.DEBUG,
				},
				args: args{
					sinkPath:    "stdout",
					errSinkPath: "xxx:///",
				},
				want: want{
					err: err,
				},
				checkFunc: func(w want, l *logger, err error) error {
					if l.logger != nil {
						return errors.New("l.logger is not empty")
					}

					if l.sugar != nil {
						return errors.New("l.sugar is not empty")
					}

					return defaultCheckFunc(w, l, err)
				},
			}
		}(),
		{
			name: "returns nil and initialize the logger and sugar fields correctly",
			fields: fields{
				format: format.RAW,
				level:  level.DEBUG,
			},
			args: args{
				sinkPath:    "stdout",
				errSinkPath: "stderr",
			},
			want: want{
				err: nil,
			},
			checkFunc: func(w want, l *logger, err error) error {
				if l.logger == nil {
					return errors.New("l.logger is nil")
				}

				if l.sugar == nil {
					return errors.New("l.sugar is nil")
				}

				return defaultCheckFunc(w, l, err)
			},
		},
		{
			name: "returns nil and initialize the logger and sugar fields correctly (with enableCaller)",
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: true,
			},
			args: args{
				sinkPath:    "stdout",
				errSinkPath: "stderr",
			},
			want: want{
				err: nil,
			},
			checkFunc: func(w want, l *logger, err error) error {
				if l.logger == nil {
					return errors.New("l.logger is nil")
				}

				if l.sugar == nil {
					return errors.New("l.sugar is nil")
				}

				return defaultCheckFunc(w, l, err)
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			err := l.initialize(test.args.sinkPath, test.args.errSinkPath)
			if err := checkFunc(test.want, l, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_toZapLevel(t *testing.T) {
	type args struct {
		lv level.Level
	}
	type want struct {
		want zapcore.Level
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, zapcore.Level) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got zapcore.Level) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns DebugLevel if lv is DEBUG",
			args: args{
				lv: level.DEBUG,
			},
			want: want{
				want: zapcore.DebugLevel,
			},
		},
		{
			name: "returns InfoLevel if lv is INFO",
			args: args{
				lv: level.INFO,
			},
			want: want{
				want: zapcore.InfoLevel,
			},
		},
		{
			name: "returns WarnLevel if lv is WARN",
			args: args{
				lv: level.WARN,
			},
			want: want{
				want: zapcore.WarnLevel,
			},
		},
		{
			name: "returns ErrorLevel if lv is Error",
			args: args{
				lv: level.ERROR,
			},
			want: want{
				want: zapcore.ErrorLevel,
			},
		},
		{
			name: "returns FatalLevel if lv is FATAL",
			args: args{
				lv: level.FATAL,
			},
			want: want{
				want: zapcore.FatalLevel,
			},
		},
		{
			name: "returns defaultLevel if lv is Unknown",
			args: args{
				lv: level.Unknown,
			},
			want: want{
				want: defaultLevel,
			},
		},
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

			got := toZapLevel(test.args.lv)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_toZapEncoder(t *testing.T) {
	type args struct {
		fmt format.Format
	}
	type want struct {
		want zapcore.Encoder
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, zapcore.Encoder) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got zapcore.Encoder) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	consoleEnc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	zapcore_NewConsoleEncoder = func(_ zapcore.EncoderConfig) zapcore.Encoder {
		return consoleEnc
	}

	jsonEnc := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	zapcore_NewJSONEncoder = func(_ zapcore.EncoderConfig) zapcore.Encoder {
		return jsonEnc
	}

	tests := []test{
		{
			name: "returns ConsoleEncoder if fmt is RAW",
			args: args{
				fmt: format.RAW,
			},
			want: want{
				want: consoleEnc,
			},
		},
		{
			name: "returns JSONEncoder if fmt is JSON",
			args: args{
				fmt: format.JSON,
			},
			want: want{
				want: jsonEnc,
			},
		},
		{
			name: "returns JSONEncoder if fmt is Unknown",
			args: args{
				fmt: format.Unknown,
			},
			want: want{
				want: jsonEnc,
			},
		},
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

			got := toZapEncoder(test.args.fmt)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Debug(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Debug",
			args: args{
				vals: []interface{}{"value"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Debug(test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Debugf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Debugf",
			args: args{
				format: "%s",
				vals:   []interface{}{"value"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Debugf(test.args.format, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Info(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Info",
			args: args{
				vals: []interface{}{"value"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Info(test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Infof(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Infof",
			args: args{
				format: "%s",
				vals:   []interface{}{"value"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Infof(test.args.format, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Warn(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Warn",
			args: args{
				vals: []interface{}{"value"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Warn(test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Warnf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Warnf",
			args: args{
				format: "%s",
				vals:   []interface{}{"value"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Warnf(test.args.format, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Error(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Error",
			args: args{
				vals: []interface{}{"value"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Error(test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Errorf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Errorf",
			args: args{
				format: "%s",
				vals:   []interface{}{"value"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Errorf(test.args.format, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Fatal(t *testing.T) {
	type args struct {
		vals []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		// {
		// 	name: "just call Fatal",
		// 	args: args{
		// 		vals: []interface{}{"value"},
		// 	},
		// 	fields: fields{
		// 		format:       format.RAW,
		// 		level:        level.DEBUG,
		// 		enableCaller: false,
		// 		logger:       zap.L(),
		// 		sugar:        zap.L().Sugar(),
		// 	},
		// 	want:      want{},
		// 	checkFunc: defaultCheckFunc,
		// },
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Fatal(test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Fatalf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		// {
		// 	name: "just call Fatalf",
		// 	args: args{
		// 		format: "%s",
		// 		vals:   []interface{}{"value"},
		// 	},
		// 	fields: fields{
		// 		format:       format.RAW,
		// 		level:        level.DEBUG,
		// 		enableCaller: false,
		// 		logger:       zap.L(),
		// 		sugar:        zap.L().Sugar(),
		// 	},
		// 	want:      want{},
		// 	checkFunc: defaultCheckFunc,
		// },
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Fatalf(test.args.format, test.args.vals...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Debugd(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Debugd",
			args: args{
				msg:     "message",
				details: []interface{}{"detail"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "call Debugd with two details",
			args: args{
				msg: "message",
				details: []interface{}{
					"detail1",
					"detail2",
				},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Debugd(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Infod(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Infod",
			args: args{
				msg:     "message",
				details: []interface{}{"detail"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "call Infod with two details",
			args: args{
				msg: "message",
				details: []interface{}{
					"detail1",
					"detail2",
				},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Infod(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Warnd(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Warnd",
			args: args{
				msg:     "message",
				details: []interface{}{"detail"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "call Warnd with two details",
			args: args{
				msg: "message",
				details: []interface{}{
					"detail1",
					"detail2",
				},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Warnd(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Errord(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		{
			name: "just call Errord",
			args: args{
				msg:     "message",
				details: []interface{}{"detail"},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "call Errord with two details",
			args: args{
				msg: "message",
				details: []interface{}{
					"detail1",
					"detail2",
				},
			},
			fields: fields{
				format:       format.RAW,
				level:        level.DEBUG,
				enableCaller: false,
				logger:       zap.L(),
				sugar:        zap.L().Sugar(),
			},
			want:      want{},
			checkFunc: defaultCheckFunc,
		},
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Errord(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Fatald(t *testing.T) {
	t.Parallel()
	type args struct {
		msg     string
		details []interface{}
	}
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct{}
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
		// {
		// 	name: "just call Fatald",
		// 	args: args{
		// 		msg:     "message",
		// 		details: []interface{}{"detail"},
		// 	},
		// 	fields: fields{
		// 		format:       format.RAW,
		// 		level:        level.DEBUG,
		// 		enableCaller: false,
		// 		logger:       zap.L(),
		// 		sugar:        zap.L().Sugar(),
		// 	},
		// 	want:      want{},
		// 	checkFunc: defaultCheckFunc,
		// },
		// {
		// 	name: "call Fatald with two details",
		// 	args: args{
		// 		msg: "message",
		// 		details: []interface{}{
		// 			"detail1",
		// 			"detail2",
		// 		},
		// 	},
		// 	fields: fields{
		// 		format:       format.RAW,
		// 		level:        level.DEBUG,
		// 		enableCaller: false,
		// 		logger:       zap.L(),
		// 		sugar:        zap.L().Sugar(),
		// 	},
		// 	want:      want{},
		// 	checkFunc: defaultCheckFunc,
		// },
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
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			l.Fatald(test.args.msg, test.args.details...)
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_logger_Close(t *testing.T) {
	t.Parallel()
	type fields struct {
		format       format.Format
		level        level.Level
		enableCaller bool
		logger       *zap.Logger
		sugar        *zap.SugaredLogger
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
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
		       fields: fields {
		           format: nil,
		           level: nil,
		           enableCaller: false,
		           logger: nil,
		           sugar: nil,
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
		           fields: fields {
		           format: nil,
		           level: nil,
		           enableCaller: false,
		           logger: nil,
		           sugar: nil,
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
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			l := &logger{
				format:       test.fields.format,
				level:        test.fields.level,
				enableCaller: test.fields.enableCaller,
				logger:       test.fields.logger,
				sugar:        test.fields.sugar,
			}

			err := l.Close()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
