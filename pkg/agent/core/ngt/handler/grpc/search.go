// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package grpc

import (
	"context"
	"fmt"
	"sync"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/pkg/agent/core/ngt/model"
	"go.opentelemetry.io/otel/attribute"
)

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.SearchRPCName)
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
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
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
		var attrs []attribute.KeyValue
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
			attrs = trace.StatusCodeAborted(err.Error())
		case errors.Is(err, errors.ErrFlushingIsInProgress):
			err = status.WrapWithAborted("Search API aborted to process search request due to flushing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Search",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			attrs = trace.StatusCodeAborted(err.Error())
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
			attrs = trace.StatusCodeNotFound(err.Error())
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
			attrs = trace.StatusCodeInternal(err.Error())
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
			attrs = trace.StatusCodeInvalidArgument(err.Error())
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
			attrs = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (res *payload.Search_Response, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.SearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetId()
	if len(uuid) == 0 {
		err = errors.ErrInvalidUUID(uuid)
		err = status.WrapWithInvalidArgument(fmt.Sprintf("SearchByID API invalid argument for uuid \"%s\" detected", uuid), err,
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
				ResourceType: ngtResourceType + "/ngt.SearchByID",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec, dst, err := s.ngt.SearchByID(
		uuid,
		req.GetConfig().GetNum(),
		req.GetConfig().GetEpsilon(),
		req.GetConfig().GetRadius())
	res, err = toSearchResponse(dst, err)
	if err != nil || res == nil {
		var attrs []attribute.KeyValue
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
			attrs = trace.StatusCodeAborted(err.Error())
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
			attrs = trace.StatusCodeNotFound(err.Error())
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
			attrs = trace.StatusCodeNotFound(err.Error())
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
			attrs = trace.StatusCodeInternal(err.Error())
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
			attrs = trace.StatusCodeInvalidArgument(err.Error())
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
			attrs = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	res.RequestId = req.GetConfig().GetRequestId()
	return res, nil
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Search(ctx, data.(*payload.Search_Request))
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Search gRPC error response")
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_IDRequest)
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamSearchByIDRPCName+"/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.SearchByID(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse SearchByID gRPC error response")
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiSearchRPCName)
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
			ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s/%s/errgroup.Go/id-%d", apiName, vald.MultiSearchRPCName, idx))
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
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) MultiSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiSearchByIDRPCName)
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
			ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s/%s/errgroup.Go/id-%d", apiName, vald.MultiSearchByIDRPCName, idx))
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
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
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
