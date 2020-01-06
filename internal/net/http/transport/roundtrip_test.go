package transport

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
)

func TestNewExpBackoff(t *testing.T) {
	type args struct {
		opts []Option
	}

	type test struct {
		name        string
		args        args
		initialized bool
	}

	tests := []test{
		{
			name:        "initialize success",
			initialized: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewExpBackoff(tt.args.opts...)

			if (got != nil) != tt.initialized {
				t.Error("New() is wrong")
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	type args struct {
		req *http.Request
	}

	type field struct {
		bo        backoff.Backoff
		transport http.RoundTripper
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(*http.Response, error) error
	}

	tests := []test{
		func() test {
			wres := new(http.Response)

			tr := &roundTripMock{
				RoundTripFunc: func(*http.Request) (*http.Response, error) {
					return wres, nil
				},
			}

			return test{
				name: "returns not error when backoff object is nil",
				field: field{
					transport: tr,
				},
				checkFunc: func(res *http.Response, err error) error {
					if err != nil {
						return fmt.Errorf("error not nil. err: %v", err)
					}

					if !reflect.DeepEqual(res, wres) {
						return fmt.Errorf("res not equals. want: %v, got: %v", wres, err)
					}

					return nil
				},
			}
		}(),

		func() test {
			wres := new(http.Response)

			tr := &roundTripMock{
				RoundTripFunc: func(*http.Request) (*http.Response, error) {
					return wres, nil
				},
			}

			bo := &backoffMock{
				DoFunc: func(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
					return fn()
				},
			}

			return test{
				name: "returns not error when backoff object is not nil",
				args: args{
					req: new(http.Request),
				},
				field: field{
					transport: tr,
					bo:        bo,
				},
				checkFunc: func(res *http.Response, err error) error {
					if err != nil {
						return fmt.Errorf("error not nil. err: %v", err)
					}

					if !reflect.DeepEqual(res, wres) {
						return fmt.Errorf("res not equals. want: %v, got: %v", wres, err)
					}

					return nil
				},
			}
		}(),

		func() test {
			res := &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Body:       ioutil.NopCloser(new(bytes.Buffer)),
			}

			tr := &roundTripMock{
				RoundTripFunc: func(*http.Request) (*http.Response, error) {
					return res, fmt.Errorf("faild")
				},
			}

			bo := &backoffMock{
				DoFunc: func(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
					return fn()
				},
			}

			return test{
				name: "backoff returns an error",
				args: args{
					req: new(http.Request),
				},
				field: field{
					transport: tr,
					bo:        bo,
				},
				checkFunc: func(res *http.Response, err error) error {
					if err == nil {
						return fmt.Errorf("err is nil")
					}

					if res != nil {
						return fmt.Errorf("res not nil. res: %v", res)
					}

					return nil
				},
			}
		}(),

		func() test {
			tr := &roundTripMock{
				RoundTripFunc: func(*http.Request) (*http.Response, error) {
					return nil, nil
				},
			}

			bo := &backoffMock{
				DoFunc: func(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
					_, err := fn()
					return "dumy", err
				},
			}

			return test{
				name: "returns type conversion error",
				args: args{
					req: new(http.Request),
				},
				field: field{
					transport: tr,
					bo:        bo,
				},
				checkFunc: func(res *http.Response, err error) error {
					if err == nil {
						return fmt.Errorf("err is nil")
					}

					if res != nil {
						return fmt.Errorf("res not nil. res: %v", res)
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ert{
				transport: tt.field.transport,
				bo:        tt.field.bo,
			}

			res, err := e.RoundTrip(tt.args.req)
			if err := tt.checkFunc(res, err); err != nil {
				t.Error(err)
			}
		})
	}
}
