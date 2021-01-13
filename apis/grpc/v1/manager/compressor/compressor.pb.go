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

package compressor

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
	proto.RegisterFile("apis/proto/v1/manager/compressor/compressor.proto", fileDescriptor_65a3baf9652f4ae9)
}

var fileDescriptor_65a3baf9652f4ae9 = []byte{
	// 466 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0xdf, 0x8a, 0xd3, 0x40,
	0x18, 0xc5, 0x89, 0x62, 0xd9, 0x8e, 0x6c, 0xd6, 0x1d, 0xbb, 0xe2, 0x06, 0x29, 0x18, 0x17, 0x91,
	0x5e, 0xcc, 0x47, 0x15, 0x6f, 0xf6, 0xb2, 0x22, 0x52, 0x50, 0x29, 0xbd, 0x28, 0xb2, 0x5e, 0x4d,
	0x93, 0x31, 0x0e, 0x26, 0x99, 0x31, 0x33, 0x09, 0x2c, 0x22, 0x82, 0xaf, 0xe0, 0x4b, 0x79, 0x29,
	0xf8, 0x02, 0x52, 0x7c, 0x01, 0xdf, 0x40, 0x66, 0xa6, 0x89, 0x11, 0xa7, 0xf6, 0x2a, 0x7f, 0xce,
	0x99, 0xdf, 0x39, 0xf3, 0xc1, 0x87, 0xa6, 0x54, 0x72, 0x05, 0xb2, 0x12, 0x5a, 0x40, 0x33, 0x85,
	0x82, 0x96, 0x34, 0x63, 0x15, 0x24, 0xa2, 0x90, 0x15, 0x53, 0x4a, 0xf4, 0x5f, 0x89, 0xb5, 0xe1,
	0x93, 0xad, 0x89, 0xf4, 0x94, 0x66, 0x1a, 0xdd, 0xfb, 0x9b, 0x24, 0xe9, 0x65, 0x2e, 0x68, 0xda,
	0x3e, 0xdd, 0xd9, 0xe8, 0x4e, 0x26, 0x44, 0x96, 0x33, 0xa0, 0x92, 0x03, 0x2d, 0x4b, 0xa1, 0xa9,
	0xe6, 0xa2, 0x54, 0x4e, 0x7d, 0xf8, 0xeb, 0x1a, 0x1a, 0xcc, 0x68, 0xf2, 0xae, 0x96, 0x38, 0x45,
	0xc3, 0x67, 0x4c, 0xaf, 0x58, 0xa2, 0x45, 0x85, 0xcf, 0x48, 0x4b, 0x69, 0xa6, 0xc4, 0x19, 0x48,
	0xa7, 0x92, 0x25, 0x7b, 0x5f, 0x33, 0xa5, 0xa3, 0x53, 0x8f, 0xcb, 0x59, 0xe2, 0x5b, 0x9f, 0xbf,
	0xff, 0xfc, 0x72, 0xe5, 0x06, 0x0e, 0xa1, 0xb1, 0x3f, 0xe0, 0x43, 0x5d, 0xf3, 0xf4, 0x23, 0x5e,
	0xa3, 0xe1, 0x73, 0x91, 0xb8, 0x0e, 0xde, 0x94, 0x4e, 0xed, 0x52, 0x46, 0x7d, 0xd7, 0xbc, 0x7c,
	0x23, 0xc8, 0x7c, 0xa1, 0xe2, 0x53, 0x1b, 0x70, 0x13, 0x1f, 0x43, 0xde, 0x9e, 0x68, 0x33, 0x5e,
	0xa2, 0x83, 0x25, 0xcb, 0xb8, 0xd2, 0xac, 0xc2, 0xbb, 0x2b, 0x46, 0xc7, 0x7d, 0xe9, 0x69, 0x21,
	0xf5, 0x65, 0x3c, 0xb2, 0xd0, 0x30, 0x1e, 0x42, 0xb5, 0x05, 0x9c, 0x07, 0x13, 0x7c, 0x81, 0x0e,
	0x5b, 0xde, 0x8b, 0x3a, 0xd7, 0x1c, 0x47, 0x3b, 0xa1, 0xca, 0x47, 0x8d, 0x2c, 0x75, 0x14, 0x1f,
	0x75, 0x54, 0x28, 0x0c, 0xc7, 0xb0, 0x57, 0x68, 0xb0, 0x64, 0x85, 0x68, 0x18, 0xbe, 0xeb, 0x81,
	0x3a, 0xa9, 0x9b, 0x84, 0x87, 0xbd, 0x9d, 0xf3, 0x24, 0x84, 0x94, 0xe5, 0x4c, 0xb3, 0x3f, 0x73,
	0xbe, 0xee, 0x0e, 0xbb, 0xc6, 0xf7, 0xf7, 0xc2, 0xad, 0xcf, 0x97, 0x70, 0xdb, 0x26, 0xe0, 0xf8,
	0xb0, 0x4d, 0xe8, 0xba, 0xbf, 0x32, 0x19, 0xee, 0x42, 0xf3, 0x85, 0xf2, 0x66, 0xcc, 0x17, 0xa4,
	0xb5, 0xfc, 0xef, 0x16, 0xa1, 0xcd, 0x38, 0x88, 0xaf, 0x02, 0x97, 0x86, 0xfc, 0x1a, 0x0d, 0x5d,
	0x3b, 0xc3, 0x3d, 0xdb, 0xc5, 0xdd, 0x37, 0x9b, 0x13, 0x4b, 0x3d, 0x8a, 0x11, 0x70, 0xb9, 0x2d,
	0x7f, 0x1e, 0x4c, 0x66, 0x9f, 0xbe, 0x6e, 0xc6, 0xc1, 0xb7, 0xcd, 0x38, 0xf8, 0xb1, 0x19, 0x07,
	0xe8, 0x81, 0xa8, 0x32, 0xd2, 0xa4, 0x94, 0x2a, 0xd2, 0xd0, 0x3c, 0x25, 0x54, 0x72, 0x43, 0xf8,
	0x77, 0xe3, 0x66, 0xe1, 0x8a, 0xe6, 0xe9, 0x93, 0xee, 0x7b, 0x11, 0x5c, 0x3c, 0xce, 0xb8, 0x7e,
	0x5b, 0xaf, 0x8d, 0x09, 0x2c, 0x02, 0x0c, 0x02, 0xec, 0x5a, 0x66, 0x95, 0x4c, 0xfc, 0xfb, 0xbd,
	0x1e, 0xd8, 0xdd, 0x7b, 0xf4, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x65, 0xd7, 0x29, 0x74, 0x0a, 0x04,
	0x00, 0x00,
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
	GetVector(ctx context.Context, in *payload.Backup_GetVector_Request, opts ...grpc.CallOption) (*payload.Backup_Vector, error)
	Locations(ctx context.Context, in *payload.Backup_Locations_Request, opts ...grpc.CallOption) (*payload.Info_IPs, error)
	Register(ctx context.Context, in *payload.Backup_Vector, opts ...grpc.CallOption) (*payload.Empty, error)
	RegisterMulti(ctx context.Context, in *payload.Backup_Vectors, opts ...grpc.CallOption) (*payload.Empty, error)
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

func (c *backupClient) GetVector(ctx context.Context, in *payload.Backup_GetVector_Request, opts ...grpc.CallOption) (*payload.Backup_Vector, error) {
	out := new(payload.Backup_Vector)
	err := c.cc.Invoke(ctx, "/manager.compressor.v1.Backup/GetVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Locations(ctx context.Context, in *payload.Backup_Locations_Request, opts ...grpc.CallOption) (*payload.Info_IPs, error) {
	out := new(payload.Info_IPs)
	err := c.cc.Invoke(ctx, "/manager.compressor.v1.Backup/Locations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Register(ctx context.Context, in *payload.Backup_Vector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.compressor.v1.Backup/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RegisterMulti(ctx context.Context, in *payload.Backup_Vectors, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.compressor.v1.Backup/RegisterMulti", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Remove(ctx context.Context, in *payload.Backup_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.compressor.v1.Backup/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RemoveMulti(ctx context.Context, in *payload.Backup_Remove_RequestMulti, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.compressor.v1.Backup/RemoveMulti", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RegisterIPs(ctx context.Context, in *payload.Backup_IP_Register_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.compressor.v1.Backup/RegisterIPs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RemoveIPs(ctx context.Context, in *payload.Backup_IP_Remove_Request, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.compressor.v1.Backup/RemoveIPs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BackupServer is the server API for Backup service.
type BackupServer interface {
	GetVector(context.Context, *payload.Backup_GetVector_Request) (*payload.Backup_Vector, error)
	Locations(context.Context, *payload.Backup_Locations_Request) (*payload.Info_IPs, error)
	Register(context.Context, *payload.Backup_Vector) (*payload.Empty, error)
	RegisterMulti(context.Context, *payload.Backup_Vectors) (*payload.Empty, error)
	Remove(context.Context, *payload.Backup_Remove_Request) (*payload.Empty, error)
	RemoveMulti(context.Context, *payload.Backup_Remove_RequestMulti) (*payload.Empty, error)
	RegisterIPs(context.Context, *payload.Backup_IP_Register_Request) (*payload.Empty, error)
	RemoveIPs(context.Context, *payload.Backup_IP_Remove_Request) (*payload.Empty, error)
}

// UnimplementedBackupServer can be embedded to have forward compatible implementations.
type UnimplementedBackupServer struct {
}

func (*UnimplementedBackupServer) GetVector(ctx context.Context, req *payload.Backup_GetVector_Request) (*payload.Backup_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVector not implemented")
}
func (*UnimplementedBackupServer) Locations(ctx context.Context, req *payload.Backup_Locations_Request) (*payload.Info_IPs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Locations not implemented")
}
func (*UnimplementedBackupServer) Register(ctx context.Context, req *payload.Backup_Vector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedBackupServer) RegisterMulti(ctx context.Context, req *payload.Backup_Vectors) (*payload.Empty, error) {
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
		FullMethod: "/manager.compressor.v1.Backup/GetVector",
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
		FullMethod: "/manager.compressor.v1.Backup/Locations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).Locations(ctx, req.(*payload.Backup_Locations_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Vector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.compressor.v1.Backup/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).Register(ctx, req.(*payload.Backup_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_RegisterMulti_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Backup_Vectors)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).RegisterMulti(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.compressor.v1.Backup/RegisterMulti",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).RegisterMulti(ctx, req.(*payload.Backup_Vectors))
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
		FullMethod: "/manager.compressor.v1.Backup/Remove",
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
		FullMethod: "/manager.compressor.v1.Backup/RemoveMulti",
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
		FullMethod: "/manager.compressor.v1.Backup/RegisterIPs",
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
		FullMethod: "/manager.compressor.v1.Backup/RemoveIPs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).RemoveIPs(ctx, req.(*payload.Backup_IP_Remove_Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _Backup_serviceDesc = grpc.ServiceDesc{
	ServiceName: "manager.compressor.v1.Backup",
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
	Metadata: "apis/proto/v1/manager/compressor/compressor.proto",
}
