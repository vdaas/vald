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

package grpc

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/net/grpc/status"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// recordingBreaker is a circuitbreaker.CircuitBreaker that runs fn and records
// the key and whether fn reported its error as breaker-ignorable.
type recordingBreaker struct {
	keys      []string
	ignorable []bool
}

func (b *recordingBreaker) Do(
	ctx context.Context, key string, fn func(ctx context.Context) (any, error),
) (any, error) {
	b.keys = append(b.keys, key)
	val, err := fn(ctx)
	var igerr *errors.ErrCircuitBreakerIgnorable
	isIgnorable := errors.As(err, &igerr)
	b.ignorable = append(b.ignorable, isIgnorable)
	if isIgnorable {
		return nil, igerr.Unwrap()
	}
	return val, err
}

// plainConnPool is a non-IP, unhealthy fixedConnPool so executeRPC's error
// classification can run without a real pool.
type plainConnPool struct {
	fixedConnPool
}

func (plainConnPool) IsIPConn() bool                 { return false }
func (plainConnPool) IsHealthy(context.Context) bool { return false }

// Test_gRPCClient_do_circuitBreakerWithoutBackoff verifies that do applies the
// circuit breaker even when no backoff is configured, with the same key and
// ignorable-error semantics as the backoff path and RoundRobin.
func Test_gRPCClient_do_circuitBreakerWithoutBackoff(t *testing.T) {
	t.Parallel()
	const addr = "10.0.0.1:8081"
	rpcErr := status.Error(codes.InvalidArgument, "bad request")
	tests := []struct {
		want          any
		wantErr       error
		f             func(context.Context, *ClientConn, ...CallOption) (any, error)
		name          string
		wantIgnorable bool
	}{
		{
			name: "success goes through the breaker",
			f: func(context.Context, *ClientConn, ...CallOption) (any, error) {
				return "ok", nil
			},
			want: "ok",
		},
		{
			name: "non-retryable error is marked ignorable for the breaker",
			f: func(context.Context, *ClientConn, ...CallOption) (any, error) {
				return nil, rpcErr
			},
			wantErr:       rpcErr,
			wantIgnorable: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cb := new(recordingBreaker)
			g, ok := New("test", WithCircuitBreaker(cb)).(*gRPCClient)
			if !ok {
				t.Fatal("New did not return *gRPCClient")
			}
			if g.bo != nil {
				t.Fatal("backoff must be unset for this test")
			}
			conn, err := ggrpc.NewClient("passthrough:///do-cb-test",
				ggrpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if cerr := conn.Close(); cerr != nil {
					t.Error(cerr)
				}
			})
			p := plainConnPool{fixedConnPool{conn: conn}}
			// The breaker applies only to a pool that is still connected.
			g.conns.Store(addr, p)
			got, err := g.do(t.Context(), p, addr, tc.f)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("do() err = %v", err)
				}
			} else if !errors.Is(err, tc.wantErr) {
				t.Fatalf("do() err = %v, want wrapping %v", err, tc.wantErr)
			}
			if tc.wantErr == nil && got != tc.want {
				t.Errorf("do() = %v, want %v", got, tc.want)
			}
			if len(cb.keys) != 1 || cb.keys[0] != addr {
				t.Fatalf("breaker keys = %v, want [%s]", cb.keys, addr)
			}
			if cb.ignorable[0] != tc.wantIgnorable {
				t.Errorf("breaker saw ignorable = %v, want %v", cb.ignorable[0], tc.wantIgnorable)
			}
		})
	}
}

// Test_gRPCClient_do_skipsBreakerForDisconnectedPool verifies that a call on a
// pool that is no longer connected to addr does not go through the breaker. A
// call that obtained the pool before Disconnect evicted the breakers of addr
// used to re-create one that nothing evicted again.
func Test_gRPCClient_do_skipsBreakerForDisconnectedPool(t *testing.T) {
	t.Parallel()
	const addr = "10.0.0.2:8081"
	tests := []struct {
		stored    func(p pool.Conn) pool.Conn
		name      string
		wantBreak bool
	}{
		{
			name:      "connected pool goes through the breaker",
			stored:    func(p pool.Conn) pool.Conn { return p },
			wantBreak: true,
		},
		{
			name:   "disconnected pool skips the breaker",
			stored: func(pool.Conn) pool.Conn { return nil },
		},
		{
			name:   "pool replaced by a successor skips the breaker",
			stored: func(pool.Conn) pool.Conn { return &hookConn{} },
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cb := new(recordingBreaker)
			g, ok := New("test", WithCircuitBreaker(cb)).(*gRPCClient)
			if !ok {
				t.Fatal("New did not return *gRPCClient")
			}
			conn, err := ggrpc.NewClient("passthrough:///do-cb-evict-test",
				ggrpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if cerr := conn.Close(); cerr != nil {
					t.Error(cerr)
				}
			})
			p := plainConnPool{fixedConnPool{conn: conn}}
			if cur := tc.stored(p); cur != nil {
				g.conns.Store(addr, cur)
			}
			got, err := g.do(t.Context(), p, addr, func(context.Context, *ClientConn, ...CallOption) (any, error) {
				return "ok", nil
			})
			if err != nil {
				t.Fatalf("do() err = %v", err)
			}
			if got != "ok" {
				t.Errorf("do() = %v, want ok", got)
			}
			if gotBreak := len(cb.keys) != 0; gotBreak != tc.wantBreak {
				t.Errorf("breaker keys = %v, want breaker used = %v", cb.keys, tc.wantBreak)
			}
		})
	}
}
