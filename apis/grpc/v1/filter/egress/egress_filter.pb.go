//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.2
// source: apis/proto/v1/filter/egress/egress_filter.proto

package egress

import (
	reflect "reflect"

	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_apis_proto_v1_filter_egress_egress_filter_proto protoreflect.FileDescriptor

var file_apis_proto_v1_filter_egress_egress_filter_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2f, 0x65, 0x67, 0x72, 0x65, 0x73, 0x73, 0x2f, 0x65, 0x67,
	0x72, 0x65, 0x73, 0x73, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x10, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2e, 0x65, 0x67, 0x72, 0x65, 0x73, 0x73,
	0x2e, 0x76, 0x31, 0x1a, 0x23, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x76, 0x31, 0x2f, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2f, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xe0, 0x01, 0x0a, 0x06, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x12, 0x6e, 0x0a, 0x0e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x44, 0x69, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x12, 0x1b, 0x2e, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76,
	0x31, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x1a, 0x1b, 0x2e, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x4f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x22,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1c, 0x3a, 0x01, 0x2a, 0x22, 0x17, 0x2f, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x2f, 0x65, 0x67, 0x72, 0x65, 0x73, 0x73, 0x2f, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x12, 0x66, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x56, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x12, 0x19, 0x2e, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76, 0x31, 0x2e,
	0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x1a, 0x19, 0x2e,
	0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x2e, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a,
	0x3a, 0x01, 0x2a, 0x22, 0x15, 0x2f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2f, 0x65, 0x67, 0x72,
	0x65, 0x73, 0x73, 0x2f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x42, 0x6b, 0x0a, 0x23, 0x6f, 0x72,
	0x67, 0x2e, 0x76, 0x64, 0x61, 0x61, 0x73, 0x2e, 0x76, 0x61, 0x6c, 0x64, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2e, 0x65, 0x67, 0x72, 0x65, 0x73,
	0x73, 0x42, 0x10, 0x56, 0x61, 0x6c, 0x64, 0x45, 0x67, 0x72, 0x65, 0x73, 0x73, 0x46, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x50, 0x01, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x76, 0x64, 0x61, 0x61, 0x73, 0x2f, 0x76, 0x61, 0x6c, 0x64, 0x2f, 0x61, 0x70, 0x69,
	0x73, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x2f, 0x65, 0x67, 0x72, 0x65, 0x73, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_apis_proto_v1_filter_egress_egress_filter_proto_goTypes = []interface{}{
	(*payload.Object_Distance)(nil), // 0: payload.v1.Object.Distance
	(*payload.Object_Vector)(nil),   // 1: payload.v1.Object.Vector
}
var file_apis_proto_v1_filter_egress_egress_filter_proto_depIdxs = []int32{
	0, // 0: filter.egress.v1.Filter.FilterDistance:input_type -> payload.v1.Object.Distance
	1, // 1: filter.egress.v1.Filter.FilterVector:input_type -> payload.v1.Object.Vector
	0, // 2: filter.egress.v1.Filter.FilterDistance:output_type -> payload.v1.Object.Distance
	1, // 3: filter.egress.v1.Filter.FilterVector:output_type -> payload.v1.Object.Vector
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_apis_proto_v1_filter_egress_egress_filter_proto_init() }
func file_apis_proto_v1_filter_egress_egress_filter_proto_init() {
	if File_apis_proto_v1_filter_egress_egress_filter_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_apis_proto_v1_filter_egress_egress_filter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_apis_proto_v1_filter_egress_egress_filter_proto_goTypes,
		DependencyIndexes: file_apis_proto_v1_filter_egress_egress_filter_proto_depIdxs,
	}.Build()
	File_apis_proto_v1_filter_egress_egress_filter_proto = out.File
	file_apis_proto_v1_filter_egress_egress_filter_proto_rawDesc = nil
	file_apis_proto_v1_filter_egress_egress_filter_proto_goTypes = nil
	file_apis_proto_v1_filter_egress_egress_filter_proto_depIdxs = nil
}
