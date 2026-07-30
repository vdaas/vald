// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/strings"
)

func (s *server) Update(
	ctx context.Context, req *payload.Update_Request,
) (res *payload.Object_Location, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.UpdateRPCName)
	defer trace.End(span)
	vec := req.GetVector()
	if err = s.validateVectorDimension(span, vald.UpdateRPCName, ngtResourceType+"/ngt.Update",
		vec.GetId(), req, len(vec.GetVector()), s.ngt.GetDimensionSize()); err != nil {
		return nil, err
	}
	uuid := vec.GetId()
	// validateUUID builds the exact same InvalidArgument status this block used
	// to build inline, and additionally records it on the span (the inline
	// version never did, leaving these failures invisible in traces).
	if err = s.validateUUID(span, vald.UpdateRPCName, ngtResourceType+"/ngt.Update", uuid, req); err != nil {
		return nil, err
	}
	err = s.ngt.UpdateWithTime(uuid, vec.GetVector(), req.GetConfig().GetTimestamp())
	if err != nil {
		var attrs []attribute.KeyValue
		if errors.Is(err, errors.ErrFlushingIsInProgress) {
			err = status.WrapWithAborted("Update API aborted to process update request due to flushing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Update"))
			log.Warn(err)
			attrs = trace.StatusCodeAborted(err.Error())
		} else if errors.Is(err, errors.ErrObjectIDNotFound(vec.GetId())) {
			err = status.WrapWithNotFound(fmt.Sprintf("Update API uuid %s not found", vec.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Update"))
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
				s.resourceInfo(ngtResourceType+"/ngt.Update"))
			log.Warn(err)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		} else if errors.Is(err, errors.ErrUUIDAlreadyExists(vec.GetId())) {
			err = status.WrapWithAlreadyExists(fmt.Sprintf("Update API uuid %s's same data already exists", vec.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Update"))
			log.Warn(err)
			attrs = trace.StatusCodeAlreadyExists(err.Error())
		} else {
			err = status.WrapWithInternal("Update API failed", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Update"), info.Get())
			log.Error(err)
			attrs = trace.StatusCodeInternal(err.Error())
		}
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, err
	}
	return s.newLocation(vec.GetId()), nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamUpdateRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Update_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamUpdateRPCName+"/id-"+req.GetVector().GetId())
			defer trace.End(sspan)
			res, err := s.Update(ctx, req)
			if err != nil {
				st, _ := status.FromError(err)
				errhandler.RecordSpanStatus(sspan, st, err)
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
		st, _ := status.FromError(err)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

func (s *server) MultiUpdate(
	ctx context.Context, reqs *payload.Update_MultiRequest,
) (res *payload.Object_Locations, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiUpdateRPCName)
	defer trace.End(span)

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
				s.resourceInfo(ngtResourceType+"/ngt.MultiUpdate"))
			log.Warn(err)
			return errhandler.HandleError[payload.Object_Locations](span, codes.InvalidArgument, err)
		}
		vmap[vec.GetId()] = vec.GetVector()
		uuids = append(uuids, vec.GetId())
	}

	err = s.ngt.UpdateMultiple(vmap)
	if err != nil {
		var attrs []attribute.KeyValue
		if errors.Is(err, errors.ErrFlushingIsInProgress) {
			err = status.WrapWithAborted("MultiUpdate API aborted to process update request due to flushing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				s.resourceInfo(ngtResourceType+"/ngt.MultiUpdate"))
			log.Warn(err)
			attrs = trace.StatusCodeAborted(err.Error())
		} else if notFoundIDs := func() []string {
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
				s.resourceInfo(ngtResourceType+"/ngt.MultiUpdate"))
			log.Warn(err)
			attrs = trace.StatusCodeNotFound(err.Error())
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
				s.resourceInfo(ngtResourceType+"/ngt.MultiUpdate"))
			log.Warn(err)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
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
				s.resourceInfo(ngtResourceType+"/ngt.MultiUpdate"))
			log.Warn(err)
			attrs = trace.StatusCodeAlreadyExists(err.Error())
		} else {
			err = status.WrapWithInternal("MultiUpdate API failed", err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				s.resourceInfo(ngtResourceType+"/ngt.MultiUpdate"), info.Get())
			log.Error(err)
			attrs = trace.StatusCodeInternal(err.Error())
		}
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, err
	}
	return s.newLocations(uuids...), nil
}

func (s *server) UpdateTimestamp(
	ctx context.Context, req *payload.Update_TimestampRequest,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateTimestampRPCName), apiName+"/"+vald.UpdateTimestampRPCName)
	defer trace.End(span)
	uuid := req.GetId()
	reqInfo := &errdetails.RequestInfo{
		RequestId:   uuid,
		ServingData: errdetails.Serialize(req),
	}
	resInfo := s.resourceInfo(errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateTimestampRPCName + "." + vald.GetObjectRPCName)
	if len(uuid) == 0 {
		err = errors.ErrInvalidUUID(uuid)
		err = status.WrapWithInvalidArgument(vald.UpdateTimestampRPCName+" API invalid uuid", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			})
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, err)
	}
	ts := req.GetTimestamp()
	if !req.GetForce() && ts < 0 {
		err = errors.ErrInvalidTimestamp(ts)
		err = status.WrapWithInvalidArgument(vald.UpdateTimestampRPCName+" API invalid vector argument", err, reqInfo, resInfo,
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "timestamp",
						Description: err.Error(),
					},
				},
			}, info.Get())
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, err)
	}
	err = s.ngt.UpdateTimestamp(uuid, ts, req.GetForce())
	if err != nil {
		var attrs []attribute.KeyValue
		if errors.Is(err, errors.ErrFlushingIsInProgress) {
			err = status.WrapWithAborted(vald.UpdateTimestampRPCName+" API aborted to process update request due to flushing indices is in progress", err, reqInfo, resInfo)
			log.Warn(err)
			attrs = trace.StatusCodeAborted(err.Error())
		} else if errors.Is(err, errors.ErrObjectNotFound(nil, uuid)) {
			err = status.WrapWithNotFound(fmt.Sprintf(vald.UpdateTimestampRPCName+" API uuid %s's data not found", uuid), err, reqInfo, resInfo)
			log.Warn(err)
			attrs = trace.StatusCodeNotFound(err.Error())
		} else if errors.Is(err, errors.ErrZeroTimestamp) || errors.Is(err, errors.ErrUUIDNotFound(0)) {
			err = status.WrapWithInvalidArgument(fmt.Sprintf(vald.UpdateTimestampRPCName+" API invalid argument for uuid \"%s\" detected", uuid), err, reqInfo, resInfo,
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "uuid, timestamp",
							Description: err.Error(),
						},
					},
				})
			log.Warn(err)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		} else if errors.Is(err, errors.ErrNewerTimestampObjectAlreadyExists(uuid, ts)) {
			err = status.WrapWithAlreadyExists(fmt.Sprintf(vald.UpdateTimestampRPCName+" API uuid %s's newer timestamp already exists", uuid), err, reqInfo, resInfo)
			log.Warn(err)
			attrs = trace.StatusCodeAlreadyExists(err.Error())
		} else {
			err = status.WrapWithInternal(vald.UpdateTimestampRPCName+" API failed", err, reqInfo, resInfo, info.Get())
			log.Error(err)
			attrs = trace.StatusCodeInternal(err.Error())
		}
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, err
	}
	return s.newLocation(uuid), nil
}
