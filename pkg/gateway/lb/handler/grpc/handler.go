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
	"maps"
	"slices"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	lbmetrics "github.com/vdaas/vald/internal/observability/metrics/gateway/lb"
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
	replica           int
	streamConcurrency int
	multiConcurrency  int
}

const (
	apiName = "vald/gateway/lb"

	// uuidField is the errdetails.BadRequest field name of a request's uuid.
	uuidField = "uuid"
)

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

// joinTargetErrors joins the per-agent failures sorted by target, so the leading
// error (whose code a joined error reports) does not depend on map order.
func joinTargetErrors(errs map[string]error) error {
	if len(errs) == 0 {
		return nil
	}
	list := make([]error, 0, len(errs))
	for _, target := range slices.Sorted(maps.Keys(errs)) {
		list = append(list, errs[target])
	}
	return errors.Join(list...)
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
	if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) || ctx.Err() != nil {
		return err
	}
	ferrs := joinTargetErrors(ferr)
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
