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

func init() { proto.RegisterFile("index/index_manager.proto", fileDescriptor_11357116787cb271) }

var fileDescriptor_11357116787cb271 = []byte{
	// 246 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcc, 0xcc, 0x4b, 0x49,
	0xad, 0xd0, 0x07, 0x93, 0xf1, 0xb9, 0x89, 0x79, 0x89, 0xe9, 0xa9, 0x45, 0x7a, 0x05, 0x45, 0xf9,
	0x25, 0xf9, 0x42, 0xbc, 0x28, 0x82, 0x52, 0xbc, 0x05, 0x89, 0x95, 0x39, 0xf9, 0x89, 0x29, 0x10,
	0x59, 0x29, 0x99, 0xf4, 0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0xfd, 0xc4, 0x82, 0x4c, 0xfd, 0xc4, 0xbc,
	0xbc, 0xfc, 0x92, 0xc4, 0x92, 0xcc, 0xfc, 0xbc, 0x62, 0xa8, 0x2c, 0x4f, 0x41, 0x92, 0x7e, 0x7a,
	0x61, 0x0e, 0x84, 0x67, 0x14, 0xc5, 0xc5, 0xea, 0x09, 0x32, 0x4b, 0xc8, 0x9b, 0x8b, 0x13, 0xcc,
	0xf0, 0xcc, 0x4b, 0xcb, 0x17, 0xe2, 0xd3, 0x83, 0x99, 0xe8, 0x9a, 0x5b, 0x50, 0x52, 0x29, 0x25,
	0x09, 0xe7, 0x83, 0xa4, 0xf5, 0xc0, 0x0a, 0xf5, 0x9c, 0xf3, 0x4b, 0xf3, 0x4a, 0x94, 0x84, 0x9b,
	0x2e, 0x3f, 0x99, 0xcc, 0xc4, 0x2b, 0xc4, 0xad, 0x0f, 0x73, 0x6e, 0x5a, 0xbe, 0x14, 0xcb, 0x86,
	0x07, 0xf2, 0x4c, 0x4e, 0xb9, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91,
	0x1c, 0x23, 0x97, 0x4c, 0x7e, 0x51, 0xba, 0x5e, 0x59, 0x4a, 0x62, 0x62, 0xb1, 0x5e, 0x59, 0x62,
	0x4e, 0x8a, 0x1e, 0xcc, 0x47, 0x60, 0x6d, 0x4e, 0x02, 0x61, 0x89, 0x39, 0x29, 0x60, 0x73, 0x7d,
	0x21, 0xe2, 0x01, 0x8c, 0x51, 0xba, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9,
	0xfa, 0x60, 0x8d, 0xfa, 0x20, 0x8d, 0x20, 0x6f, 0x15, 0xeb, 0xa7, 0x17, 0x15, 0x24, 0xeb, 0x43,
	0x8d, 0x80, 0xd8, 0x9c, 0xc4, 0x06, 0xf6, 0x91, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x98, 0x1d,
	0x6e, 0xfd, 0x38, 0x01, 0x00, 0x00,
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
	Metadata: "index/index_manager.proto",
}
