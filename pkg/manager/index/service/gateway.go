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
	"net"
	"reflect"
	"runtime"
	"sort"
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
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/manager/index/model"
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
	agentPort    int
	agentARecord string
	agents       atomic.Value
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
	ech := make(chan error, 10)
	dech, err := g.dscClient.StartConnectionMonitor(ctx)
	discover := g.discover
	if err != nil {
		g.dscClient.Close()
		discover = g.discoverByDNS
	}
	_, err = discover(ctx, ech)
	if err != nil {
		g.dscClient.Close()
		discover = g.discoverByDNS
		_, err = discover(ctx, ech)
		if err != nil {
			log.Error(err)
			close(ech)
			return nil, err
		}
	}

	as := g.agents.Load().(model.Agents)
	addrs := make([]string, 0, len(as))
	for _, a := range as {
		addrs = append(addrs,
			fmt.Sprintf("%s:%d", a.IP, g.agentPort),
		)
	}

	g.acClient = grpc.New(
		append(
			g.agentOpts,
			grpc.WithAddrs(addrs...),
			grpc.WithErrGroup(g.eg),
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
				_, err = discover(ctx, ech)
				if err != nil {
					ech <- err
					err = nil
				}
			case <-dt.C:
				_, err = discover(ctx, ech)
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

func (g *gateway) discoverByDNS(ctx context.Context, ech chan<- error) (ret interface{}, err error) {
	log.Info("starting dns discovery")
	ips, err := net.LookupIP(g.agentARecord)
	if err != nil {
		ech <- err
		return nil, err
	}

	if len(ips) == 0 {
		return nil, errors.ErrAgentAddrCouldNotDiscover
	}

	as := make(model.Agents, 0, len(ips))
	cur := make(map[string]struct{}, len(ips))
	for _, ip := range ips {
		host, err := net.LookupAddr(ip.String())
		if err != nil {
			return nil, err
		}
		as = append(as, model.Agent{
			IP:   ip.String(),
			Name: host[0],
		})
		cur[fmt.Sprintf("%s:%d", ip.String(), g.agentPort)] = struct{}{}
	}

	g.agents.Store(as)

	if g.acClient != nil {
		err = g.acClient.Range(ctx,
			func(addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
				_, ok := cur[addr]
				delete(cur, addr)
				if !ok {
					return g.acClient.Disconnect(addr)
				}
				return nil
			})
		if err != nil {
			ech <- err
			err = nil
		}

		for addr := range cur {
			err = g.acClient.Connect(ctx, addr)
			if err != nil {
				ech <- err
				err = nil
			}
		}
	}
	log.Info("finished dns discovery")
	return nil, nil
}

func (g *gateway) discover(ctx context.Context, ech chan<- error) (ret interface{}, err error) {
	log.Info("starting discoverer discovery")
	var res *payload.Info_Servers
	_, err = g.dscClient.Do(ctx, g.dscAddr,
		func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			res, err = discoverer.NewDiscovererClient(conn).
				Discover(ctx, &payload.Discoverer_Request{
					Name: g.agentName,

					Node: "",
				}, copts...)
			if err != nil {
				return nil, err
			}
			return res, nil
		})
	if err != nil {
		runtime.Gosched()
		return nil, err
	}

	srvs := res.GetServers()
	if len(srvs) == 0 {
		return nil, errors.ErrAgentAddrCouldNotDiscover
	}
	as := make(model.Agents, 0, len(srvs))
	cur := make(map[string]struct{}, len(srvs))
	for _, srv := range srvs {
		host := srv.GetServer()
		if host == nil {
			host = new(payload.Info_Server)
		}
		addr := fmt.Sprintf("%s:%d", srv.GetIp(), g.agentPort)
		as = append(as, model.Agent{
			IP:       srv.GetIp(),
			Name:     srv.GetName(),
			CPU:      srv.GetCpu(),
			Mem:      srv.GetMem(),
			HostIP:   host.GetIp(),
			HostName: host.GetName(),
			HostCPU:  host.GetCpu(),
			HostMem:  host.GetMem(),
		})
		cur[addr] = struct{}{}
	}

	sort.Sort(as)
	g.agents.Store(as)

	if g.acClient != nil {
		err = g.acClient.Range(ctx,
			func(addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
				_, ok := cur[addr]
				delete(cur, addr)
				if !ok {
					return g.acClient.Disconnect(addr)
				}
				return nil
			})
		if err != nil {
			ech <- err
			err = nil
		}

		for addr := range cur {
			err = g.acClient.Connect(ctx, addr)
			if err != nil {
				ech <- err
				err = nil
			}
		}
	}
	log.Info("starting discoverer discovery")
	return nil, nil
}

func (g *gateway) BroadCast(ctx context.Context,
	f func(ctx context.Context, target string, ac agent.AgentClient, copts ...grpc.CallOption) error) (err error) {
	return g.DoMulti(ctx, g.GetAgentCount(), f)
}

func (g *gateway) Do(ctx context.Context,
	f func(ctx context.Context, target string, ac agent.AgentClient, copts ...grpc.CallOption) error) (err error) {
	addr := fmt.Sprintf("%s:%d", g.agents.Load().(model.Agents)[0].IP, g.agentPort)
	_, err = g.acClient.Do(ctx, addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		if conn == nil {
			return nil, errors.ErrAgentClientNotConnected
		}
		return nil, f(ctx, addr, agent.NewAgentClient(conn), copts...)
	})
	return err
}

func (g *gateway) DoMulti(ctx context.Context,
	num int, f func(ctx context.Context, target string, ac agent.AgentClient, copts ...grpc.CallOption) error) error {
	var cur uint32
	limit := uint32(num)
	ctx, cancel := context.WithCancel(ctx)
	var once sync.Once
	return g.acClient.RangeConcurrent(ctx, int(limit), func(addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) (err error) {
		if conn == nil {
			return errors.ErrAgentClientNotConnected
		}
		if atomic.LoadUint32(&cur) > limit {
			once.Do(func() {
				cancel()
			})
			return nil
		}
		err = f(ctx, addr, agent.NewAgentClient(conn), copts...)
		if err != nil {
			return err
		}
		atomic.AddUint32(&cur, 1)
		return nil
	})
}

func (g *gateway) GetAgentCount() int {
	return g.agents.Load().(model.Agents).Len()
}
