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

func init() {
	proto.RegisterFile("apis/proto/v1/manager/backup/backup_manager.proto", fileDescriptor_a861c800442e9f9a)
}

var fileDescriptor_a861c800442e9f9a = []byte{
	// 474 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xdf, 0x8a, 0xd3, 0x40,
	0x14, 0x87, 0x89, 0x42, 0xd9, 0x8e, 0x6c, 0xd7, 0x8e, 0x55, 0x34, 0x48, 0xc1, 0x58, 0xbd, 0xe8,
	0xc5, 0x1c, 0xa2, 0x77, 0x7b, 0x59, 0x11, 0x29, 0xb8, 0x50, 0x7a, 0x51, 0xfc, 0x73, 0x21, 0xd3,
	0xe4, 0x18, 0x07, 0x93, 0xcc, 0x98, 0x99, 0x04, 0x56, 0xf1, 0xc6, 0x47, 0xd0, 0x97, 0xf2, 0x52,
	0xf0, 0x05, 0xa4, 0xf8, 0x20, 0x92, 0x99, 0x24, 0xee, 0xb2, 0xd1, 0xee, 0x55, 0xc8, 0x9c, 0xdf,
	0x7c, 0xdf, 0x99, 0x03, 0x87, 0x84, 0x5c, 0x09, 0x0d, 0xaa, 0x90, 0x46, 0x42, 0x15, 0x42, 0xc6,
	0x73, 0x9e, 0x60, 0x01, 0x5b, 0x1e, 0xbd, 0x2f, 0x55, 0xf3, 0x79, 0xd3, 0x9c, 0x32, 0x1b, 0xa3,
	0xe3, 0xf6, 0xd7, 0x55, 0x59, 0x15, 0xfa, 0xf7, 0xcf, 0x53, 0x14, 0x3f, 0x4d, 0x25, 0x8f, 0xdb,
	0xaf, 0xbb, 0xe7, 0xdf, 0x4d, 0xa4, 0x4c, 0x52, 0x04, 0xae, 0x04, 0xf0, 0x3c, 0x97, 0x86, 0x1b,
	0x21, 0x73, 0xed, 0xaa, 0x8f, 0xbe, 0x0e, 0xc8, 0x60, 0x61, 0x81, 0x54, 0x92, 0xe1, 0x33, 0x34,
	0x1b, 0x8c, 0x8c, 0x2c, 0xe8, 0x8c, 0xb5, 0x94, 0x2a, 0x64, 0x2e, 0xc0, 0xba, 0x2a, 0x5b, 0xe3,
	0x87, 0x12, 0xb5, 0xf1, 0xfb, 0x52, 0x4f, 0x64, 0xa6, 0x0a, 0xd4, 0x1a, 0x63, 0xe6, 0xd2, 0xc1,
	0xad, 0x2f, 0x3f, 0x7f, 0x7f, 0xbb, 0x72, 0x9d, 0x8e, 0xa0, 0xb2, 0x07, 0xf0, 0xa9, 0x2c, 0x45,
	0xfc, 0x99, 0x6e, 0xc9, 0xf0, 0xb9, 0x8c, 0x5c, 0x3b, 0xbd, 0xc2, 0xae, 0xda, 0x09, 0x27, 0x67,
	0x53, 0xcb, 0xfc, 0xad, 0x64, 0xcb, 0x95, 0x0e, 0xee, 0x58, 0xc1, 0x0d, 0x3a, 0x86, 0xb4, 0xbd,
	0xd1, 0x3a, 0x5e, 0x92, 0x83, 0x35, 0x26, 0x42, 0x1b, 0xec, 0x7f, 0xd3, 0x85, 0x6e, 0xfd, 0xf1,
	0xd9, 0xd4, 0xd3, 0x4c, 0x99, 0xd3, 0x60, 0x62, 0xf9, 0xa3, 0x60, 0x08, 0x45, 0xc3, 0x3a, 0xf6,
	0xe6, 0x14, 0xc9, 0x61, 0x8b, 0x3e, 0x29, 0x53, 0x23, 0xe8, 0x83, 0xcb, 0xf0, 0x75, 0x9f, 0xc0,
	0xb7, 0x82, 0x49, 0x70, 0xd4, 0x09, 0x20, 0xab, 0x91, 0xb5, 0x66, 0x43, 0x06, 0x6b, 0xcc, 0x64,
	0x85, 0xf4, 0x5e, 0x0f, 0xdf, 0x95, 0xba, 0xf9, 0xf4, 0xb0, 0x9b, 0xe9, 0xcf, 0x47, 0x10, 0x63,
	0x8a, 0x06, 0xff, 0x4e, 0xff, 0x9a, 0xbb, 0xec, 0x9a, 0x7f, 0xb8, 0x17, 0x6e, 0x73, 0x7d, 0x86,
	0xdb, 0xd6, 0x40, 0x83, 0xc3, 0xd6, 0xd0, 0xf5, 0xfe, 0xa2, 0x76, 0xb8, 0x07, 0x2d, 0x57, 0xba,
	0xd7, 0xb1, 0x5c, 0xb1, 0x36, 0xf2, 0xbf, 0x57, 0x8c, 0xac, 0xe3, 0x20, 0xb8, 0x0a, 0x42, 0xd5,
	0xe4, 0xd7, 0x64, 0xe8, 0xba, 0xab, 0xb9, 0xb3, 0x7f, 0x71, 0xf7, 0xcd, 0xe6, 0xa6, 0xa5, 0x1e,
	0x05, 0x04, 0x84, 0x6a, 0x9a, 0x3f, 0xf6, 0xe6, 0x8b, 0x8f, 0xdf, 0x77, 0x53, 0xef, 0xc7, 0x6e,
	0xea, 0xfd, 0xda, 0x4d, 0x3d, 0x32, 0x93, 0x45, 0xc2, 0xaa, 0x98, 0x73, 0xcd, 0x2a, 0x9e, 0xc6,
	0x8c, 0x2b, 0x51, 0x13, 0xce, 0xaf, 0xe3, 0x62, 0xbc, 0xe1, 0x69, 0xec, 0xdc, 0x27, 0xae, 0xb2,
	0xf2, 0x5e, 0x85, 0x89, 0x30, 0xef, 0xca, 0x2d, 0x8b, 0x64, 0x06, 0x96, 0x00, 0x35, 0x01, 0xec,
	0xda, 0x26, 0x85, 0x8a, 0x2e, 0xee, 0xfe, 0x76, 0x60, 0xf7, 0xf2, 0xf1, 0x9f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x5c, 0xcd, 0x50, 0x7e, 0x22, 0x04, 0x00, 0x00,
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
	GetVector(ctx context.Context, in *payload.Backup_GetVector_Request, opts ...grpc.CallOption) (*payload.Backup_Compressed_Vector, error)
	Locations(ctx context.Context, in *payload.Backup_Locations_Request, opts ...grpc.CallOption) (*payload.Info_IPs, error)
	Register(ctx context.Context, in *payload.Backup_Compressed_Vector, opts ...grpc.CallOption) (*payload.Empty, error)
	RegisterMulti(ctx context.Context, in *payload.Backup_Compressed_Vectors, opts ...grpc.CallOption) (*payload.Empty, error)
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

func (c *backupClient) GetVector(ctx context.Context, in *payload.Backup_GetVector_Request, opts ...grpc.CallOption) (*payload.Backup_Compressed_Vector, error) {
	out := new(payload.Backup_Compressed_Vector)
	err := c.cc.Invoke(ctx, "/manager.backup.v1.Backup/GetVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Locations(ctx context.Context, in *payload.Backup_Locations_Request, opts ...grpc.CallOption) (*payload.Info_IPs, error) {
	out := new(payload.Info_IPs)
	err := c.cc.Invoke(ctx, "/manager.backup.v1.Backup/Locations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Register(ctx context.Context, in *payload.Backup_Compressed_Vector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.backup.v1.Backup/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RegisterMulti(ctx context.Context, in *payload.Backup_Compressed_Vectors, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.backup.v1.Backup/RegisterMulti", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Remove(ctx context.Context, in *payload.Backup_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.backup.v1.Backup/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RemoveMulti(ctx context.Context, in *payload.Backup_Remove_RequestMulti, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.backup.v1.Backup/RemoveMulti", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RegisterIPs(ctx context.Context, in *payload.Backup_IP_Register_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.backup.v1.Backup/RegisterIPs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RemoveIPs(ctx context.Context, in *payload.Backup_IP_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.backup.v1.Backup/RemoveIPs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BackupServer is the server API for Backup service.
type BackupServer interface {
	GetVector(context.Context, *payload.Backup_GetVector_Request) (*payload.Backup_Compressed_Vector, error)
	Locations(context.Context, *payload.Backup_Locations_Request) (*payload.Info_IPs, error)
	Register(context.Context, *payload.Backup_Compressed_Vector) (*payload.Empty, error)
	RegisterMulti(context.Context, *payload.Backup_Compressed_Vectors) (*payload.Empty, error)
	Remove(context.Context, *payload.Backup_Remove_Request) (*payload.Empty, error)
	RemoveMulti(context.Context, *payload.Backup_Remove_RequestMulti) (*payload.Empty, error)
	RegisterIPs(context.Context, *payload.Backup_IP_Register_Request) (*payload.Empty, error)
	RemoveIPs(context.Context, *payload.Backup_IP_Remove_Request) (*payload.Empty, error)
}

// UnimplementedBackupServer can be embedded to have forward compatible implementations.
type UnimplementedBackupServer struct {
}

func (*UnimplementedBackupServer) GetVector(ctx context.Context, req *payload.Backup_GetVector_Request) (*payload.Backup_Compressed_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVector not implemented")
}
func (*UnimplementedBackupServer) Locations(ctx context.Context, req *payload.Backup_Locations_Request) (*payload.Info_IPs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Locations not implemented")
}
func (*UnimplementedBackupServer) Register(ctx context.Context, req *payload.Backup_Compressed_Vector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedBackupServer) RegisterMulti(ctx context.Context, req *payload.Backup_Compressed_Vectors) (*payload.Empty, error) {
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
		FullMethod: "/manager.backup.v1.Backup/GetVector",
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
		FullMethod: "/manager.backup.v1.Backup/Locations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).Locations(ctx, req.(*payload.Backup_Locations_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Compressed_Vector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.backup.v1.Backup/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).Register(ctx, req.(*payload.Backup_Compressed_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_RegisterMulti_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Compressed_Vectors)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).RegisterMulti(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.backup.v1.Backup/RegisterMulti",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).RegisterMulti(ctx, req.(*payload.Backup_Compressed_Vectors))
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
		FullMethod: "/manager.backup.v1.Backup/Remove",
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
		FullMethod: "/manager.backup.v1.Backup/RemoveMulti",
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
		FullMethod: "/manager.backup.v1.Backup/RegisterIPs",
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
		FullMethod: "/manager.backup.v1.Backup/RemoveIPs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).RemoveIPs(ctx, req.(*payload.Backup_IP_Remove_Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Backup_serviceDesc = grpc.ServiceDesc{
	ServiceName: "manager.backup.v1.Backup",
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
