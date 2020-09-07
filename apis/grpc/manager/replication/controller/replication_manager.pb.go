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
	proto.RegisterFile("apis/proto/manager/replication/controller/replication_manager.proto", fileDescriptor_1ce88998074a19aa)
}

var fileDescriptor_1ce88998074a19aa = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x72, 0x4e, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0xcc, 0x4b, 0x4c, 0x4f, 0x2d, 0xd2, 0x2f,
	0x4a, 0x2d, 0xc8, 0xc9, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x4f, 0xce, 0xcf, 0x2b, 0x29,
	0xca, 0xcf, 0xc9, 0x41, 0x15, 0x8e, 0x87, 0x2a, 0xd5, 0x03, 0x6b, 0x14, 0x12, 0xc6, 0x22, 0x25,
	0x65, 0x96, 0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c, 0x9f, 0xab, 0x5f, 0x96, 0x92, 0x98,
	0x58, 0xac, 0x5f, 0x96, 0x98, 0x93, 0xa2, 0x8f, 0x64, 0x5f, 0x41, 0x62, 0x65, 0x4e, 0x7e, 0x62,
	0x0a, 0x8c, 0x86, 0x18, 0x26, 0x25, 0x93, 0x9e, 0x9f, 0x9f, 0x9e, 0x93, 0x0a, 0x52, 0xa8, 0x9f,
	0x98, 0x97, 0x97, 0x5f, 0x02, 0x36, 0xb6, 0x18, 0x22, 0x6b, 0x94, 0xc1, 0xc5, 0x1d, 0x84, 0xb0,
	0x4c, 0x28, 0x92, 0x8b, 0x1f, 0x89, 0xeb, 0x99, 0x97, 0x96, 0x2f, 0xc4, 0xa7, 0x07, 0x33, 0xcf,
	0x35, 0xb7, 0xa0, 0xa4, 0x52, 0x4a, 0x1a, 0xce, 0x47, 0x52, 0xa9, 0xe7, 0x98, 0x9e, 0x9a, 0x57,
	0x52, 0xac, 0x24, 0xd9, 0x74, 0xf9, 0xc9, 0x64, 0x26, 0x61, 0x21, 0x41, 0x14, 0x4f, 0x67, 0xe6,
	0xa5, 0xe5, 0x3b, 0xcd, 0x67, 0x3c, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f,
	0xe4, 0x18, 0xb9, 0x74, 0xf3, 0x8b, 0xd2, 0xf5, 0xc0, 0xfe, 0xd0, 0x03, 0xf9, 0x43, 0x0f, 0x16,
	0x02, 0x48, 0xfa, 0xf4, 0x10, 0x81, 0xe5, 0xa4, 0x10, 0x96, 0x98, 0x93, 0x82, 0x64, 0xa1, 0x2f,
	0x44, 0xb9, 0x33, 0x5c, 0x45, 0x00, 0x63, 0x94, 0x1d, 0x9e, 0x10, 0x4a, 0x2f, 0x2a, 0x48, 0x26,
	0x10, 0x21, 0x49, 0x6c, 0xe0, 0x20, 0x31, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x57, 0xc2, 0xa3,
	0xe2, 0xc4, 0x01, 0x00, 0x00,
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
	Metadata: "apis/proto/manager/replication/controller/replication_manager.proto",
}
