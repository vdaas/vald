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

package sidecar

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

// SidecarClient is the client API for Sidecar service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SidecarClient interface {
}

type sidecarClient struct {
	cc grpc.ClientConnInterface
}

func NewSidecarClient(cc grpc.ClientConnInterface) SidecarClient {
	return &sidecarClient{cc}
}

// SidecarServer is the server API for Sidecar service.
// All implementations must embed UnimplementedSidecarServer
// for forward compatibility
type SidecarServer interface {
	mustEmbedUnimplementedSidecarServer()
}

// UnimplementedSidecarServer must be embedded to have forward compatible implementations.
type UnimplementedSidecarServer struct {
}

func (UnimplementedSidecarServer) mustEmbedUnimplementedSidecarServer() {}

// UnsafeSidecarServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SidecarServer will
// result in compilation errors.
type UnsafeSidecarServer interface {
	mustEmbedUnimplementedSidecarServer()
}

func RegisterSidecarServer(s grpc.ServiceRegistrar, srv SidecarServer) {
	s.RegisterService(&Sidecar_ServiceDesc, srv)
}

// Sidecar_ServiceDesc is the grpc.ServiceDesc for Sidecar service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sidecar_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sidecar.v1.Sidecar",
	HandlerType: (*SidecarServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "v1/agent/sidecar/sidecar.proto",
}
