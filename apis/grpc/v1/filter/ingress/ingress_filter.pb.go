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

package ingress

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
	proto.RegisterFile("apis/proto/v1/filter/ingress/ingress_filter.proto", fileDescriptor_8b82e91ce4fe335b)
}

var fileDescriptor_8b82e91ce4fe335b = []byte{
	// 303 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0x31, 0x4a, 0x04, 0x31,
	0x14, 0x86, 0xcd, 0x16, 0x8a, 0x61, 0x55, 0x0c, 0x62, 0xb1, 0xc8, 0x16, 0xab, 0x85, 0x58, 0x24,
	0x8e, 0x76, 0x96, 0x0b, 0x2a, 0x0a, 0xa2, 0x28, 0x6c, 0x61, 0x23, 0x6f, 0x66, 0x62, 0x8c, 0x64,
	0xe7, 0x85, 0x4c, 0x1c, 0xd0, 0xd2, 0x2b, 0x78, 0x26, 0xc1, 0x52, 0xf0, 0x02, 0xb2, 0x78, 0x10,
	0xd9, 0x64, 0xa6, 0x18, 0x45, 0x10, 0xad, 0x1e, 0xfc, 0x79, 0xf9, 0xf8, 0xc8, 0x1f, 0x9a, 0x80,
	0xd5, 0xa5, 0xb0, 0x0e, 0x3d, 0x8a, 0x2a, 0x11, 0xd7, 0xda, 0x78, 0xe9, 0x84, 0x2e, 0x94, 0x93,
	0x65, 0xd9, 0xcc, 0xab, 0x18, 0xf3, 0xb0, 0xc6, 0x16, 0xdb, 0x69, 0x6f, 0xbd, 0x8d, 0xb0, 0x70,
	0x6f, 0x10, 0xf2, 0x66, 0xc6, 0x4b, 0xbd, 0x35, 0x85, 0xa8, 0x8c, 0x14, 0x60, 0xb5, 0x80, 0xa2,
	0x40, 0x0f, 0x5e, 0x63, 0x51, 0xc6, 0xd3, 0x9d, 0xe7, 0x0e, 0x5d, 0x38, 0x8a, 0xd4, 0x83, 0x00,
	0x65, 0x27, 0x74, 0xfe, 0x50, 0x16, 0x23, 0x99, 0x79, 0x74, 0x6c, 0x85, 0x37, 0xb0, 0xd3, 0xf4,
	0x56, 0x66, 0x9e, 0x0f, 0x0d, 0xa6, 0xbd, 0xd5, 0xaf, 0x69, 0xdc, 0x1e, 0xb0, 0xc7, 0xb7, 0x8f,
	0xa7, 0x4e, 0x77, 0x30, 0x27, 0x30, 0xe4, 0x7b, 0x64, 0x8b, 0xed, 0xd3, 0xa5, 0x0b, 0xef, 0x24,
	0x8c, 0xff, 0x0a, 0x9d, 0xd9, 0x24, 0xdb, 0x84, 0x9d, 0xd3, 0x6e, 0xf4, 0xab, 0x19, 0x3f, 0x6c,
	0xff, 0x42, 0xad, 0x0a, 0xc1, 0x54, 0xed, 0x98, 0xb2, 0xa8, 0xf6, 0x2f, 0x72, 0xf0, 0x1b, 0x3e,
	0xbc, 0x4c, 0xfa, 0xe4, 0x75, 0xd2, 0x27, 0xef, 0x93, 0x3e, 0xa1, 0x1b, 0xe8, 0x14, 0xaf, 0x72,
	0x80, 0x92, 0x57, 0x60, 0x72, 0x0e, 0x56, 0xf3, 0x2a, 0xe1, 0x75, 0x99, 0x75, 0x8b, 0xc3, 0xe5,
	0x11, 0x98, 0xbc, 0xf5, 0xf8, 0x67, 0xe4, 0x32, 0x51, 0xda, 0xdf, 0xdc, 0xa5, 0x3c, 0xc3, 0xb1,
	0x08, 0x04, 0x31, 0x25, 0x88, 0xd0, 0xb4, 0x72, 0x36, 0xfb, 0xfe, 0x57, 0xd2, 0xd9, 0x50, 0xe5,
	0xee, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0x14, 0x8e, 0x33, 0xa2, 0x52, 0x02, 0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/ingress_filter.IngressFilter/GenVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingressFilterClient) StreamGenVector(ctx context.Context, opts ...grpc.CallOption) (IngressFilter_StreamGenVectorClient, error) {
	stream, err := c.cc.NewStream(ctx, &_IngressFilter_serviceDesc.Streams[0], "/ingress_filter.IngressFilter/StreamGenVector", opts...)
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
	err := c.cc.Invoke(ctx, "/ingress_filter.IngressFilter/FilterVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingressFilterClient) StreamFilterVector(ctx context.Context, opts ...grpc.CallOption) (IngressFilter_StreamFilterVectorClient, error) {
	stream, err := c.cc.NewStream(ctx, &_IngressFilter_serviceDesc.Streams[1], "/ingress_filter.IngressFilter/StreamFilterVector", opts...)
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
		FullMethod: "/ingress_filter.IngressFilter/GenVector",
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
		FullMethod: "/ingress_filter.IngressFilter/FilterVector",
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
	ServiceName: "ingress_filter.IngressFilter",
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
