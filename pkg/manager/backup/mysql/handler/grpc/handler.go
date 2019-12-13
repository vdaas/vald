//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

	"github.com/vdaas/vald/apis/grpc/manager/backup"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/model"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/service"
)

type Server backup.BackupServer

type server struct {
	mySQL service.MySQL
}

type errDetail struct {
	method string
	uuid   string
	uuids  []string
}

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) GetVector(ctx context.Context, req *payload.Backup_GetVector_Request) (res *payload.Backup_MetaVector, err error) {
	meta, err := s.mySQL.GetMeta(ctx, req.Uuid)
	if err != nil {
		detail := errDetail{method: "GetVector", uuid: req.Uuid}
		if errors.IsErrMySQLNotFound(errors.UnWrapAll(err)) {
			return nil, status.WrapWithNotFound("MySQL entry not found", &detail, err)
		}
		return nil, status.WrapWithUnknown("Unknown error occurred", &detail, err)
	}

	return toBackupMetaVector(meta)
}

func (s *server) Locations(ctx context.Context, req *payload.Backup_Locations_Request) (res *payload.Info_IPs, err error) {
	ips, err := s.mySQL.GetIPs(ctx, req.Uuid)
	if err != nil {
		return nil, status.WrapWithUnknown("Unknown error occurred", &errDetail{method: "Locations", uuid: req.Uuid}, err)
	}

	return &payload.Info_IPs{
		Ip: ips,
	}, nil
}

func (s *server) Register(ctx context.Context, meta *payload.Backup_MetaVector) (res *payload.Empty, err error) {
	m, err := toModelMetaVector(meta)
	if err != nil {
		return nil, status.WrapWithUnknown("Unknown error occurred", &errDetail{method: "Register", uuid: meta.Uuid}, err)
	}

	err = s.mySQL.SetMeta(ctx, *m)
	if err != nil {
		detail := errDetail{method: "Register", uuid: meta.Uuid}
		if errors.IsErrMySQLInvalidArgument(errors.UnWrapAll(err)) {
			return nil, status.WrapWithInvalidArgument("MySQL invalid argument", &detail, err)
		}
		return nil, status.WrapWithUnknown("Unknown error occurred", &detail, err)
	}

	return new(payload.Empty), nil
}

func (s *server) RegisterMulti(ctx context.Context, metas *payload.Backup_MetaVectors) (res *payload.Empty, err error) {
	ms := make([]model.MetaVector, 0, len(metas.Vectors))
	for _, meta := range metas.Vectors {
		var m *model.MetaVector
		m, err = toModelMetaVector(meta)
		if err != nil {
			return nil, status.WrapWithUnknown("Unknown error occurred", &errDetail{method: "RegisterMulti", uuid: meta.Uuid}, err)
		}
		ms = append(ms, *m)
	}

	err = s.mySQL.SetMetas(ctx, ms...)
	if err != nil {
		detail := errDetail{method: "RegisterMulti"}
		if errors.IsErrMySQLInvalidArgument(errors.UnWrapAll(err)) {
			return nil, status.WrapWithInvalidArgument("MySQL invalid argument", &detail, err)
		}
		return nil, status.WrapWithUnknown("Unknown error occurred", &detail, err)
	}

	return new(payload.Empty), nil
}

func (s *server) Remove(ctx context.Context, req *payload.Backup_Remove_Request) (res *payload.Empty, err error) {
	err = s.mySQL.DeleteMeta(ctx, req.Uuid)
	if err != nil {
		return nil, status.WrapWithUnknown("Unknown error occurred", &errDetail{method: "Remove", uuid: req.Uuid}, err)
	}

	return new(payload.Empty), nil
}

func (s *server) RemoveMulti(ctx context.Context, req *payload.Backup_Remove_RequestMulti) (res *payload.Empty, err error) {
	err = s.mySQL.DeleteMetas(ctx, req.GetUuid()...)
	if err != nil {
		return nil, status.WrapWithUnknown("Unknown error occurred", &errDetail{method: "RemoveMulti", uuids: req.GetUuid()}, err)
	}

	return new(payload.Empty), nil
}

func (s *server) RegisterIPs(ctx context.Context, req *payload.Backup_IP_Register_Request) (res *payload.Empty, err error) {
	err = s.mySQL.SetIPs(ctx, req.Uuid, req.Ips...)
	if err != nil {
		return nil, status.WrapWithUnknown("Unknown error occurred", &errDetail{method: "RegisterIPs", uuid: req.Uuid}, err)
	}

	return new(payload.Empty), nil
}

func (s *server) RemoveIPs(ctx context.Context, req *payload.Backup_IP_Remove_Request) (res *payload.Empty, err error) {
	err = s.mySQL.RemoveIPs(ctx, req.Ips...)
	if err != nil {
		return nil, status.WrapWithUnknown("Unknown error occurred", &errDetail{method: "RemoveIPs"}, err)
	}

	return new(payload.Empty), nil
}

func toBackupMetaVector(meta *model.MetaVector) (res *payload.Backup_MetaVector, err error) {
	vector, err := meta.GetVector()
	if err != nil {
		return nil, err
	}

	return &payload.Backup_MetaVector{
		Uuid:   meta.GetUUID(),
		Meta:   meta.GetMeta(),
		Vector: vector,
		Ips:    meta.GetIPs(),
	}, nil
}

func toModelMetaVector(obj *payload.Backup_MetaVector) (res *model.MetaVector, err error) {
	return &model.MetaVector{
		UUID:   obj.Uuid,
		Vector: obj.Vector,
		Meta:   obj.Meta,
		IPs:    obj.Ips,
	}, nil
}
