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
	math "math"

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

func init() {
	proto.RegisterFile("apis/proto/manager/backup/backup_manager.proto", fileDescriptor_e078ce592fdab8c1)
}

var fileDescriptor_e078ce592fdab8c1 = []byte{
	// 464 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xc1, 0x8a, 0xd4, 0x40,
	0x10, 0x86, 0x89, 0xe2, 0xb0, 0xd3, 0xb2, 0x59, 0xa7, 0x5d, 0x45, 0x83, 0x0e, 0x6e, 0x44, 0x0f,
	0x73, 0xe8, 0x02, 0xbd, 0xed, 0x71, 0x44, 0x64, 0xc0, 0x81, 0x30, 0xe0, 0x0a, 0x22, 0x48, 0x4f,
	0x52, 0xc6, 0xc6, 0x24, 0xdd, 0xa6, 0x3b, 0x03, 0x8b, 0x78, 0xf1, 0x15, 0x7c, 0x29, 0x8f, 0x82,
	0x2f, 0x20, 0x83, 0x67, 0x9f, 0x41, 0xd2, 0x9d, 0xb4, 0x3a, 0x71, 0x97, 0x3d, 0x35, 0xd4, 0xff,
	0xd7, 0x57, 0x7f, 0x17, 0x14, 0x61, 0x5c, 0x09, 0x0d, 0xaa, 0x96, 0x46, 0x42, 0xc9, 0x2b, 0x9e,
	0x63, 0x0d, 0x6b, 0x9e, 0xbe, 0x6f, 0x54, 0xf7, 0xbc, 0xe9, 0xaa, 0xcc, 0x7a, 0x68, 0xf8, 0x6f,
	0x35, 0xba, 0xf7, 0x57, 0xbf, 0xe2, 0xa7, 0x85, 0xe4, 0x59, 0xff, 0xba, 0x8e, 0xe8, 0x4e, 0x2e,
	0x65, 0x5e, 0x20, 0x70, 0x25, 0x80, 0x57, 0x95, 0x34, 0xdc, 0x08, 0x59, 0x69, 0xa7, 0x3e, 0xfa,
	0x75, 0x85, 0x8c, 0xe6, 0x16, 0x49, 0x4b, 0x32, 0x7e, 0x86, 0xe6, 0x04, 0x53, 0x23, 0x6b, 0x7a,
	0xc4, 0x7a, 0x8a, 0x53, 0x99, 0x97, 0xd8, 0x0a, 0x3f, 0x34, 0xa8, 0x4d, 0xf4, 0x60, 0xd7, 0xf2,
	0x44, 0x96, 0xaa, 0x46, 0xad, 0x31, 0x63, 0x4b, 0x34, 0xdc, 0xd9, 0xe3, 0x9b, 0x9f, 0xbf, 0xff,
	0xfc, 0x72, 0xe9, 0x1a, 0x0d, 0x61, 0x63, 0x0b, 0xf0, 0xb1, 0x69, 0x44, 0xf6, 0x89, 0xbe, 0x26,
	0xe3, 0xe7, 0x32, 0x75, 0x61, 0x86, 0xe3, 0xbc, 0xe4, 0xc7, 0x4d, 0xbc, 0x65, 0x51, 0xbd, 0x95,
	0x6c, 0x91, 0xe8, 0xf8, 0xb6, 0x45, 0x5f, 0xa7, 0x13, 0x28, 0x7a, 0x7b, 0x4f, 0x7f, 0x49, 0xf6,
	0x56, 0x98, 0x0b, 0x6d, 0xb0, 0xa6, 0x17, 0x0b, 0x1a, 0x85, 0xde, 0xf6, 0xb4, 0x54, 0xe6, 0x34,
	0x3e, 0xb4, 0xf4, 0x30, 0x1e, 0x43, 0xdd, 0x91, 0x8e, 0x83, 0x19, 0x4d, 0xc9, 0x7e, 0x0f, 0x5e,
	0x36, 0x85, 0x11, 0xf4, 0xe1, 0x85, 0xe8, 0x7a, 0x80, 0x8f, 0x2c, 0xfe, 0x30, 0x3e, 0xf0, 0x78,
	0x28, 0x5b, 0x60, 0x3b, 0x24, 0x21, 0xa3, 0x15, 0x96, 0x72, 0x83, 0x74, 0xba, 0x4b, 0x77, 0x75,
	0xbf, 0x95, 0x5d, 0x6a, 0xb7, 0xed, 0x59, 0x08, 0x19, 0x16, 0x68, 0xf0, 0xcf, 0xb6, 0xaf, 0xba,
	0x4e, 0x17, 0xfa, 0xfe, 0xf9, 0x58, 0x6b, 0x1a, 0xb0, 0x6f, 0x59, 0x36, 0x8d, 0xf7, 0x7b, 0xb6,
	0xcf, 0xbb, 0x6a, 0xe9, 0xee, 0x13, 0x8b, 0x44, 0x0f, 0xe9, 0x8b, 0x84, 0xf5, 0xfa, 0x99, 0xc9,
	0x43, 0x4b, 0xdf, 0x8b, 0x2f, 0x83, 0x50, 0x2d, 0xf3, 0x05, 0x19, 0xbb, 0x50, 0x2d, 0xf1, 0xe8,
	0xbf, 0xc4, 0x73, 0x37, 0x71, 0xc3, 0xf2, 0x0e, 0x62, 0x02, 0x42, 0x75, 0x81, 0x8f, 0x83, 0xd9,
	0x5c, 0x7d, 0xdd, 0x4e, 0x83, 0x6f, 0xdb, 0x69, 0xf0, 0x63, 0x3b, 0x0d, 0xc8, 0x5d, 0x59, 0xe7,
	0x6c, 0x93, 0x71, 0xae, 0xd9, 0x86, 0x17, 0x19, 0xeb, 0x8f, 0xcd, 0x5d, 0xd9, 0x7c, 0x72, 0xc2,
	0x8b, 0xcc, 0x4d, 0x5c, 0x3a, 0x25, 0x09, 0x5e, 0xb1, 0x5c, 0x98, 0x77, 0xcd, 0x9a, 0xa5, 0xb2,
	0x04, 0xdb, 0x0a, 0x6d, 0x2b, 0xd8, 0x43, 0xcc, 0x6b, 0x95, 0xee, 0xdc, 0xf1, 0x7a, 0x64, 0x2f,
	0xed, 0xf1, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8c, 0x93, 0x88, 0x8a, 0xeb, 0x03, 0x00, 0x00,
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
	Metadata: "apis/proto/manager/backup/backup_manager.proto",
}
