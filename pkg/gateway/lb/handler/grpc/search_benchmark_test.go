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
	"fmt"
	"math/rand/v2"
	"runtime"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/bench"
	"github.com/vdaas/vald/internal/test/data/strings"
)

const (
	fastDataVariation = 64
	fullDataVariation = 10000
	dataLength        = 10000
)

var (
	datas  []*payload.Search_Response
	delays []time.Duration
	//nolint:gochecknoglobals // lazy benchmark dataset cache
	benchDataOnce sync.Once
)

func ensureBenchData(full bool) {
	benchDataOnce.Do(func() {
		variation := fastDataVariation
		if full {
			variation = fullDataVariation
		}
		var mu sync.Mutex
		var wg sync.WaitGroup
		datas = make([]*payload.Search_Response, 0, variation)
		delays = make([]time.Duration, 0, variation)
		for range variation {
			wg.Go(func() {
				rr := newRandomResponse()
				delay := time.Duration(rand.N(uint32(time.Millisecond * 2))) //nolint:gosec // synthetic benchmark jitter does not require crypto/rand
				mu.Lock()
				datas = append(datas, rr)
				delays = append(delays, delay)
				mu.Unlock()
			})
		}
		wg.Wait()
	})
}

func newRandomResponse() (res *payload.Search_Response) {
	res = &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, dataLength),
	}
	for range dataLength {
		res.Results = append(res.Results, &payload.Object_Distance{
			Id:       strings.Random(20),
			Distance: rand.Float32(),
		})
	}
	slices.SortFunc(res.Results, func(left, right *payload.Object_Distance) int {
		return cmp.Compare(left.GetDistance(), right.GetDistance())
	})
	return res
}

func fastParallelisms() []int {
	pMax := min(8, runtime.NumCPU())
	if pMax <= 2 {
		return []int{2}
	}
	return []int{2, pMax}
}

func benchmark(b *testing.B, anew func(n, f, r int) Aggregator) {
	b.Helper()
	full := bench.Full()
	ensureBenchData(full)
	results := datas
	l := len(results)
	ks := []int{10, 100}
	replicas := []int{5, 50}
	var parallelisms []int
	if full {
		ks = []int{10, 100, 1000}
		replicas = []int{5, 50, 500, 5000}
		for parallelism := 2; parallelism <= runtime.NumCPU(); parallelism *= 2 {
			parallelisms = append(parallelisms, parallelism)
		}
	} else {
		parallelisms = fastParallelisms()
	}
	for _, k := range ks {
		for _, replica := range replicas {
			for _, parallelism := range parallelisms {
				b.Run(fmt.Sprintf("Top-K:%d/Replica:%d/Parallelism:%d/Thread:", k, replica, parallelism), func(bb *testing.B) {
					bb.Helper()
					bb.SetParallelism(parallelism)
					bb.ReportAllocs()
					bb.ResetTimer()
					bb.RunParallel(func(pb *testing.PB) {
						for pb.Next() {
							var cnt atomic.Uint64
							if _, err := doSearchWithAggregator(bb.Context(), k, replica, anew, func(ctx context.Context) (res *payload.Search_Response) {
								idx := cnt.Add(1) % uint64(l)
								time.Sleep(delays[idx])
								return results[idx]
							}); err != nil {
								bb.Errorf("doSearchWithAggregator() error = %v", err)
							}
						}
					})
				})
			}
		}
	}
}

func doSearchWithAggregator(
	ctx context.Context,
	k, concurrency int,
	anew func(n, f, r int) Aggregator,
	f func(ctx context.Context) *payload.Search_Response,
) (res *payload.Search_Response, err error) {
	eg, ectx := errgroup.New(ctx)
	eg.SetLimit(concurrency)
	aggr := anew(k, k, concurrency)
	for range concurrency {
		eg.Go(func() error {
			r := f(ectx)
			if r != nil && len(r.GetResults()) != 0 {
				aggr.Send(ectx, r)
			}
			return nil
		})
	}
	err = eg.Wait()
	return aggr.Result(), err
}

func BenchmarkStandard(b *testing.B) {
	benchmark(b, newStd)
}

func BenchmarkPairingHeap(b *testing.B) {
	benchmark(b, newPairingHeap)
}

func BenchmarkSortSlice(b *testing.B) {
	benchmark(b, newSlice)
}

func BenchmarkSortPoolSlice(b *testing.B) {
	benchmark(b, newPoolSlice)
}
