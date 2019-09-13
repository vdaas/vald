//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	_ "github.com/danielvladco/go-proto-gql/pb"
	proto "github.com/gogo/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/payload"
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

func init() { proto.RegisterFile("discoverer.proto", fileDescriptor_9fa655cb815aa581) }

var fileDescriptor_9fa655cb815aa581 = []byte{
	// 230 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0xc9, 0x2c, 0x4e,
	0xce, 0x2f, 0x4b, 0x2d, 0x4a, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x42, 0x88,
	0x48, 0xf1, 0x16, 0x24, 0x56, 0xe6, 0xe4, 0x27, 0xa6, 0x40, 0xa4, 0xa4, 0x64, 0xd2, 0xf3, 0xf3,
	0xd3, 0x73, 0x52, 0xf5, 0x13, 0x0b, 0x32, 0xf5, 0x13, 0xf3, 0xf2, 0xf2, 0x4b, 0x12, 0x4b, 0x32,
	0xf3, 0xf3, 0x8a, 0xa1, 0xb2, 0x3c, 0x05, 0x49, 0xfa, 0xe9, 0x85, 0x39, 0x10, 0x9e, 0x51, 0x1c,
	0x17, 0x97, 0x0b, 0xdc, 0x20, 0x21, 0x2f, 0x2e, 0x0e, 0x18, 0x4f, 0x48, 0x54, 0x0f, 0x66, 0xaa,
	0x73, 0x7e, 0x6e, 0x6e, 0x7e, 0x9e, 0x9e, 0x6b, 0x6e, 0x41, 0x49, 0xa5, 0x94, 0x08, 0x5c, 0xd8,
	0x33, 0x2f, 0x2d, 0x5f, 0xcf, 0x31, 0x3d, 0x35, 0xaf, 0xa4, 0x58, 0x49, 0xb0, 0xe9, 0xf2, 0x93,
	0xc9, 0x4c, 0xdc, 0x42, 0x9c, 0xfa, 0x30, 0x67, 0x49, 0xb1, 0x6c, 0x78, 0x20, 0xcf, 0xe4, 0x94,
	0x78, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0x72, 0x49, 0xe6,
	0x17, 0xa5, 0xeb, 0x95, 0xa5, 0x24, 0x26, 0x16, 0xeb, 0x95, 0x25, 0xe6, 0xa4, 0xe8, 0x21, 0xfc,
	0xe0, 0x84, 0xe4, 0x8c, 0x00, 0xc6, 0x28, 0xad, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4,
	0xfc, 0x5c, 0x7d, 0xb0, 0x7a, 0x7d, 0x90, 0x7a, 0x90, 0x8f, 0x8a, 0xf5, 0xd3, 0x8b, 0x0a, 0x92,
	0xf5, 0x11, 0x3a, 0x93, 0xd8, 0xc0, 0x3e, 0x31, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xe7, 0x3b,
	0x2c, 0xb6, 0x24, 0x01, 0x00, 0x00,
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
	Discover(ctx context.Context, in *payload.Common_Empty, opts ...grpc.CallOption) (*payload.Info_Agents, error)
}

type discovererClient struct {
	cc *grpc.ClientConn
}

func NewDiscovererClient(cc *grpc.ClientConn) DiscovererClient {
	return &discovererClient{cc}
}

func (c *discovererClient) Discover(ctx context.Context, in *payload.Common_Empty, opts ...grpc.CallOption) (*payload.Info_Agents, error) {
	out := new(payload.Info_Agents)
	err := c.cc.Invoke(ctx, "/discoverer.Discoverer/Discover", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiscovererServer is the server API for Discoverer service.
type DiscovererServer interface {
	Discover(context.Context, *payload.Common_Empty) (*payload.Info_Agents, error)
}

// UnimplementedDiscovererServer can be embedded to have forward compatible implementations.
type UnimplementedDiscovererServer struct {
}

func (*UnimplementedDiscovererServer) Discover(ctx context.Context, req *payload.Common_Empty) (*payload.Info_Agents, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Discover not implemented")
}

func RegisterDiscovererServer(s *grpc.Server, srv DiscovererServer) {
	s.RegisterService(&_Discoverer_serviceDesc, srv)
}

func _Discoverer_Discover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Common_Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscovererServer).Discover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/discoverer.Discoverer/Discover",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscovererServer).Discover(ctx, req.(*payload.Common_Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Discoverer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "discoverer.Discoverer",
	HandlerType: (*DiscovererServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Discover",
			Handler:    _Discoverer_Discover_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "discoverer.proto",
}
