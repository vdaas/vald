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
	"encoding/json"
	"slices"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

// centroid represents a centroid in the t-digest.
// It is unexported to encapsulate the implementation details of the TDigest.
type centroid struct {
	Mean   float64 `json:"mean"`
	Weight float64 `json:"weight"`
}

const (
	defaultBufferCapacity = 128
	// maxBufferCapacity defines a threshold to shrink buffers if they grow too large.
	// 4096 * 8 bytes (float64) = 32KB.
	maxBufferCapacity = 4096

	centroidWeightFactor = 2.0
)

// tdigest is a custom implementation of the t-digest algorithm (used as a shard).
//
// It is designed for high-performance, concurrent metric recording by using a
// mutex to protect the internal state.
//
// Invariants (must be preserved by all mutating methods):
//   - centroids is always sorted by Mean in ascending order.
//   - count is the sum of all centroid.Weight values.
//   - count > 0 if and only if len(centroids) > 0.
type tdigest struct {
	centroids                []centroid
	buffer                   []float64
	backBuffer               []float64
	scratch                  []centroid
	swap                     []centroid
	quantiles                []float64
	id                       uint64
	compression              float64
	compressionTriggerFactor float64
	count                    float64
	mu                       sync.Mutex
}

// shardedTDigest is a sharded wrapper around tdigest.
type shardedTDigest struct {
	shards    []*tdigest
	quantiles []float64
}

// tdigestConfig holds configuration for TDigest.
type tdigestConfig struct {
	Quantiles                []float64 `json:"quantiles"                  yaml:"quantiles"`
	Compression              float64   `json:"compression"                yaml:"compression"`
	CompressionTriggerFactor float64   `json:"compression_trigger_factor" yaml:"compression_trigger_factor"`
	NumShards                int       `json:"num_shards"                 yaml:"num_shards"`
}

// NewTDigest creates a new TDigest.
func NewTDigest(opts ...TDigestOption) (TDigest, error) {
	cfg := tdigestConfig{}
	// Apply options
	for _, opt := range append(defaultTDigestOpts, opts...) {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	if cfg.NumShards <= 1 {
		t := &tdigest{
			id:                       collectorIDCounter.Add(1),
			buffer:                   make([]float64, 0, defaultBufferCapacity),
			compression:              cfg.Compression,
			compressionTriggerFactor: cfg.CompressionTriggerFactor,
			quantiles:                cfg.Quantiles,
		}
		return t, nil
	}

	t := &shardedTDigest{
		shards:    make([]*tdigest, cfg.NumShards),
		quantiles: cfg.Quantiles,
	}
	for i := range t.shards {
		t.shards[i] = &tdigest{
			id:                       collectorIDCounter.Add(1),
			buffer:                   make([]float64, 0, defaultBufferCapacity),
			compression:              cfg.Compression,
			compressionTriggerFactor: cfg.CompressionTriggerFactor,
			// Shards don't need quantiles, only the wrapper or main TDigest does,
			// but we keep it for consistency if accessed directly.
			quantiles: cfg.Quantiles,
		}
	}
	return t, nil
}

// Reset resets the t-digest to its initial state.
func (t *shardedTDigest) Reset() {
	for _, shard := range t.shards {
		shard.Reset()
	}
}

// Reset resets the t-digest shard to its initial state.
func (t *tdigest) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Release memory if buffers are excessively large
	if cap(t.buffer) > maxBufferCapacity {
		t.buffer = make([]float64, 0, defaultBufferCapacity)
	} else {
		t.buffer = t.buffer[:0]
	}

	if cap(t.backBuffer) > maxBufferCapacity {
		t.backBuffer = make([]float64, 0, defaultBufferCapacity)
	} else {
		t.backBuffer = t.backBuffer[:0]
	}

	// For centroids, scratch, and swap, we use a similar heuristic
	// assuming typical usage won't require massive centroid lists if compressed.
	// But if they grew large (e.g. before compression or during heavy merge), shrink them.
	// Using maxBufferCapacity as a proxy for acceptable size.
	if cap(t.centroids) > maxBufferCapacity {
		t.centroids = make([]centroid, 0, defaultBufferCapacity)
	} else {
		t.centroids = t.centroids[:0]
	}

	if cap(t.scratch) > maxBufferCapacity {
		t.scratch = nil // Lazy re-allocation on next use
	} else {
		t.scratch = t.scratch[:0]
	}

	if cap(t.swap) > maxBufferCapacity {
		t.swap = nil // Lazy re-allocation on next use
	} else {
		t.swap = t.swap[:0]
	}

	t.count = 0
}

// Quantile returns the estimated quantile.
func (t *shardedTDigest) Quantile(q float64) float64 {
	merged := t.mergeAllShards()
	return merged.Quantile(q)
}

// CDF returns the estimated cumulative distribution function value for the given value.
func (t *shardedTDigest) CDF(value, min, max float64) float64 {
	merged := t.mergeAllShards()
	return merged.CDF(value, min, max)
}

// ForEachCentroid iterates over all centroids in the t-digest.
func (t *shardedTDigest) ForEachCentroid(f func(mean, weight float64) bool) {
	merged := t.mergeAllShards()
	merged.ForEachCentroid(f)
}

// mergeAllShards merges all shards into a single tdigest for querying.
func (t *shardedTDigest) mergeAllShards() *tdigest {
	// Create a new temporary shard
	merged := &tdigest{
		compression:              t.shards[0].compression,
		compressionTriggerFactor: t.shards[0].compressionTriggerFactor,
		quantiles:                t.quantiles,
	}

	// Lock all shards and collect centroids
	for _, shard := range t.shards {
		shard.mu.Lock()
		shard.flush() // flush buffer before reading centroids
		if len(shard.centroids) > 0 {
			merged.mergeCentroids(shard.centroids)
		}
		shard.mu.Unlock()
	}
	return merged
}

const (
	quantileScaleMax = 4.0 // so that max of q*(1-q) becomes 1 at q=0.5
)

// maxWeightForQuantile returns the maximum allowed combined weight for a
// centroid located at quantile q, given the total weight and compression
// parameter.
func (t *tdigest) maxWeightForQuantile(q, total float64) float64 {
	if total <= 0 || t.compression <= 0 {
		return total
	}
	return quantileScaleMax * total * q * (1 - q) / t.compression
}

// flush merges the buffered values into the centroids.
// It assumes the caller holds t.mu.
// CAUTION: It unlocks t.mu during sorting and relocks it before merging.
// Use flushLocked if you need to hold the lock throughout (e.g., to prevent deadlock in Merge).
func (t *tdigest) flush() {
	if len(t.buffer) == 0 {
		return
	}

	// Double buffering: swap buffer with backBuffer to sort outside the lock
	processing := t.buffer
	if cap(t.backBuffer) >= defaultBufferCapacity {
		t.buffer = t.backBuffer[:0]
	} else {
		t.buffer = make([]float64, 0, defaultBufferCapacity)
	}
	// Mark backBuffer as in-use (or empty)
	t.backBuffer = nil

	// Unlock to perform expensive sort
	t.mu.Unlock()

	// Sort the buffer
	slices.Sort(processing)

	// Re-lock to merge
	t.mu.Lock()

	// Convert buffer to centroids using scratch
	if cap(t.scratch) < len(processing) {
		t.scratch = make([]centroid, len(processing))
	}
	t.scratch = t.scratch[:len(processing)]

	for i, v := range processing {
		t.scratch[i] = centroid{Mean: v, Weight: 1}
	}

	// Recycle processing buffer to backBuffer if needed
	if cap(t.backBuffer) < cap(processing) {
		t.backBuffer = processing[:0]
	}

	t.mergeCentroids(t.scratch)
}

// flushLocked merges the buffered values into the centroids without releasing the lock.
// It assumes the caller holds t.mu.
func (t *tdigest) flushLocked() {
	if len(t.buffer) == 0 {
		return
	}

	// Sort the buffer to enable linear merge
	slices.Sort(t.buffer)

	// Convert buffer to centroids
	if cap(t.scratch) < len(t.buffer) {
		t.scratch = make([]centroid, len(t.buffer))
	}
	t.scratch = t.scratch[:len(t.buffer)]

	for i, v := range t.buffer {
		t.scratch[i] = centroid{Mean: v, Weight: 1}
	}
	t.buffer = t.buffer[:0]

	t.mergeCentroids(t.scratch)
}

// mergeCentroids merges a sorted slice of centroids into t.centroids.
// It assumes the caller holds t.mu.
func (t *tdigest) mergeCentroids(incoming []centroid) {
	if len(t.centroids) == 0 {
		needed := len(incoming)
		// Reuse t.centroids buffer if possible
		if cap(t.centroids) < needed {
			// Or reuse t.swap if it's large enough
			if cap(t.swap) >= needed {
				t.centroids, t.swap = t.swap, t.centroids
			} else {
				t.centroids = make([]centroid, needed)
			}
		}
		t.centroids = t.centroids[:needed]
		copy(t.centroids, incoming)

		for _, c := range incoming {
			t.count += c.Weight
		}
		if float64(len(t.centroids)) > t.compression*t.compressionTriggerFactor {
			t.compress()
		}
		return
	}

	// Merge two sorted centroid slices in linear time.
	// reuse t.swap as destination
	n1, n2 := len(t.centroids), len(incoming)
	needed := n1 + n2

	if cap(t.swap) < needed {
		t.swap = make([]centroid, needed)
	}
	t.swap = t.swap[:0]

	i, j := 0, 0
	for i < n1 && j < n2 {
		if t.centroids[i].Mean <= incoming[j].Mean {
			t.swap = append(t.swap, t.centroids[i])
			i++
		} else {
			t.swap = append(t.swap, incoming[j])
			j++
		}
	}
	if i < n1 {
		t.swap = append(t.swap, t.centroids[i:]...)
	}
	if j < n2 {
		t.swap = append(t.swap, incoming[j:]...)
	}

	// swap buffer and centroids
	t.centroids, t.swap = t.swap, t.centroids
	for _, c := range incoming {
		t.count += c.Weight
	}

	// Compress if the number of centroids exceeds the configured trigger.
	if float64(len(t.centroids)) > t.compression*t.compressionTriggerFactor {
		t.compress()
	}
}

// shardIndexForValue returns a shard index for the given value.
func (t *shardedTDigest) shardIndexForValue(val float64) int {
	return shardIndex(computeHash(val), len(t.shards))
}

// Add adds a value to the t-digest.
func (t *shardedTDigest) Add(value float64) {
	idx := t.shardIndexForValue(value)
	t.shards[idx].Add(value)
}

// Add adds a value to the t-digest shard.
func (t *tdigest) Add(value float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.buffer = append(t.buffer, value)
	if len(t.buffer) >= defaultBufferCapacity {
		t.flush()
	}
}

// Quantile returns the estimated quantile from a single shard (or merged digest).
func (t *tdigest) Quantile(q float64) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Flush any pending updates to ensure accuracy
	t.flushLocked()

	if t.count == 0 {
		return 0
	}

	if q <= 0 {
		return t.centroids[0].Mean
	}
	if q >= 1 {
		return t.centroids[len(t.centroids)-1].Mean
	}

	target := q * t.count
	var cumulative float64

	// Check if target is before the first centroid's center
	if len(t.centroids) == 0 {
		return 0
	}

	firstCenter := t.centroids[0].Weight / centroidWeightFactor
	if target < firstCenter {
		return t.centroids[0].Mean
	}

	cumulative = 0
	for i := 0; i < len(t.centroids); i++ {
		c := t.centroids[i]
		center := cumulative + c.Weight/centroidWeightFactor

		if target < center {
			if i == 0 {
				return c.Mean
			}
			prev := t.centroids[i-1]
			prevCenter := (cumulative - prev.Weight) + prev.Weight/centroidWeightFactor

			// Interpolate
			// fraction of the way from prevCenter to center
			frac := (target - prevCenter) / (center - prevCenter)
			return prev.Mean + frac*(c.Mean-prev.Mean)
		}
		cumulative += c.Weight
	}

	return t.centroids[len(t.centroids)-1].Mean
}

// CDF returns the estimated cumulative distribution function value for the given value.
func (t *tdigest) CDF(value, minVal, maxVal float64) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.flushLocked()

	if t.count == 0 || len(t.centroids) == 0 {
		return 0
	}

	if value < minVal {
		return 0
	}
	if value >= maxVal {
		return 1
	}

	// Single centroid case
	if len(t.centroids) == 1 {
		if maxVal > minVal {
			return (value - minVal) / (maxVal - minVal)
		}
		if value >= minVal {
			return 1
		}
		return 0
	}

	first := t.centroids[0]
	// Interpolate lower tail [minVal, first.Mean]
	if value < first.Mean {
		rangeWidth := first.Mean - minVal
		if rangeWidth <= 0 {
			return 0
		}
		fraction := (value - minVal) / rangeWidth
		// Weight in this range is half of first centroid
		weight := fraction * (first.Weight / centroidWeightFactor)
		return weight / t.count
	}

	last := t.centroids[len(t.centroids)-1]
	// Interpolate upper tail [last.Mean, maxVal]
	if value > last.Mean {
		rangeWidth := maxVal - last.Mean
		if rangeWidth <= 0 {
			return 1
		}
		fraction := (value - last.Mean) / rangeWidth
		weight := (last.Weight / centroidWeightFactor) * fraction
		// Base is Total - last.Weight/2
		cdfAtLast := t.count - (last.Weight / centroidWeightFactor)
		return (cdfAtLast + weight) / t.count
	}

	cumulative := 0.0
	for i := 1; i < len(t.centroids); i++ {
		prev := t.centroids[i-1]
		curr := t.centroids[i]

		if value < curr.Mean {
			// Interpolate between prev and curr
			prevCenter := cumulative + prev.Weight/centroidWeightFactor
			currCenter := (cumulative + prev.Weight) + curr.Weight/centroidWeightFactor

			fraction := (value - prev.Mean) / (curr.Mean - prev.Mean)
			weight := prevCenter + fraction*(currCenter-prevCenter)

			return weight / t.count
		}
		cumulative += prev.Weight
	}

	return 1
}

// ForEachCentroid iterates over all centroids in the t-digest.
func (t *tdigest) ForEachCentroid(f func(mean, weight float64) bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.flushLocked()

	for _, c := range t.centroids {
		if !f(c.Mean, c.Weight) {
			break
		}
	}
}

// Merge merges another t-digest into this one.
func (t *shardedTDigest) Merge(other TDigest) error {
	o, ok := other.(*shardedTDigest)
	if !ok {
		// Try to merge if other is *tdigest (shards=1)?
		if single, ok := other.(*tdigest); ok {
			// Merge single into all shards? Or distribute?
			// Distribute centroids.
			single.mu.Lock()
			defer single.mu.Unlock()
			single.flushLocked()

			// Batch centroids per shard to minimize lock contention
			batches := make([][]centroid, len(t.shards))
			for _, c := range single.centroids {
				idx := t.shardIndexForValue(c.Mean)
				batches[idx] = append(batches[idx], c)
			}

			// Apply batches to shards
			for idx, batch := range batches {
				if len(batch) > 0 {
					t.shards[idx].mu.Lock()
					t.shards[idx].mergeCentroids(batch)
					t.shards[idx].mu.Unlock()
				}
			}
			return nil
		}
		return errors.New("incompatible sketch type for merging")
	}
	if t == o {
		return nil
	}
	return mergeShards[TDigest](t.shards, o.shards)
}

// Merge merges another t-digest shard into this one.
func (t *tdigest) Merge(other TDigest) error {
	o, ok := other.(*tdigest)
	if !ok {
		// If other is sharded, we can merge all its shards into this one.
		if sharded, ok := other.(*shardedTDigest); ok {
			for _, shard := range sharded.shards {
				if err := t.Merge(shard); err != nil {
					return err
				}
			}
			return nil
		}
		return errors.New("incompatible sketch type for merging")
	}
	if t == o {
		return nil
	}

	// To prevent deadlocks, always lock in a consistent order.
	if t.id < o.id {
		t.mu.Lock()
		o.mu.Lock()
	} else {
		o.mu.Lock()
		t.mu.Lock()
	}
	defer t.mu.Unlock()
	defer o.mu.Unlock()

	// Flush buffers before merging
	t.flushLocked()
	o.flushLocked()

	if len(o.centroids) == 0 {
		return nil
	}

	// Use helper to merge
	t.mergeCentroids(o.centroids)
	return nil
}

// compress merges centroids to reduce their number while preserving accuracy.
func (t *tdigest) compress() {
	n := len(t.centroids)
	if n <= 1 {
		return
	}
	if t.count <= 0 {
		return
	}

	total := t.count
	out := t.centroids[:0]
	var cumulative float64
	current := t.centroids[0]

	for i := 1; i < n; i++ {
		next := t.centroids[i]
		mergedWeight := current.Weight + next.Weight
		mergedMean := (current.Mean*current.Weight + next.Mean*next.Weight) / mergedWeight
		// Quantile of the merged centroid center.
		// We use the cumulative weight to estimate the quantile `q` of the centroid.
		// The value `q` is clamped to the range [0, 1] to ensure validity.
		q := max(min((cumulative+mergedWeight/2)/total, 1.0), 0.0)
		// Maximum allowed weight for this quantile (shared with Add.tryMerge).
		k := t.maxWeightForQuantile(q, total)

		if mergedWeight <= k || len(out) == 0 {
			current.Mean = mergedMean
			current.Weight = mergedWeight
		} else {
			out = append(out, current)
			cumulative += current.Weight
			current = next
		}
	}

	t.centroids = append(out, current)
}

// Clone returns a deep copy of the t-digest.
func (t *shardedTDigest) Clone() TDigest {
	newT := &shardedTDigest{
		shards:    make([]*tdigest, len(t.shards)),
		quantiles: slices.Clone(t.quantiles),
	}
	for i, shard := range t.shards {
		// Fix forcetypeassert
		cloned := shard.Clone()
		if c, ok := cloned.(*tdigest); ok {
			newT.shards[i] = c
		}
	}
	return newT
}

// Clone returns a deep copy of the t-digest shard.
func (t *tdigest) Clone() TDigest {
	t.mu.Lock()
	defer t.mu.Unlock()

	newT := &tdigest{
		id:                       collectorIDCounter.Add(1),
		compression:              t.compression,
		compressionTriggerFactor: t.compressionTriggerFactor,
		count:                    t.count,
		buffer:                   make([]float64, len(t.buffer), cap(t.buffer)),
	}
	if len(t.centroids) > 0 {
		newT.centroids = slices.Clone(t.centroids)
	}
	if len(t.quantiles) > 0 {
		newT.quantiles = slices.Clone(t.quantiles)
	}
	if len(t.buffer) > 0 {
		copy(newT.buffer, t.buffer)
	}
	return newT
}

// Quantiles returns the configured quantiles.
func (t *shardedTDigest) Quantiles() []float64 {
	return slices.Clone(t.quantiles)
}

// Quantiles returns the configured quantiles.
func (t *tdigest) Quantiles() []float64 {
	return slices.Clone(t.quantiles)
}

// tdigestJSON is a DTO for JSON serialization of tdigest.
type tdigestJSON struct {
	Centroids                []centroid `json:"centroids"`
	Quantiles                []float64  `json:"quantiles"`
	Compression              float64    `json:"compression"`
	CompressionTriggerFactor float64    `json:"compression_trigger_factor"`
	Count                    float64    `json:"count"`
}

// MarshalJSON implements the json.Marshaler interface.
func (t *tdigest) MarshalJSON() ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Ensure everything is merged into centroids
	t.flushLocked()

	return json.Marshal(tdigestJSON{
		Centroids:                t.centroids,
		Quantiles:                t.quantiles,
		Compression:              t.compression,
		CompressionTriggerFactor: t.compressionTriggerFactor,
		Count:                    t.count,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *tdigest) UnmarshalJSON(data []byte) error {
	var v tdigestJSON
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.centroids = v.Centroids
	t.quantiles = v.Quantiles
	t.compression = v.Compression
	t.compressionTriggerFactor = v.CompressionTriggerFactor
	t.count = v.Count

	// Initialize buffers
	if t.buffer == nil {
		t.buffer = make([]float64, 0, defaultBufferCapacity)
	}
	if t.id == 0 {
		t.id = collectorIDCounter.Add(1)
	}
	return nil
}

// shardedTDigestJSON is a DTO for JSON serialization of shardedTDigest.
type shardedTDigestJSON struct {
	Shards    []*tdigest `json:"shards"`
	Quantiles []float64  `json:"quantiles"`
}

// MarshalJSON implements the json.Marshaler interface.
func (t *shardedTDigest) MarshalJSON() ([]byte, error) {
	return json.Marshal(shardedTDigestJSON{
		Shards:    t.shards,
		Quantiles: t.quantiles,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *shardedTDigest) UnmarshalJSON(data []byte) error {
	var v shardedTDigestJSON
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t.shards = v.Shards
	t.quantiles = v.Quantiles
	return nil
}
