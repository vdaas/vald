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
	math "math"

	_ "github.com/danielvladco/go-proto-gql/pb"
	proto "github.com/gogo/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/payload"
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

func init() { proto.RegisterFile("upsert.proto", fileDescriptor_056e166e1581c01d) }

var fileDescriptor_056e166e1581c01d = []byte{
	// 278 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x2d, 0x28, 0x4e,
	0x2d, 0x2a, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29, 0x4b, 0xcc, 0x49, 0x91, 0xe2,
	0x2d, 0x48, 0xac, 0xcc, 0xc9, 0x4f, 0x4c, 0x81, 0x08, 0x4a, 0xc9, 0xa4, 0xe7, 0xe7, 0xa7, 0xe7,
	0xa4, 0xea, 0x27, 0x16, 0x64, 0xea, 0x27, 0xe6, 0xe5, 0xe5, 0x97, 0x24, 0x96, 0x64, 0xe6, 0xe7,
	0x15, 0x43, 0x65, 0x79, 0x0a, 0x92, 0xf4, 0xd3, 0x0b, 0x73, 0x20, 0x3c, 0xa3, 0x6f, 0x8c, 0x5c,
	0x6c, 0xa1, 0x60, 0x13, 0x85, 0x82, 0xe1, 0x2c, 0x71, 0x3d, 0x98, 0x81, 0x10, 0x01, 0xbd, 0xa0,
	0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x29, 0x09, 0xb8, 0x84, 0x7f, 0x52, 0x56, 0x6a, 0x72, 0x89,
	0x9e, 0x4f, 0x7e, 0x32, 0xd8, 0x70, 0x25, 0xb1, 0x0d, 0x0f, 0xe4, 0x19, 0x9b, 0x2e, 0x3f, 0x99,
	0xcc, 0xc4, 0xa3, 0xc4, 0xae, 0x0f, 0x71, 0xa3, 0x15, 0xa3, 0x96, 0x90, 0x3b, 0x17, 0x4f, 0x70,
	0x49, 0x51, 0x6a, 0x62, 0x2e, 0xf9, 0x46, 0x33, 0x68, 0x30, 0x1a, 0x30, 0x0a, 0x79, 0x70, 0x71,
	0xfb, 0x96, 0xe6, 0x94, 0x64, 0x42, 0xcd, 0x91, 0x41, 0x32, 0x27, 0x25, 0xb1, 0x24, 0x55, 0x0f,
	0x2c, 0x09, 0x33, 0x4c, 0x12, 0x97, 0x61, 0xc5, 0x4a, 0x0c, 0x52, 0x2c, 0x1b, 0x1e, 0xc8, 0x33,
	0x39, 0x05, 0x9c, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x5c,
	0x42, 0xf9, 0x45, 0xe9, 0x7a, 0x65, 0x29, 0x89, 0x89, 0xc5, 0x7a, 0xa0, 0x10, 0xd5, 0x4b, 0x2c,
	0xc8, 0x74, 0x82, 0x86, 0x46, 0x00, 0x63, 0x94, 0x4a, 0x7a, 0x66, 0x49, 0x46, 0x69, 0x92, 0x5e,
	0x72, 0x7e, 0xae, 0x3e, 0x58, 0x91, 0x3e, 0x48, 0x11, 0x28, 0x7c, 0x8b, 0xf5, 0xd3, 0x8b, 0x0a,
	0x92, 0xc1, 0xdc, 0x24, 0x36, 0x70, 0x88, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x73, 0xec,
	0xaa, 0xf4, 0xa2, 0x01, 0x00, 0x00,
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
	MultiUpsert(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type upsertClient struct {
	cc *grpc.ClientConn
}

func NewUpsertClient(cc *grpc.ClientConn) UpsertClient {
	return &upsertClient{cc}
}

func (c *upsertClient) Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/vald.Upsert/Upsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *upsertClient) StreamUpsert(ctx context.Context, opts ...grpc.CallOption) (Upsert_StreamUpsertClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Upsert_serviceDesc.Streams[0], "/vald.Upsert/StreamUpsert", opts...)
	if err != nil {
		return nil, err
	}
	x := &upsertStreamUpsertClient{stream}
	return x, nil
}

type Upsert_StreamUpsertClient interface {
	Send(*payload.Upsert_Request) error
	Recv() (*payload.Object_Location, error)
	grpc.ClientStream
}

type upsertStreamUpsertClient struct {
	grpc.ClientStream
}

func (x *upsertStreamUpsertClient) Send(m *payload.Upsert_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *upsertStreamUpsertClient) Recv() (*payload.Object_Location, error) {
	m := new(payload.Object_Location)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *upsertClient) MultiUpsert(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.Upsert/MultiUpsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpsertServer is the server API for Upsert service.
type UpsertServer interface {
	Upsert(context.Context, *payload.Upsert_Request) (*payload.Object_Location, error)
	StreamUpsert(Upsert_StreamUpsertServer) error
	MultiUpsert(context.Context, *payload.Update_MultiRequest) (*payload.Object_Locations, error)
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
func (*UnimplementedUpsertServer) MultiUpsert(ctx context.Context, req *payload.Update_MultiRequest) (*payload.Object_Locations, error) {
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
		FullMethod: "/vald.Upsert/Upsert",
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
	Send(*payload.Object_Location) error
	Recv() (*payload.Upsert_Request, error)
	grpc.ServerStream
}

type upsertStreamUpsertServer struct {
	grpc.ServerStream
}

func (x *upsertStreamUpsertServer) Send(m *payload.Object_Location) error {
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
	in := new(payload.Update_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpsertServer).MultiUpsert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Upsert/MultiUpsert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UpsertServer).MultiUpsert(ctx, req.(*payload.Update_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Upsert_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vald.Upsert",
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
	Metadata: "upsert.proto",
}
