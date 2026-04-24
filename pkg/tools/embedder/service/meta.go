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

package service

import (
	"context"

	metapb "github.com/vdaas/vald/apis/grpc/v1/meta"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	vgrpc "github.com/vdaas/vald/internal/net/grpc"
)

type MetaClient interface {
	Start(context.Context) (<-chan error, error)
	Stop(context.Context) error
	Set(context.Context, *payload.Meta_KeyValue, ...vgrpc.CallOption) (*payload.Empty, error)
	Delete(context.Context, *payload.Meta_Key, ...vgrpc.CallOption) (*payload.Empty, error)
}

type metaClient struct {
	client vgrpc.Client
}

func NewMetaClient(c vgrpc.Client) (MetaClient, error) {
	if c == nil {
		return nil, errors.New("grpc client is nil")
	}
	return &metaClient{client: c}, nil
}

func (m *metaClient) Start(ctx context.Context) (<-chan error, error) {
	return m.client.StartConnectionMonitor(ctx)
}

func (m *metaClient) Stop(ctx context.Context) error {
	return m.client.Close(ctx)
}

func (m *metaClient) Set(
	ctx context.Context, in *payload.Meta_KeyValue, opts ...vgrpc.CallOption,
) (*payload.Empty, error) {
	return vgrpc.RoundRobin(ctx, m.client, func(ctx context.Context, conn *vgrpc.ClientConn, copts ...vgrpc.CallOption) (*payload.Empty, error) {
		return metapb.NewMetaClient(conn).Set(ctx, in, append(copts, opts...)...)
	})
}

func (m *metaClient) Delete(
	ctx context.Context, in *payload.Meta_Key, opts ...vgrpc.CallOption,
) (*payload.Empty, error) {
	return vgrpc.RoundRobin(ctx, m.client, func(ctx context.Context, conn *vgrpc.ClientConn, copts ...vgrpc.CallOption) (*payload.Empty, error) {
		return metapb.NewMetaClient(conn).Delete(ctx, in, append(copts, opts...)...)
	})
}
