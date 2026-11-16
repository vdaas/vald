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

package service

import (
	"context"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func searchResponse(ids ...string) *payload.Search_Response {
	res := &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, len(ids)),
	}
	for _, id := range ids {
		res.Results = append(res.Results, &payload.Object_Distance{Id: id})
	}
	return res
}

// Test_calcRecall pins the delegation to internal/algorithm/recall with the
// linear results as ground truth and k = len(linear results): the values
// reported for the len(search) <= len(linear) shapes produced in practice
// (both requests share the same config Num) are identical to the historical
// matches/len(linear) implementation, and the one divergent edge case —
// more search results than linear results — now follows the normalized
// recall@k definition, ignoring matches beyond position len(linear).
func Test_calcRecall(t *testing.T) {
	for _, tc := range []struct {
		linear *payload.Search_Response
		search *payload.Search_Response
		name   string
		want   float64
	}{
		{name: "nil linear", linear: nil, search: searchResponse("a"), want: 0},
		{name: "nil search", linear: searchResponse("a"), search: nil, want: 0},
		{name: "perfect match", linear: searchResponse("a", "b", "c"), search: searchResponse("c", "a", "b"), want: 1},
		{name: "partial match", linear: searchResponse("a", "b", "c", "d"), search: searchResponse("a", "x", "c", "y"), want: 0.5},
		{
			name:   "search shorter than linear keeps linear denominator",
			linear: searchResponse("a", "b", "c", "d"),
			search: searchResponse("a", "b"),
			want:   0.5,
		},
		{
			name:   "search longer than linear ignores matches beyond the ground-truth count",
			linear: searchResponse("a", "b"),
			search: searchResponse("x", "y", "a", "b"),
			want:   0,
		},
	} {
		t.Run(tc.name, func(tt *testing.T) {
			if got := calcRecall(tc.linear, tc.search); got != tc.want {
				tt.Errorf("calcRecall = %f, want %f", got, tc.want)
			}
		})
	}
}

// Test_meanRecall pins that queries whose search or linear search failed (nil
// slot) are excluded from the average and counted as failed, instead of being
// averaged in as recall 0.
func Test_meanRecall(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name         string
		linear       []*payload.Search_Response
		search       []*payload.Search_Response
		wantMean     float64
		wantMeasured int
		wantFailed   int
	}{
		{name: "empty", wantMean: 0, wantMeasured: 0, wantFailed: 0},
		{
			name:         "all measured",
			linear:       []*payload.Search_Response{searchResponse("a", "b"), searchResponse("c", "d")},
			search:       []*payload.Search_Response{searchResponse("a", "b"), searchResponse("c", "x")},
			wantMean:     0.75,
			wantMeasured: 2,
			wantFailed:   0,
		},
		{
			name:         "failed search and failed linear search are excluded",
			linear:       []*payload.Search_Response{searchResponse("a", "b"), searchResponse("c", "d"), nil},
			search:       []*payload.Search_Response{searchResponse("a", "b"), nil, searchResponse("e", "f")},
			wantMean:     1,
			wantMeasured: 1,
			wantFailed:   2,
		},
		{
			name:         "all failed",
			linear:       []*payload.Search_Response{nil, nil},
			search:       []*payload.Search_Response{nil, searchResponse("a")},
			wantMean:     0,
			wantMeasured: 0,
			wantFailed:   2,
		},
		{
			name:         "linear shorter than search counts the missing slots as failed",
			linear:       []*payload.Search_Response{searchResponse("a")},
			search:       []*payload.Search_Response{searchResponse("a"), searchResponse("b")},
			wantMean:     1,
			wantMeasured: 1,
			wantFailed:   1,
		},
	} {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			mean, measured, failed := meanRecall(tc.linear, tc.search)
			if mean != tc.wantMean || measured != tc.wantMeasured || failed != tc.wantFailed {
				tt.Errorf("meanRecall = (%f, %d, %d), want (%f, %d, %d)",
					mean, measured, failed, tc.wantMean, tc.wantMeasured, tc.wantFailed)
			}
		})
	}
}

// Test_annSearchResult pins that an ANN search Vald answered with NotFound
// (empty result, no response) counts as a recall-0 query in the total search
// recall, while any other error leaves the query excluded as failed.
func Test_annSearchResult(t *testing.T) {
	t.Parallel()
	notFound := status.WrapWithNotFound("error search result length is 0", errors.ErrEmptySearchResult)
	for _, tc := range []struct {
		res          *payload.Search_Response
		err          error
		name         string
		wantMean     float64
		wantMeasured int
		wantFailed   int
	}{
		{name: "result", res: searchResponse("a"), wantMean: 1, wantMeasured: 1},
		{name: "NotFound is recall 0", err: notFound, wantMeasured: 1},
		{name: "wrapped NotFound is recall 0", err: errors.Wrap(notFound, "round robin failed"), wantMeasured: 1},
		{name: "Unavailable is excluded", err: status.WrapWithUnavailable("unavailable", errors.New("conn refused")), wantFailed: 1},
		{name: "canceled is excluded", err: context.Canceled, wantFailed: 1},
		{name: "non-status error is excluded", err: errors.New("boom"), wantFailed: 1},
	} {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			sres := []*payload.Search_Response{annSearchResult(tc.res, tc.err)}
			lres := []*payload.Search_Response{searchResponse("a")}
			mean, measured, failed := meanRecall(lres, sres)
			if mean != tc.wantMean || measured != tc.wantMeasured || failed != tc.wantFailed {
				tt.Errorf("meanRecall = (%f, %d, %d), want (%f, %d, %d)",
					mean, measured, failed, tc.wantMean, tc.wantMeasured, tc.wantFailed)
			}
		})
	}
}
