// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
)

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (loc *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	vec := req.GetVector()
	if len(vec.GetVector()) != s.ngt.GetDimensionSize() {
		err = errors.ErrIncompatibleDimensionSize(len(vec.GetVector()), int(s.ngt.GetDimensionSize()))
		err = status.WrapWithInvalidArgument("Upsert API Incompatible Dimension Size detected",
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
				ResourceType: ngtResourceType + "/ngt.Upsert",
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
		err = status.WrapWithInvalidArgument(fmt.Sprintf("Upsert API invalid argument for uuid \"%s\" detected", uuid), err,
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
				ResourceType: ngtResourceType + "/ngt.Upsert",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		log.Warn(err)
		return nil, err
	}

	rtName := "/ngt.Upsert"
	_, exists := s.ngt.Exists(req.GetVector().GetId())
	if exists {
		ureq := payload.Update_RequestFromVTPool()
		ureq.Vector = req.GetVector()
		ureq.Config = payload.Update_ConfigFromVTPool()
		ureq.Config.Timestamp = req.GetConfig().GetTimestamp()
		ureq.Config.SkipStrictExistCheck = true
		loc, err = s.Update(ctx, ureq)
		rtName += "/ngt.Update"
	} else {
		ireq := payload.Insert_RequestFromVTPool()
		ireq.Vector = req.GetVector()
		ireq.Config = payload.Insert_ConfigFromVTPool()
		ireq.Config.Timestamp = req.GetConfig().GetTimestamp()
		ireq.Config.SkipStrictExistCheck = true
		loc, err = s.Insert(ctx, ireq)
		rtName += "/ngt.Insert"
	}
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Upsert gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + rtName,
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/"+vald.StreamUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Upsert_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(ctx, apiName+"/"+vald.StreamUpsertRPCName+"/id-"+req.GetVector().GetId())
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := s.Upsert(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Upsert gRPC error response")
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
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse StreamUpsert gRPC error response")
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return err
	}
	return nil
}

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/"+vald.MultiUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	insertReqs := make([]*payload.Insert_Request, 0, len(reqs.GetRequests()))
	updateReqs := make([]*payload.Update_Request, 0, len(reqs.GetRequests()))

	ids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		vec := req.GetVector()
		if len(vec.GetVector()) != s.ngt.GetDimensionSize() {
			err = errors.ErrIncompatibleDimensionSize(len(vec.GetVector()), int(s.ngt.GetDimensionSize()))
			err = status.WrapWithInvalidArgument("MultiUpsert API Incompatible Dimension Size detected",
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
					ResourceType: ngtResourceType + "/ngt.MultiUpsert",
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
		ids = append(ids, vec.GetId())
		_, exists := s.ngt.Exists(vec.GetId())
		if exists {
			ureq := payload.Update_RequestFromVTPool()
			ureq.Vector = vec
			ureq.Config = payload.Update_ConfigFromVTPool()
			ureq.Config.Timestamp = req.GetConfig().GetTimestamp()
			ureq.Config.SkipStrictExistCheck = true
			updateReqs = append(updateReqs, ureq)
		} else {
			ireq := payload.Insert_RequestFromVTPool()
			ireq.Vector = vec
			ireq.Config = payload.Insert_ConfigFromVTPool()
			ireq.Config.Timestamp = req.GetConfig().GetTimestamp()
			ireq.Config.SkipStrictExistCheck = true
			insertReqs = append(insertReqs, ireq)
		}
	}

	switch {
	case len(insertReqs) <= 0:
		ureqs := payload.Update_MultiRequestFromVTPool()
		ureqs.Requests = updateReqs
		res, err = s.MultiUpdate(ctx, ureqs)
	case len(updateReqs) <= 0:
		ireqs := payload.Insert_MultiRequestFromVTPool()
		ireqs.Requests = insertReqs
		res, err = s.MultiInsert(ctx, ireqs)
	default:
		var (
			ures, ires *payload.Object_Locations
			errs       error
			mu         sync.Mutex
			wg         sync.WaitGroup
		)
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() (err error) {
			defer wg.Done()
			ureqs := payload.Update_MultiRequestFromVTPool()
			ureqs.Requests = updateReqs
			ures, err = s.MultiUpdate(ctx, ureqs)
			if err != nil {
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Join(errs, err)
				}
				mu.Unlock()
			}
			return nil
		}))
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() (err error) {
			defer wg.Done()
			ireqs := payload.Insert_MultiRequestFromVTPool()
			ireqs.Requests = insertReqs
			ires, err = s.MultiInsert(ctx, ireqs)
			if err != nil {
				mu.Lock()
				if errs == nil {
					errs = err
				} else {
					errs = errors.Join(errs, err)
				}
				mu.Unlock()
			}
			return nil
		}))
		wg.Wait()

		if errs == nil {
			res = payload.Object_LocationsFromVTPool()
			switch {
			case ures.GetLocations() == nil:
				res.Locations = ires.GetLocations()
			case ires.GetLocations() == nil:
				res.Locations = ures.GetLocations()
			default:
				res.Locations = append(ures.GetLocations(), ires.GetLocations()...)
			}
		} else {
			err = errs
		}

	}
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal,
			"failed to parse MultiUpsert gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   strings.Join(ids, ","),
				ServingData: errdetails.Serialize(reqs),
			},
			&errdetails.ResourceInfo{
				ResourceType: ngtResourceType + "/ngt.MultiUpsert",
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			}, info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return res, nil
}
