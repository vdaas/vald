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
	"sort"
	"sync"
	"testing"
)

func TestTDigest_AddAndQuantile(t *testing.T) {
	td, _ := NewTDigest(100, 10)

	// Add some values
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, v := range values {
		td.Add(v)
	}

	// Simple quantile checks
	if q := td.Quantile(0.5); q < 5 || q > 6 {
		t.Errorf("Quantile(0.5) = %v, want ~5.5", q)
	}
	if q := td.Quantile(0.9); q < 9 || q > 10 {
		t.Errorf("Quantile(0.9) = %v, want ~9.5", q)
	}

	// Edge cases
	if q := td.Quantile(0); q != 1 {
		t.Errorf("Quantile(0) = %v, want 1", q)
	}
	if q := td.Quantile(1); q != 10 {
		t.Errorf("Quantile(1) = %v, want 10", q)
	}
}

func TestTDigest_Merge(t *testing.T) {
	td1, _ := NewTDigest(100, 10)
	td2, _ := NewTDigest(100, 10)

	for i := 1; i <= 10; i++ {
		td1.Add(float64(i))
	}
	for i := 11; i <= 20; i++ {
		td2.Add(float64(i))
	}

	if err := td1.Merge(td2); err != nil {
		t.Fatalf("Merge() error = %v", err)
	}

	if td1.count != 20 {
		t.Errorf("td1.count = %v, want 20", td1.count)
	}

	if q := td1.Quantile(0.5); q < 10 || q > 11 {
		t.Errorf("Quantile(0.5) after merge = %v, want ~10.5", q)
	}
}

func TestTDigest_Compression(t *testing.T) {
	td, _ := NewTDigest(20, 1.1) // Aggressive compression

	for i := 0; i < 1000; i++ {
		td.Add(float64(i))
	}

	if len(td.centroids) > 25 { // Should be around 20
		t.Errorf("len(td.centroids) = %v, want <= 25", len(td.centroids))
	}

	// Check if quantiles are still reasonable
	if q := td.Quantile(0.5); math.Abs(q-500) > 50 {
		t.Errorf("Quantile(0.5) after compression = %v, want ~500", q)
	}
}

func TestTDigest_Empty(t *testing.T) {
	td, _ := NewTDigest(100, 10)
	if q := td.Quantile(0.5); q != 0 {
		t.Errorf("Quantile(0.5) on empty t-digest = %v, want 0", q)
	}
}

func TestTDigest_Concurrency(t *testing.T) {
	td, _ := NewTDigest(100, 10)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(v float64) {
			defer wg.Done()
			td.Add(v)
		}(float64(i))
	}

	wg.Wait()

	if td.count != 100 {
		t.Errorf("td.count = %v, want 100", td.count)
	}

	// Check quantiles after concurrent adds
	values := make([]float64, 100)
	for i := 0; i < 100; i++ {
		values[i] = float64(i)
	}
	sort.Float64s(values)

	p50 := values[49] // Approximate
	if q := td.Quantile(0.5); math.Abs(q-p50) > 5 {
		t.Errorf("Quantile(0.5) after concurrent adds = %v, want ~%v", q, p50)
	}
}
