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

func init() { proto.RegisterFile("search.proto", fileDescriptor_453745cff914010e) }

var fileDescriptor_453745cff914010e = []byte{
	// 322 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x4e, 0x4d, 0x2c,
	0x4a, 0xce, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29, 0x4b, 0xcc, 0x49, 0x91, 0xe2,
	0x2d, 0x48, 0xac, 0xcc, 0xc9, 0x4f, 0x4c, 0x81, 0x08, 0x4a, 0xc9, 0xa4, 0xe7, 0xe7, 0xa7, 0xe7,
	0xa4, 0xea, 0x27, 0x16, 0x64, 0xea, 0x27, 0xe6, 0xe5, 0xe5, 0x97, 0x24, 0x96, 0x64, 0xe6, 0xe7,
	0x15, 0x43, 0x65, 0x79, 0x0a, 0x92, 0xf4, 0xd3, 0x0b, 0x73, 0x20, 0x3c, 0xa3, 0xf7, 0xcc, 0x5c,
	0x6c, 0xc1, 0x60, 0x13, 0x85, 0xfc, 0xe1, 0x2c, 0x71, 0x3d, 0x98, 0x81, 0x10, 0x01, 0xbd, 0xa0,
	0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x29, 0x09, 0x4c, 0x89, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54,
	0x25, 0xa1, 0xa6, 0xcb, 0x4f, 0x26, 0x33, 0xf1, 0x28, 0xb1, 0xeb, 0x43, 0xdc, 0x67, 0xc5, 0xa8,
	0x25, 0x14, 0xc1, 0xc5, 0x05, 0x51, 0xe6, 0x54, 0xe9, 0xe9, 0x22, 0x24, 0x89, 0xae, 0xd7, 0xd3,
	0x85, 0xb0, 0xb1, 0xa2, 0x60, 0x63, 0xf9, 0x95, 0xb8, 0xa0, 0xc6, 0xea, 0x67, 0xa6, 0x80, 0x4c,
	0x76, 0xe7, 0xe2, 0x09, 0x2e, 0x29, 0x4a, 0x4d, 0xcc, 0x25, 0xdf, 0xc1, 0x0c, 0x1a, 0x8c, 0x06,
	0x8c, 0x42, 0xbe, 0x5c, 0x02, 0xc8, 0x06, 0x91, 0xef, 0x50, 0x88, 0x71, 0x5e, 0x5c, 0xdc, 0xbe,
	0xa5, 0x39, 0x25, 0x99, 0x50, 0x67, 0xc9, 0xa0, 0x2b, 0x07, 0x4b, 0xc2, 0x0c, 0x93, 0xc4, 0x65,
	0x58, 0x31, 0xc8, 0x34, 0xa1, 0x00, 0x2e, 0x7e, 0x24, 0xb3, 0xc0, 0x2e, 0x93, 0xc3, 0x6a, 0x1e,
	0xc2, 0x79, 0xf8, 0x4d, 0x94, 0x62, 0xd9, 0xf0, 0x40, 0x9e, 0xc9, 0xc9, 0xf7, 0xc4, 0x23, 0x39,
	0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0xe4, 0xe2, 0xcb, 0x2f, 0x4a, 0xd7, 0x2b,
	0x4b, 0x49, 0x4c, 0x2c, 0xd6, 0x03, 0x25, 0x25, 0x27, 0xf6, 0xb0, 0xc4, 0x9c, 0x14, 0xc7, 0x82,
	0xcc, 0x00, 0xc6, 0x28, 0x95, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d,
	0xb0, 0x0a, 0x7d, 0x90, 0x0a, 0x50, 0xaa, 0x2a, 0xd6, 0x4f, 0x2f, 0x2a, 0x48, 0x06, 0x73, 0x93,
	0xd8, 0xc0, 0xe9, 0xc8, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xb8, 0xf1, 0xd7, 0x57, 0x98, 0x02,
	0x00, 0x00,
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
	MultiSearch(ctx context.Context, opts ...grpc.CallOption) (Search_MultiSearchClient, error)
	MultiSearchByID(ctx context.Context, opts ...grpc.CallOption) (Search_MultiSearchByIDClient, error)
}

type searchClient struct {
	cc *grpc.ClientConn
}

func NewSearchClient(cc *grpc.ClientConn) SearchClient {
	return &searchClient{cc}
}

func (c *searchClient) Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/vald.Search/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/vald.Search/SearchByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (Search_StreamSearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Search_serviceDesc.Streams[0], "/vald.Search/StreamSearch", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchStreamSearchClient{stream}
	return x, nil
}

type Search_StreamSearchClient interface {
	Send(*payload.Search_Request) error
	Recv() (*payload.Search_Response, error)
	grpc.ClientStream
}

type searchStreamSearchClient struct {
	grpc.ClientStream
}

func (x *searchStreamSearchClient) Send(m *payload.Search_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchStreamSearchClient) Recv() (*payload.Search_Response, error) {
	m := new(payload.Search_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (Search_StreamSearchByIDClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Search_serviceDesc.Streams[1], "/vald.Search/StreamSearchByID", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchStreamSearchByIDClient{stream}
	return x, nil
}

type Search_StreamSearchByIDClient interface {
	Send(*payload.Search_IDRequest) error
	Recv() (*payload.Search_Response, error)
	grpc.ClientStream
}

type searchStreamSearchByIDClient struct {
	grpc.ClientStream
}

func (x *searchStreamSearchByIDClient) Send(m *payload.Search_IDRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchStreamSearchByIDClient) Recv() (*payload.Search_Response, error) {
	m := new(payload.Search_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) MultiSearch(ctx context.Context, opts ...grpc.CallOption) (Search_MultiSearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Search_serviceDesc.Streams[2], "/vald.Search/MultiSearch", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchMultiSearchClient{stream}
	return x, nil
}

type Search_MultiSearchClient interface {
	Send(*payload.Search_MultiRequest) error
	CloseAndRecv() (*payload.Search_Responses, error)
	grpc.ClientStream
}

type searchMultiSearchClient struct {
	grpc.ClientStream
}

func (x *searchMultiSearchClient) Send(m *payload.Search_MultiRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchMultiSearchClient) CloseAndRecv() (*payload.Search_Responses, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(payload.Search_Responses)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) MultiSearchByID(ctx context.Context, opts ...grpc.CallOption) (Search_MultiSearchByIDClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Search_serviceDesc.Streams[3], "/vald.Search/MultiSearchByID", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchMultiSearchByIDClient{stream}
	return x, nil
}

type Search_MultiSearchByIDClient interface {
	Send(*payload.Search_MultiIDRequest) error
	CloseAndRecv() (*payload.Search_Responses, error)
	grpc.ClientStream
}

type searchMultiSearchByIDClient struct {
	grpc.ClientStream
}

func (x *searchMultiSearchByIDClient) Send(m *payload.Search_MultiIDRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchMultiSearchByIDClient) CloseAndRecv() (*payload.Search_Responses, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(payload.Search_Responses)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SearchServer is the server API for Search service.
type SearchServer interface {
	Search(context.Context, *payload.Search_Request) (*payload.Search_Response, error)
	SearchByID(context.Context, *payload.Search_IDRequest) (*payload.Search_Response, error)
	StreamSearch(Search_StreamSearchServer) error
	StreamSearchByID(Search_StreamSearchByIDServer) error
	MultiSearch(Search_MultiSearchServer) error
	MultiSearchByID(Search_MultiSearchByIDServer) error
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
func (*UnimplementedSearchServer) MultiSearch(srv Search_MultiSearchServer) error {
	return status.Errorf(codes.Unimplemented, "method MultiSearch not implemented")
}
func (*UnimplementedSearchServer) MultiSearchByID(srv Search_MultiSearchByIDServer) error {
	return status.Errorf(codes.Unimplemented, "method MultiSearchByID not implemented")
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
		FullMethod: "/vald.Search/Search",
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
		FullMethod: "/vald.Search/SearchByID",
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
	Send(*payload.Search_Response) error
	Recv() (*payload.Search_Request, error)
	grpc.ServerStream
}

type searchStreamSearchServer struct {
	grpc.ServerStream
}

func (x *searchStreamSearchServer) Send(m *payload.Search_Response) error {
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
	Send(*payload.Search_Response) error
	Recv() (*payload.Search_IDRequest, error)
	grpc.ServerStream
}

type searchStreamSearchByIDServer struct {
	grpc.ServerStream
}

func (x *searchStreamSearchByIDServer) Send(m *payload.Search_Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchStreamSearchByIDServer) Recv() (*payload.Search_IDRequest, error) {
	m := new(payload.Search_IDRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Search_MultiSearch_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SearchServer).MultiSearch(&searchMultiSearchServer{stream})
}

type Search_MultiSearchServer interface {
	SendAndClose(*payload.Search_Responses) error
	Recv() (*payload.Search_MultiRequest, error)
	grpc.ServerStream
}

type searchMultiSearchServer struct {
	grpc.ServerStream
}

func (x *searchMultiSearchServer) SendAndClose(m *payload.Search_Responses) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchMultiSearchServer) Recv() (*payload.Search_MultiRequest, error) {
	m := new(payload.Search_MultiRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Search_MultiSearchByID_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SearchServer).MultiSearchByID(&searchMultiSearchByIDServer{stream})
}

type Search_MultiSearchByIDServer interface {
	SendAndClose(*payload.Search_Responses) error
	Recv() (*payload.Search_MultiIDRequest, error)
	grpc.ServerStream
}

type searchMultiSearchByIDServer struct {
	grpc.ServerStream
}

func (x *searchMultiSearchByIDServer) SendAndClose(m *payload.Search_Responses) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchMultiSearchByIDServer) Recv() (*payload.Search_MultiIDRequest, error) {
	m := new(payload.Search_MultiIDRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Search_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vald.Search",
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
		{
			StreamName:    "MultiSearch",
			Handler:       _Search_MultiSearch_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "MultiSearchByID",
			Handler:       _Search_MultiSearchByID_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "search.proto",
}
