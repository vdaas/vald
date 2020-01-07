package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerToRestFunc(t *testing.T) {
	type test struct {
		name      string
		hfn       http.HandlerFunc
		checkFunc func(Func) error
	}

	tests := []test{
		func() test {
			cnt := 0

			hfn := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				cnt++
			})

			return test{
				name: "returns 200 status code",
				hfn:  hfn,
				checkFunc: func(fn Func) error {
					code, err := fn(httptest.NewRecorder(), new(http.Request))
					if err != nil {
						return fmt.Errorf("err is not nil. err: %v", err)
					}

					if code != http.StatusOK {
						return fmt.Errorf("status code is wrong. want: %v, got: %v", http.StatusOK, code)
					}

					if cnt != 1 {
						return fmt.Errorf("called count is wrong. want: %v, got: %v", 1, cnt)
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := HandlerToRestFunc(tt.hfn)
			if err := tt.checkFunc(fn); err != nil {
				t.Error(err)
			}
		})
	}
}
