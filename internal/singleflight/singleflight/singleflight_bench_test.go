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
	g        Group
	sleepDur time.Duration
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

func (h *helper) Do(parallel int, b *testing.B) Result {
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

	b.StopTimer()

	b.Logf("Parallel: %d\tTotal Goroutine Count: %d\tHit Count: %d\tHit Rate: %f",
		parallel,
		totalCnt,
		hitCnt,
		float64(hitCnt)/float64(totalCnt),
	)
	return Result{
		Goroutine: parallel,
		HitRate:   float64(hitCnt) / float64(totalCnt),
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

			b.StopTimer()
			b.ReportAllocs()
			b.ResetTimer()
			b.StartTimer()

			b.Run(fmt.Sprintf("%d %s", i, dur), func(b *testing.B) {
				res := h.Do(i, b)
				res.Duration = dur
				results = append(results, res)
			})
			b.StopTimer()
		}
	}
	b.StopTimer()
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
				res := h.Do(i, b)
				res.Duration = dur
				results = append(results, res)
			})
			b.StopTimer()
		}
	}
	b.StopTimer()
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
