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

package mirror

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

// MirrorClient is the client API for Mirror service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MirrorClient interface {
	// Register is the RPC to register other mirror servers.
	Register(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (*payload.Mirror_Targets, error)
}

type mirrorClient struct {
	cc grpc.ClientConnInterface
}

func NewMirrorClient(cc grpc.ClientConnInterface) MirrorClient {
	return &mirrorClient{cc}
}

func (c *mirrorClient) Register(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (*payload.Mirror_Targets, error) {
	out := new(payload.Mirror_Targets)
	err := c.cc.Invoke(ctx, "/mirror.v1.Mirror/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MirrorServer is the server API for Mirror service.
// All implementations must embed UnimplementedMirrorServer
// for forward compatibility
type MirrorServer interface {
	// Register is the RPC to register other mirror servers.
	Register(context.Context, *payload.Mirror_Targets) (*payload.Mirror_Targets, error)
	mustEmbedUnimplementedMirrorServer()
}

// UnimplementedMirrorServer must be embedded to have forward compatible implementations.
type UnimplementedMirrorServer struct {
}

func (UnimplementedMirrorServer) Register(context.Context, *payload.Mirror_Targets) (*payload.Mirror_Targets, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedMirrorServer) mustEmbedUnimplementedMirrorServer() {}

// UnsafeMirrorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MirrorServer will
// result in compilation errors.
type UnsafeMirrorServer interface {
	mustEmbedUnimplementedMirrorServer()
}

func RegisterMirrorServer(s grpc.ServiceRegistrar, srv MirrorServer) {
	s.RegisterService(&Mirror_ServiceDesc, srv)
}

func _Mirror_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Mirror_Targets)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MirrorServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mirror.v1.Mirror/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MirrorServer).Register(ctx, req.(*payload.Mirror_Targets))
	}
	return interceptor(ctx, in, info, handler)
}

// Mirror_ServiceDesc is the grpc.ServiceDesc for Mirror service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Mirror_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mirror.v1.Mirror",
	HandlerType: (*MirrorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Mirror_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/mirror/mirror.proto",
}
