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
	// 331 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xcd, 0x4a, 0x03, 0x31,
	0x1c, 0xc4, 0xd9, 0xa2, 0x42, 0x53, 0xaa, 0x10, 0x3f, 0xc0, 0x52, 0x8a, 0x58, 0xaf, 0x26, 0x54,
	0x5f, 0xc0, 0xb6, 0x78, 0x28, 0x5e, 0x44, 0x41, 0xd0, 0xdb, 0xbf, 0x9b, 0x34, 0x46, 0xb6, 0xf9,
	0xc7, 0xdd, 0x74, 0xb1, 0x57, 0x5f, 0xc1, 0x97, 0xf2, 0x28, 0x78, 0xf3, 0x24, 0xc5, 0x07, 0x91,
	0x64, 0xdb, 0x6a, 0xdd, 0x9b, 0xa7, 0xfd, 0x98, 0x99, 0x5f, 0x06, 0x26, 0xe4, 0x08, 0xac, 0xce,
	0xb8, 0x4d, 0xd1, 0x21, 0xcf, 0x3b, 0x1c, 0x94, 0x34, 0x8e, 0xc7, 0x98, 0xca, 0xe2, 0x95, 0x05,
	0x85, 0xae, 0xf9, 0x3f, 0x8d, 0xf6, 0xaa, 0xd7, 0xc2, 0x34, 0x41, 0x10, 0x8b, 0x67, 0x61, 0x6d,
	0x34, 0x15, 0xa2, 0x4a, 0x24, 0x07, 0xab, 0x39, 0x18, 0x83, 0x0e, 0x9c, 0x46, 0x93, 0x15, 0xea,
	0xc9, 0x47, 0x85, 0xac, 0x77, 0x3d, 0x98, 0xde, 0x92, 0x5a, 0x3f, 0x95, 0xe0, 0xe4, 0xc0, 0x08,
	0xf9, 0x44, 0xdb, 0x6c, 0x81, 0xe9, 0xa3, 0x71, 0x29, 0x26, 0xec, 0x97, 0x7a, 0x25, 0x1f, 0x27,
	0x32, 0x73, 0x8d, 0xcd, 0xa5, 0xe9, 0x7c, 0x6c, 0xdd, 0xf4, 0x70, 0xf7, 0xf9, 0xfd, 0xeb, 0xa5,
	0xb2, 0x45, 0xeb, 0x5c, 0x7b, 0x1b, 0x8f, 0x43, 0x84, 0x9e, 0x91, 0xea, 0x35, 0xe4, 0x73, 0xf0,
	0x9f, 0x4c, 0x89, 0xb1, 0x1d, 0x18, 0x75, 0x5a, 0x9b, 0x33, 0x32, 0xc8, 0x25, 0x55, 0x84, 0x16,
	0xc7, 0x77, 0x8d, 0xf8, 0x41, 0xfd, 0xab, 0x63, 0x33, 0xf0, 0xf7, 0xe8, 0xce, 0x4a, 0x47, 0x30,
	0x22, 0x1c, 0x74, 0x41, 0xaa, 0x21, 0x3d, 0x30, 0x23, 0x2c, 0x55, 0xdd, 0x5f, 0x7e, 0x7b, 0x99,
	0x05, 0x23, 0xeb, 0xe3, 0xc4, 0xb8, 0x52, 0x6b, 0x6d, 0x46, 0xd8, 0x7b, 0x78, 0x9d, 0xb5, 0xa2,
	0xb7, 0x59, 0x2b, 0xfa, 0x9c, 0xb5, 0x22, 0x72, 0x80, 0xa9, 0x62, 0xb9, 0x00, 0xc8, 0x58, 0x0e,
	0x89, 0x60, 0x60, 0x35, 0xcb, 0x3b, 0xac, 0x98, 0xd5, 0xef, 0xd9, 0xab, 0xde, 0x40, 0x22, 0xc2,
	0x1a, 0x97, 0xd1, 0xdd, 0xb1, 0xd2, 0xee, 0x7e, 0x32, 0x64, 0x31, 0x8e, 0x79, 0x48, 0x71, 0x9f,
	0xe2, 0x61, 0x72, 0x95, 0xda, 0x78, 0xf5, 0x76, 0x0c, 0x37, 0xc2, 0x9e, 0xa7, 0xdf, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x95, 0xd4, 0x4c, 0xa3, 0x40, 0x02, 0x00, 0x00,
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
