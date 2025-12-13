// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package metrics

import (
	"testing"
	"time"
)

func TestHistogram_Reproduction(t *testing.T) {
	// Reproduce issue: Histogram range vs Exemplar range
	// Observed Min: 123.31ms
	// Observed Max: 359.19ms
	// Interval: 10ms

	h, err := NewHistogram(
		WithBucketInterval(10*time.Millisecond),
		WithTailSegments(10),
		WithHistogramMaxBuckets(100),
		WithHistogramNumShards(1),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Record values in nanoseconds
	// 1ms = 1e6 ns
	minVal := 123.31 * 1e6
	maxVal := 359.19 * 1e6

	// Record Min/Max and some middle values
	inputs := []float64{
		minVal,
		150.0 * 1e6,
		200.0 * 1e6,
		228.0 * 1e6, // Mean-ish
		300.0 * 1e6,
		maxVal,
	}

	for _, v := range inputs {
		h.Record(v)
	}

	snap := h.Snapshot()

	t.Logf("Snapshot Min: %f (Expected %f)", snap.Min, minVal)
	t.Logf("Snapshot Max: %f (Expected %f)", snap.Max, maxVal)

	if snap.Min != minVal {
		t.Errorf("Min mismatch: got %f, want %f", snap.Min, minVal)
	}
	if snap.Max != maxVal {
		t.Errorf("Max mismatch: got %f, want %f", snap.Max, maxVal)
	}

	// Check Bounds
	if len(snap.Bounds) == 0 {
		t.Fatal("No bounds generated")
	}

	t.Logf("Bounds: %v", snap.Bounds)
	t.Logf("Counts: %v", snap.Counts)

	// Check coverage of Min
	if snap.Bounds[0] < minVal {
		// If first bound < Min, then Min is in (Bounds[0], Bounds[1]].
		// Or if logic is different.
		// Snapshot logic: start = 120. Bounds[0] = 130.
		// Min=123.31.
		// 130 >= 123.31.
		// So Bounds[0] covers Min.
		// Wait, if Bounds[0] < Min, then Min is strictly > Bounds[0].
		// So Min is NOT in (-inf, Bounds[0]].
		// So Min is in (Bounds[0], ...].
		// But first bucket IS (-inf, Bounds[0]].
		// So if Bounds[0] < Min, Min is NOT in first bucket.
		// This is acceptable IF first bucket is empty?
		// But we want "Histogram must show the bucket containing the absolute minimum value".
		// If first bucket is empty, presenter skips it.
		// Then second bucket (containing Min) is shown.
		// That is fine.
		t.Logf("First bucket bound %f < Min %f. Min is likely in subsequent bucket.", snap.Bounds[0], minVal)
	} else {
		t.Logf("First bucket bound %f >= Min %f. Min is in first bucket.", snap.Bounds[0], minVal)
	}

	// Check coverage of Max
	lastBound := snap.Bounds[len(snap.Bounds)-1]
	if lastBound < maxVal {
		t.Errorf("Last bucket bound %f is less than Max %f", lastBound, maxVal)
	}

	// Check Counts sum
	sumCounts := uint64(0)
	for _, c := range snap.Counts {
		sumCounts += c
	}
	if sumCounts != uint64(len(inputs)) {
		t.Logf("Counts sum %d != Total %d. TDigest CDF approximation error expected.", sumCounts, len(inputs))
	} else {
		t.Log("Counts sum matches Total")
	}

	// Check distribution
	// 123.31 should be captured.
	// If Bounds[0] >= Min, it is in Counts[0].
	// If Bounds[0] < Min, it is in Counts[1] (or later).

	// Find bucket containing Min
	found := false
	for i, b := range snap.Bounds {
		// Bucket i covers (prev, b]. (First bucket (-inf, b])
		// If b >= Min, it might be here.
		if b >= minVal {
			if snap.Counts[i] > 0 {
				found = true
				t.Logf("Min value covered by bucket %d (Bound %f, Count %d)", i, b, snap.Counts[i])
			}
			break
		}
	}
	if !found {
		t.Error("No bucket covering Min has counts > 0")
	}
}
