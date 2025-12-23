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

	sync "github.com/vdaas/vald/internal/sync"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ErrorType int32

const (
	ErrorType_OK                      ErrorType = 0
	ErrorType_UNKNOWN                 ErrorType = 1
	ErrorType_NOT_BOOTSTRAPPED        ErrorType = 2
	ErrorType_STORE_TOMBSTONE         ErrorType = 3
	ErrorType_ALREADY_BOOTSTRAPPED    ErrorType = 4
	ErrorType_INCOMPATIBLE_VERSION    ErrorType = 5
	ErrorType_REGION_NOT_FOUND        ErrorType = 6
	ErrorType_GLOBAL_CONFIG_NOT_FOUND ErrorType = 7
	ErrorType_DUPLICATED_ENTRY        ErrorType = 8
	ErrorType_ENTRY_NOT_FOUND         ErrorType = 9
	ErrorType_INVALID_VALUE           ErrorType = 10
	// required watch revision is smaller than current compact/min revision.
	ErrorType_DATA_COMPACTED                    ErrorType = 11
	ErrorType_REGIONS_NOT_CONTAIN_ALL_KEY_RANGE ErrorType = 12
)

// Enum value maps for ErrorType.
var (
	ErrorType_name = map[int32]string{
		0:  "OK",
		1:  "UNKNOWN",
		2:  "NOT_BOOTSTRAPPED",
		3:  "STORE_TOMBSTONE",
		4:  "ALREADY_BOOTSTRAPPED",
		5:  "INCOMPATIBLE_VERSION",
		6:  "REGION_NOT_FOUND",
		7:  "GLOBAL_CONFIG_NOT_FOUND",
		8:  "DUPLICATED_ENTRY",
		9:  "ENTRY_NOT_FOUND",
		10: "INVALID_VALUE",
		11: "DATA_COMPACTED",
		12: "REGIONS_NOT_CONTAIN_ALL_KEY_RANGE",
	}
	ErrorType_value = map[string]int32{
		"OK":                                0,
		"UNKNOWN":                           1,
		"NOT_BOOTSTRAPPED":                  2,
		"STORE_TOMBSTONE":                   3,
		"ALREADY_BOOTSTRAPPED":              4,
		"INCOMPATIBLE_VERSION":              5,
		"REGION_NOT_FOUND":                  6,
		"GLOBAL_CONFIG_NOT_FOUND":           7,
		"DUPLICATED_ENTRY":                  8,
		"ENTRY_NOT_FOUND":                   9,
		"INVALID_VALUE":                     10,
		"DATA_COMPACTED":                    11,
		"REGIONS_NOT_CONTAIN_ALL_KEY_RANGE": 12,
	}
)

func (x ErrorType) Enum() *ErrorType {
	p := new(ErrorType)
	*p = x
	return p
}

func (x ErrorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorType) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_tikv_pdpb_proto_enumTypes[0].Descriptor()
}

func (ErrorType) Type() protoreflect.EnumType {
	return &file_v1_tikv_pdpb_proto_enumTypes[0]
}

func (x ErrorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorType.Descriptor instead.
func (ErrorType) EnumDescriptor() ([]byte, []int) {
	return file_v1_tikv_pdpb_proto_rawDescGZIP(), []int{0}
}

type Error struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          ErrorType              `                   protobuf:"varint,1,opt,name=type,proto3,enum=pdpb.ErrorType" json:"type,omitempty"`
	Message       string                 `                   protobuf:"bytes,2,opt,name=message,proto3"                   json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Error) Reset() {
	*x = Error{}
	mi := &file_v1_tikv_pdpb_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_pdpb_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_v1_tikv_pdpb_proto_rawDescGZIP(), []int{0}
}

func (x *Error) GetType() ErrorType {
	if x != nil {
		return x.Type
	}
	return ErrorType_OK
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GetAllStoresRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// RequestHeader header = 1;
	// Do NOT return tombstone stores if set to true.
	ExcludeTombstoneStores bool `                   protobuf:"varint,2,opt,name=exclude_tombstone_stores,json=excludeTombstoneStores,proto3" json:"exclude_tombstone_stores,omitempty"`
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *GetAllStoresRequest) Reset() {
	*x = GetAllStoresRequest{}
	mi := &file_v1_tikv_pdpb_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAllStoresRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllStoresRequest) ProtoMessage() {}

func (x *GetAllStoresRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_pdpb_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllStoresRequest.ProtoReflect.Descriptor instead.
func (*GetAllStoresRequest) Descriptor() ([]byte, []int) {
	return file_v1_tikv_pdpb_proto_rawDescGZIP(), []int{1}
}

func (x *GetAllStoresRequest) GetExcludeTombstoneStores() bool {
	if x != nil {
		return x.ExcludeTombstoneStores
	}
	return false
}

type GetAllStoresResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Stores        []*Store               `                   protobuf:"bytes,2,rep,name=stores,proto3" json:"stores,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAllStoresResponse) Reset() {
	*x = GetAllStoresResponse{}
	mi := &file_v1_tikv_pdpb_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAllStoresResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllStoresResponse) ProtoMessage() {}

func (x *GetAllStoresResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_pdpb_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllStoresResponse.ProtoReflect.Descriptor instead.
func (*GetAllStoresResponse) Descriptor() ([]byte, []int) {
	return file_v1_tikv_pdpb_proto_rawDescGZIP(), []int{2}
}

func (x *GetAllStoresResponse) GetStores() []*Store {
	if x != nil {
		return x.Stores
	}
	return nil
}

type Region struct {
	state  protoimpl.MessageState `protogen:"open.v1"`
	Region *Region                `                   protobuf:"bytes,1,opt,name=region,proto3"                          json:"region,omitempty"`
	Leader *Peer                  `                   protobuf:"bytes,2,opt,name=leader,proto3"                          json:"leader,omitempty"`
	// Leader considers that these peers are down.
	// repeated PeerStats down_peers = 3;
	// Pending peers are the peers that the leader can't consider as
	// working followers.
	PendingPeers  []*Peer `                   protobuf:"bytes,4,rep,name=pending_peers,json=pendingPeers,proto3" json:"pending_peers,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Region) Reset() {
	*x = Region{}
	mi := &file_v1_tikv_pdpb_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Region) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Region) ProtoMessage() {}

func (x *Region) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_pdpb_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Region.ProtoReflect.Descriptor instead.
func (*Region) Descriptor() ([]byte, []int) {
	return file_v1_tikv_pdpb_proto_rawDescGZIP(), []int{3}
}

func (x *Region) GetRegion() *Region {
	if x != nil {
		return x.Region
	}
	return nil
}

func (x *Region) GetLeader() *Peer {
	if x != nil {
		return x.Leader
	}
	return nil
}

func (x *Region) GetPendingPeers() []*Peer {
	if x != nil {
		return x.PendingPeers
	}
	return nil
}

type KeyRange struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	StartKey      []byte                 `                   protobuf:"bytes,1,opt,name=start_key,json=startKey,proto3" json:"start_key,omitempty"`
	EndKey        []byte                 `                   protobuf:"bytes,2,opt,name=end_key,json=endKey,proto3"     json:"end_key,omitempty"` // end_key is +inf when it is empty.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KeyRange) Reset() {
	*x = KeyRange{}
	mi := &file_v1_tikv_pdpb_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KeyRange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyRange) ProtoMessage() {}

func (x *KeyRange) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_pdpb_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyRange.ProtoReflect.Descriptor instead.
func (*KeyRange) Descriptor() ([]byte, []int) {
	return file_v1_tikv_pdpb_proto_rawDescGZIP(), []int{4}
}

func (x *KeyRange) GetStartKey() []byte {
	if x != nil {
		return x.StartKey
	}
	return nil
}

func (x *KeyRange) GetEndKey() []byte {
	if x != nil {
		return x.EndKey
	}
	return nil
}

type BatchScanRegionsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// RequestHeader header = 1;
	NeedBuckets bool        `                   protobuf:"varint,2,opt,name=need_buckets,json=needBuckets,proto3"                 json:"need_buckets,omitempty"`
	Ranges      []*KeyRange `                   protobuf:"bytes,3,rep,name=ranges,proto3"                                         json:"ranges,omitempty"` // the given ranges must be in order.
	Limit       int32       `                   protobuf:"varint,4,opt,name=limit,proto3"                                         json:"limit,omitempty"`  // limit the total number of regions to scan.
	// If contain_all_key_range is true, the output must contain all
	// key ranges in the request.
	// If the output does not contain all key ranges, the request is considered
	// failed and returns an error(REGIONS_NOT_CONTAIN_ALL_KEY_RANGE).
	ContainAllKeyRange bool `                   protobuf:"varint,5,opt,name=contain_all_key_range,json=containAllKeyRange,proto3" json:"contain_all_key_range,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *BatchScanRegionsRequest) Reset() {
	*x = BatchScanRegionsRequest{}
	mi := &file_v1_tikv_pdpb_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BatchScanRegionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchScanRegionsRequest) ProtoMessage() {}

func (x *BatchScanRegionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_pdpb_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchScanRegionsRequest.ProtoReflect.Descriptor instead.
func (*BatchScanRegionsRequest) Descriptor() ([]byte, []int) {
	return file_v1_tikv_pdpb_proto_rawDescGZIP(), []int{5}
}

func (x *BatchScanRegionsRequest) GetNeedBuckets() bool {
	if x != nil {
		return x.NeedBuckets
	}
	return false
}

func (x *BatchScanRegionsRequest) GetRanges() []*KeyRange {
	if x != nil {
		return x.Ranges
	}
	return nil
}

func (x *BatchScanRegionsRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *BatchScanRegionsRequest) GetContainAllKeyRange() bool {
	if x != nil {
		return x.ContainAllKeyRange
	}
	return false
}

type BatchScanRegionsResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// the returned regions are flattened into a list, because the given ranges can located in the same range, we do not return duplicated regions then.
	Regions       []*Region `                   protobuf:"bytes,2,rep,name=regions,proto3" json:"regions,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BatchScanRegionsResponse) Reset() {
	*x = BatchScanRegionsResponse{}
	mi := &file_v1_tikv_pdpb_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BatchScanRegionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchScanRegionsResponse) ProtoMessage() {}

func (x *BatchScanRegionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_pdpb_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchScanRegionsResponse.ProtoReflect.Descriptor instead.
func (*BatchScanRegionsResponse) Descriptor() ([]byte, []int) {
	return file_v1_tikv_pdpb_proto_rawDescGZIP(), []int{6}
}

func (x *BatchScanRegionsResponse) GetRegions() []*Region {
	if x != nil {
		return x.Regions
	}
	return nil
}

var File_v1_tikv_pdpb_proto protoreflect.FileDescriptor

const file_v1_tikv_pdpb_proto_rawDesc = "" +
	"\n" +
	"\x12v1/tikv/pdpb.proto\x12\x04pdpb\x1a\x14v1/tikv/metapb.proto\"F\n" +
	"\x05Error\x12#\n" +
	"\x04type\x18\x01 \x01(\x0e2\x0f.pdpb.ErrorTypeR\x04type\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\"O\n" +
	"\x13GetAllStoresRequest\x128\n" +
	"\x18exclude_tombstone_stores\x18\x02 \x01(\bR\x16excludeTombstoneStores\"=\n" +
	"\x14GetAllStoresResponse\x12%\n" +
	"\x06stores\x18\x02 \x03(\v2\r.metapb.StoreR\x06stores\"\x89\x01\n" +
	"\x06Region\x12&\n" +
	"\x06region\x18\x01 \x01(\v2\x0e.metapb.RegionR\x06region\x12$\n" +
	"\x06leader\x18\x02 \x01(\v2\f.metapb.PeerR\x06leader\x121\n" +
	"\rpending_peers\x18\x04 \x03(\v2\f.metapb.PeerR\fpendingPeers\"@\n" +
	"\bKeyRange\x12\x1b\n" +
	"\tstart_key\x18\x01 \x01(\fR\bstartKey\x12\x17\n" +
	"\aend_key\x18\x02 \x01(\fR\x06endKey\"\xad\x01\n" +
	"\x17BatchScanRegionsRequest\x12!\n" +
	"\fneed_buckets\x18\x02 \x01(\bR\vneedBuckets\x12&\n" +
	"\x06ranges\x18\x03 \x03(\v2\x0e.pdpb.KeyRangeR\x06ranges\x12\x14\n" +
	"\x05limit\x18\x04 \x01(\x05R\x05limit\x121\n" +
	"\x15contain_all_key_range\x18\x05 \x01(\bR\x12containAllKeyRange\"B\n" +
	"\x18BatchScanRegionsResponse\x12&\n" +
	"\aregions\x18\x02 \x03(\v2\f.pdpb.RegionR\aregions*\xab\x02\n" +
	"\tErrorType\x12\x06\n" +
	"\x02OK\x10\x00\x12\v\n" +
	"\aUNKNOWN\x10\x01\x12\x14\n" +
	"\x10NOT_BOOTSTRAPPED\x10\x02\x12\x13\n" +
	"\x0fSTORE_TOMBSTONE\x10\x03\x12\x18\n" +
	"\x14ALREADY_BOOTSTRAPPED\x10\x04\x12\x18\n" +
	"\x14INCOMPATIBLE_VERSION\x10\x05\x12\x14\n" +
	"\x10REGION_NOT_FOUND\x10\x06\x12\x1b\n" +
	"\x17GLOBAL_CONFIG_NOT_FOUND\x10\a\x12\x14\n" +
	"\x10DUPLICATED_ENTRY\x10\b\x12\x13\n" +
	"\x0fENTRY_NOT_FOUND\x10\t\x12\x11\n" +
	"\rINVALID_VALUE\x10\n" +
	"\x12\x12\n" +
	"\x0eDATA_COMPACTED\x10\v\x12%\n" +
	"!REGIONS_NOT_CONTAIN_ALL_KEY_RANGE\x10\f2\xa2\x01\n" +
	"\x02PD\x12G\n" +
	"\fGetAllStores\x12\x19.pdpb.GetAllStoresRequest\x1a\x1a.pdpb.GetAllStoresResponse\"\x00\x12S\n" +
	"\x10BatchScanRegions\x12\x1d.pdpb.BatchScanRegionsRequest\x1a\x1e.pdpb.BatchScanRegionsResponse\"\x00B)Z'github.com/vdaas/vald/apis/grpc/v1/tikvb\x06proto3"

var (
	file_v1_tikv_pdpb_proto_rawDescOnce sync.Once
	file_v1_tikv_pdpb_proto_rawDescData []byte
)

func file_v1_tikv_pdpb_proto_rawDescGZIP() []byte {
	file_v1_tikv_pdpb_proto_rawDescOnce.Do(func() {
		file_v1_tikv_pdpb_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_tikv_pdpb_proto_rawDesc), len(file_v1_tikv_pdpb_proto_rawDesc)))
	})
	return file_v1_tikv_pdpb_proto_rawDescData
}

var (
	file_v1_tikv_pdpb_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
	file_v1_tikv_pdpb_proto_msgTypes  = make([]protoimpl.MessageInfo, 7)
	file_v1_tikv_pdpb_proto_goTypes   = []any{
		(ErrorType)(0),                   // 0: pdpb.ErrorType
		(*Error)(nil),                    // 1: pdpb.Error
		(*GetAllStoresRequest)(nil),      // 2: pdpb.GetAllStoresRequest
		(*GetAllStoresResponse)(nil),     // 3: pdpb.GetAllStoresResponse
		(*Region)(nil),                   // 4: pdpb.Region
		(*KeyRange)(nil),                 // 5: pdpb.KeyRange
		(*BatchScanRegionsRequest)(nil),  // 6: pdpb.BatchScanRegionsRequest
		(*BatchScanRegionsResponse)(nil), // 7: pdpb.BatchScanRegionsResponse
		(*Store)(nil),                    // 8: metapb.Store
		(*Region)(nil),                   // 9: metapb.Region
		(*Peer)(nil),                     // 10: metapb.Peer
	}
)
var file_v1_tikv_pdpb_proto_depIdxs = []int32{
	0,  // 0: pdpb.Error.type:type_name -> pdpb.ErrorType
	8,  // 1: pdpb.GetAllStoresResponse.stores:type_name -> metapb.Store
	9,  // 2: pdpb.Region.region:type_name -> metapb.Region
	10, // 3: pdpb.Region.leader:type_name -> metapb.Peer
	10, // 4: pdpb.Region.pending_peers:type_name -> metapb.Peer
	5,  // 5: pdpb.BatchScanRegionsRequest.ranges:type_name -> pdpb.KeyRange
	4,  // 6: pdpb.BatchScanRegionsResponse.regions:type_name -> pdpb.Region
	2,  // 7: pdpb.PD.GetAllStores:input_type -> pdpb.GetAllStoresRequest
	6,  // 8: pdpb.PD.BatchScanRegions:input_type -> pdpb.BatchScanRegionsRequest
	3,  // 9: pdpb.PD.GetAllStores:output_type -> pdpb.GetAllStoresResponse
	7,  // 10: pdpb.PD.BatchScanRegions:output_type -> pdpb.BatchScanRegionsResponse
	9,  // [9:11] is the sub-list for method output_type
	7,  // [7:9] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_v1_tikv_pdpb_proto_init() }
func file_v1_tikv_pdpb_proto_init() {
	if File_v1_tikv_pdpb_proto != nil {
		return
	}
	file_v1_tikv_metapb_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_tikv_pdpb_proto_rawDesc), len(file_v1_tikv_pdpb_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_v1_tikv_pdpb_proto_goTypes,
		DependencyIndexes: file_v1_tikv_pdpb_proto_depIdxs,
		EnumInfos:         file_v1_tikv_pdpb_proto_enumTypes,
		MessageInfos:      file_v1_tikv_pdpb_proto_msgTypes,
	}.Build()
	File_v1_tikv_pdpb_proto = out.File
	file_v1_tikv_pdpb_proto_goTypes = nil
	file_v1_tikv_pdpb_proto_depIdxs = nil
}
