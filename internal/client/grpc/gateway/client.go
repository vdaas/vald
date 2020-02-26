package gateway

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/vald"
	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/config"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
)

type Client interface {
	client.Client
	client.MetaObjectReader
}

type gatewayClient struct {
	addr              string
	cfg               *config.GRPCClient
	grpcClient        igrpc.Client
	streamConcurrency int
}

func New(opts ...Option) (Client, error) {
	c := new(gatewayClient)

	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	c.grpcClient = igrpc.New(c.cfg.Opts()...)

	return c, nil
}

func (c *gatewayClient) Exists(ctx context.Context, objectID *client.ObjectID) (*client.ObjectID, error) {
	res, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return vald.NewValdClient(conn).Exists(ctx, objectID, copts...)
	})
	if err != nil {
		return nil, err
	}
	return res.(*client.ObjectID), nil
}

func (c *gatewayClient) Search(ctx context.Context, searchRequest *client.SearchRequest) (*client.SearchResponse, error) {
	res, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return vald.NewValdClient(conn).Search(ctx, searchRequest, copts...)
	})
	if err != nil {
		return nil, err
	}
	return res.(*client.SearchResponse), nil
}

func (c *gatewayClient) SearchByID(ctx context.Context, searchIDRequest *client.SearchIDRequest) (*client.SearchResponse, error) {
	res, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return vald.NewValdClient(conn).SearchByID(ctx, searchIDRequest, copts...)
	})
	if err != nil {
		return nil, err
	}
	return res.(*client.SearchResponse), nil
}

func (c *gatewayClient) StreamSearch(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st vald.Vald_StreamSearchClient

		st, err = vald.NewValdClient(conn).StreamSearch(ctx, copts...)
		if err != nil {
			return nil, err
		}
		defer st.CloseSend()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return dataProvider()
		}, func(res interface{}, err error) {
			f(res.(*client.SearchResponse), err)
		})
	})
	return err
}

func (c *gatewayClient) StreamSearchByID(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st vald.Vald_StreamSearchByIDClient

		st, err = vald.NewValdClient(conn).StreamSearchByID(ctx, copts...)
		if err != nil {
			return nil, err
		}
		defer st.CloseSend()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return dataProvider()
		}, func(res interface{}, err error) {
			f(res.(*client.SearchResponse), err)
		})
	})

	return err
}

func (c *gatewayClient) Insert(ctx context.Context, objectVector *client.ObjectVector) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return vald.NewValdClient(conn).Insert(ctx, objectVector, copts...)
	})
	return err
}

func (c *gatewayClient) StreamInsert(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st vald.Vald_StreamInsertClient

		st, err = vald.NewValdClient(conn).StreamInsert(ctx, copts...)
		if err != nil {
			return nil, err
		}
		defer st.CloseSend()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return dataProvider()
		}, func(_ interface{}, err error) {
			f(err)
		})

	})

	return err
}

func (c *gatewayClient) MultiInsert(ctx context.Context, objectVectors *client.ObjectVectors) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return vald.NewValdClient(conn).MultiInsert(ctx, objectVectors, copts...)
	})
	return err
}

func (c *gatewayClient) Update(ctx context.Context, objectVector *client.ObjectVector) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return vald.NewValdClient(conn).Update(ctx, objectVector, copts...)
	})
	return err
}

func (c *gatewayClient) StreamUpdate(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st vald.Vald_StreamUpdateClient

		st, err = vald.NewValdClient(conn).StreamUpdate(ctx, copts...)
		if err != nil {
			return nil, err
		}
		defer st.CloseSend()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return dataProvider()
		}, func(_ interface{}, err error) {
			f(err)
		})
	})
	return err
}

func (c *gatewayClient) MultiUpdate(ctx context.Context, objectVectors *client.ObjectVectors) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return vald.NewValdClient(conn).MultiUpdate(ctx, objectVectors, copts...)
	})
	return err
}

func (c *gatewayClient) Remove(ctx context.Context, objectID *client.ObjectID) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return vald.NewValdClient(conn).Remove(ctx, objectID, copts...)
	})
	return err
}

func (c *gatewayClient) StreamRemove(ctx context.Context, dataProvider func() *client.ObjectID, f func(error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st vald.Vald_StreamRemoveClient

		st, err = vald.NewValdClient(conn).StreamRemove(ctx, copts...)
		if err != nil {
			return nil, err
		}
		defer st.CloseSend()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return dataProvider()
		}, func(_ interface{}, err error) {
			f(err)
		})
	})
	return err
}

func (c *gatewayClient) MultiRemove(ctx context.Context, objectIDs *client.ObjectIDs) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return vald.NewValdClient(conn).MultiRemove(ctx, objectIDs, copts...)
	})
	return err
}

func (c *gatewayClient) GetObject(ctx context.Context, objectID *client.ObjectID) (*client.MetaObject, error) {
	res, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (interface{}, error) {
		return vald.NewValdClient(conn).GetObject(ctx, objectID, copts...)
	})
	if err != nil {
		return nil, err
	}
	return res.(*client.MetaObject), err
}

func (c *gatewayClient) StreamGetObject(ctx context.Context, dataProvider func() *client.ObjectID, f func(*client.MetaObject, error)) error {
	_, err := c.grpcClient.Do(ctx, c.addr, func(conn *igrpc.ClientConn, copts ...igrpc.CallOption) (res interface{}, err error) {
		var st vald.Vald_StreamGetObjectClient

		st, err = vald.NewValdClient(conn).StreamGetObject(ctx, copts...)
		if err != nil {
			return nil, err
		}
		defer st.CloseSend()

		return nil, igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
			return dataProvider()
		}, func(res interface{}, err error) {
			f(res.(*client.MetaObject), err)
		})
	})
	return err
}
