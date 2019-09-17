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

	"github.com/kpango/fastime"
	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/discoverer/k8s/model"
	"github.com/vdaas/vald/pkg/discoverer/k8s/service"
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

func (s *server) Exists(ctx context.Context, oid *payload.Object_ID) (*payload.Object_ID, error) {
	id, ok := s.ngt.Exists(oid.GetId())
	if !ok {
		return nil, errors.ErrObjectIDNotFound(oid.GetId())
	}
	return &payload.Object_ID{
		Id: id,
	}, nil
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

func toSearchResponse(dists []model.Distance, err error) (*payload.Search_Response, error) {
	if err != nil {
		return &payload.Search_Response{
			Error: &payload.Common_Error{
				Msg:       err.Error(),
				Timestamp: fastime.UnixNanoNow(),
			},
		}, err
	}

	res := &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, len(dists)),
	}

	for _, dist := range dists {
		// res.Results = append(res.Results, (*payload.Object_Distance)(unsafe.Pointer(&dist)))
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
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return s.Search(ctx, data.(*payload.Search_Request))
	})
}

func (s *server) StreamSearchByID(stream agent.Agent_StreamSearchByIDServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return s.SearchByID(ctx, data.(*payload.Search_IDRequest))
	})
}

func (s *server) Insert(ctx context.Context, vec *payload.Object_Vector) (*payload.Common_Error, error) {
	err := s.ngt.Insert(vec.GetId().GetId(), vec.GetVector())
	if err != nil {
		return &payload.Common_Error{
			Msg:       err.Error(),
			Timestamp: fastime.UnixNanoNow(),
		}, err
	}
	return nil, nil
}

func (s *server) StreamInsert(stream agent.Agent_StreamInsertServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return s.Insert(ctx, data.(*payload.Object_Vector))
	})
}

func (s *server) MultiInsert(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Common_Errors, err error) {
	res = new(payload.Common_Errors)
	for _, vec := range vecs.GetVectors() {
		r, ierr := s.Insert(ctx, vec)
		if ierr != nil {
			err = errors.Wrap(err, ierr.Error())
			res.Errors = append(res.Errors, r)
		}
	}
	return
}

func (s *server) Update(ctx context.Context, vec *payload.Object_Vector) (*payload.Common_Error, error) {
	err := s.ngt.Update(vec.GetId().GetId(), vec.GetVector())
	if err != nil {
		return &payload.Common_Error{
			Msg:       err.Error(),
			Timestamp: fastime.UnixNanoNow(),
		}, err
	}
	return nil, nil
}

func (s *server) StreamUpdate(stream agent.Agent_StreamUpdateServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return s.Update(ctx, data.(*payload.Object_Vector))
	})
}

func (s *server) MultiUpdate(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Common_Errors, err error) {
	res = new(payload.Common_Errors)
	for _, vec := range vecs.GetVectors() {
		r, ierr := s.Update(ctx, vec)
		if ierr != nil {
			err = errors.Wrap(err, ierr.Error())
			res.Errors = append(res.Errors, r)
		}
	}
	return
}

func (s *server) Remove(ctx context.Context, id *payload.Object_ID) (*payload.Common_Error, error) {
	err := s.ngt.Delete(id.GetId())
	if err != nil {
		return &payload.Common_Error{
			Msg:       err.Error(),
			Timestamp: fastime.UnixNanoNow(),
		}, err
	}
	return nil, nil
}

func (s *server) StreamRemove(stream agent.Agent_StreamRemoveServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return s.Remove(ctx, data.(*payload.Object_ID))
	})
}

func (s *server) MultiRemove(ctx context.Context, ids *payload.Object_IDs) (res *payload.Common_Errors, err error) {
	res = new(payload.Common_Errors)
	for _, id := range ids.GetIds() {
		r, ierr := s.Remove(ctx, id)
		if ierr != nil {
			err = errors.Wrap(err, ierr.Error())
			res.Errors = append(res.Errors, r)
		}
	}
	return
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (*payload.Object_Vector, error) {
	vec, err := s.ngt.GetObject(id.GetId())
	if err != nil {
		return nil, err
	}
	return &payload.Object_Vector{
		Id: &payload.Object_ID{
			Id: id.GetId(),
		},
		Vector: vec,
	}, nil

}

func (s *server) StreamGetObject(stream agent.Agent_StreamGetObjectServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return s.GetObject(ctx, data.(*payload.Object_ID))
	})
}

func (s *server) CreateIndex(ctx context.Context, c *payload.Controll_CreateIndexRequest) (*payload.Common_Empty, error) {
	return nil, s.ngt.CreateIndex(c.GetPoolSize())
}

func (s *server) SaveIndex(context.Context, *payload.Common_Empty) (*payload.Common_Empty, error) {
	return nil, s.ngt.SaveIndex()
}
