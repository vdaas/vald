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

package stats

import (
	context "context"

	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	codes "github.com/vdaas/vald/internal/net/grpc/codes"
	status "github.com/vdaas/vald/internal/net/grpc/status"
	grpc "google.golang.org/grpc"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StatsClient is the client API for Stats service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatsClient interface {
	// Represent the RPC to get the resource stats.
	ResourceStats(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Stats_ResourceStats, error)
}

type statsClient struct {
	cc grpc.ClientConnInterface
}

func NewStatsClient(cc grpc.ClientConnInterface) StatsClient {
	return &statsClient{cc}
}

func (c *statsClient) ResourceStats(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Stats_ResourceStats, error) {
	out := new(payload.Info_Stats_ResourceStats)
	err := c.cc.Invoke(ctx, "/rpc.v1.Stats/ResourceStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatsServer is the server API for Stats service.
// All implementations must embed UnimplementedStatsServer
// for forward compatibility
type StatsServer interface {
	// Represent the RPC to get the resource stats.
	ResourceStats(context.Context, *payload.Empty) (*payload.Info_Stats_ResourceStats, error)
	mustEmbedUnimplementedStatsServer()
}

// UnimplementedStatsServer must be embedded to have forward compatible implementations.
type UnimplementedStatsServer struct{}

func (UnimplementedStatsServer) ResourceStats(
	context.Context, *payload.Empty,
) (*payload.Info_Stats_ResourceStats, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResourceStats not implemented")
}
func (UnimplementedStatsServer) mustEmbedUnimplementedStatsServer() {}

// UnsafeStatsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatsServer will
// result in compilation errors.
type UnsafeStatsServer interface {
	mustEmbedUnimplementedStatsServer()
}

func RegisterStatsServer(s grpc.ServiceRegistrar, srv StatsServer) {
	s.RegisterService(&Stats_ServiceDesc, srv)
}

func _Stats_ResourceStats_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsServer).ResourceStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.v1.Stats/ResourceStats",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(StatsServer).ResourceStats(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Stats_ServiceDesc is the grpc.ServiceDesc for Stats service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stats_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.v1.Stats",
	HandlerType: (*StatsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResourceStats",
			Handler:    _Stats_ResourceStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/rpc/stats/stats.proto",
}

// StatsDetailClient is the client API for StatsDetail service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatsDetailClient interface {
	// Represent the RPC to get the resource stats for each agents.
	ResourceStatsDetail(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Stats_ResourceStatsDetail, error)
}

type statsDetailClient struct {
	cc grpc.ClientConnInterface
}

func NewStatsDetailClient(cc grpc.ClientConnInterface) StatsDetailClient {
	return &statsDetailClient{cc}
}

func (c *statsDetailClient) ResourceStatsDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Stats_ResourceStatsDetail, error) {
	out := new(payload.Info_Stats_ResourceStatsDetail)
	err := c.cc.Invoke(ctx, "/rpc.v1.StatsDetail/ResourceStatsDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatsDetailServer is the server API for StatsDetail service.
// All implementations must embed UnimplementedStatsDetailServer
// for forward compatibility
type StatsDetailServer interface {
	// Represent the RPC to get the resource stats for each agents.
	ResourceStatsDetail(context.Context, *payload.Empty) (*payload.Info_Stats_ResourceStatsDetail, error)
	mustEmbedUnimplementedStatsDetailServer()
}

// UnimplementedStatsDetailServer must be embedded to have forward compatible implementations.
type UnimplementedStatsDetailServer struct{}

func (UnimplementedStatsDetailServer) ResourceStatsDetail(
	context.Context, *payload.Empty,
) (*payload.Info_Stats_ResourceStatsDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResourceStatsDetail not implemented")
}
func (UnimplementedStatsDetailServer) mustEmbedUnimplementedStatsDetailServer() {}

// UnsafeStatsDetailServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatsDetailServer will
// result in compilation errors.
type UnsafeStatsDetailServer interface {
	mustEmbedUnimplementedStatsDetailServer()
}

func RegisterStatsDetailServer(s grpc.ServiceRegistrar, srv StatsDetailServer) {
	s.RegisterService(&StatsDetail_ServiceDesc, srv)
}

func _StatsDetail_ResourceStatsDetail_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatsDetailServer).ResourceStatsDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.v1.StatsDetail/ResourceStatsDetail",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(StatsDetailServer).ResourceStatsDetail(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// StatsDetail_ServiceDesc is the grpc.ServiceDesc for StatsDetail service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StatsDetail_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.v1.StatsDetail",
	HandlerType: (*StatsDetailServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResourceStatsDetail",
			Handler:    _StatsDetail_ResourceStatsDetail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/rpc/stats/stats.proto",
}
