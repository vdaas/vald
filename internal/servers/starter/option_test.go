package starter

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/servers/server"
)

func TestWithConfig(t *testing.T) {
	type args struct {
		cfg *config.Servers
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			cfg := new(config.Servers)

			return test{
				name: "set success",
				args: args{
					cfg: cfg,
				},
				checkFunc: func(opt Option) error {
					got := new(srvs)
					opt(got)

					if !reflect.DeepEqual(got.cfg, cfg) {
						return fmt.Errorf("invalid param was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithConfig(tt.args.cfg)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGRPC(t *testing.T) {
	type args struct {
		opts func(cfg *config.Server) []server.Option
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			srvOpts := func(cfg *config.Server) []server.Option {
				return nil
			}

			return test{
				name: "set success",
				args: args{
					opts: srvOpts,
				},
				checkFunc: func(opt Option) error {
					got := new(srvs)
					opt(got)

					if reflect.ValueOf(got.grpc).Pointer() != reflect.ValueOf(srvOpts).Pointer() {
						return fmt.Errorf("invalid param was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGRPC(tt.args.opts)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithREST(t *testing.T) {
	type args struct {
		opts func(cfg *config.Server) []server.Option
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			srvOpts := func(cfg *config.Server) []server.Option {
				return nil
			}

			return test{
				name: "set success",
				args: args{
					opts: srvOpts,
				},
				checkFunc: func(opt Option) error {
					got := new(srvs)
					opt(got)

					if reflect.ValueOf(got.rest).Pointer() != reflect.ValueOf(srvOpts).Pointer() {
						return fmt.Errorf("invalid param was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithREST(tt.args.opts)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithGQL(t *testing.T) {
	type args struct {
		opts func(cfg *config.Server) []server.Option
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			srvOpts := func(cfg *config.Server) []server.Option {
				return nil
			}

			return test{
				name: "set success",
				args: args{
					opts: srvOpts,
				},
				checkFunc: func(opt Option) error {
					got := new(srvs)
					opt(got)

					if reflect.ValueOf(got.gql).Pointer() != reflect.ValueOf(srvOpts).Pointer() {
						return fmt.Errorf("invalid param was set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithGQL(tt.args.opts)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithPreStartFunc(t *testing.T) {
	type args struct {
		name string
		fn   func() error
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			name := "name"
			fn := func() error { return nil }
			return test{
				name: "set success",
				args: args{
					name: name,
					fn:   fn,
				},
				checkFunc: func(opt Option) error {
					got := &srvs{
						pstartf: make(map[string]func() error, 1),
					}
					opt(got)

					if fn, ok := got.pstartf[name]; ok {
						if reflect.ValueOf(fn).Pointer() != reflect.ValueOf(fn).Pointer() {
							return fmt.Errorf("invalid param was set")
						}
					} else {
						return fmt.Errorf("param was not set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPreStartFunc(tt.args.name, tt.args.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWithPreStopFunc(t *testing.T) {
	type args struct {
		name string
		fn   func() error
	}

	type test struct {
		name      string
		args      args
		checkFunc func(Option) error
	}

	tests := []test{
		func() test {
			name := "name"
			fn := func() error { return nil }
			return test{
				name: "set success",
				args: args{
					name: name,
					fn:   fn,
				},
				checkFunc: func(opt Option) error {
					got := &srvs{
						pstopf: make(map[string]func() error, 1),
					}
					opt(got)

					if fn, ok := got.pstopf[name]; ok {
						if reflect.ValueOf(fn).Pointer() != reflect.ValueOf(fn).Pointer() {
							return fmt.Errorf("invalid param was set")
						}
					} else {
						return fmt.Errorf("param was not set")
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithPreStopFunc(tt.args.name, tt.args.fn)
			if err := tt.checkFunc(opt); err != nil {
				t.Error(err)
			}
		})
	}
}
