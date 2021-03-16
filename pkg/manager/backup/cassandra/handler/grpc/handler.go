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

	"github.com/vdaas/vald/apis/grpc/v1/manager/backup"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/pkg/manager/backup/cassandra/model"
	"github.com/vdaas/vald/pkg/manager/backup/cassandra/service"
)

type Server backup.BackupServer

type server struct {
	cassandra service.Cassandra
}

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) GetVector(ctx context.Context, req *payload.Backup_GetVector_Request) (res *payload.Backup_Compressed_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-backup-cassandra.GetVector")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetUuid()
	vector, err := s.cassandra.GetVector(ctx, uuid)
	if err != nil {
		switch {
		case errors.IsErrCassandraNotFound(err):
			log.Warnf("[GetVector]\tnot found\t%v\t%s", req.Uuid, err.Error())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, status.WrapWithNotFound(fmt.Sprintf("GetVector API cassandra uuid %s's object not found", uuid), err, info.Get())
		case errors.IsErrCassandraUnavailable(err):
			log.Warnf("[GetVector]\tunavailable\t%+v", err)
			if span != nil {
				span.SetStatus(trace.StatusCodeUnavailable(err.Error()))
			}
			return nil, status.WrapWithUnavailable("GetVector API Cassandra unavailable", err, info.Get())

		default:
			log.Errorf("[GetVector]\tunknown error\t%+v", err)
			if span != nil {
				span.SetStatus(trace.StatusCodeUnknown(err.Error()))
			}
			return nil, status.WrapWithUnknown(fmt.Sprintf("GetVector API cassandra uuid %s's unknown error occurred", uuid), err, info.Get())
		}
	}

	return toBackupVector(vector)
}

func (s *server) Locations(ctx context.Context, req *payload.Backup_Locations_Request) (res *payload.Info_IPs, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-backup-cassandra.Locations")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetUuid()
	ips, err := s.cassandra.GetIPs(ctx, uuid)
	if err != nil {
		log.Errorf("[Locations]\tnot found\t%s", err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Locations API uuid %s's location not found", uuid), err, info.Get())
	}

	return &payload.Info_IPs{
		Ip: ips,
	}, nil
}

func (s *server) Register(ctx context.Context, vector *payload.Backup_Compressed_Vector) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-backup-cassandra.Register")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := vector.GetUuid()
	m, err := toModelVector(vector)
	if err != nil {
		log.Errorf("[Register]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Register API uuid %s's could not convert vector to backup format", uuid), err, info.Get())
	}

	err = s.cassandra.SetVector(ctx, m)
	if err != nil {
		log.Errorf("[Register]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Register API uuid %s's failed to backup vector", uuid), err, info.Get())
	}

	return new(payload.Empty), nil
}

func (s *server) RegisterMulti(ctx context.Context, vectors *payload.Backup_Compressed_Vectors) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-backup-cassandra.RegisterMulti")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ms := make([]*model.Vector, 0, len(vectors.GetVectors()))
	for _, vector := range vectors.Vectors {
		var m *model.Vector
		m, err = toModelVector(vector)
		if err != nil {
			log.Errorf("[RegisterMulti]\tunknown error\t%+v", err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, status.WrapWithInternal(fmt.Sprintf("RegisterMulti API uuids %s's could not convert vector to backup format", vector.GetUuid()), err, info.Get())
		}
		ms = append(ms, m)
	}

	err = s.cassandra.SetVectors(ctx, ms...)
	if err != nil {
		log.Errorf("[RegisterMulti]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("RegisterMulti API failed to backup vectors %#v", ms), err, info.Get())
	}

	return new(payload.Empty), nil
}

func (s *server) Remove(ctx context.Context, req *payload.Backup_Remove_Request) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-backup-cassandra.Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetUuid()
	err = s.cassandra.DeleteVector(ctx, uuid)
	if err != nil {
		log.Errorf("[Remove]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API uuid %s's could not DeleteVector", uuid), err, info.Get())
	}

	return new(payload.Empty), nil
}

func (s *server) RemoveMulti(ctx context.Context, req *payload.Backup_Remove_RequestMulti) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-backup-cassandra.RemoveMulti")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuids := req.GetUuids()
	err = s.cassandra.DeleteVectors(ctx, uuids...)
	if err != nil {
		log.Errorf("[RemoveMulti]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("RemoveMulti API uuids %#v could not DeleteVectors", uuids), err, info.Get())
	}

	return new(payload.Empty), nil
}

func (s *server) RegisterIPs(ctx context.Context, req *payload.Backup_IP_Register_Request) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-backup-cassandra.RegisterIPs")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetUuid()
	err = s.cassandra.SetIPs(ctx, uuid, req.Ips...)
	if err != nil {
		log.Errorf("[RegisterIPs]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("RegisterIPs API uuid %s's could not SetIPs", uuid), err, info.Get())
	}

	return new(payload.Empty), nil
}

func (s *server) RemoveIPs(ctx context.Context, req *payload.Backup_IP_Remove_Request) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-backup-cassandra.RemoveIPs")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ips := req.GetIps()
	err = s.cassandra.RemoveIPs(ctx, ips...)
	if err != nil {
		log.Errorf("[RemoveIPs]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("RemoveIPs API uuid %s's could not RemoveIPs", ips), err, info.Get())
	}

	return new(payload.Empty), nil
}

func toBackupVector(vector *model.Vector) (res *payload.Backup_Compressed_Vector, err error) {
	return &payload.Backup_Compressed_Vector{
		Uuid:   vector.UUID,
		Vector: vector.Vector,
		Ips:    vector.IPs,
	}, nil
}

func toModelVector(obj *payload.Backup_Compressed_Vector) (res *model.Vector, err error) {
	return &model.Vector{
		UUID:   obj.Uuid,
		Vector: obj.Vector,
		IPs:    obj.Ips,
	}, nil
}
