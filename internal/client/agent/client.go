package agent

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/config"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
)

type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type agentClient struct {
	addr              string
	cfg               *config.GRPCClient
	streamConcurrency int
	grpcClient        igrpc.Client
}

func New(ctx context.Context, opts ...Option) (Client, error) {
	c := new(agentClient)

	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	c.grpcClient = igrpc.New(c.cfg.Opts()...)

	return c, nil
}

func (c *agentClient) Exists(ctx context.Context, objectID *client.ObjectID) (*client.ObjectID, error) {
	res, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).Exists(ctx, objectID)
	})
	if err != nil {
		return nil, err
	}
	return res.(*client.ObjectID), nil
}

func (c *agentClient) Search(ctx context.Context, searchRequest *client.SearchRequest) (*client.SearchResponse, error) {
	res, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).Search(ctx, searchRequest)
	})
	if err != nil {
		return nil, err
	}
	return res.(*client.SearchResponse), nil
}

func (c *agentClient) SearchByID(ctx context.Context, searchIDRequest *client.SearchIDRequest) (*client.SearchResponse, error) {
	res, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).SearchByID(ctx, searchIDRequest)
	})
	if err != nil {
		return nil, err
	}
	return res.(*client.SearchResponse), nil
}

func (c *agentClient) StreamSearch(ctx context.Context, newData func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st agent.Agent_StreamSearchClient

		st, err = agent.NewAgentClient(conn).StreamSearch(ctx)
		if err != nil {
			return nil, err
		}

		defer func() {
			err = st.CloseSend()
		}()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return newData()
		}, func(res interface{}, err error) {
			f(res.(*client.SearchResponse), err)
		})

	})
	return err
}

func (c *agentClient) StreamSearchByID(ctx context.Context, newData func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st agent.Agent_StreamSearchByIDClient

		st, err = agent.NewAgentClient(conn).StreamSearchByID(ctx)
		if err != nil {
			return nil, err
		}

		defer func() {
			err = st.CloseSend()
		}()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return newData()
		}, func(res interface{}, err error) {
			f(res.(*client.SearchResponse), err)
		})

	})
	return err
}

func (c *agentClient) Insert(ctx context.Context, objectVector *client.ObjectVector) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).Insert(ctx, objectVector)
	})
	return err
}

func (c *agentClient) StreamInsert(ctx context.Context, newData func() *client.ObjectVector, f func(error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st agent.Agent_StreamInsertClient

		st, err = agent.NewAgentClient(conn).StreamInsert(ctx)
		if err != nil {
			return nil, err
		}

		defer func() {
			err = st.CloseSend()
		}()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return newData()
		}, func(_ interface{}, err error) {
			f(err)
		})

	})
	return err
}

func (c *agentClient) MultiInsert(ctx context.Context, objectVectors *client.ObjectVectors) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).MultiInsert(ctx, objectVectors)
	})
	return err
}

func (c *agentClient) Update(ctx context.Context, objectVector *client.ObjectVector) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).Update(ctx, objectVector)
	})
	return err
}

func (c *agentClient) StreamUpdate(ctx context.Context, newData func() *client.ObjectVector, f func(error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st agent.Agent_StreamUpdateClient

		st, err = agent.NewAgentClient(conn).StreamUpdate(ctx)
		if err != nil {
			return nil, err
		}

		defer func() {
			err = st.CloseSend()
		}()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return newData()
		}, func(_ interface{}, err error) {
			f(err)
		})
	})
	return err
}

func (c *agentClient) MultiUpdate(ctx context.Context, objectVectors *client.ObjectVectors) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).MultiUpdate(ctx, objectVectors)
	})
	return err
}

func (c *agentClient) Remove(ctx context.Context, objectID *client.ObjectID) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).Remove(ctx, objectID)
	})
	return err
}

func (c *agentClient) StreamRemove(ctx context.Context, newData func() *client.ObjectID, f func(error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st agent.Agent_StreamRemoveClient

		st, err = agent.NewAgentClient(conn).StreamRemove(ctx)
		if err != nil {
			return nil, err
		}

		defer func() {
			err = st.CloseSend()
		}()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return newData()
		}, func(_ interface{}, err error) {
			f(err)
		})
	})
	return err
}

func (c *agentClient) MultiRemove(ctx context.Context, objectIDs *client.ObjectIDs) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).MultiRemove(ctx, objectIDs)
	})
	return err
}

func (c *agentClient) GetObject(ctx context.Context, objectID *client.ObjectID) (*client.ObjectVector, error) {
	res, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).GetObject(ctx, objectID)
	})
	if err != nil {
		return nil, err
	}
	return res.(*client.ObjectVector), nil
}

func (c *agentClient) StreamGetObject(ctx context.Context, newData func() *client.ObjectID, f func(*client.ObjectVector, error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st agent.Agent_StreamGetObjectClient

		st, err = agent.NewAgentClient(conn).StreamGetObject(ctx)
		if err != nil {
			return nil, err
		}

		defer func() {
			err = st.CloseSend()
		}()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return newData()
		}, func(res interface{}, err error) {
			f(res.(*client.ObjectVector), err)
		})
	})
	return err
}

func (c *agentClient) CreateIndex(ctx context.Context, controlCreateIndexRequest *client.ControlCreateIndexRequest) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).CreateIndex(ctx, controlCreateIndexRequest)
	})
	return err
}

func (c *agentClient) SaveIndex(ctx context.Context) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).SaveIndex(ctx, new(payload.Empty))
	})
	return err
}

func (c *agentClient) CreateAndSaveIndex(ctx context.Context, controlCreateIndexRequest *client.ControlCreateIndexRequest) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).CreateAndSaveIndex(ctx, controlCreateIndexRequest)
	})
	return err
}

func (c *agentClient) IndexInfo(ctx context.Context) (*client.InfoIndex, error) {
	res, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return agent.NewAgentClient(conn).IndexInfo(ctx, new(payload.Empty))
	})
	if err != nil {
		return nil, err
	}
	return res.(*client.InfoIndex), err
}
