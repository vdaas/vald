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

func init() { proto.RegisterFile("apis/proto/v1/vald/update.proto", fileDescriptor_a564bbf4b2600403) }
func init() {
	golang_proto.RegisterFile("apis/proto/v1/vald/update.proto", fileDescriptor_a564bbf4b2600403)
}

var fileDescriptor_a564bbf4b2600403 = []byte{
	// 337 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xcd, 0x4a, 0xf3, 0x40,
	0x14, 0xfd, 0xa6, 0x8b, 0x16, 0xf2, 0x75, 0x21, 0x01, 0x37, 0xa9, 0xa4, 0x18, 0x17, 0x8a, 0xe0,
	0x8c, 0xd5, 0x85, 0xe0, 0xb2, 0x6b, 0x45, 0xb1, 0xd4, 0x85, 0xbb, 0x9b, 0x1f, 0xc7, 0x91, 0x69,
	0xef, 0x98, 0x4c, 0x02, 0x6e, 0x7d, 0x05, 0x1f, 0xc0, 0x57, 0x71, 0xd9, 0xa5, 0xe0, 0x0b, 0x48,
	0xeb, 0x83, 0xc8, 0xcc, 0x24, 0xa5, 0x4a, 0x11, 0x57, 0x39, 0xb9, 0x73, 0xce, 0x99, 0x33, 0x9c,
	0xeb, 0xf5, 0x41, 0x89, 0x82, 0xa9, 0x1c, 0x35, 0xb2, 0x6a, 0xc0, 0x2a, 0x90, 0x29, 0x2b, 0x55,
	0x0a, 0x3a, 0xa3, 0x76, 0xe8, 0x77, 0xcc, 0x88, 0x56, 0x83, 0x60, 0xe7, 0x3b, 0x53, 0xc1, 0xa3,
	0x44, 0x48, 0x9b, 0xaf, 0x63, 0x07, 0x07, 0x5c, 0xe8, 0xbb, 0x32, 0xa6, 0x09, 0x4e, 0x18, 0x47,
	0x8e, 0x8e, 0x1f, 0x97, 0xb7, 0xf6, 0xcf, 0x89, 0x0d, 0xaa, 0xe9, 0x27, 0x3f, 0xe9, 0x1c, 0x91,
	0xcb, 0xcc, 0xde, 0xe4, 0x20, 0x03, 0x25, 0x18, 0x4c, 0xa7, 0xa8, 0x41, 0x0b, 0x9c, 0x16, 0x4e,
	0x78, 0xf4, 0xd2, 0xf2, 0xda, 0x63, 0x1b, 0xd3, 0x1f, 0x2f, 0x51, 0x40, 0x9b, 0x30, 0xd5, 0x80,
	0xba, 0x19, 0xbd, 0xca, 0x1e, 0xca, 0xac, 0xd0, 0x41, 0x6f, 0xf5, 0xec, 0x22, 0xbe, 0xcf, 0x12,
	0x4d, 0xcf, 0x30, 0xb1, 0xa6, 0x91, 0xff, 0xf4, 0xfe, 0xf9, 0xdc, 0xea, 0x46, 0x9d, 0xfa, 0xe9,
	0xa7, 0x64, 0xdf, 0x1f, 0x79, 0xdd, 0x91, 0xce, 0x33, 0x98, 0xfc, 0xc1, 0x7c, 0x7b, 0x8d, 0xb9,
	0x13, 0x2f, 0xaf, 0xf8, 0xb7, 0x47, 0x0e, 0x89, 0x2f, 0xbc, 0xff, 0xe7, 0xa5, 0xd4, 0xa2, 0xf6,
	0xec, 0xaf, 0xf1, 0xb4, 0xe7, 0x8d, 0xf1, 0xd6, 0x2f, 0xa9, 0x8b, 0xa8, 0x67, 0x63, 0x6f, 0x46,
	0x1b, 0x75, 0x6c, 0x36, 0x31, 0x5a, 0x25, 0x4d, 0xfe, 0x21, 0x9f, 0xcd, 0x43, 0xf2, 0x36, 0x0f,
	0xc9, 0xc7, 0x3c, 0x24, 0xaf, 0x8b, 0x90, 0xcc, 0x16, 0x21, 0xf1, 0x02, 0xcc, 0x39, 0xad, 0x52,
	0x80, 0x82, 0xda, 0x5e, 0x41, 0x09, 0x63, 0x6d, 0xf0, 0xd0, 0xbb, 0x06, 0x99, 0xba, 0x14, 0x97,
	0xe4, 0x66, 0x77, 0xa5, 0x1a, 0x2b, 0x70, 0xbb, 0xe1, 0xaa, 0xc9, 0x55, 0xd2, 0x6c, 0x4b, 0xdc,
	0xb6, 0x8d, 0x1c, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0xca, 0x6a, 0x3e, 0x13, 0x4a, 0x02, 0x00,
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
