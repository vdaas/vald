// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package tikv

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/tikv"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

const (
	pdApiName = "vald/internal/client/v1/client/meta/pd"
)

type pdClient struct {
	addrs []string
	c     grpc.Client
}

func (c *pdClient) Start(ctx context.Context) (<-chan error, error) {
	return c.c.StartConnectionMonitor(ctx)
}

func (c *pdClient) Stop(ctx context.Context) error {
	return c.c.Close(ctx)
}

func (c *pdClient) GRPCClient() grpc.Client {
	return c.c
}

func (c *pdClient) GetAllStores(
	ctx context.Context, in *tikv.GetAllStoresRequest,
) (res *tikv.GetAllStoresResponse, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/GetAllStores"), pdApiName+"/GetAllStores")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.GetAllStoresResponse, error) {
		return tikv.NewPDClient(conn).GetAllStores(ctx, in, copts...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *pdClient) BatchScanRegions(
	ctx context.Context, in *tikv.BatchScanRegionsRequest,
) (res *tikv.BatchScanRegionsResponse, err error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/BatchScanRegions"), pdApiName+"/BatchScanRegions")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (*tikv.BatchScanRegionsResponse, error) {
		return tikv.NewPDClient(conn).BatchScanRegions(ctx, in, copts...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
