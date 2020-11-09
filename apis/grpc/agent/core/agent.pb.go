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
	// 579 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x95, 0xc1, 0x6a, 0x13, 0x41,
	0x18, 0xc7, 0x5d, 0xd1, 0x48, 0x26, 0x49, 0x5b, 0xa6, 0xb5, 0xb5, 0xa1, 0xe4, 0x10, 0x2f, 0x92,
	0xc3, 0x8c, 0xa8, 0x87, 0x22, 0x1e, 0x6c, 0x92, 0x5a, 0x16, 0x09, 0x95, 0x06, 0x8b, 0x7a, 0x72,
	0xb2, 0xfb, 0x75, 0xbb, 0xb2, 0x99, 0x59, 0x67, 0x26, 0xa1, 0x41, 0xbc, 0xf8, 0x0a, 0x3e, 0x85,
	0x6f, 0xe2, 0x51, 0xf0, 0x05, 0x24, 0xf8, 0x20, 0x32, 0x33, 0x9b, 0x34, 0xc9, 0x46, 0x4a, 0x72,
	0xfd, 0xe6, 0xfb, 0xff, 0xf6, 0x37, 0x03, 0xdf, 0x7e, 0x68, 0x2b, 0x10, 0x12, 0x28, 0x8b, 0x80,
	0x6b, 0x92, 0x4a, 0xa1, 0x05, 0xbe, 0x63, 0x2a, 0xd5, 0x4a, 0xca, 0x46, 0x89, 0x60, 0xa1, 0x2b,
	0x56, 0x0f, 0x22, 0x21, 0xa2, 0x04, 0x28, 0x4b, 0x63, 0xca, 0x38, 0x17, 0x9a, 0xe9, 0x58, 0x70,
	0xe5, 0x4e, 0x9f, 0xfc, 0x28, 0xa1, 0xbb, 0x47, 0x06, 0x81, 0x5f, 0xa1, 0xc2, 0xf1, 0x55, 0xac,
	0xb4, 0xc2, 0x98, 0x4c, 0x08, 0xa7, 0xbd, 0x4f, 0x10, 0x68, 0xe2, 0xb7, 0xab, 0x4b, 0x6a, 0xf5,
	0x9d, 0x6f, 0xbf, 0xff, 0x7e, 0xbf, 0xbd, 0x81, 0xcb, 0x14, 0x6c, 0x90, 0x7e, 0x89, 0xc3, 0xaf,
	0xf8, 0x14, 0x15, 0xba, 0xc0, 0x64, 0x70, 0x89, 0xf7, 0xa6, 0x19, 0x57, 0x20, 0x67, 0xf0, 0x79,
	0x00, 0x4a, 0x57, 0x1f, 0xe4, 0x0f, 0x54, 0x2a, 0xb8, 0x82, 0x3a, 0xb6, 0xc8, 0x72, 0xfd, 0x1e,
	0x55, 0xf6, 0xe4, 0xb9, 0xd7, 0xc0, 0xef, 0x10, 0x72, 0x6d, 0xcd, 0x91, 0xdf, 0xc6, 0xfb, 0x8b,
	0x59, 0xbf, 0x7d, 0x33, 0xf6, 0xbe, 0xc5, 0x6e, 0xd6, 0x51, 0x86, 0xa5, 0x71, 0x68, 0xc8, 0x27,
	0xa8, 0xdc, 0xd5, 0x12, 0x58, 0x7f, 0x7d, 0xe1, 0x5b, 0x8f, 0xbc, 0xc7, 0x1e, 0xee, 0xa0, 0xad,
	0x59, 0xd0, 0xfa, 0xa2, 0x0e, 0xd7, 0x46, 0x05, 0x9f, 0x2b, 0x90, 0x1a, 0xef, 0x2e, 0x3e, 0xfb,
	0x39, 0x04, 0x5a, 0xc8, 0xea, 0xc6, 0xb4, 0x7e, 0xdc, 0x4f, 0xf5, 0x68, 0xe6, 0xdd, 0x62, 0x1b,
	0x34, 0xb7, 0x7b, 0x31, 0xb9, 0xdd, 0x8a, 0x2c, 0xe7, 0x70, 0x88, 0x4a, 0x9d, 0x41, 0xa2, 0xe3,
	0x2c, 0xbc, 0xb7, 0x3c, 0xac, 0xf2, 0x69, 0x63, 0xff, 0x36, 0x0d, 0x99, 0x86, 0x35, 0xec, 0x07,
	0x36, 0x38, 0x67, 0xbf, 0x22, 0x6b, 0xde, 0x3e, 0x0b, 0xaf, 0x60, 0xdf, 0x44, 0x85, 0x33, 0xe8,
	0x8b, 0x21, 0x2c, 0x1d, 0x83, 0xc5, 0xfe, 0x6c, 0x04, 0x1a, 0x65, 0x2a, 0x6d, 0xc8, 0x8d, 0xc0,
	0xe1, 0xc4, 0x7d, 0x05, 0x92, 0xf3, 0x7e, 0x96, 0x79, 0x67, 0xc1, 0xed, 0x7c, 0x70, 0x99, 0x73,
	0x07, 0x15, 0x4f, 0x40, 0xbb, 0x96, 0xa5, 0x1f, 0xfb, 0xcf, 0xe3, 0xcd, 0x4c, 0xb0, 0xb0, 0x75,
	0xa7, 0xdf, 0x42, 0x9b, 0x4e, 0x7f, 0x3d, 0xa8, 0xbb, 0xc9, 0x7b, 0x54, 0x6a, 0x49, 0x60, 0x1a,
	0x7c, 0x1e, 0xc2, 0x15, 0x7e, 0x38, 0x6d, 0x6e, 0x09, 0xae, 0xa5, 0x48, 0xc8, 0xcc, 0xe9, 0x64,
	0x2e, 0x16, 0x6f, 0x96, 0x8d, 0x2d, 0xae, 0xd0, 0xd8, 0xb4, 0xd1, 0xc0, 0x46, 0xf0, 0x4b, 0x54,
	0xec, 0xb2, 0x61, 0x06, 0x5e, 0xc8, 0xe4, 0x18, 0xdb, 0x96, 0x51, 0xc1, 0xa5, 0x8c, 0xa1, 0xd8,
	0x10, 0x70, 0x84, 0xb0, 0xfb, 0xfc, 0x11, 0x0f, 0xaf, 0x51, 0x6b, 0x39, 0x1e, 0x58, 0xfe, 0x2e,
	0xde, 0x99, 0x73, 0x64, 0x3c, 0xb4, 0x1f, 0x7a, 0x8d, 0x8a, 0x36, 0xed, 0xf3, 0x0b, 0x91, 0x53,
	0xbd, 0xfe, 0x43, 0x98, 0x63, 0x62, 0x1b, 0x49, 0x4b, 0x0c, 0xb8, 0xce, 0x59, 0xc7, 0xfc, 0x42,
	0x34, 0x3f, 0xfe, 0x1c, 0xd7, 0xbc, 0x5f, 0xe3, 0x9a, 0xf7, 0x67, 0x5c, 0xf3, 0xd0, 0xbe, 0x90,
	0x11, 0x19, 0x86, 0x8c, 0x29, 0x32, 0x64, 0x49, 0x48, 0xdc, 0x22, 0x30, 0x1b, 0xa0, 0x59, 0x3c,
	0x67, 0x49, 0x68, 0xff, 0xea, 0x6f, 0xbc, 0x0f, 0x8d, 0x28, 0xd6, 0x97, 0x83, 0x1e, 0x09, 0x44,
	0x9f, 0xda, 0x76, 0x6a, 0xda, 0xcd, 0x3a, 0x50, 0x34, 0x92, 0x69, 0xe0, 0x36, 0x08, 0x35, 0xc1,
	0x5e, 0xc1, 0x2e, 0x85, 0xa7, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0xe8, 0x7f, 0xd1, 0xbe, 0x5b,
	0x06, 0x00, 0x00,
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
	Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_ID, error)
	Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error)
	SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)
	StreamSearch(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamSearchClient, error)
	StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamSearchByIDClient, error)
	Insert(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamInsert(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamInsertClient, error)
	MultiInsert(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error)
	Update(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamUpdateClient, error)
	MultiUpdate(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error)
	Remove(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Empty, error)
	StreamRemove(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamRemoveClient, error)
	MultiRemove(ctx context.Context, in *payload.Object_IDs, opts ...grpc.CallOption) (*payload.Empty, error)
	GetObject(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_Vector, error)
	StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamGetObjectClient, error)
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

func (c *agentClient) Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_ID, error) {
	out := new(payload.Object_ID)
	err := c.cc.Invoke(ctx, "/core.Agent/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/core.Agent/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/core.Agent/SearchByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamSearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Agent_serviceDesc.Streams[0], "/core.Agent/StreamSearch", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentStreamSearchClient{stream}
	return x, nil
}

type Agent_StreamSearchClient interface {
	Send(*payload.Search_Request) error
	Recv() (*payload.Search_Response, error)
	grpc.ClientStream
}

type agentStreamSearchClient struct {
	grpc.ClientStream
}

func (x *agentStreamSearchClient) Send(m *payload.Search_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *agentStreamSearchClient) Recv() (*payload.Search_Response, error) {
	m := new(payload.Search_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *agentClient) StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamSearchByIDClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Agent_serviceDesc.Streams[1], "/core.Agent/StreamSearchByID", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentStreamSearchByIDClient{stream}
	return x, nil
}

type Agent_StreamSearchByIDClient interface {
	Send(*payload.Search_IDRequest) error
	Recv() (*payload.Search_Response, error)
	grpc.ClientStream
}

type agentStreamSearchByIDClient struct {
	grpc.ClientStream
}

func (x *agentStreamSearchByIDClient) Send(m *payload.Search_IDRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *agentStreamSearchByIDClient) Recv() (*payload.Search_Response, error) {
	m := new(payload.Search_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *agentClient) Insert(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/core.Agent/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamInsertClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Agent_serviceDesc.Streams[2], "/core.Agent/StreamInsert", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentStreamInsertClient{stream}
	return x, nil
}

type Agent_StreamInsertClient interface {
	Send(*payload.Object_Vector) error
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type agentStreamInsertClient struct {
	grpc.ClientStream
}

func (x *agentStreamInsertClient) Send(m *payload.Object_Vector) error {
	return x.ClientStream.SendMsg(m)
}

func (x *agentStreamInsertClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *agentClient) MultiInsert(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/core.Agent/MultiInsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) Update(ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/core.Agent/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamUpdateClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Agent_serviceDesc.Streams[3], "/core.Agent/StreamUpdate", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentStreamUpdateClient{stream}
	return x, nil
}

type Agent_StreamUpdateClient interface {
	Send(*payload.Object_Vector) error
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type agentStreamUpdateClient struct {
	grpc.ClientStream
}

func (x *agentStreamUpdateClient) Send(m *payload.Object_Vector) error {
	return x.ClientStream.SendMsg(m)
}

func (x *agentStreamUpdateClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *agentClient) MultiUpdate(ctx context.Context, in *payload.Object_Vectors, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/core.Agent/MultiUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) Remove(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/core.Agent/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamRemoveClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Agent_serviceDesc.Streams[4], "/core.Agent/StreamRemove", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentStreamRemoveClient{stream}
	return x, nil
}

type Agent_StreamRemoveClient interface {
	Send(*payload.Object_ID) error
	Recv() (*payload.Empty, error)
	grpc.ClientStream
}

type agentStreamRemoveClient struct {
	grpc.ClientStream
}

func (x *agentStreamRemoveClient) Send(m *payload.Object_ID) error {
	return x.ClientStream.SendMsg(m)
}

func (x *agentStreamRemoveClient) Recv() (*payload.Empty, error) {
	m := new(payload.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *agentClient) MultiRemove(ctx context.Context, in *payload.Object_IDs, opts ...grpc.CallOption) (*payload.Empty, error) {
	out := new(payload.Empty)
	err := c.cc.Invoke(ctx, "/core.Agent/MultiRemove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) GetObject(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/core.Agent/GetObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (Agent_StreamGetObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Agent_serviceDesc.Streams[5], "/core.Agent/StreamGetObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &agentStreamGetObjectClient{stream}
	return x, nil
}

type Agent_StreamGetObjectClient interface {
	Send(*payload.Object_ID) error
	Recv() (*payload.Object_Vector, error)
	grpc.ClientStream
}

type agentStreamGetObjectClient struct {
	grpc.ClientStream
}

func (x *agentStreamGetObjectClient) Send(m *payload.Object_ID) error {
	return x.ClientStream.SendMsg(m)
}

func (x *agentStreamGetObjectClient) Recv() (*payload.Object_Vector, error) {
	m := new(payload.Object_Vector)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
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
	Exists(context.Context, *payload.Object_ID) (*payload.Object_ID, error)
	Search(context.Context, *payload.Search_Request) (*payload.Search_Response, error)
	SearchByID(context.Context, *payload.Search_IDRequest) (*payload.Search_Response, error)
	StreamSearch(Agent_StreamSearchServer) error
	StreamSearchByID(Agent_StreamSearchByIDServer) error
	Insert(context.Context, *payload.Object_Vector) (*payload.Empty, error)
	StreamInsert(Agent_StreamInsertServer) error
	MultiInsert(context.Context, *payload.Object_Vectors) (*payload.Empty, error)
	Update(context.Context, *payload.Object_Vector) (*payload.Empty, error)
	StreamUpdate(Agent_StreamUpdateServer) error
	MultiUpdate(context.Context, *payload.Object_Vectors) (*payload.Empty, error)
	Remove(context.Context, *payload.Object_ID) (*payload.Empty, error)
	StreamRemove(Agent_StreamRemoveServer) error
	MultiRemove(context.Context, *payload.Object_IDs) (*payload.Empty, error)
	GetObject(context.Context, *payload.Object_ID) (*payload.Object_Vector, error)
	StreamGetObject(Agent_StreamGetObjectServer) error
	CreateIndex(context.Context, *payload.Control_CreateIndexRequest) (*payload.Empty, error)
	SaveIndex(context.Context, *payload.Empty) (*payload.Empty, error)
	CreateAndSaveIndex(context.Context, *payload.Control_CreateIndexRequest) (*payload.Empty, error)
	IndexInfo(context.Context, *payload.Empty) (*payload.Info_Index_Count, error)
}

// UnimplementedAgentServer can be embedded to have forward compatible implementations.
type UnimplementedAgentServer struct {
}

func (*UnimplementedAgentServer) Exists(ctx context.Context, req *payload.Object_ID) (*payload.Object_ID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (*UnimplementedAgentServer) Search(ctx context.Context, req *payload.Search_Request) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (*UnimplementedAgentServer) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchByID not implemented")
}
func (*UnimplementedAgentServer) StreamSearch(srv Agent_StreamSearchServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearch not implemented")
}
func (*UnimplementedAgentServer) StreamSearchByID(srv Agent_StreamSearchByIDServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearchByID not implemented")
}
func (*UnimplementedAgentServer) Insert(ctx context.Context, req *payload.Object_Vector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (*UnimplementedAgentServer) StreamInsert(srv Agent_StreamInsertServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamInsert not implemented")
}
func (*UnimplementedAgentServer) MultiInsert(ctx context.Context, req *payload.Object_Vectors) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiInsert not implemented")
}
func (*UnimplementedAgentServer) Update(ctx context.Context, req *payload.Object_Vector) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedAgentServer) StreamUpdate(srv Agent_StreamUpdateServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpdate not implemented")
}
func (*UnimplementedAgentServer) MultiUpdate(ctx context.Context, req *payload.Object_Vectors) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpdate not implemented")
}
func (*UnimplementedAgentServer) Remove(ctx context.Context, req *payload.Object_ID) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (*UnimplementedAgentServer) StreamRemove(srv Agent_StreamRemoveServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamRemove not implemented")
}
func (*UnimplementedAgentServer) MultiRemove(ctx context.Context, req *payload.Object_IDs) (*payload.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiRemove not implemented")
}
func (*UnimplementedAgentServer) GetObject(ctx context.Context, req *payload.Object_ID) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetObject not implemented")
}
func (*UnimplementedAgentServer) StreamGetObject(srv Agent_StreamGetObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamGetObject not implemented")
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

func _Agent_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Agent/Exists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).Exists(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Search_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Agent/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).Search(ctx, req.(*payload.Search_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_SearchByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Search_IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).SearchByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Agent/SearchByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).SearchByID(ctx, req.(*payload.Search_IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_StreamSearch_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AgentServer).StreamSearch(&agentStreamSearchServer{stream})
}

type Agent_StreamSearchServer interface {
	Send(*payload.Search_Response) error
	Recv() (*payload.Search_Request, error)
	grpc.ServerStream
}

type agentStreamSearchServer struct {
	grpc.ServerStream
}

func (x *agentStreamSearchServer) Send(m *payload.Search_Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *agentStreamSearchServer) Recv() (*payload.Search_Request, error) {
	m := new(payload.Search_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Agent_StreamSearchByID_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AgentServer).StreamSearchByID(&agentStreamSearchByIDServer{stream})
}

type Agent_StreamSearchByIDServer interface {
	Send(*payload.Search_Response) error
	Recv() (*payload.Search_IDRequest, error)
	grpc.ServerStream
}

type agentStreamSearchByIDServer struct {
	grpc.ServerStream
}

func (x *agentStreamSearchByIDServer) Send(m *payload.Search_Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *agentStreamSearchByIDServer) Recv() (*payload.Search_IDRequest, error) {
	m := new(payload.Search_IDRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Agent_Insert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Agent/Insert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).Insert(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_StreamInsert_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AgentServer).StreamInsert(&agentStreamInsertServer{stream})
}

type Agent_StreamInsertServer interface {
	Send(*payload.Empty) error
	Recv() (*payload.Object_Vector, error)
	grpc.ServerStream
}

type agentStreamInsertServer struct {
	grpc.ServerStream
}

func (x *agentStreamInsertServer) Send(m *payload.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *agentStreamInsertServer) Recv() (*payload.Object_Vector, error) {
	m := new(payload.Object_Vector)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Agent_MultiInsert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vectors)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).MultiInsert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Agent/MultiInsert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).MultiInsert(ctx, req.(*payload.Object_Vectors))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vector)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Agent/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).Update(ctx, req.(*payload.Object_Vector))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_StreamUpdate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AgentServer).StreamUpdate(&agentStreamUpdateServer{stream})
}

type Agent_StreamUpdateServer interface {
	Send(*payload.Empty) error
	Recv() (*payload.Object_Vector, error)
	grpc.ServerStream
}

type agentStreamUpdateServer struct {
	grpc.ServerStream
}

func (x *agentStreamUpdateServer) Send(m *payload.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *agentStreamUpdateServer) Recv() (*payload.Object_Vector, error) {
	m := new(payload.Object_Vector)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Agent_MultiUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_Vectors)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).MultiUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Agent/MultiUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).MultiUpdate(ctx, req.(*payload.Object_Vectors))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Agent/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).Remove(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_StreamRemove_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AgentServer).StreamRemove(&agentStreamRemoveServer{stream})
}

type Agent_StreamRemoveServer interface {
	Send(*payload.Empty) error
	Recv() (*payload.Object_ID, error)
	grpc.ServerStream
}

type agentStreamRemoveServer struct {
	grpc.ServerStream
}

func (x *agentStreamRemoveServer) Send(m *payload.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *agentStreamRemoveServer) Recv() (*payload.Object_ID, error) {
	m := new(payload.Object_ID)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Agent_MultiRemove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_IDs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).MultiRemove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Agent/MultiRemove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).MultiRemove(ctx, req.(*payload.Object_IDs))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_GetObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Object_ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).GetObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.Agent/GetObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).GetObject(ctx, req.(*payload.Object_ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_StreamGetObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AgentServer).StreamGetObject(&agentStreamGetObjectServer{stream})
}

type Agent_StreamGetObjectServer interface {
	Send(*payload.Object_Vector) error
	Recv() (*payload.Object_ID, error)
	grpc.ServerStream
}

type agentStreamGetObjectServer struct {
	grpc.ServerStream
}

func (x *agentStreamGetObjectServer) Send(m *payload.Object_Vector) error {
	return x.ServerStream.SendMsg(m)
}

func (x *agentStreamGetObjectServer) Recv() (*payload.Object_ID, error) {
	m := new(payload.Object_ID)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
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
			MethodName: "Exists",
			Handler:    _Agent_Exists_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _Agent_Search_Handler,
		},
		{
			MethodName: "SearchByID",
			Handler:    _Agent_SearchByID_Handler,
		},
		{
			MethodName: "Insert",
			Handler:    _Agent_Insert_Handler,
		},
		{
			MethodName: "MultiInsert",
			Handler:    _Agent_MultiInsert_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Agent_Update_Handler,
		},
		{
			MethodName: "MultiUpdate",
			Handler:    _Agent_MultiUpdate_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _Agent_Remove_Handler,
		},
		{
			MethodName: "MultiRemove",
			Handler:    _Agent_MultiRemove_Handler,
		},
		{
			MethodName: "GetObject",
			Handler:    _Agent_GetObject_Handler,
		},
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
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamSearch",
			Handler:       _Agent_StreamSearch_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamSearchByID",
			Handler:       _Agent_StreamSearchByID_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamInsert",
			Handler:       _Agent_StreamInsert_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamUpdate",
			Handler:       _Agent_StreamUpdate_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamRemove",
			Handler:       _Agent_StreamRemove_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamGetObject",
			Handler:       _Agent_StreamGetObject_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "core/agent.proto",
}
