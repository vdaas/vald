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
)

// TestAggregators_PreallocIsCapped pins that the capacities preallocated from
// the client-supplied num (num is not bounded by the server) are capped at
// maxAggrPrealloc, while the results for a large num are unchanged. num is
// chosen so the uncapped sizes are only a few MB, i.e. larger than the cap but
// still safe to allocate if the cap regresses.
func TestAggregators_PreallocIsCapped(t *testing.T) {
	t.Parallel()
	const num, replica = 1 << 16, 4
	want := []string{"a", "b", "c"}
	tests := []struct {
		newAggr func(num, fnum, replica int) Aggregator
		caps    func(Aggregator) []uint64 // preallocated capacities
		name    string
	}{
		{
			name:    "std aggregator caps buffers",
			newAggr: newStd,
			caps: func(a Aggregator) []uint64 {
				v, ok := a.(*valdStdAggr)
				if !ok {
					return nil
				}
				return []uint64{uint64(cap(v.dch)), uint64(cap(v.result))}
			},
		},
		{
			name:    "pairing heap aggregator caps buffers",
			newAggr: newPairingHeap,
			caps: func(a Aggregator) []uint64 {
				v, ok := a.(*valdPairingHeapAggr)
				if !ok {
					return nil
				}
				return []uint64{uint64(cap(v.result))}
			},
		},
		{
			name:    "slice aggregator caps buffers",
			newAggr: newSlice,
			caps: func(a Aggregator) []uint64 {
				v, ok := a.(*valdSliceAggr)
				if !ok {
					return nil
				}
				return []uint64{uint64(cap(v.result))}
			},
		},
		{
			name:    "pool slice aggregator caps poolLen",
			newAggr: newPoolSlice,
			caps: func(Aggregator) []uint64 {
				return []uint64{poolLen.Load()}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			aggr := tc.newAggr(num, num, replica)
			caps := tc.caps(aggr)
			if len(caps) == 0 {
				t.Fatalf("unexpected aggregator type %T", aggr)
			}
			for i, c := range caps {
				if c > maxAggrPrealloc {
					t.Errorf("preallocated capacity #%d = %d, want <= %d", i, c, maxAggrPrealloc)
				}
			}
			if got := aggr.GetNum(); got != num {
				t.Errorf("GetNum() = %d, want %d", got, num)
			}
			aggr.Start(ctx)
			aggr.Send(ctx, newDistResponse(dist32("c", 0.3), dist32("a", 0.1), dist32("b", 0.2)))
			res := aggr.Result().GetResults()
			got := make([]string, 0, len(res))
			for _, r := range res {
				got = append(got, r.GetId())
			}
			if !slices.Equal(got, want) {
				t.Errorf("ids = %v, want %v", got, want)
			}
			if c := cap(res); c > maxAggrPrealloc {
				t.Errorf("cap(results) = %d, want <= %d", c, maxAggrPrealloc)
			}
		})
	}
}
