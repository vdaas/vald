//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
)

func TestServerMode_String(t *testing.T) {
	type test struct {
		name string
		m    ServerMode
		want string
	}

	tests := []test{
		{
			name: "returns REST when in REST mode",
			m:    1,
			want: "REST",
		},

		{
			name: "returns gRPC when in gRPC mode",
			m:    2,
			want: "gRPC",
		},

		{
			name: "returns GraphQL when in GraphQL mode",
			m:    3,
			want: "GraphQL",
		},

		{
			name: "returns unknown when in unknown mode",
			m:    5,
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.String()
			if tt.want != got {
				t.Errorf("String is wrong. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestMode(t *testing.T) {
	type test struct {
		name string
		str  string
		want ServerMode
	}

	tests := []test{
		{
			name: "returns REST when in REST mode (rest)",
			str:  "REST",
			want: REST,
		},

		{
			name: "returns HTTP when in REST mode (http)",
			str:  "HTTP",
			want: REST,
		},

		{
			name: "returns GRPC when in gRPC mode",
			str:  "GRPC",
			want: GRPC,
		},

		{
			name: "returns GraphQL when in GraphQL mode (graphql)",
			str:  "GraphQL",
			want: GQL,
		},

		{
			name: "returns GQL when in GraphQL mode (gql)",
			str:  "GQL",
			want: GQL,
		},

		{
			name: "returns 0 when in unknown mode",
			str:  "Vald",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mode(tt.str)
			if tt.want != got {
				t.Errorf("Mode is wrong. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type test struct {
		name      string
		opts      []Option
		checkFunc func(got *server, gotErr, wantErr error) error
		wantErr   error
	}

	tests := []test{
		func() test {
			type handler struct {
				http.Handler
			}

			hdr := new(handler)

			return test{
				name: "returns REST server instance when in REST mode",
				opts: []Option{
					WithHTTPHandler(hdr),
					WithErrorGroup(nil),
					WithReadHeaderTimeout("1s"),
					WithReadTimeout("2s"),
					WithWriteTimeout("3s"),
					WithIdleTimeout("4s"),
					WithTLSConfig(new(tls.Config)),
				},
				checkFunc: func(got *server, gotErr, wantErr error) error {
					if got.http.srv == nil {
						return errors.New("http srv is nil")
					}

					if !errors.Is(gotErr, wantErr) {
						return errors.Errorf("err is not equals. want: %v, got: %v", wantErr, gotErr)
					}

					return nil
				},
				wantErr: nil,
			}
		}(),

		func() test {
			return test{
				name: "returns nil and error when REST server returns invalid api config error",
				opts: []Option{},
				checkFunc: func(got *server, gotErr, wantErr error) error {
					if got != nil {
						return errors.Errorf("New return not nil: %v", got)
					}

					if !errors.Is(gotErr, wantErr) {
						return errors.Errorf("err is not equals. want: %v, got: %v", wantErr, gotErr)
					}

					return nil
				},
				wantErr: errors.ErrInvalidAPIConfig,
			}
		}(),

		func() test {
			fn := func(g *grpc.Server) {}

			return test{
				name: "returns gRPC server instance when in gRPC mode",
				opts: []Option{
					WithServerMode(GRPC),
					WithGRPCRegistFunc(fn),
					WithGRPCKeepaliveTime("1s"),
					WithGRPCOption([]grpc.ServerOption{}...),
					WithTLSConfig(new(tls.Config)),
				},
				checkFunc: func(got *server, gotErr, wantErr error) error {
					if got.grpc.srv == nil {
						return errors.New("grpc srv is nil")
					}

					if !errors.Is(gotErr, wantErr) {
						return errors.Errorf("err is not equals. want: %v, got: %v", wantErr, gotErr)
					}

					return nil
				},
				wantErr: nil,
			}
		}(),

		func() test {
			return test{
				name: "returns nil and error when gRPC server returns invalid api config error",
				opts: []Option{
					WithServerMode(GRPC),
				},
				checkFunc: func(got *server, gotErr, wantErr error) error {
					if got != nil {
						return errors.Errorf("New return not nil: %v", got)
					}

					if !errors.Is(gotErr, wantErr) {
						return errors.Errorf("err is not equals. want: %v, got: %v", wantErr, gotErr)
					}

					return nil
				},
				wantErr: errors.ErrInvalidAPIConfig,
			}
		}(),

		func() test {
			type handler struct {
				http.Handler
			}

			hdr := new(handler)

			return test{
				name: "returns GQL server instance when in GraphQL mode",
				opts: []Option{
					WithServerMode(GQL),
					WithHTTPHandler(hdr),
					WithErrorGroup(nil),
					WithReadHeaderTimeout("1s"),
					WithReadTimeout("2s"),
					WithWriteTimeout("3s"),
					WithIdleTimeout("4s"),
					WithTLSConfig(new(tls.Config)),
				},
				checkFunc: func(got *server, gotErr, wantErr error) error {
					if got.http.srv == nil {
						return errors.New("http srv is nil")
					}

					if !errors.Is(gotErr, wantErr) {
						return errors.Errorf("err is not equals. want: %v, got: %v", wantErr, gotErr)
					}

					return nil
				},
				wantErr: nil,
			}
		}(),

		func() test {
			return test{
				name: "returns nil and error when GQL server returns invalid api config error",
				opts: []Option{
					WithServerMode(GQL),
				},
				checkFunc: func(got *server, gotErr, wantErr error) error {
					if got != nil {
						return errors.Errorf("New return not nil: %v", got)
					}

					if !errors.Is(gotErr, wantErr) {
						return errors.Errorf("err is not equals. want: %v, got: %v", wantErr, gotErr)
					}

					return nil
				},
				wantErr: errors.ErrInvalidAPIConfig,
			}
		}(),
	}

	log.Init(log.WithLoggerType(logger.NOP.String()))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(tt.opts...)

			if tt.wantErr == nil {
				if err := tt.checkFunc(s.(*server), err, tt.wantErr); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func Test_server_IsRunning(t *testing.T) {
	type test struct {
		name string
		s    *server
		want bool
	}

	tests := []test{
		{
			name: "returns true when server is running",
			s: &server{
				running: true,
			},
			want: true,
		},

		{
			name: "returns false when server is not running",
			s: &server{
				running: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.IsRunning()
			if tt.want != got {
				t.Errorf("IsRunning is wrong. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func Test_server_Name(t *testing.T) {
	type test struct {
		name string
		s    *server
		want string
	}

	tests := []test{
		{
			name: "returns name of server instance field",
			s: &server{
				name: "name",
			},
			want: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Name()
			if tt.want != got {
				t.Errorf("Name is wrong. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func Test_server_ListenAndServe(t *testing.T) {
	type args struct {
		ctx   context.Context
		errCh chan error
	}

	type field struct {
		running        bool
		eg             errgroup.Group
		mode           ServerMode
		pwt            time.Duration
		sddur          time.Duration
		httpSrvStarter func(net.Listener) error
		grpcSrv        *grpc.Server
		lc             *net.ListenConfig
		host           string
		port           uint
		preStartFunc   func() error
	}

	type test struct {
		name      string
		args      args
		field     field
		afterFunc func()
		want      error
	}

	tests := []test{
		func() test {
			return test{
				name: "returns nil when server is already running",
				field: field{
					running: true,
				},
				want: nil,
			}
		}(),
		func() test {
			err := errors.New("faild to prestart")

			return test{
				name: "returns error when prestart function return error",
				field: field{
					running: false,
					preStartFunc: func() error {
						return err
					},
				},
				want: err,
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.New(ctx)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})

			srv := &http.Server{
				Handler: handler,
			}

			return test{
				name: "returns nil when serving of REST server is successes",
				field: field{
					mode:           REST,
					eg:             eg,
					httpSrvStarter: srv.Serve,
					host:           "vald",
					port:           8081,
					lc:             &net.ListenConfig{},
					preStartFunc: func() error {
						return nil
					},
					running: false,
				},
				afterFunc: func() {
					srv.Shutdown(ctx)
					cancel()
					eg.Wait()
				},
				want: nil,
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, _ := errgroup.New(ctx)

			srv := new(grpc.Server)

			return test{
				name: "returns nil when serving of gRPC server is successes",
				field: field{
					mode:           GRPC,
					eg:             eg,
					httpSrvStarter: srv.Serve,
					grpcSrv:        srv,
					host:           "vald",
					port:           8082,
					lc:             &net.ListenConfig{},
					preStartFunc: func() error {
						return nil
					},
					running: false,
				},
				afterFunc: func() {
					cancel()
					eg.Wait()
				},
				want: nil,
			}
		}(),
	}

	log.Init(log.WithLoggerType(logger.NOP.String()))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.afterFunc != nil {
					defer tt.afterFunc()
				}
			}()

			s := &server{
				mode: tt.field.mode,
				eg:   tt.field.eg,
				http: struct {
					srv     *http.Server
					h       http.Handler
					starter func(net.Listener) error
				}{
					starter: tt.field.httpSrvStarter,
				},
				grpc: struct {
					srv       *grpc.Server
					keepAlive *grpcKeepalive
					opts      []grpc.ServerOption
					regs      []func(*grpc.Server)
				}{
					srv: tt.field.grpcSrv,
				},
				lc:           tt.field.lc,
				pwt:          tt.field.pwt,
				sddur:        tt.field.sddur,
				running:      tt.field.running,
				preStartFunc: tt.field.preStartFunc,
			}

			got := s.ListenAndServe(tt.args.ctx, tt.args.errCh)
			if !errors.Is(got, tt.want) {
				t.Errorf("ListenAndServe returns error: %v", got)
			}
		})
	}
}

func Test_server_Shutdown(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type field struct {
		running     bool
		eg          errgroup.Group
		mode        ServerMode
		pwt         time.Duration
		sddur       time.Duration
		httpSrv     *http.Server
		grpcSrv     *grpc.Server
		preStopFunc func() error
	}

	type test struct {
		name      string
		args      args
		field     field
		afterFunc func()
		checkFunc func(s *server, got, want error) error
		want      error
	}

	defaultCheckFunc := func(s *server, got, want error) error {
		if want != got {
			return errors.Errorf("Shutdown returns error: %v", got)
		}
		var running bool
		s.mu.RLock()
		running = s.running
		s.mu.RUnlock()

		if running {
			return errors.New("server is running")
		}
		return nil
	}

	tests := []test{
		func() test {
			return test{
				name: "returns nil when server is not running",
				want: nil,
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.New(ctx)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			testSrv := httptest.NewServer(handler)

			return test{
				name: "returns nil when shutdown of REST server is successes",
				args: args{
					ctx: ctx,
				},
				field: field{
					mode:    REST,
					eg:      eg,
					pwt:     10 * time.Millisecond,
					sddur:   1 * time.Second,
					running: true,
					httpSrv: testSrv.Config,
					preStopFunc: func() error {
						return nil
					},
				},
				afterFunc: func() {
					testSrv.Close()
					cancel()
					eg.Wait()
				},
				want: nil,
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(context.Background())
			eg, ctx := errgroup.New(ctx)

			grpcSrv := grpc.NewServer()

			return test{
				name: "returns nil when shutdown of gRPC server is successes",
				args: args{
					ctx: ctx,
				},
				field: field{
					mode:    GRPC,
					eg:      eg,
					pwt:     10 * time.Millisecond,
					sddur:   1 * time.Second,
					running: true,
					grpcSrv: grpcSrv,
					preStopFunc: func() error {
						return nil
					},
				},
				afterFunc: func() {
					cancel()
					eg.Wait()
				},
				want: nil,
			}
		}(),
	}

	log.Init(log.WithLoggerType(logger.NOP.String()))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.afterFunc != nil {
					defer tt.afterFunc()
				}
			}()
			if tt.checkFunc == nil {
				tt.checkFunc = defaultCheckFunc
			}

			s := &server{
				mode: tt.field.mode,
				eg:   tt.field.eg,
				http: struct {
					srv     *http.Server
					h       http.Handler
					starter func(net.Listener) error
				}{
					srv: tt.field.httpSrv,
				},
				grpc: struct {
					srv       *grpc.Server
					keepAlive *grpcKeepalive
					opts      []grpc.ServerOption
					regs      []func(*grpc.Server)
				}{
					srv: tt.field.grpcSrv,
				},
				pwt:         tt.field.pwt,
				sddur:       tt.field.sddur,
				running:     tt.field.running,
				preStopFunc: tt.field.preStopFunc,
			}

			got := s.Shutdown(tt.args.ctx)
			if err := tt.checkFunc(s, got, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}
