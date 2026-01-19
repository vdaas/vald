//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

const (
	StatsRPCServiceName        = "Stats"
	ResourceStatsRPCName       = "ResourceStats"
	ResourceStatsDetailRPCName = "ResourceStatsDetail"
)

func (s *server) ResourceStatsDetail(
	ctx context.Context, _ *payload.Empty,
) (detail *payload.Info_Stats_ResourceStatsDetail, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, "rpc.v1."+StatsRPCServiceName+"/"+ResourceStatsDetailRPCName), apiName+"/"+ResourceStatsDetailRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ech := make(chan error, 1)
	var mu sync.Mutex
	detail = &payload.Info_Stats_ResourceStatsDetail{
		Details: make(map[string]*payload.Info_Stats_ResourceStats, s.gateway.GetAgentCount(ctx)),
	}
	s.eg.Go(safety.RecoverFunc(func() error {
		defer close(ech)
		ech <- s.gateway.BroadCast(ctx, service.READ, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
			sctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+ResourceStatsDetailRPCName+"/"+target)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			var stats *payload.Info_Stats_ResourceStats
			stats, err = vc.ResourceStats(sctx, new(payload.Empty), copts...)
			if err != nil {
				var (
					attrs trace.Attributes
					st    *status.Status
					msg   string
					code  codes.Code
				)
				switch {
				case errors.Is(err, context.Canceled),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
					attrs = trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/rpc.v1." + ResourceStatsDetailRPCName + ".BroadCast/" +
							target + " canceled: " + err.Error())
					code = codes.Canceled
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
					attrs = trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/rpc.v1." + ResourceStatsDetailRPCName + ".BroadCast/" +
							target + " deadline_exceeded: " + err.Error())
					code = codes.DeadlineExceeded
				default:
					st, msg, err = status.ParseError(err, codes.NotFound, "error "+ResourceStatsDetailRPCName+" API",
						&errdetails.ResourceInfo{
							ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/rpc.v1." + ResourceStatsDetailRPCName + ".BroadCast/" + target,
							ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
						})
					if st != nil {
						code = st.Code()
					} else {
						code = codes.NotFound
					}
					attrs = trace.FromGRPCStatus(code, msg)
				}
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(attrs...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				if err != nil && st != nil &&
					code != codes.Canceled &&
					code != codes.DeadlineExceeded &&
					code != codes.InvalidArgument &&
					code != codes.NotFound &&
					code != codes.OK &&
					code != codes.Unimplemented {
					return err
				}
				return nil
			}
			if stats != nil {
				mu.Lock()
				detail.Details[target] = stats
				mu.Unlock()
			}
			return nil
		})
		return nil
	}))
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ech:
	}
	if err != nil {
		resInfo := &errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/rpc.v1." + ResourceStatsDetailRPCName,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		}
		var attrs trace.Attributes
		switch {
		case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
			err = status.WrapWithInternal(ResourceStatsDetailRPCName+" API connection not found", err, resInfo)
			attrs = trace.StatusCodeInternal(err.Error())
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(ResourceStatsDetailRPCName+" API canceled", err, resInfo)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(ResourceStatsDetailRPCName+" API deadline exceeded", err, resInfo)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			st, _ := status.FromError(err)
			if st != nil {
				attrs = trace.FromGRPCStatus(st.Code(), st.Message())
			}
		}
		log.Debug(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		return nil, err
	}

	return detail, nil
}
