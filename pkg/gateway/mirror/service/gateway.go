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

// Package service
package service

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	mclient "github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

type contextKey string

const (
	forwardedContextKey   contextKey = "forwarded-for"
	forwardedContextValue            = "gateway mirror"
)

type Gateway interface {
	Start(ctx context.Context) (<-chan error, error)
	ConnectedTargets() ([]*payload.Mirror_Target, error)
	Connect(ctx context.Context, targets ...*payload.Mirror_Target) ([]*payload.Mirror_Target, error)
	ForwardedContext(ctx context.Context) context.Context
	FromForwardedContext(ctx context.Context) string
	BroadCast(ctx context.Context,
		f func(ctx context.Context, tgt string, conn *grpc.ClientConn, copts ...grpc.CallOption) error) error
}

type gateway struct {
	client       mclient.Client // Mirror Gateway client for other cluster.
	iclient      mclient.Client // Mirror Gateway client for the same cluster.
	eg           errgroup.Group
	advertiseDur time.Duration
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
	ech := make(chan error, 100)

	cech, err := g.client.Start(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}
	icech, err := g.iclient.Start(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}
	aech := g.startAdvertise(ctx)

	g.eg.Go(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return
			case err = <-cech:
			case err = <-icech:
			case err = <-aech:
			}
			if err != nil {
				select {
				case <-ctx.Done():
				case ech <- err:
				}
				err = nil
			}
		}
	})
	return ech, nil
}

func (g *gateway) startAdvertise(ctx context.Context) <-chan error {
	tic := time.NewTicker(g.advertiseDur)
	defer tic.Stop()

	ech := make(chan error, 1)
	g.eg.Go(func() error {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-tic.C:
				targets, err := g.ConnectedTargets()
				if err != nil {
					ech <- err
					continue
				}
				req := &payload.Mirror_Targets{
					Targets: targets,
				}
				err = g.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
					_, err := mirror.NewMirrorClient(conn).Advertise(ctx, req, copts...)
					return err
				})
				if err != nil {

				}
			}
		}
	})
	return ech
}

func (g *gateway) ForwardedContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, forwardedContextKey, forwardedContextValue)
}

func (g *gateway) FromForwardedContext(ctx context.Context) string {
	if v := ctx.Value(forwardedContextKey); v != nil {
		if name, ok := v.(string); ok {
			return name
		}
	}
	return ""
}

func (g *gateway) BroadCast(ctx context.Context,
	f func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) (err error) {
	fctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Gateway.BroadCast")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	imAddrs := g.internlMirrorAddrs()
	return g.client.GRPCClient().RangeConcurrent(g.ForwardedContext(fctx), -1, func(ictx context.Context,
		addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (err error) {
		select {
		case <-ictx.Done():
			return nil
		default:
			if _, ok := imAddrs[addr]; ok {
				return
			}
			err = f(ictx, addr, conn, copts...)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// internlAddrs returns the addresses of Mirror Gateway on the same cluster.
func (g *gateway) internlMirrorAddrs() (m map[string]struct{}) {
	for _, addr := range g.iclient.GRPCClient().ConnectedAddrs() {
		m[addr] = struct{}{}
	}
	return m
}

func (g *gateway) ConnectedTargets() ([]*payload.Mirror_Target, error) {
	addrs := g.client.GRPCClient().ConnectedAddrs()
	targets := make([]*payload.Mirror_Target, 0, len(addrs))

	for _, addr := range addrs {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		targets = append(targets, &payload.Mirror_Target{
			Ip:   host,
			Port: uint32(port),
		})
	}
	return targets, nil
}

func (g *gateway) Connect(ctx context.Context, targets ...*payload.Mirror_Target) ([]*payload.Mirror_Target, error) {
	imAddrs := g.internlMirrorAddrs()
	for _, target := range targets {
		addr := fmt.Sprintf("%s:%d", target.GetIp(), target.GetPort())
		if _, ok := imAddrs[addr]; ok {
			continue
		}

		if !g.client.GRPCClient().IsConnected(ctx, addr) {
			_, err := g.client.GRPCClient().Connect(ctx, addr, g.client.GRPCClient().GetDialOption()...)
			if err != nil {
				return nil, err
			}
		}
	}
	tgts, err := g.ConnectedTargets()
	if err != nil {
		return nil, err
	}
	return tgts, nil
}
