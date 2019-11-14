//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
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

func init() { proto.RegisterFile("ingress/ingress_filter.proto", fileDescriptor_8f5342c46835d3ee) }

var fileDescriptor_8f5342c46835d3ee = []byte{
	// 192 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xc9, 0xcc, 0x4b, 0x2f,
	0x4a, 0x2d, 0x2e, 0xd6, 0x87, 0xd2, 0xf1, 0x69, 0x99, 0x39, 0x25, 0xa9, 0x45, 0x7a, 0x05, 0x45,
	0xf9, 0x25, 0xf9, 0x42, 0x7c, 0xa8, 0xa2, 0x52, 0xbc, 0x05, 0x89, 0x95, 0x39, 0xf9, 0x89, 0x29,
	0x10, 0x69, 0x29, 0x99, 0xf4, 0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0xfd, 0xc4, 0x82, 0x4c, 0xfd, 0xc4,
	0xbc, 0xbc, 0xfc, 0x92, 0xc4, 0x92, 0xcc, 0xfc, 0xbc, 0x62, 0xa8, 0x2c, 0x4f, 0x41, 0x92, 0x7e,
	0x7a, 0x61, 0x0e, 0x84, 0x67, 0xc4, 0xcf, 0xc5, 0xeb, 0x09, 0x31, 0xcc, 0x0d, 0x6c, 0x96, 0x53,
	0xc1, 0x89, 0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0xc8, 0x25, 0x9b,
	0x5f, 0x94, 0xae, 0x57, 0x96, 0x92, 0x98, 0x58, 0xac, 0x57, 0x96, 0x98, 0x93, 0xa2, 0x07, 0x75,
	0x06, 0xd4, 0x7e, 0x27, 0xc1, 0xb0, 0xc4, 0x9c, 0x14, 0x14, 0xfd, 0x01, 0x8c, 0x51, 0x7a, 0xe9,
	0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0x60, 0xad, 0xfa, 0x20, 0xad, 0x20,
	0xd7, 0x14, 0xeb, 0xa7, 0x17, 0x15, 0x24, 0xeb, 0x43, 0x0c, 0x81, 0x79, 0x2d, 0x89, 0x0d, 0xec,
	0x12, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x72, 0x78, 0xf1, 0xa6, 0xf4, 0x00, 0x00, 0x00,
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
	Metadata:    "ingress/ingress_filter.proto",
}
