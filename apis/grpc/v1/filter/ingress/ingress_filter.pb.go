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

	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	golang_proto "github.com/golang/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	codes "github.com/vdaas/vald/internal/net/grpc/codes"
	status "github.com/vdaas/vald/internal/net/grpc/status"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
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
func init() {
	golang_proto.RegisterFile("apis/proto/v1/filter/ingress/ingress_filter.proto", fileDescriptor_8b82e91ce4fe335b)
}

var fileDescriptor_8b82e91ce4fe335b = []byte{
	// 346 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xc1, 0x4a, 0x33, 0x31,
	0x14, 0x85, 0xff, 0x74, 0xf1, 0x8b, 0xa1, 0x22, 0xcd, 0x46, 0xec, 0x62, 0x84, 0xea, 0x42, 0x04,
	0x13, 0x47, 0x17, 0x82, 0xcb, 0x2e, 0x14, 0x57, 0x2a, 0x85, 0x22, 0x6e, 0x24, 0x33, 0x93, 0xc6,
	0xc8, 0x74, 0xee, 0x90, 0x49, 0x03, 0x6e, 0xfb, 0x0a, 0xbe, 0x90, 0xcb, 0x2e, 0x05, 0x5f, 0x40,
	0x5a, 0x37, 0xbe, 0x85, 0x4c, 0x32, 0x05, 0x47, 0x2b, 0xcc, 0x2a, 0xc9, 0xbd, 0xf7, 0x9c, 0xfb,
	0xc1, 0x09, 0x0e, 0x79, 0xae, 0x0a, 0x96, 0x6b, 0x30, 0xc0, 0x6c, 0xc8, 0x46, 0x2a, 0x35, 0x42,
	0x33, 0x95, 0x49, 0x2d, 0x8a, 0x62, 0x79, 0xde, 0xfb, 0x32, 0x75, 0x63, 0xa4, 0x53, 0xbd, 0xaa,
	0x26, 0xb5, 0x61, 0x77, 0xb7, 0xee, 0x92, 0xf3, 0xa7, 0x14, 0x78, 0xb2, 0x3c, 0xbd, 0xae, 0x7b,
	0x28, 0x95, 0x79, 0x98, 0x44, 0x34, 0x86, 0x31, 0x93, 0x20, 0xc1, 0xcf, 0x47, 0x93, 0x91, 0x7b,
	0x79, 0x71, 0x79, 0xab, 0xc6, 0x4f, 0x7f, 0x8e, 0x4b, 0x00, 0x99, 0x0a, 0xb7, 0xc9, 0x5f, 0x19,
	0xcf, 0x15, 0xe3, 0x59, 0x06, 0x86, 0x1b, 0x05, 0x59, 0xe1, 0x85, 0xc7, 0x9f, 0x2d, 0xbc, 0x71,
	0xe9, 0xd9, 0xce, 0x1d, 0x29, 0x19, 0xe0, 0xf5, 0x0b, 0x91, 0x0d, 0x45, 0x6c, 0x40, 0x93, 0x2d,
	0xba, 0xc4, 0xb2, 0x21, 0xbd, 0x8a, 0x1e, 0x45, 0x6c, 0x68, 0x3f, 0x85, 0xa8, 0xbb, 0xbd, 0xa2,
	0xe1, 0x35, 0x3d, 0x32, 0x7d, 0xfb, 0x78, 0x6e, 0xb5, 0x7b, 0x6b, 0x0c, 0x5c, 0xfd, 0x0c, 0x1d,
	0x90, 0x1b, 0xbc, 0x39, 0x30, 0x5a, 0xf0, 0x71, 0x03, 0xeb, 0x9d, 0x15, 0x0d, 0x2f, 0xae, 0x16,
	0xfc, 0xdb, 0x47, 0x47, 0x88, 0xdc, 0xe2, 0xb6, 0x27, 0xae, 0xfc, 0xfe, 0x26, 0x6a, 0x06, 0x6b,
	0x5d, 0xa1, 0x84, 0x1d, 0x62, 0xe2, 0xf7, 0x35, 0xf5, 0x6f, 0x46, 0xdc, 0x9f, 0xa2, 0xd9, 0x3c,
	0x40, 0xaf, 0xf3, 0x00, 0xbd, 0xcf, 0x03, 0xf4, 0xb2, 0x08, 0xd0, 0x6c, 0x11, 0x20, 0xbc, 0x07,
	0x5a, 0x52, 0x9b, 0x70, 0x5e, 0x50, 0xcb, 0xd3, 0x84, 0xf2, 0x5c, 0x95, 0x3e, 0xf5, 0x8f, 0xd3,
	0xef, 0x0c, 0x79, 0x9a, 0xd4, 0x92, 0xba, 0x46, 0x77, 0xe1, 0xb7, 0xd8, 0x9d, 0x03, 0x2b, 0x1d,
	0x98, 0x8f, 0x5d, 0xe7, 0xf1, 0xef, 0x5f, 0x1a, 0xfd, 0x77, 0xb9, 0x9f, 0x7c, 0x05, 0x00, 0x00,
	0xff, 0xff, 0x72, 0x3c, 0xa6, 0x24, 0xcc, 0x02, 0x00, 0x00,
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
