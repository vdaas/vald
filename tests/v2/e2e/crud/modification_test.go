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
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

// newModifyRequest is a function type that creates a new modify request.
type newModifyRequest[R proto.Message] func(id string, vec []float32, ts int64, skip bool) R

// Predefined request builder functions for unary modify requests.
var (
	insertRequest newModifyRequest[*payload.Insert_Request] = func(id string, vec []float32, ts int64, skip bool) *payload.Insert_Request {
		return &payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
			},
			Config: &payload.Insert_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: skip,
			},
		}
	}
	insertMultipleRequest newMultiRequest[*payload.Insert_Request, *payload.Insert_MultiRequest] = func(reqs []*payload.Insert_Request) *payload.Insert_MultiRequest {
		return &payload.Insert_MultiRequest{
			Requests: reqs,
		}
	}
	updateRequest newModifyRequest[*payload.Update_Request] = func(id string, vec []float32, ts int64, skip bool) *payload.Update_Request {
		return &payload.Update_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
			},
			Config: &payload.Update_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: skip,
			},
		}
	}
	updateMultipleRequest newMultiRequest[*payload.Update_Request, *payload.Update_MultiRequest] = func(reqs []*payload.Update_Request) *payload.Update_MultiRequest {
		return &payload.Update_MultiRequest{
			Requests: reqs,
		}
	}
	upsertRequest newModifyRequest[*payload.Upsert_Request] = func(id string, vec []float32, ts int64, skip bool) *payload.Upsert_Request {
		return &payload.Upsert_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
			},
			Config: &payload.Upsert_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: skip,
			},
		}
	}
	upsertMultipleRequest newMultiRequest[*payload.Upsert_Request, *payload.Upsert_MultiRequest] = func(reqs []*payload.Upsert_Request) *payload.Upsert_MultiRequest {
		return &payload.Upsert_MultiRequest{
			Requests: reqs,
		}
	}
	removeRequest newModifyRequest[*payload.Remove_Request] = func(id string, vec []float32, ts int64, skip bool) *payload.Remove_Request {
		return &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: id,
			},
			Config: &payload.Remove_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: skip,
			},
		}
	}
	removeMultipleRequest newMultiRequest[*payload.Remove_Request, *payload.Remove_MultiRequest] = func(reqs []*payload.Remove_Request) *payload.Remove_MultiRequest {
		return &payload.Remove_MultiRequest{
			Requests: reqs,
		}
	}
	removeByTimestampRequest newModifyRequest[*payload.Remove_TimestampRequest] = func(id string, vec []float32, ts int64, skip bool) *payload.Remove_TimestampRequest {
		if ts == 0 {
			ts = time.Now().UnixNano()
		}
		return &payload.Remove_TimestampRequest{
			Timestamps: []*payload.Remove_Timestamp{
				{
					Timestamp: ts,
					Operator:  payload.Remove_Timestamp_Le,
				},
			},
		}
	}
)

// unaryModify performs unary modification operations.
// It sends a single request to the server and checks the response.
func (r *runner) processModification(
	t *testing.T,
	ctx context.Context,
	train iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
) {
	t.Helper()
	if plan == nil {
		t.Fatal("modification plan is nil")
		return
	}
	switch plan.Type {
	case config.OpInsert:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unaryModify(t, ctx, train, plan, r.client.Insert, insertRequest)
		case config.OperationMultiple:
			multiModify(t, ctx, train, plan, r.client.MultiInsert, insertRequest, insertMultipleRequest)
		case config.OperationStream:
			streamModify(t, ctx, train, plan, r.client.StreamInsert, insertRequest)
		}
	case config.OpUpdate:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unaryModify(t, ctx, train, plan, r.client.Update, updateRequest)
		case config.OperationMultiple:
			multiModify(t, ctx, train, plan, r.client.MultiUpdate, updateRequest, updateMultipleRequest)
		case config.OperationStream:
			streamModify(t, ctx, train, plan, r.client.StreamUpdate, updateRequest)
		}
	case config.OpUpsert:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unaryModify(t, ctx, train, plan, r.client.Upsert, upsertRequest)
		case config.OperationMultiple:
			multiModify(t, ctx, train, plan, r.client.MultiUpsert, upsertRequest, upsertMultipleRequest)
		case config.OperationStream:
			streamModify(t, ctx, train, plan, r.client.StreamUpsert, upsertRequest)
		}
	case config.OpRemove:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unaryModify(t, ctx, train, plan, r.client.Remove, removeRequest)
		case config.OperationMultiple:
			multiModify(t, ctx, train, plan, r.client.MultiRemove, removeRequest, removeMultipleRequest)
		case config.OperationStream:
			streamModify(t, ctx, train, plan, r.client.StreamRemove, removeRequest)
		}
	case config.OpRemoveByTimestamp:
		unaryModify(t, ctx, train, plan, r.client.RemoveByTimestamp, removeByTimestampRequest)
	}
}

// unaryModify handles unary modification requests. It iterates over the data and for each vector,
// it sends a modification request with the specified operation and configuration.
// The function logs the result of each modification.
// The function is used for insert, update, upsert, and remove operations.
func unaryModify[Q, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
	call grpcCall[Q, R],
	newReq newModifyRequest[Q],
) {
	t.Helper()
	// Create an error group to manage concurrent requests.
	eg, ctx := errgroup.New(ctx)
	// Set the concurrency limit from the plan configuration.
	if plan != nil && plan.BaseConfig != nil {
		// Set the concurrency limit from the plan configuration.
		eg.SetLimit(int(plan.Parallelism))
	}
	var (
		ts   int64
		skip bool
	)
	if plan.Modification != nil {
		ts = plan.Modification.Timestamp
		skip = plan.Modification.SkipStrictExistCheck
	}
	for i, vec := range data.Seq2(ctx) {
		// For each test vector, iterate over all modification configurations.
		id := strconv.FormatUint(i, 10)
		// Launch the index modify request in a goroutine.
		eg.Go(func() error {
			req := newReq(id, vec, ts, skip)
			// Execute the modify gRPC call.
			_, err := call(ctx, req)
			if err != nil {
				// Handle the error using the centralized error handler.
				handleGRPCCallError(t, err, plan)
				return nil
			}
			return nil
		})
	}
	// Wait for all goroutines to complete.
	eg.Wait()
}

// multiModify handles bulk modify requests by grouping individual requests up to BulkSize.
// Once the bulk size is reached, it sends the grouped requests and logs the responses.
// It uses the provided builder functions to create the individual requests and the bulk request.
// The bulk request is sent using the provided gRPC call function.
// The function logs the response for each batch of requests.
// The function is used for insert, update, upsert, and remove operations.
func multiModify[Q, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
	call grpcCall[R, *payload.Object_Locations],
	addReqs newModifyRequest[Q],
	toReq newMultiRequest[Q, R],
) {
	t.Helper()
	eg, ctx := errgroup.New(ctx)
	eg.SetLimit(int(plan.Parallelism))
	exec := func(batch []Q) {
		// Convert the slice of individual requests into a bulk request.
		_, err := call(ctx, toReq(batch))
		if err != nil {
			handleGRPCCallError(t, err, plan)
			return
		}
	}

	// Initialize a slice to hold the bulk requests.
	reqs := make([]Q, 0, plan.BulkSize)
	var (
		ts   int64
		skip bool
	)
	if plan.Modification != nil {
		ts = plan.Modification.Timestamp
		skip = plan.Modification.SkipStrictExistCheck
	}
	for i, vec := range data.Seq2(ctx) {
		id := strconv.FormatUint(i, 10)
		// Append a new request to the bulk slice using the provided builder.
		reqs = append(reqs, addReqs(id, vec, ts, skip))
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
	exec(reqs)
	eg.Wait()
}

// streamModify handles bidirectional streaming modify requests.
// It repeatedly sends modify requests from the data slice using the provided builder,
// and processes each response received from the stream.
func streamModify[S grpc.ClientStream, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
	newStream newStream[S],
	newReq newModifyRequest[R],
) {
	t.Helper()
	// Create a new stream using the provided stream function.
	stream, err := newStream(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	var (
		ts   int64
		skip bool
	)
	if plan.Modification != nil {
		ts = plan.Modification.Timestamp
		skip = plan.Modification.SkipStrictExistCheck
	}
	// qidx tracks the current index within the modify configuration slice.
	// idx tracks the current vector index.
	idx := uint64(0)
	// Use a bidirectional stream client to send requests and receive responses.
	err = grpc.BidirectionalStreamClient(stream, int(plan.Parallelism), func() *R {
		// If we have processed all vectors, return nil to close the stream.
		if idx >= data.Len() {
			return nil
		}
		id := strconv.FormatUint(idx, 10)
		// Build the modify configuration and return the request.
		req := newReq(id, data.At(idx), ts, skip)
		return &req
	}, func(res *payload.Object_Location, err error) bool {
		// This function is called for each response received.
		if err != nil {
			handleGRPCCallError(t, err, plan)
			return true
		}
		return true
	})
	if err != nil {
		t.Errorf("failed to complete stream: %v", err)
	}
}
