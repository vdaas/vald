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

package sidecar

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
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
	proto.RegisterFile("apis/proto/agent/sidecar/sidecar.proto", fileDescriptor_bdde3421b750feec)
}

var fileDescriptor_bdde3421b750feec = []byte{
	// 148 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4b, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x4f, 0x4c, 0x4f, 0xcd, 0x2b, 0xd1, 0x2f, 0xce, 0x4c,
	0x49, 0x4d, 0x4e, 0x2c, 0x82, 0xd1, 0x7a, 0x60, 0x39, 0x21, 0x76, 0x28, 0xd7, 0x88, 0x93, 0x8b,
	0x3d, 0x18, 0xc2, 0x74, 0xca, 0x3d, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f,
	0xe4, 0x18, 0xb9, 0x64, 0xf2, 0x8b, 0xd2, 0xf5, 0xca, 0x52, 0x12, 0x13, 0x8b, 0xf5, 0xca, 0x12,
	0x73, 0x52, 0xf4, 0xc0, 0xa6, 0xe9, 0x41, 0xb5, 0x39, 0x09, 0x84, 0x25, 0xe6, 0xa4, 0x38, 0x82,
	0x84, 0xa0, 0xba, 0x03, 0x18, 0xa3, 0x74, 0xd3, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3,
	0x73, 0xf5, 0xc1, 0x1a, 0xf5, 0x41, 0x1a, 0xf5, 0xc1, 0x2e, 0x4a, 0x2f, 0x2a, 0x48, 0x46, 0x75,
	0x50, 0x12, 0x1b, 0xd8, 0x25, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x05, 0xa7, 0xf3, 0x3c,
	0xb3, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SidecarClient is the client API for Sidecar service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SidecarClient interface {
}

type sidecarClient struct {
	cc *grpc.ClientConn
}

func NewSidecarClient(cc *grpc.ClientConn) SidecarClient {
	return &sidecarClient{cc}
}

// SidecarServer is the server API for Sidecar service.
type SidecarServer interface {
}

// UnimplementedSidecarServer can be embedded to have forward compatible implementations.
type UnimplementedSidecarServer struct {
}

func RegisterSidecarServer(s *grpc.Server, srv SidecarServer) {
	s.RegisterService(&_Sidecar_serviceDesc, srv)
}

var _Sidecar_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sidecar.Sidecar",
	HandlerType: (*SidecarServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "apis/proto/agent/sidecar/sidecar.proto",
}
