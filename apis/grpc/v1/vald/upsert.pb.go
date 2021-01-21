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

func init() { proto.RegisterFile("apis/proto/v1/vald/upsert.proto", fileDescriptor_792e000853e2404f) }
func init() {
	golang_proto.RegisterFile("apis/proto/v1/vald/upsert.proto", fileDescriptor_792e000853e2404f)
}

var fileDescriptor_792e000853e2404f = []byte{
	// 339 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xc1, 0x4a, 0xf3, 0x40,
	0x18, 0xfc, 0xb7, 0x87, 0x16, 0xf2, 0xf7, 0x20, 0x01, 0x2f, 0xa9, 0xa4, 0x18, 0x0f, 0x4a, 0xc1,
	0x5d, 0xab, 0x07, 0xa1, 0xc7, 0x9e, 0x15, 0xc5, 0x52, 0x0f, 0xde, 0xbe, 0x24, 0xeb, 0xba, 0x92,
	0xf6, 0x5b, 0x93, 0x4d, 0xc0, 0xab, 0xaf, 0xe0, 0x03, 0xf8, 0x2a, 0x1e, 0x7b, 0x14, 0x7c, 0x01,
	0x49, 0x7d, 0x10, 0xc9, 0x6e, 0x52, 0xaa, 0x14, 0xf1, 0x94, 0xc9, 0xb7, 0x33, 0xb3, 0xb3, 0xcc,
	0xe7, 0xf4, 0x41, 0xc9, 0x8c, 0xa9, 0x14, 0x35, 0xb2, 0x62, 0xc8, 0x0a, 0x48, 0x62, 0x96, 0xab,
	0x8c, 0xa7, 0x9a, 0x9a, 0xa1, 0xdb, 0xa9, 0x46, 0xb4, 0x18, 0x7a, 0x7b, 0xdf, 0x99, 0x0a, 0x1e,
	0x13, 0x84, 0xb8, 0xf9, 0x5a, 0xb6, 0x77, 0x28, 0xa4, 0xbe, 0xcb, 0x43, 0x1a, 0xe1, 0x8c, 0x09,
	0x14, 0x68, 0xf9, 0x61, 0x7e, 0x6b, 0xfe, 0xac, 0xb8, 0x42, 0x35, 0xfd, 0xf4, 0x27, 0x5d, 0x20,
	0x8a, 0x84, 0x9b, 0x9b, 0x2c, 0x64, 0xa0, 0x24, 0x83, 0xf9, 0x1c, 0x35, 0x68, 0x89, 0xf3, 0xcc,
	0x0a, 0x8f, 0x5f, 0x5a, 0x4e, 0x7b, 0x6a, 0x62, 0xba, 0xd3, 0x15, 0xf2, 0x68, 0x13, 0xa6, 0x18,
	0x52, 0x3b, 0xa3, 0x57, 0xfc, 0x21, 0xe7, 0x99, 0xf6, 0x7a, 0xeb, 0x67, 0x17, 0xe1, 0x3d, 0x8f,
	0x34, 0x3d, 0xc3, 0xc8, 0x98, 0x06, 0xee, 0xd3, 0xfb, 0xe7, 0x73, 0xab, 0x1b, 0x74, 0xea, 0xa7,
	0x8f, 0xc8, 0xc0, 0x9d, 0x38, 0xdd, 0x89, 0x4e, 0x39, 0xcc, 0xfe, 0x60, 0xbe, 0xbb, 0xc1, 0xdc,
	0x8a, 0x57, 0x57, 0xfc, 0x3b, 0x20, 0x47, 0xc4, 0x95, 0xce, 0xff, 0xf3, 0x3c, 0xd1, 0xb2, 0xf6,
	0xec, 0x6f, 0xf0, 0x34, 0xe7, 0x8d, 0xf1, 0xce, 0x2f, 0xa9, 0xb3, 0xa0, 0x67, 0x62, 0x6f, 0x07,
	0x5b, 0x75, 0x6c, 0x36, 0xab, 0xb4, 0x2a, 0xe1, 0x23, 0x32, 0x18, 0x8b, 0x45, 0xe9, 0x93, 0xb7,
	0xd2, 0x27, 0x1f, 0xa5, 0x4f, 0x5e, 0x97, 0x3e, 0x59, 0x2c, 0x7d, 0xe2, 0x78, 0x98, 0x0a, 0x5a,
	0xc4, 0x00, 0x19, 0x35, 0xbd, 0x82, 0x92, 0x95, 0x75, 0x85, 0xc7, 0xce, 0x35, 0x24, 0xb1, 0x4d,
	0x71, 0x49, 0x6e, 0xf6, 0xd7, 0xaa, 0x31, 0x02, 0xbb, 0x1b, 0xb6, 0x9a, 0x54, 0x45, 0xcd, 0xb6,
	0x84, 0x6d, 0xd3, 0xc8, 0xc9, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x35, 0x6d, 0xdd, 0x8d, 0x4a,
	0x02, 0x00, 0x00,
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
