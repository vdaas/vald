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

package controller

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
	proto.RegisterFile("apis/proto/v1/manager/replication/controller/replication_manager.proto", fileDescriptor_7996d9fdae0b086a)
}
func init() {
	golang_proto.RegisterFile("apis/proto/v1/manager/replication/controller/replication_manager.proto", fileDescriptor_7996d9fdae0b086a)
}

var fileDescriptor_7996d9fdae0b086a = []byte{
	// 311 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x51, 0xbd, 0x4e, 0xf3, 0x30,
	0x14, 0x95, 0xbf, 0xe1, 0x1b, 0xb2, 0xa0, 0x06, 0x31, 0xd0, 0xc1, 0xe2, 0x67, 0xc6, 0x56, 0x01,
	0x89, 0xb9, 0x45, 0x20, 0x31, 0x20, 0x21, 0x06, 0x06, 0x96, 0xea, 0x36, 0x49, 0x8d, 0x25, 0xd7,
	0xd7, 0x72, 0xdc, 0x88, 0xae, 0xbc, 0x02, 0x4f, 0xc2, 0x1b, 0x30, 0x76, 0x44, 0xe2, 0x05, 0x50,
	0xcb, 0x83, 0xa0, 0xd8, 0x0d, 0x71, 0x18, 0x60, 0xf2, 0xbd, 0xd7, 0xe7, 0xdc, 0x73, 0x7c, 0x9c,
	0x5c, 0x82, 0x91, 0x25, 0x37, 0x16, 0x1d, 0xf2, 0x6a, 0xc0, 0x67, 0xa0, 0x41, 0x14, 0x96, 0xdb,
	0xc2, 0x28, 0x99, 0x81, 0x93, 0xa8, 0x79, 0x86, 0xda, 0x59, 0x54, 0xaa, 0x3b, 0x1e, 0x6f, 0xa0,
	0xcc, 0x73, 0xd3, 0xfd, 0xa6, 0x8d, 0x20, 0xac, 0x65, 0xb2, 0x6a, 0xd0, 0x3f, 0xec, 0x4a, 0x19,
	0x58, 0x28, 0x84, 0xbc, 0x39, 0xc3, 0x9e, 0xfe, 0x91, 0x90, 0xee, 0x61, 0x3e, 0x61, 0x19, 0xce,
	0xb8, 0x40, 0x81, 0x01, 0x3f, 0x99, 0x4f, 0x7d, 0x17, 0xc8, 0x75, 0xb5, 0x81, 0x9f, 0xfd, 0x84,
	0x0b, 0x44, 0xa1, 0x0a, 0xaf, 0x14, 0x4a, 0x0e, 0x46, 0x72, 0xd0, 0x1a, 0x9d, 0xf7, 0x54, 0x06,
	0xe2, 0xf1, 0x63, 0xb2, 0x73, 0xdb, 0x3a, 0x3d, 0xff, 0x36, 0x9a, 0x8e, 0x93, 0xad, 0xe8, 0xe2,
	0x4a, 0x4f, 0x31, 0xed, 0xb1, 0xc6, 0x63, 0x35, 0x60, 0x17, 0x33, 0xe3, 0x16, 0x7d, 0x1a, 0x8f,
	0x22, 0x3c, 0x1b, 0x8a, 0x42, 0xbb, 0xf2, 0x60, 0xf7, 0xe9, 0xfd, 0xf3, 0xf9, 0xdf, 0x76, 0xda,
	0xeb, 0x24, 0x29, 0xf5, 0x14, 0x47, 0x2f, 0x64, 0xb9, 0xa2, 0xe4, 0x6d, 0x45, 0xc9, 0xc7, 0x8a,
	0x92, 0xd7, 0x35, 0x25, 0xcb, 0x35, 0x25, 0xc9, 0x29, 0x5a, 0xc1, 0xaa, 0x1c, 0xa0, 0x64, 0x15,
	0xa8, 0x9c, 0x81, 0x91, 0xf5, 0xee, 0xdf, 0x63, 0x1d, 0xed, 0xdd, 0x81, 0xca, 0x23, 0xfd, 0xeb,
	0x00, 0x6f, 0xdf, 0x73, 0x43, 0xee, 0x87, 0x51, 0x46, 0x5e, 0x80, 0xd7, 0x02, 0x3c, 0x64, 0x64,
	0x4d, 0xf6, 0xf7, 0xbf, 0x4f, 0xfe, 0xfb, 0xd0, 0x4e, 0xbe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x9b,
	0x4e, 0xd2, 0x83, 0x2e, 0x02, 0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/manager.replication.controller.v1.ReplicationController/ReplicationInfo", in, out, opts...)
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
		FullMethod: "/manager.replication.controller.v1.ReplicationController/ReplicationInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationControllerServer).ReplicationInfo(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _ReplicationController_serviceDesc = grpc.ServiceDesc{
	ServiceName: "manager.replication.controller.v1.ReplicationController",
	HandlerType: (*ReplicationControllerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReplicationInfo",
			Handler:    _ReplicationController_ReplicationInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/proto/v1/manager/replication/controller/replication_manager.proto",
}
