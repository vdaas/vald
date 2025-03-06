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
	"github.com/vdaas/vald/internal/net/grpc/proto"
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

func (r *runner) processSearch(
	t *testing.T,
	ctx context.Context,
	test, train [][]float32,
	neighbors [][]int,
	plan *config.Execution,
) {
	t.Helper()
	switch plan.Type {
	case config.OpSearch:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			r.search(t, ctx, test, neighbors, plan)
		case config.OperationMultiple:
			r.multiSearch(t, ctx, test, neighbors, plan)
		case config.OperationStream:
			r.streamSearch(t, ctx, test, neighbors, plan)
		}
	case config.OpSearchByID:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			r.searchByID(t, ctx, train, neighbors, plan)
		case config.OperationMultiple:
			r.multiSearchByID(t, ctx, train, neighbors, plan)
		case config.OperationStream:
			r.streamSearchByID(t, ctx, train, neighbors, plan)
		}
	case config.OpLinearSearch:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			r.linearSearch(t, ctx, test, neighbors, plan)
		case config.OperationMultiple:
			r.multiLinearSearch(t, ctx, test, neighbors, plan)
		case config.OperationStream:
			r.streamLinearSearch(t, ctx, test, neighbors, plan)
		}
	case config.OpLinearSearchByID:
		switch plan.Mode {
		case config.OperationUnary, config.OperationOther:
			r.linearSearchByID(t, ctx, train, neighbors, plan)
		case config.OperationMultiple:
			r.multiLinearSearchByID(t, ctx, train, neighbors, plan)
		case config.OperationStream:
			r.streamLinearSearchByID(t, ctx, train, neighbors, plan)
		}
	}
}

func (r *runner) search(
	t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	unarySearch(t, ctx, test, neighbors, plan, r.client.Search,
		func(_ string, vec []float32, scfg *payload.Search_Config) *payload.Search_Request {
			return &payload.Search_Request{
				Vector: vec,
				Config: scfg,
			}
		})
}

func (r *runner) linearSearch(
	t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	unarySearch(t, ctx, test, neighbors, plan, r.client.LinearSearch,
		func(id string, vec []float32, scfg *payload.Search_Config) *payload.Search_Request {
			return &payload.Search_Request{
				Vector: vec,
				Config: scfg,
			}
		})
}

func (r *runner) searchByID(
	t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	unarySearch(t, ctx, train, neighbors, plan, r.client.SearchByID,
		func(id string, _ []float32, scfg *payload.Search_Config) *payload.Search_IDRequest {
			return &payload.Search_IDRequest{
				Id:     id,
				Config: scfg,
			}
		})
}

func (r *runner) linearSearchByID(
	t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	unarySearch(t, ctx, train, neighbors, plan, r.client.LinearSearchByID,
		func(id string, _ []float32, scfg *payload.Search_Config) *payload.Search_IDRequest {
			return &payload.Search_IDRequest{
				Id:     id,
				Config: scfg,
			}
		})
}

func (r *runner) multiSearch(
	t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	multiSearch(t, ctx, test, neighbors, plan, r.client.MultiSearch,
		func(id string, vec []float32, scfg *payload.Search_Config) *payload.Search_Request {
			return &payload.Search_Request{
				Vector: vec,
				Config: scfg,
			}
		},
		func(reqs []*payload.Search_Request) *payload.Search_MultiRequest {
			return &payload.Search_MultiRequest{
				Requests: reqs,
			}
		})
}

func (r *runner) multiLinearSearch(
	t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	multiSearch(t, ctx, test, neighbors, plan, r.client.MultiLinearSearch,
		func(id string, vec []float32, scfg *payload.Search_Config) *payload.Search_Request {
			return &payload.Search_Request{
				Vector: vec,
				Config: scfg,
			}
		},
		func(reqs []*payload.Search_Request) *payload.Search_MultiRequest {
			return &payload.Search_MultiRequest{
				Requests: reqs,
			}
		})
}

func (r *runner) multiSearchByID(
	t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	multiSearch(t, ctx, train, neighbors, plan, r.client.MultiSearchByID,
		func(id string, _ []float32, scfg *payload.Search_Config) *payload.Search_IDRequest {
			return &payload.Search_IDRequest{
				Id:     id,
				Config: scfg,
			}
		},
		func(reqs []*payload.Search_IDRequest) *payload.Search_MultiIDRequest {
			return &payload.Search_MultiIDRequest{
				Requests: reqs,
			}
		})
}

func (r *runner) multiLinearSearchByID(
	t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	multiSearch(t, ctx, train, neighbors, plan, r.client.MultiLinearSearchByID,
		func(id string, _ []float32, scfg *payload.Search_Config) *payload.Search_IDRequest {
			return &payload.Search_IDRequest{
				Id:     id,
				Config: scfg,
			}
		},
		func(reqs []*payload.Search_IDRequest) *payload.Search_MultiIDRequest {
			return &payload.Search_MultiIDRequest{
				Requests: reqs,
			}
		})
}

func (r *runner) streamSearch(
	t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	streamSearch(t, ctx, test, neighbors, plan, r.client.StreamSearch,
		func(_ string, vec []float32, scfg *payload.Search_Config) *payload.Search_Request {
			return &payload.Search_Request{
				Vector: vec,
				Config: scfg,
			}
		})
}

func (r *runner) streamLinearSearch(
	t *testing.T, ctx context.Context, test [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	streamSearch(t, ctx, test, neighbors, plan, r.client.StreamLinearSearch,
		func(_ string, vec []float32, scfg *payload.Search_Config) *payload.Search_Request {
			return &payload.Search_Request{
				Vector: vec,
				Config: scfg,
			}
		})
}

func (r *runner) streamSearchByID(
	t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	streamSearch(t, ctx, train, neighbors, plan, r.client.StreamSearchByID,
		func(id string, _ []float32, scfg *payload.Search_Config) *payload.Search_IDRequest {
			return &payload.Search_IDRequest{
				Id:     id,
				Config: scfg,
			}
		})
}

func (r *runner) streamLinearSearchByID(
	t *testing.T, ctx context.Context, train [][]float32, neighbors [][]int, plan *config.Execution,
) {
	t.Helper()
	streamSearch(t, ctx, train, neighbors, plan, r.client.StreamLinearSearchByID,
		func(id string, _ []float32, scfg *payload.Search_Config) *payload.Search_IDRequest {
			return &payload.Search_IDRequest{
				Id:     id,
				Config: scfg,
			}
		})
}

func unarySearch[R proto.Message](
	t *testing.T,
	ctx context.Context,
	data [][]float32,
	neighbors [][]int,
	plan *config.Execution,
	do func(ctx context.Context, in R, opts ...grpc.CallOption) (*payload.Search_Response, error),
	newReq func(id string, vec []float32, scfg *payload.Search_Config) R,
) {
	t.Helper()
	eg, ctx := errgroup.New(ctx)
	eg.SetLimit(int(plan.Concurrency))
	for i, vec := range data {
		for _, q := range plan.SearchConfig {
			query := q
			id := strconv.Itoa(i)
			eg.Go(safety.RecoverFunc(func() (err error) {
				rid := id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)]
				res, err := do(ctx, newReq(id, vec, &payload.Search_Config{
					RequestId: rid,
					Num:       query.K,
					Radius:    query.Radius,
					Epsilon:   query.Epsilon,
					Timeout: func() int64 {
						if query.Timeout != "" {
							to, err := query.Timeout.Duration()
							if err == nil {
								return to.Nanoseconds()
							}
						}
						return time.Second.Nanoseconds()
					}(),
					AggregationAlgorithm: query.Algorithm,
					MinNum:               query.MinNum,
					Ratio: func() *wrapperspb.FloatValue {
						if query.Ratio != 0 {
							return wrapperspb.Float(query.Ratio)
						}
						return nil
					}(),
					Nprobe: query.Nprobe,
				}))
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

func multiSearch[R, S proto.Message](
	t *testing.T,
	ctx context.Context,
	data [][]float32,
	neighbors [][]int,
	plan *config.Execution,
	do func(ctx context.Context, in S, opts ...grpc.CallOption) (*payload.Search_Responses, error),
	addReqs func(id string, vec []float32, scfg *payload.Search_Config) R,
	toReq func([]R) S,
) {
	t.Helper()
	eg, ctx := errgroup.New(ctx)
	eg.SetLimit(int(plan.Concurrency))
	reqs := make([]R, 0, plan.BulkSize)
	msreq := &payload.Search_MultiRequest{
		Requests: make([]*payload.Search_Request, 0, plan.BulkSize),
	}
	for i, vec := range data {
		for _, query := range plan.SearchConfig {
			id := strconv.Itoa(i)
			reqs = append(reqs, addReqs(id, vec, &payload.Search_Config{
				RequestId: id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)],
				Num:       query.K,
				Radius:    query.Radius,
				Epsilon:   query.Epsilon,
				Timeout: func() int64 {
					if query.Timeout != "" {
						to, err := query.Timeout.Duration()
						if err == nil {
							return to.Nanoseconds()
						}
					}
					return time.Second.Nanoseconds()
				}(),
				AggregationAlgorithm: query.Algorithm,
				MinNum:               query.MinNum,
				Ratio: func() *wrapperspb.FloatValue {
					if query.Ratio != 0 {
						return wrapperspb.Float(query.Ratio)
					}
					return nil
				}(),
				Nprobe: query.Nprobe,
			}))
			if len(reqs) >= int(plan.BulkSize) {
				req := msreq.CloneVT()
				msreq.Reset()
				eg.Go(safety.RecoverFunc(func() error {
					res, err := do(ctx, toReq(reqs))
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

func streamSearch[S grpc.ClientStream, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data [][]float32,
	neighbors [][]int,
	plan *config.Execution,
	newStream func(ctx context.Context, opts ...grpc.CallOption) (S, error),
	newReq func(id string, vec []float32, scfg *payload.Search_Config) R,
) {
	t.Helper()
	stream, err := newStream(ctx)
	if err != nil {
		t.Error(err)
	}
	qidx := 0
	idx := 0
	err = grpc.BidirectionalStreamClient(stream, func() *R {
		if len(data) <= idx {
			return nil
		}
		id := strconv.Itoa(idx)
		if len(plan.SearchConfig) < qidx {
			qidx = 0
			idx++
		}
		query := plan.SearchConfig[qidx]
		qidx++
		req := newReq(id, data[idx], &payload.Search_Config{
			RequestId: id + "-" + payload.Search_AggregationAlgorithm_name[int32(query.Algorithm)],
			Num:       query.K,
			Radius:    query.Radius,
			Epsilon:   query.Epsilon,
			Timeout: func() int64 {
				if query.Timeout != "" {
					to, err := query.Timeout.Duration()
					if err == nil {
						return to.Nanoseconds()
					}
				}
				return time.Second.Nanoseconds()
			}(),
			AggregationAlgorithm: query.Algorithm,
			MinNum:               query.MinNum,
			Ratio: func() *wrapperspb.FloatValue {
				if query.Ratio != 0 {
					return wrapperspb.Float(query.Ratio)
				}
				return nil
			}(),
			Nprobe: query.Nprobe,
		})
		return &req
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
}
