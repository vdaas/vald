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
	"hash/fnv"
	"math"
	"strings"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errors"
)

// Histogram is a thread-safe, sharded histogram that uses geometric bucketing.
// It is designed for high-performance, concurrent metric recording by distributing
// updates across multiple shards, reducing lock contention.
type Histogram struct {
	shards      []histogramShard // array of shards
	bounds      []float64        // bucket boundaries
	min, max    float64          // expected lower and upper bounds
	growth      float64          // geometric growth factor for bucket widths
	numBuckets  int              // number of buckets
	numShards   int              // number of shards
	boundsCRC32 uint32           // for merge validation
}

// histogramShard is a single shard of the histogram.
// It is unexported to encapsulate the internal implementation of the Histogram.
type histogramShard struct {
	counts []atomic.Uint64 // number of values in each bucket
	total  atomic.Uint64   // total number of values in this shard
	sum    atomic.Uint64   // sum of values in this shard
	sumSq  atomic.Uint64   // sum of squares of values in this shard
	min    atomic.Uint64   // minimum value in this shard
	max    atomic.Uint64   // maximum value in this shard
}

// NewHistogram creates a new sharded histogram with geometric bucketing.
// It takes a variable number of HistogramOption functions to configure the histogram.
func NewHistogram(opts ...HistogramOption) (*Histogram, error) {
	cfg := defaultHistogramConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.numBuckets < 2 {
		return nil, errors.New("numBuckets must be at least 2")
	}
	if cfg.numShards <= 0 {
		return nil, errors.New("numShards must be positive")
	}
	if cfg.growth <= 1 {
		return nil, errors.New("histogram growth must be > 1 for geometric buckets")
	}
	if cfg.min <= 0 {
		return nil, errors.New("histogram min must be > 0 for geometric buckets")
	}

	h := &Histogram{
		shards:     make([]histogramShard, cfg.numShards),
		bounds:     make([]float64, cfg.numBuckets-1),
		min:        cfg.min,
		max:        cfg.max,
		growth:     cfg.growth,
		numBuckets: cfg.numBuckets,
		numShards:  cfg.numShards,
	}

	h.bounds[0] = cfg.min
	for i := 1; i < cfg.numBuckets-1; i++ {
		h.bounds[i] = cfg.min * math.Pow(cfg.growth, float64(i))
	}

	for i := range h.shards {
		h.shards[i].counts = make([]atomic.Uint64, cfg.numBuckets)
		h.shards[i].min.Store(math.Float64bits(math.Inf(1)))
		h.shards[i].max.Store(math.Float64bits(math.Inf(-1)))
	}

	// Calculate boundsCRC32
	buf := make([]byte, 8*len(h.bounds))
	for i, b := range h.bounds {
		binary.LittleEndian.PutUint64(buf[i*8:], math.Float64bits(b))
	}
	h.boundsCRC32 = crc32.ChecksumIEEE(buf)

	return h, nil
}

// Record adds a value to the histogram. It is thread-safe.
// It hashes the value to select a shard, then atomically updates the shard's statistics.
// This approach minimizes contention and allows for high-throughput recording.
func (h *Histogram) Record(val float64) {
	hasher := fnv.New64a()
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(val))
	hasher.Write(buf)
	shardIdx := int(hasher.Sum64() % uint64(h.numShards))
	s := &h.shards[shardIdx]

	bucketIdx := h.findBucket(val)

	s.counts[bucketIdx].Add(1)
	s.total.Add(1)

	// Atomically update sum and sumSq using CAS loops
	valBits := math.Float64bits(val)
	for {
		oldSumBits := s.sum.Load()
		newSum := math.Float64frombits(oldSumBits) + val
		if s.sum.CompareAndSwap(oldSumBits, math.Float64bits(newSum)) {
			break
		}
	}
	for {
		oldSumSqBits := s.sumSq.Load()
		newSumSq := math.Float64frombits(oldSumSqBits) + val*val
		if s.sumSq.CompareAndSwap(oldSumSqBits, math.Float64bits(newSumSq)) {
			break
		}
	}

	// Atomically update min and max
	for {
		oldMinBits := s.min.Load()
		if val >= math.Float64frombits(oldMinBits) {
			break
		}
		if s.min.CompareAndSwap(oldMinBits, valBits) {
			break
		}
	}
	for {
		oldMaxBits := s.max.Load()
		if val <= math.Float64frombits(oldMaxBits) {
			break
		}
		if s.max.CompareAndSwap(oldMaxBits, valBits) {
			break
		}
	}
}

// findBucket determines the correct bucket index for a given value using binary search.
// This is efficient for a large number of buckets.
func (h *Histogram) findBucket(val float64) int {
	if val <= h.bounds[0] {
		return 0
	}
	if val > h.bounds[len(h.bounds)-1] {
		return h.numBuckets - 1
	}

	// Binary search for the bucket
	low, high := 0, len(h.bounds)-1
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

// Merge merges another histogram into this one.
// It requires that both histograms have the same bucket boundaries and shard count.
// Merging is done atomically, shard by shard, to ensure thread safety.
func (h *Histogram) Merge(other *Histogram) error {
	if h.boundsCRC32 != other.boundsCRC32 {
		return errors.New("incompatible histograms")
	}
	if len(h.shards) != len(other.shards) {
		return errors.New("incompatible histograms: shard count mismatch")
	}
	for i := range h.shards {
		s := &h.shards[i]
		o := &other.shards[i]

		// Load atomic values from the other shard.
		otherTotal := o.total.Load()
		if otherTotal == 0 {
			continue
		}
		otherSumBits := o.sum.Load()
		otherSumSqBits := o.sumSq.Load()
		otherMinBits := o.min.Load()
		otherMaxBits := o.max.Load()

		// Add total and counts atomically.
		s.total.Add(otherTotal)
		for j := range s.counts {
			s.counts[j].Add(o.counts[j].Load())
		}

		// Atomically update sum and sumSq using CAS loops.
		for {
			oldSumBits := s.sum.Load()
			newSum := math.Float64frombits(oldSumBits) + math.Float64frombits(otherSumBits)
			if s.sum.CompareAndSwap(oldSumBits, math.Float64bits(newSum)) {
				break
			}
		}
		for {
			oldSumSqBits := s.sumSq.Load()
			newSumSq := math.Float64frombits(oldSumSqBits) + math.Float64frombits(otherSumSqBits)
			if s.sumSq.CompareAndSwap(oldSumSqBits, math.Float64bits(newSumSq)) {
				break
			}
		}

		// Atomically update min and max.
		for {
			oldMinBits := s.min.Load()
			if math.Float64frombits(otherMinBits) >= math.Float64frombits(oldMinBits) {
				break
			}
			if s.min.CompareAndSwap(oldMinBits, otherMinBits) {
				break
			}
		}
		for {
			oldMaxBits := s.max.Load()
			if math.Float64frombits(otherMaxBits) <= math.Float64frombits(oldMaxBits) {
				break
			}
			if s.max.CompareAndSwap(oldMaxBits, otherMaxBits) {
				break
			}
		}
	}
	return nil
}

// Snapshot returns a merged, consistent view of the histogram's data.
// It iterates through all shards and aggregates their statistics into a single snapshot.
// This operation is read-only and does not block new writes to the histogram.
func (h *Histogram) Snapshot() *HistogramSnapshot {
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

// String implements the fmt.Stringer interface.
func (s *HistogramSnapshot) String() string {
	if s == nil {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("    Mean: %.2f, StdDev: %.2f, Min: %.2f, Max: %.2f, Total: %d\n", s.Mean, s.StdDev, s.Min, s.Max, s.Total))
	return sb.String()
}

// Merge merges another snapshot into this one.
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
