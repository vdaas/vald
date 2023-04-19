//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"strconv"
	"sync"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/filter/egress"
	"github.com/vdaas/vald/internal/client/v1/client/filter/ingress"
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
)

type server struct {
	eg                errgroup.Group
	defaultVectorizer string
	defaultFilters    []string
	name              string
	ip                string
	ingress           ingress.Client
	egress            egress.Client
	gateway           client.Client
	copts             []grpc.CallOption
	streamConcurrency int
	Vectorizer        string
	DistanceFilters   []string
	ObjectFilters     []string
	SearchFilters     []string
	InsertFilters     []string
	UpdateFilters     []string
	UpsertFilters     []string
	vald.UnimplementedValdServerWithFilter
}

const apiName = "vald/gateway-filter"

func New(opts ...Option) vald.ServerWithFilter {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) SearchObject(ctx context.Context, req *payload.Search_ObjectRequest) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.SearchObjectRPCName), apiName+"/"+vald.SearchObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vr := req.GetVectorizer()
	if vr == nil || vr.GetPort() == 0 {
		err = errors.ErrFilterNotFound
		err = status.WrapWithInvalidArgument(vald.SearchObjectRPCName+" API vectorizer configuration is invalid", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer port",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchObjectRPCName,
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
	if vr.GetHost() == "" {
		vr.Host = "localhost"
	}
	target := fmt.Sprintf("%s:%d", vr.GetHost(), vr.GetPort())
	if len(target) == 0 {
		if len(s.Vectorizer) == 0 {
			err := errors.ErrFilterNotFound
			err = status.WrapWithInvalidArgument(vald.SearchObjectRPCName+" API vectorizer configuration is invalid", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchObjectRPCName,
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
		target = s.Vectorizer
	}
	c, err := s.ingress.Target(ctx, target)
	if err != nil {
		err = status.WrapWithUnavailable(
			fmt.Sprintf(vald.SearchRPCName+" API ingress filter targets %v not found", target),
			err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec, err := c.GenVector(ctx, &payload.Object_Blob{
		Object: req.GetObject(),
	})
	if err != nil {
		err = status.WrapWithInternal(vald.SearchObjectRPCName+" API failed to extract vector from filter", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return s.Search(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: req.GetConfig(),
	})
}

func (s *server) MultiSearchObject(ctx context.Context, reqs *payload.Search_MultiObjectRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiSearchObjectRPCName), apiName+"/"+vald.MultiSearchObjectRPCName)
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s.%s/errgroup.Go/id-%d", apiName, vald.MultiSearchObjectRPCName, idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.SearchObject(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.NotFound,
					vald.MultiSearchObjectRPCName+" API object "+string(query.GetObject())+"'s search request result not found",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequestFieldViolation{
							{
								Field:       "vectorizer targets",
								Description: err.Error(),
							},
						},
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get())
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf(vald.MultiSearchObjectRPCName+" API object %s's search request result not found",
							string(query.GetObject())), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf(vald.MultiSearchObjectRPCName+" API object %s's search request result not found",
								string(query.GetObject())), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		}))
	}
	wg.Wait()
	return res, errs
}

func (s *server) StreamSearchObject(stream vald.Filter_StreamSearchObjectServer) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamSearchObjectRPCName), apiName+"/"+vald.StreamSearchObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_ObjectRequest) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamSearchObjectRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()

			res, err := s.SearchObject(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchObjectRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
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
}

func (s *server) LinearSearchObject(ctx context.Context, req *payload.Search_ObjectRequest) (*payload.Search_Response, error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.LinearSearchObjectRPCName), apiName+"/"+vald.LinearSearchObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vr := req.GetVectorizer()
	if vr == nil || vr.GetPort() == 0 {
		err := errors.ErrInvalidAPIConfig
		err = status.WrapWithInvalidArgument(vald.LinearSearchObjectRPCName+" API vectorizer configuration is invalid", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer port",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchObjectRPCName,
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
	if vr.GetHost() == "" {
		vr.Host = "localhost"
	}
	target := fmt.Sprintf("%s:%d", vr.GetHost(), vr.GetPort())
	if len(target) == 0 {
		if len(s.Vectorizer) == 0 {
			err := errors.ErrFilterNotFound
			err = status.WrapWithInvalidArgument(vald.LinearSearchObjectRPCName+" API vectorizer configuration is invalid", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchObjectRPCName,
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
		target = s.Vectorizer
	}
	c, err := s.ingress.Target(ctx, target)
	if err != nil {
		err = status.WrapWithUnavailable(vald.LinearSearchObjectRPCName+" API target filter API unavailable", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec, err := c.GenVector(ctx, &payload.Object_Blob{
		Object: req.GetObject(),
	})
	if err != nil {
		err = status.WrapWithInternal(vald.LinearSearchObjectRPCName+" API failed to extract vector from filter", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return s.LinearSearch(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: req.GetConfig(),
	})
}

func (s *server) MultiLinearSearchObject(ctx context.Context, reqs *payload.Search_MultiObjectRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiLinearSearchObjectRPCName), apiName+"/"+vald.MultiLinearSearchObjectRPCName)
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiLinearSearchObjectRPCName+"/requestID-"+query.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()

			r, err := s.LinearSearchObject(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.NotFound,
					vald.MultiLinearSearchObjectRPCName+" API object "+string(query.GetObject())+"'s search request result not found",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get())
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}

				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf(vald.LinearSearchObjectRPCName+" API object %s's search request result not found",
							string(query.GetObject())), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf(vald.LinearSearchObjectRPCName+" API object %s's search request result not found",
								string(query.GetObject())), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		}))
	}
	wg.Wait()
	return res, errs
}

func (s *server) StreamLinearSearchObject(stream vald.Filter_StreamSearchObjectServer) error {
	ctx, span := trace.StartSpan(
		grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamLinearSearchObjectRPCName),
		apiName+"/"+vald.StreamLinearSearchObjectRPCName,
	)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err := grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_ObjectRequest) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamLinearSearchObjectRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()

			res, err := s.LinearSearchObject(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchObjectRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
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
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
		return err
	}
	return err
}

func (s *server) InsertObject(ctx context.Context, req *payload.Insert_ObjectRequest) (*payload.Object_Location, error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.InsertObjectRPCName), apiName+"/"+vald.InsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vr := req.GetVectorizer()
	if vr == nil || vr.GetPort() == 0 {
		err := errors.ErrFilterNotFound
		err = status.WrapWithInvalidArgument(vald.InsertObjectRPCName+" API vectorizer configuration is invalid", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetObject().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer port",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertObjectRPCName,
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
	if vr.GetHost() == "" {
		vr.Host = "localhost"
	}
	target := fmt.Sprintf("%s:%d", vr.GetHost(), vr.GetPort())
	if len(target) == 0 {
		if len(s.Vectorizer) == 0 {
			err := errors.ErrFilterNotFound
			err = status.WrapWithInvalidArgument(vald.InsertObjectRPCName+" API vectorizer configuration is invalid", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetObject().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertObjectRPCName,
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
		target = s.Vectorizer
	}
	c, err := s.ingress.Target(ctx, target)
	if err != nil {
		err = status.WrapWithUnavailable(vald.InsertObjectRPCName+" API target filter API unavailable", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetObject().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec, err := c.GenVector(ctx, req.GetObject())
	if err != nil {
		err = status.WrapWithInternal(vald.InsertObjectRPCName+" API failed to extract vector from filter", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetObject().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return s.Insert(ctx, &payload.Insert_Request{
		Vector: &payload.Object_Vector{
			Vector: vec.GetVector(),
			Id:     req.GetObject().GetId(),
		},
		Config: req.GetConfig(),
	})
}

func (s *server) StreamInsertObject(stream vald.Filter_StreamInsertObjectServer) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamInsertObjectRPCName), apiName+"/"+vald.StreamInsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err := grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Insert_ObjectRequest) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamInsertObjectRPCName+"/requestID-"+req.GetObject().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()

			loc, err := s.InsertObject(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.InsertObjectRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetObject().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: loc,
				},
			}, nil
		})
	if err != nil {
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiInsertObject(ctx context.Context, reqs *payload.Insert_MultiObjectRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiInsertObjectRPCName), apiName+"/"+vald.MultiInsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiInsertObjectRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()

			loc, err := s.InsertObject(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.InsertObjectRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetObject().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get())
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}

				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf(vald.MultiInsertObjectRPCName+" API object id: %s's insert failed",
							query.GetObject().GetId()), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf(vald.MultiInsertObjectRPCName+" API object id: %s's insert failed",
								query.GetObject().GetId()), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = loc
			return nil
		}))
	}
	wg.Wait()
	return locs, errs
}

func (s *server) UpdateObject(ctx context.Context, req *payload.Update_ObjectRequest) (*payload.Object_Location, error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.UpdateObjectRPCName), apiName+"/"+vald.UpdateObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vr := req.GetVectorizer()
	if vr == nil || vr.GetPort() == 0 {
		err := errors.ErrFilterNotFound
		err = status.WrapWithInvalidArgument(vald.UpdateObjectRPCName+" API vectorizer configuration is invalid", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetObject().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer port",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateObjectRPCName,
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
	if vr.GetHost() == "" {
		vr.Host = "localhost"
	}

	target := fmt.Sprintf("%s:%d", vr.GetHost(), vr.GetPort())
	if len(target) == 0 {
		if len(s.Vectorizer) == 0 {
			err := errors.ErrFilterNotFound
			err = status.WrapWithInvalidArgument(vald.UpdateObjectRPCName+" API vectorizer configuration is invalid", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetObject().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateObjectRPCName,
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
		target = s.Vectorizer
	}
	c, err := s.ingress.Target(ctx, target)
	if err != nil {
		err = status.WrapWithUnavailable(vald.UpdateObjectRPCName+" API target filter API unavailable", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetObject().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec, err := c.GenVector(ctx, req.GetObject())
	if err != nil {
		err = status.WrapWithInternal(vald.UpdateObjectRPCName+" API failed to extract vector from filter", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetObject().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return s.Update(ctx, &payload.Update_Request{
		Vector: &payload.Object_Vector{
			Vector: vec.GetVector(),
			Id:     req.GetObject().GetId(),
		},
		Config: req.GetConfig(),
	})
}

func (s *server) StreamUpdateObject(stream vald.Filter_StreamUpdateObjectServer) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamUpdateObjectRPCName), apiName+"/"+vald.StreamUpdateObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err := grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Update_ObjectRequest) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamUpdateObjectRPCName+"/id-"+req.GetObject().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			loc, err := s.UpdateObject(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpdateObjectRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetObject().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: loc,
				},
			}, nil
		})
	if err != nil {
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiUpdateObject(ctx context.Context, reqs *payload.Update_MultiObjectRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiUpdateObjectRPCName), apiName+"/"+vald.MultiUpdateObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiUpdateObjectRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			loc, err := s.UpdateObject(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.NotFound, "failed to parse "+vald.UpdateObjectRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetObject().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				log.Warn(err)
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiUpdateObject API object id: %s's insert failed",
							query.GetObject().GetId()), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiUpdateObject API object id: %s's insert failed",
								query.GetObject().GetId()), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = loc
			return nil
		}))
	}
	wg.Wait()
	return locs, errs
}

func (s *server) UpsertObject(ctx context.Context, req *payload.Upsert_ObjectRequest) (*payload.Object_Location, error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.UpsertObjectRPCName), apiName+"/"+vald.UpsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vr := req.GetVectorizer()
	if vr == nil || vr.GetPort() == 0 {
		err := errors.ErrFilterNotFound
		err = status.WrapWithInvalidArgument(vald.UpsertObjectRPCName+" API vectorizer configuration is invalid", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetObject().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer port",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertObjectRPCName,
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
	if vr.GetHost() == "" {
		vr.Host = "localhost"
	}
	target := fmt.Sprintf("%s:%d", vr.GetHost(), vr.GetPort())
	if len(target) == 0 {
		if len(s.Vectorizer) == 0 {
			err := errors.ErrFilterNotFound
			err = status.WrapWithInvalidArgument(vald.UpsertObjectRPCName+" API vectorizer configuration is invalid", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetObject().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertObjectRPCName,
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
		target = s.Vectorizer
	}
	c, err := s.ingress.Target(ctx, target)
	if err != nil {
		err = status.WrapWithUnavailable(vald.UpsertObjectRPCName+" API target filter API unavailable", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetObject().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec, err := c.GenVector(ctx, req.GetObject())
	if err != nil {
		err = status.WrapWithInternal(vald.UpsertObjectRPCName+" API failed to extract vector from filter", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetObject().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vectorizer targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return s.Upsert(ctx, &payload.Upsert_Request{
		Vector: &payload.Object_Vector{
			Vector: vec.GetVector(),
			Id:     req.GetObject().GetId(),
		},
		Config: req.GetConfig(),
	})
}

func (s *server) StreamUpsertObject(stream vald.Filter_StreamUpsertObjectServer) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamUpsertObjectRPCName), apiName+"/"+vald.StreamUpsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err := grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Upsert_ObjectRequest) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamUpdateObjectRPCName+"/id-"+req.GetObject().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()

			loc, err := s.UpsertObject(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpsertObjectRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetObject().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: loc,
				},
			}, nil
		})
	if err != nil {
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
		return err
	}
	return err
}

func (s *server) MultiUpsertObject(ctx context.Context, reqs *payload.Upsert_MultiObjectRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiUpsertObjectRPCName), apiName+"/"+vald.MultiUpsertObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiUpsertObjectRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			loc, err := s.UpsertObject(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpsertObjectRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetObject().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiUpsertObject API object id: %s's insert failed",
							query.GetObject().GetId()), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiUpsertObject API object id: %s's insert failed",
								query.GetObject().GetId()), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = loc
			return nil
		}))
	}
	wg.Wait()
	return locs, errs
}

func (s *server) Exists(ctx context.Context, meta *payload.Object_ID) (*payload.Object_ID, error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.ExistsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return s.gateway.Exists(ctx, meta, s.copts...)
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.SearchRPCName), apiName+"/"+vald.SearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	targets := req.GetConfig().GetIngressFilters().GetTargets()
	if targets != nil || s.SearchFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.SearchFilters))
		addrs = append(addrs, s.SearchFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.ingress.Target(ctx, addrs...)
		if err != nil {
			err = status.WrapWithUnavailable(
				fmt.Sprintf(vald.SearchRPCName+" API ingress filter targets %v not found", addrs),
				err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		vec, err := c.FilterVector(ctx, &payload.Object_Vector{
			Vector: req.GetVector(),
		})
		if err != nil {
			err = status.WrapWithInternal(
				fmt.Sprintf(vald.SearchRPCName+" API ingress filter request to %v failure on vec %v", addrs, req.GetVector()),
				err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		req.Vector = vec.GetVector()
	}
	res, err = s.gateway.Search(ctx, req, s.copts...)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	targets = req.GetConfig().GetEgressFilters().GetTargets()
	if targets != nil || s.DistanceFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.DistanceFilters))
		addrs = append(addrs, s.DistanceFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.egress.Target(ctx, addrs...)
		if err != nil {
			err = status.WrapWithUnavailable(
				fmt.Sprintf(vald.SearchRPCName+" API egress filter targets %v not found", addrs),
				err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		for i, dist := range res.GetResults() {
			d, err := c.FilterDistance(ctx, dist)
			if err != nil {
				err = status.WrapWithInternal(
					fmt.Sprintf(vald.SearchRPCName+" API egress filter request to %v failure on id %s", addrs, dist.GetId()),
					err,
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get())
				log.Warn(err)
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil, err
			}
			res.Results[i] = d
		}
	}
	return res, nil
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.SearchByIDRPCName), apiName+"/"+vald.SearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.gateway.SearchByID(ctx, req, s.copts...)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchByIDRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	targets := req.GetConfig().GetEgressFilters().GetTargets()
	if targets != nil || s.DistanceFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.DistanceFilters))
		addrs = append(addrs, s.DistanceFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.egress.Target(ctx, addrs...)
		if err != nil {
			err = status.WrapWithUnavailable(
				fmt.Sprintf(vald.SearchByIDRPCName+" API egress filter targets %v not found", addrs),
				err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		for i, dist := range res.GetResults() {
			d, err := c.FilterDistance(ctx, dist)
			if err != nil {
				err = status.WrapWithInternal(
					fmt.Sprintf(vald.SearchByIDRPCName+" API egress filter request to %v failure on id %s", addrs, dist.GetId()),
					err,
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get())
				log.Warn(err)
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil, err
			}
			res.Results[i] = d
		}
	}
	return res, nil
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamSearchRPCName), apiName+"/"+vald.StreamSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_Request) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamSearchRPCName+"/id-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Search(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.StreamSearchRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
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
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamSearchByIDRPCName), apiName+"/"+vald.StreamSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamSearchByIDRPCName+"/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.SearchByID(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchByIDRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequestFieldViolation{
							{
								Field:       "vectorizer targets",
								Description: err.Error(),
							},
						},
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
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
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiSearchRPCName), apiName+"/"+vald.MultiSearchRPCName)
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiSearchRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.Search(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.NotFound, "failed to parse "+vald.SearchRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequestFieldViolation{
							{
								Field:       "vectorizer targets",
								Description: err.Error(),
							},
						},
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiSearch API vector %v's search request result not found",
							query.GetVector()), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiSearch API vector %v's search request result not found",
								query.GetVector()), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		}))
	}
	wg.Wait()
	return res, errs
}

func (s *server) MultiSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiSearchByIDRPCName), apiName+"/"+vald.MultiSearchByIDRPCName)
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiSearchByIDRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.SearchByID(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.NotFound, "failed to parse "+vald.SearchByIDRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequestFieldViolation{
							{
								Field:       "vectorizer targets",
								Description: err.Error(),
							},
						},
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiSearchByID API id %s's search request result not found",
							query.GetId()), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiSearchByID API id %s's search request result not found",
								query.GetId()), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		}))
	}
	wg.Wait()
	return res, errs
}

func (s *server) LinearSearch(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.LinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	targets := req.GetConfig().GetIngressFilters().GetTargets()
	if targets != nil || s.SearchFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.SearchFilters))
		addrs = append(addrs, s.SearchFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.ingress.Target(ctx, addrs...)
		if err != nil {
			err = status.WrapWithUnavailable(
				fmt.Sprintf(vald.LinearSearchRPCName+" API ingress filter targets %v not found", addrs),
				err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		vec, err := c.FilterVector(ctx, &payload.Object_Vector{
			Vector: req.GetVector(),
		})
		if err != nil {
			err = status.WrapWithInternal(
				fmt.Sprintf(vald.LinearSearchRPCName+" API ingress filter request to %v failure on vec %v", addrs, req.GetVector()),
				err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		req.Vector = vec.GetVector()
	}
	res, err = s.gateway.LinearSearch(ctx, req, s.copts...)
	if err != nil {
		return nil, err
	}
	targets = req.GetConfig().GetEgressFilters().GetTargets()
	if targets != nil || s.DistanceFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.DistanceFilters))
		addrs = append(addrs, s.DistanceFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.egress.Target(ctx, addrs...)
		if err != nil {
			err = status.WrapWithUnavailable(
				fmt.Sprintf(vald.LinearSearchRPCName+" API ingress filter targets %v not found", addrs),
				err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		for i, dist := range res.GetResults() {
			d, err := c.FilterDistance(ctx, dist)
			if err != nil {
				err = status.WrapWithInternal(
					fmt.Sprintf(vald.LinearSearchRPCName+" API egress filter request to %v failure on id %s", addrs, dist.GetId()),
					err,
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get())
				log.Warn(err)
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil, err
			}
			res.Results[i] = d
		}
	}
	return res, nil
}

func (s *server) LinearSearchByID(ctx context.Context, req *payload.Search_IDRequest) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.LinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.gateway.LinearSearchByID(ctx, req, s.copts...)
	if err != nil {
		return nil, err
	}
	targets := req.GetConfig().GetEgressFilters().GetTargets()
	if targets != nil || s.DistanceFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.DistanceFilters))
		addrs = append(addrs, s.DistanceFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.egress.Target(ctx, addrs...)
		if err != nil {
			err = status.WrapWithUnavailable(
				fmt.Sprintf(vald.LinearSearchByIDRPCName+" API egress filter targets %v not found", addrs),
				err,
				&errdetails.RequestInfo{
					RequestId:   req.GetConfig().GetRequestId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		for i, dist := range res.GetResults() {
			d, err := c.FilterDistance(ctx, dist)
			if err != nil {
				err = status.WrapWithInternal(
					fmt.Sprintf(vald.LinearSearchByIDRPCName+" API egress filter request to %v failure on id %s", addrs, dist.GetId()),
					err,
					&errdetails.RequestInfo{
						RequestId:   dist.GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get())
				log.Warn(err)
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return nil, err
			}
			res.Results[i] = d
		}
	}
	return res, nil
}

func (s *server) StreamLinearSearch(stream vald.Search_StreamLinearSearchServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamLinearSearchRPCName), apiName+"/"+vald.StreamLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_Request) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamLinearSearchRPCName+"/id-"+req.Config.RequestId)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearch(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequestFieldViolation{
							{
								Field:       "vectorizer targets",
								Description: err.Error(),
							},
						},
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
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
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) StreamLinearSearchByID(stream vald.Search_StreamLinearSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(
		grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamLinearSearchByIDRPCName),
		apiName+"/"+vald.StreamLinearSearchByIDRPCName,
	)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_StreamResponse, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamLinearSearchByIDRPCName+"/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearchByID(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchByIDRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
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
			"failed to parse "+vald.StreamLinearSearchRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiLinearSearchRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.LinearSearch(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchByIDRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiLinearSearch API vector %v's search request result not found",
							query.GetVector()), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiLinearSearch API vector %v's search request result not found",
								query.GetVector()), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		}))
	}
	wg.Wait()
	return res, errs
}

func (s *server) MultiLinearSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiLinearSearchByIDRPCName), apiName+"/"+vald.MultiLinearSearchByIDRPCName)
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
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiLinearSearchByIDRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.LinearSearchByID(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchByIDRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiLinearSearchByID API id %s's search request result not found",
							query.GetId()), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiLinearSearchByID API id %s's search request result not found",
								query.GetId()), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			res.Responses[idx] = r
			return nil
		}))
	}
	wg.Wait()
	return res, errs
}

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.InsertRPCName), apiName+"/"+vald.InsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	uuid := vec.GetId()
	if len(vec.GetVector()) < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(len(vec.GetVector()), 0)
		err = status.WrapWithInvalidArgument(vald.InsertRPCName+" API invalid vector argument", err,
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
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, _ := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if id != nil || len(id.GetId()) > 0 {
			err = errors.ErrMetaDataAlreadyExists(uuid)
			err = status.WrapWithAlreadyExists(vald.InsertRPCName+" API ID = "+uuid+" already exists", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + "." + vald.ExistsRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Insert_Config{SkipStrictExistCheck: true}
		}
	}
	targets := req.GetConfig().GetFilters().GetTargets()
	if len(targets) == 0 && len(s.InsertFilters) == 0 {
		return s.gateway.Insert(ctx, req)
	}
	addrs := make([]string, 0, len(targets)+len(s.InsertFilters))
	addrs = append(addrs, s.InsertFilters...)
	for _, target := range targets {
		addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
	}
	c, err := s.ingress.Target(ctx, addrs...)
	if err != nil {
		err = status.WrapWithUnavailable(
			fmt.Sprintf(vald.InsertRPCName+" API ingress filter filter targets %v not found", addrs), err,
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
	}
	vec, err = c.FilterVector(ctx, req.GetVector())
	if err != nil {
		err = status.WrapWithInternal(
			fmt.Sprintf(vald.InsertRPCName+" API ingress filter request to %v failure on id: %s\tvec: %v", addrs, req.GetVector().GetId(), req.GetVector().GetVector()), err,
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if vec.GetId() == "" {
		vec.Id = req.GetVector().GetId()
	}
	req.Vector = vec
	loc, err = s.gateway.Insert(ctx, req, s.copts...)
	if err != nil {
		err = status.WrapWithInternal(
			vald.InsertRPCName+" API failed to Execute DoMulti ID = "+uuid, err,
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get(),
		)
		log.Error(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamInsertRPCName), apiName+"/"+vald.StreamInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Insert_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamInsertRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Insert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.InsertRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequestFieldViolation{
							{
								Field:       "vectorizer targets",
								Description: err.Error(),
							},
						},
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)

				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiInsertRPCName), apiName+"/"+vald.MultiInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiInsertRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.Insert(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.InsertRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequestFieldViolation{
							{
								Field:       "vectorizer targets",
								Description: err.Error(),
							},
						},
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiInsert API request %#v's Insert request result not found",
							query), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiInsert API request %#v's Insert request result not found",
								query), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = r
			return nil
		}))
	}
	wg.Wait()
	return locs, errs
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.UpdateRPCName), apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	uuid := vec.GetId()
	if len(vec.GetVector()) < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(len(vec.GetVector()), 0)
		err = status.WrapWithInvalidArgument(vald.UpdateRPCName+" API invalid vector argument", err,
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
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, _ := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if id != nil || len(id.GetId()) > 0 {
			err = status.WrapWithAlreadyExists(vald.UpdateRPCName+" API ID = "+uuid+"'s same vector data already exists", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.ExistsRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Update_Config{SkipStrictExistCheck: true}
		}
	}
	targets := req.GetConfig().GetFilters().GetTargets()
	if len(targets) == 0 && len(s.UpdateFilters) == 0 {
		return s.gateway.Update(ctx, req)
	}
	addrs := make([]string, 0, len(targets)+len(s.UpdateFilters))
	addrs = append(addrs, s.UpdateFilters...)
	for _, target := range targets {
		addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
	}
	c, err := s.ingress.Target(ctx, addrs...)
	if err != nil {
		err = status.WrapWithUnavailable(
			fmt.Sprintf(vald.UpdateRPCName+" API ingress filter filter targets %v not found", addrs), err,
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec, err = c.FilterVector(ctx, req.GetVector())
	if err != nil {
		err = status.WrapWithInternal(
			fmt.Sprintf(vald.UpdateRPCName+" API ingress filter request to %v failure on id: %s\tvec: %v", addrs, req.GetVector().GetId(), req.GetVector().GetVector()), err,
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if vec.GetId() == "" {
		vec.Id = req.GetVector().GetId()
	}
	req.Vector = vec
	loc, err = s.gateway.Update(ctx, req, s.copts...)
	if err != nil {
		err = status.WrapWithInternal(
			vald.UpdateRPCName+" API failed to Execute DoMulti ID = "+uuid, err,
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamUpdateRPCName), apiName+"/"+vald.StreamUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Update_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamUpdateRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Update(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpdateRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
		return err
	}
	return nil
}

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiUpdateRPCName), apiName+"/"+vald.MultiUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiUpdateRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.Update(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpdateRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.BadRequest{
						FieldViolations: []*errdetails.BadRequestFieldViolation{
							{
								Field:       "vectorizer targets",
								Description: err.Error(),
							},
						},
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiUpdate API request %#v's Update request result not found",
							query), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiUpdate API request %#v's Update request result not found",
								query), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = r
			return nil
		}))
	}
	wg.Wait()
	return locs, errs
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.UpsertRPCName), apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	uuid := vec.GetId()
	if len(vec.GetVector()) < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(len(vec.GetVector()), 0)
		err = status.WrapWithInvalidArgument(vald.UpsertRPCName+" API invalid vector argument", err,
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
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, _ := s.Exists(ctx, &payload.Object_ID{
			Id: uuid,
		})
		if id != nil || len(id.GetId()) > 0 {
			err = status.WrapWithAlreadyExists(vald.UpsertRPCName+" API ID = "+uuid+"'s same vector data already exists", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + vald.ExistsRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeAlreadyExists(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if req.GetConfig() != nil {
			req.GetConfig().SkipStrictExistCheck = true
		} else {
			req.Config = &payload.Upsert_Config{SkipStrictExistCheck: true}
		}
	}
	targets := req.GetConfig().GetFilters().GetTargets()
	if len(targets) == 0 && len(s.UpsertFilters) == 0 {
		return s.gateway.Upsert(ctx, req)
	}
	addrs := make([]string, 0, len(targets))
	addrs = append(addrs, s.UpsertFilters...)
	for _, target := range targets {
		addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
	}
	c, err := s.ingress.Target(ctx, addrs...)
	if err != nil {
		err = status.WrapWithUnavailable(
			fmt.Sprintf(vald.UpsertRPCName+" API ingress filter filter targets %v not found", addrs), err,
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	vec, err = c.FilterVector(ctx, req.GetVector())
	if err != nil {
		err = status.WrapWithInternal(
			fmt.Sprintf(vald.UpsertRPCName+" API ingress filter request to %v failure on id: %s\tvec: %v", addrs, req.GetVector().GetId(), req.GetVector().GetVector()), err,
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	if vec.GetId() == "" {
		vec.Id = req.GetVector().GetId()
	}
	req.Vector = vec
	loc, err = s.gateway.Upsert(ctx, req, s.copts...)
	if err != nil {
		err = status.WrapWithInternal(vald.UpsertRPCName+" API failed to Execute DoMulti ID = "+uuid, err,
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get(),
		)

		log.Error(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamUpsertRPCName), apiName+"/"+vald.StreamUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Upsert_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamUpsertRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Upsert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpsertRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.StreamUpsertRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiUpsertRPCName), apiName+"/"+vald.MultiUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiUpsertRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()

			r, err := s.Upsert(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpsertRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiUpsert API request %#v's Upsert request result not found",
							query), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiUpsert API request %#v's Upsert request result not found",
								query), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = r
			return nil
		}))
	}
	wg.Wait()
	return locs, errs
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return s.gateway.Remove(ctx, req, s.copts...)
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamRemoveRPCName), apiName+"/"+vald.StreamRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Remove_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.StreamRemoveRPCName+"/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Remove(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.RemoveRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetId().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.StreamRemoveRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
		return err
	}
	return nil
}

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, req := range reqs.Requests {
		idx, query := i, req
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ctx, sspan := trace.StartSpan(ctx, apiName+"."+vald.MultiRemoveRPCName+"/errgroup.Go/id-"+strconv.Itoa(idx))
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			r, err := s.Remove(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.NotFound,
					fmt.Sprintf(vald.MultiRemoveRPCName+" API ID = %v not found", query.GetId().GetId()),
					&errdetails.RequestInfo{
						RequestId:   query.GetId().GetId(),
						ServingData: errdetails.Serialize(reqs),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get())
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}

				mu.Lock()
				if errs == nil {
					errs = status.WrapWithNotFound(
						fmt.Sprintf("MultiRemove API id %s's Remove request result not found",
							query.GetId()), err, info.Get())
				} else {
					errs = errors.Join(errs,
						status.WrapWithNotFound(
							fmt.Sprintf("MultiRemove API id %s's Remove request result not found",
								query.GetId()), err, info.Get()))
				}
				mu.Unlock()
				return nil
			}
			locs.Locations[idx] = r
			return nil
		}))
	}
	wg.Wait()
	return locs, errs
}

func (s *server) GetObject(ctx context.Context, req *payload.Object_VectorRequest) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec, err = s.gateway.GetObject(ctx, req)
	if err != nil {
		err = status.WrapWithNotFound(vald.GetObjectRPCName+" API failed to extract vector from filter", err,
			&errdetails.RequestInfo{
				RequestId:   req.GetId().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	targets := req.GetFilters().GetTargets()
	if targets != nil || s.ObjectFilters != nil {
		addrs := make([]string, 0, len(targets)+len(s.ObjectFilters))
		addrs = append(addrs, s.ObjectFilters...)
		for _, target := range targets {
			addrs = append(addrs, fmt.Sprintf("%s:%d", target.GetHost(), target.GetPort()))
		}
		c, err := s.egress.Target(ctx, addrs...)
		if err != nil {
			err = status.WrapWithUnavailable(vald.SearchObjectRPCName+" API target filter API unavailable", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetId().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "vectorizer targets",
							Description: err.Error(),
						},
					},
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeUnavailable(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		vec, err = c.FilterVector(ctx, vec)
		if err != nil {
			err = status.WrapWithInternal(vald.GetObjectRPCName+" API egress filter API failed", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetId().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
	}
	return vec, nil
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.FilterRPCServiceName+"/"+vald.StreamSearchObjectRPCName), apiName+"/"+vald.StreamGetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Object_VectorRequest) (*payload.Object_StreamVector, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+".StreamGetObject/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.GetObject(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.GetObjectRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetId().GetId(),
						ServingData: errdetails.Serialize(req),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					}, info.Get(),
				)
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, err.Error())
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
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Error(err)
		return err
	}
	return nil
}
