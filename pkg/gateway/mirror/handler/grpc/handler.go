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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	vclient "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

type server struct {
	eg                errgroup.Group
	gateway           service.Gateway // Mirror Gateway client service.
	discoverer        service.Discoverer
	vc                vclient.Client // Vald gateway client (LB gateway) for the same cluster.
	timeout           time.Duration
	replica           int
	streamConcurrency int
	name              string
	ip                string
	mirror.UnimplementedValdServerWithMirror
}

const apiName = "vald/gateway/mirror"

func New(opts ...Option) mirror.Server {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Register(ctx context.Context, req *payload.Mirror_Targets) (*payload.Mirror_Targets, error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, mirror.PackageName+"."+mirror.MirrorRPCServiceName+"/"+mirror.RegisterRPCName), apiName+"/"+mirror.RegisterRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err := s.discoverer.Connect(ctx, req.GetTargets()...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + mirror.RegisterRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				mirror.RegisterRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + mirror.RegisterRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithCanceled(
				mirror.RegisterRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + mirror.RegisterRPCName,
			)
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				mirror.RegisterRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, errors.ErrTargetNotFound):
			err = status.WrapWithInvalidArgument(
				mirror.RegisterRPCName+" API target not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+mirror.RegisterRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn("failed to process register request\terror: %s", err.Error())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return req, nil
}

func (s *server) Advertise(ctx context.Context, req *payload.Mirror_Targets) (res *payload.Mirror_Targets, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, mirror.PackageName+"."+mirror.MirrorRPCServiceName+"/"+mirror.AdvertiseRPCName), apiName+"/"+mirror.AdvertiseRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = s.Register(ctx, req)
	if err != nil {
		return nil, err
	}
	tgts, err := s.discoverer.MirrorTargets()
	if err != nil {
		err = status.WrapWithInternal(mirror.AdvertiseRPCName+" API failed to get connected mirror gateway targets", err,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "mirror gateway targets",
						Description: err.Error(),
					},
				},
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + mirror.AdvertiseRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return &payload.Mirror_Targets{
		Targets: tgts,
	}, nil
}

func (s *server) Exists(ctx context.Context, meta *payload.Object_ID) (id *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.ExistsRPCName), apiName+"/"+vald.ExistsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	id, err = s.vc.Exists(ctx, meta, s.vc.GRPCClient().GetCallOption()...)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.ExistsRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   meta.GetId(),
				ServingData: errdetails.Serialize(meta),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.ExistsRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return id, nil
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.SearchRPCName), apiName+"/"+vald.SearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.vc.Search(ctx, req)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.SearchRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error,
) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.SearchByIDRPCName), apiName+"/"+vald.SearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.vc.SearchByID(ctx, req)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.SearchByIDRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.SearchByIDRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamSearchRPCName), apiName+"/"+vald.StreamSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_Request)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Search(ctx, data.(*payload.Search_Request))
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchRPCName+" gRPC error response")
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
			"failed to parse "+vald.StreamSearchRPCName+" gRPC error response")
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
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamSearchByIDRPCName), apiName+"/"+vald.StreamSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_IDRequest)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"."+vald.StreamSearchByIDRPCName+"/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.SearchByID(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchByIDRPCName+" gRPC error response")
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
			"failed to parse "+vald.StreamSearchByIDRPCName+" gRPC error response")
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
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiSearchRPCName), apiName+"/"+vald.MultiSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup
	rids := make([]string, 0, len(reqs.GetRequests()))

	for i, req := range reqs.GetRequests() {
		idx, query := i, req
		rids = append(rids, query.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + strconv.Itoa(idx)
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiSearchRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			r, err := s.Search(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Responses[idx] = r
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiSearchRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiSearchRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) MultiSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiSearchByIDRPCName), apiName+"/"+vald.MultiSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup
	rids := make([]string, 0, len(reqs.GetRequests()))

	for i, req := range reqs.GetRequests() {
		idx, query := i, req
		rids = append(rids, query.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + strconv.Itoa(idx)
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiSearchByIDRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			r, err := s.SearchByID(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.SearchByIDRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Responses[idx] = r
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiSearchByIDRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiSearchByIDRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) LinearSearch(ctx context.Context, req *payload.Search_Request) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.LinearSearchRPCName), apiName+"/"+vald.LinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.vc.LinearSearch(ctx, req)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.LinearSearchRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) LinearSearchByID(ctx context.Context, req *payload.Search_IDRequest) (
	res *payload.Search_Response, err error,
) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.LinearSearchByIDRPCName), apiName+"/"+vald.LinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = s.vc.LinearSearchByID(ctx, req)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.LinearSearchByIDRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetConfig().GetRequestId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.LinearSearchByIDRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *server) StreamLinearSearch(stream vald.Search_StreamLinearSearchServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamLinearSearchRPCName), apiName+"/"+vald.StreamLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_Request)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamLinearSearchRPCName+"/requestID-"+req.GetConfig().GetRequestId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearch(ctx, data.(*payload.Search_Request))
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchRPCName+" gRPC error response")
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
		return err
	}
	return nil
}

func (s *server) StreamLinearSearchByID(stream vald.Search_StreamLinearSearchByIDServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.StreamLinearSearchByIDRPCName), apiName+"/"+vald.StreamLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Search_IDRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Search_IDRequest)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"."+vald.StreamLinearSearchByIDRPCName+"/id-"+req.GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.LinearSearchByID(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchByIDRPCName+" gRPC error response")
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
			"failed to parse "+vald.StreamLinearSearchByIDRPCName+" gRPC error response")
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
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiLinearSearchRPCName), apiName+"/"+vald.MultiLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup
	rids := make([]string, 0, len(reqs.GetRequests()))

	for i, req := range reqs.GetRequests() {
		idx, query := i, req
		rids = append(rids, query.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + strconv.Itoa(idx)
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiLinearSearchRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			r, err := s.LinearSearch(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Responses[idx] = r
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiLinearSearchRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiLinearSearchRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) MultiLinearSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.SearchRPCServiceName+"/"+vald.MultiLinearSearchByIDRPCName), apiName+"/"+vald.MultiLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Search_Responses{
		Responses: make([]*payload.Search_Response, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup
	rids := make([]string, 0, len(reqs.GetRequests()))

	for i, req := range reqs.GetRequests() {
		idx, query := i, req
		rids = append(rids, query.GetConfig().GetRequestId())
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + strconv.Itoa(idx)
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiLinearSearchByIDRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			r, err := s.LinearSearchByID(ctx, query)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.LinearSearchByIDRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   query.GetConfig().GetRequestId(),
						ServingData: errdetails.Serialize(query),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Responses[idx] = r
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiLinearSearchByIDRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(rids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiLinearSearchByIDRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (ce *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.InsertRPCName), apiName+"/"+vald.InsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	reqSrcPodName := s.gateway.FromForwardedContext(ctx)

	// When this condition is matched, the request is proxied to another Mirror gateway.
	// So this component sends the request only to its cluster's Vald gateway (LB gateway).
	if len(reqSrcPodName) != 0 {
		ce, err = s.insert(ctx, s.vc, req, s.vc.GRPCClient().GetCallOption()...)
		if err != nil {
			return nil, err
		}
		log.Debugf("Insert API succeeded to %#v", ce)
		return ce, nil
	}

	// When this condition is matched, this Mirror gateway is the starting point of the mirror process.
	// This component sends requests to the Mirror Gateways of other clusters and its cluster's Vald gateway (LB gateway).
	var insertErrs error
	var errMutex, mutex sync.Mutex
	var successTargets sync.Map
	ce = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}

	// This process sends the request to the Mirror Gateways of other clusters and its cluster's Vald gateway (LB gateway).
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		loc, err := s.insert(ctx, vald.NewValdClient(conn), req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.InsertRPCName+" API gRPC error response for "+target,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if st.Code() == codes.AlreadyExists {
				// NOTE: If it is strictly necessary to check, fix this logic.
				return nil
			}
			log.Warn(err)
			errMutex.Lock()
			if insertErrs != nil {
				insertErrs = errors.Wrap(insertErrs, err.Error())
			} else {
				insertErrs = err
			}
			errMutex.Unlock()
			return err
		}
		mutex.Lock()
		ce.Name = loc.GetName()
		ce.Ips = append(ce.Ips, loc.GetIps()...)
		mutex.Unlock()
		successTargets.Store(target, struct{}{})
		return nil
	})
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.InsertRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + ".BroadCast",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if insertErrs == nil {
			insertErrs = err
		} else {
			insertErrs = errors.Wrap(insertErrs, err.Error())
		}
	}
	if insertErrs == nil {
		log.Debugf("Insert API succeeded to %#v", ce)
		return ce, nil
	}
	log.Warnf("failed to Insert API: %s", insertErrs.Error())

	var rollbackErrs error
	rmReq := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
	}

	// This process sends the rollback request to the Mirror Gateways of other clusters and its cluster's Vald gateway (LB gateway).
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "rollback/BroadCast/"+target), apiName+"/"+vald.InsertRPCName+"/rollback/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		if _, ok := successTargets.Load(target); !ok {
			return nil
		}

		_, err := s.remove(ctx, vald.NewValdClient(conn), rmReq, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveRPCName+" for "+vald.InsertRPCName+" error response for "+target,
				&errdetails.RequestInfo{
					RequestId:   rmReq.GetId().GetId(),
					ServingData: errdetails.Serialize(rmReq),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + "." + vald.RemoveRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if st.Code() == codes.NotFound {
				return nil
			}
			log.Warn(err)
			errMutex.Lock()
			if rollbackErrs == nil {
				rollbackErrs = err
			} else {
				rollbackErrs = errors.Wrap(rollbackErrs, err.Error())
			}
			errMutex.Unlock()
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.RemoveRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + ".BroadCast",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if rollbackErrs == nil {
			rollbackErrs = err
		} else {
			rollbackErrs = errors.Wrap(rollbackErrs, err.Error())
		}
	}
	if rollbackErrs != nil {
		log.Warn("failed to rollback for Insert API: %s", rollbackErrs.Error())
		st, msg, err := status.ParseError(rollbackErrs, codes.Internal,
			"failed to parse "+vald.RemoveRPCName+" for "+vald.InsertRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   rmReq.GetId().GetId(),
				ServingData: errdetails.Serialize(rmReq),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + "." + vald.RemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	st, msg, err := status.ParseError(insertErrs, codes.Internal,
		"failed to parse "+vald.InsertRPCName+" gRPC error response",
		&errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		},
		&errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		},
	)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
}

func (s *server) insert(ctx context.Context, client vald.InsertClient, req *payload.Insert_Request, opts ...grpc.CallOption) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "insert"), apiName+"/insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err = client.Insert(ctx, req, opts...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.InsertRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.InsertRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
			)
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.InsertRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.InsertRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) remove(ctx context.Context, client vald.RemoveClient, req *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "remove"), apiName+"/remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err := client.Remove(ctx, req, opts...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetId().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.RemoveRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.RemoveRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
			)
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.RemoveRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.StreamInsertRPCName), apiName+"/"+vald.StreamInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Insert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Insert_Request)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamInsertRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Insert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.InsertRPCName+" gRPC error response")
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
		},
	)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.StreamInsertRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.MultiInsertRPCName), apiName+"/"+vald.MultiInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup

	for i, r := range reqs.GetRequests() {
		idx, req := i, r
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + req.GetVector().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiInsertRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			loc, err := s.Insert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.InsertRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Locations[idx] = loc
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiInsertRPCName+" gRPC error response",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiInsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateRPCName), apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	reqSrcPodName := s.gateway.FromForwardedContext(ctx)

	if len(reqSrcPodName) != 0 {
		loc, err = s.update(ctx, s.vc, req, s.vc.GRPCClient().GetCallOption()...)
		if err != nil {
			return nil, err
		}
		log.Debugf("Update API succeeded to %#v", loc)
		return loc, nil
	}

	objReq := &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
	}
	oldVecs, err := s.getObjects(ctx, objReq)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.GetObjectRPCName+" gRPC error response for "+vald.UpdateRPCName+" API",
			&errdetails.RequestInfo{
				RequestId:   objReq.GetId().GetId(),
				ServingData: errdetails.Serialize(objReq),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.GetObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	var updateErrs error
	var errMutex, mutex sync.Mutex
	var successTargets sync.Map
	ce := &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpdateRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		loc, err := s.update(ctx, vald.NewUpdateClient(conn), req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpdateRPCName+" API error response for "+target,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if st.Code() == codes.NotFound || st.Code() == codes.AlreadyExists {
				return nil
			}
			log.Warn(err)
			errMutex.Lock()
			if updateErrs == nil {
				updateErrs = err
			} else {
				updateErrs = errors.Wrap(updateErrs, err.Error())
			}
			errMutex.Unlock()
			return err
		}
		mutex.Lock()
		ce.Name = loc.GetName()
		ce.Ips = append(loc.Ips, loc.GetIps()...)
		mutex.Unlock()
		successTargets.Store(target, struct{}{})
		return nil
	})
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpdateRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + ".BroadCast",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if updateErrs == nil {
			updateErrs = err
		} else {
			updateErrs = errors.Wrap(updateErrs, err.Error())
		}
	}
	if updateErrs == nil {
		log.Debugf("Update API succeeded to %#v", ce)
		return ce, nil
	}
	log.Warnf("failed to Update API: %s", updateErrs.Error())

	// Rollback for Update API.
	var rollbackErrs error
	rmReq := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "rollback/BroadCast/"+target), apiName+"."+vald.UpdateRPCName+"/rollback/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		if _, ok := successTargets.Load(target); !ok {
			return nil
		}
		oldVec, ok := oldVecs.Load(target)
		if !ok || oldVec == nil {
			_, err := s.remove(ctx, vald.NewRemoveClient(conn), rmReq, copts...)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal,
					"failed to parse "+vald.RemoveRPCName+" for "+vald.UpdateRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   rmReq.GetId().GetId(),
						ServingData: errdetails.Serialize(rmReq),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.RemoveRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					},
				)
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				if st.Code() == codes.NotFound {
					return nil
				}
				errMutex.Lock()
				if rollbackErrs == nil {
					rollbackErrs = err
				} else {
					rollbackErrs = errors.Wrap(rollbackErrs, err.Error())
				}
				errMutex.Unlock()
				return err
			}
			return nil
		}

		req := &payload.Update_Request{
			Vector: oldVec.(*payload.Object_Vector),
		}
		_, err := s.update(ctx, vald.NewUpdateClient(conn), req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpdateRPCName+" for "+vald.UpdateRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.UpdateRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if st.Code() == codes.NotFound || st.Code() == codes.AlreadyExists {
				return nil
			}
			errMutex.Lock()
			if rollbackErrs == nil {
				rollbackErrs = err
			} else {
				rollbackErrs = errors.Wrap(rollbackErrs, err.Error())
			}
			errMutex.Unlock()
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpdateRPCName+" for "+vald.UpdateRPCName+" API connection not found", err,
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + ".Rollback.BroadCast",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if rollbackErrs == nil {
			rollbackErrs = err
		} else {
			rollbackErrs = errors.Wrap(rollbackErrs, err.Error())
		}
	}
	if rollbackErrs != nil {
		st, msg, err := status.ParseError(rollbackErrs, codes.Internal,
			"failed to parse "+vald.UpdateRPCName+" for "+vald.UpdateRPCName+" gRPC error response",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.UpdateRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	st, msg, err := status.ParseError(updateErrs, codes.Internal,
		"failed to parse "+vald.UpdateRPCName+" gRPC error response",
		&errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		},
	)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
}

func (s *server) update(ctx context.Context, client vald.UpdateClient, req *payload.Update_Request, opts ...grpc.CallOption) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "update"), apiName+"/update")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err = client.Update(ctx, req, opts...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.UpdateRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.UpdateRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName,
			)
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.UpdateRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpdateRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.StreamUpdateRPCName), apiName+"/"+vald.StreamUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Update_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Update_Request)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamUpdateRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Update(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpdateRPCName+" gRPC error response")
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.StreamUpdateRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.MultiUpdateRPCName), apiName+"/"+vald.MultiUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup

	for i, r := range reqs.GetRequests() {
		idx, req := i, r
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + req.GetVector().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiUpdateRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			loc, err := s.Update(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpdateRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Locations[idx] = loc
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiUpdateRPCName+" gRPC error response",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpdateRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.UpsertRPCName), apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	objReq := &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
	}
	oldVecs, err := s.getObjects(ctx, objReq)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.GetObjectRPCName+" gRPC error response for "+vald.UpsertRPCName+" API",
			&errdetails.RequestInfo{
				RequestId:   objReq.GetId().GetId(),
				ServingData: errdetails.Serialize(objReq),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + vald.GetObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	var upsertErrs error
	var errMutex, mutex sync.Mutex
	var successTargets sync.Map
	ce := &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0),
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.UpdateRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		loc, err := s.upsert(ctx, vald.NewUpsertClient(conn), req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpsertRPCName+" API error response for "+target,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if st.Code() == codes.NotFound || st.Code() == codes.AlreadyExists {
				return nil
			}
			log.Warn(err)
			errMutex.Lock()
			if upsertErrs == nil {
				upsertErrs = err
			} else {
				upsertErrs = errors.Wrap(upsertErrs, err.Error())
			}
			errMutex.Unlock()
			return err
		}
		mutex.Lock()
		ce.Name = loc.GetName()
		ce.Ips = append(loc.Ips, loc.GetIps()...)
		mutex.Unlock()
		successTargets.Store(target, struct{}{})
		return nil
	})
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpsertRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + ".BroadCast",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if upsertErrs == nil {
			upsertErrs = err
		} else {
			upsertErrs = errors.Wrap(upsertErrs, err.Error())
		}
	}
	if upsertErrs == nil {
		log.Debugf("Upsert API succeeded to %#v", ce)
		return ce, nil
	}
	log.Warnf("failed to Upsert API: %s", upsertErrs.Error())

	// Rollback for Update API.
	var rollbackErrs error
	rmReq := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "rollback/BroadCast/"+target), apiName+"."+vald.UpsertRPCName+"/rollback/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		if _, ok := successTargets.Load(target); !ok {
			return nil
		}
		oldVec, ok := oldVecs.Load(target)
		if !ok || oldVec == nil {
			_, err := s.remove(ctx, vald.NewRemoveClient(conn), rmReq, copts...)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal,
					"failed to parse "+vald.RemoveRPCName+" for "+vald.UpsertRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   rmReq.GetId().GetId(),
						ServingData: errdetails.Serialize(rmReq),
					},
					&errdetails.ResourceInfo{
						ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + vald.RemoveRPCName,
						ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
					},
				)
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				if st.Code() == codes.NotFound {
					return nil
				}
				errMutex.Lock()
				if rollbackErrs == nil {
					rollbackErrs = err
				} else {
					rollbackErrs = errors.Wrap(rollbackErrs, err.Error())
				}
				errMutex.Unlock()
				return err
			}
			return nil
		}

		req := &payload.Update_Request{
			Vector: oldVec.(*payload.Object_Vector),
		}
		_, err := s.update(ctx, vald.NewUpdateClient(conn), req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpdateRPCName+" for "+vald.UpdateRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.UpdateRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if st.Code() == codes.NotFound || st.Code() == codes.AlreadyExists {
				return nil
			}
			errMutex.Lock()
			if rollbackErrs == nil {
				rollbackErrs = err
			} else {
				rollbackErrs = errors.Wrap(rollbackErrs, err.Error())
			}
			errMutex.Unlock()
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpdateRPCName+" for "+vald.UpdateRPCName+" API connection not found", err,
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + ".Rollback.BroadCast",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if rollbackErrs == nil {
			rollbackErrs = err
		} else {
			rollbackErrs = errors.Wrap(rollbackErrs, err.Error())
		}
	}
	if rollbackErrs != nil {
		st, msg, err := status.ParseError(rollbackErrs, codes.Internal,
			"failed to parse "+vald.UpdateRPCName+" for "+vald.UpsertRPCName+" gRPC error response",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + vald.UpdateRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	st, msg, err := status.ParseError(upsertErrs, codes.Internal,
		"failed to parse "+vald.UpsertRPCName+" gRPC error response",
		&errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		},
	)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
}

func (s *server) upsert(ctx context.Context, client vald.UpsertClient, req *payload.Upsert_Request, opts ...grpc.CallOption) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "upsert"), apiName+"/upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	loc, err = client.Upsert(ctx, req, opts...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetVector().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.UpsertRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.UpsertRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName,
			)
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.UpsertRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpsertRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn("failed to process Upsert request\terror: %s", err.Error())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.StreamUpsertRPCName), apiName+"/"+vald.StreamUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Upsert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Upsert_Request)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamUpsertRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Upsert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpsertRPCName+" gRPC error response")
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
		return err
	}
	return nil
}

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.MultiUpsertRPCName), apiName+"/"+vald.MultiUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup

	for i, r := range reqs.GetRequests() {
		idx, req := i, r
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + req.GetVector().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiUpsertRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			loc, err := s.Upsert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.UpsertRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetVector().GetId(),
						ServingData: errdetails.Serialize(req),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Locations[idx] = loc
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiUpsertRPCName+" gRPC error response",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiUpsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	reqSrcPodName := s.gateway.FromForwardedContext(ctx)

	// When this condition is matched, the request is proxied to another Mirror gateway.
	// So this component sends the request only to the Vald gateway (LB gateway) of own cluster.
	if len(reqSrcPodName) != 0 {
		loc, err = s.remove(ctx, s.vc, req, s.vc.GRPCClient().GetCallOption()...)
		if err != nil {
			return nil, err
		}
		log.Debugf("Remove API remove succeeded to %#v", loc)
		return loc, nil
	}

	objReq := &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetId().GetId(),
		},
	}
	oldVecs, err := s.getObjects(ctx, objReq)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.GetObjectRPCName+" for "+vald.UpsertRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   objReq.GetId().GetId(),
				ServingData: errdetails.Serialize(objReq),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + vald.GetObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	var removeErrs error
	var errMutex, mutex sync.Mutex
	var successTargets sync.Map
	ce := &payload.Object_Location{
		Uuid: req.GetId().GetId(),
		Ips:  make([]string, 0),
	}

	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.RemoveRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		loc, err := s.remove(ctx, vald.NewValdClient(conn), req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveRPCName+" gRPC error response for "+target,
				&errdetails.RequestInfo{
					RequestId:   req.GetId().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + ".BroadCast/" + target,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				},
			)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if st.Code() == codes.NotFound {
				return nil
			}
			errMutex.Lock()
			if removeErrs == nil {
				removeErrs = err
			} else {
				removeErrs = errors.Wrap(removeErrs, err.Error())
			}
			errMutex.Unlock()
			return err
		}
		mutex.Lock()
		ce.Name = loc.GetName()
		ce.Ips = append(ce.Ips, loc.GetIps()...)
		mutex.Unlock()
		successTargets.Store(target, struct{}{})
		return nil
	})
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.RemoveRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetId().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + ".BroadCast",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if removeErrs == nil {
			removeErrs = err
		} else {
			removeErrs = errors.Wrap(removeErrs, err.Error())
		}
	}
	if removeErrs == nil {
		log.Debugf("Remove API succeeded to %#v", ce)
		return ce, nil
	}
	log.Warnf("failed to Remove API: %s", removeErrs.Error())

	// Rollback for Remove RPC.
	var upsertErrs error
	err = s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "rollback/BroadCast/"+target), apiName+"/"+vald.RemoveRPCName+"/rollback/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		if _, ok := successTargets.Load(target); !ok {
			return nil
		}
		objv, ok := oldVecs.Load(target)
		if !ok || objv == nil {
			log.Debug("failed to get old vector from  %s", target)
			return nil
		}

		req := &payload.Upsert_Request{
			Vector: objv.(*payload.Object_Vector),
		}
		_, err := s.upsert(ctx, vald.NewUpsertClient(conn), req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.UpsertRPCName+" for "+vald.RemoveRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + "." + vald.UpsertRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if st.Code() == codes.AlreadyExists && st.Code() == codes.NotFound {
				return nil
			}
			errMutex.Lock()
			if upsertErrs == nil {
				upsertErrs = err
			} else {
				upsertErrs = errors.Wrap(upsertErrs, err.Error())
			}
			errMutex.Unlock()
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.UpsertRPCName+" API for "+vald.RemoveRPCName+" API connection not found", err,
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + "." + vald.UpsertRPCName + ".BroadCast",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if upsertErrs == nil {
			err = upsertErrs
		} else {
			err = errors.Wrap(upsertErrs, err.Error())
		}
	}
	if upsertErrs != nil {
		st, msg, err := status.ParseError(upsertErrs, codes.Internal,
			"failed to parse "+vald.UpsertRPCName+" for "+vald.RemoveRPCName+" gRPC error response",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName + "." + vald.UpsertRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	st, msg, err := status.ParseError(removeErrs, codes.Internal,
		"failed to parse "+vald.RemoveRPCName+" gRPC error response",
		&errdetails.RequestInfo{
			RequestId:   req.GetId().GetId(),
			ServingData: errdetails.Serialize(req),
		},
		&errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RemoveRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		},
	)
	if span != nil {
		span.RecordError(err)
		span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
		span.SetStatus(trace.StatusError, err.Error())
	}
	return nil, err
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.StreamRemoveRPCName), apiName+"/"+vald.StreamRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Remove_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Remove_Request)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamRemoveRPCName+"/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Remove(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.RemoveRPCName+" gRPC error response")
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
		return err
	}
	return nil
}

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (res *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}

	var mutex, errMutex sync.Mutex
	var wg sync.WaitGroup

	for i, r := range reqs.GetRequests() {
		idx, req := i, r
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() error {
			defer wg.Done()
			ti := "errgroup.Go/id-" + req.GetId().GetId()
			ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, ti), apiName+"/"+vald.MultiRemoveRPCName+"/"+ti)
			defer func() {
				if span != nil {
					span.End()
				}
			}()

			loc, err := s.Remove(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.RemoveRPCName+" gRPC error response",
					&errdetails.RequestInfo{
						RequestId:   req.GetId().GetId(),
						ServingData: errdetails.Serialize(req),
					})
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				errMutex.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Wrap(errs, err.Error())
				}
				errMutex.Unlock()
				return nil
			}
			mutex.Lock()
			res.Locations[idx] = loc
			mutex.Unlock()
			return nil
		}))
	}
	wg.Wait()
	if errs != nil {
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse "+vald.MultiRemoveRPCName+" gRPC error response",
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return res, err
	}
	return res, nil
}

func (s *server) GetObject(ctx context.Context, req *payload.Object_VectorRequest) (vec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec, err = s.getObject(ctx, s.vc, req, s.vc.GRPCClient().GetCallOption()...)
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse "+vald.GetObjectRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetId().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return vec, nil
}

func (s *server) getObjects(ctx context.Context, req *payload.Object_VectorRequest) (*sync.Map, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "getObjects"), apiName+"/"+vald.GetObjectRPCName+"/getObjects")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var objErrs error
	var mutex sync.Mutex
	vecs := new(sync.Map)
	err := s.gateway.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.GetObjectRPCName+"/getObjects/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		vec, err := s.getObject(ctx, vald.NewObjectClient(conn), req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.GetObjectRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   req.GetId().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if st.Code() == codes.NotFound {
				return nil
			}
			log.Warn(err)
			mutex.Lock()
			if objErrs == nil {
				objErrs = err
			} else {
				objErrs = errors.Wrap(objErrs, err.Error())
			}
			mutex.Unlock()
			return err
		}
		vecs.Store(target, vec)
		return nil
	})
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(
				vald.GetObjectRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetId().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName + ".BroadCast",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				},
			)
			log.Warn(err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		if objErrs == nil {
			objErrs = err
		} else {
			objErrs = errors.Wrap(objErrs, err.Error())
		}
	}
	if objErrs != nil {
		st, msg, err := status.ParseError(objErrs, codes.Internal,
			"failed to parse "+vald.GetObjectRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetId().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName + "." + "BroadCast",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return vecs, nil
}

func (s *server) getObject(ctx context.Context, client vald.ObjectClient, req *payload.Object_VectorRequest, opts ...grpc.CallOption) (ovec *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "getObject"), apiName+"/"+vald.GetObjectRPCName+"/getObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	ovec, err = client.GetObject(ctx, req, opts...)
	if err != nil {
		reqInfo := &errdetails.RequestInfo{
			RequestId:   req.GetId().GetId(),
			ServingData: errdetails.Serialize(req),
		}
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
		}
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.GetObjectRPCName+" API canceld", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeCancelled(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
			)
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.GetObjectRPCName+" API deadline exceeded", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeDeadlineExceeded(
				errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.GetObjectRPCName,
			)
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(
				vald.GetObjectRPCName+" API connection not found", err, reqInfo, resInfo,
			)
			attrs = trace.StatusCodeInternal(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				"failed to parse "+vald.GetObjectRPCName+" gRPC error response", reqInfo, resInfo,
			)
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		log.Warn("failed to process GetObject request\terror: %s", err.Error())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return ovec, nil
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.ObjectRPCServiceName+"/"+vald.StreamGetObjectRPCName), apiName+"/"+vald.StreamGetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Object_VectorRequest) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Object_VectorRequest)
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamInsertRPCName+"/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.GetObject(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.GetObjectRPCName+" gRPC error response")
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse "+vald.StreamGetObjectRPCName+" gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}
