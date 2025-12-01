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
	"context"
	"math"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test"
)

func TestNewHistogram(t *testing.T) {
	type args struct {
		opts []HistogramOption
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (Histogram, error) {
		return NewHistogram(args.opts...)
	}, []test.Case[Histogram, args]{
		{
			Name: "initialize with valid options",
			Args: args{
				opts: []HistogramOption{
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				},
			},
			CheckFunc: func(t *testing.T, want test.Result[Histogram], got test.Result[Histogram]) error {
				t.Helper()
				if got.Err != nil {
					return errors.Errorf("unexpected error: %v", got.Err)
				}
				if got.Val == nil {
					return errors.New("got nil histogram")
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestHistogram_Record_And_Snapshot(t *testing.T) {
	type args struct {
		opts    []HistogramOption
		records []float64
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (*HistogramSnapshot, error) {
		h, err := NewHistogram(args.opts...)
		if err != nil {
			return nil, err
		}
		for _, r := range args.records {
			h.Record(r)
		}
		return h.Snapshot(), nil
	}, []test.Case[*HistogramSnapshot, args]{
		{
			Name: "record values and check snapshot",
			Args: args{
				opts: []HistogramOption{
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				},
				records: []float64{10, 20, 30, 40, 50},
			},
			CheckFunc: func(t *testing.T, want test.Result[*HistogramSnapshot], got test.Result[*HistogramSnapshot]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				if snap.Total != 5 {
					return errors.Errorf("expected total 5, got %d", snap.Total)
				}
				if snap.Mean != 30 {
					return errors.Errorf("expected mean 30, got %f", snap.Mean)
				}
				return nil
			},
		},
		{
			Name: "empty histogram",
			Args: args{
				opts: []HistogramOption{
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				},
			},
			CheckFunc: func(t *testing.T, want test.Result[*HistogramSnapshot], got test.Result[*HistogramSnapshot]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				if snap.Total != 0 {
					return errors.Errorf("expected total 0, got %d", snap.Total)
				}
				if snap.Mean != 0 {
					return errors.Errorf("expected mean 0, got %f", snap.Mean)
				}
				return nil
			},
		},
		{
			Name: "record a single value",
			Args: args{
				opts: []HistogramOption{
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				},
				records: []float64{10},
			},
			CheckFunc: func(t *testing.T, want test.Result[*HistogramSnapshot], got test.Result[*HistogramSnapshot]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				if snap.Total != 1 {
					return errors.Errorf("expected total 1, got %d", snap.Total)
				}
				if snap.Mean != 10 {
					return errors.Errorf("expected mean 10, got %f", snap.Mean)
				}
				return nil
			},
		},
		{
			Name: "record NaN",
			Args: args{
				opts: []HistogramOption{
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				},
				records: []float64{math.NaN(), 10},
			},
			CheckFunc: func(t *testing.T, want test.Result[*HistogramSnapshot], got test.Result[*HistogramSnapshot]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				if snap.Total != 1 { // NaN should be ignored
					return errors.Errorf("expected total 1, got %d", snap.Total)
				}
				if snap.Mean != 10 {
					return errors.Errorf("expected mean 10, got %f", snap.Mean)
				}
				return nil
			},
		},
		{
			Name: "record Inf",
			Args: args{
				opts: []HistogramOption{
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				},
				records: []float64{math.Inf(1), math.Inf(-1), 10},
			},
			CheckFunc: func(t *testing.T, want test.Result[*HistogramSnapshot], got test.Result[*HistogramSnapshot]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				if snap.Total != 1 { // Inf should be ignored
					return errors.Errorf("expected total 1, got %d", snap.Total)
				}
				if snap.Mean != 10 {
					return errors.Errorf("expected mean 10, got %f", snap.Mean)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestHistogram_Merge(t *testing.T) {
	type args struct {
		h1Opts    []HistogramOption
		h1Records []float64
		h2Opts    []HistogramOption
		h2Records []float64
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (*HistogramSnapshot, error) {
		h1, err := NewHistogram(args.h1Opts...)
		if err != nil {
			return nil, err
		}
		for _, r := range args.h1Records {
			h1.Record(r)
		}

		h2, err := NewHistogram(args.h2Opts...)
		if err != nil {
			return nil, err
		}
		for _, r := range args.h2Records {
			h2.Record(r)
		}

		// Merge h1 INTO h2.
		if err := h2.Merge(h1); err != nil {
			return nil, err
		}
		return h2.Snapshot(), nil
	}, []test.Case[*HistogramSnapshot, args]{
		{
			Name: "merge two histograms",
			Args: args{
				h1Opts: []HistogramOption{
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				},
				h1Records: []float64{10, 20},
				h2Opts: []HistogramOption{
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				},
				h2Records: []float64{30, 40},
			},
			CheckFunc: func(t *testing.T, want test.Result[*HistogramSnapshot], got test.Result[*HistogramSnapshot]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				if snap.Total != 4 {
					return errors.Errorf("expected total 4, got %d", snap.Total)
				}
				if snap.Mean != 25 {
					return errors.Errorf("expected mean 25, got %f", snap.Mean)
				}
				return nil
			},
		},
		{
			Name: "merge incompatible histograms (shard count)",
			Args: args{
				h1Opts: []HistogramOption{
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(2),
				},
				h1Records: []float64{10, 20},
				h2Opts: []HistogramOption{
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(3), // Different shard count
				},
				h2Records: []float64{30, 40},
			},
			CheckFunc: func(t *testing.T, want test.Result[*HistogramSnapshot], got test.Result[*HistogramSnapshot]) error {
				t.Helper()
				if got.Err == nil {
					return errors.New("expected error, got nil")
				}
				if got.Err.Error() != "incompatible shards: count mismatch" {
					return errors.Errorf("unexpected error message: %v", got.Err)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestHistogram_Concurrent(t *testing.T) {
	type args struct {
		count int
	}
	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (*HistogramSnapshot, error) {
		h, err := NewHistogram(
			WithHistogramMin(1),
			WithHistogramGrowth(2),
			WithHistogramNumBuckets(10),
			WithHistogramNumShards(4), // multiple shards for concurrency test
		)
		if err != nil {
			return nil, err
		}
		var wg sync.WaitGroup
		for i := 0; i < args.count; i++ {
			wg.Add(1)
			go func(v float64) {
				defer wg.Done()
				h.Record(v)
			}(float64(i))
		}
		wg.Wait()
		return h.Snapshot(), nil
	}, []test.Case[*HistogramSnapshot, args]{
		{
			Name: "concurrent record",
			Args: args{count: 100},
			CheckFunc: func(t *testing.T, want test.Result[*HistogramSnapshot], got test.Result[*HistogramSnapshot]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				if snap.Total != 100 {
					return errors.Errorf("expected total 100, got %d", snap.Total)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestHistogram_Clone(t *testing.T) {
	// Test clone
	type args struct {
		opts    []HistogramOption
		records []float64
	}
	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (*HistogramSnapshot, error) {
		h, err := NewHistogram(args.opts...)
		if err != nil {
			return nil, err
		}
		for _, r := range args.records {
			h.Record(r)
		}
		cloned := h.Clone()
		// Modify original to ensure deep copy
		h.Record(1000)
		return cloned.Snapshot(), nil
	}, []test.Case[*HistogramSnapshot, args]{
		{
			Name: "clone deep copy",
			Args: args{
				opts: []HistogramOption{
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				},
				records: []float64{10, 20},
			},
			CheckFunc: func(t *testing.T, want test.Result[*HistogramSnapshot], got test.Result[*HistogramSnapshot]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				if snap.Total != 2 {
					return errors.Errorf("expected total 2, got %d", snap.Total)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func BenchmarkHistogram_Record(b *testing.B) {
	h, _ := NewHistogram(
		WithHistogramMin(1),
		WithHistogramGrowth(1.1),
		WithHistogramNumBuckets(100),
		WithHistogramNumShards(16),
	)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0.0
		for pb.Next() {
			h.Record(i)
			i += 1.0
		}
	})
}

func BenchmarkHistogram_Snapshot(b *testing.B) {
	h, _ := NewHistogram(
		WithHistogramMin(1),
		WithHistogramGrowth(1.1),
		WithHistogramNumBuckets(100),
		WithHistogramNumShards(16),
	)
	for i := range 10000 {
		h.Record(float64(i))
	}
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = h.Snapshot()
	}
}

func TestHistogram_StdDevAccuracy(t *testing.T) {
	// Reproduce StdDev = 0 issue
	h, err := NewHistogram(WithHistogramNumBuckets(10))
	if err != nil {
		t.Fatal(err)
	}

	// Record some values (using Nanoseconds unit: 10ms, 20ms, ...)
	// 10ms = 10,000,000ns
	values := []float64{1e7, 2e7, 3e7, 4e7, 5e7}
	for _, v := range values {
		h.Record(v)
	}

	snap := h.Snapshot()
	t.Logf("Snapshot: %+v", snap)

	// Mean is 30ms = 3e7
	// Variance = ((10-30)^2 + ... ) / 5 = 200 * (1e6)^2 = 2e14
	// StdDev = sqrt(2e14) = 1.4142e7 = 14.14ms

	if snap.Mean != 3e7 {
		t.Errorf("Expected Mean 3e7, got %f", snap.Mean)
	}

	expectedStdDev := 1.41421356e7
	if math.Abs(snap.StdDev-expectedStdDev) > 1e2 { // tolerance 100ns
		t.Errorf("Expected StdDev ~%f, got %f", expectedStdDev, snap.StdDev)
	}
}

func TestHistogram_BucketsValidation(t *testing.T) {
	// Reproduce "All buckets 0" (likely all in last bucket due to unit mismatch)
	// Default options: min=1, growth=1.2, buckets=50.
	// Latency recorded in Nanoseconds. 100ms = 1e8 ns.

	c, err := NewCollector() // uses defaults
	if err != nil {
		t.Fatal(err)
	}

	rr := &RequestResult{
		Latency: 100 * time.Millisecond,
	}
	c.Record(context.Background(), rr)

	snap := c.GlobalSnapshot()
	lat := snap.Latencies

	// Check bucket counts
	nonZeroBuckets := 0
	lastBucketIdx := len(lat.Counts) - 1
	for i, count := range lat.Counts {
		if count > 0 {
			nonZeroBuckets++
			t.Logf("Bucket %d has count %d (Bound: %f)", i, count, lat.Bounds[min(i, len(lat.Bounds)-1)])
		}
	}

	if nonZeroBuckets == 0 {
		t.Error("No buckets have data")
	} else if lat.Counts[lastBucketIdx] > 0 {
		t.Log("Data is in the last bucket, confirming unit mismatch (nanoseconds vs small bounds)")
	} else {
		t.Log("Data is distributed properly?")
	}
}

func TestHistogramBucketsConfiguration(t *testing.T) {
	t.Run("Default buckets should be used when no options provided", func(t *testing.T) {
		// NewHistogram uses defaultHistogramOpts which we modified to use WithHistogramBuckets
		// Using 2 shards to ensure shardedHistogram
		h, err := NewHistogram(WithHistogramNumShards(2))
		if err != nil {
			t.Fatalf("NewHistogram failed: %v", err)
		}

		sh, ok := h.(*shardedHistogram)
		if !ok {
			t.Fatalf("Expected shardedHistogram, got %T", h)
		}
		hist := sh.shards[0]

		// Check if bounds match defaultHistogramBuckets
		if len(hist.bounds) != len(defaultHistogramBuckets) {
			t.Errorf("Expected %d bounds, got %d", len(defaultHistogramBuckets), len(hist.bounds))
		}

		// Check a few values
		if hist.bounds[0] != defaultHistogramBuckets[0] {
			t.Errorf("Expected first bound %v, got %v", defaultHistogramBuckets[0], hist.bounds[0])
		}

		// Verify bucket finder is binary search
		// 5e6 is 5ms.
		// Record 4ms -> bucket 0 (<= 5ms)
		// Record 6ms -> bucket 1 (5ms - 10ms]

		hist.Record(4e6)
		hist.Record(6e6)

		snap := h.Snapshot()
		if snap.Counts[0] != 1 {
			t.Errorf("Expected count 1 in bucket 0, got %d", snap.Counts[0])
		}
		if snap.Counts[1] != 1 {
			t.Errorf("Expected count 1 in bucket 1, got %d", snap.Counts[1])
		}
	})

	t.Run("Explicit buckets should be used when provided", func(t *testing.T) {
		customBuckets := []float64{10, 20, 30}
		h, err := NewHistogram(
			WithHistogramBuckets(customBuckets),
			WithHistogramNumShards(2),
		)
		if err != nil {
			t.Fatalf("NewHistogram failed: %v", err)
		}

		sh, ok := h.(*shardedHistogram)
		if !ok {
			t.Fatalf("Expected shardedHistogram, got %T", h)
		}
		hist := sh.shards[0]

		if len(hist.bounds) != 3 {
			t.Errorf("Expected 3 bounds, got %d", len(hist.bounds))
		}
		if hist.bounds[0] != 10 {
			t.Errorf("Expected first bound 10, got %v", hist.bounds[0])
		}

		// Test recording
		// <= 10 -> bucket 0
		// 15 -> bucket 1
		// 25 -> bucket 2
		// 35 -> bucket 3 (last)

		hist.Record(5)
		hist.Record(15)
		hist.Record(25)
		hist.Record(35)

		snap := h.Snapshot()
		if snap.Counts[0] != 1 {
			t.Errorf("Bucket 0 count mismatch")
		}
		if snap.Counts[1] != 1 {
			t.Errorf("Bucket 1 count mismatch")
		}
		if snap.Counts[2] != 1 {
			t.Errorf("Bucket 2 count mismatch")
		}
		if snap.Counts[3] != 1 {
			t.Errorf("Bucket 3 count mismatch")
		}
	})

	t.Run("Geometric buckets should be used when Min/Growth provided", func(t *testing.T) {
		// Providing WithHistogramMin should clear explicit buckets
		h, err := NewHistogram(
			WithHistogramMin(10),
			WithHistogramGrowth(2),
			WithHistogramNumBuckets(5),
			WithHistogramNumShards(2),
		)
		if err != nil {
			t.Fatalf("NewHistogram failed: %v", err)
		}

		sh, ok := h.(*shardedHistogram)
		if !ok {
			t.Fatalf("Expected shardedHistogram, got %T", h)
		}
		hist := sh.shards[0]

		// Bounds: 10, 20, 40, 80
		// numBuckets: 5
		// len(bounds): 4

		if len(hist.bounds) != 4 {
			t.Errorf("Expected 4 bounds, got %d", len(hist.bounds))
		}

		if hist.bounds[0] != 10 {
			t.Errorf("Expected bound[0] 10, got %v", hist.bounds[0])
		}
		if hist.bounds[1] != 20 {
			t.Errorf("Expected bound[1] 20, got %v", hist.bounds[1])
		}

		// Verify geometric finder logic
		// Record 15 -> 10 * 2^0=10, 10*2^1=20. 15 is in (10, 20]. Bucket 1.
		// Record 5 -> <= 10. Bucket 0.

		hist.Record(5)
		hist.Record(15)

		snap := h.Snapshot()
		if snap.Counts[0] != 1 {
			t.Errorf("Bucket 0 count mismatch")
		}
		if snap.Counts[1] != 1 {
			t.Errorf("Bucket 1 count mismatch")
		}
	})

	t.Run("Mixed configuration priority", func(t *testing.T) {
		// Buckets provided AFTER Min should win
		customBuckets := []float64{100, 200}
		h, err := NewHistogram(
			WithHistogramMin(10),
			WithHistogramBuckets(customBuckets),
			WithHistogramNumShards(2),
		)
		if err != nil {
			t.Fatalf("NewHistogram failed: %v", err)
		}
		sh, ok := h.(*shardedHistogram)
		if !ok {
			t.Fatalf("Expected shardedHistogram, got %T", h)
		}
		hist := sh.shards[0]
		if len(hist.bounds) != 2 {
			t.Errorf("Expected explicit buckets to win, got %d bounds", len(hist.bounds))
		}
		if hist.bounds[0] != 100 {
			t.Errorf("Expected first bound 100")
		}

		// Min provided AFTER Buckets should win (and use geometric)
		h2, err := NewHistogram(
			WithHistogramBuckets(customBuckets),
			WithHistogramMin(10),
			WithHistogramGrowth(2),
			WithHistogramNumBuckets(4),
			WithHistogramNumShards(2),
		)
		if err != nil {
			t.Fatalf("NewHistogram failed: %v", err)
		}
		sh2, ok := h2.(*shardedHistogram)
		if !ok {
			t.Fatalf("Expected shardedHistogram, got %T", h2)
		}
		hist2 := sh2.shards[0]
		if len(hist2.bounds) != 3 { // NumBuckets 4 -> 3 bounds
			t.Errorf("Expected geometric to win, got %d bounds", len(hist2.bounds))
		}
		if hist2.bounds[0] != 10 {
			t.Errorf("Expected first bound 10")
		}
	})
}
