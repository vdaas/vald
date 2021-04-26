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

	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	golang_proto "github.com/golang/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	codes "github.com/vdaas/vald/internal/net/grpc/codes"
	status "github.com/vdaas/vald/internal/net/grpc/status"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("apis/proto/v1/vald/insert.proto", fileDescriptor_7c50984be03265a6) }
func init() {
	golang_proto.RegisterFile("apis/proto/v1/vald/insert.proto", fileDescriptor_7c50984be03265a6)
}

var fileDescriptor_7c50984be03265a6 = []byte{
	// 339 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xc1, 0x4a, 0xf3, 0x40,
	0x18, 0xfc, 0xb7, 0x87, 0x16, 0xf2, 0xf7, 0x20, 0x01, 0x2f, 0xa9, 0xa4, 0x18, 0x0f, 0x4a, 0xc1,
	0x5d, 0xab, 0x07, 0xa1, 0xc7, 0xde, 0x04, 0x45, 0xb1, 0xe8, 0xc1, 0xdb, 0x97, 0x64, 0x5d, 0x57,
	0xd2, 0x7c, 0x6b, 0xb2, 0x09, 0x78, 0xf5, 0x15, 0x7c, 0x00, 0x5f, 0xc5, 0x63, 0x8f, 0x82, 0x2f,
	0x20, 0xad, 0x0f, 0x22, 0xd9, 0x4d, 0x4a, 0x95, 0x22, 0x9e, 0x32, 0xf9, 0x76, 0x66, 0x76, 0x96,
	0xf9, 0x9c, 0x3e, 0x28, 0x99, 0x33, 0x95, 0xa1, 0x46, 0x56, 0x0e, 0x59, 0x09, 0x49, 0xcc, 0x64,
	0x9a, 0xf3, 0x4c, 0x53, 0x33, 0x74, 0x3b, 0xd5, 0x88, 0x96, 0x43, 0x6f, 0xe7, 0x3b, 0x53, 0xc1,
	0x63, 0x82, 0x10, 0x37, 0x5f, 0xcb, 0xf6, 0xf6, 0x85, 0xd4, 0x77, 0x45, 0x48, 0x23, 0x9c, 0x32,
	0x81, 0x02, 0x2d, 0x3f, 0x2c, 0x6e, 0xcd, 0x9f, 0x15, 0x57, 0xa8, 0xa6, 0x1f, 0xff, 0xa4, 0x0b,
	0x44, 0x91, 0x70, 0x73, 0x93, 0x85, 0x0c, 0x94, 0x64, 0x90, 0xa6, 0xa8, 0x41, 0x4b, 0x4c, 0x73,
	0x2b, 0x3c, 0x7c, 0x69, 0x39, 0xed, 0x13, 0x13, 0xd3, 0xbd, 0x5a, 0x22, 0x8f, 0x36, 0x61, 0xca,
	0x21, 0xb5, 0x33, 0x7a, 0xc9, 0x1f, 0x0a, 0x9e, 0x6b, 0xaf, 0xb7, 0x7a, 0x76, 0x1e, 0xde, 0xf3,
	0x48, 0xd3, 0x53, 0x8c, 0x8c, 0x69, 0xe0, 0x3e, 0xbd, 0x7f, 0x3e, 0xb7, 0xba, 0x41, 0xa7, 0x7e,
	0xfa, 0x88, 0x0c, 0xdc, 0x89, 0xd3, 0x9d, 0xe8, 0x8c, 0xc3, 0xf4, 0x0f, 0xe6, 0xdb, 0x6b, 0xcc,
	0xad, 0x78, 0x79, 0xc5, 0xbf, 0x3d, 0x72, 0x40, 0x5c, 0xe9, 0xfc, 0x3f, 0x2b, 0x12, 0x2d, 0x6b,
	0xcf, 0xfe, 0x1a, 0x4f, 0x73, 0xde, 0x18, 0x6f, 0xfd, 0x92, 0x3a, 0x0f, 0x7a, 0x26, 0xf6, 0x66,
	0xb0, 0x51, 0xc7, 0x66, 0xd3, 0x4a, 0xab, 0x12, 0x3e, 0x22, 0x83, 0xb1, 0x98, 0xcd, 0x7d, 0xf2,
	0x36, 0xf7, 0xc9, 0xc7, 0xdc, 0x27, 0xaf, 0x0b, 0x9f, 0xcc, 0x16, 0x3e, 0x71, 0x3c, 0xcc, 0x04,
	0x2d, 0x63, 0x80, 0x9c, 0x9a, 0x5e, 0x41, 0xc9, 0xca, 0xba, 0xc2, 0x63, 0xe7, 0x1a, 0x92, 0xd8,
	0xa6, 0xb8, 0x20, 0x37, 0xbb, 0x2b, 0xd5, 0x18, 0x81, 0xdd, 0x0d, 0x5b, 0x4d, 0xa6, 0xa2, 0x66,
	0x5b, 0xc2, 0xb6, 0x69, 0xe4, 0xe8, 0x2b, 0x00, 0x00, 0xff, 0xff, 0xcc, 0x1c, 0x60, 0x0e, 0x4a,
	0x02, 0x00, 0x00,
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
