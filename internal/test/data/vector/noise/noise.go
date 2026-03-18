//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

// package noise provides a noise generator for adding noise to vectors.
package noise

import (
	"math"
	"math/bits"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/vdaas/vald/internal/log"
)

type Func func(i uint64, vec []float32) (res []float32)

// -------------------------------
// Public Interface for  Addition |
// -------------------------------

// Modifier defines the interface for adding noise on‑the‑fly.
// External packages obtain the noise‑adding function via the Mod method.
type Modifier interface {
	// Mod returns a function which, given a sample index and a vector,
	// produces a modified vector with noise added.
	Mod() Func
}

// -------------------------------------------
// noiseGenerator Struct and Receiver Methods |
// -------------------------------------------

// noiseGenerator encapsulates the dataset and all pre‑computed noise parameters.
// During construction, it calculates the optimal noise level, estimates the required noise table size,
// and precomputes a noise table using a fast Gaussian generator. These values remain constant as long as the dataset is unchanged.
type noiseGenerator struct {
	noiseTable               []float32
	noiseTableDivisionFactor uint64
	minNoiseTableSize        uint64
	noiseLevelFactor         float32
}

// New constructs a new noiseGenerator instance using the Functional Option Pattern.
// It applies the default options first, then any user‑provided options.
// It precomputes the noise level and noise table size based on the input dataset and test sample count,
// and then precomputes the noise table.
// It returns a Modifier interface.
func New(data [][]float32, num uint64, opts ...Option) Modifier {
	ng := new(noiseGenerator)
	// Apply default options first, then any additional options.
	for _, opt := range append(defaultOptions, opts...) {
		opt(ng)
	}

	start := time.Now()
	log.Infof("started at %v to precomputes Noise Table: noiseLevelFactor: %v, noiseTableDivisionFactor: %v, minNoiseTableSize: %v",
		start, ng.noiseLevelFactor, ng.noiseTableDivisionFactor, ng.minNoiseTableSize)

	// Precompute the noise level based on the dataset.
	// The noise level is computed as the average standard deviation of all vectors multiplied by noiseLevelFactor.
	noiseLevel := func() float32 {
		if len(data) == 0 {
			return 0.01 // Default if dataset is empty.
		}
		var totalStd float64
		var count int
		for _, vec := range data {
			lv := float32(len(vec))
			if lv == 0 {
				continue
			}
			// Compute the mean of the vector.
			var sum float32
			for _, v := range vec {
				sum += v
			}
			mean := sum / lv
			// Compute variance and standard deviation.
			var varSum float32
			for _, v := range vec {
				diff := v - mean
				varSum += diff * diff
			}
			totalStd += math.Sqrt(float64(varSum / lv))
			count++
		}
		return float32(totalStd/float64(count)) * ng.noiseLevelFactor
	}()

	// Estimate the optimal noise table size.
	// Heuristic: required unique noise samples = (num * vectorDim) / noiseTableDivisionFactor.
	// The size is rounded up to the next power of two, ensuring it is at least minNoiseTableSize.
	noiseTableSize := func() int {
		if len(data) == 0 || len(data[0]) == 0 {
			return 1 << 20 // Fallback default.
		}
		required := num * uint64(len(data[0])) // Total required unique noise samples.
		// Reduce the required noise samples by the division factor.
		required /= ng.noiseTableDivisionFactor
		// Ensure the noise table size is at least minNoiseTableSize.
		if required < ng.minNoiseTableSize {
			required = ng.minNoiseTableSize
		}
		return 1 << bits.Len64(required-1)
	}()

	// Precompute the noise table using fastGaussian32.
	// The noise table is an array of noise samples (each already scaled by the computed noise level),
	// and a larger table reduces periodic artifacts when the same values are reused.
	var (
		haveSpare32 bool
		spare32     float32
	)
	// Preallocate the noise table.
	// The noise table is precomputed to avoid generating noise on‑the‑fly during the test.
	// This is faster and ensures that the same noise values are used for the same sample index.
	// The noise table is a power of two in size to allow for fast modulo indexing.
	ng.noiseTable = make([]float32, noiseTableSize)
	for i := range noiseTableSize {
		ng.noiseTable[i] = func() float32 {
			if haveSpare32 {
				haveSpare32 = false
				return spare32
			}
			var u, v, s float32
			// Generate two random numbers in the range [-1, 1] until s = u*u + v*v is in (0,1).
			for {
				// Use Box-Muller transform to generate two independent standard normal variables.
				// This is faster than using the standard library's Gaussian generator.
				// rand.Float32() returns a random number in the range [0, 1).
				u = rand.Float32()*2 - 1
				v = rand.Float32()*2 - 1
				s = u*u + v*v
				if s > 0 && s < 1 {
					break
				}
			}
			fs := float64(s)
			// Compute multiplier = sqrt(-2 * ln(s) / s) and scale it by the computed noise level.
			multiplier := float32(math.Sqrt(-2*math.Log(fs)/fs)) * noiseLevel
			// Cache a spare sample.
			spare32 = v * multiplier
			// Indicate that a spare sample is available.
			haveSpare32 = true
			// Return the first sample.
			return u * multiplier
		}()
	}
	log.Infof("finished at %v to precomputes Noise Table: noiseTableSize: %d, noiseLevel: %f, noiseTable: %d",
		func() time.Duration {
			return time.Since(start)
		}(), noiseTableSize, noiseLevel, len(ng.noiseTable))

	return ng
}

// Mod implements the Modifier interface.
// It returns a function that, when called with a sample index and a vector,
// produces a modified vector by adding noise values from the precomputed noise table.
// The noise is selected deterministically based on the sample index.
func (ng *noiseGenerator) Mod() Func {
	// Clone the noise table so that the mod function uses a copy (if needed).
	noiseTable := slices.Clone(ng.noiseTable)
	tableSize := uint64(len(noiseTable))
	return func(i uint64, vec []float32) (res []float32) {
		// Clone the input vector to avoid mutating the original.
		res = slices.Clone(vec)
		n := uint64(len(res))
		baseIdx := i * n // Precompute the base index.
		for j := range n {
			res[j] += noiseTable[(baseIdx+j)%tableSize]
		}
		return res
	}
}
