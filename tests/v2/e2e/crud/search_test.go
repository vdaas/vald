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
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/iter"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
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
		t.Error("search query is nil")
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
	searchRequest newSearchRequest[*payload.Search_Request] = func(_ string, vec []float32, scfg *payload.Search_Config) *payload.Search_Request {
		return &payload.Search_Request{
			Vector: vec,
			Config: scfg,
		}
	}

	// searchIDRequest builds a Search_IDRequest given an id and search configuration.
	// The vector is ignored for search-by-ID requests.
	searchIDRequest newSearchRequest[*payload.Search_IDRequest] = func(id string, _ []float32, scfg *payload.Search_Config) *payload.Search_IDRequest {
		return &payload.Search_IDRequest{
			Id:     id,
			Config: scfg,
		}
	}

	// searchMultiRequest builds a Search_MultiRequest from a slice of Search_Request.
	searchMultiRequest newMultiRequest[*payload.Search_Request, *payload.Search_MultiRequest] = func(reqs []*payload.Search_Request) *payload.Search_MultiRequest {
		return &payload.Search_MultiRequest{
			Requests: reqs,
		}
	}

	// searchMultiIDRequest builds a Search_MultiIDRequest from a slice of Search_IDRequest.
	searchMultiIDRequest newMultiRequest[*payload.Search_IDRequest, *payload.Search_MultiIDRequest] = func(reqs []*payload.Search_IDRequest) *payload.Search_MultiIDRequest {
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
			unarySearch(t, ctx, test, neighbors, plan, r.client.Search, searchRequest)
		case config.OperationMultiple:
			// For bulk search requests, use the generic multiSearch function with searchRequest and searchMultiRequest builders.
			multiSearch(t, ctx, test, neighbors, plan, r.client.MultiSearch, searchRequest, searchMultiRequest)
		case config.OperationStream:
			// For streaming search requests, use the generic streamSearch function with the searchRequest builder.
			streamSearch(t, ctx, test, neighbors, plan, r.client.StreamSearch, searchRequest)
		}
	case config.OpSearchByID:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unarySearch(t, ctx, train, neighbors, plan, r.client.SearchByID, searchIDRequest)
		case config.OperationMultiple:
			multiSearch(t, ctx, train, neighbors, plan, r.client.MultiSearchByID, searchIDRequest, searchMultiIDRequest)
		case config.OperationStream:
			streamSearch(t, ctx, train, neighbors, plan, r.client.StreamSearchByID, searchIDRequest)
		}
	case config.OpLinearSearch:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unarySearch(t, ctx, test, neighbors, plan, r.client.LinearSearch, searchRequest)
		case config.OperationMultiple:
			multiSearch(t, ctx, test, neighbors, plan, r.client.MultiLinearSearch, searchRequest, searchMultiRequest)
		case config.OperationStream:
			streamSearch(t, ctx, test, neighbors, plan, r.client.StreamLinearSearch, searchRequest)
		}
	case config.OpLinearSearchByID:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unarySearch(t, ctx, train, neighbors, plan, r.client.LinearSearchByID, searchIDRequest)
		case config.OperationMultiple:
			multiSearch(t, ctx, train, neighbors, plan, r.client.MultiLinearSearchByID, searchIDRequest, searchMultiIDRequest)
		case config.OperationStream:
			streamSearch(t, ctx, train, neighbors, plan, r.client.StreamLinearSearchByID, searchIDRequest)
		}
	}
}

// unarySearch handles unary search requests. It iterates over the data and for each vector,
// it sends a search request with each search configuration specified in the plan.
func unarySearch[R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	neighbors iter.Cycle[[][]int, []int],
	plan *config.Execution,
	call grpcCall[R, *payload.Search_Response],
	newReq newSearchRequest[R],
) {
	t.Helper()
	// Create an error group to manage concurrent requests.
	eg, ctx := errgroup.New(ctx)
	// Set the concurrency limit from the plan configuration.
	if plan != nil && plan.BaseConfig != nil {
		// Set the concurrency limit from the plan configuration.
		eg.SetLimit(int(plan.Parallelism))
	}
	for i, vec := range data.Seq2(ctx) {
		idx := i
		// For each test vector, iterate over all search configurations.
		for _, q := range plan.Search {
			id := strconv.FormatUint(idx, 10)
			// Build the search configuration for this request.
			scfg := newSearchConfig(t, id, q)
			// Launch the search request in a goroutine.
			eg.Go(func() error {
				// Execute the search gRPC call.
				res, err := call(ctx, newReq(id, vec, scfg))
				if err != nil {
					// Handle the error using the centralized error handler.
					handleGRPCCallError(t, err, plan)
					return nil
				}
				// Log the result including calculated recall.
				t.Logf("vector %v id %s searched recall: %f, payload %s",
					vec, scfg.RequestId, calculateRecall(t, neighbors.At(idx), res), res.String())
				return nil
			})
		}
	}
	// Wait for all goroutines to complete.
	eg.Wait()
}

// multiSearch handles bulk search requests by grouping individual requests up to BulkSize.
// Once the bulk size is reached, it sends the grouped requests and logs the responses.
// It uses the provided builder functions to create the individual requests and the bulk request.
// The bulk request is sent using the provided gRPC call function.
// The function logs the response for each batch of requests.
// The function is used for vector search, searchByID, linearSearch, and linearSearchByID operations.
func multiSearch[Q, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	neighbors iter.Cycle[[][]int, []int],
	plan *config.Execution,
	call grpcCall[R, *payload.Search_Responses],
	addReqs newSearchRequest[Q],
	toReq newMultiRequest[Q, R],
) {
	t.Helper()
	eg, ctx := errgroup.New(ctx)
	eg.SetLimit(int(plan.Parallelism))
	exec := func(batch []Q) {
		// Convert the slice of individual requests into a bulk request.
		// Execute the bulk request using the provided gRPC call function.
		res, err := call(ctx, toReq(batch))
		if err != nil {
			handleGRPCCallError(t, err, plan)
			return
		}
		// For each response in the bulk response, log the recall.
		for _, r := range res.GetResponses() {
			id, _, _ := strings.Cut(r.GetRequestId(), "-")
			idx, _ := strconv.Atoi(id)
			t.Logf("id %s searched recall: %f, payload %s",
				r.GetRequestId(),
				calculateRecall(t, neighbors.At(uint64(idx)), &payload.Search_Response{
					RequestId: r.GetRequestId(),
					Results:   r.GetResults(),
				}), res.String())
		}
	}

	// Initialize a slice to hold the bulk requests.
	reqs := make([]Q, 0, plan.BulkSize)
	for i, vec := range data.Seq2(ctx) {
		for _, query := range plan.Search {
			id := strconv.FormatUint(i, 10)
			// Append a new request to the bulk slice using the provided builder.
			reqs = append(reqs, addReqs(id, vec, newSearchConfig(t, id, query)))
			// If the bulk size is reached, send the batch.
			if len(reqs) >= int(plan.BulkSize) {
				// Capture the current batch.
				batch := slices.Clone(reqs)
				// Reset the bulk request slice for the next batch.
				reqs = reqs[:0]
				eg.Go(func() error {
					exec(batch)
					return nil
				})
			}
		}
	}
	exec(reqs)

	eg.Wait()
}

// streamSearch handles bidirectional streaming search requests.
// It repeatedly sends search requests from the data slice using the provided builder,
// and processes each response received from the stream.
func streamSearch[S grpc.ClientStream, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	neighbors iter.Cycle[[][]int, []int],
	plan *config.Execution,
	newStream newStream[S],
	newReq newSearchRequest[R],
) {
	t.Helper()
	// Create a new stream using the provided stream function.
	stream, err := newStream(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	// qidx tracks the current index within the search configuration slice.
	// idx tracks the current vector index.
	qidx, idx := uint64(0), uint64(0)
	// Use a bidirectional stream client to send requests and receive responses.
	err = grpc.BidirectionalStreamClient(stream, int(plan.Parallelism), func() *R {
		// If we have processed all vectors, return nil to close the stream.
		if idx >= data.Len() {
			return nil
		}
		id := strconv.FormatUint(idx, 10)
		// If we have exhausted the search configurations, move to the next vector.
		if qidx >= uint64(len(plan.Search)) {
			qidx = 0
			idx++
			if idx >= data.Len() {
				return nil
			}
			id = strconv.FormatUint(idx, 10)
		}
		// Select the current search configuration.
		query := plan.Search[qidx]
		qidx++
		// Build the search configuration and return the request.
		req := newReq(id, data.At(idx), newSearchConfig(t, id, query))
		return &req
	}, func(res *payload.Search_Response, err error) bool {
		// This function is called for each response received.
		if err != nil {
			handleGRPCCallError(t, err, plan)
			return true
		}
		// Extract the vector index from the request ID.
		id, _, _ := strings.Cut(res.GetRequestId(), "-")
		idx, _ := strconv.Atoi(id)
		t.Logf("request id %s searched recall: %f, payload %s", res.GetRequestId(), calculateRecall(t, neighbors.At(uint64(idx)), res), res.String())
		return true
	})
	if err != nil {
		t.Errorf("failed to complete stream: %v", err)
	}
}
