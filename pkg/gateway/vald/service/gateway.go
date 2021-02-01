//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"reflect"
	"sync/atomic"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
)

type Gateway interface {
	Start(ctx context.Context) (<-chan error, error)
	GetAgentCount(ctx context.Context) int
	Do(ctx context.Context,
		f func(ctx context.Context, tgt string, vc vald.Client, copts ...grpc.CallOption) error) error
	DoMulti(ctx context.Context, num int,
		f func(ctx context.Context, tgt string, vc vald.Client, copts ...grpc.CallOption) error) error
	BroadCast(ctx context.Context,
		f func(ctx context.Context, tgt string, vc vald.Client, copts ...grpc.CallOption) error) error
}

type gateway struct {
	client discoverer.Client
	eg     errgroup.Group
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
	return g.client.Start(ctx)
}

func (g *gateway) BroadCast(ctx context.Context,
	f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) (err error) {
	return g.client.GetClient().RangeConcurrent(ctx, -1, func(ctx context.Context,
		addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) (err error) {
		select {
		case <-ctx.Done():
			return nil
		default:
			err = f(ctx, addr, vald.NewValdClient(conn), copts...)
			if err != nil {
				log.Debugf("an error occurred while calling RPC of %s: %s", addr, err)
				return err
			}
		}
		return nil
	})
}

func (g *gateway) Do(ctx context.Context,
	f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) (err error) {
	addr := g.client.GetAddrs(ctx)[0]
	_, err = g.client.GetClient().Do(ctx, addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		return nil, f(ctx, addr, vald.NewValdClient(conn), copts...)
	})
	return err
}

func (g *gateway) DoMulti(ctx context.Context, num int,
	f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) (err error) {
	var cur uint32 = 0
	limit := uint32(num)
	addrs := g.client.GetAddrs(ctx)
	log.Debug("executing DoMulti for addrs:", addrs)
	err = g.client.GetClient().OrderedRange(ctx, addrs, func(ictx context.Context,
		addr string,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (err error) {
		if atomic.LoadUint32(&cur) < limit {
			err = f(ictx, addr, vald.NewValdClient(conn), copts...)
			if err != nil {
				log.Debugf("an error occurred while calling RPC of %s: %s", addr, err)
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

func (g *gateway) GetAgentCount(ctx context.Context) int {
	return len(g.client.GetAddrs(ctx))
}
