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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	client "github.com/vdaas/vald/internal/client/gateway/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/gateway/tool/location"
	"github.com/vdaas/vald/pkg/gateway/vald/service"
)

type server struct {
	eg                errgroup.Group
	gateway           client.Client
	backup            service.Backup
	copts             []grpc.CallOption
	strict            bool
	streamConcurrency int
}

const apiName = "vald/gateway-backup"

func New(opts ...Option) vald.Server {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Exists(ctx context.Context, meta *payload.Object_ID) (*payload.Object_ID, error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Exists")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return s.gateway.Exists(ctx, meta, s.copts...)
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(req.Vector) < 2 {
		return nil, errors.ErrInvalidDimensionSize(len(req.Vector), 0)
	}
	return s.gateway.Search(ctx, req, s.copts...)

}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".SearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return s.gateway.SearchByID(ctx, req, s.copts...)
}

func (s *server) StreamSearch(stream vald.Vald_StreamSearchServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Search(ctx, data.(*payload.Search_Request))
		})
}

func (s *server) StreamSearchByID(stream vald.Vald_StreamSearchByIDServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.SearchByID(ctx, data.(*payload.Search_IDRequest))
		})
}

func (s *server) Insert(ctx context.Context, vec *payload.Object_Vector) (ce *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(vec.Vector) < 2 {
		return nil, errors.ErrInvalidDimensionSize(len(vec.Vector), 0)
	}

	if s.strict {
		locs, err := s.backup.GetLocation(ctx, vec.GetId())
		if err != nil {
			log.Debug("an error occurred while calling meta Exists:", err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, status.WrapWithInternal(
				fmt.Sprintf("Insert API ID %s couldn't check meta already exists or not", vec.GetId()), err, info.Get())
		}
		if len(locs) > 0 {
			err = errors.Wrap(err, errors.ErrMetaDataAlreadyExists(vec.GetId()).Error())
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, status.WrapWithAlreadyExists(fmt.Sprintf("Insert API ID %s already exists", vec.GetId()), err, info.Get())
		}
	}

	ce, err = s.gateway.Insert(ctx, vec, s.copts...)
	if err != nil {
		err = errors.Wrapf(err, "Insert API failed to Insert uuid = %s\tinfo = %#v", vec.GetId(), info.Get())
		log.Debug(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API failed to Execute DoMulti error = %s", err.Error()), err, info.Get())
	}
	vecs := &payload.Backup_MetaVector{
		Uuid: vec.GetId(),
		Ips:  ce.GetIps(),
	}
	if vec != nil {
		vecs.Vector = vec.GetVector()
	}
	err = s.backup.Register(ctx, vecs)
	if err != nil {
		err = errors.Wrapf(err, "Insert API (backup.Register) failed to Backup Vectors = %#v\t info = %#v", vecs, info.Get())
		log.Debug(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(err.Error(), err)
	}
	return ce, nil
}

func (s *server) StreamInsert(stream vald.Vald_StreamInsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Insert(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiInsert(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if s.strict {
		for _, vec := range vecs.GetVectors() {
			loc, err := s.backup.GetLocation(ctx, vec.GetId())
			if err != nil {
				log.Debug("an error occurred during calling meta Exists:", err)
				if span != nil {
					span.SetStatus(trace.StatusCodeInternal(err.Error()))
				}
				return nil, status.WrapWithInternal(
					fmt.Sprintf("MultiInsert API couldn't check metadata exists or not metas = %v", vec.GetId()), err, info.Get())
			}
			if len(loc) > 0 {
				if span != nil {
					span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
				}
				return nil, status.WrapWithAlreadyExists(
					fmt.Sprintf("MultiInsert API failed metadata already exists meta = %s", vec.GetId()), err, info.Get())
			}
		}
	}

	res, err = s.gateway.MultiInsert(ctx, vecs, s.copts...)
	if err != nil {
		err = errors.Wrapf(err, "MultiInsert API failed to Insert info = %#v", info.Get())
		log.Debug(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API failed to Execute DoMulti error = %s", err.Error()), err, info.Get())
	}

	mvecs := &payload.Backup_MetaVectors{
		Vectors: make([]*payload.Backup_MetaVector, 0, len(vecs.GetVectors())),
	}
	for i, vec := range vecs.GetVectors() {
		uuid := vec.GetId()
		mvecs.Vectors = append(mvecs.Vectors, &payload.Backup_MetaVector{
			Uuid:   uuid,
			Vector: vec.GetVector(),
			Ips:    res.Locations[i].GetIps(),
		})
	}
	err = s.backup.RegisterMultiple(ctx, mvecs)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiInsert API failed RegisterMultiple %#v", mvecs), err, info.Get())
	}
	return res, nil
}

func (s *server) Update(ctx context.Context, vec *payload.Object_Vector) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Update")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.Remove(ctx, &payload.Object_ID{
		Id: vec.GetId(),
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Update API failed to remove exsisting data for update %#v", vec), err, info.Get())
	}
	res, err = s.Insert(ctx, vec)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Update API failed to remove exsisting data for update %#v", vec), err, info.Get())
	}
	return res, nil
}

func (s *server) StreamUpdate(stream vald.Vald_StreamUpdateServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Update(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiUpdate(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ids := make([]string, 0, len(vecs.GetVectors()))
	for _, vec := range vecs.GetVectors() {
		ids = append(ids, vec.GetId())
	}
	_, err = s.MultiRemove(ctx, &payload.Object_IDs{
		Ids: ids,
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed Remove request %#v", ids), err, info.Get())
	}
	_, err = s.MultiInsert(ctx, vecs)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed Insert request %#v", vecs), err, info.Get())
	}
	return new(payload.Object_Locations), nil
}

func (s *server) Upsert(ctx context.Context, vec *payload.Object_Vector) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	ips, err := s.backup.GetLocation(ctx, vec.GetId())
	if err != nil {
		log.Debug("an error occurred during calling meta Exists:", err)
	}
	if len(ips) > 0 {
		loc, err = s.Update(ctx, vec)
	} else {
		loc, err = s.Insert(ctx, vec)
	}
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Upsert API failed to Upsert request %#v", vec), err, info.Get())
	}

	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Vald_StreamUpsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_Vector) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Upsert(ctx, data.(*payload.Object_Vector))
		})
}

func (s *server) MultiUpsert(ctx context.Context, vecs *payload.Object_Vectors) (*payload.Object_Locations, error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	insertVecs := make([]*payload.Object_Vector, 0, len(vecs.GetVectors()))
	updateVecs := make([]*payload.Object_Vector, 0, len(vecs.GetVectors()))

	ids := make([]string, 0, len(vecs.GetVectors()))
	for _, vec := range vecs.GetVectors() {
		ids = append(ids, vec.GetId())
		ips, err := s.backup.GetLocation(ctx, vec.GetId())
		if err != nil {
			log.Debug("an error occurred during calling meta Exists:", err)
		}
		if len(ips) > 0 {
			updateVecs = append(updateVecs, vec)
		} else {
			insertVecs = append(insertVecs, vec)
		}
	}

	eg, ectx := errgroup.New(ctx)
	insertLocs := make([]*payload.Object_Location, 0, len(insertVecs))
	updateLocs := make([]*payload.Object_Location, 0, len(updateVecs))
	eg.Go(safety.RecoverFunc(func() error {
		ectx, span := trace.StartSpan(ectx, apiName+".MultiUpsert/Go-MultiUpdate")
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		var err error
		if len(updateVecs) > 0 {
			loc, err := s.MultiUpdate(ectx, &payload.Object_Vectors{
				Vectors: updateVecs,
			})
			if err == nil {
				updateLocs = loc.GetLocations()
			}
		}
		return err
	}))
	eg.Go(safety.RecoverFunc(func() error {
		ectx, span := trace.StartSpan(ectx, apiName+".MultiUpsert/Go-MultiInsert")

		defer func() {
			if span != nil {
				span.End()
			}
		}()
		var err error
		if len(insertVecs) > 0 {
			loc, err := s.MultiInsert(ectx, &payload.Object_Vectors{
				Vectors: insertVecs,
			})
			if err == nil {
				insertLocs = loc.GetLocations()
			}
		}
		return err
	}))
	err := eg.Wait()
	if err != nil {
		log.Debugf("MultiUpsert API failed to process request uuids:\t%v\terror:\t%s", ids, err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpsert API failed to process request %v", ids), err, info.Get())
	}

	return location.ReStructure(ids, &payload.Object_Locations{
		Locations: append(insertLocs, updateLocs...),
	}), nil
}

func (s *server) Remove(ctx context.Context, id *payload.Object_ID) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if s.strict {
		ips, err := s.backup.GetLocation(ctx, id.GetId())
		if err != nil {
			log.Debug("an error occurred while calling meta Exists:", err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, status.WrapWithInternal(
				fmt.Sprintf("Remove API ID %s couldn't check meta already exists or not", id.GetId()), err, info.Get())
		}
		if len(ips) <= 0 {
			err = errors.Wrap(err, errors.ErrMetaDataAlreadyExists(id.GetId()).Error())
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, status.WrapWithAlreadyExists(fmt.Sprintf("Remove API ID %s not found", id.GetId()), err, info.Get())
		}
	}

	loc, err = s.gateway.Remove(ctx, id, s.copts...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed to Remove backup uuid = %s", id.GetId()), err, info.Get())
	}
	err = s.backup.Remove(ctx, id.GetId())
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed to Remove backup uuid = %s", id.GetId()), err, info.Get())
	}
	return loc, nil
}

func (s *server) StreamRemove(stream vald.Vald_StreamRemoveServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Remove(ctx, data.(*payload.Object_ID))
		})
}

func (s *server) MultiRemove(ctx context.Context, ids *payload.Object_IDs) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if s.strict {
		for _, id := range ids.GetIds() {
			ips, err := s.backup.GetLocation(ctx, id)
			if err != nil {
				log.Debug("an error occurred while calling meta Exists:", err)
				if span != nil {
					span.SetStatus(trace.StatusCodeInternal(err.Error()))
				}
				return nil, status.WrapWithInternal(
					fmt.Sprintf("MultiRemove API ID %s couldn't check meta already exists or not", id), err, info.Get())
			}
			if len(ips) <= 0 {
				err = errors.Wrap(err, errors.ErrMetaDataAlreadyExists(id).Error())
				if span != nil {
					span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
				}
				return nil, status.WrapWithAlreadyExists(fmt.Sprintf("MultiRemove API ID %s not found", id), err, info.Get())
			}
		}
	}

	locs, err = s.gateway.MultiRemove(ctx, ids, s.copts...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiRemove API failed to Remove backup uuids = %v", ids.GetIds()), err, info.Get())
	}
	err = s.backup.RemoveMultiple(ctx, ids.GetIds()...)
	if err != nil {

		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiRemove API failed to Remove backup uuids %v ", ids.GetIds()), err, info.Get())
	}
	return locs, nil
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".GetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	mvec, err := s.backup.GetObject(ctx, id.GetId())
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("GetObject API uuid %s Object not found", id.GetId()), err, info.Get())
	}
	return &payload.Object_Vector{
		Id:     mvec.GetUuid(),
		Vector: mvec.GetVector(),
	}, nil
}

func (s *server) StreamGetObject(stream vald.Vald_StreamGetObjectServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamGetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_ID) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.GetObject(ctx, data.(*payload.Object_ID))
		})
}
