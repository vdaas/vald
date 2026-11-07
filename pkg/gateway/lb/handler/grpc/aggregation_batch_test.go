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
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

func batchResponse(n int) *payload.Search_Response {
	ds := make([]*payload.Object_Distance, 0, n)
	for i := range n {
		ds = append(ds, dist32("id-"+strconv.Itoa(i), float32(i+1)/10))
	}
	return newDistResponse(ds...)
}

// TestSliceAggregators_SendAllocsPerResponse pins that the slice aggregators
// copy one agent response into a single DistPayload slab (one allocation)
// instead of allocating one *DistPayload per result.
func TestSliceAggregators_SendAllocsPerResponse(t *testing.T) {
	// not parallel: AllocsPerRun must not race with other allocating tests.
	const n = 16
	resp := batchResponse(n)
	ctx := t.Context()
	for name, newAggr := range map[string]func(num, fnum, replica int) Aggregator{
		"valdSliceAggr":     newSlice,
		"valdPoolSliceAggr": newPoolSlice,
	} {
		aggr := newAggr(n, n, 1)
		reset := func() {}
		switch v := aggr.(type) {
		case *valdSliceAggr:
			v.result = make([]*DistPayload, 0, n)
			reset = func() { v.result = v.result[:0] }
		case *valdPoolSliceAggr:
			v.result = make([]*DistPayload, 0, n)
			reset = func() { v.result = v.result[:0] }
		default:
			t.Fatalf("%s: unexpected aggregator %T", name, aggr)
		}
		if got := testing.AllocsPerRun(100, func() {
			reset()
			aggr.Send(ctx, resp)
		}); got > 1 {
			t.Errorf("%s: Send allocated %v times for %d results, want <= 1", name, got, n)
		}
		if got := len(aggr.Result().GetResults()); got != n {
			t.Errorf("%s: %d results, want %d", name, got, n)
		}
	}
}

// TestPairingHeapAggr_DuplicateSendAllocs pins that re-sending already visited
// IDs (the replica case) is deduplicated through the plain visited map without
// a per-result allocation.
func TestPairingHeapAggr_DuplicateSendAllocs(t *testing.T) {
	// not parallel: AllocsPerRun must not race with other allocating tests.
	const n = 16
	resp := batchResponse(n)
	ctx := t.Context()
	aggr := newPairingHeap(n, n, 1)
	aggr.Send(ctx, resp)
	if got := testing.AllocsPerRun(100, func() { aggr.Send(ctx, resp) }); got > 1 {
		t.Errorf("duplicate Send allocated %v times for %d results, want <= 1", got, n)
	}
	res := aggr.Result().GetResults()
	got := make([]string, 0, len(res))
	for _, r := range res {
		got = append(got, r.GetId())
	}
	want := make([]string, 0, n)
	for i := range n {
		want = append(want, "id-"+strconv.Itoa(i))
	}
	if !slices.Equal(got, want) {
		t.Errorf("ids = %v, want %v", got, want)
	}
}

// TestBatchAggregators_CanceledSendAddsNothing pins that a Send whose ctx is
// already done contributes no result, and that the next live Send still does.
func TestBatchAggregators_CanceledSendAddsNothing(t *testing.T) {
	t.Parallel()
	for name, newAggr := range map[string]func(num, fnum, replica int) Aggregator{
		"valdPairingHeapAggr": newPairingHeap,
		"valdSliceAggr":       newSlice,
		"valdPoolSliceAggr":   newPoolSlice,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			canceled, cancel := context.WithCancel(t.Context())
			cancel()
			aggr := newAggr(4, 4, 1)
			aggr.Start(t.Context())
			aggr.Send(canceled, newDistResponse(dist32("a", 0.1), dist32("b", 0.2)))
			// "b" was in the canceled Send too; the live Send must still add it.
			aggr.Send(t.Context(), newDistResponse(dist32("b", 0.3), dist32("c", 0.4)))
			res := aggr.Result().GetResults()
			got := make([]string, 0, len(res))
			for _, r := range res {
				got = append(got, r.GetId())
			}
			if want := []string{"b", "c"}; !slices.Equal(got, want) {
				t.Errorf("ids = %v, want %v", got, want)
			}
		})
	}
}
