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
	"strings"

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
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
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

const apiName = "vald/gateway/backup"

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
	if err == nil && len(ips) > 0 {
		return meta, nil
	}
	id, err := s.gateway.Exists(ctx, meta, s.copts...)
	if err != nil {
		err = status.WrapWithNotFound(fmt.Sprintf("Exists API meta %s's uuid not found", meta.GetId()), err,
			&errdetails.RequestInfo{
				RequestId:   meta.GetId(),
				ServingData: errdetails.Serialize(meta),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}
	return id, nil
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
		err = status.WrapWithInvalidArgument("Search API invalid vector argument", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	res, err = s.gateway.Search(ctx, req, s.copts...)
	if err != nil {
		err = status.WrapWithInternal("Search API failed to process search request", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Search",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
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
	il := len(req.GetId())
	if il == 0 {
		err = errors.ErrUUIDNotFound(0)
		err = status.WrapWithInvalidArgument("SearchByID API invalid id argument", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "id is empty",
						Description: err.Error(),
					},
				},
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	res, err = s.gateway.SearchByID(ctx, req, s.copts...)
	if err != nil {
		err = status.WrapWithInternal("SearchByID API failed to process search request", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.SearchByID",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			})
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
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

	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_IDRequest)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamSearchByID/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.SearchByID(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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

	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	rids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.Requests {
		rids = append(rids, req.GetConfig().GetRequestId())
	}

	res, errs = s.gateway.MultiSearch(ctx, reqs, s.copts...)
	if errs != nil {
		err := errs
		err = status.WrapWithInternal("MultiSearch API failed to search", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiSearch",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) MultiSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	rids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.Requests {
		rids = append(rids, req.GetConfig().GetRequestId())
	}

	res, errs = s.gateway.MultiSearchByID(ctx, reqs, s.copts...)
	if errs != nil {
		err := errs
		err = status.WrapWithInternal("MultiSearchByID API failed to search", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiSearchByID",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return res, nil
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
	vl := len(vec.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument("Insert API invalid vector argument", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, err := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if err == nil && id != nil && len(id.GetId()) != 0 {
			err = errors.ErrMetaDataAlreadyExists(uuid)
			err = status.WrapWithAlreadyExists(fmt.Sprintf("Insert API ID = %v already exists", uuid), err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
					ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Insert_Config{SkipStrictExistCheck: true}
		}
	}

	loc, err = s.gateway.Insert(ctx, req, s.copts...)
	if err != nil {
		err = status.WrapWithInternal("Insert API failed to insert to next gateway", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Insert",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}

	bv := &payload.Backup_Vector{
		Uuid:   uuid,
		Ips:    loc.GetIps(),
		Vector: req.GetVector().GetVector(),
	}
	err = s.backup.Register(ctx, bv)
	if err != nil {
		err = status.WrapWithInternal("insert API failed to register vector to backup component", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(bv),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/manager.compressor.v1.Register",
				ResourceName: strings.Join(s.backup.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		rr := &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		}
		_, rerr := s.gateway.Remove(ctx, rr)
		if rerr != nil {
			err = errors.Wrap(err, status.WrapWithInternal("insert API failed to remove miss inserted data caused by failed to store backup component", rerr,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(rr),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Remove",
					ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}).Error())
			log.Error(err)
		}
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Insert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Insert_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamInsert/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Insert(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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

	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	ids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.GetRequests() {
		uuid := req.GetVector().GetId()
		vl := len(req.GetVector().GetVector())
		if vl < algorithm.MinimumVectorDimensionSize {
			err = errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument("MultiInsert API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       fmt.Sprintf("vector dimension size for id: %s", uuid),
							Description: err.Error(),
						},
					},
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
			}
			return nil, err
		}
		if !req.GetConfig().GetSkipStrictExistCheck() {
			id, err := s.Exists(ctx, &payload.Object_ID{
				Id: uuid,
			})
			if err == nil && id != nil && len(id.GetId()) != 0 {
				err = errors.ErrMetaDataAlreadyExists(uuid)
				err = status.WrapWithAlreadyExists(fmt.Sprintf("MultiInsert API ID = %v already exists", uuid), err,
					&errdetails.RequestInfo{
						RequestId:   uuid,
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
						ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
						Owner:        errdetails.ValdResourceOwner,
						Description:  err.Error(),
					}, info.Get())
				if span != nil {
					span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
				}
				return nil, err
			}
			if req.GetConfig() != nil {
				reqs.Requests[i].GetConfig().SkipStrictExistCheck = true
			} else {
				reqs.Requests[i].Config = &payload.Insert_Config{SkipStrictExistCheck: true}
			}
		}
		ids = append(ids, uuid)
	}

	res, err = s.gateway.MultiInsert(ctx, reqs, s.copts...)
	if err != nil {
		err = status.WrapWithInternal("MultiInsert API failed to insert to next gateway", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiInsert.DoMulti",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}

	mvecs := &payload.Backup_Vectors{
		Vectors: make([]*payload.Backup_Vector, 0, len(reqs.GetRequests())),
	}
	for i, req := range reqs.GetRequests() {
		vec := req.GetVector()
		uuid := vec.GetId()
		mvecs.Vectors = append(mvecs.GetVectors(), &payload.Backup_Vector{
			Uuid:   uuid,
			Vector: vec.GetVector(),
			Ips:    res.Locations[i].GetIps(),
		})
	}
	err = s.backup.RegisterMultiple(ctx, mvecs)
	if err != nil {
		err = status.WrapWithInternal("MultiInsert API failed to register vector to backup component", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(mvecs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/manager.compressor.v1.RegisterMultiple",
				ResourceName: strings.Join(s.backup.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		rmr := &payload.Remove_MultiRequest{
			Requests: make([]*payload.Remove_Request, 0, len(reqs.GetRequests())),
		}
		for _, req := range reqs.GetRequests() {
			rmr.Requests = append(rmr.GetRequests(), &payload.Remove_Request{
				Id: &payload.Object_ID{
					Id: req.GetVector().GetId(),
				},
			})
		}
		_, rerr := s.gateway.MultiRemove(ctx, rmr, s.copts...)
		if rerr != nil {
			err = errors.Wrap(err, status.WrapWithInternal("MultiInsert API failed to remove miss inserted data caused by failed to store backup component", rerr,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(ids, ","),
					ServingData: errdetails.Serialize(rmr),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiRemove",
					ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}).Error())
			log.Error(err)
		}
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
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
	uuid := req.GetVector().GetId()
	vl := len(req.GetVector().GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument("Update API invalid vector argument", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, err := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if err != nil || id == nil || len(id.GetId()) == 0 {
			if err == nil {
				err = errors.ErrObjectIDNotFound(uuid)
			}
			err = status.WrapWithNotFound(fmt.Sprintf("Update API ID = %v not found", uuid), err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
					ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Update_Config{SkipStrictExistCheck: true}
		}
	}

	res, err = s.Remove(ctx, &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: uuid,
		},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: true,
		},
	})
	if err != nil {
		err = status.WrapWithInternal("Update API failed to Remove for Insert", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Remove",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}

	res, err = s.Insert(ctx, &payload.Insert_Request{
		Vector: req.GetVector(),
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: true,
			Filters:              req.GetConfig().GetFilters(),
		},
	})
	if err != nil {
		err = status.WrapWithInternal("Update API failed to Insert for Update", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Insert",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}

	return res, nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Update_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Update_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamUpdate/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Update(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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

	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vecs := reqs.GetRequests()
	ids := make([]string, 0, len(vecs))
	ireqs := make([]*payload.Insert_Request, 0, len(vecs))
	rreqs := make([]*payload.Remove_Request, 0, len(vecs))
	for _, vec := range vecs {
		vl := len(vec.GetVector().GetVector())
		if vl < algorithm.MinimumVectorDimensionSize {
			err = errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument("MultiUpdate API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   vec.GetVector().GetId(),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vector dimension size",
							Description: err.Error(),
						},
					},
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
			}
			return nil, err
		}
		uuid := vec.GetVector().GetId()
		if !vec.GetConfig().GetSkipStrictExistCheck() {
			id, err := s.Exists(ctx, &payload.Object_ID{
				Id: uuid,
			})
			if err != nil || id == nil || len(id.GetId()) == 0 {
				if err == nil {
					err = errors.ErrObjectIDNotFound(uuid)
				}
				err = status.WrapWithNotFound(fmt.Sprintf("MultiUpdate API ID = %v not found", uuid), err,
					&errdetails.RequestInfo{
						RequestId:   uuid,
						ServingData: errdetails.Serialize(reqs),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
						ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
						Owner:        errdetails.ValdResourceOwner,
						Description:  err.Error(),
					}, info.Get())
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				return nil, err
			}
			if vec.GetConfig() != nil {
				vec.GetConfig().SkipStrictExistCheck = true
			} else {
				vec.Config = &payload.Update_Config{SkipStrictExistCheck: true}
			}
		}
		ids = append(ids, vec.GetVector().GetId())
		ireqs = append(ireqs, &payload.Insert_Request{
			Vector: vec.GetVector(),
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Filters:              vec.GetConfig().GetFilters(),
			},
		})
		rreqs = append(rreqs, &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: vec.GetVector().GetId(),
			},
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: true,
			},
		})
	}
	locs, err := s.MultiRemove(ctx, &payload.Remove_MultiRequest{
		Requests: rreqs,
	})
	if err != nil {
		err = status.WrapWithInternal("MultiUpdate API failed to Execute MultiRemove for update", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiRemove",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	log.Debugf("uuids %v were removed from %v due to MultiUpdate. MultiInsert will be executed for them soon. Please see detail %#v", ids, locs.GetLocations(), locs)
	locs, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
		Requests: ireqs,
	})
	if err != nil {
		err = status.WrapWithInternal("MultiUpdate API failed to Execute MultiInsert for update", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiInsert",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return locs, nil
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	vec := req.GetVector()
	uuid := vec.GetId()
	vl := len(vec.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument("Upsert API invalid vector argument", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	filters := req.GetConfig().GetFilters()
	id, err := s.Exists(ctx, &payload.Object_ID{
		Id: uuid,
	})
	if err != nil || id == nil || len(id.GetId()) == 0 {
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
		err = status.WrapWithInternal("Upsert API failed to upsert", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Upsert",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Upsert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Upsert_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamUpsert/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Upsert(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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

	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		log.Error(err)
		return err
	}
	return nil
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
		vl := len(vec.GetVector())
		if vl < algorithm.MinimumVectorDimensionSize {
			err = errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument("MultiUpsert API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vector dimension size",
							Description: err.Error(),
						},
					},
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
			}
			return nil, err
		}
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
			if len(updateReqs) <= 0 {
				return nil
			}

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
			if len(insertReqs) <= 0 {
				return nil
			}

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
		err = status.WrapWithInternal("MultiUpsert API failed to upsert", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiUpsert",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
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
		id, err := s.Exists(ctx, id)
		if err != nil || id == nil || len(id.GetId()) == 0 {
			if err == nil {
				err = errors.ErrObjectIDNotFound(id.GetId())
			}
			err = status.WrapWithNotFound(fmt.Sprintf("Remove API ID = %v not found", id.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   id.GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
					ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeNotFound(err.Error()))
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Remove_Config{SkipStrictExistCheck: true}
		}
	}

	loc, err = s.gateway.Remove(ctx, req, s.copts...)
	if err != nil {
		err = status.WrapWithInternal("Remove API failed to remove request", err,
			&errdetails.RequestInfo{
				RequestId:   id.GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Remove",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	err = s.backup.Remove(ctx, id.GetId())
	if err != nil {
		err = status.WrapWithInternal("Remove API failed to remove from backup store", err,
			&errdetails.RequestInfo{
				RequestId:   id.GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/manager.compressor.v1.Remove",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err

	}
	return loc, nil
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Remove_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Remove_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamRemove/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Remove(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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

	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (locs *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.GetRequests() {
		id := req.GetId()
		ids = append(ids, id.GetId())
		if !req.GetConfig().GetSkipStrictExistCheck() {
			sid, err := s.Exists(ctx, id)
			if err != nil || sid == nil || len(sid.GetId()) == 0 {
				if err == nil {
					err = errors.ErrObjectIDNotFound(id.GetId())
				}
				err = status.WrapWithNotFound(fmt.Sprintf("MultiRemove API ID = %v not found", id.GetId()), err,
					&errdetails.RequestInfo{
						RequestId:   id.GetId(),
						ServingData: errdetails.Serialize(reqs),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
						ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
						Owner:        errdetails.ValdResourceOwner,
						Description:  err.Error(),
					}, info.Get())
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				return nil, err
			}
			if reqs.Requests[i].GetConfig() != nil {
				reqs.Requests[i].GetConfig().SkipStrictExistCheck = true
			} else {
				reqs.Requests[i].Config = &payload.Remove_Config{SkipStrictExistCheck: true}
			}
		}
	}
	locs, err = s.gateway.MultiRemove(ctx, reqs, s.copts...)
	if err != nil {
		err = status.WrapWithInternal("MultiRemove API failed to remove", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiRemove",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	err = s.backup.RemoveMultiple(ctx, ids...)
	if err != nil {
		err = status.WrapWithInternal("MultiRemove API failed to remove", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/manager.compressor.v1.Remove",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return locs, nil
}

func (s *server) GetObject(ctx context.Context, req *payload.Object_VectorRequest) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".GetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	mvec, err := s.backup.GetVector(ctx, req.GetId().GetId())
	if err == nil && mvec != nil {
		return &payload.Object_Vector{
			Id:     mvec.GetUuid(),
			Vector: mvec.GetVector(),
		}, nil
	}
	vec, err = s.gateway.GetObject(ctx, req, s.copts...)
	if err != nil {
		err = status.WrapWithNotFound("GetObject API failed to get Object", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetId().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.GetObject",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}
	return vec, nil
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamGetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_VectorRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Object_VectorRequest)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamGetObject/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.GetObject(ctx, req)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
					err = errors.Wrap(st.Err(), err.Error())
				}
				if sspan != nil {
					sspan.SetStatus(trace.StatusCodeInternal(err.Error()))
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

	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		log.Error(err)
		return err
	}
	return nil
}
