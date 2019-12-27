package middleware

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
)

func TestWithErrorGroup(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(TimeoutOption) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "10s",
			},
			checkFunc: func(opt TimeoutOption) error {
				got := new(timeout)
				opt(got)

				if got.dur != 10*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
		{
			name: "set default value",
			args: args{
				dur: "ok",
			},
			checkFunc: func(opt TimeoutOption) error {
				got := new(timeout)
				opt(got)

				if got.dur != 3*time.Second {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithTimeout(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	type args struct {
		eg errgroup.Group
	}

	type test struct {
		name      string
		args      args
		checkFunc func(TimeoutOption) error
	}

	tests := []test{
		func() test {
			eg := errgroup.Get()

			return test{
				name: "set success",
				args: args{
					eg: eg,
				},
				checkFunc: func(opt TimeoutOption) error {
					got := new(timeout)
					opt(got)

					if !reflect.DeepEqual(got.eg, eg) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithErrorGroup(tt.args.eg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
