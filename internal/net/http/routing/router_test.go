package routing

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/http/rest"
)

func TestNew(t *testing.T) {
	type test struct {
		name        string
		opts        []Option
		initialized bool
	}

	tests := []test{
		func() test {
			mw := &middlewareMock{
				WrapFunc: func(r rest.Func) rest.Func {
					return r
				},
			}

			return test{
				name: "initialize success",
				opts: []Option{
					WithMiddleware(mw),
					WithRoutes(
						Route{},
					),
				},
				initialized: true,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.opts...)
			if (got != nil) != tt.initialized {
				t.Error("New() is wrong.")
			}
		})
	}
}

func TestRouting(t *testing.T) {
	type args struct {
		name string
		path string
		m    []string
		h    rest.Func
	}

	type test struct {
		name      string
		args      args
		checkFunc func(http.Handler) error
	}

	tests := []test{
		func() test {
			w := new(httptest.ResponseRecorder)
			r := httptest.NewRequest(http.MethodGet, "/", new(bytes.Buffer))

			cnt := 0
			h := func(w http.ResponseWriter, req *http.Request) (code int, err error) {
				cnt++
				w.WriteHeader(http.StatusOK)
				return http.StatusOK, nil
			}

			return test{
				name: "status ok",
				args: args{
					m: []string{
						http.MethodGet,
					},
					h: h,
				},
				checkFunc: func(hdr http.Handler) error {
					hdr.ServeHTTP(w, r)

					if cnt != 1 {
						return fmt.Errorf("call count is wrong. want: %v, got: %v", 1, cnt)
					}

					if got, want := w.Code, http.StatusOK; got != want {
						return fmt.Errorf("status code not equals. want: %v, got: %v", want, got)
					}
					return nil
				},
			}
		}(),

		func() test {
			w := new(httptest.ResponseRecorder)
			r := httptest.NewRequest(http.MethodGet, "/", new(bytes.Buffer))

			return test{
				name: "invalid request method",
				checkFunc: func(hdr http.Handler) error {
					hdr.ServeHTTP(w, r)

					if got, want := w.Code, http.StatusMethodNotAllowed; got != want {
						return fmt.Errorf("status code not equals. want: %v, got: %v", want, got)
					}
					return nil
				},
			}
		}(),

		func() test {
			w := new(httptest.ResponseRecorder)
			r := httptest.NewRequest(http.MethodGet, "/", new(bytes.Buffer))

			cnt := 0
			h := func(w http.ResponseWriter, req *http.Request) (code int, err error) {
				cnt++
				w.WriteHeader(http.StatusBadRequest)
				return http.StatusOK, fmt.Errorf("faild")
			}

			return test{
				name: "internal call error",
				args: args{
					m: []string{
						http.MethodGet,
					},
					h: h,
				},
				checkFunc: func(hdr http.Handler) error {
					hdr.ServeHTTP(w, r)

					if cnt != 1 {
						return fmt.Errorf("call count is wrong. want: %v, got: %v", 1, cnt)
					}

					if got, want := w.Code, http.StatusBadRequest; got != want {
						return fmt.Errorf("status code not equals. want: %v, got: %v", want, got)
					}
					return nil
				},
			}
		}(),
	}

	log.Init(log.DefaultGlg())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdr := new(router).routing(tt.args.name, tt.args.path, tt.args.m, tt.args.h)
			if err := tt.checkFunc(hdr); err != nil {
				t.Error(err)
			}
		})
	}
}
