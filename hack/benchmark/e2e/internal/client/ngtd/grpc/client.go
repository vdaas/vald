package grpc

import (
	"context"

	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	proto "github.com/yahoojapan/ngtd/proto"
)

type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type ngtdClient struct {
	addr              string
	streamConcurrency int
	cfg               *config.GRPCClient
	grpc.Client
	opts []grpc.Option
}

func New(ctx context.Context, opts ...Option) (Client, error) {
	c := new(ngtdClient)

	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	c.Client = grpc.New(c.opts...)

	if err := c.Client.Connect(ctx, c.addr); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *ngtdClient) Exists(ctx context.Context, req *client.ObjectID) (*client.ObjectID, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) Search(ctx context.Context, req *client.SearchRequest) (*client.SearchResponse, error) {
	resp, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		return proto.NewNGTDClient(conn).Search(ctx, searchRequestToNGTDSearchRequest(req), copts...)
	})
	if err != nil {
		return nil, err
	}
	return toSearchResponse(resp.(*proto.SearchResponse)), nil
}

func (c *ngtdClient) SearchByID(ctx context.Context, req *client.SearchIDRequest) (*client.SearchResponse, error) {
	resp, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		return proto.NewNGTDClient(conn).SearchByID(ctx, searchIDRequestToNGTDSearchRequest(req), copts...)
	})
	if err != nil {
		return nil, err
	}
	return toSearchResponse(resp.(*proto.SearchResponse)), nil
}

func (c *ngtdClient) StreamSearch(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	_, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		st, err := proto.NewNGTDClient(conn).StreamSearch(ctx, copts...)
		if err != nil {
			return nil, err
		}

		return nil, grpc.BidirectionalStreamClient(st, c.streamConcurrency,
			func() interface{} {
				if d := dataProvider(); d != nil {
					return searchRequestToNGTDSearchRequest(d)
				}
				return nil
			}, func() interface{} {
				return new(proto.SearchResponse)
			}, func(res interface{}, err error) {
				r := res.(*proto.SearchResponse)
				f(toSearchResponse(r), err)
			})

	})
	return err
}

func (c *ngtdClient) StreamSearchByID(ctx context.Context, dataProvider func() *client.SearchIDRequest, f func(*client.SearchResponse, error)) error {
	_, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		st, err := proto.NewNGTDClient(conn).StreamSearchByID(ctx, copts...)
		if err != nil {
			return nil, err
		}
		defer st.CloseSend()

		return nil, grpc.BidirectionalStreamClient(st, c.streamConcurrency,
			func() interface{} {
				if d := dataProvider(); d != nil {
					return searchIDRequestToNGTDSearchRequest(d)
				}
				return nil
			}, func() interface{} {
				return new(proto.SearchResponse)
			}, func(res interface{}, err error) {
				r := res.(*proto.SearchResponse)
				f(toSearchResponse(r), err)
			})

	})
	return err
}

func (c *ngtdClient) Insert(ctx context.Context, req *client.ObjectVector) error {
	_, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		return proto.NewNGTDClient(conn).Insert(ctx, &proto.InsertRequest{
			Id:     []byte(req.GetId()),
			Vector: tofloat64(req.GetVector()),
		}, copts...)
	})
	return err
}

func (c *ngtdClient) StreamInsert(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) error {
	_, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		st, err := proto.NewNGTDClient(conn).StreamInsert(ctx, copts...)
		if err != nil {
			return nil, err
		}
		defer st.CloseSend()

		return nil, grpc.BidirectionalStreamClient(st, c.streamConcurrency,
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

	})
	return err
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
	_, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		return proto.NewNGTDClient(conn).Remove(ctx, &proto.RemoveRequest{
			Id: []byte(req.GetId()),
		}, copts...)
	})
	return err
}

func (c *ngtdClient) StreamRemove(ctx context.Context, dataProvider func() *client.ObjectID, f func(error)) error {
	_, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		st, err := proto.NewNGTDClient(conn).StreamRemove(ctx, copts...)
		if err != nil {
			return nil, err
		}
		defer st.CloseSend()

		return nil, grpc.BidirectionalStreamClient(st, c.streamConcurrency,
			func() interface{} {
				if d := dataProvider(); d != nil {
					return &proto.RemoveRequest{
						Id: []byte(d.GetId()),
					}
				}
				return nil
			}, func() interface{} {
				return new(proto.RemoveResponse)
			}, func(_ interface{}, err error) {
				f(err)
			})
	})
	return err
}

func (c *ngtdClient) MultiRemove(ctx context.Context, req *client.ObjectIDs) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) GetObject(ctx context.Context, req *client.ObjectID) (*client.ObjectVector, error) {
	resp, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		resp, err := proto.NewNGTDClient(conn).GetObject(ctx, &proto.GetObjectRequest{
			Id: []byte(req.GetId()),
		}, copts...)
		if err != nil {
			return nil, err
		}

		return &client.ObjectVector{
			Id:     string(resp.GetId()),
			Vector: resp.GetVector(),
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*client.ObjectVector), err
}

func (c *ngtdClient) StreamGetObject(ctx context.Context, dataProvider func() *client.ObjectID, f func(*client.ObjectVector, error)) error {
	_, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		st, err := proto.NewNGTDClient(conn).StreamGetObject(ctx, copts...)
		if err != nil {
			return nil, err
		}
		defer st.CloseSend()

		return nil, grpc.BidirectionalStreamClient(st, c.streamConcurrency,
			func() interface{} {
				if d := dataProvider(); d != nil {
					return &proto.GetObjectRequest{
						Id: []byte(d.GetId()),
					}
				}
				return nil
			}, func() interface{} {
				return new(proto.InsertResponse)
			}, func(res interface{}, err error) {
				r := res.(*proto.GetObjectResponse)
				f(&client.ObjectVector{
					Id:     string(r.GetId()),
					Vector: r.GetVector(),
				}, err)
			})

	})
	return err
}

func (c *ngtdClient) CreateIndex(ctx context.Context, req *client.ControlCreateIndexRequest) error {
	_, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		return proto.NewNGTDClient(conn).CreateIndex(ctx, &proto.CreateIndexRequest{
			PoolSize: req.GetPoolSize(),
		}, copts...)
	})
	return err
}

func (c *ngtdClient) SaveIndex(ctx context.Context) error {
	_, err := c.Client.Do(ctx, c.addr, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		return proto.NewNGTDClient(conn).SaveIndex(ctx, new(proto.Empty), copts...)
	})
	return err
}

func (c *ngtdClient) CreateAndSaveIndex(ctx context.Context, req *client.ControlCreateIndexRequest) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) IndexInfo(ctx context.Context) (*client.InfoIndex, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func searchIDRequestToNGTDSearchRequest(req *client.SearchIDRequest) *proto.SearchRequest {
	size, epsilon := getSizeAndEpsilon(req.GetConfig())
	return &proto.SearchRequest{
		Id:      []byte(req.GetId()),
		Epsilon: epsilon,
		Size_:   size,
	}
}

func searchRequestToNGTDSearchRequest(req *client.SearchRequest) *proto.SearchRequest {
	size, epsilon := getSizeAndEpsilon(req.GetConfig())
	return &proto.SearchRequest{
		Vector:  tofloat64(req.GetVector()),
		Epsilon: epsilon,
		Size_:   size,
	}
}

func getSizeAndEpsilon(cfg *client.SearchConfig) (size int32, epsilon float32) {
	if cfg != nil {
		size = int32(cfg.GetNum())
		epsilon = float32(cfg.GetEpsilon())
	}
	return
}

func toSearchResponse(in *proto.SearchResponse) (to *client.SearchResponse) {
	results := make([]*client.ObjectDistance, 0, len(in.GetResult()))
	for _, r := range results {
		results = append(results, &client.ObjectDistance{
			Id:       string(r.GetId()),
			Distance: r.GetDistance(),
		})
	}
	return &client.SearchResponse{
		Results: results,
	}
}

func toNGTDSearchResponse(in *client.SearchResponse) (to *proto.SearchResponse) {
	results := make([]*proto.ObjectDistance, 0, len(in.GetResults()))
	for _, r := range results {
		results = append(results, &proto.ObjectDistance{
			Id:       r.GetId(),
			Distance: r.GetDistance(),
		})
	}
	return &proto.SearchResponse{
		Result: results,
	}
}

func tofloat64(in []float32) (out []float64) {
	out = make([]float64, len(in))
	for i := range in {
		out[i] = float64(in[i])
	}
	return
}
