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

func init() { proto.RegisterFile("apis/proto/v1/vald/object.proto", fileDescriptor_f3068a4c11e32302) }
func init() {
	golang_proto.RegisterFile("apis/proto/v1/vald/object.proto", fileDescriptor_f3068a4c11e32302)
}

var fileDescriptor_f3068a4c11e32302 = []byte{
	// 334 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x31, 0x4b, 0x33, 0x41,
	0x10, 0xfd, 0x36, 0x45, 0x3e, 0x5c, 0xc4, 0xe0, 0xa1, 0x88, 0x57, 0x5c, 0x20, 0x16, 0xda, 0xb8,
	0x6b, 0xb4, 0xb0, 0x0f, 0x8a, 0x08, 0x82, 0xa2, 0x92, 0x42, 0xab, 0xb9, 0xbb, 0x75, 0x5d, 0xb9,
	0x64, 0xce, 0xdb, 0xcd, 0xa1, 0x88, 0x8d, 0x9d, 0xb5, 0x7f, 0xc8, 0x32, 0xa5, 0xe0, 0x1f, 0x90,
	0xc4, 0x1f, 0x22, 0xbb, 0x7b, 0x01, 0x95, 0x2b, 0xac, 0x76, 0x76, 0xe6, 0xbd, 0x37, 0xc3, 0xbc,
	0xa1, 0x6d, 0xc8, 0x95, 0xe6, 0x79, 0x81, 0x06, 0x79, 0xd9, 0xe5, 0x25, 0x64, 0x29, 0xc7, 0xf8,
	0x46, 0x24, 0x86, 0xb9, 0x64, 0xf0, 0xdf, 0xa6, 0x58, 0xd9, 0x0d, 0xd7, 0x7e, 0x22, 0x73, 0xb8,
	0xcf, 0x10, 0xd2, 0xd9, 0xeb, 0xd1, 0xe1, 0xa6, 0x54, 0xe6, 0x7a, 0x14, 0xb3, 0x04, 0x07, 0x5c,
	0xa2, 0x44, 0x8f, 0x8f, 0x47, 0x57, 0xee, 0xe7, 0xc9, 0x36, 0xaa, 0xe0, 0xbb, 0xbf, 0xe1, 0x12,
	0x51, 0x66, 0xc2, 0x75, 0xf2, 0x21, 0x87, 0x5c, 0x71, 0x18, 0x0e, 0xd1, 0x80, 0x51, 0x38, 0xd4,
	0x9e, 0xb8, 0xfd, 0xdc, 0xa0, 0xcd, 0x63, 0x37, 0x66, 0x70, 0x44, 0x9b, 0xfb, 0x77, 0x4a, 0x1b,
	0x1d, 0x2c, 0xb3, 0xd9, 0x30, 0x65, 0x97, 0xf9, 0x2a, 0x3b, 0xdc, 0x0b, 0xeb, 0xd3, 0x9d, 0xa5,
	0xa7, 0xf7, 0xcf, 0x97, 0xc6, 0x42, 0x30, 0xcf, 0x85, 0xa3, 0xf3, 0x07, 0x95, 0x3e, 0x06, 0x97,
	0x74, 0xee, 0x40, 0x98, 0x4a, 0x3a, 0xac, 0x61, 0x9e, 0x8a, 0xdb, 0x91, 0xd0, 0x26, 0x5c, 0xad,
	0xa9, 0xf5, 0x45, 0x62, 0xb0, 0xe8, 0xac, 0x38, 0xe5, 0xc5, 0xa0, 0x55, 0x6d, 0xd2, 0x2a, 0x33,
	0x2b, 0x7e, 0x4e, 0x5b, 0x67, 0xa6, 0x10, 0x30, 0xf8, 0x5b, 0x8b, 0x76, 0x4d, 0xcd, 0xf3, 0xab,
	0x46, 0xff, 0x36, 0xc8, 0x16, 0xe9, 0xc9, 0xf1, 0x24, 0x22, 0x6f, 0x93, 0x88, 0x7c, 0x4c, 0x22,
	0xf2, 0x3a, 0x8d, 0xc8, 0x78, 0x1a, 0x11, 0x1a, 0x62, 0x21, 0x59, 0x99, 0x02, 0x68, 0xe6, 0x1c,
	0x84, 0x5c, 0x59, 0x19, 0x1b, 0xf7, 0x68, 0x1f, 0xb2, 0xd4, 0x0b, 0x9e, 0x90, 0x8b, 0xf5, 0x6f,
	0x26, 0x38, 0x82, 0xbf, 0x02, 0x6f, 0x42, 0x91, 0x27, 0xb3, 0xbb, 0x88, 0x9b, 0x6e, 0xf7, 0x3b,
	0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xd0, 0xfa, 0x7f, 0x69, 0x34, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ObjectClient is the client API for Object service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ObjectClient interface {
	Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_ID, error)
	GetObject(ctx context.Context, in *payload.Object_Request, opts ...grpc.CallOption) (*payload.Object_Vector, error)
	StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (Object_StreamGetObjectClient, error)
}

type objectClient struct {
	cc *grpc.ClientConn
}

func NewObjectClient(cc *grpc.ClientConn) ObjectClient {
	return &objectClient{cc}
}

func (c *objectClient) Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_ID, error) {
	out := new(payload.Object_ID)
	err := c.cc.Invoke(ctx, "/vald.v1.Object/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectClient) GetObject(ctx context.Context, in *payload.Object_Request, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/vald.v1.Object/GetObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectClient) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (Object_StreamGetObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Object_serviceDesc.Streams[0], "/vald.v1.Object/StreamGetObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &objectStreamGetObjectClient{stream}
	return x, nil
}

type Object_StreamGetObjectClient interface {
	Send(*payload.Object_Request) error
	Recv() (*payload.Object_StreamVector, error)
	grpc.ClientStream
}

type objectStreamGetObjectClient struct {
	grpc.ClientStream
}

func (x *objectStreamGetObjectClient) Send(m *payload.Object_Request) error {
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
type ObjectServer interface {
	Exists(context.Context, *payload.Object_ID) (*payload.Object_ID, error)
	GetObject(context.Context, *payload.Object_Request) (*payload.Object_Vector, error)
	StreamGetObject(Object_StreamGetObjectServer) error
}

// UnimplementedObjectServer can be embedded to have forward compatible implementations.
type UnimplementedObjectServer struct {
}

func (*UnimplementedObjectServer) Exists(ctx context.Context, req *payload.Object_ID) (*payload.Object_ID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (*UnimplementedObjectServer) GetObject(ctx context.Context, req *payload.Object_Request) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetObject not implemented")
}
func (*UnimplementedObjectServer) StreamGetObject(srv Object_StreamGetObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamGetObject not implemented")
}

func RegisterObjectServer(s *grpc.Server, srv ObjectServer) {
	s.RegisterService(&_Object_serviceDesc, srv)
}

func _Object_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
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
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ObjectServer).Exists(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Object_GetObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Request)
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
		return srv.(ObjectServer).GetObject(ctx, req.(*payload.Object_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Object_StreamGetObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ObjectServer).StreamGetObject(&objectStreamGetObjectServer{stream})
}

type Object_StreamGetObjectServer interface {
	Send(*payload.Object_StreamVector) error
	Recv() (*payload.Object_Request, error)
	grpc.ServerStream
}

type objectStreamGetObjectServer struct {
	grpc.ServerStream
}

func (x *objectStreamGetObjectServer) Send(m *payload.Object_StreamVector) error {
	return x.ServerStream.SendMsg(m)
}

func (x *objectStreamGetObjectServer) Recv() (*payload.Object_Request, error) {
	m := new(payload.Object_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Object_serviceDesc = grpc.ServiceDesc{
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
