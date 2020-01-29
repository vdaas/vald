package glg

import (
	"reflect"
	"testing"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/format"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/log/mock"
	"github.com/vdaas/vald/internal/log/retry"
)

func TestWithGlg(t *testing.T) {
	type test struct {
		name      string
		g         *glg.Glg
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			g := glg.New()

			return test{
				name: "set success when glg is not nil",
				g:    g,
				checkFunc: func(opt Option) error {
					got := new(logger)
					opt(got)

					if !reflect.DeepEqual(got.glg, g) {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),

		func() test {
			g := glg.New()

			return test{
				name: "set success when glg is not nil",
				g:    nil,
				checkFunc: func(opt Option) error {
					got := &logger{
						glg: g,
					}
					opt(got)

					if !reflect.DeepEqual(got.glg, g) {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGlg(tt.g)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithFormat(t *testing.T) {
	type test struct {
		name      string
		str       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success when str is JSON",
			str:  format.JSON.String(),
			checkFunc: func(opt Option) error {
				got := new(logger)
				opt(got)

				if got.format != format.JSON {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},

		{
			name: "returns nothing when str is empty",
			checkFunc: func(opt Option) error {
				got := &logger{
					format: format.RAW,
				}
				opt(got)

				if got.format != format.RAW {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithFormat(tt.str)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithLevel(t *testing.T) {
	type test struct {
		name      string
		str       string
		checkFunc func(Option) error
	}

	tests := []test{
		{
			name: "set success when str is DEBUG",
			str:  "DEBUG",
			checkFunc: func(opt Option) error {
				got := new(logger)
				opt(got)

				if got.level != level.DEBUG {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},

		{
			name: "returns nothing when str is empty",
			checkFunc: func(opt Option) error {
				got := &logger{
					level: level.ERROR,
				}
				opt(got)

				if got.level != level.ERROR {
					return errors.New("invalid params was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithLevel(tt.str)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithRetry(t *testing.T) {
	type test struct {
		name      string
		rt        retry.Retry
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			rt := new(mock.Retry)
			return test{
				name: "set success when rt is not nil",
				rt:   rt,
				checkFunc: func(opt Option) error {
					got := new(logger)
					opt(got)

					if !reflect.DeepEqual(got.retry, rt) {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),

		func() test {
			rt := new(mock.Retry)

			return test{
				name: "set success when rt is not nil",
				rt:   nil,
				checkFunc: func(opt Option) error {
					got := &logger{
						retry: rt,
					}
					opt(got)

					if !reflect.DeepEqual(got.retry, rt) {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRetry(tt.rt)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
