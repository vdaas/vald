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
	"fmt"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type fanOutConnPool struct {
	conn    *ClientConn
	healthy bool
	isIP    bool
}

func (p fanOutConnPool) Connect(context.Context) (pool.Conn, error)      { return p, nil }
func (fanOutConnPool) Disconnect(context.Context) error                  { return nil }
func (fanOutConnPool) Do(context.Context, func(*ClientConn) error) error { return nil }
func (p fanOutConnPool) Get(context.Context) (*ClientConn, bool)         { return p.conn, p.conn != nil }
func (p fanOutConnPool) IsHealthy(context.Context) bool                  { return p.healthy }
func (p fanOutConnPool) IsIPConn() bool                                  { return p.isIP }
func (fanOutConnPool) Len() uint64                                       { return 1 }
func (fanOutConnPool) Size() uint64                                      { return 1 }
func (p fanOutConnPool) Reconnect(context.Context, bool) (pool.Conn, error) {
	return p, nil
}
func (fanOutConnPool) String() string { return "fanOutConnPool" }

type fanOutRecordingBackoff struct {
	calls atomic.Int64
}

func (b *fanOutRecordingBackoff) Do(
	ctx context.Context, fn func(context.Context) (any, bool, error),
) (any, error) {
	b.calls.Add(1)
	res, _, err := fn(ctx)
	return res, err
}

func (*fanOutRecordingBackoff) Close() {}

func newFanOutTestClientConn(tb testing.TB) *ClientConn {
	tb.Helper()
	conn, err := ggrpc.NewClient("passthrough:///fanout-test",
		ggrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		tb.Fatalf("NewClient() error = %v", err)
	}
	tb.Cleanup(func() {
		if cerr := conn.Close(); cerr != nil {
			tb.Error(cerr)
		}
	})
	return conn
}

func newFanOutTestClient(tb testing.TB, name string, opts ...Option) *gRPCClient {
	tb.Helper()
	g, ok := New(name, opts...).(*gRPCClient)
	if !ok {
		tb.Fatal("New did not return *gRPCClient")
	}
	tb.Cleanup(func() {
		if err := g.Close(context.Background()); err != nil {
			tb.Error(err)
		}
	})
	return g
}

func newFanOutBenchClient(b *testing.B, n int) (*gRPCClient, []string) {
	b.Helper()
	conn := newFanOutTestClientConn(b)
	g := newFanOutTestClient(b, "fanout-bench", WithErrGroup(errgroup.Get()))
	addrs := make([]string, n)
	for i := range addrs {
		addrs[i] = fmt.Sprintf("10.0.%d.%d:8081", i/250, i%250+1)
		g.storeConn(addrs[i], fanOutConnPool{conn: conn, healthy: true})
	}
	return g, addrs
}

func Test_gRPCClient_FanOut(t *testing.T) {
	t.Parallel()

	for _, n := range []int{1, 2, 5, 16} {
		t.Run("RangeConcurrent_n="+strconv.Itoa(n), func(t *testing.T) {
			t.Parallel()
			conn := newFanOutTestClientConn(t)
			g := newFanOutTestClient(t, "fanout-range", WithErrGroup(errgroup.Get()))

			want := make(map[string]int, n)
			for i := range n {
				addr := "10.0.0." + strconv.Itoa(i+1) + ":8081"
				want[addr] = 1
				g.storeConn(addr, fanOutConnPool{conn: conn, healthy: true})
			}

			var mu sync.Mutex
			got := make(map[string]int, n)
			err := g.RangeConcurrent(t.Context(), 0, func(_ context.Context, addr string, c *ClientConn, _ ...CallOption) error {
				if c == nil {
					t.Errorf("nil conn for %s", addr)
				}
				mu.Lock()
				got[addr]++
				mu.Unlock()
				return nil
			})
			if err != nil {
				t.Fatalf("RangeConcurrent() error = %v", err)
			}
			for addr, c := range want {
				if got[addr] != c {
					t.Errorf("addr %s count = %d, want %d", addr, got[addr], c)
				}
			}
		})

		t.Run("OrderedRangeConcurrent_n="+strconv.Itoa(n), func(t *testing.T) {
			t.Parallel()
			conn := newFanOutTestClientConn(t)
			g := newFanOutTestClient(t, "fanout-ordered", WithErrGroup(errgroup.Get()))

			orders := make([]string, 0, n)
			for i := range n {
				addr := "10.0.1." + strconv.Itoa(i+1) + ":8081"
				orders = append(orders, addr)
				g.storeConn(addr, fanOutConnPool{conn: conn, healthy: true})
			}
			missing := "10.0.1.254:8081"
			allOrders := append(slices.Clone(orders), missing)

			var mu sync.Mutex
			got := make(map[string]int, n)
			err := g.OrderedRangeConcurrent(t.Context(), allOrders, 0, func(_ context.Context, addr string, _ *ClientConn, _ ...CallOption) error {
				mu.Lock()
				got[addr]++
				mu.Unlock()
				return nil
			})
			if err != nil {
				t.Fatalf("OrderedRangeConcurrent() error = %v", err)
			}
			for _, addr := range orders {
				if got[addr] != 1 {
					t.Errorf("addr %s count = %d, want 1", addr, got[addr])
				}
			}
			if _, ok := g.crl.Load(missing); !ok {
				t.Errorf("missing addr %s was not queued for reconnect", missing)
			}
		})
	}

	t.Run("bounded concurrency caps in-flight callbacks", func(t *testing.T) {
		t.Parallel()
		conn := newFanOutTestClientConn(t)
		g := newFanOutTestClient(t, "fanout-bounded", WithErrGroup(errgroup.Get()))

		const (
			targets = 8
			limit   = 2
		)
		for i := range targets {
			addr := "10.0.2." + strconv.Itoa(i+1) + ":8081"
			g.storeConn(addr, fanOutConnPool{conn: conn, healthy: true})
		}

		var cur, peak, calls atomic.Int64
		err := g.RangeConcurrent(t.Context(), limit, func(context.Context, string, *ClientConn, ...CallOption) error {
			n := cur.Add(1)
			for {
				prev := peak.Load()
				if n <= prev || peak.CompareAndSwap(prev, n) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			cur.Add(-1)
			calls.Add(1)
			return nil
		})
		if err != nil {
			t.Fatalf("RangeConcurrent(limit=%d) error = %v", limit, err)
		}
		if peak.Load() > limit {
			t.Errorf("peak concurrency = %d, want <= %d", peak.Load(), limit)
		}
		if calls.Load() != targets {
			t.Errorf("calls = %d, want %d", calls.Load(), targets)
		}
	})

	t.Run("panic recovery keeps sibling targets running", func(t *testing.T) {
		t.Parallel()
		conn := newFanOutTestClientConn(t)
		g := newFanOutTestClient(t, "fanout-panic", WithErrGroup(errgroup.Get()))

		const targets = 5
		for i := range targets {
			g.storeConn("10.0.3."+strconv.Itoa(i+1)+":8081", fanOutConnPool{conn: conn, healthy: true})
		}
		var calls atomic.Int64
		err := g.RangeConcurrent(t.Context(), 0, func(_ context.Context, addr string, _ *ClientConn, _ ...CallOption) error {
			calls.Add(1)
			if addr == "10.0.3.1:8081" {
				panic("simulated fan-out panic")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("RangeConcurrent() error = %v", err)
		}
		if calls.Load() != targets {
			t.Errorf("calls = %d, want %d", calls.Load(), targets)
		}
	})

	t.Run("custom limited canceled errgroup returns admission error", func(t *testing.T) {
		t.Parallel()
		conn := newFanOutTestClientConn(t)
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		eg, _ := errgroup.New(ctx)
		eg.SetLimit(0)

		g := newFanOutTestClient(t, "fanout-custom-rejected", WithErrGroup(eg))
		g.storeConn("10.0.4.1:8081", fanOutConnPool{conn: conn, healthy: true})
		g.storeConn("10.0.4.2:8081", fanOutConnPool{conn: conn, healthy: true})

		err := g.RangeConcurrent(t.Context(), 0, func(context.Context, string, *ClientConn, ...CallOption) error {
			return nil
		})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("RangeConcurrent() error = %v, want wrapping context.Canceled", err)
		}
	})
}

func Test_gRPCClient_cachedBreakerName(t *testing.T) {
	t.Parallel()
	const addr = "agent:8081"
	tests := []struct {
		name   string
		method string
		want   string
	}{
		{
			name:   "empty method returns addr",
			method: "",
			want:   addr,
		},
		{
			name:   "method equals addr",
			method: addr,
			want:   addr,
		},
		{
			name:   "method already suffixed by slash addr",
			method: "BroadCast/" + addr,
			want:   "BroadCast/" + addr,
		},
		{
			name:   "partial hostname suffix is not mistaken for addr",
			method: "other" + addr,
			want:   "other" + addr + "/" + addr,
		},
	}
	g := newFanOutTestClient(t, "breaker-name")
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			if tc.method != "" {
				ctx = WrapGRPCMethod(ctx, tc.method)
			}
			if got := g.cachedBreakerName(ctx, addr); got != tc.want {
				t.Errorf("cachedBreakerName() = %q, want %q", got, tc.want)
			}
		})
	}
}

func Test_gRPCClient_breakerName_and_policy(t *testing.T) {
	t.Parallel()

	t.Run("Disconnect evicts suffixed partial hostname breaker key", func(t *testing.T) {
		const (
			addr        = "partial-agent:8081"
			otherMethod = "otherpartial-agent:8081"
		)
		ctx := t.Context()
		cb, err := circuitbreaker.New()
		if err != nil {
			t.Fatal(err)
		}
		g := newFanOutTestClient(t, "evict-segment", WithCircuitBreaker(cb))

		p := fanOutConnPool{conn: newFanOutTestClientConn(t), healthy: true}
		g.storeConn(addr, p)
		g.clientCount = 1

		err = g.doVoid(WrapGRPCMethod(ctx, otherMethod), p, addr, func(context.Context, string, *ClientConn, ...CallOption) error {
			return nil
		})
		if err != nil {
			t.Fatalf("doVoid() error = %v", err)
		}
		if _, err := cb.Do(ctx, otherMethod, func(context.Context) (any, error) { return true, nil }); err != nil {
			t.Fatal(err)
		}

		wantEvicted := otherMethod + "/" + addr
		if err := g.Disconnect(ctx, addr); err != nil {
			t.Fatalf("Disconnect() error = %v", err)
		}
		ms := circuitbreaker.Metrics(ctx)
		if _, exists := ms[wantEvicted]; exists {
			t.Errorf("breaker metrics for %q were not evicted", wantEvicted)
		}
		if _, exists := ms[otherMethod]; !exists {
			t.Errorf("unrelated breaker metrics for %q were erroneously evicted", otherMethod)
		}
	})

	t.Run("do and doVoid share backoff and breaker policy", func(t *testing.T) {
		t.Parallel()
		conn := newFanOutTestClientConn(t)
		bo := &fanOutRecordingBackoff{}
		cb := new(recordingBreaker)
		g := newFanOutTestClient(t, "policy-shared", WithBackoff(bo), WithCircuitBreaker(cb))

		const addr = "10.0.5.1:8081"
		p := fanOutConnPool{conn: conn, healthy: true}
		g.storeConn(addr, p)
		ctx := WrapGRPCMethod(t.Context(), "Search")

		got, err := g.do(ctx, p, addr, func(context.Context, *ClientConn, ...CallOption) (any, error) {
			return "ok", nil
		})
		if err != nil || got != "ok" {
			t.Fatalf("do() = (%v, %v), want (ok, nil)", got, err)
		}

		err = g.doVoid(ctx, p, addr, func(context.Context, string, *ClientConn, ...CallOption) error {
			return status.Error(codes.InvalidArgument, "bad")
		})
		if err == nil {
			t.Fatal("doVoid() error = nil, want non-nil")
		}
		if bo.calls.Load() != 2 {
			t.Errorf("backoff calls = %d, want 2", bo.calls.Load())
		}
		if !slices.Equal(cb.keys, []string{"Search/" + addr, "Search/" + addr}) {
			t.Errorf("breaker keys = %v, want [Search/%s Search/%s]", cb.keys, addr, addr)
		}
		if !slices.Equal(cb.ignorable, []bool{false, true}) {
			t.Errorf("breaker ignorable = %v, want [false true]", cb.ignorable)
		}
	})

	t.Run("cachedBreakerName caches and Disconnect evicts only target addr", func(t *testing.T) {
		t.Parallel()
		const (
			method = "Search"
			addr1  = "10.0.6.1:8081"
			addr2  = "10.0.6.2:8081"
		)
		g := newFanOutTestClient(t, "breaker-cache")
		ctx := WithGRPCMethod(t.Context(), method)

		got1 := g.cachedBreakerName(ctx, addr1)
		got2 := g.cachedBreakerName(ctx, addr2)
		if got1 != method+"/"+addr1 || got2 != method+"/"+addr2 {
			t.Fatalf("cachedBreakerName = (%q, %q), want (%s/%s, %s/%s)", got1, got2, method, addr1, method, addr2)
		}
		if _, ok := g.breakerKeys.Load(breakerKeyPair{method: method, addr: addr1}); !ok {
			t.Fatalf("expected breakerKeys to cache %s", addr1)
		}
		g.evictAddr(addr1)
		if _, ok := g.breakerKeys.Load(breakerKeyPair{method: method, addr: addr1}); ok {
			t.Errorf("expected %s to be evicted from breakerKeys", addr1)
		}
		if _, ok := g.breakerKeys.Load(breakerKeyPair{method: method, addr: addr2}); !ok {
			t.Errorf("expected %s to remain in breakerKeys", addr2)
		}
	})

	t.Run("disconnectGen epoch detects stale pool snapshot", func(t *testing.T) {
		t.Parallel()
		const addr = "10.0.7.1:8081"
		cb := new(recordingBreaker)
		g := newFanOutTestClient(t, "disconnect-gen", WithCircuitBreaker(cb))
		p1 := fanOutConnPool{conn: newFanOutTestClientConn(t), healthy: true}
		g.storeConn(addr, p1)

		gen := g.disconnectGen.Load()
		if !g.useBreaker(addr, p1, true, gen) {
			t.Fatal("useBreaker should return true for matching epoch")
		}
		if err := g.Disconnect(t.Context(), addr); err != nil {
			t.Fatalf("Disconnect() error = %v", err)
		}
		if g.useBreaker(addr, p1, true, gen) {
			t.Error("useBreaker should return false after Disconnect increments epoch")
		}
	})

	t.Run("rpcCallFailedError matches ErrRPCCallFailed semantics", func(t *testing.T) {
		t.Parallel()
		const addr = "10.0.8.1:8081"
		base := status.Error(codes.NotFound, "missing")
		lazy := error(&rpcCallFailedError{addr: addr, err: base})
		eager := errors.ErrRPCCallFailed(addr, base)
		if lazy.Error() != eager.Error() {
			t.Errorf("Error() = %q, want %q", lazy.Error(), eager.Error())
		}
		if !errors.Is(lazy, eager) {
			t.Errorf("errors.Is(lazy, eager) = false")
		}
		if !errors.Is(lazy, base) {
			t.Errorf("errors.Is(lazy, base) = false")
		}
		nilWrap := error(&rpcCallFailedError{addr: addr})
		if nilWrap.Error() != "addr: "+addr {
			t.Errorf("nil err Error() = %q, want %q", nilWrap.Error(), "addr: "+addr)
		}
	})
}

func Test_gRPCClient_handleExecuteErr_fastPath(t *testing.T) {
	t.Parallel()
	const addr = "10.0.9.1:8081"
	g := newFanOutTestClient(t, "handle-err-fastpath")
	p := fanOutConnPool{conn: newFanOutTestClientConn(t), healthy: true}

	for _, tc := range []struct {
		err     error
		name    string
		quiet   bool
		noRetry bool
		wantRet bool
	}{
		{name: "context.Canceled", quiet: true, err: context.Canceled, wantRet: false},
		{name: "context.DeadlineExceeded", quiet: true, err: context.DeadlineExceeded, wantRet: false},
		{name: "NotFound quiet", quiet: true, err: status.Error(codes.NotFound, "nf"), wantRet: false},
		{name: "NotFound normal", err: status.Error(codes.NotFound, "nf"), wantRet: false},
		{name: "Unavailable retryable", quiet: true, err: status.Error(codes.Unavailable, "unavail"), wantRet: true},
		{name: "ResourceExhausted retryable", quiet: true, err: status.Error(codes.ResourceExhausted, "re"), wantRet: true},
		{name: "ResourceExhausted suppressed", quiet: true, noRetry: true, err: status.Error(codes.ResourceExhausted, "re"), wantRet: false},
		{name: "wrapped Unavailable", quiet: true, err: errors.Wrap(status.Error(codes.Unavailable, "unavail"), "wrapped"), wantRet: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			if tc.quiet {
				ctx = WithQuietErrorLog(ctx)
			}
			if tc.noRetry {
				ctx = WithoutResourceExhaustedRetry(ctx)
			}
			_, ret, gotErr := g.handleExecuteErr(ctx, p, addr, tc.err)
			if ret != tc.wantRet {
				t.Errorf("retryable = %v, want %v", ret, tc.wantRet)
			}
			if !errors.Is(gotErr, tc.err) {
				t.Errorf("err = %v, want %v", gotErr, tc.err)
			}
		})
	}
}

func Test_gRPCClient_RangeSliceConcurrent_and_FastPolicy(t *testing.T) {
	t.Parallel()

	t.Run("RangeSliceConcurrent covers cyclic windows and bounded concurrency", func(t *testing.T) {
		t.Parallel()
		const (
			addrA = "a:8081"
			addrB = "b:8081"
			addrC = "c:8081"
			addrD = "d:8081"
			addrE = "e:8081"
		)
		conn := newFanOutTestClientConn(t)
		g := newFanOutTestClient(t, "range-slice", WithErrGroup(errgroup.Get()))
		addrs := []string{addrA, addrB, addrC, addrD, addrE}
		for _, a := range addrs {
			g.storeConn(a, fanOutConnPool{conn: conn, healthy: true})
		}

		tests := []struct {
			name        string
			want        []string
			offset      int
			count       int
			concurrency int
		}{
			{
				name:        "wrapped cyclic slice",
				offset:      3,
				count:       3,
				concurrency: 0,
				want:        []string{addrD, addrE, addrA},
			},
			{
				name:        "single item",
				offset:      4,
				count:       1,
				concurrency: 0,
				want:        []string{addrE},
			},
			{
				name:        "sequential concurrency=1",
				offset:      4,
				count:       3,
				concurrency: 1,
				want:        []string{addrE, addrA, addrB},
			},
			{
				name:        "bounded concurrency=2",
				offset:      2,
				count:       4,
				concurrency: 2,
				want:        []string{addrC, addrD, addrE, addrA},
			},
			{
				name:        "count clamped to len(addrs)",
				offset:      -1,
				count:       10,
				concurrency: 0,
				want:        []string{addrE, addrA, addrB, addrC, addrD},
			},
			{
				name:        "count<=0 is noop",
				offset:      0,
				count:       0,
				concurrency: 0,
				want:        nil,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				var mu sync.Mutex
				got := make([]string, 0, len(tc.want))
				err := g.RangeSliceConcurrent(t.Context(), addrs, tc.offset, tc.count, tc.concurrency, func(_ context.Context, addr string, c *ClientConn, _ ...CallOption) error {
					if c == nil {
						t.Errorf("nil conn for %s", addr)
					}
					mu.Lock()
					got = append(got, addr)
					mu.Unlock()
					return nil
				})
				if err != nil {
					t.Fatalf("RangeSliceConcurrent() error = %v", err)
				}
				slices.Sort(got)
				wantSorted := slices.Clone(tc.want)
				slices.Sort(wantSorted)
				if !slices.Equal(got, wantSorted) {
					t.Errorf("visited = %v, want %v", got, wantSorted)
				}
			})
		}
	})

	t.Run("execVoidFast retries transient errors and respects circuit breaker", func(t *testing.T) {
		t.Parallel()
		const addr = "10.0.10.1:8081"
		bo := backoff.New(
			backoff.WithInitialDuration("1ms"),
			backoff.WithMaximumDuration("2ms"),
			backoff.WithJitterLimit("1ms"),
			backoff.WithBackOffTimeLimit("1s"),
			backoff.WithRetryCount(3),
		)
		cb, err := circuitbreaker.New(
			circuitbreaker.WithBreakerOpts(
				circuitbreaker.WithClosedErrorRate(0.5),
				circuitbreaker.WithMinSamples(3),
				circuitbreaker.WithOpenTimeout("200ms"),
			),
		)
		if err != nil {
			t.Fatal(err)
		}
		g := newFanOutTestClient(t, "fast-policy", WithBackoff(bo), WithCircuitBreaker(cb))
		p := fanOutConnPool{conn: newFanOutTestClientConn(t), healthy: true}
		g.storeConn(addr, p)

		// 1. Healthy call with empty method uses t.addrBr fast path.
		ctx := WithSearchBroadcastPolicy(t.Context())
		if rerr := g.RangeConcurrent(ctx, 0, func(context.Context, string, *ClientConn, ...CallOption) error {
			return nil
		}); rerr != nil {
			t.Fatalf("healthy RangeConcurrent error = %v", rerr)
		}

		// 2. Retryable error retries via RetryVoid and succeeds on attempt 2.
		var attempts atomic.Int32
		if rerr := g.doVoid(ctx, p, addr, func(context.Context, string, *ClientConn, ...CallOption) error {
			if attempts.Add(1) < 2 {
				return status.Error(codes.Unavailable, "transient")
			}
			return nil
		}); rerr != nil {
			t.Fatalf("retryable doVoid error = %v", rerr)
		}
		if attempts.Load() != 2 {
			t.Errorf("attempts = %d, want 2", attempts.Load())
		}

		// 3. Non-retryable NotFound does not retry and does not trip breaker.
		var nfCalls atomic.Int32
		err = g.doVoid(ctx, p, addr, func(context.Context, string, *ClientConn, ...CallOption) error {
			nfCalls.Add(1)
			return status.Error(codes.NotFound, "not found")
		})
		if err == nil || nfCalls.Load() != 1 {
			t.Fatalf("NotFound doVoid = (%v, calls=%d), want non-nil and calls=1", err, nfCalls.Load())
		}

		// 4. Persistent retryable failures trip the breaker to Open.
		for range 5 {
			if derr := g.doVoid(ctx, p, addr, func(context.Context, string, *ClientConn, ...CallOption) error {
				return status.Error(codes.Internal, "boom")
			}); derr == nil {
				t.Fatal("doVoid() returned nil error, want boom")
			}
		}
		err = g.doVoid(ctx, p, addr, func(context.Context, string, *ClientConn, ...CallOption) error {
			t.Fatal("breaker open should reject before invoking f")
			return nil
		})
		if !errors.Is(err, errors.ErrCircuitBreakerOpenState) {
			t.Fatalf("open breaker doVoid error = %v, want ErrCircuitBreakerOpenState", err)
		}
	})
}

const (
	fanOutRangeLabel   = "Range"
	fanOutOrderedLabel = "Ordered"
)

// BenchmarkClientFanOut measures one unbounded fan-out to n pools with a no-op callback.
func BenchmarkClientFanOut(b *testing.B) {
	log.Init(log.WithLoggerType("nop"))
	noop := func(context.Context, string, *ClientConn, ...CallOption) error { return nil }
	for _, ordered := range []bool{false, true} {
		method := fanOutRangeLabel
		if ordered {
			method = fanOutOrderedLabel
		}
		for _, n := range []int{1, 5, 50} {
			b.Run(fmt.Sprintf("%s/targets=%d", method, n), func(b *testing.B) {
				g, addrs := newFanOutBenchClient(b, n)
				ctx := context.Background()
				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					var err error
					if ordered {
						err = g.OrderedRangeConcurrent(ctx, addrs, 0, noop)
					} else {
						err = g.RangeConcurrent(ctx, 0, noop)
					}
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

//go:noinline

func fanOutStackWork(level int) int {
	var buf [16]int
	buf[level%len(buf)] = level
	if level == 0 {
		return buf[0]
	}
	return fanOutStackWork(level-1) + buf[level%len(buf)]
}

// BenchmarkClientFanOutStack measures bounded/pooled fan-out with deep recursion so retained stacks matter.
func BenchmarkClientFanOutStack(b *testing.B) {
	log.Init(log.WithLoggerType("nop"))
	const levels = 32
	want := fanOutStackWork(levels)
	for _, ordered := range []bool{false, true} {
		method := fanOutRangeLabel
		if ordered {
			method = fanOutOrderedLabel
		}
		for _, n := range []int{1, 5, 50} {
			b.Run(fmt.Sprintf("%s/targets=%d", method, n), func(b *testing.B) {
				g, addrs := newFanOutBenchClient(b, n)
				var bad atomic.Int32
				work := func(context.Context, string, *ClientConn, ...CallOption) error {
					if got := fanOutStackWork(levels); got != want {
						bad.Add(1)
						return errors.Errorf("stack work = %d, want %d", got, want)
					}
					return nil
				}
				ctx := context.Background()
				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					var err error
					if ordered {
						err = g.OrderedRangeConcurrent(ctx, addrs, 0, work)
					} else {
						err = g.RangeConcurrent(ctx, 0, work)
					}
					if err != nil {
						b.Fatal(err)
					}
				}
				if v := bad.Load(); v != 0 {
					b.Fatalf("%d callbacks computed a wrong result", v)
				}
			})
		}
	}
}

// BenchmarkClientFanOutParallel measures concurrent fan-out calls across parallel goroutines.
func BenchmarkClientFanOutParallel(b *testing.B) {
	log.Init(log.WithLoggerType("nop"))
	noop := func(context.Context, string, *ClientConn, ...CallOption) error { return nil }
	for _, ordered := range []bool{false, true} {
		method := fanOutRangeLabel
		if ordered {
			method = fanOutOrderedLabel
		}
		for _, n := range []int{3, 5, 16} {
			b.Run(fmt.Sprintf("%s/targets=%d", method, n), func(b *testing.B) {
				g, addrs := newFanOutBenchClient(b, n)
				ctx := context.Background()
				b.ReportAllocs()
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						var err error
						if ordered {
							err = g.OrderedRangeConcurrent(ctx, addrs, 0, noop)
						} else {
							err = g.RangeConcurrent(ctx, 0, noop)
						}
						if err != nil {
							b.Fatal(err)
						}
					}
				})
			})
		}
	}
}
