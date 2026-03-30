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

// IndexClient is the client API for Index service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IndexClient interface {
	// Overview
	// Represent the RPC to get the index information.
	IndexInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Index_Count, error)
	// Overview
	// Represent the RPC to get the index information for each agents.
	IndexDetail(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Index_Detail, error)
	// Overview
	// Represent the RPC to get the index statistics.
	IndexStatistics(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Index_Statistics, error)
	// Overview
	// Represent the RPC to get the index statistics for each agents.
	IndexStatisticsDetail(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Index_StatisticsDetail, error)
	// Overview
	// Represent the RPC to get the index property.
	IndexProperty(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Index_PropertyDetail, error)
}

type indexClient struct {
	cc grpc.ClientConnInterface
}

func NewIndexClient(cc grpc.ClientConnInterface) IndexClient {
	return &indexClient{cc}
}

func (c *indexClient) IndexInfo(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_Count, error) {
	out := new(payload.Info_Index_Count)
	err := c.cc.Invoke(ctx, "/vald.v1.Index/IndexInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) IndexDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_Detail, error) {
	out := new(payload.Info_Index_Detail)
	err := c.cc.Invoke(ctx, "/vald.v1.Index/IndexDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) IndexStatistics(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_Statistics, error) {
	out := new(payload.Info_Index_Statistics)
	err := c.cc.Invoke(ctx, "/vald.v1.Index/IndexStatistics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) IndexStatisticsDetail(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_StatisticsDetail, error) {
	out := new(payload.Info_Index_StatisticsDetail)
	err := c.cc.Invoke(ctx, "/vald.v1.Index/IndexStatisticsDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indexClient) IndexProperty(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Info_Index_PropertyDetail, error) {
	out := new(payload.Info_Index_PropertyDetail)
	err := c.cc.Invoke(ctx, "/vald.v1.Index/IndexProperty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IndexServer is the server API for Index service.
// All implementations must embed UnimplementedIndexServer
// for forward compatibility
type IndexServer interface {
	// Overview
	// Represent the RPC to get the index information.
	IndexInfo(context.Context, *payload.Empty) (*payload.Info_Index_Count, error)
	// Overview
	// Represent the RPC to get the index information for each agents.
	IndexDetail(context.Context, *payload.Empty) (*payload.Info_Index_Detail, error)
	// Overview
	// Represent the RPC to get the index statistics.
	IndexStatistics(context.Context, *payload.Empty) (*payload.Info_Index_Statistics, error)
	// Overview
	// Represent the RPC to get the index statistics for each agents.
	IndexStatisticsDetail(context.Context, *payload.Empty) (*payload.Info_Index_StatisticsDetail, error)
	// Overview
	// Represent the RPC to get the index property.
	IndexProperty(context.Context, *payload.Empty) (*payload.Info_Index_PropertyDetail, error)
	mustEmbedUnimplementedIndexServer()
}

// UnimplementedIndexServer must be embedded to have forward compatible implementations.
type UnimplementedIndexServer struct{}

func (UnimplementedIndexServer) IndexInfo(
	context.Context, *payload.Empty,
) (*payload.Info_Index_Count, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IndexInfo not implemented")
}

func (UnimplementedIndexServer) IndexDetail(
	context.Context, *payload.Empty,
) (*payload.Info_Index_Detail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IndexDetail not implemented")
}

func (UnimplementedIndexServer) IndexStatistics(
	context.Context, *payload.Empty,
) (*payload.Info_Index_Statistics, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IndexStatistics not implemented")
}

func (UnimplementedIndexServer) IndexStatisticsDetail(
	context.Context, *payload.Empty,
) (*payload.Info_Index_StatisticsDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IndexStatisticsDetail not implemented")
}

func (UnimplementedIndexServer) IndexProperty(
	context.Context, *payload.Empty,
) (*payload.Info_Index_PropertyDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IndexProperty not implemented")
}
func (UnimplementedIndexServer) mustEmbedUnimplementedIndexServer() {}

// UnsafeIndexServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IndexServer will
// result in compilation errors.
type UnsafeIndexServer interface {
	mustEmbedUnimplementedIndexServer()
}

func RegisterIndexServer(s grpc.ServiceRegistrar, srv IndexServer) {
	s.RegisterService(&Index_ServiceDesc, srv)
}

func _Index_IndexInfo_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).IndexInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Index/IndexInfo",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(IndexServer).IndexInfo(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_IndexDetail_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).IndexDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Index/IndexDetail",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(IndexServer).IndexDetail(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_IndexStatistics_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).IndexStatistics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Index/IndexStatistics",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(IndexServer).IndexStatistics(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_IndexStatisticsDetail_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).IndexStatisticsDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Index/IndexStatisticsDetail",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(IndexServer).IndexStatisticsDetail(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Index_IndexProperty_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).IndexProperty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Index/IndexProperty",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(IndexServer).IndexProperty(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Index_ServiceDesc is the grpc.ServiceDesc for Index service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Index_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vald.v1.Index",
	HandlerType: (*IndexServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IndexInfo",
			Handler:    _Index_IndexInfo_Handler,
		},
		{
			MethodName: "IndexDetail",
			Handler:    _Index_IndexDetail_Handler,
		},
		{
			MethodName: "IndexStatistics",
			Handler:    _Index_IndexStatistics_Handler,
		},
		{
			MethodName: "IndexStatisticsDetail",
			Handler:    _Index_IndexStatisticsDetail_Handler,
		},
		{
			MethodName: "IndexProperty",
			Handler:    _Index_IndexProperty_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/vald/index.proto",
}
