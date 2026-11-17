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
	"math"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/goleak"
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

// TestPool_IsIPLiteral_Classification verifies that IP literals are correctly classified
// as isIPAddr = true while hostnames (e.g. Kubernetes headless services) are classified
// as isIPAddr = false.
func TestPool_IsIPLiteral_Classification(t *testing.T) {
	t.Parallel()
	tests := []struct {
		addr     string
		wantIsIP bool
	}{
		{addr: "127.0.0.1:5001", wantIsIP: true},
		{addr: "10.0.0.1:8081", wantIsIP: true},
		{addr: "[::1]:5001", wantIsIP: true},
		{addr: "[2001:db8::1]:8080", wantIsIP: true},
		{addr: "localhost:5001", wantIsIP: false},
		{addr: "vald-agent-headless:8081", wantIsIP: false},
		{addr: "agent.vald.svc.cluster.local:8081", wantIsIP: false},
	}

	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			p, err := New(ctx,
				WithAddr(tt.addr),
				WithSize(1),
				WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
			)
			if err != nil {
				t.Fatalf("New() failed for %s: %v", tt.addr, err)
			}
			po, ok := p.(*pool)
			if !ok {
				t.Fatalf("p is not *pool: %T", p)
			}
			if got := po.IsIPConn(); got != tt.wantIsIP {
				t.Errorf("IsIPConn() for %s = %v, want %v", tt.addr, got, tt.wantIsIP)
			}
			if got := po.isIPAddr; got != tt.wantIsIP {
				t.Errorf("isIPAddr for %s = %v, want %v", tt.addr, got, tt.wantIsIP)
			}
		})
	}
}

// TestPool_RoundRobin_UniformDistribution verifies that getHealthyConn distributes
// selections uniformly across all pool slots under concurrent contention without CAS bias.
func TestPool_RoundRobin_UniformDistribution(t *testing.T) {
	t.Parallel()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer()
	go func() {
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			t.Errorf("Serve failed: %v", err)
		}
	}()
	t.Cleanup(func() {
		s.Stop()
		if err := lis.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			t.Errorf("Close failed: %v", err)
		}
	})

	addr := lis.Addr().String()
	ctx := t.Context()
	const poolSize = 4
	p, err := New(ctx,
		WithAddr(addr),
		WithSize(poolSize),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		t.Fatal(err)
	}
	p, err = p.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)

	po := p.(*pool)
	const totalCalls = 40000
	const goroutines = 20
	const callsPerGoroutine = totalCalls / goroutines

	counts := make([]atomic.Uint64, poolSize)
	eg, egctx := errgroup.New(ctx)

	for range goroutines {
		eg.Go(func() error {
			for range callsPerGoroutine {
				pc, ok := po.getHealthyConn(egctx)
				if !ok || pc == nil {
					return errors.New("getHealthyConn returned not ok")
				}
				// Identify which slot was returned
				slots := *po.connSlots.Load()
				for idx := range poolSize {
					if slots[idx].Load() == pc {
						counts[idx].Add(1)
						break
					}
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	// Verify each slot received roughly 25% of requests (tolerance +/- 5%)
	for i := range poolSize {
		got := counts[i].Load()
		expected := uint64(totalCalls / poolSize)
		diff := float64(got) - float64(expected)
		if diff < 0 {
			diff = -diff
		}
		if diff > float64(expected)*0.05 {
			t.Errorf("slot %d received %d calls, expected ~%d (diff %f > 5%%)", i, got, expected, diff)
		}
	}
}

// TestPool_getHealthyConn_concurrentChurn stress-tests
// getHealthyConn and currentIndex across 128 concurrent goroutines while
// simultaneous connection churn (IsHealthy, Reconnect, and redials) occurs.
func TestPool_getHealthyConn_concurrentChurn(t *testing.T) {
	t.Parallel()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer()
	go func() {
		_ = s.Serve(lis)
	}()
	t.Cleanup(func() {
		s.Stop()
		_ = lis.Close()
	})

	addr := lis.Addr().String()
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()

	const poolSize = 8
	p, err := New(ctx,
		WithAddr(addr),
		WithSize(poolSize),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		t.Fatal(err)
	}
	p, err = p.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}

	po, ok := p.(*pool)
	if !ok {
		t.Fatalf("expected *pool, got %T", p)
	}

	const numWorkers = 128
	const opsPerWorker = 2000
	slotCounts := make([]atomic.Uint64, poolSize)

	eg, egctx := errgroup.New(ctx)

	// Churn goroutine: repeatedly checks health and triggers reconnects/refreshes
	churnStop := make(chan struct{})
	t.Cleanup(func() { close(churnStop) })
	go func() {
		ticker := time.NewTicker(5 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-churnStop:
				return
			case <-egctx.Done():
				return
			case <-ticker.C:
				_ = po.IsHealthy(egctx)
				_, _ = po.Reconnect(egctx, false)
			}
		}
	}()

	// 128 workers hammering getHealthyConn and Get
	for range numWorkers {
		eg.Go(func() error {
			for range opsPerWorker {
				pc, ok := po.getHealthyConn(egctx)
				if !ok || pc == nil {
					// During rapid reconnect/churn, a slot may briefly be transitioning
					continue
				}
				if pc.conn == nil {
					return errors.New("getHealthyConn returned poolConn with nil ClientConn")
				}
				// Verify Get() also succeeds
				conn, gok := po.Get(egctx)
				if gok && conn == nil {
					return errors.New("po.Get() returned ok=true but nil conn")
				}

				// Find which slot holds this pc
				slots := *po.connSlots.Load()
				for idx := range poolSize {
					if idx < len(slots) && slots[idx].Load() == pc {
						slotCounts[idx].Add(1)
						break
					}
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("concurrent churn test failed: %v", err)
	}

	totalSelected := uint64(0)
	for i := range poolSize {
		cnt := slotCounts[i].Load()
		totalSelected += cnt
		t.Logf("slot %d selected %d times", i, cnt)
	}
	if totalSelected == 0 {
		t.Errorf("expected selections across slots, got 0")
	}
}

// TestPool_getHealthyConn_indexWrapAround tests that round-robin index
// distribution maintains uniformity across MaxUint64 overflow without panic or bias.
func TestPool_getHealthyConn_indexWrapAround(t *testing.T) {
	t.Parallel()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer()
	go func() {
		_ = s.Serve(lis)
	}()
	t.Cleanup(func() {
		s.Stop()
		_ = lis.Close()
	})

	addr := lis.Addr().String()
	ctx := t.Context()
	const poolSize = 4
	p, err := New(ctx,
		WithAddr(addr),
		WithSize(poolSize),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		t.Fatal(err)
	}
	p, err = p.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}

	po := p.(*pool)

	// Set currentIndex right before MaxUint64 to force uint64 wrap-around
	po.currentIndex.Store(math.MaxUint64 - 50)

	const totalCalls = 40000
	const goroutines = 20
	const callsPerWorker = totalCalls / goroutines
	counts := make([]atomic.Uint64, poolSize)

	eg, egctx := errgroup.New(ctx)
	for range goroutines {
		eg.Go(func() error {
			for range callsPerWorker {
				pc, ok := po.getHealthyConn(egctx)
				if !ok || pc == nil {
					return errors.New("getHealthyConn returned not ok during wrap-around")
				}
				slots := *po.connSlots.Load()
				for idx := range poolSize {
					if slots[idx].Load() == pc {
						counts[idx].Add(1)
						break
					}
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("wrap-around test failed: %v", err)
	}

	expectedPerSlot := uint64(totalCalls / poolSize)
	for i := range poolSize {
		got := counts[i].Load()
		diff := float64(got) - float64(expectedPerSlot)
		if diff < 0 {
			diff = -diff
		}
		tolerance := float64(expectedPerSlot) * 0.08 // 8% tolerance under high concurrency
		if diff > tolerance {
			t.Errorf("slot %d got %d selections, expected ~%d (diff %f > tolerance %f)",
				i, got, expectedPerSlot, diff, tolerance)
		}
	}
}

// TestPool_connect_flappingEndpoints tests that changing IPs
// via connect/Reconnect under concurrent load safely replaces connections
// without race conditions, nil dereferences, or dropped calls.
func TestPool_connect_flappingEndpoints(t *testing.T) {
	t.Parallel()

	// Server 1
	lis1, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s1 := grpc.NewServer()
	go func() { _ = s1.Serve(lis1) }()
	t.Cleanup(func() { s1.Stop(); _ = lis1.Close() })

	// Server 2
	lis2, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s2 := grpc.NewServer()
	go func() { _ = s2.Serve(lis2) }()
	t.Cleanup(func() { s2.Stop(); _ = lis2.Close() })

	addr1 := lis1.Addr().String()
	addr2 := lis2.Addr().String()

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	const poolSize = 4
	p, err := New(ctx,
		WithAddr(addr1),
		WithSize(poolSize),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		t.Fatal(err)
	}
	p, err = p.Connect(ctx)
	if err != nil {
		t.Fatal(err)
	}

	po := p.(*pool)

	const numWorkers = 32
	const iterations = 1000
	eg, egctx := errgroup.New(ctx)

	// Workers continuously call Get() and getHealthyConn()
	for range numWorkers {
		eg.Go(func() error {
			for range iterations {
				pc, ok := po.getHealthyConn(egctx)
				if ok && pc != nil && pc.conn == nil {
					return errors.New("getHealthyConn returned poolConn with nil ClientConn during flapping")
				}
				time.Sleep(100 * time.Microsecond)
			}
			return nil
		})
	}

	// Churner switches between addr1 and addr2
	eg.Go(func() error {
		targets := [][]string{
			{"127.0.0.1"},
			{"127.0.0.1", "127.0.0.1"},
		}
		for i := range 10 {
			_ = i
			_, cerr := po.singleTargetConnect(egctx, addr2)
			if cerr != nil {
				return cerr
			}
			time.Sleep(5 * time.Millisecond)
			_, cerr = po.connect(egctx, targets[i%len(targets)]...)
			if cerr != nil {
				return cerr
			}
			time.Sleep(5 * time.Millisecond)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		t.Fatalf("flapping endpoints test failed: %v", err)
	}
}

// Test_poolConn_Close_invalidatesStateCache verifies that
// pc.Close() immediately invalidates the state cache to connectivity.Shutdown
// so subsequent getState() calls do not return stale connectivity.Ready.
func Test_poolConn_Close_invalidatesStateCache(t *testing.T) {
	t.Parallel()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer()
	go func() { _ = s.Serve(lis) }()
	t.Cleanup(func() { s.Stop(); _ = lis.Close() })

	addr := lis.Addr().String()
	ctx := t.Context()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	pc := &poolConn{conn: conn, addr: addr}
	pc.state.Store(int32(connectivity.Ready) + 1)

	// Verify cached state is Ready
	if got := pc.getState(); got != connectivity.Ready {
		t.Fatalf("expected cached Ready, got %v", got)
	}

	// Close connection
	err = pc.Close(ctx, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Immediately check getState() - must NOT be Ready
	if got := pc.getState(); got == connectivity.Ready {
		t.Errorf("getState() returned stale Ready after Close()")
	}
	if got := pc.getState(); got != connectivity.Shutdown {
		t.Errorf("getState() = %v, want Shutdown", got)
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_poolConn_Close(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		delay time.Duration
// 	}
// 	type fields struct {
// 		conn *ClientConn
// 		addr string
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           delay:nil,
// 		       },
// 		       fields: fields {
// 		           conn:nil,
// 		           addr:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           delay:nil,
// 		           },
// 		           fields: fields {
// 		           conn:nil,
// 		           addr:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			pc := &poolConn{
// 				conn: test.fields.conn,
// 				addr: test.fields.addr,
// 			}
//
// 			err := pc.Close(test.args.ctx, test.args.delay)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		opts []Option
// 	}
// 	type want struct {
// 		wantC Conn
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotC Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotC, w.wantC) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotC, w.wantC)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotC, err := New(test.args.ctx, test.args.opts...)
// 			if err := checkFunc(test.want, gotC, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_init(t *testing.T) {
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			p.init()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_getSlots(t *testing.T) {
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want []atomic.Pointer[poolConn]
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []atomic.Pointer[poolConn]) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got []atomic.Pointer[poolConn]) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got := p.getSlots()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_grow(t *testing.T) {
// 	type args struct {
// 		newSize uint64
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           newSize:0,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           newSize:0,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			p.grow(test.args.newSize)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_load(t *testing.T) {
// 	type args struct {
// 		idx uint64
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		wantRidx uint64
// 		wantPc   *poolConn
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64, *poolConn) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRidx uint64, gotPc *poolConn) error {
// 		if !reflect.DeepEqual(gotRidx, w.wantRidx) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRidx, w.wantRidx)
// 		}
// 		if !reflect.DeepEqual(gotPc, w.wantPc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPc, w.wantPc)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           idx:0,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           idx:0,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			gotRidx, gotPc := p.load(test.args.idx)
// 			if err := checkFunc(test.want, gotRidx, gotPc); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_store(t *testing.T) {
// 	type args struct {
// 		idx uint64
// 		pc  *poolConn
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           idx:0,
// 		           pc:poolConn{},
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           idx:0,
// 		           pc:poolConn{},
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			p.store(test.args.idx, test.args.pc)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_loop(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		fn  func(ctx context.Context, idx uint64, pc *poolConn) bool
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           fn:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           fn:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			err := p.loop(test.args.ctx, test.args.fn)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_slotCount(t *testing.T) {
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got := p.slotCount()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_flush(t *testing.T) {
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			p.flush()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_refreshConn(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		eg   errgroup.Group
// 		idx  uint64
// 		pc   *poolConn
// 		addr string
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           eg:nil,
// 		           idx:0,
// 		           pc:poolConn{},
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           eg:nil,
// 		           idx:0,
// 		           pc:poolConn{},
// 		           addr:"",
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			p.refreshConn(test.args.ctx, test.args.eg, test.args.idx, test.args.pc, test.args.addr)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Connect(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want Conn
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got, err := p.Connect(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_connect(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		ips []string
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		wantC Conn
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotC Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotC, w.wantC) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotC, w.wantC)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           ips:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           ips:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			gotC, err := p.connect(test.args.ctx, test.args.ips...)
// 			if err := checkFunc(test.want, gotC, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_singleTargetConnect(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want Conn
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           addr:"",
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got, err := p.singleTargetConnect(test.args.ctx, test.args.addr)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Reconnect(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		force bool
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want Conn
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           force:false,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           force:false,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got, err := p.Reconnect(test.args.ctx, test.args.force)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Disconnect(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			err := p.Disconnect(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_dial(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		idx  uint64
// 		addr string
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want *ClientConn
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *ClientConn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got *ClientConn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           idx:0,
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           idx:0,
// 		           addr:"",
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got, err := p.dial(test.args.ctx, test.args.idx, test.args.addr)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_getHealthyConn(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		wantPc *poolConn
// 		wantOk bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *poolConn, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotPc *poolConn, gotOk bool) error {
// 		if !reflect.DeepEqual(gotPc, w.wantPc) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPc, w.wantPc)
// 		}
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			gotPc, gotOk := p.getHealthyConn(test.args.ctx)
// 			if err := checkFunc(test.want, gotPc, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Do(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		f   func(conn *ClientConn) error
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			err := p.Do(test.args.ctx, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Get(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		wantConn *ClientConn
// 		wantOk   bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *ClientConn, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotConn *ClientConn, gotOk bool) error {
// 		if !reflect.DeepEqual(gotConn, w.wantConn) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotConn, w.wantConn)
// 		}
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			gotConn, gotOk := p.Get(test.args.ctx)
// 			if err := checkFunc(test.want, gotConn, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_IsHealthy(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got := p.IsHealthy(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Len(t *testing.T) {
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got := p.Len()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_Size(t *testing.T) {
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got := p.Size()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_IsIPConn(t *testing.T) {
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got := p.IsIPConn()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_String(t *testing.T) {
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want string
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got string) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got := p.String()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_lookupIPAddr(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		want []string
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			got, err := p.lookupIPAddr(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_scanGRPCPort(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		wantPort uint16
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint16, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotPort uint16, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotPort, w.wantPort) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPort, w.wantPort)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			gotPort, err := p.scanGRPCPort(test.args.ctx)
// 			if err := checkFunc(test.want, gotPort, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestMetrics(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type want struct {
// 		want map[string]uint64
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, map[string]uint64) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got map[string]uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           in0:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           in0:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := Metrics(test.args.in0)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_pool_isHealthy(t *testing.T) {
// 	type args struct {
// 		idx  uint64
// 		conn *ClientConn
// 	}
// 	type fields struct {
// 		errGroup          errgroup.Group
// 		bo                backoff.Backoff
// 		dnsHash           atomic.Pointer[string]
// 		connSlots         atomic.Pointer[[]atomic.Pointer[poolConn]]
// 		host              string
// 		addr              string
// 		dialOpts          []DialOption
// 		oldConnCloseDelay time.Duration
// 		poolSize          atomic.Uint64
// 		currentIndex      atomic.Uint64
// 		dialTimeout       time.Duration
// 		closing           atomic.Bool
// 		port              uint16
// 		endPort           uint16
// 		startPort         uint16
// 		enableDNSLookup   bool
// 		isIPAddr          bool
// 		enableMetrics     bool
// 	}
// 	type want struct {
// 		wantState   connectivity.State
// 		wantHealthy bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, connectivity.State, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotState connectivity.State, gotHealthy bool) error {
// 		if !reflect.DeepEqual(gotState, w.wantState) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotState, w.wantState)
// 		}
// 		if !reflect.DeepEqual(gotHealthy, w.wantHealthy) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotHealthy, w.wantHealthy)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           idx:0,
// 		           conn:nil,
// 		       },
// 		       fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           idx:0,
// 		           conn:nil,
// 		           },
// 		           fields: fields {
// 		           errGroup:nil,
// 		           bo:nil,
// 		           dnsHash:nil,
// 		           connSlots:nil,
// 		           host:"",
// 		           addr:"",
// 		           dialOpts:nil,
// 		           oldConnCloseDelay:nil,
// 		           poolSize:nil,
// 		           currentIndex:nil,
// 		           dialTimeout:nil,
// 		           closing:nil,
// 		           port:0,
// 		           endPort:0,
// 		           startPort:0,
// 		           enableDNSLookup:false,
// 		           isIPAddr:false,
// 		           enableMetrics:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			p := &pool{
// 				errGroup:          test.fields.errGroup,
// 				bo:                test.fields.bo,
// 				dnsHash:           test.fields.dnsHash,
// 				connSlots:         test.fields.connSlots,
// 				host:              test.fields.host,
// 				addr:              test.fields.addr,
// 				dialOpts:          test.fields.dialOpts,
// 				oldConnCloseDelay: test.fields.oldConnCloseDelay,
// 				poolSize:          test.fields.poolSize,
// 				currentIndex:      test.fields.currentIndex,
// 				dialTimeout:       test.fields.dialTimeout,
// 				closing:           test.fields.closing,
// 				port:              test.fields.port,
// 				endPort:           test.fields.endPort,
// 				startPort:         test.fields.startPort,
// 				enableDNSLookup:   test.fields.enableDNSLookup,
// 				isIPAddr:          test.fields.isIPAddr,
// 				enableMetrics:     test.fields.enableMetrics,
// 			}
//
// 			gotState, gotHealthy := p.isHealthy(test.args.idx, test.args.conn)
// 			if err := checkFunc(test.want, gotState, gotHealthy); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

// Test_newPoolConn_watch verifies that the state published by the watcher follows
// the connection (Ready, then not Ready once the server is gone) and that the
// watcher goroutine ends when the connection is closed.
func Test_newPoolConn_watch(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer()
	go func() { _ = s.Serve(lis) }()
	addr := lis.Addr().String()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.Stop()
		t.Fatal(err)
	}
	conn.Connect()
	pc := newPoolConn(conn, addr)

	waitState := func(desc string, ok func(connectivity.State) bool) {
		t.Helper()
		deadline := time.Now().Add(10 * time.Second)
		for !ok(pc.getState()) {
			if time.Now().After(deadline) {
				t.Fatalf("%s: state is still %v", desc, pc.getState())
			}
			time.Sleep(time.Millisecond)
		}
	}
	waitState("connection becomes Ready", func(st connectivity.State) bool { return st == connectivity.Ready })

	s.Stop()
	_ = lis.Close()
	waitState("connection leaves Ready after the server stops", func(st connectivity.State) bool { return st != connectivity.Ready })

	if err := pc.Close(t.Context(), 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if got := pc.getState(); got != connectivity.Shutdown {
		t.Errorf("getState() after Close = %v, want Shutdown", got)
	}
}

// Test_pool_healthChecksDoNotStoreState checks that the connection state the
// watcher publishes is not overwritten by IsHealthy or by getHealthyConn: both
// read the connection first, and a store derived from that read could land after
// a newer state published by the watcher. (refreshConn's old store is not covered:
// it repeated newPoolConn's store within a very narrow window and cannot be observed
// from a test.) The state
// of a stable Ready connection is set to a marker the watcher never produces for
// it, and must be left alone.
func Test_pool_healthChecksDoNotStoreState(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer()
	go func() { _ = s.Serve(lis) }()
	defer s.Stop()
	addr := lis.Addr().String()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	conn.Connect()
	pc := newPoolConn(conn, addr)
	defer func() {
		if err := pc.Close(t.Context(), 100*time.Millisecond); err != nil {
			t.Error(err)
		}
	}()

	deadline := time.Now().Add(10 * time.Second)
	for pc.getState() != connectivity.Ready {
		if time.Now().After(deadline) {
			t.Fatalf("connection is not Ready: %v", pc.getState())
		}
		time.Sleep(time.Millisecond)
	}

	p := &pool{addr: addr}
	p.poolSize.Store(1)
	slots := make([]atomic.Pointer[poolConn], 1)
	slots[0].Store(pc)
	p.connSlots.Store(&slots)
	p.currentIndex.Store(0)

	tests := []struct {
		run    func(t *testing.T)
		name   string
		marker connectivity.State
	}{
		{
			name:   "IsHealthy does not store the state it read",
			marker: connectivity.TransientFailure,
			run: func(t *testing.T) {
				t.Helper()
				if !p.IsHealthy(t.Context()) {
					t.Error("IsHealthy() = false for a Ready connection")
				}
			},
		},
		{
			name:   "getHealthyConn does not store Connecting when the published state is Idle",
			marker: connectivity.Idle,
			run: func(t *testing.T) {
				t.Helper()
				if got, ok := p.getHealthyConn(t.Context()); !ok || got != pc {
					t.Errorf("getHealthyConn() = %v, %v, want the Idle connection", got, ok)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc.state.Store(int32(tt.marker) + 1) //nolint:gosec // connectivity.State is 0..4
			tt.run(t)
			if got := pc.getState(); got != tt.marker {
				t.Errorf("state = %v after the call, want it left at %v", got, tt.marker)
			}
		})
	}
}
