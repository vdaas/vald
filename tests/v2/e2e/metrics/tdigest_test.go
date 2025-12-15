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
	"encoding/json"
	"math"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test"
)

func TestNewTDigest(t *testing.T) {
	if err := test.Run(t.Context(), t, func(tt *testing.T, opts []TDigestOption) (TDigest, error) {
		return NewTDigest(opts...)
	}, []test.Case[TDigest, []TDigestOption]{
		{
			Name: "initialize with valid options",
			Args: []TDigestOption{
				WithTDigestCompression(100),
				WithTDigestCompressionTriggerFactor(10),
			},
			CheckFunc: func(tt *testing.T, want test.Result[TDigest], got test.Result[TDigest]) error {
				if got.Err != nil {
					return got.Err
				}
				if got.Val == nil {
					return errors.New("got nil TDigest")
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestTDigest_Add_And_Quantile(t *testing.T) {
	type args struct {
		opts   []TDigestOption
		values []float64
		q      float64
	}

	if err := test.Run(t.Context(), t, func(tt *testing.T, args args) (float64, error) {
		td, err := NewTDigest(args.opts...)
		if err != nil {
			return 0, err
		}
		for _, v := range args.values {
			td.Add(v)
		}
		return td.Quantile(args.q), nil
	}, []test.Case[float64, args]{
		{
			Name: "quantile 0.5 check",
			Args: args{
				opts:   []TDigestOption{WithTDigestCompression(100)},
				values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				q:      0.5,
			},
			CheckFunc: func(tt *testing.T, want test.Result[float64], got test.Result[float64]) error {
				if got.Val < 5 || got.Val > 6 {
					return errors.Errorf("Quantile(0.5) = %v, want ~5.5", got.Val)
				}
				return nil
			},
		},
		{
			Name: "quantile 0.9 check",
			Args: args{
				opts:   []TDigestOption{WithTDigestCompression(100)},
				values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				q:      0.9,
			},
			CheckFunc: func(tt *testing.T, want test.Result[float64], got test.Result[float64]) error {
				if got.Val < 9 || got.Val > 10 {
					return errors.Errorf("Quantile(0.9) = %v, want ~9.5", got.Val)
				}
				return nil
			},
		},
		{
			Name: "quantile 0 (min)",
			Args: args{
				opts:   []TDigestOption{WithTDigestCompression(100)},
				values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				q:      0,
			},
			Want: test.Result[float64]{
				Val: 1.0,
			},
		},
		{
			Name: "quantile 1 (max)",
			Args: args{
				opts:   []TDigestOption{WithTDigestCompression(100)},
				values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				q:      1,
			},
			Want: test.Result[float64]{
				Val: 10.0,
			},
		},
		{
			Name: "empty tdigest",
			Args: args{
				opts:   []TDigestOption{WithTDigestCompression(100)},
				values: []float64{},
				q:      0.5,
			},
			Want: test.Result[float64]{
				Val: 0,
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestTDigest_Merge(t *testing.T) {
	type args struct {
		vals1 []float64
		vals2 []float64
	}

	if err := test.Run(t.Context(), t, func(tt *testing.T, args args) (TDigest, error) {
		td1, err := NewTDigest(WithTDigestCompression(100), WithTDigestCompressionTriggerFactor(10))
		if err != nil {
			return nil, err
		}
		for _, v := range args.vals1 {
			td1.Add(v)
		}

		td2, err := NewTDigest(WithTDigestCompression(100), WithTDigestCompressionTriggerFactor(10))
		if err != nil {
			return nil, err
		}
		for _, v := range args.vals2 {
			td2.Add(v)
		}

		if err := td1.Merge(td2); err != nil {
			return nil, err
		}
		return td1, nil
	}, []test.Case[TDigest, args]{
		{
			Name: "merge disjoint sets",
			Args: args{
				vals1: []float64{1, 2, 3, 4, 5},
				vals2: []float64{6, 7, 8, 9, 10},
			},
			CheckFunc: func(tt *testing.T, want test.Result[TDigest], got test.Result[TDigest]) error {
				if got.Err != nil {
					return got.Err
				}
				td := got.Val
				// Need to verify count. But sharded TDigest doesn't expose count directly, only internally.
				// We can check quantiles.
				q := td.Quantile(0.5)
				if q < 5 || q > 6 {
					return errors.Errorf("Quantile(0.5) after merge = %v, want ~5.5", q)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestTDigest_Compression(t *testing.T) {
	type args struct {
		count       int
		compression float64
	}

	if err := test.Run(t.Context(), t, func(tt *testing.T, args args) (TDigest, error) {
		td, err := NewTDigest(WithTDigestCompression(args.compression), WithTDigestCompressionTriggerFactor(1.1))
		if err != nil {
			return nil, err
		}
		for i := 0; i < args.count; i++ {
			td.Add(float64(i))
		}
		// Force flush by calling Quantile
		td.Quantile(0)
		return td, nil
	}, []test.Case[TDigest, args]{
		{
			Name: "aggressive compression",
			Args: args{
				count:       1000,
				compression: 20,
			},
			CheckFunc: func(tt *testing.T, want test.Result[TDigest], got test.Result[TDigest]) error {
				var totalCentroids int
				if td, ok := got.Val.(*shardedTDigest); ok {
					for _, shard := range td.shards {
						shard.mu.Lock()
						totalCentroids += len(shard.centroids)
						shard.mu.Unlock()
					}
				} else {
					td := got.Val.(*tdigest)
					totalCentroids = len(td.centroids)
				}

				// Since we have multiple shards (or not), each shard is compressed independently.
				// Each shard should be small. But total number might be large if many shards.
				// 16 shards. 1000 items. Each shard gets ~62 items.
				// Compression 20 means each shard should have ~20 centroids or less.
				// 16 * 20 = 320 max centroids total.

				if totalCentroids > 500 {
					return errors.Errorf("total centroids = %v, want <= 500 (approx)", totalCentroids)
				}
				q := got.Val.Quantile(0.5)
				if math.Abs(q-500) > 50 {
					return errors.Errorf("Quantile(0.5) after compression = %v, want ~500", q)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestTDigest_Concurrency(t *testing.T) {
	if err := test.Run(t.Context(), t, func(tt *testing.T, count int) (TDigest, error) {
		td, err := NewTDigest(WithTDigestCompression(100), WithTDigestCompressionTriggerFactor(10))
		if err != nil {
			return nil, err
		}
		var wg sync.WaitGroup
		for i := range count {
			wg.Add(1)
			go func(v float64) {
				defer wg.Done()
				td.Add(v)
			}(float64(i))
		}
		wg.Wait()
		td.Quantile(0) // Flush
		return td, nil
	}, []test.Case[TDigest, int]{
		{
			Name: "concurrent adds",
			Args: 100,
			CheckFunc: func(tt *testing.T, want test.Result[TDigest], got test.Result[TDigest]) error {
				td := got.Val
				q := td.Quantile(0.5)
				if math.Abs(q-49.5) > 5 { // 0..99, mean ~49.5
					return errors.Errorf("Quantile(0.5) after concurrent adds = %v, want ~49.5", q)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestTDigest_Quantile_Accuracy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		values    []float64
		weights   []float64 // If empty, assume weight 1
		quantile  float64
		want      float64
		tolerance float64
	}{
		{
			name:      "Simple Median",
			values:    []float64{10, 20},
			weights:   []float64{1, 1},
			quantile:  0.5,
			want:      15.0,
			tolerance: 0.001,
		},
		{
			name:      "Three values",
			values:    []float64{10, 20, 30},
			weights:   []float64{1, 1, 1},
			quantile:  0.5,
			want:      20.0,
			tolerance: 0.001,
		},
		{
			name:    "Weighted Median",
			values:  []float64{10, 30},
			weights: []float64{1, 3}, // Total 4. Mean of C0=10 (w=1), C1=30 (w=3).
			// Centers: C0 at 0.5. C1 at 1 + 1.5 = 2.5.
			// Target q=0.5 -> 2.0.
			// 2.0 is between 0.5 and 2.5.
			// Interpolate: prev=10, prevCenter=0.5. curr=30, center=2.5.
			// frac = (2.0 - 0.5) / (2.5 - 0.5) = 1.5 / 2.0 = 0.75.
			// Val = 10 + 0.75 * (30 - 10) = 10 + 15 = 25.
			quantile:  0.5,
			want:      25.0,
			tolerance: 0.001,
		},
		{
			name:      "Min",
			values:    []float64{10, 20},
			quantile:  0.0,
			want:      10.0,
			tolerance: 0.001,
		},
		{
			name:      "Max",
			values:    []float64{10, 20},
			quantile:  1.0,
			want:      20.0,
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td, err := NewTDigest(WithTDigestCompression(100))
			if err != nil {
				t.Fatal(err)
			}
			impl := td.(*tdigest)

			for i, v := range tt.values {
				w := 1.0
				if len(tt.weights) > i {
					w = tt.weights[i]
				}
				// Manually add centroid to bypass buffer/merge logic for precise test control
				impl.centroids = append(impl.centroids, centroid{Mean: v, Weight: w})
				impl.count += w
			}

			got := td.Quantile(tt.quantile)
			if diff := got - tt.want; diff < -tt.tolerance || diff > tt.tolerance {
				t.Errorf("Quantile(%v) = %v, want %v (diff %v)", tt.quantile, got, tt.want, diff)
			}
		})
	}
}

func TestGlobalSnapshot_JSON(t *testing.T) {
	t.Parallel()

	// 1. Create a collector and populate it
	c, err := NewCollector()
	if err != nil {
		t.Fatal(err)
	}

	// Add some data
	c.Record(nil, 0, &RequestResult{
		Latency:   100, // ns
		QueueWait: 50,  // ns
	})
	c.Record(nil, 0, &RequestResult{
		Latency:   200, // ns
		QueueWait: 60,  // ns
	})

	// 2. Snapshot
	snap := c.GlobalSnapshot()

	// 3. Marshal
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}

	// 4. Unmarshal into new snapshot
	var restored GlobalSnapshot
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatal(err)
	}

	// 5. Verify TDigests are not nil and contain data
	if restored.LatPercentiles == nil {
		t.Fatal("LatPercentiles is nil after unmarshal")
	}
	if restored.QWPercentiles == nil {
		t.Fatal("QWPercentiles is nil after unmarshal")
	}

	// Check counts via interface
	// Since TDigest interface doesn't have Count(), we check via Quantile
	// With 2 samples, min/max should be preserved roughly or at least Quantile works.
	// Actually we know it's *tdigest so we can check.

	latTD, ok := restored.LatPercentiles.(*tdigest)
	if !ok {
		t.Fatal("Restored LatPercentiles is not *tdigest")
	}
	if latTD.count != 2 {
		t.Errorf("Restored Lat count = %v, want 2", latTD.count)
	}
	// Check quantiles
	// With {100, 200}, median should be 150
	if got := restored.LatPercentiles.Quantile(0.5); got != 150 {
		t.Errorf("Restored Lat Median = %v, want 150", got)
	}
}

func TestTDigest_MarshalJSON(t *testing.T) {
	t.Parallel()

	td, _ := NewTDigest()
	td.Add(10)
	td.Add(20)

	data, err := json.Marshal(td)
	if err != nil {
		t.Fatal(err)
	}

	var td2 tdigest
	if err := json.Unmarshal(data, &td2); err != nil {
		t.Fatal(err)
	}

	if td2.count != 2 {
		t.Errorf("Count = %v, want 2", td2.count)
	}
	if q := td2.Quantile(0.5); q != 15 {
		t.Errorf("Quantile(0.5) = %v, want 15", q)
	}
}

// Benchmark for memory allocation.
// Can be run with -bench=.
func BenchmarkTDigest_Add(b *testing.B) {
	td, _ := NewTDigest()
	b.ResetTimer()
	var i uint64
	for b.Loop() {
		td.Add(float64(i))
		i++
	}
}
