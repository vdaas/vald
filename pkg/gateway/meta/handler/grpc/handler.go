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
	"sync"
	"time"

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
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
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
		err = status.WrapWithNotFound(fmt.Sprintf("Exists API meta %s's uuid not found", meta.GetId()), err,
			&errdetails.RequestInfo{
				RequestId:   meta.GetId(),
				ServingData: errdetails.Serialize(meta),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.GetUUID",
				ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
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
	res, err = s.search(ctx, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
		return vc.Search(ctx, req, copts...)
	})
	if err != nil {
		st, msg, err := status.ParseError(err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.SearchByID",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
	meta := req.GetId()
	req.Id, err = s.metadata.GetUUID(ctx, meta)
	if err != nil {
		req.Id = meta
		err = status.WrapWithNotFound("SearchByID API uuid could not found", err,
			&errdetails.RequestInfo{
				RequestId:   meta,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.GetUUID",
				ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}
	res, err = s.search(ctx, func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Search_Response, error) {
		return vc.SearchByID(ctx, req, copts...)
	})
	if err != nil {
		st, msg, err := status.ParseError(err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.SearchByID",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
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
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
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
		metas, merr := s.metadata.GetMetas(ctx, uuids...)
		if merr != nil {
			merr = status.WrapWithNotFound("search API search result's metadata could not found", merr,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ","),
					ServingData: errdetails.Serialize(uuids),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.GetMetas",
					ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  merr.Error(),
				}, info.Get())
			log.Warn(merr)
			if err == nil {
				err = merr
			} else {
				err = errors.Wrap(merr, err.Error())
			}
		} else {
			for i, k := range metas {
				if len(k) != 0 {
					res.Results[i].Id = k
				}
			}
		}
	}
	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
	}
	return res, err
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

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.Requests)),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		vl := len(req.GetVector())
		if vl < algorithm.MinimumVectorDimensionSize {
			err := errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument("MultiSearch API invalid vector argument", err,
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
		wg.Add(1)
		s.eg.Go(func() error {
			ctx, span := trace.StartSpan(ctx, fmt.Sprintf("%s.MultiSearch/goroutine/id-%d", apiName, idx))
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			defer wg.Done()
			r, err := s.Search(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		})
	}
	wg.Wait()
	if errs != nil {
		err := errs
		err = status.WrapWithInternal("MultiSearch API failed to search", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
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
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(func() error {
			ctx, span := trace.StartSpan(ctx, fmt.Sprintf("%s.MultiSearchByID/goroutine/id-%d", apiName, idx))
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			defer wg.Done()
			r, err := s.SearchByID(ctx, query)
			if err != nil {
				if span != nil {
					span.SetStatus(trace.StatusCodeNotFound(err.Error()))
				}
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		})
	}
	wg.Wait()
	if errs != nil {
		err := errs
		err = status.WrapWithInternal("MultiSearchByID API failed to search", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.SearchByID",
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
	meta := vec.GetId()
	vl := len(vec.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument("Insert API invalid vector argument", err,
			&errdetails.RequestInfo{
				RequestId:   meta,
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
		exists, err := s.metadata.Exists(ctx, meta)
		if err != nil {
			err = status.WrapWithInternal("Insert API failed to check metadata exsisting or not", err,
				&errdetails.RequestInfo{
					RequestId:   meta,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.GetMetaInverse",
					ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			log.Error(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInternal(err.Error()))
			}
			return nil, err
		}
		if exists {
			err = errors.ErrMetaDataAlreadyExists(meta)
			err = status.WrapWithAlreadyExists(fmt.Sprintf("Insert API ID = %v already exists", meta), err,
				&errdetails.RequestInfo{
					RequestId:   meta,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
					ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  err.Error(),
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.Config.SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Insert_Config{SkipStrictExistCheck: true}
		}
	}
	uuid := fuid.String()
	req.Vector.Id = uuid
	if req.GetConfig().GetTimestamp() == 0 {
		now := time.Now().UnixNano()
		if req.GetConfig() == nil {
			req.Config = &payload.Insert_Config{
				Timestamp: now,
			}
		} else {
			req.Config.Timestamp = now
		}
	}
	loc, err = s.gateway.Insert(ctx, req, s.copts...)
	if err != nil {
		err = status.WrapWithInternal("Insert API failed to insert next gateway", err,
			&errdetails.RequestInfo{
				RequestId:   meta,
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
	err = s.metadata.SetUUIDandMeta(ctx, uuid, meta)
	if err != nil {
		err = status.WrapWithInternal("Insert API failed to store metadata", err,
			&errdetails.RequestInfo{
				RequestId:   meta,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.SetUUIDandMeta",
				ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)

		_, rerr := s.gateway.Remove(ctx, &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		if rerr != nil {
			err = errors.Wrap(err, status.WrapWithInternal("Insert API failed to remove inserted data caused by metadata store failure", rerr,
				&errdetails.RequestInfo{
					RequestId:   meta,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Remove",
					ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
					Owner:        errdetails.ValdResourceOwner,
					Description:  rerr.Error(),
				}, info.Get()).Error())
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
	vecs := reqs.GetRequests()
	metaMap := make(map[string]string, len(vecs))
	metas := make([]string, 0, len(vecs))
	now := time.Now().UnixNano()
	for i, req := range vecs {
		vec := req.GetVector()
		vl := len(vec.GetVector())
		meta := vec.GetId()
		if vl < algorithm.MinimumVectorDimensionSize {
			err = errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument("MultiInsert API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   meta,
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       fmt.Sprintf("vector dimension size for id: %s", meta),
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
			exists, err := s.metadata.Exists(ctx, meta)
			if err != nil {
				err = status.WrapWithInternal("MultiInsert API failed to check metadata exsisting or not", err,
					&errdetails.RequestInfo{
						RequestId:   meta,
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.GetMetaInverse",
						ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
						Owner:        errdetails.ValdResourceOwner,
						Description:  err.Error(),
					}, info.Get())
				log.Error(err)
				if span != nil {
					span.SetStatus(trace.StatusCodeInternal(err.Error()))
				}
				return nil, err
			}
			if exists {
				err = errors.ErrMetaDataAlreadyExists(meta)
				err = status.WrapWithAlreadyExists(fmt.Sprintf("MultiInsert API ID = %v already exists", meta), err,
					&errdetails.RequestInfo{
						RequestId:   meta,
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Exists",
						ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
						Owner:        errdetails.ValdResourceOwner,
						Description:  err.Error(),
					}, info.Get())
				if span != nil {
					span.SetStatus(trace.StatusCodeAlreadyExists(err.Error()))
				}
				return nil, err
			}
			if req.GetConfig() != nil {
				reqs.Requests[i].Config.SkipStrictExistCheck = true
			} else {
				reqs.Requests[i].Config = &payload.Insert_Config{SkipStrictExistCheck: true}
			}
		}

		uuid := fuid.String()
		metaMap[uuid] = meta
		metas = append(metas, meta)
		reqs.Requests[i].Vector.Id = uuid
		if reqs.Requests[i].GetConfig().GetTimestamp() == 0 {
			if reqs.Requests[i].GetConfig() == nil {
				reqs.Requests[i].Config = &payload.Insert_Config{
					Timestamp: now,
				}
			} else {
				reqs.Requests[i].Config.Timestamp = now
			}
		}
	}

	res, err = s.gateway.MultiInsert(ctx, reqs, s.copts...)
	if err != nil {
		err = status.WrapWithInternal("MultiInsert API failed to insert to next gateway", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(metas, ","),
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

	err = s.metadata.SetUUIDandMetas(ctx, metaMap)
	if err != nil {
		err = status.WrapWithInternal("MultiInsert API failed to store metadata to meta component", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(metas, ","),
				ServingData: errdetails.Serialize(metaMap),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.SetUUIDandMetas",
				ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		log.Error(err)
		rmr := &payload.Remove_MultiRequest{
			Requests: make([]*payload.Remove_Request, 0, len(reqs.GetRequests())),
		}
		for _, req := range reqs.GetRequests() {
			rmr.Requests = append(rmr.Requests, &payload.Remove_Request{
				Id: &payload.Object_ID{
					Id: req.GetVector().GetId(),
				},
			})
		}
		_, rerr := s.gateway.MultiRemove(ctx, rmr, s.copts...)
		if rerr != nil {
			err = errors.Wrap(err, status.WrapWithInternal("MultiInsert API failed to remove miss inserted data caused by failed to store meta component", rerr,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(metas, ","),
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
	meta := req.GetVector().GetId()
	vl := len(req.GetVector().GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument("Update API invalid vector argument", err,
			&errdetails.RequestInfo{
				RequestId:   meta,
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
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil || len(uuid) == 0 {
		if err == nil {
			err = errors.ErrObjectIDNotFound(meta)
		}
		err = status.WrapWithNotFound(fmt.Sprintf("Update API ID = %v not found", uuid), err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.GetUUID",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}
	req.Vector.Id = uuid
	if req.GetConfig().GetTimestamp() == 0 {
		now := time.Now().UnixNano()
		if req.GetConfig() == nil {
			req.Config = &payload.Update_Config{
				Timestamp: now,
			}
		} else {
			req.Config.Timestamp = now
		}
	}
	res, err = s.gateway.Update(ctx, req, s.copts...)
	if err != nil {
		err = status.WrapWithInternal("Update API failed to update next gateway", err,
			&errdetails.RequestInfo{
				RequestId:   meta,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.Update",
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
	ids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		vl := len(req.GetVector().GetVector())
		if vl < algorithm.MinimumVectorDimensionSize {
			err = errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument("MultiUpdate API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
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
		ids = append(ids, req.GetVector().GetId())
	}

	now := time.Now().UnixNano()
	uuids, err := s.metadata.GetUUIDs(ctx, ids...)
	if err != nil {
		err = status.WrapWithNotFound(fmt.Sprintf("MultiUpdate API ID = %v not found", uuids), err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(uuids, ", "),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.GetUUIDs",
				ResourceName: strings.Join(s.gateway.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}
	for i, uuid := range uuids {
		reqs.Requests[i].Vector.Id = uuid
		if reqs.Requests[i].GetConfig().GetTimestamp() == 0 {
			if reqs.GetRequests()[i].GetConfig() == nil {
				reqs.Requests[i].Config = &payload.Update_Config{
					Timestamp: now,
				}
			} else {
				reqs.Requests[i].Config.Timestamp = now
			}
		}
	}
	res, err = s.gateway.MultiUpdate(ctx, reqs)
	if err != nil {
		err = status.WrapWithInternal("MultiUpdate API failed to update next gateway", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(uuids, ", "),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1.MultiUpdate",
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

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Upsert")
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
		err = status.WrapWithInvalidArgument("Upsert API invalid vector argument", err,
			&errdetails.RequestInfo{
				RequestId:   meta,
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

	var now int64
	if req.GetConfig().GetTimestamp() > 0 {
		now = req.GetConfig().GetTimestamp()
	} else {
		now = time.Now().UnixNano()
	}
	filters := req.GetConfig().GetFilters()
	exists, err := s.metadata.Exists(ctx, meta)
	if err != nil {
		log.Debugf("Upsert API metadata exists check error:\t%s", err.Error())
	}
	var operation string
	if err != nil || !exists {
		operation = "Insert"
		loc, err = s.Insert(ctx, &payload.Insert_Request{
			Vector: vec,
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Filters:              filters,
				Timestamp:            now,
			},
		})
	} else {
		operation = "Update"
		loc, err = s.Update(ctx, &payload.Update_Request{
			Vector: vec,
			Config: &payload.Update_Config{
				SkipStrictExistCheck: true,
				Filters:              filters,
				Timestamp:            now,
			},
		})
	}
	if err != nil {
		err = status.WrapWithInternal("Upsert API failed to "+operation+" next gateway", err,
			&errdetails.RequestInfo{
				RequestId:   meta,
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
	now := time.Now().UnixNano()
	for _, req := range reqs.GetRequests() {
		vec := req.GetVector()
		meta := vec.GetId()
		vl := len(vec.GetVector())
		if vl < algorithm.MinimumVectorDimensionSize {
			err = errors.ErrInvalidDimensionSize(vl, 0)
			err = status.WrapWithInvalidArgument("MultiUpsert API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   meta,
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
		if req.GetConfig().GetTimestamp() == 0 {
			req.Config.Timestamp = now
		}
		ids = append(ids, meta)
		filters := req.GetConfig().GetFilters()
		exists, err := s.metadata.Exists(ctx, meta)
		if !exists || err != nil {
			insertReqs = append(insertReqs, &payload.Insert_Request{
				Vector: vec,
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
					Filters:              filters,
					Timestamp:            req.GetConfig().GetTimestamp(),
				},
			})
		} else {
			updateReqs = append(updateReqs, &payload.Update_Request{
				Vector: vec,
				Config: &payload.Update_Config{
					SkipStrictExistCheck: true,
					Filters:              filters,
					Timestamp:            req.GetConfig().GetTimestamp(),
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
	now := time.Now().UnixNano()
	meta := req.GetId().GetId()
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil {
		err = status.WrapWithNotFound(fmt.Sprintf("Remove API ID = %v not found", meta), err,
			&errdetails.RequestInfo{
				RequestId:   meta,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.GetUUID",
				ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}

	if req.GetConfig() != nil {
		req.Config.SkipStrictExistCheck = true
	} else {
		req.Config = &payload.Remove_Config{SkipStrictExistCheck: true}
	}

	if req.GetId() != nil {
		req.Id.Id = uuid
	} else {
		req.Id = &payload.Object_ID{Id: uuid}
	}
	if req.GetConfig().GetTimestamp() == 0 {
		req.Config.Timestamp = now
	}
	loc, err = s.gateway.Remove(ctx, req, s.copts...)
	if err != nil {
		err = status.WrapWithInternal("Remove API failed to remove request", err,
			&errdetails.RequestInfo{
				RequestId:   meta,
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
	_, err = s.metadata.DeleteMeta(ctx, uuid)
	if err != nil {
		err = status.WrapWithInternal("Remove API failed to remove from meta store", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.DeleteMeta",
				ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
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
	for _, req := range reqs.GetRequests() {
		ids = append(ids, req.GetId().GetId())
	}
	now := time.Now().UnixNano()
	uuids, err := s.metadata.GetUUIDs(ctx, ids...)
	if err != nil {
		err = status.WrapWithNotFound(fmt.Sprintf("MultiRemove API ID = %v not found", ids), err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ", "),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.GetUUIDs",
				ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}
	for i, id := range uuids {
		if reqs.Requests[i].GetId() != nil {
			reqs.Requests[i].Id.Id = id
		} else {
			reqs.Requests[i].Id = &payload.Object_ID{Id: id}
		}
		req := reqs.Requests[i]
		if req.GetConfig() != nil {
			reqs.Requests[i].Config.SkipStrictExistCheck = true
		} else {
			reqs.Requests[i].Config = &payload.Remove_Config{SkipStrictExistCheck: true}
		}

		if req.GetId() != nil {
			reqs.Requests[i].Id.Id = id
		} else {
			reqs.Requests[i].Id = &payload.Object_ID{Id: id}
		}
		if req.GetConfig().GetTimestamp() == 0 {
			reqs.Requests[i].Config.Timestamp = now
		}
	}
	locs, err = s.gateway.MultiRemove(ctx, reqs, s.copts...)
	if err != nil {
		err = status.WrapWithInternal("MultiRemove API failed to remove request", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ", "),
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
	_, err = s.metadata.DeleteMetas(ctx, uuids...)
	if err != nil {
		err = status.WrapWithInternal("MultiRemove API failed to remove from meta store", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ", "),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.DeleteMetas",
				ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
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
	meta := req.GetId().GetId()
	uuid, err := s.metadata.GetUUID(ctx, meta)
	if err != nil {
		err = status.WrapWithNotFound(fmt.Sprintf("GetObject API ID = %v not found", meta), err,
			&errdetails.RequestInfo{
				RequestId:   meta,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/meta.v1.GetUUID",
				ResourceName: strings.Join(s.metadata.GRPCClient().ConnectedAddrs(), ", "),
				Owner:        errdetails.ValdResourceOwner,
				Description:  err.Error(),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}
	if req.GetId() != nil {
		req.Id.Id = uuid
	} else {
		req.Id = &payload.Object_ID{Id: uuid}
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
	vec.Id = meta
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
