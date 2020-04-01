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
	proto.RegisterFile("replication/agent/replication_manager.proto", fileDescriptor_2c09480608bbf428)
}

var fileDescriptor_2c09480608bbf428 = []byte{
	// 329 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0xd1, 0x4f, 0x4a, 0xfb, 0x40,
	0x14, 0x07, 0x70, 0x52, 0x7e, 0xfc, 0xa4, 0xe3, 0x9f, 0xc5, 0x28, 0x55, 0xd3, 0x9a, 0x42, 0x56,
	0x82, 0x30, 0x03, 0xba, 0x72, 0x69, 0xc1, 0x85, 0x0b, 0x41, 0xba, 0x28, 0xa8, 0x0b, 0x79, 0x49,
	0xa6, 0x63, 0x60, 0x3a, 0x6f, 0x9c, 0x8c, 0x85, 0x6e, 0x3d, 0x81, 0xe0, 0x11, 0xbc, 0x80, 0xc7,
	0x70, 0x29, 0x78, 0x81, 0x52, 0x3c, 0x88, 0x74, 0x9a, 0x6a, 0x68, 0xeb, 0x32, 0xf9, 0x7e, 0xf3,
	0x79, 0x8f, 0x3c, 0x72, 0x64, 0x85, 0x51, 0x79, 0x0a, 0x2e, 0x47, 0xcd, 0x41, 0x0a, 0xed, 0x78,
	0xe5, 0xcd, 0xdd, 0x00, 0x34, 0x48, 0x61, 0x99, 0xb1, 0xe8, 0x90, 0x6e, 0xaf, 0x88, 0xc2, 0x4d,
	0x03, 0x23, 0x85, 0x90, 0xcd, 0x3a, 0x61, 0x4b, 0x22, 0x4a, 0x25, 0x38, 0x98, 0x9c, 0x83, 0xd6,
	0xe8, 0x7c, 0xbb, 0x28, 0xd3, 0x0d, 0x93, 0x70, 0xf9, 0xa0, 0x66, 0x4f, 0xc7, 0xaf, 0x35, 0xb2,
	0xde, 0xfd, 0x25, 0x69, 0x8f, 0xac, 0x75, 0x45, 0x8a, 0x43, 0x61, 0xe9, 0x01, 0x9b, 0xb3, 0x95,
	0x02, 0x2b, 0xd3, 0x51, 0xb8, 0xf5, 0x13, 0x9f, 0x0f, 0x8c, 0x1b, 0xc5, 0xad, 0xa7, 0xcf, 0xaf,
	0x97, 0x5a, 0x23, 0xde, 0xa9, 0x6e, 0xcf, 0x6d, 0x89, 0xdd, 0x92, 0x7a, 0x57, 0x24, 0xa0, 0x40,
	0xa7, 0x82, 0x46, 0x7f, 0xc8, 0x65, 0xbe, 0x44, 0x47, 0x9e, 0xde, 0x8b, 0x1b, 0x0b, 0xf4, 0xdc,
	0xbb, 0x26, 0xf5, 0xb3, 0xe9, 0x7f, 0xbb, 0xd0, 0x7d, 0xa4, 0x0b, 0x1f, 0x87, 0xcd, 0x95, 0xc3,
	0x7c, 0xbf, 0x88, 0xdb, 0x5e, 0xde, 0xa7, 0xbb, 0x7c, 0xf9, 0x08, 0xb9, 0xee, 0x63, 0xf8, 0xef,
	0x6d, 0xdc, 0xae, 0x75, 0x9e, 0x83, 0xf7, 0x49, 0x14, 0x7c, 0x4c, 0xa2, 0x60, 0x3c, 0x89, 0x02,
	0x72, 0x88, 0x56, 0xb2, 0x61, 0x06, 0x50, 0xb0, 0x21, 0xa8, 0x8c, 0xcd, 0x4f, 0x54, 0x31, 0x98,
	0x37, 0x3a, 0xcd, 0x1e, 0xa8, 0xac, 0x32, 0xf7, 0x72, 0xd6, 0xf4, 0xd3, 0xaf, 0x82, 0x9b, 0x53,
	0x99, 0xbb, 0xfb, 0xc7, 0x84, 0xa5, 0x38, 0xe0, 0xde, 0xe3, 0x53, 0x6f, 0x7a, 0xb8, 0x82, 0x4b,
	0x6b, 0x52, 0x5e, 0xca, 0xcb, 0xdb, 0x25, 0xff, 0xfd, 0xfd, 0x4e, 0xbe, 0x03, 0x00, 0x00, 0xff,
	0xff, 0x8c, 0x5e, 0xa1, 0xb7, 0x3e, 0x02, 0x00, 0x00,
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
	Metadata: "replication/agent/replication_manager.proto",
}
