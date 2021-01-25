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
	// 332 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0x3f, 0x4b, 0x03, 0x31,
	0x14, 0x37, 0x1d, 0x2a, 0x06, 0xb1, 0x78, 0x28, 0xe2, 0x0d, 0x57, 0xa9, 0x83, 0x2e, 0x26, 0x56,
	0x07, 0xf7, 0xa2, 0x88, 0x20, 0x28, 0x0a, 0x1d, 0xba, 0xbd, 0xbb, 0x8b, 0x31, 0x72, 0xed, 0x3b,
	0x2f, 0xe9, 0xa1, 0x88, 0x8b, 0x5f, 0x41, 0xfc, 0x3e, 0x8e, 0x1d, 0x05, 0xbf, 0x80, 0xb4, 0x7e,
	0x10, 0x49, 0x72, 0x05, 0x95, 0x3a, 0x38, 0xe5, 0xe5, 0xe5, 0xf7, 0xe7, 0x91, 0xdf, 0xa3, 0x4d,
	0xc8, 0x95, 0xe6, 0x79, 0x81, 0x06, 0x79, 0xd9, 0xe6, 0x25, 0x64, 0x29, 0xc7, 0xf8, 0x46, 0x24,
	0x86, 0xb9, 0x66, 0x30, 0x6f, 0x5b, 0xac, 0x6c, 0x87, 0x9b, 0x3f, 0x91, 0x39, 0xdc, 0x67, 0x08,
	0xe9, 0xf4, 0xf4, 0xe8, 0x70, 0x47, 0x2a, 0x73, 0x3d, 0x8c, 0x59, 0x82, 0x7d, 0x2e, 0x51, 0xa2,
	0xc7, 0xc7, 0xc3, 0x2b, 0x77, 0xf3, 0x64, 0x5b, 0x55, 0xf0, 0x83, 0xdf, 0x70, 0x89, 0x28, 0x33,
	0xe1, 0x9c, 0x7c, 0xc9, 0x21, 0x57, 0x1c, 0x06, 0x03, 0x34, 0x60, 0x14, 0x0e, 0xb4, 0x27, 0xee,
	0xbd, 0xd4, 0x68, 0xfd, 0xcc, 0x8d, 0x19, 0x9c, 0xd2, 0xfa, 0xd1, 0x9d, 0xd2, 0x46, 0x07, 0xab,
	0x6c, 0x3a, 0x4c, 0xd9, 0x66, 0xfe, 0x95, 0x9d, 0x1c, 0x86, 0xb3, 0xdb, 0xad, 0x95, 0xa7, 0xf7,
	0xcf, 0xe7, 0xda, 0x52, 0xb0, 0xc8, 0x85, 0xa3, 0xf3, 0x07, 0x95, 0x3e, 0x06, 0x40, 0x17, 0x8e,
	0x85, 0xa9, 0xa4, 0x37, 0x66, 0x30, 0xbb, 0x22, 0x31, 0x58, 0x5c, 0x88, 0xdb, 0xa1, 0xd0, 0x26,
	0x5c, 0xff, 0x13, 0xd1, 0x5a, 0x73, 0xfa, 0xcb, 0x41, 0xa3, 0xfa, 0x4f, 0xab, 0xcf, 0xac, 0x45,
	0x8f, 0x36, 0x2e, 0x4d, 0x21, 0xa0, 0xff, 0x1f, 0xa3, 0xe6, 0x0c, 0x84, 0x57, 0xa9, 0xec, 0xe6,
	0xb6, 0xc9, 0x2e, 0xe9, 0xc8, 0xd1, 0x38, 0x22, 0x6f, 0xe3, 0x88, 0x7c, 0x8c, 0x23, 0xf2, 0x3a,
	0x89, 0xc8, 0x68, 0x12, 0x11, 0x1a, 0x62, 0x21, 0x59, 0x99, 0x02, 0x68, 0xe6, 0xd2, 0x84, 0x5c,
	0x59, 0x19, 0x5b, 0x77, 0x68, 0x17, 0xb2, 0xd4, 0x0b, 0x9e, 0x93, 0xde, 0xd6, 0xb7, 0x40, 0x1c,
	0xc1, 0x6f, 0x84, 0x0f, 0xa4, 0xc8, 0x93, 0xe9, 0x8e, 0xc4, 0x75, 0x97, 0xc3, 0xfe, 0x57, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x67, 0xb5, 0x39, 0xb1, 0x40, 0x02, 0x00, 0x00,
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
	GetObject(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (*payload.Object_Vector, error)
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

func (c *objectClient) GetObject(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
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
type ObjectServer interface {
	Exists(context.Context, *payload.Object_ID) (*payload.Object_ID, error)
	GetObject(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error)
	StreamGetObject(Object_StreamGetObjectServer) error
}

// UnimplementedObjectServer can be embedded to have forward compatible implementations.
type UnimplementedObjectServer struct {
}

func (*UnimplementedObjectServer) Exists(ctx context.Context, req *payload.Object_ID) (*payload.Object_ID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (*UnimplementedObjectServer) GetObject(ctx context.Context, req *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
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
	m := new(payload.Object_VectorRequest)
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
