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
	"strings"

	"github.com/vdaas/vald/apis/grpc/discoverer"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/pkg/discoverer/k8s/service"
)

type DiscovererServer interface {
	discoverer.DiscovererServer
	Start(context.Context)
}

type server struct {
	dsc                 service.Discoverer
	cache               cache.Cache
	enableCache         bool
	expireCheckDuration string
	expireDuration      string
}

const (
	podPrefix  = "pods"
	nodePrefix = "nodes"
	keyDelim   = "-"
)

func New(opts ...Option) (ds DiscovererServer, err error) {
	s := new(server)

	for _, opt := range append(defaultOpts, opts...) {
		err = opt(s)
		if err != nil {
			return nil, err
		}
	}

	if s.enableCache {
		if s.cache == nil {
			c, err := cache.New(
				cache.WithExpireDuration(s.expireDuration),
				cache.WithExpireCheckDuration(s.expireCheckDuration),
				cache.WithExpiredHook(func(ctx context.Context, key string) {
					pref, req := toRequest(key)
					switch pref {
					case podPrefix:
						pods, err := s.dsc.GetPods(req)
						if err != nil {
							log.Error(err)
							return
						}
						s.cache.Set(key, pods)
					case nodePrefix:
						nodes, err := s.dsc.GetNodes(req)
						if err != nil {
							log.Error(err)
							return
						}
						s.cache.Set(key, nodes)
					}
				}),
			)
			if err != nil {
				return nil, err
			}
			s.cache = c
		}
	}

	return s, nil
}

func (s *server) Start(ctx context.Context) {
	if s.enableCache && s.cache != nil {
		s.cache.Start(ctx)
	}
}

func (s *server) Pods(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Pods, error) {
	ctx, span := trace.StartSpan(ctx, "vald/discoverer-k8s.Pods")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var cacheKey string
	if s.cache != nil {
		cacheKey = toCacheKey(podPrefix, req)
		cp, ok := s.cache.Get(cacheKey)
		if ok && cp != nil {
			return proto.Clone(cp.(*payload.Info_Pods)).(*payload.Info_Pods), nil
		}
	}
	pods, err := s.dsc.GetPods(req)
	if err != nil {
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Pods API request %#v pods not found", req), err, info.Get())
	}
	if s.cache != nil && len(cacheKey) != 0 {
		s.cache.Set(cacheKey, pods)
	}
	return proto.Clone(pods).(*payload.Info_Pods), nil
}

func (s *server) Nodes(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	ctx, span := trace.StartSpan(ctx, "vald/discoverer-k8s.Nodes")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var cacheKey string
	if s.cache != nil {
		cacheKey = toCacheKey(nodePrefix, req)
		cn, ok := s.cache.Get(cacheKey)
		if ok && cn != nil {
			return proto.Clone(cn.(*payload.Info_Nodes)).(*payload.Info_Nodes), nil
		}
	}
	nodes, err := s.dsc.GetNodes(req)
	if err != nil {
		log.Error(err)
		if span != nil {
			span.SetStatus(trace.StatusCodeNotFound(err.Error()))
		}
		return nil, status.WrapWithNotFound(fmt.Sprintf("Nodes API request %#v nodes not found", req), err, info.Get())
	}
	if s.cache != nil && len(cacheKey) != 0 {
		s.cache.Set(cacheKey, nodes)
	}
	return proto.Clone(nodes).(*payload.Info_Nodes), nil
}

func toRequest(key string) (pref string, req *payload.Discoverer_Request) {
	infos := strings.Split(key, keyDelim)
	return infos[0], &payload.Discoverer_Request{
		Node:      infos[1],
		Namespace: infos[2],
		Name:      infos[3],
	}
}

func toCacheKey(pref string, req *payload.Discoverer_Request) string {
	return strings.Join([]string{
		pref,
		req.GetNode(),
		req.GetNamespace(),
		req.GetName(),
	}, keyDelim)
}
