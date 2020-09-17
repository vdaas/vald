//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	proto "github.com/gogo/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
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
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0x2f, 0x4b, 0xcc, 0x49, 0xd1, 0xcf,
	0x4f, 0xca, 0x4a, 0x4d, 0x2e, 0xd1, 0x03, 0x0b, 0x0a, 0xb1, 0x80, 0x84, 0xa4, 0x94, 0x51, 0x95,
	0x15, 0x24, 0x56, 0xe6, 0xe4, 0x27, 0xa6, 0xc0, 0x68, 0x88, 0x52, 0x29, 0x99, 0xf4, 0xfc, 0xfc,
	0xf4, 0x9c, 0x54, 0xfd, 0xc4, 0x82, 0x4c, 0xfd, 0xc4, 0xbc, 0xbc, 0xfc, 0x92, 0xc4, 0x92, 0xcc,
	0xfc, 0xbc, 0x62, 0x88, 0xac, 0xd1, 0x13, 0x46, 0x2e, 0x36, 0x7f, 0xb0, 0xc9, 0x42, 0x6e, 0x5c,
	0x6c, 0xae, 0x15, 0x99, 0xc5, 0x25, 0xc5, 0x42, 0x42, 0x7a, 0x30, 0x23, 0x20, 0x52, 0x7a, 0x9e,
	0x2e, 0x52, 0x58, 0xc4, 0x94, 0x44, 0x9a, 0x2e, 0x3f, 0x99, 0xcc, 0xc4, 0x27, 0xc4, 0xa3, 0x9f,
	0x0a, 0xd6, 0xa8, 0x5f, 0x9d, 0x99, 0x52, 0x2b, 0xe4, 0xcb, 0xc5, 0xe9, 0x9e, 0x5a, 0x02, 0x35,
	0x14, 0x9b, 0x51, 0x62, 0xe8, 0x62, 0x61, 0xa9, 0xc9, 0x25, 0xf9, 0x45, 0x48, 0xc6, 0x41, 0xfc,
	0x0a, 0x31, 0xce, 0x99, 0x8b, 0x3f, 0xb8, 0xa4, 0x28, 0x35, 0x31, 0x97, 0x3c, 0x43, 0x19, 0x34,
	0x18, 0x0d, 0x18, 0x9d, 0xa2, 0x4f, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23,
	0x39, 0x46, 0x2e, 0xa9, 0xfc, 0xa2, 0x74, 0xbd, 0xb2, 0x94, 0xc4, 0xc4, 0x62, 0x3d, 0x50, 0x38,
	0xea, 0x25, 0x16, 0x64, 0xea, 0x95, 0x19, 0x82, 0xd9, 0x4e, 0xd0, 0xd0, 0x08, 0x60, 0x8c, 0x52,
	0x4f, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x07, 0x2b, 0x86, 0xc4, 0x03,
	0x38, 0xcc, 0xd3, 0x8b, 0x0a, 0x92, 0x61, 0x31, 0x93, 0xc4, 0x06, 0x0e, 0x4a, 0x63, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x21, 0x7c, 0x8b, 0xfe, 0xb6, 0x01, 0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/vald.Object/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectClient) GetObject(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/vald.Object/GetObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectClient) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (Object_StreamGetObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Object_serviceDesc.Streams[0], "/vald.Object/StreamGetObject", opts...)
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
		FullMethod: "/vald.Object/Exists",
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
		FullMethod: "/vald.Object/GetObject",
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
	ServiceName: "vald.Object",
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
