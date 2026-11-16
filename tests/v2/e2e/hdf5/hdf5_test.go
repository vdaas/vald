//go:build e2e

// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package hdf5

import (
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/sync"
)

func newNoiseTestDataset() *Dataset {
	train := [][]float32{
		{1, 2, 3, 4, 5, 6, 7, 8},
		{8, 7, 6, 5, 4, 3, 2, 1},
	}
	test := [][]float32{
		{2, 4, 6, 8, 1, 3, 5, 7},
		{7, 5, 3, 1, 8, 6, 4, 2},
	}
	return New(train, test, [][]int{{0}, {1}}, 2)
}

func TestDataset_NoiseOnlyWhenNumExceedsMaxLen(t *testing.T) {
	t.Parallel()
	d := newNoiseTestDataset()

	// A request larger than the dataset is noised.
	if got := d.TrainCycle(10, 0).At(0); slices.Equal(got, d.Train[0]) {
		t.Fatalf("TrainCycle(num > maxLen).At(0) = %v, want noised vector", got)
	}
	// Later requests that fit in the dataset must not inherit the noise.
	if got := d.TrainCycle(2, 0).At(0); !slices.Equal(got, d.Train[0]) {
		t.Errorf("TrainCycle(num <= maxLen).At(0) = %v, want %v", got, d.Train[0])
	}
	if got := d.TestCycle(2, 0).At(1); !slices.Equal(got, d.Test[1]) {
		t.Errorf("TestCycle(num <= maxLen).At(1) = %v, want %v", got, d.Test[1])
	}
	// TestCycle above maxLen is noised as well.
	if got := d.TestCycle(10, 0).At(1); slices.Equal(got, d.Test[1]) {
		t.Errorf("TestCycle(num > maxLen).At(1) = %v, want noised vector", got)
	}
}

func TestDataset_ConcurrentNoiseInit(t *testing.T) {
	t.Parallel()
	d := newNoiseTestDataset()

	const workers = 8
	results := make([][]float32, workers)
	var wg sync.WaitGroup
	for w := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if w%2 == 0 {
				results[w] = d.TrainCycle(10, 0).At(3)
			} else {
				d.TestCycle(10, 0)
				results[w] = d.TrainCycle(10, 0).At(3)
			}
		}()
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Minute):
		t.Fatal("timed out waiting for concurrent TrainCycle/TestCycle calls")
	}

	// Every caller must observe the same (single) noise function.
	for w := 1; w < workers; w++ {
		if !slices.Equal(results[w], results[0]) {
			t.Errorf("worker %d got %v, want %v", w, results[w], results[0])
		}
	}
	if slices.Equal(results[0], d.Train[1]) {
		t.Errorf("TrainCycle(num > maxLen).At(3) = %v, want noised vector", results[0])
	}
}
