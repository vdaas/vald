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

package ingress

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
	proto.RegisterFile("apis/proto/v1/filter/ingress/ingress_filter.proto", fileDescriptor_8b82e91ce4fe335b)
}

var fileDescriptor_8b82e91ce4fe335b = []byte{
	// 305 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0x31, 0x4b, 0x43, 0x31,
	0x14, 0x85, 0x4d, 0x07, 0xc5, 0x50, 0x91, 0x66, 0x11, 0x8b, 0x74, 0xa8, 0x0e, 0xe2, 0x90, 0xf8,
	0x74, 0x73, 0xec, 0xa0, 0x88, 0x83, 0x85, 0x42, 0x11, 0x17, 0xb9, 0xef, 0xbd, 0x18, 0x23, 0xe9,
	0xbb, 0x21, 0x2f, 0x06, 0x74, 0xf4, 0x2f, 0xf8, 0xa7, 0x1c, 0x05, 0x67, 0x41, 0x8a, 0x3f, 0x44,
	0x9a, 0xb4, 0xc3, 0x43, 0x05, 0x41, 0xa7, 0x4b, 0xee, 0xc9, 0xf9, 0x38, 0x70, 0x2e, 0xcd, 0xc0,
	0xea, 0x5a, 0x58, 0x87, 0x1e, 0x45, 0xc8, 0xc4, 0xb5, 0x36, 0x5e, 0x3a, 0xa1, 0x2b, 0xe5, 0x64,
	0x5d, 0x2f, 0xe6, 0x55, 0x5a, 0xf3, 0xf8, 0x8d, 0x75, 0xe6, 0xaf, 0xb9, 0xc8, 0x43, 0xd6, 0xdd,
	0x6e, 0x52, 0x2c, 0xdc, 0x1b, 0x84, 0x72, 0x31, 0x93, 0xaf, 0xbb, 0xa5, 0x10, 0x95, 0x91, 0x02,
	0xac, 0x16, 0x50, 0x55, 0xe8, 0xc1, 0x6b, 0xac, 0xea, 0xa4, 0x1e, 0xbc, 0xb5, 0xe8, 0xda, 0x69,
	0x22, 0x1e, 0x47, 0x3e, 0x1b, 0xd1, 0xd5, 0x13, 0x59, 0x8d, 0x65, 0xe1, 0xd1, 0xb1, 0x0d, 0xbe,
	0x80, 0x85, 0x8c, 0x9f, 0xe7, 0xb7, 0xb2, 0xf0, 0x7c, 0x60, 0x30, 0xef, 0x6e, 0x7e, 0x23, 0x24,
	0x4f, 0x9f, 0x3d, 0xbe, 0x7e, 0x3c, 0xb5, 0xda, 0xfd, 0x15, 0x81, 0x71, 0x7f, 0x44, 0xf6, 0xd8,
	0x19, 0x5d, 0x1f, 0x79, 0x27, 0x61, 0xf2, 0x37, 0xf4, 0xd2, 0x2e, 0xd9, 0x27, 0xec, 0x82, 0xb6,
	0x53, 0xd6, 0x39, 0xe9, 0x67, 0xc3, 0xef, 0x62, 0x86, 0xb8, 0x98, 0xc5, 0x1c, 0x52, 0x96, 0x62,
	0xfe, 0x03, 0x3f, 0x66, 0x1d, 0x3c, 0x3c, 0x4f, 0x7b, 0xe4, 0x65, 0xda, 0x23, 0xef, 0xd3, 0x1e,
	0xa1, 0x3b, 0xe8, 0x14, 0x0f, 0x25, 0x40, 0xcd, 0x03, 0x98, 0x92, 0x83, 0xd5, 0x33, 0x63, 0xb3,
	0xd9, 0x41, 0x67, 0x0c, 0xa6, 0x6c, 0x94, 0x32, 0x24, 0x97, 0x99, 0xd2, 0xfe, 0xe6, 0x2e, 0xe7,
	0x05, 0x4e, 0x44, 0x24, 0x88, 0x19, 0x41, 0xc4, 0x0b, 0x50, 0xce, 0x16, 0x5f, 0xcf, 0x28, 0x5f,
	0x8e, 0x15, 0x1f, 0x7e, 0x06, 0x00, 0x00, 0xff, 0xff, 0x06, 0x3e, 0x7c, 0x2a, 0x6d, 0x02, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// IngressFilterClient is the client API for IngressFilter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IngressFilterClient interface {
	GenVector(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Vector, error)
	StreamGenVector(ctx context.Context, opts ...grpc.CallOption) (IngressFilter_StreamGenVectorClient, error)
	FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Object_Vector, error)
	StreamFilterVector(ctx context.Context, opts ...grpc.CallOption) (IngressFilter_StreamFilterVectorClient, error)
}

type ingressFilterClient struct {
	cc *grpc.ClientConn
}

func NewIngressFilterClient(cc *grpc.ClientConn) IngressFilterClient {
	return &ingressFilterClient{cc}
}

func (c *ingressFilterClient) GenVector(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/filter.ingress.v1.IngressFilter/GenVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingressFilterClient) StreamGenVector(ctx context.Context, opts ...grpc.CallOption) (IngressFilter_StreamGenVectorClient, error) {
	stream, err := c.cc.NewStream(ctx, &_IngressFilter_serviceDesc.Streams[0], "/filter.ingress.v1.IngressFilter/StreamGenVector", opts...)
	if err != nil {
		return nil, err
	}
	x := &ingressFilterStreamGenVectorClient{stream}
	return x, nil
}

type IngressFilter_StreamGenVectorClient interface {
	Send(*payload.Object_Blob) error
	Recv() (*payload.Object_Vector, error)
	grpc.ClientStream
}

type ingressFilterStreamGenVectorClient struct {
	grpc.ClientStream
}

func (x *ingressFilterStreamGenVectorClient) Send(m *payload.Object_Blob) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ingressFilterStreamGenVectorClient) Recv() (*payload.Object_Vector, error) {
	m := new(payload.Object_Vector)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *ingressFilterClient) FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/filter.ingress.v1.IngressFilter/FilterVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingressFilterClient) StreamFilterVector(ctx context.Context, opts ...grpc.CallOption) (IngressFilter_StreamFilterVectorClient, error) {
	stream, err := c.cc.NewStream(ctx, &_IngressFilter_serviceDesc.Streams[1], "/filter.ingress.v1.IngressFilter/StreamFilterVector", opts...)
	if err != nil {
		return nil, err
	}
	x := &ingressFilterStreamFilterVectorClient{stream}
	return x, nil
}

type IngressFilter_StreamFilterVectorClient interface {
	Send(*payload.Object_Vector) error
	Recv() (*payload.Object_Vector, error)
	grpc.ClientStream
}

type ingressFilterStreamFilterVectorClient struct {
	grpc.ClientStream
}

func (x *ingressFilterStreamFilterVectorClient) Send(m *payload.Object_Vector) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ingressFilterStreamFilterVectorClient) Recv() (*payload.Object_Vector, error) {
	m := new(payload.Object_Vector)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// IngressFilterServer is the server API for IngressFilter service.
type IngressFilterServer interface {
	GenVector(context.Context, *payload.Object_Blob) (*payload.Object_Vector, error)
	StreamGenVector(IngressFilter_StreamGenVectorServer) error
	FilterVector(context.Context, *payload.Object_Vector) (*payload.Object_Vector, error)
	StreamFilterVector(IngressFilter_StreamFilterVectorServer) error
}

// UnimplementedIngressFilterServer can be embedded to have forward compatible implementations.
type UnimplementedIngressFilterServer struct {
}

func (*UnimplementedIngressFilterServer) GenVector(ctx context.Context, req *payload.Object_Blob) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenVector not implemented")
}
func (*UnimplementedIngressFilterServer) StreamGenVector(srv IngressFilter_StreamGenVectorServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamGenVector not implemented")
}
func (*UnimplementedIngressFilterServer) FilterVector(ctx context.Context, req *payload.Object_Vector) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FilterVector not implemented")
}
func (*UnimplementedIngressFilterServer) StreamFilterVector(srv IngressFilter_StreamFilterVectorServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamFilterVector not implemented")
}

func RegisterIngressFilterServer(s *grpc.Server, srv IngressFilterServer) {
	s.RegisterService(&_IngressFilter_serviceDesc, srv)
}

func _IngressFilter_GenVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngressFilterServer).GenVector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.ingress.v1.IngressFilter/GenVector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngressFilterServer).GenVector(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngressFilter_StreamGenVector_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IngressFilterServer).StreamGenVector(&ingressFilterStreamGenVectorServer{stream})
}

type IngressFilter_StreamGenVectorServer interface {
	Send(*payload.Object_Vector) error
	Recv() (*payload.Object_Blob, error)
	grpc.ServerStream
}

type ingressFilterStreamGenVectorServer struct {
	grpc.ServerStream
}

func (x *ingressFilterStreamGenVectorServer) Send(m *payload.Object_Vector) error {
	return x.ServerStream.SendMsg(m)
}

func (x *ingressFilterStreamGenVectorServer) Recv() (*payload.Object_Blob, error) {
	m := new(payload.Object_Blob)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _IngressFilter_FilterVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngressFilterServer).FilterVector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.ingress.v1.IngressFilter/FilterVector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngressFilterServer).FilterVector(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngressFilter_StreamFilterVector_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IngressFilterServer).StreamFilterVector(&ingressFilterStreamFilterVectorServer{stream})
}

type IngressFilter_StreamFilterVectorServer interface {
	Send(*payload.Object_Vector) error
	Recv() (*payload.Object_Vector, error)
	grpc.ServerStream
}

type ingressFilterStreamFilterVectorServer struct {
	grpc.ServerStream
}

func (x *ingressFilterStreamFilterVectorServer) Send(m *payload.Object_Vector) error {
	return x.ServerStream.SendMsg(m)
}

func (x *ingressFilterStreamFilterVectorServer) Recv() (*payload.Object_Vector, error) {
	m := new(payload.Object_Vector)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _IngressFilter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "filter.ingress.v1.IngressFilter",
	HandlerType: (*IngressFilterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenVector",
			Handler:    _IngressFilter_GenVector_Handler,
		},
		{
			MethodName: "FilterVector",
			Handler:    _IngressFilter_FilterVector_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamGenVector",
			Handler:       _IngressFilter_StreamGenVector_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamFilterVector",
			Handler:       _IngressFilter_StreamFilterVector_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "apis/proto/v1/filter/ingress/ingress_filter.proto",
}
