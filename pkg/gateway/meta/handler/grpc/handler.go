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
	"time"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
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
)

type server struct {
	eg                errgroup.Group
	metadata          service.Meta
	gateway           client.Client
	copts             []grpc.CallOption
	timeout           time.Duration
	replica           int
	streamConcurrency int
}

func New(opts ...Option) vald.ValdServer {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}

	return s
}

func (s *server) Exists(ctx context.Context, meta *payload.Object_ID) (*payload.Object_ID, error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.Exists")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid, err := s.metadata.GetUUID(ctx, meta.GetId())
	if err != nil {
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
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.Search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return s.search(ctx, func(ctx context.Context, vc vald.ValdClient, copts ...grpc.CallOption) (*payload.Search_Response, error) {
		return vc.Search(ctx, req, copts...)
	})
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.SearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	metaID := req.GetId()
	req.Id, err = s.metadata.GetUUID(ctx, metaID)
	if err != nil {
		req.Id = metaID
		log.Errorf("error at SearchByID\t%v", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("SearchByID API meta %s's uuid not found", metaID), err, req, info.Get())
	}
	return s.search(ctx, func(ctx context.Context, vc vald.ValdClient, copts ...grpc.CallOption) (*payload.Search_Response, error) {
		return vc.SearchByID(ctx, req, copts...)
	})
}

func (s *server) search(ctx context.Context,
	f func(ctx context.Context, vc vald.ValdClient, copts ...grpc.CallOption) (*payload.Search_Response, error)) (
	res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.search")
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
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-meta.StreamSearch")
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
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-meta.StreamSearchByID")
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

func (s *server) Insert(ctx context.Context, vec *payload.Object_Vector) (ce *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.Insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	meta := vec.GetId()
	exists, err := s.metadata.Exists(ctx, meta)
	if err != nil {
		log.Error(err)
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
	_, err = s.gateway.Insert(ctx, vec, s.copts...)
	if err != nil {
		err = errors.Wrapf(err, "Insert API (do multiple) failed to Insert uuid = %s\tmeta = %s\t info = %#v", uuid, meta, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API failed to Execute DoMulti error = %s", err.Error()), err, info.Get())
	}
	err = s.metadata.SetUUIDandMeta(ctx, uuid, meta)
	if err != nil {
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API meta %s & uuid %s couldn't store", meta, uuid), err, info.Get())
	}
	return new(payload.Empty), nil
}

func (s *server) StreamInsert(stream vald.Vald_StreamInsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-meta.StreamInsert")
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

func (s *server) MultiInsert(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	metaMap := make(map[string]string)
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
			log.Error(err)
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
		log.Error(err)
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

func (s *server) Update(ctx context.Context, vec *payload.Object_Vector) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.Update")
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

	return new(payload.Empty), nil
}

func (s *server) StreamUpdate(stream vald.Vald_StreamUpdateServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-meta.StreamUpdate")
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

func (s *server) MultiUpdate(ctx context.Context, vecs *payload.Object_Vectors) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.MultiUpdate")
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
	return new(payload.Empty), nil
}

func (s *server) Upsert(ctx context.Context, vec *payload.Object_Vector) (*payload.Empty, error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	meta := vec.GetId()
	exists, errs := s.metadata.Exists(ctx, meta)
	if errs != nil {
		log.Error(errs)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(errs.Error()))
		}
	}

	if exists {
		_, err := s.Update(ctx, vec)
		if err != nil {
			errs = errors.Wrap(errs, err.Error())
		}
	} else {
		_, err := s.Insert(ctx, vec)
		if err != nil {
			errs = errors.Wrap(errs, err.Error())
		}
	}

	return new(payload.Empty), errs
}

func (s *server) StreamUpsert(stream vald.Vald_StreamUpsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-meta.StreamUpsert")
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

func (s *server) MultiUpsert(ctx context.Context, vecs *payload.Object_Vectors) (*payload.Empty, error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.MultiUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	insertVecs := make([]*payload.Object_Vector, 0, len(vecs.GetVectors()))
	updateVecs := make([]*payload.Object_Vector, 0, len(vecs.GetVectors()))

	var errs error
	for _, vec := range vecs.GetVectors() {
		exists, err := s.metadata.Exists(ctx, vec.GetId())
		if err != nil {
			log.Error(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			errs = errors.Wrap(errs, err.Error())
		}

		if exists {
			updateVecs = append(updateVecs, vec)
		} else {
			insertVecs = append(insertVecs, vec)
		}
	}

	eg, ectx := errgroup.New(ctx)

	eg.Go(safety.RecoverFunc(func() error {
		var err error
		if len(updateVecs) > 0 {
			_, err = s.MultiUpdate(ectx, &payload.Object_Vectors{
				Vectors: updateVecs,
			})
		}
		return err
	}))

	eg.Go(safety.RecoverFunc(func() error {
		var err error
		if len(insertVecs) > 0 {
			_, err = s.MultiInsert(ectx, &payload.Object_Vectors{
				Vectors: insertVecs,
			})
		}
		return err
	}))

	err := eg.Wait()
	if err != nil {
		errs = errors.Wrap(errs, err.Error())
		return nil, status.WrapWithInternal("MultiUpsert API failed", errs, info.Get())
	}

	return new(payload.Empty), errs
}

func (s *server) Remove(ctx context.Context, id *payload.Object_ID) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.Remove")
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
		return nil, status.WrapWithNotFound(fmt.Sprintf("Remove API meta %s's uuid not found", meta), err, info.Get())
	}

	id.Id = uuid
	res, err = s.gateway.Remove(ctx, id, s.copts...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed request uuid %s", uuid), err, info.Get())
	}
	_, err = s.metadata.DeleteMeta(ctx, uuid)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed Delete metadata uuid = %s", uuid), err, info.Get())
	}
	return res, nil
}

func (s *server) StreamRemove(stream vald.Vald_StreamRemoveServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-meta.StreamRemove")
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

func (s *server) MultiRemove(ctx context.Context, ids *payload.Object_IDs) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuids, err := s.metadata.GetUUIDs(ctx, ids.GetIds()...)
	if err != nil {
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

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (vec *payload.Backup_MetaVector, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-meta.GetObject")
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
	return vec, nil
}

func (s *server) StreamGetObject(stream vald.Vald_StreamGetObjectServer) error {
	ctx, span := trace.StartSpan(stream.Context(), "vald/gateway-meta.StreamGetObject")
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
