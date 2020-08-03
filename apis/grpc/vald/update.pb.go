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

func init() { proto.RegisterFile("update.proto", fileDescriptor_3f0fa214029f1c21) }

var fileDescriptor_3f0fa214029f1c21 = []byte{
	// 275 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x2d, 0x48, 0x49,
	0x2c, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29, 0x4b, 0xcc, 0x49, 0x91, 0xe2,
	0x2d, 0x48, 0xac, 0xcc, 0xc9, 0x4f, 0x4c, 0x81, 0x08, 0x4a, 0xc9, 0xa4, 0xe7, 0xe7, 0xa7, 0xe7,
	0xa4, 0xea, 0x27, 0x16, 0x64, 0xea, 0x27, 0xe6, 0xe5, 0xe5, 0x97, 0x24, 0x96, 0x64, 0xe6, 0xe7,
	0x15, 0x43, 0x65, 0x79, 0x0a, 0x92, 0xf4, 0xd3, 0x0b, 0x73, 0x20, 0x3c, 0xa3, 0x6f, 0x8c, 0x5c,
	0x6c, 0xa1, 0x60, 0x13, 0x85, 0x82, 0xe1, 0x2c, 0x71, 0x3d, 0x98, 0x81, 0x10, 0x01, 0xbd, 0xa0,
	0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x29, 0x09, 0xb8, 0x84, 0x7f, 0x52, 0x56, 0x6a, 0x72, 0x89,
	0x9e, 0x4f, 0x7e, 0x32, 0xd8, 0x70, 0x25, 0xb1, 0x0d, 0x0f, 0xe4, 0x19, 0x9b, 0x2e, 0x3f, 0x99,
	0xcc, 0xc4, 0xa3, 0xc4, 0xae, 0x0f, 0x71, 0xa3, 0x15, 0xa3, 0x96, 0x90, 0x3b, 0x17, 0x4f, 0x70,
	0x49, 0x51, 0x6a, 0x62, 0x2e, 0xf9, 0x46, 0x33, 0x68, 0x30, 0x1a, 0x30, 0x0a, 0x79, 0x70, 0x71,
	0xfb, 0x96, 0xe6, 0x94, 0x64, 0x42, 0xcd, 0x91, 0x41, 0x37, 0x07, 0x2c, 0x09, 0x33, 0x4c, 0x12,
	0x97, 0x61, 0xc5, 0x4a, 0x0c, 0x52, 0x2c, 0x1b, 0x1e, 0xc8, 0x33, 0x39, 0xf9, 0x9e, 0x78, 0x24,
	0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x5c, 0x7c, 0xf9, 0x45, 0xe9, 0x7a,
	0x65, 0x29, 0x89, 0x89, 0xc5, 0x7a, 0xa0, 0x10, 0x75, 0x62, 0x0f, 0x4b, 0xcc, 0x49, 0x71, 0x2c,
	0xc8, 0x0c, 0x60, 0x8c, 0x52, 0x49, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5,
	0x07, 0xab, 0xd0, 0x07, 0xa9, 0x00, 0x05, 0x6e, 0xb1, 0x7e, 0x7a, 0x51, 0x41, 0x32, 0x98, 0x9b,
	0xc4, 0x06, 0x0e, 0x4e, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9c, 0xd5, 0x1e, 0x66, 0x9f,
	0x01, 0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/vald.Update/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *updateClient) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (Update_StreamUpdateClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Update_serviceDesc.Streams[0], "/vald.Update/StreamUpdate", opts...)
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
	err := c.cc.Invoke(ctx, "/vald.Update/MultiUpdate", in, out, opts...)
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
		FullMethod: "/vald.Update/Update",
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
		FullMethod: "/vald.Update/MultiUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UpdateServer).MultiUpdate(ctx, req.(*payload.Update_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Update_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vald.Update",
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
	Metadata: "update.proto",
}
