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

// Package discoverer
package discoverer

import (
	"cmp"
	"context"
	"reflect"
	"slices"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Client interface {
	Start(ctx context.Context) (<-chan error, error)
	GetAddrs(ctx context.Context) []string

	// GetClient returns the grpc.Client for both read and write.
	GetClient() grpc.Client

	// GetReadClient returns the grpc.Client only for read. If there's no readreplica, this returns the grpc.Client for the primary agent.
	// Use this API only for getting client for agent. For other use cases, use GetClient() instead.
	// Internally, this API round robin between c.client and c.readClient with the ratio of
	// agent replicas and read replica agent replicas.
	GetReadClient() grpc.Client
}

type client struct {
	autoconn     bool
	onDiscover   func(ctx context.Context, c Client, addrs []string) error
	onConnect    func(ctx context.Context, c Client, addr string) error
	onDisconnect func(ctx context.Context, c Client, addr string) error
	client       grpc.Client
	dns          string
	opts         []grpc.Option
	port         int
	addrs        atomic.Pointer[[]string]
	dscClient    grpc.Client
	dscDur       time.Duration
	eg           errgroup.Group
	name         string
	namespace    string
	nodeName     string
	// read replica related members below
	readClient          grpc.Client
	readReplicaReplicas uint64
	roundRobin          atomic.Uint64
}

func New(opts ...Option) (d Client, err error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return c, nil
}

// Start starts the discoverer client.
// skipcq: GO-R1005
func (c *client) Start(ctx context.Context) (<-chan error, error) {
	dech, err := c.dscClient.StartConnectionMonitor(ctx)
	if err != nil {
		return nil, err
	}

	var rrech <-chan error
	if c.readClient != nil {
		rrech, err = c.readClient.StartConnectionMonitor(ctx)
		if err != nil {
			return nil, err
		}
	}

	addrs, err := c.updateDiscoveryInfo(ctx)
	if err != nil || len(addrs) == 0 {
		addrs, err = c.dnsDiscovery(ctx)
		if err != nil {
			return nil, err
		}
	}

	var aech <-chan error
	if c.client == nil {
		c.client = grpc.New(
			"Agent Client",
			append(
				c.opts,
				grpc.WithAddrs(addrs...),
				grpc.WithErrGroup(c.eg),
			)...,
		)
		aech, err = c.client.StartConnectionMonitor(ctx)
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			if c.onConnect != nil {
				err = c.onConnect(ctx, c, addr)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		for _, addr := range addrs {
			err = c.connect(ctx, addr)
			if err != nil {
				return nil, err
			}
		}
		aech, err = c.client.StartConnectionMonitor(ctx)
	}
	if err != nil {
		return nil, err
	}
	c.addrs.Store(&addrs)

	err = c.discover(ctx)
	if err != nil {
		return nil, errors.Join(c.dscClient.Close(ctx), err)
	}
	ech := make(chan error, 100)
	c.eg.Go(safety.RecoverFunc(func() (err error) {
		dt := time.NewTicker(c.dscDur)
		defer dt.Stop()
		finalize := func() (err error) {
			var errs error
			err = c.dscClient.Close(ctx)
			if err != nil {
				errs = errors.Join(errs, err)
			}
			if c.autoconn && c.client != nil {
				err = c.client.Close(ctx)
				if err != nil {
					errs = errors.Join(errs, err)
				}
			}
			err = ctx.Err()
			if err != nil && err != context.Canceled {
				errs = errors.Join(errs, err)
			}
			return errs
		}
		for {
			select {
			case <-ctx.Done():
				return finalize()
			case err = <-dech:
			case err = <-aech:
			case err = <-rrech:
			case <-dt.C:
				err = c.discover(ctx)
			}
			if err != nil {
				log.Error(err)
				select {
				case <-ctx.Done():
					return finalize()
				case ech <- err:
				}
				err = nil
			}
		}
	}))
	return ech, nil
}

func (c *client) GetAddrs(ctx context.Context) (addrs []string) {
	a := c.addrs.Load()
	if a == nil {
		var err error
		addrs, err = c.dnsDiscovery(ctx)
		if err != nil {
			return nil
		}
	} else {
		addrs = *a
	}
	return addrs
}

func (c *client) GetClient() grpc.Client {
	return c.client
}

func (c *client) GetReadClient() grpc.Client {
	// just return write client when there is no read replica
	if c.readClient == nil {
		return c.client
	}

	var next uint64
	for {
		cur := c.roundRobin.Load()
		next = (cur + 1) % (c.readReplicaReplicas + 1)
		if c.roundRobin.CompareAndSwap(cur, next) {
			break
		}
	}
	if next == 0 {
		return c.client
	}
	return c.readClient
}

func (c *client) connect(ctx context.Context, addr string) (err error) {
	if c.autoconn && c.client != nil {
		_, err = c.client.Connect(ctx, addr, c.client.GetDialOption()...)
		if err != nil {
			return err
		}
		if c.onConnect != nil {
			err = c.onConnect(ctx, c, addr)
		}
	}
	return
}

func (c *client) disconnect(ctx context.Context, addr string) (err error) {
	if c.autoconn && c.client != nil {
		err = c.client.Disconnect(ctx, addr)
		if err == nil && c.onDisconnect != nil {
			err = c.onDisconnect(ctx, c, addr)
		}
	}
	return
}

func (c *client) dnsDiscovery(ctx context.Context) (addrs []string, err error) {
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, c.dns)
	if err != nil || len(ips) == 0 {
		return nil, errors.ErrAddrCouldNotDiscover(err, c.dns)
	}
	log.Debugf("dns discovery succeeded for dns = %s", c.dns)
	addrs = make([]string, 0, len(ips))
	for _, ip := range ips {
		addr := net.JoinHostPort(ip.String(), uint16(c.port))
		if err = c.connect(ctx, addr); err != nil {
			log.Debugf("dns discovery connect for addr = %s from dns = %s failed %v", addr, c.dns, err)
		} else {
			log.Debugf("dns discovery connect for addr = %s from dns = %s succeeded", addr, c.dns)
			addrs = append(addrs, addr)
		}
	}
	if len(addrs) == 0 {
		return nil, errors.ErrAddrCouldNotDiscover(err, c.dns)
	}
	if len(addrs) != 0 && c.onDiscover != nil {
		return addrs, c.onDiscover(ctx, c, addrs)
	}
	return addrs, nil
}

func (c *client) discover(ctx context.Context) (err error) {
	if c.dscClient == nil || (c.autoconn && c.client == nil) {
		return errors.ErrGRPCClientNotFound
	}

	var connected []string
	if bo := c.client.GetBackoff(); bo != nil {
		_, err = bo.Do(ctx, func(ctx context.Context) (any, bool, error) {
			connected, err = c.updateDiscoveryInfo(ctx)
			if err != nil {
				if errors.IsNot(err, errors.ErrGRPCClientNotFound, errors.ErrGRPCClientConnNotFound("*")) {
					return nil, true, err
				}
				return nil, false, err
			}
			return nil, false, nil
		})
	} else {
		connected, err = c.updateDiscoveryInfo(ctx)
	}
	if err != nil {
		log.Warnf("failed to discover addrs from discoverer API, error: %v,\ttrying to dns discovery from %s...", err, c.dns)
		connected, err = c.dnsDiscovery(ctx)
		if err != nil {
			return err
		}
	}

	oldAddrs := c.GetAddrs(ctx)
	c.addrs.Store(&connected)
	return c.disconnectOldAddrs(ctx, oldAddrs, connected)
}

func (c *client) updateDiscoveryInfo(ctx context.Context) (connected []string, err error) {
	nodes, err := c.discoverNodes(ctx)
	if err != nil {
		log.Warnf("error detected when discovering nodes,\terrors: %v", err)
		return nil, err
	}
	if nodes == nil {
		log.Warn("no nodes found")
		return nil, errors.ErrNodeNotFound("all")
	}
	connected, err = c.discoverAddrs(ctx, nodes)
	if err != nil {
		return nil, err
	}
	if len(connected) == 0 {
		log.Warnf("connected addr is zero,\tnodes: %s", nodes.String())
		return nil, errors.ErrAddrCouldNotDiscover(err, c.dns)
	}
	if c.onDiscover != nil {
		err = c.onDiscover(ctx, c, connected)
		if err != nil {
			log.Warnf("error detected when onDiscover execution,\tnodes: %s,\tconnected: %v,\terrors: %v", nodes.String(), connected, err)
			return nil, err
		}
	}
	return connected, nil
}

func (c *client) discoverNodes(ctx context.Context) (nodes *payload.Info_Nodes, err error) {
	nodes, err = grpc.RoundRobin(c.dscClient, grpc.WithGRPCMethod(ctx, "discoverer.v1.Discoverer/Nodes"), func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (*payload.Info_Nodes, error) {
		return discoverer.NewDiscovererClient(conn).
			Nodes(ctx, &payload.Discoverer_Request{
				Namespace: c.namespace,
				Name:      c.name,
				Node:      c.nodeName,
			}, copts...)
	})
	if err != nil {
		return nil, err
	}
	slices.SortFunc(nodes.Nodes, func(left, right *payload.Info_Node) int {
		if left.GetMemory() == nil || right.GetMemory() == nil {
			return 0 // Default comparison value; adjust as needed.
		}
		return cmp.Compare(left.GetMemory().GetUsage(), right.GetMemory().GetUsage())
	})
	return nodes, nil
}

func (c *client) discoverAddrs(
	ctx context.Context, nodes *payload.Info_Nodes,
) (addrs []string, err error) {
	if nodes == nil {
		return nil, errors.ErrAddrCouldNotDiscover(err, c.dns)
	}
	maxPodLen := 0
	podLength := 0
	for i, node := range nodes.GetNodes() {
		if node != nil && node.GetPods() != nil && node.GetPods().GetPods() != nil {
			l := len(node.GetPods().GetPods())
			podLength += l
			if l > maxPodLen {
				maxPodLen = l
			}
			slices.SortFunc(nodes.Nodes[i].Pods.Pods, func(left, right *payload.Info_Pod) int {
				if left.GetMemory() == nil || right.GetMemory() == nil {
					return 0 // Default comparison value; adjust as needed.
				}
				return cmp.Compare(left.GetMemory().GetUsage(), right.GetMemory().GetUsage())
			})
		}
	}
	addrs = make([]string, 0, podLength)
	for i := 0; i < maxPodLen; i++ {
		for _, node := range nodes.GetNodes() {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				if node != nil &&
					node.GetPods() != nil &&
					len(node.GetPods().GetPods()) > i &&
					len(node.GetPods().GetPods()[i].GetIp()) != 0 {
					addr := net.JoinHostPort(node.GetPods().GetPods()[i].GetIp(), uint16(c.port))
					if err = c.connect(ctx, addr); err != nil {
						log.Debugf("resource based discovery connect from discoverer API for addr = %s failed %v", addr, errors.ErrAddrCouldNotDiscover(err, addr))
						err = nil
					} else {
						addrs = append(addrs, addr)
					}

				}
			}
		}
	}
	return addrs, nil
}

func (c *client) disconnectOldAddrs(
	ctx context.Context, oldAddrs, connectedAddrs []string,
) (err error) {
	if !c.autoconn {
		return nil
	}
	var cur sync.Map[string, any]
	for _, addr := range connectedAddrs {
		cur.Store(addr, struct{}{})
	}

	for _, old := range oldAddrs {
		_, ok := cur.Load(old)
		if !ok {
			c.eg.Go(safety.RecoverFunc(func() error {
				err = c.disconnect(ctx, old)
				if err != nil {
					log.Error(err)
				}
				return nil
			}))
		}
	}
	if c.client != nil {
		if err = c.client.RangeConcurrent(ctx, len(connectedAddrs)/3, func(ctx context.Context,
			addr string,
			conn *grpc.ClientConn,
			copts ...grpc.CallOption,
		) (err error) {
			_, ok := cur.Load(addr)
			if !ok {
				err = c.disconnect(ctx, addr)
				if err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			log.Error(err)
		}
	}
	return nil
}
