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
	"context"
	"math"
	"testing"
	"time"
)

func TestReproduction_StdDev(t *testing.T) {
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

func TestReproduction_HistogramBuckets(t *testing.T) {
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
