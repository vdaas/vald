package log

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/mock"
)

func Test_option_apply(t *testing.T) {
	type args struct {
		opt []Option
	}

	type test struct {
		name string
		opts []Option
		want *option
	}

	tests := []test{
		func() test {
			logger := new(mock.Logger)

			return test{
				name: "returns option",
				opts: []Option{
					WithLogger(logger),
				},
				want: &option{
					logger: logger,
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := (&option{}).apply(tt.opts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	type test struct {
		name      string
		logger    Logger
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			logger := new(mock.Logger)

			return test{
				name:   "set success",
				logger: logger,
				checkFunc: func(opt Option) error {
					got := new(option)
					opt(got)

					if !reflect.DeepEqual(got.logger, logger) {
						return errors.New("invalid param was set")
					}

					return nil
				},
			}
		}(),

		func() test {
			return test{
				name: "set default",
				checkFunc: func(opt Option) error {
					got := new(option)
					opt(got)

					if got.logger != nil {
						return errors.New("invalid param was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithLogger(tt.logger)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
