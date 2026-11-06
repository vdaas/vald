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
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
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
