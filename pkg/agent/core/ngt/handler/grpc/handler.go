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
	"strconv"
	"sync"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

type Server interface {
	agent.AgentServer
	vald.Server
}

type server struct {
	name              string
	ip                string
	ngt               service.NGT
	eg                errgroup.Group
	streamConcurrency int
}

const apiName = "vald/agent-ngt"

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) newLocations(uuids ...string) (locs *payload.Object_Locations) {
	if len(uuids) == 0 {
		return nil
	}
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, 0, len(uuids)),
	}
	for _, uuid := range uuids {
		locs.Locations = append(locs.Locations, &payload.Object_Location{
			Name: s.name,
			Uuid: uuid,
			Ips:  []string{s.ip},
		})
	}
	return locs
}

func (s *server) newLocation(uuid string) *payload.Object_Location {
	locs := s.newLocations(uuid)
	if locs != nil && locs.Locations != nil && len(locs.Locations) > 0 {
		return locs.Locations[0]
	}
	return nil
}

func (s *server) Exists(ctx context.Context, uid *payload.Object_ID) (res *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Exists")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := uid.GetId()
	oid, ok := s.ngt.Exists(uuid)
	if !ok {
		err = errors.ErrObjectIDNotFound(uuid)
		log.Warn("[Exists] an error occurred:", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Exists API uuid %s's oid not found", uuid), err, info.Get())
	}
	return &payload.Object_ID{
		Id: strconv.Itoa(int(oid)),
	}, nil
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (*payload.Search_Response, error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return toSearchResponse(
		s.ngt.Search(
			req.GetVector(),
			req.GetConfig().GetNum(),
			req.GetConfig().GetEpsilon(),
			req.GetConfig().GetRadius()))
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_Response, error) {
	ctx, span := trace.StartSpan(ctx, apiName+".SearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return toSearchResponse(
		s.ngt.SearchByID(
			req.GetId(),
			req.GetConfig().GetNum(),
			req.GetConfig().GetEpsilon(),
			req.GetConfig().GetRadius()))
}

func toSearchResponse(dists []model.Distance, err error) (res *payload.Search_Response, rerr error) {
	res = new(payload.Search_Response)
	if err != nil {
		log.Errorf("[toSearchResponse]\tUnknown error\t%+v", err)
		err = status.WrapWithInternal("Search API error occurred", err, info.Get())
	}
	res.Results = make([]*payload.Object_Distance, 0, len(dists))
	for _, dist := range dists {
		res.Results = append(res.Results, &payload.Object_Distance{
			Id:       dist.ID,
			Distance: dist.Distance,
		})
	}
	return res, err
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
				return &payload.Search_StreamResponse{
					Payload: &payload.Search_StreamResponse_Error{
						Error: status.FromError(err),
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
				return &payload.Search_StreamResponse{
					Payload: &payload.Search_StreamResponse_Error{
						Error: status.FromError(err),
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
				errs = errors.Wrap(errs, status.WrapWithNotFound(fmt.Sprintf("MultiSearch API vector %v's search request result not found", query.GetVector()), err, info.Get()).Error())
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
				errs = errors.Wrap(errs, status.WrapWithNotFound(fmt.Sprintf("MultiSearchByID API uuid %v's search by id request result not found", query.GetId()), err, info.Get()).Error())
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

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	err = s.ngt.Insert(vec.GetId(), vec.GetVector())
	if err != nil {
		log.Errorf("[Insert]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API failed to insert %#v", vec), err, info.Get())
	}
	return s.newLocation(vec.GetId()), nil
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
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Error{
						Error: status.FromError(err),
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

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuids := make([]string, 0, len(reqs.GetRequests()))
	vmap := make(map[string][]float32, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		vec := req.GetVector()
		vmap[vec.GetId()] = vec.GetVector()
		uuids = append(uuids, vec.GetId())
	}
	err = s.ngt.InsertMultiple(vmap)
	if err != nil {
		log.Errorf("[MultiInsert]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiInsert API failed insert %#v", vmap), err, info.Get())
	}
	return s.newLocations(uuids...), nil
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Update")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	err = s.ngt.Update(vec.GetId(), vec.GetVector())
	if err != nil {
		log.Errorf("[Update]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Update API failed to update %#v", vec), err, info.Get())
	}
	return s.newLocation(vec.GetId()), nil
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
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Error{
						Error: status.FromError(err),
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

	uuids := make([]string, 0, len(reqs.GetRequests()))
	vmap := make(map[string][]float32, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		vec := req.GetVector()
		vmap[vec.GetId()] = vec.GetVector()
		uuids = append(uuids, vec.GetId())
	}

	err = s.ngt.UpdateMultiple(vmap)
	if err != nil {
		log.Errorf("[MultiUpdate]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed to update %#v", vmap), err, info.Get())
	}
	return s.newLocations(uuids...), nil
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (*payload.Object_Location, error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, exists := s.ngt.Exists(req.GetVector().GetId())
	if exists {
		return s.Update(ctx, &payload.Update_Request{
			Vector: req.GetVector(),
		})
	}
	return s.Insert(ctx, &payload.Insert_Request{
		Vector: req.GetVector(),
	})
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
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Error{
						Error: status.FromError(err),
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

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	insertReqs := make([]*payload.Insert_Request, 0, len(reqs.GetRequests()))
	updateReqs := make([]*payload.Update_Request, 0, len(reqs.GetRequests()))

	for _, req := range reqs.GetRequests() {
		vec := req.GetVector()
		_, exists := s.ngt.Exists(vec.GetId())
		if exists {
			updateReqs = append(updateReqs, &payload.Update_Request{
				Vector: vec,
			})
		} else {
			insertReqs = append(insertReqs, &payload.Insert_Request{
				Vector: vec,
			})
		}
	}

	var ures, ires *payload.Object_Locations

	eg, ectx := errgroup.New(ctx)
	eg.Go(safety.RecoverFunc(func() error {
		var err error
		if len(updateReqs) > 0 {
			ures, err = s.MultiUpdate(ectx, &payload.Update_MultiRequest{
				Requests: updateReqs,
			})
		}
		return err
	}))

	eg.Go(safety.RecoverFunc(func() error {
		var err error
		if len(insertReqs) > 0 {
			ires, err = s.MultiInsert(ectx, &payload.Insert_MultiRequest{
				Requests: insertReqs,
			})
		}
		return err
	}))

	if err = eg.Wait(); err != nil {
		return nil, status.WrapWithInternal("MultiUpsert API failed", err, info.Get())
	}

	return &payload.Object_Locations{
		Locations: append(ures.Locations, ires.Locations...),
	}, nil
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	id := req.GetId()
	uuid := id.GetId()
	err = s.ngt.Delete(uuid)
	if err != nil {
		log.Errorf("[Remove]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed to delete uuid %s", uuid), err, info.Get())
	}
	return s.newLocation(uuid), nil
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
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Error{
						Error: status.FromError(err),
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

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		uuids = append(uuids, req.GetId().GetId())
	}
	err = s.ngt.DeleteMultiple(uuids...)
	if err != nil {
		log.Errorf("[MultiRemove]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed to delete %#v", uuids), err, info.Get())
	}
	return s.newLocations(uuids...), nil
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".GetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := id.GetId()
	vec, err := s.ngt.GetObject(uuid)
	if err != nil {
		log.Warnf("[GetObject]\tUUID not found\t%v", uuid)
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("GetObject API uuid %s Object not found", uuid), err, info.Get())
	}
	return &payload.Object_Vector{
		Id:     uuid,
		Vector: vec,
	}, nil
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
				return &payload.Object_StreamVector{
					Payload: &payload.Object_StreamVector_Error{
						Error: status.FromError(err),
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

func (s *server) CreateIndex(ctx context.Context, c *payload.Control_CreateIndexRequest) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".CreateIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res = new(payload.Empty)
	err = s.ngt.CreateIndex(ctx, c.GetPoolSize())
	if err != nil {
		if err == errors.ErrUncommittedIndexNotFound {
			log.Warnf("[CreateIndex]\tfailed precondition error\t%s", err.Error())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, status.WrapWithFailedPrecondition(fmt.Sprintf("CreateIndex API failed: %s", err), err)
		}

		log.Errorf("[CreateIndex]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("CreateIndex API failed to create indexes pool_size = %d", c.GetPoolSize()), err, info.Get())
	}
	return res, nil
}

func (s *server) SaveIndex(ctx context.Context, _ *payload.Empty) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".SaveIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res = new(payload.Empty)
	err = s.ngt.SaveIndex(ctx)
	if err != nil {
		log.Errorf("[SaveIndex]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal("SaveIndex API failed to save indexes ", err, info.Get())
	}
	return res, nil
}

func (s *server) CreateAndSaveIndex(ctx context.Context, c *payload.Control_CreateIndexRequest) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".CreateAndSaveIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res = new(payload.Empty)
	err = s.ngt.CreateAndSaveIndex(ctx, c.GetPoolSize())
	if err != nil {
		log.Errorf("[CreateAndSaveIndex]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("CreateAndSaveIndex API failed to create and save indexes pool_size = %d", c.GetPoolSize()), err, info.Get())
	}
	return res, nil
}

func (s *server) IndexInfo(ctx context.Context, _ *payload.Empty) (res *payload.Info_Index_Count, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".IndexInfo")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return &payload.Info_Index_Count{
		Stored:      uint32(s.ngt.Len()),
		Uncommitted: uint32(s.ngt.InsertVCacheLen()),
		Indexing:    s.ngt.IsIndexing(),
	}, nil
}
