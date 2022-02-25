//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: apis/proto/v1/agent/sidecar/sidecar.proto

package sidecar

import (
	reflect "reflect"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_apis_proto_v1_agent_sidecar_sidecar_proto protoreflect.FileDescriptor

var file_apis_proto_v1_agent_sidecar_sidecar_proto_rawDesc = []byte{
	0x0a, 0x29, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x69, 0x64, 0x65, 0x63, 0x61, 0x72, 0x2f, 0x73, 0x69,
	0x64, 0x65, 0x63, 0x61, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x73, 0x69, 0x64,
	0x65, 0x63, 0x61, 0x72, 0x2e, 0x76, 0x31, 0x32, 0x09, 0x0a, 0x07, 0x53, 0x69, 0x64, 0x65, 0x63,
	0x61, 0x72, 0x42, 0x6b, 0x0a, 0x23, 0x6f, 0x72, 0x67, 0x2e, 0x76, 0x64, 0x61, 0x61, 0x73, 0x2e,
	0x76, 0x61, 0x6c, 0x64, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x73, 0x69, 0x64, 0x65, 0x63, 0x61, 0x72, 0x42, 0x10, 0x56, 0x61, 0x6c, 0x64, 0x41,
	0x67, 0x65, 0x6e, 0x74, 0x53, 0x69, 0x64, 0x65, 0x63, 0x61, 0x72, 0x50, 0x01, 0x5a, 0x30, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x64, 0x61, 0x61, 0x73, 0x2f,
	0x76, 0x61, 0x6c, 0x64, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x76,
	0x31, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x69, 0x64, 0x65, 0x63, 0x61, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_apis_proto_v1_agent_sidecar_sidecar_proto_goTypes = []interface{}{}
var file_apis_proto_v1_agent_sidecar_sidecar_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_apis_proto_v1_agent_sidecar_sidecar_proto_init() }
func file_apis_proto_v1_agent_sidecar_sidecar_proto_init() {
	if File_apis_proto_v1_agent_sidecar_sidecar_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_apis_proto_v1_agent_sidecar_sidecar_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_apis_proto_v1_agent_sidecar_sidecar_proto_goTypes,
		DependencyIndexes: file_apis_proto_v1_agent_sidecar_sidecar_proto_depIdxs,
	}.Build()
	File_apis_proto_v1_agent_sidecar_sidecar_proto = out.File
	file_apis_proto_v1_agent_sidecar_sidecar_proto_rawDesc = nil
	file_apis_proto_v1_agent_sidecar_sidecar_proto_goTypes = nil
	file_apis_proto_v1_agent_sidecar_sidecar_proto_depIdxs = nil
}
