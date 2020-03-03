package grpc

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	proto "github.com/yahoojapan/ngtd/proto"
	ggrpc "google.golang.org/grpc"
)

type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type ngtdClient struct {
	addr              string
	streamConcurrency int
	proto.NGTDClient
}

func New(ctx context.Context, opts ...Option) (Client, error) {
	c := new(ngtdClient)

	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	conn, err := ggrpc.DialContext(ctx, c.addr, ggrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c.NGTDClient = proto.NewNGTDClient(conn)

	return c, nil
}

func (c *ngtdClient) Exists(ctx context.Context, req *client.ObjectID) (*client.ObjectID, error) {
	// TODO: errors.NotSupportedClientMethod
	return nil, nil
}

func (c *ngtdClient) Search(ctx context.Context, req *client.SearchRequest) (*client.SearchResponse, error) {
	resp, err := c.NGTDClient.Search(ctx, &proto.SearchRequest{
		Vector:  tofloat64(req.GetVector()),
		Epsilon: req.GetConfig().GetEpsilon(),
		Size_:   int32(req.GetConfig().GetNum()),
	})
	if err != nil {
		return nil, err
	}

	return &client.SearchResponse{
		Results: toObjectDistances(resp.GetResult()),
	}, nil
}

func (c *ngtdClient) SearchByID(ctx context.Context, req *client.SearchIDRequest) (*client.SearchResponse, error) {
	resp, err := c.NGTDClient.SearchByID(ctx, &proto.SearchRequest{
		Id:      []byte(req.GetId()),
		Epsilon: req.GetConfig().GetEpsilon(),
		Size_:   int32(req.GetConfig().GetNum()),
	})
	if err != nil {
		return nil, err
	}

	return &client.SearchResponse{
		Results: toObjectDistances(resp.GetResult()),
	}, nil
}

func (c *ngtdClient) StreamSearch(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	st, err := c.NGTDClient.StreamSearch(ctx)
	if err != nil {
		return err
	}
	defer st.CloseSend()

	/**
	grpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return dataProvider()
	}, func(res interface{}, err error) {
		_ = client.SearchResponse{
			Results: toObjectDistances(res.GetResult()),
		}
	})
	*/

	return nil
}

func (c *ngtdClient) StreamSearchByID(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	st, err := c.NGTDClient.StreamSearchByID(ctx)
	if err != nil {
		return err
	}
	defer st.CloseSend()

	/**
	grpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return dataProvider()
	}, func(res interface{}, err error) {
		_ = client.SearchResponse{
			Results: toObjectDistances(res.GetResult()),
		}
	})
	*/

	return nil
}

func (c *ngtdClient) Insert(ctx context.Context, req *client.ObjectVector) error {
	_, err := c.NGTDClient.Insert(ctx, &proto.InsertRequest{
		Id:     []byte(req.GetId()),
		Vector: tofloat64(req.GetVector()),
	})
	return err
}

func (c *ngtdClient) StreamInsert(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) error {
	st, err := c.NGTDClient.StreamInsert(ctx)
	if err != nil {
		return err
	}
	defer st.CloseSend()

	return grpc.BidirectionalStreamClient(st, c.streamConcurrency,
		func() interface{} {
			if d := dataProvider(); d != nil {
				return &proto.InsertRequest{
					Id:     []byte(d.GetId()),
					Vector: tofloat64(d.GetVector()),
				}
			}
			return nil
		}, func() interface{} {
			return new(proto.InsertResponse)
		}, func(_ interface{}, err error) {
			f(err)
		})
}

func (c *ngtdClient) MultiInsert(ctx context.Context, req *client.ObjectVectors) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) Update(ctx context.Context, req *client.ObjectVector) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) StreamUpdate(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) MultiUpdate(ctx context.Context, req *client.ObjectVectors) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) Remove(ctx context.Context, req *client.ObjectID) error {
	_, err := c.NGTDClient.Remove(ctx, &proto.RemoveRequest{
		Id: []byte(req.GetId()),
	})
	return err
}

func (c *ngtdClient) StreamRemove(ctx context.Context, dataProvider func() *client.ObjectID, f func(error)) error {
	st, err := c.NGTDClient.StreamRemove(ctx)
	if err != nil {
		return err
	}
	defer st.CloseSend()

	return nil
	/**
	grpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return dataProvider()
	}, func(_ interface{}, err error) {
		f(err)
	})
	*/
}

func (c *ngtdClient) MultiRemove(ctx context.Context, req *client.ObjectIDs) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) GetObject(ctx context.Context, req *client.ObjectID) (*client.ObjectVector, error) {
	resp, err := c.NGTDClient.GetObject(ctx, &proto.GetObjectRequest{
		Id: []byte(req.GetId()),
	})
	if err != nil {
		return nil, err
	}

	return &client.ObjectVector{
		Id:     string(resp.GetId()),
		Vector: resp.GetVector(),
	}, nil
}

func (c *ngtdClient) StreamGetObject(ctx context.Context, dataProvider func() *client.ObjectID, f func(*client.ObjectVector, error)) error {
	st, err := c.NGTDClient.StreamGetObject(ctx)
	if err != nil {
		return err
	}
	defer st.CloseSend()

	/**
	grpc.BidirectionalStreamClient(st, c.streamConcurrency, func() interface{} {
		return dataProvider()
	}, func(res interface{}, err error) {
		r := res.(*proto.GetObjectResponse)
		f(&client.ObjectVector {
			Id: string(resp.GetId()),
			Vector: resp.GetVector(),
		},err)
	})
	*/
	return nil
}

func (c *ngtdClient) CreateIndex(ctx context.Context, req *client.ControlCreateIndexRequest) error {
	_, err := c.NGTDClient.CreateIndex(ctx, &proto.CreateIndexRequest{
		PoolSize: req.GetPoolSize(),
	})
	return err
}

func (c *ngtdClient) SaveIndex(ctx context.Context) error {
	_, err := c.NGTDClient.SaveIndex(ctx, new(proto.Empty))
	return err
}

func (c *ngtdClient) CreateAndSaveIndex(ctx context.Context, req *client.ControlCreateIndexRequest) error {
	// TODO: errors.NotSupportedClientMethod
	return nil
}

func (c *ngtdClient) IndexInfo(ctx context.Context) (*client.InfoIndex, error) {
	// TODO: errors.NotSupportedClientMethod
	return nil, nil
}

func tofloat64(in []float32) (out []float64) {
	out = make([]float64, len(in))
	for i := range in {
		out[i] = float64(in[i])
	}
	return
}

func toObjectDistances(in []*proto.ObjectDistance) (to []*payload.Object_Distance) {
	to = make([]*payload.Object_Distance, 0, len(in))

	for _, elm := range in {
		to = append(to, &payload.Object_Distance{
			Id:       string(elm.GetId()),
			Distance: elm.GetDistance(),
		})
	}
	return nil
}
