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
const (
	paddingSize         = 128
	defaultMaxBuckets   = 1000
	defaultTailSegments = 10
	defaultNumShards    = 16
	defaultInterval     = 10 * 1e6 // 10ms
)

// histogram is a thread-safe, sharded histogram that uses dynamic bucketing.
// It uses a TDigest to store samples and generates buckets at report time.
type histogram struct {
	digest TDigest

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
	NumShards int `json:"num_shards" yaml:"num_shards"`
}

// Init initializes the histogram.
func (h *histogram) Init() error {
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

// NewHistogram creates a new sharded histogram with dynamic bucketing.
func NewHistogram(opts ...HistogramOption) (Histogram, error) {
	cfg := histogramConfig{
		NumShards: defaultNumShards,
	}
	// Apply default options first
	for _, opt := range defaultHistogramOpts {
		_ = opt(&cfg)
	}
	// Apply user options
	for _, opt := range opts {
		_ = opt(&cfg)
	}

	if cfg.NumShards <= 1 {
		h := new(histogram)
		if err := h.Init(); err != nil {
			return nil, err
		}
		return h, nil
	}

	sh := &shardedHistogram{
		shards: make([]*histogram, cfg.NumShards),
	}
	for i := range sh.shards {
		h := new(histogram)
		if err := h.Init(); err != nil {
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
	// hash modulo length is always within int bounds
	return int(computeHash(val) % uint64(len(sh.shards))) //nolint:gosec
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
	// Use standard arithmetic instead of math.FMA for performance portability
	// FMA can be slow if emulated in software
	h.mean += delta / float64(h.total)
	h.m2 += delta * (val - h.mean)
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
		snap := GetHistogramSnapshot()
		snap.Min = math.Inf(1)
		snap.Max = math.Inf(-1)
		return snap
	}

	merged := sh.shards[0].Clone().(*histogram)
	for i := 1; i < len(sh.shards); i++ {
		_ = merged.mergeHistogram(sh.shards[i])
	}
	return merged.Snapshot()
}

func (h *histogram) Snapshot() *HistogramSnapshot {
	h.mu.Lock()
	snap := GetHistogramSnapshot()
	snap.Total = h.total
	snap.Mean = h.mean
	snap.M2 = h.m2
	snap.Min = h.min
	snap.Max = h.max

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

	digest := h.digest.Clone()
	// Reuse existing buffer for bounds if available
	snap.Bounds = calculateBucketBoundaries(digest, snap.Min, snap.Max, snap.Bounds[:0])

	neededCounts := len(snap.Bounds)
	if cap(snap.Counts) < neededCounts {
		snap.Counts = make([]uint64, neededCounts)
	} else {
		snap.Counts = snap.Counts[:neededCounts]
		// Zero out reused memory
		for i := range snap.Counts {
			snap.Counts[i] = 0
		}
	}

	bucketWeights := make([]float64, len(snap.Bounds))

	digest.ForEachCentroid(func(mean, weight float64) bool {
		idx, _ := slices.BinarySearch(snap.Bounds, mean)
		if idx < len(bucketWeights) {
			bucketWeights[idx] += weight
		} else {
			if len(bucketWeights) > 0 {
				bucketWeights[len(bucketWeights)-1] += weight
			}
		}
		return true
	})

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

	if snap.Total > 0 && len(snap.Counts) > 0 {
		if snap.Counts[0] == 0 {
			for k := 1; k < len(snap.Counts); k++ {
				if snap.Counts[k] > 0 {
					snap.Counts[k]--
					snap.Counts[0]++
					break
				}
			}
		}
		lastIdx := len(snap.Counts) - 1
		if snap.Counts[lastIdx] == 0 {
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

// calculateBucketBoundaries determines the optimal bucket boundaries based on the distribution.
func calculateBucketBoundaries(digest TDigest, minVal, maxVal float64, dst []float64) []float64 {
	p99 := digest.Quantile(0.99)
	if math.IsInf(maxVal, 0) {
		maxVal = 0
	}
	if math.IsInf(minVal, 0) {
		minVal = 0
	}

	// Automatic Interval Calculation
	// Determine scale from P90 (or P99/Max if P90 is 0)
	scaleVal := digest.Quantile(0.90)
	if scaleVal <= 0 {
		scaleVal = p99
	}
	if scaleVal <= 0 {
		scaleVal = maxVal
	}

	interval := float64(defaultInterval) // Default 10ms

	// Check lower quantiles to detect microsecond-level distributions
	p01 := digest.Quantile(0.01)
	p05 := digest.Quantile(0.05)

	// If a significant portion of requests are very fast (e.g., < 1ms), reduce interval.
	// 1ms = 1e6 ns.
	if (p01 > 0 && p01 < 1e6) || (p05 > 0 && p05 < 1e6) || (scaleVal > 0 && scaleVal < 10*1e6) {
		// Use a finer grain interval, e.g., 100us or 10us depending on min.
		if minVal > 0 && minVal < 100e3 {
			interval = 10e3 // 10us
		} else {
			interval = 100e3 // 100us
		}
	} else if scaleVal > 0 {
		// Target ~20 buckets for the body
		targetRes := scaleVal / 20.0
		if targetRes < 1.0 {
			targetRes = 1.0
		}
		interval = snapInterval(targetRes)
	}

	// Align start to interval
	start := math.Floor(minVal/interval) * interval
	current := start

	// Reuse dst as bounds
	bounds := dst
	limit := defaultMaxBuckets

	for current <= p99 {
		current += interval
		bounds = append(bounds, current)
		if len(bounds) >= limit {
			break
		}
	}

	ceilMax := math.Ceil(maxVal)
	if ceilMax > current {
		remainingBuckets := defaultTailSegments
		tailRange := ceilMax - current
		step := tailRange / float64(remainingBuckets)
		step = snapInterval(step)
		if step < interval {
			step = interval
		}

		safetyLimit := defaultMaxBuckets
		for i := 0; i < safetyLimit && current < ceilMax; i++ {
			current += step
			if current >= ceilMax {
				current = ceilMax
			}
			bounds = append(bounds, current)
		}
	}

	// Deduplicate bounds
	validBounds := make([]float64, 0, len(bounds))
	prev := 0.0
	for _, b := range bounds {
		if b > prev {
			validBounds = append(validBounds, b)
			prev = b
		}
	}
	return validBounds
}

// GetHistogramSnapshot returns a HistogramSnapshot from the pool.
func GetHistogramSnapshot() *HistogramSnapshot {
	val, ok := histogramSnapshotPool.Get().(*HistogramSnapshot)
	if !ok {
		return new(HistogramSnapshot)
	}
	return val
}

// PutHistogramSnapshot returns a HistogramSnapshot to the pool.
func PutHistogramSnapshot(hs *HistogramSnapshot) {
	hs.Reset()
	histogramSnapshotPool.Put(hs)
}

// nolint:gochecknoglobals
var histogramSnapshotPool = sync.Pool{
	New: func() any {
		return new(HistogramSnapshot)
	},
}

// Reset resets the snapshot to its zero state, preserving slice capacity.
func (s *HistogramSnapshot) Reset() {
	if s.Counts != nil {
		s.Counts = s.Counts[:0]
	}
	if s.Bounds != nil {
		s.Bounds = s.Bounds[:0]
	}
	s.Total = 0
	s.Sum = 0
	s.SumSq = 0
	s.M2 = 0
	s.Mean = 0
	s.StdDev = 0
	s.Min = 0
	s.Max = 0
}

func (sh *shardedHistogram) BoundsHash() uint64 {
	return 0
}

func (h *histogram) BoundsHash() uint64 {
	return 0
}

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

func (s *HistogramSnapshot) Merge(other *HistogramSnapshot) error {
	if other == nil || other.Total == 0 {
		return nil
	}

	if len(other.Counts) > 0 {
		if len(s.Counts) == 0 {
			s.Counts = slices.Clone(other.Counts)
			s.Bounds = slices.Clone(other.Bounds)
		} else {
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

func (s *HistogramSnapshot) EnforceExemplarConsistency(details *ExemplarDetails) {
	if s == nil || details == nil || len(s.Bounds) == 0 {
		return
	}

	process := func(items []*ExemplarItem) {
		minCounts := make(map[int]uint64)
		for _, item := range items {
			val := float64(item.Latency)
			idx, _ := slices.BinarySearch(s.Bounds, val)
			if idx < len(s.Counts) {
				minCounts[idx]++
			}
		}

		for idx, minCount := range minCounts {
			if s.Counts[idx] < minCount {
				diff := minCount - s.Counts[idx]
				s.Counts[idx] += diff

				maxC := uint64(0)
				maxIdx := -1
				for i, c := range s.Counts {
					if c > maxC {
						maxC = c
						maxIdx = i
					}
				}
				if maxIdx != -1 && maxIdx != idx && s.Counts[maxIdx] >= diff {
					s.Counts[maxIdx] -= diff
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

func snapInterval(val float64) float64 {
	if val <= 0 {
		return 1.0
	}
	mag := math.Pow(10, math.Floor(math.Log10(val)))
	norm := val / mag

	var snapped float64
	if norm < 1.5 {
		snapped = 1.0
	} else if norm < 3.5 {
		snapped = 2.0
	} else if norm < 7.5 {
		snapped = 5.0
	} else {
		snapped = 10.0
	}

	return snapped * mag
}
