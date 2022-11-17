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
	"go.opentelemetry.io/otel/attribute"
)

func (s *server) LinearSearch(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.LinearSearchRPCName)
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
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	res, err = toSearchResponse(
		s.ngt.LinearSearch(
			req.GetVector(),
			req.GetConfig().GetNum()))
	if err != nil || res == nil {
		var attrs []attribute.KeyValue
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
			attrs = trace.StatusCodeAborted(err.Error())
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
			attrs = trace.StatusCodeNotFound(err.Error())
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
			attrs = trace.StatusCodeInternal(err.Error())
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
			attrs = trace.StatusCodeInvalidArgument(err.Error())
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

func (s *server) LinearSearchByID(ctx context.Context, req *payload.Search_IDRequest) (res *payload.Search_Response, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.LinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	uuid := req.GetId()
	if len(uuid) == 0 {
		err = errors.ErrInvalidUUID(uuid)
		err = status.WrapWithInvalidArgument(fmt.Sprintf("LinearSearchByID API invalid argument for uuid \"%s\" detected", uuid), err,
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
				ResourceType: ngtResourceType + "/ngt.LinearSearchByID",
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
	vec, dst, err := s.ngt.LinearSearchByID(
		uuid,
		req.GetConfig().GetNum())
	res, err = toSearchResponse(dst, err)
	if err != nil || res == nil {
		var attrs []attribute.KeyValue
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
			attrs = trace.StatusCodeAborted(err.Error())
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
			attrs = trace.StatusCodeNotFound(err.Error())
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
			attrs = trace.StatusCodeNotFound(err.Error())
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
			attrs = trace.StatusCodeInternal(err.Error())
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
			attrs = trace.StatusCodeInvalidArgument(err.Error())
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

func (s *server) StreamLinearSearch(stream vald.Search_StreamLinearSearchServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamLinearSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearch(ctx, data.(*payload.Search_Request))
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse LinearSearch gRPC error response")
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
			"failed to parse StreamLinearSearch gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) StreamLinearSearchByID(stream vald.Search_StreamLinearSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_IDRequest)
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamLinearSearchByIDRPCName+"/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearchByID(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse LinearSearchByID gRPC error response")
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamLinearSearchByID gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiLinearSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiLinearSearchRPCName)
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
			ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s/%s/errgroup.Go/id-%d", apiName, vald.MultiLinearSearchRPCName, idx))
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
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) MultiLinearSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiLinearSearchByIDRPCName)
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
			ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s/%s/errgroup.Go/id-%d", apiName, vald.MultiLinearSearchByIDRPCName, idx))
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
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}
