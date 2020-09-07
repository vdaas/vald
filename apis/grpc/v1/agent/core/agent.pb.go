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
	proto.RegisterFile("apis/proto/v1/agent/core/agent.proto", fileDescriptor_dc5722b42aaec2d2)
}

var fileDescriptor_dc5722b42aaec2d2 = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xcf, 0x4a, 0x33, 0x31,
	0x14, 0xc5, 0x99, 0xf2, 0x7d, 0x42, 0x53, 0xaa, 0x10, 0xff, 0x80, 0xa5, 0x14, 0x51, 0xb7, 0x26,
	0x54, 0x57, 0xee, 0x6c, 0x8b, 0x8b, 0xe2, 0x46, 0x74, 0xa5, 0xbb, 0xdb, 0x99, 0x34, 0x0e, 0x4c,
	0x73, 0x63, 0x26, 0x0d, 0x76, 0xeb, 0x2b, 0xf8, 0x52, 0x2e, 0x05, 0x77, 0xae, 0xa4, 0xf8, 0x20,
	0x92, 0xa4, 0xad, 0xd6, 0x82, 0x0b, 0x57, 0x33, 0xc9, 0x39, 0xf7, 0x77, 0x0f, 0x9c, 0x90, 0x43,
	0xd0, 0x79, 0xc9, 0xb5, 0x41, 0x8b, 0xdc, 0xb5, 0x39, 0x48, 0xa1, 0x2c, 0x4f, 0xd1, 0x88, 0xf8,
	0xcb, 0x82, 0x42, 0xff, 0xf9, 0x9b, 0xc6, 0xa9, 0xcc, 0xed, 0xdd, 0x78, 0xc0, 0x52, 0x1c, 0x71,
	0x97, 0x01, 0x94, 0xdc, 0x41, 0x91, 0xf1, 0x65, 0x82, 0x86, 0x49, 0x81, 0x90, 0xcd, 0xbf, 0x11,
	0xd0, 0x68, 0x4a, 0x44, 0x59, 0x08, 0xef, 0xe5, 0xa0, 0x14, 0x5a, 0xb0, 0x39, 0xaa, 0x32, 0xaa,
	0xc7, 0x6f, 0x15, 0xf2, 0xbf, 0xe3, 0xd7, 0xd1, 0x1b, 0x52, 0xeb, 0x19, 0x01, 0x56, 0xf4, 0x55,
	0x26, 0x1e, 0xe8, 0x01, 0x9b, 0x63, 0x7a, 0xa8, 0xac, 0xc1, 0x82, 0x7d, 0x53, 0xaf, 0xc4, 0xfd,
	0x58, 0x94, 0xb6, 0xb1, 0xbe, 0x30, 0x9d, 0x8f, 0xb4, 0x9d, 0xec, 0x6f, 0x3f, 0xbe, 0x7e, 0x3c,
	0x55, 0x36, 0x68, 0x9d, 0xe7, 0xde, 0xc6, 0xd3, 0x30, 0x42, 0xcf, 0x48, 0xf5, 0x1a, 0xdc, 0x0c,
	0xfc, 0x63, 0x66, 0x85, 0xb1, 0x19, 0x18, 0x75, 0x5a, 0x9b, 0x31, 0x4a, 0x70, 0x82, 0x4a, 0x42,
	0xe3, 0xfa, 0x8e, 0xca, 0xbe, 0x50, 0x7f, 0xca, 0xd8, 0x0c, 0xfc, 0x1d, 0xba, 0xb5, 0x94, 0x11,
	0x54, 0x16, 0x16, 0x5d, 0x90, 0x6a, 0x98, 0xee, 0xab, 0x21, 0xae, 0x44, 0xdd, 0x5d, 0x9c, 0xbd,
	0xcc, 0x82, 0x91, 0xf5, 0x70, 0xac, 0xec, 0x4a, 0xea, 0x5c, 0x0d, 0xb1, 0x3b, 0x7c, 0x9e, 0xb6,
	0x92, 0x97, 0x69, 0x2b, 0x79, 0x9f, 0xb6, 0x12, 0xb2, 0x87, 0x46, 0xb2, 0x50, 0x1e, 0xf3, 0xe5,
	0x31, 0xd0, 0x39, 0x73, 0x6d, 0x16, 0xcb, 0xf6, 0x2d, 0x77, 0x63, 0x13, 0x97, 0xc9, 0xed, 0xd1,
	0x2f, 0x75, 0x4b, 0xa3, 0xd3, 0xe5, 0xf7, 0x32, 0x58, 0x0b, 0x5d, 0x9e, 0x7c, 0x06, 0x00, 0x00,
	0xff, 0xff, 0x1e, 0x85, 0x84, 0x80, 0x52, 0x02, 0x00, 0x00,
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
	Metadata: "apis/proto/v1/agent/core/agent.proto",
}
