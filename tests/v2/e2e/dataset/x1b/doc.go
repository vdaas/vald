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

// package x1b provides cluster-independent parsers for the fvecs/bvecs/ivecs
// raw binary vector formats used by billion-scale ANN benchmark datasets
// (SIFT1B, DEEP1B, etc.), ported from hack/benchmark/assets/x1b for use by
// tests/v2/e2e.
//
// Unlike hack/benchmark/assets/x1b, which mmaps a file for lazy, random-access
// reads suited to out-of-core benchmarking of billion-scale datasets, this
// package eagerly decodes an io.Reader into in-memory [][]T slices. That
// trades memory-boundedness for a pure, deterministic, cgo-free function
// signature that plugs into the same unit-test contract as its sibling
// tests/v2/e2e/metrics package: no `//go:build e2e` tag and no cluster or
// gonum/hdf5 cgo dependency, so `go test ./tests/v2/e2e/dataset/x1b/...` runs
// standalone. Wiring the parsed slices into tests/v2/e2e/hdf5.Dataset (or a
// billion-scale-aware alternative) is left to the implementer of the
// corresponding GREEN-phase change.
package x1b
