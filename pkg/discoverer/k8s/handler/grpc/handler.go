//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

	"github.com/vdaas/vald/apis/grpc/discoverer"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/pkg/discoverer/k8s/service"
)

type server struct {
	dsc service.Discoverer
}

func New(opts ...Option) discoverer.DiscovererServer {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		opt(s)
	}
	return s
}

func (s *server) Pods(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Pods, error) {
	ctx, span := trace.StartSpan(ctx, "vald/discoverer-k8s.Pods")
	defer span.End()
	pods, err := s.dsc.GetPods(req)
	if err != nil {
		return nil, status.WrapWithNotFound(fmt.Sprintf("Pods API request %#v pods not found", req), err, info.Get())
	}
	return pods, nil
}

func (s *server) Nodes(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	ctx, span := trace.StartSpan(ctx, "vald/discoverer-k8s.Nodes")
	defer span.End()
	nodes, err := s.dsc.GetNodes(req)
	if err != nil {
		return nil, status.WrapWithNotFound(fmt.Sprintf("Nodes API request %#v nodes not found", req), err, info.Get())
	}
	return nodes, nil
}
