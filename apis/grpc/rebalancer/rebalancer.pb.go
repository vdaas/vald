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

package rebalancer

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

func init() { proto.RegisterFile("rebalancer.proto", fileDescriptor_060d7ad7e5fab0a2) }

var fileDescriptor_060d7ad7e5fab0a2 = []byte{
	// 182 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x8d, 0x4d, 0xaa, 0xc2, 0x30,
	0x14, 0x46, 0x09, 0x8f, 0xe7, 0x20, 0x28, 0x88, 0x33, 0x8b, 0x74, 0x01, 0x0e, 0x12, 0xa8, 0x3b,
	0xa8, 0x33, 0x47, 0xe2, 0xd0, 0xd9, 0xcd, 0x0f, 0x31, 0x10, 0x73, 0xc3, 0x4d, 0x2c, 0xb8, 0x43,
	0x87, 0x2e, 0x41, 0xba, 0x12, 0x69, 0x0b, 0x76, 0xf6, 0xc1, 0x39, 0x1f, 0x87, 0xaf, 0xc9, 0x2a,
	0x08, 0x10, 0xb5, 0x25, 0x91, 0x08, 0x0b, 0x6e, 0xb8, 0xf1, 0x59, 0x63, 0x67, 0xc9, 0x52, 0xb5,
	0x4a, 0xf0, 0x0c, 0x08, 0x66, 0x42, 0xd5, 0xce, 0x21, 0xba, 0x60, 0x25, 0x24, 0x2f, 0x21, 0x46,
	0x2c, 0x50, 0x3c, 0xc6, 0x3c, 0xd1, 0x66, 0xc9, 0xf9, 0x11, 0x63, 0x21, 0x0c, 0xc1, 0x52, 0xf3,
	0xcf, 0xff, 0x4e, 0xa8, 0x5a, 0x78, 0xf5, 0x35, 0x7b, 0xf7, 0x35, 0xfb, 0xf4, 0x35, 0xe3, 0x5b,
	0x24, 0x27, 0x3a, 0x03, 0x90, 0x45, 0x07, 0xc1, 0x88, 0x39, 0xde, 0xf2, 0xcb, 0x6f, 0x9f, 0xd9,
	0x75, 0xef, 0x7c, 0xb9, 0x3d, 0x94, 0xd0, 0x78, 0x97, 0xa3, 0x2f, 0x07, 0x7f, 0x08, 0x67, 0xe9,
	0x28, 0x69, 0x39, 0x3f, 0xd5, 0x62, 0xcc, 0x1f, 0xbe, 0x01, 0x00, 0x00, 0xff, 0xff, 0xfb, 0x2c,
	0x25, 0xbf, 0xcb, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ControllerClient is the client API for Controller service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ControllerClient interface {
}

type controllerClient struct {
	cc *grpc.ClientConn
}

func NewControllerClient(cc *grpc.ClientConn) ControllerClient {
	return &controllerClient{cc}
}

// ControllerServer is the server API for Controller service.
type ControllerServer interface {
}

// UnimplementedControllerServer can be embedded to have forward compatible implementations.
type UnimplementedControllerServer struct {
}

func RegisterControllerServer(s *grpc.Server, srv ControllerServer) {
	s.RegisterService(&_Controller_serviceDesc, srv)
}

var _Controller_serviceDesc = grpc.ServiceDesc{
	ServiceName: "discoverer.Controller",
	HandlerType: (*ControllerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "rebalancer.proto",
}

// JobClient is the client API for Job service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type JobClient interface {
}

type jobClient struct {
	cc *grpc.ClientConn
}

func NewJobClient(cc *grpc.ClientConn) JobClient {
	return &jobClient{cc}
}

// JobServer is the server API for Job service.
type JobServer interface {
}

// UnimplementedJobServer can be embedded to have forward compatible implementations.
type UnimplementedJobServer struct {
}

func RegisterJobServer(s *grpc.Server, srv JobServer) {
	s.RegisterService(&_Job_serviceDesc, srv)
}

var _Job_serviceDesc = grpc.ServiceDesc{
	ServiceName: "discoverer.Job",
	HandlerType: (*JobServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "rebalancer.proto",
}
