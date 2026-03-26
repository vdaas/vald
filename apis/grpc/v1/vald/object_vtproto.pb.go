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

// ObjectClient is the client API for Object service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ObjectClient interface {
	// Overview
	// Exists RPC is the method to check that a vector exists in the `vald-agent`.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_ID, error)
	// Overview
	// GetObject RPC is the method to get the metadata of a vector inserted into the `vald-agent`.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	GetObject(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (*payload.Object_Vector, error)
	// Overview
	// StreamGetObject RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the GetObject request can be communicated in any order between client and server.
	// Each Upsert request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (Object_StreamGetObjectClient, error)
	// Overview
	// A method to get all the vectors with server streaming
	// ---
	// Status Code
	// TODO
	// ---
	// Troubleshooting
	// TODO
	StreamListObject(ctx context.Context, in *payload.Object_List_Request, opts ...grpc.CallOption) (Object_StreamListObjectClient, error)
	// Overview
	// Represent the RPC to get the vector metadata. This RPC is mainly used for index correction process
	// ---
	// Status Code
	// TODO
	// ---
	// Troubleshooting
	// TODO
	GetTimestamp(ctx context.Context, in *payload.Object_TimestampRequest, opts ...grpc.CallOption) (*payload.Object_Timestamp, error)
}

type objectClient struct {
	cc grpc.ClientConnInterface
}

func NewObjectClient(cc grpc.ClientConnInterface) ObjectClient {
	return &objectClient{cc}
}

func (c *objectClient) Exists(
	ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption,
) (*payload.Object_ID, error) {
	out := new(payload.Object_ID)
	err := c.cc.Invoke(ctx, "/vald.v1.Object/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectClient) GetObject(
	ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/vald.v1.Object/GetObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectClient) StreamGetObject(
	ctx context.Context, opts ...grpc.CallOption,
) (Object_StreamGetObjectClient, error) {
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

func (c *objectClient) StreamListObject(
	ctx context.Context, in *payload.Object_List_Request, opts ...grpc.CallOption,
) (Object_StreamListObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &Object_ServiceDesc.Streams[1], "/vald.v1.Object/StreamListObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &objectStreamListObjectClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Object_StreamListObjectClient interface {
	Recv() (*payload.Object_List_Response, error)
	grpc.ClientStream
}

type objectStreamListObjectClient struct {
	grpc.ClientStream
}

func (x *objectStreamListObjectClient) Recv() (*payload.Object_List_Response, error) {
	m := new(payload.Object_List_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *objectClient) GetTimestamp(
	ctx context.Context, in *payload.Object_TimestampRequest, opts ...grpc.CallOption,
) (*payload.Object_Timestamp, error) {
	out := new(payload.Object_Timestamp)
	err := c.cc.Invoke(ctx, "/vald.v1.Object/GetTimestamp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ObjectServer is the server API for Object service.
// All implementations must embed UnimplementedObjectServer
// for forward compatibility
type ObjectServer interface {
	// Overview
	// Exists RPC is the method to check that a vector exists in the `vald-agent`.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	Exists(context.Context, *payload.Object_ID) (*payload.Object_ID, error)
	// Overview
	// GetObject RPC is the method to get the metadata of a vector inserted into the `vald-agent`.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	GetObject(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error)
	// Overview
	// StreamGetObject RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the GetObject request can be communicated in any order between client and server.
	// Each Upsert request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	StreamGetObject(Object_StreamGetObjectServer) error
	// Overview
	// A method to get all the vectors with server streaming
	// ---
	// Status Code
	// TODO
	// ---
	// Troubleshooting
	// TODO
	StreamListObject(*payload.Object_List_Request, Object_StreamListObjectServer) error
	// Overview
	// Represent the RPC to get the vector metadata. This RPC is mainly used for index correction process
	// ---
	// Status Code
	// TODO
	// ---
	// Troubleshooting
	// TODO
	GetTimestamp(context.Context, *payload.Object_TimestampRequest) (*payload.Object_Timestamp, error)
	mustEmbedUnimplementedObjectServer()
}

// UnimplementedObjectServer must be embedded to have forward compatible implementations.
type UnimplementedObjectServer struct{}

func (UnimplementedObjectServer) Exists(
	context.Context, *payload.Object_ID,
) (*payload.Object_ID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}

func (UnimplementedObjectServer) GetObject(
	context.Context, *payload.Object_VectorRequest,
) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetObject not implemented")
}

func (UnimplementedObjectServer) StreamGetObject(Object_StreamGetObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamGetObject not implemented")
}

func (UnimplementedObjectServer) StreamListObject(
	*payload.Object_List_Request, Object_StreamListObjectServer,
) error {
	return status.Errorf(codes.Unimplemented, "method StreamListObject not implemented")
}

func (UnimplementedObjectServer) GetTimestamp(
	context.Context, *payload.Object_TimestampRequest,
) (*payload.Object_Timestamp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTimestamp not implemented")
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

func _Object_Exists_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Object_ID)
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
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(ObjectServer).Exists(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Object_GetObject_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Object_VectorRequest)
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
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(ObjectServer).GetObject(ctx, req.(*payload.Object_VectorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Object_StreamGetObject_Handler(srv any, stream grpc.ServerStream) error {
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
	m := new(payload.Object_VectorRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Object_StreamListObject_Handler(srv any, stream grpc.ServerStream) error {
	m := new(payload.Object_List_Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ObjectServer).StreamListObject(m, &objectStreamListObjectServer{stream})
}

type Object_StreamListObjectServer interface {
	Send(*payload.Object_List_Response) error
	grpc.ServerStream
}

type objectStreamListObjectServer struct {
	grpc.ServerStream
}

func (x *objectStreamListObjectServer) Send(m *payload.Object_List_Response) error {
	return x.ServerStream.SendMsg(m)
}

func _Object_GetTimestamp_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Object_TimestampRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectServer).GetTimestamp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Object/GetTimestamp",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(ObjectServer).GetTimestamp(ctx, req.(*payload.Object_TimestampRequest))
	}
	return interceptor(ctx, in, info, handler)
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
		{
			MethodName: "GetTimestamp",
			Handler:    _Object_GetTimestamp_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamGetObject",
			Handler:       _Object_StreamGetObject_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamListObject",
			Handler:       _Object_StreamListObject_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "v1/vald/object.proto",
}
