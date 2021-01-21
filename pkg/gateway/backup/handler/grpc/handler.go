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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/compressor"
	client "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/gateway/internal/location"
)

type server struct {
	eg                errgroup.Group
	backup            compressor.Client
	gateway           client.Client
	copts             []grpc.CallOption
	streamConcurrency int
}

const apiName = "vald/gateway-backup"

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
	ips, err := s.backup.GetLocation(ctx, meta.GetId())
	if err != nil {
		log.Debug("an error occurred during calling meta Exists:", err)
		return s.gateway.Exists(ctx, meta, s.copts...)
	}
	if len(ips) > 0 {
		return meta, nil
	}
	return nil, status.WrapWithNotFound(fmt.Sprintf("Exists API meta %s's uuid not found", meta.GetId()), err, meta.GetId(), info.Get())
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
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
	vec, err := s.backup.GetVector(ctx, req.GetId())
	if err != nil {
		return s.gateway.SearchByID(ctx, req, s.copts...)
	}
	return s.gateway.Search(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: req.GetConfig(),
	}, s.copts...)
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
			res, err := s.Search(ctx, data.(*payload.Search_Request))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				return &payload.Search_StreamResponse{
					Payload: &payload.Search_StreamResponse_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Search_StreamResponse{
				Payload: &payload.Search_StreamResponse_Response{
					Response: res,
				},
			}, nil
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
			res, err := s.SearchByID(ctx, data.(*payload.Search_IDRequest))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				return &payload.Search_StreamResponse{
					Payload: &payload.Search_StreamResponse_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Search_StreamResponse{
				Payload: &payload.Search_StreamResponse_Response{
					Response: res,
				},
			}, nil
		})
}

func (s *server) MultiSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return s.gateway.MultiSearch(ctx, reqs, s.copts...)
}

func (s *server) MultiSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	return s.gateway.MultiSearchByID(ctx, reqs, s.copts...)
}

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	uuid := vec.GetId()
	if len(vec.GetVector()) < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(len(vec.GetVector()), 0)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, status.WrapWithInvalidArgument("Insert API invalid vector argument", err, req, info.Get())
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		locs, err := s.backup.GetLocation(ctx, uuid)
		if err != nil {
			log.Debug("an error occurred while calling meta Exists:", err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, status.WrapWithInternal(
				fmt.Sprintf("Insert API ID %s couldn't check meta already exists or not", uuid), err, info.Get())
		}
		if len(locs) > 0 {
			err = errors.Wrap(err, errors.ErrMetaDataAlreadyExists(uuid).Error())
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, status.WrapWithAlreadyExists(fmt.Sprintf("Insert API ID %s already exists", vec.GetId()), err, info.Get())
		}
		req.Config.SkipStrictExistCheck = true
	}

	loc, err = s.gateway.Insert(ctx, req, s.copts...)
	if err != nil {
		err = errors.Wrapf(err, "Insert API failed to Insert uuid = %s\tinfo = %#v", uuid, info.Get())
		log.Debug(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API failed to Execute DoMulti error = %s", err.Error()), err, info.Get())
	}
	vecs := &payload.Backup_Vector{
		Uuid: uuid,
		Ips:  loc.GetIps(),
	}
	if vec != nil {
		vecs.Vector = vec.GetVector()
	}
	err = s.backup.Register(ctx, vecs)
	if err != nil {
		_, rerr := s.gateway.Remove(ctx, &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		if rerr != nil {
			err = errors.Wrap(err, rerr.Error())
		}
		err = errors.Wrapf(err, "Insert API (backup.Register) failed to Backup Vectors = %#v\t info = %#v", vecs, info.Get())
		log.Debug(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Insert API uuid %s couldn't store", uuid), err, info.Get())
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
			res, err := s.Insert(ctx, data.(*payload.Insert_Request))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: res,
				},
			}, nil
		})
}

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	for i, req := range reqs.GetRequests() {
		if !req.GetConfig().GetSkipStrictExistCheck() {
			id := req.GetVector().GetId()
			loc, err := s.backup.GetLocation(ctx, id)
			if err != nil {
				log.Debug("an error occurred during calling meta Exists:", err)
				if span != nil {
					span.SetStatus(trace.StatusCodeInternal(err.Error()))
				}
				return nil, status.WrapWithInternal(
					fmt.Sprintf("MultiInsert API couldn't check metadata exists or not metas = %v", id), err, info.Get())
			}
			if len(loc) > 0 {
				if span != nil {
					span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
				}
				return nil, status.WrapWithAlreadyExists(
					fmt.Sprintf("MultiInsert API failed metadata already exists meta = %s", id), err, info.Get())
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

	mvecs := &payload.Backup_Vectors{
		Vectors: make([]*payload.Backup_Vector, 0, len(reqs.GetRequests())),
	}
	for i, req := range reqs.GetRequests() {
		vec := req.GetVector()
		uuid := vec.GetId()
		mvecs.Vectors = append(mvecs.Vectors, &payload.Backup_Vector{
			Uuid:   uuid,
			Vector: vec.GetVector(),
			Ips:    res.Locations[i].GetIps(),
		})
	}
	err = s.backup.RegisterMultiple(ctx, mvecs)
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
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiInsert API failed RegisterMultiple %#v", mvecs), err, info.Get())
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
	id := req.GetVector().GetId()
	res, err = s.Remove(ctx, &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: id,
		},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: true,
		},
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Update API failed to remove exsisting data for update %#v", req), err, info.Get())
	}
	res, err = s.Insert(ctx, &payload.Insert_Request{
		Vector: &payload.Object_Vector{
			Id:     id,
			Vector: req.GetVector().GetVector(),
		},
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: true,
			Filters:              req.GetConfig().GetFilters(),
		},
	})
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
			res, err := s.Update(ctx, data.(*payload.Update_Request))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: res,
				},
			}, nil
		})
}

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	rreqs := make([]*payload.Remove_Request, 0, len(reqs.GetRequests()))
	ireqs := make([]*payload.Insert_Request, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		rreqs = append(rreqs, &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: req.GetVector().GetId(),
			},
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: true,
			},
		})
		ireqs = append(ireqs, &payload.Insert_Request{
			Vector: req.GetVector(),
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Filters:              req.GetConfig().GetFilters(),
			},
		})
	}
	_, err = s.MultiRemove(ctx, &payload.Remove_MultiRequest{
		Requests: rreqs,
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed Remove request %#v", rreqs), err, info.Get())
	}
	res, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
		Requests: ireqs,
	})
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiUpdate API failed Insert request %#v", ireqs), err, info.Get())
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
	id := vec.GetId()
	filters := req.GetConfig().GetFilters()
	ips, err := s.backup.GetLocation(ctx, req.GetVector().GetId())
	if err != nil {
		log.Debug("an error occurred during calling meta Exists:", err)
	}
	if len(ips) <= 0 {
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
		log.Debugf("Upsert API failed to process request uuid:\t%s\terror:\t%s", id, err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Upsert API failed to Upsert request %#v", req), err, info.Get())
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
			res, err := s.Upsert(ctx, data.(*payload.Upsert_Request))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: res,
				},
			}, nil
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
	id := req.GetId()
	if !req.GetConfig().GetSkipStrictExistCheck() {
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
		req.Config.SkipStrictExistCheck = true
	}

	loc, err = s.gateway.Remove(ctx, req, s.copts...)
	if err != nil {
		log.Debugf("Remove API failed to process request uuid:\t%s\terror:\t%s", id.GetId(), err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed to Remove backup uuid = %s", id.GetId()), err, info.Get())
	}
	err = s.backup.Remove(ctx, id.GetId())
	if err != nil {
		log.Debugf("Remove API failed to remove backup data\tid:\t%s\terror:\t%s", id.GetId(), err.Error())
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("Remove API failed to Remove backup uuid = %s", id.GetId()), err, info.Get())
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
			res, err := s.Remove(ctx, data.(*payload.Remove_Request))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: res,
				},
			}, nil
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
		id := req.GetId().GetId()
		ids = append(ids, id)
		if !req.GetConfig().GetSkipStrictExistCheck() {
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
	locs, err = s.gateway.MultiRemove(ctx, reqs, s.copts...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiRemove API failed to Remove backup uuids = %v", ids), err, info.Get())
	}
	err = s.backup.RemoveMultiple(ctx, ids...)
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("MultiRemove API failed to Remove backup uuids %v ", ids), err, info.Get())
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
	mvec, err := s.backup.GetVector(ctx, id.GetId())
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
			res, err := s.GetObject(ctx, data.(*payload.Object_ID))
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				return &payload.Object_StreamVector{
					Payload: &payload.Object_StreamVector_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamVector{
				Payload: &payload.Object_StreamVector_Vector{
					Vector: res,
				},
			}, nil
		})
}
