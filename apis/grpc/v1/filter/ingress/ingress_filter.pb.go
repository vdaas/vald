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
var (
	_ = proto.Marshal
	_ = golang_proto.Marshal
	_ = fmt.Errorf
	_ = math.Inf
)

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
	// 311 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0x3f, 0x4b, 0xc4, 0x30,
	0x14, 0x27, 0x0e, 0x07, 0x16, 0x97, 0xeb, 0xa0, 0x58, 0xa1, 0xe0, 0xe9, 0x74, 0x60, 0x42, 0x75,
	0x10, 0x6e, 0xec, 0xa0, 0x38, 0xe9, 0x74, 0x83, 0x8b, 0x24, 0x6d, 0x2e, 0x46, 0x62, 0x5f, 0x48,
	0x73, 0x01, 0xd7, 0xfb, 0x0a, 0x7e, 0x21, 0xc7, 0x1b, 0x05, 0x71, 0x97, 0x9e, 0x1f, 0x44, 0x9a,
	0xb4, 0xe0, 0x9d, 0xba, 0x38, 0xbd, 0xf7, 0xfa, 0x7e, 0x7f, 0x1e, 0xbf, 0x26, 0xca, 0xa8, 0x96,
	0x35, 0xd1, 0x06, 0x2c, 0x10, 0x97, 0x91, 0x99, 0x54, 0x96, 0x1b, 0x22, 0x2b, 0x61, 0x78, 0x5d,
	0xf7, 0xf5, 0x2e, 0x7c, 0xc6, 0x1e, 0x16, 0x0f, 0xbb, 0xa9, 0x5b, 0x62, 0x97, 0x25, 0x47, 0xeb,
	0x2a, 0x9a, 0x3e, 0x29, 0xa0, 0x65, 0x5f, 0x03, 0x2f, 0x39, 0x11, 0xd2, 0xde, 0xcf, 0x19, 0x2e,
	0xe0, 0x91, 0x08, 0x10, 0x10, 0xf0, 0x6c, 0x3e, 0xf3, 0x53, 0x20, 0xb7, 0x5d, 0x07, 0x3f, 0xdf,
	0x84, 0x0b, 0x00, 0xa1, 0xb8, 0x77, 0x0a, 0x2d, 0xa1, 0x5a, 0x12, 0x5a, 0x55, 0x60, 0xa9, 0x95,
	0x50, 0xd5, 0x81, 0x78, 0xfa, 0x8e, 0xa2, 0xc1, 0x85, 0x3f, 0x31, 0x66, 0xd1, 0xf6, 0x25, 0xaf,
	0xa6, 0xbc, 0xb0, 0x60, 0xe2, 0x3d, 0xdc, 0xdf, 0xe3, 0x32, 0x7c, 0xcd, 0x1e, 0x78, 0x61, 0x71,
	0xae, 0x80, 0x25, 0xfb, 0xbf, 0x2c, 0x02, 0x67, 0x74, 0xb8, 0x78, 0xfb, 0x7c, 0xde, 0x3a, 0x18,
	0xed, 0x6e, 0x26, 0x03, 0x1e, 0x36, 0x41, 0xe3, 0x58, 0x44, 0x3b, 0xc1, 0xad, 0xb3, 0xf9, 0x5b,
	0xed, 0x5f, 0x46, 0xce, 0xef, 0x27, 0x68, 0x9c, 0x2f, 0xd0, 0xb2, 0x49, 0xd1, 0x6b, 0x93, 0xa2,
	0x8f, 0x26, 0x45, 0x2f, 0xab, 0x14, 0x2d, 0x57, 0x29, 0x8a, 0x8e, 0xc1, 0x08, 0xec, 0x4a, 0x4a,
	0x6b, 0xec, 0xa8, 0x2a, 0x31, 0xd5, 0xb2, 0x95, 0x5f, 0xff, 0x49, 0xf9, 0x70, 0x4a, 0x55, 0x79,
	0x15, 0x86, 0x70, 0xee, 0x0d, 0xba, 0xcd, 0xbe, 0x45, 0xec, 0x15, 0x48, 0xab, 0x40, 0x42, 0xc4,
	0x46, 0x17, 0x3f, 0x5f, 0x04, 0x1b, 0xf8, 0x8c, 0xcf, 0xbe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x57,
	0x46, 0xab, 0x11, 0x38, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ context.Context
	_ grpc.ClientConn
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FilterClient is the client API for Filter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FilterClient interface {
	GenVector(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Vector, error)
	FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Object_Vector, error)
}

type filterClient struct {
	cc *grpc.ClientConn
}

func NewFilterClient(cc *grpc.ClientConn) FilterClient {
	return &filterClient{cc}
}

func (c *filterClient) GenVector(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/filter.ingress.v1.Filter/GenVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/filter.ingress.v1.Filter/FilterVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilterServer is the server API for Filter service.
type FilterServer interface {
	GenVector(context.Context, *payload.Object_Blob) (*payload.Object_Vector, error)
	FilterVector(context.Context, *payload.Object_Vector) (*payload.Object_Vector, error)
}

// UnimplementedFilterServer can be embedded to have forward compatible implementations.
type UnimplementedFilterServer struct{}

func (*UnimplementedFilterServer) GenVector(ctx context.Context, req *payload.Object_Blob) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenVector not implemented")
}

func (*UnimplementedFilterServer) FilterVector(ctx context.Context, req *payload.Object_Vector) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FilterVector not implemented")
}

func RegisterFilterServer(s *grpc.Server, srv FilterServer) {
	s.RegisterService(&_Filter_serviceDesc, srv)
}

func _Filter_GenVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).GenVector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.ingress.v1.Filter/GenVector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).GenVector(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_FilterVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).FilterVector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.ingress.v1.Filter/FilterVector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).FilterVector(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

var _Filter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "filter.ingress.v1.Filter",
	HandlerType: (*FilterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenVector",
			Handler:    _Filter_GenVector_Handler,
		},
		{
			MethodName: "FilterVector",
			Handler:    _Filter_FilterVector_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/proto/v1/filter/ingress/ingress_filter.proto",
}
