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

package discoverer

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

// DiscovererClient is the client API for Discoverer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DiscovererClient interface {
	// Represent the RPC to get the agent pods information.
	Pods(ctx context.Context, in *payload.Discoverer_Request, opts ...grpc.CallOption) (*payload.Info_Pods, error)
	// Represent the RPC to get the node information.
	Nodes(ctx context.Context, in *payload.Discoverer_Request, opts ...grpc.CallOption) (*payload.Info_Nodes, error)
	// Represent the RPC to get the readreplica svc information.
	Services(ctx context.Context, in *payload.Discoverer_Request, opts ...grpc.CallOption) (*payload.Info_Services, error)
}

type discovererClient struct {
	cc grpc.ClientConnInterface
}

func NewDiscovererClient(cc grpc.ClientConnInterface) DiscovererClient {
	return &discovererClient{cc}
}

func (c *discovererClient) Pods(ctx context.Context, in *payload.Discoverer_Request, opts ...grpc.CallOption) (*payload.Info_Pods, error) {
	out := new(payload.Info_Pods)
	err := c.cc.Invoke(ctx, "/discoverer.v1.Discoverer/Pods", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discovererClient) Nodes(ctx context.Context, in *payload.Discoverer_Request, opts ...grpc.CallOption) (*payload.Info_Nodes, error) {
	out := new(payload.Info_Nodes)
	err := c.cc.Invoke(ctx, "/discoverer.v1.Discoverer/Nodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discovererClient) Services(ctx context.Context, in *payload.Discoverer_Request, opts ...grpc.CallOption) (*payload.Info_Services, error) {
	out := new(payload.Info_Services)
	err := c.cc.Invoke(ctx, "/discoverer.v1.Discoverer/Services", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiscovererServer is the server API for Discoverer service.
// All implementations must embed UnimplementedDiscovererServer
// for forward compatibility
type DiscovererServer interface {
	// Represent the RPC to get the agent pods information.
	Pods(context.Context, *payload.Discoverer_Request) (*payload.Info_Pods, error)
	// Represent the RPC to get the node information.
	Nodes(context.Context, *payload.Discoverer_Request) (*payload.Info_Nodes, error)
	// Represent the RPC to get the readreplica svc information.
	Services(context.Context, *payload.Discoverer_Request) (*payload.Info_Services, error)
	mustEmbedUnimplementedDiscovererServer()
}

// UnimplementedDiscovererServer must be embedded to have forward compatible implementations.
type UnimplementedDiscovererServer struct {
}

func (UnimplementedDiscovererServer) Pods(context.Context, *payload.Discoverer_Request) (*payload.Info_Pods, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pods not implemented")
}
func (UnimplementedDiscovererServer) Nodes(context.Context, *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Nodes not implemented")
}
func (UnimplementedDiscovererServer) Services(context.Context, *payload.Discoverer_Request) (*payload.Info_Services, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Services not implemented")
}
func (UnimplementedDiscovererServer) mustEmbedUnimplementedDiscovererServer() {}

// UnsafeDiscovererServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DiscovererServer will
// result in compilation errors.
type UnsafeDiscovererServer interface {
	mustEmbedUnimplementedDiscovererServer()
}

func RegisterDiscovererServer(s grpc.ServiceRegistrar, srv DiscovererServer) {
	s.RegisterService(&Discoverer_ServiceDesc, srv)
}

func _Discoverer_Pods_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Discoverer_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscovererServer).Pods(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discoverer.v1.Discoverer/Pods",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscovererServer).Pods(ctx, req.(*payload.Discoverer_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Discoverer_Nodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Discoverer_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscovererServer).Nodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discoverer.v1.Discoverer/Nodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscovererServer).Nodes(ctx, req.(*payload.Discoverer_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Discoverer_Services_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Discoverer_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscovererServer).Services(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discoverer.v1.Discoverer/Services",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscovererServer).Services(ctx, req.(*payload.Discoverer_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Discoverer_ServiceDesc is the grpc.ServiceDesc for Discoverer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Discoverer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "discoverer.v1.Discoverer",
	HandlerType: (*DiscovererServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Pods",
			Handler:    _Discoverer_Pods_Handler,
		},
		{
			MethodName: "Nodes",
			Handler:    _Discoverer_Nodes_Handler,
		},
		{
			MethodName: "Services",
			Handler:    _Discoverer_Services_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/discoverer/discoverer.proto",
}
