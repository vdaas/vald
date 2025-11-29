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
	"fmt"
	"math"
	"slices"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

// paddingSize is the size of the padding used to prevent false sharing.
// It is set to 128 bytes, which covers two cache lines on most architectures (64 bytes each).
const paddingSize = 128

// bucketGrowthOptimized is the growth factor that allows for optimized log2 calculation (2.0).
const bucketGrowthOptimized = 2.0

// binarySearchThreshold is the number of buckets below which binary search is preferred.
const binarySearchThreshold = 100

// histogram is a thread-safe, sharded histogram that uses geometric bucketing.
// It is designed for high-performance, concurrent metric recording by distributing
// updates across multiple shards, reducing false-sharing and contention.
type histogram struct {
	bucketFinder func(float64) int
	bounds       []float64
	counts       []uint64
	numBuckets   int
	boundsHash   uint64
	total        uint64
	maxVal       float64
	m2           float64
	min          float64
	minVal       float64
	mean         float64
	growth       float64
	invLogGrowth float64
	max          float64
	mu           sync.Mutex
	_            [paddingSize]byte
}

// shardedHistogram is a sharded wrapper around histogram.
type shardedHistogram struct {
	shards []*histogram
}

// histogramConfig holds configuration for Histogram.
type histogramConfig struct {
	Min        float64 `json:"min"         yaml:"min"`
	Max        float64 `json:"max"         yaml:"max"`
	Growth     float64 `json:"growth"      yaml:"growth"`
	NumBuckets int     `json:"num_buckets" yaml:"num_buckets"`
	NumShards  int     `json:"num_shards"  yaml:"num_shards"`
}

// Init initializes the histogram with the provided options.
func (h *histogram) Init(opts ...HistogramOption) error {
	cfg := histogramConfig{}
	// Apply default options first
	for _, opt := range defaultHistogramOpts {
		if err := opt(&cfg); err != nil {
			return err
		}
	}
	// Apply user options
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return err
		}
	}

	h.minVal = cfg.Min
	h.maxVal = cfg.Max
	h.growth = cfg.Growth
	h.numBuckets = cfg.NumBuckets
	h.invLogGrowth = 1.0 / math.Log(h.growth)

	h.mu.Lock()
	if len(h.counts) < h.numBuckets {
		h.counts = make([]uint64, h.numBuckets)
	} else {
		for j := range h.counts {
			h.counts[j] = 0
		}
		h.counts = h.counts[:h.numBuckets]
	}
	h.total = 0
	h.mean = 0
	h.m2 = 0
	h.min = math.Inf(1)
	h.max = math.Inf(-1)
	h.mu.Unlock()

	if len(h.bounds) < h.numBuckets-1 {
		h.bounds = make([]float64, h.numBuckets-1)
	} else {
		h.bounds = h.bounds[:h.numBuckets-1]
	}

	// Build geometric bucket boundaries.
	// Geometric buckets grow exponentially based on the growth factor.
	// bounds[i] = min * growth^i
	h.bounds[0] = h.minVal
	for i := 1; i < h.numBuckets-1; i++ {
		h.bounds[i] = h.minVal * math.Pow(h.growth, float64(i))
	}

	// Precompute boundsHash for compatibility checks on merge.
	h.boundsHash = computeHash(h.bounds...)

	// Select optimal bucket finding strategy
	h.setBucketFinder()

	return nil
}

func (h *histogram) setBucketFinder() {
	if h.growth == bucketGrowthOptimized {
		h.bucketFinder = h.findBucketGrowth2
	} else if h.numBuckets < binarySearchThreshold {
		h.bucketFinder = h.findBucketBinarySearch
	} else {
		h.bucketFinder = h.findBucketLog
	}
}

// NewHistogram creates a new sharded histogram with geometric bucketing.
// It takes a variable number of HistogramOption functions to configure the histogram.
func NewHistogram(opts ...HistogramOption) (Histogram, error) {
	cfg := histogramConfig{}
	for _, opt := range append(defaultHistogramOpts, opts...) {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	if cfg.NumShards <= 1 {
		h := new(histogram)
		if err := h.Init(opts...); err != nil {
			return nil, err
		}
		return h, nil
	}

	sh := &shardedHistogram{
		shards: make([]*histogram, cfg.NumShards),
	}
	for i := range sh.shards {
		h := new(histogram)
		if err := h.Init(opts...); err != nil {
			return nil, err
		}
		sh.shards[i] = h
	}
	return sh, nil
}

// Reset resets the sharded histogram.
func (sh *shardedHistogram) Reset() {
	for _, h := range sh.shards {
		h.Reset()
	}
}

// Reset resets the histogram to its initial state.
func (h *histogram) Reset() {
	h.mu.Lock()
	h.total = 0
	h.mean = 0
	h.m2 = 0
	h.min = math.Inf(1)
	h.max = math.Inf(-1)
	for j := range h.counts {
		h.counts[j] = 0
	}
	h.mu.Unlock()
}

// shardIndexForValue selects a shard index for the given value.
func (sh *shardedHistogram) shardIndexForValue(val float64) int {
	if len(sh.shards) <= 1 {
		return 0
	}
	return int(computeHash(val) % uint64(len(sh.shards)))
}

// Record adds a value to the sharded histogram.
// It distributes values across shards to reduce lock contention.
func (sh *shardedHistogram) Record(val float64) {
	idx := sh.shardIndexForValue(val)
	// Ensure index is within bounds, although computeHash % len guarantees it (if len > 0)
	if idx >= 0 && idx < len(sh.shards) {
		sh.shards[idx].Record(val)
	}
}

// Record adds a value to the histogram. It is thread-safe.
//
// It hashes the value to select a shard, then updates the shard's
// bucket counts and summary statistics using Welford's algorithm.
func (h *histogram) Record(val float64) {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return
	}

	// Determine bucket index for this value.
	bucketIdx := h.findBucket(val)

	h.mu.Lock()
	defer h.mu.Unlock()

	// Update bucket count and total count.
	if bucketIdx >= 0 && bucketIdx < len(h.counts) {
		h.counts[bucketIdx]++
	}
	h.total++

	// Update min and max.
	if val < h.min {
		h.min = val
	}
	if val > h.max {
		h.max = val
	}

	// Update mean and m2 using Welford's online algorithm for variance.
	// This method is numerically stable and avoids catastrophic cancellation.
	//
	// Let M_n be the sum of squares of differences from the current mean:
	//   M_n = Σ_{i=1 to n} (x_i - mean_n)^2
	//
	// The recurrence relations are:
	//   mean_n = mean_{n-1} + (x_n - mean_{n-1}) / n
	//   M_n = M_{n-1} + (x_n - mean_{n-1}) * (x_n - mean_n)
	//
	// Here, `s.mean` is mean_n, `s.m2` is M_n, and `val` is x_n.
	// We use FMA (Fused Multiply-Add) for better precision.
	delta := val - h.mean
	h.mean = math.FMA(delta, 1.0/float64(h.total), h.mean)
	h.m2 = math.FMA(delta, val-h.mean, h.m2)
}

// findBucket delegates to the selected strategy.
func (h *histogram) findBucket(val float64) int {
	return h.bucketFinder(val)
}

// findBucketLog determines the correct bucket index for a given value using
// a logarithmic formula for geometric buckets.
//
// Formula: index = log_{growth}(val / min) + 1
//
//	index = ln(val/min) / ln(growth) + 1
func (h *histogram) findBucketLog(val float64) int {
	if val <= h.minVal {
		return 0
	}

	// Use precomputed inverse log growth to avoid division
	idx := int(math.Ceil(math.Log(val/h.minVal) * h.invLogGrowth))

	if idx < 0 {
		return 0
	}
	if idx >= h.numBuckets {
		return h.numBuckets - 1
	}
	return idx
}

// findBucketBinarySearch uses binary search on bounds.
func (h *histogram) findBucketBinarySearch(val float64) int {
	if val <= h.minVal {
		return 0
	}
	if len(h.bounds) > 0 && val > h.bounds[len(h.bounds)-1] {
		return h.numBuckets - 1
	}

	// slices.BinarySearch finds the smallest index i such that bounds[i] >= val
	idx, _ := slices.BinarySearch(h.bounds, val)
	// The buckets are:
	// 0: <= bounds[0] (which is minVal, covered by check above)
	// 1: (bounds[0], bounds[1]]
	// ...
	// i: (bounds[i-1], bounds[i]]
	//
	// If val is in bucket i, then bounds[i-1] < val <= bounds[i].
	// BinarySearch returns i such that bounds[i] >= val.
	// So the index returned matches the bucket index, except for the 0th bucket check.
	// But wait, bounds[0] == minVal.
	// if val <= minVal, returned 0.
	// if val > minVal:
	//   if val <= bounds[1], idx will be 1. Correct.
	//   if val > bounds[1] and val <= bounds[2], idx will be 2. Correct.
	return idx
}

// findBucketGrowth2 uses math.Log2 which is potentially faster than math.Log.
func (h *histogram) findBucketGrowth2(val float64) int {
	if val <= h.minVal {
		return 0
	}

	// index = log2(val / min) + 1
	idx := int(math.Ceil(math.Log2(val / h.minVal)))

	// Check bounds
	if idx < 0 {
		return 0
	}
	if idx >= h.numBuckets {
		return h.numBuckets - 1
	}
	return idx
}

// Clone returns a deep copy of the sharded histogram.
func (sh *shardedHistogram) Clone() Histogram {
	newSH := &shardedHistogram{
		shards: make([]*histogram, len(sh.shards)),
	}
	for i, h := range sh.shards {
		// Fix forcetypeassert
		cloned := h.Clone()
		c, ok := cloned.(*histogram)
		if !ok {
			panic(fmt.Sprintf("histogram: failed to cast cloned histogram: %T", cloned))
		}
		newSH.shards[i] = c
	}
	return newSH
}

// Clone returns a deep copy of the histogram.
func (h *histogram) Clone() Histogram {
	newH := new(histogram)
	newH.minVal = h.minVal
	newH.maxVal = h.maxVal
	newH.growth = h.growth
	newH.invLogGrowth = h.invLogGrowth
	newH.numBuckets = h.numBuckets
	newH.boundsHash = h.boundsHash
	newH.setBucketFinder()

	// Copy bounds
	// Bounds are immutable after initialization, so we can share the underlying array.
	newH.bounds = h.bounds

	h.mu.Lock()
	newH.total = h.total
	newH.mean = h.mean
	newH.m2 = h.m2
	newH.min = h.min
	newH.max = h.max

	if cap(newH.counts) < len(h.counts) {
		newH.counts = make([]uint64, len(h.counts))
	}
	newH.counts = newH.counts[:len(h.counts)]
	copy(newH.counts, h.counts)
	h.mu.Unlock()

	return newH
}

// Merge merges this histogram into the provided Histogram.
func (sh *shardedHistogram) Merge(other Histogram) error {
	if other == nil {
		return nil
	}
	if o, ok := other.(*shardedHistogram); ok {
		return mergeShards(sh.shards, o.shards)
	}
	if _, ok := other.(*histogram); ok {
		return errors.New("cannot merge single histogram into sharded histogram")
	}
	return errors.New("unknown histogram type")
}

// Merge merges this histogram into the provided Histogram.
func (h *histogram) Merge(other Histogram) error {
	if other == nil {
		return nil
	}
	if o, ok := other.(*histogram); ok {
		return h.mergeHistogram(o)
	}
	if _, ok := other.(*shardedHistogram); ok {
		return errors.New("cannot merge sharded histogram into single histogram")
	}
	return errors.New("unknown histogram type")
}

// mergeHistogram merges data from src into h.
func (h *histogram) mergeHistogram(src *histogram) error {
	if h.boundsHash != src.boundsHash {
		return errors.New("incompatible histograms: bounds checksum mismatch")
	}

	src.mu.Lock()
	srcTotal := src.total
	if srcTotal == 0 {
		src.mu.Unlock()
		return nil
	}
	srcMean := src.mean
	srcM2 := src.m2
	srcMin := src.min
	srcMax := src.max
	srcCounts := slices.Clone(src.counts)
	src.mu.Unlock()

	h.mu.Lock()
	defer h.mu.Unlock()

	for j := range h.counts {
		if j < len(srcCounts) {
			h.counts[j] += srcCounts[j]
		}
	}

	// Merge statistics using Welford's method for combined variance
	if h.total == 0 {
		h.total = srcTotal
		h.mean = srcMean
		h.m2 = srcM2
		h.min = srcMin
		h.max = srcMax
	} else {
		n1 := float64(h.total)
		n2 := float64(srcTotal)
		delta := srcMean - h.mean
		newTotal := n1 + n2

		h.mean = h.mean + delta*n2/newTotal
		h.m2 = h.m2 + srcM2 + delta*delta*n1*n2/newTotal
		h.total += srcTotal

		if srcMin < h.min {
			h.min = srcMin
		}
		if srcMax > h.max {
			h.max = srcMax
		}
	}
	return nil
}

// Snapshot returns a merged, consistent view of the histogram's data.
func (sh *shardedHistogram) Snapshot() *HistogramSnapshot {
	snap := &HistogramSnapshot{
		Min: math.Inf(1),
		Max: math.Inf(-1),
	}

	if len(sh.shards) == 0 {
		return snap
	}

	first := sh.shards[0]
	snap.Bounds = first.bounds
	snap.Counts = make([]uint64, first.numBuckets)

	var totalCount uint64
	var grandMean float64
	var grandM2 float64

	for _, h := range sh.shards {
		h.mu.Lock()
		total := h.total
		if total == 0 {
			h.mu.Unlock()
			continue
		}

		for j := range h.counts {
			snap.Counts[j] += h.counts[j]
		}

		if h.min < snap.Min {
			snap.Min = h.min
		}
		if h.max > snap.Max {
			snap.Max = h.max
		}

		// Aggregate Mean/M2 using Welford's algorithm
		if totalCount == 0 {
			grandMean = h.mean
			grandM2 = h.m2
			totalCount = total
		} else {
			n1 := float64(totalCount)
			n2 := float64(total)
			delta := h.mean - grandMean
			newTotal := n1 + n2

			grandMean = grandMean + delta*n2/newTotal
			grandM2 = grandM2 + h.m2 + delta*delta*n1*n2/newTotal
			totalCount += total
		}
		h.mu.Unlock()
	}

	snap.Total = totalCount
	snap.Mean = grandMean
	snap.M2 = grandM2
	snap.Sum = grandMean * float64(totalCount)

	if totalCount > 0 {
		snap.SumSq = grandM2 + (grandMean*grandMean)*float64(totalCount)
		snap.StdDev = math.Sqrt(grandM2 / float64(totalCount))
	}

	return snap
}

func (h *histogram) Snapshot() *HistogramSnapshot {
	h.mu.Lock()
	defer h.mu.Unlock()

	snap := &HistogramSnapshot{
		Counts: slices.Clone(h.counts),
		Bounds: h.bounds,
		Total:  h.total,
		Mean:   h.mean,
		M2:     h.m2,
		Min:    h.min,
		Max:    h.max,
	}

	snap.Sum = h.mean * float64(h.total)
	if h.total > 0 {
		snap.SumSq = h.m2 + (h.mean*h.mean)*float64(h.total)
		snap.StdDev = math.Sqrt(h.m2 / float64(h.total))
	} else {
		snap.Min = math.Inf(1)
		snap.Max = math.Inf(-1)
	}
	return snap
}

func (sh *shardedHistogram) BoundsHash() uint64 {
	if len(sh.shards) > 0 {
		return sh.shards[0].BoundsHash()
	}
	return 0
}

func (h *histogram) BoundsHash() uint64 {
	return h.boundsHash
}

// HistogramSnapshot represents a consistent point-in-time view of a Histogram.
type HistogramSnapshot struct {
	Counts []uint64  `json:"counts"`
	Bounds []float64 `json:"bounds"`
	Total  uint64    `json:"total"`
	Sum    float64   `json:"sum"`
	SumSq  float64   `json:"sum_sq"`
	M2     float64   `json:"m2"`
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
func (s *HistogramSnapshot) Merge(other *HistogramSnapshot) error {
	if other == nil || other.Total == 0 {
		return nil
	}

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

	if s.Total == 0 {
		s.Min = other.Min
		s.Max = other.Max
		s.Total = other.Total
		s.Sum = other.Sum
		s.SumSq = other.SumSq
		s.M2 = other.M2
		s.Mean = other.Mean
		s.StdDev = other.StdDev
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

		s.Mean = s.Mean + delta*n2/newTotal
		s.M2 = s.M2 + other.M2 + delta*delta*n1*n2/newTotal
		s.StdDev = math.Sqrt(s.M2 / newTotal)

		s.Sum += other.Sum
		s.SumSq += other.SumSq
		s.Total += other.Total
	}

	if len(s.Bounds) == 0 {
		s.Bounds = other.Bounds
	}
	return nil
}
