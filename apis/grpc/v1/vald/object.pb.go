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
	// 317 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xb1, 0x4a, 0x03, 0x41,
	0x10, 0x75, 0x2d, 0x22, 0x2e, 0xa2, 0x70, 0x68, 0xe1, 0x15, 0x17, 0xd0, 0x42, 0x1b, 0x77, 0x8d,
	0x16, 0xf6, 0x41, 0x11, 0x41, 0x88, 0x10, 0x48, 0x61, 0x37, 0x77, 0xb7, 0xae, 0x2b, 0x97, 0xcc,
	0x72, 0xb7, 0x39, 0x14, 0xb1, 0xf1, 0x17, 0xfc, 0x21, 0xcb, 0x94, 0x82, 0x3f, 0x20, 0x89, 0x9f,
	0x61, 0x21, 0xb7, 0x93, 0x83, 0x28, 0x27, 0x56, 0x3b, 0xfb, 0xf6, 0xbd, 0x37, 0xc3, 0x9b, 0xe5,
	0x6d, 0xb0, 0xa6, 0x90, 0x36, 0x47, 0x87, 0xb2, 0xec, 0xc8, 0x12, 0xb2, 0x54, 0x62, 0x7c, 0xa7,
	0x12, 0x27, 0x3c, 0x18, 0xac, 0x54, 0x90, 0x28, 0x3b, 0xe1, 0xee, 0x4f, 0xa6, 0x85, 0x87, 0x0c,
	0x21, 0xad, 0x4f, 0x62, 0x87, 0x07, 0xda, 0xb8, 0xdb, 0x71, 0x2c, 0x12, 0x1c, 0x4a, 0x8d, 0x1a,
	0x89, 0x1f, 0x8f, 0x6f, 0xfc, 0x8d, 0xc4, 0x55, 0x35, 0xa7, 0x9f, 0xfc, 0xa6, 0x6b, 0x44, 0x9d,
	0x29, 0xdf, 0x89, 0x4a, 0x09, 0xd6, 0x48, 0x18, 0x8d, 0xd0, 0x81, 0x33, 0x38, 0x2a, 0x48, 0x78,
	0xf4, 0xc5, 0x78, 0xab, 0xe7, 0xc7, 0x0c, 0x2e, 0x79, 0xeb, 0xec, 0xde, 0x14, 0xae, 0x08, 0xb6,
	0x44, 0x3d, 0x4c, 0xd9, 0x11, 0xf4, 0x2a, 0x2e, 0x4e, 0xc3, 0x66, 0x78, 0x67, 0xf3, 0xf9, 0xfd,
	0xf3, 0x65, 0x79, 0x3d, 0x58, 0x93, 0xca, 0xcb, 0xe5, 0xa3, 0x49, 0x9f, 0x82, 0x3e, 0x5f, 0x3d,
	0x57, 0x6e, 0x6e, 0xfd, 0x87, 0xe1, 0x76, 0x03, 0x3c, 0x50, 0x89, 0xc3, 0x7c, 0xc1, 0x94, 0x42,
	0x24, 0xd3, 0x1e, 0xdf, 0xe8, 0xbb, 0x5c, 0xc1, 0xf0, 0x5f, 0xeb, 0x76, 0x03, 0x4c, 0xd2, 0x79,
	0x83, 0xa5, 0x7d, 0x76, 0xc8, 0xba, 0x7a, 0x32, 0x8d, 0xd8, 0xdb, 0x34, 0x62, 0x1f, 0xd3, 0x88,
	0xbd, 0xce, 0x22, 0x36, 0x99, 0x45, 0x8c, 0x87, 0x98, 0x6b, 0x51, 0xa6, 0x00, 0x85, 0xf0, 0x4b,
	0x03, 0x6b, 0x2a, 0x9b, 0xaa, 0xee, 0xf2, 0x01, 0x64, 0x29, 0x19, 0x5e, 0xb1, 0xeb, 0xbd, 0x85,
	0xdc, 0xbd, 0x80, 0x16, 0x4f, 0xb9, 0xe7, 0x36, 0xa9, 0xbf, 0x42, 0xdc, 0xf2, 0x71, 0x1f, 0x7f,
	0x07, 0x00, 0x00, 0xff, 0xff, 0x5d, 0x8d, 0x64, 0xa7, 0x27, 0x02, 0x00, 0x00,
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
	Recv() (*payload.Object_StreamVector, error)
	grpc.ClientStream
}

type objectStreamGetObjectClient struct {
	grpc.ClientStream
}

func (x *objectStreamGetObjectClient) Send(m *payload.Object_ID) error {
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
	Send(*payload.Object_StreamVector) error
	Recv() (*payload.Object_ID, error)
	grpc.ServerStream
}

type objectStreamGetObjectServer struct {
	grpc.ServerStream
}

func (x *objectStreamGetObjectServer) Send(m *payload.Object_StreamVector) error {
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
