// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package algorithm

import "math"

const (
	// MaximumVectorDimensionSize is maximum value of vector dimension.
	MaximumVectorDimensionSize = math.MaxUint32
	// MinimumVectorDimensionSize is minimum value of vector dimension.
	MinimumVectorDimensionSize = 2
)

// HasNaNOrInf reports whether vec contains NaN, +Inf or -Inf.
// v-v is 0 for every finite v and NaN for NaN and +-Inf, so a single
// comparison covers all of them without a function call per element.
func HasNaNOrInf(vec []float32) bool {
	for _, v := range vec {
		if v-v != 0 {
			return true
		}
	}
	return false
}

// SearchResult is struct for comfortable use in Go.
type SearchResult struct {
	Error    error
	ID       uint32
	Distance float32
}
