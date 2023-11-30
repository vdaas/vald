// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
package service

import (
	"context"

	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/net/grpc"
)

type mockDiscovererClient struct {
	discoverer.Client
	GetAddrsFunc  func(ctx context.Context) []string
	GetClientFunc func() grpc.Client
}

func (mc *mockDiscovererClient) GetAddrs(ctx context.Context) []string {
	return mc.GetAddrsFunc(ctx)
}

func (mc *mockDiscovererClient) GetClient() grpc.Client {
	return mc.GetClientFunc()
}

type mockGrpcClient struct {
	grpc.Client
	OrderedRangeConcurrentFunc func(ctx context.Context,
		order []string,
		concurrency int,
		f func(ctx context.Context,
			addr string,
			conn *grpc.ClientConn,
			copts ...grpc.CallOption) error) error
}

func (mc *mockGrpcClient) OrderedRangeConcurrent(ctx context.Context,
	order []string,
	concurrency int,
	f func(ctx context.Context,
		addr string,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) error,
) error {
	return mc.OrderedRangeConcurrentFunc(ctx, order, concurrency, f)
}
