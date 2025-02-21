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
	// "github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/tests/v2/e2e/kubernetes"
	k8sclient "github.com/vdaas/vald/tests/v2/e2e/kubernetes"
)

func TestE2EUnaryRolloutRestartAgentCRUD(t *testing.T) {
	timestamp := time.Now().UnixNano()

	t.Log(cfg, ctx)

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

	// TODO: inifinite search
	done := make(chan struct{})
	if cfg.Search.Num != 0 {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
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
			}
		}()

		eg, _ = errgroup.New(ctx)
		eg.Go(safety.RecoverFunc(func() error {
			kclient, err := k8sclient.New(cfg.Kubernetes.KubeConfig)
			if err != nil {
				t.Logf("failed to create kubernetes client: %v", err)
				return err
			}
			r := kubernetes.StatefulSet{
				Name: "vald-agent",
				Namespace: "defualt",
			}
			clientset := kclient.Clientset()
			err = kubernetes.RolloutRestart(ctx, clientset, r)
			if err != nil {
				t.Logf("failed to rollout restart: %v", err)
				return err
			}
			return kubernetes.WaitForRestart(ctx, clientset, r)
		}))
		err := eg.Wait()
		close(done)
		if err != nil {
			t.Fatalf("failed to rollout restart: %s", err.Error())
		}
	}

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
