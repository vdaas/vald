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
package service

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	grpcmock "github.com/vdaas/vald/internal/test/mock/grpc"
)

// One unreachable target must not stop Connect from connecting the rest.
func Test_mirr_Connect_continuesPastFailure(t *testing.T) {
	t.Parallel()
	const (
		badAddr  = "192.168.10.1:8081"
		goodAddr = "192.168.10.2:8081"
	)
	failErr := errors.New("dial failure")
	var (
		mu        sync.Mutex
		attempted []string
	)
	m := &mirr{
		gateway: &GatewayMock{
			GRPCClientFunc: func() grpc.Client {
				return &grpcmock.GRPCClientMock{
					IsConnectedFunc: func(context.Context, string) bool { return false },
					ConnectFunc: func(_ context.Context, addr string, _ ...grpc.DialOption) (pool.Conn, error) {
						mu.Lock()
						attempted = append(attempted, addr)
						mu.Unlock()
						if addr == badAddr {
							return nil, failErr
						}
						return nil, nil
					},
					SetDisableResolveDNSAddrFunc: func(string, bool) {},
				}
			},
		},
	}

	err := m.Connect(t.Context(),
		&payload.Mirror_Target{Host: "192.168.10.1", Port: 8081},
		&payload.Mirror_Target{Host: "192.168.10.2", Port: 8081},
	)
	if !errors.Is(err, failErr) {
		t.Errorf("error = %v, want %v", err, failErr)
	}
	mu.Lock()
	got := append([]string(nil), attempted...)
	mu.Unlock()
	if len(got) != 2 || got[0] != badAddr || got[1] != goodAddr {
		t.Errorf("attempted = %v, want [%s %s]", got, badAddr, goodAddr)
	}
	if _, ok := m.addrs.Load(goodAddr); !ok {
		t.Errorf("%s not stored after successful connect", goodAddr)
	}
	if _, ok := m.addrs.Load(badAddr); ok {
		t.Errorf("%s stored after failed connect", badAddr)
	}
}

// Start must not run the register broadcast while no peer is connected
// (it can only fail with ErrTargetNotFound), and must run it once one is.
func Test_mirr_Start_registerOnlyWithPeers(t *testing.T) {
	t.Parallel()
	const peerAddr = "192.168.20.1:8081"
	tests := []struct {
		name          string
		peerConnected bool
	}{
		{name: "no connected peer skips register", peerConnected: false},
		{name: "connected peer runs register", peerConnected: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithTimeout(tt.Context(), 10*time.Second)
			defer cancel()

			var (
				isConnectedCalls atomic.Int64
				doMultiCalls     atomic.Int64
				once             sync.Once
			)
			// Signals that the loop has ticked enough (or reached register).
			reached := make(chan struct{})
			signal := func() { once.Do(func() { close(reached) }) }

			eg, egctx := errgroup.New(ctx)
			m := &mirr{
				eg:          eg,
				registerDur: time.Millisecond,
				gateway: &GatewayMock{
					GRPCClientFunc: func() grpc.Client {
						return &grpcmock.GRPCClientMock{
							IsConnectedFunc: func(context.Context, string) bool {
								// MirrorTargets and connectedOtherMirrorAddrs each
								// check the peer once per tick: 6 calls = 3 ticks.
								if isConnectedCalls.Add(1) >= 6 && !tc.peerConnected {
									signal()
								}
								return tc.peerConnected
							},
							ConnectedAddrsFunc: func() []string { return nil },
						}
					},
					DoMultiFunc: func(_ context.Context, targets []string, _ func(context.Context, string, MirrorClient, ...grpc.CallOption) error) error {
						doMultiCalls.Add(1)
						signal()
						if len(targets) == 0 {
							return errors.ErrTargetNotFound
						}
						return nil
					},
				},
			}
			m.addrs.Store(peerAddr, struct{}{})

			ech := m.Start(egctx)
			select {
			case <-reached:
			case <-ctx.Done():
				tt.Fatal("mirror register loop did not progress before timeout")
			}
			cancel()

			var errs []error
			for done := false; !done; {
				select {
				case err, ok := <-ech:
					if !ok {
						done = true
						break
					}
					errs = append(errs, err)
				case <-time.After(10 * time.Second):
					tt.Fatal("error channel was not closed after cancel")
				}
			}
			if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				tt.Errorf("unexpected errgroup error: %v", err)
			}

			if tc.peerConnected {
				if doMultiCalls.Load() == 0 {
					tt.Error("register was not called although a peer is connected")
				}
				return
			}
			if n := doMultiCalls.Load(); n != 0 {
				tt.Errorf("register called %d times without connected peers", n)
			}
			if len(errs) != 0 {
				tt.Errorf("errors reported without connected peers: %v", errs)
			}
		})
	}
}
