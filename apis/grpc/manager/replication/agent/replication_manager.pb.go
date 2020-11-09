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
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xc1, 0x4a, 0xf3, 0x40,
	0x14, 0x85, 0x49, 0x17, 0xff, 0x4f, 0x47, 0x74, 0x31, 0x4a, 0xd5, 0xb4, 0x46, 0xc8, 0x4a, 0x10,
	0x66, 0x40, 0x57, 0x2e, 0x2d, 0xb8, 0x70, 0x21, 0x48, 0x17, 0x05, 0x75, 0x21, 0xb7, 0xc9, 0x74,
	0x0c, 0x4c, 0xe7, 0x0e, 0x93, 0xb1, 0xd0, 0xad, 0x4f, 0x20, 0xb8, 0xf6, 0x7d, 0x5c, 0x0a, 0xbe,
	0x80, 0x14, 0x1f, 0x44, 0x32, 0x49, 0x34, 0x34, 0x75, 0x3b, 0xe7, 0x9c, 0xef, 0x0c, 0xf7, 0x90,
	0x63, 0x2b, 0x8c, 0xca, 0x12, 0x70, 0x19, 0x6a, 0x0e, 0x52, 0x68, 0xc7, 0x1b, 0x2f, 0xf7, 0x33,
	0xd0, 0x20, 0x85, 0x65, 0xc6, 0xa2, 0x43, 0xba, 0xbd, 0x46, 0x0a, 0x37, 0x0d, 0x2c, 0x14, 0x42,
	0x5a, 0x7a, 0xc2, 0x81, 0x44, 0x94, 0x4a, 0x70, 0x30, 0x19, 0x07, 0xad, 0xd1, 0x79, 0x77, 0x5e,
	0xaa, 0x27, 0xaf, 0x1d, 0xb2, 0x31, 0xfa, 0x85, 0xd0, 0x31, 0xf9, 0x3f, 0x12, 0x09, 0xce, 0x85,
	0xa5, 0x07, 0xac, 0x06, 0x35, 0x0c, 0xac, 0x52, 0x17, 0xe1, 0xd6, 0x8f, 0x7c, 0x31, 0x33, 0x6e,
	0x11, 0x0f, 0x9e, 0x3e, 0xbe, 0x5e, 0x3a, 0xbd, 0x78, 0xa7, 0xf9, 0x5f, 0x6e, 0x2b, 0xd8, 0x1d,
	0xe9, 0x8e, 0xc4, 0x04, 0x14, 0xe8, 0x44, 0xd0, 0xe8, 0x0f, 0x72, 0xa5, 0xb7, 0xd0, 0x91, 0x47,
	0xef, 0xc5, 0xbd, 0x15, 0x74, 0xcd, 0xbb, 0x21, 0xdd, 0xf3, 0xe2, 0x52, 0x97, 0x7a, 0x8a, 0x74,
	0x25, 0x1c, 0xf6, 0xd7, 0x96, 0x79, 0x7f, 0x1e, 0x1f, 0x7a, 0xf2, 0x3e, 0xdd, 0xe5, 0xed, 0xb3,
	0x67, 0x7a, 0x8a, 0xc3, 0xe7, 0xe0, 0x6d, 0x19, 0x05, 0xef, 0xcb, 0x28, 0xf8, 0x5c, 0x46, 0x01,
	0x39, 0x42, 0x2b, 0xd9, 0x3c, 0x05, 0xc8, 0xd9, 0x1c, 0x54, 0xca, 0xea, 0x39, 0x1a, 0x69, 0xe6,
	0xd3, 0xc3, 0xfe, 0x18, 0x54, 0xda, 0x68, 0xbc, 0x2a, 0x9d, 0xbe, 0xf7, 0x3a, 0xb8, 0x3d, 0x93,
	0x99, 0x7b, 0x78, 0x9c, 0xb0, 0x04, 0x67, 0xdc, 0xf3, 0x78, 0xc1, 0x2b, 0x46, 0xca, 0xb9, 0xb4,
	0x26, 0xe1, 0x15, 0xb9, 0xfd, 0xaf, 0xc9, 0x3f, 0xbf, 0xdc, 0xe9, 0x77, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x95, 0x4a, 0x66, 0x91, 0x2a, 0x02, 0x00, 0x00,
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
