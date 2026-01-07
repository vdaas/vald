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
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "github.com/vdaas/vald/internal/sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The API version the server and the client is using.
// See more details in https://github.com/tikv/rfcs/blob/master/text/0069-api-v2.md.
type APIVersion int32

const (
	// `V1` is mainly for TiDB & TxnKV, and is not safe to use RawKV along with the others.
	// V1 server only accepts V1 requests. V1 raw requests with TTL will be rejected.
	APIVersion_V1 APIVersion = 0
	// `V1TTL` is only available to RawKV, and 8 bytes representing the unix timestamp in
	// seconds for expiring time will be append to the value of all RawKV entries. For example:
	// ------------------------------------------------------------
	// | User value     | Expire Ts                               |
	// ------------------------------------------------------------
	// | 0x12 0x34 0x56 | 0x00 0x00 0x00 0x00 0x00 0x00 0xff 0xff |
	// ------------------------------------------------------------
	// V1TTL server only accepts V1 raw requests.
	// V1 client should not use `V1TTL` in request. V1 client should always send `V1`.
	APIVersion_V1TTL APIVersion = 1
	// `V2` use new encoding for RawKV & TxnKV to support more features.
	//
	// Key Encoding:
	//
	//	TiDB: start with `m` or `t`, the same as `V1`.
	//	TxnKV: prefix with `x`, encoded as `MCE( x{keyspace id} + {user key} ) + timestamp`.
	//	RawKV: prefix with `r`, encoded as `MCE( r{keyspace id} + {user key} ) + timestamp`.
	//	Where the `{keyspace id}` is fixed-length of 3 bytes in network byte order.
	//	Besides, RawKV entires must be in `default` CF.
	//
	// Value Encoding:
	//
	//	TiDB & TxnKV: the same as `V1`.
	//	RawKV: `{user value} + {optional fields} + {meta flag}`. The last byte in the
	//	raw value must be meta flags. For example:
	//	--------------------------------------
	//	| User value     | Meta flags        |
	//	--------------------------------------
	//	| 0x12 0x34 0x56 | 0x00 (0b00000000) |
	//	--------------------------------------
	//	Bit 0 of meta flags is for TTL. If set, the value contains 8 bytes expiring time as
	//	unix timestamp in seconds at the very left to the meta flags.
	//	--------------------------------------------------------------------------------
	//	| User value     | Expiring time                           | Meta flags        |
	//	--------------------------------------------------------------------------------
	//	| 0x12 0x34 0x56 | 0x00 0x00 0x00 0x00 0x00 0x00 0xff 0xff | 0x01 (0b00000001) |
	//	--------------------------------------------------------------------------------
	//	Bit 1 is for deletion. If set, the entry is logical deleted.
	//	---------------------
	//	| Meta flags        |
	//	---------------------
	//	| 0x02 (0b00000010) |
	//	---------------------
	//
	// V2 server accpets V2 requests and V1 transactional requests that statrts with TiDB key
	// prefix (`m` and `t`).
	APIVersion_V2 APIVersion = 2
)

// Enum value maps for APIVersion.
var (
	APIVersion_name = map[int32]string{
		0: "V1",
		1: "V1TTL",
		2: "V2",
	}
	APIVersion_value = map[string]int32{
		"V1":    0,
		"V1TTL": 1,
		"V2":    2,
	}
)

func (x APIVersion) Enum() *APIVersion {
	p := new(APIVersion)
	*p = x
	return p
}

func (x APIVersion) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (APIVersion) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_tikv_kvrpcpb_proto_enumTypes[0].Descriptor()
}

func (APIVersion) Type() protoreflect.EnumType {
	return &file_v1_tikv_kvrpcpb_proto_enumTypes[0]
}

func (x APIVersion) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use APIVersion.Descriptor instead.
func (APIVersion) EnumDescriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{0}
}

type CommandPri int32

const (
	CommandPri_Normal CommandPri = 0 // Normal is the default value.
	CommandPri_Low    CommandPri = 1
	CommandPri_High   CommandPri = 2
)

// Enum value maps for CommandPri.
var (
	CommandPri_name = map[int32]string{
		0: "Normal",
		1: "Low",
		2: "High",
	}
	CommandPri_value = map[string]int32{
		"Normal": 0,
		"Low":    1,
		"High":   2,
	}
)

func (x CommandPri) Enum() *CommandPri {
	p := new(CommandPri)
	*p = x
	return p
}

func (x CommandPri) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CommandPri) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_tikv_kvrpcpb_proto_enumTypes[1].Descriptor()
}

func (CommandPri) Type() protoreflect.EnumType {
	return &file_v1_tikv_kvrpcpb_proto_enumTypes[1]
}

func (x CommandPri) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CommandPri.Descriptor instead.
func (CommandPri) EnumDescriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{1}
}

type IsolationLevel int32

const (
	IsolationLevel_SI        IsolationLevel = 0 // SI = snapshot isolation
	IsolationLevel_RC        IsolationLevel = 1 // RC = read committed
	IsolationLevel_RCCheckTS IsolationLevel = 2 // RC read and it's needed to check if there exists more recent versions.
)

// Enum value maps for IsolationLevel.
var (
	IsolationLevel_name = map[int32]string{
		0: "SI",
		1: "RC",
		2: "RCCheckTS",
	}
	IsolationLevel_value = map[string]int32{
		"SI":        0,
		"RC":        1,
		"RCCheckTS": 2,
	}
)

func (x IsolationLevel) Enum() *IsolationLevel {
	p := new(IsolationLevel)
	*p = x
	return p
}

func (x IsolationLevel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (IsolationLevel) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_tikv_kvrpcpb_proto_enumTypes[2].Descriptor()
}

func (IsolationLevel) Type() protoreflect.EnumType {
	return &file_v1_tikv_kvrpcpb_proto_enumTypes[2]
}

func (x IsolationLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use IsolationLevel.Descriptor instead.
func (IsolationLevel) EnumDescriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{2}
}

// Operation allowed info during each TiKV storage threshold.
type DiskFullOpt int32

const (
	DiskFullOpt_NotAllowedOnFull     DiskFullOpt = 0 // The default value, means operations are not allowed either under almost full or already full.
	DiskFullOpt_AllowedOnAlmostFull  DiskFullOpt = 1 // Means operations will be allowed when disk is almost full.
	DiskFullOpt_AllowedOnAlreadyFull DiskFullOpt = 2 // Means operations will be allowed when disk is already full.
)

// Enum value maps for DiskFullOpt.
var (
	DiskFullOpt_name = map[int32]string{
		0: "NotAllowedOnFull",
		1: "AllowedOnAlmostFull",
		2: "AllowedOnAlreadyFull",
	}
	DiskFullOpt_value = map[string]int32{
		"NotAllowedOnFull":     0,
		"AllowedOnAlmostFull":  1,
		"AllowedOnAlreadyFull": 2,
	}
)

func (x DiskFullOpt) Enum() *DiskFullOpt {
	p := new(DiskFullOpt)
	*p = x
	return p
}

func (x DiskFullOpt) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DiskFullOpt) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_tikv_kvrpcpb_proto_enumTypes[3].Descriptor()
}

func (DiskFullOpt) Type() protoreflect.EnumType {
	return &file_v1_tikv_kvrpcpb_proto_enumTypes[3]
}

func (x DiskFullOpt) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DiskFullOpt.Descriptor instead.
func (DiskFullOpt) EnumDescriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{3}
}

type Op int32

const (
	Op_Put      Op = 0
	Op_Del      Op = 1
	Op_Lock     Op = 2
	Op_Rollback Op = 3
	// insert operation has a constraint that key should not exist before.
	Op_Insert          Op = 4
	Op_PessimisticLock Op = 5
	Op_CheckNotExists  Op = 6
)

// Enum value maps for Op.
var (
	Op_name = map[int32]string{
		0: "Put",
		1: "Del",
		2: "Lock",
		3: "Rollback",
		4: "Insert",
		5: "PessimisticLock",
		6: "CheckNotExists",
	}
	Op_value = map[string]int32{
		"Put":             0,
		"Del":             1,
		"Lock":            2,
		"Rollback":        3,
		"Insert":          4,
		"PessimisticLock": 5,
		"CheckNotExists":  6,
	}
)

func (x Op) Enum() *Op {
	p := new(Op)
	*p = x
	return p
}

func (x Op) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Op) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_tikv_kvrpcpb_proto_enumTypes[4].Descriptor()
}

func (Op) Type() protoreflect.EnumType {
	return &file_v1_tikv_kvrpcpb_proto_enumTypes[4]
}

func (x Op) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Op.Descriptor instead.
func (Op) EnumDescriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{4}
}

type Assertion int32

const (
	Assertion_None     Assertion = 0
	Assertion_Exist    Assertion = 1
	Assertion_NotExist Assertion = 2
)

// Enum value maps for Assertion.
var (
	Assertion_name = map[int32]string{
		0: "None",
		1: "Exist",
		2: "NotExist",
	}
	Assertion_value = map[string]int32{
		"None":     0,
		"Exist":    1,
		"NotExist": 2,
	}
)

func (x Assertion) Enum() *Assertion {
	p := new(Assertion)
	*p = x
	return p
}

func (x Assertion) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Assertion) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_tikv_kvrpcpb_proto_enumTypes[5].Descriptor()
}

func (Assertion) Type() protoreflect.EnumType {
	return &file_v1_tikv_kvrpcpb_proto_enumTypes[5]
}

func (x Assertion) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Assertion.Descriptor instead.
func (Assertion) EnumDescriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{5}
}

type WriteConflict_Reason int32

const (
	WriteConflict_Unknown              WriteConflict_Reason = 0
	WriteConflict_Optimistic           WriteConflict_Reason = 1 // in optimistic transactions.
	WriteConflict_PessimisticRetry     WriteConflict_Reason = 2 // a lock acquisition request waits for a lock and awakes, or meets a newer version of data, let TiDB retry.
	WriteConflict_SelfRolledBack       WriteConflict_Reason = 3 // the transaction itself has been rolled back when it tries to prewrite.
	WriteConflict_RcCheckTs            WriteConflict_Reason = 4 // RcCheckTs failure by meeting a newer version, let TiDB retry.
	WriteConflict_LazyUniquenessCheck  WriteConflict_Reason = 5 // write conflict found when deferring constraint checks in pessimistic transactions. Deprecated in next-gen (cloud-storage-engine).
	WriteConflict_NotLockedKeyConflict WriteConflict_Reason = 6 // write conflict found on keys that do not acquire pessimistic locks in pessimistic transactions.
)

// Enum value maps for WriteConflict_Reason.
var (
	WriteConflict_Reason_name = map[int32]string{
		0: "Unknown",
		1: "Optimistic",
		2: "PessimisticRetry",
		3: "SelfRolledBack",
		4: "RcCheckTs",
		5: "LazyUniquenessCheck",
		6: "NotLockedKeyConflict",
	}
	WriteConflict_Reason_value = map[string]int32{
		"Unknown":              0,
		"Optimistic":           1,
		"PessimisticRetry":     2,
		"SelfRolledBack":       3,
		"RcCheckTs":            4,
		"LazyUniquenessCheck":  5,
		"NotLockedKeyConflict": 6,
	}
)

func (x WriteConflict_Reason) Enum() *WriteConflict_Reason {
	p := new(WriteConflict_Reason)
	*p = x
	return p
}

func (x WriteConflict_Reason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WriteConflict_Reason) Descriptor() protoreflect.EnumDescriptor {
	return file_v1_tikv_kvrpcpb_proto_enumTypes[6].Descriptor()
}

func (WriteConflict_Reason) Type() protoreflect.EnumType {
	return &file_v1_tikv_kvrpcpb_proto_enumTypes[6]
}

func (x WriteConflict_Reason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WriteConflict_Reason.Descriptor instead.
func (WriteConflict_Reason) EnumDescriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{17, 0}
}

type RawGetRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Context       *Context               `protobuf:"bytes,1,opt,name=context,proto3" json:"context,omitempty"`
	Key           []byte                 `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Cf            string                 `protobuf:"bytes,3,opt,name=cf,proto3" json:"cf,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawGetRequest) Reset() {
	*x = RawGetRequest{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawGetRequest) ProtoMessage() {}

func (x *RawGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawGetRequest.ProtoReflect.Descriptor instead.
func (*RawGetRequest) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{0}
}

func (x *RawGetRequest) GetContext() *Context {
	if x != nil {
		return x.Context
	}
	return nil
}

func (x *RawGetRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *RawGetRequest) GetCf() string {
	if x != nil {
		return x.Cf
	}
	return ""
}

type RawGetResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RegionError   *Error                 `protobuf:"bytes,1,opt,name=region_error,json=regionError,proto3" json:"region_error,omitempty"`
	Error         string                 `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	Value         []byte                 `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	NotFound      bool                   `protobuf:"varint,4,opt,name=not_found,json=notFound,proto3" json:"not_found,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawGetResponse) Reset() {
	*x = RawGetResponse{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawGetResponse) ProtoMessage() {}

func (x *RawGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawGetResponse.ProtoReflect.Descriptor instead.
func (*RawGetResponse) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{1}
}

func (x *RawGetResponse) GetRegionError() *Error {
	if x != nil {
		return x.RegionError
	}
	return nil
}

func (x *RawGetResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *RawGetResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *RawGetResponse) GetNotFound() bool {
	if x != nil {
		return x.NotFound
	}
	return false
}

type RawBatchGetRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Context       *Context               `protobuf:"bytes,1,opt,name=context,proto3" json:"context,omitempty"`
	Keys          [][]byte               `protobuf:"bytes,2,rep,name=keys,proto3" json:"keys,omitempty"`
	Cf            string                 `protobuf:"bytes,3,opt,name=cf,proto3" json:"cf,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawBatchGetRequest) Reset() {
	*x = RawBatchGetRequest{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawBatchGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawBatchGetRequest) ProtoMessage() {}

func (x *RawBatchGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawBatchGetRequest.ProtoReflect.Descriptor instead.
func (*RawBatchGetRequest) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{2}
}

func (x *RawBatchGetRequest) GetContext() *Context {
	if x != nil {
		return x.Context
	}
	return nil
}

func (x *RawBatchGetRequest) GetKeys() [][]byte {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *RawBatchGetRequest) GetCf() string {
	if x != nil {
		return x.Cf
	}
	return ""
}

type RawBatchGetResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RegionError   *Error                 `protobuf:"bytes,1,opt,name=region_error,json=regionError,proto3" json:"region_error,omitempty"`
	Pairs         []*KvPair              `protobuf:"bytes,2,rep,name=pairs,proto3" json:"pairs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawBatchGetResponse) Reset() {
	*x = RawBatchGetResponse{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawBatchGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawBatchGetResponse) ProtoMessage() {}

func (x *RawBatchGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawBatchGetResponse.ProtoReflect.Descriptor instead.
func (*RawBatchGetResponse) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{3}
}

func (x *RawBatchGetResponse) GetRegionError() *Error {
	if x != nil {
		return x.RegionError
	}
	return nil
}

func (x *RawBatchGetResponse) GetPairs() []*KvPair {
	if x != nil {
		return x.Pairs
	}
	return nil
}

type RawPutRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Context       *Context               `protobuf:"bytes,1,opt,name=context,proto3" json:"context,omitempty"`
	Key           []byte                 `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value         []byte                 `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	Cf            string                 `protobuf:"bytes,4,opt,name=cf,proto3" json:"cf,omitempty"`
	Ttl           uint64                 `protobuf:"varint,5,opt,name=ttl,proto3" json:"ttl,omitempty"`
	ForCas        bool                   `protobuf:"varint,6,opt,name=for_cas,json=forCas,proto3" json:"for_cas,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawPutRequest) Reset() {
	*x = RawPutRequest{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawPutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawPutRequest) ProtoMessage() {}

func (x *RawPutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawPutRequest.ProtoReflect.Descriptor instead.
func (*RawPutRequest) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{4}
}

func (x *RawPutRequest) GetContext() *Context {
	if x != nil {
		return x.Context
	}
	return nil
}

func (x *RawPutRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *RawPutRequest) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *RawPutRequest) GetCf() string {
	if x != nil {
		return x.Cf
	}
	return ""
}

func (x *RawPutRequest) GetTtl() uint64 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

func (x *RawPutRequest) GetForCas() bool {
	if x != nil {
		return x.ForCas
	}
	return false
}

type RawPutResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RegionError   *Error                 `protobuf:"bytes,1,opt,name=region_error,json=regionError,proto3" json:"region_error,omitempty"`
	Error         string                 `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawPutResponse) Reset() {
	*x = RawPutResponse{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawPutResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawPutResponse) ProtoMessage() {}

func (x *RawPutResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawPutResponse.ProtoReflect.Descriptor instead.
func (*RawPutResponse) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{5}
}

func (x *RawPutResponse) GetRegionError() *Error {
	if x != nil {
		return x.RegionError
	}
	return nil
}

func (x *RawPutResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type RawBatchPutRequest struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	Context *Context               `protobuf:"bytes,1,opt,name=context,proto3" json:"context,omitempty"`
	Pairs   []*KvPair              `protobuf:"bytes,2,rep,name=pairs,proto3" json:"pairs,omitempty"`
	Cf      string                 `protobuf:"bytes,3,opt,name=cf,proto3" json:"cf,omitempty"`
	// Deprecated: Marked as deprecated in v1/tikv/kvrpcpb.proto.
	Ttl    uint64 `protobuf:"varint,4,opt,name=ttl,proto3" json:"ttl,omitempty"`
	ForCas bool   `protobuf:"varint,5,opt,name=for_cas,json=forCas,proto3" json:"for_cas,omitempty"`
	// The time-to-live for each keys in seconds, and if the length of `ttls`
	// is exactly one, the ttl will be applied to all keys. Otherwise, the length
	// mismatch between `ttls` and `pairs` will return an error.
	Ttls          []uint64 `protobuf:"varint,6,rep,packed,name=ttls,proto3" json:"ttls,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawBatchPutRequest) Reset() {
	*x = RawBatchPutRequest{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawBatchPutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawBatchPutRequest) ProtoMessage() {}

func (x *RawBatchPutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawBatchPutRequest.ProtoReflect.Descriptor instead.
func (*RawBatchPutRequest) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{6}
}

func (x *RawBatchPutRequest) GetContext() *Context {
	if x != nil {
		return x.Context
	}
	return nil
}

func (x *RawBatchPutRequest) GetPairs() []*KvPair {
	if x != nil {
		return x.Pairs
	}
	return nil
}

func (x *RawBatchPutRequest) GetCf() string {
	if x != nil {
		return x.Cf
	}
	return ""
}

// Deprecated: Marked as deprecated in v1/tikv/kvrpcpb.proto.
func (x *RawBatchPutRequest) GetTtl() uint64 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

func (x *RawBatchPutRequest) GetForCas() bool {
	if x != nil {
		return x.ForCas
	}
	return false
}

func (x *RawBatchPutRequest) GetTtls() []uint64 {
	if x != nil {
		return x.Ttls
	}
	return nil
}

type RawBatchPutResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RegionError   *Error                 `protobuf:"bytes,1,opt,name=region_error,json=regionError,proto3" json:"region_error,omitempty"`
	Error         string                 `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawBatchPutResponse) Reset() {
	*x = RawBatchPutResponse{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawBatchPutResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawBatchPutResponse) ProtoMessage() {}

func (x *RawBatchPutResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawBatchPutResponse.ProtoReflect.Descriptor instead.
func (*RawBatchPutResponse) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{7}
}

func (x *RawBatchPutResponse) GetRegionError() *Error {
	if x != nil {
		return x.RegionError
	}
	return nil
}

func (x *RawBatchPutResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type RawDeleteRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Context       *Context               `protobuf:"bytes,1,opt,name=context,proto3" json:"context,omitempty"`
	Key           []byte                 `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Cf            string                 `protobuf:"bytes,3,opt,name=cf,proto3" json:"cf,omitempty"`
	ForCas        bool                   `protobuf:"varint,4,opt,name=for_cas,json=forCas,proto3" json:"for_cas,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawDeleteRequest) Reset() {
	*x = RawDeleteRequest{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawDeleteRequest) ProtoMessage() {}

func (x *RawDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawDeleteRequest.ProtoReflect.Descriptor instead.
func (*RawDeleteRequest) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{8}
}

func (x *RawDeleteRequest) GetContext() *Context {
	if x != nil {
		return x.Context
	}
	return nil
}

func (x *RawDeleteRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *RawDeleteRequest) GetCf() string {
	if x != nil {
		return x.Cf
	}
	return ""
}

func (x *RawDeleteRequest) GetForCas() bool {
	if x != nil {
		return x.ForCas
	}
	return false
}

type RawDeleteResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RegionError   *Error                 `protobuf:"bytes,1,opt,name=region_error,json=regionError,proto3" json:"region_error,omitempty"`
	Error         string                 `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawDeleteResponse) Reset() {
	*x = RawDeleteResponse{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawDeleteResponse) ProtoMessage() {}

func (x *RawDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawDeleteResponse.ProtoReflect.Descriptor instead.
func (*RawDeleteResponse) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{9}
}

func (x *RawDeleteResponse) GetRegionError() *Error {
	if x != nil {
		return x.RegionError
	}
	return nil
}

func (x *RawDeleteResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type RawBatchDeleteRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Context       *Context               `protobuf:"bytes,1,opt,name=context,proto3" json:"context,omitempty"`
	Keys          [][]byte               `protobuf:"bytes,2,rep,name=keys,proto3" json:"keys,omitempty"`
	Cf            string                 `protobuf:"bytes,3,opt,name=cf,proto3" json:"cf,omitempty"`
	ForCas        bool                   `protobuf:"varint,4,opt,name=for_cas,json=forCas,proto3" json:"for_cas,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawBatchDeleteRequest) Reset() {
	*x = RawBatchDeleteRequest{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawBatchDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawBatchDeleteRequest) ProtoMessage() {}

func (x *RawBatchDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawBatchDeleteRequest.ProtoReflect.Descriptor instead.
func (*RawBatchDeleteRequest) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{10}
}

func (x *RawBatchDeleteRequest) GetContext() *Context {
	if x != nil {
		return x.Context
	}
	return nil
}

func (x *RawBatchDeleteRequest) GetKeys() [][]byte {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *RawBatchDeleteRequest) GetCf() string {
	if x != nil {
		return x.Cf
	}
	return ""
}

func (x *RawBatchDeleteRequest) GetForCas() bool {
	if x != nil {
		return x.ForCas
	}
	return false
}

type RawBatchDeleteResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RegionError   *Error                 `protobuf:"bytes,1,opt,name=region_error,json=regionError,proto3" json:"region_error,omitempty"`
	Error         string                 `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RawBatchDeleteResponse) Reset() {
	*x = RawBatchDeleteResponse{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RawBatchDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawBatchDeleteResponse) ProtoMessage() {}

func (x *RawBatchDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawBatchDeleteResponse.ProtoReflect.Descriptor instead.
func (*RawBatchDeleteResponse) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{11}
}

func (x *RawBatchDeleteResponse) GetRegionError() *Error {
	if x != nil {
		return x.RegionError
	}
	return nil
}

func (x *RawBatchDeleteResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// Miscellaneous metadata attached to most requests.
type Context struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	RegionId       uint64                 `protobuf:"varint,1,opt,name=region_id,json=regionId,proto3" json:"region_id,omitempty"`
	RegionEpoch    *RegionEpoch           `protobuf:"bytes,2,opt,name=region_epoch,json=regionEpoch,proto3" json:"region_epoch,omitempty"`
	Peer           *Peer                  `protobuf:"bytes,3,opt,name=peer,proto3" json:"peer,omitempty"`
	Term           uint64                 `protobuf:"varint,5,opt,name=term,proto3" json:"term,omitempty"`
	Priority       CommandPri             `protobuf:"varint,6,opt,name=priority,proto3,enum=tikv.CommandPri" json:"priority,omitempty"`
	IsolationLevel IsolationLevel         `protobuf:"varint,7,opt,name=isolation_level,json=isolationLevel,proto3,enum=tikv.IsolationLevel" json:"isolation_level,omitempty"`
	NotFillCache   bool                   `protobuf:"varint,8,opt,name=not_fill_cache,json=notFillCache,proto3" json:"not_fill_cache,omitempty"`
	SyncLog        bool                   `protobuf:"varint,9,opt,name=sync_log,json=syncLog,proto3" json:"sync_log,omitempty"`
	// True means execution time statistics should be recorded and returned.
	RecordTimeStat bool `protobuf:"varint,10,opt,name=record_time_stat,json=recordTimeStat,proto3" json:"record_time_stat,omitempty"`
	// True means RocksDB scan statistics should be recorded and returned.
	RecordScanStat bool `protobuf:"varint,11,opt,name=record_scan_stat,json=recordScanStat,proto3" json:"record_scan_stat,omitempty"`
	ReplicaRead    bool `protobuf:"varint,12,opt,name=replica_read,json=replicaRead,proto3" json:"replica_read,omitempty"`
	// Read requests can ignore locks belonging to these transactions because either
	// these transactions are rolled back or theirs commit_ts > read request's start_ts.
	ResolvedLocks          []uint64 `protobuf:"varint,13,rep,packed,name=resolved_locks,json=resolvedLocks,proto3" json:"resolved_locks,omitempty"`
	MaxExecutionDurationMs uint64   `protobuf:"varint,14,opt,name=max_execution_duration_ms,json=maxExecutionDurationMs,proto3" json:"max_execution_duration_ms,omitempty"`
	// After a region applies to `applied_index`, we can get a
	// snapshot for the region even if the peer is a follower.
	AppliedIndex uint64 `protobuf:"varint,15,opt,name=applied_index,json=appliedIndex,proto3" json:"applied_index,omitempty"`
	// A hint for TiKV to schedule tasks more fairly. Query with same task ID
	// may share same priority and resource quota.
	TaskId uint64 `protobuf:"varint,16,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	// Not required to read the most up-to-date data, replicas with `safe_ts` >= `start_ts`
	// can handle read request directly
	StaleRead bool `protobuf:"varint,17,opt,name=stale_read,json=staleRead,proto3" json:"stale_read,omitempty"`
	// Any additional serialized information about the request.
	ResourceGroupTag []byte `protobuf:"bytes,18,opt,name=resource_group_tag,json=resourceGroupTag,proto3" json:"resource_group_tag,omitempty"`
	// Used to tell TiKV whether operations are allowed or not on different disk usages.
	DiskFullOpt DiskFullOpt `protobuf:"varint,19,opt,name=disk_full_opt,json=diskFullOpt,proto3,enum=tikv.DiskFullOpt" json:"disk_full_opt,omitempty"`
	// Indicates the request is a retry request and the same request may have been sent before.
	IsRetryRequest bool `protobuf:"varint,20,opt,name=is_retry_request,json=isRetryRequest,proto3" json:"is_retry_request,omitempty"`
	// API version implies the encode of the key and value.
	ApiVersion APIVersion `protobuf:"varint,21,opt,name=api_version,json=apiVersion,proto3,enum=tikv.APIVersion" json:"api_version,omitempty"`
	// Read request should read through locks belonging to these transactions because these
	// transactions are committed and theirs commit_ts <= read request's start_ts.
	CommittedLocks []uint64 `protobuf:"varint,22,rep,packed,name=committed_locks,json=committedLocks,proto3" json:"committed_locks,omitempty"`
	// The source of the request, will be used as the tag of the metrics reporting.
	// This field can be set for any requests that require to report metrics with any extra labels.
	RequestSource string `protobuf:"bytes,24,opt,name=request_source,json=requestSource,proto3" json:"request_source,omitempty"`
	// The source of the current transaction.
	TxnSource uint64 `protobuf:"varint,25,opt,name=txn_source,json=txnSource,proto3" json:"txn_source,omitempty"`
	// If `busy_threshold_ms` is given, TiKV can reject the request and return a `ServerIsBusy`
	// error before processing if the estimated waiting duration exceeds the threshold.
	BusyThresholdMs uint32 `protobuf:"varint,27,opt,name=busy_threshold_ms,json=busyThresholdMs,proto3" json:"busy_threshold_ms,omitempty"`
	// Some information used for resource control.
	ResourceControlContext *ResourceControlContext `protobuf:"bytes,28,opt,name=resource_control_context,json=resourceControlContext,proto3" json:"resource_control_context,omitempty"`
	// The keyspace that the request is sent to.
	// NOTE: This field is only meaningful while the api_version is V2.
	KeyspaceName string `protobuf:"bytes,31,opt,name=keyspace_name,json=keyspaceName,proto3" json:"keyspace_name,omitempty"`
	// The keyspace that the request is sent to.
	// NOTE: This field is only meaningful while the api_version is V2.
	KeyspaceId uint32 `protobuf:"varint,32,opt,name=keyspace_id,json=keyspaceId,proto3" json:"keyspace_id,omitempty"`
	// The buckets version that the request is sent to.
	// NOTE: This field is only meaningful while enable buckets.
	BucketsVersion uint64 `protobuf:"varint,33,opt,name=buckets_version,json=bucketsVersion,proto3" json:"buckets_version,omitempty"`
	// It tells us where the request comes from in TiDB. If it isn't from TiDB, leave it blank.
	// This is for tests only and thus can be safely changed/removed without affecting compatibility.
	SourceStmt *SourceStmt `protobuf:"bytes,34,opt,name=source_stmt,json=sourceStmt,proto3" json:"source_stmt,omitempty"`
	// The cluster id of the request
	ClusterId uint64 `protobuf:"varint,35,opt,name=cluster_id,json=clusterId,proto3" json:"cluster_id,omitempty"`
	// The trace id of the request, will be used for tracing the request's execution's inner steps.
	TraceId []byte `protobuf:"bytes,36,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	// Control flags for trace logging behavior.
	// Bit 0: immediate_log - Force immediate logging without buffering
	// Bit 1: category_req_resp - Enable request/response tracing
	// Bit 2: category_write_details - Enable detailed write tracing
	// Bit 3: category_read_details - Enable detailed read tracing
	// Bits 4-63: Reserved for future use
	// This field is set by client-go based on an extractor function provided by TiDB.
	TraceControlFlags uint64 `protobuf:"varint,37,opt,name=trace_control_flags,json=traceControlFlags,proto3" json:"trace_control_flags,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *Context) Reset() {
	*x = Context{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Context) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Context) ProtoMessage() {}

func (x *Context) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Context.ProtoReflect.Descriptor instead.
func (*Context) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{12}
}

func (x *Context) GetRegionId() uint64 {
	if x != nil {
		return x.RegionId
	}
	return 0
}

func (x *Context) GetRegionEpoch() *RegionEpoch {
	if x != nil {
		return x.RegionEpoch
	}
	return nil
}

func (x *Context) GetPeer() *Peer {
	if x != nil {
		return x.Peer
	}
	return nil
}

func (x *Context) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *Context) GetPriority() CommandPri {
	if x != nil {
		return x.Priority
	}
	return CommandPri_Normal
}

func (x *Context) GetIsolationLevel() IsolationLevel {
	if x != nil {
		return x.IsolationLevel
	}
	return IsolationLevel_SI
}

func (x *Context) GetNotFillCache() bool {
	if x != nil {
		return x.NotFillCache
	}
	return false
}

func (x *Context) GetSyncLog() bool {
	if x != nil {
		return x.SyncLog
	}
	return false
}

func (x *Context) GetRecordTimeStat() bool {
	if x != nil {
		return x.RecordTimeStat
	}
	return false
}

func (x *Context) GetRecordScanStat() bool {
	if x != nil {
		return x.RecordScanStat
	}
	return false
}

func (x *Context) GetReplicaRead() bool {
	if x != nil {
		return x.ReplicaRead
	}
	return false
}

func (x *Context) GetResolvedLocks() []uint64 {
	if x != nil {
		return x.ResolvedLocks
	}
	return nil
}

func (x *Context) GetMaxExecutionDurationMs() uint64 {
	if x != nil {
		return x.MaxExecutionDurationMs
	}
	return 0
}

func (x *Context) GetAppliedIndex() uint64 {
	if x != nil {
		return x.AppliedIndex
	}
	return 0
}

func (x *Context) GetTaskId() uint64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *Context) GetStaleRead() bool {
	if x != nil {
		return x.StaleRead
	}
	return false
}

func (x *Context) GetResourceGroupTag() []byte {
	if x != nil {
		return x.ResourceGroupTag
	}
	return nil
}

func (x *Context) GetDiskFullOpt() DiskFullOpt {
	if x != nil {
		return x.DiskFullOpt
	}
	return DiskFullOpt_NotAllowedOnFull
}

func (x *Context) GetIsRetryRequest() bool {
	if x != nil {
		return x.IsRetryRequest
	}
	return false
}

func (x *Context) GetApiVersion() APIVersion {
	if x != nil {
		return x.ApiVersion
	}
	return APIVersion_V1
}

func (x *Context) GetCommittedLocks() []uint64 {
	if x != nil {
		return x.CommittedLocks
	}
	return nil
}

func (x *Context) GetRequestSource() string {
	if x != nil {
		return x.RequestSource
	}
	return ""
}

func (x *Context) GetTxnSource() uint64 {
	if x != nil {
		return x.TxnSource
	}
	return 0
}

func (x *Context) GetBusyThresholdMs() uint32 {
	if x != nil {
		return x.BusyThresholdMs
	}
	return 0
}

func (x *Context) GetResourceControlContext() *ResourceControlContext {
	if x != nil {
		return x.ResourceControlContext
	}
	return nil
}

func (x *Context) GetKeyspaceName() string {
	if x != nil {
		return x.KeyspaceName
	}
	return ""
}

func (x *Context) GetKeyspaceId() uint32 {
	if x != nil {
		return x.KeyspaceId
	}
	return 0
}

func (x *Context) GetBucketsVersion() uint64 {
	if x != nil {
		return x.BucketsVersion
	}
	return 0
}

func (x *Context) GetSourceStmt() *SourceStmt {
	if x != nil {
		return x.SourceStmt
	}
	return nil
}

func (x *Context) GetClusterId() uint64 {
	if x != nil {
		return x.ClusterId
	}
	return 0
}

func (x *Context) GetTraceId() []byte {
	if x != nil {
		return x.TraceId
	}
	return nil
}

func (x *Context) GetTraceControlFlags() uint64 {
	if x != nil {
		return x.TraceControlFlags
	}
	return 0
}

type ResourceControlContext struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// It's used to identify which resource group the request belongs to.
	ResourceGroupName string `protobuf:"bytes,1,opt,name=resource_group_name,json=resourceGroupName,proto3" json:"resource_group_name,omitempty"`
	// This priority would override the original priority of the resource group for the request.
	// Used to deprioritize the runaway queries.
	OverridePriority uint64 `protobuf:"varint,3,opt,name=override_priority,json=overridePriority,proto3" json:"override_priority,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *ResourceControlContext) Reset() {
	*x = ResourceControlContext{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResourceControlContext) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceControlContext) ProtoMessage() {}

func (x *ResourceControlContext) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceControlContext.ProtoReflect.Descriptor instead.
func (*ResourceControlContext) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{13}
}

func (x *ResourceControlContext) GetResourceGroupName() string {
	if x != nil {
		return x.ResourceGroupName
	}
	return ""
}

func (x *ResourceControlContext) GetOverridePriority() uint64 {
	if x != nil {
		return x.OverridePriority
	}
	return 0
}

type SourceStmt struct {
	state        protoimpl.MessageState `protogen:"open.v1"`
	StartTs      uint64                 `protobuf:"varint,1,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	ConnectionId uint64                 `protobuf:"varint,2,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
	StmtId       uint64                 `protobuf:"varint,3,opt,name=stmt_id,json=stmtId,proto3" json:"stmt_id,omitempty"`
	// session alias set by user
	SessionAlias  string `protobuf:"bytes,4,opt,name=session_alias,json=sessionAlias,proto3" json:"session_alias,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SourceStmt) Reset() {
	*x = SourceStmt{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SourceStmt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SourceStmt) ProtoMessage() {}

func (x *SourceStmt) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SourceStmt.ProtoReflect.Descriptor instead.
func (*SourceStmt) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{14}
}

func (x *SourceStmt) GetStartTs() uint64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *SourceStmt) GetConnectionId() uint64 {
	if x != nil {
		return x.ConnectionId
	}
	return 0
}

func (x *SourceStmt) GetStmtId() uint64 {
	if x != nil {
		return x.StmtId
	}
	return 0
}

func (x *SourceStmt) GetSessionAlias() string {
	if x != nil {
		return x.SessionAlias
	}
	return ""
}

type LockInfo struct {
	state       protoimpl.MessageState `protogen:"open.v1"`
	PrimaryLock []byte                 `protobuf:"bytes,1,opt,name=primary_lock,json=primaryLock,proto3" json:"primary_lock,omitempty"`
	LockVersion uint64                 `protobuf:"varint,2,opt,name=lock_version,json=lockVersion,proto3" json:"lock_version,omitempty"`
	Key         []byte                 `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	LockTtl     uint64                 `protobuf:"varint,4,opt,name=lock_ttl,json=lockTtl,proto3" json:"lock_ttl,omitempty"`
	// How many keys this transaction involves in this region.
	TxnSize         uint64 `protobuf:"varint,5,opt,name=txn_size,json=txnSize,proto3" json:"txn_size,omitempty"`
	LockType        Op     `protobuf:"varint,6,opt,name=lock_type,json=lockType,proto3,enum=tikv.Op" json:"lock_type,omitempty"`
	LockForUpdateTs uint64 `protobuf:"varint,7,opt,name=lock_for_update_ts,json=lockForUpdateTs,proto3" json:"lock_for_update_ts,omitempty"`
	// Fields for transactions that are using Async Commit.
	UseAsyncCommit bool     `protobuf:"varint,8,opt,name=use_async_commit,json=useAsyncCommit,proto3" json:"use_async_commit,omitempty"`
	MinCommitTs    uint64   `protobuf:"varint,9,opt,name=min_commit_ts,json=minCommitTs,proto3" json:"min_commit_ts,omitempty"`
	Secondaries    [][]byte `protobuf:"bytes,10,rep,name=secondaries,proto3" json:"secondaries,omitempty"`
	// The time elapsed since last update of lock wait info when waiting.
	// It's used in timeout errors. 0 means unknown or not applicable.
	// It can be used to help the client decide whether to try resolving the lock.
	DurationToLastUpdateMs uint64 `protobuf:"varint,11,opt,name=duration_to_last_update_ms,json=durationToLastUpdateMs,proto3" json:"duration_to_last_update_ms,omitempty"`
	// Reserved for file based transaction.
	IsTxnFile     bool `protobuf:"varint,100,opt,name=is_txn_file,json=isTxnFile,proto3" json:"is_txn_file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LockInfo) Reset() {
	*x = LockInfo{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LockInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LockInfo) ProtoMessage() {}

func (x *LockInfo) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LockInfo.ProtoReflect.Descriptor instead.
func (*LockInfo) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{15}
}

func (x *LockInfo) GetPrimaryLock() []byte {
	if x != nil {
		return x.PrimaryLock
	}
	return nil
}

func (x *LockInfo) GetLockVersion() uint64 {
	if x != nil {
		return x.LockVersion
	}
	return 0
}

func (x *LockInfo) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *LockInfo) GetLockTtl() uint64 {
	if x != nil {
		return x.LockTtl
	}
	return 0
}

func (x *LockInfo) GetTxnSize() uint64 {
	if x != nil {
		return x.TxnSize
	}
	return 0
}

func (x *LockInfo) GetLockType() Op {
	if x != nil {
		return x.LockType
	}
	return Op_Put
}

func (x *LockInfo) GetLockForUpdateTs() uint64 {
	if x != nil {
		return x.LockForUpdateTs
	}
	return 0
}

func (x *LockInfo) GetUseAsyncCommit() bool {
	if x != nil {
		return x.UseAsyncCommit
	}
	return false
}

func (x *LockInfo) GetMinCommitTs() uint64 {
	if x != nil {
		return x.MinCommitTs
	}
	return 0
}

func (x *LockInfo) GetSecondaries() [][]byte {
	if x != nil {
		return x.Secondaries
	}
	return nil
}

func (x *LockInfo) GetDurationToLastUpdateMs() uint64 {
	if x != nil {
		return x.DurationToLastUpdateMs
	}
	return 0
}

func (x *LockInfo) GetIsTxnFile() bool {
	if x != nil {
		return x.IsTxnFile
	}
	return false
}

type KeyError struct {
	state        protoimpl.MessageState `protogen:"open.v1"`
	Locked       *LockInfo              `protobuf:"bytes,1,opt,name=locked,proto3" json:"locked,omitempty"`                                 // Client should backoff or cleanup the lock then retry.
	Retryable    string                 `protobuf:"bytes,2,opt,name=retryable,proto3" json:"retryable,omitempty"`                           // Client may restart the txn. e.g write conflict.
	Abort        string                 `protobuf:"bytes,3,opt,name=abort,proto3" json:"abort,omitempty"`                                   // Client should abort the txn.
	Conflict     *WriteConflict         `protobuf:"bytes,4,opt,name=conflict,proto3" json:"conflict,omitempty"`                             // Write conflict is moved from retryable to here.
	AlreadyExist *AlreadyExist          `protobuf:"bytes,5,opt,name=already_exist,json=alreadyExist,proto3" json:"already_exist,omitempty"` // Key already exists
	// Deadlock deadlock = 6; // Deadlock is used in pessimistic transaction for single statement rollback.
	CommitTsExpired  *CommitTsExpired  `protobuf:"bytes,7,opt,name=commit_ts_expired,json=commitTsExpired,proto3" json:"commit_ts_expired,omitempty"`      // Commit ts is earlier than min commit ts of a transaction.
	TxnNotFound      *TxnNotFound      `protobuf:"bytes,8,opt,name=txn_not_found,json=txnNotFound,proto3" json:"txn_not_found,omitempty"`                  // Txn not found when checking txn status.
	CommitTsTooLarge *CommitTsTooLarge `protobuf:"bytes,9,opt,name=commit_ts_too_large,json=commitTsTooLarge,proto3" json:"commit_ts_too_large,omitempty"` // Calculated commit TS exceeds the limit given by the user.
	AssertionFailed  *AssertionFailed  `protobuf:"bytes,10,opt,name=assertion_failed,json=assertionFailed,proto3" json:"assertion_failed,omitempty"`       // Assertion of a `Mutation` is evaluated as a failure.
	PrimaryMismatch  *PrimaryMismatch  `protobuf:"bytes,11,opt,name=primary_mismatch,json=primaryMismatch,proto3" json:"primary_mismatch,omitempty"`       // CheckTxnStatus is sent to a lock that's not the primary.
	TxnLockNotFound  *TxnLockNotFound  `protobuf:"bytes,12,opt,name=txn_lock_not_found,json=txnLockNotFound,proto3" json:"txn_lock_not_found,omitempty"`   // TxnLockNotFound indicates the txn lock is not found.
	// Extra information for error debugging
	DebugInfo     *DebugInfo `protobuf:"bytes,100,opt,name=debug_info,json=debugInfo,proto3" json:"debug_info,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KeyError) Reset() {
	*x = KeyError{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KeyError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyError) ProtoMessage() {}

func (x *KeyError) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyError.ProtoReflect.Descriptor instead.
func (*KeyError) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{16}
}

func (x *KeyError) GetLocked() *LockInfo {
	if x != nil {
		return x.Locked
	}
	return nil
}

func (x *KeyError) GetRetryable() string {
	if x != nil {
		return x.Retryable
	}
	return ""
}

func (x *KeyError) GetAbort() string {
	if x != nil {
		return x.Abort
	}
	return ""
}

func (x *KeyError) GetConflict() *WriteConflict {
	if x != nil {
		return x.Conflict
	}
	return nil
}

func (x *KeyError) GetAlreadyExist() *AlreadyExist {
	if x != nil {
		return x.AlreadyExist
	}
	return nil
}

func (x *KeyError) GetCommitTsExpired() *CommitTsExpired {
	if x != nil {
		return x.CommitTsExpired
	}
	return nil
}

func (x *KeyError) GetTxnNotFound() *TxnNotFound {
	if x != nil {
		return x.TxnNotFound
	}
	return nil
}

func (x *KeyError) GetCommitTsTooLarge() *CommitTsTooLarge {
	if x != nil {
		return x.CommitTsTooLarge
	}
	return nil
}

func (x *KeyError) GetAssertionFailed() *AssertionFailed {
	if x != nil {
		return x.AssertionFailed
	}
	return nil
}

func (x *KeyError) GetPrimaryMismatch() *PrimaryMismatch {
	if x != nil {
		return x.PrimaryMismatch
	}
	return nil
}

func (x *KeyError) GetTxnLockNotFound() *TxnLockNotFound {
	if x != nil {
		return x.TxnLockNotFound
	}
	return nil
}

func (x *KeyError) GetDebugInfo() *DebugInfo {
	if x != nil {
		return x.DebugInfo
	}
	return nil
}

type WriteConflict struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	StartTs          uint64                 `protobuf:"varint,1,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	ConflictTs       uint64                 `protobuf:"varint,2,opt,name=conflict_ts,json=conflictTs,proto3" json:"conflict_ts,omitempty"`
	Key              []byte                 `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Primary          []byte                 `protobuf:"bytes,4,opt,name=primary,proto3" json:"primary,omitempty"`
	ConflictCommitTs uint64                 `protobuf:"varint,5,opt,name=conflict_commit_ts,json=conflictCommitTs,proto3" json:"conflict_commit_ts,omitempty"`
	Reason           WriteConflict_Reason   `protobuf:"varint,6,opt,name=reason,proto3,enum=tikv.WriteConflict_Reason" json:"reason,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *WriteConflict) Reset() {
	*x = WriteConflict{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[17]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WriteConflict) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteConflict) ProtoMessage() {}

func (x *WriteConflict) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[17]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteConflict.ProtoReflect.Descriptor instead.
func (*WriteConflict) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{17}
}

func (x *WriteConflict) GetStartTs() uint64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *WriteConflict) GetConflictTs() uint64 {
	if x != nil {
		return x.ConflictTs
	}
	return 0
}

func (x *WriteConflict) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *WriteConflict) GetPrimary() []byte {
	if x != nil {
		return x.Primary
	}
	return nil
}

func (x *WriteConflict) GetConflictCommitTs() uint64 {
	if x != nil {
		return x.ConflictCommitTs
	}
	return 0
}

func (x *WriteConflict) GetReason() WriteConflict_Reason {
	if x != nil {
		return x.Reason
	}
	return WriteConflict_Unknown
}

type AlreadyExist struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           []byte                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AlreadyExist) Reset() {
	*x = AlreadyExist{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[18]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AlreadyExist) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlreadyExist) ProtoMessage() {}

func (x *AlreadyExist) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[18]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlreadyExist.ProtoReflect.Descriptor instead.
func (*AlreadyExist) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{18}
}

func (x *AlreadyExist) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

type CommitTsExpired struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	StartTs           uint64                 `protobuf:"varint,1,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	AttemptedCommitTs uint64                 `protobuf:"varint,2,opt,name=attempted_commit_ts,json=attemptedCommitTs,proto3" json:"attempted_commit_ts,omitempty"`
	Key               []byte                 `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	MinCommitTs       uint64                 `protobuf:"varint,4,opt,name=min_commit_ts,json=minCommitTs,proto3" json:"min_commit_ts,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *CommitTsExpired) Reset() {
	*x = CommitTsExpired{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[19]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommitTsExpired) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitTsExpired) ProtoMessage() {}

func (x *CommitTsExpired) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[19]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitTsExpired.ProtoReflect.Descriptor instead.
func (*CommitTsExpired) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{19}
}

func (x *CommitTsExpired) GetStartTs() uint64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *CommitTsExpired) GetAttemptedCommitTs() uint64 {
	if x != nil {
		return x.AttemptedCommitTs
	}
	return 0
}

func (x *CommitTsExpired) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *CommitTsExpired) GetMinCommitTs() uint64 {
	if x != nil {
		return x.MinCommitTs
	}
	return 0
}

type TxnNotFound struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	StartTs       uint64                 `protobuf:"varint,1,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	PrimaryKey    []byte                 `protobuf:"bytes,2,opt,name=primary_key,json=primaryKey,proto3" json:"primary_key,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TxnNotFound) Reset() {
	*x = TxnNotFound{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[20]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TxnNotFound) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TxnNotFound) ProtoMessage() {}

func (x *TxnNotFound) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[20]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TxnNotFound.ProtoReflect.Descriptor instead.
func (*TxnNotFound) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{20}
}

func (x *TxnNotFound) GetStartTs() uint64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *TxnNotFound) GetPrimaryKey() []byte {
	if x != nil {
		return x.PrimaryKey
	}
	return nil
}

type CommitTsTooLarge struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CommitTs      uint64                 `protobuf:"varint,1,opt,name=commit_ts,json=commitTs,proto3" json:"commit_ts,omitempty"` // The calculated commit TS.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommitTsTooLarge) Reset() {
	*x = CommitTsTooLarge{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[21]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommitTsTooLarge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitTsTooLarge) ProtoMessage() {}

func (x *CommitTsTooLarge) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[21]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitTsTooLarge.ProtoReflect.Descriptor instead.
func (*CommitTsTooLarge) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{21}
}

func (x *CommitTsTooLarge) GetCommitTs() uint64 {
	if x != nil {
		return x.CommitTs
	}
	return 0
}

type AssertionFailed struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	StartTs          uint64                 `protobuf:"varint,1,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	Key              []byte                 `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Assertion        Assertion              `protobuf:"varint,3,opt,name=assertion,proto3,enum=tikv.Assertion" json:"assertion,omitempty"`
	ExistingStartTs  uint64                 `protobuf:"varint,4,opt,name=existing_start_ts,json=existingStartTs,proto3" json:"existing_start_ts,omitempty"`
	ExistingCommitTs uint64                 `protobuf:"varint,5,opt,name=existing_commit_ts,json=existingCommitTs,proto3" json:"existing_commit_ts,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *AssertionFailed) Reset() {
	*x = AssertionFailed{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[22]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AssertionFailed) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssertionFailed) ProtoMessage() {}

func (x *AssertionFailed) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[22]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssertionFailed.ProtoReflect.Descriptor instead.
func (*AssertionFailed) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{22}
}

func (x *AssertionFailed) GetStartTs() uint64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *AssertionFailed) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *AssertionFailed) GetAssertion() Assertion {
	if x != nil {
		return x.Assertion
	}
	return Assertion_None
}

func (x *AssertionFailed) GetExistingStartTs() uint64 {
	if x != nil {
		return x.ExistingStartTs
	}
	return 0
}

func (x *AssertionFailed) GetExistingCommitTs() uint64 {
	if x != nil {
		return x.ExistingCommitTs
	}
	return 0
}

type PrimaryMismatch struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	LockInfo      *LockInfo              `protobuf:"bytes,1,opt,name=lock_info,json=lockInfo,proto3" json:"lock_info,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PrimaryMismatch) Reset() {
	*x = PrimaryMismatch{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[23]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PrimaryMismatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrimaryMismatch) ProtoMessage() {}

func (x *PrimaryMismatch) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[23]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrimaryMismatch.ProtoReflect.Descriptor instead.
func (*PrimaryMismatch) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{23}
}

func (x *PrimaryMismatch) GetLockInfo() *LockInfo {
	if x != nil {
		return x.LockInfo
	}
	return nil
}

type TxnLockNotFound struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           []byte                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TxnLockNotFound) Reset() {
	*x = TxnLockNotFound{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[24]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TxnLockNotFound) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TxnLockNotFound) ProtoMessage() {}

func (x *TxnLockNotFound) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[24]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TxnLockNotFound.ProtoReflect.Descriptor instead.
func (*TxnLockNotFound) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{24}
}

func (x *TxnLockNotFound) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

type MvccDebugInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           []byte                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Mvcc          *MvccInfo              `protobuf:"bytes,2,opt,name=mvcc,proto3" json:"mvcc,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MvccDebugInfo) Reset() {
	*x = MvccDebugInfo{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[25]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MvccDebugInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MvccDebugInfo) ProtoMessage() {}

func (x *MvccDebugInfo) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[25]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MvccDebugInfo.ProtoReflect.Descriptor instead.
func (*MvccDebugInfo) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{25}
}

func (x *MvccDebugInfo) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *MvccDebugInfo) GetMvcc() *MvccInfo {
	if x != nil {
		return x.Mvcc
	}
	return nil
}

type DebugInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MvccInfo      []*MvccDebugInfo       `protobuf:"bytes,1,rep,name=mvcc_info,json=mvccInfo,proto3" json:"mvcc_info,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DebugInfo) Reset() {
	*x = DebugInfo{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[26]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DebugInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DebugInfo) ProtoMessage() {}

func (x *DebugInfo) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[26]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DebugInfo.ProtoReflect.Descriptor instead.
func (*DebugInfo) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{26}
}

func (x *DebugInfo) GetMvccInfo() []*MvccDebugInfo {
	if x != nil {
		return x.MvccInfo
	}
	return nil
}

type TimeDetail struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Off-cpu wall time elapsed in TiKV side. Usually this includes queue waiting time and
	// other kind of waitings in series. (Wait time in the raftstore is not included.)
	WaitWallTimeMs uint64 `protobuf:"varint,1,opt,name=wait_wall_time_ms,json=waitWallTimeMs,proto3" json:"wait_wall_time_ms,omitempty"`
	// Off-cpu and on-cpu wall time elapsed to actually process the request payload. It does not
	// include `wait_wall_time`.
	// This field is very close to the CPU time in most cases. Some wait time spend in RocksDB
	// cannot be excluded for now, like Mutex wait time, which is included in this field, so that
	// this field is called wall time instead of CPU time.
	ProcessWallTimeMs uint64 `protobuf:"varint,2,opt,name=process_wall_time_ms,json=processWallTimeMs,proto3" json:"process_wall_time_ms,omitempty"`
	// KV read wall Time means the time used in key/value scan and get.
	KvReadWallTimeMs uint64 `protobuf:"varint,3,opt,name=kv_read_wall_time_ms,json=kvReadWallTimeMs,proto3" json:"kv_read_wall_time_ms,omitempty"`
	// Total wall clock time spent on this RPC in TiKV .
	TotalRpcWallTimeNs uint64 `protobuf:"varint,4,opt,name=total_rpc_wall_time_ns,json=totalRpcWallTimeNs,proto3" json:"total_rpc_wall_time_ns,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *TimeDetail) Reset() {
	*x = TimeDetail{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[27]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TimeDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeDetail) ProtoMessage() {}

func (x *TimeDetail) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[27]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeDetail.ProtoReflect.Descriptor instead.
func (*TimeDetail) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{27}
}

func (x *TimeDetail) GetWaitWallTimeMs() uint64 {
	if x != nil {
		return x.WaitWallTimeMs
	}
	return 0
}

func (x *TimeDetail) GetProcessWallTimeMs() uint64 {
	if x != nil {
		return x.ProcessWallTimeMs
	}
	return 0
}

func (x *TimeDetail) GetKvReadWallTimeMs() uint64 {
	if x != nil {
		return x.KvReadWallTimeMs
	}
	return 0
}

func (x *TimeDetail) GetTotalRpcWallTimeNs() uint64 {
	if x != nil {
		return x.TotalRpcWallTimeNs
	}
	return 0
}

type TimeDetailV2 struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Off-cpu wall time elapsed in TiKV side. Usually this includes queue waiting time and
	// other kind of waitings in series. (Wait time in the raftstore is not included.)
	WaitWallTimeNs uint64 `protobuf:"varint,1,opt,name=wait_wall_time_ns,json=waitWallTimeNs,proto3" json:"wait_wall_time_ns,omitempty"`
	// Off-cpu and on-cpu wall time elapsed to actually process the request payload. It does not
	// include `wait_wall_time` and `suspend_wall_time`.
	// This field is very close to the CPU time in most cases. Some wait time spend in RocksDB
	// cannot be excluded for now, like Mutex wait time, which is included in this field, so that
	// this field is called wall time instead of CPU time.
	ProcessWallTimeNs uint64 `protobuf:"varint,2,opt,name=process_wall_time_ns,json=processWallTimeNs,proto3" json:"process_wall_time_ns,omitempty"`
	// Cpu wall time elapsed that task is waiting in queue.
	ProcessSuspendWallTimeNs uint64 `protobuf:"varint,3,opt,name=process_suspend_wall_time_ns,json=processSuspendWallTimeNs,proto3" json:"process_suspend_wall_time_ns,omitempty"`
	// KV read wall Time means the time used in key/value scan and get.
	KvReadWallTimeNs uint64 `protobuf:"varint,4,opt,name=kv_read_wall_time_ns,json=kvReadWallTimeNs,proto3" json:"kv_read_wall_time_ns,omitempty"`
	// Total wall clock time spent on this RPC in TiKV .
	TotalRpcWallTimeNs uint64 `protobuf:"varint,5,opt,name=total_rpc_wall_time_ns,json=totalRpcWallTimeNs,proto3" json:"total_rpc_wall_time_ns,omitempty"`
	// Time spent on the gRPC layer.
	KvGrpcProcessTimeNs uint64 `protobuf:"varint,6,opt,name=kv_grpc_process_time_ns,json=kvGrpcProcessTimeNs,proto3" json:"kv_grpc_process_time_ns,omitempty"`
	// Time spent on waiting for run again in grpc pool from other executor pool.
	KvGrpcWaitTimeNs uint64 `protobuf:"varint,7,opt,name=kv_grpc_wait_time_ns,json=kvGrpcWaitTimeNs,proto3" json:"kv_grpc_wait_time_ns,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *TimeDetailV2) Reset() {
	*x = TimeDetailV2{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[28]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TimeDetailV2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeDetailV2) ProtoMessage() {}

func (x *TimeDetailV2) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[28]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeDetailV2.ProtoReflect.Descriptor instead.
func (*TimeDetailV2) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{28}
}

func (x *TimeDetailV2) GetWaitWallTimeNs() uint64 {
	if x != nil {
		return x.WaitWallTimeNs
	}
	return 0
}

func (x *TimeDetailV2) GetProcessWallTimeNs() uint64 {
	if x != nil {
		return x.ProcessWallTimeNs
	}
	return 0
}

func (x *TimeDetailV2) GetProcessSuspendWallTimeNs() uint64 {
	if x != nil {
		return x.ProcessSuspendWallTimeNs
	}
	return 0
}

func (x *TimeDetailV2) GetKvReadWallTimeNs() uint64 {
	if x != nil {
		return x.KvReadWallTimeNs
	}
	return 0
}

func (x *TimeDetailV2) GetTotalRpcWallTimeNs() uint64 {
	if x != nil {
		return x.TotalRpcWallTimeNs
	}
	return 0
}

func (x *TimeDetailV2) GetKvGrpcProcessTimeNs() uint64 {
	if x != nil {
		return x.KvGrpcProcessTimeNs
	}
	return 0
}

func (x *TimeDetailV2) GetKvGrpcWaitTimeNs() uint64 {
	if x != nil {
		return x.KvGrpcWaitTimeNs
	}
	return 0
}

type ScanInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Total         int64                  `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Processed     int64                  `protobuf:"varint,2,opt,name=processed,proto3" json:"processed,omitempty"`
	ReadBytes     int64                  `protobuf:"varint,3,opt,name=read_bytes,json=readBytes,proto3" json:"read_bytes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ScanInfo) Reset() {
	*x = ScanInfo{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[29]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScanInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanInfo) ProtoMessage() {}

func (x *ScanInfo) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[29]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanInfo.ProtoReflect.Descriptor instead.
func (*ScanInfo) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{29}
}

func (x *ScanInfo) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *ScanInfo) GetProcessed() int64 {
	if x != nil {
		return x.Processed
	}
	return 0
}

func (x *ScanInfo) GetReadBytes() int64 {
	if x != nil {
		return x.ReadBytes
	}
	return 0
}

// Only reserved for compatibility.
type ScanDetail struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Write         *ScanInfo              `protobuf:"bytes,1,opt,name=write,proto3" json:"write,omitempty"`
	Lock          *ScanInfo              `protobuf:"bytes,2,opt,name=lock,proto3" json:"lock,omitempty"`
	Data          *ScanInfo              `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ScanDetail) Reset() {
	*x = ScanDetail{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[30]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScanDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanDetail) ProtoMessage() {}

func (x *ScanDetail) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[30]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanDetail.ProtoReflect.Descriptor instead.
func (*ScanDetail) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{30}
}

func (x *ScanDetail) GetWrite() *ScanInfo {
	if x != nil {
		return x.Write
	}
	return nil
}

func (x *ScanDetail) GetLock() *ScanInfo {
	if x != nil {
		return x.Lock
	}
	return nil
}

func (x *ScanDetail) GetData() *ScanInfo {
	if x != nil {
		return x.Data
	}
	return nil
}

type ScanDetailV2 struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Number of user keys scanned from the storage.
	// It does not include deleted version or RocksDB tombstone keys.
	// For Coprocessor requests, it includes keys that has been filtered out by
	// Selection.
	ProcessedVersions uint64 `protobuf:"varint,1,opt,name=processed_versions,json=processedVersions,proto3" json:"processed_versions,omitempty"`
	// Number of bytes of user key-value pairs scanned from the storage, i.e.
	// total size of data returned from MVCC layer.
	ProcessedVersionsSize uint64 `protobuf:"varint,8,opt,name=processed_versions_size,json=processedVersionsSize,proto3" json:"processed_versions_size,omitempty"`
	// Approximate number of MVCC keys meet during scanning. It includes
	// deleted versions, but does not include RocksDB tombstone keys.
	//
	// When this field is notably larger than `processed_versions`, it means
	// there are a lot of deleted MVCC keys.
	TotalVersions uint64 `protobuf:"varint,2,opt,name=total_versions,json=totalVersions,proto3" json:"total_versions,omitempty"`
	// Total number of deletes and single deletes skipped over during
	// iteration, i.e. how many RocksDB tombstones are skipped.
	RocksdbDeleteSkippedCount uint64 `protobuf:"varint,3,opt,name=rocksdb_delete_skipped_count,json=rocksdbDeleteSkippedCount,proto3" json:"rocksdb_delete_skipped_count,omitempty"`
	// Total number of internal keys skipped over during iteration.
	// See https://github.com/facebook/rocksdb/blob/9f1c84ca471d8b1ad7be9f3eebfc2c7e07dfd7a7/include/rocksdb/perf_context.h#L84 for details.
	RocksdbKeySkippedCount uint64 `protobuf:"varint,4,opt,name=rocksdb_key_skipped_count,json=rocksdbKeySkippedCount,proto3" json:"rocksdb_key_skipped_count,omitempty"`
	// Total number of RocksDB block cache hits.
	RocksdbBlockCacheHitCount uint64 `protobuf:"varint,5,opt,name=rocksdb_block_cache_hit_count,json=rocksdbBlockCacheHitCount,proto3" json:"rocksdb_block_cache_hit_count,omitempty"`
	// Total number of block reads (with IO).
	RocksdbBlockReadCount uint64 `protobuf:"varint,6,opt,name=rocksdb_block_read_count,json=rocksdbBlockReadCount,proto3" json:"rocksdb_block_read_count,omitempty"`
	// Total number of bytes from block reads.
	RocksdbBlockReadByte uint64 `protobuf:"varint,7,opt,name=rocksdb_block_read_byte,json=rocksdbBlockReadByte,proto3" json:"rocksdb_block_read_byte,omitempty"`
	// Total time used for block reads.
	RocksdbBlockReadNanos uint64 `protobuf:"varint,9,opt,name=rocksdb_block_read_nanos,json=rocksdbBlockReadNanos,proto3" json:"rocksdb_block_read_nanos,omitempty"`
	// Time used for getting a raftstore snapshot (including proposing read index, leader confirmation and getting the RocksDB snapshot).
	GetSnapshotNanos uint64 `protobuf:"varint,10,opt,name=get_snapshot_nanos,json=getSnapshotNanos,proto3" json:"get_snapshot_nanos,omitempty"`
	// Time used for proposing read index from read pool to store pool, equals 0 when performing lease read.
	ReadIndexProposeWaitNanos uint64 `protobuf:"varint,11,opt,name=read_index_propose_wait_nanos,json=readIndexProposeWaitNanos,proto3" json:"read_index_propose_wait_nanos,omitempty"`
	// Time used for leader confirmation, equals 0 when performing lease read.
	ReadIndexConfirmWaitNanos uint64 `protobuf:"varint,12,opt,name=read_index_confirm_wait_nanos,json=readIndexConfirmWaitNanos,proto3" json:"read_index_confirm_wait_nanos,omitempty"`
	// Time used for read pool scheduling.
	ReadPoolScheduleWaitNanos uint64 `protobuf:"varint,13,opt,name=read_pool_schedule_wait_nanos,json=readPoolScheduleWaitNanos,proto3" json:"read_pool_schedule_wait_nanos,omitempty"`
	unknownFields             protoimpl.UnknownFields
	sizeCache                 protoimpl.SizeCache
}

func (x *ScanDetailV2) Reset() {
	*x = ScanDetailV2{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[31]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ScanDetailV2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanDetailV2) ProtoMessage() {}

func (x *ScanDetailV2) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[31]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanDetailV2.ProtoReflect.Descriptor instead.
func (*ScanDetailV2) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{31}
}

func (x *ScanDetailV2) GetProcessedVersions() uint64 {
	if x != nil {
		return x.ProcessedVersions
	}
	return 0
}

func (x *ScanDetailV2) GetProcessedVersionsSize() uint64 {
	if x != nil {
		return x.ProcessedVersionsSize
	}
	return 0
}

func (x *ScanDetailV2) GetTotalVersions() uint64 {
	if x != nil {
		return x.TotalVersions
	}
	return 0
}

func (x *ScanDetailV2) GetRocksdbDeleteSkippedCount() uint64 {
	if x != nil {
		return x.RocksdbDeleteSkippedCount
	}
	return 0
}

func (x *ScanDetailV2) GetRocksdbKeySkippedCount() uint64 {
	if x != nil {
		return x.RocksdbKeySkippedCount
	}
	return 0
}

func (x *ScanDetailV2) GetRocksdbBlockCacheHitCount() uint64 {
	if x != nil {
		return x.RocksdbBlockCacheHitCount
	}
	return 0
}

func (x *ScanDetailV2) GetRocksdbBlockReadCount() uint64 {
	if x != nil {
		return x.RocksdbBlockReadCount
	}
	return 0
}

func (x *ScanDetailV2) GetRocksdbBlockReadByte() uint64 {
	if x != nil {
		return x.RocksdbBlockReadByte
	}
	return 0
}

func (x *ScanDetailV2) GetRocksdbBlockReadNanos() uint64 {
	if x != nil {
		return x.RocksdbBlockReadNanos
	}
	return 0
}

func (x *ScanDetailV2) GetGetSnapshotNanos() uint64 {
	if x != nil {
		return x.GetSnapshotNanos
	}
	return 0
}

func (x *ScanDetailV2) GetReadIndexProposeWaitNanos() uint64 {
	if x != nil {
		return x.ReadIndexProposeWaitNanos
	}
	return 0
}

func (x *ScanDetailV2) GetReadIndexConfirmWaitNanos() uint64 {
	if x != nil {
		return x.ReadIndexConfirmWaitNanos
	}
	return 0
}

func (x *ScanDetailV2) GetReadPoolScheduleWaitNanos() uint64 {
	if x != nil {
		return x.ReadPoolScheduleWaitNanos
	}
	return 0
}

type ExecDetails struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Available when ctx.record_time_stat = true or meet slow query.
	TimeDetail *TimeDetail `protobuf:"bytes,1,opt,name=time_detail,json=timeDetail,proto3" json:"time_detail,omitempty"`
	// Available when ctx.record_scan_stat = true or meet slow query.
	ScanDetail    *ScanDetail `protobuf:"bytes,2,opt,name=scan_detail,json=scanDetail,proto3" json:"scan_detail,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExecDetails) Reset() {
	*x = ExecDetails{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[32]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExecDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecDetails) ProtoMessage() {}

func (x *ExecDetails) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[32]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecDetails.ProtoReflect.Descriptor instead.
func (*ExecDetails) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{32}
}

func (x *ExecDetails) GetTimeDetail() *TimeDetail {
	if x != nil {
		return x.TimeDetail
	}
	return nil
}

func (x *ExecDetails) GetScanDetail() *ScanDetail {
	if x != nil {
		return x.ScanDetail
	}
	return nil
}

type ExecDetailsV2 struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Available when ctx.record_time_stat = true or meet slow query.
	// deprecated. Should use `time_detail_v2` instead.
	TimeDetail *TimeDetail `protobuf:"bytes,1,opt,name=time_detail,json=timeDetail,proto3" json:"time_detail,omitempty"`
	// Available when ctx.record_scan_stat = true or meet slow query.
	ScanDetailV2 *ScanDetailV2 `protobuf:"bytes,2,opt,name=scan_detail_v2,json=scanDetailV2,proto3" json:"scan_detail_v2,omitempty"`
	// Raftstore writing durations of the request. Only available for some write requests.
	WriteDetail *WriteDetail `protobuf:"bytes,3,opt,name=write_detail,json=writeDetail,proto3" json:"write_detail,omitempty"`
	// Available when ctx.record_time_stat = true or meet slow query.
	TimeDetailV2  *TimeDetailV2 `protobuf:"bytes,4,opt,name=time_detail_v2,json=timeDetailV2,proto3" json:"time_detail_v2,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExecDetailsV2) Reset() {
	*x = ExecDetailsV2{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[33]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExecDetailsV2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecDetailsV2) ProtoMessage() {}

func (x *ExecDetailsV2) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[33]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecDetailsV2.ProtoReflect.Descriptor instead.
func (*ExecDetailsV2) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{33}
}

func (x *ExecDetailsV2) GetTimeDetail() *TimeDetail {
	if x != nil {
		return x.TimeDetail
	}
	return nil
}

func (x *ExecDetailsV2) GetScanDetailV2() *ScanDetailV2 {
	if x != nil {
		return x.ScanDetailV2
	}
	return nil
}

func (x *ExecDetailsV2) GetWriteDetail() *WriteDetail {
	if x != nil {
		return x.WriteDetail
	}
	return nil
}

func (x *ExecDetailsV2) GetTimeDetailV2() *TimeDetailV2 {
	if x != nil {
		return x.TimeDetailV2
	}
	return nil
}

type WriteDetail struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Wait duration in the store loop.
	StoreBatchWaitNanos uint64 `protobuf:"varint,1,opt,name=store_batch_wait_nanos,json=storeBatchWaitNanos,proto3" json:"store_batch_wait_nanos,omitempty"`
	// Wait duration before sending proposal to peers.
	ProposeSendWaitNanos uint64 `protobuf:"varint,2,opt,name=propose_send_wait_nanos,json=proposeSendWaitNanos,proto3" json:"propose_send_wait_nanos,omitempty"`
	// Total time spent on persisting the log.
	PersistLogNanos uint64 `protobuf:"varint,3,opt,name=persist_log_nanos,json=persistLogNanos,proto3" json:"persist_log_nanos,omitempty"`
	// Wait time until the Raft log write leader begins to write.
	RaftDbWriteLeaderWaitNanos uint64 `protobuf:"varint,4,opt,name=raft_db_write_leader_wait_nanos,json=raftDbWriteLeaderWaitNanos,proto3" json:"raft_db_write_leader_wait_nanos,omitempty"`
	// Time spent on synchronizing the Raft log to the disk.
	RaftDbSyncLogNanos uint64 `protobuf:"varint,5,opt,name=raft_db_sync_log_nanos,json=raftDbSyncLogNanos,proto3" json:"raft_db_sync_log_nanos,omitempty"`
	// Time spent on writing the Raft log to the Raft memtable.
	RaftDbWriteMemtableNanos uint64 `protobuf:"varint,6,opt,name=raft_db_write_memtable_nanos,json=raftDbWriteMemtableNanos,proto3" json:"raft_db_write_memtable_nanos,omitempty"`
	// Time waiting for peers to confirm the proposal (counting from the instant when the leader sends the proposal message).
	CommitLogNanos uint64 `protobuf:"varint,7,opt,name=commit_log_nanos,json=commitLogNanos,proto3" json:"commit_log_nanos,omitempty"`
	// Wait duration in the apply loop.
	ApplyBatchWaitNanos uint64 `protobuf:"varint,8,opt,name=apply_batch_wait_nanos,json=applyBatchWaitNanos,proto3" json:"apply_batch_wait_nanos,omitempty"`
	// Total time spend to applying the log.
	ApplyLogNanos uint64 `protobuf:"varint,9,opt,name=apply_log_nanos,json=applyLogNanos,proto3" json:"apply_log_nanos,omitempty"`
	// Wait time until the KV RocksDB lock is acquired.
	ApplyMutexLockNanos uint64 `protobuf:"varint,10,opt,name=apply_mutex_lock_nanos,json=applyMutexLockNanos,proto3" json:"apply_mutex_lock_nanos,omitempty"`
	// Wait time until becoming the KV RocksDB write leader.
	ApplyWriteLeaderWaitNanos uint64 `protobuf:"varint,11,opt,name=apply_write_leader_wait_nanos,json=applyWriteLeaderWaitNanos,proto3" json:"apply_write_leader_wait_nanos,omitempty"`
	// Time spent on writing the KV DB WAL to the disk.
	ApplyWriteWalNanos uint64 `protobuf:"varint,12,opt,name=apply_write_wal_nanos,json=applyWriteWalNanos,proto3" json:"apply_write_wal_nanos,omitempty"`
	// Time spent on writing to the memtable of the KV RocksDB.
	ApplyWriteMemtableNanos uint64 `protobuf:"varint,13,opt,name=apply_write_memtable_nanos,json=applyWriteMemtableNanos,proto3" json:"apply_write_memtable_nanos,omitempty"`
	// Time spent on waiting in the latch.
	LatchWaitNanos uint64 `protobuf:"varint,14,opt,name=latch_wait_nanos,json=latchWaitNanos,proto3" json:"latch_wait_nanos,omitempty"`
	// Processing time in the transaction layer.
	ProcessNanos uint64 `protobuf:"varint,15,opt,name=process_nanos,json=processNanos,proto3" json:"process_nanos,omitempty"`
	// Wait time because of the scheduler flow control or quota limiter throttling.
	ThrottleNanos uint64 `protobuf:"varint,16,opt,name=throttle_nanos,json=throttleNanos,proto3" json:"throttle_nanos,omitempty"`
	// Wait time in the waiter manager for pessimistic locking.
	PessimisticLockWaitNanos uint64 `protobuf:"varint,17,opt,name=pessimistic_lock_wait_nanos,json=pessimisticLockWaitNanos,proto3" json:"pessimistic_lock_wait_nanos,omitempty"`
	unknownFields            protoimpl.UnknownFields
	sizeCache                protoimpl.SizeCache
}

func (x *WriteDetail) Reset() {
	*x = WriteDetail{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[34]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WriteDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteDetail) ProtoMessage() {}

func (x *WriteDetail) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[34]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteDetail.ProtoReflect.Descriptor instead.
func (*WriteDetail) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{34}
}

func (x *WriteDetail) GetStoreBatchWaitNanos() uint64 {
	if x != nil {
		return x.StoreBatchWaitNanos
	}
	return 0
}

func (x *WriteDetail) GetProposeSendWaitNanos() uint64 {
	if x != nil {
		return x.ProposeSendWaitNanos
	}
	return 0
}

func (x *WriteDetail) GetPersistLogNanos() uint64 {
	if x != nil {
		return x.PersistLogNanos
	}
	return 0
}

func (x *WriteDetail) GetRaftDbWriteLeaderWaitNanos() uint64 {
	if x != nil {
		return x.RaftDbWriteLeaderWaitNanos
	}
	return 0
}

func (x *WriteDetail) GetRaftDbSyncLogNanos() uint64 {
	if x != nil {
		return x.RaftDbSyncLogNanos
	}
	return 0
}

func (x *WriteDetail) GetRaftDbWriteMemtableNanos() uint64 {
	if x != nil {
		return x.RaftDbWriteMemtableNanos
	}
	return 0
}

func (x *WriteDetail) GetCommitLogNanos() uint64 {
	if x != nil {
		return x.CommitLogNanos
	}
	return 0
}

func (x *WriteDetail) GetApplyBatchWaitNanos() uint64 {
	if x != nil {
		return x.ApplyBatchWaitNanos
	}
	return 0
}

func (x *WriteDetail) GetApplyLogNanos() uint64 {
	if x != nil {
		return x.ApplyLogNanos
	}
	return 0
}

func (x *WriteDetail) GetApplyMutexLockNanos() uint64 {
	if x != nil {
		return x.ApplyMutexLockNanos
	}
	return 0
}

func (x *WriteDetail) GetApplyWriteLeaderWaitNanos() uint64 {
	if x != nil {
		return x.ApplyWriteLeaderWaitNanos
	}
	return 0
}

func (x *WriteDetail) GetApplyWriteWalNanos() uint64 {
	if x != nil {
		return x.ApplyWriteWalNanos
	}
	return 0
}

func (x *WriteDetail) GetApplyWriteMemtableNanos() uint64 {
	if x != nil {
		return x.ApplyWriteMemtableNanos
	}
	return 0
}

func (x *WriteDetail) GetLatchWaitNanos() uint64 {
	if x != nil {
		return x.LatchWaitNanos
	}
	return 0
}

func (x *WriteDetail) GetProcessNanos() uint64 {
	if x != nil {
		return x.ProcessNanos
	}
	return 0
}

func (x *WriteDetail) GetThrottleNanos() uint64 {
	if x != nil {
		return x.ThrottleNanos
	}
	return 0
}

func (x *WriteDetail) GetPessimisticLockWaitNanos() uint64 {
	if x != nil {
		return x.PessimisticLockWaitNanos
	}
	return 0
}

type KvPair struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Error *KeyError              `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Key   []byte                 `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte                 `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	// The commit timestamp of the key.
	// If it is zero, it means the commit timestamp is unknown.
	CommitTs      uint64 `protobuf:"varint,4,opt,name=commit_ts,json=commitTs,proto3" json:"commit_ts,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KvPair) Reset() {
	*x = KvPair{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[35]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KvPair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KvPair) ProtoMessage() {}

func (x *KvPair) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[35]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KvPair.ProtoReflect.Descriptor instead.
func (*KvPair) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{35}
}

func (x *KvPair) GetError() *KeyError {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *KvPair) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *KvPair) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *KvPair) GetCommitTs() uint64 {
	if x != nil {
		return x.CommitTs
	}
	return 0
}

type MvccWrite struct {
	state                 protoimpl.MessageState `protogen:"open.v1"`
	Type                  Op                     `protobuf:"varint,1,opt,name=type,proto3,enum=tikv.Op" json:"type,omitempty"`
	StartTs               uint64                 `protobuf:"varint,2,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	CommitTs              uint64                 `protobuf:"varint,3,opt,name=commit_ts,json=commitTs,proto3" json:"commit_ts,omitempty"`
	ShortValue            []byte                 `protobuf:"bytes,4,opt,name=short_value,json=shortValue,proto3" json:"short_value,omitempty"`
	HasOverlappedRollback bool                   `protobuf:"varint,5,opt,name=has_overlapped_rollback,json=hasOverlappedRollback,proto3" json:"has_overlapped_rollback,omitempty"`
	HasGcFence            bool                   `protobuf:"varint,6,opt,name=has_gc_fence,json=hasGcFence,proto3" json:"has_gc_fence,omitempty"`
	GcFence               uint64                 `protobuf:"varint,7,opt,name=gc_fence,json=gcFence,proto3" json:"gc_fence,omitempty"`
	LastChangeTs          uint64                 `protobuf:"varint,8,opt,name=last_change_ts,json=lastChangeTs,proto3" json:"last_change_ts,omitempty"`
	VersionsToLastChange  uint64                 `protobuf:"varint,9,opt,name=versions_to_last_change,json=versionsToLastChange,proto3" json:"versions_to_last_change,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *MvccWrite) Reset() {
	*x = MvccWrite{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[36]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MvccWrite) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MvccWrite) ProtoMessage() {}

func (x *MvccWrite) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[36]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MvccWrite.ProtoReflect.Descriptor instead.
func (*MvccWrite) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{36}
}

func (x *MvccWrite) GetType() Op {
	if x != nil {
		return x.Type
	}
	return Op_Put
}

func (x *MvccWrite) GetStartTs() uint64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *MvccWrite) GetCommitTs() uint64 {
	if x != nil {
		return x.CommitTs
	}
	return 0
}

func (x *MvccWrite) GetShortValue() []byte {
	if x != nil {
		return x.ShortValue
	}
	return nil
}

func (x *MvccWrite) GetHasOverlappedRollback() bool {
	if x != nil {
		return x.HasOverlappedRollback
	}
	return false
}

func (x *MvccWrite) GetHasGcFence() bool {
	if x != nil {
		return x.HasGcFence
	}
	return false
}

func (x *MvccWrite) GetGcFence() uint64 {
	if x != nil {
		return x.GcFence
	}
	return 0
}

func (x *MvccWrite) GetLastChangeTs() uint64 {
	if x != nil {
		return x.LastChangeTs
	}
	return 0
}

func (x *MvccWrite) GetVersionsToLastChange() uint64 {
	if x != nil {
		return x.VersionsToLastChange
	}
	return 0
}

type MvccValue struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	StartTs       uint64                 `protobuf:"varint,1,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	Value         []byte                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MvccValue) Reset() {
	*x = MvccValue{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[37]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MvccValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MvccValue) ProtoMessage() {}

func (x *MvccValue) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[37]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MvccValue.ProtoReflect.Descriptor instead.
func (*MvccValue) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{37}
}

func (x *MvccValue) GetStartTs() uint64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *MvccValue) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type MvccLock struct {
	state                protoimpl.MessageState `protogen:"open.v1"`
	Type                 Op                     `protobuf:"varint,1,opt,name=type,proto3,enum=tikv.Op" json:"type,omitempty"`
	StartTs              uint64                 `protobuf:"varint,2,opt,name=start_ts,json=startTs,proto3" json:"start_ts,omitempty"`
	Primary              []byte                 `protobuf:"bytes,3,opt,name=primary,proto3" json:"primary,omitempty"`
	ShortValue           []byte                 `protobuf:"bytes,4,opt,name=short_value,json=shortValue,proto3" json:"short_value,omitempty"`
	Ttl                  uint64                 `protobuf:"varint,5,opt,name=ttl,proto3" json:"ttl,omitempty"`
	ForUpdateTs          uint64                 `protobuf:"varint,6,opt,name=for_update_ts,json=forUpdateTs,proto3" json:"for_update_ts,omitempty"`
	TxnSize              uint64                 `protobuf:"varint,7,opt,name=txn_size,json=txnSize,proto3" json:"txn_size,omitempty"`
	UseAsyncCommit       bool                   `protobuf:"varint,8,opt,name=use_async_commit,json=useAsyncCommit,proto3" json:"use_async_commit,omitempty"`
	Secondaries          [][]byte               `protobuf:"bytes,9,rep,name=secondaries,proto3" json:"secondaries,omitempty"`
	RollbackTs           []uint64               `protobuf:"varint,10,rep,packed,name=rollback_ts,json=rollbackTs,proto3" json:"rollback_ts,omitempty"`
	LastChangeTs         uint64                 `protobuf:"varint,11,opt,name=last_change_ts,json=lastChangeTs,proto3" json:"last_change_ts,omitempty"`
	VersionsToLastChange uint64                 `protobuf:"varint,12,opt,name=versions_to_last_change,json=versionsToLastChange,proto3" json:"versions_to_last_change,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *MvccLock) Reset() {
	*x = MvccLock{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[38]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MvccLock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MvccLock) ProtoMessage() {}

func (x *MvccLock) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[38]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MvccLock.ProtoReflect.Descriptor instead.
func (*MvccLock) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{38}
}

func (x *MvccLock) GetType() Op {
	if x != nil {
		return x.Type
	}
	return Op_Put
}

func (x *MvccLock) GetStartTs() uint64 {
	if x != nil {
		return x.StartTs
	}
	return 0
}

func (x *MvccLock) GetPrimary() []byte {
	if x != nil {
		return x.Primary
	}
	return nil
}

func (x *MvccLock) GetShortValue() []byte {
	if x != nil {
		return x.ShortValue
	}
	return nil
}

func (x *MvccLock) GetTtl() uint64 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

func (x *MvccLock) GetForUpdateTs() uint64 {
	if x != nil {
		return x.ForUpdateTs
	}
	return 0
}

func (x *MvccLock) GetTxnSize() uint64 {
	if x != nil {
		return x.TxnSize
	}
	return 0
}

func (x *MvccLock) GetUseAsyncCommit() bool {
	if x != nil {
		return x.UseAsyncCommit
	}
	return false
}

func (x *MvccLock) GetSecondaries() [][]byte {
	if x != nil {
		return x.Secondaries
	}
	return nil
}

func (x *MvccLock) GetRollbackTs() []uint64 {
	if x != nil {
		return x.RollbackTs
	}
	return nil
}

func (x *MvccLock) GetLastChangeTs() uint64 {
	if x != nil {
		return x.LastChangeTs
	}
	return 0
}

func (x *MvccLock) GetVersionsToLastChange() uint64 {
	if x != nil {
		return x.VersionsToLastChange
	}
	return 0
}

type MvccInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Lock          *MvccLock              `protobuf:"bytes,1,opt,name=lock,proto3" json:"lock,omitempty"`
	Writes        []*MvccWrite           `protobuf:"bytes,2,rep,name=writes,proto3" json:"writes,omitempty"`
	Values        []*MvccValue           `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MvccInfo) Reset() {
	*x = MvccInfo{}
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[39]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MvccInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MvccInfo) ProtoMessage() {}

func (x *MvccInfo) ProtoReflect() protoreflect.Message {
	mi := &file_v1_tikv_kvrpcpb_proto_msgTypes[39]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MvccInfo.ProtoReflect.Descriptor instead.
func (*MvccInfo) Descriptor() ([]byte, []int) {
	return file_v1_tikv_kvrpcpb_proto_rawDescGZIP(), []int{39}
}

func (x *MvccInfo) GetLock() *MvccLock {
	if x != nil {
		return x.Lock
	}
	return nil
}

func (x *MvccInfo) GetWrites() []*MvccWrite {
	if x != nil {
		return x.Writes
	}
	return nil
}

func (x *MvccInfo) GetValues() []*MvccValue {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_v1_tikv_kvrpcpb_proto protoreflect.FileDescriptor

const file_v1_tikv_kvrpcpb_proto_rawDesc = "" +
	"\n" +
	"\x15v1/tikv/kvrpcpb.proto\x12\x04tikv\x1a\x15v1/tikv/errorpb.proto\x1a\x14v1/tikv/metapb.proto\"Z\n" +
	"\rRawGetRequest\x12'\n" +
	"\acontext\x18\x01 \x01(\v2\r.tikv.ContextR\acontext\x12\x10\n" +
	"\x03key\x18\x02 \x01(\fR\x03key\x12\x0e\n" +
	"\x02cf\x18\x03 \x01(\tR\x02cf\"\x89\x01\n" +
	"\x0eRawGetResponse\x12.\n" +
	"\fregion_error\x18\x01 \x01(\v2\v.tikv.ErrorR\vregionError\x12\x14\n" +
	"\x05error\x18\x02 \x01(\tR\x05error\x12\x14\n" +
	"\x05value\x18\x03 \x01(\fR\x05value\x12\x1b\n" +
	"\tnot_found\x18\x04 \x01(\bR\bnotFound\"a\n" +
	"\x12RawBatchGetRequest\x12'\n" +
	"\acontext\x18\x01 \x01(\v2\r.tikv.ContextR\acontext\x12\x12\n" +
	"\x04keys\x18\x02 \x03(\fR\x04keys\x12\x0e\n" +
	"\x02cf\x18\x03 \x01(\tR\x02cf\"i\n" +
	"\x13RawBatchGetResponse\x12.\n" +
	"\fregion_error\x18\x01 \x01(\v2\v.tikv.ErrorR\vregionError\x12\"\n" +
	"\x05pairs\x18\x02 \x03(\v2\f.tikv.KvPairR\x05pairs\"\x9b\x01\n" +
	"\rRawPutRequest\x12'\n" +
	"\acontext\x18\x01 \x01(\v2\r.tikv.ContextR\acontext\x12\x10\n" +
	"\x03key\x18\x02 \x01(\fR\x03key\x12\x14\n" +
	"\x05value\x18\x03 \x01(\fR\x05value\x12\x0e\n" +
	"\x02cf\x18\x04 \x01(\tR\x02cf\x12\x10\n" +
	"\x03ttl\x18\x05 \x01(\x04R\x03ttl\x12\x17\n" +
	"\afor_cas\x18\x06 \x01(\bR\x06forCas\"V\n" +
	"\x0eRawPutResponse\x12.\n" +
	"\fregion_error\x18\x01 \x01(\v2\v.tikv.ErrorR\vregionError\x12\x14\n" +
	"\x05error\x18\x02 \x01(\tR\x05error\"\xb4\x01\n" +
	"\x12RawBatchPutRequest\x12'\n" +
	"\acontext\x18\x01 \x01(\v2\r.tikv.ContextR\acontext\x12\"\n" +
	"\x05pairs\x18\x02 \x03(\v2\f.tikv.KvPairR\x05pairs\x12\x0e\n" +
	"\x02cf\x18\x03 \x01(\tR\x02cf\x12\x14\n" +
	"\x03ttl\x18\x04 \x01(\x04B\x02\x18\x01R\x03ttl\x12\x17\n" +
	"\afor_cas\x18\x05 \x01(\bR\x06forCas\x12\x12\n" +
	"\x04ttls\x18\x06 \x03(\x04R\x04ttls\"[\n" +
	"\x13RawBatchPutResponse\x12.\n" +
	"\fregion_error\x18\x01 \x01(\v2\v.tikv.ErrorR\vregionError\x12\x14\n" +
	"\x05error\x18\x02 \x01(\tR\x05error\"v\n" +
	"\x10RawDeleteRequest\x12'\n" +
	"\acontext\x18\x01 \x01(\v2\r.tikv.ContextR\acontext\x12\x10\n" +
	"\x03key\x18\x02 \x01(\fR\x03key\x12\x0e\n" +
	"\x02cf\x18\x03 \x01(\tR\x02cf\x12\x17\n" +
	"\afor_cas\x18\x04 \x01(\bR\x06forCas\"Y\n" +
	"\x11RawDeleteResponse\x12.\n" +
	"\fregion_error\x18\x01 \x01(\v2\v.tikv.ErrorR\vregionError\x12\x14\n" +
	"\x05error\x18\x02 \x01(\tR\x05error\"}\n" +
	"\x15RawBatchDeleteRequest\x12'\n" +
	"\acontext\x18\x01 \x01(\v2\r.tikv.ContextR\acontext\x12\x12\n" +
	"\x04keys\x18\x02 \x03(\fR\x04keys\x12\x0e\n" +
	"\x02cf\x18\x03 \x01(\tR\x02cf\x12\x17\n" +
	"\afor_cas\x18\x04 \x01(\bR\x06forCas\"^\n" +
	"\x16RawBatchDeleteResponse\x12.\n" +
	"\fregion_error\x18\x01 \x01(\v2\v.tikv.ErrorR\vregionError\x12\x14\n" +
	"\x05error\x18\x02 \x01(\tR\x05error\"\xd2\n" +
	"\n" +
	"\aContext\x12\x1b\n" +
	"\tregion_id\x18\x01 \x01(\x04R\bregionId\x126\n" +
	"\fregion_epoch\x18\x02 \x01(\v2\x13.metapb.RegionEpochR\vregionEpoch\x12 \n" +
	"\x04peer\x18\x03 \x01(\v2\f.metapb.PeerR\x04peer\x12\x12\n" +
	"\x04term\x18\x05 \x01(\x04R\x04term\x12,\n" +
	"\bpriority\x18\x06 \x01(\x0e2\x10.tikv.CommandPriR\bpriority\x12=\n" +
	"\x0fisolation_level\x18\a \x01(\x0e2\x14.tikv.IsolationLevelR\x0eisolationLevel\x12$\n" +
	"\x0enot_fill_cache\x18\b \x01(\bR\fnotFillCache\x12\x19\n" +
	"\bsync_log\x18\t \x01(\bR\asyncLog\x12(\n" +
	"\x10record_time_stat\x18\n" +
	" \x01(\bR\x0erecordTimeStat\x12(\n" +
	"\x10record_scan_stat\x18\v \x01(\bR\x0erecordScanStat\x12!\n" +
	"\freplica_read\x18\f \x01(\bR\vreplicaRead\x12%\n" +
	"\x0eresolved_locks\x18\r \x03(\x04R\rresolvedLocks\x129\n" +
	"\x19max_execution_duration_ms\x18\x0e \x01(\x04R\x16maxExecutionDurationMs\x12#\n" +
	"\rapplied_index\x18\x0f \x01(\x04R\fappliedIndex\x12\x17\n" +
	"\atask_id\x18\x10 \x01(\x04R\x06taskId\x12\x1d\n" +
	"\n" +
	"stale_read\x18\x11 \x01(\bR\tstaleRead\x12,\n" +
	"\x12resource_group_tag\x18\x12 \x01(\fR\x10resourceGroupTag\x125\n" +
	"\rdisk_full_opt\x18\x13 \x01(\x0e2\x11.tikv.DiskFullOptR\vdiskFullOpt\x12(\n" +
	"\x10is_retry_request\x18\x14 \x01(\bR\x0eisRetryRequest\x121\n" +
	"\vapi_version\x18\x15 \x01(\x0e2\x10.tikv.APIVersionR\n" +
	"apiVersion\x12'\n" +
	"\x0fcommitted_locks\x18\x16 \x03(\x04R\x0ecommittedLocks\x12%\n" +
	"\x0erequest_source\x18\x18 \x01(\tR\rrequestSource\x12\x1d\n" +
	"\n" +
	"txn_source\x18\x19 \x01(\x04R\ttxnSource\x12*\n" +
	"\x11busy_threshold_ms\x18\x1b \x01(\rR\x0fbusyThresholdMs\x12V\n" +
	"\x18resource_control_context\x18\x1c \x01(\v2\x1c.tikv.ResourceControlContextR\x16resourceControlContext\x12#\n" +
	"\rkeyspace_name\x18\x1f \x01(\tR\fkeyspaceName\x12\x1f\n" +
	"\vkeyspace_id\x18  \x01(\rR\n" +
	"keyspaceId\x12'\n" +
	"\x0fbuckets_version\x18! \x01(\x04R\x0ebucketsVersion\x121\n" +
	"\vsource_stmt\x18\" \x01(\v2\x10.tikv.SourceStmtR\n" +
	"sourceStmt\x12\x1d\n" +
	"\n" +
	"cluster_id\x18# \x01(\x04R\tclusterId\x12\x19\n" +
	"\btrace_id\x18$ \x01(\fR\atraceId\x12.\n" +
	"\x13trace_control_flags\x18% \x01(\x04R\x11traceControlFlagsJ\x04\b\x04\x10\x05J\x04\b\x1a\x10\x1bR\vread_quorum\"u\n" +
	"\x16ResourceControlContext\x12.\n" +
	"\x13resource_group_name\x18\x01 \x01(\tR\x11resourceGroupName\x12+\n" +
	"\x11override_priority\x18\x03 \x01(\x04R\x10overridePriority\"\x8a\x01\n" +
	"\n" +
	"SourceStmt\x12\x19\n" +
	"\bstart_ts\x18\x01 \x01(\x04R\astartTs\x12#\n" +
	"\rconnection_id\x18\x02 \x01(\x04R\fconnectionId\x12\x17\n" +
	"\astmt_id\x18\x03 \x01(\x04R\x06stmtId\x12#\n" +
	"\rsession_alias\x18\x04 \x01(\tR\fsessionAlias\"\xb8\x03\n" +
	"\bLockInfo\x12!\n" +
	"\fprimary_lock\x18\x01 \x01(\fR\vprimaryLock\x12!\n" +
	"\flock_version\x18\x02 \x01(\x04R\vlockVersion\x12\x10\n" +
	"\x03key\x18\x03 \x01(\fR\x03key\x12\x19\n" +
	"\block_ttl\x18\x04 \x01(\x04R\alockTtl\x12\x19\n" +
	"\btxn_size\x18\x05 \x01(\x04R\atxnSize\x12%\n" +
	"\tlock_type\x18\x06 \x01(\x0e2\b.tikv.OpR\blockType\x12+\n" +
	"\x12lock_for_update_ts\x18\a \x01(\x04R\x0flockForUpdateTs\x12(\n" +
	"\x10use_async_commit\x18\b \x01(\bR\x0euseAsyncCommit\x12\"\n" +
	"\rmin_commit_ts\x18\t \x01(\x04R\vminCommitTs\x12 \n" +
	"\vsecondaries\x18\n" +
	" \x03(\fR\vsecondaries\x12:\n" +
	"\x1aduration_to_last_update_ms\x18\v \x01(\x04R\x16durationToLastUpdateMs\x12\x1e\n" +
	"\vis_txn_file\x18d \x01(\bR\tisTxnFile\"\x89\x05\n" +
	"\bKeyError\x12&\n" +
	"\x06locked\x18\x01 \x01(\v2\x0e.tikv.LockInfoR\x06locked\x12\x1c\n" +
	"\tretryable\x18\x02 \x01(\tR\tretryable\x12\x14\n" +
	"\x05abort\x18\x03 \x01(\tR\x05abort\x12/\n" +
	"\bconflict\x18\x04 \x01(\v2\x13.tikv.WriteConflictR\bconflict\x127\n" +
	"\ralready_exist\x18\x05 \x01(\v2\x12.tikv.AlreadyExistR\falreadyExist\x12A\n" +
	"\x11commit_ts_expired\x18\a \x01(\v2\x15.tikv.CommitTsExpiredR\x0fcommitTsExpired\x125\n" +
	"\rtxn_not_found\x18\b \x01(\v2\x11.tikv.TxnNotFoundR\vtxnNotFound\x12E\n" +
	"\x13commit_ts_too_large\x18\t \x01(\v2\x16.tikv.CommitTsTooLargeR\x10commitTsTooLarge\x12@\n" +
	"\x10assertion_failed\x18\n" +
	" \x01(\v2\x15.tikv.AssertionFailedR\x0fassertionFailed\x12@\n" +
	"\x10primary_mismatch\x18\v \x01(\v2\x15.tikv.PrimaryMismatchR\x0fprimaryMismatch\x12B\n" +
	"\x12txn_lock_not_found\x18\f \x01(\v2\x15.tikv.TxnLockNotFoundR\x0ftxnLockNotFound\x12.\n" +
	"\n" +
	"debug_info\x18d \x01(\v2\x0f.tikv.DebugInfoR\tdebugInfo\"\xed\x02\n" +
	"\rWriteConflict\x12\x19\n" +
	"\bstart_ts\x18\x01 \x01(\x04R\astartTs\x12\x1f\n" +
	"\vconflict_ts\x18\x02 \x01(\x04R\n" +
	"conflictTs\x12\x10\n" +
	"\x03key\x18\x03 \x01(\fR\x03key\x12\x18\n" +
	"\aprimary\x18\x04 \x01(\fR\aprimary\x12,\n" +
	"\x12conflict_commit_ts\x18\x05 \x01(\x04R\x10conflictCommitTs\x122\n" +
	"\x06reason\x18\x06 \x01(\x0e2\x1a.tikv.WriteConflict.ReasonR\x06reason\"\x91\x01\n" +
	"\x06Reason\x12\v\n" +
	"\aUnknown\x10\x00\x12\x0e\n" +
	"\n" +
	"Optimistic\x10\x01\x12\x14\n" +
	"\x10PessimisticRetry\x10\x02\x12\x12\n" +
	"\x0eSelfRolledBack\x10\x03\x12\r\n" +
	"\tRcCheckTs\x10\x04\x12\x17\n" +
	"\x13LazyUniquenessCheck\x10\x05\x12\x18\n" +
	"\x14NotLockedKeyConflict\x10\x06\" \n" +
	"\fAlreadyExist\x12\x10\n" +
	"\x03key\x18\x01 \x01(\fR\x03key\"\x92\x01\n" +
	"\x0fCommitTsExpired\x12\x19\n" +
	"\bstart_ts\x18\x01 \x01(\x04R\astartTs\x12.\n" +
	"\x13attempted_commit_ts\x18\x02 \x01(\x04R\x11attemptedCommitTs\x12\x10\n" +
	"\x03key\x18\x03 \x01(\fR\x03key\x12\"\n" +
	"\rmin_commit_ts\x18\x04 \x01(\x04R\vminCommitTs\"I\n" +
	"\vTxnNotFound\x12\x19\n" +
	"\bstart_ts\x18\x01 \x01(\x04R\astartTs\x12\x1f\n" +
	"\vprimary_key\x18\x02 \x01(\fR\n" +
	"primaryKey\"/\n" +
	"\x10CommitTsTooLarge\x12\x1b\n" +
	"\tcommit_ts\x18\x01 \x01(\x04R\bcommitTs\"\xc7\x01\n" +
	"\x0fAssertionFailed\x12\x19\n" +
	"\bstart_ts\x18\x01 \x01(\x04R\astartTs\x12\x10\n" +
	"\x03key\x18\x02 \x01(\fR\x03key\x12-\n" +
	"\tassertion\x18\x03 \x01(\x0e2\x0f.tikv.AssertionR\tassertion\x12*\n" +
	"\x11existing_start_ts\x18\x04 \x01(\x04R\x0fexistingStartTs\x12,\n" +
	"\x12existing_commit_ts\x18\x05 \x01(\x04R\x10existingCommitTs\">\n" +
	"\x0fPrimaryMismatch\x12+\n" +
	"\tlock_info\x18\x01 \x01(\v2\x0e.tikv.LockInfoR\blockInfo\"#\n" +
	"\x0fTxnLockNotFound\x12\x10\n" +
	"\x03key\x18\x01 \x01(\fR\x03key\"E\n" +
	"\rMvccDebugInfo\x12\x10\n" +
	"\x03key\x18\x01 \x01(\fR\x03key\x12\"\n" +
	"\x04mvcc\x18\x02 \x01(\v2\x0e.tikv.MvccInfoR\x04mvcc\"=\n" +
	"\tDebugInfo\x120\n" +
	"\tmvcc_info\x18\x01 \x03(\v2\x13.tikv.MvccDebugInfoR\bmvccInfo\"\xcc\x01\n" +
	"\n" +
	"TimeDetail\x12)\n" +
	"\x11wait_wall_time_ms\x18\x01 \x01(\x04R\x0ewaitWallTimeMs\x12/\n" +
	"\x14process_wall_time_ms\x18\x02 \x01(\x04R\x11processWallTimeMs\x12.\n" +
	"\x14kv_read_wall_time_ms\x18\x03 \x01(\x04R\x10kvReadWallTimeMs\x122\n" +
	"\x16total_rpc_wall_time_ns\x18\x04 \x01(\x04R\x12totalRpcWallTimeNs\"\xf4\x02\n" +
	"\fTimeDetailV2\x12)\n" +
	"\x11wait_wall_time_ns\x18\x01 \x01(\x04R\x0ewaitWallTimeNs\x12/\n" +
	"\x14process_wall_time_ns\x18\x02 \x01(\x04R\x11processWallTimeNs\x12>\n" +
	"\x1cprocess_suspend_wall_time_ns\x18\x03 \x01(\x04R\x18processSuspendWallTimeNs\x12.\n" +
	"\x14kv_read_wall_time_ns\x18\x04 \x01(\x04R\x10kvReadWallTimeNs\x122\n" +
	"\x16total_rpc_wall_time_ns\x18\x05 \x01(\x04R\x12totalRpcWallTimeNs\x124\n" +
	"\x17kv_grpc_process_time_ns\x18\x06 \x01(\x04R\x13kvGrpcProcessTimeNs\x12.\n" +
	"\x14kv_grpc_wait_time_ns\x18\a \x01(\x04R\x10kvGrpcWaitTimeNs\"]\n" +
	"\bScanInfo\x12\x14\n" +
	"\x05total\x18\x01 \x01(\x03R\x05total\x12\x1c\n" +
	"\tprocessed\x18\x02 \x01(\x03R\tprocessed\x12\x1d\n" +
	"\n" +
	"read_bytes\x18\x03 \x01(\x03R\treadBytes\"z\n" +
	"\n" +
	"ScanDetail\x12$\n" +
	"\x05write\x18\x01 \x01(\v2\x0e.tikv.ScanInfoR\x05write\x12\"\n" +
	"\x04lock\x18\x02 \x01(\v2\x0e.tikv.ScanInfoR\x04lock\x12\"\n" +
	"\x04data\x18\x03 \x01(\v2\x0e.tikv.ScanInfoR\x04data\"\xf7\x05\n" +
	"\fScanDetailV2\x12-\n" +
	"\x12processed_versions\x18\x01 \x01(\x04R\x11processedVersions\x126\n" +
	"\x17processed_versions_size\x18\b \x01(\x04R\x15processedVersionsSize\x12%\n" +
	"\x0etotal_versions\x18\x02 \x01(\x04R\rtotalVersions\x12?\n" +
	"\x1crocksdb_delete_skipped_count\x18\x03 \x01(\x04R\x19rocksdbDeleteSkippedCount\x129\n" +
	"\x19rocksdb_key_skipped_count\x18\x04 \x01(\x04R\x16rocksdbKeySkippedCount\x12@\n" +
	"\x1drocksdb_block_cache_hit_count\x18\x05 \x01(\x04R\x19rocksdbBlockCacheHitCount\x127\n" +
	"\x18rocksdb_block_read_count\x18\x06 \x01(\x04R\x15rocksdbBlockReadCount\x125\n" +
	"\x17rocksdb_block_read_byte\x18\a \x01(\x04R\x14rocksdbBlockReadByte\x127\n" +
	"\x18rocksdb_block_read_nanos\x18\t \x01(\x04R\x15rocksdbBlockReadNanos\x12,\n" +
	"\x12get_snapshot_nanos\x18\n" +
	" \x01(\x04R\x10getSnapshotNanos\x12@\n" +
	"\x1dread_index_propose_wait_nanos\x18\v \x01(\x04R\x19readIndexProposeWaitNanos\x12@\n" +
	"\x1dread_index_confirm_wait_nanos\x18\f \x01(\x04R\x19readIndexConfirmWaitNanos\x12@\n" +
	"\x1dread_pool_schedule_wait_nanos\x18\r \x01(\x04R\x19readPoolScheduleWaitNanos\"\x7f\n" +
	"\vExecDetails\x121\n" +
	"\vtime_detail\x18\x01 \x01(\v2\x10.tikv.TimeDetailR\n" +
	"timeDetail\x121\n" +
	"\vscan_detail\x18\x02 \x01(\v2\x10.tikv.ScanDetailR\n" +
	"scanDetailJ\x04\b\x03\x10\x04J\x04\b\x04\x10\x05\"\xec\x01\n" +
	"\rExecDetailsV2\x121\n" +
	"\vtime_detail\x18\x01 \x01(\v2\x10.tikv.TimeDetailR\n" +
	"timeDetail\x128\n" +
	"\x0escan_detail_v2\x18\x02 \x01(\v2\x12.tikv.ScanDetailV2R\fscanDetailV2\x124\n" +
	"\fwrite_detail\x18\x03 \x01(\v2\x11.tikv.WriteDetailR\vwriteDetail\x128\n" +
	"\x0etime_detail_v2\x18\x04 \x01(\v2\x12.tikv.TimeDetailV2R\ftimeDetailV2\"\x81\a\n" +
	"\vWriteDetail\x123\n" +
	"\x16store_batch_wait_nanos\x18\x01 \x01(\x04R\x13storeBatchWaitNanos\x125\n" +
	"\x17propose_send_wait_nanos\x18\x02 \x01(\x04R\x14proposeSendWaitNanos\x12*\n" +
	"\x11persist_log_nanos\x18\x03 \x01(\x04R\x0fpersistLogNanos\x12C\n" +
	"\x1fraft_db_write_leader_wait_nanos\x18\x04 \x01(\x04R\x1araftDbWriteLeaderWaitNanos\x122\n" +
	"\x16raft_db_sync_log_nanos\x18\x05 \x01(\x04R\x12raftDbSyncLogNanos\x12>\n" +
	"\x1craft_db_write_memtable_nanos\x18\x06 \x01(\x04R\x18raftDbWriteMemtableNanos\x12(\n" +
	"\x10commit_log_nanos\x18\a \x01(\x04R\x0ecommitLogNanos\x123\n" +
	"\x16apply_batch_wait_nanos\x18\b \x01(\x04R\x13applyBatchWaitNanos\x12&\n" +
	"\x0fapply_log_nanos\x18\t \x01(\x04R\rapplyLogNanos\x123\n" +
	"\x16apply_mutex_lock_nanos\x18\n" +
	" \x01(\x04R\x13applyMutexLockNanos\x12@\n" +
	"\x1dapply_write_leader_wait_nanos\x18\v \x01(\x04R\x19applyWriteLeaderWaitNanos\x121\n" +
	"\x15apply_write_wal_nanos\x18\f \x01(\x04R\x12applyWriteWalNanos\x12;\n" +
	"\x1aapply_write_memtable_nanos\x18\r \x01(\x04R\x17applyWriteMemtableNanos\x12(\n" +
	"\x10latch_wait_nanos\x18\x0e \x01(\x04R\x0elatchWaitNanos\x12#\n" +
	"\rprocess_nanos\x18\x0f \x01(\x04R\fprocessNanos\x12%\n" +
	"\x0ethrottle_nanos\x18\x10 \x01(\x04R\rthrottleNanos\x12=\n" +
	"\x1bpessimistic_lock_wait_nanos\x18\x11 \x01(\x04R\x18pessimisticLockWaitNanos\"s\n" +
	"\x06KvPair\x12$\n" +
	"\x05error\x18\x01 \x01(\v2\x0e.tikv.KeyErrorR\x05error\x12\x10\n" +
	"\x03key\x18\x02 \x01(\fR\x03key\x12\x14\n" +
	"\x05value\x18\x03 \x01(\fR\x05value\x12\x1b\n" +
	"\tcommit_ts\x18\x04 \x01(\x04R\bcommitTs\"\xd4\x02\n" +
	"\tMvccWrite\x12\x1c\n" +
	"\x04type\x18\x01 \x01(\x0e2\b.tikv.OpR\x04type\x12\x19\n" +
	"\bstart_ts\x18\x02 \x01(\x04R\astartTs\x12\x1b\n" +
	"\tcommit_ts\x18\x03 \x01(\x04R\bcommitTs\x12\x1f\n" +
	"\vshort_value\x18\x04 \x01(\fR\n" +
	"shortValue\x126\n" +
	"\x17has_overlapped_rollback\x18\x05 \x01(\bR\x15hasOverlappedRollback\x12 \n" +
	"\fhas_gc_fence\x18\x06 \x01(\bR\n" +
	"hasGcFence\x12\x19\n" +
	"\bgc_fence\x18\a \x01(\x04R\agcFence\x12$\n" +
	"\x0elast_change_ts\x18\b \x01(\x04R\flastChangeTs\x125\n" +
	"\x17versions_to_last_change\x18\t \x01(\x04R\x14versionsToLastChange\"<\n" +
	"\tMvccValue\x12\x19\n" +
	"\bstart_ts\x18\x01 \x01(\x04R\astartTs\x12\x14\n" +
	"\x05value\x18\x02 \x01(\fR\x05value\"\x99\x03\n" +
	"\bMvccLock\x12\x1c\n" +
	"\x04type\x18\x01 \x01(\x0e2\b.tikv.OpR\x04type\x12\x19\n" +
	"\bstart_ts\x18\x02 \x01(\x04R\astartTs\x12\x18\n" +
	"\aprimary\x18\x03 \x01(\fR\aprimary\x12\x1f\n" +
	"\vshort_value\x18\x04 \x01(\fR\n" +
	"shortValue\x12\x10\n" +
	"\x03ttl\x18\x05 \x01(\x04R\x03ttl\x12\"\n" +
	"\rfor_update_ts\x18\x06 \x01(\x04R\vforUpdateTs\x12\x19\n" +
	"\btxn_size\x18\a \x01(\x04R\atxnSize\x12(\n" +
	"\x10use_async_commit\x18\b \x01(\bR\x0euseAsyncCommit\x12 \n" +
	"\vsecondaries\x18\t \x03(\fR\vsecondaries\x12\x1f\n" +
	"\vrollback_ts\x18\n" +
	" \x03(\x04R\n" +
	"rollbackTs\x12$\n" +
	"\x0elast_change_ts\x18\v \x01(\x04R\flastChangeTs\x125\n" +
	"\x17versions_to_last_change\x18\f \x01(\x04R\x14versionsToLastChange\"\x80\x01\n" +
	"\bMvccInfo\x12\"\n" +
	"\x04lock\x18\x01 \x01(\v2\x0e.tikv.MvccLockR\x04lock\x12'\n" +
	"\x06writes\x18\x02 \x03(\v2\x0f.tikv.MvccWriteR\x06writes\x12'\n" +
	"\x06values\x18\x03 \x03(\v2\x0f.tikv.MvccValueR\x06values*'\n" +
	"\n" +
	"APIVersion\x12\x06\n" +
	"\x02V1\x10\x00\x12\t\n" +
	"\x05V1TTL\x10\x01\x12\x06\n" +
	"\x02V2\x10\x02*+\n" +
	"\n" +
	"CommandPri\x12\n" +
	"\n" +
	"\x06Normal\x10\x00\x12\a\n" +
	"\x03Low\x10\x01\x12\b\n" +
	"\x04High\x10\x02*/\n" +
	"\x0eIsolationLevel\x12\x06\n" +
	"\x02SI\x10\x00\x12\x06\n" +
	"\x02RC\x10\x01\x12\r\n" +
	"\tRCCheckTS\x10\x02*V\n" +
	"\vDiskFullOpt\x12\x14\n" +
	"\x10NotAllowedOnFull\x10\x00\x12\x17\n" +
	"\x13AllowedOnAlmostFull\x10\x01\x12\x18\n" +
	"\x14AllowedOnAlreadyFull\x10\x02*c\n" +
	"\x02Op\x12\a\n" +
	"\x03Put\x10\x00\x12\a\n" +
	"\x03Del\x10\x01\x12\b\n" +
	"\x04Lock\x10\x02\x12\f\n" +
	"\bRollback\x10\x03\x12\n" +
	"\n" +
	"\x06Insert\x10\x04\x12\x13\n" +
	"\x0fPessimisticLock\x10\x05\x12\x12\n" +
	"\x0eCheckNotExists\x10\x06*.\n" +
	"\tAssertion\x12\b\n" +
	"\x04None\x10\x00\x12\t\n" +
	"\x05Exist\x10\x01\x12\f\n" +
	"\bNotExist\x10\x02B)Z'github.com/vdaas/vald/apis/grpc/v1/tikvb\x06proto3"

var (
	file_v1_tikv_kvrpcpb_proto_rawDescOnce sync.Once
	file_v1_tikv_kvrpcpb_proto_rawDescData []byte
)

func file_v1_tikv_kvrpcpb_proto_rawDescGZIP() []byte {
	file_v1_tikv_kvrpcpb_proto_rawDescOnce.Do(func() {
		file_v1_tikv_kvrpcpb_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_v1_tikv_kvrpcpb_proto_rawDesc), len(file_v1_tikv_kvrpcpb_proto_rawDesc)))
	})
	return file_v1_tikv_kvrpcpb_proto_rawDescData
}

var file_v1_tikv_kvrpcpb_proto_enumTypes = make([]protoimpl.EnumInfo, 7)
var file_v1_tikv_kvrpcpb_proto_msgTypes = make([]protoimpl.MessageInfo, 40)
var file_v1_tikv_kvrpcpb_proto_goTypes = []any{
	(APIVersion)(0),                // 0: tikv.APIVersion
	(CommandPri)(0),                // 1: tikv.CommandPri
	(IsolationLevel)(0),            // 2: tikv.IsolationLevel
	(DiskFullOpt)(0),               // 3: tikv.DiskFullOpt
	(Op)(0),                        // 4: tikv.Op
	(Assertion)(0),                 // 5: tikv.Assertion
	(WriteConflict_Reason)(0),      // 6: tikv.WriteConflict.Reason
	(*RawGetRequest)(nil),          // 7: tikv.RawGetRequest
	(*RawGetResponse)(nil),         // 8: tikv.RawGetResponse
	(*RawBatchGetRequest)(nil),     // 9: tikv.RawBatchGetRequest
	(*RawBatchGetResponse)(nil),    // 10: tikv.RawBatchGetResponse
	(*RawPutRequest)(nil),          // 11: tikv.RawPutRequest
	(*RawPutResponse)(nil),         // 12: tikv.RawPutResponse
	(*RawBatchPutRequest)(nil),     // 13: tikv.RawBatchPutRequest
	(*RawBatchPutResponse)(nil),    // 14: tikv.RawBatchPutResponse
	(*RawDeleteRequest)(nil),       // 15: tikv.RawDeleteRequest
	(*RawDeleteResponse)(nil),      // 16: tikv.RawDeleteResponse
	(*RawBatchDeleteRequest)(nil),  // 17: tikv.RawBatchDeleteRequest
	(*RawBatchDeleteResponse)(nil), // 18: tikv.RawBatchDeleteResponse
	(*Context)(nil),                // 19: tikv.Context
	(*ResourceControlContext)(nil), // 20: tikv.ResourceControlContext
	(*SourceStmt)(nil),             // 21: tikv.SourceStmt
	(*LockInfo)(nil),               // 22: tikv.LockInfo
	(*KeyError)(nil),               // 23: tikv.KeyError
	(*WriteConflict)(nil),          // 24: tikv.WriteConflict
	(*AlreadyExist)(nil),           // 25: tikv.AlreadyExist
	(*CommitTsExpired)(nil),        // 26: tikv.CommitTsExpired
	(*TxnNotFound)(nil),            // 27: tikv.TxnNotFound
	(*CommitTsTooLarge)(nil),       // 28: tikv.CommitTsTooLarge
	(*AssertionFailed)(nil),        // 29: tikv.AssertionFailed
	(*PrimaryMismatch)(nil),        // 30: tikv.PrimaryMismatch
	(*TxnLockNotFound)(nil),        // 31: tikv.TxnLockNotFound
	(*MvccDebugInfo)(nil),          // 32: tikv.MvccDebugInfo
	(*DebugInfo)(nil),              // 33: tikv.DebugInfo
	(*TimeDetail)(nil),             // 34: tikv.TimeDetail
	(*TimeDetailV2)(nil),           // 35: tikv.TimeDetailV2
	(*ScanInfo)(nil),               // 36: tikv.ScanInfo
	(*ScanDetail)(nil),             // 37: tikv.ScanDetail
	(*ScanDetailV2)(nil),           // 38: tikv.ScanDetailV2
	(*ExecDetails)(nil),            // 39: tikv.ExecDetails
	(*ExecDetailsV2)(nil),          // 40: tikv.ExecDetailsV2
	(*WriteDetail)(nil),            // 41: tikv.WriteDetail
	(*KvPair)(nil),                 // 42: tikv.KvPair
	(*MvccWrite)(nil),              // 43: tikv.MvccWrite
	(*MvccValue)(nil),              // 44: tikv.MvccValue
	(*MvccLock)(nil),               // 45: tikv.MvccLock
	(*MvccInfo)(nil),               // 46: tikv.MvccInfo
	(*Error)(nil),                  // 47: tikv.Error
	(*RegionEpoch)(nil),            // 48: metapb.RegionEpoch
	(*Peer)(nil),                   // 49: metapb.Peer
}
var file_v1_tikv_kvrpcpb_proto_depIdxs = []int32{
	19, // 0: tikv.RawGetRequest.context:type_name -> tikv.Context
	47, // 1: tikv.RawGetResponse.region_error:type_name -> tikv.Error
	19, // 2: tikv.RawBatchGetRequest.context:type_name -> tikv.Context
	47, // 3: tikv.RawBatchGetResponse.region_error:type_name -> tikv.Error
	42, // 4: tikv.RawBatchGetResponse.pairs:type_name -> tikv.KvPair
	19, // 5: tikv.RawPutRequest.context:type_name -> tikv.Context
	47, // 6: tikv.RawPutResponse.region_error:type_name -> tikv.Error
	19, // 7: tikv.RawBatchPutRequest.context:type_name -> tikv.Context
	42, // 8: tikv.RawBatchPutRequest.pairs:type_name -> tikv.KvPair
	47, // 9: tikv.RawBatchPutResponse.region_error:type_name -> tikv.Error
	19, // 10: tikv.RawDeleteRequest.context:type_name -> tikv.Context
	47, // 11: tikv.RawDeleteResponse.region_error:type_name -> tikv.Error
	19, // 12: tikv.RawBatchDeleteRequest.context:type_name -> tikv.Context
	47, // 13: tikv.RawBatchDeleteResponse.region_error:type_name -> tikv.Error
	48, // 14: tikv.Context.region_epoch:type_name -> metapb.RegionEpoch
	49, // 15: tikv.Context.peer:type_name -> metapb.Peer
	1,  // 16: tikv.Context.priority:type_name -> tikv.CommandPri
	2,  // 17: tikv.Context.isolation_level:type_name -> tikv.IsolationLevel
	3,  // 18: tikv.Context.disk_full_opt:type_name -> tikv.DiskFullOpt
	0,  // 19: tikv.Context.api_version:type_name -> tikv.APIVersion
	20, // 20: tikv.Context.resource_control_context:type_name -> tikv.ResourceControlContext
	21, // 21: tikv.Context.source_stmt:type_name -> tikv.SourceStmt
	4,  // 22: tikv.LockInfo.lock_type:type_name -> tikv.Op
	22, // 23: tikv.KeyError.locked:type_name -> tikv.LockInfo
	24, // 24: tikv.KeyError.conflict:type_name -> tikv.WriteConflict
	25, // 25: tikv.KeyError.already_exist:type_name -> tikv.AlreadyExist
	26, // 26: tikv.KeyError.commit_ts_expired:type_name -> tikv.CommitTsExpired
	27, // 27: tikv.KeyError.txn_not_found:type_name -> tikv.TxnNotFound
	28, // 28: tikv.KeyError.commit_ts_too_large:type_name -> tikv.CommitTsTooLarge
	29, // 29: tikv.KeyError.assertion_failed:type_name -> tikv.AssertionFailed
	30, // 30: tikv.KeyError.primary_mismatch:type_name -> tikv.PrimaryMismatch
	31, // 31: tikv.KeyError.txn_lock_not_found:type_name -> tikv.TxnLockNotFound
	33, // 32: tikv.KeyError.debug_info:type_name -> tikv.DebugInfo
	6,  // 33: tikv.WriteConflict.reason:type_name -> tikv.WriteConflict.Reason
	5,  // 34: tikv.AssertionFailed.assertion:type_name -> tikv.Assertion
	22, // 35: tikv.PrimaryMismatch.lock_info:type_name -> tikv.LockInfo
	46, // 36: tikv.MvccDebugInfo.mvcc:type_name -> tikv.MvccInfo
	32, // 37: tikv.DebugInfo.mvcc_info:type_name -> tikv.MvccDebugInfo
	36, // 38: tikv.ScanDetail.write:type_name -> tikv.ScanInfo
	36, // 39: tikv.ScanDetail.lock:type_name -> tikv.ScanInfo
	36, // 40: tikv.ScanDetail.data:type_name -> tikv.ScanInfo
	34, // 41: tikv.ExecDetails.time_detail:type_name -> tikv.TimeDetail
	37, // 42: tikv.ExecDetails.scan_detail:type_name -> tikv.ScanDetail
	34, // 43: tikv.ExecDetailsV2.time_detail:type_name -> tikv.TimeDetail
	38, // 44: tikv.ExecDetailsV2.scan_detail_v2:type_name -> tikv.ScanDetailV2
	41, // 45: tikv.ExecDetailsV2.write_detail:type_name -> tikv.WriteDetail
	35, // 46: tikv.ExecDetailsV2.time_detail_v2:type_name -> tikv.TimeDetailV2
	23, // 47: tikv.KvPair.error:type_name -> tikv.KeyError
	4,  // 48: tikv.MvccWrite.type:type_name -> tikv.Op
	4,  // 49: tikv.MvccLock.type:type_name -> tikv.Op
	45, // 50: tikv.MvccInfo.lock:type_name -> tikv.MvccLock
	43, // 51: tikv.MvccInfo.writes:type_name -> tikv.MvccWrite
	44, // 52: tikv.MvccInfo.values:type_name -> tikv.MvccValue
	53, // [53:53] is the sub-list for method output_type
	53, // [53:53] is the sub-list for method input_type
	53, // [53:53] is the sub-list for extension type_name
	53, // [53:53] is the sub-list for extension extendee
	0,  // [0:53] is the sub-list for field type_name
}

func init() { file_v1_tikv_kvrpcpb_proto_init() }
func file_v1_tikv_kvrpcpb_proto_init() {
	if File_v1_tikv_kvrpcpb_proto != nil {
		return
	}
	file_v1_tikv_errorpb_proto_init()
	file_v1_tikv_metapb_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_v1_tikv_kvrpcpb_proto_rawDesc), len(file_v1_tikv_kvrpcpb_proto_rawDesc)),
			NumEnums:      7,
			NumMessages:   40,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_v1_tikv_kvrpcpb_proto_goTypes,
		DependencyIndexes: file_v1_tikv_kvrpcpb_proto_depIdxs,
		EnumInfos:         file_v1_tikv_kvrpcpb_proto_enumTypes,
		MessageInfos:      file_v1_tikv_kvrpcpb_proto_msgTypes,
	}.Build()
	File_v1_tikv_kvrpcpb_proto = out.File
	file_v1_tikv_kvrpcpb_proto_goTypes = nil
	file_v1_tikv_kvrpcpb_proto_depIdxs = nil
}
