//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package core provides agent ngt gRPC client functions
package core

import (
	"context"
	"reflect"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/internal/client/v1/client"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

// Client represents agent NGT client interface.
type Client interface {
	vald.Client
	client.ObjectReader
	client.Indexer
}

type agentClient struct {
	vald.Client
	addrs []string
	c     grpc.Client
}

type singleAgentClient struct {
	vald.Client
	ac agent.AgentClient
}

const (
	apiName = "vald/internal/client/v1/client/agent/core"
)

// New returns Client implementation if no error occurs.
func New(opts ...Option) (Client, error) {
	c := new(agentClient)
	for _, opt := range append(defaultOptions, opts...) {
		err := opt(c)
		if err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	if c.c == nil {
		if c.Client != nil {
			c.c = c.Client.GRPCClient()
		} else {
			if c.addrs == nil {
				return nil, errors.ErrGRPCTargetAddrNotFound
			}
			c.c = grpc.New(grpc.WithAddrs(c.addrs...))
		}
	}
	if c.Client == nil {
		if c.addrs == nil {
			return nil, errors.ErrGRPCTargetAddrNotFound
		}
		var err error
		c.Client, err = vald.New(
			vald.WithAddrs(c.addrs...),
			vald.WithClient(c.c),
		)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func NewAgentClient(cc *grpc.ClientConn) interface {
	vald.Client
	client.ObjectReader
	client.Indexer
} {
	return &singleAgentClient{
		Client: vald.NewValdClient(cc),
		ac:     agent.NewAgentClient(cc),
	}
}

func (c *agentClient) CreateIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
	_ ...grpc.CallOption,
) (*client.Empty, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+agent.CreateIndexRPCName), apiName+"/"+agent.CreateIndexRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err := c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (interface{}, error) {
		return agent.NewAgentClient(conn).CreateIndex(ctx, req, copts...)
	})
	return nil, err
}

func (c *agentClient) SaveIndex(
	ctx context.Context,
	_ *client.Empty,
	_ ...grpc.CallOption,
) (*client.Empty, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+agent.SaveIndexRPCName), apiName+"/"+agent.SaveIndexRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err := c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (interface{}, error) {
		return agent.NewAgentClient(conn).SaveIndex(ctx, new(client.Empty), copts...)
	})
	return nil, err
}

func (c *agentClient) CreateAndSaveIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
	_ ...grpc.CallOption,
) (*client.Empty, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+agent.CreateAndSaveIndexRPCName), apiName+"/"+agent.CreateAndSaveIndexRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err := c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (interface{}, error) {
		return agent.NewAgentClient(conn).CreateAndSaveIndex(ctx, req, copts...)
	})
	return nil, err
}

func (c *agentClient) IndexInfo(
	ctx context.Context,
	_ *client.Empty,
	_ ...grpc.CallOption,
) (res *client.InfoIndexCount, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+agent.IndexInfoRPCName), apiName+"/"+agent.IndexInfoRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err := agent.NewAgentClient(conn).IndexInfo(ctx, new(client.Empty), copts...)
		if err != nil {
			return nil, err
		}
		return res, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) GetObjectMeta(
	ctx context.Context,
	req *client.ObjectMetaRequest,
	_ ...grpc.CallOption,
) (res *client.ObjectMeta, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+agent.GetObjectMetaRPCName), apiName+"/"+agent.GetObjectMetaRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err := agent.NewAgentClient(conn).GetObjectMeta(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return res, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *singleAgentClient) CreateIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+agent.CreateIndexRPCName), apiName+"/"+agent.CreateIndexRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.ac.CreateIndex(ctx, req, opts...)
}

func (c *singleAgentClient) SaveIndex(
	ctx context.Context,
	_ *client.Empty,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+agent.SaveIndexRPCName), apiName+"/"+agent.SaveIndexRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.ac.SaveIndex(ctx, new(client.Empty), opts...)
}

func (c *singleAgentClient) CreateAndSaveIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+agent.CreateAndSaveIndexRPCName), apiName+"/"+agent.CreateAndSaveIndexRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.ac.CreateAndSaveIndex(ctx, req, opts...)
}

func (c *singleAgentClient) IndexInfo(
	ctx context.Context,
	_ *client.Empty,
	opts ...grpc.CallOption,
) (res *client.InfoIndexCount, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+agent.IndexInfoRPCName), apiName+"/"+agent.IndexInfoRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.ac.IndexInfo(ctx, new(client.Empty), opts...)
}

func (c *singleAgentClient) GetObjectMeta(
	ctx context.Context,
	req *client.ObjectMetaRequest,
	opts ...grpc.CallOption,
) (res *client.ObjectMeta, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+agent.IndexInfoRPCName), apiName+"/"+agent.IndexInfoRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.ac.GetObjectMeta(ctx, req, opts...)
}
