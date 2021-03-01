//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
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
	"github.com/vdaas/vald/pkg/gateway/internal/location"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type server struct {
	eg                errgroup.Group
	gateway           service.Gateway
	timeout           time.Duration
	replica           int
	streamConcurrency int
}

const apiName = "vald/gateway/lb"

func New(opts ...Option) vald.Server {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Exists(ctx context.Context, meta *payload.Object_ID) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Exists")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	var once sync.Once
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(ctx, apiName+".Exists/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		oid, err := vc.Exists(ctx, meta, copts...)
		if err != nil {
			err = status.WrapWithNotFound(fmt.Sprintf("Exists API meta %s's uuid not found in %s", meta.GetId(), target), err,
				&errdetails.RequestInfo{
					RequestId:   meta.GetId(),
					ServingData: errdetails.Serialize(meta),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
					ResourceName: target,
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				})
// 			 log.Debug(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil
		}
		if oid != nil && oid.GetId() != "" {
			once.Do(func() {
				id = oid
				cancel()
			})
		}
		return nil
	})
	if err != nil || id == nil || id.GetId() == "" {
		if err == nil {
			err = errors.ErrObjectIDNotFound(id.GetId())
		}
		err = status.WrapWithNotFound(fmt.Sprintf("Exists API meta %s's uuid not found", meta.GetId()), err,
			&errdetails.RequestInfo{
				RequestId:   meta.GetId(),
				ServingData: errdetails.Serialize(meta),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Warn(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}
	return id, nil
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vl := len(req.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument("Search API invalid vector argument", err,
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
// 		 log.Warn(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	res, err = s.search(ctx, req.GetConfig(),
		func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.Search(ctx, req, copts...)
		})

	if err != nil {
		err = status.WrapWithInternal("Search API failed to process search request", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Search",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".SearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(req.GetId()) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument("SearchByID API invalid uuid", err,
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
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetId(),
		},
		Filters: req.GetConfig().GetEgressFilters(),
	})
	if err != nil {
		err = status.WrapWithNotFound(fmt.Sprintf("SearchByID API uuid %s's object not found", req.GetId()), err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.SearchByID",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			},
			info.Get())
// 		 log.Error(err)
		var serr error
		res, serr = s.search(ctx, req.GetConfig(),
			func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
				return vc.SearchByID(ctx, req, copts...)
			})

		if serr == nil {
			return res, nil
		}
		err = errors.Wrap(err, status.WrapWithInternal("SearchByID API failed to process search request", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.SearchByID",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}).Error())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}
	res, err = s.Search(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: req.GetConfig(),
	})
	if err != nil {
		err = status.WrapWithInternal("SearchByID API failed to process search request", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.SearchByID",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		var serr error
		res, serr = s.search(ctx, req.GetConfig(),
			func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
				return vc.SearchByID(ctx, req, copts...)
			})

		if serr == nil {
			return res, nil
		}
		err = errors.Wrap(err, status.WrapWithInternal("SearchByID API failed to process search request", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.SearchByID",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}).Error())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) search(ctx context.Context, cfg *payload.Search_Config,
	f func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error)) (
	res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	num := int(cfg.GetNum())
	res = new(payload.Search_Response)
	res.Results = make([]*payload.Object_Distance, 0, s.gateway.GetAgentCount(ctx)*num)
	dch := make(chan *payload.Object_Distance, cap(res.GetResults())/2)
	eg, ectx := errgroup.New(ctx)
	var cancel context.CancelFunc
	var timeout time.Duration
	if to := cfg.GetTimeout(); to != 0 {
		timeout = time.Duration(to)
	} else {
		timeout = s.timeout
	}

	var maxDist uint32
	atomic.StoreUint32(&maxDist, math.Float32bits(math.MaxFloat32))
	ectx, cancel = context.WithTimeout(ectx, timeout)
	eg.Go(safety.RecoverFunc(func() error {
		defer cancel()
		visited := new(sync.Map)
		return s.gateway.BroadCast(ectx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(ctx, apiName+".search/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			r, err := f(ctx, vc, copts...)
			if err != nil {
				if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
					err = status.WrapWithInternal("failed to process search request", err,
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Search",
							ResourceName: target,
							Owner:        errdetails.ValdResourceOwner,
							Description:  err.Error(),
						})
// 					 log.Warn(err)
				}
				if span != nil {
					span.SetStatus(trace.StatusCodeInternal(err.Error()))
				}
				return nil
			}
			if r == nil || len(r.GetResults()) == 0 {
				err = errors.ErrIndexNotFound
				err = status.WrapWithNotFound("failed to process search request", err,
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Search",
						ResourceName: target,
						Owner:        errdetails.ValdResourceOwner,
						Description:  err.Error(),
					})
// 				 log.Debug(err)
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				return nil
			}
			for _, dist := range r.GetResults() {
				if dist == nil {
					continue
				}
				if dist.GetDistance() >= math.Float32frombits(atomic.LoadUint32(&maxDist)) {
					return nil
				}
				if _, already := visited.LoadOrStore(dist.GetId(), struct{}{}); !already {
					select {
					case <-ectx.Done():
						return nil
					case dch <- dist:
					}
				}
			}
			return nil
		})
	}))
	for {
		select {
		case <-ectx.Done():
			err = eg.Wait()
			if err != nil {
// 				 log.Error(err)
			}
			close(dch)
			if num != 0 && len(res.GetResults()) > num {
				res.Results = res.Results[:num]
			}
			if len(res.Results) == 0 {
				err = status.WrapWithNotFound("search result not found", err,
					&errdetails.RequestInfo{
						RequestId:   cfg.GetRequestId(),
						ServingData: errdetails.Serialize(cfg),
					},
					&errdetails.ResourceInfo{
						ResourceType: "vald agent server Search API",
						ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
					}, info.Get())
// 				 log.Warn(err)
				return nil, err
			}
			return res, nil
		case dist := <-dch:
			rl := len(res.GetResults()) // result length
			if rl >= num && dist.GetDistance() >= math.Float32frombits(atomic.LoadUint32(&maxDist)) {
				continue
			}
			switch rl {
			case 0:
				res.Results = append(res.Results, dist)
			case 1:
				if res.GetResults()[0].GetDistance() <= dist.GetDistance() {
					res.Results = append(res.GetResults(), dist)
				} else {
					res.Results = []*payload.Object_Distance{dist, res.GetResults()[0]}
				}
			default:
				pos := rl
				for idx := rl; idx >= 1; idx-- {
					if res.GetResults()[idx-1].GetDistance() <= dist.GetDistance() {
						pos = idx - 1
						break
					}
				}
				switch {
				case pos == rl:
					res.Results = append([]*payload.Object_Distance{dist}, res.Results...)
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
			if distEnd := res.GetResults()[rl-1].GetDistance(); rl >= num &&
				distEnd < math.Float32frombits(atomic.LoadUint32(&maxDist)) {
				atomic.StoreUint32(&maxDist, math.Float32bits(distEnd))
			}
		}
	}
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			res, err := s.Search(ctx, data.(*payload.Search_Request))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
// 		 log.Error(err)
		return err
	}
	return nil
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_IDRequest)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamSearchByID/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.SearchByID(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
// 		 log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.Requests)),
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
			err = status.WrapWithInvalidArgument("MultiSearch API invalid vector argument", err,
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
// 			 log.Warn(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
			}
			return nil, err
		}
		wg.Add(1)
		s.eg.Go(func() error {
			ctx, span := trace.StartSpan(ctx, fmt.Sprintf("%s.MultiSearch/goroutine/id-%d", apiName, idx))
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			defer wg.Done()
			r, err := s.Search(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
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
		err := errs
		err = status.WrapWithInternal("MultiSearch API failed to search", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Search",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) MultiSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.Requests)),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(func() error {
			ctx, span := trace.StartSpan(ctx, fmt.Sprintf("%s.MultiSearchByID/goroutine/id-%d", apiName, idx))
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			defer wg.Done()
			r, err := s.SearchByID(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
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
		err := errs
		err = status.WrapWithInternal("MultiSearchByID API failed to search", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.SearchByID",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (ce *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector().GetVector()
	uuid := req.GetVector().GetId()
	vl := len(vec)
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument("Insert API invalid vector argument", err,
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
// 		 log.Warn(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, err := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if err == nil && id != nil && len(id.GetId()) != 0 {
			err = errors.ErrMetaDataAlreadyExists(uuid)
			err = status.WrapWithAlreadyExists(fmt.Sprintf("Insert API ID = %v already exists", uuid), err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
					ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
// 			 log.Error(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.Config.SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Insert_Config{SkipStrictExistCheck: true}
		}
	}

	mu := new(sync.Mutex)
	ce = &payload.Object_Location{
		Uuid: uuid,
		Ips:  make([]string, 0, s.replica),
	}
	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(ctx, apiName+".Insert/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.Insert(ctx, req, copts...)
		if err != nil {
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			rpcerr := errors.ErrRPCCallFailed(target, context.Canceled)
			if err == rpcerr || errors.Is(err, rpcerr) {
				return nil
			}
			return err
		}
		mu.Lock()
		ce.Ips = append(ce.GetIps(), loc.GetIps()...)
		ce.Name = loc.GetName()
		mu.Unlock()
		return nil
	})
	if err != nil {
		err = status.WrapWithInternal("Insert API failed to Execute DoMulti", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Insert.DoMulti",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
// 	 log.Debugf("Insert API insert succeeded to %#v", ce)
	return ce, nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Insert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Insert_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamInsert/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Insert(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
// 		 log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vecs := reqs.GetRequests()
	ids := make([]string, 0, len(vecs))
	for i, vec := range vecs {
		uuid := vec.GetVector().GetId()
		vector := vec.GetVector().GetVector()
		vl := len(vector)
		if vl < algorithm.MinimumVectorDimensionSize {
			err = errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument("MultiInsert API invalid vector argument", err,
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
// 			 log.Warn(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
			}
			return nil, err
		}
		if !vec.GetConfig().GetSkipStrictExistCheck() {
			id, err := s.Exists(ctx, &payload.Object_ID{
				Id: uuid,
			})
			if err == nil && id != nil && len(id.GetId()) != 0 {
				err = errors.ErrMetaDataAlreadyExists(uuid)
				err = status.WrapWithAlreadyExists(fmt.Sprintf("MultiInsert API ID = %v already exists", uuid), err,
					&errdetails.RequestInfo{
						RequestId:   uuid,
						ServingData: errdetails.Serialize(reqs),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
						ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
						Owner:        errdetails.ValdResourceOwner,
						Description:  err.Error(),
					}, info.Get())
// 				 log.Error(err)
				if span != nil {
					span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
				}
				return nil, err
			}
			if reqs.Requests[i] != nil {
				reqs.Requests[i].Config.SkipStrictExistCheck = true
			} else {
				reqs.Requests[i].Config = &payload.Insert_Config{SkipStrictExistCheck: true}
			}
		}
		ids = append(ids, uuid)
	}

	mu := new(sync.Mutex)
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, 0, s.replica),
	}
	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(ctx, apiName+".MultiInsert/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.MultiInsert(ctx, reqs, copts...)
		if err != nil {
			err = status.WrapWithInternal("MultiInsert API failed to insert", err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(ids, ","),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiInsert",
					ResourceName: target,
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				})
// 			 log.Debug(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return err
		}
		mu.Lock()
		locs.Locations = append(locs.Locations, loc.Locations...)
		mu.Unlock()
		return nil
	})
	if err != nil {
		err = status.WrapWithInternal("MultiInsert API failed to Execute DoMulti", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiInsert.DoMulti",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return location.ReStructure(ids, locs), nil
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Update")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector().GetVector()
	uuid := req.GetVector().GetId()
	vl := len(vec)
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument("Update API invalid vector argument", err,
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
// 		 log.Warn(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}

	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, err := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if err != nil || id == nil || len(id.GetId()) == 0 {
			if err == nil {
				err = errors.ErrObjectIDNotFound(uuid)
			}
			err = status.WrapWithNotFound(fmt.Sprintf("Update API ID = %v not found", uuid), err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
					ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
// 			 log.Error(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.Config.SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Update_Config{SkipStrictExistCheck: true}
		}
	}

	res, err = s.Remove(ctx, &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: uuid,
		},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: true,
		},
	})
	if err != nil {
		err = status.WrapWithInternal("Update API failed to Remove for Insert", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Remove",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}

	res, err = s.Insert(ctx, &payload.Insert_Request{
		Vector: req.GetVector(),
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: true,
			Filters:              req.GetConfig().GetFilters(),
		},
	})
	if err != nil {
		err = status.WrapWithInternal("Update API failed to Insert for Update", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Insert",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Update_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Update_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamUpdate/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Update(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
// 		 log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vecs := reqs.GetRequests()
	ids := make([]string, 0, len(vecs))
	ireqs := make([]*payload.Insert_Request, 0, len(vecs))
	rreqs := make([]*payload.Remove_Request, 0, len(vecs))
	for _, vec := range vecs {
		vl := len(vec.GetVector().GetVector())
		if vl < algorithm.MinimumVectorDimensionSize {
			err = errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument("MultiUpdate API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   vec.GetVector().GetId(),
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
// 			 log.Warn(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
			}
			return nil, err
		}
		uuid := vec.GetVector().GetId()
		if !vec.GetConfig().GetSkipStrictExistCheck() {
			id, err := s.Exists(ctx, &payload.Object_ID{
				Id: uuid,
			})
			if err != nil || id == nil || len(id.GetId()) == 0 {
				if err == nil {
					err = errors.ErrObjectIDNotFound(uuid)
				}
				err = status.WrapWithNotFound(fmt.Sprintf("MultiUpdate API ID = %v not found", uuid), err,
					&errdetails.RequestInfo{
						RequestId:   uuid,
						ServingData: errdetails.Serialize(reqs),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
						ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
						Owner:        errdetails.ValdResourceOwner,
						Description:  err.Error(),
					}, info.Get())
// 				 log.Error(err)
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				return nil, err
			}
			if vec.GetConfig() != nil {
				vec.Config.SkipStrictExistCheck = true
			} else {
				vec.Config = &payload.Update_Config{SkipStrictExistCheck: true}
			}
		}
		ids = append(ids, vec.GetVector().GetId())
		ireqs = append(ireqs, &payload.Insert_Request{
			Vector: vec.GetVector(),
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Filters:              vec.GetConfig().GetFilters(),
			},
		})
		rreqs = append(rreqs, &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: vec.GetVector().GetId(),
			},
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: true,
			},
		})
	}
	locs, err := s.MultiRemove(ctx, &payload.Remove_MultiRequest{
		Requests: rreqs,
	})
	if err != nil {
		err = status.WrapWithInternal("MultiUpdate API failed to Execute MultiRemove for update", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiRemove",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
// 	 log.Debugf("uuids %v were removed from %v for MultiUpdate it will execute MultiInsert soon, see detailt %#v", ids, locs.GetLocations(), locs)
	locs, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
		Requests: ireqs,
	})
	if err != nil {
		err = status.WrapWithInternal("MultiUpdate API failed to Execute MultiInsert for update", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiInsert",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return locs, nil
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	vec := req.GetVector()
	uuid := vec.GetId()
	vl := len(vec.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument("Upsert API invalid vector argument", err,
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
// 		 log.Warn(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	filters := req.GetConfig().GetFilters()
	id, err := s.Exists(ctx, &payload.Object_ID{
		Id: uuid,
	})
	if err != nil || id == nil || len(id.GetId()) == 0 {
		loc, err = s.Insert(ctx, &payload.Insert_Request{
			Vector: vec,
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Filters:              filters,
			},
		})
	} else {
		loc, err = s.Update(ctx, &payload.Update_Request{
			Vector: vec,
			Config: &payload.Update_Config{
				SkipStrictExistCheck: true,
				Filters:              filters,
			},
		})
	}

	if err != nil {
		err = status.WrapWithInternal("Upsert API failed to upsert", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Upsert",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Upsert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Upsert_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamUpsert/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Upsert(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
// 		 log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiUpsert")
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
			err = status.WrapWithInvalidArgument("MultiUpsert API invalid vector argument", err,
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
// 			 log.Warn(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
			}
			return nil, err
		}
		ids = append(ids, uuid)
		_, err = s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		filters := req.GetConfig().GetFilters()
		if err != nil {
			insertReqs = append(insertReqs, &payload.Insert_Request{
				Vector: vec,
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
					Filters:              filters,
				},
			})
		} else {
			updateReqs = append(updateReqs, &payload.Update_Request{
				Vector: vec,
				Config: &payload.Update_Config{
					SkipStrictExistCheck: true,
					Filters:              filters,
				},
			})
		}
	}

	insertLocs := make([]*payload.Object_Location, 0, len(insertReqs))
	updateLocs := make([]*payload.Object_Location, 0, len(updateReqs))

	eg, ectx := errgroup.New(ctx)
	eg.Go(safety.RecoverFunc(func() error {
		if len(updateReqs) <= 0 {
			return nil
		}

		ectx, span := trace.StartSpan(ectx, apiName+".MultiUpsert/Go-MultiUpdate")
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		var err error
		loc, err := s.MultiUpdate(ectx, &payload.Update_MultiRequest{
			Requests: updateReqs,
		})
		if err == nil {
			updateLocs = loc.GetLocations()
		}
		return err
	}))
	eg.Go(safety.RecoverFunc(func() error {
		if len(insertReqs) <= 0 {
			return nil
		}

		ectx, span := trace.StartSpan(ectx, apiName+".MultiUpsert/Go-MultiInsert")
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		var err error
		loc, err := s.MultiInsert(ectx, &payload.Insert_MultiRequest{
			Requests: insertReqs,
		})
		if err == nil {
			insertLocs = loc.GetLocations()
		}
		return err
	}))

	err = eg.Wait()
	if err != nil {
		err = status.WrapWithInternal("MultiUpsert API failed to upsert", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Upsert",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return location.ReStructure(ids, &payload.Object_Locations{
		Locations: append(insertLocs, updateLocs...),
	}), nil
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (locs *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Remove")
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
			err = status.WrapWithNotFound(fmt.Sprintf("Remove API ID = %v not found", id.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   id.GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
					ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
// 			 log.Error(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.Config.SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Remove_Config{SkipStrictExistCheck: true}
		}
	}
	var mu sync.Mutex
	locs = &payload.Object_Location{
		Uuid: id.GetId(),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(ctx, apiName+".Remove/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.Remove(ctx, req, copts...)
		if err != nil {
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil
		}
		mu.Lock()
		locs.Ips = append(locs.Ips, loc.GetIps()...)
		locs.Name = loc.GetName()
		mu.Unlock()
		return nil
	})
	if err != nil {
		err = status.WrapWithInternal("Remove API failed to remove request", err,
			&errdetails.RequestInfo{
				RequestId:   id.GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Remove",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return locs, nil
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Remove_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Remove_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamRemove/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Remove(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
// 		 log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
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
				err = status.WrapWithNotFound(fmt.Sprintf("MultiRemove API ID = %v not found", id.GetId()), err,
					&errdetails.RequestInfo{
						RequestId:   id.GetId(),
						ServingData: errdetails.Serialize(reqs),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
						ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
						Owner:        errdetails.ValdResourceOwner,
						Description:  err.Error(),
					}, info.Get())
// 				 log.Error(err)
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				return nil, err
			}
			if reqs.Requests[i].GetConfig() != nil {
				reqs.Requests[i].Config.SkipStrictExistCheck = true
			} else {
				reqs.Requests[i].Config = &payload.Remove_Config{SkipStrictExistCheck: true}
			}

		}
	}
	var mu sync.Mutex
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, 0, len(reqs.GetRequests())),
	}
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(ctx, apiName+".MultiRemove/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		loc, err := vc.MultiRemove(ctx, reqs, copts...)
		if err != nil {
// 			 log.Debug(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil
		}
		mu.Lock()
		locs.Locations = append(locs.Locations, loc.Locations...)
		mu.Unlock()
		return nil
	})
	if err != nil {
		err = status.WrapWithInternal("MultiRemove API failed to remove", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiRemove",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return location.ReStructure(ids, locs), nil
}

func (s *server) GetObject(ctx context.Context, req *payload.Object_VectorRequest) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".GetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	var once sync.Once
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(ctx, apiName+".GetObject/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		ovec, err := vc.GetObject(ctx, req, copts...)
		if err != nil {
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil
		}
		if ovec != nil && ovec.GetId() != "" && ovec.GetVector() != nil {
			once.Do(func() {
				vec = ovec
				cancel()
			})
		}
		return nil
	})
	if err != nil || vec == nil || vec.GetId() == "" || vec.GetVector() == nil {
		err = errors.ErrObjectNotFound(err, req.GetId().GetId())
		err = status.WrapWithNotFound("GetObject API failed to remove request", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetId().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.GetObject",
				ResourceName: strings.Join(s.gateway.Addrs(ctx), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
// 		 log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}
	return vec, nil
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamGetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_VectorRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Object_VectorRequest)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamGetObject/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.GetObject(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
// 		 log.Error(err)
		return err
	}
	return nil
}
