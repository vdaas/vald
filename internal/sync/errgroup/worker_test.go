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

package errgroup

import (
	"context"
	"math/bits"
	"strconv"
	"testing"
	"testing/synctest"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
)

type foreignGroup struct {
	calls atomic.Int64
}

func (f *foreignGroup) Go(fn func() error) {
	f.calls.Add(1)
	go func() {
		if err := fn(); err != nil {
			panic(err)
		}
	}()
}

func (*foreignGroup) SetLimit(int) {}

func (f *foreignGroup) TryGo(fn func() error) bool {
	f.Go(fn)
	return true
}

func (*foreignGroup) Wait() error { return nil }

func withIsolatedGlobalGroup(t *testing.T, fn func(g Group)) {
	t.Helper()
	origInst := instance
	instance = nil
	once = sync.Once{}
	defer func() {
		instance = origInst
		once = sync.Once{}
		if origInst != nil {
			once.Do(func() {})
		}
	}()
	fn(Get())
}

func idleCount(p *WorkerPool) int {
	total := 0
	for i := range p.shards {
		total += bits.OnesCount64(p.shards[i].idleMask.Load())
	}
	return total
}

func runningCount(p *WorkerPool) int {
	total := 0
	for i := range p.shards {
		total += bits.OnesCount64(p.shards[i].runningMask.Load())
	}
	return total
}

func snapshotIdleMasks(p *WorkerPool) [poolShards]uint64 {
	var masks [poolShards]uint64
	for i := range p.shards {
		masks[i] = p.shards[i].idleMask.Load()
	}
	return masks
}

func snapshotRunningMasks(p *WorkerPool) [poolShards]uint64 {
	var masks [poolShards]uint64
	for i := range p.shards {
		masks[i] = p.shards[i].runningMask.Load()
	}
	return masks
}

func warmParkedWorkers(t *testing.T, g Group, p *WorkerPool, workers int) {
	t.Helper()
	arrived := make(chan struct{}, workers)
	release := make(chan struct{})
	if err := p.FanOut(g, workers+1, func(i int) {
		if i < workers {
			arrived <- struct{}{}
			<-release
			return
		}
		for range workers {
			<-arrived
		}
		close(release)
	}); err != nil {
		t.Fatalf("warm FanOut(%d) error = %v", workers+1, err)
	}
	synctest.Wait()
}

func TestWorkerPool_FanOut(t *testing.T) {
	t.Run("noop edge cases", func(t *testing.T) {
		withIsolatedGlobalGroup(t, func(g Group) {
			var p WorkerPool
			defer p.Close()

			if err := p.FanOut(g, 0, func(int) { t.Fatal("unexpected call for n=0") }); err != nil {
				t.Fatalf("FanOut(0) error = %v", err)
			}
			if err := p.FanOut(g, -3, func(int) { t.Fatal("unexpected call for n=-3") }); err != nil {
				t.Fatalf("FanOut(-3) error = %v", err)
			}
			if err := p.FanOut(g, 4, nil); err != nil {
				t.Fatalf("FanOut(nil) error = %v", err)
			}

			var called int
			if err := p.FanOut(g, 1, func(i int) {
				called++
				if i != 0 {
					t.Errorf("index = %d, want 0", i)
				}
			}); err != nil {
				t.Fatalf("FanOut(1) error = %v", err)
			}
			if called != 1 {
				t.Fatalf("called = %d, want 1", called)
			}
		})
	})

	for _, n := range []int{2, 5, 16, 50} {
		t.Run("pooled_completeness/n="+strconv.Itoa(n), func(t *testing.T) {
			withIsolatedGlobalGroup(t, func(g Group) {
				var p WorkerPool
				defer p.Close()

				seen := make([]atomic.Int64, n)
				if err := p.FanOut(g, n, func(i int) {
					seen[i].Add(1)
				}); err != nil {
					t.Fatalf("FanOut(%d) error = %v", n, err)
				}
				for i := range seen {
					if got := seen[i].Load(); got != 1 {
						t.Errorf("n=%d index %d count = %d, want 1", n, i, got)
					}
				}
			})
		})
	}

	t.Run("overflows beyond 512 pooled slots cleanly", func(t *testing.T) {
		withIsolatedGlobalGroup(t, func(g Group) {
			var p WorkerPool
			defer p.Close()

			const total = poolShards*slotsPerShard + 22
			seen := make([]atomic.Int64, total)
			if err := p.FanOut(g, total, func(i int) {
				seen[i].Add(1)
			}); err != nil {
				t.Fatalf("FanOut(%d) error = %v", total, err)
			}
			for i := range seen {
				if got := seen[i].Load(); got != 1 {
					t.Errorf("overflow index %d count = %d, want 1", i, got)
				}
			}

			// Saturate all 512 pooled slots and verify a concurrent direct-dispatch batch overflows cleanly.
			const capWorkers = poolShards * slotsPerShard
			arrived := make(chan struct{}, capWorkers)
			release := make(chan struct{})
			saturatedDone := make(chan error, 1)
			go func() {
				saturatedDone <- p.FanOut(g, capWorkers+1, func(i int) {
					if i < capWorkers {
						arrived <- struct{}{}
						<-release
					}
				})
			}()
			for range capWorkers {
				select {
				case <-arrived:
				case <-time.After(5 * time.Second):
					close(release)
					t.Fatal("timed out waiting to saturate worker pool")
				}
			}

			directSeen := make([]atomic.Int64, 5)
			if err := p.FanOut(g, len(directSeen), func(i int) {
				directSeen[i].Add(1)
			}); err != nil {
				close(release)
				t.Fatalf("direct overflow FanOut error = %v", err)
			}
			close(release)
			if err := <-saturatedDone; err != nil {
				t.Fatalf("saturated FanOut error = %v", err)
			}
			for i := range directSeen {
				if got := directSeen[i].Load(); got != 1 {
					t.Errorf("direct overflow index %d count = %d, want 1", i, got)
				}
			}
		})
	})

	t.Run("async workers run concurrently with caller inline last index", func(t *testing.T) {
		withIsolatedGlobalGroup(t, func(g Group) {
			const n = 5
			var p WorkerPool
			defer p.Close()

			arrived := make(chan int, n-1)
			release := make(chan struct{})
			if err := p.FanOut(g, n, func(i int) {
				if i < n-1 {
					arrived <- i
					<-release
					return
				}
				for range n - 1 {
					select {
					case <-arrived:
					case <-time.After(5 * time.Second):
						close(release)
						t.Fatal("timed out waiting for async workers")
					}
				}
				close(release)
			}); err != nil {
				t.Fatalf("FanOut(%d) error = %v", n, err)
			}
		})
	})
}

func TestWorkerPool_Synctest(t *testing.T) {
	t.Run("reuses idle workers in lowest-slot bitmask order", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			withIsolatedGlobalGroup(t, func(g Group) {
				var p WorkerPool
				defer func() {
					p.Close()
					synctest.Wait()
				}()

				const workers = 16
				warmParkedWorkers(t, g, &p, workers)

				firstIdle := idleCount(&p)
				if firstIdle != workers {
					t.Fatalf("idle workers = %d, want %d", firstIdle, workers)
				}
				if firstRunning := runningCount(&p); firstRunning != workers {
					t.Fatalf("running workers = %d, want %d", firstRunning, workers)
				}
				beforeIdle := snapshotIdleMasks(&p)
				beforeRunning := snapshotRunningMasks(&p)
				hasMultiSlotShard := false
				for _, m := range beforeIdle {
					if bits.OnesCount64(m) >= 2 {
						hasMultiSlotShard = true
						break
					}
				}
				if !hasMultiSlotShard {
					t.Fatal("expected at least one shard with depth >= 2")
				}

				poppedShard := -1
				poppedSlot := -1
				if err := p.FanOut(g, 2, func(i int) {
					if i == 0 {
						if got := idleCount(&p); got != firstIdle-1 {
							t.Errorf("idle workers during 1-worker batch = %d, want %d", got, firstIdle-1)
						}
						if got := runningCount(&p); got != workers {
							t.Errorf("running workers during reuse = %d, want %d", got, workers)
						}
						during := snapshotIdleMasks(&p)
						for s := range during {
							diff := beforeIdle[s] &^ during[s]
							if diff != 0 {
								poppedShard = s
								poppedSlot = bits.TrailingZeros64(diff)
								break
							}
						}
					}
				}); err != nil {
					t.Fatalf("second FanOut error = %v", err)
				}
				synctest.Wait()

				if got := idleCount(&p); got != firstIdle {
					t.Fatalf("idle workers after reuse = %d, want %d", got, firstIdle)
				}
				if got := runningCount(&p); got != workers {
					t.Fatalf("running workers after reuse = %d, want %d", got, workers)
				}
				if poppedShard < 0 {
					t.Fatal("expected one existing idle slot to be popped")
				}
				if wantSlot := bits.TrailingZeros64(beforeIdle[poppedShard]); poppedSlot != wantSlot {
					t.Errorf("popped slot %d != lowest idle slot %d in shard %d", poppedSlot, wantSlot, poppedShard)
				}
				if afterIdle := snapshotIdleMasks(&p); afterIdle != beforeIdle {
					t.Errorf("after return idle masks = %v, want %v", afterIdle, beforeIdle)
				}
				if afterRunning := snapshotRunningMasks(&p); afterRunning != beforeRunning {
					t.Errorf("after return running masks = %v, want %v", afterRunning, beforeRunning)
				}
			})
		})
	})

	t.Run("reaps idle workers after timeout", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			withIsolatedGlobalGroup(t, func(g Group) {
				var p WorkerPool
				defer func() {
					p.Close()
					synctest.Wait()
				}()

				warmParkedWorkers(t, g, &p, 4)

				time.Sleep(workerIdleTimeout - time.Millisecond)
				synctest.Wait()
				if beforeExpiry := idleCount(&p); beforeExpiry != 4 {
					t.Fatalf("idle before timeout = %d, want 4", beforeExpiry)
				}

				time.Sleep(2 * time.Millisecond)
				synctest.Wait()
				afterExpiry := idleCount(&p)
				afterRunning := runningCount(&p)
				p.reaper.mu.Lock()
				timer := p.reapTimer
				p.reaper.mu.Unlock()
				if afterExpiry != 0 || afterRunning != 0 || timer != nil {
					t.Fatalf("after timeout: idle=%d running=%d timer=%v, want 0, 0, and nil", afterExpiry, afterRunning, timer)
				}
			})
		})
	})

	t.Run("close is idempotent and post-close runs cleanly", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			withIsolatedGlobalGroup(t, func(g Group) {
				var p WorkerPool
				if err := p.FanOut(g, 4, func(int) {}); err != nil {
					t.Fatalf("FanOut error = %v", err)
				}
				synctest.Wait()

				p.Close()
				p.Close()
				synctest.Wait()
				if running := runningCount(&p); running != 0 {
					t.Errorf("post-close running = %d, want 0", running)
				}

				seen := make([]atomic.Int64, 6)
				if err := p.FanOut(g, 6, func(i int) { seen[i].Add(1) }); err != nil {
					t.Fatalf("post-close FanOut error = %v", err)
				}
				synctest.Wait()
				for i := range seen {
					if got := seen[i].Load(); got != 1 {
						t.Errorf("post-close index %d count = %d, want 1", i, got)
					}
				}
				if idle := idleCount(&p); idle != 0 {
					t.Errorf("post-close idle = %d, want 0", idle)
				}
				if running := runningCount(&p); running != 0 {
					t.Errorf("post-close running = %d, want 0", running)
				}
			})
		})
	})
}

// Not parallel: nil-group fallback reads the package-level instance.
func TestWorkerPool_FanOut_FallbackGroups(t *testing.T) {
	t.Run("serial limit=1 runs without pooling", func(t *testing.T) {
		var p WorkerPool
		t.Cleanup(p.Close)
		g, _ := New(t.Context())
		g.SetLimit(1)

		seen := make([]atomic.Int64, 5)
		if err := p.FanOut(g, 5, func(i int) { seen[i].Add(1) }); err != nil {
			t.Fatalf("FanOut error = %v", err)
		}
		if err := g.Wait(); err != nil {
			t.Fatalf("Wait error = %v", err)
		}
		for i := range seen {
			if got := seen[i].Load(); got != 1 {
				t.Errorf("index %d count = %d, want 1", i, got)
			}
		}
	})

	t.Run("rejected semaphore returns admission error and runs inline index", func(t *testing.T) {
		var p WorkerPool
		t.Cleanup(p.Close)
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		g, _ := New(ctx)
		g.SetLimit(0)

		var calls atomic.Int64
		var inlineCalled atomic.Bool
		err := p.FanOut(g, 5, func(i int) {
			calls.Add(1)
			if i == 4 {
				inlineCalled.Store(true)
			}
		})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("FanOut error = %v, want wrapping context.Canceled", err)
		}
		if !inlineCalled.Load() || calls.Load() != 1 {
			t.Fatalf("inlineCalled=%v calls=%d, want true and 1", inlineCalled.Load(), calls.Load())
		}
	})

	t.Run("foreign group uses Go directly", func(t *testing.T) {
		var p WorkerPool
		t.Cleanup(p.Close)
		var fg foreignGroup
		seen := make([]atomic.Int64, 4)
		if err := p.FanOut(&fg, 4, func(i int) { seen[i].Add(1) }); err != nil {
			t.Fatalf("FanOut error = %v", err)
		}
		for i := range seen {
			if got := seen[i].Load(); got != 1 {
				t.Errorf("index %d count = %d, want 1", i, got)
			}
		}
		if fg.calls.Load() == 0 {
			t.Errorf("expected foreignGroup.Go to be called at least once")
		}
	})

	t.Run("nil group falls back to global pool", func(t *testing.T) {
		withIsolatedGlobalGroup(t, func(Group) {
			var p WorkerPool
			t.Cleanup(p.Close)
			seen := make([]atomic.Int64, 5)
			if err := p.FanOut(nil, len(seen), func(i int) { seen[i].Add(1) }); err != nil {
				t.Fatalf("FanOut(nil) error = %v", err)
			}
			for i := range seen {
				if got := seen[i].Load(); got != 1 {
					t.Errorf("index %d count = %d, want 1", i, got)
				}
			}
		})
	})
}

func BenchmarkWorkerPoolFanOut(b *testing.B) {
	noop := func(int) {}
	for _, n := range []int{3, 5, 16, 50} {
		b.Run("Serial/n="+strconv.Itoa(n), func(b *testing.B) {
			var p WorkerPool
			b.Cleanup(p.Close)
			g := Get()
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				if err := p.FanOut(g, n, noop); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("Parallel/n="+strconv.Itoa(n), func(b *testing.B) {
			var p WorkerPool
			b.Cleanup(p.Close)
			g := Get()
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					if err := p.FanOut(g, n, noop); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}
