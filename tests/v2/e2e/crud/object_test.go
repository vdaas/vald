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
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/iter"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

var (
	objectRequest newRequest[*payload.Object_VectorRequest] = func(t *testing.T, _ uint64, id string, _ []float32, _ *config.Execution) *payload.Object_VectorRequest {
		return &payload.Object_VectorRequest{
			Id: existsRequest(t, 0, id, nil, nil),
		}
	}

	existsRequest newRequest[*payload.Object_ID] = func(t *testing.T, _ uint64, id string, _ []float32, _ *config.Execution) *payload.Object_ID {
		return &payload.Object_ID{
			Id: id,
		}
	}

	timestampRequest newRequest[*payload.Object_TimestampRequest] = func(t *testing.T, _ uint64, id string, _ []float32, _ *config.Execution) *payload.Object_TimestampRequest {
		return &payload.Object_TimestampRequest{
			Id: existsRequest(t, 0, id, nil, nil),
		}
	}
)

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
			unary(t, ctx, train, plan, r.client.GetObject, objectRequest)
		case config.OperationMultiple:
			t.Errorf("unsupported Object operation %s for %s", plan.Mode, plan.Type)
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamGetObject, objectRequest)
		}
	case config.OpTimestamp:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unary(t, ctx, train, plan, r.client.GetTimestamp, timestampRequest)
		case config.OperationMultiple, config.OperationStream:
			t.Errorf("unsupported Timestamp operation %s for %s", plan.Mode, plan.Type)
		}
	case config.OpExists:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			unary(t, ctx, train, plan, r.client.Exists, existsRequest)
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
					if res != nil && res.GetStatus() != nil {
						_ = handleGRPCStatusCodeError(t, err, codes.Code(res.GetStatus().GetCode()), plan)
					} else {
						handleGRPCCallError(t, err, plan)
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
