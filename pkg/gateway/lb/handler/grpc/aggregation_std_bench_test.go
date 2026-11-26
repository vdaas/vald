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
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

// BenchmarkStdAggregator_Agents aggregates one search (top-10, each agent
// answering its sorted top-10) as aggregationSearch does with the default Std
// aggregator, for 10, 50 and 200 agents, without network or sleep noise.
func BenchmarkStdAggregator_Agents(b *testing.B) {
	const num = 10
	for _, agents := range []int{10, 50, 200} {
		r := rand.New(rand.NewPCG(1, uint64(agents))) //nolint:gosec // deterministic benchmark data, not security sensitive
		resps := make([]*payload.Search_Response, agents)
		for a := range resps {
			res := make([]*payload.Object_Distance, num)
			for i := range res {
				res[i] = &payload.Object_Distance{Id: strconv.Itoa(a*num + i), Distance: r.Float32()}
			}
			slices.SortFunc(res, func(x, y *payload.Object_Distance) int { return cmp.Compare(x.GetDistance(), y.GetDistance()) })
			resps[a] = &payload.Search_Response{Results: res}
		}
		b.Run("agents="+strconv.Itoa(agents), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				aggr := newStd(num, num, agents)
				for _, res := range resps {
					aggr.Send(b.Context(), res)
				}
				if got := len(aggr.Result().GetResults()); got != num {
					b.Fatalf("Result() returned %d results, want %d", got, num)
				}
			}
		})
	}
}

// BenchmarkAggregators_Agents compares all 4 gateway aggregators under realistic
// replica=3 overlapping IDs for 5, 10, and 50 agents.
func BenchmarkAggregators_Agents(b *testing.B) {
	const (
		num     = 10
		replica = 3
	)
	algos := []struct {
		newFn func(num, fnum, replica int) Aggregator
		name  string
	}{
		{name: "Std", newFn: newStd},
		{name: "Slice", newFn: newSlice},
		{name: "PoolSlice", newFn: newPoolSlice},
		{name: "PairingHeap", newFn: newPairingHeap},
	}
	for _, agents := range []int{5, 10, 50} {
		r := rand.New(rand.NewPCG(7, uint64(agents))) //nolint:gosec // deterministic benchmark data, not security sensitive
		totalUnique := agents * num / replica
		pool := make([]*payload.Object_Distance, totalUnique)
		for i := range pool {
			pool[i] = &payload.Object_Distance{Id: "id-" + strconv.Itoa(i), Distance: r.Float32()}
		}
		resps := make([]*payload.Search_Response, agents)
		for a := range resps {
			res := make([]*payload.Object_Distance, num)
			for i := range res {
				res[i] = pool[(a*num/replica+i)%totalUnique]
			}
			slices.SortFunc(res, func(x, y *payload.Object_Distance) int { return cmp.Compare(x.GetDistance(), y.GetDistance()) })
			resps[a] = &payload.Search_Response{Results: res}
		}
		for _, algo := range algos {
			b.Run("algo="+algo.name+"/agents="+strconv.Itoa(agents), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					aggr := algo.newFn(num, num, replica)
					for _, res := range resps {
						aggr.Send(b.Context(), res)
					}
					if got := len(aggr.Result().GetResults()); got != num {
						b.Fatalf("%s Result() returned %d results, want %d", algo.name, got, num)
					}
				}
			})
		}
	}
}
