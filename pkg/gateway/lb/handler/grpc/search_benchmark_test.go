// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"flag"
	"fmt"
	"runtime"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/rand"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/data/strings"
)

const (
	dataVariation = 10000
	dataLength    = 10000
)

var (
	datas  []*payload.Search_Response
	delays []time.Duration
)

func TestMain(m *testing.M) {
	testing.Init()
	flag.Parse()
	if testing.Short() {
		m.Run()
		return
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	datas = make([]*payload.Search_Response, 0, dataVariation)
	delays = make([]time.Duration, 0, dataVariation)
	for i := 0; i < dataVariation; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rr := newRandomResponse()
			delay := time.Duration(rand.LimitedUint32(uint64(time.Millisecond * 2)))
			mu.Lock()
			datas = append(datas, rr)
			delays = append(delays, delay)
			mu.Unlock()
		}()
	}
	wg.Wait()
	m.Run()
	datas = nil
	delays = nil
}

func newRandomResponse() (res *payload.Search_Response) {
	res = &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, dataLength),
	}
	for i := 0; i < dataLength; i++ {
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

func benchmark(b *testing.B, results []*payload.Search_Response, anew func(n, r int) Aggregator) {
	ctx := context.Background()
	l := len(results)
	for k := 10; k < dataLength; k *= 10 {
		for replica := 5; replica <= dataVariation; replica *= 10 {
			for parallelism := 2; parallelism <= runtime.NumCPU(); parallelism *= 2 {
				b.Run(fmt.Sprintf("Top-K:%d/Replica:%d/Parallelism:%d/Thread:", k, replica, parallelism), func(bb *testing.B) {
					bb.Helper()
					bb.SetParallelism(parallelism)
					bb.ReportAllocs()
					bb.ResetTimer()
					bb.RunParallel(func(pb *testing.PB) {
						for pb.Next() {
							var cnt atomic.Uint64
							_, _ = doSearchWithAggregator(ctx, k, replica, anew, func(ctx context.Context) (res *payload.Search_Response) {
								idx := cnt.Add(1) % uint64(l)
								time.Sleep(delays[idx])
								return results[idx]
							})
						}
					})
				})
			}
		}
	}
}

func doSearchWithAggregator(ctx context.Context, k, concurrency int, anew func(n, r int) Aggregator,
	f func(ctx context.Context) *payload.Search_Response,
) (res *payload.Search_Response, err error) {
	eg, ectx := errgroup.New(ctx)
	eg.SetLimit(concurrency)
	aggr := anew(k, concurrency)
	aggr.Start(ectx)
	for i := 0; i < concurrency; i++ {
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
	benchmark(b, datas, newStd)
}

func BenchmarkPairingHeap(b *testing.B) {
	benchmark(b, datas, newPairingHeap)
}

func BenchmarkSortSlice(b *testing.B) {
	benchmark(b, datas, newSlice)
}

func BenchmarkSortPoolSlice(b *testing.B) {
	benchmark(b, datas, newPoolSlice)
}
