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
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func TestE2EUnaryCRUD(t *testing.T) {
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

	var eg errgroup.Group
	if cfg.Insert.Num != 0 {
		eg, _ := errgroup.New(ctx)
		eg.SetLimit(int(cfg.Insert.Concurrency))
		for i, vec := range ds.Train[cfg.Insert.Offset : cfg.Insert.Offset+cfg.Insert.Num] {
			id := strconv.Itoa(i)
			ts := cfg.Insert.Timestamp
			if ts == 0 {
				ts = timestamp
			}
			eg.Go(safety.RecoverFunc(func() error {
				res, err := client.Insert(ctx, &payload.Insert_Request{
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
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil {
						t.Errorf("failed to insert vector: %v, status: %s", err, st.String())
					} else {
						t.Errorf("failed to insert vector: %v", err)
					}
				}
				t.Logf("vector %v id %s inserted to %s", vec, id, res.String())
				return nil
			}))
		}
		eg.Wait()

		sleep(t, cfg.Index.WaitAfterInsert)

	}

	indexStatus(t, ctx)

	if cfg.Search.Num != 0 {
		eg, _ = errgroup.New(ctx)
		eg.SetLimit(int(cfg.Search.Concurrency))
		for i, vec := range ds.Test[cfg.Search.Offset : cfg.Search.Offset+cfg.Search.Num] {
			for _, query := range cfg.Search.Queries {
				id := strconv.Itoa(i)
				rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
				eg.Go(safety.RecoverFunc(func() error {
					var ratio *wrapperspb.FloatValue
					if query.Ratio != 0 {
						ratio = wrapperspb.Float(query.Ratio)
					} else {
						ratio = nil
					}
					res, err := client.Search(ctx, &payload.Search_Request{
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
					if err != nil {
						st, ok := status.FromError(err)
						if ok && st != nil {
							t.Errorf("failed to search vector: %v, status: %s", err, st.String())
						} else {
							t.Errorf("failed to search vector: %v", err)
						}
					}
					t.Logf("vector %v id %s searched recall: %f, payload %s", vec, rid, calculateRecall(t, res, i), res.String())
					return nil
				}))
			}
		}
		eg.Wait()
	}

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.SearchByID.Concurrency))
	for i, vec := range ds.Train[cfg.SearchByID.Offset : cfg.SearchByID.Offset+cfg.SearchByID.Num] {
		for _, query := range cfg.SearchByID.Queries {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			eg.Go(safety.RecoverFunc(func() error {
				var ratio *wrapperspb.FloatValue
				if query.Ratio != 0 {
					ratio = wrapperspb.Float(query.Ratio)
				} else {
					ratio = nil
				}

				res, err := client.SearchByID(ctx, &payload.Search_IDRequest{
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
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil {
						t.Errorf("failed to search vector: %v, status: %s", err, st.String())
					} else {
						t.Errorf("failed to search vector: %v", err)
					}
				}
				t.Logf("vector %v id %s searched recall: %f, payload %s", vec, rid, calculateRecall(t, res, i), res.String())
				return nil
			}))
		}
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.LinearSearch.Concurrency))
	for i, vec := range ds.Test[cfg.LinearSearch.Offset : cfg.LinearSearch.Offset+cfg.LinearSearch.Num] {
		for _, query := range cfg.LinearSearch.Queries {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			eg.Go(safety.RecoverFunc(func() error {
				var ratio *wrapperspb.FloatValue
				if query.Ratio != 0 {
					ratio = wrapperspb.Float(query.Ratio)
				} else {
					ratio = nil
				}

				res, err := client.LinearSearch(ctx, &payload.Search_Request{
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
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil {
						t.Errorf("failed to search vector: %v, status: %s", err, st.String())
					} else {
						t.Errorf("failed to search vector: %v", err)
					}
				}
				t.Logf("vector %v id %s searched recall: %f, payload %s", vec, rid, calculateRecall(t, res, i), res.String())
				return nil
			}))
		}
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.LinearSearchByID.Concurrency))
	for i, vec := range ds.Train[cfg.LinearSearchByID.Offset : cfg.LinearSearchByID.Offset+cfg.LinearSearchByID.Num] {
		for _, query := range cfg.LinearSearchByID.Queries {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			eg.Go(safety.RecoverFunc(func() error {
				var ratio *wrapperspb.FloatValue
				if query.Ratio != 0 {
					ratio = wrapperspb.Float(query.Ratio)
				} else {
					ratio = nil
				}

				res, err := client.LinearSearchByID(ctx, &payload.Search_IDRequest{
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
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil {
						t.Errorf("failed to search vector: %v, status: %s", err, st.String())
					} else {
						t.Errorf("failed to search vector: %v", err)
					}
				}
				t.Logf("vector %v id %s searched recall: %f, payload %s", vec, rid, calculateRecall(t, res, i), res.String())
				return nil
			}))
		}
	}
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
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to get object: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to get object: %v", err)
				}
			}
			t.Logf("id %s got object: %v", id, obj.String())

			exists, err := client.Exists(ctx, &payload.Object_ID{Id: id})
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to check object exists: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to check object exitst: %v", err)
				}
			}
			t.Logf("id %s exists: %v", id, exists.String())

			res, err := client.GetTimestamp(ctx, &payload.Object_TimestampRequest{
				Id: &payload.Object_ID{Id: id},
			})
			if err != nil {
				t.Errorf("failed to get timestamp: %v", err)
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to get object timestamp: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to get object timestamp: %v", err)
				}
			}
			t.Logf("id %s got timestamp: %v", id, res.String())
			return nil
		}))
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Update.Concurrency))
	for i, vec := range ds.Train[cfg.Update.Offset : cfg.Update.Offset+cfg.Update.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Update.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Update(ctx, &payload.Update_Request{
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
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to update vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to update vector: %v", err)
				}
			}
			t.Logf("vector %v id %s updated to %s", vec, id, res.String())
			return nil
		}))
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Remove.Concurrency))
	for i := range ds.Train[cfg.Remove.Offset : cfg.Remove.Offset+cfg.Remove.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Remove.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Remove(ctx, &payload.Remove_Request{
				Id: &payload.Object_ID{Id: id},
				Config: &payload.Remove_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: cfg.Remove.SkipStrictExistCheck,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to remove vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to remove vector: %v", err)
				}
			}
			t.Logf("id %s'd vector removed to %s", id, res.String())
			return nil
		}))
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Upsert.Concurrency))
	for i, vec := range ds.Train[cfg.Upsert.Offset : cfg.Upsert.Offset+cfg.Upsert.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Upsert.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Upsert(ctx, &payload.Upsert_Request{
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
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to upsert vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to upsert vector: %v", err)
				}
			}
			t.Logf("vector %v id %s upserted to %s", vec, id, res.String())
			return nil
		}))
	}
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

func TestE2EUnarySkipStrictExistsCheckCRUD(t *testing.T) {
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

	t.Log("starting test #1 run Update with SkipStrictExistCheck=true and check that it fails.")
	eg, _ := errgroup.New(ctx)
	eg.SetLimit(int(cfg.Update.Concurrency))
	for i, vec := range ds.Train[cfg.Update.Offset : cfg.Update.Offset+cfg.Update.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Update.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			_, err := client.Update(ctx, &payload.Update_Request{
				Vector: &payload.Object_Vector{
					Id:        id,
					Vector:    vec,
					Timestamp: ts,
				},
				Config: &payload.Update_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: true,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if !ok || st == nil || st.Code() != codes.NotFound {
					t.Errorf("update vector response is not NotFound: %v with SkipStrictExistCheck=true", err)
				}
			}
			t.Logf("received a NotFound error on #1: %s", err.Error())
			return nil
		}))
	}
	eg.Wait()

	t.Log("starting test #2 run Update with SkipStrictExistCheck=false, and check that the internal Remove Operation returns a NotFound error.")
	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Update.Concurrency))
	for i, vec := range ds.Train[cfg.Update.Offset : cfg.Update.Offset+cfg.Update.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Update.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			_, err := client.Update(ctx, &payload.Update_Request{
				Vector: &payload.Object_Vector{
					Id:        id,
					Vector:    vec,
					Timestamp: ts,
				},
				Config: &payload.Update_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: false,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if !ok || st == nil || st.Code() != codes.NotFound {
					t.Errorf("update vector response is not NotFound: %v with SkipStrictExistCheck=false", err)
				}
			}
			t.Logf("received a NotFound error on #2: %s", err.Error())
			return nil
		}))
	}
	eg.Wait()

	t.Log("starting test #3 run Insert with SkipStrictExistCheck=false and confirmed that it succeeded")
	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Insert.Concurrency))
	for i, vec := range ds.Train[cfg.Insert.Offset : cfg.Insert.Offset+cfg.Insert.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Insert.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Insert(ctx, &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:        id,
					Vector:    vec,
					Timestamp: ts,
				},
				Config: &payload.Insert_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: false,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to insert vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to insert vector: %v", err)
				}
			}
			t.Logf("vector %v id %s inserted on #3 to %s", vec, id, res.String())
			return nil
		}))
	}
	eg.Wait()

	sleep(t, cfg.Index.WaitAfterInsert)

	indexStatus(t, ctx)

	t.Log("starting test #4 run Update with SkipStrictExistCheck=false & a different vector, and check that it succeeds")
	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Update.Concurrency))
	for i, vec := range ds.Train[cfg.Update.Offset : cfg.Update.Offset+cfg.Update.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Update.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Update(ctx, &payload.Update_Request{
				Vector: &payload.Object_Vector{
					Id:        id,
					Vector:    vec,
					Timestamp: ts,
				},
				Config: &payload.Update_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: false,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to update vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to update vector: %v", err)
				}
			}
			t.Logf("vector %v id %s updated on #4 to %s", vec, id, res.String())
			return nil
		}))
	}
	eg.Wait()

	sleep(t, cfg.Index.WaitAfterInsert)

	t.Log("starting test #5 run Update with SkipStrictExistCheck=false, and check that the internal Remove Operation returns a NotFound error.")
	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Update.Concurrency))
	for i, vec := range ds.Train[cfg.Update.Offset : cfg.Update.Offset+cfg.Update.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Update.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			_, err := client.Update(ctx, &payload.Update_Request{
				Vector: &payload.Object_Vector{
					Id:        id,
					Vector:    vec,
					Timestamp: ts,
				},
				Config: &payload.Update_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: false,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if !ok || st == nil || st.Code() != codes.AlreadyExists {
					t.Errorf("update vector response is not AlreadyExists: %v with SkipStrictExistCheck=false", err)
				}
			}
			t.Logf("received a NotFound error on #5: %s", err.Error())
			return nil
		}))
	}
	eg.Wait()

	t.Log("starting test #6 run Update with SkipStrictExistCheck=true & 4 and check that it succeess")
	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Update.Concurrency))
	for i, vec := range ds.Train[cfg.Update.Offset : cfg.Update.Offset+cfg.Update.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Update.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Update(ctx, &payload.Update_Request{
				Vector: &payload.Object_Vector{
					Id:        id,
					Vector:    vec,
					Timestamp: ts,
				},
				Config: &payload.Update_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: true,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to update vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to update vector: %v", err)
				}
			}
			t.Logf("vector %v id %s updated on #6 to %s", vec, id, res.String())
			return nil
		}))
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Remove.Concurrency))
	for i := range ds.Train[cfg.Remove.Offset : cfg.Remove.Offset+cfg.Remove.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Remove.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Remove(ctx, &payload.Remove_Request{
				Id: &payload.Object_ID{Id: id},
				Config: &payload.Remove_Config{
					Timestamp:            ts,
					SkipStrictExistCheck: cfg.Remove.SkipStrictExistCheck,
				},
			})
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to remove vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to remove vector: %v", err)
				}
			}
			t.Logf("id %s'd vector removed to %s", id, res.String())
			return nil
		}))
	}
	eg.Wait()

	eg, _ = errgroup.New(ctx)
	eg.SetLimit(int(cfg.Upsert.Concurrency))
	for i, vec := range ds.Train[cfg.Upsert.Offset : cfg.Upsert.Offset+cfg.Upsert.Num] {
		id := strconv.Itoa(i)
		ts := cfg.Upsert.Timestamp
		if ts == 0 {
			ts = timestamp
		}
		eg.Go(safety.RecoverFunc(func() error {
			res, err := client.Upsert(ctx, &payload.Upsert_Request{
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
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					t.Errorf("failed to upsert vector: %v, status: %s", err, st.String())
				} else {
					t.Errorf("failed to upsert vector: %v", err)
				}
			}
			t.Logf("vector %v id %s upserted to %s", vec, id, res.String())
			return nil
		}))
	}
	eg.Wait()

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
