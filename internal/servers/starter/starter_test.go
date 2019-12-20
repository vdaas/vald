package starter

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/servers"
	"github.com/vdaas/vald/internal/servers/server"
	"github.com/vdaas/vald/internal/tls"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}

	type test struct {
		name      string
		args      args
		checkFunc func(s Server, err error) error
	}

	tests := []test{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestSetupAPIs(t *testing.T) {
	type handler struct {
		http.Handler
	}

	type args struct {
		cfg *tls.Config
	}

	type field struct {
		cfg  *config.Servers
		rest func(cfg *config.Server) []server.Option
		grpc func(cfg *config.Server) []server.Option
		gql  func(cfg *config.Server) []server.Option
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func([]servers.Option, error) error
	}

	tests := []test{
		func() test {
			fn := func(srv *grpc.Server) {}
			return test{
				name: "setup is success",
				args: args{
					cfg: new(tls.Config),
				},
				field: field{
					cfg: &config.Servers{
						Servers: []*config.Server{
							{
								Mode: server.REST.String(),
								Name: "rest_srv",
								Host: "rest_srv_host",
								Port: 8080,
							},
							{
								Mode: server.GRPC.String(),
								Name: "grpc_srv",
								Host: "grpc_srv_host",
								Port: 8081,
							},
							{
								Mode: server.GQL.String(),
								Name: "gql_srv",
								Host: "gql_srv_host",
								Port: 8082,
							},
						},
					},
					rest: func(cfg *config.Server) []server.Option {
						return []server.Option{
							server.WithHTTPHandler(new(handler)),
						}
					},
					grpc: func(cfg *config.Server) []server.Option {
						return []server.Option{
							server.WithGRPCRegistFunc(fn),
						}
					},
					gql: func(cfg *config.Server) []server.Option {
						return []server.Option{
							server.WithHTTPHandler(new(handler)),
						}
					},
				},
				checkFunc: func(opts []servers.Option, err error) error {
					if err != nil {
						return fmt.Errorf("returns an error: %v", err)
					}

					if len(opts) != 3 {
						return fmt.Errorf("length of options is wrong. want: %v got: %v", 3, len(opts))
					}

					return nil
				},
			}
		}(),
		{
			name: "faild to setup of RESR server",
			args: args{
				cfg: new(tls.Config),
			},
			field: field{
				cfg: &config.Servers{
					Servers: []*config.Server{
						{
							Mode: server.REST.String(),
							Name: "rest_srv",
							Host: "rest_srv_host",
							Port: 8080,
						},
					},
				},
				rest: func(cfg *config.Server) []server.Option { return nil },
			},
			checkFunc: func(opts []servers.Option, err error) error {
				if err == nil {
					return fmt.Errorf("error is nil")
				} else if got, want := err, errors.ErrInvalidAPIConfig; got.Error() != want.Error() {
					return fmt.Errorf("error is not equals. want: %v, got: %v", want, got)
				}
				return nil
			},
		},
		{
			name: "faild to setup of gRPC server",
			args: args{
				cfg: new(tls.Config),
			},
			field: field{
				cfg: &config.Servers{
					Servers: []*config.Server{
						{
							Mode: server.GRPC.String(),
							Name: "grpc_srv",
							Host: "grpc_srv_host",
							Port: 8080,
						},
					},
				},
				grpc: func(cfg *config.Server) []server.Option { return nil },
			},
			checkFunc: func(opts []servers.Option, err error) error {
				if err == nil {
					return fmt.Errorf("error is nil")
				} else if got, want := err, errors.ErrInvalidAPIConfig; got.Error() != want.Error() {
					return fmt.Errorf("error is not equals. want: %v, got: %v", want, got)
				}
				return nil
			},
		},
		{
			name: "faild to setup of GQL server",
			args: args{
				cfg: new(tls.Config),
			},
			field: field{
				cfg: &config.Servers{
					Servers: []*config.Server{
						{
							Mode: server.GQL.String(),
							Name: "gql_srv",
							Host: "gql_srv_host",
							Port: 8080,
						},
					},
				},
				gql: func(cfg *config.Server) []server.Option { return nil },
			},
			checkFunc: func(opts []servers.Option, err error) error {
				if err == nil {
					return fmt.Errorf("error is nil")
				} else if got, want := err, errors.ErrInvalidAPIConfig; got.Error() != want.Error() {
					return fmt.Errorf("error is not equals. want: %v, got: %v", want, got)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srvs := &srvs{
				cfg:  tt.field.cfg,
				rest: tt.field.rest,
				grpc: tt.field.grpc,
				gql:  tt.field.gql,
			}

			opts, err := srvs.setupAPIs(tt.args.cfg)
			if err := tt.checkFunc(opts, err); err != nil {
				t.Error(err)
			}
		})
	}

}

func TestSetupHealthCheck(t *testing.T) {
	type args struct {
		cfg *tls.Config
	}

	type field struct {
		cfg *config.Servers
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func([]servers.Option, error) error
	}

	tests := []test{
		func() test {
			return test{
				name: "setup is success",
				args: args{
					cfg: new(tls.Config),
				},
				field: field{
					cfg: &config.Servers{
						HealthCheckServers: []*config.Server{
							{
								Name: "name",
								Host: "host",
								Port: 8080,
							},
						},
					},
				},
				checkFunc: func(opts []servers.Option, err error) error {
					if err != nil {
						return fmt.Errorf("returns an error: %v", err)
					}

					if len(opts) != 1 {
						return fmt.Errorf("length of options is wrong. want: %v got: %v", 1, len(opts))
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srvs := &srvs{
				cfg: tt.field.cfg,
			}

			opts, err := srvs.setupHealthCheck(tt.args.cfg)
			if err := tt.checkFunc(opts, err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSetupMetrics(t *testing.T) {
	type args struct {
		cfg *tls.Config
	}

	type field struct {
		cfg *config.Servers
	}

	type test struct {
		name      string
		args      args
		field     field
		checkFunc func([]servers.Option, error) error
	}

	tests := []test{
		func() test {
			return test{
				name: "setup is success",
				args: args{
					cfg: new(tls.Config),
				},
				field: field{
					cfg: &config.Servers{
						MetricsServers: []*config.Server{
							{
								Name: "",
							},
							{
								Name: "pprof",
								Host: "host",
								Port: 8080,
							},
						},
					},
				},
				checkFunc: func(opts []servers.Option, err error) error {
					if err != nil {
						return fmt.Errorf("returns an error: %v", err)
					}

					if len(opts) != 1 {
						return fmt.Errorf("length of options is wrong. want: %v got: %v", 1, len(opts))
					}

					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srvs := &srvs{
				cfg: tt.field.cfg,
			}

			opts, err := srvs.setupMetrics(tt.args.cfg)
			if err := tt.checkFunc(opts, err); err != nil {
				t.Error(err)
			}
		})
	}
}
