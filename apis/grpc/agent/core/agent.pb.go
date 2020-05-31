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

func init() { proto.RegisterFile("core/agent.proto", fileDescriptor_30864f15308ac822) }

var fileDescriptor_30864f15308ac822 = []byte{
	// 324 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xcd, 0x4a, 0x3b, 0x31,
	0x14, 0xc5, 0x99, 0xf2, 0xff, 0x0b, 0x4d, 0xad, 0x4a, 0xfc, 0x80, 0x96, 0xd2, 0x85, 0xee, 0x5c,
	0x24, 0xa0, 0x2f, 0x60, 0x5b, 0x5c, 0x14, 0x37, 0xa2, 0x20, 0xe8, 0xca, 0xdb, 0x49, 0x1a, 0x07,
	0xd2, 0xdc, 0x34, 0x93, 0x0e, 0x76, 0xeb, 0x2b, 0xf8, 0x52, 0x2e, 0x05, 0x77, 0xae, 0xa4, 0xf8,
	0x20, 0x92, 0x4c, 0xad, 0x1f, 0xb3, 0x73, 0x99, 0x7b, 0xcf, 0xf9, 0xdd, 0x03, 0x27, 0x64, 0x2b,
	0x45, 0x27, 0x39, 0x28, 0x69, 0x3c, 0xb3, 0x0e, 0x3d, 0xd2, 0x7f, 0x61, 0xd2, 0x6e, 0x5a, 0x98,
	0x6b, 0x04, 0x51, 0x0e, 0xdb, 0x1d, 0x85, 0xa8, 0xb4, 0xe4, 0x60, 0x33, 0x0e, 0xc6, 0xa0, 0x07,
	0x9f, 0xa1, 0xc9, 0x97, 0xdb, 0x75, 0x3b, 0xe2, 0x6a, 0xaa, 0xcb, 0xd7, 0xd1, 0x6b, 0x8d, 0xfc,
	0xef, 0x05, 0x20, 0xbd, 0x26, 0x8d, 0x81, 0x93, 0xe0, 0xe5, 0xd0, 0x08, 0x79, 0x4f, 0x0f, 0xd8,
	0x27, 0x74, 0x80, 0xc6, 0x3b, 0xd4, 0xec, 0xdb, 0xf6, 0x42, 0x4e, 0x67, 0x32, 0xf7, 0xed, 0x8d,
	0x95, 0xe8, 0x74, 0x62, 0xfd, 0x7c, 0x7f, 0xf7, 0xe1, 0xe5, 0xfd, 0xb1, 0xb6, 0x49, 0x9b, 0x3c,
	0x0b, 0x32, 0x9e, 0x46, 0x0b, 0x3d, 0x21, 0xf5, 0x4b, 0x28, 0x96, 0xe0, 0x5f, 0x9e, 0x0a, 0x63,
	0x3b, 0x32, 0x9a, 0xb4, 0xb1, 0x64, 0xe4, 0x50, 0x48, 0xaa, 0x08, 0x2d, 0xcf, 0xf7, 0x8c, 0xf8,
	0x42, 0xfd, 0x29, 0x63, 0x27, 0xf2, 0xf7, 0xe8, 0xce, 0x8f, 0x8c, 0x60, 0x44, 0x3c, 0x74, 0x46,
	0xea, 0xd1, 0x3d, 0x34, 0x63, 0xac, 0x44, 0x6d, 0xad, 0xde, 0x61, 0xcd, 0xa2, 0x90, 0x0d, 0x70,
	0x66, 0x7c, 0x25, 0x75, 0x66, 0xc6, 0xd8, 0xbf, 0x7d, 0x5a, 0x74, 0x93, 0xe7, 0x45, 0x37, 0x79,
	0x5b, 0x74, 0x13, 0xd2, 0x42, 0xa7, 0x58, 0x21, 0x00, 0x72, 0x56, 0x80, 0x16, 0xac, 0xec, 0x31,
	0x14, 0xd8, 0xaf, 0x5f, 0x81, 0x16, 0xb1, 0x86, 0xf3, 0xe4, 0xe6, 0x50, 0x65, 0xfe, 0x6e, 0x36,
	0x62, 0x29, 0x4e, 0x78, 0x94, 0xf3, 0x20, 0x0f, 0x6d, 0xe6, 0x5c, 0x39, 0x9b, 0x96, 0x1f, 0x80,
	0x07, 0xe3, 0x68, 0x2d, 0xb6, 0x78, 0xfc, 0x11, 0x00, 0x00, 0xff, 0xff, 0x24, 0x47, 0x3a, 0xa6,
	0x1a, 0x02, 0x00, 0x00,
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
	Metadata: "core/agent.proto",
}
