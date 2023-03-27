//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

type server struct {
	eg                errgroup.Group
	gateway           service.Gateway // Mirror Gateway client service.
	mirror            service.Mirror
	vAddr             string // Vald Gateway address (lb-gateway).
	streamConcurrency int
	name              string
	ip                string
	mirror.UnimplementedValdServerWithMirror
}

const apiName = "vald/gateway/mirror"

func New(opts ...Option) (mirror.Server, error) {
	s := new(server)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(s); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(err, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	return s, nil
}

func (s *server) Register(ctx context.Context, req *payload.Mirror_Targets) (*payload.Mirror_Targets, error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, mirror.PackageName+"."+mirror.MirrorRPCServiceName+"/"+mirror.RegisterRPCName), apiName+"/"+mirror.RegisterRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err := s.mirror.Connect(ctx, req.GetTargets()...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + mirror.RegisterRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				mirror.RegisterRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + mirror.RegisterRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithCanceled(
				mirror.RegisterRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + mirror.RegisterRPCName,
			)
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				mirror.RegisterRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInvalidArgument(
				mirror.RegisterRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+mirror.RegisterRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return req, nil
}

func (s *server) Advertise(ctx context.Context, req *payload.Mirror_Targets) (res *payload.Mirror_Targets, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, mirror.PackageName+"."+mirror.MirrorRPCServiceName+"/"+mirror.AdvertiseRPCName), apiName+"/"+mirror.AdvertiseRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = s.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	tgts, err := s.mirror.MirrorTargets()
	if err != nil {
		err = status.WrapWithInternal(mirror.AdvertiseRPCName+" API failed to get connected mirror gateway targets", err,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "mirror gateway targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + mirror.AdvertiseRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return &payload.Mirror_Targets{
		Targets: tgts,
	}, nil
}

func (s *server) Exists(ctx context.Context, meta *payload.Object_ID) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.ExistsRPCName), apiName+"/"+vald.ExistsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
		id, err = vc.Exists(ctx, meta, copts...)
		if err != nil {
			return nil, err
		}
		return id, nil
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(meta),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.ExistsRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.ExistsRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.ExistsRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.ExistsRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.ExistsRPCName,
			)
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInternal(
				vald.ExistsRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.ExistsRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.ExistsRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return id, nil
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.SearchRPCName), apiName+"/"+vald.SearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.Search(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetConfig().GetRequestId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.SearchRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.SearchRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName,
			)
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInternal(
				vald.SearchRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.SearchRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.SearchRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error,
) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.SearchByIDRPCName), apiName+"/"+vald.SearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.SearchByID(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetConfig().GetRequestId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.SearchByIDRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.SearchByIDRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
			)
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInternal(
				vald.SearchByIDRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.SearchByIDRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.SearchByIDRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamSearchRPCName), apiName+"/"+vald.StreamSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_Request) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Search(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchRPCName+" gRPC error response")
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
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
		},
	)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.StreamSearchRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamSearchByIDRPCName), apiName+"/"+vald.StreamSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"."+vald.StreamSearchByIDRPCName+"/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.SearchByID(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchByIDRPCName+" gRPC error response")
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
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
		},
	)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.StreamSearchByIDRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiSearch(ctx context.Context, req *payload.Search_MultiRequest) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiSearchRPCName), apiName+"/"+vald.MultiSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.MultiSearch(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiSearchRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.MultiSearchRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiSearchRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.MultiSearchRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiSearchRPCName,
			)
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInternal(
				vald.MultiSearchRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.MultiSearchRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.MultiSearchRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) MultiSearchByID(ctx context.Context, req *payload.Search_MultiIDRequest) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiSearchByIDRPCName), apiName+"/"+vald.MultiSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.MultiSearchByID(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiSearchByIDRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.MultiSearchByIDRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiSearchByIDRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.MultiSearchByIDRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiSearchByIDRPCName,
			)
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInternal(
				vald.MultiSearchByIDRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.MultiSearchByIDRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.MultiSearchByIDRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) LinearSearch(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.LinearSearchRPCName), apiName+"/"+vald.LinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.LinearSearch(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetConfig().GetRequestId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.LinearSearchRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.LinearSearchRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchRPCName,
			)
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInternal(
				vald.LinearSearchRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.LinearSearchRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.LinearSearchRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) LinearSearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error,
) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.LinearSearchByIDRPCName), apiName+"/"+vald.LinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.LinearSearchByID(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetConfig().GetRequestId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.LinearSearchByIDRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.LinearSearchByIDRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
			)
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInternal(
				vald.LinearSearchByIDRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.LinearSearchByIDRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.LinearSearchByIDRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) StreamLinearSearch(stream vald.Search_StreamLinearSearchServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamLinearSearchRPCName), apiName+"/"+vald.StreamLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_Request) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamLinearSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearch(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchRPCName+" gRPC error response")
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
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
		},
	)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.StreamLinearSearchRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) StreamLinearSearchByID(stream vald.Search_StreamLinearSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamLinearSearchByIDRPCName), apiName+"/"+vald.StreamLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"."+vald.StreamLinearSearchByIDRPCName+"/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearchByID(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchByIDRPCName+" gRPC error response")
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
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
		},
	)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.StreamLinearSearchByIDRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiLinearSearch(ctx context.Context, req *payload.Search_MultiRequest) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiLinearSearchRPCName), apiName+"/"+vald.MultiLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.MultiLinearSearch(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiLinearSearchRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.MultiLinearSearchRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiLinearSearchRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.MultiLinearSearchRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiLinearSearchRPCName,
			)
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInternal(
				vald.MultiLinearSearchRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.MultiLinearSearchRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.MultiLinearSearchRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) MultiLinearSearchByID(ctx context.Context, req *payload.Search_MultiIDRequest) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiLinearSearchByIDRPCName), apiName+"/"+vald.MultiLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.MultiLinearSearchByID(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiLinearSearchByIDRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.MultiLinearSearchByIDRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiLinearSearchByIDRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.MultiLinearSearchByIDRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiLinearSearchByIDRPCName,
			)
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInternal(
				vald.MultiLinearSearchByIDRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.MultiLinearSearchByIDRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.MultiLinearSearchByIDRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (ce *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.InsertRPCName), apiName+"/"+vald.InsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	reqSrcPodName := s.gateway.FromForwardedContext(ctx)

	// When this condition is matched, the request is proxied to another Mirror gateway.
	// So this component sends requests only to the Vald gateway (LB gateway) of its own cluster.
	if len(reqSrcPodName) != 0 {
		_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
			ce, err = vc.Insert(ctx, req, copts...)
			if err != nil {
				return nil, err
			}
			return ce, nil
		})
		if err != nil {
			reqInfo := &errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			}
			resInfo := &errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
			}
			var attrs trace.Attributes

			switch {
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(
					vald.InsertRPCName+" API canceld", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeCancelled(
					errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
				)
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(
					vald.InsertRPCName+" API deadline exceeded", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeDeadlineExceeded(
					errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
				)
			case errors.Is(err, errors.ErrTargetNotFound):
				err = status.WrapWithInternal(
					vald.InsertRPCName+" API target not found", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeInternal(err.Error())
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(
					vald.InsertRPCName+" API connection not found", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeInternal(err.Error())
			default:
				var (
					st  *status.Status
					msg string
				)
				st, msg, err = status.ParseError(err, codes.Internal,
					"failed to parse "+vald.InsertRPCName+" gRPC error response", reqInfo, resInfo,
				)
				attrs = trace.FromGRPCStatus(st.Code(), msg)
			}
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		log.Debugf("Insert API succeeded to %#v", ce)
		return ce, nil
	}

	var mu sync.Mutex
	var result sync.Map
	ce = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		loc, err := s.insert(ctx, vc, req, copts...)
		if err != nil {
			st, _, _ := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.InsertRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			if st.Code() == codes.AlreadyExists {
				// NOTE: If it is strictly necessary to check, fix this logic.
				return nil
			}
		}
		if loc != nil {
			mu.Lock()
			ce.Name = loc.GetName()
			ce.Ips = append(ce.Ips, loc.GetIps()...)
			mu.Unlock()
		}
		result.Store(target, err)
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + ".BroadCast",
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.InsertRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.InsertRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	var errs error
	targets := make([]string, 0, 10)
	result.Range(func(target, err any) bool {
		if err == nil {
			targets = append(targets, target.(string))
		} else {
			if err, ok := err.(error); ok && err != nil {
				errs = errors.Join(errs, err)
			}
		}
		return true
	})
	if errs == nil {
		log.Debugf("Insert API mirror request succeeded to %#v", ce)
		return ce, nil
	}
	log.Error("failed to Insert API mirror request: %v, so starts the rollback request", errs)

	var emu sync.Mutex
	var rerrs error
	rmReq := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
	}
	err = s.gateway.DoMulti(ctx, targets,
		func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "rollback/BroadCast/"+target), apiName+"/"+vald.InsertRPCName+"/rollback/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			_, err := s.remove(ctx, vc, rmReq, copts...)
			if err != nil {
				st, _, err := status.ParseError(err, codes.Internal,
					"failed to parse "+vald.RemoveRPCName+" for "+vald.InsertRPCName+" error response for "+target,
					&errdetails.RequestInfo{
						RequestId:   rmReq.GetId().GetId(),
						ServingData: errdetails.Serialize(rmReq),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + "." + vald.RemoveRPCName + ".BroadCast/" + target,
						ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
					},
				)
				if st.Code() == codes.NotFound {
					return nil
				}
				emu.Lock()
				rerrs = errors.Join(rerrs, err)
				emu.Unlock()
				return err
			}
			return nil
		},
	)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   rmReq.GetId().GetId(),
			ServingData: errdetails.Serialize(rmReq),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + "." + vald.RemoveRPCName + ".BroadCast",
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.RemoveRPCName+" for "+vald.InsertRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveRPCName+" for "+vald.InsertRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if rerrs == nil {
		log.Debugf("rollback for Insert API mirror request succeeded to %v", targets)
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse "+vald.InsertRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	log.Debugf("failed to rollback for Insert API mirror request succeeded to %v: %v", targets, rerrs)
	st, msg, err := status.ParseError(rerrs, codes.Internal,
		"failed to parse "+vald.RemoveRPCName+" for "+vald.InsertRPCName+" gRPC error response",
		&errdetails.RequestInfo{
			RequestId:   rmReq.GetId().GetId(),
			ServingData: errdetails.Serialize(rmReq),
		},
		&errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + "." + vald.RemoveRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) %v", apiName, s.name, s.ip, targets),
		},
	)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
}

func (s *server) insert(ctx context.Context, client vald.InsertClient, req *payload.Insert_Request, opts ...grpc.CallOption) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "insert"), apiName+"/insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err = client.Insert(ctx, req, opts...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.InsertRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.InsertRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
			)
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.InsertRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.InsertRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.StreamInsertRPCName), apiName+"/"+vald.StreamInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Insert_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamInsertRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Insert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.InsertRPCName+" gRPC error response")
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: res,
				},
			}, nil
		},
	)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.StreamInsertRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.MultiInsertRPCName), apiName+"/"+vald.MultiInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup

	for i, r := range reqs.GetRequests() {
		idx, req := i, r
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + req.GetVector().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiInsertRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			loc, err := s.Insert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.InsertRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Locations[idx] = loc
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiInsertRPCName+" gRPC error response",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiInsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateRPCName), apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	reqSrcPodName := s.gateway.FromForwardedContext(ctx)

	// When this condition is matched, the request is proxied to another Mirror gateway.
	// So this component sends requests only to the Vald gateway (LB gateway) of its own cluster.
	if len(reqSrcPodName) != 0 {
		_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
			loc, err = s.update(ctx, vc, req, copts...)
			if err != nil {
				return nil, err
			}
			return loc, nil
		})
		if err != nil {
			reqInfo := &errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			}
			resInfo := &errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
			}
			var attrs trace.Attributes

			switch {
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(
					vald.UpdateRPCName+" API canceld", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeCancelled(
					errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
				)
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(
					vald.UpdateRPCName+" API deadline exceeded", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeDeadlineExceeded(
					errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
				)
			case errors.Is(err, errors.ErrTargetNotFound):
				err = status.WrapWithInternal(
					vald.UpdateRPCName+" API target not found", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeInternal(err.Error())
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(
					vald.UpdateRPCName+" API connection not found", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeInternal(err.Error())
			default:
				var (
					st  *status.Status
					msg string
				)
				st, msg, err = status.ParseError(err, codes.Internal,
					"failed to parse "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo,
				)
				attrs = trace.FromGRPCStatus(st.Code(), msg)
			}
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		log.Debugf("Update API succeeded to %#v", loc)
		return loc, nil
	}

	objReq := &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
	}
	oldVecs, err := s.getObjects(ctx, objReq)
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	var result sync.Map
	ce := &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}

	err = s.gateway.BroadCast(ctx,
		func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpdateRPCName+"/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			loc, err := s.update(ctx, vc, req, copts...)
			if err != nil {
				st, _, _ := status.ParseError(err, codes.Internal,
					"failed to parse "+vald.UpdateRPCName+" API error response for "+target,
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + ".BroadCast/" + target,
						ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
					},
				)
				if st.Code() == codes.AlreadyExists {
					// NOTE: If it is strictly necessary to check, fix this logic.
					return nil
				}
			}
			if loc != nil {
				mu.Lock()
				ce.Name = loc.GetName()
				ce.Ips = append(ce.Ips, loc.GetIps()...)
				mu.Unlock()
			}
			result.Store(target, err)
			return err
		},
	)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + ".BroadCast",
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpdateRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	var errs error
	targets := make([]string, 0, 10)
	result.Range(func(target, err any) bool {
		if err == nil {
			targets = append(targets, target.(string))
		} else {
			if err, ok := err.(error); ok && err != nil {
				errs = errors.Join(errs, err)
			}
		}
		return true
	})
	if errs == nil {
		log.Debugf("Update API mirror request succeeded to %#v", ce)
		return ce, nil
	}
	log.Error("failed to Update API mirror request: %v, so starts the rollback request", errs)

	var emu sync.Mutex
	var rerrs error
	rmReq := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
	}

	err = s.gateway.DoMulti(ctx, targets,
		func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "rollback/BroadCast/"+target), apiName+"/"+vald.RemoveRPCName+"/rollback/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			oldVec, ok := oldVecs.Load(target)
			if !ok || oldVec == nil {
				_, err := s.remove(ctx, vc, rmReq, copts...)
				if err != nil {
					st, _, _ := status.ParseError(err, codes.Internal,
						"failed to parse "+vald.RemoveRPCName+" for "+vald.UpdateRPCName+" gRPC error response",
						&errdetails.RequestInfo{
							RequestId:   rmReq.GetId().GetId(),
							ServingData: errdetails.Serialize(rmReq),
						},
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.RemoveRPCName + ".BroadCast/" + target,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						},
					)
					if st.Code() == codes.NotFound {
						return nil
					}
					emu.Lock()
					errs = errors.Join(errs, err)
					emu.Unlock()
					return err
				}
				return nil
			}

			req := &payload.Update_Request{
				Vector: oldVec.(*payload.Object_Vector),
				Config: &payload.Update_Config{
					SkipStrictExistCheck: true,
				},
			}
			_, err := s.update(ctx, vc, req, copts...)
			if err != nil {
				st, _, _ := status.ParseError(err, codes.Internal,
					"failed to parse "+vald.UpdateRPCName+" for "+vald.UpdateRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.UpdateRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
					},
				)
				if st.Code() == codes.AlreadyExists {
					return nil
				}
				emu.Lock()
				errs = errors.Join(errs, err)
				emu.Unlock()
				return err
			}
			return nil
		},
	)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId: req.GetVector().GetId(),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + ".Rollback.BroadCast",
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpdateRPCName+" for Rollback connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpdateRPCName+" for Rollback gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if rerrs == nil {
		log.Debugf("rollback for Update API mirror request succeeded to %v", targets)
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse "+vald.UpdateRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	log.Debugf("failed to rollback for Update API mirror request succeeded to %v: %v", targets, rerrs)
	st, msg, err := status.ParseError(rerrs, codes.Internal,
		"failed to parse "+vald.UpdateRPCName+" for Rollback gRPC error response",
		&errdetails.RequestInfo{
			RequestId: req.GetVector().GetId(),
		},
		&errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + ".Rollback",
			ResourceName: fmt.Sprintf("%s: %s(%s) %v", apiName, s.name, s.ip, targets),
		},
	)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
}

func (s *server) update(ctx context.Context, client vald.UpdateClient, req *payload.Update_Request, opts ...grpc.CallOption) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "update"), apiName+"/update")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err = client.Update(ctx, req, opts...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.UpdateRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.UpdateRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
			)
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.UpdateRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.StreamUpdateRPCName), apiName+"/"+vald.StreamUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Update_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamUpdateRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Update(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpdateRPCName+" gRPC error response")
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: res,
				},
			}, nil
		},
	)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.StreamUpdateRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.MultiUpdateRPCName), apiName+"/"+vald.MultiUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup

	for i, r := range reqs.GetRequests() {
		idx, req := i, r
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + req.GetVector().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiUpdateRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			loc, err := s.Update(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpdateRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Locations[idx] = loc
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiUpdateRPCName+" gRPC error response",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpdateRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.UpsertRPCName), apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	reqSrcPodName := s.gateway.FromForwardedContext(ctx)

	// When this condition is matched, the request is proxied to another Mirror gateway.
	// So this component sends requests only to the Vald gateway (LB gateway) of its own cluster.
	if len(reqSrcPodName) != 0 {
		_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
			loc, err = vc.Upsert(ctx, req, copts...)
			if err != nil {
				return nil, err
			}
			return loc, nil
		})
		if err != nil {
			reqInfo := &errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			}
			resInfo := &errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
			}
			var attrs trace.Attributes

			switch {
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(
					vald.UpsertRPCName+" API canceld", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeCancelled(
					errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
				)
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(
					vald.UpsertRPCName+" API deadline exceeded", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeDeadlineExceeded(
					errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
				)
			case errors.Is(err, errors.ErrTargetNotFound):
				err = status.WrapWithInternal(
					vald.UpsertRPCName+" API target not found", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeInternal(err.Error())
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(
					vald.UpsertRPCName+" API connection not found", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeInternal(err.Error())
			default:
				var (
					st  *status.Status
					msg string
				)
				st, msg, err = status.ParseError(err, codes.Internal,
					"failed to parse "+vald.UpsertRPCName+" gRPC error response", reqInfo, resInfo,
				)
				attrs = trace.FromGRPCStatus(st.Code(), msg)
			}
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		log.Debugf("Upsert API succeeded to %#v", loc)
		return loc, nil
	}

	objReq := &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
	}
	oldVecs, err := s.getObjects(ctx, objReq)
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	var result sync.Map
	loc = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpsertRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		ce, err := s.upsert(ctx, vc, req, copts...)
		if err != nil {
			st, _, _ := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpsertRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			if st.Code() == codes.AlreadyExists {
				// NOTE: If it is strictly necessary to check, fix this logic.
				return nil
			}
		}
		if ce != nil {
			mu.Lock()
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
			mu.Unlock()
		}
		result.Store(target, err)
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + ".BroadCast",
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpsertRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpsertRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	var errs error
	targets := make([]string, 0, 10)
	result.Range(func(target, err any) bool {
		if err == nil {
			targets = append(targets, target.(string))
		} else {
			if err, ok := err.(error); ok && err != nil {
				errs = errors.Join(errs, err)
			}
		}
		return true
	})
	if errs == nil {
		log.Debugf("Upsert API mirror request succeeded to %#v", loc)
		return loc, nil
	}
	log.Error("failed to Upsert API mirror request: %v, so starts the rollback request", errs)

	var emu sync.Mutex
	var rerrs error
	rmReq := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
	}
	err = s.gateway.DoMulti(ctx, targets,
		func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "rollback/BroadCast/"+target), apiName+"/"+vald.UpsertRPCName+"/rollback/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			oldVec, ok := oldVecs.Load(target)
			if !ok || oldVec == nil {
				_, err := s.remove(ctx, vc, rmReq, copts...)
				if err != nil {
					st, _, _ := status.ParseError(err, codes.Internal,
						"failed to parse "+vald.RemoveRPCName+" for "+vald.UpsertRPCName+" gRPC error response",
						&errdetails.RequestInfo{
							RequestId:   rmReq.GetId().GetId(),
							ServingData: errdetails.Serialize(rmReq),
						},
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + vald.RemoveRPCName + ".BroadCast/" + target,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						},
					)
					if st.Code() == codes.NotFound {
						return nil
					}
					emu.Lock()
					errs = errors.Join(errs, err)
					emu.Unlock()
					return err
				}
				return nil
			}

			req := &payload.Update_Request{
				Vector: oldVec.(*payload.Object_Vector),
				Config: &payload.Update_Config{
					SkipStrictExistCheck: true,
				},
			}
			_, err := s.update(ctx, vc, req, copts...)
			if err != nil {
				st, _, _ := status.ParseError(err, codes.Internal,
					"failed to parse "+vald.UpdateRPCName+" for "+vald.UpsertRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + vald.UpdateRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
					},
				)
				if st.Code() == codes.AlreadyExists {
					return nil
				}
				emu.Lock()
				errs = errors.Join(errs, err)
				emu.Unlock()
				return err
			}
			return nil
		},
	)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId: req.GetVector().GetId(),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + ".Rollback.BroadCast",
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpsertRPCName+" for Rollback connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpsertRPCName+" for Rollback gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if rerrs == nil {
		log.Debugf("rollback for Upsert API mirror request succeeded to %v", targets)
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse "+vald.UpsertRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	log.Debugf("failed to rollback for Upsert API mirror request succeeded to %v: %v", targets, rerrs)
	st, msg, err := status.ParseError(rerrs, codes.Internal,
		"failed to parse "+vald.UpsertRPCName+" for Rollback gRPC error response",
		&errdetails.RequestInfo{
			RequestId: req.GetVector().GetId(),
		},
		&errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + ".Rollback",
			ResourceName: fmt.Sprintf("%s: %s(%s) %v", apiName, s.name, s.ip, targets),
		},
	)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
}

func (s *server) upsert(ctx context.Context, client vald.UpsertClient, req *payload.Upsert_Request, opts ...grpc.CallOption) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "upsert"), apiName+"/upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err = client.Upsert(ctx, req, opts...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.UpsertRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.UpsertRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
			)
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.UpsertRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpsertRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn("failed to process Upsert request\terror: %s", err.Error())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.StreamUpsertRPCName), apiName+"/"+vald.StreamUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Upsert_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamUpsertRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Upsert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpsertRPCName+" gRPC error response")
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: res,
				},
			}, nil
		},
	)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.StreamUpsertRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.MultiUpsertRPCName), apiName+"/"+vald.MultiUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup

	for i, r := range reqs.GetRequests() {
		idx, req := i, r
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + req.GetVector().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiUpsertRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			loc, err := s.Upsert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpsertRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Locations[idx] = loc
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiUpsertRPCName+" gRPC error response",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	reqSrcPodName := s.gateway.FromForwardedContext(ctx)

	// When this condition is matched, the request is proxied to another Mirror gateway.
	// So this component sends the request only to the Vald gateway (LB gateway) of own cluster.
	if len(reqSrcPodName) != 0 {
		_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
			loc, err = vc.Remove(ctx, req, copts...)
			if err != nil {
				return nil, err
			}
			return loc, nil
		})
		if err != nil {
			reqInfo := &errdetails.RequestInfo{
				RequestId:   req.GetId().GetId(),
				ServingData: errdetails.Serialize(req),
			}
			resInfo := &errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
			}
			var attrs trace.Attributes

			switch {
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(
					vald.RemoveRPCName+" API canceld", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeCancelled(
					errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
				)
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(
					vald.RemoveRPCName+" API deadline exceeded", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeDeadlineExceeded(
					errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
				)
			case errors.Is(err, errors.ErrTargetNotFound):
				err = status.WrapWithInternal(
					vald.RemoveRPCName+" API target not found", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeInternal(err.Error())
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(
					vald.RemoveRPCName+" API connection not found", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeInternal(err.Error())
			default:
				var (
					st  *status.Status
					msg string
				)
				st, msg, err = status.ParseError(err, codes.Internal,
					"failed to parse "+vald.RemoveRPCName+" gRPC error response", reqInfo, resInfo,
				)
				attrs = trace.FromGRPCStatus(st.Code(), msg)
			}
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		log.Debugf("Remove API remove succeeded to %#v", loc)
		return loc, nil
	}

	objReq := &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetId().GetId(),
		},
	}
	oldVecs, err := s.getObjects(ctx, objReq)
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	var result sync.Map
	ce := &payload.Object_Location{
		Uuid: req.GetId().GetId(),
		Ips:  make([]string, 0),
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.RemoveRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		loc, err := s.remove(ctx, vc, req, copts...)
		if err != nil {
			st, _, _ := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveRPCName+" gRPC error response for "+target,
				&errdetails.RequestInfo{
					RequestId:   req.GetId().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			if st.Code() == codes.NotFound {
				// NOTE: If it is strictly necessary to check, fix this logic.
				return nil
			}
		}
		if loc != nil {
			mu.Lock()
			ce.Name = loc.GetName()
			ce.Ips = append(ce.Ips, loc.GetIps()...)
			mu.Unlock()
		}
		result.Store(target, err)
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetId().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + ".BroadCast",
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.RemoveRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	var errs error
	targets := make([]string, 0, 10)
	result.Range(func(target, err any) bool {
		if err == nil {
			targets = append(targets, target.(string))
		} else {
			if err, ok := err.(error); ok && err != nil {
				errs = errors.Join(errs, err)
			}
		}
		return true
	})
	if errs == nil {
		log.Debugf("Remove API mirror request succeeded to %#v", ce)
		return ce, nil
	}
	log.Error("failed to Remove API mirror request: %v, so starts the rollback request", errs)

	var emu sync.Mutex
	var rerrs error
	err = s.gateway.DoMulti(ctx, targets,
		func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "rollback/BroadCast/"+target), apiName+"/"+vald.RemoveRPCName+"/rollback/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			objv, ok := oldVecs.Load(target)
			if !ok || objv == nil {
				log.Debug("failed to load old vector from  %s", target)
				return nil
			}
			req := &payload.Upsert_Request{
				Vector: objv.(*payload.Object_Vector),
				Config: &payload.Upsert_Config{
					SkipStrictExistCheck: true,
				},
			}
			_, err := s.upsert(ctx, vc, req, copts...)
			if err != nil {
				st, _, _ := status.ParseError(err, codes.Internal,
					"failed to parse "+vald.UpsertRPCName+" for "+vald.RemoveRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + "." + vald.UpsertRPCName + ".BroadCast/" + target,
						ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
					},
				)
				if st.Code() == codes.AlreadyExists {
					return nil
				}
				emu.Lock()
				rerrs = errors.Join(rerrs, err)
				emu.Unlock()
				return err
			}
			return nil
		},
	)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId: req.GetId().GetId(),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + "." + vald.UpsertRPCName + ".BroadCast",
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpsertRPCName+" for "+vald.RemoveRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpsertRPCName+" for "+vald.RemoveRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if rerrs == nil {
		log.Debugf("rollback for Remove API mirror request succeeded to %v", targets)
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse "+vald.RemoveRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetId().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	log.Debugf("failed to rollback for Remove API mirror request succeeded to %v: %v", targets, rerrs)
	st, msg, err := status.ParseError(rerrs, codes.Internal,
		"failed to parse "+vald.UpsertRPCName+" for "+vald.RemoveRPCName+" gRPC error response",
		&errdetails.RequestInfo{
			RequestId: req.GetId().GetId(),
		},
		&errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + "." + vald.UpsertRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) %v", apiName, s.name, s.ip, targets),
		},
	)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
}

func (s *server) remove(ctx context.Context, client vald.RemoveClient, req *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "remove"), apiName+"/remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err := client.Remove(ctx, req, opts...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetId().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.RemoveRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.RemoveRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
			)
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.RemoveRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.StreamRemoveRPCName), apiName+"/"+vald.StreamRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Remove_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamRemoveRPCName+"/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Remove(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.RemoveRPCName+" gRPC error response")
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: res,
				},
			}, nil
		},
	)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.StreamRemoveRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup

	for i, r := range reqs.GetRequests() {
		idx, req := i, r
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + req.GetId().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiRemoveRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			loc, err := s.Remove(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.RemoveRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetId().GetId(),
						ServingData: errdetails.Serialize(req),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Locations[idx] = loc
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiRemoveRPCName+" gRPC error response",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) GetObject(ctx context.Context, req *payload.Object_VectorRequest) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (interface{}, error) {
		vec, err = vc.GetObject(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return vec, nil
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetId().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.GetObjectRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.GetObjectRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
			)
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInternal(
				vald.GetObjectRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.GetObjectRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.GetObjectRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return vec, nil
}

func (s *server) getObjects(ctx context.Context, req *payload.Object_VectorRequest) (vecs *sync.Map, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "getObjects"), apiName+"/"+vald.GetObjectRPCName+"/getObjects")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var errs error
	var emu sync.Mutex
	vecs = new(sync.Map)
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.GetObjectRPCName+"/getObjects/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		vec, err := vc.GetObject(ctx, req, copts...)
		if err != nil {
			reqInfo := &errdetails.RequestInfo{
				RequestId:   req.GetId().GetId(),
				ServingData: errdetails.Serialize(req),
			}
			resInfo := &errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
			}
			var attrs trace.Attributes
			var code codes.Code

			switch {
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(
					vald.GetObjectRPCName+" API canceld", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeCancelled(
					errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
				)
				code = codes.Canceled
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(
					vald.GetObjectRPCName+" API deadline exceeded", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeDeadlineExceeded(
					errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
				)
				code = codes.DeadlineExceeded
			case errors.Is(err, errors.ErrTargetNotFound):
				err = status.WrapWithInternal(
					vald.GetObjectRPCName+" API target not found", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeInternal(err.Error())
				code = codes.Internal
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(
					vald.GetObjectRPCName+" API connection not found", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeInternal(err.Error())
				code = codes.Internal
			default:
				var (
					st  *status.Status
					msg string
				)
				st, msg, err = status.ParseError(err, codes.Internal,
					"failed to parse "+vald.GetObjectRPCName+" gRPC error response", reqInfo, resInfo,
				)
				attrs = trace.FromGRPCStatus(st.Code(), msg)
				code = st.Code()
			}
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if code == codes.NotFound {
				return nil
			}
			emu.Lock()
			if errs == nil {
				errs = err
			} else {
				errs = errors.Join(errs, err)
			}
			emu.Unlock()
			return err
		}
		vecs.Store(target, vec)
		return nil
	})
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.GetObjectRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetId().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName + ".BroadCast",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		errs = errors.Join(errs, err)
	}
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse "+vald.GetObjectRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetId().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName + "." + "BroadCast",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return vecs, nil
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamGetObjectRPCName), apiName+"/"+vald.StreamGetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Object_VectorRequest) (*payload.Object_StreamVector, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamInsertRPCName+"/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.GetObject(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.GetObjectRPCName+" gRPC error response")
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				return &payload.Object_StreamVector{
					Payload: &payload.Object_StreamVector_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamVector{
				Payload: &payload.Object_StreamVector_Vector{
					Vector: res,
				},
			}, nil
		},
	)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.StreamGetObjectRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}
