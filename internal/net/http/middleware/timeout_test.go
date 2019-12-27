package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	errgroup "github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/http/rest"
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
						return fmt.Errorf("err is not nil. err: %v", err)
					}

					if code != http.StatusOK {
						return fmt.Errorf("code is not equals. want: %v, got: %v", http.StatusOK, code)
					}

					if cnt != 1 {
						return fmt.Errorf("called cnt is equals. want: %v, got: %v", 1, cnt)
					}

					return nil
				},
			}
		}(),
		func() test {
			wantErr := fmt.Errorf("faild")

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
						return fmt.Errorf("err not equals. want: %v, got: %v", wantErr, err)
					}

					if code != http.StatusInternalServerError {
						return fmt.Errorf("code is not equals. want: %v, got: %v", http.StatusInternalServerError, code)
					}

					if cnt != 1 {
						return fmt.Errorf("called cnt is equals. want: %v, got: %v", 1, cnt)
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
						return fmt.Errorf("err is nil")
					}

					if code != http.StatusRequestTimeout {
						return fmt.Errorf("code is not equals. want: %v, got: %v", http.StatusRequestTimeout, code)
					}

					if !strings.Contains(err.Error(), "handler timeout") {
						return fmt.Errorf("err string no contains word of `handler timeout`")
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
