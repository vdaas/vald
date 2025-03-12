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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/iter"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

type newObjectRequest[Q proto.Message] func(id string) Q

var (
	objectRequest newObjectRequest[*payload.Object_VectorRequest] = func(id string) *payload.Object_VectorRequest {
		return &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: id,
			},
		}
	}

	existsRequest newObjectRequest[*payload.Object_ID] = func(id string) *payload.Object_ID {
		return &payload.Object_ID{
			Id: id,
		}
	}

	timestampRequest newObjectRequest[*payload.Object_TimestampRequest] = func(id string) *payload.Object_TimestampRequest {
		return &payload.Object_TimestampRequest{
			Id: &payload.Object_ID{
				Id: id,
			},
		}
	}
)

// unaryModify performs unary modification operations.
// It sends a single request to the server and checks the response.
func (r *runner) processObject(
	t *testing.T,
	ctx context.Context,
	train iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
) {
	t.Helper()
	if plan == nil {
		t.Fatal("object operation plan is nil")
		return
	}
	switch plan.Type {
	case config.OpObject:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unaryObject(t, ctx, train, plan, r.client.GetObject, objectRequest)
		case config.OperationMultiple:
			t.Errorf("unsupported Object operation %s for %s", plan.Mode, plan.Type)
		case config.OperationStream:
			streamObject(t, ctx, train, plan, r.client.StreamGetObject, objectRequest)
		}
	case config.OpTimestamp:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unaryObject(t, ctx, train, plan, r.client.GetTimestamp, timestampRequest)
		case config.OperationMultiple, config.OperationStream:
			t.Errorf("unsupported Timestamp operation %s for %s", plan.Mode, plan.Type)
		}
	case config.OpExists:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unaryObject(t, ctx, train, plan, r.client.Exists, existsRequest)
		case config.OperationMultiple, config.OperationStream:
			t.Errorf("unsupported Exists operation %s for %s", plan.Mode, plan.Type)
		}
	case config.OpListObject:
		switch plan.Mode {
		case config.OperationMultiple, config.OperationStream:
			t.Errorf("unsupported ListObject operation %s for %s", plan.Mode, plan.Type)
		case config.OperationUnary, config.OperationOther:
			stream, err := r.client.StreamListObject(ctx, new(payload.Object_List_Request))
			if err != nil {
				t.Error(err)
				return
			}
			cnt := uint64(0)
			defer stream.CloseSend()
			for {
				cnt++
				res, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						return
					}
					if plan.ExpectedStatusCodes != nil && plan.ExpectedStatusCodes.Equals(codes.ToString(res.GetStatus().GetCode())) {
						t.Logf("expected error: %v", err)
					} else {
						t.Errorf("unexpected error: %v", err)
					}

					break
				}
				t.Logf("successfully get vector %v", res.GetVector())
				if cnt >= train.Len() {
					return
				}

			}
		}
	}
}

// unaryObject handles unary modification requests. It iterates over the data and for each vector,
// it sends a modification request with the specified operation and configuration.
// The function logs the result of each modification.
// The function is used for insert, update, upsert, and remove operations.
func unaryObject[Q, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
	call grpcCall[Q, R],
	newReq newObjectRequest[Q],
) {
	t.Helper()
	// Create an error group to manage concurrent requests.
	eg, ctx := errgroup.New(ctx)
	if plan != nil && plan.BaseConfig != nil {
		// Set the concurrency limit from the plan configuration.
		eg.SetLimit(int(plan.Parallelism))
	}
	for i := range data.Indexes(ctx) {
		// For each test vector, iterate over all modification configurations.
		id := strconv.FormatUint(i, 10)
		// Launch the index modify request in a goroutine.
		eg.Go(func() error {
			// Execute the modify gRPC call.
			res, err := call(ctx, newReq(id))
			if err != nil {
				log.Errorf("object request id %s returned %v", id, err)
				// Handle the error using the centralized error handler.
				handleGRPCCallError(t, err, plan)
				return nil
			}
			log.Debugf("object request id %s returned %v", id, res)
			return nil
		})
	}
	// Wait for all goroutines to complete.
	eg.Wait()
}

// streamObject handles bidirectional streaming modify requests.
// It repeatedly sends modify requests from the data slice using the provided builder,
// and processes each response received from the stream.
func streamObject[S grpc.ClientStream, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
	newStream newStream[S],
	newReq newObjectRequest[R],
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
	idx := uint64(0)
	// Use a bidirectional stream client to send requests and receive responses.
	err = grpc.BidirectionalStreamClient(stream, int(plan.Parallelism), func() *R {
		// If we have processed all vectors, return nil to close the stream.
		if idx >= data.Len() {
			return nil
		}
		// Build the modify configuration and return the request.
		req := newReq(strconv.FormatUint(idx, 10))
		return &req
	}, func(res *payload.Object_Location, err error) bool {
		// This function is called for each response received.
		if err != nil {
			handleGRPCCallError(t, err, plan)
			return true
		}
		t.Logf("vector id %s inserted to %s", res.GetUuid(), res.String())
		return true
	})
	if err != nil {
		t.Errorf("failed to complete stream: %v", err)
	}
}
