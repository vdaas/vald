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
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	_ "github.com/vdaas/vald/internal/net/grpc/logger"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	ocodes "go.opentelemetry.io/otel/codes"
	otrace "go.opentelemetry.io/otel/trace"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

// fixedConnPool is a pool.Conn whose Get always returns conn.
type fixedConnPool struct {
	pool.Conn
	conn *ClientConn
}

func (p fixedConnPool) Get(context.Context) (*ClientConn, bool) {
	return p.conn, p.conn != nil
}

// Test_gRPCClient_executeRPC_noPerCallOptions verifies that executeRPC does not
// pass g.copts per call: they are already installed as the dial-level default
// call options of every pooled conn, and repeating them per call made grpc
// allocate a combined slice and apply every option twice for each RPC.
func Test_gRPCClient_executeRPC_noPerCallOptions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		opts []Option
	}{
		{
			name: "with call options",
			opts: []Option{WithWaitForReady(true), WithMaxRecvMsgSize(1 << 20)},
		},
		{
			name: "without call options",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g, ok := New("test", tc.opts...).(*gRPCClient)
			if !ok {
				t.Fatal("New did not return *gRPCClient")
			}
			conn, err := ggrpc.NewClient("passthrough:///executeRPC-test",
				ggrpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if cerr := conn.Close(); cerr != nil {
					t.Error(cerr)
				}
			})
			got := -1
			_, _, err = g.executeRPC(t.Context(), fixedConnPool{conn: conn}, "addr",
				func(_ context.Context, _ *ClientConn, copts ...CallOption) (any, error) {
					got = len(copts)
					return struct{}{}, nil
				})
			if err != nil {
				t.Fatalf("executeRPC() err = %v", err)
			}
			if got != 0 {
				t.Errorf("f received %d per-call options, want 0 (defaults come from the conn)", got)
			}
		})
	}
}

// TestGRPCClient_mergedDialOptions verifies that per-call dial options never
// leak into the shared g.dopts backing array, which would let concurrent
// Connect calls overwrite each other's options.
func TestGRPCClient_mergedDialOptions(t *testing.T) {
	t.Parallel()
	base := ggrpc.WithUserAgent("base")
	first := ggrpc.WithUserAgent("first")
	second := ggrpc.WithUserAgent("second")

	tests := []struct {
		name string
		// spare is the spare capacity of g.dopts.
		spare int
	}{
		{name: "g.dopts without spare capacity", spare: 0},
		{name: "g.dopts with spare capacity", spare: 4},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			dopts := make([]DialOption, 1, 1+test.spare)
			dopts[0] = base
			g := &gRPCClient{dopts: dopts}

			got1 := g.mergedDialOptions([]DialOption{first})
			got2 := g.mergedDialOptions([]DialOption{second})

			if len(got1) != 2 || got1[0] != base || got1[1] != first {
				t.Errorf("first call result was modified or wrong: %v", got1)
			}
			if len(got2) != 2 || got2[0] != base || got2[1] != second {
				t.Errorf("second call result is wrong: %v", got2)
			}
			if len(g.dopts) != 1 {
				t.Errorf("g.dopts length changed: %d", len(g.dopts))
			}
			if full := g.dopts[:cap(g.dopts)]; test.spare > 0 && (full[1] != nil || full[2] != nil) {
				t.Errorf("g.dopts spare capacity was written: %v", full)
			}
		})
	}
}

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

// evictRecordingDialer is a net.Dialer fake that records EvictAddr calls. The
// embedded Dialer is nil: these tests never dial.
type evictRecordingDialer struct {
	net.Dialer
	evicted []string
	mu      sync.Mutex
}

func (d *evictRecordingDialer) EvictAddr(addr string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evicted = append(d.evicted, addr)
}

// Test_gRPCClient_Disconnect_evictsAddrState verifies that Disconnect drops the
// circuit breaker metrics and the dialer state of the removed address, and of
// no other address. They used to grow with every address ever connected.
func Test_gRPCClient_Disconnect_evictsAddrState(t *testing.T) {
	t.Parallel()
	const (
		method = "Test_gRPCClient_Disconnect_evictsAddrState"
		addr   = "192.0.2.166:8081"
		other  = "192.0.2.167:8081"
	)
	tests := []struct {
		name        string
		wantEvicted []string
		wantGone    []string
		wantKept    []string
		stored      bool
	}{
		{
			name:        "connected address",
			stored:      true,
			wantEvicted: []string{addr},
			wantGone:    []string{method + "/" + addr, addr},
			wantKept:    []string{method + "/" + other},
		},
		{
			name:     "unknown address",
			stored:   false,
			wantKept: []string{method + "/" + addr, addr, method + "/" + other},
		},
	}
	// Subtests run in sequence: the breaker metrics are package-global.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			cb, err := circuitbreaker.New()
			if err != nil {
				t.Fatal(err)
			}
			der := new(evictRecordingDialer)
			g, ok := New("evict-test", WithCircuitBreaker(cb), WithDialer("tcp", der)).(*gRPCClient)
			if !ok {
				t.Fatal("unexpected client type")
			}
			for _, key := range []string{method + "/" + addr, addr, method + "/" + other} {
				if _, derr := cb.Do(ctx, key, func(context.Context) (any, error) {
					return true, nil
				}); derr != nil {
					t.Fatal(derr)
				}
			}
			if test.stored {
				g.conns.Store(addr, &hookConn{})
				g.clientCount = 1
			}

			err = g.Disconnect(ctx, addr)
			if (err == nil) != test.stored {
				t.Fatalf("Disconnect returned error %v for stored=%v", err, test.stored)
			}

			ms := circuitbreaker.Metrics(ctx)
			for _, key := range test.wantGone {
				if _, ok := ms[key]; ok {
					t.Errorf("breaker metrics for %q remain after Disconnect", key)
				}
			}
			for _, key := range test.wantKept {
				if _, ok := ms[key]; !ok {
					t.Errorf("breaker metrics for %q were removed", key)
				}
			}
			der.mu.Lock()
			evicted := slices.Clone(der.evicted)
			der.mu.Unlock()
			// Disconnect evicts before and after closing the pool, so the same
			// address may be evicted more than once.
			slices.Sort(evicted)
			evicted = slices.Compact(evicted)
			if !slices.Equal(evicted, test.wantEvicted) {
				t.Errorf("dialer evicted %v, want %v", evicted, test.wantEvicted)
			}
		})
	}
}

// TestGRPCClient_StartConnectionMonitor_HostnameNotDisconnected verifies that
// hostnames are not removed by the bulk disconnection monitor, preserving retries.
func TestGRPCClient_StartConnectionMonitor_HostnameNotDisconnected(t *testing.T) {
	t.Parallel()
	c := New("hostname-disconnect-test",
		WithAddrs("vald-agent-headless:8081", "127.0.0.1:8081"),
	)
	g, ok := c.(*gRPCClient)
	if !ok {
		t.Fatalf("unexpected client type: %T", c)
	}
	if g == nil {
		t.Fatal("client is nil")
	}
	// Simulate connection monitor evaluation for both targets
	disconnectTargets := []string{"vald-agent-headless:8081", "127.0.0.1:8081"}
	for _, addr := range disconnectTargets {
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			t.Fatalf("net.SplitHostPort failed for %s: %v", addr, err)
		}
		isDirectIP := net.IsIPLiteral(host)
		if addr == "vald-agent-headless:8081" && isDirectIP {
			t.Errorf("hostname %s misclassified as direct IP connection", addr)
		}
		if addr == "127.0.0.1:8081" && !isDirectIP {
			t.Errorf("IP %s not classified as direct IP connection", addr)
		}
	}
}

// TestGRPCClient_StartConnectionMonitor_initialConnectionFailure verifies that a failed
// initial connection does not leave the monitor marked as running, so that
// following calls retry and report the error instead of returning (nil, nil).
func TestGRPCClient_StartConnectionMonitor_initialConnectionFailure(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		calls int
	}{
		{name: "first call fails", calls: 1},
		{name: "every retry reports the failure again", calls: 3},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			// an empty address is skipped on the initial connection, so no connection is
			// established (no dial, no port scan) and the monitor start fails.
			c := New("monitor-failure-test", WithAddrs(""))
			g, ok := c.(*gRPCClient)
			if !ok {
				t.Fatalf("unexpected client type: %T", c)
			}
			for i := range test.calls {
				ech, err := g.StartConnectionMonitor(ctx)
				if err == nil {
					t.Fatalf("call %d: StartConnectionMonitor must return an error, got ech: %v", i, ech)
				}
				if ech != nil {
					t.Errorf("call %d: error channel must be nil on error", i)
				}
				if g.monitorRunning.Load() {
					t.Errorf("call %d: monitorRunning must be reset after the failure", i)
				}
			}
		})
	}
}

// TestGRPCClient_StartConnectionMonitor_alreadyRunning verifies that a call while the monitor
// is running does not start another one and returns the existing channel.
func TestGRPCClient_StartConnectionMonitor_alreadyRunning(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	g := &gRPCClient{}
	ech := make(chan error)
	g.ech = ech
	g.monitorRunning.Store(true)

	got, err := g.StartConnectionMonitor(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got != (<-chan error)(ech) {
		t.Errorf("got channel %v, want %v", got, ech)
	}
	if !g.monitorRunning.Load() {
		t.Error("monitorRunning must stay true")
	}
}

// TestGRPCClient_StartConnectionMonitor_concurrentFailure verifies that concurrent callers
// never observe (nil, nil) while the initial connection of the starting caller fails.
func TestGRPCClient_StartConnectionMonitor_concurrentFailure(t *testing.T) {
	t.Parallel()
	const callers = 8
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	c := New("monitor-concurrent-failure-test", WithAddrs(""))
	g, ok := c.(*gRPCClient)
	if !ok {
		t.Fatalf("unexpected client type: %T", c)
	}
	type result struct {
		ech <-chan error
		err error
	}
	results := make(chan result, callers)
	for range callers {
		go func() {
			ech, err := g.StartConnectionMonitor(ctx)
			results <- result{ech: ech, err: err}
		}()
	}
	for range callers {
		select {
		case r := <-results:
			if r.err == nil {
				t.Errorf("StartConnectionMonitor must return an error, got ech: %v", r.ech)
			}
			if r.ech != nil {
				t.Error("error channel must be nil on error")
			}
		case <-ctx.Done():
			t.Fatal("timed out waiting for StartConnectionMonitor callers")
		}
	}
	if g.monitorRunning.Load() {
		t.Error("monitorRunning must be reset after the failure")
	}
}

// Test_gRPCClient_executeRPC_noResourceExhaustedRetry verifies that only
// ResourceExhausted becomes non-retryable under WithoutResourceExhaustedRetry,
// and that every code keeps its default classification without it.
func Test_gRPCClient_executeRPC_noResourceExhaustedRetry(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		code          codes.Code
		marked        bool
		wantRetryable bool
	}{
		{name: "ResourceExhausted is retried by default", code: codes.ResourceExhausted, wantRetryable: true},
		{name: "ResourceExhausted is not retried when marked", code: codes.ResourceExhausted, marked: true},
		{name: "Unavailable is still retried when marked", code: codes.Unavailable, marked: true, wantRetryable: true},
		{name: "Internal is still retried when marked", code: codes.Internal, marked: true, wantRetryable: true},
		{name: "InvalidArgument is never retried", code: codes.InvalidArgument, marked: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g, ok := New("test").(*gRPCClient)
			if !ok {
				t.Fatal("New did not return *gRPCClient")
			}
			conn, err := ggrpc.NewClient("passthrough:///re-test",
				ggrpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if cerr := conn.Close(); cerr != nil {
					t.Error(cerr)
				}
			})
			rpcErr := status.Error(tc.code, "rpc failed")
			ctx := WithQuietErrorLog(t.Context())
			if tc.marked {
				ctx = WithoutResourceExhaustedRetry(ctx)
			}
			_, retryable, err := g.executeRPC(ctx, healthyConnPool{fixedConnPool{conn: conn}}, "10.0.0.1:8081",
				func(context.Context, *ClientConn, ...CallOption) (any, error) { return nil, rpcErr })
			if !errors.Is(err, rpcErr) {
				t.Fatalf("executeRPC() err = %v, want %v", err, rpcErr)
			}
			if retryable != tc.wantRetryable {
				t.Errorf("executeRPC() retryable = %v, want %v", retryable, tc.wantRetryable)
			}
		})
	}
}

// recordingSpan records the calls recordSpanError makes. The embedded nil Span
// panics on any other method, so unexpected calls fail the test.
type recordingSpan struct {
	trace.Span
	desc   string
	errs   []error
	attrs  []attribute.KeyValue
	status int
	code   ocodes.Code
}

func (s *recordingSpan) RecordError(err error, _ ...otrace.EventOption) {
	s.errs = append(s.errs, err)
}

func (s *recordingSpan) SetAttributes(kv ...attribute.KeyValue) {
	s.attrs = append(s.attrs, kv...)
}

func (s *recordingSpan) SetStatus(code ocodes.Code, desc string) {
	s.status++
	s.code, s.desc = code, desc
}

func TestRecordSpanError(t *testing.T) {
	t.Parallel()
	grpcErr := status.Error(codes.Unavailable, "unavailable")
	plainErr := errors.New("plain")
	tests := []struct {
		err       error
		name      string
		wantAttrs []attribute.KeyValue
		wantCalls bool
	}{
		{
			name: "nil error records nothing",
		},
		{
			name:      "gRPC status error records error, status attributes and status",
			err:       grpcErr,
			wantCalls: true,
			wantAttrs: trace.FromGRPCStatus(codes.Unavailable, grpcErr.Error()),
		},
		{
			name:      "non-status error records error and status without attributes",
			err:       plainErr,
			wantCalls: true,
		},
	}
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			span := new(recordingSpan)
			recordSpanError(span, test.err)
			if !test.wantCalls {
				if len(span.errs) != 0 || len(span.attrs) != 0 || span.status != 0 {
					tt.Errorf("span was modified: %+v", span)
				}
				return
			}
			if len(span.errs) != 1 || !errors.Is(span.errs[0], test.err) {
				tt.Errorf("RecordError calls = %v, want [%v]", span.errs, test.err)
			}
			if len(span.attrs) != len(test.wantAttrs) {
				tt.Fatalf("attributes = %v, want %v", span.attrs, test.wantAttrs)
			}
			for i := range span.attrs {
				if span.attrs[i] != test.wantAttrs[i] {
					tt.Errorf("attribute[%d] = %v, want %v", i, span.attrs[i], test.wantAttrs[i])
				}
			}
			if span.status != 1 || span.code != trace.StatusError || span.desc != test.err.Error() {
				tt.Errorf("SetStatus calls = %d (%v, %q), want 1 (%v, %q)",
					span.status, span.code, span.desc, trace.StatusError, test.err.Error())
			}
		})
	}
	// a nil span must not panic.
	recordSpanError(nil, plainErr)
}

// TestStartAddrSpan_TracingDisabled pins that, with tracing off, startAddrSpan
// returns the given ctx and a nil span without building the span name.
// Not parallel: AllocsPerRun changes GOMAXPROCS.
func TestStartAddrSpan_TracingDisabled(t *testing.T) {
	if trace.Enabled() {
		t.Skip("tracing is enabled in this process")
	}
	type ctxKey struct{}
	ctx := context.WithValue(t.Context(), ctxKey{}, "v")
	addr := string([]byte("10.0.0.1:8081")) // not a constant, like a real target
	gotCtx, span := startAddrSpan(ctx, apiName+"/Client.do/", addr)
	if gotCtx != ctx {
		t.Error("startAddrSpan returned a different ctx")
	}
	if span != nil {
		t.Errorf("startAddrSpan returned span %v, want nil", span)
	}
	if allocs := testing.AllocsPerRun(100, func() {
		_, span = startAddrSpan(ctx, apiName+"/Client.do/", addr)
	}); allocs != 0 {
		t.Errorf("startAddrSpan allocated %v times per call, want 0", allocs)
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		name string
// 		opts []Option
// 	}
// 	type want struct {
// 		wantC Client
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Client) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotC Client) error {
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
// 		           name:"",
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
// 		           name:"",
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
// 			gotC := New(test.args.name, test.args.opts...)
// 			if err := checkFunc(test.want, gotC); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_StartConnectionMonitor(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
// 	}
// 	type want struct {
// 		want <-chan error
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error, err error) error {
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			got, err := g.StartConnectionMonitor(test.args.ctx)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_Range(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		f   func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			err := g.Range(test.args.ctx, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_RangeConcurrent(t *testing.T) {
// 	type args struct {
// 		ctx         context.Context
// 		concurrency int
// 		f           func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
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
// 		           concurrency:0,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           concurrency:0,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			err := g.RangeConcurrent(test.args.ctx, test.args.concurrency, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_OrderedRange(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		orders []string
// 		f      func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
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
// 		           orders:nil,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           orders:nil,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			err := g.OrderedRange(test.args.ctx, test.args.orders, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_OrderedRangeConcurrent(t *testing.T) {
// 	type args struct {
// 		ctx         context.Context
// 		orders      []string
// 		concurrency int
// 		f           func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
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
// 		           orders:nil,
// 		           concurrency:0,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           orders:nil,
// 		           concurrency:0,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			err := g.OrderedRangeConcurrent(test.args.ctx, test.args.orders, test.args.concurrency, test.args.f)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestRoundRobin(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		c   Client
// 		f   func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (R, error)
// 	}
// 	type want struct {
// 		wantData R
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, R, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotData R, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotData, w.wantData) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
// 		           c:nil,
// 		           f:nil,
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
// 		           c:nil,
// 		           f:nil,
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
// 			gotData, err := RoundRobin(test.args.ctx, test.args.c, test.args.f)
// 			if err := checkFunc(test.want, gotData, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_RoundRobin(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		f   func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error)
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
// 	}
// 	type want struct {
// 		wantData any
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, any, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotData any, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotData, w.wantData) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			gotData, err := g.RoundRobin(test.args.ctx, test.args.f)
// 			if err := checkFunc(test.want, gotData, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_Do(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 		f    func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error)
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
// 	}
// 	type want struct {
// 		wantData any
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, any, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotData any, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotData, w.wantData) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			gotData, err := g.Do(test.args.ctx, test.args.addr, test.args.f)
// 			if err := checkFunc(test.want, gotData, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_executeRPC(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		p    pool.Conn
// 		addr string
// 		f    func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error)
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
// 	}
// 	type want struct {
// 		wantRes       any
// 		wantRetryable bool
// 		err           error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, any, bool, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes any, gotRetryable bool, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		if !reflect.DeepEqual(gotRetryable, w.wantRetryable) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRetryable, w.wantRetryable)
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
// 		           p:nil,
// 		           addr:"",
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           p:nil,
// 		           addr:"",
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			gotRes, gotRetryable, err := g.executeRPC(test.args.ctx, test.args.p, test.args.addr, test.args.f)
// 			if err := checkFunc(test.want, gotRes, gotRetryable, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_do(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		p    pool.Conn
// 		addr string
// 		f    func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error)
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
// 	}
// 	type want struct {
// 		wantData any
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, any, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotData any, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotData, w.wantData) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotData, w.wantData)
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
// 		           p:nil,
// 		           addr:"",
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           p:nil,
// 		           addr:"",
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			gotData, err := g.do(test.args.ctx, test.args.p, test.args.addr, test.args.f)
// 			if err := checkFunc(test.want, gotData, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_GetDialOption(t *testing.T) {
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
// 	}
// 	type want struct {
// 		want []DialOption
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []DialOption) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got []DialOption) error {
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			got := g.GetDialOption()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_GetCallOption(t *testing.T) {
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
// 	}
// 	type want struct {
// 		want []CallOption
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []CallOption) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got []CallOption) error {
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			got := g.GetCallOption()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_GetBackoff(t *testing.T) {
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
// 	}
// 	type want struct {
// 		want backoff.Backoff
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, backoff.Backoff) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got backoff.Backoff) error {
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			got := g.GetBackoff()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_SetDisableResolveDNSAddr(t *testing.T) {
// 	type args struct {
// 		addr     string
// 		disabled bool
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
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
// 		           addr:"",
// 		           disabled:false,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           addr:"",
// 		           disabled:false,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			g.SetDisableResolveDNSAddr(test.args.addr, test.args.disabled)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_Connect(t *testing.T) {
// 	type args struct {
// 		ctx   context.Context
// 		addr  string
// 		dopts []DialOption
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
// 	}
// 	type want struct {
// 		wantConn pool.Conn
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, pool.Conn, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotConn pool.Conn, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotConn, w.wantConn) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotConn, w.wantConn)
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
// 		           dopts:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           dopts:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			gotConn, err := g.Connect(test.args.ctx, test.args.addr, test.args.dopts...)
// 			if err := checkFunc(test.want, gotConn, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_IsConnected(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
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
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			got := g.IsConnected(test.args.ctx, test.args.addr)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_Disconnect(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		addr string
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
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
// 		           addr:"",
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			err := g.Disconnect(test.args.ctx, test.args.addr)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_ConnectedAddrs(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
// 	}
// 	type want struct {
// 		wantAddrs []string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotAddrs []string) error {
// 		if !reflect.DeepEqual(gotAddrs, w.wantAddrs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotAddrs, w.wantAddrs)
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
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			gotAddrs := g.ConnectedAddrs(test.args.in0)
// 			if err := checkFunc(test.want, gotAddrs); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_Close(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			err := g.Close(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gRPCClient_rangeConns(t *testing.T) {
// 	type args struct {
// 		action string
// 		fn     func(addr string, p pool.Conn) bool
// 	}
// 	type fields struct {
// 		dialer                 net.Dialer
// 		eg                     errgroup.Group
// 		group                  singleflight.Group[pool.Conn]
// 		cb                     circuitbreaker.CircuitBreaker
// 		bo                     backoff.Backoff
// 		addrs                  map[string]struct{}
// 		stopMonitor            context.CancelFunc
// 		ech                    <-chan error
// 		crl                    sync.Map[string, bool]
// 		disableResolveDNSAddrs sync.Map[string, bool]
// 		conns                  sync.Map[string, pool.Conn]
// 		name                   string
// 		roccd                  string
// 		dopts                  []DialOption
// 		copts                  []CallOption
// 		gbo                    gbackoff.Config
// 		hcDur                  time.Duration
// 		mcd                    time.Duration
// 		prDur                  time.Duration
// 		clientCount            uint64
// 		poolSize               uint64
// 		monitorRunning         atomic.Bool
// 		resolveDNS             bool
// 		enablePoolRebalance    bool
// 		enablePoolMetrics      bool
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
// 		           action:"",
// 		           fn:nil,
// 		       },
// 		       fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 		           action:"",
// 		           fn:nil,
// 		           },
// 		           fields: fields {
// 		           dialer:nil,
// 		           eg:nil,
// 		           group:nil,
// 		           cb:nil,
// 		           bo:nil,
// 		           addrs:nil,
// 		           stopMonitor:nil,
// 		           ech:nil,
// 		           crl:nil,
// 		           disableResolveDNSAddrs:nil,
// 		           conns:nil,
// 		           name:"",
// 		           roccd:"",
// 		           dopts:nil,
// 		           copts:nil,
// 		           gbo:nil,
// 		           hcDur:nil,
// 		           mcd:nil,
// 		           prDur:nil,
// 		           clientCount:0,
// 		           poolSize:0,
// 		           monitorRunning:nil,
// 		           resolveDNS:false,
// 		           enablePoolRebalance:false,
// 		           enablePoolMetrics:false,
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
// 			g := &gRPCClient{
// 				dialer:                 test.fields.dialer,
// 				eg:                     test.fields.eg,
// 				group:                  test.fields.group,
// 				cb:                     test.fields.cb,
// 				bo:                     test.fields.bo,
// 				addrs:                  test.fields.addrs,
// 				stopMonitor:            test.fields.stopMonitor,
// 				ech:                    test.fields.ech,
// 				crl:                    test.fields.crl,
// 				disableResolveDNSAddrs: test.fields.disableResolveDNSAddrs,
// 				conns:                  test.fields.conns,
// 				name:                   test.fields.name,
// 				roccd:                  test.fields.roccd,
// 				dopts:                  test.fields.dopts,
// 				copts:                  test.fields.copts,
// 				gbo:                    test.fields.gbo,
// 				hcDur:                  test.fields.hcDur,
// 				mcd:                    test.fields.mcd,
// 				prDur:                  test.fields.prDur,
// 				clientCount:            test.fields.clientCount,
// 				poolSize:               test.fields.poolSize,
// 				monitorRunning:         test.fields.monitorRunning,
// 				resolveDNS:             test.fields.resolveDNS,
// 				enablePoolRebalance:    test.fields.enablePoolRebalance,
// 				enablePoolMetrics:      test.fields.enablePoolMetrics,
// 			}
//
// 			err := g.rangeConns(test.args.action, test.args.fn)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
