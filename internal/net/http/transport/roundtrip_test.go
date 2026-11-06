// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transport

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	goleak.VerifyTestMain(m)
}

func TestNewExpBackoff(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
	}
	type want struct {
		want http.RoundTripper
	}
	type test struct {
		want       want
		checkFunc  func(want, http.RoundTripper) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, got http.RoundTripper) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := NewExpBackoff(test.args.opts...)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ert_RoundTrip(t *testing.T) {
	t.Parallel()
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
		fields     fields
		want       want
		args       args
		checkFunc  func(want, *http.Response, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
	}
	defaultCheckFunc := func(w want, gotRes *http.Response, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
						return &http.Response{
							Status: "200",
						}, nil
					},
				},
				bo: &backoffMock{
					DoFunc: func(ctx context.Context, fn func(context.Context) (any, bool, error)) (any, error) {
						val, _, err := fn(ctx)
						return val, err
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
			name: "return default roundtrip response if backoff is not nil",
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
				bo: &backoffMock{
					DoFunc: func(ctx context.Context, fn func(context.Context) (any, bool, error)) (any, error) {
						val, _, err := fn(ctx)
						return val, err
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
			name: "return default roundtrip response if backoff use the default roundtrip",
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
				bo: &backoffMock{
					DoFunc: func(ctx context.Context, fn func(context.Context) (any, bool, error)) (any, error) {
						val, _, err := fn(ctx)
						return val, err
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
					DoFunc: func(ctx context.Context, fn func(context.Context) (any, bool, error)) (any, error) {
						return nil, errors.New("error")
					},
				},
			},
			want: want{
				err: errors.New("error"),
			},
		},
		{
			name: "return default roundtrip error if backoff use the default roundtrip",
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
					DoFunc: func(ctx context.Context, fn func(context.Context) (any, bool, error)) (any, error) {
						val, _, err := fn(ctx)
						return val, err
					},
				},
			},
			want: want{
				err: errors.New("error"),
			},
		},
		{
			name: "return retryable error",
			args: args{
				req: &http.Request{},
			},
			fields: fields{
				transport: &roundTripMock{
					RoundTripFunc: func(*http.Request) (*http.Response, error) {
						return nil, errors.Wrap(errors.ErrTransportRetryable, "error")
					},
				},
				bo: &backoffMock{
					DoFunc: func(ctx context.Context, fn func(context.Context) (any, bool, error)) (any, error) {
						val, _, err := fn(ctx)
						return val, err
					},
				},
			},
			want: want{
				err: errors.Wrap(errors.ErrTransportRetryable, "error"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			e := &ert{
				transport: test.fields.transport,
				bo:        test.fields.bo,
			}

			gotRes, err := e.RoundTrip(test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ert_doRoundTrip(t *testing.T) {
	t.Parallel()
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
		fields     fields
		want       want
		args       args
		checkFunc  func(want, *http.Response, error) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
	}
	defaultCheckFunc := func(w want, gotRes *http.Response, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotRes, w.wantRes) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
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
							Status: strconv.Itoa(http.StatusOK),
						}, nil
					},
				},
			},
			want: want{
				wantRes: &http.Response{
					Status: strconv.Itoa(http.StatusOK),
				},
			},
		},
		{
			name: "roundtrip return error",
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
		{
			name: "roundtrip return retryable error when status code is 502",
			args: args{
				req: &http.Request{},
			},
			fields: fields{
				transport: &roundTripMock{
					RoundTripFunc: func(*http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusBadGateway,
							Body:       io.NopCloser(bytes.NewBufferString("abc")),
						}, nil
					},
				},
			},
			want: want{
				err: errors.ErrTransportRetryable,
			},
		},
		{
			name: "roundtrip return success when status code is 301",
			args: args{
				req: &http.Request{},
			},
			fields: fields{
				transport: &roundTripMock{
					RoundTripFunc: func(*http.Request) (*http.Response, error) {
						return &http.Response{
							Status: strconv.Itoa(http.StatusMovedPermanently),
						}, nil
					},
				},
			},
			want: want{
				wantRes: &http.Response{
					Status: strconv.Itoa(http.StatusMovedPermanently),
				},
			},
		},
		{
			name: "roundtrip return retryable error even when error occurred and roundtrip response is not nil",
			args: args{
				req: &http.Request{},
			},
			fields: fields{
				transport: &roundTripMock{
					RoundTripFunc: func(*http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusBadGateway,
							Body:       io.NopCloser(bytes.NewBufferString("abc")),
						}, errors.New("dummy")
					},
				},
			},
			want: want{
				err: errors.ErrTransportRetryable,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			e := &ert{
				transport: test.fields.transport,
				bo:        test.fields.bo,
			}

			gotRes, err := e.doRoundTrip(test.args.req)
			if err := checkFunc(test.want, gotRes, err); err != nil {
				tt.Errorf("error = %v", err)
			}

			if gotRes != nil {
				defer closeBody(gotRes.Body)
			}
			if test.args.req != nil {
				defer closeBody(test.args.req.Body)
			}
		})
	}
}

func Test_retryableStatusCode(t *testing.T) {
	t.Parallel()
	type args struct {
		status int
	}
	type want struct {
		want bool
	}
	type test struct {
		checkFunc  func(want, bool) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
		args       args
		want       want
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return true when response status is 429(TooManyRequest)",
			args: args{
				status: http.StatusTooManyRequests,
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false when response status is 301(MovedPermanently)",
			args: args{
				status: http.StatusMovedPermanently,
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false when response status is 200(OK)",
			args: args{
				status: http.StatusOK,
			},
			want: want{
				want: false,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := retryableStatusCode(test.args.status)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_closeBody(t *testing.T) {
	t.Parallel()
	type args struct {
		rc io.ReadCloser
	}
	type want struct{}
	type test struct {
		want       want
		args       args
		checkFunc  func(io.ReadCloser, want) error
		beforeFunc func(args)
		afterFunc  func(args)
		name       string
	}
	defaultCheckFunc := func(rc io.ReadCloser, w want) error {
		if i, err := rc.Read([]byte{}); i != 0 || !errors.Is(err, io.EOF) {
			return errors.Errorf("connection not closed, num: %d, err: %v", i, err)
		}
		return nil
	}
	tests := []test{
		func() test {
			rr := httptest.NewRecorder()
			rr.WriteString("abc")
			res := rr.Result()

			return test{
				name: "close response body",
				args: args{
					rc: res.Body,
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			closeBody(test.args.rc)
			if err := checkFunc(test.args.rc, test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// trackBody records reads and closes of a request body.
type trackBody struct {
	reads, closes int
}

func (b *trackBody) Read([]byte) (int, error) {
	b.reads++
	return 0, io.EOF
}

func (b *trackBody) Close() error {
	b.closes++
	return nil
}

// RoundTrip drained and closed the original req.Body after the transport had
// taken it (the transport may still be using it from its own goroutine), and it
// retried requests whose body could not be rewound (GetBody nil).
func Test_ert_RoundTrip_requestBodyOwnership(t *testing.T) {
	t.Parallel()
	// retryingBackoff calls fn while it reports retryable, up to 3 attempts.
	retryingBackoff := &backoffMock{
		DoFunc: func(ctx context.Context, fn func(context.Context) (any, bool, error)) (any, error) {
			var (
				val   any
				retry bool
				err   error
			)
			for range 3 {
				val, retry, err = fn(ctx)
				if !retry {
					break
				}
			}
			return val, err
		},
	}
	tests := []struct {
		bo          backoff.Backoff
		name        string
		wantCalls   int
		wantReads   int
		wantCloses  int
		rewindable  bool
		retryStatus bool
	}{
		{
			name:      "no backoff: handed body is left to the transport",
			wantCalls: 1,
		},
		{
			name:        "non-rewindable body is not retried and left to the transport",
			bo:          retryingBackoff,
			retryStatus: true,
			wantCalls:   1,
		},
		{
			name:        "rewindable body is retried; the unsent original is closed without draining",
			bo:          retryingBackoff,
			rewindable:  true,
			retryStatus: true,
			wantCalls:   3,
			wantCloses:  1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			orig := new(trackBody)
			req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://example.invalid", orig)
			if err != nil {
				t.Fatal(err)
			}
			if tc.rewindable {
				req.GetBody = func() (io.ReadCloser, error) {
					return io.NopCloser(bytes.NewReader(nil)), nil
				}
			}
			var calls int
			e := &ert{
				// The mock never touches the body, like a transport that is still
				// writing it from its own goroutine.
				transport: &roundTripMock{
					RoundTripFunc: func(*http.Request) (*http.Response, error) {
						calls++
						if tc.retryStatus {
							return nil, errors.ErrTransportRetryable
						}
						return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
					},
				},
				bo: tc.bo,
			}
			res, err := e.RoundTrip(req)
			if res != nil && res.Body != nil {
				if cerr := res.Body.Close(); cerr != nil {
					t.Error(cerr)
				}
			}
			if tc.retryStatus && !errors.Is(err, errors.ErrTransportRetryable) {
				t.Errorf("RoundTrip() error = %v, want ErrTransportRetryable", err)
			}
			if !tc.retryStatus && err != nil {
				t.Errorf("RoundTrip() error = %v, want nil", err)
			}
			if calls != tc.wantCalls {
				t.Errorf("transport calls = %d, want %d", calls, tc.wantCalls)
			}
			if orig.reads != tc.wantReads {
				t.Errorf("original body reads by RoundTrip = %d, want %d", orig.reads, tc.wantReads)
			}
			if orig.closes != tc.wantCloses {
				t.Errorf("original body closes by RoundTrip = %d, want %d", orig.closes, tc.wantCloses)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
