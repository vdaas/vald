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

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (res *payload.Object_Location, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	id := req.GetId()
	uuid := id.GetId()
	if len(uuid) == 0 {
		err = errors.ErrInvalidUUID(uuid)
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
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Warn(err)
		return nil, err
	}
	err = s.ngt.DeleteWithTime(uuid, req.GetConfig().GetTimestamp())
	if err != nil {
		var attrs []attribute.KeyValue
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
			attrs = trace.StatusCodeNotFound(err.Error())
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
			attrs = trace.StatusCodeInvalidArgument(err.Error())
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
			attrs = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return s.newLocation(uuid), nil
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func() interface{} { return new(payload.Remove_Request) },
		func(ctx context.Context, data interface{}) (interface{}, error) {
			req := data.(*payload.Remove_Request)
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamRemoveRPCName+"/id-"+req.GetId().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Remove(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Remove gRPC error response")
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamRemove gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (res *payload.Object_Locations, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiRemoveRPCName)
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
		var attrs []attribute.KeyValue
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
			attrs = trace.StatusCodeNotFound(err.Error())
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
			attrs = trace.StatusCodeInvalidArgument(err.Error())
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
