package rest

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHandlerToRestFunc(t *testing.T) {
	type responseWriter struct {
		http.ResponseWriter
	}

	type args struct {
		hfn http.HandlerFunc
	}

	type test struct {
		name      string
		args      args
		checkFunc func(f Func) error
	}

	tests := []test{
		func() test {
			cnt := 0

			hfn := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				cnt++
			})

			return test{
				name: "success",
				args: args{
					hfn: hfn,
				},
				checkFunc: func(fn Func) error {
					code, err := fn(new(responseWriter), new(http.Request))
					if err != nil {
						return fmt.Errorf("err is not nil: %v", err)
					}

					if code != http.StatusOK {
						return fmt.Errorf("code is wrong. want: %v, got: %v", http.StatusOK, code)
					}

					if cnt != 1 {
						return fmt.Errorf("internal call count is wrong. want: %v, got: %v", 1, cnt)
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := HandlerToRestFunc(tt.args.hfn)
			if err := tt.checkFunc(fn); err != nil {
				t.Error(err)
			}
		})
	}
}
