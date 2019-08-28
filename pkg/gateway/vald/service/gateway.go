//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"io"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/discoverer"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/gateway/vald/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Gateway interface {
	GetAgentCount() int
	Do(ctx context.Context,
		f func(ctx context.Context, tgt string, ac agent.AgentClient) error) error
	DoMulti(ctx context.Context, num int,
		f func(ctx context.Context, tgt string, ac agent.AgentClient) error) error
	BroadCast(ctx context.Context,
		f func(ctx context.Context, tgt string, ac agent.AgentClient) error) error
}

type gateway struct {
	agentName   string
	agentExists map[string]struct{}
	port        int
	dscDur      time.Duration
	gopts       []grpc.DialOption
	dsc         discoverer.DiscovererClient
	dscHost     string
	dscPort     int
	conns       sync.Map
	agents      atomic.Value
	bo          backoff.Backoff
	eg          errgroup.Group
}

type client struct {
	closer io.Closer
}

func New(opts ...GWOption) (gw Gateway, err error) {
	g := new(gateway)
	for _, opt := range opts {
		err = opt(g)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

func (g *gateway) StartDiscoverd(ctx context.Context) <-chan error {
	ech := make(chan error, 100)

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", g.dscHost, g.dscPort), g.gopts...)
	g.dsc = discoverer.NewDiscovererClient(conn)
	if err != nil {
		ech <- err
	}

	g.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		dt := time.NewTicker(g.dscDur)
		defer dt.Stop()
		for {
			select {
			case <-ctx.Done():
				g.conns.Range(func(key interface{}, c interface{}) bool {
					err = c.(*grpc.ClientConn).Close()
					if err != nil {
						ech <- err
					}
					g.conns.Delete(key)
					return true
				})

				return nil
			case <-dt.C:
				_, err = g.bo.Do(ctx, func() (_ interface{}, err error) {
					return g.discover(ctx, ech)
				})
				if err != nil {
					ech <- err
				}
			}
		}
	}))
	return ech
}

func (g *gateway) discover(ctx context.Context, ech chan<- error) (ret interface{}, err error) {
	res, err := g.dsc.Discover(ctx, &payload.Discoverer_Request{
		Name: g.agentName,
		Node: "",
	}, nil)
	if err != nil {
		return nil, err
	}
	srvs := res.GetServers()
	as := make(model.Agents, 0, len(srvs))
	cur := make(map[string]struct{}, len(srvs))
	for _, srv := range srvs {
		host := srv.GetServer()
		if host == nil {
			host = new(payload.Info_Server)
		}
		a := model.Agent{
			IP:       srv.GetIp(),
			Name:     srv.GetName(),
			CPU:      srv.GetCpu(),
			Mem:      srv.GetMem(),
			HostIP:   host.GetIp(),
			HostName: host.GetName(),
			HostCPU:  host.GetCpu(),
			HostMem:  host.GetMem(),
		}
		c, ok := g.conns.Load(a.Name)
		var conn *grpc.ClientConn
		if ok {
			conn = c.(*grpc.ClientConn)
			if conn.GetState() == connectivity.Shutdown ||
				conn.GetState() == connectivity.TransientFailure {
				g.conns.Delete(a.Name)
				err = conn.Close()
				if err != nil {
					ech <- err
				}
			} else {
				cur[a.Name] = struct{}{}
				as = append(as, a)
				continue
			}
		}
		conn, err = grpc.DialContext(ctx, fmt.Sprintf("%s:%d", srv.GetIp(), g.port), g.gopts...)
		if err != nil {
			ech <- err
			continue
		} else {
			cur[a.Name] = struct{}{}
			g.conns.Store(srv.GetName(), conn)
			as = append(as, a)
		}
	}
	sort.Sort(as)
	g.agents.Store(as)
	g.conns.Range(func(key interface{}, c interface{}) bool {
		_, ok := cur[key.(string)]
		if !ok {
			err = c.(*grpc.ClientConn).Close()
			if err != nil {
				ech <- err
			}
			g.conns.Delete(key)
		}
		return true
	})
	return nil, nil
}

func (g *gateway) BroadCast(ctx context.Context,
	f func(ctx context.Context, target string, ac agent.AgentClient) error) (err error) {
	return g.DoMulti(ctx, g.GetAgentCount(), f)
}

func (g *gateway) Do(ctx context.Context,
	f func(ctx context.Context, target string, ac agent.AgentClient) error) (err error) {
	_, err = g.bo.Do(ctx, func() (_ interface{}, err error) {
		var ac agent.AgentClient
		var tgt string
		for _, a := range g.agents.Load().(model.Agents) {
			c, ok := g.conns.Load(a.Name)
			if ok {
				ac = agent.NewAgentClient(c.(*grpc.ClientConn))
				tgt = a.IP
				break
			}
		}
		if ac == nil {
			return nil, errors.ErrAgentClientNotConnected
		}
		return nil, f(ctx, tgt, ac)
	})
	return nil
}

func (g *gateway) DoMulti(ctx context.Context,
	num int, f func(ctx context.Context, target string, ac agent.AgentClient) error) error {
	eg, ctx := errgroup.New(ctx)
	agents := g.agents.Load().(model.Agents)
	if num > agents.Len() {
		num = agents.Len()
	}
	for _, a := range agents[:num] {
		var ac agent.AgentClient
		var tgt string
		c, ok := g.conns.Load(a.Name)
		if ok {
			tgt = a.IP
			ac = agent.NewAgentClient(c.(*grpc.ClientConn))
		}
		eg.Go(safety.RecoverFunc(func() (err error) {
			if ac == nil {
				return errors.ErrAgentClientNotConnected
			}
			_, err = g.bo.Do(ctx, func() (_ interface{}, err error) {
				err = f(ctx, tgt, ac)
				return
			})
			return
		}))
	}
	return eg.Wait()
}

func (g *gateway) GetAgentCount() int {
	return g.agents.Load().(model.Agents).Len()
}
