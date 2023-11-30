// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package slices

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func Sort[E constraints.Ordered](x []E) {
	slices.Sort(x)
}

func SortFunc[E any](x []E, less func(left, right E) bool) {
	slices.SortFunc(x, less)
}

func SortStableFunc[E any](x []E, less func(left, right E) bool) {
	slices.SortStableFunc(x, less)
}

func RemoveDuplicates[E comparable](x []E, less func(left, right E) bool) []E {
	if len(x) < 2 {
		return x
	}
	SortStableFunc(x, less)
	up := 0 // uniqPointer
	for i := 1; i < len(x); i++ {
		if x[up] != x[i] {
			up++
			x[up] = x[i]
		}
	}
	return x[:up+1]
}
