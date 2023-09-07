//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	codes "github.com/vdaas/vald/internal/net/grpc/codes"
	status "github.com/vdaas/vald/internal/net/grpc/status"
	grpc "google.golang.org/grpc"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FilterClient is the client API for Filter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FilterClient interface {
	// A method to search object.
	SearchObject(ctx context.Context, in *payload.Search_ObjectRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)
	// A method to search multiple objects.
	MultiSearchObject(ctx context.Context, in *payload.Search_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	// A method to search object by bidirectional streaming.
	StreamSearchObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamSearchObjectClient, error)
	// A method insert object.
	InsertObject(ctx context.Context, in *payload.Insert_ObjectRequest, opts ...grpc.CallOption) (*payload.Object_Location, error)
	// Represent the streaming RPC to insert object by bidirectional streaming.
	StreamInsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamInsertObjectClient, error)
	// A method to insert multiple objects.
	MultiInsertObject(ctx context.Context, in *payload.Insert_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
	// A method to update object.
	UpdateObject(ctx context.Context, in *payload.Update_ObjectRequest, opts ...grpc.CallOption) (*payload.Object_Location, error)
	// A method to update object by bidirectional streaming.
	StreamUpdateObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpdateObjectClient, error)
	// A method to update multiple objects.
	MultiUpdateObject(ctx context.Context, in *payload.Update_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
	// A method to upsert object.
	UpsertObject(ctx context.Context, in *payload.Upsert_ObjectRequest, opts ...grpc.CallOption) (*payload.Object_Location, error)
	// A method to upsert object by bidirectional streaming.
	StreamUpsertObject(ctx context.Context, opts ...grpc.CallOption) (Filter_StreamUpsertObjectClient, error)
	// A method to upsert multiple objects.
	MultiUpsertObject(ctx context.Context, in *payload.Upsert_MultiObjectRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type filterClient struct {
	cc grpc.ClientConnInterface
}

func NewFilterClient(cc grpc.ClientConnInterface) FilterClient {
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
	stream, err := c.cc.NewStream(ctx, &Filter_ServiceDesc.Streams[0], "/vald.v1.Filter/StreamSearchObject", opts...)
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
	stream, err := c.cc.NewStream(ctx, &Filter_ServiceDesc.Streams[1], "/vald.v1.Filter/StreamInsertObject", opts...)
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
	stream, err := c.cc.NewStream(ctx, &Filter_ServiceDesc.Streams[2], "/vald.v1.Filter/StreamUpdateObject", opts...)
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
	stream, err := c.cc.NewStream(ctx, &Filter_ServiceDesc.Streams[3], "/vald.v1.Filter/StreamUpsertObject", opts...)
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
// All implementations must embed UnimplementedFilterServer
// for forward compatibility
type FilterServer interface {
	// A method to search object.
	SearchObject(context.Context, *payload.Search_ObjectRequest) (*payload.Search_Response, error)
	// A method to search multiple objects.
	MultiSearchObject(context.Context, *payload.Search_MultiObjectRequest) (*payload.Search_Responses, error)
	// A method to search object by bidirectional streaming.
	StreamSearchObject(Filter_StreamSearchObjectServer) error
	// A method insert object.
	InsertObject(context.Context, *payload.Insert_ObjectRequest) (*payload.Object_Location, error)
	// Represent the streaming RPC to insert object by bidirectional streaming.
	StreamInsertObject(Filter_StreamInsertObjectServer) error
	// A method to insert multiple objects.
	MultiInsertObject(context.Context, *payload.Insert_MultiObjectRequest) (*payload.Object_Locations, error)
	// A method to update object.
	UpdateObject(context.Context, *payload.Update_ObjectRequest) (*payload.Object_Location, error)
	// A method to update object by bidirectional streaming.
	StreamUpdateObject(Filter_StreamUpdateObjectServer) error
	// A method to update multiple objects.
	MultiUpdateObject(context.Context, *payload.Update_MultiObjectRequest) (*payload.Object_Locations, error)
	// A method to upsert object.
	UpsertObject(context.Context, *payload.Upsert_ObjectRequest) (*payload.Object_Location, error)
	// A method to upsert object by bidirectional streaming.
	StreamUpsertObject(Filter_StreamUpsertObjectServer) error
	// A method to upsert multiple objects.
	MultiUpsertObject(context.Context, *payload.Upsert_MultiObjectRequest) (*payload.Object_Locations, error)
	mustEmbedUnimplementedFilterServer()
}

// UnimplementedFilterServer must be embedded to have forward compatible implementations.
type UnimplementedFilterServer struct {
}

func (UnimplementedFilterServer) SearchObject(context.Context, *payload.Search_ObjectRequest) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchObject not implemented")
}
func (UnimplementedFilterServer) MultiSearchObject(context.Context, *payload.Search_MultiObjectRequest) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiSearchObject not implemented")
}
func (UnimplementedFilterServer) StreamSearchObject(Filter_StreamSearchObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearchObject not implemented")
}
func (UnimplementedFilterServer) InsertObject(context.Context, *payload.Insert_ObjectRequest) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertObject not implemented")
}
func (UnimplementedFilterServer) StreamInsertObject(Filter_StreamInsertObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamInsertObject not implemented")
}
func (UnimplementedFilterServer) MultiInsertObject(context.Context, *payload.Insert_MultiObjectRequest) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiInsertObject not implemented")
}
func (UnimplementedFilterServer) UpdateObject(context.Context, *payload.Update_ObjectRequest) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateObject not implemented")
}
func (UnimplementedFilterServer) StreamUpdateObject(Filter_StreamUpdateObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpdateObject not implemented")
}
func (UnimplementedFilterServer) MultiUpdateObject(context.Context, *payload.Update_MultiObjectRequest) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpdateObject not implemented")
}
func (UnimplementedFilterServer) UpsertObject(context.Context, *payload.Upsert_ObjectRequest) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertObject not implemented")
}
func (UnimplementedFilterServer) StreamUpsertObject(Filter_StreamUpsertObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpsertObject not implemented")
}
func (UnimplementedFilterServer) MultiUpsertObject(context.Context, *payload.Upsert_MultiObjectRequest) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpsertObject not implemented")
}
func (UnimplementedFilterServer) mustEmbedUnimplementedFilterServer() {}

// UnsafeFilterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FilterServer will
// result in compilation errors.
type UnsafeFilterServer interface {
	mustEmbedUnimplementedFilterServer()
}

func RegisterFilterServer(s grpc.ServiceRegistrar, srv FilterServer) {
	s.RegisterService(&Filter_ServiceDesc, srv)
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

// Filter_ServiceDesc is the grpc.ServiceDesc for Filter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Filter_ServiceDesc = grpc.ServiceDesc{
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
