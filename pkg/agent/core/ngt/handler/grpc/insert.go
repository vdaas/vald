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
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/strings"
)

// Insert inserts a vector to the NGT.
func (s *server) Insert(
	ctx context.Context, req *payload.Insert_Request,
) (res *payload.Object_Location, err error) {
	var span trace.Span
	if trace.Enabled() {
		_, span = trace.StartSpan(ctx, apiName+"/"+vald.InsertRPCName)
		defer trace.End(span)
	}
	vec := req.GetVector()
	if err = s.validateVectorDimension(span, vald.InsertRPCName, ngtResourceType+"/ngt.Insert",
		vec.GetId(), req, len(vec.GetVector()), s.ngt.GetDimensionSize()); err != nil {
		return nil, err
	}
	if err = s.validateVectorValue(span, vald.InsertRPCName, ngtResourceType+"/ngt.Insert",
		vec.GetId(), req, vec.GetVector()); err != nil {
		return nil, err
	}
	return s.insertValidated(span, req)
}

func (s *server) insertValidated(
	span trace.Span, req *payload.Insert_Request,
) (res *payload.Object_Location, err error) {
	vec := req.GetVector()

	ts := vec.GetTimestamp()
	if ts == 0 {
		ts = req.GetConfig().GetTimestamp()
	}
	if req.GetConfig().GetSkipStrictExistCheck() {
		err = s.ngt.InsertWithConfig(vec.GetId(), vec.GetVector(), ts, true)
	} else if ts != 0 {
		err = s.ngt.InsertWithTime(vec.GetId(), vec.GetVector(), ts)
	} else {
		err = s.ngt.Insert(vec.GetId(), vec.GetVector())
	}
	if err != nil {
		var attrs []attribute.KeyValue
		if errors.Is(err, errors.ErrFlushingIsInProgress) {
			err = status.WrapWithAborted("Insert API aborted to process insert request due to flushing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Insert"))
			log.Warn(err)
			if span != nil {
				attrs = trace.StatusCodeAborted(err.Error())
			}
		} else if errors.Is(err, errors.ErrUUIDAlreadyExists(vec.GetId())) {
			err = status.WrapWithAlreadyExists(fmt.Sprintf("Insert API uuid %s already exists", vec.GetId()), err,
				&errdetails.RequestInfo{
					RequestId:   req.GetVector().GetId(),
					ServingData: errdetails.Serialize(req),
				},
				s.resourceInfo(ngtResourceType+"/ngt.Insert"))
			log.Warn(err)
			if span != nil {
				attrs = trace.StatusCodeAlreadyExists(err.Error())
			}
		} else if errors.Is(err, errors.ErrUUIDNotFound(0)) {
			err = status.WrapWithInvalidArgument(fmt.Sprintf("Insert API invalid id: \"%s\" or vector: %v was given", vec.GetId(), vec.GetVector()), err,
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
				s.resourceInfo(ngtResourceType+"/ngt.Insert"))
			log.Warn(err)
			if span != nil {
				attrs = trace.StatusCodeInvalidArgument(err.Error())
			}
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
				s.resourceInfo(ngtResourceType+"/ngt.Insert"), info.Get())
			if span != nil {
				attrs = trace.FromGRPCStatus(st.Code(), msg)
			}
		}
		if span != nil {
			errhandler.RecordSpanAttrs(span, attrs, err)
		}
		return nil, err
	}
	return s.newLocation(vec.GetId()), nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) (err error) {
	return streamLocation(s, stream, vald.StreamInsertRPCName, func(req *payload.Insert_Request) string {
		return req.GetVector().GetId()
	}, s.Insert)
}

func (s *server) MultiInsert(
	ctx context.Context, reqs *payload.Insert_MultiRequest,
) (res *payload.Object_Locations, err error) {
	var span trace.Span
	if trace.Enabled() {
		_, span = trace.StartSpan(ctx, apiName+"/"+vald.MultiInsertRPCName)
		defer trace.End(span)
	}
	nReqs := len(reqs.GetRequests())
	uuids := make([]string, 0, nReqs)
	var (
		strictMap map[string][]float32
		skipMap   map[string][]float32
		ts        int64
	)
	for i, req := range reqs.GetRequests() {
		vec := req.GetVector()
		if i == 0 {
			ts = vec.GetTimestamp()
			if ts == 0 {
				ts = req.GetConfig().GetTimestamp()
			}
		}
		if err = s.validateVectorDimension(span, vald.MultiInsertRPCName, ngtResourceType+"/ngt.MultiInsert",
			vec.GetId(), req, len(vec.GetVector()), s.ngt.GetDimensionSize()); err != nil {
			return nil, err
		}
		if err = s.validateVectorValue(span, vald.MultiInsertRPCName, ngtResourceType+"/ngt.MultiInsert",
			vec.GetId(), req, vec.GetVector()); err != nil {
			return nil, err
		}
		id := vec.GetId()
		if req.GetConfig().GetSkipStrictExistCheck() {
			if skipMap == nil {
				skipMap = make(map[string][]float32, nReqs)
			}
			skipMap[id] = vec.GetVector()
			delete(strictMap, id)
		} else {
			if strictMap == nil {
				strictMap = make(map[string][]float32, nReqs)
			}
			strictMap[id] = vec.GetVector()
			delete(skipMap, id)
		}
		uuids = append(uuids, id)
	}
	strictInsert := func(m map[string][]float32) error {
		if ts != 0 {
			return s.ngt.InsertMultipleWithTime(m, ts)
		}
		return s.ngt.InsertMultiple(m)
	}
	switch {
	case len(skipMap) == 0:
		err = strictInsert(strictMap)
	case len(strictMap) == 0:
		err = s.ngt.InsertMultipleWithConfig(skipMap, ts, true)
	default:
		if ts == 0 {
			ts = time.Now().UnixNano()
		}
		if err = strictInsert(strictMap); err == nil {
			err = s.ngt.InsertMultipleWithConfig(skipMap, ts, true)
		}
	}
	if err != nil {
		var attrs []attribute.KeyValue
		isErr := errMatcher(err)
		if errors.Is(err, errors.ErrFlushingIsInProgress) {
			err = status.WrapWithAborted("MultiInsert API aborted to process insert request due to flushing indices is in progress", err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				s.resourceInfo(ngtResourceType+"/ngt.MultiInsert"))
			log.Warn(err)
			if span != nil {
				attrs = trace.StatusCodeAborted(err.Error())
			}
		} else if alreadyExistsIDs := func() []string {
			aids := make([]string, 0, len(uuids))
			for _, id := range uuids {
				if isErr(errors.ErrUUIDAlreadyExists(id)) {
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
				s.resourceInfo(ngtResourceType+"/ngt.MultiInsert"))
			log.Warn(err)
			if span != nil {
				attrs = trace.StatusCodeAlreadyExists(err.Error())
			}
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
				s.resourceInfo(ngtResourceType+"/ngt.MultiInsert"))
			log.Warn(err)
			if span != nil {
				attrs = trace.StatusCodeInvalidArgument(err.Error())
			}
		} else {
			err = status.WrapWithInternal("MultiInsert API failed", err,
				&errdetails.RequestInfo{
					RequestId:   strings.Join(uuids, ", "),
					ServingData: errdetails.Serialize(reqs),
				},
				s.resourceInfo(ngtResourceType+"/ngt.MultiInsert"), info.Get())
			log.Error(err)
			if span != nil {
				attrs = trace.StatusCodeInternal(err.Error())
			}
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
