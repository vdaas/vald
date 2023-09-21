// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package mirror

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

const (
	apiName = "vald/internal/client/v1/client/mirror"
)

type Client interface {
	mirror.MirrorClient
	GRPCClient() grpc.Client
	Start(context.Context) (<-chan error, error)
	Stop(context.Context) error
}

type client struct {
	addrs []string
	c     grpc.Client
}

func New(opts ...Option) (Client, error) {
	c := new(client)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if c.c == nil {
		if len(c.addrs) == 0 {
			return nil, errors.ErrGRPCTargetAddrNotFound
		}
		c.c = grpc.New(grpc.WithAddrs(c.addrs...))
	}
	return c, nil
}

func (c *client) Start(ctx context.Context) (<-chan error, error) {
	return c.c.StartConnectionMonitor(ctx)
}

func (c *client) Stop(ctx context.Context) error {
	return c.Stop(ctx)
}

func (c *client) GRPCClient() grpc.Client {
	return c.c
}

func (c *client) Register(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (res *payload.Mirror_Targets, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.RegisterRPCName), apiName+"/"+vald.RegisterRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = c.c.RoundRobin(ctx, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		res, err = mirror.NewMirrorClient(conn).Register(ctx, in, append(copts, opts...)...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Advertise(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (res *payload.Mirror_Targets, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.AdvertiseRPCName), apiName+"/"+vald.AdvertiseRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = c.c.RoundRobin(ctx, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		res, err = mirror.NewMirrorClient(conn).Advertise(ctx, in, append(copts, opts...)...)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
