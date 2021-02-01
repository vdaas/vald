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
	proto.RegisterFile("apis/proto/v1/manager/index/index_manager.proto", fileDescriptor_0152ec67984b188e)
}
func init() {
	golang_proto.RegisterFile("apis/proto/v1/manager/index/index_manager.proto", fileDescriptor_0152ec67984b188e)
}

var fileDescriptor_0152ec67984b188e = []byte{
	// 286 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0xbd, 0x4a, 0xc4, 0x40,
	0x10, 0xc7, 0x59, 0x41, 0xc1, 0x88, 0x70, 0xc6, 0x2e, 0x48, 0x0a, 0xaf, 0x76, 0xd6, 0x68, 0x61,
	0x7f, 0x62, 0x71, 0x85, 0xa0, 0x8d, 0xc5, 0x35, 0x32, 0xb9, 0x24, 0x6b, 0x20, 0xd9, 0x59, 0x92,
	0xcd, 0xe2, 0x55, 0x82, 0xaf, 0xe0, 0x0b, 0x59, 0x5e, 0x29, 0xf8, 0x02, 0x92, 0xf3, 0x41, 0x64,
	0x77, 0x2f, 0x72, 0xa7, 0x4d, 0x32, 0x33, 0xfb, 0xfb, 0xcf, 0x67, 0xc0, 0x51, 0x95, 0x2d, 0x57,
	0x0d, 0x69, 0xe2, 0x26, 0xe1, 0x35, 0x4a, 0x14, 0x79, 0xc3, 0x4b, 0x99, 0xe5, 0xcf, 0xfe, 0xfb,
	0xb8, 0x8e, 0x81, 0x83, 0xc2, 0xd1, 0xe0, 0xba, 0x47, 0x30, 0x49, 0x34, 0xde, 0x4e, 0xa1, 0x70,
	0x51, 0x11, 0x66, 0xc3, 0xdf, 0xcb, 0xa2, 0x33, 0x51, 0xea, 0xa7, 0x2e, 0x85, 0x39, 0xd5, 0x5c,
	0x90, 0x20, 0xcf, 0xa7, 0x5d, 0xe1, 0x3c, 0x2f, 0xb6, 0xd6, 0x1a, 0xbf, 0xfa, 0x8b, 0x0b, 0x22,
	0x51, 0xe5, 0xae, 0x92, 0x37, 0x6d, 0xe3, 0x1c, 0xa5, 0x24, 0x8d, 0xba, 0x24, 0xd9, 0x7a, 0xe1,
	0xc5, 0x2c, 0xd8, 0x9d, 0xda, 0xc6, 0xc2, 0xfb, 0x60, 0xdf, 0x19, 0x53, 0x59, 0x50, 0x78, 0x04,
	0x43, 0x37, 0x26, 0x81, 0x9b, 0x5a, 0xe9, 0x45, 0x74, 0xb2, 0x19, 0xb2, 0x10, 0x38, 0x1c, 0xae,
	0xa9, 0x93, 0xfa, 0xf4, 0xf8, 0xf5, 0xf3, 0xfb, 0x6d, 0xe7, 0x30, 0x3c, 0xf8, 0xdd, 0x44, 0x41,
	0x93, 0x97, 0x65, 0x1f, 0xb3, 0x8f, 0x3e, 0x66, 0x5f, 0x7d, 0xcc, 0xde, 0x57, 0x31, 0x5b, 0xae,
	0x62, 0x16, 0x8c, 0xa9, 0x11, 0x60, 0x32, 0xc4, 0x16, 0x0c, 0x56, 0x19, 0xa0, 0x2a, 0x6d, 0xca,
	0xad, 0x35, 0x4d, 0x46, 0x0f, 0x58, 0x65, 0xae, 0xc0, 0xad, 0x8f, 0xdf, 0xb1, 0xd9, 0xf9, 0xc6,
	0x8c, 0x4e, 0xcf, 0xad, 0xde, 0x1f, 0x44, 0x34, 0x6a, 0xfe, 0xef, 0x1e, 0xe9, 0x9e, 0x9b, 0xf1,
	0xf2, 0x27, 0x00, 0x00, 0xff, 0xff, 0x8b, 0x11, 0x3e, 0x32, 0xb5, 0x01, 0x00, 0x00,
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
