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

package replication

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

func init() {
	proto.RegisterFile("replication/replication_manager.proto", fileDescriptor_f84c0b45f233a237)
}

var fileDescriptor_f84c0b45f233a237 = []byte{
	// 195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2d, 0x4a, 0x2d, 0xc8,
	0xc9, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x47, 0x62, 0xc7, 0xe7, 0x26, 0xe6, 0x25, 0xa6,
	0xa7, 0x16, 0xe9, 0x15, 0x14, 0xe5, 0x97, 0xe4, 0x0b, 0x09, 0x63, 0x91, 0x92, 0xe2, 0x2d, 0x48,
	0xac, 0xcc, 0xc9, 0x4f, 0x4c, 0x81, 0xa8, 0x91, 0x92, 0x49, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5,
	0x4f, 0x2c, 0xc8, 0xd4, 0x4f, 0xcc, 0xcb, 0xcb, 0x2f, 0x01, 0xab, 0x2e, 0x86, 0xca, 0xf2, 0x14,
	0x24, 0xe9, 0xa7, 0x17, 0xe6, 0x40, 0x78, 0x46, 0xbc, 0x5c, 0xdc, 0x41, 0x08, 0x13, 0x9d, 0xea,
	0x4f, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0x46, 0x2e, 0xa5, 0xfc,
	0xa2, 0x74, 0xbd, 0xb2, 0x94, 0xc4, 0xc4, 0x62, 0xbd, 0xb2, 0xc4, 0x9c, 0x14, 0x3d, 0x98, 0x53,
	0x90, 0xdc, 0xe0, 0x24, 0x16, 0x96, 0x98, 0x93, 0x82, 0x64, 0x84, 0x2f, 0x44, 0x4d, 0x00, 0x63,
	0x94, 0x71, 0x7a, 0x66, 0x49, 0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae, 0x3e, 0xd8, 0x10, 0x7d,
	0x90, 0x21, 0x20, 0x57, 0x15, 0xeb, 0xa7, 0x17, 0x15, 0x24, 0xeb, 0x43, 0x8d, 0x43, 0xf6, 0x6d,
	0x12, 0x1b, 0xd8, 0x59, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb8, 0x7c, 0x98, 0x2c, 0x0f,
	0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ReplicationClient is the client API for Replication service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ReplicationClient interface {
}

type replicationClient struct {
	cc *grpc.ClientConn
}

func NewReplicationClient(cc *grpc.ClientConn) ReplicationClient {
	return &replicationClient{cc}
}

// ReplicationServer is the server API for Replication service.
type ReplicationServer interface {
}

// UnimplementedReplicationServer can be embedded to have forward compatible implementations.
type UnimplementedReplicationServer struct {
}

func RegisterReplicationServer(s *grpc.Server, srv ReplicationServer) {
	s.RegisterService(&_Replication_serviceDesc, srv)
}

var _Replication_serviceDesc = grpc.ServiceDesc{
	ServiceName: "replication_manager.Replication",
	HandlerType: (*ReplicationServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "replication/replication_manager.proto",
}
