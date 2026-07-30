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
	"fmt"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

// validateVectorDimension folds the up-front "incompatible dimension size"
// guard that every unary handler repeats at the top of its method into a single
// call. It returns nil when got == want (the vector dimension matches the
// index), otherwise it builds the exact same InvalidArgument status the call
// sites used to build inline — the "<rpcName> API Incompatible Dimension Size
// detected" message, the RequestInfo carrying id and the serialized request,
// the "vector dimension size" BadRequest field violation, and this server's
// resourceInfo(resourceType) — logs it at Warn, records it on span, and returns
// it. Callers keep ownership of the return decision:
//
//	if err := s.validateVectorDimension(span, vald.InsertRPCName,
//		faissResourceType+"/faiss.Insert", vec.GetId(), req,
//		len(vec.GetVector()), s.faiss.GetDimensionSize()); err != nil {
//		return nil, err
//	}
func (s *server) validateVectorDimension(
	span trace.Span, rpcName, resourceType, id string, servingData any, got, want int,
) error {
	if got == want {
		return nil
	}
	err := errors.ErrIncompatibleDimensionSize(got, want)
	err = status.WrapWithInvalidArgument(rpcName+" API Incompatible Dimension Size detected", err,
		&errdetails.RequestInfo{
			RequestId:   id,
			ServingData: errdetails.Serialize(servingData),
		},
		&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       "vector dimension size",
					Description: err.Error(),
				},
			},
		},
		s.resourceInfo(resourceType))
	log.Warn(err)
	errhandler.RecordSpanError(span, codes.InvalidArgument, err)
	return err
}

// validateUUID folds the up-front "empty uuid" guard into a single call. It
// returns nil when uuid is non-empty, otherwise it builds the exact same
// InvalidArgument status the call sites used to build inline — the
// fmt.Sprintf("%s API invalid argument for uuid %q detected", rpcName, uuid)
// message, the RequestInfo carrying uuid as RequestId and the serialized
// request, the "uuid" BadRequest field violation, and this server's
// resourceInfo(resourceType) — logs it at Warn, records it on span, and returns
// it. Callers keep ownership of the return decision.
func (s *server) validateUUID(
	span trace.Span, rpcName, resourceType, uuid string, servingData any,
) error {
	if len(uuid) != 0 {
		return nil
	}
	err := errors.ErrInvalidUUID(uuid)
	err = status.WrapWithInvalidArgument(fmt.Sprintf("%s API invalid argument for uuid \"%s\" detected", rpcName, uuid), err,
		&errdetails.RequestInfo{
			RequestId:   uuid,
			ServingData: errdetails.Serialize(servingData),
		},
		&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequestFieldViolation{
				{
					Field:       "uuid",
					Description: err.Error(),
				},
			},
		},
		s.resourceInfo(resourceType))
	log.Warn(err)
	errhandler.RecordSpanError(span, codes.InvalidArgument, err)
	return err
}
