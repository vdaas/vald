//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
		c.c = grpc.New("Vald Client", grpc.WithAddrs(c.addrs...))
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

func (c *client) Exists(
	ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption,
) (oid *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.ExistsRPCName), apiName+"/"+vald.ExistsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	oid, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_ID, error) {
		return vald.NewValdClient(conn).Exists(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return oid, nil
}

func (c *client) Search(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.SearchRPCName), apiName+"/"+vald.SearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Response, error) {
		return vald.NewValdClient(conn).Search(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) SearchWithMetadata(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.SearchWithMetadataRPCName), apiName+"/"+vald.SearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Response, error) {
		return vald.NewValdClient(conn).SearchWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) SearchByID(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.SearchByIDRPCName), apiName+"/"+vald.SearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Response, error) {
		return vald.NewValdClient(conn).SearchByID(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) SearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.SearchByIDWithMetadataRPCName), apiName+"/"+vald.SearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Response, error) {
		return vald.NewValdClient(conn).SearchByIDWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamSearch(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Search_StreamSearchClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamSearchRPCName), apiName+"/"+vald.StreamSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.Search_StreamSearchClient, error) {
		return vald.NewValdClient(conn).StreamSearch(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamSearchWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.SearchWithMetadata_StreamSearchWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamSearchWithMetadataRPCName), apiName+"/"+vald.StreamSearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.SearchWithMetadata_StreamSearchWithMetadataClient, error) {
		return vald.NewValdClient(conn).StreamSearchWithMetadata(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamSearchByID(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Search_StreamSearchByIDClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamSearchByIDRPCName), apiName+"/"+vald.StreamSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.Search_StreamSearchByIDClient, error) {
		return vald.NewValdClient(conn).StreamSearchByID(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamSearchByIDWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.SearchWithMetadata_StreamSearchByIDWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamSearchByIDWithMetadataRPCName), apiName+"/"+vald.StreamSearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.SearchWithMetadata_StreamSearchByIDWithMetadataClient, error) {
		return vald.NewValdClient(conn).StreamSearchByIDWithMetadata(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiSearch(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiSearchRPCName), apiName+"/"+vald.MultiSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Responses, error) {
		return vald.NewValdClient(conn).MultiSearch(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiSearchWithMetadata(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiSearchWithMetadataRPCName), apiName+"/"+vald.MultiSearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Responses, error) {
		return vald.NewValdClient(conn).MultiSearchWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiSearchByID(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiSearchByIDRPCName), apiName+"/"+vald.MultiSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Responses, error) {
		return vald.NewValdClient(conn).MultiSearchByID(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiSearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiSearchByIDWithMetadataRPCName), apiName+"/"+vald.MultiSearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Responses, error) {
		return vald.NewValdClient(conn).MultiSearchByIDWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) LinearSearch(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.LinearSearchRPCName), apiName+"/"+vald.LinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Response, error) {
		return vald.NewValdClient(conn).LinearSearch(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) LinearSearchWithMetadata(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.LinearSearchWithMetadataRPCName), apiName+"/"+vald.LinearSearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Response, error) {
		return vald.NewValdClient(conn).LinearSearchWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) LinearSearchByID(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.LinearSearchByIDRPCName), apiName+"/"+vald.LinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Response, error) {
		return vald.NewValdClient(conn).LinearSearchByID(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) LinearSearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.LinearSearchByIDWithMetadataRPCName), apiName+"/"+vald.LinearSearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Response, error) {
		return vald.NewValdClient(conn).LinearSearchByIDWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamLinearSearch(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Search_StreamLinearSearchClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamLinearSearchRPCName), apiName+"/"+vald.StreamLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.Search_StreamLinearSearchClient, error) {
		return vald.NewValdClient(conn).StreamLinearSearch(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamLinearSearchWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.SearchWithMetadata_StreamLinearSearchWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamLinearSearchWithMetadataRPCName), apiName+"/"+vald.StreamLinearSearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.SearchWithMetadata_StreamLinearSearchWithMetadataClient, error) {
		return vald.NewValdClient(conn).StreamLinearSearchWithMetadata(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamLinearSearchByID(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Search_StreamLinearSearchByIDClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamLinearSearchByIDRPCName), apiName+"/"+vald.StreamLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.Search_StreamSearchByIDClient, error) {
		return vald.NewValdClient(conn).StreamLinearSearchByID(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamLinearSearchByIDWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.SearchWithMetadata_StreamLinearSearchByIDWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamLinearSearchByIDWithMetadataRPCName), apiName+"/"+vald.StreamLinearSearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.SearchWithMetadata_StreamSearchByIDWithMetadataClient, error) {
		return vald.NewValdClient(conn).StreamLinearSearchByIDWithMetadata(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiLinearSearch(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiLinearSearchRPCName), apiName+"/"+vald.MultiLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Responses, error) {
		return vald.NewValdClient(conn).MultiLinearSearch(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiLinearSearchWithMetadata(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiLinearSearchWithMetadataRPCName), apiName+"/"+vald.MultiLinearSearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Responses, error) {
		return vald.NewValdClient(conn).MultiLinearSearchWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiLinearSearchByID(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiLinearSearchByIDRPCName), apiName+"/"+vald.MultiLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Responses, error) {
		return vald.NewValdClient(conn).MultiLinearSearchByID(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiLinearSearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiLinearSearchByIDWithMetadataRPCName), apiName+"/"+vald.MultiLinearSearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Search_Responses, error) {
		return vald.NewValdClient(conn).MultiLinearSearchByIDWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Insert(
	ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.InsertRPCName), apiName+"/"+vald.InsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Location, error) {
		return vald.NewValdClient(conn).Insert(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) InsertWithMetadata(
	ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.InsertWithMetadataRPCName), apiName+"/"+vald.InsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Location, error) {
		return vald.NewValdClient(conn).InsertWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamInsert(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Insert_StreamInsertClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamInsertRPCName), apiName+"/"+vald.StreamInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.Insert_StreamInsertClient, error) {
		return vald.NewValdClient(conn).StreamInsert(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamInsertWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.InsertWithMetadata_StreamInsertWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamInsertWithMetadataRPCName), apiName+"/"+vald.StreamInsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.InsertWithMetadata_StreamInsertWithMetadataClient, error) {
		return vald.NewValdClient(conn).StreamInsertWithMetadata(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiInsert(
	ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiInsertRPCName), apiName+"/"+vald.MultiInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Locations, error) {
		return vald.NewValdClient(conn).MultiInsert(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiInsertWithMetadata(
	ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiInsertWithMetadataRPCName), apiName+"/"+vald.MultiInsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Locations, error) {
		return vald.NewValdClient(conn).MultiInsertWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Update(
	ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.UpdateRPCName), apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Location, error) {
		return vald.NewValdClient(conn).Update(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) UpdateWithMetadata(
	ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.UpdateWithMetadataRPCName), apiName+"/"+vald.UpdateWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Location, error) {
		return vald.NewValdClient(conn).UpdateWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamUpdate(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Update_StreamUpdateClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamUpdateRPCName), apiName+"/"+vald.StreamUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.Update_StreamUpdateClient, error) {
		return vald.NewValdClient(conn).StreamUpdate(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamUpdateWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.UpdateWithMetadata_StreamUpdateWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamUpdateWithMetadataRPCName), apiName+"/"+vald.StreamUpdateWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.UpdateWithMetadata_StreamUpdateWithMetadataClient, error) {
		return vald.NewValdClient(conn).StreamUpdateWithMetadata(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiUpdate(
	ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiUpdateRPCName), apiName+"/"+vald.MultiUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Locations, error) {
		return vald.NewValdClient(conn).MultiUpdate(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiUpdateWithMetadata(
	ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiUpdateWithMetadataRPCName), apiName+"/"+vald.MultiUpdateWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Locations, error) {
		return vald.NewValdClient(conn).MultiUpdateWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) UpdateTimestamp(
	ctx context.Context, in *payload.Update_TimestampRequest, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.UpdateTimestampRPCName), apiName+"/"+vald.UpdateTimestampRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Location, error) {
		return vald.NewValdClient(conn).UpdateTimestamp(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) UpdateTimestampWithMetadata(
	ctx context.Context, in *payload.Update_TimestampRequest, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.UpdateTimestampWithMetadataRPCName), apiName+"/"+vald.UpdateTimestampWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Location, error) {
		return vald.NewValdClient(conn).UpdateTimestampWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Upsert(
	ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.UpsertRPCName), apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Location, error) {
		return vald.NewValdClient(conn).Upsert(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) UpsertWithMetadata(
	ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.UpsertWithMetadataRPCName), apiName+"/"+vald.UpsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Location, error) {
		return vald.NewValdClient(conn).UpsertWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamUpsert(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Upsert_StreamUpsertClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamUpsertRPCName), apiName+"/"+vald.StreamUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.Upsert_StreamUpsertClient, error) {
		return vald.NewValdClient(conn).StreamUpsert(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamUpsertWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.UpsertWithMetadata_StreamUpsertWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamUpsertWithMetadataRPCName), apiName+"/"+vald.StreamUpsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.UpsertWithMetadata_StreamUpsertWithMetadataClient, error) {
		return vald.NewValdClient(conn).StreamUpsertWithMetadata(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiUpsert(
	ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiUpsertRPCName), apiName+"/"+vald.MultiUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Locations, error) {
		return vald.NewValdClient(conn).MultiUpsert(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiUpsertWithMetadata(
	ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiUpsertWithMetadataRPCName), apiName+"/"+vald.MultiUpsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Locations, error) {
		return vald.NewValdClient(conn).MultiUpsertWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Remove(
	ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Location, error) {
		return vald.NewValdClient(conn).Remove(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) RemoveWithMetadata(
	ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.RemoveWithMetadataRPCName), apiName+"/"+vald.RemoveWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Location, error) {
		return vald.NewValdClient(conn).RemoveWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamRemove(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Remove_StreamRemoveClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamRemoveRPCName), apiName+"/"+vald.StreamRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.Remove_StreamRemoveClient, error) {
		return vald.NewValdClient(conn).StreamRemove(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamRemoveWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.RemoveWithMetadata_StreamRemoveWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamRemoveWithMetadataRPCName), apiName+"/"+vald.StreamRemoveWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.RemoveWithMetadata_StreamRemoveWithMetadataClient, error) {
		return vald.NewValdClient(conn).StreamRemoveWithMetadata(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiRemove(
	ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Locations, error) {
		return vald.NewValdClient(conn).MultiRemove(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) MultiRemoveWithMetadata(
	ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.MultiRemoveWithMetadataRPCName), apiName+"/"+vald.MultiRemoveWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Locations, error) {
		return vald.NewValdClient(conn).MultiRemoveWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) Flush(
	ctx context.Context, in *payload.Flush_Request, opts ...grpc.CallOption,
) (res *payload.Info_Index_Count, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.FlushRPCName), apiName+"/"+vald.FlushRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Info_Index_Count, error) {
		return vald.NewValdClient(conn).Flush(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) RemoveByTimestamp(
	ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.RemoveByTimestampRPCName), apiName+"/"+vald.RemoveByTimestampRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Locations, error) {
		return vald.NewValdClient(conn).RemoveByTimestamp(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) RemoveByTimestampWithMetadata(
	ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.RemoveByTimestampWithMetadataRPCName), apiName+"/"+vald.RemoveByTimestampWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Locations, error) {
		return vald.NewValdClient(conn).RemoveByTimestampWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) GetObject(
	ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption,
) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Vector, error) {
		return vald.NewValdClient(conn).GetObject(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) GetObjectWithMetadata(
	ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption,
) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.GetObjectWithMetadataRPCName), apiName+"/"+vald.GetObjectWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Vector, error) {
		return vald.NewValdClient(conn).GetObjectWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamGetObject(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Object_StreamGetObjectClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamGetObjectRPCName), apiName+"/"+vald.StreamGetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.Object_StreamGetObjectClient, error) {
		return vald.NewValdClient(conn).StreamGetObject(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamGetObjectWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.ObjectWithMetadata_StreamGetObjectWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamGetObjectWithMetadataRPCName), apiName+"/"+vald.StreamGetObjectWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.ObjectWithMetadata_StreamGetObjectWithMetadataClient, error) {
		return vald.NewValdClient(conn).StreamGetObjectWithMetadata(ctx, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamListObject(
	ctx context.Context, in *payload.Object_List_Request, opts ...grpc.CallOption,
) (res vald.Object_StreamListObjectClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamListObjectRPCName), apiName+"/"+vald.StreamListObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.Object_StreamListObjectClient, error) {
		return vald.NewValdClient(conn).StreamListObject(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) StreamListObjectWithMetadata(
	ctx context.Context, in *payload.Object_List_Request, opts ...grpc.CallOption,
) (res vald.ObjectWithMetadata_StreamListObjectWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.StreamListObjectWithMetadataRPCName), apiName+"/"+vald.StreamListObjectWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (vald.ObjectWithMetadata_StreamListObjectWithMetadataClient, error) {
		return vald.NewValdClient(conn).StreamListObjectWithMetadata(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) IndexInfo(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (res *payload.Info_Index_Count, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.IndexInfoRPCName), apiName+"/"+vald.IndexInfoRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Info_Index_Count, error) {
		return vald.NewValdClient(conn).IndexInfo(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) IndexDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (res *payload.Info_Index_Detail, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.IndexDetailRPCName), apiName+"/"+vald.IndexDetailRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Info_Index_Detail, error) {
		return vald.NewValdClient(conn).IndexDetail(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) IndexStatistics(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (res *payload.Info_Index_Statistics, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.IndexStatisticsRPCName), apiName+"/"+vald.IndexStatisticsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Info_Index_Statistics, error) {
		return vald.NewValdClient(conn).IndexStatistics(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) IndexStatisticsDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (res *payload.Info_Index_StatisticsDetail, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.IndexStatisticsDetailRPCName), apiName+"/"+vald.IndexStatisticsDetailRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Info_Index_StatisticsDetail, error) {
		return vald.NewValdClient(conn).IndexStatisticsDetail(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) IndexProperty(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (res *payload.Info_Index_PropertyDetail, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.IndexPropertyRPCName), apiName+"/"+vald.IndexPropertyRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Info_Index_PropertyDetail, error) {
		return vald.NewValdClient(conn).IndexProperty(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) GetTimestamp(
	ctx context.Context, in *payload.Object_TimestampRequest, opts ...grpc.CallOption,
) (res *payload.Object_Timestamp, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.GetTimestampRPCName), apiName+"/"+vald.GetTimestampRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Timestamp, error) {
		return vald.NewValdClient(conn).GetTimestamp(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (*singleClient) Start(context.Context) (<-chan error, error) {
	return nil, nil
}

func (*singleClient) Stop(context.Context) error {
	return nil
}

func (*singleClient) GRPCClient() grpc.Client {
	return nil
}

func (c *singleClient) Exists(
	ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption,
) (oid *payload.Object_ID, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.ExistsRPCName), apiName+"/"+vald.ExistsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Exists(ctx, in, opts...)
}

func (c *singleClient) Search(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.SearchRPCName), apiName+"/"+vald.SearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Search(ctx, in, opts...)
}

func (c *singleClient) SearchWithMetadata(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.SearchWithMetadataRPCName), apiName+"/"+vald.SearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.SearchWithMetadata(ctx, in, opts...)
}

func (c *singleClient) SearchByID(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.SearchByIDRPCName), apiName+"/"+vald.SearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.SearchByID(ctx, in, opts...)
}

func (c *singleClient) SearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.SearchByIDWithMetadataRPCName), apiName+"/"+vald.SearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.SearchByIDWithMetadata(ctx, in, opts...)
}

func (c *singleClient) StreamSearch(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Search_StreamSearchClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamSearchRPCName), apiName+"/"+vald.StreamSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamSearch(ctx, opts...)
}

func (c *singleClient) StreamSearchWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.SearchWithMetadata_StreamSearchWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamSearchWithMetadataRPCName), apiName+"/"+vald.StreamSearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamSearchWithMetadata(ctx, opts...)
}

func (c *singleClient) StreamSearchByID(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Search_StreamSearchByIDClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamSearchByIDRPCName), apiName+"/"+vald.StreamSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamSearchByID(ctx, opts...)
}

func (c *singleClient) StreamSearchByIDWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.SearchWithMetadata_StreamSearchByIDWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamSearchByIDWithMetadataRPCName), apiName+"/"+vald.StreamSearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamSearchByIDWithMetadata(ctx, opts...)
}

func (c *singleClient) MultiSearch(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiSearchRPCName), apiName+"/"+vald.MultiSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiSearch(ctx, in, opts...)
}

func (c *singleClient) MultiSearchWithMetadata(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiSearchWithMetadataRPCName), apiName+"/"+vald.MultiSearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiSearchWithMetadata(ctx, in, opts...)
}

func (c *singleClient) MultiSearchByID(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiSearchByIDRPCName), apiName+"/"+vald.MultiSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiSearchByID(ctx, in, opts...)
}

func (c *singleClient) MultiSearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiSearchByIDWithMetadataRPCName), apiName+"/"+vald.MultiSearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiSearchByIDWithMetadata(ctx, in, opts...)
}

func (c *singleClient) LinearSearch(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.LinearSearchRPCName), apiName+"/"+vald.LinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.LinearSearch(ctx, in, opts...)
}

func (c *singleClient) LinearSearchWithMetadata(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.LinearSearchWithMetadataRPCName), apiName+"/"+vald.LinearSearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.LinearSearchWithMetadata(ctx, in, opts...)
}

func (c *singleClient) LinearSearchByID(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.LinearSearchByIDRPCName), apiName+"/"+vald.LinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.LinearSearchByID(ctx, in, opts...)
}

func (c *singleClient) LinearSearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Response, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.LinearSearchByIDWithMetadataRPCName), apiName+"/"+vald.LinearSearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.LinearSearchByIDWithMetadata(ctx, in, opts...)
}

func (c *singleClient) StreamLinearSearch(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Search_StreamLinearSearchClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamLinearSearchRPCName), apiName+"/"+vald.StreamLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamLinearSearch(ctx, opts...)
}

func (c *singleClient) StreamLinearSearchWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.SearchWithMetadata_StreamLinearSearchWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamLinearSearchWithMetadataRPCName), apiName+"/"+vald.StreamLinearSearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamLinearSearchWithMetadata(ctx, opts...)
}

func (c *singleClient) StreamLinearSearchByID(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Search_StreamLinearSearchByIDClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamLinearSearchByIDRPCName), apiName+"/"+vald.StreamLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamLinearSearchByID(ctx, opts...)
}

func (c *singleClient) StreamLinearSearchByIDWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.SearchWithMetadata_StreamLinearSearchByIDWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamLinearSearchByIDWithMetadataRPCName), apiName+"/"+vald.StreamLinearSearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamLinearSearchByIDWithMetadata(ctx, opts...)
}

func (c *singleClient) MultiLinearSearch(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiLinearSearchRPCName), apiName+"/"+vald.MultiLinearSearchRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiLinearSearch(ctx, in, opts...)
}

func (c *singleClient) MultiLinearSearchWithMetadata(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiLinearSearchWithMetadataRPCName), apiName+"/"+vald.MultiLinearSearchWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiLinearSearchWithMetadata(ctx, in, opts...)
}

func (c *singleClient) MultiLinearSearchByID(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiLinearSearchByIDRPCName), apiName+"/"+vald.MultiLinearSearchByIDRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiLinearSearchByID(ctx, in, opts...)
}

func (c *singleClient) MultiLinearSearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (res *payload.Search_Responses, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiLinearSearchByIDWithMetadataRPCName), apiName+"/"+vald.MultiLinearSearchByIDWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiLinearSearchByIDWithMetadata(ctx, in, opts...)
}

func (c *singleClient) Insert(
	ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.InsertRPCName), apiName+"/"+vald.InsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Insert(ctx, in, opts...)
}

func (c *singleClient) InsertWithMetadata(
	ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.InsertWithMetadataRPCName), apiName+"/"+vald.InsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.InsertWithMetadata(ctx, in, opts...)
}

func (c *singleClient) StreamInsert(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Insert_StreamInsertClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamInsertRPCName), apiName+"/"+vald.StreamInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamInsert(ctx, opts...)
}

func (c *singleClient) StreamInsertWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.InsertWithMetadata_StreamInsertWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamInsertWithMetadataRPCName), apiName+"/"+vald.StreamInsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamInsertWithMetadata(ctx, opts...)
}

func (c *singleClient) MultiInsert(
	ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiInsertRPCName), apiName+"/"+vald.MultiInsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiInsert(ctx, in, opts...)
}

func (c *singleClient) MultiInsertWithMetadata(
	ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiInsertWithMetadataRPCName), apiName+"/"+vald.MultiInsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiInsertWithMetadata(ctx, in, opts...)
}

func (c *singleClient) Update(
	ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.UpdateRPCName), apiName+"/"+vald.UpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Update(ctx, in, opts...)
}

func (c *singleClient) UpdateWithMetadata(
	ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.UpdateWithMetadataRPCName), apiName+"/"+vald.UpdateWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.UpdateWithMetadata(ctx, in, opts...)
}

func (c *singleClient) UpdateTimestamp(
	ctx context.Context, in *payload.Update_TimestampRequest, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.UpdateTimestampRPCName), apiName+"/"+vald.UpdateTimestampRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.UpdateTimestamp(ctx, in, opts...)
}

func (c *singleClient) UpdateTimestampWithMetadata(
	ctx context.Context, in *payload.Update_TimestampRequest, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.UpdateTimestampWithMetadataRPCName), apiName+"/"+vald.UpdateTimestampWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.UpdateTimestampWithMetadata(ctx, in, opts...)
}

func (c *singleClient) StreamUpdate(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Update_StreamUpdateClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamUpdateRPCName), apiName+"/"+vald.StreamUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamUpdate(ctx, opts...)
}

func (c *singleClient) StreamUpdateWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.UpdateWithMetadata_StreamUpdateWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamUpdateWithMetadataRPCName), apiName+"/"+vald.StreamUpdateWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamUpdateWithMetadata(ctx, opts...)
}

func (c *singleClient) MultiUpdate(
	ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiUpdateRPCName), apiName+"/"+vald.MultiUpdateRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiUpdate(ctx, in, opts...)
}

func (c *singleClient) MultiUpdateWithMetadata(
	ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiUpdateWithMetadataRPCName), apiName+"/"+vald.MultiUpdateWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiUpdateWithMetadata(ctx, in, opts...)
}

func (c *singleClient) Upsert(
	ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.UpsertRPCName), apiName+"/"+vald.UpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Upsert(ctx, in, opts...)
}

func (c *singleClient) UpsertWithMetadata(
	ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.UpsertWithMetadataRPCName), apiName+"/"+vald.UpsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.UpsertWithMetadata(ctx, in, opts...)
}

func (c *singleClient) StreamUpsert(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Upsert_StreamUpsertClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamUpsertRPCName), apiName+"/"+vald.StreamUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamUpsert(ctx, opts...)
}

func (c *singleClient) StreamUpsertWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.UpsertWithMetadata_StreamUpsertWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamUpsertWithMetadataRPCName), apiName+"/"+vald.StreamUpsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamUpsertWithMetadata(ctx, opts...)
}

func (c *singleClient) MultiUpsert(
	ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiUpsertRPCName), apiName+"/"+vald.MultiUpsertRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiUpsert(ctx, in, opts...)
}

func (c *singleClient) MultiUpsertWithMetadata(
	ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiUpsertWithMetadataRPCName), apiName+"/"+vald.MultiUpsertWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiUpsertWithMetadata(ctx, in, opts...)
}

func (c *singleClient) Remove(
	ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.RemoveRPCName), apiName+"/"+vald.RemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Remove(ctx, in, opts...)
}

func (c *singleClient) RemoveWithMetadata(
	ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption,
) (res *payload.Object_Location, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.RemoveWithMetadataRPCName), apiName+"/"+vald.RemoveWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.RemoveWithMetadata(ctx, in, opts...)
}

func (c *singleClient) StreamRemove(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Remove_StreamRemoveClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamRemoveRPCName), apiName+"/"+vald.StreamRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamRemove(ctx, opts...)
}

func (c *singleClient) StreamRemoveWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.RemoveWithMetadata_StreamRemoveWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamRemoveWithMetadataRPCName), apiName+"/"+vald.StreamRemoveWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamRemoveWithMetadata(ctx, opts...)
}

func (c *singleClient) MultiRemove(
	ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiRemoveRPCName), apiName+"/"+vald.MultiRemoveRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiRemove(ctx, in, opts...)
}

func (c *singleClient) MultiRemoveWithMetadata(
	ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.MultiRemoveWithMetadataRPCName), apiName+"/"+vald.MultiRemoveWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.MultiRemoveWithMetadata(ctx, in, opts...)
}

func (c *singleClient) Flush(
	ctx context.Context, in *payload.Flush_Request, opts ...grpc.CallOption,
) (res *payload.Info_Index_Count, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/singleClient.Flush")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.Flush(ctx, in, opts...)
}

func (c *singleClient) RemoveByTimestamp(
	ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.RemoveByTimestampRPCName), apiName+"/"+vald.RemoveByTimestampRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.RemoveByTimestamp(ctx, in, opts...)
}

func (c *singleClient) RemoveByTimestampWithMetadata(
	ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption,
) (res *payload.Object_Locations, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.RemoveByTimestampWithMetadataRPCName), apiName+"/"+vald.RemoveByTimestampWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.RemoveByTimestampWithMetadata(ctx, in, opts...)
}

func (c *singleClient) GetObject(
	ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption,
) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.GetObjectRPCName), apiName+"/"+vald.GetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.GetObject(ctx, in, opts...)
}

func (c *singleClient) GetObjectWithMetadata(
	ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption,
) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.GetObjectWithMetadataRPCName), apiName+"/"+vald.GetObjectWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.GetObjectWithMetadata(ctx, in, opts...)
}

func (c *singleClient) StreamGetObject(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.Object_StreamGetObjectClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamGetObjectRPCName), apiName+"/"+vald.StreamGetObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamGetObject(ctx, opts...)
}

func (c *singleClient) StreamGetObjectWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (res vald.ObjectWithMetadata_StreamGetObjectWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamGetObjectWithMetadataRPCName), apiName+"/"+vald.StreamGetObjectWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamGetObjectWithMetadata(ctx, opts...)
}

func (c *singleClient) StreamListObject(
	ctx context.Context, in *payload.Object_List_Request, opts ...grpc.CallOption,
) (res vald.Object_StreamListObjectClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamListObjectRPCName), apiName+"/"+vald.StreamListObjectRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamListObject(ctx, in, opts...)
}

func (c *singleClient) StreamListObjectWithMetadata(
	ctx context.Context, in *payload.Object_List_Request, opts ...grpc.CallOption,
) (res vald.ObjectWithMetadata_StreamListObjectWithMetadataClient, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.StreamListObjectWithMetadataRPCName), apiName+"/"+vald.StreamListObjectWithMetadataRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.StreamListObjectWithMetadata(ctx, in, opts...)
}

func (c *singleClient) IndexInfo(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (res *payload.Info_Index_Count, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.IndexInfoRPCName), apiName+"/"+vald.IndexInfoRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.IndexInfo(ctx, in, opts...)
}

func (c *singleClient) IndexDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (res *payload.Info_Index_Detail, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.IndexDetailRPCName), apiName+"/"+vald.IndexDetailRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.IndexDetail(ctx, in, opts...)
}

func (c *singleClient) IndexStatistics(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (res *payload.Info_Index_Statistics, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.IndexStatisticsRPCName), apiName+"/"+vald.IndexStatisticsRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.IndexStatistics(ctx, in, opts...)
}

func (c *singleClient) IndexStatisticsDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (res *payload.Info_Index_StatisticsDetail, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.IndexStatisticsDetailRPCName), apiName+"/"+vald.IndexStatisticsDetailRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.IndexStatisticsDetail(ctx, in, opts...)
}

func (c *singleClient) IndexProperty(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (res *payload.Info_Index_PropertyDetail, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/singleClient/"+vald.IndexPropertyRPCName), apiName+"/"+vald.IndexPropertyRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.IndexProperty(ctx, in, opts...)
}

func (c *singleClient) GetTimestamp(
	ctx context.Context, in *payload.Object_TimestampRequest, opts ...grpc.CallOption,
) (res *payload.Object_Timestamp, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/"+vald.GetTimestampRPCName), apiName+"/"+vald.GetTimestampRPCName)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.vc.GetTimestamp(ctx, in, opts...)
}
