package info

import (
	"strings"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func TestShowVersionInfo(t *testing.T) {
	type args struct {
		name   string
		extras map[string]string
	}
	type test struct {
		name      string
		logger    log.Logger
		args      args
		checkFunc func() error
	}

	tests := []test{
		func() test {
			got := make([]string, 0, 10)

			logger := &loggerMock{
				InfoFunc: func(vals ...interface{}) {
					for _, v := range vals {
						if str, ok := v.(string); ok {
							sp := strings.Split(str, "\n")
							got = append(got, sp...)
						}
					}
				},
			}

			return test{
				name:   "show version success",
				logger: logger,
				args: args{
					name: "team",
					extras: map[string]string{
						"test": "vald",
					},
				},
				checkFunc: func() error {
					if len(got) != 6 {
						return errors.Errorf("item count is wrong. want: %v, got: %v", 6, len(got))
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Init(tt.logger)
			ShowVersionInfo(tt.args.extras)(tt.args.name)
			if err := tt.checkFunc(); err != nil {
				t.Error(err)
			}
		})
	}
}
