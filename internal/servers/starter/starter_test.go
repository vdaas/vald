// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package starter

import (
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
	type test struct {
		name      string
		opts      []Option
		checkFunc func(s Server, err error) error
	}

	tests := []test{
		{
			name: "initialize is success",
			opts: []Option{
				WithConfig(&config.Servers{
					TLS: &config.TLS{
						Enabled: true,
						Cert:    "./testdata/dummyServer.crt",
						CA:      "./testdata/dummyCa.pem",
						Key:     "./testdata/dummyServer.key",
					},
				}),
			},
			checkFunc: func(s Server, err error) error {
				if err != nil {
					return errors.Errorf("return an error: %v", err)
				}

				if s == nil {
					return errors.New("server is nil")
				}

				return nil
			},
		},

		{
			name: "initialize is faild when tls.New returns error",
			opts: []Option{
				WithConfig(&config.Servers{
					TLS: &config.TLS{
						Enabled: true,
					},
				}),
			},
			checkFunc: func(s Server, err error) error {
				if err == nil {
					return errors.New("error is nil")
				}

				if s != nil {
					return errors.New("server is not nil")
				}

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, err := New(tt.opts...)
			if err := tt.checkFunc(srv, err); err != nil {
				t.Error(err)
			}
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
				name: "returns options and nil",
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
						return errors.Errorf("returns an error: %v", err)
					}

					if len(opts) != 3 {
						return errors.Errorf("length of options is wrong. want: %v got: %v", 3, len(opts))
					}

					return nil
				},
			}
		}(),

		{
			name: "returns nil options and error when setup of REST server fails",
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
					return errors.New("error is nil")
				} else if got, want := err, errors.ErrInvalidAPIConfig; got.Error() != want.Error() {
					return errors.Errorf("error is not equals. want: %v, got: %v", want, got)
				}
				return nil
			},
		},

		{
			name: "returns nil options and error when setup of gRPC server fails",
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
					return errors.New("error is nil")
				} else if got, want := err, errors.ErrInvalidAPIConfig; got.Error() != want.Error() {
					return errors.Errorf("error is not equals. want: %v, got: %v", want, got)
				}
				return nil
			},
		},

		{
			name: "returns nil options and error when setup of GQL server fails",
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
					return errors.New("error is nil")
				} else if got, want := err, errors.ErrInvalidAPIConfig; got.Error() != want.Error() {
					return errors.Errorf("error is not equals. want: %v, got: %v", want, got)
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
		{
			name: "returns options and nil",
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
					return errors.Errorf("returns an error: %v", err)
				}

				if len(opts) != 1 {
					return errors.Errorf("length of options is wrong. want: %v got: %v", 1, len(opts))
				}

				return nil
			},
		},

		{
			name: "returns nil option and error when server.New returns error",
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
							Mode: server.GRPC.String(),
						},
					},
				},
			},
			checkFunc: func(opts []servers.Option, err error) error {
				if got, want := err, errors.ErrInvalidAPIConfig; !errors.Is(want, got) {
					return errors.Errorf("error is wrong. want: %v, got: %v", want, got)
				}

				if len(opts) != 0 {
					return errors.Errorf("length of options is wrong. want: %v got: %v", 0, len(opts))
				}

				return nil
			},
		},
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
		{
			name: "returns options and nil",
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
					return errors.Errorf("returns an error: %v", err)
				}

				if len(opts) != 1 {
					return errors.Errorf("length of options is wrong. want: %v got: %v", 1, len(opts))
				}
				return nil
			},
		},

		{
			name: "returns nil option and error when server.New returns error",
			args: args{
				cfg: new(tls.Config),
			},
			field: field{
				cfg: &config.Servers{
					MetricsServers: []*config.Server{
						{
							Name: "pprof",
							Host: "host",
							Port: 8080,
							Mode: server.GRPC.String(),
						},
					},
				},
			},
			checkFunc: func(opts []servers.Option, err error) error {
				if got, want := err, errors.ErrInvalidAPIConfig; !errors.Is(want, got) {
					return errors.Errorf("error is wrong. want: %v, got: %v", want, got)
				}

				if len(opts) != 0 {
					return errors.Errorf("length of options is wrong. want: %v got: %v", 0, len(opts))
				}
				return nil
			},
		},
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
