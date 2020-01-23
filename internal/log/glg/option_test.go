package glg

import (
	"errors"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/log/retry"
)

func TestWithEnableJSON(t *testing.T) {
	type test struct {
		name      string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			checkFunc: func(Option) error {
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithEnableJSON()
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithLevel(t *testing.T) {
	type test struct {
		name      string
		lv        string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success",
			lv:   "debug",
			checkFunc: func(opt Option) error {
				got := new(GlgLogger)
				opt(got)

				if got.lv != DEBUG {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},

		{
			name: "set nothing",
			checkFunc: func(opt Option) error {
				got := &GlgLogger{
					lv: ERROR,
				}
				opt(got)

				if got.lv != ERROR {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithLevel(tt.lv)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithRetryOut(t *testing.T) {
	type test struct {
		name      string
		fn        retry.Out
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			fn := func(fn func(vals ...interface{}) error, vals ...interface{}) {}
			return test{
				name: "set success",
				fn:   fn,
				checkFunc: func(opt Option) error {
					got := new(GlgLogger)
					opt(got)

					if reflect.ValueOf(got.rout).Pointer() != reflect.ValueOf(fn).Pointer() {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),

		func() test {
			fn := func(fn func(vals ...interface{}) error, vals ...interface{}) {}
			return test{
				name: "set nothing",
				fn:   nil,
				checkFunc: func(opt Option) error {
					got := &GlgLogger{
						rout: fn,
					}
					opt(got)

					if reflect.ValueOf(got.rout).Pointer() != reflect.ValueOf(fn).Pointer() {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRetryOut(tt.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithRetryOutf(t *testing.T) {
	type test struct {
		name      string
		fn        retry.Outf
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			fn := func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{}) {}
			return test{
				name: "set success",
				fn:   fn,
				checkFunc: func(opt Option) error {
					got := new(GlgLogger)
					opt(got)

					if reflect.ValueOf(got.routf).Pointer() != reflect.ValueOf(fn).Pointer() {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),

		func() test {
			fn := func(fn func(format string, vals ...interface{}) error, format string, vals ...interface{}) {}
			return test{
				name: "set nothing",
				fn:   nil,
				checkFunc: func(opt Option) error {
					got := &GlgLogger{
						routf: fn,
					}
					opt(got)

					if reflect.ValueOf(got.routf).Pointer() != reflect.ValueOf(fn).Pointer() {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRetryOutf(tt.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
