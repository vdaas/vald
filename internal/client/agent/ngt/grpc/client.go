// Package grpc provides gRPC client functions
package grpc

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/internal/client"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
	"google.golang.org/grpc"
)

// Client represents agent NGT client interface.
type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type agentClient struct {
	addr string
	opts []igrpc.Option
	igrpc.Client
}

// New returns Client implementation if no error occurs.
func New(ctx context.Context, opts ...Option) (Client, error) {
	c := new(agentClient)
	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	c.Client = igrpc.New(c.opts...)

	if err := c.Client.Connect(ctx, c.addr); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *agentClient) Exists(
	ctx context.Context,
	req *client.ObjectID,
) (*client.ObjectID, error) {
	res, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).Exists(ctx, req, copts...)
		},
	)
	if err != nil {
		return nil, err
	}
	return res.(*client.ObjectID), nil
}

func (c *agentClient) Search(
	ctx context.Context,
	req *client.SearchRequest,
) (*client.SearchResponse, error) {
	res, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).Search(ctx, req, copts...)
		},
	)
	if err != nil {
		return nil, err
	}
	return res.(*client.SearchResponse), nil
}

func (c *agentClient) SearchByID(
	ctx context.Context,
	req *client.SearchIDRequest,
) (*client.SearchResponse, error) {
	res, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).SearchByID(ctx, req, copts...)
		},
	)
	if err != nil {
		return nil, err
	}
	return res.(*client.SearchResponse), nil
}

func (c *agentClient) StreamSearch(
	ctx context.Context,
	dataProvider func() *client.SearchRequest,
	f func(*client.SearchResponse, error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
			var st agent.Agent_StreamSearchClient

			st, err = agent.NewAgentClient(conn).StreamSearch(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, streamSearch(st,
				func() interface{} {
					if d := dataProvider(); d != nil {
						return d
					}
					return nil
				}, f)
		},
	)
	return err
}

func (c *agentClient) StreamSearchByID(
	ctx context.Context,
	dataProvider func() *client.SearchIDRequest,
	f func(*client.SearchResponse, error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
			var st agent.Agent_StreamSearchByIDClient

			st, err = agent.NewAgentClient(conn).StreamSearchByID(ctx, copts...)
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

func (c *agentClient) Insert(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).Insert(ctx, req, copts...)
		},
	)
	return err
}

func (c *agentClient) StreamInsert(
	ctx context.Context,
	dataProvider func() *client.ObjectVector,
	f func(error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
			var st agent.Agent_StreamInsertClient

			st, err = agent.NewAgentClient(conn).StreamInsert(ctx, copts...)
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

func (c *agentClient) MultiInsert(
	ctx context.Context,
	req *client.ObjectVectors,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).MultiInsert(ctx, req, copts...)
		},
	)
	return err
}

func (c *agentClient) Update(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).Update(ctx, req, copts...)
		},
	)
	return err
}

func (c *agentClient) StreamUpdate(
	ctx context.Context,
	dataProvider func() *client.ObjectVector,
	f func(error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
			var st agent.Agent_StreamUpdateClient

			st, err = agent.NewAgentClient(conn).StreamUpdate(ctx, copts...)
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

func (c *agentClient) MultiUpdate(
	ctx context.Context,
	req *client.ObjectVectors,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).MultiUpdate(ctx, req, copts...)
		},
	)
	return err
}

func (c *agentClient) Remove(
	ctx context.Context,
	req *client.ObjectID,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).Remove(ctx, req, copts...)
		},
	)
	return err
}

func (c *agentClient) StreamRemove(
	ctx context.Context,
	dataProvider func() *client.ObjectID,
	f func(error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			st, err := agent.NewAgentClient(conn).StreamRemove(ctx, copts...)
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

func (c *agentClient) MultiRemove(
	ctx context.Context,
	req *client.ObjectIDs,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).MultiRemove(ctx, req, copts...)
		},
	)
	return err
}

func (c *agentClient) GetObject(
	ctx context.Context,
	req *client.ObjectID,
) (*client.ObjectVector, error) {
	res, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).GetObject(ctx, req, copts...)
		},
	)
	if err != nil {
		return nil, err
	}
	return res.(*client.ObjectVector), nil
}

func (c *agentClient) StreamGetObject(
	ctx context.Context,
	dataProvider func() *client.ObjectID,
	f func(*client.ObjectVector, error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
			var st agent.Agent_StreamGetObjectClient

			st, err = agent.NewAgentClient(conn).StreamGetObject(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, igrpc.BidirectionalStreamClient(st,
				func() interface{} {
					if d := dataProvider(); d != nil {
						return d
					}
					return nil
				}, func() interface{} {
					return new(client.ObjectVector)
				}, func(res interface{}, err error) {
					f(res.(*client.ObjectVector), err)
				})
		},
	)
	return err
}

func (c *agentClient) CreateIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).CreateIndex(ctx, req, copts...)
		},
	)
	return err
}

func (c *agentClient) SaveIndex(ctx context.Context) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).SaveIndex(ctx, new(client.Empty), copts...)
		},
	)
	return err
}

func (c *agentClient) CreateAndSaveIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).CreateAndSaveIndex(ctx, req, copts...)
		},
	)
	return err
}

func (c *agentClient) IndexInfo(ctx context.Context) (*client.InfoIndex, error) {
	res, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
			return agent.NewAgentClient(conn).IndexInfo(ctx, new(client.Empty), copts...)
		},
	)
	if err != nil {
		return nil, err
	}
	return res.(*client.InfoIndex), err
}

func streamSearch(
	st grpc.ClientStream,
	dataProvider func() interface{},
	f func(*client.SearchResponse, error),
) error {
	return igrpc.BidirectionalStreamClient(st, dataProvider,
		func() interface{} {
			return new(client.SearchResponse)
		}, func(res interface{}, err error) {
			f(res.(*client.SearchResponse), err)
		})
}

func stream(
	st grpc.ClientStream,
	dataProvider func() interface{},
	f func(error),
) error {
	return igrpc.BidirectionalStreamClient(st, dataProvider,
		func() interface{} {
			return new(client.Empty)
		}, func(_ interface{}, err error) {
			f(err)
		})
}
