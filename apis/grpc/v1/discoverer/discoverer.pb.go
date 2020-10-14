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
	// 277 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0x31, 0x4a, 0x34, 0x31,
	0x14, 0x80, 0xc9, 0xcf, 0xaf, 0x45, 0x8a, 0x15, 0x66, 0x6d, 0x1c, 0x65, 0x90, 0xb1, 0x72, 0xc1,
	0x17, 0x56, 0xbb, 0x2d, 0x17, 0x1b, 0x1b, 0x59, 0x16, 0xb4, 0xb0, 0x7b, 0x3b, 0x89, 0x63, 0x60,
	0x9c, 0x17, 0x27, 0xd9, 0x80, 0xad, 0x57, 0xf0, 0x30, 0x5e, 0xc1, 0x52, 0xf0, 0x02, 0x32, 0x78,
	0x10, 0x99, 0x2c, 0x3b, 0x71, 0x2a, 0xab, 0x84, 0x7c, 0x7c, 0x5f, 0xc2, 0x0b, 0x3f, 0x45, 0xa3,
	0xad, 0x30, 0x0d, 0x39, 0x12, 0x7e, 0x2a, 0xa4, 0xb6, 0x05, 0x79, 0xd5, 0xa8, 0xe6, 0xd7, 0x16,
	0x02, 0x4e, 0x78, 0x3c, 0x49, 0x4f, 0x86, 0x9a, 0xc1, 0xe7, 0x8a, 0x50, 0x6e, 0xd7, 0x8d, 0x90,
	0x1e, 0x95, 0x44, 0x65, 0xa5, 0x04, 0x1a, 0x2d, 0xb0, 0xae, 0xc9, 0xa1, 0xd3, 0x54, 0xdb, 0x0d,
	0x3d, 0x7f, 0x63, 0x9c, 0x5f, 0xf6, 0xc5, 0x64, 0xc9, 0xff, 0x2f, 0x48, 0xda, 0xe4, 0x10, 0xb6,
	0x91, 0x08, 0x61, 0xa9, 0x9e, 0xd6, 0xca, 0xba, 0x34, 0xe9, 0xe1, 0x55, 0x7d, 0x4f, 0xd0, 0x09,
	0xf9, 0xc1, 0xcb, 0xe7, 0xf7, 0xeb, 0xbf, 0x71, 0x3e, 0xea, 0x5f, 0x2c, 0x0c, 0x49, 0x3b, 0x63,
	0x93, 0xe4, 0x86, 0xef, 0x5c, 0x93, 0x54, 0x7f, 0x44, 0xc7, 0xc3, 0x68, 0x30, 0xf2, 0x34, 0x54,
	0xf7, 0xf3, 0xbd, 0x58, 0xad, 0x3b, 0x30, 0x63, 0x93, 0x39, 0xbd, 0xb7, 0x19, 0xfb, 0x68, 0x33,
	0xf6, 0xd5, 0x66, 0x8c, 0x1f, 0x53, 0x53, 0x82, 0x97, 0x88, 0x16, 0x3c, 0x56, 0x12, 0xd0, 0x68,
	0xf0, 0x53, 0x88, 0xc3, 0x9a, 0x8f, 0x6e, 0xb1, 0x92, 0xf1, 0xe2, 0x05, 0xbb, 0x3b, 0x2b, 0xb5,
	0x7b, 0x58, 0xaf, 0xa0, 0xa0, 0x47, 0x11, 0x54, 0xd1, 0xa9, 0x22, 0x0c, 0xb5, 0x6c, 0x4c, 0x31,
	0xfc, 0x8a, 0xd5, 0x6e, 0x98, 0xd8, 0xc5, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc3, 0x7d, 0x80,
	0x44, 0xad, 0x01, 0x00, 0x00,
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
	Metadata: "apis/proto/v1/discoverer/discoverer.proto",
}
