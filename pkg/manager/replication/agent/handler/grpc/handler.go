//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"fmt"

	"github.com/vdaas/vald/apis/grpc/v1/manager/replication/agent"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/pkg/manager/replication/agent/service"
)

type Server agent.ReplicationServer

type server struct {
	rep service.Replicator
}

func New(opts ...Option) Server {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Recover(ctx context.Context, req *payload.Replication_Recovery) (_ *payload.Empty, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-replication-agent.Recover")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if err != nil {
		if span != nil {
			span.SetStatus(trace.StatusCodeInternal(err.Error()))
		}
		return nil, status.WrapWithInternal(fmt.Sprintf("RemoveIPs API uuid %s's could not RemoveIPs", ""), err, info.Get())
	}
	return new(payload.Empty), nil
}

func (s *server) Rebalance(ctx context.Context, req *payload.Replication_Rebalance) (*payload.Empty, error) {
	ctx, span := trace.StartSpan(ctx, "vald/manager-replication-agent.Rebalance")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return new(payload.Empty), nil
}

func (s *server) AgentInfo(ctx context.Context, req *payload.Empty) (*payload.Replication_Agents, error) {
	// TODO implement this later
	ctx, span := trace.StartSpan(ctx, "vald/manager-replication-agent.AgentInfo")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return new(payload.Replication_Agents), nil
}
