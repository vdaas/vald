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

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

// validateVectorValue returns an InvalidArgument error when vec contains NaN or Inf,
// nil otherwise. Update removes the stored vector from every agent before it inserts
// the new one, so relying on the agents' own NaN/Inf rejection would surface only
// after the old vector is already deleted; the gateway has to reject the request
// before it touches any agent.
func (s *server) validateVectorValue(
	ctx context.Context, rpcName, uuid string, req any, vec []float32,
) error {
	if !algorithm.HasNaNOrInf(vec) {
		return nil
	}
	err := errors.ErrNaNOrInfVector
	return status.WrapWithInvalidArgument(rpcName+" API vector with NaN or Inf value detected", err,
		append(s.errDetails(ctx, rpcName, uuid, req),
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       vectorField,
						Description: err.Error(),
					},
				},
			}, info.Get())...)
}

func (s *server) invalidUUIDError(
	ctx context.Context, rpcName, errDetailsRPC, uuid string, req any,
) error {
	err := errors.ErrInvalidMetaDataConfig
	return status.WrapWithInvalidArgument(rpcName+" API invalid uuid", err,
		append(s.errDetails(ctx, errDetailsRPC, uuid, req),
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			})...)
}

func (s *server) invalidVectorDimError(
	ctx context.Context, rpcName, errDetailsRPC, uuid string, req any, vl int,
) error {
	err := errors.ErrInvalidDimensionSize(vl, 0)
	return status.WrapWithInvalidArgument(rpcName+" API invalid vector argument", err,
		append(s.errDetails(ctx, errDetailsRPC, uuid, req),
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			}, info.Get())...)
}

func (s *server) handleExistsCheckError(
	ctx context.Context, span trace.Span, rpcName, uuid string, req any, cerr error,
) error {
	var attrs trace.Attributes
	switch {
	case errors.Is(cerr, errGRPCClientConnNotFoundStar):
		cerr = status.WrapWithInternal(vald.ExistsRPCName+" API for "+rpcName+" API connection not found", cerr, s.errDetails(ctx, rpcName, uuid, req)...)
		if span != nil {
			attrs = trace.StatusCodeInternal(cerr.Error())
		}
	case errors.Is(cerr, context.Canceled):
		cerr = status.WrapWithCanceled(vald.ExistsRPCName+" API for "+rpcName+" API canceled", cerr, s.errDetails(ctx, rpcName, uuid, req)...)
		if span != nil {
			attrs = trace.StatusCodeCancelled(cerr.Error())
		}
	case errors.Is(cerr, context.DeadlineExceeded):
		cerr = status.WrapWithDeadlineExceeded(vald.ExistsRPCName+" API for "+rpcName+" API deadline exceeded", cerr, s.errDetails(ctx, rpcName, uuid, req)...)
		if span != nil {
			attrs = trace.StatusCodeDeadlineExceeded(cerr.Error())
		}
	default:
		return nil
	}
	if span != nil {
		errhandler.RecordSpanAttrs(span, attrs, cerr)
	}
	return cerr
}
