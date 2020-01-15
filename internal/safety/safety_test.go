package safety

import (
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func TestRecoverFunc(t *testing.T) {
	type test struct {
		name       string
		fn         func() error
		runtimeErr bool
		want       error
	}

	tests := []test{
		{
			name: "runtime error",
			fn: func() error {
				_ = []string{}[10]
				return nil
			},
			runtimeErr: true,
			want:       errors.New("system paniced caused by runtime error: runtime error: index out of range [10] with length 0"),
		},

		{
			name: "panic string",
			fn: func() error {
				panic("panic")
			},
			want: errors.New("panic recovered: panic"),
		},

		{
			name: "panic error",
			fn: func() error {
				panic(fmt.Errorf("error"))
			},
			want: errors.New("error"),
		},

		{
			name: "default case panic",
			fn: func() error {
				panic(10)
			},
			want: errors.New("panic recovered: 10"),
		},
	}

	log.Init(log.DefaultGlg())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if ok := tt.runtimeErr; ok {
					if want, got := tt.want, recover().(error); !errors.Is(got, want) {
						t.Errorf("not equals. want: %v, got: %v", want, got)
					}
				}
			}()

			got := RecoverFunc(tt.fn)()
			if !errors.Is(got, tt.want) {
				t.Errorf("not equals. want: %v, got: %v", tt.want, got)
			}
		})
	}
}
