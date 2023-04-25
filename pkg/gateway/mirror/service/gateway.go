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

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	client "github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

const (
	forwardedContextKey   = "forwarded-for"
	forwardedContextValue = "gateway mirror"
)

type Gateway interface {
	ForwardedContext(ctx context.Context, podName string) context.Context
	FromForwardedContext(ctx context.Context) string
	BroadCast(ctx context.Context,
		f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error
	Do(ctx context.Context, target string,
		f func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error)) (interface{}, error)
	DoMulti(ctx context.Context, targets []string,
		f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error
	GRPCClient() grpc.Client
}

type gateway struct {
	client  client.Client // Mirror Gateway client for other clusters and to the Vald gateway (LB gateway) client for own cluster.
	eg      errgroup.Group
	podName string
}

func NewGateway(opts ...Option) (Gateway, error) {
	g := new(gateway)
	for _, opt := range append(defaultGWOpts, opts...) {
		if err := opt(g); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(err, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
			return nil, oerr
		}
	}
	return g, nil
}

func (g *gateway) GRPCClient() grpc.Client {
	return g.client.GRPCClient()
}

func (g *gateway) ForwardedContext(ctx context.Context, podName string) context.Context {
	return grpc.NewOutgoingContext(ctx, grpc.MD{
		forwardedContextKey: []string{
			podName,
		},
	})
}

func (g *gateway) FromForwardedContext(ctx context.Context) string {
	md, ok := grpc.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	vals, ok := md[forwardedContextKey]
	if !ok {
		return ""
	}
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (g *gateway) BroadCast(ctx context.Context,
	f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error,
) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Gateway.BroadCast")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return g.client.GRPCClient().RangeConcurrent(g.ForwardedContext(ctx, g.podName), -1, func(ictx context.Context,
		addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (err error) {
		select {
		case <-ictx.Done():
			return nil
		default:
			err = f(ictx, addr, vald.NewValdClientWithMirror(conn), copts...)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *gateway) Do(ctx context.Context, target string,
	f func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error),
) (res interface{}, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Gateway.Do")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(target) == 0 {
		return nil, errors.ErrTargetNotFound
	}
	return g.client.GRPCClient().Do(g.ForwardedContext(ctx, g.podName), target,
		func(ictx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return f(ictx, vald.NewValdClientWithMirror(conn), copts...)
		},
	)
}

func (g *gateway) DoMulti(ctx context.Context, targets []string,
	f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error,
) error {
	ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Gateway.DoMulti")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(targets) == 0 {
		return errors.ErrTargetNotFound
	}
	return g.client.GRPCClient().OrderedRangeConcurrent(g.ForwardedContext(ctx, g.podName), targets, -1,
		func(ictx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) (err error) {
			select {
			case <-ictx.Done():
				return nil
			default:
				err = f(ictx, addr, vald.NewValdClientWithMirror(conn), copts...)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
}
