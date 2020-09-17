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

package index

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
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

func init() {
	proto.RegisterFile("apis/proto/v1/manager/index/index_manager.proto", fileDescriptor_0152ec67984b188e)
}

var fileDescriptor_0152ec67984b188e = []byte{
	// 248 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4f, 0x2c, 0xc8, 0x2c,
	0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x33, 0xd4, 0xcf, 0x4d, 0xcc, 0x4b, 0x4c, 0x4f,
	0x2d, 0xd2, 0xcf, 0xcc, 0x4b, 0x49, 0xad, 0x80, 0x90, 0xf1, 0x50, 0x31, 0x3d, 0xb0, 0x22, 0x21,
	0x5e, 0x14, 0x41, 0x29, 0x65, 0x54, 0xfd, 0x05, 0x89, 0x95, 0x39, 0xf9, 0x89, 0x29, 0x30, 0x1a,
	0xa2, 0x47, 0x4a, 0x26, 0x3d, 0x3f, 0x3f, 0x3d, 0x27, 0x15, 0x64, 0x97, 0x7e, 0x62, 0x5e, 0x5e,
	0x7e, 0x49, 0x62, 0x49, 0x66, 0x7e, 0x5e, 0x31, 0x44, 0xd6, 0x28, 0x84, 0x8b, 0xd5, 0x13, 0x64,
	0xa6, 0x90, 0x37, 0x17, 0x27, 0x98, 0xe1, 0x99, 0x97, 0x96, 0x2f, 0xc4, 0xa7, 0x07, 0x33, 0xc3,
	0x35, 0xb7, 0xa0, 0xa4, 0x52, 0x4a, 0x12, 0xce, 0x07, 0x49, 0xeb, 0x81, 0x15, 0xea, 0x39, 0xe7,
	0x97, 0xe6, 0x95, 0x28, 0x09, 0x37, 0x5d, 0x7e, 0x32, 0x99, 0x89, 0x57, 0x88, 0x1b, 0xee, 0xec,
	0xb4, 0x7c, 0xa7, 0xf2, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e,
	0x91, 0x4b, 0x39, 0xbf, 0x28, 0x5d, 0xaf, 0x2c, 0x25, 0x31, 0xb1, 0x58, 0xaf, 0x2c, 0x31, 0x27,
	0x45, 0x2f, 0xb1, 0x20, 0x53, 0xaf, 0xcc, 0x50, 0x0f, 0xe6, 0x35, 0xb0, 0x3e, 0x27, 0x81, 0xb0,
	0xc4, 0x9c, 0x14, 0xb0, 0xc1, 0xbe, 0x10, 0xf1, 0x00, 0xc6, 0x28, 0x83, 0xf4, 0xcc, 0x92, 0x8c,
	0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d, 0xb0, 0x7e, 0x7d, 0x90, 0x7e, 0x48, 0xa8, 0xa5, 0x17,
	0x15, 0x24, 0x63, 0x04, 0x5a, 0x12, 0x1b, 0xd8, 0x57, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x7c, 0x30, 0xa5, 0x15, 0x5a, 0x01, 0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/index_manager.Index/IndexInfo", in, out, opts...)
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
		FullMethod: "/index_manager.Index/IndexInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndexServer).IndexInfo(ctx, req.(*payload.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Index_serviceDesc = grpc.ServiceDesc{
	ServiceName: "index_manager.Index",
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
