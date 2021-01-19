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

package core

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
	proto.RegisterFile("apis/proto/v1/agent/core/agent.proto", fileDescriptor_dc5722b42aaec2d2)
}
func init() {
	golang_proto.RegisterFile("apis/proto/v1/agent/core/agent.proto", fileDescriptor_dc5722b42aaec2d2)
}

var fileDescriptor_dc5722b42aaec2d2 = []byte{
	// 368 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0xc6, 0x49, 0x41, 0xa5, 0x29, 0x45, 0x5c, 0xff, 0x1c, 0x4a, 0x09, 0xa2, 0xe2, 0xcd, 0x5d,
	0xaa, 0x07, 0xcf, 0x6d, 0x51, 0xe8, 0xcd, 0x3f, 0xe0, 0xc1, 0x8b, 0x4c, 0x93, 0x6d, 0x0c, 0xa4,
	0x3b, 0x71, 0xb3, 0x59, 0xec, 0xd5, 0x57, 0xf0, 0x85, 0x3c, 0xf6, 0x28, 0xf8, 0x02, 0xd2, 0xfa,
	0x04, 0x3e, 0x81, 0xec, 0x6e, 0x8b, 0xa9, 0xed, 0xcd, 0x53, 0x26, 0x93, 0xdf, 0xf7, 0xcd, 0x17,
	0x66, 0xfc, 0x23, 0xc8, 0x92, 0x9c, 0x65, 0x12, 0x15, 0x32, 0xdd, 0x62, 0x10, 0x73, 0xa1, 0x58,
	0x88, 0x92, 0xbb, 0x92, 0xda, 0x2f, 0x64, 0xc3, 0x74, 0xa8, 0x6e, 0x35, 0x0e, 0x17, 0xf1, 0x0c,
	0x46, 0x29, 0x42, 0x34, 0x7f, 0x3a, 0xba, 0x71, 0x12, 0x27, 0xea, 0xb1, 0xe8, 0xd3, 0x10, 0x87,
	0x2c, 0xc6, 0x18, 0x1d, 0xdf, 0x2f, 0x06, 0xf6, 0xcd, 0x89, 0x4d, 0x35, 0xc3, 0xcf, 0xff, 0xe2,
	0x31, 0x62, 0x9c, 0x72, 0x3b, 0xc9, 0x95, 0x0c, 0xb2, 0x84, 0x81, 0x10, 0xa8, 0x40, 0x25, 0x28,
	0x72, 0x27, 0x3c, 0xfd, 0xae, 0xf8, 0x6b, 0x6d, 0x93, 0x92, 0x3c, 0xf8, 0xb5, 0xae, 0xe4, 0xa0,
	0x78, 0x4f, 0x44, 0xfc, 0x99, 0x1c, 0xd3, 0x79, 0x20, 0xdd, 0xa2, 0x5d, 0x14, 0x4a, 0x62, 0x4a,
	0x4b, 0xc0, 0x0d, 0x7f, 0x2a, 0x78, 0xae, 0x1a, 0x5b, 0x65, 0xee, 0x62, 0x98, 0xa9, 0xd1, 0xc1,
	0xee, 0xcb, 0xc7, 0xd7, 0x6b, 0x65, 0x93, 0xd4, 0x59, 0x62, 0x48, 0x16, 0x5a, 0x15, 0xb9, 0xf4,
	0xab, 0xb7, 0xa0, 0x67, 0xf6, 0xcb, 0xb2, 0x55, 0x4e, 0xdb, 0xd6, 0xa9, 0x4e, 0x6a, 0x33, 0xa7,
	0x1c, 0x34, 0x27, 0x43, 0x9f, 0xb8, 0x1c, 0x6d, 0x11, 0xfd, 0x1a, 0xfe, 0x23, 0x6f, 0xd3, 0x4e,
	0xd9, 0x23, 0x3b, 0x0b, 0x79, 0x41, 0x44, 0x76, 0xdc, 0xb5, 0x5f, 0xb5, 0x06, 0x3d, 0x31, 0xc0,
	0x55, 0xb1, 0x9b, 0xe5, 0x96, 0x81, 0xa8, 0xc5, 0x69, 0x17, 0x0b, 0xa1, 0x96, 0xfe, 0x20, 0x11,
	0x03, 0xec, 0xc8, 0xf1, 0x24, 0xf0, 0xde, 0x27, 0x81, 0xf7, 0x39, 0x09, 0xbc, 0xb7, 0x69, 0xe0,
	0x8d, 0xa7, 0x81, 0xe7, 0xef, 0xa3, 0x8c, 0xa9, 0x8e, 0x00, 0x72, 0xaa, 0x21, 0x8d, 0x28, 0x64,
	0x89, 0xb1, 0x74, 0x37, 0x64, 0x8e, 0xa7, 0x53, 0xbd, 0x83, 0x34, 0xb2, 0xdb, 0xba, 0xf2, 0xee,
	0xcb, 0x17, 0x62, 0x55, 0xcc, 0xa8, 0x98, 0x5b, 0xb9, 0xcc, 0xc2, 0xc5, 0x53, 0xec, 0xaf, 0xdb,
	0x7d, 0x9f, 0xfd, 0x04, 0x00, 0x00, 0xff, 0xff, 0x86, 0x3e, 0xed, 0x35, 0xad, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AgentClient is the client API for Agent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AgentClient interface {
	CreateIndex(ctx context.Context, in *payload.Control_CreateIndexRequest, opts ...grpc.CallOption) (*payload.Empty, error)
	SaveIndex(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Empty, error)
	CreateAndSaveIndex(ctx context.Context, in *payload.Control_CreateIndexRequest, opts ...grpc.CallOption) (*payload.Empty, error)
	IndexInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Index_Count, error)
}

type agentClient struct {
	cc *grpc.ClientConn
}

func NewAgentClient(cc *grpc.ClientConn) AgentClient {
	return &agentClient{cc}
}

func (c *agentClient) CreateIndex(ctx context.Context, in *payload.Control_CreateIndexRequest, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/core.v1.Agent/CreateIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) SaveIndex(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/core.v1.Agent/SaveIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) CreateAndSaveIndex(ctx context.Context, in *payload.Control_CreateIndexRequest, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/core.v1.Agent/CreateAndSaveIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) IndexInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Index_Count, error) {
	out := new(payload.Info_Index_Count)
	err := c.cc.Invoke(ctx, "/core.v1.Agent/IndexInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentServer is the server API for Agent service.
type AgentServer interface {
	CreateIndex(context.Context, *payload.Control_CreateIndexRequest) (*payload.Empty, error)
	SaveIndex(context.Context, *payload.Empty) (*payload.Empty, error)
	CreateAndSaveIndex(context.Context, *payload.Control_CreateIndexRequest) (*payload.Empty, error)
	IndexInfo(context.Context, *payload.Empty) (*payload.Info_Index_Count, error)
}

// UnimplementedAgentServer can be embedded to have forward compatible implementations.
type UnimplementedAgentServer struct {
}

func (*UnimplementedAgentServer) CreateIndex(ctx context.Context, req *payload.Control_CreateIndexRequest) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateIndex not implemented")
}
func (*UnimplementedAgentServer) SaveIndex(ctx context.Context, req *payload.Empty) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveIndex not implemented")
}
func (*UnimplementedAgentServer) CreateAndSaveIndex(ctx context.Context, req *payload.Control_CreateIndexRequest) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAndSaveIndex not implemented")
}
func (*UnimplementedAgentServer) IndexInfo(ctx context.Context, req *payload.Empty) (*payload.Info_Index_Count, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IndexInfo not implemented")
}

func RegisterAgentServer(s *grpc.Server, srv AgentServer) {
	s.RegisterService(&_Agent_serviceDesc, srv)
}

func _Agent_CreateIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Control_CreateIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).CreateIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.v1.Agent/CreateIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).CreateIndex(ctx, req.(*payload.Control_CreateIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_SaveIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).SaveIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.v1.Agent/SaveIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).SaveIndex(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_CreateAndSaveIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Control_CreateIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).CreateAndSaveIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.v1.Agent/CreateAndSaveIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).CreateAndSaveIndex(ctx, req.(*payload.Control_CreateIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_IndexInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).IndexInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.v1.Agent/IndexInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).IndexInfo(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Agent_serviceDesc = grpc.ServiceDesc{
	ServiceName: "core.v1.Agent",
	HandlerType: (*AgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateIndex",
			Handler:    _Agent_CreateIndex_Handler,
		},
		{
			MethodName: "SaveIndex",
			Handler:    _Agent_SaveIndex_Handler,
		},
		{
			MethodName: "CreateAndSaveIndex",
			Handler:    _Agent_CreateAndSaveIndex_Handler,
		},
		{
			MethodName: "IndexInfo",
			Handler:    _Agent_IndexInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/proto/v1/agent/core/agent.proto",
}
