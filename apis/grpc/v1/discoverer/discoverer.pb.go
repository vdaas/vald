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
	proto "github.com/gogo/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
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
	// 273 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4c, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0x4f, 0xc9, 0x2c, 0x4e, 0xce, 0x2f,
	0x4b, 0x2d, 0x4a, 0x2d, 0x42, 0x62, 0xea, 0x81, 0xa5, 0x85, 0xb8, 0x10, 0x22, 0x52, 0xca, 0xa8,
	0xda, 0x0a, 0x12, 0x2b, 0x73, 0xf2, 0x13, 0x53, 0x60, 0x34, 0x44, 0x83, 0x94, 0x4c, 0x7a, 0x7e,
	0x7e, 0x7a, 0x4e, 0xaa, 0x7e, 0x62, 0x41, 0xa6, 0x7e, 0x62, 0x5e, 0x5e, 0x7e, 0x49, 0x62, 0x49,
	0x66, 0x7e, 0x5e, 0x31, 0x44, 0xd6, 0x68, 0x3b, 0x23, 0x17, 0x97, 0x0b, 0xdc, 0x44, 0xa1, 0x20,
	0x2e, 0x96, 0x80, 0xfc, 0x94, 0x62, 0x21, 0x69, 0x3d, 0x98, 0x21, 0x08, 0x49, 0xbd, 0xa0, 0xd4,
	0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x29, 0x21, 0xb8, 0xa4, 0x67, 0x5e, 0x5a, 0xbe, 0x1e, 0x48, 0x83,
	0x92, 0x64, 0xd3, 0xe5, 0x27, 0x93, 0x99, 0x84, 0x95, 0xf8, 0xe0, 0x2e, 0xd6, 0x2f, 0xc8, 0x4f,
	0x29, 0xb6, 0x62, 0xd4, 0x12, 0x0a, 0xe5, 0x62, 0xf5, 0xcb, 0x4f, 0x49, 0x25, 0x60, 0xa8, 0x30,
	0xaa, 0xa1, 0x60, 0x1d, 0x4a, 0x52, 0x60, 0x53, 0x45, 0x94, 0xf8, 0x11, 0xa6, 0xe6, 0x81, 0x24,
	0xac, 0x18, 0xb5, 0x9c, 0xb2, 0x4f, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23,
	0x39, 0x46, 0x2e, 0x85, 0xfc, 0xa2, 0x74, 0xbd, 0xb2, 0x94, 0xc4, 0xc4, 0x62, 0xbd, 0xb2, 0xc4,
	0x9c, 0x14, 0xbd, 0xc4, 0x82, 0x4c, 0xbd, 0x32, 0x43, 0x3d, 0x44, 0x60, 0x39, 0x21, 0x79, 0x33,
	0x80, 0x31, 0x4a, 0x37, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x1f, 0xac,
	0x4d, 0x1f, 0xa4, 0x4d, 0x1f, 0x1c, 0xa0, 0xe9, 0x45, 0x05, 0xc9, 0xa8, 0xd1, 0x90, 0xc4, 0x06,
	0x0e, 0x2d, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3a, 0x8a, 0x3f, 0xfd, 0xa9, 0x01, 0x00,
	0x00,
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
