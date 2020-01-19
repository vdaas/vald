//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
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
	"github.com/vdaas/vald/internal/errors"
)

func TestNewExpBackoff(t *testing.T) {
	type test struct {
		name        string
		opts        []Option
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
			got := NewExpBackoff(tt.opts...)

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
			wantRes := new(http.Response)

			tr := &roundTripMock{
				RoundTripFunc: func(*http.Request) (*http.Response, error) {
					return wantRes, nil
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

					if !reflect.DeepEqual(res, wantRes) {
						return errors.Errorf("res not equals. want: %v, got: %v", wantRes, err)
					}

					return nil
				},
			}
		}(),

		func() test {
			wantRes := new(http.Response)

			tr := &roundTripMock{
				RoundTripFunc: func(*http.Request) (*http.Response, error) {
					return wantRes, nil
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
						return errors.Errorf("error not nil. err: %v", err)
					}

					if !reflect.DeepEqual(res, wantRes) {
						return errors.Errorf("res not equals. want: %v, got: %v", wantRes, err)
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
					return res, errors.New("faild")
				},
			}

			bo := &backoffMock{
				DoFunc: func(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
					return fn()
				},
			}

			return test{
				name: "returns error when Do function returns error",
				args: args{
					req: new(http.Request),
				},
				field: field{
					transport: tr,
					bo:        bo,
				},
				checkFunc: func(res *http.Response, err error) error {
					if err == nil {
						return errors.New("err is nil")
					}

					if res != nil {
						return errors.Errorf("res not nil. res: %v", res)
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
				name: "returns error when type conversion error occurs",
				args: args{
					req: new(http.Request),
				},
				field: field{
					transport: tr,
					bo:        bo,
				},
				checkFunc: func(res *http.Response, err error) error {
					if err == nil {
						return errors.New("err is nil")
					}

					if res != nil {
						return errors.Errorf("res not nil. res: %v", res)
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
