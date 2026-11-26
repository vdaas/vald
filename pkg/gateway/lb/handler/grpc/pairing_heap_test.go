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
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

func TestPairingHeapNodes_EdgeCasesAndOrdering(t *testing.T) {
	t.Parallel()

	if got := mergeMinNodes(nil, pairingNull, 0); got != 0 {
		t.Errorf("mergeMinNodes(pairingNull, 0) = %d, want 0", got)
	}
	if got := mergeMinNodes(nil, 2, pairingNull); got != 2 {
		t.Errorf("mergeMinNodes(2, pairingNull) = %d, want 2", got)
	}
	if got := mergeMaxNodes(nil, pairingNull, 1); got != 1 {
		t.Errorf("mergeMaxNodes(pairingNull, 1) = %d, want 1", got)
	}
	if got := mergeMaxNodes(nil, 3, pairingNull); got != 3 {
		t.Errorf("mergeMaxNodes(3, pairingNull) = %d, want 3", got)
	}
	if got := mergePairsMinNodes(nil, pairingNull); got != pairingNull {
		t.Errorf("mergePairsMinNodes(pairingNull) = %d, want %d", got, pairingNull)
	}
	if got := mergePairsMaxNodes(nil, pairingNull); got != pairingNull {
		t.Errorf("mergePairsMaxNodes(pairingNull) = %d, want %d", got, pairingNull)
	}

	for _, count := range []int{1, 2, 3, 1000} {
		t.Run(fmt.Sprintf("extract_count_%d", count), func(t *testing.T) {
			t.Parallel()
			rng := rand.New(rand.NewPCG(42, uint64(count))) //nolint:gosec // deterministic test seed
			minNodes := make([]pairingNode, count)
			maxNodes := make([]pairingNode, count)
			minRoot := pairingNull
			maxRoot := pairingNull
			wantAsc := make([]float64, count)
			for i := range count {
				d := rng.Float64()
				wantAsc[i] = d
				dp := DistPayload{
					raw:      dist32(fmt.Sprintf("id-%d", i), float32(d)),
					distance: d,
				}
				minNodes[i] = pairingNode{dp: dp, child: pairingNull, sibling: pairingNull}
				maxNodes[i] = pairingNode{dp: dp, child: pairingNull, sibling: pairingNull}
				idx := int32(i) //nolint:gosec // test count bounded to 1000
				minRoot = mergeMinNodes(minNodes, minRoot, idx)
				maxRoot = mergeMaxNodes(maxNodes, maxRoot, idx)
			}
			slices.Sort(wantAsc)

			gotMin := make([]float64, 0, count)
			for minRoot != pairingNull {
				gotMin = append(gotMin, minNodes[minRoot].dp.distance)
				minRoot = mergePairsMinNodes(minNodes, minNodes[minRoot].child)
			}
			if !slices.Equal(gotMin, wantAsc) {
				t.Errorf("min-heap sequence mismatch for count=%d", count)
			}

			gotMax := make([]float64, 0, count)
			for maxRoot != pairingNull {
				gotMax = append(gotMax, maxNodes[maxRoot].dp.distance)
				maxRoot = mergePairsMaxNodes(maxNodes, maxNodes[maxRoot].child)
			}
			slices.Reverse(gotMax)
			if !slices.Equal(gotMax, wantAsc) {
				t.Errorf("max-heap sequence mismatch for count=%d", count)
			}
		})
	}
}

func TestPairingHeapAggr_UnboundedSort(t *testing.T) {
	t.Parallel()
	const count = 100
	rng := rand.New(rand.NewPCG(7, count)) //nolint:gosec // deterministic test seed
	ds := make([]*payload.Object_Distance, count)
	want := make([]float32, count)
	for i := range count {
		d := rng.Float32()
		ds[i] = dist32(fmt.Sprintf("id-%d", i), d)
		want[i] = d
	}
	slices.Sort(want)

	aggr := newPairingHeap(0, 0, 1)
	aggr.Send(t.Context(), newDistResponse(ds...))
	res := aggr.Result().GetResults()
	if len(res) != count {
		t.Fatalf("len(results) = %d, want %d", len(res), count)
	}
	for i, r := range res {
		if r.GetDistance() != want[i] {
			t.Fatalf("res[%d].Distance = %v, want %v", i, r.GetDistance(), want[i])
		}
	}
}
