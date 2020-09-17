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

package egress

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
	proto.RegisterFile("apis/proto/v1/filter/egress/egress_filter.proto", fileDescriptor_7f3e67472eb32d70)
}

var fileDescriptor_7f3e67472eb32d70 = []byte{
	// 247 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4f, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0x4f, 0xcb, 0xcc, 0x29, 0x49, 0x2d,
	0xd2, 0x4f, 0x4d, 0x2f, 0x4a, 0x2d, 0x2e, 0x86, 0x52, 0xf1, 0x10, 0x41, 0x3d, 0xb0, 0x22, 0x21,
	0x5e, 0x14, 0x41, 0x29, 0x65, 0x54, 0xfd, 0x05, 0x89, 0x95, 0x39, 0xf9, 0x89, 0x29, 0x30, 0x1a,
	0xa2, 0x47, 0x4a, 0x26, 0x3d, 0x3f, 0x3f, 0x3d, 0x27, 0x15, 0x64, 0x97, 0x7e, 0x62, 0x5e, 0x5e,
	0x7e, 0x49, 0x62, 0x49, 0x66, 0x7e, 0x5e, 0x31, 0x44, 0xd6, 0x68, 0x06, 0x23, 0x17, 0x8f, 0x2b,
	0xd8, 0x50, 0x37, 0xb0, 0x99, 0x42, 0x76, 0x5c, 0x6c, 0x50, 0x96, 0x84, 0x1e, 0xcc, 0x20, 0xff,
	0xa4, 0xac, 0xd4, 0xe4, 0x12, 0x3d, 0x97, 0xcc, 0xe2, 0x92, 0xc4, 0xbc, 0xe4, 0x54, 0x29, 0x9c,
	0x32, 0x4a, 0x0c, 0x42, 0x1e, 0x5c, 0x3c, 0xc1, 0x25, 0x45, 0xa9, 0x89, 0xb9, 0x94, 0x98, 0xa2,
	0xc1, 0x68, 0xc0, 0xe8, 0x54, 0x7e, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e,
	0xc9, 0x31, 0x72, 0x29, 0xe7, 0x17, 0xa5, 0xeb, 0x95, 0xa5, 0x24, 0x26, 0x16, 0xeb, 0x95, 0x25,
	0xe6, 0xa4, 0xe8, 0x25, 0x16, 0x64, 0xea, 0x95, 0x19, 0xea, 0x41, 0x83, 0x07, 0x12, 0x2e, 0x4e,
	0x02, 0x61, 0x89, 0x39, 0x29, 0xc8, 0xde, 0x09, 0x60, 0x8c, 0x32, 0x48, 0xcf, 0x2c, 0xc9, 0x28,
	0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x07, 0xeb, 0xd7, 0x07, 0xe9, 0x87, 0x04, 0x7d, 0x7a, 0x51,
	0x41, 0x32, 0x46, 0xc8, 0x27, 0xb1, 0x81, 0x83, 0xc6, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x1c,
	0xbd, 0x87, 0x2f, 0x9f, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EgressFilterClient is the client API for EgressFilter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EgressFilterClient interface {
	Filter(ctx context.Context, in *payload.Object_Distance, opts ...grpc.CallOption) (*payload.Object_Distance, error)
	StreamFilter(ctx context.Context, opts ...grpc.CallOption) (EgressFilter_StreamFilterClient, error)
}

type egressFilterClient struct {
	cc *grpc.ClientConn
}

func NewEgressFilterClient(cc *grpc.ClientConn) EgressFilterClient {
	return &egressFilterClient{cc}
}

func (c *egressFilterClient) Filter(ctx context.Context, in *payload.Object_Distance, opts ...grpc.CallOption) (*payload.Object_Distance, error) {
	out := new(payload.Object_Distance)
	err := c.cc.Invoke(ctx, "/egress_filter.EgressFilter/Filter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *egressFilterClient) StreamFilter(ctx context.Context, opts ...grpc.CallOption) (EgressFilter_StreamFilterClient, error) {
	stream, err := c.cc.NewStream(ctx, &_EgressFilter_serviceDesc.Streams[0], "/egress_filter.EgressFilter/StreamFilter", opts...)
	if err != nil {
		return nil, err
	}
	x := &egressFilterStreamFilterClient{stream}
	return x, nil
}

type EgressFilter_StreamFilterClient interface {
	Send(*payload.Object_Distance) error
	Recv() (*payload.Object_Distance, error)
	grpc.ClientStream
}

type egressFilterStreamFilterClient struct {
	grpc.ClientStream
}

func (x *egressFilterStreamFilterClient) Send(m *payload.Object_Distance) error {
	return x.ClientStream.SendMsg(m)
}

func (x *egressFilterStreamFilterClient) Recv() (*payload.Object_Distance, error) {
	m := new(payload.Object_Distance)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EgressFilterServer is the server API for EgressFilter service.
type EgressFilterServer interface {
	Filter(context.Context, *payload.Object_Distance) (*payload.Object_Distance, error)
	StreamFilter(EgressFilter_StreamFilterServer) error
}

// UnimplementedEgressFilterServer can be embedded to have forward compatible implementations.
type UnimplementedEgressFilterServer struct {
}

func (*UnimplementedEgressFilterServer) Filter(ctx context.Context, req *payload.Object_Distance) (*payload.Object_Distance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Filter not implemented")
}
func (*UnimplementedEgressFilterServer) StreamFilter(srv EgressFilter_StreamFilterServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamFilter not implemented")
}

func RegisterEgressFilterServer(s *grpc.Server, srv EgressFilterServer) {
	s.RegisterService(&_EgressFilter_serviceDesc, srv)
}

func _EgressFilter_Filter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Distance)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EgressFilterServer).Filter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/egress_filter.EgressFilter/Filter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EgressFilterServer).Filter(ctx, req.(*payload.Object_Distance))
	}
	return interceptor(ctx, in, info, handler)
}

func _EgressFilter_StreamFilter_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EgressFilterServer).StreamFilter(&egressFilterStreamFilterServer{stream})
}

type EgressFilter_StreamFilterServer interface {
	Send(*payload.Object_Distance) error
	Recv() (*payload.Object_Distance, error)
	grpc.ServerStream
}

type egressFilterStreamFilterServer struct {
	grpc.ServerStream
}

func (x *egressFilterStreamFilterServer) Send(m *payload.Object_Distance) error {
	return x.ServerStream.SendMsg(m)
}

func (x *egressFilterStreamFilterServer) Recv() (*payload.Object_Distance, error) {
	m := new(payload.Object_Distance)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _EgressFilter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "egress_filter.EgressFilter",
	HandlerType: (*EgressFilterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Filter",
			Handler:    _EgressFilter_Filter_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamFilter",
			Handler:       _EgressFilter_StreamFilter_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "apis/proto/v1/filter/egress/egress_filter.proto",
}
