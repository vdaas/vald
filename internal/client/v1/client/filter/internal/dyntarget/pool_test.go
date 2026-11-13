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

package dyntarget

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

// setDown makes addr unreachable while keeping any pool already open for it,
// like grpc Connect leaving the old unhealthy pool in place on failure.
func (ct *connTracker) setDown(addr string) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.unreachable[addr] = true
}

func (ct *connTracker) state(addr string) (connected bool, disconnects int) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	return ct.conns[addr], ct.disconnects[addr]
}

func (ct *connTracker) connectCount(addr string) int {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	return ct.connects[addr]
}

// pending reports the entry of a request-derived addr and whether an idle
// disconnect is scheduled for it.
func (p *Pool) pending(addr string) (dt *dynamicTarget, gen uint64, scheduled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	dt = p.dyn[addr]
	if dt == nil {
		return nil, 0, false
	}
	return dt, dt.gen, dt.evict != nil
}

func newTestPool(ct *connTracker, idle time.Duration, configured ...string) *Pool {
	return New(ct.mock(), configured, idle)
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

// TestPool_Release pins that a filter target named in a request (not one of
// the configured addrs) keeps its connection after its last user releases it
// and only gets an idle disconnect scheduled, while configured targets never
// do. The idle timeout is long so no scheduled disconnect runs here.
func TestPool_Release(t *testing.T) {
	t.Parallel()
	const (
		configured = "configured:8081"
		dynamic    = "dynamic:8081"
		dynamic2   = "dynamic2:8081"
		down       = "down:8081"
	)
	type want struct {
		connected   map[string]bool
		connects    map[string]int
		disconnects map[string]int
		scheduled   map[string]bool
	}
	tests := []struct {
		run  func(context.Context, *Pool) error
		want want
		name string
	}{
		{
			name: "request-derived target stays connected with an idle disconnect scheduled after release",
			run: func(ctx context.Context, c *Pool) error {
				if err := c.Acquire(ctx, dynamic); err != nil {
					return err
				}
				c.Release(ctx, dynamic)
				return nil
			},
			want: want{
				connected:   map[string]bool{dynamic: true},
				disconnects: map[string]int{dynamic: 0},
				scheduled:   map[string]bool{dynamic: true},
			},
		},
		{
			name: "configured target stays connected without an idle disconnect after release",
			run: func(ctx context.Context, c *Pool) error {
				if err := c.Acquire(ctx, configured, dynamic); err != nil {
					return err
				}
				c.Release(ctx, configured, dynamic)
				return nil
			},
			want: want{
				connected:   map[string]bool{configured: true, dynamic: true},
				disconnects: map[string]int{configured: 0, dynamic: 0},
				scheduled:   map[string]bool{configured: false, dynamic: true},
			},
		},
		{
			name: "request-derived target stays connected while another user holds it",
			run: func(ctx context.Context, c *Pool) error {
				for range 2 {
					if err := c.Acquire(ctx, dynamic); err != nil {
						return err
					}
				}
				c.Release(ctx, dynamic)
				return nil
			},
			want: want{
				connected:   map[string]bool{dynamic: true},
				disconnects: map[string]int{dynamic: 0},
				scheduled:   map[string]bool{dynamic: false},
			},
		},
		{
			name: "hot request-derived target reuses its connection across requests",
			run: func(ctx context.Context, c *Pool) error {
				for range 3 {
					if err := c.Acquire(ctx, dynamic); err != nil {
						return err
					}
					c.Release(ctx, dynamic)
				}
				return nil
			},
			want: want{
				connected:   map[string]bool{dynamic: true},
				connects:    map[string]int{dynamic: 1},
				disconnects: map[string]int{dynamic: 0},
				scheduled:   map[string]bool{dynamic: true},
			},
		},
		{
			name: "targets acquired before a failed connect are released",
			run: func(ctx context.Context, c *Pool) error {
				if err := c.Acquire(ctx, dynamic, dynamic2, down); err == nil {
					return errors.New("Acquire of an unreachable addr succeeded")
				}
				return nil
			},
			want: want{
				connected:   map[string]bool{dynamic: true, dynamic2: true, down: false},
				disconnects: map[string]int{dynamic: 0, dynamic2: 0, down: 0},
				scheduled:   map[string]bool{dynamic: true, dynamic2: true, down: true},
			},
		},
		{
			name: "release of an unknown addr is a no-op",
			run: func(ctx context.Context, c *Pool) error {
				c.Release(ctx, dynamic)
				return nil
			},
			want: want{
				disconnects: map[string]int{dynamic: 0},
				scheduled:   map[string]bool{dynamic: false},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			ct := newConnTracker(down)
			c := newTestPool(ct, time.Hour, configured)
			if err := tc.run(tt.Context(), c); err != nil {
				tt.Fatal(err)
			}
			for addr, want := range tc.want.connected {
				if got, _ := ct.state(addr); got != want {
					tt.Errorf("%s connected = %v, want %v", addr, got, want)
				}
			}
			for addr, want := range tc.want.connects {
				if got := ct.connectCount(addr); got != want {
					tt.Errorf("%s Connect calls = %d, want %d", addr, got, want)
				}
			}
			for addr, want := range tc.want.disconnects {
				if _, got := ct.state(addr); got != want {
					tt.Errorf("%s Disconnect calls = %d, want %d", addr, got, want)
				}
			}
			for addr, want := range tc.want.scheduled {
				if _, _, got := c.pending(addr); got != want {
					tt.Errorf("%s idle disconnect scheduled = %v, want %v", addr, got, want)
				}
			}
			c.Stop()
		})
	}
}

// TestPool_disconnectIdle runs the scheduled idle disconnect directly to pin
// that it disconnects only a target that stayed unused since it was scheduled
// on a running pool.
func TestPool_disconnectIdle(t *testing.T) {
	t.Parallel()
	const addr = "dynamic:8081"
	type want struct {
		disconnects int
		connected   bool
		kept        bool // the entry is still tracked
	}
	tests := []struct {
		// run returns the generation whose idle disconnect is then run.
		run  func(context.Context, *Pool) (uint64, error)
		name string
		want want
	}{
		{
			name: "target unused for the idle timeout is disconnected",
			run: func(ctx context.Context, c *Pool) (uint64, error) {
				if err := c.Acquire(ctx, addr); err != nil {
					return 0, err
				}
				c.Release(ctx, addr)
				_, gen, _ := c.pending(addr)
				return gen, nil
			},
			want: want{connected: false, disconnects: 1, kept: false},
		},
		{
			name: "target used again during the idle window keeps its connection",
			run: func(ctx context.Context, c *Pool) (uint64, error) {
				if err := c.Acquire(ctx, addr); err != nil {
					return 0, err
				}
				c.Release(ctx, addr)
				_, gen, _ := c.pending(addr)
				if err := c.Acquire(ctx, addr); err != nil {
					return 0, err
				}
				return gen, nil
			},
			want: want{connected: true, disconnects: 0, kept: true},
		},
		{
			name: "superseded idle disconnect of a re-released target is skipped",
			run: func(ctx context.Context, c *Pool) (uint64, error) {
				if err := c.Acquire(ctx, addr); err != nil {
					return 0, err
				}
				c.Release(ctx, addr)
				_, gen, _ := c.pending(addr)
				if err := c.Acquire(ctx, addr); err != nil {
					return 0, err
				}
				c.Release(ctx, addr)
				return gen, nil
			},
			want: want{connected: true, disconnects: 0, kept: true},
		},
		{
			name: "stopped pool does not disconnect",
			run: func(ctx context.Context, c *Pool) (uint64, error) {
				if err := c.Acquire(ctx, addr); err != nil {
					return 0, err
				}
				c.Release(ctx, addr)
				_, gen, _ := c.pending(addr)
				c.Stop()
				return gen, nil
			},
			want: want{connected: true, disconnects: 0, kept: true},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			ct := newConnTracker()
			c := newTestPool(ct, time.Hour)
			gen, err := tc.run(tt.Context(), c)
			if err != nil {
				tt.Fatal(err)
			}
			dt, _, _ := c.pending(addr)
			if dt == nil {
				tt.Fatal("no entry tracked for the request-derived target")
			}
			c.disconnectIdle(tt.Context(), addr, dt, gen)
			connected, disconnects := ct.state(addr)
			if connected != tc.want.connected {
				tt.Errorf("connected = %v, want %v", connected, tc.want.connected)
			}
			if disconnects != tc.want.disconnects {
				tt.Errorf("Disconnect calls = %d, want %d", disconnects, tc.want.disconnects)
			}
			if got, _, _ := c.pending(addr); (got != nil) != tc.want.kept {
				tt.Errorf("entry kept = %v, want %v", got != nil, tc.want.kept)
			}
			c.Stop()
		})
	}
}

// TestPool_Release_idleTimer pins that the scheduled idle disconnect fires
// on its own once the idle timeout passes.
func TestPool_Release_idleTimer(t *testing.T) {
	t.Parallel()
	const addr = "dynamic:8081"
	ct := newConnTracker()
	c := newTestPool(ct, time.Millisecond)
	defer c.Stop()
	if err := c.Acquire(t.Context(), addr); err != nil {
		t.Fatal(err)
	}
	c.Release(t.Context(), addr)
	waitFor(t, "the idle disconnect", func() bool {
		connected, disconnects := ct.state(addr)
		return !connected && disconnects == 1
	})
	waitFor(t, "the entry to be dropped", func() bool {
		dt, _, _ := c.pending(addr)
		return dt == nil
	})
}

// TestPool_Start_stopsIdleDisconnects pins that pending idle disconnects are
// cancelled when the ctx passed to Start is done, so no timer outlives it.
func TestPool_Start_stopsIdleDisconnects(t *testing.T) {
	t.Parallel()
	const addr = "dynamic:8081"
	ct := newConnTracker()
	c := newTestPool(ct, time.Hour)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	c.Start(ctx)
	if err := c.Acquire(t.Context(), addr); err != nil {
		t.Fatal(err)
	}
	c.Release(t.Context(), addr)
	if _, _, scheduled := c.pending(addr); !scheduled {
		t.Fatal("no idle disconnect scheduled after release")
	}
	cancel()
	waitFor(t, "the idle disconnect to be cancelled", func() bool {
		_, _, scheduled := c.pending(addr)
		return !scheduled
	})
	// A release after the stop schedules nothing.
	if err := c.Acquire(t.Context(), addr); err != nil {
		t.Fatal(err)
	}
	c.Release(t.Context(), addr)
	if _, _, scheduled := c.pending(addr); scheduled {
		t.Error("idle disconnect scheduled after the Start ctx was done")
	}
	if connected, disconnects := ct.state(addr); !connected || disconnects != 0 {
		t.Errorf("connected = %v, Disconnect calls = %d, want true, 0", connected, disconnects)
	}
}

// TestPool_Release_concurrent pins that a request-derived target is never
// disconnected while a concurrent user still holds it, even with idle
// disconnects firing between uses, and that it is disconnected once all users
// are done.
func TestPool_Release_concurrent(t *testing.T) {
	t.Parallel()
	const (
		addr    = "dynamic:8081"
		workers = 8
		rounds  = 200
	)
	ct := newConnTracker()
	c := newTestPool(ct, time.Microsecond, "configured:8081")
	defer c.Stop()
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	errs := make(chan error, workers)
	for range workers {
		go func() {
			for range rounds {
				if ctx.Err() != nil {
					errs <- ctx.Err()
					return
				}
				if err := c.Acquire(ctx, addr); err != nil {
					errs <- err
					return
				}
				if connected, _ := ct.state(addr); !connected {
					errs <- errors.New("target disconnected while held")
					return
				}
				c.Release(ctx, addr)
			}
			errs <- nil
		}()
	}
	for range workers {
		select {
		case err := <-errs:
			if err != nil {
				t.Fatal(err)
			}
		case <-ctx.Done():
			t.Fatal("timed out waiting for the workers")
		}
	}
	waitFor(t, "the target to be disconnected after every user released it", func() bool {
		connected, _ := ct.state(addr)
		return !connected
	})
}

// TestPool_acquire_failedReconnectKeepsIdleDisconnect pins that a target
// whose pending idle disconnect is cancelled by a new request still gets one
// when that request's reconnect fails; otherwise the pool left behind for it is
// never disconnected.
func TestPool_acquire_failedReconnectKeepsIdleDisconnect(t *testing.T) {
	t.Parallel()
	const dynamic = "dynamic:8081"
	ct := newConnTracker()
	c := newTestPool(ct, time.Hour)
	defer c.Stop()
	if err := c.Acquire(t.Context(), dynamic); err != nil {
		t.Fatal(err)
	}
	c.Release(t.Context(), dynamic)
	if _, _, scheduled := c.pending(dynamic); !scheduled {
		t.Fatal("no idle disconnect scheduled after release")
	}
	ct.setDown(dynamic)
	if err := c.Acquire(t.Context(), dynamic); err == nil {
		t.Fatal("Acquire of a down addr succeeded")
	}
	dt, _, scheduled := c.pending(dynamic)
	if dt == nil || !scheduled {
		t.Errorf("entry present = %v, idle disconnect scheduled = %v, want true, true", dt != nil, scheduled)
	}
	if dt != nil && dt.refs != 0 {
		t.Errorf("refs = %d, want 0", dt.refs)
	}
}
