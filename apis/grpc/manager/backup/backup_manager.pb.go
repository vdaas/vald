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

package backup

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

func init() { proto.RegisterFile("backup_manager.proto", fileDescriptor_4f75347abe93af04) }

var fileDescriptor_4f75347abe93af04 = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xcf, 0x4e, 0x32, 0x31,
	0x14, 0xc5, 0x19, 0xbe, 0x2f, 0x44, 0x2a, 0x90, 0x50, 0x5d, 0x98, 0x51, 0x21, 0x99, 0x25, 0x8b,
	0x36, 0x51, 0x56, 0x26, 0xba, 0x20, 0xfe, 0x09, 0x89, 0x44, 0xc2, 0x82, 0x85, 0x1b, 0xd3, 0x99,
	0xa9, 0x75, 0x74, 0x66, 0x6e, 0xed, 0x14, 0x12, 0x62, 0xdc, 0xf8, 0x0a, 0xbc, 0x08, 0x8f, 0xe1,
	0xd2, 0xc4, 0x17, 0x20, 0xc4, 0x07, 0x31, 0x4c, 0x81, 0x44, 0x19, 0x59, 0xf6, 0x9c, 0x9e, 0xdf,
	0xbd, 0x4d, 0x0f, 0xda, 0x75, 0x99, 0xf7, 0x34, 0x90, 0x77, 0x11, 0x8b, 0x99, 0xe0, 0x8a, 0x48,
	0x05, 0x1a, 0x70, 0xe5, 0xa7, 0x6a, 0x97, 0x25, 0x1b, 0x85, 0xc0, 0x7c, 0x63, 0xdb, 0x07, 0x02,
	0x40, 0x84, 0x9c, 0x32, 0x19, 0x50, 0x16, 0xc7, 0xa0, 0x99, 0x0e, 0x20, 0x4e, 0x16, 0x6e, 0x49,
	0xba, 0x54, 0x3c, 0x87, 0xe6, 0x74, 0x34, 0xfe, 0x87, 0x0a, 0xad, 0x94, 0x86, 0x4f, 0x51, 0xf1,
	0x8a, 0xeb, 0x3e, 0xf7, 0x34, 0x28, 0x8c, 0xc9, 0x92, 0x79, 0xe3, 0x3e, 0x72, 0x4f, 0x93, 0xf6,
	0xb9, 0x6d, 0xff, 0xd6, 0x3a, 0x5c, 0x33, 0x73, 0xdf, 0xc9, 0xe1, 0x26, 0x2a, 0x5e, 0x83, 0x67,
	0x46, 0x65, 0xc6, 0xab, 0x2b, 0xad, 0x1d, 0xdf, 0x03, 0x69, 0x77, 0x13, 0x27, 0x87, 0xbb, 0x68,
	0xab, 0xc7, 0x45, 0x90, 0x68, 0xae, 0xf0, 0x06, 0xbe, 0x5d, 0x59, 0x79, 0x17, 0x91, 0xd4, 0x23,
	0x67, 0x6f, 0x32, 0xad, 0x5b, 0x6f, 0x9f, 0x5f, 0xe3, 0x7c, 0xc5, 0x29, 0x52, 0xb5, 0x40, 0x9c,
	0x58, 0x0d, 0x7c, 0x86, 0xca, 0x4b, 0x62, 0x67, 0x10, 0xea, 0x00, 0xef, 0xff, 0x8d, 0x4d, 0xd6,
	0xb8, 0x39, 0x7c, 0x89, 0x0a, 0x3d, 0x1e, 0xc1, 0x90, 0x67, 0x3e, 0x62, 0xc3, 0x1e, 0x8d, 0x12,
	0x55, 0x69, 0x90, 0xbe, 0x04, 0xfe, 0x2b, 0x6e, 0xa2, 0x6d, 0xc3, 0x31, 0x5b, 0xec, 0xac, 0xc3,
	0x32, 0xa6, 0xdb, 0xff, 0x27, 0xd3, 0x7a, 0xbe, 0x25, 0xdf, 0x67, 0x35, 0xeb, 0x63, 0x56, 0xb3,
	0xa6, 0xb3, 0x9a, 0x85, 0x0e, 0x41, 0x09, 0x32, 0xf4, 0x19, 0x4b, 0xc8, 0x90, 0x85, 0x3e, 0x59,
	0x96, 0xc1, 0xb4, 0xa0, 0x55, 0xed, 0xb3, 0xd0, 0x37, 0x7f, 0xd8, 0x31, 0x4e, 0xd7, 0xba, 0x25,
	0x22, 0xd0, 0x0f, 0x03, 0x97, 0x78, 0x10, 0xd1, 0x34, 0x4a, 0xe7, 0xd1, 0x79, 0x25, 0x12, 0x2a,
	0x94, 0xf4, 0xe8, 0x02, 0x42, 0x0d, 0xc4, 0x2d, 0xa4, 0x75, 0x38, 0xfe, 0x0e, 0x00, 0x00, 0xff,
	0xff, 0x8d, 0x6e, 0x9f, 0x37, 0x71, 0x02, 0x00, 0x00,
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
	GetVector(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_MetaVector, error)
	Locations(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Info_IPs, error)
	Register(ctx context.Context, in *payload.Object_MetaVector, opts ...grpc.CallOption) (*payload.Empty, error)
	RegisterMulti(ctx context.Context, in *payload.Object_MetaVectors, opts ...grpc.CallOption) (*payload.Empty, error)
	Remove(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Empty, error)
	RemoveMulti(ctx context.Context, in *payload.Object_IDs, opts ...grpc.CallOption) (*payload.Empty, error)
}

type backupClient struct {
	cc *grpc.ClientConn
}

func NewBackupClient(cc *grpc.ClientConn) BackupClient {
	return &backupClient{cc}
}

func (c *backupClient) GetVector(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_MetaVector, error) {
	out := new(payload.Object_MetaVector)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/GetVector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Locations(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Info_IPs, error) {
	out := new(payload.Info_IPs)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/Locations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Register(ctx context.Context, in *payload.Object_MetaVector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RegisterMulti(ctx context.Context, in *payload.Object_MetaVectors, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/RegisterMulti", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Remove(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RemoveMulti(ctx context.Context, in *payload.Object_IDs, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/RemoveMulti", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BackupServer is the server API for Backup service.
type BackupServer interface {
	GetVector(context.Context, *payload.Object_ID) (*payload.Object_MetaVector, error)
	Locations(context.Context, *payload.Object_ID) (*payload.Info_IPs, error)
	Register(context.Context, *payload.Object_MetaVector) (*payload.Empty, error)
	RegisterMulti(context.Context, *payload.Object_MetaVectors) (*payload.Empty, error)
	Remove(context.Context, *payload.Object_ID) (*payload.Empty, error)
	RemoveMulti(context.Context, *payload.Object_IDs) (*payload.Empty, error)
}

// UnimplementedBackupServer can be embedded to have forward compatible implementations.
type UnimplementedBackupServer struct {
}

func (*UnimplementedBackupServer) GetVector(ctx context.Context, req *payload.Object_ID) (*payload.Object_MetaVector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVector not implemented")
}
func (*UnimplementedBackupServer) Locations(ctx context.Context, req *payload.Object_ID) (*payload.Info_IPs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Locations not implemented")
}
func (*UnimplementedBackupServer) Register(ctx context.Context, req *payload.Object_MetaVector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedBackupServer) RegisterMulti(ctx context.Context, req *payload.Object_MetaVectors) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterMulti not implemented")
}
func (*UnimplementedBackupServer) Remove(ctx context.Context, req *payload.Object_ID) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (*UnimplementedBackupServer) RemoveMulti(ctx context.Context, req *payload.Object_IDs) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveMulti not implemented")
}

func RegisterBackupServer(s *grpc.Server, srv BackupServer) {
	s.RegisterService(&_Backup_serviceDesc, srv)
}

func _Backup_GetVector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_ID)
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
		return srv.(BackupServer).GetVector(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_Locations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_ID)
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
		return srv.(BackupServer).Locations(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_MetaVector)
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
		return srv.(BackupServer).Register(ctx, req.(*payload.Object_MetaVector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_RegisterMulti_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_MetaVectors)
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
		return srv.(BackupServer).RegisterMulti(ctx, req.(*payload.Object_MetaVectors))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_ID)
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
		return srv.(BackupServer).Remove(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_RemoveMulti_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_IDs)
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
		return srv.(BackupServer).RemoveMulti(ctx, req.(*payload.Object_IDs))
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "backup_manager.proto",
}
