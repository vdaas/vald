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
	"math"
	"math/big"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	vclient "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/pkg/gateway/internal/location"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type MirrorServer interface {
	vald.Server
	mirror.MirrorServer
}

type server struct {
	eg                errgroup.Group
	gateway           service.Gateway
	client            vclient.Client
	timeout           time.Duration
	replica           int
	streamConcurrency int
	name              string
	ip                string
	vald.UnimplementedValdServer
	mirror.MirrorServer
}

const apiName = "vald/gateway/lb"

func New(opts ...Option) MirrorServer {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
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
	res, err = s.search(ctx, &payload.Search_Config{
		RequestId:      cfg.GetRequestId(),
		Num:            cfg.GetNum(),
		MinNum:         mn,
		Radius:         cfg.GetRadius(),
		Epsilon:        cfg.GetEpsilon(),
		Timeout:        cfg.GetTimeout(),
		IngressFilters: cfg.GetIngressFilters(),
		EgressFilters:  cfg.GetEgressFilters(),
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
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.SearchByIDRPCName+" API invalid uuid", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	oreq := &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: uuid,
		},
		Filters: req.GetConfig().GetEgressFilters(),
	}
	vec, err := s.GetObject(ctx, oreq)
	cfg := req.GetConfig()
	mn := cfg.GetMinNum()
	if req.Config != nil {
		req.Config.MinNum = 0
	}
	scfg := &payload.Search_Config{
		RequestId:      cfg.GetRequestId(),
		Num:            cfg.GetNum(),
		MinNum:         mn,
		Radius:         cfg.GetRadius(),
		Epsilon:        cfg.GetEpsilon(),
		Timeout:        cfg.GetTimeout(),
		IngressFilters: cfg.GetIngressFilters(),
		EgressFilters:  cfg.GetEgressFilters(),
	}
	if err != nil {
		_, _, err := status.ParseError(err, codes.NotFound, vald.SearchByIDRPCName+" API failed to get uuid "+req.GetId()+"'s object",
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(oreq),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName + "." + vald.GetObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		var serr error
		res, serr = s.search(ctx, scfg, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.SearchByID(ctx, req, copts...)
		})
		if serr == nil {
			return res, nil
		}
		err = errors.Wrap(err, serr.Error())
		st, msg, serr := status.ParseError(err, codes.Internal, vald.SearchByIDRPCName+" API failed to process search request",
			&errdetails.RequestInfo{
				RequestId: req.GetConfig().GetRequestId(),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, errors.Wrap(err, serr.Error())
	}
	res, err = s.Search(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: scfg,
	})
	if err != nil {
		_, _, err := status.ParseError(err, codes.Internal, vald.SearchByIDRPCName+" API failed to process search request",
			&errdetails.RequestInfo{
				RequestId: req.GetConfig().GetRequestId(),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName + "." + vald.SearchRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		var serr error
		res, serr = s.search(ctx, scfg, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.SearchByID(ctx, req, copts...)
		})
		if serr == nil {
			return res, nil
		}
		st, msg, serr := status.ParseError(serr, codes.Internal, vald.SearchByIDRPCName+" API failed to process search request",
			&errdetails.RequestInfo{
				RequestId: req.GetConfig().GetRequestId(),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		err = errors.Wrap(err, serr.Error())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

type DistPayload struct {
	raw      *payload.Object_Distance
	distance *big.Float
}

func (s *server) search(ctx context.Context, cfg *payload.Search_Config,
	f func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error)) (
	res *payload.Search_Response, err error,
) {
	ctx, span := trace.StartSpan(ctx, apiName+".search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	num := int(cfg.GetNum())
	min := int(cfg.GetMinNum())
	res = new(payload.Search_Response)
	res.Results = make([]*payload.Object_Distance, 0, s.gateway.GetAgentCount(ctx)*num)
	dch := make(chan DistPayload, cap(res.GetResults())/2)
	eg, ectx := errgroup.New(ctx)
	var cancel context.CancelFunc
	var timeout time.Duration
	if to := cfg.GetTimeout(); to != 0 {
		timeout = time.Duration(to)
	} else {
		timeout = s.timeout
	}

	var maxDist atomic.Value
	maxDist.Store(big.NewFloat(math.MaxFloat64))
	ectx, cancel = context.WithTimeout(ectx, timeout)
	eg.Go(safety.RecoverFunc(func() error {
		defer cancel()
		visited := new(sync.Map)
		return s.gateway.BroadCast(ectx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+".search/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := f(sctx, vc, copts...)
			switch {
			case errors.Is(err, context.Canceled),
				errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1.search.BroadCast/" +
							target + " canceled: " + err.Error())...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
			case errors.Is(err, context.DeadlineExceeded),
				errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1.search.BroadCast/" +
							target + " deadline_exceeded: " + err.Error())...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
			case err != nil:
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse search gRPC error response",
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
						ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
					})
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				switch st.Code() {
				case codes.Internal,
					codes.Unavailable,
					codes.ResourceExhausted:
					return err
				}
			case r == nil || len(r.GetResults()) == 0:
				err = status.WrapWithNotFound("failed to process search request", errors.ErrEmptySearchResult,
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
						ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
					})
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
			}
			for _, dist := range r.GetResults() {
				if dist == nil {
					continue
				}
				fdist := big.NewFloat(float64(dist.GetDistance()))
				bf, ok := maxDist.Load().(*big.Float)
				if !ok || fdist.Cmp(bf) >= 0 {
					return nil
				}
				if _, already := visited.LoadOrStore(dist.GetId(), struct{}{}); !already {
					select {
					case <-ectx.Done():
						return nil
					case dch <- DistPayload{raw: dist, distance: fdist}:
					}
				}
			}
			return nil
		})
	}))
	add := func(distance *big.Float, dist *payload.Object_Distance) {
		rl := len(res.GetResults()) // result length
		fmax, ok := maxDist.Load().(*big.Float)
		if !ok {
			return
		}
		if rl >= num && distance.Cmp(fmax) >= 0 {
			return
		}
		switch rl {
		case 0:
			res.Results = append(res.GetResults(), dist)
		case 1:

			if distance.Cmp(big.NewFloat(float64(res.GetResults()[0].GetDistance()))) >= 0 {
				res.Results = append(res.GetResults(), dist)
			} else {
				res.Results = []*payload.Object_Distance{dist, res.GetResults()[0]}
			}
		default:
			pos := rl
			for idx := rl; idx >= 1; idx-- {
				if distance.Cmp(big.NewFloat(float64(res.GetResults()[idx-1].GetDistance()))) >= 0 {
					pos = idx - 1
					break
				}
			}
			switch {
			case pos == rl:
				res.Results = append([]*payload.Object_Distance{dist}, res.GetResults()...)
			case pos == rl-1:
				res.Results = append(res.GetResults(), dist)
			case pos >= 0:
				res.Results = append(res.GetResults()[:pos+1], res.GetResults()[pos:]...)
				res.Results[pos+1] = dist
			}
		}
		rl = len(res.GetResults())
		if rl > num && num != 0 {
			res.Results = res.GetResults()[:num]
			rl = len(res.GetResults())
		}
		if distEnd := big.NewFloat(float64(res.GetResults()[rl-1].GetDistance())); rl >= num &&
			distEnd.Cmp(fmax) < 0 {
			maxDist.Store(distEnd)
		}
	}
	for {
		select {
		case <-ectx.Done():
			err = eg.Wait()
			close(dch)
			// range over channel patter to check remaining channel's data for vald's search accuracy
			for dist := range dch {
				add(dist.distance, dist.raw)
			}
			if num != 0 && len(res.GetResults()) > num {
				res.Results = res.GetResults()[:num]
			}

			if errors.Is(ectx.Err(), context.DeadlineExceeded) {
				if len(res.GetResults()) == 0 {
					err = status.WrapWithDeadlineExceeded(
						"error search result length is 0 due to the timeoutage limit",
						errors.ErrEmptySearchResult,
						&errdetails.RequestInfo{
							RequestId:   cfg.GetRequestId(),
							ServingData: errdetails.Serialize(cfg),
						},
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
							ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
						}, info.Get(),
					)
					if span != nil {
						span.RecordError(err)
						span.SetAttributes(trace.StatusCodeDeadlineExceeded(err.Error())...)
						span.SetStatus(trace.StatusError, err.Error())
					}
					return nil, err
				}
				if 0 < min && len(res.GetResults()) < min {
					err = status.WrapWithDeadlineExceeded(
						fmt.Sprintf("error search result length is not enough due to the timeoutage limit, required: %d, found: %d", min, len(res.GetResults())),
						errors.ErrInsuffcientSearchResult,
						&errdetails.RequestInfo{
							RequestId:   cfg.GetRequestId(),
							ServingData: errdetails.Serialize(cfg),
						},
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
							ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
						}, info.Get(),
					)
					if span != nil {
						span.RecordError(err)
						span.SetAttributes(trace.StatusCodeDeadlineExceeded(err.Error())...)
						span.SetStatus(trace.StatusError, err.Error())
					}
					return nil, err
				}
			}

			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal,
					"failed to parse search gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   cfg.GetRequestId(),
						ServingData: errdetails.Serialize(cfg),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
						ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
					}, info.Get())
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				log.Warn(err)
				if len(res.GetResults()) == 0 {
					return nil, err
				}
			}
			if num != 0 && len(res.GetResults()) == 0 {
				if err == nil {
					err = errors.ErrEmptySearchResult
				}
				err = status.WrapWithNotFound("error search result length is 0", err,
					&errdetails.RequestInfo{
						RequestId:   cfg.GetRequestId(),
						ServingData: errdetails.Serialize(cfg),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
						ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
					}, info.Get())
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil, err
			}

			if 0 < min && len(res.GetResults()) < min {
				if err == nil {
					err = errors.ErrInsuffcientSearchResult
				}
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				err = status.WrapWithNotFound(
					fmt.Sprintf("error search result length is not enough required: %d, found: %d", min, len(res.GetResults())),
					errors.ErrInsuffcientSearchResult,
					&errdetails.RequestInfo{
						RequestId:   cfg.GetRequestId(),
						ServingData: errdetails.Serialize(cfg),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.search",
						ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
					}, info.Get(),
				)
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil, err
			}
			res.RequestId = cfg.GetRequestId()
			return res, nil
		case dist := <-dch:
			add(dist.distance, dist.raw)
		}
	}
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

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		vl := len(req.GetVector())
		if vl < algorithm.MinimumVectorDimensionSize {
			err := errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument(vald.MultiSearchRPCName+" API invalid vector argument", err,
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
				}, info.Get())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiSearchRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
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
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		})
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
		return nil, err
	}
	return res, nil
}

func (s *server) MultiSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(func() error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiSearchByIDRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			defer wg.Done()
			r, err := s.SearchByID(sctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal,
					"failed to parse "+vald.SearchByIDRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		})
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse "+vald.MultiSearchByIDRPCName+" gRPC error response",
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
	res, err = s.search(ctx, &payload.Search_Config{
		RequestId:      cfg.GetRequestId(),
		Num:            cfg.GetNum(),
		MinNum:         mn,
		Timeout:        cfg.GetTimeout(),
		IngressFilters: cfg.GetIngressFilters(),
		EgressFilters:  cfg.GetEgressFilters(),
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
	if len(req.GetId()) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.LinearSearchByIDRPCName+" API invalid uuid", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
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
	oreq := &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetId(),
		},
		Filters: req.GetConfig().GetEgressFilters(),
	}
	vec, err := s.GetObject(ctx, oreq)
	cfg := req.GetConfig()
	mn := cfg.GetMinNum()
	if req.Config != nil {
		req.Config.MinNum = 0
	}
	scfg := &payload.Search_Config{
		RequestId:      cfg.GetRequestId(),
		Num:            cfg.GetNum(),
		MinNum:         mn,
		Timeout:        cfg.GetTimeout(),
		IngressFilters: cfg.GetIngressFilters(),
		EgressFilters:  cfg.GetEgressFilters(),
	}
	if err != nil {
		_, _, err := status.ParseError(err, codes.NotFound, fmt.Sprintf(vald.LinearSearchByIDRPCName+" API failed to get uuid %s's object", req.GetId()),
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(oreq),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName + "." + vald.GetObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		var serr error
		res, serr = s.search(ctx, scfg, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.LinearSearchByID(ctx, req, copts...)
		})
		if serr == nil {
			return res, nil
		}
		err = errors.Wrap(err, serr.Error())
		st, msg, serr := status.ParseError(err, codes.Internal, vald.LinearSearchByIDRPCName+" API failed to process search request",
			&errdetails.RequestInfo{
				RequestId: req.GetConfig().GetRequestId(),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, errors.Wrap(err, serr.Error())
	}

	res, err = s.LinearSearch(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: scfg,
	})
	if err != nil {
		_, _, err := status.ParseError(err, codes.Internal, vald.LinearSearchByIDRPCName+" API failed to process search request",
			&errdetails.RequestInfo{
				RequestId: req.GetConfig().GetRequestId(),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName + "." + vald.LinearSearchRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		var serr error
		res, serr = s.search(ctx, scfg, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.LinearSearchByID(ctx, req, copts...)
		})
		if serr == nil {
			return res, nil
		}
		st, msg, serr := status.ParseError(serr, codes.Internal, vald.LinearSearchByIDRPCName+" API failed to process search request",
			&errdetails.RequestInfo{
				RequestId: req.GetConfig().GetRequestId(),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		err = errors.Wrap(err, serr.Error())
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

func (s *server) MultiLinearSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		vl := len(req.GetVector())
		if vl < algorithm.MinimumVectorDimensionSize {
			err := errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument(vald.MultiLinearSearchRPCName+" API invalid vector argument", err,
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
				}, info.Get())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiLinearSearchRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.LinearSearch(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		})
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
		return nil, err
	}
	return res, nil
}

func (s *server) MultiLinearSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(func() error {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiLinearSearchByIDRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			defer wg.Done()
			r, err := s.LinearSearchByID(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal,
					"failed to parse "+vald.LinearSearchByIDRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		})
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse "+vald.MultiLinearSearchByIDRPCName+" gRPC error response",
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
	uuid := req.GetVector().GetId()
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.InsertRPCName+" API invalid uuid", err,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
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
		err = status.WrapWithInvalidArgument(vald.InsertRPCName+" API invalid vector argument", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
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
		id, err := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if err == nil && id != nil && len(id.GetId()) != 0 {
			if err == nil {
				err = errors.ErrMetaDataAlreadyExists(uuid)
			}
			st, msg, err := status.ParseError(err, codes.AlreadyExists,
				"error "+vald.InsertRPCName+" API ID = "+uuid+" already exists",
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + "." + vald.ExistsRPCName,
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
			req.Config = &payload.Insert_Config{SkipStrictExistCheck: true}
		}
	}

	mu := new(sync.Mutex)
	ce = &payload.Object_Location{
		Uuid: uuid,
		Ips:  make([]string, 0, s.replica),
	}
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
		ctx, span := trace.StartSpan(ctx, apiName+"."+vald.InsertRPCName+"/"+target)
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
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				emu.Unlock()
				return err
			}
			return nil
		}
		mu.Lock()
		ce.Ips = append(ce.GetIps(), loc.GetIps()...)
		ce.Name = loc.GetName()
		mu.Unlock()
		return nil
	})
	if err != nil {
		if errs == nil {
			errs = err
		} else {
			errs = errors.Wrap(errs, err.Error())
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
	log.Debugf("Insert API insert succeeded to %#v", ce)
	return ce, nil
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
	vecs := reqs.GetRequests()
	ids := make([]string, 0, len(vecs))
	now := time.Now().UnixNano()
	for i, req := range vecs {
		uuid := req.GetVector().GetId()
		vector := req.GetVector().GetVector()
		vl := len(vector)
		if vl < algorithm.MinimumVectorDimensionSize {
			err = errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument(vald.MultiInsertRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       fmt.Sprintf("vector dimension size for id: %s", uuid),
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
			id, err := s.Exists(ctx, &payload.Object_ID{
				Id: uuid,
			})
			if err == nil && id != nil && len(id.GetId()) != 0 {
				if err == nil {
					err = errors.ErrMetaDataAlreadyExists(uuid)
				}
				st, msg, err := status.ParseError(err, codes.AlreadyExists,
					"error "+vald.MultiInsertRPCName+" API ID = "+uuid+" already exists",
					&errdetails.RequestInfo{
						RequestId:   uuid,
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiInsertRPCName + "." + vald.ExistsRPCName,
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
				reqs.GetRequests()[i].GetConfig().SkipStrictExistCheck = true
			} else {
				reqs.GetRequests()[i].Config = &payload.Insert_Config{SkipStrictExistCheck: true}
			}
		}
		if reqs.GetRequests()[i].GetConfig().GetTimestamp() == 0 {
			if reqs.GetRequests()[i].GetConfig() == nil {
				reqs.GetRequests()[i].Config = &payload.Insert_Config{
					Timestamp: now,
				}
			} else {
				reqs.GetRequests()[i].GetConfig().Timestamp = now
			}
		}
		ids = append(ids, uuid)
	}

	mu := new(sync.Mutex)
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, 0, s.replica),
	}

	emu := new(sync.Mutex)
	var errs error
	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(ctx, apiName+"."+vald.MultiInsertRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.MultiInsert(ctx, reqs, copts...)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled),
				errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.MultiInsertRPCName + ".DoMulti/" +
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
							"/vald.v1." + vald.MultiInsertRPCName + ".DoMulti/" +
							target + " deadline_exceeded: " + err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil
			}
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.MultiInsertRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   strings.Join(ids, ","),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiInsertRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				})
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}

			if err != nil {
				emu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				emu.Unlock()
			}
			return err
		}
		mu.Lock()
		locs.Locations = append(locs.GetLocations(), loc.Locations...)
		mu.Unlock()
		return nil
	})
	if err != nil {
		if errs == nil {
			errs = err
		} else {
			errs = errors.Wrap(errs, err.Error())
		}
	}

	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse "+vald.MultiInsertRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ", "),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiInsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return location.ReStructure(ids, locs), nil
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetVector().GetId()
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.UpdateRPCName+" API invalid uuid", err,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
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
		err = status.WrapWithInvalidArgument(vald.UpdateRPCName+" API invalid vector argument", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
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
		vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		if err != nil || vec == nil || len(vec.GetId()) == 0 {
			if err == nil {
				err = errors.ErrObjectIDNotFound(uuid)
			}
			st, msg, err := status.ParseError(err, codes.NotFound,
				"error "+vald.UpdateRPCName+" API ID = "+uuid+" not found",
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
	vecs := reqs.GetRequests()
	ids := make([]string, 0, len(vecs))
	ireqs := make([]*payload.Insert_Request, 0, len(vecs))
	rreqs := make([]*payload.Remove_Request, 0, len(vecs))
	now := time.Now().UnixNano()
	for _, req := range vecs {
		vl := len(req.GetVector().GetVector())
		if vl < algorithm.MinimumVectorDimensionSize {
			err = errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument(vald.MultiUpdateRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(reqs),
				},
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
		uuid := req.GetVector().GetId()
		if !req.GetConfig().GetSkipStrictExistCheck() {
			vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
				Id: &payload.Object_ID{
					Id: uuid,
				},
			})
			if err != nil || vec == nil || len(vec.GetId()) == 0 {
				if err == nil {
					err = errors.ErrObjectIDNotFound(uuid)
				}
				st, msg, err := status.ParseError(err, codes.NotFound,
					"error "+vald.MultiUpdateRPCName+" API ID = "+uuid+" not fount",
					&errdetails.RequestInfo{
						RequestId:   uuid,
						ServingData: errdetails.Serialize(reqs),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpdateRPCName + "." + vald.GetObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
					}, info.Get())
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil, err
			}
			if conv.F32stos(vec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
				log.Warn(errors.ErrSameVectorAlreadyExists(uuid, vec.GetVector(), req.GetVector().GetVector()))
				continue
			}
			if req.GetConfig() != nil {
				req.Config.SkipStrictExistCheck = true
			} else {
				req.Config = &payload.Update_Config{SkipStrictExistCheck: true}
			}
		}
		var n int64
		if req.GetConfig().GetTimestamp() != 0 {
			n = req.GetConfig().GetTimestamp()
		} else {
			n = now
		}
		ids = append(ids, req.GetVector().GetId())
		rreqs = append(rreqs, &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: req.GetVector().GetId(),
			},
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: true,
				Timestamp:            n,
			},
		})
		n++
		ireqs = append(ireqs, &payload.Insert_Request{
			Vector: req.GetVector(),
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Filters:              req.GetConfig().GetFilters(),
				Timestamp:            n,
			},
		})
	}
	locs, err := s.MultiRemove(ctx, &payload.Remove_MultiRequest{
		Requests: rreqs,
	})
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.MultiRemoveRPCName+" for "+vald.MultiUpdateRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(rreqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpdateRPCName + "." + vald.MultiRemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	log.Debugf("uuids %v were removed from %v due to MultiUpdate. "+vald.MultiInsertRPCName+" will be executed for them soon. Please see detail %#v", ids, locs.GetLocations(), locs)
	locs, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
		Requests: ireqs,
	})
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.MultiInsertRPCName+" for "+vald.MultiUpdateRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(ireqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpdateRPCName + "." + vald.MultiInsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return locs, nil
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	vec := req.GetVector()
	uuid := vec.GetId()
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.UpsertRPCName+" API invalid uuid", err,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
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
		err = status.WrapWithInvalidArgument(vald.UpsertRPCName+" API invalid vector argument", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
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
		if err != nil || vec == nil || len(vec.GetId()) == 0 {
			shouldInsert = true
		} else if conv.F32stos(vec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
			if err == nil {
				err = errors.ErrSameVectorAlreadyExists(uuid, vec.GetVector(), req.GetVector().GetVector())
			}
			st, msg, err := status.ParseError(err, codes.AlreadyExists,
				"error "+vald.UpdateRPCName+" for "+vald.UpsertRPCName+" API ID = "+uuid+"'s same vector data already exists",
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + vald.GetObjectRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				}, info.Get())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
	} else {
		id, err := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
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
				SkipStrictExistCheck: true,
				Filters:              req.GetConfig().GetFilters(),
				Timestamp:            req.GetConfig().GetTimestamp(),
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

	insertReqs := make([]*payload.Insert_Request, 0, len(reqs.GetRequests()))
	updateReqs := make([]*payload.Update_Request, 0, len(reqs.GetRequests()))

	ids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		vec := req.GetVector()
		uuid := vec.GetId()
		vl := len(vec.GetVector())
		if vl < algorithm.MinimumVectorDimensionSize {
			err = errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument(vald.MultiUpsertRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
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
			if err != nil || vec == nil || len(vec.GetId()) == 0 {
				shouldInsert = true
			} else if conv.F32stos(vec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
				continue
			}
		} else {
			id, err := s.Exists(ctx, &payload.Object_ID{
				Id: uuid,
			})
			shouldInsert = err != nil || id == nil || len(id.GetId()) == 0
		}
		ids = append(ids, uuid)
		if shouldInsert {
			insertReqs = append(insertReqs, &payload.Insert_Request{
				Vector: vec,
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
					Filters:              req.GetConfig().GetFilters(),
					Timestamp:            req.GetConfig().GetTimestamp(),
				},
			})
		} else {
			updateReqs = append(updateReqs, &payload.Update_Request{
				Vector: vec,
				Config: &payload.Update_Config{
					SkipStrictExistCheck: true,
					Filters:              req.GetConfig().GetFilters(),
					Timestamp:            req.GetConfig().GetTimestamp(),
				},
			})
		}
	}

	switch {
	case len(insertReqs) <= 0:
		res, err = s.MultiUpdate(ctx, &payload.Update_MultiRequest{
			Requests: updateReqs,
		})
	case len(updateReqs) <= 0:
		res, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
			Requests: insertReqs,
		})
	default:
		var (
			ures, ires *payload.Object_Locations
			errs       error
			mu         sync.Mutex
		)
		eg, ectx := errgroup.New(ctx)
		eg.Go(safety.RecoverFunc(func() (err error) {
			ures, err = s.MultiUpdate(ectx, &payload.Update_MultiRequest{
				Requests: updateReqs,
			})
			if err != nil {
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				mu.Unlock()
			}
			return nil
		}))
		eg.Go(safety.RecoverFunc(func() (err error) {
			ires, err = s.MultiInsert(ectx, &payload.Insert_MultiRequest{
				Requests: insertReqs,
			})
			if err != nil {
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				mu.Unlock()
			}
			return nil
		}))
		err = eg.Wait()
		if err != nil {
			if errs == nil {
				errs = err
			} else {
				errs = errors.Wrap(errs, err.Error())
			}
		}
		if errs != nil {
			err = errs
		}
		switch {
		case ures.GetLocations() == nil && ires.GetLocations() != nil:
			res = ires
		case ures.GetLocations() != nil && ires.GetLocations() == nil:
			res = ures
		case ures.GetLocations() != nil && ires.GetLocations() != nil:
			res = &payload.Object_Locations{
				Locations: append(ures.GetLocations(), ires.GetLocations()...),
			}
		default:
			res = new(payload.Object_Locations)
		}

	}

	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.MultiUpsertRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return location.ReStructure(ids, res), nil
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (locs *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	id := req.GetId()
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, err := s.Exists(ctx, id)
		if err != nil || id == nil || len(id.GetId()) == 0 {
			if err == nil {
				err = errors.ErrObjectIDNotFound(id.GetId())
			}
			st, msg, err := status.ParseError(err, codes.NotFound,
				"error "+vald.RemoveRPCName+" API ID = "+id.GetId()+" not found",
				&errdetails.RequestInfo{
					RequestId:   id.GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + "." + vald.ExistsRPCName,
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
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(ctx, apiName+"."+vald.RemoveRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.Remove(ctx, req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   id.GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
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
		mu.Lock()
		locs.Ips = append(locs.GetIps(), loc.GetIps()...)
		locs.Name = loc.GetName()
		mu.Unlock()
		return nil
	})
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.RemoveRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   id.GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if len(locs.Ips) <= 0 {
		err = errors.ErrIndexNotFound
		err = status.WrapWithNotFound(vald.RemoveRPCName+" API remove target not found", err,
			&errdetails.RequestInfo{
				RequestId:   id.GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return locs, nil
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

	now := time.Now().UnixNano()
	ids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.GetRequests() {
		id := req.GetId()
		ids = append(ids, id.GetId())
		if !req.GetConfig().GetSkipStrictExistCheck() {
			sid, err := s.Exists(ctx, id)
			if err != nil || sid == nil || len(sid.GetId()) == 0 {
				if err == nil {
					err = errors.ErrObjectIDNotFound(id.GetId())
				}
				st, msg, err := status.ParseError(err, codes.NotFound,
					fmt.Sprintf(vald.MultiRemoveRPCName+" API ID = %v not found", id.GetId()),
					&errdetails.RequestInfo{
						RequestId:   id.GetId(),
						ServingData: errdetails.Serialize(reqs),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName + "." + vald.ExistsRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
					}, info.Get())
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil, err
			}
			if reqs.GetRequests()[i].GetConfig() != nil {
				reqs.GetRequests()[i].GetConfig().SkipStrictExistCheck = true
			} else {
				reqs.GetRequests()[i].Config = &payload.Remove_Config{SkipStrictExistCheck: true}
			}

		}
		if req.GetConfig().GetTimestamp() == 0 {
			if req.GetConfig() == nil {
				reqs.GetRequests()[i].Config = &payload.Remove_Config{
					Timestamp: now,
				}
			} else {
				reqs.GetRequests()[i].GetConfig().Timestamp = now
			}
		}
	}
	var mu sync.Mutex
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, 0, len(reqs.GetRequests())),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(ctx, apiName+"."+vald.MultiRemoveRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.MultiRemove(ctx, reqs, copts...)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled),
				errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.MultiRemoveRPCName + ".BroadCast/" +
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
							"/vald.v1." + vald.MultiRemoveRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil
			}
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse MultiRemove gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   strings.Join(ids, ","),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName,
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
		mu.Lock()
		locs.Locations = append(locs.GetLocations(), loc.GetLocations()...)
		mu.Unlock()
		return nil
	})
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.MultiRemoveRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if len(locs.Locations) <= 0 {
		err = errors.ErrIndexNotFound
		err = status.WrapWithNotFound(vald.MultiRemoveRPCName+" API remove target not found", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return location.ReStructure(ids, locs), nil
}

func (s *server) GetObject(ctx context.Context, req *payload.Object_VectorRequest) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vch := make(chan *payload.Object_Vector, 1)
	ech := make(chan error, 1)
	s.eg.Go(func() error {
		defer close(vch)
		defer close(ech)
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		var once sync.Once
		ech <- s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.GetObjectRPCName+"/"+target)
			defer func() {
				if span != nil {
					sspan.End()
				}
			}()
			ovec, err := vc.GetObject(sctx, req, copts...)
			if err != nil {
				switch {
				case errors.Is(err, context.Canceled),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.StatusCodeCancelled(
							errdetails.ValdGRPCResourceTypePrefix +
								"/vald.v1." + vald.GetObjectRPCName + ".BroadCast/" +
								target + " canceled: " + err.Error())...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
					return nil
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.StatusCodeDeadlineExceeded(
							errdetails.ValdGRPCResourceTypePrefix +
								"/vald.v1." + vald.GetObjectRPCName + ".BroadCast/" +
								target + " deadline_exceeded: " + err.Error())...)
						sspan.SetStatus(trace.StatusError, err.Error())
					}
					return nil
				}
				uuid := req.GetId().GetId()
				st, msg, err := status.ParseError(err, codes.NotFound,
					vald.GetObjectRPCName+" API ID = "+uuid+"'s not found",
					&errdetails.RequestInfo{
						RequestId:   uuid,
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
					}, info.Get())
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				if err != nil && st.Code() != codes.NotFound {
					return err
				}
				return nil
			}
			if ovec != nil && ovec.GetId() != "" && ovec.GetVector() != nil {
				once.Do(func() {
					vch <- ovec
					cancel()
				})
			}
			return nil
		})
		return nil
	})
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case vec = <-vch:
	case err = <-ech:
	}
	if err != nil || vec == nil || vec.GetId() == "" || vec.GetVector() == nil {
		err = errors.ErrObjectNotFound(err, req.GetId().GetId())
		st, msg, err := status.ParseError(err, codes.NotFound,
			"failed to parse "+vald.GetObjectRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetId().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
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
