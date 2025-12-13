//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package metrics

import (
	"context"
	"time"
	"unsafe"

	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/zeebo/xxh3"
)

// sizeOfFloat64 is the size of float64 in bytes.
const sizeOfFloat64 = 8

// Recorder is the interface for recording metrics.
type Recorder interface {
	Record(ctx context.Context, key uint64, rr *RequestResult)
}

// Reporter is the interface for reporting and exporting data.
type Reporter interface {
	GlobalSnapshot() *GlobalSnapshot
	RangeScalesSnapshot() map[string]*ScaleSnapshot
	TimeScalesSnapshot() map[string]*ScaleSnapshot
}

// CounterManager is the interface for managing custom counters.
type CounterManager interface {
	CounterHandle(name string) (*CounterHandle, error)
	IncCounter(name string, val int64)
}

// Collector is the interface for the metrics collector.
type Collector interface {
	Recorder
	Reporter
	CounterManager
	MergeInto(dest Collector) error
	Clone() (Collector, error)
	merge(other *collector) error
	Reset()
}

// Histogram is the interface for a histogram.
type Histogram interface {
	Record(val float64)
	Merge(other Histogram) error
	Snapshot() *HistogramSnapshot
	BoundsHash() uint64
	Clone() Histogram
	Reset()
}

// Exemplar is the interface for an exemplar.
type Exemplar interface {
	Offer(latency time.Duration, requestID string, err error, msg string)
	Snapshot() []*ExemplarItem
	DetailedSnapshot() (*ExemplarDetails, error)
	Merge(other Exemplar) error
	Clone() Exemplar
	Reset()
}

// TDigest is the interface for approximate percentile estimators like t-digest.
type TDigest interface {
	Add(value float64)
	Quantile(q float64) float64
	CDF(value, min, max float64) float64
	ForEachCentroid(f func(mean, weight float64) bool)
	Quantiles() []float64
	Merge(other TDigest) error
	Clone() TDigest
	Reset()
}

// Slot is the interface for a time/range slot in a Scale.
type Slot interface {
	Record(rr *RequestResult, windowIdx uint64)
	Merge(other Slot) error
	Snapshot() *SlotSnapshot
	Reset()
	Clone() Slot
}

// Scale is the interface for a metrics scale (ring buffer of slots).
type Scale interface {
	Record(ctx context.Context, key uint64, rr *RequestResult)
	Merge(other Scale) error
	Snapshot() *ScaleSnapshot
	Reset()
	Clone() Scale
	Type() ScaleType
	Name() string
}

// ScaleType represents the type of scale (Range or Time).
type ScaleType uint8

const (
	// RangeScale buckets metrics by request ID.
	RangeScale ScaleType = iota
	// TimeScale buckets metrics by time.
	TimeScale
)

// computeHash computes a xxh3 checksum for a single float64 value.
func computeHash(val float64) uint64 {
	// Reinterpret the float64 as a byte slice to avoid allocation using unsafe.Slice.
	// using unsafe for zero-copy slice conversion for performance
	return xxh3.Hash(unsafe.Slice((*byte)(unsafe.Pointer(&val)), sizeOfFloat64)) //nolint:gosec
}

// shardIndex calculates the shard index for a given hash and number of shards.
func shardIndex(hash uint64, n int) int {
	if n <= 1 {
		return 0
	}
	// Optimized for power of 2
	if (n & (n - 1)) == 0 {
		// Fix G115: hash & uint64(n-1) will be in [0, n-1].
		return int(hash & uint64(n-1)) //nolint:gosec // bitwise and with positive int fits in int
	}
	// Fix G115: hash % uint64(n) will be in [0, n-1]. Since n is int and n > 0, the result fits in int.
	return int(hash % uint64(n)) //nolint:gosec // hash modulo length is always within int bounds
}

// shardedErrorCounts is a sharded map for error counts to reduce lock contention.
type shardedErrorCounts struct {
	shards []*errorCountShard
}

type errorCountShard struct {
	counts map[string]*atomic.Uint64
	mu     sync.RWMutex
}

func newShardedErrorCounts(numShards int) *shardedErrorCounts {
	if numShards <= 0 {
		numShards = 16
	}
	s := &shardedErrorCounts{
		shards: make([]*errorCountShard, numShards),
	}
	for i := range s.shards {
		s.shards[i] = &errorCountShard{
			counts: make(map[string]*atomic.Uint64),
		}
	}
	return s
}

func (s *shardedErrorCounts) shardIndex(key string) int {
	return shardIndex(xxh3.HashString(key), len(s.shards))
}

func (s *shardedErrorCounts) Add(key string, val uint64) {
	idx := s.shardIndex(key)
	shard := s.shards[idx]

	shard.mu.RLock()
	c, ok := shard.counts[key]
	shard.mu.RUnlock()

	if ok {
		c.Add(val)
		return
	}

	shard.mu.Lock()
	c, ok = shard.counts[key]
	if !ok {
		c = new(atomic.Uint64)
		shard.counts[key] = c
	}
	shard.mu.Unlock()
	c.Add(val)
}

func (s *shardedErrorCounts) Reset() {
	for _, shard := range s.shards {
		shard.mu.Lock()
		clear(shard.counts)
		shard.mu.Unlock()
	}
}

func (s *shardedErrorCounts) Snapshot() map[string]uint64 {
	res := make(map[string]uint64)
	for _, shard := range s.shards {
		shard.mu.RLock()
		for k, v := range shard.counts {
			if val := v.Load(); val > 0 {
				res[k] += val
			}
		}
		shard.mu.RUnlock()
	}
	return res
}

func (s *shardedErrorCounts) Merge(other *shardedErrorCounts) {
	if other == nil {
		return
	}
	// Iterate over other's shards and merge into this
	// Since keys might hash to different shards if shard count differs (though unlikely here if default used),
	// we should iterate conceptually.
	// Assuming shard count matches or we just use Add which handles routing.
	for _, shard := range other.shards {
		shard.mu.RLock()
		for k, v := range shard.counts {
			s.Add(k, v.Load())
		}
		shard.mu.RUnlock()
	}
}

func (s *shardedErrorCounts) Clone() *shardedErrorCounts {
	newS := newShardedErrorCounts(len(s.shards))
	for i, shard := range s.shards {
		targetShard := newS.shards[i]
		shard.mu.RLock()
		targetShard.mu.Lock() // Should be safe on new object
		for k, v := range shard.counts {
			nv := new(atomic.Uint64)
			nv.Store(v.Load())
			targetShard.counts[k] = nv
		}
		targetShard.mu.Unlock()
		shard.mu.RUnlock()
	}
	return newS
}
