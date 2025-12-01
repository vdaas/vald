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
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

// paddingSize is the size of the padding used to prevent false sharing.
// It is set to 128 bytes, which covers two cache lines on most architectures (64 bytes each).
const paddingSize = 128

// histogram is a thread-safe, sharded histogram that uses dynamic bucketing.
// It uses a TDigest to store samples and generates buckets at report time.
type histogram struct {
	digest         TDigest
	bucketInterval float64
	tailSegments   int
	maxBuckets     int

	total uint64
	mean  float64
	m2    float64
	min   float64
	max   float64

	mu sync.Mutex
	_  [paddingSize]byte
}

// shardedHistogram is a sharded wrapper around histogram.
type shardedHistogram struct {
	shards []*histogram
}

// histogramConfig holds configuration for Histogram.
type histogramConfig struct {
	NumShards      int           `json:"num_shards"      yaml:"num_shards"`
	BucketInterval time.Duration `json:"bucket_interval" yaml:"bucket_interval"`
	TailSegments   int           `json:"tail_segments"   yaml:"tail_segments"`
	MaxBuckets     int           `json:"max_buckets"     yaml:"max_buckets"`
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

	h.bucketInterval = float64(cfg.BucketInterval.Nanoseconds())
	h.tailSegments = cfg.TailSegments
	h.maxBuckets = cfg.MaxBuckets

	var err error
	h.digest, err = NewTDigest(defaultTDigestOpts...)
	if err != nil {
		return err
	}

	h.mu.Lock()
	h.total = 0
	h.mean = 0
	h.m2 = 0
	h.min = math.Inf(1)
	h.max = math.Inf(-1)
	h.mu.Unlock()

	return nil
}

// EnforceExemplarConsistency ensures that histogram counts are at least as high as the observed exemplars in those buckets.
func (s *HistogramSnapshot) EnforceExemplarConsistency(details *ExemplarDetails) {
	if s == nil || details == nil || len(s.Bounds) == 0 {
		return
	}

	// Helper to process exemplars
	process := func(items []*ExemplarItem) {
		// Map bucket index -> count of exemplars in that bucket
		minCounts := make(map[int]uint64)
		for _, item := range items {
			val := float64(item.Latency)
			// Find bucket
			idx, _ := slices.BinarySearch(s.Bounds, val)
			// Bucket i is (Bounds[i-1], Bounds[i]].
			// BinarySearch returns first index where Bounds[i] >= val.
			if idx < len(s.Counts) {
				minCounts[idx]++
			}
		}

		// Enforce
		for idx, minCount := range minCounts {
			if s.Counts[idx] < minCount {
				diff := minCount - s.Counts[idx]
				s.Counts[idx] += diff

				// Steal from largest bucket to preserve Total
				maxC := uint64(0)
				maxIdx := -1
				for i, c := range s.Counts {
					if c > maxC {
						maxC = c
						maxIdx = i
					}
				}
				// If we can steal from a different bucket that has enough count
				if maxIdx != -1 && maxIdx != idx && s.Counts[maxIdx] >= diff {
					s.Counts[maxIdx] -= diff
				} else {
					// If we can't steal (unlikely in large N), we just increase Total count implicitly for the histogram view.
					// Note: This creates a small divergence from s.Total, but presentation usually relies on sum(Counts) or s.Total independently.
				}
			}
		}
	}

	if details.Slowest != nil {
		process(details.Slowest)
	}
	if details.Fastest != nil {
		process(details.Fastest)
	}
}

// NewHistogram creates a new sharded histogram with dynamic bucketing.
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
	defer h.mu.Unlock()
	h.total = 0
	h.mean = 0
	h.m2 = 0
	h.min = math.Inf(1)
	h.max = math.Inf(-1)
	if h.digest != nil {
		h.digest.Reset()
	}
}

// shardIndexForValue selects a shard index for the given value.
func (sh *shardedHistogram) shardIndexForValue(val float64) int {
	if len(sh.shards) <= 1 {
		return 0
	}
	return int(computeHash(val) % uint64(len(sh.shards))) //nolint:gosec // hash modulo length is always within int bounds
}

// Record adds a value to the sharded histogram.
func (sh *shardedHistogram) Record(val float64) {
	idx := sh.shardIndexForValue(val)
	if idx >= 0 && idx < len(sh.shards) {
		sh.shards[idx].Record(val)
	}
}

// Record adds a value to the histogram.
func (h *histogram) Record(val float64) {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return
	}

	h.digest.Add(val)

	h.mu.Lock()
	defer h.mu.Unlock()

	h.total++
	if val < h.min {
		h.min = val
	}
	if val > h.max {
		h.max = val
	}

	// Update mean and m2 using Welford's algorithm
	delta := val - h.mean
	h.mean = math.FMA(delta, 1.0/float64(h.total), h.mean)
	h.m2 = math.FMA(delta, val-h.mean, h.m2)
}

// Clone returns a deep copy of the sharded histogram.
func (sh *shardedHistogram) Clone() Histogram {
	newSH := &shardedHistogram{
		shards: make([]*histogram, len(sh.shards)),
	}
	for i, h := range sh.shards {
		cloned := h.Clone()
		c, ok := cloned.(*histogram)
		if ok {
			newSH.shards[i] = c
		}
	}
	return newSH
}

// Clone returns a deep copy of the histogram.
func (h *histogram) Clone() Histogram {
	newH := new(histogram)
	newH.bucketInterval = h.bucketInterval
	newH.tailSegments = h.tailSegments
	newH.maxBuckets = h.maxBuckets
	if h.digest != nil {
		newH.digest = h.digest.Clone()
	}

	h.mu.Lock()
	newH.total = h.total
	newH.mean = h.mean
	newH.m2 = h.m2
	newH.min = h.min
	newH.max = h.max
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
	// Check configuration compatibility?
	// For now, assume consistent config from same collector factory.

	if err := h.digest.Merge(src.digest); err != nil {
		return err
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
	src.mu.Unlock()

	h.mu.Lock()
	defer h.mu.Unlock()

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
	if len(sh.shards) == 0 {
		return &HistogramSnapshot{
			Min: math.Inf(1),
			Max: math.Inf(-1),
		}
	}

	// Merge all shards into one histogram to generate buckets correctly.
	merged := sh.shards[0].Clone().(*histogram)
	for i := 1; i < len(sh.shards); i++ {
		// Ignore error as shards are compatible
		_ = merged.mergeHistogram(sh.shards[i])
	}
	return merged.Snapshot()
}

func (h *histogram) Snapshot() *HistogramSnapshot {
	h.mu.Lock()
	snap := &HistogramSnapshot{
		Total: h.total,
		Mean:  h.mean,
		M2:    h.m2,
		Min:   h.min,
		Max:   h.max,
	}
	snap.Sum = h.mean * float64(h.total)
	if h.total > 0 {
		snap.SumSq = h.m2 + (h.mean*h.mean)*float64(h.total)
		snap.StdDev = math.Sqrt(h.m2 / float64(h.total))
	} else {
		snap.Min = math.Inf(1)
		snap.Max = math.Inf(-1)
	}
	h.mu.Unlock()

	if snap.Total == 0 {
		return snap
	}

	// Dynamic Bucketing Strategy
	// Clone digest to use for analysis (thread-safe after clone)
	// Actually Clone takes lock, so we should do it outside lock above?
	// h.digest is safe to Clone if h.digest is thread safe. TDigest impl is thread safe.
	// However, merging might be happening. TDigest locking handles it.
	digest := h.digest.Clone()

	// Step A: Analyze Distribution
	p99 := digest.Quantile(0.99)
	maxVal := snap.Max
	if math.IsInf(maxVal, 0) {
		maxVal = 0 // Should not happen if Total > 0
	}

	// Step B: Phase 1 - Main Body (Min to P99)
	var bounds []float64
	minVal := snap.Min
	if math.IsInf(minVal, 0) {
		minVal = 0
	}

	// Use h.bucketInterval (float64 ns)
	interval := h.bucketInterval
	if interval <= 0 {
		interval = 10 * 1e6 // 10ms fallback
	}

	// Adaptive Interval: If P99 is small compared to configured interval, reduce interval
	// to ensure the main body of the distribution is visualized with sufficient resolution (e.g. ~20 buckets).
	if p99 > 0 {
		targetRes := p99 / 20.0
		if targetRes < interval {
			if targetRes < 1.0 {
				targetRes = 1.0 // Minimum 1ns
			}
			interval = targetRes
		}
	}

	// Align start to interval to facilitate merging of histograms from different sources
	start := math.Floor(minVal/interval) * interval
	current := start

	// Safety check loop count
	limit := h.maxBuckets
	if limit <= 0 {
		limit = 1000
	}

	// If P99 is very large compared to interval, increase interval
	// Estimate buckets: P99 / interval.
	estimated := p99 / interval
	if estimated > float64(limit) {
		interval = p99 / float64(limit)
		// Round interval up to nicer number? keeping it simple for now.
	}

	for current <= p99 {
		current += interval
		bounds = append(bounds, current)
		if len(bounds) >= limit {
			break
		}
	}

	// Step C: Phase 2 - Long Tail (P99 to Max)
	// Ensure we cover up to Ceil(Max)
	ceilMax := math.Ceil(maxVal)
	if ceilMax > current {
		remainingBuckets := h.tailSegments
		if remainingBuckets <= 0 {
			remainingBuckets = 10
		}

		// If we already hit limit, we might only add one bucket to Max?
		// Or if we adjusted interval, we should be fine.
		// If we hit limit in Phase 1, current is at limit * interval (~P99).
		// We should add at least one bucket to Max if Max > current.

		tailRange := ceilMax - current
		step := tailRange / float64(remainingBuckets)

		for i := 0; i < remainingBuckets; i++ {
			current += step
			// Ensure we don't exceed Max due to float precision, or ensure last is Max.
			if i == remainingBuckets-1 {
				current = ceilMax
			}
			bounds = append(bounds, current)
		}
	}

	// Deduplicate bounds if any (e.g. if P99 > Max due to estimation error, or step is 0)
	// Also ensure monotonic increasing
	// Sanitize bounds
	validBounds := make([]float64, 0, len(bounds))
	prev := 0.0 // Start from 0
	for _, b := range bounds {
		if b > prev {
			validBounds = append(validBounds, b)
			prev = b
		}
	}
	snap.Bounds = validBounds

	// Calculate counts by iterating centroids (point masses)
	snap.Counts = make([]uint64, len(validBounds))
	bucketWeights := make([]float64, len(validBounds))

	digest.ForEachCentroid(func(mean, weight float64) bool {
		idx, _ := slices.BinarySearch(validBounds, mean)
		if idx < len(bucketWeights) {
			bucketWeights[idx] += weight
		} else {
			// If mean > last bound, put in last bucket (or should we expand?)
			// Given bounds cover Max (ceilMax), this should only happen for floating point epsilon issues
			if len(bucketWeights) > 0 {
				bucketWeights[len(bucketWeights)-1] += weight
			}
		}
		return true
	})

	// Convert weights to counts with cumulative rounding to preserve Total
	currentSum := 0.0
	prevIntSum := uint64(0)
	for i, w := range bucketWeights {
		currentSum += w
		targetIntSum := uint64(math.Round(currentSum))
		if targetIntSum >= prevIntSum {
			snap.Counts[i] = targetIntSum - prevIntSum
		} else {
			snap.Counts[i] = 0
		}
		prevIntSum = targetIntSum
	}

	// Ensure Min and Max are visible (count >= 1) even if probability mass is low.
	// We steal count from neighbors to preserve Total.
	if snap.Total > 0 && len(snap.Counts) > 0 {
		// Min is always in the first bucket (by alignment logic)
		if snap.Counts[0] == 0 {
			// Find nearest non-zero bucket to steal from
			for k := 1; k < len(snap.Counts); k++ {
				if snap.Counts[k] > 0 {
					snap.Counts[k]--
					snap.Counts[0]++
					break
				}
			}
		}
		// Max is always in the last bucket (by ceilMax logic)
		lastIdx := len(snap.Counts) - 1
		if snap.Counts[lastIdx] == 0 {
			// Find nearest non-zero bucket to steal from (searching backwards)
			for k := lastIdx - 1; k >= 0; k-- {
				if snap.Counts[k] > 0 {
					snap.Counts[k]--
					snap.Counts[lastIdx]++
					break
				}
			}
		}
	}

	return snap
}

func (sh *shardedHistogram) BoundsHash() uint64 {
	// Not applicable for dynamic histogram, return 0 or consistent value
	return 0
}

func (h *histogram) BoundsHash() uint64 {
	return 0
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

	// For dynamic histograms, bounds might differ.
	// Simple merging of counts is not possible if bounds mismatch.
	// However, MergeSnapshots is used for aggregating snapshots from different collectors.
	// If collectors are configured identically and see similar data, bounds might be close but likely not identical.
	// Dynamic histogram snapshot merging requires re-binning or just erroring out?
	// If we use TDigest in GlobalSnapshot, we can regenerate buckets from merged TDigest?
	// But GlobalSnapshot stores *HistogramSnapshot.
	// If HistogramSnapshot stores Counts/Bounds, we can't merge them if bounds differ.
	//
	// In the original code, bounds were static, so merging was easy.
	// With dynamic bounds, merging snapshots is hard.
	// However, GlobalSnapshot also has `LatPercentiles` (TDigest).
	// Ideally, we should merge TDigests and then generate HistogramSnapshot from merged TDigest.
	//
	// `metrics.go` `MergeSnapshots` merges histograms.
	// If we cannot merge histograms, the `Latencies` field in merged snapshot will be invalid.
	//
	// For now, I will implement a basic merge that only works if bounds match, else errors or clears buckets.
	// Or, if this is for display only, maybe we don't need to merge HistogramSnapshots?
	// But `MergeSnapshots` is used.
	//
	// Given the scope, I will implement strict check.

	if len(other.Counts) > 0 {
		if len(s.Counts) == 0 {
			s.Counts = slices.Clone(other.Counts)
			s.Bounds = slices.Clone(other.Bounds)
		} else {
			// Re-binning Merge: Combine bounds and distribute counts
			// 1. Identify union of bounds
			boundMap := make(map[float64]struct{}, len(s.Bounds)+len(other.Bounds))
			for _, b := range s.Bounds {
				boundMap[b] = struct{}{}
			}
			for _, b := range other.Bounds {
				boundMap[b] = struct{}{}
			}
			newBounds := make([]float64, 0, len(boundMap))
			for b := range boundMap {
				newBounds = append(newBounds, b)
			}
			slices.Sort(newBounds)

			// 2. Distribute counts into new buckets
			newCounts := make([]float64, len(newBounds))

			distribute := func(srcBounds []float64, srcCounts []uint64) {
				for i, count := range srcCounts {
					if count == 0 {
						continue
					}
					low := 0.0
					if i > 0 {
						low = srcBounds[i-1]
					}
					high := srcBounds[i]
					width := high - low
					if width <= 0 {
						width = 1.0
					}

					for j, nb := range newBounds {
						nLow := 0.0
						if j > 0 {
							nLow = newBounds[j-1]
						}
						nHigh := nb

						overlapStart := max(low, nLow)
						overlapEnd := min(high, nHigh)

						if overlapEnd > overlapStart {
							fraction := (overlapEnd - overlapStart) / width
							newCounts[j] += float64(count) * fraction
						}
					}
				}
			}

			distribute(s.Bounds, s.Counts)
			distribute(other.Bounds, other.Counts)

			s.Bounds = newBounds
			s.Counts = make([]uint64, len(newCounts))
			// Convert float counts to integer counts using cumulative rounding to preserve total
			currentSum := 0.0
			prevIntSum := uint64(0)
			for i, v := range newCounts {
				currentSum += v
				targetIntSum := uint64(math.Round(currentSum))
				if targetIntSum >= prevIntSum {
					s.Counts[i] = targetIntSum - prevIntSum
				} else {
					s.Counts[i] = 0
				}
				prevIntSum = targetIntSum
			}
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
	return nil
}
