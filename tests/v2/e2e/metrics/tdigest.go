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
	"slices"
	"unsafe"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/zeebo/xxh3"
)

// centroid represents a centroid in the t-digest.
// It is unexported to encapsulate the implementation details of the TDigest.
type centroid struct {
	Mean   float64
	Weight float64
}

const defaultBufferCapacity = 128

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
	mu                       sync.Mutex
	id                       uint64
	centroids                []centroid
	buffer                   []float64
	scratch                  []centroid
	swap                     []centroid
	compression              float64
	compressionTriggerFactor float64
	count                    float64
	quantiles                []float64
}

// shardedTDigest is a sharded wrapper around tdigest.
type shardedTDigest struct {
	shards    []*tdigest
	quantiles []float64
}

// NewTDigest creates a new TDigest.
func NewTDigest(opts ...TDigestOption) (TDigest, error) {
	cfg := TDigestConfig{}
	// Apply defaults via option functions first?
	// The default options in option.go are applied in NewTDigest usually by appending.
	// We need to apply defaults to cfg.
	// But `defaultTDigestOpts` are `func(*TDigestConfig) error`.

	// Apply defaults
	for _, opt := range defaultTDigestOpts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}
	// Apply user options
	for _, opt := range opts {
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
			quantiles:                cfg.Quantiles,
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
	t.centroids = t.centroids[:0]
	t.buffer = t.buffer[:0]
	t.count = 0
	t.mu.Unlock()
}

// String implements the fmt.Stringer interface.
func (t *shardedTDigest) String() string {
	if t == nil {
		return "No data collected for percentiles.\n"
	}
	quantiles := t.quantiles
	if len(quantiles) == 0 {
		return ""
	}

	merged := t.mergeAllShards()

	var sb strings.Builder
	for _, q := range quantiles {
		fmt.Fprintf(&sb, "\tp%d:\t%.2f", uint(q*100), merged.Quantile(q))
	}
	fmt.Fprint(&sb, "\n")
	return sb.String()
}

// String implements the fmt.Stringer interface for tdigest (shard).
func (t *tdigest) String() string {
	if t == nil {
		return "No data collected for percentiles.\n"
	}
	quantiles := t.quantiles

	var sb strings.Builder
	for _, q := range quantiles {
		fmt.Fprintf(&sb, "\tp%d:\t%.2f", uint(q*100), t.Quantile(q))
	}
	fmt.Fprint(&sb, "\n")
	return sb.String()
}

// Quantile returns the estimated quantile.
func (t *shardedTDigest) Quantile(q float64) float64 {
	merged := t.mergeAllShards()
	return merged.Quantile(q)
}

// mergeAllShards merges all shards into a single tdigest for querying.
func (t *shardedTDigest) mergeAllShards() *tdigest {
	// Create a new temporary shard
	merged := &tdigest{
		compression:              t.shards[0].compression,
		compressionTriggerFactor: t.shards[0].compressionTriggerFactor,
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
func (t *tdigest) flush() {
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
		t.centroids = slices.Clone(incoming)
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
	if len(t.shards) <= 1 {
		return 0
	}
	// Use xxh3 for hashing the float64 value to ensure good distribution across shards.
	// We interpret the float64 as a byte slice without allocation using unsafe.
	// This avoids allocating a new byte slice for every Record call.
	// The sliceHeader struct is defined in histogram.go (shared in package metrics).
	//nolint:gosec
	h := xxh3.Hash(*(*[]byte)(unsafe.Pointer(&sliceHeader{
		Data: unsafe.Pointer(&val),
		Len:  8,
		Cap:  8,
	})))
	return int(h % uint64(len(t.shards)))
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
	t.flush()

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
	var sum float64

	for i, c := range t.centroids {
		nextSum := sum + c.Weight
		if target <= nextSum {
			if i == 0 {
				return c.Mean
			}
			prev := t.centroids[i-1]
			if c.Weight <= 0 {
				return c.Mean
			}
			fraction := (target - sum) / c.Weight
			if fraction < 0 {
				fraction = 0
			} else if fraction > 1 {
				fraction = 1
			}
			return prev.Mean + (c.Mean-prev.Mean)*fraction
		}
		sum = nextSum
	}
	return t.centroids[len(t.centroids)-1].Mean
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
			single.flush()
			for _, c := range single.centroids {
				// We need to Add weighted value? TDigest doesn't support adding weighted value easily via Add.
				// But we can merge centroids into specific shards based on Mean.
				idx := t.shardIndexForValue(c.Mean)
				t.shards[idx].mu.Lock()
				t.shards[idx].mergeCentroids([]centroid{c})
				t.shards[idx].mu.Unlock()
			}
			return nil
		}
		return errors.New("incompatible sketch type for merging")
	}
	if t == o {
		return nil
	}
	if len(t.shards) != len(o.shards) {
		return errors.New("incompatible shard count for merging")
	}

	// Merge shard by shard
	for i := range t.shards {
		if err := t.shards[i].Merge(o.shards[i]); err != nil {
			return err
		}
	}
	return nil
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
	t.flush()
	o.flush()

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

	t.centroids = slices.Clip(append(out, current))
}

// Clone returns a deep copy of the t-digest.
func (t *shardedTDigest) Clone() TDigest {
	newT := &shardedTDigest{
		shards:    make([]*tdigest, len(t.shards)),
		quantiles: slices.Clone(t.quantiles),
	}
	for i, shard := range t.shards {
		newT.shards[i] = shard.Clone().(*tdigest)
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
