//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package tikv

import (
	context "context"

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

// TikvClient is the client API for Tikv service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TikvClient interface {
	RawGet(ctx context.Context, in *RawGetRequest, opts ...grpc.CallOption) (*RawGetResponse, error)
	RawBatchGet(ctx context.Context, in *RawBatchGetRequest, opts ...grpc.CallOption) (*RawBatchGetResponse, error)
	RawPut(ctx context.Context, in *RawPutRequest, opts ...grpc.CallOption) (*RawPutResponse, error)
	RawBatchPut(ctx context.Context, in *RawBatchPutRequest, opts ...grpc.CallOption) (*RawBatchPutResponse, error)
	RawDelete(ctx context.Context, in *RawDeleteRequest, opts ...grpc.CallOption) (*RawDeleteResponse, error)
	RawBatchDelete(ctx context.Context, in *RawBatchDeleteRequest, opts ...grpc.CallOption) (*RawBatchDeleteResponse, error)
}

type tikvClient struct {
	cc grpc.ClientConnInterface
}

func NewTikvClient(cc grpc.ClientConnInterface) TikvClient {
	return &tikvClient{cc}
}

func (c *tikvClient) RawGet(
	ctx context.Context, in *RawGetRequest, opts ...grpc.CallOption,
) (*RawGetResponse, error) {
	out := new(RawGetResponse)
	err := c.cc.Invoke(ctx, "/tikv.Tikv/RawGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tikvClient) RawBatchGet(
	ctx context.Context, in *RawBatchGetRequest, opts ...grpc.CallOption,
) (*RawBatchGetResponse, error) {
	out := new(RawBatchGetResponse)
	err := c.cc.Invoke(ctx, "/tikv.Tikv/RawBatchGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tikvClient) RawPut(
	ctx context.Context, in *RawPutRequest, opts ...grpc.CallOption,
) (*RawPutResponse, error) {
	out := new(RawPutResponse)
	err := c.cc.Invoke(ctx, "/tikv.Tikv/RawPut", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tikvClient) RawBatchPut(
	ctx context.Context, in *RawBatchPutRequest, opts ...grpc.CallOption,
) (*RawBatchPutResponse, error) {
	out := new(RawBatchPutResponse)
	err := c.cc.Invoke(ctx, "/tikv.Tikv/RawBatchPut", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tikvClient) RawDelete(
	ctx context.Context, in *RawDeleteRequest, opts ...grpc.CallOption,
) (*RawDeleteResponse, error) {
	out := new(RawDeleteResponse)
	err := c.cc.Invoke(ctx, "/tikv.Tikv/RawDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tikvClient) RawBatchDelete(
	ctx context.Context, in *RawBatchDeleteRequest, opts ...grpc.CallOption,
) (*RawBatchDeleteResponse, error) {
	out := new(RawBatchDeleteResponse)
	err := c.cc.Invoke(ctx, "/tikv.Tikv/RawBatchDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TikvServer is the server API for Tikv service.
// All implementations must embed UnimplementedTikvServer
// for forward compatibility
type TikvServer interface {
	RawGet(context.Context, *RawGetRequest) (*RawGetResponse, error)
	RawBatchGet(context.Context, *RawBatchGetRequest) (*RawBatchGetResponse, error)
	RawPut(context.Context, *RawPutRequest) (*RawPutResponse, error)
	RawBatchPut(context.Context, *RawBatchPutRequest) (*RawBatchPutResponse, error)
	RawDelete(context.Context, *RawDeleteRequest) (*RawDeleteResponse, error)
	RawBatchDelete(context.Context, *RawBatchDeleteRequest) (*RawBatchDeleteResponse, error)
	mustEmbedUnimplementedTikvServer()
}

// UnimplementedTikvServer must be embedded to have forward compatible implementations.
type UnimplementedTikvServer struct{}

func (UnimplementedTikvServer) RawGet(context.Context, *RawGetRequest) (*RawGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawGet not implemented")
}

func (UnimplementedTikvServer) RawBatchGet(
	context.Context, *RawBatchGetRequest,
) (*RawBatchGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawBatchGet not implemented")
}

func (UnimplementedTikvServer) RawPut(context.Context, *RawPutRequest) (*RawPutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawPut not implemented")
}

func (UnimplementedTikvServer) RawBatchPut(
	context.Context, *RawBatchPutRequest,
) (*RawBatchPutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawBatchPut not implemented")
}

func (UnimplementedTikvServer) RawDelete(
	context.Context, *RawDeleteRequest,
) (*RawDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawDelete not implemented")
}

func (UnimplementedTikvServer) RawBatchDelete(
	context.Context, *RawBatchDeleteRequest,
) (*RawBatchDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RawBatchDelete not implemented")
}
func (UnimplementedTikvServer) mustEmbedUnimplementedTikvServer() {}

// UnsafeTikvServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TikvServer will
// result in compilation errors.
type UnsafeTikvServer interface {
	mustEmbedUnimplementedTikvServer()
}

func RegisterTikvServer(s grpc.ServiceRegistrar, srv TikvServer) {
	s.RegisterService(&Tikv_ServiceDesc, srv)
}

func _Tikv_RawGet_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(RawGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TikvServer).RawGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tikv.Tikv/RawGet",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(TikvServer).RawGet(ctx, req.(*RawGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tikv_RawBatchGet_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(RawBatchGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TikvServer).RawBatchGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tikv.Tikv/RawBatchGet",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(TikvServer).RawBatchGet(ctx, req.(*RawBatchGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tikv_RawPut_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(RawPutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TikvServer).RawPut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tikv.Tikv/RawPut",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(TikvServer).RawPut(ctx, req.(*RawPutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tikv_RawBatchPut_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(RawBatchPutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TikvServer).RawBatchPut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tikv.Tikv/RawBatchPut",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(TikvServer).RawBatchPut(ctx, req.(*RawBatchPutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tikv_RawDelete_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(RawDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TikvServer).RawDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tikv.Tikv/RawDelete",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(TikvServer).RawDelete(ctx, req.(*RawDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tikv_RawBatchDelete_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(RawBatchDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TikvServer).RawBatchDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tikv.Tikv/RawBatchDelete",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(TikvServer).RawBatchDelete(ctx, req.(*RawBatchDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Tikv_ServiceDesc is the grpc.ServiceDesc for Tikv service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tikv_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tikv.Tikv",
	HandlerType: (*TikvServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RawGet",
			Handler:    _Tikv_RawGet_Handler,
		},
		{
			MethodName: "RawBatchGet",
			Handler:    _Tikv_RawBatchGet_Handler,
		},
		{
			MethodName: "RawPut",
			Handler:    _Tikv_RawPut_Handler,
		},
		{
			MethodName: "RawBatchPut",
			Handler:    _Tikv_RawBatchPut_Handler,
		},
		{
			MethodName: "RawDelete",
			Handler:    _Tikv_RawDelete_Handler,
		},
		{
			MethodName: "RawBatchDelete",
			Handler:    _Tikv_RawBatchDelete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/tikv/tikvpb.proto",
}
