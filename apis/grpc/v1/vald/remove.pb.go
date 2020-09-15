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

func init() { proto.RegisterFile("apis/proto/v1/vald/remove.proto", fileDescriptor_5b638f34e0c25c81) }

var fileDescriptor_5b638f34e0c25c81 = []byte{
	// 257 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0x2f, 0x4b, 0xcc, 0x49, 0xd1, 0x2f,
	0x4a, 0xcd, 0xcd, 0x2f, 0x4b, 0xd5, 0x03, 0x0b, 0x0a, 0xb1, 0x80, 0x84, 0xa4, 0x2c, 0xd3, 0x33,
	0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0xcb, 0x52, 0x12, 0x13, 0x8b, 0x21, 0x2a,
	0x51, 0x35, 0x17, 0x24, 0x56, 0xe6, 0xe4, 0x27, 0xa6, 0xc0, 0x68, 0x88, 0x01, 0x52, 0x32, 0xe9,
	0xf9, 0xf9, 0xe9, 0x39, 0xa9, 0x20, 0xb5, 0xfa, 0x89, 0x79, 0x79, 0xf9, 0x25, 0x89, 0x25, 0x99,
	0xf9, 0x79, 0xc5, 0x10, 0x59, 0xa3, 0x5b, 0x8c, 0x5c, 0x6c, 0x41, 0x60, 0xfb, 0x84, 0x6c, 0xe1,
	0x2c, 0x71, 0x3d, 0x98, 0x11, 0x10, 0x01, 0xbd, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x29,
	0x09, 0xb8, 0x84, 0x7f, 0x52, 0x56, 0x6a, 0x72, 0x89, 0x9e, 0x4f, 0x7e, 0x32, 0xd8, 0x38, 0x25,
	0x06, 0x21, 0x77, 0x2e, 0x9e, 0xe0, 0x92, 0xa2, 0xd4, 0xc4, 0x5c, 0x0a, 0x0c, 0xd1, 0x60, 0x34,
	0x60, 0x14, 0xf2, 0xe0, 0xe2, 0xf6, 0x2d, 0xcd, 0x29, 0xc9, 0x84, 0x9a, 0x23, 0x83, 0x6e, 0x0e,
	0x54, 0x12, 0x62, 0x98, 0x24, 0x2e, 0xc3, 0x8a, 0x95, 0x18, 0x9c, 0xa2, 0x4f, 0x3c, 0x92, 0x63,
	0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0x46, 0x2e, 0xa9, 0xfc, 0xa2, 0x74, 0x3d, 0x70,
	0xe0, 0xe9, 0x81, 0x02, 0x4f, 0x2f, 0xb1, 0x20, 0x53, 0xaf, 0xcc, 0x10, 0xcc, 0x76, 0x82, 0xfa,
	0x3c, 0x80, 0x31, 0x4a, 0x1d, 0x4f, 0x48, 0xa7, 0x17, 0x15, 0x24, 0xc3, 0x62, 0x29, 0x89, 0x0d,
	0x1c, 0x80, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xbe, 0x90, 0xb3, 0x40, 0xc2, 0x01, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RemoveClient is the client API for Remove service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RemoveClient interface {
	Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamRemove(ctx context.Context, opts ...grpc.CallOption) (Remove_StreamRemoveClient, error)
	MultiRemove(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type removeClient struct {
	cc *grpc.ClientConn
}

func NewRemoveClient(cc *grpc.ClientConn) RemoveClient {
	return &removeClient{cc}
}

func (c *removeClient) Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/vald.Remove/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *removeClient) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (Remove_StreamRemoveClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Remove_serviceDesc.Streams[0], "/vald.Remove/StreamRemove", opts...)
	if err != nil {
		return nil, err
	}
	x := &removeStreamRemoveClient{stream}
	return x, nil
}

type Remove_StreamRemoveClient interface {
	Send(*payload.Remove_Request) error
	Recv() (*payload.Object_Location, error)
	grpc.ClientStream
}

type removeStreamRemoveClient struct {
	grpc.ClientStream
}

func (x *removeStreamRemoveClient) Send(m *payload.Remove_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *removeStreamRemoveClient) Recv() (*payload.Object_Location, error) {
	m := new(payload.Object_Location)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *removeClient) MultiRemove(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.Remove/MultiRemove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveServer is the server API for Remove service.
type RemoveServer interface {
	Remove(context.Context, *payload.Remove_Request) (*payload.Object_Location, error)
	StreamRemove(Remove_StreamRemoveServer) error
	MultiRemove(context.Context, *payload.Remove_MultiRequest) (*payload.Object_Locations, error)
}

// UnimplementedRemoveServer can be embedded to have forward compatible implementations.
type UnimplementedRemoveServer struct {
}

func (*UnimplementedRemoveServer) Remove(ctx context.Context, req *payload.Remove_Request) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (*UnimplementedRemoveServer) StreamRemove(srv Remove_StreamRemoveServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamRemove not implemented")
}
func (*UnimplementedRemoveServer) MultiRemove(ctx context.Context, req *payload.Remove_MultiRequest) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiRemove not implemented")
}

func RegisterRemoveServer(s *grpc.Server, srv RemoveServer) {
	s.RegisterService(&_Remove_serviceDesc, srv)
}

func _Remove_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Remove_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoveServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Remove/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoveServer).Remove(ctx, req.(*payload.Remove_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Remove_StreamRemove_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RemoveServer).StreamRemove(&removeStreamRemoveServer{stream})
}

type Remove_StreamRemoveServer interface {
	Send(*payload.Object_Location) error
	Recv() (*payload.Remove_Request, error)
	grpc.ServerStream
}

type removeStreamRemoveServer struct {
	grpc.ServerStream
}

func (x *removeStreamRemoveServer) Send(m *payload.Object_Location) error {
	return x.ServerStream.SendMsg(m)
}

func (x *removeStreamRemoveServer) Recv() (*payload.Remove_Request, error) {
	m := new(payload.Remove_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Remove_MultiRemove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Remove_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoveServer).MultiRemove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Remove/MultiRemove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoveServer).MultiRemove(ctx, req.(*payload.Remove_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Remove_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vald.Remove",
	HandlerType: (*RemoveServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Remove",
			Handler:    _Remove_Remove_Handler,
		},
		{
			MethodName: "MultiRemove",
			Handler:    _Remove_MultiRemove_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamRemove",
			Handler:       _Remove_StreamRemove_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "apis/proto/v1/vald/remove.proto",
}
