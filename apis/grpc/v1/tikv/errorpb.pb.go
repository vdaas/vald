//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

// NotLeader is the error variant that tells a request be handle by raft leader
// is sent to raft follower or learner.
type NotLeader struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The requested region ID
	RegionId uint64 `                   protobuf:"varint,1,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`
	// Region leader of the requested region
	Leader        *Peer `                   protobuf:"bytes,2,opt,name=leader,proto3"                   json:"leader,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NotLeader) Reset() {
	*x = NotLeader{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NotLeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotLeader) ProtoMessage() {}

func (x *NotLeader) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotLeader.ProtoReflect.Descriptor instead.
func (*NotLeader) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{0}
}

func (x *NotLeader) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

func (x *NotLeader) GetLeader() *Peer {
	if x != nil {
		return x.Leader
	}
	return nil
}

// IsWitness is the error variant that tells a request be handle by witness
// which should be forbidden and retry.
type IsWitness struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The requested region ID
	RegionId      uint64 `                   protobuf:"varint,1,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *IsWitness) Reset() {
	*x = IsWitness{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IsWitness) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsWitness) ProtoMessage() {}

func (x *IsWitness) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsWitness.ProtoReflect.Descriptor instead.
func (*IsWitness) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{1}
}

func (x *IsWitness) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

// BucketVersionNotMatch is the error variant that tells the request buckets version is not match.
// client should update the buckets version and retry.
type BucketVersionNotMatch struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Version       uint64                 `                   protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Keys          [][]byte               `                   protobuf:"bytes,2,rep,name=keys,proto3"     json:"keys,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BucketVersionNotMatch) Reset() {
	*x = BucketVersionNotMatch{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BucketVersionNotMatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BucketVersionNotMatch) ProtoMessage() {}

func (x *BucketVersionNotMatch) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BucketVersionNotMatch.ProtoReflect.Descriptor instead.
func (*BucketVersionNotMatch) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{2}
}

func (x *BucketVersionNotMatch) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *BucketVersionNotMatch) GetKeys() [][]byte {
	if x != nil {
		return x.Keys
	}
	return nil
}

type DiskFull struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The requested store ID
	StoreId []uint64 `                   protobuf:"varint,1,rep,packed,name=store_id,json=storeId,proto3" json:"store_id,omitempty"`
	// The detailed info
	Reason        string `                   protobuf:"bytes,2,opt,name=reason,proto3"                        json:"reason,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DiskFull) Reset() {
	*x = DiskFull{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DiskFull) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiskFull) ProtoMessage() {}

func (x *DiskFull) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiskFull.ProtoReflect.Descriptor instead.
func (*DiskFull) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{3}
}

func (x *DiskFull) GetStoreId() []uint64 {
	if x != nil {
		return x.StoreId
	}
	return nil
}

func (x *DiskFull) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

// StoreNotMatch is the error variant that tells the request is sent to wrong store.
// (i.e. inconsistency of the store ID that request shows and the real store ID of this server.)
type StoreNotMatch struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Store id in request
	RequestStoreId uint64 `                   protobuf:"varint,1,opt,name=request_store_id,json=requestStoreId,proto3" json:"request_store_id,omitempty"`
	// Actual store id
	ActualStoreId uint64 `                   protobuf:"varint,2,opt,name=actual_store_id,json=actualStoreId,proto3"   json:"actual_store_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StoreNotMatch) Reset() {
	*x = StoreNotMatch{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StoreNotMatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreNotMatch) ProtoMessage() {}

func (x *StoreNotMatch) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreNotMatch.ProtoReflect.Descriptor instead.
func (*StoreNotMatch) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{4}
}

func (x *StoreNotMatch) GetRequestStoreId() uint64 {
	if x != nil {
		return x.RequestStoreId
	}
	return 0
}

func (x *StoreNotMatch) GetActualStoreId() uint64 {
	if x != nil {
		return x.ActualStoreId
	}
	return 0
}

// RegionNotFound is the error variant that tells there isn't any region in this TiKV
// matches the requested region ID.
type RegionNotFound struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The requested region ID
	RegionId      uint64 `                   protobuf:"varint,1,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegionNotFound) Reset() {
	*x = RegionNotFound{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegionNotFound) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegionNotFound) ProtoMessage() {}

func (x *RegionNotFound) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegionNotFound.ProtoReflect.Descriptor instead.
func (*RegionNotFound) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{5}
}

func (x *RegionNotFound) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

// RegionNotInitialized is the error variant that tells there isn't any initialized peer
// matchesthe request region ID.
type RegionNotInitialized struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The request region ID
	RegionId      uint64 `                   protobuf:"varint,1,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegionNotInitialized) Reset() {
	*x = RegionNotInitialized{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegionNotInitialized) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegionNotInitialized) ProtoMessage() {}

func (x *RegionNotInitialized) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegionNotInitialized.ProtoReflect.Descriptor instead.
func (*RegionNotInitialized) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{6}
}

func (x *RegionNotInitialized) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

// KeyNotInRegion is the error variant that tells the key the request requires isn't present in
// this region.
type KeyNotInRegion struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The requested key
	Key []byte `                   protobuf:"bytes,1,opt,name=key,proto3"                      json:"key,omitempty"`
	// The requested region ID
	RegionId uint64 `                   protobuf:"varint,2,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`
	// Start key of the requested region
	StartKey []byte `                   protobuf:"bytes,3,opt,name=start_key,json=startKey,proto3"  json:"start_key,omitempty"`
	// Snd key of the requested region
	EndKey        []byte `                   protobuf:"bytes,4,opt,name=end_key,json=endKey,proto3"      json:"end_key,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KeyNotInRegion) Reset() {
	*x = KeyNotInRegion{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KeyNotInRegion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyNotInRegion) ProtoMessage() {}

func (x *KeyNotInRegion) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyNotInRegion.ProtoReflect.Descriptor instead.
func (*KeyNotInRegion) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{7}
}

func (x *KeyNotInRegion) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *KeyNotInRegion) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

func (x *KeyNotInRegion) GetStartKey() []byte {
	if x != nil {
		return x.StartKey
	}
	return nil
}

func (x *KeyNotInRegion) GetEndKey() []byte {
	if x != nil {
		return x.EndKey
	}
	return nil
}

// EpochNotMatch is the error variant that tells a region has been updated.
// (e.g. by splitting / merging, or raft Confchange.)
// Hence, a command is based on a stale version of a region.
type EpochNotMatch struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Available regions that may be siblings of the requested one.
	CurrentRegions []*Region2 `                   protobuf:"bytes,1,rep,name=current_regions,json=currentRegions,proto3" json:"current_regions,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *EpochNotMatch) Reset() {
	*x = EpochNotMatch{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EpochNotMatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EpochNotMatch) ProtoMessage() {}

func (x *EpochNotMatch) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EpochNotMatch.ProtoReflect.Descriptor instead.
func (*EpochNotMatch) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{8}
}

func (x *EpochNotMatch) GetCurrentRegions() []*Region2 {
	if x != nil {
		return x.CurrentRegions
	}
	return nil
}

// ServerIsBusy is the error variant that tells the server is too busy to response.
type ServerIsBusy struct {
	state  protoimpl.MessageState `protogen:"open.v1"`
	Reason string                 `                   protobuf:"bytes,1,opt,name=reason,proto3"                                  json:"reason,omitempty"`
	// The suggested backoff time
	BackoffMs       uint64 `                   protobuf:"varint,2,opt,name=backoff_ms,json=backoffMs,proto3"              json:"backoff_ms,omitempty"`
	EstimatedWaitMs uint32 `                   protobuf:"varint,3,opt,name=estimated_wait_ms,json=estimatedWaitMs,proto3" json:"estimated_wait_ms,omitempty"`
	// Current applied_index at the leader, may be used in replica read.
	AppliedIndex  uint64 `                   protobuf:"varint,4,opt,name=applied_index,json=appliedIndex,proto3"        json:"applied_index,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ServerIsBusy) Reset() {
	*x = ServerIsBusy{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServerIsBusy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerIsBusy) ProtoMessage() {}

func (x *ServerIsBusy) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerIsBusy.ProtoReflect.Descriptor instead.
func (*ServerIsBusy) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{9}
}

func (x *ServerIsBusy) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *ServerIsBusy) GetBackoffMs() uint64 {
	if x != nil {
		return x.BackoffMs
	}
	return 0
}

func (x *ServerIsBusy) GetEstimatedWaitMs() uint32 {
	if x != nil {
		return x.EstimatedWaitMs
	}
	return 0
}

func (x *ServerIsBusy) GetAppliedIndex() uint64 {
	if x != nil {
		return x.AppliedIndex
	}
	return 0
}

// StaleCommand is the error variant that tells the command is stale, that is,
// the current request term is lower than current raft term.
// This can be retried at most time.
type StaleCommand struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StaleCommand) Reset() {
	*x = StaleCommand{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StaleCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StaleCommand) ProtoMessage() {}

func (x *StaleCommand) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StaleCommand.ProtoReflect.Descriptor instead.
func (*StaleCommand) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{10}
}

// RaftEntryTooLarge is the error variant that tells the request is too large to be serialized to a
// reasonable small raft entry.
// (i.e. greater than the configured value `raft_entry_max_size` in `raftstore`)
type RaftEntryTooLarge struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The requested region ID
	RegionId uint64 `                   protobuf:"varint,1,opt,name=region_id,json=regionId,proto3"   json:"region_id,omitempty"`
	// Size of the raft entry
	EntrySize     uint64 `                   protobuf:"varint,2,opt,name=entry_size,json=entrySize,proto3" json:"entry_size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RaftEntryTooLarge) Reset() {
	*x = RaftEntryTooLarge{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RaftEntryTooLarge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftEntryTooLarge) ProtoMessage() {}

func (x *RaftEntryTooLarge) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftEntryTooLarge.ProtoReflect.Descriptor instead.
func (*RaftEntryTooLarge) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{11}
}

func (x *RaftEntryTooLarge) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

func (x *RaftEntryTooLarge) GetEntrySize() uint64 {
	if x != nil {
		return x.EntrySize
	}
	return 0
}

// MaxTimestampNotSynced is the error variant that tells the peer has just become a leader and
// updating the max timestamp in the concurrency manager from PD TSO is ongoing. In this case,
// the prewrite of an async commit transaction cannot succeed. The client can backoff and
// resend the request.
type MaxTimestampNotSynced struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MaxTimestampNotSynced) Reset() {
	*x = MaxTimestampNotSynced{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MaxTimestampNotSynced) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MaxTimestampNotSynced) ProtoMessage() {}

func (x *MaxTimestampNotSynced) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MaxTimestampNotSynced.ProtoReflect.Descriptor instead.
func (*MaxTimestampNotSynced) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{12}
}

// ReadIndexNotReady is the error variant that tells the read index request is not ready, that is,
// the current region is in a status that not ready to serve the read index request. For example,
// region is in splitting or merging status.
// This can be retried at most time.
type ReadIndexNotReady struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The reason why the region is not ready to serve read index request
	Reason string `                   protobuf:"bytes,1,opt,name=reason,proto3"                   json:"reason,omitempty"`
	// The requested region ID
	RegionId      uint64 `                   protobuf:"varint,2,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReadIndexNotReady) Reset() {
	*x = ReadIndexNotReady{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadIndexNotReady) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadIndexNotReady) ProtoMessage() {}

func (x *ReadIndexNotReady) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadIndexNotReady.ProtoReflect.Descriptor instead.
func (*ReadIndexNotReady) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{13}
}

func (x *ReadIndexNotReady) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *ReadIndexNotReady) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

// ProposalInMergingMode is the error variant that tells the proposal is rejected because raft is
// in the merging mode. This may happen when BR/Lightning try to ingest SST.
// This can be retried at most time.
type ProposalInMergingMode struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The requested region ID
	RegionId      uint64 `                   protobuf:"varint,1,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProposalInMergingMode) Reset() {
	*x = ProposalInMergingMode{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProposalInMergingMode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProposalInMergingMode) ProtoMessage() {}

func (x *ProposalInMergingMode) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProposalInMergingMode.ProtoReflect.Descriptor instead.
func (*ProposalInMergingMode) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{14}
}

func (x *ProposalInMergingMode) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

type DataIsNotReady struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The requested region ID
	RegionId      uint64 `                   protobuf:"varint,1,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`
	PeerId        uint64 `                   protobuf:"varint,2,opt,name=peer_id,json=peerId,proto3"     json:"peer_id,omitempty"`
	SafeTs        uint64 `                   protobuf:"varint,3,opt,name=safe_ts,json=safeTs,proto3"     json:"safe_ts,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DataIsNotReady) Reset() {
	*x = DataIsNotReady{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DataIsNotReady) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataIsNotReady) ProtoMessage() {}

func (x *DataIsNotReady) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataIsNotReady.ProtoReflect.Descriptor instead.
func (*DataIsNotReady) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{15}
}

func (x *DataIsNotReady) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

func (x *DataIsNotReady) GetPeerId() uint64 {
	if x != nil {
		return x.PeerId
	}
	return 0
}

func (x *DataIsNotReady) GetSafeTs() uint64 {
	if x != nil {
		return x.SafeTs
	}
	return 0
}

type RecoveryInProgress struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The requested region ID
	RegionId      uint64 `                   protobuf:"varint,1,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RecoveryInProgress) Reset() {
	*x = RecoveryInProgress{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecoveryInProgress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecoveryInProgress) ProtoMessage() {}

func (x *RecoveryInProgress) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecoveryInProgress.ProtoReflect.Descriptor instead.
func (*RecoveryInProgress) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{16}
}

func (x *RecoveryInProgress) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

type FlashbackInProgress struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The requested region ID
	RegionId         uint64 `                   protobuf:"varint,1,opt,name=region_id,json=regionId,proto3"                  json:"region_id,omitempty"`
	FlashbackStartTs uint64 `                   protobuf:"varint,2,opt,name=flashback_start_ts,json=flashbackStartTs,proto3" json:"flashback_start_ts,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *FlashbackInProgress) Reset() {
	*x = FlashbackInProgress{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[17]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FlashbackInProgress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlashbackInProgress) ProtoMessage() {}

func (x *FlashbackInProgress) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[17]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlashbackInProgress.ProtoReflect.Descriptor instead.
func (*FlashbackInProgress) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{17}
}

func (x *FlashbackInProgress) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

func (x *FlashbackInProgress) GetFlashbackStartTs() uint64 {
	if x != nil {
		return x.FlashbackStartTs
	}
	return 0
}

type FlashbackNotPrepared struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The requested region ID
	RegionId      uint64 `                   protobuf:"varint,1,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FlashbackNotPrepared) Reset() {
	*x = FlashbackNotPrepared{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[18]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FlashbackNotPrepared) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlashbackNotPrepared) ProtoMessage() {}

func (x *FlashbackNotPrepared) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[18]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlashbackNotPrepared.ProtoReflect.Descriptor instead.
func (*FlashbackNotPrepared) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{18}
}

func (x *FlashbackNotPrepared) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

// MismatchPeerId is the error variant that tells the request is sent to wrong peer.
// Client receives this error should reload the region info and retry.
type MismatchPeerId struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RequestPeerId uint64                 `                   protobuf:"varint,1,opt,name=request_peer_id,json=requestPeerId,proto3" json:"request_peer_id,omitempty"`
	StorePeerId   uint64                 `                   protobuf:"varint,2,opt,name=store_peer_id,json=storePeerId,proto3"     json:"store_peer_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MismatchPeerId) Reset() {
	*x = MismatchPeerId{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[19]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MismatchPeerId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MismatchPeerId) ProtoMessage() {}

func (x *MismatchPeerId) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[19]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MismatchPeerId.ProtoReflect.Descriptor instead.
func (*MismatchPeerId) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{19}
}

func (x *MismatchPeerId) GetRequestPeerId() uint64 {
	if x != nil {
		return x.RequestPeerId
	}
	return 0
}

func (x *MismatchPeerId) GetStorePeerId() uint64 {
	if x != nil {
		return x.StorePeerId
	}
	return 0
}

// UndeterminedResult is the error variant that tells the result is not determined yet.
// For example, the raft protocol timed out and the apply result is unknown.
type UndeterminedResult struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `                   protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UndeterminedResult) Reset() {
	*x = UndeterminedResult{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[20]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UndeterminedResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UndeterminedResult) ProtoMessage() {}

func (x *UndeterminedResult) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[20]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UndeterminedResult.ProtoReflect.Descriptor instead.
func (*UndeterminedResult) Descriptor() ([]byte, []int) {
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{20}
}

func (x *UndeterminedResult) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Error wraps all region errors, indicates an error encountered by a request.
type Error struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The error message
	Message               string                 `                   protobuf:"bytes,1,opt,name=message,proto3"                                              json:"message,omitempty"`
	NotLeader             *NotLeader             `                   protobuf:"bytes,2,opt,name=not_leader,json=notLeader,proto3"                            json:"not_leader,omitempty"`
	RegionNotFound        *RegionNotFound        `                   protobuf:"bytes,3,opt,name=region_not_found,json=regionNotFound,proto3"                 json:"region_not_found,omitempty"`
	KeyNotInRegion        *KeyNotInRegion        `                   protobuf:"bytes,4,opt,name=key_not_in_region,json=keyNotInRegion,proto3"                json:"key_not_in_region,omitempty"`
	EpochNotMatch         *EpochNotMatch         `                   protobuf:"bytes,5,opt,name=epoch_not_match,json=epochNotMatch,proto3"                   json:"epoch_not_match,omitempty"`
	ServerIsBusy          *ServerIsBusy          `                   protobuf:"bytes,6,opt,name=server_is_busy,json=serverIsBusy,proto3"                     json:"server_is_busy,omitempty"`
	StaleCommand          *StaleCommand          `                   protobuf:"bytes,7,opt,name=stale_command,json=staleCommand,proto3"                      json:"stale_command,omitempty"`
	StoreNotMatch         *StoreNotMatch         `                   protobuf:"bytes,8,opt,name=store_not_match,json=storeNotMatch,proto3"                   json:"store_not_match,omitempty"`
	RaftEntryTooLarge     *RaftEntryTooLarge     `                   protobuf:"bytes,9,opt,name=raft_entry_too_large,json=raftEntryTooLarge,proto3"          json:"raft_entry_too_large,omitempty"`
	MaxTimestampNotSynced *MaxTimestampNotSynced `                   protobuf:"bytes,10,opt,name=max_timestamp_not_synced,json=maxTimestampNotSynced,proto3" json:"max_timestamp_not_synced,omitempty"`
	ReadIndexNotReady     *ReadIndexNotReady     `                   protobuf:"bytes,11,opt,name=read_index_not_ready,json=readIndexNotReady,proto3"         json:"read_index_not_ready,omitempty"`
	ProposalInMergingMode *ProposalInMergingMode `                   protobuf:"bytes,12,opt,name=proposal_in_merging_mode,json=proposalInMergingMode,proto3" json:"proposal_in_merging_mode,omitempty"`
	DataIsNotReady        *DataIsNotReady        `                   protobuf:"bytes,13,opt,name=data_is_not_ready,json=dataIsNotReady,proto3"               json:"data_is_not_ready,omitempty"`
	RegionNotInitialized  *RegionNotInitialized  `                   protobuf:"bytes,14,opt,name=region_not_initialized,json=regionNotInitialized,proto3"    json:"region_not_initialized,omitempty"`
	DiskFull              *DiskFull              `                   protobuf:"bytes,15,opt,name=disk_full,json=diskFull,proto3"                             json:"disk_full,omitempty"`
	// Online recovery is still in performing, reject writes to avoid potential issues
	RecoveryInProgress *RecoveryInProgress `                   protobuf:"bytes,16,opt,name=RecoveryInProgress,proto3"                                  json:"RecoveryInProgress,omitempty"`
	// Flashback is still in performing, reject any read or write to avoid potential issues.
	// NOTICE: this error is non-retryable, the request should fail ASAP when it meets this error.
	FlashbackInProgress *FlashbackInProgress `                   protobuf:"bytes,17,opt,name=FlashbackInProgress,proto3"                                 json:"FlashbackInProgress,omitempty"`
	// If the second phase flashback request is sent to a region that is not prepared for the flashback,
	// this error will be returned.
	// NOTICE: this error is non-retryable, the client should retry the first phase flashback request when it meets this error.
	FlashbackNotPrepared *FlashbackNotPrepared `                   protobuf:"bytes,18,opt,name=FlashbackNotPrepared,proto3"                                json:"FlashbackNotPrepared,omitempty"`
	// IsWitness is the error variant that tells a request be handle by witness
	// which should be forbidden and retry.
	IsWitness      *IsWitness      `                   protobuf:"bytes,19,opt,name=is_witness,json=isWitness,proto3"                           json:"is_witness,omitempty"`
	MismatchPeerId *MismatchPeerId `                   protobuf:"bytes,20,opt,name=mismatch_peer_id,json=mismatchPeerId,proto3"                json:"mismatch_peer_id,omitempty"`
	// BucketVersionNotMatch is the error variant that tells the request buckets version is not match.
	BucketVersionNotMatch *BucketVersionNotMatch `                   protobuf:"bytes,21,opt,name=bucket_version_not_match,json=bucketVersionNotMatch,proto3" json:"bucket_version_not_match,omitempty"`
	// UndeterminedResult is the error variant that tells the result is not determined yet.
	UndeterminedResult *UndeterminedResult `                   protobuf:"bytes,22,opt,name=undetermined_result,json=undeterminedResult,proto3"         json:"undetermined_result,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *Error) Reset() {
	*x = Error{}
	mi := &file_v1_tikv_errorpb_proto_msgTypes[21]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_errorpb_proto_msgTypes[21]
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
	return file_v1_tikv_errorpb_proto_rawDescGZIP(), []int{21}
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Error) GetNotLeader() *NotLeader {
	if x != nil {
		return x.NotLeader
	}
	return nil
}

func (x *Error) GetRegionNotFound() *RegionNotFound {
	if x != nil {
		return x.RegionNotFound
	}
	return nil
}

func (x *Error) GetKeyNotInRegion() *KeyNotInRegion {
	if x != nil {
		return x.KeyNotInRegion
	}
	return nil
}

func (x *Error) GetEpochNotMatch() *EpochNotMatch {
	if x != nil {
		return x.EpochNotMatch
	}
	return nil
}

func (x *Error) GetServerIsBusy() *ServerIsBusy {
	if x != nil {
		return x.ServerIsBusy
	}
	return nil
}

func (x *Error) GetStaleCommand() *StaleCommand {
	if x != nil {
		return x.StaleCommand
	}
	return nil
}

func (x *Error) GetStoreNotMatch() *StoreNotMatch {
	if x != nil {
		return x.StoreNotMatch
	}
	return nil
}

func (x *Error) GetRaftEntryTooLarge() *RaftEntryTooLarge {
	if x != nil {
		return x.RaftEntryTooLarge
	}
	return nil
}

func (x *Error) GetMaxTimestampNotSynced() *MaxTimestampNotSynced {
	if x != nil {
		return x.MaxTimestampNotSynced
	}
	return nil
}

func (x *Error) GetReadIndexNotReady() *ReadIndexNotReady {
	if x != nil {
		return x.ReadIndexNotReady
	}
	return nil
}

func (x *Error) GetProposalInMergingMode() *ProposalInMergingMode {
	if x != nil {
		return x.ProposalInMergingMode
	}
	return nil
}

func (x *Error) GetDataIsNotReady() *DataIsNotReady {
	if x != nil {
		return x.DataIsNotReady
	}
	return nil
}

func (x *Error) GetRegionNotInitialized() *RegionNotInitialized {
	if x != nil {
		return x.RegionNotInitialized
	}
	return nil
}

func (x *Error) GetDiskFull() *DiskFull {
	if x != nil {
		return x.DiskFull
	}
	return nil
}

func (x *Error) GetRecoveryInProgress() *RecoveryInProgress {
	if x != nil {
		return x.RecoveryInProgress
	}
	return nil
}

func (x *Error) GetFlashbackInProgress() *FlashbackInProgress {
	if x != nil {
		return x.FlashbackInProgress
	}
	return nil
}

func (x *Error) GetFlashbackNotPrepared() *FlashbackNotPrepared {
	if x != nil {
		return x.FlashbackNotPrepared
	}
	return nil
}

func (x *Error) GetIsWitness() *IsWitness {
	if x != nil {
		return x.IsWitness
	}
	return nil
}

func (x *Error) GetMismatchPeerId() *MismatchPeerId {
	if x != nil {
		return x.MismatchPeerId
	}
	return nil
}

func (x *Error) GetBucketVersionNotMatch() *BucketVersionNotMatch {
	if x != nil {
		return x.BucketVersionNotMatch
	}
	return nil
}

func (x *Error) GetUndeterminedResult() *UndeterminedResult {
	if x != nil {
		return x.UndeterminedResult
	}
	return nil
}

var File_v1_tikv_errorpb_proto protoreflect.FileDescriptor

const file_v1_tikv_errorpb_proto_rawDesc = "" +
	"\n" +
	"\x15v1/tikv/errorpb.proto\x12\x04tikv\x1a\x14v1/tikv/metapb.proto\"N\n" +
	"\tNotLeader\x12\x1b\n" +
	"\tregion_id\x18\x01 \x01(\x04R\bregionId\x12$\n" +
	"\x06leader\x18\x02 \x01(\v2\f.metapb.PeerR\x06leader\"(\n" +
	"\tIsWitness\x12\x1b\n" +
	"\tregion_id\x18\x01 \x01(\x04R\bregionId\"E\n" +
	"\x15BucketVersionNotMatch\x12\x18\n" +
	"\aversion\x18\x01 \x01(\x04R\aversion\x12\x12\n" +
	"\x04keys\x18\x02 \x03(\fR\x04keys\"=\n" +
	"\bDiskFull\x12\x19\n" +
	"\bstore_id\x18\x01 \x03(\x04R\astoreId\x12\x16\n" +
	"\x06reason\x18\x02 \x01(\tR\x06reason\"a\n" +
	"\rStoreNotMatch\x12(\n" +
	"\x10request_store_id\x18\x01 \x01(\x04R\x0erequestStoreId\x12&\n" +
	"\x0factual_store_id\x18\x02 \x01(\x04R\ractualStoreId\"-\n" +
	"\x0eRegionNotFound\x12\x1b\n" +
	"\tregion_id\x18\x01 \x01(\x04R\bregionId\"3\n" +
	"\x14RegionNotInitialized\x12\x1b\n" +
	"\tregion_id\x18\x01 \x01(\x04R\bregionId\"u\n" +
	"\x0eKeyNotInRegion\x12\x10\n" +
	"\x03key\x18\x01 \x01(\fR\x03key\x12\x1b\n" +
	"\tregion_id\x18\x02 \x01(\x04R\bregionId\x12\x1b\n" +
	"\tstart_key\x18\x03 \x01(\fR\bstartKey\x12\x17\n" +
	"\aend_key\x18\x04 \x01(\fR\x06endKey\"I\n" +
	"\rEpochNotMatch\x128\n" +
	"\x0fcurrent_regions\x18\x01 \x03(\v2\x0f.metapb.Region2R\x0ecurrentRegions\"\x96\x01\n" +
	"\fServerIsBusy\x12\x16\n" +
	"\x06reason\x18\x01 \x01(\tR\x06reason\x12\x1d\n" +
	"\n" +
	"backoff_ms\x18\x02 \x01(\x04R\tbackoffMs\x12*\n" +
	"\x11estimated_wait_ms\x18\x03 \x01(\rR\x0festimatedWaitMs\x12#\n" +
	"\rapplied_index\x18\x04 \x01(\x04R\fappliedIndex\"\x0e\n" +
	"\fStaleCommand\"O\n" +
	"\x11RaftEntryTooLarge\x12\x1b\n" +
	"\tregion_id\x18\x01 \x01(\x04R\bregionId\x12\x1d\n" +
	"\n" +
	"entry_size\x18\x02 \x01(\x04R\tentrySize\"\x17\n" +
	"\x15MaxTimestampNotSynced\"H\n" +
	"\x11ReadIndexNotReady\x12\x16\n" +
	"\x06reason\x18\x01 \x01(\tR\x06reason\x12\x1b\n" +
	"\tregion_id\x18\x02 \x01(\x04R\bregionId\"4\n" +
	"\x15ProposalInMergingMode\x12\x1b\n" +
	"\tregion_id\x18\x01 \x01(\x04R\bregionId\"_\n" +
	"\x0eDataIsNotReady\x12\x1b\n" +
	"\tregion_id\x18\x01 \x01(\x04R\bregionId\x12\x17\n" +
	"\apeer_id\x18\x02 \x01(\x04R\x06peerId\x12\x17\n" +
	"\asafe_ts\x18\x03 \x01(\x04R\x06safeTs\"1\n" +
	"\x12RecoveryInProgress\x12\x1b\n" +
	"\tregion_id\x18\x01 \x01(\x04R\bregionId\"`\n" +
	"\x13FlashbackInProgress\x12\x1b\n" +
	"\tregion_id\x18\x01 \x01(\x04R\bregionId\x12,\n" +
	"\x12flashback_start_ts\x18\x02 \x01(\x04R\x10flashbackStartTs\"3\n" +
	"\x14FlashbackNotPrepared\x12\x1b\n" +
	"\tregion_id\x18\x01 \x01(\x04R\bregionId\"\\\n" +
	"\x0eMismatchPeerId\x12&\n" +
	"\x0frequest_peer_id\x18\x01 \x01(\x04R\rrequestPeerId\x12\"\n" +
	"\rstore_peer_id\x18\x02 \x01(\x04R\vstorePeerId\".\n" +
	"\x12UndeterminedResult\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"\xc4\v\n" +
	"\x05Error\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\x12.\n" +
	"\n" +
	"not_leader\x18\x02 \x01(\v2\x0f.tikv.NotLeaderR\tnotLeader\x12>\n" +
	"\x10region_not_found\x18\x03 \x01(\v2\x14.tikv.RegionNotFoundR\x0eregionNotFound\x12?\n" +
	"\x11key_not_in_region\x18\x04 \x01(\v2\x14.tikv.KeyNotInRegionR\x0ekeyNotInRegion\x12;\n" +
	"\x0fepoch_not_match\x18\x05 \x01(\v2\x13.tikv.EpochNotMatchR\repochNotMatch\x128\n" +
	"\x0eserver_is_busy\x18\x06 \x01(\v2\x12.tikv.ServerIsBusyR\fserverIsBusy\x127\n" +
	"\rstale_command\x18\a \x01(\v2\x12.tikv.StaleCommandR\fstaleCommand\x12;\n" +
	"\x0fstore_not_match\x18\b \x01(\v2\x13.tikv.StoreNotMatchR\rstoreNotMatch\x12H\n" +
	"\x14raft_entry_too_large\x18\t \x01(\v2\x17.tikv.RaftEntryTooLargeR\x11raftEntryTooLarge\x12T\n" +
	"\x18max_timestamp_not_synced\x18\n" +
	" \x01(\v2\x1b.tikv.MaxTimestampNotSyncedR\x15maxTimestampNotSynced\x12H\n" +
	"\x14read_index_not_ready\x18\v \x01(\v2\x17.tikv.ReadIndexNotReadyR\x11readIndexNotReady\x12T\n" +
	"\x18proposal_in_merging_mode\x18\f \x01(\v2\x1b.tikv.ProposalInMergingModeR\x15proposalInMergingMode\x12?\n" +
	"\x11data_is_not_ready\x18\r \x01(\v2\x14.tikv.DataIsNotReadyR\x0edataIsNotReady\x12P\n" +
	"\x16region_not_initialized\x18\x0e \x01(\v2\x1a.tikv.RegionNotInitializedR\x14regionNotInitialized\x12+\n" +
	"\tdisk_full\x18\x0f \x01(\v2\x0e.tikv.DiskFullR\bdiskFull\x12H\n" +
	"\x12RecoveryInProgress\x18\x10 \x01(\v2\x18.tikv.RecoveryInProgressR\x12RecoveryInProgress\x12K\n" +
	"\x13FlashbackInProgress\x18\x11 \x01(\v2\x19.tikv.FlashbackInProgressR\x13FlashbackInProgress\x12N\n" +
	"\x14FlashbackNotPrepared\x18\x12 \x01(\v2\x1a.tikv.FlashbackNotPreparedR\x14FlashbackNotPrepared\x12.\n" +
	"\n" +
	"is_witness\x18\x13 \x01(\v2\x0f.tikv.IsWitnessR\tisWitness\x12>\n" +
	"\x10mismatch_peer_id\x18\x14 \x01(\v2\x14.tikv.MismatchPeerIdR\x0emismatchPeerId\x12T\n" +
	"\x18bucket_version_not_match\x18\x15 \x01(\v2\x1b.tikv.BucketVersionNotMatchR\x15bucketVersionNotMatch\x12I\n" +
	"\x13undetermined_result\x18\x16 \x01(\v2\x18.tikv.UndeterminedResultR\x12undeterminedResultR\vstale_epochB)Z'github.com/vdaas/vald/apis/grpc/v1/tikvb\x06proto3"

var (
	file_v1_tikv_errorpb_proto_rawDescOnce sync.Once
	file_v1_tikv_errorpb_proto_rawDescData []byte
)

func file_v1_tikv_errorpb_proto_rawDescGZIP() []byte {
	file_v1_tikv_errorpb_proto_rawDescOnce.Do(func() {
		file_v1_tikv_errorpb_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_tikv_errorpb_proto_rawDesc), len(file_v1_tikv_errorpb_proto_rawDesc)))
	})
	return file_v1_tikv_errorpb_proto_rawDescData
}

var (
	file_v1_tikv_errorpb_proto_msgTypes = make([]protoimpl.MessageInfo, 22)
	file_v1_tikv_errorpb_proto_goTypes  = []any{
		(*NotLeader)(nil),             // 0: tikv.NotLeader
		(*IsWitness)(nil),             // 1: tikv.IsWitness
		(*BucketVersionNotMatch)(nil), // 2: tikv.BucketVersionNotMatch
		(*DiskFull)(nil),              // 3: tikv.DiskFull
		(*StoreNotMatch)(nil),         // 4: tikv.StoreNotMatch
		(*RegionNotFound)(nil),        // 5: tikv.RegionNotFound
		(*RegionNotInitialized)(nil),  // 6: tikv.RegionNotInitialized
		(*KeyNotInRegion)(nil),        // 7: tikv.KeyNotInRegion
		(*EpochNotMatch)(nil),         // 8: tikv.EpochNotMatch
		(*ServerIsBusy)(nil),          // 9: tikv.ServerIsBusy
		(*StaleCommand)(nil),          // 10: tikv.StaleCommand
		(*RaftEntryTooLarge)(nil),     // 11: tikv.RaftEntryTooLarge
		(*MaxTimestampNotSynced)(nil), // 12: tikv.MaxTimestampNotSynced
		(*ReadIndexNotReady)(nil),     // 13: tikv.ReadIndexNotReady
		(*ProposalInMergingMode)(nil), // 14: tikv.ProposalInMergingMode
		(*DataIsNotReady)(nil),        // 15: tikv.DataIsNotReady
		(*RecoveryInProgress)(nil),    // 16: tikv.RecoveryInProgress
		(*FlashbackInProgress)(nil),   // 17: tikv.FlashbackInProgress
		(*FlashbackNotPrepared)(nil),  // 18: tikv.FlashbackNotPrepared
		(*MismatchPeerId)(nil),        // 19: tikv.MismatchPeerId
		(*UndeterminedResult)(nil),    // 20: tikv.UndeterminedResult
		(*Error)(nil),                 // 21: tikv.Error
		(*Peer)(nil),                  // 22: metapb.Peer
		(*Region2)(nil),               // 23: metapb.Region2
	}
)

var file_v1_tikv_errorpb_proto_depIdxs = []int32{
	22, // 0: tikv.NotLeader.leader:type_name -> metapb.Peer
	23, // 1: tikv.EpochNotMatch.current_regions:type_name -> metapb.Region2
	0,  // 2: tikv.Error.not_leader:type_name -> tikv.NotLeader
	5,  // 3: tikv.Error.region_not_found:type_name -> tikv.RegionNotFound
	7,  // 4: tikv.Error.key_not_in_region:type_name -> tikv.KeyNotInRegion
	8,  // 5: tikv.Error.epoch_not_match:type_name -> tikv.EpochNotMatch
	9,  // 6: tikv.Error.server_is_busy:type_name -> tikv.ServerIsBusy
	10, // 7: tikv.Error.stale_command:type_name -> tikv.StaleCommand
	4,  // 8: tikv.Error.store_not_match:type_name -> tikv.StoreNotMatch
	11, // 9: tikv.Error.raft_entry_too_large:type_name -> tikv.RaftEntryTooLarge
	12, // 10: tikv.Error.max_timestamp_not_synced:type_name -> tikv.MaxTimestampNotSynced
	13, // 11: tikv.Error.read_index_not_ready:type_name -> tikv.ReadIndexNotReady
	14, // 12: tikv.Error.proposal_in_merging_mode:type_name -> tikv.ProposalInMergingMode
	15, // 13: tikv.Error.data_is_not_ready:type_name -> tikv.DataIsNotReady
	6,  // 14: tikv.Error.region_not_initialized:type_name -> tikv.RegionNotInitialized
	3,  // 15: tikv.Error.disk_full:type_name -> tikv.DiskFull
	16, // 16: tikv.Error.RecoveryInProgress:type_name -> tikv.RecoveryInProgress
	17, // 17: tikv.Error.FlashbackInProgress:type_name -> tikv.FlashbackInProgress
	18, // 18: tikv.Error.FlashbackNotPrepared:type_name -> tikv.FlashbackNotPrepared
	1,  // 19: tikv.Error.is_witness:type_name -> tikv.IsWitness
	19, // 20: tikv.Error.mismatch_peer_id:type_name -> tikv.MismatchPeerId
	2,  // 21: tikv.Error.bucket_version_not_match:type_name -> tikv.BucketVersionNotMatch
	20, // 22: tikv.Error.undetermined_result:type_name -> tikv.UndeterminedResult
	23, // [23:23] is the sub-list for method output_type
	23, // [23:23] is the sub-list for method input_type
	23, // [23:23] is the sub-list for extension type_name
	23, // [23:23] is the sub-list for extension extendee
	0,  // [0:23] is the sub-list for field type_name
}

func init() { file_v1_tikv_errorpb_proto_init() }
func file_v1_tikv_errorpb_proto_init() {
	if File_v1_tikv_errorpb_proto != nil {
		return
	}
	file_v1_tikv_metapb_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_tikv_errorpb_proto_rawDesc), len(file_v1_tikv_errorpb_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   22,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_tikv_errorpb_proto_goTypes,
		DependencyIndexes: file_v1_tikv_errorpb_proto_depIdxs,
		MessageInfos:      file_v1_tikv_errorpb_proto_msgTypes,
	}.Build()
	File_v1_tikv_errorpb_proto = out.File
	file_v1_tikv_errorpb_proto_goTypes = nil
	file_v1_tikv_errorpb_proto_depIdxs = nil
}
