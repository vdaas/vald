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

// InsertClient is the client API for Insert service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InsertClient interface {
	// A method to add a new single vector.
	Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	// A method to add new multiple vectors by bidirectional streaming.
	StreamInsert(ctx context.Context, opts ...grpc.CallOption) (Insert_StreamInsertClient, error)
	// A method to add new multiple vectors in a single request.
	MultiInsert(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type insertClient struct {
	cc grpc.ClientConnInterface
}

func NewInsertClient(cc grpc.ClientConnInterface) InsertClient {
	return &insertClient{cc}
}

func (c *insertClient) Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := payload.Object_LocationFromVTPool()
	err := c.cc.Invoke(ctx, "/vald.v1.Insert/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *insertClient) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (Insert_StreamInsertClient, error) {
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

func (c *insertClient) MultiInsert(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := payload.Object_LocationsFromVTPool()
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
	// A method to add a new single vector.
	Insert(context.Context, *payload.Insert_Request) (*payload.Object_Location, error)
	// A method to add new multiple vectors by bidirectional streaming.
	StreamInsert(Insert_StreamInsertServer) error
	// A method to add new multiple vectors in a single request.
	MultiInsert(context.Context, *payload.Insert_MultiRequest) (*payload.Object_Locations, error)
	mustEmbedUnimplementedInsertServer()
}

// UnimplementedInsertServer must be embedded to have forward compatible implementations.
type UnimplementedInsertServer struct {
}

func (UnimplementedInsertServer) Insert(context.Context, *payload.Insert_Request) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (UnimplementedInsertServer) StreamInsert(Insert_StreamInsertServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamInsert not implemented")
}
func (UnimplementedInsertServer) MultiInsert(context.Context, *payload.Insert_MultiRequest) (*payload.Object_Locations, error) {
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

func _Insert_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := payload.Insert_RequestFromVTPool()
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
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InsertServer).Insert(ctx, req.(*payload.Insert_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Insert_StreamInsert_Handler(srv interface{}, stream grpc.ServerStream) error {
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
	m := payload.Insert_RequestFromVTPool()
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Insert_MultiInsert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := payload.Insert_MultiRequestFromVTPool()
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
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
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
	Metadata: "apis/proto/v1/vald/insert.proto",
}
