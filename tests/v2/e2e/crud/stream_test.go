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

// package crud provides e2e tests using ann-benchmarks datasets
package crud

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func TestE2EStreamCRUD(t *testing.T) {
	timestamp := time.Now().UnixNano()

	{
		res, err := client.IndexProperty(ctx, &payload.Empty{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to get IndexProperty %v status: %s", err, st.String())
			} else {
				t.Errorf("failed to get IndexProperty %v", err)
			}
		}
		t.Logf("IndexProperty: %v", res.String())
	}

	var (
		stream grpc.ClientStream
		err    error
	)
	stream, err = client.StreamInsert(ctx)
	if err != nil {
		t.Error(err)
	}
	var idx int
	ts := cfg.Insert.Timestamp
	if ts == 0 {
		ts = timestamp
	}
	datas := ds.Train[cfg.Insert.Offset : cfg.Insert.Offset+cfg.Insert.Num]
	err = grpc.BidirectionalStreamClient(stream, func() *payload.Insert_Request {
		id := strconv.Itoa(idx)
		if len(datas) < idx {
			return nil
		}
		vec := datas[idx]
		idx++
		return &payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
			},
			Config: &payload.Insert_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: cfg.Insert.SkipStrictExistCheck,
			},
		}
	}, func(res *payload.Object_Location, err error) bool {
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to insert vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to insert vector: %v", err)
			}
		}
		t.Logf("vector id %s inserted to %s", res.GetUuid(), res.String())
		return true
	})
	if err != nil {
		t.Errorf("failed to complete insert stream %v", err)
	}

	sleep(t, cfg.Index.WaitAfterInsert)

	indexStatus(t, ctx)

	stream, err = client.StreamSearch(ctx)
	if err != nil {
		t.Error(err)
	}
	qidx := 0
	idx = 0
	datas = ds.Test[cfg.Search.Offset : cfg.Search.Offset+cfg.Search.Num]
	err = grpc.BidirectionalStreamClient(stream, func() *payload.Search_Request {
		id := strconv.Itoa(idx)
		if len(datas) < idx {
			return nil
		}
		if len(cfg.Search.Queries) < qidx {
			qidx = 0
			idx++
		}
		query := cfg.Search.Queries[qidx]
		rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
		vec := datas[idx]
		qidx++
		var ratio *wrapperspb.FloatValue
		if query.Ratio != 0 {
			ratio = wrapperspb.Float(query.Ratio)
		} else {
			ratio = nil
		}
		return &payload.Search_Request{
			Vector: vec,
			Config: &payload.Search_Config{
				RequestId:            rid,
				Num:                  query.K,
				Radius:               query.Radius,
				Epsilon:              query.Epsilon,
				Timeout:              query.Timeout.Nanoseconds(),
				AggregationAlgorithm: query.Algorithm,
				MinNum:               query.MinNum,
				Ratio:                ratio,
				Nprobe:               query.Nprobe,
			},
		}
	}, func(res *payload.Search_Response, err error) bool {
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to search vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to search vector: %v", err)
			}
		}
		id, _, _ := strings.Cut(res.GetRequestId(), "-")
		idx, _ := strconv.Atoi(id)
		t.Logf("request id %s searched recall: %f, payload %s", res.GetRequestId(), calculateRecall(t, &payload.Search_Response{
			RequestId: res.GetRequestId(),
			Results:   res.GetResults(),
		}, idx), res.String())

		return true
	})
	if err != nil {
		t.Errorf("failed to complete insert stream %v", err)
	}

	stream, err = client.StreamSearchByID(ctx)
	if err != nil {
		t.Error(err)
	}
	qidx = 0
	idx = 0
	datas = ds.Train[cfg.SearchByID.Offset : cfg.SearchByID.Offset+cfg.SearchByID.Num]
	err = grpc.BidirectionalStreamClient(stream, func() *payload.Search_IDRequest {
		id := strconv.Itoa(idx)
		if len(datas) < idx {
			return nil
		}
		if len(cfg.Search.Queries) < qidx {
			qidx = 0
			idx++
		}
		query := cfg.Search.Queries[qidx]
		rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
		qidx++
		var ratio *wrapperspb.FloatValue
		if query.Ratio != 0 {
			ratio = wrapperspb.Float(query.Ratio)
		} else {
			ratio = nil
		}
		return &payload.Search_IDRequest{
			Id: id,
			Config: &payload.Search_Config{
				RequestId:            rid,
				Num:                  query.K,
				Radius:               query.Radius,
				Epsilon:              query.Epsilon,
				Timeout:              query.Timeout.Nanoseconds(),
				AggregationAlgorithm: query.Algorithm,
				MinNum:               query.MinNum,
				Ratio:                ratio,
				Nprobe:               query.Nprobe,
			},
		}
	}, func(res *payload.Search_Response, err error) bool {
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to search vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to search vector: %v", err)
			}
		}
		id, _, _ := strings.Cut(res.GetRequestId(), "-")
		idx, _ := strconv.Atoi(id)
		t.Logf("request id %s searched recall: %f, payload %s", res.GetRequestId(), calculateRecall(t, &payload.Search_Response{
			RequestId: res.GetRequestId(),
			Results:   res.GetResults(),
		}, idx), res.String())

		return true
	})
	if err != nil {
		t.Errorf("failed to complete insert stream %v", err)
	}
	stream, err = client.StreamLinearSearch(ctx)
	if err != nil {
		t.Error(err)
	}
	qidx = 0
	idx = 0
	datas = ds.Test[cfg.LinearSearch.Offset : cfg.LinearSearch.Offset+cfg.LinearSearch.Num]
	err = grpc.BidirectionalStreamClient(stream, func() *payload.Search_Request {
		id := strconv.Itoa(idx)
		if len(datas) < idx {
			return nil
		}
		if len(cfg.LinearSearch.Queries) < qidx {
			qidx = 0
			idx++
		}
		query := cfg.LinearSearch.Queries[qidx]
		rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
		vec := datas[idx]
		qidx++
		var ratio *wrapperspb.FloatValue
		if query.Ratio != 0 {
			ratio = wrapperspb.Float(query.Ratio)
		} else {
			ratio = nil
		}
		return &payload.Search_Request{
			Vector: vec,
			Config: &payload.Search_Config{
				RequestId:            rid,
				Num:                  query.K,
				Radius:               query.Radius,
				Epsilon:              query.Epsilon,
				Timeout:              query.Timeout.Nanoseconds(),
				AggregationAlgorithm: query.Algorithm,
				MinNum:               query.MinNum,
				Ratio:                ratio,
				Nprobe:               query.Nprobe,
			},
		}
	}, func(res *payload.Search_Response, err error) bool {
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to search vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to search vector: %v", err)
			}
		}
		id, _, _ := strings.Cut(res.GetRequestId(), "-")
		idx, _ := strconv.Atoi(id)
		t.Logf("request id %s searched recall: %f, payload %s", res.GetRequestId(), calculateRecall(t, &payload.Search_Response{
			RequestId: res.GetRequestId(),
			Results:   res.GetResults(),
		}, idx), res.String())

		return true
	})
	if err != nil {
		t.Errorf("failed to complete insert stream %v", err)
	}

	stream, err = client.StreamLinearSearchByID(ctx)
	if err != nil {
		t.Error(err)
	}
	qidx = 0
	idx = 0
	datas = ds.Train[cfg.LinearSearchByID.Offset : cfg.LinearSearchByID.Offset+cfg.LinearSearchByID.Num]
	err = grpc.BidirectionalStreamClient(stream, func() *payload.Search_IDRequest {
		id := strconv.Itoa(idx)
		if len(datas) < idx {
			return nil
		}
		if len(cfg.LinearSearch.Queries) < qidx {
			qidx = 0
			idx++
		}
		query := cfg.LinearSearch.Queries[qidx]
		rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
		qidx++
		var ratio *wrapperspb.FloatValue
		if query.Ratio != 0 {
			ratio = wrapperspb.Float(query.Ratio)
		} else {
			ratio = nil
		}
		return &payload.Search_IDRequest{
			Id: id,
			Config: &payload.Search_Config{
				RequestId:            rid,
				Num:                  query.K,
				Radius:               query.Radius,
				Epsilon:              query.Epsilon,
				Timeout:              query.Timeout.Nanoseconds(),
				AggregationAlgorithm: query.Algorithm,
				MinNum:               query.MinNum,
				Ratio:                ratio,
				Nprobe:               query.Nprobe,
			},
		}
	}, func(res *payload.Search_Response, err error) bool {
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to search vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to search vector: %v", err)
			}
		}
		id, _, _ := strings.Cut(res.GetRequestId(), "-")
		idx, _ := strconv.Atoi(id)
		t.Logf("request id %s searched recall: %f, payload %s", res.GetRequestId(), calculateRecall(t, &payload.Search_Response{
			RequestId: res.GetRequestId(),
			Results:   res.GetResults(),
		}, idx), res.String())

		return true
	})
	if err != nil {
		t.Errorf("failed to complete insert stream %v", err)
	}

	stream, err = client.StreamGetObject(ctx)
	if err != nil {
		t.Error(err)
	}
	idx = 0
	datas = ds.Train[cfg.Object.Offset : cfg.Object.Offset+cfg.Object.Num]
	err = grpc.BidirectionalStreamClient(stream, func() *payload.Object_VectorRequest {
		id := strconv.Itoa(idx)
		if len(datas) < idx {
			return nil
		}
		idx++
		return &payload.Object_VectorRequest{
			Id: &payload.Object_ID{Id: id},
		}
	}, func(res *payload.Object_Vector, err error) bool {
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to get vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to get vector: %v", err)
			}
		}
		t.Logf("vector id %s loaded %s", res.GetId(), res.String())
		return true
	})
	if err != nil {
		t.Errorf("failed to complete get object stream %v", err)
	}

	stream, err = client.StreamUpdate(ctx)
	if err != nil {
		t.Error(err)
	}
	ts = cfg.Update.Timestamp
	if ts == 0 {
		ts = timestamp
	}
	idx = 0
	datas = ds.Train[cfg.Update.Offset : cfg.Update.Offset+cfg.Update.Num]
	err = grpc.BidirectionalStreamClient(stream, func() *payload.Update_Request {
		id := strconv.Itoa(idx)
		if len(datas) < idx {
			return nil
		}
		vec := datas[idx]
		idx++
		return &payload.Update_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
			},
			Config: &payload.Update_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: cfg.Update.SkipStrictExistCheck,
			},
		}
	}, func(res *payload.Object_Location, err error) bool {
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to update vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to update vector: %v", err)
			}
		}
		t.Logf("vector id %s updated to %s", res.GetUuid(), res.String())
		return true
	})
	if err != nil {
		t.Errorf("failed to complete update stream %v", err)
	}

	stream, err = client.StreamRemove(ctx)
	if err != nil {
		t.Error(err)
	}
	ts = cfg.Remove.Timestamp
	if ts == 0 {
		ts = timestamp
	}
	idx = 0
	datas = ds.Train[cfg.Remove.Offset : cfg.Update.Offset+cfg.Update.Num]
	err = grpc.BidirectionalStreamClient(stream, func() *payload.Remove_Request {
		id := strconv.Itoa(idx)
		if len(datas) < idx {
			return nil
		}
		idx++
		return &payload.Remove_Request{
			Id: &payload.Object_ID{Id: id},
			Config: &payload.Remove_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: cfg.Remove.SkipStrictExistCheck,
			},
		}
	}, func(res *payload.Object_Location, err error) bool {
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to remove vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to remove vector: %v", err)
			}
		}
		t.Logf("vector id %s removed to %s", res.GetUuid(), res.String())
		return true
	})
	if err != nil {
		t.Errorf("failed to complete remove stream %v", err)
	}

	stream, err = client.StreamUpsert(ctx)
	if err != nil {
		t.Error(err)
	}
	ts = cfg.Upsert.Timestamp
	if ts == 0 {
		ts = timestamp
	}
	idx = 0
	datas = ds.Train[cfg.Upsert.Offset : cfg.Upsert.Offset+cfg.Upsert.Num]
	err = grpc.BidirectionalStreamClient(stream, func() *payload.Upsert_Request {
		id := strconv.Itoa(idx)
		if len(datas) < idx {
			return nil
		}
		vec := datas[idx]
		idx++
		return &payload.Upsert_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
			},
			Config: &payload.Upsert_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: cfg.Upsert.SkipStrictExistCheck,
			},
		}
	}, func(res *payload.Object_Location, err error) bool {
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to upsert vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to upsert vector: %v", err)
			}
		}
		t.Logf("vector id %s upserted to %s", res.GetUuid(), res.String())
		return true
	})
	if err != nil {
		t.Errorf("failed to complete upsert stream %v", err)
	}

	{
		rts := time.Now().Add(-time.Hour).UnixNano()
		res, err := client.RemoveByTimestamp(ctx, &payload.Remove_TimestampRequest{
			Timestamps: []*payload.Remove_Timestamp{
				{
					Timestamp: rts,
					Operator:  payload.Remove_Timestamp_Le,
				},
			},
		})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to remove by timestamp vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to remove by timestamp vector: %v", err)
			}
		}
		t.Logf("removed by timestamp %s to %s", time.Unix(0, rts).String(), res.String())
	}

	{
		res, err := client.Flush(ctx, &payload.Flush_Request{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to flush %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to flush %v", err)
			}
		}
		t.Logf("flushed %s", res.String())
	}

	indexStatus(t, ctx)
}
