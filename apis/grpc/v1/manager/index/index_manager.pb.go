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

package index

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
	proto.RegisterFile("apis/proto/v1/manager/index/index_manager.proto", fileDescriptor_0152ec67984b188e)
}

var fileDescriptor_0152ec67984b188e = []byte{
	// 252 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4f, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0xcf, 0x4d, 0xcc, 0x4b, 0x4c, 0x4f,
	0x2d, 0xd2, 0xcf, 0xcc, 0x4b, 0x49, 0xad, 0x80, 0x90, 0xf1, 0x50, 0x31, 0x3d, 0xb0, 0x22, 0x21,
	0x01, 0x18, 0x17, 0x2c, 0xa9, 0x57, 0x66, 0x28, 0xa5, 0x8c, 0x6a, 0x44, 0x41, 0x62, 0x65, 0x4e,
	0x7e, 0x62, 0x0a, 0x8c, 0x86, 0x68, 0x93, 0x92, 0x49, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0x05, 0x59,
	0xa7, 0x9f, 0x98, 0x97, 0x97, 0x5f, 0x92, 0x58, 0x92, 0x99, 0x9f, 0x57, 0x0c, 0x91, 0x35, 0x8a,
	0xe2, 0x62, 0xf5, 0x04, 0x19, 0x27, 0x14, 0xc8, 0xc5, 0x09, 0x66, 0x78, 0xe6, 0xa5, 0xe5, 0x0b,
	0x09, 0xea, 0xc1, 0xcc, 0x28, 0x33, 0xd4, 0x73, 0xcd, 0x2d, 0x28, 0xa9, 0x94, 0x92, 0x41, 0x16,
	0x02, 0x29, 0xd2, 0x03, 0x2b, 0xd7, 0x73, 0xce, 0x2f, 0xcd, 0x2b, 0x51, 0x12, 0x6e, 0xba, 0xfc,
	0x64, 0x32, 0x13, 0xaf, 0x10, 0x37, 0xdc, 0xfd, 0x69, 0xf9, 0x4e, 0xe5, 0x27, 0x1e, 0xc9, 0x31,
	0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0x23, 0x97, 0x72, 0x7e, 0x51, 0xba, 0x5e, 0x59,
	0x4a, 0x62, 0x62, 0xb1, 0x5e, 0x59, 0x62, 0x4e, 0x8a, 0x5e, 0x62, 0x41, 0x26, 0xc8, 0x28, 0x14,
	0x4f, 0x39, 0x09, 0x84, 0x25, 0xe6, 0xa4, 0x80, 0x0d, 0xf6, 0x85, 0x88, 0x07, 0x30, 0x46, 0x19,
	0xa4, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea, 0x83, 0xf5, 0xeb, 0x83, 0xf4,
	0x43, 0x82, 0x2f, 0xbd, 0xa8, 0x20, 0x19, 0x23, 0xf4, 0x92, 0xd8, 0xc0, 0x7e, 0x33, 0x06, 0x04,
	0x00, 0x00, 0xff, 0xff, 0x35, 0xf2, 0xc1, 0x66, 0x63, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// IndexClient is the client API for Index service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IndexClient interface {
	IndexInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Index_Count, error)
}

type indexClient struct {
	cc *grpc.ClientConn
}

func NewIndexClient(cc *grpc.ClientConn) IndexClient {
	return &indexClient{cc}
}

func (c *indexClient) IndexInfo(ctx context.Context, in *payload.Empty, opts ...grpc.CallOption) (*payload.Info_Index_Count, error) {
	out := new(payload.Info_Index_Count)
	err := c.cc.Invoke(ctx, "/manager.index.v1.Index/IndexInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IndexServer is the server API for Index service.
type IndexServer interface {
	IndexInfo(context.Context, *payload.Empty) (*payload.Info_Index_Count, error)
}

// UnimplementedIndexServer can be embedded to have forward compatible implementations.
type UnimplementedIndexServer struct {
}

func (*UnimplementedIndexServer) IndexInfo(ctx context.Context, req *payload.Empty) (*payload.Info_Index_Count, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IndexInfo not implemented")
}

func RegisterIndexServer(s *grpc.Server, srv IndexServer) {
	s.RegisterService(&_Index_serviceDesc, srv)
}

func _Index_IndexInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(payload.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndexServer).IndexInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/manager.index.v1.Index/IndexInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).IndexInfo(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Index_serviceDesc = grpc.ServiceDesc{
	ServiceName: "manager.index.v1.Index",
	HandlerType: (*IndexServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IndexInfo",
			Handler:    _Index_IndexInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apis/proto/v1/manager/index/index_manager.proto",
}
