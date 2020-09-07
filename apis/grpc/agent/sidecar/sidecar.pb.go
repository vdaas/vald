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
	_ "github.com/vdaas/vald/apis/grpc/payload"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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
	// 188 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4b, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x4f, 0x4c, 0x4f, 0xcd, 0x2b, 0xd1, 0x2f, 0xce, 0x4c,
	0x49, 0x4d, 0x4e, 0x2c, 0x82, 0xd1, 0x7a, 0x60, 0x39, 0x21, 0x76, 0x28, 0x57, 0xca, 0x2c, 0x3d,
	0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0xbf, 0x2c, 0x25, 0x31, 0xb1, 0x58, 0xbf,
	0x2c, 0x31, 0x27, 0x45, 0x1f, 0xc9, 0x98, 0x82, 0xc4, 0xca, 0x9c, 0xfc, 0xc4, 0x14, 0x18, 0x0d,
	0x31, 0x40, 0x4a, 0x26, 0x3d, 0x3f, 0x3f, 0x3d, 0x27, 0x15, 0xa4, 0x50, 0x3f, 0x31, 0x2f, 0x2f,
	0xbf, 0x24, 0xb1, 0x24, 0x33, 0x3f, 0xaf, 0x18, 0x22, 0x6b, 0xc4, 0xc9, 0xc5, 0x1e, 0x0c, 0xb1,
	0xc0, 0x29, 0xf7, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0xe4,
	0x92, 0xc9, 0x2f, 0x4a, 0xd7, 0x03, 0xdb, 0xa3, 0x07, 0xb2, 0x47, 0x0f, 0xec, 0x46, 0x3d, 0xa8,
	0x63, 0x9c, 0x04, 0xc2, 0x12, 0x73, 0x52, 0x1c, 0x41, 0x42, 0x50, 0xdd, 0x01, 0x8c, 0x51, 0xba,
	0x78, 0x1c, 0x98, 0x5e, 0x54, 0x90, 0x8c, 0xea, 0xcd, 0x24, 0x36, 0xb0, 0x03, 0x8c, 0x01, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x50, 0x37, 0x79, 0xba, 0x09, 0x01, 0x00, 0x00,
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
