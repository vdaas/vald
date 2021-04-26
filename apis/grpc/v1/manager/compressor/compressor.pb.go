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
	proto.RegisterFile("apis/proto/v1/manager/compressor/compressor.proto", fileDescriptor_65a3baf9652f4ae9)
}
func init() {
	golang_proto.RegisterFile("apis/proto/v1/manager/compressor/compressor.proto", fileDescriptor_65a3baf9652f4ae9)
}

var fileDescriptor_65a3baf9652f4ae9 = []byte{
	// 501 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x94, 0xc1, 0x8a, 0xd4, 0x30,
	0x18, 0xc7, 0x89, 0xe2, 0xb0, 0x13, 0xd9, 0xae, 0x1b, 0x67, 0xc5, 0xed, 0xa1, 0x60, 0x5d, 0x44,
	0x06, 0x4c, 0x18, 0x45, 0x84, 0x3d, 0x8e, 0x88, 0x0c, 0xa8, 0x0c, 0x73, 0x18, 0x64, 0x3d, 0x65,
	0xda, 0x6c, 0x0d, 0xb6, 0x4d, 0x6c, 0xd2, 0xc2, 0x22, 0x5e, 0xbc, 0xf8, 0x00, 0xbe, 0x90, 0xc7,
	0x3d, 0x0a, 0xbe, 0x80, 0xcc, 0xfa, 0x02, 0xbe, 0x81, 0x24, 0x69, 0x6b, 0xd5, 0xce, 0xce, 0x69,
	0x92, 0xf9, 0xfe, 0xdf, 0xef, 0xff, 0xcf, 0x07, 0x5f, 0xe1, 0x84, 0x4a, 0xae, 0x88, 0x2c, 0x84,
	0x16, 0xa4, 0x9a, 0x90, 0x8c, 0xe6, 0x34, 0x61, 0x05, 0x89, 0x44, 0x26, 0x0b, 0xa6, 0x94, 0xe8,
	0x1e, 0xb1, 0x95, 0xa1, 0x83, 0x5a, 0x84, 0x3b, 0x95, 0x6a, 0xe2, 0xdf, 0xfd, 0x9b, 0x24, 0xe9,
	0x59, 0x2a, 0x68, 0xdc, 0xfc, 0xba, 0x5e, 0xff, 0x41, 0xc2, 0xf5, 0xdb, 0x72, 0x65, 0x5a, 0x49,
	0x22, 0x12, 0xe1, 0xf4, 0xab, 0xf2, 0xd4, 0xde, 0x5c, 0xb3, 0x39, 0xd5, 0xf2, 0x27, 0xff, 0xca,
	0x13, 0x21, 0x92, 0x94, 0x59, 0x27, 0x77, 0x24, 0x54, 0x72, 0x42, 0xf3, 0x5c, 0x68, 0xaa, 0xb9,
	0xc8, 0x95, 0x6b, 0x7c, 0xf8, 0xeb, 0x1a, 0x1c, 0x4c, 0x69, 0xf4, 0xae, 0x94, 0x28, 0x86, 0xc3,
	0xe7, 0x4c, 0x2f, 0x59, 0xa4, 0x45, 0x81, 0x8e, 0x70, 0x93, 0xa7, 0x9a, 0x60, 0x27, 0xc0, 0x6d,
	0x15, 0x2f, 0xd8, 0xfb, 0x92, 0x29, 0xed, 0x1f, 0xf6, 0xa8, 0x9c, 0x24, 0xbc, 0xf5, 0xe9, 0xfb,
	0xcf, 0x2f, 0x57, 0x6e, 0x20, 0x8f, 0x54, 0xf6, 0x0f, 0xf2, 0xa1, 0x2c, 0x79, 0xfc, 0x11, 0xad,
	0xe0, 0xf0, 0x85, 0x88, 0x5c, 0x86, 0x5e, 0x97, 0xb6, 0xda, 0xba, 0x8c, 0xba, 0xaa, 0x59, 0x7e,
	0x2a, 0xf0, 0x6c, 0xae, 0xc2, 0x43, 0x6b, 0x70, 0x13, 0xed, 0x93, 0xb4, 0xe9, 0x68, 0x3c, 0x5e,
	0xc1, 0x9d, 0x05, 0x4b, 0xb8, 0xd2, 0xac, 0x40, 0x9b, 0x23, 0xfa, 0xfb, 0xdd, 0xd2, 0xb3, 0x4c,
	0xea, 0xb3, 0x70, 0x64, 0xa1, 0x5e, 0x38, 0x24, 0x45, 0x0d, 0x38, 0x06, 0x63, 0x74, 0x02, 0x77,
	0x1b, 0xde, 0xcb, 0x32, 0xd5, 0x1c, 0xf9, 0x1b, 0xa1, 0xaa, 0x8f, 0xea, 0x5b, 0xea, 0x28, 0xdc,
	0x6b, 0xa9, 0x24, 0x33, 0x1c, 0xc3, 0x5e, 0xc2, 0xc1, 0x82, 0x65, 0xa2, 0x62, 0xe8, 0x4e, 0x0f,
	0xd4, 0x95, 0xda, 0x49, 0xf4, 0xb0, 0xeb, 0x39, 0x8f, 0x3d, 0x12, 0xb3, 0x94, 0x69, 0xf6, 0x67,
	0xce, 0xd7, 0x5d, 0xb3, 0x4b, 0x7c, 0x6f, 0x2b, 0xdc, 0xea, 0xfa, 0x1c, 0x6e, 0x5b, 0x07, 0x14,
	0xee, 0x36, 0x0e, 0x6d, 0xf6, 0xd7, 0xc6, 0xc3, 0x3d, 0x68, 0x36, 0x57, 0xbd, 0x1e, 0xb3, 0x39,
	0x6e, 0x24, 0x97, 0xbd, 0xc2, 0xb3, 0x1e, 0x3b, 0xe1, 0x55, 0xc2, 0xa5, 0x21, 0xbf, 0x81, 0x43,
	0x97, 0xce, 0x70, 0x8f, 0x36, 0x71, 0xb7, 0xcd, 0xe6, 0xc0, 0x52, 0xf7, 0x42, 0x48, 0xb8, 0xac,
	0xc3, 0x1f, 0x83, 0xf1, 0xf4, 0x33, 0x38, 0x5f, 0x07, 0xe0, 0xdb, 0x3a, 0x00, 0x3f, 0xd6, 0x01,
	0xf8, 0x7a, 0x11, 0x80, 0xf3, 0x8b, 0x00, 0xc0, 0xfb, 0xa2, 0x48, 0x70, 0x15, 0x53, 0xaa, 0x70,
	0x45, 0xd3, 0x18, 0x53, 0xc9, 0x0d, 0xea, 0xff, 0x25, 0x9e, 0x7a, 0x4b, 0x9a, 0xc6, 0x4f, 0xdb,
	0xfb, 0x1c, 0x9c, 0x3c, 0xee, 0xec, 0x9f, 0x45, 0x10, 0x83, 0x20, 0x6e, 0xff, 0x0a, 0x19, 0xf5,
	0x7f, 0x32, 0x56, 0x03, 0xbb, 0x84, 0x8f, 0x7e, 0x07, 0x00, 0x00, 0xff, 0xff, 0x4f, 0x8f, 0xa0,
	0x45, 0x5d, 0x04, 0x00, 0x00,
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
