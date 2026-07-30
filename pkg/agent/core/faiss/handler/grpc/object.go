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
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

func (s *server) Exists(
	ctx context.Context, uid *payload.Object_ID,
) (res *payload.Object_ID, err error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.ExistsRPCName)
	defer trace.End(span)
	uuid := uid.GetId()
	if err = s.validateUUID(span, vald.ExistsRPCName, faissResourceType+"/faiss.Exists",
		uuid, uid); err != nil {
		return nil, err
	}
	if _, ok := s.faiss.Exists(uuid); !ok {
		err = errors.ErrObjectIDNotFound(uid.GetId())
		err = status.WrapWithNotFound(fmt.Sprintf("Exists API meta %s's uuid not found", uid.GetId()), err,
			&errdetails.RequestInfo{
				RequestId:   uid.GetId(),
				ServingData: errdetails.Serialize(uid),
			},
			s.resourceInfo(faissResourceType+"/faiss.Exists"),
			uid.GetId())
		return errhandler.HandleError[payload.Object_ID](span, codes.NotFound, err)
	}
	return uid, nil
}

func (s *server) GetObject(
	ctx context.Context, id *payload.Object_VectorRequest,
) (res *payload.Object_Vector, err error) {
	return s.UnimplementedObjectServer.GetObject(ctx, id)
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) (err error) {
	return s.UnimplementedObjectServer.StreamGetObject(stream)
}
