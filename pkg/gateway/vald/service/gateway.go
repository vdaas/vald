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

// Package service
package service

import (
	"context"
	"fmt"
	"math"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/discoverer"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/metric"
	"github.com/vdaas/vald/internal/safety"
)

type Gateway interface {
	Start(ctx context.Context) (<-chan error, error)
	GetAgentCount() int
	Do(ctx context.Context,
		f func(ctx context.Context, tgt string, ac agent.AgentClient, copts ...grpc.CallOption) error) error
	DoMulti(ctx context.Context, num int,
		f func(ctx context.Context, tgt string, ac agent.AgentClient, copts ...grpc.CallOption) error) error
	BroadCast(ctx context.Context,
		f func(ctx context.Context, tgt string, ac agent.AgentClient, copts ...grpc.CallOption) error) error
}

type gateway struct {
	agentName    string
	namespace    string
	nodeName     string
	agentPort    int
	agentARecord string
	agents       atomic.Value // []string ips
	dscAddr      string
	dscDur       time.Duration
	dscClient    grpc.Client
	acClient     grpc.Client
	agentOpts    []grpc.Option
	eg           errgroup.Group
}

func NewGateway(opts ...GWOption) (gw Gateway, err error) {
	g := new(gateway)
	for _, opt := range append(defaultGWOpts, opts...) {
		if err := opt(g); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return g, nil
}

func (g *gateway) Start(ctx context.Context) (<-chan error, error) {
	dech, err := g.dscClient.StartConnectionMonitor(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, 100)
	err = g.discover(ctx, ech)
	if err != nil {
		close(ech)
		g.dscClient.Close()
		return nil, err
	}
	g.acClient = grpc.New(
		append(
			g.agentOpts,
			grpc.WithAddrs(g.agents.Load().([]string)...),
			grpc.WithErrGroup(g.eg),
			grpc.WithDialOptions(
				metric.WithStatsHandler(metric.NewClientHandler()),
			),
		)...,
	)

	aech, err := g.acClient.StartConnectionMonitor(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	g.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		fch := make(chan struct{}, 1)
		defer close(fch)
		dt := time.NewTicker(g.dscDur)
		defer dt.Stop()
		finalize := func() (err error) {
			var errs error
			err = g.dscClient.Close()
			if err != nil {
				errs = errors.Wrap(errs, err.Error())
			}
			err = g.acClient.Close()
			if err != nil {
				errs = errors.Wrap(errs, err.Error())
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
				err = g.discover(ctx, ech)
				if err != nil {
					ech <- err
					err = nil
				}
			case <-dt.C:
				err = g.discover(ctx, ech)
				if err != nil {
					ech <- err
					log.Error(err)
					err = nil
					time.Sleep(g.dscDur / 5)
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
				log.Debug(g.acClient.GetAddrs())
			}
		}
	}))
	return ech, nil
}

func (g *gateway) discover(ctx context.Context, ech chan<- error) (err error) {
	log.Info("starting discoverer discovery")
	addrs := make([]string, 0, 100)
	_, err = g.dscClient.Do(ctx, g.dscAddr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		nodes, err := discoverer.NewDiscovererClient(conn).
			Nodes(ctx, &payload.Discoverer_Request{
				Namespace: g.namespace,
				Name:      g.agentName,
				Node:      g.nodeName,
			}, copts...)
		if err != nil {
			return nil, err
		}
		for i := 0; i < (math.MaxInt32); i++ {
			visited := false
			for i, node := range nodes.GetNodes() {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
					if node != nil && node.GetPods() != nil {
						pods := node.GetPods().GetPods()
						if i < len(pods) {
							addrs = append(addrs, fmt.Sprintf("%s:%d", pods[i].GetIp(), g.agentPort))
							if !visited {
								visited = true
							}
							break
						}
					}
				}
			}
			if !visited {
				return nil, nil
			}
		}
		return nil, nil
	})
	if err != nil {
		log.Warn("failed to discover agents from discoverer API, trying to discover from dns...")
		ips, err := net.DefaultResolver.LookupIPAddr(ctx, g.agentARecord)
		if err != nil {
			return errors.ErrAgentAddrCouldNotDiscover(err, g.agentARecord)
		}
		addrs = make([]string, 0, len(ips))
		for _, ip := range ips {
			addrs = append(addrs, fmt.Sprintf("%s:%d", ip.String(), g.agentPort))
		}
	}
	if g.acClient != nil {
		cur := make(map[string]struct{}, len(addrs))
		connected := make([]string, 0, len(addrs))
		for _, addr := range addrs {
			err = g.acClient.Connect(ctx, addr)
			if err != nil {
				ech <- err
				err = nil
			} else {
				cur[addr] = struct{}{}
				connected = append(connected, addr)
			}
		}
		g.agents.Store(connected)
		err = g.acClient.Range(ctx,
			func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
				_, ok := cur[addr]
				if !ok {
					return g.acClient.Disconnect(addr)
				}
				delete(cur, addr)
				return nil
			})
		if err != nil {
			ech <- err
			err = nil
		}
	}
	log.Info("finished discoverer discovery")
	return nil
}

func (g *gateway) BroadCast(ctx context.Context,
	f func(ctx context.Context, target string, ac agent.AgentClient, copts ...grpc.CallOption) error) (err error) {
	return g.acClient.RangeConcurrent(ctx, -1, func(ctx context.Context,
		addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) (err error) {
		select {
		case <-ctx.Done():
			return nil
		default:
			err = f(ctx, addr, agent.NewAgentClient(conn), copts...)
			if err != nil {
				log.Debug(addr, err)
				return err
			}
		}
		return nil
	})
}

func (g *gateway) Do(ctx context.Context,
	f func(ctx context.Context, target string, ac agent.AgentClient, copts ...grpc.CallOption) error) (err error) {
	addr := g.agents.Load().([]string)[0]
	_, err = g.acClient.Do(ctx, addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		return nil, f(ctx, addr, agent.NewAgentClient(conn), copts...)
	})
	return err
}

func (g *gateway) DoMulti(ctx context.Context,
	num int, f func(ctx context.Context, target string, ac agent.AgentClient, copts ...grpc.CallOption) error) (err error) {
	var cur uint32 = 0
	limit := uint32(num)
	cctx, cancel := context.WithCancel(ctx)
	var once sync.Once
	err = g.acClient.OrderedRangeConcurrent(cctx, g.agents.Load().([]string), num,
		func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) (err error) {
			select {
			case <-cctx.Done():
				return nil
			default:
				if atomic.LoadUint32(&cur) >= limit {
					once.Do(func() {
						cancel()
					})
					return nil
				}
				err = f(cctx, addr, agent.NewAgentClient(conn), copts...)
				if err != nil {
					log.Debug(addr, err)
					return err
				}
				atomic.AddUint32(&cur, 1)
			}
			return nil
		})
	if err != nil && cur < limit {
		return err
	}
	return nil
}

func (g *gateway) GetAgentCount() int {
	return len(g.agents.Load().([]string))
}
