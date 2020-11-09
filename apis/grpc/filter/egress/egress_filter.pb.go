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
	// 248 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xb1, 0x4a, 0x04, 0x31,
	0x14, 0x45, 0x37, 0xcd, 0x16, 0x61, 0x17, 0x64, 0x2a, 0x19, 0x96, 0x29, 0xb6, 0xb2, 0x31, 0x11,
	0xed, 0x2d, 0x16, 0x15, 0x3b, 0xc5, 0x05, 0x0b, 0x1b, 0x79, 0x93, 0x79, 0x66, 0x23, 0x99, 0xbc,
	0x90, 0xc4, 0x05, 0xff, 0xc4, 0x4f, 0xb2, 0xf4, 0x13, 0x64, 0xbe, 0x44, 0x36, 0x19, 0x61, 0x55,
	0xac, 0x12, 0xde, 0xe5, 0x1e, 0xb8, 0x87, 0xd7, 0xa8, 0x03, 0xc6, 0x28, 0xcb, 0xf3, 0xf8, 0x64,
	0x6c, 0xc2, 0x20, 0x7c, 0xa0, 0x44, 0xd5, 0xfc, 0xc7, 0xb1, 0x9e, 0x7b, 0x78, 0xb5, 0x04, 0x5d,
	0x49, 0xeb, 0x85, 0x26, 0xd2, 0x16, 0x25, 0x78, 0x23, 0xc1, 0x39, 0x4a, 0x90, 0x0c, 0xb9, 0x58,
	0xd2, 0xd3, 0x37, 0xc6, 0x67, 0x97, 0xb9, 0x7e, 0x95, 0xdb, 0xd5, 0x39, 0x9f, 0x8e, 0xbf, 0x43,
	0xf1, 0x0d, 0x5a, 0x23, 0x04, 0xb5, 0x11, 0x77, 0x18, 0x3d, 0xb9, 0x88, 0xf5, 0xbf, 0xc9, 0x72,
	0x52, 0x5d, 0xf3, 0xd9, 0x3a, 0x05, 0x84, 0xfe, 0x0f, 0xe5, 0xa6, 0x7d, 0x46, 0x95, 0xc4, 0x85,
	0x89, 0x09, 0x9c, 0xda, 0xa7, 0xfc, 0x4a, 0x96, 0x93, 0x23, 0x76, 0xc2, 0x56, 0xfd, 0xfb, 0xd0,
	0xb0, 0x8f, 0xa1, 0x61, 0x9f, 0x43, 0xc3, 0xf8, 0x82, 0x82, 0x16, 0xdb, 0x0e, 0x20, 0x8a, 0x2d,
	0xd8, 0x4e, 0x8c, 0x06, 0xca, 0xf4, 0xd5, 0xc1, 0x3d, 0xd8, 0x6e, 0x7f, 0xc7, 0x2d, 0x7b, 0x38,
	0xd6, 0x26, 0x6d, 0x5e, 0x5a, 0xa1, 0xa8, 0x97, 0xb9, 0x28, 0x77, 0xc5, 0x9d, 0x87, 0x28, 0x75,
	0xf0, 0x4a, 0x16, 0xc4, 0xa8, 0xb4, 0x9d, 0x66, 0x21, 0x67, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff,
	0xa6, 0x83, 0x64, 0xdf, 0x6a, 0x01, 0x00, 0x00,
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
