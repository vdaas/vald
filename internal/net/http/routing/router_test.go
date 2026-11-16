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
package routing

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net/http/middleware"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/sync"
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
		h    rest.Func
		name string
		path string
		m    []string
	}

	type test struct {
		checkFunc func(http.Handler) error
		name      string
		args      args
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

		func() test {
			w := httptest.NewRecorder()
			release := make(chan struct{})
			r := httptest.NewRequest(http.MethodPost, "/", &blockingReader{release: release})
			h := func(http.ResponseWriter, *http.Request) (int, error) {
				return http.StatusRequestTimeout, errors.New("handler timeout")
			}

			return test{
				name: "serveHTTP does not read the request body after handler timeout",
				args: args{
					m: []string{
						http.MethodPost,
					},
					h: h,
				},
				checkFunc: func(hdr http.Handler) error {
					defer close(release)
					done := make(chan struct{})
					go func() {
						defer close(done)
						hdr.ServeHTTP(w, r)
					}()
					select {
					case <-done:
					case <-time.After(10 * time.Second):
						return errors.New("serveHTTP blocked on reading the request body")
					}
					if got, want := w.Code, http.StatusRequestTimeout; got != want {
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

// blockingReader blocks every Read until release is closed.
type blockingReader struct {
	release chan struct{}
}

func (b *blockingReader) Read([]byte) (int, error) {
	<-b.release
	return 0, io.EOF
}

// NOT IMPLEMENTED BELOW
//
// func Test_router_routing(t *testing.T) {
// 	type args struct {
// 		name string
// 		in1  string
// 		m    []string
// 		h    rest.Func
// 	}
// 	type fields struct {
// 		middlewares []middleware.Wrapper
// 		routes      []Route
// 	}
// 	type want struct {
// 		want http.Handler
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, http.Handler) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got http.Handler) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           name:"",
// 		           in1:"",
// 		           m:nil,
// 		           h:nil,
// 		       },
// 		       fields: fields {
// 		           middlewares:nil,
// 		           routes:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           name:"",
// 		           in1:"",
// 		           m:nil,
// 		           h:nil,
// 		           },
// 		           fields: fields {
// 		           middlewares:nil,
// 		           routes:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			r := &router{
// 				middlewares: test.fields.middlewares,
// 				routes:      test.fields.routes,
// 			}
//
// 			got := r.routing(test.args.name, test.args.in1, test.args.m, test.args.h)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

func TestNew_concurrentDoesNotRaceOnDefaultTransport(t *testing.T) {
	// Not parallel: resets the package-level once and swaps http.DefaultTransport.
	orig := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = orig
		tuneOnce = sync.Once{}
	})

	tr := &http.Transport{}
	http.DefaultTransport = tr
	tuneOnce = sync.Once{}

	const n = 16
	var wg sync.WaitGroup
	done := make(chan struct{})
	for range n {
		wg.Go(func() {
			if New() == nil {
				t.Error("New returned nil")
			}
		})
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for concurrent New calls")
	}

	// All New calls have returned, so reading is ordered after the writes.
	if got, want := tr.MaxIdleConnsPerHost, defaultMaxIdleConnsPerHost; got != want {
		t.Errorf("MaxIdleConnsPerHost = %d, want %d", got, want)
	}
}

type notTransport struct{}

func (notTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("unused")
}

func TestNew_defaultTransportNotHTTPTransport(t *testing.T) {
	// Not parallel: resets the package-level once and swaps http.DefaultTransport.
	orig := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = orig
		tuneOnce = sync.Once{}
	})

	http.DefaultTransport = notTransport{}
	tuneOnce = sync.Once{}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("New panicked when http.DefaultTransport is not *http.Transport: %v", r)
		}
	}()
	if New() == nil {
		t.Error("New returned nil")
	}
}

func TestWithRoutesMiddlewares_doNotAliasCallerSlice(t *testing.T) {
	t.Parallel()

	mws := make([]middleware.Wrapper, 1, 4)
	mws[0] = new(middlewareMock)
	r := new(router)
	WithMiddlewares(mws...)(r)
	if len(r.middlewares) != 1 || &r.middlewares[0] == &mws[0] {
		t.Error("WithMiddlewares must copy the caller's slice")
	}

	routes := make([]Route, 1, 4)
	r = new(router)
	WithRoutes(routes...)(r)
	if len(r.routes) != 1 || &r.routes[0] == &routes[0] {
		t.Error("WithRoutes must copy the caller's slice")
	}
}
