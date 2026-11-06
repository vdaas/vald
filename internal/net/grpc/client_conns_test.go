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
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/sync/atomic"
)

// hookConn is a pool.Conn fake whose Connect runs a hook, so that a test can
// interleave a Disconnect at a fixed point of gRPCClient.Connect.
type hookConn struct {
	onConnect func(ctx context.Context, c *hookConn) (pool.Conn, error)
	healthy   atomic.Bool
}

func (c *hookConn) Connect(ctx context.Context) (pool.Conn, error) { return c.onConnect(ctx, c) }
func (*hookConn) Disconnect(context.Context) error                 { return nil }
func (*hookConn) Do(context.Context, func(*pool.ClientConn) error) error {
	return errors.ErrGRPCClientConnNotFound("*")
}
func (*hookConn) Get(context.Context) (*pool.ClientConn, bool) { return nil, false }
func (c *hookConn) IsHealthy(context.Context) bool             { return c.healthy.Load() }
func (*hookConn) IsIPConn() bool                               { return true }
func (*hookConn) Len() uint64                                  { return 0 }
func (*hookConn) Size() uint64                                 { return 0 }
func (c *hookConn) Reconnect(ctx context.Context, _ bool) (pool.Conn, error) {
	return c.Connect(ctx)
}
func (*hookConn) String() string { return "hookConn" }

// TestGRPCClient_Connect_concurrentDisconnect verifies that g.conns and
// clientCount stay consistent when Disconnect removes the unhealthy pool while
// Connect is reconnecting it: Connect used to store the closed pool back into
// g.conns, and to skip counting a new pool because its stale Load still saw
// the removed one.
func TestGRPCClient_Connect_concurrentDisconnect(t *testing.T) {
	t.Parallel()
	const addr = "127.0.0.10:1"
	tests := []struct {
		// reconnectErr is returned by the old pool's Connect after Disconnect ran.
		reconnectErr error
		name         string
		wantCount    uint64
		wantStored   bool
	}{
		{
			name:         "reconnected old pool is not stored back after Disconnect",
			reconnectErr: nil,
			wantStored:   false,
			wantCount:    0,
		},
		{
			name:         "new pool stored after Disconnect is counted",
			reconnectErr: errors.New("reconnect failed"),
			wantStored:   true,
			wantCount:    1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			c := New("connect-disconnect-test", WithInsecure(true), WithConnectionPoolSize(1))
			g, ok := c.(*gRPCClient)
			if !ok {
				t.Fatalf("unexpected client type: %T", c)
			}
			old := &hookConn{
				onConnect: func(ctx context.Context, hc *hookConn) (pool.Conn, error) {
					if err := g.Disconnect(ctx, addr); err != nil {
						t.Errorf("Disconnect returned error: %v", err)
					}
					if test.reconnectErr != nil {
						return nil, test.reconnectErr
					}
					hc.healthy.Store(true)
					return hc, nil
				},
			}
			g.conns.Store(addr, old)
			atomic.AddUint64(&g.clientCount, 1)

			if _, err := g.Connect(ctx, addr); err != nil {
				t.Fatalf("Connect returned error: %v", err)
			}

			got, stored := g.conns.Load(addr)
			if stored != test.wantStored {
				t.Errorf("conn stored = %v, want %v", stored, test.wantStored)
			}
			if stored && got == pool.Conn(old) {
				t.Error("the pool removed by Disconnect must not be stored back")
			}
			if cnt := atomic.LoadUint64(&g.clientCount); cnt != test.wantCount {
				t.Errorf("clientCount = %d, want %d", cnt, test.wantCount)
			}
			if stored {
				if err := g.Disconnect(ctx, addr); err != nil {
					t.Errorf("cleanup Disconnect returned error: %v", err)
				}
			}
		})
	}
}
