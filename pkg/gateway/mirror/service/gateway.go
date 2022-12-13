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
	"reflect"

	mclient "github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
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
	ForwardedContext(ctx context.Context) context.Context
	FromForwardedContext(ctx context.Context) string
	Addrs(ctx context.Context) []string
	BroadCast(ctx context.Context,
		f func(ctx context.Context, tgt string, conn *grpc.ClientConn, copts ...grpc.CallOption) error) error
}

type gateway struct {
	client  mclient.Client // Mirror Gateway client for other cluster.
	iclient mclient.Client // Mirror Gateway client for the same cluster.
	eg      errgroup.Group
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
	fctx, span := trace.StartSpan(ctx, "vald/gateway-mirror/service/Gateway.BroadCast")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return g.client.GRPCClient().RangeConcurrent(g.ForwardedContext(fctx), -1, func(ictx context.Context,
		addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (err error) {
		select {
		case <-ictx.Done():
			return nil
		default:
			err = f(ictx, addr, conn, copts...)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// func (g *gateway) Connect(ctx context.Context, addrs ...string) error {
// 	return nil
// }

func (g *gateway) Addrs(ctx context.Context) []string {
	return g.client.GRPCClient().ConnectedAddrs()
}
