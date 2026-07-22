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

package metrics

import algorithm "github.com/vdaas/vald/internal/algorithm/recall"

// CalcRecall reports the recall@k score of an approximate nearest-neighbor
// search result against the ground-truth neighbor IDs for a single query.
//
// got is the ID list returned by the (approximate) search, ordered from the
// nearest neighbor to the farthest. truth is the corresponding ground-truth
// neighbor ID list for the same query, e.g. one row of
// tests/v2/e2e/hdf5.Dataset.Neighbors (also ordered nearest to farthest, and
// typically longer than the k actually requested from the search, since
// ann-benchmarks-style hdf5 files usually store the top-100 ground truth
// regardless of the benchmark's k).
//
// recall@k depends only on set membership, not on order:
//
//	effectiveK := min(k, len(truth))
//	recall@k   := |top(got, effectiveK) ∩ top(truth, effectiveK)| / effectiveK
//
// where top(x, n) denotes the first min(n, len(x)) elements of x. In other
// words, both got and truth are truncated down to effectiveK entries before
// being compared as unordered sets — a matching ID that only appears beyond
// position effectiveK in got (e.g. because the search returned more than k
// results) must NOT count towards the score, and ground-truth entries beyond
// effectiveK must NOT be considered "correct" either.
//
// If effectiveK <= 0 (k <= 0, or truth is empty/nil — e.g. no ground-truth
// dataset is available), CalcRecall returns 0 without an error, since recall
// is undefined for an empty ground truth and 0 is the conservative answer.
//
// CalcRecall delegates to internal/algorithm/recall.Calc — the repository's
// single recall implementation, shared with pkg/tools/benchmark — retaining
// this package's historical signature over integer hdf5 ground-truth
// indices.
func CalcRecall(got, truth []int, k int) (recall float64) {
	return algorithm.Calc(got, truth, k)
}
