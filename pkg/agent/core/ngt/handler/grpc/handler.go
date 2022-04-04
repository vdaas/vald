//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
	"reflect"
	"strconv"
	"sync"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
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
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

type Server interface {
	agent.AgentServer
	vald.Server
}

type server struct {
	name              string
	ip                string
	ngt               service.NGT
	eg                errgroup.Group
	streamConcurrency int
	agent.UnimplementedAgentServer
	vald.UnimplementedValdServer
}

const (
	apiName         = "vald/agent/core/ngt"
	ngtResourceType = "vald/internal/core/algorithm"
)

var errNGT = new(errors.NGTError)

func New(opts ...Option) (Server, error) {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(s); err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))

			e := new(errors.ErrCriticalOption)
			if errors.As(err, &e) {
				log.Error(werr)
				return nil, werr
			}
			log.Warn(werr)
		}
	}
	return s, nil
}

func (s *server) newLocations(uuids ...string) (locs *payload.Object_Locations) {
	if len(uuids) == 0 {
		return nil
	}
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, 0, len(uuids)),
	}
	for _, uuid := range uuids {
		locs.Locations = append(locs.GetLocations(), &payload.Object_Location{
			Name: s.name,
			Uuid: uuid,
			Ips:  []string{s.ip},
		})
	}
	return locs
}

func (s *server) newLocation(uuid string) *payload.Object_Location {
	locs := s.newLocations(uuid)
	if locs != nil && locs.GetLocations() != nil && len(locs.GetLocations()) > 0 {
		return locs.Locations[0]
	}
	return nil
}

func (s *server) Exists(ctx context.Context, uid *payload.Object_ID) (res *payload.Object_ID, err error) {
	_, span := trace.StartSpan(ctx, apiName+".Exists")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := uid.GetId()
	oid, ok := s.ngt.Exists(uuid)
	if !ok {
		err = errors.ErrObjectIDNotFound(uid.GetId())
		err = status.WrapWithNotFound(fmt.Sprintf("Exists API meta %s's uuid not found", uid.GetId()), err,
			&errdetails.RequestInfo{
				RequestId:   uid.GetId(),
				ServingData: errdetails.Serialize(uid),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.Exists",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
			uid.GetId(), info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}
	return &payload.Object_ID{
		Id: strconv.Itoa(int(oid)),
	}, nil
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	_, span := trace.StartSpan(ctx, apiName+".Search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(req.GetVector()) != s.ngt.GetDimensionSize() {
		err = errors.ErrIncompatibleDimensionSize(len(req.GetVector()), int(s.ngt.GetDimensionSize()))
		err = status.WrapWithInvalidArgument("Search API Incompatible Dimension Size detected",
			err,
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
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.Search",
			})
		log.Warn(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	res, err = toSearchResponse(
		s.ngt.Search(
			req.GetVector(),
			req.GetConfig().GetNum(),
			req.GetConfig().GetEpsilon(),
			req.GetConfig().GetRadius()))
	if err != nil || res == nil {
		var stat trace.Status
		switch {
		case errors.Is(err, errors.ErrCreateIndexingIsInProgress):
			err = status.WrapWithAborted("Search API aborted to process search request due to createing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Search",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			stat = trace.StatusCodeAborted(err.Error())
		case errors.Is(err, errors.ErrEmptySearchResult),
			err == nil && res == nil,
			0 < req.GetConfig().GetMinNum() && len(res.GetResults()) < int(req.GetConfig().GetMinNum()):
			err = status.WrapWithNotFound(fmt.Sprintf("Search API requestID %s's search result not found", req.GetConfig().GetRequestId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Search",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			stat = trace.StatusCodeNotFound(err.Error())
		case errors.As(err, &errNGT):
			log.Errorf("ngt core process returned error: %v", err)
			err = status.WrapWithInternal("Search API failed to process search request due to ngt core process returned error", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Search/core.ngt",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			stat = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrIncompatibleDimensionSize(len(req.GetVector()), int(s.ngt.GetDimensionSize()))):
			err = status.WrapWithInvalidArgument("Search API Incompatible Dimension Size detected",
				err,
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
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Search",
				})
			log.Warn(err)
			stat = trace.StatusCodeInvalidArgument(err.Error())
		default:
			err = status.WrapWithInternal("Search API failed to process search request", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Search",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			stat = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.SetStatus(stat)
		}
		return nil, err
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (res *payload.Search_Response, err error) {
	_, span := trace.StartSpan(ctx, apiName+".SearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec, dst, err := s.ngt.SearchByID(
		req.GetId(),
		req.GetConfig().GetNum(),
		req.GetConfig().GetEpsilon(),
		req.GetConfig().GetRadius())
	res, err = toSearchResponse(dst, err)
	if err != nil || res == nil {
		var stat trace.Status
		switch {
		case errors.Is(err, errors.ErrCreateIndexingIsInProgress):
			err = status.WrapWithAborted("SearchByID API aborted to process search request due to createing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.SearchByID",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			stat = trace.StatusCodeAborted(err.Error())
		case errors.Is(err, errors.ErrEmptySearchResult),
			err == nil && res == nil,
			0 < req.GetConfig().GetMinNum() && len(res.GetResults()) < int(req.GetConfig().GetMinNum()):
			err = status.WrapWithNotFound(fmt.Sprintf("SearchByID API uuid %s's search result not found", req.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.SearchByID",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			stat = trace.StatusCodeNotFound(err.Error())
		case errors.Is(err, errors.ErrObjectIDNotFound(req.GetId())),
			strings.Contains(err.Error(), fmt.Sprintf("ngt uuid %s's object not found", req.GetId())):
			err = status.WrapWithNotFound(fmt.Sprintf("SearchByID API uuid %s's object not found", req.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.SearchByID",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			stat = trace.StatusCodeNotFound(err.Error())
		case errors.As(err, &errNGT):
			log.Errorf("ngt core process returned error: %v", err)
			err = status.WrapWithInternal("SearchByID API failed to process search request due to ngt core process returned error", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.SearchByID/core.ngt",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			stat = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrIncompatibleDimensionSize(len(vec), int(s.ngt.GetDimensionSize()))):
			err = status.WrapWithInvalidArgument("SearchByID API Incompatible Dimension Size detected",
				err,
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
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.SearchByID",
				})
			log.Warn(err)
			stat = trace.StatusCodeInvalidArgument(err.Error())
		default:
			err = status.WrapWithInternal("SearchByID API failed to process search request", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.SearchByID",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			stat = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.SetStatus(stat)
		}
		return nil, err
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func toSearchResponse(dists []model.Distance, err error) (res *payload.Search_Response, rerr error) {
	if err != nil {
		return nil, err
	}
	if len(dists) == 0 {
		return nil, errors.ErrEmptySearchResult
	}
	res = new(payload.Search_Response)
	res.Results = make([]*payload.Object_Distance, 0, len(dists))
	for _, dist := range dists {
		res.Results = append(res.GetResults(), &payload.Object_Distance{
			Id:       dist.ID,
			Distance: dist.Distance,
		})
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
			req := data.(*payload.Search_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamSearch/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Search(ctx, data.(*payload.Search_Request))
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Search gRPC error response")
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse StreamSearch gRPC error response")
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
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
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse SearchByID gRPC error response")
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamSearchByID gRPC error response")
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
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
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(func() (err error) {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s.MultiSearch/errgroup.Go/id-%d", apiName, idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.Search(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal,
					"failed to parse Search gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse MultiSearch gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.MultiSearch",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(func() error {
			ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s.MultiSearchByID/errgroup.Go/id-%d", apiName, idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			defer wg.Done()
			r, err := s.SearchByID(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal,
					"failed to parse SearchByID gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse MultiSearchByID gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.MultiSearchByID",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) LinearSearch(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	_, span := trace.StartSpan(ctx, apiName+".LinearSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(req.GetVector()) != s.ngt.GetDimensionSize() {
		err = errors.ErrIncompatibleDimensionSize(len(req.GetVector()), int(s.ngt.GetDimensionSize()))
		err = status.WrapWithInvalidArgument("LinearSearch API Incompatible Dimension Size detected",
			err,
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
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.LinearSearch",
			})
		log.Warn(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	res, err = toSearchResponse(
		s.ngt.LinearSearch(
			req.GetVector(),
			req.GetConfig().GetNum()))
	if err != nil || res == nil {
		var stat trace.Status
		switch {
		case errors.Is(err, errors.ErrCreateIndexingIsInProgress):
			err = status.WrapWithAborted("LinearSearch API aborted to process search request due to createing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.LinearSearch",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			stat = trace.StatusCodeAborted(err.Error())
		case errors.Is(err, errors.ErrEmptySearchResult),
			err == nil && res == nil,
			0 < req.GetConfig().GetMinNum() && len(res.GetResults()) < int(req.GetConfig().GetMinNum()):
			err = status.WrapWithNotFound(fmt.Sprintf("LinearSearch API requestID %s's search result not found", req.GetConfig().GetRequestId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.LinearSearch",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			stat = trace.StatusCodeNotFound(err.Error())
		case errors.As(err, &errNGT):
			log.Errorf("ngt core process returned error: %v", err)
			err = status.WrapWithInternal("LinearSearch API failed to process search request due to ngt core process returned error", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.LinearSearch/core.ngt",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			stat = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrIncompatibleDimensionSize(len(req.GetVector()), int(s.ngt.GetDimensionSize()))):
			err = status.WrapWithInvalidArgument("LinearSearch API Incompatible Dimension Size detected",
				err,
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
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.LinearSearch",
				})
			log.Warn(err)
			stat = trace.StatusCodeInvalidArgument(err.Error())
		default:
			err = status.WrapWithInternal("LinearSearch API failed to process search request", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.LinearSearch",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			stat = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.SetStatus(stat)
		}
		return nil, err
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func (s *server) LinearSearchByID(ctx context.Context, req *payload.Search_IDRequest) (res *payload.Search_Response, err error) {
	_, span := trace.StartSpan(ctx, apiName+".LinearSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec, dst, err := s.ngt.LinearSearchByID(
		req.GetId(),
		req.GetConfig().GetNum())
	res, err = toSearchResponse(dst, err)
	if err != nil || res == nil {
		var stat trace.Status
		switch {
		case errors.Is(err, errors.ErrCreateIndexingIsInProgress):
			err = status.WrapWithAborted("LinearSearchByID API aborted to process search request due to createing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.LinearSearchByID",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			stat = trace.StatusCodeAborted(err.Error())
		case errors.Is(err, errors.ErrEmptySearchResult),
			err == nil && res == nil,
			0 < req.GetConfig().GetMinNum() && len(res.GetResults()) < int(req.GetConfig().GetMinNum()):
			err = status.WrapWithNotFound(fmt.Sprintf("LinearSearchByID API uuid %s's search result not found", req.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.LinearSearchByID",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			stat = trace.StatusCodeNotFound(err.Error())
		case errors.Is(err, errors.ErrObjectIDNotFound(req.GetId())),
			strings.Contains(err.Error(), fmt.Sprintf("ngt uuid %s's object not found", req.GetId())):
			err = status.WrapWithNotFound(fmt.Sprintf("LinearSearchByID API uuid %s's object not found", req.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.LinearSearchByID",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			stat = trace.StatusCodeNotFound(err.Error())
		case errors.As(err, &errNGT):
			log.Errorf("ngt core process returned error: %v", err)
			err = status.WrapWithInternal("LinearSearchByID API failed to process search request due to ngt core process returned error", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.LinearSearchByID/core.ngt",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			stat = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrIncompatibleDimensionSize(len(vec), int(s.ngt.GetDimensionSize()))):
			err = status.WrapWithInvalidArgument("LinearSearchByID API Incompatible Dimension Size detected",
				err,
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
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.LinearSearchByID",
				})
			log.Warn(err)
			stat = trace.StatusCodeInvalidArgument(err.Error())
		default:
			err = status.WrapWithInternal("LinearSearchByID API failed to process search request", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.LinearSearchByID",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			stat = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.SetStatus(stat)
		}
		return nil, err
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func (s *server) StreamLinearSearch(stream vald.Search_StreamLinearSearchServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamLinearSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamLinearSearch/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearch(ctx, data.(*payload.Search_Request))
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse LinearSearch gRPC error response")
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse StreamLinearSearch gRPC error response")
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
		return err
	}
	return nil
}

func (s *server) StreamLinearSearchByID(stream vald.Search_StreamLinearSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+".StreamLinearSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_IDRequest)
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamLinearSearchByID/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearchByID(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse LinearSearchByID gRPC error response")
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamLinearSearchByID gRPC error response")
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
		return err
	}
	return nil
}

func (s *server) MultiLinearSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiLinearSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(func() (err error) {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s.MultiLinearSearch/errgroup.Go/id-%d", apiName, idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.LinearSearch(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal,
					"failed to parse LinearSearch gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse MultiLinearSearch gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.MultiLinearSearch",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) MultiLinearSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiLinearSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	rids := make([]string, 0, len(reqs.GetRequests()))
	for i, req := range reqs.Requests {
		idx, query := i, req
		rids = append(rids, req.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(func() error {
			ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s.MultiLinearSearchByID/errgroup.Go/id-%d", apiName, idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			defer wg.Done()
			r, err := s.LinearSearchByID(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal,
					"failed to parse LinearSearchByID gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(errs, codes.Internal,
			"failed to parse MultiLinearSearchByID gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.MultiLinearSearchByID",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
		return nil, err
	}
	return res, nil
}

// Insert inserts a vector to the NGT.
func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (res *payload.Object_Location, err error) {
	_, span := trace.StartSpan(ctx, apiName+".Insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	if len(vec.GetVector()) != s.ngt.GetDimensionSize() {
		err = errors.ErrIncompatibleDimensionSize(len(vec.GetVector()), int(s.ngt.GetDimensionSize()))
		err = status.WrapWithInvalidArgument("Insert API Incompatible Dimension Size detected",
			err,
			&errdetails.RequestInfo{
				RequestId:   vec.GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.Insert",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		log.Warn(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}

	err = s.ngt.InsertWithTime(vec.GetId(), vec.GetVector(), req.GetConfig().GetTimestamp())
	if err != nil {
		var code trace.Status

		if errors.Is(err, errors.ErrUUIDAlreadyExists(vec.GetId())) {
			err = status.WrapWithAlreadyExists(fmt.Sprintf("Insert API uuid %s already exists", vec.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Insert",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeAlreadyExists(err.Error())
		} else if errors.Is(err, errors.ErrUUIDNotFound(0)) {
			err = status.WrapWithInvalidArgument(fmt.Sprintf("Insert API empty uuid \"%s\" was given", vec.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "uuid",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Insert",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeInvalidArgument(err.Error())
		} else {
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse Insert gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Insert",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			code = trace.FromGRPCStatus(st.Code(), msg)
		}
		if span != nil {
			span.SetStatus(code)
		}
		return nil, err
	}
	return s.newLocation(vec.GetId()), nil
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
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Insert gRPC error response")
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamInsert gRPC error response")
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
		return err
	}
	return nil
}

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (res *payload.Object_Locations, err error) {
	_, span := trace.StartSpan(ctx, apiName+".MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuids := make([]string, 0, len(reqs.GetRequests()))
	vmap := make(map[string][]float32, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		vec := req.GetVector()
		if len(vec.GetVector()) != s.ngt.GetDimensionSize() {
			err = errors.ErrIncompatibleDimensionSize(len(vec.GetVector()), int(s.ngt.GetDimensionSize()))
			err = status.WrapWithInvalidArgument("MultiInsert API Incompatible Dimension Size detected",
				err,
				&errdetails.RequestInfo{
					RequestId:   vec.GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vector dimension size",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiInsert",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
			}
			return nil, err
		}
		vmap[vec.GetId()] = vec.GetVector()
		uuids = append(uuids, vec.GetId())
	}
	err = s.ngt.InsertMultiple(vmap)
	if err != nil {
		var code trace.Status
		if alreadyExistsIDs := func() []string {
			aids := make([]string, 0, len(uuids))
			for _, id := range uuids {
				if errors.Is(err, errors.ErrUUIDAlreadyExists(id)) {
					aids = append(aids, id)
				}
			}
			return aids
		}(); len(alreadyExistsIDs) != 0 {
			err = status.WrapWithAlreadyExists(fmt.Sprintf("MultiInsert API uuids %v already exists", alreadyExistsIDs), err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiInsert",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeAlreadyExists(err.Error())
		} else if errors.Is(err, errors.ErrUUIDNotFound(0)) {
			err = status.WrapWithInvalidArgument(fmt.Sprintf("MultiInsert API invalid uuids \"%v\" detected", uuids), err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "uuid",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiInsert",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeInvalidArgument(err.Error())
		} else {
			err = status.WrapWithInternal("MultiInsert API failed", err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiInsert",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			code = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.SetStatus(code)
		}
		return nil, err
	}
	return s.newLocations(uuids...), nil
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (res *payload.Object_Location, err error) {
	_, span := trace.StartSpan(ctx, apiName+".Update")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	if len(vec.GetVector()) != s.ngt.GetDimensionSize() {
		err = errors.ErrIncompatibleDimensionSize(len(vec.GetVector()), int(s.ngt.GetDimensionSize()))
		err = status.WrapWithInvalidArgument("Update API Incompatible Dimension Size detected",
			err,
			&errdetails.RequestInfo{
				RequestId:   vec.GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.Update",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		log.Warn(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}
	err = s.ngt.UpdateWithTime(vec.GetId(), vec.GetVector(), req.GetConfig().GetTimestamp())
	if err != nil {
		var code trace.Status
		if errors.Is(err, errors.ErrObjectIDNotFound(vec.GetId())) {
			err = status.WrapWithNotFound(fmt.Sprintf("Update API uuid %s not found", vec.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Update",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeNotFound(err.Error())
		} else if errors.Is(err, errors.ErrUUIDNotFound(0)) || errors.Is(err, errors.ErrInvalidDimensionSize(len(vec.GetVector()), s.ngt.GetDimensionSize())) {
			err = status.WrapWithInvalidArgument(fmt.Sprintf("Update API invalid argument for uuid \"%s\" vec \"%v\" detected", vec.GetId(), vec.GetVector()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "uuid or vector",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Update",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeInvalidArgument(err.Error())
		} else if errors.Is(err, errors.ErrUUIDAlreadyExists(vec.GetId())) {
			err = status.WrapWithAlreadyExists(fmt.Sprintf("Update API uuid %s's same data already exists", vec.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Update",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeAlreadyExists(err.Error())
		} else {
			err = status.WrapWithInternal("Update API failed", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Update",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			code = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.SetStatus(code)
		}
		return nil, err
	}
	return s.newLocation(vec.GetId()), nil
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
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Update gRPC error response")
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamUpdate gRPC error response")
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
		return err
	}
	return nil
}

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (res *payload.Object_Locations, err error) {
	_, span := trace.StartSpan(ctx, apiName+".MultiUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	uuids := make([]string, 0, len(reqs.GetRequests()))
	vmap := make(map[string][]float32, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		vec := req.GetVector()
		if len(vec.GetVector()) != s.ngt.GetDimensionSize() {
			err = errors.ErrIncompatibleDimensionSize(len(vec.GetVector()), int(s.ngt.GetDimensionSize()))
			err = status.WrapWithInvalidArgument("MultiUpdate API Incompatible Dimension Size detected",
				err,
				&errdetails.RequestInfo{
					RequestId:   vec.GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vector dimension size",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiUpdate",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
			}
			return nil, err
		}
		vmap[vec.GetId()] = vec.GetVector()
		uuids = append(uuids, vec.GetId())
	}

	err = s.ngt.UpdateMultiple(vmap)
	if err != nil {
		var code trace.Status
		if notFoundIDs := func() []string {
			aids := make([]string, 0, len(uuids))
			for _, id := range uuids {
				if errors.Is(err, errors.ErrObjectIDNotFound(id)) {
					aids = append(aids, id)
				}
			}
			return aids
		}(); len(notFoundIDs) != 0 {
			err = status.WrapWithNotFound(fmt.Sprintf("MultiUpdate API uuids %v not found", notFoundIDs), err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiUpdate",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeNotFound(err.Error())
		} else if invalidDimensionIDs := func() []string {
			idis := make([]string, 0, len(uuids))
			for id, vec := range vmap {
				if errors.Is(err, errors.ErrInvalidDimensionSize(len(vec), s.ngt.GetDimensionSize())) {
					idis = append(idis, id)
				}
			}
			return idis
		}(); len(invalidDimensionIDs) != 0 || errors.Is(err, errors.ErrUUIDNotFound(0)) {
			err = status.WrapWithInvalidArgument(fmt.Sprintf("MultiUpdate API invalid argument for uuids \"%v\" detected", invalidDimensionIDs), err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "uuid or vector",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiUpdate",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeInvalidArgument(err.Error())
		} else if alreadyExistsIDs := func() []string {
			aids := make([]string, 0, len(uuids))
			for _, id := range uuids {
				if errors.Is(err, errors.ErrUUIDAlreadyExists(id)) {
					aids = append(aids, id)
				}
			}
			return aids
		}(); len(alreadyExistsIDs) != 0 {
			err = status.WrapWithAlreadyExists(fmt.Sprintf("MultiUpdate API uuids %v already exists", alreadyExistsIDs), err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiUpdate",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeAlreadyExists(err.Error())
		} else {
			err = status.WrapWithInternal("Update API failed", err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Update",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			code = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.SetStatus(code)
		}
		return nil, err
	}
	return s.newLocations(uuids...), nil
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	if len(vec.GetVector()) != s.ngt.GetDimensionSize() {
		err = errors.ErrIncompatibleDimensionSize(len(vec.GetVector()), int(s.ngt.GetDimensionSize()))
		err = status.WrapWithInvalidArgument("Upsert API Incompatible Dimension Size detected",
			err,
			&errdetails.RequestInfo{
				RequestId:   vec.GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.Upsert",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		log.Warn(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
		}
		return nil, err
	}

	rtName := "/ngt.Upsert"
	_, exists := s.ngt.Exists(req.GetVector().GetId())
	if exists {
		loc, err = s.Update(ctx, &payload.Update_Request{
			Vector: req.GetVector(),
			Config: &payload.Update_Config{
				Timestamp:            req.GetConfig().GetTimestamp(),
				SkipStrictExistCheck: true,
			},
		})
		rtName += "/ngt.Update"
	} else {
		loc, err = s.Insert(ctx, &payload.Insert_Request{
			Vector: req.GetVector(),
			Config: &payload.Insert_Config{
				Timestamp:            req.GetConfig().GetTimestamp(),
				SkipStrictExistCheck: true,
			},
		})
		rtName += "/ngt.Insert"
	}
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Upsert gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + rtName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Upsert gRPC error response")
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamUpsert gRPC error response")
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
		return err
	}
	return nil
}

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (res *payload.Object_Locations, err error) {
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
		if len(vec.GetVector()) != s.ngt.GetDimensionSize() {
			err = errors.ErrIncompatibleDimensionSize(len(vec.GetVector()), int(s.ngt.GetDimensionSize()))
			err = status.WrapWithInvalidArgument("MultiUpsert API Incompatible Dimension Size detected",
				err,
				&errdetails.RequestInfo{
					RequestId:   vec.GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vector dimension size",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiUpsert",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			if span != nil {
				span.SetStatus(trace.StatusCodeInvalidArgument(err.Error()))
			}
			return nil, err
		}
		ids = append(ids, vec.GetId())
		_, exists := s.ngt.Exists(vec.GetId())
		if exists {
			updateReqs = append(updateReqs, &payload.Update_Request{
				Vector: vec,
				Config: &payload.Update_Config{
					Timestamp:            req.GetConfig().GetTimestamp(),
					SkipStrictExistCheck: true,
				},
			})
		} else {
			insertReqs = append(insertReqs, &payload.Insert_Request{
				Vector: vec,
				Config: &payload.Insert_Config{
					Timestamp:            req.GetConfig().GetTimestamp(),
					SkipStrictExistCheck: true,
				},
			})
		}
	}

	switch {
	case len(insertReqs) <= 0:
		res, err = s.MultiUpdate(ctx, &payload.Update_MultiRequest{
			Requests: updateReqs,
		})
	case len(updateReqs) <= 0:
		res, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
			Requests: insertReqs,
		})
	default:
		var (
			ures, ires *payload.Object_Locations
			errs       error
			mu         sync.Mutex
			wg         sync.WaitGroup
		)
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() (err error) {
			defer wg.Done()
			ures, err = s.MultiUpdate(ctx, &payload.Update_MultiRequest{
				Requests: updateReqs,
			})
			if err != nil {
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				mu.Unlock()
			}
			return nil
		}))
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() (err error) {
			defer wg.Done()
			ires, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
				Requests: insertReqs,
			})
			if err != nil {
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				mu.Unlock()
			}
			return nil
		}))
		wg.Wait()

		if errs == nil {
			var locs []*payload.Object_Location
			switch {
			case ures.GetLocations() == nil:
				locs = ires.GetLocations()
			case ires.GetLocations() == nil:
				locs = ures.GetLocations()
			default:
				locs = append(ures.GetLocations(), ires.GetLocations()...)
			}
			res = &payload.Object_Locations{
				Locations: locs,
			}
		} else {
			err = errs
		}

	}
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse MultiUpsert gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.MultiUpsert",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	id := req.GetId()
	uuid := id.GetId()
	err = s.ngt.DeleteWithTime(uuid, req.GetConfig().GetTimestamp())
	if err != nil {
		var code trace.Status
		if errors.Is(err, errors.ErrObjectIDNotFound(uuid)) {
			err = status.WrapWithNotFound(fmt.Sprintf("Remove API uuid %s not found", uuid), err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Remove",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeNotFound(err.Error())
		} else if errors.Is(err, errors.ErrUUIDNotFound(0)) {
			err = status.WrapWithInvalidArgument(fmt.Sprintf("Remove API invalid argument for uuid \"%s\" detected", uuid), err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "uuid",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Remove",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeInvalidArgument(err.Error())
		} else {
			err = status.WrapWithInternal("Remove API failed", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Remove",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			code = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.SetStatus(code)
		}
		return nil, err
	}
	return s.newLocation(uuid), nil
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
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Remove gRPC error response")
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamRemove gRPC error response")
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}
		return err
	}
	return nil
}

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		uuids = append(uuids, req.GetId().GetId())
	}
	err = s.ngt.DeleteMultiple(uuids...)
	if err != nil {
		var code trace.Status
		if notFoundIDs := func() []string {
			aids := make([]string, 0, len(uuids))
			for _, id := range uuids {
				if errors.Is(err, errors.ErrObjectIDNotFound(id)) {
					aids = append(aids, id)
				}
			}
			return aids
		}(); len(notFoundIDs) != 0 {
			err = status.WrapWithNotFound(fmt.Sprintf("MultiRemove API uuids %v not found", notFoundIDs), err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiRemove",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeNotFound(err.Error())
		} else if errors.Is(err, errors.ErrUUIDNotFound(0)) {
			err = status.WrapWithInvalidArgument(fmt.Sprintf("MultiRemove API invalid argument for uuids \"%v\" detected", uuids), err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "uuid",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiRemove",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			code = trace.StatusCodeInvalidArgument(err.Error())
		} else {
			err = status.WrapWithInternal("MultiRemove API failed", err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.MultiRemove",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			code = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.SetStatus(code)
		}
		return nil, err
	}
	return s.newLocations(uuids...), nil
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_VectorRequest) (res *payload.Object_Vector, err error) {
	_, span := trace.StartSpan(ctx, apiName+".GetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := id.GetId().GetId()
	vec, err := s.ngt.GetObject(uuid)
	if err != nil || vec == nil {
		err = errors.ErrObjectNotFound(err, uuid)
		err = status.WrapWithNotFound("GetObject API failed to remove request", err,
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(id),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.GetObject",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, err
	}

	return &payload.Object_Vector{
		Id:     uuid,
		Vector: vec,
	}, nil
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
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse GetObject gRPC error response")
				if sspan != nil {
					sspan.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamGetObject gRPC error response")
		if span != nil {
			span.SetStatus(trace.FromGRPCStatus(st.Code(), msg))
		}

		log.Error(err)
		return err
	}
	return nil
}

func (s *server) CreateIndex(ctx context.Context, c *payload.Control_CreateIndexRequest) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".CreateIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res = new(payload.Empty)
	err = s.ngt.CreateIndex(ctx, c.GetPoolSize())
	if err != nil {
		if errors.Is(err, errors.ErrUncommittedIndexNotFound) {
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("CreateIndex API failed to create indexes pool_size = %d", c.GetPoolSize()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(c),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.CreateIndex",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
				&errdetails.PreconditionFailure{
					Violations: []*errdetails.PreconditionFailureViolation{
						{
							Type:    "uncommited index is empty",
							Subject: "failed to CreateIndex operation caused by empty uncommited indices",
						},
					},
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		}
		log.Error(err)
		err = status.WrapWithInternal(fmt.Sprintf("CreateIndex API failed to create indexes pool_size = %d", c.GetPoolSize()), err,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(c),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.CreateIndex",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) SaveIndex(ctx context.Context, _ *payload.Empty) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".SaveIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res = new(payload.Empty)
	err = s.ngt.SaveIndex(ctx)
	if err != nil {
		log.Error(err)
		err = status.WrapWithInternal("SaveIndex API failed to save indices", err,
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.SaveIndex",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) CreateAndSaveIndex(ctx context.Context, c *payload.Control_CreateIndexRequest) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".CreateAndSaveIndex")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res = new(payload.Empty)
	err = s.ngt.CreateAndSaveIndex(ctx, c.GetPoolSize())
	if err != nil {
		if errors.Is(err, errors.ErrUncommittedIndexNotFound) {
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("CreateAndSaveIndex API failed to create indexes pool_size = %d", c.GetPoolSize()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(c),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.CreateAndSaveIndex",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
				&errdetails.PreconditionFailure{
					Violations: []*errdetails.PreconditionFailureViolation{
						{
							Type:    "uncommited index is empty",
							Subject: "failed to CreateAndSaveIndex operation caused by empty uncommited indices",
						},
					},
				}, info.Get())
			if span != nil {
				span.SetStatus(trace.StatusCodeFailedPrecondition(err.Error()))
			}
			return nil, err
		}
		err = status.WrapWithInternal(fmt.Sprintf("CreateAndSaveIndex API failed to create indexes pool_size = %d", c.GetPoolSize()), err,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(c),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.CreateAndSaveIndex",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, err
	}
	return res, nil
}

func (s *server) IndexInfo(ctx context.Context, _ *payload.Empty) (res *payload.Info_Index_Count, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".IndexInfo")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return &payload.Info_Index_Count{
		Stored:      uint32(s.ngt.Len()),
		Uncommitted: uint32(s.ngt.InsertVQueueBufferLen() + s.ngt.DeleteVQueueBufferLen()),
		Indexing:    s.ngt.IsIndexing(),
		Saving:      s.ngt.IsSaving(),
	}, nil
}
