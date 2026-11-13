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
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
)

func (s *server) LinearSearch(
	ctx context.Context, req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.LinearSearchRPCName), apiName+"/"+vald.LinearSearchRPCName)
	defer trace.End(span)
	vl := len(req.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument(vald.LinearSearchRPCName+" API invalid vector argument", err,
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
		vald.LinearSearchRPCName,
		req.GetConfig(),
		func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.LinearSearch(ctx, &payload.Search_Request{Vector: req.GetVector(), Config: fcfg}, copts...)
		},
	)
	if err != nil {
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, err
	}
	return res, nil
}

func (s *server) LinearSearchByID(
	ctx context.Context, req *payload.Search_IDRequest,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.LinearSearchByIDRPCName), apiName+"/"+vald.LinearSearchByIDRPCName)
	defer trace.End(span)
	uuid := req.GetId()
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.LinearSearchByIDRPCName+" API invalid uuid", err, append(s.errDetails(ctx, vald.LinearSearchByIDRPCName, uuid, req),
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
		// try search by using agent's LinearSearchByID method this operation is emergency fallback, the search quality is not same as usual LinearSearchByID operation.
		var attrs []attribute.KeyValue
		res, attrs, err = s.doSearch(
			ctx,
			vald.LinearSearchByIDRPCName,
			req.GetConfig(),
			func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
				return vc.LinearSearchByID(ctx, &payload.Search_IDRequest{Id: req.GetId(), Config: fcfg}, copts...)
			},
		)
		if err == nil {
			return res, nil
		}
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, err
	}
	res, err = s.LinearSearch(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: req.GetConfig(),
	})
	if err != nil {
		var attrs []attribute.KeyValue
		res, attrs, err = s.doSearch(
			ctx,
			vald.LinearSearchByIDRPCName,
			req.GetConfig(),
			func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
				return vc.LinearSearchByID(ctx, &payload.Search_IDRequest{Id: req.GetId(), Config: fcfg}, copts...)
			},
		)
		if err != nil {
			errhandler.RecordSpanAttrs(span, attrs, err)
			return nil, err
		}
	}
	return res, nil
}

func (s *server) StreamLinearSearch(stream vald.Search_StreamLinearSearchServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamLinearSearchRPCName), apiName+"/"+vald.StreamLinearSearchRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_Request) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamLinearSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer trace.End(sspan)
			res, err := s.LinearSearch(ctx, req)
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

func (s *server) StreamLinearSearchByID(
	stream vald.Search_StreamLinearSearchByIDServer,
) (err error) {
	ctx, span := trace.StartSpan(
		grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamLinearSearchByIDRPCName),
		apiName+"/"+vald.StreamLinearSearchByIDRPCName,
	)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamLinearSearchByIDRPCName+"/id-"+req.GetId())
			defer trace.End(sspan)
			res, err := s.LinearSearchByID(ctx, req)
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
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiLinearSearch(
	ctx context.Context, reqs *payload.Search_MultiRequest,
) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiLinearSearchRPCName), apiName+"/"+vald.MultiLinearSearchRPCName)
	defer trace.End(span)
	res, errs = multiSearch(ctx, s, vald.MultiLinearSearchRPCName, reqs.GetRequests(), s.LinearSearch)
	if errs != nil {
		st, _ := status.FromError(errs)
		errhandler.RecordSpanStatus(span, st, errs)
	}
	return res, errs
}

func (s *server) MultiLinearSearchByID(
	ctx context.Context, reqs *payload.Search_MultiIDRequest,
) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiLinearSearchByIDRPCName), apiName+"/"+vald.MultiLinearSearchByIDRPCName)
	defer trace.End(span)
	res, errs = multiSearch(ctx, s, vald.MultiLinearSearchByIDRPCName, reqs.GetRequests(), s.LinearSearchByID)
	if errs != nil {
		st, _ := status.FromError(errs)
		errhandler.RecordSpanStatus(span, st, errs)
	}
	return res, errs
}
