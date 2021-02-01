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

package sidecar

import (
	context "context"
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	golang_proto "github.com/golang/protobuf/proto"
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

func init() {
	proto.RegisterFile("apis/proto/v1/agent/sidecar/sidecar.proto", fileDescriptor_c78d66f1184a1433)
}
func init() {
	golang_proto.RegisterFile("apis/proto/v1/agent/sidecar/sidecar.proto", fileDescriptor_c78d66f1184a1433)
}

var fileDescriptor_c78d66f1184a1433 = []byte{
	// 190 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4c, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0x4f, 0x4c, 0x4f, 0xcd, 0x2b, 0xd1,
	0x2f, 0xce, 0x4c, 0x49, 0x4d, 0x4e, 0x2c, 0x82, 0xd1, 0x7a, 0x60, 0x69, 0x21, 0x2e, 0x18, 0xb7,
	0xcc, 0x50, 0x4a, 0x37, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x3f, 0x3d,
	0x3f, 0x3d, 0x1f, 0x62, 0x42, 0x52, 0x69, 0x1a, 0x98, 0x07, 0x31, 0x0e, 0xc4, 0x82, 0x68, 0x35,
	0xe2, 0xe4, 0x62, 0x0f, 0x86, 0x68, 0x76, 0xaa, 0x3f, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39,
	0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x0f, 0x3c, 0x96, 0x63, 0x3c, 0xf1, 0x58, 0x8e, 0x91, 0x4b, 0x39,
	0xbf, 0x28, 0x5d, 0xaf, 0x2c, 0x25, 0x31, 0xb1, 0x58, 0xaf, 0x2c, 0x31, 0x27, 0x45, 0x2f, 0xb1,
	0x20, 0x53, 0xaf, 0xcc, 0x50, 0x0f, 0xec, 0x1c, 0x3d, 0xa8, 0xbd, 0x4e, 0x02, 0x61, 0x89, 0x39,
	0x29, 0x8e, 0x20, 0x21, 0xa8, 0x61, 0x01, 0x8c, 0x51, 0x06, 0x48, 0x0e, 0x01, 0xeb, 0xd7, 0x07,
	0xe9, 0xd7, 0x07, 0xfb, 0x2a, 0xbd, 0xa8, 0x20, 0x19, 0xc3, 0x53, 0x49, 0x6c, 0x60, 0x27, 0x19,
	0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x35, 0x31, 0xcd, 0x2f, 0xfa, 0x00, 0x00, 0x00,
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
	ServiceName: "sidecar.v1.Sidecar",
	HandlerType: (*SidecarServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "apis/proto/v1/agent/sidecar/sidecar.proto",
}
