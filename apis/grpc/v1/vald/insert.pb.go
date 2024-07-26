//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: v1/vald/insert.proto

package vald

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

var File_v1_vald_insert_proto protoreflect.FileDescriptor

var file_v1_vald_insert_proto_rawDesc = []byte{
	0x0a, 0x14, 0x76, 0x31, 0x2f, 0x76, 0x61, 0x6c, 0x64, 0x2f, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x76, 0x61, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x76,
	0x31, 0x2f, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2f, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x9f, 0x02, 0x0a, 0x06, 0x49, 0x6e, 0x73, 0x65,
	0x72, 0x74, 0x12, 0x55, 0x0a, 0x06, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x1a, 0x2e, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74,
	0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x4c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x12, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x3a, 0x01, 0x2a,
	0x22, 0x07, 0x2f, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x53, 0x0a, 0x0c, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x1a, 0x2e, 0x70, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e,
	0x76, 0x31, 0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x69,
	0x0a, 0x0b, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x1f, 0x2e,
	0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73, 0x65, 0x72,
	0x74, 0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c,
	0x2e, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x1b, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x15, 0x3a, 0x01, 0x2a, 0x22, 0x10, 0x2f, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x74,
	0x2f, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x65, 0x42, 0x53, 0x0a, 0x1a, 0x6f, 0x72, 0x67,
	0x2e, 0x76, 0x64, 0x61, 0x61, 0x73, 0x2e, 0x76, 0x61, 0x6c, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x76, 0x61, 0x6c, 0x64, 0x42, 0x0a, 0x56, 0x61, 0x6c, 0x64, 0x49, 0x6e, 0x73,
	0x65, 0x72, 0x74, 0x50, 0x01, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x76, 0x64, 0x61, 0x61, 0x73, 0x2f, 0x76, 0x61, 0x6c, 0x64, 0x2f, 0x61, 0x70, 0x69,
	0x73, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x61, 0x6c, 0x64, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_v1_vald_insert_proto_goTypes = []any{
	(*payload.Insert_Request)(nil),        // 0: payload.v1.Insert.Request
	(*payload.Insert_MultiRequest)(nil),   // 1: payload.v1.Insert.MultiRequest
	(*payload.Object_Location)(nil),       // 2: payload.v1.Object.Location
	(*payload.Object_StreamLocation)(nil), // 3: payload.v1.Object.StreamLocation
	(*payload.Object_Locations)(nil),      // 4: payload.v1.Object.Locations
}

var file_v1_vald_insert_proto_depIdxs = []int32{
	0, // 0: vald.v1.Insert.Insert:input_type -> payload.v1.Insert.Request
	0, // 1: vald.v1.Insert.StreamInsert:input_type -> payload.v1.Insert.Request
	1, // 2: vald.v1.Insert.MultiInsert:input_type -> payload.v1.Insert.MultiRequest
	2, // 3: vald.v1.Insert.Insert:output_type -> payload.v1.Object.Location
	3, // 4: vald.v1.Insert.StreamInsert:output_type -> payload.v1.Object.StreamLocation
	4, // 5: vald.v1.Insert.MultiInsert:output_type -> payload.v1.Object.Locations
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v1_vald_insert_proto_init() }
func file_v1_vald_insert_proto_init() {
	if File_v1_vald_insert_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_v1_vald_insert_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_vald_insert_proto_goTypes,
		DependencyIndexes: file_v1_vald_insert_proto_depIdxs,
	}.Build()
	File_v1_vald_insert_proto = out.File
	file_v1_vald_insert_proto_rawDesc = nil
	file_v1_vald_insert_proto_goTypes = nil
	file_v1_vald_insert_proto_depIdxs = nil
}
