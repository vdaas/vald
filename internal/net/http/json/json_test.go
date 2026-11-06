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
package json

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	goleak.VerifyTestMain(m)
}

func TestEncodeResponse(t *testing.T) {
	t.Parallel()
	type args struct {
		w            http.ResponseWriter
		data         any
		contentTypes []string
		status       int
	}
	type test struct {
		checkFunc func(err error) error
		name      string
		args      args
	}
	tests := []test{
		func() test {
			w := new(httptest.ResponseRecorder)
			return test{
				name: "returns nil",
				args: args{
					w:      w,
					data:   []byte(`{"name":"vald"}`),
					status: http.StatusOK,
					contentTypes: []string{
						"application/json",
					},
				},
				checkFunc: func(err error) error {
					if err != nil {
						return errors.Errorf("err not equals. want: %v, got: %v", nil, err)
					}
					if got, want := w.Header().Get(rest.ContentType), "application/json"; got != want {
						return errors.Errorf("content-type not equals. want: %v, got: %v", want, got)
					}
					if got, want := w.Code, 200; got != want {
						return errors.Errorf("code not equals. want: %v, got: %v", want, got)
					}
					return nil
				},
			}
		}(),
		{
			name: "returns error when type is invalid",
			args: args{
				w:      new(httptest.ResponseRecorder),
				data:   make(chan struct{}),
				status: http.StatusOK,
				contentTypes: []string{
					"application/json",
				},
			},
			checkFunc: func(err error) error {
				if err == nil {
					return errors.New("err is nil")
				}
				return nil
			},
		},
		func() test {
			w := new(httptest.ResponseRecorder)
			return test{
				name: "does not write the header when encoding fails",
				args: args{
					w:      w,
					data:   make(chan struct{}),
					status: http.StatusOK,
					contentTypes: []string{
						"application/json",
					},
				},
				checkFunc: func(err error) error {
					if err == nil {
						return errors.New("err is nil")
					}
					if w.Code != 0 {
						return errors.Errorf("header was written. code: %v", w.Code)
					}
					if got := w.Header().Get(rest.ContentType); got != "" {
						return errors.Errorf("content-type was set: %v", got)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			err := EncodeResponse(test.args.w, test.args.data, test.args.status, test.args.contentTypes...)
			if err := test.checkFunc(err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDecodeRequest(t *testing.T) {
	t.Parallel()
	type args struct {
		r    *http.Request
		data map[string]string
	}
	type test struct {
		args      args
		checkFunc func(err error, data map[string]string) error
		name      string
	}
	tests := []test{
		func() test {
			buf := new(bytes.Buffer)
			buf.WriteString(`{"name":"vald"}`)
			r := httptest.NewRequest(http.MethodPost, "/", buf)

			return test{
				name: "return nil",
				args: args{
					r:    r,
					data: make(map[string]string, 1),
				},
				checkFunc: func(err error, data map[string]string) error {
					if err != nil {
						return errors.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if got, want := data, map[string]string{
						"name": "vald",
					}; !reflect.DeepEqual(got, want) {
						return errors.Errorf("read data not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),
		func() test {
			buf := new(bytes.Buffer)
			buf.WriteString(`10`)
			r := httptest.NewRequest(http.MethodPost, "/", buf)

			return test{
				name: "returns error because of faild to decode",
				args: args{
					r: r,
				},
				checkFunc: func(err error, data map[string]string) error {
					if err == nil {
						return errors.New("err is nil")
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			err := DecodeRequest(test.args.r, &test.args.data)
			if err := test.checkFunc(err, test.args.data); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		w     http.ResponseWriter
		r     *http.Request
		data  map[string]string
		logic func() (any, error)
	}
	type test struct {
		args      args
		checkFunc func(code int, err error, data map[string]string) error
		name      string
	}
	tests := []test{
		func() test {
			buf := new(bytes.Buffer)
			buf.WriteString(`{"name":"vald"}`)
			r := httptest.NewRequest(http.MethodPost, "/", buf)

			w := &httptest.ResponseRecorder{
				Body: new(bytes.Buffer),
			}

			data := make(map[string]string, 1)

			return test{
				name: "returns 200 and nil",
				args: args{
					r:    r,
					w:    w,
					data: data,
					logic: func() (any, error) {
						return "hello", nil
					},
				},
				checkFunc: func(code int, err error, data map[string]string) error {
					if err != nil {
						return errors.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if code != http.StatusOK {
						return errors.Errorf("code not equals. want: %v, got: %v", http.StatusOK, code)
					}

					if got, want := data, map[string]string{
						"name": "vald",
					}; !reflect.DeepEqual(got, want) {
						return errors.Errorf("data not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),
		func() test {
			buf := new(bytes.Buffer)
			buf.WriteString(`2`)
			r := httptest.NewRequest(http.MethodPost, "/", buf)

			return test{
				name: "returns error because of faild to decode",
				args: args{
					r: r,
				},
				checkFunc: func(code int, err error, data map[string]string) error {
					if err == nil {
						return errors.New("err is nil")
					}

					if code != http.StatusBadRequest {
						return errors.Errorf("code not equals. want: %v, got: %v", http.StatusBadRequest, code)
					}

					return nil
				},
			}
		}(),
		func() test {
			wantErr := errors.New("logic error")

			buf := new(bytes.Buffer)
			buf.WriteString(`{"name":"vald"}`)
			r := httptest.NewRequest(http.MethodPost, "/", buf)

			return test{
				name: "faild to logic",
				args: args{
					r:    r,
					data: make(map[string]string),
					logic: func() (any, error) {
						return nil, wantErr
					},
				},
				checkFunc: func(code int, err error, data map[string]string) error {
					if !errors.Is(err, wantErr) {
						return errors.Errorf("err not equals. want: %v, got: %v", wantErr, err)
					}

					if code != http.StatusInternalServerError {
						return errors.Errorf("code not equals. want: %v, got: %v", http.StatusInternalServerError, code)
					}

					if got, want := data, map[string]string{
						"name": "vald",
					}; !reflect.DeepEqual(got, want) {
						return errors.Errorf("data not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),
		{
			name: "returns error because of faild to encode",
			args: args{
				r: func() *http.Request {
					buf := new(bytes.Buffer)
					buf.WriteString(`{"name":"vald"}`)
					return httptest.NewRequest(http.MethodPost, "/", buf)
				}(),
				w:    new(httptest.ResponseRecorder),
				data: make(map[string]string),
				logic: func() (any, error) {
					return func() {}, nil
				},
			},
			checkFunc: func(code int, err error, data map[string]string) error {
				if code != http.StatusServiceUnavailable {
					return errors.Errorf("code not equals. want: %v, got: %v", http.StatusServiceUnavailable, code)
				}

				if got, want := data, map[string]string{
					"name": "vald",
				}; !reflect.DeepEqual(got, want) {
					return errors.Errorf("data not equals. want: %v, got: %v", want, got)
				}

				return nil
			},
		},
		func() test {
			w := new(httptest.ResponseRecorder)
			buf := new(bytes.Buffer)
			buf.WriteString(`{"name":"vald"}`)
			r := httptest.NewRequest(http.MethodPost, "/", buf)

			return test{
				name: "sends 503 to the client when encoding fails",
				args: args{
					r:    r,
					w:    w,
					data: make(map[string]string),
					logic: func() (any, error) {
						return func() {}, nil
					},
				},
				checkFunc: func(code int, err error, data map[string]string) error {
					if err == nil {
						return errors.New("err is nil")
					}
					if code != http.StatusServiceUnavailable {
						return errors.Errorf("code not equals. want: %v, got: %v", http.StatusServiceUnavailable, code)
					}
					if w.Code != http.StatusServiceUnavailable {
						return errors.Errorf("written code not equals. want: %v, got: %v", http.StatusServiceUnavailable, w.Code)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			code, err := Handler(test.args.w, test.args.r, &test.args.data, test.args.logic)
			if cerr := test.checkFunc(code, err, test.args.data); cerr != nil {
				t.Error(cerr)
			}
		})
	}
}

func TestHandler_logicErrorStatus(t *testing.T) {
	t.Parallel()
	type args struct {
		w     *httptest.ResponseRecorder
		r     *http.Request
		data  map[string]string
		logic func() (any, error)
	}
	type test struct {
		args      args
		checkFunc func(code int, err error, data map[string]string) error
		name      string
	}
	tests := []test{
		func() test {
			w := httptest.NewRecorder()
			wantErr := status.Error(codes.NotFound, "not found")
			buf := new(bytes.Buffer)
			buf.WriteString(`{"name":"vald"}`)
			r := httptest.NewRequest(http.MethodPost, "/", buf)

			return test{
				name: "returns 404 and leaves the body to the router when logic returns NotFound",
				args: args{
					r:    r,
					w:    w,
					data: make(map[string]string),
					logic: func() (any, error) {
						return nil, wantErr
					},
				},
				checkFunc: func(code int, err error, data map[string]string) error {
					if !errors.Is(err, wantErr) {
						return errors.Errorf("err not equals. want: %v, got: %v", wantErr, err)
					}
					if code != http.StatusNotFound {
						return errors.Errorf("code not equals. want: %v, got: %v", http.StatusNotFound, code)
					}
					if w.Body.Len() != 0 {
						return errors.Errorf("body written by Handler: %q", w.Body.String())
					}
					return nil
				},
			}
		}(),
		func() test {
			w := httptest.NewRecorder()
			wantErr := status.Error(codes.Unavailable, "unavailable")
			buf := new(bytes.Buffer)
			buf.WriteString(`{"name":"vald"}`)
			r := httptest.NewRequest(http.MethodPost, "/", buf)

			return test{
				name: "writes 503 problem json when logic returns Unavailable",
				args: args{
					r:    r,
					w:    w,
					data: make(map[string]string),
					logic: func() (any, error) {
						return nil, wantErr
					},
				},
				checkFunc: func(code int, err error, data map[string]string) error {
					if !errors.Is(err, wantErr) {
						return errors.Errorf("err not equals. want: %v, got: %v", wantErr, err)
					}
					if code != http.StatusServiceUnavailable {
						return errors.Errorf("code not equals. want: %v, got: %v", http.StatusServiceUnavailable, code)
					}
					if w.Code != http.StatusServiceUnavailable {
						return errors.Errorf("written code not equals. want: %v, got: %v", http.StatusServiceUnavailable, w.Code)
					}
					var got RFC7807Error
					if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
						return errors.Errorf("body is not problem json: %v, body: %q", err, w.Body.String())
					}
					if got.Status != http.StatusServiceUnavailable {
						return errors.Errorf("body status not equals. want: %v, got: %v", http.StatusServiceUnavailable, got.Status)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			code, err := Handler(test.args.w, test.args.r, &test.args.data, test.args.logic)
			if cerr := test.checkFunc(code, err, test.args.data); cerr != nil {
				tt.Error(cerr)
			}
		})
	}
}

func TestErrorHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		w    http.ResponseWriter
		err  error
		r    *http.Request
		msg  string
		code int
	}
	type test struct {
		checkFunc func(err error) error
		name      string
		args      args
	}
	tests := []test{
		func() test {
			rbuf := new(bytes.Buffer)
			rbuf.WriteString(`{"name":"vald"}`)
			r := httptest.NewRequest(http.MethodGet, "/", rbuf)

			wbuf := new(bytes.Buffer)
			w := &httptest.ResponseRecorder{
				Body: wbuf,
			}

			return test{
				name: "returns nil",
				args: args{
					r:    r,
					w:    w,
					msg:  "hello",
					code: http.StatusInternalServerError,
					err:  errors.New("faild to handler"),
				},
				checkFunc: func(err error) error {
					if err != nil {
						return errors.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if got, want := w.Header()[rest.ContentType],
						[]string{rest.ProblemJSON, rest.CharsetUTF8}; !reflect.DeepEqual(got, want) {
						return errors.Errorf("resp %v header not equals. want: %v, got: %v", rest.ContentType, want, got)
					}

					if got, want := w.Code, http.StatusInternalServerError; got != want {
						return errors.Errorf("response code not equals. want: %v, got: %v", http.StatusInternalServerError, got)
					}
					return nil
				},
			}
		}(),
		func() test {
			rbuf := new(bytes.Buffer)
			rbuf.WriteString(`{"name":"vald"}`)
			r := httptest.NewRequest(http.MethodPost, "/", rbuf)
			// The body has already been consumed and closed, as it is after Handler.
			if err := DecodeRequest(r, &map[string]any{}); err != nil {
				t.Fatalf("failed to decode request: %v", err)
			}

			w := &httptest.ResponseRecorder{
				Body: new(bytes.Buffer),
			}

			return test{
				name: "returns nil when the request body is already consumed",
				args: args{
					r:    r,
					w:    w,
					msg:  "hello",
					code: http.StatusNotFound,
					err:  errors.New("not found"),
				},
				checkFunc: func(err error) error {
					if err != nil {
						return errors.Errorf("err not equals. want: %v, got: %v", nil, err)
					}
					if got, want := w.Header()[rest.ContentType],
						[]string{rest.ProblemJSON, rest.CharsetUTF8}; !reflect.DeepEqual(got, want) {
						return errors.Errorf("resp %v header not equals. want: %v, got: %v", rest.ContentType, want, got)
					}
					if got, want := w.Code, http.StatusNotFound; got != want {
						return errors.Errorf("response code not equals. want: %v, got: %v", want, got)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			got := ErrorHandler(test.args.w, test.args.r, test.args.msg, test.args.code, test.args.err)
			if err := test.checkFunc(got); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_isConsumedBody(t *testing.T) {
	t.Parallel()
	decodeErr := func(t *testing.T, r *http.Request) error {
		t.Helper()
		return json.Decode(r.Body, &map[string]any{})
	}
	tests := []struct {
		err  func(t *testing.T) error
		name string
		want bool
	}{
		{
			name: "returns true for a body drained by DecodeRequest",
			err: func(t *testing.T) error {
				t.Helper()
				r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"vald"}`))
				if err := DecodeRequest(r, &map[string]any{}); err != nil {
					t.Fatalf("failed to decode request: %v", err)
				}
				return decodeErr(t, r)
			},
			want: true,
		},
		{
			name: "returns true for an empty body",
			err: func(t *testing.T) error {
				t.Helper()
				return decodeErr(t, httptest.NewRequest(http.MethodGet, "/", nil))
			},
			want: true,
		},
		{
			name: "returns true for a closed server request body",
			err: func(t *testing.T) error {
				t.Helper()
				return errors.Wrap(http.ErrBodyReadAfterClose, "read")
			},
			want: true,
		},
		{
			name: "returns false for a malformed body",
			err: func(t *testing.T) error {
				t.Helper()
				return decodeErr(t, httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":`)))
			},
			want: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.err(t)
			if err == nil {
				t.Fatal("decode error is nil")
			}
			if got := isConsumedBody(err); got != tc.want {
				t.Errorf("isConsumedBody(%v) = %v, want %v", err, got, tc.want)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW

func Test_httpStatusFromError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		err  error
		name string
		want int
	}{
		{name: "plain error is 500", err: errors.New("logic error"), want: http.StatusInternalServerError},
		{name: "Unknown is 500", err: status.Error(codes.Unknown, "x"), want: http.StatusInternalServerError},
		{name: "Internal is 500", err: status.Error(codes.Internal, "x"), want: http.StatusInternalServerError},
		{name: "DataLoss is 500", err: status.Error(codes.DataLoss, "x"), want: http.StatusInternalServerError},
		{name: "InvalidArgument is 400", err: status.Error(codes.InvalidArgument, "x"), want: http.StatusBadRequest},
		{name: "FailedPrecondition is 400", err: status.Error(codes.FailedPrecondition, "x"), want: http.StatusBadRequest},
		{name: "OutOfRange is 400", err: status.Error(codes.OutOfRange, "x"), want: http.StatusBadRequest},
		{name: "Unauthenticated is 401", err: status.Error(codes.Unauthenticated, "x"), want: http.StatusUnauthorized},
		{name: "PermissionDenied is 403", err: status.Error(codes.PermissionDenied, "x"), want: http.StatusForbidden},
		{name: "NotFound is 404", err: status.Error(codes.NotFound, "x"), want: http.StatusNotFound},
		{name: "AlreadyExists is 409", err: status.Error(codes.AlreadyExists, "x"), want: http.StatusConflict},
		{name: "Aborted is 409", err: status.Error(codes.Aborted, "x"), want: http.StatusConflict},
		{name: "ResourceExhausted is 429", err: status.Error(codes.ResourceExhausted, "x"), want: http.StatusTooManyRequests},
		{name: "Canceled is 499", err: status.Error(codes.Canceled, "x"), want: 499},
		{name: "Unimplemented is 501", err: status.Error(codes.Unimplemented, "x"), want: http.StatusNotImplemented},
		{name: "Unavailable is 503", err: status.Error(codes.Unavailable, "x"), want: http.StatusServiceUnavailable},
		{name: "DeadlineExceeded is 504", err: status.Error(codes.DeadlineExceeded, "x"), want: http.StatusGatewayTimeout},
		{name: "wrapped NotFound is 404", err: errors.Wrap(status.Error(codes.NotFound, "x"), "wrapped"), want: http.StatusNotFound},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			if got := httpStatusFromError(test.err); got != test.want {
				tt.Errorf("httpStatusFromError(%v) = %d, want %d", test.err, got, test.want)
			}
		})
	}
}
