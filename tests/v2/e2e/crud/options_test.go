//go:build e2e

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

package crud

import (
	"context"
	"slices"
	"strconv"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/iter"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"google.golang.org/protobuf/types/known/anypb"
)

func algorithmOptions(t testing.TB, options config.AlgorithmOptions) []*anypb.Any {
	t.Helper()
	values, err := options.Any()
	if err != nil {
		t.Fatalf("invalid algorithm options: %v", err)
	}
	return values
}

func validateSearchIDs(res *payload.Search_Response, expected []string) error {
	if res == nil {
		return errors.New("missing search response")
	}
	if len(res.Results) != len(expected) {
		return errors.Errorf("expected result IDs %v, got %v", expected, res.Results)
	}
	remaining := make(map[string]struct{}, len(expected))
	for _, id := range expected {
		if _, duplicate := remaining[id]; duplicate {
			return errors.Errorf("duplicate expected result ID %q", id)
		}
		remaining[id] = struct{}{}
	}
	for _, result := range res.Results {
		if result == nil {
			return errors.New("nil search result")
		}
		if _, ok := remaining[result.Id]; !ok {
			return errors.Errorf("unexpected or duplicate result ID %q; expected %v", result.Id, expected)
		}
		delete(remaining, result.Id)
	}
	return nil
}

func TestValidateSearchIDs(t *testing.T) {
	tests := []struct {
		name        string
		expected    []string
		ids         []string
		nilResponse bool
		wantError   bool
	}{
		{name: "order independent", expected: []string{"0", "2"}, ids: []string{"2", "0"}},
		{name: "missing result", expected: []string{"0", "2"}, ids: []string{"0"}, wantError: true},
		{name: "filter ignored", expected: []string{"0"}, ids: []string{"0", "1"}, wantError: true},
		{name: "wrong metadata group", expected: []string{"0"}, ids: []string{"1"}, wantError: true},
		{name: "duplicate result", expected: []string{"0", "2"}, ids: []string{"0", "0"}, wantError: true},
		{name: "duplicate expectation", expected: []string{"0", "0"}, ids: []string{"0", "0"}, wantError: true},
		{name: "empty match", expected: []string{}, ids: []string{}},
		{name: "missing response", nilResponse: true, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var res *payload.Search_Response
			if !tt.nilResponse {
				res = new(payload.Search_Response)
				for _, id := range tt.ids {
					res.Results = append(res.Results, &payload.Object_Distance{Id: id})
				}
			}
			if err := validateSearchIDs(res, tt.expected); (err != nil) != tt.wantError {
				t.Errorf("validateSearchIDs: %v, wantError=%t", err, tt.wantError)
			}
		})
	}
}

func TestAlgorithmRequestBuilders(t *testing.T) {
	t.Run("insert forwards document fields", func(t *testing.T) {
		options := config.AlgorithmOptions{{"@type": "type.googleapis.com/payload.v1.ZVec.DocumentOptions", "fields": map[string]any{"category": "news"}}}
		plan := &config.Execution{Modification: &config.ModificationConfig{ID: "metadata-fixture", Options: options}}
		req := insertRequest(t, 0, "fixture", []float32{1, 2}, plan)
		if req.GetVector().GetId() != "metadata-fixture" {
			t.Fatal("explicit fixture ID was dropped")
		}
		if len(req.GetConfig().GetOptions()) != 1 {
			t.Fatal("document options were dropped")
		}
		doc := new(payload.ZVec_DocumentOptions)
		if err := req.Config.Options[0].UnmarshalTo(doc); err != nil {
			t.Fatal(err)
		}
		if doc.Fields["category"] != "news" {
			t.Errorf("document metadata was lost: %v", doc.Fields)
		}
	})
	t.Run("search forwards scalar and full text options", func(t *testing.T) {
		query := &config.SearchQuery{
			K: 10,
			Options: config.AlgorithmOptions{
				{
					"@type":          "type.googleapis.com/payload.v1.ZVec.SearchOptions",
					"pre_filter":     "category = 'news'",
					"hybrid_queries": []any{map[string]any{"field_name": "content", "fts": map[string]any{"match_string": "galaxy"}}},
				},
			},
		}
		req := newSearchConfig(t, "fixture", query)
		if len(req.Options) != 1 {
			t.Fatal("search options were dropped")
		}
		options := new(payload.ZVec_SearchOptions)
		if err := req.Options[0].UnmarshalTo(options); err != nil {
			t.Fatal(err)
		}
		if options.PreFilter != "category = 'news'" || len(options.HybridQueries) != 1 || options.HybridQueries[0].GetFts().GetMatchString() != "galaxy" {
			t.Errorf("search options were lost: %v", options)
		}
	})
}

func TestMultiBatches(t *testing.T) {
	for _, tt := range []struct {
		name      string
		num, bulk uint64
		sizes     []int
	}{
		{name: "no requests", bulk: 10},
		{name: "full batch without empty tail", num: 10, bulk: 10, sizes: []int{10}},
		{name: "partial batch", num: 9, bulk: 10, sizes: []int{9}},
		{name: "full and partial batches", num: 11, bulk: 10, sizes: []int{1, 10}},
		{name: "explicit single-request batches", num: 3, bulk: 1, sizes: []int{1, 1, 1}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			batches := make(chan []string, int(tt.num)+1)
			call := func(_ context.Context, request *payload.Search_MultiRequest, _ ...grpc.CallOption) (*payload.Search_Responses, error) {
				ids := make([]string, 0, len(request.Requests))
				for _, query := range request.Requests {
					ids = append(ids, query.GetConfig().GetRequestId())
				}
				batches <- ids
				return &payload.Search_Responses{}, nil
			}
			plan := &config.Execution{
				BaseConfig: &config.BaseConfig{Num: tt.num, BulkSize: tt.bulk, Parallelism: 2},
				Search:     &config.SearchQuery{K: 1},
			}
			data := iter.NewCycle([][]float32{{1, 2}}, tt.num, 0, nil)
			if err := multi(t, context.Background(), data, plan, call, searchRequest, searchMultiRequest); err != nil {
				t.Fatal(err)
			}
			close(batches)
			var sizes []int
			var ids []string
			for batch := range batches {
				sizes = append(sizes, len(batch))
				ids = append(ids, batch...)
			}
			slices.Sort(sizes)
			if !slices.Equal(sizes, tt.sizes) {
				t.Errorf("batch sizes=%v, want %v", sizes, tt.sizes)
			}
			slices.Sort(ids)
			expected := make([]string, 0, tt.num)
			for i := uint64(0); i < tt.num; i++ {
				expected = append(expected, strconv.FormatUint(i, 10)+"-Unknown")
			}
			slices.Sort(expected)
			if !slices.Equal(ids, expected) {
				t.Errorf("request IDs=%v, want %v", ids, expected)
			}
		})
	}
}
