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
				td := got.Val.(*tdigest)
				if td.count != 10 {
					return errors.Errorf("td1.count = %v, want 10", td.count)
				}
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

	if err := test.Run(t.Context(), t, func(tt *testing.T, args args) (*tdigest, error) {
		td, err := NewTDigest(WithTDigestCompression(args.compression), WithTDigestCompressionTriggerFactor(1.1))
		if err != nil {
			return nil, err
		}
		for i := 0; i < args.count; i++ {
			td.Add(float64(i))
		}
		// Force flush by calling Quantile
		td.Quantile(0)
		return td.(*tdigest), nil
	}, []test.Case[*tdigest, args]{
		{
			Name: "aggressive compression",
			Args: args{
				count:       1000,
				compression: 20,
			},
			CheckFunc: func(tt *testing.T, want test.Result[*tdigest], got test.Result[*tdigest]) error {
				td := got.Val
				if len(td.centroids) > 25 {
					// Should be roughly compression/2 to compression size, 25 is reasonable check for 20
					// Wait, K centroids. If compression is delta.
					// The number of centroids is roughly proportional to delta.
					// Just keep the check loose but verifying compression happened.
					// Without compression, it would be 1000 centroids.
					return errors.Errorf("len(td.centroids) = %v, want <= 25", len(td.centroids))
				}
				q := td.Quantile(0.5)
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
	if err := test.Run(t.Context(), t, func(tt *testing.T, count int) (*tdigest, error) {
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
		return td.(*tdigest), nil
	}, []test.Case[*tdigest, int]{
		{
			Name: "concurrent adds",
			Args: 100,
			CheckFunc: func(tt *testing.T, want test.Result[*tdigest], got test.Result[*tdigest]) error {
				td := got.Val
				if td.count != 100 {
					return errors.Errorf("td.count = %v, want 100", td.count)
				}
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
