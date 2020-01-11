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

package index

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

func init() { proto.RegisterFile("index/index_manager.proto", fileDescriptor_11357116787cb271) }

var fileDescriptor_11357116787cb271 = []byte{
	// 190 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcc, 0xcc, 0x4b, 0x49,
	0xad, 0xd0, 0x07, 0x93, 0xf1, 0xb9, 0x89, 0x79, 0x89, 0xe9, 0xa9, 0x45, 0x7a, 0x05, 0x45, 0xf9,
	0x25, 0xf9, 0x42, 0xbc, 0x28, 0x82, 0x52, 0xbc, 0x05, 0x89, 0x95, 0x39, 0xf9, 0x89, 0x29, 0x10,
	0x59, 0x29, 0x99, 0xf4, 0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0xfd, 0xc4, 0x82, 0x4c, 0xfd, 0xc4, 0xbc,
	0xbc, 0xfc, 0x92, 0xc4, 0x92, 0xcc, 0xfc, 0xbc, 0x62, 0xa8, 0x2c, 0x4f, 0x41, 0x92, 0x7e, 0x7a,
	0x61, 0x0e, 0x84, 0x67, 0xc4, 0xce, 0xc5, 0xea, 0x09, 0x32, 0xcb, 0x29, 0xf7, 0xc4, 0x23, 0x39,
	0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0xe4, 0x92, 0xc9, 0x2f, 0x4a, 0xd7, 0x2b,
	0x4b, 0x49, 0x4c, 0x2c, 0xd6, 0x2b, 0x4b, 0xcc, 0x49, 0xd1, 0x83, 0x59, 0x0f, 0xb6, 0xd7, 0x49,
	0x20, 0x2c, 0x31, 0x27, 0x05, 0xcc, 0xf4, 0x85, 0x88, 0x07, 0x30, 0x46, 0xe9, 0xa6, 0x67, 0x96,
	0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea, 0x83, 0x35, 0xea, 0x83, 0x34, 0x82, 0xdc, 0x50,
	0xac, 0x9f, 0x5e, 0x54, 0x90, 0xac, 0x0f, 0x35, 0x02, 0xe2, 0x9f, 0x24, 0x36, 0xb0, 0xf5, 0xc6,
	0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x69, 0x1c, 0xae, 0xa3, 0xe5, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// IndexClient is the client API for Index service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IndexClient interface {
}

type indexClient struct {
	cc *grpc.ClientConn
}

func NewIndexClient(cc *grpc.ClientConn) IndexClient {
	return &indexClient{cc}
}

// IndexServer is the server API for Index service.
type IndexServer interface {
}

// UnimplementedIndexServer can be embedded to have forward compatible implementations.
type UnimplementedIndexServer struct {
}

func RegisterIndexServer(s *grpc.Server, srv IndexServer) {
	s.RegisterService(&_Index_serviceDesc, srv)
}

var _Index_serviceDesc = grpc.ServiceDesc{
	ServiceName: "index_manager.Index",
	HandlerType: (*IndexServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "index/index_manager.proto",
}
