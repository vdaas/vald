//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/singleflight"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/pkg/discoverer/k8s/service"
)

type DiscovererServer interface {
	discoverer.DiscovererServer
	Start(context.Context)
}

type server struct {
	dsc   service.Discoverer
	group singleflight.Group
	ip    string
	name  string
	discoverer.UnimplementedDiscovererServer
}

const (
	apiName    = "vald/discoverer/k8s"
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

func (*server) Start(context.Context) {
}

func (s *server) Pods(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Pods, error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Pods")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	key := singleflightKey(podPrefix, req)
	res, _, err := s.group.Do(ctx, key, func() (interface{}, error) {
		return s.dsc.GetPods(req)
	})
	if err != nil {
		err = status.WrapWithInternal(fmt.Sprintf("Pods API request (name: %s, namespace: %s, node: %s) failed", req.GetName(), req.GetNamespace(), req.GetNode()), err,
			&errdetails.RequestInfo{
				RequestId:   key,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/discoverer.v1.Pods",
				ResourceName: fmt.Sprintf("%s(%s)", s.name, s.ip),
			},
			info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Warnf("GetPods returned error: %v", err)
		return nil, err
	}
	if res == nil {
		err = status.WrapWithNotFound(fmt.Sprintf("Pods API request (name: %s, namespace: %s, node: %s) pods not found", req.GetName(), req.GetNamespace(), req.GetNode()), err,
			&errdetails.RequestInfo{
				RequestId:   key,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/discoverer.v1.Pods",
				ResourceName: fmt.Sprintf("%s(%s)", s.name, s.ip),
			},
			info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Warnf("Pods not found: %#v, error: %v", res, err)
		return nil, err
	}
	cp := proto.Clone(res.(*payload.Info_Pods))
	if cp == nil {
		err = status.WrapWithNotFound(fmt.Sprintf("Pods API request (name: %s, namespace: %s, node: %s) pods not found, cloned response is nil", req.GetName(), req.GetNamespace(), req.GetNode()), err,
			&errdetails.RequestInfo{
				RequestId:   key,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/discoverer.v1.Pods",
				ResourceName: fmt.Sprintf("%s(%s)", s.name, s.ip),
			},
			info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Warnf("Pods not found: %#v, error: %v", res, err)
		return nil, err
	}
	in, ok := cp.(*payload.Info_Pods)
	if in == nil || !ok {
		err = status.WrapWithNotFound(fmt.Sprintf("Pods API request (name: %s, namespace: %s, node: %s) pods not found, cloned response is nil", req.GetName(), req.GetNamespace(), req.GetNode()), err,
			&errdetails.RequestInfo{
				RequestId:   key,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/discoverer.v1.Pods",
				ResourceName: fmt.Sprintf("%s(%s)", s.name, s.ip),
			},
			info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Warnf("Pods not found: %#v, error: %v", res, err)
		return nil, err
	}
	return in, nil
}

func (s *server) Nodes(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	ctx, span := trace.StartSpan(ctx, apiName+".Nodes")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	key := singleflightKey(nodePrefix, req)
	res, _, err := s.group.Do(ctx, key, func() (interface{}, error) {
		return s.dsc.GetNodes(req)
	})
	if err != nil {
		err = status.WrapWithInternal(fmt.Sprintf("Nodes API request (name: %s, namespace: %s, node: %s) failed", req.GetName(), req.GetNamespace(), req.GetNode()), err,
			&errdetails.RequestInfo{
				RequestId:   key,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/discoverer.v1.Nodes",
				ResourceName: fmt.Sprintf("%s(%s)", s.name, s.ip),
			},
			info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeInternal(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Warnf("GetNodes returned error: %v", err)
		return nil, err
	}
	if res == nil {
		err = status.WrapWithNotFound(fmt.Sprintf("Nodes API request (name: %s, namespace: %s, node: %s) nodes not found", req.GetName(), req.GetNamespace(), req.GetNode()), err,
			&errdetails.RequestInfo{
				RequestId:   key,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/discoverer.v1.Nodes",
				ResourceName: fmt.Sprintf("%s(%s)", s.name, s.ip),
			},
			info.Get())
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Warnf("Nodes not found: %#v, error: %v", res, err)
		return nil, err
	}
	cp := proto.Clone(res.(*payload.Info_Nodes))
	if cp == nil {
		err = status.WrapWithNotFound(
			fmt.Sprintf("Nodes API request (name: %s, namespace: %s, node: %s) nodes not found, cloned response is nil", req.GetName(), req.GetNamespace(), req.GetNode()),
			err,
			&errdetails.RequestInfo{
				RequestId:   key,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/discoverer.v1.Nodes",
				ResourceName: fmt.Sprintf("%s(%s)", s.name, s.ip),
			},
			info.Get(),
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Warnf("Nodes not found: %#v, error: %v", res, err)
		return nil, err
	}
	in, ok := cp.(*payload.Info_Nodes)
	if in == nil || !ok {
		err = status.WrapWithNotFound(
			fmt.Sprintf("Nodes API request (name: %s, namespace: %s, node: %s) nodes not found, cloned response is nil", req.GetName(), req.GetNamespace(), req.GetNode()),
			err,
			&errdetails.RequestInfo{
				RequestId:   key,
				ServingData: errdetails.Serialize(req),
			},
			&errdetails.ResourceInfo{
				ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/discoverer.v1.Nodes",
				ResourceName: fmt.Sprintf("%s(%s)", s.name, s.ip),
			},
			info.Get(),
		)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(trace.StatusCodeNotFound(err.Error())...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		log.Warnf("Nodes not found: %#v, error: %v", res, err)
		return nil, err
	}
	return in, nil
}

func singleflightKey(pref string, req *payload.Discoverer_Request) string {
	return strings.Join([]string{
		pref,
		req.GetNode(),
		req.GetNamespace(),
		req.GetName(),
	}, keyDelim)
}
