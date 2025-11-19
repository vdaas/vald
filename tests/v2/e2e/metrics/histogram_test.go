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
	"sync"
	"testing"
)

func TestHistogram(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name    string
		h1      func() (Histogram, error)
		h2      func() (Histogram, error)
		records []float64
		check   func(t *testing.T, h Histogram)
	}

	tests := []testCase{
		{
			name: "record values and check snapshot",
			h1: func() (Histogram, error) {
				return NewHistogram(
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				)
			},
			records: []float64{10, 20, 30, 40, 50},
			check: func(t *testing.T, h Histogram) {
				snap := h.Snapshot()
				if snap.Total != 5 {
					t.Errorf("expected total 5, got %d", snap.Total)
				}
				if snap.Mean != 30 {
					t.Errorf("expected mean 30, got %f", snap.Mean)
				}
			},
		},
		{
			name: "merge two histograms",
			h1: func() (Histogram, error) {
				h, err := NewHistogram(
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				)
				if err != nil {
					return nil, err
				}
				h.Record(10)
				h.Record(20)
				return h, nil
			},
			h2: func() (Histogram, error) {
				h, err := NewHistogram(
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				)
				if err != nil {
					return nil, err
				}
				h.Record(30)
				h.Record(40)
				return h, nil
			},
			check: func(t *testing.T, h Histogram) {
				snap := h.Snapshot()
				if snap.Total != 4 {
					t.Errorf("expected total 4, got %d", snap.Total)
				}
				if snap.Mean != 25 {
					t.Errorf("expected mean 25, got %f", snap.Mean)
				}
			},
		},
		{
			name: "empty histogram",
			h1: func() (Histogram, error) {
				return NewHistogram(
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				)
			},
			check: func(t *testing.T, h Histogram) {
				snap := h.Snapshot()
				if snap.Total != 0 {
					t.Errorf("expected total 0, got %d", snap.Total)
				}
				if snap.Mean != 0 {
					t.Errorf("expected mean 0, got %f", snap.Mean)
				}
			},
		},
		{
			name: "record a single value",
			h1: func() (Histogram, error) {
				return NewHistogram(
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				)
			},
			records: []float64{10},
			check: func(t *testing.T, h Histogram) {
				snap := h.Snapshot()
				if snap.Total != 1 {
					t.Errorf("expected total 1, got %d", snap.Total)
				}
				if snap.Mean != 10 {
					t.Errorf("expected mean 10, got %f", snap.Mean)
				}
			},
		},
		{
			name: "concurrent record",
			h1: func() (Histogram, error) {
				return NewHistogram(
					WithHistogramMin(1),
					WithHistogramGrowth(2),
					WithHistogramNumBuckets(10),
					WithHistogramNumShards(1),
				)
			},
			check: func(t *testing.T, h Histogram) {
				var wg sync.WaitGroup
				for i := 0; i < 100; i++ {
					wg.Add(1)
					go func(v float64) {
						defer wg.Done()
						h.Record(v)
					}(float64(i))
				}
				wg.Wait()

				snap := h.Snapshot()
				if snap.Total != 100 {
					t.Errorf("expected total 100, got %d", snap.Total)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h1, err := tt.h1()
			if err != nil {
				t.Fatalf("failed to create h1: %v", err)
			}

			for _, r := range tt.records {
				h1.Record(r)
			}

			h_to_check := h1
			if tt.h2 != nil {
				h2, err := tt.h2()
				if err != nil {
					t.Fatalf("failed to create h2: %v", err)
				}
				if err := h1.Merge(h2); err != nil {
					t.Fatalf("failed to merge histograms: %v", err)
				}
				h_to_check = h2
			}

			tt.check(t, h_to_check)
		})
	}
}
