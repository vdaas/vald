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

package discoverer

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

func init() {
	proto.RegisterFile("apis/proto/v1/discoverer/discoverer.proto", fileDescriptor_374200cbacdb4f39)
}
func init() {
	golang_proto.RegisterFile("apis/proto/v1/discoverer/discoverer.proto", fileDescriptor_374200cbacdb4f39)
}

var fileDescriptor_374200cbacdb4f39 = []byte{
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0x3f, 0x4b, 0x03, 0x31,
	0x18, 0xc6, 0x89, 0xa8, 0x43, 0xc0, 0x0a, 0xe7, 0x1f, 0xf0, 0x86, 0x20, 0xe7, 0x64, 0xc1, 0x84,
	0xea, 0x20, 0x74, 0x2c, 0x2e, 0x2e, 0x52, 0x1c, 0x44, 0x74, 0x4a, 0x9b, 0x34, 0x06, 0xce, 0x7b,
	0xe3, 0xe5, 0x2e, 0xe2, 0xea, 0x57, 0xf0, 0x0b, 0x09, 0x2e, 0x1d, 0x05, 0xbf, 0x80, 0x5c, 0xfd,
	0x20, 0x72, 0x89, 0xed, 0x5d, 0x75, 0x70, 0xca, 0x9b, 0x3c, 0xef, 0xef, 0xc9, 0xcb, 0xf3, 0xe2,
	0x43, 0x6e, 0xb4, 0x65, 0x26, 0x87, 0x02, 0x98, 0xeb, 0x31, 0xa1, 0xed, 0x18, 0x9c, 0xcc, 0x65,
	0xde, 0x2a, 0xa9, 0x97, 0xa3, 0x8d, 0xd6, 0x8b, 0xeb, 0xc5, 0x07, 0xcb, 0xa4, 0xe1, 0x4f, 0x29,
	0x70, 0x31, 0x3f, 0x03, 0x13, 0x1f, 0x29, 0x5d, 0xdc, 0x95, 0x23, 0x3a, 0x86, 0x7b, 0xa6, 0x40,
	0x41, 0xe8, 0x1f, 0x95, 0x13, 0x7f, 0x0b, 0x70, 0x5d, 0xfd, 0xb4, 0x9f, 0xfe, 0x6e, 0x57, 0x00,
	0x2a, 0x95, 0xfe, 0xa7, 0x50, 0x32, 0x6e, 0x34, 0xe3, 0x59, 0x06, 0x05, 0x2f, 0x34, 0x64, 0x36,
	0x80, 0xc7, 0x6f, 0x08, 0xe3, 0xb3, 0xc5, 0x78, 0xd1, 0x35, 0x5e, 0x1d, 0x82, 0xb0, 0x11, 0xa1,
	0xf3, 0x71, 0x5c, 0x8f, 0x36, 0x3a, 0xbd, 0x94, 0x0f, 0xa5, 0xb4, 0x45, 0xbc, 0xd3, 0xd6, 0xcf,
	0xb3, 0x09, 0xd0, 0x1a, 0x4b, 0xf6, 0x9e, 0x3f, 0xbe, 0x5e, 0x56, 0xb6, 0x92, 0xce, 0x22, 0x04,
	0x66, 0x40, 0xd8, 0x3e, 0xea, 0x46, 0xb7, 0x78, 0xed, 0x02, 0x84, 0xfc, 0xdf, 0x7a, 0xf7, 0x8f,
	0xb5, 0xe7, 0x92, 0xd8, 0x7b, 0x6f, 0x27, 0x9b, 0x8d, 0x77, 0x56, 0x0b, 0x7d, 0xd4, 0x1d, 0x3c,
	0x4e, 0x2b, 0x82, 0xde, 0x2b, 0x82, 0x3e, 0x2b, 0x82, 0x5e, 0x67, 0x04, 0x4d, 0x67, 0x04, 0xe1,
	0x7d, 0xc8, 0x15, 0x75, 0x82, 0x73, 0x4b, 0x1d, 0x4f, 0x05, 0xe5, 0x46, 0xd7, 0x9e, 0xcd, 0x36,
	0x06, 0x9d, 0x2b, 0x9e, 0x8a, 0x66, 0x86, 0x21, 0xba, 0x69, 0xe7, 0xee, 0x51, 0x56, 0xa3, 0x2c,
	0x04, 0x99, 0x9b, 0xf1, 0xf2, 0xae, 0x47, 0xeb, 0x3e, 0xc5, 0x93, 0xef, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x78, 0x1f, 0xc3, 0x18, 0x0e, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DiscovererClient is the client API for Discoverer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DiscovererClient interface {
	Pods(ctx context.Context, in *payload.Discoverer_Request, opts ...grpc.CallOption) (*payload.Info_Pods, error)
	Nodes(ctx context.Context, in *payload.Discoverer_Request, opts ...grpc.CallOption) (*payload.Info_Nodes, error)
}

type discovererClient struct {
	cc *grpc.ClientConn
}

func NewDiscovererClient(cc *grpc.ClientConn) DiscovererClient {
	return &discovererClient{cc}
}

func (c *discovererClient) Pods(ctx context.Context, in *payload.Discoverer_Request, opts ...grpc.CallOption) (*payload.Info_Pods, error) {
	out := new(payload.Info_Pods)
	err := c.cc.Invoke(ctx, "/discoverer.v1.Discoverer/Pods", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discovererClient) Nodes(ctx context.Context, in *payload.Discoverer_Request, opts ...grpc.CallOption) (*payload.Info_Nodes, error) {
	out := new(payload.Info_Nodes)
	err := c.cc.Invoke(ctx, "/discoverer.v1.Discoverer/Nodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiscovererServer is the server API for Discoverer service.
type DiscovererServer interface {
	Pods(context.Context, *payload.Discoverer_Request) (*payload.Info_Pods, error)
	Nodes(context.Context, *payload.Discoverer_Request) (*payload.Info_Nodes, error)
}

// UnimplementedDiscovererServer can be embedded to have forward compatible implementations.
type UnimplementedDiscovererServer struct {
}

func (*UnimplementedDiscovererServer) Pods(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Pods, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pods not implemented")
}
func (*UnimplementedDiscovererServer) Nodes(ctx context.Context, req *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Nodes not implemented")
}

func RegisterDiscovererServer(s *grpc.Server, srv DiscovererServer) {
	s.RegisterService(&_Discoverer_serviceDesc, srv)
}

func _Discoverer_Pods_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Discoverer_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscovererServer).Pods(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discoverer.v1.Discoverer/Pods",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscovererServer).Pods(ctx, req.(*payload.Discoverer_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Discoverer_Nodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Discoverer_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscovererServer).Nodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discoverer.v1.Discoverer/Nodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscovererServer).Nodes(ctx, req.(*payload.Discoverer_Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Discoverer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "discoverer.v1.Discoverer",
	HandlerType: (*DiscovererServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Pods",
			Handler:    _Discoverer_Pods_Handler,
		},
		{
			MethodName: "Nodes",
			Handler:    _Discoverer_Nodes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/proto/v1/discoverer/discoverer.proto",
}
