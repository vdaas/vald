// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	"sync/atomic"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"go.opentelemetry.io/otel/attribute"
)

// Flush removes all vectors that are indexed and uncommitted in the `vald-agent`.
func (s *server) Flush(ctx context.Context, req *payload.Flush_Request) (*payload.Info_Index_Count, error) {
	_, span := trace.StartSpan(ctx, apiName+"/"+vald.FlushRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err := s.ngt.RegenerateIndexes(ctx)
	if err != nil {
		var attrs []attribute.KeyValue
		if errors.Is(err, errors.ErrFlushingIsInProgress) {
			err = status.WrapWithAborted("Flush API aborted to process search request due to flushing indices is in progress", err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Flush",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				})
			log.Debug(err)
			attrs = trace.StatusCodeAborted(err.Error())
		} else {
			err = status.WrapWithInternal("Flush API failed", err,
				&errdetails.RequestInfo{
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: ngtResourceType + "/ngt.Flush",
					ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
				}, info.Get())
			log.Error(err)
			attrs = trace.StatusCodeInternal(err.Error())
		}
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	var (
		stored     uint32
		uncommited uint32
		indexing   atomic.Value
		saving     atomic.Value
	)
	stored = 0
	uncommited = 0
	indexing.Store(false)
	saving.Store(false)

	cnts := &payload.Info_Index_Count{
		Stored:      atomic.LoadUint32(&stored),
		Uncommitted: atomic.LoadUint32(&uncommited),
		Indexing:    indexing.Load().(bool),
		Saving:      saving.Load().(bool),
	}

	return cnts, nil
}
