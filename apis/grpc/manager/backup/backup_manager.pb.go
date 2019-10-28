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
	// 336 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xcd, 0x4a, 0x03, 0x31,
	0x14, 0x85, 0x9d, 0x22, 0x55, 0x63, 0x5b, 0x30, 0x4a, 0x95, 0x41, 0x5b, 0xe8, 0xb2, 0x8b, 0x04,
	0xfc, 0xd9, 0x74, 0x59, 0x54, 0x50, 0x2c, 0x4a, 0x17, 0x5d, 0xb8, 0x91, 0xcc, 0x4c, 0x88, 0xd1,
	0x99, 0xb9, 0x31, 0x49, 0x0b, 0x45, 0xdc, 0xf8, 0x0a, 0xbe, 0x48, 0x1f, 0x43, 0x70, 0x23, 0xf8,
	0x02, 0xa5, 0xf8, 0x20, 0xd2, 0x49, 0x5b, 0xd0, 0xce, 0xc2, 0x65, 0xee, 0xb9, 0xe7, 0xcb, 0xe1,
	0x1e, 0xb4, 0x13, 0xb0, 0xf0, 0xb1, 0xaf, 0xee, 0x12, 0x96, 0x32, 0xc1, 0x35, 0x51, 0x1a, 0x2c,
	0xe0, 0xca, 0xef, 0xa9, 0x5f, 0x56, 0x6c, 0x18, 0x03, 0x8b, 0x9c, 0xec, 0xef, 0x0b, 0x00, 0x11,
	0x73, 0xca, 0x94, 0xa4, 0x2c, 0x4d, 0xc1, 0x32, 0x2b, 0x21, 0x35, 0x33, 0xb5, 0xa4, 0x02, 0x2a,
	0x9e, 0x62, 0xf7, 0x3a, 0xfc, 0x28, 0xa0, 0x62, 0x3b, 0xa3, 0xe1, 0x13, 0xb4, 0x71, 0x05, 0xa1,
	0xdb, 0xc5, 0x55, 0x32, 0x67, 0x5e, 0x07, 0x0f, 0x3c, 0xb4, 0xa4, 0xc7, 0x43, 0x0b, 0xda, 0xaf,
	0x2c, 0xe6, 0x67, 0x89, 0xb2, 0xc3, 0xc6, 0x0a, 0xbe, 0x44, 0xeb, 0x5d, 0x2e, 0xa4, 0xb1, 0x5c,
	0xff, 0xdb, 0x55, 0x1d, 0x8d, 0xeb, 0xde, 0xeb, 0xd7, 0xf7, 0x5b, 0xa1, 0xd4, 0x58, 0xa3, 0x32,
	0x35, 0x5c, 0xdb, 0x96, 0xd7, 0xc4, 0x2d, 0x54, 0x9e, 0xb3, 0x3a, 0xfd, 0xd8, 0x4a, 0xbc, 0x9b,
	0x0f, 0x34, 0x39, 0x39, 0xce, 0x51, 0xb1, 0xcb, 0x13, 0x18, 0x70, 0x8c, 0xff, 0x9a, 0x2e, 0x4e,
	0x97, 0xf6, 0xf7, 0x16, 0x09, 0x2a, 0xcd, 0x12, 0xd5, 0x99, 0x91, 0x3e, 0xcb, 0xe8, 0x05, 0x1f,
	0xa3, 0x4d, 0xc7, 0x71, 0x09, 0xb6, 0x97, 0x61, 0x39, 0xbf, 0xfb, 0xab, 0xa3, 0x71, 0xbd, 0xd0,
	0x56, 0xef, 0x93, 0x9a, 0xf7, 0x39, 0xa9, 0x79, 0xe3, 0x49, 0xcd, 0x43, 0x07, 0xa0, 0x05, 0x19,
	0x44, 0x8c, 0x19, 0x32, 0x60, 0x71, 0x44, 0xe6, 0x25, 0xba, 0xf6, 0xda, 0x5b, 0x3d, 0x16, 0x47,
	0xee, 0xf6, 0x1d, 0xa7, 0xdc, 0x78, 0xb7, 0x44, 0x48, 0x7b, 0xdf, 0x0f, 0x48, 0x08, 0x09, 0xcd,
	0xac, 0x74, 0x6a, 0x9d, 0x56, 0x69, 0xa8, 0xd0, 0x2a, 0xa4, 0x33, 0x08, 0x75, 0x90, 0xa0, 0x98,
	0xd5, 0x78, 0xf4, 0x13, 0x00, 0x00, 0xff, 0xff, 0x0e, 0x7e, 0x1d, 0x98, 0x29, 0x02, 0x00, 0x00,
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
	Locations(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error)
	Register(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error)
	RegisterMulti(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error)
	Remove(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Empty, error)
	RemoveMulti(ctx context.Context, in *payload.Object_IDs, opts ...grpc.CallOption) (*payload.Empty, error)
}

type backupClient struct {
	cc *grpc.ClientConn
}

func NewBackupClient(cc *grpc.ClientConn) BackupClient {
	return &backupClient{cc}
}

func (c *backupClient) Locations(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/Locations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Register(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) RegisterMulti(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error) {
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
	Locations(context.Context, *payload.Object_Vector) (*payload.Empty, error)
	Register(context.Context, *payload.Object_Vector) (*payload.Empty, error)
	RegisterMulti(context.Context, *payload.Object_Vectors) (*payload.Empty, error)
	Remove(context.Context, *payload.Object_ID) (*payload.Empty, error)
	RemoveMulti(context.Context, *payload.Object_IDs) (*payload.Empty, error)
}

// UnimplementedBackupServer can be embedded to have forward compatible implementations.
type UnimplementedBackupServer struct {
}

func (*UnimplementedBackupServer) Locations(ctx context.Context, req *payload.Object_Vector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Locations not implemented")
}
func (*UnimplementedBackupServer) Register(ctx context.Context, req *payload.Object_Vector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedBackupServer) RegisterMulti(ctx context.Context, req *payload.Object_Vectors) (*payload.Empty, error) {
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

func _Backup_Locations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vector)
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
		return srv.(BackupServer).Locations(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vector)
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
		return srv.(BackupServer).Register(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_RegisterMulti_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vectors)
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
		return srv.(BackupServer).RegisterMulti(ctx, req.(*payload.Object_Vectors))
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
