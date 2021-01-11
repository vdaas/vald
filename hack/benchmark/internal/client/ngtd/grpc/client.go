//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package grpc provides grpc client functions
package grpc

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client"
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
	c    grpc.Client
	opts []grpc.Option
}

func New(ctx context.Context, opts ...Option) (Client, error) {
	c := new(ngtdClient)

	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	c.c = grpc.New(c.opts...)

	if _, err := c.c.Connect(ctx, c.addr); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *ngtdClient) Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (oid *payload.Object_ID, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		id, err := proto.NewNGTDClient(conn).GetObject(ctx, &proto.GetObjectRequest{
			Id: []byte(in.GetId()),
		})
		if err != nil {
			return nil, err
		}
		if len(id.GetError()) != 0 {
			return nil, errors.New(id.GetError())
		}
		oid = &payload.Object_ID{
			Id: string(id.GetId()),
		}
		return oid, nil
	})
	if err != nil {
		return nil, err
	}
	return oid, nil
}

func (c *ngtdClient) Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	res, err := c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		r, err := proto.NewNGTDClient(conn).Search(ctx, searchRequestToNgtdSearchRequest(in), copts...)
		if err != nil {
			return nil, err
		}
		if len(r.GetError()) != 0 {
			return nil, errors.New(r.GetError())
		}
		return r, nil
	})
	if err != nil {
		return nil, err
	}
	return ngtdSearchResponseToSearchResponse(res.(*proto.SearchResponse)), nil
}

func (c *ngtdClient) SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	res, err := c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		r, err := proto.NewNGTDClient(conn).SearchByID(ctx, searchIDRequestToNgtdSearchRequest(in), copts...)
		if err != nil {
			return nil, err
		}
		if len(r.GetError()) != 0 {
			return nil, errors.New(r.GetError())
		}
		return r, nil
	})
	if err != nil {
		return nil, err
	}
	return ngtdSearchResponseToSearchResponse(res.(*proto.SearchResponse)), nil
}

func (c *ngtdClient) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		st, err := proto.NewNGTDClient(conn).StreamSearch(ctx, copts...)
		if err != nil {
			return nil, err
		}
		res = NewStreamSearchClient(st)
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchByIDClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		st, err := proto.NewNGTDClient(conn).StreamSearch(ctx, copts...)
		if err != nil {
			return nil, err
		}
		res = NewStreamSearchByIDClient(st)
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) MultiSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res = &payload.Search_Responses{
			Responses: make([]*payload.Search_Response, 0, len(in.GetRequests())),
		}
		for _, req := range in.GetRequests() {
			sres, err := c.Search(ctx, req, opts...)
			if err != nil {
				return nil, err
			}
			res.Responses = append(res.Responses, sres)
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) MultiSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res = &payload.Search_Responses{
			Responses: make([]*payload.Search_Response, 0, len(in.GetRequests())),
		}
		for _, req := range in.GetRequests() {
			sres, err := c.SearchByID(ctx, req, opts...)
			if err != nil {
				return nil, err
			}
			res.Responses = append(res.Responses, sres)
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		data, err := proto.NewNGTDClient(conn).Insert(ctx, &proto.InsertRequest{
			Id:     []byte(in.GetVector().GetId()),
			Vector: tofloat64(in.GetVector().GetVector()),
		}, copts...)
		if err != nil {
			return nil, err
		}
		if len(data.GetError()) != 0 {
			return nil, errors.New(data.GetError())
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Insert_StreamInsertClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		st, err := proto.NewNGTDClient(conn).StreamInsert(ctx, copts...)
		if err != nil {
			return nil, err
		}
		res = NewStreamInsertClient(st)
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) MultiInsert(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res = &payload.Object_Locations{
			Locations: make([]*payload.Object_Location, 0, len(in.GetRequests())),
		}
		for _, req := range in.GetRequests() {
			sres, err := c.Insert(ctx, req, opts...)
			if err != nil {
				return nil, err
			}
			res.Locations = append(res.Locations, sres)
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	_, err = c.Remove(ctx, &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: in.GetVector().GetId(),
		},
	}, opts...)
	if err != nil {
		return nil, err
	}
	_, err = c.Insert(ctx, &payload.Insert_Request{
		Vector: in.GetVector(),
	}, opts...)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *ngtdClient) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (res vald.Update_StreamUpdateClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		ist, err := proto.NewNGTDClient(conn).StreamInsert(ctx, copts...)
		if err != nil {
			return nil, err
		}
		rst, err := proto.NewNGTDClient(conn).StreamRemove(ctx, copts...)
		if err != nil {
			return nil, err
		}
		res = NewStreamUpdateClient(ist, rst)
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) MultiUpdate(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res = &payload.Object_Locations{
			Locations: make([]*payload.Object_Location, 0, len(in.GetRequests())),
		}
		for _, req := range in.GetRequests() {
			sres, err := c.Update(ctx, req, opts...)
			if err != nil {
				return nil, err
			}
			res.Locations = append(res.Locations, sres)
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	id, err := c.Exists(ctx, &payload.Object_ID{
		Id: in.GetVector().GetId(),
	}, opts...)
	if err == nil || len(id.GetId()) != 0 {
		return c.Update(ctx, &payload.Update_Request{
			Vector: in.GetVector(),
		}, opts...)
	}
	return c.Insert(ctx, &payload.Insert_Request{
		Vector: in.GetVector(),
	}, opts...)
}

func (c *ngtdClient) StreamUpsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Upsert_StreamUpsertClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		st, err := proto.NewNGTDClient(conn).StreamInsert(ctx, copts...)
		if err != nil {
			return nil, err
		}
		res = NewStreamUpsertClient(c, st)
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) MultiUpsert(ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res = &payload.Object_Locations{
			Locations: make([]*payload.Object_Location, 0, len(in.GetRequests())),
		}
		for _, req := range in.GetRequests() {
			sres, err := c.Upsert(ctx, req, opts...)
			if err != nil {
				return nil, err
			}
			res.Locations = append(res.Locations, sres)
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err := proto.NewNGTDClient(conn).Remove(ctx, &proto.RemoveRequest{
			Id: []byte(in.GetId().GetId()),
		}, copts...)
		if err != nil {
			return nil, err
		}

		if len(res.GetError()) != 0 {
			return nil, errors.New(res.GetError())
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (res vald.Remove_StreamRemoveClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		st, err := proto.NewNGTDClient(conn).StreamRemove(ctx, copts...)
		if err != nil {
			return nil, err
		}
		res = NewStreamRemoveClient(st)
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) MultiRemove(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	_, err = c.c.Do(ctx, c.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		for _, req := range in.GetRequests() {
			id, err := proto.NewNGTDClient(conn).Remove(ctx, &proto.RemoveRequest{
				Id: []byte(req.GetId().GetId()),
			}, append(copts, opts...)...)
			if err != nil {
				return nil, err
			}
			if len(id.GetError()) != 0 {
				return nil, errors.New(id.GetError())
			}

		}
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) GetObject(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	res, err := c.c.Do(ctx, c.addr, func(
		ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		res, err := proto.NewNGTDClient(conn).GetObject(ctx, &proto.GetObjectRequest{
			Id: []byte(in.GetId()),
		}, copts...)
		if err != nil {
			return nil, err
		}

		if len(res.GetError()) != 0 {
			return nil, errors.New(res.GetError())
		}
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	r, ok := res.(*proto.GetObjectResponse)
	if !ok {
		return nil, errors.ErrInvalidAPIConfig
	}
	if len(r.GetError()) != 0 {
		return nil, errors.New(r.GetError())
	}
	return &client.ObjectVector{
		Id:     string(r.GetId()),
		Vector: r.GetVector(),
	}, nil
}

func (c *ngtdClient) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (res vald.Object_StreamGetObjectClient, err error) {
	_, err = c.c.Do(ctx, c.addr, func(
		ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) (interface{}, error) {
		st, err := proto.NewNGTDClient(conn).StreamGetObject(ctx, copts...)
		if err != nil {
			return nil, err
		}
		res = NewStreamObjectClient(st)
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ngtdClient) CreateIndex(
	ctx context.Context,
	in *client.ControlCreateIndexRequest,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	_, err := c.c.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return proto.NewNGTDClient(conn).CreateIndex(ctx, &proto.CreateIndexRequest{
				PoolSize: in.GetPoolSize(),
			}, copts...)
		},
	)
	return nil, err
}

func (c *ngtdClient) SaveIndex(
	ctx context.Context,
	req *client.Empty,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	_, err := c.c.Do(ctx, c.addr,
		func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
			return proto.NewNGTDClient(conn).SaveIndex(ctx, new(proto.Empty), copts...)
		},
	)
	return nil, err
}

func (c *ngtdClient) CreateAndSaveIndex(
	ctx context.Context,
	req *client.ControlCreateIndexRequest,
	opts ...grpc.CallOption,
) (*client.Empty, error) {
	_, err := c.CreateIndex(ctx, req)
	if err != nil {
		return nil, err
	}
	_, err = c.SaveIndex(ctx, nil)
	return nil, err
}

func (c *ngtdClient) IndexInfo(
	ctx context.Context,
	req *client.Empty,
	opts ...grpc.CallOption,
) (res *client.InfoIndexCount, err error) {
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
