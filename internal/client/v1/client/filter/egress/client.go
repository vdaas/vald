//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package egress provides egress filter client logic
package egress

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/apis/grpc/v1/filter/egress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/strings"
	valdsync "github.com/vdaas/vald/internal/sync"
)

type Client interface {
	egress.FilterClient
	Target(ctx context.Context, targets ...string) (egress.FilterClient, error)
	GRPCClient() grpc.Client
	Start(context.Context) (<-chan error, error)
	Stop(context.Context) error
}

type client struct {
	addrs []string
	cl    valdsync.Map[string, any]
	c     grpc.Client
}

type specificAddrClient struct {
	addr string
	c    grpc.Client
}

type multipleAddrsClient struct {
	addrs []string
	c     grpc.Client
}

const (
	apiName = "vald/internal/client/v1/client/filter/egress"
)

func New(opts ...Option) (Client, error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		err := opt(c)
		if err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	if c.c == nil {
		if c.addrs == nil {
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
	return c.c.Close(ctx)
}

func (c *client) GRPCClient() grpc.Client {
	return c.c
}

func (c *client) Target(ctx context.Context, targets ...string) (egress.FilterClient, error) {
	if len(targets) == 0 {
		return nil, errors.ErrTargetNotFound
	}
	if len(targets) == 1 {
		addr := targets[0]
		_, ok := c.cl.Load(addr)
		if !ok || !c.c.IsConnected(ctx, addr) {
			_, err := c.c.Connect(ctx, addr)
			if err != nil {
				if ok {
					c.cl.Delete(addr)
				}
				return nil, err
			}
			c.cl.Store(addr, struct{}{})
		}
		return &specificAddrClient{
			addr: addr,
			c:    c.c,
		}, nil
	}
	addrs := make([]string, 0, len(targets))
	for _, addr := range targets {
		_, ok := c.cl.Load(addr)
		if !ok || !c.c.IsConnected(ctx, addr) {
			_, err := c.c.Connect(ctx, addr)
			if err != nil {
				if ok {
					c.cl.Delete(addr)
				}
				return nil, err
			}
			c.cl.Store(addr, struct{}{})
		}
		addrs = append(addrs, addr)
	}
	return &multipleAddrsClient{
		addrs: addrs,
		c:     c.c,
	}, nil
}

func (c *client) FilterDistance(ctx context.Context, in *payload.Object_Distance, opts ...grpc.CallOption) (res *payload.Object_Distance, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.FilterDistance")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = egress.NewFilterClient(conn).FilterDistance(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *specificAddrClient) FilterDistance(ctx context.Context, in *payload.Object_Distance, opts ...grpc.CallOption) (res *payload.Object_Distance, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.FilterDistance/"+s.addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = s.c.Do(ctx, s.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = egress.NewFilterClient(conn).FilterDistance(ctx, in, append(copts, opts...)...)
		if err != nil {
			return nil, err
		}
		in = res
		return in, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *multipleAddrsClient) FilterDistance(ctx context.Context, in *payload.Object_Distance, opts ...grpc.CallOption) (res *payload.Object_Distance, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.FilterDistance/["+strings.Join(m.addrs, ",")+"]")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = m.c.OrderedRange(ctx, m.addrs, func(ctx context.Context, addr string,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) error {
		res, err = egress.NewFilterClient(conn).FilterDistance(ctx, in, append(copts, opts...)...)
		if err != nil {
			return err
		}
		in = res
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.FilterVector")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = egress.NewFilterClient(conn).FilterVector(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *specificAddrClient) FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.FilterVector/"+s.addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = s.c.Do(ctx, s.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = egress.NewFilterClient(conn).FilterVector(ctx, in, append(copts, opts...)...)
		if err != nil {
			return nil, err
		}
		in = res
		return in, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *multipleAddrsClient) FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.FilterVector/["+strings.Join(m.addrs, ",")+"]")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	err = m.c.OrderedRange(ctx, m.addrs, func(ctx context.Context, addr string,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) error {
		res, err = egress.NewFilterClient(conn).FilterVector(ctx, in, append(copts, opts...)...)
		if err != nil {
			return err
		}
		in = res
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
