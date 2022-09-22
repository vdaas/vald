//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

package mirroring

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

// MirroringClient is the client API for Mirroring service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MirroringClient interface {
	// Register the RPC to register other mirroring servers.
	Register(ctx context.Context, in *payload.Mirroring_Request, opts ...grpc.CallOption) (*payload.Empty, error)
}

type mirroringClient struct {
	cc grpc.ClientConnInterface
}

func NewMirroringClient(cc grpc.ClientConnInterface) MirroringClient {
	return &mirroringClient{cc}
}

func (c *mirroringClient) Register(ctx context.Context, in *payload.Mirroring_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/mirroring.v1.Mirroring/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MirroringServer is the server API for Mirroring service.
// All implementations must embed UnimplementedMirroringServer
// for forward compatibility
type MirroringServer interface {
	// Register the RPC to register other mirroring servers.
	Register(context.Context, *payload.Mirroring_Request) (*payload.Empty, error)
	mustEmbedUnimplementedMirroringServer()
}

// UnimplementedMirroringServer must be embedded to have forward compatible implementations.
type UnimplementedMirroringServer struct {
}

func (UnimplementedMirroringServer) Register(context.Context, *payload.Mirroring_Request) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedMirroringServer) mustEmbedUnimplementedMirroringServer() {}

// UnsafeMirroringServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MirroringServer will
// result in compilation errors.
type UnsafeMirroringServer interface {
	mustEmbedUnimplementedMirroringServer()
}

func RegisterMirroringServer(s grpc.ServiceRegistrar, srv MirroringServer) {
	s.RegisterService(&Mirroring_ServiceDesc, srv)
}

func _Mirroring_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Mirroring_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MirroringServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mirroring.v1.Mirroring/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MirroringServer).Register(ctx, req.(*payload.Mirroring_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Mirroring_ServiceDesc is the grpc.ServiceDesc for Mirroring service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Mirroring_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mirroring.v1.Mirroring",
	HandlerType: (*MirroringServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Mirroring_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/proto/v1/mirroring/mirroring.proto",
}
