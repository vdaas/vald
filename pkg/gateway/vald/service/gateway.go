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
	"runtime"
	"sort"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/discoverer"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/pkg/gateway/vald/model"
)

type Gateway interface {
	Start(ctx context.Context) <-chan error
	GetAgentCount() int
	Do(ctx context.Context,
		f func(ctx context.Context, tgt string, ac agent.AgentClient) error) error
	DoMulti(ctx context.Context, num int,
		f func(ctx context.Context, tgt string, ac agent.AgentClient) error) error
	BroadCast(ctx context.Context,
		f func(ctx context.Context, tgt string, ac agent.AgentClient) error) error
}

type gateway struct {
	agentName  string
	agentPort  int
	agentAddrs []string
	agentHcDur string
	agents     atomic.Value
	dscAddr    string
	dscDur     time.Duration
	dscClient  grpc.Client
	acClient   grpc.Client
	gopts      []grpc.DialOption
	copts      []grpc.CallOption
	bo         backoff.Backoff
	eg         errgroup.Group
}

type client struct {
	closer io.Closer
}

func NewGateway(opts ...GWOption) (gw Gateway, err error) {
	g := new(gateway)
	for _, opt := range opts {
		err = opt(g)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

func (g *gateway) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 3)
	dech := g.dscClient.StartConnectionMonitor(ctx)

	var err error
	_, err = g.discover(ctx, ech)
	if err != nil {
		ech <- err
	}

	g.acClient, err = grpc.New(
		grpc.WithAddrs(g.agentAddrs...),
		grpc.WithErrGroup(g.eg),
		grpc.WithBackoff(g.bo),
		grpc.WithGRPCCallOptions(g.copts...),
		grpc.WithGRPCDialOptions(g.gopts...),
		grpc.WithHealthCheckDuration(g.agentHcDur),
	)
	if err != nil {
		ech <- err
	}

	aech := g.acClient.StartConnectionMonitor(ctx)

	g.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		fch := make(chan struct{}, 1)
		defer close(fch)
		dt := time.NewTicker(g.dscDur)
		defer dt.Stop()
		for {
			select {
			case <-ctx.Done():
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
			case ech <- <-dech:
			case ech <- <-aech:
			case <-fch:
				_, err = g.discover(ctx, ech)
				if err != nil {
					ech <- err
				}
			case <-dt.C:
				_, err = g.discover(ctx, ech)
				if err != nil {
					log.Error(err)
					time.Sleep(g.dscDur / 5)
					fch <- struct{}{}
				}
			}
		}
	}))
	return ech
}

func (g *gateway) discover(ctx context.Context, ech chan<- error) (ret interface{}, err error) {
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

	g.agentAddrs = g.agentAddrs[:0]

	srvs := res.GetServers()
	as := make(model.Agents, 0, len(srvs))
	cur := make(map[string]struct{}, len(srvs))
	for _, srv := range srvs {
		host := srv.GetServer()
		if host == nil {
			host = new(payload.Info_Server)
		}
		addr := fmt.Sprintf("%s:%d", srv.GetIp(), g.agentPort)
		g.agentAddrs = append(g.agentAddrs, addr)
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
	return nil, nil
}

func (g *gateway) BroadCast(ctx context.Context,
	f func(ctx context.Context, target string, ac agent.AgentClient) error) (err error) {
	return g.DoMulti(ctx, g.GetAgentCount(), f)
}

func (g *gateway) Do(ctx context.Context,
	f func(ctx context.Context, target string, ac agent.AgentClient) error) (err error) {
	addr := fmt.Sprintf("%s:%d", g.agents.Load().(model.Agents)[0].IP, g.agentPort)
	_, err = g.acClient.Do(ctx, addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		if conn == nil {
			return nil, errors.ErrAgentClientNotConnected
		}
		return nil, f(ctx, addr, agent.NewAgentClient(conn))
	})
	return err
}

func (g *gateway) DoMulti(ctx context.Context,
	num int, f func(ctx context.Context, target string, ac agent.AgentClient) error) error {
	var cur uint32
	limit := uint32(num)
	return g.acClient.RangeConcurrent(ctx, 0, func(addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		if conn == nil {
			return errors.ErrAgentClientNotConnected
		}
		if atomic.AddUint32(&cur, 1) > limit {
			return nil
		}
		return f(ctx, addr, agent.NewAgentClient(conn))
	})
}

func (g *gateway) GetAgentCount() int {
	return g.agents.Load().(model.Agents).Len()
}
