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
	"github.com/vdaas/vald/tests/v2/e2e/metrics"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// resultIDsToInts converts the string result IDs returned by a search
// response into ints so they can be compared against
// tests/v2/e2e/hdf5.Dataset.Neighbors (which are already []int) via
// metrics.CalcRecall. IDs that fail to parse are dropped (and reported via
// t.Error) rather than aborting the whole recall calculation, since a single
// malformed ID must not hide the recall signal of the rest of the response.
func resultIDsToInts(t testing.TB, ids []string) []int {
	t.Helper()
	out := make([]int, 0, len(ids))
	for _, id := range ids {
		n, err := strconv.Atoi(id)
		if err != nil {
			t.Errorf("failed to parse result ID %q as int for recall calculation: %v", id, err)
			continue
		}
		out = append(out, n)
	}
	return out
}

// calculateRecall extracts the topK result IDs from the search response and
// computes the recall@k against neighbors (one row of
// tests/v2/e2e/hdf5.Dataset.Neighbors) using metrics.CalcRecall, which
// clamps k down to len(neighbors) and truncates both sides consistently
// (see tests/v2/e2e/metrics/recall.go for the exact semantics).
func calculateRecall(
	t testing.TB, idx uint64, neighbors []int, res *payload.Search_Response,
) float64 {
	t.Helper()
	results := res.GetResults()
	// If no results are returned, log an error. The request index is the only
	// reliable identifier here: on a failed request the response payload (and
	// therefore its request id) is empty.
	if len(results) == 0 {
		t.Errorf("empty search result for request idx %d, request id %q, response: %s", idx, res.GetRequestId(), res.String())
		return 0
	}
	topKIDs := make([]string, 0, len(results))
	for _, d := range results {
		topKIDs = append(topKIDs, d.GetId())
	}
	got := resultIDsToInts(t, topKIDs)
	return metrics.CalcRecall(got, neighbors, len(results))
}

// newSearchConfig creates a new Search_Config instance based on the provided search query and test ID.
// It parses the timeout string into nanoseconds, sets a default timeout if needed, and conditionally sets the ratio.
func newSearchConfig(t testing.TB, id string, query *config.SearchQuery) *payload.Search_Config {
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
	searchRequest newRequest[*payload.Search_Request] = func(t testing.TB, _ uint64, id string, vec []float32, e *config.Execution) *payload.Search_Request {
		return &payload.Search_Request{
			Vector: vec,
			Config: newSearchConfig(t, id, e.Search),
		}
	}

	// searchIDRequest builds a Search_IDRequest given an id and search configuration.
	// The vector is ignored for search-by-ID requests.
	searchIDRequest newRequest[*payload.Search_IDRequest] = func(t testing.TB, _ uint64, id string, _ []float32, e *config.Execution) *payload.Search_IDRequest {
		return &payload.Search_IDRequest{
			Id:     id,
			Config: newSearchConfig(t, id, e.Search),
		}
	}

	// searchMultiRequest builds a Search_MultiRequest from a slice of Search_Request.
	searchMultiRequest newMultiRequest[*payload.Search_Request, *payload.Search_MultiRequest] = func(t testing.TB, reqs ...*payload.Search_Request) *payload.Search_MultiRequest {
		return &payload.Search_MultiRequest{
			Requests: reqs,
		}
	}

	// searchMultiIDRequest builds a Search_MultiIDRequest from a slice of Search_IDRequest.
	searchMultiIDRequest newMultiRequest[*payload.Search_IDRequest, *payload.Search_MultiIDRequest] = func(t testing.TB, reqs ...*payload.Search_IDRequest) *payload.Search_MultiIDRequest {
		return &payload.Search_MultiIDRequest{
			Requests: reqs,
		}
	}
)

// processSearch dispatches the search operation based on the type and mode specified in the plan.
// It supports unary, multiple (bulk), and stream operations for both vector search and search-by-ID.
func (r *runner) processSearch(
	t testing.TB,
	ctx context.Context,
	test, train iter.Cycle[[][]float32, []float32],
	neighbors iter.Cycle[[][]int, []int],
	plan *config.Execution,
) error {
	t.Helper()
	if plan == nil {
		t.Fatal("search operation plan is nil")
		return nil
	}

	if plan.BaseConfig == nil {
		t.Fatal("base configuration is nil")
		return nil
	}
	if plan.Search == nil {
		t.Fatal("search configuration is nil")
		return nil
	}

	switch plan.Type {
	case config.OpSearch:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			// For unary search requests, use the generic unarySearch function with the searchRequest builder.
			return unary(t, ctx, test, plan, r.client.Search, searchRequest, checkUnarySearchResponse(neighbors, plan))
		case config.OperationMultiple:
			// For bulk search requests, use the generic multiSearch function with searchRequest and searchMultiRequest builders.
			return multi(t, ctx, test, plan, r.client.MultiSearch, searchRequest, searchMultiRequest, checkMultiSearchResponse(neighbors, plan))
		case config.OperationStream:
			// For streaming search requests, use the generic streamSearch function with the searchRequest builder.
			stream(t, ctx, test, plan, r.client.StreamSearch, searchRequest, checkStreamSearchResponse(neighbors, plan))
		}
	case config.OpSearchByID:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			return unary(t, ctx, train, plan, r.client.SearchByID, searchIDRequest, checkUnarySearchResponse(neighbors, plan))
		case config.OperationMultiple:
			return multi(t, ctx, train, plan, r.client.MultiSearchByID, searchIDRequest, searchMultiIDRequest, checkMultiSearchResponse(neighbors, plan))
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamSearchByID, searchIDRequest, checkStreamSearchResponse(neighbors, plan))
		}
	case config.OpLinearSearch:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			return unary(t, ctx, test, plan, r.client.LinearSearch, searchRequest, checkUnarySearchResponse(neighbors, plan))
		case config.OperationMultiple:
			return multi(t, ctx, test, plan, r.client.MultiLinearSearch, searchRequest, searchMultiRequest, checkMultiSearchResponse(neighbors, plan))
		case config.OperationStream:
			stream(t, ctx, test, plan, r.client.StreamLinearSearch, searchRequest, checkStreamSearchResponse(neighbors, plan))
		}
	case config.OpLinearSearchByID:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			// The train cycle is passed for consistency with the other ByID
			// branches; searchIDRequest ignores the vector and derives the ID
			// from the request index alone.
			return unary(t, ctx, train, plan, r.client.LinearSearchByID, searchIDRequest, checkUnarySearchResponse(neighbors, plan))
		case config.OperationMultiple:
			return multi(t, ctx, train, plan, r.client.MultiLinearSearchByID, searchIDRequest, searchMultiIDRequest, checkMultiSearchResponse(neighbors, plan))
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamLinearSearchByID, searchIDRequest, checkStreamSearchResponse(neighbors, plan))
		}
	}
	return nil
}

// checkUnarySearchResponse returns a callback that logs the recall@k of each
// search response and, when plan.Metrics is enabled, records it into
// plan.Collector via recordRecall (see recall_qps_test.go) so it can be
// exposed alongside QPS once the strategy/operation/execution finishes.
func checkUnarySearchResponse(
	neighbors iter.Cycle[[][]int, []int], plan *config.Execution,
) func(t testing.TB, idx uint64, res *payload.Search_Response, err error) bool {
	return func(t testing.TB, idx uint64, res *payload.Search_Response, err error) bool {
		t.Helper()
		// A transport/gRPC error means there is no payload to score; report
		// the real error with its request context instead of letting the
		// empty response masquerade as a recall-0 result.
		if err != nil {
			t.Errorf("%s request idx %d (type: %s, mode: %s) failed: %v", plan.Name, idx, plan.Type, plan.Mode, err)
			recordRecall(plan, 0)
			return true
		}
		rc := calculateRecall(t, idx, neighbors.At(idx), res)
		t.Logf("request id %s searched recall: %f, payload %s", res.GetRequestId(), rc, res.String())
		recordRecall(plan, rc)
		return true
	}
}

func checkMultiSearchResponse(
	neighbors iter.Cycle[[][]int, []int], plan *config.Execution,
) func(t testing.TB, idx uint64, res *payload.Search_Responses, err error) bool {
	return func(t testing.TB, idx uint64, res *payload.Search_Responses, err error) bool {
		t.Helper()
		// For each response in the bulk response, log the recall.
		for _, r := range res.GetResponses() {
			if !checkUnarySearchResponse(neighbors, plan)(t, getIndexFromSearchResponse(t, r), r, err) {
				return false
			}
		}
		return true
	}
}

func checkStreamSearchResponse(
	neighbors iter.Cycle[[][]int, []int], plan *config.Execution,
) func(t testing.TB, idx uint64, res *payload.Search_StreamResponse, err error) bool {
	return func(t testing.TB, idx uint64, res *payload.Search_StreamResponse, err error) bool {
		t.Helper()
		if err != nil {
			t.Errorf("%s stream request idx %d (type: %s, mode: %s) failed: %v, status: %s",
				plan.Name, idx, plan.Type, plan.Mode, err, res.GetStatus().String())
			// A failed stream item has no payload to score; record it as
			// recall 0 so it stays in the denominator (mirrors the unary
			// path), otherwise dropped failures inflate the reported recall.
			recordRecall(plan, 0)
			return true
		}
		r := res.GetResponse()
		if r == nil {
			// The gateway's StreamSearch handler reports a per-item failure
			// (e.g. a real per-request timeout) by sending a
			// Search_StreamResponse whose oneof holds a Status rather than a
			// Response, while the RPC itself stays healthy (err above is
			// nil in this case). Surface that real status/code instead of
			// guessing "it can be timeout".
			if st := res.GetStatus(); st != nil {
				t.Errorf("%s stream request idx %d (type: %s, mode: %s) returned no result, status: %s",
					plan.Name, idx, plan.Type, plan.Mode, st.String())
			} else {
				t.Errorf("%s stream request idx %d (type: %s, mode: %s) returned a nil response, it can be timeout",
					plan.Name, idx, plan.Type, plan.Mode)
			}
			// No scorable payload here either; count it as recall 0.
			recordRecall(plan, 0)
			return true
		}
		return checkUnarySearchResponse(neighbors, plan)(t, getIndexFromSearchResponse(t, r), r, err)
	}
}

func getIndexFromSearchResponse(t testing.TB, res *payload.Search_Response) (idx uint64) {
	t.Helper()
	if res == nil {
		t.Error("search response is nil")
		return idx
	}
	id, _, _ := strings.Cut(res.GetRequestId(), "-")
	var err error
	idx, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		t.Errorf("failed to parse request id %s: %v", id, err)
		return idx
	}
	return idx
}
