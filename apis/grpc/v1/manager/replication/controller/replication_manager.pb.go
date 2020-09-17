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
	proto "github.com/gogo/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
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
	proto.RegisterFile("apis/proto/v1/manager/replication/controller/replication_manager.proto", fileDescriptor_7996d9fdae0b086a)
}

var fileDescriptor_7996d9fdae0b086a = []byte{
	// 275 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x72, 0x4b, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0xcf, 0x4d, 0xcc, 0x4b, 0x4c, 0x4f,
	0x2d, 0xd2, 0x2f, 0x4a, 0x2d, 0xc8, 0xc9, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x4f, 0xce,
	0xcf, 0x2b, 0x29, 0xca, 0xcf, 0xc9, 0x41, 0x15, 0x8e, 0x87, 0x2a, 0xd5, 0x03, 0xeb, 0x15, 0x12,
	0xc6, 0x22, 0x25, 0xa5, 0x8c, 0x6a, 0x78, 0x41, 0x62, 0x65, 0x4e, 0x7e, 0x62, 0x0a, 0x8c, 0x86,
	0xe8, 0x94, 0x92, 0x49, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x4f, 0x2c, 0xc8, 0xd4, 0x4f, 0xcc,
	0xcb, 0xcb, 0x2f, 0x01, 0x9b, 0x51, 0x0c, 0x91, 0x35, 0x2a, 0xe2, 0x12, 0x0d, 0x42, 0x98, 0xec,
	0x0c, 0x77, 0x8a, 0x50, 0x24, 0x17, 0x3f, 0x92, 0x84, 0x67, 0x5e, 0x5a, 0xbe, 0x10, 0x9f, 0x1e,
	0xcc, 0x64, 0xd7, 0xdc, 0x82, 0x92, 0x4a, 0x29, 0x69, 0x38, 0x1f, 0x49, 0xa5, 0x9e, 0x63, 0x7a,
	0x6a, 0x5e, 0x49, 0xb1, 0x92, 0x64, 0xd3, 0xe5, 0x27, 0x93, 0x99, 0x84, 0x85, 0x04, 0x51, 0xfc,
	0x9a, 0x99, 0x97, 0x96, 0xef, 0xb4, 0x92, 0xf1, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18,
	0x1f, 0x3c, 0x92, 0x63, 0xe4, 0x32, 0xc9, 0x2f, 0x4a, 0xd7, 0x2b, 0x4b, 0x49, 0x4c, 0x2c, 0xd6,
	0x2b, 0x4b, 0xcc, 0x49, 0xd1, 0x4b, 0x2c, 0xc8, 0xd4, 0x2b, 0x33, 0xd4, 0x83, 0xf9, 0x1f, 0x49,
	0xbb, 0x1e, 0x22, 0xa8, 0x9c, 0x14, 0xc2, 0x12, 0x73, 0x52, 0x90, 0xec, 0xf5, 0x85, 0x28, 0x47,
	0xf8, 0x20, 0x80, 0x31, 0xca, 0x31, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57,
	0x1f, 0x6c, 0x81, 0x3e, 0xc8, 0x02, 0x7d, 0x70, 0xa8, 0xa5, 0x17, 0x15, 0x24, 0x13, 0x8e, 0x91,
	0x24, 0x36, 0x70, 0x30, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x3e, 0xcb, 0xef, 0x05, 0xc8,
	0x01, 0x00, 0x00,
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
	Metadata: "apis/proto/v1/manager/replication/controller/replication_manager.proto",
}
