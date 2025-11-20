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
	"slices"
	"sync"

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
// It uses a mutex to protect its state and implements Welford's algorithm
// for numerically stable variance calculation.
type histogramShard struct {
	mu     sync.Mutex
	counts []uint64  // number of values in each bucket
	total  uint64    // total number of values in this shard
	mean   float64   // mean of values in this shard
	m2     float64   // sum of squares of differences from the mean
	min    float64   // minimum value in this shard
	max    float64   // maximum value in this shard
	_      [128]byte // Padding to prevent false sharing
}

// Init initializes the histogram with the provided options.
func (h *histogram) Init(opts ...HistogramOption) error {
	// Apply default options first
	for _, opt := range defaultHistogramOpts {
		if err := opt(h); err != nil {
			return err
		}
	}
	// Apply user options
	for _, opt := range opts {
		if err := opt(h); err != nil {
			return err
		}
	}

	if len(h.shards) < h.numShards {
		h.shards = make([]histogramShard, h.numShards)
	} else {
		h.shards = h.shards[:h.numShards]
	}

	// Initialize shards.
	for i := range h.shards {
		h.shards[i].mu.Lock()
		if len(h.shards[i].counts) < h.numBuckets {
			h.shards[i].counts = make([]uint64, h.numBuckets)
		} else {
			// Reset existing counts
			for j := range h.shards[i].counts {
				h.shards[i].counts[j] = 0
			}
			h.shards[i].counts = h.shards[i].counts[:h.numBuckets]
		}
		h.shards[i].total = 0
		h.shards[i].mean = 0
		h.shards[i].m2 = 0
		h.shards[i].min = math.Inf(1)
		h.shards[i].max = math.Inf(-1)
		h.shards[i].mu.Unlock()
	}

	if len(h.bounds) < h.numBuckets-1 {
		h.bounds = make([]float64, h.numBuckets-1)
	} else {
		h.bounds = h.bounds[:h.numBuckets-1]
	}

	// Build geometric bucket boundaries.
	h.bounds[0] = h.min
	for i := 1; i < h.numBuckets-1; i++ {
		h.bounds[i] = h.min * math.Pow(h.growth, float64(i))
	}

	// Precompute boundsCRC32 for compatibility checks on merge.
	h.boundsCRC32 = computeBoundsCRC32(h.bounds)

	return nil
}

// NewHistogram creates a new sharded histogram with geometric bucketing.
// It takes a variable number of HistogramOption functions to configure the histogram.
func NewHistogram(opts ...HistogramOption) (Histogram, error) {
	h := histogramPool.Get()
	if err := h.Init(opts...); err != nil {
		histogramPool.Put(h)
		return nil, err
	}
	return h, nil
}

// Reset resets the histogram to its initial state, clearing all data but keeping capacity.
func (h *histogram) Reset() {
	for i := range h.shards {
		s := &h.shards[i]
		s.mu.Lock()
		s.total = 0
		s.mean = 0
		s.m2 = 0
		s.min = math.Inf(1)
		s.max = math.Inf(-1)
		for j := range s.counts {
			s.counts[j] = 0
		}
		s.mu.Unlock()
	}
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

// shardIndexForValue selects a shard index for the given value.
//
// It uses a lightweight hash derived from the float64 bits to distribute
// values across shards without allocations or external hashers.
func (h *histogram) shardIndexForValue(val float64) int {
	if h.numShards <= 1 {
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
// It hashes the value to select a shard, then updates the shard's
// bucket counts and summary statistics using Welford's algorithm.
func (h *histogram) Record(val float64) {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return
	}

	// Select shard for this value.
	shardIdx := h.shardIndexForValue(val)
	s := &h.shards[shardIdx]

	// Determine bucket index for this value.
	bucketIdx := h.findBucket(val)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Update bucket count and total count for this shard.
	s.counts[bucketIdx]++
	s.total++

	// Update min and max.
	if val < s.min {
		s.min = val
	}
	if val > s.max {
		s.max = val
	}

	// Update mean and m2 using Welford's algorithm.
	delta := val - s.mean
	s.mean = math.FMA(delta, 1.0/float64(s.total), s.mean)
	s.m2 = math.FMA(delta, val-s.mean, s.m2)
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
	idx, _ := slices.BinarySearch(h.bounds, val)
	return idx
}

// Clone returns a deep copy of the histogram.
func (h *histogram) Clone() Histogram {
	newH := histogramPool.Get()
	newH.Reset()

	newH.min = h.min
	newH.max = h.max
	newH.growth = h.growth
	newH.numBuckets = h.numBuckets
	newH.numShards = h.numShards
	newH.boundsCRC32 = h.boundsCRC32

	// Copy bounds
	if cap(newH.bounds) < len(h.bounds) {
		newH.bounds = make([]float64, len(h.bounds))
	}
	newH.bounds = newH.bounds[:len(h.bounds)]
	copy(newH.bounds, h.bounds)

	// Copy shards
	if cap(newH.shards) < len(h.shards) {
		newH.shards = make([]histogramShard, len(h.shards))
	}
	newH.shards = newH.shards[:len(h.shards)]

	for i := range h.shards {
		src := &h.shards[i]
		dst := &newH.shards[i]

		// We must lock the source shard to get a consistent snapshot
		src.mu.Lock()
		dst.total = src.total
		dst.mean = src.mean
		dst.m2 = src.m2
		dst.min = src.min
		dst.max = src.max

		if cap(dst.counts) < len(src.counts) {
			dst.counts = make([]uint64, len(src.counts))
		}
		dst.counts = dst.counts[:len(src.counts)]
		copy(dst.counts, src.counts)
		src.mu.Unlock()
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

		// Read source shard atomically
		srcShard.mu.Lock()
		srcTotal := srcShard.total
		if srcTotal == 0 {
			srcShard.mu.Unlock()
			continue
		}
		srcMean := srcShard.mean
		srcM2 := srcShard.m2
		srcMin := srcShard.min
		srcMax := srcShard.max
		srcCounts := slices.Clone(srcShard.counts)
		srcShard.mu.Unlock()

		// Merge into dest shard
		dstShard.mu.Lock()

		// Merge counts
		for j := range dstShard.counts {
			if j < len(srcCounts) {
				dstShard.counts[j] += srcCounts[j]
			}
		}

		// Merge statistics using Welford's method
		if dstShard.total == 0 {
			dstShard.total = srcTotal
			dstShard.mean = srcMean
			dstShard.m2 = srcM2
			dstShard.min = srcMin
			dstShard.max = srcMax
		} else {
			n1 := float64(dstShard.total)
			n2 := float64(srcTotal)
			delta := srcMean - dstShard.mean
			newTotal := n1 + n2

			dstShard.mean = dstShard.mean + delta*n2/newTotal
			dstShard.m2 = dstShard.m2 + srcM2 + delta*delta*n1*n2/newTotal
			dstShard.total += srcTotal

			if srcMin < dstShard.min {
				dstShard.min = srcMin
			}
			if srcMax > dstShard.max {
				dstShard.max = srcMax
			}
		}
		dstShard.mu.Unlock()
	}
	return nil
}

// Snapshot returns a merged, consistent view of the histogram's data.
//
// It iterates through all shards and aggregates their statistics into a
// single snapshot.
func (h *histogram) Snapshot() *HistogramSnapshot {
	snap := &HistogramSnapshot{
		Counts: make([]uint64, h.numBuckets),
		Bounds: h.bounds,
		Min:    math.Inf(1),
		Max:    math.Inf(-1),
	}

	// Temporary variables for aggregating Welford stats
	var totalCount uint64
	var grandMean float64
	var grandM2 float64

	for i := range h.shards {
		s := &h.shards[i]
		s.mu.Lock()
		total := s.total
		if total == 0 {
			s.mu.Unlock()
			continue
		}

		// Aggregate counts
		for j := range s.counts {
			snap.Counts[j] += s.counts[j]
		}

		// Aggregate Min/Max
		if s.min < snap.Min {
			snap.Min = s.min
		}
		if s.max > snap.Max {
			snap.Max = s.max
		}

		// Aggregate Mean/M2
		if totalCount == 0 {
			grandMean = s.mean
			grandM2 = s.m2
			totalCount = total
		} else {
			n1 := float64(totalCount)
			n2 := float64(total)
			delta := s.mean - grandMean
			newTotal := n1 + n2

			grandMean = grandMean + delta*n2/newTotal
			grandM2 = grandM2 + s.m2 + delta*delta*n1*n2/newTotal
			totalCount += total
		}
		s.mu.Unlock()
	}

	snap.Total = totalCount
	snap.Mean = grandMean
	snap.Sum = grandMean * float64(totalCount) // Back-calculate sum if needed
	snap.SumSq = grandM2 + snap.Sum*snap.Mean  // Approx back-calculate sumSq if needed?
	// Actually, Snapshot struct has Sum and SumSq.
	// But with Welford, we track Mean and M2.
	// Variance = M2 / Total (or Total-1).
	// snap.StdDev = sqrt(Variance).
	// We can set snap.SumSq if required by consumers, but Mean/StdDev are more important.
	// M2 = SumSq - Sum^2/N => SumSq = M2 + Sum^2/N = M2 + Mean^2 * N.

	if totalCount > 0 {
		snap.SumSq = grandM2 + (grandMean*grandMean)*float64(totalCount)
		snap.StdDev = math.Sqrt(grandM2 / float64(totalCount))
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
	// Note: Snapshot merge is still using Sum/SumSq because Snapshot struct uses them.
	// If we want Welford precision here, we'd need Mean/M2 in Snapshot.
	// But Snapshot has Mean/StdDev.
	// We can reconstruct M2 from StdDev: M2 = StdDev^2 * Total.
	// Let's try to be more precise if possible.

	otherM2 := other.StdDev * other.StdDev * float64(other.Total)
	sM2 := s.StdDev * s.StdDev * float64(s.Total)

	if s.Total == 0 {
		s.Min = other.Min
		s.Max = other.Max
		s.Mean = other.Mean
		s.StdDev = other.StdDev
		s.Sum = other.Sum
		s.SumSq = other.SumSq
	} else {
		if other.Min < s.Min {
			s.Min = other.Min
		}
		if other.Max > s.Max {
			s.Max = other.Max
		}

		n1 := float64(s.Total)
		n2 := float64(other.Total)
		delta := other.Mean - s.Mean
		newTotal := n1 + n2

		newMean := s.Mean + delta*n2/newTotal
		newM2 := sM2 + otherM2 + delta*delta*n1*n2/newTotal

		s.Mean = newMean
		s.StdDev = math.Sqrt(newM2 / newTotal)
		s.Sum += other.Sum
		// Reconstruct SumSq
		s.SumSq = newM2 + (newMean * newMean * newTotal)
	}
	s.Total += other.Total

	if len(s.Bounds) == 0 {
		s.Bounds = other.Bounds
	}
	return nil
}
