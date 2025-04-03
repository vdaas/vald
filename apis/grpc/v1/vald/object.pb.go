//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: v1/vald/object.proto

package vald

import (
	reflect "reflect"
	unsafe "unsafe"

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

var File_v1_vald_object_proto protoreflect.FileDescriptor

const file_v1_vald_object_proto_rawDesc = "" +
	"\n" +
	"\x14v1/vald/object.proto\x12\avald.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x18v1/payload/payload.proto2\xf5\x03\n" +
	"\x06Object\x12L\n" +
	"\x06Exists\x12\x15.payload.v1.Object.ID\x1a\x15.payload.v1.Object.ID\"\x14\x82\xd3\xe4\x93\x02\x0e\x12\f/exists/{id}\x12a\n" +
	"\tGetObject\x12 .payload.v1.Object.VectorRequest\x1a\x19.payload.v1.Object.Vector\"\x17\x82\xd3\xe4\x93\x02\x11\x12\x0f/object/{id.id}\x12Z\n" +
	"\x0fStreamGetObject\x12 .payload.v1.Object.VectorRequest\x1a\x1f.payload.v1.Object.StreamVector\"\x00(\x010\x01\x12m\n" +
	"\x10StreamListObject\x12\x1f.payload.v1.Object.List.Request\x1a .payload.v1.Object.List.Response\"\x14\x82\xd3\xe4\x93\x02\x0e\x12\f/object/list0\x01\x12o\n" +
	"\fGetTimestamp\x12#.payload.v1.Object.TimestampRequest\x1a\x1c.payload.v1.Object.Timestamp\"\x1c\x82\xd3\xe4\x93\x02\x16\x12\x14/object/meta/{id.id}BS\n" +
	"\x1aorg.vdaas.vald.api.v1.valdB\n" +
	"ValdObjectP\x01Z'github.com/vdaas/vald/apis/grpc/v1/valdb\x06proto3"

var file_v1_vald_object_proto_goTypes = []any{
	(*payload.Object_ID)(nil),               // 0: payload.v1.Object.ID
	(*payload.Object_VectorRequest)(nil),    // 1: payload.v1.Object.VectorRequest
	(*payload.Object_List_Request)(nil),     // 2: payload.v1.Object.List.Request
	(*payload.Object_TimestampRequest)(nil), // 3: payload.v1.Object.TimestampRequest
	(*payload.Object_Vector)(nil),           // 4: payload.v1.Object.Vector
	(*payload.Object_StreamVector)(nil),     // 5: payload.v1.Object.StreamVector
	(*payload.Object_List_Response)(nil),    // 6: payload.v1.Object.List.Response
	(*payload.Object_Timestamp)(nil),        // 7: payload.v1.Object.Timestamp
}

var file_v1_vald_object_proto_depIdxs = []int32{
	0, // 0: vald.v1.Object.Exists:input_type -> payload.v1.Object.ID
	1, // 1: vald.v1.Object.GetObject:input_type -> payload.v1.Object.VectorRequest
	1, // 2: vald.v1.Object.StreamGetObject:input_type -> payload.v1.Object.VectorRequest
	2, // 3: vald.v1.Object.StreamListObject:input_type -> payload.v1.Object.List.Request
	3, // 4: vald.v1.Object.GetTimestamp:input_type -> payload.v1.Object.TimestampRequest
	0, // 5: vald.v1.Object.Exists:output_type -> payload.v1.Object.ID
	4, // 6: vald.v1.Object.GetObject:output_type -> payload.v1.Object.Vector
	5, // 7: vald.v1.Object.StreamGetObject:output_type -> payload.v1.Object.StreamVector
	6, // 8: vald.v1.Object.StreamListObject:output_type -> payload.v1.Object.List.Response
	7, // 9: vald.v1.Object.GetTimestamp:output_type -> payload.v1.Object.Timestamp
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v1_vald_object_proto_init() }
func file_v1_vald_object_proto_init() {
	if File_v1_vald_object_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_vald_object_proto_rawDesc), len(file_v1_vald_object_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_vald_object_proto_goTypes,
		DependencyIndexes: file_v1_vald_object_proto_depIdxs,
	}.Build()
	File_v1_vald_object_proto = out.File
	file_v1_vald_object_proto_goTypes = nil
	file_v1_vald_object_proto_depIdxs = nil
}
