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
	"sync"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/filter/egress"
	"github.com/vdaas/vald/internal/client/v1/client/filter/ingress"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

type server struct {
	eg                errgroup.Group
	defaultVectorizer string
	defaultFilters    []string
	ingress           ingress.Client
	egress            egress.Client
	gateway           client.Client
	copts             []grpc.CallOption
	streamConcurrency int
	Vectorizer        string
	DistanceFilters   []string
	ObjectFilters     []string
	SearchFilters     []string
	InsertFilters     []string
	UpdateFilters     []string
	UpsertFilters     []string
	vald.UnimplementedValdServerWithFilter
}

const apiName = "vald/gateway-filter"

func New(opts ...Option) vald.ServerWithFilter {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) SearchObject(ctx context.Context, req *payload.Search_ObjectRequest) (*payload.Search_Response, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.SearchObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vr := req.GetVectorizer()
	if vr == nil || vr.GetPort() == 0 {
		return nil, status.WrapWithInvalidArgument("SearchObject API vectorizer configuration is invalid", errors.ErrFilterNotFound, info.Get())
	}
	if vr.GetHost() == "" {
		vr.Host = "localhost"
	}
	target := fmt.Sprintf("%s:%d", vr.GetHost(), vr.GetPort())
	if len(target) == 0 {
		if len(s.Vectorizer) == 0 {
			return nil, status.WrapWithInvalidArgument("SearchObject API vectorizer configuration is invalid", errors.ErrFilterNotFound, info.Get())
		}
		target = s.Vectorizer
	}
	c, err := s.ingress.Target(ctx, target)
	if err != nil {
		return nil, status.WrapWithUnavailable("SearchObject API target filter API unavailable", err, req, info.Get())
	}
	vec, err := c.GenVector(ctx, &payload.Object_Blob{
		Object: req.GetObject(),
	})
	if err != nil {
		return nil, status.WrapWithInternal("SearchObject API failed to extract vector from filter", err, req, info.Get())
	}
	return s.Search(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: req.GetConfig(),
	})
}

func (s *server) MultiSearchObject(ctx context.Context, reqs *payload.Search_MultiObjectRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiSearchObjectRPCName)
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.SearchObject(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiSearchObject API object %s's search request result not found",
							string(query.GetObject())), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiSearchObject API object %s's search request result not found",
								string(query.GetObject())), err, info.Get()).Error())
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

func (s *server) StreamSearchObject(stream vald.Filter_StreamSearchObjectServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamSearchObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_ObjectRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			res, err := s.SearchObject(ctx, data.(*payload.Search_ObjectRequest))
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

func (s *server) LinearSearchObject(ctx context.Context, req *payload.Search_ObjectRequest) (*payload.Search_Response, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.LinearSearchObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vr := req.GetVectorizer()
	if vr == nil || vr.GetPort() == 0 {
		return nil, status.WrapWithInvalidArgument("LinearSearchObject API vectorizer configuration is invalid", errors.ErrFilterNotFound, info.Get())
	}
	if vr.GetHost() == "" {
		vr.Host = "localhost"
	}
	target := fmt.Sprintf("%s:%d", vr.GetHost(), vr.GetPort())
	if len(target) == 0 {
		if len(s.Vectorizer) == 0 {
			return nil, status.WrapWithInvalidArgument("LinearSearchObject API vectorizer configuration is invalid", errors.ErrFilterNotFound, info.Get())
		}
		target = s.Vectorizer
	}
	c, err := s.ingress.Target(ctx, target)
	if err != nil {
		return nil, status.WrapWithUnavailable("LinearSearchObject API target filter API unavailable", err, req, info.Get())
	}
	vec, err := c.GenVector(ctx, &payload.Object_Blob{
		Object: req.GetObject(),
	})
	if err != nil {
		return nil, status.WrapWithInternal("LinearSearchObject API failed to extract vector from filter", err, req, info.Get())
	}
	return s.LinearSearch(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: req.GetConfig(),
	})
}

func (s *server) MultiLinearSearchObject(ctx context.Context, reqs *payload.Search_MultiObjectRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiLinearSearchObjectRPCName)
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.LinearSearchObject(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiLinearSearchObject API object %s's search request result not found",
							string(query.GetObject())), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiLinearSearchObject API object %s's search request result not found",
								string(query.GetObject())), err, info.Get()).Error())
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

func (s *server) StreamLinearSearchObject(stream vald.Filter_StreamSearchObjectServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamLinearSearchObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_ObjectRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			res, err := s.LinearSearchObject(ctx, data.(*payload.Search_ObjectRequest))
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

func (s *server) InsertObject(ctx context.Context, req *payload.Insert_ObjectRequest) (*payload.Object_Location, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.InsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vr := req.GetVectorizer()
	if vr == nil || vr.GetPort() == 0 {
		return nil, status.WrapWithInvalidArgument("InsertObject API vectorizer configuration is invalid", errors.ErrFilterNotFound, info.Get())
	}
	if vr.GetHost() == "" {
		vr.Host = "localhost"
	}
	target := fmt.Sprintf("%s:%d", vr.GetHost(), vr.GetPort())
	if len(target) == 0 {
		if len(s.Vectorizer) == 0 {
			return nil, status.WrapWithInvalidArgument("InsertObject API vectorizer configuration is invalid", errors.ErrFilterNotFound, info.Get())
		}
		target = s.Vectorizer
	}
	c, err := s.ingress.Target(ctx, target)
	if err != nil {
		return nil, status.WrapWithUnavailable("InsertObject API target filter API unavailable", err, req, info.Get())
	}
	vec, err := c.GenVector(ctx, req.GetObject())
	if err != nil {
		return nil, status.WrapWithInternal("InsertObject API failed to extract vector from filter", err, req, info.Get())
	}
	return s.Insert(ctx, &payload.Insert_Request{
		Vector: &payload.Object_Vector{
			Vector: vec.GetVector(),
			Id:     req.GetObject().GetId(),
		},
		Config: req.GetConfig(),
	})
}

func (s *server) StreamInsertObject(stream vald.Filter_StreamInsertObjectServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamInsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Insert_ObjectRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Insert_ObjectRequest)
			loc, err := s.InsertObject(ctx, req)
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
					Location: loc,
				},
			}, nil
		})
}

func (s *server) MultiInsertObject(ctx context.Context, reqs *payload.Insert_MultiObjectRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiInsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			loc, err := s.InsertObject(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiInsertObject API object id: %s's insert failed",
							query.GetObject().GetId()), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiInsertObject API object id: %s's insert failed",
								query.GetObject().GetId()), err, info.Get()).Error())
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = loc
			return nil
		})
	}
	wg.Wait()
	return locs, errs
}

func (s *server) UpdateObject(ctx context.Context, req *payload.Update_ObjectRequest) (*payload.Object_Location, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.UpdateObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vr := req.GetVectorizer()
	if vr == nil || vr.GetPort() == 0 {
		return nil, status.WrapWithInvalidArgument("UpdateObject API vectorizer configuration is invalid", errors.ErrFilterNotFound, info.Get())
	}
	if vr.GetHost() == "" {
		vr.Host = "localhost"
	}

	target := fmt.Sprintf("%s:%d", vr.GetHost(), vr.GetPort())
	if len(target) == 0 {
		if len(s.Vectorizer) == 0 {
			return nil, status.WrapWithInvalidArgument("UpdateObject API vectorizer configuration is invalid", errors.ErrFilterNotFound, info.Get())
		}
		target = s.Vectorizer
	}
	c, err := s.ingress.Target(ctx, target)
	if err != nil {
		return nil, status.WrapWithUnavailable("UpdateObject API target filter API unavailable", err, req, info.Get())
	}
	vec, err := c.GenVector(ctx, req.GetObject())
	if err != nil {
		return nil, status.WrapWithInternal("UpdateObject API failed to extract vector from filter", err, req, info.Get())
	}
	return s.Update(ctx, &payload.Update_Request{
		Vector: &payload.Object_Vector{
			Vector: vec.GetVector(),
			Id:     req.GetObject().GetId(),
		},
		Config: req.GetConfig(),
	})
}

func (s *server) StreamUpdateObject(stream vald.Filter_StreamUpdateObjectServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamUpdateObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Update_ObjectRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			loc, err := s.UpdateObject(ctx, data.(*payload.Update_ObjectRequest))
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
					Location: loc,
				},
			}, nil
		})
}

func (s *server) MultiUpdateObject(ctx context.Context, reqs *payload.Update_MultiObjectRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiUpdateObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			loc, err := s.UpdateObject(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiUpdateObject API object id: %s's insert failed",
							query.GetObject().GetId()), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiUpdateObject API object id: %s's insert failed",
								query.GetObject().GetId()), err, info.Get()).Error())
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = loc
			return nil
		})
	}
	wg.Wait()
	return locs, errs
}

func (s *server) UpsertObject(ctx context.Context, req *payload.Upsert_ObjectRequest) (*payload.Object_Location, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.UpsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vr := req.GetVectorizer()
	if vr == nil || vr.GetPort() == 0 {
		return nil, status.WrapWithInvalidArgument("UpsertObject API vectorizer configuration is invalid", errors.ErrFilterNotFound, info.Get())
	}
	if vr.GetHost() == "" {
		vr.Host = "localhost"
	}
	target := fmt.Sprintf("%s:%d", vr.GetHost(), vr.GetPort())
	if len(target) == 0 {
		if len(s.Vectorizer) == 0 {
			return nil, status.WrapWithInvalidArgument("UpsertObject API vectorizer configuration is invalid", errors.ErrFilterNotFound, info.Get())
		}
		target = s.Vectorizer
	}
	c, err := s.ingress.Target(ctx, target)
	if err != nil {
		return nil, status.WrapWithUnavailable("UpsertObject API target filter API unavailable", err, req, info.Get())
	}
	vec, err := c.GenVector(ctx, req.GetObject())
	if err != nil {
		return nil, status.WrapWithInternal("UpsertObject API failed to extract vector from filter", err, req, info.Get())
	}
	return s.Upsert(ctx, &payload.Upsert_Request{
		Vector: &payload.Object_Vector{
			Vector: vec.GetVector(),
			Id:     req.GetObject().GetId(),
		},
		Config: req.GetConfig(),
	})
}

func (s *server) StreamUpsertObject(stream vald.Filter_StreamUpsertObjectServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamUpsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Upsert_ObjectRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			loc, err := s.UpsertObject(ctx, data.(*payload.Upsert_ObjectRequest))
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
					Location: loc,
				},
			}, nil
		})
}

func (s *server) MultiUpsertObject(ctx context.Context, reqs *payload.Upsert_MultiObjectRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiUpsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			loc, err := s.UpsertObject(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiUpsertObject API object id: %s's insert failed",
							query.GetObject().GetId()), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiUpsertObject API object id: %s's insert failed",
								query.GetObject().GetId()), err, info.Get()).Error())
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = loc
			return nil
		})
	}
	wg.Wait()
	return locs, errs
}

func (s *server) Exists(ctx context.Context, meta *payload.Object_ID) (*payload.Object_ID, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.ExistsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return s.gateway.Exists(ctx, meta, s.copts...)
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.SearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	targets := req.GetConfig().GetIngressFilters().GetTargets()
	if targets != nil || s.SearchFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.SearchFilters))
		addrs = append(addrs, s.SearchFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.ingress.Target(ctx, addrs...)
		if err != nil {
			return nil, status.WrapWithUnavailable(fmt.Sprintf("Search API ingress filter targets %v not found", addrs), err, info.Get())
		}
		vec, err := c.FilterVector(ctx, &payload.Object_Vector{
			Vector: req.GetVector(),
		})
		if err != nil {
			return nil, status.WrapWithInternal(fmt.Sprintf("Search API ingress filter request to %v failure on vec %v", addrs, req.GetVector()), err, info.Get())
		}
		req.Vector = vec.GetVector()
	}
	res, err = s.gateway.Search(ctx, req, s.copts...)
	if err != nil {
		return nil, err
	}
	targets = req.GetConfig().GetEgressFilters().GetTargets()
	if targets != nil || s.DistanceFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.DistanceFilters))
		addrs = append(addrs, s.DistanceFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.egress.Target(ctx, addrs...)
		if err != nil {
			return nil, status.WrapWithUnavailable(fmt.Sprintf("Search API egress filter targets %v not found", addrs), err, info.Get())
		}
		for i, dist := range res.GetResults() {
			d, err := c.FilterDistance(ctx, dist)
			if err != nil {
				return nil, status.WrapWithInternal(fmt.Sprintf("Search API egress filter request to %v failure on id %s", addrs, dist.GetId()), err, info.Get())
			}
			res.Results[i] = d
		}
	}
	return res, nil
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.SearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.gateway.SearchByID(ctx, req, s.copts...)
	if err != nil {
		return nil, err
	}
	targets := req.GetConfig().GetEgressFilters().GetTargets()
	if targets != nil || s.DistanceFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.DistanceFilters))
		addrs = append(addrs, s.DistanceFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.egress.Target(ctx, addrs...)
		if err != nil {
			return nil, status.WrapWithUnavailable(fmt.Sprintf("SearchByID API egress filter targets %v not found", addrs), err, info.Get())
		}
		for i, dist := range res.GetResults() {
			d, err := c.FilterDistance(ctx, dist)
			if err != nil {
				return nil, status.WrapWithInternal(fmt.Sprintf("SearchByID API egress filter request to %v failure on id %s", addrs, dist.GetId()), err, info.Get())
			}
			res.Results[i] = d
		}
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
		log.Error(err)
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
		log.Error(err)
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
						fmt.Sprintf("MultiSearchByID API id %s's search request result not found",
							query.GetId()), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiSearchByID API id %s's search request result not found",
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

func (s *server) LinearSearch(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.LinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	targets := req.GetConfig().GetIngressFilters().GetTargets()
	if targets != nil || s.SearchFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.SearchFilters))
		addrs = append(addrs, s.SearchFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.ingress.Target(ctx, addrs...)
		if err != nil {
			return nil, status.WrapWithUnavailable(fmt.Sprintf("LinearSearch API ingress filter targets %v not found", addrs), err, info.Get())
		}
		vec, err := c.FilterVector(ctx, &payload.Object_Vector{
			Vector: req.GetVector(),
		})
		if err != nil {
			return nil, status.WrapWithInternal(fmt.Sprintf("LinearSearch API ingress filter request to %v failure on vec %v", addrs, req.GetVector()), err, info.Get())
		}
		req.Vector = vec.GetVector()
	}
	res, err = s.gateway.LinearSearch(ctx, req, s.copts...)
	if err != nil {
		return nil, err
	}
	targets = req.GetConfig().GetEgressFilters().GetTargets()
	if targets != nil || s.DistanceFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.DistanceFilters))
		addrs = append(addrs, s.DistanceFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.egress.Target(ctx, addrs...)
		if err != nil {
			return nil, status.WrapWithUnavailable(fmt.Sprintf("LinearSearch API egress filter targets %v not found", addrs), err, info.Get())
		}
		for i, dist := range res.GetResults() {
			d, err := c.FilterDistance(ctx, dist)
			if err != nil {
				return nil, status.WrapWithInternal(fmt.Sprintf("LinearSearch API egress filter request to %v failure on id %s", addrs, dist.GetId()), err, info.Get())
			}
			res.Results[i] = d
		}
	}
	return res, nil
}

func (s *server) LinearSearchByID(ctx context.Context, req *payload.Search_IDRequest) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.LinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.gateway.LinearSearchByID(ctx, req, s.copts...)
	if err != nil {
		return nil, err
	}
	targets := req.GetConfig().GetEgressFilters().GetTargets()
	if targets != nil || s.DistanceFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.DistanceFilters))
		addrs = append(addrs, s.DistanceFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.egress.Target(ctx, addrs...)
		if err != nil {
			return nil, status.WrapWithUnavailable(fmt.Sprintf("LinearSearchByID API egress filter targets %v not found", addrs), err, info.Get())
		}
		for i, dist := range res.GetResults() {
			d, err := c.FilterDistance(ctx, dist)
			if err != nil {
				return nil, status.WrapWithInternal(fmt.Sprintf("LinearSearchByID API egress filter request to %v failure on id %s", addrs, dist.GetId()), err, info.Get())
			}
			res.Results[i] = d
		}
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
			res, err := s.LinearSearch(ctx, data.(*payload.Search_Request))
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
		log.Error(err)
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
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamLinearSearchByID/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearchByID(ctx, req)
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
		log.Error(err)
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.LinearSearch(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiLinearSearch API vector %v's search request result not found",
							query.GetVector()), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiLinearSearch API vector %v's search request result not found",
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.LinearSearchByID(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiLinearSearchByID API id %s's search request result not found",
							query.GetId()), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiLinearSearchByID API id %s's search request result not found",
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

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.InsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	uuid := vec.GetId()
	if len(vec.GetVector()) < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(len(vec.GetVector()), 0)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("Insert API invalid vector argument", err, req, info.Get())
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, _ := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if id != nil || len(id.GetId()) > 0 {
			err = errors.Wrap(err, errors.ErrMetaDataAlreadyExists(uuid).Error())
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, status.WrapWithAlreadyExists(fmt.Sprintf("Insert API ID %s already exists", vec.GetId()), err, info.Get())
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Insert_Config{SkipStrictExistCheck: true}
		}
	}
	targets := req.GetConfig().GetFilters().GetTargets()
	if len(targets) == 0 && len(s.InsertFilters) == 0 {
		return s.gateway.Insert(ctx, req)
	}
	addrs := make([]string, 0, len(targets)+len(s.InsertFilters))
	addrs = append(addrs, s.InsertFilters...)
	for _, target := range targets {
		addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
	}
	c, err := s.ingress.Target(ctx, addrs...)
	if err != nil {
		return nil, status.WrapWithUnavailable(fmt.Sprintf("Insert API ingress filter targets %v not found", addrs), err, info.Get())
	}
	vec, err = c.FilterVector(ctx, req.GetVector())
	if err != nil {
		return nil, status.WrapWithInternal(
			fmt.Sprintf("Insert API ingress filter request to %v failure on id: %s\tvec: %v", addrs, req.GetVector().GetId(), req.GetVector().GetVector()),
			err,
			info.Get(),
		)
	}
	if vec.GetId() == "" {
		vec.Id = req.GetVector().GetId()
	}
	req.Vector = vec
	loc, err = s.gateway.Insert(ctx, req, s.copts...)
	if err != nil {
		err = errors.Wrapf(err, "Insert API failed to Insert uuid = %s\tinfo = %#v", uuid, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API failed to Execute DoMulti error = %s", err.Error()), err, info.Get())
	}
	return loc, nil
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
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.Insert(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiInsert API request %#v's Insert request result not found",
							query), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiInsert API request %#v's Insert request result not found",
								query), err, info.Get()).Error())
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = r
			return nil
		})
	}
	wg.Wait()
	return locs, errs
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	uuid := vec.GetId()
	if len(vec.GetVector()) < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(len(vec.GetVector()), 0)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("Update API invalid vector argument", err, req, info.Get())
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, _ := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if id != nil || len(id.GetId()) > 0 {
			err = errors.Wrap(err, errors.ErrMetaDataAlreadyExists(uuid).Error())
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, status.WrapWithAlreadyExists(fmt.Sprintf("Update API ID %s already exists", vec.GetId()), err, info.Get())
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Update_Config{SkipStrictExistCheck: true}
		}
	}
	targets := req.GetConfig().GetFilters().GetTargets()
	if len(targets) == 0 && len(s.UpdateFilters) == 0 {
		return s.gateway.Update(ctx, req)
	}
	addrs := make([]string, 0, len(targets)+len(s.UpdateFilters))
	addrs = append(addrs, s.UpdateFilters...)
	for _, target := range targets {
		addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
	}
	c, err := s.ingress.Target(ctx, addrs...)
	if err != nil {
		return nil, status.WrapWithUnavailable(fmt.Sprintf("Update API ingress filter targets %v not found", addrs), err, info.Get())
	}
	vec, err = c.FilterVector(ctx, req.GetVector())
	if err != nil {
		return nil, status.WrapWithInternal(
			fmt.Sprintf("Update API ingress filter request to %v failure on id: %s\tvec: %v", addrs, req.GetVector().GetId(), req.GetVector().GetVector()),
			err,
			info.Get(),
		)
	}
	if vec.GetId() == "" {
		vec.Id = req.GetVector().GetId()
	}
	req.Vector = vec
	loc, err = s.gateway.Update(ctx, req, s.copts...)
	if err != nil {
		err = errors.Wrapf(err, "Update API failed to Update uuid = %s\tinfo = %#v", uuid, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Update API failed to Execute DoMulti error = %s", err.Error()), err, info.Get())
	}
	return loc, nil
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
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.Update(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiUpdate API request %#v's Update request result not found",
							query), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiUpdate API request %#v's Update request result not found",
								query), err, info.Get()).Error())
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = r
			return nil
		})
	}
	wg.Wait()
	return locs, errs
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
	if len(vec.GetVector()) < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(len(vec.GetVector()), 0)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("Upsert API invalid vector argument", err, req, info.Get())
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, _ := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if id != nil || len(id.GetId()) > 0 {
			err = errors.Wrap(err, errors.ErrMetaDataAlreadyExists(uuid).Error())
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, status.WrapWithAlreadyExists(fmt.Sprintf("Upsert API ID %s already exists", vec.GetId()), err, info.Get())
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Upsert_Config{SkipStrictExistCheck: true}
		}
	}
	targets := req.GetConfig().GetFilters().GetTargets()
	if len(targets) == 0 && len(s.UpsertFilters) == 0 {
		return s.gateway.Upsert(ctx, req)
	}
	addrs := make([]string, 0, len(targets))
	addrs = append(addrs, s.UpsertFilters...)
	for _, target := range targets {
		addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
	}
	c, err := s.ingress.Target(ctx, addrs...)
	if err != nil {
		return nil, status.WrapWithUnavailable(fmt.Sprintf("Upsert API ingress filter targets %v not found", addrs), err, info.Get())
	}
	vec, err = c.FilterVector(ctx, req.GetVector())
	if err != nil {
		return nil, status.WrapWithInternal(
			fmt.Sprintf("Upsert API ingress filter request to %v failure on id: %s\tvec: %v", addrs, req.GetVector().GetId(), req.GetVector().GetVector()),
			err,
			info.Get(),
		)
	}
	if vec.GetId() == "" {
		vec.Id = req.GetVector().GetId()
	}
	req.Vector = vec
	loc, err = s.gateway.Upsert(ctx, req, s.copts...)
	if err != nil {
		err = errors.Wrapf(err, "Upsert API failed to Upsert uuid = %s\tinfo = %#v", uuid, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Upsert API failed to Execute DoMulti error = %s", err.Error()), err, info.Get())
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
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.Upsert(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiUpsert API request %#v's Upsert request result not found",
							query), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiUpsert API request %#v's Upsert request result not found",
								query), err, info.Get()).Error())
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = r
			return nil
		})
	}
	wg.Wait()
	return locs, errs
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return s.gateway.Remove(ctx, req, s.copts...)
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
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.Remove(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiRemove API id %s's Remove request result not found",
							query.GetId()), err, info.Get())
				} else {
					errs = errors.Wrap(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiRemove API id %s's Remove request result not found",
								query.GetId()), err, info.Get()).Error())
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = r
			return nil
		})
	}
	wg.Wait()
	return locs, errs
}

func (s *server) GetObject(ctx context.Context, req *payload.Object_VectorRequest) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.GetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec, err = s.gateway.GetObject(ctx, req)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("GetObject API uuid %s Object not found", req.GetId().GetId()), err, info.Get())
	}
	targets := req.GetFilters().GetTargets()
	if targets != nil || s.ObjectFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.ObjectFilters))
		addrs = append(addrs, s.ObjectFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.egress.Target(ctx, addrs...)
		if err != nil {
			return nil, status.WrapWithUnavailable(fmt.Sprintf("GetObject API egress filter targets %v not found on id %s", addrs, req.GetId().GetId()), err, info.Get())
		}
		vec, err = c.FilterVector(ctx, vec)
		if err != nil {
			return nil, status.WrapWithInternal(fmt.Sprintf("GetObject API egress filter request to %v failure on id %s", addrs, req.GetId().GetId()), err, info.Get())
		}
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
		log.Error(err)
		return err
	}
	return nil
}
