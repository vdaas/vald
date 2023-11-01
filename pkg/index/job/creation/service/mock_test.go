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
