//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package vald provides vald grpc client library
package vald

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

type Client interface {
	vald.Client
	GRPCClient() grpc.Client
	Start(context.Context) (<-chan error, error)
	Stop(context.Context) error
}

type client struct {
	addrs []string
	c     grpc.Client
}

type singleClient struct {
	vc vald.Client
}

const (
	apiName = "vald/internal/client/v1/client/vald"
)

func New(opts ...Option) (Client, error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		err := opt(c)
		if err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	if c.c == nil {
		if c.addrs == nil {
			return nil, errors.ErrGRPCTargetAddrNotFound
		}
		c.c = grpc.New(grpc.WithAddrs(c.addrs...))
	}
	return c, nil
}

func NewValdClient(cc *grpc.ClientConn) Client {
	return &singleClient{vc: vald.NewValdClient(cc)}
}

func (c *client) Start(ctx context.Context) (<-chan error, error) {
	return c.c.StartConnectionMonitor(ctx)
}

func (c *client) Stop(ctx context.Context) error {
	return c.c.Close(ctx)
}

func (c *client) GRPCClient() grpc.Client {
	return c.c
}

func (c *client) Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (oid *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Exists")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		oid, err = vald.NewValdClient(conn).Exists(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return oid, nil
}

func (c *client) Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).Search(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.SearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).SearchByID(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.StreamSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamSearch(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchByIDClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.StreamSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamSearchByID(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.MultiSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiSearch(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.MultiSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiSearchByID(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) LinearSearch(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.LinearSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).LinearSearch(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) LinearSearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.LinearSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).LinearSearchByID(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamLinearSearch(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamLinearSearchClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.LinearStreamSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamLinearSearch(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamLinearSearchByID(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamLinearSearchByIDClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.StreamLinearSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamLinearSearchByID(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiLinearSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.MultiLinearSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiLinearSearch(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiLinearSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.MultiLinearSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiLinearSearchByID(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).Insert(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Insert_StreamInsertClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.StreamInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamInsert(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiInsert(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiInsert(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Update")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).Update(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (res vald.Update_StreamUpdateClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.StreamUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamUpdate(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiUpdate(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.MultiUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiUpdate(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).Upsert(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamUpsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Upsert_StreamUpsertClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.StreamUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamUpsert(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiUpsert(ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.MultiUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiUpsert(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).Remove(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (res vald.Remove_StreamRemoveClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.StreamRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamRemove(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiRemove(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).MultiRemove(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) GetObject(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.GetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).GetObject(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (res vald.Object_StreamGetObjectClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.StreamGetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.c.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (interface{}, error) {
		res, err = vald.NewValdClient(conn).StreamGetObject(ctx, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (*singleClient) Start(_ context.Context) (<-chan error, error) {
	return nil, nil
}

func (*singleClient) Stop(_ context.Context) error {
	return nil
}

func (*singleClient) GRPCClient() grpc.Client {
	return nil
}

func (c *singleClient) Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (oid *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.Exists")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Exists(ctx, in, opts...)
}

func (c *singleClient) Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.Search")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Search(ctx, in, opts...)
}

func (c *singleClient) SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.SearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.SearchByID(ctx, in, opts...)
}

func (c *singleClient) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.StreamSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamSearch(ctx, opts...)
}

func (c *singleClient) StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchByIDClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.StreamSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamSearchByID(ctx, opts...)
}

func (c *singleClient) MultiSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.MultiSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiSearch(ctx, in, opts...)
}

func (c *singleClient) MultiSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.MultiSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiSearchByID(ctx, in, opts...)
}

func (c *singleClient) LinearSearch(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.LinearSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.LinearSearch(ctx, in, opts...)
}

func (c *singleClient) LinearSearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.LinearSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.LinearSearchByID(ctx, in, opts...)
}

func (c *singleClient) StreamLinearSearch(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamLinearSearchClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.StreamLinearSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamLinearSearch(ctx, opts...)
}

func (c *singleClient) StreamLinearSearchByID(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamLinearSearchByIDClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.StreamLinearSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamLinearSearchByID(ctx, opts...)
}

func (c *singleClient) MultiLinearSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.MultiLinearSearch")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiLinearSearch(ctx, in, opts...)
}

func (c *singleClient) MultiLinearSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.MultiLinearSearchByID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiLinearSearchByID(ctx, in, opts...)
}

func (c *singleClient) Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.Insert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Insert(ctx, in, opts...)
}

func (c *singleClient) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Insert_StreamInsertClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.StreamInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamInsert(ctx, opts...)
}

func (c *singleClient) MultiInsert(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.MultiInsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiInsert(ctx, in, opts...)
}

func (c *singleClient) Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.Update")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Update(ctx, in, opts...)
}

func (c *singleClient) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (res vald.Update_StreamUpdateClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.StreamUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamUpdate(ctx, opts...)
}

func (c *singleClient) MultiUpdate(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.MultiUpdate")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiUpdate(ctx, in, opts...)
}

func (c *singleClient) Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.Upsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Upsert(ctx, in, opts...)
}

func (c *singleClient) StreamUpsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Upsert_StreamUpsertClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.StreamUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamUpsert(ctx, opts...)
}

func (c *singleClient) MultiUpsert(ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.MultiUpsert")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiUpsert(ctx, in, opts...)
}

func (c *singleClient) Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Remove(ctx, in, opts...)
}

func (c *singleClient) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (res vald.Remove_StreamRemoveClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.StreamRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamRemove(ctx, opts...)
}

func (c *singleClient) MultiRemove(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.MultiRemove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiRemove(ctx, in, opts...)
}

func (c *singleClient) GetObject(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.GetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.GetObject(ctx, in, opts...)
}

func (c *singleClient) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (res vald.Object_StreamGetObjectClient, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.StreamGetObject")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamGetObject(ctx, opts...)
}
