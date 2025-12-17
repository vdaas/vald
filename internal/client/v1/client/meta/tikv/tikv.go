// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
package tikv

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/tikv"
	"github.com/vdaas/vald/internal/client/v1/client/meta"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

const (
	apiName = "vald/internal/client/v1/client/meta/tikv"
)

type Client interface {
	meta.MetadataClient
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
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if c.c == nil {
		if len(c.addrs) == 0 {
			return nil, errors.ErrGRPCTargetAddrNotFound
		}
		c.c = grpc.New("TiKV Client", grpc.WithAddrs(c.addrs...))
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

func (c *client) RawGet(
	ctx context.Context, in *tikv.RawGetRequest, opts ...grpc.CallOption,
) (res *tikv.RawGetResponse, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawGet"), apiName+"/RawGet")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawGetResponse, error) {
		return tikv.NewTikvClient(conn).RawGet(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) RawBatchGet(
	ctx context.Context, in *tikv.RawBatchGetRequest, opts ...grpc.CallOption,
) (res *tikv.RawBatchGetResponse, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawBatchGet"), apiName+"/RawBatchGet")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchGetResponse, error) {
		return tikv.NewTikvClient(conn).RawBatchGet(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) RawPut(
	ctx context.Context, in *tikv.RawPutRequest, opts ...grpc.CallOption,
) (res *tikv.RawPutResponse, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawPut"), apiName+"/RawPut")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawPutResponse, error) {
		return tikv.NewTikvClient(conn).RawPut(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) RawBatchPut(
	ctx context.Context, in *tikv.RawBatchPutRequest, opts ...grpc.CallOption,
) (res *tikv.RawBatchPutResponse, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawBatchPut"), apiName+"/RawBatchPut")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchPutResponse, error) {
		return tikv.NewTikvClient(conn).RawBatchPut(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) RawDelete(
	ctx context.Context, in *tikv.RawDeleteRequest, opts ...grpc.CallOption,
) (res *tikv.RawDeleteResponse, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawDelete"), apiName+"/RawDelete")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawDeleteResponse, error) {
		return tikv.NewTikvClient(conn).RawDelete(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) RawBatchDelete(
	ctx context.Context, in *tikv.RawBatchDeleteRequest, opts ...grpc.CallOption,
) (res *tikv.RawBatchDeleteResponse, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawBatchDelete"), apiName+"/RawBatchDelete")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchDeleteResponse, error) {
		return tikv.NewTikvClient(conn).RawBatchDelete(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
