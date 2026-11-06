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
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

// TestAggregators_DuplicateIDs sends overlapping result IDs from two agents
// (as happens with replicas) into every aggregator. Each ID must be returned
// once, the first-arrived entry is the one kept (as the std and pairing heap
// aggregators do via their visited maps) and the result is truncated to num
// only after de-duplication.
func TestAggregators_DuplicateIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		wantDist map[string]float32
		name     string
		respA    []*payload.Object_Distance
		respB    []*payload.Object_Distance
		// wantIDs is the ascending-distance order; ties keep arrival order.
		// It is checked for the slice aggregators only (heap ties are unordered).
		wantIDs []string
		num     int
	}{
		{
			name:     "equal distances for the same ID",
			num:      10,
			respA:    []*payload.Object_Distance{dist32("x", 0.1), dist32("y", 0.2), dist32("z", 0.4)},
			respB:    []*payload.Object_Distance{dist32("x", 0.1), dist32("y", 0.2), dist32("w", 0.3)},
			wantIDs:  []string{"x", "y", "w", "z"},
			wantDist: map[string]float32{"x": 0.1, "y": 0.2, "w": 0.3, "z": 0.4},
		},
		{
			name:     "different distances for the same ID keep the first arrival",
			num:      10,
			respA:    []*payload.Object_Distance{dist32("x", 0.5), dist32("y", 0.2)},
			respB:    []*payload.Object_Distance{dist32("x", 0.1), dist32("z", 0.3)},
			wantIDs:  []string{"y", "z", "x"},
			wantDist: map[string]float32{"x": 0.5, "y": 0.2, "z": 0.3},
		},
		{
			name:     "truncation counts distinct IDs only",
			num:      3,
			respA:    []*payload.Object_Distance{dist32("a", 0.1), dist32("b", 0.2), dist32("c", 0.3)},
			respB:    []*payload.Object_Distance{dist32("a", 0.1), dist32("b", 0.2), dist32("c", 0.3), dist32("d", 0.4)},
			wantIDs:  []string{"a", "b", "c"},
			wantDist: map[string]float32{"a": 0.1, "b": 0.2, "c": 0.3},
		},
		{
			name:     "fewer distinct IDs than num",
			num:      5,
			respA:    []*payload.Object_Distance{dist32("a", 0.1), dist32("b", 0.2)},
			respB:    []*payload.Object_Distance{dist32("b", 0.2), dist32("a", 0.1)},
			wantIDs:  []string{"a", "b"},
			wantDist: map[string]float32{"a": 0.1, "b": 0.2},
		},
	}
	aggrs := []struct {
		newFn   func(num, fnum, replica int) Aggregator
		name    string
		ordered bool
	}{
		{name: "std", newFn: newStd},
		{name: "slice", newFn: newSlice, ordered: true},
		{name: "pool slice", newFn: newPoolSlice, ordered: true},
		{name: "pairing heap", newFn: newPairingHeap},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			for _, ag := range aggrs {
				tt.Run(ag.name, func(ttt *testing.T) {
					ttt.Parallel()
					ctx, cancel := context.WithTimeout(ttt.Context(), 10*time.Second)
					defer cancel()
					aggr := ag.newFn(test.num, 10, 2)
					aggr.Start(ctx)
					aggr.Send(ctx, newDistResponse(test.respA...))
					aggr.Send(ctx, newDistResponse(test.respB...))
					res := aggr.Result().GetResults()

					if len(res) != len(test.wantDist) {
						ttt.Errorf("len(results) = %d, want %d: %v", len(res), len(test.wantDist), res)
					}
					gotDist := make(map[string]float32, len(res))
					for _, r := range res {
						if _, dup := gotDist[r.GetId()]; dup {
							ttt.Errorf("id %q returned more than once: %v", r.GetId(), res)
						}
						gotDist[r.GetId()] = r.GetDistance()
					}
					for id, d := range test.wantDist {
						if got, ok := gotDist[id]; !ok || got != d {
							ttt.Errorf("id %q: distance = %v (present: %t), want %v", id, got, ok, d)
						}
					}
					if !slices.IsSortedFunc(res, func(l, r *payload.Object_Distance) int {
						return cmpDist(&DistPayload{distance: float64(l.GetDistance())}, &DistPayload{distance: float64(r.GetDistance())})
					}) {
						ttt.Errorf("results are not sorted by ascending distance: %v", res)
					}
					if ag.ordered {
						got := make([]string, 0, len(res))
						for _, r := range res {
							got = append(got, r.GetId())
						}
						if !slices.Equal(got, test.wantIDs) {
							ttt.Errorf("ids = %v, want %v", got, test.wantIDs)
						}
					}
				})
			}
		})
	}
}

// TestDedupSortByID_ClearsVacatedTail pins that the dropped duplicates are not
// kept reachable through the backing array (the pooled slice is reused).
func TestDedupSortByID_ClearsVacatedTail(t *testing.T) {
	t.Parallel()
	x := []*DistPayload{
		{raw: dist32("a", 0.3), distance: float64(float32(0.3))},
		{raw: dist32("a", 0.1), distance: float64(float32(0.1))},
		{raw: dist32("b", 0.2), distance: float64(float32(0.2))},
		{raw: dist32("b", 0.2), distance: float64(float32(0.2))},
	}
	got := dedupSortByID(x)
	if len(got) != 2 || got[0].raw.GetId() != "b" || got[1].raw.GetId() != "a" || got[1].raw.GetDistance() != 0.3 {
		t.Errorf("dedupSortByID = %v, want [b@0.2 a@0.3 (first arrival)]", got)
	}
	for i, p := range x[len(got):] {
		if p != nil {
			t.Errorf("x[%d] = %v, want nil (vacated tail must be cleared)", len(got)+i, p)
		}
	}
}
