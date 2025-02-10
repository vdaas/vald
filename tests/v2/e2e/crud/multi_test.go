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
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func TestE2EMultiCRUD(t *testing.T) {
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

	eg, _ := errgroup.New(ctx)
	eg.SetLimit(int(cfg.Insert.Concurrency))
	mireq := &payload.Insert_MultiRequest{
		Requests: make([]*payload.Insert_Request, 0, cfg.Insert.BulkSize),
	}
	for i, vec := range ds.Train[cfg.Insert.Offset : cfg.Insert.Offset+cfg.Insert.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Insert.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		mireq.Requests = append(mireq.Requests, &payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
			},
			Config: &payload.Insert_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: cfg.Insert.SkipStrictExistCheck,
			},
		})
		if len(mireq.GetRequests()) >= cfg.Insert.BulkSize {
			req := mireq.CloneVT()
			mireq.Reset()
			eg.Go(safety.RecoverFunc(func() error {
				res, err := client.MultiInsert(ctx, req)
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil {
						t.Errorf("failed to insert vector: %v, status: %s", err, st.String())
					} else {
						t.Errorf("failed to insert vector: %v", err)
					}
				}
				t.Logf("vectors %s inserted %s", req.String(), res.String())
				return nil
			}))
		}
	}
	eg.Go(safety.RecoverFunc(func() error {
		res, err := client.MultiInsert(ctx, mireq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to insert vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to insert vector: %v", err)
			}
		}
		t.Logf("vectors %s inserted %s", mireq.String(), res.String())
		return nil
	}))
	eg.Wait()

	sleep(t, cfg.Index.WaitAfterInsert)

	indexStatus(t, ctx)

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Search.Concurrency))
	msreq := &payload.Search_MultiRequest{
		Requests: make([]*payload.Search_Request, 0, cfg.Search.BulkSize),
	}
	for i, vec := range ds.Test[cfg.Search.Offset : cfg.Search.Offset+cfg.Search.Num] {
		for _, query := range cfg.Search.Queries {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			var ratio *wrapperspb.FloatValue
			if query.Ratio != 0 {
				ratio = wrapperspb.Float(query.Ratio)
			} else {
				ratio = nil
			}
			msreq.Requests = append(msreq.Requests, &payload.Search_Request{
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
			})
			if len(msreq.GetRequests()) >= cfg.Search.BulkSize {
				req := msreq.CloneVT()
				msreq.Reset()
				eg.Go(safety.RecoverFunc(func() error {
					res, err := client.MultiSearch(ctx, req)
					if err != nil {
						st, ok := status.FromError(err)
						if ok && st != nil {
							t.Errorf("failed to search vector: %v, status: %s", err, st.String())
						} else {
							t.Errorf("failed to search vector: %v", err)
						}
					}
					for i, r := range res.GetResponses() {
						id, _, _ := strings.Cut(r.GetRequestId(), "-")
						idx, _ := strconv.Atoi(id)
						t.Logf("vector %v id %s searched recall: %f, payload %s", req.GetRequests()[i].GetVector(), r.GetRequestId(), calculateRecall(t, &payload.Search_Response{
							RequestId: r.GetRequestId(),
							Results:   r.GetResults(),
						}, idx), res.String())
					}
					return nil
				}))
			}
		}
	}
	eg.Go(safety.RecoverFunc(func() error {
		res, err := client.MultiSearch(ctx, msreq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to search vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to search vector: %v", err)
			}
		}
		for i, r := range res.GetResponses() {
			id, _, _ := strings.Cut(r.GetRequestId(), "-")
			idx, _ := strconv.Atoi(id)
			t.Logf("vector %v id %s searched recall: %f, payload %s", msreq.GetRequests()[i].GetVector(), r.GetRequestId(), calculateRecall(t, &payload.Search_Response{
				RequestId: r.GetRequestId(),
				Results:   r.GetResults(),
			}, idx), res.String())
		}
		return nil
	}))
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.SearchByID.Concurrency))
	msireq := &payload.Search_MultiIDRequest{
		Requests: make([]*payload.Search_IDRequest, 0, cfg.SearchByID.BulkSize),
	}
	for i := range ds.Train[cfg.SearchByID.Offset : cfg.SearchByID.Offset+cfg.SearchByID.Num] {
		for _, query := range cfg.SearchByID.Queries {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			var ratio *wrapperspb.FloatValue
			if query.Ratio != 0 {
				ratio = wrapperspb.Float(query.Ratio)
			} else {
				ratio = nil
			}
			msireq.Requests = append(msireq.Requests, &payload.Search_IDRequest{
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
			})
			if len(msireq.GetRequests()) >= cfg.SearchByID.BulkSize {
				req := msireq.CloneVT()
				msireq.Reset()
				eg.Go(safety.RecoverFunc(func() error {
					res, err := client.MultiSearchByID(ctx, req)
					if err != nil {
						st, ok := status.FromError(err)
						if ok && st != nil {
							t.Errorf("failed to search vector: %v, status: %s", err, st.String())
						} else {
							t.Errorf("failed to search vector: %v", err)
						}
					}
					for _, r := range res.GetResponses() {
						id, _, _ := strings.Cut(r.GetRequestId(), "-")
						idx, _ := strconv.Atoi(id)
						t.Logf("id %s searched recall: %f, payload %s", r.GetRequestId(), calculateRecall(t, &payload.Search_Response{
							RequestId: r.GetRequestId(),
							Results:   r.GetResults(),
						}, idx), res.String())
					}
					return nil
				}))
			}
		}
	}
	eg.Go(safety.RecoverFunc(func() error {
		res, err := client.MultiSearchByID(ctx, msireq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to search vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to search vector: %v", err)
			}
		}
		for _, r := range res.GetResponses() {
			id, _, _ := strings.Cut(r.GetRequestId(), "-")
			idx, _ := strconv.Atoi(id)
			t.Logf("id %s searched recall: %f, payload %s", r.GetRequestId(), calculateRecall(t, &payload.Search_Response{
				RequestId: r.GetRequestId(),
				Results:   r.GetResults(),
			}, idx), res.String())
		}
		return nil
	}))
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Search.Concurrency))
	msreq = &payload.Search_MultiRequest{
		Requests: make([]*payload.Search_Request, 0, cfg.Search.BulkSize),
	}
	for i, vec := range ds.Test[cfg.Search.Offset : cfg.Search.Offset+cfg.Search.Num] {
		for _, query := range cfg.Search.Queries {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			var ratio *wrapperspb.FloatValue
			if query.Ratio != 0 {
				ratio = wrapperspb.Float(query.Ratio)
			} else {
				ratio = nil
			}
			msreq.Requests = append(msreq.Requests, &payload.Search_Request{
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
			})
			if len(msreq.GetRequests()) >= cfg.Search.BulkSize {
				req := msreq.CloneVT()
				msreq.Reset()
				eg.Go(safety.RecoverFunc(func() error {
					res, err := client.MultiLinearSearch(ctx, req)
					if err != nil {
						st, ok := status.FromError(err)
						if ok && st != nil {
							t.Errorf("failed to search vector: %v, status: %s", err, st.String())
						} else {
							t.Errorf("failed to search vector: %v", err)
						}
					}
					for i, r := range res.GetResponses() {
						id, _, _ := strings.Cut(r.GetRequestId(), "-")
						idx, _ := strconv.Atoi(id)
						t.Logf("vector %v id %s searched recall: %f, payload %s", req.GetRequests()[i].GetVector(), r.GetRequestId(), calculateRecall(t, &payload.Search_Response{
							RequestId: r.GetRequestId(),
							Results:   r.GetResults(),
						}, idx), res.String())
					}
					return nil
				}))
			}
		}
	}
	eg.Go(safety.RecoverFunc(func() error {
		res, err := client.MultiLinearSearch(ctx, msreq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to search vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to search vector: %v", err)
			}
		}
		for i, r := range res.GetResponses() {
			id, _, _ := strings.Cut(r.GetRequestId(), "-")
			idx, _ := strconv.Atoi(id)
			t.Logf("vector %v id %s searched recall: %f, payload %s", msreq.GetRequests()[i].GetVector(), r.GetRequestId(), calculateRecall(t, &payload.Search_Response{
				RequestId: r.GetRequestId(),
				Results:   r.GetResults(),
			}, idx), res.String())
		}
		return nil
	}))
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.SearchByID.Concurrency))
	msireq = &payload.Search_MultiIDRequest{
		Requests: make([]*payload.Search_IDRequest, 0, cfg.SearchByID.BulkSize),
	}
	for i := range ds.Train[cfg.SearchByID.Offset : cfg.SearchByID.Offset+cfg.SearchByID.Num] {
		for _, query := range cfg.SearchByID.Queries {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			var ratio *wrapperspb.FloatValue
			if query.Ratio != 0 {
				ratio = wrapperspb.Float(query.Ratio)
			} else {
				ratio = nil
			}
			msireq.Requests = append(msireq.Requests, &payload.Search_IDRequest{
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
			})
			if len(msireq.GetRequests()) >= cfg.SearchByID.BulkSize {
				req := msireq.CloneVT()
				msireq.Reset()
				eg.Go(safety.RecoverFunc(func() error {
					res, err := client.MultiLinearSearchByID(ctx, req)
					if err != nil {
						st, ok := status.FromError(err)
						if ok && st != nil {
							t.Errorf("failed to search vector: %v, status: %s", err, st.String())
						} else {
							t.Errorf("failed to search vector: %v", err)
						}
					}
					for _, r := range res.GetResponses() {
						id, _, _ := strings.Cut(r.GetRequestId(), "-")
						idx, _ := strconv.Atoi(id)
						t.Logf("id %s searched recall: %f, payload %s", r.GetRequestId(), calculateRecall(t, &payload.Search_Response{
							RequestId: r.GetRequestId(),
							Results:   r.GetResults(),
						}, idx), res.String())
					}
					return nil
				}))
			}
		}
	}
	eg.Go(safety.RecoverFunc(func() error {
		res, err := client.MultiLinearSearchByID(ctx, msireq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to search vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to search vector: %v", err)
			}
		}
		for _, r := range res.GetResponses() {
			id, _, _ := strings.Cut(r.GetRequestId(), "-")
			idx, _ := strconv.Atoi(id)
			t.Logf("id %s searched recall: %f, payload %s", r.GetRequestId(), calculateRecall(t, &payload.Search_Response{
				RequestId: r.GetRequestId(),
				Results:   r.GetResults(),
			}, idx), res.String())
		}
		return nil
	}))
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Object.Concurrency))
	for i := range ds.Train[cfg.Object.Offset : cfg.Object.Offset+cfg.Object.Num] {
		id := strconv.Itoa(i)
		eg.Go(safety.RecoverFunc(func() error {
			obj, err := client.GetObject(ctx, &payload.Object_VectorRequest{
				Id: &payload.Object_ID{Id: id},
			})
			if err != nil {
				t.Errorf("failed to get object: %v", err)
			}
			t.Logf("id %s got object: %v", id, obj.String())

			exists, err := client.Exists(ctx, &payload.Object_ID{Id: id})
			if err != nil {
				t.Errorf("failed to check object exists: %v", err)
			}
			t.Logf("id %s exists: %v", id, exists.String())

			res, err := client.GetTimestamp(ctx, &payload.Object_TimestampRequest{
				Id: &payload.Object_ID{Id: id},
			})
			if err != nil {
				t.Errorf("failed to get timestamp: %v", err)
			}
			t.Logf("id %s got timestamp: %v", id, res.String())
			return nil
		}))
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Update.Concurrency))
	mupreq := &payload.Update_MultiRequest{
		Requests: make([]*payload.Update_Request, 0, cfg.Update.BulkSize),
	}
	for i, vec := range ds.Train[cfg.Update.Offset : cfg.Update.Offset+cfg.Update.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Update.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		mupreq.Requests = append(mupreq.Requests, &payload.Update_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
			},
			Config: &payload.Update_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: cfg.Update.SkipStrictExistCheck,
			},
		})
		if len(mupreq.GetRequests()) >= cfg.Update.BulkSize {
			req := mupreq.CloneVT()
			mupreq.Reset()
			eg.Go(safety.RecoverFunc(func() error {
				res, err := client.MultiUpdate(ctx, req)
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil {
						t.Errorf("failed to update vector: %v, status: %s", err, st.String())
					} else {
						t.Errorf("failed to update vector: %v", err)
					}
				}
				t.Logf("vectors %s updated %s", req.String(), res.String())
				return nil
			}))
		}
	}
	eg.Go(safety.RecoverFunc(func() error {
		res, err := client.MultiUpdate(ctx, mupreq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to update vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to update vector: %v", err)
			}
		}
		t.Logf("vectors %s updated %s", mupreq.String(), res.String())
		return nil
	}))
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Remove.Concurrency))
	mrreq := &payload.Remove_MultiRequest{
		Requests: make([]*payload.Remove_Request, 0, cfg.Remove.BulkSize),
	}
	for i := range ds.Train[cfg.Remove.Offset : cfg.Remove.Offset+cfg.Remove.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Remove.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		mrreq.Requests = append(mrreq.Requests, &payload.Remove_Request{
			Id: &payload.Object_ID{Id: id},
			Config: &payload.Remove_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: cfg.Remove.SkipStrictExistCheck,
			},
		})
		if len(mrreq.GetRequests()) >= cfg.Remove.BulkSize {
			req := mrreq.CloneVT()
			mrreq.Reset()
			eg.Go(safety.RecoverFunc(func() error {
				res, err := client.MultiRemove(ctx, req)
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil {
						t.Errorf("failed to remove vector: %v, status: %s", err, st.String())
					} else {
						t.Errorf("failed to remove vector: %v", err)
					}
				}
				t.Logf("vectors %s removed %s", req.String(), res.String())
				return nil
			}))
		}
	}
	eg.Go(safety.RecoverFunc(func() error {
		res, err := client.MultiRemove(ctx, mrreq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to remove vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to remove vector: %v", err)
			}
		}
		t.Logf("vectors %s removed %s", mrreq.String(), res.String())
		return nil
	}))
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Upsert.Concurrency))
	musreq := &payload.Upsert_MultiRequest{
		Requests: make([]*payload.Upsert_Request, 0, cfg.Upsert.BulkSize),
	}
	for i, vec := range ds.Train[cfg.Upsert.Offset : cfg.Upsert.Offset+cfg.Upsert.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Upsert.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		musreq.Requests = append(musreq.Requests, &payload.Upsert_Request{
			Vector: &payload.Object_Vector{
				Id:        id,
				Vector:    vec,
				Timestamp: ts,
			},
			Config: &payload.Upsert_Config{
				Timestamp:            ts,
				SkipStrictExistCheck: cfg.Upsert.SkipStrictExistCheck,
			},
		})
		if len(musreq.GetRequests()) >= cfg.Upsert.BulkSize {
			req := musreq.CloneVT()
			musreq.Reset()
			eg.Go(safety.RecoverFunc(func() error {
				res, err := client.MultiUpsert(ctx, req)
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil {
						t.Errorf("failed to upsert vector: %v, status: %s", err, st.String())
					} else {
						t.Errorf("failed to upsert vector: %v", err)
					}
				}
				t.Logf("vectors %s upsertd %s", req.String(), res.String())
				return nil
			}))
		}
	}
	eg.Go(safety.RecoverFunc(func() error {
		res, err := client.MultiUpsert(ctx, musreq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st != nil {
				t.Errorf("failed to upsert vector: %v, status: %s", err, st.String())
			} else {
				t.Errorf("failed to upsert vector: %v", err)
			}
		}
		t.Logf("vectors %s upsertd %s", musreq.String(), res.String())
		return nil
	}))
	eg.Wait()

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
