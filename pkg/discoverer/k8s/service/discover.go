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
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	mnode "github.com/vdaas/vald/internal/k8s/metrics/node"
	mpod "github.com/vdaas/vald/internal/k8s/metrics/pod"
	"github.com/vdaas/vald/internal/k8s/node"
	"github.com/vdaas/vald/internal/k8s/pod"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
)

type Discoverer interface {
	Start(context.Context) (<-chan error, error)
	GetPods(*payload.Discoverer_Request) (*payload.Info_Pods, error)
	GetNodes(*payload.Discoverer_Request) (*payload.Info_Nodes, error)
}

type discoverer struct {
	maxPods         int
	nodes           nodeMap
	nodeMetrics     nodeMetricsMap
	pods            podsMap
	podMetrics      podMetricsMap
	podsByNode      atomic.Value
	podsByNamespace atomic.Value
	podsByName      atomic.Value
	nodeByName      atomic.Value
	ctrl            k8s.Controller
	namespace       string
	name            string
	csd             time.Duration
	eg              errgroup.Group
}

func New(opts ...Option) (dsc Discoverer, err error) {
	d := new(discoverer)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(d); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	d.podsByNode.Store(make(map[string]map[string]map[string][]*payload.Info_Pod))
	d.podsByNamespace.Store(make(map[string]map[string][]*payload.Info_Pod))
	d.podsByName.Store(make(map[string][]*payload.Info_Pod))
	d.nodeByName.Store(make(map[string]*payload.Info_Node))
	d.ctrl, err = k8s.New(
		k8s.WithControllerName("vald k8s agent discoverer"),
		k8s.WithDisableLeaderElection(),
		k8s.WithResourceController(mnode.New(
			mnode.WithControllerName("node metrics discoverer"),
			mnode.WithOnErrorFunc(func(err error) {
				log.Error("failed to reconcile:", err)
			}),
			mnode.WithOnReconcileFunc(func(nodes map[string]mnode.Node) {
				log.Debugf("node metrics reconciled\t%#v", nodes)
				for name, metrics := range nodes {
					d.nodeMetrics.Store(name, metrics)
				}
				d.nodeMetrics.Range(func(name string, _ mnode.Node) bool {
					_, ok := nodes[name]
					if !ok {
						d.nodeMetrics.Delete(name)
					}
					return true
				})
			}),
		)),
		k8s.WithResourceController(mpod.New(
			mpod.WithControllerName("pod metrics discoverer"),
			mpod.WithOnErrorFunc(func(err error) {
				log.Error("failed to reconcile:", err)
			}),
			mpod.WithOnReconcileFunc(func(podList map[string]mpod.Pod) {
				log.Debugf("pod metrics reconciled\t%#v", podList)
				for name, pods := range podList {
					d.podMetrics.Store(name, pods)
				}
				d.podMetrics.Range(func(name string, _ mpod.Pod) bool {
					_, ok := podList[name]
					if !ok {
						d.podMetrics.Delete(name)
					}
					return true
				})
			}),
		)),
		k8s.WithResourceController(pod.New(
			pod.WithControllerName("pod discoverer"),
			pod.WithOnErrorFunc(func(err error) {
				log.Error("failed to reconcile:", err)
			}),
			pod.WithOnReconcileFunc(func(podList map[string][]pod.Pod) {
				log.Debugf("pod resource reconciled\t%#v", podList)
				for name, pods := range podList {
					if len(pods) > d.maxPods {
						d.maxPods = len(pods)
					}
					d.pods.Store(name, pods)
				}
				d.pods.Range(func(name string, _ []pod.Pod) bool {
					_, ok := podList[name]
					if !ok {
						d.pods.Delete(name)
					}
					return true
				})
			}),
		)),
		k8s.WithResourceController(node.New(
			node.WithControllerName("node discoverer"),
			node.WithOnErrorFunc(func(err error) {
				log.Error("failed to reconcile:", err)
			}),
			node.WithOnReconcileFunc(func(nodes []node.Node) {
				log.Debugf("node resource reconciled\t%#v", nodes)
				nm := make(map[string]struct{}, len(nodes))
				for _, n := range nodes {
					nm[n.Name] = struct{}{}
					d.nodes.Store(n.Name, n)
				}
				d.nodes.Range(func(name string, _ node.Node) bool {
					_, ok := nm[name]
					if !ok {
						d.nodes.Delete(name)
					}
					return true
				})
			}),
		)),
	)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *discoverer) Start(ctx context.Context) (<-chan error, error) {
	dech, err := d.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, 2)
	d.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		dt := time.NewTicker(d.csd)
		defer dt.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-dt.C:
				var (
					podsByNode      = make(map[string]map[string]map[string][]*payload.Info_Pod) // map[node][namespace][name][]pod
					podsByNamespace = make(map[string]map[string][]*payload.Info_Pod)            // map[namespace][name][]pod
					podsByName      = make(map[string][]*payload.Info_Pod)                       // map[name][]pod
					nodeByName      = make(map[string]*payload.Info_Node)                        // map[name]node
				)

				d.nodes.Range(func(nodeName string, n node.Node) bool {
					select {
					case <-ctx.Done():
						return false
					default:
						ni := &payload.Info_Node{
							Name:         n.Name,
							InternalAddr: n.InternalAddr,
							ExternalAddr: n.ExternalAddr,
							Cpu: &payload.Info_CPU{
								Limit:   n.CPUCapacity,
								Request: n.CPUCapacity - n.CPURemain,
							},
							Memory: &payload.Info_Memory{
								Limit:   n.MemCapacity,
								Request: n.MemCapacity - n.MemRemain,
							},
							Pods: &payload.Info_Pods{
								Pods: make([]*payload.Info_Pod, d.maxPods),
							},
						}
						nm, ok := d.nodeMetrics.Load(nodeName)
						if ok {
							ni.Cpu.Usage = nm.CPU
							ni.Memory.Usage = nm.Mem
						}
						nodeByName[nodeName] = ni
						return true
					}
				})
				d.pods.Range(func(appName string, pods []pod.Pod) bool {
					select {
					case <-ctx.Done():
						return false
					default:
						for _, p := range pods {
							select {
							case <-ctx.Done():
								return false
							default:
								pi := &payload.Info_Pod{
									AppName:   appName,
									Name:      p.Name,
									Namespace: p.Namespace,
									Ip:        p.IP,
									Cpu: &payload.Info_CPU{
										Limit:   p.CPULimit,
										Request: p.CPURequest,
									},
									Memory: &payload.Info_Memory{
										Limit:   p.MemLimit,
										Request: p.MemRequest,
									},
								}
								pm, ok := d.podMetrics.Load(p.Name)
								if ok {
									pi.Cpu.Usage = pm.CPU
									pi.Memory.Usage = pm.Mem
								}
								n, ok := nodeByName[p.NodeName]
								if ok {
									pi.Node = n
								}
								_, ok = podsByNode[p.NodeName]
								if !ok {
									podsByNode[p.NodeName] = make(map[string]map[string][]*payload.Info_Pod, len(nodeByName))
								}
								_, ok = podsByNode[p.NodeName][p.Namespace]
								if !ok {
									podsByNode[p.NodeName][p.Namespace] = make(map[string][]*payload.Info_Pod, len(pods))
								}
								_, ok = podsByNamespace[p.Namespace]
								if !ok {
									podsByNamespace[p.Namespace] = make(map[string][]*payload.Info_Pod, len(pods))
								}
								_, ok = podsByNode[p.NodeName][p.Namespace][appName]
								if !ok {
									podsByNode[p.NodeName][p.Namespace][appName] = make([]*payload.Info_Pod, 0, len(pods))
								}
								_, ok = podsByNamespace[p.Namespace][appName]
								if !ok {
									podsByNamespace[p.Namespace][appName] = make([]*payload.Info_Pod, 0, len(pods))
								}
								_, ok = podsByName[appName]
								if !ok {
									podsByName[appName] = make([]*payload.Info_Pod, 0, len(pods))
								}
								podsByNode[p.NodeName][p.Namespace][appName] = append(podsByNode[p.NodeName][p.Namespace][appName], pi)
								podsByNamespace[p.Namespace][appName] = append(podsByNamespace[p.Namespace][appName], pi)
								podsByName[appName] = append(podsByName[appName], pi)
							}
						}
						return true
					}
				})
				var wg sync.WaitGroup
				wg.Add(1)
				d.eg.Go(safety.RecoverFunc(func() error {
					defer wg.Done()
					for nodeName := range podsByNode {
						for namespace := range podsByNode[nodeName] {
							for appName, p := range podsByNode[nodeName][namespace] {
								sort.Slice(p, func(i, j int) bool {
									return p[i].GetMemory().GetUsage() < p[j].GetMemory().GetUsage()
								})
								podsByNode[nodeName][namespace][appName] = p
								nn, ok := nodeByName[nodeName]
								if !ok || nn == nil {
									nodeByName[nodeName] = new(payload.Info_Node)
									nn, ok = nodeByName[nodeName]
									if !ok {
										continue
									}
								}
								if nn.Pods == nil {
									nodeByName[nodeName].Pods = new(payload.Info_Pods)
								}
								if nn.Pods.Pods == nil {
									nodeByName[nodeName].Pods.Pods = make([]*payload.Info_Pod, 0, len(p))
								}
								nn, ok = nodeByName[nodeName]
								if ok && nn.Pods != nil && nn.Pods.Pods != nil {
									nodeByName[nodeName].Pods.Pods = append(nodeByName[nodeName].Pods.Pods, p...)
								}
							}
						}
						nn, ok := nodeByName[nodeName]
						if ok && nn.Pods != nil && nn.Pods.Pods != nil {
							p := nn.Pods.Pods
							sort.Slice(p, func(i, j int) bool {
								return p[i].GetMemory().GetUsage() < p[j].GetMemory().GetUsage()
							})
							nodeByName[nodeName].Pods.Pods = p
						}
					}
					d.nodeByName.Store(nodeByName)
					d.podsByNode.Store(podsByNode)
					return nil
				}))
				wg.Add(1)
				d.eg.Go(safety.RecoverFunc(func() error {
					defer wg.Done()
					for namespace := range podsByNamespace {
						for appName, p := range podsByNamespace[namespace] {
							sort.Slice(p, func(i, j int) bool {
								return p[i].GetMemory().GetUsage() < p[j].GetMemory().GetUsage()
							})
							podsByNamespace[namespace][appName] = p
						}
					}
					d.podsByNamespace.Store(podsByNamespace)
					return nil
				}))
				wg.Add(1)
				d.eg.Go(safety.RecoverFunc(func() error {
					defer wg.Done()
					for appName, p := range podsByName {
						sort.Slice(p, func(i, j int) bool {
							return p[i].GetMemory().GetUsage() < p[j].GetMemory().GetUsage()
						})
						podsByName[appName] = p
					}
					d.podsByName.Store(podsByName)
					return nil
				}))
				wg.Wait()
			case err = <-dech:
				if err != nil {
					ech <- err
				}
			}

		}
	}))
	return ech, nil
}

func (d *discoverer) GetPods(req *payload.Discoverer_Request) (pods *payload.Info_Pods, err error) {
	var (
		podsByNamespace map[string]map[string][]*payload.Info_Pod
		podsByName      map[string][]*payload.Info_Pod
		ok              bool
	)
	pods = new(payload.Info_Pods)
	if req.GetNode() != "" && req.GetNode() != "*" {
		pbn, ok := d.podsByNode.Load().(map[string]map[string]map[string][]*payload.Info_Pod)
		if !ok {
			return nil, errors.ErrInvalidDiscoveryCache
		}
		podsByNamespace, ok = pbn[req.GetNode()]
		if !ok {
			return nil, errors.ErrNodeNotFound(req.GetNode())
		}
	}
	if req.GetNamespace() != "" && req.GetNamespace() != "*" {
		if podsByNamespace == nil {
			podsByNamespace, ok = d.podsByNamespace.Load().(map[string]map[string][]*payload.Info_Pod)
			if !ok {
				return nil, errors.ErrInvalidDiscoveryCache
			}
		}
		podsByName, ok = podsByNamespace[req.GetNamespace()]
		if !ok {
			return nil, errors.ErrNamespaceNotFound(req.GetNamespace())
		}
	}
	if podsByName == nil {
		if podsByNamespace != nil {
			podsByName = make(map[string][]*payload.Info_Pod)
			for _, pbn := range podsByNamespace {
				for appName, pb := range pbn {
					podsByName[appName] = append(podsByName[appName], pb...)
				}
			}
		} else {
			podsByName, ok = d.podsByName.Load().(map[string][]*payload.Info_Pod)
			if !ok {
				return nil, errors.ErrInvalidDiscoveryCache
			}
		}
	}
	if req.GetName() != "" && req.GetName() != "*" {
		pods.Pods, ok = podsByName[req.GetName()]
		if !ok {
			return nil, errors.ErrPodNameNotFound(req.GetName())
		}
	} else {
		for _, ps := range podsByName {
			pods.Pods = append(pods.Pods, ps...)
		}
	}
	for i := range pods.GetPods() {
		if pods.Pods[i].Node != nil {
			pods.Pods[i].Node.Pods = nil
		}
	}
	return pods, nil
}

func (d *discoverer) GetNodes(req *payload.Discoverer_Request) (nodes *payload.Info_Nodes, err error) {
	nodes = new(payload.Info_Nodes)
	nbn, ok := d.nodeByName.Load().(map[string]*payload.Info_Node)
	if !ok {
		return nil, errors.ErrInvalidDiscoveryCache
	}
	if req.GetNode() != "" && req.GetNode() != "*" {
		n, ok := nbn[req.GetNode()]
		if !ok {
			return nil, errors.ErrNodeNotFound(req.GetNode())
		}
		ps, err := d.GetPods(req)
		if err == nil {
			n.Pods = ps
		}
		nodes.Nodes = append(nodes.Nodes, n)
		return nodes, nil
	}
	ns := nodes.Nodes
	for name, n := range nbn {
		req.Node = name
		if n.Pods != nil {
			n.Pods.Pods = nil
			ps, err := d.GetPods(req)
			if err == nil && ps != nil {
				for i := range ps.Pods {
					ps.Pods[i].Node = nil
				}
				n.Pods = ps
			}
		}
		ns = append(ns, n)
	}
	sort.Slice(ns, func(i, j int) bool {
		return ns[i].GetMemory().GetUsage() < ns[j].GetMemory().GetUsage()
	})
	nodes.Nodes = ns
	return nodes, nil
}
