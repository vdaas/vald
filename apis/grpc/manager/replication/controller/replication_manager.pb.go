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

package controller

import (
	context "context"
	fmt "fmt"
	math "math"

	_ "github.com/danielvladco/go-proto-gql/pb"
	proto "github.com/gogo/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/payload"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
	proto.RegisterFile("replication/controller/replication_manager.proto", fileDescriptor_ca7dddd9e8833d57)
}

var fileDescriptor_ca7dddd9e8833d57 = []byte{
	// 271 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x28, 0x4a, 0x2d, 0xc8,
	0xc9, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x4f, 0xce, 0xcf, 0x2b, 0x29, 0xca, 0xcf, 0xc9,
	0x49, 0x2d, 0xd2, 0x47, 0x12, 0x8e, 0xcf, 0x4d, 0xcc, 0x4b, 0x4c, 0x4f, 0x2d, 0xd2, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0x12, 0xc6, 0x22, 0x25, 0xc5, 0x5b, 0x90, 0x58, 0x99, 0x93, 0x9f, 0x98,
	0x02, 0x51, 0x23, 0x25, 0x93, 0x9e, 0x9f, 0x9f, 0x9e, 0x93, 0xaa, 0x9f, 0x58, 0x90, 0xa9, 0x9f,
	0x98, 0x97, 0x97, 0x5f, 0x02, 0x56, 0x5d, 0x0c, 0x95, 0xe5, 0x29, 0x48, 0xd2, 0x4f, 0x2f, 0xcc,
	0x81, 0xf0, 0x8c, 0xf2, 0xb8, 0xb8, 0x83, 0x10, 0x26, 0x0a, 0x45, 0x72, 0xf1, 0x23, 0x71, 0x3d,
	0xf3, 0xd2, 0xf2, 0x85, 0xf8, 0xf4, 0x60, 0xa6, 0xbb, 0xe6, 0x16, 0x94, 0x54, 0x4a, 0x49, 0xc3,
	0xf9, 0x48, 0x2a, 0xf5, 0x1c, 0xd3, 0x53, 0xf3, 0x4a, 0x8a, 0x95, 0x24, 0x9b, 0x2e, 0x3f, 0x99,
	0xcc, 0x24, 0x2c, 0x24, 0x88, 0xec, 0x05, 0xfd, 0xcc, 0xbc, 0xb4, 0x7c, 0x29, 0x96, 0x0d, 0x0f,
	0xe4, 0x99, 0x9c, 0xe6, 0x33, 0x9e, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47,
	0x72, 0x8c, 0x5c, 0xba, 0xf9, 0x45, 0xe9, 0x7a, 0x65, 0x29, 0x89, 0x89, 0xc5, 0x7a, 0x65, 0x89,
	0x39, 0x29, 0x7a, 0x30, 0xcf, 0x22, 0xe9, 0xd6, 0x43, 0x84, 0x8b, 0x93, 0x42, 0x58, 0x62, 0x4e,
	0x0a, 0x92, 0xb5, 0xbe, 0x10, 0xe5, 0xce, 0x70, 0x15, 0x01, 0x8c, 0x51, 0x76, 0xe9, 0x99, 0x25,
	0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0x60, 0x93, 0xf5, 0x41, 0x26, 0x83, 0x02, 0xa3,
	0x58, 0x3f, 0xbd, 0xa8, 0x20, 0x59, 0x1f, 0x6a, 0x87, 0x3e, 0xf6, 0xb0, 0x4f, 0x62, 0x03, 0x07,
	0x8c, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xd7, 0xde, 0x29, 0x2e, 0x9c, 0x01, 0x00, 0x00,
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
	ReplicationInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Replication_Agents, error)
}

type replicationClient struct {
	cc *grpc.ClientConn
}

func NewReplicationClient(cc *grpc.ClientConn) ReplicationClient {
	return &replicationClient{cc}
}

func (c *replicationClient) ReplicationInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Replication_Agents, error) {
	out := new(payload.Replication_Agents)
	err := c.cc.Invoke(ctx, "/replication_manager.Replication/ReplicationInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplicationServer is the server API for Replication service.
type ReplicationServer interface {
	ReplicationInfo(context.Context, *payload.Empty) (*payload.Replication_Agents, error)
}

// UnimplementedReplicationServer can be embedded to have forward compatible implementations.
type UnimplementedReplicationServer struct {
}

func (*UnimplementedReplicationServer) ReplicationInfo(ctx context.Context, req *payload.Empty) (*payload.Replication_Agents, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplicationInfo not implemented")
}

func RegisterReplicationServer(s *grpc.Server, srv ReplicationServer) {
	s.RegisterService(&_Replication_serviceDesc, srv)
}

func _Replication_ReplicationInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationServer).ReplicationInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/replication_manager.Replication/ReplicationInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServer).ReplicationInfo(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Replication_serviceDesc = grpc.ServiceDesc{
	ServiceName: "replication_manager.Replication",
	HandlerType: (*ReplicationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReplicationInfo",
			Handler:    _Replication_ReplicationInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "replication/controller/replication_manager.proto",
}
