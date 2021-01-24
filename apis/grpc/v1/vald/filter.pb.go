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

package vald

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

func init() { proto.RegisterFile("apis/proto/v1/vald/filter.proto", fileDescriptor_a46f8d8eee988c86) }
func init() {
	golang_proto.RegisterFile("apis/proto/v1/vald/filter.proto", fileDescriptor_a46f8d8eee988c86)
}

var fileDescriptor_a46f8d8eee988c86 = []byte{
	// 479 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x55, 0xc1, 0x8a, 0xd4, 0x30,
	0x18, 0x36, 0x1e, 0x56, 0x08, 0x8b, 0x60, 0x3d, 0x88, 0x75, 0xa9, 0x6e, 0x65, 0x50, 0x04, 0x13,
	0x47, 0x0f, 0x82, 0xc7, 0x3d, 0x08, 0x82, 0xa2, 0xec, 0xa2, 0x07, 0x2f, 0x92, 0xb6, 0xd9, 0x36,
	0xd2, 0x69, 0x62, 0x93, 0x16, 0x04, 0x4f, 0xbe, 0x82, 0xaf, 0xe1, 0x43, 0x78, 0xdc, 0xa3, 0xe0,
	0x0b, 0xc8, 0xac, 0x0f, 0x22, 0xcd, 0xbf, 0x1d, 0x9a, 0x31, 0x1d, 0xe9, 0xcc, 0x69, 0x32, 0xfd,
	0xbf, 0xff, 0xfb, 0xfa, 0x7d, 0x7f, 0xf8, 0x8b, 0x6f, 0x33, 0x25, 0x34, 0x55, 0xb5, 0x34, 0x92,
	0xb6, 0x73, 0xda, 0xb2, 0x32, 0xa3, 0xa7, 0xa2, 0x34, 0xbc, 0x26, 0xf6, 0x61, 0x70, 0xa5, 0x7b,
	0x44, 0xda, 0x79, 0x78, 0xd7, 0x45, 0x2a, 0xf6, 0xb9, 0x94, 0x2c, 0xeb, 0x7f, 0x01, 0x1d, 0x3e,
	0xcc, 0x85, 0x29, 0x9a, 0x84, 0xa4, 0x72, 0x41, 0x73, 0x99, 0x4b, 0xc0, 0x27, 0xcd, 0xa9, 0xfd,
	0x07, 0xcd, 0xdd, 0xe9, 0x02, 0xfe, 0x74, 0x1d, 0x9e, 0x4b, 0x99, 0x97, 0xdc, 0x2a, 0xc1, 0x91,
	0x32, 0x25, 0x28, 0xab, 0x2a, 0x69, 0x98, 0x11, 0xb2, 0xd2, 0xd0, 0xf8, 0xf8, 0x3b, 0xc6, 0x7b,
	0xcf, 0xed, 0x6b, 0x06, 0x05, 0xde, 0x3f, 0xe1, 0xac, 0x4e, 0x8b, 0xd7, 0xc9, 0x47, 0x9e, 0x9a,
	0xe0, 0x0e, 0xe9, 0x5f, 0xa9, 0x9d, 0x13, 0xa8, 0x10, 0x28, 0x1d, 0xf3, 0x4f, 0x0d, 0xd7, 0x26,
	0xbc, 0xe5, 0x41, 0x1c, 0x73, 0xad, 0x64, 0xa5, 0x79, 0x7c, 0xf3, 0xeb, 0xaf, 0x3f, 0xdf, 0x2e,
	0x5f, 0x8f, 0xaf, 0x52, 0x6d, 0x2b, 0x54, 0xda, 0xde, 0x67, 0xe8, 0x41, 0xf0, 0x05, 0x5f, 0x7b,
	0xd5, 0x94, 0x46, 0x38, 0x72, 0x33, 0x0f, 0x99, 0x45, 0xb9, 0x9a, 0x07, 0x1b, 0x34, 0x75, 0x1c,
	0x5b, 0xd1, 0x83, 0xf8, 0x86, 0x2b, 0x4a, 0x17, 0x1d, 0x91, 0x2a, 0x79, 0xa7, 0xfe, 0x01, 0x07,
	0x27, 0xa6, 0xe6, 0x6c, 0x31, 0xd1, 0xed, 0xa1, 0x07, 0x01, 0x44, 0x2b, 0xcf, 0x97, 0xee, 0xa3,
	0x47, 0xa8, 0x0b, 0xf2, 0x45, 0xa5, 0x79, 0x6d, 0x7c, 0xd4, 0x50, 0xd9, 0x14, 0x24, 0x94, 0xc8,
	0x4b, 0x99, 0xda, 0x49, 0x0d, 0x82, 0x14, 0xb6, 0x77, 0x10, 0xe4, 0xca, 0xca, 0x44, 0xbd, 0x43,
	0x8f, 0x1e, 0x10, 0xad, 0x54, 0xc1, 0x4a, 0x3f, 0x29, 0x87, 0x7f, 0xe6, 0xe1, 0xff, 0xdf, 0xa4,
	0xd6, 0x4c, 0x0d, 0x27, 0xe5, 0xb8, 0x72, 0x26, 0x55, 0xe0, 0xfd, 0xb7, 0x2a, 0x63, 0x86, 0xfb,
	0x8c, 0x41, 0x65, 0xbb, 0x20, 0x1b, 0xdb, 0xeb, 0x0b, 0x72, 0xa2, 0xde, 0xe4, 0x20, 0x1d, 0xfe,
	0x99, 0x87, 0x7f, 0x87, 0x20, 0x1d, 0x57, 0xff, 0x06, 0x39, 0x76, 0x43, 0xa0, 0xb2, 0x6d, 0x90,
	0x63, 0x37, 0x72, 0xa2, 0xde, 0x16, 0x41, 0x8e, 0xdd, 0xc8, 0x0b, 0xfe, 0x9d, 0x82, 0x1c, 0xb9,
	0x91, 0x47, 0xf9, 0xd9, 0x32, 0x42, 0x3f, 0x97, 0x11, 0xfa, 0xbd, 0x8c, 0xd0, 0x8f, 0xf3, 0x08,
	0x9d, 0x9d, 0x47, 0x08, 0x87, 0xb2, 0xce, 0x49, 0x9b, 0x31, 0xa6, 0x89, 0x5d, 0xf2, 0x4c, 0x89,
	0x4e, 0xa1, 0x3b, 0x1f, 0xe1, 0x77, 0xac, 0xcc, 0x60, 0xc3, 0xbe, 0x41, 0xef, 0xef, 0x0d, 0xf6,
	0xb4, 0x6d, 0x80, 0x0f, 0x05, 0xec, 0xe9, 0x5a, 0xa5, 0xfd, 0xa7, 0x23, 0xd9, 0xb3, 0xeb, 0xf9,
	0xc9, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd7, 0xe9, 0xe7, 0xb6, 0x57, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FilterClient is the client API for Filter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FilterClient interface {
	SearchObject(ctx context.Context, in *payload.Search_ObjectRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)
	MultiSearchObject(ctx context.Context, in *payload.Search_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	StreamSearchObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamSearchObjectClient, error)
	InsertObject(ctx context.Context, in *payload.Insert_ObjectRequest, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamInsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamInsertObjectClient, error)
	MultiInsertObject(ctx context.Context, in *payload.Insert_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
	UpdateObject(ctx context.Context, in *payload.Update_ObjectRequest, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamUpdateObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpdateObjectClient, error)
	MultiUpdateObject(ctx context.Context, in *payload.Update_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
	UpsertObject(ctx context.Context, in *payload.Upsert_ObjectRequest, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamUpsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpsertObjectClient, error)
	MultiUpsertObject(ctx context.Context, in *payload.Upsert_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type filterClient struct {
	cc *grpc.ClientConn
}

func NewFilterClient(cc *grpc.ClientConn) FilterClient {
	return &filterClient{cc}
}

func (c *filterClient) SearchObject(ctx context.Context, in *payload.Search_ObjectRequest, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/vald.v1.Filter/SearchObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) MultiSearchObject(ctx context.Context, in *payload.Search_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error) {
	out := new(payload.Search_Responses)
	err := c.cc.Invoke(ctx, "/vald.v1.Filter/MultiSearchObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) StreamSearchObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamSearchObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Filter_serviceDesc.Streams[0], "/vald.v1.Filter/StreamSearchObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &filterStreamSearchObjectClient{stream}
	return x, nil
}

type Filter_StreamSearchObjectClient interface {
	Send(*payload.Search_ObjectRequest) error
	Recv() (*payload.Search_StreamResponse, error)
	grpc.ClientStream
}

type filterStreamSearchObjectClient struct {
	grpc.ClientStream
}

func (x *filterStreamSearchObjectClient) Send(m *payload.Search_ObjectRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamSearchObjectClient) Recv() (*payload.Search_StreamResponse, error) {
	m := new(payload.Search_StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) InsertObject(ctx context.Context, in *payload.Insert_ObjectRequest, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/vald.v1.Filter/InsertObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) StreamInsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamInsertObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Filter_serviceDesc.Streams[1], "/vald.v1.Filter/StreamInsertObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &filterStreamInsertObjectClient{stream}
	return x, nil
}

type Filter_StreamInsertObjectClient interface {
	Send(*payload.Insert_ObjectRequest) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type filterStreamInsertObjectClient struct {
	grpc.ClientStream
}

func (x *filterStreamInsertObjectClient) Send(m *payload.Insert_ObjectRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamInsertObjectClient) Recv() (*payload.Object_StreamLocation, error) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) MultiInsertObject(ctx context.Context, in *payload.Insert_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.v1.Filter/MultiInsertObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) UpdateObject(ctx context.Context, in *payload.Update_ObjectRequest, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/vald.v1.Filter/UpdateObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) StreamUpdateObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpdateObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Filter_serviceDesc.Streams[2], "/vald.v1.Filter/StreamUpdateObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &filterStreamUpdateObjectClient{stream}
	return x, nil
}

type Filter_StreamUpdateObjectClient interface {
	Send(*payload.Update_ObjectRequest) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type filterStreamUpdateObjectClient struct {
	grpc.ClientStream
}

func (x *filterStreamUpdateObjectClient) Send(m *payload.Update_ObjectRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamUpdateObjectClient) Recv() (*payload.Object_StreamLocation, error) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) MultiUpdateObject(ctx context.Context, in *payload.Update_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.v1.Filter/MultiUpdateObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) UpsertObject(ctx context.Context, in *payload.Upsert_ObjectRequest, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/vald.v1.Filter/UpsertObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) StreamUpsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpsertObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Filter_serviceDesc.Streams[3], "/vald.v1.Filter/StreamUpsertObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &filterStreamUpsertObjectClient{stream}
	return x, nil
}

type Filter_StreamUpsertObjectClient interface {
	Send(*payload.Upsert_ObjectRequest) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type filterStreamUpsertObjectClient struct {
	grpc.ClientStream
}

func (x *filterStreamUpsertObjectClient) Send(m *payload.Upsert_ObjectRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamUpsertObjectClient) Recv() (*payload.Object_StreamLocation, error) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) MultiUpsertObject(ctx context.Context, in *payload.Upsert_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/vald.v1.Filter/MultiUpsertObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilterServer is the server API for Filter service.
type FilterServer interface {
	SearchObject(context.Context, *payload.Search_ObjectRequest) (*payload.Search_Response, error)
	MultiSearchObject(context.Context, *payload.Search_MultiObjectRequest) (*payload.Search_Responses, error)
	StreamSearchObject(Filter_StreamSearchObjectServer) error
	InsertObject(context.Context, *payload.Insert_ObjectRequest) (*payload.Object_Location, error)
	StreamInsertObject(Filter_StreamInsertObjectServer) error
	MultiInsertObject(context.Context, *payload.Insert_MultiObjectRequest) (*payload.Object_Locations, error)
	UpdateObject(context.Context, *payload.Update_ObjectRequest) (*payload.Object_Location, error)
	StreamUpdateObject(Filter_StreamUpdateObjectServer) error
	MultiUpdateObject(context.Context, *payload.Update_MultiObjectRequest) (*payload.Object_Locations, error)
	UpsertObject(context.Context, *payload.Upsert_ObjectRequest) (*payload.Object_Location, error)
	StreamUpsertObject(Filter_StreamUpsertObjectServer) error
	MultiUpsertObject(context.Context, *payload.Upsert_MultiObjectRequest) (*payload.Object_Locations, error)
}

// UnimplementedFilterServer can be embedded to have forward compatible implementations.
type UnimplementedFilterServer struct {
}

func (*UnimplementedFilterServer) SearchObject(ctx context.Context, req *payload.Search_ObjectRequest) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchObject not implemented")
}
func (*UnimplementedFilterServer) MultiSearchObject(ctx context.Context, req *payload.Search_MultiObjectRequest) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiSearchObject not implemented")
}
func (*UnimplementedFilterServer) StreamSearchObject(srv Filter_StreamSearchObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearchObject not implemented")
}
func (*UnimplementedFilterServer) InsertObject(ctx context.Context, req *payload.Insert_ObjectRequest) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertObject not implemented")
}
func (*UnimplementedFilterServer) StreamInsertObject(srv Filter_StreamInsertObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamInsertObject not implemented")
}
func (*UnimplementedFilterServer) MultiInsertObject(ctx context.Context, req *payload.Insert_MultiObjectRequest) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiInsertObject not implemented")
}
func (*UnimplementedFilterServer) UpdateObject(ctx context.Context, req *payload.Update_ObjectRequest) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateObject not implemented")
}
func (*UnimplementedFilterServer) StreamUpdateObject(srv Filter_StreamUpdateObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpdateObject not implemented")
}
func (*UnimplementedFilterServer) MultiUpdateObject(ctx context.Context, req *payload.Update_MultiObjectRequest) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpdateObject not implemented")
}
func (*UnimplementedFilterServer) UpsertObject(ctx context.Context, req *payload.Upsert_ObjectRequest) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertObject not implemented")
}
func (*UnimplementedFilterServer) StreamUpsertObject(srv Filter_StreamUpsertObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpsertObject not implemented")
}
func (*UnimplementedFilterServer) MultiUpsertObject(ctx context.Context, req *payload.Upsert_MultiObjectRequest) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpsertObject not implemented")
}

func RegisterFilterServer(s *grpc.Server, srv FilterServer) {
	s.RegisterService(&_Filter_serviceDesc, srv)
}

func _Filter_SearchObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Search_ObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).SearchObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Filter/SearchObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).SearchObject(ctx, req.(*payload.Search_ObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_MultiSearchObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Search_MultiObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).MultiSearchObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Filter/MultiSearchObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).MultiSearchObject(ctx, req.(*payload.Search_MultiObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_StreamSearchObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilterServer).StreamSearchObject(&filterStreamSearchObjectServer{stream})
}

type Filter_StreamSearchObjectServer interface {
	Send(*payload.Search_StreamResponse) error
	Recv() (*payload.Search_ObjectRequest, error)
	grpc.ServerStream
}

type filterStreamSearchObjectServer struct {
	grpc.ServerStream
}

func (x *filterStreamSearchObjectServer) Send(m *payload.Search_StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *filterStreamSearchObjectServer) Recv() (*payload.Search_ObjectRequest, error) {
	m := new(payload.Search_ObjectRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Filter_InsertObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Insert_ObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).InsertObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Filter/InsertObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).InsertObject(ctx, req.(*payload.Insert_ObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_StreamInsertObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilterServer).StreamInsertObject(&filterStreamInsertObjectServer{stream})
}

type Filter_StreamInsertObjectServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Insert_ObjectRequest, error)
	grpc.ServerStream
}

type filterStreamInsertObjectServer struct {
	grpc.ServerStream
}

func (x *filterStreamInsertObjectServer) Send(m *payload.Object_StreamLocation) error {
	return x.ServerStream.SendMsg(m)
}

func (x *filterStreamInsertObjectServer) Recv() (*payload.Insert_ObjectRequest, error) {
	m := new(payload.Insert_ObjectRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Filter_MultiInsertObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Insert_MultiObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).MultiInsertObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Filter/MultiInsertObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).MultiInsertObject(ctx, req.(*payload.Insert_MultiObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_UpdateObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Update_ObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).UpdateObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Filter/UpdateObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).UpdateObject(ctx, req.(*payload.Update_ObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_StreamUpdateObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilterServer).StreamUpdateObject(&filterStreamUpdateObjectServer{stream})
}

type Filter_StreamUpdateObjectServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Update_ObjectRequest, error)
	grpc.ServerStream
}

type filterStreamUpdateObjectServer struct {
	grpc.ServerStream
}

func (x *filterStreamUpdateObjectServer) Send(m *payload.Object_StreamLocation) error {
	return x.ServerStream.SendMsg(m)
}

func (x *filterStreamUpdateObjectServer) Recv() (*payload.Update_ObjectRequest, error) {
	m := new(payload.Update_ObjectRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Filter_MultiUpdateObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Update_MultiObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).MultiUpdateObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Filter/MultiUpdateObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).MultiUpdateObject(ctx, req.(*payload.Update_MultiObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_UpsertObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Upsert_ObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).UpsertObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Filter/UpsertObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).UpsertObject(ctx, req.(*payload.Upsert_ObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_StreamUpsertObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilterServer).StreamUpsertObject(&filterStreamUpsertObjectServer{stream})
}

type Filter_StreamUpsertObjectServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Upsert_ObjectRequest, error)
	grpc.ServerStream
}

type filterStreamUpsertObjectServer struct {
	grpc.ServerStream
}

func (x *filterStreamUpsertObjectServer) Send(m *payload.Object_StreamLocation) error {
	return x.ServerStream.SendMsg(m)
}

func (x *filterStreamUpsertObjectServer) Recv() (*payload.Upsert_ObjectRequest, error) {
	m := new(payload.Upsert_ObjectRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Filter_MultiUpsertObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Upsert_MultiObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).MultiUpsertObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Filter/MultiUpsertObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).MultiUpsertObject(ctx, req.(*payload.Upsert_MultiObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Filter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vald.v1.Filter",
	HandlerType: (*FilterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchObject",
			Handler:    _Filter_SearchObject_Handler,
		},
		{
			MethodName: "MultiSearchObject",
			Handler:    _Filter_MultiSearchObject_Handler,
		},
		{
			MethodName: "InsertObject",
			Handler:    _Filter_InsertObject_Handler,
		},
		{
			MethodName: "MultiInsertObject",
			Handler:    _Filter_MultiInsertObject_Handler,
		},
		{
			MethodName: "UpdateObject",
			Handler:    _Filter_UpdateObject_Handler,
		},
		{
			MethodName: "MultiUpdateObject",
			Handler:    _Filter_MultiUpdateObject_Handler,
		},
		{
			MethodName: "UpsertObject",
			Handler:    _Filter_UpsertObject_Handler,
		},
		{
			MethodName: "MultiUpsertObject",
			Handler:    _Filter_MultiUpsertObject_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamSearchObject",
			Handler:       _Filter_StreamSearchObject_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamInsertObject",
			Handler:       _Filter_StreamInsertObject_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamUpdateObject",
			Handler:       _Filter_StreamUpdateObject_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamUpsertObject",
			Handler:       _Filter_StreamUpsertObject_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "apis/proto/v1/vald/filter.proto",
}
