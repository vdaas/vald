//go:build e2e

//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// package crud provides end-to-end tests using ann-benchmarks datasets.
package crud

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/iter"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// recall calculates the recall ratio by comparing the list of result IDs
// with the expected neighbors provided as a slice of integers.
// It returns the ratio of matching neighbor IDs to the total number of expected neighbors.
func recall(t *testing.T, resultIDs []string, neighbors []int) float64 {
	t.Helper()
	// Create a set of expected neighbor IDs for fast lookup.
	ns := make(map[string]struct{})
	for _, n := range neighbors {
		ns[strconv.Itoa(n)] = struct{}{}
	}

	// Count how many resultIDs exist in the set of expected neighbor IDs.
	var count int
	for _, r := range resultIDs {
		if _, ok := ns[r]; ok {
			count++
		}
	}
	// Return the recall as a ratio.
	return float64(count) / float64(len(neighbors))
}

// calculateRecall extracts the topK result IDs from the search response and computes the recall.
// It uses the provided index to select the expected neighbor IDs from a global source (ds.Neighbors).
func calculateRecall(t *testing.T, neighbors []int, res *payload.Search_Response) float64 {
	t.Helper()
	// Extract the IDs from the results.
	topKIDs := make([]string, 0, len(res.GetResults()))
	for _, d := range res.GetResults() {
		topKIDs = append(topKIDs, d.GetId())
	}

	// If no results are returned, log an error.
	if len(topKIDs) == 0 {
		t.Errorf("empty result is returned for test ID %s: %#v", res.GetRequestId(), topKIDs)
		return 0
	}
	// ds.Neighbors is assumed to be defined globally with expected neighbor IDs.
	return recall(t, topKIDs, neighbors[:len(topKIDs)])
}

// newSearchConfig creates a new Search_Config instance based on the provided search query and test ID.
// It parses the timeout string into nanoseconds, sets a default timeout if needed, and conditionally sets the ratio.
func newSearchConfig(t *testing.T, id string, query *config.SearchQuery) *payload.Search_Config {
	t.Helper()
	if query == nil {
		t.Errorf("search query is nil")
	}
	return &payload.Search_Config{
		// The RequestId is composed of the test ID and the name of the aggregation algorithm.
		RequestId: id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)],
		Num:       query.K,
		Radius:    query.Radius,
		Epsilon:   query.Epsilon,
		// Parse the timeout value; use 1 second as default if parsing fails or timeout is empty.
		Timeout: func() int64 {
			if query.Timeout != "" {
				timeout, err := query.Timeout.Duration()
				if err == nil {
					return timeout.Nanoseconds()
				}
			}
			return time.Second.Nanoseconds()
		}(),
		AggregationAlgorithm: query.Algorithm,
		MinNum:               query.MinNum,
		// Conditionally set the ratio if it is non-zero.
		Ratio: func() *wrapperspb.FloatValue {
			if query.Ratio != 0 {
				return wrapperspb.Float(query.Ratio)
			}
			return nil
		}(),
		Nprobe: query.Nprobe,
	}
}

// newSearchRequest is a generic type for functions that create search requests.
type newSearchRequest[R proto.Message] func(id string, vec []float32, scfg *payload.Search_Config) R

// Predefined request builder functions for unary and multi search requests.
var (
	// searchRequest builds a Search_Request given a vector and search configuration.
	// The id parameter is ignored in this case.
	searchRequest newRequest[*payload.Search_Request] = func(t *testing.T, _ uint64, id string, vec []float32, e *config.Execution) *payload.Search_Request {
		return &payload.Search_Request{
			Vector: vec,
			Config: newSearchConfig(t, id, e.Search),
		}
	}

	// searchIDRequest builds a Search_IDRequest given an id and search configuration.
	// The vector is ignored for search-by-ID requests.
	searchIDRequest newRequest[*payload.Search_IDRequest] = func(t *testing.T, _ uint64, id string, _ []float32, e *config.Execution) *payload.Search_IDRequest {
		return &payload.Search_IDRequest{
			Id:     id,
			Config: newSearchConfig(t, id, e.Search),
		}
	}

	// searchMultiRequest builds a Search_MultiRequest from a slice of Search_Request.
	searchMultiRequest newMultiRequest[*payload.Search_Request, *payload.Search_MultiRequest] = func(t *testing.T, reqs ...*payload.Search_Request) *payload.Search_MultiRequest {
		return &payload.Search_MultiRequest{
			Requests: reqs,
		}
	}

	// searchMultiIDRequest builds a Search_MultiIDRequest from a slice of Search_IDRequest.
	searchMultiIDRequest newMultiRequest[*payload.Search_IDRequest, *payload.Search_MultiIDRequest] = func(t *testing.T, reqs ...*payload.Search_IDRequest) *payload.Search_MultiIDRequest {
		return &payload.Search_MultiIDRequest{
			Requests: reqs,
		}
	}
)

// processSearch dispatches the search operation based on the type and mode specified in the plan.
// It supports unary, multiple (bulk), and stream operations for both vector search and search-by-ID.
func (r *runner) processSearch(
	t *testing.T,
	ctx context.Context,
	test, train iter.Cycle[[][]float32, []float32],
	neighbors iter.Cycle[[][]int, []int],
	plan *config.Execution,
) {
	t.Helper()
	if plan == nil {
		t.Fatal("search operation plan is nil")
		return
	}

	if plan.BaseConfig == nil {
		t.Fatal("base configuration is nil")
		return
	}
	if plan.Search == nil {
		t.Fatal("search configuration is nil")
		return
	}

	switch plan.Type {
	case config.OpSearch:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			// For unary search requests, use the generic unarySearch function with the searchRequest builder.
			unary(t, ctx, test, plan, r.client.Search, searchRequest, checkUnarySearchResponse(neighbors))
		case config.OperationMultiple:
			// For bulk search requests, use the generic multiSearch function with searchRequest and searchMultiRequest builders.
			multi(t, ctx, test, plan, r.client.MultiSearch, searchRequest, searchMultiRequest, checkMultiSearchResponse(neighbors))
		case config.OperationStream:
			// For streaming search requests, use the generic streamSearch function with the searchRequest builder.
			stream(t, ctx, test, plan, r.client.StreamSearch, searchRequest, checkStreamSearchResponse(neighbors))
		}
	case config.OpSearchByID:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unary(t, ctx, train, plan, r.client.SearchByID, searchIDRequest, checkUnarySearchResponse(neighbors))
		case config.OperationMultiple:
			multi(t, ctx, train, plan, r.client.MultiSearchByID, searchIDRequest, searchMultiIDRequest, checkMultiSearchResponse(neighbors))
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamSearchByID, searchIDRequest, checkStreamSearchResponse(neighbors))
		}
	case config.OpLinearSearch:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unary(t, ctx, test, plan, r.client.LinearSearch, searchRequest, checkUnarySearchResponse(neighbors))
		case config.OperationMultiple:
			multi(t, ctx, test, plan, r.client.MultiLinearSearch, searchRequest, searchMultiRequest, checkMultiSearchResponse(neighbors))
		case config.OperationStream:
			stream(t, ctx, test, plan, r.client.StreamLinearSearch, searchRequest, checkStreamSearchResponse(neighbors))
		}
	case config.OpLinearSearchByID:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unary(t, ctx, test, plan, r.client.LinearSearchByID, searchIDRequest, checkUnarySearchResponse(neighbors))
		case config.OperationMultiple:
			multi(t, ctx, train, plan, r.client.MultiLinearSearchByID, searchIDRequest, searchMultiIDRequest, checkMultiSearchResponse(neighbors))
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamLinearSearchByID, searchIDRequest, checkStreamSearchResponse(neighbors))
		}
	}
}

func checkUnarySearchResponse(
	neighbors iter.Cycle[[][]int, []int],
) func(t *testing.T, idx uint64, res *payload.Search_Response, err error) bool {
	return func(t *testing.T, idx uint64, res *payload.Search_Response, err error) bool {
		rc := calculateRecall(t, neighbors.At(idx), res)
		t.Logf("request id %s searched recall: %f, payload %s", res.GetRequestId(), rc, res.String())
		return true
	}
}

func checkMultiSearchResponse(
	neighbors iter.Cycle[[][]int, []int],
) func(t *testing.T, idx uint64, res *payload.Search_Responses, err error) bool {
	return func(t *testing.T, idx uint64, res *payload.Search_Responses, err error) bool {
		// For each response in the bulk response, log the recall.
		for _, r := range res.GetResponses() {
			if !checkUnarySearchResponse(neighbors)(t, getIndexFromSearchResponse(t, r), r, err) {
				return false
			}
		}
		return true
	}
}

func checkStreamSearchResponse(
	neighbors iter.Cycle[[][]int, []int],
) func(t *testing.T, idx uint64, res *payload.Search_Response, err error) bool {
	return func(t *testing.T, idx uint64, res *payload.Search_Response, err error) bool {
		return checkUnarySearchResponse(neighbors)(t, getIndexFromSearchResponse(t, res), res, err)
	}
}

func getIndexFromSearchResponse(t *testing.T, res *payload.Search_Response) (idx uint64) {
	t.Helper()
	id, _, _ := strings.Cut(res.GetRequestId(), "-")
	var err error
	idx, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		t.Error(err)
	}
	return idx
}
