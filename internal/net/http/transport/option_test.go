package transport

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
)

func TestWithRoundTripper(t *testing.T) {
	type test struct {
		name      string
		tr        http.RoundTripper
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			tr := new(roundTripMock)

			return test{
				name: "set success",
				tr:   tr,
				checkFunc: func(opt Option) error {
					got := new(ert)
					opt(got)

					if got, want := got.transport, tr; !reflect.DeepEqual(got, want) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithRoundTripper(tt.tr)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithBackoff(t *testing.T) {
	type test struct {
		name      string
		bo        backoff.Backoff
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			bo := new(backoffMock)

			return test{
				name: "set success",
				bo:   bo,
				checkFunc: func(opt Option) error {
					got := new(ert)
					opt(got)

					if got, want := got.bo, bo; !reflect.DeepEqual(got, want) {
						return errors.New("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithBackoff(tt.bo)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
