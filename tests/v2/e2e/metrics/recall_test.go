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

package metrics

import (
	"math"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test"
)

// recallTolerance is the floating point comparison tolerance used across the
// table-driven CalcRecall cases below. float64 division of small integers
// (e.g. 2/3) does not always compare exactly equal across expression forms,
// so an epsilon compare is used instead of reflect.DeepEqual/==.
const recallTolerance = 1e-9

func TestCalcRecall(t *testing.T) {
	type args struct {
		got   []int
		truth []int
		k     int
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (float64, error) {
		t.Helper()
		return CalcRecall(args.got, args.truth, args.k), nil
	}, []test.Case[float64, args]{
		{
			// Full ID-set overlap between got and truth must yield a perfect
			// recall@k score of 1.0.
			Name: "exact match returns full recall",
			Args: args{
				got:   []int{1, 2, 3},
				truth: []int{1, 2, 3},
				k:     3,
			},
			Want: test.Result[float64]{Val: 1.0},
		},
		{
			// recall@k is a set-membership metric: search results and ground
			// truth are both ordered by distance, but permuting either order
			// must not change the score.
			Name: "order does not affect recall",
			Args: args{
				got:   []int{3, 1, 2},
				truth: []int{1, 2, 3},
				k:     3,
			},
			Want: test.Result[float64]{Val: 1.0},
		},
		{
			// 2 of the 3 ground-truth IDs are present in got: recall = 2/3.
			Name: "partial match returns m over k",
			Args: args{
				got:   []int{1, 2, 99},
				truth: []int{1, 2, 3},
				k:     3,
			},
			Want: test.Result[float64]{Val: 2.0 / 3.0},
		},
		{
			// No overlap at all between got and truth must yield 0.
			Name: "complete mismatch returns zero",
			Args: args{
				got:   []int{100, 101, 102},
				truth: []int{1, 2, 3},
				k:     3,
			},
			Want: test.Result[float64]{Val: 0.0},
		},
		{
			// truth is longer than k (as real hdf5 ground-truth rows usually
			// are: they store far more neighbors than any single benchmark's
			// k). Only the first k=2 ground-truth IDs ({9, 1}) may count as
			// correct; "2" is truth[2] and must be ignored even though got
			// contains it.
			Name: "k truncates a longer ground truth list",
			Args: args{
				got:   []int{1, 2},
				truth: []int{9, 1, 2, 3, 4},
				k:     2,
			},
			Want: test.Result[float64]{Val: 0.5},
		},
		{
			// got is longer than k (the search returned more results than
			// requested/expected). Only got's first k=2 entries ({100, 101})
			// may be counted, so the later "1" and "2" matches — which do
			// appear in truth — must NOT be credited.
			Name: "match beyond effective k in got must not count",
			Args: args{
				got:   []int{100, 101, 1, 2},
				truth: []int{1, 2},
				k:     2,
			},
			Want: test.Result[float64]{Val: 0.0},
		},
		{
			// got has fewer entries than k: only 1 of the 3 ground-truth IDs
			// can possibly be matched, so recall = 1/3, not 1/1.
			Name: "fewer got entries than k lowers the score",
			Args: args{
				got:   []int{1},
				truth: []int{1, 2, 3},
				k:     3,
			},
			Want: test.Result[float64]{Val: 1.0 / 3.0},
		},
		{
			// k larger than len(truth) must clamp effectiveK down to
			// len(truth) = 2 (both for truth and for truncating got), not
			// divide by the requested k=10.
			Name: "k larger than truth length is clamped down",
			Args: args{
				got:   []int{1, 2, 3, 4, 5},
				truth: []int{1, 2},
				k:     10,
			},
			Want: test.Result[float64]{Val: 1.0},
		},
		{
			Name: "k is zero returns zero",
			Args: args{
				got:   []int{1, 2, 3},
				truth: []int{1, 2, 3},
				k:     0,
			},
			Want: test.Result[float64]{Val: 0.0},
		},
		{
			Name: "negative k returns zero",
			Args: args{
				got:   []int{1, 2, 3},
				truth: []int{1, 2, 3},
				k:     -5,
			},
			Want: test.Result[float64]{Val: 0.0},
		},
		{
			// No ground-truth dataset available for this query (e.g. hdf5
			// neighbors missing/nil) must not panic and must return 0.
			Name: "nil ground truth returns zero",
			Args: args{
				got:   []int{1, 2, 3},
				truth: nil,
				k:     3,
			},
			Want: test.Result[float64]{Val: 0.0},
		},
		{
			// An empty search result against a non-empty ground truth is
			// zero overlap, not an error.
			Name: "nil got returns zero overlap",
			Args: args{
				got:   nil,
				truth: []int{1, 2, 3},
				k:     3,
			},
			Want: test.Result[float64]{Val: 0.0},
		},
		{
			Name: "both got and truth nil returns zero",
			Args: args{
				got:   nil,
				truth: nil,
				k:     5,
			},
			Want: test.Result[float64]{Val: 0.0},
		},
	}...); err != nil {
		t.Error(err)
	}
}

// TestCalcRecall_ToleranceCompare exercises CalcRecall's CheckFunc plumbing
// directly (outside of test.Run's DefaultCheck, which uses reflect.DeepEqual
// and would be too strict for float64 fractions such as 2/3). It documents
// that CalcRecall must return values within recallTolerance of the expected
// fraction, not necessarily bit-identical values.
func TestCalcRecall_ToleranceCompare(t *testing.T) {
	type args struct {
		got   []int
		truth []int
		k     int
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (float64, error) {
		t.Helper()
		return CalcRecall(args.got, args.truth, args.k), nil
	}, []test.Case[float64, args]{
		{
			Name: "one third recall within tolerance",
			Args: args{
				got:   []int{1},
				truth: []int{1, 2, 3},
				k:     3,
			},
			Want: test.Result[float64]{Val: 1.0 / 3.0},
			CheckFunc: func(t *testing.T, want test.Result[float64], got test.Result[float64]) error {
				t.Helper()
				if got.Err != nil {
					return errors.Errorf("unexpected error: %v", got.Err)
				}
				if diff := math.Abs(got.Val - want.Val); diff > recallTolerance {
					return errors.Errorf("CalcRecall() = %v, want ~%v (diff %v > tolerance %v)", got.Val, want.Val, diff, recallTolerance)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}
