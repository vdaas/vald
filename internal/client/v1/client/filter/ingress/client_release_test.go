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

package ingress

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/sync"
	grpcmock "github.com/vdaas/vald/internal/test/mock/grpc"
)

// connTracker is a fake gRPC client connection table.
type connTracker struct {
	conns       map[string]bool
	unreachable map[string]bool
	connects    map[string]int
	disconnects map[string]int
	mu          sync.Mutex
}

func newConnTracker(unreachable ...string) *connTracker {
	ct := &connTracker{
		conns:       make(map[string]bool),
		unreachable: make(map[string]bool),
		connects:    make(map[string]int),
		disconnects: make(map[string]int),
	}
	for _, addr := range unreachable {
		ct.unreachable[addr] = true
	}
	return ct
}

func (ct *connTracker) mock() *grpcmock.GRPCClientMock {
	return &grpcmock.GRPCClientMock{
		ConnectFunc: func(_ context.Context, addr string, _ ...grpc.DialOption) (pool.Conn, error) {
			ct.mu.Lock()
			defer ct.mu.Unlock()
			if ct.unreachable[addr] {
				return nil, errors.New("unreachable")
			}
			ct.connects[addr]++
			ct.conns[addr] = true
			return nil, nil
		},
		DisconnectFunc: func(_ context.Context, addr string) error {
			ct.mu.Lock()
			defer ct.mu.Unlock()
			ct.disconnects[addr]++
			if !ct.conns[addr] {
				return errors.ErrGRPCClientConnNotFound(addr)
			}
			delete(ct.conns, addr)
			return nil
		},
		IsConnectedFunc: func(_ context.Context, addr string) bool {
			ct.mu.Lock()
			defer ct.mu.Unlock()
			// a pool left behind for an addr that went down is not connected.
			return ct.conns[addr] && !ct.unreachable[addr]
		},
	}
}

func (ct *connTracker) state(addr string) (connected bool, disconnects int) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	return ct.conns[addr], ct.disconnects[addr]
}

// waitFor polls cond until it holds; it fails the test after a bounded time.
func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()
	for !cond() {
		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for %s", what)
		case <-tick.C:
		}
	}
}

func TestWithIdleTimeout(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		dur  string
		want time.Duration
	}{
		{name: "unset keeps the default", dur: "", want: defaultIdleTimeout},
		{name: "valid duration is used", dur: "5s", want: 5 * time.Second},
		{name: "invalid duration keeps the default", dur: "invalid", want: defaultIdleTimeout},
		{name: "non-positive duration keeps the default", dur: "-1s", want: defaultIdleTimeout},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			c, err := New(WithAddrs("configured:8081"), WithClient(newConnTracker().mock()), WithIdleTimeout(tc.dur))
			if err != nil {
				tt.Fatal(err)
			}
			cl, ok := c.(*client)
			if !ok {
				tt.Fatalf("New returned %T, want *client", c)
			}
			if cl.idle != tc.want {
				tt.Errorf("idle = %v, want %v", cl.idle, tc.want)
			}
		})
	}
}

// Test_client_Target_Release pins that Target and Release go through the
// request-derived target pool with the configured idle timeout: a released
// request-derived target is disconnected after it, configured targets stay
// connected, and the targets acquired before a failed connect are released.
// The pool itself is tested in package dyntarget.
func Test_client_Target_Release(t *testing.T) {
	t.Parallel()
	const (
		configured = "configured:8081"
		dynamic    = "dynamic:8081"
		down       = "down:8081"
	)
	ct := newConnTracker(down)
	c, err := New(WithAddrs(configured), WithClient(ct.mock()), WithIdleTimeout("1ms"))
	if err != nil {
		t.Fatal(err)
	}
	cl, ok := c.(*client)
	if !ok {
		t.Fatalf("New returned %T, want *client", c)
	}
	defer cl.targets.Stop()
	ctx := t.Context()
	if _, err = c.Target(ctx); !errors.Is(err, errors.ErrTargetNotFound) {
		t.Errorf("Target() error = %v, want %v", err, errors.ErrTargetNotFound)
	}
	fc, err := c.Target(ctx, configured, dynamic)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := fc.(*multipleAddrsClient); !ok {
		t.Errorf("Target returned %T, want *multipleAddrsClient", fc)
	}
	c.Release(ctx, configured, dynamic)
	waitFor(t, "the released request-derived target to be disconnected", func() bool {
		connected, disconnects := ct.state(dynamic)
		return !connected && disconnects == 1
	})
	if connected, disconnects := ct.state(configured); !connected || disconnects != 0 {
		t.Errorf("configured connected = %v, Disconnect calls = %d, want true, 0", connected, disconnects)
	}
	if _, err = c.Target(ctx, dynamic, down); err == nil {
		t.Fatal("Target to an unreachable addr succeeded")
	}
	waitFor(t, "the target acquired before the failed connect to be disconnected", func() bool {
		connected, disconnects := ct.state(dynamic)
		return !connected && disconnects == 2
	})
}

// lifecycleClient adds the connection monitor and Close that the client's
// Start and Stop need to the mock.
type lifecycleClient struct {
	*grpcmock.GRPCClientMock
}

func (lifecycleClient) StartConnectionMonitor(context.Context) (<-chan error, error) {
	ech := make(chan error)
	close(ech)
	return ech, nil
}

func (lifecycleClient) Close(context.Context) error { return nil }

// releaseKeepsConnection reports whether a released request-derived target is
// still connected, with no further Disconnect, well after the idle timeout.
func releaseKeepsConnection(
	ctx context.Context, t *testing.T, c Client, ct *connTracker, addr string,
) bool {
	t.Helper()
	_, before := ct.state(addr)
	if _, err := c.Target(ctx, addr); err != nil {
		t.Fatal(err)
	}
	c.Release(ctx, addr)
	time.Sleep(50 * time.Millisecond) // 50x the 1ms idle timeout
	connected, after := ct.state(addr)
	return connected && after == before
}

// Test_client_StartStop_wireIdleDisconnects pins that the client's Start ties
// the request-derived target pool to its ctx and that Stop stops it: once
// either has happened, a released request-derived target is no longer
// disconnected after the idle timeout.
func Test_client_StartStop_wireIdleDisconnects(t *testing.T) {
	t.Parallel()
	const dynamic = "dynamic:8081"
	newClient := func(t *testing.T) (Client, *connTracker) {
		t.Helper()
		ct := newConnTracker()
		c, err := New(WithAddrs("configured:8081"),
			WithClient(lifecycleClient{ct.mock()}), WithIdleTimeout("1ms"))
		if err != nil {
			t.Fatal(err)
		}
		return c, ct
	}
	t.Run("Start ctx done", func(t *testing.T) {
		t.Parallel()
		c, ct := newClient(t)
		ctx, cancel := context.WithCancel(t.Context())
		if _, err := c.Start(ctx); err != nil {
			t.Fatal(err)
		}
		cancel()
		// Start stops the pool asynchronously (context.AfterFunc) once ctx is done.
		waitFor(t, "a released target to stay connected after the Start ctx is done", func() bool {
			return releaseKeepsConnection(t.Context(), t, c, ct, dynamic)
		})
	})
	t.Run("Stop", func(t *testing.T) {
		t.Parallel()
		c, ct := newClient(t)
		if err := c.Stop(t.Context()); err != nil {
			t.Fatal(err)
		}
		if !releaseKeepsConnection(t.Context(), t, c, ct, dynamic) {
			t.Error("a released target was disconnected after Stop")
		}
	})
}
