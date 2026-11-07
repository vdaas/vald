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
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func (s *server) Search(
	ctx context.Context, req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.SearchRPCName), apiName+"/"+vald.SearchRPCName)
	defer trace.End(span)
	vl := len(req.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument(vald.SearchRPCName+" API invalid vector argument", err,
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
	res, attrs, err := s.doSearch(
		ctx,
		vald.SearchRPCName,
		req.GetConfig(),
		func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.Search(ctx, &payload.Search_Request{Vector: req.GetVector(), Config: fcfg}, copts...)
		},
	)
	if err != nil {
		if attrs != nil && span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) SearchByID(
	ctx context.Context, req *payload.Search_IDRequest,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.SearchByIDRPCName), apiName+"/"+vald.SearchByIDRPCName)
	defer trace.End(span)
	uuid := req.GetId()
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.SearchByIDRPCName+" API invalid uuid", err, append(s.errDetails(ctx, vald.SearchByIDRPCName, uuid, req),
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			})...)
		return errhandler.HandleError[payload.Search_Response](span, codes.InvalidArgument, err)
	}
	vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: uuid,
		},
	})
	if err != nil {
		st, _ := status.FromError(err)
		if span != nil && st != nil && st.Code() != codes.NotFound {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		if skipSearchByIDFallback(ctx, err) {
			return nil, err
		}
		// try search by using agent's SearchByID method this operation is emergency fallback, the search quality is not same as usual SearchByID operation.
		var attrs []attribute.KeyValue
		res, attrs, err = s.doSearch(
			ctx,
			vald.SearchByIDRPCName,
			req.GetConfig(),
			func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
				return vc.SearchByID(ctx, &payload.Search_IDRequest{Id: req.GetId(), Config: fcfg}, copts...)
			},
		)
		if err == nil {
			return res, nil
		}
		if st != nil && span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	res, err = s.Search(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: req.GetConfig(),
	})
	if err != nil {
		if skipSearchByIDFallback(ctx, err) {
			return nil, err
		}
		var attrs []attribute.KeyValue
		res, attrs, err = s.doSearch(
			ctx,
			vald.SearchByIDRPCName,
			req.GetConfig(),
			func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
				return vc.SearchByID(ctx, &payload.Search_IDRequest{Id: req.GetId(), Config: fcfg}, copts...)
			},
		)
		if err != nil {
			errhandler.RecordSpanAttrs(span, attrs, err)
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
	if ctx.Err() != nil {
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

// doSearch broadcasts the search to every agent and aggregates the results. rpc
// is the agent RPC name f calls; it labels the partial-failure metrics.
func (s *server) doSearch(
	ctx context.Context,
	rpc string,
	cfg *payload.Search_Config,
	f func(ctx context.Context, cfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error),
) (res *payload.Search_Response, attrs []attribute.KeyValue, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doSearch"), apiName+"/doSearch")
	defer trace.End(span)

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
		return nil, nil, err
	}

	var (
		num     = int(cfg.GetNum())
		fnum    int
		replica = s.gateway.GetAgentCount(ctx)
	)

	if cfg.GetRatio() != nil {
		fnum = int(s.calculateNum(cfg.GetNum(), cfg.GetRatio().GetValue(), replica))
	}
	if fnum <= 0 {
		fnum = num
	}

	return s.aggregationSearch(ctx, rpc, replica, selectAggregator(cfg.GetAggregationAlgorithm(), num, fnum, replica), cfg, f)
}

func selectAggregator(algo payload.Search_AggregationAlgorithm, num, fnum, replica int) Aggregator {
	switch algo {
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

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamSearchRPCName), apiName+"/"+vald.StreamSearchRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_Request) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer trace.End(sspan)
			res, err := s.Search(ctx, req)
			if err != nil {
				st, _ := status.FromError(err)
				errhandler.RecordSpanStatus(sspan, st, err)
				return &payload.Search_StreamResponse{
					Payload: &payload.Search_StreamResponse_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Search_StreamResponse{
				Payload: &payload.Search_StreamResponse_Response{
					Response: res,
				},
			}, nil
		})
	if err != nil {
		st, _ := status.FromError(err)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamSearchByIDRPCName), apiName+"/"+vald.StreamSearchByIDRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamSearchByIDRPCName+"/id-"+req.GetId())
			defer trace.End(sspan)
			res, err := s.SearchByID(ctx, req)
			if err != nil {
				st, _ := status.FromError(err)
				errhandler.RecordSpanStatus(sspan, st, err)
				return &payload.Search_StreamResponse{
					Payload: &payload.Search_StreamResponse_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Search_StreamResponse{
				Payload: &payload.Search_StreamResponse_Response{
					Response: res,
				},
			}, nil
		})
	if err != nil {
		st, _ := status.FromError(err)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

func (s *server) MultiSearch(
	ctx context.Context, reqs *payload.Search_MultiRequest,
) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiSearchRPCName), apiName+"/"+vald.MultiSearchRPCName)
	defer trace.End(span)
	res, errs = multiSearch(ctx, s, vald.MultiSearchRPCName, reqs.GetRequests(), s.Search)
	if errs != nil {
		st, _ := status.FromError(errs)
		errhandler.RecordSpanStatus(span, st, errs)
	}
	return res, errs
}

// multiSearch is the shared fan-out of MultiSearch, MultiSearchByID,
// MultiLinearSearch and MultiLinearSearchByID: it runs search for each request
// under the multi-request concurrency limit, records each failure on a per-item
// child span of rpcName, and returns the positional responses together with the
// joined errors. Each goroutine writes its own index of the pre-sized slice and
// eg.Wait orders those writes before the return, so only errList needs a lock.
func multiSearch[Q any](
	ctx context.Context,
	s *server,
	rpcName string,
	reqs []Q,
	search func(context.Context, Q) (*payload.Search_Response, error),
) (*payload.Search_Responses, error) {
	res := &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs)),
	}
	eg := errgroup.NewBoundedFanOut(ctx, s.multiConcurrency)
	var emu sync.Mutex
	var errList []error
	for i, req := range reqs {
		idx, query := i, req
		eg.Go(safety.RecoverFunc(func() error {
			ti := "errgroup.Go/id-" + strconv.Itoa(idx)
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+rpcName+"/"+ti)
			defer trace.End(sspan)
			r, err := search(sctx, query)
			if err != nil {
				st, _ := status.FromError(err)
				errhandler.RecordSpanStatus(sspan, st, err)
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
	return res, errors.Join(errList...)
}

func (s *server) MultiSearchByID(
	ctx context.Context, reqs *payload.Search_MultiIDRequest,
) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiSearchByIDRPCName), apiName+"/"+vald.MultiSearchByIDRPCName)
	defer trace.End(span)
	res, errs = multiSearch(ctx, s, vald.MultiSearchByIDRPCName, reqs.GetRequests(), s.SearchByID)
	if errs != nil {
		st, _ := status.FromError(errs)
		errhandler.RecordSpanStatus(span, st, errs)
	}
	return res, errs
}
