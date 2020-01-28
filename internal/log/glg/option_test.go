package glg

import (
	"reflect"
	"testing"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/internal/errors"
)

func TestWithGlg(t *testing.T) {
	type test struct {
		name      string
		g         *glg.Glg
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			g := glg.New()

			return test{
				name: "set success when glg is not nil",
				g:    g,
				checkFunc: func(opt Option) error {
					got := new(logger)
					opt(got)

					if !reflect.DeepEqual(got.glg, g) {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),

		func() test {
			g := glg.New()

			return test{
				name: "set success when glg is not nil",
				g:    nil,
				checkFunc: func(opt Option) error {
					got := &logger{
						glg: g,
					}
					opt(got)

					if !reflect.DeepEqual(got.glg, g) {
						return errors.New("invalid params was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGlg(tt.g)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithEnableJSON(t *testing.T) {}
func TestWithFormat(t *testing.T)     {}
func TestWithLevel(t *testing.T)      {}
func TestWithRetry(t *testing.T)      {}
