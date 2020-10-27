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

// Package grpc provides agent ngt gRPC client functions
package grpc

import (
	"context"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client"
	"github.com/vdaas/vald/internal/net/grpc"
)

// Client represents agent NGT client interface.
type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type agentClient struct {
	addr string
	c    grpc.Client
}

// New returns Client implementation if no error occurs.
func New(opts ...Option) Client {
	c := new(agentClient)
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *agentClient) Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (oid *payload.Object_ID, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		oid, err = vald.NewValdClient(conn).Exists(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return oid, nil
}

func (c *agentClient) Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).Search(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).SearchByID(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamSearch(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchByIDClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamSearchByID(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) MultiSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiSearch(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) MultiSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiSearchByID(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).Insert(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Insert_StreamInsertClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamInsert(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) MultiInsert(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiInsert(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).Update(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (res vald.Update_StreamUpdateClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamUpdate(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) MultiUpdate(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiUpdate(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).Upsert(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) StreamUpsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Upsert_StreamUpsertClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamUpsert(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) MultiUpsert(ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiUpsert(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).Remove(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (res vald.Remove_StreamRemoveClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamRemove(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) MultiRemove(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiRemove(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) GetObject(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (res *payload.Object_Vector, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).GetObject(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *agentClient) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (res vald.Object_StreamGetObjectClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamGetObject(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
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
