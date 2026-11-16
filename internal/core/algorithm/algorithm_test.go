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

import (
	"math"
	"testing"
)

func TestHasNaNOrInf(t *testing.T) {
	t.Parallel()
	nan := float32(math.NaN())
	pinf := float32(math.Inf(1))
	ninf := float32(math.Inf(-1))
	tests := []struct {
		name string
		vec  []float32
		want bool
	}{
		{name: "nil vector", vec: nil, want: false},
		{name: "empty vector", vec: []float32{}, want: false},
		{name: "finite vector", vec: []float32{0, 1, -1, 1.5}, want: false},
		{name: "max and min finite values", vec: []float32{math.MaxFloat32, -math.MaxFloat32, math.SmallestNonzeroFloat32}, want: false},
		{name: "negative zero", vec: []float32{float32(math.Copysign(0, -1))}, want: false},
		{name: "NaN only", vec: []float32{nan}, want: true},
		{name: "+Inf only", vec: []float32{pinf}, want: true},
		{name: "-Inf only", vec: []float32{ninf}, want: true},
		{name: "NaN at head", vec: []float32{nan, 1, 2}, want: true},
		{name: "+Inf at middle", vec: []float32{1, pinf, 2}, want: true},
		{name: "-Inf at tail", vec: []float32{1, 2, ninf}, want: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := HasNaNOrInf(tc.vec); got != tc.want {
				t.Errorf("HasNaNOrInf(%v) = %v, want %v", tc.vec, got, tc.want)
			}
		})
	}
}

func BenchmarkHasNaNOrInf(b *testing.B) {
	vec := make([]float32, 1024)
	for i := range vec {
		vec[i] = float32(i)
	}
	b.ReportAllocs()
	for b.Loop() {
		if HasNaNOrInf(vec) {
			b.Fatal("unexpected NaN or Inf")
		}
	}
}
