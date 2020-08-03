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

func init() { proto.RegisterFile("egress/egress_filter.proto", fileDescriptor_8d8e16edf70dd8e8) }

var fileDescriptor_8d8e16edf70dd8e8 = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4a, 0x4d, 0x2f, 0x4a,
	0x2d, 0x2e, 0xd6, 0x87, 0x50, 0xf1, 0x69, 0x99, 0x39, 0x25, 0xa9, 0x45, 0x7a, 0x05, 0x45, 0xf9,
	0x25, 0xf9, 0x42, 0xbc, 0x28, 0x82, 0x52, 0xbc, 0x05, 0x89, 0x95, 0x39, 0xf9, 0x89, 0x29, 0x10,
	0x59, 0x29, 0x99, 0xf4, 0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0xfd, 0xc4, 0x82, 0x4c, 0xfd, 0xc4, 0xbc,
	0xbc, 0xfc, 0x92, 0xc4, 0x92, 0xcc, 0xfc, 0xbc, 0x62, 0xa8, 0x2c, 0x4f, 0x41, 0x92, 0x7e, 0x7a,
	0x61, 0x0e, 0x84, 0x67, 0x34, 0x8f, 0x91, 0x8b, 0xc7, 0x15, 0x6c, 0x98, 0x1b, 0xd8, 0x2c, 0x21,
	0x3b, 0x2e, 0x36, 0x28, 0x4b, 0x42, 0x0f, 0x66, 0xac, 0x7f, 0x52, 0x56, 0x6a, 0x72, 0x89, 0x9e,
	0x4b, 0x66, 0x71, 0x49, 0x62, 0x5e, 0x72, 0xaa, 0x14, 0x4e, 0x19, 0x25, 0x06, 0x21, 0x0f, 0x2e,
	0x9e, 0xe0, 0x92, 0xa2, 0xd4, 0xc4, 0x5c, 0x4a, 0x4c, 0xd1, 0x60, 0x34, 0x60, 0x94, 0x62, 0xd9,
	0xf0, 0x40, 0x9e, 0xc9, 0x29, 0xf7, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c,
	0x92, 0x63, 0xe4, 0x92, 0xc9, 0x2f, 0x4a, 0xd7, 0x2b, 0x4b, 0x49, 0x4c, 0x2c, 0xd6, 0x2b, 0x4b,
	0xcc, 0x49, 0xd1, 0x83, 0x86, 0x0a, 0x24, 0x38, 0x9c, 0x04, 0xc2, 0x12, 0x73, 0x52, 0x90, 0x7d,
	0x13, 0xc0, 0x18, 0xa5, 0x9b, 0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c, 0x9f, 0xab, 0x0f,
	0xd6, 0xa8, 0x0f, 0xd2, 0x08, 0x0a, 0x9b, 0x62, 0xfd, 0xf4, 0xa2, 0x82, 0x64, 0x7d, 0x88, 0x11,
	0xd0, 0x60, 0x4e, 0x62, 0x03, 0x07, 0x8b, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xdf, 0x68, 0x6c,
	0x9d, 0x7e, 0x01, 0x00, 0x00,
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
	Metadata: "egress/egress_filter.proto",
}
