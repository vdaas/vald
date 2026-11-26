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
	"slices"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
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

func (s *server) Remove(
	ctx context.Context, req *payload.Remove_Request,
) (locs *payload.Object_Location, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
		defer trace.End(span)
	}

	id := req.GetId()
	uuid := id.GetId()
	if !req.GetConfig().GetSkipStrictExistCheck() {
		_, err := s.exists(ctx, uuid)
		if err != nil {
			var attrs trace.Attributes
			switch {
			case errors.Is(err, errors.ErrInvalidUUID(uuid)):
				err = status.WrapWithInvalidArgument(
					vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API invalid argument for uuid \""+uuid+"\" detected",
					err,
					append(s.errDetails(ctx, vald.RemoveRPCName, uuid, req),
						&errdetails.BadRequest{
							FieldViolations: []*errdetails.BadRequestFieldViolation{
								{
									Field:       uuidField,
									Description: err.Error(),
								},
							},
						},
					)...,
				)
				if span != nil {
					attrs = trace.StatusCodeInvalidArgument(err.Error())
				}
			case errors.Is(err, errors.ErrObjectIDNotFound(uuid)):
				err = status.WrapWithNotFound(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API id "+uuid+"'s data not found", err, s.errDetails(ctx, vald.RemoveRPCName, uuid, req)...)
				if span != nil {
					attrs = trace.StatusCodeNotFound(err.Error())
				}
			case errors.Is(err, errGRPCClientConnNotFoundStar):
				err = status.WrapWithInternal(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API connection not found", err, s.errDetails(ctx, vald.RemoveRPCName, uuid, req)...)
				if span != nil {
					attrs = trace.StatusCodeInternal(err.Error())
				}
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API canceled", err, s.errDetails(ctx, vald.RemoveRPCName, uuid, req)...)
				if span != nil {
					attrs = trace.StatusCodeCancelled(err.Error())
				}
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API deadline exceeded", err, s.errDetails(ctx, vald.RemoveRPCName, uuid, req)...)
				if span != nil {
					attrs = trace.StatusCodeDeadlineExceeded(err.Error())
				}
			default:
				if span != nil {
					st, _ := status.FromError(err)
					if st != nil {
						attrs = trace.FromGRPCStatus(st.Code(), st.Message())
					}
				}
			}
			if err != nil {
				if span != nil {
					errhandler.RecordSpanAttrs(span, attrs, err)
				}
				return nil, err
			}
		}
	}
	// The caller's request is read by others while this runs (MultiRemove serializes the
	// whole MultiRequest while sibling Removes are in flight), so never write to it:
	// forward a new request that skips the agents' exist check (done above) and is timestamped.
	freq := req
	if cfg := req.GetConfig(); cfg == nil || !cfg.GetSkipStrictExistCheck() || cfg.GetTimestamp() == 0 {
		ts := cfg.GetTimestamp()
		if ts == 0 {
			ts = time.Now().UnixNano()
		}
		freq = &payload.Remove_Request{
			Id: id,
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: true,
				Timestamp:            ts,
			},
		}
	}
	locs, _, _, err = s.broadcastRemove(ctx, uuid, req, freq)
	if err != nil {
		errhandler.RecordSpanFromError(span, err)
		return nil, err
	}
	if len(locs.Ips) <= 0 {
		err = errors.ErrIndexNotFound
		err = status.WrapWithNotFound(vald.RemoveRPCName+" API remove target not found", err, s.errDetails(ctx, vald.RemoveRPCName, uuid, req)...)
		return errhandler.HandleError[payload.Object_Location](span, codes.NotFound, err)
	}
	return locs, nil
}

// broadcastRemove sends freq to every agent and returns the removed locations
// (named after the sorted agent names), the number of agents that removed uuid,
// each failed agent's latest error by target, and BroadCast's error, which joins those failures
// (an agent that does not hold uuid is not a failure). req is the caller's
// request, used for the error details.
func (s *server) broadcastRemove(
	ctx context.Context, uuid string, req, freq *payload.Remove_Request,
) (locs *payload.Object_Location, removed int, failed map[string]error, err error) {
	var (
		mu            sync.Mutex
		failedTargets map[string]error
		hasErr        atomic.Bool
	)
	locs = &payload.Object_Location{
		Uuid: freq.GetId().GetId(),
		Ips:  make([]string, 0, s.replica),
	}
	ls := make([]string, 0, s.replica)
	err = s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		var span trace.Span
		if trace.Enabled() {
			ctx, span = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.RemoveRPCName+"/"+target)
			defer trace.End(span)
		}
		loc, err := vc.Remove(ctx, freq, copts...)
		if err != nil {
			// ParseError never attaches details to a NotFound status, and an agent that
			// does not hold the object is tolerated below, so skip building them for it.
			if pst, ok := status.FromError(err); ok && pst != nil && pst.Code() == codes.NotFound {
				if span != nil {
					errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(codes.NotFound, pst.Message()), err)
				}
				if hasErr.Load() {
					mu.Lock()
					if failedTargets != nil {
						delete(failedTargets, target)
					}
					mu.Unlock()
				}
				return nil
			}
			details := s.errDetails(ctx, vald.RemoveRPCName, uuid, req)
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveRPCName+" gRPC error response", details...)
			if span != nil {
				errhandler.RecordSpanAttrs(span, trace.FromGRPCStatus(st.Code(), msg), err)
			}
			if err != nil && st.Code() != codes.NotFound {
				log.Error(err)
				mu.Lock()
				if failedTargets == nil {
					failedTargets = make(map[string]error, s.replica)
				}
				failedTargets[target] = err
				hasErr.Store(true)
				mu.Unlock()
				return err
			}
			if hasErr.Load() {
				mu.Lock()
				if failedTargets != nil {
					delete(failedTargets, target)
				}
				mu.Unlock()
			}
			return nil
		}
		mu.Lock()
		if failedTargets != nil {
			delete(failedTargets, target)
		}
		locs.Ips = append(locs.GetIps(), loc.GetIps()...)
		ls = append(ls, loc.GetName())
		mu.Unlock()
		return nil
	})
	slices.Sort(ls)
	locs.Name = strings.Join(ls, ",")
	return locs, len(ls), failedTargets, err
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	return streamLocation(s, stream, vald.RemoveRPCServiceName, vald.StreamRemoveRPCName, func(req *payload.Remove_Request) string {
		return req.GetId().GetId()
	}, s.Remove)
}

func (s *server) MultiRemove(
	ctx context.Context, reqs *payload.Remove_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
	return multiWrite(ctx, s, vald.RemoveRPCServiceName, vald.MultiRemoveRPCName, reqs, reqs.GetRequests(),
		func(r *payload.Remove_Request) (string, string, error) {
			id := r.GetId().GetId()
			if r == nil || id == "" {
				return id, uuidField, errors.ErrInvalidUUID(id)
			}
			return id, "", nil
		}, s.Remove)
}

func (s *server) RemoveByTimestamp(
	ctx context.Context, req *payload.Remove_TimestampRequest,
) (locs *payload.Object_Locations, errs error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveByTimestampRPCName), apiName+"/"+vald.RemoveByTimestampRPCName)
		defer trace.End(span)
	}

	var mu sync.Mutex
	var emu sync.Mutex
	visited := make(map[string]int) // map[uuid: position of locs]
	locs = new(payload.Object_Locations)

	err := s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		sctx := ctx
		var sspan trace.Span
		if trace.Enabled() {
			sctx, sspan = trace.StartSpan(grpc.WithGRPCMethod(ctx, "BroadCast/"+target), apiName+"/removeByTimestamp/BroadCast/"+target)
			defer trace.End(sspan)
		}

		res, err := vc.RemoveByTimestamp(sctx, req, copts...)
		if err != nil {
			if errors.Is(err, errGRPCClientConnNotFoundStar) {
				err = status.WrapWithInternal(
					vald.RemoveByTimestampRPCName+" API connection not found", err,
				)
				errhandler.RecordSpanError(sspan, codes.Internal, err)
				log.Error(err)
				emu.Lock()
				errs = errors.Join(errs, err)
				emu.Unlock()
				return nil
			}
			if pst, ok := status.FromError(err); ok && pst != nil && pst.Code() == codes.NotFound {
				if sspan != nil {
					errhandler.RecordSpanAttrs(sspan, trace.FromGRPCStatus(codes.NotFound, pst.Message()), err)
				}
			} else {
				var (
					st  *status.Status
					msg string
				)
				st, msg, err = status.ParseError(err, codes.Internal,
					vald.RemoveByTimestampRPCName+" gRPC error response",
				)
				if sspan != nil {
					errhandler.RecordSpanAttrs(sspan, trace.FromGRPCStatus(st.Code(), msg), err)
				}
				if err != nil && st.Code() != codes.NotFound {
					log.Error(err)
					emu.Lock()
					errs = errors.Join(errs, err)
					emu.Unlock()
					return nil
				}
			}
		}

		if res != nil && len(res.GetLocations()) > 0 {
			mu.Lock()
			for _, loc := range res.GetLocations() {
				if pos, ok := visited[loc.GetUuid()]; !ok {
					locs.Locations = append(locs.GetLocations(), loc)
					visited[loc.GetUuid()] = len(locs.Locations) - 1
				} else {
					if pos < len(locs.GetLocations()) {
						locs.GetLocations()[pos].Ips = append(locs.GetLocations()[pos].Ips, loc.GetIps()...)
						if s := locs.GetLocations()[pos].Name; len(s) == 0 {
							locs.GetLocations()[pos].Name = loc.GetName()
						} else {
							locs.GetLocations()[pos].Name = s + "," + loc.GetName()
						}
					}
				}
			}
			mu.Unlock()
		}
		return nil
	})
	if errs != nil {
		err = errors.Join(err, errs)
	}
	if err != nil {
		log.Error(err)
		errhandler.RecordSpanFromError(span, err)
		return nil, err
	}
	if len(locs.GetLocations()) == 0 {
		err = status.WrapWithNotFound(
			vald.RemoveByTimestampRPCName+" API remove target not found", errors.ErrIndexNotFound,
			&errdetails.RequestInfo{
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceName: fmt.Sprintf("%s: %s(%s)", apiName, s.name, s.ip),
			},
		)
		log.Error(err)
		return errhandler.HandleError[payload.Object_Locations](span, codes.NotFound, err)
	}
	return locs, nil
}
