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
	// 297 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x91, 0xb1, 0x4a, 0xc4, 0x40,
	0x10, 0x86, 0xdd, 0x2b, 0xee, 0x60, 0xbd, 0x42, 0x02, 0x36, 0xb9, 0x23, 0x07, 0xb1, 0x50, 0x2c,
	0x76, 0x8d, 0x76, 0x96, 0x57, 0x7b, 0x28, 0xca, 0x59, 0xd8, 0xc8, 0x24, 0x59, 0xe2, 0xca, 0x26,
	0xb3, 0x26, 0x9b, 0x80, 0xad, 0xaf, 0x60, 0xe5, 0x1b, 0x59, 0x0a, 0xbe, 0x80, 0x04, 0x1f, 0x44,
	0xb2, 0x9b, 0x88, 0x82, 0xd8, 0x5c, 0x95, 0xe1, 0xff, 0x67, 0xbe, 0xfc, 0xec, 0x4f, 0x17, 0xa0,
	0x65, 0xc5, 0x75, 0x89, 0x06, 0x79, 0x13, 0xf1, 0x06, 0x54, 0xca, 0x6b, 0x9d, 0x82, 0x11, 0xcc,
	0x8a, 0xde, 0xa4, 0x93, 0x58, 0x13, 0xf9, 0x7b, 0xbf, 0x37, 0x35, 0x3c, 0x2a, 0x84, 0x74, 0xf8,
	0xba, 0x6d, 0x7f, 0x9e, 0x21, 0x66, 0x4a, 0x70, 0xd0, 0x92, 0x43, 0x51, 0xa0, 0x01, 0x23, 0xb1,
	0xa8, 0x9c, 0x7b, 0xfc, 0x32, 0xa2, 0xe3, 0xb5, 0x85, 0x7b, 0xeb, 0xef, 0xc9, 0x67, 0x03, 0xa2,
	0x89, 0x98, 0xd3, 0xd8, 0xa5, 0x78, 0xa8, 0x45, 0x65, 0xfc, 0xd9, 0x4f, 0xef, 0x3c, 0xbe, 0x17,
	0x89, 0x61, 0x67, 0x98, 0x58, 0x68, 0xe8, 0x3d, 0xbd, 0x7f, 0x3e, 0x8f, 0xa6, 0xe1, 0xa4, 0x0f,
	0x7c, 0x4a, 0x0e, 0xbd, 0x15, 0x9d, 0x5e, 0x99, 0x52, 0x40, 0xbe, 0x29, 0x7c, 0xeb, 0x80, 0x1c,
	0x11, 0x4f, 0xd2, 0xed, 0x55, 0xad, 0x8c, 0xec, 0x69, 0x8b, 0x3f, 0x68, 0xd6, 0x1f, 0x90, 0xf3,
	0x7f, 0x90, 0x55, 0x38, 0xb3, 0x81, 0x77, 0xc3, 0x9d, 0x3e, 0x30, 0xcf, 0xbb, 0x5b, 0xad, 0xba,
	0xe4, 0xcb, 0xdb, 0xd7, 0x36, 0x20, 0x6f, 0x6d, 0x40, 0x3e, 0xda, 0x80, 0x50, 0x1f, 0xcb, 0x8c,
	0x35, 0x29, 0x40, 0xc5, 0xec, 0xfb, 0x83, 0x96, 0x1d, 0xb2, 0x9b, 0x97, 0xf4, 0x1a, 0x54, 0xea,
	0xfe, 0x7e, 0x41, 0x6e, 0xf6, 0x33, 0x69, 0xee, 0xea, 0x98, 0x25, 0x98, 0x73, 0x7b, 0xe0, 0x3a,
	0xb4, 0x65, 0x65, 0xa5, 0x4e, 0x86, 0x56, 0xe3, 0xb1, 0xed, 0xe0, 0xe4, 0x2b, 0x00, 0x00, 0xff,
	0xff, 0x08, 0x6a, 0x73, 0x66, 0xf2, 0x01, 0x00, 0x00,
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
	Recv() (*payload.Object_Location, error)
	grpc.ClientStream
}

type updateStreamUpdateClient struct {
	grpc.ClientStream
}

func (x *updateStreamUpdateClient) Send(m *payload.Update_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *updateStreamUpdateClient) Recv() (*payload.Object_Location, error) {
	m := new(payload.Object_Location)
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
	Send(*payload.Object_Location) error
	Recv() (*payload.Update_Request, error)
	grpc.ServerStream
}

type updateStreamUpdateServer struct {
	grpc.ServerStream
}

func (x *updateStreamUpdateServer) Send(m *payload.Object_Location) error {
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
