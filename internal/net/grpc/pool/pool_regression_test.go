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
package pool

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"google.golang.org/grpc"
	gbackoff "google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// onceBackoff is a backoff.Backoff stub that runs the function exactly once
// and returns its result as is.
type onceBackoff struct{}

func (onceBackoff) Do(
	ctx context.Context, f func(ctx context.Context) (any, bool, error),
) (any, error) {
	res, _, err := f(ctx)
	return res, err
}

func (onceBackoff) Close() {}

// TestPool_dial_backoffError verifies that dial with a backoff strategy returns
// the dial error instead of (nil, nil) when the dial fails and no connection exists.
func TestPool_dial_backoffError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
	}{
		{name: "dial failure with backoff returns error"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			// no dial option (transport credentials) is set, so grpc.NewClient always fails.
			p := &pool{bo: onceBackoff{}}
			p.poolSize.Store(1)
			p.init()

			conn, err := p.dial(ctx, 0, "127.0.0.1:1")
			if err == nil {
				t.Errorf("dial must return an error, got conn: %v", conn)
			}
			if conn != nil {
				t.Errorf("dial must not return a connection on error, got: %v", conn)
			}
		})
	}
}

// TestPool_refreshConn_dialFailure verifies that a failed re-dial never stores a poolConn
// with a nil connection into the slot.
func TestPool_refreshConn_dialFailure(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	p := &pool{bo: onceBackoff{}, addr: "127.0.0.1:1"}
	p.poolSize.Store(1)
	p.init()
	eg, ctx := errgroup.New(ctx)
	p.errGroup = eg

	p.refreshConn(ctx, eg, 0, nil, p.addr)
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
	if pc := p.getSlots()[0].Load(); pc != nil {
		t.Errorf("slot must stay empty when dial fails, got: %+v", pc)
	}
}

// TestPool_IsHealthy_metrics verifies that Metrics reports the latest healthy
// connection count per address and does not grow on repeated health checks.
func TestPool_IsHealthy_metrics(t *testing.T) {
	t.Parallel()
	type want struct {
		poolAddr string
		slotAddr string
		poolCnt  uint64
		slotCnt  uint64
	}
	tests := []struct {
		name  string
		want  want
		calls int
	}{
		{
			name:  "single health check",
			calls: 1,
			want: want{
				poolAddr: "metrics-test-pool-1.invalid:8081",
				poolCnt:  2,
				slotAddr: "127.0.0.1:1",
				slotCnt:  2,
			},
		},
		{
			name:  "repeated health checks do not accumulate",
			calls: 5,
			want: want{
				poolAddr: "metrics-test-pool-2.invalid:8081",
				poolCnt:  2,
				slotAddr: "127.0.0.2:1",
				slotCnt:  2,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			p := &pool{addr: test.want.poolAddr, enableMetrics: true}
			p.poolSize.Store(2)
			p.init()
			for i := range p.getSlots() {
				conn, err := grpc.NewClient(test.want.slotAddr,
					grpc.WithTransportCredentials(insecure.NewCredentials()))
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := conn.Close(); err != nil {
						t.Errorf("failed to close connection: %v", err)
					}
				})
				p.getSlots()[i].Store(&poolConn{conn: conn, addr: test.want.slotAddr})
			}

			for range test.calls {
				if !p.IsHealthy(ctx) {
					t.Fatal("pool must be healthy")
				}
			}

			got := Metrics(ctx)
			if got[test.want.poolAddr] != test.want.poolCnt {
				t.Errorf("pool addr count = %d, want %d", got[test.want.poolAddr], test.want.poolCnt)
			}
			if got[test.want.slotAddr] != test.want.slotCnt {
				t.Errorf("slot addr count = %d, want %d", got[test.want.slotAddr], test.want.slotCnt)
			}
		})
	}
}

// TestPool_IsHealthy_metrics_reset verifies that the gauge of an address whose connections
// all became unhealthy is reset to 0 instead of keeping the last non-zero value.
func TestPool_IsHealthy_metrics_reset(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	const (
		poolAddr = "metrics-test-pool-3.invalid:8081"
		liveAddr = "127.0.0.3:1"
		downAddr = "127.0.0.4:1"
	)
	p := &pool{addr: poolAddr, enableMetrics: true}
	p.poolSize.Store(2)
	p.init()
	conns := make([]*ClientConn, 0, 2)
	for i, addr := range []string{liveAddr, downAddr} {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatal(err)
		}
		conns = append(conns, conn)
		p.getSlots()[i].Store(&poolConn{conn: conn, addr: addr})
	}
	t.Cleanup(func() {
		if err := conns[0].Close(); err != nil {
			t.Errorf("failed to close connection: %v", err)
		}
	})

	if !p.IsHealthy(ctx) {
		t.Fatal("pool must be healthy")
	}
	if got := Metrics(ctx); got[liveAddr] != 1 || got[downAddr] != 1 {
		t.Fatalf("metrics before shutdown = %v, want %s=1 and %s=1", got, liveAddr, downAddr)
	}

	// a closed connection is in Shutdown state and reported as unhealthy.
	if err := conns[1].Close(); err != nil {
		t.Fatal(err)
	}
	if !p.IsHealthy(ctx) {
		t.Fatal("pool must stay healthy with one live connection")
	}
	got := Metrics(ctx)
	if got[liveAddr] != 1 {
		t.Errorf("live addr count = %d, want 1", got[liveAddr])
	}
	if got[downAddr] != 0 {
		t.Errorf("down addr count = %d, want 0", got[downAddr])
	}
	if got[poolAddr] != 1 {
		t.Errorf("pool addr count = %d, want 1", got[poolAddr])
	}
}

// TestPool_probeGRPCPorts_ctxDone verifies that a done context surfaces as an error
// instead of (0, nil), which Connect would take as a found port and dial port 0.
func TestPool_probeGRPCPorts_ctxDone(t *testing.T) {
	t.Parallel()
	tests := []struct {
		ctx     func(t *testing.T) context.Context
		name    string
		wantErr error
		ports   []uint16
	}{
		{
			name: "canceled context returns context.Canceled",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				return ctx
			},
			ports:   []uint16{1},
			wantErr: context.Canceled,
		},
		{
			name: "expired deadline returns context.DeadlineExceeded",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
				t.Cleanup(cancel)
				return ctx
			},
			ports:   []uint16{1},
			wantErr: context.DeadlineExceeded,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			p := &pool{host: "127.0.0.1"}
			port, err := p.probeGRPCPorts(test.ctx(t), test.ports)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("err = %v, want %v", err, test.wantErr)
			}
			if port != 0 {
				t.Errorf("port = %d, want 0", port)
			}
		})
	}
}

// TestPool_probeGRPCPorts_noPort verifies that having no port to probe keeps returning
// the invalid port error, which is not a context error.
func TestPool_probeGRPCPorts_noPort(t *testing.T) {
	t.Parallel()
	p := &pool{host: "127.0.0.1"}
	port, err := p.probeGRPCPorts(t.Context(), nil)
	if err == nil {
		t.Fatal("err = nil, want an invalid gRPC port error")
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, must not be a context error", err)
	}
	if port != 0 {
		t.Errorf("port = %d, want 0", port)
	}
}

// load used to index slots[len(slots)] (out of range) when the pool size ran ahead
// of the slots snapshot during grow; slotAt clamps to the last slot instead.
func TestSlotAt(t *testing.T) {
	t.Parallel()
	pc := &poolConn{addr: "last"}
	slots := make([]atomic.Pointer[poolConn], 2)
	slots[1].Store(pc)
	tests := []struct {
		wantConn *poolConn
		name     string
		slots    []atomic.Pointer[poolConn]
		idx      uint64
		wantIdx  uint64
	}{
		{name: "empty", slots: nil, idx: 3, wantIdx: 0, wantConn: nil},
		{name: "in range", slots: slots, idx: 1, wantIdx: 1, wantConn: pc},
		{name: "equal to length is clamped", slots: slots, idx: 2, wantIdx: 1, wantConn: pc},
		{name: "beyond length is clamped", slots: slots, idx: 9, wantIdx: 1, wantConn: pc},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotIdx, gotConn := slotAt(tt.slots, tt.idx)
			if gotIdx != tt.wantIdx || gotConn != tt.wantConn {
				t.Errorf("slotAt(%d) = (%d, %v), want (%d, %v)", tt.idx, gotIdx, gotConn, tt.wantIdx, tt.wantConn)
			}
		})
	}
}

// newIdleTestConn returns an unconnected (Idle) client conn; tests either hand it to
// the pool under test or close it in Cleanup.
func newIdleTestConn(t *testing.T, addr string) *ClientConn {
	t.Helper()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if conn.GetState() != connectivity.Shutdown {
			if err := conn.Close(); err != nil {
				t.Errorf("failed to close connection: %v", err)
			}
		}
	})
	return conn
}

// waitGroupBounded waits for eg but fails the test instead of hanging forever.
func waitGroupBounded(t *testing.T, eg errgroup.Group) {
	t.Helper()
	done := make(chan error, 1)
	go func() {
		done <- eg.Wait()
	}()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for refreshConn")
	}
}

// TestPool_Disconnect_doneContext verifies that Disconnect closes every slot even when the
// caller's ctx is already done; it used to skip all of them and then flush the slots,
// leaking the conns.
func TestPool_Disconnect_doneContext(t *testing.T) {
	t.Parallel()
	tests := []struct {
		ctx  func(t *testing.T) context.Context
		name string
	}{
		{
			name: "canceled context",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				return ctx
			},
		},
		{
			name: "expired deadline",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
				t.Cleanup(cancel)
				return ctx
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			p := &pool{addr: "127.0.0.7:1", oldConnCloseDelay: time.Second}
			p.poolSize.Store(2)
			p.init()
			conns := make([]*ClientConn, 0, 2)
			for i := range p.getSlots() {
				conn := newIdleTestConn(t, p.addr)
				conns = append(conns, conn)
				p.getSlots()[i].Store(&poolConn{conn: conn, addr: p.addr})
			}

			if err := p.Disconnect(test.ctx(t)); err != nil {
				t.Errorf("Disconnect returned error: %v", err)
			}
			for i, conn := range conns {
				if st := conn.GetState(); st != connectivity.Shutdown {
					t.Errorf("conn %d state = %s, want %s", i, st, connectivity.Shutdown)
				}
			}
		})
	}
}

// TestPool_refreshConn_afterDisconnect verifies that a refresh whose dial completes after
// Disconnect does not leave its new conn in the flushed slots, where nobody would close it.
func TestPool_refreshConn_afterDisconnect(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	p := &pool{
		addr:              "127.0.0.8:1",
		oldConnCloseDelay: time.Second,
		dialOpts:          []DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	}
	p.poolSize.Store(2)
	p.init()
	if err := p.Disconnect(ctx); err != nil {
		t.Fatalf("Disconnect returned error: %v", err)
	}

	eg, ectx := errgroup.New(ctx)
	p.errGroup = eg
	for idx := range uint64(2) {
		p.refreshConn(ectx, eg, idx, nil, p.addr)
	}
	waitGroupBounded(t, eg)

	for i := range p.getSlots() {
		if pc := p.getSlots()[i].Load(); pc != nil {
			t.Errorf("slot %d must stay empty after Disconnect, got conn in state %s", i, pc.conn.GetState())
		}
	}
}

// countingGroup records how many tasks are handed to Go.
type countingGroup struct {
	errgroup.Group

	n atomic.Int64
}

func (g *countingGroup) Go(f func() error) {
	g.n.Add(1)
	g.Group.Go(f)
}

// TestPool_closedPoolDoesNotDial verifies that refreshes on a disconnected pool
// do not dial: the pool stays closed, so every dialed conn was only torn down again.
func TestPool_closedPoolDoesNotDial(t *testing.T) {
	t.Parallel()
	tests := []struct {
		run  func(t *testing.T, ctx context.Context, p *pool, eg errgroup.Group)
		name string
	}{
		{
			name: "refreshConn",
			run: func(t *testing.T, ctx context.Context, p *pool, eg errgroup.Group) {
				t.Helper()
				for idx := range uint64(2) {
					p.refreshConn(ctx, eg, idx, nil, p.addr)
				}
			},
		},
		{
			name: "IsHealthy",
			run: func(t *testing.T, ctx context.Context, p *pool, _ errgroup.Group) {
				t.Helper()
				if p.IsHealthy(ctx) {
					t.Error("a disconnected pool must not be healthy")
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			p := &pool{
				addr:              "127.0.0.9:1",
				oldConnCloseDelay: time.Second,
				dialOpts:          []DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
			}
			p.poolSize.Store(2)
			p.init()
			if err := p.Disconnect(ctx); err != nil {
				t.Fatalf("Disconnect returned error: %v", err)
			}

			eg, ectx := errgroup.New(ctx)
			cg := &countingGroup{Group: eg}
			p.errGroup = cg
			test.run(t, ectx, p, cg)
			waitGroupBounded(t, eg)

			if n := cg.n.Load(); n != 0 {
				t.Errorf("%d dials were started on a disconnected pool, want 0", n)
			}
		})
	}
}

// childTrackingCtx is a parent context implementing AfterFunc, so
// context.WithCancel* registers each child through AfterFunc and calls the
// returned stop when the child is canceled. A registration without a stop is a
// child that stays attached to the parent until the parent is done.
type childTrackingCtx struct {
	done       chan struct{}
	registered atomic.Int64
	stopped    atomic.Int64
}

func (*childTrackingCtx) Deadline() (time.Time, bool) { return time.Time{}, false }

func (c *childTrackingCtx) Done() <-chan struct{} { return c.done }

func (*childTrackingCtx) Err() error { return nil }

func (*childTrackingCtx) Value(any) any { return nil }

func (c *childTrackingCtx) AfterFunc(func()) func() bool {
	c.registered.Add(1)
	var once atomic.Bool
	return func() bool {
		if once.CompareAndSwap(false, true) {
			c.stopped.Add(1)
			return true
		}
		return false
	}
}

// TestPool_getHealthyConn_transientFailureNoCtxLeak verifies that returning a
// TransientFailure conn as the last resort does not leave a child context
// registered on the caller's ctx; errgroup.New's context is only canceled by Wait.
func TestPool_getHealthyConn_transientFailureNoCtxLeak(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	const addr = "127.0.0.1:1" // connection refused
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// a long backoff keeps the conn in TransientFailure for the whole test.
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: gbackoff.Config{BaseDelay: time.Minute, Multiplier: 1, MaxDelay: time.Minute},
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Errorf("failed to close connection: %v", err)
		}
	})
	conn.Connect()
	for st := conn.GetState(); st != connectivity.TransientFailure; st = conn.GetState() {
		if !conn.WaitForStateChange(ctx, st) {
			t.Fatalf("conn did not reach %s, last state %s", connectivity.TransientFailure, st)
		}
	}

	p := &pool{addr: addr}
	p.poolSize.Store(1)
	p.init()
	p.getSlots()[0].Store(&poolConn{conn: conn, addr: addr})

	pctx := &childTrackingCtx{done: make(chan struct{})}
	pc, ok := p.getHealthyConn(pctx)
	if !ok || pc == nil || pc.conn != conn {
		t.Fatalf("getHealthyConn() = (%v, %t), want the TransientFailure conn", pc, ok)
	}
	if reg, stop := pctx.registered.Load(), pctx.stopped.Load(); reg != stop {
		t.Errorf("%d child contexts registered on the caller's ctx but only %d released", reg, stop)
	}
}

// TestPool_grow_afterDisconnect verifies that a grow from a connect that passed its
// closing check before Disconnect does not republish slots into the disconnected
// pool; it used to replace the flushed slots and resize the closed pool.
func TestPool_grow_afterDisconnect(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		growTo  uint64
		size    uint64
		withPCs bool
	}{
		{name: "empty pool", size: 2, growTo: 4},
		{name: "pool with conns", size: 2, growTo: 8, withPCs: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			p := &pool{addr: "127.0.0.10:1", oldConnCloseDelay: time.Second}
			p.poolSize.Store(test.size)
			p.init()
			var conns []*ClientConn
			if test.withPCs {
				for i := range p.getSlots() {
					conn := newIdleTestConn(t, p.addr)
					conns = append(conns, conn)
					p.getSlots()[i].Store(&poolConn{conn: conn, addr: p.addr})
				}
			}
			if err := p.Disconnect(ctx); err != nil {
				t.Fatalf("Disconnect returned error: %v", err)
			}
			flushed := p.connSlots.Load()

			p.grow(test.growTo)

			if got := p.connSlots.Load(); got != flushed {
				t.Errorf("grow replaced the slots of a disconnected pool (len %d -> %d)", len(*flushed), len(*got))
			}
			if got := p.Size(); got != test.size {
				t.Errorf("Size() = %d after grow on a disconnected pool, want %d", got, test.size)
			}
			for i := range p.getSlots() {
				if pc := p.getSlots()[i].Load(); pc != nil {
					t.Errorf("slot %d must stay empty after Disconnect", i)
				}
			}
			for i, conn := range conns {
				if st := conn.GetState(); st != connectivity.Shutdown {
					t.Errorf("conn %d state = %s, want %s", i, st, connectivity.Shutdown)
				}
			}
		})
	}
}
