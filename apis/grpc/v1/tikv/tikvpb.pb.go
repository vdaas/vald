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

package tikv

import (
	reflect "reflect"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_v1_tikv_tikvpb_proto protoreflect.FileDescriptor

const file_v1_tikv_tikvpb_proto_rawDesc = "" +
	"\n" +
	"\x14v1/tikv/tikvpb.proto\x12\x06tikvpb\x1a\x15v1/tikv/kvrpcpb.proto2\xb3\x03\n" +
	"\x04Tikv\x12;\n" +
	"\x06RawGet\x12\x16.kvrpcpb.RawGetRequest\x1a\x17.kvrpcpb.RawGetResponse\"\x00\x12J\n" +
	"\vRawBatchGet\x12\x1b.kvrpcpb.RawBatchGetRequest\x1a\x1c.kvrpcpb.RawBatchGetResponse\"\x00\x12;\n" +
	"\x06RawPut\x12\x16.kvrpcpb.RawPutRequest\x1a\x17.kvrpcpb.RawPutResponse\"\x00\x12J\n" +
	"\vRawBatchPut\x12\x1b.kvrpcpb.RawBatchPutRequest\x1a\x1c.kvrpcpb.RawBatchPutResponse\"\x00\x12D\n" +
	"\tRawDelete\x12\x19.kvrpcpb.RawDeleteRequest\x1a\x1a.kvrpcpb.RawDeleteResponse\"\x00\x12S\n" +
	"\x0eRawBatchDelete\x12\x1e.kvrpcpb.RawBatchDeleteRequest\x1a\x1f.kvrpcpb.RawBatchDeleteResponse\"\x00B)Z'github.com/vdaas/vald/apis/grpc/v1/tikvb\x06proto3"

var file_v1_tikv_tikvpb_proto_goTypes = []any{
	(*RawGetRequest)(nil),          // 0: kvrpcpb.RawGetRequest
	(*RawBatchGetRequest)(nil),     // 1: kvrpcpb.RawBatchGetRequest
	(*RawPutRequest)(nil),          // 2: kvrpcpb.RawPutRequest
	(*RawBatchPutRequest)(nil),     // 3: kvrpcpb.RawBatchPutRequest
	(*RawDeleteRequest)(nil),       // 4: kvrpcpb.RawDeleteRequest
	(*RawBatchDeleteRequest)(nil),  // 5: kvrpcpb.RawBatchDeleteRequest
	(*RawGetResponse)(nil),         // 6: kvrpcpb.RawGetResponse
	(*RawBatchGetResponse)(nil),    // 7: kvrpcpb.RawBatchGetResponse
	(*RawPutResponse)(nil),         // 8: kvrpcpb.RawPutResponse
	(*RawBatchPutResponse)(nil),    // 9: kvrpcpb.RawBatchPutResponse
	(*RawDeleteResponse)(nil),      // 10: kvrpcpb.RawDeleteResponse
	(*RawBatchDeleteResponse)(nil), // 11: kvrpcpb.RawBatchDeleteResponse
}
var file_v1_tikv_tikvpb_proto_depIdxs = []int32{
	0,  // 0: tikvpb.Tikv.RawGet:input_type -> kvrpcpb.RawGetRequest
	1,  // 1: tikvpb.Tikv.RawBatchGet:input_type -> kvrpcpb.RawBatchGetRequest
	2,  // 2: tikvpb.Tikv.RawPut:input_type -> kvrpcpb.RawPutRequest
	3,  // 3: tikvpb.Tikv.RawBatchPut:input_type -> kvrpcpb.RawBatchPutRequest
	4,  // 4: tikvpb.Tikv.RawDelete:input_type -> kvrpcpb.RawDeleteRequest
	5,  // 5: tikvpb.Tikv.RawBatchDelete:input_type -> kvrpcpb.RawBatchDeleteRequest
	6,  // 6: tikvpb.Tikv.RawGet:output_type -> kvrpcpb.RawGetResponse
	7,  // 7: tikvpb.Tikv.RawBatchGet:output_type -> kvrpcpb.RawBatchGetResponse
	8,  // 8: tikvpb.Tikv.RawPut:output_type -> kvrpcpb.RawPutResponse
	9,  // 9: tikvpb.Tikv.RawBatchPut:output_type -> kvrpcpb.RawBatchPutResponse
	10, // 10: tikvpb.Tikv.RawDelete:output_type -> kvrpcpb.RawDeleteResponse
	11, // 11: tikvpb.Tikv.RawBatchDelete:output_type -> kvrpcpb.RawBatchDeleteResponse
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_v1_tikv_tikvpb_proto_init() }
func file_v1_tikv_tikvpb_proto_init() {
	if File_v1_tikv_tikvpb_proto != nil {
		return
	}
	file_v1_tikv_kvrpcpb_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_tikv_tikvpb_proto_rawDesc), len(file_v1_tikv_tikvpb_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_tikv_tikvpb_proto_goTypes,
		DependencyIndexes: file_v1_tikv_tikvpb_proto_depIdxs,
	}.Build()
	File_v1_tikv_tikvpb_proto = out.File
	file_v1_tikv_tikvpb_proto_goTypes = nil
	file_v1_tikv_tikvpb_proto_depIdxs = nil
}
