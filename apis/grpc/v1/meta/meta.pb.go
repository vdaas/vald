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

package meta

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

func init() { proto.RegisterFile("apis/proto/v1/meta/meta.proto", fileDescriptor_f506bb68c7e24dcc) }

var fileDescriptor_f506bb68c7e24dcc = []byte{
	// 391 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0x5f, 0x4b, 0xeb, 0x30,
	0x18, 0x87, 0xe9, 0x61, 0xff, 0xc8, 0x39, 0xeb, 0xce, 0xc2, 0xd9, 0xb9, 0x28, 0xba, 0x8b, 0xdc,
	0x08, 0xbd, 0x48, 0x98, 0x5e, 0x79, 0x23, 0x6c, 0x38, 0x44, 0x64, 0x30, 0x99, 0x08, 0x7a, 0x23,
	0xd9, 0x16, 0x6a, 0xa1, 0x6d, 0x4a, 0x13, 0x0b, 0xbb, 0xf5, 0x2b, 0xf8, 0xa5, 0xbc, 0x14, 0xfc,
	0x02, 0x3a, 0xfc, 0x20, 0x92, 0xb4, 0x71, 0xae, 0x1d, 0xc2, 0x76, 0x93, 0xa6, 0x69, 0xde, 0xa7,
	0x0f, 0xef, 0x8f, 0x17, 0xec, 0xd3, 0xd8, 0x17, 0x24, 0x4e, 0xb8, 0xe4, 0x24, 0xed, 0x91, 0x90,
	0x49, 0xaa, 0x17, 0xac, 0x8f, 0xe0, 0x1f, 0xb5, 0xbf, 0x0b, 0x69, 0x44, 0x3d, 0x96, 0x38, 0xc7,
	0x9e, 0x2f, 0xef, 0x1f, 0xa6, 0x78, 0xc6, 0x43, 0x92, 0xce, 0x29, 0x15, 0x24, 0xa5, 0xc1, 0x9c,
	0xac, 0x23, 0x62, 0xba, 0x08, 0x38, 0x9d, 0x9b, 0x67, 0x06, 0x72, 0xf6, 0x3c, 0xce, 0xbd, 0x80,
	0xa9, 0xbb, 0x84, 0x46, 0x11, 0x97, 0x54, 0xfa, 0x3c, 0x12, 0xd9, 0xd7, 0xc3, 0xf7, 0x2a, 0xa8,
	0x8c, 0x98, 0xa4, 0xf0, 0x04, 0xd4, 0xcf, 0x98, 0xd4, 0xdb, 0x36, 0x36, 0x04, 0xf5, 0x8a, 0x2f,
	0xd8, 0xc2, 0x29, 0x1c, 0x5d, 0xd3, 0x00, 0x35, 0x1f, 0x5f, 0x3f, 0x9e, 0x7e, 0xd5, 0x51, 0x55,
	0x5b, 0xc3, 0x01, 0x68, 0xe4, 0xf5, 0x02, 0xc2, 0x12, 0x40, 0x38, 0xb0, 0x44, 0x10, 0xc8, 0xd6,
	0x88, 0x06, 0xaa, 0x69, 0x84, 0x80, 0x23, 0x60, 0xe7, 0x8c, 0xf3, 0x28, 0x65, 0x89, 0x60, 0xb0,
	0xfc, 0x5f, 0xa7, 0x6c, 0x87, 0x3a, 0x9a, 0xd3, 0x42, 0x4d, 0xe2, 0x67, 0x75, 0x99, 0xd2, 0x25,
	0x68, 0x19, 0x25, 0xc3, 0xdb, 0x60, 0xe1, 0x6c, 0xb0, 0x45, 0xff, 0x35, 0xf1, 0x2f, 0xb2, 0xd7,
	0x88, 0x42, 0x75, 0x69, 0x92, 0x77, 0xe9, 0x5f, 0xa9, 0x4c, 0xd9, 0xd9, 0x5f, 0xa7, 0xc3, 0x30,
	0x96, 0x8b, 0x62, 0x97, 0xfa, 0xa0, 0x31, 0x31, 0x5d, 0xea, 0x6c, 0x02, 0x88, 0x12, 0xa1, 0xd8,
	0xa4, 0x3e, 0x00, 0xa7, 0x2c, 0x60, 0x92, 0x6d, 0x9f, 0x95, 0x9b, 0x5b, 0x0c, 0xc1, 0xef, 0x15,
	0x62, 0xeb, 0xb8, 0x5c, 0x63, 0x32, 0x06, 0xed, 0x15, 0x66, 0xa7, 0xc4, 0xdc, 0x42, 0x62, 0x57,
	0x00, 0x7e, 0x13, 0xdb, 0x31, 0x34, 0xb7, 0x10, 0xda, 0xe0, 0xe6, 0x79, 0xd9, 0xb5, 0x5e, 0x96,
	0x5d, 0xeb, 0x6d, 0xd9, 0xb5, 0x80, 0xc3, 0x13, 0x0f, 0xeb, 0x19, 0xc2, 0x6a, 0x86, 0x30, 0x8d,
	0x7d, 0x9c, 0xf6, 0xb0, 0xba, 0x3a, 0xa8, 0xa8, 0x75, 0x6c, 0xdd, 0x1e, 0xfc, 0x30, 0x6e, 0x5e,
	0x12, 0xcf, 0xcc, 0xc0, 0x4e, 0x6b, 0x7a, 0x8a, 0x8e, 0x3e, 0x03, 0x00, 0x00, 0xff, 0xff, 0xf3,
	0xab, 0x55, 0xa6, 0xcd, 0x03, 0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/meta_manager.Meta/GetMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) GetMetas(ctx context.Context, in *payload.Meta_Keys, opts ...grpc.CallOption) (*payload.Meta_Vals, error) {
	out := new(payload.Meta_Vals)
	err := c.cc.Invoke(ctx, "/meta_manager.Meta/GetMetas", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) GetMetaInverse(ctx context.Context, in *payload.Meta_Val, opts ...grpc.CallOption) (*payload.Meta_Key, error) {
	out := new(payload.Meta_Key)
	err := c.cc.Invoke(ctx, "/meta_manager.Meta/GetMetaInverse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) GetMetasInverse(ctx context.Context, in *payload.Meta_Vals, opts ...grpc.CallOption) (*payload.Meta_Keys, error) {
	out := new(payload.Meta_Keys)
	err := c.cc.Invoke(ctx, "/meta_manager.Meta/GetMetasInverse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) SetMeta(ctx context.Context, in *payload.Meta_KeyVal, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/meta_manager.Meta/SetMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) SetMetas(ctx context.Context, in *payload.Meta_KeyVals, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/meta_manager.Meta/SetMetas", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) DeleteMeta(ctx context.Context, in *payload.Meta_Key, opts ...grpc.CallOption) (*payload.Meta_Val, error) {
	out := new(payload.Meta_Val)
	err := c.cc.Invoke(ctx, "/meta_manager.Meta/DeleteMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) DeleteMetas(ctx context.Context, in *payload.Meta_Keys, opts ...grpc.CallOption) (*payload.Meta_Vals, error) {
	out := new(payload.Meta_Vals)
	err := c.cc.Invoke(ctx, "/meta_manager.Meta/DeleteMetas", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) DeleteMetaInverse(ctx context.Context, in *payload.Meta_Val, opts ...grpc.CallOption) (*payload.Meta_Key, error) {
	out := new(payload.Meta_Key)
	err := c.cc.Invoke(ctx, "/meta_manager.Meta/DeleteMetaInverse", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metaClient) DeleteMetasInverse(ctx context.Context, in *payload.Meta_Vals, opts ...grpc.CallOption) (*payload.Meta_Keys, error) {
	out := new(payload.Meta_Keys)
	err := c.cc.Invoke(ctx, "/meta_manager.Meta/DeleteMetasInverse", in, out, opts...)
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
		FullMethod: "/meta_manager.Meta/GetMeta",
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
		FullMethod: "/meta_manager.Meta/GetMetas",
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
		FullMethod: "/meta_manager.Meta/GetMetaInverse",
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
		FullMethod: "/meta_manager.Meta/GetMetasInverse",
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
		FullMethod: "/meta_manager.Meta/SetMeta",
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
		FullMethod: "/meta_manager.Meta/SetMetas",
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
		FullMethod: "/meta_manager.Meta/DeleteMeta",
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
		FullMethod: "/meta_manager.Meta/DeleteMetas",
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
		FullMethod: "/meta_manager.Meta/DeleteMetaInverse",
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
		FullMethod: "/meta_manager.Meta/DeleteMetasInverse",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetaServer).DeleteMetasInverse(ctx, req.(*payload.Meta_Vals))
	}
	return interceptor(ctx, in, info, handler)
}

var _Meta_serviceDesc = grpc.ServiceDesc{
	ServiceName: "meta_manager.Meta",
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
