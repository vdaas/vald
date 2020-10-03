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
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x91, 0xb1, 0x4b, 0xc4, 0x30,
	0x18, 0xc5, 0x2f, 0x20, 0x37, 0x44, 0xa7, 0x2e, 0x6a, 0x39, 0x4e, 0x38, 0x07, 0x9d, 0xbe, 0x78,
	0x3a, 0xbb, 0xdc, 0xa2, 0x83, 0xa2, 0x9c, 0x9b, 0x4e, 0x5f, 0xdb, 0x10, 0x23, 0x6d, 0xbf, 0x98,
	0xa4, 0x01, 0xff, 0x43, 0x47, 0x67, 0x27, 0xe9, 0x5f, 0x22, 0x97, 0xb4, 0x07, 0x0a, 0x4e, 0x4e,
	0x79, 0xbc, 0xf7, 0xf8, 0x41, 0xde, 0xc7, 0x8f, 0xd0, 0x68, 0x27, 0x8c, 0x25, 0x4f, 0x22, 0x2c,
	0x45, 0xc0, 0xba, 0x12, 0x56, 0x36, 0x14, 0x24, 0x44, 0x33, 0xdb, 0xd9, 0x58, 0xf9, 0xf1, 0xcf,
	0x9a, 0xc1, 0xb7, 0x9a, 0xb0, 0x1a, 0xdf, 0x54, 0xcd, 0x67, 0x8a, 0x48, 0xd5, 0x52, 0xa0, 0xd1,
	0x02, 0xdb, 0x96, 0x3c, 0x7a, 0x4d, 0xad, 0x4b, 0xe9, 0xf9, 0x27, 0xe3, 0xd3, 0x75, 0x24, 0x67,
	0x97, 0x5b, 0xb5, 0x0f, 0x23, 0x22, 0x19, 0xb0, 0x96, 0xaf, 0x9d, 0x74, 0x3e, 0x3f, 0xd8, 0x06,
	0x77, 0xc5, 0x8b, 0x2c, 0x3d, 0xdc, 0x50, 0x19, 0x71, 0x8b, 0x49, 0x76, 0xc5, 0xf7, 0x1e, 0xbc,
	0x95, 0xd8, 0xfc, 0x03, 0x72, 0xca, 0xce, 0x58, 0x76, 0xcd, 0x77, 0x6f, 0xbb, 0xda, 0xeb, 0x81,
	0x33, 0xfb, 0xcd, 0x19, 0xc2, 0x04, 0x3b, 0xfc, 0x0b, 0xe6, 0x16, 0x93, 0xd5, 0xd3, 0x7b, 0x3f,
	0x67, 0x1f, 0xfd, 0x9c, 0x7d, 0xf5, 0x73, 0xc6, 0x73, 0xb2, 0x0a, 0x42, 0x85, 0xe8, 0x60, 0xb3,
	0x1e, 0xa0, 0xd1, 0x10, 0x96, 0x51, 0xaf, 0x86, 0x9f, 0xdf, 0xb3, 0xc7, 0x13, 0xa5, 0xfd, 0x73,
	0x57, 0x40, 0x49, 0x8d, 0x88, 0xe5, 0xb4, 0x7e, 0x5c, 0x5a, 0x59, 0x53, 0x8e, 0xf7, 0x28, 0xa6,
	0x71, 0xc0, 0x8b, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6c, 0x0e, 0xfa, 0x77, 0xac, 0x01, 0x00,
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
