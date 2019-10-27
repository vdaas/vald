//
// Copyright (C) 2019 kpango (Yusuke Kato)
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

package egress

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

func init() { proto.RegisterFile("egress_filter.proto", fileDescriptor_8cce019dd038e049) }

var fileDescriptor_8cce019dd038e049 = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xb1, 0x4a, 0x03, 0x31,
	0x1c, 0x87, 0x1b, 0x91, 0x0e, 0xe1, 0x0a, 0x72, 0x2e, 0x72, 0x94, 0x13, 0x3a, 0xb9, 0x98, 0x88,
	0xee, 0x0e, 0x45, 0xc5, 0x4d, 0xb1, 0xe0, 0xe0, 0x22, 0xff, 0xcb, 0xc5, 0x34, 0x92, 0xcb, 0x3f,
	0x26, 0xb1, 0xe0, 0xcb, 0xf8, 0x0c, 0x3e, 0x86, 0xa3, 0x8f, 0x50, 0xee, 0x49, 0xa4, 0x97, 0x13,
	0x4e, 0xa5, 0x5b, 0x92, 0x8f, 0xdf, 0x07, 0xf9, 0xe8, 0xbe, 0x54, 0x5e, 0x86, 0xf0, 0xf8, 0xa4,
	0x4d, 0x94, 0x9e, 0x39, 0x8f, 0x11, 0xf3, 0xc9, 0xaf, 0xc7, 0x62, 0xe2, 0xe0, 0xcd, 0x20, 0xd4,
	0x89, 0x16, 0x53, 0x85, 0xa8, 0x8c, 0xe4, 0xe0, 0x34, 0x07, 0x6b, 0x31, 0x42, 0xd4, 0x68, 0x43,
	0x4f, 0x33, 0x57, 0x71, 0xf5, 0x62, 0xd2, 0xed, 0xf4, 0x9d, 0xd0, 0xec, 0xb2, 0x93, 0x5d, 0x75,
	0xae, 0xfc, 0x9c, 0x8e, 0xfb, 0xd3, 0x01, 0xfb, 0xd1, 0x2e, 0x24, 0x78, 0xb1, 0x64, 0x77, 0x32,
	0x38, 0xb4, 0x41, 0x16, 0x5b, 0xc9, 0x6c, 0x94, 0x5f, 0xd3, 0x6c, 0x11, 0xbd, 0x84, 0xe6, 0x9f,
	0xe5, 0xa6, 0x7a, 0x96, 0x22, 0xb2, 0x0b, 0x1d, 0x22, 0x58, 0x31, 0xb4, 0xfc, 0x21, 0xb3, 0xd1,
	0x11, 0x39, 0x21, 0xc5, 0xee, 0xc7, 0xfa, 0x70, 0x67, 0xde, 0x7c, 0xb6, 0x25, 0xf9, 0x6a, 0x4b,
	0xb2, 0x6e, 0x4b, 0x42, 0xa7, 0xe8, 0x15, 0x5b, 0xd5, 0x00, 0x81, 0xad, 0xc0, 0xd4, 0xac, 0xaf,
	0x92, 0x72, 0xcc, 0xf7, 0xee, 0xc1, 0xd4, 0xc3, 0xdf, 0xdc, 0x92, 0x87, 0x63, 0xa5, 0xe3, 0xf2,
	0xb5, 0x62, 0x02, 0x1b, 0xde, 0x0d, 0xf9, 0x66, 0xb8, 0x69, 0x13, 0xb8, 0xf2, 0x4e, 0xf0, 0xa4,
	0xe0, 0x49, 0x51, 0x8d, 0xbb, 0x2c, 0x67, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x7b, 0xcd, 0x0f,
	0xbe, 0x77, 0x01, 0x00, 0x00,
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
	Filter(ctx context.Context, in *payload.Search_Response, opts ...grpc.CallOption) (*payload.Search_Response, error)
	StreamFilter(ctx context.Context, opts ...grpc.CallOption) (EgressFilter_StreamFilterClient, error)
}

type egressFilterClient struct {
	cc *grpc.ClientConn
}

func NewEgressFilterClient(cc *grpc.ClientConn) EgressFilterClient {
	return &egressFilterClient{cc}
}

func (c *egressFilterClient) Filter(ctx context.Context, in *payload.Search_Response, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
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
	Filter(context.Context, *payload.Search_Response) (*payload.Search_Response, error)
	StreamFilter(EgressFilter_StreamFilterServer) error
}

// UnimplementedEgressFilterServer can be embedded to have forward compatible implementations.
type UnimplementedEgressFilterServer struct {
}

func (*UnimplementedEgressFilterServer) Filter(ctx context.Context, req *payload.Search_Response) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Filter not implemented")
}
func (*UnimplementedEgressFilterServer) StreamFilter(srv EgressFilter_StreamFilterServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamFilter not implemented")
}

func RegisterEgressFilterServer(s *grpc.Server, srv EgressFilterServer) {
	s.RegisterService(&_EgressFilter_serviceDesc, srv)
}

func _EgressFilter_Filter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Search_Response)
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
		return srv.(EgressFilterServer).Filter(ctx, req.(*payload.Search_Response))
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
	Metadata: "egress_filter.proto",
}
