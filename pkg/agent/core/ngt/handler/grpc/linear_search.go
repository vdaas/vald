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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/observability/trace"
)

func (s *server) LinearSearch(
	ctx context.Context, req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(ctx, apiName+"/"+vald.LinearSearchRPCName)
		defer trace.End(span)
	}
	if err = s.validateVectorDimension(span, vald.LinearSearchRPCName, ngtResourceType+"/ngt.LinearSearch",
		req.GetConfig().GetRequestId(), req, len(req.GetVector()), s.ngt.GetDimensionSize()); err != nil {
		return nil, err
	}
	if err = s.validateVectorValue(span, vald.LinearSearchRPCName, ngtResourceType+"/ngt.LinearSearch",
		req.GetConfig().GetRequestId(), req, req.GetVector()); err != nil {
		return nil, err
	}
	res, err = s.ngt.LinearSearchValidated(ctx,
		req.GetVector(),
		req.GetConfig().GetNum())
	if err == nil && res == nil {
		return nil, nil
	}
	if err != nil {
		return nil, s.handleSearchError(span, vald.LinearSearchRPCName, "", false, req.GetConfig(), req, len(req.GetVector()), err)
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func (s *server) LinearSearchByID(
	ctx context.Context, req *payload.Search_IDRequest,
) (res *payload.Search_Response, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(ctx, apiName+"/"+vald.LinearSearchByIDRPCName)
		defer trace.End(span)
	}
	uuid := req.GetId()
	if err = s.validateUUID(span, vald.LinearSearchByIDRPCName, ngtResourceType+"/ngt.LinearSearchByID",
		uuid, req); err != nil {
		return nil, err
	}
	vec, res, err := s.ngt.LinearSearchByID(ctx,
		uuid,
		req.GetConfig().GetNum())
	if err == nil && res == nil {
		return nil, nil
	}
	if err != nil {
		return nil, s.handleSearchError(span, vald.LinearSearchByIDRPCName, uuid, true, req.GetConfig(), req, len(vec), err)
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func (s *server) StreamLinearSearch(stream vald.Search_StreamLinearSearchServer) error {
	return streamSearch(stream, vald.StreamLinearSearchRPCName, s.streamConcurrency, searchReqSpanSuffix, s.LinearSearch)
}

func (s *server) StreamLinearSearchByID(stream vald.Search_StreamLinearSearchByIDServer) error {
	return streamSearch(stream, vald.StreamLinearSearchByIDRPCName, s.streamConcurrency, idReqSpanSuffix, s.LinearSearchByID)
}

func (s *server) MultiLinearSearch(
	ctx context.Context, reqs *payload.Search_MultiRequest,
) (res *payload.Search_Responses, errs error) {
	return multiSearch(ctx, s, vald.MultiLinearSearchRPCName, reqs.Requests, s.LinearSearch)
}

func (s *server) MultiLinearSearchByID(
	ctx context.Context, reqs *payload.Search_MultiIDRequest,
) (res *payload.Search_Responses, errs error) {
	return multiSearch(ctx, s, vald.MultiLinearSearchByIDRPCName, reqs.Requests, s.LinearSearchByID)
}
