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
	// 365 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x41, 0x4a, 0xf3, 0x40,
	0x18, 0x86, 0xff, 0x29, 0x25, 0x3f, 0xc4, 0x58, 0x70, 0x94, 0x2a, 0x41, 0x5b, 0xe8, 0x4a, 0xba,
	0x98, 0x11, 0x75, 0x51, 0xc4, 0x55, 0x51, 0xb1, 0x8b, 0xa2, 0x28, 0x76, 0xe1, 0x46, 0xbe, 0x24,
	0x43, 0x8c, 0x26, 0x99, 0x71, 0x32, 0x29, 0x14, 0x71, 0xe3, 0x15, 0xbc, 0x48, 0x8f, 0x21, 0xae,
	0x04, 0x2f, 0x50, 0x8a, 0x07, 0x91, 0x66, 0xd2, 0x42, 0x6d, 0x37, 0x71, 0x99, 0xf9, 0xf2, 0x3c,
	0xbc, 0x2f, 0xbc, 0xe6, 0x86, 0x03, 0xee, 0x63, 0x2a, 0xee, 0x22, 0x88, 0xc1, 0x67, 0x92, 0x08,
	0xc9, 0x15, 0xc7, 0x95, 0xf9, 0x57, 0x7b, 0x55, 0xc0, 0x20, 0xe4, 0xe0, 0xe9, 0xb3, 0xbd, 0xed,
	0x73, 0xee, 0x87, 0x8c, 0x82, 0x08, 0x28, 0xc4, 0x31, 0x57, 0xa0, 0x02, 0x1e, 0x27, 0xf9, 0xd5,
	0x12, 0x0e, 0xf5, 0x9f, 0x42, 0xfd, 0xb5, 0xff, 0x51, 0x36, 0x8d, 0x76, 0x66, 0xc3, 0xe7, 0xa6,
	0xd1, 0x89, 0x13, 0x26, 0x15, 0xae, 0x92, 0xa9, 0xf0, 0xc2, 0x79, 0x60, 0xae, 0x22, 0x3d, 0xe6,
	0x2a, 0x2e, 0xed, 0xca, 0xec, 0xfd, 0x34, 0x12, 0x6a, 0xd0, 0xa8, 0x0e, 0x47, 0x75, 0xf4, 0xfa,
	0xf5, 0xfd, 0x56, 0xb2, 0x1a, 0xff, 0x69, 0x90, 0xc1, 0x47, 0xa8, 0x89, 0x8f, 0x4d, 0xeb, 0x5a,
	0x49, 0x06, 0x51, 0x41, 0xdf, 0xbf, 0x5d, 0xb4, 0x87, 0x70, 0xcb, 0x5c, 0xe9, 0xa6, 0xa1, 0x0a,
	0x72, 0x78, 0x73, 0x39, 0x9c, 0x2c, 0xd2, 0x93, 0x06, 0x37, 0xc2, 0x03, 0xc5, 0xfe, 0xd8, 0x20,
	0xcd, 0xe0, 0xb9, 0x06, 0x05, 0x7d, 0xf3, 0x0d, 0x72, 0xb8, 0x40, 0x83, 0x33, 0xd3, 0xb8, 0x62,
	0x11, 0xef, 0x33, 0x8c, 0x7f, 0x43, 0x9d, 0x93, 0x85, 0xff, 0xb7, 0x66, 0xe9, 0x2b, 0x4d, 0x8b,
	0xca, 0x0c, 0xa4, 0xcf, 0x81, 0xf7, 0x82, 0x5b, 0xd3, 0xfc, 0x05, 0x6c, 0x3a, 0xfb, 0x61, 0x9e,
	0x3d, 0x07, 0xd7, 0x17, 0xc1, 0x25, 0xb9, 0xed, 0xf2, 0x70, 0x54, 0x2f, 0xb5, 0xc5, 0xfb, 0xb8,
	0x86, 0x3e, 0xc7, 0x35, 0x34, 0x1a, 0xd7, 0x90, 0xb9, 0xc3, 0xa5, 0x4f, 0xfa, 0x1e, 0x40, 0x42,
	0xfa, 0x10, 0x7a, 0x64, 0xba, 0x61, 0x3d, 0xde, 0xf6, 0x5a, 0x0f, 0x42, 0x4f, 0x4f, 0xaf, 0xab,
	0x2f, 0x97, 0xe8, 0x96, 0xf8, 0x81, 0xba, 0x4f, 0x1d, 0xe2, 0xf2, 0x88, 0x66, 0x28, 0x9d, 0xa0,
	0x93, 0x25, 0x27, 0xd4, 0x97, 0xc2, 0xa5, 0xb9, 0x84, 0x6a, 0x89, 0x63, 0x64, 0x2b, 0x3e, 0xf8,
	0x09, 0x00, 0x00, 0xff, 0xff, 0x9a, 0x20, 0x7f, 0x73, 0x28, 0x03, 0x00, 0x00,
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
	Insert(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamInsert(ctx context.Context, opts ...grpc.CallOption) (Backup_StreamInsertClient, error)
	MultiInsert(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error)
	Update(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (Backup_StreamUpdateClient, error)
	MultiUpdate(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error)
	Remove(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamRemove(ctx context.Context, opts ...grpc.CallOption) (Backup_StreamRemoveClient, error)
	MultiRemove(ctx context.Context, in *payload.Object_IDs, opts ...grpc.CallOption) (*payload.Empty, error)
}

type backupClient struct {
	cc *grpc.ClientConn
}

func NewBackupClient(cc *grpc.ClientConn) BackupClient {
	return &backupClient{cc}
}

func (c *backupClient) Insert(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (Backup_StreamInsertClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Backup_serviceDesc.Streams[0], "/backup_manager.Backup/StreamInsert", opts...)
	if err != nil {
		return nil, err
	}
	x := &backupStreamInsertClient{stream}
	return x, nil
}

type Backup_StreamInsertClient interface {
	Send(*payload.Object_Vector) error
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type backupStreamInsertClient struct {
	grpc.ClientStream
}

func (x *backupStreamInsertClient) Send(m *payload.Object_Vector) error {
	return x.ClientStream.SendMsg(m)
}

func (x *backupStreamInsertClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *backupClient) MultiInsert(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/MultiInsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) Update(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backupClient) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (Backup_StreamUpdateClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Backup_serviceDesc.Streams[1], "/backup_manager.Backup/StreamUpdate", opts...)
	if err != nil {
		return nil, err
	}
	x := &backupStreamUpdateClient{stream}
	return x, nil
}

type Backup_StreamUpdateClient interface {
	Send(*payload.Object_Vector) error
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type backupStreamUpdateClient struct {
	grpc.ClientStream
}

func (x *backupStreamUpdateClient) Send(m *payload.Object_Vector) error {
	return x.ClientStream.SendMsg(m)
}

func (x *backupStreamUpdateClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *backupClient) MultiUpdate(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/MultiUpdate", in, out, opts...)
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

func (c *backupClient) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (Backup_StreamRemoveClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Backup_serviceDesc.Streams[2], "/backup_manager.Backup/StreamRemove", opts...)
	if err != nil {
		return nil, err
	}
	x := &backupStreamRemoveClient{stream}
	return x, nil
}

type Backup_StreamRemoveClient interface {
	Send(*payload.Object_ID) error
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type backupStreamRemoveClient struct {
	grpc.ClientStream
}

func (x *backupStreamRemoveClient) Send(m *payload.Object_ID) error {
	return x.ClientStream.SendMsg(m)
}

func (x *backupStreamRemoveClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *backupClient) MultiRemove(ctx context.Context, in *payload.Object_IDs, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/backup_manager.Backup/MultiRemove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BackupServer is the server API for Backup service.
type BackupServer interface {
	Insert(context.Context, *payload.Object_Vector) (*payload.Empty, error)
	StreamInsert(Backup_StreamInsertServer) error
	MultiInsert(context.Context, *payload.Object_Vectors) (*payload.Empty, error)
	Update(context.Context, *payload.Object_Vector) (*payload.Empty, error)
	StreamUpdate(Backup_StreamUpdateServer) error
	MultiUpdate(context.Context, *payload.Object_Vectors) (*payload.Empty, error)
	Remove(context.Context, *payload.Object_ID) (*payload.Empty, error)
	StreamRemove(Backup_StreamRemoveServer) error
	MultiRemove(context.Context, *payload.Object_IDs) (*payload.Empty, error)
}

// UnimplementedBackupServer can be embedded to have forward compatible implementations.
type UnimplementedBackupServer struct {
}

func (*UnimplementedBackupServer) Insert(ctx context.Context, req *payload.Object_Vector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (*UnimplementedBackupServer) StreamInsert(srv Backup_StreamInsertServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamInsert not implemented")
}
func (*UnimplementedBackupServer) MultiInsert(ctx context.Context, req *payload.Object_Vectors) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiInsert not implemented")
}
func (*UnimplementedBackupServer) Update(ctx context.Context, req *payload.Object_Vector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedBackupServer) StreamUpdate(srv Backup_StreamUpdateServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpdate not implemented")
}
func (*UnimplementedBackupServer) MultiUpdate(ctx context.Context, req *payload.Object_Vectors) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpdate not implemented")
}
func (*UnimplementedBackupServer) Remove(ctx context.Context, req *payload.Object_ID) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (*UnimplementedBackupServer) StreamRemove(srv Backup_StreamRemoveServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamRemove not implemented")
}
func (*UnimplementedBackupServer) MultiRemove(ctx context.Context, req *payload.Object_IDs) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiRemove not implemented")
}

func RegisterBackupServer(s *grpc.Server, srv BackupServer) {
	s.RegisterService(&_Backup_serviceDesc, srv)
}

func _Backup_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/Insert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).Insert(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_StreamInsert_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BackupServer).StreamInsert(&backupStreamInsertServer{stream})
}

type Backup_StreamInsertServer interface {
	Send(*payload.Empty) error
	Recv() (*payload.Object_Vector, error)
	grpc.ServerStream
}

type backupStreamInsertServer struct {
	grpc.ServerStream
}

func (x *backupStreamInsertServer) Send(m *payload.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *backupStreamInsertServer) Recv() (*payload.Object_Vector, error) {
	m := new(payload.Object_Vector)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Backup_MultiInsert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vectors)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).MultiInsert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/MultiInsert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).MultiInsert(ctx, req.(*payload.Object_Vectors))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).Update(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backup_StreamUpdate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BackupServer).StreamUpdate(&backupStreamUpdateServer{stream})
}

type Backup_StreamUpdateServer interface {
	Send(*payload.Empty) error
	Recv() (*payload.Object_Vector, error)
	grpc.ServerStream
}

type backupStreamUpdateServer struct {
	grpc.ServerStream
}

func (x *backupStreamUpdateServer) Send(m *payload.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *backupStreamUpdateServer) Recv() (*payload.Object_Vector, error) {
	m := new(payload.Object_Vector)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Backup_MultiUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vectors)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).MultiUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/MultiUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).MultiUpdate(ctx, req.(*payload.Object_Vectors))
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

func _Backup_StreamRemove_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BackupServer).StreamRemove(&backupStreamRemoveServer{stream})
}

type Backup_StreamRemoveServer interface {
	Send(*payload.Empty) error
	Recv() (*payload.Object_ID, error)
	grpc.ServerStream
}

type backupStreamRemoveServer struct {
	grpc.ServerStream
}

func (x *backupStreamRemoveServer) Send(m *payload.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *backupStreamRemoveServer) Recv() (*payload.Object_ID, error) {
	m := new(payload.Object_ID)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Backup_MultiRemove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_IDs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackupServer).MultiRemove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backup_manager.Backup/MultiRemove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackupServer).MultiRemove(ctx, req.(*payload.Object_IDs))
	}
	return interceptor(ctx, in, info, handler)
}

var _Backup_serviceDesc = grpc.ServiceDesc{
	ServiceName: "backup_manager.Backup",
	HandlerType: (*BackupServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Insert",
			Handler:    _Backup_Insert_Handler,
		},
		{
			MethodName: "MultiInsert",
			Handler:    _Backup_MultiInsert_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Backup_Update_Handler,
		},
		{
			MethodName: "MultiUpdate",
			Handler:    _Backup_MultiUpdate_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Backup_Remove_Handler,
		},
		{
			MethodName: "MultiRemove",
			Handler:    _Backup_MultiRemove_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamInsert",
			Handler:       _Backup_StreamInsert_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamUpdate",
			Handler:       _Backup_StreamUpdate_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamRemove",
			Handler:       _Backup_StreamRemove_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "backup_manager.proto",
}
