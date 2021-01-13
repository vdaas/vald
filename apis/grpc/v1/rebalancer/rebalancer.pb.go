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

package rebalancer

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/gogo/protobuf/proto"
	_ "github.com/vdaas/vald/apis/grpc/v1/payload"
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
	proto.RegisterFile("apis/proto/v1/rebalancer/rebalancer.proto", fileDescriptor_c446ac80040fa7f9)
}

var fileDescriptor_c446ac80040fa7f9 = []byte{
	// 177 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4c, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0x2f, 0x4a, 0x4d, 0x4a, 0xcc, 0x49,
	0xcc, 0x4b, 0x4e, 0x2d, 0x42, 0x62, 0xea, 0x81, 0xa5, 0x85, 0x78, 0x91, 0x44, 0xca, 0x0c, 0xa5,
	0x94, 0x51, 0x75, 0x16, 0x24, 0x56, 0xe6, 0xe4, 0x27, 0xa6, 0xc0, 0x68, 0x88, 0x1e, 0x23, 0x1e,
	0x2e, 0x2e, 0xe7, 0xfc, 0xbc, 0x92, 0xa2, 0xfc, 0x9c, 0x9c, 0xd4, 0x22, 0x23, 0x56, 0x2e, 0x66,
	0xaf, 0xfc, 0x24, 0xa7, 0xec, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48,
	0x8e, 0x91, 0x4b, 0x21, 0xbf, 0x28, 0x5d, 0xaf, 0x2c, 0x25, 0x31, 0xb1, 0x58, 0xaf, 0x2c, 0x31,
	0x27, 0x45, 0x2f, 0xb1, 0x20, 0x53, 0xaf, 0xcc, 0x50, 0x0f, 0x61, 0x99, 0x13, 0x57, 0x10, 0x9c,
	0x1d, 0xc0, 0x18, 0xa5, 0x9b, 0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c, 0x9f, 0xab, 0x0f,
	0xd6, 0xa6, 0x0f, 0xd2, 0xa6, 0x0f, 0x76, 0x4d, 0x7a, 0x51, 0x41, 0x32, 0xaa, 0x37, 0x92, 0xd8,
	0xc0, 0x0e, 0x31, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x6b, 0xe8, 0x23, 0xd4, 0xe9, 0x00, 0x00,
	0x00,
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
	ServiceName: "rebalancer.v1.Controller",
	HandlerType: (*ControllerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "apis/proto/v1/rebalancer/rebalancer.proto",
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
	ServiceName: "rebalancer.v1.Job",
	HandlerType: (*JobServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "apis/proto/v1/rebalancer/rebalancer.proto",
}
