//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"io"
	"math"
	"slices"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/core/algorithm"
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
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
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
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
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

func (s *server) Exists(
	ctx context.Context, meta *payload.Object_ID,
) (id *payload.Object_ID, err error) {
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

func (s *server) Search(
	ctx context.Context, req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
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
	res, err = s.doSearch(ctx, req.GetConfig(), func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
		req.Config = fcfg
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

func (s *server) SearchByID(
	ctx context.Context, req *payload.Search_IDRequest,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.SearchByIDRPCName), apiName+"/"+vald.SearchByIDRPCName)
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
	vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: uuid,
		},
	})
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Unknown, vald.GetObjectRPCName+" API for "+vald.SearchByIDRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
		if span != nil && st != nil && st.Code() != codes.NotFound {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		// try search by using agent's SearchByID method this operation is emergency fallback, the search quality is not same as usual SearchByID operation.
		res, err = s.doSearch(ctx, req.GetConfig(), func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			req.Config = fcfg
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
		Config: req.GetConfig(),
	})
	if err != nil {
		res, err = s.doSearch(ctx, req.GetConfig(), func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			req.Config = fcfg
			return vc.SearchByID(ctx, req, copts...)
		})
		if err == nil {
			return res, nil
		}
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

// calculateNum adjusts the number of search results based on the ratio and the number of replicas.
// It ensures that the number of results is not less than the minimum required and adjusts based on the provided ratio.
func (s *server) calculateNum(ctx context.Context, num uint32, ratio float32) (n uint32) {
	min := float64(s.replica) / float64(s.gateway.GetAgentCount(ctx))
	if ratio <= 0.0 {
		return uint32(math.Ceil(float64(num) * min))
	}
	n = uint32(math.Ceil(float64(num) * (min + ((1 - min) * float64(ratio)))))
	sn := uint32(math.Ceil(float64(num) * min))
	if n-1 < sn {
		return sn
	}
	return n - 1
}

func (s *server) doSearch(
	ctx context.Context,
	cfg *payload.Search_Config,
	f func(ctx context.Context, cfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error),
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "doSearch"), apiName+"/doSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var (
		num     = int(cfg.GetNum())
		fnum    int
		replica = s.gateway.GetAgentCount(ctx)
	)

	if cfg.GetRatio() != nil {
		fnum = int(s.calculateNum(ctx, cfg.GetNum(), cfg.GetRatio().GetValue()))
	}
	if fnum <= 0 {
		fnum = num
	}

	return s.aggregationSearch(ctx, selectAggregator(cfg.GetAggregationAlgorithm(), num, fnum, replica), cfg, f)
}

func selectAggregator(algo payload.Search_AggregationAlgorithm, num, fnum, replica int) Aggregator {
	switch algo {
	case payload.Search_Unknown, payload.Search_ConcurrentQueue:
		return newStd(num, fnum, replica)
	case payload.Search_SortSlice:
		return newSlice(num, fnum, replica)
	case payload.Search_SortPoolSlice:
		return newPoolSlice(num, fnum, replica)
	case payload.Search_PairingHeap:
		return newPairingHeap(num, fnum, replica)
	default:
		return newStd(num, fnum, replica)
	}
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

func (s *server) MultiSearch(
	ctx context.Context, reqs *payload.Search_MultiRequest,
) (res *payload.Search_Responses, errs error) {
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

func (s *server) MultiSearchByID(
	ctx context.Context, reqs *payload.Search_MultiIDRequest,
) (res *payload.Search_Responses, errs error) {
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

func (s *server) LinearSearch(
	ctx context.Context, req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
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
	res, err = s.doSearch(ctx, req.GetConfig(), func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
		req.Config = fcfg
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

func (s *server) LinearSearchByID(
	ctx context.Context, req *payload.Search_IDRequest,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.LinearSearchByIDRPCName), apiName+"/"+vald.LinearSearchByIDRPCName)
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
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	if len(uuid) == 0 {
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
	vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: uuid,
		},
	})
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Unknown, vald.GetObjectRPCName+" API for "+vald.LinearSearchByIDRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
		if span != nil && st != nil && st.Code() != codes.NotFound {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		// try search by using agent's LinearSearchByID method this operation is emergency fallback, the search quality is not same as usual LinearSearchByID operation.
		res, err = s.doSearch(ctx, req.GetConfig(), func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			req.Config = fcfg
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
		Config: req.GetConfig(),
	})
	if err != nil {
		res, err = s.doSearch(ctx, req.GetConfig(), func(ctx context.Context, fcfg *payload.Search_Config, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			req.Config = fcfg
			return vc.LinearSearchByID(ctx, req, copts...)
		})
		if err == nil {
			return res, nil
		}
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal, vald.LinearSearchByIDRPCName+" API failed to process search request", reqInfo, resInfo)
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

func (s *server) StreamLinearSearchByID(
	stream vald.Search_StreamLinearSearchByIDServer,
) (err error) {
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

func (s *server) MultiLinearSearch(
	ctx context.Context, reqs *payload.Search_MultiRequest,
) (res *payload.Search_Responses, errs error) {
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

func (s *server) MultiLinearSearchByID(
	ctx context.Context, reqs *payload.Search_MultiIDRequest,
) (res *payload.Search_Responses, errs error) {
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

func (s *server) Insert(
	ctx context.Context, req *payload.Insert_Request,
) (ce *payload.Object_Location, err error) {
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

func (s *server) MultiInsert(
	ctx context.Context, reqs *payload.Insert_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
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
	eg.SetLimit(s.multiConcurrency)
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

func (s *server) Update(
	ctx context.Context, req *payload.Update_Request,
) (res *payload.Object_Location, err error) {
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
		err = s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
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
		vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Unknown, vald.GetObjectRPCName+" API for "+vald.UpdateRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
			if span != nil && st != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if conv.F32stos(vec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
			if vec.GetTimestamp() < req.GetVector().GetTimestamp() {
				return s.UpdateTimestamp(ctx, &payload.Update_TimestampRequest{
					Id:        uuid,
					Timestamp: req.GetVector().GetTimestamp(),
				})
			}
			err = errors.ErrSameVectorAlreadyExists(uuid, vec.GetVector(), req.GetVector().GetVector())
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
	if req.GetConfig().GetTimestamp() >= 0 {
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

func (s *server) MultiUpdate(
	ctx context.Context, reqs *payload.Update_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
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
	eg.SetLimit(s.multiConcurrency)
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

func (s *server) UpdateTimestamp(
	ctx context.Context, req *payload.Update_TimestampRequest,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateTimestampRPCName), apiName+"/"+vald.UpdateTimestampRPCName)
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
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateTimestampRPCName + "." + vald.GetObjectRPCName,
		ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
	}
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.UpdateTimestampRPCName+" API invalid uuid", err, reqInfo, resInfo,
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
	ts := req.GetTimestamp()
	if ts < 0 {
		err = errors.ErrInvalidTimestamp(ts)
		err = status.WrapWithInvalidArgument(vald.UpdateTimestampRPCName+" API invalid vector argument", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "timestamp",
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
	err = s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpdateTimestampRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.UpdateTimestamp(ctx, req, copts...)
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
						span.SetAttributes(trace.FromGRPCStatus(st.Code(), fmt.Sprintf("UpdateTimestamp operation for Agent %s failed,\terror: %v", target, err))...)
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
			"failed to parse "+vald.UpdateTimestampRPCName+" gRPC error response", reqInfo, resInfo, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	case len(locs.Ips) <= 0:
		err = errors.ErrIndexNotFound
		err = status.WrapWithNotFound(vald.UpdateTimestampRPCName+" API update target not found", err, reqInfo, resInfo)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	case updated.Load()+aeCount.Load() < uint64(s.replica):
		shortage := s.replica - int(updated.Load()+aeCount.Load())
		vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Unknown, vald.GetObjectRPCName+" API for "+vald.UpdateTimestampRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
			if span != nil && st != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}

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
				Vector: vec,
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
					Timestamp:            ts,
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
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Unknown, vald.InsertRPCName+" API for "+vald.UpdateTimestampRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
			if span != nil && st != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
	case updated.Load() == 0 && aeCount.Load() > 0:
		err = status.WrapWithAlreadyExists(vald.UpdateTimestampRPCName+" API update target same vector already exists", errors.ErrSameVectorAlreadyExists(uuid, nil, nil), reqInfo, resInfo)
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

func (s *server) Upsert(
	ctx context.Context, req *payload.Upsert_Request,
) (loc *payload.Object_Location, err error) {
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
		vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		var attrs trace.Attributes
		if err != nil || vec == nil {
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Unknown, vald.GetObjectRPCName+" API for "+vald.UpsertRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
			if st != nil {
				attrs = trace.FromGRPCStatus(st.Code(), msg)
				if st.Code() == codes.NotFound {
					shouldInsert = true
					err = nil
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

func (s *server) MultiUpsert(
	ctx context.Context, reqs *payload.Upsert_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
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
	eg.SetLimit(s.multiConcurrency)
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

func (s *server) Remove(
	ctx context.Context, req *payload.Remove_Request,
) (locs *payload.Object_Location, err error) {
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
	err = s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
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

func (s *server) MultiRemove(
	ctx context.Context, reqs *payload.Remove_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
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
	eg.SetLimit(s.multiConcurrency)
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

func (s *server) RemoveByTimestamp(
	ctx context.Context, req *payload.Remove_TimestampRequest,
) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveByTimestampRPCName), apiName+"/"+vald.RemoveByTimestampRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var mu sync.Mutex
	var emu sync.Mutex
	visited := make(map[string]int) // map[uuid: position of locs]
	locs = new(payload.Object_Locations)

	err := s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		sctx, sspan := trace.StartSpan(grpc.WithGRPCMethod(ctx, "BroadCast/"+target), apiName+"/removeByTimestamp/BroadCast/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()

		res, err := vc.RemoveByTimestamp(sctx, req, copts...)
		if err != nil {
			if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
				err = status.WrapWithInternal(
					vald.RemoveByTimestampRPCName+" API connection not found", err,
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.StatusCodeInternal(err.Error())...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				log.Error(err)
				emu.Lock()
				errs = errors.Join(errs, err)
				emu.Unlock()
				return nil
			}
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				vald.RemoveByTimestampRPCName+" gRPC error response",
			)
			if sspan != nil {
				sspan.RecordError(err)
				sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				sspan.SetStatus(trace.StatusError, err.Error())
			}
			if err != nil && st.Code() != codes.NotFound {
				log.Error(err)
				emu.Lock()
				errs = errors.Join(errs, err)
				emu.Unlock()
				return nil
			}
		}

		if res != nil && len(res.GetLocations()) > 0 {
			for _, loc := range res.GetLocations() {
				mu.Lock()
				if pos, ok := visited[loc.GetUuid()]; !ok {
					locs.Locations = append(locs.GetLocations(), loc)
					visited[loc.GetUuid()] = len(locs.Locations) - 1
				} else {
					if pos < len(locs.GetLocations()) {
						locs.GetLocations()[pos].Ips = append(locs.GetLocations()[pos].Ips, loc.GetIps()...)
						if s := locs.GetLocations()[pos].Name; len(s) == 0 {
							locs.GetLocations()[pos].Name = loc.GetName()
						} else {
							// strings.Join is used because '+=' causes performance degradation when the number of characters is large.
							locs.GetLocations()[pos].Name = strings.Join([]string{
								s, loc.GetName(),
							}, ",")
						}
					}
				}
				mu.Unlock()
			}
		}
		return nil
	})
	if errs != nil {
		err = errors.Join(err, errs)
	}
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveByTimestampRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Error(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if locs == nil || len(locs.GetLocations()) == 0 {
		err = status.WrapWithNotFound(
			vald.RemoveByTimestampRPCName+" API remove target not found", errors.ErrIndexNotFound,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Error(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return locs, nil
}

func (s *server) getObject(
	ctx context.Context, uuid string,
) (vec *payload.Object_Vector, err error) {
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
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/getObject/BroadCast/"+target)
			defer func() {
				if sspan != nil {
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

func (s *server) Flush(
	ctx context.Context, req *payload.Flush_Request,
) (cnts *payload.Info_Index_Count, err error) {
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
	err = s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
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
			"stored index: %d, uncommitted: %d, indexing: %t, saving: %t",
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

func (s *server) GetObject(
	ctx context.Context, req *payload.Object_VectorRequest,
) (vec *payload.Object_Vector, err error) {
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

func (s *server) StreamListObject(
	req *payload.Object_List_Request, stream vald.Object_StreamListObjectServer,
) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamListObjectRPCName), apiName+"/"+vald.StreamListObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var rmu, smu sync.Mutex
	err := s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.StreamListObjectRPCName+"/"+target)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()

		client, err := vc.StreamListObject(ctx, req, copts...)
		if err != nil {
			log.Errorf("failed to get StreamListObject client for agent(%s): %v", target, err)
			return err
		}

		eg, ctx := errgroup.WithContext(ctx)
		ectx, ecancel := context.WithCancel(ctx)
		defer ecancel()
		eg.SetLimit(s.streamConcurrency)

		for {
			select {
			case <-ectx.Done():
				var err error
				if !errors.Is(ctx.Err(), context.Canceled) {
					err = errors.Join(err, ctx.Err())
				}
				if egerr := eg.Wait(); err != nil {
					err = errors.Join(err, egerr)
				}
				return err
			default:
				eg.Go(safety.RecoverFunc(func() error {
					rmu.Lock()
					res, err := client.Recv()
					rmu.Unlock()
					if err != nil {
						if errors.Is(err, io.EOF) {
							ecancel()
							return nil
						}
						return errors.ErrServerStreamClientRecv(err)
					}

					vec := res.GetVector()
					if vec == nil {
						st := res.GetStatus()
						log.Warnf("received empty vector: code %v: details %v: message %v",
							st.GetCode(),
							st.GetDetails(),
							st.GetMessage(),
						)
						return nil
					}

					smu.Lock()
					err = stream.Send(res)
					smu.Unlock()
					if err != nil {
						if sspan != nil {
							st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamListObject send gRPC error response")
							sspan.RecordError(err)
							sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
							sspan.SetStatus(trace.StatusError, err.Error())
						}
						return errors.ErrServerStreamServerSend(err)
					}

					return nil
				}))
			}
		}
	})
	return err
}

func (s *server) IndexInfo(
	ctx context.Context, _ *payload.Empty,
) (vec *payload.Info_Index_Count, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+vald.IndexInfoRPCName), apiName+"/"+vald.IndexInfoRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ech := make(chan error, 1)
	var (
		stored, uncommitted atomic.Uint32
		indexing, saving    atomic.Bool
	)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.IndexInfoRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			info, err := vc.IndexInfo(sctx, new(payload.Empty), copts...)
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
							"/vald.v1." + vald.IndexInfoRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.IndexInfoRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.IndexInfoRPCName+" API",
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexInfoRPCName + ".BroadCase/" + target,
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
			if info != nil {
				stored.Add(info.GetStored())
				uncommitted.Add(info.GetUncommitted())
				if info.GetIndexing() {
					indexing.Store(true)
				}
				if info.GetSaving() {
					saving.Store(true)
				}
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ech:
	}
	if err != nil {
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexInfoRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.IndexInfoRPCName+" API connection not found", err, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.IndexInfoRPCName+" API canceled", err, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.IndexInfoRPCName+" API deadline exceeded", err, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Unknown, vald.IndexInfoRPCName+" API request returned error", resInfo)
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
	return &payload.Info_Index_Count{
		Stored:      stored.Load(),
		Uncommitted: uncommitted.Load(),
		Indexing:    indexing.Load(),
		Saving:      saving.Load(),
	}, nil
}

func (s *server) IndexDetail(
	ctx context.Context, _ *payload.Empty,
) (vec *payload.Info_Index_Detail, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+vald.IndexDetailRPCName), apiName+"/"+vald.IndexDetailRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ech := make(chan error, 1)
	var (
		mu     sync.Mutex
		detail = &payload.Info_Index_Detail{
			Counts:     make(map[string]*payload.Info_Index_Count),
			Replica:    uint32(s.replica),
			LiveAgents: uint32(s.gateway.GetAgentCount(ctx)),
		}
	)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.IndexDetailRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			info, err := vc.IndexInfo(sctx, new(payload.Empty), copts...)
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
							"/vald.v1." + vald.IndexDetailRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.IndexDetailRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.IndexDetailRPCName+" API",
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexDetailRPCName + ".BroadCase/" + target,
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
			if info != nil {
				mu.Lock()
				detail.Counts[target] = info
				mu.Unlock()
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ech:
	}
	if err != nil {
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexDetailRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.IndexDetailRPCName+" API connection not found", err, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.IndexDetailRPCName+" API canceled", err, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.IndexDetailRPCName+" API deadline exceeded", err, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Unknown, vald.IndexDetailRPCName+" API request returned error", resInfo)
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
	return detail, nil
}

func (s *server) GetTimestamp(
	ctx context.Context, req *payload.Object_TimestampRequest,
) (ts *payload.Object_Timestamp, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetTimestampRPCName), apiName+"/"+vald.GetTimestampRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetId().GetId()
	tch := make(chan *payload.Object_Timestamp, 1)
	ech := make(chan error, 1)
	doneErr := errors.New("done getTimestamp")
	ctx, cancel := context.WithCancelCause(ctx)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(tch)
		defer close(ech)
		var once sync.Once
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/getTimestamp/BroadCast/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			req := &payload.Object_TimestampRequest{
				Id: &payload.Object_ID{
					Id: uuid,
				},
			}
			ots, err := vc.GetTimestamp(sctx, req, copts...)
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
							"/vald.v1." + vald.GetTimestampRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.GetTimestampRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.GetTimestampRPCName+" API meta "+uuid+"'s uuid not found",
						&errdetails.RequestInfo{
							RequestId:   uuid,
							ServingData: errdetails.Serialize(req),
						},
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetTimestampRPCName,
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
			if ots != nil && ots.GetId() != "" {
				once.Do(func() {
					tch <- ots
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
			case ts = <-tch:
				if ts == nil || ts.GetId() == "" {
					err = errors.ErrObjectNotFound(nil, uuid)
				} else {
					err = nil
				}
			default:
			}
		}
	case ts = <-tch:
		if ts == nil || ts.GetId() == "" {
			err = errors.ErrObjectNotFound(nil, uuid)
		}
	case err = <-ech:
	}
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   uuid,
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetTimestampRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrInvalidUUID(uuid)):
			err = status.WrapWithInvalidArgument(vald.GetTimestampRPCName+" API invalid argument for uuid \""+uuid+"\" detected", err, reqInfo, resInfo, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "uuid",
						Description: err.Error(),
					},
				},
			})
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		case errors.Is(err, errors.ErrObjectIDNotFound(uuid)), errors.Is(err, errors.ErrObjectNotFound(nil, uuid)):
			err = status.WrapWithNotFound(vald.GetTimestampRPCName+" API id "+uuid+"'s object not found", err, reqInfo, resInfo)
			attrs = trace.StatusCodeNotFound(err.Error())
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.GetTimestampRPCName+" API connection not found", err, reqInfo, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.GetTimestampRPCName+" API canceled", err, reqInfo, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.GetTimestampRPCName+" API deadline exceeded", err, reqInfo, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Unknown, vald.GetTimestampRPCName+" API uuid "+uuid+"'s request returned error", reqInfo, resInfo)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return ts, nil
}

func (s *server) IndexStatistics(
	ctx context.Context, req *payload.Empty,
) (vec *payload.Info_Index_Statistics, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+vald.IndexStatisticsRPCName), apiName+"/"+vald.IndexStatisticsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	details, err := s.IndexStatisticsDetail(ctx, req)
	if err != nil || details == nil {
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexStatisticsRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.IndexStatisticsRPCName+" API connection not found", err, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.IndexStatisticsRPCName+" API canceled", err, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.IndexStatisticsRPCName+" API deadline exceeded", err, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Unknown, vald.IndexStatisticsRPCName+" API request returned error", resInfo)
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
	return mergeInfoIndexStatistics(details.GetDetails()), nil
}

func (s *server) IndexStatisticsDetail(
	ctx context.Context, _ *payload.Empty,
) (vec *payload.Info_Index_StatisticsDetail, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+vald.IndexStatisticsDetailRPCName), apiName+"/"+vald.IndexStatisticsDetailRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ech := make(chan error, 1)
	var (
		mu     sync.Mutex
		detail = &payload.Info_Index_StatisticsDetail{
			Details: make(map[string]*payload.Info_Index_Statistics, s.gateway.GetAgentCount(ctx)),
		}
	)
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.IndexStatisticsDetailRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			var stats *payload.Info_Index_Statistics
			stats, err = vc.IndexStatistics(sctx, new(payload.Empty), copts...)
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
							"/vald.v1." + vald.IndexStatisticsDetailRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.IndexStatisticsDetailRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.IndexStatisticsDetailRPCName+" API",
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexStatisticsDetailRPCName + ".BroadCase/" + target,
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
			if stats != nil {
				mu.Lock()
				detail.Details[target] = stats
				mu.Unlock()
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ech:
	}
	if err != nil {
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexStatisticsDetailRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.IndexStatisticsDetailRPCName+" API connection not found", err, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.IndexStatisticsDetailRPCName+" API canceled", err, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.IndexStatisticsDetailRPCName+" API deadline exceeded", err, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Unknown, vald.IndexStatisticsDetailRPCName+" API request returned error", resInfo)
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
	return detail, nil
}

func calculateMedian(data []int32) int32 {
	slices.Sort(data)
	n := len(data)
	if n%2 == 0 {
		return (data[n/2-1] + data[n/2]) / 2
	}
	return data[n/2]
}

func sumHistograms(hist1, hist2 []uint64) []uint64 {
	if len(hist1) < len(hist2) {
		hist1, hist2 = hist2, hist1
	}
	for i := range hist2 {
		hist1[i] += hist2[i]
	}
	return hist1
}

func mergeInfoIndexStatistics(
	stats map[string]*payload.Info_Index_Statistics,
) (merged *payload.Info_Index_Statistics) {
	merged = new(payload.Info_Index_Statistics)

	if len(stats) == 0 {
		return merged
	}

	var indegrees, outdegrees []int32
	var indegreeCounts [][]int64
	var outdegreeHistograms, indegreeHistograms [][]uint64
	merged.Valid = true

	for _, stat := range stats {
		if !stat.Valid {
			continue
		}
		indegrees = append(indegrees, stat.MedianIndegree)
		outdegrees = append(outdegrees, stat.MedianOutdegree)

		indegreeCounts = append(indegreeCounts, stat.IndegreeCount)
		outdegreeHistograms = append(outdegreeHistograms, stat.OutdegreeHistogram)
		indegreeHistograms = append(indegreeHistograms, stat.IndegreeHistogram)

		if stat.MaxNumberOfIndegree > merged.MaxNumberOfIndegree {
			merged.MaxNumberOfIndegree = stat.MaxNumberOfIndegree
		}
		if stat.MaxNumberOfOutdegree > merged.MaxNumberOfOutdegree {
			merged.MaxNumberOfOutdegree = stat.MaxNumberOfOutdegree
		}
		if stat.MinNumberOfIndegree < merged.MinNumberOfIndegree || merged.MinNumberOfIndegree == 0 {
			merged.MinNumberOfIndegree = stat.MinNumberOfIndegree
		}
		if stat.MinNumberOfOutdegree < merged.MinNumberOfOutdegree || merged.MinNumberOfOutdegree == 0 {
			merged.MinNumberOfOutdegree = stat.MinNumberOfOutdegree
		}
		merged.ModeIndegree += stat.ModeIndegree
		merged.ModeOutdegree += stat.ModeOutdegree
		merged.NodesSkippedFor10Edges += stat.NodesSkippedFor10Edges
		merged.NodesSkippedForIndegreeDistance += stat.NodesSkippedForIndegreeDistance
		merged.NumberOfEdges += stat.NumberOfEdges
		merged.NumberOfIndexedObjects += stat.NumberOfIndexedObjects
		merged.NumberOfNodes += stat.NumberOfNodes
		merged.NumberOfNodesWithoutEdges += stat.NumberOfNodesWithoutEdges
		merged.NumberOfNodesWithoutIndegree += stat.NumberOfNodesWithoutIndegree
		merged.NumberOfObjects += stat.NumberOfObjects
		merged.NumberOfRemovedObjects += stat.NumberOfRemovedObjects
		merged.SizeOfObjectRepository += stat.SizeOfObjectRepository
		merged.SizeOfRefinementObjectRepository += stat.SizeOfRefinementObjectRepository

		merged.VarianceOfIndegree += stat.VarianceOfIndegree
		merged.VarianceOfOutdegree += stat.VarianceOfOutdegree
		merged.MeanEdgeLength += stat.MeanEdgeLength
		merged.MeanEdgeLengthFor10Edges += stat.MeanEdgeLengthFor10Edges
		merged.MeanIndegreeDistanceFor10Edges += stat.MeanIndegreeDistanceFor10Edges
		merged.MeanNumberOfEdgesPerNode += stat.MeanNumberOfEdgesPerNode

		merged.C1Indegree += stat.C1Indegree
		merged.C5Indegree += stat.C5Indegree
		merged.C95Outdegree += stat.C95Outdegree
		merged.C99Outdegree += stat.C99Outdegree
	}

	merged.MedianIndegree = calculateMedian(indegrees)
	merged.MedianOutdegree = calculateMedian(outdegrees)
	merged.IndegreeCount = make([]int64, len(indegreeCounts[0]))
	for i := range merged.IndegreeCount {
		var (
			alen int64
			sum  int64
		)
		for _, count := range indegreeCounts {
			if i < len(count) {
				alen++
				sum += count[i]
			}
		}
		merged.IndegreeCount[i] = sum / alen
	}

	for _, hist := range outdegreeHistograms {
		merged.OutdegreeHistogram = sumHistograms(merged.OutdegreeHistogram, hist)
	}

	for _, hist := range indegreeHistograms {
		merged.IndegreeHistogram = sumHistograms(merged.IndegreeHistogram, hist)
	}

	merged.ModeIndegree /= uint64(len(stats))
	merged.ModeOutdegree /= uint64(len(stats))
	merged.VarianceOfIndegree /= float64(len(stats))
	merged.VarianceOfOutdegree /= float64(len(stats))
	merged.MeanEdgeLength /= float64(len(stats))
	merged.MeanEdgeLengthFor10Edges /= float64(len(stats))
	merged.MeanIndegreeDistanceFor10Edges /= float64(len(stats))
	merged.MeanNumberOfEdgesPerNode /= float64(len(stats))
	merged.C1Indegree /= float64(len(stats))
	merged.C5Indegree /= float64(len(stats))
	merged.C95Outdegree /= float64(len(stats))
	merged.C99Outdegree /= float64(len(stats))

	return merged
}

func (s *server) IndexProperty(
	ctx context.Context, _ *payload.Empty,
) (detail *payload.Info_Index_PropertyDetail, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.IndexRPCServiceName+"/"+vald.IndexPropertyRPCName), apiName+"/"+vald.IndexPropertyRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ech := make(chan error, 1)
	var mu sync.Mutex
	detail = &payload.Info_Index_PropertyDetail{
		Details: make(map[string]*payload.Info_Index_Property, s.gateway.GetAgentCount(ctx)),
	}

	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.IndexStatisticsDetailRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			var prop *payload.Info_Index_PropertyDetail
			prop, err = vc.IndexProperty(sctx, new(payload.Empty), copts...)
			if err != nil {
				var (
					attrs trace.Attributes
					st    *status.Status
					msg   string
					code  codes.Code
				)
				switch {
				case errors.Is(err, context.Canceled), errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					attrs = trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexPropertyRPCName + ".BroadCast/" + target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded), errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexPropertyRPCName + ".BroadCast/" + target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+vald.IndexPropertyRPCName+" API",
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexPropertyRPCName + ".BroadCast/" + target,
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
				if err != nil && st != nil && code != codes.Canceled && code != codes.DeadlineExceeded && code != codes.InvalidArgument && code != codes.NotFound && code != codes.OK && code != codes.Unimplemented {
					return err
				}
				return nil
			}
			if prop != nil {
				mu.Lock()
				for key, value := range prop.Details {
					detail.Details[target+"-"+key] = value
				}
				mu.Unlock()
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ech:
	}
	if err != nil {
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.IndexPropertyRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(vald.IndexPropertyRPCName+" API connection not found", err, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(vald.IndexPropertyRPCName+" API canceled", err, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(vald.IndexPropertyRPCName+" API deadline exceeded", err, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Unknown, vald.IndexPropertyRPCName+" API request returned error", resInfo)
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
	return detail, nil
}
