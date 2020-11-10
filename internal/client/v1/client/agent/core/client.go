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

// Package core provides agent ngt gRPC client functions
package core

import (
	"context"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/internal/client/v1/client"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/net/grpc"
)

// Client represents agent NGT client interface.
type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type agentClient struct {
	vald.Client
	addr string
	c    grpc.Client
}

// New returns Client implementation if no error occurs.
func New(opts ...Option) Client {
	c := new(agentClient)
	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}
	return &agentClient{
		Client: vald.New(
			vald.WithAddr(c.addr),
			vald.WithClient(c.c),
		),
		addr: c.addr,
		c:    c.c,
	}
}

func (c *agentClient) CreateIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	_, err := c.c.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).CreateIndex(ctx, req, copts...)
		},
	)
	return nil, err
}

func (c *agentClient) SaveIndex(
	ctx context.Context,
	req *client.Empty,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	_, err := c.c.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).SaveIndex(ctx, new(client.Empty), copts...)
		},
	)
	return nil, err
}

func (c *agentClient) CreateAndSaveIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	_, err := c.c.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).CreateAndSaveIndex(ctx, req, copts...)
		},
	)
	return nil, err
}

func (c *agentClient) IndexInfo(
	ctx context.Context,
	req *client.Empty,
	opts ...grpc.CallOption,
) (res *client.InfoIndexCount, err error) {
	_, err = c.c.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			res, err := agent.NewAgentClient(conn).IndexInfo(ctx, new(client.Empty), copts...)
			if err != nil {
				return nil, err
			}
			return res, err
		},
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
