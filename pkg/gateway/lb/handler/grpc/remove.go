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
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

func (s *server) Remove(
	ctx context.Context, req *payload.Remove_Request,
) (locs *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer trace.End(span)

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
				attrs = trace.StatusCodeInvalidArgument(err.Error())
			case errors.Is(err, errors.ErrObjectIDNotFound(uuid)):
				err = status.WrapWithNotFound(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API id "+uuid+"'s data not found", err, s.errDetails(ctx, vald.RemoveRPCName, uuid, req)...)
				attrs = trace.StatusCodeNotFound(err.Error())
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API connection not found", err, s.errDetails(ctx, vald.RemoveRPCName, uuid, req)...)
				attrs = trace.StatusCodeInternal(err.Error())
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API canceled", err, s.errDetails(ctx, vald.RemoveRPCName, uuid, req)...)
				attrs = trace.StatusCodeCancelled(err.Error())
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(vald.ExistsRPCName+" API for "+vald.RemoveRPCName+" API deadline exceeded", err, s.errDetails(ctx, vald.RemoveRPCName, uuid, req)...)
				attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			default:
				st, _ := status.FromError(err)
				if st != nil {
					attrs = trace.FromGRPCStatus(st.Code(), st.Message())
				}
			}
			if err != nil {
				errhandler.RecordSpanAttrs(span, attrs, err)
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
		st, _ := status.FromError(err)
		errhandler.RecordSpanStatus(span, st, err)
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
	var mu sync.Mutex
	// the client backoff may call f again for a target: keep its latest outcome.
	failedTargets := make(map[string]error)
	locs = &payload.Object_Location{
		Uuid: freq.GetId().GetId(),
		Ips:  make([]string, 0, s.replica),
	}
	ls := make([]string, 0, s.replica)
	err = s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BroadCast/"+target), apiName+"/"+vald.RemoveRPCName+"/"+target)
		defer trace.End(span)
		loc, err := vc.Remove(ctx, freq, copts...)
		if err != nil {
			// ParseError never attaches details to a NotFound status, and an agent that
			// does not hold the object is tolerated below, so skip building them for it.
			var details []any
			if pst, ok := status.FromError(err); !ok || pst.Code() != codes.NotFound {
				details = s.errDetails(ctx, vald.RemoveRPCName, uuid, req)
			}
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.RemoveRPCName+" gRPC error response", details...)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if err != nil && st.Code() != codes.NotFound {
				log.Error(err)
				mu.Lock()
				failedTargets[target] = err
				mu.Unlock()
				return err
			}
			mu.Lock()
			delete(failedTargets, target)
			mu.Unlock()
			return nil
		}
		mu.Lock()
		delete(failedTargets, target)
		locs.Ips = append(locs.GetIps(), loc.GetIps()...)
		ls = append(ls, loc.GetName())
		mu.Unlock()
		return nil
	})
	mu.Lock()
	defer mu.Unlock()
	slices.Sort(ls)
	locs.Name = strings.Join(ls, ",")
	return locs, len(ls), failedTargets, err
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.StreamRemoveRPCName), apiName+"/"+vald.StreamRemoveRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Remove_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamRemoveRPCName+"/id-"+req.GetId().GetId())
			defer trace.End(sspan)
			res, err := s.Remove(ctx, req)
			if err != nil {
				st, _ := status.FromError(err)
				errhandler.RecordSpanStatus(sspan, st, err)
				return &payload.Object_StreamLocation{
					Payload: &payload.Object_StreamLocation_Status{
						Status: st.Proto(),
					},
				}, err
			}
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{
					Location: res,
				},
			}, nil
		})
	if err != nil {
		st, _ := status.FromError(err)
		errhandler.RecordSpanStatus(span, st, err)
		return err
	}
	return nil
}

func (s *server) MultiRemove(
	ctx context.Context, reqs *payload.Remove_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
	defer trace.End(span)
	var (
		errList []error
		emu     sync.Mutex
		lmu     sync.Mutex
	)
	eg, ectx := errgroup.New(ctx)
	eg.SetLimit(s.multiLimit(ctx, false))
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs.GetRequests())),
	}
	// The invalid-argument details are identical for every invalid request of
	// the batch, so serialize the whole batch and resolve Addrs once, lazily.
	var (
		detailsBuilt bool
		servingData  string
		resourceName string
	)
	for i, r := range reqs.GetRequests() {
		if r != nil && r.GetId().GetId() != "" {
			idx := i
			req := r
			eg.Go(safety.RecoverFunc(func() (err error) {
				ectx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ectx, "eg.Go"), apiName+"/"+vald.MultiRemoveRPCName+"/id-"+req.GetId().GetId())
				defer trace.End(sspan)
				res, err := s.Remove(ectx, req)
				if err != nil {
					st, _ := status.FromError(err)
					errhandler.RecordSpanStatus(sspan, st, err)
					emu.Lock()
					errList = append(errList, err)
					emu.Unlock()
				} else if res != nil && res.GetUuid() == req.GetId().GetId() && res.GetIps() != nil {
					lmu.Lock()
					locs.Locations[idx] = res
					lmu.Unlock()
				}
				return nil
			}))
		} else {
			err := errors.ErrInvalidUUID(r.GetId().GetId())
			if !detailsBuilt {
				servingData = errdetails.Serialize(reqs)
				resourceName = fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx))
				detailsBuilt = true
			}
			err = status.WrapWithInvalidArgument(vald.MultiRemoveRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   r.GetId().GetId(),
					ServingData: servingData,
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiRemoveRPCName,
					ResourceName: resourceName,
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       uuidField,
							Description: err.Error(),
						},
					},
				}, info.Get())
			emu.Lock()
			errList = append(errList, err)
			emu.Unlock()

		}
	}
	if err := eg.Wait(); err != nil {
		errList = append(errList, err)
	}
	errs = errors.Join(errList...)

	if errs != nil {
		st, _ := status.FromError(errs)
		errhandler.RecordSpanStatus(span, st, errs)
	}

	return locs, errs
}

func (s *server) RemoveByTimestamp(
	ctx context.Context, req *payload.Remove_TimestampRequest,
) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.RemoveRPCServiceName+"/"+vald.RemoveByTimestampRPCName), apiName+"/"+vald.RemoveByTimestampRPCName)
	defer trace.End(span)

	var mu sync.Mutex
	var emu sync.Mutex
	visited := make(map[string]int) // map[uuid: position of locs]
	locs = new(payload.Object_Locations)

	err := s.gateway.BroadCast(ctx, service.WRITE, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		sctx, sspan := trace.StartSpan(grpc.WithGRPCMethod(ctx, "BroadCast/"+target), apiName+"/removeByTimestamp/BroadCast/"+target)
		defer trace.End(sspan)

		res, err := vc.RemoveByTimestamp(sctx, req, copts...)
		if err != nil {
			if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
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
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal,
				vald.RemoveByTimestampRPCName+" gRPC error response",
			)
			if sspan != nil {
				sspan.RecordError(err)
				sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				sspan.SetStatus(trace.StatusError, err.Error())
			}
			if err != nil && st.Code() != codes.NotFound {
				log.Error(err)
				emu.Lock()
				errs = errors.Join(errs, err)
				emu.Unlock()
				return nil
			}
		}

		if res != nil && len(res.GetLocations()) > 0 {
			for _, loc := range res.GetLocations() {
				mu.Lock()
				if pos, ok := visited[loc.GetUuid()]; !ok {
					locs.Locations = append(locs.GetLocations(), loc)
					visited[loc.GetUuid()] = len(locs.Locations) - 1
				} else {
					if pos < len(locs.GetLocations()) {
						locs.GetLocations()[pos].Ips = append(locs.GetLocations()[pos].Ips, loc.GetIps()...)
						if s := locs.GetLocations()[pos].Name; len(s) == 0 {
							locs.GetLocations()[pos].Name = loc.GetName()
						} else {
							// strings.Join is used because '+=' causes performance degradation when the number of characters is large.
							locs.GetLocations()[pos].Name = strings.Join([]string{
								s, loc.GetName(),
							}, ",")
						}
					}
				}
				mu.Unlock()
			}
		}
		return nil
	})
	if errs != nil {
		err = errors.Join(err, errs)
	}
	if err != nil {
		st, _ := status.FromError(err)
		log.Error(err)
		errhandler.RecordSpanStatus(span, st, err)
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
