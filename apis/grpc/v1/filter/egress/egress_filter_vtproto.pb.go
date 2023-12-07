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

package egress

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
	// Represent the RPC to filter the distance.
	FilterDistance(ctx context.Context, in *payload.Filter_DistanceRequest, opts ...grpc.CallOption) (*payload.Filter_DistanceResponse, error)
	// Represent the RPC to filter the vector.
	FilterVector(ctx context.Context, in *payload.Filter_VectorRequest, opts ...grpc.CallOption) (*payload.Filter_VectorResponse, error)
}

type filterClient struct {
	cc grpc.ClientConnInterface
}

func NewFilterClient(cc grpc.ClientConnInterface) FilterClient {
	return &filterClient{cc}
}

func (c *filterClient) FilterDistance(ctx context.Context, in *payload.Filter_DistanceRequest, opts ...grpc.CallOption) (*payload.Filter_DistanceResponse, error) {
	out := new(payload.Filter_DistanceResponse)
	err := c.cc.Invoke(ctx, "/filter.egress.v1.Filter/FilterDistance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) FilterVector(ctx context.Context, in *payload.Filter_VectorRequest, opts ...grpc.CallOption) (*payload.Filter_VectorResponse, error) {
	out := new(payload.Filter_VectorResponse)
	err := c.cc.Invoke(ctx, "/filter.egress.v1.Filter/FilterVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilterServer is the server API for Filter service.
// All implementations must embed UnimplementedFilterServer
// for forward compatibility
type FilterServer interface {
	// Represent the RPC to filter the distance.
	FilterDistance(context.Context, *payload.Filter_DistanceRequest) (*payload.Filter_DistanceResponse, error)
	// Represent the RPC to filter the vector.
	FilterVector(context.Context, *payload.Filter_VectorRequest) (*payload.Filter_VectorResponse, error)
	mustEmbedUnimplementedFilterServer()
}

// UnimplementedFilterServer must be embedded to have forward compatible implementations.
type UnimplementedFilterServer struct {
}

func (UnimplementedFilterServer) FilterDistance(context.Context, *payload.Filter_DistanceRequest) (*payload.Filter_DistanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FilterDistance not implemented")
}
func (UnimplementedFilterServer) FilterVector(context.Context, *payload.Filter_VectorRequest) (*payload.Filter_VectorResponse, error) {
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

func _Filter_FilterDistance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Filter_DistanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).FilterDistance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.egress.v1.Filter/FilterDistance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).FilterDistance(ctx, req.(*payload.Filter_DistanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_FilterVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Filter_VectorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).FilterVector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.egress.v1.Filter/FilterVector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).FilterVector(ctx, req.(*payload.Filter_VectorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Filter_ServiceDesc is the grpc.ServiceDesc for Filter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Filter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "filter.egress.v1.Filter",
	HandlerType: (*FilterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FilterDistance",
			Handler:    _Filter_FilterDistance_Handler,
		},
		{
			MethodName: "FilterVector",
			Handler:    _Filter_FilterVector_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/filter/egress/egress_filter.proto",
}
