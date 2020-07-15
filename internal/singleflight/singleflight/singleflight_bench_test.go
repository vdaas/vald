package singleflight

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/singleflight"
)

type helper struct {
	g        Group
	sleepDur time.Duration
}

func (h *helper) Do(parallel int, b *testing.B) {
	b.Helper()

	var (
		calledCnt, totalCnt int64

		fn = func() (interface{}, error) {
			time.Sleep(h.sleepDur)
			atomic.AddInt64(&calledCnt, 1)
			return "", nil
		}
	)

	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
		h.g.Do(context.Background(), "key", fn)
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
			atomic.AddInt64(&totalCnt, 1)
			go func() {
				defer wg.Done()
				h.g.Do(context.Background(), "key", fn)
			}()
		}
		wg.Wait()
	})

	hitCnt := totalCnt - calledCnt

	b.Logf("Parallel: %d\tTotal Goroutine Count: %d\tHit Count: %d\tHit Rate: %f",
		parallel,
		totalCnt,
		hitCnt,
		float64(hitCnt)/float64(totalCnt),
	)
}

func Benchmark_group_Do_with_mutex_1(b *testing.B) {
	dursFn := func() (durs []time.Duration) {
		for i := 10; i <= 1000; i *= 10 {
			durs = append(durs, time.Duration(i)*time.Microsecond)
		}
		return
	}
	for i := 100; i <= 10000; i *= 10 {
		for _, dur := range dursFn() {
			h := &helper{
				g:        singleflight.New(10),
				sleepDur: dur,
			}

			b.StopTimer()
			b.ReportAllocs()
			b.ResetTimer()
			b.StartTimer()

			b.Run(fmt.Sprintf("%d %s", i, dur), func(b *testing.B) {
				h.Do(i, b)
			})
		}
	}
}

func Benchmark_group_Do_with_syncMap(b *testing.B) {
	dursFn := func() (durs []time.Duration) {
		for i := 10; i <= 1000; i *= 10 {
			durs = append(durs, time.Duration(i)*time.Microsecond)
		}
		return
	}
	for i := 100; i <= 10000; i *= 10 {
		for _, dur := range dursFn() {
			h := &helper{
				g:        New(),
				sleepDur: dur,
			}

			b.StopTimer()
			b.ReportAllocs()
			b.ResetTimer()
			b.StartTimer()

			b.Run(fmt.Sprintf("%d %s", i, dur), func(b *testing.B) {
				h.Do(i, b)
			})
		}
	}
}
