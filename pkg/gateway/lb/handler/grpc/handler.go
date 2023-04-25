//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/slices"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type server struct {
	eg                errgroup.Group
	gateway           service.Gateway
	timeout           time.Duration
	replica           int
	streamConcurrency int
	multiConcurrency  int
	name              string
	ip                string
	vald.UnimplementedValdServer
}

const apiName = "vald/gateway/lb"

func New(opts ...Option) vald.Server {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) exists(ctx context.Context, uuid string) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "exists"), apiName+"/"+vald.ExistsRPCName+"/exists")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(uuid) == 0 {
		err = errors.ErrInvalidUUID(uuid)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	ich := make(chan *payload.Object_ID, 1)
	ech := make(chan error, 1)
	doneErr := errors.New("done exists")
	ctx, cancel := context.WithCancelCause(ctx)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(ich)
		defer close(ech)
		var once sync.Once
		ech <- s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/exists/BroadCast/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			meta := &payload.Object_ID{
				Id: uuid,
			}
			oid, err := vc.Exists(sctx, meta, copts...)
			if err != nil {
				var (
					attrs trace.Attributes
					st    *status.Status
					msg   string
					code  codes.Code
				)
				switch {
				case errors.Is(err, context.Canceled),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					attrs = trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.ExistsRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.ExistsRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.ExistsRPCName+" API meta "+uuid+"'s uuid not found",
						&errdetails.RequestInfo{
							RequestId:   uuid,
							ServingData: errdetails.Serialize(meta),
						},
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.ExistsRPCName,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if st != nil {
						code = st.Code()
					} else {
						code = codes.NotFound
					}
					attrs = trace.FromGRPCStatus(code, msg)
				}
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(attrs...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				if err != nil && st != nil &&
					code != codes.Canceled &&
					code != codes.DeadlineExceeded &&
					code != codes.InvalidArgument &&
					code != codes.NotFound &&
					code != codes.OK &&
					code != codes.Unimplemented {
					return err
				}
				return nil
			}
			if oid != nil && oid.GetId() != "" {
				once.Do(func() {
					ich <- oid
					cancel(doneErr)
				})
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
		if errors.Is(err, context.Canceled) && errors.Is(context.Cause(ctx), doneErr) {
			select {
			case id = <-ich:
				if id == nil || id.GetId() == "" {
					err = errors.ErrObjectIDNotFound(uuid)
				} else {
					err = nil
				}
			default:
			}
		}
	case id = <-ich:
		if id == nil || id.GetId() == "" {
			err = errors.ErrObjectIDNotFound(uuid)
		}
	case err = <-ech:
	}
	if err == nil && (id == nil || id.GetId() == "") {
		err = errors.ErrObjectIDNotFound(uuid)
	}
	if err != nil {
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return id, nil
}

func (s *server) Exists(ctx context.Context, meta *payload.Object_ID) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.ExistsRPCName), apiName+"/"+vald.ExistsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := meta.GetId()
	id, err = s.exists(ctx, uuid)
	if err == nil {
		return id, nil
	}
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(meta),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.ExistsRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	var attrs trace.Attributes
	switch {
	case errors.Is(err, errors.ErrInvalidUUID(uuid)):
		err = status.WrapWithInvalidArgument(vald.ExistsRPCName+" API invalid argument for uuid \""+uuid+"\" detected", err, reqInfo, resInfo, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       "uuid",
					Description: err.Error(),
				},
			},
		})
		attrs = trace.StatusCodeInvalidArgument(err.Error())
	case errors.Is(err, errors.ErrObjectIDNotFound(uuid)):
		err = status.WrapWithNotFound(vald.ExistsRPCName+" API id "+uuid+"'s data not found", err, reqInfo, resInfo)
		attrs = trace.StatusCodeNotFound(err.Error())
	case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
		err = status.WrapWithInternal(vald.ExistsRPCName+" API connection not found", err, reqInfo, resInfo)
		attrs = trace.StatusCodeInternal(err.Error())
	case errors.Is(err, context.Canceled):
		err = status.WrapWithCanceled(vald.ExistsRPCName+" API canceled", err, reqInfo, resInfo)
		attrs = trace.StatusCodeCancelled(err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		err = status.WrapWithDeadlineExceeded(vald.ExistsRPCName+" API deadline exceeded", err, reqInfo, resInfo)
		attrs = trace.StatusCodeDeadlineExceeded(err.Error())
	default:
		var (
			st  *status.Status
			msg string
		)
		st, msg, err = status.ParseError(err, codes.Unknown, vald.ExistsRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
		attrs = trace.FromGRPCStatus(st.Code(), msg)
	}
	log.Debug(err)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(attrs...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.SearchRPCName), apiName+"/"+vald.SearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vl := len(req.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument(vald.SearchRPCName+" API invalid vector argument", err,
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
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	cfg := req.GetConfig()
	mn := cfg.GetMinNum()
	if req.Config != nil {
		req.Config.MinNum = 0
	}
	res, err = s.doSearch(ctx, &payload.Search_Config{
		RequestId:            cfg.GetRequestId(),
		Num:                  cfg.GetNum(),
		MinNum:               mn,
		Radius:               cfg.GetRadius(),
		Epsilon:              cfg.GetEpsilon(),
		Timeout:              cfg.GetTimeout(),
		IngressFilters:       cfg.GetIngressFilters(),
		EgressFilters:        cfg.GetEgressFilters(),
		AggregationAlgorithm: cfg.GetAggregationAlgorithm(),
	}, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
		return vc.Search(ctx, req, copts...)
	})
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.SearchRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId: req.GetConfig().GetRequestId(),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
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
	uuid := req.GetId()
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(req),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.SearchByIDRPCName+" API invalid uuid", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec, err := s.getObject(ctx, uuid)
	cfg := req.GetConfig()
	mn := cfg.GetMinNum()
	if req.Config != nil {
		req.Config.MinNum = 0
	}
	scfg := &payload.Search_Config{
		RequestId:            cfg.GetRequestId(),
		Num:                  cfg.GetNum(),
		MinNum:               mn,
		Radius:               cfg.GetRadius(),
		Epsilon:              cfg.GetEpsilon(),
		Timeout:              cfg.GetTimeout(),
		IngressFilters:       cfg.GetIngressFilters(),
		EgressFilters:        cfg.GetEgressFilters(),
		AggregationAlgorithm: cfg.GetAggregationAlgorithm(),
	}
	if err != nil {
		var (
			attrs trace.Attributes
			st    *status.Status
			msg   string
		)
		switch {
		case errors.Is(err, errors.ErrInvalidUUID(uuid)):
			err = status.WrapWithInvalidArgument(
				vald.GetObjectRPCName+" API for "+vald.SearchByIDRPCName+" API invalid argument for uuid \""+uuid+"\" detected",
				err,
				reqInfo,
				resInfo,
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "uuid",
							Description: err.Error(),
						},
					},
				},
			)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.GetObjectRPCName+" API for "+vald.SearchByIDRPCName+" API connection not found", err, reqInfo, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.GetObjectRPCName+" API for "+vald.SearchByIDRPCName+" API canceled", err, reqInfo, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.GetObjectRPCName+" API for "+vald.SearchByIDRPCName+" API deadline exceeded", err, reqInfo, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		case errors.Is(err, errors.ErrObjectIDNotFound(uuid)), errors.Is(err, errors.ErrObjectNotFound(nil, uuid)):
			err = nil
		default:
			st, msg, err = status.ParseError(err, codes.Unknown, vald.GetObjectRPCName+" API for "+vald.SearchByIDRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
			if st == nil || st.Code() == codes.NotFound {
				err = nil
			}
		}
		if err != nil {
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		res, err = s.doSearch(ctx, scfg, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.SearchByID(ctx, req, copts...)
		})
		if err == nil {
			return res, nil
		}
		st, msg, err = status.ParseError(err, codes.Internal, vald.SearchByIDRPCName+" API failed to process search request", reqInfo, resInfo)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	res, err = s.Search(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: scfg,
	})
	if err != nil {
		res, err = s.doSearch(ctx, scfg, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.SearchByID(ctx, req, copts...)
		})
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal, vald.SearchByIDRPCName+" API failed to process search request", reqInfo, resInfo)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
	}
	return res, nil
}

func (s *server) doSearch(ctx context.Context, cfg *payload.Search_Config,
	f func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error)) (
	res *payload.Search_Response, err error,
) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doSearch"), apiName+"/doSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var (
		aggr    Aggregator
		num     = int(cfg.GetNum())
		replica = s.gateway.GetAgentCount(ctx)
	)
	switch cfg.GetAggregationAlgorithm() {
	case payload.Search_Unknown:
		aggr = newStd(num, replica)
	case payload.Search_ConcurrentQueue:
		aggr = newStd(num, replica)
	case payload.Search_SortSlice:
		aggr = newSlice(num, replica)
	case payload.Search_SortPoolSlice:
		aggr = newPoolSlice(num, replica)
	case payload.Search_PairingHeap:
		aggr = newPairingHeap(num, replica)
	default:
		aggr = newStd(num, replica)
	}
	return s.aggregationSearch(ctx, aggr, cfg, f)
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
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamSearchByIDRPCName), apiName+"/"+vald.StreamSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamSearchByIDRPCName+"/id-"+req.GetId())
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
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiSearchRPCName), apiName+"/"+vald.MultiSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu, emu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.GetRequests() {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + strconv.Itoa(idx)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiSearchRPCName+"/"+ti)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.Search(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				emu.Lock()
				if errs != nil {
					errs = errors.Join(errs, err)
				} else {
					errs = err
				}
				emu.Unlock()
				return nil
			}
			mu.Lock()
			res.Responses[idx] = r
			mu.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiSearchRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiSearchRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
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

func (s *server) MultiSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiSearchByIDRPCName), apiName+"/"+vald.MultiSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu, emu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.GetRequests() {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + strconv.Itoa(idx)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiSearchByIDRPCName+"/"+ti)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.SearchByID(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchByIDRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				emu.Lock()
				if errs != nil {
					errs = errors.Join(errs, err)
				} else {
					errs = err
				}
				emu.Unlock()
				return nil
			}
			mu.Lock()
			res.Responses[idx] = r
			mu.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiSearchByIDRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiSearchByIDRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
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

func (s *server) LinearSearch(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.LinearSearchRPCName), apiName+"/"+vald.LinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
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
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	cfg := req.GetConfig()
	mn := cfg.GetMinNum()
	if req.Config != nil {
		req.Config.MinNum = 0
	}
	res, err = s.doSearch(ctx, &payload.Search_Config{
		RequestId:            cfg.GetRequestId(),
		Num:                  cfg.GetNum(),
		MinNum:               mn,
		Timeout:              cfg.GetTimeout(),
		IngressFilters:       cfg.GetIngressFilters(),
		EgressFilters:        cfg.GetEgressFilters(),
		AggregationAlgorithm: cfg.GetAggregationAlgorithm(),
	}, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
		return vc.LinearSearch(ctx, req, copts...)
	})
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.LinearSearchRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId: req.GetConfig().GetRequestId(),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
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

	uuid := req.GetId()
	reqInfo := &errdetails.RequestInfo{
		RequestId: uuid,
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	if len(req.GetId()) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.LinearSearchByIDRPCName+" API invalid uuid", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec, err := s.getObject(ctx, uuid)
	cfg := req.GetConfig()
	mn := cfg.GetMinNum()
	if req.Config != nil {
		req.Config.MinNum = 0
	}
	scfg := &payload.Search_Config{
		RequestId:            cfg.GetRequestId(),
		Num:                  cfg.GetNum(),
		MinNum:               mn,
		Timeout:              cfg.GetTimeout(),
		IngressFilters:       cfg.GetIngressFilters(),
		EgressFilters:        cfg.GetEgressFilters(),
		AggregationAlgorithm: cfg.GetAggregationAlgorithm(),
	}
	if err != nil {
		var (
			attrs trace.Attributes
			st    *status.Status
			msg   string
		)
		switch {
		case errors.Is(err, errors.ErrInvalidUUID(uuid)):
			err = status.WrapWithInvalidArgument(
				vald.GetObjectRPCName+" API for "+vald.SearchByIDRPCName+" API invalid argument for uuid \""+uuid+"\" detected",
				err,
				reqInfo,
				resInfo,
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "uuid",
							Description: err.Error(),
						},
					},
				},
			)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.GetObjectRPCName+" API for "+vald.SearchByIDRPCName+" API connection not found", err, reqInfo, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.GetObjectRPCName+" API for "+vald.SearchByIDRPCName+" API canceled", err, reqInfo, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.GetObjectRPCName+" API for "+vald.SearchByIDRPCName+" API deadline exceeded", err, reqInfo, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		case errors.Is(err, errors.ErrObjectIDNotFound(uuid)), errors.Is(err, errors.ErrObjectNotFound(nil, uuid)):
			err = nil
		default:
			st, msg, err = status.ParseError(err, codes.Unknown, vald.GetObjectRPCName+" API for "+vald.SearchByIDRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
			if st == nil || st.Code() == codes.NotFound {
				err = nil
			}
		}
		if err != nil {
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		res, err = s.doSearch(ctx, scfg, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.LinearSearchByID(ctx, req, copts...)
		})
		if err == nil {
			return res, nil
		}
		st, msg, err = status.ParseError(err, codes.Internal, vald.LinearSearchByIDRPCName+" API failed to process search request", reqInfo, resInfo)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	res, err = s.LinearSearch(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: scfg,
	})
	if err != nil {
		res, err = s.doSearch(ctx, scfg, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.LinearSearchByID(ctx, req, copts...)
		})
		if err == nil {
			return res, nil
		}
		st, msg, err := status.ParseError(err, codes.Internal, vald.LinearSearchByIDRPCName+" API failed to process search request", reqInfo, resInfo)
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
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamLinearSearchByIDRPCName+"/id-"+req.GetId())
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

func (s *server) MultiLinearSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiLinearSearchRPCName), apiName+"/"+vald.MultiLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu, emu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.GetRequests() {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + strconv.Itoa(idx)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiLinearSearchRPCName+"/"+ti)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.LinearSearch(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				emu.Lock()
				if errs != nil {
					errs = errors.Join(errs, err)
				} else {
					errs = err
				}
				emu.Unlock()
				return nil
			}
			mu.Lock()
			res.Responses[idx] = r
			mu.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiLinearSearchRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiLinearSearchRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
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

func (s *server) MultiLinearSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiLinearSearchByIDRPCName), apiName+"/"+vald.MultiLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu, emu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.GetRequests() {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + strconv.Itoa(idx)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiLinearSearchByIDRPCName+"/"+ti)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.LinearSearchByID(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchByIDRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				emu.Lock()
				if errs != nil {
					errs = errors.Join(errs, err)
				} else {
					errs = err
				}
				emu.Unlock()
				return nil
			}
			mu.Lock()
			res.Responses[idx] = r
			mu.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiLinearSearchByIDRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiLinearSearchByIDRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
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

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (ce *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.InsertRPCName), apiName+"/"+vald.InsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetVector().GetId()
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(req),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.InsertRPCName+" API invalid uuid", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec := req.GetVector().GetVector()
	vl := len(vec)
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument(vald.InsertRPCName+" API invalid vector argument", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, err := s.exists(ctx, uuid)
		var attrs trace.Attributes
		if err != nil {
			switch {
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(vald.ExistsRPCName+" API for "+vald.InsertRPCName+" API connection not found", err, reqInfo, resInfo)
				attrs = trace.StatusCodeInternal(err.Error())
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(vald.ExistsRPCName+" API for "+vald.InsertRPCName+" API canceled", err, reqInfo, resInfo)
				attrs = trace.StatusCodeCancelled(err.Error())
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(vald.ExistsRPCName+" API for "+vald.InsertRPCName+" API deadline exceeded", err, reqInfo, resInfo)
				attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			default:
				err = nil
			}
		} else if id != nil && len(id.GetId()) != 0 {
			err = status.WrapWithAlreadyExists(vald.InsertRPCName+" API uuid "+uuid+"'s data already exists", errors.ErrMetaDataAlreadyExists(uuid), reqInfo, resInfo, info.Get())
			attrs = trace.StatusCodeAlreadyExists(err.Error())
		}
		if err != nil {
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Insert_Config{SkipStrictExistCheck: true}
		}
	}

	mu := new(sync.Mutex)
	ce = &payload.Object_Location{
		Uuid: uuid,
		Ips:  make([]string, 0, s.replica),
	}
	locs := make([]string, 0, s.replica)
	if req.GetConfig().GetTimestamp() == 0 {
		now := time.Now().UnixNano()
		if req.GetConfig() == nil {
			req.Config = &payload.Insert_Config{
				Timestamp: now,
			}
		} else {
			req.GetConfig().Timestamp = now
		}
	}
	emu := new(sync.Mutex)
	var errs error
	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "DoMulti/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.Insert(ctx, req, copts...)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled),
				errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.InsertRPCName + ".DoMulti/" +
							target + " canceled: " + err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil
			case errors.Is(err, context.DeadlineExceeded),
				errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.InsertRPCName + ".DoMulti/" +
							target + " deadline_exceeded: " + err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil
			}
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.InsertRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				})
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if err != nil && st.Code() != codes.AlreadyExists {
				emu.Lock()
				if errs != nil {
					errs = errors.Join(errs, err)
				} else {
					errs = err
				}
				emu.Unlock()
				return err
			}
			return nil
		}
		mu.Lock()
		ce.Ips = append(ce.GetIps(), loc.GetIps()...)
		locs = append(locs, loc.GetName())
		mu.Unlock()
		return nil
	})
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(vald.InsertRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + ".DoMulti",
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				})
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if errs == nil {
			errs = err
		} else {
			errs = errors.Join(errs, err)
		}
	}
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse "+vald.InsertRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + ".DoMulti",
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	slices.Sort(locs)
	ce.Name = strings.Join(locs, ",")
	return ce, nil
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

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.MultiInsertRPCName), apiName+"/"+vald.MultiInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var (
		emu sync.Mutex
		lmu sync.Mutex
	)
	eg, ectx := errgroup.New(ctx)
	eg.Limitation(s.multiConcurrency)
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	for i, r := range reqs.GetRequests() {
		if r != nil && r.GetVector() != nil && len(r.GetVector().GetVector()) >= algorithm.MinimumVectorDimensionSize && r.GetVector().GetId() != "" {
			idx := i
			req := r
			eg.Go(safety.RecoverFunc(func() (err error) {
				ectx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ectx, "eg.Go"), apiName+"/"+vald.MultiInsertRPCName+"/id-"+req.GetVector().GetId())
				defer func() {
					if sspan != nil {
						sspan.End()
					}
				}()
				res, err := s.Insert(ectx, req)
				if err != nil {
					st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.InsertRPCName+" gRPC error response")
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
					emu.Lock()
					if errs == nil {
						errs = err
					} else {
						errs = errors.Join(errs, err)
					}
					emu.Unlock()
				} else if res != nil && res.GetUuid() == req.GetVector().GetId() && res.GetIps() != nil {
					lmu.Lock()
					locs.Locations[idx] = res
					lmu.Unlock()
				}
				return nil
			}))
		} else {
			var (
				err   error
				field string
			)
			switch {
			case r.GetVector() == nil, len(r.GetVector().GetVector()) < algorithm.MinimumVectorDimensionSize:
				err = errors.ErrInvalidDimensionSize(len(r.GetVector().GetVector()), 0)
				field = "vector"
			case r.GetVector().GetId() == "":
				err = errors.ErrInvalidUUID(r.GetVector().GetId())
				field = "uuid"
			}
			err = status.WrapWithInvalidArgument(vald.MultiInsertRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   r.GetVector().GetId(),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiInsertRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       field,
							Description: err.Error(),
						},
					},
				}, info.Get())
			emu.Lock()
			if errs == nil {
				errs = err
			} else {
				errs = errors.Join(errs, err)
			}
			emu.Unlock()

		}
	}
	err := eg.Wait()
	if err != nil {
		emu.Lock()
		if errs == nil {
			errs = err
		} else {
			errs = errors.Join(errs, err)
		}
		emu.Unlock()
	}

	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "error detected"+vald.MultiInsertRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		errs = err
	}

	return locs, errs
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateRPCName), apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetVector().GetId()
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(req),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.GetObjectRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.UpdateRPCName+" API invalid uuid", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec := req.GetVector().GetVector()
	vl := len(vec)
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument(vald.UpdateRPCName+" API invalid vector argument", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	if req.GetConfig().GetDisableBalancedUpdate() {
		var (
			mu      sync.RWMutex
			aeCount atomic.Uint64
			updated atomic.Uint64
			ls      = make([]string, 0, s.replica)
			visited = make(map[string]bool, s.replica)
			locs    = &payload.Object_Location{
				Uuid: uuid,
				Ips:  make([]string, 0, s.replica),
			}
		)
		err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpdateRPCName+"/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			loc, err := vc.Update(ctx, req, copts...)
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st != nil {
					if st.Code() != codes.AlreadyExists &&
						st.Code() != codes.Canceled &&
						st.Code() != codes.DeadlineExceeded &&
						st.Code() != codes.InvalidArgument &&
						st.Code() != codes.NotFound &&
						st.Code() != codes.OK &&
						st.Code() != codes.Unimplemented {
						if span != nil {
							span.RecordError(err)
							span.SetAttributes(trace.FromGRPCStatus(st.Code(), fmt.Sprintf("Update operation for Agent %s failed,\terror: %v", target, err))...)
							span.SetStatus(trace.StatusError, err.Error())
						}
						return err
					}
					if st.Code() == codes.AlreadyExists {
						host, _, err := net.SplitHostPort(target)
						if err != nil {
							host = target
						}
						aeCount.Add(1)
						mu.Lock()
						visited[target] = true
						locs.Ips = append(locs.GetIps(), host)
						ls = append(ls, host)
						mu.Unlock()

					}
				}
				return nil
			}
			if loc != nil {
				updated.Add(1)
				mu.Lock()
				visited[target] = true
				locs.Ips = append(locs.GetIps(), loc.GetIps()...)
				ls = append(ls, loc.GetName())
				mu.Unlock()
			}
			return nil
		})
		switch {
		case err != nil:
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo, info.Get())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		case len(locs.Ips) <= 0:
			err = errors.ErrIndexNotFound
			err = status.WrapWithNotFound(vald.UpdateRPCName+" API update target not found", err, reqInfo, resInfo)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		case updated.Load()+aeCount.Load() < uint64(s.replica):
			shortage := s.replica - int(updated.Load()+aeCount.Load())
			err = s.gateway.DoMulti(ctx, shortage, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
				mu.RLock()
				tf, ok := visited[target]
				mu.RUnlock()
				if tf && ok {
					return errors.Errorf("target: %s already inserted will skip", target)
				}
				ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "DoMulti/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
				defer func() {
					if span != nil {
						span.End()
					}
				}()
				loc, err := vc.Insert(ctx, &payload.Insert_Request{
					Vector: req.GetVector(),
					Config: &payload.Insert_Config{
						SkipStrictExistCheck: true,
						Filters:              req.GetConfig().GetFilters(),
						Timestamp:            req.GetConfig().GetTimestamp(),
					},
				}, copts...)
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st != nil && span != nil {
						span.RecordError(err)
						span.SetAttributes(trace.FromGRPCStatus(st.Code(), fmt.Sprintf("Shortage index Insert for Update operation for Agent %s failed,\terror: %v", target, err))...)
						span.SetStatus(trace.StatusError, err.Error())
					}
					return err
				}
				if loc != nil {
					updated.Add(1)
					mu.Lock()
					locs.Ips = append(locs.GetIps(), loc.GetIps()...)
					ls = append(ls, loc.GetName())
					mu.Unlock()
				}
				return nil
			})
		case updated.Load() == 0 && aeCount.Load() > 0:
			err = errors.ErrSameVectorAlreadyExists(uuid, vec, vec)
			err = status.WrapWithAlreadyExists(vald.UpdateRPCName+" API update target same vector already exists", err, reqInfo, resInfo)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err

		}
		slices.Sort(ls)
		locs.Name = strings.Join(ls, ",")
		return locs, nil
	}

	if !req.GetConfig().GetSkipStrictExistCheck() {
		vec, err := s.getObject(ctx, uuid)
		if err != nil || vec == nil {
			var (
				attrs trace.Attributes
				st    *status.Status
				msg   string
			)
			switch {
			case errors.Is(err, errors.ErrInvalidUUID(uuid)):
				err = status.WrapWithInvalidArgument(
					vald.GetObjectRPCName+" API for "+vald.UpdateRPCName+" API invalid argument for uuid \""+uuid+"\" detected",
					err,
					reqInfo,
					resInfo,
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequestFieldViolation{
							{
								Field:       "uuid",
								Description: err.Error(),
							},
						},
					},
				)
				attrs = trace.StatusCodeInvalidArgument(err.Error())
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(vald.GetObjectRPCName+" API for "+vald.UpdateRPCName+" API connection not found", err, reqInfo, resInfo)
				attrs = trace.StatusCodeInternal(err.Error())
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(vald.GetObjectRPCName+" API for "+vald.UpdateRPCName+" API canceled", err, reqInfo, resInfo)
				attrs = trace.StatusCodeCancelled(err.Error())
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(vald.GetObjectRPCName+" API for "+vald.UpdateRPCName+" API deadline exceeded", err, reqInfo, resInfo)
				attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			case errors.Is(err, errors.ErrObjectIDNotFound(uuid)), errors.Is(err, errors.ErrObjectNotFound(nil, uuid)):
				err = status.WrapWithNotFound(vald.GetObjectRPCName+" API for "+vald.UpdateRPCName+" API uuid "+uuid+"'s object not found", err, reqInfo, resInfo)
				attrs = trace.StatusCodeNotFound(err.Error())
			default:
				code := codes.Unknown
				if err == nil {
					err = errors.ErrObjectIDNotFound(uuid)
					code = codes.NotFound
				}
				st, msg, err = status.ParseError(err, code, vald.GetObjectRPCName+" API for "+vald.UpdateRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
				attrs = trace.FromGRPCStatus(st.Code(), msg)
			}
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if conv.F32stos(vec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
			if err == nil {
				err = errors.ErrSameVectorAlreadyExists(uuid, vec.GetVector(), req.GetVector().GetVector())
			}
			st, msg, err := status.ParseError(err, codes.AlreadyExists,
				"error "+vald.UpdateRPCName+" API ID = "+uuid+"'s same vector data already exists",
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.GetObjectRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				}, info.Get())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Update_Config{SkipStrictExistCheck: true}
		}
	}
	var now int64
	if req.GetConfig().GetTimestamp() != 0 {
		now = req.GetConfig().GetTimestamp()
	} else {
		now = time.Now().UnixNano()
	}

	rreq := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: uuid,
		},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: true,
			Timestamp:            now,
		},
	}
	res, err = s.Remove(ctx, rreq)
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(vald.RemoveRPCName+" for "+vald.UpdateRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(rreq),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.RemoveRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				})
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveRPCName+" for "+vald.UpdateRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(rreq),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.RemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())

		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	now++
	ireq := &payload.Insert_Request{
		Vector: req.GetVector(),
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: true,
			Filters:              req.GetConfig().GetFilters(),
			Timestamp:            now,
		},
	}
	res, err = s.Insert(ctx, ireq)
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(ireq),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.InsertRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				})
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.InsertRPCName+" for "+vald.UpdateRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(ireq),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.InsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
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

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.MultiUpdateRPCName), apiName+"/"+vald.MultiUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var (
		emu sync.Mutex
		lmu sync.Mutex
	)
	eg, ectx := errgroup.New(ctx)
	eg.Limitation(s.multiConcurrency)
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	for i, r := range reqs.GetRequests() {
		if r != nil && r.GetVector() != nil && len(r.GetVector().GetVector()) >= algorithm.MinimumVectorDimensionSize && r.GetVector().GetId() != "" {
			idx := i
			req := r
			eg.Go(safety.RecoverFunc(func() (err error) {
				ectx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ectx, "eg.Go"), apiName+"/"+vald.MultiUpdateRPCName+"/id-"+req.GetVector().GetId())
				defer func() {
					if sspan != nil {
						sspan.End()
					}
				}()
				res, err := s.Update(ectx, req)
				if err != nil {
					st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpdateRPCName+" gRPC error response")
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
					emu.Lock()
					if errs == nil {
						errs = err
					} else {
						errs = errors.Join(errs, err)
					}
					emu.Unlock()
				} else if res != nil && res.GetUuid() == req.GetVector().GetId() && res.GetIps() != nil {
					lmu.Lock()
					locs.Locations[idx] = res
					lmu.Unlock()
				}
				return nil
			}))
		} else {
			var (
				err   error
				field string
			)
			switch {
			case r.GetVector() == nil, len(r.GetVector().GetVector()) < algorithm.MinimumVectorDimensionSize:
				err = errors.ErrInvalidDimensionSize(len(r.GetVector().GetVector()), 0)
				field = "vector"
			case r.GetVector().GetId() == "":
				err = errors.ErrInvalidUUID(r.GetVector().GetId())
				field = "uuid"
			}
			err = status.WrapWithInvalidArgument(vald.MultiUpdateRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   r.GetVector().GetId(),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpdateRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       field,
							Description: err.Error(),
						},
					},
				}, info.Get())
			emu.Lock()
			if errs == nil {
				errs = err
			} else {
				errs = errors.Join(errs, err)
			}
			emu.Unlock()

		}
	}
	err := eg.Wait()
	if err != nil {
		emu.Lock()
		if errs == nil {
			errs = err
		} else {
			errs = errors.Join(errs, err)
		}
		emu.Unlock()
	}

	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "error detected"+vald.MultiUpdateRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		errs = err
	}

	return locs, errs
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.UpsertRPCName), apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	vec := req.GetVector()
	uuid := vec.GetId()
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(req),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	if len(uuid) == 0 {
		err = status.WrapWithInvalidArgument(vald.UpsertRPCName+" API invalid uuid", errors.ErrInvalidMetaDataConfig, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	vl := len(vec.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument(vald.UpsertRPCName+" API invalid vector argument", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	var shouldInsert bool
	if !req.GetConfig().GetSkipStrictExistCheck() {
		vec, err := s.getObject(ctx, uuid)
		var (
			attrs trace.Attributes
			st    *status.Status
			msg   string
		)
		if err != nil || vec == nil {
			switch {
			case errors.Is(err, errors.ErrInvalidUUID(uuid)):
				err = status.WrapWithInvalidArgument(
					vald.GetObjectRPCName+" API for "+vald.UpsertRPCName+" API invalid argument for uuid \""+uuid+"\" detected",
					err,
					reqInfo,
					resInfo,
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequestFieldViolation{
							{
								Field:       "uuid",
								Description: err.Error(),
							},
						},
					},
				)
				attrs = trace.StatusCodeInvalidArgument(err.Error())
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(vald.GetObjectRPCName+" API for "+vald.UpsertRPCName+" API connection not found", err, reqInfo, resInfo)
				attrs = trace.StatusCodeInternal(err.Error())
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(vald.GetObjectRPCName+" API for "+vald.UpsertRPCName+" API canceled", err, reqInfo, resInfo)
				attrs = trace.StatusCodeCancelled(err.Error())
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(vald.GetObjectRPCName+" API for "+vald.UpsertRPCName+" API deadline exceeded", err, reqInfo, resInfo)
				attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			case errors.Is(err, errors.ErrObjectIDNotFound(uuid)), errors.Is(err, errors.ErrObjectNotFound(nil, uuid)):
				err = nil
				shouldInsert = true
			default:
				st, msg, err = status.ParseError(err, codes.Unknown, vald.GetObjectRPCName+" API for "+vald.UpsertRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
				attrs = trace.FromGRPCStatus(st.Code(), msg)
				if st != nil && st.Code() == codes.NotFound {
					err = nil
					shouldInsert = true
				}
			}
		} else if conv.F32stos(vec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
			err = status.WrapWithAlreadyExists(vald.GetObjectRPCName+" API for "+vald.UpsertRPCName+" API ID = "+uuid+"'s same vector data already exists", errors.ErrSameVectorAlreadyExists(uuid, vec.GetVector(), req.GetVector().GetVector()), reqInfo, resInfo)
			attrs = trace.StatusCodeAlreadyExists(err.Error())
		}
		if err != nil {
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

	} else {
		id, err := s.exists(ctx, uuid)
		if err != nil {
			var attrs trace.Attributes
			switch {
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(vald.ExistsRPCName+" API for "+vald.UpsertRPCName+" API connection not found", err, reqInfo, resInfo)
				attrs = trace.StatusCodeInternal(err.Error())
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(vald.ExistsRPCName+" API for "+vald.UpsertRPCName+" API canceled", err, reqInfo, resInfo)
				attrs = trace.StatusCodeCancelled(err.Error())
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(vald.ExistsRPCName+" API for "+vald.UpsertRPCName+" API deadline exceeded", err, reqInfo, resInfo)
				attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			default:
				err = nil
			}
			if err != nil {
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(attrs...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil, err
			}
		}
		shouldInsert = err != nil || id == nil || len(id.GetId()) == 0
	}

	var operation string
	if shouldInsert {
		operation = vald.InsertRPCName
		loc, err = s.Insert(ctx, &payload.Insert_Request{
			Vector: vec,
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Filters:              req.GetConfig().GetFilters(),
				Timestamp:            req.GetConfig().GetTimestamp(),
			},
		})
	} else {
		operation = vald.UpdateRPCName
		loc, err = s.Update(ctx, &payload.Update_Request{
			Vector: vec,
			Config: &payload.Update_Config{
				SkipStrictExistCheck:  true,
				Filters:               req.GetConfig().GetFilters(),
				Timestamp:             req.GetConfig().GetTimestamp(),
				DisableBalancedUpdate: req.GetConfig().GetDisableBalancedUpdate(),
			},
		})
	}

	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+operation+" for "+vald.UpsertRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + operation,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
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

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.MultiUpsertRPCName), apiName+"/"+vald.MultiUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var (
		emu sync.Mutex
		lmu sync.Mutex
	)
	eg, ectx := errgroup.New(ctx)
	eg.Limitation(s.multiConcurrency)
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	for i, r := range reqs.GetRequests() {
		if r != nil && r.GetVector() != nil && len(r.GetVector().GetVector()) >= algorithm.MinimumVectorDimensionSize && r.GetVector().GetId() != "" {
			idx := i
			req := r
			eg.Go(safety.RecoverFunc(func() (err error) {
				ectx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ectx, "eg.Go"), apiName+"/"+vald.MultiUpsertRPCName+"/id-"+req.GetVector().GetId())
				defer func() {
					if sspan != nil {
						sspan.End()
					}
				}()
				res, err := s.Upsert(ectx, req)
				if err != nil {
					st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpsertRPCName+" gRPC error response")
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
					emu.Lock()
					if errs == nil {
						errs = err
					} else {
						errs = errors.Join(errs, err)
					}
					emu.Unlock()
				} else if res != nil && res.GetUuid() == req.GetVector().GetId() && res.GetIps() != nil {
					lmu.Lock()
					locs.Locations[idx] = res
					lmu.Unlock()
				}
				return nil
			}))
		} else {
			var (
				err   error
				field string
			)
			switch {
			case r.GetVector() == nil, len(r.GetVector().GetVector()) < algorithm.MinimumVectorDimensionSize:
				err = errors.ErrInvalidDimensionSize(len(r.GetVector().GetVector()), 0)
				field = "vector"
			case r.GetVector().GetId() == "":
				err = errors.ErrInvalidUUID(r.GetVector().GetId())
				field = "uuid"
			}
			err = status.WrapWithInvalidArgument(vald.MultiUpsertRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   r.GetVector().GetId(),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpsertRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       field,
							Description: err.Error(),
						},
					},
				}, info.Get())
			emu.Lock()
			if errs == nil {
				errs = err
			} else {
				errs = errors.Join(errs, err)
			}
			emu.Unlock()

		}
	}
	err := eg.Wait()
	if err != nil {
		emu.Lock()
		if errs == nil {
			errs = err
		} else {
			errs = errors.Join(errs, err)
		}
		emu.Unlock()
	}

	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "error detected"+vald.MultiUpsertRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		errs = err
	}

	return locs, errs
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (locs *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	id := req.GetId()
	uuid := id.GetId()
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(req),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		_, err := s.exists(ctx, uuid)
		if err != nil {
			var attrs trace.Attributes
			switch {
			case errors.Is(err, errors.ErrInvalidUUID(uuid)):
				err = status.WrapWithInvalidArgument(
					vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API invalid argument for uuid \""+uuid+"\" detected",
					err,
					reqInfo,
					resInfo,
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequestFieldViolation{
							{
								Field:       "uuid",
								Description: err.Error(),
							},
						},
					},
				)
				attrs = trace.StatusCodeInvalidArgument(err.Error())
			case errors.Is(err, errors.ErrObjectIDNotFound(uuid)):
				err = status.WrapWithNotFound(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API id "+uuid+"'s data not found", err, reqInfo, resInfo)
				attrs = trace.StatusCodeNotFound(err.Error())
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API connection not found", err, reqInfo, resInfo)
				attrs = trace.StatusCodeInternal(err.Error())
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API canceled", err, reqInfo, resInfo)
				attrs = trace.StatusCodeCancelled(err.Error())
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API deadline exceeded", err, reqInfo, resInfo)
				attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			default:
				var (
					st  *status.Status
					msg string
				)
				st, msg, err = status.ParseError(err, codes.Unknown, vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
				attrs = trace.FromGRPCStatus(st.Code(), msg)
			}
			if err != nil {
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(attrs...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil, err
			}
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Remove_Config{SkipStrictExistCheck: true}
		}
	}
	if req.GetConfig().GetTimestamp() == 0 {
		now := time.Now().UnixNano()
		if req.GetConfig() == nil {
			req.Config = &payload.Remove_Config{
				Timestamp: now,
			}
		} else {
			req.GetConfig().Timestamp = now
		}
	}
	var mu sync.Mutex
	locs = &payload.Object_Location{
		Uuid: id.GetId(),
		Ips:  make([]string, 0, s.replica),
	}
	ls := make([]string, 0, s.replica)
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.RemoveRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.Remove(ctx, req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveRPCName+" gRPC error response", reqInfo, resInfo)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if err != nil && st.Code() != codes.NotFound {
				log.Error(err)
				return err
			}
			return nil
		}
		mu.Lock()
		locs.Ips = append(locs.GetIps(), loc.GetIps()...)
		ls = append(ls, loc.GetName())
		mu.Unlock()
		return nil
	})
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveRPCName+" gRPC error response", reqInfo, resInfo, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if len(locs.Ips) <= 0 {
		err = errors.ErrIndexNotFound
		err = status.WrapWithNotFound(vald.RemoveRPCName+" API remove target not found", err, reqInfo, resInfo)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	slices.Sort(ls)
	locs.Name = strings.Join(ls, ",")
	return locs, nil
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamSearchRPCName), apiName+"/"+vald.StreamSearchRPCName)
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

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var (
		emu sync.Mutex
		lmu sync.Mutex
	)
	eg, ectx := errgroup.New(ctx)
	eg.Limitation(s.multiConcurrency)
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	for i, r := range reqs.GetRequests() {
		if r != nil && r.GetId().GetId() != "" {
			idx := i
			req := r
			eg.Go(safety.RecoverFunc(func() (err error) {
				ectx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ectx, "eg.Go"), apiName+"/"+vald.MultiRemoveRPCName+"/id-"+req.GetId().GetId())
				defer func() {
					if sspan != nil {
						sspan.End()
					}
				}()
				res, err := s.Remove(ectx, req)
				if err != nil {
					st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.RemoveRPCName+" gRPC error response")
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
					emu.Lock()
					if errs == nil {
						errs = err
					} else {
						errs = errors.Join(errs, err)
					}
					emu.Unlock()
				} else if res != nil && res.GetUuid() == req.GetId().GetId() && res.GetIps() != nil {
					lmu.Lock()
					locs.Locations[idx] = res
					lmu.Unlock()
				}
				return nil
			}))
		} else {
			err := errors.ErrInvalidUUID(r.GetId().GetId())
			err = status.WrapWithInvalidArgument(vald.MultiRemoveRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   r.GetId().GetId(),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "uuid",
							Description: err.Error(),
						},
					},
				}, info.Get())
			emu.Lock()
			if errs == nil {
				errs = err
			} else {
				errs = errors.Join(errs, err)
			}
			emu.Unlock()

		}
	}
	err := eg.Wait()
	if err != nil {
		emu.Lock()
		if errs == nil {
			errs = err
		} else {
			errs = errors.Join(errs, err)
		}
		emu.Unlock()
	}

	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "error detected"+vald.MultiRemoveRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		errs = err
	}

	return locs, errs
}

func (s *server) getObject(ctx context.Context, uuid string) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "getObject"), apiName+"/"+vald.GetObjectRPCName+"/getObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vch := make(chan *payload.Object_Vector, 1)
	ech := make(chan error, 1)
	doneErr := errors.New("done getObject")
	ctx, cancel := context.WithCancelCause(ctx)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(vch)
		defer close(ech)
		var once sync.Once
		ech <- s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/getObject/BroadCast/"+target)
			defer func() {
				if span != nil {
					sspan.End()
				}
			}()
			req := &payload.Object_VectorRequest{
				Id: &payload.Object_ID{
					Id: uuid,
				},
			}
			ovec, err := vc.GetObject(sctx, req, copts...)
			if err != nil {
				var (
					attrs trace.Attributes
					st    *status.Status
					msg   string
					code  codes.Code
				)
				switch {
				case errors.Is(err, context.Canceled),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					attrs = trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.GetObjectRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.GetObjectRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.GetObjectRPCName+" API meta "+uuid+"'s uuid not found",
						&errdetails.RequestInfo{
							RequestId:   uuid,
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if st != nil {
						code = st.Code()
					} else {
						code = codes.NotFound
					}
					attrs = trace.FromGRPCStatus(code, msg)
				}
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(attrs...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				if err != nil && st != nil &&
					code != codes.Canceled &&
					code != codes.DeadlineExceeded &&
					code != codes.InvalidArgument &&
					code != codes.NotFound &&
					code != codes.OK &&
					code != codes.Unimplemented {
					return err
				}
				return nil
			}
			if ovec != nil && ovec.GetId() != "" && ovec.GetVector() != nil {
				once.Do(func() {
					vch <- ovec
					cancel(doneErr)
				})
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
		if errors.Is(err, context.Canceled) && errors.Is(context.Cause(ctx), doneErr) {
			select {
			case vec = <-vch:
				if vec == nil || vec.GetId() == "" || vec.GetVector() == nil {
					err = errors.ErrObjectNotFound(nil, uuid)
				} else {
					err = nil
				}
			default:
			}
		}
	case vec = <-vch:
		if vec == nil || vec.GetId() == "" || vec.GetVector() == nil {
			err = errors.ErrObjectNotFound(nil, uuid)
		}
	case err = <-ech:
	}
	if err == nil && (vec == nil || vec.GetId() == "" || vec.GetVector() == nil) {
		err = errors.ErrObjectNotFound(nil, uuid)
	}
	if err != nil {
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	return vec, nil
}

func (s *server) Flush(ctx context.Context, req *payload.Flush_Request) (cnts *payload.Info_Index_Count, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FlushRPCServiceName+"/"+vald.FlushRPCName), apiName+"/"+vald.FlushRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var (
		stored      uint32
		uncommitted uint32
		indexing    atomic.Value
		saving      atomic.Value
	)
	indexing.Store(false)
	saving.Store(false)
	now := time.Now().UnixNano()
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(ctx, apiName+"."+vald.FlushRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		cnt, err := vc.Flush(ctx, req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.FlushRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   strconv.FormatInt(now, 10),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.FlushRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				})
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if err != nil && st.Code() != codes.NotFound {
				log.Error(err)
				return err
			}
			return nil
		}
		atomic.AddUint32(&stored, cnt.Stored)
		atomic.AddUint32(&uncommitted, cnt.Uncommitted)
		if cnt.Indexing {
			indexing.Store(cnt.Indexing)
		}
		if cnt.Saving {
			saving.Store(cnt.Saving)
		}
		return nil
	})
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.FlushRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strconv.FormatInt(now, 10),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.FlushRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	cnts = &payload.Info_Index_Count{
		Stored:      atomic.LoadUint32(&stored),
		Uncommitted: atomic.LoadUint32(&uncommitted),
		Indexing:    indexing.Load().(bool),
		Saving:      saving.Load().(bool),
	}
	if cnts.Stored > 0 || cnts.Uncommitted > 0 || cnts.Indexing || cnts.Saving {
		err = errors.Errorf(
			"stored index: %d, uncommited: %d, indexing: %t, saving: %t",
			cnts.Stored, cnts.Uncommitted, cnts.Indexing, cnts.Saving,
		)
		err = status.WrapWithInternal(vald.FlushRPCName+" API flush failed", err,
			&errdetails.RequestInfo{
				RequestId:   strconv.FormatInt(now, 10),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.FlushRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	return cnts, nil
}

func (s *server) GetObject(ctx context.Context, req *payload.Object_VectorRequest) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetId().GetId()
	vec, err = s.getObject(ctx, uuid)
	if err == nil {
		if vec != nil && vec.GetId() != "" && vec.GetVector() != nil {
			return vec, nil
		}
		err = errors.ErrObjectNotFound(nil, uuid)
	}
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(req),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	var attrs trace.Attributes
	switch {
	case errors.Is(err, errors.ErrInvalidUUID(uuid)):
		err = status.WrapWithInvalidArgument(vald.GetObjectRPCName+" API invalid argument for uuid \""+uuid+"\" detected", err, reqInfo, resInfo, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       "uuid",
					Description: err.Error(),
				},
			},
		})
		attrs = trace.StatusCodeInvalidArgument(err.Error())
	case errors.Is(err, errors.ErrObjectIDNotFound(uuid)), errors.Is(err, errors.ErrObjectNotFound(nil, uuid)):
		err = status.WrapWithNotFound(vald.GetObjectRPCName+" API id "+uuid+"'s object not found", err, reqInfo, resInfo)
		attrs = trace.StatusCodeNotFound(err.Error())
	case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
		err = status.WrapWithInternal(vald.GetObjectRPCName+" API connection not found", err, reqInfo, resInfo)
		attrs = trace.StatusCodeInternal(err.Error())
	case errors.Is(err, context.Canceled):
		err = status.WrapWithCanceled(vald.GetObjectRPCName+" API canceled", err, reqInfo, resInfo)
		attrs = trace.StatusCodeCancelled(err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		err = status.WrapWithDeadlineExceeded(vald.GetObjectRPCName+" API deadline exceeded", err, reqInfo, resInfo)
		attrs = trace.StatusCodeDeadlineExceeded(err.Error())
	default:
		var (
			st  *status.Status
			msg string
		)
		st, msg, err = status.ParseError(err, codes.Unknown, vald.GetObjectRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
		attrs = trace.FromGRPCStatus(st.Code(), msg)
	}
	log.Debug(err)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(attrs...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
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
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamGetObjectRPCName+"/id-"+req.GetId().GetId())
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
