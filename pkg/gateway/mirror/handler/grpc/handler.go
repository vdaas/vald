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
	"sync"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	vclient "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

type MirrorServer interface {
	vald.Server
	mirror.MirrorServer
}

type server struct {
	eg                errgroup.Group
	gateway           service.Gateway
	client            vclient.Client // LB Gateway client in the same cluster.
	timeout           time.Duration
	replica           int
	streamConcurrency int
	name              string
	ip                string
	vald.UnimplementedValdServer
	mirror.MirrorServer
}

const (
	apiName      = "vald/gateway/mirror"
	rollbackName = "Rollback"
)

func New(opts ...Option) MirrorServer {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Register(context.Context, *payload.Mirror_Targets) (*payload.Mirror_Targets, error) {
	return nil, nil
}

func (s *server) Advertise(context.Context, *payload.Mirror_Targets) (*payload.Mirror_Targets, error) {
	return nil, nil
}

func (s *server) Exists(ctx context.Context, meta *payload.Object_ID) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.ExistsRPCName), apiName+"/"+vald.ExistsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	id, err = s.client.Exists(ctx, meta, s.client.GRPCClient().GetCallOption()...)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.ExistsRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
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
	res, err = s.client.Search(ctx, req)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error,
) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.SearchByIDRPCName), apiName+"/"+vald.SearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.client.SearchByID(ctx, req)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchByIDRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Search(ctx, data.(*payload.Search_Request))
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
		})
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
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_IDRequest)
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamSearchByIDRPCName+"/id-"+req.GetId())
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
		})
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

func (s *server) MultiSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err := s.client.MultiSearch(ctx, reqs)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.MultiSearchRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) MultiSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.client.MultiSearchByID(ctx, reqs)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.MultiSearchByIDRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
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
	res, err = s.client.LinearSearch(ctx, req)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
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
	res, err = s.client.LinearSearchByID(ctx, req)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchByIDRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) StreamLinearSearch(stream vald.Search_StreamLinearSearchServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamLinearSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearch(ctx, data.(*payload.Search_Request))
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
		})

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
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_IDRequest)
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamLinearSearchByIDRPCName+"/id-"+req.GetId())
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
		})
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

func (s *server) MultiLinearSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.client.MultiLinearSearch(ctx, reqs)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.MultiSearchRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) MultiLinearSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.client.MultiLinearSearchByID(ctx, reqs)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.MultiLinearSearchRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
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
	handleSpan := func(rpcName string, span trace.Span, err error) error {
		if err == nil {
			return nil
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.InsertRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).Insert(sctx, req, copts...)
		return handleSpan(vald.InsertRPCName, sspan, err)
	})
	if err != nil {
		if err := handleSpan(rollbackName+" for "+vald.InsertRPCName, span, s.insertRollback(ctx, req)); err != nil {
			return nil, err
		}
		return nil, handleSpan(vald.InsertRPCName, span, err)
	}

	ce, err = s.client.Insert(ctx, req)
	if err := handleSpan(vald.InsertRPCName, span, err); err != nil {
		return nil, err
	}
	return ce, nil
}

// insertRollback executes a Remove RPC for rollback of Insert RPC.
func (s *server) insertRollback(ctx context.Context, req *payload.Insert_Request) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	handleSpan := func(rpcName string, span trace.Span, err error) error {
		if err == nil {
			return nil
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	removeReq := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: true,
		},
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.RemoveRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).Remove(sctx, removeReq)
		return handleSpan(vald.RemoveRPCName, sspan, err)
	})
	if err := handleSpan(vald.RemoveRPCName, span, err); err != nil {
		return err
	}
	return nil
}

// func (s *server) _Insert(ctx context.Context, req *payload.Insert_Request) (ce *payload.Object_Location, err error) {
// 	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.InsertRPCName), apiName+"/"+vald.InsertRPCName)
// 	defer func() {
// 		if span != nil {
// 			span.End()
// 		}
// 	}()
// 	uuid := req.GetVector().GetId()
// 	if len(uuid) == 0 {
// 		err = errors.ErrInvalidMetaDataConfig
// 		err = status.WrapWithInvalidArgument(vald.InsertRPCName+" API invalid uuid", err,
// 			&errdetails.RequestInfo{
// 				ServingData: errdetails.Serialize(req),
// 			},
// 			&errdetails.BadRequest{
// 				FieldViolations: []*errdetails.BadRequestFieldViolation{
// 					{
// 						Field:       "invalid id",
// 						Description: err.Error(),
// 					},
// 				},
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
// 			})
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	vec := req.GetVector().GetVector()
// 	vl := len(vec)
// 	if vl < algorithm.MinimumVectorDimensionSize {
// 		err = errors.ErrInvalidDimensionSize(vl, 0)
// 		err = status.WrapWithInvalidArgument(vald.InsertRPCName+" API invalid vector argument", err,
// 			&errdetails.RequestInfo{
// 				RequestId:   uuid,
// 				ServingData: errdetails.Serialize(req),
// 			},
// 			&errdetails.BadRequest{
// 				FieldViolations: []*errdetails.BadRequestFieldViolation{
// 					{
// 						Field:       "vector dimension size",
// 						Description: err.Error(),
// 					},
// 				},
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	if !req.GetConfig().GetSkipStrictExistCheck() {
// 		id, err := s.Exists(ctx, &payload.Object_ID{
// 			Id: uuid,
// 		})
// 		if err == nil && id != nil && len(id.GetId()) != 0 {
// 			if err == nil {
// 				err = errors.ErrMetaDataAlreadyExists(uuid)
// 			}
// 			st, msg, err := status.ParseError(err, codes.AlreadyExists,
// 				"error "+vald.InsertRPCName+" API ID = "+uuid+" already exists",
// 				&errdetails.RequestInfo{
// 					RequestId:   uuid,
// 					ServingData: errdetails.Serialize(req),
// 				},
// 				&errdetails.ResourceInfo{
// 					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + "." + vald.ExistsRPCName,
// 					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 				}, info.Get())
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
// 			return nil, err
// 		}
// 		if req.GetConfig() != nil {
// 			req.GetConfig().SkipStrictExistCheck = true
// 		} else {
// 			req.Config = &payload.Insert_Config{SkipStrictExistCheck: true}
// 		}
// 	}
//
// 	mu := new(sync.Mutex)
// 	ce = &payload.Object_Location{
// 		Uuid: uuid,
// 		Ips:  make([]string, 0, s.replica),
// 	}
// 	if req.GetConfig().GetTimestamp() == 0 {
// 		now := time.Now().UnixNano()
// 		if req.GetConfig() == nil {
// 			req.Config = &payload.Insert_Config{
// 				Timestamp: now,
// 			}
// 		} else {
// 			req.GetConfig().Timestamp = now
// 		}
// 	}
// 	emu := new(sync.Mutex)
// 	var errs error
// 	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
// 		ctx, span := trace.StartSpan(ctx, apiName+"."+vald.InsertRPCName+"/"+target)
// 		defer func() {
// 			if span != nil {
// 				span.End()
// 			}
// 		}()
// 		loc, err := vc.Insert(ctx, req, copts...)
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, context.Canceled),
// 				errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
// 				if span != nil {
// 					span.RecordError(err)
// 					span.SetAttributes(trace.StatusCodeCancelled(
// 						errdetails.ValdGRPCResourceTypePrefix +
// 							"/vald.v1." + vald.InsertRPCName + ".DoMulti/" +
// 							target + " canceled: " + err.Error())...)
// 					span.SetStatus(trace.StatusError, err.Error())
// 				}
// 				return nil
// 			case errors.Is(err, context.DeadlineExceeded),
// 				errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
// 				if span != nil {
// 					span.RecordError(err)
// 					span.SetAttributes(trace.StatusCodeDeadlineExceeded(
// 						errdetails.ValdGRPCResourceTypePrefix +
// 							"/vald.v1." + vald.InsertRPCName + ".DoMulti/" +
// 							target + " deadline_exceeded: " + err.Error())...)
// 					span.SetStatus(trace.StatusError, err.Error())
// 				}
// 				return nil
// 			}
// 			st, msg, err := status.ParseError(err, codes.Internal,
// 				"failed to parse "+vald.InsertRPCName+" gRPC error response",
// 				&errdetails.RequestInfo{
// 					RequestId:   uuid,
// 					ServingData: errdetails.Serialize(req),
// 				},
// 				&errdetails.ResourceInfo{
// 					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
// 					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
// 				})
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
// 			if err != nil && st.Code() != codes.AlreadyExists {
// 				emu.Lock()
// 				if errs == nil {
// 					errs = err
// 				} else {
// 					errs = errors.Wrap(errs, err.Error())
// 				}
// 				emu.Unlock()
// 				return err
// 			}
// 			return nil
// 		}
// 		mu.Lock()
// 		ce.Ips = append(ce.GetIps(), loc.GetIps()...)
// 		ce.Name = loc.GetName()
// 		mu.Unlock()
// 		return nil
// 	})
// 	if err != nil {
// 		if errs == nil {
// 			errs = err
// 		} else {
// 			errs = errors.Wrap(errs, err.Error())
// 		}
// 	}
// 	if errs != nil {
// 		st, msg, err := status.ParseError(errs, codes.Internal,
// 			"failed to parse "+vald.InsertRPCName+" gRPC error response",
// 			&errdetails.RequestInfo{
// 				RequestId:   uuid,
// 				ServingData: errdetails.Serialize(req),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + ".DoMulti",
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	log.Debugf("Insert API insert succeeded to %#v", ce)
// 	return ce, nil
// }

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Insert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Insert_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamInsertRPCName+"/id-"+req.GetVector().GetId())
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
		})

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

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.MultiInsertRPCName), apiName+"/"+vald.MultiInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	handleSpan := func(rpcName string, span trace.Span, err error) error {
		if err == nil {
			return nil
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiInsertRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).MultiInsert(sctx, reqs, copts...)
		return handleSpan(vald.MultiInsertRPCName, sspan, err)
	})
	if err != nil {
		removeReqs := make([]*payload.Remove_Request, 0, len(reqs.Requests))
		for _, req := range reqs.Requests {
			removeReqs = append(removeReqs, &payload.Remove_Request{
				Id: &payload.Object_ID{
					Id: req.GetVector().GetId(),
				},
				Config: &payload.Remove_Config{
					SkipStrictExistCheck: true,
				},
			})
		}
		multiRemoveReq := &payload.Remove_MultiRequest{
			Requests: removeReqs,
		}
		err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiRemoveRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err := vald.NewValdClient(conn).MultiRemove(sctx, multiRemoveReq)
			return handleSpan(vald.MultiRemoveRPCName, sspan, err)
		})
		if err := handleSpan(vald.MultiRemoveRPCName, span, err); err != nil {
			return nil, err
		}
		return nil, err
	}

	ces, err := s.client.MultiInsert(ctx, reqs, s.client.GRPCClient().GetCallOption()...)
	if err := handleSpan(vald.MultiInsertRPCName, span, err); err != nil {
		return nil, err
	}
	return ces, nil
}

// func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (locs *payload.Object_Locations, err error) {
// 	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.MultiInsertRPCName), apiName+"/"+vald.MultiInsertRPCName)
// 	defer func() {
// 		if span != nil {
// 			span.End()
// 		}
// 	}()
// 	vecs := reqs.GetRequests()
// 	ids := make([]string, 0, len(vecs))
// 	now := time.Now().UnixNano()
// 	for i, req := range vecs {
// 		uuid := req.GetVector().GetId()
// 		vector := req.GetVector().GetVector()
// 		vl := len(vector)
// 		if vl < algorithm.MinimumVectorDimensionSize {
// 			err = errors.ErrInvalidDimensionSize(vl, 0)
// 			err = status.WrapWithInvalidArgument(vald.MultiInsertRPCName+" API invalid vector argument", err,
// 				&errdetails.RequestInfo{
// 					RequestId:   uuid,
// 					ServingData: errdetails.Serialize(reqs),
// 				},
// 				&errdetails.BadRequest{
// 					FieldViolations: []*errdetails.BadRequestFieldViolation{
// 						{
// 							Field:       fmt.Sprintf("vector dimension size for id: %s", uuid),
// 							Description: err.Error(),
// 						},
// 					},
// 				}, info.Get())
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
// 			return nil, err
// 		}
// 		if !req.GetConfig().GetSkipStrictExistCheck() {
// 			id, err := s.Exists(ctx, &payload.Object_ID{
// 				Id: uuid,
// 			})
// 			if err == nil && id != nil && len(id.GetId()) != 0 {
// 				if err == nil {
// 					err = errors.ErrMetaDataAlreadyExists(uuid)
// 				}
// 				st, msg, err := status.ParseError(err, codes.AlreadyExists,
// 					"error "+vald.MultiInsertRPCName+" API ID = "+uuid+" already exists",
// 					&errdetails.RequestInfo{
// 						RequestId:   uuid,
// 						ServingData: errdetails.Serialize(req),
// 					},
// 					&errdetails.ResourceInfo{
// 						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiInsertRPCName + "." + vald.ExistsRPCName,
// 						ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 					}, info.Get())
// 				if span != nil {
// 					span.RecordError(err)
// 					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 					span.SetStatus(trace.StatusError, err.Error())
// 				}
// 				return nil, err
// 			}
// 			if req.GetConfig() != nil {
// 				reqs.GetRequests()[i].GetConfig().SkipStrictExistCheck = true
// 			} else {
// 				reqs.GetRequests()[i].Config = &payload.Insert_Config{SkipStrictExistCheck: true}
// 			}
// 		}
// 		if reqs.GetRequests()[i].GetConfig().GetTimestamp() == 0 {
// 			if reqs.GetRequests()[i].GetConfig() == nil {
// 				reqs.GetRequests()[i].Config = &payload.Insert_Config{
// 					Timestamp: now,
// 				}
// 			} else {
// 				reqs.GetRequests()[i].GetConfig().Timestamp = now
// 			}
// 		}
// 		ids = append(ids, uuid)
// 	}
//
// 	mu := new(sync.Mutex)
// 	locs = &payload.Object_Locations{
// 		Locations: make([]*payload.Object_Location, 0, s.replica),
// 	}
//
// 	emu := new(sync.Mutex)
// 	var errs error
// 	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
// 		ctx, span := trace.StartSpan(ctx, apiName+"."+vald.MultiInsertRPCName+"/"+target)
// 		defer func() {
// 			if span != nil {
// 				span.End()
// 			}
// 		}()
// 		loc, err := vc.MultiInsert(ctx, reqs, copts...)
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, context.Canceled),
// 				errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
// 				if span != nil {
// 					span.RecordError(err)
// 					span.SetAttributes(trace.StatusCodeCancelled(
// 						errdetails.ValdGRPCResourceTypePrefix +
// 							"/vald.v1." + vald.MultiInsertRPCName + ".DoMulti/" +
// 							target + " canceled: " + err.Error())...)
// 					span.SetStatus(trace.StatusError, err.Error())
// 				}
// 				return nil
// 			case errors.Is(err, context.DeadlineExceeded),
// 				errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
// 				if span != nil {
// 					span.RecordError(err)
// 					span.SetAttributes(trace.StatusCodeDeadlineExceeded(
// 						errdetails.ValdGRPCResourceTypePrefix +
// 							"/vald.v1." + vald.MultiInsertRPCName + ".DoMulti/" +
// 							target + " deadline_exceeded: " + err.Error())...)
// 					span.SetStatus(trace.StatusError, err.Error())
// 				}
// 				return nil
// 			}
// 			st, msg, err := status.ParseError(err, codes.Internal,
// 				"failed to parse "+vald.MultiInsertRPCName+" gRPC error response",
// 				&errdetails.RequestInfo{
// 					RequestId:   strings.Join(ids, ","),
// 					ServingData: errdetails.Serialize(reqs),
// 				},
// 				&errdetails.ResourceInfo{
// 					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiInsertRPCName,
// 					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
// 				})
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
//
// 			if err != nil {
// 				emu.Lock()
// 				if errs == nil {
// 					errs = err
// 				} else {
// 					errs = errors.Wrap(errs, err.Error())
// 				}
// 				emu.Unlock()
// 			}
// 			return err
// 		}
// 		mu.Lock()
// 		locs.Locations = append(locs.GetLocations(), loc.Locations...)
// 		mu.Unlock()
// 		return nil
// 	})
// 	if err != nil {
// 		if errs == nil {
// 			errs = err
// 		} else {
// 			errs = errors.Wrap(errs, err.Error())
// 		}
// 	}
//
// 	if errs != nil {
// 		st, msg, err := status.ParseError(errs, codes.Internal,
// 			"failed to parse "+vald.MultiInsertRPCName+" gRPC error response",
// 			&errdetails.RequestInfo{
// 				RequestId:   strings.Join(ids, ", "),
// 				ServingData: errdetails.Serialize(reqs),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiInsertRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	return location.ReStructure(ids, locs), nil
// }

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	handleSpan := func(rpcName string, span trace.Span, err error) error {
		if err == nil {
			return nil
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.UpdateRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).Update(sctx, req)
		return handleSpan(vald.UpdateRPCName, sspan, err)
	})
	if err != nil {
		ce, err := s.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: req.GetVector().GetId(),
			},
		})
		if err = handleSpan(vald.UpdateRPCName+" for "+vald.GetObjectRPCName, span, err); err != nil {
			return nil, err
		}

		newReq := &payload.Update_Request{
			Vector: &payload.Object_Vector{
				Id:     ce.GetId(),
				Vector: ce.GetVector(),
			},
			Config: &payload.Update_Config{
				SkipStrictExistCheck: true,
			},
		}
		err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.UpdateRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err := vald.NewValdClient(conn).Update(sctx, newReq)
			return handleSpan(vald.UpdateRPCName, sspan, err)
		})
		if err := handleSpan(vald.UpdateRPCName, span, err); err != nil {
			return nil, err
		}
		return nil, err
	}

	ce, err := s.client.Update(ctx, req, s.client.GRPCClient().GetCallOption()...)
	if err = handleSpan(vald.UpdateRPCName, span, err); err != nil {
		return nil, err
	}
	return ce, nil
}

// func (s *server) Update(ctx context.Context, req *payload.Update_Request) (res *payload.Object_Location, err error) {
// 	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.UpdateRPCName)
// 	defer func() {
// 		if span != nil {
// 			span.End()
// 		}
// 	}()
// 	uuid := req.GetVector().GetId()
// 	if len(uuid) == 0 {
// 		err = errors.ErrInvalidMetaDataConfig
// 		err = status.WrapWithInvalidArgument(vald.UpdateRPCName+" API invalid uuid", err,
// 			&errdetails.RequestInfo{
// 				ServingData: errdetails.Serialize(req),
// 			},
// 			&errdetails.BadRequest{
// 				FieldViolations: []*errdetails.BadRequestFieldViolation{
// 					{
// 						Field:       "invalid id",
// 						Description: err.Error(),
// 					},
// 				},
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
// 			})
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	vec := req.GetVector().GetVector()
// 	vl := len(vec)
// 	if vl < algorithm.MinimumVectorDimensionSize {
// 		err = errors.ErrInvalidDimensionSize(vl, 0)
// 		err = status.WrapWithInvalidArgument(vald.UpdateRPCName+" API invalid vector argument", err,
// 			&errdetails.RequestInfo{
// 				RequestId:   uuid,
// 				ServingData: errdetails.Serialize(req),
// 			},
// 			&errdetails.BadRequest{
// 				FieldViolations: []*errdetails.BadRequestFieldViolation{
// 					{
// 						Field:       "vector dimension size",
// 						Description: err.Error(),
// 					},
// 				},
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
//
// 	if !req.GetConfig().GetSkipStrictExistCheck() {
// 		vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
// 			Id: &payload.Object_ID{
// 				Id: uuid,
// 			},
// 		})
// 		if err != nil || vec == nil || len(vec.GetId()) == 0 {
// 			if err == nil {
// 				err = errors.ErrObjectIDNotFound(uuid)
// 			}
// 			st, msg, err := status.ParseError(err, codes.NotFound,
// 				"error "+vald.UpdateRPCName+" API ID = "+uuid+" not found",
// 				&errdetails.RequestInfo{
// 					RequestId:   uuid,
// 					ServingData: errdetails.Serialize(req),
// 				},
// 				&errdetails.ResourceInfo{
// 					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.GetObjectRPCName,
// 					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 				}, info.Get())
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
// 			return nil, err
// 		}
// 		if conv.F32stos(vec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
// 			if err == nil {
// 				err = errors.ErrSameVectorAlreadyExists(uuid, vec.GetVector(), req.GetVector().GetVector())
// 			}
// 			st, msg, err := status.ParseError(err, codes.AlreadyExists,
// 				"error "+vald.UpdateRPCName+" API ID = "+uuid+"'s same vector data already exists",
// 				&errdetails.RequestInfo{
// 					RequestId:   uuid,
// 					ServingData: errdetails.Serialize(req),
// 				},
// 				&errdetails.ResourceInfo{
// 					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.GetObjectRPCName,
// 					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 				}, info.Get())
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
// 			return nil, err
// 		}
// 		if req.GetConfig() != nil {
// 			req.GetConfig().SkipStrictExistCheck = true
// 		} else {
// 			req.Config = &payload.Update_Config{SkipStrictExistCheck: true}
// 		}
// 	}
// 	var now int64
// 	if req.GetConfig().GetTimestamp() != 0 {
// 		now = req.GetConfig().GetTimestamp()
// 	} else {
// 		now = time.Now().UnixNano()
// 	}
//
// 	rreq := &payload.Remove_Request{
// 		Id: &payload.Object_ID{
// 			Id: uuid,
// 		},
// 		Config: &payload.Remove_Config{
// 			SkipStrictExistCheck: true,
// 			Timestamp:            now,
// 		},
// 	}
// 	res, err = s.Remove(ctx, rreq)
// 	if err != nil {
// 		st, msg, err := status.ParseError(err, codes.Internal,
// 			"failed to parse "+vald.RemoveRPCName+" for "+vald.UpdateRPCName+" gRPC error response",
// 			&errdetails.RequestInfo{
// 				RequestId:   uuid,
// 				ServingData: errdetails.Serialize(rreq),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.RemoveRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	now++
// 	ireq := &payload.Insert_Request{
// 		Vector: req.GetVector(),
// 		Config: &payload.Insert_Config{
// 			SkipStrictExistCheck: true,
// 			Filters:              req.GetConfig().GetFilters(),
// 			Timestamp:            now,
// 		},
// 	}
// 	res, err = s.Insert(ctx, ireq)
// 	if err != nil {
// 		st, msg, err := status.ParseError(err, codes.Internal,
// 			"failed to parse "+vald.InsertRPCName+" for "+vald.UpdateRPCName+" gRPC error response",
// 			&errdetails.RequestInfo{
// 				RequestId:   uuid,
// 				ServingData: errdetails.Serialize(ireq),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.InsertRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
//
// 		return nil, err
// 	}
// 	return res, nil
// }

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Update_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Update_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamUpdateRPCName+"/id-"+req.GetVector().GetId())
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
		})

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

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	handleSpan := func(rpcName string, span trace.Span, err error) error {
		if err == nil {
			return nil
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiUpdateRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).MultiUpdate(sctx, reqs)
		return handleSpan(vald.MultiUpdateRPCName, sspan, err)
	})
	if err != nil {
		mu := new(sync.Mutex)
		newReqs := &payload.Update_MultiRequest{
			Requests: make([]*payload.Update_Request, 0, len(reqs.GetRequests())),
		}
		eg, egctx := errgroup.New(ctx)

		for _, req := range reqs.GetRequests() {
			req := req
			eg.Go(func() error {
				ovec, err := s.GetObject(egctx, &payload.Object_VectorRequest{
					Id: &payload.Object_ID{
						Id: req.GetVector().GetId(),
					},
				})
				if err != nil {
					return err
				}
				mu.Lock()
				newReqs.Requests = append(newReqs.Requests, &payload.Update_Request{
					Vector: &payload.Object_Vector{
						Id:     ovec.GetId(),
						Vector: ovec.GetVector(),
					},
					Config: &payload.Update_Config{
						SkipStrictExistCheck: true,
					},
				})
				mu.Unlock()
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			return nil, err
		}

		err := s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiUpdateRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err := vald.NewValdClient(conn).MultiUpdate(sctx, newReqs, copts...)
			return handleSpan(vald.MultiUpdateRPCName, sspan, err)
		})
		if err := handleSpan(vald.MultiUpdateRPCName, span, err); err != nil {
			return nil, err
		}
		return nil, err
	}

	ces, err := s.client.MultiUpdate(ctx, reqs, s.client.GRPCClient().GetCallOption()...)
	if err = handleSpan(vald.MultiUpdateRPCName, span, err); err != nil {
		return nil, err
	}
	return ces, nil
}

// func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (res *payload.Object_Locations, err error) {
// 	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiUpdateRPCName)
// 	defer func() {
// 		if span != nil {
// 			span.End()
// 		}
// 	}()
// 	vecs := reqs.GetRequests()
// 	ids := make([]string, 0, len(vecs))
// 	ireqs := make([]*payload.Insert_Request, 0, len(vecs))
// 	rreqs := make([]*payload.Remove_Request, 0, len(vecs))
// 	now := time.Now().UnixNano()
// 	for _, req := range vecs {
// 		vl := len(req.GetVector().GetVector())
// 		if vl < algorithm.MinimumVectorDimensionSize {
// 			err = errors.ErrInvalidDimensionSize(vl, 0)
// 			err = status.WrapWithInvalidArgument(vald.MultiUpdateRPCName+" API invalid vector argument", err,
// 				&errdetails.RequestInfo{
// 					RequestId:   req.GetVector().GetId(),
// 					ServingData: errdetails.Serialize(reqs),
// 				},
// 				&errdetails.BadRequest{
// 					FieldViolations: []*errdetails.BadRequestFieldViolation{
// 						{
// 							Field:       "vector dimension size",
// 							Description: err.Error(),
// 						},
// 					},
// 				}, info.Get())
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
// 			return nil, err
// 		}
// 		uuid := req.GetVector().GetId()
// 		if !req.GetConfig().GetSkipStrictExistCheck() {
// 			vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
// 				Id: &payload.Object_ID{
// 					Id: uuid,
// 				},
// 			})
// 			if err != nil || vec == nil || len(vec.GetId()) == 0 {
// 				if err == nil {
// 					err = errors.ErrObjectIDNotFound(uuid)
// 				}
// 				st, msg, err := status.ParseError(err, codes.NotFound,
// 					"error "+vald.MultiUpdateRPCName+" API ID = "+uuid+" not fount",
// 					&errdetails.RequestInfo{
// 						RequestId:   uuid,
// 						ServingData: errdetails.Serialize(reqs),
// 					},
// 					&errdetails.ResourceInfo{
// 						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpdateRPCName + "." + vald.GetObjectRPCName,
// 						ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 					}, info.Get())
// 				if span != nil {
// 					span.RecordError(err)
// 					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 					span.SetStatus(trace.StatusError, err.Error())
// 				}
// 				return nil, err
// 			}
// 			if conv.F32stos(vec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
// 				log.Warn(errors.ErrSameVectorAlreadyExists(uuid, vec.GetVector(), req.GetVector().GetVector()))
// 				continue
// 			}
// 			if req.GetConfig() != nil {
// 				req.Config.SkipStrictExistCheck = true
// 			} else {
// 				req.Config = &payload.Update_Config{SkipStrictExistCheck: true}
// 			}
// 		}
// 		var n int64
// 		if req.GetConfig().GetTimestamp() != 0 {
// 			n = req.GetConfig().GetTimestamp()
// 		} else {
// 			n = now
// 		}
// 		ids = append(ids, req.GetVector().GetId())
// 		rreqs = append(rreqs, &payload.Remove_Request{
// 			Id: &payload.Object_ID{
// 				Id: req.GetVector().GetId(),
// 			},
// 			Config: &payload.Remove_Config{
// 				SkipStrictExistCheck: true,
// 				Timestamp:            n,
// 			},
// 		})
// 		n++
// 		ireqs = append(ireqs, &payload.Insert_Request{
// 			Vector: req.GetVector(),
// 			Config: &payload.Insert_Config{
// 				SkipStrictExistCheck: true,
// 				Filters:              req.GetConfig().GetFilters(),
// 				Timestamp:            n,
// 			},
// 		})
// 	}
// 	locs, err := s.MultiRemove(ctx, &payload.Remove_MultiRequest{
// 		Requests: rreqs,
// 	})
// 	if err != nil {
// 		st, msg, err := status.ParseError(err, codes.Internal,
// 			"failed to parse "+vald.MultiRemoveRPCName+" for "+vald.MultiUpdateRPCName+" gRPC error response",
// 			&errdetails.RequestInfo{
// 				RequestId:   strings.Join(ids, ","),
// 				ServingData: errdetails.Serialize(rreqs),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpdateRPCName + "." + vald.MultiRemoveRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	log.Debugf("uuids %v were removed from %v due to MultiUpdate. "+vald.MultiInsertRPCName+" will be executed for them soon. Please see detail %#v", ids, locs.GetLocations(), locs)
// 	locs, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
// 		Requests: ireqs,
// 	})
// 	if err != nil {
// 		st, msg, err := status.ParseError(err, codes.Internal,
// 			"failed to parse "+vald.MultiInsertRPCName+" for "+vald.MultiUpdateRPCName+" gRPC error response",
// 			&errdetails.RequestInfo{
// 				RequestId:   strings.Join(ids, ","),
// 				ServingData: errdetails.Serialize(ireqs),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpdateRPCName + "." + vald.MultiInsertRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	return locs, nil
// }

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	handleSpan := func(rpcName string, span trace.Span, err error) error {
		if err == nil {
			return nil
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.UpsertRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).Upsert(sctx, req, copts...)
		return handleSpan(vald.UpsertRPCName, sspan, err)
	})
	if err != nil {
		if err := s.upsertRollback(ctx, req); err != nil {
			return nil, err
		}
		return nil, err
	}

	ce, err := s.client.Upsert(ctx, req, s.client.GRPCClient().GetCallOption()...)
	if err := handleSpan(vald.UpsertRPCName, span, err); err != nil {
		return nil, err
	}
	return ce, nil
}

func (s *server) upsertRollback(ctx context.Context, req *payload.Upsert_Request) error {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	handleSpan := func(rpcName string, span trace.Span, err error) error {
		if err == nil {
			return nil
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	oid := &payload.Object_ID{
		Id: req.GetVector().GetId(),
	}
	ce, err := s.GetObject(ctx, &payload.Object_VectorRequest{
		Id: oid,
	})
	if err != nil {
		st, _, err := status.ParseError(err, codes.Internal, "error "+vald.GetObjectRPCName+" API")
		if err != nil && st.Code() == codes.NotFound {
			removeReq := &payload.Remove_Request{
				Id: oid,
			}
			err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
				sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.RemoveRPCName+"/"+target)
				defer func() {
					if sspan != nil {
						sspan.End()
					}
				}()
				_, err := vald.NewValdClient(conn).Remove(sctx, removeReq, copts...)
				return handleSpan(vald.RemoveRPCName, sspan, err)
			})
			if err := handleSpan(vald.RemoveRPCName, span, err); err != nil {
				return err
			}
			return err
		}
		return err
	}
	updateReq := &payload.Update_Request{
		Vector: &payload.Object_Vector{
			Id:     ce.GetId(),
			Vector: ce.GetVector(),
		},
		Config: &payload.Update_Config{
			SkipStrictExistCheck: true,
		},
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.UpdateRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).Update(sctx, updateReq, copts...)
		return handleSpan(vald.UpdateRPCName, sspan, err)
	})
	if err := handleSpan(vald.UpdateRPCName, span, err); err != nil {
		return err
	}
	return nil
}

// func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
// 	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.UpsertRPCName)
// 	defer func() {
// 		if span != nil {
// 			span.End()
// 		}
// 	}()
//
// 	vec := req.GetVector()
// 	uuid := vec.GetId()
// 	if len(uuid) == 0 {
// 		err = errors.ErrInvalidMetaDataConfig
// 		err = status.WrapWithInvalidArgument(vald.UpsertRPCName+" API invalid uuid", err,
// 			&errdetails.RequestInfo{
// 				ServingData: errdetails.Serialize(req),
// 			},
// 			&errdetails.BadRequest{
// 				FieldViolations: []*errdetails.BadRequestFieldViolation{
// 					{
// 						Field:       "invalid id",
// 						Description: err.Error(),
// 					},
// 				},
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
// 			})
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
//
// 	vl := len(vec.GetVector())
// 	if vl < algorithm.MinimumVectorDimensionSize {
// 		err = errors.ErrInvalidDimensionSize(vl, 0)
// 		err = status.WrapWithInvalidArgument(vald.UpsertRPCName+" API invalid vector argument", err,
// 			&errdetails.RequestInfo{
// 				RequestId:   uuid,
// 				ServingData: errdetails.Serialize(req),
// 			},
// 			&errdetails.BadRequest{
// 				FieldViolations: []*errdetails.BadRequestFieldViolation{
// 					{
// 						Field:       "vector dimension size",
// 						Description: err.Error(),
// 					},
// 				},
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	var shouldInsert bool
// 	if !req.GetConfig().GetSkipStrictExistCheck() {
// 		vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
// 			Id: &payload.Object_ID{
// 				Id: uuid,
// 			},
// 		})
// 		if err != nil || vec == nil || len(vec.GetId()) == 0 {
// 			shouldInsert = true
// 		} else if conv.F32stos(vec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
// 			if err == nil {
// 				err = errors.ErrSameVectorAlreadyExists(uuid, vec.GetVector(), req.GetVector().GetVector())
// 			}
// 			st, msg, err := status.ParseError(err, codes.AlreadyExists,
// 				"error "+vald.UpdateRPCName+" for "+vald.UpsertRPCName+" API ID = "+uuid+"'s same vector data already exists",
// 				&errdetails.RequestInfo{
// 					RequestId:   uuid,
// 					ServingData: errdetails.Serialize(req),
// 				},
// 				&errdetails.ResourceInfo{
// 					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + vald.GetObjectRPCName,
// 					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 				}, info.Get())
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
// 			return nil, err
// 		}
// 	} else {
// 		id, err := s.Exists(ctx, &payload.Object_ID{
// 			Id: uuid,
// 		})
// 		shouldInsert = err != nil || id == nil || len(id.GetId()) == 0
// 	}
//
// 	var operation string
// 	if shouldInsert {
// 		operation = vald.InsertRPCName
// 		loc, err = s.Insert(ctx, &payload.Insert_Request{
// 			Vector: vec,
// 			Config: &payload.Insert_Config{
// 				SkipStrictExistCheck: true,
// 				Filters:              req.GetConfig().GetFilters(),
// 				Timestamp:            req.GetConfig().GetTimestamp(),
// 			},
// 		})
// 	} else {
// 		operation = vald.UpdateRPCName
// 		loc, err = s.Update(ctx, &payload.Update_Request{
// 			Vector: vec,
// 			Config: &payload.Update_Config{
// 				SkipStrictExistCheck: true,
// 				Filters:              req.GetConfig().GetFilters(),
// 				Timestamp:            req.GetConfig().GetTimestamp(),
// 			},
// 		})
// 	}
//
// 	if err != nil {
// 		st, msg, err := status.ParseError(err, codes.Internal,
// 			"failed to parse "+operation+" for "+vald.UpsertRPCName+" gRPC error response",
// 			&errdetails.RequestInfo{
// 				RequestId:   uuid,
// 				ServingData: errdetails.Serialize(req),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + operation,
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	return loc, nil
// }

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Upsert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Upsert_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamUpsertRPCName+"/id-"+req.GetVector().GetId())
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
		})

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

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	handleSpan := func(rpcName string, span trace.Span, err error) error {
		if err == nil {
			return nil
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.UpsertRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).MultiUpsert(sctx, reqs, copts...)
		return handleSpan(vald.UpsertRPCName, sspan, err)
	})
	if err != nil {
		eg, egctx := errgroup.New(ctx)
		for _, req := range reqs.Requests {
			req := req
			eg.Go(func() error {
				return s.upsertRollback(egctx, req)
			})
		}
		if err := eg.Wait(); err != nil {
			return nil, err
		}
		return nil, err
	}

	res, err = s.client.MultiUpsert(ctx, reqs, s.client.GRPCClient().GetCallOption()...)
	if err := handleSpan(vald.MultiUpsertRPCName, span, err); err != nil {
		return nil, err
	}
	return res, nil
}

// func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (res *payload.Object_Locations, err error) {
// 	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiUpsertRPCName)
// 	defer func() {
// 		if span != nil {
// 			span.End()
// 		}
// 	}()
//
// 	insertReqs := make([]*payload.Insert_Request, 0, len(reqs.GetRequests()))
// 	updateReqs := make([]*payload.Update_Request, 0, len(reqs.GetRequests()))
//
// 	ids := make([]string, 0, len(reqs.GetRequests()))
// 	for _, req := range reqs.GetRequests() {
// 		vec := req.GetVector()
// 		uuid := vec.GetId()
// 		vl := len(vec.GetVector())
// 		if vl < algorithm.MinimumVectorDimensionSize {
// 			err = errors.ErrInvalidDimensionSize(vl, 0)
// 			err = status.WrapWithInvalidArgument(vald.MultiUpsertRPCName+" API invalid vector argument", err,
// 				&errdetails.RequestInfo{
// 					RequestId:   uuid,
// 					ServingData: errdetails.Serialize(req),
// 				},
// 				&errdetails.BadRequest{
// 					FieldViolations: []*errdetails.BadRequestFieldViolation{
// 						{
// 							Field:       "vector dimension size",
// 							Description: err.Error(),
// 						},
// 					},
// 				}, info.Get())
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
// 			return nil, err
// 		}
// 		var shouldInsert bool
// 		if !req.GetConfig().GetSkipStrictExistCheck() {
// 			vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
// 				Id: &payload.Object_ID{
// 					Id: uuid,
// 				},
// 			})
// 			if err != nil || vec == nil || len(vec.GetId()) == 0 {
// 				shouldInsert = true
// 			} else if conv.F32stos(vec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
// 				continue
// 			}
// 		} else {
// 			id, err := s.Exists(ctx, &payload.Object_ID{
// 				Id: uuid,
// 			})
// 			shouldInsert = err != nil || id == nil || len(id.GetId()) == 0
// 		}
// 		ids = append(ids, uuid)
// 		if shouldInsert {
// 			insertReqs = append(insertReqs, &payload.Insert_Request{
// 				Vector: vec,
// 				Config: &payload.Insert_Config{
// 					SkipStrictExistCheck: true,
// 					Filters:              req.GetConfig().GetFilters(),
// 					Timestamp:            req.GetConfig().GetTimestamp(),
// 				},
// 			})
// 		} else {
// 			updateReqs = append(updateReqs, &payload.Update_Request{
// 				Vector: vec,
// 				Config: &payload.Update_Config{
// 					SkipStrictExistCheck: true,
// 					Filters:              req.GetConfig().GetFilters(),
// 					Timestamp:            req.GetConfig().GetTimestamp(),
// 				},
// 			})
// 		}
// 	}
//
// 	switch {
// 	case len(insertReqs) <= 0:
// 		res, err = s.MultiUpdate(ctx, &payload.Update_MultiRequest{
// 			Requests: updateReqs,
// 		})
// 	case len(updateReqs) <= 0:
// 		res, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
// 			Requests: insertReqs,
// 		})
// 	default:
// 		var (
// 			ures, ires *payload.Object_Locations
// 			errs       error
// 			mu         sync.Mutex
// 		)
// 		eg, ectx := errgroup.New(ctx)
// 		eg.Go(safety.RecoverFunc(func() (err error) {
// 			ures, err = s.MultiUpdate(ectx, &payload.Update_MultiRequest{
// 				Requests: updateReqs,
// 			})
// 			if err != nil {
// 				mu.Lock()
// 				if errs == nil {
// 					errs = err
// 				} else {
// 					errs = errors.Wrap(errs, err.Error())
// 				}
// 				mu.Unlock()
// 			}
// 			return nil
// 		}))
// 		eg.Go(safety.RecoverFunc(func() (err error) {
// 			ires, err = s.MultiInsert(ectx, &payload.Insert_MultiRequest{
// 				Requests: insertReqs,
// 			})
// 			if err != nil {
// 				mu.Lock()
// 				if errs == nil {
// 					errs = err
// 				} else {
// 					errs = errors.Wrap(errs, err.Error())
// 				}
// 				mu.Unlock()
// 			}
// 			return nil
// 		}))
// 		err = eg.Wait()
// 		if err != nil {
// 			if errs == nil {
// 				errs = err
// 			} else {
// 				errs = errors.Wrap(errs, err.Error())
// 			}
// 		}
// 		if errs != nil {
// 			err = errs
// 		}
// 		switch {
// 		case ures.GetLocations() == nil && ires.GetLocations() != nil:
// 			res = ires
// 		case ures.GetLocations() != nil && ires.GetLocations() == nil:
// 			res = ures
// 		case ures.GetLocations() != nil && ires.GetLocations() != nil:
// 			res = &payload.Object_Locations{
// 				Locations: append(ures.GetLocations(), ires.GetLocations()...),
// 			}
// 		default:
// 			res = new(payload.Object_Locations)
// 		}
//
// 	}
//
// 	if err != nil {
// 		st, msg, err := status.ParseError(err, codes.Internal,
// 			"failed to parse "+vald.MultiUpsertRPCName+" gRPC error response",
// 			&errdetails.RequestInfo{
// 				RequestId:   strings.Join(ids, ","),
// 				ServingData: errdetails.Serialize(reqs),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpsertRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	return location.ReStructure(ids, res), nil
// }

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	handleSpan := func(rpcName string, span trace.Span, err error) error {
		if err == nil {
			return nil
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.RemoveRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).Remove(sctx, req, copts...)
		return handleSpan(vald.RemoveRPCName, sspan, err)
	})
	if err != nil {
		if err := handleSpan(vald.RemoveRPCName, span, s.removeRollback(ctx, req)); err != nil {
			return nil, err
		}
		return nil, err
	}

	loc, err = s.client.Remove(ctx, req, s.client.GRPCClient().GetCallOption()...)
	if err := handleSpan(vald.RemoveRPCName, span, err); err != nil {
		return nil, err
	}
	return loc, nil
}

func (s *server) removeRollback(ctx context.Context, req *payload.Remove_Request) (err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	handleSpan := func(rpcName string, span trace.Span, err error) error {
		if err == nil {
			return nil
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}

	loc, err := s.GetObject(ctx, &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetId().GetId(),
		},
	})
	if err = handleSpan(vald.RemoveRPCName, span, err); err != nil {
		return err
	}

	upsertReq := &payload.Upsert_Request{
		Vector: &payload.Object_Vector{
			Id:     loc.GetId(),
			Vector: loc.GetVector(),
		},
		Config: &payload.Upsert_Config{
			SkipStrictExistCheck: true,
		},
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.UpsertRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()

		_, err := vald.NewValdClient(conn).Upsert(sctx, upsertReq, copts...)
		return handleSpan(vald.UpsertRPCName, sspan, err)
	})
	if err = handleSpan(vald.RemoveRPCName, span, err); err != nil {
		return err
	}
	return nil
}

// func (s *server) _Remove(ctx context.Context, req *payload.Remove_Request) (locs *payload.Object_Location, err error) {
// 	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
// 	defer func() {
// 		if span != nil {
// 			span.End()
// 		}
// 	}()
//
// 	id := req.GetId()
// 	if !req.GetConfig().GetSkipStrictExistCheck() {
// 		id, err := s.Exists(ctx, id)
// 		if err != nil || id == nil || len(id.GetId()) == 0 {
// 			if err == nil {
// 				err = errors.ErrObjectIDNotFound(id.GetId())
// 			}
// 			st, msg, err := status.ParseError(err, codes.NotFound,
// 				"error "+vald.RemoveRPCName+" API ID = "+id.GetId()+" not found",
// 				&errdetails.RequestInfo{
// 					RequestId:   id.GetId(),
// 					ServingData: errdetails.Serialize(req),
// 				},
// 				&errdetails.ResourceInfo{
// 					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + "." + vald.ExistsRPCName,
// 					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 				}, info.Get())
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
// 			return nil, err
// 		}
// 		if req.GetConfig() != nil {
// 			req.GetConfig().SkipStrictExistCheck = true
// 		} else {
// 			req.Config = &payload.Remove_Config{SkipStrictExistCheck: true}
// 		}
// 	}
// 	if req.GetConfig().GetTimestamp() == 0 {
// 		now := time.Now().UnixNano()
// 		if req.GetConfig() == nil {
// 			req.Config = &payload.Remove_Config{
// 				Timestamp: now,
// 			}
// 		} else {
// 			req.GetConfig().Timestamp = now
// 		}
// 	}
// 	var mu sync.Mutex
// 	locs = &payload.Object_Location{
// 		Uuid: id.GetId(),
// 		Ips:  make([]string, 0, s.replica),
// 	}
// 	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
// 		ctx, span := trace.StartSpan(ctx, apiName+"."+vald.RemoveRPCName+"/"+target)
// 		defer func() {
// 			if span != nil {
// 				span.End()
// 			}
// 		}()
// 		loc, err := vc.Remove(ctx, req, copts...)
// 		if err != nil {
// 			st, msg, err := status.ParseError(err, codes.Internal,
// 				"failed to parse "+vald.RemoveRPCName+" gRPC error response",
// 				&errdetails.RequestInfo{
// 					RequestId:   id.GetId(),
// 					ServingData: errdetails.Serialize(req),
// 				},
// 				&errdetails.ResourceInfo{
// 					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
// 					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
// 				})
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
// 			if err != nil && st.Code() != codes.NotFound {
// 				log.Error(err)
// 				return err
// 			}
// 			return nil
// 		}
// 		mu.Lock()
// 		locs.Ips = append(locs.GetIps(), loc.GetIps()...)
// 		locs.Name = loc.GetName()
// 		mu.Unlock()
// 		return nil
// 	})
// 	if err != nil {
// 		st, msg, err := status.ParseError(err, codes.Internal,
// 			"failed to parse "+vald.RemoveRPCName+" gRPC error response",
// 			&errdetails.RequestInfo{
// 				RequestId:   id.GetId(),
// 				ServingData: errdetails.Serialize(req),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	if len(locs.Ips) <= 0 {
// 		err = errors.ErrIndexNotFound
// 		err = status.WrapWithNotFound(vald.RemoveRPCName+" API remove target not found", err,
// 			&errdetails.RequestInfo{
// 				RequestId:   id.GetId(),
// 				ServingData: errdetails.Serialize(req),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			})
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	return locs, nil
// }

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Remove_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Remove_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamRemoveRPCName+"/id-"+req.GetId().GetId())
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
		})

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

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	handleSpan := func(rpcName string, span trace.Span, err error) error {
		if err == nil {
			return nil
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+rpcName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.RemoveRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).MultiRemove(sctx, reqs, copts...)
		return handleSpan(vald.RemoveRPCName, sspan, err)
	})
	if err != nil {
		eg, egctx := errgroup.New(ctx)

		for _, req := range reqs.GetRequests() {
			req := req
			eg.Go(func() error {
				return s.removeRollback(egctx, req)
			})
		}
		if err := eg.Wait(); err != nil {
			return nil, handleSpan(vald.RemoveRPCName, span, err)
		}
		return nil, handleSpan(vald.RemoveRPCName, span, err)
	}

	locs, err = s.client.MultiRemove(ctx, reqs, s.client.GRPCClient().GetCallOption()...)
	if err := handleSpan(vald.MultiRemoveRPCName, span, err); err != nil {
		return nil, err
	}
	return locs, nil
}

// func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (locs *payload.Object_Locations, err error) {
// 	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
// 	defer func() {
// 		if span != nil {
// 			span.End()
// 		}
// 	}()
//
// 	now := time.Now().UnixNano()
// 	ids := make([]string, 0, len(reqs.GetRequests()))
// 	for i, req := range reqs.GetRequests() {
// 		id := req.GetId()
// 		ids = append(ids, id.GetId())
// 		if !req.GetConfig().GetSkipStrictExistCheck() {
// 			sid, err := s.Exists(ctx, id)
// 			if err != nil || sid == nil || len(sid.GetId()) == 0 {
// 				if err == nil {
// 					err = errors.ErrObjectIDNotFound(id.GetId())
// 				}
// 				st, msg, err := status.ParseError(err, codes.NotFound,
// 					fmt.Sprintf(vald.MultiRemoveRPCName+" API ID = %v not found", id.GetId()),
// 					&errdetails.RequestInfo{
// 						RequestId:   id.GetId(),
// 						ServingData: errdetails.Serialize(reqs),
// 					},
// 					&errdetails.ResourceInfo{
// 						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName + "." + vald.ExistsRPCName,
// 						ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 					}, info.Get())
// 				if span != nil {
// 					span.RecordError(err)
// 					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 					span.SetStatus(trace.StatusError, err.Error())
// 				}
// 				return nil, err
// 			}
// 			if reqs.GetRequests()[i].GetConfig() != nil {
// 				reqs.GetRequests()[i].GetConfig().SkipStrictExistCheck = true
// 			} else {
// 				reqs.GetRequests()[i].Config = &payload.Remove_Config{SkipStrictExistCheck: true}
// 			}
//
// 		}
// 		if req.GetConfig().GetTimestamp() == 0 {
// 			if req.GetConfig() == nil {
// 				reqs.GetRequests()[i].Config = &payload.Remove_Config{
// 					Timestamp: now,
// 				}
// 			} else {
// 				reqs.GetRequests()[i].GetConfig().Timestamp = now
// 			}
// 		}
// 	}
// 	var mu sync.Mutex
// 	locs = &payload.Object_Locations{
// 		Locations: make([]*payload.Object_Location, 0, len(reqs.GetRequests())),
// 	}
// 	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
// 		ctx, span := trace.StartSpan(ctx, apiName+"."+vald.MultiRemoveRPCName+"/"+target)
// 		defer func() {
// 			if span != nil {
// 				span.End()
// 			}
// 		}()
// 		loc, err := vc.MultiRemove(ctx, reqs, copts...)
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, context.Canceled),
// 				errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
// 				if span != nil {
// 					span.RecordError(err)
// 					span.SetAttributes(trace.StatusCodeCancelled(
// 						errdetails.ValdGRPCResourceTypePrefix +
// 							"/vald.v1." + vald.MultiRemoveRPCName + ".BroadCast/" +
// 							target + " canceled: " + err.Error())...)
// 					span.SetStatus(trace.StatusError, err.Error())
// 				}
// 				return nil
// 			case errors.Is(err, context.DeadlineExceeded),
// 				errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
// 				if span != nil {
// 					span.RecordError(err)
// 					span.SetAttributes(trace.StatusCodeDeadlineExceeded(
// 						errdetails.ValdGRPCResourceTypePrefix +
// 							"/vald.v1." + vald.MultiRemoveRPCName + ".BroadCast/" +
// 							target + " deadline_exceeded: " + err.Error())...)
// 					span.SetStatus(trace.StatusError, err.Error())
// 				}
// 				return nil
// 			}
// 			st, msg, err := status.ParseError(err, codes.Internal,
// 				"failed to parse MultiRemove gRPC error response",
// 				&errdetails.RequestInfo{
// 					RequestId:   strings.Join(ids, ","),
// 					ServingData: errdetails.Serialize(reqs),
// 				},
// 				&errdetails.ResourceInfo{
// 					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName,
// 					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
// 				})
// 			if span != nil {
// 				span.RecordError(err)
// 				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 				span.SetStatus(trace.StatusError, err.Error())
// 			}
//
// 			if err != nil && st.Code() != codes.NotFound {
// 				log.Error(err)
// 				return err
// 			}
// 			return nil
// 		}
// 		mu.Lock()
// 		locs.Locations = append(locs.GetLocations(), loc.GetLocations()...)
// 		mu.Unlock()
// 		return nil
// 	})
// 	if err != nil {
// 		st, msg, err := status.ParseError(err, codes.Internal,
// 			"failed to parse "+vald.MultiRemoveRPCName+" gRPC error response",
// 			&errdetails.RequestInfo{
// 				RequestId:   strings.Join(ids, ","),
// 				ServingData: errdetails.Serialize(reqs),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			}, info.Get())
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	if len(locs.Locations) <= 0 {
// 		err = errors.ErrIndexNotFound
// 		err = status.WrapWithNotFound(vald.MultiRemoveRPCName+" API remove target not found", err,
// 			&errdetails.RequestInfo{
// 				RequestId:   strings.Join(ids, ","),
// 				ServingData: errdetails.Serialize(reqs),
// 			},
// 			&errdetails.ResourceInfo{
// 				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName,
// 				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
// 			})
// 		if span != nil {
// 			span.RecordError(err)
// 			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
// 			span.SetStatus(trace.StatusError, err.Error())
// 		}
// 		return nil, err
// 	}
// 	return location.ReStructure(ids, locs), nil
// }
//
func (s *server) GetObject(ctx context.Context, req *payload.Object_VectorRequest) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec, err = s.client.GetObject(ctx, req)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.GetObjectRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return vec, nil
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamGetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_VectorRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Object_VectorRequest)
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamGetObjectRPCName+"/id-"+req.GetId().GetId())
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
		})

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
