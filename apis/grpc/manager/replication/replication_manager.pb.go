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

package replication

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
	proto.RegisterFile("replication/replication_manager.proto", fileDescriptor_f84c0b45f233a237)
}

var fileDescriptor_f84c0b45f233a237 = []byte{
	// 247 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2d, 0x4a, 0x2d, 0xc8,
	0xc9, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x47, 0x62, 0xc7, 0xe7, 0x26, 0xe6, 0x25, 0xa6,
	0xa7, 0x16, 0xe9, 0x15, 0x14, 0xe5, 0x97, 0xe4, 0x0b, 0x09, 0x63, 0x91, 0x92, 0xe2, 0x2d, 0x48,
	0xac, 0xcc, 0xc9, 0x4f, 0x4c, 0x81, 0xa8, 0x91, 0x92, 0x49, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5,
	0x4f, 0x2c, 0xc8, 0xd4, 0x4f, 0xcc, 0xcb, 0xcb, 0x2f, 0x01, 0xab, 0x2e, 0x86, 0xca, 0xf2, 0x14,
	0x24, 0xe9, 0xa7, 0x17, 0xe6, 0x40, 0x78, 0x46, 0x39, 0x5c, 0xdc, 0x41, 0x08, 0x13, 0x85, 0xc2,
	0xb9, 0xf8, 0x91, 0xb8, 0x9e, 0x79, 0x69, 0xf9, 0x42, 0x7c, 0x7a, 0x30, 0xd3, 0x5d, 0x73, 0x0b,
	0x4a, 0x2a, 0xa5, 0x24, 0xe1, 0x7c, 0x90, 0xb4, 0x1e, 0x92, 0x72, 0x25, 0xc9, 0xa6, 0xcb, 0x4f,
	0x26, 0x33, 0x09, 0x0b, 0x09, 0x22, 0x7b, 0x40, 0x3f, 0x33, 0x2f, 0x2d, 0x5f, 0x8a, 0x65, 0xc3,
	0x03, 0x79, 0x26, 0xa7, 0xfa, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48,
	0x8e, 0x91, 0x4b, 0x29, 0xbf, 0x28, 0x5d, 0xaf, 0x2c, 0x25, 0x31, 0xb1, 0x58, 0xaf, 0x2c, 0x31,
	0x27, 0x45, 0x0f, 0xe6, 0x53, 0x24, 0xcd, 0x4e, 0x62, 0x61, 0x89, 0x39, 0x29, 0x48, 0x76, 0xf8,
	0x42, 0xd4, 0x04, 0x30, 0x46, 0x19, 0xa7, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7, 0xe7,
	0xea, 0x83, 0x0d, 0xd1, 0x07, 0x19, 0x02, 0xf2, 0x74, 0xb1, 0x7e, 0x7a, 0x51, 0x41, 0xb2, 0x3e,
	0xd4, 0x38, 0x64, 0xb7, 0x24, 0xb1, 0x81, 0x7d, 0x6d, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x8c,
	0x5a, 0x1f, 0xc0, 0x6e, 0x01, 0x00, 0x00,
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
	ReplicationInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Replication, error)
}

type replicationClient struct {
	cc *grpc.ClientConn
}

func NewReplicationClient(cc *grpc.ClientConn) ReplicationClient {
	return &replicationClient{cc}
}

func (c *replicationClient) ReplicationInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Replication, error) {
	out := new(payload.Info_Replication)
	err := c.cc.Invoke(ctx, "/replication_manager.Replication/ReplicationInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplicationServer is the server API for Replication service.
type ReplicationServer interface {
	ReplicationInfo(context.Context, *payload.Empty) (*payload.Info_Replication, error)
}

// UnimplementedReplicationServer can be embedded to have forward compatible implementations.
type UnimplementedReplicationServer struct {
}

func (*UnimplementedReplicationServer) ReplicationInfo(ctx context.Context, req *payload.Empty) (*payload.Info_Replication, error) {
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
	Metadata: "replication/replication_manager.proto",
}
