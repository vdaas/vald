//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	_ "github.com/danielvladco/go-proto-gql/pb"
	proto "github.com/gogo/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/payload"
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

func init() { proto.RegisterFile("meta.proto", fileDescriptor_3b5ea8fe65782bcc) }

var fileDescriptor_3b5ea8fe65782bcc = []byte{
	// 386 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xcd, 0x4a, 0xf3, 0x40,
	0x14, 0x86, 0x49, 0xe9, 0x1f, 0xe7, 0x6b, 0xd3, 0xaf, 0xa3, 0x75, 0x11, 0xa4, 0xc2, 0xe0, 0xaa,
	0x8b, 0x0c, 0xe8, 0x5e, 0x68, 0xb1, 0x88, 0x48, 0xb1, 0x52, 0x71, 0xe1, 0x46, 0xa6, 0xed, 0x10,
	0x03, 0x49, 0x26, 0x66, 0xc6, 0x42, 0xb7, 0xde, 0x82, 0x37, 0xe2, 0x65, 0xb8, 0x14, 0xbc, 0x81,
	0x52, 0x04, 0x6f, 0x43, 0x66, 0x92, 0x58, 0x4d, 0xda, 0x85, 0xdd, 0x84, 0xf9, 0x39, 0xe7, 0xc9,
	0xc3, 0x79, 0x19, 0x00, 0x9f, 0x49, 0x6a, 0x87, 0x11, 0x97, 0x1c, 0xd5, 0xd4, 0xfa, 0xce, 0xa7,
	0x01, 0x75, 0x58, 0x64, 0xd5, 0x43, 0x3a, 0xf7, 0x38, 0x9d, 0xc6, 0x97, 0xd6, 0xbe, 0xc3, 0xb9,
	0xe3, 0x31, 0x42, 0x43, 0x97, 0xd0, 0x20, 0xe0, 0x92, 0x4a, 0x97, 0x07, 0x22, 0xb9, 0xad, 0x85,
	0x63, 0xe2, 0x3c, 0x78, 0xf1, 0xee, 0xe8, 0xb3, 0x04, 0xc5, 0x01, 0x93, 0x14, 0x9d, 0x40, 0xe5,
	0x8c, 0x49, 0xbd, 0x6c, 0xda, 0x29, 0x4f, 0x6d, 0xed, 0x0b, 0x36, 0xb7, 0x32, 0x47, 0x37, 0xd4,
	0xc3, 0xf5, 0xa7, 0xf7, 0x8f, 0xe7, 0x42, 0x05, 0x97, 0x88, 0x72, 0x41, 0x3d, 0xa8, 0x26, 0xfd,
	0x02, 0xa1, 0x1c, 0x40, 0x58, 0x28, 0x47, 0x10, 0xd8, 0xd4, 0x88, 0x2a, 0x2e, 0x6b, 0x84, 0x40,
	0x03, 0x30, 0x13, 0xc6, 0x79, 0x30, 0x63, 0x91, 0x60, 0x28, 0xff, 0x5f, 0x2b, 0x6f, 0x87, 0x5b,
	0x9a, 0xd3, 0xc0, 0x75, 0xe2, 0xc6, 0x7d, 0xb1, 0xd2, 0x15, 0x34, 0x52, 0xa5, 0x94, 0xb7, 0xc6,
	0xc2, 0x5a, 0x63, 0x8b, 0xf7, 0x34, 0xf1, 0x3f, 0x36, 0x7f, 0x11, 0x85, 0x9a, 0xd2, 0x28, 0x99,
	0xd2, 0x6e, 0xae, 0x4d, 0xd9, 0x99, 0xdf, 0xa7, 0x7d, 0x3f, 0x94, 0xf3, 0xec, 0x94, 0xba, 0x50,
	0x1d, 0xa5, 0x53, 0x6a, 0xad, 0x03, 0x88, 0x1c, 0x21, 0x3b, 0xa4, 0x2e, 0xc0, 0x29, 0xf3, 0x98,
	0x64, 0x7f, 0xcf, 0x0a, 0x25, 0x16, 0x7d, 0xf8, 0xb7, 0x42, 0x6c, 0x1f, 0xd7, 0x10, 0x9a, 0x2b,
	0xcc, 0x56, 0x89, 0x75, 0x32, 0x89, 0x5d, 0x03, 0xfa, 0x21, 0xb6, 0x21, 0xb4, 0x8d, 0x7e, 0x49,
	0x68, 0x9d, 0x4c, 0x68, 0x56, 0xf1, 0x65, 0x71, 0x50, 0xe8, 0x5d, 0xbe, 0x2e, 0xdb, 0xc6, 0xdb,
	0xb2, 0x6d, 0x2c, 0x96, 0x6d, 0x03, 0x76, 0x78, 0xe4, 0xd8, 0xb3, 0x29, 0xa5, 0xc2, 0x9e, 0x51,
	0x6f, 0x6a, 0xab, 0xca, 0x5e, 0x51, 0x7d, 0x87, 0xc6, 0xed, 0xa1, 0xe3, 0xca, 0xfb, 0xc7, 0xb1,
	0x3d, 0xe1, 0x3e, 0xd1, 0x35, 0x44, 0xd5, 0xa8, 0xf7, 0x24, 0x88, 0x13, 0x85, 0x13, 0xcd, 0x1d,
	0x97, 0xf5, 0x0b, 0x3a, 0xfe, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x74, 0x4c, 0xb8, 0x20, 0x98, 0x03,
	0x00, 0x00,
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
	DeleteMetasInverse(ctx context.Context, in *payload.Meta_Keys, opts ...grpc.CallOption) (*payload.Meta_Vals, error)
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

func (c *metaClient) DeleteMetasInverse(ctx context.Context, in *payload.Meta_Keys, opts ...grpc.CallOption) (*payload.Meta_Vals, error) {
	out := new(payload.Meta_Vals)
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
	DeleteMetasInverse(context.Context, *payload.Meta_Keys) (*payload.Meta_Vals, error)
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
func (*UnimplementedMetaServer) DeleteMetasInverse(ctx context.Context, req *payload.Meta_Keys) (*payload.Meta_Vals, error) {
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
	in := new(payload.Meta_Keys)
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
		return srv.(MetaServer).DeleteMetasInverse(ctx, req.(*payload.Meta_Keys))
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
	Metadata: "meta.proto",
}
