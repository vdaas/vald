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

	"github.com/vdaas/vald/apis/grpc/v1/manager/compressor"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/pkg/manager/compressor/service"
)

type Server compressor.BackupServer

type server struct {
	backup     service.Backup
	compressor service.Compressor
	registerer service.Registerer
}

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) GetVector(ctx context.Context, req *payload.Backup_GetVector_Request) (res *payload.Backup_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor.GetVector")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetUuid()
	r, err := s.backup.GetObject(ctx, uuid)
	if err != nil {
		log.Errorf("[GetVector]\tnot found\t%s", err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("GetVector API uuid %s's object not found", uuid), err, info.Get())
	}

	vector, err := s.compressor.Decompress(ctx, r.GetVector())
	if err != nil {
		log.Errorf("[GetVector]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("GetVector API uuid %s's object failed to decompress %#v", uuid, r), err, info.Get())
	}

	return &payload.Backup_Vector{
		Uuid:   r.GetUuid(),
		Vector: vector,
		Ips:    r.GetIps(),
	}, nil
}

func (s *server) Locations(ctx context.Context, req *payload.Backup_Locations_Request) (res *payload.Info_IPs, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor.Locations")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetUuid()
	r, err := s.backup.GetLocation(ctx, uuid)
	if err != nil {
		log.Errorf("[Locations]\tnot found\t%s", err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Locations API uuid %s's location not found", uuid), err, info.Get())
	}

	return &payload.Info_IPs{
		Ip: r,
	}, nil
}

func (s *server) Register(ctx context.Context, vec *payload.Backup_Vector) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor.Register")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err = s.registerer.Register(ctx, vec)
	if err != nil {
		log.Errorf("[Register]\tregisterer returns error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(
			fmt.Sprintf("Register API uuid %s could not processed", vec.GetUuid()), err, info.Get())
	}

	return new(payload.Empty), nil
}

func (s *server) RegisterMulti(ctx context.Context, vecs *payload.Backup_Vectors) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor.RegisterMulti")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	err = s.registerer.RegisterMulti(ctx, vecs)
	if err != nil {
		log.Errorf("[RegisterMulti]\tregisterer returns error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal("RegisterMulti API could not processed", err, info.Get())
	}

	return new(payload.Empty), nil
}

func (s *server) Remove(ctx context.Context, req *payload.Backup_Remove_Request) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor.Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetUuid()
	err = s.backup.Remove(ctx, uuid)
	if err != nil {
		log.Errorf("[Remove]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API uuid %s could not remove", uuid), err, info.Get())
	}

	return new(payload.Empty), nil
}

func (s *server) RemoveMulti(ctx context.Context, req *payload.Backup_Remove_RequestMulti) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor.RemoveMulti")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuids := req.GetUuids()
	err = s.backup.RemoveMultiple(ctx, uuids...)
	if err != nil {
		log.Errorf("[RemoveMulti]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API uuids %#v could not remove", uuids), err, info.Get())
	}

	return new(payload.Empty), nil
}

func (s *server) RegisterIPs(ctx context.Context, req *payload.Backup_IP_Register_Request) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor.RegisterIPs")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetUuid()
	ips := req.GetIps()
	err = s.backup.RegisterIPs(ctx, uuid, ips)
	if err != nil {
		log.Errorf("[RegisterIPs]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("RegisterIPs API uuid %s ips %#v could not register", uuid, ips), err, info.Get())
	}

	return new(payload.Empty), nil
}

func (s *server) RemoveIPs(ctx context.Context, req *payload.Backup_IP_Remove_Request) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor.RemoveIPs")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ips := req.GetIps()
	err = s.backup.RemoveIPs(ctx, ips)
	if err != nil {
		log.Errorf("[RemoveIPs]\tunknown error\t%+v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("RemoveIPs API ips %#v could not remove", ips), err, info.Get())
	}

	return new(payload.Empty), nil
}
