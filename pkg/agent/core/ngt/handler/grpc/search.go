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
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

func (s *server) Search(
	ctx context.Context, req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.SearchRPCName)
	defer trace.End(span)
	if len(req.GetVector()) != s.ngt.GetDimensionSize() {
		err = errors.ErrIncompatibleDimensionSize(len(req.GetVector()), int(s.ngt.GetDimensionSize()))
		err = status.WrapWithInvalidArgument("Search API Incompatible Dimension Size detected",
			err,
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
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.Search",
			})
		log.Warn(err)
		return errhandler.HandleError[payload.Search_Response](span, codes.InvalidArgument, err)
	}
	res, err = s.ngt.Search(ctx,
		req.GetVector(),
		req.GetConfig().GetNum(),
		req.GetConfig().GetEpsilon(),
		req.GetConfig().GetRadius(),
		req.GetConfig().GetEdgeSize())
	if err == nil && res == nil {
		return nil, nil
	}
	if err != nil || res == nil {
		var attrs []attribute.KeyValue
		switch {
		case errors.Is(err, errors.ErrCreateIndexingIsInProgress):
			err = status.WrapWithAborted("Search API aborted to process search request due to creating indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Search"))
			log.Debug(err)
			attrs = trace.StatusCodeAborted(err.Error())
		case errors.Is(err, errors.ErrFlushingIsInProgress):
			err = status.WrapWithAborted("Search API aborted to process search request due to flushing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Search"))
			log.Debug(err)
			attrs = trace.StatusCodeAborted(err.Error())
		case errors.Is(err, errors.ErrEmptySearchResult),
			err == nil && res == nil,
			0 < req.GetConfig().GetMinNum() && len(res.GetResults()) < int(req.GetConfig().GetMinNum()):
			err = status.WrapWithNotFound(fmt.Sprintf("Search API requestID %s's search result not found", req.GetConfig().GetRequestId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Search"))
			log.Debug(err)
			attrs = trace.StatusCodeNotFound(err.Error())
		case errors.As(err, &errNGT):
			log.Errorf("ngt core process returned error: %v", err)
			err = status.WrapWithInternal("Search API failed to process search request due to ngt core process returned error", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Search/core.ngt"), info.Get())
			log.Error(err)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrIncompatibleDimensionSize(len(req.GetVector()), int(s.ngt.GetDimensionSize()))):
			err = status.WrapWithInvalidArgument("Search API Incompatible Dimension Size detected",
				err,
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
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Search",
				})
			log.Warn(err)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		default:
			err = status.WrapWithInternal("Search API failed to process search request", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Search"), info.Get())
			log.Error(err)
			attrs = trace.StatusCodeInternal(err.Error())
		}
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, err
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func (s *server) SearchByID(
	ctx context.Context, req *payload.Search_IDRequest,
) (res *payload.Search_Response, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.SearchByIDRPCName)
	defer trace.End(span)
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
		return nil, nil
	}
	if err != nil || res == nil {
		var attrs []attribute.KeyValue
		switch {
		case errors.Is(err, errors.ErrCreateIndexingIsInProgress):
			err = status.WrapWithAborted("SearchByID API aborted to process search request due to creating indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.SearchByID"))
			log.Debug(err)
			attrs = trace.StatusCodeAborted(err.Error())
		case errors.Is(err, errors.ErrFlushingIsInProgress):
			err = status.WrapWithAborted("SearchByID API aborted to process search request due to flushing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.SearchByID"))
			log.Debug(err)
			attrs = trace.StatusCodeAborted(err.Error())
		case errors.Is(err, errors.ErrEmptySearchResult),
			err == nil && res == nil,
			0 < req.GetConfig().GetMinNum() && len(res.GetResults()) < int(req.GetConfig().GetMinNum()):
			err = status.WrapWithNotFound(fmt.Sprintf("SearchByID API uuid %s's search result not found", req.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.SearchByID"))
			log.Debug(err)
			attrs = trace.StatusCodeNotFound(err.Error())
		case errors.Is(err, errors.ErrObjectIDNotFound(req.GetId())),
			strings.Contains(err.Error(), fmt.Sprintf("ngt uuid %s's object not found", req.GetId())):
			err = status.WrapWithNotFound(fmt.Sprintf("SearchByID API uuid %s's object not found", req.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.SearchByID"))
			log.Debug(err)
			attrs = trace.StatusCodeNotFound(err.Error())
		case errors.As(err, &errNGT):
			log.Errorf("ngt core process returned error: %v", err)
			err = status.WrapWithInternal("SearchByID API failed to process search request due to ngt core process returned error", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.SearchByID/core.ngt"), info.Get())
			log.Error(err)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrIncompatibleDimensionSize(len(vec), int(s.ngt.GetDimensionSize()))):
			err = status.WrapWithInvalidArgument("SearchByID API Incompatible Dimension Size detected",
				err,
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
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.SearchByID",
				})
			log.Warn(err)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		default:
			err = status.WrapWithInternal("SearchByID API failed to process search request", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.SearchByID"), info.Get())
			log.Error(err)
			attrs = trace.StatusCodeInternal(err.Error())
		}
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, err
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamSearchRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_Request) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
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
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamSearchByIDRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamSearchByIDRPCName+"/id-"+req.GetId())
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
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiSearchRPCName)
	defer trace.End(span)

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() (err error) {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s/%s/errgroup.Go/id-%d", apiName, vald.MultiSearchRPCName, idx))
			defer trace.End(sspan)
			r, err := s.Search(ctx, query)
			if err != nil {
				st, _ := status.FromError(err)
				errhandler.RecordSpanStatus(sspan, st, err)
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Join(errs, err)
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, _ := status.FromError(errs)
		errhandler.RecordSpanStatus(span, st, errs)
		return nil, errs
	}
	return res, nil
}

func (s *server) MultiSearchByID(
	ctx context.Context, reqs *payload.Search_MultiIDRequest,
) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiSearchByIDRPCName)
	defer trace.End(span)

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s/%s/errgroup.Go/id-%d", apiName, vald.MultiSearchByIDRPCName, idx))
			defer trace.End(sspan)
			defer wg.Done()
			r, err := s.SearchByID(ctx, query)
			if err != nil {
				st, _ := status.FromError(err)
				errhandler.RecordSpanStatus(sspan, st, err)
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Join(errs, err)
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, _ := status.FromError(errs)
		errhandler.RecordSpanStatus(span, st, errs)
		return nil, errs
	}
	return res, nil
}
