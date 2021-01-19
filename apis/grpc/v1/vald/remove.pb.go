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

func init() { proto.RegisterFile("apis/proto/v1/vald/remove.proto", fileDescriptor_5b638f34e0c25c81) }
func init() {
	golang_proto.RegisterFile("apis/proto/v1/vald/remove.proto", fileDescriptor_5b638f34e0c25c81)
}

var fileDescriptor_5b638f34e0c25c81 = []byte{
	// 337 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xc1, 0x4a, 0xf3, 0x40,
	0x18, 0xfc, 0xb7, 0x87, 0x16, 0xf2, 0xf7, 0x20, 0x01, 0x2f, 0xa9, 0xa4, 0x18, 0x0f, 0x8a, 0xe0,
	0xae, 0xd5, 0x83, 0xe0, 0xb1, 0x67, 0x45, 0x69, 0xd1, 0x83, 0xb7, 0x2f, 0xc9, 0xba, 0xae, 0x24,
	0xfd, 0xd6, 0x64, 0x13, 0xf0, 0xea, 0x2b, 0xf8, 0x00, 0xbe, 0x8a, 0xc7, 0x1e, 0x05, 0x5f, 0x40,
	0x5a, 0x1f, 0x44, 0x76, 0x37, 0x2d, 0x55, 0x8a, 0x78, 0xca, 0xe4, 0xdb, 0x99, 0xd9, 0x59, 0xe6,
	0xf3, 0xfa, 0xa0, 0x64, 0xc9, 0x54, 0x81, 0x1a, 0x59, 0x3d, 0x60, 0x35, 0x64, 0x29, 0x2b, 0x78,
	0x8e, 0x35, 0xa7, 0x76, 0xe8, 0x77, 0xcc, 0x88, 0xd6, 0x83, 0x60, 0xe7, 0x3b, 0x53, 0xc1, 0x63,
	0x86, 0x90, 0x2e, 0xbe, 0x8e, 0x1d, 0x1c, 0x08, 0xa9, 0xef, 0xaa, 0x98, 0x26, 0x98, 0x33, 0x81,
	0x02, 0x1d, 0x3f, 0xae, 0x6e, 0xed, 0x9f, 0x13, 0x1b, 0xd4, 0xd0, 0x4f, 0x7e, 0xd2, 0x05, 0xa2,
	0xc8, 0xb8, 0xbd, 0xc9, 0x41, 0x06, 0x4a, 0x32, 0x98, 0x4c, 0x50, 0x83, 0x96, 0x38, 0x29, 0x9d,
	0xf0, 0xe8, 0xa5, 0xe5, 0xb5, 0x47, 0x36, 0xa6, 0x7f, 0xb5, 0x44, 0x01, 0x5d, 0x84, 0xa9, 0x07,
	0xd4, 0xcd, 0xe8, 0x88, 0x3f, 0x54, 0xbc, 0xd4, 0x41, 0x6f, 0xf5, 0xec, 0x22, 0xbe, 0xe7, 0x89,
	0xa6, 0x67, 0x98, 0x58, 0xd3, 0xc8, 0x7f, 0x7a, 0xff, 0x7c, 0x6e, 0x75, 0xa3, 0x4e, 0xf3, 0xf4,
	0x53, 0xb2, 0xef, 0x8f, 0xbd, 0xee, 0x58, 0x17, 0x1c, 0xf2, 0x3f, 0x98, 0x6f, 0xaf, 0x31, 0x77,
	0xe2, 0xe5, 0x15, 0xff, 0xf6, 0xc8, 0x21, 0xf1, 0xa5, 0xf7, 0xff, 0xbc, 0xca, 0xb4, 0x6c, 0x3c,
	0xfb, 0x6b, 0x3c, 0x9b, 0x73, 0x67, 0xbc, 0xf5, 0x4b, 0xea, 0x32, 0xea, 0xd9, 0xd8, 0x9b, 0xd1,
	0x46, 0x13, 0x9b, 0xe5, 0x46, 0xab, 0x32, 0x93, 0x7f, 0x28, 0xa6, 0xb3, 0x90, 0xbc, 0xcd, 0x42,
	0xf2, 0x31, 0x0b, 0xc9, 0xeb, 0x3c, 0x24, 0xd3, 0x79, 0x48, 0xbc, 0x00, 0x0b, 0x41, 0xeb, 0x14,
	0xa0, 0xa4, 0xb6, 0x57, 0x50, 0xd2, 0x58, 0x1b, 0x3c, 0xf4, 0xae, 0x21, 0x4b, 0x5d, 0x8a, 0x4b,
	0x72, 0xb3, 0xbb, 0x52, 0x8d, 0x15, 0xb8, 0xdd, 0x70, 0xd5, 0x14, 0x2a, 0x59, 0x6c, 0x4b, 0xdc,
	0xb6, 0x8d, 0x1c, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0xc4, 0x33, 0xc1, 0x2d, 0x4a, 0x02, 0x00,
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
	err := c.cc.Invoke(ctx, "/vald.v1.Remove/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *removeClient) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (Remove_StreamRemoveClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Remove_serviceDesc.Streams[0], "/vald.v1.Remove/StreamRemove", opts...)
	if err != nil {
		return nil, err
	}
	x := &removeStreamRemoveClient{stream}
	return x, nil
}

type Remove_StreamRemoveClient interface {
	Send(*payload.Remove_Request) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type removeStreamRemoveClient struct {
	grpc.ClientStream
}

func (x *removeStreamRemoveClient) Send(m *payload.Remove_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *removeStreamRemoveClient) Recv() (*payload.Object_StreamLocation, error) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *removeClient) MultiRemove(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.v1.Remove/MultiRemove", in, out, opts...)
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
		FullMethod: "/vald.v1.Remove/Remove",
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
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Remove_Request, error)
	grpc.ServerStream
}

type removeStreamRemoveServer struct {
	grpc.ServerStream
}

func (x *removeStreamRemoveServer) Send(m *payload.Object_StreamLocation) error {
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
		FullMethod: "/vald.v1.Remove/MultiRemove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemoveServer).MultiRemove(ctx, req.(*payload.Remove_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Remove_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vald.v1.Remove",
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
