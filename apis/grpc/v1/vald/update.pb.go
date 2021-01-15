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

func init() { proto.RegisterFile("apis/proto/v1/vald/update.proto", fileDescriptor_a564bbf4b2600403) }

var fileDescriptor_a564bbf4b2600403 = []byte{
	// 305 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x31, 0x4a, 0xc4, 0x40,
	0x14, 0x86, 0x9d, 0x2d, 0x76, 0x61, 0xdc, 0x42, 0x02, 0x36, 0xd9, 0x65, 0x17, 0x63, 0xa1, 0x58,
	0xcc, 0x18, 0xed, 0x2c, 0xb7, 0x56, 0x14, 0x97, 0xb5, 0xb0, 0x91, 0x97, 0x64, 0x88, 0x23, 0x93,
	0xbc, 0x31, 0x99, 0x04, 0x6c, 0xbd, 0x82, 0x07, 0xf0, 0x3a, 0x96, 0x82, 0x17, 0x90, 0xe0, 0x41,
	0x24, 0x33, 0xc9, 0xa2, 0xb0, 0x88, 0xd5, 0x0c, 0xef, 0x7f, 0xff, 0xf7, 0x7e, 0xf8, 0xe9, 0x1c,
	0xb4, 0x2c, 0xb9, 0x2e, 0xd0, 0x20, 0xaf, 0x43, 0x5e, 0x83, 0x4a, 0x78, 0xa5, 0x13, 0x30, 0x82,
	0xd9, 0xa1, 0x37, 0x6a, 0x47, 0xac, 0x0e, 0xfd, 0xfd, 0xdf, 0x9b, 0x1a, 0x9e, 0x14, 0x42, 0xd2,
	0xbf, 0x6e, 0xdb, 0x9f, 0xa6, 0x88, 0xa9, 0x12, 0x1c, 0xb4, 0xe4, 0x90, 0xe7, 0x68, 0xc0, 0x48,
	0xcc, 0x4b, 0xa7, 0x9e, 0xbc, 0x0e, 0xe8, 0x70, 0x65, 0xe1, 0xde, 0x6a, 0xfd, 0xf3, 0x59, 0x8f,
	0xa8, 0x43, 0xe6, 0x66, 0xec, 0x5a, 0x3c, 0x56, 0xa2, 0x34, 0xfe, 0xe4, 0xa7, 0x76, 0x19, 0x3d,
	0x88, 0xd8, 0xb0, 0x73, 0x8c, 0x2d, 0x34, 0xf0, 0x9e, 0x3f, 0xbe, 0x5e, 0x06, 0xe3, 0x60, 0xd4,
	0x05, 0x3e, 0x23, 0x47, 0xde, 0x92, 0x8e, 0x97, 0xa6, 0x10, 0x90, 0xfd, 0x03, 0xbe, 0xb7, 0x01,
	0xee, 0xcc, 0xeb, 0x13, 0x5b, 0x87, 0xe4, 0x98, 0x78, 0x92, 0x6e, 0x5f, 0x54, 0xca, 0xc8, 0x8e,
	0x39, 0xdf, 0xc0, 0xb4, 0x7a, 0x0f, 0x9e, 0xfe, 0x91, 0xba, 0x0c, 0x26, 0x36, 0xf6, 0x6e, 0xb0,
	0xd3, 0xc5, 0xe6, 0x59, 0xeb, 0xd5, 0xaa, 0xcd, 0xbf, 0xb8, 0x7b, 0x6b, 0x66, 0xe4, 0xbd, 0x99,
	0x91, 0xcf, 0x66, 0x46, 0xa8, 0x8f, 0x45, 0xca, 0xea, 0x04, 0xa0, 0x64, 0xb6, 0x05, 0xd0, 0xb2,
	0x45, 0xb6, 0xff, 0x05, 0xbd, 0x01, 0x95, 0xb8, 0xeb, 0x57, 0xe4, 0xf6, 0x20, 0x95, 0xe6, 0xbe,
	0x8a, 0x58, 0x8c, 0x19, 0xb7, 0x06, 0xd7, 0xa4, 0xad, 0x2c, 0x2d, 0x74, 0xdc, 0x77, 0x1b, 0x0d,
	0x6d, 0x13, 0xa7, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa7, 0x5d, 0x12, 0xbc, 0xf8, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UpdateClient is the client API for Update service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UpdateClient interface {
	Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (Update_StreamUpdateClient, error)
	MultiUpdate(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type updateClient struct {
	cc *grpc.ClientConn
}

func NewUpdateClient(cc *grpc.ClientConn) UpdateClient {
	return &updateClient{cc}
}

func (c *updateClient) Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/vald.v1.Update/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *updateClient) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (Update_StreamUpdateClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Update_serviceDesc.Streams[0], "/vald.v1.Update/StreamUpdate", opts...)
	if err != nil {
		return nil, err
	}
	x := &updateStreamUpdateClient{stream}
	return x, nil
}

type Update_StreamUpdateClient interface {
	Send(*payload.Update_Request) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type updateStreamUpdateClient struct {
	grpc.ClientStream
}

func (x *updateStreamUpdateClient) Send(m *payload.Update_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *updateStreamUpdateClient) Recv() (*payload.Object_StreamLocation, error) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *updateClient) MultiUpdate(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.v1.Update/MultiUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateServer is the server API for Update service.
type UpdateServer interface {
	Update(context.Context, *payload.Update_Request) (*payload.Object_Location, error)
	StreamUpdate(Update_StreamUpdateServer) error
	MultiUpdate(context.Context, *payload.Update_MultiRequest) (*payload.Object_Locations, error)
}

// UnimplementedUpdateServer can be embedded to have forward compatible implementations.
type UnimplementedUpdateServer struct {
}

func (*UnimplementedUpdateServer) Update(ctx context.Context, req *payload.Update_Request) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedUpdateServer) StreamUpdate(srv Update_StreamUpdateServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpdate not implemented")
}
func (*UnimplementedUpdateServer) MultiUpdate(ctx context.Context, req *payload.Update_MultiRequest) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpdate not implemented")
}

func RegisterUpdateServer(s *grpc.Server, srv UpdateServer) {
	s.RegisterService(&_Update_serviceDesc, srv)
}

func _Update_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Update_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpdateServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Update/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UpdateServer).Update(ctx, req.(*payload.Update_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Update_StreamUpdate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UpdateServer).StreamUpdate(&updateStreamUpdateServer{stream})
}

type Update_StreamUpdateServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Update_Request, error)
	grpc.ServerStream
}

type updateStreamUpdateServer struct {
	grpc.ServerStream
}

func (x *updateStreamUpdateServer) Send(m *payload.Object_StreamLocation) error {
	return x.ServerStream.SendMsg(m)
}

func (x *updateStreamUpdateServer) Recv() (*payload.Update_Request, error) {
	m := new(payload.Update_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Update_MultiUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Update_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpdateServer).MultiUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Update/MultiUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UpdateServer).MultiUpdate(ctx, req.(*payload.Update_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Update_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vald.v1.Update",
	HandlerType: (*UpdateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Update",
			Handler:    _Update_Update_Handler,
		},
		{
			MethodName: "MultiUpdate",
			Handler:    _Update_MultiUpdate_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamUpdate",
			Handler:       _Update_StreamUpdate_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "apis/proto/v1/vald/update.proto",
}
