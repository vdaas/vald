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
	// 339 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0xcd, 0x4a, 0x3b, 0x31,
	0x14, 0xc5, 0x99, 0xc2, 0xff, 0x2f, 0x4d, 0x29, 0x62, 0xfc, 0x58, 0x94, 0x52, 0x44, 0xc5, 0x9d,
	0x09, 0xd5, 0x27, 0x68, 0x8b, 0x42, 0x77, 0x7e, 0x80, 0x0b, 0x37, 0x72, 0x3b, 0x49, 0xc7, 0xc8,
	0x34, 0x37, 0x66, 0xd2, 0x60, 0xb7, 0xbe, 0x82, 0x2f, 0xe5, 0x52, 0xf0, 0x05, 0xa4, 0xf8, 0x04,
	0x3e, 0x81, 0x24, 0xd3, 0x62, 0x4b, 0xbb, 0x73, 0x35, 0x43, 0x72, 0xce, 0xef, 0x9e, 0x70, 0x2e,
	0x39, 0x02, 0xa3, 0x0a, 0x6e, 0x2c, 0x3a, 0xe4, 0xbe, 0xcd, 0x21, 0x93, 0xda, 0xf1, 0x14, 0xad,
	0x2c, 0x7f, 0x59, 0xbc, 0xa1, 0x1b, 0xe1, 0x84, 0xf9, 0x76, 0xe3, 0x70, 0x59, 0x6e, 0x60, 0x92,
	0x23, 0x88, 0xf9, 0xb7, 0x54, 0x37, 0x9a, 0x19, 0x62, 0x96, 0x4b, 0x0e, 0x46, 0x71, 0xd0, 0x1a,
	0x1d, 0x38, 0x85, 0xba, 0x28, 0x6f, 0x4f, 0xbf, 0x2b, 0xe4, 0x5f, 0x27, 0xb0, 0xe9, 0x3d, 0xa9,
	0xf5, 0xac, 0x04, 0x27, 0xfb, 0x5a, 0xc8, 0x67, 0x7a, 0xcc, 0xe6, 0x18, 0xdf, 0x66, 0x3d, 0xd4,
	0xce, 0x62, 0xce, 0x16, 0x04, 0xd7, 0xf2, 0x69, 0x2c, 0x0b, 0xd7, 0xd8, 0x5a, 0xd4, 0x9d, 0x8f,
	0x8c, 0x9b, 0x1c, 0xec, 0xbe, 0x7c, 0x7c, 0xbd, 0x56, 0x36, 0x69, 0x9d, 0xab, 0xa0, 0xe4, 0x69,
	0x74, 0xd1, 0x0b, 0x52, 0xbd, 0x01, 0x3f, 0xc3, 0xaf, 0xda, 0xd6, 0x91, 0xb6, 0x23, 0xa9, 0x4e,
	0x6b, 0x33, 0x52, 0x01, 0x5e, 0xd2, 0x11, 0xa1, 0x65, 0x8e, 0x8e, 0x16, 0xbf, 0xc0, 0x3f, 0xe4,
	0x6d, 0xc6, 0x29, 0x7b, 0x74, 0x67, 0x29, 0x2f, 0x68, 0x11, 0xc7, 0x5d, 0x91, 0x6a, 0x04, 0xf4,
	0xf5, 0x10, 0xd7, 0xc5, 0x6e, 0x2e, 0x1e, 0x05, 0x11, 0x8b, 0x72, 0xd6, 0xc3, 0xb1, 0x76, 0x2b,
	0x2f, 0x50, 0x7a, 0x88, 0xdd, 0xc7, 0xb7, 0x69, 0x2b, 0x79, 0x9f, 0xb6, 0x92, 0xcf, 0x69, 0x2b,
	0x21, 0xfb, 0x68, 0x33, 0xe6, 0x05, 0x40, 0xc1, 0x3c, 0xe4, 0x82, 0x81, 0x51, 0x01, 0x55, 0x36,
	0x1e, 0xaa, 0xee, 0x56, 0x6f, 0x21, 0x17, 0xb1, 0xa5, 0xcb, 0xe4, 0xee, 0x24, 0x53, 0xee, 0x61,
	0x3c, 0x60, 0x29, 0x8e, 0x78, 0x74, 0xf1, 0xe0, 0xe2, 0x71, 0x15, 0x32, 0x6b, 0xd2, 0xe5, 0xc5,
	0x19, 0xfc, 0x8f, 0x3d, 0x9f, 0xfd, 0x04, 0x00, 0x00, 0xff, 0xff, 0x41, 0xe0, 0x80, 0x4b, 0x5b,
	0x02, 0x00, 0x00,
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
