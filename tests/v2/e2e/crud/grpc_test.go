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

// Package crud provides end-to-end tests using ann-benchmarks datasets.
package crud

import (
	"context"
	"slices"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/internal/iter"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

// Type aliases for generic search functions.
type (
	// grpcCall is a generic function type for making gRPC calls.
	grpcCall[Q, R proto.Message] func(ctx context.Context, query Q, opts ...grpc.CallOption) (response R, err error)
	// newStream is a generic type for functions that create a new gRPC stream.
	newStream[S grpc.ClientStream] func(ctx context.Context, opts ...grpc.CallOption) (S, error)
	// newRequest is a function type that creates a new request.
	newRequest[Q proto.Message] func(t *testing.T, idx uint64, id string, vec []float32, e *config.Execution) Q
	// newMultiRequest is a generic type for functions that build bulk search requests.
	newMultiRequest[R, S proto.Message] func(t *testing.T, reqs []R) S
	// callback is a function type that processes the response and error from a gRPC call.
	callback[R proto.Message] func(t *testing.T, idx uint64, res R, err error) bool
)

// handleGRPCCallError centralizes the gRPC error handling and logging.
// It compares the error's status code with the expected codes from the plan.
// If the error is expected, it logs a message; otherwise, it logs an error.
func handleGRPCCallError(t *testing.T, err error, plan *config.Execution) {
	t.Helper()
	if err != nil {
		if st, ok := status.FromError(err); ok && st != nil {
			if len(plan.ExpectedStatusCodes) != 0 && !plan.ExpectedStatusCodes.Equals(st.Code().String()) {
				t.Errorf("unexpected error: %v", st)
			}
			return
		}
		t.Errorf("failed to execute gRPC call error: %v", err)
	}
}

func single[Q, R proto.Message](
	t *testing.T,
	ctx context.Context,
	idx uint64,
	plan *config.Execution,
	req Q,
	call grpcCall[Q, R],
	callback ...callback[R],
) {
	t.Helper()
	if plan.BaseConfig != nil && plan.BaseConfig.Limiter != nil {
		plan.BaseConfig.Limiter.Wait(ctx)
	}
	// Execute the modify gRPC call.
	res, err := call(ctx, req)
	if err != nil {
		// Handle the error using the centralized error handler.
		handleGRPCCallError(t, err, plan)
		return
	}

	for _, cb := range callback {
		if cb != nil {
			if !cb(t, idx, res, err) {
				return
			}
		}
	}
	return
}

func unary[Q, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
	call grpcCall[Q, R],
	newReq newRequest[Q],
	callback ...callback[R],
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
		// Copy id to avoid data race.
		idx := i
		// Execute request in a goroutine.
		eg.Go(func() error {
			single(t, ctx, idx, plan, newReq(t, idx, strconv.FormatUint(idx, 10), vec, plan), call, callback...)
			return nil
		})
	}
	// Wait for all goroutines to complete.
	eg.Wait()
}

func multi[Q, M, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
	call grpcCall[M, R],
	addReqs newRequest[Q],
	toReq newMultiRequest[Q, M],
	callbacks ...callback[R],
) {
	t.Helper()
	eg, ctx := errgroup.New(ctx)
	// Set the concurrency limit from the plan configuration.
	if plan != nil && plan.BaseConfig != nil {
		// Set the concurrency limit from the plan configuration.
		eg.SetLimit(int(plan.Parallelism))
	}
	// Initialize a slice to hold the bulk requests.
	reqs := make([]Q, 0, plan.BulkSize)
	for i, vec := range data.Seq2(ctx) {
		id := strconv.FormatUint(i, 10)
		// Append a new request to the bulk slice using the provided builder.
		reqs = append(reqs, addReqs(t, i, id, vec, plan))
		// If the bulk size is reached, send the batch.
		if len(reqs) >= int(plan.BulkSize) {
			// Capture the current batch.
			batch := slices.Clone(reqs)
			idx := i
			// Meset the bulk request slice for the next batch.
			reqs = reqs[:0]
			eg.Go(func() error {
				single(t, ctx, idx, plan, toReq(t, batch), call, callbacks...)
				return nil
			})
		}
	}
	eg.Go(func() error {
		single(t, ctx, data.Len(), plan, toReq(t, reqs), call, callbacks...)
		return nil
	})
	eg.Wait()
}

func stream[S grpc.ClientStream, Q, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
	newStream newStream[S],
	newReq newRequest[Q],
	callbacks ...callback[R],
) {
	t.Helper()
	// Create a new stream using the provided stream function.
	stream, err := newStream(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	// qidx tracks the current index within the modify configuration slice.
	// idx tracks the current vector index.
	var idx atomic.Uint64
	var sidx atomic.Uint64
	// Use a bidirectional stream client to send requests and receive responses.
	err = grpc.BidirectionalStreamClient(stream, int(plan.Parallelism), func() *Q {
		// If we have processed all vectors, return nil to close the stream.
		if idx.Load() >= data.Len() {
			return nil
		}
		// Build the modify configuration and return the request.
		req := newReq(t, idx.Load(), strconv.FormatUint(idx.Load(), 10), data.At(idx.Load()), plan)
		idx.Add(1)
		return &req
	}, func(res *R, err error) bool {
		id := sidx.Add(1) - 1
		for _, cb := range callbacks {
			if cb != nil {
				if !cb(t, id, *res, err) {
					return false
				}
			}
		}
		return true
	})
	if err != nil {
		t.Errorf("failed to complete stream: %v", err)
	}
}
