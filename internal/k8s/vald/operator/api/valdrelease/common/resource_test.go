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
package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestCalcResource(t *testing.T) {
	tests := []struct {
		name      string
		value     int64
		ratio     float64
		divs      []float64
		wantMilli int64
	}{
		// ratio=0.5 is exactly representable in float64 (=2^-1), so results are exact.
		{"no div", 4, 0.5, nil, 2000},
		{"div=2", 4, 0.5, []float64{2}, 1000},
		{"div=4", 8, 0.5, []float64{4}, 1000},
		// ratio=0.25 is also exact.
		{"0.25 no div", 8, 0.25, nil, 2000},
		{"0.25 div=2", 8, 0.25, []float64{2}, 1000},
		// Scale-up: large CPU values typical of real nodes (int64 safe range).
		{"16 cores × 0.5", 16, 0.5, nil, 8000},
		{"16 cores × 0.5 / 2", 16, 0.5, []float64{2}, 4000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcResource(tt.value, tt.ratio, tt.divs...)
			assert.Equal(t, tt.wantMilli, got.MilliValue())
		})
	}
}

func TestNormalizeResourceList_CPU(t *testing.T) {
	tests := []struct {
		name       string
		inputMilli int64
		wantStr    string
	}{
		{"whole core: 3000m → 3", 3000, "3"},
		{"whole core: 1000m → 1", 1000, "1"},
		{"fractional: 1500m stays", 1500, "1500m"},
		{"fractional: 600m stays", 600, "600m"},
		{"fractional: 4800m → 4800m (not whole)", 4800, "4800m"},
		{"whole core: 6000m → 6", 6000, "6"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := resource.Quantity{}
			q.SetMilli(tt.inputMilli)
			rl := NormalizeResourceList(v1.ResourceList{v1.ResourceCPU: q})
			got := rl[v1.ResourceCPU]
			assert.Equal(t, tt.wantStr, got.String())
		})
	}
}

func TestNormalizeResourceList_Memory(t *testing.T) {
	tests := []struct {
		name       string
		inputBytes int64 // passed via SetMilli(inputBytes * 1000) to simulate CalcResource output
		wantBytes  int64
	}{
		// Values should be rounded up to nearest Mega (10^6 bytes).
		{"3000M exact", 3_000_000_000, 3_000_000_000},
		{"rounds up to next M", 2_999_999_001, 3_000_000_000},
		{"1500M exact", 1_500_000_000, 1_500_000_000},
		{"1500M with sub-M remainder rounds up", 1_499_999_001, 1_500_000_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := resource.Quantity{}
			q.SetMilli(tt.inputBytes * 1000)
			rl := NormalizeResourceList(v1.ResourceList{v1.ResourceMemory: q})
			got := rl[v1.ResourceMemory]
			assert.Equal(t, tt.wantBytes, got.Value())
		})
	}
}
