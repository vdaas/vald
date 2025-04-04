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
// source: v1/discoverer/discoverer.proto

package discoverer

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

var File_v1_discoverer_discoverer_proto protoreflect.FileDescriptor

const file_v1_discoverer_discoverer_proto_rawDesc = "" +
	"\n" +
	"\x1ev1/discoverer/discoverer.proto\x12\rdiscoverer.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x18v1/payload/payload.proto2\xa9\x02\n" +
	"\n" +
	"Discoverer\x12X\n" +
	"\x04Pods\x12\x1e.payload.v1.Discoverer.Request\x1a\x15.payload.v1.Info.Pods\"\x19\x82\xd3\xe4\x93\x02\x13:\x01*\"\x0e/discover/pods\x12[\n" +
	"\x05Nodes\x12\x1e.payload.v1.Discoverer.Request\x1a\x16.payload.v1.Info.Nodes\"\x1a\x82\xd3\xe4\x93\x02\x14:\x01*\"\x0f/discover/nodes\x12d\n" +
	"\bServices\x12\x1e.payload.v1.Discoverer.Request\x1a\x19.payload.v1.Info.Services\"\x1d\x82\xd3\xe4\x93\x02\x17:\x01*\"\x12/discover/servicesBc\n" +
	" org.vdaas.vald.api.v1.discovererB\x0eValdDiscovererP\x01Z-github.com/vdaas/vald/apis/grpc/v1/discovererb\x06proto3"

var file_v1_discoverer_discoverer_proto_goTypes = []any{
	(*payload.Discoverer_Request)(nil), // 0: payload.v1.Discoverer.Request
	(*payload.Info_Pods)(nil),          // 1: payload.v1.Info.Pods
	(*payload.Info_Nodes)(nil),         // 2: payload.v1.Info.Nodes
	(*payload.Info_Services)(nil),      // 3: payload.v1.Info.Services
}

var file_v1_discoverer_discoverer_proto_depIdxs = []int32{
	0, // 0: discoverer.v1.Discoverer.Pods:input_type -> payload.v1.Discoverer.Request
	0, // 1: discoverer.v1.Discoverer.Nodes:input_type -> payload.v1.Discoverer.Request
	0, // 2: discoverer.v1.Discoverer.Services:input_type -> payload.v1.Discoverer.Request
	1, // 3: discoverer.v1.Discoverer.Pods:output_type -> payload.v1.Info.Pods
	2, // 4: discoverer.v1.Discoverer.Nodes:output_type -> payload.v1.Info.Nodes
	3, // 5: discoverer.v1.Discoverer.Services:output_type -> payload.v1.Info.Services
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_v1_discoverer_discoverer_proto_init() }
func file_v1_discoverer_discoverer_proto_init() {
	if File_v1_discoverer_discoverer_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_discoverer_discoverer_proto_rawDesc), len(file_v1_discoverer_discoverer_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_discoverer_discoverer_proto_goTypes,
		DependencyIndexes: file_v1_discoverer_discoverer_proto_depIdxs,
	}.Build()
	File_v1_discoverer_discoverer_proto = out.File
	file_v1_discoverer_discoverer_proto_goTypes = nil
	file_v1_discoverer_discoverer_proto_depIdxs = nil
}
