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

package replication

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

func init() { proto.RegisterFile("replication_manager.proto", fileDescriptor_f983d5516f249038) }

var fileDescriptor_f983d5516f249038 = []byte{
	// 193 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2c, 0x4a, 0x2d, 0xc8,
	0xc9, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0x8b, 0xcf, 0x4d, 0xcc, 0x4b, 0x4c, 0x4f, 0x2d, 0xd2,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0xc6, 0x22, 0x25, 0xc5, 0x5b, 0x90, 0x58, 0x99, 0x93,
	0x9f, 0x98, 0x02, 0x51, 0x23, 0x25, 0x93, 0x9e, 0x9f, 0x9f, 0x9e, 0x93, 0xaa, 0x9f, 0x58, 0x90,
	0xa9, 0x9f, 0x98, 0x97, 0x97, 0x5f, 0x02, 0x56, 0x5d, 0x0c, 0x95, 0xe5, 0x29, 0x48, 0xd2, 0x4f,
	0x2f, 0xcc, 0x81, 0xf0, 0x8c, 0x78, 0xb9, 0xb8, 0x83, 0x10, 0x26, 0x3a, 0xd5, 0x9f, 0x78, 0x24,
	0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x5c, 0x4a, 0xf9, 0x45, 0xe9, 0x7a,
	0x65, 0x29, 0x89, 0x89, 0xc5, 0x7a, 0x65, 0x89, 0x39, 0x29, 0x7a, 0x30, 0xa7, 0x20, 0xb9, 0xc1,
	0x49, 0x2c, 0x2c, 0x31, 0x27, 0x05, 0xc9, 0x08, 0x5f, 0x88, 0x9a, 0x00, 0xc6, 0x28, 0xe3, 0xf4,
	0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d, 0xb0, 0x21, 0xfa, 0x20, 0x43, 0x40,
	0xae, 0x2a, 0xd6, 0x4f, 0x2f, 0x2a, 0x48, 0xd6, 0x87, 0x1a, 0xa7, 0x8f, 0x64, 0x5c, 0x12, 0x1b,
	0xd8, 0x59, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xae, 0x4e, 0x0c, 0xe8, 0x03, 0x01, 0x00,
	0x00,
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
	Metadata:    "replication_manager.proto",
}
