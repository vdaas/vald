//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"cmp"
	"context"
	"reflect"
	"slices"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	mnode "github.com/vdaas/vald/internal/k8s/metrics/node"
	mpod "github.com/vdaas/vald/internal/k8s/metrics/pod"
	"github.com/vdaas/vald/internal/k8s/node"
	"github.com/vdaas/vald/internal/k8s/pod"
	"github.com/vdaas/vald/internal/k8s/service"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Discoverer interface {
	Start(context.Context) (<-chan error, error)
	GetPods(*payload.Discoverer_Request) (*payload.Info_Pods, error)
	GetNodes(*payload.Discoverer_Request) (*payload.Info_Nodes, error)
	GetServices(*payload.Discoverer_Request) (*payload.Info_Services, error)
}

type discoverer struct {
	maxPods         int
	nodes           sync.Map[string, *node.Node]
	nodeMetrics     sync.Map[string, mnode.Node]
	pods            sync.Map[string, *[]pod.Pod]
	podMetrics      sync.Map[string, mpod.Pod]
	services        sync.Map[string, *service.Service]
	podsByNode      atomic.Pointer[map[string]map[string]map[string][]*payload.Info_Pod]
	podsByNamespace atomic.Pointer[map[string]map[string][]*payload.Info_Pod]
	podsByName      atomic.Pointer[map[string][]*payload.Info_Pod]
	nodeByName      atomic.Pointer[map[string]*payload.Info_Node]
	svcsByName      atomic.Pointer[map[string]*payload.Info_Service]
	ctrl            k8s.Controller
	namespace       string
	name            string
	csd             time.Duration
	der             net.Dialer
	eg              errgroup.Group
}

// New returns Discoverer implementation.
// skipcq: GO-R1005
func New(selector *config.Selectors, opts ...Option) (dsc Discoverer, err error) {
	d := new(discoverer)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(d); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	var (
		podsByNode      = make(map[string]map[string]map[string][]*payload.Info_Pod) // map[node][namespace][name][]pod
		podsByNamespace = make(map[string]map[string][]*payload.Info_Pod)            // map[namespace][name][]pod
		podsByName      = make(map[string][]*payload.Info_Pod)                       // map[name][]pod
		nodeByName      = make(map[string]*payload.Info_Node)                        // map[name]node
		svcsByName      = make(map[string]*payload.Info_Service)                     // map[name]svc
	)
	d.podsByNode.Store(&podsByNode)
	d.podsByNamespace.Store(&podsByNamespace)
	d.podsByName.Store(&podsByName)
	d.nodeByName.Store(&nodeByName)
	d.svcsByName.Store(&svcsByName)

	var k8sOpts []k8s.Option
	k8sOpts = append(k8sOpts,
		k8s.WithDialer(d.der),
		k8s.WithControllerName("vald k8s agent discoverer"),
		k8s.WithLeaderElection(false, "", ""),
		k8s.WithResourceController(mnode.New(
			mnode.WithControllerName("node metrics discoverer"),
			mnode.WithOnErrorFunc(func(err error) {
				log.Error("failed to reconcile node metrics:", err)
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
			mnode.WithNamespace(d.namespace),
			mnode.WithFields(selector.GetNodeMetricsFields()),
			mnode.WithLabels(selector.GetNodeMetricsLabels()),
		)),
		k8s.WithResourceController(mpod.New(
			mpod.WithControllerName("pod metrics discoverer"),
			mpod.WithOnErrorFunc(func(err error) {
				log.Error("failed to reconcile pod metrics:", err)
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
			mpod.WithNamespace(d.namespace),
			mpod.WithFields(selector.GetPodMetricsFields()),
			mpod.WithLabels(selector.GetPodMetricsLabels()),
		)),
		k8s.WithResourceController(pod.New(
			pod.WithControllerName("pod discoverer"),
			pod.WithOnErrorFunc(func(err error) {
				log.Error("failed to reconcile pod resource:", err)
			}),
			pod.WithOnReconcileFunc(func(_ context.Context, podList map[string][]pod.Pod) {
				log.Debugf("pod resource reconciled\t%#v", podList)
				for name, pods := range podList {
					if len(pods) > d.maxPods {
						d.maxPods = len(pods)
					}
					d.pods.Store(name, &pods)
				}
				d.pods.Range(func(name string, _ *[]pod.Pod) bool {
					_, ok := podList[name]
					if !ok {
						d.pods.Delete(name)
					}
					return true
				})
			}),
			pod.WithNamespace(d.namespace),
			pod.WithFields(selector.GetPodFields()),
			pod.WithLabels(selector.GetPodLabels()),
		)),
		k8s.WithResourceController(node.New(
			node.WithControllerName("node discoverer"),
			node.WithOnErrorFunc(func(err error) {
				log.Error("failed to reconcile node resource:", err)
			}),
			node.WithOnReconcileFunc(func(nodes []node.Node) {
				log.Debugf("node resource reconciled\t%#v", nodes)
				nm := make(map[string]struct{}, len(nodes))
				for _, n := range nodes {
					nm[n.Name] = struct{}{}
					d.nodes.Store(n.Name, &n)
				}
				d.nodes.Range(func(name string, _ *node.Node) bool {
					_, ok := nm[name]
					if !ok {
						d.nodes.Delete(name)
					}
					return true
				})
			}),
			node.WithNamespace(d.namespace),
			node.WithFields(selector.GetNodeFields()),
			node.WithLabels(selector.GetNodeLabels()),
		)),
		// Only required when service reconciation is required like read replica.
		// k8s.WithResourceController(service.New(
		// 	service.WithControllerName("service discoverer"),
		// 	service.WithOnErrorFunc(func(err error) {
		// 		log.Error("failed to reconcile:", err)
		// 	}),
		// 	service.WithOnReconcileFunc(func(svcs []service.Service) {
		// 		log.Debugf("svc resource reconciled\t%#v", svcs)
		// 		svcsmap := make(map[string]struct{}, len(svcs))
		// 		for i := range svcs {
		// 			svc := &svcs[i]
		// 			svcsmap[svc.Name] = struct{}{}
		// 			d.services.Store(svc.Name, svc)
		// 		}
		// 		d.services.Range(func(name string, _ *service.Service) bool {
		// 			_, ok := svcsmap[name]
		// 			if !ok {
		// 				d.services.Delete(name)
		// 			}
		// 			return true
		// 		})
		// 	}),
		// 	service.WithNamespace(d.namespace),
		// 	service.WithFields(selector.GetServiceFields()),
		// 	service.WithLabels(selector.GetServiceLabels()),
		// )),
	)

	d.ctrl, err = k8s.New(k8sOpts...)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Start starts the discoverer.
// skipcq: GO-R1005
func (d *discoverer) Start(ctx context.Context) (<-chan error, error) {
	dech, err := d.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, 2)
	d.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		d.der.StartDialerCache(ctx)
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
					svcsByName      = make(map[string]*payload.Info_Service)                     // map[name]svc
				)

				d.nodes.Range(func(nodeName string, n *node.Node) bool {
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
							ni.GetCpu().Usage = nm.CPU
							ni.GetMemory().Usage = nm.Mem
						}
						nodeByName[nodeName] = ni
						return true
					}
				})
				d.pods.Range(func(appName string, pods *[]pod.Pod) bool {
					select {
					case <-ctx.Done():
						return false
					default:
						for _, p := range *pods {
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
									pi.GetCpu().Usage = pm.CPU
									pi.GetMemory().Usage = pm.Mem
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
									podsByNode[p.NodeName][p.Namespace] = make(map[string][]*payload.Info_Pod, len(*pods))
								}
								_, ok = podsByNamespace[p.Namespace]
								if !ok {
									podsByNamespace[p.Namespace] = make(map[string][]*payload.Info_Pod, len(*pods))
								}
								_, ok = podsByNode[p.NodeName][p.Namespace][appName]
								if !ok {
									podsByNode[p.NodeName][p.Namespace][appName] = make([]*payload.Info_Pod, 0, len(*pods))
								}
								_, ok = podsByNamespace[p.Namespace][appName]
								if !ok {
									podsByNamespace[p.Namespace][appName] = make([]*payload.Info_Pod, 0, len(*pods))
								}
								_, ok = podsByName[appName]
								if !ok {
									podsByName[appName] = make([]*payload.Info_Pod, 0, len(*pods))
								}
								podsByNode[p.NodeName][p.Namespace][appName] = append(podsByNode[p.NodeName][p.Namespace][appName], pi)
								podsByNamespace[p.Namespace][appName] = append(podsByNamespace[p.Namespace][appName], pi)
								podsByName[appName] = append(podsByName[appName], pi)
							}
						}
						return true
					}
				})
				d.services.Range(func(key string, svc *service.Service) bool {
					select {
					case <-ctx.Done():
						return false
					default:
						var ports []*payload.Info_ServicePort
						for _, p := range svc.Ports {
							ports = append(ports, &payload.Info_ServicePort{
								Name: p.Name,
								Port: p.Port,
							})
						}
						ni := &payload.Info_Service{
							Name:       svc.Name,
							ClusterIp:  svc.ClusterIP,
							ClusterIps: svc.ClusterIPs,
							Ports:      ports,
							Labels: &payload.Info_Labels{
								Labels: svc.Labels,
							},
							Annotations: &payload.Info_Annotations{
								Annotations: svc.Annotations,
							},
						}
						svcsByName[svc.Name] = ni
						return true
					}
				})
				d.svcsByName.Store(&svcsByName)

				var wg sync.WaitGroup
				wg.Add(1)
				d.eg.Go(safety.RecoverFunc(func() error {
					defer wg.Done()
					for nodeName := range podsByNode {
						for namespace := range podsByNode[nodeName] {
							for appName, p := range podsByNode[nodeName][namespace] {
								slices.SortFunc(p, func(left, right *payload.Info_Pod) int {
									return cmp.Compare(left.GetMemory().GetUsage(), right.GetMemory().GetUsage())
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
								if nn.GetPods() == nil {
									nodeByName[nodeName].Pods = new(payload.Info_Pods)
								}
								if nn.GetPods().GetPods() == nil {
									nodeByName[nodeName].GetPods().Pods = make([]*payload.Info_Pod, 0, len(p))
								}
								nn, ok = nodeByName[nodeName]
								if ok && nn.GetPods() != nil && nn.GetPods().GetPods() != nil {
									nodeByName[nodeName].GetPods().Pods = append(nodeByName[nodeName].GetPods().GetPods(), p...)
								}
							}
						}
						nn, ok := nodeByName[nodeName]
						if ok && nn.GetPods() != nil && nn.GetPods().GetPods() != nil {
							p := nn.GetPods().Pods
							slices.SortFunc(p, func(left, right *payload.Info_Pod) int {
								return cmp.Compare(left.GetMemory().GetUsage(), right.GetMemory().GetUsage())
							})
							nodeByName[nodeName].GetPods().Pods = p
						}
					}
					d.nodeByName.Store(&nodeByName)
					d.podsByNode.Store(&podsByNode)
					return nil
				}))
				wg.Add(1)
				d.eg.Go(safety.RecoverFunc(func() error {
					defer wg.Done()
					for namespace := range podsByNamespace {
						for appName, p := range podsByNamespace[namespace] {
							slices.SortFunc(p, func(left, right *payload.Info_Pod) int {
								return cmp.Compare(left.GetMemory().GetUsage(), right.GetMemory().GetUsage())
							})
							podsByNamespace[namespace][appName] = p
						}
					}
					d.podsByNamespace.Store(&podsByNamespace)
					return nil
				}))
				wg.Add(1)
				d.eg.Go(safety.RecoverFunc(func() error {
					defer wg.Done()
					for appName, p := range podsByName {
						slices.SortFunc(p, func(left, right *payload.Info_Pod) int {
							return cmp.Compare(left.GetMemory().GetUsage(), right.GetMemory().GetUsage())
						})
						podsByName[appName] = p
					}
					d.podsByName.Store(&podsByName)
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
		pbn := *d.podsByNode.Load()
		if pbn == nil {
			return nil, errors.ErrInvalidDiscoveryCache
		}
		podsByNamespace, ok = pbn[req.GetNode()]
		if !ok {
			return nil, errors.ErrNodeNotFound(req.GetNode())
		}
	}
	if req.GetNamespace() != "" && req.GetNamespace() != "*" {
		if podsByNamespace == nil {
			podsByNamespace = *d.podsByNamespace.Load()
			if podsByNamespace == nil {
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
			podsByName = *d.podsByName.Load()
			if podsByName == nil {
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
			pods.Pods = append(pods.GetPods(), ps...)
		}
	}
	for i := range pods.GetPods() {
		if pods.GetPods()[i].GetNode() != nil {
			pods.GetPods()[i].GetNode().Pods = nil
		}
	}
	slices.SortFunc(pods.Pods, func(left, right *payload.Info_Pod) int {
		return cmp.Compare(left.GetMemory().GetUsage(), right.GetMemory().GetUsage())
	})
	return pods, nil
}

func (d *discoverer) GetNodes(
	req *payload.Discoverer_Request,
) (nodes *payload.Info_Nodes, err error) {
	nbn := *d.nodeByName.Load()
	if nbn == nil {
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
		return &payload.Info_Nodes{
			Nodes: []*payload.Info_Node{
				n,
			},
		}, nil
	}
	nodes = &payload.Info_Nodes{
		Nodes: make([]*payload.Info_Node, 0, len(nbn)),
	}
	for name, n := range nbn {
		req.Node = name
		if n.GetPods() != nil {
			n.GetPods().Pods = nil
			ps, err := d.GetPods(req)
			if err == nil && ps != nil {
				for i := range ps.Pods {
					ps.GetPods()[i].Node = nil
				}
				slices.SortFunc(ps.Pods, func(left, right *payload.Info_Pod) int {
					return cmp.Compare(left.GetMemory().GetUsage(), right.GetMemory().GetUsage())
				})
				n.Pods = ps
			}
		}
		nodes.Nodes = append(nodes.Nodes, n)
	}
	slices.SortFunc(nodes.Nodes, func(left, right *payload.Info_Node) int {
		return cmp.Compare(left.GetMemory().GetUsage(), right.GetMemory().GetUsage())
	})
	return nodes, nil
}

// Get Services returns the services that matches the request.
func (d *discoverer) GetServices(
	req *payload.Discoverer_Request,
) (svcs *payload.Info_Services, err error) {
	svcs = new(payload.Info_Services)
	sbn := *d.svcsByName.Load()
	if sbn == nil {
		return nil, errors.ErrInvalidDiscoveryCache
	}

	if req.GetName() != "" && req.GetName() != "*" {
		v, ok := sbn[req.GetName()]
		if !ok {
			return nil, errors.ErrSvcNameNotFound(req.GetName())
		}
		svcs.Services = append(svcs.Services, v)
	} else {
		for _, svc := range sbn {
			svcs.Services = append(svcs.Services, svc)
		}
	}

	return svcs, nil
}
