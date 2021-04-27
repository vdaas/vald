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

func init() { proto.RegisterFile("apis/proto/v1/vald/search.proto", fileDescriptor_f8168beed818734d) }
func init() {
	golang_proto.RegisterFile("apis/proto/v1/vald/search.proto", fileDescriptor_f8168beed818734d)
}

var fileDescriptor_f8168beed818734d = []byte{
	// 397 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xcf, 0x4a, 0xf3, 0x40,
	0x10, 0xff, 0xf6, 0xe3, 0xa3, 0x85, 0xfd, 0x0a, 0x2d, 0x91, 0x82, 0xa4, 0x25, 0xa5, 0xf1, 0xa0,
	0x08, 0xee, 0x5a, 0x3d, 0x08, 0x1e, 0x4b, 0x2f, 0x3d, 0x08, 0x62, 0x51, 0xd0, 0x83, 0xb0, 0x6d,
	0x62, 0xba, 0x90, 0x66, 0xd7, 0x6c, 0x12, 0xe8, 0xd5, 0x57, 0xf0, 0x85, 0x3c, 0xf6, 0x24, 0x82,
	0x2f, 0x20, 0xad, 0x0f, 0x22, 0xbb, 0x9b, 0xb4, 0x51, 0x62, 0xc9, 0xa9, 0xd3, 0x9d, 0xdf, 0x9f,
	0x99, 0xc9, 0x0c, 0xec, 0x10, 0x4e, 0x05, 0xe6, 0x21, 0x8b, 0x18, 0x4e, 0x7a, 0x38, 0x21, 0xbe,
	0x83, 0x85, 0x4b, 0xc2, 0xc9, 0x14, 0xa9, 0x47, 0xa3, 0x2a, 0x9f, 0x50, 0xd2, 0x33, 0xf7, 0xbe,
	0x23, 0x39, 0x99, 0xfb, 0x8c, 0x38, 0xd9, 0xaf, 0x46, 0x9b, 0x47, 0x1e, 0x8d, 0xa6, 0xf1, 0x18,
	0x4d, 0xd8, 0x0c, 0x7b, 0xcc, 0x63, 0x1a, 0x3f, 0x8e, 0x1f, 0xd4, 0x3f, 0x4d, 0x96, 0x51, 0x0a,
	0x3f, 0xfb, 0x09, 0xf7, 0x18, 0xf3, 0x7c, 0x57, 0x39, 0xe9, 0x10, 0x13, 0x4e, 0x31, 0x09, 0x02,
	0x16, 0x91, 0x88, 0xb2, 0x40, 0x68, 0xe2, 0xc9, 0xeb, 0x3f, 0x58, 0x19, 0xa9, 0x32, 0x8d, 0xeb,
	0x75, 0x64, 0xa2, 0xac, 0x98, 0xa4, 0x87, 0xf4, 0x1b, 0xba, 0x72, 0x1f, 0x63, 0x57, 0x44, 0x66,
	0xab, 0x30, 0x27, 0x38, 0x0b, 0x84, 0x6b, 0x1b, 0x4f, 0xef, 0x9f, 0xcf, 0x7f, 0x6b, 0x76, 0x35,
	0x6d, 0xfd, 0x1c, 0x1c, 0x1a, 0xf7, 0x10, 0x6a, 0x58, 0x7f, 0x3e, 0x1c, 0x18, 0xed, 0x02, 0xfa,
	0x70, 0x50, 0x4a, 0xbc, 0xa9, 0xc4, 0xeb, 0x36, 0x4c, 0xc5, 0x31, 0x75, 0xa4, 0xfe, 0x08, 0xd6,
	0x46, 0x51, 0xe8, 0x92, 0x59, 0x89, 0xe2, 0xbb, 0x05, 0x39, 0x4d, 0x5e, 0xbb, 0xfc, 0x39, 0x00,
	0xc7, 0xc0, 0xb8, 0x85, 0x8d, 0xbc, 0x68, 0x89, 0xd2, 0x4b, 0x4b, 0x53, 0xf8, 0xff, 0x22, 0xf6,
	0x23, 0x9a, 0x96, 0xdb, 0x29, 0xe0, 0xa9, 0x7c, 0x26, 0xdc, 0xde, 0x32, 0x13, 0x61, 0xb7, 0xd4,
	0x50, 0x9a, 0x76, 0x23, 0x1b, 0xca, 0x4c, 0x72, 0xb9, 0xef, 0xca, 0xd1, 0x84, 0xb0, 0x9e, 0xb3,
	0x52, 0x4d, 0x74, 0x7f, 0xb3, 0xdb, 0x74, 0xb2, 0xdd, 0xd0, 0x52, 0x86, 0xbb, 0xf6, 0xce, 0xe6,
	0x2b, 0xe4, 0x3d, 0xfb, 0xde, 0x62, 0x69, 0x81, 0xb7, 0xa5, 0x05, 0x3e, 0x96, 0x16, 0x78, 0x59,
	0x59, 0x60, 0xb1, 0xb2, 0x00, 0x34, 0x59, 0xe8, 0xa1, 0xc4, 0x21, 0x44, 0x20, 0x75, 0x06, 0x84,
	0x53, 0xa9, 0x2e, 0xe3, 0x3e, 0xbc, 0x21, 0xbe, 0xa3, 0x7d, 0x2e, 0xc1, 0xdd, 0x7e, 0x6e, 0x93,
	0x15, 0x41, 0x9f, 0x92, 0xde, 0xe4, 0x90, 0x4f, 0xb2, 0xe3, 0x1a, 0x57, 0xd4, 0x02, 0x9f, 0x7e,
	0x05, 0x00, 0x00, 0xff, 0xff, 0x27, 0x0b, 0xb8, 0xfb, 0x79, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SearchClient is the client API for Search service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SearchClient interface {
	Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error)
	SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)
	StreamSearch(ctx context.Context, opts ...grpc.CallOption) (Search_StreamSearchClient, error)
	StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (Search_StreamSearchByIDClient, error)
	MultiSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	MultiSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
}

type searchClient struct {
	cc *grpc.ClientConn
}

func NewSearchClient(cc *grpc.ClientConn) SearchClient {
	return &searchClient{cc}
}

func (c *searchClient) Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/SearchByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (Search_StreamSearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Search_serviceDesc.Streams[0], "/vald.v1.Search/StreamSearch", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchStreamSearchClient{stream}
	return x, nil
}

type Search_StreamSearchClient interface {
	Send(*payload.Search_Request) error
	Recv() (*payload.Search_StreamResponse, error)
	grpc.ClientStream
}

type searchStreamSearchClient struct {
	grpc.ClientStream
}

func (x *searchStreamSearchClient) Send(m *payload.Search_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchStreamSearchClient) Recv() (*payload.Search_StreamResponse, error) {
	m := new(payload.Search_StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (Search_StreamSearchByIDClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Search_serviceDesc.Streams[1], "/vald.v1.Search/StreamSearchByID", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchStreamSearchByIDClient{stream}
	return x, nil
}

type Search_StreamSearchByIDClient interface {
	Send(*payload.Search_IDRequest) error
	Recv() (*payload.Search_StreamResponse, error)
	grpc.ClientStream
}

type searchStreamSearchByIDClient struct {
	grpc.ClientStream
}

func (x *searchStreamSearchByIDClient) Send(m *payload.Search_IDRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchStreamSearchByIDClient) Recv() (*payload.Search_StreamResponse, error) {
	m := new(payload.Search_StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) MultiSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error) {
	out := new(payload.Search_Responses)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/MultiSearch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) MultiSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error) {
	out := new(payload.Search_Responses)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/MultiSearchByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchServer is the server API for Search service.
type SearchServer interface {
	Search(context.Context, *payload.Search_Request) (*payload.Search_Response, error)
	SearchByID(context.Context, *payload.Search_IDRequest) (*payload.Search_Response, error)
	StreamSearch(Search_StreamSearchServer) error
	StreamSearchByID(Search_StreamSearchByIDServer) error
	MultiSearch(context.Context, *payload.Search_MultiRequest) (*payload.Search_Responses, error)
	MultiSearchByID(context.Context, *payload.Search_MultiIDRequest) (*payload.Search_Responses, error)
}

// UnimplementedSearchServer can be embedded to have forward compatible implementations.
type UnimplementedSearchServer struct {
}

func (*UnimplementedSearchServer) Search(ctx context.Context, req *payload.Search_Request) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (*UnimplementedSearchServer) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchByID not implemented")
}
func (*UnimplementedSearchServer) StreamSearch(srv Search_StreamSearchServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearch not implemented")
}
func (*UnimplementedSearchServer) StreamSearchByID(srv Search_StreamSearchByIDServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearchByID not implemented")
}
func (*UnimplementedSearchServer) MultiSearch(ctx context.Context, req *payload.Search_MultiRequest) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiSearch not implemented")
}
func (*UnimplementedSearchServer) MultiSearchByID(ctx context.Context, req *payload.Search_MultiIDRequest) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiSearchByID not implemented")
}

func RegisterSearchServer(s *grpc.Server, srv SearchServer) {
	s.RegisterService(&_Search_serviceDesc, srv)
}

func _Search_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Search_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).Search(ctx, req.(*payload.Search_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_SearchByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Search_IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).SearchByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/SearchByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).SearchByID(ctx, req.(*payload.Search_IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_StreamSearch_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SearchServer).StreamSearch(&searchStreamSearchServer{stream})
}

type Search_StreamSearchServer interface {
	Send(*payload.Search_StreamResponse) error
	Recv() (*payload.Search_Request, error)
	grpc.ServerStream
}

type searchStreamSearchServer struct {
	grpc.ServerStream
}

func (x *searchStreamSearchServer) Send(m *payload.Search_StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchStreamSearchServer) Recv() (*payload.Search_Request, error) {
	m := new(payload.Search_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Search_StreamSearchByID_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SearchServer).StreamSearchByID(&searchStreamSearchByIDServer{stream})
}

type Search_StreamSearchByIDServer interface {
	Send(*payload.Search_StreamResponse) error
	Recv() (*payload.Search_IDRequest, error)
	grpc.ServerStream
}

type searchStreamSearchByIDServer struct {
	grpc.ServerStream
}

func (x *searchStreamSearchByIDServer) Send(m *payload.Search_StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchStreamSearchByIDServer) Recv() (*payload.Search_IDRequest, error) {
	m := new(payload.Search_IDRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Search_MultiSearch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Search_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).MultiSearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/MultiSearch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).MultiSearch(ctx, req.(*payload.Search_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_MultiSearchByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Search_MultiIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).MultiSearchByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/MultiSearchByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).MultiSearchByID(ctx, req.(*payload.Search_MultiIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Search_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vald.v1.Search",
	HandlerType: (*SearchServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _Search_Search_Handler,
		},
		{
			MethodName: "SearchByID",
			Handler:    _Search_SearchByID_Handler,
		},
		{
			MethodName: "MultiSearch",
			Handler:    _Search_MultiSearch_Handler,
		},
		{
			MethodName: "MultiSearchByID",
			Handler:    _Search_MultiSearchByID_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamSearch",
			Handler:       _Search_StreamSearch_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamSearchByID",
			Handler:       _Search_StreamSearchByID_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "apis/proto/v1/vald/search.proto",
}
