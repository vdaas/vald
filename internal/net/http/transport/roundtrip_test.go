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

// Package transport provides http transport roundtrip option
package transport

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

var (
	// Goroutine leak is detected by `fastime`, but it should be ignored in the test because it is an external package.
	goleakIgnoreOptions = []goleak.Option{
		goleak.IgnoreTopFunction("github.com/kpango/fastime.(*Fastime).StartTimerD.func1"),
	}
)

func TestNewExpBackoff(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want http.RoundTripper
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, http.RoundTripper) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got http.RoundTripper) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "initialize success",
			want: want{
				want: &ert{
					transport: http.DefaultTransport,
				},
			},
		},
		func() test {
			b := backoff.New()
			return test{
				name: "initialize success with option",
				args: args{
					opts: []Option{
						WithBackoff(b),
					},
				},
				want: want{
					want: &ert{
						transport: http.DefaultTransport,
						bo:        b,
					},
				},
			}
		}(),
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got := NewExpBackoff(test.args.opts...)
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ert_RoundTrip(t *testing.T) {
	type args struct {
		req *http.Request
	}
	type fields struct {
		transport http.RoundTripper
		bo        backoff.Backoff
	}
	type want struct {
		wantRes *http.Response
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *http.Response, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *http.Response, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got = %v, want %v", gotRes, w.wantRes)
		}
		return nil
	}
	tests := []test{
		{
			name: "return roundtrip response if backoff is nil",
			args: args{
				req: nil,
			},
			fields: fields{
				transport: &roundTripMock{
					RoundTripFunc: func(*http.Request) (*http.Response, error) {
						return &http.Response{
							Status: "200",
						}, nil
					},
				},
			},
			want: want{
				wantRes: &http.Response{
					Status: "200",
				},
			},
		},
		{
			name: "return backoff response if backoff is not nil",
			args: args{
				req: &http.Request{},
			},
			fields: fields{
				transport: &roundTripMock{
					RoundTripFunc: func(*http.Request) (*http.Response, error) {
						return nil, errors.New("error")
					},
				},
				bo: &backoffMock{
					DoFunc: func(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
						return &http.Response{
							Status: "200",
						}, nil
					},
				},
			},
			want: want{
				wantRes: &http.Response{
					Status: "200",
				},
			},
		},
		{
			name: "return backoff response if backoff is not nil",
			args: args{
				req: &http.Request{},
			},
			fields: fields{
				transport: &roundTripMock{
					RoundTripFunc: func(*http.Request) (*http.Response, error) {
						return nil, errors.New("error")
					},
				},
				bo: &backoffMock{
					DoFunc: func(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
						return &http.Response{
							Status: "200",
						}, nil
					},
				},
			},
			want: want{
				wantRes: &http.Response{
					Status: "200",
				},
			},
		},
		{
			name: "return backoff error",
			args: args{
				req: &http.Request{},
			},
			fields: fields{
				bo: &backoffMock{
					DoFunc: func(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
						return nil, errors.New("error")
					},
				},
			},
			want: want{
				err: errors.New("error"),
			},
		},
		{
			name: "return error when backoff return invalid type result",
			args: args{
				req: &http.Request{},
			},
			fields: fields{
				bo: &backoffMock{
					DoFunc: func(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
						return struct{}{}, nil
					},
				},
			},
			want: want{
				err: errors.ErrInvalidTypeConversion(struct{}{}, &http.Response{}),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &ert{
				transport: test.fields.transport,
				bo:        test.fields.bo,
			}

			gotRes, err := e.RoundTrip(test.args.req)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ert_roundTrip(t *testing.T) {
	type args struct {
		req *http.Request
	}
	type fields struct {
		transport http.RoundTripper
		bo        backoff.Backoff
	}
	type want struct {
		wantRes *http.Response
		err     error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, *http.Response, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotRes *http.Response, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got = %v, want %v", gotRes, w.wantRes)
		}
		return nil
	}
	tests := []test{
		{
			name: "roundtrip return success",
			args: args{
				req: &http.Request{},
			},
			fields: fields{
				transport: &roundTripMock{
					RoundTripFunc: func(*http.Request) (*http.Response, error) {
						return &http.Response{
							Status: "200",
						}, nil
					},
				},
			},
			want: want{
				wantRes: &http.Response{
					Status: "200",
				},
			},
		},
		{
			name: "roundtrip return empty response with error",
			args: args{
				req: &http.Request{},
			},
			fields: fields{
				transport: &roundTripMock{
					RoundTripFunc: func(*http.Request) (*http.Response, error) {
						return nil, errors.New("error")
					},
				},
			},
			want: want{
				err: errors.New("error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			e := &ert{
				transport: test.fields.transport,
				bo:        test.fields.bo,
			}

			gotRes, err := e.roundTrip(test.args.req)
			if err := test.checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
