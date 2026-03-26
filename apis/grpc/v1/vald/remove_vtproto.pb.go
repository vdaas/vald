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

// RemoveClient is the client API for Remove service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RemoveClient interface {
	// Overview
	// Remove RPC is the method to remove a single vector.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
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
	Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	// Overview
	// RemoveByTimestamp RPC is the method to remove vectors based on timestamp.
	//
	// <div class="notice">
	// In the TimestampRequest message, the 'timestamps' field is repeated, allowing the inclusion of multiple Timestamp.<br>
	// When multiple Timestamps are provided, it results in an `AND` condition, enabling the realization of deletions with specified ranges.<br>
	// This design allows for versatile deletion operations, facilitating tasks such as removing data within a specific time range.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                                                       |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.                              |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed.                             |
	// | NOT_FOUND         | No vectors in the system match the specified timestamp conditions.                              | Check whether vectors matching the specified timestamp conditions exist in the system, and fix conditions if needed. |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.
	RemoveByTimestamp(ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
	// Overview
	// A method to remove multiple indexed vectors by bidirectional streaming.
	//
	// StreamRemove RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the remove request can be communicated in any order between client and server.
	// Each Remove request and response are independent.
	// It's the recommended method to remove a large number of vectors.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	//
	//	The request process may not be completed when the response code is NOT `0 (OK)`.
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
	StreamRemove(ctx context.Context, opts ...grpc.CallOption) (Remove_StreamRemoveClient, error)
	// Overview
	// MultiRemove is the method to remove multiple vectors in **1** request.
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
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
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
	MultiRemove(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type removeClient struct {
	cc grpc.ClientConnInterface
}

func NewRemoveClient(cc grpc.ClientConnInterface) RemoveClient {
	return &removeClient{cc}
}

func (c *removeClient) Remove(
	ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/vald.v1.Remove/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *removeClient) RemoveByTimestamp(
	ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.v1.Remove/RemoveByTimestamp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *removeClient) StreamRemove(
	ctx context.Context, opts ...grpc.CallOption,
) (Remove_StreamRemoveClient, error) {
	stream, err := c.cc.NewStream(ctx, &Remove_ServiceDesc.Streams[0], "/vald.v1.Remove/StreamRemove", opts...)
	if err != nil {
		return nil, err
	}
	x := &removeStreamRemoveClient{stream}
	return x, nil
}

type Remove_StreamRemoveClient interface {
	Send(*payload.Remove_Request) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type removeStreamRemoveClient struct {
	grpc.ClientStream
}

func (x *removeStreamRemoveClient) Send(m *payload.Remove_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *removeStreamRemoveClient) Recv() (*payload.Object_StreamLocation, error) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *removeClient) MultiRemove(
	ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.v1.Remove/MultiRemove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveServer is the server API for Remove service.
// All implementations must embed UnimplementedRemoveServer
// for forward compatibility
type RemoveServer interface {
	// Overview
	// Remove RPC is the method to remove a single vector.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
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
	Remove(context.Context, *payload.Remove_Request) (*payload.Object_Location, error)
	// Overview
	// RemoveByTimestamp RPC is the method to remove vectors based on timestamp.
	//
	// <div class="notice">
	// In the TimestampRequest message, the 'timestamps' field is repeated, allowing the inclusion of multiple Timestamp.<br>
	// When multiple Timestamps are provided, it results in an `AND` condition, enabling the realization of deletions with specified ranges.<br>
	// This design allows for versatile deletion operations, facilitating tasks such as removing data within a specific time range.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                                                       |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.                              |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed.                             |
	// | NOT_FOUND         | No vectors in the system match the specified timestamp conditions.                              | Check whether vectors matching the specified timestamp conditions exist in the system, and fix conditions if needed. |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.
	RemoveByTimestamp(context.Context, *payload.Remove_TimestampRequest) (*payload.Object_Locations, error)
	// Overview
	// A method to remove multiple indexed vectors by bidirectional streaming.
	//
	// StreamRemove RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the remove request can be communicated in any order between client and server.
	// Each Remove request and response are independent.
	// It's the recommended method to remove a large number of vectors.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	//
	//	The request process may not be completed when the response code is NOT `0 (OK)`.
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
	StreamRemove(Remove_StreamRemoveServer) error
	// Overview
	// MultiRemove is the method to remove multiple vectors in **1** request.
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
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
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
	MultiRemove(context.Context, *payload.Remove_MultiRequest) (*payload.Object_Locations, error)
	mustEmbedUnimplementedRemoveServer()
}

// UnimplementedRemoveServer must be embedded to have forward compatible implementations.
type UnimplementedRemoveServer struct{}

func (UnimplementedRemoveServer) Remove(
	context.Context, *payload.Remove_Request,
) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}

func (UnimplementedRemoveServer) RemoveByTimestamp(
	context.Context, *payload.Remove_TimestampRequest,
) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveByTimestamp not implemented")
}

func (UnimplementedRemoveServer) StreamRemove(Remove_StreamRemoveServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamRemove not implemented")
}

func (UnimplementedRemoveServer) MultiRemove(
	context.Context, *payload.Remove_MultiRequest,
) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiRemove not implemented")
}
func (UnimplementedRemoveServer) mustEmbedUnimplementedRemoveServer() {}

// UnsafeRemoveServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RemoveServer will
// result in compilation errors.
type UnsafeRemoveServer interface {
	mustEmbedUnimplementedRemoveServer()
}

func RegisterRemoveServer(s grpc.ServiceRegistrar, srv RemoveServer) {
	s.RegisterService(&Remove_ServiceDesc, srv)
}

func _Remove_Remove_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Remove_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoveServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Remove/Remove",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(RemoveServer).Remove(ctx, req.(*payload.Remove_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Remove_RemoveByTimestamp_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Remove_TimestampRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoveServer).RemoveByTimestamp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Remove/RemoveByTimestamp",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(RemoveServer).RemoveByTimestamp(ctx, req.(*payload.Remove_TimestampRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Remove_StreamRemove_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(RemoveServer).StreamRemove(&removeStreamRemoveServer{stream})
}

type Remove_StreamRemoveServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Remove_Request, error)
	grpc.ServerStream
}

type removeStreamRemoveServer struct {
	grpc.ServerStream
}

func (x *removeStreamRemoveServer) Send(m *payload.Object_StreamLocation) error {
	return x.ServerStream.SendMsg(m)
}

func (x *removeStreamRemoveServer) Recv() (*payload.Remove_Request, error) {
	m := new(payload.Remove_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Remove_MultiRemove_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Remove_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoveServer).MultiRemove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Remove/MultiRemove",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(RemoveServer).MultiRemove(ctx, req.(*payload.Remove_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Remove_ServiceDesc is the grpc.ServiceDesc for Remove service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Remove_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vald.v1.Remove",
	HandlerType: (*RemoveServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Remove",
			Handler:    _Remove_Remove_Handler,
		},
		{
			MethodName: "RemoveByTimestamp",
			Handler:    _Remove_RemoveByTimestamp_Handler,
		},
		{
			MethodName: "MultiRemove",
			Handler:    _Remove_MultiRemove_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamRemove",
			Handler:       _Remove_StreamRemove_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "v1/vald/remove.proto",
}
