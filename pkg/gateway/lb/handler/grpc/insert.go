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
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

func (s *server) Insert(
	ctx context.Context, req *payload.Insert_Request,
) (ce *payload.Object_Location, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.InsertRPCServiceName+"/"+vald.InsertRPCName), apiName+"/"+vald.InsertRPCName)
		defer trace.End(span)
	}
	uuid := req.GetVector().GetId()
	if len(uuid) == 0 {
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, s.invalidUUIDError(ctx, vald.InsertRPCName, vald.InsertRPCName, uuid, req))
	}
	vec := req.GetVector().GetVector()
	vl := len(vec)
	if vl < algorithm.MinimumVectorDimensionSize {
		return errhandler.HandleError[payload.Object_Location](span, codes.InvalidArgument, s.invalidVectorDimError(ctx, vald.InsertRPCName, vald.InsertRPCName, uuid, req, vl))
	}
	if !req.GetConfig().GetSkipStrictExistCheck() {
		ok, cerr := s.existsCheck(ctx, uuid)
		if cerr != nil {
			if cerr = s.handleExistsCheckError(ctx, span, vald.InsertRPCName, uuid, req, cerr); cerr != nil {
				return nil, cerr
			}
		} else if ok {
			cerr = status.WrapWithAlreadyExists(vald.InsertRPCName+" API uuid "+uuid+"'s data already exists", errors.ErrMetaDataAlreadyExists(uuid),
				append(s.errDetails(ctx, vald.InsertRPCName, uuid, req), info.Get())...)
			if span != nil {
				errhandler.RecordSpanAttrs(span, trace.StatusCodeAlreadyExists(cerr.Error()), cerr)
			}
			return nil, cerr
		}
	}

	type insertState struct {
		ce          *payload.Object_Location
		errByTarget map[string]error
		locs        []string
		mu          sync.Mutex
	}
	ins := &insertState{
		ce: &payload.Object_Location{
			Uuid: uuid,
			Ips:  make([]string, 0, s.replica),
		},
		locs: make([]string, 0, s.replica),
	}
	ce = ins.ce
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
	// The client's backoff and DoMulti's second pass call f again for the same target,
	// so keep only each target's latest failure: a later success clears it.
	err = s.gateway.DoMulti(ctx, s.replica, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) (err error) {
		var sspan trace.Span
		if trace.Enabled() {
			ctx, sspan = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "DoMulti/"+target), apiName+"/"+vald.InsertRPCName+"/"+target)
			defer trace.End(sspan)
		}
		loc, err := vc.Insert(ctx, freq, copts...)
		if err != nil {
			// A canceled or timed-out agent stored nothing: return its error so that
			// DoMulti does not count it as a replica. It does not abort the request.
			switch {
			case errors.Is(err, context.Canceled):
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.StatusCodeCancelled(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.InsertRPCName + ".DoMulti/" +
							target + " canceled: " + err.Error())...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				err = status.WrapWithCanceled(vald.InsertRPCName+" API canceled for "+target, err)
				ins.mu.Lock()
				if ins.errByTarget == nil {
					ins.errByTarget = make(map[string]error, s.replica)
				}
				ins.errByTarget[target] = err
				ins.mu.Unlock()
				return err
			case errors.Is(err, context.DeadlineExceeded):
				if sspan != nil {
					sspan.RecordError(err)
					sspan.SetAttributes(trace.StatusCodeDeadlineExceeded(
						errdetails.ValdGRPCResourceTypePrefix +
							"/vald.v1." + vald.InsertRPCName + ".DoMulti/" +
							target + " deadline_exceeded: " + err.Error())...)
					sspan.SetStatus(trace.StatusError, err.Error())
				}
				err = status.WrapWithDeadlineExceeded(vald.InsertRPCName+" API deadline exceeded for "+target, err)
				ins.mu.Lock()
				if ins.errByTarget == nil {
					ins.errByTarget = make(map[string]error, s.replica)
				}
				ins.errByTarget[target] = err
				ins.mu.Unlock()
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
			if sspan != nil {
				errhandler.RecordSpanAttrs(sspan, trace.FromGRPCStatus(st.Code(), msg), err)
			}
			if err != nil && st.Code() != codes.AlreadyExists {
				ins.mu.Lock()
				if ins.errByTarget == nil {
					ins.errByTarget = make(map[string]error, s.replica)
				}
				ins.errByTarget[target] = err
				ins.mu.Unlock()
				return err
			}
			// The agent already holds uuid: report it in the location as Update does.
			host, _, herr := net.SplitHostPort(target)
			if herr != nil {
				host = target
			}
			ins.mu.Lock()
			if ins.errByTarget != nil {
				delete(ins.errByTarget, target)
			}
			ins.ce.Ips = append(ins.ce.GetIps(), host)
			ins.locs = append(ins.locs, host)
			ins.mu.Unlock()
			return nil
		}
		ins.mu.Lock()
		if ins.errByTarget != nil {
			delete(ins.errByTarget, target)
		}
		ins.ce.Ips = append(ins.ce.GetIps(), loc.GetIps()...)
		ins.locs = append(ins.locs, loc.GetName())
		ins.mu.Unlock()
		return nil
	})
	errs := errors.JoinMap(ins.errByTarget)
	if err != nil {
		if errors.Is(err, errGRPCClientConnNotFoundStar) {
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
		errhandler.RecordSpanFromError(span, errs)
		return nil, errs
	}
	// The request succeeds once any agent stored the vector; the failed agents only
	// leave it under-replicated. With no stored location it fails even when no agent
	// reported an error: DoMulti skips agents whose circuit breaker is open without
	// calling f, and does not report that.
	if len(ins.locs) == 0 {
		if errs != nil {
			errhandler.RecordSpanFromError(span, errs)
			return nil, errs
		}
		err = status.WrapWithUnavailable(vald.InsertRPCName+" API stored the vector on no agent",
			errors.ErrNoAvailableAgentToInsert, s.errDetails(ctx, vald.InsertRPCName, uuid, req)...)
		return errhandler.HandleError[payload.Object_Location](span, codes.Unavailable, err)
	}
	if errs != nil || (len(ins.locs) < s.replica && len(ins.locs) < s.gateway.GetAgentCount(ctx)) {
		s.warnPartialWrite(vald.InsertRPCName, uuid, len(ins.locs), ins.errByTarget, errs)
	}
	slices.Sort(ins.locs)
	ce.Name = strings.Join(ins.locs, ",")
	return ce, nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) (err error) {
	return streamLocation(s, stream, vald.InsertRPCServiceName, vald.StreamInsertRPCName, func(req *payload.Insert_Request) string {
		return req.GetVector().GetId()
	}, s.Insert)
}

func (s *server) MultiInsert(
	ctx context.Context, reqs *payload.Insert_MultiRequest,
) (locs *payload.Object_Locations, errs error) {
	return multiWrite(ctx, s, vald.InsertRPCServiceName, vald.MultiInsertRPCName, reqs, reqs.GetRequests(),
		func(r *payload.Insert_Request) (string, string, error) {
			return validateMultiVector(r.GetVector())
		}, s.Insert)
}
