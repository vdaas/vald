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
// source: v1/vald/search.proto

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

var File_v1_vald_search_proto protoreflect.FileDescriptor

const file_v1_vald_search_proto_rawDesc = "" +
	"\n" +
	"\x14v1/vald/search.proto\x12\avald.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x18v1/payload/payload.proto2\xd0\t\n" +
	"\x06Search\x12U\n" +
	"\x06Search\x12\x1a.payload.v1.Search.Request\x1a\x1b.payload.v1.Search.Response\"\x12\x82\xd3\xe4\x93\x02\f:\x01*\"\a/search\x12^\n" +
	"\n" +
	"SearchByID\x12\x1c.payload.v1.Search.IDRequest\x1a\x1b.payload.v1.Search.Response\"\x15\x82\xd3\xe4\x93\x02\x0f:\x01*\"\n" +
	"/search/id\x12S\n" +
	"\fStreamSearch\x12\x1a.payload.v1.Search.Request\x1a!.payload.v1.Search.StreamResponse\"\x00(\x010\x01\x12Y\n" +
	"\x10StreamSearchByID\x12\x1c.payload.v1.Search.IDRequest\x1a!.payload.v1.Search.StreamResponse\"\x00(\x010\x01\x12i\n" +
	"\vMultiSearch\x12\x1f.payload.v1.Search.MultiRequest\x1a\x1c.payload.v1.Search.Responses\"\x1b\x82\xd3\xe4\x93\x02\x15:\x01*\"\x10/search/multiple\x12r\n" +
	"\x0fMultiSearchByID\x12!.payload.v1.Search.MultiIDRequest\x1a\x1c.payload.v1.Search.Responses\"\x1e\x82\xd3\xe4\x93\x02\x18:\x01*\"\x13/search/id/multiple\x12a\n" +
	"\fLinearSearch\x12\x1a.payload.v1.Search.Request\x1a\x1b.payload.v1.Search.Response\"\x18\x82\xd3\xe4\x93\x02\x12:\x01*\"\r/linearsearch\x12j\n" +
	"\x10LinearSearchByID\x12\x1c.payload.v1.Search.IDRequest\x1a\x1b.payload.v1.Search.Response\"\x1b\x82\xd3\xe4\x93\x02\x15:\x01*\"\x10/linearsearch/id\x12Y\n" +
	"\x12StreamLinearSearch\x12\x1a.payload.v1.Search.Request\x1a!.payload.v1.Search.StreamResponse\"\x00(\x010\x01\x12_\n" +
	"\x16StreamLinearSearchByID\x12\x1c.payload.v1.Search.IDRequest\x1a!.payload.v1.Search.StreamResponse\"\x00(\x010\x01\x12u\n" +
	"\x11MultiLinearSearch\x12\x1f.payload.v1.Search.MultiRequest\x1a\x1c.payload.v1.Search.Responses\"!\x82\xd3\xe4\x93\x02\x1b:\x01*\"\x16/linearsearch/multiple\x12~\n" +
	"\x15MultiLinearSearchByID\x12!.payload.v1.Search.MultiIDRequest\x1a\x1c.payload.v1.Search.Responses\"$\x82\xd3\xe4\x93\x02\x1e:\x01*\"\x19/linearsearch/id/multipleBS\n" +
	"\x1aorg.vdaas.vald.api.v1.valdB\n" +
	"ValdSearchP\x01Z'github.com/vdaas/vald/apis/grpc/v1/valdb\x06proto3"

var file_v1_vald_search_proto_goTypes = []any{
	(*payload.Search_Request)(nil),        // 0: payload.v1.Search.Request
	(*payload.Search_IDRequest)(nil),      // 1: payload.v1.Search.IDRequest
	(*payload.Search_MultiRequest)(nil),   // 2: payload.v1.Search.MultiRequest
	(*payload.Search_MultiIDRequest)(nil), // 3: payload.v1.Search.MultiIDRequest
	(*payload.Search_Response)(nil),       // 4: payload.v1.Search.Response
	(*payload.Search_StreamResponse)(nil), // 5: payload.v1.Search.StreamResponse
	(*payload.Search_Responses)(nil),      // 6: payload.v1.Search.Responses
}

var file_v1_vald_search_proto_depIdxs = []int32{
	0,  // 0: vald.v1.Search.Search:input_type -> payload.v1.Search.Request
	1,  // 1: vald.v1.Search.SearchByID:input_type -> payload.v1.Search.IDRequest
	0,  // 2: vald.v1.Search.StreamSearch:input_type -> payload.v1.Search.Request
	1,  // 3: vald.v1.Search.StreamSearchByID:input_type -> payload.v1.Search.IDRequest
	2,  // 4: vald.v1.Search.MultiSearch:input_type -> payload.v1.Search.MultiRequest
	3,  // 5: vald.v1.Search.MultiSearchByID:input_type -> payload.v1.Search.MultiIDRequest
	0,  // 6: vald.v1.Search.LinearSearch:input_type -> payload.v1.Search.Request
	1,  // 7: vald.v1.Search.LinearSearchByID:input_type -> payload.v1.Search.IDRequest
	0,  // 8: vald.v1.Search.StreamLinearSearch:input_type -> payload.v1.Search.Request
	1,  // 9: vald.v1.Search.StreamLinearSearchByID:input_type -> payload.v1.Search.IDRequest
	2,  // 10: vald.v1.Search.MultiLinearSearch:input_type -> payload.v1.Search.MultiRequest
	3,  // 11: vald.v1.Search.MultiLinearSearchByID:input_type -> payload.v1.Search.MultiIDRequest
	4,  // 12: vald.v1.Search.Search:output_type -> payload.v1.Search.Response
	4,  // 13: vald.v1.Search.SearchByID:output_type -> payload.v1.Search.Response
	5,  // 14: vald.v1.Search.StreamSearch:output_type -> payload.v1.Search.StreamResponse
	5,  // 15: vald.v1.Search.StreamSearchByID:output_type -> payload.v1.Search.StreamResponse
	6,  // 16: vald.v1.Search.MultiSearch:output_type -> payload.v1.Search.Responses
	6,  // 17: vald.v1.Search.MultiSearchByID:output_type -> payload.v1.Search.Responses
	4,  // 18: vald.v1.Search.LinearSearch:output_type -> payload.v1.Search.Response
	4,  // 19: vald.v1.Search.LinearSearchByID:output_type -> payload.v1.Search.Response
	5,  // 20: vald.v1.Search.StreamLinearSearch:output_type -> payload.v1.Search.StreamResponse
	5,  // 21: vald.v1.Search.StreamLinearSearchByID:output_type -> payload.v1.Search.StreamResponse
	6,  // 22: vald.v1.Search.MultiLinearSearch:output_type -> payload.v1.Search.Responses
	6,  // 23: vald.v1.Search.MultiLinearSearchByID:output_type -> payload.v1.Search.Responses
	12, // [12:24] is the sub-list for method output_type
	0,  // [0:12] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_v1_vald_search_proto_init() }
func file_v1_vald_search_proto_init() {
	if File_v1_vald_search_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_vald_search_proto_rawDesc), len(file_v1_vald_search_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_vald_search_proto_goTypes,
		DependencyIndexes: file_v1_vald_search_proto_depIdxs,
	}.Build()
	File_v1_vald_search_proto = out.File
	file_v1_vald_search_proto_goTypes = nil
	file_v1_vald_search_proto_depIdxs = nil
}
