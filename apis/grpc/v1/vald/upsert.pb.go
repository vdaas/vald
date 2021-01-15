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

func init() { proto.RegisterFile("apis/proto/v1/vald/upsert.proto", fileDescriptor_792e000853e2404f) }

var fileDescriptor_792e000853e2404f = []byte{
	// 306 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xc1, 0x4a, 0xc3, 0x30,
	0x18, 0xc7, 0xcd, 0x0e, 0x1b, 0xc4, 0x1d, 0xa4, 0xe0, 0xa5, 0x1b, 0x1d, 0xd6, 0x83, 0xb2, 0x43,
	0x62, 0xf5, 0xb6, 0xe3, 0xce, 0x8a, 0xe2, 0x98, 0x07, 0x2f, 0xf2, 0xb5, 0x0d, 0x35, 0x92, 0x36,
	0xb1, 0x49, 0x0b, 0x5e, 0x7d, 0x05, 0x1f, 0xc0, 0xd7, 0xf1, 0x28, 0xf8, 0x02, 0x52, 0x7c, 0x10,
	0x69, 0xd2, 0x0e, 0x85, 0x21, 0x9e, 0x12, 0xbe, 0xff, 0xf7, 0xff, 0x7d, 0x7f, 0xf8, 0xe3, 0x19,
	0x28, 0xae, 0xa9, 0x2a, 0xa5, 0x91, 0xb4, 0x8e, 0x68, 0x0d, 0x22, 0xa5, 0x95, 0xd2, 0xac, 0x34,
	0xc4, 0x0e, 0xbd, 0x51, 0x3b, 0x22, 0x75, 0xe4, 0x1f, 0xfe, 0xde, 0x54, 0xf0, 0x24, 0x24, 0xa4,
	0xfd, 0xeb, 0xb6, 0xfd, 0x69, 0x26, 0x65, 0x26, 0x18, 0x05, 0xc5, 0x29, 0x14, 0x85, 0x34, 0x60,
	0xb8, 0x2c, 0xb4, 0x53, 0x4f, 0x5f, 0x07, 0x78, 0xb8, 0xb6, 0x70, 0x6f, 0xbd, 0xf9, 0xf9, 0xa4,
	0x47, 0xd4, 0x11, 0x71, 0x33, 0x72, 0xcd, 0x1e, 0x2b, 0xa6, 0x8d, 0x3f, 0xf9, 0xa9, 0x5d, 0xc6,
	0x0f, 0x2c, 0x31, 0xe4, 0x5c, 0x26, 0x16, 0x1a, 0x7a, 0xcf, 0x1f, 0x5f, 0x2f, 0x83, 0x71, 0x38,
	0xea, 0x02, 0x2f, 0xd0, 0xdc, 0x5b, 0xe1, 0xf1, 0xca, 0x94, 0x0c, 0xf2, 0x7f, 0xc0, 0x0f, 0xb6,
	0xc0, 0x9d, 0x79, 0x73, 0x62, 0xe7, 0x18, 0x9d, 0x20, 0x8f, 0xe3, 0xdd, 0x8b, 0x4a, 0x18, 0xde,
	0x31, 0x67, 0x5b, 0x98, 0x56, 0xef, 0xc1, 0xd3, 0x3f, 0x52, 0xeb, 0x70, 0x62, 0x63, 0xef, 0x87,
	0x7b, 0x5d, 0x6c, 0x9a, 0xb7, 0x5e, 0x25, 0xd8, 0x02, 0xcd, 0x97, 0x77, 0x6f, 0x4d, 0x80, 0xde,
	0x9b, 0x00, 0x7d, 0x36, 0x01, 0xc2, 0xbe, 0x2c, 0x33, 0x52, 0xa7, 0x00, 0x9a, 0xd8, 0x16, 0x40,
	0xf1, 0x16, 0xd9, 0xfe, 0x97, 0xf8, 0x06, 0x44, 0xea, 0xae, 0x5f, 0xa1, 0xdb, 0xa3, 0x8c, 0x9b,
	0xfb, 0x2a, 0x26, 0x89, 0xcc, 0xa9, 0x35, 0xb8, 0x26, 0x6d, 0x65, 0x59, 0xa9, 0x92, 0xbe, 0xdb,
	0x78, 0x68, 0x9b, 0x38, 0xfb, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xbe, 0x92, 0xcb, 0x48, 0xf8, 0x01,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UpsertClient is the client API for Upsert service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UpsertClient interface {
	Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamUpsert(ctx context.Context, opts ...grpc.CallOption) (Upsert_StreamUpsertClient, error)
	MultiUpsert(ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type upsertClient struct {
	cc *grpc.ClientConn
}

func NewUpsertClient(cc *grpc.ClientConn) UpsertClient {
	return &upsertClient{cc}
}

func (c *upsertClient) Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/vald.v1.Upsert/Upsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *upsertClient) StreamUpsert(ctx context.Context, opts ...grpc.CallOption) (Upsert_StreamUpsertClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Upsert_serviceDesc.Streams[0], "/vald.v1.Upsert/StreamUpsert", opts...)
	if err != nil {
		return nil, err
	}
	x := &upsertStreamUpsertClient{stream}
	return x, nil
}

type Upsert_StreamUpsertClient interface {
	Send(*payload.Upsert_Request) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type upsertStreamUpsertClient struct {
	grpc.ClientStream
}

func (x *upsertStreamUpsertClient) Send(m *payload.Upsert_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *upsertStreamUpsertClient) Recv() (*payload.Object_StreamLocation, error) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *upsertClient) MultiUpsert(ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.v1.Upsert/MultiUpsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpsertServer is the server API for Upsert service.
type UpsertServer interface {
	Upsert(context.Context, *payload.Upsert_Request) (*payload.Object_Location, error)
	StreamUpsert(Upsert_StreamUpsertServer) error
	MultiUpsert(context.Context, *payload.Upsert_MultiRequest) (*payload.Object_Locations, error)
}

// UnimplementedUpsertServer can be embedded to have forward compatible implementations.
type UnimplementedUpsertServer struct {
}

func (*UnimplementedUpsertServer) Upsert(ctx context.Context, req *payload.Upsert_Request) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Upsert not implemented")
}
func (*UnimplementedUpsertServer) StreamUpsert(srv Upsert_StreamUpsertServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpsert not implemented")
}
func (*UnimplementedUpsertServer) MultiUpsert(ctx context.Context, req *payload.Upsert_MultiRequest) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpsert not implemented")
}

func RegisterUpsertServer(s *grpc.Server, srv UpsertServer) {
	s.RegisterService(&_Upsert_serviceDesc, srv)
}

func _Upsert_Upsert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Upsert_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpsertServer).Upsert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Upsert/Upsert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UpsertServer).Upsert(ctx, req.(*payload.Upsert_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Upsert_StreamUpsert_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UpsertServer).StreamUpsert(&upsertStreamUpsertServer{stream})
}

type Upsert_StreamUpsertServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Upsert_Request, error)
	grpc.ServerStream
}

type upsertStreamUpsertServer struct {
	grpc.ServerStream
}

func (x *upsertStreamUpsertServer) Send(m *payload.Object_StreamLocation) error {
	return x.ServerStream.SendMsg(m)
}

func (x *upsertStreamUpsertServer) Recv() (*payload.Upsert_Request, error) {
	m := new(payload.Upsert_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Upsert_MultiUpsert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Upsert_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpsertServer).MultiUpsert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Upsert/MultiUpsert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UpsertServer).MultiUpsert(ctx, req.(*payload.Upsert_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Upsert_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vald.v1.Upsert",
	HandlerType: (*UpsertServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Upsert",
			Handler:    _Upsert_Upsert_Handler,
		},
		{
			MethodName: "MultiUpsert",
			Handler:    _Upsert_MultiUpsert_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamUpsert",
			Handler:       _Upsert_StreamUpsert_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "apis/proto/v1/vald/upsert.proto",
}
