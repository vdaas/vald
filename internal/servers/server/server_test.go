package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	time "time"

	errgroup "github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/tcp"

	"google.golang.org/grpc"
)

func TestString(t *testing.T) {
	type test struct {
		name string
		m    mode
		want string
	}

	tests := []test{
		{
			name: "REST mode",
			m:    1,
			want: "REST",
		},

		{
			name: "gRPC mode",
			m:    2,
			want: "gRPC",
		},

		{
			name: "GraphQL mode",
			m:    3,
			want: "GraphQL",
		},

		{
			name: "unknown mode",
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
	type args struct {
		m string
	}

	type test struct {
		name string
		args args
		want mode
	}

	tests := []test{
		{
			name: "REST mode (rest)",
			args: args{
				m: "REST",
			},
			want: REST,
		},
		{
			name: "REST mode (http)",
			args: args{
				m: "HTTP",
			},
			want: REST,
		},
		{
			name: "gRPC mode",
			args: args{
				m: "GRPC",
			},
			want: GRPC,
		},
		{
			name: "GraphQL mode (graphql)",
			args: args{
				m: "GraphQL",
			},
			want: GQL,
		},
		{
			name: "GraphQL mode (gql)",
			args: args{
				m: "GQL",
			},
			want: GQL,
		},
		{
			name: "unknown mode",
			args: args{
				m: "Vald",
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mode(tt.args.m)

			if tt.want != got {
				t.Errorf("Mode is wrong. want: %v, got: %v", tt.want, got)
			}
		})
	}
}

func aTestNew(t *testing.T) {
	type args struct {
		opts []Option
	}

	type test struct {
		name      string
		args      args
		checkFunc func(got, want *server) error
		want      *server
	}

	tests := []test{
		func() test {
			type handler struct {
				http.Handler
			}

			hdr := new(handler)

			srv := &http.Server{
				ReadHeaderTimeout: 1 * time.Second,
				ReadTimeout:       2 * time.Second,
				WriteTimeout:      3 * time.Second,
				IdleTimeout:       4 * time.Second,
				Handler:           nil,
			}
			srv.SetKeepAlivesEnabled(true)

			lsCfg := &net.ListenConfig{
				Control: tcp.Control,
			}

			return test{
				name: "initialize REST server",
				args: args{
					opts: []Option{
						WithHTTPHandler(hdr),
						WithHTTPServer(srv),
						WithListenConfig(lsCfg),
						WithReadHeaderTimeout("1s"),
						WithReadTimeout("2s"),
						WithWriteTimeout("3s"),
						WithIdleTimeout("4s"),
					},
				},
				want: &server{
					enableRestart: false,
					mode:          REST,
					eg:            errgroup.Get(),
					http: struct {
						srv     *http.Server
						h       http.Handler
						starter func(net.Listener) error
					}{
						srv:     srv,
						h:       hdr,
						starter: srv.Serve,
					},
					lc:  lsCfg,
					rht: 1 * time.Second,
					rt:  2 * time.Second,
					wt:  3 * time.Second,
					it:  4 * time.Second,
				},
				checkFunc: func(got *server, want *server) error {
					if got, want := got.http.srv, want.http.srv; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("server.http.srv is wrong. want: %v, got: %v", want, got)
					}

					if got, want := got.http.h, want.http.h; !reflect.DeepEqual(got, want) {
						return fmt.Errorf("server.http.h is wrong. want: %v, got: %v", want, got)
					}

					if got, want := reflect.ValueOf(got.http.starter).Pointer(), reflect.ValueOf(want.http.starter).Pointer(); !reflect.DeepEqual(got, want) {
						return fmt.Errorf("server.http.starter is wrong. want: %v, got: %v", want, got)
					}

					if got, want := reflect.ValueOf(got.http.srv).Pointer(), reflect.ValueOf(want.http.srv).Pointer(); !reflect.DeepEqual(got, want) {
						return fmt.Errorf("server.http.srv is wrong. want: %v, got: %v", want, got)
					}

					h := struct {
						srv     *http.Server
						h       http.Handler
						starter func(net.Listener) error
					}{}
					got.http, want.http = h, h

					if !reflect.DeepEqual(got, want) {
						return fmt.Errorf("want: %v, got: %v", want, got)
					}

					return nil
				},
			}
		}(),

		func() test {
			type handler struct {
				http.Handler
			}

			return test{
				name: "initialize gRPC server",
				args: args{
					opts: []Option{},
				},
				want: &server{
					enableRestart: false,
					mode:          REST,
					eg:            errgroup.Get(),
					grpc: struct {
						srv       *grpc.Server
						keepAlive *grpcKeepAlive
						opts      []grpc.ServerOption
						reg       func(*grpc.Server)
					}{},
				},
				checkFunc: func(got *server, want *server) error {
					if !reflect.DeepEqual(got, want) {
						return fmt.Errorf("want: %v, got: %v", want, got)
					}
					return nil
				},
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := New(tt.args.opts...)

			if err := tt.checkFunc(got.(*server), tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestIsRunning(t *testing.T) {
	type test struct {
		name string
		s    *server
		want bool
	}

	tests := []test{
		{
			name: "server is running",
			s: &server{
				running: true,
			},
			want: true,
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

func TestName(t *testing.T) {
	type test struct {
		name string
		s    *server
		want string
	}

	tests := []test{
		{
			name: "server is running",
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

func TestListenAndServe(t *testing.T) {
	type args struct {
		ctx   context.Context
		errCh chan error
	}

	type field struct {
		running        bool
		eg             errgroup.Group
		mode           mode
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
		checkFunc func(s *server, got, want error) error
		want      error
	}

	tests := []test{
		{
			name: "server already running",
			field: field{
				running: true,
			},
			checkFunc: func(s *server, got, want error) error {
				if want != got {
					t.Errorf("ListenAndServe returns error: %v", got)
				}
				return nil
			},
			want: nil,
		},

		func() test {
			err := fmt.Errorf("faild to prestart")

			return test{
				name: "prestart error",
				field: field{
					running: false,
					preStartFunc: func() error {
						return err
					},
				},
				checkFunc: func(s *server, got, want error) error {
					if want != got {
						t.Errorf("ListenAndServe returns error: %v", got)
					}
					return nil
				},
				want: err,
			}
		}(),

		func() test {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})

			srv := &http.Server{
				Handler: handler,
			}

			return test{
				name: "serving of REST server is successful",
				field: field{
					mode:           REST,
					eg:             errgroup.Get(),
					httpSrvStarter: srv.Serve,
					host:           "vald",
					port:           8081,
					lc: &net.ListenConfig{
						Control: tcp.Control,
					},
					preStartFunc: func() error {
						return nil
					},
					running: false,
				},
				checkFunc: func(s *server, got, want error) error {
					if want != got {
						t.Errorf("ListenAndServe returns error: %v", got)
					}
					return nil
				},
				want: nil,
			}
		}(),

		func() test {
			srv := new(grpc.Server)

			return test{
				name: "serving of gRPC server is successful",
				field: field{
					mode:           GRPC,
					eg:             errgroup.Get(),
					httpSrvStarter: srv.Serve,
					grpcSrv:        srv,
					host:           "vald",
					port:           8082,
					lc: &net.ListenConfig{
						Control: tcp.Control,
					},
					preStartFunc: func() error {
						return nil
					},
					running: false,
				},
				checkFunc: func(s *server, got, want error) error {
					if want != got {
						t.Errorf("ListenAndServe returns error: %v", got)
					}
					return nil
				},
				want: nil,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.afterFunc != nil {
					defer tt.afterFunc()
				}
			}()
			log.Init(log.DefaultGlg())

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
					keepAlive *grpcKeepAlive
					opts      []grpc.ServerOption
					reg       func(*grpc.Server)
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
			if err := tt.checkFunc(s, got, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestShutdown(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	type field struct {
		running     bool
		eg          errgroup.Group
		mode        mode
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

	tests := []test{
		{
			name: "server not running",
			checkFunc: func(s *server, got, want error) error {
				if want != got {
					t.Errorf("Shutdown returns error: %v", got)
				}
				return nil
			},
			want: nil,
		},

		func() test {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			testSrv := httptest.NewServer(handler)

			return test{
				name: "shutdown REST server",
				args: args{
					ctx: context.Background(),
				},
				field: field{
					mode:    REST,
					eg:      errgroup.Get(),
					pwt:     10 * time.Millisecond,
					sddur:   1 * time.Second,
					running: true,
					httpSrv: testSrv.Config,
					preStopFunc: func() error {
						return nil
					},
				},
				checkFunc: func(s *server, got, want error) error {
					if want != got {
						t.Errorf("Shutdown returns error: %v", got)
					}
					return nil
				},
				afterFunc: func() {
					testSrv.Close()
				},
				want: nil,
			}
		}(),

		func() test {
			grpcSrv := grpc.NewServer()

			return test{
				name: "shutdown gRPC server",
				args: args{
					ctx: context.Background(),
				},
				field: field{
					mode:    GRPC,
					eg:      errgroup.Get(),
					pwt:     10 * time.Millisecond,
					sddur:   1 * time.Second,
					running: true,
					grpcSrv: grpcSrv,
					preStopFunc: func() error {
						return nil
					},
				},
				checkFunc: func(s *server, got, want error) error {
					if want != got {
						t.Errorf("Shutdown returns error: %v", got)
					}
					return nil
				},
				afterFunc: func() {
				},
				want: nil,
			}
		}(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.afterFunc != nil {
					defer tt.afterFunc()
				}
			}()
			log.Init(log.DefaultGlg())

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
					keepAlive *grpcKeepAlive
					opts      []grpc.ServerOption
					reg       func(*grpc.Server)
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
