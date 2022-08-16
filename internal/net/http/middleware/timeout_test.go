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
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/rest"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestNewTimeout(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTimeout()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		h rest.Func
		w http.ResponseWriter
		r *http.Request
	}

	type field struct {
		dur time.Duration
		eg  errgroup.Group
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(code int, err error) error
	}

	tests := []test{
		func() test {
			var cnt int
			h := func(w http.ResponseWriter, req *http.Request) (code int, err error) {
				cnt++
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
					dur: 2 * time.Second,
					eg:  errgroup.Get(),
				},
				checkFunc: func(code int, err error) error {
					if err != nil {
						return errors.Errorf("err is not nil. err: %v", err)
					}

					if code != http.StatusOK {
						return errors.Errorf("code is not equals. want: %v, got: %v", http.StatusOK, code)
					}

					if cnt != 1 {
						return errors.Errorf("called cnt is equals. want: %v, got: %v", 1, cnt)
					}

					return nil
				},
			}
		}(),
		func() test {
			wantErr := errors.New("faild")

			var cnt int
			h := func(w http.ResponseWriter, req *http.Request) (code int, err error) {
				cnt++
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
					eg:  errgroup.Get(),
				},
				checkFunc: func(code int, err error) error {
					if !errors.Is(err, wantErr) {
						return errors.Errorf("err not equals. want: %v, got: %v", wantErr, err)
					}

					if code != http.StatusInternalServerError {
						return errors.Errorf("code is not equals. want: %v, got: %v", http.StatusInternalServerError, code)
					}

					if cnt != 1 {
						return errors.Errorf("called cnt is equals. want: %v, got: %v", 1, cnt)
					}

					return nil
				},
			}
		}(),
		func() test {
			h := func(w http.ResponseWriter, req *http.Request) (code int, err error) {
				time.Sleep(10 * time.Second)
				return http.StatusOK, nil
			}

			return test{
				name: "timeout processing of internally called handler",
				args: args{
					h: h,
					w: new(httptest.ResponseRecorder),
					r: new(http.Request),
				},
				field: field{
					dur: 1 * time.Second,
					eg:  errgroup.Get(),
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
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			to := &timeout{
				dur: tt.field.dur,
				eg:  tt.field.eg,
			}

			code, err := to.Wrap(tt.args.h)(tt.args.w, tt.args.r)
			if err := tt.checkFunc(code, err); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_timeout_Wrap(t *testing.T) {
	type args struct {
		h rest.Func
	}
	type fields struct {
		dur time.Duration
		eg  errgroup.Group
	}
	type want struct {
		want rest.Func
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, rest.Func) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got rest.Func) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           h: nil,
		       },
		       fields: fields {
		           dur: nil,
		           eg: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           h: nil,
		           },
		           fields: fields {
		           dur: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
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
			t := &timeout{
				dur: test.fields.dur,
				eg:  test.fields.eg,
			}

			got := t.Wrap(test.args.h)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
