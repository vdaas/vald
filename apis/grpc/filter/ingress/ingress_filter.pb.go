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

func init() { proto.RegisterFile("ingress/ingress_filter.proto", fileDescriptor_8f5342c46835d3ee) }

var fileDescriptor_8f5342c46835d3ee = []byte{
	// 299 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xc9, 0xcc, 0x4b, 0x2f,
	0x4a, 0x2d, 0x2e, 0xd6, 0x87, 0xd2, 0xf1, 0x69, 0x99, 0x39, 0x25, 0xa9, 0x45, 0x7a, 0x05, 0x45,
	0xf9, 0x25, 0xf9, 0x42, 0x7c, 0xa8, 0xa2, 0x52, 0xbc, 0x05, 0x89, 0x95, 0x39, 0xf9, 0x89, 0x29,
	0x10, 0x69, 0x29, 0x99, 0xf4, 0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0xfd, 0xc4, 0x82, 0x4c, 0xfd, 0xc4,
	0xbc, 0xbc, 0xfc, 0x92, 0xc4, 0x92, 0xcc, 0xfc, 0xbc, 0x62, 0xa8, 0x2c, 0x4f, 0x41, 0x92, 0x7e,
	0x7a, 0x61, 0x0e, 0x84, 0x67, 0x74, 0x92, 0x89, 0x8b, 0xd7, 0x13, 0x62, 0x9a, 0x1b, 0xd8, 0x30,
	0x21, 0x5f, 0x2e, 0x4e, 0xf7, 0xd4, 0xbc, 0xb0, 0xd4, 0xe4, 0x92, 0xfc, 0x22, 0x21, 0x11, 0x3d,
	0x98, 0xd1, 0xfe, 0x49, 0x59, 0xa9, 0xc9, 0x25, 0x7a, 0x4e, 0x39, 0xf9, 0x49, 0x52, 0x62, 0xe8,
	0xa2, 0x10, 0xd5, 0x4a, 0x42, 0x4d, 0x97, 0x9f, 0x4c, 0x66, 0xe2, 0x51, 0x62, 0xd7, 0xcf, 0x07,
	0x8b, 0x5b, 0x31, 0x6a, 0x09, 0xb9, 0x72, 0xf1, 0x07, 0x97, 0x14, 0xa5, 0x26, 0xe6, 0x92, 0x6b,
	0x28, 0x83, 0x06, 0xa3, 0x01, 0xa3, 0x50, 0x18, 0x17, 0x0f, 0xc4, 0x7d, 0x50, 0x33, 0x70, 0xa8,
	0xc6, 0x69, 0x8a, 0xd8, 0x86, 0x07, 0xf2, 0x8c, 0x70, 0xe7, 0x95, 0x81, 0x05, 0x41, 0xce, 0xf3,
	0xe2, 0x12, 0x82, 0x38, 0x8f, 0x22, 0xd3, 0xc1, 0x6e, 0x74, 0x2a, 0x38, 0xf1, 0x48, 0x8e, 0xf1,
	0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0xb9, 0x64, 0xf3, 0x8b, 0xd2, 0xf5, 0xca, 0x52,
	0x12, 0x13, 0x8b, 0xf5, 0xca, 0x12, 0x73, 0x52, 0xf4, 0xa0, 0x31, 0x08, 0x8d, 0x3a, 0x27, 0xc1,
	0xb0, 0xc4, 0x9c, 0x14, 0x94, 0x90, 0x0f, 0x60, 0x8c, 0xd2, 0x4b, 0xcf, 0x2c, 0xc9, 0x28, 0x4d,
	0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x07, 0x6b, 0xd5, 0x07, 0x69, 0x05, 0x45, 0x64, 0xb1, 0x7e, 0x7a,
	0x51, 0x41, 0xb2, 0x3e, 0xc4, 0x10, 0x58, 0xaa, 0x48, 0x62, 0x03, 0x47, 0xa2, 0x31, 0x20, 0x00,
	0x00, 0xff, 0xff, 0xed, 0xb1, 0xae, 0xad, 0x2f, 0x02, 0x00, 0x00,
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
	Metadata: "ingress/ingress_filter.proto",
}
