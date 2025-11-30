//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package metrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test"
)

func TestNewExemplar(t *testing.T) {
	if err := test.Run(t.Context(), t, func(t *testing.T, opts []ExemplarOption) (Exemplar, error) {
		t.Helper()
		return NewExemplar(opts...), nil
	}, []test.Case[Exemplar, []ExemplarOption]{
		{
			Name: "initialized with default options",
			Args: nil,
			CheckFunc: func(t *testing.T, want test.Result[Exemplar], got test.Result[Exemplar]) error {
				t.Helper()
				if got.Val == nil {
					return errors.New("got nil exemplar")
				}
				// Check interface implementation indirectly via Snapshot
				snap := got.Val.Snapshot()
				if snap == nil {
					return errors.New("got nil snapshot")
				}

				// Optional: Check type if we care about it being sharded by default
				if _, ok := got.Val.(*shardedExemplar); !ok {
					return errors.New("expected default exemplar to be sharded")
				}
				return nil
			},
		},
		{
			Name: "initialized with capacity",
			Args: []ExemplarOption{
				WithExemplarCapacity(10),
				WithExemplarNumShards(1), // Force single shard for capacity check simplicity
			},
			CheckFunc: func(t *testing.T, want test.Result[Exemplar], got test.Result[Exemplar]) error {
				t.Helper()
				if got.Val == nil {
					return errors.New("got nil exemplar")
				}
				e, ok := got.Val.(*exemplar)
				if !ok {
					return errors.New("expected single exemplar with shards=1")
				}
				if e.k != 10 {
					return errors.Errorf("expected k=10, got %d", e.k)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestExemplar_Offer(t *testing.T) {
	type offer struct {
		id      string
		latency time.Duration
	}
	type args struct {
		opts   []ExemplarOption
		offers []offer
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, args args) ([]*ExemplarItem, error) {
		t.Helper()
		e := NewExemplar(args.opts...)
		for _, o := range args.offers {
			e.Offer(o.latency, o.id, nil, "")
		}
		return e.Snapshot(), nil
	}, []test.Case[[]*ExemplarItem, args]{
		{
			Name: "offer requests and check snapshot",
			Args: args{
				// Use 1 shard to ensure deterministic "Top K" behavior when counts match capacity exactly
				opts: []ExemplarOption{WithExemplarCapacity(3), WithExemplarNumShards(1)},
				offers: []offer{
					{id: "req-1", latency: 100 * time.Millisecond},
					{id: "req-2", latency: 200 * time.Millisecond},
					{id: "req-3", latency: 50 * time.Millisecond},
					{id: "req-4", latency: 300 * time.Millisecond},
				},
			},
			CheckFunc: func(t *testing.T, want test.Result[[]*ExemplarItem], got test.Result[[]*ExemplarItem]) error {
				t.Helper()
				snap := got.Val
				if len(snap) != 3 {
					return errors.Errorf("expected snapshot length 3, got %d", len(snap))
				}
				expectedIDs := []string{"req-4", "req-2", "req-1"}
				for i, id := range expectedIDs {
					if snap[i].RequestID != id {
						return errors.Errorf("expected snapshot[%d] ID %s, got %s", i, id, snap[i].RequestID)
					}
				}
				return nil
			},
		},
		{
			Name: "offer requests with same latency",
			Args: args{
				opts: []ExemplarOption{WithExemplarCapacity(3)},
				offers: []offer{
					{id: "req-1", latency: 100 * time.Millisecond},
					{id: "req-2", latency: 200 * time.Millisecond},
					{id: "req-3", latency: 100 * time.Millisecond},
					{id: "req-4", latency: 300 * time.Millisecond},
				},
			},
			CheckFunc: func(t *testing.T, want test.Result[[]*ExemplarItem], got test.Result[[]*ExemplarItem]) error {
				snap := got.Val
				if len(snap) != 3 {
					return errors.Errorf("expected snapshot length 3, got %d", len(snap))
				}
				return nil
			},
		},
		{
			Name: "empty exemplar",
			Args: args{
				opts: []ExemplarOption{WithExemplarCapacity(3)},
			},
			CheckFunc: func(t *testing.T, want test.Result[[]*ExemplarItem], got test.Result[[]*ExemplarItem]) error {
				snap := got.Val
				if len(snap) != 0 {
					return errors.Errorf("expected snapshot length 0, got %d", len(snap))
				}
				return nil
			},
		},
		{
			Name: "snapshot is sorted by latency",
			Args: args{
				opts: []ExemplarOption{WithExemplarCapacity(3), WithExemplarNumShards(1)},
				offers: []offer{
					{id: "req-2", latency: 200 * time.Millisecond},
					{id: "req-1", latency: 100 * time.Millisecond},
					{id: "req-3", latency: 300 * time.Millisecond},
				},
			},
			CheckFunc: func(t *testing.T, want test.Result[[]*ExemplarItem], got test.Result[[]*ExemplarItem]) error {
				snap := got.Val
				if len(snap) != 3 {
					return errors.Errorf("expected snapshot length 3, got %d", len(snap))
				}
				expectedIDs := []string{"req-3", "req-2", "req-1"}
				for i, id := range expectedIDs {
					if snap[i].RequestID != id {
						return errors.Errorf("expected snapshot[%d] ID %s, got %s", i, id, snap[i].RequestID)
					}
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestExemplar_Reset(t *testing.T) {
	type args struct {
		opts   []ExemplarOption
		offers []struct {
			id      string
			latency time.Duration
		}
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (Exemplar, error) {
		e := NewExemplar(args.opts...)
		for _, o := range args.offers {
			e.Offer(o.latency, o.id, nil, "")
		}
		e.Reset()
		return e, nil
	}, []test.Case[Exemplar, args]{
		{
			Name: "reset clears all data",
			Args: args{
				opts: []ExemplarOption{WithExemplarCapacity(3)},
				offers: []struct {
					id      string
					latency time.Duration
				}{
					{id: "req-1", latency: 100 * time.Millisecond},
					{id: "req-2", latency: 200 * time.Millisecond},
				},
			},
			CheckFunc: func(t *testing.T, want test.Result[Exemplar], got test.Result[Exemplar]) error {
				if got.Val == nil {
					return errors.New("got nil exemplar")
				}
				snap := got.Val.Snapshot()
				if len(snap) != 0 {
					return errors.Errorf("expected snapshot length 0 after reset, got %d", len(snap))
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestExemplar_Clone(t *testing.T) {
	type args struct {
		opts   []ExemplarOption
		offers []struct {
			id      string
			latency time.Duration
		}
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (Exemplar, error) {
		e := NewExemplar(args.opts...)
		for _, o := range args.offers {
			e.Offer(o.latency, o.id, nil, "")
		}
		return e.Clone(), nil
	}, []test.Case[Exemplar, args]{
		{
			Name: "clone copies data",
			Args: args{
				opts: []ExemplarOption{WithExemplarCapacity(3)},
				offers: []struct {
					id      string
					latency time.Duration
				}{
					{id: "req-1", latency: 100 * time.Millisecond},
					{id: "req-2", latency: 200 * time.Millisecond},
				},
			},
			CheckFunc: func(t *testing.T, want test.Result[Exemplar], got test.Result[Exemplar]) error {
				if got.Val == nil {
					return errors.New("got nil exemplar")
				}
				snap := got.Val.Snapshot()
				if len(snap) != 2 {
					return errors.Errorf("expected snapshot length 2, got %d", len(snap))
				}
				// Verify content
				hasReq1 := false
				hasReq2 := false
				for _, item := range snap {
					if item.RequestID == "req-1" {
						hasReq1 = true
					}
					if item.RequestID == "req-2" {
						hasReq2 = true
					}
				}
				if !hasReq1 || !hasReq2 {
					return errors.New("snapshot missing expected items")
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestExemplar_Concurrent(t *testing.T) {
	type args struct {
		capacity          int
		workers           int
		requestsPerWorker int
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (Exemplar, error) {
		e := NewExemplar(WithExemplarCapacity(args.capacity))
		var wg sync.WaitGroup
		for i := 0; i < args.workers; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < args.requestsPerWorker; j++ {
					e.Offer(time.Duration(j)*time.Millisecond, fmt.Sprintf("req-%d-%d", i, j), nil, "")
				}
			}(i)
		}
		wg.Wait()
		return e, nil
	}, []test.Case[Exemplar, args]{
		{
			Name: "concurrent offers fill capacity",
			Args: args{
				capacity:          10,
				workers:           10,
				requestsPerWorker: 100,
			},
			CheckFunc: func(t *testing.T, want test.Result[Exemplar], got test.Result[Exemplar]) error {
				snap := got.Val.Snapshot()
				if len(snap) != 10 {
					return errors.Errorf("expected 10 exemplars, got %d", len(snap))
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestExemplar_Race(t *testing.T) {
	type args struct {
		capacity          int
		workers           int
		requestsPerWorker int
	}
	// This test verifies no race conditions when Offer and Snapshot are called concurrently.
	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (Exemplar, error) {
		e := NewExemplar(WithExemplarCapacity(args.capacity))
		eg, _ := errgroup.New(t.Context())
		for i := 0; i < args.workers; i++ {
			eg.Go(func() error {
				for j := 0; j < args.requestsPerWorker; j++ {
					e.Offer(time.Duration(j)*time.Millisecond, "req", nil, "")
					e.Snapshot()
				}
				return nil
			})
		}
		return e, eg.Wait()
	}, []test.Case[Exemplar, args]{
		{
			Name: "race detection",
			Args: args{
				capacity:          10,
				workers:           10,
				requestsPerWorker: 100,
			},
			CheckFunc: func(t *testing.T, want test.Result[Exemplar], got test.Result[Exemplar]) error {
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func BenchmarkExemplar_Offer(b *testing.B) {
	e := NewExemplar(WithExemplarCapacity(100))
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			e.Offer(time.Duration(i)*time.Millisecond, "req", nil, "")
			i++
		}
	})
}

func TestExemplar_Categories(t *testing.T) {
	type args struct {
		opts []ExemplarOption
	}
	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (*ExemplarDetails, error) {
		e := NewExemplar(args.opts...)

		// Offer 5 items with varying latencies
		// 10ms, 50ms, 30ms, 90ms, 20ms
		e.Offer(10*time.Millisecond, "req-10", nil, "")
		e.Offer(50*time.Millisecond, "req-50", nil, "")
		e.Offer(30*time.Millisecond, "req-30", nil, "")
		e.Offer(90*time.Millisecond, "req-90", nil, "")
		e.Offer(20*time.Millisecond, "req-20", nil, "")

		// Offer failures
		e.Offer(100*time.Millisecond, "fail-100", errors.New("failed"), "failed")
		e.Offer(40*time.Millisecond, "fail-40", errors.New("failed"), "failed")

		return e.DetailedSnapshot()
	}, []test.Case[*ExemplarDetails, args]{
		{
			Name: "check categories with k=3",
			Args: args{
				// Use 1 shard to simplify verification of Exact K behavior
				opts: []ExemplarOption{WithExemplarCapacity(3), WithExemplarNumShards(1)},
			},
			CheckFunc: func(t *testing.T, want test.Result[*ExemplarDetails], got test.Result[*ExemplarDetails]) error {
				d := got.Val

				// Slowest (Top 3 Max): 100 (fail), 90, 50
				if len(d.Slowest) != 3 {
					return errors.Errorf("expected 3 slowest, got %d", len(d.Slowest))
				}
				if d.Slowest[0].Latency != 100*time.Millisecond {
					return errors.Errorf("expected slowest[0] 100ms, got %v", d.Slowest[0].Latency)
				}
				if d.Slowest[2].Latency != 50*time.Millisecond {
					return errors.Errorf("expected slowest[2] 50ms, got %v", d.Slowest[2].Latency)
				}

				// Fastest (Top 3 Min): 10, 20, 30
				if len(d.Fastest) != 3 {
					return errors.Errorf("expected 3 fastest, got %d", len(d.Fastest))
				}
				if d.Fastest[0].Latency != 10*time.Millisecond {
					return errors.Errorf("expected fastest[0] 10ms, got %v", d.Fastest[0].Latency)
				}
				if d.Fastest[2].Latency != 30*time.Millisecond {
					return errors.Errorf("expected fastest[2] 30ms, got %v", d.Fastest[2].Latency)
				}

				// Failures: fail-100, fail-40
				if len(d.Failures) != 2 {
					return errors.Errorf("expected 2 failures, got %d", len(d.Failures))
				}
				// Sorted desc
				if d.Failures[0].Latency != 100*time.Millisecond {
					return errors.Errorf("expected failure[0] 100ms, got %v", d.Failures[0].Latency)
				}

				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestExemplar_ProbabilisticSampling(t *testing.T) {
	// This test verifies that even with probabilistic sampling, we eventually collect average samples,
	// and that edge cases (Slowest/Fastest) are always recorded accurately.
	type args struct {
		opts []ExemplarOption
	}
	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (*ExemplarDetails, error) {
		e := NewExemplar(args.opts...)

		// 1. Offer distinct Slowest/Fastest candidates.
		// Fastest: 1ms
		e.Offer(1*time.Millisecond, "fastest", nil, "")
		// Slowest: 10s
		e.Offer(10*time.Second, "slowest", nil, "")
		// Error: should always be recorded
		e.Offer(500*time.Millisecond, "error-req", errors.New("some error"), "err msg")

		// 2. Offer many "average" requests.
		// These should be filtered by the sampling (1/16).
		// Latencies between 10ms and 1s.
		// We offer enough to statistically guarantee hitting the reservoir at least once.
		for i := 0; i < 1000; i++ {
			e.Offer(100*time.Millisecond, fmt.Sprintf("avg-%d", i), nil, "")
		}

		return e.DetailedSnapshot()
	}, []test.Case[*ExemplarDetails, args]{
		{
			Name: "ensures probabilistic sampling retains critical samples and some average samples",
			Args: args{
				opts: []ExemplarOption{WithExemplarCapacity(10), WithExemplarNumShards(1)},
			},
			CheckFunc: func(t *testing.T, want test.Result[*ExemplarDetails], got test.Result[*ExemplarDetails]) error {
				d := got.Val
				if d == nil {
					return errors.New("got nil details")
				}

				// Check Fastest
				if len(d.Fastest) == 0 {
					return errors.New("fastest sample missing")
				}
				if d.Fastest[0].Latency != 1*time.Millisecond {
					return errors.Errorf("expected fastest 1ms, got %v", d.Fastest[0].Latency)
				}

				// Check Slowest
				if len(d.Slowest) == 0 {
					return errors.New("slowest sample missing")
				}
				if d.Slowest[0].Latency != 10*time.Second {
					return errors.Errorf("expected slowest 10s, got %v", d.Slowest[0].Latency)
				}

				// Check Failures
				if len(d.Failures) == 0 {
					return errors.New("failure sample missing")
				}
				if d.Failures[0].RequestID != "error-req" {
					return errors.Errorf("expected failure req-id 'error-req', got %s", d.Failures[0].RequestID)
				}

				// Check Average
				// With 1000 samples and 1/16 rate, we expect ~62 samples.
				// Capacity is 10. So reservoir should be full (10 items).
				if len(d.Average) == 0 {
					return errors.New("average samples missing (probabilistic failure?)")
				}

				// Verify that we captured at least some "average" samples (IDs starting with "avg-")
				foundAvg := false
				for _, item := range d.Average {
					// The Average reservoir may contain Fastest/Slowest samples too because
					// they fall through the fast-path check and update all categories.
					// We just want to ensure we didn't miss the bulk of "avg-*" requests.
					if len(item.RequestID) >= 4 && item.RequestID[:4] == "avg-" {
						foundAvg = true
						break
					}
				}
				if !foundAvg {
					return errors.New("no 'avg-*' samples found in average reservoir")
				}

				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}
