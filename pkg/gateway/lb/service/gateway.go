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

// Package service
package service

import (
	"context"
	"reflect"
	"sync/atomic"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	vc "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Gateway interface {
	Start(ctx context.Context) (<-chan error, error)
	GetAgentCount(ctx context.Context) int
	Addrs(ctx context.Context) []string
	DoMulti(ctx context.Context, num int,
		f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error) error
	BroadCast(ctx context.Context, kind BroadCastKind,
		f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error) error
}

type BroadCastKind int

const (
	READ BroadCastKind = iota
	WRITE
)

type gateway struct {
	client discoverer.Client
	eg     errgroup.Group
}

func NewGateway(opts ...Option) (gw Gateway, err error) {
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

func (g *gateway) BroadCast(
	ctx context.Context,
	kind BroadCastKind,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) (err error) {
	fctx, span := trace.StartSpan(ctx, "vald/gateway-lb/service/Gateway.BroadCast")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var client grpc.Client
	switch kind {
	case READ:
		client = g.client.GetReadClient()
	case WRITE:
		client = g.client.GetClient()
	}

	return client.RangeConcurrent(fctx, -1, func(ictx context.Context,
		addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (err error) {
		select {
		case <-ictx.Done():
			return nil
		default:
			err = f(ictx, addr, vc.NewValdClient(conn), copts...)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *gateway) DoMulti(
	ctx context.Context,
	num int,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) (err error) {
	sctx, span := trace.StartSpan(ctx, "vald/gateway-lb/service/Gateway.DoMulti")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var cur uint32 = 0
	addrs := g.client.GetAddrs(sctx)
	var limit uint32
	if len(addrs) < num {
		limit = uint32(len(addrs))
	} else {
		limit = uint32(num)
	}
	var visited sync.Map[string, any]
	err = g.client.GetClient().OrderedRange(sctx, addrs, func(ictx context.Context,
		addr string,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (err error) {
		if atomic.LoadUint32(&cur) < limit {
			err = f(ictx, addr, vc.NewValdClient(conn), copts...)
			if err != nil {
				return err
			}
			atomic.AddUint32(&cur, 1)
			visited.Store(addr, struct{}{})
		}
		return nil
	})
	if err != nil || cur < limit {
		err = g.client.GetClient().OrderedRange(sctx, addrs, func(ictx context.Context,
			addr string,
			conn *grpc.ClientConn,
			copts ...grpc.CallOption,
		) (err error) {
			if atomic.LoadUint32(&cur) < limit {
				_, ok := visited.Load(addr)
				if !ok {
					err = f(ictx, addr, vc.NewValdClient(conn), copts...)
					if err != nil {
						return err
					}
					atomic.AddUint32(&cur, 1)
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *gateway) GetAgentCount(ctx context.Context) int {
	return len(g.Addrs(ctx))
}

func (g *gateway) Addrs(ctx context.Context) []string {
	return g.client.GetAddrs(ctx)
}
