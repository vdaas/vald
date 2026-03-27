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

// InsertClient is the client API for Insert service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InsertClient interface {
	// Overview
	// Insert RPC is the method to add a new single vector.
	// ---
	// Status Code
	// | 0    | OK                |
	// | 1    | CANCELLED         |
	// | 3    | INVALID_ARGUMENT  |
	// | 4    | DEADLINE_EXCEEDED |
	// | 5    | NOT_FOUND         |
	// | 13   | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	// Overview
	// StreamInsert RPC is the method to add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
	// Each Insert request and response are independent.
	// It's the recommended method to insert a large number of vectors.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	StreamInsert(ctx context.Context, opts ...grpc.CallOption) (Insert_StreamInsertClient, error)
	// Overview
	// MultiInsert RPC is the method to add multiple new vectors in **1** request.
	//
	// <div class="notice">
	// gRPC has a message size limitation.<br>
	// Please be careful that the size of the request exceeds the limit.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	MultiInsert(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type insertClient struct {
	cc grpc.ClientConnInterface
}

func NewInsertClient(cc grpc.ClientConnInterface) InsertClient {
	return &insertClient{cc}
}

func (c *insertClient) Insert(
	ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/vald.v1.Insert/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *insertClient) StreamInsert(
	ctx context.Context, opts ...grpc.CallOption,
) (Insert_StreamInsertClient, error) {
	stream, err := c.cc.NewStream(ctx, &Insert_ServiceDesc.Streams[0], "/vald.v1.Insert/StreamInsert", opts...)
	if err != nil {
		return nil, err
	}
	x := &insertStreamInsertClient{stream}
	return x, nil
}

type Insert_StreamInsertClient interface {
	Send(*payload.Insert_Request) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type insertStreamInsertClient struct {
	grpc.ClientStream
}

func (x *insertStreamInsertClient) Send(m *payload.Insert_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *insertStreamInsertClient) Recv() (*payload.Object_StreamLocation, error) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *insertClient) MultiInsert(
	ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.v1.Insert/MultiInsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InsertServer is the server API for Insert service.
// All implementations must embed UnimplementedInsertServer
// for forward compatibility
type InsertServer interface {
	// Overview
	// Insert RPC is the method to add a new single vector.
	// ---
	// Status Code
	// | 0    | OK                |
	// | 1    | CANCELLED         |
	// | 3    | INVALID_ARGUMENT  |
	// | 4    | DEADLINE_EXCEEDED |
	// | 5    | NOT_FOUND         |
	// | 13   | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	Insert(context.Context, *payload.Insert_Request) (*payload.Object_Location, error)
	// Overview
	// StreamInsert RPC is the method to add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
	// Each Insert request and response are independent.
	// It's the recommended method to insert a large number of vectors.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	StreamInsert(Insert_StreamInsertServer) error
	// Overview
	// MultiInsert RPC is the method to add multiple new vectors in **1** request.
	//
	// <div class="notice">
	// gRPC has a message size limitation.<br>
	// Please be careful that the size of the request exceeds the limit.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	MultiInsert(context.Context, *payload.Insert_MultiRequest) (*payload.Object_Locations, error)
	mustEmbedUnimplementedInsertServer()
}

// UnimplementedInsertServer must be embedded to have forward compatible implementations.
type UnimplementedInsertServer struct{}

func (UnimplementedInsertServer) Insert(
	context.Context, *payload.Insert_Request,
) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}

func (UnimplementedInsertServer) StreamInsert(Insert_StreamInsertServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamInsert not implemented")
}

func (UnimplementedInsertServer) MultiInsert(
	context.Context, *payload.Insert_MultiRequest,
) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiInsert not implemented")
}
func (UnimplementedInsertServer) mustEmbedUnimplementedInsertServer() {}

// UnsafeInsertServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InsertServer will
// result in compilation errors.
type UnsafeInsertServer interface {
	mustEmbedUnimplementedInsertServer()
}

func RegisterInsertServer(s grpc.ServiceRegistrar, srv InsertServer) {
	s.RegisterService(&Insert_ServiceDesc, srv)
}

func _Insert_Insert_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Insert_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InsertServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Insert/Insert",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(InsertServer).Insert(ctx, req.(*payload.Insert_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Insert_StreamInsert_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(InsertServer).StreamInsert(&insertStreamInsertServer{stream})
}

type Insert_StreamInsertServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Insert_Request, error)
	grpc.ServerStream
}

type insertStreamInsertServer struct {
	grpc.ServerStream
}

func (x *insertStreamInsertServer) Send(m *payload.Object_StreamLocation) error {
	return x.ServerStream.SendMsg(m)
}

func (x *insertStreamInsertServer) Recv() (*payload.Insert_Request, error) {
	m := new(payload.Insert_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Insert_MultiInsert_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Insert_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InsertServer).MultiInsert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Insert/MultiInsert",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(InsertServer).MultiInsert(ctx, req.(*payload.Insert_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Insert_ServiceDesc is the grpc.ServiceDesc for Insert service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Insert_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vald.v1.Insert",
	HandlerType: (*InsertServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Insert",
			Handler:    _Insert_Insert_Handler,
		},
		{
			MethodName: "MultiInsert",
			Handler:    _Insert_MultiInsert_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamInsert",
			Handler:       _Insert_StreamInsert_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "v1/vald/insert.proto",
}
