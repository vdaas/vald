//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// ObjectClient is the client API for Object service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ObjectClient interface {
	// A method to check whether a specified ID is indexed or not.
	Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_ID, error)
	// A method to fetch a vector.
	GetObject(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (*payload.Object_Vector, error)
	// A method to fetch vectors by bidirectional streaming.
	StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (Object_StreamGetObjectClient, error)
}

type objectClient struct {
	cc grpc.ClientConnInterface
}

func NewObjectClient(cc grpc.ClientConnInterface) ObjectClient {
	return &objectClient{cc}
}

func (c *objectClient) Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_ID, error) {
	out := payload.Object_IDFromVTPool()
	err := c.cc.Invoke(ctx, "/vald.v1.Object/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectClient) GetObject(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	out := payload.Object_VectorFromVTPool()
	err := c.cc.Invoke(ctx, "/vald.v1.Object/GetObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectClient) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (Object_StreamGetObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &Object_ServiceDesc.Streams[0], "/vald.v1.Object/StreamGetObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &objectStreamGetObjectClient{stream}
	return x, nil
}

type Object_StreamGetObjectClient interface {
	Send(*payload.Object_VectorRequest) error
	Recv() (*payload.Object_StreamVector, error)
	grpc.ClientStream
}

type objectStreamGetObjectClient struct {
	grpc.ClientStream
}

func (x *objectStreamGetObjectClient) Send(m *payload.Object_VectorRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *objectStreamGetObjectClient) Recv() (*payload.Object_StreamVector, error) {
	m := new(payload.Object_StreamVector)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ObjectServer is the server API for Object service.
// All implementations must embed UnimplementedObjectServer
// for forward compatibility
type ObjectServer interface {
	// A method to check whether a specified ID is indexed or not.
	Exists(context.Context, *payload.Object_ID) (*payload.Object_ID, error)
	// A method to fetch a vector.
	GetObject(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error)
	// A method to fetch vectors by bidirectional streaming.
	StreamGetObject(Object_StreamGetObjectServer) error
	mustEmbedUnimplementedObjectServer()
}

// UnimplementedObjectServer must be embedded to have forward compatible implementations.
type UnimplementedObjectServer struct {
}

func (UnimplementedObjectServer) Exists(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (UnimplementedObjectServer) GetObject(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetObject not implemented")
}
func (UnimplementedObjectServer) StreamGetObject(Object_StreamGetObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamGetObject not implemented")
}
func (UnimplementedObjectServer) mustEmbedUnimplementedObjectServer() {}

// UnsafeObjectServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ObjectServer will
// result in compilation errors.
type UnsafeObjectServer interface {
	mustEmbedUnimplementedObjectServer()
}

func RegisterObjectServer(s grpc.ServiceRegistrar, srv ObjectServer) {
	s.RegisterService(&Object_ServiceDesc, srv)
}

func _Object_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := payload.Object_IDFromVTPool()
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Object/Exists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectServer).Exists(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Object_GetObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := payload.Object_VectorRequestFromVTPool()
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectServer).GetObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Object/GetObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectServer).GetObject(ctx, req.(*payload.Object_VectorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Object_StreamGetObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ObjectServer).StreamGetObject(&objectStreamGetObjectServer{stream})
}

type Object_StreamGetObjectServer interface {
	Send(*payload.Object_StreamVector) error
	Recv() (*payload.Object_VectorRequest, error)
	grpc.ServerStream
}

type objectStreamGetObjectServer struct {
	grpc.ServerStream
}

func (x *objectStreamGetObjectServer) Send(m *payload.Object_StreamVector) error {
	return x.ServerStream.SendMsg(m)
}

func (x *objectStreamGetObjectServer) Recv() (*payload.Object_VectorRequest, error) {
	m := payload.Object_VectorRequestFromVTPool()
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Object_ServiceDesc is the grpc.ServiceDesc for Object service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Object_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vald.v1.Object",
	HandlerType: (*ObjectServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Exists",
			Handler:    _Object_Exists_Handler,
		},
		{
			MethodName: "GetObject",
			Handler:    _Object_GetObject_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamGetObject",
			Handler:       _Object_StreamGetObject_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "apis/proto/v1/vald/object.proto",
}
