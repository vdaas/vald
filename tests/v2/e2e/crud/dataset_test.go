//go:build e2e

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

// package crud provides e2e tests using ann-benchmarks datasets
package crud

import (
	"testing"

	"github.com/vdaas/vald/internal/iter"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

// --------------------------------------------------
// Dataset Slices Construction Using Cycle Iterators |
// --------------------------------------------------

// getDatasetSlices constructs iterators for the Train, Test, and Neighbors datasets.
// For Train and Test, if more samples are requested than are available,
// a NoiseModifier (via a noiseGenerator) is used to add noise on‑the‑fly.
func getDatasetSlices(
	t *testing.T, e *config.Execution,
) (train, test iter.Cycle[[][]float32, []float32], neighbors iter.Cycle[[][]int, []int]) {
	t.Helper()
	if ds == nil || e == nil || e.BaseConfig == nil {
		return nil, nil, nil
	}
	return ds.TrainCycle(e.Num, e.Offset),
		ds.TestCycle(e.Num, e.Offset),
		ds.NeighborsCycle(e.Num, e.Offset)
}
