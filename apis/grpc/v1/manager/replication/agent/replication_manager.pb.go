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
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0xd2, 0x41, 0x4b, 0xc3, 0x30,
	0x14, 0x07, 0x70, 0xba, 0x83, 0xb2, 0x0a, 0x1e, 0xa2, 0x4c, 0xed, 0x66, 0x85, 0x7a, 0x4f, 0x98,
	0x5e, 0xbd, 0x6c, 0xe0, 0xc1, 0x83, 0x20, 0x3b, 0x0c, 0xd4, 0x83, 0xbc, 0xb5, 0x59, 0x0c, 0x64,
	0x79, 0x21, 0x8d, 0x85, 0x5d, 0xfd, 0x0a, 0xde, 0x04, 0xbf, 0x8f, 0x47, 0xc1, 0x2f, 0x20, 0xc3,
	0x0f, 0x22, 0xcb, 0x5a, 0xed, 0xe6, 0xf4, 0x54, 0xe8, 0xff, 0xbd, 0xdf, 0x23, 0x79, 0x09, 0x7b,
	0x60, 0x64, 0xce, 0x8c, 0x45, 0x87, 0xac, 0xe8, 0xb2, 0x09, 0x68, 0x10, 0xdc, 0x32, 0xcb, 0x8d,
	0x92, 0x29, 0x38, 0x89, 0x9a, 0x81, 0xe0, 0xda, 0xd5, 0xff, 0xdc, 0x95, 0x55, 0xd4, 0xb7, 0x91,
	0x9d, 0x35, 0x51, 0x74, 0xbc, 0xec, 0x1a, 0x98, 0x2a, 0x84, 0xac, 0xfa, 0x2e, 0x3a, 0xa3, 0x8e,
	0x40, 0x14, 0x8a, 0x33, 0x30, 0x92, 0x81, 0xd6, 0xe8, 0xbc, 0x91, 0x2f, 0xd2, 0x93, 0x97, 0x46,
	0xb8, 0x35, 0xf8, 0xa1, 0xc9, 0x30, 0xdc, 0x1c, 0xf0, 0x14, 0x0b, 0x6e, 0xc9, 0x21, 0xad, 0xa0,
	0x5a, 0x01, 0x2d, 0xd3, 0x69, 0xb4, 0xfd, 0x1d, 0x9f, 0x4f, 0x8c, 0x9b, 0x26, 0x9d, 0xc7, 0xf7,
	0xcf, 0xa7, 0x46, 0x2b, 0xd9, 0x5d, 0x3a, 0x97, 0x2d, 0xb1, 0xdb, 0xb0, 0x39, 0xe0, 0x23, 0x50,
	0xa0, 0x53, 0x4e, 0xe2, 0x3f, 0xe4, 0x32, 0xff, 0x45, 0xc7, 0x9e, 0xde, 0x4f, 0x5a, 0x2b, 0x74,
	0xe5, 0x5d, 0x87, 0xcd, 0xde, 0xfc, 0xfe, 0x2e, 0xf4, 0x18, 0xc9, 0x4a, 0x73, 0xd4, 0x5e, 0x3b,
	0xcc, 0xd7, 0xe7, 0xc9, 0x91, 0x97, 0x0f, 0xc8, 0xde, 0x9a, 0x65, 0x48, 0x3d, 0xc6, 0xfe, 0x73,
	0xf0, 0x3a, 0x8b, 0x83, 0xb7, 0x59, 0x1c, 0x7c, 0xcc, 0xe2, 0x20, 0x64, 0x68, 0x05, 0x2d, 0x32,
	0x80, 0x9c, 0x16, 0xa0, 0x32, 0x0a, 0x46, 0xd2, 0xa2, 0x4b, 0xab, 0x5d, 0xd5, 0x10, 0xea, 0x91,
	0x7e, 0x7b, 0x08, 0x2a, 0xab, 0x0d, 0xbe, 0x5c, 0x54, 0xfa, 0xf1, 0x57, 0xc1, 0xcd, 0x99, 0x90,
	0xee, 0xfe, 0x61, 0x44, 0x53, 0x9c, 0x30, 0xcf, 0xb2, 0x39, 0xcb, 0xfc, 0x5e, 0x85, 0x35, 0xe9,
	0xbf, 0xcf, 0x65, 0xb4, 0xe1, 0x77, 0x78, 0xfa, 0x15, 0x00, 0x00, 0xff, 0xff, 0xb3, 0xca, 0x7f,
	0x0c, 0x60, 0x02, 0x00, 0x00,
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
