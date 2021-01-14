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
	// 312 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xb1, 0x4a, 0x04, 0x31,
	0x10, 0x86, 0xcd, 0x15, 0x8a, 0xe1, 0x44, 0x2e, 0x8d, 0x78, 0xc8, 0x09, 0xa7, 0x85, 0x58, 0x24,
	0xae, 0x76, 0x96, 0x57, 0x28, 0x56, 0x2a, 0x07, 0x87, 0xd8, 0xc8, 0xec, 0x6e, 0x5c, 0x23, 0xb9,
	0x9d, 0x90, 0x8d, 0x01, 0x2d, 0x7d, 0x05, 0x5f, 0xca, 0x52, 0xf0, 0x05, 0xe4, 0xb0, 0xf1, 0x2d,
	0xe4, 0x92, 0xbd, 0x62, 0xf1, 0x84, 0xad, 0x86, 0xcc, 0x9f, 0xff, 0x9b, 0x1f, 0x66, 0x68, 0x02,
	0x46, 0x55, 0xc2, 0x58, 0x74, 0x28, 0x7c, 0x22, 0xee, 0x95, 0x76, 0xd2, 0x0a, 0x55, 0x16, 0x56,
	0x56, 0xd5, 0xa2, 0xde, 0xc5, 0x36, 0x0f, 0xdf, 0x58, 0xaf, 0x7e, 0xd5, 0x22, 0xf7, 0x49, 0x7f,
	0xaf, 0x49, 0x31, 0xf0, 0xac, 0x11, 0xf2, 0x45, 0x8d, 0xbe, 0xfe, 0x4e, 0x81, 0x58, 0x68, 0x29,
	0xc0, 0x28, 0x01, 0x65, 0x89, 0x0e, 0x9c, 0xc2, 0xb2, 0x8a, 0xea, 0xf1, 0x4f, 0x87, 0x6e, 0x5c,
	0x44, 0xe2, 0x59, 0xe0, 0xb3, 0x31, 0x5d, 0x3f, 0x97, 0xe5, 0x44, 0x66, 0x0e, 0x2d, 0xdb, 0xe2,
	0x0b, 0x98, 0x4f, 0xf8, 0x65, 0xfa, 0x28, 0x33, 0xc7, 0x47, 0x1a, 0xd3, 0xfe, 0xf6, 0x12, 0x21,
	0x7a, 0x86, 0xec, 0xf5, 0xf3, 0xfb, 0xad, 0xd3, 0x1d, 0xae, 0x09, 0x0c, 0xfd, 0x53, 0x72, 0xc8,
	0xae, 0xe9, 0xe6, 0xd8, 0x59, 0x09, 0xd3, 0x16, 0xe8, 0xdd, 0x25, 0x42, 0x34, 0xd7, 0x03, 0x56,
	0x0e, 0xc8, 0x11, 0x61, 0x37, 0xb4, 0x1b, 0x13, 0xd7, 0xbc, 0xff, 0x13, 0xb5, 0x0b, 0xeb, 0x43,
	0x63, 0x1e, 0x76, 0x42, 0x59, 0x9c, 0xd7, 0x96, 0xdf, 0x2e, 0xf1, 0xe8, 0xe5, 0x7d, 0x36, 0x20,
	0x1f, 0xb3, 0x01, 0xf9, 0x9a, 0x0d, 0x08, 0xdd, 0x47, 0x5b, 0x70, 0x9f, 0x03, 0x54, 0xdc, 0x83,
	0xce, 0x39, 0x18, 0x35, 0xb7, 0x37, 0xb7, 0x3c, 0xea, 0x4d, 0x40, 0xe7, 0x8d, 0x05, 0x5d, 0x91,
	0xdb, 0xa4, 0x50, 0xee, 0xe1, 0x29, 0xe5, 0x19, 0x4e, 0x45, 0x20, 0x88, 0x39, 0x41, 0x84, 0x6b,
	0x28, 0xac, 0xc9, 0xfe, 0x9e, 0x54, 0xba, 0x1a, 0xd6, 0x7d, 0xf2, 0x1b, 0x00, 0x00, 0xff, 0xff,
	0xf7, 0x7d, 0x5c, 0x94, 0x79, 0x02, 0x00, 0x00,
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
	Recv() (*payload.Object_StreamVector, error)
	grpc.ClientStream
}

type ingressFilterStreamGenVectorClient struct {
	grpc.ClientStream
}

func (x *ingressFilterStreamGenVectorClient) Send(m *payload.Object_Blob) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ingressFilterStreamGenVectorClient) Recv() (*payload.Object_StreamVector, error) {
	m := new(payload.Object_StreamVector)
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
	Recv() (*payload.Object_StreamVector, error)
	grpc.ClientStream
}

type ingressFilterStreamFilterVectorClient struct {
	grpc.ClientStream
}

func (x *ingressFilterStreamFilterVectorClient) Send(m *payload.Object_Vector) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ingressFilterStreamFilterVectorClient) Recv() (*payload.Object_StreamVector, error) {
	m := new(payload.Object_StreamVector)
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
	Send(*payload.Object_StreamVector) error
	Recv() (*payload.Object_Blob, error)
	grpc.ServerStream
}

type ingressFilterStreamGenVectorServer struct {
	grpc.ServerStream
}

func (x *ingressFilterStreamGenVectorServer) Send(m *payload.Object_StreamVector) error {
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
	Send(*payload.Object_StreamVector) error
	Recv() (*payload.Object_Vector, error)
	grpc.ServerStream
}

type ingressFilterStreamFilterVectorServer struct {
	grpc.ServerStream
}

func (x *ingressFilterStreamFilterVectorServer) Send(m *payload.Object_StreamVector) error {
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
