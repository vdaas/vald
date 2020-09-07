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
	// 329 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xc1, 0x4a, 0x33, 0x31,
	0x14, 0x85, 0x99, 0x2e, 0xfe, 0x9f, 0x8e, 0xe0, 0x22, 0x4a, 0xd5, 0x69, 0x1d, 0x61, 0x1e, 0x20,
	0xa1, 0xba, 0x12, 0xdc, 0xb4, 0xe0, 0xc2, 0x85, 0x20, 0x5d, 0x14, 0xd4, 0x85, 0xdc, 0xce, 0xa4,
	0x31, 0x90, 0xe6, 0x86, 0x4c, 0x1c, 0xe8, 0xd6, 0x57, 0x70, 0x27, 0xf8, 0x3e, 0x2e, 0x05, 0x5f,
	0x40, 0x8a, 0x0f, 0x22, 0x4d, 0x67, 0x74, 0x5a, 0x6b, 0x57, 0x81, 0x9c, 0x73, 0xbf, 0x43, 0x72,
	0x6e, 0xd8, 0x03, 0x23, 0x73, 0x66, 0x2c, 0x3a, 0x64, 0x45, 0x97, 0x4d, 0x40, 0x83, 0xe0, 0x96,
	0x59, 0x6e, 0x94, 0x4c, 0xc1, 0x49, 0xd4, 0x0c, 0x04, 0xd7, 0xae, 0x7e, 0x73, 0x57, 0xba, 0xa8,
	0x1f, 0x23, 0x3b, 0x6b, 0xa4, 0xe8, 0x54, 0x48, 0x77, 0xff, 0x30, 0xa2, 0x29, 0x4e, 0x58, 0x91,
	0x01, 0xe4, 0xac, 0x00, 0x95, 0xb1, 0xe5, 0x34, 0x03, 0x53, 0x85, 0x90, 0x55, 0xe7, 0x82, 0x17,
	0x75, 0x04, 0xa2, 0x50, 0x7c, 0xee, 0x65, 0xa0, 0x35, 0x3a, 0x4f, 0xce, 0x17, 0xea, 0xf1, 0x4b,
	0x23, 0xdc, 0x1a, 0xfc, 0x04, 0x92, 0x61, 0xf8, 0x7f, 0xc0, 0x53, 0x2c, 0xb8, 0x25, 0x87, 0xb4,
	0x02, 0xd5, 0x0c, 0xb4, 0x54, 0xa7, 0xd1, 0xf6, 0xb7, 0x7c, 0x3e, 0x31, 0x6e, 0x9a, 0x74, 0x1e,
	0xdf, 0x3f, 0x9f, 0x1a, 0xad, 0x64, 0x77, 0xe9, 0xb5, 0xb6, 0x84, 0xdd, 0x86, 0xcd, 0x01, 0x1f,
	0x81, 0x02, 0x9d, 0x72, 0x12, 0xff, 0x41, 0x2e, 0xf5, 0x5f, 0xe8, 0xd8, 0xa3, 0xf7, 0x93, 0xd6,
	0x0a, 0xba, 0xe2, 0x5d, 0x87, 0xcd, 0xde, 0xfc, 0x57, 0x2f, 0xf4, 0x18, 0xc9, 0xca, 0x70, 0xd4,
	0x5e, 0x1b, 0xe6, 0xfd, 0x79, 0x72, 0xe4, 0xc9, 0x07, 0x64, 0x6f, 0x4d, 0x45, 0x52, 0x8f, 0xb1,
	0xff, 0x1c, 0xbc, 0xce, 0xe2, 0xe0, 0x6d, 0x16, 0x07, 0x1f, 0xb3, 0x38, 0x08, 0x19, 0x5a, 0x41,
	0x7d, 0x01, 0x74, 0x5e, 0x00, 0x05, 0x23, 0x69, 0xd1, 0xa5, 0x55, 0x83, 0x35, 0x08, 0xf5, 0x90,
	0x7e, 0x7b, 0x08, 0x2a, 0xab, 0x05, 0x5f, 0x2e, 0x9c, 0x3e, 0xfe, 0x2a, 0xb8, 0x39, 0xdb, 0xd0,
	0xab, 0xb0, 0x26, 0xdd, 0xb8, 0x44, 0xa3, 0x7f, 0xbe, 0xc3, 0x93, 0xaf, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x76, 0x9d, 0xe9, 0x94, 0x76, 0x02, 0x00, 0x00,
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
	Metadata: "apis/proto/v1/manager/replication/agent/replication_manager.proto",
}
