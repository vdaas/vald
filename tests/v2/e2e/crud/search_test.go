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
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func recall(t *testing.T, resultIDs []string, neighbors []int) (recall float64) {
	t.Helper()
	ns := map[string]struct{}{}
	for _, n := range neighbors {
		ns[strconv.Itoa(n)] = struct{}{}
	}

	for _, r := range resultIDs {
		if _, ok := ns[r]; ok {
			recall++
		}
	}

	return recall / float64(len(neighbors))
}

func calculateRecall(t *testing.T, res *payload.Search_Response, idx int) (rc float64) {
	t.Helper()
	topKIDs := make([]string, 0, len(res.GetResults()))
	for _, d := range res.GetResults() {
		topKIDs = append(topKIDs, d.GetId())
	}

	if len(topKIDs) == 0 {
		t.Errorf("empty result is returned for test ID %s: %#v", res.GetRequestId(), topKIDs)
		return
	}
	rc = recall(t, topKIDs, ds.Neighbors[idx][:len(topKIDs)])
	return rc
}

func (r *runner) search(t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution) {
	t.Helper()
	unarySearch(t, ctx, test, neighbors, plan, r.client.Search)
}

func (r *runner) linearSearch(t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution) {
	t.Helper()
	unarySearch(t, ctx, test, neighbors, plan, r.client.LinearSearch)
}

func (r *runner) searchByID(t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution) {
	t.Helper()
	unarySearchByID(t, ctx, train, neighbors, plan, r.client.SearchByID)
}

func (r *runner) linearSearchByID(t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution) {
	t.Helper()
	unarySearchByID(t, ctx, train, neighbors, plan, r.client.LinearSearchByID)
}

func (r *runner) multiSearch(t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution) {
	t.Helper()
	multiSearch(t, ctx, test, neighbors, plan, r.client.MultiSearch)
}

func (r *runner) multiLinearSearch(t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution) {
	t.Helper()
	multiSearch(t, ctx, test, neighbors, plan, r.client.MultiLinearSearch)
}

func (r *runner) multiSearchByID(t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution) {
	t.Helper()
	multiSearchByID(t, ctx, train, neighbors, plan, r.client.MultiSearchByID)
}

func (r *runner) multiLinearSearchByID(t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution) {
	t.Helper()
	multiSearchByID(t, ctx, train, neighbors, plan, r.client.MultiLinearSearchByID)
}

func unarySearch(t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution, do func(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error)) {
	t.Helper()
	eg, ctx := errgroup.New(ctx)
	eg.SetLimit(int(plan.Concurrency))
	for i, vec := range test {
		for _, query := range plan.SearchConfig {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			eg.Go(safety.RecoverFunc(func() (err error) {
				var ratio *wrapperspb.FloatValue
				if query.Ratio != 0 {
					ratio = wrapperspb.Float(query.Ratio)
				} else {
					ratio = nil
				}
				var to time.Duration
				if query.Timeout != "" {
					to, err = query.Timeout.Duration()
					if err != nil {
						t.Errorf("failed to parse timeout duration: %s", err)
						return nil
					}
				}
				res, err := do(ctx, &payload.Search_Request{
					Vector: vec,
					Config: &payload.Search_Config{
						RequestId:            rid,
						Num:                  query.K,
						Radius:               query.Radius,
						Epsilon:              query.Epsilon,
						Timeout:              to.Nanoseconds(),
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

func unarySearchByID(t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution, do func(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)) {
	t.Helper()
	eg, ctx := errgroup.New(ctx)
	eg.SetLimit(int(plan.Concurrency))
	for i, vec := range train {
		for _, query := range plan.SearchConfig {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			eg.Go(safety.RecoverFunc(func() (err error) {
				var ratio *wrapperspb.FloatValue
				if query.Ratio != 0 {
					ratio = wrapperspb.Float(query.Ratio)
				} else {
					ratio = nil
				}
				var to time.Duration
				if query.Timeout != "" {
					to, err = query.Timeout.Duration()
					if err != nil {
						t.Errorf("failed to parse timeout duration: %s", err)
						return nil
					}
				}
				res, err := do(ctx, &payload.Search_IDRequest{
					Id: id,
					Config: &payload.Search_Config{
						RequestId:            rid,
						Num:                  query.K,
						Radius:               query.Radius,
						Epsilon:              query.Epsilon,
						Timeout:              to.Nanoseconds(),
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

func multiSearch(t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution, do func(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)) {
	t.Helper()
	eg, ctx := errgroup.New(ctx)
	eg.SetLimit(int(plan.Concurrency))
	msreq := &payload.Search_MultiRequest{
		Requests: make([]*payload.Search_Request, 0, plan.BulkSize),
	}
	for i, vec := range test {
		for _, query := range plan.SearchConfig {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			var ratio *wrapperspb.FloatValue
			if query.Ratio != 0 {
				ratio = wrapperspb.Float(query.Ratio)
			} else {
				ratio = nil
			}
			var (
				to  time.Duration
				err error
			)

			if query.Timeout != "" {
				to, err = query.Timeout.Duration()
				if err != nil {
					t.Errorf("failed to parse timeout duration: %s", err)
					continue
				}
			}
			msreq.Requests = append(msreq.Requests, &payload.Search_Request{
				Vector: vec,
				Config: &payload.Search_Config{
					RequestId:            rid,
					Num:                  query.K,
					Radius:               query.Radius,
					Epsilon:              query.Epsilon,
					Timeout:              to.Nanoseconds(),
					AggregationAlgorithm: query.Algorithm,
					MinNum:               query.MinNum,
					Ratio:                ratio,
					Nprobe:               query.Nprobe,
				},
			})
			if len(msreq.GetRequests()) >= int(plan.BulkSize) {
				req := msreq.CloneVT()
				msreq.Reset()
				eg.Go(safety.RecoverFunc(func() error {
					res, err := do(ctx, req)
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
}

func multiSearchByID(t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution, do func(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)) {
	t.Helper()
	eg, ctx := errgroup.New(ctx)
	eg.SetLimit(int(plan.Concurrency))
	msreq := &payload.Search_MultiIDRequest{
		Requests: make([]*payload.Search_IDRequest, 0, plan.BulkSize),
	}
	for i := range train {
		for _, query := range plan.SearchConfig {
			id := strconv.Itoa(i)
			rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
			var ratio *wrapperspb.FloatValue
			if query.Ratio != 0 {
				ratio = wrapperspb.Float(query.Ratio)
			} else {
				ratio = nil
			}
			var (
				to  time.Duration
				err error
			)

			if query.Timeout != "" {
				to, err = query.Timeout.Duration()
				if err != nil {
					t.Errorf("failed to parse timeout duration: %s", err)
					continue
				}
			}
			msreq.Requests = append(msreq.Requests, &payload.Search_IDRequest{
				Id: id,
				Config: &payload.Search_Config{
					RequestId:            rid,
					Num:                  query.K,
					Radius:               query.Radius,
					Epsilon:              query.Epsilon,
					Timeout:              to.Nanoseconds(),
					AggregationAlgorithm: query.Algorithm,
					MinNum:               query.MinNum,
					Ratio:                ratio,
					Nprobe:               query.Nprobe,
				},
			})
			if len(msreq.GetRequests()) >= int(plan.BulkSize) {
				req := msreq.CloneVT()
				msreq.Reset()
				eg.Go(safety.RecoverFunc(func() error {
					res, err := do(ctx, req)
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
}
