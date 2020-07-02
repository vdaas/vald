//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package payload

import (
	encoding_binary "encoding/binary"
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Search struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Search) Reset()         { *m = Search{} }
func (m *Search) String() string { return proto.CompactTextString(m) }
func (*Search) ProtoMessage()    {}
func (*Search) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{0}
}
func (m *Search) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Search) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Search.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Search) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Search.Merge(m, src)
}
func (m *Search) XXX_Size() int {
	return m.Size()
}
func (m *Search) XXX_DiscardUnknown() {
	xxx_messageInfo_Search.DiscardUnknown(m)
}

var xxx_messageInfo_Search proto.InternalMessageInfo

type Search_Request struct {
	Vector               []float32      `protobuf:"fixed32,1,rep,packed,name=vector,proto3" json:"vector,omitempty"`
	Config               *Search_Config `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Search_Request) Reset()         { *m = Search_Request{} }
func (m *Search_Request) String() string { return proto.CompactTextString(m) }
func (*Search_Request) ProtoMessage()    {}
func (*Search_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{0, 0}
}
func (m *Search_Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Search_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Search_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Search_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Search_Request.Merge(m, src)
}
func (m *Search_Request) XXX_Size() int {
	return m.Size()
}
func (m *Search_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Search_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Search_Request proto.InternalMessageInfo

func (m *Search_Request) GetVector() []float32 {
	if m != nil {
		return m.Vector
	}
	return nil
}

func (m *Search_Request) GetConfig() *Search_Config {
	if m != nil {
		return m.Config
	}
	return nil
}

type Search_IDRequest struct {
	Id                   string         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Config               *Search_Config `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Search_IDRequest) Reset()         { *m = Search_IDRequest{} }
func (m *Search_IDRequest) String() string { return proto.CompactTextString(m) }
func (*Search_IDRequest) ProtoMessage()    {}
func (*Search_IDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{0, 1}
}
func (m *Search_IDRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Search_IDRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Search_IDRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Search_IDRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Search_IDRequest.Merge(m, src)
}
func (m *Search_IDRequest) XXX_Size() int {
	return m.Size()
}
func (m *Search_IDRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_Search_IDRequest.DiscardUnknown(m)
}

var xxx_messageInfo_Search_IDRequest proto.InternalMessageInfo

func (m *Search_IDRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Search_IDRequest) GetConfig() *Search_Config {
	if m != nil {
		return m.Config
	}
	return nil
}

type Search_ObjectRequest struct {
	Object               []byte         `protobuf:"bytes,1,opt,name=object,proto3" json:"object,omitempty"`
	Config               *Search_Config `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Search_ObjectRequest) Reset()         { *m = Search_ObjectRequest{} }
func (m *Search_ObjectRequest) String() string { return proto.CompactTextString(m) }
func (*Search_ObjectRequest) ProtoMessage()    {}
func (*Search_ObjectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{0, 2}
}
func (m *Search_ObjectRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Search_ObjectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Search_ObjectRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Search_ObjectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Search_ObjectRequest.Merge(m, src)
}
func (m *Search_ObjectRequest) XXX_Size() int {
	return m.Size()
}
func (m *Search_ObjectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_Search_ObjectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_Search_ObjectRequest proto.InternalMessageInfo

func (m *Search_ObjectRequest) GetObject() []byte {
	if m != nil {
		return m.Object
	}
	return nil
}

func (m *Search_ObjectRequest) GetConfig() *Search_Config {
	if m != nil {
		return m.Config
	}
	return nil
}

type Search_Config struct {
	Num                  uint32   `protobuf:"varint,1,opt,name=num,proto3" json:"num,omitempty"`
	Radius               float32  `protobuf:"fixed32,2,opt,name=radius,proto3" json:"radius,omitempty"`
	Epsilon              float32  `protobuf:"fixed32,3,opt,name=epsilon,proto3" json:"epsilon,omitempty"`
	Timeout              int64    `protobuf:"varint,4,opt,name=timeout,proto3" json:"timeout,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Search_Config) Reset()         { *m = Search_Config{} }
func (m *Search_Config) String() string { return proto.CompactTextString(m) }
func (*Search_Config) ProtoMessage()    {}
func (*Search_Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{0, 3}
}
func (m *Search_Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Search_Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Search_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Search_Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Search_Config.Merge(m, src)
}
func (m *Search_Config) XXX_Size() int {
	return m.Size()
}
func (m *Search_Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Search_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Search_Config proto.InternalMessageInfo

func (m *Search_Config) GetNum() uint32 {
	if m != nil {
		return m.Num
	}
	return 0
}

func (m *Search_Config) GetRadius() float32 {
	if m != nil {
		return m.Radius
	}
	return 0
}

func (m *Search_Config) GetEpsilon() float32 {
	if m != nil {
		return m.Epsilon
	}
	return 0
}

func (m *Search_Config) GetTimeout() int64 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

type Search_Response struct {
	Results              []*Object_Distance `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Search_Response) Reset()         { *m = Search_Response{} }
func (m *Search_Response) String() string { return proto.CompactTextString(m) }
func (*Search_Response) ProtoMessage()    {}
func (*Search_Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{0, 4}
}
func (m *Search_Response) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Search_Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Search_Response.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Search_Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Search_Response.Merge(m, src)
}
func (m *Search_Response) XXX_Size() int {
	return m.Size()
}
func (m *Search_Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Search_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Search_Response proto.InternalMessageInfo

func (m *Search_Response) GetResults() []*Object_Distance {
	if m != nil {
		return m.Results
	}
	return nil
}

type Meta struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Meta) Reset()         { *m = Meta{} }
func (m *Meta) String() string { return proto.CompactTextString(m) }
func (*Meta) ProtoMessage()    {}
func (*Meta) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{1}
}
func (m *Meta) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Meta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Meta.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Meta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta.Merge(m, src)
}
func (m *Meta) XXX_Size() int {
	return m.Size()
}
func (m *Meta) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta.DiscardUnknown(m)
}

var xxx_messageInfo_Meta proto.InternalMessageInfo

type Meta_Key struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Meta_Key) Reset()         { *m = Meta_Key{} }
func (m *Meta_Key) String() string { return proto.CompactTextString(m) }
func (*Meta_Key) ProtoMessage()    {}
func (*Meta_Key) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{1, 0}
}
func (m *Meta_Key) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Meta_Key) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Meta_Key.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Meta_Key) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta_Key.Merge(m, src)
}
func (m *Meta_Key) XXX_Size() int {
	return m.Size()
}
func (m *Meta_Key) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta_Key.DiscardUnknown(m)
}

var xxx_messageInfo_Meta_Key proto.InternalMessageInfo

func (m *Meta_Key) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type Meta_Keys struct {
	Keys                 []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Meta_Keys) Reset()         { *m = Meta_Keys{} }
func (m *Meta_Keys) String() string { return proto.CompactTextString(m) }
func (*Meta_Keys) ProtoMessage()    {}
func (*Meta_Keys) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{1, 1}
}
func (m *Meta_Keys) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Meta_Keys) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Meta_Keys.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Meta_Keys) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta_Keys.Merge(m, src)
}
func (m *Meta_Keys) XXX_Size() int {
	return m.Size()
}
func (m *Meta_Keys) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta_Keys.DiscardUnknown(m)
}

var xxx_messageInfo_Meta_Keys proto.InternalMessageInfo

func (m *Meta_Keys) GetKeys() []string {
	if m != nil {
		return m.Keys
	}
	return nil
}

type Meta_Val struct {
	Val                  string   `protobuf:"bytes,1,opt,name=val,proto3" json:"val,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Meta_Val) Reset()         { *m = Meta_Val{} }
func (m *Meta_Val) String() string { return proto.CompactTextString(m) }
func (*Meta_Val) ProtoMessage()    {}
func (*Meta_Val) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{1, 2}
}
func (m *Meta_Val) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Meta_Val) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Meta_Val.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Meta_Val) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta_Val.Merge(m, src)
}
func (m *Meta_Val) XXX_Size() int {
	return m.Size()
}
func (m *Meta_Val) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta_Val.DiscardUnknown(m)
}

var xxx_messageInfo_Meta_Val proto.InternalMessageInfo

func (m *Meta_Val) GetVal() string {
	if m != nil {
		return m.Val
	}
	return ""
}

type Meta_Vals struct {
	Vals                 []string `protobuf:"bytes,1,rep,name=vals,proto3" json:"vals,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Meta_Vals) Reset()         { *m = Meta_Vals{} }
func (m *Meta_Vals) String() string { return proto.CompactTextString(m) }
func (*Meta_Vals) ProtoMessage()    {}
func (*Meta_Vals) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{1, 3}
}
func (m *Meta_Vals) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Meta_Vals) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Meta_Vals.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Meta_Vals) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta_Vals.Merge(m, src)
}
func (m *Meta_Vals) XXX_Size() int {
	return m.Size()
}
func (m *Meta_Vals) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta_Vals.DiscardUnknown(m)
}

var xxx_messageInfo_Meta_Vals proto.InternalMessageInfo

func (m *Meta_Vals) GetVals() []string {
	if m != nil {
		return m.Vals
	}
	return nil
}

type Meta_KeyVal struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Val                  string   `protobuf:"bytes,2,opt,name=val,proto3" json:"val,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Meta_KeyVal) Reset()         { *m = Meta_KeyVal{} }
func (m *Meta_KeyVal) String() string { return proto.CompactTextString(m) }
func (*Meta_KeyVal) ProtoMessage()    {}
func (*Meta_KeyVal) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{1, 4}
}
func (m *Meta_KeyVal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Meta_KeyVal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Meta_KeyVal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Meta_KeyVal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta_KeyVal.Merge(m, src)
}
func (m *Meta_KeyVal) XXX_Size() int {
	return m.Size()
}
func (m *Meta_KeyVal) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta_KeyVal.DiscardUnknown(m)
}

var xxx_messageInfo_Meta_KeyVal proto.InternalMessageInfo

func (m *Meta_KeyVal) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Meta_KeyVal) GetVal() string {
	if m != nil {
		return m.Val
	}
	return ""
}

type Meta_KeyVals struct {
	Kvs                  []*Meta_KeyVal `protobuf:"bytes,1,rep,name=kvs,proto3" json:"kvs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Meta_KeyVals) Reset()         { *m = Meta_KeyVals{} }
func (m *Meta_KeyVals) String() string { return proto.CompactTextString(m) }
func (*Meta_KeyVals) ProtoMessage()    {}
func (*Meta_KeyVals) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{1, 5}
}
func (m *Meta_KeyVals) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Meta_KeyVals) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Meta_KeyVals.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Meta_KeyVals) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta_KeyVals.Merge(m, src)
}
func (m *Meta_KeyVals) XXX_Size() int {
	return m.Size()
}
func (m *Meta_KeyVals) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta_KeyVals.DiscardUnknown(m)
}

var xxx_messageInfo_Meta_KeyVals proto.InternalMessageInfo

func (m *Meta_KeyVals) GetKvs() []*Meta_KeyVal {
	if m != nil {
		return m.Kvs
	}
	return nil
}

type Object struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Object) Reset()         { *m = Object{} }
func (m *Object) String() string { return proto.CompactTextString(m) }
func (*Object) ProtoMessage()    {}
func (*Object) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{2}
}
func (m *Object) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object.Merge(m, src)
}
func (m *Object) XXX_Size() int {
	return m.Size()
}
func (m *Object) XXX_DiscardUnknown() {
	xxx_messageInfo_Object.DiscardUnknown(m)
}

var xxx_messageInfo_Object proto.InternalMessageInfo

type Object_Distance struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Distance             float32  `protobuf:"fixed32,2,opt,name=distance,proto3" json:"distance,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Object_Distance) Reset()         { *m = Object_Distance{} }
func (m *Object_Distance) String() string { return proto.CompactTextString(m) }
func (*Object_Distance) ProtoMessage()    {}
func (*Object_Distance) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{2, 0}
}
func (m *Object_Distance) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_Distance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_Distance.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_Distance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_Distance.Merge(m, src)
}
func (m *Object_Distance) XXX_Size() int {
	return m.Size()
}
func (m *Object_Distance) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_Distance.DiscardUnknown(m)
}

var xxx_messageInfo_Object_Distance proto.InternalMessageInfo

func (m *Object_Distance) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Object_Distance) GetDistance() float32 {
	if m != nil {
		return m.Distance
	}
	return 0
}

type Object_ID struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Object_ID) Reset()         { *m = Object_ID{} }
func (m *Object_ID) String() string { return proto.CompactTextString(m) }
func (*Object_ID) ProtoMessage()    {}
func (*Object_ID) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{2, 1}
}
func (m *Object_ID) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_ID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_ID.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_ID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_ID.Merge(m, src)
}
func (m *Object_ID) XXX_Size() int {
	return m.Size()
}
func (m *Object_ID) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_ID.DiscardUnknown(m)
}

var xxx_messageInfo_Object_ID proto.InternalMessageInfo

func (m *Object_ID) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type Object_IDs struct {
	Ids                  []string `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Object_IDs) Reset()         { *m = Object_IDs{} }
func (m *Object_IDs) String() string { return proto.CompactTextString(m) }
func (*Object_IDs) ProtoMessage()    {}
func (*Object_IDs) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{2, 2}
}
func (m *Object_IDs) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_IDs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_IDs.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_IDs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_IDs.Merge(m, src)
}
func (m *Object_IDs) XXX_Size() int {
	return m.Size()
}
func (m *Object_IDs) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_IDs.DiscardUnknown(m)
}

var xxx_messageInfo_Object_IDs proto.InternalMessageInfo

func (m *Object_IDs) GetIds() []string {
	if m != nil {
		return m.Ids
	}
	return nil
}

type Object_Vector struct {
	Id                   string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Vector               []float32 `protobuf:"fixed32,2,rep,packed,name=vector,proto3" json:"vector,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Object_Vector) Reset()         { *m = Object_Vector{} }
func (m *Object_Vector) String() string { return proto.CompactTextString(m) }
func (*Object_Vector) ProtoMessage()    {}
func (*Object_Vector) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{2, 3}
}
func (m *Object_Vector) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_Vector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_Vector.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_Vector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_Vector.Merge(m, src)
}
func (m *Object_Vector) XXX_Size() int {
	return m.Size()
}
func (m *Object_Vector) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_Vector.DiscardUnknown(m)
}

var xxx_messageInfo_Object_Vector proto.InternalMessageInfo

func (m *Object_Vector) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Object_Vector) GetVector() []float32 {
	if m != nil {
		return m.Vector
	}
	return nil
}

type Object_Vectors struct {
	Vectors              []*Object_Vector `protobuf:"bytes,1,rep,name=vectors,proto3" json:"vectors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Object_Vectors) Reset()         { *m = Object_Vectors{} }
func (m *Object_Vectors) String() string { return proto.CompactTextString(m) }
func (*Object_Vectors) ProtoMessage()    {}
func (*Object_Vectors) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{2, 4}
}
func (m *Object_Vectors) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_Vectors) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_Vectors.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_Vectors) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_Vectors.Merge(m, src)
}
func (m *Object_Vectors) XXX_Size() int {
	return m.Size()
}
func (m *Object_Vectors) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_Vectors.DiscardUnknown(m)
}

var xxx_messageInfo_Object_Vectors proto.InternalMessageInfo

func (m *Object_Vectors) GetVectors() []*Object_Vector {
	if m != nil {
		return m.Vectors
	}
	return nil
}

type Object_Blob struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Object               []byte   `protobuf:"bytes,2,opt,name=object,proto3" json:"object,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Object_Blob) Reset()         { *m = Object_Blob{} }
func (m *Object_Blob) String() string { return proto.CompactTextString(m) }
func (*Object_Blob) ProtoMessage()    {}
func (*Object_Blob) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{2, 5}
}
func (m *Object_Blob) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_Blob) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_Blob.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_Blob) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_Blob.Merge(m, src)
}
func (m *Object_Blob) XXX_Size() int {
	return m.Size()
}
func (m *Object_Blob) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_Blob.DiscardUnknown(m)
}

var xxx_messageInfo_Object_Blob proto.InternalMessageInfo

func (m *Object_Blob) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Object_Blob) GetObject() []byte {
	if m != nil {
		return m.Object
	}
	return nil
}

type Object_Location struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Uuid                 string   `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Ips                  []string `protobuf:"bytes,3,rep,name=ips,proto3" json:"ips,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Object_Location) Reset()         { *m = Object_Location{} }
func (m *Object_Location) String() string { return proto.CompactTextString(m) }
func (*Object_Location) ProtoMessage()    {}
func (*Object_Location) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{2, 6}
}
func (m *Object_Location) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_Location) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_Location.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_Location) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_Location.Merge(m, src)
}
func (m *Object_Location) XXX_Size() int {
	return m.Size()
}
func (m *Object_Location) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_Location.DiscardUnknown(m)
}

var xxx_messageInfo_Object_Location proto.InternalMessageInfo

func (m *Object_Location) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Object_Location) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Object_Location) GetIps() []string {
	if m != nil {
		return m.Ips
	}
	return nil
}

type Object_Locations struct {
	Locations            []*Object_Location `protobuf:"bytes,1,rep,name=locations,proto3" json:"locations,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Object_Locations) Reset()         { *m = Object_Locations{} }
func (m *Object_Locations) String() string { return proto.CompactTextString(m) }
func (*Object_Locations) ProtoMessage()    {}
func (*Object_Locations) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{2, 7}
}
func (m *Object_Locations) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_Locations) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_Locations.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_Locations) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_Locations.Merge(m, src)
}
func (m *Object_Locations) XXX_Size() int {
	return m.Size()
}
func (m *Object_Locations) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_Locations.DiscardUnknown(m)
}

var xxx_messageInfo_Object_Locations proto.InternalMessageInfo

func (m *Object_Locations) GetLocations() []*Object_Location {
	if m != nil {
		return m.Locations
	}
	return nil
}

type Control struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Control) Reset()         { *m = Control{} }
func (m *Control) String() string { return proto.CompactTextString(m) }
func (*Control) ProtoMessage()    {}
func (*Control) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{3}
}
func (m *Control) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Control) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Control.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Control) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Control.Merge(m, src)
}
func (m *Control) XXX_Size() int {
	return m.Size()
}
func (m *Control) XXX_DiscardUnknown() {
	xxx_messageInfo_Control.DiscardUnknown(m)
}

var xxx_messageInfo_Control proto.InternalMessageInfo

type Control_CreateIndexRequest struct {
	PoolSize             uint32   `protobuf:"varint,1,opt,name=pool_size,json=poolSize,proto3" json:"pool_size,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Control_CreateIndexRequest) Reset()         { *m = Control_CreateIndexRequest{} }
func (m *Control_CreateIndexRequest) String() string { return proto.CompactTextString(m) }
func (*Control_CreateIndexRequest) ProtoMessage()    {}
func (*Control_CreateIndexRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{3, 0}
}
func (m *Control_CreateIndexRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Control_CreateIndexRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Control_CreateIndexRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Control_CreateIndexRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Control_CreateIndexRequest.Merge(m, src)
}
func (m *Control_CreateIndexRequest) XXX_Size() int {
	return m.Size()
}
func (m *Control_CreateIndexRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_Control_CreateIndexRequest.DiscardUnknown(m)
}

var xxx_messageInfo_Control_CreateIndexRequest proto.InternalMessageInfo

func (m *Control_CreateIndexRequest) GetPoolSize() uint32 {
	if m != nil {
		return m.PoolSize
	}
	return 0
}

type Replication struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Replication) Reset()         { *m = Replication{} }
func (m *Replication) String() string { return proto.CompactTextString(m) }
func (*Replication) ProtoMessage()    {}
func (*Replication) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{4}
}
func (m *Replication) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Replication) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Replication.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Replication) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Replication.Merge(m, src)
}
func (m *Replication) XXX_Size() int {
	return m.Size()
}
func (m *Replication) XXX_DiscardUnknown() {
	xxx_messageInfo_Replication.DiscardUnknown(m)
}

var xxx_messageInfo_Replication proto.InternalMessageInfo

type Replication_Recovery struct {
	DeletedAgents        []string `protobuf:"bytes,1,rep,name=deleted_agents,json=deletedAgents,proto3" json:"deleted_agents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Replication_Recovery) Reset()         { *m = Replication_Recovery{} }
func (m *Replication_Recovery) String() string { return proto.CompactTextString(m) }
func (*Replication_Recovery) ProtoMessage()    {}
func (*Replication_Recovery) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{4, 0}
}
func (m *Replication_Recovery) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Replication_Recovery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Replication_Recovery.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Replication_Recovery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Replication_Recovery.Merge(m, src)
}
func (m *Replication_Recovery) XXX_Size() int {
	return m.Size()
}
func (m *Replication_Recovery) XXX_DiscardUnknown() {
	xxx_messageInfo_Replication_Recovery.DiscardUnknown(m)
}

var xxx_messageInfo_Replication_Recovery proto.InternalMessageInfo

func (m *Replication_Recovery) GetDeletedAgents() []string {
	if m != nil {
		return m.DeletedAgents
	}
	return nil
}

type Replication_Rebalance struct {
	HighUsageAgents      []string `protobuf:"bytes,1,rep,name=high_usage_agents,json=highUsageAgents,proto3" json:"high_usage_agents,omitempty"`
	LowUsageAgents       []string `protobuf:"bytes,2,rep,name=low_usage_agents,json=lowUsageAgents,proto3" json:"low_usage_agents,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Replication_Rebalance) Reset()         { *m = Replication_Rebalance{} }
func (m *Replication_Rebalance) String() string { return proto.CompactTextString(m) }
func (*Replication_Rebalance) ProtoMessage()    {}
func (*Replication_Rebalance) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{4, 1}
}
func (m *Replication_Rebalance) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Replication_Rebalance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Replication_Rebalance.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Replication_Rebalance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Replication_Rebalance.Merge(m, src)
}
func (m *Replication_Rebalance) XXX_Size() int {
	return m.Size()
}
func (m *Replication_Rebalance) XXX_DiscardUnknown() {
	xxx_messageInfo_Replication_Rebalance.DiscardUnknown(m)
}

var xxx_messageInfo_Replication_Rebalance proto.InternalMessageInfo

func (m *Replication_Rebalance) GetHighUsageAgents() []string {
	if m != nil {
		return m.HighUsageAgents
	}
	return nil
}

func (m *Replication_Rebalance) GetLowUsageAgents() []string {
	if m != nil {
		return m.LowUsageAgents
	}
	return nil
}

type Replication_Agents struct {
	Agents               []string `protobuf:"bytes,1,rep,name=agents,proto3" json:"agents,omitempty"`
	RemovedAgents        []string `protobuf:"bytes,2,rep,name=removed_agents,json=removedAgents,proto3" json:"removed_agents,omitempty"`
	ReplicatingAgent     []string `protobuf:"bytes,3,rep,name=replicating_agent,json=replicatingAgent,proto3" json:"replicating_agent,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Replication_Agents) Reset()         { *m = Replication_Agents{} }
func (m *Replication_Agents) String() string { return proto.CompactTextString(m) }
func (*Replication_Agents) ProtoMessage()    {}
func (*Replication_Agents) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{4, 2}
}
func (m *Replication_Agents) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Replication_Agents) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Replication_Agents.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Replication_Agents) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Replication_Agents.Merge(m, src)
}
func (m *Replication_Agents) XXX_Size() int {
	return m.Size()
}
func (m *Replication_Agents) XXX_DiscardUnknown() {
	xxx_messageInfo_Replication_Agents.DiscardUnknown(m)
}

var xxx_messageInfo_Replication_Agents proto.InternalMessageInfo

func (m *Replication_Agents) GetAgents() []string {
	if m != nil {
		return m.Agents
	}
	return nil
}

func (m *Replication_Agents) GetRemovedAgents() []string {
	if m != nil {
		return m.RemovedAgents
	}
	return nil
}

func (m *Replication_Agents) GetReplicatingAgent() []string {
	if m != nil {
		return m.ReplicatingAgent
	}
	return nil
}

type Discoverer struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Discoverer) Reset()         { *m = Discoverer{} }
func (m *Discoverer) String() string { return proto.CompactTextString(m) }
func (*Discoverer) ProtoMessage()    {}
func (*Discoverer) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{5}
}
func (m *Discoverer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Discoverer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Discoverer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Discoverer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Discoverer.Merge(m, src)
}
func (m *Discoverer) XXX_Size() int {
	return m.Size()
}
func (m *Discoverer) XXX_DiscardUnknown() {
	xxx_messageInfo_Discoverer.DiscardUnknown(m)
}

var xxx_messageInfo_Discoverer proto.InternalMessageInfo

type Discoverer_Request struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace            string   `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Node                 string   `protobuf:"bytes,3,opt,name=node,proto3" json:"node,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Discoverer_Request) Reset()         { *m = Discoverer_Request{} }
func (m *Discoverer_Request) String() string { return proto.CompactTextString(m) }
func (*Discoverer_Request) ProtoMessage()    {}
func (*Discoverer_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{5, 0}
}
func (m *Discoverer_Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Discoverer_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Discoverer_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Discoverer_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Discoverer_Request.Merge(m, src)
}
func (m *Discoverer_Request) XXX_Size() int {
	return m.Size()
}
func (m *Discoverer_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Discoverer_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Discoverer_Request proto.InternalMessageInfo

func (m *Discoverer_Request) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Discoverer_Request) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Discoverer_Request) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

type Backup struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup) Reset()         { *m = Backup{} }
func (m *Backup) String() string { return proto.CompactTextString(m) }
func (*Backup) ProtoMessage()    {}
func (*Backup) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6}
}
func (m *Backup) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup.Merge(m, src)
}
func (m *Backup) XXX_Size() int {
	return m.Size()
}
func (m *Backup) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup.DiscardUnknown(m)
}

var xxx_messageInfo_Backup proto.InternalMessageInfo

type Backup_GetVector struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_GetVector) Reset()         { *m = Backup_GetVector{} }
func (m *Backup_GetVector) String() string { return proto.CompactTextString(m) }
func (*Backup_GetVector) ProtoMessage()    {}
func (*Backup_GetVector) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 0}
}
func (m *Backup_GetVector) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_GetVector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_GetVector.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_GetVector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_GetVector.Merge(m, src)
}
func (m *Backup_GetVector) XXX_Size() int {
	return m.Size()
}
func (m *Backup_GetVector) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_GetVector.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_GetVector proto.InternalMessageInfo

type Backup_GetVector_Request struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_GetVector_Request) Reset()         { *m = Backup_GetVector_Request{} }
func (m *Backup_GetVector_Request) String() string { return proto.CompactTextString(m) }
func (*Backup_GetVector_Request) ProtoMessage()    {}
func (*Backup_GetVector_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 0, 0}
}
func (m *Backup_GetVector_Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_GetVector_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_GetVector_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_GetVector_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_GetVector_Request.Merge(m, src)
}
func (m *Backup_GetVector_Request) XXX_Size() int {
	return m.Size()
}
func (m *Backup_GetVector_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_GetVector_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_GetVector_Request proto.InternalMessageInfo

func (m *Backup_GetVector_Request) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type Backup_GetVector_Owner struct {
	Ip                   string   `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_GetVector_Owner) Reset()         { *m = Backup_GetVector_Owner{} }
func (m *Backup_GetVector_Owner) String() string { return proto.CompactTextString(m) }
func (*Backup_GetVector_Owner) ProtoMessage()    {}
func (*Backup_GetVector_Owner) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 0, 1}
}
func (m *Backup_GetVector_Owner) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_GetVector_Owner) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_GetVector_Owner.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_GetVector_Owner) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_GetVector_Owner.Merge(m, src)
}
func (m *Backup_GetVector_Owner) XXX_Size() int {
	return m.Size()
}
func (m *Backup_GetVector_Owner) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_GetVector_Owner.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_GetVector_Owner proto.InternalMessageInfo

func (m *Backup_GetVector_Owner) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

type Backup_Locations struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_Locations) Reset()         { *m = Backup_Locations{} }
func (m *Backup_Locations) String() string { return proto.CompactTextString(m) }
func (*Backup_Locations) ProtoMessage()    {}
func (*Backup_Locations) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 1}
}
func (m *Backup_Locations) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Locations) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Locations.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Locations) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Locations.Merge(m, src)
}
func (m *Backup_Locations) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Locations) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Locations.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Locations proto.InternalMessageInfo

type Backup_Locations_Request struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_Locations_Request) Reset()         { *m = Backup_Locations_Request{} }
func (m *Backup_Locations_Request) String() string { return proto.CompactTextString(m) }
func (*Backup_Locations_Request) ProtoMessage()    {}
func (*Backup_Locations_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 1, 0}
}
func (m *Backup_Locations_Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Locations_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Locations_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Locations_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Locations_Request.Merge(m, src)
}
func (m *Backup_Locations_Request) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Locations_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Locations_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Locations_Request proto.InternalMessageInfo

func (m *Backup_Locations_Request) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type Backup_Remove struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_Remove) Reset()         { *m = Backup_Remove{} }
func (m *Backup_Remove) String() string { return proto.CompactTextString(m) }
func (*Backup_Remove) ProtoMessage()    {}
func (*Backup_Remove) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 2}
}
func (m *Backup_Remove) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Remove) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Remove.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Remove) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Remove.Merge(m, src)
}
func (m *Backup_Remove) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Remove) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Remove.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Remove proto.InternalMessageInfo

type Backup_Remove_Request struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_Remove_Request) Reset()         { *m = Backup_Remove_Request{} }
func (m *Backup_Remove_Request) String() string { return proto.CompactTextString(m) }
func (*Backup_Remove_Request) ProtoMessage()    {}
func (*Backup_Remove_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 2, 0}
}
func (m *Backup_Remove_Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Remove_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Remove_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Remove_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Remove_Request.Merge(m, src)
}
func (m *Backup_Remove_Request) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Remove_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Remove_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Remove_Request proto.InternalMessageInfo

func (m *Backup_Remove_Request) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type Backup_Remove_RequestMulti struct {
	Uuids                []string `protobuf:"bytes,1,rep,name=uuids,proto3" json:"uuids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_Remove_RequestMulti) Reset()         { *m = Backup_Remove_RequestMulti{} }
func (m *Backup_Remove_RequestMulti) String() string { return proto.CompactTextString(m) }
func (*Backup_Remove_RequestMulti) ProtoMessage()    {}
func (*Backup_Remove_RequestMulti) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 2, 1}
}
func (m *Backup_Remove_RequestMulti) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Remove_RequestMulti) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Remove_RequestMulti.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Remove_RequestMulti) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Remove_RequestMulti.Merge(m, src)
}
func (m *Backup_Remove_RequestMulti) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Remove_RequestMulti) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Remove_RequestMulti.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Remove_RequestMulti proto.InternalMessageInfo

func (m *Backup_Remove_RequestMulti) GetUuids() []string {
	if m != nil {
		return m.Uuids
	}
	return nil
}

type Backup_IP struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_IP) Reset()         { *m = Backup_IP{} }
func (m *Backup_IP) String() string { return proto.CompactTextString(m) }
func (*Backup_IP) ProtoMessage()    {}
func (*Backup_IP) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 3}
}
func (m *Backup_IP) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_IP) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_IP.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_IP) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_IP.Merge(m, src)
}
func (m *Backup_IP) XXX_Size() int {
	return m.Size()
}
func (m *Backup_IP) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_IP.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_IP proto.InternalMessageInfo

type Backup_IP_Register struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_IP_Register) Reset()         { *m = Backup_IP_Register{} }
func (m *Backup_IP_Register) String() string { return proto.CompactTextString(m) }
func (*Backup_IP_Register) ProtoMessage()    {}
func (*Backup_IP_Register) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 3, 0}
}
func (m *Backup_IP_Register) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_IP_Register) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_IP_Register.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_IP_Register) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_IP_Register.Merge(m, src)
}
func (m *Backup_IP_Register) XXX_Size() int {
	return m.Size()
}
func (m *Backup_IP_Register) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_IP_Register.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_IP_Register proto.InternalMessageInfo

type Backup_IP_Register_Request struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Ips                  []string `protobuf:"bytes,2,rep,name=ips,proto3" json:"ips,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_IP_Register_Request) Reset()         { *m = Backup_IP_Register_Request{} }
func (m *Backup_IP_Register_Request) String() string { return proto.CompactTextString(m) }
func (*Backup_IP_Register_Request) ProtoMessage()    {}
func (*Backup_IP_Register_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 3, 0, 0}
}
func (m *Backup_IP_Register_Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_IP_Register_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_IP_Register_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_IP_Register_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_IP_Register_Request.Merge(m, src)
}
func (m *Backup_IP_Register_Request) XXX_Size() int {
	return m.Size()
}
func (m *Backup_IP_Register_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_IP_Register_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_IP_Register_Request proto.InternalMessageInfo

func (m *Backup_IP_Register_Request) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Backup_IP_Register_Request) GetIps() []string {
	if m != nil {
		return m.Ips
	}
	return nil
}

type Backup_IP_Remove struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_IP_Remove) Reset()         { *m = Backup_IP_Remove{} }
func (m *Backup_IP_Remove) String() string { return proto.CompactTextString(m) }
func (*Backup_IP_Remove) ProtoMessage()    {}
func (*Backup_IP_Remove) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 3, 1}
}
func (m *Backup_IP_Remove) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_IP_Remove) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_IP_Remove.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_IP_Remove) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_IP_Remove.Merge(m, src)
}
func (m *Backup_IP_Remove) XXX_Size() int {
	return m.Size()
}
func (m *Backup_IP_Remove) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_IP_Remove.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_IP_Remove proto.InternalMessageInfo

type Backup_IP_Remove_Request struct {
	Ips                  []string `protobuf:"bytes,1,rep,name=ips,proto3" json:"ips,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_IP_Remove_Request) Reset()         { *m = Backup_IP_Remove_Request{} }
func (m *Backup_IP_Remove_Request) String() string { return proto.CompactTextString(m) }
func (*Backup_IP_Remove_Request) ProtoMessage()    {}
func (*Backup_IP_Remove_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 3, 1, 0}
}
func (m *Backup_IP_Remove_Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_IP_Remove_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_IP_Remove_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_IP_Remove_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_IP_Remove_Request.Merge(m, src)
}
func (m *Backup_IP_Remove_Request) XXX_Size() int {
	return m.Size()
}
func (m *Backup_IP_Remove_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_IP_Remove_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_IP_Remove_Request proto.InternalMessageInfo

func (m *Backup_IP_Remove_Request) GetIps() []string {
	if m != nil {
		return m.Ips
	}
	return nil
}

type Backup_MetaVector struct {
	Uuid                 string    `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Vector               []float32 `protobuf:"fixed32,3,rep,packed,name=vector,proto3" json:"vector,omitempty"`
	Ips                  []string  `protobuf:"bytes,4,rep,name=ips,proto3" json:"ips,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Backup_MetaVector) Reset()         { *m = Backup_MetaVector{} }
func (m *Backup_MetaVector) String() string { return proto.CompactTextString(m) }
func (*Backup_MetaVector) ProtoMessage()    {}
func (*Backup_MetaVector) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 4}
}
func (m *Backup_MetaVector) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_MetaVector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_MetaVector.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_MetaVector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_MetaVector.Merge(m, src)
}
func (m *Backup_MetaVector) XXX_Size() int {
	return m.Size()
}
func (m *Backup_MetaVector) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_MetaVector.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_MetaVector proto.InternalMessageInfo

func (m *Backup_MetaVector) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Backup_MetaVector) GetVector() []float32 {
	if m != nil {
		return m.Vector
	}
	return nil
}

func (m *Backup_MetaVector) GetIps() []string {
	if m != nil {
		return m.Ips
	}
	return nil
}

type Backup_MetaVectors struct {
	Vectors              []*Backup_MetaVector `protobuf:"bytes,1,rep,name=vectors,proto3" json:"vectors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Backup_MetaVectors) Reset()         { *m = Backup_MetaVectors{} }
func (m *Backup_MetaVectors) String() string { return proto.CompactTextString(m) }
func (*Backup_MetaVectors) ProtoMessage()    {}
func (*Backup_MetaVectors) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 5}
}
func (m *Backup_MetaVectors) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_MetaVectors) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_MetaVectors.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_MetaVectors) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_MetaVectors.Merge(m, src)
}
func (m *Backup_MetaVectors) XXX_Size() int {
	return m.Size()
}
func (m *Backup_MetaVectors) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_MetaVectors.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_MetaVectors proto.InternalMessageInfo

func (m *Backup_MetaVectors) GetVectors() []*Backup_MetaVector {
	if m != nil {
		return m.Vectors
	}
	return nil
}

type Backup_Compressed struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_Compressed) Reset()         { *m = Backup_Compressed{} }
func (m *Backup_Compressed) String() string { return proto.CompactTextString(m) }
func (*Backup_Compressed) ProtoMessage()    {}
func (*Backup_Compressed) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 6}
}
func (m *Backup_Compressed) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Compressed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Compressed.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Compressed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Compressed.Merge(m, src)
}
func (m *Backup_Compressed) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Compressed) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Compressed.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Compressed proto.InternalMessageInfo

type Backup_Compressed_MetaVector struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Vector               []byte   `protobuf:"bytes,3,opt,name=vector,proto3" json:"vector,omitempty"`
	Ips                  []string `protobuf:"bytes,4,rep,name=ips,proto3" json:"ips,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_Compressed_MetaVector) Reset()         { *m = Backup_Compressed_MetaVector{} }
func (m *Backup_Compressed_MetaVector) String() string { return proto.CompactTextString(m) }
func (*Backup_Compressed_MetaVector) ProtoMessage()    {}
func (*Backup_Compressed_MetaVector) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 6, 0}
}
func (m *Backup_Compressed_MetaVector) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Compressed_MetaVector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Compressed_MetaVector.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Compressed_MetaVector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Compressed_MetaVector.Merge(m, src)
}
func (m *Backup_Compressed_MetaVector) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Compressed_MetaVector) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Compressed_MetaVector.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Compressed_MetaVector proto.InternalMessageInfo

func (m *Backup_Compressed_MetaVector) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Backup_Compressed_MetaVector) GetVector() []byte {
	if m != nil {
		return m.Vector
	}
	return nil
}

func (m *Backup_Compressed_MetaVector) GetIps() []string {
	if m != nil {
		return m.Ips
	}
	return nil
}

type Backup_Compressed_MetaVectors struct {
	Vectors              []*Backup_Compressed_MetaVector `protobuf:"bytes,1,rep,name=vectors,proto3" json:"vectors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *Backup_Compressed_MetaVectors) Reset()         { *m = Backup_Compressed_MetaVectors{} }
func (m *Backup_Compressed_MetaVectors) String() string { return proto.CompactTextString(m) }
func (*Backup_Compressed_MetaVectors) ProtoMessage()    {}
func (*Backup_Compressed_MetaVectors) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{6, 6, 1}
}
func (m *Backup_Compressed_MetaVectors) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Compressed_MetaVectors) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Compressed_MetaVectors.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Compressed_MetaVectors) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Compressed_MetaVectors.Merge(m, src)
}
func (m *Backup_Compressed_MetaVectors) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Compressed_MetaVectors) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Compressed_MetaVectors.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Compressed_MetaVectors proto.InternalMessageInfo

func (m *Backup_Compressed_MetaVectors) GetVectors() []*Backup_Compressed_MetaVector {
	if m != nil {
		return m.Vectors
	}
	return nil
}

type Info struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Info) Reset()         { *m = Info{} }
func (m *Info) String() string { return proto.CompactTextString(m) }
func (*Info) ProtoMessage()    {}
func (*Info) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7}
}
func (m *Info) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info.Merge(m, src)
}
func (m *Info) XXX_Size() int {
	return m.Size()
}
func (m *Info) XXX_DiscardUnknown() {
	xxx_messageInfo_Info.DiscardUnknown(m)
}

var xxx_messageInfo_Info proto.InternalMessageInfo

type Info_Index struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Info_Index) Reset()         { *m = Info_Index{} }
func (m *Info_Index) String() string { return proto.CompactTextString(m) }
func (*Info_Index) ProtoMessage()    {}
func (*Info_Index) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 0}
}
func (m *Info_Index) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_Index) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_Index.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_Index) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_Index.Merge(m, src)
}
func (m *Info_Index) XXX_Size() int {
	return m.Size()
}
func (m *Info_Index) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_Index.DiscardUnknown(m)
}

var xxx_messageInfo_Info_Index proto.InternalMessageInfo

type Info_Index_Count struct {
	Stored               uint32   `protobuf:"varint,1,opt,name=stored,proto3" json:"stored,omitempty"`
	Uncommitted          uint32   `protobuf:"varint,2,opt,name=uncommitted,proto3" json:"uncommitted,omitempty"`
	Indexing             bool     `protobuf:"varint,3,opt,name=indexing,proto3" json:"indexing,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Info_Index_Count) Reset()         { *m = Info_Index_Count{} }
func (m *Info_Index_Count) String() string { return proto.CompactTextString(m) }
func (*Info_Index_Count) ProtoMessage()    {}
func (*Info_Index_Count) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 0, 0}
}
func (m *Info_Index_Count) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_Index_Count) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_Index_Count.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_Index_Count) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_Index_Count.Merge(m, src)
}
func (m *Info_Index_Count) XXX_Size() int {
	return m.Size()
}
func (m *Info_Index_Count) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_Index_Count.DiscardUnknown(m)
}

var xxx_messageInfo_Info_Index_Count proto.InternalMessageInfo

func (m *Info_Index_Count) GetStored() uint32 {
	if m != nil {
		return m.Stored
	}
	return 0
}

func (m *Info_Index_Count) GetUncommitted() uint32 {
	if m != nil {
		return m.Uncommitted
	}
	return 0
}

func (m *Info_Index_Count) GetIndexing() bool {
	if m != nil {
		return m.Indexing
	}
	return false
}

type Info_Index_UUID struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Info_Index_UUID) Reset()         { *m = Info_Index_UUID{} }
func (m *Info_Index_UUID) String() string { return proto.CompactTextString(m) }
func (*Info_Index_UUID) ProtoMessage()    {}
func (*Info_Index_UUID) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 0, 1}
}
func (m *Info_Index_UUID) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_Index_UUID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_Index_UUID.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_Index_UUID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_Index_UUID.Merge(m, src)
}
func (m *Info_Index_UUID) XXX_Size() int {
	return m.Size()
}
func (m *Info_Index_UUID) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_Index_UUID.DiscardUnknown(m)
}

var xxx_messageInfo_Info_Index_UUID proto.InternalMessageInfo

type Info_Index_UUID_Committed struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Info_Index_UUID_Committed) Reset()         { *m = Info_Index_UUID_Committed{} }
func (m *Info_Index_UUID_Committed) String() string { return proto.CompactTextString(m) }
func (*Info_Index_UUID_Committed) ProtoMessage()    {}
func (*Info_Index_UUID_Committed) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 0, 1, 0}
}
func (m *Info_Index_UUID_Committed) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_Index_UUID_Committed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_Index_UUID_Committed.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_Index_UUID_Committed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_Index_UUID_Committed.Merge(m, src)
}
func (m *Info_Index_UUID_Committed) XXX_Size() int {
	return m.Size()
}
func (m *Info_Index_UUID_Committed) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_Index_UUID_Committed.DiscardUnknown(m)
}

var xxx_messageInfo_Info_Index_UUID_Committed proto.InternalMessageInfo

func (m *Info_Index_UUID_Committed) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type Info_Index_UUID_Uncommitted struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Info_Index_UUID_Uncommitted) Reset()         { *m = Info_Index_UUID_Uncommitted{} }
func (m *Info_Index_UUID_Uncommitted) String() string { return proto.CompactTextString(m) }
func (*Info_Index_UUID_Uncommitted) ProtoMessage()    {}
func (*Info_Index_UUID_Uncommitted) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 0, 1, 1}
}
func (m *Info_Index_UUID_Uncommitted) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_Index_UUID_Uncommitted) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_Index_UUID_Uncommitted.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_Index_UUID_Uncommitted) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_Index_UUID_Uncommitted.Merge(m, src)
}
func (m *Info_Index_UUID_Uncommitted) XXX_Size() int {
	return m.Size()
}
func (m *Info_Index_UUID_Uncommitted) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_Index_UUID_Uncommitted.DiscardUnknown(m)
}

var xxx_messageInfo_Info_Index_UUID_Uncommitted proto.InternalMessageInfo

func (m *Info_Index_UUID_Uncommitted) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type Info_Pod struct {
	AppName              string       `protobuf:"bytes,1,opt,name=app_name,json=appName,proto3" json:"app_name,omitempty"`
	Name                 string       `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Namespace            string       `protobuf:"bytes,3,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Ip                   string       `protobuf:"bytes,4,opt,name=ip,proto3" json:"ip,omitempty"`
	Cpu                  *Info_CPU    `protobuf:"bytes,5,opt,name=cpu,proto3" json:"cpu,omitempty"`
	Memory               *Info_Memory `protobuf:"bytes,6,opt,name=memory,proto3" json:"memory,omitempty"`
	Node                 *Info_Node   `protobuf:"bytes,7,opt,name=node,proto3" json:"node,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Info_Pod) Reset()         { *m = Info_Pod{} }
func (m *Info_Pod) String() string { return proto.CompactTextString(m) }
func (*Info_Pod) ProtoMessage()    {}
func (*Info_Pod) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 1}
}
func (m *Info_Pod) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_Pod) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_Pod.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_Pod) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_Pod.Merge(m, src)
}
func (m *Info_Pod) XXX_Size() int {
	return m.Size()
}
func (m *Info_Pod) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_Pod.DiscardUnknown(m)
}

var xxx_messageInfo_Info_Pod proto.InternalMessageInfo

func (m *Info_Pod) GetAppName() string {
	if m != nil {
		return m.AppName
	}
	return ""
}

func (m *Info_Pod) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Info_Pod) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Info_Pod) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *Info_Pod) GetCpu() *Info_CPU {
	if m != nil {
		return m.Cpu
	}
	return nil
}

func (m *Info_Pod) GetMemory() *Info_Memory {
	if m != nil {
		return m.Memory
	}
	return nil
}

func (m *Info_Pod) GetNode() *Info_Node {
	if m != nil {
		return m.Node
	}
	return nil
}

type Info_Node struct {
	Name                 string       `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	InternalAddr         string       `protobuf:"bytes,2,opt,name=internal_addr,json=internalAddr,proto3" json:"internal_addr,omitempty"`
	ExternalAddr         string       `protobuf:"bytes,3,opt,name=external_addr,json=externalAddr,proto3" json:"external_addr,omitempty"`
	Cpu                  *Info_CPU    `protobuf:"bytes,4,opt,name=cpu,proto3" json:"cpu,omitempty"`
	Memory               *Info_Memory `protobuf:"bytes,5,opt,name=memory,proto3" json:"memory,omitempty"`
	Pods                 *Info_Pods   `protobuf:"bytes,6,opt,name=Pods,proto3" json:"Pods,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Info_Node) Reset()         { *m = Info_Node{} }
func (m *Info_Node) String() string { return proto.CompactTextString(m) }
func (*Info_Node) ProtoMessage()    {}
func (*Info_Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 2}
}
func (m *Info_Node) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_Node.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_Node.Merge(m, src)
}
func (m *Info_Node) XXX_Size() int {
	return m.Size()
}
func (m *Info_Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Info_Node proto.InternalMessageInfo

func (m *Info_Node) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Info_Node) GetInternalAddr() string {
	if m != nil {
		return m.InternalAddr
	}
	return ""
}

func (m *Info_Node) GetExternalAddr() string {
	if m != nil {
		return m.ExternalAddr
	}
	return ""
}

func (m *Info_Node) GetCpu() *Info_CPU {
	if m != nil {
		return m.Cpu
	}
	return nil
}

func (m *Info_Node) GetMemory() *Info_Memory {
	if m != nil {
		return m.Memory
	}
	return nil
}

func (m *Info_Node) GetPods() *Info_Pods {
	if m != nil {
		return m.Pods
	}
	return nil
}

type Info_CPU struct {
	Limit                float64  `protobuf:"fixed64,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Request              float64  `protobuf:"fixed64,2,opt,name=request,proto3" json:"request,omitempty"`
	Usage                float64  `protobuf:"fixed64,3,opt,name=usage,proto3" json:"usage,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Info_CPU) Reset()         { *m = Info_CPU{} }
func (m *Info_CPU) String() string { return proto.CompactTextString(m) }
func (*Info_CPU) ProtoMessage()    {}
func (*Info_CPU) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 3}
}
func (m *Info_CPU) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_CPU) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_CPU.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_CPU) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_CPU.Merge(m, src)
}
func (m *Info_CPU) XXX_Size() int {
	return m.Size()
}
func (m *Info_CPU) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_CPU.DiscardUnknown(m)
}

var xxx_messageInfo_Info_CPU proto.InternalMessageInfo

func (m *Info_CPU) GetLimit() float64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *Info_CPU) GetRequest() float64 {
	if m != nil {
		return m.Request
	}
	return 0
}

func (m *Info_CPU) GetUsage() float64 {
	if m != nil {
		return m.Usage
	}
	return 0
}

type Info_Memory struct {
	Limit                float64  `protobuf:"fixed64,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Request              float64  `protobuf:"fixed64,2,opt,name=request,proto3" json:"request,omitempty"`
	Usage                float64  `protobuf:"fixed64,3,opt,name=usage,proto3" json:"usage,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Info_Memory) Reset()         { *m = Info_Memory{} }
func (m *Info_Memory) String() string { return proto.CompactTextString(m) }
func (*Info_Memory) ProtoMessage()    {}
func (*Info_Memory) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 4}
}
func (m *Info_Memory) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_Memory) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_Memory.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_Memory) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_Memory.Merge(m, src)
}
func (m *Info_Memory) XXX_Size() int {
	return m.Size()
}
func (m *Info_Memory) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_Memory.DiscardUnknown(m)
}

var xxx_messageInfo_Info_Memory proto.InternalMessageInfo

func (m *Info_Memory) GetLimit() float64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *Info_Memory) GetRequest() float64 {
	if m != nil {
		return m.Request
	}
	return 0
}

func (m *Info_Memory) GetUsage() float64 {
	if m != nil {
		return m.Usage
	}
	return 0
}

type Info_Pods struct {
	Pods                 []*Info_Pod `protobuf:"bytes,1,rep,name=pods,proto3" json:"pods,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Info_Pods) Reset()         { *m = Info_Pods{} }
func (m *Info_Pods) String() string { return proto.CompactTextString(m) }
func (*Info_Pods) ProtoMessage()    {}
func (*Info_Pods) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 5}
}
func (m *Info_Pods) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_Pods) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_Pods.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_Pods) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_Pods.Merge(m, src)
}
func (m *Info_Pods) XXX_Size() int {
	return m.Size()
}
func (m *Info_Pods) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_Pods.DiscardUnknown(m)
}

var xxx_messageInfo_Info_Pods proto.InternalMessageInfo

func (m *Info_Pods) GetPods() []*Info_Pod {
	if m != nil {
		return m.Pods
	}
	return nil
}

type Info_Nodes struct {
	Nodes                []*Info_Node `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Info_Nodes) Reset()         { *m = Info_Nodes{} }
func (m *Info_Nodes) String() string { return proto.CompactTextString(m) }
func (*Info_Nodes) ProtoMessage()    {}
func (*Info_Nodes) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 6}
}
func (m *Info_Nodes) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_Nodes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_Nodes.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_Nodes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_Nodes.Merge(m, src)
}
func (m *Info_Nodes) XXX_Size() int {
	return m.Size()
}
func (m *Info_Nodes) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_Nodes.DiscardUnknown(m)
}

var xxx_messageInfo_Info_Nodes proto.InternalMessageInfo

func (m *Info_Nodes) GetNodes() []*Info_Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type Info_IPs struct {
	Ip                   []string `protobuf:"bytes,1,rep,name=ip,proto3" json:"ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Info_IPs) Reset()         { *m = Info_IPs{} }
func (m *Info_IPs) String() string { return proto.CompactTextString(m) }
func (*Info_IPs) ProtoMessage()    {}
func (*Info_IPs) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{7, 7}
}
func (m *Info_IPs) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Info_IPs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Info_IPs.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Info_IPs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Info_IPs.Merge(m, src)
}
func (m *Info_IPs) XXX_Size() int {
	return m.Size()
}
func (m *Info_IPs) XXX_DiscardUnknown() {
	xxx_messageInfo_Info_IPs.DiscardUnknown(m)
}

var xxx_messageInfo_Info_IPs proto.InternalMessageInfo

func (m *Info_IPs) GetIp() []string {
	if m != nil {
		return m.Ip
	}
	return nil
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_678c914f1bee6d56, []int{8}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return m.Size()
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Search)(nil), "payload.Search")
	proto.RegisterType((*Search_Request)(nil), "payload.Search.Request")
	proto.RegisterType((*Search_IDRequest)(nil), "payload.Search.IDRequest")
	proto.RegisterType((*Search_ObjectRequest)(nil), "payload.Search.ObjectRequest")
	proto.RegisterType((*Search_Config)(nil), "payload.Search.Config")
	proto.RegisterType((*Search_Response)(nil), "payload.Search.Response")
	proto.RegisterType((*Meta)(nil), "payload.Meta")
	proto.RegisterType((*Meta_Key)(nil), "payload.Meta.Key")
	proto.RegisterType((*Meta_Keys)(nil), "payload.Meta.Keys")
	proto.RegisterType((*Meta_Val)(nil), "payload.Meta.Val")
	proto.RegisterType((*Meta_Vals)(nil), "payload.Meta.Vals")
	proto.RegisterType((*Meta_KeyVal)(nil), "payload.Meta.KeyVal")
	proto.RegisterType((*Meta_KeyVals)(nil), "payload.Meta.KeyVals")
	proto.RegisterType((*Object)(nil), "payload.Object")
	proto.RegisterType((*Object_Distance)(nil), "payload.Object.Distance")
	proto.RegisterType((*Object_ID)(nil), "payload.Object.ID")
	proto.RegisterType((*Object_IDs)(nil), "payload.Object.IDs")
	proto.RegisterType((*Object_Vector)(nil), "payload.Object.Vector")
	proto.RegisterType((*Object_Vectors)(nil), "payload.Object.Vectors")
	proto.RegisterType((*Object_Blob)(nil), "payload.Object.Blob")
	proto.RegisterType((*Object_Location)(nil), "payload.Object.Location")
	proto.RegisterType((*Object_Locations)(nil), "payload.Object.Locations")
	proto.RegisterType((*Control)(nil), "payload.Control")
	proto.RegisterType((*Control_CreateIndexRequest)(nil), "payload.Control.CreateIndexRequest")
	proto.RegisterType((*Replication)(nil), "payload.Replication")
	proto.RegisterType((*Replication_Recovery)(nil), "payload.Replication.Recovery")
	proto.RegisterType((*Replication_Rebalance)(nil), "payload.Replication.Rebalance")
	proto.RegisterType((*Replication_Agents)(nil), "payload.Replication.Agents")
	proto.RegisterType((*Discoverer)(nil), "payload.Discoverer")
	proto.RegisterType((*Discoverer_Request)(nil), "payload.Discoverer.Request")
	proto.RegisterType((*Backup)(nil), "payload.Backup")
	proto.RegisterType((*Backup_GetVector)(nil), "payload.Backup.GetVector")
	proto.RegisterType((*Backup_GetVector_Request)(nil), "payload.Backup.GetVector.Request")
	proto.RegisterType((*Backup_GetVector_Owner)(nil), "payload.Backup.GetVector.Owner")
	proto.RegisterType((*Backup_Locations)(nil), "payload.Backup.Locations")
	proto.RegisterType((*Backup_Locations_Request)(nil), "payload.Backup.Locations.Request")
	proto.RegisterType((*Backup_Remove)(nil), "payload.Backup.Remove")
	proto.RegisterType((*Backup_Remove_Request)(nil), "payload.Backup.Remove.Request")
	proto.RegisterType((*Backup_Remove_RequestMulti)(nil), "payload.Backup.Remove.RequestMulti")
	proto.RegisterType((*Backup_IP)(nil), "payload.Backup.IP")
	proto.RegisterType((*Backup_IP_Register)(nil), "payload.Backup.IP.Register")
	proto.RegisterType((*Backup_IP_Register_Request)(nil), "payload.Backup.IP.Register.Request")
	proto.RegisterType((*Backup_IP_Remove)(nil), "payload.Backup.IP.Remove")
	proto.RegisterType((*Backup_IP_Remove_Request)(nil), "payload.Backup.IP.Remove.Request")
	proto.RegisterType((*Backup_MetaVector)(nil), "payload.Backup.MetaVector")
	proto.RegisterType((*Backup_MetaVectors)(nil), "payload.Backup.MetaVectors")
	proto.RegisterType((*Backup_Compressed)(nil), "payload.Backup.Compressed")
	proto.RegisterType((*Backup_Compressed_MetaVector)(nil), "payload.Backup.Compressed.MetaVector")
	proto.RegisterType((*Backup_Compressed_MetaVectors)(nil), "payload.Backup.Compressed.MetaVectors")
	proto.RegisterType((*Info)(nil), "payload.Info")
	proto.RegisterType((*Info_Index)(nil), "payload.Info.Index")
	proto.RegisterType((*Info_Index_Count)(nil), "payload.Info.Index.Count")
	proto.RegisterType((*Info_Index_UUID)(nil), "payload.Info.Index.UUID")
	proto.RegisterType((*Info_Index_UUID_Committed)(nil), "payload.Info.Index.UUID.Committed")
	proto.RegisterType((*Info_Index_UUID_Uncommitted)(nil), "payload.Info.Index.UUID.Uncommitted")
	proto.RegisterType((*Info_Pod)(nil), "payload.Info.Pod")
	proto.RegisterType((*Info_Node)(nil), "payload.Info.Node")
	proto.RegisterType((*Info_CPU)(nil), "payload.Info.CPU")
	proto.RegisterType((*Info_Memory)(nil), "payload.Info.Memory")
	proto.RegisterType((*Info_Pods)(nil), "payload.Info.Pods")
	proto.RegisterType((*Info_Nodes)(nil), "payload.Info.Nodes")
	proto.RegisterType((*Info_IPs)(nil), "payload.Info.IPs")
	proto.RegisterType((*Empty)(nil), "payload.Empty")
}

func init() { proto.RegisterFile("payload.proto", fileDescriptor_678c914f1bee6d56) }

var fileDescriptor_678c914f1bee6d56 = []byte{
	// 1387 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x57, 0x4f, 0x6f, 0x1b, 0x45,
	0x1b, 0x7f, 0x77, 0xd7, 0x5e, 0xdb, 0x4f, 0xe2, 0xbe, 0xc9, 0xa8, 0x6f, 0xea, 0xce, 0x0b, 0x21,
	0xb8, 0xb4, 0x44, 0xa5, 0x72, 0x68, 0x8a, 0x5a, 0x89, 0x4a, 0xa0, 0xda, 0xae, 0x90, 0x1b, 0x9a,
	0x5a, 0x53, 0x1c, 0x10, 0x7f, 0x14, 0x26, 0xde, 0xa9, 0xb3, 0x64, 0xbd, 0xb3, 0xec, 0xec, 0xba,
	0x75, 0xef, 0x7c, 0x01, 0x6e, 0x7c, 0x03, 0x8e, 0x88, 0x33, 0xe2, 0xc4, 0x81, 0x23, 0x7c, 0x02,
	0x50, 0x0f, 0x7c, 0x88, 0x9e, 0xd0, 0xfc, 0x59, 0xef, 0xda, 0x75, 0x44, 0x2a, 0x71, 0xf2, 0x3e,
	0xcf, 0xfc, 0x9e, 0xdf, 0x3c, 0xff, 0x77, 0x0d, 0xf5, 0x88, 0x4e, 0x03, 0x4e, 0xbd, 0x56, 0x14,
	0xf3, 0x84, 0xa3, 0x8a, 0x11, 0xf1, 0x85, 0x09, 0x0d, 0x7c, 0x8f, 0x26, 0x6c, 0x27, 0x7b, 0xd0,
	0x88, 0xe6, 0x4f, 0x0e, 0xb8, 0x0f, 0x19, 0x8d, 0x87, 0xc7, 0xf8, 0x33, 0xa8, 0x10, 0xf6, 0x75,
	0xca, 0x44, 0x82, 0xb6, 0xc0, 0x9d, 0xb0, 0x61, 0xc2, 0xe3, 0x86, 0xb5, 0xe5, 0x6c, 0xdb, 0xed,
	0xea, 0xf3, 0x76, 0xf9, 0x5b, 0xcb, 0xae, 0xda, 0xc4, 0xe8, 0x51, 0x0b, 0xdc, 0x21, 0x0f, 0x1f,
	0xf9, 0xa3, 0x86, 0xbd, 0x65, 0x6d, 0xaf, 0xec, 0x6e, 0xb4, 0xb2, 0x9b, 0x35, 0x5b, 0xab, 0xa3,
	0x4e, 0x89, 0x41, 0xe1, 0x3d, 0xa8, 0xf5, 0xba, 0x19, 0xfd, 0x39, 0xb0, 0x7d, 0xaf, 0x61, 0x6d,
	0x59, 0xdb, 0x35, 0x62, 0xfb, 0xde, 0x4b, 0x93, 0x7d, 0x0c, 0xf5, 0x07, 0x47, 0x5f, 0xb1, 0x61,
	0x92, 0x11, 0x6e, 0x80, 0xcb, 0x95, 0x42, 0x91, 0xae, 0x12, 0x23, 0xbd, 0x34, 0x31, 0x07, 0x57,
	0x6b, 0xd0, 0x45, 0x70, 0xc2, 0x74, 0xac, 0xe8, 0xea, 0xed, 0xca, 0xf3, 0x76, 0xe9, 0xaa, 0xbd,
	0x6d, 0x11, 0xa9, 0x93, 0x97, 0xc5, 0xd4, 0xf3, 0x53, 0xa1, 0x48, 0x6d, 0x62, 0x24, 0xd4, 0x80,
	0x0a, 0x8b, 0x84, 0x1f, 0xf0, 0xb0, 0xe1, 0xa8, 0x83, 0x4c, 0x94, 0x27, 0x89, 0x3f, 0x66, 0x3c,
	0x4d, 0x1a, 0xa5, 0x2d, 0x6b, 0xdb, 0x21, 0x99, 0x88, 0xdf, 0x83, 0x2a, 0x61, 0x22, 0xe2, 0xa1,
	0x60, 0x68, 0x17, 0x2a, 0x31, 0x13, 0x69, 0x90, 0x08, 0x95, 0xf5, 0x95, 0xdd, 0xc6, 0xcc, 0x5b,
	0x1d, 0x6d, 0xab, 0xeb, 0x8b, 0x84, 0x86, 0x43, 0x46, 0x32, 0x60, 0xf3, 0x77, 0x0b, 0x4a, 0xf7,
	0x59, 0x42, 0xf1, 0x05, 0x70, 0xf6, 0xd8, 0x14, 0xad, 0x81, 0x73, 0xc2, 0xa6, 0x26, 0xb5, 0xf2,
	0x11, 0x63, 0x28, 0xed, 0xb1, 0xa9, 0x40, 0x08, 0x4a, 0x27, 0x6c, 0xaa, 0xa9, 0x6b, 0x44, 0x3d,
	0x4b, 0xa3, 0x03, 0x1a, 0x48, 0xa3, 0x09, 0x0d, 0x32, 0xa3, 0x09, 0x0d, 0xa4, 0xd1, 0x01, 0x0d,
	0x94, 0xd1, 0x84, 0x06, 0x33, 0x23, 0xf9, 0x8c, 0xaf, 0x81, 0xbb, 0xc7, 0xa6, 0xc6, 0x6e, 0xfe,
	0xb2, 0x8c, 0xc9, 0xce, 0x99, 0xae, 0x43, 0x45, 0xa3, 0x05, 0xba, 0x02, 0xce, 0xc9, 0x24, 0x8b,
	0xed, 0xfc, 0x2c, 0x36, 0xe9, 0x7e, 0x4b, 0x63, 0x88, 0x04, 0x34, 0x7f, 0x76, 0xc0, 0xd5, 0x01,
	0xe3, 0x9b, 0x50, 0xcd, 0x62, 0x7e, 0xa1, 0x69, 0x30, 0x54, 0x3d, 0x73, 0x66, 0x0a, 0x31, 0x93,
	0xf1, 0xab, 0x60, 0xf7, 0xba, 0xe8, 0x42, 0x6e, 0xa1, 0x4a, 0x18, 0xdb, 0x6b, 0x96, 0x34, 0x95,
	0x71, 0xf7, 0xba, 0x42, 0x7a, 0xeb, 0x7b, 0x59, 0x70, 0xf2, 0x11, 0x77, 0xc0, 0x3d, 0xd0, 0xfd,
	0x7d, 0x9a, 0x6d, 0x61, 0x34, 0xec, 0xe5, 0xa3, 0x81, 0x6f, 0x43, 0x45, 0x93, 0x08, 0xf4, 0x36,
	0x54, 0xb4, 0x32, 0x0b, 0x7b, 0x63, 0xb1, 0xa4, 0x1a, 0x49, 0x32, 0x18, 0xbe, 0x05, 0xa5, 0x76,
	0xc0, 0x8f, 0x4e, 0xbf, 0x3f, 0x6f, 0x75, 0xbb, 0xd8, 0xea, 0xb8, 0x0b, 0xd5, 0x0f, 0xf9, 0x90,
	0x26, 0x3e, 0x0f, 0x65, 0xd9, 0x42, 0x3a, 0x66, 0x26, 0x59, 0xea, 0x59, 0xea, 0xd2, 0xd4, 0xf7,
	0x4c, 0x6d, 0xd4, 0xb3, 0x4a, 0x40, 0x24, 0x1a, 0x8e, 0x49, 0x40, 0x24, 0x13, 0x50, 0xcb, 0x58,
	0x04, 0xba, 0x09, 0xb5, 0x20, 0x13, 0x4e, 0x6b, 0xc9, 0x0c, 0x4d, 0x72, 0x68, 0xf3, 0x2e, 0x54,
	0x3a, 0x3c, 0x4c, 0x62, 0x1e, 0xe0, 0x77, 0x01, 0x75, 0x62, 0x46, 0x13, 0xd6, 0x0b, 0x3d, 0xf6,
	0x24, 0x1b, 0xd7, 0x37, 0xa0, 0x16, 0x71, 0x1e, 0x1c, 0x0a, 0xff, 0x29, 0x9b, 0x1f, 0xb1, 0xff,
	0x90, 0xaa, 0x3c, 0x79, 0xe8, 0x3f, 0x65, 0xcd, 0xef, 0x6c, 0x58, 0x21, 0x2c, 0x0a, 0x7c, 0xcd,
	0x8b, 0xaf, 0xcb, 0x59, 0x19, 0xf2, 0x09, 0x8b, 0xa7, 0xe8, 0x32, 0x9c, 0xf3, 0x58, 0xc0, 0x12,
	0xe6, 0x1d, 0xd2, 0x11, 0x0b, 0x93, 0xac, 0x8a, 0x75, 0xa3, 0xbd, 0xa3, 0x94, 0x98, 0x42, 0x8d,
	0xb0, 0x23, 0x1a, 0xa8, 0x06, 0xba, 0x0a, 0xeb, 0xc7, 0xfe, 0xe8, 0xf8, 0x30, 0x15, 0x74, 0xc4,
	0xe6, 0xcd, 0xfe, 0x2b, 0x0f, 0x06, 0x52, 0xaf, 0x0d, 0xd1, 0x36, 0xac, 0x05, 0xfc, 0xf1, 0x3c,
	0xd4, 0x56, 0xd0, 0x73, 0x01, 0x7f, 0x5c, 0x40, 0xe2, 0x04, 0x5c, 0x63, 0xb3, 0x01, 0xee, 0x1c,
	0xa9, 0x91, 0xa4, 0xaf, 0x31, 0x1b, 0xf3, 0x49, 0xee, 0xab, 0x66, 0xaa, 0x1b, 0xad, 0x31, 0x7f,
	0x0b, 0xd6, 0xe3, 0x2c, 0xda, 0x70, 0xa4, 0xa1, 0xa6, 0x34, 0x6b, 0x85, 0x03, 0x85, 0x6e, 0x3e,
	0x02, 0xe8, 0xfa, 0x42, 0x25, 0x83, 0xc5, 0xf8, 0x93, 0x7c, 0x73, 0xff, 0xbf, 0x58, 0xfa, 0xbc,
	0x73, 0x74, 0x0f, 0xbc, 0x02, 0x35, 0xf9, 0x2b, 0x22, 0x6a, 0x66, 0xa6, 0x46, 0x72, 0x85, 0xea,
	0x1a, 0xee, 0x31, 0xb5, 0xbc, 0x64, 0xd7, 0x70, 0x8f, 0x35, 0x7f, 0x29, 0x83, 0xdb, 0xa6, 0xc3,
	0x93, 0x34, 0xc2, 0x03, 0xa8, 0x7d, 0xc0, 0x12, 0xdd, 0xaf, 0xf8, 0xca, 0xdc, 0x8d, 0xaa, 0xb1,
	0x16, 0x6f, 0x94, 0x4a, 0xbc, 0x05, 0xe5, 0x07, 0x8f, 0x43, 0xa6, 0xe7, 0x29, 0x7a, 0xb1, 0x9f,
	0x23, 0x7c, 0xa3, 0xd0, 0x71, 0x67, 0xa6, 0xfd, 0x12, 0x5c, 0xa2, 0x92, 0x77, 0x66, 0x8b, 0x16,
	0xac, 0x1a, 0xdc, 0xfd, 0x34, 0x48, 0x7c, 0xb4, 0x09, 0x65, 0xa9, 0x37, 0xb5, 0x9a, 0x4d, 0xb1,
	0x45, 0xb4, 0x1a, 0x7f, 0x63, 0x81, 0xdd, 0xeb, 0xe3, 0x7d, 0xd9, 0x73, 0x23, 0x5f, 0x24, 0x2c,
	0xc6, 0xed, 0xb3, 0x5d, 0x85, 0xb0, 0x9e, 0x2a, 0x7b, 0x81, 0x58, 0xcd, 0xd7, 0xce, 0xcc, 0xf1,
	0xcb, 0x39, 0x9b, 0x31, 0xb0, 0x96, 0x19, 0x7c, 0x04, 0x20, 0x17, 0xa4, 0xd9, 0x4a, 0xa8, 0x78,
	0xaf, 0xb9, 0x2e, 0x5f, 0x48, 0xce, 0x29, 0xef, 0x6a, 0x33, 0xe6, 0xa5, 0xe2, 0x98, 0xaf, 0xe4,
	0xac, 0x02, 0xbd, 0xb3, 0xb8, 0xa6, 0xf0, 0x6c, 0xcc, 0x75, 0xf1, 0x5b, 0x39, 0x3a, 0x5f, 0x55,
	0xdf, 0x5b, 0x00, 0x1d, 0x3e, 0x8e, 0x62, 0x26, 0x04, 0xf3, 0xf0, 0xbd, 0x7f, 0xf4, 0x74, 0xa3,
	0xe0, 0xa9, 0x5a, 0x5d, 0xa7, 0xfa, 0xb7, 0x3f, 0xef, 0xdf, 0xfb, 0x8b, 0xfe, 0x5d, 0x5e, 0xf4,
	0x2f, 0xf7, 0x63, 0x99, 0xab, 0xcd, 0x1f, 0x2a, 0x50, 0xea, 0x85, 0x8f, 0x38, 0xfe, 0xd1, 0x82,
	0xb2, 0x5a, 0x45, 0xf8, 0x0b, 0x28, 0x77, 0x78, 0x1a, 0xaa, 0x6f, 0x07, 0x91, 0xf0, 0x98, 0x69,
	0x5f, 0xeb, 0xc4, 0x48, 0x68, 0x0b, 0x56, 0xd2, 0x70, 0xc8, 0xc7, 0x63, 0x3f, 0x49, 0x98, 0xde,
	0x9b, 0x75, 0x52, 0x54, 0xc9, 0x37, 0x90, 0x2f, 0xb9, 0xfc, 0x70, 0xa4, 0x22, 0xaa, 0x92, 0x99,
	0x8c, 0xef, 0x41, 0x69, 0x30, 0xe8, 0x75, 0xf1, 0x6b, 0x50, 0xeb, 0xcc, 0x0c, 0x96, 0x24, 0x05,
	0xbf, 0x0e, 0x2b, 0x83, 0x02, 0xe7, 0x32, 0xc8, 0x5f, 0x16, 0x38, 0x7d, 0xee, 0xa1, 0x8b, 0x50,
	0xa5, 0x51, 0x74, 0x58, 0x58, 0xed, 0x15, 0x1a, 0x45, 0xfb, 0x66, 0xbb, 0x2b, 0xb5, 0x5d, 0xd8,
	0xf8, 0x73, 0xd3, 0xee, 0x2c, 0x4e, 0xbb, 0x1e, 0xc8, 0x52, 0xa1, 0x81, 0x9f, 0xa8, 0x81, 0x44,
	0x97, 0xc0, 0x19, 0x46, 0x69, 0xa3, 0xac, 0x3e, 0x98, 0xd6, 0x67, 0x89, 0x96, 0xe9, 0x6b, 0x75,
	0xfa, 0x03, 0x22, 0x4f, 0xd1, 0x35, 0x70, 0xc7, 0x6c, 0xcc, 0xe3, 0x69, 0xc3, 0x55, 0xb8, 0xf3,
	0xf3, 0xb8, 0xfb, 0xea, 0x8c, 0x18, 0x0c, 0xba, 0x62, 0x36, 0x4b, 0x45, 0x61, 0xd1, 0x3c, 0x76,
	0x9f, 0x7b, 0x4c, 0x6f, 0x1b, 0xfc, 0x87, 0x05, 0x25, 0x29, 0x2e, 0x7d, 0x81, 0x5d, 0x82, 0xba,
	0x1f, 0x26, 0x2c, 0x0e, 0x69, 0x70, 0x48, 0x3d, 0x2f, 0x36, 0xb1, 0xae, 0x66, 0xca, 0x3b, 0x9e,
	0x17, 0x4b, 0x10, 0x7b, 0x52, 0x04, 0xe9, 0xb8, 0x57, 0x33, 0xa5, 0x01, 0xa9, 0x08, 0x4b, 0x67,
	0x8c, 0xb0, 0x7c, 0xb6, 0x08, 0xfb, 0xdc, 0x13, 0x26, 0x1b, 0x0b, 0x11, 0xca, 0x13, 0xa2, 0xce,
	0xf1, 0x1e, 0x38, 0x9d, 0xfe, 0x00, 0x9d, 0x87, 0x72, 0xe0, 0x8f, 0x7d, 0xfd, 0xb9, 0x6a, 0x11,
	0x2d, 0xc8, 0xcf, 0xc4, 0x58, 0xaf, 0x04, 0x15, 0x9b, 0x45, 0x32, 0x51, 0xe2, 0xd5, 0xab, 0x48,
	0x85, 0x23, 0x77, 0x94, 0x14, 0xf0, 0x3e, 0xb8, 0xda, 0x8d, 0x7f, 0x89, 0xef, 0x96, 0x0e, 0x02,
	0xed, 0x40, 0x29, 0xe2, 0x5e, 0x36, 0x6b, 0xeb, 0x2f, 0x04, 0x53, 0xd8, 0x51, 0x0a, 0x88, 0x6f,
	0x43, 0x59, 0x96, 0x4d, 0xa0, 0x5d, 0x28, 0xcb, 0x42, 0x66, 0xa6, 0x4b, 0x2a, 0x5d, 0xdc, 0xb4,
	0x0a, 0x8a, 0xff, 0x07, 0x4e, 0xaf, 0x2f, 0xd4, 0xe7, 0x5d, 0x64, 0xde, 0x9c, 0xb6, 0x1f, 0x35,
	0x2b, 0x50, 0xbe, 0x3b, 0x8e, 0x92, 0x69, 0xfb, 0xf3, 0x5f, 0x9f, 0x6d, 0x5a, 0xbf, 0x3d, 0xdb,
	0xb4, 0xfe, 0x7c, 0xb6, 0x69, 0xc1, 0x06, 0x8f, 0x47, 0xad, 0x89, 0x47, 0xa9, 0x68, 0x4d, 0x68,
	0xe0, 0x65, 0x17, 0xb4, 0x57, 0x0e, 0x68, 0xe0, 0xf5, 0xb5, 0xd0, 0xb7, 0x3e, 0x7d, 0x73, 0xe4,
	0x27, 0xc7, 0xe9, 0x51, 0x6b, 0xc8, 0xc7, 0x3b, 0x0a, 0x2d, 0xff, 0xf7, 0x78, 0x3b, 0x34, 0xf2,
	0xc5, 0xce, 0x28, 0x8e, 0x86, 0x3b, 0xc6, 0xee, 0xc8, 0x55, 0x7f, 0x83, 0x6e, 0xfc, 0x1d, 0x00,
	0x00, 0xff, 0xff, 0x2b, 0xc7, 0x7e, 0x5a, 0x39, 0x0d, 0x00, 0x00,
}

func (m *Search) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Search) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Search_Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Search_Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search_Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Config != nil {
		{
			size, err := m.Config.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Vector) > 0 {
		for iNdEx := len(m.Vector) - 1; iNdEx >= 0; iNdEx-- {
			f2 := math.Float32bits(float32(m.Vector[iNdEx]))
			i -= 4
			encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(f2))
		}
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Vector)*4))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Search_IDRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Search_IDRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search_IDRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Config != nil {
		{
			size, err := m.Config.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Search_ObjectRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Search_ObjectRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search_ObjectRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Config != nil {
		{
			size, err := m.Config.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Object) > 0 {
		i -= len(m.Object)
		copy(dAtA[i:], m.Object)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Object)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Search_Config) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Search_Config) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search_Config) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Timeout != 0 {
		i = encodeVarintPayload(dAtA, i, uint64(m.Timeout))
		i--
		dAtA[i] = 0x20
	}
	if m.Epsilon != 0 {
		i -= 4
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.Epsilon))))
		i--
		dAtA[i] = 0x1d
	}
	if m.Radius != 0 {
		i -= 4
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.Radius))))
		i--
		dAtA[i] = 0x15
	}
	if m.Num != 0 {
		i = encodeVarintPayload(dAtA, i, uint64(m.Num))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Search_Response) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Search_Response) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search_Response) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Results) > 0 {
		for iNdEx := len(m.Results) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Results[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPayload(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Meta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Meta) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Meta) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Meta_Key) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Meta_Key) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Meta_Key) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Meta_Keys) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Meta_Keys) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Meta_Keys) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Keys) > 0 {
		for iNdEx := len(m.Keys) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Keys[iNdEx])
			copy(dAtA[i:], m.Keys[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Keys[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Meta_Val) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Meta_Val) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Meta_Val) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Val) > 0 {
		i -= len(m.Val)
		copy(dAtA[i:], m.Val)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Val)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Meta_Vals) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Meta_Vals) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Meta_Vals) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Vals) > 0 {
		for iNdEx := len(m.Vals) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Vals[iNdEx])
			copy(dAtA[i:], m.Vals[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Vals[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Meta_KeyVal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Meta_KeyVal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Meta_KeyVal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Val) > 0 {
		i -= len(m.Val)
		copy(dAtA[i:], m.Val)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Val)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Key) > 0 {
		i -= len(m.Key)
		copy(dAtA[i:], m.Key)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Key)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Meta_KeyVals) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Meta_KeyVals) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Meta_KeyVals) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Kvs) > 0 {
		for iNdEx := len(m.Kvs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Kvs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPayload(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Object) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Object_Distance) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_Distance) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_Distance) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Distance != 0 {
		i -= 4
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.Distance))))
		i--
		dAtA[i] = 0x15
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Object_ID) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_ID) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_ID) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Object_IDs) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_IDs) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_IDs) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Ids) > 0 {
		for iNdEx := len(m.Ids) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Ids[iNdEx])
			copy(dAtA[i:], m.Ids[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Ids[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Object_Vector) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_Vector) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_Vector) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Vector) > 0 {
		for iNdEx := len(m.Vector) - 1; iNdEx >= 0; iNdEx-- {
			f5 := math.Float32bits(float32(m.Vector[iNdEx]))
			i -= 4
			encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(f5))
		}
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Vector)*4))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Object_Vectors) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_Vectors) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_Vectors) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Vectors) > 0 {
		for iNdEx := len(m.Vectors) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Vectors[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPayload(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Object_Blob) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_Blob) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_Blob) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Object) > 0 {
		i -= len(m.Object)
		copy(dAtA[i:], m.Object)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Object)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Object_Location) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_Location) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_Location) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Ips) > 0 {
		for iNdEx := len(m.Ips) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Ips[iNdEx])
			copy(dAtA[i:], m.Ips[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Ips[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Uuid) > 0 {
		i -= len(m.Uuid)
		copy(dAtA[i:], m.Uuid)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Uuid)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Object_Locations) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_Locations) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_Locations) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Locations) > 0 {
		for iNdEx := len(m.Locations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Locations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPayload(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Control) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Control) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Control) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Control_CreateIndexRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Control_CreateIndexRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Control_CreateIndexRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.PoolSize != 0 {
		i = encodeVarintPayload(dAtA, i, uint64(m.PoolSize))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Replication) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Replication) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Replication) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Replication_Recovery) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Replication_Recovery) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Replication_Recovery) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.DeletedAgents) > 0 {
		for iNdEx := len(m.DeletedAgents) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.DeletedAgents[iNdEx])
			copy(dAtA[i:], m.DeletedAgents[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.DeletedAgents[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Replication_Rebalance) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Replication_Rebalance) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Replication_Rebalance) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.LowUsageAgents) > 0 {
		for iNdEx := len(m.LowUsageAgents) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.LowUsageAgents[iNdEx])
			copy(dAtA[i:], m.LowUsageAgents[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.LowUsageAgents[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.HighUsageAgents) > 0 {
		for iNdEx := len(m.HighUsageAgents) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.HighUsageAgents[iNdEx])
			copy(dAtA[i:], m.HighUsageAgents[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.HighUsageAgents[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Replication_Agents) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Replication_Agents) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Replication_Agents) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ReplicatingAgent) > 0 {
		for iNdEx := len(m.ReplicatingAgent) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ReplicatingAgent[iNdEx])
			copy(dAtA[i:], m.ReplicatingAgent[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.ReplicatingAgent[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.RemovedAgents) > 0 {
		for iNdEx := len(m.RemovedAgents) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.RemovedAgents[iNdEx])
			copy(dAtA[i:], m.RemovedAgents[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.RemovedAgents[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Agents) > 0 {
		for iNdEx := len(m.Agents) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Agents[iNdEx])
			copy(dAtA[i:], m.Agents[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Agents[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Discoverer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Discoverer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Discoverer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Discoverer_Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Discoverer_Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Discoverer_Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Node) > 0 {
		i -= len(m.Node)
		copy(dAtA[i:], m.Node)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Node)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Namespace) > 0 {
		i -= len(m.Namespace)
		copy(dAtA[i:], m.Namespace)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Namespace)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Backup) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Backup_GetVector) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_GetVector) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_GetVector) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Backup_GetVector_Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_GetVector_Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_GetVector_Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Uuid) > 0 {
		i -= len(m.Uuid)
		copy(dAtA[i:], m.Uuid)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Uuid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Backup_GetVector_Owner) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_GetVector_Owner) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_GetVector_Owner) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Ip) > 0 {
		i -= len(m.Ip)
		copy(dAtA[i:], m.Ip)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Ip)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Backup_Locations) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Locations) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Locations) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Backup_Locations_Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Locations_Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Locations_Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Uuid) > 0 {
		i -= len(m.Uuid)
		copy(dAtA[i:], m.Uuid)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Uuid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Backup_Remove) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Remove) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Remove) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Backup_Remove_Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Remove_Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Remove_Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Uuid) > 0 {
		i -= len(m.Uuid)
		copy(dAtA[i:], m.Uuid)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Uuid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Backup_Remove_RequestMulti) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Remove_RequestMulti) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Remove_RequestMulti) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Uuids) > 0 {
		for iNdEx := len(m.Uuids) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Uuids[iNdEx])
			copy(dAtA[i:], m.Uuids[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Uuids[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Backup_IP) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_IP) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_IP) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Backup_IP_Register) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_IP_Register) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_IP_Register) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Backup_IP_Register_Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_IP_Register_Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_IP_Register_Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Ips) > 0 {
		for iNdEx := len(m.Ips) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Ips[iNdEx])
			copy(dAtA[i:], m.Ips[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Ips[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Uuid) > 0 {
		i -= len(m.Uuid)
		copy(dAtA[i:], m.Uuid)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Uuid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Backup_IP_Remove) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_IP_Remove) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_IP_Remove) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Backup_IP_Remove_Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_IP_Remove_Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_IP_Remove_Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Ips) > 0 {
		for iNdEx := len(m.Ips) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Ips[iNdEx])
			copy(dAtA[i:], m.Ips[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Ips[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Backup_MetaVector) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_MetaVector) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_MetaVector) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Ips) > 0 {
		for iNdEx := len(m.Ips) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Ips[iNdEx])
			copy(dAtA[i:], m.Ips[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Ips[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Vector) > 0 {
		for iNdEx := len(m.Vector) - 1; iNdEx >= 0; iNdEx-- {
			f6 := math.Float32bits(float32(m.Vector[iNdEx]))
			i -= 4
			encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(f6))
		}
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Vector)*4))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Uuid) > 0 {
		i -= len(m.Uuid)
		copy(dAtA[i:], m.Uuid)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Uuid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Backup_MetaVectors) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_MetaVectors) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_MetaVectors) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Vectors) > 0 {
		for iNdEx := len(m.Vectors) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Vectors[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPayload(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Backup_Compressed) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Compressed) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Compressed) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Backup_Compressed_MetaVector) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Compressed_MetaVector) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Compressed_MetaVector) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Ips) > 0 {
		for iNdEx := len(m.Ips) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Ips[iNdEx])
			copy(dAtA[i:], m.Ips[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Ips[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Vector) > 0 {
		i -= len(m.Vector)
		copy(dAtA[i:], m.Vector)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Vector)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Uuid) > 0 {
		i -= len(m.Uuid)
		copy(dAtA[i:], m.Uuid)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Uuid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Backup_Compressed_MetaVectors) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Compressed_MetaVectors) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Compressed_MetaVectors) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Vectors) > 0 {
		for iNdEx := len(m.Vectors) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Vectors[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPayload(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Info) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Info_Index) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_Index) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_Index) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Info_Index_Count) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_Index_Count) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_Index_Count) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Indexing {
		i--
		if m.Indexing {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.Uncommitted != 0 {
		i = encodeVarintPayload(dAtA, i, uint64(m.Uncommitted))
		i--
		dAtA[i] = 0x10
	}
	if m.Stored != 0 {
		i = encodeVarintPayload(dAtA, i, uint64(m.Stored))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Info_Index_UUID) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_Index_UUID) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_Index_UUID) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func (m *Info_Index_UUID_Committed) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_Index_UUID_Committed) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_Index_UUID_Committed) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Uuid) > 0 {
		i -= len(m.Uuid)
		copy(dAtA[i:], m.Uuid)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Uuid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Info_Index_UUID_Uncommitted) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_Index_UUID_Uncommitted) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_Index_UUID_Uncommitted) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Uuid) > 0 {
		i -= len(m.Uuid)
		copy(dAtA[i:], m.Uuid)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Uuid)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Info_Pod) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_Pod) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_Pod) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Node != nil {
		{
			size, err := m.Node.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x3a
	}
	if m.Memory != nil {
		{
			size, err := m.Memory.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.Cpu != nil {
		{
			size, err := m.Cpu.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Ip) > 0 {
		i -= len(m.Ip)
		copy(dAtA[i:], m.Ip)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Ip)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Namespace) > 0 {
		i -= len(m.Namespace)
		copy(dAtA[i:], m.Namespace)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Namespace)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.AppName) > 0 {
		i -= len(m.AppName)
		copy(dAtA[i:], m.AppName)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.AppName)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Info_Node) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_Node) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_Node) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Pods != nil {
		{
			size, err := m.Pods.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.Memory != nil {
		{
			size, err := m.Memory.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Cpu != nil {
		{
			size, err := m.Cpu.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.ExternalAddr) > 0 {
		i -= len(m.ExternalAddr)
		copy(dAtA[i:], m.ExternalAddr)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.ExternalAddr)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.InternalAddr) > 0 {
		i -= len(m.InternalAddr)
		copy(dAtA[i:], m.InternalAddr)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.InternalAddr)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Info_CPU) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_CPU) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_CPU) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Usage != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Usage))))
		i--
		dAtA[i] = 0x19
	}
	if m.Request != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Request))))
		i--
		dAtA[i] = 0x11
	}
	if m.Limit != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Limit))))
		i--
		dAtA[i] = 0x9
	}
	return len(dAtA) - i, nil
}

func (m *Info_Memory) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_Memory) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_Memory) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Usage != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Usage))))
		i--
		dAtA[i] = 0x19
	}
	if m.Request != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Request))))
		i--
		dAtA[i] = 0x11
	}
	if m.Limit != 0 {
		i -= 8
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Limit))))
		i--
		dAtA[i] = 0x9
	}
	return len(dAtA) - i, nil
}

func (m *Info_Pods) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_Pods) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_Pods) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Pods) > 0 {
		for iNdEx := len(m.Pods) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Pods[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPayload(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Info_Nodes) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_Nodes) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_Nodes) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Nodes) > 0 {
		for iNdEx := len(m.Nodes) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Nodes[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPayload(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Info_IPs) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Info_IPs) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Info_IPs) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Ip) > 0 {
		for iNdEx := len(m.Ip) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Ip[iNdEx])
			copy(dAtA[i:], m.Ip[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Ip[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Empty) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Empty) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Empty) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func encodeVarintPayload(dAtA []byte, offset int, v uint64) int {
	offset -= sovPayload(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Search) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Search_Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Vector) > 0 {
		n += 1 + sovPayload(uint64(len(m.Vector)*4)) + len(m.Vector)*4
	}
	if m.Config != nil {
		l = m.Config.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Search_IDRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.Config != nil {
		l = m.Config.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Search_ObjectRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Object)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.Config != nil {
		l = m.Config.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Search_Config) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Num != 0 {
		n += 1 + sovPayload(uint64(m.Num))
	}
	if m.Radius != 0 {
		n += 5
	}
	if m.Epsilon != 0 {
		n += 5
	}
	if m.Timeout != 0 {
		n += 1 + sovPayload(uint64(m.Timeout))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Search_Response) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Results) > 0 {
		for _, e := range m.Results {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Meta) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Meta_Key) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Meta_Keys) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Keys) > 0 {
		for _, s := range m.Keys {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Meta_Val) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Val)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Meta_Vals) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Vals) > 0 {
		for _, s := range m.Vals {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Meta_KeyVal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	l = len(m.Val)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Meta_KeyVals) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Kvs) > 0 {
		for _, e := range m.Kvs {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_Distance) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.Distance != 0 {
		n += 5
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_ID) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_IDs) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Ids) > 0 {
		for _, s := range m.Ids {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_Vector) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if len(m.Vector) > 0 {
		n += 1 + sovPayload(uint64(len(m.Vector)*4)) + len(m.Vector)*4
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_Vectors) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Vectors) > 0 {
		for _, e := range m.Vectors {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_Blob) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	l = len(m.Object)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_Location) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	l = len(m.Uuid)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if len(m.Ips) > 0 {
		for _, s := range m.Ips {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_Locations) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Locations) > 0 {
		for _, e := range m.Locations {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Control) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Control_CreateIndexRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolSize != 0 {
		n += 1 + sovPayload(uint64(m.PoolSize))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Replication) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Replication_Recovery) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.DeletedAgents) > 0 {
		for _, s := range m.DeletedAgents {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Replication_Rebalance) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.HighUsageAgents) > 0 {
		for _, s := range m.HighUsageAgents {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if len(m.LowUsageAgents) > 0 {
		for _, s := range m.LowUsageAgents {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Replication_Agents) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Agents) > 0 {
		for _, s := range m.Agents {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if len(m.RemovedAgents) > 0 {
		for _, s := range m.RemovedAgents {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if len(m.ReplicatingAgent) > 0 {
		for _, s := range m.ReplicatingAgent {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Discoverer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Discoverer_Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	l = len(m.Namespace)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	l = len(m.Node)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_GetVector) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_GetVector_Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Uuid)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_GetVector_Owner) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Ip)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_Locations) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_Locations_Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Uuid)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_Remove) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_Remove_Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Uuid)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_Remove_RequestMulti) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Uuids) > 0 {
		for _, s := range m.Uuids {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_IP) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_IP_Register) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_IP_Register_Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Uuid)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if len(m.Ips) > 0 {
		for _, s := range m.Ips {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_IP_Remove) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_IP_Remove_Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Ips) > 0 {
		for _, s := range m.Ips {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_MetaVector) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Uuid)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if len(m.Vector) > 0 {
		n += 1 + sovPayload(uint64(len(m.Vector)*4)) + len(m.Vector)*4
	}
	if len(m.Ips) > 0 {
		for _, s := range m.Ips {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_MetaVectors) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Vectors) > 0 {
		for _, e := range m.Vectors {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_Compressed) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_Compressed_MetaVector) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Uuid)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	l = len(m.Vector)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if len(m.Ips) > 0 {
		for _, s := range m.Ips {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Backup_Compressed_MetaVectors) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Vectors) > 0 {
		for _, e := range m.Vectors {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_Index) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_Index_Count) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Stored != 0 {
		n += 1 + sovPayload(uint64(m.Stored))
	}
	if m.Uncommitted != 0 {
		n += 1 + sovPayload(uint64(m.Uncommitted))
	}
	if m.Indexing {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_Index_UUID) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_Index_UUID_Committed) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Uuid)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_Index_UUID_Uncommitted) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Uuid)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_Pod) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.AppName)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	l = len(m.Namespace)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	l = len(m.Ip)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.Cpu != nil {
		l = m.Cpu.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.Memory != nil {
		l = m.Memory.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.Node != nil {
		l = m.Node.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_Node) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	l = len(m.InternalAddr)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	l = len(m.ExternalAddr)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.Cpu != nil {
		l = m.Cpu.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.Memory != nil {
		l = m.Memory.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.Pods != nil {
		l = m.Pods.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_CPU) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Limit != 0 {
		n += 9
	}
	if m.Request != 0 {
		n += 9
	}
	if m.Usage != 0 {
		n += 9
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_Memory) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Limit != 0 {
		n += 9
	}
	if m.Request != 0 {
		n += 9
	}
	if m.Usage != 0 {
		n += 9
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_Pods) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Pods) > 0 {
		for _, e := range m.Pods {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_Nodes) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Nodes) > 0 {
		for _, e := range m.Nodes {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Info_IPs) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Ip) > 0 {
		for _, s := range m.Ip {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Empty) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovPayload(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPayload(x uint64) (n int) {
	return sovPayload(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Search) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Search: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Search: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Search_Request) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Request: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Request: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 5 {
				var v uint32
				if (iNdEx + 4) > l {
					return io.ErrUnexpectedEOF
				}
				v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
				iNdEx += 4
				v2 := float32(math.Float32frombits(v))
				m.Vector = append(m.Vector, v2)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPayload
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthPayload
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthPayload
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen / 4
				if elementCount != 0 && len(m.Vector) == 0 {
					m.Vector = make([]float32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint32
					if (iNdEx + 4) > l {
						return io.ErrUnexpectedEOF
					}
					v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
					iNdEx += 4
					v2 := float32(math.Float32frombits(v))
					m.Vector = append(m.Vector, v2)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Vector", wireType)
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Config == nil {
				m.Config = &Search_Config{}
			}
			if err := m.Config.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Search_IDRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IDRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IDRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Config == nil {
				m.Config = &Search_Config{}
			}
			if err := m.Config.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Search_ObjectRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ObjectRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ObjectRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Object", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Object = append(m.Object[:0], dAtA[iNdEx:postIndex]...)
			if m.Object == nil {
				m.Object = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Config == nil {
				m.Config = &Search_Config{}
			}
			if err := m.Config.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Search_Config) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Config: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Config: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Num", wireType)
			}
			m.Num = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Num |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field Radius", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
			iNdEx += 4
			m.Radius = float32(math.Float32frombits(v))
		case 3:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field Epsilon", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
			iNdEx += 4
			m.Epsilon = float32(math.Float32frombits(v))
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timeout", wireType)
			}
			m.Timeout = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timeout |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Search_Response) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Response: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Response: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Results", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Results = append(m.Results, &Object_Distance{})
			if err := m.Results[len(m.Results)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Meta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Meta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Meta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Meta_Key) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Key: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Key: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Meta_Keys) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Keys: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Keys: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Keys", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Keys = append(m.Keys, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Meta_Val) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Val: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Val: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Val", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Val = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Meta_Vals) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Vals: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Vals: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vals", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vals = append(m.Vals, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Meta_KeyVal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: KeyVal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: KeyVal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Val", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Val = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Meta_KeyVals) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: KeyVals: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: KeyVals: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Kvs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Kvs = append(m.Kvs, &Meta_KeyVal{})
			if err := m.Kvs[len(m.Kvs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Object) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Object: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Object: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Object_Distance) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Distance: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Distance: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field Distance", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
			iNdEx += 4
			m.Distance = float32(math.Float32frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Object_ID) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ID: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ID: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Object_IDs) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IDs: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IDs: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ids", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ids = append(m.Ids, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Object_Vector) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Vector: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Vector: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType == 5 {
				var v uint32
				if (iNdEx + 4) > l {
					return io.ErrUnexpectedEOF
				}
				v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
				iNdEx += 4
				v2 := float32(math.Float32frombits(v))
				m.Vector = append(m.Vector, v2)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPayload
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthPayload
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthPayload
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen / 4
				if elementCount != 0 && len(m.Vector) == 0 {
					m.Vector = make([]float32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint32
					if (iNdEx + 4) > l {
						return io.ErrUnexpectedEOF
					}
					v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
					iNdEx += 4
					v2 := float32(math.Float32frombits(v))
					m.Vector = append(m.Vector, v2)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Vector", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Object_Vectors) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Vectors: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Vectors: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vectors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vectors = append(m.Vectors, &Object_Vector{})
			if err := m.Vectors[len(m.Vectors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Object_Blob) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Blob: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Blob: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Object", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Object = append(m.Object[:0], dAtA[iNdEx:postIndex]...)
			if m.Object == nil {
				m.Object = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Object_Location) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Location: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Location: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uuid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uuid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ips", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ips = append(m.Ips, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Object_Locations) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Locations: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Locations: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Locations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Locations = append(m.Locations, &Object_Location{})
			if err := m.Locations[len(m.Locations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Control) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Control: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Control: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Control_CreateIndexRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CreateIndexRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CreateIndexRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolSize", wireType)
			}
			m.PoolSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolSize |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Replication) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Replication: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Replication: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Replication_Recovery) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Recovery: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Recovery: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeletedAgents", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DeletedAgents = append(m.DeletedAgents, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Replication_Rebalance) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Rebalance: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Rebalance: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HighUsageAgents", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HighUsageAgents = append(m.HighUsageAgents, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LowUsageAgents", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LowUsageAgents = append(m.LowUsageAgents, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Replication_Agents) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Agents: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Agents: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Agents", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Agents = append(m.Agents, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RemovedAgents", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RemovedAgents = append(m.RemovedAgents, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReplicatingAgent", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReplicatingAgent = append(m.ReplicatingAgent, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Discoverer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Discoverer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Discoverer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Discoverer_Request) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Request: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Request: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Namespace", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Namespace = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Node", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Node = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Backup: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Backup: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_GetVector) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GetVector: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GetVector: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_GetVector_Request) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Request: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Request: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uuid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uuid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_GetVector_Owner) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Owner: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Owner: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ip", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ip = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_Locations) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Locations: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Locations: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_Locations_Request) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Request: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Request: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uuid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uuid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_Remove) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Remove: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Remove: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_Remove_Request) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Request: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Request: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uuid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uuid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_Remove_RequestMulti) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RequestMulti: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RequestMulti: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uuids", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uuids = append(m.Uuids, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_IP) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IP: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IP: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_IP_Register) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Register: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Register: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_IP_Register_Request) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Request: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Request: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uuid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uuid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ips", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ips = append(m.Ips, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_IP_Remove) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Remove: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Remove: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_IP_Remove_Request) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Request: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Request: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ips", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ips = append(m.Ips, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_MetaVector) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetaVector: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetaVector: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uuid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uuid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType == 5 {
				var v uint32
				if (iNdEx + 4) > l {
					return io.ErrUnexpectedEOF
				}
				v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
				iNdEx += 4
				v2 := float32(math.Float32frombits(v))
				m.Vector = append(m.Vector, v2)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowPayload
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthPayload
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthPayload
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				elementCount = packedLen / 4
				if elementCount != 0 && len(m.Vector) == 0 {
					m.Vector = make([]float32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint32
					if (iNdEx + 4) > l {
						return io.ErrUnexpectedEOF
					}
					v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
					iNdEx += 4
					v2 := float32(math.Float32frombits(v))
					m.Vector = append(m.Vector, v2)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Vector", wireType)
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ips", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ips = append(m.Ips, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_MetaVectors) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetaVectors: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetaVectors: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vectors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vectors = append(m.Vectors, &Backup_MetaVector{})
			if err := m.Vectors[len(m.Vectors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_Compressed) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Compressed: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Compressed: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_Compressed_MetaVector) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetaVector: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetaVector: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uuid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uuid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vector", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vector = append(m.Vector[:0], dAtA[iNdEx:postIndex]...)
			if m.Vector == nil {
				m.Vector = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ips", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ips = append(m.Ips, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Backup_Compressed_MetaVectors) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetaVectors: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetaVectors: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vectors", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vectors = append(m.Vectors, &Backup_Compressed_MetaVector{})
			if err := m.Vectors[len(m.Vectors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Info: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Info: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_Index) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Index: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Index: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_Index_Count) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Count: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Count: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stored", wireType)
			}
			m.Stored = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Stored |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uncommitted", wireType)
			}
			m.Uncommitted = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Uncommitted |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Indexing", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Indexing = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_Index_UUID) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: UUID: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UUID: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_Index_UUID_Committed) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Committed: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Committed: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uuid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uuid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_Index_UUID_Uncommitted) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Uncommitted: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Uncommitted: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uuid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uuid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_Pod) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Pod: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Pod: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AppName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AppName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Namespace", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Namespace = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ip", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ip = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cpu", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Cpu == nil {
				m.Cpu = &Info_CPU{}
			}
			if err := m.Cpu.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Memory", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Memory == nil {
				m.Memory = &Info_Memory{}
			}
			if err := m.Memory.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Node", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Node == nil {
				m.Node = &Info_Node{}
			}
			if err := m.Node.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_Node) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Node: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Node: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InternalAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InternalAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExternalAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExternalAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cpu", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Cpu == nil {
				m.Cpu = &Info_CPU{}
			}
			if err := m.Cpu.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Memory", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Memory == nil {
				m.Memory = &Info_Memory{}
			}
			if err := m.Memory.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pods", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pods == nil {
				m.Pods = &Info_Pods{}
			}
			if err := m.Pods.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_CPU) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: CPU: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CPU: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Limit", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Limit = float64(math.Float64frombits(v))
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Request", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Request = float64(math.Float64frombits(v))
		case 3:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Usage", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Usage = float64(math.Float64frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_Memory) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Memory: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Memory: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Limit", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Limit = float64(math.Float64frombits(v))
		case 2:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Request", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Request = float64(math.Float64frombits(v))
		case 3:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Usage", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Usage = float64(math.Float64frombits(v))
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_Pods) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Pods: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Pods: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pods", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pods = append(m.Pods, &Info_Pod{})
			if err := m.Pods[len(m.Pods)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_Nodes) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Nodes: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Nodes: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nodes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Nodes = append(m.Nodes, &Info_Node{})
			if err := m.Nodes[len(m.Nodes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Info_IPs) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IPs: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IPs: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ip", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPayload
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPayload
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ip = append(m.Ip, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Empty) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Empty: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Empty: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPayload
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipPayload(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPayload
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthPayload
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPayload
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPayload
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPayload        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPayload          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPayload = fmt.Errorf("proto: unexpected end of group")
)
