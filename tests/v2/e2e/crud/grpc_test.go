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

// Package crud provides end-to-end tests using ann-benchmarks datasets.
package crud

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/iter"
	"github.com/vdaas/vald/internal/jsonpath"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"github.com/vdaas/vald/tests/v2/e2e/metrics"
	"google.golang.org/protobuf/encoding/protojson"
)

// Type aliases for generic search functions.
type (
	// grpcCall is a generic function type for making gRPC calls.
	grpcCall[Q, R proto.Message] func(ctx context.Context, query Q, opts ...grpc.CallOption) (response R, err error)
	// newStream is a generic type for functions that create a new gRPC stream.
	newStream[S grpc.ClientStream] func(ctx context.Context, opts ...grpc.CallOption) (S, error)
	// newRequest is a function type that creates a new request.
	newRequest[Q proto.Message] func(t *testing.T, idx uint64, id string, vec []float32, e *config.Execution) Q
	// newMultiRequest is a generic type for functions that build bulk search requests.
	newMultiRequest[R, S proto.Message] func(t *testing.T, reqs ...R) S
	// callback is a function type that processes the response and error from a gRPC call.
	callback[R proto.Message] func(t *testing.T, idx uint64, res R, err error) bool
)

func passThrough[M proto.Message](t *testing.T, msg M) any {
	t.Helper()
	return msg
}

func emptyCallback[M proto.Message](name string) callback[M] {
	return func(t *testing.T, _ uint64, _ M, err error) bool {
		t.Helper()
		if err != nil {
			log.Errorf("%s operation returned error: %v", name, err)
			return false
		}
		return true
	}
}

func printCallback[M proto.Message](unwrap func(t *testing.T, msg M) any) callback[M] {
	return func(t *testing.T, idx uint64, msg M, err error) bool {
		t.Helper()
		if err != nil {
			log.Errorf("idx: %d operation returned error: %v", idx, err)
			return true
		}
		log.Infof("idx: %d operation returned result: %v", idx, unwrap(t, msg))
		return true
	}
}

func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint64:
		// float64 can only precisely represent integers up to 2^53.
		if val <= 1<<53 {
			return float64(val), true
		}
		return 0, false
	default:
		return 0, false
	}
}

func compare(a, b any) (float64, float64, bool) {
	aT, ok1 := toFloat64(a)
	bT, ok2 := toFloat64(b)
	return aT, bT, ok1 && ok2
}

func handleGRPCWithStatusCode(
	t *testing.T, err error, code codes.Code, res proto.Message, plan *config.Execution,
) error {
	t.Helper()
	if len(plan.Expect) == 0 {
		return err
	}

	var protoJSON []byte
	if res != nil {
		marshaller := protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: false,
		}
		protoJSON, err = marshaller.Marshal(res)
		if err != nil {
			return fmt.Errorf("failed to marshal proto: %w", err)
		}
	}

	errs := make([]error, 0, len(plan.Expect)+1)
	if err != nil {
		errs = append(errs, err)
	}
	for _, expect := range plan.Expect {
		if expect.StatusCode != "" && !expect.StatusCode.Equals(code.String()) {
			err := fmt.Errorf("unexpected gRPC response received expected: %s, got: %s", expect.StatusCode, code)
			log.Errorf("❌ assert failed, err: %v", err)
			errs = append(errs, err)
			continue
		}
		if expect.Value != nil {
			val, err := jsonpath.JSONPathEval(protoJSON, expect.Path)
			if err != nil {
				log.Errorf("❌ assert failed, err: %v", err)
				errs = append(errs, fmt.Errorf("failed to evaluate JSONPath: %s, JSON: %s, err: %s", expect.Path, protoJSON, err))
				continue
			}
			commonErr := fmt.Errorf("assert_%v failed, JSONPath: %s, expected: %v actual: %v", expect.Op, expect.Path, expect.Value, val)
			switch op := expect.Op; op {
			case config.Eq, config.Ne:
				isMatched := reflect.DeepEqual(val, expect.Value) || fmt.Sprintf("%v", val) == fmt.Sprintf("%v", expect.Value)
				if isMatched && op == config.Ne || !isMatched && op == config.Eq {
					errs = append(errs, commonErr)
					continue
				}
			case config.Gt, config.Ge, config.Lt, config.Le:
				a, b, ok := compare(val, expect.Value)
				if !ok {
					errs = append(errs, fmt.Errorf("assert_%v failed, cannot compare values, JSONPath: %s, expected: %v actual: %v", expect.Op, expect.Path, expect.Value, val))
					continue
				}
				switch op {
				case config.Gt:
					if a <= b {
						errs = append(errs, commonErr)
						continue
					}
				case config.Ge:
					if a < b {
						errs = append(errs, commonErr)
						continue
					}
				case config.Lt:
					if a >= b {
						errs = append(errs, commonErr)
						continue
					}
				case config.Le:
					if a > b {
						errs = append(errs, commonErr)
						continue
					}
				}
			default:
				errs = append(errs, fmt.Errorf("unsupported operator '%s' for JSONPath %s", expect.Op, expect.Path))
				continue
			}
			log.Infof("✅ assert_%v passed, expected: %v actual: %v", expect.Op, expect.Value, val)
		}
		return nil
	}

	err = errors.Join(errs...)
	log.Errorf("❌ assert failed, err: %v", err)
	return err
}

// handleGRPCCall centralizes the gRPC error handling, logging and assertion.
// It compares the error's status code with the expected codes from the plan.
// If the error is expected, it logs a message; otherwise, it logs an error.
// If the results do not match, it logs an error.
func handleGRPCCall(
	t *testing.T, err error, res proto.Message, plan *config.Execution,
) (code codes.Code, msg string, rerr error) {
	t.Helper()
	if err != nil {
		if st, ok := status.FromError(err); ok && st != nil {
			msg = st.String()
			code = st.Code()
			rerr = errors.Wrapf(err, "gRPC request received: %s", msg)
		}
	} else {
		code = codes.OK
	}
	rerr = handleGRPCWithStatusCode(t, err, code, res, plan)
	return code, msg, rerr
}

func single[Q, R proto.Message](
	t *testing.T,
	ctx context.Context,
	idx uint64,
	plan *config.Execution,
	req Q,
	call grpcCall[Q, R],
	callback ...callback[R],
) (err error) {
	t.Helper()
	if plan == nil {
		return nil
	}

	queuedAt := time.Now()
	if plan.BaseConfig != nil && plan.BaseConfig.Limiter != nil {
		plan.BaseConfig.Limiter.Wait(ctx)
	}
	startedAt := time.Now()
	res, err := call(ctx, req)
	endedAt := time.Now()

	st, msg, rerr := handleGRPCCall(t, err, res, plan)
	if plan.Metrics != nil && plan.Metrics.Enabled && plan.Collector != nil {
		rr := metrics.GetRequestResult()
		defer metrics.PutRequestResult(rr)
		rr.RequestID = strconv.FormatUint(idx, 10)
		rr.Status = st
		rr.Err = err
		rr.Msg = msg
		rr.QueuedAt = queuedAt
		rr.StartedAt = startedAt
		rr.EndedAt = endedAt
		plan.Collector.Record(ctx, idx, rr)
	}
	if rerr != nil && errors.IsNot(err, rerr) {
		return rerr
	}
	for _, cb := range callback {
		if cb != nil {
			if !cb(t, idx, res, err) {
				return fmt.Errorf("callback failed for idx: %d, err: %v", idx, err)
			}
		}
	}
	return nil
}

func unary[Q, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
	call grpcCall[Q, R],
	newReq newRequest[Q],
	callback ...callback[R],
) error {
	t.Helper()
	// Create an error group to manage concurrent requests.
	eg, ctx := errgroup.New(ctx)
	// Set the concurrency limit from the plan configuration.
	if plan != nil && plan.BaseConfig != nil {
		// Set the concurrency limit from the plan configuration.
		eg.SetLimit(int(plan.Parallelism))
	}
	for i, vec := range data.Seq2(ctx) {
		// Copy id to avoid data race.
		idx := i
		// Execute request in a goroutine.
		eg.Go(func() error {
			return single(t, ctx, idx, plan, newReq(t, idx, strconv.FormatUint(idx, 10), vec, plan), call, callback...)
		})
	}
	// Wait for all goroutines to complete.
	return eg.Wait()
}

func multi[Q, M, R proto.Message](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
	call grpcCall[M, R],
	addReqs newRequest[Q],
	toReq newMultiRequest[Q, M],
	callbacks ...callback[R],
) error {
	t.Helper()
	eg, ctx := errgroup.New(ctx)
	// Set the concurrency limit from the plan configuration.
	if plan != nil && plan.BaseConfig != nil {
		// Set the concurrency limit from the plan configuration.
		eg.SetLimit(int(plan.Parallelism))
	}
	var bulkSize uint64
	if plan.BulkSize < 2 {
		bulkSize = 10
	} else {
		bulkSize = plan.BulkSize
	}

	// Initialize a slice to hold the bulk requests.
	reqs := make([]Q, 0, bulkSize)
	for i, vec := range data.Seq2(ctx) {
		id := strconv.FormatUint(i, 10)
		// Append a new request to the bulk slice using the provided builder.
		reqs = append(reqs, addReqs(t, i, id, vec, plan))
		// If the bulk size is reached, send the batch.
		if len(reqs) >= int(bulkSize) {
			// Capture the current batch.
			batch := slices.Clone(reqs)
			idx := i
			// Meset the bulk request slice for the next batch.
			reqs = reqs[:0]
			eg.Go(func() error {
				return single(t, ctx, idx, plan, toReq(t, batch...), call, callbacks...)
			})
		}
	}
	eg.Go(func() error {
		return single(t, ctx, data.Len(), plan, toReq(t, reqs...), call, callbacks...)
	})
	return eg.Wait()
}

func stream[Q, R proto.Message, S grpc.TypedClientStream[Q, R]](
	t *testing.T,
	ctx context.Context,
	data iter.Cycle[[][]float32, []float32],
	plan *config.Execution,
	newStream newStream[S],
	newReq newRequest[Q],
	callbacks ...callback[R],
) {
	t.Helper()
	// Create a new stream using the provided stream function.
	stream, err := newStream(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	if any(stream) == nil {
		err = errors.ErrGRPCClientStreamNotFound
		t.Error(err)
		return
	}
	rchan := make(chan *metrics.RequestResult, int(plan.Parallelism)*2)
	var idx atomic.Uint64
	// Use a bidirectional stream client to send requests and receive responses.
	err = grpc.BidirectionalStreamClient(stream, int(plan.Parallelism), func() (query Q, ok bool) {
		id := idx.Load()
		idx.Add(1)
		// If we have processed all vectors, return nil to close the stream.
		if id >= data.Len() {
			close(rchan)
			return query, false
		}
		rr := metrics.GetRequestResult()
		rr.RequestID = strconv.FormatUint(id, 10)

		// Build the modify configuration and return the request.
		query = newReq(t, id, strconv.FormatUint(id, 10), data.At(id), plan)

		rr.QueuedAt = time.Now()
		if plan.BaseConfig != nil && plan.BaseConfig.Limiter != nil {
			plan.BaseConfig.Limiter.Wait(stream.Context())
		}
		rr.StartedAt = time.Now()
		select {
		case <-ctx.Done():
		case rchan <- rr:
		}
		return query, true
	}, func(res R, err error) bool {
		endedAt := time.Now()
		var rr *metrics.RequestResult
		select {
		case <-ctx.Done():
		case rr = <-rchan:
		}
		if rr == nil {
			rr = new(metrics.RequestResult)
		}
		defer metrics.PutRequestResult(rr)
		var rerr error
		rr.Status, rr.Msg, rerr = handleGRPCCall(t, err, res, plan)
		if plan.Metrics != nil && plan.Metrics.Enabled && plan.Collector != nil {
			rr.Err = err
			rr.EndedAt = endedAt
			id, err := strconv.ParseUint(rr.RequestID, 10, 64)
			if err != nil {
				id = 0
			}
			plan.Collector.Record(ctx, id, rr)
		}
		if rerr != nil && errors.IsNot(err, rerr) {
			return false
		}

		id, err := strconv.ParseUint(rr.RequestID, 10, 0)
		if err != nil {
			id = 0
		}

		for _, cb := range callbacks {
			if cb != nil {
				if !cb(t, id, res, err) {
					return false
				}
			}
		}
		return true
	})
	if err != nil {
		t.Errorf("failed to complete stream: %v", err)
	}
}
