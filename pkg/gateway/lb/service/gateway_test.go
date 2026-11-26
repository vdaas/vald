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

// Package service
package service

import (
	"context"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/credentials"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// rangeClient models gRPCClient.RangeConcurrent: f runs per target, up to
// attempts times while it fails (the client backoff), and every per-target
// error that is not a local ctx error is dropped from the return value.
type rangeClient struct {
	grpc.Client
	// onRange, when set, runs when the broadcast starts, to change the agent list
	// that the discoverer reports while the write is in flight.
	onRange  func()
	addrs    []string
	attempts int
}

func (c *rangeClient) OrderedRange(
	ctx context.Context,
	orders []string,
	f func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	if len(orders) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	for _, addr := range orders {
		if err := f(ctx, addr, nil); err != nil {
			return err
		}
	}
	return nil
}

func (c *rangeClient) OrderedRangeConcurrent(
	ctx context.Context,
	orders []string,
	concurrency int,
	f func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	if len(orders) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	eg, egctx := errgroup.New(ctx)
	if concurrency > 0 {
		eg.SetLimit(concurrency)
	}
	for _, addr := range orders {
		a := addr
		eg.Go(func() error {
			return f(egctx, a, nil)
		})
	}
	return eg.Wait()
}

func (c *rangeClient) RangeConcurrent(
	ctx context.Context,
	_ int,
	f func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	if c.onRange != nil {
		c.onRange()
	}
	if len(c.addrs) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	for _, addr := range c.addrs {
		for range max(c.attempts, 1) {
			if f(ctx, addr, nil) == nil {
				break
			}
		}
	}
	return nil
}

func (c *rangeClient) RangeSliceConcurrent(
	ctx context.Context,
	addrs []string,
	offset int,
	count int,
	_ int,
	f func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	if len(addrs) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	if count <= 0 {
		return nil
	}
	count = min(count, len(addrs))
	for i := range count {
		addr := addrs[(offset+i)%len(addrs)]
		if len(c.addrs) > 0 && !slices.Contains(c.addrs, addr) {
			continue
		}
		for range max(c.attempts, 1) {
			if f(ctx, addr, nil) == nil {
				break
			}
		}
	}
	return nil
}

type rangeDiscoverer struct {
	discoverer.Client
	c  grpc.Client
	rc grpc.Client
	// addrs are the agents the discoverer knows, reached or not.
	addrs []string
}

func (d *rangeDiscoverer) GetClient() grpc.Client { return d.c }

func (d *rangeDiscoverer) GetReadClient() grpc.Client {
	if d.rc != nil {
		return d.rc
	}
	return d.c
}

func (d *rangeDiscoverer) GetAddrs(context.Context) []string { return d.addrs }

// BroadCast returned whatever RangeConcurrent returned, so an agent failure
// (Unavailable/Internal/Aborted/...) reached the lb handlers as nil and they
// answered NotFound or OK instead of the agent's code.
func Test_gateway_BroadCast_propagatesCallbackErrors(t *testing.T) {
	t.Parallel()
	unavailable := status.Error(codes.Unavailable, "agent unavailable")
	aborted := status.Error(codes.Aborted, "flush in progress")
	internal := status.Error(codes.Internal, "agent internal")
	const addrA, addrB = "a:8081", "b:8081"
	tests := []struct {
		// errs maps target to the error of each successive attempt; a missing
		// entry or an exhausted list means success.
		errs     map[string][]error
		name     string
		addrs    []string
		wantIs   []error
		attempts int
		wantCode codes.Code // codes.OK: nil error expected
	}{
		{
			name:  "all agents succeed",
			addrs: []string{addrA, addrB},
		},
		{
			name:     "one agent unavailable",
			addrs:    []string{addrA, addrB},
			errs:     map[string][]error{addrB: {unavailable}},
			wantCode: codes.Unavailable,
			wantIs:   []error{unavailable},
		},
		{
			name:     "all agents aborted",
			addrs:    []string{addrA, addrB},
			errs:     map[string][]error{addrA: {aborted}, addrB: {aborted}},
			wantCode: codes.Aborted,
			wantIs:   []error{aborted},
		},
		{
			name:     "mixed codes are all kept and the first target's code leads",
			addrs:    []string{addrB, addrA},
			errs:     map[string][]error{addrA: {internal}, addrB: {aborted}},
			wantCode: codes.Internal,
			wantIs:   []error{internal, aborted},
		},
		{
			name:     "a retry that succeeds clears the earlier error",
			addrs:    []string{addrA},
			attempts: 2,
			errs:     map[string][]error{addrA: {unavailable}},
		},
		{
			name:     "a retry that fails again keeps the last error",
			addrs:    []string{addrA},
			attempts: 2,
			errs:     map[string][]error{addrA: {unavailable, internal}},
			wantCode: codes.Internal,
			wantIs:   []error{internal},
		},
	}
	for _, test := range tests {
		for _, kind := range []BroadCastKind{READ, WRITE} {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()
				g := &gateway{client: &rangeDiscoverer{c: &rangeClient{addrs: test.addrs, attempts: test.attempts}}}
				calls := make(map[string]int) // the model calls f sequentially
				err := g.BroadCast(t.Context(), kind, func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
					n := calls[target]
					calls[target]++
					if errs := test.errs[target]; n < len(errs) {
						return errs[n]
					}
					return nil
				})
				if test.wantCode == codes.OK {
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}
					return
				}
				st, ok := status.FromError(err)
				if !ok || st == nil || st.Code() != test.wantCode {
					t.Fatalf("got error: %v, want code %v", err, test.wantCode)
				}
				for _, want := range test.wantIs {
					if !errors.Is(err, want) {
						t.Errorf("error %v does not contain %v", err, want)
					}
				}
			})
		}
	}
}

// No connection at all is still reported as ErrGRPCClientConnNotFound, which the
// handlers map to Internal.
func Test_gateway_BroadCast_connNotFound(t *testing.T) {
	t.Parallel()
	g := &gateway{client: &rangeDiscoverer{c: new(rangeClient)}}
	err := g.BroadCast(t.Context(), WRITE, func(context.Context, string, vald.Client, ...grpc.CallOption) error {
		t.Error("callback must not run without connections")
		return nil
	})
	if !errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
		t.Fatalf("got %v, want connection not found", err)
	}
}

// RangeConcurrent returns nil for a target whose connection could not be acquired
// (or whose circuit breaker is open): the callback never runs and the target's
// error is dropped. A write BroadCast must still report such a target, otherwise
// Remove/Flush look successful while that agent keeps the object. The agents that
// are reachable must be written regardless.
func Test_gateway_BroadCast_unreachedTargets(t *testing.T) {
	t.Parallel()
	const addrA, addrB, addrC = "a:8081", "b:8081", "c:8081"
	known := []string{addrA, addrB, addrC}
	tests := []struct {
		cbErr     map[string]error
		name      string
		reached   []string // targets RangeConcurrent runs the callback for
		wantInErr []string // substrings that the returned error must contain
		kind      BroadCastKind
		wantCode  codes.Code
	}{
		{
			name:     "write: an unreached agent is reported as Unavailable",
			kind:     WRITE,
			reached:  []string{addrA, addrC},
			wantCode: codes.Unavailable, wantInErr: []string{addrB},
		},
		{
			name:     "write: callback failures and unreached agents are both reported",
			kind:     WRITE,
			reached:  []string{addrA, addrB},
			cbErr:    map[string]error{addrB: status.Error(codes.Internal, "agent internal")},
			wantCode: codes.Internal, wantInErr: []string{"agent internal", addrC},
		},
		{
			name:    "write: all agents reached",
			kind:    WRITE,
			reached: known,
		},
		{
			name:    "read: unreached agents keep being ignored (partial search results)",
			kind:    READ,
			reached: []string{addrA},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			g := &gateway{client: &rangeDiscoverer{c: &rangeClient{addrs: test.reached}, addrs: known}}
			var mu sync.Mutex
			called := map[string]bool{}
			err := g.BroadCast(t.Context(), test.kind, func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
				mu.Lock()
				called[target] = true
				mu.Unlock()
				return test.cbErr[target]
			})
			for _, addr := range test.reached {
				if !called[addr] {
					t.Errorf("reachable agent %s was not written", addr)
				}
			}
			if test.wantCode == codes.OK {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("got nil error, want the unreached agent reported")
			}
			if st, ok := status.FromError(err); !ok || st == nil || st.Code() != test.wantCode {
				t.Errorf("got error %v, want code %v", err, test.wantCode)
			}
			for _, sub := range test.wantInErr {
				if !strings.Contains(err.Error(), sub) {
					t.Errorf("error %q does not mention %q", err, sub)
				}
			}
		})
	}
}

// The agents a write has to reach are the ones known when it starts. An agent that
// falls out of the discoverer's list during the broadcast and was not reached must
// still be reported: it can rejoin with the object the write was meant to remove.
// An agent that joins the list during the broadcast is not a target of this write.
func Test_gateway_BroadCast_agentListChangesDuringBroadcast(t *testing.T) {
	t.Parallel()
	const addrA, addrB, addrC = "a:8081", "b:8081", "c:8081"
	tests := []struct {
		name      string
		wantInErr string   // empty: nil error expected
		before    []string // the discoverer's list when the broadcast starts
		after     []string // the list after the broadcast has started
		reached   []string // targets RangeConcurrent runs the callback for
	}{
		{
			name:      "an unreached agent that leaves the list is still reported",
			before:    []string{addrA, addrB},
			after:     []string{addrA},
			reached:   []string{addrA},
			wantInErr: addrB,
		},
		{
			name:    "an unreached agent that joins the list is not a target of this write",
			before:  []string{addrA},
			after:   []string{addrA, addrC},
			reached: []string{addrA},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			d := &rangeDiscoverer{addrs: test.before}
			d.c = &rangeClient{
				addrs:   test.reached,
				onRange: func() { d.addrs = test.after },
			}
			g := &gateway{client: d}
			err := g.BroadCast(t.Context(), WRITE, func(context.Context, string, vald.Client, ...grpc.CallOption) error { return nil })
			if test.wantInErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("got nil error, want the agent that left the list reported as Unavailable")
			}
			if st, ok := status.FromError(err); !ok || st == nil || st.Code() != codes.Unavailable {
				t.Errorf("got error %v, want code %v", err, codes.Unavailable)
			}
			if !strings.Contains(err.Error(), test.wantInErr) {
				t.Errorf("error %q does not mention %q", err, test.wantInErr)
			}
		})
	}
}

func Test_gateway_DoMulti(t *testing.T) {
	t.Parallel()
	tests := []struct {
		wantErr   error
		f         func(visited *sync.Map[string, int]) func(context.Context, string, vald.Client, ...grpc.CallOption) error
		wantCalls map[string]int
		name      string
		addrs     []string
		num       int
	}{
		{
			name:  "concurrent pass 1 with 3 replicas out of 5 agents succeeds",
			num:   3,
			addrs: []string{"a", "b", "c", "d", "e"},
			f: func(visited *sync.Map[string, int]) func(context.Context, string, vald.Client, ...grpc.CallOption) error {
				return func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
					v, _ := visited.LoadOrStore(target, 0)
					visited.Store(target, v+1)
					return nil
				}
			},
			wantCalls: map[string]int{
				"a": 1,
				"b": 1,
				"c": 1,
			},
		},
		{
			name:  "second pass retries failed target to satisfy replica",
			num:   2,
			addrs: []string{"a", "b", "c"},
			f: func(visited *sync.Map[string, int]) func(context.Context, string, vald.Client, ...grpc.CallOption) error {
				var aCount atomic.Int32
				return func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
					if target == "a" && aCount.Add(1) == 1 {
						return errors.New("agent a temporary error")
					}
					v, _ := visited.LoadOrStore(target, 0)
					visited.Store(target, v+1)
					return nil
				}
			},
			wantCalls: map[string]int{
				"b": 1,
				"a": 1,
			},
		},
		{
			name:  "empty addrs returns ErrGRPCClientConnNotFound",
			num:   2,
			addrs: []string{},
			f: func(visited *sync.Map[string, int]) func(context.Context, string, vald.Client, ...grpc.CallOption) error {
				return func(_ context.Context, _ string, _ vald.Client, _ ...grpc.CallOption) error {
					return nil
				}
			},
			wantErr: errors.ErrGRPCClientConnNotFound("*"),
		},
		{
			name:  "num is 0 returns nil with 0 calls",
			num:   0,
			addrs: []string{"a", "b"},
			f: func(visited *sync.Map[string, int]) func(context.Context, string, vald.Client, ...grpc.CallOption) error {
				return func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
					v, _ := visited.LoadOrStore(target, 0)
					visited.Store(target, v+1)
					return nil
				}
			},
			wantCalls: map[string]int{},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx := tt.Context()
			var visited sync.Map[string, int]
			rc := &rangeClient{addrs: test.addrs}
			rd := &rangeDiscoverer{c: rc, addrs: test.addrs}
			g := &gateway{client: rd}

			err := g.DoMulti(ctx, test.num, test.f(&visited))
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					tt.Fatalf("err = %v, want %v", err, test.wantErr)
				}
				return
			}
			if err != nil {
				tt.Fatalf("unexpected error = %v", err)
			}
			for target, wantCount := range test.wantCalls {
				val, ok := visited.Load(target)
				if !ok || val != wantCount {
					tt.Errorf("target %s called %d times, want %d", target, val, wantCount)
				}
			}
			visited.Range(func(key string, val int) bool {
				if _, ok := test.wantCalls[key]; !ok {
					tt.Errorf("unexpected target %s called %d times", key, val)
				}
				return true
			})
		})
	}
}

func Test_gateway_DoMulti_RoundRobin(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	addrs := []string{"a", "b", "c", "d", "e"}
	rc := &rangeClient{addrs: addrs}
	rd := &rangeDiscoverer{c: rc, addrs: addrs}
	g := &gateway{client: rd}

	var mu sync.Mutex
	counts := make(map[string]int, len(addrs))
	const rounds = 5
	for round := range rounds {
		err := g.DoMulti(ctx, 3, func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
			mu.Lock()
			counts[target]++
			mu.Unlock()
			return nil
		})
		if err != nil {
			t.Fatalf("round %d: unexpected error = %v", round, err)
		}
	}
	for _, addr := range addrs {
		if got := counts[addr]; got != 3 {
			t.Errorf("target %s called %d times over %d rounds, want 3 (counts=%v)", addr, got, rounds, counts)
		}
	}
}

func Test_gateway_getValdClient_caching(t *testing.T) {
	t.Parallel()
	g := new(gateway)
	conn := new(grpc.ClientConn)
	c1 := g.getValdClient(conn)
	c2 := g.getValdClient(conn)
	if c1 == nil || c2 == nil {
		t.Fatal("expected non-nil client")
	}
	if c1 != c2 {
		t.Errorf("expected cached client instance to be reused, got %p != %p", c1, c2)
	}
}

type concurrentRangeClient struct {
	rangeClient
}

func (c *concurrentRangeClient) RangeConcurrent(
	ctx context.Context,
	concurrency int,
	f func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	if len(c.addrs) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	eg, egctx := errgroup.New(ctx)
	if concurrency > 0 {
		eg.SetLimit(concurrency)
	}
	attempts := max(c.attempts, 1)
	for _, addr := range c.addrs {
		a := addr
		eg.Go(func() error {
			for range attempts {
				if f(egctx, a, nil) == nil {
					break
				}
			}
			return nil
		})
	}
	return eg.Wait()
}

func Test_gateway_getValdClient_nilSafe(t *testing.T) {
	t.Parallel()
	g := new(gateway)
	ac := g.getValdClient(nil)
	if ac == nil {
		t.Fatal("getValdClient(nil) returned nil client")
	}
	if _, err := ac.Search(t.Context(), nil); !errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
		t.Fatalf("getValdClient(nil).Search() error = %v, want %v", err, errors.ErrGRPCClientConnNotFound("*"))
	}
}

func Test_gateway_BroadCast_concurrentReadAndWrite(t *testing.T) {
	t.Parallel()
	const addrA, addrB, addrC = "a:8081", "b:8081", "c:8081"
	unavailable := status.Error(codes.Unavailable, "agent unavailable")

	tests := []struct {
		checkErr func(t *testing.T, err error, calls int32)
		fn       func(bAttempts *atomic.Int32) func(context.Context, string, vald.Client, ...grpc.CallOption) error
		name     string
		addrs    []string
		rcAddrs  []string
		attempts int
		kind     BroadCastKind
	}{
		{
			name:    "concurrent READ succeeds without locking ferr",
			addrs:   []string{addrA, addrB, addrC},
			rcAddrs: []string{addrA, addrB, addrC},
			kind:    READ,
			fn: func(_ *atomic.Int32) func(context.Context, string, vald.Client, ...grpc.CallOption) error {
				return func(context.Context, string, vald.Client, ...grpc.CallOption) error { return nil }
			},
			checkErr: func(t *testing.T, err error, calls int32) {
				t.Helper()
				if err != nil {
					t.Fatalf("BroadCast(READ) error = %v, want nil", err)
				}
				if calls != 3 {
					t.Fatalf("calls = %d, want 3", calls)
				}
			},
		},
		{
			name:     "concurrent READ retry clears transient error",
			addrs:    []string{addrA, addrB, addrC},
			rcAddrs:  []string{addrA, addrB, addrC},
			attempts: 2,
			kind:     READ,
			fn: func(bAttempts *atomic.Int32) func(context.Context, string, vald.Client, ...grpc.CallOption) error {
				return func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
					if target == addrB && bAttempts.Add(1) == 1 {
						return unavailable
					}
					return nil
				}
			},
			checkErr: func(t *testing.T, err error, _ int32) {
				t.Helper()
				if err != nil {
					t.Fatalf("BroadCast(READ) error = %v, want nil after retry recovery", err)
				}
			},
		},
		{
			name:    "concurrent READ preserves failure under race",
			addrs:   []string{addrA, addrB, addrC},
			rcAddrs: []string{addrA, addrB, addrC},
			kind:    READ,
			fn: func(_ *atomic.Int32) func(context.Context, string, vald.Client, ...grpc.CallOption) error {
				return func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
					if target == addrB {
						return unavailable
					}
					return nil
				}
			},
			checkErr: func(t *testing.T, err error, _ int32) {
				t.Helper()
				if !errors.Is(err, unavailable) {
					t.Fatalf("BroadCast(READ) error = %v, want %v", err, unavailable)
				}
			},
		},
		{
			name:    "concurrent WRITE reports unreached agent via pooled ran map",
			addrs:   []string{addrA, addrB, addrC},
			rcAddrs: []string{addrA, addrC},
			kind:    WRITE,
			fn: func(_ *atomic.Int32) func(context.Context, string, vald.Client, ...grpc.CallOption) error {
				return func(context.Context, string, vald.Client, ...grpc.CallOption) error { return nil }
			},
			checkErr: func(t *testing.T, err error, _ int32) {
				t.Helper()
				if st, ok := status.FromError(err); !ok || st.Code() != codes.Unavailable || !strings.Contains(err.Error(), addrB) {
					t.Fatalf("BroadCast(WRITE) error = %v, want Unavailable referencing %s", err, addrB)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			rc := &concurrentRangeClient{addrs: tc.rcAddrs, attempts: tc.attempts}
			rd := &rangeDiscoverer{c: rc, addrs: tc.addrs}
			g := &gateway{client: rd}

			var calls atomic.Int32
			var bAttempts atomic.Int32
			cb := tc.fn(&bAttempts)
			err := g.BroadCast(tt.Context(), tc.kind, func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error {
				calls.Add(1)
				return cb(ctx, target, ac, copts...)
			})
			tc.checkErr(tt, err, calls.Load())
		})
	}
}

func Test_gateway_DoMulti_concurrentAndEdgeCases(t *testing.T) {
	t.Parallel()
	t.Run("negative num is treated as zero replicas", func(tt *testing.T) {
		tt.Parallel()
		rc := &rangeClient{addrs: []string{"a", "b"}}
		rd := &rangeDiscoverer{c: rc, addrs: rc.addrs}
		g := &gateway{client: rd}
		var calls atomic.Int32
		err := g.DoMulti(tt.Context(), -1, func(_ context.Context, _ string, _ vald.Client, _ ...grpc.CallOption) error {
			calls.Add(1)
			return nil
		})
		if err != nil {
			tt.Fatalf("DoMulti(-1) error = %v, want nil", err)
		}
		if got := calls.Load(); got != 0 {
			tt.Fatalf("DoMulti(-1) calls = %d, want 0", got)
		}
	})

	t.Run("repeated concurrent DoMulti with fallback resets pooled visited slice", func(tt *testing.T) {
		tt.Parallel()
		addrs := []string{"a", "b", "c", "d"}
		for iter := range 16 {
			var failA atomic.Bool
			rc := &rangeClient{addrs: addrs}
			rd := &rangeDiscoverer{c: rc, addrs: addrs}
			g := &gateway{client: rd}
			var visited sync.Map[string, int]
			err := g.DoMulti(tt.Context(), 3, func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
				if iter%2 == 1 && target == "a" && !failA.Swap(true) {
					return errors.New("transient failure")
				}
				visited.Store(target, 1)
				return nil
			})
			if err != nil {
				tt.Fatalf("iter %d: DoMulti error = %v", iter, err)
			}
			var total int
			visited.Range(func(_ string, _ int) bool {
				total++
				return true
			})
			if total != 3 {
				tt.Fatalf("iter %d: reached %d distinct targets, want 3", iter, total)
			}
		}
	})
}

type sliceRangeClient struct {
	rangeClient
	sliceCalls atomic.Int64
	lastOffset atomic.Int64
	lastCount  atomic.Int64
}

func (c *sliceRangeClient) RangeSliceConcurrent(
	ctx context.Context,
	addrs []string,
	offset int,
	count int,
	concurrency int,
	f func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	c.sliceCalls.Add(1)
	c.lastOffset.Store(int64(offset))
	c.lastCount.Store(int64(count))
	if len(addrs) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	if count <= 0 {
		return nil
	}
	count = min(count, len(addrs))
	eg, egctx := errgroup.New(ctx)
	if concurrency > 0 {
		eg.SetLimit(concurrency)
	}
	for i := range count {
		addr := addrs[(offset+i)%len(addrs)]
		if len(c.addrs) > 0 && !slices.Contains(c.addrs, addr) {
			continue
		}
		eg.Go(func() error {
			for range max(c.attempts, 1) {
				if f(egctx, addr, nil) == nil {
					break
				}
			}
			return nil
		})
	}
	return eg.Wait()
}

func Test_gateway_vcTable_and_pruning(t *testing.T) {
	t.Parallel()
	g := new(gateway)
	conns := make([]*grpc.ClientConn, 12)
	for i := range conns {
		cc, err := grpc.NewClient("passthrough:///vc-test-"+strconv.Itoa(i),
			grpc.WithTransportCredentials(credentials.NewInsecure()))
		if err != nil {
			t.Fatal(err)
		}
		conns[i] = cc
	}
	t.Cleanup(func() {
		for _, cc := range conns[1:] {
			if err := cc.Close(); err != nil {
				t.Errorf("close conn: %v", err)
			}
		}
	})

	clients := make([]vald.Client, len(conns))
	for i, cc := range conns {
		clients[i] = g.getValdClient(cc)
		if clients[i] == nil {
			t.Fatalf("getValdClient(%d) returned nil", i)
		}
		if got := g.getValdClient(cc); got != clients[i] {
			t.Fatalf("getValdClient(%d) did not return cached instance", i)
		}
	}

	// Close first conn (moves to Shutdown) and insert another to trigger pruning.
	if err := conns[0].Close(); err != nil {
		t.Fatal(err)
	}
	extra, err := grpc.NewClient("passthrough:///vc-test-extra",
		grpc.WithTransportCredentials(credentials.NewInsecure()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := extra.Close(); err != nil {
			t.Errorf("close extra conn: %v", err)
		}
	})
	_ = g.getValdClient(extra)
	if _, ok := g.vcCache.Load(extra); !ok {
		t.Fatal("expected vcCache to contain extra conn")
	}
	if _, ok := g.vcCache.Load(conns[0]); ok {
		t.Error("expected Shutdown conn to be pruned from vcCache")
	}
}

func Test_gateway_DoMulti_RangeSliceClient(t *testing.T) {
	t.Parallel()
	addrs := []string{"a", "b", "c", "d", "e"}
	src := &sliceRangeClient{addrs: addrs}
	rd := &rangeDiscoverer{c: src, addrs: addrs}
	g := &gateway{client: rd}

	counts := make(map[string]int, len(addrs))
	var mu sync.Mutex
	for round := range len(addrs) {
		err := g.DoMulti(t.Context(), 3, func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
			mu.Lock()
			counts[target]++
			mu.Unlock()
			return nil
		})
		if err != nil {
			t.Fatalf("round %d DoMulti error = %v", round, err)
		}
	}
	if src.sliceCalls.Load() != int64(len(addrs)) {
		t.Errorf("RangeSliceConcurrent calls = %d, want %d", src.sliceCalls.Load(), len(addrs))
	}
	for _, a := range addrs {
		if counts[a] != 3 {
			t.Errorf("target %s called %d times, want 3", a, counts[a])
		}
	}
}

func Test_gateway_SubsetGateway(t *testing.T) {
	t.Parallel()
	const (
		addrA = "a:8081"
		addrB = "b:8081"
		addrC = "c:8081"
		addrD = "d:8081"
		addrE = "e:8081"
	)
	addrs := []string{addrA, addrB, addrC, addrD, addrE}

	t.Run("NextRotateOffset and BroadCastSlice READ/WRITE", func(t *testing.T) {
		t.Parallel()
		src := &sliceRangeClient{addrs: addrs}
		rd := &rangeDiscoverer{c: src, addrs: addrs}
		var gw Gateway = &gateway{client: rd}
		if off := gw.NextRotateOffset(1); off != 0 {
			t.Errorf("NextRotateOffset(1) = %d, want 0", off)
		}
		for i := range len(addrs) {
			if off := gw.NextRotateOffset(len(addrs)); off != i {
				t.Errorf("NextRotateOffset(%d) = %d, want %d", len(addrs), off, i)
			}
		}

		// Wrapped slice: offset=3, count=3 -> d, e, a.
		var mu sync.Mutex
		var visited []string
		sliced, err := gw.BroadCastSlice(t.Context(), READ, addrs, 3, 3, func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
			mu.Lock()
			visited = append(visited, target)
			mu.Unlock()
			return nil
		})
		if err != nil || !sliced {
			t.Fatalf("BroadCastSlice(READ) = (%v, %v), want (true, nil)", sliced, err)
		}
		slices.Sort(visited)
		if !slices.Equal(visited, []string{addrA, addrD, addrE}) {
			t.Errorf("visited = %v, want [%s %s %s]", visited, addrA, addrD, addrE)
		}

		// WRITE reports unreached target within slice only.
		partial := &sliceRangeClient{addrs: []string{addrD, addrA}}
		gwPartial := &gateway{client: &rangeDiscoverer{c: partial, addrs: addrs}}
		sliced, err = gwPartial.BroadCastSlice(t.Context(), WRITE, addrs, 3, 3, func(context.Context, string, vald.Client, ...grpc.CallOption) error {
			return nil
		})
		if !sliced {
			t.Fatalf("BroadCastSlice(WRITE) sliced = false, want true")
		}
		if st, ok := status.FromError(err); !ok || st.Code() != codes.Unavailable || !strings.Contains(err.Error(), addrE) {
			t.Fatalf("BroadCastSlice(WRITE) error = %v, want Unavailable mentioning %s", err, addrE)
		}
	})

	t.Run("BroadCastSlice dedicated ReadReplica returns sliced=false and visits read targets once", func(t *testing.T) {
		t.Parallel()
		readTargets := []string{"rr-0:8081", "rr-1:8081"}
		primary := &sliceRangeClient{addrs: addrs}
		readOnly := &sliceRangeClient{addrs: readTargets}
		var gw Gateway = &gateway{client: &rangeDiscoverer{c: primary, rc: readOnly, addrs: addrs}}

		var mu sync.Mutex
		var visited []string
		sliced, err := gw.BroadCastSlice(t.Context(), READ, addrs, 0, 3, func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
			mu.Lock()
			visited = append(visited, target)
			mu.Unlock()
			return nil
		})
		if err != nil || sliced {
			t.Fatalf("BroadCastSlice(READ with ReadReplica) = (%v, %v), want (false, nil)", sliced, err)
		}
		slices.Sort(visited)
		if !slices.Equal(visited, readTargets) {
			t.Errorf("visited = %v, want %v", visited, readTargets)
		}

		visited = nil
		sliced, err = gw.BroadCastSlice(t.Context(), READ_PRIMARY, addrs, 3, 2, func(_ context.Context, target string, _ vald.Client, _ ...grpc.CallOption) error {
			mu.Lock()
			visited = append(visited, target)
			mu.Unlock()
			return nil
		})
		if err != nil || !sliced {
			t.Fatalf("BroadCastSlice(READ_PRIMARY) = (%v, %v), want (true, nil)", sliced, err)
		}
		slices.Sort(visited)
		if !slices.Equal(visited, []string{addrD, addrE}) {
			t.Errorf("READ_PRIMARY visited = %v, want [%s %s]", visited, addrD, addrE)
		}
	})
}

func BenchmarkGateway_getValdClient(b *testing.B) {
	log.Init(log.WithLoggerType("nop"))
	g := new(gateway)
	conns := make([]*grpc.ClientConn, 5)
	for i := range conns {
		cc, err := grpc.NewClient("passthrough:///vc-bench-"+strconv.Itoa(i),
			grpc.WithTransportCredentials(credentials.NewInsecure()))
		if err != nil {
			b.Fatal(err)
		}
		conns[i] = cc
		_ = g.getValdClient(cc)
	}
	b.Cleanup(func() {
		for _, cc := range conns {
			if err := cc.Close(); err != nil {
				b.Errorf("close conn: %v", err)
			}
		}
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		if c := g.getValdClient(conns[i%len(conns)]); c == nil {
			b.Fatal("nil client")
		}
	}
}

func BenchmarkGateway_DoMulti_RangeSlice(b *testing.B) {
	const (
		addrA = "a:8081"
		addrB = "b:8081"
		addrC = "c:8081"
		addrD = "d:8081"
		addrE = "e:8081"
	)
	addrs := []string{addrA, addrB, addrC, addrD, addrE}
	src := &sliceRangeClient{addrs: addrs}
	g := &gateway{client: &rangeDiscoverer{c: src, addrs: addrs}}
	ctx := context.Background()
	noop := func(context.Context, string, vald.Client, ...grpc.CallOption) error { return nil }

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := g.DoMulti(ctx, 3, noop); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGateway_BroadCastSlice_READ(b *testing.B) {
	const (
		addrA = "a:8081"
		addrB = "b:8081"
		addrC = "c:8081"
		addrD = "d:8081"
		addrE = "e:8081"
	)
	addrs := []string{addrA, addrB, addrC, addrD, addrE}
	src := &sliceRangeClient{addrs: addrs}
	g := &gateway{client: &rangeDiscoverer{c: src, addrs: addrs}}
	ctx := context.Background()
	noop := func(context.Context, string, vald.Client, ...grpc.CallOption) error { return nil }

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		if _, err := g.BroadCastSlice(ctx, READ, addrs, i%len(addrs), 3, noop); err != nil {
			b.Fatal(err)
		}
	}
}

// NOT IMPLEMENTED BELOW
//
//	func TestNewGateway(t *testing.T) {
//		type args struct {
//			opts []Option
//		}
//		type want struct {
//			wantGw Gateway
//			err    error
//		}
//		type test struct {
//			name       string
//			args
//			want
//			checkFunc  func(want, Gateway, error) error
//			beforeFunc func(*testing.T, args)
//			afterFunc  func(*testing.T, args)
//		}
//		defaultCheckFunc := func(w want, gotGw Gateway, err error) error {
//			if !errors.Is(err, w.err) {
//				return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
//			}
//			if !reflect.DeepEqual(gotGw, w.wantGw) {
//				return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotGw, w.wantGw)
//			}
//			return nil
//		}
//		tests := []test{
//			// TODO test cases
//			/*
//			   {
//			       name: "test_case_1",
//			       args: args {
//			           opts:nil,
//			       },
//			       want: want{},
//			       checkFunc: defaultCheckFunc,
//			       beforeFunc: func(t *testing.T, args args) {
//			           t.Helper()
//			       },
//			       afterFunc: func(t *testing.T, args args) {
//			           t.Helper()
//			       },
//			   },
//			*/
//
//			// TODO test cases
//			/*
//			   func() test {
//			       return test {
//			           name: "test_case_2",
//			           args: args {
//			           opts:nil,
//			           },
//			           want: want{},
//			           checkFunc: defaultCheckFunc,
//			           beforeFunc: func(t *testing.T, args args) {
//			               t.Helper()
//			           },
//			           afterFunc: func(t *testing.T, args args) {
//			               t.Helper()
//			           },
//			       }
//			   }(),
//			*/
//		}
//
//		for _, tc := range tests {
//			test := tc
//			t.Run(test.name, func(tt *testing.T) {
//				tt.Parallel()
//				defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
//				if test.beforeFunc != nil {
//					test.beforeFunc(tt, test.args)
//				}
//				if test.afterFunc != nil {
//					defer test.afterFunc(tt, test.args)
//				}
//				checkFunc := test.checkFunc
//				if test.checkFunc == nil {
//					checkFunc = defaultCheckFunc
//				}
//
//				gotGw, err := NewGateway(test.args.opts...)
//				if err := checkFunc(test.want, gotGw, err); err != nil {
//					tt.Errorf("error = %v", err)
//				}
//			})
//		}
//	}
//
//	func Test_gateway_Start(t *testing.T) {
//		type args struct {
//			ctx context.Context
//		}
//		type fields struct {
//			client discoverer.Client
//			eg     errgroup.Group
//		}
//		type want struct {
//			want <-chan error
//			err  error
//		}
//		type test struct {
//			name       string
//			args
//			fields
//			want
//			checkFunc  func(want, <-chan error, error) error
//			beforeFunc func(*testing.T, args)
//			afterFunc  func(*testing.T, args)
//		}
//		defaultCheckFunc := func(w want, got <-chan error, err error) error {
//			if !errors.Is(err, w.err) {
//				return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
//			}
//			if !reflect.DeepEqual(got, w.want) {
//				return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
//			}
//			return nil
//		}
//		tests := []test{
//			// TODO test cases
//			/*
//			   {
//			       name: "test_case_1",
//			       args: args {
//			           ctx:nil,
//			       },
//			       fields: fields {
//			           client:nil,
//			           eg:nil,
//			       },
//			       want: want{},
//			       checkFunc: defaultCheckFunc,
//			       beforeFunc: func(t *testing.T, args args) {
//			           t.Helper()
//			       },
//			       afterFunc: func(t *testing.T, args args) {
//			           t.Helper()
//			       },
//			   },
//			*/
//
//			// TODO test cases
//			/*
//			   func() test {
//			       return test {
//			           name: "test_case_2",
//			           args: args {
//			           ctx:nil,
//			           },
//			           fields: fields {
//			           client:nil,
//			           eg:nil,
//			           },
//			           want: want{},
//			           checkFunc: defaultCheckFunc,
//			           beforeFunc: func(t *testing.T, args args) {
//			               t.Helper()
//			           },
//			           afterFunc: func(t *testing.T, args args) {
//			               t.Helper()
//			           },
//			       }
//			   }(),
//			*/
//		}
//
//		for _, tc := range tests {
//			test := tc
//			t.Run(test.name, func(tt *testing.T) {
//				tt.Parallel()
//				defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
//				if test.beforeFunc != nil {
//					test.beforeFunc(tt, test.args)
//				}
//				if test.afterFunc != nil {
//					defer test.afterFunc(tt, test.args)
//				}
//				checkFunc := test.checkFunc
//				if test.checkFunc == nil {
//					checkFunc = defaultCheckFunc
//				}
//				g := &gateway{
//					client: test.fields.client,
//					eg:     test.fields.eg,
//				}
//
//				got, err := g.Start(test.args.ctx)
//				if err := checkFunc(test.want, got, err); err != nil {
//					tt.Errorf("error = %v", err)
//				}
//			})
//		}
//	}
//
//	func Test_gateway_BroadCast(t *testing.T) {
//		type args struct {
//			ctx  context.Context
//			kind BroadCastKind
//			f    func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error
//		}
//		type fields struct {
//			client discoverer.Client
//			eg     errgroup.Group
//		}
//		type want struct {
//			err error
//		}
//		type test struct {
//			name       string
//			args
//			fields
//			want
//			checkFunc  func(want, error) error
//			beforeFunc func(*testing.T, args)
//			afterFunc  func(*testing.T, args)
//		}
//		defaultCheckFunc := func(w want, err error) error {
//			if !errors.Is(err, w.err) {
//				return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
//			}
//			return nil
//		}
//		tests := []test{
//			// TODO test cases
//			/*
//			   {
//			       name: "test_case_1",
//			       args: args {
//			           ctx:nil,
//			           kind:nil,
//			           f:nil,
//			       },
//			       fields: fields {
//			           client:nil,
//			           eg:nil,
//			       },
//			       want: want{},
//			       checkFunc: defaultCheckFunc,
//			       beforeFunc: func(t *testing.T, args args) {
//			           t.Helper()
//			       },
//			       afterFunc: func(t *testing.T, args args) {
//			           t.Helper()
//			       },
//			   },
//			*/
//
//			// TODO test cases
//			/*
//			   func() test {
//			       return test {
//			           name: "test_case_2",
//			           args: args {
//			           ctx:nil,
//			           kind:nil,
//			           f:nil,
//			           },
//			           fields: fields {
//			           client:nil,
//			           eg:nil,
//			           },
//			           want: want{},
//			           checkFunc: defaultCheckFunc,
//			           beforeFunc: func(t *testing.T, args args) {
//			               t.Helper()
//			           },
//			           afterFunc: func(t *testing.T, args args) {
//			               t.Helper()
//			           },
//			       }
//			   }(),
//			*/
//		}
//
//		for _, tc := range tests {
//			test := tc
//			t.Run(test.name, func(tt *testing.T) {
//				tt.Parallel()
//				defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
//				if test.beforeFunc != nil {
//					test.beforeFunc(tt, test.args)
//				}
//				if test.afterFunc != nil {
//					defer test.afterFunc(tt, test.args)
//				}
//				checkFunc := test.checkFunc
//				if test.checkFunc == nil {
//					checkFunc = defaultCheckFunc
//				}
//				g := &gateway{
//					client: test.fields.client,
//					eg:     test.fields.eg,
//				}
//
//				err := g.BroadCast(test.args.ctx, test.args.kind, test.args.f)
//				if err := checkFunc(test.want, err); err != nil {
//					tt.Errorf("error = %v", err)
//				}
//			})
//		}
//	}
//
// func Test_gateway_GetAgentCount(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client discoverer.Client
// 		eg     errgroup.Group
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		fields
// 		want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
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
// 		           client:nil,
// 		           eg:nil,
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
// 		           client:nil,
// 		           eg:nil,
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
// 			g := &gateway{
// 				client: test.fields.client,
// 				eg:     test.fields.eg,
// 			}
//
// 			got := g.GetAgentCount(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_gateway_Addrs(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		client discoverer.Client
// 		eg     errgroup.Group
// 	}
// 	type want struct {
// 		want []string
// 	}
// 	type test struct {
// 		name       string
// 		args
// 		fields
// 		want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string) error {
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
// 		           client:nil,
// 		           eg:nil,
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
// 		           client:nil,
// 		           eg:nil,
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
// 			g := &gateway{
// 				client: test.fields.client,
// 				eg:     test.fields.eg,
// 			}
//
// 			got := g.Addrs(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
