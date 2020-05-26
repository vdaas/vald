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

func init() { proto.RegisterFile("filter/filter.proto", fileDescriptor_f45b4a1fb0c46f6e) }

var fileDescriptor_f45b4a1fb0c46f6e = []byte{
	// 374 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xb1, 0x4e, 0xf3, 0x30,
	0x14, 0x85, 0x7f, 0xff, 0x43, 0x06, 0x2b, 0x54, 0x6a, 0xca, 0x00, 0x11, 0x2d, 0x52, 0x27, 0xc4,
	0x60, 0x23, 0xd8, 0x60, 0x41, 0x95, 0x00, 0x31, 0x54, 0xa0, 0x56, 0x30, 0xb0, 0xdd, 0x24, 0x26,
	0x0d, 0x72, 0x63, 0x37, 0x76, 0x8a, 0xba, 0xf2, 0x0a, 0xbc, 0x08, 0x1b, 0xaf, 0xc0, 0x88, 0xc4,
	0x0b, 0x54, 0x15, 0x0f, 0x82, 0x9a, 0x1b, 0x2a, 0xa8, 0x60, 0x20, 0x9d, 0x22, 0x9f, 0x9b, 0xfb,
	0xf9, 0x1c, 0x59, 0x87, 0x36, 0x6e, 0x13, 0x69, 0x45, 0xc6, 0xf1, 0xc3, 0x74, 0xa6, 0xac, 0xf2,
	0x1c, 0x3c, 0xf9, 0x6b, 0x1a, 0x26, 0x52, 0x41, 0x84, 0xb2, 0xbf, 0x15, 0x2b, 0x15, 0x4b, 0xc1,
	0x41, 0x27, 0x1c, 0xd2, 0x54, 0x59, 0xb0, 0x89, 0x4a, 0x4d, 0x39, 0x75, 0x75, 0xc0, 0xe3, 0x91,
	0xc4, 0xd3, 0xfe, 0xb3, 0x43, 0x9d, 0xd3, 0x82, 0xe2, 0x05, 0xd4, 0xed, 0x0b, 0xc8, 0xc2, 0xc1,
	0x45, 0x70, 0x27, 0x42, 0xeb, 0x35, 0xd9, 0x27, 0x16, 0x65, 0x86, 0x7a, 0x4f, 0x8c, 0x72, 0x61,
	0xac, 0xbf, 0xb1, 0x3c, 0xee, 0x09, 0xa3, 0x55, 0x6a, 0x44, 0x7b, 0xf3, 0xe1, 0xed, 0xfd, 0xf1,
	0x7f, 0xa3, 0x5d, 0xe3, 0xa6, 0x98, 0x70, 0x55, 0x2c, 0x1e, 0x92, 0x5d, 0xaf, 0x4b, 0xdd, 0xbe,
	0xcd, 0x04, 0x0c, 0x71, 0xa7, 0xfa, 0x1d, 0xff, 0x76, 0xc8, 0x1e, 0xf1, 0xfa, 0xd4, 0x3d, 0x4f,
	0x8d, 0xc8, 0x6c, 0x69, 0x79, 0x7d, 0xf1, 0x3f, 0x0a, 0xac, 0x23, 0x55, 0xe0, 0xd7, 0x16, 0xea,
	0xc9, 0x50, 0xdb, 0x49, 0xbb, 0xf9, 0x34, 0xdd, 0x26, 0x0b, 0x8f, 0x49, 0x01, 0xf8, 0xe2, 0xf1,
	0x98, 0x7a, 0xe8, 0xb1, 0x02, 0x1a, 0x6d, 0x1d, 0xd1, 0x7a, 0x37, 0x97, 0x36, 0xa9, 0x02, 0x98,
	0x67, 0xba, 0xd2, 0x11, 0x58, 0xb1, 0x42, 0xa6, 0xbc, 0x00, 0xfc, 0x94, 0xa9, 0x02, 0xfa, 0x7b,
	0xa6, 0x2a, 0x00, 0xcc, 0xb4, 0xe2, 0x3b, 0xe5, 0xfa, 0xb7, 0x77, 0xaa, 0x80, 0x5e, 0xce, 0xf4,
	0x77, 0x40, 0x47, 0xbf, 0xcc, 0x5a, 0xe4, 0x75, 0xd6, 0x22, 0xd3, 0x59, 0x8b, 0xd0, 0xa6, 0xca,
	0x62, 0x36, 0x8e, 0x00, 0x0c, 0x1b, 0x83, 0x8c, 0x58, 0x0c, 0x56, 0xdc, 0xc3, 0x84, 0x61, 0x43,
	0x3b, 0xf5, 0x6b, 0x90, 0x11, 0xf6, 0xec, 0x0c, 0x27, 0x97, 0xe4, 0x86, 0xc5, 0x89, 0x1d, 0xe4,
	0x01, 0x0b, 0xd5, 0x90, 0x17, 0xab, 0x7c, 0xbe, 0x3a, 0xaf, 0xad, 0xe1, 0x71, 0xa6, 0x43, 0x5e,
	0x42, 0xca, 0xd2, 0x07, 0x4e, 0x51, 0xd9, 0x83, 0x8f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2c, 0xf0,
	0xcb, 0xd3, 0x0c, 0x04, 0x00, 0x00,
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
	StreamSearch(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamSearchClient, error)
	InsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamInsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamInsertObjectClient, error)
	MultiInsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error)
	UpdateObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamUpdateObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpdateObjectClient, error)
	MultiUpdateObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error)
	UpsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamUpsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpsertObjectClient, error)
	MultiUpsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error)
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

func (c *filterClient) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamSearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Filter_serviceDesc.Streams[0], "/filter.Filter/StreamSearch", opts...)
	if err != nil {
		return nil, err
	}
	x := &filterStreamSearchClient{stream}
	return x, nil
}

type Filter_StreamSearchClient interface {
	Send(*payload.Search_ObjectRequest) error
	Recv() (*payload.Search_Response, error)
	grpc.ClientStream
}

type filterStreamSearchClient struct {
	grpc.ClientStream
}

func (x *filterStreamSearchClient) Send(m *payload.Search_ObjectRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamSearchClient) Recv() (*payload.Search_Response, error) {
	m := new(payload.Search_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) InsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
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
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type filterStreamInsertObjectClient struct {
	grpc.ClientStream
}

func (x *filterStreamInsertObjectClient) Send(m *payload.Object_Blob) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamInsertObjectClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) MultiInsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/filter.Filter/MultiInsertObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) UpdateObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
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
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type filterStreamUpdateObjectClient struct {
	grpc.ClientStream
}

func (x *filterStreamUpdateObjectClient) Send(m *payload.Object_Blob) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamUpdateObjectClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) MultiUpdateObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/filter.Filter/MultiUpdateObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filterClient) UpsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
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
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type filterStreamUpsertObjectClient struct {
	grpc.ClientStream
}

func (x *filterStreamUpsertObjectClient) Send(m *payload.Object_Blob) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filterStreamUpsertObjectClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filterClient) MultiUpsertObject(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/filter.Filter/MultiUpsertObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilterServer is the server API for Filter service.
type FilterServer interface {
	SearchObject(context.Context, *payload.Search_ObjectRequest) (*payload.Search_Response, error)
	StreamSearch(Filter_StreamSearchServer) error
	InsertObject(context.Context, *payload.Object_Blob) (*payload.Empty, error)
	StreamInsertObject(Filter_StreamInsertObjectServer) error
	MultiInsertObject(context.Context, *payload.Object_Blob) (*payload.Empty, error)
	UpdateObject(context.Context, *payload.Object_Blob) (*payload.Empty, error)
	StreamUpdateObject(Filter_StreamUpdateObjectServer) error
	MultiUpdateObject(context.Context, *payload.Object_Blob) (*payload.Empty, error)
	UpsertObject(context.Context, *payload.Object_Blob) (*payload.Empty, error)
	StreamUpsertObject(Filter_StreamUpsertObjectServer) error
	MultiUpsertObject(context.Context, *payload.Object_Blob) (*payload.Empty, error)
}

// UnimplementedFilterServer can be embedded to have forward compatible implementations.
type UnimplementedFilterServer struct {
}

func (*UnimplementedFilterServer) SearchObject(ctx context.Context, req *payload.Search_ObjectRequest) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchObject not implemented")
}
func (*UnimplementedFilterServer) StreamSearch(srv Filter_StreamSearchServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearch not implemented")
}
func (*UnimplementedFilterServer) InsertObject(ctx context.Context, req *payload.Object_Blob) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertObject not implemented")
}
func (*UnimplementedFilterServer) StreamInsertObject(srv Filter_StreamInsertObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamInsertObject not implemented")
}
func (*UnimplementedFilterServer) MultiInsertObject(ctx context.Context, req *payload.Object_Blob) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiInsertObject not implemented")
}
func (*UnimplementedFilterServer) UpdateObject(ctx context.Context, req *payload.Object_Blob) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateObject not implemented")
}
func (*UnimplementedFilterServer) StreamUpdateObject(srv Filter_StreamUpdateObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpdateObject not implemented")
}
func (*UnimplementedFilterServer) MultiUpdateObject(ctx context.Context, req *payload.Object_Blob) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpdateObject not implemented")
}
func (*UnimplementedFilterServer) UpsertObject(ctx context.Context, req *payload.Object_Blob) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertObject not implemented")
}
func (*UnimplementedFilterServer) StreamUpsertObject(srv Filter_StreamUpsertObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpsertObject not implemented")
}
func (*UnimplementedFilterServer) MultiUpsertObject(ctx context.Context, req *payload.Object_Blob) (*payload.Empty, error) {
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

func _Filter_StreamSearch_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilterServer).StreamSearch(&filterStreamSearchServer{stream})
}

type Filter_StreamSearchServer interface {
	Send(*payload.Search_Response) error
	Recv() (*payload.Search_ObjectRequest, error)
	grpc.ServerStream
}

type filterStreamSearchServer struct {
	grpc.ServerStream
}

func (x *filterStreamSearchServer) Send(m *payload.Search_Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *filterStreamSearchServer) Recv() (*payload.Search_ObjectRequest, error) {
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
	Send(*payload.Empty) error
	Recv() (*payload.Object_Blob, error)
	grpc.ServerStream
}

type filterStreamInsertObjectServer struct {
	grpc.ServerStream
}

func (x *filterStreamInsertObjectServer) Send(m *payload.Empty) error {
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
	Send(*payload.Empty) error
	Recv() (*payload.Object_Blob, error)
	grpc.ServerStream
}

type filterStreamUpdateObjectServer struct {
	grpc.ServerStream
}

func (x *filterStreamUpdateObjectServer) Send(m *payload.Empty) error {
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
	Send(*payload.Empty) error
	Recv() (*payload.Object_Blob, error)
	grpc.ServerStream
}

type filterStreamUpsertObjectServer struct {
	grpc.ServerStream
}

func (x *filterStreamUpsertObjectServer) Send(m *payload.Empty) error {
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
			StreamName:    "StreamSearch",
			Handler:       _Filter_StreamSearch_Handler,
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
	Metadata: "filter/filter.proto",
}
