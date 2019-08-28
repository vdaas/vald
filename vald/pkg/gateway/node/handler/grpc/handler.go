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

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/apis/grpc/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/gateway/vald/service"
)

type Server vald.ValdServer

type server struct {
	gateway service.ValdProxy
}

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Exists(ctx context.Context, oid *payload.Object_ID) (*payload.Object_ID, error) {
	return nil, nil
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (*payload.Search_Response, error) {
	return nil, nil
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_Response, error) {
	return nil, nil
}

func (s *server) StreamSearch(stream vald.Vald_StreamSearchServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return nil, nil
	})
}

func (s *server) StreamSearchByID(stream vald.Vald_StreamSearchByIDServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return nil, nil
	})
}

func (s *server) Insert(ctx context.Context, vec *payload.Object_Vector) (*payload.Common_Error, error) {
	return nil, nil
}

func (s *server) StreamInsert(stream vald.Vald_StreamInsertServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return nil, nil
	})
}

func (s *server) MultiInsert(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Common_Errors, err error) {
	return nil, nil
}

func (s *server) Update(ctx context.Context, vec *payload.Object_Vector) (*payload.Common_Error, error) {
	return nil, nil
}

func (s *server) StreamUpdate(stream vald.Vald_StreamUpdateServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return nil, nil
	})
}

func (s *server) MultiUpdate(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Common_Errors, err error) {
	return nil, nil
}

func (s *server) Remove(ctx context.Context, id *payload.Object_ID) (*payload.Common_Error, error) {
	return nil, nil
}

func (s *server) StreamRemove(stream vald.Vald_StreamRemoveServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return nil, nil
	})
}

func (s *server) MultiRemove(ctx context.Context, ids *payload.Object_IDs) (res *payload.Common_Errors, err error) {
	return nil, nil
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (*payload.Object_Vector, error) {
	return nil, nil
}

func (s *server) StreamGetObject(stream vald.Vald_StreamGetObjectServer) error {
	return grpc.BidirectionalStream(stream, func(ctx context.Context, data interface{}) (interface{}, error) {
		return nil, nil
	})
}
