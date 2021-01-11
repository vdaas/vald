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
	// 300 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0x2f, 0x4b, 0xcc, 0x49, 0xd1, 0xcf,
	0xcc, 0x2b, 0x4e, 0x2d, 0x2a, 0xd1, 0x03, 0x0b, 0x0a, 0xb1, 0x83, 0x84, 0xf4, 0xca, 0x0c, 0xa5,
	0x94, 0x51, 0x55, 0x16, 0x24, 0x56, 0xe6, 0xe4, 0x27, 0xa6, 0xc0, 0x68, 0x88, 0x6a, 0x29, 0x99,
	0xf4, 0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0xfd, 0xc4, 0x82, 0x4c, 0xfd, 0xc4, 0xbc, 0xbc, 0xfc, 0x92,
	0xc4, 0x92, 0xcc, 0xfc, 0xbc, 0x62, 0x88, 0xac, 0xd1, 0x4c, 0x26, 0x2e, 0x36, 0x4f, 0xb0, 0xe1,
	0x42, 0xa1, 0x70, 0x96, 0x94, 0x1e, 0xcc, 0x88, 0x32, 0x43, 0x3d, 0x88, 0x98, 0x5e, 0x50, 0x6a,
	0x61, 0x69, 0x6a, 0x71, 0x89, 0x94, 0x34, 0xb2, 0x9c, 0x7f, 0x52, 0x56, 0x6a, 0x72, 0x89, 0x9e,
	0x4f, 0x7e, 0x32, 0xd8, 0x50, 0x25, 0xa1, 0xa6, 0xcb, 0x4f, 0x26, 0x33, 0xf1, 0x28, 0xb1, 0x43,
	0x1d, 0x6c, 0xc5, 0xa8, 0x25, 0xe4, 0xcb, 0xc5, 0x13, 0x5c, 0x52, 0x94, 0x9a, 0x98, 0x4b, 0xa9,
	0xe1, 0x0c, 0x1a, 0x8c, 0x06, 0x8c, 0x42, 0x99, 0x5c, 0xdc, 0xbe, 0xa5, 0x39, 0x25, 0x99, 0x50,
	0xd3, 0xe4, 0xb1, 0x98, 0x06, 0x96, 0x87, 0x19, 0x29, 0x83, 0xc7, 0xc8, 0x62, 0x25, 0x69, 0xb0,
	0x83, 0x45, 0x95, 0x04, 0xa0, 0x0e, 0xd6, 0xcf, 0x05, 0xe9, 0x2d, 0xc8, 0x49, 0xb5, 0x62, 0xd4,
	0x72, 0x8a, 0x3f, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0xb9,
	0xa4, 0xf2, 0x8b, 0xd2, 0xf5, 0xca, 0x52, 0x12, 0x13, 0x8b, 0xf5, 0xc0, 0xe1, 0x9f, 0x58, 0x90,
	0x09, 0x32, 0x12, 0xc4, 0x76, 0xe2, 0x0a, 0x4b, 0xcc, 0x49, 0x81, 0xd8, 0x1e, 0xc0, 0x18, 0xa5,
	0x9e, 0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c, 0x9f, 0xab, 0x0f, 0xd6, 0x00, 0x89, 0x43,
	0x70, 0x64, 0xa5, 0x17, 0x15, 0x24, 0xc3, 0x62, 0x35, 0x89, 0x0d, 0x1c, 0x07, 0xc6, 0x80, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xc8, 0x14, 0x43, 0x96, 0xf2, 0x01, 0x00, 0x00,
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
	Recv() (*payload.Object_Location, error)
	grpc.ClientStream
}

type insertStreamInsertClient struct {
	grpc.ClientStream
}

func (x *insertStreamInsertClient) Send(m *payload.Insert_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *insertStreamInsertClient) Recv() (*payload.Object_Location, error) {
	m := new(payload.Object_Location)
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
	Send(*payload.Object_Location) error
	Recv() (*payload.Insert_Request, error)
	grpc.ServerStream
}

type insertStreamInsertServer struct {
	grpc.ServerStream
}

func (x *insertStreamInsertServer) Send(m *payload.Object_Location) error {
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
