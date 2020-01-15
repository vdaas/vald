package log

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestInit(t *testing.T) {
	type test struct {
		name string
		l    Logger
	}

	tests := []test{
		func() test {
			l := &loggerMock{}

			return test{
				name: "set success",
				l:    l,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.l)
			if !reflect.DeepEqual(tt.l, logger) {
				t.Errorf("logger is not equals. want: %v, got: %v", tt.l, logger)
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

			l := &loggerMock{
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
			Init(tt.global.l)
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

			l := &loggerMock{
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
			Init(tt.global.l)
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

			l := &loggerMock{
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
			Init(tt.global.l)
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

			l := &loggerMock{
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
			Init(tt.global.l)
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

			l := &loggerMock{
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
			Init(tt.global.l)
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

			l := &loggerMock{
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
			Init(tt.global.l)
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

			l := &loggerMock{
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
			Init(tt.global.l)
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

			l := &loggerMock{
				ErrorfFunc: func(format string, vals ...interface{}) {
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
			Init(tt.global.l)
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

			l := &loggerMock{
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
			Init(tt.global.l)
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

			l := &loggerMock{
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
			Init(tt.global.l)
		})
	}
}
