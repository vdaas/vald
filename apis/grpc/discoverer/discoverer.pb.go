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

package discoverer

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

func init() { proto.RegisterFile("discoverer.proto", fileDescriptor_9fa655cb815aa581) }

var fileDescriptor_9fa655cb815aa581 = []byte{
	// 270 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0xc9, 0x2c, 0x4e,
	0xce, 0x2f, 0x4b, 0x2d, 0x4a, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x42, 0x88,
	0x48, 0xf1, 0x16, 0x24, 0x56, 0xe6, 0xe4, 0x27, 0xa6, 0x40, 0xa4, 0xa4, 0x64, 0xd2, 0xf3, 0xf3,
	0xd3, 0x73, 0x52, 0xf5, 0x13, 0x0b, 0x32, 0xf5, 0x13, 0xf3, 0xf2, 0xf2, 0x4b, 0x12, 0x4b, 0x32,
	0xf3, 0xf3, 0x8a, 0xa1, 0xb2, 0x3c, 0x05, 0x49, 0xfa, 0xe9, 0x85, 0x39, 0x10, 0x9e, 0xd1, 0x51,
	0x46, 0x2e, 0x2e, 0x17, 0xb8, 0x49, 0x42, 0x61, 0x5c, 0x2c, 0x01, 0xf9, 0x29, 0xc5, 0x42, 0xd2,
	0x7a, 0x30, 0x23, 0x11, 0x92, 0x7a, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x52, 0x42, 0x70,
	0x49, 0xcf, 0xbc, 0xb4, 0x7c, 0x3d, 0x90, 0x06, 0x25, 0xd9, 0xa6, 0xcb, 0x4f, 0x26, 0x33, 0x89,
	0x0b, 0x89, 0xea, 0xc3, 0xdc, 0xa5, 0x5f, 0x90, 0x9f, 0x52, 0xac, 0x5f, 0x9d, 0x97, 0x98, 0x9b,
	0x5a, 0x2b, 0x14, 0xc9, 0xc5, 0xea, 0x97, 0x9f, 0x92, 0x4a, 0xc0, 0x60, 0x61, 0x54, 0x83, 0xc1,
	0x3a, 0x94, 0xe4, 0xc0, 0x26, 0x4b, 0x08, 0x89, 0x21, 0x4c, 0xce, 0x03, 0x49, 0x40, 0x8d, 0x96,
	0x62, 0xd9, 0xf0, 0x40, 0x9e, 0xc9, 0x29, 0xf1, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18,
	0x1f, 0x3c, 0x92, 0x63, 0xe4, 0x92, 0xcc, 0x2f, 0x4a, 0xd7, 0x2b, 0x4b, 0x49, 0x4c, 0x2c, 0xd6,
	0x2b, 0x4b, 0xcc, 0x49, 0xd1, 0x43, 0x84, 0x95, 0x13, 0x92, 0x6f, 0x03, 0x18, 0xa3, 0xb4, 0xd2,
	0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0xc1, 0xea, 0xf5, 0x41, 0xea, 0x41,
	0x21, 0x57, 0xac, 0x9f, 0x5e, 0x54, 0x90, 0xac, 0x8f, 0xd0, 0x99, 0xc4, 0x06, 0x0e, 0x31, 0x63,
	0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x06, 0xbc, 0x75, 0x19, 0x8c, 0x01, 0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/discoverer.Discoverer/Pods", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discovererClient) Nodes(ctx context.Context, in *payload.Discoverer_Request, opts ...grpc.CallOption) (*payload.Info_Nodes, error) {
	out := new(payload.Info_Nodes)
	err := c.cc.Invoke(ctx, "/discoverer.Discoverer/Nodes", in, out, opts...)
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
		FullMethod: "/discoverer.Discoverer/Pods",
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
		FullMethod: "/discoverer.Discoverer/Nodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscovererServer).Nodes(ctx, req.(*payload.Discoverer_Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Discoverer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "discoverer.Discoverer",
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
	Metadata: "discoverer.proto",
}
