// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package grpc

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/slices"
)

type Aggregator interface {
	Start(ctx context.Context)
	Send(ctx context.Context, data *payload.Search_Response)
	Result() *payload.Search_Response
}

type DistPayload struct {
	raw      *payload.Object_Distance
	distance *big.Float
}

func (s *server) aggregationSearch(ctx context.Context, aggr Aggregator, cfg *payload.Search_Config,
	f func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error)) (
	res *payload.Search_Response, err error,
) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "aggregationSearch"), apiName+"/aggregationSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	num := int(cfg.GetNum())
	min := int(cfg.GetMinNum())
	eg, ectx := errgroup.New(ctx)
	var cancel context.CancelFunc
	var timeout time.Duration
	if to := cfg.GetTimeout(); to != 0 {
		timeout = time.Duration(to)
	} else {
		timeout = s.timeout
	}

	ectx, cancel = context.WithTimeout(ectx, timeout)
	aggr.Start(ectx)
	eg.Go(safety.RecoverFunc(func() error {
		defer cancel()
		return s.gateway.BroadCast(ectx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/aggregationSearch/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := f(sctx, vc, copts...)
			if err != nil {
				switch {
				case errors.Is(err, context.Canceled),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.StatusCodeCancelled(
							errdetails.ValdGRPCResourceTypePrefix +
								"/vald.v1.search.BroadCast/" +
								target + " canceled: " + err.Error())...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.StatusCodeDeadlineExceeded(
							errdetails.ValdGRPCResourceTypePrefix +
								"/vald.v1.search.BroadCast/" +
								target + " deadline_exceeded: " + err.Error())...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
				default:
					st, msg, err := status.ParseError(err, codes.Internal, "failed to parse search gRPC error response",
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
					switch st.Code() {
					case codes.Internal,
						codes.Unavailable,
						codes.ResourceExhausted:
						log.Warn(err)
						return err
					case codes.NotFound,
						codes.Aborted,
						codes.InvalidArgument:
						return nil
					}
				}
				log.Debug(err)
				return nil
			}
			if r == nil || len(r.GetResults()) == 0 {
				select {
				case <-sctx.Done():
					err = status.WrapWithNotFound("failed to process search request", errors.ErrEmptySearchResult,
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
					log.Debug(err)
					return nil
				default:
					r, err = f(sctx, vc, copts...)
					if err != nil {
						switch {
						case errors.Is(err, context.Canceled),
							errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
							if sspan != nil {
								sspan.RecordError(err)
								sspan.SetAttributes(trace.StatusCodeCancelled(
									errdetails.ValdGRPCResourceTypePrefix +
										"/vald.v1.search.BroadCast/" +
										target + " canceled: " + err.Error())...)
								sspan.SetStatus(trace.StatusError, err.Error())
							}
						case errors.Is(err, context.DeadlineExceeded),
							errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
							if sspan != nil {
								sspan.RecordError(err)
								sspan.SetAttributes(trace.StatusCodeDeadlineExceeded(
									errdetails.ValdGRPCResourceTypePrefix +
										"/vald.v1.search.BroadCast/" +
										target + " deadline_exceeded: " + err.Error())...)
								sspan.SetStatus(trace.StatusError, err.Error())
							}
						default:
							st, msg, err := status.ParseError(err, codes.Internal, "failed to parse search gRPC error response",
								&errdetails.ResourceInfo{
									ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
									ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
								})
							if sspan != nil {
								sspan.RecordError(err)
								sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
								sspan.SetStatus(trace.StatusError, err.Error())
							}
							switch st.Code() {
							case codes.Internal,
								codes.Unavailable,
								codes.ResourceExhausted:
								log.Warn(err)
								return err
							case codes.NotFound,
								codes.Aborted,
								codes.InvalidArgument:
								return nil
							}
						}
						log.Debug(err)
						return nil
					}
					if r == nil || len(r.GetResults()) == 0 {
						err = status.WrapWithNotFound("failed to process search request", errors.ErrEmptySearchResult,
							&errdetails.ResourceInfo{
								ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
								ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
							})
						if sspan != nil {
							sspan.RecordError(err)
							sspan.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
							sspan.SetStatus(trace.StatusError, err.Error())
						}
						log.Debug(err)
						return nil
					}
				}
			}
			aggr.Send(sctx, r)
			return nil
		})
	}))

	<-ectx.Done() // Blocking here

	err = eg.Wait()
	if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
		err = status.WrapWithInternal("search API connection not found", err,
			&errdetails.RequestInfo{
				RequestId:   cfg.GetRequestId(),
				ServingData: errdetails.Serialize(cfg),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	res = aggr.Result()
	if num != 0 && len(res.GetResults()) > num {
		res.Results = res.GetResults()[:num]
	}

	if errors.Is(ectx.Err(), context.DeadlineExceeded) {
		if len(res.GetResults()) == 0 {
			err = status.WrapWithDeadlineExceeded(
				"error search result length is 0 due to the timeoutage limit",
				errors.ErrEmptySearchResult,
				&errdetails.RequestInfo{
					RequestId:   cfg.GetRequestId(),
					ServingData: errdetails.Serialize(cfg),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				}, info.Get(),
			)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeDeadlineExceeded(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if 0 < min && len(res.GetResults()) < min {
			err = status.WrapWithDeadlineExceeded(
				fmt.Sprintf("error search result length is not enough due to the timeoutage limit, required: %d, found: %d", min, len(res.GetResults())),
				errors.ErrInsuffcientSearchResult,
				&errdetails.RequestInfo{
					RequestId:   cfg.GetRequestId(),
					ServingData: errdetails.Serialize(cfg),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				}, info.Get(),
			)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeDeadlineExceeded(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
	}

	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse search gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   cfg.GetRequestId(),
				ServingData: errdetails.Serialize(cfg),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Warn(err)
		if len(res.GetResults()) == 0 {
			return nil, err
		}
	}
	if num != 0 && len(res.GetResults()) == 0 {
		if err == nil {
			err = errors.ErrEmptySearchResult
		}
		err = status.WrapWithNotFound("error search result length is 0", err,
			&errdetails.RequestInfo{
				RequestId:   cfg.GetRequestId(),
				ServingData: errdetails.Serialize(cfg),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	if 0 < min && len(res.GetResults()) < min {
		if err == nil {
			err = errors.ErrInsuffcientSearchResult
		}
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		err = status.WrapWithNotFound(
			fmt.Sprintf("error search result length is not enough required: %d, found: %d", min, len(res.GetResults())),
			errors.ErrInsuffcientSearchResult,
			&errdetails.RequestInfo{
				RequestId:   cfg.GetRequestId(),
				ServingData: errdetails.Serialize(cfg),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get(),
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	res.RequestId = cfg.GetRequestId()
	return res, nil
}

// vald standard algorithm
type valdStdAggr struct {
	num     int
	wg      sync.WaitGroup
	dch     chan DistPayload
	closed  atomic.Bool
	maxDist atomic.Value
	visited sync.Map
	result  []*payload.Object_Distance
	cancel  context.CancelFunc
}

func newStd(num, replica int) Aggregator {
	vsa := &valdStdAggr{
		num: num,
		dch: make(chan DistPayload, num*replica),
		maxDist: func() (av atomic.Value) {
			av.Store(big.NewFloat(math.MaxFloat64))
			return av
		}(),
		result: make([]*payload.Object_Distance, 0, num*replica),
	}
	vsa.closed.Store(false)
	return vsa
}

func (v *valdStdAggr) Start(ctx context.Context) {
	ctx, v.cancel = context.WithCancel(ctx)
	add := func(distance *big.Float, dist *payload.Object_Distance) {
		rl := len(v.result) // result length
		fmax, ok := v.maxDist.Load().(*big.Float)
		if !ok {
			return
		}
		if rl >= v.num && distance.Cmp(fmax) >= 0 {
			return
		}
		switch rl {
		case 0:
			v.result = append(v.result, dist)
		case 1:

			if distance.Cmp(big.NewFloat(float64(v.result[0].GetDistance()))) >= 0 {
				v.result = append(v.result, dist)
			} else {
				v.result = []*payload.Object_Distance{dist, v.result[0]}
			}
		default:
			pos := rl
			for idx := rl; idx >= 1; idx-- {
				if distance.Cmp(big.NewFloat(float64(v.result[idx-1].GetDistance()))) >= 0 {
					pos = idx - 1
					break
				}
			}
			switch {
			case pos == rl:
				v.result = append([]*payload.Object_Distance{dist}, v.result...)
			case pos == rl-1:
				v.result = append(v.result, dist)
			case pos >= 0:
				// skipcq: CRT-D0001
				v.result = append(v.result[:pos+1], v.result[pos:]...)
				v.result[pos+1] = dist
			}
		}
		rl = len(v.result)
		if rl > v.num && v.num != 0 {
			v.result = v.result[:v.num]
			rl = len(v.result)
		}
		if distEnd := big.NewFloat(float64(v.result[rl-1].GetDistance())); rl >= v.num &&
			distEnd.Cmp(fmax) < 0 {
			v.maxDist.Store(distEnd)
		}
	}
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()
		for {
			select {
			case <-ctx.Done():
				v.closed.Store(true)
				close(v.dch)
				for dist := range v.dch {
					add(dist.distance, dist.raw)
				}
				return
			case dist := <-v.dch:
				add(dist.distance, dist.raw)
			}
		}
	}()
}

func (v *valdStdAggr) Send(ctx context.Context, data *payload.Search_Response) {
	result := data.GetResults()
	if len(result) > v.num {
		result = result[:v.num]
	}
	for _, dist := range result {
		if dist != nil {
			fdist := big.NewFloat(float64(dist.GetDistance()))
			bf, ok := v.maxDist.Load().(*big.Float)
			if !ok || fdist.Cmp(bf) >= 0 {
				return
			}
			if _, already := v.visited.LoadOrStore(dist.GetId(), struct{}{}); !already {
				if v.closed.Load() {
					return
				}
				select {
				case <-ctx.Done():
					return
				case v.dch <- DistPayload{raw: dist, distance: fdist}:
				}
			}
		}
	}
}

func (v *valdStdAggr) Result() *payload.Search_Response {
	v.cancel()
	v.wg.Wait()
	if len(v.result) > v.num {
		v.result = v.result[:v.num]
	}
	return &payload.Search_Response{
		Results: v.result,
	}
}

// pairing heap
type valdPairingHeapAggr struct {
	num     int
	ph      *PairingHeap
	mu      sync.Mutex
	visited sync.Map
	result  []*payload.Object_Distance
}

func newPairingHeap(num, replica int) Aggregator {
	return &valdPairingHeapAggr{
		num:    num,
		ph:     new(PairingHeap),
		result: make([]*payload.Object_Distance, 0, num),
	}
}

func (v *valdPairingHeapAggr) Start(_ context.Context) {}

func (v *valdPairingHeapAggr) Send(ctx context.Context, data *payload.Search_Response) {
	result := data.GetResults()
	if len(result) > v.num {
		result = result[:v.num]
	}
	for _, dist := range result {
		if dist != nil {
			if _, already := v.visited.LoadOrStore(dist.GetId(), struct{}{}); !already {
				select {
				case <-ctx.Done():
					return
				default:
					dp := &DistPayload{raw: dist, distance: big.NewFloat(float64(dist.GetDistance()))}
					v.mu.Lock()
					v.ph = v.ph.Insert(dp)
					v.mu.Unlock()
				}
			}
		}
	}
}

func (v *valdPairingHeapAggr) Result() *payload.Search_Response {
	for !v.ph.IsEmpty() && len(v.result) <= v.num {
		var min *DistPayload
		min, v.ph = v.ph.ExtractMin()
		v.result = append(v.result, min.raw)
	}
	if len(v.result) > v.num {
		v.result = v.result[:v.num]
	}
	return &payload.Search_Response{
		Results: v.result,
	}
}

// plane sort
type valdSliceAggr struct {
	num    int
	mu     sync.Mutex
	result []*DistPayload
}

func newSlice(num, replica int) Aggregator {
	return &valdSliceAggr{
		num:    num,
		result: make([]*DistPayload, 0, num*replica),
	}
}

func (_ *valdSliceAggr) Start(_ context.Context) {}

func (v *valdSliceAggr) Send(ctx context.Context, data *payload.Search_Response) {
	result := data.GetResults()
	if len(result) > v.num {
		result = result[:v.num]
	}
	for _, dist := range result {
		if dist != nil {
			select {
			case <-ctx.Done():
				return
			default:
				dp := &DistPayload{raw: dist, distance: big.NewFloat(float64(dist.GetDistance()))}
				v.mu.Lock()
				v.result = append(v.result, dp)
				v.mu.Unlock()
			}
		}
	}
}

func (v *valdSliceAggr) Result() (res *payload.Search_Response) {
	slices.RemoveDuplicates(v.result, func(l, r *DistPayload) bool {
		return l.distance.Cmp(r.distance) < 0
	})

	if len(v.result) > v.num {
		v.result = v.result[:v.num]
	}
	res = &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, v.num),
	}
	for _, r := range v.result {
		res.Results = append(res.GetResults(), r.raw)
	}
	return res
}

// plane sort
type valdPoolSliceAggr struct {
	num    int
	mu     sync.Mutex
	result []*DistPayload
}

var (
	poolDist = sync.Pool{
		New: func() interface{} {
			return make([]*DistPayload, 0, poolLen.Load())
		},
	}
	poolLen atomic.Uint64
)

func newPoolSlice(num, replica int) Aggregator {
	l := uint64(num * replica)
	if poolLen.Load() < l {
		poolLen.Store(l)
	}
	return &valdPoolSliceAggr{
		num:    num,
		result: poolDist.Get().([]*DistPayload),
	}
}

func (_ *valdPoolSliceAggr) Start(_ context.Context) {}

func (v *valdPoolSliceAggr) Send(ctx context.Context, data *payload.Search_Response) {
	result := data.GetResults()
	if len(result) > v.num {
		result = result[:v.num]
	}
	for _, dist := range result {
		if dist != nil {
			select {
			case <-ctx.Done():
				return
			default:
				dp := &DistPayload{raw: dist, distance: big.NewFloat(float64(dist.GetDistance()))}
				v.mu.Lock()
				v.result = append(v.result, dp)
				v.mu.Unlock()
			}
		}
	}
}

func (v *valdPoolSliceAggr) Result() (res *payload.Search_Response) {
	slices.RemoveDuplicates(v.result, func(l, r *DistPayload) bool {
		return l.distance.Cmp(r.distance) < 0
	})

	if len(v.result) > v.num {
		v.result = v.result[:v.num]
	}
	res = &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, v.num),
	}
	for _, r := range v.result {
		res.Results = append(res.GetResults(), r.raw)
	}
	poolDist.Put(v.result[:0])
	return res
}
