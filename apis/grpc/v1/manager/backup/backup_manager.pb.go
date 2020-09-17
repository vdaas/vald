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

package backup

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
	proto.RegisterFile("apis/proto/v1/manager/backup/backup_manager.proto", fileDescriptor_a861c800442e9f9a)
}

var fileDescriptor_a861c800442e9f9a = []byte{
	// 471 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0xd3, 0x51, 0x8b, 0xd3, 0x40,
	0x10, 0x00, 0x60, 0xa2, 0x58, 0xae, 0x2b, 0x97, 0xb3, 0xeb, 0x29, 0x1a, 0xa4, 0x70, 0x39, 0xf5,
	0xa1, 0x0f, 0x3b, 0x54, 0xdf, 0xee, 0xb1, 0x22, 0x52, 0xb0, 0x10, 0x0a, 0x9e, 0x20, 0x82, 0x6c,
	0x93, 0x31, 0x2e, 0x26, 0xd9, 0x35, 0xbb, 0x09, 0x9c, 0xe2, 0x8b, 0x7f, 0xc1, 0x3f, 0xe5, 0xa3,
	0xe0, 0x1f, 0x90, 0xe2, 0xb3, 0xbf, 0x41, 0xb2, 0x9b, 0x44, 0xae, 0xd1, 0xe3, 0x9e, 0x02, 0x33,
	0xb3, 0xdf, 0x4c, 0x06, 0x86, 0xcc, 0xb9, 0x12, 0x1a, 0x54, 0x29, 0x8d, 0x84, 0x7a, 0x0e, 0x39,
	0x2f, 0x78, 0x8a, 0x25, 0x6c, 0x78, 0xfc, 0xbe, 0x52, 0xed, 0xe7, 0x4d, 0x1b, 0x65, 0xb6, 0x8c,
	0xfa, 0xe7, 0xa3, 0xc1, 0xf1, 0x79, 0x42, 0xf1, 0xb3, 0x4c, 0xf2, 0xa4, 0xfb, 0xba, 0x47, 0xc1,
	0xbd, 0x54, 0xca, 0x34, 0x43, 0xe0, 0x4a, 0x00, 0x2f, 0x0a, 0x69, 0xb8, 0x11, 0xb2, 0xd0, 0x2e,
	0xfb, 0xe8, 0xf7, 0x35, 0x32, 0x5a, 0x58, 0x95, 0xe6, 0x64, 0xfc, 0x0c, 0xcd, 0x29, 0xc6, 0x46,
	0x96, 0xf4, 0x88, 0x75, 0x8a, 0xcb, 0xb2, 0x3e, 0xc5, 0xd6, 0xf8, 0xa1, 0x42, 0x6d, 0x82, 0x07,
	0xbb, 0x25, 0x4f, 0x64, 0xae, 0x4a, 0xd4, 0x1a, 0x13, 0xb6, 0x42, 0xc3, 0x5d, 0x79, 0x78, 0xfb,
	0xcb, 0x8f, 0x5f, 0x5f, 0xaf, 0xdc, 0xa0, 0x3e, 0xd4, 0x36, 0x00, 0x9f, 0xaa, 0x4a, 0x24, 0x9f,
	0xe9, 0x6b, 0x32, 0x7e, 0x2e, 0x63, 0x37, 0xcc, 0xb0, 0x5d, 0x9f, 0xea, 0xdb, 0x4d, 0xfa, 0x92,
	0x65, 0xf1, 0x56, 0xb2, 0x65, 0xa4, 0xc3, 0xbb, 0x96, 0xbe, 0x49, 0x27, 0x90, 0x75, 0xe5, 0x9d,
	0xfe, 0x92, 0xec, 0xad, 0x31, 0x15, 0xda, 0x60, 0x49, 0x2f, 0x37, 0x68, 0xe0, 0xf7, 0x65, 0x4f,
	0x73, 0x65, 0xce, 0xc2, 0x43, 0xab, 0xfb, 0xe1, 0x18, 0xca, 0x56, 0x3a, 0xf1, 0x66, 0x34, 0x26,
	0xfb, 0x1d, 0xbc, 0xaa, 0x32, 0x23, 0xe8, 0xc3, 0x4b, 0xe9, 0x7a, 0xc0, 0x07, 0x96, 0x3f, 0x0c,
	0x0f, 0x7a, 0x1e, 0xf2, 0x06, 0x6c, 0x9a, 0x44, 0x64, 0xb4, 0xc6, 0x5c, 0xd6, 0x48, 0xa7, 0xbb,
	0xba, 0x8b, 0xf7, 0x5b, 0xd9, 0x55, 0xdb, 0x6d, 0xcf, 0x7c, 0x48, 0x30, 0x43, 0x83, 0x7f, 0xb7,
	0x7d, 0xdd, 0xbd, 0x74, 0x43, 0x1f, 0x5f, 0xcc, 0xda, 0xa2, 0x81, 0x7d, 0xc7, 0xda, 0x34, 0xdc,
	0xef, 0xec, 0x7e, 0xde, 0x75, 0xa3, 0xbb, 0x9f, 0x58, 0x46, 0x7a, 0xa8, 0x2f, 0x23, 0xd6, 0xe5,
	0xff, 0x3b, 0xb9, 0x6f, 0xf5, 0xbd, 0xf0, 0x2a, 0x08, 0xd5, 0x98, 0x2f, 0xc8, 0xd8, 0x0d, 0xd5,
	0x88, 0x47, 0xff, 0x14, 0x2f, 0xdc, 0xc4, 0x2d, 0xeb, 0x1d, 0x84, 0x04, 0x84, 0x6a, 0x07, 0x3e,
	0xf1, 0x66, 0x8b, 0x8f, 0xdf, 0xb6, 0x53, 0xef, 0xfb, 0x76, 0xea, 0xfd, 0xdc, 0x4e, 0x3d, 0x72,
	0x5f, 0x96, 0x29, 0xab, 0x13, 0xce, 0x35, 0xab, 0x79, 0x96, 0x30, 0xae, 0x04, 0xab, 0xe7, 0xac,
	0x3b, 0x3b, 0x77, 0x6f, 0x8b, 0xc9, 0x29, 0xcf, 0x12, 0xd7, 0x78, 0xe5, 0x32, 0x91, 0xf7, 0x6a,
	0x9e, 0x0a, 0xf3, 0xae, 0xda, 0xb0, 0x58, 0xe6, 0x60, 0x05, 0x68, 0x04, 0xb0, 0x27, 0x99, 0x96,
	0x2a, 0x1e, 0x1e, 0xf5, 0x66, 0x64, 0x6f, 0xee, 0xf1, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xcd,
	0x8e, 0x91, 0x98, 0xfb, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

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
	Metadata: "apis/proto/v1/manager/backup/backup_manager.proto",
}
