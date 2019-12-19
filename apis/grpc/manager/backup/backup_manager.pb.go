//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

package backup

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

func init() { proto.RegisterFile("backup/backup_manager.proto", fileDescriptor_d3d7e5699810d1ca) }

var fileDescriptor_d3d7e5699810d1ca = []byte{
	// 524 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x95, 0xc1, 0x8a, 0xd4, 0x30,
	0x18, 0xc7, 0xe9, 0xaa, 0xc3, 0x4c, 0x74, 0xba, 0x3b, 0x51, 0xc4, 0xed, 0xae, 0x5d, 0xb6, 0xa2,
	0x87, 0x39, 0x24, 0xb0, 0xde, 0xf6, 0x38, 0x22, 0x32, 0xe0, 0x40, 0x19, 0x64, 0x40, 0x5d, 0xd0,
	0xb4, 0x8d, 0xb5, 0xd8, 0x36, 0xd9, 0x26, 0x1d, 0x58, 0xc4, 0x8b, 0x67, 0x6f, 0xbe, 0xc8, 0x3e,
	0x86, 0x47, 0xc1, 0x17, 0x18, 0x06, 0xaf, 0xbe, 0x83, 0xb4, 0x69, 0x82, 0x76, 0xdc, 0x61, 0xe7,
	0xa4, 0x07, 0x4f, 0xa5, 0xf9, 0x7f, 0xf9, 0xfd, 0xff, 0x09, 0xf9, 0xf8, 0xc0, 0x5e, 0x40, 0xc2,
	0x77, 0x25, 0xc7, 0xea, 0xf3, 0x2a, 0x23, 0x39, 0x89, 0x69, 0x81, 0x78, 0xc1, 0x24, 0x83, 0xf6,
	0xef, 0xab, 0x4e, 0x9f, 0x93, 0xb3, 0x94, 0x91, 0x48, 0xc9, 0xce, 0x7e, 0xcc, 0x58, 0x9c, 0x52,
	0x4c, 0x78, 0x82, 0x49, 0x9e, 0x33, 0x49, 0x64, 0xc2, 0x72, 0xd1, 0xa8, 0x37, 0x78, 0x80, 0xe3,
	0xd3, 0x54, 0xfd, 0x1d, 0xfd, 0xb8, 0x06, 0xba, 0x8f, 0x58, 0xc6, 0x0b, 0x2a, 0x04, 0x0c, 0x40,
	0xef, 0x09, 0x95, 0x33, 0x1a, 0x4a, 0x56, 0xc0, 0x43, 0xa4, 0xa9, 0xa3, 0xda, 0x0d, 0x19, 0x09,
	0x4d, 0xe9, 0x69, 0x49, 0x85, 0x74, 0x9c, 0x76, 0xc9, 0x84, 0x4a, 0xa2, 0x6a, 0xbc, 0xdb, 0x1f,
	0xbf, 0x7d, 0xff, 0xbc, 0xb5, 0x03, 0x6d, 0x3c, 0xaf, 0x17, 0xf0, 0xfb, 0xb2, 0x4c, 0xa2, 0x0f,
	0xf0, 0x04, 0xf4, 0x9e, 0xb2, 0x50, 0x25, 0x5a, 0xf5, 0x30, 0x92, 0xf1, 0x18, 0x98, 0x92, 0x71,
	0xfe, 0x86, 0xa1, 0xb1, 0x2f, 0xbc, 0xdd, 0x1a, 0x7d, 0x13, 0x0e, 0x70, 0xaa, 0xcb, 0x35, 0xdd,
	0x07, 0xdd, 0x29, 0x8d, 0x13, 0x21, 0x69, 0x01, 0xd7, 0xa4, 0x73, 0x6c, 0xa3, 0x3d, 0xce, 0xb8,
	0x3c, 0xf3, 0xee, 0x9c, 0x2f, 0x0e, 0xac, 0x1a, 0x6b, 0x7b, 0x3d, 0x5c, 0x34, 0x88, 0x63, 0x6b,
	0x08, 0x4f, 0x40, 0x5f, 0x13, 0x27, 0x65, 0x2a, 0x13, 0xb8, 0x77, 0x31, 0x56, 0xac, 0x70, 0x5d,
	0xc3, 0xbd, 0xe5, 0x6d, 0x1b, 0x2e, 0xce, 0x2a, 0x52, 0x45, 0x7f, 0x06, 0x3a, 0x53, 0x9a, 0xb1,
	0x39, 0x85, 0x6e, 0x1b, 0xab, 0xd6, 0xcd, 0x3d, 0xb4, 0xc9, 0x8e, 0x21, 0xef, 0x0c, 0x6d, 0x1c,
	0xd1, 0x94, 0x4a, 0xaa, 0x6f, 0xe1, 0x35, 0xb8, 0xae, 0x76, 0xab, 0xc4, 0xf7, 0xd6, 0xa3, 0xeb,
	0xa2, 0x15, 0xfe, 0xbe, 0xe1, 0x43, 0xaf, 0xaf, 0xf9, 0x26, 0xf7, 0xac, 0x72, 0x50, 0x87, 0x19,
	0xfb, 0x62, 0xd5, 0x61, 0xec, 0x23, 0xad, 0x5f, 0x78, 0x02, 0x68, 0x1c, 0xba, 0xde, 0x15, 0x9c,
	0xf0, 0x8a, 0xfb, 0x1c, 0xf4, 0x54, 0xb8, 0x8a, 0x7a, 0xf8, 0x47, 0xea, 0xda, 0x5b, 0xd9, 0x35,
	0xcc, 0x6d, 0x0f, 0xe0, 0x84, 0x37, 0xc1, 0x8f, 0xad, 0xa1, 0x73, 0xf5, 0x7c, 0x71, 0xb0, 0x75,
	0xf4, 0xa9, 0x03, 0x3a, 0x8a, 0x06, 0xb3, 0x0d, 0x5f, 0xfb, 0xfd, 0x76, 0x89, 0xee, 0x1b, 0x1a,
	0xfd, 0xfd, 0x87, 0xff, 0xf2, 0x97, 0x87, 0x7f, 0xb9, 0xa0, 0x1b, 0xf4, 0x40, 0xdc, 0xee, 0x81,
	0x07, 0x97, 0x72, 0xf8, 0xdf, 0x0e, 0xff, 0x7a, 0x3b, 0x8c, 0xf8, 0x97, 0xa5, 0x6b, 0x7d, 0x5d,
	0xba, 0xd6, 0x62, 0xe9, 0x5a, 0xe0, 0x2e, 0x2b, 0x62, 0x34, 0x8f, 0x08, 0x11, 0x68, 0x4e, 0xd2,
	0x08, 0xe9, 0xa9, 0xa3, 0xc6, 0xcd, 0x68, 0x30, 0x23, 0x69, 0xa4, 0xbc, 0x27, 0x4a, 0xf1, 0xad,
	0x17, 0x28, 0x4e, 0xe4, 0xdb, 0x32, 0x40, 0x21, 0xcb, 0x70, 0xbd, 0x15, 0x57, 0x5b, 0xab, 0xd9,
	0x23, 0x70, 0x5c, 0xf0, 0x10, 0x37, 0x90, 0x66, 0x92, 0x05, 0x9d, 0x7a, 0xee, 0x3c, 0xfc, 0x19,
	0x00, 0x00, 0xff, 0xff, 0xd4, 0xa5, 0xd2, 0x31, 0xe1, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CompressClient is the client API for Compress service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CompressClient interface {
	GetVector(ctx context.Context, in *payload.Backup_GetVector_Request, opts ...grpc.CallOption) (*payload.Backup_MetaVector, error)
	Locations(ctx context.Context, in *payload.Backup_Locations_Request, opts ...grpc.CallOption) (*payload.Info_IPs, error)
	Register(ctx context.Context, in *payload.Backup_MetaVector, opts ...grpc.CallOption) (*payload.Empty, error)
	RegisterMulti(ctx context.Context, in *payload.Backup_MetaVectors, opts ...grpc.CallOption) (*payload.Empty, error)
	Remove(ctx context.Context, in *payload.Backup_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error)
	RemoveMulti(ctx context.Context, in *payload.Backup_Remove_RequestMulti, opts ...grpc.CallOption) (*payload.Empty, error)
	RegisterIPs(ctx context.Context, in *payload.Backup_IP_Register_Request, opts ...grpc.CallOption) (*payload.Empty, error)
	RemoveIPs(ctx context.Context, in *payload.Backup_IP_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error)
}

type compressClient struct {
	cc *grpc.ClientConn
}

func NewCompressClient(cc *grpc.ClientConn) CompressClient {
	return &compressClient{cc}
}

func (c *compressClient) GetVector(ctx context.Context, in *payload.Backup_GetVector_Request, opts ...grpc.CallOption) (*payload.Backup_MetaVector, error) {
	out := new(payload.Backup_MetaVector)
	err := c.cc.Invoke(ctx, "/backup_manager.Compress/GetVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *compressClient) Locations(ctx context.Context, in *payload.Backup_Locations_Request, opts ...grpc.CallOption) (*payload.Info_IPs, error) {
	out := new(payload.Info_IPs)
	err := c.cc.Invoke(ctx, "/backup_manager.Compress/Locations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *compressClient) Register(ctx context.Context, in *payload.Backup_MetaVector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Compress/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *compressClient) RegisterMulti(ctx context.Context, in *payload.Backup_MetaVectors, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Compress/RegisterMulti", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *compressClient) Remove(ctx context.Context, in *payload.Backup_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Compress/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *compressClient) RemoveMulti(ctx context.Context, in *payload.Backup_Remove_RequestMulti, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Compress/RemoveMulti", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *compressClient) RegisterIPs(ctx context.Context, in *payload.Backup_IP_Register_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Compress/RegisterIPs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *compressClient) RemoveIPs(ctx context.Context, in *payload.Backup_IP_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Compress/RemoveIPs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompressServer is the server API for Compress service.
type CompressServer interface {
	GetVector(context.Context, *payload.Backup_GetVector_Request) (*payload.Backup_MetaVector, error)
	Locations(context.Context, *payload.Backup_Locations_Request) (*payload.Info_IPs, error)
	Register(context.Context, *payload.Backup_MetaVector) (*payload.Empty, error)
	RegisterMulti(context.Context, *payload.Backup_MetaVectors) (*payload.Empty, error)
	Remove(context.Context, *payload.Backup_Remove_Request) (*payload.Empty, error)
	RemoveMulti(context.Context, *payload.Backup_Remove_RequestMulti) (*payload.Empty, error)
	RegisterIPs(context.Context, *payload.Backup_IP_Register_Request) (*payload.Empty, error)
	RemoveIPs(context.Context, *payload.Backup_IP_Remove_Request) (*payload.Empty, error)
}

// UnimplementedCompressServer can be embedded to have forward compatible implementations.
type UnimplementedCompressServer struct {
}

func (*UnimplementedCompressServer) GetVector(ctx context.Context, req *payload.Backup_GetVector_Request) (*payload.Backup_MetaVector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVector not implemented")
}
func (*UnimplementedCompressServer) Locations(ctx context.Context, req *payload.Backup_Locations_Request) (*payload.Info_IPs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Locations not implemented")
}
func (*UnimplementedCompressServer) Register(ctx context.Context, req *payload.Backup_MetaVector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedCompressServer) RegisterMulti(ctx context.Context, req *payload.Backup_MetaVectors) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterMulti not implemented")
}
func (*UnimplementedCompressServer) Remove(ctx context.Context, req *payload.Backup_Remove_Request) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (*UnimplementedCompressServer) RemoveMulti(ctx context.Context, req *payload.Backup_Remove_RequestMulti) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveMulti not implemented")
}
func (*UnimplementedCompressServer) RegisterIPs(ctx context.Context, req *payload.Backup_IP_Register_Request) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterIPs not implemented")
}
func (*UnimplementedCompressServer) RemoveIPs(ctx context.Context, req *payload.Backup_IP_Remove_Request) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveIPs not implemented")
}

func RegisterCompressServer(s *grpc.Server, srv CompressServer) {
	s.RegisterService(&_Compress_serviceDesc, srv)
}

func _Compress_GetVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_GetVector_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompressServer).GetVector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Compress/GetVector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompressServer).GetVector(ctx, req.(*payload.Backup_GetVector_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Compress_Locations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Locations_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompressServer).Locations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Compress/Locations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompressServer).Locations(ctx, req.(*payload.Backup_Locations_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Compress_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_MetaVector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompressServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Compress/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompressServer).Register(ctx, req.(*payload.Backup_MetaVector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Compress_RegisterMulti_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_MetaVectors)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompressServer).RegisterMulti(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Compress/RegisterMulti",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompressServer).RegisterMulti(ctx, req.(*payload.Backup_MetaVectors))
	}
	return interceptor(ctx, in, info, handler)
}

func _Compress_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Remove_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompressServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Compress/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompressServer).Remove(ctx, req.(*payload.Backup_Remove_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Compress_RemoveMulti_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Remove_RequestMulti)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompressServer).RemoveMulti(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Compress/RemoveMulti",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompressServer).RemoveMulti(ctx, req.(*payload.Backup_Remove_RequestMulti))
	}
	return interceptor(ctx, in, info, handler)
}

func _Compress_RegisterIPs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_IP_Register_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompressServer).RegisterIPs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Compress/RegisterIPs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompressServer).RegisterIPs(ctx, req.(*payload.Backup_IP_Register_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Compress_RemoveIPs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_IP_Remove_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompressServer).RemoveIPs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Compress/RemoveIPs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompressServer).RemoveIPs(ctx, req.(*payload.Backup_IP_Remove_Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Compress_serviceDesc = grpc.ServiceDesc{
	ServiceName: "backup_manager.Compress",
	HandlerType: (*CompressServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVector",
			Handler:    _Compress_GetVector_Handler,
		},
		{
			MethodName: "Locations",
			Handler:    _Compress_Locations_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _Compress_Register_Handler,
		},
		{
			MethodName: "RegisterMulti",
			Handler:    _Compress_RegisterMulti_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Compress_Remove_Handler,
		},
		{
			MethodName: "RemoveMulti",
			Handler:    _Compress_RemoveMulti_Handler,
		},
		{
			MethodName: "RegisterIPs",
			Handler:    _Compress_RegisterIPs_Handler,
		},
		{
			MethodName: "RemoveIPs",
			Handler:    _Compress_RemoveIPs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "backup/backup_manager.proto",
}

// BackupClient is the client API for Backup service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BackupClient interface {
	GetVector(ctx context.Context, in *payload.Backup_GetVector_Request, opts ...grpc.CallOption) (*payload.Backup_Compressed_MetaVector, error)
	Locations(ctx context.Context, in *payload.Backup_Locations_Request, opts ...grpc.CallOption) (*payload.Info_IPs, error)
	Register(ctx context.Context, in *payload.Backup_Compressed_MetaVector, opts ...grpc.CallOption) (*payload.Empty, error)
	RegisterMulti(ctx context.Context, in *payload.Backup_Compressed_MetaVectors, opts ...grpc.CallOption) (*payload.Empty, error)
	Remove(ctx context.Context, in *payload.Backup_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error)
	RemoveMulti(ctx context.Context, in *payload.Backup_Remove_RequestMulti, opts ...grpc.CallOption) (*payload.Empty, error)
	RegisterIPs(ctx context.Context, in *payload.Backup_IP_Register_Request, opts ...grpc.CallOption) (*payload.Empty, error)
	RemoveIPs(ctx context.Context, in *payload.Backup_IP_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error)
}

type backupClient struct {
	cc *grpc.ClientConn
}

func NewBackupClient(cc *grpc.ClientConn) BackupClient {
	return &backupClient{cc}
}

func (c *backupClient) GetVector(ctx context.Context, in *payload.Backup_GetVector_Request, opts ...grpc.CallOption) (*payload.Backup_Compressed_MetaVector, error) {
	out := new(payload.Backup_Compressed_MetaVector)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/GetVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Locations(ctx context.Context, in *payload.Backup_Locations_Request, opts ...grpc.CallOption) (*payload.Info_IPs, error) {
	out := new(payload.Info_IPs)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/Locations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Register(ctx context.Context, in *payload.Backup_Compressed_MetaVector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RegisterMulti(ctx context.Context, in *payload.Backup_Compressed_MetaVectors, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/RegisterMulti", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Remove(ctx context.Context, in *payload.Backup_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RemoveMulti(ctx context.Context, in *payload.Backup_Remove_RequestMulti, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/RemoveMulti", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RegisterIPs(ctx context.Context, in *payload.Backup_IP_Register_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/RegisterIPs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RemoveIPs(ctx context.Context, in *payload.Backup_IP_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/RemoveIPs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BackupServer is the server API for Backup service.
type BackupServer interface {
	GetVector(context.Context, *payload.Backup_GetVector_Request) (*payload.Backup_Compressed_MetaVector, error)
	Locations(context.Context, *payload.Backup_Locations_Request) (*payload.Info_IPs, error)
	Register(context.Context, *payload.Backup_Compressed_MetaVector) (*payload.Empty, error)
	RegisterMulti(context.Context, *payload.Backup_Compressed_MetaVectors) (*payload.Empty, error)
	Remove(context.Context, *payload.Backup_Remove_Request) (*payload.Empty, error)
	RemoveMulti(context.Context, *payload.Backup_Remove_RequestMulti) (*payload.Empty, error)
	RegisterIPs(context.Context, *payload.Backup_IP_Register_Request) (*payload.Empty, error)
	RemoveIPs(context.Context, *payload.Backup_IP_Remove_Request) (*payload.Empty, error)
}

// UnimplementedBackupServer can be embedded to have forward compatible implementations.
type UnimplementedBackupServer struct {
}

func (*UnimplementedBackupServer) GetVector(ctx context.Context, req *payload.Backup_GetVector_Request) (*payload.Backup_Compressed_MetaVector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVector not implemented")
}
func (*UnimplementedBackupServer) Locations(ctx context.Context, req *payload.Backup_Locations_Request) (*payload.Info_IPs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Locations not implemented")
}
func (*UnimplementedBackupServer) Register(ctx context.Context, req *payload.Backup_Compressed_MetaVector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedBackupServer) RegisterMulti(ctx context.Context, req *payload.Backup_Compressed_MetaVectors) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterMulti not implemented")
}
func (*UnimplementedBackupServer) Remove(ctx context.Context, req *payload.Backup_Remove_Request) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (*UnimplementedBackupServer) RemoveMulti(ctx context.Context, req *payload.Backup_Remove_RequestMulti) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveMulti not implemented")
}
func (*UnimplementedBackupServer) RegisterIPs(ctx context.Context, req *payload.Backup_IP_Register_Request) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterIPs not implemented")
}
func (*UnimplementedBackupServer) RemoveIPs(ctx context.Context, req *payload.Backup_IP_Remove_Request) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveIPs not implemented")
}

func RegisterBackupServer(s *grpc.Server, srv BackupServer) {
	s.RegisterService(&_Backup_serviceDesc, srv)
}

func _Backup_GetVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_GetVector_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).GetVector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/GetVector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).GetVector(ctx, req.(*payload.Backup_GetVector_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_Locations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Locations_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).Locations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/Locations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).Locations(ctx, req.(*payload.Backup_Locations_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Compressed_MetaVector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).Register(ctx, req.(*payload.Backup_Compressed_MetaVector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_RegisterMulti_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Compressed_MetaVectors)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).RegisterMulti(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/RegisterMulti",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).RegisterMulti(ctx, req.(*payload.Backup_Compressed_MetaVectors))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Remove_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).Remove(ctx, req.(*payload.Backup_Remove_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_RemoveMulti_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Remove_RequestMulti)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).RemoveMulti(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/RemoveMulti",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).RemoveMulti(ctx, req.(*payload.Backup_Remove_RequestMulti))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_RegisterIPs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_IP_Register_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).RegisterIPs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/RegisterIPs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).RegisterIPs(ctx, req.(*payload.Backup_IP_Register_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_RemoveIPs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_IP_Remove_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).RemoveIPs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/RemoveIPs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).RemoveIPs(ctx, req.(*payload.Backup_IP_Remove_Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Backup_serviceDesc = grpc.ServiceDesc{
	ServiceName: "backup_manager.Backup",
	HandlerType: (*BackupServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVector",
			Handler:    _Backup_GetVector_Handler,
		},
		{
			MethodName: "Locations",
			Handler:    _Backup_Locations_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _Backup_Register_Handler,
		},
		{
			MethodName: "RegisterMulti",
			Handler:    _Backup_RegisterMulti_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Backup_Remove_Handler,
		},
		{
			MethodName: "RemoveMulti",
			Handler:    _Backup_RemoveMulti_Handler,
		},
		{
			MethodName: "RegisterIPs",
			Handler:    _Backup_RegisterIPs_Handler,
		},
		{
			MethodName: "RemoveIPs",
			Handler:    _Backup_RemoveIPs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "backup/backup_manager.proto",
}
