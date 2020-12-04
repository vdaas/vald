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
	// 274 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x28, 0x4a, 0x2d, 0xc8,
	0xc9, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x4f, 0xce, 0xcf, 0x2b, 0x29, 0xca, 0xcf, 0xc9,
	0x49, 0x2d, 0xd2, 0x47, 0x12, 0x8e, 0xcf, 0x4d, 0xcc, 0x4b, 0x4c, 0x4f, 0x2d, 0xd2, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0x12, 0xc6, 0x22, 0x25, 0xc5, 0x5b, 0x90, 0x58, 0x99, 0x93, 0x9f, 0x98,
	0x02, 0x51, 0x23, 0x25, 0x93, 0x9e, 0x9f, 0x9f, 0x9e, 0x93, 0xaa, 0x9f, 0x58, 0x90, 0xa9, 0x9f,
	0x98, 0x97, 0x97, 0x5f, 0x02, 0x56, 0x5d, 0x0c, 0x95, 0xe5, 0x29, 0x48, 0xd2, 0x4f, 0x2f, 0xcc,
	0x81, 0xf0, 0x8c, 0x2a, 0xb8, 0x44, 0x83, 0x10, 0x26, 0x3a, 0xc3, 0x9d, 0x20, 0x14, 0xc9, 0xc5,
	0x8f, 0x24, 0xe1, 0x99, 0x97, 0x96, 0x2f, 0xc4, 0xa7, 0x07, 0xb3, 0xc7, 0x35, 0xb7, 0xa0, 0xa4,
	0x52, 0x4a, 0x1a, 0xce, 0x47, 0x52, 0xa9, 0xe7, 0x98, 0x9e, 0x9a, 0x57, 0x52, 0xac, 0x24, 0xd9,
	0x74, 0xf9, 0xc9, 0x64, 0x26, 0x61, 0x21, 0x41, 0x64, 0xcf, 0xe8, 0x67, 0xe6, 0xa5, 0xe5, 0x4b,
	0xb1, 0x6c, 0x78, 0x20, 0xcf, 0xe4, 0x34, 0x9f, 0xf1, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4,
	0x18, 0x1f, 0x3c, 0x92, 0x63, 0xe4, 0xd2, 0xcd, 0x2f, 0x4a, 0xd7, 0x2b, 0x4b, 0x49, 0x4c, 0x2c,
	0xd6, 0x2b, 0x4b, 0xcc, 0x49, 0xd1, 0x83, 0x79, 0x1b, 0x49, 0xb7, 0x1e, 0x22, 0x84, 0x9c, 0x14,
	0xc2, 0x12, 0x73, 0x52, 0x90, 0xac, 0xf5, 0x85, 0x28, 0x47, 0x78, 0x20, 0x80, 0x31, 0xca, 0x2e,
	0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x1f, 0x6c, 0xb2, 0x3e, 0xc8, 0x64,
	0x50, 0xb0, 0x14, 0xeb, 0xa7, 0x17, 0x15, 0x24, 0xeb, 0x43, 0xed, 0xd0, 0xc7, 0x1e, 0x0b, 0x49,
	0x6c, 0xe0, 0x20, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xc5, 0xa0, 0x9e, 0x80, 0xa6, 0x01,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ReplicationControllerClient is the client API for ReplicationController service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ReplicationControllerClient interface {
	ReplicationInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Replication_Agents, error)
}

type replicationControllerClient struct {
	cc *grpc.ClientConn
}

func NewReplicationControllerClient(cc *grpc.ClientConn) ReplicationControllerClient {
	return &replicationControllerClient{cc}
}

func (c *replicationControllerClient) ReplicationInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Replication_Agents, error) {
	out := new(payload.Replication_Agents)
	err := c.cc.Invoke(ctx, "/replication_manager.ReplicationController/ReplicationInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplicationControllerServer is the server API for ReplicationController service.
type ReplicationControllerServer interface {
	ReplicationInfo(context.Context, *payload.Empty) (*payload.Replication_Agents, error)
}

// UnimplementedReplicationControllerServer can be embedded to have forward compatible implementations.
type UnimplementedReplicationControllerServer struct {
}

func (*UnimplementedReplicationControllerServer) ReplicationInfo(ctx context.Context, req *payload.Empty) (*payload.Replication_Agents, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplicationInfo not implemented")
}

func RegisterReplicationControllerServer(s *grpc.Server, srv ReplicationControllerServer) {
	s.RegisterService(&_ReplicationController_serviceDesc, srv)
}

func _ReplicationController_ReplicationInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationControllerServer).ReplicationInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/replication_manager.ReplicationController/ReplicationInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationControllerServer).ReplicationInfo(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _ReplicationController_serviceDesc = grpc.ServiceDesc{
	ServiceName: "replication_manager.ReplicationController",
	HandlerType: (*ReplicationControllerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReplicationInfo",
			Handler:    _ReplicationController_ReplicationInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "replication/controller/replication_manager.proto",
}
