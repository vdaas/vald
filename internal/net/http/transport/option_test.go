package transport

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
)

func TestWithRoundTripper(t *testing.T) {
	type args struct {
		tr http.RoundTripper
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			tr := new(roundTripMock)

			return test{
				name: "set success",
				args: args{
					tr: tr,
				},
				checkFunc: func(opt Option) error {
					got := new(ert)
					opt(got)

					if got, want := got.transport, tr; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRoundTripper(tt.args.tr)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithBackoff(t *testing.T) {
	type args struct {
		bo backoff.Backoff
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			bo := new(backoffMock)

			return test{
				name: "set success",
				args: args{
					bo: bo,
				},
				checkFunc: func(opt Option) error {
					got := new(ert)
					opt(got)

					if got, want := got.bo, bo; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithBackoff(tt.args.bo)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
