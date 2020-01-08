package json

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/http/rest"
)

func TestEncode(t *testing.T) {
	type args struct {
		w    io.Writer
		data interface{}
	}

	type test struct {
		name      string
		args      args
		checkFunc func(err error) error
	}

	tests := []test{
		func() test {
			buf := new(bytes.Buffer)
			data := map[string]string{
				"name": "vald",
			}

			return test{
				name: "returns nil",
				args: args{
					w:    buf,
					data: data,
				},
				checkFunc: func(err error) error {
					if err != nil {
						return fmt.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if got, want := buf.String(), "{\"name\":\"vald\"}\n"; got != want {
						return fmt.Errorf("output data not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),

		{
			name: "returns error",
			args: args{
				w:    new(bytes.Buffer),
				data: make(chan struct{}),
			},
			checkFunc: func(err error) error {
				if err == nil {
					return fmt.Errorf("err is nil")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Encode(tt.args.w, tt.args.data)
			if err := tt.checkFunc(err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		r    io.Reader
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

			return test{
				name: "returns nil",
				args: args{
					r:    buf,
					data: make(map[string]string, 1),
				},
				checkFunc: func(err error, data map[string]string) error {
					if err != nil {
						return fmt.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if got, want := data, map[string]string{
						"name": "vald",
					}; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("read data not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),

		func() test {
			buf := new(bytes.Buffer)
			buf.WriteString(`1`)

			wantData := make(map[string]string, 1)

			return test{
				name: "returns error",
				args: args{
					r:    buf,
					data: wantData,
				},
				checkFunc: func(err error, data map[string]string) error {
					if err == nil {
						return fmt.Errorf("err is nil")
					}

					if !reflect.DeepEqual(data, wantData) {
						return fmt.Errorf("data not equals. want: %v, got: %v", wantData, data)
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Decode(tt.args.r, &tt.args.data)
			if err := tt.checkFunc(err, tt.args.data); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMarshalIndent(t *testing.T) {
	type args struct {
		data interface{}
		pref string
		ind  string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(data []byte, err error) error
	}

	tests := []test{
		func() test {
			return test{
				name: "returns data and nil",
				args: args{
					data: map[string]string{
						"name": "vald",
					},
					pref: "",
					ind:  "",
				},
				checkFunc: func(data []byte, err error) error {
					if err != nil {
						return fmt.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if got, want := data, []byte(`{"name":"vald"}`); !reflect.DeepEqual(got, want) {
						return fmt.Errorf("data not equals. want: %v, got: %v", string(got), string(want))
					}

					return nil
				},
			}
		}(),

		{
			name: "returns error",
			args: args{
				data: make(chan struct{}),
				pref: "",
				ind:  "",
			},
			checkFunc: func(data []byte, err error) error {
				if err == nil {
					return fmt.Errorf("err is nil")
				}

				if len(data) != 0 {
					return fmt.Errorf("data is not empty")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalIndent(tt.args.data, tt.args.pref, tt.args.ind)
			if err := tt.checkFunc(data, err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestEncodeResponse(t *testing.T) {
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
						return fmt.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if got, want := w.Header().Get(rest.ContentType), "application/json"; got != want {
						return fmt.Errorf("content-type not equals. want: %v, got: %v", want, got)
					}

					if got, want := w.Code, 200; got != want {
						return fmt.Errorf("code not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),
		func() test {
			w := new(httptest.ResponseRecorder)

			return test{
				name: "returns error",
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
						return fmt.Errorf("err is nil")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EncodeResponse(tt.args.w, tt.args.data, tt.args.status, tt.args.contentTypes...)
			if err := tt.checkFunc(err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestDecodeRequest(t *testing.T) {
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
						return fmt.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if got, want := data, map[string]string{
						"name": "vald",
					}; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("read data not equals. want: %v, got: %v", want, got)
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
				name: "faild to decode",
				args: args{
					r: r,
				},
				checkFunc: func(err error, data map[string]string) error {
					if err == nil {
						return fmt.Errorf("err is nil.")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DecodeRequest(tt.args.r, &tt.args.data)
			if err := tt.checkFunc(err, tt.args.data); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestHandler(t *testing.T) {
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
		want      int
		wantErr   error
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
				name: "returns 200 status code and nil",
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
						return fmt.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if code != http.StatusOK {
						return fmt.Errorf("code not equals. want: %v, got: %v", http.StatusOK, code)
					}

					if got, want := data, map[string]string{
						"name": "vald",
					}; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("data not equals. want: %v, got: %v", want, got)
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
				name: "faild to decode request body",
				args: args{
					r: r,
				},
				checkFunc: func(code int, err error, data map[string]string) error {
					if err == nil {
						return fmt.Errorf("err is nil")
					}

					if code != http.StatusBadRequest {
						return fmt.Errorf("code not equals. want: %v, got: %v", http.StatusBadRequest, code)
					}

					return nil
				},
			}
		}(),

		func() test {
			wantErr := fmt.Errorf("logic error")
			return test{
				name: "faild to logic",
				args: args{
					r: func() *http.Request {
						buf := new(bytes.Buffer)
						buf.WriteString(`{"name":"vald"}`)
						return httptest.NewRequest(http.MethodPost, "/", buf)
					}(),
					data: make(map[string]string),
					logic: func() (interface{}, error) {
						return nil, wantErr
					},
				},
				checkFunc: func(code int, err error, data map[string]string) error {
					if !errors.Is(err, wantErr) {
						return fmt.Errorf("err not equals. want: %v, got: %v", wantErr, err)
					}

					if code != http.StatusInternalServerError {
						return fmt.Errorf("code not equals. want: %v, got: %v", http.StatusInternalServerError, code)
					}

					if got, want := data, map[string]string{
						"name": "vald",
					}; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("data not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),

		func() test {
			return test{
				name: "faild to encode",
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
						return fmt.Errorf("code not equals. want: %v, got: %v", http.StatusServiceUnavailable, code)
					}

					if got, want := data, map[string]string{
						"name": "vald",
					}; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("data not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := Handler(tt.args.w, tt.args.r, &tt.args.data, tt.args.logic)
			if err := tt.checkFunc(code, err, tt.args.data); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestErrorHandler(t *testing.T) {
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
					err:  fmt.Errorf("faild to handler"),
				},
				checkFunc: func(err error) error {
					if err != nil {
						return fmt.Errorf("err not equals. want: %v, got: %v", nil, err)
					}

					if got, want := w.Header()[rest.ContentType],
						[]string{rest.ProblemJSON, rest.CharsetUTF8}; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("resp %v header not equals. want: %v, got: %v", rest.ContentType, want, got)
					}

					if got, want := w.Code, http.StatusInternalServerError; got != want {
						return fmt.Errorf("reso code not equals. want: %v, got: %v", http.StatusInternalServerError, got)
					}
					return nil
				},
			}
		}(),
	}

	log.Init(log.DefaultGlg())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorHandler(tt.args.w, tt.args.r, tt.args.msg, tt.args.code, tt.args.err)
			if err := tt.checkFunc(got); err != nil {
				t.Error(err)
			}
		})
	}
}
