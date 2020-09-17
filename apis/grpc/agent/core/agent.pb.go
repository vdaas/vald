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

package core

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

func init() { proto.RegisterFile("apis/proto/agent/core/agent.proto", fileDescriptor_000d4e288381ff88) }

var fileDescriptor_000d4e288381ff88 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xcf, 0x4a, 0x03, 0x31,
	0x10, 0xc6, 0xd9, 0xa2, 0x42, 0x53, 0xaa, 0x10, 0xff, 0x40, 0x4b, 0x29, 0xfe, 0xb9, 0x79, 0x48,
	0x40, 0x5f, 0xc0, 0xb6, 0x78, 0x28, 0x5e, 0x44, 0x41, 0xd0, 0x93, 0xd3, 0x4d, 0x1a, 0x03, 0xdb,
	0x4c, 0xcc, 0xa6, 0x8b, 0xbd, 0xfa, 0x0a, 0xbe, 0x94, 0x47, 0xc1, 0x9b, 0x27, 0x29, 0x3e, 0x88,
	0x24, 0xdb, 0xd6, 0x6a, 0x6f, 0x9e, 0x92, 0xcc, 0x7c, 0xf3, 0x9b, 0x0f, 0xbe, 0x90, 0x03, 0xb0,
	0x3a, 0xe7, 0xd6, 0xa1, 0x47, 0x0e, 0x4a, 0x1a, 0xcf, 0x53, 0x74, 0xb2, 0xbc, 0xb2, 0x58, 0xa6,
	0x6b, 0xa1, 0xd2, 0xdc, 0x5f, 0x12, 0x5a, 0x98, 0x64, 0x08, 0x62, 0x7e, 0x96, 0xba, 0x66, 0x4b,
	0x21, 0xaa, 0x4c, 0x72, 0xb0, 0x9a, 0x83, 0x31, 0xe8, 0xc1, 0x6b, 0x34, 0x79, 0xd9, 0x3d, 0xf9,
	0xa8, 0x90, 0xf5, 0x4e, 0xa0, 0xd2, 0x5b, 0x52, 0xeb, 0x39, 0x09, 0x5e, 0xf6, 0x8d, 0x90, 0x4f,
	0xf4, 0x88, 0xcd, 0x31, 0x3d, 0x34, 0xde, 0x61, 0xc6, 0x96, 0xba, 0x57, 0xf2, 0x71, 0x2c, 0x73,
	0xdf, 0xdc, 0x5c, 0x88, 0xce, 0x47, 0xd6, 0x4f, 0x0e, 0x77, 0x9f, 0xdf, 0xbf, 0x5e, 0x2a, 0x5b,
	0xb4, 0xce, 0x75, 0x90, 0xf1, 0x34, 0x8e, 0xd0, 0x33, 0x52, 0xbd, 0x86, 0x62, 0x06, 0xfe, 0x33,
	0xb3, 0xc2, 0xd8, 0x8e, 0x8c, 0x3a, 0xad, 0xcd, 0x18, 0x39, 0x14, 0x92, 0x2a, 0x42, 0xcb, 0xf5,
	0x1d, 0x23, 0x7e, 0x50, 0xff, 0xf2, 0xd8, 0x8a, 0xfc, 0x3d, 0xba, 0xf3, 0xcb, 0x23, 0x18, 0x11,
	0x17, 0x5d, 0x90, 0x6a, 0x9c, 0xee, 0x9b, 0x21, 0xae, 0x58, 0x6d, 0x2c, 0xde, 0xa1, 0xcd, 0xa2,
	0x90, 0xf5, 0x70, 0x6c, 0xfc, 0x8a, 0x6b, 0x6d, 0x86, 0xd8, 0xbd, 0x7f, 0x9d, 0xb6, 0x93, 0xb7,
	0x69, 0x3b, 0xf9, 0x9c, 0xb6, 0x13, 0xd2, 0x40, 0xa7, 0x58, 0x21, 0x00, 0x72, 0x56, 0x40, 0x26,
	0x58, 0x19, 0x66, 0x48, 0xb1, 0x5b, 0xbd, 0x81, 0x4c, 0xc4, 0x18, 0x2e, 0x93, 0xbb, 0x63, 0xa5,
	0xfd, 0xc3, 0x78, 0xc0, 0x52, 0x1c, 0xf1, 0x28, 0xe7, 0x41, 0xce, 0x63, 0xd0, 0xca, 0xd9, 0x74,
	0xe9, 0x43, 0x0c, 0x36, 0x62, 0x8a, 0xa7, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x8d, 0x19, 0x58,
	0x1d, 0x30, 0x02, 0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/core.Agent/CreateIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) SaveIndex(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/core.Agent/SaveIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) CreateAndSaveIndex(ctx context.Context, in *payload.Control_CreateIndexRequest, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/core.Agent/CreateAndSaveIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) IndexInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Index_Count, error) {
	out := new(payload.Info_Index_Count)
	err := c.cc.Invoke(ctx, "/core.Agent/IndexInfo", in, out, opts...)
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
		FullMethod: "/core.Agent/CreateIndex",
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
		FullMethod: "/core.Agent/SaveIndex",
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
		FullMethod: "/core.Agent/CreateAndSaveIndex",
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
		FullMethod: "/core.Agent/IndexInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).IndexInfo(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Agent_serviceDesc = grpc.ServiceDesc{
	ServiceName: "core.Agent",
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
	Metadata: "apis/proto/agent/core/agent.proto",
}
