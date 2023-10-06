// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"reflect"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
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
	vald.UnimplementedValdServerWithMirror
}

const apiName = "vald/gateway/mirror"

func New(opts ...Option) (vald.ServerWithMirror, error) {
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
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.MirrorRPCServiceName+"/"+vald.RegisterRPCName), apiName+"/"+vald.RegisterRPCName)
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
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RegisterRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.RegisterRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithCanceled(
				vald.RegisterRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.RegisterRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInvalidArgument(
				vald.RegisterRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RegisterRPCName+" gRPC error response", reqInfo, resInfo,
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

	// Get own address and the addresses of other mirror gateways to which this gateway is currently connected.
	tgts, err := s.mirror.MirrorTargets()
	if err != nil {
		err = status.WrapWithInternal(vald.RegisterRPCName+" API failed to get connected vald gateway targets", err,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "mirror gateway targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RegisterRPCName,
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

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
		id, err = vc.Exists(ctx, meta, copts...)
		return id, err
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.ExistsRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.Search(ctx, req, copts...)
		return res, err
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.SearchRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.SearchByID(ctx, req, copts...)
		return res, err
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.SearchByIDRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.MultiSearch(ctx, req, copts...)
		return res, err
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.MultiSearchRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.MultiSearchByID(ctx, req, copts...)
		return res, err
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.MultiSearchByIDRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.LinearSearch(ctx, req, copts...)
		return res, err
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.LinearSearchRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.LinearSearchByID(ctx, req, copts...)
		return res, err
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.LinearSearchByIDRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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
	ctx, span := trace.StartSpan(
		grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamLinearSearchByIDRPCName),
		apiName+"/"+vald.StreamLinearSearchByIDRPCName,
	)
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

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.MultiLinearSearch(ctx, req, copts...)
		return res, err
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.MultiLinearSearchRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
		res, err = vc.MultiLinearSearchByID(ctx, req, copts...)
		return res, err
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.MultiLinearSearchByIDRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (loc *payload.Object_Location, err error) {
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
		loc, err = s.doInsert(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			_, derr := s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
				loc, err = vc.Insert(ctx, req, copts...)
				return loc, err
			})
			return loc, errors.Join(derr, err)
		})
		if err != nil {
			reqInfo := &errdetails.RequestInfo{
				RequestId: req.GetVector().GetId(),
			}
			resInfo := &errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
			}
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
		log.Debugf("Insert API succeeded to %#v", loc)
		return loc, nil
	}

	var mu sync.Mutex
	var result sync.Map[string, *errorState] // map[target: errorState]
	loc = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		code := codes.OK
		ce, err := s.doInsert(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Insert(ctx, req, copts...)
		})
		if err != nil {
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.InsertRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId: req.GetVector().GetId(),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			code = st.Code()
		}
		if err == nil && ce != nil {
			mu.Lock()
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
			mu.Unlock()
		}
		result.Store(target, &errorState{err, code})
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId: req.GetVector().GetId(),
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

	alreadyExistsTgts := make([]string, 0, result.Len()/2)
	successTgts := make([]string, 0, result.Len()/2)
	result.Range(func(target string, es *errorState) bool {
		switch {
		case es.err == nil:
			successTgts = append(successTgts, target)
		case es.code == codes.AlreadyExists:
			alreadyExistsTgts = append(alreadyExistsTgts, target)
			err = errors.Join(err, es.err)
		default:
			err = errors.Join(es.err, err)
		}
		return true
	})
	if err == nil {
		log.Debugf(vald.InsertRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}

	reqInfo := &errdetails.RequestInfo{
		RequestId: req.GetVector().GetId(),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
	}

	switch {
	case result.Len() == len(alreadyExistsTgts):
		err = status.WrapWithAlreadyExists(vald.InsertRPCName+" API target same vector already exists", err, reqInfo, resInfo)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	case result.Len() > len(successTgts)+len(alreadyExistsTgts): // Contains errors except for ALREADY_EXIST.
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

	// In this case, the status code in the result object contains only OK or ALREADY_EXIST.
	// And send Update API requst to ALREADY_EXIST cluster using query requested by the user.
	log.Warnf("failed to "+vald.InsertRPCName+" API: %#v", err)

	updateReq := &payload.Update_Request{
		Vector: req.GetVector(),
		Config: &payload.Update_Config{
			Timestamp: req.GetConfig().GetTimestamp(),
		},
	}
	err = s.gateway.DoMulti(ctx, alreadyExistsTgts, func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "DoMulti/"+target), apiName+"/"+vald.UpdateRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		code := codes.OK
		ce, err := s.doUpdate(ctx, updateReq, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Update(ctx, updateReq, copts...)
		})
		if err != nil {
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpdateRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId: req.GetVector().GetId(),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + ".DoMulti/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			code = st.Code()
		}
		if err == nil && ce != nil {
			mu.Lock()
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
			mu.Unlock()
		}
		result.Store(target, &errorState{err, code})
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId: updateReq.GetVector().GetId(),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.InsertRPCName + ".DoMulti",
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpdateRPCName+" for "+vald.InsertRPCName+" API connection not found", err, reqInfo, resInfo,
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
			"failed to parse "+vald.UpdateRPCName+" for "+vald.InsertRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	alreadyExistsTgts = alreadyExistsTgts[0:0]
	successTgts = successTgts[0:0]
	result.Range(func(target string, es *errorState) bool {
		switch {
		case es.err == nil:
			successTgts = append(successTgts, target)
		case es.code == codes.AlreadyExists:
			alreadyExistsTgts = append(alreadyExistsTgts, target)
			err = errors.Join(err, es.err)
		default:
			err = errors.Join(es.err, err)
		}
		return true
	})
	if err == nil {
		log.Debugf(vald.UpdateRPCName+"for "+vald.InsertRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}

	reqInfo = &errdetails.RequestInfo{
		RequestId: req.GetVector().GetId(),
	}
	resInfo = &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + "." + vald.UpdateRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
	}

	switch {
	case result.Len() == len(alreadyExistsTgts):
		err = status.WrapWithAlreadyExists(vald.UpdateRPCName+" for "+vald.InsertRPCName+" API target same vector already exists", err, reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	case result.Len() > len(successTgts)+len(alreadyExistsTgts): // Contains errors except for NOT_FOUND and ALREADY_EXIST.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpdateRPCName+" for "+vald.InsertRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	log.Debugf(vald.UpdateRPCName+"for "+vald.InsertRPCName+" API request succeeded to %#v, err: %v", loc, err)
	return loc, nil
}

func (s *server) doInsert(ctx context.Context, req *payload.Insert_Request, f func(ctx context.Context) (*payload.Object_Location, error)) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doInsert"), apiName+"/doInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err = f(ctx)
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.InsertRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

	var mu, emu sync.Mutex
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
				emu.Lock()
				errs = errors.Join(errs, err)
				emu.Unlock()
				return nil
			}
			mu.Lock()
			res.Locations[idx] = loc
			mu.Unlock()
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
		loc, err = s.doUpdate(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			_, derr := s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
				loc, err = vc.Update(ctx, req, copts...)
				return loc, err
			})
			return loc, errors.Join(derr, err)
		})
		if err != nil {
			reqInfo := &errdetails.RequestInfo{
				RequestId: req.GetVector().GetId(),
			}
			resInfo := &errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
			}
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
		log.Debugf("Update API succeeded to %#v", loc)
		return loc, nil
	}

	var mu sync.Mutex
	var result sync.Map[string, *errorState]
	loc = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpdateRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		code := codes.OK
		ce, err := s.doUpdate(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Update(ctx, req, copts...)
		})
		if err != nil {
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpdateRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			code = st.Code()
		}
		if err == nil && ce != nil {
			mu.Lock()
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
			mu.Unlock()
		}
		result.Store(target, &errorState{err, code})
		return err
	})
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

	var alreadyExistsCnt int
	notFoundTgts := make([]string, 0, result.Len()/2)
	successTgts := make([]string, 0, result.Len()/2)
	result.Range(func(target string, es *errorState) bool {
		switch {
		case es.err == nil:
			successTgts = append(successTgts, target)
		case es.code == codes.AlreadyExists:
			alreadyExistsCnt++
			err = errors.Join(err, es.err)
		case es.code == codes.NotFound:
			notFoundTgts = append(notFoundTgts, target)
			err = errors.Join(err, es.err)
		default:
			err = errors.Join(es.err, err)
		}
		return true
	})
	if err == nil || (len(successTgts) > 0 && result.Len() == len(successTgts)+alreadyExistsCnt) {
		log.Debugf(vald.UpdateRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}

	reqInfo := &errdetails.RequestInfo{
		RequestId: req.GetVector().GetId(),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
	}

	switch {
	case result.Len() == len(notFoundTgts):
		err = status.WrapWithNotFound(vald.UpdateRPCName+" API target not found", err, reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	case result.Len() == alreadyExistsCnt:
		err = status.WrapWithAlreadyExists(vald.UpdateRPCName+" API target same vector already exists", err, reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	case result.Len() > len(successTgts)+len(notFoundTgts)+alreadyExistsCnt: // Contains errors except for NOT_FOUND and ALREADY_EXIST.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	// In this case, the status code in the result object contains only OK or ALREADY_EXIST.
	// And send Insert API requst to NOT_FOUND cluster using query requested by the user.
	log.Warnf("failed to "+vald.UpdateRPCName+" API: %#v", err)

	insReq := &payload.Insert_Request{
		Vector: req.GetVector(),
		Config: &payload.Insert_Config{
			Timestamp: req.GetConfig().GetTimestamp(),
		},
	}
	err = s.gateway.DoMulti(ctx, notFoundTgts, func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		code := codes.OK
		ce, err := s.doInsert(ctx, insReq, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Insert(ctx, insReq, copts...)
		})
		if err != nil {
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
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
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			code = st.Code()
		}
		if err == nil && ce != nil {
			mu.Lock()
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
			mu.Unlock()
		}
		result.Store(target, &errorState{err, code})
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.InsertRPCName + ".BroadCast",
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.InsertRPCName+" for "+vald.UpdateRPCName+" API connection not found", err, reqInfo, resInfo,
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
			"failed to parse "+vald.InsertRPCName+" for "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	alreadyExistsCnt = 0
	notFoundTgts = notFoundTgts[0:0]
	successTgts = successTgts[0:0]
	result.Range(func(target string, em *errorState) bool {
		switch {
		case em.err == nil:
			successTgts = append(successTgts, target)
		case em.code == codes.AlreadyExists:
			alreadyExistsCnt++
			err = errors.Join(err, em.err)
		case em.code == codes.NotFound:
			notFoundTgts = append(notFoundTgts, target)
			err = errors.Join(err, em.err)
		default:
			err = errors.Join(em.err, err)
		}
		return true
	})
	if err == nil || (len(successTgts) > 0 && result.Len() == len(successTgts)+alreadyExistsCnt) {
		log.Debugf(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}

	reqInfo = &errdetails.RequestInfo{
		RequestId:   req.GetVector().GetId(),
		ServingData: errdetails.Serialize(req),
	}
	resInfo = &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.InsertRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
	}

	switch {
	case result.Len() == len(notFoundTgts):
		err = status.WrapWithNotFound(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API target not found", err, reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	case result.Len() == alreadyExistsCnt:
		err = status.WrapWithAlreadyExists(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API target same vector already exists", err, reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	case result.Len() > len(successTgts)+len(notFoundTgts)+alreadyExistsCnt: // Contains errors except for NOT_FOUND and ALREADY_EXIST.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.InsertRPCName+" for "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	log.Debugf(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API request succeeded to %#v, err: %v", loc, err)
	return loc, nil
}

func (s *server) doUpdate(ctx context.Context, req *payload.Update_Request, f func(ctx context.Context) (*payload.Object_Location, error)) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doUpdate"), apiName+"/doUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err = f(ctx)
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.UpdateRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

	var mu, emu sync.Mutex
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
				emu.Lock()
				errs = errors.Join(errs, err)
				emu.Unlock()
				return nil
			}
			mu.Lock()
			res.Locations[idx] = loc
			mu.Unlock()
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
		loc, err = s.doUpsert(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
				loc, err = vc.Upsert(ctx, req, copts...)
				return loc, err
			})
			return loc, err
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
		log.Debugf("Upsert API succeeded to %#v", loc)
		return loc, nil
	}

	var mu sync.Mutex
	var result sync.Map[string, *errorState]
	loc = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpsertRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		code := codes.OK
		ce, err := s.doUpsert(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Upsert(ctx, req, copts...)
		})
		if err != nil {
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
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
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			code = st.Code()
		}
		if err == nil && ce != nil {
			mu.Lock()
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
			mu.Unlock()
		}
		result.Store(target, &errorState{err, code})
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

	var alreadyExistsCnt int
	successTgts := make([]string, 0, result.Len()/2)
	result.Range(func(target string, es *errorState) bool {
		switch {
		case es.err == nil:
			successTgts = append(successTgts, target)
		case es.code == codes.AlreadyExists:
			alreadyExistsCnt++
			err = errors.Join(err, es.err)
		default:
			err = errors.Join(es.err, err)
		}
		return true
	})
	if err == nil || (len(successTgts) > 0 && result.Len() == len(successTgts)+alreadyExistsCnt) {
		log.Debugf(vald.UpsertRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}
	reqInfo := &errdetails.RequestInfo{
		RequestId:   req.GetVector().GetId(),
		ServingData: errdetails.Serialize(req),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
	}

	switch {
	case result.Len() == alreadyExistsCnt:
		err = status.WrapWithAlreadyExists(vald.UpsertRPCName+" API target same vector already exists", err, reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	default:
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpsertRPCName+" gRPC error response", reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
}

func (s *server) doUpsert(ctx context.Context, req *payload.Upsert_Request, f func(ctx context.Context) (*payload.Object_Location, error)) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doUpsert"), apiName+"/doUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err = f(ctx)
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.UpsertRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

	var mu, emu sync.Mutex
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
				emu.Lock()
				errs = errors.Join(errs, err)
				emu.Unlock()
				return nil
			}
			mu.Lock()
			res.Locations[idx] = loc
			mu.Unlock()
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
		loc, err = s.doRemove(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
				loc, err = vc.Remove(ctx, req, copts...)
				return loc, err
			})
			return loc, err
		})
		if err != nil {
			reqInfo := &errdetails.RequestInfo{
				RequestId: req.GetId().GetId(),
			}
			resInfo := &errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
			}
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
		log.Debugf("Remove API remove succeeded to %#v", loc)
		return loc, nil
	}

	var mu sync.Mutex
	var result sync.Map[string, *errorState]
	loc = &payload.Object_Location{
		Uuid: req.GetId().GetId(),
		Ips:  make([]string, 0),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.RemoveRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		code := codes.OK
		ce, err := s.doRemove(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Remove(ctx, req, copts...)
		})
		if err != nil {
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId: req.GetId().GetId(),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			code = st.Code()
		}
		if err == nil && ce != nil {
			mu.Lock()
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
			mu.Unlock()
		}
		result.Store(target, &errorState{err, code})
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId: req.GetId().GetId(),
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

	var notFoundCnt int
	successTgts := make([]string, 0, result.Len()/2)
	result.Range(func(target string, es *errorState) bool {
		switch {
		case es.err == nil:
			successTgts = append(successTgts, target)
		case es.code == codes.NotFound:
			notFoundCnt++
			err = errors.Join(err, es.err)
		default:
			err = errors.Join(es.err, err)
		}
		return true
	})
	if err == nil || (len(successTgts) > 0 && result.Len() == len(successTgts)+notFoundCnt) {
		log.Debugf(vald.RemoveRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}

	reqInfo := &errdetails.RequestInfo{
		RequestId: req.GetId().GetId(),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
	}

	switch {
	case result.Len() == notFoundCnt:
		err = status.WrapWithNotFound(vald.RemoveRPCName+" API id "+req.GetId().GetId()+" not found", err, reqInfo, resInfo)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	default:
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveRPCName+" gRPC error response", reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
}

func (s *server) doRemove(ctx context.Context, req *payload.Remove_Request, f func(ctx context.Context) (*payload.Object_Location, error)) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doRemove"), apiName+"/doRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err = f(ctx)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId: req.GetId().GetId(),
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.RemoveRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

	var mu, emu sync.Mutex
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
				emu.Lock()
				errs = errors.Join(errs, err)
				emu.Unlock()
				return nil
			}
			mu.Lock()
			res.Locations[idx] = loc
			mu.Unlock()
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

func (s *server) RemoveByTimestamp(ctx context.Context, req *payload.Remove_TimestampRequest) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveByTimestampRPCName), apiName+"/"+vald.RemoveByTimestampRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	reqSrcPodName := s.gateway.FromForwardedContext(ctx)

	// When this condition is matched, the request is proxied to another Mirror gateway.
	// So this component sends requests only to the Vald gateway (LB gateway) of its own cluster.
	if len(reqSrcPodName) != 0 {
		locs, err = s.doRemoveByTimestamp(ctx, req, func(ctx context.Context) (*payload.Object_Locations, error) {
			_, derr := s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
				locs, err = vc.RemoveByTimestamp(ctx, req, copts...)
				return locs, err
			})
			return locs, errors.Join(derr, err)
		})
		if err != nil {
			reqInfo := &errdetails.RequestInfo{
				ServingData: errdetails.Serialize(req),
			}
			resInfo := &errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveByTimestampRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
			}
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
		log.Debugf("RemoveByTimestamp API remove succeeded to %#v", locs)
		return locs, nil
	}

	var mu sync.Mutex
	var result sync.Map[string, *errorState]
	locs = new(payload.Object_Locations)

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.RemoveByTimestampRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		code := codes.OK
		res, err := s.doRemoveByTimestamp(ctx, req, func(ctx context.Context) (*payload.Object_Locations, error) {
			return vc.RemoveByTimestamp(ctx, req, copts...)
		})
		if err != nil {
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveByTimestampRPCName+" gRPC error response",
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveByTimestampRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			code = st.Code()
		}
		if err == nil && res != nil {
			mu.Lock()
			locs.Locations = append(locs.Locations, res.GetLocations()...)
			mu.Unlock()
		}
		result.Store(target, &errorState{err, code})
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveByTimestampRPCName + ".BroadCast",
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.RemoveByTimestampRPCName+" API connection not found", err, reqInfo, resInfo,
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
			"failed to parse "+vald.RemoveByTimestampRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	var notFoundCnt int
	successTgts := make([]string, 0, result.Len()/2)
	result.Range(func(target string, es *errorState) bool {
		switch {
		case es.err == nil:
			successTgts = append(successTgts, target)
		case es.code == codes.NotFound:
			notFoundCnt++
			err = errors.Join(err, es.err)
		default:
			err = errors.Join(es.err, err)
		}
		return true
	})
	if err == nil || (len(successTgts) > 0 && result.Len() == len(successTgts)+notFoundCnt) {
		log.Debugf(vald.RemoveByTimestampRPCName+" API request succeeded to %#v", locs)
		return locs, nil
	}

	reqInfo := &errdetails.RequestInfo{
		ServingData: errdetails.Serialize(req),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveByTimestampRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
	}

	switch {
	case result.Len() == notFoundCnt:
		err = status.WrapWithNotFound(vald.RemoveByTimestampRPCName+" API target not found", err, reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	default:
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveByTimestampRPCName+" gRPC error response", reqInfo, resInfo)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
}

func (s *server) doRemoveByTimestamp(ctx context.Context, req *payload.Remove_TimestampRequest, f func(ctx context.Context) (*payload.Object_Locations, error)) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doRemoveByTimestamp"), apiName+"/doRemoveByTimestamp")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	locs, err = f(ctx)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveByTimestampRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.RemoveByTimestampRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.RemoveByTimestampRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInternal(
				vald.RemoveByTimestampRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.RemoveByTimestampRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveByTimestampRPCName+" gRPC error response", reqInfo, resInfo,
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
	return locs, nil
}

func (s *server) GetObject(ctx context.Context, req *payload.Object_VectorRequest) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error) {
		vec, err = vc.GetObject(ctx, req, copts...)
		return vec, err
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
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.GetObjectRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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

type errorState struct {
	err  error
	code codes.Code
}
