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
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	testdata "github.com/vdaas/vald/internal/test"
)

func TestNewExemplar(t *testing.T) {
	type args struct {
		opts []ExemplarOption
	}
	type want struct {
		k int
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Exemplar) error
		beforeFunc func(*testing.T, args) args
		afterFunc  func(*testing.T, args)
	}

	if err := testdata.Run(t.Context(), t, func(tt *testing.T, args args) (Exemplar, error) {
		return NewExemplar(args.opts...), nil
	}, []testdata.Case[Exemplar, args]{
		{
			Name: "initialized with default options",
			Args: args{
				opts: nil,
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[Exemplar], got testdata.Result[Exemplar]) error {
				if got.Val == nil {
					return errors.New("got nil exemplar")
				}
				e := got.Val.(*exemplar)
				if e.k != 1 {
					return errors.Errorf("expected k=1, got %d", e.k)
				}
				return nil
			},
		},
		{
			Name: "initialized with capacity",
			Args: args{
				opts: []ExemplarOption{
					WithExemplarCapacity(10),
				},
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[Exemplar], got testdata.Result[Exemplar]) error {
				if got.Val == nil {
					return errors.New("got nil exemplar")
				}
				e := got.Val.(*exemplar)
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
		latency time.Duration
		id      string
	}
	type args struct {
		opts   []ExemplarOption
		offers []offer
	}

	if err := testdata.Run(t.Context(), t, func(tt *testing.T, args args) ([]*item, error) {
		e := NewExemplar(args.opts...)
		for _, o := range args.offers {
			e.Offer(o.latency, o.id, false)
		}
		return e.Snapshot(), nil
	}, []testdata.Case[[]*item, args]{
		{
			Name: "offer requests and check snapshot",
			Args: args{
				opts: []ExemplarOption{WithExemplarCapacity(3)},
				offers: []offer{
					{100 * time.Millisecond, "req-1"},
					{200 * time.Millisecond, "req-2"},
					{50 * time.Millisecond, "req-3"},
					{300 * time.Millisecond, "req-4"},
				},
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[[]*item], got testdata.Result[[]*item]) error {
				snap := got.Val
				if len(snap) != 3 {
					return errors.Errorf("expected snapshot length 3, got %d", len(snap))
				}
				expectedIDs := []string{"req-4", "req-2", "req-1"}
				for i, id := range expectedIDs {
					if snap[i].requestID != id {
						return errors.Errorf("expected snapshot[%d] ID %s, got %s", i, id, snap[i].requestID)
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
					{100 * time.Millisecond, "req-1"},
					{200 * time.Millisecond, "req-2"},
					{100 * time.Millisecond, "req-3"},
					{300 * time.Millisecond, "req-4"},
				},
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[[]*item], got testdata.Result[[]*item]) error {
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
			CheckFunc: func(tt *testing.T, want testdata.Result[[]*item], got testdata.Result[[]*item]) error {
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
				opts: []ExemplarOption{WithExemplarCapacity(3)},
				offers: []offer{
					{200 * time.Millisecond, "req-2"},
					{100 * time.Millisecond, "req-1"},
					{300 * time.Millisecond, "req-3"},
				},
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[[]*item], got testdata.Result[[]*item]) error {
				snap := got.Val
				if len(snap) != 3 {
					return errors.Errorf("expected snapshot length 3, got %d", len(snap))
				}
				expectedIDs := []string{"req-3", "req-2", "req-1"}
				for i, id := range expectedIDs {
					if snap[i].requestID != id {
						return errors.Errorf("expected snapshot[%d] ID %s, got %s", i, id, snap[i].requestID)
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
			latency time.Duration
			id      string
		}
	}

	if err := testdata.Run(t.Context(), t, func(tt *testing.T, args args) (Exemplar, error) {
		e := NewExemplar(args.opts...)
		for _, o := range args.offers {
			e.Offer(o.latency, o.id, false)
		}
		e.Reset()
		return e, nil
	}, []testdata.Case[Exemplar, args]{
		{
			Name: "reset clears all data",
			Args: args{
				opts: []ExemplarOption{WithExemplarCapacity(3)},
				offers: []struct {
					latency time.Duration
					id      string
				}{
					{100 * time.Millisecond, "req-1"},
					{200 * time.Millisecond, "req-2"},
				},
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[Exemplar], got testdata.Result[Exemplar]) error {
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
			latency time.Duration
			id      string
		}
	}

	if err := testdata.Run(t.Context(), t, func(tt *testing.T, args args) (Exemplar, error) {
		e := NewExemplar(args.opts...)
		for _, o := range args.offers {
			e.Offer(o.latency, o.id, false)
		}
		return e.Clone(), nil
	}, []testdata.Case[Exemplar, args]{
		{
			Name: "clone copies data",
			Args: args{
				opts: []ExemplarOption{WithExemplarCapacity(3)},
				offers: []struct {
					latency time.Duration
					id      string
				}{
					{100 * time.Millisecond, "req-1"},
					{200 * time.Millisecond, "req-2"},
				},
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[Exemplar], got testdata.Result[Exemplar]) error {
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
					if item.requestID == "req-1" {
						hasReq1 = true
					}
					if item.requestID == "req-2" {
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

	if err := testdata.Run(t.Context(), t, func(tt *testing.T, args args) (Exemplar, error) {
		e := NewExemplar(WithExemplarCapacity(args.capacity))
		var wg sync.WaitGroup
		for i := 0; i < args.workers; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < args.requestsPerWorker; j++ {
					e.Offer(time.Duration(j)*time.Millisecond, fmt.Sprintf("req-%d-%d", i, j), false)
				}
			}(i)
		}
		wg.Wait()
		return e, nil
	}, []testdata.Case[Exemplar, args]{
		{
			Name: "concurrent offers fill capacity",
			Args: args{
				capacity:          10,
				workers:           10,
				requestsPerWorker: 100,
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[Exemplar], got testdata.Result[Exemplar]) error {
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
	if err := testdata.Run(t.Context(), t, func(tt *testing.T, args args) (Exemplar, error) {
		e := NewExemplar(WithExemplarCapacity(args.capacity))
		var wg sync.WaitGroup
		for i := 0; i < args.workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < args.requestsPerWorker; j++ {
					e.Offer(time.Duration(j)*time.Millisecond, "req", false)
					e.Snapshot()
				}
			}()
		}
		wg.Wait()
		return e, nil
	}, []testdata.Case[Exemplar, args]{
		{
			Name: "race detection",
			Args: args{
				capacity:          10,
				workers:           10,
				requestsPerWorker: 100,
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[Exemplar], got testdata.Result[Exemplar]) error {
				// Just completion is enough to prove no panic/race (race detector needed).
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
			e.Offer(time.Duration(i)*time.Millisecond, "req", false)
			i++
		}
	})
}

func TestExemplar_Categories(t *testing.T) {
	type args struct {
		opts []ExemplarOption
	}
	if err := testdata.Run(t.Context(), t, func(tt *testing.T, args args) (*ExemplarDetails, error) {
		e := NewExemplar(args.opts...)

		// Offer 5 items with varying latencies
		// 10ms, 50ms, 30ms, 90ms, 20ms
		e.Offer(10*time.Millisecond, "req-10", false)
		e.Offer(50*time.Millisecond, "req-50", false)
		e.Offer(30*time.Millisecond, "req-30", false)
		e.Offer(90*time.Millisecond, "req-90", false)
		e.Offer(20*time.Millisecond, "req-20", false)

		// Offer failures
		e.Offer(100*time.Millisecond, "fail-100", true)
		e.Offer(40*time.Millisecond, "fail-40", true)

		return e.DetailedSnapshot()
	}, []testdata.Case[*ExemplarDetails, args]{
		{
			Name: "check categories with k=3",
			Args: args{
				opts: []ExemplarOption{WithExemplarCapacity(3)},
			},
			CheckFunc: func(tt *testing.T, want testdata.Result[*ExemplarDetails], got testdata.Result[*ExemplarDetails]) error {
				d := got.Val

				// Slowest (Top 3 Max): 100 (fail), 90, 50
				if len(d.Slowest) != 3 {
					return errors.Errorf("expected 3 slowest, got %d", len(d.Slowest))
				}
				if d.Slowest[0].latency != 100*time.Millisecond {
					return errors.Errorf("expected slowest[0] 100ms, got %v", d.Slowest[0].latency)
				}
				if d.Slowest[2].latency != 50*time.Millisecond {
					return errors.Errorf("expected slowest[2] 50ms, got %v", d.Slowest[2].latency)
				}

				// Fastest (Top 3 Min): 10, 20, 30
				if len(d.Fastest) != 3 {
					return errors.Errorf("expected 3 fastest, got %d", len(d.Fastest))
				}
				if d.Fastest[0].latency != 10*time.Millisecond {
					return errors.Errorf("expected fastest[0] 10ms, got %v", d.Fastest[0].latency)
				}
				if d.Fastest[2].latency != 30*time.Millisecond {
					return errors.Errorf("expected fastest[2] 30ms, got %v", d.Fastest[2].latency)
				}

				// Failures: fail-100, fail-40
				if len(d.Failures) != 2 {
					return errors.Errorf("expected 2 failures, got %d", len(d.Failures))
				}
				// Sorted desc
				if d.Failures[0].latency != 100*time.Millisecond {
					return errors.Errorf("expected failure[0] 100ms, got %v", d.Failures[0].latency)
				}

				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}
