package retry

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestNew(t *testing.T) {
	type test struct {
		name string
		opts []Option
		want *retry
	}

	tests := []test{
		func() test {
			return test{
				name: "returns retry object when options is empty",
				want: &retry{
					warnFn:  nopFunc,
					errorFn: nopFunc,
				},
			}
		}(),

		func() test {
			errorFn := func(vals ...interface{}) {}

			return test{
				name: "returns retry object when WithError options is on",
				opts: []Option{
					WithError(errorFn),
				},
				want: &retry{
					warnFn:  nopFunc,
					errorFn: errorFn,
				},
			}
		}(),

		func() test {
			warnFn := func(vals ...interface{}) {}

			return test{
				name: "returns retry object when WithWarn options is on",
				opts: []Option{
					WithWarn(warnFn),
				},
				want: &retry{
					warnFn:  warnFn,
					errorFn: nopFunc,
				},
			}
		}(),

		func() test {
			warnFn := func(vals ...interface{}) {}
			errorFn := func(vals ...interface{}) {}

			return test{
				name: "returns retry object when WithError and WithWarn options is on",
				opts: []Option{
					WithWarn(warnFn),
					WithError(errorFn),
				},
				want: &retry{
					warnFn:  warnFn,
					errorFn: errorFn,
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := New(tt.opts...).(*retry)
			if !ok {
				t.Errorf("type is invalid")
			}

			if reflect.ValueOf(got.errorFn).Pointer() != reflect.ValueOf(tt.want.errorFn).Pointer() {
				t.Error("errorfn is not equals")
			}

			if reflect.ValueOf(got.warnFn).Pointer() != reflect.ValueOf(tt.want.warnFn).Pointer() {
				t.Error("warnfn is not equals")
			}
		})
	}
}

func TestOut(t *testing.T) {
	type args struct {
		fn     func(vals ...interface{}) error
		format string
		vals   []interface{}
	}

	type field struct {
		warnFn  func(...interface{})
		errorFn func(...interface{})
	}

	type test struct {
		name      string
		args      args
		field     field
		panicked  bool
		checkFunc func() error
	}

	tests := []test{
		func() test {
			cnt := 0
			fn := func(vals ...interface{}) error {
				cnt++
				return nil
			}
			return test{
				name: "returns nothing when fn returns nil",
				args: args{
					fn: fn,
				},
				checkFunc: func() error {
					if cnt != 1 {
						return errors.Errorf("called cnt is wrong. want: %v, but got: %v", 1, cnt)
					}
					return nil
				},
			}
		}(),

		func() test {
			return test{
				name: "returns nothing when fn is nil",
				checkFunc: func() error {
					return nil
				},
			}
		}(),

		func() test {
			err := errors.New("fn error")
			fn := func(vals ...interface{}) error {
				return err
			}

			var gotWarnErr error
			warnFn := func(vals ...interface{}) {
				if len(vals) == 1 {
					gotWarnErr = vals[0].(error)
				}
			}

			var gotError error
			errorFn := func(vals ...interface{}) {
				if len(vals) == 1 {
					gotError = vals[0].(error)
				}
			}

			return test{
				name: "panic when fn returns error",
				args: args{
					fn: fn,
				},
				field: field{
					warnFn:  warnFn,
					errorFn: errorFn,
				},
				checkFunc: func() error {
					if !errors.Is(gotWarnErr, err) {
						return errors.New("warnFn argument is not wrong")
					}

					if !errors.Is(gotError, err) {
						return errors.New("errorFn argument is not wrong")
					}
					return nil
				},
				panicked: true,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.panicked {
					if e := recover(); e != nil {
						if err := tt.checkFunc(); err != nil {
							t.Error(err)
						}
					} else {
						t.Error("panic not occurs")
					}
				}
			}()

			r := &retry{
				warnFn:  tt.field.warnFn,
				errorFn: tt.field.errorFn,
			}
			r.Out(tt.args.fn, tt.args.vals...)

			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}

}

func TestOutf(t *testing.T) {
	type args struct {
		fn     func(format string, vals ...interface{}) error
		format string
		vals   []interface{}
	}

	type field struct {
		warnFn  func(...interface{})
		errorFn func(...interface{})
	}

	type test struct {
		name      string
		args      args
		field     field
		panicked  bool
		checkFunc func() error
	}

	tests := []test{
		func() test {
			cnt := 0
			fn := func(format string, vals ...interface{}) error {
				cnt++
				return nil
			}
			return test{
				name: "returns nothing when fn returns nil",
				args: args{
					fn: fn,
				},
				checkFunc: func() error {
					if cnt != 1 {
						return errors.Errorf("called cnt is wrong. want: %v, but got: %v", 1, cnt)
					}
					return nil
				},
			}
		}(),

		func() test {
			return test{
				name: "returns nothing when fn is nil",
				checkFunc: func() error {
					return nil
				},
			}
		}(),

		func() test {
			err := errors.New("fn error")
			fn := func(format string, vals ...interface{}) error {
				return err
			}

			var gotWarnErr error
			warnFn := func(vals ...interface{}) {
				if len(vals) == 1 {
					gotWarnErr = vals[0].(error)
				}
			}

			var gotError error
			errorFn := func(vals ...interface{}) {
				if len(vals) == 1 {
					gotError = vals[0].(error)
				}
			}

			return test{
				name: "panic when fn returns error",
				args: args{
					fn: fn,
				},
				field: field{
					warnFn:  warnFn,
					errorFn: errorFn,
				},
				checkFunc: func() error {
					if !errors.Is(gotWarnErr, err) {
						return errors.New("warnFn argument is not wrong")
					}

					if !errors.Is(gotError, err) {
						return errors.New("errorFn argument is not wrong")
					}
					return nil
				},
				panicked: true,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.panicked {
					if e := recover(); e != nil {
						if err := tt.checkFunc(); err != nil {
							t.Error(err)
						}
					} else {
						t.Error("panic not occurs")
					}
				}
			}()

			r := &retry{
				warnFn:  tt.field.warnFn,
				errorFn: tt.field.errorFn,
			}
			r.Outf(tt.args.fn, tt.args.format, tt.args.vals...)

			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}
