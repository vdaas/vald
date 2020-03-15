// Package grpc provides grpc client functions
package grpc

import (
	"context"

	"github.com/vdaas/vald/internal/client"
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
	addr string
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

func (c *ngtdClient) Exists(
	ctx context.Context,
	req *client.ObjectID,
) (*client.ObjectID, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) Search(
	ctx context.Context,
	req *client.SearchRequest,
) (*client.SearchResponse, error) {
	res, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			res, err := proto.NewNGTDClient(conn).Search(ctx, searchRequestToNgtdSearchRequest(req), copts...)
			if err != nil {
				return nil, err
			}

			if len(res.GetError()) != 0 {
				return nil, errors.New(res.GetError())
			}
			return res, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return ngtdSearchResponseToSearchResponse(res.(*proto.SearchResponse)), nil
}

func (c *ngtdClient) SearchByID(
	ctx context.Context,
	req *client.SearchIDRequest,
) (*client.SearchResponse, error) {
	res, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			res, err := proto.NewNGTDClient(conn).SearchByID(ctx, searchIDRequestToNgtdSearchRequest(req), copts...)
			if err != nil {
				return nil, err
			}

			if len(res.GetError()) != 0 {
				return nil, errors.New(res.GetError())
			}
			return res, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return ngtdSearchResponseToSearchResponse(res.(*proto.SearchResponse)), nil
}

func (c *ngtdClient) StreamSearch(
	ctx context.Context,
	dataProvider func() *client.SearchRequest,
	f func(*client.SearchResponse, error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			st, err := proto.NewNGTDClient(conn).StreamSearch(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, grpc.BidirectionalStreamClient(st,
				func() interface{} {
					if d := dataProvider(); d != nil {
						return searchRequestToNgtdSearchRequest(d)
					}
					return nil
				}, func() interface{} {
					return new(proto.SearchResponse)
				}, func(intr interface{}, err error) {
					if err != nil {
						f(nil, err)
						return
					}

					res := intr.(*proto.SearchResponse)
					if len(res.GetError()) != 0 {
						f(nil, errors.New(res.GetError()))
						return
					}

					f(ngtdSearchResponseToSearchResponse(res), err)
				})
		},
	)
	return err
}

func (c *ngtdClient) StreamSearchByID(
	ctx context.Context,
	dataProvider func() *client.SearchIDRequest,
	f func(*client.SearchResponse, error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			st, err := proto.NewNGTDClient(conn).StreamSearchByID(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, grpc.BidirectionalStreamClient(st,
				func() interface{} {
					if d := dataProvider(); d != nil {
						return searchIDRequestToNgtdSearchRequest(d)
					}
					return nil
				}, func() interface{} {
					return new(proto.SearchResponse)
				}, func(intr interface{}, err error) {
					if err != nil {
						f(nil, err)
						return
					}

					res := intr.(*proto.SearchResponse)
					if len(res.GetError()) != 0 {
						f(nil, errors.New(res.GetError()))
						return
					}

					f(ngtdSearchResponseToSearchResponse(res), err)
				})
		},
	)
	return err
}

func (c *ngtdClient) Insert(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			res, err := proto.NewNGTDClient(conn).Insert(ctx, objectVectorToNGTDInsertRequest(req), copts...)
			if err != nil {
				return nil, err
			}

			if len(res.GetError()) != 0 {
				return nil, errors.New(res.GetError())
			}
			return res, nil
		},
	)
	return err
}

func (c *ngtdClient) StreamInsert(
	ctx context.Context,
	dataProvider func() *client.ObjectVector,
	f func(error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			st, err := proto.NewNGTDClient(conn).StreamInsert(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, grpc.BidirectionalStreamClient(st,
				func() interface{} {
					if d := dataProvider(); d != nil {
						return objectVectorToNGTDInsertRequest(d)
					}
					return nil
				}, func() interface{} {
					return new(proto.InsertResponse)
				}, func(intr interface{}, err error) {
					if err != nil {
						f(err)
						return
					}

					res := intr.(*proto.InsertResponse)
					if len(res.GetError()) != 0 {
						f(errors.New(res.GetError()))
						return
					}

					f(err)
				})
		},
	)
	return err
}

func (c *ngtdClient) MultiInsert(
	ctx context.Context,
	req *client.ObjectVectors,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) Update(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) StreamUpdate(
	ctx context.Context,
	dataProvider func() *client.ObjectVector,
	f func(error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) MultiUpdate(
	ctx context.Context,
	req *client.ObjectVectors,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) Remove(
	ctx context.Context,
	req *client.ObjectID,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			res, err := proto.NewNGTDClient(conn).Remove(ctx, objectIDToNGTDRemoveRequest(req), copts...)
			if err != nil {
				return nil, err
			}

			if len(res.GetError()) != 0 {
				return nil, errors.New(res.GetError())
			}
			return res, nil
		},
	)
	return err
}

func (c *ngtdClient) StreamRemove(
	ctx context.Context,
	dataProvider func() *client.ObjectID,
	f func(error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			st, err := proto.NewNGTDClient(conn).StreamRemove(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, grpc.BidirectionalStreamClient(st,
				func() interface{} {
					if d := dataProvider(); d != nil {
						return objectIDToNGTDRemoveRequest(d)
					}
					return nil
				}, func() interface{} {
					return new(proto.RemoveResponse)
				}, func(intr interface{}, err error) {
					if err != nil {
						f(err)
						return
					}

					res := intr.(*proto.RemoveResponse)
					if len(res.GetError()) != 0 {
						f(errors.New(res.GetError()))
						return
					}

					f(err)
				})
		},
	)
	return err
}

func (c *ngtdClient) MultiRemove(
	ctx context.Context,
	req *client.ObjectIDs,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) GetObject(
	ctx context.Context,
	req *client.ObjectID,
) (*client.ObjectVector, error) {
	resp, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			res, err := proto.NewNGTDClient(conn).GetObject(ctx, objectIDToNGTDGetObjectRequest(req), copts...)
			if err != nil {
				return nil, err
			}

			if len(res.GetError()) != 0 {
				return nil, errors.New(res.GetError())
			}
			return res, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return ngtdGetObjectResponseToObjectVector(resp.(*proto.GetObjectResponse)), err
}

func (c *ngtdClient) StreamGetObject(
	ctx context.Context,
	dataProvider func() *client.ObjectID,
	f func(*client.ObjectVector, error),
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			st, err := proto.NewNGTDClient(conn).StreamGetObject(ctx, copts...)
			if err != nil {
				return nil, err
			}

			return nil, grpc.BidirectionalStreamClient(st,
				func() interface{} {
					if d := dataProvider(); d != nil {
						return objectIDToNGTDGetObjectRequest(d)
					}
					return nil
				}, func() interface{} {
					return new(proto.InsertResponse)
				}, func(intr interface{}, err error) {
					if err != nil {
						f(nil, err)
					}

					res := intr.(*proto.GetObjectResponse)
					if len(res.GetError()) != 0 {
						f(nil, errors.New(res.GetError()))
					}

					f(ngtdGetObjectResponseToObjectVector(res), err)
				})
		},
	)
	return err
}

func (c *ngtdClient) CreateIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return proto.NewNGTDClient(conn).CreateIndex(ctx, controlCreateIndexRequestToCreateIndexRequest(req), copts...)
		},
	)
	return err
}

func (c *ngtdClient) SaveIndex(ctx context.Context) error {
	_, err := c.Client.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return proto.NewNGTDClient(conn).SaveIndex(ctx, new(proto.Empty), copts...)
		},
	)
	return err
}

func (c *ngtdClient) CreateAndSaveIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *ngtdClient) IndexInfo(ctx context.Context) (*client.InfoIndex, error) {
	return nil, errors.ErrUnsupportedClientMethod
}

func searchRequestToNgtdSearchRequest(in *client.SearchRequest) *proto.SearchRequest {
	size, epsilon := getSizeAndEpsilon(in.GetConfig())
	return &proto.SearchRequest{
		Vector:  tofloat64(in.GetVector()),
		Size_:   size,
		Epsilon: epsilon,
	}
}

func searchIDRequestToNgtdSearchRequest(in *client.SearchIDRequest) *proto.SearchRequest {
	size, epsilon := getSizeAndEpsilon(in.GetConfig())
	return &proto.SearchRequest{
		Id:      []byte(in.GetId()),
		Size_:   size,
		Epsilon: epsilon,
	}
}

func ngtdSearchResponseToSearchResponse(in *proto.SearchResponse) *client.SearchResponse {
	if len(in.GetError()) != 0 {
		return nil
	}

	results := make([]*client.ObjectDistance, len(in.GetResult()))

	for _, r := range in.GetResult() {
		if len(r.GetError()) == 0 {
			results = append(results, &client.ObjectDistance{
				Id:       string(r.GetId()),
				Distance: r.GetDistance(),
			})
		}
	}

	return &client.SearchResponse{
		Results: results,
	}
}

func ngtdGetObjectResponseToObjectVector(in *proto.GetObjectResponse) *client.ObjectVector {
	if len(in.GetError()) != 0 {
		return nil
	}

	return &client.ObjectVector{
		Id:     string(in.GetId()),
		Vector: in.GetVector(),
	}
}

func objectVectorToNGTDInsertRequest(in *client.ObjectVector) *proto.InsertRequest {
	return &proto.InsertRequest{
		Id:     []byte(in.GetId()),
		Vector: tofloat64(in.GetVector()),
	}
}

func objectIDToNGTDRemoveRequest(in *client.ObjectID) *proto.RemoveRequest {
	return &proto.RemoveRequest{
		Id: []byte(in.GetId()),
	}
}

func objectIDToNGTDGetObjectRequest(in *client.ObjectID) *proto.GetObjectRequest {
	return &proto.GetObjectRequest{
		Id: []byte(in.GetId()),
	}
}

func controlCreateIndexRequestToCreateIndexRequest(in *client.ControlCreateIndexRequest) *proto.CreateIndexRequest {
	return &proto.CreateIndexRequest{
		PoolSize: in.GetPoolSize(),
	}
}

func getSizeAndEpsilon(cfg *client.SearchConfig) (size int32, epsilon float32) {
	if cfg != nil {
		size = int32(cfg.GetNum())
		epsilon = float32(cfg.GetEpsilon())
	}
	return
}

func tofloat64(in []float32) (out []float64) {
	out = make([]float64, len(in))
	for i := range in {
		out[i] = float64(in[i])
	}
	return
}
