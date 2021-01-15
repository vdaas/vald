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

func init() { proto.RegisterFile("apis/proto/v1/vald/search.proto", fileDescriptor_f8168beed818734d) }

var fileDescriptor_f8168beed818734d = []byte{
	// 356 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4f, 0x4a, 0xc3, 0x40,
	0x14, 0xc6, 0x3b, 0x22, 0x2d, 0x8c, 0x85, 0x96, 0x40, 0x37, 0x69, 0x69, 0x69, 0x5c, 0x28, 0x2e,
	0x66, 0xac, 0xee, 0x5c, 0x96, 0x6e, 0xba, 0x10, 0xc4, 0xaa, 0xa0, 0x0b, 0xe5, 0xb5, 0x19, 0xd2,
	0x81, 0x34, 0x33, 0x66, 0xa6, 0x81, 0x6e, 0xbd, 0x82, 0x17, 0xf1, 0x18, 0x2e, 0x05, 0x2f, 0x20,
	0xc5, 0x83, 0x48, 0x66, 0x92, 0x5a, 0x21, 0x96, 0xac, 0x66, 0x78, 0xef, 0x7d, 0xbf, 0xf7, 0x87,
	0x0f, 0xf7, 0x40, 0x72, 0x45, 0x65, 0x2c, 0xb4, 0xa0, 0xc9, 0x80, 0x26, 0x10, 0xfa, 0x54, 0x31,
	0x88, 0x67, 0x73, 0x62, 0x82, 0x4e, 0x2d, 0x0d, 0x91, 0x64, 0xe0, 0x1e, 0xfe, 0xad, 0x94, 0xb0,
	0x0a, 0x05, 0xf8, 0xf9, 0x6b, 0xab, 0xdd, 0x4e, 0x20, 0x44, 0x10, 0x32, 0x0a, 0x92, 0x53, 0x88,
	0x22, 0xa1, 0x41, 0x73, 0x11, 0x29, 0x9b, 0x3d, 0x7b, 0xdb, 0xc7, 0xd5, 0x89, 0x81, 0x3b, 0xb7,
	0x9b, 0x9f, 0x4b, 0x72, 0x44, 0x32, 0x20, 0x36, 0x46, 0xae, 0xd9, 0xf3, 0x92, 0x29, 0xed, 0xb6,
	0x0b, 0x73, 0x4a, 0x8a, 0x48, 0x31, 0xcf, 0x79, 0xf9, 0xfc, 0x7e, 0xdd, 0xab, 0x7b, 0xb5, 0x6c,
	0xe0, 0x0b, 0x74, 0xe2, 0x3c, 0x62, 0x6c, 0xcb, 0x86, 0xab, 0xf1, 0xc8, 0xe9, 0x14, 0xc8, 0xc7,
	0xa3, 0x52, 0xf0, 0x96, 0x81, 0x37, 0x3c, 0x9c, 0xc1, 0x29, 0xf7, 0x53, 0xfe, 0x04, 0xd7, 0x27,
	0x3a, 0x66, 0xb0, 0x28, 0x31, 0x7c, 0xbf, 0x20, 0x67, 0xc5, 0x9b, 0x2e, 0x95, 0x63, 0x74, 0x8a,
	0x9c, 0x7b, 0xdc, 0xdc, 0x86, 0x96, 0x18, 0xbd, 0x34, 0x9a, 0xe3, 0x83, 0xcb, 0x65, 0xa8, 0x79,
	0x36, 0x6e, 0xaf, 0x40, 0x67, 0xf2, 0x39, 0xb8, 0xb3, 0xe3, 0x26, 0xca, 0x6b, 0x9b, 0xa3, 0xb4,
	0xbc, 0x66, 0x7e, 0x94, 0x45, 0xaa, 0x95, 0x21, 0x4b, 0x4f, 0x73, 0x83, 0x1b, 0x5b, 0xad, 0xcc,
	0x12, 0xfd, 0xff, 0xda, 0xfd, 0x6e, 0xb2, 0xbb, 0x61, 0x65, 0xf8, 0xf4, 0xbe, 0xee, 0xa2, 0x8f,
	0x75, 0x17, 0x7d, 0xad, 0xbb, 0x08, 0xbb, 0x22, 0x0e, 0x48, 0xe2, 0x03, 0x28, 0x62, 0x6c, 0x09,
	0x92, 0xa7, 0xba, 0xf4, 0x3f, 0xc4, 0x77, 0x10, 0xfa, 0x96, 0x70, 0x85, 0x1e, 0x8e, 0x02, 0xae,
	0xe7, 0xcb, 0x29, 0x99, 0x89, 0x05, 0x35, 0x02, 0x6b, 0x6d, 0xe3, 0xe1, 0x20, 0x96, 0xb3, 0xdc,
	0xec, 0xd3, 0xaa, 0xb1, 0xe6, 0xf9, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x70, 0x1f, 0x76,
	0x09, 0x03, 0x00, 0x00,
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
