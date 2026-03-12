//go:build e2e

//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	"fmt"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/iter"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

// Predefined request builder functions for unary modify requests.
var (
	insertRequest newRequest[*payload.Insert_Request] = func(t *testing.T, idx uint64, id string, vec []float32, plan *config.Execution) *payload.Insert_Request {
		ts, skip := toModificationConfig(plan)
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
	insertWithMetadataRequest newRequest[*payload.Insert_Request] = func(t *testing.T, idx uint64, id string, vec []float32, plan *config.Execution) *payload.Insert_Request {
		ts, skip := toModificationConfig(plan)

		// Generate pre-vector metadata for testing
		// Format: idx,id,operation
		metadataValue := fmt.Sprintf("%d,%s,insert", idx, id)
		metadata := []byte(metadataValue)

		return &payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
				Metadata: metadata,
			},
			Config: &payload.Insert_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: skip,
			},
		}
	}
	insertMultipleRequest newMultiRequest[*payload.Insert_Request, *payload.Insert_MultiRequest] = func(t *testing.T, reqs ...*payload.Insert_Request) *payload.Insert_MultiRequest {
		return &payload.Insert_MultiRequest{
			Requests: reqs,
		}
	}
	updateRequest newRequest[*payload.Update_Request] = func(t *testing.T, idx uint64, id string, vec []float32, plan *config.Execution) *payload.Update_Request {
		ts, skip := toModificationConfig(plan)
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
	updateWithMetadataRequest newRequest[*payload.Update_Request] = func(t *testing.T, idx uint64, id string, vec []float32, plan *config.Execution) *payload.Update_Request {
		ts, skip := toModificationConfig(plan)

		// Generate pre-vector metadata for testing
		// Format: idx,id,operation
		metadataValue := fmt.Sprintf("%d,%s,update", idx, id)
		metadata := []byte(metadataValue)

		return &payload.Update_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
				Metadata: metadata,
			},
			Config: &payload.Update_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: skip,
			},
		}
	}
	updateMultipleRequest newMultiRequest[*payload.Update_Request, *payload.Update_MultiRequest] = func(t *testing.T, reqs ...*payload.Update_Request) *payload.Update_MultiRequest {
		return &payload.Update_MultiRequest{
			Requests: reqs,
		}
	}
	upsertRequest newRequest[*payload.Upsert_Request] = func(t *testing.T, idx uint64, id string, vec []float32, plan *config.Execution) *payload.Upsert_Request {
		ts, skip := toModificationConfig(plan)
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
	upsertWithMetadataRequest newRequest[*payload.Upsert_Request] = func(t *testing.T, idx uint64, id string, vec []float32, plan *config.Execution) *payload.Upsert_Request {
		ts, skip := toModificationConfig(plan)

		// Generate pre-vector metadata for testing
		// Format: idx,id,operation
		metadataValue := fmt.Sprintf("%d,%s,upsert", idx, id)
		metadata := []byte(metadataValue)

		return &payload.Upsert_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
				Metadata: metadata,
			},
			Config: &payload.Upsert_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: skip,
			},
		}
	}
	upsertMultipleRequest newMultiRequest[*payload.Upsert_Request, *payload.Upsert_MultiRequest] = func(t *testing.T, reqs ...*payload.Upsert_Request) *payload.Upsert_MultiRequest {
		return &payload.Upsert_MultiRequest{
			Requests: reqs,
		}
	}
	removeRequest newRequest[*payload.Remove_Request] = func(t *testing.T, idx uint64, id string, vec []float32, plan *config.Execution) *payload.Remove_Request {
		ts, skip := toModificationConfig(plan)
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
	removeMultipleRequest newMultiRequest[*payload.Remove_Request, *payload.Remove_MultiRequest] = func(t *testing.T, reqs ...*payload.Remove_Request) *payload.Remove_MultiRequest {
		return &payload.Remove_MultiRequest{
			Requests: reqs,
		}
	}
	removeByTimestampRequest newRequest[*payload.Remove_TimestampRequest] = func(t *testing.T, idx uint64, id string, vec []float32, plan *config.Execution) *payload.Remove_TimestampRequest {
		ts, _ := toModificationConfig(plan)
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

func (r *runner) processModification(
	t *testing.T,
	ctx context.Context,
	train iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
) error {
	t.Helper()
	if plan == nil {
		t.Fatal("modification plan is nil")
		return nil
	}
	switch plan.Type {
	case config.OpInsert:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			return unary(t, ctx, train, plan, r.client.Insert, insertRequest)
		case config.OperationMultiple:
			return multi(t, ctx, train, plan, r.client.MultiInsert, insertRequest, insertMultipleRequest)
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamInsert, insertRequest)
		}
	case config.OpInsertMeta:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			return unary(t, ctx, train, plan, r.client.InsertWithMetadata, insertWithMetadataRequest)
		case config.OperationMultiple:
			return multi(t, ctx, train, plan, r.client.MultiInsertWithMetadata, insertWithMetadataRequest, insertMultipleRequest)
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamInsertWithMetadata, insertWithMetadataRequest)
		}
	case config.OpUpdate:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			return unary(t, ctx, train, plan, r.client.Update, updateRequest)
		case config.OperationMultiple:
			return multi(t, ctx, train, plan, r.client.MultiUpdate, updateRequest, updateMultipleRequest)
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamUpdate, updateRequest)
		}
	case config.OpUpdateMeta:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			return unary(t, ctx, train, plan, r.client.UpdateWithMetadata, updateWithMetadataRequest)
		case config.OperationMultiple:
			return multi(t, ctx, train, plan, r.client.MultiUpdateWithMetadata, updateWithMetadataRequest, updateMultipleRequest)
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamUpdateWithMetadata, updateWithMetadataRequest)
		}
	case config.OpUpsert:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			return unary(t, ctx, train, plan, r.client.Upsert, upsertRequest)
		case config.OperationMultiple:
			return multi(t, ctx, train, plan, r.client.MultiUpsert, upsertRequest, upsertMultipleRequest)
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamUpsert, upsertRequest)
		}
	case config.OpUpsertMeta:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			return unary(t, ctx, train, plan, r.client.UpsertWithMetadata, upsertWithMetadataRequest)
		case config.OperationMultiple:
			return multi(t, ctx, train, plan, r.client.MultiUpsertWithMetadata, upsertWithMetadataRequest, upsertMultipleRequest)
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamUpsertWithMetadata, upsertWithMetadataRequest)
		}
	case config.OpRemove:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			return unary(t, ctx, train, plan, r.client.Remove, removeRequest)
		case config.OperationMultiple:
			return multi(t, ctx, train, plan, r.client.MultiRemove, removeRequest, removeMultipleRequest)
		case config.OperationStream:
			stream(t, ctx, train, plan, r.client.StreamRemove, removeRequest)
		}
	case config.OpRemoveByTimestamp:
		return single(t, ctx, 0, plan, removeByTimestampRequest(t, 0, "", nil, plan), r.client.RemoveByTimestamp)
	}
	return nil
}

func toModificationConfig(plan *config.Execution) (ts int64, skip bool) {
	if plan != nil && plan.Modification != nil {
		ts = plan.Modification.Timestamp
		skip = plan.Modification.SkipStrictExistCheck
	}
	return ts, skip
}
