package glg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/format"
)

func TestNew(t *testing.T) {
	type test struct {
		name string
		opts []Option
		want *logger
	}

	tests := []test{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestSetLevelMode(t *testing.T) {}

func TestSetLogFormat(t *testing.T) {
	type args struct {
		fmt format.Format
	}

	type field struct {
		glg *glg.Glg
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func(got *logger) error
	}

	tests := []test{
		func() test {
			return test{
				name: "returns logger object updated the glg object when format is JSON",
				args: args{
					fmt: format.JSON,
				},
				field: field{
					glg: glg.New().SetMode(glg.BOTH),
				},
				checkFunc: func(got *logger) error {
					buf := new(bytes.Buffer)
					got.glg.SetLevelWriter(glg.INFO, buf)
					got.glg.Info("vald")

					var obj map[string]interface{}
					if err := json.NewDecoder(buf).Decode(&obj); err != nil {
						return errors.New("not in JSON output mode")
					}
					return nil
				},
			}
		}(),

		func() test {
			return test{
				name: "returns logger object without updating the glg object when format is invalid",
				args: args{
					fmt: format.Unknown,
				},
				field: field{
					glg: glg.New().SetMode(glg.BOTH),
				},
				checkFunc: func(got *logger) error {
					buf := new(bytes.Buffer)
					got.glg.AddLevelWriter(glg.INFO, buf)
					got.glg.Info("vald")

					var obj map[string]interface{}
					if err := json.NewDecoder(buf).Decode(&obj); err == nil {
						fmt.Println(obj)
						return errors.New("not in RAW output mode")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				tt.field.glg.DisableJSON()
			}()

			l := (&logger{
				glg: tt.field.glg,
			}).setLogFormat(tt.args.fmt)

			if err := tt.checkFunc(l); err != nil {
				t.Error(err)
			}
		})
	}
}
