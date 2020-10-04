//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package grpc provides vald gRPC client functions
package grpc

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/net/grpc"
)

// Client represents gateway client interface.
type Client interface {
	client.Client
	client.MetaObjectReader
	client.Upserter
}

type gatewayClient struct {
	addr string
	cfg  *config.GRPCClient
	grpc.Client
}

// New returns Client implementation if no error occurs.
func New(ctx context.Context, opts ...Option) (Client, error) {
	c := new(gatewayClient)

	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	c.Client = grpc.New(c.cfg.Opts()...)

	if err := c.Client.Connect(ctx, c.addr); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *gatewayClient) Exists(
	ctx context.Context,
	req *client.ObjectID,
) (*client.ObjectID, error) {
	res, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).Exists(ctx, req, copts...)
		},
	)
	if err != nil {
		return nil, err
	}
	return res.(*client.ObjectID), nil
}

func (c *gatewayClient) Search(
	ctx context.Context,
	req *client.SearchRequest,
) (*client.SearchResponse, error) {
	res, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).Search(ctx, req, copts...)
		},
	)
	if err != nil {
		return nil, err
	}
	return res.(*client.SearchResponse), nil
}

func (c *gatewayClient) SearchByID(
	ctx context.Context,
	req *client.SearchIDRequest,
) (*client.SearchResponse, error) {
	res, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).SearchByID(ctx, req, copts...)
		},
	)
	if err != nil {
		return nil, err
	}
	return res.(*client.SearchResponse), nil
}

func (c *gatewayClient) StreamSearch(
	ctx context.Context,
	dataProvider func() *client.SearchRequest,
	f func(*client.SearchResponse, error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (res interface{}, err error) {
			var st vald.Vald_StreamSearchClient

			st, err = vald.NewValdClient(conn).StreamSearch(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, streamSearch(st,
				func() interface{} {
					if d := dataProvider(); d != nil {
						return d
					}
					return nil
				}, f,
			)
		},
	)
	return err
}

func (c *gatewayClient) StreamSearchByID(
	ctx context.Context,
	dataProvider func() *client.SearchIDRequest,
	f func(*client.SearchResponse, error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (res interface{}, err error) {
			var st vald.Vald_StreamSearchByIDClient

			st, err = vald.NewValdClient(conn).StreamSearchByID(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, streamSearch(st,
				func() interface{} {
					if d := dataProvider(); d != nil {
						return d
					}
					return nil
				}, f,
			)
		},
	)
	return err
}

func (c *gatewayClient) Insert(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).Insert(ctx, req, copts...)
		},
	)
	return err
}

func (c *gatewayClient) StreamInsert(
	ctx context.Context,
	dataProvider func() *client.ObjectVector,
	f func(error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (res interface{}, err error) {
			var st vald.Vald_StreamInsertClient

			st, err = vald.NewValdClient(conn).StreamInsert(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, stream(st, func() interface{} {
				if d := dataProvider(); d != nil {
					return d
				}
				return nil
			}, f)
		},
	)
	return err
}

func (c *gatewayClient) MultiInsert(
	ctx context.Context,
	req *client.ObjectVectors,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).MultiInsert(ctx, req, copts...)
		},
	)
	return err
}

func (c *gatewayClient) Update(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).Update(ctx, req, copts...)
		},
	)
	return err
}

func (c *gatewayClient) StreamUpdate(
	ctx context.Context,
	dataProvider func() *client.ObjectVector,
	f func(error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (res interface{}, err error) {
			var st vald.Vald_StreamUpdateClient

			st, err = vald.NewValdClient(conn).StreamUpdate(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, stream(st, func() interface{} {
				if d := dataProvider(); d != nil {
					return d
				}
				return nil
			}, f)
		},
	)
	return err
}

func (c *gatewayClient) MultiUpdate(
	ctx context.Context,
	req *client.ObjectVectors,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).MultiUpdate(ctx, req, copts...)
		},
	)
	return err
}

func (c *gatewayClient) Upsert(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).Upsert(ctx, req, copts...)
		},
	)
	return err
}

func (c *gatewayClient) MultiUpsert(
	ctx context.Context,
	req *client.ObjectVectors,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).MultiUpsert(ctx, req, copts...)
		},
	)
	return err
}

func (c *gatewayClient) StreamUpsert(
	ctx context.Context,
	dataProvider func() *client.ObjectVector,
	f func(error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			var st vald.Vald_StreamUpsertClient

			st, err := vald.NewValdClient(conn).StreamUpsert(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, stream(st,
				func() interface{} {
					if d := dataProvider(); d != nil {
						return d
					}
					return nil
				}, f,
			)
		},
	)
	return err
}

func (c *gatewayClient) Remove(
	ctx context.Context,
	req *client.ObjectID,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).Remove(ctx, req, copts...)
		},
	)
	return err
}

func (c *gatewayClient) StreamRemove(
	ctx context.Context,
	dataProvider func() *client.ObjectID,
	f func(error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (res interface{}, err error) {
			var st vald.Vald_StreamRemoveClient

			st, err = vald.NewValdClient(conn).StreamRemove(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, stream(st,
				func() interface{} {
					return dataProvider()
				}, f,
			)
		},
	)
	return err
}

func (c *gatewayClient) MultiRemove(
	ctx context.Context,
	req *client.ObjectIDs,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).MultiRemove(ctx, req, copts...)
		},
	)
	return err
}

func (c *gatewayClient) GetObject(
	ctx context.Context,
	req *client.ObjectID,
) (*client.MetaObject, error) {
	res, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).GetObject(ctx, req, copts...)
		},
	)
	if err != nil {
		return nil, err
	}
	return res.(*client.MetaObject), err
}

func (c *gatewayClient) StreamGetObject(
	ctx context.Context,
	dataProvider func() *client.ObjectID,
	f func(*client.MetaObject, error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (res interface{}, err error) {
			var st vald.Vald_StreamGetObjectClient

			st, err = vald.NewValdClient(conn).StreamGetObject(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, grpc.BidirectionalStreamClient(st,
				func() interface{} {
					return dataProvider()
				}, func() interface{} {
					return new(client.BackupMetaVector)
				}, func(res interface{}, err error) {
					f(res.(*client.MetaObject), err)
				})
		},
	)
	return err
}

func streamSearch(
	st grpc.ClientStream,
	dataProvider func() interface{},
	f func(*client.SearchResponse, error),
) error {
	return grpc.BidirectionalStreamClient(st, dataProvider,
		func() interface{} {
			return new(client.SearchResponse)
		}, func(res interface{}, err error) {
			f(res.(*client.SearchResponse), err)
		},
	)
}

func stream(
	st grpc.ClientStream,
	dataProvider func() interface{},
	f func(error),
) error {
	return grpc.BidirectionalStreamClient(st, dataProvider,
		func() interface{} {
			return new(client.Empty)
		}, func(_ interface{}, err error) {
			f(err)
		},
	)
}
