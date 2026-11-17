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
	"cmp"
	"context"
	"math"
	"math/big"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func newDistResponse(ds ...*payload.Object_Distance) *payload.Search_Response {
	return &payload.Search_Response{Results: ds}
}

func dist32(id string, d float32) *payload.Object_Distance {
	return &payload.Object_Distance{Id: id, Distance: d}
}

// TestCmpDist_MatchesBigFloat pins that the float64 comparison orders every
// non-NaN float32 pair exactly like the former big.Float.Cmp did.
func TestCmpDist_MatchesBigFloat(t *testing.T) {
	t.Parallel()
	negZero := float32(math.Copysign(0, -1))
	vals := []float32{
		float32(math.Inf(-1)), -math.MaxFloat32, -1, -math.SmallestNonzeroFloat32,
		negZero, 0, math.SmallestNonzeroFloat32, 1.17549435e-38 /* min normal */, 0.5, 1,
		math.MaxFloat32, float32(math.Inf(1)),
	}
	for _, a := range vals {
		for _, b := range vals {
			want := big.NewFloat(float64(a)).Cmp(big.NewFloat(float64(b)))
			got := cmpDist(
				&DistPayload{distance: float64(a)},
				&DistPayload{distance: float64(b)},
			)
			if got != want {
				t.Errorf("cmpDist(%v, %v) = %d, big.Float.Cmp = %d", a, b, got, want)
			}
			// the operators used by the valdStdAggr insertion path
			if ge, wantGE := float64(a) >= float64(b), want >= 0; ge != wantGE {
				t.Errorf("%v >= %v = %t, big.Float.Cmp >= 0 = %t", a, b, ge, wantGE)
			}
		}
	}
}

func TestCmpDist_ZeroAlloc(t *testing.T) {
	// not parallel: AllocsPerRun must not race with other allocating tests.
	l, r := &DistPayload{distance: 0.25}, &DistPayload{distance: 0.5}
	var sink int
	if n := testing.AllocsPerRun(1000, func() { sink += cmpDist(l, r) }); n != 0 {
		t.Errorf("cmpDist allocated %v times per run, want 0", n)
	}
	if sink >= 0 {
		t.Errorf("sink = %d, want negative (l < r)", sink)
	}
}

// TestAggregators_DistanceOrdering feeds the same two agent responses (equal
// distances from different agents, +-0, +-Inf, subnormal) sequentially into
// every aggregator and checks ordering and tie behaviour.
func TestAggregators_DistanceOrdering(t *testing.T) {
	t.Parallel()
	negZero := float32(math.Copysign(0, -1))
	inf := float32(math.Inf(1))
	respA := func() *payload.Search_Response {
		return newDistResponse(
			dist32("a-ninf", float32(math.Inf(-1))),
			dist32("a-pz", 0),
			dist32("a-half", 0.5),
			dist32("a-inf", inf),
		)
	}
	respB := func() *payload.Search_Response {
		return newDistResponse(
			dist32("b-neg", -1),
			dist32("b-nz", negZero),
			dist32("b-sub", math.SmallestNonzeroFloat32),
			dist32("b-half", 0.5),
		)
	}
	// ties (a-pz/b-nz at +-0, a-half/b-half at 0.5) keep arrival order (A then B).
	stable := []string{"a-ninf", "b-neg", "a-pz", "b-nz", "b-sub", "a-half", "b-half"}
	tests := []struct {
		newFn func(num, fnum, replica int) Aggregator
		name  string
		want  []string // nil: only the distance sequence is checked (heap ties are unordered)
		// wantDists is compared with ==, so -0 == +0.
		wantDists []float32
	}{
		{
			name:  "std drops +Inf (>= initial MaxFloat64 bound, as before)",
			newFn: newStd,
			want:  stable,
		},
		{
			name:  "slice keeps +Inf last",
			newFn: newSlice,
			want:  append(slices.Clone(stable), "a-inf"),
		},
		{
			name:  "pool slice keeps +Inf last",
			newFn: newPoolSlice,
			want:  append(slices.Clone(stable), "a-inf"),
		},
		{
			name:  "pairing heap sorted by distance",
			newFn: newPairingHeap,
			wantDists: []float32{
				float32(math.Inf(-1)), -1, 0, 0, math.SmallestNonzeroFloat32, 0.5, 0.5, inf,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithTimeout(tt.Context(), 10*time.Second)
			defer cancel()
			aggr := test.newFn(10, 10, 2)
			aggr.Start(ctx)
			aggr.Send(ctx, respA())
			aggr.Send(ctx, respB())
			res := aggr.Result().GetResults()
			if test.want != nil {
				got := make([]string, 0, len(res))
				for _, r := range res {
					got = append(got, r.GetId())
				}
				if !slices.Equal(got, test.want) {
					tt.Errorf("ids = %v, want %v", got, test.want)
				}
				return
			}
			got := make([]float32, 0, len(res))
			for _, r := range res {
				got = append(got, r.GetDistance())
			}
			if !slices.EqualFunc(got, test.wantDists, func(a, b float32) bool { return cmp.Compare(a, b) == 0 }) {
				tt.Errorf("distances = %v, want %v", got, test.wantDists)
			}
		})
	}
}

// TestAggregators_NaNDistanceSkipped pins the NaN policy: NaN-distance results
// are skipped (not ranked, not returned), the other results are unaffected and
// nothing panics (big.NewFloat(NaN) used to).
func TestAggregators_NaNDistanceSkipped(t *testing.T) {
	t.Parallel()
	nan := float32(math.NaN())
	tests := []struct {
		newFn func(num, fnum, replica int) Aggregator
		name  string
	}{
		{name: "std", newFn: newStd},
		{name: "slice", newFn: newSlice},
		{name: "pool slice", newFn: newPoolSlice},
		{name: "pairing heap", newFn: newPairingHeap},
	}
	want := []string{"ok1", "ok2", "ok3"}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithTimeout(tt.Context(), 10*time.Second)
			defer cancel()
			aggr := test.newFn(10, 10, 2)
			aggr.Start(ctx)
			// NaN first, in the middle (must not stop processing of later entries) and last.
			aggr.Send(ctx, newDistResponse(
				dist32("nan1", nan), dist32("ok1", 0.1), dist32("nan2", nan), dist32("ok2", 0.2),
			))
			aggr.Send(ctx, newDistResponse(dist32("ok3", 0.3), dist32("nan3", nan)))
			// a NaN id is not marked visited: a later valid result with the same id is kept.
			aggr.Send(ctx, newDistResponse(dist32("nan1", 0.4)))
			res := aggr.Result().GetResults()
			got := make([]string, 0, len(res))
			for _, r := range res {
				if math.IsNaN(float64(r.GetDistance())) {
					tt.Errorf("NaN distance result %q was returned", r.GetId())
				}
				got = append(got, r.GetId())
			}
			if wantAll := append(slices.Clone(want), "nan1"); !slices.Equal(got, wantAll) {
				tt.Errorf("ids = %v, want %v", got, wantAll)
			}
		})
	}
}

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

func Test_searchRecallAtRisk(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		failures int
		replica  int
		want     bool
	}{
		{name: "no failure", failures: 0, replica: 3, want: false},
		{name: "fewer failures than replica", failures: 2, replica: 3, want: false},
		{name: "failures equal replica", failures: 3, replica: 3, want: true},
		{name: "failures above replica", failures: 4, replica: 3, want: true},
		{name: "single replica, one failure", failures: 1, replica: 1, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := searchRecallAtRisk(test.failures, test.replica); got != test.want {
				t.Errorf("searchRecallAtRisk(%d, %d) = %v, want %v", test.failures, test.replica, got, test.want)
			}
		})
	}
}

func Test_searchAgentFailure(t *testing.T) {
	t.Parallel()
	const target = "a:8081"
	tests := []struct {
		err      error
		name     string
		wantCode codes.Code
		want     bool
	}{
		{name: "nil", err: nil},
		{name: "context canceled", err: context.Canceled},
		{name: "context deadline", err: context.DeadlineExceeded},
		{name: "rpc call canceled", err: errors.ErrRPCCallFailed(target, context.Canceled)},
		{name: "rpc call deadline", err: errors.ErrRPCCallFailed(target, context.DeadlineExceeded)},
		{name: "status Canceled", err: status.Error(codes.Canceled, "c")},
		{name: "status DeadlineExceeded", err: status.Error(codes.DeadlineExceeded, "d")},
		{name: "status NotFound", err: status.Error(codes.NotFound, "nf")},
		{name: "status InvalidArgument is client-caused", err: status.Error(codes.InvalidArgument, "dim")},
		{name: "status Unavailable", err: status.Error(codes.Unavailable, "u"), wantCode: codes.Unavailable, want: true},
		{name: "status FailedPrecondition", err: status.Error(codes.FailedPrecondition, "fp"), wantCode: codes.FailedPrecondition, want: true},
		{name: "status Unimplemented", err: status.Error(codes.Unimplemented, "ui"), wantCode: codes.Unimplemented, want: true},
		{name: "non-status error", err: errors.New("boom"), wantCode: codes.Unknown, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			code, ok := searchAgentFailure(target, test.err)
			if ok != test.want || (ok && code != test.wantCode) {
				t.Errorf("searchAgentFailure(%v) = (%v, %v), want (%v, %v)", test.err, code, ok, test.wantCode, test.want)
			}
		})
	}
}

// The per-agent errors are logged at Debug; this per-request level decides
// whether the agent failures are a Warn.
func Test_partialSearchLogLevel(t *testing.T) {
	t.Parallel()
	bcErr := status.Error(codes.Unavailable, "agent down")
	tests := []struct {
		bcErr    error
		reqErr   error
		name     string
		failures int
		want     level.Level
	}{
		{name: "nothing failed", want: level.Unknown},
		{name: "failures below replica", failures: 1, bcErr: bcErr, want: level.DEBUG},
		{name: "failures below replica without broadcast error", failures: 1, want: level.DEBUG},
		{name: "failures reach replica", failures: 2, bcErr: bcErr, want: level.WARN},
		{name: "broadcast error without agent failure", bcErr: context.Canceled, want: level.WARN},
		{name: "failed request with agent failures", failures: 1, reqErr: bcErr, want: level.WARN},
		{name: "failed request without agent failures", reqErr: bcErr, want: level.Unknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := partialSearchLogLevel(test.failures, 2, test.bcErr, test.reqErr); got != test.want {
				t.Errorf("partialSearchLogLevel(%d, 2, %v, %v) = %v, want %v", test.failures, test.bcErr, test.reqErr, got, test.want)
			}
		})
	}
}

// A request that failed with agent failures is Warn-logged once: when the deferred summary logs
// the request error, aggregationSearch must not log it (or bcErr) again. A failing request without
// agent failures is not logged by the summary.
func Test_searchSummaryLogsRequestErr(t *testing.T) {
	t.Parallel()
	reqErr := status.Error(codes.Unavailable, "agent down")
	tests := []struct {
		reqErr   error
		name     string
		failures int
		want     bool
	}{
		{name: "successful request", failures: 2},
		{name: "failed request without agent failures", reqErr: reqErr},
		{name: "failed request with agent failures", failures: 1, reqErr: reqErr, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := searchSummaryLogsRequestErr(test.failures, test.reqErr)
			if got != test.want {
				t.Errorf("searchSummaryLogsRequestErr(%d, %v) = %v, want %v", test.failures, test.reqErr, got, test.want)
			}
			// the guard must agree with the summary's own decision, for any bcErr
			for _, bcErr := range []error{nil, reqErr} {
				summary := partialSearchLogLevel(test.failures, 2, bcErr, test.reqErr) == level.WARN && test.reqErr != nil
				if summary != got {
					t.Errorf("summary Warns the request error = %v, guard = %v (bcErr %v)", summary, got, bcErr)
				}
			}
		})
	}
}

func Test_server_searchQuorum(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		grace   time.Duration
		replica int
		agents  int
		want    int
	}{
		{name: "disabled waits for every agent", replica: 3, agents: 6},
		{name: "no replication waits for every agent", grace: time.Millisecond, replica: 1, agents: 6},
		{name: "agents not above replica waits for every agent", grace: time.Millisecond, replica: 3, agents: 3},
		{name: "6 agents with replica 3 need 4", grace: time.Millisecond, replica: 3, agents: 6, want: 4},
		{name: "20 agents with replica 2 need 19", grace: time.Millisecond, replica: 2, agents: 20, want: 19},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := &server{searchQuorumGrace: tc.grace, replica: tc.replica}
			if got := s.searchQuorum(tc.agents); got != tc.want {
				t.Errorf("searchQuorum(%d) = %d, want %d", tc.agents, got, tc.want)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_server_aggregationSearch(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		aggr Aggregator
// 		bcfg *payload.Search_Config
// 		f    func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error)
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantRes   *payload.Search_Response
// 		wantAttrs []attribute.KeyValue
// 		err       error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, []attribute.KeyValue, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, gotAttrs []attribute.KeyValue, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		if !reflect.DeepEqual(gotAttrs, w.wantAttrs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotAttrs, w.wantAttrs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           aggr:nil,
// 		           bcfg:nil,
// 		           f:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           aggr:nil,
// 		           bcfg:nil,
// 		           f:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotRes, gotAttrs, err := s.aggregationSearch(test.args.ctx, test.args.aggr, test.args.bcfg, test.args.f)
// 			if err := checkFunc(test.want, gotRes, gotAttrs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_newStd(t *testing.T) {
// 	type args struct {
// 		num     int
// 		fnum    int
// 		replica int
// 	}
// 	type want struct {
// 		want Aggregator
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Aggregator) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Aggregator) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           num:0,
// 		           fnum:0,
// 		           replica:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           num:0,
// 		           fnum:0,
// 		           replica:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := newStd(test.args.num, test.args.fnum, test.args.replica)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdStdAggr_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		maxDist atomic.Value
// 		eg      errgroup.Group
// 		dch     chan DistPayload
// 		cancel  context.CancelFunc
// 		visited sync.Map[string, any]
// 		result  []*payload.Object_Distance
// 		num     int
// 		fnum    int
// 		closed  atomic.Bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           maxDist:nil,
// 		           eg:nil,
// 		           dch:nil,
// 		           cancel:nil,
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           closed:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           maxDist:nil,
// 		           eg:nil,
// 		           dch:nil,
// 		           cancel:nil,
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           closed:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdStdAggr{
// 				maxDist: test.fields.maxDist,
// 				eg:      test.fields.eg,
// 				dch:     test.fields.dch,
// 				cancel:  test.fields.cancel,
// 				visited: test.fields.visited,
// 				result:  test.fields.result,
// 				num:     test.fields.num,
// 				fnum:    test.fields.fnum,
// 				closed:  test.fields.closed,
// 			}
//
// 			v.Start(test.args.ctx)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdStdAggr_Send(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		data *payload.Search_Response
// 	}
// 	type fields struct {
// 		maxDist atomic.Value
// 		eg      errgroup.Group
// 		dch     chan DistPayload
// 		cancel  context.CancelFunc
// 		visited sync.Map[string, any]
// 		result  []*payload.Object_Distance
// 		num     int
// 		fnum    int
// 		closed  atomic.Bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           data:nil,
// 		       },
// 		       fields: fields {
// 		           maxDist:nil,
// 		           eg:nil,
// 		           dch:nil,
// 		           cancel:nil,
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           closed:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           data:nil,
// 		           },
// 		           fields: fields {
// 		           maxDist:nil,
// 		           eg:nil,
// 		           dch:nil,
// 		           cancel:nil,
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           closed:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdStdAggr{
// 				maxDist: test.fields.maxDist,
// 				eg:      test.fields.eg,
// 				dch:     test.fields.dch,
// 				cancel:  test.fields.cancel,
// 				visited: test.fields.visited,
// 				result:  test.fields.result,
// 				num:     test.fields.num,
// 				fnum:    test.fields.fnum,
// 				closed:  test.fields.closed,
// 			}
//
// 			v.Send(test.args.ctx, test.args.data)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdStdAggr_Result(t *testing.T) {
// 	type fields struct {
// 		maxDist atomic.Value
// 		eg      errgroup.Group
// 		dch     chan DistPayload
// 		cancel  context.CancelFunc
// 		visited sync.Map[string, any]
// 		result  []*payload.Object_Distance
// 		num     int
// 		fnum    int
// 		closed  atomic.Bool
// 	}
// 	type want struct {
// 		want *payload.Search_Response
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Search_Response) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           maxDist:nil,
// 		           eg:nil,
// 		           dch:nil,
// 		           cancel:nil,
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           closed:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           maxDist:nil,
// 		           eg:nil,
// 		           dch:nil,
// 		           cancel:nil,
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           closed:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdStdAggr{
// 				maxDist: test.fields.maxDist,
// 				eg:      test.fields.eg,
// 				dch:     test.fields.dch,
// 				cancel:  test.fields.cancel,
// 				visited: test.fields.visited,
// 				result:  test.fields.result,
// 				num:     test.fields.num,
// 				fnum:    test.fields.fnum,
// 				closed:  test.fields.closed,
// 			}
//
// 			got := v.Result()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdStdAggr_GetNum(t *testing.T) {
// 	type fields struct {
// 		maxDist atomic.Value
// 		eg      errgroup.Group
// 		dch     chan DistPayload
// 		cancel  context.CancelFunc
// 		visited sync.Map[string, any]
// 		result  []*payload.Object_Distance
// 		num     int
// 		fnum    int
// 		closed  atomic.Bool
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           maxDist:nil,
// 		           eg:nil,
// 		           dch:nil,
// 		           cancel:nil,
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           closed:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           maxDist:nil,
// 		           eg:nil,
// 		           dch:nil,
// 		           cancel:nil,
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           closed:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdStdAggr{
// 				maxDist: test.fields.maxDist,
// 				eg:      test.fields.eg,
// 				dch:     test.fields.dch,
// 				cancel:  test.fields.cancel,
// 				visited: test.fields.visited,
// 				result:  test.fields.result,
// 				num:     test.fields.num,
// 				fnum:    test.fields.fnum,
// 				closed:  test.fields.closed,
// 			}
//
// 			got := v.GetNum()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdStdAggr_GetFnum(t *testing.T) {
// 	type fields struct {
// 		maxDist atomic.Value
// 		eg      errgroup.Group
// 		dch     chan DistPayload
// 		cancel  context.CancelFunc
// 		visited sync.Map[string, any]
// 		result  []*payload.Object_Distance
// 		num     int
// 		fnum    int
// 		closed  atomic.Bool
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           maxDist:nil,
// 		           eg:nil,
// 		           dch:nil,
// 		           cancel:nil,
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           closed:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           maxDist:nil,
// 		           eg:nil,
// 		           dch:nil,
// 		           cancel:nil,
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           closed:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdStdAggr{
// 				maxDist: test.fields.maxDist,
// 				eg:      test.fields.eg,
// 				dch:     test.fields.dch,
// 				cancel:  test.fields.cancel,
// 				visited: test.fields.visited,
// 				result:  test.fields.result,
// 				num:     test.fields.num,
// 				fnum:    test.fields.fnum,
// 				closed:  test.fields.closed,
// 			}
//
// 			got := v.GetFnum()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_newPairingHeap(t *testing.T) {
// 	type args struct {
// 		num  int
// 		fnum int
// 		in2  int
// 	}
// 	type want struct {
// 		want Aggregator
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Aggregator) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Aggregator) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           num:0,
// 		           fnum:0,
// 		           in2:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           num:0,
// 		           fnum:0,
// 		           in2:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := newPairingHeap(test.args.num, test.args.fnum, test.args.in2)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdPairingHeapAggr_Start(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		ph      *PairingHeap
// 		visited sync.Map[string, any]
// 		result  []*payload.Object_Distance
// 		num     int
// 		fnum    int
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           in0:nil,
// 		       },
// 		       fields: fields {
// 		           ph:PairingHeap{},
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           in0:nil,
// 		           },
// 		           fields: fields {
// 		           ph:PairingHeap{},
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdPairingHeapAggr{
// 				ph:      test.fields.ph,
// 				visited: test.fields.visited,
// 				result:  test.fields.result,
// 				num:     test.fields.num,
// 				fnum:    test.fields.fnum,
// 			}
//
// 			v.Start(test.args.in0)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdPairingHeapAggr_Send(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		data *payload.Search_Response
// 	}
// 	type fields struct {
// 		ph      *PairingHeap
// 		visited sync.Map[string, any]
// 		result  []*payload.Object_Distance
// 		num     int
// 		fnum    int
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           data:nil,
// 		       },
// 		       fields: fields {
// 		           ph:PairingHeap{},
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           data:nil,
// 		           },
// 		           fields: fields {
// 		           ph:PairingHeap{},
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdPairingHeapAggr{
// 				ph:      test.fields.ph,
// 				visited: test.fields.visited,
// 				result:  test.fields.result,
// 				num:     test.fields.num,
// 				fnum:    test.fields.fnum,
// 			}
//
// 			v.Send(test.args.ctx, test.args.data)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdPairingHeapAggr_Result(t *testing.T) {
// 	type fields struct {
// 		ph      *PairingHeap
// 		visited sync.Map[string, any]
// 		result  []*payload.Object_Distance
// 		num     int
// 		fnum    int
// 	}
// 	type want struct {
// 		want *payload.Search_Response
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *payload.Search_Response) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           ph:PairingHeap{},
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           ph:PairingHeap{},
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdPairingHeapAggr{
// 				ph:      test.fields.ph,
// 				visited: test.fields.visited,
// 				result:  test.fields.result,
// 				num:     test.fields.num,
// 				fnum:    test.fields.fnum,
// 			}
//
// 			got := v.Result()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdPairingHeapAggr_GetNum(t *testing.T) {
// 	type fields struct {
// 		ph      *PairingHeap
// 		visited sync.Map[string, any]
// 		result  []*payload.Object_Distance
// 		num     int
// 		fnum    int
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           ph:PairingHeap{},
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           ph:PairingHeap{},
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdPairingHeapAggr{
// 				ph:      test.fields.ph,
// 				visited: test.fields.visited,
// 				result:  test.fields.result,
// 				num:     test.fields.num,
// 				fnum:    test.fields.fnum,
// 			}
//
// 			got := v.GetNum()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdPairingHeapAggr_GetFnum(t *testing.T) {
// 	type fields struct {
// 		ph      *PairingHeap
// 		visited sync.Map[string, any]
// 		result  []*payload.Object_Distance
// 		num     int
// 		fnum    int
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           ph:PairingHeap{},
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           ph:PairingHeap{},
// 		           visited:nil,
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdPairingHeapAggr{
// 				ph:      test.fields.ph,
// 				visited: test.fields.visited,
// 				result:  test.fields.result,
// 				num:     test.fields.num,
// 				fnum:    test.fields.fnum,
// 			}
//
// 			got := v.GetFnum()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_newSlice(t *testing.T) {
// 	type args struct {
// 		num     int
// 		fnum    int
// 		replica int
// 	}
// 	type want struct {
// 		want Aggregator
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Aggregator) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Aggregator) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           num:0,
// 		           fnum:0,
// 		           replica:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           num:0,
// 		           fnum:0,
// 		           replica:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := newSlice(test.args.num, test.args.fnum, test.args.replica)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdSliceAggr_Start(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		result []*DistPayload
// 		num    int
// 		fnum   int
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           in0:nil,
// 		       },
// 		       fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           in0:nil,
// 		           },
// 		           fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdSliceAggr{
// 				result: test.fields.result,
// 				num:    test.fields.num,
// 				fnum:   test.fields.fnum,
// 			}
//
// 			v.Start(test.args.in0)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdSliceAggr_Send(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		data *payload.Search_Response
// 	}
// 	type fields struct {
// 		result []*DistPayload
// 		num    int
// 		fnum   int
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           data:nil,
// 		       },
// 		       fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           data:nil,
// 		           },
// 		           fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdSliceAggr{
// 				result: test.fields.result,
// 				num:    test.fields.num,
// 				fnum:   test.fields.fnum,
// 			}
//
// 			v.Send(test.args.ctx, test.args.data)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdSliceAggr_Result(t *testing.T) {
// 	type fields struct {
// 		result []*DistPayload
// 		num    int
// 		fnum   int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response) error {
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdSliceAggr{
// 				result: test.fields.result,
// 				num:    test.fields.num,
// 				fnum:   test.fields.fnum,
// 			}
//
// 			gotRes := v.Result()
// 			if err := checkFunc(test.want, gotRes); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdSliceAggr_GetNum(t *testing.T) {
// 	type fields struct {
// 		result []*DistPayload
// 		num    int
// 		fnum   int
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdSliceAggr{
// 				result: test.fields.result,
// 				num:    test.fields.num,
// 				fnum:   test.fields.fnum,
// 			}
//
// 			got := v.GetNum()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdSliceAggr_GetFnum(t *testing.T) {
// 	type fields struct {
// 		result []*DistPayload
// 		num    int
// 		fnum   int
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdSliceAggr{
// 				result: test.fields.result,
// 				num:    test.fields.num,
// 				fnum:   test.fields.fnum,
// 			}
//
// 			got := v.GetFnum()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_newPoolSlice(t *testing.T) {
// 	type args struct {
// 		num     int
// 		fnum    int
// 		replica int
// 	}
// 	type want struct {
// 		want Aggregator
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Aggregator) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Aggregator) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           num:0,
// 		           fnum:0,
// 		           replica:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           num:0,
// 		           fnum:0,
// 		           replica:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := newPoolSlice(test.args.num, test.args.fnum, test.args.replica)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdPoolSliceAggr_Start(t *testing.T) {
// 	type args struct {
// 		in0 context.Context
// 	}
// 	type fields struct {
// 		result []*DistPayload
// 		num    int
// 		fnum   int
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           in0:nil,
// 		       },
// 		       fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           in0:nil,
// 		           },
// 		           fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdPoolSliceAggr{
// 				result: test.fields.result,
// 				num:    test.fields.num,
// 				fnum:   test.fields.fnum,
// 			}
//
// 			v.Start(test.args.in0)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdPoolSliceAggr_Send(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		data *payload.Search_Response
// 	}
// 	type fields struct {
// 		result []*DistPayload
// 		num    int
// 		fnum   int
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           data:nil,
// 		       },
// 		       fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           data:nil,
// 		           },
// 		           fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdPoolSliceAggr{
// 				result: test.fields.result,
// 				num:    test.fields.num,
// 				fnum:   test.fields.fnum,
// 			}
//
// 			v.Send(test.args.ctx, test.args.data)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdPoolSliceAggr_Result(t *testing.T) {
// 	type fields struct {
// 		result []*DistPayload
// 		num    int
// 		fnum   int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response) error {
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdPoolSliceAggr{
// 				result: test.fields.result,
// 				num:    test.fields.num,
// 				fnum:   test.fields.fnum,
// 			}
//
// 			gotRes := v.Result()
// 			if err := checkFunc(test.want, gotRes); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdPoolSliceAggr_GetNum(t *testing.T) {
// 	type fields struct {
// 		result []*DistPayload
// 		num    int
// 		fnum   int
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdPoolSliceAggr{
// 				result: test.fields.result,
// 				num:    test.fields.num,
// 				fnum:   test.fields.fnum,
// 			}
//
// 			got := v.GetNum()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_valdPoolSliceAggr_GetFnum(t *testing.T) {
// 	type fields struct {
// 		result []*DistPayload
// 		num    int
// 		fnum   int
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           result:nil,
// 		           num:0,
// 		           fnum:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			v := &valdPoolSliceAggr{
// 				result: test.fields.result,
// 				num:    test.fields.num,
// 				fnum:   test.fields.fnum,
// 			}
//
// 			got := v.GetFnum()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
