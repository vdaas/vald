//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/agent/ngt/model"
	"github.com/vdaas/vald/pkg/agent/ngt/service"
)

type Server agent.AgentServer

type server struct {
	ngt service.NGT
}

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Exists(ctx context.Context, oid *payload.Object_ID) (res *payload.Object_ID, err error) {
	res = new(payload.Object_ID)
	var ok bool
	rid, ok := s.ngt.Exists(oid.GetId())
	if !ok {
		err = errors.ErrObjectIDNotFound(oid.GetId())
		return nil, err
	}
	res.Id = strconv.Itoa(int(rid))
	return res, nil
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (*payload.Search_Response, error) {
	return toSearchResponse(
		s.ngt.Search(
			req.GetVector().GetVector(),
			req.GetConfig().GetNum(),
			req.GetConfig().GetEpsilon(),
			req.GetConfig().GetRadius()))
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_Response, error) {
	return toSearchResponse(
		s.ngt.SearchByID(
			req.GetId().GetId(),
			req.GetConfig().GetNum(),
			req.GetConfig().GetEpsilon(),
			req.GetConfig().GetRadius()))
}

func toSearchResponse(dists []model.Distance, err error) (res *payload.Search_Response, rerr error) {
	res = new(payload.Search_Response)
	if err != nil {
		return nil, err
	}
	res.Results = make([]*payload.Object_Distance, 0, len(dists))
	for _, dist := range dists {
		res.Results = append(res.Results, &payload.Object_Distance{
			Id: &payload.Object_ID{
				Id: dist.ID,
			},
			Distance: dist.Distance,
		})
	}
	return res, nil
}

func (s *server) StreamSearch(stream agent.Agent_StreamSearchServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Search(ctx, data.(*payload.Search_Request))
		})
}

func (s *server) StreamSearchByID(stream agent.Agent_StreamSearchByIDServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.SearchByID(ctx, data.(*payload.Search_IDRequest))
		})
}

func (s *server) Insert(ctx context.Context, vec *payload.Object_Vector) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	err = s.ngt.Insert(vec.GetId().GetId(), vec.GetVector())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *server) StreamInsert(stream agent.Agent_StreamInsertServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Insert(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiInsert(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	for _, vec := range vecs.GetVectors() {
		_, ierr := s.Insert(ctx, vec)
		if ierr != nil {
			err = errors.Wrap(err, ierr.Error())
		}
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *server) Update(ctx context.Context, vec *payload.Object_Vector) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	err = s.ngt.Update(vec.GetId().GetId(), vec.GetVector())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *server) StreamUpdate(stream agent.Agent_StreamUpdateServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Update(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiUpdate(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	for _, vec := range vecs.GetVectors() {
		_, ierr := s.Update(ctx, vec)
		if ierr != nil {
			err = errors.Wrap(err, ierr.Error())
		}
	}
	if err != nil {
		return nil, err
	}
	return res, err
}

func (s *server) Remove(ctx context.Context, id *payload.Object_ID) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	err = s.ngt.Delete(id.GetId())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *server) StreamRemove(stream agent.Agent_StreamRemoveServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Remove(ctx, data.(*payload.Object_ID))
		})
}

func (s *server) MultiRemove(ctx context.Context, ids *payload.Object_IDs) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	for _, id := range ids.GetIds() {
		_, ierr := s.Remove(ctx, id)
		if ierr != nil {
			err = errors.Wrap(err, ierr.Error())
		}
	}
	if err != nil {
		return nil, err
	}
	return res, err
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (res *payload.Object_Vector, err error) {
	res = new(payload.Object_Vector)
	res.Id = &payload.Object_ID{
		Id: id.GetId(),
	}
	res.Vector, err = s.ngt.GetObject(id.GetId())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *server) StreamGetObject(stream agent.Agent_StreamGetObjectServer) error {
	return grpc.BidirectionalStream(stream,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.GetObject(ctx, data.(*payload.Object_ID))
		})
}

func (s *server) CreateIndex(ctx context.Context, c *payload.Controll_CreateIndexRequest) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	err = s.ngt.CreateIndex(c.GetPoolSize())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *server) SaveIndex(context.Context, *payload.Empty) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	err = s.ngt.SaveIndex()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *server) CreateAndSaveIndex(ctx context.Context, c *payload.Controll_CreateIndexRequest) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	err = s.ngt.CreateAndSaveIndex(c.GetPoolSize())
	if err != nil {
		return nil, err
	}
	return res, nil
}
