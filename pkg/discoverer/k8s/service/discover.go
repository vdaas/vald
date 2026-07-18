//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/k8s"
	"github.com/vdaas/vald/internal/k8s/reconciler"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Discoverer interface {
	Start(context.Context) (<-chan error, error)
	GetPods(*payload.Discoverer_Request) (*payload.Info_Pods, error)
	GetNodes(*payload.Discoverer_Request) (*payload.Info_Nodes, error)
	GetServices(*payload.Discoverer_Request) (*payload.Info_Services, error)
}

type discoverer struct {
	eg              errgroup.Group
	der             net.Dialer
	ctrl            k8s.Controller
	podsByName      atomic.Pointer[map[string][]*payload.Info_Pod]
	svcsByName      atomic.Pointer[map[string]*payload.Info_Service]
	nodeByName      atomic.Pointer[map[string]*payload.Info_Node]
	podsByNode      atomic.Pointer[map[string]map[string]map[string][]*payload.Info_Pod]
	podsByNamespace atomic.Pointer[map[string]map[string][]*payload.Info_Pod]
	namespace       string
	name            string
	pods            sync.Map[string, *[]Pod]
	nodeMetrics     sync.Map[string, NodeMetrics]
	nodes           sync.Map[string, *Node]
	podMetrics      sync.Map[string, PodMetrics]
	maxPods         int
	csd             time.Duration
}

// New returns Discoverer implementation.
// skipcq: GO-R1005
func New(selector *config.Selectors, opts ...Option) (Discoverer, error) {
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

	k8sOpts := []k8s.Option{
		k8s.WithDialer(d.der),
		k8s.WithControllerName("vald k8s agent discoverer"),
		k8s.WithLeaderElection(false, "", ""),
		listSyncController("node metrics discoverer", "node metrics",
			new(k8s.NodeMetrics), &d.nodeMetrics,
			debugLogged("node metrics", toNodeMetricsMap),
			d.namespace, selector.GetNodeMetricsFields(), selector.GetNodeMetricsLabels(),
			reconciler.WithAddToScheme(k8s.AddMetricsToScheme),
		),
		listSyncController("pod metrics discoverer", "pod metrics",
			new(k8s.PodMetrics), &d.podMetrics,
			debugLogged("pod metrics", toPodMetricsMap),
			d.namespace, selector.GetPodMetricsFields(), selector.GetPodMetricsLabels(),
			reconciler.WithAddToScheme(k8s.AddMetricsToScheme),
			reconciler.WithFieldIndex("containers.name", podMetricsContainersNameIndexer),
		),
		listSyncController("pod discoverer", "pod resource",
			new(k8s.Pod), &d.pods,
			d.podEntries,
			d.namespace, selector.GetPodFields(), selector.GetPodLabels(),
			reconciler.WithFieldIndex("status.phase", podStatusPhaseIndexer),
		),
		listSyncController("node discoverer", "node resource",
			new(k8s.Node), &d.nodes,
			nodeEntries,
			d.namespace, selector.GetNodeFields(), selector.GetNodeLabels(),
			reconciler.WithFieldIndex("status.phase", nodeStatusPhaseIndexer),
		),
	}

	ctrl, err := k8s.New(k8sOpts...)
	if err != nil {
		return nil, err
	}
	d.ctrl = ctrl
	return d, nil
}

// startErrChanBufferSize buffers the controller error stream plus this
// goroutine's own terminal error so neither send can block shutdown.
const startErrChanBufferSize = 2

// Start starts the discoverer.
func (d *discoverer) Start(ctx context.Context) (<-chan error, error) {
	dech, err := d.ctrl.Start(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, startErrChanBufferSize)
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
				d.refreshSnapshots(ctx)
			case err = <-dech:
				if err != nil {
					ech <- err
				}
			}
		}
	}))
	return ech, nil
}

// refreshSnapshots rebuilds the discovery snapshots (podsByNode,
// podsByNamespace, podsByName and nodeByName) from the reconciled watch state
// and atomically publishes them. It is called once per discover tick and
// returns only after every snapshot has been sorted and stored.
// skipcq: GO-R1005
//
//nolint:maintidx // complexity comes from the enumerated resource branches in the sequential relay loop, not tangled logic
func (d *discoverer) refreshSnapshots(ctx context.Context) {
	var (
		podsByNode      = make(map[string]map[string]map[string][]*payload.Info_Pod) // map[node][namespace][name][]pod
		podsByNamespace = make(map[string]map[string][]*payload.Info_Pod)            // map[namespace][name][]pod
		podsByName      = make(map[string][]*payload.Info_Pod)                       // map[name][]pod
		nodeByName      = make(map[string]*payload.Info_Node)                        // map[name]node
	)

	d.nodes.Range(func(nodeName string, n *Node) bool {
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
	d.pods.Range(func(appName string, pods *[]Pod) bool {
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
					// No capacity hints for the per-node maps and the group
					// slices below: len(*pods) is the application's
					// cluster-wide pod count, so reserving it per
					// (node, namespace, app) or (namespace, app) group
					// over-allocated by the number of nodes/namespaces the
					// app spans (~70% of the tick's allocated bytes at
					// 50 nodes / 1000 pods); let append grow them naturally.
					nodeNs, ok := podsByNode[p.NodeName]
					if !ok {
						nodeNs = make(map[string]map[string][]*payload.Info_Pod)
						podsByNode[p.NodeName] = nodeNs
					}
					nodeApps, ok := nodeNs[p.Namespace]
					if !ok {
						nodeApps = make(map[string][]*payload.Info_Pod)
						nodeNs[p.Namespace] = nodeApps
					}
					nsApps, ok := podsByNamespace[p.Namespace]
					if !ok {
						nsApps = make(map[string][]*payload.Info_Pod, len(*pods))
						podsByNamespace[p.Namespace] = nsApps
					}
					group, ok := podsByName[appName]
					if !ok {
						group = make([]*payload.Info_Pod, 0, len(*pods))
					}
					nodeApps[appName] = append(nodeApps[appName], pi)
					nsApps[appName] = append(nsApps[appName], pi)
					podsByName[appName] = append(group, pi)
				}
			}
			return true
		}
	})

	// The three goroutines below each sort one snapshot's group slices. They
	// never share a backing array: every group above is built by appending to
	// its own allocation, so the concurrent in-place sorts only share the
	// read-only *payload.Info_Pod elements.
	var wg sync.WaitGroup
	wg.Add(1)
	d.eg.Go(safety.RecoverFunc(func() error {
		defer wg.Done()
		for nodeName, nodeNs := range podsByNode {
			nn, ok := nodeByName[nodeName]
			if !ok || nn == nil {
				nn = new(payload.Info_Node)
				nodeByName[nodeName] = nn
			}
			if nn.GetPods() == nil {
				nn.Pods = new(payload.Info_Pods)
			}
			aggregated := nn.GetPods().GetPods()
			for _, nodeApps := range nodeNs {
				for _, p := range nodeApps {
					slices.SortFunc(p, byMemoryUsage)
					aggregated = append(aggregated, p...)
				}
			}
			// Sorting the aggregate is not redundant with the per-group
			// sorts above: it orders the concatenation of the sorted groups.
			slices.SortFunc(aggregated, byMemoryUsage)
			nn.GetPods().Pods = aggregated
		}
		d.nodeByName.Store(&nodeByName)
		d.podsByNode.Store(&podsByNode)
		return nil
	}))
	wg.Add(1)
	d.eg.Go(safety.RecoverFunc(func() error {
		defer wg.Done()
		for _, nsApps := range podsByNamespace {
			for _, p := range nsApps {
				slices.SortFunc(p, byMemoryUsage)
			}
		}
		d.podsByNamespace.Store(&podsByNamespace)
		return nil
	}))
	wg.Add(1)
	d.eg.Go(safety.RecoverFunc(func() error {
		defer wg.Done()
		for _, p := range podsByName {
			slices.SortFunc(p, byMemoryUsage)
		}
		d.podsByName.Store(&podsByName)
		return nil
	}))
	wg.Wait()
}

// byMemoryUsage orders pods by ascending memory usage. Every published pod
// snapshot group keeps this order as the snapshot's distribution-order
// contract; read handlers sort response-only copies of the groups, so
// concurrent requests never re-sort (or otherwise mutate) the shared slices.
func byMemoryUsage(left, right *payload.Info_Pod) int {
	return cmp.Compare(left.GetMemory().GetUsage(), right.GetMemory().GetUsage())
}

// responseNode returns a response-only shallow copy of the shared snapshot
// node n with the recursive per-node pod aggregate dropped. Snapshot nodes
// are shared by every pod on the node and by every concurrent request, so
// read handlers must never write to them. The leaf Cpu/Memory messages stay
// shared: nothing mutates them after publication.
func responseNode(n *payload.Info_Node) *payload.Info_Node {
	if n == nil {
		return nil
	}
	return &payload.Info_Node{
		Name:         n.GetName(),
		InternalAddr: n.GetInternalAddr(),
		ExternalAddr: n.GetExternalAddr(),
		Cpu:          n.GetCpu(),
		Memory:       n.GetMemory(),
	}
}

// responsePod returns a response-only shallow copy of the shared snapshot pod
// p whose node link is replaced by a response-only copy without the recursive
// pod aggregate.
func responsePod(p *payload.Info_Pod) *payload.Info_Pod {
	if p == nil {
		return nil
	}
	return &payload.Info_Pod{
		AppName:   p.GetAppName(),
		Name:      p.GetName(),
		Namespace: p.GetNamespace(),
		Ip:        p.GetIp(),
		Cpu:       p.GetCpu(),
		Memory:    p.GetMemory(),
		Node:      responseNode(p.GetNode()),
	}
}

// stripResponsePodNodes drops the redundant node link from every pod of a
// node-scoped response: the pods are already nested under their node, and
// keeping the link would nest the node twice (the previous in-place
// implementation even made the single-node GetNodes response cyclic). ps
// holds response-only copies from GetPods, so the writes stay inside the
// response.
func stripResponsePodNodes(ps *payload.Info_Pods) {
	for _, p := range ps.GetPods() {
		p.Node = nil
	}
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
	// The snapshot groups and the *Info_Pod / *Info_Node messages they point
	// to are shared with every other concurrent request, so the response is
	// built from response-only copies: both the recursion-stripping inside
	// responsePod and the sort below must never touch the shared snapshot.
	if req.GetName() != "" && req.GetName() != "*" {
		group, ok := podsByName[req.GetName()]
		if !ok {
			return nil, errors.ErrPodNameNotFound(req.GetName())
		}
		pods.Pods = make([]*payload.Info_Pod, 0, len(group))
		for _, p := range group {
			pods.Pods = append(pods.Pods, responsePod(p))
		}
	} else {
		for _, ps := range podsByName {
			for _, p := range ps {
				pods.Pods = append(pods.Pods, responsePod(p))
			}
		}
	}
	slices.SortFunc(pods.Pods, byMemoryUsage)
	return pods, nil
}

func (d *discoverer) GetNodes(
	req *payload.Discoverer_Request,
) (nodes *payload.Info_Nodes, err error) {
	nbn := *d.nodeByName.Load()
	if nbn == nil {
		return nil, errors.ErrInvalidDiscoveryCache
	}
	// Snapshot nodes are shared across every concurrent request and with the
	// pods' node links, so each response node is a response-only copy whose
	// nested pod list comes from GetPods' response-only, pre-sorted copies.
	if req.GetNode() != "" && req.GetNode() != "*" {
		n, ok := nbn[req.GetNode()]
		if !ok {
			return nil, errors.ErrNodeNotFound(req.GetNode())
		}
		rn := responseNode(n)
		ps, err := d.GetPods(req)
		if err == nil {
			stripResponsePodNodes(ps)
			rn.Pods = ps
		}
		return &payload.Info_Nodes{
			Nodes: []*payload.Info_Node{
				rn,
			},
		}, nil
	}
	nodes = &payload.Info_Nodes{
		Nodes: make([]*payload.Info_Node, 0, len(nbn)),
	}
	for name, n := range nbn {
		req.Node = name
		rn := responseNode(n)
		if n.GetPods() != nil {
			// An empty, non-nil aggregate mirrors the previous response shape
			// for nodes without pods in the snapshot (GetPods returns
			// node-not-found for them).
			rn.Pods = new(payload.Info_Pods)
			ps, err := d.GetPods(req)
			if err == nil && ps != nil {
				stripResponsePodNodes(ps)
				rn.Pods = ps
			}
		}
		nodes.Nodes = append(nodes.Nodes, rn)
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
