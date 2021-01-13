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

// Package grpc provides grpc client functions
package grpc

import (
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	proto "github.com/yahoojapan/ngtd/proto"
)

type (
	StreamSearch     vald.Search_StreamSearchClient
	StreamSearchByID vald.Search_StreamSearchByIDClient
	StreamInsert     vald.Insert_StreamInsertClient
	StreamUpdate     vald.Update_StreamUpdateClient
	StreamUpsert     vald.Upsert_StreamUpsertClient
	StreamRemove     vald.Remove_StreamRemoveClient
	StreamObject     vald.Object_StreamGetObjectClient
)

type streamSearch struct {
	grpc.ClientStream
	ngtd proto.NGTD_StreamSearchClient
}

func NewStreamSearchClient(ngtd proto.NGTD_StreamSearchClient) StreamSearch {
	return &streamSearch{
		ClientStream: ngtd,
		ngtd:         ngtd,
	}
}

func (s *streamSearch) Send(req *payload.Search_Request) error {
	vec := make([]float64, 0, len(req.GetVector()))
	for _, v := range req.GetVector() {
		vec = append(vec, float64(v))
	}
	return s.ngtd.Send(&proto.SearchRequest{
		Vector: vec,
	})
}

func (s *streamSearch) Recv() (*payload.Search_Response, error) {
	data, err := s.ngtd.Recv()
	if err != nil {
		return nil, err
	}
	if len(data.GetError()) != 0 {
		return nil, errors.New(data.GetError())
	}
	res := &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, len(data.GetResult())),
	}
	for _, dist := range data.GetResult() {
		res.Results = append(res.Results, &payload.Object_Distance{
			Distance: dist.GetDistance(),
			Id:       string(dist.GetId()),
		})
	}
	return res, nil
}

type streamSearchByID struct {
	grpc.ClientStream
	ngtd proto.NGTD_StreamSearchByIDClient
}

func NewStreamSearchByIDClient(ngtd proto.NGTD_StreamSearchByIDClient) StreamSearchByID {
	return &streamSearchByID{
		ClientStream: ngtd,
		ngtd:         ngtd,
	}
}

func (s *streamSearchByID) Send(req *payload.Search_IDRequest) error {
	return s.ngtd.Send(&proto.SearchRequest{
		Id: []byte(req.GetId()),
	})
}

func (s *streamSearchByID) Recv() (*payload.Search_Response, error) {
	data, err := s.ngtd.Recv()
	if err != nil {
		return nil, err
	}
	if len(data.GetError()) != 0 {
		return nil, errors.New(data.GetError())
	}
	res := &payload.Search_Response{
		Results: make([]*payload.Object_Distance, 0, len(data.GetResult())),
	}
	for _, dist := range data.GetResult() {
		res.Results = append(res.Results, &payload.Object_Distance{
			Distance: dist.GetDistance(),
			Id:       string(dist.GetId()),
		})
	}
	return res, nil
}

type streamInsert struct {
	grpc.ClientStream
	ngtd proto.NGTD_StreamInsertClient
}

func NewStreamInsertClient(ngtd proto.NGTD_StreamInsertClient) StreamInsert {
	return &streamInsert{
		ClientStream: ngtd,
		ngtd:         ngtd,
	}
}

func (s *streamInsert) Send(req *payload.Insert_Request) error {
	vec := make([]float64, 0, len(req.GetVector().GetVector()))
	for _, v := range req.GetVector().GetVector() {
		vec = append(vec, float64(v))
	}
	return s.ngtd.Send(&proto.InsertRequest{
		Id:     []byte(req.GetVector().GetId()),
		Vector: vec,
	})
}

func (s *streamInsert) Recv() (*payload.Object_Location, error) {
	data, err := s.ngtd.Recv()
	if err != nil {
		return nil, err
	}
	if len(data.GetError()) != 0 {
		return nil, errors.New(data.GetError())
	}
	return nil, nil
}

type streamUpdate struct {
	grpc.ClientStream
	ic proto.NGTD_StreamInsertClient
	rc proto.NGTD_StreamRemoveClient
}

func NewStreamUpdateClient(ic proto.NGTD_StreamInsertClient, rc proto.NGTD_StreamRemoveClient) StreamUpdate {
	return &streamUpdate{
		ClientStream: ic,
		ic:           ic,
		rc:           rc,
	}
}

func (s *streamUpdate) Send(req *payload.Update_Request) error {
	vec := make([]float64, 0, len(req.GetVector().GetVector()))
	for _, v := range req.GetVector().GetVector() {
		vec = append(vec, float64(v))
	}
	err := s.rc.Send(&proto.RemoveRequest{
		Id: []byte(req.GetVector().GetId()),
	})
	if err != nil {
		return err
	}
	err = s.ic.Send(&proto.InsertRequest{
		Id:     []byte(req.GetVector().GetId()),
		Vector: vec,
	})
	return err
}

func (s *streamUpdate) Recv() (*payload.Object_Location, error) {
	rdata, err := s.rc.Recv()
	if err != nil {
		return nil, err
	}
	if len(rdata.GetError()) != 0 {
		return nil, errors.New(rdata.GetError())
	}
	idata, err := s.ic.Recv()
	if err != nil {
		return nil, err
	}
	if len(idata.GetError()) != 0 {
		return nil, errors.New(idata.GetError())
	}
	return nil, nil
}

type streamUpsert struct {
	grpc.ClientStream
	cc Client
	ch chan *payload.Object_Location
}

func NewStreamUpsertClient(c Client, ic proto.NGTD_StreamInsertClient) StreamUpsert {
	return &streamUpsert{
		ClientStream: ic,
		cc:           c,
		ch:           make(chan *payload.Object_Location, 10),
	}
}

func (s *streamUpsert) Send(req *payload.Upsert_Request) error {
	go func() {
		ctx := s.ClientStream.Context()
		id, err := s.cc.Exists(ctx, &payload.Object_ID{
			Id: req.GetVector().GetId(),
		})
		var loc *payload.Object_Location
		if err == nil || len(id.GetId()) != 0 {
			loc, err = s.cc.Update(ctx, &payload.Update_Request{
				Vector: req.GetVector(),
			})
		} else {
			loc, err = s.cc.Insert(ctx, &payload.Insert_Request{
				Vector: req.GetVector(),
			})
		}
		if err == nil {
			s.ch <- loc
		}
	}()
	return nil
}

func (s *streamUpsert) Recv() (loc *payload.Object_Location, err error) {
	ctx := s.ClientStream.Context()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case loc := <-s.ch:
		return loc, err
	}
}

type streamRemove struct {
	grpc.ClientStream
	ngtd proto.NGTD_StreamRemoveClient
}

func NewStreamRemoveClient(ngtd proto.NGTD_StreamRemoveClient) StreamRemove {
	return &streamRemove{
		ClientStream: ngtd,
		ngtd:         ngtd,
	}
}

func (s *streamRemove) Send(req *payload.Remove_Request) error {
	return s.ngtd.Send(&proto.RemoveRequest{
		Id: []byte(req.GetId().GetId()),
	})
}

func (s *streamRemove) Recv() (*payload.Object_Location, error) {
	data, err := s.ngtd.Recv()
	if err != nil {
		return nil, err
	}
	if len(data.GetError()) != 0 {
		return nil, errors.New(data.GetError())
	}
	return nil, nil
}

type streamGetObject struct {
	grpc.ClientStream
	ngtd proto.NGTD_StreamGetObjectClient
}

func NewStreamObjectClient(ngtd proto.NGTD_StreamGetObjectClient) StreamObject {
	return &streamGetObject{
		ClientStream: ngtd,
		ngtd:         ngtd,
	}
}

func (s *streamGetObject) Send(req *payload.Object_ID) error {
	return s.ngtd.Send(&proto.GetObjectRequest{
		Id: []byte(req.GetId()),
	})
}

func (s *streamGetObject) Recv() (*payload.Object_Vector, error) {
	data, err := s.ngtd.Recv()
	if err != nil {
		return nil, err
	}
	if len(data.GetError()) != 0 {
		return nil, errors.New(data.GetError())
	}
	return &payload.Object_Vector{
		Id:     string(data.GetId()),
		Vector: data.GetVector(),
	}, nil
}
