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
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/sync/errgroup"
)

func TestExemplar_Offer(t *testing.T) {
	e := NewExemplar(3)

	e.Offer(100*time.Millisecond, "req-1")
	e.Offer(200*time.Millisecond, "req-2")
	e.Offer(50*time.Millisecond, "req-3")

	snap := e.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("Snapshot length = %d, want 3", len(snap))
	}

	e.Offer(300*time.Millisecond, "req-4")
	snap = e.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("Snapshot length = %d, want 3", len(snap))
	}

	minLatency := snap[0].latency
	for _, item := range snap {
		if item.latency < minLatency {
			minLatency = item.latency
		}
	}
	if minLatency != 100*time.Millisecond {
		t.Errorf("min latency = %v, want 100ms", minLatency)
	}
}

func TestExemplar_Concurrency(t *testing.T) {
	k := 10
	e := NewExemplar(k)
	numGoroutines := 100
	numOffersPerG := 20

	eg, _ := errgroup.New(t.Context())
	for i := range numGoroutines {
		idx := i
		eg.Go(func() error {
			for j := range numOffersPerG {
				lat := time.Duration(rand.Intn(1000)) * time.Millisecond
				reqID := fmt.Sprintf("req-%d-%d", idx, j)
				e.Offer(lat, reqID)
			}
			return nil
		})
	}

	snapshots := make([][]*item, 10)
	for i := range len(snapshots) {
		idx := i
		eg.Go(func() error {
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			snapshots[idx] = e.Snapshot()
			return nil
		})
	}

	_ = eg.Wait()

	finalSnap := e.Snapshot()
	if len(finalSnap) > k {
		t.Errorf("Final snapshot length = %d, want <= %d", len(finalSnap), k)
	}

	for _, snap := range snapshots {
		if len(snap) > k {
			t.Errorf("Intermediate snapshot length = %d, want <= %d", len(snap), k)
		}
	}
}

func TestExemplar_Race(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping race test in short mode")
	}

	k := 5
	e := NewExemplar(k)
	eg, _ := errgroup.New(t.Context())

	// Writer goroutines
	for range k {
		eg.Go(func() error {
			for j := 0; j < 100; j++ {
				e.Offer(time.Duration(rand.Intn(100))*time.Millisecond, "req")
			}
			return nil
		})
	}

	// Reader goroutines
	for range k {
		eg.Go(func() error {
			for j := 0; j < 100; j++ {
				snap := e.Snapshot()
				if len(snap) > k {
					t.Errorf("snapshot too large: %d", len(snap))
				}
			}
			return nil
		})
	}

	eg.Wait()
}
