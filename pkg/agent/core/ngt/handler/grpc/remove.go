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
	"github.com/vdaas/vald/internal/sync"
)

func (s *server) Remove(
	ctx context.Context, req *payload.Remove_Request,
) (res *payload.Object_Location, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.RemoveRPCName)
	defer trace.End(span)
	id := req.GetId()
	uuid := id.GetId()
	if err = s.validateUUID(span, vald.RemoveRPCName, ngtResourceType+"/ngt.Remove",
		uuid, req); err != nil {
		return nil, err
	}
	err = s.ngt.DeleteWithTime(uuid, req.GetConfig().GetTimestamp())
	if err != nil {
		var attrs []attribute.KeyValue
		if errors.Is(err, errors.ErrFlushingIsInProgress) {
			err = status.WrapWithAborted("Remove API aborted to process remove request due to flushing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Remove"))
			log.Warn(err)
			attrs = trace.StatusCodeAborted(err.Error())
		} else if errors.Is(err, errors.ErrObjectIDNotFound(uuid)) {
			err = status.WrapWithNotFound(fmt.Sprintf("Remove API uuid %s not found", uuid), err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Remove"))
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
				s.resourceInfo(ngtResourceType+"/ngt.Remove"))
			log.Warn(err)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		} else {
			err = status.WrapWithInternal("Remove API failed", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Remove"), info.Get())
			log.Error(err)
			attrs = trace.StatusCodeInternal(err.Error())
		}
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, err
	}
	return s.newLocation(uuid), nil
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamRemoveRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Remove_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamRemoveRPCName+"/id-"+req.GetId().GetId())
			defer trace.End(sspan)
			res, err := s.Remove(ctx, req)
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

func (s *server) MultiRemove(
	ctx context.Context, reqs *payload.Remove_MultiRequest,
) (res *payload.Object_Locations, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiRemoveRPCName)
	defer trace.End(span)
	uuids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		uuids = append(uuids, req.GetId().GetId())
	}
	err = s.ngt.DeleteMultiple(uuids...)
	if err != nil {
		var attrs []attribute.KeyValue
		if errors.Is(err, errors.ErrFlushingIsInProgress) {
			err = status.WrapWithAborted("MultiRemove API aborted to process remove request due to flushing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				s.resourceInfo(ngtResourceType+"/ngt.MultiRemove"))
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
			err = status.WrapWithNotFound(fmt.Sprintf("MultiRemove API uuids %v not found", notFoundIDs), err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				s.resourceInfo(ngtResourceType+"/ngt.MultiRemove"))
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
				s.resourceInfo(ngtResourceType+"/ngt.MultiRemove"))
			log.Warn(err)
			attrs = trace.StatusCodeInvalidArgument(err.Error())
		} else {
			err = status.WrapWithInternal("MultiRemove API failed", err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				s.resourceInfo(ngtResourceType+"/ngt.MultiRemove"), info.Get())
			log.Error(err)
			attrs = trace.StatusCodeInternal(err.Error())
		}
		errhandler.RecordSpanAttrs(span, attrs, err)
		return nil, err
	}
	return s.newLocations(uuids...), nil
}

func (s *server) RemoveByTimestamp(
	ctx context.Context, req *payload.Remove_TimestampRequest,
) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.RemoveByTimestampRPCName)
	defer trace.End(span)

	var mu, emu sync.Mutex
	locs = new(payload.Object_Locations)

	timestampOpFn := timestampOpsFunc(req.GetTimestamps())
	s.ngt.ListObjectFunc(ctx, func(uuid string, oid uint32, timestamp int64) bool {
		if !timestampOpFn(timestamp) {
			return true
		}
		res, err := s.Remove(ctx, &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		if err != nil {
			emu.Lock()
			errs = errors.Join(errs, err)
			emu.Unlock()
		}
		if res != nil {
			mu.Lock()
			locs.Locations = append(locs.Locations, res)
			mu.Unlock()
		}
		return true
	})
	if errs != nil {
		st, _ := status.FromError(errs)
		log.Error(errs)
		errhandler.RecordSpanStatus(span, st, errs)
		return nil, errs
	}
	if locs == nil || len(locs.GetLocations()) == 0 {
		err := status.WrapWithNotFound(
			vald.RemoveByTimestampRPCName+" API remove target not found", errors.ErrIndexNotFound,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(req),
			},
			s.resourceInfo(ngtResourceType+"/ngt.Remove"),
		)
		log.Error(err)
		return errhandler.HandleError[payload.Object_Locations](span, codes.NotFound, err)
	}
	return locs, nil
}

func timestampOpsFunc(ts []*payload.Remove_Timestamp) func(int64) bool {
	fns := make([]func(int64) bool, 0, len(ts))
	for _, t := range ts {
		fns = append(fns, timestampOpFunc(t))
	}
	return func(t int64) bool {
		for _, fn := range fns {
			if !fn(t) {
				return false
			}
		}
		return true
	}
}

func timestampOpFunc(ts *payload.Remove_Timestamp) func(int64) bool {
	switch ts.GetOperator() {
	case payload.Remove_Timestamp_EQ:
		return func(t int64) bool {
			return ts.GetTimestamp() == t
		}
	case payload.Remove_Timestamp_NE:
		return func(t int64) bool {
			return ts.GetTimestamp() != t
		}
	case payload.Remove_Timestamp_GE:
		return func(t int64) bool {
			return ts.GetTimestamp() <= t
		}
	case payload.Remove_Timestamp_GT:
		return func(t int64) bool {
			return ts.GetTimestamp() < t
		}
	case payload.Remove_Timestamp_LE:
		return func(t int64) bool {
			return ts.GetTimestamp() >= t
		}
	case payload.Remove_Timestamp_LT:
		return func(t int64) bool {
			return ts.GetTimestamp() > t
		}
	default:
		return func(timestamp int64) bool {
			return false
		}
	}
}
