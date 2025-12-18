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

// $ tiup playground --mode tikv-slim
// $ TIKV_STORE_ADDRS=127.0.0.1:20160 go test ./internal/client/v1/client/meta/tikv/ -run=^$ -bench=. -benchmem
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

var errNotFound = errors.New("tikv: key not found")

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

func (c *client) Get(
	ctx context.Context, key []byte,
) (val []byte, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawGet"), apiName+"/RawGet")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err := grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawGetResponse, error) {
		return tikv.NewTikvClient(conn).RawGet(ctx, &tikv.RawGetRequest{
			Key: key,
		}, copts...)
	})
	if err != nil {
		return nil, err
	}
	// TODO
	// if res.RegionError != nil {
	if res.Error != "" {
		return nil, errors.New(res.Error)
	}
	if res.NotFound {
		return nil, errNotFound
	}
	return res.Value, nil
}

func (c *client) BatchGet(
	ctx context.Context, keys [][]byte,
) (kv [][]byte, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawBatchGet"), apiName+"/RawBatchGet")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err := grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchGetResponse, error) {
		return tikv.NewTikvClient(conn).RawBatchGet(ctx, &tikv.RawBatchGetRequest{
			Keys: keys,
		}, copts...)
	})
	if err != nil {
		return nil, err
	}
	// TODO
	// if res.RegionError != nil {
	kv = make([][]byte, len(res.Pairs))
	for i, pair := range res.Pairs {
		if pair.Error != nil {
			return nil, errors.Errorf("KeyError happened %+v", pair.Error)
		}
		kv[i] = pair.Value
	}
	return kv, nil
}

func (c *client) Put(
	ctx context.Context, key, val []byte,
) (err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawPut"), apiName+"/RawPut")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err := grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawPutResponse, error) {
		return tikv.NewTikvClient(conn).RawPut(ctx, &tikv.RawPutRequest{
			Key:	 key,
			Value: val,
		}, copts...)
	})
	if err != nil {
		return err
	}
	// TODO
	// if res.RegionError != nil {
	_ = res.RegionError
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

func (c *client) BatchPut(
	ctx context.Context, keys [][]byte, vals [][]byte,
) (err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawBatchPut"), apiName+"/RawBatchPut")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	pairs := make([]*tikv.KvPair, len(keys))
	for i := range keys {
		pairs[i] = &tikv.KvPair{
			Key:   keys[i],
			Value: vals[i],
		}
	}
	res, err := grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchPutResponse, error) {
		return tikv.NewTikvClient(conn).RawBatchPut(ctx, &tikv.RawBatchPutRequest{
			Pairs: pairs,
		}, copts...)
	})
	if err != nil {
		return err
	}
	// TODO
	// if res.RegionError != nil {
	_ = res.RegionError
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

func (c *client) Delete(
	ctx context.Context, key []byte,
) (err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawDelete"), apiName+"/RawDelete")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err := grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawDeleteResponse, error) {
		return tikv.NewTikvClient(conn).RawDelete(ctx, &tikv.RawDeleteRequest{
			Key: key,
		}, copts...)
	})
	if err != nil {
		return nil
	}
	// TODO
	// if res.RegionError != nil {
	_ = res.RegionError
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

func (c *client) BatchDelete(
	ctx context.Context, keys [][]byte,
) (err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/RawBatchDelete"), apiName+"/RawBatchDelete")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err := grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.RawBatchDeleteResponse, error) {
		return tikv.NewTikvClient(conn).RawBatchDelete(ctx, &tikv.RawBatchDeleteRequest{
			Keys: keys,
		}, copts...)
	})
	if err != nil {
		return err
	}
	// TODO
	// if res.RegionError != nil {
	_ = res.RegionError
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}
