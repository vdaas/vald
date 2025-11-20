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
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"math"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errors"
)

// histogram is a thread-safe, sharded histogram that uses geometric bucketing.
// It is designed for high-performance, concurrent metric recording by distributing
// updates across multiple shards, reducing false-sharing and contention.
type histogram struct {
	shards      []histogramShard // array of shards
	bounds      []float64        // bucket boundaries (length = numBuckets-1)
	min, max    float64          // expected lower and upper bounds (configuration hint)
	growth      float64          // geometric growth factor for bucket widths
	numBuckets  int              // number of buckets
	numShards   int              // number of shards
	boundsCRC32 uint32           // checksum of bounds for merge validation
}

// histogramShard is a single shard of the histogram.
// It is unexported to encapsulate the internal implementation of the Histogram.
type histogramShard struct {
	counts []atomic.Uint64 // number of values in each bucket
	total  atomic.Uint64   // total number of values in this shard
	sum    atomic.Uint64   // sum of values in this shard (stored as float64 bits)
	sumSq  atomic.Uint64   // sum of squares of values in this shard (stored as float64 bits)
	min    atomic.Uint64   // minimum value in this shard (stored as float64 bits)
	max    atomic.Uint64   // maximum value in this shard (stored as float64 bits)
}

// NewHistogram creates a new sharded histogram with geometric bucketing.
// It takes a variable number of HistogramOption functions to configure the histogram.
func NewHistogram(opts ...HistogramOption) (Histogram, error) {
	h := new(histogram)
	for _, opt := range append(defaultHistogramOpts, opts...) {
		err := opt(h)
		if err != nil {
			return nil, err
		}
	}
	if len(h.shards) < h.numShards {
		h.shards = make([]histogramShard, h.numShards)
	}

	// Initialize shards.
	for i := range h.shards {
		// counts length must match numBuckets (buckets = boundaries+2 edges).
		h.shards[i].counts = make([]atomic.Uint64, h.numBuckets)
		h.shards[i].min.Store(math.Float64bits(math.Inf(1)))
		h.shards[i].max.Store(math.Float64bits(math.Inf(-1)))
	}

	if len(h.bounds) < h.numBuckets-1 {
		h.bounds = make([]float64, h.numBuckets-1)
	}

	// Build geometric bucket boundaries.
	h.bounds[0] = h.min
	for i := 1; i < h.numBuckets-1; i++ {
		h.bounds[i] = h.min * math.Pow(h.growth, float64(i))
	}

	// Precompute boundsCRC32 for compatibility checks on merge.
	h.boundsCRC32 = computeBoundsCRC32(h.bounds)

	return h, nil
}

// computeBoundsCRC32 computes a CRC32 checksum for the given slice of
// float64 bounds. It is used to validate histogram compatibility on merge.
func computeBoundsCRC32(bounds []float64) uint32 {
	if len(bounds) == 0 {
		return 0
	}
	buf := make([]byte, 8*len(bounds))
	for i, b := range bounds {
		binary.LittleEndian.PutUint64(buf[i*8:], math.Float64bits(b))
	}
	return crc32.ChecksumIEEE(buf)
}

// atomicAddFloat64 adds delta to the float64 value stored in dst using a CAS loop.
//
// The value is stored as a float64 encoded in the Uint64, so this helper
// centralizes the unsafe-but-necessary conversion logic.
func atomicAddFloat64(dst *atomic.Uint64, delta float64) {
	for {
		oldBits := dst.Load()
		oldVal := math.Float64frombits(oldBits)
		newVal := oldVal + delta
		if dst.CompareAndSwap(oldBits, math.Float64bits(newVal)) {
			return
		}
	}
}

// atomicUpdateMinFloat64 updates dst with val if val is smaller than the
// current value stored in dst. The value is stored as float64 bits.
func atomicUpdateMinFloat64(dst *atomic.Uint64, val float64) {
	valBits := math.Float64bits(val)
	for {
		oldBits := dst.Load()
		oldVal := math.Float64frombits(oldBits)
		if val >= oldVal {
			return
		}
		if dst.CompareAndSwap(oldBits, valBits) {
			return
		}
	}
}

// atomicUpdateMaxFloat64 updates dst with val if val is larger than the
// current value stored in dst. The value is stored as float64 bits.
func atomicUpdateMaxFloat64(dst *atomic.Uint64, val float64) {
	valBits := math.Float64bits(val)
	for {
		oldBits := dst.Load()
		oldVal := math.Float64frombits(oldBits)
		if val <= oldVal {
			return
		}
		if dst.CompareAndSwap(oldBits, valBits) {
			return
		}
	}
}

// shardIndexForValue selects a shard index for the given value.
//
// It uses a lightweight hash derived from the float64 bits to distribute
// values across shards without allocations or external hashers.
func (h *histogram) shardIndexForValue(val float64) int {
	if h.numShards == 1 {
		return 0
	}
	bits := math.Float64bits(val)
	// Simple bit-mixing to avoid pathological patterns.
	// This is not cryptographically strong, just enough to spread values.
	x := bits ^ (bits >> 33)
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	x *= 0xc4ceb9fe1a85ec53
	x ^= x >> 33
	return int(x % uint64(h.numShards))
}

// Record adds a value to the histogram. It is thread-safe.
//
// It hashes the value to select a shard, then atomically updates the shard's
// bucket counts and summary statistics. This approach minimizes contention
// and allows for high-throughput recording on multi-core systems.
func (h *histogram) Record(val float64) {
	// Select shard for this value.
	shardIdx := h.shardIndexForValue(val)
	s := &h.shards[shardIdx]

	// Determine bucket index for this value.
	bucketIdx := h.findBucket(val)

	// Update bucket count and total count for this shard.
	s.counts[bucketIdx].Add(1)
	s.total.Add(1)

	// Atomically update sum and sum of squares.
	atomicAddFloat64(&s.sum, val)
	atomicAddFloat64(&s.sumSq, val*val)

	// Atomically update min and max.
	atomicUpdateMinFloat64(&s.min, val)
	atomicUpdateMaxFloat64(&s.max, val)
}

// findBucket determines the correct bucket index for a given value using
// binary search over the precomputed bounds. This is efficient for a large
// number of buckets.
//
// Buckets are interpreted as:
//
//	bucket 0:           (-inf, bounds[0]]
//	bucket i (1..N-2):  (bounds[i-1], bounds[i]]
//	bucket N-1:         (bounds[N-2], +inf)
func (h *histogram) findBucket(val float64) int {
	if val <= h.bounds[0] {
		return 0
	}
	lastIdx := len(h.bounds) - 1
	if val > h.bounds[lastIdx] {
		return h.numBuckets - 1
	}

	// Standard binary search to find the first index where bounds[idx] >= val.
	low, high := 0, lastIdx
	for low <= high {
		mid := (low + high) / 2
		if h.bounds[mid] < val {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return low
}

// Clone returns a deep copy of the histogram.
func (h *histogram) Clone() Histogram {
	newH := &histogram{
		min:         h.min,
		max:         h.max,
		growth:      h.growth,
		numBuckets:  h.numBuckets,
		numShards:   h.numShards,
		boundsCRC32: h.boundsCRC32,
	}

	if len(h.bounds) > 0 {
		newH.bounds = make([]float64, len(h.bounds))
		copy(newH.bounds, h.bounds)
	}

	if len(h.shards) > 0 {
		newH.shards = make([]histogramShard, len(h.shards))
		for i := range h.shards {
			src := &h.shards[i]
			dst := &newH.shards[i]

			if len(src.counts) > 0 {
				dst.counts = make([]atomic.Uint64, len(src.counts))
				for j := range src.counts {
					dst.counts[j].Store(src.counts[j].Load())
				}
			}
			dst.total.Store(src.total.Load())
			dst.sum.Store(src.sum.Load())
			dst.sumSq.Store(src.sumSq.Load())
			dst.min.Store(src.min.Load())
			dst.max.Store(src.max.Load())
		}
	}
	return newH
}

// Merge merges this histogram into the provided Histogram.
//
// Semantics:
//
//	h.Merge(other) merges the data from h into other.
//	Internally, this is implemented as other.merge(h).
func (h *histogram) Merge(other Histogram) error {
	return other.merge(h)
}

// merge merges src into the receiver histogram (dest).
//
// Precondition:
//   - dest.boundsCRC32 == src.boundsCRC32
//   - len(dest.shards) == len(src.shards)
//
// This method aggregates all shards of src into the corresponding shards
// of dest using atomic operations. It assumes that concurrent writes may
// still happen on src and dest, but the merge itself is safe.
func (dest *histogram) merge(src *histogram) error {
	if dest.boundsCRC32 != src.boundsCRC32 {
		return errors.New("incompatible histograms: bounds checksum mismatch")
	}
	if len(dest.shards) != len(src.shards) {
		return errors.New("incompatible histograms: shard count mismatch")
	}

	for i := range dest.shards {
		dstShard := &dest.shards[i]
		srcShard := &src.shards[i]

		srcTotal := srcShard.total.Load()
		if srcTotal == 0 {
			continue
		}

		// Merge total count and bucket counts.
		dstShard.total.Add(srcTotal)
		for j := range dstShard.counts {
			dstShard.counts[j].Add(srcShard.counts[j].Load())
		}

		// Merge scalar statistics (sum, sumSq, min, max).
		srcSum := math.Float64frombits(srcShard.sum.Load())
		srcSumSq := math.Float64frombits(srcShard.sumSq.Load())
		srcMin := math.Float64frombits(srcShard.min.Load())
		srcMax := math.Float64frombits(srcShard.max.Load())

		if srcTotal > 0 {
			atomicAddFloat64(&dstShard.sum, srcSum)
			atomicAddFloat64(&dstShard.sumSq, srcSumSq)
			atomicUpdateMinFloat64(&dstShard.min, srcMin)
			atomicUpdateMaxFloat64(&dstShard.max, srcMax)
		}
	}
	return nil
}

// Snapshot returns a merged, consistent view of the histogram's data.
//
// It iterates through all shards and aggregates their statistics into a
// single snapshot. This operation is read-only with respect to the histogram
// abstraction and does not block new writes, but the snapshot is necessarily
// approximate under concurrent updates.
func (h *histogram) Snapshot() *HistogramSnapshot {
	snap := &HistogramSnapshot{
		Counts: make([]uint64, h.numBuckets),
		Bounds: h.bounds,
		Min:    math.Inf(1),
		Max:    math.Inf(-1),
	}

	for i := range h.shards {
		s := &h.shards[i]
		total := s.total.Load()
		if total == 0 {
			continue
		}

		for j := range s.counts {
			snap.Counts[j] += s.counts[j].Load()
		}
		snap.Total += total
		snap.Sum += math.Float64frombits(s.sum.Load())
		snap.SumSq += math.Float64frombits(s.sumSq.Load())

		minVal := math.Float64frombits(s.min.Load())
		if minVal < snap.Min {
			snap.Min = minVal
		}
		maxVal := math.Float64frombits(s.max.Load())
		if maxVal > snap.Max {
			snap.Max = maxVal
		}
	}

	if snap.Total > 0 {
		snap.Mean = snap.Sum / float64(snap.Total)
		variance := (snap.SumSq / float64(snap.Total)) - (snap.Mean * snap.Mean)
		if variance > 0 {
			snap.StdDev = math.Sqrt(variance)
		}
	}
	return snap
}

// BoundsCRC32 returns the precomputed CRC32 checksum of the histogram bounds.
// It can be used to cheaply check compatibility before attempting a merge.
func (h *histogram) BoundsCRC32() uint32 {
	return h.boundsCRC32
}

// HistogramSnapshot represents a consistent point-in-time view of a Histogram.
type HistogramSnapshot struct {
	Counts []uint64  `json:"counts"`
	Bounds []float64 `json:"bounds"`
	Total  uint64    `json:"total"`
	Sum    float64   `json:"sum"`
	SumSq  float64   `json:"sum_sq"`
	Mean   float64   `json:"mean"`
	StdDev float64   `json:"std_dev"`
	Min    float64   `json:"min"`
	Max    float64   `json:"max"`
}

// String implements the fmt.Stringer interface for HistogramSnapshot.
func (s *HistogramSnapshot) String() string {
	if s == nil || s.Total == 0 {
		return "No data collected.\n"
	}
	return fmt.Sprintf(
		"\tMean:\t%.2f\tStdDev:\t%.2f\tMin:\t%.2f\tMax:\t%.2f\tTotal:\t%d\n",
		s.Mean,
		s.StdDev,
		s.Min,
		s.Max,
		s.Total,
	)
}

// Merge merges another snapshot into this one.
//
// The bucket structure (Counts length and Bounds) must be compatible.
// If the current snapshot is empty, it adopts the other's structure.
func (s *HistogramSnapshot) Merge(other *HistogramSnapshot) error {
	if other == nil || other.Total == 0 {
		return nil
	}

	// Initialize or validate bucket structure.
	if len(other.Counts) > 0 {
		if len(s.Counts) == 0 {
			s.Counts = make([]uint64, len(other.Counts))
		} else if len(s.Counts) != len(other.Counts) {
			return errors.New("cannot merge histograms with different bucket counts")
		}
		for i, c := range other.Counts {
			s.Counts[i] += c
		}
	}

	// Merge scalar stats.
	if s.Total == 0 {
		s.Min = other.Min
		s.Max = other.Max
	} else {
		if other.Min < s.Min {
			s.Min = other.Min
		}
		if other.Max > s.Max {
			s.Max = other.Max
		}
	}
	s.Total += other.Total
	s.Sum += other.Sum
	s.SumSq += other.SumSq

	if s.Total > 0 {
		s.Mean = s.Sum / float64(s.Total)
		variance := (s.SumSq / float64(s.Total)) - (s.Mean * s.Mean)
		if variance > 0 {
			s.StdDev = math.Sqrt(variance)
		}
	}
	if len(s.Bounds) == 0 {
		s.Bounds = other.Bounds
	}
	return nil
}
