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

package vald

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

func init() { proto.RegisterFile("vald.proto", fileDescriptor_33e68360ff34c546) }

var fileDescriptor_33e68360ff34c546 = []byte{
	// 481 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xc1, 0x6e, 0xd3, 0x30,
	0x18, 0xc7, 0xf1, 0x54, 0x65, 0x60, 0x42, 0x41, 0x06, 0x6d, 0x23, 0x42, 0x9d, 0x54, 0x71, 0x40,
	0x3b, 0xd8, 0x08, 0x38, 0x4c, 0x88, 0x0b, 0x55, 0xb7, 0xd2, 0x43, 0x35, 0xd8, 0xb4, 0x1d, 0xb8,
	0xb9, 0x89, 0x95, 0x19, 0xd2, 0xd8, 0xb3, 0x9d, 0x8a, 0x0a, 0x71, 0xe1, 0x15, 0x78, 0x91, 0x3d,
	0x06, 0x47, 0x24, 0x5e, 0xa0, 0xaa, 0x38, 0xf2, 0x10, 0x53, 0xec, 0x34, 0x5a, 0xdb, 0x5c, 0x92,
	0x5b, 0xfc, 0xf9, 0xfb, 0xfd, 0xf3, 0xff, 0x5b, 0xfa, 0x3e, 0x08, 0xa7, 0x34, 0x89, 0xb0, 0x54,
	0xc2, 0x08, 0xd4, 0xca, 0xbf, 0x83, 0x07, 0x92, 0xce, 0x12, 0x41, 0x8b, 0x62, 0xf0, 0x2c, 0x16,
	0x22, 0x4e, 0x18, 0xa1, 0x92, 0x13, 0x9a, 0xa6, 0xc2, 0x50, 0xc3, 0x45, 0xaa, 0x8b, 0x5b, 0x5f,
	0x8e, 0x49, 0x7c, 0x95, 0xb8, 0xd3, 0xab, 0xff, 0x77, 0x61, 0xeb, 0x82, 0x26, 0x11, 0x3a, 0x86,
	0xde, 0xd1, 0x37, 0xae, 0x8d, 0x46, 0x08, 0x2f, 0xe5, 0x4e, 0xc6, 0x5f, 0x58, 0x68, 0xf0, 0xb0,
	0x1f, 0x54, 0xd4, 0xba, 0x4f, 0x7e, 0xfe, 0xfd, 0xf7, 0x6b, 0xab, 0x8d, 0x7c, 0xc2, 0x2c, 0x48,
	0xbe, 0xf3, 0xe8, 0x07, 0x3a, 0x81, 0xde, 0x19, 0xa3, 0x2a, 0xbc, 0x44, 0xbb, 0x25, 0xe3, 0x0a,
	0xf8, 0x94, 0x5d, 0x65, 0x4c, 0x9b, 0x60, 0x6f, 0xf3, 0x42, 0x4b, 0x91, 0x6a, 0xd6, 0x45, 0x56,
	0xd2, 0xef, 0x6e, 0x13, 0x6d, 0x6f, 0xde, 0x82, 0x03, 0x74, 0x0e, 0xa1, 0x6b, 0xeb, 0xcd, 0x86,
	0x7d, 0xf4, 0x74, 0x9d, 0x1d, 0xf6, 0x9b, 0xc9, 0x0e, 0xa0, 0x7f, 0x66, 0x14, 0xa3, 0x93, 0xe6,
	0x6e, 0xef, 0xbc, 0x00, 0x2f, 0x01, 0x1a, 0xc1, 0x47, 0xb7, 0x85, 0x9a, 0xbb, 0x74, 0x72, 0x1f,
	0xa0, 0x37, 0x4c, 0x35, 0x53, 0x06, 0xed, 0xac, 0xbf, 0xf9, 0x05, 0x0b, 0x8d, 0x50, 0x41, 0xbb,
	0xac, 0x1f, 0x4d, 0xa4, 0x99, 0x75, 0x77, 0xae, 0xe7, 0xfb, 0xa0, 0x4c, 0xc8, 0x2d, 0x9c, 0x27,
	0x7c, 0xb7, 0x4c, 0x58, 0x53, 0xcf, 0xf9, 0x38, 0x84, 0xf7, 0x47, 0x59, 0x62, 0x78, 0x01, 0xef,
	0x56, 0xc3, 0x7a, 0x93, 0xce, 0x13, 0x9c, 0xcb, 0x88, 0x1a, 0xd6, 0x30, 0x41, 0x66, 0xe1, 0x95,
	0x04, 0x35, 0xf5, 0x56, 0x13, 0x14, 0x70, 0x8d, 0x04, 0xc7, 0xd0, 0x3b, 0x65, 0x13, 0x31, 0x65,
	0x95, 0xb3, 0xb0, 0xde, 0xbf, 0x57, 0xba, 0x6f, 0x1f, 0xf8, 0x44, 0x59, 0xd0, 0xcd, 0xc2, 0xe1,
	0xd2, 0x7f, 0x0d, 0x35, 0xe7, 0xfd, 0x4d, 0xe1, 0xbd, 0x00, 0x1f, 0x6f, 0x82, 0x55, 0xbe, 0x3f,
	0xc1, 0x7b, 0x03, 0x66, 0x5c, 0x4b, 0xe5, 0xcf, 0x82, 0xb2, 0xd6, 0xa3, 0xe1, 0xd7, 0x4c, 0xe2,
	0x11, 0x33, 0xd4, 0xbd, 0xc3, 0xad, 0x71, 0x16, 0xb6, 0xdf, 0x45, 0x18, 0xc0, 0x87, 0x2e, 0x42,
	0x73, 0x61, 0x9b, 0x28, 0x68, 0x5d, 0xcf, 0xf7, 0xb7, 0x7a, 0xa3, 0xdf, 0x8b, 0x0e, 0xf8, 0xb3,
	0xe8, 0x80, 0xf9, 0xa2, 0x03, 0x60, 0x5b, 0xa8, 0x18, 0x4f, 0x23, 0x4a, 0x35, 0xce, 0xf7, 0x58,
	0x6f, 0x3b, 0xdf, 0x44, 0xef, 0x25, 0xff, 0x08, 0x3e, 0x3f, 0x8f, 0xb9, 0xb9, 0xcc, 0xc6, 0x38,
	0x14, 0x13, 0x62, 0x3b, 0x48, 0xde, 0x91, 0xaf, 0x34, 0x4d, 0x62, 0x25, 0x43, 0x7b, 0x1c, 0x7b,
	0x76, 0x89, 0xbd, 0xbe, 0x09, 0x00, 0x00, 0xff, 0xff, 0x98, 0x5a, 0xf3, 0xfc, 0x13, 0x05, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ValdClient is the client API for Vald service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ValdClient interface {
	Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_ID, error)
	Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error)
	SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)
	StreamSearch(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamSearchClient, error)
	StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamSearchByIDClient, error)
	Insert(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamInsert(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamInsertClient, error)
	MultiInsert(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error)
	Update(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamUpdateClient, error)
	MultiUpdate(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error)
	Remove(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamRemove(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamRemoveClient, error)
	MultiRemove(ctx context.Context, in *payload.Object_IDs, opts ...grpc.CallOption) (*payload.Empty, error)
	GetObject(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Backup_MetaVector, error)
	StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamGetObjectClient, error)
}

type valdClient struct {
	cc *grpc.ClientConn
}

func NewValdClient(cc *grpc.ClientConn) ValdClient {
	return &valdClient{cc}
}

func (c *valdClient) Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_ID, error) {
	out := new(payload.Object_ID)
	err := c.cc.Invoke(ctx, "/vald.Vald/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *valdClient) Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/vald.Vald/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *valdClient) SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/vald.Vald/SearchByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *valdClient) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamSearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Vald_serviceDesc.Streams[0], "/vald.Vald/StreamSearch", opts...)
	if err != nil {
		return nil, err
	}
	x := &valdStreamSearchClient{stream}
	return x, nil
}

type Vald_StreamSearchClient interface {
	Send(*payload.Search_Request) error
	Recv() (*payload.Search_Response, error)
	grpc.ClientStream
}

type valdStreamSearchClient struct {
	grpc.ClientStream
}

func (x *valdStreamSearchClient) Send(m *payload.Search_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *valdStreamSearchClient) Recv() (*payload.Search_Response, error) {
	m := new(payload.Search_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *valdClient) StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamSearchByIDClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Vald_serviceDesc.Streams[1], "/vald.Vald/StreamSearchByID", opts...)
	if err != nil {
		return nil, err
	}
	x := &valdStreamSearchByIDClient{stream}
	return x, nil
}

type Vald_StreamSearchByIDClient interface {
	Send(*payload.Search_IDRequest) error
	Recv() (*payload.Search_Response, error)
	grpc.ClientStream
}

type valdStreamSearchByIDClient struct {
	grpc.ClientStream
}

func (x *valdStreamSearchByIDClient) Send(m *payload.Search_IDRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *valdStreamSearchByIDClient) Recv() (*payload.Search_Response, error) {
	m := new(payload.Search_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *valdClient) Insert(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/vald.Vald/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *valdClient) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamInsertClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Vald_serviceDesc.Streams[2], "/vald.Vald/StreamInsert", opts...)
	if err != nil {
		return nil, err
	}
	x := &valdStreamInsertClient{stream}
	return x, nil
}

type Vald_StreamInsertClient interface {
	Send(*payload.Object_Vector) error
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type valdStreamInsertClient struct {
	grpc.ClientStream
}

func (x *valdStreamInsertClient) Send(m *payload.Object_Vector) error {
	return x.ClientStream.SendMsg(m)
}

func (x *valdStreamInsertClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *valdClient) MultiInsert(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/vald.Vald/MultiInsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *valdClient) Update(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/vald.Vald/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *valdClient) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamUpdateClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Vald_serviceDesc.Streams[3], "/vald.Vald/StreamUpdate", opts...)
	if err != nil {
		return nil, err
	}
	x := &valdStreamUpdateClient{stream}
	return x, nil
}

type Vald_StreamUpdateClient interface {
	Send(*payload.Object_Vector) error
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type valdStreamUpdateClient struct {
	grpc.ClientStream
}

func (x *valdStreamUpdateClient) Send(m *payload.Object_Vector) error {
	return x.ClientStream.SendMsg(m)
}

func (x *valdStreamUpdateClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *valdClient) MultiUpdate(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/vald.Vald/MultiUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *valdClient) Remove(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/vald.Vald/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *valdClient) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamRemoveClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Vald_serviceDesc.Streams[4], "/vald.Vald/StreamRemove", opts...)
	if err != nil {
		return nil, err
	}
	x := &valdStreamRemoveClient{stream}
	return x, nil
}

type Vald_StreamRemoveClient interface {
	Send(*payload.Object_ID) error
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type valdStreamRemoveClient struct {
	grpc.ClientStream
}

func (x *valdStreamRemoveClient) Send(m *payload.Object_ID) error {
	return x.ClientStream.SendMsg(m)
}

func (x *valdStreamRemoveClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *valdClient) MultiRemove(ctx context.Context, in *payload.Object_IDs, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/vald.Vald/MultiRemove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *valdClient) GetObject(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Backup_MetaVector, error) {
	out := new(payload.Backup_MetaVector)
	err := c.cc.Invoke(ctx, "/vald.Vald/GetObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *valdClient) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (Vald_StreamGetObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Vald_serviceDesc.Streams[5], "/vald.Vald/StreamGetObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &valdStreamGetObjectClient{stream}
	return x, nil
}

type Vald_StreamGetObjectClient interface {
	Send(*payload.Object_ID) error
	Recv() (*payload.Backup_MetaVector, error)
	grpc.ClientStream
}

type valdStreamGetObjectClient struct {
	grpc.ClientStream
}

func (x *valdStreamGetObjectClient) Send(m *payload.Object_ID) error {
	return x.ClientStream.SendMsg(m)
}

func (x *valdStreamGetObjectClient) Recv() (*payload.Backup_MetaVector, error) {
	m := new(payload.Backup_MetaVector)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ValdServer is the server API for Vald service.
type ValdServer interface {
	Exists(context.Context, *payload.Object_ID) (*payload.Object_ID, error)
	Search(context.Context, *payload.Search_Request) (*payload.Search_Response, error)
	SearchByID(context.Context, *payload.Search_IDRequest) (*payload.Search_Response, error)
	StreamSearch(Vald_StreamSearchServer) error
	StreamSearchByID(Vald_StreamSearchByIDServer) error
	Insert(context.Context, *payload.Object_Vector) (*payload.Empty, error)
	StreamInsert(Vald_StreamInsertServer) error
	MultiInsert(context.Context, *payload.Object_Vectors) (*payload.Empty, error)
	Update(context.Context, *payload.Object_Vector) (*payload.Empty, error)
	StreamUpdate(Vald_StreamUpdateServer) error
	MultiUpdate(context.Context, *payload.Object_Vectors) (*payload.Empty, error)
	Remove(context.Context, *payload.Object_ID) (*payload.Empty, error)
	StreamRemove(Vald_StreamRemoveServer) error
	MultiRemove(context.Context, *payload.Object_IDs) (*payload.Empty, error)
	GetObject(context.Context, *payload.Object_ID) (*payload.Backup_MetaVector, error)
	StreamGetObject(Vald_StreamGetObjectServer) error
}

// UnimplementedValdServer can be embedded to have forward compatible implementations.
type UnimplementedValdServer struct {
}

func (*UnimplementedValdServer) Exists(ctx context.Context, req *payload.Object_ID) (*payload.Object_ID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (*UnimplementedValdServer) Search(ctx context.Context, req *payload.Search_Request) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (*UnimplementedValdServer) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchByID not implemented")
}
func (*UnimplementedValdServer) StreamSearch(srv Vald_StreamSearchServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearch not implemented")
}
func (*UnimplementedValdServer) StreamSearchByID(srv Vald_StreamSearchByIDServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearchByID not implemented")
}
func (*UnimplementedValdServer) Insert(ctx context.Context, req *payload.Object_Vector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (*UnimplementedValdServer) StreamInsert(srv Vald_StreamInsertServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamInsert not implemented")
}
func (*UnimplementedValdServer) MultiInsert(ctx context.Context, req *payload.Object_Vectors) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiInsert not implemented")
}
func (*UnimplementedValdServer) Update(ctx context.Context, req *payload.Object_Vector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedValdServer) StreamUpdate(srv Vald_StreamUpdateServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpdate not implemented")
}
func (*UnimplementedValdServer) MultiUpdate(ctx context.Context, req *payload.Object_Vectors) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpdate not implemented")
}
func (*UnimplementedValdServer) Remove(ctx context.Context, req *payload.Object_ID) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (*UnimplementedValdServer) StreamRemove(srv Vald_StreamRemoveServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamRemove not implemented")
}
func (*UnimplementedValdServer) MultiRemove(ctx context.Context, req *payload.Object_IDs) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiRemove not implemented")
}
func (*UnimplementedValdServer) GetObject(ctx context.Context, req *payload.Object_ID) (*payload.Backup_MetaVector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetObject not implemented")
}
func (*UnimplementedValdServer) StreamGetObject(srv Vald_StreamGetObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamGetObject not implemented")
}

func RegisterValdServer(s *grpc.Server, srv ValdServer) {
	s.RegisterService(&_Vald_serviceDesc, srv)
}

func _Vald_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValdServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Vald/Exists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValdServer).Exists(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vald_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Search_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValdServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Vald/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValdServer).Search(ctx, req.(*payload.Search_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vald_SearchByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Search_IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValdServer).SearchByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Vald/SearchByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValdServer).SearchByID(ctx, req.(*payload.Search_IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vald_StreamSearch_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ValdServer).StreamSearch(&valdStreamSearchServer{stream})
}

type Vald_StreamSearchServer interface {
	Send(*payload.Search_Response) error
	Recv() (*payload.Search_Request, error)
	grpc.ServerStream
}

type valdStreamSearchServer struct {
	grpc.ServerStream
}

func (x *valdStreamSearchServer) Send(m *payload.Search_Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *valdStreamSearchServer) Recv() (*payload.Search_Request, error) {
	m := new(payload.Search_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Vald_StreamSearchByID_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ValdServer).StreamSearchByID(&valdStreamSearchByIDServer{stream})
}

type Vald_StreamSearchByIDServer interface {
	Send(*payload.Search_Response) error
	Recv() (*payload.Search_IDRequest, error)
	grpc.ServerStream
}

type valdStreamSearchByIDServer struct {
	grpc.ServerStream
}

func (x *valdStreamSearchByIDServer) Send(m *payload.Search_Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *valdStreamSearchByIDServer) Recv() (*payload.Search_IDRequest, error) {
	m := new(payload.Search_IDRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Vald_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValdServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Vald/Insert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValdServer).Insert(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vald_StreamInsert_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ValdServer).StreamInsert(&valdStreamInsertServer{stream})
}

type Vald_StreamInsertServer interface {
	Send(*payload.Empty) error
	Recv() (*payload.Object_Vector, error)
	grpc.ServerStream
}

type valdStreamInsertServer struct {
	grpc.ServerStream
}

func (x *valdStreamInsertServer) Send(m *payload.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *valdStreamInsertServer) Recv() (*payload.Object_Vector, error) {
	m := new(payload.Object_Vector)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Vald_MultiInsert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vectors)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValdServer).MultiInsert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Vald/MultiInsert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValdServer).MultiInsert(ctx, req.(*payload.Object_Vectors))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vald_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValdServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Vald/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValdServer).Update(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vald_StreamUpdate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ValdServer).StreamUpdate(&valdStreamUpdateServer{stream})
}

type Vald_StreamUpdateServer interface {
	Send(*payload.Empty) error
	Recv() (*payload.Object_Vector, error)
	grpc.ServerStream
}

type valdStreamUpdateServer struct {
	grpc.ServerStream
}

func (x *valdStreamUpdateServer) Send(m *payload.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *valdStreamUpdateServer) Recv() (*payload.Object_Vector, error) {
	m := new(payload.Object_Vector)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Vald_MultiUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vectors)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValdServer).MultiUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Vald/MultiUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValdServer).MultiUpdate(ctx, req.(*payload.Object_Vectors))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vald_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValdServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Vald/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValdServer).Remove(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vald_StreamRemove_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ValdServer).StreamRemove(&valdStreamRemoveServer{stream})
}

type Vald_StreamRemoveServer interface {
	Send(*payload.Empty) error
	Recv() (*payload.Object_ID, error)
	grpc.ServerStream
}

type valdStreamRemoveServer struct {
	grpc.ServerStream
}

func (x *valdStreamRemoveServer) Send(m *payload.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *valdStreamRemoveServer) Recv() (*payload.Object_ID, error) {
	m := new(payload.Object_ID)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Vald_MultiRemove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_IDs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValdServer).MultiRemove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Vald/MultiRemove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValdServer).MultiRemove(ctx, req.(*payload.Object_IDs))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vald_GetObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValdServer).GetObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.Vald/GetObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValdServer).GetObject(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vald_StreamGetObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ValdServer).StreamGetObject(&valdStreamGetObjectServer{stream})
}

type Vald_StreamGetObjectServer interface {
	Send(*payload.Backup_MetaVector) error
	Recv() (*payload.Object_ID, error)
	grpc.ServerStream
}

type valdStreamGetObjectServer struct {
	grpc.ServerStream
}

func (x *valdStreamGetObjectServer) Send(m *payload.Backup_MetaVector) error {
	return x.ServerStream.SendMsg(m)
}

func (x *valdStreamGetObjectServer) Recv() (*payload.Object_ID, error) {
	m := new(payload.Object_ID)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Vald_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vald.Vald",
	HandlerType: (*ValdServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Exists",
			Handler:    _Vald_Exists_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _Vald_Search_Handler,
		},
		{
			MethodName: "SearchByID",
			Handler:    _Vald_SearchByID_Handler,
		},
		{
			MethodName: "Insert",
			Handler:    _Vald_Insert_Handler,
		},
		{
			MethodName: "MultiInsert",
			Handler:    _Vald_MultiInsert_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Vald_Update_Handler,
		},
		{
			MethodName: "MultiUpdate",
			Handler:    _Vald_MultiUpdate_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Vald_Remove_Handler,
		},
		{
			MethodName: "MultiRemove",
			Handler:    _Vald_MultiRemove_Handler,
		},
		{
			MethodName: "GetObject",
			Handler:    _Vald_GetObject_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamSearch",
			Handler:       _Vald_StreamSearch_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamSearchByID",
			Handler:       _Vald_StreamSearchByID_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamInsert",
			Handler:       _Vald_StreamInsert_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamUpdate",
			Handler:       _Vald_StreamUpdate_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamRemove",
			Handler:       _Vald_StreamRemove_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamGetObject",
			Handler:       _Vald_StreamGetObject_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "vald.proto",
}
