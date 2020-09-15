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

	"github.com/kpango/fuid"
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
	"github.com/vdaas/vald/pkg/gateway/meta/service"
	"github.com/vdaas/vald/pkg/gateway/tool/location"
)

type server struct {
	eg                errgroup.Group
	metadata          service.Meta
	gateway           client.Client
	copts             []grpc.CallOption
	streamConcurrency int
}

const apiName = "vald/gateway-meta"

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
	uuid, err := s.metadata.GetUUID(ctx, meta.GetId())
	if err != nil {
		log.Debugf("Exists API failed to get uuid:\t%s\terror:\t%s", meta.GetId(), err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Exists API meta %s's uuid not found", meta.GetId()), err, meta.GetId(), info.Get())
	}
	return &payload.Object_ID{
		Id: uuid,
	}, nil
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.search(ctx, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
		return vc.Search(ctx, req, copts...)
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal("Search API failed to process search request", err, req, info.Get())
	}
	return res, nil
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".SearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	meta := req.GetId()
	req.Id, err = s.metadata.GetUUID(ctx, meta)
	if err != nil {
		log.Debugf("MultiRemove API failed to process request uuids:\t%v\terror:\t%s", meta, err.Error())
		req.Id = meta
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("SearchByID API meta %s's uuid not found", meta), err, req, info.Get())
	}
	res, err = s.search(ctx, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
		return vc.SearchByID(ctx, req, copts...)
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound("Search API failed to process search request", err, req, info.Get())
	}
	return res, nil
}

func (s *server) search(ctx context.Context,
	f func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error)) (
	res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = f(ctx, s.gateway, s.copts...)
	if err != nil {
		return nil, err
	}
	uuids := make([]string, 0, len(res.Results))
	for _, r := range res.Results {
		uuids = append(uuids, r.GetId())
	}
	if s.metadata != nil {
		var metas []string
		metas, err = s.metadata.GetMetas(ctx, uuids...)
		for i, k := range metas {
			if len(k) != 0 {
				res.Results[i].Id = k
			}
		}
	}
	return res, err
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
	meta := vec.GetId()
	exists, err := s.metadata.Exists(ctx, meta)
	if err != nil {
		log.Debug(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(
			fmt.Sprintf("Insert API meta %s couldn't check meta already exists or not", meta), err, info.Get())
	}
	if exists {
		err = errors.Wrap(err, errors.ErrMetaDataAlreadyExists(meta).Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
		}
		return nil, status.WrapWithAlreadyExists(fmt.Sprintf("Insert API meta %s already exists", meta), err, info.Get())
	}
	uuid := fuid.String()
	vec.Id = uuid
	loc, err := s.gateway.Insert(ctx, vec, s.copts...)
	if err != nil {
		err = errors.Wrapf(err, "Insert API (do multiple) failed to Insert uuid = %s\tmeta = %s\t info = %#v", uuid, meta, info.Get())
		log.Debug(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API failed to Execute DoMulti error = %s", err.Error()), err, info.Get())
	}
	err = s.metadata.SetUUIDandMeta(ctx, uuid, meta)
	if err != nil {
		log.Debug(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API meta %s & uuid %s couldn't store", meta, uuid), err, info.Get())
	}
	return loc, nil
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
	metaMap := make(map[string]string, len(vecs.GetVectors()))
	metas := make([]string, 0, len(vecs.GetVectors()))
	for i, vec := range vecs.GetVectors() {
		uuid := fuid.String()
		meta := vec.GetId()
		metaMap[uuid] = meta
		metas = append(metas, meta)
		vecs.Vectors[i].Id = uuid
	}

	for _, meta := range metas {
		exists, err := s.metadata.Exists(ctx, meta)
		if err != nil {
			log.Debug(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, status.WrapWithInternal(
				fmt.Sprintf("MultiInsert API couldn't check metadata exists or not metas = %v", metas), err, info.Get())
		}
		if exists {
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, status.WrapWithAlreadyExists(
				fmt.Sprintf("MultiInsert API failed metadata already exists meta = %s", meta), err, info.Get())
		}
	}

	res, err = s.gateway.MultiInsert(ctx, vecs, s.copts...)
	if err != nil {
		log.Debug(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal("got error from MultiInsert API", err, info.Get())
	}

	err = s.metadata.SetUUIDandMetas(ctx, metaMap)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiInsert API failed SetUUIDandMetas %#v", metaMap), err, info.Get())
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
	meta := vec.GetId()
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Update API failed GetUUID meta = %s", meta), err, info.Get())
	}
	vec.Id = uuid
	res, err = s.gateway.Update(ctx, vec, s.copts...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Update API failed request %#v", vec), err, info.Get())
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
	metas, err := s.metadata.GetUUIDs(ctx, ids...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed MultiUpdate request %#v", ids), err, info.Get())
	}
	for i := range vecs.GetVectors() {
		vecs.Vectors[i].Id = metas[i]
	}
	res, err = s.MultiUpdate(ctx, vecs)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed MultiUpdate request %#v", ids), err, info.Get())
	}
	return res, err
}

func (s *server) Upsert(ctx context.Context, vec *payload.Object_Vector) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	exists, err := s.metadata.Exists(ctx, vec.GetId())
	if err != nil {
		log.Debugf("Upsert API metadata exists check error:\t%s", err.Error())
	}

	if exists {
		loc, err = s.Update(ctx, vec)
	} else {
		loc, err = s.Insert(ctx, vec)
	}
	if err != nil {
		log.Debugf("Upsert API failed to process request uuid:\t%s\terror:\t%s", vec.GetId(), err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Upsert API failed to process request %s", vec.GetId()), err, info.Get())
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
		exists, err := s.metadata.Exists(ctx, vec.GetId())
		if err != nil {
			log.Debugf("Upsert API Exists returns error %s", err.Error())
		}

		if exists {
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

func (s *server) Remove(ctx context.Context, id *payload.Object_ID) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	meta := id.GetId()
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil {
		log.Debugf("Remove API failed to get uuid:\t%s\terror:\t%s", meta, err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Remove API meta %s's uuid not found", meta), err, info.Get())
	}

	id.Id = uuid
	res, err = s.gateway.Remove(ctx, id, s.copts...)
	if err != nil {
		log.Debugf("Remove API failed to process request uuid:\t%s\terror:\t%s", meta, err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed request uuid %s", uuid), err, info.Get())
	}
	_, err = s.metadata.DeleteMeta(ctx, uuid)
	if err != nil {
		log.Debugf("Remove API failed to remove metadata:\t%s\terror:\t%s", meta, err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed Delete metadata uuid = %s", uuid), err, info.Get())
	}
	return res, nil
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

func (s *server) MultiRemove(ctx context.Context, ids *payload.Object_IDs) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuids, err := s.metadata.GetUUIDs(ctx, ids.GetIds()...)
	if err != nil {
		log.Debugf("MultiRemove API failed to process request uuids:\t%v\terror:\t%s", ids.GetIds(), err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("MultiRemove API meta datas %v's uuid not found", ids.GetIds()), err, info.Get())
	}
	ids.Ids = uuids
	res, err = s.gateway.MultiRemove(ctx, ids, s.copts...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiRemove API failed to request uuids %v metas %v ", uuids, ids.GetIds()), err, info.Get())
	}
	_, err = s.metadata.DeleteMetas(ctx, uuids...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiRemove API failed to DeleteMetas uuids %v ", uuids), err, info.Get())
	}
	return res, nil
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".GetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	meta := id.GetId()
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("GetObject API meta %s's uuid not found", meta), err, info.Get())
	}
	id.Id = uuid
	vec, err = s.gateway.GetObject(ctx, id, s.copts...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("GetObject API meta %s uuid %s Object not found", meta, uuid), err, info.Get())
	}
	vec.Id = meta
	return vec, nil
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
