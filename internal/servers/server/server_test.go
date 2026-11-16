// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/tls"
	"golang.org/x/net/http2"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// loopbackHost is the IPv4 loopback address the test servers bind to.
const loopbackHost = "127.0.0.1"

func TestServerMode_String(t *testing.T) {
	type test struct {
		name string
		want string
		m    ServerMode
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

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(_ *testing.T) {
			got := test.m.String()
			if test.want != got {
				t.Errorf("String is wrong. want: %v, got: %v", test.want, got)
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
		wantErr   error
		checkFunc func(got *server, gotErr, wantErr error) error
		name      string
		opts      []Option
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
			fn := func(_ *grpc.Server) {}

			return test{
				name: "returns gRPC server instance when in gRPC mode",
				opts: []Option{
					WithServerMode(GRPC),
					WithGRPCRegisterar(fn),
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
		s    *server
		name string
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
		eg             errgroup.Group
		httpSrvStarter func(net.Listener) error
		grpcSrv        *grpc.Server
		lc             *net.ListenConfig
		preStartFunc   func() error
		host           string
		pwt            time.Duration
		sddur          time.Duration
		port           uint
		running        bool
		mode           ServerMode
	}

	type test struct {
		args      args
		want      error
		afterFunc func(*testing.T)
		name      string
		field     field
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, ctx := errgroup.New(ctx)

			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			srv := &http.Server{
				Handler:           handler,
				ReadHeaderTimeout: 3 * time.Second,
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
				afterFunc: func(t *testing.T) {
					t.Helper()
					srv.Shutdown(ctx)
					cancel()
					eg.Wait()
				},
				want: nil,
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
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
				afterFunc: func(t *testing.T) {
					t.Helper()
					cancel()
					eg.Wait()
				},
				want: nil,
			}
		}(),
	}

	log.Init(log.WithLoggerType(logger.NOP.String()))
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer func() {
				if test.afterFunc != nil {
					defer test.afterFunc(tt)
				}
			}()

			s := &server{
				mode: test.field.mode,
				eg:   test.field.eg,
				http: struct {
					h        http.Handler
					srv      *http.Server
					h2srv    *http2.Server
					enableH2 bool
					starter  func(net.Listener) error
				}{
					starter: test.field.httpSrvStarter,
				},
				grpc: struct {
					srv        *grpc.Server
					keepAlive  *grpcKeepalive
					maxMsgSize int
					opts       []grpc.ServerOption
					regs       []func(*grpc.Server)
				}{
					srv: test.field.grpcSrv,
				},
				lc:           test.field.lc,
				pwt:          test.field.pwt,
				sddur:        test.field.sddur,
				running:      test.field.running,
				preStartFunc: test.field.preStartFunc,
			}

			got := s.ListenAndServe(test.args.ctx, test.args.errCh)
			if !errors.Is(got, test.want) {
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
		eg          errgroup.Group
		httpSrv     *http.Server
		grpcSrv     *grpc.Server
		preStopFunc func() error
		pwt         time.Duration
		sddur       time.Duration
		running     bool
		mode        ServerMode
	}

	type test struct {
		args      args
		want      error
		afterFunc func(*testing.T)
		checkFunc func(s *server, got, want error) error
		name      string
		field     field
	}

	defaultCheckFunc := func(s *server, got, want error) error {
		if !errors.Is(want, got) {
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
			ctx, cancel := context.WithCancel(t.Context())
			eg, ctx := errgroup.New(ctx)

			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
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
				afterFunc: func(t *testing.T) {
					t.Helper()
					testSrv.Close()
					cancel()
					eg.Wait()
				},
				want: nil,
			}
		}(),
		func() test {
			ctx, cancel := context.WithCancel(t.Context())
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
				afterFunc: func(t *testing.T) {
					t.Helper()
					cancel()
					eg.Wait()
				},
				want: nil,
			}
		}(),
	}

	log.Init(log.WithLoggerType(logger.NOP.String()))
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer func() {
				if test.afterFunc != nil {
					defer test.afterFunc(tt)
				}
			}()
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			s := &server{
				mode: test.field.mode,
				eg:   test.field.eg,
				http: struct {
					h        http.Handler
					srv      *http.Server
					h2srv    *http2.Server
					enableH2 bool
					starter  func(net.Listener) error
				}{
					srv: test.field.httpSrv,
				},
				grpc: struct {
					srv        *grpc.Server
					keepAlive  *grpcKeepalive
					maxMsgSize int
					opts       []grpc.ServerOption
					regs       []func(*grpc.Server)
				}{
					srv: test.field.grpcSrv,
				},
				pwt:         test.field.pwt,
				sddur:       test.field.sddur,
				running:     test.field.running,
				preStopFunc: test.field.preStopFunc,
			}

			got := s.Shutdown(test.args.ctx)
			if err := test.checkFunc(s, got, test.want); err != nil {
				t.Error(err)
			}
		})
	}
}

// closeErrListener is a net.Listener whose Close always fails with err.
type closeErrListener struct {
	net.Listener
	err error
}

func (l closeErrListener) Close() error {
	return l.err
}

func Test_server_Shutdown_ReturnsErrors(t *testing.T) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	preStopErr := errors.New("failed to pre stop")
	closeErr := errors.New("failed to close listener")

	type test struct {
		newServer func(t *testing.T, eg errgroup.Group) *server
		want      error
		name      string
	}

	tests := []test{
		{
			name: "returns the error of preStopFunc",
			want: preStopErr,
			newServer: func(t *testing.T, eg errgroup.Group) *server {
				t.Helper()
				testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				t.Cleanup(testSrv.Close)
				s := &server{
					mode: REST,
					eg:   eg,
					// long enough that the result of preStopFunc is always received first.
					pwt:         10 * time.Second,
					sddur:       time.Second,
					running:     true,
					preStopFunc: func() error { return preStopErr },
				}
				s.http.srv = testSrv.Config
				return s
			},
		},
		{
			name: "returns the error of closing the listener",
			want: closeErr,
			newServer: func(t *testing.T, eg errgroup.Group) *server {
				t.Helper()
				s := &server{
					mode:     GRPC,
					eg:       eg,
					sddur:    time.Second,
					running:  true,
					listener: closeErrListener{err: closeErr},
				}
				s.grpc.srv = grpc.NewServer()
				return s
			},
		},
		{
			name: "ignores the error of closing the listener which is already closed",
			want: nil,
			newServer: func(t *testing.T, eg errgroup.Group) *server {
				t.Helper()
				s := &server{
					mode:     GRPC,
					eg:       eg,
					sddur:    time.Second,
					running:  true,
					listener: closeErrListener{err: errors.Wrap(net.ErrClosed, "close tcp")},
				}
				s.grpc.srv = grpc.NewServer()
				return s
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			ctx, cancel := context.WithTimeout(tt.Context(), 10*time.Second)
			defer cancel()
			eg, ctx := errgroup.New(ctx)
			defer eg.Wait()
			defer cancel()

			s := test.newServer(tt, eg)
			got := s.Shutdown(ctx)
			if !errors.Is(got, test.want) {
				tt.Errorf("Shutdown must return the error. want: %v, got: %v", test.want, got)
			}
			if s.IsRunning() {
				tt.Error("server is running")
			}
		})
	}
}

func Test_server_ListenAndServe_ResetsRunningOnFailure(t *testing.T) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	preStartErr := errors.New("failed to pre start")
	s := &server{
		mode:         GRPC,
		preStartFunc: func() error { return preStartErr },
	}

	err := s.ListenAndServe(t.Context(), make(chan error, 1))
	if !errors.Is(err, preStartErr) {
		t.Fatalf("unexpected error. want: %v, got: %v", preStartErr, err)
	}
	if s.IsRunning() {
		t.Error("a server that failed to start must not be marked as running")
	}
}

// socketInodeCounts returns how many open fds refer to each socket inode of this process.
func socketInodeCounts(t *testing.T) map[string]int {
	t.Helper()
	const dir = "/proc/self/fd"
	ents, err := os.ReadDir(dir)
	if err != nil {
		t.Skipf("cannot read %s: %v", dir, err)
	}
	m := make(map[string]int)
	for _, e := range ents {
		link, err := os.Readlink(filepath.Join(dir, e.Name()))
		if err != nil {
			// the fd used by ReadDir itself is already closed.
			continue
		}
		if strings.HasPrefix(link, "socket:[") {
			m[link]++
		}
	}
	return m
}

func Test_server_ListenAndServe_DoesNotLeakDuplicatedFd(t *testing.T) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	eg, ctx := errgroup.New(ctx)
	defer eg.Wait()
	defer cancel()

	s := &server{
		mode:  GRPC,
		eg:    eg,
		host:  loopbackHost,
		port:  0,
		sddur: time.Second,
		lc:    &net.ListenConfig{},
	}
	s.grpc.srv = grpc.NewServer()

	before := socketInodeCounts(t)
	if err := s.ListenAndServe(ctx, make(chan error, 4)); err != nil {
		t.Fatal(err)
	}
	after := socketInodeCounts(t)

	sctx, scancel := context.WithTimeout(ctx, 10*time.Second)
	defer scancel()
	if err := s.Shutdown(sctx); err != nil {
		t.Error(err)
	}

	for inode, c := range after {
		if c > 1 && before[inode] < 2 {
			t.Errorf("socket %s is referenced by %d fds: the duplicated listener fd is leaked", inode, c)
		}
	}
}

// Test_server_Shutdown_GRPCBoundedByContext guards against gRPC Shutdown
// ignoring ctx: GracefulStop waits for every pending RPC, so an open stream
// that never ends used to block Shutdown forever. Shutdown must force Stop
// once ctx expires.
func Test_server_Shutdown_GRPCBoundedByContext(t *testing.T) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	eg, ctx := errgroup.New(ctx)
	defer eg.Wait()
	defer cancel()

	entered := make(chan struct{})
	var enterOnce sync.Once
	s := &server{
		mode:  GRPC,
		eg:    eg,
		host:  loopbackHost,
		port:  0,
		sddur: time.Second,
		lc:    &net.ListenConfig{},
	}
	s.grpc.srv = grpc.NewServer()
	s.grpc.srv.RegisterService(&ggrpc.ServiceDesc{
		ServiceName: "vald.test.Hang",
		HandlerType: (*any)(nil),
		Streams: []ggrpc.StreamDesc{{
			StreamName:    "Hang",
			ServerStreams: true,
			ClientStreams: true,
			Handler: func(_ any, stream ggrpc.ServerStream) error {
				enterOnce.Do(func() { close(entered) })
				// a long-lived stream that only ends when its connection does.
				<-stream.Context().Done()
				return stream.Context().Err()
			},
		}},
	}, struct{}{})

	if err := s.ListenAndServe(ctx, make(chan error, 4)); err != nil {
		t.Fatal(err)
	}
	s.mu.RLock()
	addr := s.listener.Addr().String()
	s.mu.RUnlock()

	conn, err := ggrpc.NewClient(addr, ggrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	cctx, ccancel := context.WithCancel(ctx)
	defer ccancel()
	if _, err := conn.NewStream(cctx, &ggrpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}, "/vald.test.Hang/Hang"); err != nil {
		t.Fatal(err)
	}
	select {
	case <-entered:
	case <-time.After(10 * time.Second):
		t.Fatal("stream handler was not invoked")
	}

	sctx, scancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer scancel()
	done := make(chan error, 1)
	go func() {
		done <- s.Shutdown(sctx)
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Shutdown() error = %v, want nil", err)
		}
	case <-time.After(10 * time.Second):
		t.Error("Shutdown did not return after its context expired: GracefulStop is not bounded by ctx")
		// release the stream so the stuck GracefulStop and Shutdown can finish.
		ccancel()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatal("Shutdown did not return even after the client stream was cancelled")
		}
	}
	if s.IsRunning() {
		t.Error("server is running")
	}
}

// Test_server_Shutdown_PreStopPanicDoesNotHang guards against preStopFunc
// panicking with a non-runtime error: RecoverFunc recovers it, so the wait
// group must still be released or Shutdown blocks forever on s.wg.Wait.
func Test_server_Shutdown_PreStopPanicDoesNotHang(t *testing.T) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	eg, ctx := errgroup.New(ctx)
	defer eg.Wait()
	defer cancel()

	s := &server{
		mode:    GRPC,
		eg:      eg,
		pwt:     10 * time.Second,
		sddur:   time.Second,
		running: true,
		preStopFunc: func() error {
			panic("pre stop panic")
		},
	}
	s.grpc.srv = grpc.NewServer()

	done := make(chan error, 1)
	go func() {
		done <- s.Shutdown(ctx)
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		// free the leader blocked on ctx so the goroutine can exit.
		cancel()
		t.Fatal("Shutdown hung after preStopFunc panicked: wait group was never released")
	}
}

// Test_server_Shutdown_ZeroShutdownDurationDrains guards against a zero
// shutdown duration turning into an already expired context, which made
// http.Server.Shutdown return immediately without draining in-flight requests.
func Test_server_Shutdown_ZeroShutdownDurationDrains(t *testing.T) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	eg, ctx := errgroup.New(ctx)
	defer eg.Wait()
	defer cancel()

	entered := make(chan struct{})
	release := make(chan struct{})
	finished := make(chan struct{})
	testSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(entered)
		select {
		case <-release:
		case <-ctx.Done():
		}
		w.WriteHeader(http.StatusOK)
		close(finished)
	}))
	defer testSrv.Close()

	reqDone := make(chan struct{})
	go func() {
		defer close(reqDone)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, testSrv.URL, nil)
		if err != nil {
			return
		}
		resp, err := testSrv.Client().Do(req)
		if err != nil {
			return
		}
		if err := resp.Body.Close(); err != nil {
			t.Logf("failed to close response body: %v", err)
		}
	}()
	select {
	case <-entered:
	case <-time.After(10 * time.Second):
		cancel()
		t.Fatal("request handler was not invoked")
	}

	s := &server{
		mode:    REST,
		eg:      eg,
		sddur:   0,
		running: true,
	}
	s.http.srv = testSrv.Config

	done := make(chan error, 1)
	go func() {
		done <- s.Shutdown(ctx)
	}()
	// Shutdown must still be draining the in-flight request.
	select {
	case <-done:
		t.Error("Shutdown returned before the in-flight request finished: zero shutdown duration skipped draining")
	case <-time.After(200 * time.Millisecond):
	}
	close(release)
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		cancel()
		t.Fatal("Shutdown did not return after the in-flight request finished")
	}
	select {
	case <-finished:
	case <-time.After(10 * time.Second):
		t.Error("in-flight request did not finish")
	}
	select {
	case <-reqDone:
	case <-time.After(10 * time.Second):
		cancel()
		t.Error("client request did not return")
	}
}

// newRestartTestServer returns a running REST server with restart enabled that listens on a random local port.
// When served is non-nil it is closed once the first Serve call returns.
func newRestartTestServer(
	ctx context.Context, t *testing.T, eg errgroup.Group, ech chan<- error, served chan struct{},
) *server {
	t.Helper()
	h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	s := &server{
		mode:          REST,
		eg:            eg,
		host:          loopbackHost,
		port:          0,
		sddur:         time.Second,
		lc:            &net.ListenConfig{},
		enableRestart: true,
	}
	s.http.srv = &http.Server{Handler: h, ReadHeaderTimeout: time.Second}
	s.http.starter = s.http.srv.Serve
	if served != nil {
		var once sync.Once
		s.http.starter = func(l net.Listener) error {
			err := s.http.srv.Serve(l)
			once.Do(func() { close(served) })
			return err
		}
	}
	if err := s.ListenAndServe(ctx, ech); err != nil {
		t.Fatal(err)
	}
	return s
}

// currentListener returns the listener the server is serving on.
func (s *server) currentListener() net.Listener {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.listener
}

// shutdownWithin runs Shutdown and fails the test when it does not return within d.
func shutdownWithin(ctx context.Context, t *testing.T, s *server, d time.Duration) {
	t.Helper()
	done := make(chan error, 1)
	go func() {
		done <- s.Shutdown(ctx)
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Shutdown() error = %v, want nil", err)
		}
	case <-time.After(d):
		t.Fatal("Shutdown did not return")
	}
}

// waitGroupWithin waits for eg and reports an error instead of hanging when it does not finish within d.
func waitGroupWithin(t *testing.T, eg errgroup.Group, d time.Duration) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			t.Logf("errgroup returned: %v", err)
		}
	}()
	select {
	case <-done:
	case <-time.After(d):
		t.Error("server goroutine did not finish")
	}
}

// Test_server_ListenAndServe_RestartsWithNewListener guards against the
// restart loop serving the listener that Serve has already closed, which made
// an enabled restart spin on "use of closed network connection" forever.
func Test_server_ListenAndServe_RestartsWithNewListener(t *testing.T) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	eg, ctx := errgroup.New(ctx)
	defer waitGroupWithin(t, eg, 10*time.Second)
	defer cancel()

	ech := make(chan error, 16)
	s := newRestartTestServer(ctx, t, eg, ech, nil)
	l1 := s.currentListener()
	if l1 == nil {
		t.Fatal("listener is not set")
	}
	// break the listener so that Serve returns with an error.
	if err := l1.Close(); err != nil {
		t.Fatal(err)
	}

	var l2 net.Listener
	deadline := time.After(10 * restartInterval)
	for l2 == nil {
		select {
		case <-deadline:
			cancel()
			t.Fatal("server was not restarted with a new listener")
		case <-time.After(10 * time.Millisecond):
		}
		if l := s.currentListener(); l != nil && l != l1 {
			l2 = l
		}
	}
	if !s.IsRunning() {
		t.Error("restarted server is not running")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+l2.Addr().String()+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("restarted server does not serve: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Logf("failed to close response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	shutdownWithin(ctx, t, s, 10*time.Second)
	if s.IsRunning() {
		t.Error("server is running after Shutdown")
	}
	if l := s.currentListener(); l != nil {
		t.Error("listener is still set after Shutdown")
	}
	// with restart enabled the Serve error is logged, not sent: repeated sends
	// would fill ech and could race its close.
	select {
	case err := <-ech:
		t.Errorf("unexpected error on ech with restart enabled: %v", err)
	default:
	}
}

// newNoRestartTestServer returns a running REST server with restart disabled
// that listens on a random local port and whose Serve fails with serveErr.
func newNoRestartTestServer(
	ctx context.Context, t *testing.T, eg errgroup.Group, ech chan<- error, serveErr error,
) *server {
	t.Helper()
	s := &server{
		mode:  REST,
		eg:    eg,
		host:  loopbackHost,
		port:  0,
		sddur: time.Second,
		lc:    &net.ListenConfig{},
	}
	s.http.srv = &http.Server{Handler: http.NotFoundHandler(), ReadHeaderTimeout: time.Second}
	s.http.starter = func(l net.Listener) error {
		if err := l.Close(); err != nil {
			t.Logf("failed to close listener: %v", err)
		}
		return serveErr
	}
	if err := s.ListenAndServe(ctx, ech); err != nil {
		t.Fatal(err)
	}
	return s
}

// Test_server_ListenAndServe_ErrorSendBoundedByContext guards against the
// serving goroutine blocking forever on sending its Serve error when nobody
// reads ech any more, and verifies the error is delivered while it is read.
func Test_server_ListenAndServe_ErrorSendBoundedByContext(t *testing.T) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	serveErr := errors.New("serve failed")
	tests := []struct {
		name    string
		receive bool
	}{
		{name: "error is delivered while ech is read", receive: true},
		{name: "goroutine exits on ctx done while ech is not read", receive: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			eg, ctx := errgroup.New(ctx)
			defer waitGroupWithin(t, eg, 10*time.Second)
			defer cancel()

			// unbuffered and, without receive, never read: the send can only finish via ctx.
			ech := make(chan error)
			newNoRestartTestServer(ctx, t, eg, ech, serveErr)
			if tc.receive {
				select {
				case err := <-ech:
					if !errors.Is(err, serveErr) {
						t.Errorf("unexpected error. want: %v, got: %v", serveErr, err)
					}
				case <-time.After(10 * time.Second):
					t.Error("Serve error was not delivered")
				}
			}
		})
	}
}

// Test_server_Shutdown_StopsPendingRestart guards against a Shutdown issued
// after Serve has returned but before the server is listening again (usually
// during the restart wait): it was skipped because running had already been
// cleared, and the server was restarted afterwards.
func Test_server_Shutdown_StopsPendingRestart(t *testing.T) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	eg, ctx := errgroup.New(ctx)
	defer waitGroupWithin(t, eg, 10*time.Second)
	defer cancel()

	ech := make(chan error, 16)
	served := make(chan struct{})
	s := newRestartTestServer(ctx, t, eg, ech, served)
	l1 := s.currentListener()
	if l1 == nil {
		t.Fatal("listener is not set")
	}
	if err := l1.Close(); err != nil {
		t.Fatal(err)
	}
	// wait until Serve has returned so that Shutdown hits the pending restart.
	select {
	case <-served:
	case <-time.After(10 * time.Second):
		t.Fatal("Serve did not return after its listener was closed")
	}
	if !s.IsRunning() {
		t.Error("server waiting for a restart must be reported as running so that Shutdown is not skipped")
	}

	shutdownWithin(ctx, t, s, 10*time.Second)
	// Shutdown waits for the serving goroutine, so no restart may happen after it returns.
	if s.IsRunning() {
		t.Error("server is running after Shutdown")
	}
	if l := s.currentListener(); l != nil {
		t.Errorf("server was restarted after Shutdown on %s", l.Addr().String())
	}
	// Shutdown waited for the serving goroutine, so any Serve error it sent is visible here.
	select {
	case err := <-ech:
		t.Errorf("unexpected error on ech with restart enabled: %v", err)
	default:
	}
}

// Test_server_ListenAndServe_TLSAppliedOnce guards against the listener being
// wrapped with TLS on top of ServeTLS / grpc.Creds: a client doing a single,
// ordinary TLS handshake must be able to talk to the server in every mode.
func Test_server_ListenAndServe_TLSAppliedOnce(t *testing.T) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	// borrow httptest's self-signed certificate (valid for 127.0.0.1) and its trust pool.
	ts := httptest.NewTLSServer(http.NotFoundHandler())
	cert := ts.TLS.Certificates[0]
	tr, ok := ts.Client().Transport.(*http.Transport)
	if !ok || tr.TLSClientConfig == nil || tr.TLSClientConfig.RootCAs == nil {
		ts.Close()
		t.Fatal("httptest client has no TLS root CAs")
	}
	roots := tr.TLSClientConfig.RootCAs
	ts.Close()

	newServer := func(t *testing.T, ctx context.Context, eg errgroup.Group, opts ...Option) *server {
		t.Helper()
		srv, err := New(append([]Option{
			WithErrorGroup(eg),
			WithHost(loopbackHost),
			WithShutdownDuration("1s"),
			WithTLSConfig(&tls.Config{Certificates: []tls.Certificate{cert}}),
		}, opts...)...)
		if err != nil {
			t.Fatal(err)
		}
		s, ok := srv.(*server)
		if !ok {
			t.Fatalf("New returned %T, want *server", srv)
		}
		if err := s.ListenAndServe(ctx, make(chan error, 4)); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			sctx, scancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
			defer scancel()
			if err := s.Shutdown(sctx); err != nil {
				t.Error(err)
			}
		})
		return s
	}
	addrOf := func(s *server) string {
		s.mu.RLock()
		defer s.mu.RUnlock()
		return s.listener.Addr().String()
	}

	httpCheck := func(t *testing.T, ctx context.Context, addr string) {
		t.Helper()
		client := &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{RootCAs: roots},
			},
		}
		defer client.CloseIdleConnections()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+addr+"/", http.NoBody)
		if err != nil {
			t.Fatal(err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("TLS request failed (listener double-wrapped with TLS?): %v", err)
		}
		if err := res.Body.Close(); err != nil {
			t.Error(err)
		}
		if res.StatusCode != http.StatusNoContent {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNoContent)
		}
		if res.TLS == nil {
			t.Error("response was not served over TLS")
		}
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	tests := []struct {
		check func(t *testing.T, ctx context.Context, addr string)
		name  string
		opts  []Option
	}{
		{
			name:  "REST",
			opts:  []Option{WithServerMode(REST), WithHTTPHandler(handler)},
			check: httpCheck,
		},
		{
			name:  "GQL",
			opts:  []Option{WithServerMode(GQL), WithHTTPHandler(handler)},
			check: httpCheck,
		},
		{
			name: "GRPC",
			opts: []Option{WithServerMode(GRPC), WithGRPCRegisterar(func(*grpc.Server) {})},
			check: func(t *testing.T, ctx context.Context, addr string) {
				t.Helper()
				conn, err := ggrpc.NewClient(addr, ggrpc.WithTransportCredentials(
					credentials.NewTLS(&tls.Config{RootCAs: roots}),
				))
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if cerr := conn.Close(); cerr != nil {
						t.Error(cerr)
					}
				}()
				res, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{})
				if err != nil {
					t.Fatalf("TLS health check failed (listener double-wrapped with TLS?): %v", err)
				}
				if res.GetStatus() != healthpb.HealthCheckResponse_SERVING {
					t.Errorf("health status = %v, want SERVING", res.GetStatus())
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			eg, ectx := errgroup.New(ctx)
			t.Cleanup(func() {
				cancel()
				done := make(chan error, 1)
				go func() {
					done <- eg.Wait()
				}()
				select {
				case err := <-done:
					if err != nil {
						t.Error(err)
					}
				case <-time.After(10 * time.Second):
					t.Error("server goroutine did not exit")
				}
			})

			s := newServer(t, ectx, eg, tc.opts...)
			cctx, ccancel := context.WithTimeout(ectx, 10*time.Second)
			defer ccancel()
			tc.check(t, cctx, addrOf(s))
		})
	}
}

// NOT IMPLEMENTED BELOW
