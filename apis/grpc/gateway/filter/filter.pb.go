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

package filter

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
	proto.RegisterFile("apis/proto/gateway/filter/filter.proto", fileDescriptor_098434da7999ba4a)
}

var fileDescriptor_098434da7999ba4a = []byte{
	// 378 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0x41, 0x4b, 0x2b, 0x31,
	0x10, 0xc7, 0x9b, 0x77, 0xe8, 0x83, 0x50, 0x1e, 0x74, 0xdf, 0x3b, 0xbc, 0x16, 0x5b, 0xa4, 0x07,
	0x11, 0x0f, 0x89, 0xe8, 0xcd, 0x63, 0x0f, 0x15, 0x45, 0x51, 0x5a, 0x14, 0xf4, 0x36, 0xbb, 0x1b,
	0xb7, 0x91, 0x74, 0x13, 0x93, 0x6c, 0xa5, 0x57, 0x3f, 0x80, 0x17, 0xbf, 0x94, 0x47, 0xc1, 0x2f,
	0x20, 0xc5, 0x0f, 0x22, 0x4d, 0xb6, 0xb2, 0x16, 0x8b, 0x85, 0xed, 0x29, 0xec, 0xfc, 0x33, 0xbf,
	0x99, 0xff, 0x0e, 0x19, 0xbc, 0x05, 0x8a, 0x1b, 0xaa, 0xb4, 0xb4, 0x92, 0x26, 0x60, 0xd9, 0x3d,
	0x4c, 0xe8, 0x0d, 0x17, 0x96, 0xe9, 0xfc, 0x20, 0x4e, 0x0b, 0xaa, 0xfe, 0xab, 0xb9, 0x59, 0xb8,
	0xaf, 0x60, 0x22, 0x24, 0xc4, 0xf3, 0xd3, 0xdf, 0x6c, 0x6e, 0x24, 0x52, 0x26, 0x82, 0x51, 0x50,
	0x9c, 0x42, 0x9a, 0x4a, 0x0b, 0x96, 0xcb, 0xd4, 0x78, 0x75, 0xef, 0xf1, 0x37, 0xae, 0xf6, 0x1c,
	0x2a, 0x08, 0x71, 0x6d, 0xc0, 0x40, 0x47, 0xc3, 0xb3, 0xf0, 0x96, 0x45, 0x36, 0x68, 0x91, 0x39,
	0xc8, 0x87, 0x89, 0x8f, 0xf7, 0xd9, 0x5d, 0xc6, 0x8c, 0x6d, 0xfe, 0x5f, 0x94, 0xfb, 0xcc, 0x28,
	0x99, 0x1a, 0xd6, 0x69, 0x3c, 0xbc, 0xbe, 0x3f, 0xfd, 0xfa, 0xdb, 0xf9, 0x43, 0x8d, 0x53, 0xa8,
	0x74, 0x89, 0x07, 0x68, 0x27, 0x18, 0xe0, 0x60, 0x60, 0x35, 0x83, 0xd1, 0x7a, 0x2a, 0x55, 0xb6,
	0xd1, 0x2e, 0x0a, 0xae, 0x70, 0xed, 0x28, 0x35, 0x4c, 0xdb, 0x1c, 0xf7, 0xef, 0xf3, 0xbe, 0x0f,
	0x90, 0xae, 0x90, 0x61, 0x81, 0x92, 0x47, 0x4f, 0x64, 0xe4, 0x7e, 0x45, 0xa1, 0x5f, 0xee, 0x30,
	0x85, 0x7e, 0x8f, 0xe7, 0xfd, 0x96, 0x2a, 0xe0, 0xdb, 0xec, 0xe1, 0xfa, 0x69, 0x26, 0x2c, 0x5f,
	0x01, 0xd5, 0x58, 0x86, 0x32, 0x9d, 0xca, 0xcc, 0xee, 0x85, 0x8a, 0xc1, 0xb2, 0xd2, 0x76, 0x33,
	0x87, 0xf9, 0xce, 0x6e, 0xa9, 0x02, 0x5f, 0xed, 0xae, 0x80, 0xfa, 0xd9, 0xee, 0x5a, 0xa6, 0x9b,
	0xa9, 0x65, 0xd3, 0x2d, 0x55, 0x60, 0xd1, 0x6e, 0xa9, 0xe9, 0x76, 0xd5, 0xf3, 0xb4, 0x8d, 0x5e,
	0xa6, 0x6d, 0xf4, 0x36, 0x6d, 0x23, 0xdc, 0x92, 0x3a, 0x21, 0xe3, 0x18, 0xc0, 0x90, 0x31, 0x88,
	0x98, 0xe4, 0x2b, 0x81, 0xf8, 0xd7, 0xdf, 0xad, 0x5f, 0x82, 0x88, 0xfd, 0xf3, 0x3d, 0xf4, 0xca,
	0x39, 0xba, 0x26, 0x09, 0xb7, 0xc3, 0x2c, 0x24, 0x91, 0x1c, 0x51, 0x97, 0x4a, 0x67, 0xa9, 0xd4,
	0x2d, 0x8a, 0x44, 0xab, 0x68, 0x61, 0xaf, 0x84, 0x55, 0xb7, 0x09, 0xf6, 0x3f, 0x02, 0x00, 0x00,
	0xff, 0xff, 0x1b, 0x2c, 0x52, 0xcf, 0x7b, 0x04, 0x00, 0x00,
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
	StreamSearchObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamSearchObjectClient, error)
	InsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamInsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamInsertObjectClient, error)
	MultiInsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Locations, error)
	UpdateObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamUpdateObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpdateObjectClient, error)
	MultiUpdateObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Locations, error)
	UpsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamUpsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpsertObjectClient, error)
	MultiUpsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type filterClient struct {
	cc *grpc.ClientConn
}

func NewFilterClient(cc *grpc.ClientConn) FilterClient {
	return &filterClient{cc}
}

func (c *filterClient) SearchObject(ctx context.Context, in *payload.Search_ObjectRequest, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/filter.Filter/SearchObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) StreamSearchObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamSearchObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Filter_serviceDesc.Streams[0], "/filter.Filter/StreamSearchObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &filterStreamSearchObjectClient{stream}
	return x, nil
}

type Filter_StreamSearchObjectClient interface {
	Send(*payload.Search_ObjectRequest) error
	Recv() (*payload.Search_Response, error)
	grpc.ClientStream
}

type filterStreamSearchObjectClient struct {
	grpc.ClientStream
}

func (x *filterStreamSearchObjectClient) Send(m *payload.Search_ObjectRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamSearchObjectClient) Recv() (*payload.Search_Response, error) {
	m := new(payload.Search_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) InsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/filter.Filter/InsertObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) StreamInsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamInsertObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Filter_serviceDesc.Streams[1], "/filter.Filter/StreamInsertObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &filterStreamInsertObjectClient{stream}
	return x, nil
}

type Filter_StreamInsertObjectClient interface {
	Send(*payload.Object_Blob) error
	Recv() (*payload.Object_Location, error)
	grpc.ClientStream
}

type filterStreamInsertObjectClient struct {
	grpc.ClientStream
}

func (x *filterStreamInsertObjectClient) Send(m *payload.Object_Blob) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamInsertObjectClient) Recv() (*payload.Object_Location, error) {
	m := new(payload.Object_Location)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) MultiInsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/filter.Filter/MultiInsertObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) UpdateObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/filter.Filter/UpdateObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) StreamUpdateObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpdateObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Filter_serviceDesc.Streams[2], "/filter.Filter/StreamUpdateObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &filterStreamUpdateObjectClient{stream}
	return x, nil
}

type Filter_StreamUpdateObjectClient interface {
	Send(*payload.Object_Blob) error
	Recv() (*payload.Object_Location, error)
	grpc.ClientStream
}

type filterStreamUpdateObjectClient struct {
	grpc.ClientStream
}

func (x *filterStreamUpdateObjectClient) Send(m *payload.Object_Blob) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamUpdateObjectClient) Recv() (*payload.Object_Location, error) {
	m := new(payload.Object_Location)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) MultiUpdateObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/filter.Filter/MultiUpdateObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) UpsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/filter.Filter/UpsertObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) StreamUpsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpsertObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Filter_serviceDesc.Streams[3], "/filter.Filter/StreamUpsertObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &filterStreamUpsertObjectClient{stream}
	return x, nil
}

type Filter_StreamUpsertObjectClient interface {
	Send(*payload.Object_Blob) error
	Recv() (*payload.Object_Location, error)
	grpc.ClientStream
}

type filterStreamUpsertObjectClient struct {
	grpc.ClientStream
}

func (x *filterStreamUpsertObjectClient) Send(m *payload.Object_Blob) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamUpsertObjectClient) Recv() (*payload.Object_Location, error) {
	m := new(payload.Object_Location)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) MultiUpsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/filter.Filter/MultiUpsertObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilterServer is the server API for Filter service.
type FilterServer interface {
	SearchObject(context.Context, *payload.Search_ObjectRequest) (*payload.Search_Response, error)
	StreamSearchObject(Filter_StreamSearchObjectServer) error
	InsertObject(context.Context, *payload.Object_Blob) (*payload.Object_Location, error)
	StreamInsertObject(Filter_StreamInsertObjectServer) error
	MultiInsertObject(context.Context, *payload.Object_Blob) (*payload.Object_Locations, error)
	UpdateObject(context.Context, *payload.Object_Blob) (*payload.Object_Location, error)
	StreamUpdateObject(Filter_StreamUpdateObjectServer) error
	MultiUpdateObject(context.Context, *payload.Object_Blob) (*payload.Object_Locations, error)
	UpsertObject(context.Context, *payload.Object_Blob) (*payload.Object_Location, error)
	StreamUpsertObject(Filter_StreamUpsertObjectServer) error
	MultiUpsertObject(context.Context, *payload.Object_Blob) (*payload.Object_Locations, error)
}

// UnimplementedFilterServer can be embedded to have forward compatible implementations.
type UnimplementedFilterServer struct {
}

func (*UnimplementedFilterServer) SearchObject(ctx context.Context, req *payload.Search_ObjectRequest) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchObject not implemented")
}
func (*UnimplementedFilterServer) StreamSearchObject(srv Filter_StreamSearchObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearchObject not implemented")
}
func (*UnimplementedFilterServer) InsertObject(ctx context.Context, req *payload.Object_Blob) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertObject not implemented")
}
func (*UnimplementedFilterServer) StreamInsertObject(srv Filter_StreamInsertObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamInsertObject not implemented")
}
func (*UnimplementedFilterServer) MultiInsertObject(ctx context.Context, req *payload.Object_Blob) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiInsertObject not implemented")
}
func (*UnimplementedFilterServer) UpdateObject(ctx context.Context, req *payload.Object_Blob) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateObject not implemented")
}
func (*UnimplementedFilterServer) StreamUpdateObject(srv Filter_StreamUpdateObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpdateObject not implemented")
}
func (*UnimplementedFilterServer) MultiUpdateObject(ctx context.Context, req *payload.Object_Blob) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpdateObject not implemented")
}
func (*UnimplementedFilterServer) UpsertObject(ctx context.Context, req *payload.Object_Blob) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertObject not implemented")
}
func (*UnimplementedFilterServer) StreamUpsertObject(srv Filter_StreamUpsertObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpsertObject not implemented")
}
func (*UnimplementedFilterServer) MultiUpsertObject(ctx context.Context, req *payload.Object_Blob) (*payload.Object_Locations, error) {
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
		FullMethod: "/filter.Filter/SearchObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).SearchObject(ctx, req.(*payload.Search_ObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_StreamSearchObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilterServer).StreamSearchObject(&filterStreamSearchObjectServer{stream})
}

type Filter_StreamSearchObjectServer interface {
	Send(*payload.Search_Response) error
	Recv() (*payload.Search_ObjectRequest, error)
	grpc.ServerStream
}

type filterStreamSearchObjectServer struct {
	grpc.ServerStream
}

func (x *filterStreamSearchObjectServer) Send(m *payload.Search_Response) error {
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
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).InsertObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.Filter/InsertObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).InsertObject(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_StreamInsertObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilterServer).StreamInsertObject(&filterStreamInsertObjectServer{stream})
}

type Filter_StreamInsertObjectServer interface {
	Send(*payload.Object_Location) error
	Recv() (*payload.Object_Blob, error)
	grpc.ServerStream
}

type filterStreamInsertObjectServer struct {
	grpc.ServerStream
}

func (x *filterStreamInsertObjectServer) Send(m *payload.Object_Location) error {
	return x.ServerStream.SendMsg(m)
}

func (x *filterStreamInsertObjectServer) Recv() (*payload.Object_Blob, error) {
	m := new(payload.Object_Blob)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Filter_MultiInsertObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).MultiInsertObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.Filter/MultiInsertObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).MultiInsertObject(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_UpdateObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).UpdateObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.Filter/UpdateObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).UpdateObject(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_StreamUpdateObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilterServer).StreamUpdateObject(&filterStreamUpdateObjectServer{stream})
}

type Filter_StreamUpdateObjectServer interface {
	Send(*payload.Object_Location) error
	Recv() (*payload.Object_Blob, error)
	grpc.ServerStream
}

type filterStreamUpdateObjectServer struct {
	grpc.ServerStream
}

func (x *filterStreamUpdateObjectServer) Send(m *payload.Object_Location) error {
	return x.ServerStream.SendMsg(m)
}

func (x *filterStreamUpdateObjectServer) Recv() (*payload.Object_Blob, error) {
	m := new(payload.Object_Blob)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Filter_MultiUpdateObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).MultiUpdateObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.Filter/MultiUpdateObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).MultiUpdateObject(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_UpsertObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).UpsertObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.Filter/UpsertObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).UpsertObject(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

func _Filter_StreamUpsertObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilterServer).StreamUpsertObject(&filterStreamUpsertObjectServer{stream})
}

type Filter_StreamUpsertObjectServer interface {
	Send(*payload.Object_Location) error
	Recv() (*payload.Object_Blob, error)
	grpc.ServerStream
}

type filterStreamUpsertObjectServer struct {
	grpc.ServerStream
}

func (x *filterStreamUpsertObjectServer) Send(m *payload.Object_Location) error {
	return x.ServerStream.SendMsg(m)
}

func (x *filterStreamUpsertObjectServer) Recv() (*payload.Object_Blob, error) {
	m := new(payload.Object_Blob)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Filter_MultiUpsertObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServer).MultiUpsertObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.Filter/MultiUpsertObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServer).MultiUpsertObject(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

var _Filter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "filter.Filter",
	HandlerType: (*FilterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchObject",
			Handler:    _Filter_SearchObject_Handler,
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
	Metadata: "apis/proto/gateway/filter/filter.proto",
}
