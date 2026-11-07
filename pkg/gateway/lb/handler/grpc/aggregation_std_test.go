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
package grpc

import (
	"context"
	"math"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

// TestValdStdAggr_FrontInsertKeepsCapacity pins that inserting a new minimum
// shifts within the num*replica buffer preallocated by newStd instead of
// reallocating a len+1 slice.
func TestValdStdAggr_FrontInsertKeepsCapacity(t *testing.T) {
	t.Parallel()
	const num, replica = 10, 2
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	aggr := newStd(num, num, replica)
	aggr.Start(ctx)
	// Descending distances: every element after the first becomes the new minimum
	// (the rl == 1 branch, then the pos == rl branch).
	aggr.Send(ctx, newDistResponse(dist32("d", 0.5), dist32("c", 0.4), dist32("b", 0.3), dist32("a", 0.2)))
	res := aggr.Result().GetResults()
	got := make([]string, 0, len(res))
	for _, r := range res {
		got = append(got, r.GetId())
	}
	if want := []string{"a", "b", "c", "d"}; !slices.Equal(got, want) {
		t.Errorf("ids = %v, want %v", got, want)
	}
	// result is preallocated for num entries plus the one add holds before trimming.
	if c := cap(res); c != num+1 {
		t.Errorf("cap(results) = %d, want %d (preallocated buffer reused)", c, num+1)
	}
}

// TestValdStdAggr_MaxDistBound pins the distance bound kept as float64 bits:
// it starts at MaxFloat64 and tightens to the num-th distance once num results
// are held.
func TestValdStdAggr_MaxDistBound(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	aggr := newStd(2, 10, 1)
	v, ok := aggr.(*valdStdAggr)
	if !ok {
		t.Fatalf("newStd returned %T, want *valdStdAggr", aggr)
	}
	if got := math.Float64frombits(v.maxDist.Load()); got != math.MaxFloat64 {
		t.Errorf("initial bound = %v, want MaxFloat64", got)
	}
	aggr.Start(ctx)
	aggr.Send(ctx, newDistResponse(dist32("a", 0.1), dist32("b", 0.2), dist32("c", 0.3)))
	res := aggr.Result().GetResults()
	if len(res) != 2 || res[0].GetId() != "a" || res[1].GetId() != "b" {
		t.Errorf("results = %v, want [a b]", res)
	}
	if got, want := math.Float64frombits(v.maxDist.Load()), float64(float32(0.2)); got != want {
		t.Errorf("bound = %v, want %v", got, want)
	}
}

// TestValdStdAggr_ConcurrentDuplicateIDs sends the same IDs from many
// concurrent Send callers (one per replica agent). De-duplication happens on
// the single aggregation goroutine; each ID must still be returned once.
func TestValdStdAggr_ConcurrentDuplicateIDs(t *testing.T) {
	t.Parallel()
	const num, senders = 10, 8
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	resp := make([]*payload.Object_Distance, 0, num)
	want := make([]string, 0, num)
	for i := range num {
		id := "id-" + strconv.Itoa(i)
		resp = append(resp, dist32(id, float32(i+1)/10))
		want = append(want, id)
	}
	aggr := newStd(num, num, senders)
	aggr.Start(ctx)
	done := make(chan struct{}, senders)
	for range senders {
		go func() {
			aggr.Send(ctx, newDistResponse(resp...))
			done <- struct{}{}
		}()
	}
	for range senders {
		select {
		case <-done:
		case <-ctx.Done():
			t.Fatal("timed out waiting for Send")
		}
	}
	res := aggr.Result().GetResults()
	got := make([]string, 0, len(res))
	for _, r := range res {
		got = append(got, r.GetId())
	}
	if !slices.Equal(got, want) {
		t.Errorf("ids = %v, want %v", got, want)
	}
}
