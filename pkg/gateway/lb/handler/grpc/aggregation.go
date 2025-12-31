// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"slices"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type Aggregator interface {
	Start(ctx context.Context)
	Send(ctx context.Context, data *payload.Search_Response)
	GetNum() int  // get top-k number
	GetFnum() int // get forwarding top-k number calculated by search ratio
	Result() *payload.Search_Response
}

type DistPayload struct {
	raw      *payload.Object_Distance
	distance *big.Float
}

func (s *server) aggregationSearch(
	ctx context.Context,
	aggr Aggregator,
	bcfg *payload.Search_Config, // Base Config of Request
	f func(ctx context.Context,
		fcfg *payload.Search_Config, // Forwarding Config to Agent
		vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error),
) (res *payload.Search_Response, attrs []attribute.KeyValue, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "aggregationSearch"), apiName+"/aggregationSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if bcfg == nil {
		return nil, nil, errors.ErrInvalidSearchConfig("bcfg is nil in aggregationSearch")
	}

	num := aggr.GetNum()
	minNum := int(bcfg.GetMinNum())

	var timeout time.Duration
	if to := bcfg.GetTimeout(); to != 0 {
		timeout = time.Duration(to)
	} else {
		timeout = s.timeout
	}

	fcfg := bcfg.CloneVT() // Forwarding Config to Agent, this config need to modify like below so it should be cloned
	fcfg.Num = uint32(aggr.GetFnum())
	fcfg.MinNum = 0

	ctx, cancel := context.WithTimeout(ctx, timeout)
	aggr.Start(ctx)
	err = s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/aggregationSearch/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		r, err := f(sctx, fcfg, vc, copts...)
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
				return nil
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
				return nil
			default:
				st, ok := status.FromError(err)
				if !ok {
					log.Debug(err)
					return nil
				}
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
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
				log.Debug(err)
			}
			return nil
		}
		if r == nil || len(r.GetResults()) == 0 {
			select {
			case <-sctx.Done():
				err = status.WrapWithNotFound(fmt.Sprintf("failed to process search request from %s", target),
					errors.ErrEmptySearchResult,
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
						ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
					})
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				return nil
			default:
				r, err = f(sctx, fcfg, vc, copts...)
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
						return nil
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
						return nil
					default:
						st, ok := status.FromError(err)
						if !ok {
							log.Debug(err)
							return nil
						}
						if sspan != nil {
							sspan.RecordError(err)
							sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
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
						log.Debug(err)
					}
					return nil
				}
				if r == nil || len(r.GetResults()) == 0 {
					err = status.WrapWithNotFound(fmt.Sprintf("failed to process search request from %s", target),
						errors.ErrEmptySearchResult,
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
					return nil
				}
			}
		}
		aggr.Send(sctx, r)
		return nil
	})
	cancel()
	if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
		err = status.WrapWithInternal("search API connection not found", err,
			&errdetails.RequestInfo{
				RequestId:   bcfg.GetRequestId(),
				ServingData: errdetails.Serialize(bcfg),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		attrs = trace.StatusCodeInternal(err.Error())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, attrs, err
	}
	res = aggr.Result()
	if num != 0 && len(res.GetResults()) > num {
		res.Results = res.GetResults()[:num]
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		if len(res.GetResults()) == 0 {
			err = status.WrapWithDeadlineExceeded(
				"error search result length is 0 due to the timeoutage limit",
				errors.ErrEmptySearchResult,
				&errdetails.RequestInfo{
					RequestId:   bcfg.GetRequestId(),
					ServingData: errdetails.Serialize(bcfg),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				}, info.Get(),
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, attrs, err
		}
		if 0 < minNum && len(res.GetResults()) < minNum {
			err = status.WrapWithDeadlineExceeded(
				fmt.Sprintf("error search result length is not enough due to the timeoutage limit, required: %d, found: %d", minNum, len(res.GetResults())),
				errors.ErrInsuffcientSearchResult,
				&errdetails.RequestInfo{
					RequestId:   bcfg.GetRequestId(),
					ServingData: errdetails.Serialize(bcfg),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				}, info.Get(),
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, attrs, err
		}
	}

	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse search gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   bcfg.GetRequestId(),
				ServingData: errdetails.Serialize(bcfg),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		attrs = trace.FromGRPCStatus(st.Code(), msg)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Warn(err)
		if len(res.GetResults()) == 0 {
			return nil, attrs, err
		}
	}
	if num != 0 && len(res.GetResults()) == 0 {
		if err == nil {
			err = errors.ErrEmptySearchResult
		}
		err = status.WrapWithNotFound("error search result length is 0", err,
			&errdetails.RequestInfo{
				RequestId:   bcfg.GetRequestId(),
				ServingData: errdetails.Serialize(bcfg),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		attrs = trace.StatusCodeNotFound(err.Error())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, attrs, err
	}

	if 0 < minNum && len(res.GetResults()) < minNum {
		if err == nil {
			err = errors.ErrInsuffcientSearchResult
		}
		attrs = trace.StatusCodeNotFound(err.Error())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		err = status.WrapWithNotFound(
			fmt.Sprintf("error search result length is not enough required: %d, found: %d", minNum, len(res.GetResults())),
			errors.ErrInsuffcientSearchResult,
			&errdetails.RequestInfo{
				RequestId:   bcfg.GetRequestId(),
				ServingData: errdetails.Serialize(bcfg),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get(),
		)
		attrs = trace.StatusCodeNotFound(err.Error())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, attrs, err
	}
	res.RequestId = bcfg.GetRequestId()
	return res, attrs, nil
}

// vald standard algorithm.
type valdStdAggr struct {
	num     int // top-k number
	fnum    int // forward top-k number
	wg      sync.WaitGroup
	dch     chan DistPayload
	closed  atomic.Bool
	maxDist atomic.Value
	visited sync.Map[string, any]
	result  []*payload.Object_Distance
	cancel  context.CancelFunc
	eg      errgroup.Group
}

// newStd returns a valdStdAggr configured for the standard aggregation algorithm.
// It sets the top-k size (`num`), the forwarding top-k size (`fnum`), and allocates internal buffers sized proportionally to `replica`.
// The returned Aggregator is initialized (including its error group) and ready to Start and receive results.
func newStd(num, fnum, replica int) Aggregator {
	vsa := &valdStdAggr{
		num:  num,
		fnum: fnum,
		dch:  make(chan DistPayload, num*replica),
		maxDist: func() (av atomic.Value) {
			av.Store(big.NewFloat(math.MaxFloat64))
			return av
		}(),
		result: make([]*payload.Object_Distance, 0, num*replica),
		eg:     errgroup.Get(),
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
	v.eg.Go(func() error {
		defer v.wg.Done()
		for {
			select {
			case <-ctx.Done():
				v.closed.Store(true)
				close(v.dch)
				for dist := range v.dch {
					add(dist.distance, dist.raw)
				}
				return nil
			case dist := <-v.dch:
				add(dist.distance, dist.raw)
			}
		}
	})
}

func (v *valdStdAggr) Send(ctx context.Context, data *payload.Search_Response) {
	result := data.GetResults()
	if len(result) > v.fnum {
		result = result[:v.fnum]
	}
	for _, dist := range result {
		if dist != nil {
			fdist := big.NewFloat(float64(dist.GetDistance()))
			bf, ok := v.maxDist.Load().(*big.Float)
			if !ok || fdist.Cmp(bf) >= 0 || v.closed.Load() {
				return
			}
			if _, already := v.visited.LoadOrStore(dist.GetId(), struct{}{}); !already {
				if v.closed.Load() {
					return
				}
				select {
				case <-ctx.Done():
					v.visited.Delete(dist.GetId())
					return
				case v.dch <- DistPayload{raw: dist, distance: fdist}:
				}
			}
		}
	}
}

func (v *valdStdAggr) Result() *payload.Search_Response {
	v.closed.Store(true)
	v.cancel()
	v.wg.Wait()
	if len(v.result) > v.num {
		v.result = v.result[:v.num]
	}
	return &payload.Search_Response{
		Results: v.result,
	}
}

func (v *valdStdAggr) GetNum() int {
	if v != nil {
		return v.num
	}
	return 0
}

func (v *valdStdAggr) GetFnum() int {
	if v != nil {
		return v.fnum
	}
	return 0
}

// pairing heap.
type valdPairingHeapAggr struct {
	num     int // top-k number
	fnum    int // forward top-k number
	ph      *PairingHeap
	mu      sync.Mutex
	visited sync.Map[string, any]
	result  []*payload.Object_Distance
}

func newPairingHeap(num, fnum, _ int) Aggregator {
	return &valdPairingHeapAggr{
		num:    num,
		fnum:   fnum,
		ph:     new(PairingHeap),
		result: make([]*payload.Object_Distance, 0, num),
	}
}

func (*valdPairingHeapAggr) Start(_ context.Context) {}

func (v *valdPairingHeapAggr) Send(ctx context.Context, data *payload.Search_Response) {
	result := data.GetResults()
	if len(result) > v.fnum {
		result = result[:v.fnum]
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
		var minDist *DistPayload
		minDist, v.ph = v.ph.ExtractMin()
		if minDist != nil {
			v.result = append(v.result, minDist.raw)
		} else if v.ph == nil {
			break
		}
	}
	if len(v.result) > v.num {
		v.result = v.result[:v.num]
	}
	return &payload.Search_Response{
		Results: v.result,
	}
}

func (v *valdPairingHeapAggr) GetNum() int {
	if v != nil {
		return v.num
	}
	return 0
}

func (v *valdPairingHeapAggr) GetFnum() int {
	if v != nil {
		return v.fnum
	}
	return 0
}

// plane sort.
type valdSliceAggr struct {
	num    int // top-k number
	fnum   int // forward top-k number
	mu     sync.Mutex
	result []*DistPayload
}

func newSlice(num, fnum, replica int) Aggregator {
	return &valdSliceAggr{
		num:    num,
		fnum:   fnum,
		result: make([]*DistPayload, 0, num*replica),
	}
}

func (*valdSliceAggr) Start(_ context.Context) {}

func (v *valdSliceAggr) Send(ctx context.Context, data *payload.Search_Response) {
	result := data.GetResults()
	if len(result) > v.fnum {
		result = result[:v.fnum]
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
	removeDuplicates(v.result, func(l, r *DistPayload) int {
		return l.distance.Cmp(r.distance)
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

func (v *valdSliceAggr) GetNum() int {
	if v != nil {
		return v.num
	}
	return 0
}

func (v *valdSliceAggr) GetFnum() int {
	if v != nil {
		return v.fnum
	}
	return 0
}

// plane sort.
type valdPoolSliceAggr struct {
	num    int // top-k number
	fnum   int // forward top-k number
	mu     sync.Mutex
	result []*DistPayload
}

var (
	poolDist = sync.Pool{
		New: func() any {
			return make([]*DistPayload, 0, poolLen.Load())
		},
	}
	poolLen atomic.Uint64
)

func newPoolSlice(num, fnum, replica int) Aggregator {
	l := uint64(num * replica)
	if poolLen.Load() < l {
		poolLen.Store(l)
	}
	return &valdPoolSliceAggr{
		num:    num,
		fnum:   fnum,
		result: poolDist.Get().([]*DistPayload),
	}
}

func (*valdPoolSliceAggr) Start(_ context.Context) {}

func (v *valdPoolSliceAggr) Send(ctx context.Context, data *payload.Search_Response) {
	result := data.GetResults()
	if len(result) > v.fnum {
		result = result[:v.fnum]
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
	removeDuplicates(v.result, func(l, r *DistPayload) int {
		return l.distance.Cmp(r.distance)
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
	// skipcq: SCC-SA6002
	poolDist.Put(v.result[:0])
	return res
}

func (v *valdPoolSliceAggr) GetNum() int {
	if v != nil {
		return v.num
	}
	return 0
}

func (v *valdPoolSliceAggr) GetFnum() int {
	if v != nil {
		return v.fnum
	}
	return 0
}

func removeDuplicates[S ~[]E, E comparable](x S, less func(left, right E) int) S {
	if len(x) < 2 {
		return x
	}
	slices.SortStableFunc(x, less)
	return slices.Compact(x)
}
