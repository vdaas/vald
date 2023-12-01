//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

package ingress

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

// FilterClient is the client API for Filter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FilterClient interface {
	// Represent the RPC to generate the vector.
	GenVector(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Vector, error)
	// Represent the RPC to filter the vector.
	FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Object_Vector, error)
}

type filterClient struct {
	cc grpc.ClientConnInterface
}

func NewFilterClient(cc grpc.ClientConnInterface) FilterClient {
	return &filterClient{cc}
}

func (c *filterClient) GenVector(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/filter.ingress.v1.Filter/GenVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/filter.ingress.v1.Filter/FilterVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilterServer is the server API for Filter service.
// All implementations must embed UnimplementedFilterServer
// for forward compatibility
type FilterServer interface {
	// Represent the RPC to generate the vector.
	GenVector(context.Context, *payload.Object_Blob) (*payload.Object_Vector, error)
	// Represent the RPC to filter the vector.
	FilterVector(context.Context, *payload.Object_Vector) (*payload.Object_Vector, error)
	mustEmbedUnimplementedFilterServer()
}

// UnimplementedFilterServer must be embedded to have forward compatible implementations.
type UnimplementedFilterServer struct {
}

func (UnimplementedFilterServer) GenVector(context.Context, *payload.Object_Blob) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenVector not implemented")
}
func (UnimplementedFilterServer) FilterVector(context.Context, *payload.Object_Vector) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FilterVector not implemented")
}
func (UnimplementedFilterServer) mustEmbedUnimplementedFilterServer() {}

// UnsafeFilterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FilterServer will
// result in compilation errors.
type UnsafeFilterServer interface {
	mustEmbedUnimplementedFilterServer()
}

func RegisterFilterServer(s grpc.ServiceRegistrar, srv FilterServer) {
	s.RegisterService(&Filter_ServiceDesc, srv)
}

func _Filter_GenVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).GenVector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.ingress.v1.Filter/GenVector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).GenVector(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_FilterVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).FilterVector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.ingress.v1.Filter/FilterVector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).FilterVector(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

// Filter_ServiceDesc is the grpc.ServiceDesc for Filter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Filter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "filter.ingress.v1.Filter",
	HandlerType: (*FilterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenVector",
			Handler:    _Filter_GenVector_Handler,
		},
		{
			MethodName: "FilterVector",
			Handler:    _Filter_FilterVector_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/proto/v1/filter/ingress/ingress_filter.proto",
}
