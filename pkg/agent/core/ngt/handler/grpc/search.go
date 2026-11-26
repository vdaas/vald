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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

func (s *server) wrapSearchContextError(
	err error, api, resType, reqID string, req any,
) ([]attribute.KeyValue, error) {
	reqInfo := &errdetails.RequestInfo{
		RequestId:   reqID,
		ServingData: errdetails.Serialize(req),
	}
	resInfo := s.resourceInfo(resType)
	if errors.Is(err, context.DeadlineExceeded) {
		err = status.WrapWithDeadlineExceeded(api+" API deadline exceeded to process search request", err, reqInfo, resInfo)
		log.Debug(err)
		return trace.StatusCodeDeadlineExceeded(err.Error()), err
	}
	err = status.WrapWithCanceled(api+" API canceled to process search request", err, reqInfo, resInfo)
	log.Debug(err)
	return trace.StatusCodeCancelled(err.Error()), err
}

func (s *server) handleSearchError(
	span trace.Span,
	rpcName string,
	uuid string,
	byID bool,
	reqCfg *payload.Search_Config,
	req any,
	vecLen int,
	err error,
) error {
	resType := ngtResourceType + "/ngt." + rpcName
	reqID := reqCfg.GetRequestId()
	var attrs []attribute.KeyValue
	// errors.As writes to its target, so it must be local to stay race-free across concurrent requests.
	var errNGT errors.NGTError
	switch {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		attrs, err = s.wrapSearchContextError(err, rpcName, resType, reqID, req)
	case errors.Is(err, errors.ErrCreateIndexingIsInProgress):
		err = status.WrapWithAborted(rpcName+" API aborted to process search request due to creating indices is in progress", err,
			&errdetails.RequestInfo{
				RequestId:   reqID,
				ServingData: errdetails.Serialize(req),
			},
			s.resourceInfo(resType))
		log.Debug(err)
		if span != nil {
			attrs = trace.StatusCodeAborted(err.Error())
		}
	case errors.Is(err, errors.ErrFlushingIsInProgress):
		err = status.WrapWithAborted(rpcName+" API aborted to process search request due to flushing indices is in progress", err,
			&errdetails.RequestInfo{
				RequestId:   reqID,
				ServingData: errdetails.Serialize(req),
			},
			s.resourceInfo(resType))
		log.Debug(err)
		if span != nil {
			attrs = trace.StatusCodeAborted(err.Error())
		}
	case errors.Is(err, errors.ErrEmptySearchResult):
		if byID {
			err = status.New(codes.NotFound, rpcName+" API uuid "+uuid+"'s search result not found").Err()
		} else {
			err = status.New(codes.NotFound, rpcName+" API requestID "+reqID+"'s search result not found").Err()
		}
		if span != nil {
			attrs = trace.StatusCodeNotFound(err.Error())
		}
	case byID && (errors.Is(err, errors.ErrObjectIDNotFound(uuid)) ||
		strings.Contains(err.Error(), "ngt uuid "+uuid+"'s object not found")):
		err = status.New(codes.NotFound, rpcName+" API uuid "+uuid+"'s object not found").Err()
		if span != nil {
			attrs = trace.StatusCodeNotFound(err.Error())
		}
	case errors.As(err, &errNGT):
		log.Errorf("ngt core process returned error: %v", err)
		err = status.WrapWithInternal(rpcName+" API failed to process search request due to ngt core process returned error", err,
			&errdetails.RequestInfo{
				RequestId:   reqID,
				ServingData: errdetails.Serialize(req),
			},
			s.resourceInfo(resType+"/core.ngt"), info.Get())
		log.Error(err)
		if span != nil {
			attrs = trace.StatusCodeInternal(err.Error())
		}
	case errors.Is(err, errors.ErrIncompatibleDimensionSize(vecLen, s.ngt.GetDimensionSize())):
		err = status.WrapWithInvalidArgument(rpcName+" API Incompatible Dimension Size detected",
			err,
			&errdetails.RequestInfo{
				RequestId:   reqID,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: resType,
			})
		log.Warn(err)
		if span != nil {
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		}
	case errors.Is(err, errors.ErrNaNOrInfVector):
		if byID {
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("%s API stored vector of %s contains NaN or Inf", rpcName, uuid),
				err,
				&errdetails.RequestInfo{
					RequestId:   reqID,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: resType,
					ResourceName: uuid,
				})
			log.Warn(err)
			if span != nil {
				attrs = trace.StatusCodeFailedPrecondition(err.Error())
			}
		} else {
			err = status.WrapWithInvalidArgument(rpcName+" API vector with NaN or Inf value detected",
				err,
				&errdetails.RequestInfo{
					RequestId:   reqID,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       vectorField,
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: resType,
				})
			log.Warn(err)
			if span != nil {
				attrs = trace.StatusCodeInvalidArgument(err.Error())
			}
		}
	default:
		err = status.WrapWithInternal(rpcName+" API failed to process search request", err,
			&errdetails.RequestInfo{
				RequestId:   reqID,
				ServingData: errdetails.Serialize(req),
			},
			s.resourceInfo(resType), info.Get())
		log.Error(err)
		if span != nil {
			attrs = trace.StatusCodeInternal(err.Error())
		}
	}
	if span != nil {
		errhandler.RecordSpanAttrs(span, attrs, err)
	}
	return err
}

func (s *server) Search(
	ctx context.Context, req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(ctx, apiName+"/"+vald.SearchRPCName)
		defer trace.End(span)
	}
	if err = s.validateVectorDimension(span, vald.SearchRPCName, ngtResourceType+"/ngt.Search",
		req.GetConfig().GetRequestId(), req, len(req.GetVector()), s.ngt.GetDimensionSize()); err != nil {
		return nil, err
	}
	if err = s.validateVectorValue(span, vald.SearchRPCName, ngtResourceType+"/ngt.Search",
		req.GetConfig().GetRequestId(), req, req.GetVector()); err != nil {
		return nil, err
	}
	res, err = s.ngt.SearchValidated(ctx,
		req.GetVector(),
		req.GetConfig().GetNum(),
		req.GetConfig().GetEpsilon(),
		req.GetConfig().GetRadius(),
		req.GetConfig().GetEdgeSize())
	if err == nil && res == nil {
		return nil, nil //nolint:nilnil // preserve existing nil-response contract when NGT returns (nil, nil)
	}
	if err != nil {
		return nil, s.handleSearchError(span, vald.SearchRPCName, "", false, req.GetConfig(), req, len(req.GetVector()), err)
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func (s *server) SearchByID(
	ctx context.Context, req *payload.Search_IDRequest,
) (res *payload.Search_Response, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(ctx, apiName+"/"+vald.SearchByIDRPCName)
		defer trace.End(span)
	}
	uuid := req.GetId()
	if err = s.validateUUID(span, vald.SearchByIDRPCName, ngtResourceType+"/ngt.SearchByID",
		uuid, req); err != nil {
		return nil, err
	}
	vec, res, err := s.ngt.SearchByID(ctx,
		uuid,
		req.GetConfig().GetNum(),
		req.GetConfig().GetEpsilon(),
		req.GetConfig().GetRadius(),
		req.GetConfig().GetEdgeSize())
	if err == nil && res == nil {
		return nil, nil //nolint:nilnil // preserve existing nil-response contract when NGT returns (nil, nil)
	}
	if err != nil {
		return nil, s.handleSearchError(span, vald.SearchByIDRPCName, uuid, true, req.GetConfig(), req, len(vec), err)
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func searchReqSpanSuffix(req *payload.Search_Request) string {
	return "requestID-" + req.GetConfig().GetRequestId()
}

func idReqSpanSuffix(req *payload.Search_IDRequest) string {
	return "id-" + req.GetId()
}

func streamSearch[
	Q proto.Message,
	S grpc.TypedServerStream[Q, *payload.Search_StreamResponse],
](
	stream S,
	rpcName string,
	streamConcurrency int,
	spanSuffix func(Q) string,
	search func(context.Context, Q) (*payload.Search_Response, error),
) error {
	return grpc.HandleBidirectionalStream(
		stream,
		streamConcurrency,
		"",
		apiName+"/"+rpcName,
		spanSuffix,
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
	return streamSearch(stream, vald.StreamSearchRPCName, s.streamConcurrency, searchReqSpanSuffix, s.Search)
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) error {
	return streamSearch(stream, vald.StreamSearchByIDRPCName, s.streamConcurrency, idReqSpanSuffix, s.SearchByID)
}

func (s *server) MultiSearch(
	ctx context.Context, reqs *payload.Search_MultiRequest,
) (res *payload.Search_Responses, errs error) {
	return multiSearch(ctx, s, vald.MultiSearchRPCName, reqs.Requests, s.Search)
}

func (s *server) MultiSearchByID(
	ctx context.Context, reqs *payload.Search_MultiIDRequest,
) (res *payload.Search_Responses, errs error) {
	return multiSearch(ctx, s, vald.MultiSearchByIDRPCName, reqs.Requests, s.SearchByID)
}

// multiSearch is the shared fan-out of MultiSearch, MultiSearchByID,
// MultiLinearSearch and MultiLinearSearchByID: each request runs search on the
// bounded multi group, per-item errors are collected and joined, and results
// keep the request order.
func multiSearch[Q any](
	ctx context.Context,
	s *server,
	rpcName string,
	reqs []Q,
	search func(context.Context, Q) (*payload.Search_Response, error),
) (res *payload.Search_Responses, errs error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(ctx, apiName+"/"+rpcName)
		defer trace.End(span)
	}

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs)),
	}
	eg := errgroup.NewBoundedFanOut(ctx, s.streamConcurrency)
	var mu sync.Mutex
	var errList []error
	for i, req := range reqs {
		idx, query := i, req
		eg.Go(safety.RecoverFunc(func() error {
			sctx := ctx
			var sspan trace.Span
			if trace.Enabled() {
				sctx, sspan = trace.StartSpan(ctx, fmt.Sprintf("%s/%s/errgroup.Go/id-%d", apiName, rpcName, idx))
				defer trace.End(sspan)
			}
			r, err := search(sctx, query)
			if err != nil {
				if sspan != nil {
					errhandler.RecordSpanFromError(sspan, err)
				}
				mu.Lock()
				errList = append(errList, err)
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		}))
	}
	if err := eg.Wait(); err != nil {
		errList = append(errList, err)
	}
	errs = errors.Join(errList...)
	if errs != nil {
		if span != nil {
			errhandler.RecordSpanFromError(span, errs)
		}
		return nil, errs
	}
	return res, nil
}
