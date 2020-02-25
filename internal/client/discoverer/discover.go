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

// Package discoverer
package discoverer

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/discoverer"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
)

type Client interface {
	Start(ctx context.Context) (<-chan error, error)
	GetAddrs() []string
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
	dscAddr      string
	dscClient    grpc.Client
	dscDur       time.Duration
	eg           errgroup.Group
	name         string
	namespace    string
	nodeName     string
}

func New(opts ...Option) (d Client, err error) {
	c := new(client)
	for _, opt := range append(defaultOpts, opts...) {
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

	if c.autoconn {
		c.client = grpc.New(
			append(
				c.opts,
				grpc.WithAddrs(addrs...),
				grpc.WithErrGroup(c.eg),
			)...,
		)
	}

	err = c.discover(ctx, ech)
	if err != nil {
		close(ech)
		return nil, errors.Wrap(c.dscClient.Close(), err.Error())
	}

	var aech <-chan error
	if c.autoconn && c.client != nil {
		aech, err = c.client.StartConnectionMonitor(ctx)
		if err != nil {
			close(ech)
			return nil, err
		}
	}

	c.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		fch := make(chan struct{}, 1)
		defer close(fch)
		dt := time.NewTicker(c.dscDur)
		defer dt.Stop()
		finalize := func() (err error) {
			var errs error
			err = c.dscClient.Close()
			if err != nil {
				errs = errors.Wrap(errs, err.Error())
			}
			if c.autoconn && c.client != nil {
				err = c.client.Close()
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
			case <-fch:
				err = c.discover(ctx, ech)
				if err != nil {
					ech <- err
					err = nil
				}
			case <-dt.C:
				err = c.discover(ctx, ech)
				if err != nil {
					ech <- err
					log.Error(err)
					err = nil
					time.Sleep(c.dscDur / 5)
					fch <- struct{}{}
				}
			}
			if err != nil {
				log.Error(err)
				select {
				case <-ctx.Done():
					return finalize()
				case ech <- err:
				}
			} else {
				log.Debug(c.GetAddrs())
			}
		}
	}))
	return ech, nil
}

func (c *client) GetAddrs() (addrs []string) {
	var ok bool
	addrs, ok = c.addrs.Load().([]string)
	if !ok {
		return nil
	}
	return addrs
}

func (c *client) GetClient() grpc.Client {
	return c.client
}

func (c *client) connect(ctx context.Context, addr string) (err error) {
	if c.autoconn && c.client != nil {
		err = c.client.Connect(ctx, addr)
		if err != nil {
			return err
		}
		if c.onConnect != nil {
			err = c.onConnect(ctx, c, addr)
		}
	}
	return
}

func (c *client) dnsDiscovery(ctx context.Context, ech chan<- error) (addrs []string, err error) {
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, c.dns)
	if err != nil {
		return nil, errors.ErrAddrCouldNotDiscover(err, c.dns)
	}
	addrs = make([]string, 0, len(ips))
	for _, ip := range ips {
		addr := fmt.Sprintf("%s:%d", ip.String(), c.port)
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
	log.Info("starting discoverer discovery")
	connected := make([]string, 0, len(c.GetAddrs()))
	var cur sync.Map
	if _, err = c.dscClient.Do(ctx, c.dscAddr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		nodes, err := discoverer.NewDiscovererClient(conn).
			Nodes(ctx, &payload.Discoverer_Request{
				Namespace: c.namespace,
				Name:      c.name,
				Node:      c.nodeName,
			}, copts...)
		if err != nil {
            log.Warn("failed to call discoverer.Node API")
			return nil, errors.ErrRPCCallFailed(c.dscAddr, err)
		}
		var wg sync.WaitGroup
		cond := sync.NewCond(new(sync.Mutex))
		cctx, cancel := context.WithCancel(ctx)
		pch := make(chan string, len(nodes.GetNodes()))
		for _, n := range nodes.GetNodes() {
			select {
			case <-cctx.Done():
				return nil, cctx.Err()
			default:
				if n != nil && n.GetPods() != nil && n.GetPods().GetPods() != nil {
					node := n
					wg.Add(1)
					c.eg.Go(safety.RecoverFunc(func() (err error) {
						defer wg.Done()
						cond.L.Lock()
						cond.Wait()
						cond.L.Unlock()
						for _, pod := range node.GetPods().GetPods() {
							select {
							case <-cctx.Done():
								return nil
							default:
								if pod != nil && pod.GetIp() != "" {
									addr := fmt.Sprintf("%s:%d", pod.GetIp(), c.port)
									if err = c.connect(ctx, addr); err != nil {
										err = errors.ErrAddrCouldNotDiscover(err, addr)
										log.Debug(err)
										ech <- err
										err = nil
									} else {
										if c.autoconn {
											cur.Store(addr, struct{}{})
										}
										pch <- addr
									}
								}
							}
						}
						return nil
					}))
				}
			}
		}
		c.eg.Go(safety.RecoverFunc(func() error {
			wg.Wait()
			cancel()
			return nil
		}))
		cond.Broadcast()
		for {
			select {
			case <-cctx.Done():
				if len(connected) == 0 {
                    log.Warn("connected addr is zero")
					return nil, errors.ErrAddrCouldNotDiscover(err, c.dns)
				}
				if c.onDiscover != nil {
					return nil, c.onDiscover(ctx, c, connected)
				}
				return nil, nil
			case addr := <-pch:
				connected = append(connected, addr)
			}
		}
	}); err != nil {
		log.Warnf("failed to discover addrs from discoverer API, trying to discover from dns..., %v", err)
		connected, err = c.dnsDiscovery(ctx, ech)
		if err != nil {
			return err
		}
		if c.autoconn {
			cur = sync.Map{}
			for _, addr := range connected {
				cur.Store(addr, struct{}{})
			}
		}
	}

	c.addrs.Store(connected)
	if c.autoconn && c.client != nil {
		if err = c.client.RangeConcurrent(ctx, len(connected)/3, func(ctx context.Context,
			addr string,
			conn *grpc.ClientConn,
			copts ...grpc.CallOption) (err error) {
			_, ok := cur.Load(addr)
			if !ok {
				err = c.client.Disconnect(addr)
				if err != nil {
					ech <- err
				} else {
					if c.onDisconnect != nil {
						err = c.onDisconnect(ctx, c, addr)
					}
				}
				return err
			}
			return nil
		}); err != nil {
			ech <- err
			return err
		}
	}

	log.Info("finished discoverer discovery")
	return nil
}
