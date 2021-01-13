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

package agent

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/gogo/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
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
	proto.RegisterFile("apis/proto/v1/manager/replication/agent/replication_manager.proto", fileDescriptor_e8f74170057978aa)
}

var fileDescriptor_e8f74170057978aa = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xc1, 0x4a, 0xfb, 0x30,
	0x1c, 0xc7, 0xe9, 0x0e, 0xff, 0x3f, 0xab, 0x27, 0x83, 0x4c, 0xad, 0xa3, 0x6a, 0xbd, 0x27, 0x4c,
	0xaf, 0x5e, 0x36, 0xf0, 0xe0, 0x41, 0x90, 0x1d, 0x3c, 0x88, 0x30, 0x7e, 0x6b, 0xb3, 0x18, 0xc8,
	0xf2, 0x0b, 0x69, 0x0c, 0xec, 0xea, 0x2b, 0x78, 0xf3, 0x25, 0x7c, 0x0d, 0x8f, 0x82, 0x2f, 0x20,
	0xc3, 0x07, 0x91, 0x65, 0x2d, 0x76, 0xba, 0x79, 0x2a, 0x6d, 0xbf, 0xdf, 0xcf, 0x87, 0xf6, 0x9b,
	0xb8, 0x0f, 0x46, 0x96, 0xcc, 0x58, 0x74, 0xc8, 0x7c, 0x8f, 0x4d, 0x41, 0x83, 0xe0, 0x96, 0x59,
	0x6e, 0x94, 0xcc, 0xc1, 0x49, 0xd4, 0x0c, 0x04, 0xd7, 0xae, 0xf9, 0x64, 0x54, 0xa5, 0x68, 0xa8,
	0x91, 0x6e, 0x7d, 0xdb, 0x88, 0xd0, 0x50, 0xa2, 0xbe, 0x97, 0x9c, 0xac, 0x0a, 0x0c, 0xcc, 0x14,
	0x42, 0x51, 0x5f, 0x97, 0x88, 0xa4, 0x2b, 0x10, 0x85, 0xe2, 0x0c, 0x8c, 0x64, 0xa0, 0x35, 0xba,
	0x00, 0x29, 0x97, 0x6f, 0x4f, 0x5f, 0x5a, 0xf1, 0xd6, 0xf0, 0x9b, 0x4d, 0xee, 0xe2, 0xff, 0x43,
	0x9e, 0xa3, 0xe7, 0x96, 0x1c, 0xd1, 0x1a, 0xe4, 0x7b, 0xb4, 0x91, 0xa1, 0x55, 0x60, 0x96, 0x6c,
	0x37, 0x13, 0x17, 0x53, 0xe3, 0x66, 0x59, 0xf7, 0xf1, 0xfd, 0xf3, 0xa9, 0xd5, 0xc9, 0x76, 0x56,
	0x3e, 0xd3, 0x56, 0x48, 0x88, 0xdb, 0x43, 0x3e, 0x06, 0x05, 0x3a, 0xe7, 0xe4, 0x78, 0x33, 0xbf,
	0x8a, 0xac, 0x13, 0xa4, 0x41, 0xb0, 0x97, 0x75, 0x7e, 0x08, 0x6a, 0xea, 0x28, 0x6e, 0xf7, 0x17,
	0xff, 0xe7, 0x52, 0x4f, 0x90, 0xfc, 0xee, 0x27, 0xe9, 0x26, 0x6b, 0x68, 0x95, 0xd9, 0x61, 0xe0,
	0xef, 0x93, 0xdd, 0x35, 0x3b, 0x49, 0x3d, 0xc1, 0xc1, 0x73, 0xf4, 0x3a, 0x4f, 0xa3, 0xb7, 0x79,
	0x1a, 0x7d, 0xcc, 0xd3, 0x28, 0x66, 0x68, 0x05, 0xf5, 0x05, 0x40, 0x49, 0x3d, 0xa8, 0x82, 0x82,
	0x91, 0x0b, 0xf2, 0xc6, 0xdd, 0x06, 0x07, 0x37, 0xa0, 0x8a, 0x86, 0xf8, 0x6a, 0x99, 0x0c, 0xfa,
	0xeb, 0xe8, 0xf6, 0x5c, 0x48, 0x77, 0xff, 0x30, 0xa6, 0x39, 0x4e, 0x59, 0xc0, 0xb2, 0x05, 0x96,
	0x85, 0xa5, 0x85, 0x35, 0xf9, 0x9f, 0x27, 0x69, 0xfc, 0x2f, 0xac, 0x7a, 0xf6, 0x15, 0x00, 0x00,
	0xff, 0xff, 0x3d, 0x03, 0xdc, 0x90, 0x7b, 0x02, 0x00, 0x00,
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
	Recover(ctx context.Context, in *payload.Replication_Recovery, opts ...grpc.CallOption) (*payload.Empty, error)
	Rebalance(ctx context.Context, in *payload.Replication_Rebalance, opts ...grpc.CallOption) (*payload.Empty, error)
	AgentInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Replication_Agents, error)
}

type replicationClient struct {
	cc *grpc.ClientConn
}

func NewReplicationClient(cc *grpc.ClientConn) ReplicationClient {
	return &replicationClient{cc}
}

func (c *replicationClient) Recover(ctx context.Context, in *payload.Replication_Recovery, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.replication.agent.v1.Replication/Recover", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicationClient) Rebalance(ctx context.Context, in *payload.Replication_Rebalance, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/manager.replication.agent.v1.Replication/Rebalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicationClient) AgentInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Replication_Agents, error) {
	out := new(payload.Replication_Agents)
	err := c.cc.Invoke(ctx, "/manager.replication.agent.v1.Replication/AgentInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplicationServer is the server API for Replication service.
type ReplicationServer interface {
	Recover(context.Context, *payload.Replication_Recovery) (*payload.Empty, error)
	Rebalance(context.Context, *payload.Replication_Rebalance) (*payload.Empty, error)
	AgentInfo(context.Context, *payload.Empty) (*payload.Replication_Agents, error)
}

// UnimplementedReplicationServer can be embedded to have forward compatible implementations.
type UnimplementedReplicationServer struct {
}

func (*UnimplementedReplicationServer) Recover(ctx context.Context, req *payload.Replication_Recovery) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Recover not implemented")
}
func (*UnimplementedReplicationServer) Rebalance(ctx context.Context, req *payload.Replication_Rebalance) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rebalance not implemented")
}
func (*UnimplementedReplicationServer) AgentInfo(ctx context.Context, req *payload.Empty) (*payload.Replication_Agents, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AgentInfo not implemented")
}

func RegisterReplicationServer(s *grpc.Server, srv ReplicationServer) {
	s.RegisterService(&_Replication_serviceDesc, srv)
}

func _Replication_Recover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Replication_Recovery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationServer).Recover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.replication.agent.v1.Replication/Recover",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServer).Recover(ctx, req.(*payload.Replication_Recovery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Replication_Rebalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Replication_Rebalance)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationServer).Rebalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.replication.agent.v1.Replication/Rebalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServer).Rebalance(ctx, req.(*payload.Replication_Rebalance))
	}
	return interceptor(ctx, in, info, handler)
}

func _Replication_AgentInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationServer).AgentInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.replication.agent.v1.Replication/AgentInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServer).AgentInfo(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Replication_serviceDesc = grpc.ServiceDesc{
	ServiceName: "manager.replication.agent.v1.Replication",
	HandlerType: (*ReplicationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Recover",
			Handler:    _Replication_Recover_Handler,
		},
		{
			MethodName: "Rebalance",
			Handler:    _Replication_Rebalance_Handler,
		},
		{
			MethodName: "AgentInfo",
			Handler:    _Replication_AgentInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/proto/v1/manager/replication/agent/replication_manager.proto",
}
