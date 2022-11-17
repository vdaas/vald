// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/strings"
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
		data         interface{}
		status       int
		contentTypes []string
	}
	type test struct {
		name      string
		args      args
		checkFunc func(err error) error
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
		name      string
		args      args
		checkFunc func(err error, data map[string]string) error
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
		logic func() (interface{}, error)
	}
	type test struct {
		name      string
		args      args
		checkFunc func(code int, err error, data map[string]string) error
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
					logic: func() (interface{}, error) {
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
					logic: func() (interface{}, error) {
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
				logic: func() (interface{}, error) {
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
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			code, err := Handler(test.args.w, test.args.r, &test.args.data, test.args.logic)
			if err := test.checkFunc(code, err, test.args.data); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestErrorHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		w    http.ResponseWriter
		r    *http.Request
		msg  string
		code int
		err  error
	}
	type test struct {
		name      string
		args      args
		checkFunc func(err error) error
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
						return errors.Errorf("reso code not equals. want: %v, got: %v", http.StatusInternalServerError, got)
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

func TestDecodeResponse(t *testing.T) {
	t.Parallel()
	type args struct {
		res  *http.Response
		data interface{}
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "returns nil when response is nil",
				args: args{
					res:  nil,
					data: new(interface{}),
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			return test{
				name: "returns nil when response body is nil",
				args: args{
					res: &http.Response{
						Body: nil,
					},
					data: new(interface{}),
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			return test{
				name: "returns nil when the data to be decoded is nil",
				args: args{
					res: &http.Response{
						Body: io.NopCloser(new(bytes.Buffer)),
					},
					data: nil,
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			return test{
				name: "returns nil when the contents length is 0",
				args: args{
					res: &http.Response{
						Body:          io.NopCloser(new(bytes.Buffer)),
						ContentLength: 0,
					},
					data: new(interface{}),
				},
				want: want{
					err: nil,
				},
			}
		}(),
		func() test {
			return test{
				name: "returns json decode error when the response body is invalid",
				args: args{
					res: &http.Response{
						Body:          io.NopCloser(strings.NewReader("1+3i")),
						ContentLength: 2,
					},
					data: new(interface{}),
				},
				want: want{
					err: &strconv.NumError{
						Func: "ParseFloat",
						Num:  "1+3",
						Err:  strconv.ErrSyntax,
					},
				},
			}
		}(),
		func() test {
			var data int
			return test{
				name: "returns nil when the decode success",
				args: args{
					res: &http.Response{
						Body:          io.NopCloser(strings.NewReader("1")),
						ContentLength: 1,
					},
					data: &data,
				},
				checkFunc: func(w want, got error) error {
					if err := defaultCheckFunc(w, got); err != nil {
						return err
					}

					if want, got := 1, data; want != got {
						return errors.Errorf("data want: %d, but got: %d", want, got)
					}

					return nil
				},
				want: want{
					err: nil,
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

			err := DecodeResponse(test.args.res, test.args.data)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestEncodeRequest(t *testing.T) {
	t.Parallel()
	type args struct {
		req          *http.Request
		data         interface{}
		contentTypes []string
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if w.err != nil && err != nil && !strings.HasPrefix(err.Error(), w.err.Error()) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			val := 1 + 3i
			return test{
				name: "returns json encode error when the json encode fails",
				args: args{
					req:  new(http.Request),
					data: val,
				},
				want: want{
					err: errors.New("json: unsupported type:"),
				},
			}
		}(),
		func() test {
			req := &http.Request{
				Header: http.Header{},
			}
			return test{
				name: "returns nil when json encode success",
				args: args{
					req: req,
					contentTypes: []string{
						"application/json",
					},
					data: 1,
				},
				checkFunc: func(w want, got error) error {
					if err := defaultCheckFunc(w, got); err != nil {
						return err
					}

					if len(req.Header) != 1 {
						return errors.Errorf("header length is wrong. want: %v, but got: %d", 1, len(req.Header))
					}

					gotHeaders, ok := req.Header[rest.ContentType]
					if !ok {
						return errors.Errorf("header not found. key: %s", rest.ContentType)
					}

					if len(gotHeaders) != 1 {
						return errors.Errorf("header value length is wrong. key:%s want: %d, but got: %d", rest.ContentType, 1, len(gotHeaders))
					}

					if want, got := "application/json", gotHeaders[0]; want != got {
						return errors.Errorf("header value is wrong. want: %s, but got: %s", want, got)
					}

					if want, got := int64(2), req.ContentLength; want != got {
						return errors.Errorf("content length is wrong. want: %v, but got: %d", want, got)
					}

					if req.Body == nil {
						return errors.New("Body is nil")
					}

					return nil
				},
				want: want{
					err: nil,
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

			err := EncodeRequest(test.args.req, test.args.data, test.args.contentTypes...)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestRequest(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context
		method  string
		url     string
		payloyd interface{}
		data    interface{}
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "returns generation of request error when method is invalid",
				args: args{
					ctx:     context.Background(),
					method:  "@",
					url:     "/",
					payloyd: nil,
					data:    nil,
				},
				want: want{
					err: errors.Errorf("net/http: invalid method %q", "@"),
				},
			}
		}(),
		func() test {
			return test{
				name: "returns json encode error when the request json encoding fails",
				args: args{
					ctx:     context.Background(),
					method:  "POST",
					url:     "/",
					payloyd: 1 + 3i,
					data:    new(interface{}),
				},
				checkFunc: func(w want, err error) error {
					if w.err != nil && err != nil && !strings.HasPrefix(err.Error(), w.err.Error()) {
						return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
					}
					return nil
				},
				want: want{
					err: errors.New("json: unsupported type:"),
				},
			}
		}(),
		func() test {
			return test{
				name: "returns http request error when sending http request fails",
				args: args{
					ctx:     context.Background(),
					method:  "POST",
					url:     "/",
					payloyd: "1",
					data:    new(interface{}),
				},
				want: want{
					err: &url.Error{
						Op:  "Post",
						URL: "/",
						Err: errors.New(`unsupported protocol scheme ""`),
					},
				},
			}
		}(),
		func() test {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				_, _ = w.Write([]byte("\"1\""))
			}))
			var got string
			return test{
				name: "returns nil when no error occurs internally",
				args: args{
					ctx:     context.Background(),
					method:  "POST",
					url:     srv.URL,
					payloyd: "1",
					data:    &got,
				},
				want: want{
					err: nil,
				},
				checkFunc: func(w want, err error) error {
					if err := defaultCheckFunc(w, err); err != nil {
						return err
					}
					if want, got := "1", got; want != got {
						return errors.Errorf("decoded data: want %s, but got: %v", want, got)
					}
					return nil
				},
				afterFunc: func(args) {
					srv.Close()
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

			err := Request(test.args.ctx, test.args.method, test.args.url, test.args.payloyd, test.args.data)
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
