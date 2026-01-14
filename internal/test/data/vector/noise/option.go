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

// ---------------------------------------
// Functional Options and Default Options |
// ---------------------------------------

// Option is a functional option for configuring a noiseGenerator.
type Option func(*noiseGenerator)

// defaultOptions holds the default configuration for noiseGenerator.
var defaultOptions = []Option{
	// Set default noise level factor to 10% (i.e., 0.1)
	WithLevelFactor(0.1),
	// Set default noise table division factor to 10.
	WithTableDivisionFactor(10),
	// Set default minimum noise table size to 1024.
	WithMinTableSize(1024),
}

// WithLevelFactor sets the fraction of the average standard deviation used as the noise level.
func WithLevelFactor(f float32) Option {
	return func(ng *noiseGenerator) {
		ng.noiseLevelFactor = f
	}
}

// WithTableDivisionFactor sets the division factor used when sizing the noise table.
func WithTableDivisionFactor(f uint64) Option {
	return func(ng *noiseGenerator) {
		ng.noiseTableDivisionFactor = f
	}
}

// WithMinTableSize sets the minimum allowed size for the noise table.
func WithMinTableSize(s uint64) Option {
	return func(ng *noiseGenerator) {
		ng.minNoiseTableSize = s
	}
}
