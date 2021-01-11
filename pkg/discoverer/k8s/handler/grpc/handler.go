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
	"strings"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/singleflight"
	"github.com/vdaas/vald/pkg/discoverer/k8s/service"
)

type DiscovererServer interface {
	discoverer.DiscovererServer
	Start(context.Context)
}

type server struct {
	dsc   service.Discoverer
	group singleflight.Group
}

const (
	podPrefix  = "pods"
	nodePrefix = "nodes"
	keyDelim   = "-"
)

func New(opts ...Option) (ds DiscovererServer, err error) {
	s := new(server)

	for _, opt := range append(defaultOptions, opts...) {
		err = opt(s)
		if err != nil {
			return nil, err
		}
	}

	s.group = singleflight.New()

	return s, nil
}

func (s *server) Start(ctx context.Context) {
}

func (s *server) Pods(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Pods, error) {
	ctx, span := trace.StartSpan(ctx, "vald/discoverer-k8s.Pods")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, _, err := s.group.Do(ctx, singleflightKey(podPrefix, req), func() (interface{}, error) {
		return s.dsc.GetPods(req)
	})
	if err != nil {
		log.Error("an error occurred during GetPods", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Pods API request %#v pods not found", req), err, info.Get())
	}
	return proto.Clone(res.(*payload.Info_Pods)).(*payload.Info_Pods), nil
}

func (s *server) Nodes(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	ctx, span := trace.StartSpan(ctx, "vald/discoverer-k8s.Nodes")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, _, err := s.group.Do(ctx, singleflightKey(nodePrefix, req), func() (interface{}, error) {
		return s.dsc.GetNodes(req)
	})
	if err != nil {
		log.Error("an error occurred during GetNodes", err)
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Nodes API request %#v nodes not found", req), err, info.Get())
	}
	return proto.Clone(res.(*payload.Info_Nodes)).(*payload.Info_Nodes), nil
}

func singleflightKey(pref string, req *payload.Discoverer_Request) string {
	return strings.Join([]string{
		pref,
		req.GetNode(),
		req.GetNamespace(),
		req.GetName(),
	}, keyDelim)
}
