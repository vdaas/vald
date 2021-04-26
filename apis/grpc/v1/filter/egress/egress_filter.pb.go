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

package egress

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
	proto.RegisterFile("apis/proto/v1/filter/egress/egress_filter.proto", fileDescriptor_7f3e67472eb32d70)
}
func init() {
	golang_proto.RegisterFile("apis/proto/v1/filter/egress/egress_filter.proto", fileDescriptor_7f3e67472eb32d70)
}

var fileDescriptor_7f3e67472eb32d70 = []byte{
	// 321 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0x4f, 0x4b, 0xc3, 0x30,
	0x1c, 0x25, 0x1e, 0x76, 0x28, 0x22, 0xa3, 0x20, 0x62, 0x95, 0x22, 0xdd, 0x6d, 0x60, 0xe2, 0xf4,
	0x20, 0xec, 0x38, 0xd4, 0xab, 0x9e, 0x76, 0xf0, 0x22, 0x69, 0x9b, 0xc5, 0x48, 0xed, 0xaf, 0x24,
	0x59, 0xc0, 0x93, 0xe0, 0x57, 0xf0, 0x0b, 0x79, 0xdc, 0x51, 0xf0, 0x0b, 0x8c, 0xcd, 0x0f, 0x22,
	0xcd, 0xaf, 0x05, 0x57, 0xff, 0x9c, 0xf2, 0xd2, 0xf7, 0xde, 0xef, 0x3d, 0x7e, 0x4d, 0xc0, 0x78,
	0xa5, 0x0c, 0xab, 0x34, 0x58, 0x60, 0x6e, 0xc4, 0x66, 0xaa, 0xb0, 0x42, 0x33, 0x21, 0xb5, 0x30,
	0xa6, 0x39, 0xee, 0xf0, 0x23, 0xf5, 0xa2, 0xb0, 0xdf, 0xdc, 0x90, 0xa3, 0x6e, 0x14, 0x0d, 0x36,
	0x47, 0x54, 0xfc, 0xa9, 0x00, 0x9e, 0xb7, 0x27, 0xda, 0xa2, 0x63, 0xa9, 0xec, 0xfd, 0x3c, 0xa5,
	0x19, 0x3c, 0x32, 0x09, 0x12, 0x50, 0x9f, 0xce, 0x67, 0xfe, 0x86, 0xe6, 0x1a, 0x35, 0xf2, 0xf3,
	0xae, 0x5c, 0x02, 0xc8, 0x42, 0xf8, 0x24, 0x84, 0x75, 0x71, 0xc6, 0xcb, 0x12, 0x2c, 0xb7, 0x0a,
	0x4a, 0x83, 0xc6, 0xd3, 0x25, 0x09, 0x7a, 0x57, 0xbe, 0x61, 0x58, 0x06, 0x3b, 0x88, 0x2e, 0x94,
	0xb1, 0xbc, 0xcc, 0x44, 0x78, 0x40, 0xdb, 0x52, 0x6e, 0x44, 0xaf, 0xd3, 0x07, 0x91, 0x59, 0xda,
	0x92, 0xd1, 0x7f, 0x64, 0x92, 0xbc, 0x7c, 0x7c, 0xbe, 0x6e, 0x1d, 0x26, 0x7b, 0x9d, 0x0d, 0xe5,
	0x8d, 0x60, 0x4c, 0x86, 0xe1, 0x2c, 0xd8, 0xc6, 0xbc, 0xa9, 0xc8, 0x2c, 0xe8, 0x70, 0xff, 0x97,
	0x81, 0x48, 0x45, 0x7f, 0x53, 0xc9, 0x91, 0x4f, 0x8a, 0x92, 0xdd, 0x4e, 0x92, 0xf3, 0xf4, 0x98,
	0x0c, 0x27, 0xcf, 0x8b, 0x55, 0x4c, 0xde, 0x57, 0x31, 0x59, 0xae, 0x62, 0xf2, 0xb6, 0x8e, 0xc9,
	0x62, 0x1d, 0x93, 0x60, 0x00, 0x5a, 0x52, 0x97, 0x73, 0x6e, 0xa8, 0xe3, 0x45, 0x4e, 0x79, 0xa5,
	0xea, 0xe1, 0x1b, 0x7f, 0x6b, 0xd2, 0x9f, 0xf2, 0x22, 0xbf, 0xf4, 0x18, 0xab, 0xde, 0x90, 0xdb,
	0x93, 0x6f, 0xab, 0xf6, 0x7e, 0x56, 0xfb, 0xf1, 0x5d, 0x48, 0x5d, 0x65, 0x3f, 0x9e, 0x45, 0xda,
	0xf3, 0xab, 0x3e, 0xfb, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x3d, 0x91, 0xa0, 0x48, 0x3c, 0x02, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FilterClient is the client API for Filter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FilterClient interface {
	FilterDistance(ctx context.Context, in *payload.Object_Distance, opts ...grpc.CallOption) (*payload.Object_Distance, error)
	FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Object_Vector, error)
}

type filterClient struct {
	cc *grpc.ClientConn
}

func NewFilterClient(cc *grpc.ClientConn) FilterClient {
	return &filterClient{cc}
}

func (c *filterClient) FilterDistance(ctx context.Context, in *payload.Object_Distance, opts ...grpc.CallOption) (*payload.Object_Distance, error) {
	out := new(payload.Object_Distance)
	err := c.cc.Invoke(ctx, "/filter.egress.v1.Filter/FilterDistance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) FilterVector(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/filter.egress.v1.Filter/FilterVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilterServer is the server API for Filter service.
type FilterServer interface {
	FilterDistance(context.Context, *payload.Object_Distance) (*payload.Object_Distance, error)
	FilterVector(context.Context, *payload.Object_Vector) (*payload.Object_Vector, error)
}

// UnimplementedFilterServer can be embedded to have forward compatible implementations.
type UnimplementedFilterServer struct {
}

func (*UnimplementedFilterServer) FilterDistance(ctx context.Context, req *payload.Object_Distance) (*payload.Object_Distance, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FilterDistance not implemented")
}
func (*UnimplementedFilterServer) FilterVector(ctx context.Context, req *payload.Object_Vector) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FilterVector not implemented")
}

func RegisterFilterServer(s *grpc.Server, srv FilterServer) {
	s.RegisterService(&_Filter_serviceDesc, srv)
}

func _Filter_FilterDistance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Distance)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).FilterDistance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.egress.v1.Filter/FilterDistance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).FilterDistance(ctx, req.(*payload.Object_Distance))
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
		FullMethod: "/filter.egress.v1.Filter/FilterVector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).FilterVector(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

var _Filter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "filter.egress.v1.Filter",
	HandlerType: (*FilterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FilterDistance",
			Handler:    _Filter_FilterDistance_Handler,
		},
		{
			MethodName: "FilterVector",
			Handler:    _Filter_FilterVector_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/proto/v1/filter/egress/egress_filter.proto",
}
