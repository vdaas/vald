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
	"reflect"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

type Server interface {
	vald.Server
	mirror.MirrorServer
}

type server struct {
	vald.UnimplementedValdServer
	mirror.UnimplementedMirrorServer
	eg                errgroup.Group
	gateway           service.Gateway
	mirror            service.Mirror
	vAddr             string
	name              string
	ip                string
	streamConcurrency int
}

const (
	apiName = "vald/gateway/mirror"

	// targetCapHintDivisor sizes the per-outcome target slices built from a
	// broadcast result at a fraction of the result count, as a capacity hint.
	targetCapHintDivisor = 2
)

// resourceInfo builds the errdetails.ResourceInfo shared by every handler error
// path: the resourceType is caller-specific, the resource name is always this
// server's apiName/name/ip triple.
func (s *server) resourceInfo(resourceType string) *errdetails.ResourceInfo {
	return &errdetails.ResourceInfo{
		ResourceType: resourceType,
		ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
	}
}

// broadCastNoTargetError builds the error for a BroadCast that returned nil
// without invoking its callback for any target. RangeConcurrent skips a target
// whose pool has no usable conn or whose circuit breaker is open and drops that
// error, so an empty result means the request reached no cluster at all and
// must not be reported as OK.
func (s *server) broadCastNoTargetError(rpcName string, reqInfo *errdetails.RequestInfo) error {
	return status.WrapWithUnavailable(
		rpcName+" API no reachable target", errors.ErrGRPCClientConnNotFound("*"),
		reqInfo, s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix+"/vald.v1."+rpcName+".BroadCast"),
	)
}

// vAddrResourceInfo builds the errdetails.ResourceInfo for an RPC that is
// forwarded to this cluster's Vald gateway (LB gateway) at s.vAddr.
func (s *server) vAddrResourceInfo(rpcName string) *errdetails.ResourceInfo {
	return &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + rpcName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, s.vAddr),
	}
}

// wrapTargetError converts the error of a call to a single target into the
// gRPC status error returned for rpcName, logs it and records it on span.
func wrapTargetError(
	span trace.Span,
	rpcName string,
	err error,
	reqInfo *errdetails.RequestInfo,
	resInfo *errdetails.ResourceInfo,
) error {
	var attrs trace.Attributes
	switch {
	case errors.Is(err, context.Canceled):
		err = status.WrapWithCanceled(
			rpcName+" API canceld", err, reqInfo, resInfo,
		)
		attrs = trace.StatusCodeCancelled(err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		err = status.WrapWithDeadlineExceeded(
			rpcName+" API deadline exceeded", err, reqInfo, resInfo,
		)
		attrs = trace.StatusCodeDeadlineExceeded(err.Error())
	case errors.Is(err, errors.ErrTargetNotFound):
		err = status.WrapWithInternal(
			rpcName+" API target not found", err, reqInfo, resInfo,
		)
		attrs = trace.StatusCodeInternal(err.Error())
	case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
		err = status.WrapWithInternal(
			rpcName+" API connection not found", err, reqInfo, resInfo,
		)
		attrs = trace.StatusCodeInternal(err.Error())
	default:
		var (
			st  *status.Status
			msg string
		)
		st, msg, err = status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response", reqInfo, resInfo,
		)
		attrs = trace.FromGRPCStatus(st.Code(), msg)
	}
	log.Warn(err)
	errhandler.RecordSpanAttrs(span, attrs, err)
	return err
}

// statusOf returns the gRPC status of err. When err carries no status (or an
// empty message) it logs that and falls back to an Unknown status for rpcName,
// so callers can always use st.Code() and st.Proto() without a nil check.
func statusOf(err error, rpcName string) *status.Status {
	st, _ := status.FromError(err)
	if st == nil || st.Message() == "" {
		log.Errorf("gRPC call returned not a gRPC status error: %v", err)
		st = status.New(codes.Unknown, "failed to parse "+rpcName+" gRPC error response")
	}
	return st
}

// New returns a Vald server as gRPC handler with mirror using the provided options.
func New(opts ...Option) (Server, error) {
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

// Register handles the registration of mirror targets.
// The function connects to the mirror using the provided targets, and if successful,
// returns the addresses of connected Mirror gateways.
func (s *server) Register(
	ctx context.Context, req *payload.Mirror_Targets,
) (*payload.Mirror_Targets, error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+mirror.RPCServiceName+"/"+mirror.RegisterRPCName), apiName+"/"+mirror.RegisterRPCName)
	defer trace.End(span)
	err := s.mirror.Connect(ctx, req.GetTargets()...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + mirror.RegisterRPCName)
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				mirror.RegisterRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				mirror.RegisterRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
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
			err = status.WrapWithInternal(
				mirror.RegisterRPCName+" API failed to connect mirror gateway targets", err, reqInfo, resInfo,
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "mirror gateway targets",
							Description: err.Error(),
						},
					},
				},
			)
			attrs = trace.StatusCodeInternal(err.Error())
		}
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, err
	}

	// Get own address and the addresses of other mirror gateways to which this gateway is currently connected.
	tgts, err := s.mirror.MirrorTargets(ctx)
	if err != nil {
		err = status.WrapWithInternal(mirror.RegisterRPCName+" API failed to get connected mirror gateway targets", err,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "mirror gateway targets",
						Description: err.Error(),
					},
				},
			},
			s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix+"/vald.v1."+mirror.RegisterRPCName),
		)
		log.Warn(err)
		return errhandler.HandleError[payload.Mirror_Targets](span, codes.Internal, err)
	}
	return &payload.Mirror_Targets{
		Targets: tgts,
	}, nil
}

// Exists bypasses the incoming Exist request to Vald gateway (LB gateway) in its own cluster.
func (s *server) Exists(
	ctx context.Context, meta *payload.Object_ID,
) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.ExistsRPCName), apiName+"/"+vald.ExistsRPCName)
	defer trace.End(span)

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
		id, err = vc.Exists(ctx, meta, copts...)
		return id, err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(meta),
		}
		resInfo := s.vAddrResourceInfo(vald.ExistsRPCName)
		return nil, wrapTargetError(span, vald.ExistsRPCName, err, reqInfo, resInfo)
	}
	return id, nil
}

// Search bypasses the incoming Search request to Vald gateway (LB gateway) in its own cluster.
func (s *server) Search(
	ctx context.Context, req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.SearchRPCName), apiName+"/"+vald.SearchRPCName)
	defer trace.End(span)

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
		res, err = vc.Search(ctx, req, copts...)
		return res, err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetConfig().GetRequestId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.vAddrResourceInfo(vald.SearchRPCName)
		return nil, wrapTargetError(span, vald.SearchRPCName, err, reqInfo, resInfo)
	}
	return res, nil
}

// SearchByID bypasses the incoming SearchByID request to Vald gateway (LB gateway) in its own cluster.
func (s *server) SearchByID(
	ctx context.Context, req *payload.Search_IDRequest,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.SearchByIDRPCName), apiName+"/"+vald.SearchByIDRPCName)
	defer trace.End(span)

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
		res, err = vc.SearchByID(ctx, req, copts...)
		return res, err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetConfig().GetRequestId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.vAddrResourceInfo(vald.SearchByIDRPCName)
		return nil, wrapTargetError(span, vald.SearchByIDRPCName, err, reqInfo, resInfo)
	}
	return res, nil
}

// StreamSearch bypasses it as a Search request to Vald gateway (LB gateway) in its own cluster.
func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamSearchRPCName), apiName+"/"+vald.StreamSearchRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_Request) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer trace.End(sspan)
			res, err := s.Search(ctx, req)
			if err != nil {
				st := statusOf(err, vald.SearchRPCName)
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
		},
	)
	if err != nil {
		st := statusOf(err, vald.StreamSearchRPCName)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

// StreamSearchByID bypasses it as a SearchByID request to Vald gateway (LB gateway) in its own cluster.
func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamSearchByIDRPCName), apiName+"/"+vald.StreamSearchByIDRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"."+vald.StreamSearchByIDRPCName+"/id-"+req.GetId())
			defer trace.End(sspan)
			res, err := s.SearchByID(ctx, req)
			if err != nil {
				st := statusOf(err, vald.SearchByIDRPCName)
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
		},
	)
	if err != nil {
		st := statusOf(err, vald.StreamSearchByIDRPCName)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

// MultiSearch bypasses the incoming MultiSearch request to Vald gateway (LB gateway) in its own cluster.
func (s *server) MultiSearch(
	ctx context.Context, req *payload.Search_MultiRequest,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiSearchRPCName), apiName+"/"+vald.MultiSearchRPCName)
	defer trace.End(span)

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
		res, err = vc.MultiSearch(ctx, req, copts...)
		return res, err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.vAddrResourceInfo(vald.MultiSearchRPCName)
		return nil, wrapTargetError(span, vald.MultiSearchRPCName, err, reqInfo, resInfo)
	}
	return res, nil
}

// MultiSearchByID bypasses the incoming MultiSearchByID request to Vald gateway (LB gateway) in its own cluster.
func (s *server) MultiSearchByID(
	ctx context.Context, req *payload.Search_MultiIDRequest,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiSearchByIDRPCName), apiName+"/"+vald.MultiSearchByIDRPCName)
	defer trace.End(span)

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
		res, err = vc.MultiSearchByID(ctx, req, copts...)
		return res, err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.vAddrResourceInfo(vald.MultiSearchByIDRPCName)
		return nil, wrapTargetError(span, vald.MultiSearchByIDRPCName, err, reqInfo, resInfo)
	}
	return res, nil
}

// LinearSearch bypasses the incoming LinearSearch request to Vald gateway (LB gateway) in its own cluster.
func (s *server) LinearSearch(
	ctx context.Context, req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.LinearSearchRPCName), apiName+"/"+vald.LinearSearchRPCName)
	defer trace.End(span)

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
		res, err = vc.LinearSearch(ctx, req, copts...)
		return res, err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetConfig().GetRequestId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.vAddrResourceInfo(vald.LinearSearchRPCName)
		return nil, wrapTargetError(span, vald.LinearSearchRPCName, err, reqInfo, resInfo)
	}
	return res, nil
}

// LinearSearchByID bypasses the incoming LinearSearchByID request to Vald gateway (LB gateway) in its own cluster.
func (s *server) LinearSearchByID(
	ctx context.Context, req *payload.Search_IDRequest,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.LinearSearchByIDRPCName), apiName+"/"+vald.LinearSearchByIDRPCName)
	defer trace.End(span)

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
		res, err = vc.LinearSearchByID(ctx, req, copts...)
		return res, err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetConfig().GetRequestId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.vAddrResourceInfo(vald.LinearSearchByIDRPCName)
		return nil, wrapTargetError(span, vald.LinearSearchByIDRPCName, err, reqInfo, resInfo)
	}
	return res, nil
}

// StreamLinearSearch bypasses it as a LinearSearch request to Vald gateway (LB gateway) in its own cluster.
func (s *server) StreamLinearSearch(stream vald.Search_StreamLinearSearchServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamLinearSearchRPCName), apiName+"/"+vald.StreamLinearSearchRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_Request) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamLinearSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer trace.End(sspan)
			res, err := s.LinearSearch(ctx, req)
			if err != nil {
				st := statusOf(err, vald.LinearSearchRPCName)
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
		},
	)
	if err != nil {
		st := statusOf(err, vald.StreamLinearSearchRPCName)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

// StreamLinearSearchByID bypasses it as a LinearSearchByID request to Vald gateway (LB gateway) in its own cluster.
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
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"."+vald.StreamLinearSearchByIDRPCName+"/id-"+req.GetId())
			defer trace.End(sspan)
			res, err := s.LinearSearchByID(ctx, req)
			if err != nil {
				st := statusOf(err, vald.LinearSearchByIDRPCName)
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
		},
	)
	if err != nil {
		st := statusOf(err, vald.StreamLinearSearchByIDRPCName)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

// MultiLinearSearch bypasses the incoming MultiLinearSearch request to Vald gateway (LB gateway) in its own cluster.
func (s *server) MultiLinearSearch(
	ctx context.Context, req *payload.Search_MultiRequest,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiLinearSearchRPCName), apiName+"/"+vald.MultiLinearSearchRPCName)
	defer trace.End(span)

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
		res, err = vc.MultiLinearSearch(ctx, req, copts...)
		return res, err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.vAddrResourceInfo(vald.MultiLinearSearchRPCName)
		return nil, wrapTargetError(span, vald.MultiLinearSearchRPCName, err, reqInfo, resInfo)
	}
	return res, nil
}

// MultiLinearSearchByID bypasses the incoming MultiLinearSearchByID request to Vald gateway (LB gateway) in its own cluster.
func (s *server) MultiLinearSearchByID(
	ctx context.Context, req *payload.Search_MultiIDRequest,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiLinearSearchByIDRPCName), apiName+"/"+vald.MultiLinearSearchByIDRPCName)
	defer trace.End(span)

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
		res, err = vc.MultiLinearSearchByID(ctx, req, copts...)
		return res, err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.vAddrResourceInfo(vald.MultiLinearSearchByIDRPCName)
		return nil, wrapTargetError(span, vald.MultiLinearSearchByIDRPCName, err, reqInfo, resInfo)
	}
	return res, nil
}

// Insert handles the insertion of an object with the given request.
// If the request is proxied from another Mirror gateway, the request is forwarded to the Vald gateway (LB gateway) of its own cluster.
// If the request is from a user, it is sent to other Mirror gateways and the Vald gateway (LB gateway) of its own cluster.
// The result is a location of the inserted object.
func (s *server) Insert(
	ctx context.Context, req *payload.Insert_Request,
) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.InsertRPCName), apiName+"/"+vald.InsertRPCName)
	defer trace.End(span)

	// If this condition is matched, it means that the request was proxied from another Mirror gateway.
	// So this component sends requests only to the Vald gateway (LB gateway) of its own cluster.
	if s.isProxied(ctx) {
		loc, err = s.doInsert(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			_, derr := s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
				loc, err = vc.Insert(ctx, req, copts...)
				return loc, err
			})
			return loc, errors.Join(derr, err)
		})
		if err != nil {
			st := statusOf(err, vald.InsertRPCName)
			errhandler.RecordSpanStatus(span, st, err)
			log.Warn(err)
			return nil, err
		}
		log.Debugf("Insert API succeeded to %#v", loc)
		return loc, nil
	}

	// If this condition is matched, it means the request from user.
	// So this component sends requests to other Mirror gateways and the Vald gateway (LB gateway) of its own cluster.
	return s.handleInsert(ctx, req)
}

func (s *server) handleInsert(
	ctx context.Context, req *payload.Insert_Request,
) (loc *payload.Object_Location, err error) { // skipcq: GO-R1005
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, "handleInsert"), apiName+"/handleInsert")
	defer trace.End(span)

	var mu sync.Mutex
	result := make(map[string]*errorState) // map[target host: error state]
	loc = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
		defer trace.End(span)

		code := codes.OK
		ce, err := s.doInsert(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Insert(ctx, req, copts...)
		})
		if err != nil {
			st := statusOf(err, vald.InsertRPCName)
			log.Warn(err)
			errhandler.RecordSpanStatus(span, st, err)
			code = st.Code()
		}
		mu.Lock()
		if err == nil && ce != nil {
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
		}
		result[target] = &errorState{err, code}
		mu.Unlock()
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId: req.GetVector().GetId(),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + ".BroadCast")
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.InsertRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			return errhandler.HandleError[payload.Object_Location](span, codes.Internal, err)
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.InsertRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}
	if len(result) == 0 {
		err = s.broadCastNoTargetError(vald.InsertRPCName, &errdetails.RequestInfo{
			RequestId: req.GetVector().GetId(),
		})
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Location](span, codes.Unavailable, err)
	}

	alreadyExistsTgts := make([]string, 0, len(result)/targetCapHintDivisor)
	successTgts := make([]string, 0, len(result)/targetCapHintDivisor)
	for target, es := range result {
		switch {
		case es.err == nil:
			successTgts = append(successTgts, target)
		case es.code == codes.AlreadyExists:
			alreadyExistsTgts = append(alreadyExistsTgts, target)
			err = errors.Join(err, es.err)
		default:
			err = errors.Join(es.err, err)
		}
	}
	if err == nil {
		log.Debugf(vald.InsertRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}

	reqInfo := &errdetails.RequestInfo{
		RequestId: req.GetVector().GetId(),
	}
	resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName)

	switch {
	case len(result) == len(alreadyExistsTgts):
		err = status.WrapWithAlreadyExists(vald.InsertRPCName+" API target same vector already exists", err, reqInfo, resInfo)
		return errhandler.HandleError[payload.Object_Location](span, codes.AlreadyExists, err)
	case len(result) > len(successTgts)+len(alreadyExistsTgts): // Contains errors except for ALREADY_EXIST.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.InsertRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}

	// In this case, the status code in the result object contains only OK or ALREADY_EXIST.
	// And send Update API request to ALREADY_EXIST cluster using the query requested by the user.
	log.Warnf("failed to "+vald.InsertRPCName+" API: %#v", err)

	resLoc, err := s.handleInsertResult(ctx, alreadyExistsTgts, &payload.Update_Request{
		Vector: req.GetVector(),
		Config: &payload.Update_Config{
			Timestamp: req.GetConfig().GetTimestamp(),
		},
	}, result)
	if err != nil {
		return nil, err
	}
	loc.Name = resLoc.Name
	loc.Ips = append(loc.Ips, resLoc.Ips...)
	return loc, nil
}

func (s *server) handleInsertResult(
	// skipcq: GO-R1005
	ctx context.Context,
	alreadyExistsTgts []string,
	req *payload.Update_Request,
	result map[string]*errorState,
) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, "handleInsertResult"), apiName+"/handleInsertResult")
	defer trace.End(span)

	var mu sync.Mutex
	loc = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}

	err = s.gateway.DoMulti(ctx, alreadyExistsTgts, func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "DoMulti/"+target), apiName+"/"+vald.UpdateRPCName+"/"+target)
		defer trace.End(span)

		code := codes.OK
		ce, err := s.doUpdate(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Update(ctx, req, copts...)
		})
		if err != nil {
			st := statusOf(err, vald.UpdateRPCName)
			log.Warn(err)
			errhandler.RecordSpanStatus(span, st, err)
			code = st.Code()
		}
		mu.Lock()
		if err == nil && ce != nil {
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
		}
		result[target] = &errorState{err, code}
		mu.Unlock()
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId: req.GetVector().GetId(),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.InsertRPCName + ".DoMulti")
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpdateRPCName+" for "+vald.InsertRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			return errhandler.HandleError[payload.Object_Location](span, codes.Internal, err)
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpdateRPCName+" for "+vald.InsertRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}

	alreadyExistsTgts = alreadyExistsTgts[0:0]
	successTgts := make([]string, 0, len(result)/targetCapHintDivisor)
	for target, es := range result {
		switch {
		case es.err == nil:
			successTgts = append(successTgts, target)
		case es.code == codes.AlreadyExists:
			alreadyExistsTgts = append(alreadyExistsTgts, target)
			err = errors.Join(err, es.err)
		default:
			err = errors.Join(es.err, err)
		}
	}
	if err == nil || (len(successTgts) > 0 && len(result) == len(successTgts)+len(alreadyExistsTgts)) {
		log.Debugf(vald.UpdateRPCName+" for "+vald.InsertRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}

	reqInfo := &errdetails.RequestInfo{
		RequestId: req.GetVector().GetId(),
	}
	resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + "." + vald.UpdateRPCName)

	switch {
	case len(result) == len(alreadyExistsTgts):
		err = status.WrapWithAlreadyExists(vald.UpdateRPCName+" for "+vald.InsertRPCName+" API target same vector already exists", err, reqInfo, resInfo)
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Location](span, codes.AlreadyExists, err)
	case len(result) > len(successTgts)+len(alreadyExistsTgts): // Contains errors except for ALREADY_EXIST.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpdateRPCName+" for "+vald.InsertRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}
	log.Debugf(vald.UpdateRPCName+"for "+vald.InsertRPCName+" API request succeeded to %#v, err: %v", loc, err)
	return loc, nil
}

func (s *server) doInsert(
	ctx context.Context,
	req *payload.Insert_Request,
	f func(ctx context.Context) (*payload.Object_Location, error),
) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doInsert"), apiName+"/doInsert")
	defer trace.End(span)

	loc, err = f(ctx)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName)
		return nil, wrapTargetError(span, vald.InsertRPCName, err, reqInfo, resInfo)
	}
	return loc, nil
}

// StreamInsert handles bidirectional streaming for inserting objects.
// It wraps the bidirectional stream logic for the Insert RPC method.
// For each incoming request in the bidirectional stream, it calls the Insert function.
// The response is then sent back through the stream with the corresponding status or location information.
func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.StreamInsertRPCName), apiName+"/"+vald.StreamInsertRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Insert_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamInsertRPCName+"/id-"+req.GetVector().GetId())
			defer trace.End(sspan)
			res, err := s.Insert(ctx, req)
			if err != nil {
				st := statusOf(err, vald.InsertRPCName)
				errhandler.RecordSpanStatus(sspan, st, err)
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
		st := statusOf(err, vald.StreamInsertRPCName)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

// MultiInsert handles the insertion of multiple objects with the given requests.
// For each request in parallel, it calls the Insert function to insert an object.
// If an error occurs during any of the insertions, it accumulates the errors and returns them along with the successfully inserted locations.
func (s *server) MultiInsert(
	ctx context.Context, reqs *payload.Insert_MultiRequest,
) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.MultiInsertRPCName), apiName+"/"+vald.MultiInsertRPCName)
	defer trace.End(span)

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var emu sync.Mutex
	var errList []error
	eg := errgroup.NewBoundedFanOut(ctx, s.streamConcurrency)

	for idx, req := range reqs.GetRequests() {
		eg.Go(safety.RecoverFunc(func() error {
			ti := "errgroup.Go/id-" + req.GetVector().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiInsertRPCName+"/"+ti)
			defer trace.End(span)

			loc, err := s.Insert(ctx, req)
			if err != nil {
				st := statusOf(err, vald.InsertRPCName)
				errhandler.RecordSpanStatus(span, st, err)
				emu.Lock()
				errList = append(errList, err)
				emu.Unlock()
				return nil
			}
			res.Locations[idx] = loc
			return nil
		}))
	}
	if err := eg.Wait(); err != nil {
		errList = append(errList, err)
	}
	errs = errors.Join(errList...)
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiInsertRPCName+" gRPC error response",
			s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix+"/vald.v1."+vald.MultiInsertRPCName))
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return res, err
	}
	return res, nil
}

// Update handles the update of an object with the given request.
// If the request is proxied from another Mirror gateway, it sends the request only to the Vald gateway (LB gateway) of its own cluster.
// If the request is from a user, it sends requests to other Mirror gateways and the Vald gateway (LB gateway) of its own cluster.
// The result is a location of the updated object.
func (s *server) Update(
	ctx context.Context, req *payload.Update_Request,
) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateRPCName), apiName+"/"+vald.UpdateRPCName)
	defer trace.End(span)

	// If this condition is matched, it means that the request was proxied from another Mirror gateway.
	// So this component sends requests only to the Vald gateway (LB gateway) of its own cluster.
	if s.isProxied(ctx) {
		loc, err = s.doUpdate(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			_, derr := s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
				loc, err = vc.Update(ctx, req, copts...)
				return loc, err
			})
			return loc, errors.Join(derr, err)
		})
		if err != nil {
			st := statusOf(err, vald.UpdateRPCName)
			errhandler.RecordSpanStatus(span, st, err)
			log.Warn(err)
			return nil, err
		}
		log.Debugf("Update API succeeded to %#v", loc)
		return loc, nil
	}

	// If this condition is matched, it means the request from user.
	// So this component sends requests to other Mirror gateways and the Vald gateway (LB gateway) of its own cluster.
	return s.handleUpdate(ctx, req)
}

func (s *server) handleUpdate(
	ctx context.Context, req *payload.Update_Request,
) (loc *payload.Object_Location, err error) { // skipcq: GO-R1005
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, "handleUpdate"), apiName+"/handleUpdate")
	defer trace.End(span)

	var mu sync.Mutex
	result := make(map[string]*errorState) // map[target host: error state]
	loc = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpdateRPCName+"/"+target)
		defer trace.End(span)

		code := codes.OK
		ce, err := s.doUpdate(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Update(ctx, req, copts...)
		})
		if err != nil {
			st := statusOf(err, vald.UpdateRPCName)
			log.Warn(err)
			errhandler.RecordSpanStatus(span, st, err)
			code = st.Code()
		}
		mu.Lock()
		if err == nil && ce != nil {
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
		}
		result[target] = &errorState{err, code}
		mu.Unlock()
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + ".BroadCast")
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpdateRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			return errhandler.HandleError[payload.Object_Location](span, codes.Internal, err)
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}
	if len(result) == 0 {
		err = s.broadCastNoTargetError(vald.UpdateRPCName, &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		})
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Location](span, codes.Unavailable, err)
	}

	var alreadyExistsCnt int
	notFoundTgts := make([]string, 0, len(result)/targetCapHintDivisor)
	successTgts := make([]string, 0, len(result)/targetCapHintDivisor)
	for target, es := range result {
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
	}
	if err == nil || (len(successTgts) > 0 && len(result) == len(successTgts)+alreadyExistsCnt) {
		log.Debugf(vald.UpdateRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}

	reqInfo := &errdetails.RequestInfo{
		RequestId: req.GetVector().GetId(),
	}
	resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName)

	switch {
	case len(result) == len(notFoundTgts):
		err = status.WrapWithNotFound(vald.UpdateRPCName+" API target not found", err, reqInfo, resInfo)
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Location](span, codes.NotFound, err)
	case len(result) == alreadyExistsCnt:
		err = status.WrapWithAlreadyExists(vald.UpdateRPCName+" API target same vector already exists", err, reqInfo, resInfo)
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Location](span, codes.AlreadyExists, err)
	case len(result) > len(successTgts)+len(notFoundTgts)+alreadyExistsCnt: // Contains errors except for NOT_FOUND and ALREADY_EXIST.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}

	// In this case, the status code in the result object contains only OK or ALREADY_EXIST or NOT_FOUND.
	// And send Insert API request to NOT_FOUND cluster using query requested by the user.
	log.Warnf("failed to "+vald.UpdateRPCName+" API: %#v", err)

	resLoc, err := s.handleUpdateResult(ctx, notFoundTgts, &payload.Insert_Request{
		Vector: req.GetVector(),
		Config: &payload.Insert_Config{
			Timestamp: req.GetConfig().GetTimestamp(),
		},
	}, result)
	if err != nil {
		return nil, err
	}
	loc.Name = resLoc.Name
	loc.Ips = append(loc.Ips, resLoc.Ips...)
	return loc, nil
}

func (s *server) handleUpdateResult(
	// skipcq: GO-R1005
	ctx context.Context,
	notFoundTgts []string,
	req *payload.Insert_Request,
	result map[string]*errorState,
) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, "handleUpdateResult"), apiName+"/handleUpdateResult")
	defer trace.End(span)

	var mu sync.Mutex
	loc = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}

	err = s.gateway.DoMulti(ctx, notFoundTgts, func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
		defer trace.End(span)

		code := codes.OK
		ce, err := s.doInsert(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Insert(ctx, req, copts...)
		})
		if err != nil {
			st := statusOf(err, vald.InsertRPCName)
			log.Warn(err)
			errhandler.RecordSpanStatus(span, st, err)
			code = st.Code()
		}
		mu.Lock()
		if err == nil && ce != nil {
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
		}
		result[target] = &errorState{err, code}
		mu.Unlock()
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.InsertRPCName + ".BroadCast")
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.InsertRPCName+" for "+vald.UpdateRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			return errhandler.HandleError[payload.Object_Location](span, codes.Internal, err)
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.InsertRPCName+" for "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}

	alreadyExistsCnt := 0
	notFoundTgts = notFoundTgts[0:0]
	successTgts := make([]string, 0, len(result)/targetCapHintDivisor)
	for target, em := range result {
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
	}
	if err == nil || (len(successTgts) > 0 && len(result) == len(successTgts)+alreadyExistsCnt) {
		log.Debugf(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}

	reqInfo := &errdetails.RequestInfo{
		RequestId:   req.GetVector().GetId(),
		ServingData: errdetails.Serialize(req),
	}
	resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.InsertRPCName)

	switch {
	case len(result) == len(notFoundTgts):
		err = status.WrapWithNotFound(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API target not found", err, reqInfo, resInfo)
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Location](span, codes.NotFound, err)
	case len(result) == alreadyExistsCnt:
		err = status.WrapWithAlreadyExists(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API target same vector already exists", err, reqInfo, resInfo)
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Location](span, codes.AlreadyExists, err)
	case len(result) > len(successTgts)+len(notFoundTgts)+alreadyExistsCnt: // Contains errors except for NOT_FOUND and ALREADY_EXIST.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.InsertRPCName+" for "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}
	if len(notFoundTgts) > 0 { // The Insert was never applied to some NOT_FOUND targets (e.g. skipped by DoMulti).
		err = status.WrapWithNotFound(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API target not found", err, reqInfo, resInfo)
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Location](span, codes.NotFound, err)
	}
	log.Debugf(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API request succeeded to %#v, err: %v", loc, err)
	return loc, nil
}

func (s *server) doUpdate(
	ctx context.Context,
	req *payload.Update_Request,
	f func(ctx context.Context) (*payload.Object_Location, error),
) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doUpdate"), apiName+"/doUpdate")
	defer trace.End(span)

	loc, err = f(ctx)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName)
		return nil, wrapTargetError(span, vald.UpdateRPCName, err, reqInfo, resInfo)
	}
	return loc, nil
}

// StreamUpdate handles bidirectional streaming for updating objects.
// It wraps the bidirectional stream logic for the Update RPC method.
// For each incoming request in the bidirectional stream, it calls the Update function.
// The response is then sent back through the stream with the corresponding status or location information.
func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.StreamUpdateRPCName), apiName+"/"+vald.StreamUpdateRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Update_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamUpdateRPCName+"/id-"+req.GetVector().GetId())
			defer trace.End(sspan)
			res, err := s.Update(ctx, req)
			if err != nil {
				st := statusOf(err, vald.UpdateRPCName)
				errhandler.RecordSpanStatus(sspan, st, err)
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
		st := statusOf(err, vald.StreamUpdateRPCName)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

// MultiUpdate handles the update of multiple objects with the given requests.
// For each request in parallel, it calls the Update function to update an object.
// If an error occurs during any of the insertions, it accumulates the errors and returns them along with the successfully updated locations.
func (s *server) MultiUpdate(
	ctx context.Context, reqs *payload.Update_MultiRequest,
) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.MultiUpdateRPCName), apiName+"/"+vald.MultiUpdateRPCName)
	defer trace.End(span)

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var emu sync.Mutex
	var errList []error
	eg := errgroup.NewBoundedFanOut(ctx, s.streamConcurrency)

	for idx, req := range reqs.GetRequests() {
		eg.Go(safety.RecoverFunc(func() error {
			ti := "errgroup.Go/id-" + req.GetVector().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiUpdateRPCName+"/"+ti)
			defer trace.End(span)

			loc, err := s.Update(ctx, req)
			if err != nil {
				st := statusOf(err, vald.UpdateRPCName)
				errhandler.RecordSpanStatus(span, st, err)
				emu.Lock()
				errList = append(errList, err)
				emu.Unlock()
				return nil
			}
			res.Locations[idx] = loc
			return nil
		}))
	}
	if err := eg.Wait(); err != nil {
		errList = append(errList, err)
	}
	errs = errors.Join(errList...)
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiUpdateRPCName+" gRPC error response",
			s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix+"/vald.v1."+vald.MultiUpdateRPCName))
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return res, err
	}
	return res, nil
}

// Upsert handles the upsert of an object with the given request.
// If the request is proxied from another Mirror gateway, the request is forwarded to the Vald gateway (LB gateway) of its own cluster.
// If the request is from a user, it is sent to other Mirror gateways and the Vald gateway (LB gateway) of its own cluster.
// The result is a location of the upserted object.
func (s *server) Upsert(
	ctx context.Context, req *payload.Upsert_Request,
) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.UpsertRPCName), apiName+"/"+vald.UpsertRPCName)
	defer trace.End(span)

	// If this condition is matched, it means that the request was proxied from another Mirror gateway.
	// So this component sends requests only to the Vald gateway (LB gateway) of its own cluster.
	if s.isProxied(ctx) {
		loc, err = s.doUpsert(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			_, derr := s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
				loc, err = vc.Upsert(ctx, req, copts...)
				return loc, err
			})
			return loc, errors.Join(derr, err)
		})
		if err != nil {
			st := statusOf(err, vald.UpsertRPCName)
			errhandler.RecordSpanStatus(span, st, err)
			log.Warn(err)
			return nil, err
		}
		log.Debugf("Upsert API succeeded to %#v", loc)
		return loc, nil
	}

	// If this condition is matched, it means the request from user.
	// So this component sends requests to other Mirror gateways and the Vald gateway (LB gateway) of its own cluster.
	return s.handleUpsert(ctx, req)
}

func (s *server) handleUpsert(
	ctx context.Context, req *payload.Upsert_Request,
) (loc *payload.Object_Location, err error) { // skipcq: GO-R1005
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, "handleUpsert"), apiName+"/handleUpsert")
	defer trace.End(span)

	var mu sync.Mutex
	result := make(map[string]*errorState) // map[target host: error state]
	loc = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpsertRPCName+"/"+target)
		defer trace.End(span)

		code := codes.OK
		ce, err := s.doUpsert(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Upsert(ctx, req, copts...)
		})
		if err != nil {
			st := statusOf(err, vald.UpsertRPCName)
			log.Warn(err)
			errhandler.RecordSpanStatus(span, st, err)
			code = st.Code()
		}
		mu.Lock()
		if err == nil && ce != nil {
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
		}
		result[target] = &errorState{err, code}
		mu.Unlock()
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + ".BroadCast")
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpsertRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			return errhandler.HandleError[payload.Object_Location](span, codes.Internal, err)
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpsertRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}
	if len(result) == 0 {
		err = s.broadCastNoTargetError(vald.UpsertRPCName, &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		})
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Location](span, codes.Unavailable, err)
	}

	var alreadyExistsCnt int
	successTgts := make([]string, 0, len(result)/targetCapHintDivisor)
	for target, es := range result {
		switch {
		case es.err == nil:
			successTgts = append(successTgts, target)
		case es.code == codes.AlreadyExists:
			alreadyExistsCnt++
			err = errors.Join(err, es.err)
		default:
			err = errors.Join(es.err, err)
		}
	}
	if err == nil || (len(successTgts) > 0 && len(result) == len(successTgts)+alreadyExistsCnt) {
		log.Debugf(vald.UpsertRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}
	reqInfo := &errdetails.RequestInfo{
		RequestId:   req.GetVector().GetId(),
		ServingData: errdetails.Serialize(req),
	}
	resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName)

	switch {
	case len(result) == alreadyExistsCnt:
		err = status.WrapWithAlreadyExists(vald.UpsertRPCName+" API target same vector already exists", err, reqInfo, resInfo)
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Location](span, codes.AlreadyExists, err)
	default:
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.UpsertRPCName+" gRPC error response", reqInfo, resInfo)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}
}

func (s *server) doUpsert(
	ctx context.Context,
	req *payload.Upsert_Request,
	f func(ctx context.Context) (*payload.Object_Location, error),
) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doUpsert"), apiName+"/doUpsert")
	defer trace.End(span)

	loc, err = f(ctx)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName)
		return nil, wrapTargetError(span, vald.UpsertRPCName, err, reqInfo, resInfo)
	}
	return loc, nil
}

// StreamUpsert handles bidirectional streaming for upsert objects.
// It wraps the bidirectional stream logic for the Upsert RPC method.
// For each incoming request in the bidirectional stream, it calls the Upsert function.
// The response is then sent back through the stream with the corresponding status or location information.
func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.StreamUpsertRPCName), apiName+"/"+vald.StreamUpsertRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Upsert_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamUpsertRPCName+"/id-"+req.GetVector().GetId())
			defer trace.End(sspan)
			res, err := s.Upsert(ctx, req)
			if err != nil {
				st := statusOf(err, vald.UpsertRPCName)
				errhandler.RecordSpanStatus(sspan, st, err)
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
		st := statusOf(err, vald.StreamUpsertRPCName)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

// MultiUpsert handles the upsert of multiple objects with the given requests.
// For each request in parallel, it calls the Upsert function to upsert an object.
// If an error occurs during any of the insertions, it accumulates the errors and returns them along with the successfully upserted locations.
func (s *server) MultiUpsert(
	ctx context.Context, reqs *payload.Upsert_MultiRequest,
) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.MultiUpsertRPCName), apiName+"/"+vald.MultiUpsertRPCName)
	defer trace.End(span)

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var emu sync.Mutex
	var errList []error
	eg := errgroup.NewBoundedFanOut(ctx, s.streamConcurrency)

	for idx, req := range reqs.GetRequests() {
		eg.Go(safety.RecoverFunc(func() error {
			ti := "errgroup.Go/id-" + req.GetVector().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiUpsertRPCName+"/"+ti)
			defer trace.End(span)

			loc, err := s.Upsert(ctx, req)
			if err != nil {
				st := statusOf(err, vald.UpsertRPCName)
				errhandler.RecordSpanStatus(span, st, err)
				emu.Lock()
				errList = append(errList, err)
				emu.Unlock()
				return nil
			}
			res.Locations[idx] = loc
			return nil
		}))
	}
	if err := eg.Wait(); err != nil {
		errList = append(errList, err)
	}
	errs = errors.Join(errList...)
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiUpsertRPCName+" gRPC error response",
			s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix+"/vald.v1."+vald.MultiUpsertRPCName))
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return res, err
	}
	return res, nil
}

// Remove handles the remove of an object with the given request.
// If the request is proxied from another Mirror gateway, the request is forwarded to the Vald gateway (LB gateway) of its own cluster.
// If the request is from a user, it is sent to other Mirror gateways and the Vald gateway (LB gateway) of its own cluster.
// The result is a location of the removed object.
func (s *server) Remove(
	ctx context.Context, req *payload.Remove_Request,
) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer trace.End(span)

	// If this condition is matched, it means that the request was proxied from another Mirror gateway.
	// So this component sends requests only to the Vald gateway (LB gateway) of its own cluster.
	if s.isProxied(ctx) {
		loc, err = s.doRemove(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			_, derr := s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
				loc, err = vc.Remove(ctx, req, copts...)
				return loc, err
			})
			return loc, errors.Join(derr, err)
		})
		if err != nil {
			st := statusOf(err, vald.RemoveRPCName)
			errhandler.RecordSpanStatus(span, st, err)
			log.Warn(err)
			return nil, err
		}
		log.Debugf("Remove API remove succeeded to %#v", loc)
		return loc, nil
	}

	// If this condition is matched, it means the request from user.
	// So this component sends requests to other Mirror gateways and the Vald gateway (LB gateway) of its own cluster.
	return s.handleRemove(ctx, req)
}

func (s *server) handleRemove(
	ctx context.Context, req *payload.Remove_Request,
) (loc *payload.Object_Location, err error) { // skipcq: GO-R1005
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, "handleRemove"), apiName+"/handleRemove")
	defer trace.End(span)

	var mu sync.Mutex
	result := make(map[string]*errorState) // map[target host: error state]
	loc = &payload.Object_Location{
		Uuid: req.GetId().GetId(),
		Ips:  make([]string, 0),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.RemoveRPCName+"/"+target)
		defer trace.End(span)

		code := codes.OK
		ce, err := s.doRemove(ctx, req, func(ctx context.Context) (*payload.Object_Location, error) {
			return vc.Remove(ctx, req, copts...)
		})
		if err != nil {
			st := statusOf(err, vald.RemoveRPCName)
			log.Warn(err)
			errhandler.RecordSpanStatus(span, st, err)
			code = st.Code()
		}
		mu.Lock()
		if err == nil && ce != nil {
			loc.Name = ce.GetName()
			loc.Ips = append(loc.Ips, ce.GetIps()...)
		}
		result[target] = &errorState{err, code}
		mu.Unlock()
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId: req.GetId().GetId(),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + ".BroadCast")
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.RemoveRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			return errhandler.HandleError[payload.Object_Location](span, codes.Internal, err)
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}
	if len(result) == 0 {
		err = s.broadCastNoTargetError(vald.RemoveRPCName, &errdetails.RequestInfo{
			RequestId: req.GetId().GetId(),
		})
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Location](span, codes.Unavailable, err)
	}

	var notFoundCnt int
	successTgts := make([]string, 0, len(result)/targetCapHintDivisor)
	for target, es := range result {
		switch {
		case es.err == nil:
			successTgts = append(successTgts, target)
		case es.code == codes.NotFound:
			notFoundCnt++
			err = errors.Join(err, es.err)
		default:
			err = errors.Join(es.err, err)
		}
	}
	if err == nil || (len(successTgts) > 0 && len(result) == len(successTgts)+notFoundCnt) {
		log.Debugf(vald.RemoveRPCName+" API request succeeded to %#v", loc)
		return loc, nil
	}

	reqInfo := &errdetails.RequestInfo{
		RequestId: req.GetId().GetId(),
	}
	resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName)

	switch {
	case len(result) == notFoundCnt:
		err = status.WrapWithNotFound(vald.RemoveRPCName+" API id "+req.GetId().GetId()+" not found", err, reqInfo, resInfo)
		return errhandler.HandleError[payload.Object_Location](span, codes.NotFound, err)
	default:
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveRPCName+" gRPC error response", reqInfo, resInfo)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}
}

func (s *server) doRemove(
	ctx context.Context,
	req *payload.Remove_Request,
	f func(ctx context.Context) (*payload.Object_Location, error),
) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doRemove"), apiName+"/doRemove")
	defer trace.End(span)

	loc, err = f(ctx)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId: req.GetId().GetId(),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName)
		return nil, wrapTargetError(span, vald.RemoveRPCName, err, reqInfo, resInfo)
	}
	return loc, nil
}

// StreamRemove handles bidirectional streaming for removing objects.
// It wraps the bidirectional stream logic for the Remove RPC method.
// For each incoming request in the bidirectional stream, it calls the Remove function.
// The response is then sent back through the stream with the corresponding status or location information.
func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.StreamRemoveRPCName), apiName+"/"+vald.StreamRemoveRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Remove_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamRemoveRPCName+"/id-"+req.GetId().GetId())
			defer trace.End(sspan)
			res, err := s.Remove(ctx, req)
			if err != nil {
				st := statusOf(err, vald.RemoveRPCName)
				errhandler.RecordSpanStatus(sspan, st, err)
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
		st := statusOf(err, vald.StreamRemoveRPCName)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

// MultiRemove handles the remove of multiple objects with the given requests.
// For each request in parallel, it calls the Remove function to insert an object.
// If an error occurs during any of the insertions, it accumulates the errors and returns them along with the successfully removed locations.
func (s *server) MultiRemove(
	ctx context.Context, reqs *payload.Remove_MultiRequest,
) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
	defer trace.End(span)

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var emu sync.Mutex
	var errList []error
	eg := errgroup.NewBoundedFanOut(ctx, s.streamConcurrency)

	for idx, req := range reqs.GetRequests() {
		eg.Go(safety.RecoverFunc(func() error {
			ti := "errgroup.Go/id-" + req.GetId().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiRemoveRPCName+"/"+ti)
			defer trace.End(span)

			loc, err := s.Remove(ctx, req)
			if err != nil {
				st := statusOf(err, vald.RemoveRPCName)
				errhandler.RecordSpanStatus(span, st, err)
				emu.Lock()
				errList = append(errList, err)
				emu.Unlock()
				return nil
			}
			res.Locations[idx] = loc
			return nil
		}))
	}
	if err := eg.Wait(); err != nil {
		errList = append(errList, err)
	}
	errs = errors.Join(errList...)
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiRemoveRPCName+" gRPC error response",
			s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix+"/vald.v1."+vald.MultiRemoveRPCName))
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return res, err
	}
	return res, nil
}

// RemoveByTimestamp handles the remove of an object with the given request.
// If the request is proxied from another Mirror gateway, the request is forwarded to the Vald gateway (LB gateway) of its own cluster.
// If the request is from a user, it is sent to other Mirror gateways and the Vald gateway (LB gateway) of its own cluster.
// The result is a location of the removed object.
func (s *server) RemoveByTimestamp(
	ctx context.Context, req *payload.Remove_TimestampRequest,
) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveByTimestampRPCName), apiName+"/"+vald.RemoveByTimestampRPCName)
	defer trace.End(span)

	// If this condition is matched, it means that the request was proxied from another Mirror gateway.
	// So this component sends requests only to the Vald gateway (LB gateway) of its own cluster.
	if s.isProxied(ctx) {
		locs, err = s.doRemoveByTimestamp(ctx, req, func(ctx context.Context) (*payload.Object_Locations, error) {
			_, derr := s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
				locs, err = vc.RemoveByTimestamp(ctx, req, copts...)
				return locs, err
			})
			return locs, errors.Join(derr, err)
		})
		if err != nil {
			st := statusOf(err, vald.RemoveByTimestampRPCName)
			log.Warn(err)
			errhandler.RecordSpanStatus(span, st, err)
			return nil, err
		}
		log.Debugf("RemoveByTimestamp API remove succeeded to %#v", locs)
		return locs, nil
	}

	// If this condition is matched, it means the request from user.
	// So this component sends requests to other Mirror gateways and the Vald gateway (LB gateway) of its own cluster.
	return s.handleRemoveByTimestamp(ctx, req)
}

func (s *server) handleRemoveByTimestamp(
	ctx context.Context, req *payload.Remove_TimestampRequest,
) (locs *payload.Object_Locations, err error) { // skipcq: GO-R1005
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, "handleRemoveByTimestamp"), apiName+"/handleRemoveByTimestamp")
	defer trace.End(span)

	var mu sync.Mutex
	result := make(map[string]*errorState) // map[target host: error state]
	locs = new(payload.Object_Locations)

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.RemoveByTimestampRPCName+"/"+target)
		defer trace.End(span)

		code := codes.OK
		res, err := s.doRemoveByTimestamp(ctx, req, func(ctx context.Context) (*payload.Object_Locations, error) {
			return vc.RemoveByTimestamp(ctx, req, copts...)
		})
		if err != nil {
			st := statusOf(err, vald.RemoveByTimestampRPCName)
			log.Warn(err)
			errhandler.RecordSpanStatus(span, st, err)
			code = st.Code()
		}
		mu.Lock()
		if err == nil && res != nil {
			locs.Locations = append(locs.Locations, res.GetLocations()...)
		}
		result[target] = &errorState{err, code}
		mu.Unlock()
		return err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveByTimestampRPCName + ".BroadCast")
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.RemoveByTimestampRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			log.Warn(err)
			return errhandler.HandleError[payload.Object_Locations](span, codes.Internal, err)
		}

		// There is no possibility to reach this part, but we add error handling just in case.
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveByTimestampRPCName+" gRPC error response", reqInfo, resInfo,
		)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}
	if len(result) == 0 {
		err = s.broadCastNoTargetError(vald.RemoveByTimestampRPCName, &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		})
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Locations](span, codes.Unavailable, err)
	}

	var notFoundCnt int
	successTgts := make([]string, 0, len(result)/targetCapHintDivisor)
	for target, es := range result {
		switch {
		case es.err == nil:
			successTgts = append(successTgts, target)
		case es.code == codes.NotFound:
			notFoundCnt++
			err = errors.Join(err, es.err)
		default:
			err = errors.Join(es.err, err)
		}
	}
	if err == nil || (len(successTgts) > 0 && len(result) == len(successTgts)+notFoundCnt) {
		log.Debugf(vald.RemoveByTimestampRPCName+" API request succeeded to %#v", locs)
		return locs, nil
	}

	reqInfo := &errdetails.RequestInfo{
		ServingData: errdetails.Serialize(req),
	}
	resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveByTimestampRPCName)

	switch {
	case len(result) == notFoundCnt:
		err = status.WrapWithNotFound(vald.RemoveByTimestampRPCName+" API target not found", err, reqInfo, resInfo)
		log.Warn(err)
		return errhandler.HandleError[payload.Object_Locations](span, codes.NotFound, err)
	default:
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveByTimestampRPCName+" gRPC error response", reqInfo, resInfo)
		log.Warn(err)
		errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
		return nil, err
	}
}

func (s *server) doRemoveByTimestamp(
	ctx context.Context,
	req *payload.Remove_TimestampRequest,
	f func(ctx context.Context) (*payload.Object_Locations, error),
) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doRemoveByTimestamp"), apiName+"/doRemoveByTimestamp")
	defer trace.End(span)

	locs, err = f(ctx)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveByTimestampRPCName)
		return nil, wrapTargetError(span, vald.RemoveByTimestampRPCName, err, reqInfo, resInfo)
	}
	return locs, nil
}

// GetObject bypasses the incoming GetObject request to Vald LB gateway in its own cluster.
func (s *server) GetObject(
	ctx context.Context, req *payload.Object_VectorRequest,
) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer trace.End(span)

	_, err = s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, _ string, vc service.MirrorClient, copts ...grpc.CallOption) (any, error) {
		vec, err = vc.GetObject(ctx, req, copts...)
		return vec, err
	})
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetId().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := s.vAddrResourceInfo(vald.GetObjectRPCName)
		return nil, wrapTargetError(span, vald.GetObjectRPCName, err, reqInfo, resInfo)
	}
	return vec, nil
}

// StreamGetObject bypasses it as a GetObject request to the Vald LB gateway in its own cluster.
func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamGetObjectRPCName), apiName+"/"+vald.StreamGetObjectRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Object_VectorRequest) (*payload.Object_StreamVector, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamInsertRPCName+"/id-"+req.GetId().GetId())
			defer trace.End(sspan)
			res, err := s.GetObject(ctx, req)
			if err != nil {
				st, _ := status.FromError(err)
				if st == nil || st.Message() == "" {
					// This condition is implemented just in case to prevent nil pointer errors when retrieving st.Code(), although it is unlikely to match this condition.
					log.Errorf("gRPC call returned not a gRPC status error: %v", err)
					st = status.New(codes.Internal, "failed to parse "+vald.GetObjectRPCName+" gRPC error response")
				}
				errhandler.RecordSpanStatus(sspan, st, err)
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
		st := statusOf(err, vald.StreamGetObjectRPCName)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

// StreamListObject bypasses it as a StreamListObject request to the Vald gateway (LB gateway) in its own cluster.
func (s *server) StreamListObject(
	req *payload.Object_List_Request, stream vald.Object_StreamListObjectServer,
) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamListObjectRPCName), apiName+"/"+vald.StreamListObjectRPCName)
	defer trace.End(span)

	// forwarded records whether any message has already reached the caller. It
	// spans every invocation of the callback below (the gRPC client backoff
	// retry and the gateway.Do fallback to other mirrors), each of which would
	// restart the listing from the beginning on the same server stream.
	var forwarded atomic.Bool
	_, err := s.gateway.Do(ctx, s.vAddr, func(ctx context.Context, target string, vc service.MirrorClient, copts ...grpc.CallOption) (obj any, err error) {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "Do/"+target), apiName+"/"+vald.StreamListObjectRPCName+"/"+target)
		defer trace.End(span)
		return obj, s.doStreamListObject(ctx, vc, req, stream, &forwarded, copts...)
	})
	if err != nil {
		// gateway.Do's RoundRobin fallback returns the error of the last conn it
		// visited (e.g. ErrTargetNotFound for the own LB address), not the
		// callback's Aborted, so re-assert Aborted once anything was delivered.
		if forwarded.Load() && !status.Is(err, codes.Aborted) {
			err = status.WrapWithAborted(vald.StreamListObjectRPCName+" API stream failed after partial delivery", err)
		}
		st := statusOf(err, vald.StreamListObjectRPCName)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

func (s *server) doStreamListObject(
	ctx context.Context,
	vc service.MirrorClient,
	req *payload.Object_List_Request,
	server vald.Object_StreamListObjectServer,
	forwarded *atomic.Bool,
	copts ...grpc.CallOption,
) (err error) { // skipcq: GO-R1005
	// A failure after partial delivery is returned as Aborted (not retried) and
	// a re-invocation must not open a new stream either, to avoid sending
	// duplicates to the caller.
	if forwarded.Load() {
		return status.WrapWithAborted(vald.StreamListObjectRPCName+" API stream failed after partial delivery", nil)
	}
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, egctx := errgroup.WithContext(cctx)
	eg.SetLimit(s.streamConcurrency)
	// Bind the client stream to egctx (created before it here), which is cancelled
	// both on EOF (cancel() -> cctx -> egctx) and on an eg.Go worker error
	// (errgroup), so client.Recv() below actually aborts on either. Binding it to
	// the caller's ancestor ctx (the previous ordering) made the cancellation a
	// no-op and silently defeated the early-abort on a downstream send failure.
	client, err := vc.StreamListObject(egctx, req, copts...)
	if err != nil {
		return err
	}

	var mu, rmu sync.Mutex
	var egCnt atomic.Int64
	for {
		select {
		case <-egctx.Done():
			// If the root context is not canceld error, it is treated as an error.
			if ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) {
				err = errors.Join(ctx.Err(), err)
			}
			if egerr := eg.Wait(); egerr != nil {
				err = errors.Join(err, egerr)
			}
			if err != nil && forwarded.Load() {
				err = status.WrapWithAborted(vald.StreamListObjectRPCName+" API stream failed after partial delivery", err)
			}
			return err
		default:
			eg.Go(safety.RecoverFunc(func() (err error) {
				id := fmt.Sprintf("stream-%020d", egCnt.Add(1))
				_, span := trace.StartSpan(egctx, apiName+"/streamListObject/"+id)
				defer trace.End(span)

				rmu.Lock()
				res, err := client.Recv()
				rmu.Unlock()
				if err != nil {
					if errors.Is(err, io.EOF) {
						cancel()
						return nil
					}
					err = errors.ErrServerStreamClientRecv(err)
					var attr trace.Attributes
					switch {
					case errors.Is(err, context.Canceled):
						err = status.WrapWithCanceled("Stream Recv returned canceld error at "+id, err)
						attr = trace.StatusCodeCancelled(err.Error())
					case errors.Is(err, context.DeadlineExceeded):
						err = status.WrapWithDeadlineExceeded("Stream Recv returned deadline exceeded error at "+id, err)
						attr = trace.StatusCodeDeadlineExceeded(err.Error())
					default:
						var (
							st  *status.Status
							msg string
						)
						st, msg, err = status.ParseError(err, codes.Internal, "Stream Recv returned an error at "+id)
						if st != nil {
							attr = trace.FromGRPCStatus(st.Code(), msg)
						}
					}
					log.Warn(err)
					errhandler.RecordSpanAttrs(span, attr, err)
					return err
				}
				if res.GetVector() == nil {
					return nil
				}

				mu.Lock()
				err = server.Send(res)
				mu.Unlock()
				if err != nil {
					if errors.Is(err, io.EOF) {
						cancel()
						return nil
					}
					err = errors.ErrServerStreamServerSend(err)
					var attr trace.Attributes
					switch {
					case errors.Is(err, context.Canceled):
						err = status.WrapWithCanceled("Stream Send returned canceld error at "+id, err)
						attr = trace.StatusCodeCancelled(err.Error())
					case errors.Is(err, context.DeadlineExceeded):
						err = status.WrapWithDeadlineExceeded("Stream Send returned deadline exceeded error at "+id, err)
						attr = trace.StatusCodeDeadlineExceeded(err.Error())
					default:
						var (
							st  *status.Status
							msg string
						)
						st, msg, err = status.ParseError(err, codes.Internal, "Stream Send returned an error at "+id)
						if st != nil {
							attr = trace.FromGRPCStatus(st.Code(), msg)
						}
					}
					log.Warn(err)
					errhandler.RecordSpanAttrs(span, attr, err)
					return err
				}
				forwarded.Store(true)
				return nil
			}))
		}
	}
}

// TODO: implement Flush handler.
func (s *server) Flush(
	ctx context.Context, req *payload.Flush_Request,
) (*payload.Info_Index_Count, error) {
	return s.UnimplementedFlushServer.Flush(ctx, req)
}

func (s *server) isProxied(ctx context.Context) bool {
	return s.gateway.FromForwardedContext(ctx) != ""
}

type errorState struct {
	err  error
	code codes.Code
}
