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
package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/rest"
)

func TestNewTimeout(t *testing.T) {
	t.Parallel()
	type test struct {
		name string
		want Wrapper
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
		name      string
		args      args
		field     field
		checkFunc func(code int, err error) error
	}
	tests := []test{
		func() test {
			var cnt int32 = 0
			h := func(w http.ResponseWriter, req *http.Request) (code int, err error) {
				atomic.AddInt32(&cnt, 1)
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
					if atomic.LoadInt32(&cnt) != 1 {
						return errors.Errorf("called cnt is equals. want: %v, got: %v", 1, cnt)
					}
					return nil
				},
			}
		}(),
		func() test {
			wantErr := errors.New("failed")
			var cnt int32 = 0
			h := func(w http.ResponseWriter, req *http.Request) (code int, err error) {
				atomic.AddInt32(&cnt, 1)
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
					if atomic.LoadInt32(&cnt) != 1 {
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

			eg, _ := errgroup.New(context.Background())
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
