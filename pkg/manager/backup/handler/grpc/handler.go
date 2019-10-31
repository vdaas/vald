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

	"github.com/vdaas/vald/apis/grpc/manager/backup"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/pkg/manager/backup/model"
	"github.com/vdaas/vald/pkg/manager/backup/service"
)

type Server backup.BackupServer

type server struct {
	mySQL service.MySQL
}

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) GetVector(ctx context.Context, oid *payload.Object_ID) (res *payload.Object_MetaVector, err error) {
	meta, err := s.mySQL.GetMeta(ctx, oid.Id)
	if err != nil {
		return nil, err
	}

	return toObjectMetaVector(meta)
}

func (s *server) Locations(ctx context.Context, oid *payload.Object_ID) (res *payload.Info_IPs, err error) {
	ips, err := s.mySQL.GetIPs(ctx, oid.Id)
	if err != nil {
		return nil, err
	}

	return &payload.Info_IPs{
		Ip: ips,
	}, nil
}

func (s *server) Register(ctx context.Context, meta *payload.Object_MetaVector) (res *payload.Empty, err error) {
	m, err := toModelMetaVector(meta)
	if err != nil {
		return nil, err
	}

	err = s.mySQL.SetMeta(ctx, *m)
	if err != nil {
		return nil, err
	}

	return new(payload.Empty), nil
}

func (s *server) RegisterMulti(ctx context.Context, metas *payload.Object_MetaVectors) (res *payload.Empty, err error) {
	var m *model.MetaVector
	ms := make([]model.MetaVector, 0, len(metas.Vectors))
	for _, meta := range metas.Vectors {
		m, err = toModelMetaVector(meta)
		if err != nil {
			return nil, err
		}
		ms = append(ms, *m)
	}

	err = s.mySQL.SetMetas(ctx, ms...)
	if err != nil {
		return nil, err
	}

	return new(payload.Empty), nil
}

func (s *server) Remove(ctx context.Context, oid *payload.Object_ID) (res *payload.Empty, err error) {
	err = s.mySQL.DeleteMeta(ctx, oid.Id)
	if err != nil {
		return nil, err
	}

	return new(payload.Empty), nil
}

func (s *server) RemoveMulti(ctx context.Context, oids *payload.Object_IDs) (res *payload.Empty, err error) {
	uuids := make([]string, 0, len(oids.Ids))
	for _, oid := range oids.Ids {
		uuids = append(uuids, oid.Id)
	}

	err = s.mySQL.DeleteMetas(ctx, uuids...)
	if err != nil {
		return nil, err
	}

	return new(payload.Empty), nil
}

func toObjectMetaVector(meta *model.MetaVector) (res *payload.Object_MetaVector, err error) {
	vector, err := meta.GetVector()
	if err != nil {
		return nil, err
	}

	return &payload.Object_MetaVector{
		Uuid: meta.GetUUID(),
		Meta: meta.GetMeta(),
		Vector: &payload.Object_Vector{
			Id: &payload.Object_ID{
				Id: meta.GetObjectID(),
			},
			Vector: vector,
		},
		Ips: meta.GetIPs(),
	}, nil
}

func toModelMetaVector(obj *payload.Object_MetaVector) (res *model.MetaVector, err error) {
	return &model.MetaVector{
		UUID:     obj.Uuid,
		ObjectID: obj.Vector.Id.Id,
		Vector:   obj.Vector.Vector,
		Meta:     obj.Meta,
		IPs:      obj.Ips,
	}, nil
}
