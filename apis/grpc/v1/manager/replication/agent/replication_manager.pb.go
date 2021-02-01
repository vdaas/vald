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

	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	golang_proto "github.com/golang/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	codes "github.com/vdaas/vald/internal/net/grpc/codes"
	status "github.com/vdaas/vald/internal/net/grpc/status"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
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
func init() {
	golang_proto.RegisterFile("apis/proto/v1/manager/replication/agent/replication_manager.proto", fileDescriptor_e8f74170057978aa)
}

var fileDescriptor_e8f74170057978aa = []byte{
	// 365 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xcd, 0x4a, 0xeb, 0x40,
	0x14, 0xc7, 0x49, 0x17, 0xf7, 0xd2, 0xdc, 0xd5, 0x0d, 0x97, 0x5e, 0x8d, 0x25, 0x6a, 0x5c, 0x3b,
	0x43, 0x75, 0xe1, 0xc6, 0x4d, 0x0b, 0x2e, 0x5c, 0x08, 0xd2, 0x85, 0x0b, 0x11, 0xca, 0x49, 0x32,
	0x1d, 0x07, 0xd2, 0x39, 0xc3, 0x24, 0x0d, 0x74, 0xeb, 0x2b, 0xf8, 0x06, 0xbe, 0x80, 0xaf, 0xe0,
	0xb2, 0x4b, 0xc1, 0x17, 0x90, 0xd6, 0x07, 0x91, 0xcc, 0x24, 0x98, 0xd6, 0xd6, 0x55, 0x66, 0xc2,
	0xef, 0xff, 0xc1, 0x9c, 0xe3, 0xf6, 0x41, 0x89, 0x8c, 0x2a, 0x8d, 0x39, 0xd2, 0xa2, 0x47, 0x27,
	0x20, 0x81, 0x33, 0x4d, 0x35, 0x53, 0xa9, 0x88, 0x21, 0x17, 0x28, 0x29, 0x70, 0x26, 0xf3, 0xe6,
	0x9f, 0x51, 0x45, 0x11, 0x23, 0xf3, 0xba, 0xf5, 0xb5, 0x81, 0x10, 0x23, 0x22, 0x45, 0xcf, 0x3f,
	0x5a, 0x0d, 0x50, 0x30, 0x4b, 0x11, 0x92, 0xfa, 0x6b, 0x2d, 0xfc, 0x63, 0x2e, 0xf2, 0xfb, 0x69,
	0x44, 0x62, 0x9c, 0x50, 0x8e, 0x1c, 0x2d, 0x1f, 0x4d, 0xc7, 0xe6, 0x66, 0xc5, 0xe5, 0xa9, 0xc2,
	0xcf, 0xd6, 0x71, 0x8e, 0xc8, 0x53, 0x66, 0x92, 0xec, 0x91, 0x82, 0x12, 0x14, 0xa4, 0xc4, 0xdc,
	0xd4, 0xc9, 0xac, 0xf0, 0xe4, 0xb9, 0xe5, 0xfe, 0x19, 0x7e, 0xb5, 0xf4, 0xee, 0xdc, 0xdf, 0x43,
	0x16, 0x63, 0xc1, 0xb4, 0x77, 0x40, 0xea, 0x4a, 0x45, 0x8f, 0x34, 0x18, 0x52, 0x01, 0x33, 0xff,
	0x6f, 0x93, 0xb8, 0x98, 0xa8, 0x7c, 0x16, 0x76, 0x1f, 0xde, 0x3e, 0x1e, 0x5b, 0x9d, 0xf0, 0xdf,
	0xca, 0x83, 0xe9, 0xca, 0x12, 0xdc, 0xf6, 0x90, 0x45, 0x90, 0x82, 0x8c, 0x99, 0x77, 0xb8, 0xdd,
	0xbf, 0x42, 0x36, 0x05, 0x04, 0x26, 0x60, 0x27, 0xec, 0xac, 0x05, 0xd4, 0xae, 0x23, 0xb7, 0xdd,
	0x2f, 0x5f, 0xfa, 0x52, 0x8e, 0xd1, 0xfb, 0xae, 0xf7, 0x83, 0x6d, 0xa9, 0x46, 0x95, 0x85, 0xfb,
	0xc6, 0x7f, 0xd7, 0xfb, 0xbf, 0x61, 0xe2, 0x42, 0x8e, 0x71, 0xf0, 0xe4, 0xcc, 0x17, 0x81, 0xf3,
	0xba, 0x08, 0x9c, 0xf7, 0x45, 0xe0, 0xbc, 0x2c, 0x03, 0x67, 0xbe, 0x0c, 0x1c, 0x97, 0xa2, 0xe6,
	0xa4, 0x48, 0x00, 0x32, 0x52, 0x40, 0x9a, 0x10, 0x50, 0xa2, 0x4c, 0xd8, 0xba, 0x09, 0x83, 0xbd,
	0x1b, 0x48, 0x93, 0x46, 0x81, 0x2b, 0x4b, 0x9a, 0x1a, 0xd7, 0xce, 0xed, 0x79, 0x63, 0xa2, 0xc6,
	0x96, 0x96, 0xb6, 0xd4, 0x4e, 0x54, 0xab, 0xf8, 0xc7, 0xdd, 0x8c, 0x7e, 0x99, 0xe9, 0x9e, 0x7e,
	0x06, 0x00, 0x00, 0xff, 0xff, 0x3b, 0x2d, 0x92, 0xad, 0xcd, 0x02, 0x00, 0x00,
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
