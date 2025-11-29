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
