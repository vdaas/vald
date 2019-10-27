//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ingress

import (
	context "context"
	fmt "fmt"
	_ "github.com/danielvladco/go-proto-gql/pb"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/vdaas/vald/apis/grpc/payload"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	math "math"
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

func init() { proto.RegisterFile("ingress_filter.proto", fileDescriptor_840104afa3c66d9d) }

var fileDescriptor_840104afa3c66d9d = []byte{
	// 189 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8e, 0xb1, 0xca, 0xc2, 0x30,
	0x14, 0x46, 0xc9, 0xf2, 0x0f, 0xe5, 0xaf, 0xa2, 0x38, 0x15, 0xed, 0x23, 0xe4, 0x82, 0xbe, 0x41,
	0x07, 0xc1, 0xcd, 0xc9, 0xc1, 0x45, 0x6e, 0x9b, 0x1a, 0x03, 0x31, 0x37, 0x26, 0xb1, 0xe0, 0x1b,
	0x3a, 0xfa, 0x08, 0xd2, 0x27, 0x91, 0x36, 0x59, 0x3a, 0x26, 0x1f, 0xe7, 0xdc, 0x93, 0xad, 0x94,
	0x91, 0xae, 0xf5, 0xfe, 0x72, 0x55, 0x3a, 0xb4, 0x8e, 0x5b, 0x47, 0x81, 0x96, 0xb3, 0xe9, 0x6f,
	0x91, 0x5b, 0x7c, 0x69, 0x42, 0x11, 0xe7, 0x62, 0x2d, 0x89, 0xa4, 0x6e, 0x01, 0xad, 0x02, 0x34,
	0x86, 0x02, 0x06, 0x45, 0xc6, 0xa7, 0xf5, 0xdf, 0xd6, 0x20, 0x1f, 0x3a, 0xbe, 0xb6, 0xf3, 0x2c,
	0x3f, 0x44, 0xd9, 0x7e, 0x74, 0x55, 0xf6, 0xdd, 0x97, 0xec, 0xd3, 0x97, 0xec, 0xdb, 0x97, 0x2c,
	0xdb, 0x90, 0x93, 0xbc, 0x13, 0x88, 0x9e, 0x77, 0xa8, 0x05, 0x4f, 0x19, 0xe9, 0x7e, 0xb5, 0x38,
	0xa1, 0x16, 0x13, 0xfe, 0xc8, 0xce, 0x5c, 0xaa, 0x70, 0x7b, 0xd6, 0xbc, 0xa1, 0x3b, 0x8c, 0x28,
	0x0c, 0xe8, 0x50, 0xe3, 0x41, 0x3a, 0xdb, 0x40, 0x94, 0x40, 0x92, 0xd4, 0x7f, 0x63, 0xc9, 0xee,
	0x17, 0x00, 0x00, 0xff, 0xff, 0xab, 0xa5, 0x88, 0x0b, 0xec, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// IngressFilterClient is the client API for IngressFilter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IngressFilterClient interface {
}

type ingressFilterClient struct {
	cc *grpc.ClientConn
}

func NewIngressFilterClient(cc *grpc.ClientConn) IngressFilterClient {
	return &ingressFilterClient{cc}
}

// IngressFilterServer is the server API for IngressFilter service.
type IngressFilterServer interface {
}

// UnimplementedIngressFilterServer can be embedded to have forward compatible implementations.
type UnimplementedIngressFilterServer struct {
}

func RegisterIngressFilterServer(s *grpc.Server, srv IngressFilterServer) {
	s.RegisterService(&_IngressFilter_serviceDesc, srv)
}

var _IngressFilter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ingress_filter.IngressFilter",
	HandlerType: (*IngressFilterServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "ingress_filter.proto",
}
