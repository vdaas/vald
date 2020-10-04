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

// Package vald provides vald grpc client library
package vald

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
)

type Client vald.Client

type client struct {
	addr string
	c    grpc.Client
}

func New(opts ...Option) Client {
	c := new(client)
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *client) Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (oid *payload.Object_ID, err error) {
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

func (c *client) Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
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

func (c *client) SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (res *payload.Search_Response, err error) {
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

func (c *client) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchClient, err error) {
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

func (c *client) StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (res vald.Search_StreamSearchByIDClient, err error) {
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

func (c *client) MultiSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
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

func (c *client) MultiSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (res *payload.Search_Responses, err error) {
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

func (c *client) Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
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

func (c *client) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Insert_StreamInsertClient, err error) {
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

func (c *client) MultiInsert(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
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

func (c *client) Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
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

func (c *client) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (res vald.Update_StreamUpdateClient, err error) {
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

func (c *client) MultiUpdate(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
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

func (c *client) Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
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

func (c *client) StreamUpsert(ctx context.Context, opts ...grpc.CallOption) (res vald.Upsert_StreamUpsertClient, err error) {
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

func (c *client) MultiUpsert(ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
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

func (c *client) Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (res *payload.Object_Location, err error) {
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

func (c *client) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (res vald.Remove_StreamRemoveClient, err error) {
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

func (c *client) MultiRemove(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (res *payload.Object_Locations, err error) {
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

func (c *client) GetObject(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (res *payload.Object_Vector, err error) {
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

func (c *client) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (res vald.Object_StreamGetObjectClient, err error) {
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
