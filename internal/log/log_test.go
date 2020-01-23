package log

import (
	"reflect"
	"sync"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/mock"
)

func TestInit(t *testing.T) {
	type test struct {
		name       string
		beforeFunc func()
		opts       []Option
		checkFunc  func() error
	}

	beforeFunc := func() {
		logger = nil
		once = sync.Once{}
	}

	tests := []test{
		func() test {
			mlogger := new(mock.Logger)

			return test{
				name:       "set success when options is not empty",
				beforeFunc: beforeFunc,
				opts: []Option{
					WithLogger(mlogger),
				},
				checkFunc: func() error {
					if !reflect.DeepEqual(logger, mlogger) {
						return errors.Errorf("invalid object was set. want: %v, got: %v", mlogger, logger)
					}
					return nil
				},
			}
		}(),

		{
			name:       "set success when options is empty",
			beforeFunc: beforeFunc,
			checkFunc: func() error {
				if logger == nil {
					return errors.New("logger is nil")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.beforeFunc != nil {
				tt.beforeFunc()
			}
			Init(tt.opts...)

			if err := tt.checkFunc(); err != nil {
				t.Error(err)
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
			want: "\033[1mVald\033[0m",
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
				name: "set success",
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
			logger = tt.global.l

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
				name: "set success",
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
			logger = tt.global.l

			Infof(tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
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
				name: "set success",
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
			logger = tt.global.l

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
				name: "set success",
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
			logger = tt.global.l

			Debugf(tt.args.format, tt.args.vals...)
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
				name: "set success",
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
			logger = tt.global.l

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
				name: "set success",
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
			logger = tt.global.l

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
				name: "set success",
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
			logger = tt.global.l

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
				name: "set success",
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
			logger = tt.global.l

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
				name: "set success",
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
			logger = tt.global.l

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
				name: "set success",
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
			logger = tt.global.l

			Fatalf(tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestOut(t *testing.T) {
	type args struct {
		fn   func(...interface{}) error
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
		recovery  bool
	}

	tests := []test{
		func() test {
			var cnt int
			fn := func(vals ...interface{}) error {
				cnt++
				return nil
			}

			return test{
				name: "processing is successes when fn return nil",
				args: args{
					vals: []interface{}{
						"name",
					},
					fn: fn,
				},
				checkFunc: func() error {
					if cnt != 1 {
						return errors.Errorf("called cnt is wrong. want: %v, got: %v", 1, cnt)
					}
					return nil
				},
				recovery: false,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.recovery {
					if err := recover(); err == nil {
						t.Error("panic is nil")
					}
				}
			}()

			logger = tt.global.l

			retryOut(tt.args.fn, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestOutf(t *testing.T) {
	type args struct {
		fn     func(string, ...interface{}) error
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
		recovery  bool
	}

	tests := []test{
		func() test {
			var cnt int
			fn := func(format string, vals ...interface{}) error {
				cnt++
				return nil
			}

			return test{
				name: "processing is successes when fn return nil",
				args: args{
					format: "format",
					vals: []interface{}{
						"name",
					},
					fn: fn,
				},
				checkFunc: func() error {
					if cnt != 1 {
						return errors.Errorf("called cnt is wrong. want: %v, got: %v", 1, cnt)
					}
					return nil
				},
				recovery: false,
			}
		}(),

		func() test {
			fnErr := errors.New("fail")

			var cnt int
			fn := func(format string, vals ...interface{}) error {
				cnt++
				return fnErr
			}

			var (
				warn error
				err  error
			)

			l := &mock.Logger{
				WarnFunc: func(vals ...interface{}) {
					warn = vals[0].(error)
				},
				ErrorFunc: func(vals ...interface{}) {
					err = vals[0].(error)
				},
			}

			return test{
				name: "processing is fails when fn return error",
				args: args{
					format: "format",
					vals: []interface{}{
						"name",
					},
					fn: fn,
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if cnt != 3 {
						return errors.Errorf("called cnt is wrong. want: %v, got: %v", 3, cnt)
					}

					if !errors.Is(warn, fnErr) {
						return errors.Errorf("argument of warn funcion is wrong. want: %v, got: %v", warn, fnErr)
					}

					if !errors.Is(warn, fnErr) {
						return errors.Errorf("argument of error function is wrong. want: %v, got: %v", err, fnErr)
					}

					return nil
				},
				recovery: true,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.recovery {
					if err := recover(); err == nil {
						t.Error("panic is nil")
					}
				}
			}()

			logger = tt.global.l

			retryOutf(tt.args.fn, tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}
