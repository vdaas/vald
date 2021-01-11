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

func init() { proto.RegisterFile("apis/proto/v1/vald/object.proto", fileDescriptor_f3068a4c11e32302) }

var fileDescriptor_f3068a4c11e32302 = []byte{
	// 276 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0x2f, 0x4b, 0xcc, 0x49, 0xd1, 0xcf,
	0x4f, 0xca, 0x4a, 0x4d, 0x2e, 0xd1, 0x03, 0x0b, 0x0a, 0xb1, 0x83, 0x84, 0xf4, 0xca, 0x0c, 0xa5,
	0x94, 0x51, 0x55, 0x16, 0x24, 0x56, 0xe6, 0xe4, 0x27, 0xa6, 0xc0, 0x68, 0x88, 0x6a, 0x29, 0x99,
	0xf4, 0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0xfd, 0xc4, 0x82, 0x4c, 0xfd, 0xc4, 0xbc, 0xbc, 0xfc, 0x92,
	0xc4, 0x92, 0xcc, 0xfc, 0xbc, 0x62, 0x88, 0xac, 0xd1, 0x37, 0x46, 0x2e, 0x36, 0x7f, 0xb0, 0xe1,
	0x42, 0x3e, 0x5c, 0x6c, 0xae, 0x15, 0x99, 0xc5, 0x25, 0xc5, 0x42, 0xa2, 0x7a, 0x30, 0x23, 0xca,
	0x0c, 0xf5, 0x20, 0xb2, 0x7a, 0x9e, 0x2e, 0x52, 0xd8, 0x85, 0x95, 0x44, 0x9a, 0x2e, 0x3f, 0x99,
	0xcc, 0xc4, 0x27, 0xc4, 0xa3, 0x9f, 0x0a, 0xd6, 0xae, 0x5f, 0x9d, 0x99, 0x52, 0x2b, 0x14, 0xcc,
	0xc5, 0xe9, 0x9e, 0x5a, 0x02, 0x35, 0x1a, 0x87, 0x81, 0x92, 0x58, 0x84, 0xc3, 0x52, 0x93, 0x4b,
	0xf2, 0x8b, 0x90, 0x0c, 0x85, 0x78, 0x1d, 0x62, 0xa8, 0x27, 0x17, 0x7f, 0x70, 0x49, 0x51, 0x6a,
	0x62, 0x2e, 0x25, 0x46, 0x33, 0x68, 0x30, 0x1a, 0x30, 0x3a, 0xc5, 0x9f, 0x78, 0x24, 0xc7, 0x78,
	0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x5c, 0x52, 0xf9, 0x45, 0xe9, 0x7a, 0x65, 0x29,
	0x89, 0x89, 0xc5, 0x7a, 0xe0, 0xc0, 0x4d, 0x2c, 0xc8, 0x04, 0x69, 0x04, 0xb1, 0x9d, 0xb8, 0xc2,
	0x12, 0x73, 0x52, 0x20, 0x46, 0x04, 0x30, 0x46, 0xa9, 0xa7, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9,
	0x25, 0xe7, 0xe7, 0xea, 0x83, 0x35, 0x40, 0x22, 0x08, 0x1c, 0x13, 0xe9, 0x45, 0x05, 0xc9, 0xb0,
	0x28, 0x4b, 0x62, 0x03, 0x07, 0xb0, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xe5, 0x83, 0x50, 0xc6,
	0xcf, 0x01, 0x00, 0x00,
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
	GetObject(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_Vector, error)
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

func (c *objectClient) GetObject(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
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
	Send(*payload.Object_ID) error
	Recv() (*payload.Object_Vector, error)
	grpc.ClientStream
}

type objectStreamGetObjectClient struct {
	grpc.ClientStream
}

func (x *objectStreamGetObjectClient) Send(m *payload.Object_ID) error {
	return x.ClientStream.SendMsg(m)
}

func (x *objectStreamGetObjectClient) Recv() (*payload.Object_Vector, error) {
	m := new(payload.Object_Vector)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ObjectServer is the server API for Object service.
type ObjectServer interface {
	Exists(context.Context, *payload.Object_ID) (*payload.Object_ID, error)
	GetObject(context.Context, *payload.Object_ID) (*payload.Object_Vector, error)
	StreamGetObject(Object_StreamGetObjectServer) error
}

// UnimplementedObjectServer can be embedded to have forward compatible implementations.
type UnimplementedObjectServer struct {
}

func (*UnimplementedObjectServer) Exists(ctx context.Context, req *payload.Object_ID) (*payload.Object_ID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (*UnimplementedObjectServer) GetObject(ctx context.Context, req *payload.Object_ID) (*payload.Object_Vector, error) {
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
	in := new(payload.Object_ID)
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
		return srv.(ObjectServer).GetObject(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Object_StreamGetObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ObjectServer).StreamGetObject(&objectStreamGetObjectServer{stream})
}

type Object_StreamGetObjectServer interface {
	Send(*payload.Object_Vector) error
	Recv() (*payload.Object_ID, error)
	grpc.ServerStream
}

type objectStreamGetObjectServer struct {
	grpc.ServerStream
}

func (x *objectStreamGetObjectServer) Send(m *payload.Object_Vector) error {
	return x.ServerStream.SendMsg(m)
}

func (x *objectStreamGetObjectServer) Recv() (*payload.Object_ID, error) {
	m := new(payload.Object_ID)
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
