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
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net"
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
)

func (s *server) Insert(
	ctx context.Context, req *payload.Insert_Request,
) (ce *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.InsertRPCName), apiName+"/"+vald.InsertRPCName)
	defer trace.End(span)
	uuid := req.GetVector().GetId()
	if len(uuid) == 0 {
		err = errors.ErrInvalidMetaDataConfig
		err = status.WrapWithInvalidArgument(vald.InsertRPCName+" API invalid uuid", err, append(s.errDetails(ctx, vald.InsertRPCName, uuid, req),
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "invalid id",
						Description: err.Error(),
					},
				},
			})...)
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, err)
	}
	vec := req.GetVector().GetVector()
	vl := len(vec)
	if vl < algorithm.MinimumVectorDimensionSize {
		err = errors.ErrInvalidDimensionSize(vl, 0)
		err = status.WrapWithInvalidArgument(vald.InsertRPCName+" API invalid vector argument", err, append(s.errDetails(ctx, vald.InsertRPCName, uuid, req),
			&errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequestFieldViolation{
					{
						Field:       "vector dimension size",
						Description: err.Error(),
					},
				},
			}, info.Get())...)
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, err)
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		id, err := s.exists(ctx, uuid)
		var attrs trace.Attributes
		if err != nil {
			switch {
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(vald.ExistsRPCName+" API for "+vald.InsertRPCName+" API connection not found", err, s.errDetails(ctx, vald.InsertRPCName, uuid, req)...)
				attrs = trace.StatusCodeInternal(err.Error())
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(vald.ExistsRPCName+" API for "+vald.InsertRPCName+" API canceled", err, s.errDetails(ctx, vald.InsertRPCName, uuid, req)...)
				attrs = trace.StatusCodeCancelled(err.Error())
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithDeadlineExceeded(vald.ExistsRPCName+" API for "+vald.InsertRPCName+" API deadline exceeded", err, s.errDetails(ctx, vald.InsertRPCName, uuid, req)...)
				attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			default:
				err = nil
			}
		} else if id != nil && len(id.GetId()) != 0 {
			err = status.WrapWithAlreadyExists(vald.InsertRPCName+" API uuid "+uuid+"'s data already exists", errors.ErrMetaDataAlreadyExists(uuid),
				append(s.errDetails(ctx, vald.InsertRPCName, uuid, req), info.Get())...)
			attrs = trace.StatusCodeAlreadyExists(err.Error())
		}
		if err != nil {
			errhandler.RecordSpanAttrs(span, attrs, err)
			return nil, err
		}
	}

	mu := new(sync.Mutex)
	ce = &payload.Object_Location{
		Uuid: uuid,
		Ips:  make([]string, 0, s.replica),
	}
	locs := make([]string, 0, s.replica)
	// The caller's request is read by others while this runs (MultiInsert serializes the
	// whole MultiRequest while sibling Inserts are in flight), so never write to it:
	// forward a new request that skips the agents' exist check (done above) and is timestamped.
	freq := req
	if cfg := req.GetConfig(); cfg == nil || !cfg.GetSkipStrictExistCheck() || cfg.GetTimestamp() == 0 {
		ts := cfg.GetTimestamp()
		if ts == 0 {
			ts = time.Now().UnixNano()
		}
		freq = &payload.Insert_Request{
			Vector: req.GetVector(),
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Filters:              cfg.GetFilters(),
				Timestamp:            ts,
			},
		}
	}
	emu := new(sync.Mutex)
	// The client's backoff and DoMulti's second pass call f again for the same target,
	// so keep only each target's latest failure: a later success clears it.
	errByTarget := make(map[string]error, s.replica)
	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "DoMulti/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
		defer trace.End(span)
		loc, err := vc.Insert(ctx, freq, copts...)
		if err != nil {
			// A canceled or timed-out agent stored nothing: return its error so that
			// DoMulti does not count it as a replica. It does not abort the request.
			switch {
			case errors.Is(err, context.Canceled),
				errors.Is(err, errors.ErrRPCCallFailed(target, context.Canceled)):
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.InsertRPCName + ".DoMulti/" +
							target + " canceled: " + err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				err = status.WrapWithCanceled(vald.InsertRPCName+" API canceled for "+target, err)
				emu.Lock()
				errByTarget[target] = err
				emu.Unlock()
				return err
			case errors.Is(err, context.DeadlineExceeded),
				errors.Is(err, errors.ErrRPCCallFailed(target, context.DeadlineExceeded)):
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.InsertRPCName + ".DoMulti/" +
							target + " deadline_exceeded: " + err.Error())...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				err = status.WrapWithDeadlineExceeded(vald.InsertRPCName+" API deadline exceeded for "+target, err)
				emu.Lock()
				errByTarget[target] = err
				emu.Unlock()
				return err
			}
			st, msg, err := status.ParseError(err, codes.Internal,
				"failed to parse "+vald.InsertRPCName+" gRPC error response",
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(freq),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %s", apiName, s.name, s.ip, target),
				})
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			if err != nil && st.Code() != codes.AlreadyExists {
				emu.Lock()
				errByTarget[target] = err
				emu.Unlock()
				return err
			}
			// The agent already holds uuid: report it in the location as Update does.
			host, _, herr := net.SplitHostPort(target)
			if herr != nil {
				host = target
			}
			emu.Lock()
			delete(errByTarget, target)
			emu.Unlock()
			mu.Lock()
			ce.Ips = append(ce.GetIps(), host)
			locs = append(locs, host)
			mu.Unlock()
			return nil
		}
		emu.Lock()
		delete(errByTarget, target)
		emu.Unlock()
		mu.Lock()
		ce.Ips = append(ce.GetIps(), loc.GetIps()...)
		locs = append(locs, loc.GetName())
		mu.Unlock()
		return nil
	})
	errs := joinTargetErrors(errByTarget)
	if err != nil {
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			err = status.WrapWithInternal(vald.InsertRPCName+" API connection not found", err,
				&errdetails.RequestInfo{
					RequestId:   uuid,
					ServingData: errdetails.Serialize(freq),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.InsertRPCName + ".DoMulti",
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				})
			return errhandler.HandleError[payload.Object_Location](span, codes.Internal, err)
		}
		if errs == nil {
			errs = err
		} else {
			errs = errors.Join(errs, err)
		}
		st, _ := status.FromError(errs)
		errhandler.RecordSpanStatus(span, st, errs)
		return nil, errs
	}
	// The request succeeds once any agent stored the vector; the failed agents only
	// leave it under-replicated. With no stored location it fails even when no agent
	// reported an error: DoMulti skips agents whose circuit breaker is open without
	// calling f, and does not report that.
	if len(locs) == 0 {
		if errs != nil {
			st, _ := status.FromError(errs)
			errhandler.RecordSpanStatus(span, st, errs)
			return nil, errs
		}
		err = status.WrapWithUnavailable(vald.InsertRPCName+" API stored the vector on no agent",
			errors.ErrNoAvailableAgentToInsert, s.errDetails(ctx, vald.InsertRPCName, uuid, req)...)
		return errhandler.HandleError[payload.Object_Location](span, codes.Unavailable, err)
	}
	if errs != nil || (len(locs) < s.replica && len(locs) < s.gateway.GetAgentCount(ctx)) {
		s.warnPartialWrite(vald.InsertRPCName, uuid, len(locs), errByTarget, errs)
	}
	slices.Sort(locs)
	ce.Name = strings.Join(locs, ",")
	return ce, nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) (err error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(stream.Context(), vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.StreamInsertRPCName), apiName+"/"+vald.StreamInsertRPCName)
	defer trace.End(span)
	err = grpc.BidirectionalStream(ctx, stream, s.streamConcurrency,
		func(ctx context.Context, req *payload.Insert_Request) (*payload.Object_StreamLocation, error) {
			ctx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "BidirectionalStream"), apiName+"/"+vald.StreamInsertRPCName+"/id-"+req.GetVector().GetId())
			defer trace.End(sspan)
			res, err := s.Insert(ctx, req)
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

func (s *server) MultiInsert(
	ctx context.Context, reqs *payload.Insert_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.MultiInsertRPCName), apiName+"/"+vald.MultiInsertRPCName)
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
		if r != nil && r.GetVector() != nil && len(r.GetVector().GetVector()) >= algorithm.MinimumVectorDimensionSize && r.GetVector().GetId() != "" {
			idx := i
			req := r
			eg.Go(safety.RecoverFunc(func() (err error) {
				ectx, sspan := trace.StartSpan(grpc.WrapGRPCMethod(ectx, "eg.Go"), apiName+"/"+vald.MultiInsertRPCName+"/id-"+req.GetVector().GetId())
				defer trace.End(sspan)
				res, err := s.Insert(ectx, req)
				if err != nil {
					st, _ := status.FromError(err)
					errhandler.RecordSpanStatus(sspan, st, err)
					emu.Lock()
					errList = append(errList, err)
					emu.Unlock()
				} else if res != nil && res.GetUuid() == req.GetVector().GetId() && res.GetIps() != nil {
					lmu.Lock()
					locs.Locations[idx] = res
					lmu.Unlock()
				}
				return nil
			}))
		} else {
			var (
				err   error
				field string
			)
			switch {
			case r.GetVector() == nil, len(r.GetVector().GetVector()) < algorithm.MinimumVectorDimensionSize:
				err = errors.ErrInvalidDimensionSize(len(r.GetVector().GetVector()), 0)
				field = vectorField
			case r.GetVector().GetId() == "":
				err = errors.ErrInvalidUUID(r.GetVector().GetId())
				field = "uuid"
			}
			if !detailsBuilt {
				servingData = errdetails.Serialize(reqs)
				resourceName = fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx))
				detailsBuilt = true
			}
			err = status.WrapWithInvalidArgument(vald.MultiInsertRPCName+" API invalid vector argument", err,
				&errdetails.RequestInfo{
					RequestId:   r.GetVector().GetId(),
					ServingData: servingData,
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.MultiInsertRPCName,
					ResourceName: resourceName,
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       field,
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
