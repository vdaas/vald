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

package agent

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/payload"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

func init() {
	proto.RegisterFile("apis/proto/manager/replication/agent/replication_manager.proto", fileDescriptor_86c5f765f4865443)
}

var fileDescriptor_86c5f765f4865443 = []byte{
	// 319 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x31, 0x4b, 0x03, 0x31,
	0x1c, 0xc5, 0xb9, 0x0e, 0x4a, 0x23, 0x38, 0x44, 0xa9, 0x7a, 0xad, 0xa7, 0xdc, 0xe4, 0x94, 0x80,
	0x4e, 0x2e, 0x82, 0x05, 0x07, 0x07, 0x41, 0x3a, 0x14, 0xd4, 0x41, 0xfe, 0xbd, 0x4b, 0x63, 0x20,
	0xcd, 0x3f, 0xe4, 0x62, 0xa1, 0xab, 0x9f, 0x40, 0x70, 0xf6, 0xfb, 0x38, 0x0a, 0x7e, 0x01, 0x29,
	0x7e, 0x10, 0x69, 0x7a, 0xa7, 0x67, 0x6d, 0xa7, 0x83, 0x7b, 0xef, 0xfd, 0x1e, 0xc9, 0x0b, 0x39,
	0x03, 0xab, 0x0a, 0x6e, 0x1d, 0x7a, 0xe4, 0x23, 0x30, 0x20, 0x85, 0xe3, 0x4e, 0x58, 0xad, 0x32,
	0xf0, 0x0a, 0x0d, 0x07, 0x29, 0x8c, 0xaf, 0xff, 0xb9, 0x2f, 0x5d, 0x2c, 0x64, 0xe8, 0xd6, 0x12,
	0x29, 0x3e, 0xac, 0x41, 0x2d, 0x4c, 0x34, 0x42, 0x5e, 0x7d, 0xe7, 0xb1, 0xb8, 0x23, 0x11, 0xa5,
	0x16, 0x1c, 0xac, 0xe2, 0x60, 0x0c, 0xfa, 0x00, 0x28, 0xe6, 0xea, 0xf1, 0x6b, 0x83, 0x6c, 0xf4,
	0x7e, 0xb9, 0xb4, 0x4f, 0xd6, 0x7b, 0x22, 0xc3, 0xb1, 0x70, 0x74, 0x9f, 0x55, 0xa0, 0x9a, 0x81,
	0x95, 0xea, 0x24, 0xde, 0xfc, 0x91, 0x2f, 0x46, 0xd6, 0x4f, 0xd2, 0xce, 0xd3, 0xc7, 0xd7, 0x4b,
	0xa3, 0x95, 0x6e, 0xff, 0x39, 0x94, 0x2b, 0x61, 0x77, 0xa4, 0xd9, 0x13, 0x03, 0xd0, 0x60, 0x32,
	0x41, 0x93, 0x15, 0xe4, 0x52, 0xff, 0x87, 0x4e, 0x02, 0x7a, 0x37, 0x6d, 0x2d, 0xa0, 0x2b, 0xde,
	0x0d, 0x69, 0x9e, 0xcf, 0x2e, 0xef, 0xd2, 0x0c, 0x91, 0x2e, 0x84, 0xe3, 0xf6, 0xd2, 0xb2, 0xe0,
	0x2f, 0xd2, 0x83, 0x40, 0xde, 0xa3, 0x3b, 0x4b, 0x96, 0x50, 0x66, 0x88, 0xdd, 0xe7, 0xe8, 0x6d,
	0x9a, 0x44, 0xef, 0xd3, 0x24, 0xfa, 0x9c, 0x26, 0x11, 0x39, 0x42, 0x27, 0xd9, 0x38, 0x07, 0x28,
	0xd8, 0x18, 0x74, 0xce, 0xaa, 0x85, 0x6a, 0x69, 0x16, 0xd2, 0xdd, 0x76, 0x1f, 0x74, 0x5e, 0x6b,
	0xbc, 0x9a, 0x3b, 0x43, 0xef, 0x75, 0x74, 0x7b, 0x2a, 0x95, 0x7f, 0x78, 0x1c, 0xb0, 0x0c, 0x47,
	0x3c, 0xf0, 0xf8, 0x8c, 0xc7, 0xc3, 0x9a, 0xd2, 0xd9, 0x6c, 0xf5, 0x0b, 0x19, 0xac, 0x85, 0xe5,
	0x4e, 0xbe, 0x03, 0x00, 0x00, 0xff, 0xff, 0xb3, 0xd9, 0x67, 0x7e, 0x50, 0x02, 0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/replication_manager.Replication/Recover", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicationClient) Rebalance(ctx context.Context, in *payload.Replication_Rebalance, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/replication_manager.Replication/Rebalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicationClient) AgentInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Replication_Agents, error) {
	out := new(payload.Replication_Agents)
	err := c.cc.Invoke(ctx, "/replication_manager.Replication/AgentInfo", in, out, opts...)
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
		FullMethod: "/replication_manager.Replication/Recover",
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
		FullMethod: "/replication_manager.Replication/Rebalance",
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
		FullMethod: "/replication_manager.Replication/AgentInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServer).AgentInfo(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Replication_serviceDesc = grpc.ServiceDesc{
	ServiceName: "replication_manager.Replication",
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
	Metadata: "apis/proto/manager/replication/agent/replication_manager.proto",
}
