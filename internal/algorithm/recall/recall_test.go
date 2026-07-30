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

package recall

import "testing"

func TestCalc_Ints(t *testing.T) {
	for _, tc := range []struct {
		name  string
		got   []int
		truth []int
		k     int
		want  float64
	}{
		{name: "perfect match", got: []int{1, 2, 3}, truth: []int{3, 2, 1}, k: 3, want: 1},
		{name: "partial match", got: []int{1, 2, 9}, truth: []int{1, 2, 3}, k: 3, want: 2.0 / 3.0},
		{name: "no match", got: []int{7, 8, 9}, truth: []int{1, 2, 3}, k: 3, want: 0},
		{name: "truth longer than k truncates", got: []int{1, 99}, truth: []int{1, 2, 3, 4}, k: 2, want: 0.5},
		{name: "got beyond effectiveK ignored", got: []int{9, 8, 1}, truth: []int{1, 2}, k: 2, want: 0},
		{name: "empty truth", got: []int{1}, truth: nil, k: 3, want: 0},
		{name: "k zero", got: []int{1}, truth: []int{1}, k: 0, want: 0},
		{name: "short truth caps denominator", got: []int{1, 2, 3}, truth: []int{1}, k: 10, want: 1},
	} {
		t.Run(tc.name, func(tt *testing.T) {
			if got := Calc(tc.got, tc.truth, tc.k); got != tc.want {
				tt.Errorf("Calc(%v, %v, %d) = %f, want %f", tc.got, tc.truth, tc.k, got, tc.want)
			}
		})
	}
}

// TestCalc_Strings pins the generic instantiation the benchmark job uses:
// string result IDs scored against linear-search ground truth with
// k = len(truth), which reproduces the historical matches/len(linear)
// denominator.
func TestCalc_Strings(t *testing.T) {
	truth := []string{"a", "b", "c", "d"}
	got := []string{"a", "c", "x", "y"}
	if r := Calc(got, truth, len(truth)); r != 0.5 {
		t.Errorf("Calc over string IDs = %f, want 0.5", r)
	}
}
