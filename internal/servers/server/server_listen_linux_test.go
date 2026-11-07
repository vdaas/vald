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

//go:build linux

package server

import (
	"context"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Test_server_ListenAndServe_ListensWithoutSpareFd lowers RLIMIT_NOFILE to the lowest free fd right after
// the listening socket is created (inside the listen control hook), so any extra fd that listen allocates
// (e.g. a dup via lt.File(), which used to fail here with EMFILE) makes the start fail. The server must
// start anyway and serve on a non-blocking listener.
// This test must not call t.Parallel(): the rlimit is process-wide.
func Test_server_ListenAndServe_ListensWithoutSpareFd(t *testing.T) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	// initialize the netpoller before the rlimit is lowered: its first use allocates fds.
	warm, err := net.Listen(net.TCP.String(), "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	if err := warm.Close(); err != nil {
		t.Fatal(err)
	}

	type test struct {
		// bind returns the host/port/socketPath of the server and the address to re-listen on.
		bind    func(t *testing.T) (s *server, relisten func() (net.Listener, error))
		name    string
		network net.NetworkType
	}
	tests := []test{
		{
			name:    "tcp listener needs no spare fd",
			network: net.TCP,
			bind: func(t *testing.T) (*server, func() (net.Listener, error)) {
				t.Helper()
				// reserve a free port, then release it for the server under test.
				tmp, err := net.Listen(net.TCP.String(), "127.0.0.1:0")
				if err != nil {
					t.Fatal(err)
				}
				addr := tmp.Addr().String()
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					t.Fatal(err)
				}
				if err := tmp.Close(); err != nil {
					t.Fatal(err)
				}
				return &server{host: host, port: port}, func() (net.Listener, error) {
					return net.Listen(net.TCP.String(), addr)
				}
			},
		},
		{
			name:    "unix listener needs no spare fd",
			network: net.UNIX,
			bind: func(t *testing.T) (*server, func() (net.Listener, error)) {
				t.Helper()
				// keep the path short: sun_path is limited to 108 bytes and t.TempDir embeds the test name.
				dir, err := os.MkdirTemp("", "vs")
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.RemoveAll(dir); err != nil {
						t.Error(err)
					}
				})
				path := filepath.Join(dir, "s.sock")
				return &server{socketPath: path}, func() (net.Listener, error) {
					return net.Listen(net.UNIX.String(), path)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			ctx, cancel := context.WithTimeout(tt.Context(), 30*time.Second)
			defer cancel()
			eg, ctx := errgroup.New(ctx)
			defer eg.Wait()
			defer cancel()

			s, relisten := test.bind(tt)
			s.mode = GRPC
			s.network = test.network
			s.eg = eg
			s.sddur = time.Second
			s.grpc.srv = grpc.NewServer()

			var orig syscall.Rlimit
			if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &orig); err != nil {
				tt.Fatal(err)
			}
			var (
				limited bool
				ctrlErr error
			)
			s.lc = &net.ListenConfig{
				Control: func(_, _ string, c syscall.RawConn) error {
					if limited {
						return nil
					}
					limited = true
					var derr error
					if err := c.Control(func(fd uintptr) {
						// the lowest free fd is what dup(2) will try to allocate in lt.File().
						var nfd int
						nfd, derr = syscall.Dup(int(fd))
						if derr != nil {
							return
						}
						if derr = syscall.Close(nfd); derr != nil {
							return
						}
						if nfd < 0 {
							derr = syscall.EBADF
							return
						}
						derr = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{Cur: uint64(nfd), Max: orig.Max})
					}); err != nil {
						ctrlErr = err
					} else {
						ctrlErr = derr
					}
					return nil
				},
			}

			err := s.ListenAndServe(ctx, make(chan error, 4))
			// restore the limit before any assertion allocates fds.
			if rerr := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &orig); rerr != nil {
				tt.Fatalf("failed to restore RLIMIT_NOFILE: %v", rerr)
			}
			if ctrlErr != nil {
				tt.Fatalf("failed to lower RLIMIT_NOFILE: %v", ctrlErr)
			}
			if err != nil {
				tt.Fatalf("ListenAndServe needed a spare fd beyond the listening socket: %v", err)
			}
			defer func() {
				sctx, scancel := context.WithTimeout(ctx, 10*time.Second)
				defer scancel()
				if serr := s.Shutdown(sctx); serr != nil {
					tt.Error(serr)
				}
			}()

			s.mu.RLock()
			l := s.listener
			s.mu.RUnlock()
			sc, ok := l.(syscall.Conn)
			if !ok {
				tt.Fatalf("listener %T does not expose its raw fd", l)
			}
			rc, err := sc.SyscallConn()
			if err != nil {
				tt.Fatal(err)
			}
			var (
				flags uintptr
				errno syscall.Errno
			)
			if err := rc.Control(func(fd uintptr) {
				flags, _, errno = syscall.Syscall(syscall.SYS_FCNTL, fd, syscall.F_GETFL, 0)
			}); err != nil {
				tt.Fatal(err)
			}
			if errno != 0 {
				tt.Fatalf("fcntl(F_GETFL) failed: %v", errno)
			}
			if flags&syscall.O_NONBLOCK == 0 {
				tt.Error("listener is in blocking mode, want O_NONBLOCK")
			}
			// the address is held by the running server: re-listening must fail.
			if rl, err := relisten(); err == nil {
				tt.Error("re-listen succeeded while the server holds the address")
				if err := rl.Close(); err != nil {
					tt.Error(err)
				}
			}
		})
	}
}
