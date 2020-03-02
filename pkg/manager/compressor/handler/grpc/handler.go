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

	"github.com/vdaas/vald/apis/grpc/manager/compressor"
	"github.com/vdaas/vald/apis/grpc/payload"
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
}

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) GetVector(ctx context.Context, req *payload.Backup_GetVector_Request) (res *payload.Backup_MetaVector, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor.GetVector")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetUuid()
	r, err := s.backup.GetObject(ctx, uuid)
	if err != nil {
		log.Errorf("[GetVector]\tunknown error\t%+v", err)
		return nil, status.WrapWithNotFound(fmt.Sprintf("GetVector API uuid %s's object not found", uuid), err, info.Get())
	}

	vector, err := s.compressor.Decompress(ctx, r.GetVector())
	if err != nil {
		log.Errorf("[GetVector]\tunknown error\t%+v", err)
		return nil, status.WrapWithInternal(fmt.Sprintf("GetVector API uuid %s's object failed to decompress %#v", uuid, r), err, info.Get())
	}

	return &payload.Backup_MetaVector{
		Uuid:   r.GetUuid(),
		Meta:   r.GetMeta(),
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
		log.Errorf("[Locations]\tunknown error\t%+v", err)
		return nil, status.WrapWithNotFound(fmt.Sprintf("Locations API uuid %s's location not found", uuid), err, info.Get())
	}

	return &payload.Info_IPs{
		Ip: r,
	}, nil
}

func (s *server) Register(ctx context.Context, meta *payload.Backup_MetaVector) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor.Register")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := meta.GetUuid()
	vector, err := s.compressor.Compress(ctx, meta.GetVector())
	if err != nil {
		log.Errorf("[Register]\tunknown error\t%+v", err)
		return nil, status.WrapWithInternal(fmt.Sprintf("Register API uuid %s's could not compress", uuid), err, info.Get())
	}

	mvec := &payload.Backup_Compressed_MetaVector{
		Uuid:   meta.GetUuid(),
		Meta:   meta.GetMeta(),
		Vector: vector,
		Ips:    meta.GetIps(),
	}

	err = s.backup.Register(ctx, mvec)
	if err != nil {
		log.Errorf("[Register]\tunknown error\t%+v", err)
		return nil, status.WrapWithInternal(fmt.Sprintf("Register API uuid %s's could not register %#v", uuid, mvec), err, info.Get())
	}

	return new(payload.Empty), nil
}

func (s *server) RegisterMulti(ctx context.Context, metas *payload.Backup_MetaVectors) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-compressor.RegisterMulti")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	mvs := metas.GetVectors()
	vectors := make([][]float32, 0, len(mvs))
	for _, mv := range mvs {
		vectors = append(vectors, mv.GetVector())
	}

	compressedVecs, err := s.compressor.MultiCompress(ctx, vectors)
	if err != nil {
		log.Errorf("[RegisterMulti]\tinternal error\t%+v", err)
		uuids := make([]string, 0, len(mvs))
		for _, mv := range mvs {
			uuids = append(uuids, mv.GetUuid())
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("RegisterMulti API uuids %#v's could not compress", uuids), err, info.Get())
	}

	compressedMVs := make([]*payload.Backup_Compressed_MetaVector, 0, len(mvs))
	for i, mv := range mvs {
		compressedMVs = append(compressedMVs, &payload.Backup_Compressed_MetaVector{
			Uuid:   mv.GetUuid(),
			Meta:   mv.GetMeta(),
			Vector: compressedVecs[i],
			Ips:    mv.GetIps(),
		})
	}

	err = s.backup.RegisterMultiple(ctx, &payload.Backup_Compressed_MetaVectors{
		Vectors: compressedMVs,
	})
	if err != nil {
		log.Errorf("[RegisterMulti]\tunknown error\t%+v", err)
		uuids := make([]string, 0, len(mvs))
		for _, mv := range mvs {
			uuids = append(uuids, mv.GetUuid())
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("RegisterMulti API uuids %#v's could not register %#v", uuids, compressedMVs), err, info.Get())
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
		return nil, status.WrapWithInternal(fmt.Sprintf("RemoveIPs API ips %#v could not remove", ips), err, info.Get())
	}

	return new(payload.Empty), nil
}
