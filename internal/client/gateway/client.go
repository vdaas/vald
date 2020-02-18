package gateway

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/vald"
	"github.com/vdaas/vald/internal/client"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
	"google.golang.org/grpc"
)

type Client interface {
	client.Client
	client.MetaObjectReader
}

type gatewayClient struct {
	streamConcurrency int
	vald.ValdClient
}

func New(ctx context.Context, addr string) (Client, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &gatewayClient{
		ValdClient: vald.NewValdClient(conn),
	}, nil
}

func (c *gatewayClient) Exists(ctx context.Context, objectID *client.ObjectID) (*client.ObjectID, error) {
	return c.ValdClient.Exists(ctx, objectID)
}

func (c *gatewayClient) Search(ctx context.Context, searchRequest *client.SearchRequest) (*client.SearchResponse, error) {
	return c.ValdClient.Search(ctx, searchRequest)
}

func (c *gatewayClient) SearchByID(ctx context.Context, searchIDRequest *client.SearchIDRequest) (*client.SearchResponse, error) {
	return c.ValdClient.SearchByID(ctx, searchIDRequest)
}

func (c *gatewayClient) StreamSearch(ctx context.Context, newData func() *client.SearchRequest, f func(*client.SearchResponse, error)) (err error) {
	var st vald.Vald_StreamSearchClient

	st, err = c.ValdClient.StreamSearch(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(res interface{}, err error) {
		f(res.(*client.SearchResponse), err)
	})
}

func (c *gatewayClient) StreamSearchByID(ctx context.Context, newData func() *client.SearchRequest, f func(*client.SearchResponse, error)) (err error) {
	var st vald.Vald_StreamSearchByIDClient

	st, err = c.ValdClient.StreamSearchByID(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(res interface{}, err error) {
		f(res.(*client.SearchResponse), err)
	})
}

func (c *gatewayClient) Insert(ctx context.Context, objectVector *client.ObjectVector) error {
	_, err := c.ValdClient.Insert(ctx, objectVector)
	return err
}

func (c *gatewayClient) StreamInsert(ctx context.Context, newData func() *client.ObjectVector, f func(error)) (err error) {
	var st vald.Vald_StreamInsertClient

	st, err = c.ValdClient.StreamInsert(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(_ interface{}, err error) {
		f(err)
	})
}

func (c *gatewayClient) MultiInsert(ctx context.Context, objectVectors *client.ObjectVectors) error {
	_, err := c.ValdClient.MultiInsert(ctx, objectVectors)
	return err
}

func (c *gatewayClient) Update(ctx context.Context, objectVector *client.ObjectVector) error {
	_, err := c.ValdClient.Update(ctx, objectVector)
	return err
}

func (c *gatewayClient) StreamUpdate(ctx context.Context, newData func() *client.ObjectVector, f func(error)) (err error) {
	var st vald.Vald_StreamUpdateClient

	st, err = c.ValdClient.StreamUpdate(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(_ interface{}, err error) {
		f(err)
	})
}

func (c *gatewayClient) MultiUpdate(ctx context.Context, objectVectors *client.ObjectVectors) error {
	_, err := c.ValdClient.MultiUpdate(ctx, objectVectors)
	return err
}

func (c *gatewayClient) Remove(ctx context.Context, objectID *client.ObjectID) error {
	_, err := c.ValdClient.Remove(ctx, objectID)
	return err
}

func (c *gatewayClient) StreamRemove(ctx context.Context, newData func() *client.ObjectID, f func(error)) (err error) {
	var st vald.Vald_StreamRemoveClient

	st, err = c.ValdClient.StreamRemove(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(_ interface{}, err error) {
		f(err)
	})
}

func (c *gatewayClient) MultiRemove(ctx context.Context, objectIDs *client.ObjectIDs) error {
	_, err := c.ValdClient.MultiRemove(ctx, objectIDs)
	return err
}

func (c *gatewayClient) GetObject(ctx context.Context, objectID *client.ObjectID) (*client.MetaObject, error) {
	return c.ValdClient.GetObject(ctx, objectID)
}

func (c *gatewayClient) StreamGetObject(ctx context.Context, newData func() *client.ObjectID, f func(*client.MetaObject, error)) (err error) {
	var st vald.Vald_StreamGetObjectClient

	st, err = c.ValdClient.StreamGetObject(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = st.CloseSend()
	}()

	return igrpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return newData()
	}, func(res interface{}, err error) {
		f(res.(*client.MetaObject), err)
	})
}
