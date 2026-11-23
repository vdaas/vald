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
	"math"
	"slices"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

func (s *server) broadcastUpdate(
	ctx context.Context,
	rpcName, uuid string,
	call func(ctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Object_Location, error),
) (
	locs *payload.Object_Location,
	ls []string,
	visited map[string]bool,
	updated, aeCount int,
	err error,
) {
	var (
		mu     sync.Mutex
		ferr   map[string]error
		hasErr atomic.Bool
	)
	ls = make([]string, 0, s.replica)
	visited = make(map[string]bool, s.replica)
	locs = &payload.Object_Location{
		Uuid: uuid,
		Ips:  make([]string, 0, s.replica),
	}
	err = s.gateway.BroadCast(ctx, service.WRITE, func(bctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		var sspan trace.Span
		if trace.Enabled() {
			bctx, sspan = trace.StartSpan(grpc.WrapGRPCMethod(bctx, "BroadCast/"+target), apiName+"/"+rpcName+"/"+target)
			defer trace.End(sspan)
		}
		loc, callErr := call(bctx, vc, copts...)
		if callErr != nil {
			st, ok := status.FromError(callErr)
			if ok && st != nil {
				if st.Code() != codes.AlreadyExists &&
					st.Code() != codes.Canceled &&
					st.Code() != codes.DeadlineExceeded &&
					st.Code() != codes.InvalidArgument &&
					st.Code() != codes.NotFound &&
					st.Code() != codes.OK &&
					st.Code() != codes.Unimplemented {
					if sspan != nil {
						errhandler.RecordSpanAttrs(sspan, trace.FromGRPCStatus(st.Code(), fmt.Sprintf("%s operation for Agent %s failed,\terror: %v", rpcName, target, callErr)), callErr)
					}
					mu.Lock()
					if ferr == nil {
						ferr = make(map[string]error, s.replica)
					}
					ferr[target] = callErr
					hasErr.Store(true)
					mu.Unlock()
					return callErr
				}
				if st.Code() == codes.AlreadyExists {
					host, _, herr := net.SplitHostPort(target)
					if herr != nil {
						host = target
					}
					mu.Lock()
					if ferr != nil {
						delete(ferr, target)
					}
					aeCount++
					visited[target] = true
					locs.Ips = append(locs.GetIps(), host)
					ls = append(ls, host)
					mu.Unlock()
					return nil
				}
			}
			if hasErr.Load() {
				mu.Lock()
				if ferr != nil {
					delete(ferr, target)
				}
				mu.Unlock()
			}
			return nil
		}
		if loc != nil {
			mu.Lock()
			if ferr != nil {
				delete(ferr, target)
			}
			updated++
			visited[target] = true
			locs.Ips = append(locs.GetIps(), loc.GetIps()...)
			ls = append(ls, loc.GetName())
			mu.Unlock()
		} else if hasErr.Load() {
			mu.Lock()
			if ferr != nil {
				delete(ferr, target)
			}
			mu.Unlock()
		}
		return nil
	})
	err = s.partialBroadCastResult(ctx, rpcName, uuid, err, updated+aeCount, ferr)
	return locs, ls, visited, updated, aeCount, err
}

func (s *server) topUpUpdateReplicas(
	ctx context.Context,
	shortage int,
	visited map[string]bool,
	ireq *payload.Insert_Request,
	locs *payload.Object_Location,
	ls *[]string,
) error {
	var mu sync.Mutex
	return s.gateway.DoMulti(ctx, shortage, func(dctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error {
		if visited[target] {
			return errors.Errorf("target: %s already inserted will skip", target)
		}
		var sspan trace.Span
		if trace.Enabled() {
			dctx, sspan = trace.StartSpan(grpc.WrapGRPCMethod(dctx, "DoMulti/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
			defer trace.End(sspan)
		}
		loc, err := vc.Insert(dctx, ireq, copts...)
		if err != nil {
			if st, ok := status.FromError(err); ok && st != nil && sspan != nil {
				errhandler.RecordSpanAttrs(sspan, trace.FromGRPCStatus(st.Code(), fmt.Sprintf("Shortage index Insert for Update operation for Agent %s failed,\terror: %v", target, err)), err)
			}
			return err
		}
		if loc != nil {
			mu.Lock()
			locs.Ips = append(locs.GetIps(), loc.GetIps()...)
			*ls = append(*ls, loc.GetName())
			mu.Unlock()
		}
		return nil
	})
}

func (s *server) Update(
	ctx context.Context, req *payload.Update_Request,
) (res *payload.Object_Location, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateRPCName), apiName+"/"+vald.UpdateRPCName)
		defer trace.End(span)
	}
	uuid := req.GetVector().GetId()
	if len(uuid) == 0 {
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, s.invalidUUIDError(ctx, vald.UpdateRPCName, vald.UpdateRPCName+"."+vald.GetObjectRPCName, uuid, req))
	}
	vec := req.GetVector().GetVector()
	vl := len(vec)
	if vl < algorithm.MinimumVectorDimensionSize {
		return errhandler.HandleError[payload.Object_Location](
			span,
			codes.InvalidArgument,
			s.invalidVectorDimError(ctx, vald.UpdateRPCName, vald.UpdateRPCName+"."+vald.GetObjectRPCName, uuid, req, vl),
		)
	}
	if err = s.validateVectorValue(ctx, vald.UpdateRPCName, uuid, req, vec); err != nil {
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, err)
	}
	// Update inserts at timestamp+1 after removing at timestamp, which would wrap
	// to a negative value for MaxInt64 and silently drop the insert.
	if ts := req.GetConfig().GetTimestamp(); ts == math.MaxInt64 {
		err = errors.ErrInvalidTimestamp(ts)
		err = status.WrapWithInvalidArgument(vald.UpdateRPCName+" API invalid timestamp", err, append(s.errDetails(ctx, vald.UpdateRPCName, uuid, req),
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "timestamp",
						Description: err.Error(),
					},
				},
			})...)
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, err)
	}

	if req.GetConfig().GetDisableBalancedUpdate() {
		locs, ls, visited, updated, aeCount, berr := s.broadcastUpdate(
			ctx,
			vald.UpdateRPCName,
			uuid,
			func(bctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Object_Location, error) {
				return vc.Update(bctx, req, copts...)
			},
		)
		switch {
		case berr != nil:
			st, msg, perr := status.ParseError(berr, codes.Internal,
				"failed to parse "+vald.UpdateRPCName+" gRPC error response", append(s.errDetails(ctx, vald.UpdateRPCName+"."+vald.GetObjectRPCName, uuid, req), info.Get())...)
			if span != nil {
				errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), perr)
			}
			return nil, perr
		case len(locs.Ips) <= 0:
			berr = status.WrapWithNotFound(vald.UpdateRPCName+" API update target not found", errors.ErrIndexNotFound, s.errDetails(ctx, vald.UpdateRPCName+"."+vald.GetObjectRPCName, uuid, req)...)
			return errhandler.HandleError[payload.Object_Location](span, codes.NotFound, berr)
		case updated+aeCount < s.replica && func() bool {
			agents := s.gateway.GetAgentCount(ctx)
			return agents == 0 || updated+aeCount < agents
		}():
			// Agents' Update deletes at the config timestamp and inserts at +1
			// (MaxInt64 is rejected above); 0 lets each agent pick its own clock.
			its := req.GetConfig().GetTimestamp()
			if its > 0 {
				its++
			}
			// Agents prefer the vector's timestamp; give the top-up replicas the
			// timestamp (config timestamp + 1) the Update call above stored on the others.
			ireq := &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:        uuid,
					Vector:    vec,
					Timestamp: its,
				},
				Config: &payload.Insert_Config{
					SkipStrictExistCheck: true,
					Filters:              req.GetConfig().GetFilters(),
					Timestamp:            its,
				},
			}
			if berr = s.topUpUpdateReplicas(ctx, s.replica-(updated+aeCount), visited, ireq, locs, &ls); berr != nil {
				errhandler.RecordSpanFromError(span, berr)
				return nil, berr
			}
		case updated == 0 && aeCount > 0:
			berr = errors.ErrSameVectorAlreadyExists(uuid, vec, vec)
			berr = status.WrapWithAlreadyExists(vald.UpdateRPCName+" API update target same vector already exists", berr, s.errDetails(ctx, vald.UpdateRPCName+"."+vald.GetObjectRPCName, uuid, req)...)
			return errhandler.HandleError[payload.Object_Location](span, codes.AlreadyExists, berr)
		}
		slices.Sort(ls)
		locs.Name = strings.Join(ls, ",")
		return locs, nil
	}

	if !req.GetConfig().GetSkipStrictExistCheck() {
		ovec, gerr := s.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		if gerr != nil {
			errhandler.RecordSpanFromError(span, gerr)
			return nil, gerr
		}
		if conv.F32stos(ovec.GetVector()) == conv.F32stos(req.GetVector().GetVector()) {
			if ovec.GetTimestamp() < req.GetVector().GetTimestamp() {
				return s.UpdateTimestamp(ctx, &payload.Update_TimestampRequest{
					Id:        uuid,
					Timestamp: req.GetVector().GetTimestamp(),
				})
			}
			serr := errors.ErrSameVectorAlreadyExists(uuid, ovec.GetVector(), req.GetVector().GetVector())
			st, msg, perr := status.ParseError(serr, codes.AlreadyExists,
				"error "+vald.UpdateRPCName+" API ID = "+uuid+"'s same vector data already exists",
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(req),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.GetObjectRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				}, info.Get())
			if span != nil {
				errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), perr)
			}
			return nil, perr
		}
	}
	var now int64
	if req.GetConfig().GetTimestamp() > 0 {
		now = req.GetConfig().GetTimestamp()
	} else {
		now = time.Now().UnixNano()
	}

	rreq := &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: uuid,
		},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: true,
			Timestamp:            now,
		},
	}
	// An agent failing to remove the old copy does not abort the update: the copy
	// inserted below at now+1 takes precedence, and the index correction removes
	// the stale one. Only a missing connection or a done request context is fatal.
	rlocs, removed, failed, err := s.broadcastRemove(ctx, uuid, rreq, rreq)
	switch {
	case err == nil && len(rlocs.GetIps()) <= 0:
		err = errors.ErrIndexNotFound
		err = status.WrapWithNotFound(vald.RemoveRPCName+" API remove target not found", err, s.errDetails(ctx, vald.RemoveRPCName, uuid, rreq)...)
	case err != nil && !errors.Is(err, errGRPCClientConnNotFoundStar) && ctx.Err() == nil:
		if removed > 0 {
			s.warnPartialWrite(vald.RemoveRPCName+" for "+vald.UpdateRPCName, uuid, removed, failed, err)
		} else {
			// Not a partial write (no agent removed the old copy), so it is not counted;
			// the insert below still decides the Update's result.
			log.Warnf("%s API for uuid %s failed on every agent holding it (%d agent(s)), errors: %v",
				vald.RemoveRPCName+" for "+vald.UpdateRPCName, uuid, len(failed), err)
		}
		err = nil
	}
	if err != nil {
		if errors.Is(err, errGRPCClientConnNotFoundStar) {
			err = status.WrapWithInternal(vald.RemoveRPCName+" for "+vald.UpdateRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(rreq),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.RemoveRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				})
			return errhandler.HandleError[payload.Object_Location](span, codes.Internal, err)
		}

		errhandler.RecordSpanFromError(span, err)
		return nil, err
	}
	now++
	// Agents prefer the vector's own timestamp over the config's, so the client's
	// vector timestamp must not reach them: an insert older than the delete queued
	// above would be dropped while still reporting success.
	ireq := &payload.Insert_Request{
		Vector: &payload.Object_Vector{
			Id:        uuid,
			Vector:    req.GetVector().GetVector(),
			Timestamp: now,
		},
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: true,
			Filters:              req.GetConfig().GetFilters(),
			Timestamp:            now,
		},
	}
	res, err = s.Insert(ctx, ireq)
	if err != nil {
		if errors.Is(err, errGRPCClientConnNotFoundStar) {
			err = status.WrapWithInternal(vald.InsertRPCName+" for "+vald.UpdateRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(ireq),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.UpdateRPCName + "." + vald.InsertRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				})
			return errhandler.HandleError[payload.Object_Location](span, codes.Internal, err)
		}
		errhandler.RecordSpanFromError(span, err)
		return nil, err
	}
	return res, nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) (err error) {
	return streamLocation(s, stream, vald.UpdateRPCServiceName, vald.StreamUpdateRPCName, func(req *payload.Update_Request) string {
		return req.GetVector().GetId()
	}, s.Update)
}

func (s *server) MultiUpdate(
	ctx context.Context, reqs *payload.Update_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
	return multiWrite(ctx, s, vald.UpdateRPCServiceName, vald.MultiUpdateRPCName, reqs, reqs.GetRequests(),
		func(r *payload.Update_Request) (string, string, error) {
			return validateMultiVector(r.GetVector())
		}, s.Update)
}

func (s *server) UpdateTimestamp(
	ctx context.Context, req *payload.Update_TimestampRequest,
) (res *payload.Object_Location, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.UpdateRPCServiceName+"/"+vald.UpdateTimestampRPCName), apiName+"/"+vald.UpdateTimestampRPCName)
		defer trace.End(span)
	}
	uuid := req.GetId()
	if len(uuid) == 0 {
		return errhandler.HandleError[payload.Object_Location](
			span,
			codes.InvalidArgument,
			s.invalidUUIDError(ctx, vald.UpdateTimestampRPCName, vald.UpdateTimestampRPCName+"."+vald.GetObjectRPCName, uuid, req),
		)
	}
	ts := req.GetTimestamp()
	if ts < 0 {
		err = errors.ErrInvalidTimestamp(ts)
		err = status.WrapWithInvalidArgument(vald.UpdateTimestampRPCName+" API invalid vector argument", err,
			append(s.errDetails(ctx, vald.UpdateTimestampRPCName+"."+vald.GetObjectRPCName, uuid, req),
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       "timestamp",
							Description: err.Error(),
						},
					},
				}, info.Get())...)
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, err)
	}
	locs, ls, visited, updated, aeCount, berr := s.broadcastUpdate(
		ctx,
		vald.UpdateTimestampRPCName,
		uuid,
		func(bctx context.Context, vc vald.Client, copts ...grpc.CallOption) (*payload.Object_Location, error) {
			return vc.UpdateTimestamp(bctx, req, copts...)
		},
	)
	switch {
	case berr != nil:
		errhandler.RecordSpanFromError(span, berr)
		return nil, berr
	case len(locs.Ips) <= 0:
		berr = status.WrapWithNotFound(
			vald.UpdateTimestampRPCName+" API update target not found",
			errors.ErrIndexNotFound,
			s.errDetails(ctx, vald.UpdateTimestampRPCName+"."+vald.GetObjectRPCName, uuid, req)...)
		return errhandler.HandleError[payload.Object_Location](span, codes.NotFound, berr)
	case updated+aeCount < s.replica && func() bool {
		agents := s.gateway.GetAgentCount(ctx)
		return agents == 0 || updated+aeCount < agents
	}():
		ovec, gerr := s.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: uuid,
			},
		})
		if gerr != nil {
			errhandler.RecordSpanFromError(span, gerr)
			return nil, gerr
		}
		// ovec carries the stored (old) timestamp, which agents prefer over the config's.
		ireq := &payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:        uuid,
				Vector:    ovec.GetVector(),
				Timestamp: ts,
			},
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Timestamp:            ts,
			},
		}
		if berr = s.topUpUpdateReplicas(ctx, s.replica-(updated+aeCount), visited, ireq, locs, &ls); berr != nil {
			errhandler.RecordSpanFromError(span, berr)
			return nil, berr
		}
	case updated == 0 && aeCount > 0:
		berr = status.WrapWithAlreadyExists(vald.UpdateTimestampRPCName+" API update target same vector already exists", errors.ErrSameVectorAlreadyExists(uuid, nil, nil),
			s.errDetails(ctx, vald.UpdateTimestampRPCName+"."+vald.GetObjectRPCName, uuid, req)...)
		return errhandler.HandleError[payload.Object_Location](span, codes.AlreadyExists, berr)
	}
	slices.Sort(ls)
	locs.Name = strings.Join(ls, ",")
	return locs, nil
}
