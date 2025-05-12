//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package embedder

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

// EmbedderClient is the client API for Embedder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmbedderClient interface {
	Insert(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Location, error)
	Search(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Search_Response, error)
	Embedding(ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption) (*payload.Object_Vector, error)
	Commit(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Empty, error)
}

type embedderClient struct {
	cc grpc.ClientConnInterface
}

func NewEmbedderClient(cc grpc.ClientConnInterface) EmbedderClient {
	return &embedderClient{cc}
}

func (c *embedderClient) Insert(
	ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/embedder.v1.Embedder/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *embedderClient) Search(
	ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/embedder.v1.Embedder/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *embedderClient) Embedding(
	ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/embedder.v1.Embedder/Embedding", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *embedderClient) Commit(
	ctx context.Context, in *payload.Empty, opts ...grpc.CallOption,
) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/embedder.v1.Embedder/Commit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmbedderServer is the server API for Embedder service.
// All implementations must embed UnimplementedEmbedderServer
// for forward compatibility
type EmbedderServer interface {
	Insert(context.Context, *payload.Object_Blob) (*payload.Object_Location, error)
	Search(context.Context, *payload.Object_Blob) (*payload.Search_Response, error)
	Embedding(context.Context, *payload.Object_Blob) (*payload.Object_Vector, error)
	Commit(context.Context, *payload.Empty) (*payload.Empty, error)
	mustEmbedUnimplementedEmbedderServer()
}

// UnimplementedEmbedderServer must be embedded to have forward compatible implementations.
type UnimplementedEmbedderServer struct{}

func (UnimplementedEmbedderServer) Insert(
	context.Context, *payload.Object_Blob,
) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}

func (UnimplementedEmbedderServer) Search(
	context.Context, *payload.Object_Blob,
) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}

func (UnimplementedEmbedderServer) Embedding(
	context.Context, *payload.Object_Blob,
) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Embedding not implemented")
}

func (UnimplementedEmbedderServer) Commit(context.Context, *payload.Empty) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Commit not implemented")
}
func (UnimplementedEmbedderServer) mustEmbedUnimplementedEmbedderServer() {}

// UnsafeEmbedderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmbedderServer will
// result in compilation errors.
type UnsafeEmbedderServer interface {
	mustEmbedUnimplementedEmbedderServer()
}

func RegisterEmbedderServer(s grpc.ServiceRegistrar, srv EmbedderServer) {
	s.RegisterService(&Embedder_ServiceDesc, srv)
}

func _Embedder_Insert_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmbedderServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/embedder.v1.Embedder/Insert",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(EmbedderServer).Insert(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

func _Embedder_Search_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmbedderServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/embedder.v1.Embedder/Search",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(EmbedderServer).Search(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

func _Embedder_Embedding_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Object_Blob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmbedderServer).Embedding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/embedder.v1.Embedder/Embedding",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(EmbedderServer).Embedding(ctx, req.(*payload.Object_Blob))
	}
	return interceptor(ctx, in, info, handler)
}

func _Embedder_Commit_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmbedderServer).Commit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/embedder.v1.Embedder/Commit",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(EmbedderServer).Commit(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Embedder_ServiceDesc is the grpc.ServiceDesc for Embedder service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Embedder_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "embedder.v1.Embedder",
	HandlerType: (*EmbedderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Insert",
			Handler:    _Embedder_Insert_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _Embedder_Search_Handler,
		},
		{
			MethodName: "Embedding",
			Handler:    _Embedder_Embedding_Handler,
		},
		{
			MethodName: "Commit",
			Handler:    _Embedder_Commit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/embedder/embedder.proto",
}
