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

// Package recall is the repository's single implementation of the recall@k
// accuracy metric for nearest-neighbor search results. Both the e2e test
// harness (tests/v2/e2e/metrics, integer hdf5 ground-truth indices) and the
// benchmark job (pkg/tools/benchmark/job/service, string IDs against
// linear-search ground truth) delegate here.
package recall

// Calc reports the recall@k score of an approximate nearest-neighbor search
// result against the ground-truth neighbor IDs for a single query.
//
// got is the ID list returned by the (approximate) search, ordered from the
// nearest neighbor to the farthest. truth is the corresponding ground-truth
// neighbor ID list for the same query — e.g. one row of an
// ann-benchmarks-style hdf5 neighbors dataset, or the IDs returned by an
// exhaustive linear search — also ordered nearest to farthest and possibly
// longer than the k actually requested from the search.
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
// dataset is available), Calc returns 0, since recall is undefined for an
// empty ground truth and 0 is the conservative answer.
func Calc[T comparable](got, truth []T, k int) float64 {
	effectiveK := min(len(truth), k)
	if effectiveK <= 0 {
		return 0
	}
	if len(got) > effectiveK {
		got = got[:effectiveK]
	}
	truth = truth[:effectiveK]

	truthIDs := make(map[T]struct{}, effectiveK)
	for _, id := range truth {
		truthIDs[id] = struct{}{}
	}

	var matched float64
	for _, id := range got {
		if _, ok := truthIDs[id]; ok {
			matched++
		}
	}
	return matched / float64(effectiveK)
}
