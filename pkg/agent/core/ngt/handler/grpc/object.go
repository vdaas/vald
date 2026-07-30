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
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
)

func (s *server) Exists(
	ctx context.Context, uid *payload.Object_ID,
) (res *payload.Object_ID, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.ExistsRPCName)
	defer trace.End(span)
	uuid := uid.GetId()
	if err = s.validateUUID(span, vald.ExistsRPCName, ngtResourceType+"/ngt.Exists",
		uuid, uid); err != nil {
		return nil, err
	}
	if _, ok := s.ngt.Exists(uuid); !ok {
		err = status.New(codes.NotFound, errors.ErrObjectIDNotFound(uid.GetId()).Error()).Err()
		return errhandler.HandleError[payload.Object_ID](span, codes.NotFound, err)
	}
	return uid, nil
}

func (s *server) GetObject(
	ctx context.Context, id *payload.Object_VectorRequest,
) (res *payload.Object_Vector, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.GetObjectRPCName)
	defer trace.End(span)
	uuid := id.GetId().GetId()
	if err = s.validateUUID(span, vald.GetObjectRPCName, ngtResourceType+"/ngt.GetObject",
		uuid, id); err != nil {
		return nil, err
	}
	vec, ts, err := s.ngt.GetObject(uuid)
	if err != nil || vec == nil {
		err = status.New(codes.NotFound, errors.ErrObjectNotFound(err, uuid).Error()).Err()
		return errhandler.HandleError[payload.Object_Vector](span, codes.NotFound, err)
	}

	return &payload.Object_Vector{
		Id:        uuid,
		Vector:    vec,
		Timestamp: ts,
	}, nil
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamGetObjectRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Object_VectorRequest) (*payload.Object_StreamVector, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamGetObjectRPCName+"/id-"+req.GetId().GetId())
			defer trace.End(sspan)
			res, err := s.GetObject(ctx, req)
			if err != nil {
				st, _ := status.FromError(err)
				errhandler.RecordSpanStatus(sspan, st, err)
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
		st, _ := status.FromError(err)
		errhandler.RecordSpanStatus(span, st, err)

		log.Error(err)
		return err
	}
	return nil
}

func (s *server) StreamListObject(
	_ *payload.Object_List_Request, stream vald.Object_StreamListObjectServer,
) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamListObjectRPCName)
	defer trace.End(span)

	var (
		mu   sync.Mutex
		emu  sync.RWMutex
		errs = make([]error, 0, s.ngt.Len())
		emap = make(map[string]struct{})
	)
	s.ngt.ListObjectFunc(ctx, func(uuid string, _ uint32, _ int64) bool {
		vec, ts, err := s.ngt.GetObject(uuid)
		var res *payload.Object_List_Response
		if err != nil {
			st := status.CreateWithNotFound(fmt.Sprintf("failed to get object with uuid: %s", uuid), err)
			res = &payload.Object_List_Response{
				Payload: &payload.Object_List_Response_Status{
					Status: st.Proto(),
				},
			}
		} else {
			res = &payload.Object_List_Response{
				Payload: &payload.Object_List_Response_Vector{
					Vector: &payload.Object_Vector{
						Id:        uuid,
						Vector:    vec,
						Timestamp: ts,
					},
				},
			}
		}

		mu.Lock()
		err = stream.Send(res)
		mu.Unlock()

		if err != nil {
			emu.RLock()
			_, ok := emap[err.Error()]
			emu.RUnlock()
			if !ok {
				emu.Lock()
				errs = append(errs, err)
				emap[err.Error()] = struct{}{}
				emu.Unlock()
			}
		}

		// always return true to continue streaming and let the context cancel the Range when stream closes.
		return true
	})

	if len(errs) != 0 {
		// Register all the gRPC codes to the span. Doing this because ParseError cannot parse joined error.
		if span != nil {
			for _, e := range errs {
				st, msg, err := status.ParseError(e, codes.Internal, "failed to parse StreamListObject final gRPC error response")
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, msg)
			}
		}

		// now join all the errors to return
		return errors.Join(errs...)
	}
	return nil
}

// GetTimestamp returns meta information of the object specified by uuid.
// This rpc is only served in AgentServer and not served in LB. Only for internal use mainly for index correction to reduce
// network bandwidth(because vector itself is not required for index correction logic) while processing.
func (s *server) GetTimestamp(
	ctx context.Context, id *payload.Object_TimestampRequest,
) (res *payload.Object_Timestamp, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.GetTimestampRPCName)
	defer trace.End(span)
	uuid := id.GetId().GetId()
	if err = s.validateUUID(span, vald.GetTimestampRPCName, ngtResourceType+"/ngt.GetTimestamp",
		uuid, id); err != nil {
		return nil, err
	}
	_, ts, err := s.ngt.GetObject(uuid)
	if err != nil {
		err = status.New(codes.NotFound, errors.ErrObjectNotFound(err, uuid).Error()).Err()
		return errhandler.HandleError[payload.Object_Timestamp](span, codes.NotFound, err)
	}
	return &payload.Object_Timestamp{
		Id:        uuid,
		Timestamp: ts,
	}, nil
}
