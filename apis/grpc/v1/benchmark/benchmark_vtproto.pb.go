//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

package benchmark

import (
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

// ControllerClient is the client API for Controller service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ControllerClient interface {
}

type controllerClient struct {
	cc grpc.ClientConnInterface
}

func NewControllerClient(cc grpc.ClientConnInterface) ControllerClient {
	return &controllerClient{cc}
}

// ControllerServer is the server API for Controller service.
// All implementations must embed UnimplementedControllerServer
// for forward compatibility
type ControllerServer interface {
	mustEmbedUnimplementedControllerServer()
}

// UnimplementedControllerServer must be embedded to have forward compatible implementations.
type UnimplementedControllerServer struct {
}

func (UnimplementedControllerServer) mustEmbedUnimplementedControllerServer() {}

// UnsafeControllerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ControllerServer will
// result in compilation errors.
type UnsafeControllerServer interface {
	mustEmbedUnimplementedControllerServer()
}

func RegisterControllerServer(s grpc.ServiceRegistrar, srv ControllerServer) {
	s.RegisterService(&Controller_ServiceDesc, srv)
}

// Controller_ServiceDesc is the grpc.ServiceDesc for Controller service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Controller_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "benchmark.v1.Controller",
	HandlerType: (*ControllerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "apis/proto/v1/benchmark/benchmark.proto",
}

// SearchJobClient is the client API for SearchJob service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchJobClient interface {
}

type searchJobClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchJobClient(cc grpc.ClientConnInterface) SearchJobClient {
	return &searchJobClient{cc}
}

// SearchJobServer is the server API for SearchJob service.
// All implementations must embed UnimplementedSearchJobServer
// for forward compatibility
type SearchJobServer interface {
	mustEmbedUnimplementedSearchJobServer()
}

// UnimplementedSearchJobServer must be embedded to have forward compatible implementations.
type UnimplementedSearchJobServer struct {
}

func (UnimplementedSearchJobServer) mustEmbedUnimplementedSearchJobServer() {}

// UnsafeSearchJobServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchJobServer will
// result in compilation errors.
type UnsafeSearchJobServer interface {
	mustEmbedUnimplementedSearchJobServer()
}

func RegisterSearchJobServer(s grpc.ServiceRegistrar, srv SearchJobServer) {
	s.RegisterService(&SearchJob_ServiceDesc, srv)
}

// SearchJob_ServiceDesc is the grpc.ServiceDesc for SearchJob service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SearchJob_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "benchmark.v1.SearchJob",
	HandlerType: (*SearchJobServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "apis/proto/v1/benchmark/benchmark.proto",
}
