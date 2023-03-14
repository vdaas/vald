// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package routing

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	goleak.VerifyTestMain(m)
}

func TestNew(t *testing.T) {
	t.Parallel()
	type test struct {
		name        string
		opts        []Option
		initialized bool
	}

	tests := []test{
		{
			name: "initialize success",
			opts: []Option{
				WithMiddleware(&middlewareMock{
					WrapFunc: func(r rest.Func) rest.Func {
						return r
					},
				}),
				WithRoutes(
					Route{},
				),
			},
			initialized: true,
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			got := New(test.opts...)
			if (got != nil) != test.initialized {
				t.Error("New() is wrong.")
			}
		})
	}
}

func TestRouting(t *testing.T) {
	t.Parallel()
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

			var cnt int32 = 0
			h := func(w http.ResponseWriter, req *http.Request) (code int, err error) {
				atomic.AddInt32(&cnt, 1)
				w.WriteHeader(http.StatusOK)
				return http.StatusOK, nil
			}

			return test{
				name: "serveHTTP is success when handler returns status ok",
				args: args{
					m: []string{
						http.MethodGet,
					},
					h: h,
				},
				checkFunc: func(hdr http.Handler) error {
					hdr.ServeHTTP(w, r)

					if atomic.LoadInt32(&cnt) != 1 {
						return errors.Errorf("call count is wrong. want: %v, got: %v", 1, cnt)
					}

					if got, want := w.Code, http.StatusOK; got != want {
						return errors.Errorf("status code not equals. want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),

		func() test {
			w := new(httptest.ResponseRecorder)
			r := httptest.NewRequest(http.MethodGet, "/", new(bytes.Buffer))

			return test{
				name: "serveHTTP is fail when handler returns invalid request method",
				checkFunc: func(hdr http.Handler) error {
					hdr.ServeHTTP(w, r)

					if got, want := w.Code, http.StatusMethodNotAllowed; got != want {
						return errors.Errorf("status code not equals. want: %v, got: %v", want, got)
					}
					return nil
				},
			}
		}(),

		func() test {
			w := new(httptest.ResponseRecorder)
			r := httptest.NewRequest(http.MethodGet, "/", new(bytes.Buffer))

			var cnt int32 = 0
			h := func(w http.ResponseWriter, req *http.Request) (code int, err error) {
				atomic.AddInt32(&cnt, 1)
				w.WriteHeader(http.StatusBadRequest)
				return http.StatusOK, errors.New("faild")
			}

			return test{
				name: "serveHTTP is fail when handler returns error",
				args: args{
					m: []string{
						http.MethodGet,
					},
					h: h,
				},
				checkFunc: func(hdr http.Handler) error {
					hdr.ServeHTTP(w, r)

					if atomic.LoadInt32(&cnt) != 1 {
						return errors.Errorf("call count is wrong. want: %v, got: %v", 1, cnt)
					}

					if got, want := w.Code, http.StatusBadRequest; got != want {
						return errors.Errorf("status code not equals. want: %v, got: %v", want, got)
					}
					return nil
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			// tt.Parallel()
			hdr := new(router).routing(test.args.name, test.args.path, test.args.m, test.args.h)
			if err := test.checkFunc(hdr); err != nil {
				t.Error(err)
			}
		})
	}
}
