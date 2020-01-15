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

// Package service manages the main logic of server.
package service

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/k8s"
	mpod "github.com/vdaas/vald/internal/k8s/metrics/pod"
	"github.com/vdaas/vald/internal/k8s/node"
	"github.com/vdaas/vald/internal/k8s/pod"
	"github.com/vdaas/vald/internal/log"
)

type Discoverer interface {
	GetServers(string, string) *payload.Info_Servers
	Start(context.Context) <-chan error
}

type discoverer struct {
	maxServers int
	nodes      sync.Map
	pods       sync.Map
	podMetrics sync.Map
	ctrl       k8s.Controller
}

func New() (dsc Discoverer, err error) {
	d := new(discoverer)
	d.ctrl, err = k8s.New(
		k8s.WithControllerName("vald k8s agent discoverer"),
		k8s.WithDisableLeaderElection(),
		k8s.WithResourceController(mpod.New(
			mpod.WithControllerName("pod metrics discoverer"),
			mpod.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			mpod.WithOnReconcileFunc(func(podMetricsList map[string][]mpod.Pod) {
				b, _ := json.Marshal(podMetricsList)
				log.Debug(string(b))
				for name, pods := range podMetricsList {
					if len(pods) > d.maxServers {
						d.maxServers = len(pods)
					}
					d.podMetrics.Store(name, pods)
				}
			}),
		)),
		k8s.WithResourceController(pod.New(
			pod.WithControllerName("pod discoverer"),
			pod.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			pod.WithOnReconcileFunc(func(podList map[string][]pod.Pod) {
				b, _ := json.Marshal(podList)
				log.Debug(string(b))
				for name, pods := range podList {
					if len(pods) > d.maxServers {
						d.maxServers = len(pods)
					}
					d.pods.Store(name, pods)
				}
			}),
		)),
		k8s.WithResourceController(node.New(
			node.WithControllerName("node discoverer"),
			node.WithOnErrorFunc(func(err error) {
				log.Error(err)
			}),
			node.WithOnReconcileFunc(func(nodes []node.Node) {
				log.Debug(nodes)
				for _, n := range nodes {
					d.nodes.Store(n.Name, n)
				}
			}),
		)),
	)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *discoverer) Start(ctx context.Context) <-chan error {
	return d.ctrl.Start(ctx)
}

func (d *discoverer) GetServers(name, nodeName string) (srvs *payload.Info_Servers) {
	srvs = &payload.Info_Servers{
		Servers: make([]*payload.Info_Server, 0, d.maxServers),
	}
	d.pods.Range(func(metaname, pods interface{}) bool {
		for _, p := range pods.([]pod.Pod) {
			if (name == "" && nodeName == "") ||
				(name == "" && nodeName == p.NodeName) ||
				(metaname == name && nodeName == "") ||
				(metaname == name && nodeName == p.NodeName) {
				srv := &payload.Info_Server{
					Name: p.Name,
					Ip:   p.IP,
					Cpu:  p.CPU,
					Mem:  p.Mem,
				}
				if nr, ok := d.nodes.Load(p.NodeName); ok {
					n, ok := nr.(node.Node)
					if ok {
						srv.Server.Name = n.Name
						srv.Server.Ip = n.IP
						srv.Server.Cpu = n.CPU
						srv.Server.Mem = n.Mem
					}
				}
				srvs.Servers = append(srvs.Servers, srv)
			}
		}
		return true
	})
	srvs.Servers = srvs.Servers[:len(srvs.Servers)]
	return srvs
}
