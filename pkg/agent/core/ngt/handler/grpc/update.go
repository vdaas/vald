// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (res *payload.Object_Location, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.UpdateRPCName)
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
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	uuid := vec.GetId()
	if len(uuid) == 0 {
		err = errors.ErrInvalidUUID(uuid)
		err = status.WrapWithInvalidArgument(fmt.Sprintf("Update API invalid argument for uuid \"%s\" detected", uuid), err,
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
				ResourceType: ngtResourceType + "/ngt.Update",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		log.Warn(err)
		return nil, err
	}
	cfg := req.GetConfig()
	err = s.ngt.UpdateWithTime(uuid, vec.GetVector(), cfg.GetTimestamp(), cfg.GetUpdateTimestampIfExists())
	if err != nil {
		var attrs []attribute.KeyValue
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
			attrs = trace.StatusCodeNotFound(err.Error())
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
			attrs = trace.StatusCodeInvalidArgument(err.Error())
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
			attrs = trace.StatusCodeAlreadyExists(err.Error())
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
			attrs = trace.StatusCodeInternal(err.Error())
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

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Update_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamUpdateRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Update(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Update gRPC error response")
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamUpdate gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (res *payload.Object_Locations, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	uuids := make([]string, 0, len(reqs.GetRequests()))
	reqMap := make(map[string]*payload.Update_Request, len(reqs.GetRequests()))

	// check dimension and remove duplicated ID
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
				span.RecordError(err)
				span.SetAttributes(trace.StatusCodeInvalidArgument(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		reqMap[vec.GetId()] = req
		uuids = append(uuids, vec.GetId())
	}

	if len(uuids) == 0 {
		return s.newLocations(), nil
	}

	notFoundIDs := make([]string, 0)
	invalidArgumentIDs := make([]string, 0)
	alreadyExistsIDs := make([]string, 0)
	otherIDs := make([]string, 0)

	for _, req := range reqMap {
		vec := req.GetVector()
		id := vec.GetId()
		v := vec.GetVector()

		err := s.ngt.UpdateWithTime(id, v, vec.GetTimestamp(), req.GetConfig().GetUpdateTimestampIfExists())
		if err != nil {
			if errors.Is(err, errors.ErrObjectIDNotFound(id)) {
				notFoundIDs = append(notFoundIDs, id)
			} else if errors.Is(err, errors.ErrInvalidDimensionSize(len(v), s.ngt.GetDimensionSize())) || errors.Is(err, errors.ErrUUIDNotFound(0)) {
				invalidArgumentIDs = append(invalidArgumentIDs, id)
			} else if errors.Is(err, errors.ErrUUIDAlreadyExists(id)) {
				alreadyExistsIDs = append(alreadyExistsIDs, id)
			} else {
				otherIDs = append(otherIDs, id)
			}
		}
	}

	var errs error
	handleErrIDs := func(ids []string, errFunc func() error, attrsFunc func(string) []attribute.KeyValue) {
		if len(ids) != 0 {
			err := errFunc()
			log.Warn(err)
			errs = errors.Join(errs, err)

			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrsFunc(err.Error())...)
				span.SetStatus(trace.StatusError, err.Error())
			}
		}
	}

	handleErrIDs(notFoundIDs, func() error {
		return status.WrapWithNotFound(fmt.Sprintf("MultiUpdate API uuids %v not found", notFoundIDs), err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(notFoundIDs, ", "),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.MultiUpdate",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
	}, trace.StatusCodeNotFound)

	handleErrIDs(invalidArgumentIDs, func() error {
		return status.WrapWithInvalidArgument(fmt.Sprintf("MultiUpdate API invalid argument for uuids \"%v\" detected", invalidArgumentIDs), err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(invalidArgumentIDs, ", "),
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
	}, trace.StatusCodeInvalidArgument)

	handleErrIDs(alreadyExistsIDs, func() error {
		return status.WrapWithAlreadyExists(fmt.Sprintf("MultiUpdate API uuids %v already exists", alreadyExistsIDs), err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(alreadyExistsIDs, ", "),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.MultiUpdate",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
	}, trace.StatusCodeAlreadyExists)

	handleErrIDs(otherIDs, func() error {
		return status.WrapWithInternal("Update API failed", err,
			&errdetails.RequestInfo{
				RequestId:   strings.Join(otherIDs, ", "),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.Update",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
	}, trace.StatusCodeInternal)

	if errs != nil {
		return nil, errs
	}

	return s.newLocations(uuids...), nil
}
