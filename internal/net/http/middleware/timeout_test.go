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
package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func TestNewTimeout(t *testing.T) {
	t.Parallel()
	type test struct {
		want Wrapper
		name string
	}
	tests := []test{
		{
			name: "create object",
			want: &timeout{
				dur: 3 * time.Second,
				eg:  errgroup.Get(),
			},
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			got := NewTimeout()
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("not equals. want: %v, got: %v", test.want, got)
			}
		})
	}
}

func Test_timeout_Wrap(t *testing.T) {
	t.Parallel()
	type args struct {
		h rest.Func
		w http.ResponseWriter
		r *http.Request
	}
	type field struct {
		dur time.Duration
	}
	type test struct {
		args      args
		checkFunc func(code int, err error) error
		name      string
		field     field
	}
	tests := []test{
		func() test {
			var cnt uint64 = 0
			h := func(w http.ResponseWriter, req *http.Request) (code int, err error) {
				atomic.AddUint64(&cnt, 1)
				return http.StatusOK, nil
			}

			return test{
				name: "internally called handler returns nil",
				args: args{
					h: h,
					w: new(httptest.ResponseRecorder),
					r: new(http.Request),
				},
				field: field{
					dur: 5 * time.Second,
				},
				checkFunc: func(code int, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil. err: %v", err)
					}
					if code != http.StatusOK {
						return errors.Errorf("code is not equals. want: %v, got: %v", http.StatusOK, code)
					}
					if atomic.LoadUint64(&cnt) != 1 {
						return errors.Errorf("called cnt is equals. want: %v, got: %v", 1, cnt)
					}
					return nil
				},
			}
		}(),
		func() test {
			wantErr := errors.New("failed")
			var cnt uint64 = 0
			h := func(w http.ResponseWriter, req *http.Request) (code int, err error) {
				atomic.AddUint64(&cnt, 1)
				return http.StatusInternalServerError, wantErr
			}

			return test{
				name: "internally called handler returns error",
				args: args{
					h: h,
					w: new(httptest.ResponseRecorder),
					r: new(http.Request),
				},
				field: field{
					dur: 2 * time.Second,
				},
				checkFunc: func(code int, err error) error {
					if !errors.Is(err, wantErr) {
						return errors.Errorf("err not equals. want: %v, got: %v", wantErr, err)
					}
					if code != http.StatusInternalServerError {
						return errors.Errorf("code is not equals. want: %v, got: %v", http.StatusInternalServerError, code)
					}
					if atomic.LoadUint64(&cnt) != 1 {
						return errors.Errorf("called cnt is equals. want: %v, got: %v", 1, cnt)
					}
					return nil
				},
			}
		}(),
		{
			name: "timeout processing of internally called handler",
			args: args{
				h: func(w http.ResponseWriter, req *http.Request) (code int, err error) {
					return http.StatusOK, nil
				},
				w: new(httptest.ResponseRecorder),
				r: new(http.Request),
			},
			field: field{
				dur: 1, // set to extermemly small value to let it timeout
			},
			checkFunc: func(code int, err error) error {
				if err == nil {
					return errors.Errorf("err is nil")
				}
				if code != http.StatusRequestTimeout {
					return errors.Errorf("code is not equals. want: %v, got: %v", http.StatusRequestTimeout, code)
				}
				if !strings.Contains(err.Error(), "handler timeout") {
					return errors.Errorf("err string no contains word of `handler timeout`")
				}

				return nil
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			eg, _ := errgroup.New(tt.Context())
			to := &timeout{
				dur: test.field.dur,
				eg:  eg,
			}
			code, err := to.Wrap(test.args.h)(test.args.w, test.args.r)
			if err := test.checkFunc(code, err); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_timeout_Wrap_discardsLateWrite(t *testing.T) {
	t.Parallel()
	release := make(chan struct{})
	written := make(chan error, 1)
	h := func(w http.ResponseWriter, _ *http.Request) (int, error) {
		select {
		case <-release:
		case <-t.Context().Done():
			written <- t.Context().Err()
			return http.StatusOK, nil
		}
		w.Header().Set("X-Late", "1")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("late"))
		written <- err
		return http.StatusOK, nil
	}
	eg, _ := errgroup.New(t.Context())
	to := &timeout{dur: time.Millisecond, eg: eg}
	rec := httptest.NewRecorder()
	// the handler is blocked until release is closed, so Wrap can only return by the timeout.
	code, err := to.Wrap(h)(rec, httptest.NewRequest(http.MethodPost, "/", nil))
	if code != http.StatusRequestTimeout || err == nil {
		t.Fatalf("want (%d, non-nil error), got (%d, %v)", http.StatusRequestTimeout, code, err)
	}
	close(release)
	select {
	case werr := <-written:
		if !errors.Is(werr, http.ErrHandlerTimeout) {
			t.Errorf("late write error want: %v, got: %v", http.ErrHandlerTimeout, werr)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("handler did not finish")
	}
	if rec.Body.Len() != 0 || rec.Header().Get("X-Late") != "" {
		t.Errorf("late handler output reached the real writer: header=%v body=%q", rec.Header(), rec.Body.String())
	}
}

func Test_timeout_Wrap_flushesResponse(t *testing.T) {
	t.Parallel()
	h := func(w http.ResponseWriter, _ *http.Request) (int, error) {
		w.Header().Set("X-Test", "1")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte("ok"))
		return http.StatusOK, err
	}
	eg, _ := errgroup.New(t.Context())
	to := &timeout{dur: time.Minute, eg: eg}
	rec := httptest.NewRecorder()
	code, err := to.Wrap(h)(rec, httptest.NewRequest(http.MethodPost, "/", nil))
	if code != http.StatusOK || err != nil {
		t.Fatalf("want (%d, nil), got (%d, %v)", http.StatusOK, code, err)
	}
	if rec.Code != http.StatusCreated || rec.Body.String() != "ok" || rec.Header().Get("X-Test") != "1" {
		t.Errorf("response not flushed: code=%d header=%v body=%q", rec.Code, rec.Header(), rec.Body.String())
	}
}

// Test_timeout_Wrap_inline runs the handler inline (errgroup limit 1) so that its result is
// already sent when Wrap selects, which makes the outcome independent of goroutine scheduling.
func Test_timeout_Wrap_inline(t *testing.T) {
	t.Parallel()
	wantErr := errors.New("failed")
	type test struct {
		h        func(cancel context.CancelFunc) rest.Func
		name     string
		wantMsg  string
		wantCode int
	}
	tests := []test{
		{
			name: "returns timeout when the context is done before the failed handler returns",
			h: func(cancel context.CancelFunc) rest.Func {
				return func(w http.ResponseWriter, _ *http.Request) (int, error) {
					cancel()
					return http.StatusInternalServerError, wantErr
				}
			},
			wantCode: http.StatusRequestTimeout,
			wantMsg:  "handler timeout",
		},
		{
			name: "returns timeout when the context is done before the succeeded handler returns",
			h: func(cancel context.CancelFunc) rest.Func {
				return func(w http.ResponseWriter, _ *http.Request) (int, error) {
					cancel()
					return http.StatusOK, nil
				}
			},
			wantCode: http.StatusRequestTimeout,
			wantMsg:  "handler timeout",
		},
		{
			name: "returns internal server error when the handler panics",
			h: func(context.CancelFunc) rest.Func {
				return func(w http.ResponseWriter, _ *http.Request) (int, error) {
					panic("handler panic")
				}
			},
			wantCode: http.StatusInternalServerError,
			wantMsg:  "handler panic",
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			eg, _ := errgroup.New(tt.Context())
			eg.SetLimit(1)
			to := &timeout{dur: time.Minute, eg: eg}
			ctx, cancel := context.WithCancel(tt.Context())
			defer cancel()
			r := httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx)
			// repeat because the old implementation returned a random result
			for range 100 {
				code, err := to.Wrap(test.h(cancel))(httptest.NewRecorder(), r)
				if code != test.wantCode || err == nil || !strings.Contains(err.Error(), test.wantMsg) {
					tt.Fatalf("want (%d, error containing %q), got (%d, %v)", test.wantCode, test.wantMsg, code, err)
				}
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
