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
	"strconv"
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

type server struct {
	eg                errgroup.Group
	gateway           service.Gateway // Mirror Gateway client for the other cluster.
	client            vclient.Client  // LB Gateway client for the same cluster.
	timeout           time.Duration
	replica           int
	streamConcurrency int
	name              string
	ip                string
	mirror.UnimplementedValdServerWithMirror
}

const (
	apiName      = "vald/gateway/mirror"
	rollbackName = "Rollback"
)

func New(opts ...Option) mirror.Server {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Register(ctx context.Context, req *payload.Mirror_Targets) (*payload.Mirror_Targets, error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, mirror.PackageName+"."+mirror.MirrorRPCServiceName+"/"+mirror.RegisterRPCName), apiName+"/"+mirror.RegisterRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	// TODO:

	return nil, nil
}

func (s *server) Advertise(ctx context.Context, req *payload.Mirror_Targets) (*payload.Mirror_Targets, error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, mirror.PackageName+"."+mirror.MirrorRPCServiceName+"/"+mirror.AdvertiseRPCName), apiName+"/"+mirror.AdvertiseRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	// TODO:

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

func (s *server) handleSpan(rpcName string, span trace.Span, err error) error {
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

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (ce *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.InsertRPCName), apiName+"/"+vald.InsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(s.gateway.FromForwardedContext(ctx)) == 0 {
		err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.InsertRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err := vald.NewValdClient(conn).Insert(sctx, req, copts...)
			return s.handleSpan(vald.InsertRPCName, sspan, err)
		})
		if err != nil {
			if err := s.handleSpan(rollbackName+" for "+vald.InsertRPCName, span, s.insertRollback(ctx, req)); err != nil {
				return nil, err
			}
			return nil, s.handleSpan(vald.InsertRPCName, span, err)
		}
	}

	ce, err = s.client.Insert(ctx, req)
	if err := s.handleSpan(vald.InsertRPCName, span, err); err != nil {
		return nil, err
	}
	return ce, nil
}

// insertRollback executes the Remove RPC for rollback.
func (s *server) insertRollback(ctx context.Context, req *payload.Insert_Request) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	newReq := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: false,
		},
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.RemoveRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).Remove(sctx, newReq, copts...)
		return s.handleSpan(vald.RemoveRPCName, sspan, err)
	})
	if err := s.handleSpan(vald.RemoveRPCName, span, err); err != nil {
		return err
	}
	return nil
}

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

	if len(s.gateway.FromForwardedContext(ctx)) == 0 {
		err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiInsertRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err := vald.NewValdClient(conn).MultiInsert(sctx, reqs, copts...)
			return s.handleSpan(vald.MultiInsertRPCName, sspan, err)
		})
		if err != nil {
			if err := s.handleSpan(rollbackName+" for "+vald.MultiInsertRPCName, span, s.multiInsertRollback(ctx, reqs)); err != nil {
				return nil, err
			}
			return nil, s.handleSpan(vald.MultiInsertRPCName, span, err)
		}
	}

	locs, err = s.client.MultiInsert(ctx, reqs, s.client.GRPCClient().GetCallOption()...)
	if err := s.handleSpan(vald.MultiInsertRPCName, span, err); err != nil {
		return nil, err
	}
	return locs, nil
}

// multiInsertRollback executes the MultiRemove RPC for rollback.
func (s *server) multiInsertRollback(ctx context.Context, reqs *payload.Insert_MultiRequest) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	newReq := &payload.Remove_MultiRequest{
		Requests: make([]*payload.Remove_Request, 0, len(reqs.Requests)),
	}
	for _, req := range reqs.Requests {
		newReq.Requests = append(newReq.Requests, &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: req.GetVector().GetId(),
			},
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: false,
			},
		})
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiRemoveRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).MultiRemove(sctx, newReq, copts...)
		return s.handleSpan(vald.MultiRemoveRPCName, sspan, err)
	})
	if err := s.handleSpan(vald.MultiRemoveRPCName, span, err); err != nil {
		return err
	}
	return nil
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateRPCName), apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(s.gateway.FromForwardedContext(ctx)) == 0 {
		err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.UpdateRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err := vald.NewValdClient(conn).Update(sctx, req, copts...)
			return s.handleSpan(vald.UpdateRPCName, sspan, err)
		})
		if err != nil {
			if err := s.handleSpan(rollbackName+" for "+vald.UpdateRPCName, span, s.updateRollback(ctx, req)); err != nil {
				return nil, err
			}
			return nil, s.handleSpan(vald.UpdateRPCName, span, err)
		}
	}

	ce, err := s.client.Update(ctx, req, s.client.GRPCClient().GetCallOption()...)
	if err = s.handleSpan(vald.UpdateRPCName, span, err); err != nil {
		return nil, err
	}
	return ce, nil
}

// updateRollback executes the GetObject on the same cluster to get old vector data and executes the Update RPC for rollback.
func (s *server) updateRollback(ctx context.Context, req *payload.Update_Request) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateRPCName), apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	loc, err := s.GetObject(ctx, &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
	})
	if err = s.handleSpan(vald.GetObjectRPCName+" for "+vald.UpdateRPCName, span, err); err != nil {
		return err
	}

	newReq := &payload.Update_Request{
		Vector: &payload.Object_Vector{
			Id:     loc.GetId(),
			Vector: loc.GetVector(),
		},
		Config: &payload.Update_Config{
			SkipStrictExistCheck: false,
		},
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.UpdateRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).Update(sctx, newReq, copts...)
		return s.handleSpan(vald.UpdateRPCName, sspan, err)
	})
	if err := s.handleSpan(vald.UpdateRPCName, span, err); err != nil {
		return err
	}
	return nil
}

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
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.MultiUpdateRPCName), apiName+"/"+vald.MultiUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(s.gateway.FromForwardedContext(ctx)) == 0 {
		err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiUpdateRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err := vald.NewValdClient(conn).MultiUpdate(sctx, reqs)
			return s.handleSpan(vald.MultiUpdateRPCName, sspan, err)
		})
		if err != nil {
			if err := s.handleSpan(rollbackName+" for "+vald.MultiUpdateRPCName, span, s.multiUpdateRollback(ctx, reqs)); err != nil {
				return nil, err
			}
			return nil, s.handleSpan(vald.MultiUpdateRPCName, span, err)
		}
	}

	ces, err := s.client.MultiUpdate(ctx, reqs, s.client.GRPCClient().GetCallOption()...)
	if err = s.handleSpan(vald.MultiUpdateRPCName, span, err); err != nil {
		return nil, err
	}
	return ces, nil
}

// multiUpdateRollback executes the GetObject on the same cluster to get old vector data and executes the MultiUpdate RPC for rollback.
func (s *server) multiUpdateRollback(ctx context.Context, reqs *payload.Update_MultiRequest) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.MultiUpdateRPCName), apiName+"/"+vald.MultiUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	mu := new(sync.Mutex)
	newReqs := &payload.Update_MultiRequest{
		Requests: make([]*payload.Update_Request, 0, len(reqs.GetRequests())),
	}
	eg, egctx := errgroup.New(ctx)

	for idx, req := range reqs.GetRequests() {
		idx, req := idx, req
		eg.Go(func() error {
			ctx, sspan := trace.StartSpan(egctx, apiName+"."+vald.GetObjectRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			ovec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
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
					SkipStrictExistCheck: false,
				},
			})
			mu.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	err := s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiUpdateRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).MultiUpdate(sctx, newReqs, copts...)
		return s.handleSpan(vald.MultiUpdateRPCName, sspan, err)
	})
	if err := s.handleSpan(vald.MultiUpdateRPCName, span, err); err != nil {
		return err
	}
	return nil
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.UpsertRPCName), apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(s.gateway.FromForwardedContext(ctx)) == 0 {
		err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.UpsertRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err := vald.NewValdClient(conn).Upsert(sctx, req, copts...)
			return s.handleSpan(vald.UpsertRPCName, sspan, err)
		})
		if err != nil {
			if err := s.handleSpan(rollbackName+" for "+vald.UpsertRPCName, span, s.upsertRollback(ctx, req)); err != nil {
				return nil, err
			}
			return nil, s.handleSpan(vald.UpsertRPCName, span, err)
		}
	}

	ce, err := s.client.Upsert(ctx, req, s.client.GRPCClient().GetCallOption()...)
	if err := s.handleSpan(vald.UpsertRPCName, span, err); err != nil {
		return nil, err
	}
	return ce, nil
}

// upsertRollback executes the updateRollback method. If NotFound error occurs, executes the insertRollback method.
func (s *server) upsertRollback(ctx context.Context, req *payload.Upsert_Request) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.UpsertRPCName), apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err = s.updateRollback(ctx, &payload.Update_Request{
		Vector: req.GetVector(),
		Config: &payload.Update_Config{
			SkipStrictExistCheck: false,
		},
	})
	if err != nil {
		st, _, err := status.ParseError(err, codes.Internal, "error "+vald.GetObjectRPCName+" API")
		if err != nil && st.Code() == codes.NotFound {
			if err := s.insertRollback(ctx, &payload.Insert_Request{
				Vector: req.GetVector(),
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: false,
				},
			}); err != nil {
				return s.handleSpan(rollbackName+" for "+vald.UpsertRPCName, span, err)
			}
		}
		return s.handleSpan(rollbackName+" for "+vald.UpsertRPCName, span, err)
	}
	return nil
}

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
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.MultiUpsertRPCName), apiName+"/"+vald.MultiUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(s.gateway.FromForwardedContext(ctx)) == 0 {
		err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.UpsertRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err := vald.NewValdClient(conn).MultiUpsert(sctx, reqs, copts...)
			return s.handleSpan(vald.UpsertRPCName, sspan, err)
		})
		if err != nil {
			if err := s.handleSpan(rollbackName+" for "+vald.MultiUpsertRPCName, span, s.multiUpsertRollback(ctx, reqs)); err != nil {
				return nil, err
			}
			return nil, s.handleSpan(vald.MultiUpsertRPCName, span, err)
		}
	}

	res, err = s.client.MultiUpsert(ctx, reqs, s.client.GRPCClient().GetCallOption()...)
	if err := s.handleSpan(vald.MultiUpsertRPCName, span, err); err != nil {
		return nil, err
	}
	return res, nil
}

// multiUpsertRollback executes the upsertRollback method.
func (s *server) multiUpsertRollback(ctx context.Context, reqs *payload.Upsert_MultiRequest) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.MultiUpsertRPCName), apiName+"/"+vald.MultiUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	eg, egctx := errgroup.New(ctx)
	for idx, req := range reqs.Requests {
		idx, req := idx, req
		eg.Go(func() error {
			ctx, sspan := trace.StartSpan(egctx, apiName+"."+vald.UpsertRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			if err := s.handleSpan(rollbackName+" for "+vald.UpsertRPCName, sspan, s.upsertRollback(ctx, req)); err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(s.gateway.FromForwardedContext(ctx)) == 0 {
		err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.RemoveRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err := vald.NewValdClient(conn).Remove(sctx, req, copts...)
			return s.handleSpan(vald.RemoveRPCName, sspan, err)
		})
		if err != nil {
			if err := s.handleSpan(rollbackName+" for "+vald.RemoveRPCName, span, s.removeRollback(ctx, req)); err != nil {
				return nil, err
			}
			return nil, s.handleSpan(vald.RemoveRPCName, span, err)
		}
	}

	loc, err = s.client.Remove(ctx, req, s.client.GRPCClient().GetCallOption()...)
	if err := s.handleSpan(vald.RemoveRPCName, span, err); err != nil {
		return nil, err
	}
	return loc, nil
}

// updateRollback executes the GetObject on the same cluster to get old vector data and executes the Upsert RPC for rollback.
func (s *server) removeRollback(ctx context.Context, req *payload.Remove_Request) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	loc, err := s.GetObject(ctx, &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetId().GetId(),
		},
	})
	if err = s.handleSpan(vald.GetObjectRPCName+" for "+vald.RemoveRPCName, span, err); err != nil {
		return err
	}

	newReq := &payload.Upsert_Request{
		Vector: &payload.Object_Vector{
			Id:     loc.GetId(),
			Vector: loc.GetVector(),
		},
		Config: &payload.Upsert_Config{
			SkipStrictExistCheck: false,
		},
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.UpsertRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := vald.NewValdClient(conn).Upsert(sctx, newReq, copts...)
		return s.handleSpan(vald.UpsertRPCName, sspan, err)
	})
	if err = s.handleSpan(vald.UpsertObjectRPCName+" for "+vald.RemoveRPCName, span, err); err != nil {
		return err
	}
	return nil
}

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
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(s.gateway.FromForwardedContext(ctx)) == 0 {
		err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.RemoveRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err := vald.NewValdClient(conn).MultiRemove(sctx, reqs, copts...)
			return s.handleSpan(vald.RemoveRPCName, sspan, err)
		})
		if err != nil {
			if err := s.handleSpan(rollbackName+" for "+vald.RemoveRPCName, span, s.multiRemoveRollback(ctx, reqs)); err != nil {
				return nil, err
			}
			return nil, s.handleSpan(vald.MultiRemoveRPCName, span, err)
		}
	}

	locs, err = s.client.MultiRemove(ctx, reqs, s.client.GRPCClient().GetCallOption()...)
	if err := s.handleSpan(vald.MultiRemoveRPCName, span, err); err != nil {
		return nil, err
	}
	return locs, nil
}

// updateRollback executes the removeRollback method.
func (s *server) multiRemoveRollback(ctx context.Context, reqs *payload.Remove_MultiRequest) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	eg, egctx := errgroup.New(ctx)
	for idx, req := range reqs.Requests {
		idx, req := idx, req
		eg.Go(func() error {
			ctx, sspan := trace.StartSpan(egctx, apiName+"."+vald.RemoveRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			if err := s.handleSpan(rollbackName+" for "+vald.RemoveRPCName, sspan, s.removeRollback(ctx, req)); err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

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
