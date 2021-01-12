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
	"sync"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/gateway/internal/location"
	"github.com/vdaas/vald/pkg/gateway/meta/service"
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

	for _, opt := range append(defaultOptions, opts...) {
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
	vl := len(req.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("Search API invalid vector argument", err, req, info.Get())
	}
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
		log.Debugf("SearchByID API failed to process request uuids:\t%v\terror:\t%s", meta, err.Error())
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
		return nil, status.WrapWithNotFound("SearchByID API failed to process search request", err, req, info.Get())
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

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) error {
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

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) error {
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

func (s *server) MultiSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.Requests)),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.Search(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				errs = errors.Wrap(errs, status.WrapWithNotFound(fmt.Sprintf("MultiSearch API vector %v's search request result not found", query.GetVector()), err, info.Get()).Error())
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		})
	}
	wg.Wait()
	return res, errs
}

func (s *server) MultiSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.Requests)),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(func() error {
			defer wg.Done()
			r, err := s.SearchByID(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				errs = errors.Wrap(errs, status.WrapWithNotFound(fmt.Sprintf("MultiSearchByID API uuid %v's search by id request result not found", query.GetId()), err, info.Get()).Error())
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		})
	}
	wg.Wait()
	return res, errs
}

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	meta := vec.GetId()
	vl := len(vec.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("Insert API invalid vector argument", err, req, info.Get())
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
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
		req.Config.SkipStrictExistCheck = true
	}
	uuid := fuid.String()
	req.Vector.Id = uuid
	loc, err = s.gateway.Insert(ctx, req, s.copts...)
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
		_, rerr := s.gateway.Remove(ctx, &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		if rerr != nil {
			err = errors.Wrap(err, rerr.Error())
		}
		err = errors.Wrapf(err, "Insert API (meta.SetUUIDandMeta) failed to Register Metadata Vectors = %#v\t info = %#v", req, info.Get())
		log.Debug(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API meta %s & uuid %s couldn't store", meta, uuid), err, info.Get())
	}
	return loc, nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Insert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Insert(ctx, data.(*payload.Insert_Request))
		})
}

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vecs := reqs.GetRequests()
	metaMap := make(map[string]string, len(vecs))
	metas := make([]string, 0, len(vecs))
	for i, req := range vecs {
		vec := req.GetVector()
		uuid := fuid.String()
		meta := vec.GetId()
		metaMap[uuid] = meta
		metas = append(metas, meta)
		reqs.Requests[i].Vector.Id = uuid
		if !req.GetConfig().GetSkipStrictExistCheck() {
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
			reqs.Requests[i].Config.SkipStrictExistCheck = true
		}
	}

	res, err = s.gateway.MultiInsert(ctx, reqs, s.copts...)
	if err != nil {
		err = errors.Wrapf(err, "MultiInsert API failed to Insert info = %#v", info.Get())
		log.Debug(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiInsert API failed to Insert error = %s", err.Error()), err, info.Get())
	}

	err = s.metadata.SetUUIDandMetas(ctx, metaMap)
	if err != nil {
		removeList := make([]*payload.Remove_Request, 0, len(reqs.GetRequests()))
		for _, req := range reqs.GetRequests() {
			removeList = append(removeList, &payload.Remove_Request{
				Id: &payload.Object_ID{
					Id: req.GetVector().GetId(),
				},
			})
		}
		_, rerr := s.gateway.MultiRemove(ctx, &payload.Remove_MultiRequest{
			Requests: removeList,
		}, s.copts...)
		if rerr != nil {
			err = errors.Wrap(err, rerr.Error())
		}
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiInsert API failed SetUUIDandMetas %#v", metaMap), err, info.Get())
	}
	return res, nil
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Update")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	meta := req.GetVector().GetId()
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Update API failed GetUUID meta = %s", meta), err, info.Get())
	}
	req.Vector.Id = uuid
	res, err = s.gateway.Update(ctx, req, s.copts...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Update API failed to insert data for update %#v", req), err, info.Get())
	}
	return res, nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Update_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Update(ctx, data.(*payload.Update_Request))
		})
}

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		ids = append(ids, req.GetVector().GetId())
	}
	metas, err := s.metadata.GetUUIDs(ctx, ids...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed MultiUpdate request %#v", ids), err, info.Get())
	}

	for i, meta := range metas {
		reqs.Requests[i].Vector.Id = meta
	}

	res, err = s.gateway.MultiUpdate(ctx, reqs)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed MultiUpdate request %#v", ids), err, info.Get())
	}
	return res, nil
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	vec := req.GetVector()
	meta := vec.GetId()
	filters := req.GetConfig().GetFilters()
	exists, err := s.metadata.Exists(ctx, meta)
	if err != nil {
		log.Debugf("Upsert API metadata exists check error:\t%s", err.Error())
	}
	if !exists {
		loc, err = s.Insert(ctx, &payload.Insert_Request{
			Vector: vec,
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Filters:              filters,
			},
		})
	} else {
		loc, err = s.Update(ctx, &payload.Update_Request{
			Vector: vec,
			Config: &payload.Update_Config{
				SkipStrictExistCheck: true,
				Filters:              filters,
			},
		})
	}
	if err != nil {
		log.Debugf("Upsert API failed to process request uuid:\t%s\terror:\t%s", meta, err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Upsert API failed to process request %s", meta), err, info.Get())
	}
	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Upsert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Upsert(ctx, data.(*payload.Upsert_Request))
		})
}

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	insertReqs := make([]*payload.Insert_Request, 0, len(reqs.GetRequests()))
	updateReqs := make([]*payload.Update_Request, 0, len(reqs.GetRequests()))

	ids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		vec := req.GetVector()
		uuid := vec.GetId()
		ids = append(ids, uuid)
		_, err = s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		filters := req.GetConfig().GetFilters()
		if err != nil {
			insertReqs = append(insertReqs, &payload.Insert_Request{
				Vector: vec,
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
					Filters:              filters,
				},
			})
		} else {
			updateReqs = append(updateReqs, &payload.Update_Request{
				Vector: vec,
				Config: &payload.Update_Config{
					SkipStrictExistCheck: true,
					Filters:              filters,
				},
			})
		}
	}

	insertLocs := make([]*payload.Object_Location, 0, len(insertReqs))
	updateLocs := make([]*payload.Object_Location, 0, len(updateReqs))

	eg, ectx := errgroup.New(ctx)
	if len(updateReqs) <= 0 {
		eg.Go(safety.RecoverFunc(func() error {
			ectx, span := trace.StartSpan(ectx, apiName+".MultiUpsert/Go-MultiUpdate")
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			var err error
			loc, err := s.MultiUpdate(ectx, &payload.Update_MultiRequest{
				Requests: updateReqs,
			})
			if err == nil {
				updateLocs = loc.GetLocations()
			}
			return err
		}))
	}
	if len(insertReqs) <= 0 {
		eg.Go(safety.RecoverFunc(func() error {
			ectx, span := trace.StartSpan(ectx, apiName+".MultiUpsert/Go-MultiInsert")
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			var err error
			loc, err := s.MultiInsert(ectx, &payload.Insert_MultiRequest{
				Requests: insertReqs,
			})
			if err == nil {
				insertLocs = loc.GetLocations()
			}
			return err
		}))
	}
	err = eg.Wait()
	if err != nil {
		log.Debugf("MultiUpsert API failed to process request uuids:\t%s\terror:\t%s", ids, err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpsert API failed to process request %v", ids), err, info.Get())
	}

	return location.ReStructure(ids, &payload.Object_Locations{
		Locations: append(insertLocs, updateLocs...),
	}), nil
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	meta := req.GetId().GetId()
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil {
		log.Debugf("Remove API failed to get uuid:\t%s\terror:\t%s", meta, err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Remove API meta %s's uuid not found", meta), err, info.Get())
	}

	req.Id.Id = uuid
	loc, err = s.gateway.Remove(ctx, req, s.copts...)
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
	return loc, nil
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) error {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Remove_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			return s.Remove(ctx, data.(*payload.Remove_Request))
		})
}

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		ids = append(ids, req.GetId().GetId())
	}
	uuids, err := s.metadata.GetUUIDs(ctx, ids...)
	if err != nil {
		log.Debugf("MultiRemove API failed to process request uuids:\t%v\terror:\t%s", ids, err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("MultiRemove API meta datas %v's uuid not found", ids), err, info.Get())
	}
	for i, id := range uuids {
		reqs.Requests[i].Id.Id = id
	}
	locs, err = s.gateway.MultiRemove(ctx, reqs, s.copts...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiRemove API failed to request uuids %v metas %v ", uuids, ids), err, info.Get())
	}
	_, err = s.metadata.DeleteMetas(ctx, uuids...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiRemove API failed to DeleteMetas uuids %v ", uuids), err, info.Get())
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

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) error {
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
