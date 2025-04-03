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
// source: v1/agent/core/agent.proto

package core

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

var File_v1_agent_core_agent_proto protoreflect.FileDescriptor

const file_v1_agent_core_agent_proto_rawDesc = "" +
	"\n" +
	"\x19v1/agent/core/agent.proto\x12\acore.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x18v1/payload/payload.proto2\xb7\x02\n" +
	"\x05Agent\x12k\n" +
	"\vCreateIndex\x12&.payload.v1.Control.CreateIndexRequest\x1a\x11.payload.v1.Empty\"!\x82\xd3\xe4\x93\x02\x1b\x12\x19/index/create/{pool_size}\x12F\n" +
	"\tSaveIndex\x12\x11.payload.v1.Empty\x1a\x11.payload.v1.Empty\"\x13\x82\xd3\xe4\x93\x02\r\x12\v/index/save\x12y\n" +
	"\x12CreateAndSaveIndex\x12&.payload.v1.Control.CreateIndexRequest\x1a\x11.payload.v1.Empty\"(\x82\xd3\xe4\x93\x02\"\x12 /index/createandsave/{pool_size}Bc\n" +
	" org.vdaas.vald.api.v1.agent.coreB\tValdAgentP\x01Z2github.com/vdaas/vald/apis/grpc/v1/agent/core;coreb\x06proto3"

var file_v1_agent_core_agent_proto_goTypes = []any{
	(*payload.Control_CreateIndexRequest)(nil), // 0: payload.v1.Control.CreateIndexRequest
	(*payload.Empty)(nil),                      // 1: payload.v1.Empty
}

var file_v1_agent_core_agent_proto_depIdxs = []int32{
	0, // 0: core.v1.Agent.CreateIndex:input_type -> payload.v1.Control.CreateIndexRequest
	1, // 1: core.v1.Agent.SaveIndex:input_type -> payload.v1.Empty
	0, // 2: core.v1.Agent.CreateAndSaveIndex:input_type -> payload.v1.Control.CreateIndexRequest
	1, // 3: core.v1.Agent.CreateIndex:output_type -> payload.v1.Empty
	1, // 4: core.v1.Agent.SaveIndex:output_type -> payload.v1.Empty
	1, // 5: core.v1.Agent.CreateAndSaveIndex:output_type -> payload.v1.Empty
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v1_agent_core_agent_proto_init() }
func file_v1_agent_core_agent_proto_init() {
	if File_v1_agent_core_agent_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_agent_core_agent_proto_rawDesc), len(file_v1_agent_core_agent_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_agent_core_agent_proto_goTypes,
		DependencyIndexes: file_v1_agent_core_agent_proto_depIdxs,
	}.Build()
	File_v1_agent_core_agent_proto = out.File
	file_v1_agent_core_agent_proto_goTypes = nil
	file_v1_agent_core_agent_proto_depIdxs = nil
}
