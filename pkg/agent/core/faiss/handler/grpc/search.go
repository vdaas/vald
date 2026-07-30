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
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
)

func (s *server) Search(
	ctx context.Context, req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.SearchRPCName)
	defer trace.End(span)
	if len(req.GetVector()) != s.faiss.GetDimensionSize() {
		err = errors.ErrIncompatibleDimensionSize(len(req.GetVector()), s.faiss.GetDimensionSize())
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
			// ResourceType-only ResourceInfo (no ResourceName) preserves the
			// pre-refactor wire output and keeps parity with ngt.Search. The shared
			// validateVectorDimension helper (used by Insert/Update) always sets
			// ResourceName, so Search intentionally builds this inline instead.
			&errdetails.ResourceInfo{
				ResourceType: faissResourceType + "/faiss.Search",
			})
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.StatusCodeInvalidArgument(err.Error()), err)
		return nil, err
	}
	res, err = s.faiss.Search(
		req.GetConfig().GetNum(),
		req.GetConfig().GetNprobe(),
		1,
		req.GetVector())
	if err == nil && res == nil {
		return nil, nil
	}
	if err != nil || res == nil {
		var attrs []attribute.KeyValue
		switch {
		case errors.Is(err, errors.ErrCreateIndexingIsInProgress):
			err = status.WrapWithAborted("Search API aborted to process search request due to createing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(faissResourceType+"/faiss.Search"))
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
				s.resourceInfo(faissResourceType+"/faiss.Search"))
			log.Debug(err)
			attrs = trace.StatusCodeNotFound(err.Error())
		case errors.As(err, &errFaiss):
			log.Errorf("faiss core process returned error: %v", err)
			err = status.WrapWithInternal("Search API failed to process search request due to faiss core process returned error", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(faissResourceType+"/faiss.Search/core.faiss"), info.Get())
			log.Error(err)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrIncompatibleDimensionSize(len(req.GetVector()), int(s.faiss.GetDimensionSize()))):
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
				// ResourceType-only (no ResourceName), matching the up-front check,
				// main's wire output, and ngt.Search.
				&errdetails.ResourceInfo{
					ResourceType: faissResourceType + "/faiss.Search",
				})
			log.Warn(err)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		default:
			err = status.WrapWithInternal("Search API failed to process search request", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(faissResourceType+"/faiss.Search"), info.Get())
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
	return s.UnimplementedSearchServer.SearchByID(ctx, req)
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) (err error) {
	return s.UnimplementedSearchServer.StreamSearch(stream)
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) (err error) {
	return s.UnimplementedSearchServer.StreamSearchByID(stream)
}

func (s *server) MultiSearch(
	ctx context.Context, reqs *payload.Search_MultiRequest,
) (res *payload.Search_Responses, errs error) {
	return s.UnimplementedSearchServer.MultiSearch(ctx, reqs)
}

func (s *server) MultiSearchByID(
	ctx context.Context, reqs *payload.Search_MultiIDRequest,
) (res *payload.Search_Responses, errs error) {
	return s.UnimplementedSearchServer.MultiSearchByID(ctx, reqs)
}
