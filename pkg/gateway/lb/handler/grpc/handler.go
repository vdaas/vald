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
	"runtime"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/status"
	lbmetrics "github.com/vdaas/vald/internal/observability/metrics/gateway/lb"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

type server struct {
	vald.UnimplementedValdServer
	eg                errgroup.Group
	gateway           service.Gateway
	searchMetrics     *lbmetrics.SearchMetrics
	name              string
	ip                string
	timeout           time.Duration
	searchQuorumGrace time.Duration
	replica           int
	streamConcurrency int
	multiConcurrency  int
}

const (
	apiName = "vald/gateway/lb"

	// uuidField is the errdetails.BadRequest field name of a request's uuid.
	uuidField = "uuid"

	// vectorField is the errdetails.BadRequest field name of a request's vector.
	vectorField = "vector"
)

// autoMultiPerAgent is the number of in-flight requests per agent the automatic
// Multi* concurrency targets. BenchmarkAgentConcurrencySweep (agent handler
// package) shows Insert throughput levelling off around 64 in flight whatever
// the agent's cores, and Search levelling off near 2x the agent's cores (64 for
// a 32-core agent) and staying flat beyond.
const autoMultiPerAgent = 64

// autoMultiPerCPU scales the automatic Multi* concurrency ceiling by GOMAXPROCS,
// the limit the handler used as its fixed default before the automatic mode.
const autoMultiPerCPU = 10

// multiLimit returns how many items of a Multi* request run at once. An explicit
// multiConcurrency wins; 0 selects the automatic limit: GOMAXPROCS*10, capped so
// each agent sees about autoMultiPerAgent in-flight requests. Write items reach
// index_replica agents each, so the cap is ceil(agents*autoMultiPerAgent/replica);
// a MultiSearch item reaches every agent, so the cap is autoMultiPerAgent.
func (s *server) multiLimit(ctx context.Context, fanOutAll bool) int {
	if s.multiConcurrency > 0 {
		return s.multiConcurrency
	}
	limit := runtime.GOMAXPROCS(0) * autoMultiPerCPU
	perAgent := autoMultiPerAgent
	if !fanOutAll && s.replica > 0 {
		if agents := s.gateway.GetAgentCount(ctx); agents > 0 {
			perAgent = (agents*autoMultiPerAgent + s.replica - 1) / s.replica
		}
	}
	return max(min(limit, perAgent), 1)
}

func New(opts ...Option) vald.Server {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

// errDetails builds the RequestInfo/ResourceInfo details attached to an error status.
// It JSON-encodes the whole request and queries the gateway addresses, so call it only
// on error paths; rpc names the resource, e.g. vald.InsertRPCName.
func (s *server) errDetails(ctx context.Context, rpc, uuid string, req any) []any {
	return []any{
		&errdetails.RequestInfo{
			RequestId:   uuid,
			ServingData: errdetails.Serialize(req),
		},
		&errdetails.ResourceInfo{
			ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + rpc,
			ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
		},
	}
}

// warnPartialWrite reports a write RPC that succeeded although it was applied to
// fewer agents than intended (failed agents, or fewer agents than replicas). It is
// the single place such partial writes are reported and counted. failed holds each
// failed agent's latest error by target, so every failed agent is counted once.
// Callers report only writes applied to at least one agent (succeeded > 0): a
// write applied to no agent is a failure, not a partial write.
func (s *server) warnPartialWrite(
	rpc, uuid string, succeeded int, failed map[string]error, err error,
) {
	s.searchMetrics.RecordWritePartialSuccess(rpc)
	for _, ferr := range failed {
		code := codes.Unknown
		if st, ok := status.FromError(ferr); ok && st != nil {
			code = st.Code()
		}
		s.searchMetrics.RecordWriteAgentFailure(rpc, code)
	}
	log.Warnf("%s API for uuid %s succeeded partially: applied to %d agent(s), failed on %d agent(s), errors: %v",
		rpc, uuid, succeeded, len(failed), err)
}

// partialBroadCastResult decides the result of a write BroadCast whose callback
// recorded each agent's latest failure in ferr and that succeeded on succeeded
// agents. A missing connection and a done request context stay fatal; otherwise
// the agent failures fail the request only when no agent applied the change, and
// are reported as a partial write when some did.
func (s *server) partialBroadCastResult(
	ctx context.Context, rpc, uuid string, err error, succeeded int, ferr map[string]error,
) error {
	if errors.Is(err, errGRPCClientConnNotFoundStar) || ctx.Err() != nil {
		return err
	}
	ferrs := errors.JoinMap(ferr)
	if succeeded == 0 {
		if err != nil {
			return err
		}
		return ferrs
	}
	if ferrs == nil {
		ferrs = err
	}
	if ferrs != nil {
		s.warnPartialWrite(rpc, uuid, succeeded, ferr, ferrs)
	}
	return nil
}

var errGRPCClientConnNotFoundStar = errors.ErrGRPCClientConnNotFound("*")

func streamLocation[Q proto.Message, S grpc.TypedServerStream[Q, *payload.Object_StreamLocation]](
	s *server,
	stream S,
	svcName, rpcName string,
	itemID func(Q) string,
	fn func(context.Context, Q) (*payload.Object_Location, error),
) error {
	return grpc.HandleBidirectionalStream(
		stream,
		s.streamConcurrency,
		vald.PackageName+"."+svcName+"/"+rpcName,
		apiName+"/"+rpcName,
		func(req Q) string { return "id-" + itemID(req) },
		fn,
		func(res *payload.Object_Location) *payload.Object_StreamLocation {
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Location{Location: res},
			}
		},
		func(st *status.Status) *payload.Object_StreamLocation {
			return &payload.Object_StreamLocation{
				Payload: &payload.Object_StreamLocation_Status{Status: st.Proto()},
			}
		},
	)
}

func validateMultiVector(vec *payload.Object_Vector) (uuid, field string, err error) {
	if vec == nil || len(vec.GetVector()) < algorithm.MinimumVectorDimensionSize {
		return vec.GetId(), vectorField, errors.ErrInvalidDimensionSize(len(vec.GetVector()), 0)
	}
	if vec.GetId() == "" {
		return "", uuidField, errors.ErrInvalidUUID("")
	}
	return vec.GetId(), "", nil
}

func multiWrite[Q any](
	ctx context.Context,
	s *server,
	svcName, rpcName string,
	rawReqs any,
	reqs []Q,
	validate func(Q) (uuid, field string, err error),
	fn func(context.Context, Q) (*payload.Object_Location, error),
) (locs *payload.Object_Locations, errs error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+svcName+"/"+rpcName), apiName+"/"+rpcName)
		defer trace.End(span)
	}
	if len(reqs) == 0 {
		return new(payload.Object_Locations), nil
	}
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, len(reqs)),
	}
	if len(reqs) == 1 {
		id, field, verr := validate(reqs[0])
		if verr == nil {
			reqCtx := ctx
			var sspan trace.Span
			if trace.Enabled() {
				reqCtx, sspan = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "eg.Go"), apiName+"/"+rpcName+"/id-"+id)
				defer trace.End(sspan)
			}
			res, rerr := fn(reqCtx, reqs[0])
			if rerr != nil {
				if sspan != nil {
					errhandler.RecordSpanFromError(sspan, rerr)
				}
				errs = rerr
			} else if res != nil && res.GetUuid() == id && res.GetIps() != nil {
				locs.Locations[0] = res
			}
		} else {
			errs = status.WrapWithInvalidArgument(rpcName+" API invalid vector argument", verr,
				&errdetails.RequestInfo{
					RequestId:   id,
					ServingData: errdetails.Serialize(rawReqs),
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + rpcName,
					ResourceName: fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx)),
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       field,
							Description: verr.Error(),
						},
					},
				}, info.Get())
		}
		if errs != nil && span != nil {
			errhandler.RecordSpanFromError(span, errs)
		}
		return locs, errs
	}
	var (
		errList []error
		emu     sync.Mutex
	)
	eg := errgroup.NewBoundedFanOut(ctx, s.multiLimit(ctx, false))
	// The invalid-argument details are identical for every invalid request of
	// the batch, so serialize the whole batch and resolve Addrs once, lazily.
	var (
		detailsBuilt bool
		servingData  string
		resourceName string
	)
	for i, r := range reqs {
		id, field, verr := validate(r)
		if verr == nil {
			idx := i
			req := r
			eg.Go(safety.RecoverFunc(func() error {
				reqCtx := ctx
				var sspan trace.Span
				if trace.Enabled() {
					reqCtx, sspan = trace.StartSpan(grpc.WrapGRPCMethod(ctx, "eg.Go"), apiName+"/"+rpcName+"/id-"+id)
					defer trace.End(sspan)
				}
				res, rerr := fn(reqCtx, req)
				if rerr != nil {
					if sspan != nil {
						errhandler.RecordSpanFromError(sspan, rerr)
					}
					emu.Lock()
					errList = append(errList, rerr)
					emu.Unlock()
				} else if res != nil && res.GetUuid() == id && res.GetIps() != nil {
					locs.Locations[idx] = res
				}
				return nil
			}))
		} else {
			if !detailsBuilt {
				servingData = errdetails.Serialize(rawReqs)
				resourceName = fmt.Sprintf("%s: %s(%s) to %v", apiName, s.name, s.ip, s.gateway.Addrs(ctx))
				detailsBuilt = true
			}
			ierr := status.WrapWithInvalidArgument(rpcName+" API invalid vector argument", verr,
				&errdetails.RequestInfo{
					RequestId:   id,
					ServingData: servingData,
				},
				&errdetails.ResourceInfo{
					ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + rpcName,
					ResourceName: resourceName,
				},
				&errdetails.BadRequest{
					FieldViolations: []*errdetails.BadRequestFieldViolation{
						{
							Field:       field,
							Description: verr.Error(),
						},
					},
				}, info.Get())
			emu.Lock()
			errList = append(errList, ierr)
			emu.Unlock()
		}
	}
	if werr := eg.Wait(); werr != nil {
		errList = append(errList, werr)
	}
	errs = errors.Join(errList...)

	if errs != nil && span != nil {
		errhandler.RecordSpanFromError(span, errs)
	}

	return locs, errs
}
