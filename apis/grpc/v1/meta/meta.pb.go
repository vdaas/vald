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

package meta

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

func init() { proto.RegisterFile("apis/proto/v1/meta/meta.proto", fileDescriptor_f506bb68c7e24dcc) }
func init() {
	golang_proto.RegisterFile("apis/proto/v1/meta/meta.proto", fileDescriptor_f506bb68c7e24dcc)
}

var fileDescriptor_f506bb68c7e24dcc = []byte{
	// 424 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0xc1, 0x6a, 0xe2, 0x40,
	0x18, 0xc7, 0xc9, 0xb2, 0xab, 0x61, 0x16, 0x75, 0x0d, 0xeb, 0xee, 0x12, 0xd8, 0x1c, 0x66, 0x0f,
	0x0b, 0x81, 0x9d, 0xc1, 0xed, 0xa1, 0x77, 0xb1, 0xd8, 0x22, 0x85, 0x82, 0x60, 0xc1, 0xdb, 0x68,
	0xa6, 0x69, 0x20, 0x66, 0x42, 0x66, 0x0c, 0x78, 0xed, 0x2b, 0xf4, 0xdc, 0x77, 0xe9, 0xd1, 0x63,
	0xa1, 0x2f, 0x50, 0xb4, 0x0f, 0x52, 0x66, 0x26, 0x41, 0x4d, 0xa7, 0x17, 0xbd, 0xe8, 0x97, 0xc9,
	0xff, 0xfb, 0xe5, 0x9f, 0xef, 0xcf, 0x17, 0xf0, 0x9b, 0xa4, 0x11, 0xc7, 0x69, 0xc6, 0x04, 0xc3,
	0x79, 0x17, 0xcf, 0xa9, 0x20, 0xea, 0x07, 0xa9, 0x23, 0xa7, 0xae, 0xea, 0xbc, 0xeb, 0xfe, 0xd9,
	0xd7, 0xa5, 0x64, 0x19, 0x33, 0x12, 0x94, 0xff, 0x5a, 0xed, 0xfe, 0x0b, 0x23, 0x71, 0xbb, 0x98,
	0xa2, 0x19, 0x9b, 0xe3, 0x90, 0x85, 0x4c, 0xeb, 0xa7, 0x8b, 0x1b, 0x75, 0xa5, 0x9b, 0x65, 0x55,
	0xc8, 0x4f, 0xab, 0xf2, 0x90, 0xb1, 0x30, 0xa6, 0xea, 0x49, 0xba, 0xc4, 0x24, 0x8d, 0x30, 0x49,
	0x12, 0x26, 0x88, 0x88, 0x58, 0xc2, 0x75, 0xe3, 0xff, 0x87, 0x1a, 0xf8, 0x7c, 0x49, 0x05, 0x71,
	0xfa, 0xa0, 0x3e, 0xa0, 0x42, 0x95, 0xdf, 0x51, 0xe9, 0x25, 0xef, 0x22, 0x79, 0x82, 0x86, 0x74,
	0xe9, 0xbe, 0x3f, 0x1d, 0x93, 0x18, 0x36, 0xee, 0x9e, 0x5f, 0xef, 0x3f, 0xd5, 0xe1, 0x17, 0xf5,
	0xaa, 0xce, 0x39, 0xb0, 0x0b, 0x0a, 0x77, 0x3a, 0x26, 0x0c, 0x77, 0x3b, 0x26, 0x0e, 0x87, 0x4d,
	0x05, 0xb2, 0x61, 0x4d, 0x81, 0xb8, 0x33, 0x02, 0xcd, 0x82, 0x74, 0x91, 0xe4, 0x34, 0xe3, 0xd4,
	0x31, 0x1a, 0x70, 0x8d, 0x66, 0x61, 0x47, 0xd1, 0x5a, 0xb0, 0x81, 0x23, 0xdd, 0xad, 0xed, 0x5d,
	0x83, 0x56, 0x69, 0xaf, 0xa4, 0x9a, 0xed, 0xb8, 0x66, 0xf3, 0xf0, 0x87, 0xe2, 0x7e, 0x83, 0xcd,
	0x3d, 0x2e, 0x97, 0xd3, 0x1b, 0x15, 0xd3, 0xfb, 0x69, 0xea, 0x94, 0x4e, 0xdb, 0xbb, 0x37, 0xce,
	0xe6, 0xa9, 0x58, 0x56, 0xa7, 0x37, 0x00, 0xf6, 0xa8, 0x9c, 0xde, 0xaf, 0x0f, 0x30, 0xdc, 0xc4,
	0xa9, 0x0e, 0x6f, 0x00, 0x40, 0x9f, 0xc6, 0x54, 0xd0, 0x43, 0xf3, 0xf4, 0x0b, 0x47, 0x43, 0xf0,
	0x75, 0x0b, 0x3a, 0x30, 0x52, 0xbf, 0x74, 0x35, 0x06, 0xed, 0x2d, 0xec, 0x88, 0x54, 0xfd, 0x4a,
	0xaa, 0x13, 0xe0, 0xec, 0x98, 0x3c, 0x2a, 0x58, 0xbf, 0x12, 0x6c, 0x8f, 0xae, 0xd6, 0x9e, 0xf5,
	0xb4, 0xf6, 0xac, 0x97, 0xb5, 0x67, 0x3d, 0x6e, 0x3c, 0x6b, 0xb5, 0xf1, 0x2c, 0xe0, 0xb2, 0x2c,
	0x44, 0x79, 0x40, 0x08, 0x47, 0x39, 0x89, 0x03, 0x44, 0xd2, 0x48, 0xe2, 0x64, 0x4b, 0xcf, 0x1e,
	0x93, 0x38, 0x90, 0xe0, 0x2b, 0x6b, 0xf2, 0x77, 0x67, 0x2d, 0x95, 0x1c, 0x4b, 0x39, 0xd6, 0x6b,
	0x99, 0xa5, 0xb3, 0xf2, 0x3b, 0x31, 0xad, 0xa9, 0x6d, 0x3c, 0x79, 0x0b, 0x00, 0x00, 0xff, 0xff,
	0xdd, 0x41, 0x10, 0xcf, 0x44, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MetaClient is the client API for Meta service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MetaClient interface {
	GetMeta(ctx context.Context, in *payload.Meta_Key, opts ...grpc.CallOption) (*payload.Meta_Val, error)
	GetMetas(ctx context.Context, in *payload.Meta_Keys, opts ...grpc.CallOption) (*payload.Meta_Vals, error)
	GetMetaInverse(ctx context.Context, in *payload.Meta_Val, opts ...grpc.CallOption) (*payload.Meta_Key, error)
	GetMetasInverse(ctx context.Context, in *payload.Meta_Vals, opts ...grpc.CallOption) (*payload.Meta_Keys, error)
	SetMeta(ctx context.Context, in *payload.Meta_KeyVal, opts ...grpc.CallOption) (*payload.Empty, error)
	SetMetas(ctx context.Context, in *payload.Meta_KeyVals, opts ...grpc.CallOption) (*payload.Empty, error)
	DeleteMeta(ctx context.Context, in *payload.Meta_Key, opts ...grpc.CallOption) (*payload.Meta_Val, error)
	DeleteMetas(ctx context.Context, in *payload.Meta_Keys, opts ...grpc.CallOption) (*payload.Meta_Vals, error)
	DeleteMetaInverse(ctx context.Context, in *payload.Meta_Val, opts ...grpc.CallOption) (*payload.Meta_Key, error)
	DeleteMetasInverse(ctx context.Context, in *payload.Meta_Vals, opts ...grpc.CallOption) (*payload.Meta_Keys, error)
}

type metaClient struct {
	cc *grpc.ClientConn
}

func NewMetaClient(cc *grpc.ClientConn) MetaClient {
	return &metaClient{cc}
}

func (c *metaClient) GetMeta(ctx context.Context, in *payload.Meta_Key, opts ...grpc.CallOption) (*payload.Meta_Val, error) {
	out := new(payload.Meta_Val)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/GetMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) GetMetas(ctx context.Context, in *payload.Meta_Keys, opts ...grpc.CallOption) (*payload.Meta_Vals, error) {
	out := new(payload.Meta_Vals)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/GetMetas", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) GetMetaInverse(ctx context.Context, in *payload.Meta_Val, opts ...grpc.CallOption) (*payload.Meta_Key, error) {
	out := new(payload.Meta_Key)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/GetMetaInverse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) GetMetasInverse(ctx context.Context, in *payload.Meta_Vals, opts ...grpc.CallOption) (*payload.Meta_Keys, error) {
	out := new(payload.Meta_Keys)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/GetMetasInverse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) SetMeta(ctx context.Context, in *payload.Meta_KeyVal, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/SetMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) SetMetas(ctx context.Context, in *payload.Meta_KeyVals, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/SetMetas", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) DeleteMeta(ctx context.Context, in *payload.Meta_Key, opts ...grpc.CallOption) (*payload.Meta_Val, error) {
	out := new(payload.Meta_Val)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/DeleteMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) DeleteMetas(ctx context.Context, in *payload.Meta_Keys, opts ...grpc.CallOption) (*payload.Meta_Vals, error) {
	out := new(payload.Meta_Vals)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/DeleteMetas", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) DeleteMetaInverse(ctx context.Context, in *payload.Meta_Val, opts ...grpc.CallOption) (*payload.Meta_Key, error) {
	out := new(payload.Meta_Key)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/DeleteMetaInverse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) DeleteMetasInverse(ctx context.Context, in *payload.Meta_Vals, opts ...grpc.CallOption) (*payload.Meta_Keys, error) {
	out := new(payload.Meta_Keys)
	err := c.cc.Invoke(ctx, "/meta.v1.Meta/DeleteMetasInverse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetaServer is the server API for Meta service.
type MetaServer interface {
	GetMeta(context.Context, *payload.Meta_Key) (*payload.Meta_Val, error)
	GetMetas(context.Context, *payload.Meta_Keys) (*payload.Meta_Vals, error)
	GetMetaInverse(context.Context, *payload.Meta_Val) (*payload.Meta_Key, error)
	GetMetasInverse(context.Context, *payload.Meta_Vals) (*payload.Meta_Keys, error)
	SetMeta(context.Context, *payload.Meta_KeyVal) (*payload.Empty, error)
	SetMetas(context.Context, *payload.Meta_KeyVals) (*payload.Empty, error)
	DeleteMeta(context.Context, *payload.Meta_Key) (*payload.Meta_Val, error)
	DeleteMetas(context.Context, *payload.Meta_Keys) (*payload.Meta_Vals, error)
	DeleteMetaInverse(context.Context, *payload.Meta_Val) (*payload.Meta_Key, error)
	DeleteMetasInverse(context.Context, *payload.Meta_Vals) (*payload.Meta_Keys, error)
}

// UnimplementedMetaServer can be embedded to have forward compatible implementations.
type UnimplementedMetaServer struct {
}

func (*UnimplementedMetaServer) GetMeta(ctx context.Context, req *payload.Meta_Key) (*payload.Meta_Val, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMeta not implemented")
}
func (*UnimplementedMetaServer) GetMetas(ctx context.Context, req *payload.Meta_Keys) (*payload.Meta_Vals, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetas not implemented")
}
func (*UnimplementedMetaServer) GetMetaInverse(ctx context.Context, req *payload.Meta_Val) (*payload.Meta_Key, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetaInverse not implemented")
}
func (*UnimplementedMetaServer) GetMetasInverse(ctx context.Context, req *payload.Meta_Vals) (*payload.Meta_Keys, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetasInverse not implemented")
}
func (*UnimplementedMetaServer) SetMeta(ctx context.Context, req *payload.Meta_KeyVal) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMeta not implemented")
}
func (*UnimplementedMetaServer) SetMetas(ctx context.Context, req *payload.Meta_KeyVals) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMetas not implemented")
}
func (*UnimplementedMetaServer) DeleteMeta(ctx context.Context, req *payload.Meta_Key) (*payload.Meta_Val, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMeta not implemented")
}
func (*UnimplementedMetaServer) DeleteMetas(ctx context.Context, req *payload.Meta_Keys) (*payload.Meta_Vals, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMetas not implemented")
}
func (*UnimplementedMetaServer) DeleteMetaInverse(ctx context.Context, req *payload.Meta_Val) (*payload.Meta_Key, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMetaInverse not implemented")
}
func (*UnimplementedMetaServer) DeleteMetasInverse(ctx context.Context, req *payload.Meta_Vals) (*payload.Meta_Keys, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMetasInverse not implemented")
}

func RegisterMetaServer(s *grpc.Server, srv MetaServer) {
	s.RegisterService(&_Meta_serviceDesc, srv)
}

func _Meta_GetMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Meta_Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).GetMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/GetMeta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).GetMeta(ctx, req.(*payload.Meta_Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_GetMetas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Meta_Keys)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).GetMetas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/GetMetas",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).GetMetas(ctx, req.(*payload.Meta_Keys))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_GetMetaInverse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Meta_Val)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).GetMetaInverse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/GetMetaInverse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).GetMetaInverse(ctx, req.(*payload.Meta_Val))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_GetMetasInverse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Meta_Vals)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).GetMetasInverse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/GetMetasInverse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).GetMetasInverse(ctx, req.(*payload.Meta_Vals))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_SetMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Meta_KeyVal)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).SetMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/SetMeta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).SetMeta(ctx, req.(*payload.Meta_KeyVal))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_SetMetas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Meta_KeyVals)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).SetMetas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/SetMetas",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).SetMetas(ctx, req.(*payload.Meta_KeyVals))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_DeleteMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Meta_Key)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).DeleteMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/DeleteMeta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).DeleteMeta(ctx, req.(*payload.Meta_Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_DeleteMetas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Meta_Keys)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).DeleteMetas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/DeleteMetas",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).DeleteMetas(ctx, req.(*payload.Meta_Keys))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_DeleteMetaInverse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Meta_Val)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).DeleteMetaInverse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/DeleteMetaInverse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).DeleteMetaInverse(ctx, req.(*payload.Meta_Val))
	}
	return interceptor(ctx, in, info, handler)
}

func _Meta_DeleteMetasInverse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Meta_Vals)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetaServer).DeleteMetasInverse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.Meta/DeleteMetasInverse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).DeleteMetasInverse(ctx, req.(*payload.Meta_Vals))
	}
	return interceptor(ctx, in, info, handler)
}

var _Meta_serviceDesc = grpc.ServiceDesc{
	ServiceName: "meta.v1.Meta",
	HandlerType: (*MetaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMeta",
			Handler:    _Meta_GetMeta_Handler,
		},
		{
			MethodName: "GetMetas",
			Handler:    _Meta_GetMetas_Handler,
		},
		{
			MethodName: "GetMetaInverse",
			Handler:    _Meta_GetMetaInverse_Handler,
		},
		{
			MethodName: "GetMetasInverse",
			Handler:    _Meta_GetMetasInverse_Handler,
		},
		{
			MethodName: "SetMeta",
			Handler:    _Meta_SetMeta_Handler,
		},
		{
			MethodName: "SetMetas",
			Handler:    _Meta_SetMetas_Handler,
		},
		{
			MethodName: "DeleteMeta",
			Handler:    _Meta_DeleteMeta_Handler,
		},
		{
			MethodName: "DeleteMetas",
			Handler:    _Meta_DeleteMetas_Handler,
		},
		{
			MethodName: "DeleteMetaInverse",
			Handler:    _Meta_DeleteMetaInverse_Handler,
		},
		{
			MethodName: "DeleteMetasInverse",
			Handler:    _Meta_DeleteMetasInverse_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/proto/v1/meta/meta.proto",
}
