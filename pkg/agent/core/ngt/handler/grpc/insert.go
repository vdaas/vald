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

// Insert inserts a vector to the NGT.
func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (res *payload.Object_Location, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.InsertRPCName)
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
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	err = s.ngt.InsertWithTime(vec.GetId(), vec.GetVector(), req.GetConfig().GetTimestamp())
	if err != nil {
		var attrs []attribute.KeyValue

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
			attrs = trace.StatusCodeAlreadyExists(err.Error())
		} else if errors.Is(err, errors.ErrFlushingIsInProgress()) {
			err = status.WrapWithAborted("Insert API aborted to process search request due to flushing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Insert",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Warn(err)
			attrs = trace.StatusCodeAlreadyExists(err.Error())
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
			attrs = trace.StatusCodeInvalidArgument(err.Error())
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
			attrs = trace.FromGRPCStatus(st.Code(), msg)
		}
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return s.newLocation(vec.GetId()), nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Insert_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Insert_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamInsertRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Insert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Insert gRPC error response")
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamInsert gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (res *payload.Object_Locations, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiInsertRPCName)
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
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		vmap[vec.GetId()] = vec.GetVector()
		uuids = append(uuids, vec.GetId())
	}
	err = s.ngt.InsertMultiple(vmap)
	if err != nil {
		var attrs []attribute.KeyValue
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
			attrs = trace.StatusCodeAlreadyExists(err.Error())
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
			attrs = trace.StatusCodeInvalidArgument(err.Error())
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
			attrs = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())

		}
		return nil, err
	}
	return s.newLocations(uuids...), nil
}
