//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
)

func TestHistogram_Merge(t *testing.T) {
	h1, err := NewHistogram(
		WithHistogramMin(1),
		WithHistogramMax(100),
		WithHistogramGrowth(1.6),
		WithHistogramNumBuckets(10),
		WithHistogramNumShards(2),
	)
	if err != nil {
		t.Fatalf("Failed to create h1: %v", err)
	}
	h2, err := NewHistogram(
		WithHistogramMin(1),
		WithHistogramMax(100),
		WithHistogramGrowth(1.6),
		WithHistogramNumBuckets(10),
		WithHistogramNumShards(2),
	)
	if err != nil {
		t.Fatalf("Failed to create h2: %v", err)
	}
	h3, err := NewHistogram(
		WithHistogramMin(1),
		WithHistogramMax(100),
		WithHistogramGrowth(1.6),
		WithHistogramNumBuckets(20),
		WithHistogramNumShards(2),
	) // Incompatible
	if err != nil {
		t.Fatalf("Failed to create h3: %v", err)
	}

	h1.Record(10)
	h1.Record(20)

	h2.Record(30)
	h2.Record(40)

	if err := h1.Merge(h2); err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	snap := h1.Snapshot()
	if snap.Total != 4 {
		t.Errorf("snap.Total = %d, want 4", snap.Total)
	}
	if snap.Sum != 100 {
		t.Errorf("snap.Sum = %f, want 100", snap.Sum)
	}

	if err := h1.Merge(h3); err == nil {
		t.Error("Merge with incompatible histogram should have failed")
	}
}

func TestHistogramSnapshot_Merge(t *testing.T) {
	s1 := &HistogramSnapshot{
		Total:  2,
		Sum:    30,
		SumSq:  500,
		Min:    10,
		Max:    20,
		Counts: []uint64{0, 2, 0},
	}
	s2 := &HistogramSnapshot{
		Total:  2,
		Sum:    70,
		SumSq:  2500,
		Min:    30,
		Max:    40,
		Counts: []uint64{0, 0, 2},
	}

	s1.Merge(s2)

	if s1.Total != 4 {
		t.Errorf("s1.Total = %d, want 4", s1.Total)
	}
	if s1.Sum != 100 {
		t.Errorf("s1.Sum = %f, want 100", s1.Sum)
	}
	if s1.SumSq != 3000 {
		t.Errorf("s1.SumSq = %f, want 3000", s1.SumSq)
	}
	if s1.Min != 10 {
		t.Errorf("s1.Min = %f, want 10", s1.Min)
	}
	if s1.Max != 40 {
		t.Errorf("s1.Max = %f, want 40", s1.Max)
	}
	if s1.Mean != 25 {
		t.Errorf("s1.Mean = %f, want 25", s1.Mean)
	}
	if math.Abs(s1.StdDev-11.180340) > 0.0001 {
		t.Errorf("s1.StdDev = %f, want ~11.18", s1.StdDev)
	}
}
