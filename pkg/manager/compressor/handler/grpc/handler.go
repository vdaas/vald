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

	"github.com/vdaas/vald/apis/grpc/manager/compressor"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/pkg/manager/compressor/service"
)

type Server compressor.BackupServer

type server struct {
	backup     service.Backup
	compressor service.Compressor
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
	r, err := s.backup.GetObject(ctx, req.GetUuid())
	if err != nil {
		log.Errorf("[GetVector]\tunknown error\t%+v", err)
		detail := errDetail{method: "GetVector", uuid: req.Uuid}
		return nil, status.WrapWithUnknown("Unknown error occurred", &detail, err)
	}

	vector, err := s.compressor.Decompress(ctx, r.GetVector())
	if err != nil {
		log.Errorf("[GetVector]\tunknown error\t%+v", err)
		detail := errDetail{method: "GetVector", uuid: req.Uuid}
		return nil, status.WrapWithInternal("Internal error occurred", &detail, err)
	}

	return &payload.Backup_MetaVector{
		Uuid:   r.GetUuid(),
		Meta:   r.GetMeta(),
		Vector: vector,
		Ips:    r.GetIps(),
	}, nil
}

func (s *server) Locations(ctx context.Context, req *payload.Backup_Locations_Request) (res *payload.Info_IPs, err error) {
	r, err := s.backup.GetLocation(ctx, req.GetUuid())
	if err != nil {
		log.Errorf("[Locations]\tunknown error\t%+v", err)
		detail := errDetail{method: "Locations", uuid: req.Uuid}
		return nil, status.WrapWithUnknown("Unknown error occurred", &detail, err)
	}

	return &payload.Info_IPs{
		Ip: r,
	}, nil
}

func (s *server) Register(ctx context.Context, meta *payload.Backup_MetaVector) (res *payload.Empty, err error) {
	vector, err := s.compressor.Compress(ctx, meta.GetVector())
	if err != nil {
		log.Errorf("[Register]\tunknown error\t%+v", err)
		detail := errDetail{method: "Register", uuid: meta.Uuid}
		return nil, status.WrapWithInternal("Internal error occurred", &detail, err)
	}

	err = s.backup.Register(ctx, &payload.Backup_Compressed_MetaVector{
		Uuid:   meta.GetUuid(),
		Meta:   meta.GetMeta(),
		Vector: vector,
		Ips:    meta.GetIps(),
	})
	if err != nil {
		log.Errorf("[Register]\tunknown error\t%+v", err)
		detail := errDetail{method: "Register", uuid: meta.Uuid}
		return nil, status.WrapWithUnknown("Unknown error occurred", &detail, err)
	}

	return new(payload.Empty), nil
}

func (s *server) RegisterMulti(ctx context.Context, metas *payload.Backup_MetaVectors) (res *payload.Empty, err error) {
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
		detail := errDetail{method: "RegisterMulti", uuids: uuids}
		return nil, status.WrapWithInternal("Internal error occurred", &detail, err)
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
		detail := errDetail{method: "RegisterMulti", uuids: uuids}
		return nil, status.WrapWithUnknown("Unknown error occurred", &detail, err)
	}

	return new(payload.Empty), nil
}

func (s *server) Remove(ctx context.Context, req *payload.Backup_Remove_Request) (res *payload.Empty, err error) {
	err = s.backup.Remove(ctx, req.GetUuid())
	if err != nil {
		log.Errorf("[Remove]\tunknown error\t%+v", err)
		detail := errDetail{method: "Remove", uuid: req.GetUuid()}
		return nil, status.WrapWithUnknown("Unknown error occurred", &detail, err)
	}

	return new(payload.Empty), nil
}

func (s *server) RemoveMulti(ctx context.Context, req *payload.Backup_Remove_RequestMulti) (res *payload.Empty, err error) {
	err = s.backup.RemoveMultiple(ctx, req.GetUuid()...)
	if err != nil {
		log.Errorf("[RemoveMulti]\tunknown error\t%+v", err)
		detail := errDetail{method: "RemoveMulti", uuids: req.GetUuid()}
		return nil, status.WrapWithUnknown("Unknown error occurred", &detail, err)
	}

	return new(payload.Empty), nil
}

func (s *server) RegisterIPs(ctx context.Context, req *payload.Backup_IP_Register_Request) (res *payload.Empty, err error) {
	err = s.backup.RegisterIPs(ctx, req.GetUuid(), req.GetIps())
	if err != nil {
		log.Errorf("[RegisterIPs]\tunknown error\t%+v", err)
		detail := errDetail{method: "RegisterIPs", uuid: req.GetUuid()}
		return nil, status.WrapWithUnknown("Unknown error occurred", &detail, err)
	}

	return new(payload.Empty), nil
}

func (s *server) RemoveIPs(ctx context.Context, req *payload.Backup_IP_Remove_Request) (res *payload.Empty, err error) {
	err = s.backup.RemoveIPs(ctx, req.GetIps())
	if err != nil {
		log.Errorf("[RemoveIPs]\tunknown error\t%+v", err)
		detail := errDetail{method: "RemoveIPs"}
		return nil, status.WrapWithUnknown("Unknown error occurred", &detail, err)
	}

	return new(payload.Empty), nil
}
