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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
)

func (s *server) Upsert(
	ctx context.Context, req *payload.Upsert_Request,
) (loc *payload.Object_Location, err error) {
	var span trace.Span
	if trace.Enabled() {
		_, span = trace.StartSpan(ctx, apiName+"/"+vald.UpsertRPCName)
		defer trace.End(span)
	}
	vec := req.GetVector()
	if err = s.validateVectorDimension(span, vald.UpsertRPCName, ngtResourceType+"/ngt.Upsert",
		vec.GetId(), req, len(vec.GetVector()), s.ngt.GetDimensionSize()); err != nil {
		return nil, err
	}
	uuid := vec.GetId()
	// validateUUID builds the exact same InvalidArgument status this block used
	// to build inline, and additionally records it on the span (the inline
	// version never did, leaving these failures invisible in traces).
	if err = s.validateUUID(span, vald.UpsertRPCName, ngtResourceType+"/ngt.Upsert", uuid, req); err != nil {
		return nil, err
	}
	if err = s.validateVectorValue(span, vald.UpsertRPCName, ngtResourceType+"/ngt.Upsert",
		vec.GetId(), req, vec.GetVector()); err != nil {
		return nil, err
	}

	rtName := "/ngt.Upsert"
	_, exists := s.ngt.Exists(req.GetVector().GetId())
	if exists {
		loc, err = s.updateValidated(span, &payload.Update_Request{
			Vector: req.GetVector(),
			Config: &payload.Update_Config{
				Timestamp:            req.GetConfig().GetTimestamp(),
				SkipStrictExistCheck: true,
			},
		})
		rtName += "/ngt.Update"
	} else {
		loc, err = s.insertValidated(span, &payload.Insert_Request{
			Vector: req.GetVector(),
			Config: &payload.Insert_Config{
				Timestamp:            req.GetConfig().GetTimestamp(),
				SkipStrictExistCheck: true,
			},
		})
		rtName += "/ngt.Insert"
	}
	if err != nil {
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse Upsert gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   req.GetVector().GetId(),
				ServingData: errdetails.Serialize(req),
			},
			s.resourceInfo(ngtResourceType+rtName))
		if span != nil && span.IsRecording() {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) (err error) {
	return streamLocation(s, stream, vald.StreamUpsertRPCName, func(req *payload.Upsert_Request) string {
		return req.GetVector().GetId()
	}, s.Upsert)
}

func (s *server) MultiUpsert(
	ctx context.Context, reqs *payload.Upsert_MultiRequest,
) (res *payload.Object_Locations, err error) {
	var span trace.Span
	if trace.Enabled() {
		_, span = trace.StartSpan(ctx, apiName+"/"+vald.MultiUpsertRPCName)
		defer trace.End(span)
	}

	insertReqs := make([]*payload.Insert_Request, 0, len(reqs.GetRequests()))
	updateReqs := make([]*payload.Update_Request, 0, len(reqs.GetRequests()))

	ids := make([]string, 0, len(reqs.GetRequests()))
	for _, req := range reqs.GetRequests() {
		vec := req.GetVector()
		if err = s.validateVectorDimension(span, vald.MultiUpsertRPCName, ngtResourceType+"/ngt.MultiUpsert",
			vec.GetId(), req, len(vec.GetVector()), s.ngt.GetDimensionSize()); err != nil {
			return nil, err
		}
		if err = s.validateVectorValue(span, vald.MultiUpsertRPCName, ngtResourceType+"/ngt.MultiUpsert",
			vec.GetId(), req, vec.GetVector()); err != nil {
			return nil, err
		}
		ids = append(ids, vec.GetId())
		_, exists := s.ngt.Exists(vec.GetId())
		if exists {
			updateReqs = append(updateReqs, &payload.Update_Request{
				Vector: vec,
				Config: &payload.Update_Config{
					Timestamp:            req.GetConfig().GetTimestamp(),
					SkipStrictExistCheck: true,
				},
			})
		} else {
			insertReqs = append(insertReqs, &payload.Insert_Request{
				Vector: vec,
				Config: &payload.Insert_Config{
					Timestamp:            req.GetConfig().GetTimestamp(),
					SkipStrictExistCheck: true,
				},
			})
		}
	}

	switch {
	case len(insertReqs) <= 0:
		res, err = s.MultiUpdate(ctx, &payload.Update_MultiRequest{
			Requests: updateReqs,
		})
	case len(updateReqs) <= 0:
		res, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
			Requests: insertReqs,
		})
	default:
		var (
			errs error
			mu   sync.Mutex
			wg   sync.WaitGroup
		)
		wg.Add(1)
		s.eg.Go(safety.RecoverFunc(func() (err error) {
			defer wg.Done()
			_, err = s.MultiUpdate(ctx, &payload.Update_MultiRequest{
				Requests: updateReqs,
			})
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
			_, err = s.MultiInsert(ctx, &payload.Insert_MultiRequest{
				Requests: insertReqs,
			})
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
			// every id was updated or inserted; ids keeps the request order,
			// which concatenating the update and insert results would lose.
			res = s.newLocations(ids...)
		} else {
			err = errs
		}

	}
	if err != nil {
		errhandler.RecordSpanFromError(span, err)
		return nil, err
	}
	return res, nil
}
