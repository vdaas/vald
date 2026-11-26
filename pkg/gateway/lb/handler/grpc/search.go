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
	"math"
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func (s *server) Search(
	ctx context.Context, req *payload.Search_Request,
) (*payload.Search_Response, error) {
	return s.searchVector(ctx, vald.SearchRPCName, req, vald.Client.Search)
}

func (s *server) SearchByID(
	ctx context.Context, req *payload.Search_IDRequest,
) (*payload.Search_Response, error) {
	return s.searchByID(ctx, vald.SearchByIDRPCName, vald.SearchRPCName, req, vald.Client.SearchByID, vald.Client.Search)
}

func (s *server) searchVector(
	ctx context.Context,
	rpcName string,
	req *payload.Search_Request,
	fn func(vc vald.Client, ctx context.Context, req *payload.Search_Request, copts ...grpc.CallOption) (*payload.Search_Response, error),
) (*payload.Search_Response, error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+rpcName), apiName+"/"+rpcName)
		defer trace.End(span)
	}
	vl := len(req.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err := errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument(rpcName+" API invalid vector argument", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			})
		return errhandler.HandleError[payload.Search_Response](span, codes.InvalidArgument, err)
	}
	cfg := req.GetConfig()
	res, attrs, err := doSearch(ctx, s, rpcName, cfg,
		func(fcfg *payload.Search_Config) *payload.Search_Request {
			return &payload.Search_Request{Vector: req.GetVector(), Config: fcfg}
		}, fn)
	if err != nil {
		if span != nil {
			errhandler.RecordSpanAttrs(span, attrs, err)
		}
		return nil, err
	}
	return res, nil
}

func (s *server) searchByID(
	ctx context.Context,
	idRPCName, vecRPCName string,
	req *payload.Search_IDRequest,
	idFn func(vc vald.Client, ctx context.Context, req *payload.Search_IDRequest, copts ...grpc.CallOption) (*payload.Search_Response, error),
	vecFn func(vc vald.Client, ctx context.Context, req *payload.Search_Request, copts ...grpc.CallOption) (*payload.Search_Response, error),
) (*payload.Search_Response, error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+idRPCName), apiName+"/"+idRPCName)
		defer trace.End(span)
	}
	uuid := req.GetId()
	if len(uuid) == 0 {
		return errhandler.HandleError[payload.Search_Response](span, codes.InvalidArgument, s.invalidUUIDError(ctx, idRPCName, idRPCName, uuid, req))
	}
	buildIDReq := func(fcfg *payload.Search_Config) *payload.Search_IDRequest {
		return &payload.Search_IDRequest{Id: uuid, Config: fcfg}
	}
	vec, err := s.getObject(ctx, uuid)
	if err != nil {
		if skipSearchByIDFallback(ctx, err) {
			return nil, s.wrapObjectReadError(ctx, span, idRPCName, uuid, req, "object not found", err)
		}
		if span != nil {
			if st, _ := status.FromError(err); st != nil && st.Code() != codes.NotFound {
				errhandler.RecordSpanStatus(span, st, err)
			}
		}
		// try search by using agent's SearchByID method this operation is emergency fallback, the search quality is not same as usual SearchByID operation.
		var (
			res   *payload.Search_Response
			attrs []attribute.KeyValue
		)
		res, attrs, err = doSearch(ctx, s, idRPCName, req.GetConfig(), buildIDReq, idFn)
		if err == nil {
			return res, nil
		}
		if span != nil {
			errhandler.RecordSpanAttrs(span, attrs, err)
		}
		return nil, err
	}
	var (
		res   *payload.Search_Response
		attrs []attribute.KeyValue
	)
	if vl := len(vec.GetVector()); vl < algorithm.MinimumVectorDimensionSize {
		err = status.WrapWithInvalidArgument(vecRPCName+" API invalid vector argument", errors.ErrInvalidDimensionSize(vl, 0),
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			})
	} else {
		res, attrs, err = doSearch(ctx, s, vecRPCName, req.GetConfig(),
			func(fcfg *payload.Search_Config) *payload.Search_Request {
				return &payload.Search_Request{Vector: vec.GetVector(), Config: fcfg}
			}, vecFn)
	}
	if err != nil {
		if skipSearchByIDFallback(ctx, err) {
			if span != nil {
				errhandler.RecordSpanAttrs(span, attrs, err)
			}
			return nil, err
		}
		res, attrs, err = doSearch(ctx, s, idRPCName, req.GetConfig(), buildIDReq, idFn)
		if err != nil {
			if span != nil {
				errhandler.RecordSpanAttrs(span, attrs, err)
			}
			return nil, err
		}
	}
	return res, nil
}

// skipSearchByIDFallback reports whether SearchByID must return err instead of
// falling back to the agent-side SearchByID broadcast: the request ctx is done,
// or the failed call already used up the per-request search timeout, which the
// fallback broadcast would only wait out a second time.
func skipSearchByIDFallback(ctx context.Context, err error) bool {
	if ctx.Err() != nil || errors.IsAny(err, context.Canceled, context.DeadlineExceeded) {
		return true
	}
	st, ok := status.FromError(err)
	return ok && st != nil && (st.Code() == codes.DeadlineExceeded || st.Code() == codes.Canceled)
}

// calculateNum adjusts the number of search results based on the ratio and the number of replicas.
// It ensures that the number of results is not less than the minimum required and adjusts based on the provided ratio.
func (s *server) calculateNum(num uint32, ratio float32, agents int) (n uint32) {
	minRatio := float64(s.replica) / float64(agents)
	if ratio <= 0.0 {
		return uint32(math.Ceil(float64(num) * minRatio))
	}
	n = uint32(math.Ceil(float64(num) * (minRatio + ((1 - minRatio) * float64(ratio)))))
	sn := uint32(math.Ceil(float64(num) * minRatio))
	if n == 0 || n-1 < sn {
		return sn
	}
	return n - 1
}

func (s *server) prepareSearch(
	ctx context.Context, span trace.Span, cfg *payload.Search_Config,
) (replica int, aggr Aggregator, err error) {
	if cfg == nil {
		err = errors.ErrInvalidSearchConfig("search config is nil in doSearch")
		err = status.WrapWithInvalidArgument(apiName+"/doSearch", err,
			&errdetails.RequestInfo{
				RequestId:   "unknown request ID",
				ServingData: "Search_Config is nil",
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "Search_Config is nil",
						Description: err.Error(),
					},
				},
			})
		errhandler.RecordSpanError(span, codes.InvalidArgument, err)
		return 0, nil, err
	}

	var (
		num  = int(cfg.GetNum())
		fnum int
	)
	replica = s.gateway.GetAgentCount(ctx)

	if cfg.GetRatio() != nil {
		fnum = int(s.calculateNum(cfg.GetNum(), cfg.GetRatio().GetValue(), replica))
	}
	if fnum <= 0 {
		fnum = num
	}

	return replica, selectAggregator(cfg.GetAggregationAlgorithm(), num, fnum, replica), nil
}

// doSearch prepares the aggregator, constructs the forwarding request once, and
// broadcasts it to every agent before aggregating the results.
func doSearch[Q any](
	ctx context.Context,
	s *server,
	rpc string,
	cfg *payload.Search_Config,
	buildReq func(fcfg *payload.Search_Config) Q,
	f func(vc vald.Client, ctx context.Context, req Q, copts ...grpc.CallOption) (*payload.Search_Response, error),
) (res *payload.Search_Response, attrs []attribute.KeyValue, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doSearch"), apiName+"/doSearch")
		defer trace.End(span)
	}

	replica, aggr, err := s.prepareSearch(ctx, span, cfg)
	if err != nil {
		return nil, nil, err
	}

	sreq := buildReq(forwardSearchConfig(cfg, aggr.GetFnum()))
	return s.aggregationSearch(ctx, rpc, replica, aggr, cfg,
		func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return f(vc, ctx, sreq, copts...)
		},
	)
}

func selectAggregator(algo payload.Search_AggregationAlgorithm, num, fnum, replica int) Aggregator {
	switch algo {
	case payload.Search_Unknown, payload.Search_ConcurrentQueue:
		return newStd(num, fnum, replica)
	case payload.Search_SortSlice:
		return newSlice(num, fnum, replica)
	case payload.Search_SortPoolSlice:
		return newPoolSlice(num, fnum, replica)
	case payload.Search_PairingHeap:
		return newPairingHeap(num, fnum, replica)
	default:
		return newStd(num, fnum, replica)
	}
}

func streamSearch[Q proto.Message, S grpc.TypedServerStream[Q, *payload.Search_StreamResponse]](
	s *server,
	stream S,
	rpcName string,
	itemSpan func(Q) string,
	search func(context.Context, Q) (*payload.Search_Response, error),
) error {
	return grpc.HandleBidirectionalStream(
		stream,
		s.streamConcurrency,
		vald.PackageName+"."+vald.SearchRPCServiceName+"/"+rpcName,
		apiName+"/"+rpcName,
		itemSpan,
		search,
		func(res *payload.Search_Response) *payload.Search_StreamResponse {
			return &payload.Search_StreamResponse{
				Payload: &payload.Search_StreamResponse_Response{Response: res},
			}
		},
		func(st *status.Status) *payload.Search_StreamResponse {
			return &payload.Search_StreamResponse{
				Payload: &payload.Search_StreamResponse_Status{Status: st.Proto()},
			}
		},
	)
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) error {
	return streamSearch(s, stream, vald.StreamSearchRPCName, func(req *payload.Search_Request) string {
		return "requestID-" + req.GetConfig().GetRequestId()
	}, s.Search)
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) error {
	return streamSearch(s, stream, vald.StreamSearchByIDRPCName, func(req *payload.Search_IDRequest) string {
		return "id-" + req.GetId()
	}, s.SearchByID)
}

func (s *server) MultiSearch(
	ctx context.Context, reqs *payload.Search_MultiRequest,
) (*payload.Search_Responses, error) {
	return multiSearch(ctx, s, vald.MultiSearchRPCName, reqs.GetRequests(), s.Search)
}

// multiSearch is the shared fan-out of MultiSearch, MultiSearchByID,
// MultiLinearSearch and MultiLinearSearchByID: it opens the root trace span,
// runs search for each request under the multi-request concurrency limit,
// records each failure on a per-item child span of rpcName, and returns the
// positional responses together with the joined errors. Each goroutine writes
// its own index of the pre-sized slice and eg.Wait orders those writes before
// the return, so only errList needs a lock.
func multiSearch[Q any](
	ctx context.Context,
	s *server,
	rpcName string,
	reqs []Q,
	search func(context.Context, Q) (*payload.Search_Response, error),
) (*payload.Search_Responses, error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+rpcName), apiName+"/"+rpcName)
		defer trace.End(span)
	}
	res := &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs)),
	}
	if len(reqs) == 0 {
		return res, nil
	}
	if len(reqs) == 1 {
		sctx := ctx
		var sspan trace.Span
		if trace.Enabled() {
			const ti = "errgroup.Go/id-0"
			sctx, sspan = trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+rpcName+"/"+ti)
			defer trace.End(sspan)
		}
		r, err := search(sctx, reqs[0])
		if err != nil {
			if sspan != nil {
				errhandler.RecordSpanFromError(sspan, err)
			}
			if span != nil {
				errhandler.RecordSpanFromError(span, err)
			}
			return res, err
		}
		res.Responses[0] = r
		return res, nil
	}
	eg := errgroup.NewBoundedFanOut(ctx, s.multiLimit(ctx, true))
	var emu sync.Mutex
	var errList []error
	for i, req := range reqs {
		idx, query := i, req
		eg.Go(safety.RecoverFunc(func() error {
			sctx := ctx
			var sspan trace.Span
			if trace.Enabled() {
				ti := "errgroup.Go/id-" + strconv.Itoa(idx)
				sctx, sspan = trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+rpcName+"/"+ti)
				defer trace.End(sspan)
			}
			r, err := search(sctx, query)
			if err != nil {
				if sspan != nil {
					errhandler.RecordSpanFromError(sspan, err)
				}
				emu.Lock()
				errList = append(errList, err)
				emu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		}))
	}
	if err := eg.Wait(); err != nil {
		errList = append(errList, err)
	}
	errs := errors.Join(errList...)
	if errs != nil && span != nil {
		errhandler.RecordSpanFromError(span, errs)
	}
	return res, errs
}

func (s *server) MultiSearchByID(
	ctx context.Context, reqs *payload.Search_MultiIDRequest,
) (*payload.Search_Responses, error) {
	return multiSearch(ctx, s, vald.MultiSearchByIDRPCName, reqs.GetRequests(), s.SearchByID)
}
