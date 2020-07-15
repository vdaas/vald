package singleflight

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/singleflight"
)

type Result struct {
	Goroutine int           `csv:"goroutine"`
	Duration  time.Duration `csv:"duration"`
	HitRate   float64       `csv:"hit_rate"`
}

type helper struct {
	g         Group
	sleepDur  time.Duration
	calledCnt int64
	totalCnt  int64
}

const (
	minGoroutine  = 10
	maxGoroutine  = 100000
	goroutineStep = 10
)

var (
	durs = []time.Duration{
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
)

func (h *helper) Do(parallel int, b *testing.B) {
	b.Helper()

	var (
		fn = func() (interface{}, error) {
			atomic.AddInt64(&h.calledCnt, 1)
			time.Sleep(h.sleepDur)
			return "", nil
		}
	)

	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
		atomic.AddInt64(&h.calledCnt, -1)
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
			atomic.AddInt64(&h.totalCnt, 1)
			go func() {
				defer wg.Done()
				h.g.Do(context.Background(), "key", fn)
			}()
		}
		wg.Wait()
	})
}

func Benchmark_group_Do_with_mutex(b *testing.B) {
	results := make([]Result, 0, len(durs)*maxGoroutine-minGoroutine/goroutineStep)
	for i := minGoroutine; i <= maxGoroutine; i *= goroutineStep {
		for _, dur := range durs {
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
				Duration:  dur,
				HitRate:   hitRate,
			})
		}
	}
}

func Benchmark_group_Do_with_mutex_1(b *testing.B) {
	results := make([]Result, 0, len(durs)*maxGoroutine-minGoroutine/goroutineStep)
	for i := minGoroutine; i <= maxGoroutine; i *= goroutineStep {
		for _, dur := range durs {
			h := &helper{
				g:        singleflight.New(10),
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
				Duration:  dur,
				HitRate:   hitRate,
			})
		}
	}
	toCSV("mutex.csv", results)
}

func Benchmark_group_Do_with_syncMap(b *testing.B) {
	results := make([]Result, 0, len(durs)*maxGoroutine-minGoroutine/goroutineStep)
	for i := minGoroutine; i <= maxGoroutine; i *= goroutineStep {
		for _, dur := range durs {
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
				Duration:  dur,
				HitRate:   hitRate,
			})
		}
	}
	toCSV("syncmap.csv", results)
}

func toCSV(name string, r []Result) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprint(f, "goroutine,duration,hit_rate")
	for _, res := range r {
		_, err = fmt.Fprintf(f, "%d,%v,%f\n", res.Goroutine, res.Duration, res.HitRate)
		if err != nil {
			return err
		}
	}
	return nil
}
