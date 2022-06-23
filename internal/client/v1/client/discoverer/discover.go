//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package discoverer
package discoverer

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
)

type Client interface {
	Start(ctx context.Context) (<-chan error, error)
	GetAddrs(ctx context.Context) []string
	GetClient() grpc.Client
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
	addrs        atomic.Value
	dscClient    grpc.Client
	dscDur       time.Duration
	eg           errgroup.Group
	name         string
	namespace    string
	nodeName     string
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

func (c *client) Start(ctx context.Context) (<-chan error, error) {
	dech, err := c.dscClient.StartConnectionMonitor(ctx)
	if err != nil {
		return nil, err
	}

	ech := make(chan error, 100)
	addrs, err := c.dnsDiscovery(ctx, ech)
	if err != nil {
		close(ech)
		return nil, err
	}
	c.addrs.Store(addrs)

	var aech <-chan error
	if c.autoconn {
		c.client = grpc.New(
			append(
				c.opts,
				grpc.WithAddrs(addrs...),
				grpc.WithErrGroup(c.eg),
			)...,
		)
		if c.client != nil {
			aech, err = c.client.StartConnectionMonitor(ctx)
			if err != nil {
				close(ech)
				return nil, err
			}
		}
	}

	err = c.discover(ctx, ech)
	if err != nil {
		close(ech)
		return nil, errors.Wrap(c.dscClient.Close(ctx), err.Error())
	}

	c.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		dt := time.NewTicker(c.dscDur)
		defer dt.Stop()
		finalize := func() (err error) {
			var errs error
			err = c.dscClient.Close(ctx)
			if err != nil {
				errs = errors.Wrap(errs, err.Error())
			}
			if c.autoconn && c.client != nil {
				err = c.client.Close(ctx)
				if err != nil {
					errs = errors.Wrap(errs, err.Error())
				}
			}
			err = ctx.Err()
			if err != nil && err != context.Canceled {
				errs = errors.Wrap(errs, err.Error())
			}
			return errs
		}
		for {
			select {
			case <-ctx.Done():
				return finalize()
			case err = <-dech:
			case err = <-aech:
			case <-dt.C:
				err = c.discover(ctx, ech)
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
	var ok bool
	addrs, ok = c.addrs.Load().([]string)
	if !ok {
		ips, err := net.DefaultResolver.LookupIPAddr(ctx, c.dns)
		if err != nil {
			return nil
		}
		addrs = make([]string, 0, len(ips))
		for _, ip := range ips {
			addrs = append(addrs, ip.String())
		}
	}
	return addrs
}

func (c *client) GetClient() grpc.Client {
	return c.client
}

func (c *client) connect(ctx context.Context, addr string) (err error) {
	if c.autoconn && c.client != nil {
		_, err = c.client.Connect(ctx, addr)
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

func (c *client) dnsDiscovery(ctx context.Context, ech chan<- error) (addrs []string, err error) {
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, c.dns)
	if err != nil || len(ips) == 0 {
		return nil, errors.ErrAddrCouldNotDiscover(err, c.dns)
	}
	addrs = make([]string, 0, len(ips))
	for _, ip := range ips {
		addr := net.JoinHostPort(ip.String(), uint16(c.port))
		if err = c.connect(ctx, addr); err != nil {
			ech <- err
		} else {
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

func (c *client) discover(ctx context.Context, ech chan<- error) (err error) {
	if c.dscClient == nil || (c.autoconn && c.client == nil) {
		return errors.ErrGRPCClientNotFound
	}

	var connected []string
	if bo := c.client.GetBackoff(); bo != nil {
		_, err = bo.Do(ctx, func(ctx context.Context) (interface{}, bool, error) {
			connected, err = c.updateDiscoveryInfo(ctx, ech)
			if err != nil {
				return nil, true, err
			}
			return nil, false, nil
		})
	} else {
		connected, err = c.updateDiscoveryInfo(ctx, ech)
	}
	if err != nil {
		log.Warn("failed to discover addrs from discoverer API, trying to discover from dns...\t" + err.Error())
		connected, err = c.dnsDiscovery(ctx, ech)
		if err != nil {
			return err
		}
	}

	oldAddrs := c.GetAddrs(ctx)
	c.addrs.Store(connected)
	return c.disconnectOldAddrs(ctx, oldAddrs, connected, ech)
}

func (c *client) updateDiscoveryInfo(ctx context.Context, ech chan<- error) (connected []string, err error) {
	nodes, err := c.discoverNodes(ctx)
	if err != nil {
		return nil, err
	}
	connected, err = c.discoverAddrs(ctx, nodes, ech)
	if err != nil {
		return nil, err
	}
	if len(connected) == 0 {
		log.Warn("connected addr is zero")
		return nil, errors.ErrAddrCouldNotDiscover(err, c.dns)
	}
	if c.onDiscover != nil {
		err = c.onDiscover(ctx, c, connected)
		if err != nil {
			return nil, err
		}
	}
	return connected, nil
}

func (c *client) discoverNodes(ctx context.Context) (nodes *payload.Info_Nodes, err error) {
	_, err = c.dscClient.RoundRobin(grpc.WithGRPCMethod(ctx, "discoverer.v1.Discoverer/Nodes"), func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (interface{}, error) {
		nodes, err = discoverer.NewDiscovererClient(conn).
			Nodes(ctx, &payload.Discoverer_Request{
				Namespace: c.namespace,
				Name:      c.name,
				Node:      c.nodeName,
			}, copts...)
		if err != nil {
			return nil, err
		}
		return nodes, nil
	})
	return nodes, err
}

func (c *client) discoverAddrs(ctx context.Context, nodes *payload.Info_Nodes, ech chan<- error) (addrs []string, err error) {
	maxPodLen := 0
	podLength := 0
	for _, node := range nodes.GetNodes() {
		l := len(node.GetPods().GetPods())
		podLength += l
		if l > maxPodLen {
			maxPodLen = l
		}
	}
	addrs = make([]string, 0, podLength)
	for i := 0; i < maxPodLen; i++ {
		for _, node := range nodes.GetNodes() {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				if node != nil && node.GetPods() != nil && len(node.GetPods().GetPods()) > i {
					addr := net.JoinHostPort(node.GetPods().GetPods()[i].GetIp(), uint16(c.port))
					if err = c.connect(ctx, addr); err != nil {
						err = errors.ErrAddrCouldNotDiscover(err, addr)
						select {
						case <-ctx.Done():
							return nil, ctx.Err()
						case ech <- err:
						}
						err = nil
					} else {
						addrs = append(addrs, addr)
					}
					break
				}
			}
		}
	}
	return addrs, nil
}

func (c *client) disconnectOldAddrs(ctx context.Context, oldAddrs, connectedAddrs []string, ech chan<- error) (err error) {
	if !c.autoconn {
		return nil
	}
	var cur sync.Map
	for _, addr := range connectedAddrs {
		cur.Store(addr, struct{}{})
	}

	for _, old := range oldAddrs {
		_, ok := cur.Load(old)
		if !ok {
			c.eg.Go(safety.RecoverFunc(func() error {
				err = c.disconnect(ctx, old)
				if err != nil {
					ech <- err
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
					select {
					case <-ctx.Done():
						return errors.Wrap(ctx.Err(), err.Error())
					case ech <- err:
						return err
					}
				}
			}
			return nil
		}); err != nil {
			select {
			case <-ctx.Done():
				return errors.Wrap(ctx.Err(), err.Error())
			case ech <- err:
				return err
			}
		}
	}
	return nil
}
