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

// UpdateClient is the client API for Update service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UpdateClient interface {
	// A method to update an indexed vector.
	Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	// A method to update multiple indexed vectors by bidirectional streaming.
	StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (Update_StreamUpdateClient, error)
	// A method to update multiple indexed vectors in a single request.
	MultiUpdate(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type updateClient struct {
	cc grpc.ClientConnInterface
}

func NewUpdateClient(cc grpc.ClientConnInterface) UpdateClient {
	return &updateClient{cc}
}

func (c *updateClient) Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/vald.v1.Update/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *updateClient) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (Update_StreamUpdateClient, error) {
	stream, err := c.cc.NewStream(ctx, &Update_ServiceDesc.Streams[0], "/vald.v1.Update/StreamUpdate", opts...)
	if err != nil {
		return nil, err
	}
	x := &updateStreamUpdateClient{stream}
	return x, nil
}

type Update_StreamUpdateClient interface {
	Send(*payload.Update_Request) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type updateStreamUpdateClient struct {
	grpc.ClientStream
}

func (x *updateStreamUpdateClient) Send(m *payload.Update_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *updateStreamUpdateClient) Recv() (*payload.Object_StreamLocation, error) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *updateClient) MultiUpdate(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.v1.Update/MultiUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateServer is the server API for Update service.
// All implementations must embed UnimplementedUpdateServer
// for forward compatibility
type UpdateServer interface {
	// A method to update an indexed vector.
	Update(context.Context, *payload.Update_Request) (*payload.Object_Location, error)
	// A method to update multiple indexed vectors by bidirectional streaming.
	StreamUpdate(Update_StreamUpdateServer) error
	// A method to update multiple indexed vectors in a single request.
	MultiUpdate(context.Context, *payload.Update_MultiRequest) (*payload.Object_Locations, error)
	mustEmbedUnimplementedUpdateServer()
}

// UnimplementedUpdateServer must be embedded to have forward compatible implementations.
type UnimplementedUpdateServer struct {
}

func (UnimplementedUpdateServer) Update(context.Context, *payload.Update_Request) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedUpdateServer) StreamUpdate(Update_StreamUpdateServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpdate not implemented")
}
func (UnimplementedUpdateServer) MultiUpdate(context.Context, *payload.Update_MultiRequest) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpdate not implemented")
}
func (UnimplementedUpdateServer) mustEmbedUnimplementedUpdateServer() {}

// UnsafeUpdateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UpdateServer will
// result in compilation errors.
type UnsafeUpdateServer interface {
	mustEmbedUnimplementedUpdateServer()
}

func RegisterUpdateServer(s grpc.ServiceRegistrar, srv UpdateServer) {
	s.RegisterService(&Update_ServiceDesc, srv)
}

func _Update_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Update_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpdateServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Update/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UpdateServer).Update(ctx, req.(*payload.Update_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Update_StreamUpdate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UpdateServer).StreamUpdate(&updateStreamUpdateServer{stream})
}

type Update_StreamUpdateServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Update_Request, error)
	grpc.ServerStream
}

type updateStreamUpdateServer struct {
	grpc.ServerStream
}

func (x *updateStreamUpdateServer) Send(m *payload.Object_StreamLocation) error {
	return x.ServerStream.SendMsg(m)
}

func (x *updateStreamUpdateServer) Recv() (*payload.Update_Request, error) {
	m := new(payload.Update_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Update_MultiUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Update_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpdateServer).MultiUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Update/MultiUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UpdateServer).MultiUpdate(ctx, req.(*payload.Update_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Update_ServiceDesc is the grpc.ServiceDesc for Update service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Update_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vald.v1.Update",
	HandlerType: (*UpdateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Update",
			Handler:    _Update_Update_Handler,
		},
		{
			MethodName: "MultiUpdate",
			Handler:    _Update_MultiUpdate_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamUpdate",
			Handler:       _Update_StreamUpdate_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "v1/vald/update.proto",
}
