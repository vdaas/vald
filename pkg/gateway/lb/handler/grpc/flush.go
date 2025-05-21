//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package grpc

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

func (s *server) Flush(
	ctx context.Context, req *payload.Flush_Request,
) (cnts *payload.Info_Index_Count, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.FlushRPCServiceName+"/"+vald.FlushRPCName), apiName+"/"+vald.FlushRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var (
		stored      uint32
		uncommitted uint32
		indexing    atomic.Value
		saving      atomic.Value
	)
	indexing.Store(false)
	saving.Store(false)
	now := time.Now().UnixNano()
	err = s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(ctx, apiName+"."+vald.FlushRPCName+"/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()
		cnt, err := vc.Flush(ctx, req, copts...)
		if err != nil {
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.FlushRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   strconv.FormatInt(now, 10),
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.FlushRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				})
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if err != nil && st.Code() != codes.NotFound {
				log.Error(err)
				return err
			}
			return nil
		}
		atomic.AddUint32(&stored, cnt.Stored)
		atomic.AddUint32(&uncommitted, cnt.Uncommitted)
		if cnt.Indexing {
			indexing.Store(cnt.Indexing)
		}
		if cnt.Saving {
			saving.Store(cnt.Saving)
		}
		return nil
	})
	if err != nil {
		st, _ := status.FromError(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), st.Message())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}
	cnts = &payload.Info_Index_Count{
		Stored:      atomic.LoadUint32(&stored),
		Uncommitted: atomic.LoadUint32(&uncommitted),
		Indexing:    indexing.Load().(bool),
		Saving:      saving.Load().(bool),
	}
	if cnts.Stored > 0 || cnts.Uncommitted > 0 || cnts.Indexing || cnts.Saving {
		err = errors.Errorf(
			"stored index: %d, uncommitted: %d, indexing: %t, saving: %t",
			cnts.Stored, cnts.Uncommitted, cnts.Indexing, cnts.Saving,
		)
		err = status.WrapWithInternal(vald.FlushRPCName+" API flush failed", err,
			&errdetails.RequestInfo{
				RequestId:   strconv.FormatInt(now, 10),
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.FlushRPCName,
				ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
			})
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	return cnts, nil
}
