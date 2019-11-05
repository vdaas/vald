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

func init() { proto.RegisterFile("egress/egress_filter.proto", fileDescriptor_8d8e16edf70dd8e8) }

var fileDescriptor_8d8e16edf70dd8e8 = []byte{
	// 262 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0x4d, 0x4a, 0x03, 0x31,
	0x14, 0x80, 0x1b, 0x91, 0x2e, 0xc2, 0x14, 0x64, 0x56, 0x12, 0xca, 0x08, 0x5d, 0xb9, 0x31, 0x11,
	0xdd, 0xbb, 0x28, 0x2a, 0xee, 0x14, 0x0b, 0x2e, 0xdc, 0xc8, 0x9b, 0x4c, 0x4c, 0x23, 0x99, 0xbc,
	0x98, 0xc4, 0x82, 0x97, 0xf1, 0x0c, 0x1e, 0xc3, 0xa5, 0x47, 0x28, 0x73, 0x12, 0xe9, 0x64, 0x84,
	0xaa, 0xb8, 0xca, 0xcf, 0xc7, 0xfb, 0xe0, 0x7d, 0x94, 0x29, 0x1d, 0x54, 0x8c, 0x22, 0x1f, 0x0f,
	0x8f, 0xc6, 0x26, 0x15, 0xb8, 0x0f, 0x98, 0xb0, 0x9c, 0xfc, 0xf8, 0x64, 0x13, 0x0f, 0xaf, 0x16,
	0xa1, 0xc9, 0x94, 0x4d, 0x35, 0xa2, 0xb6, 0x4a, 0x80, 0x37, 0x02, 0x9c, 0xc3, 0x04, 0xc9, 0xa0,
	0x8b, 0x03, 0x2d, 0x7c, 0x2d, 0xf4, 0xb3, 0xcd, 0xaf, 0x93, 0x37, 0x42, 0x8b, 0x8b, 0x5e, 0x76,
	0xd9, 0xbb, 0xca, 0x33, 0x3a, 0x1e, 0x6e, 0xfb, 0xfc, 0x5b, 0xbb, 0x50, 0x10, 0xe4, 0x92, 0xdf,
	0xaa, 0xe8, 0xd1, 0x45, 0xc5, 0xfe, 0x25, 0xb3, 0x51, 0x79, 0x45, 0x8b, 0x45, 0x0a, 0x0a, 0xda,
	0x3f, 0x96, 0xeb, 0xfa, 0x49, 0xc9, 0xc4, 0xcf, 0x4d, 0x4c, 0xe0, 0xe4, 0xb6, 0xe5, 0x17, 0x99,
	0x8d, 0x0e, 0xc9, 0x31, 0x61, 0xbb, 0xef, 0xeb, 0x83, 0x9d, 0x79, 0xfb, 0xd1, 0x55, 0xe4, 0xb3,
	0xab, 0xc8, 0xba, 0xab, 0x08, 0x9d, 0x62, 0xd0, 0x7c, 0xd5, 0x00, 0x44, 0xbe, 0x02, 0xdb, 0xf0,
	0xa1, 0x4a, 0xce, 0x31, 0xdf, 0xbb, 0x03, 0xdb, 0x6c, 0x6f, 0x73, 0x43, 0xee, 0x8f, 0xb4, 0x49,
	0xcb, 0x97, 0x9a, 0x4b, 0x6c, 0x45, 0x3f, 0x28, 0x36, 0x83, 0x9b, 0x36, 0x51, 0xe8, 0xe0, 0xa5,
	0xc8, 0x8a, 0x21, 0x73, 0x3d, 0xee, 0xb3, 0x9c, 0x7e, 0x05, 0x00, 0x00, 0xff, 0xff, 0xd5, 0x4a,
	0xa7, 0x4b, 0x7e, 0x01, 0x00, 0x00,
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
	Metadata: "egress/egress_filter.proto",
}
