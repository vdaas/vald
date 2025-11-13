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
	"math"
	"sort"
	"sync"

	"github.com/vdaas/vald/internal/errors"
)

// Centroid represents a centroid in the t-digest.
type Centroid struct {
	Mean   float64
	Weight float64
}

// TDigest is a custom implementation of the t-digest algorithm.
type TDigest struct {
	mu                       sync.Mutex
	centroids                []Centroid
	compression              float64
	compressionTriggerFactor float64
	count                    float64
}

// NewTDigest creates a new TDigest.
func NewTDigest(compression, compressionTriggerFactor float64) (*TDigest, error) {
	return &TDigest{
		compression:              compression,
		compressionTriggerFactor: compressionTriggerFactor,
	}, nil
}

// Add adds a value to the t-digest.
func (t *TDigest) Add(value float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Find the closest centroid
	minDist := math.Inf(1)
	closestIdx := -1
	for i, c := range t.centroids {
		dist := math.Abs(c.Mean - value)
		if dist < minDist {
			minDist = dist
			closestIdx = i
		}
	}

	// If a close enough centroid is found, merge with it
	if closestIdx != -1 {
		c := &t.centroids[closestIdx]
		// The threshold for merging is based on the quantile of the centroid.
		// This is the core idea of the t-digest algorithm.
		q := (c.Weight/2 + t.sumWeightBefore(closestIdx)) / t.count
		k := 4 * t.count * q * (1 - q) / t.compression
		if c.Weight+1 <= k {
			c.Mean = (c.Mean*c.Weight + value) / (c.Weight + 1)
			c.Weight++
			t.count++
			return
		}
	}

	// Otherwise, create a new centroid
	t.centroids = append(t.centroids, Centroid{Mean: value, Weight: 1})
	t.count++

	// If the number of centroids exceeds the compression limit, compress them
	if float64(len(t.centroids)) > t.compression*t.compressionTriggerFactor {
		t.compress()
	}

	// Sort the centroids by mean
	sort.Slice(t.centroids, func(i, j int) bool {
		return t.centroids[i].Mean < t.centroids[j].Mean
	})
}

// Quantile returns the estimated quantile.
func (t *TDigest) Quantile(q float64) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.count == 0 {
		return 0
	}

	target := q * t.count
	var sum float64
	for i, c := range t.centroids {
		if sum+c.Weight > target {
			if i == 0 {
				return c.Mean
			}
			prev := t.centroids[i-1]
			// Linear interpolation
			return prev.Mean + (c.Mean-prev.Mean)*(target-sum)/(c.Weight)
		}
		sum += c.Weight
	}
	return t.centroids[len(t.centroids)-1].Mean
}

// Merge merges another t-digest into this one.
func (t *TDigest) Merge(other QuantileSketch) error {
	if o, ok := other.(*TDigest); ok {
		t.mu.Lock()
		defer t.mu.Unlock()
		o.mu.Lock()
		defer o.mu.Unlock()

		for _, c := range o.centroids {
			// Merge centroid preserving weight
			for i := 0; i < int(c.Weight); i++ {
				t.centroids = append(t.centroids, Centroid{Mean: c.Mean, Weight: 1})
				t.count++
			}
		}
		// Sort and compress after merging all centroids
		sort.Slice(t.centroids, func(i, j int) bool {
			return t.centroids[i].Mean < t.centroids[j].Mean
		})
		if float64(len(t.centroids)) > t.compression*t.compressionTriggerFactor {
			t.compress()
		}

		}
		return nil
	}
	return errors.New("incompatible sketch type for merging")
}

// sumWeightBefore returns the sum of weights of centroids before the given index.
func (t *TDigest) sumWeightBefore(idx int) (sum float64) {
	for i := range idx {
		sum += t.centroids[i].Weight
	}
	return sum
}

// compress merges the centroids to reduce their number.
func (t *TDigest) compress() {
	if len(t.centroids) <= 1 {
		return
	}

	// Sort centroids by mean
	sort.Slice(t.centroids, func(i, j int) bool {
		return t.centroids[i].Mean < t.centroids[j].Mean
	})

	newCentroids := make([]Centroid, 0, int(t.compression))
	newCentroids = append(newCentroids, t.centroids[0])
	cumulativeWeight := t.centroids[0].Weight

	for i := 1; i < len(t.centroids); i++ {
		last := &newCentroids[len(newCentroids)-1]
		current := t.centroids[i]

		// The merging condition is the same as in the Add method.
		q := (cumulativeWeight - last.Weight/2) / t.count
		k := 4 * t.count * q * (1 - q) / t.compression

		if last.Weight+current.Weight <= k {
			last.Mean = (last.Mean*last.Weight + current.Mean*current.Weight) / (last.Weight + current.Weight)
			last.Weight += current.Weight
		} else {
			newCentroids = append(newCentroids, current)
			cumulativeWeight += current.Weight
		}
	}
	t.centroids = newCentroids
}
