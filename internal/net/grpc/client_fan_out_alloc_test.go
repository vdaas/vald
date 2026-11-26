//go:build !race

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
	"strconv"
	"testing"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func Test_gRPCClient_FanOut_Allocs(t *testing.T) {
	conn := newFanOutTestClientConn(t)
	noop := func(context.Context, string, *ClientConn, ...CallOption) error { return nil }

	for _, n := range []int{1, 5, 50} {
		addrs := make([]string, n)
		g := newFanOutTestClient(t, "alloc-guard-"+strconv.Itoa(n), WithErrGroup(errgroup.Get()))
		for i := range n {
			addr := "10.1.0." + strconv.Itoa(i+1) + ":8081"
			addrs[i] = addr
			g.storeConn(addr, fanOutConnPool{conn: conn, healthy: true})
		}

		// Warm up WorkerPool workers and fanOutBatch/batchState pools.
		for range 5 {
			if err := g.RangeConcurrent(t.Context(), 0, noop); err != nil {
				t.Fatalf("warm RangeConcurrent(%d) error = %v", n, err)
			}
			if err := g.OrderedRangeConcurrent(t.Context(), addrs, 0, noop); err != nil {
				t.Fatalf("warm OrderedRangeConcurrent(%d) error = %v", n, err)
			}
			if err := g.RangeSliceConcurrent(t.Context(), addrs, 1, n, 0, noop); err != nil {
				t.Fatalf("warm RangeSliceConcurrent(%d) error = %v", n, err)
			}
		}

		if got := testing.AllocsPerRun(50, func() {
			if err := g.RangeConcurrent(t.Context(), 0, noop); err != nil {
				t.Errorf("RangeConcurrent(%d) error = %v", n, err)
			}
		}); got != 0 {
			t.Errorf("RangeConcurrent(%d) allocs/op = %.1f, want 0", n, got)
		}

		if got := testing.AllocsPerRun(50, func() {
			if err := g.OrderedRangeConcurrent(t.Context(), addrs, 0, noop); err != nil {
				t.Errorf("OrderedRangeConcurrent(%d) error = %v", n, err)
			}
		}); got != 0 {
			t.Errorf("OrderedRangeConcurrent(%d) allocs/op = %.1f, want 0", n, got)
		}

		if got := testing.AllocsPerRun(50, func() {
			if err := g.RangeSliceConcurrent(t.Context(), addrs, 1, n, 0, noop); err != nil {
				t.Errorf("RangeSliceConcurrent(%d) error = %v", n, err)
			}
		}); got != 0 {
			t.Errorf("RangeSliceConcurrent(%d) allocs/op = %.1f, want 0", n, got)
		}
	}

	t.Run("cachedBreakerName pre-suffixed zero alloc", func(t *testing.T) {
		const addr = "agent-0:8081"
		g := newFanOutTestClient(t, "alloc-presuf")
		ctx := WrapGRPCMethod(t.Context(), "BroadCast/"+addr)
		var sink atomic.Int64
		if got := testing.AllocsPerRun(100, func() {
			sink.Add(int64(len(g.cachedBreakerName(ctx, addr))))
		}); got != 0 {
			t.Errorf("cachedBreakerName(pre-suffixed) allocs/op = %.1f, want 0", got)
		}
	})

	t.Run("cachedBreakerName and handleExecuteErr NotFound zero alloc", func(t *testing.T) {
		const addr = "10.1.9.1:8081"
		g := newFanOutTestClient(t, "alloc-policy")
		p := fanOutConnPool{conn: conn, healthy: true}
		ctx := WithQuietErrorLog(WithGRPCMethod(t.Context(), "Exists"))
		notFound := status.Error(codes.NotFound, "not found")
		var sink atomic.Int64
		// Warm breakerKeys cache once.
		sink.Add(int64(len(g.cachedBreakerName(ctx, addr))))
		if got := testing.AllocsPerRun(100, func() {
			sink.Add(int64(len(g.cachedBreakerName(ctx, addr))))
		}); got != 0 {
			t.Errorf("cachedBreakerName allocs/op = %.1f, want 0", got)
		}
		if got := testing.AllocsPerRun(100, func() {
			_, _, err := g.handleExecuteErr(ctx, p, addr, notFound)
			if err == nil {
				t.Fatal("expected error")
			}
		}); got > 1 {
			t.Errorf("handleExecuteErr NotFound allocs/op = %.1f, want <= 1", got)
		}
	})

	t.Run("RangeConcurrent with healthy backoff and circuitbreaker zero alloc", func(t *testing.T) {
		bo := backoff.New(
			backoff.WithInitialDuration("5ms"),
			backoff.WithMaximumDuration("5s"),
			backoff.WithBackOffTimeLimit("5s"),
			backoff.WithRetryCount(10),
		)
		cb, err := circuitbreaker.New(
			circuitbreaker.WithBreakerOpts(
				circuitbreaker.WithClosedErrorRate(0.7),
				circuitbreaker.WithMinSamples(1000),
				circuitbreaker.WithOpenTimeout("1s"),
				circuitbreaker.WithClosedRefreshTimeout("10s"),
			),
		)
		if err != nil {
			t.Fatal(err)
		}
		g := newFanOutTestClient(t, "alloc-healthy-policy", WithErrGroup(errgroup.Get()), WithBackoff(bo), WithCircuitBreaker(cb))
		for i := range 5 {
			addr := "10.2.0." + strconv.Itoa(i+1) + ":8081"
			g.storeConn(addr, fanOutConnPool{conn: conn, healthy: true})
		}
		ctx := WithSearchBroadcastPolicy(t.Context())
		for range 5 {
			if err := g.RangeConcurrent(ctx, 0, noop); err != nil {
				t.Fatalf("warm RangeConcurrent error = %v", err)
			}
		}
		if got := testing.AllocsPerRun(50, func() {
			if err := g.RangeConcurrent(ctx, 0, noop); err != nil {
				t.Errorf("RangeConcurrent error = %v", err)
			}
		}); got != 0 {
			t.Errorf("RangeConcurrent with policy allocs/op = %.1f, want 0", got)
		}
	})
}
