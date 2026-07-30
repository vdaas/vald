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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

func (s *server) CreateIndex(
	ctx context.Context, c *payload.Control_CreateIndexRequest,
) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".CreateIndex")
	defer trace.End(span)
	res = new(payload.Empty)
	err = s.faiss.CreateIndex(ctx)
	if err != nil {
		if errors.Is(err, errors.ErrUncommittedIndexNotFound) {
			err = status.WrapWithFailedPrecondition("CreateIndex API failed", err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(c),
				},
				s.resourceInfo(faissResourceType+"/faiss.CreateIndex"),
				&errdetails.PreconditionFailure{
					Violations: []*errdetails.PreconditionFailureViolation{
						{
							Type:    "uncommitted index is empty",
							Subject: "failed to CreateIndex operation caused by empty uncommitted indices",
						},
					},
				}, info.Get())
			return errhandler.HandleError[payload.Empty](span, codes.FailedPrecondition, err)
		}
		log.Error(err)
		err = status.WrapWithInternal("CreateIndex API failed", err,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(c),
			},
			s.resourceInfo(faissResourceType+"/faiss.CreateIndex"), info.Get())
		log.Error(err)
		return errhandler.HandleError[payload.Empty](span, codes.Internal, err)
	}
	return res, nil
}

func (s *server) SaveIndex(ctx context.Context, _ *payload.Empty) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".SaveIndex")
	defer trace.End(span)
	res = new(payload.Empty)
	err = s.faiss.SaveIndex(ctx)
	if err != nil {
		log.Error(err)
		err = status.WrapWithInternal("SaveIndex API failed to save indices", err,
			s.resourceInfo(faissResourceType+"/faiss.SaveIndex"), info.Get())
		log.Error(err)
		return errhandler.HandleError[payload.Empty](span, codes.Internal, err)
	}
	return res, nil
}

func (s *server) CreateAndSaveIndex(
	ctx context.Context, c *payload.Control_CreateIndexRequest,
) (res *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+".CreateAndSaveIndex")
	defer trace.End(span)
	res = new(payload.Empty)
	err = s.faiss.CreateAndSaveIndex(ctx)
	if err != nil {
		if errors.Is(err, errors.ErrUncommittedIndexNotFound) {
			err = status.WrapWithFailedPrecondition(fmt.Sprintf("CreateAndSaveIndex API failed to create indexes pool_size = %d", c.GetPoolSize()), err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(c),
				},
				s.resourceInfo(faissResourceType+"/faiss.CreateAndSaveIndex"),
				&errdetails.PreconditionFailure{
					Violations: []*errdetails.PreconditionFailureViolation{
						{
							Type:    "uncommitted index is empty",
							Subject: "failed to CreateAndSaveIndex operation caused by empty uncommitted indices",
						},
					},
				}, info.Get())
			return errhandler.HandleError[payload.Empty](span, codes.FailedPrecondition, err)
		}
		err = status.WrapWithInternal(fmt.Sprintf("CreateAndSaveIndex API failed to create indexes pool_size = %d", c.GetPoolSize()), err,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(c),
			},
			s.resourceInfo(faissResourceType+"/faiss.CreateAndSaveIndex"), info.Get())
		log.Error(err)
		return errhandler.HandleError[payload.Empty](span, codes.Internal, err)
	}
	return res, nil
}

func (s *server) IndexInfo(
	ctx context.Context, c *payload.Empty,
) (res *payload.Info_Index_Count, err error) {
	_, span := trace.StartSpan(ctx, apiName+".IndexInfo")
	defer trace.End(span)

	return &payload.Info_Index_Count{
		Stored:      uint32(s.faiss.Len()),
		Uncommitted: uint32(s.faiss.InsertVQueueBufferLen() + s.faiss.DeleteVQueueBufferLen()),
		Indexing:    s.faiss.IsIndexing(),
		Saving:      s.faiss.IsSaving(),
	}, nil
}
