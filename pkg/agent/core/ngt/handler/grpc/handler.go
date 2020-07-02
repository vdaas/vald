//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

	agent "github.com/vdaas/vald/apis/grpc/agent/core"
	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
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
	vald.ValdServer
}

type server struct {
	name              string
	ip                string
	ngt               service.NGT
	streamConcurrency int
}

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) newLocation(uuids ...string) (locs *payload.Object_Locations) {
	locs = new(payload.Object_Locations)
	for _, uuid := range uuids {
		locs.Locations = append(locs.Locations, &payload.Object_Location{
			Name: s.name,
			Uuid: uuid,
			Ips:  []string{s.ip},
		})
	}
	return locs
}

func (s *server) Exists(ctx context.Context, uid *payload.Object_ID) (res *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.Exists")
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
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.Search")
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
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.SearchByID")
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

func (s *server) StreamSearch(stream vald.Vald_StreamSearchServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/agent-ngt.StreamSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Search(ctx, data.(*payload.Search_Request))
		})
}

func (s *server) StreamSearchByID(stream vald.Vald_StreamSearchByIDServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/agent-ngt.StreamSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.SearchByID(ctx, data.(*payload.Search_IDRequest))
		})
}

func (s *server) Insert(ctx context.Context, vec *payload.Object_Vector) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.Insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = s.ngt.Insert(vec.GetId(), vec.GetVector())
	if err != nil {
		log.Errorf("[Insert]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API failed to insert %#v", vec), err, info.Get())
	}
	return s.newLocation(vec.GetId()).Locations[0], nil
}

func (s *server) StreamInsert(stream vald.Vald_StreamInsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/agent-ngt.StreamInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Insert(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiInsert(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuids := make([]string, 0, len(vecs.Vectors))
	vmap := make(map[string][]float32, len(vecs.GetVectors()))
	for _, vec := range vecs.GetVectors() {
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
	return s.newLocation(uuids...), nil
}

func (s *server) Update(ctx context.Context, vec *payload.Object_Vector) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.Update")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = s.ngt.Update(vec.GetId(), vec.GetVector())
	if err != nil {
		log.Errorf("[Update]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Update API failed to update %#v", vec), err, info.Get())
	}
	return s.newLocation(vec.GetId()).Locations[0], nil
}

func (s *server) StreamUpdate(stream vald.Vald_StreamUpdateServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/agent-ngt.StreamUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Update(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiUpdate(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.MultiUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	uuids := make([]string, 0, len(vecs.Vectors))
	vmap := make(map[string][]float32, len(vecs.GetVectors()))
	for _, vec := range vecs.GetVectors() {
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
	return s.newLocation(uuids...), nil
}

func (s *server) Upsert(ctx context.Context, vec *payload.Object_Vector) (*payload.Object_Location, error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, exists := s.ngt.Exists(vec.GetId())
	if exists {
		return s.Update(ctx, vec)
	}
	return s.Insert(ctx, vec)
}

func (s *server) StreamUpsert(stream vald.Vald_StreamUpsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/agent-ngt.StreamUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Upsert(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiUpsert(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.MultiUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	insertVecs := make([]*payload.Object_Vector, 0, len(vecs.GetVectors()))
	updateVecs := make([]*payload.Object_Vector, 0, len(vecs.GetVectors()))

	for _, vec := range vecs.GetVectors() {
		_, exists := s.ngt.Exists(vec.GetId())
		if exists {
			updateVecs = append(updateVecs, vec)
		} else {
			insertVecs = append(insertVecs, vec)
		}
	}

	var ures, ires *payload.Object_Locations

	eg, ectx := errgroup.New(ctx)
	eg.Go(safety.RecoverFunc(func() error {
		var err error
		if len(updateVecs) > 0 {
			ures, err = s.MultiUpdate(ectx, &payload.Object_Vectors{
				Vectors: updateVecs,
			})
		}
		return err
	}))

	eg.Go(safety.RecoverFunc(func() error {
		var err error
		if len(insertVecs) > 0 {
			ires, err = s.MultiInsert(ectx, &payload.Object_Vectors{
				Vectors: insertVecs,
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

func (s *server) Remove(ctx context.Context, id *payload.Object_ID) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := id.GetId()
	err = s.ngt.Delete(uuid)
	if err != nil {
		log.Errorf("[Remove]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed to delete uuid %s", uuid), err, info.Get())
	}
	return s.newLocation(uuid).Locations[0], nil
}

func (s *server) StreamRemove(stream vald.Vald_StreamRemoveServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/agent-ngt.StreamRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Remove(ctx, data.(*payload.Object_ID))
		})
}

func (s *server) MultiRemove(ctx context.Context, ids *payload.Object_IDs) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuids := ids.GetIds()
	err = s.ngt.DeleteMultiple(uuids...)
	if err != nil {
		log.Errorf("[MultiRemove]\tUnknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed to delete %#v", uuids), err, info.Get())
	}
	return s.newLocation(uuids...), nil
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.GetObject")
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

func (s *server) StreamGetObject(stream vald.Vald_StreamGetObjectServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/agent-ngt.StreamGetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.GetObject(ctx, data.(*payload.Object_ID))
		})
}

func (s *server) CreateIndex(ctx context.Context, c *payload.Control_CreateIndexRequest) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.CreateIndex")
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
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.SaveIndex")
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
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.CreateAndSaveIndex")
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
	ctx, span := trace.StartSpan(ctx, "vald/agent-ngt.IndexInfo")
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
