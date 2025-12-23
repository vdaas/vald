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

type StoreState int32

const (
	StoreState_Up        StoreState = 0
	StoreState_Offline   StoreState = 1
	StoreState_Tombstone StoreState = 2
)

// Enum value maps for StoreState.
var (
	StoreState_name = map[int32]string{
		0: "Up",
		1: "Offline",
		2: "Tombstone",
	}
	StoreState_value = map[string]int32{
		"Up":        0,
		"Offline":   1,
		"Tombstone": 2,
	}
)

func (x StoreState) Enum() *StoreState {
	p := new(StoreState)
	*p = x
	return p
}

func (x StoreState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StoreState) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_tikv_metapb_proto_enumTypes[0].Descriptor()
}

func (StoreState) Type() protoreflect.EnumType {
	return &file_v1_tikv_metapb_proto_enumTypes[0]
}

func (x StoreState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StoreState.Descriptor instead.
func (StoreState) EnumDescriptor() ([]byte, []int) {
	return file_v1_tikv_metapb_proto_rawDescGZIP(), []int{0}
}

// Copied from metapb.proto
type Peer struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	Id      uint64                 `                   protobuf:"varint,1,opt,name=id,proto3"                        json:"id,omitempty"`
	StoreId uint64                 `                   protobuf:"varint,2,opt,name=store_id,json=storeId,proto3"     json:"store_id,omitempty"`
	// PeerRole role = 3;
	IsWitness     bool `                   protobuf:"varint,4,opt,name=is_witness,json=isWitness,proto3" json:"is_witness,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Peer) Reset() {
	*x = Peer{}
	mi := &file_v1_tikv_metapb_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Peer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Peer) ProtoMessage() {}

func (x *Peer) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_metapb_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Peer.ProtoReflect.Descriptor instead.
func (*Peer) Descriptor() ([]byte, []int) {
	return file_v1_tikv_metapb_proto_rawDescGZIP(), []int{0}
}

func (x *Peer) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Peer) GetStoreId() uint64 {
	if x != nil {
		return x.StoreId
	}
	return 0
}

func (x *Peer) GetIsWitness() bool {
	if x != nil {
		return x.IsWitness
	}
	return false
}

type Store struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Id    uint64                 `                   protobuf:"varint,1,opt,name=id,proto3"                                             json:"id,omitempty"`
	// Address to handle client requests (kv, cop, etc.)
	Address string     `                   protobuf:"bytes,2,opt,name=address,proto3"                                         json:"address,omitempty"`
	State   StoreState `                   protobuf:"varint,3,opt,name=state,proto3,enum=metapb.StoreState"                   json:"state,omitempty"`
	// repeated StoreLabel labels = 4;
	Version string `                   protobuf:"bytes,5,opt,name=version,proto3"                                         json:"version,omitempty"`
	// Address to handle peer requests (raft messages from other store).
	// Empty means same as address.
	PeerAddress string `                   protobuf:"bytes,6,opt,name=peer_address,json=peerAddress,proto3"                   json:"peer_address,omitempty"`
	// Status address provides the HTTP service for external components
	StatusAddress string `                   protobuf:"bytes,7,opt,name=status_address,json=statusAddress,proto3"               json:"status_address,omitempty"`
	GitHash       string `                   protobuf:"bytes,8,opt,name=git_hash,json=gitHash,proto3"                           json:"git_hash,omitempty"`
	// The start timestamp of the current store
	StartTimestamp int64  `                   protobuf:"varint,9,opt,name=start_timestamp,json=startTimestamp,proto3"            json:"start_timestamp,omitempty"`
	DeployPath     string `                   protobuf:"bytes,10,opt,name=deploy_path,json=deployPath,proto3"                    json:"deploy_path,omitempty"`
	// The last heartbeat timestamp of the store.
	LastHeartbeat int64 `                   protobuf:"varint,11,opt,name=last_heartbeat,json=lastHeartbeat,proto3"             json:"last_heartbeat,omitempty"`
	// If the store is physically destroyed, which means it can never up again.
	PhysicallyDestroyed bool `                   protobuf:"varint,12,opt,name=physically_destroyed,json=physicallyDestroyed,proto3" json:"physically_destroyed,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *Store) Reset() {
	*x = Store{}
	mi := &file_v1_tikv_metapb_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Store) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Store) ProtoMessage() {}

func (x *Store) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_metapb_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Store.ProtoReflect.Descriptor instead.
func (*Store) Descriptor() ([]byte, []int) {
	return file_v1_tikv_metapb_proto_rawDescGZIP(), []int{1}
}

func (x *Store) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Store) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Store) GetState() StoreState {
	if x != nil {
		return x.State
	}
	return StoreState_Up
}

func (x *Store) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Store) GetPeerAddress() string {
	if x != nil {
		return x.PeerAddress
	}
	return ""
}

func (x *Store) GetStatusAddress() string {
	if x != nil {
		return x.StatusAddress
	}
	return ""
}

func (x *Store) GetGitHash() string {
	if x != nil {
		return x.GitHash
	}
	return ""
}

func (x *Store) GetStartTimestamp() int64 {
	if x != nil {
		return x.StartTimestamp
	}
	return 0
}

func (x *Store) GetDeployPath() string {
	if x != nil {
		return x.DeployPath
	}
	return ""
}

func (x *Store) GetLastHeartbeat() int64 {
	if x != nil {
		return x.LastHeartbeat
	}
	return 0
}

func (x *Store) GetPhysicallyDestroyed() bool {
	if x != nil {
		return x.PhysicallyDestroyed
	}
	return false
}

type Region2 struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Id    uint64                 `                   protobuf:"varint,1,opt,name=id,proto3"                                       json:"id,omitempty"`
	// Region key range [start_key, end_key).
	StartKey []byte `                   protobuf:"bytes,2,opt,name=start_key,json=startKey,proto3"                   json:"start_key,omitempty"`
	EndKey   []byte `                   protobuf:"bytes,3,opt,name=end_key,json=endKey,proto3"                       json:"end_key,omitempty"`
	// RegionEpoch region_epoch = 4;
	Peers []*Peer `                   protobuf:"bytes,5,rep,name=peers,proto3"                                     json:"peers,omitempty"`
	// Encryption metadata for start_key and end_key. encryption_meta.iv is IV for start_key.
	// IV for end_key is calculated from (encryption_meta.iv + len(start_key)).
	// The field is only used by PD and should be ignored otherwise.
	// If encryption_meta is empty (i.e. nil), it means start_key and end_key are unencrypted.
	// encryptionpb.EncryptionMeta encryption_meta = 6;
	// The flashback state indicates whether this region is in the flashback state.
	// TODO: only check by `flashback_start_ts` in the future. Keep for compatibility now.
	IsInFlashback bool `                   protobuf:"varint,7,opt,name=is_in_flashback,json=isInFlashback,proto3"       json:"is_in_flashback,omitempty"`
	// The start_ts that the current flashback progress is using.
	FlashbackStartTs uint64 `                   protobuf:"varint,8,opt,name=flashback_start_ts,json=flashbackStartTs,proto3" json:"flashback_start_ts,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *Region2) Reset() {
	*x = Region2{}
	mi := &file_v1_tikv_metapb_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Region2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Region2) ProtoMessage() {}

func (x *Region2) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_metapb_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Region2.ProtoReflect.Descriptor instead.
func (*Region2) Descriptor() ([]byte, []int) {
	return file_v1_tikv_metapb_proto_rawDescGZIP(), []int{2}
}

func (x *Region2) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Region2) GetStartKey() []byte {
	if x != nil {
		return x.StartKey
	}
	return nil
}

func (x *Region2) GetEndKey() []byte {
	if x != nil {
		return x.EndKey
	}
	return nil
}

func (x *Region2) GetPeers() []*Peer {
	if x != nil {
		return x.Peers
	}
	return nil
}

func (x *Region2) GetIsInFlashback() bool {
	if x != nil {
		return x.IsInFlashback
	}
	return false
}

func (x *Region2) GetFlashbackStartTs() uint64 {
	if x != nil {
		return x.FlashbackStartTs
	}
	return 0
}

var File_v1_tikv_metapb_proto protoreflect.FileDescriptor

const file_v1_tikv_metapb_proto_rawDesc = "" +
	"\n" +
	"\x14v1/tikv/metapb.proto\x12\x06metapb\"P\n" +
	"\x04Peer\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x19\n" +
	"\bstore_id\x18\x02 \x01(\x04R\astoreId\x12\x1d\n" +
	"\n" +
	"is_witness\x18\x04 \x01(\bR\tisWitness\"\xfe\x02\n" +
	"\x05Store\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x18\n" +
	"\aaddress\x18\x02 \x01(\tR\aaddress\x12(\n" +
	"\x05state\x18\x03 \x01(\x0e2\x12.metapb.StoreStateR\x05state\x12\x18\n" +
	"\aversion\x18\x05 \x01(\tR\aversion\x12!\n" +
	"\fpeer_address\x18\x06 \x01(\tR\vpeerAddress\x12%\n" +
	"\x0estatus_address\x18\a \x01(\tR\rstatusAddress\x12\x19\n" +
	"\bgit_hash\x18\b \x01(\tR\agitHash\x12'\n" +
	"\x0fstart_timestamp\x18\t \x01(\x03R\x0estartTimestamp\x12\x1f\n" +
	"\vdeploy_path\x18\n" +
	" \x01(\tR\n" +
	"deployPath\x12%\n" +
	"\x0elast_heartbeat\x18\v \x01(\x03R\rlastHeartbeat\x121\n" +
	"\x14physically_destroyed\x18\f \x01(\bR\x13physicallyDestroyed\"\xc9\x01\n" +
	"\aRegion2\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x1b\n" +
	"\tstart_key\x18\x02 \x01(\fR\bstartKey\x12\x17\n" +
	"\aend_key\x18\x03 \x01(\fR\x06endKey\x12\"\n" +
	"\x05peers\x18\x05 \x03(\v2\f.metapb.PeerR\x05peers\x12&\n" +
	"\x0fis_in_flashback\x18\a \x01(\bR\risInFlashback\x12,\n" +
	"\x12flashback_start_ts\x18\b \x01(\x04R\x10flashbackStartTs*0\n" +
	"\n" +
	"StoreState\x12\x06\n" +
	"\x02Up\x10\x00\x12\v\n" +
	"\aOffline\x10\x01\x12\r\n" +
	"\tTombstone\x10\x02B)Z'github.com/vdaas/vald/apis/grpc/v1/tikvb\x06proto3"

var (
	file_v1_tikv_metapb_proto_rawDescOnce sync.Once
	file_v1_tikv_metapb_proto_rawDescData []byte
)

func file_v1_tikv_metapb_proto_rawDescGZIP() []byte {
	file_v1_tikv_metapb_proto_rawDescOnce.Do(func() {
		file_v1_tikv_metapb_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_tikv_metapb_proto_rawDesc), len(file_v1_tikv_metapb_proto_rawDesc)))
	})
	return file_v1_tikv_metapb_proto_rawDescData
}

var (
	file_v1_tikv_metapb_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
	file_v1_tikv_metapb_proto_msgTypes  = make([]protoimpl.MessageInfo, 3)
	file_v1_tikv_metapb_proto_goTypes   = []any{
		(StoreState)(0), // 0: metapb.StoreState
		(*Peer)(nil),    // 1: metapb.Peer
		(*Store)(nil),   // 2: metapb.Store
		(*Region2)(nil), // 3: metapb.Region2
	}
)
var file_v1_tikv_metapb_proto_depIdxs = []int32{
	0, // 0: metapb.Store.state:type_name -> metapb.StoreState
	1, // 1: metapb.Region2.peers:type_name -> metapb.Peer
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_v1_tikv_metapb_proto_init() }
func file_v1_tikv_metapb_proto_init() {
	if File_v1_tikv_metapb_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_tikv_metapb_proto_rawDesc), len(file_v1_tikv_metapb_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_tikv_metapb_proto_goTypes,
		DependencyIndexes: file_v1_tikv_metapb_proto_depIdxs,
		EnumInfos:         file_v1_tikv_metapb_proto_enumTypes,
		MessageInfos:      file_v1_tikv_metapb_proto_msgTypes,
	}.Build()
	File_v1_tikv_metapb_proto = out.File
	file_v1_tikv_metapb_proto_goTypes = nil
	file_v1_tikv_metapb_proto_depIdxs = nil
}
