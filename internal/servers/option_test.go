package servers

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/servers/server"
)

func TestWithServer(t *testing.T) {
	type args struct {
		srv server.Server
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {

			srv := NewMockServer()
			srv.NameFunc = func() string {
				return "srv"
			}

			return test{
				name: "set success",
				args: args{
					srv: srv,
				},
				checkFunc: func(opt Option) error {
					got := new(listener)
					opt(got)

					if len(got.servers) != 1 {
						return fmt.Errorf("servers count is wrong. want: %v, got: %v", 1, len(got.servers))
					}

					if gsrv, ok := got.servers["srv"]; !ok {
						return fmt.Errorf("servers['srv'] is nothing")
					} else {
						if !reflect.DeepEqual(gsrv, srv) {
							return fmt.Errorf("servers['srv'] is not equals. want: %v, got: %b", srv, gsrv)
						}
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithServer(tt.args.srv)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithErrorGroup(t *testing.T) {
	type args struct {
		eg errgroup.Group
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			eg := errgroup.Get()

			return test{
				name: "set success",
				args: args{
					eg: eg,
				},
				checkFunc: func(opt Option) error {
					got := new(listener)
					opt(got)

					if !reflect.DeepEqual(got.eg, eg) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithErrorGroup(tt.args.eg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithShutdownDuration(t *testing.T) {
	type args struct {
		dur string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		{
			name: "set success",
			args: args{
				dur: "10s",
			},
			checkFunc: func(opt Option) error {
				got := new(listener)
				opt(got)

				if !reflect.DeepEqual(got.sddur, 10*time.Second) {
					return fmt.Errorf("invalid param was set")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithShutdownDuration(tt.args.dur)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithStartUpStrategy(t *testing.T) {
	type args struct {
		strg []string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			strg := []string{
				"strg_1",
				"strg_2",
			}

			return test{
				name: "set success",
				args: args{
					strg: strg,
				},
				checkFunc: func(opt Option) error {
					got := new(listener)
					opt(got)

					if !reflect.DeepEqual(got.sus, strg) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithStartUpStrategy(tt.args.strg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithShutdownStrategy(t *testing.T) {
	type args struct {
		strg []string
	}

	type test struct {
		name      string
		args      args
		checkFunc func(opt Option) error
	}

	tests := []test{
		func() test {
			strg := []string{
				"strg_1",
				"strg_2",
			}

			return test{
				name: "set success",
				args: args{
					strg: strg,
				},
				checkFunc: func(opt Option) error {
					got := new(listener)
					opt(got)

					if !reflect.DeepEqual(got.sds, strg) {
						return fmt.Errorf("invalid param was set")
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithShutdownStrategy(tt.args.strg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
