//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

package meta

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

// MetaClient is the client API for Meta service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetaClient interface {
	Get(ctx context.Context, in *payload.Meta_Key, opts ...grpc.CallOption) (*payload.Meta_Value, error)
	Set(ctx context.Context, in *payload.Meta_KeyValue, opts ...grpc.CallOption) (*payload.Empty, error)
	Delete(ctx context.Context, in *payload.Meta_Key, opts ...grpc.CallOption) (*payload.Empty, error)
}

type metaClient struct {
	cc grpc.ClientConnInterface
}

func NewMetaClient(cc grpc.ClientConnInterface) MetaClient {
	return &metaClient{cc}
}

func (c *metaClient) Get(
	ctx context.Context, in *payload.Meta_Key, opts ...grpc.CallOption,
) (*payload.Meta_Value, error) {
	out := new(payload.Meta_Value)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) Set(
	ctx context.Context, in *payload.Meta_KeyValue, opts ...grpc.CallOption,
) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) Delete(
	ctx context.Context, in *payload.Meta_Key, opts ...grpc.CallOption,
) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetaServer is the server API for Meta service.
// All implementations must embed UnimplementedMetaServer
// for forward compatibility
type MetaServer interface {
	Get(context.Context, *payload.Meta_Key) (*payload.Meta_Value, error)
	Set(context.Context, *payload.Meta_KeyValue) (*payload.Empty, error)
	Delete(context.Context, *payload.Meta_Key) (*payload.Empty, error)
	mustEmbedUnimplementedMetaServer()
}

// UnimplementedMetaServer must be embedded to have forward compatible implementations.
type UnimplementedMetaServer struct{}

func (UnimplementedMetaServer) Get(
	context.Context, *payload.Meta_Key,
) (*payload.Meta_Value, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func (UnimplementedMetaServer) Set(
	context.Context, *payload.Meta_KeyValue,
) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}

func (UnimplementedMetaServer) Delete(context.Context, *payload.Meta_Key) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedMetaServer) mustEmbedUnimplementedMetaServer() {}

// UnsafeMetaServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetaServer will
// result in compilation errors.
type UnsafeMetaServer interface {
	mustEmbedUnimplementedMetaServer()
}

func RegisterMetaServer(s grpc.ServiceRegistrar, srv MetaServer) {
	s.RegisterService(&Meta_ServiceDesc, srv)
}

func _Meta_Get_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Meta_Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/Get",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(MetaServer).Get(ctx, req.(*payload.Meta_Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_Set_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Meta_KeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/Set",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(MetaServer).Set(ctx, req.(*payload.Meta_KeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_Delete_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Meta_Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/Delete",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(MetaServer).Delete(ctx, req.(*payload.Meta_Key))
	}
	return interceptor(ctx, in, info, handler)
}

// Meta_ServiceDesc is the grpc.ServiceDesc for Meta service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Meta_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "meta.v1.Meta",
	HandlerType: (*MetaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Meta_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _Meta_Set_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Meta_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/meta/meta.proto",
}
