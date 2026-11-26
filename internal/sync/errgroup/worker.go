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
	"math/bits"
	"math/rand/v2"
	"time"

	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"golang.org/x/sys/cpu"
)

const (
	// workerIdleTimeout is how long a worker may stay idle before being reaped.
	workerIdleTimeout = 30 * time.Second

	// reapInterval is half of workerIdleTimeout so a worker untouched across two
	// consecutive sweeps is retired after [reapInterval, 2*reapInterval] without
	// calling time.Now() on every parkIdle.
	reapInterval = workerIdleTimeout / 2

	poolShards        = 8
	poolShardMask     = poolShards - 1
	slotsPerShard     = 64
	directDispatchMax = 16

	batchWaveInit   = 8
	batchWaveMax    = 16
	batchWaveGrowth = 2
	idleReapEpochs  = 2
)

// directBatch coordinates small FanOut batches (n <= 16). Cache-line padding
// isolates read-only fields, atomic cursor, and WaitGroup against false sharing.
type directBatch struct {
	fn     func(int)
	workFn func()
	last   int32
	_      cpu.CacheLinePad
	next   atomic.Int32
	_      cpu.CacheLinePad
	wg     sync.WaitGroup
	_      cpu.CacheLinePad
}

// directBatchPool reuses directBatch carriers across small FanOut calls.
var directBatchPool = sync.Pool[*directBatch]{New: newDirectBatch} //nolint:gochecknoglobals // batch carriers are reused across FanOut calls

func newDirectBatch() *directBatch {
	b := &directBatch{}
	b.workFn = b.work
	return b
}

func getDirectBatch() *directBatch {
	if b := directBatchPool.Get(); b != nil {
		return b
	}
	return newDirectBatch()
}

func putDirectBatch(b *directBatch) {
	b.fn = nil
	b.last = 0
	directBatchPool.Put(b)
}

func (b *directBatch) work() {
	defer b.wg.Done()
	for {
		if b.next.Load() >= b.last {
			return
		}
		i := b.next.Add(1) - 1
		if i >= b.last {
			return
		}
		b.fn(int(i))
	}
}

// batchState coordinates one wave-dispatched FanOut batch. Full cache-line-sized
// padding isolates the read-only fields, atomic work cursor, WaitGroup, and
// adjacent heap-allocated batches onto distinct cache lines even under 16-byte
// heap alignment.
type batchState struct {
	fn        func(int)
	workFn    func()
	workErrFn func() error
	last      int
	_         cpu.CacheLinePad
	next      atomic.Int64
	_         cpu.CacheLinePad
	wg        sync.WaitGroup
	_         cpu.CacheLinePad
}

// batchStatePool reuses batch descriptors across pooled FanOut calls.
var batchStatePool = sync.Pool[*batchState]{New: newBatchState} //nolint:gochecknoglobals // batch carriers are reused across FanOut calls

func newBatchState() *batchState {
	b := &batchState{}
	b.workFn = b.work
	b.workErrFn = func() error {
		b.work()
		return nil
	}
	return b
}

func getBatchState() *batchState {
	if b := batchStatePool.Get(); b != nil {
		return b
	}
	return newBatchState()
}

func putBatchState(b *batchState) {
	b.fn = nil
	b.last = 0
	b.next.Store(0)
	batchStatePool.Put(b)
}

func (b *batchState) work() {
	defer b.wg.Done()
	for {
		// Pre-check avoids an extra LOCK XADD RFO bounce after another worker drains last.
		if b.next.Load() >= int64(b.last) {
			return
		}
		i := int(b.next.Add(1) - 1)
		if i >= b.last {
			return
		}
		b.fn(i)
	}
}

// workerMsg is the mailbox payload sent to a pooled worker.
// direct == nil && batch == nil signals the worker to exit; direct != nil runs
// a single directBatch index; batch != nil runs b.work().
type workerMsg struct {
	direct *directBatch
	batch  *batchState
	idx    int
}

// workerSlot is the fixed mailbox and idle timestamp of one pooled worker.
type workerSlot struct {
	ch       chan workerMsg
	idleTick atomic.Uint64
}

// workerShard holds one shard's atomic idle and running slot bitmasks, padded to 64 bytes.
type workerShard struct {
	idleMask    atomic.Uint64
	runningMask atomic.Uint64
	_           [48]byte
}

// WorkerPool caches up to poolShards*slotsPerShard (512) idle workers across
// lock-free atomic bitmask shards and reaps expired workers on a single
// pool-level epoch timer. reapTimer and slots precede shards for minimal
// pointerBytes (govet fieldalignment) while keeping every 64-byte workerShard
// free of pointers and mutexes; reaper, nextShard, and epoch are ordered and
// padded on separate cache lines. Group.Go and Group.TryGo never use it.
// The zero value is ready to use.
type WorkerPool struct {
	reapTimer *time.Timer
	slots     [poolShards][slotsPerShard]workerSlot
	shards    [poolShards]workerShard
	reaper    struct {
		_      cpu.CacheLinePad
		mu     sync.Mutex
		armed  atomic.Bool
		closed atomic.Bool
		_      cpu.CacheLinePad
	}
	epoch struct {
		v atomic.Uint64
		_ cpu.CacheLinePad
	}
}

func takeLowestBits(mask uint64, limit int) uint64 {
	var claimed uint64
	for mask != 0 && limit > 0 {
		lsb := mask & -mask
		claimed |= lsb
		mask &^= lsb
		limit--
	}
	return claimed
}

func (p *WorkerPool) fanOutDirect(last int, fn func(i int)) {
	b := getDirectBatch()
	b.fn = fn
	b.wg.Add(last)

	if p.reaper.closed.Load() {
		b.next.Store(0)
		b.last = int32(last) //nolint:gosec // last < directDispatchMax (16)
		for range last {
			go b.workFn()
		}
		fn(last)
		b.wg.Wait()
		putDirectBatch(b)
		return
	}

	dispatched := 0
	start := rand.Uint32() //nolint:gosec // per-M runtime rand avoids shared atomic RFO contention for shard selection
	for off := range uint32(poolShards) {
		if dispatched == last {
			break
		}
		sIdx := (start + off) & poolShardMask
		sh := &p.shards[sIdx]
		for {
			m := sh.idleMask.Load()
			if m == 0 {
				break
			}
			claimed := takeLowestBits(m, last-dispatched)
			if sh.idleMask.CompareAndSwap(m, m&^claimed) {
				for claimed != 0 {
					slotIdx := bits.TrailingZeros64(claimed)
					claimed &^= uint64(1) << slotIdx
					p.slots[sIdx][slotIdx].ch <- workerMsg{direct: b, idx: dispatched}
					dispatched++
				}
				break
			}
		}
	}

	if dispatched < last && !p.reaper.closed.Load() {
		for off := range uint32(poolShards) {
			if dispatched == last || p.reaper.closed.Load() {
				break
			}
			sIdx := (start + off) & poolShardMask
			sh := &p.shards[sIdx]
			for {
				m := sh.runningMask.Load()
				free := ^m
				if free == 0 {
					break
				}
				claimed := takeLowestBits(free, last-dispatched)
				if sh.runningMask.CompareAndSwap(m, m|claimed) {
					for claimed != 0 {
						slotIdx := bits.TrailingZeros64(claimed)
						claimed &^= uint64(1) << slotIdx
						slot := &p.slots[sIdx][slotIdx]
						if slot.ch == nil {
							slot.ch = make(chan workerMsg, 1)
						}
						go p.workerLoop(sIdx, slotIdx, workerMsg{direct: b, idx: dispatched})
						dispatched++
					}
					break
				}
			}
		}
	}

	if dispatched < last {
		b.next.Store(int32(dispatched)) //nolint:gosec // dispatched < last < directDispatchMax
		b.last = int32(last)            //nolint:gosec // last < directDispatchMax
		for range last - dispatched {
			go b.workFn()
		}
	}

	fn(last)
	b.wg.Wait()
	putDirectBatch(b)
}

// FanOut runs fn(0)..fn(n-1) concurrently when n > 1, running indices 0..n-2 on
// workers and n-1 inline on the caller, and waits for all launched workers to
// return. For the unlimited process-wide group (g == Get() with no semaphore),
// idle workers of p are reused; any other group falls back to its own Group
// submission semantics without reusing p.
func (p *WorkerPool) FanOut(g Group, n int, fn func(i int)) error {
	if n <= 0 || fn == nil {
		return nil
	}
	last := n - 1
	if last == 0 {
		fn(0)
		return nil
	}
	if g == nil {
		g = Get()
	}
	gg, ok := g.(*group)
	if !ok || gg.sem != nil || g != Get() {
		return fanOutGroup(g, last, fn)
	}

	if last < directDispatchMax {
		p.fanOutDirect(last, fn)
		return nil
	}

	b := getBatchState()
	b.fn = fn
	b.last = last
	b.next.Store(0)

	waveCap := batchWaveInit
	start := rand.Uint32() //nolint:gosec // per-M runtime rand avoids shared atomic RFO contention for shard selection
	for launched := 0; launched < last; {
		if b.next.Load() >= int64(last) {
			break
		}
		want := min(last-launched, waveCap)
		waveCap = min(waveCap*batchWaveGrowth, batchWaveMax)

		if p.reaper.closed.Load() {
			for launched < last && b.next.Load() < int64(last) {
				b.wg.Add(1)
				go b.workFn()
				launched++
			}
			break
		}

		claimedSlots := 0
		for off := range uint32(poolShards) {
			if claimedSlots == want || b.next.Load() >= int64(last) {
				break
			}
			sIdx := (start + off) & poolShardMask
			sh := &p.shards[sIdx]
			for {
				m := sh.idleMask.Load()
				if m == 0 {
					break
				}
				claimed := takeLowestBits(m, want-claimedSlots)
				if sh.idleMask.CompareAndSwap(m, m&^claimed) {
					for claimed != 0 {
						slotIdx := bits.TrailingZeros64(claimed)
						claimed &^= uint64(1) << slotIdx
						if b.next.Load() < int64(last) {
							b.wg.Add(1)
							p.slots[sIdx][slotIdx].ch <- workerMsg{batch: b}
							claimedSlots++
							launched++
						} else if !p.parkIdle(sIdx, slotIdx) {
							p.slots[sIdx][slotIdx].ch <- workerMsg{}
						}
					}
					break
				}
			}
		}
		// Rotate start shard across waves by a step coprime to poolShards.
		start += poolShards - 1

		if claimedSlots < want && b.next.Load() < int64(last) && !p.reaper.closed.Load() {
			for off := range uint32(poolShards) {
				if claimedSlots == want || b.next.Load() >= int64(last) || p.reaper.closed.Load() {
					break
				}
				sIdx := (start + off) & poolShardMask
				sh := &p.shards[sIdx]
				for {
					m := sh.runningMask.Load()
					free := ^m
					if free == 0 {
						break
					}
					claimed := takeLowestBits(free, want-claimedSlots)
					if sh.runningMask.CompareAndSwap(m, m|claimed) {
						for claimed != 0 {
							if b.next.Load() >= int64(last) {
								sh.runningMask.And(^claimed)
								break
							}
							slotIdx := bits.TrailingZeros64(claimed)
							claimed &^= uint64(1) << slotIdx
							slot := &p.slots[sIdx][slotIdx]
							if slot.ch == nil {
								slot.ch = make(chan workerMsg, 1)
							}
							b.wg.Add(1)
							go p.workerLoop(sIdx, slotIdx, workerMsg{batch: b})
							claimedSlots++
							launched++
						}
						break
					}
				}
			}
		}

		for claimedSlots < want && b.next.Load() < int64(last) {
			b.wg.Add(1)
			go b.workFn()
			claimedSlots++
			launched++
		}
	}

	fn(last)
	b.wg.Wait()
	putBatchState(b)
	return nil
}

// fanOutGroup runs indices 0..last-1 via submitGroup on a non-pooled Group and
// runs index last inline on the caller. If g rejects every worker before
// starting it, the first admission error is returned.
func fanOutGroup(g Group, last int, fn func(i int)) error {
	b := getBatchState()
	b.fn = fn
	b.last = last
	b.next.Store(0)
	var (
		admitErr error
		accepted int
	)
	for range last {
		if b.next.Load() >= int64(last) {
			break
		}
		b.wg.Add(1)
		if err := submitGroup(g, b.workErrFn); err != nil {
			b.wg.Done()
			if admitErr == nil {
				admitErr = err
			}
			continue
		}
		accepted++
	}
	fn(last)
	b.wg.Wait()
	putBatchState(b)
	if accepted == 0 {
		return admitErr
	}
	return nil
}

// Close releases all idle workers and makes subsequent tasks run on plain
// goroutines. Busy and reserved workers exit when they try to park. Close is
// idempotent.
func (p *WorkerPool) Close() {
	p.reaper.closed.Store(true)
	p.reaper.mu.Lock()
	if p.reapTimer != nil {
		p.reapTimer.Stop()
		p.reapTimer = nil
	}
	p.reaper.armed.Store(false)
	p.reaper.mu.Unlock()

	for sIdx := range uint32(poolShards) {
		m := p.shards[sIdx].idleMask.Swap(0)
		for m != 0 {
			slotIdx := bits.TrailingZeros64(m)
			m &^= uint64(1) << slotIdx
			p.slots[sIdx][slotIdx].ch <- workerMsg{}
		}
	}
}

func (p *WorkerPool) hasAnyIdle() bool {
	for i := range p.shards {
		if p.shards[i].idleMask.Load() != 0 {
			return true
		}
	}
	return false
}

func (p *WorkerPool) ensureReaper() {
	p.reaper.mu.Lock()
	if !p.reaper.closed.Load() && p.reapTimer == nil {
		p.reaper.armed.Store(true)
		p.reapTimer = time.AfterFunc(reapInterval, p.reap)
	}
	p.reaper.mu.Unlock()
}

func (p *WorkerPool) parkBitOnly(sIdx uint32, slotIdx int) bool {
	sh := &p.shards[sIdx]
	bit := uint64(1) << slotIdx
	sh.idleMask.Or(bit)
	if p.reaper.closed.Load() {
		if (sh.idleMask.And(^bit) & bit) != 0 {
			return false
		}
		// Close() already swapped out this bit and sent the exit sentinel to slot.ch;
		// returning true lets workerLoop receive and drain that sentinel cleanly.
		return true
	}
	return true
}

// parkIdle publishes the slot as idle in its shard bitmask with the current sweep
// epoch and ensures the background reaper timer is armed. Returns false if closed.
func (p *WorkerPool) parkIdle(sIdx uint32, slotIdx int) bool {
	if p.reaper.closed.Load() {
		return false
	}
	slot := &p.slots[sIdx][slotIdx]
	if ep := p.epoch.v.Load(); slot.idleTick.Load() != ep {
		slot.idleTick.Store(ep)
	}
	if !p.parkBitOnly(sIdx, slotIdx) {
		return false
	}
	if !p.reaper.armed.Load() {
		p.ensureReaper()
	}
	return true
}

// reap sweeps every non-empty shard, retiring slots whose epoch delta >= 2
// after re-verifying idleTick post-CAS, and re-arming or disarming reapTimer
// cleanly under reaper.mu.
func (p *WorkerPool) reap() {
	if p.reaper.closed.Load() {
		return
	}
	nowEpoch := p.epoch.v.Add(1)
	var survivors int32
	for sIdx := range uint32(poolShards) {
		sh := &p.shards[sIdx]
		m := sh.idleMask.Load()
		for m != 0 {
			slotIdx := bits.TrailingZeros64(m)
			bit := uint64(1) << slotIdx
			m &^= bit
			slot := &p.slots[sIdx][slotIdx]
			if nowEpoch-slot.idleTick.Load() < idleReapEpochs {
				survivors++
				continue
			}
			claimed := false
			for {
				cur := sh.idleMask.Load()
				if cur&bit == 0 {
					break
				}
				if sh.idleMask.CompareAndSwap(cur, cur&^bit) {
					claimed = true
					break
				}
			}
			if !claimed {
				continue
			}
			if nowEpoch-slot.idleTick.Load() >= idleReapEpochs {
				slot.ch <- workerMsg{}
				continue
			}
			if !p.parkBitOnly(sIdx, slotIdx) {
				slot.ch <- workerMsg{}
			} else {
				survivors++
			}
		}
	}

	p.reaper.mu.Lock()
	if p.reaper.closed.Load() {
		p.reaper.mu.Unlock()
		return
	}
	if survivors > 0 {
		if p.reapTimer != nil {
			p.reapTimer.Reset(reapInterval)
		}
		p.reaper.mu.Unlock()
		return
	}
	// Dekker handoff with parkIdle: disarm before re-checking shard masks so an
	// idle worker parked after idleMask.Or(...) either stays visible here or sees
	// armed == false and calls ensureReaper().
	p.reaper.armed.Store(false)
	p.reapTimer = nil
	if p.hasAnyIdle() {
		p.reaper.armed.Store(true)
		p.reapTimer = time.AfterFunc(reapInterval, p.reap)
	}
	p.reaper.mu.Unlock()
}

// workerLoop runs msg, then publishes itself as idle in its shard bitmask and
// waits on its private mailbox until a new batch arrives or reap/Close sends
// an exit sentinel (direct == nil && batch == nil).
func (p *WorkerPool) workerLoop(sIdx uint32, slotIdx int, msg workerMsg) {
	sh := &p.shards[sIdx]
	slot := &p.slots[sIdx][slotIdx]
	bit := uint64(1) << slotIdx
	defer sh.runningMask.And(^bit)
	for {
		if msg.direct != nil {
			b := msg.direct
			b.fn(msg.idx)
			b.wg.Done()
		} else if msg.batch != nil {
			msg.batch.work()
		} else {
			return
		}
		if !p.parkIdle(sIdx, slotIdx) {
			return
		}
		msg = <-slot.ch
	}
}
