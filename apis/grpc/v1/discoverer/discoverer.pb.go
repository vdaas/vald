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

func init() {
	proto.RegisterFile("apis/proto/v1/discoverer/discoverer.proto", fileDescriptor_374200cbacdb4f39)
}

var fileDescriptor_374200cbacdb4f39 = []byte{
	// 280 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xbf, 0x4a, 0xc4, 0x30,
	0x18, 0xc0, 0x89, 0xa8, 0x43, 0xc0, 0x13, 0xea, 0x1f, 0xb0, 0x48, 0x91, 0x3a, 0x79, 0x60, 0x42,
	0x75, 0xbb, 0xf1, 0x70, 0x71, 0x91, 0xc3, 0x41, 0x44, 0xa7, 0xef, 0x9a, 0x58, 0x03, 0x35, 0x5f,
	0x4c, 0x72, 0x01, 0x57, 0x5f, 0xc1, 0x47, 0x72, 0x71, 0x14, 0x7c, 0x01, 0x29, 0x3e, 0x88, 0x34,
	0xc7, 0x5d, 0x7b, 0x38, 0x38, 0x25, 0xe4, 0xc7, 0xef, 0x97, 0xf0, 0x85, 0x9e, 0x80, 0x51, 0x8e,
	0x1b, 0x8b, 0x1e, 0x79, 0x28, 0xb8, 0x50, 0xae, 0xc4, 0x20, 0xad, 0xb4, 0xbd, 0x2d, 0x8b, 0x38,
	0xd9, 0xea, 0x9d, 0x84, 0x22, 0x3d, 0x5e, 0x35, 0x0d, 0xbc, 0xd4, 0x08, 0x62, 0xb1, 0xce, 0x9d,
	0xf4, 0xb0, 0x42, 0xac, 0x6a, 0xc9, 0xc1, 0x28, 0x0e, 0x5a, 0xa3, 0x07, 0xaf, 0x50, 0xbb, 0x39,
	0x3d, 0x7b, 0x27, 0x94, 0x5e, 0x2c, 0xa3, 0xc9, 0x2d, 0x5d, 0x9f, 0xa0, 0x70, 0x49, 0xc6, 0x16,
	0x91, 0x50, 0xb0, 0x8e, 0xb3, 0x6b, 0xf9, 0x3c, 0x93, 0xce, 0xa7, 0x7b, 0x7d, 0x7e, 0xa9, 0x1f,
	0x90, 0xb5, 0x5a, 0x7e, 0xf0, 0xfa, 0xf5, 0xf3, 0xb6, 0xb6, 0x93, 0x0f, 0x96, 0x4f, 0xe7, 0x06,
	0x85, 0x1b, 0x91, 0x61, 0x72, 0x4f, 0x37, 0xae, 0x50, 0xc8, 0xff, 0xd3, 0xfb, 0x7f, 0xd2, 0xd1,
	0xcb, 0xd3, 0xd8, 0xde, 0xcd, 0xb7, 0xbb, 0xb6, 0x6e, 0xc1, 0x88, 0x0c, 0xc7, 0xf8, 0xd1, 0x64,
	0xe4, 0xb3, 0xc9, 0xc8, 0x77, 0x93, 0x11, 0x7a, 0x84, 0xb6, 0x62, 0x41, 0x00, 0x38, 0x16, 0xa0,
	0x16, 0x0c, 0x8c, 0x6a, 0x5b, 0xdd, 0xec, 0xc6, 0x83, 0x1b, 0xa8, 0x45, 0x77, 0xf7, 0x84, 0xdc,
	0x9d, 0x56, 0xca, 0x3f, 0xce, 0xa6, 0xac, 0xc4, 0x27, 0x1e, 0x55, 0xde, 0xaa, 0x3c, 0x0e, 0xb8,
	0xb2, 0xa6, 0x5c, 0xfd, 0x99, 0xe9, 0x66, 0x9c, 0xde, 0xf9, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x2f, 0x88, 0x6b, 0xe9, 0xbc, 0x01, 0x00, 0x00,
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
