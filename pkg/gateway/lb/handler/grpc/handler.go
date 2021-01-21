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

const apiName = "vald/gateway-lb"

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
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil
		}
		if oid != nil && oid.Id != "" {
			once.Do(func() {
				id = &payload.Object_ID{
					Id: oid.Id,
				}
				cancel()
			})
		}
		return nil
	})
	if err != nil || id == nil || id.Id == "" {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Exists API meta %s's uuid not found", meta.GetId()), err, meta.GetId(), info.Get())
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("Search API invalid vector argument", err, req, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       "vector dimension size",
					Description: err.Error(),
				},
			},
		}, info.Get())
	}
	res, err = s.search(ctx, req.GetConfig(),
		func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
			return vc.Search(ctx, req, copts...)
		})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal("Search API failed to process search request", err, req, info.Get())
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("SearchByID API invalid uuid", err, req, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       "invalid id",
					Description: err.Error(),
				},
			},
		}, info.Get())
	}
	vec, err := s.GetObject(ctx, &payload.Object_ID{
		Id: req.GetId(),
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("SearchByID API uuid %s's object not found", req.GetId()), err, info.Get())
	}
	res, err = s.Search(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: req.GetConfig(),
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal("SearchByID API failed to process search request", err, req, info.Get())
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
				log.Debug(err)
				if span != nil {
					span.SetStatus(trace.StatusCodeInternal(err.Error()))
				}
				return nil
			}
			if r == nil || len(r.GetResults()) == 0 {
				err = errors.ErrIndexNotFound
				log.Debug(err)
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
				log.Error(err)
			}
			close(dch)
			if num != 0 && len(res.GetResults()) > num {
				res.Results = res.Results[:num]
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

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
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
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			res, err := s.SearchByID(ctx, data.(*payload.Search_IDRequest))
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.Search(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiSearch API vector %v's search request result not found",
							query.GetVector()), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiSearch API vector %v's search request result not found",
								query.GetVector()), err, info.Get()).Error())
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		})
	}
	wg.Wait()
	return res, errs
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.SearchByID(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiSearchByID API uuid %v's search by id request result not found",
							query.GetId()), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiSearchByID API uuid %v's search by id request result not found",
								query.GetId()), err, info.Get()).Error())
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		})
	}
	wg.Wait()
	return res, errs
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("Search API invalid vector argument", err, req, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       "vector dimension size",
					Description: err.Error(),
				},
			},
		}, info.Get())
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, err := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if err == nil && id != nil && len(id.GetId()) != 0 {
			err = errors.ErrMetaDataAlreadyExists(uuid)
			log.Error(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, status.WrapWithAlreadyExists(
				fmt.Sprintf("Insert API ID = %v already exists", uuid), err, info.Get())
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
			if err == errors.ErrRPCCallFailed(target, context.Canceled) {
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
		err = errors.Wrapf(err, "Insert API (do multiple) failed to Insert uuid = %s\t info = %#v", uuid, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API failed to Execute DoMulti error = %s", err.Error()), err, info.Get())
	}
	log.Debugf("Insert API insert succeeded to %#v", ce)
	return ce, nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Insert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			res, err := s.Insert(ctx, data.(*payload.Insert_Request))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
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
	for _, vec := range vecs {
		uuid := vec.GetVector().GetId()
		if !vec.GetConfig().GetSkipStrictExistCheck() {
			id, err := s.Exists(ctx, &payload.Object_ID{
				Id: uuid,
			})
			if err == nil && id != nil && len(id.GetId()) != 0 {
				err = errors.ErrMetaDataAlreadyExists(uuid)
				log.Error(err)
				if span != nil {
					span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
				}
				return nil, status.WrapWithAlreadyExists(
					fmt.Sprintf("MultiInsert API ID = %v already exists", uuid), err, info.Get())
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiInsert API failed request %#v", vecs), err, info.Get())
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

	res, err = s.Remove(ctx, &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: req.GetConfig().GetSkipStrictExistCheck(),
		},
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
		}
		return nil, err
	}

	res, err = s.Insert(ctx, &payload.Insert_Request{
		Vector: req.GetVector(),
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: true,
			Filters:              req.GetConfig().Filters,
		},
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Update API failed to remove exsisting data for update %#v", req), err, info.Get())
	}
	return res, nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Update_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			res, err := s.Update(ctx, data.(*payload.Update_Request))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
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
				SkipStrictExistCheck: vec.GetConfig().GetSkipStrictExistCheck(),
			},
		})
	}
	locs, err := s.MultiRemove(ctx, &payload.Remove_MultiRequest{
		Requests: rreqs,
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed Remove request %#v", ids), err, info.Get())
	}
	log.Debugf("uuids %v were removed from %v for MultiUpdate it will execute MultiInsert soon, see detailt %#v", ids, locs.GetLocations(), locs)
	locs, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
		Requests: ireqs,
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed Insert request %#v", vecs), err, info.Get())
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
	filters := req.GetConfig().GetFilters()
	_, err = s.Exists(ctx, &payload.Object_ID{
		Id: uuid,
	})
	if err != nil {
		log.Debugf("Upsert API metadata exists check to Agent error:\t%s", err.Error())
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
		log.Debugf("Upsert API failed to process request uuid:\t%s\terror:\t%s", uuid, err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Upsert API failed to process request %s", uuid), err, info.Get())
	}
	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Upsert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			res, err := s.Upsert(ctx, data.(*payload.Upsert_Request))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
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
		log.Debugf("MultiUpsert API failed to process request uuids:\t%s\terror:\t%s", ids, err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpsert API failed to process request %v", ids), err, info.Get())
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
		sid, err := s.Exists(ctx, id)
		if err != nil || sid == nil || len(sid.GetId()) == 0 {
			err = errors.ErrObjectNotFound(err, id.GetId())
			log.Error(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, status.WrapWithNotFound(
				fmt.Sprintf("Remove API ID = %v not found", id), err, info.Get())
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
			log.Debug(err)
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
		log.Debugf("Remove API failed to remove uuid:\t%s\terror:\t%s", id.GetId(), err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed request uuid %s", id.GetId()), err, info.Get())
	}
	return locs, nil
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Remove_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			res, err := s.Remove(ctx, data.(*payload.Remove_Request))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
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
}

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		id := req.GetId()
		ids = append(ids, id.GetId())
		if !req.GetConfig().GetSkipStrictExistCheck() {
			sid, err := s.Exists(ctx, id)
			if err != nil || sid == nil || len(sid.GetId()) == 0 {
				err = errors.ErrObjectNotFound(err, id.GetId())
				log.Error(err)
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				return nil, status.WrapWithNotFound(
					fmt.Sprintf("MultiRemove API ID = %v not found", id.GetId()), err, info.Get())
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
			log.Debug(err)
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
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiRemove API failed to request uuids %v ", ids), err, info.Get())
	}
	return location.ReStructure(ids, locs), nil
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (vec *payload.Object_Vector, err error) {
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
		ovec, err := vc.GetObject(ctx, id, copts...)
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
		err = errors.ErrObjectNotFound(err, id.GetId())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("GetObject API uuid %s's object not found", vec.GetId()), err, info.Get())
	}
	return vec, nil
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamGetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			res, err := s.GetObject(ctx, data.(*payload.Object_ID))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
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
}
