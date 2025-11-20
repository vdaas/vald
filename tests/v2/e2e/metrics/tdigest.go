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
	"cmp"
	"fmt"
	"math"
	"slices"
	"unsafe"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

// centroid represents a centroid in the t-digest.
// It is unexported to encapsulate the implementation details of the TDigest.
type centroid struct {
	Mean   float64
	Weight float64
}

// tdigest is a custom implementation of the t-digest algorithm.
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
	centroids                []centroid
	compression              float64
	compressionTriggerFactor float64
	count                    float64
	quantiles                []float64
}

// NewTDigest creates a new TDigest.
func NewTDigest(opts ...TDigestOption) (TDigest, error) {
	t := new(tdigest)
	for _, opt := range append(defaultTDigestOpts, opts...) {
		err := opt(t)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

// String implements the fmt.Stringer interface.
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

// buildPrefix builds prefix sums of centroid weights.
//
// prefix[i] = sum of centroids[0:i].Weight.
// The returned slice has length len(centroids)+1, and prefix[0] is always 0.
func buildPrefix(centroids []centroid) []float64 {
	prefix := make([]float64, len(centroids)+1)
	for i, c := range centroids {
		prefix[i+1] = prefix[i] + c.Weight
	}
	return prefix
}

const (
	quantileScaleMax = 4.0 // so that max of q*(1-q) becomes 1 at q=0.5
)

// maxWeightForQuantile returns the maximum allowed combined weight for a
// centroid located at quantile q, given the total weight and compression
// parameter.
//
// This encapsulates the core t-digest formula:
//
//	k = 4 * total * q * (1-q) / compression
//
// The caller is responsible for clamping q into [0,1] if needed.
func (t *tdigest) maxWeightForQuantile(q, total float64) float64 {
	if total <= 0 || t.compression <= 0 {
		// Degenerate case: no meaningful constraint, treat as "no limit".
		return total
	}
	// Scale factor quantileScaleMax normalizes q*(1-q) to [0,1].
	return quantileScaleMax * total * q * (1 - q) / t.compression
}

// tryMerge attempts to merge the given value into the centroid at index idx.
//
// Arguments:
//   - value:  the new sample value to be merged.
//   - idx:    index of the candidate centroid in t.centroids.
//   - prefix: prefix sums of centroid weights where prefix[i] is the sum of
//     t.centroids[0:i].Weight at the time this method is called.
//   - total:  total weight (i.e. t.count) at the time this method is called.
//
// This method assumes the caller holds t.mu and that prefix/total are
// consistent with the current centroids layout.
func (t *tdigest) tryMerge(value float64, idx int, prefix []float64, total float64) bool {
	if idx < 0 || idx >= len(t.centroids) {
		return false
	}

	c := &t.centroids[idx]

	// Quantile of this centroid's center.
	q := (c.Weight/2 + prefix[idx]) / total

	// At extreme quantiles, allow merging more aggressively to avoid
	// excessive number of centroids. This also gracefully handles q
	// slightly outside [0,1] due to numerical error.
	if q <= 0 || q >= 1 {
		c.Mean = (c.Mean*c.Weight + value) / (c.Weight + 1)
		c.Weight++
		t.count++
		return true
	}

	// Maximum allowed weight for this centroid based on compression.
	k := t.maxWeightForQuantile(q, total)
	if c.Weight+1 <= k {
		// Merge new value into this centroid.
		c.Mean = (c.Mean*c.Weight + value) / (c.Weight + 1)
		c.Weight++
		t.count++
		return true
	}
	return false
}

// Add adds a value to the t-digest.
//
// This implementation keeps the centroids slice always sorted by Mean and
// avoids re-sorting on every insertion. It uses slices.BinarySearchFunc to
// locate the insertion/merge position and slices.Insert to insert a new
// centroid when necessary.
func (t *tdigest) Add(value float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Fast path for the first sample.
	if t.count == 0 {
		t.count = 1
		cd := centroid{
			Mean:   value,
			Weight: 1,
		}
		if t.centroids == nil {
			t.centroids = []centroid{cd}
		} else {
			t.centroids = append(t.centroids, cd)
		}
		return
	}

	n := len(t.centroids)

	// Binary search the position where value should be inserted
	// to keep centroids sorted by Mean.
	idx, _ := slices.BinarySearchFunc(t.centroids, value,
		func(c centroid, v float64) int {
			return cmp.Compare(c.Mean, v)
		},
	)

	// Choose the nearest centroid (by Mean) among the left neighbor and
	// the element at the insertion index, then try merging into it.
	candidate := -1
	leftIdx := idx - 1
	rightIdx := idx

	if leftIdx >= 0 && rightIdx < n {
		if math.Abs(value-t.centroids[leftIdx].Mean) <= math.Abs(t.centroids[rightIdx].Mean-value) {
			candidate = leftIdx
		} else {
			candidate = rightIdx
		}
	} else if leftIdx >= 0 {
		candidate = leftIdx
	} else if rightIdx < n {
		candidate = rightIdx
	}

	if candidate >= 0 && t.tryMerge(value, candidate,
		// Precompute prefix sums of weights:
		// prefix[i] = sum of centroids[0:i].Weight.
		buildPrefix(t.centroids), t.count) {
		// Successfully merged into an existing centroid.
		return
	}

	// Could not merge into a nearby centroid: create a new one and insert
	// at the sorted position using slices.Insert.
	t.centroids = slices.Insert(t.centroids, idx, centroid{
		Mean:   value,
		Weight: 1,
	})
	t.count++

	// Compress if the number of centroids exceeds the configured trigger.
	if float64(len(t.centroids)) > t.compression*t.compressionTriggerFactor {
		t.compress()
	}
}

// Quantile returns the estimated quantile.
//
// It performs a single pass over the centroids (O(#centroids)), which is
// efficient because the centroids slice is bounded by the compression factor.
func (t *tdigest) Quantile(q float64) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

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
				// The target is in the first centroid.
				return c.Mean
			}
			// Linear interpolation between the previous and current centroid.
			prev := t.centroids[i-1]
			if c.Weight <= 0 {
				// Degenerate case; just return current mean.
				return c.Mean
			}
			return prev.Mean + (c.Mean-prev.Mean)*max(min((target-sum)/c.Weight, 0), 1)
		}
		sum = nextSum
	}

	// Fallback: return the maximum.
	return t.centroids[len(t.centroids)-1].Mean
}

// Merge merges another t-digest into this one.
//
// This implementation assumes that both digests maintain their centroids
// sorted by Mean. It merges the two sorted slices in linear time and then
// optionally triggers compression based on the configured threshold.
func (t *tdigest) Merge(other TDigest) error {
	o, ok := other.(*tdigest)
	if !ok {
		return errors.New("incompatible sketch type for merging")
	}
	if t == o {
		// Merging the same instance is a no-op.
		return nil
	}

	// To prevent deadlocks, always lock in a consistent order.
	if uintptr(unsafe.Pointer(t)) < uintptr(unsafe.Pointer(o)) {
		t.mu.Lock()
		o.mu.Lock()
	} else {
		o.mu.Lock()
		t.mu.Lock()
	}
	defer t.mu.Unlock()
	defer o.mu.Unlock()

	if len(o.centroids) == 0 {
		return nil
	}

	// Fast path: if this digest is empty, clone the other one.
	if len(t.centroids) == 0 {
		t.centroids = slices.Clone(o.centroids)
		t.count = o.count
		if float64(len(t.centroids)) > t.compression*t.compressionTriggerFactor {
			t.compress()
		}
		return nil
	}

	// Merge two sorted centroid slices in linear time.
	i, j := 0, 0
	n1, n2 := len(t.centroids), len(o.centroids)
	merged := make([]centroid, 0, n1+n2)

	for i < n1 && j < n2 {
		if t.centroids[i].Mean <= o.centroids[j].Mean {
			merged = append(merged, t.centroids[i])
			i++
		} else {
			merged = append(merged, o.centroids[j])
			j++
		}
	}
	if i < n1 {
		merged = append(merged, t.centroids[i:]...)
	}
	if j < n2 {
		merged = append(merged, o.centroids[j:]...)
	}

	t.centroids = merged
	t.count += o.count

	// Compress if the number of centroids exceeds the configured trigger.
	if float64(len(t.centroids)) > t.compression*t.compressionTriggerFactor {
		t.compress()
	}
	return nil
}

// compress merges centroids to reduce their number while preserving the
// t-digest shape.
//
// This implementation performs a single linear pass over the sorted centroids
// and greedily merges adjacent centroids as long as the merged centroid does
// not exceed the quantile-based weight limit. This reduces the complexity
// from O(n^2) to O(n) per compression.
//
// It also reuses the underlying slice of t.centroids to minimize allocations
// and applies slices.Clip at the end to trim any excess capacity.
func (t *tdigest) compress() {
	n := len(t.centroids)
	if n <= 1 {
		return
	}
	if t.count <= 0 {
		return
	}

	total := t.count

	// Reuse the existing slice backing array to avoid extra allocations.
	out := t.centroids[:0]

	// cumulative is the sum of weights of centroids already flushed into out.
	var cumulative float64

	// We assume centroids are already sorted by Mean.
	current := t.centroids[0]

	for i := 1; i < n; i++ {
		next := t.centroids[i]

		// Candidate merged centroid.
		mergedWeight := current.Weight + next.Weight
		mergedMean := (current.Mean*current.Weight + next.Mean*next.Weight) / mergedWeight

		// Quantile of the merged centroid center.
		q := max(min((cumulative+mergedWeight/2)/total, 0.0), 1.0)

		// Maximum allowed weight for this quantile (shared with Add.tryMerge).
		k := t.maxWeightForQuantile(q, total)

		if mergedWeight <= k || len(out) == 0 {
			// Safe to merge next into current.
			current.Mean = mergedMean
			current.Weight = mergedWeight
		} else {
			// Flush current and start a new centroid with next.
			out = append(out, current)
			cumulative += current.Weight
			current = next
		}
	}

	// append the last centroid.
	// Replace centroids with the compressed list and clip capacity.
	t.centroids = slices.Clip(append(out, current))
	// Note: t.count remains unchanged because we only redistributed weights.
}

// Clone returns a deep copy of the t-digest.
func (t *tdigest) Clone() TDigest {
	t.mu.Lock()
	defer t.mu.Unlock()

	newT := &tdigest{
		compression:              t.compression,
		compressionTriggerFactor: t.compressionTriggerFactor,
		count:                    t.count,
	}
	if len(t.centroids) > 0 {
		newT.centroids = slices.Clone(t.centroids)
	}
	if len(t.quantiles) > 0 {
		newT.quantiles = slices.Clone(t.quantiles)
	}
	return newT
}

// Quantiles returns the configured quantiles.
func (t *tdigest) Quantiles() []float64 {
	if t == nil || len(t.quantiles) == 0 {
		return nil
	}
	return slices.Clone(t.quantiles)
}
