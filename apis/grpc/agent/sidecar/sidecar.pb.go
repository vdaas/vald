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

	_ "github.com/danielvladco/go-proto-gql/pb"
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

func init() { proto.RegisterFile("sidecar/sidecar.proto", fileDescriptor_a79f12d5eccb8a6a) }

var fileDescriptor_a79f12d5eccb8a6a = []byte{
	// 188 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0xce, 0x4c, 0x49,
	0x4d, 0x4e, 0x2c, 0xd2, 0x87, 0xd2, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0xec, 0x50, 0xae,
	0x14, 0x6f, 0x41, 0x62, 0x65, 0x4e, 0x7e, 0x62, 0x0a, 0x44, 0x5c, 0x4a, 0x26, 0x3d, 0x3f, 0x3f,
	0x3d, 0x27, 0x55, 0x3f, 0xb1, 0x20, 0x53, 0x3f, 0x31, 0x2f, 0x2f, 0xbf, 0x24, 0xb1, 0x24, 0x33,
	0x3f, 0xaf, 0x18, 0x2a, 0xcb, 0x53, 0x90, 0xa4, 0x9f, 0x5e, 0x98, 0x03, 0xe1, 0x19, 0x71, 0x72,
	0xb1, 0x07, 0x43, 0x4c, 0x71, 0xca, 0x3d, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07,
	0x8f, 0xe4, 0x18, 0xb9, 0x64, 0xf2, 0x8b, 0xd2, 0xf5, 0xca, 0x52, 0x12, 0x13, 0x8b, 0xf5, 0xca,
	0x12, 0x73, 0x52, 0xf4, 0x12, 0xd3, 0x53, 0xf3, 0x4a, 0xf4, 0xa0, 0x36, 0x3a, 0x09, 0x84, 0x25,
	0xe6, 0xa4, 0x38, 0x82, 0x84, 0xa0, 0xba, 0x03, 0x18, 0xa3, 0x74, 0xd3, 0x33, 0x4b, 0x32, 0x4a,
	0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0xc1, 0x1a, 0xf5, 0x41, 0x1a, 0x41, 0xae, 0x28, 0xd6, 0x4f,
	0x2f, 0x2a, 0x48, 0xd6, 0x07, 0x1b, 0x01, 0xf3, 0x43, 0x12, 0x1b, 0xd8, 0x01, 0xc6, 0x80, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x1c, 0x2f, 0x76, 0x1a, 0xdd, 0x00, 0x00, 0x00,
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
	Metadata:    "sidecar/sidecar.proto",
}
