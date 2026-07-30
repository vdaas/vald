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
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
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
