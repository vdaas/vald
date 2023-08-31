// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package singleflight

import (
	"context"
	"fmt"
	"io/fs"
	"math"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/singleflight"
	stdsingleflight "golang.org/x/sync/singleflight"
)

type Result struct {
	Goroutine int     `csv:"goroutine"`
	Duration  int64   `csv:"duration"`
	HitRate   float64 `csv:"hit_rate"`
}

type helper struct {
	initDoFn  func() func(ctx context.Context, key string, fn func(context.Context) (string, error))
	sleepDur  time.Duration
	calledCnt int64
	totalCnt  int64
}

const (
	minGoroutine  = 10
	maxGoroutine  = 10000
	goroutineStep = 10
	tryCnt        = 5
)

var durs = []time.Duration{
	time.Microsecond * 10,
	time.Microsecond * 100,
	time.Microsecond * 200,
	time.Microsecond * 500,
	time.Millisecond,
	time.Millisecond * 5,
	time.Millisecond * 10,
	time.Millisecond * 25,
	time.Millisecond * 50,
	time.Millisecond * 100,
	time.Millisecond * 250,
	time.Millisecond * 500,
}

func (h *helper) Do(parallel int, b *testing.B) {
	b.Helper()

	fn := func(context.Context) (string, error) {
		atomic.AddInt64(&h.calledCnt, 1)
		time.Sleep(h.sleepDur)
		return "", nil
	}

	doFn := h.initDoFn()

	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
		atomic.AddInt64(&h.calledCnt, -1)
		doFn(context.Background(), "key", fn)
	}()
	<-ch
	close(ch)

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()

	b.SetParallelism(parallel)
	b.RunParallel(func(pb *testing.PB) {
		var wg sync.WaitGroup
		for pb.Next() {
			wg.Add(1)
			atomic.AddInt64(&h.totalCnt, 1)
			go func() {
				defer wg.Done()
				doFn(context.Background(), "key", fn)
			}()
		}
		wg.Wait()
	})
}

func Benchmark_group_Do_with_sync_singleflight(b *testing.B) {
	const (
		varianceCSV = "sync_singleflight_variance.csv"
		averageCSV  = "sync_singleflight_average.csv"
	)
	resultsmap := make(map[string][]Result)
	for i := minGoroutine; i <= maxGoroutine; i *= goroutineStep {
		for _, dur := range durs {
			results := make([]Result, 0, tryCnt)
			for j := 0; j < tryCnt; j++ {
				h := &helper{
					initDoFn: func() func(ctx context.Context, key string, fn func(context.Context) (string, error)) {
						g := new(stdsingleflight.Group)
						return func(ctx context.Context, key string, fn func(context.Context) (string, error)) {
							g.Do(key, func() (interface{}, error) { return fn(context.Background()) })
						}
					},
					sleepDur: dur,
				}

				b.Helper()
				b.StopTimer()
				b.ReportAllocs()
				b.ResetTimer()
				b.StartTimer()

				b.Run(fmt.Sprintf("%d %s", i, dur), func(b *testing.B) {
					h.Do(i, b)
				})
				b.StopTimer()

				hitCnt := h.totalCnt - h.calledCnt
				hitRate := float64(hitCnt) / float64(h.totalCnt)

				b.Logf("Parallel: %d\tTotal Goroutine Count: %d\tHit Count: %d\tHit Rate: %f",
					i,
					h.totalCnt,
					hitCnt,
					hitRate,
				)
				results = append(results, Result{
					Goroutine: i,
					Duration:  dur.Nanoseconds(),
					HitRate:   hitRate,
				})
			}

			resultsmap[averageCSV] = append(
				resultsmap[averageCSV],
				calcAverage(results),
			)
			resultsmap[varianceCSV] = append(
				resultsmap[varianceCSV],
				calcVariance(results),
			)
		}
	}

	for name, results := range resultsmap {
		if err := toCSV(name, results); err != nil {
			b.Error(err)
		}
	}
}

func calcAverage(in []Result) (out Result) {
	var sum float64
	for i, r := range in {
		if i == 0 {
			out.Goroutine = r.Goroutine
			out.Duration = r.Duration
		}
		sum += r.HitRate
	}
	out.HitRate = sum / float64(len(in))

	return
}

func calcVariance(in []Result) (out Result) {
	aveResult := calcAverage(in)

	var sum float64
	for i, r := range in {
		if i == 0 {
			out.Goroutine = r.Goroutine
			out.Duration = r.Duration
		}
		sum += math.Pow(r.HitRate-aveResult.HitRate, 2)
	}
	out.HitRate = sum / float64(len(in))

	return
}

func Benchmark_group_Do_with_vald_internal_singleflight(b *testing.B) {
	const (
		varianceCSV = "vald_internal_singlefligh_variance.csv"
		averageCSV  = "vald_internal_singlefligh_average.csv"
	)
	resultsmap := make(map[string][]Result)
	for i := minGoroutine; i <= maxGoroutine; i *= goroutineStep {
		for _, dur := range durs {
			results := make([]Result, 0, tryCnt)
			for j := 0; j < tryCnt; j++ {
				h := &helper{
					initDoFn: func() func(ctx context.Context, key string, fn func(context.Context) (string, error)) {
						g := singleflight.New[string]()
						return func(ctx context.Context, key string, fn func(context.Context) (string, error)) {
							g.Do(ctx, key, fn)
						}
					},
					sleepDur: dur,
				}

				b.StopTimer()
				b.ReportAllocs()
				b.ResetTimer()
				b.StartTimer()

				b.Run(fmt.Sprintf("%d %s", i, dur), func(b *testing.B) {
					h.Do(i, b)
				})
				b.StopTimer()

				hitCnt := h.totalCnt - h.calledCnt
				hitRate := float64(hitCnt) / float64(h.totalCnt)

				b.Logf("Parallel: %d\tTotal Goroutine Count: %d\tHit Count: %d\tHit Rate: %f",
					i,
					h.totalCnt,
					hitCnt,
					hitRate,
				)
				results = append(results, Result{
					Goroutine: i,
					Duration:  dur.Nanoseconds(),
					HitRate:   hitRate,
				})
			}

			resultsmap[averageCSV] = append(
				resultsmap[averageCSV],
				calcAverage(results),
			)
			resultsmap[varianceCSV] = append(
				resultsmap[varianceCSV],
				calcVariance(results),
			)
		}
	}
	for name, results := range resultsmap {
		if err := toCSV(name, results); err != nil {
			b.Error(err)
		}
	}
}

func toCSV(name string, r []Result) (err error) {
	var f *os.File
	f, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		e := f.Close()
		if e != nil {
			err = errors.Wrap(err, e.Error())
		}
	}()
	_, err = fmt.Fprintln(f, "goroutine,duration,hit_rate")
	if err != nil {
		return err
	}
	for _, res := range r {
		_, err = fmt.Fprintf(f, "%d,%v,%f\n", res.Goroutine, res.Duration, res.HitRate)
		if err != nil {
			return err
		}
	}
	return nil
}
