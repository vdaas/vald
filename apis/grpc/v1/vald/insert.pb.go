//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	fmt "fmt"
	math "math"

	proto "github.com/gogo/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("apis/proto/v1/vald/insert.proto", fileDescriptor_7c50984be03265a6) }

var fileDescriptor_7c50984be03265a6 = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x41, 0x4a, 0xc4, 0x30,
	0x14, 0x86, 0xcd, 0x2c, 0x66, 0x20, 0xce, 0x42, 0x0a, 0x6e, 0x3a, 0x43, 0x07, 0xeb, 0x42, 0x99,
	0x45, 0x62, 0x75, 0x37, 0xcb, 0xd9, 0x09, 0x8a, 0xe2, 0xa0, 0x0b, 0x37, 0xf2, 0xda, 0x86, 0x1a,
	0x49, 0x9b, 0xd8, 0xa4, 0x05, 0xb7, 0x5e, 0xc1, 0x03, 0x78, 0x1d, 0x97, 0x82, 0x17, 0x90, 0xe2,
	0x41, 0xa4, 0x49, 0x3b, 0x28, 0x0c, 0xe2, 0x2a, 0xe1, 0xfd, 0xef, 0xff, 0xde, 0x0f, 0x3f, 0x9e,
	0x81, 0xe2, 0x9a, 0xaa, 0x52, 0x1a, 0x49, 0xeb, 0x88, 0xd6, 0x20, 0x52, 0xca, 0x0b, 0xcd, 0x4a,
	0x43, 0xec, 0xd0, 0x1b, 0xb5, 0x23, 0x52, 0x47, 0xfe, 0xfe, 0xef, 0x4d, 0x05, 0x4f, 0x42, 0x42,
	0xda, 0xbf, 0x6e, 0xdb, 0x9f, 0x66, 0x52, 0x66, 0x82, 0x51, 0x50, 0x9c, 0x42, 0x51, 0x48, 0x03,
	0x86, 0xcb, 0x42, 0x3b, 0xf5, 0xf8, 0x75, 0x80, 0x87, 0xa7, 0x16, 0xee, 0x5d, 0xaf, 0x7f, 0x3e,
	0xe9, 0x11, 0x75, 0x44, 0xdc, 0x8c, 0x5c, 0xb1, 0xc7, 0x8a, 0x69, 0xe3, 0x4f, 0x7e, 0x6a, 0x17,
	0xf1, 0x03, 0x4b, 0x0c, 0x39, 0x93, 0x89, 0x85, 0x86, 0xde, 0xf3, 0xc7, 0xd7, 0xcb, 0x60, 0x1c,
	0x8e, 0xba, 0xc0, 0x0b, 0x34, 0xf7, 0x56, 0x78, 0xbc, 0x32, 0x25, 0x83, 0xfc, 0x1f, 0xf0, 0xbd,
	0x0d, 0x70, 0x67, 0x5e, 0x9f, 0xd8, 0x3a, 0x44, 0x47, 0xc8, 0xe3, 0x78, 0xfb, 0xbc, 0x12, 0x86,
	0x77, 0xcc, 0xd9, 0x06, 0xa6, 0xd5, 0x7b, 0xf0, 0xf4, 0x8f, 0xd4, 0x3a, 0x9c, 0xd8, 0xd8, 0xbb,
	0xe1, 0x4e, 0x17, 0x9b, 0xe6, 0xad, 0x57, 0x09, 0xb6, 0x40, 0xf3, 0xe5, 0xdd, 0x5b, 0x13, 0xa0,
	0xf7, 0x26, 0x40, 0x9f, 0x4d, 0x80, 0xb0, 0x2f, 0xcb, 0x8c, 0xd4, 0x29, 0x80, 0x26, 0xb6, 0x05,
	0x50, 0xbc, 0x45, 0xb6, 0xff, 0x25, 0xbe, 0x01, 0x91, 0xba, 0xeb, 0x97, 0xe8, 0xf6, 0x20, 0xe3,
	0xe6, 0xbe, 0x8a, 0x49, 0x22, 0x73, 0x6a, 0x0d, 0xae, 0x49, 0x5b, 0x59, 0x56, 0xaa, 0xa4, 0xef,
	0x36, 0x1e, 0xda, 0x26, 0x4e, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff, 0x95, 0xe3, 0xf2, 0xc9, 0xf8,
	0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// InsertClient is the client API for Insert service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type InsertClient interface {
	Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamInsert(ctx context.Context, opts ...grpc.CallOption) (Insert_StreamInsertClient, error)
	MultiInsert(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type insertClient struct {
	cc *grpc.ClientConn
}

func NewInsertClient(cc *grpc.ClientConn) InsertClient {
	return &insertClient{cc}
}

func (c *insertClient) Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/vald.v1.Insert/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *insertClient) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (Insert_StreamInsertClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Insert_serviceDesc.Streams[0], "/vald.v1.Insert/StreamInsert", opts...)
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
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.v1.Insert/MultiInsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InsertServer is the server API for Insert service.
type InsertServer interface {
	Insert(context.Context, *payload.Insert_Request) (*payload.Object_Location, error)
	StreamInsert(Insert_StreamInsertServer) error
	MultiInsert(context.Context, *payload.Insert_MultiRequest) (*payload.Object_Locations, error)
}

// UnimplementedInsertServer can be embedded to have forward compatible implementations.
type UnimplementedInsertServer struct {
}

func (*UnimplementedInsertServer) Insert(ctx context.Context, req *payload.Insert_Request) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (*UnimplementedInsertServer) StreamInsert(srv Insert_StreamInsertServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamInsert not implemented")
}
func (*UnimplementedInsertServer) MultiInsert(ctx context.Context, req *payload.Insert_MultiRequest) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiInsert not implemented")
}

func RegisterInsertServer(s *grpc.Server, srv InsertServer) {
	s.RegisterService(&_Insert_serviceDesc, srv)
}

func _Insert_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
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
	m := new(payload.Insert_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Insert_MultiInsert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
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
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InsertServer).MultiInsert(ctx, req.(*payload.Insert_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Insert_serviceDesc = grpc.ServiceDesc{
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
