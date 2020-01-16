package log

import (
	"reflect"
	"testing"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/mock"
)

func TestNewGlg(t *testing.T) {
	type test struct {
		name string
		log  *glg.Glg
		want *glglogger
	}

	tests := []test{
		func() test {
			log := glg.Get()

			return test{
				name: "initialize success",
				log:  log,
				want: &glglogger{
					log: log,
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGlg(tt.log)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestDefaultGlg(t *testing.T) {
	type test struct {
		name string
		want *glglogger
	}

	tests := []test{
		{
			name: "default glg success",
			want: &glglogger{
				log: glg.Get().SetMode(glg.NONE),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultGlg()
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestGlgInfo(t *testing.T) {
	type args struct {
		vals []interface{}
	}

	type field struct {
		log *glg.Glg
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		field     field
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
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
				name: "output is successes",
				args: args{
					vals: []interface{}{
						"name",
					},
				},
				field: field{
					log: glg.Get().SetMode(glg.NONE),
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if warn != nil {
						t.Errorf("argument of warn funcion is not nil. err: %v", warn)
					}

					if err != nil {
						t.Errorf("argument of error funcion is not nil. err: %v", err)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()

			logger = tt.global.l

			gl := &glglogger{
				log: tt.field.log,
			}

			gl.Info(tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGlgInfof(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}

	type field struct {
		log *glg.Glg
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		field     field
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
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
				name: "output is successes",
				args: args{
					format: "%v",
					vals: []interface{}{
						"name",
					},
				},
				field: field{
					log: glg.Get().SetMode(glg.NONE),
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if warn != nil {
						t.Errorf("argument of warn funcion is not nil. err: %v", warn)
					}

					if err != nil {
						t.Errorf("argument of error funcion is not nil. err: %v", err)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()

			logger = tt.global.l

			gl := &glglogger{
				log: tt.field.log,
			}

			gl.Infof(tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGlgDebug(t *testing.T) {
	type args struct {
		vals []interface{}
	}

	type field struct {
		log *glg.Glg
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		field     field
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
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
				name: "output is successes",
				args: args{
					vals: []interface{}{
						"name",
					},
				},
				field: field{
					log: glg.Get().SetMode(glg.NONE),
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if warn != nil {
						t.Errorf("argument of warn funcion is not nil. err: %v", warn)
					}

					if err != nil {
						t.Errorf("argument of error funcion is not nil. err: %v", err)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()

			logger = tt.global.l

			gl := &glglogger{
				log: tt.field.log,
			}

			gl.Debug(tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGlgDebugf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}

	type field struct {
		log *glg.Glg
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		field     field
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
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
				name: "output is successes",
				args: args{
					format: "%v",
					vals: []interface{}{
						"name",
					},
				},
				field: field{
					log: glg.Get().SetMode(glg.NONE),
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if warn != nil {
						t.Errorf("argument of warn funcion is not nil. err: %v", warn)
					}

					if err != nil {
						t.Errorf("argument of error funcion is not nil. err: %v", err)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()

			logger = tt.global.l

			gl := &glglogger{
				log: tt.field.log,
			}

			gl.Debugf(tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGlgWarn(t *testing.T) {
	type args struct {
		vals []interface{}
	}

	type field struct {
		log *glg.Glg
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		field     field
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
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
				name: "output is successes",
				args: args{
					vals: []interface{}{
						"name",
					},
				},
				field: field{
					log: glg.Get().SetMode(glg.NONE),
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if warn != nil {
						t.Errorf("argument of warn funcion is not nil. err: %v", warn)
					}

					if err != nil {
						t.Errorf("argument of error funcion is not nil. err: %v", err)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()

			logger = tt.global.l

			gl := &glglogger{
				log: tt.field.log,
			}

			gl.Warn(tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGlgWarnf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}

	type field struct {
		log *glg.Glg
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		field     field
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
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
				name: "output is successes",
				args: args{
					format: "%v",
					vals: []interface{}{
						"name",
					},
				},
				field: field{
					log: glg.Get().SetMode(glg.NONE),
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if warn != nil {
						t.Errorf("argument of warn funcion is not nil. err: %v", warn)
					}

					if err != nil {
						t.Errorf("argument of error funcion is not nil. err: %v", err)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()

			logger = tt.global.l

			gl := &glglogger{
				log: tt.field.log,
			}

			gl.Warnf(tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGlgError(t *testing.T) {
	type args struct {
		vals []interface{}
	}

	type field struct {
		log *glg.Glg
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		field     field
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
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
				name: "output is successes",
				args: args{
					vals: []interface{}{
						"name",
					},
				},
				field: field{
					log: glg.Get().SetMode(glg.NONE),
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if warn != nil {
						t.Errorf("argument of warn funcion is not nil. err: %v", warn)
					}

					if err != nil {
						t.Errorf("argument of error funcion is not nil. err: %v", err)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()

			logger = tt.global.l

			gl := &glglogger{
				log: tt.field.log,
			}

			gl.Error(tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGlgErrorf(t *testing.T) {
	type args struct {
		format string
		vals   []interface{}
	}

	type field struct {
		log *glg.Glg
	}

	type global struct {
		l Logger
	}

	type test struct {
		name      string
		args      args
		field     field
		global    global
		checkFunc func() error
	}

	tests := []test{
		func() test {
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
				name: "output is successes",
				args: args{
					format: "%v",
					vals: []interface{}{
						"name",
					},
				},
				field: field{
					log: glg.Get().SetMode(glg.NONE),
				},
				global: global{
					l: l,
				},
				checkFunc: func() error {
					if warn != nil {
						t.Errorf("argument of warn funcion is not nil. err: %v", warn)
					}

					if err != nil {
						t.Errorf("argument of error funcion is not nil. err: %v", err)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					t.Error(err)
				}
			}()

			logger = tt.global.l

			gl := &glglogger{
				log: tt.field.log,
			}

			gl.Errorf(tt.args.format, tt.args.vals...)
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

		func() test {
			fnErr := errors.New("fail")

			var cnt int
			fn := func(vals ...interface{}) error {
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

			out(tt.args.fn, tt.args.vals...)
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

			outf(tt.args.fn, tt.args.format, tt.args.vals...)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}
