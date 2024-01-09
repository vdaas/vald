// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
package grpc

import (
	"context"

	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/pool"
)

// GRPCClientMock is the mock for gRPC client.
type GRPCClientMock struct {
	grpc.Client
	OrderedRangeConcurrentFunc func(ctx context.Context,
		order []string,
		concurrency int,
		f func(ctx context.Context,
			addr string,
			conn *grpc.ClientConn,
			copts ...grpc.CallOption) error) error
	ConnectFunc        func(ctx context.Context, addr string, dopts ...grpc.DialOption) (pool.Conn, error)
	DisconnectFunc     func(ctx context.Context, addr string) error
	IsConnectedFunc    func(ctx context.Context, addr string) bool
	ConnectedAddrsFunc func() []string
}

// OrderedRangeConcurrent calls the OrderedRangeConcurrentFunc object.
func (gc *GRPCClientMock) OrderedRangeConcurrent(ctx context.Context,
	order []string,
	concurrency int,
	f func(ctx context.Context,
		addr string,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption) error,
) error {
	return gc.OrderedRangeConcurrentFunc(ctx, order, concurrency, f)
}

// ConnectedAddrs calls the ConnectedAddrsFunc object.
func (gc *GRPCClientMock) ConnectedAddrs() []string {
	return gc.ConnectedAddrsFunc()
}

// Connect calls the ConnectFunc object.
func (gc *GRPCClientMock) Connect(ctx context.Context, addr string, dopts ...grpc.DialOption) (pool.Conn, error) {
	return gc.ConnectFunc(ctx, addr, dopts...)
}

// Disconnect calls the DisconnectFunc object.
func (gc *GRPCClientMock) Disconnect(ctx context.Context, addr string) error {
	return gc.DisconnectFunc(ctx, addr)
}

// IsConnected calls the IsConnectedFunc object.
func (gc *GRPCClientMock) IsConnected(ctx context.Context, addr string) bool {
	return gc.IsConnectedFunc(ctx, addr)
}
