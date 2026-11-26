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
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

func (s *server) Upsert(
	ctx context.Context, req *payload.Upsert_Request,
) (loc *payload.Object_Location, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpsertRPCServiceName+"/"+vald.UpsertRPCName), apiName+"/"+vald.UpsertRPCName)
		defer trace.End(span)
	}

	vec := req.GetVector()
	uuid := vec.GetId()
	if len(uuid) == 0 {
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, s.invalidUUIDError(ctx, vald.UpsertRPCName, vald.UpsertRPCName, uuid, req))
	}

	vl := len(vec.GetVector())
	if vl < algorithm.MinimumVectorDimensionSize {
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, s.invalidVectorDimError(ctx, vald.UpsertRPCName, vald.UpsertRPCName, uuid, req, vl))
	}
	if err = s.validateVectorValue(ctx, vald.UpsertRPCName, uuid, req, vec.GetVector()); err != nil {
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, err)
	}
	var shouldInsert bool
	if !req.GetConfig().GetSkipStrictExistCheck() {
		ovec, oerr := s.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		var attrs trace.Attributes
		if oerr != nil || ovec == nil {
			st, _ := status.FromError(oerr)
			if st != nil {
				if span != nil {
					attrs = trace.FromGRPCStatus(st.Code(), st.Message())
				}
				if st.Code() == codes.NotFound {
					shouldInsert = true
					oerr = nil
				}
			}
		} else if conv.F32stos(ovec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
			oerr = status.WrapWithAlreadyExists(
				vald.GetObjectRPCName+" API for "+vald.UpsertRPCName+" API ID = "+uuid+"'s same vector data already exists",
				errors.ErrSameVectorAlreadyExists(uuid, ovec.GetVector(), req.GetVector().GetVector()),
				s.errDetails(ctx, vald.UpsertRPCName, uuid, req)...,
			)
			if span != nil {
				attrs = trace.StatusCodeAlreadyExists(oerr.Error())
			}
		}
		if oerr != nil {
			if span != nil {
				errhandler.RecordSpanAttrs(span, attrs, oerr)
			}
			return nil, oerr
		}
	} else {
		ok, cerr := s.existsCheck(ctx, uuid)
		if cerr != nil {
			if cerr = s.handleExistsCheckError(ctx, span, vald.UpsertRPCName, uuid, req, cerr); cerr != nil {
				return nil, cerr
			}
		}
		shouldInsert = !ok
	}

	var operation string
	if shouldInsert {
		operation = vald.InsertRPCName
		loc, err = s.Insert(ctx, &payload.Insert_Request{
			Vector: vec,
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Filters:              req.GetConfig().GetFilters(),
				Timestamp:            req.GetConfig().GetTimestamp(),
			},
		})
	} else {
		operation = vald.UpdateRPCName
		loc, err = s.Update(ctx, &payload.Update_Request{
			Vector: vec,
			Config: &payload.Update_Config{
				SkipStrictExistCheck:  true,
				Filters:               req.GetConfig().GetFilters(),
				Timestamp:             req.GetConfig().GetTimestamp(),
				DisableBalancedUpdate: req.GetConfig().GetDisableBalancedUpdate(),
			},
		})
	}

	if err != nil {
		// Should we use `status.FromError(err)` instead of `status.PraseError(err)` ?
		// It seems `operation` has an important role for tracing, so I don't refactor for now.
		st, msg, perr := status.ParseError(err, codes.Internal,
			"failed to parse "+operation+" for "+vald.UpsertRPCName+" gRPC error response",
			&errdetails.RequestInfo{
				RequestId:   uuid,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpsertRPCName + "." + operation,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			}, info.Get())
		if span != nil {
			errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), perr)
		}
		return nil, perr
	}
	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) (err error) {
	return streamLocation(s, stream, vald.UpsertRPCServiceName, vald.StreamUpsertRPCName, func(req *payload.Upsert_Request) string {
		return req.GetVector().GetId()
	}, s.Upsert)
}

func (s *server) MultiUpsert(
	ctx context.Context, reqs *payload.Upsert_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
	return multiWrite(ctx, s, vald.UpsertRPCServiceName, vald.MultiUpsertRPCName, reqs, reqs.GetRequests(),
		func(r *payload.Upsert_Request) (string, string, error) {
			return validateMultiVector(r.GetVector())
		}, s.Upsert)
}
