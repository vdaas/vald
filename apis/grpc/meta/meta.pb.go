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
	// 308 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xcf, 0x4a, 0xc3, 0x30,
	0x18, 0xc0, 0xe9, 0x9c, 0xdb, 0x88, 0xdb, 0xc0, 0xa8, 0x97, 0x20, 0x13, 0x8a, 0xe7, 0x04, 0xf4,
	0x2e, 0x6c, 0x38, 0x3c, 0x88, 0x28, 0x0c, 0x76, 0xf0, 0x22, 0x5f, 0xd7, 0x10, 0x0b, 0x69, 0x13,
	0x9b, 0x38, 0xe8, 0xd5, 0x57, 0xf0, 0xec, 0x3b, 0xf8, 0x18, 0x1e, 0x05, 0x5f, 0x60, 0x14, 0x1f,
	0x44, 0x92, 0x75, 0x0a, 0xad, 0x17, 0xbd, 0x84, 0xef, 0x5f, 0x7e, 0xfc, 0xf8, 0x12, 0x84, 0x52,
	0x6e, 0x81, 0xea, 0x5c, 0x59, 0x85, 0xfb, 0x2e, 0xbe, 0x4b, 0x21, 0x03, 0xc1, 0x73, 0x32, 0xd0,
	0x50, 0x48, 0x05, 0xf1, 0xba, 0x49, 0x0e, 0x85, 0x52, 0x42, 0x72, 0x06, 0x3a, 0x61, 0x90, 0x65,
	0xca, 0x82, 0x4d, 0x54, 0x66, 0xaa, 0x6e, 0x5f, 0x47, 0x4c, 0x3c, 0xc8, 0x75, 0x76, 0xf2, 0xb2,
	0x85, 0xda, 0x57, 0xdc, 0x02, 0x3e, 0x43, 0xdd, 0x0b, 0x6e, 0x7d, 0xb8, 0x4b, 0x37, 0x3c, 0x97,
	0xd2, 0x4b, 0x5e, 0x90, 0x5a, 0x69, 0x0e, 0x32, 0x1c, 0x3c, 0x7d, 0x7c, 0x3e, 0xb7, 0xba, 0xe1,
	0x36, 0x73, 0x2e, 0x78, 0x82, 0x7a, 0xd5, 0x7d, 0x83, 0x71, 0x03, 0x60, 0x08, 0x6e, 0x10, 0x4c,
	0x38, 0xf4, 0x88, 0x5e, 0xd8, 0xf1, 0x08, 0xe3, 0x1c, 0x66, 0x95, 0xc3, 0x7e, 0x03, 0x31, 0x07,
	0x49, 0x86, 0xdf, 0xd5, 0x69, 0xaa, 0x6d, 0x51, 0x77, 0x18, 0xa3, 0xde, 0x6c, 0xe3, 0x70, 0xf0,
	0x1b, 0xc0, 0x34, 0x08, 0x75, 0x85, 0x31, 0x42, 0xe7, 0x5c, 0x72, 0xcb, 0xff, 0xbe, 0x09, 0x5c,
	0x59, 0x4c, 0xd1, 0xce, 0x0f, 0xe2, 0xdf, 0xcb, 0x20, 0xed, 0xd7, 0xd5, 0x51, 0x6b, 0x72, 0xfd,
	0x56, 0x8e, 0x82, 0xf7, 0x72, 0x14, 0xac, 0xca, 0x51, 0x80, 0xf6, 0x54, 0x2e, 0xe8, 0x32, 0x06,
	0x30, 0x74, 0x09, 0x32, 0xa6, 0x6e, 0x72, 0xd2, 0x76, 0xe7, 0x4d, 0x70, 0x7b, 0x2c, 0x12, 0x7b,
	0xff, 0x18, 0xd1, 0x85, 0x4a, 0x99, 0x9f, 0x61, 0x6e, 0xc6, 0xfd, 0x02, 0xc3, 0x44, 0xae, 0x17,
	0x9e, 0x1b, 0x75, 0xfc, 0xbb, 0x9f, 0x7e, 0x05, 0x00, 0x00, 0xff, 0xff, 0x67, 0x60, 0x2d, 0xbe,
	0x4e, 0x02, 0x00, 0x00,
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
	SetMeta(ctx context.Context, in *payload.Meta_KeyVal, opts ...grpc.CallOption) (*payload.Empty, error)
	SetMetas(ctx context.Context, in *payload.Meta_KeyVals, opts ...grpc.CallOption) (*payload.Empty, error)
	DeleteMeta(ctx context.Context, in *payload.Meta_Key, opts ...grpc.CallOption) (*payload.Meta_Val, error)
	DeleteMetas(ctx context.Context, in *payload.Meta_Keys, opts ...grpc.CallOption) (*payload.Meta_Vals, error)
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

// MetaServer is the server API for Meta service.
type MetaServer interface {
	GetMeta(context.Context, *payload.Meta_Key) (*payload.Meta_Val, error)
	GetMetas(context.Context, *payload.Meta_Keys) (*payload.Meta_Vals, error)
	SetMeta(context.Context, *payload.Meta_KeyVal) (*payload.Empty, error)
	SetMetas(context.Context, *payload.Meta_KeyVals) (*payload.Empty, error)
	DeleteMeta(context.Context, *payload.Meta_Key) (*payload.Meta_Val, error)
	DeleteMetas(context.Context, *payload.Meta_Keys) (*payload.Meta_Vals, error)
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "meta.proto",
}
