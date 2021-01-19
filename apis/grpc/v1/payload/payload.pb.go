//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	rpc "github.com/gogo/googleapis/google/rpc"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	golang_proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = golang_proto.Marshal
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
	return fileDescriptor_f0518b37b4e7594b, []int{0}
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

func (*Search) XXX_MessageName() string {
	return "payload.v1.Search"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{0, 0}
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

func (*Search_Request) XXX_MessageName() string {
	return "payload.v1.Search.Request"
}

type Search_MultiRequest struct {
	Requests             []*Search_Request `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Search_MultiRequest) Reset()         { *m = Search_MultiRequest{} }
func (m *Search_MultiRequest) String() string { return proto.CompactTextString(m) }
func (*Search_MultiRequest) ProtoMessage()    {}
func (*Search_MultiRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{0, 1}
}
func (m *Search_MultiRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Search_MultiRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Search_MultiRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Search_MultiRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Search_MultiRequest.Merge(m, src)
}
func (m *Search_MultiRequest) XXX_Size() int {
	return m.Size()
}
func (m *Search_MultiRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_Search_MultiRequest.DiscardUnknown(m)
}

var xxx_messageInfo_Search_MultiRequest proto.InternalMessageInfo

func (m *Search_MultiRequest) GetRequests() []*Search_Request {
	if m != nil {
		return m.Requests
	}
	return nil
}

func (*Search_MultiRequest) XXX_MessageName() string {
	return "payload.v1.Search.MultiRequest"
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
	return fileDescriptor_f0518b37b4e7594b, []int{0, 2}
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

func (*Search_IDRequest) XXX_MessageName() string {
	return "payload.v1.Search.IDRequest"
}

type Search_MultiIDRequest struct {
	Requests             []*Search_IDRequest `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Search_MultiIDRequest) Reset()         { *m = Search_MultiIDRequest{} }
func (m *Search_MultiIDRequest) String() string { return proto.CompactTextString(m) }
func (*Search_MultiIDRequest) ProtoMessage()    {}
func (*Search_MultiIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{0, 3}
}
func (m *Search_MultiIDRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Search_MultiIDRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Search_MultiIDRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Search_MultiIDRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Search_MultiIDRequest.Merge(m, src)
}
func (m *Search_MultiIDRequest) XXX_Size() int {
	return m.Size()
}
func (m *Search_MultiIDRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_Search_MultiIDRequest.DiscardUnknown(m)
}

var xxx_messageInfo_Search_MultiIDRequest proto.InternalMessageInfo

func (m *Search_MultiIDRequest) GetRequests() []*Search_IDRequest {
	if m != nil {
		return m.Requests
	}
	return nil
}

func (*Search_MultiIDRequest) XXX_MessageName() string {
	return "payload.v1.Search.MultiIDRequest"
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
	return fileDescriptor_f0518b37b4e7594b, []int{0, 4}
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

func (*Search_ObjectRequest) XXX_MessageName() string {
	return "payload.v1.Search.ObjectRequest"
}

type Search_Config struct {
	RequestId            string         `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Num                  uint32         `protobuf:"varint,2,opt,name=num,proto3" json:"num,omitempty"`
	Radius               float32        `protobuf:"fixed32,3,opt,name=radius,proto3" json:"radius,omitempty"`
	Epsilon              float32        `protobuf:"fixed32,4,opt,name=epsilon,proto3" json:"epsilon,omitempty"`
	Timeout              int64          `protobuf:"varint,5,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Filters              *Filter_Config `protobuf:"bytes,6,opt,name=filters,proto3" json:"filters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Search_Config) Reset()         { *m = Search_Config{} }
func (m *Search_Config) String() string { return proto.CompactTextString(m) }
func (*Search_Config) ProtoMessage()    {}
func (*Search_Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{0, 5}
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

func (m *Search_Config) GetRequestId() string {
	if m != nil {
		return m.RequestId
	}
	return ""
}

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

func (m *Search_Config) GetFilters() *Filter_Config {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (*Search_Config) XXX_MessageName() string {
	return "payload.v1.Search.Config"
}

type Search_Response struct {
	RequestId            string             `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Results              []*Object_Distance `protobuf:"bytes,2,rep,name=results,proto3" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Search_Response) Reset()         { *m = Search_Response{} }
func (m *Search_Response) String() string { return proto.CompactTextString(m) }
func (*Search_Response) ProtoMessage()    {}
func (*Search_Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{0, 6}
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

func (m *Search_Response) GetRequestId() string {
	if m != nil {
		return m.RequestId
	}
	return ""
}

func (m *Search_Response) GetResults() []*Object_Distance {
	if m != nil {
		return m.Results
	}
	return nil
}

func (*Search_Response) XXX_MessageName() string {
	return "payload.v1.Search.Response"
}

type Search_Responses struct {
	Responses            []*Search_Response `protobuf:"bytes,1,rep,name=responses,proto3" json:"responses,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Search_Responses) Reset()         { *m = Search_Responses{} }
func (m *Search_Responses) String() string { return proto.CompactTextString(m) }
func (*Search_Responses) ProtoMessage()    {}
func (*Search_Responses) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{0, 7}
}
func (m *Search_Responses) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Search_Responses) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Search_Responses.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Search_Responses) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Search_Responses.Merge(m, src)
}
func (m *Search_Responses) XXX_Size() int {
	return m.Size()
}
func (m *Search_Responses) XXX_DiscardUnknown() {
	xxx_messageInfo_Search_Responses.DiscardUnknown(m)
}

var xxx_messageInfo_Search_Responses proto.InternalMessageInfo

func (m *Search_Responses) GetResponses() []*Search_Response {
	if m != nil {
		return m.Responses
	}
	return nil
}

func (*Search_Responses) XXX_MessageName() string {
	return "payload.v1.Search.Responses"
}

type Search_StreamResponse struct {
	// Types that are valid to be assigned to Payload:
	//	*Search_StreamResponse_Response
	//	*Search_StreamResponse_Status
	Payload              isSearch_StreamResponse_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *Search_StreamResponse) Reset()         { *m = Search_StreamResponse{} }
func (m *Search_StreamResponse) String() string { return proto.CompactTextString(m) }
func (*Search_StreamResponse) ProtoMessage()    {}
func (*Search_StreamResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{0, 8}
}
func (m *Search_StreamResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Search_StreamResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Search_StreamResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Search_StreamResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Search_StreamResponse.Merge(m, src)
}
func (m *Search_StreamResponse) XXX_Size() int {
	return m.Size()
}
func (m *Search_StreamResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_Search_StreamResponse.DiscardUnknown(m)
}

var xxx_messageInfo_Search_StreamResponse proto.InternalMessageInfo

type isSearch_StreamResponse_Payload interface {
	isSearch_StreamResponse_Payload()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Search_StreamResponse_Response struct {
	Response *Search_Response `protobuf:"bytes,1,opt,name=response,proto3,oneof" json:"response,omitempty"`
}
type Search_StreamResponse_Status struct {
	Status *rpc.Status `protobuf:"bytes,2,opt,name=status,proto3,oneof" json:"status,omitempty"`
}

func (*Search_StreamResponse_Response) isSearch_StreamResponse_Payload() {}
func (*Search_StreamResponse_Status) isSearch_StreamResponse_Payload()   {}

func (m *Search_StreamResponse) GetPayload() isSearch_StreamResponse_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Search_StreamResponse) GetResponse() *Search_Response {
	if x, ok := m.GetPayload().(*Search_StreamResponse_Response); ok {
		return x.Response
	}
	return nil
}

func (m *Search_StreamResponse) GetStatus() *rpc.Status {
	if x, ok := m.GetPayload().(*Search_StreamResponse_Status); ok {
		return x.Status
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Search_StreamResponse) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Search_StreamResponse_Response)(nil),
		(*Search_StreamResponse_Status)(nil),
	}
}

func (*Search_StreamResponse) XXX_MessageName() string {
	return "payload.v1.Search.StreamResponse"
}

type Filter struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Filter) Reset()         { *m = Filter{} }
func (m *Filter) String() string { return proto.CompactTextString(m) }
func (*Filter) ProtoMessage()    {}
func (*Filter) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{1}
}
func (m *Filter) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Filter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Filter.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Filter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Filter.Merge(m, src)
}
func (m *Filter) XXX_Size() int {
	return m.Size()
}
func (m *Filter) XXX_DiscardUnknown() {
	xxx_messageInfo_Filter.DiscardUnknown(m)
}

var xxx_messageInfo_Filter proto.InternalMessageInfo

func (*Filter) XXX_MessageName() string {
	return "payload.v1.Filter"
}

type Filter_Target struct {
	Host                 string   `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port                 uint32   `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Filter_Target) Reset()         { *m = Filter_Target{} }
func (m *Filter_Target) String() string { return proto.CompactTextString(m) }
func (*Filter_Target) ProtoMessage()    {}
func (*Filter_Target) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{1, 0}
}
func (m *Filter_Target) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Filter_Target) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Filter_Target.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Filter_Target) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Filter_Target.Merge(m, src)
}
func (m *Filter_Target) XXX_Size() int {
	return m.Size()
}
func (m *Filter_Target) XXX_DiscardUnknown() {
	xxx_messageInfo_Filter_Target.DiscardUnknown(m)
}

var xxx_messageInfo_Filter_Target proto.InternalMessageInfo

func (m *Filter_Target) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *Filter_Target) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (*Filter_Target) XXX_MessageName() string {
	return "payload.v1.Filter.Target"
}

type Filter_Config struct {
	Targets              []string `protobuf:"bytes,1,rep,name=targets,proto3" json:"targets,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Filter_Config) Reset()         { *m = Filter_Config{} }
func (m *Filter_Config) String() string { return proto.CompactTextString(m) }
func (*Filter_Config) ProtoMessage()    {}
func (*Filter_Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{1, 1}
}
func (m *Filter_Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Filter_Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Filter_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Filter_Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Filter_Config.Merge(m, src)
}
func (m *Filter_Config) XXX_Size() int {
	return m.Size()
}
func (m *Filter_Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Filter_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Filter_Config proto.InternalMessageInfo

func (m *Filter_Config) GetTargets() []string {
	if m != nil {
		return m.Targets
	}
	return nil
}

func (*Filter_Config) XXX_MessageName() string {
	return "payload.v1.Filter.Config"
}

type Insert struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Insert) Reset()         { *m = Insert{} }
func (m *Insert) String() string { return proto.CompactTextString(m) }
func (*Insert) ProtoMessage()    {}
func (*Insert) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{2}
}
func (m *Insert) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Insert) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Insert.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Insert) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Insert.Merge(m, src)
}
func (m *Insert) XXX_Size() int {
	return m.Size()
}
func (m *Insert) XXX_DiscardUnknown() {
	xxx_messageInfo_Insert.DiscardUnknown(m)
}

var xxx_messageInfo_Insert proto.InternalMessageInfo

func (*Insert) XXX_MessageName() string {
	return "payload.v1.Insert"
}

type Insert_Request struct {
	Vector               *Object_Vector `protobuf:"bytes,1,opt,name=vector,proto3" json:"vector,omitempty"`
	Config               *Insert_Config `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Insert_Request) Reset()         { *m = Insert_Request{} }
func (m *Insert_Request) String() string { return proto.CompactTextString(m) }
func (*Insert_Request) ProtoMessage()    {}
func (*Insert_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{2, 0}
}
func (m *Insert_Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Insert_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Insert_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Insert_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Insert_Request.Merge(m, src)
}
func (m *Insert_Request) XXX_Size() int {
	return m.Size()
}
func (m *Insert_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Insert_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Insert_Request proto.InternalMessageInfo

func (m *Insert_Request) GetVector() *Object_Vector {
	if m != nil {
		return m.Vector
	}
	return nil
}

func (m *Insert_Request) GetConfig() *Insert_Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func (*Insert_Request) XXX_MessageName() string {
	return "payload.v1.Insert.Request"
}

type Insert_MultiRequest struct {
	Requests             []*Insert_Request `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Insert_MultiRequest) Reset()         { *m = Insert_MultiRequest{} }
func (m *Insert_MultiRequest) String() string { return proto.CompactTextString(m) }
func (*Insert_MultiRequest) ProtoMessage()    {}
func (*Insert_MultiRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{2, 1}
}
func (m *Insert_MultiRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Insert_MultiRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Insert_MultiRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Insert_MultiRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Insert_MultiRequest.Merge(m, src)
}
func (m *Insert_MultiRequest) XXX_Size() int {
	return m.Size()
}
func (m *Insert_MultiRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_Insert_MultiRequest.DiscardUnknown(m)
}

var xxx_messageInfo_Insert_MultiRequest proto.InternalMessageInfo

func (m *Insert_MultiRequest) GetRequests() []*Insert_Request {
	if m != nil {
		return m.Requests
	}
	return nil
}

func (*Insert_MultiRequest) XXX_MessageName() string {
	return "payload.v1.Insert.MultiRequest"
}

type Insert_Config struct {
	SkipStrictExistCheck bool           `protobuf:"varint,1,opt,name=skip_strict_exist_check,json=skipStrictExistCheck,proto3" json:"skip_strict_exist_check,omitempty"`
	Filters              *Filter_Config `protobuf:"bytes,2,opt,name=filters,proto3" json:"filters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Insert_Config) Reset()         { *m = Insert_Config{} }
func (m *Insert_Config) String() string { return proto.CompactTextString(m) }
func (*Insert_Config) ProtoMessage()    {}
func (*Insert_Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{2, 2}
}
func (m *Insert_Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Insert_Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Insert_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Insert_Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Insert_Config.Merge(m, src)
}
func (m *Insert_Config) XXX_Size() int {
	return m.Size()
}
func (m *Insert_Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Insert_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Insert_Config proto.InternalMessageInfo

func (m *Insert_Config) GetSkipStrictExistCheck() bool {
	if m != nil {
		return m.SkipStrictExistCheck
	}
	return false
}

func (m *Insert_Config) GetFilters() *Filter_Config {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (*Insert_Config) XXX_MessageName() string {
	return "payload.v1.Insert.Config"
}

type Update struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Update) Reset()         { *m = Update{} }
func (m *Update) String() string { return proto.CompactTextString(m) }
func (*Update) ProtoMessage()    {}
func (*Update) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{3}
}
func (m *Update) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Update) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Update.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Update) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Update.Merge(m, src)
}
func (m *Update) XXX_Size() int {
	return m.Size()
}
func (m *Update) XXX_DiscardUnknown() {
	xxx_messageInfo_Update.DiscardUnknown(m)
}

var xxx_messageInfo_Update proto.InternalMessageInfo

func (*Update) XXX_MessageName() string {
	return "payload.v1.Update"
}

type Update_Request struct {
	Vector               *Object_Vector `protobuf:"bytes,1,opt,name=vector,proto3" json:"vector,omitempty"`
	Config               *Update_Config `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Update_Request) Reset()         { *m = Update_Request{} }
func (m *Update_Request) String() string { return proto.CompactTextString(m) }
func (*Update_Request) ProtoMessage()    {}
func (*Update_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{3, 0}
}
func (m *Update_Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Update_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Update_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Update_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Update_Request.Merge(m, src)
}
func (m *Update_Request) XXX_Size() int {
	return m.Size()
}
func (m *Update_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Update_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Update_Request proto.InternalMessageInfo

func (m *Update_Request) GetVector() *Object_Vector {
	if m != nil {
		return m.Vector
	}
	return nil
}

func (m *Update_Request) GetConfig() *Update_Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func (*Update_Request) XXX_MessageName() string {
	return "payload.v1.Update.Request"
}

type Update_MultiRequest struct {
	Requests             []*Update_Request `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Update_MultiRequest) Reset()         { *m = Update_MultiRequest{} }
func (m *Update_MultiRequest) String() string { return proto.CompactTextString(m) }
func (*Update_MultiRequest) ProtoMessage()    {}
func (*Update_MultiRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{3, 1}
}
func (m *Update_MultiRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Update_MultiRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Update_MultiRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Update_MultiRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Update_MultiRequest.Merge(m, src)
}
func (m *Update_MultiRequest) XXX_Size() int {
	return m.Size()
}
func (m *Update_MultiRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_Update_MultiRequest.DiscardUnknown(m)
}

var xxx_messageInfo_Update_MultiRequest proto.InternalMessageInfo

func (m *Update_MultiRequest) GetRequests() []*Update_Request {
	if m != nil {
		return m.Requests
	}
	return nil
}

func (*Update_MultiRequest) XXX_MessageName() string {
	return "payload.v1.Update.MultiRequest"
}

type Update_Config struct {
	SkipStrictExistCheck bool           `protobuf:"varint,1,opt,name=skip_strict_exist_check,json=skipStrictExistCheck,proto3" json:"skip_strict_exist_check,omitempty"`
	Filters              *Filter_Config `protobuf:"bytes,2,opt,name=filters,proto3" json:"filters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Update_Config) Reset()         { *m = Update_Config{} }
func (m *Update_Config) String() string { return proto.CompactTextString(m) }
func (*Update_Config) ProtoMessage()    {}
func (*Update_Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{3, 2}
}
func (m *Update_Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Update_Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Update_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Update_Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Update_Config.Merge(m, src)
}
func (m *Update_Config) XXX_Size() int {
	return m.Size()
}
func (m *Update_Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Update_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Update_Config proto.InternalMessageInfo

func (m *Update_Config) GetSkipStrictExistCheck() bool {
	if m != nil {
		return m.SkipStrictExistCheck
	}
	return false
}

func (m *Update_Config) GetFilters() *Filter_Config {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (*Update_Config) XXX_MessageName() string {
	return "payload.v1.Update.Config"
}

type Upsert struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Upsert) Reset()         { *m = Upsert{} }
func (m *Upsert) String() string { return proto.CompactTextString(m) }
func (*Upsert) ProtoMessage()    {}
func (*Upsert) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{4}
}
func (m *Upsert) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Upsert) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Upsert.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Upsert) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Upsert.Merge(m, src)
}
func (m *Upsert) XXX_Size() int {
	return m.Size()
}
func (m *Upsert) XXX_DiscardUnknown() {
	xxx_messageInfo_Upsert.DiscardUnknown(m)
}

var xxx_messageInfo_Upsert proto.InternalMessageInfo

func (*Upsert) XXX_MessageName() string {
	return "payload.v1.Upsert"
}

type Upsert_Request struct {
	Vector               *Object_Vector `protobuf:"bytes,1,opt,name=vector,proto3" json:"vector,omitempty"`
	Config               *Upsert_Config `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Upsert_Request) Reset()         { *m = Upsert_Request{} }
func (m *Upsert_Request) String() string { return proto.CompactTextString(m) }
func (*Upsert_Request) ProtoMessage()    {}
func (*Upsert_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{4, 0}
}
func (m *Upsert_Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Upsert_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Upsert_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Upsert_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Upsert_Request.Merge(m, src)
}
func (m *Upsert_Request) XXX_Size() int {
	return m.Size()
}
func (m *Upsert_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Upsert_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Upsert_Request proto.InternalMessageInfo

func (m *Upsert_Request) GetVector() *Object_Vector {
	if m != nil {
		return m.Vector
	}
	return nil
}

func (m *Upsert_Request) GetConfig() *Upsert_Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func (*Upsert_Request) XXX_MessageName() string {
	return "payload.v1.Upsert.Request"
}

type Upsert_MultiRequest struct {
	Requests             []*Upsert_Request `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Upsert_MultiRequest) Reset()         { *m = Upsert_MultiRequest{} }
func (m *Upsert_MultiRequest) String() string { return proto.CompactTextString(m) }
func (*Upsert_MultiRequest) ProtoMessage()    {}
func (*Upsert_MultiRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{4, 1}
}
func (m *Upsert_MultiRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Upsert_MultiRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Upsert_MultiRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Upsert_MultiRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Upsert_MultiRequest.Merge(m, src)
}
func (m *Upsert_MultiRequest) XXX_Size() int {
	return m.Size()
}
func (m *Upsert_MultiRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_Upsert_MultiRequest.DiscardUnknown(m)
}

var xxx_messageInfo_Upsert_MultiRequest proto.InternalMessageInfo

func (m *Upsert_MultiRequest) GetRequests() []*Upsert_Request {
	if m != nil {
		return m.Requests
	}
	return nil
}

func (*Upsert_MultiRequest) XXX_MessageName() string {
	return "payload.v1.Upsert.MultiRequest"
}

type Upsert_Config struct {
	SkipStrictExistCheck bool           `protobuf:"varint,1,opt,name=skip_strict_exist_check,json=skipStrictExistCheck,proto3" json:"skip_strict_exist_check,omitempty"`
	Filters              *Filter_Config `protobuf:"bytes,2,opt,name=filters,proto3" json:"filters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Upsert_Config) Reset()         { *m = Upsert_Config{} }
func (m *Upsert_Config) String() string { return proto.CompactTextString(m) }
func (*Upsert_Config) ProtoMessage()    {}
func (*Upsert_Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{4, 2}
}
func (m *Upsert_Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Upsert_Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Upsert_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Upsert_Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Upsert_Config.Merge(m, src)
}
func (m *Upsert_Config) XXX_Size() int {
	return m.Size()
}
func (m *Upsert_Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Upsert_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Upsert_Config proto.InternalMessageInfo

func (m *Upsert_Config) GetSkipStrictExistCheck() bool {
	if m != nil {
		return m.SkipStrictExistCheck
	}
	return false
}

func (m *Upsert_Config) GetFilters() *Filter_Config {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (*Upsert_Config) XXX_MessageName() string {
	return "payload.v1.Upsert.Config"
}

type Remove struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Remove) Reset()         { *m = Remove{} }
func (m *Remove) String() string { return proto.CompactTextString(m) }
func (*Remove) ProtoMessage()    {}
func (*Remove) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{5}
}
func (m *Remove) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Remove) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Remove.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Remove) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Remove.Merge(m, src)
}
func (m *Remove) XXX_Size() int {
	return m.Size()
}
func (m *Remove) XXX_DiscardUnknown() {
	xxx_messageInfo_Remove.DiscardUnknown(m)
}

var xxx_messageInfo_Remove proto.InternalMessageInfo

func (*Remove) XXX_MessageName() string {
	return "payload.v1.Remove"
}

type Remove_Request struct {
	Id                   *Object_ID     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Config               *Remove_Config `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Remove_Request) Reset()         { *m = Remove_Request{} }
func (m *Remove_Request) String() string { return proto.CompactTextString(m) }
func (*Remove_Request) ProtoMessage()    {}
func (*Remove_Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{5, 0}
}
func (m *Remove_Request) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Remove_Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Remove_Request.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Remove_Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Remove_Request.Merge(m, src)
}
func (m *Remove_Request) XXX_Size() int {
	return m.Size()
}
func (m *Remove_Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Remove_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Remove_Request proto.InternalMessageInfo

func (m *Remove_Request) GetId() *Object_ID {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Remove_Request) GetConfig() *Remove_Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func (*Remove_Request) XXX_MessageName() string {
	return "payload.v1.Remove.Request"
}

type Remove_MultiRequest struct {
	Requests             []*Remove_Request `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Remove_MultiRequest) Reset()         { *m = Remove_MultiRequest{} }
func (m *Remove_MultiRequest) String() string { return proto.CompactTextString(m) }
func (*Remove_MultiRequest) ProtoMessage()    {}
func (*Remove_MultiRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{5, 1}
}
func (m *Remove_MultiRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Remove_MultiRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Remove_MultiRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Remove_MultiRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Remove_MultiRequest.Merge(m, src)
}
func (m *Remove_MultiRequest) XXX_Size() int {
	return m.Size()
}
func (m *Remove_MultiRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_Remove_MultiRequest.DiscardUnknown(m)
}

var xxx_messageInfo_Remove_MultiRequest proto.InternalMessageInfo

func (m *Remove_MultiRequest) GetRequests() []*Remove_Request {
	if m != nil {
		return m.Requests
	}
	return nil
}

func (*Remove_MultiRequest) XXX_MessageName() string {
	return "payload.v1.Remove.MultiRequest"
}

type Remove_Config struct {
	SkipStrictExistCheck bool     `protobuf:"varint,1,opt,name=skip_strict_exist_check,json=skipStrictExistCheck,proto3" json:"skip_strict_exist_check,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Remove_Config) Reset()         { *m = Remove_Config{} }
func (m *Remove_Config) String() string { return proto.CompactTextString(m) }
func (*Remove_Config) ProtoMessage()    {}
func (*Remove_Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{5, 2}
}
func (m *Remove_Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Remove_Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Remove_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Remove_Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Remove_Config.Merge(m, src)
}
func (m *Remove_Config) XXX_Size() int {
	return m.Size()
}
func (m *Remove_Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Remove_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Remove_Config proto.InternalMessageInfo

func (m *Remove_Config) GetSkipStrictExistCheck() bool {
	if m != nil {
		return m.SkipStrictExistCheck
	}
	return false
}

func (*Remove_Config) XXX_MessageName() string {
	return "payload.v1.Remove.Config"
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
	return fileDescriptor_f0518b37b4e7594b, []int{6}
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

func (*Meta) XXX_MessageName() string {
	return "payload.v1.Meta"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{6, 0}
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

func (*Meta_Key) XXX_MessageName() string {
	return "payload.v1.Meta.Key"
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
	return fileDescriptor_f0518b37b4e7594b, []int{6, 1}
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

func (*Meta_Keys) XXX_MessageName() string {
	return "payload.v1.Meta.Keys"
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
	return fileDescriptor_f0518b37b4e7594b, []int{6, 2}
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

func (*Meta_Val) XXX_MessageName() string {
	return "payload.v1.Meta.Val"
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
	return fileDescriptor_f0518b37b4e7594b, []int{6, 3}
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

func (*Meta_Vals) XXX_MessageName() string {
	return "payload.v1.Meta.Vals"
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
	return fileDescriptor_f0518b37b4e7594b, []int{6, 4}
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

func (*Meta_KeyVal) XXX_MessageName() string {
	return "payload.v1.Meta.KeyVal"
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
	return fileDescriptor_f0518b37b4e7594b, []int{6, 5}
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

func (*Meta_KeyVals) XXX_MessageName() string {
	return "payload.v1.Meta.KeyVals"
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
	return fileDescriptor_f0518b37b4e7594b, []int{7}
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

func (*Object) XXX_MessageName() string {
	return "payload.v1.Object"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{7, 0}
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

func (*Object_Distance) XXX_MessageName() string {
	return "payload.v1.Object.Distance"
}

type Object_StreamDistance struct {
	// Types that are valid to be assigned to Payload:
	//	*Object_StreamDistance_Distance
	//	*Object_StreamDistance_Status
	Payload              isObject_StreamDistance_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *Object_StreamDistance) Reset()         { *m = Object_StreamDistance{} }
func (m *Object_StreamDistance) String() string { return proto.CompactTextString(m) }
func (*Object_StreamDistance) ProtoMessage()    {}
func (*Object_StreamDistance) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{7, 1}
}
func (m *Object_StreamDistance) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_StreamDistance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_StreamDistance.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_StreamDistance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_StreamDistance.Merge(m, src)
}
func (m *Object_StreamDistance) XXX_Size() int {
	return m.Size()
}
func (m *Object_StreamDistance) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_StreamDistance.DiscardUnknown(m)
}

var xxx_messageInfo_Object_StreamDistance proto.InternalMessageInfo

type isObject_StreamDistance_Payload interface {
	isObject_StreamDistance_Payload()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Object_StreamDistance_Distance struct {
	Distance *Object_Distance `protobuf:"bytes,1,opt,name=distance,proto3,oneof" json:"distance,omitempty"`
}
type Object_StreamDistance_Status struct {
	Status *rpc.Status `protobuf:"bytes,2,opt,name=status,proto3,oneof" json:"status,omitempty"`
}

func (*Object_StreamDistance_Distance) isObject_StreamDistance_Payload() {}
func (*Object_StreamDistance_Status) isObject_StreamDistance_Payload()   {}

func (m *Object_StreamDistance) GetPayload() isObject_StreamDistance_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Object_StreamDistance) GetDistance() *Object_Distance {
	if x, ok := m.GetPayload().(*Object_StreamDistance_Distance); ok {
		return x.Distance
	}
	return nil
}

func (m *Object_StreamDistance) GetStatus() *rpc.Status {
	if x, ok := m.GetPayload().(*Object_StreamDistance_Status); ok {
		return x.Status
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Object_StreamDistance) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Object_StreamDistance_Distance)(nil),
		(*Object_StreamDistance_Status)(nil),
	}
}

func (*Object_StreamDistance) XXX_MessageName() string {
	return "payload.v1.Object.StreamDistance"
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
	return fileDescriptor_f0518b37b4e7594b, []int{7, 2}
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

func (*Object_ID) XXX_MessageName() string {
	return "payload.v1.Object.ID"
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
	return fileDescriptor_f0518b37b4e7594b, []int{7, 3}
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

func (*Object_IDs) XXX_MessageName() string {
	return "payload.v1.Object.IDs"
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
	return fileDescriptor_f0518b37b4e7594b, []int{7, 4}
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

func (*Object_Vector) XXX_MessageName() string {
	return "payload.v1.Object.Vector"
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
	return fileDescriptor_f0518b37b4e7594b, []int{7, 5}
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

func (*Object_Vectors) XXX_MessageName() string {
	return "payload.v1.Object.Vectors"
}

type Object_StreamVector struct {
	// Types that are valid to be assigned to Payload:
	//	*Object_StreamVector_Vector
	//	*Object_StreamVector_Status
	Payload              isObject_StreamVector_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *Object_StreamVector) Reset()         { *m = Object_StreamVector{} }
func (m *Object_StreamVector) String() string { return proto.CompactTextString(m) }
func (*Object_StreamVector) ProtoMessage()    {}
func (*Object_StreamVector) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{7, 6}
}
func (m *Object_StreamVector) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_StreamVector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_StreamVector.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_StreamVector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_StreamVector.Merge(m, src)
}
func (m *Object_StreamVector) XXX_Size() int {
	return m.Size()
}
func (m *Object_StreamVector) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_StreamVector.DiscardUnknown(m)
}

var xxx_messageInfo_Object_StreamVector proto.InternalMessageInfo

type isObject_StreamVector_Payload interface {
	isObject_StreamVector_Payload()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Object_StreamVector_Vector struct {
	Vector *Object_Vector `protobuf:"bytes,1,opt,name=vector,proto3,oneof" json:"vector,omitempty"`
}
type Object_StreamVector_Status struct {
	Status *rpc.Status `protobuf:"bytes,2,opt,name=status,proto3,oneof" json:"status,omitempty"`
}

func (*Object_StreamVector_Vector) isObject_StreamVector_Payload() {}
func (*Object_StreamVector_Status) isObject_StreamVector_Payload() {}

func (m *Object_StreamVector) GetPayload() isObject_StreamVector_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Object_StreamVector) GetVector() *Object_Vector {
	if x, ok := m.GetPayload().(*Object_StreamVector_Vector); ok {
		return x.Vector
	}
	return nil
}

func (m *Object_StreamVector) GetStatus() *rpc.Status {
	if x, ok := m.GetPayload().(*Object_StreamVector_Status); ok {
		return x.Status
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Object_StreamVector) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Object_StreamVector_Vector)(nil),
		(*Object_StreamVector_Status)(nil),
	}
}

func (*Object_StreamVector) XXX_MessageName() string {
	return "payload.v1.Object.StreamVector"
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
	return fileDescriptor_f0518b37b4e7594b, []int{7, 7}
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

func (*Object_Blob) XXX_MessageName() string {
	return "payload.v1.Object.Blob"
}

type Object_StreamBlob struct {
	// Types that are valid to be assigned to Payload:
	//	*Object_StreamBlob_Blob
	//	*Object_StreamBlob_Status
	Payload              isObject_StreamBlob_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *Object_StreamBlob) Reset()         { *m = Object_StreamBlob{} }
func (m *Object_StreamBlob) String() string { return proto.CompactTextString(m) }
func (*Object_StreamBlob) ProtoMessage()    {}
func (*Object_StreamBlob) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{7, 8}
}
func (m *Object_StreamBlob) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_StreamBlob) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_StreamBlob.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_StreamBlob) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_StreamBlob.Merge(m, src)
}
func (m *Object_StreamBlob) XXX_Size() int {
	return m.Size()
}
func (m *Object_StreamBlob) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_StreamBlob.DiscardUnknown(m)
}

var xxx_messageInfo_Object_StreamBlob proto.InternalMessageInfo

type isObject_StreamBlob_Payload interface {
	isObject_StreamBlob_Payload()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Object_StreamBlob_Blob struct {
	Blob *Object_Blob `protobuf:"bytes,1,opt,name=blob,proto3,oneof" json:"blob,omitempty"`
}
type Object_StreamBlob_Status struct {
	Status *rpc.Status `protobuf:"bytes,2,opt,name=status,proto3,oneof" json:"status,omitempty"`
}

func (*Object_StreamBlob_Blob) isObject_StreamBlob_Payload()   {}
func (*Object_StreamBlob_Status) isObject_StreamBlob_Payload() {}

func (m *Object_StreamBlob) GetPayload() isObject_StreamBlob_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Object_StreamBlob) GetBlob() *Object_Blob {
	if x, ok := m.GetPayload().(*Object_StreamBlob_Blob); ok {
		return x.Blob
	}
	return nil
}

func (m *Object_StreamBlob) GetStatus() *rpc.Status {
	if x, ok := m.GetPayload().(*Object_StreamBlob_Status); ok {
		return x.Status
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Object_StreamBlob) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Object_StreamBlob_Blob)(nil),
		(*Object_StreamBlob_Status)(nil),
	}
}

func (*Object_StreamBlob) XXX_MessageName() string {
	return "payload.v1.Object.StreamBlob"
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
	return fileDescriptor_f0518b37b4e7594b, []int{7, 9}
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

func (*Object_Location) XXX_MessageName() string {
	return "payload.v1.Object.Location"
}

type Object_StreamLocation struct {
	// Types that are valid to be assigned to Payload:
	//	*Object_StreamLocation_Location
	//	*Object_StreamLocation_Status
	Payload              isObject_StreamLocation_Payload `protobuf_oneof:"payload"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *Object_StreamLocation) Reset()         { *m = Object_StreamLocation{} }
func (m *Object_StreamLocation) String() string { return proto.CompactTextString(m) }
func (*Object_StreamLocation) ProtoMessage()    {}
func (*Object_StreamLocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{7, 10}
}
func (m *Object_StreamLocation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Object_StreamLocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Object_StreamLocation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Object_StreamLocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Object_StreamLocation.Merge(m, src)
}
func (m *Object_StreamLocation) XXX_Size() int {
	return m.Size()
}
func (m *Object_StreamLocation) XXX_DiscardUnknown() {
	xxx_messageInfo_Object_StreamLocation.DiscardUnknown(m)
}

var xxx_messageInfo_Object_StreamLocation proto.InternalMessageInfo

type isObject_StreamLocation_Payload interface {
	isObject_StreamLocation_Payload()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Object_StreamLocation_Location struct {
	Location *Object_Location `protobuf:"bytes,1,opt,name=location,proto3,oneof" json:"location,omitempty"`
}
type Object_StreamLocation_Status struct {
	Status *rpc.Status `protobuf:"bytes,2,opt,name=status,proto3,oneof" json:"status,omitempty"`
}

func (*Object_StreamLocation_Location) isObject_StreamLocation_Payload() {}
func (*Object_StreamLocation_Status) isObject_StreamLocation_Payload()   {}

func (m *Object_StreamLocation) GetPayload() isObject_StreamLocation_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Object_StreamLocation) GetLocation() *Object_Location {
	if x, ok := m.GetPayload().(*Object_StreamLocation_Location); ok {
		return x.Location
	}
	return nil
}

func (m *Object_StreamLocation) GetStatus() *rpc.Status {
	if x, ok := m.GetPayload().(*Object_StreamLocation_Status); ok {
		return x.Status
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Object_StreamLocation) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Object_StreamLocation_Location)(nil),
		(*Object_StreamLocation_Status)(nil),
	}
}

func (*Object_StreamLocation) XXX_MessageName() string {
	return "payload.v1.Object.StreamLocation"
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
	return fileDescriptor_f0518b37b4e7594b, []int{7, 11}
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

func (*Object_Locations) XXX_MessageName() string {
	return "payload.v1.Object.Locations"
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
	return fileDescriptor_f0518b37b4e7594b, []int{8}
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

func (*Control) XXX_MessageName() string {
	return "payload.v1.Control"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{8, 0}
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

func (*Control_CreateIndexRequest) XXX_MessageName() string {
	return "payload.v1.Control.CreateIndexRequest"
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
	return fileDescriptor_f0518b37b4e7594b, []int{9}
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

func (*Replication) XXX_MessageName() string {
	return "payload.v1.Replication"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{9, 0}
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

func (*Replication_Recovery) XXX_MessageName() string {
	return "payload.v1.Replication.Recovery"
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
	return fileDescriptor_f0518b37b4e7594b, []int{9, 1}
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

func (*Replication_Rebalance) XXX_MessageName() string {
	return "payload.v1.Replication.Rebalance"
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
	return fileDescriptor_f0518b37b4e7594b, []int{9, 2}
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

func (*Replication_Agents) XXX_MessageName() string {
	return "payload.v1.Replication.Agents"
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
	return fileDescriptor_f0518b37b4e7594b, []int{10}
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

func (*Discoverer) XXX_MessageName() string {
	return "payload.v1.Discoverer"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{10, 0}
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

func (*Discoverer_Request) XXX_MessageName() string {
	return "payload.v1.Discoverer.Request"
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
	return fileDescriptor_f0518b37b4e7594b, []int{11}
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

func (*Backup) XXX_MessageName() string {
	return "payload.v1.Backup"
}

type Backup_GetVector struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_GetVector) Reset()         { *m = Backup_GetVector{} }
func (m *Backup_GetVector) String() string { return proto.CompactTextString(m) }
func (*Backup_GetVector) ProtoMessage()    {}
func (*Backup_GetVector) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{11, 0}
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

func (*Backup_GetVector) XXX_MessageName() string {
	return "payload.v1.Backup.GetVector"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 0, 0}
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

func (*Backup_GetVector_Request) XXX_MessageName() string {
	return "payload.v1.Backup.GetVector.Request"
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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 0, 1}
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

func (*Backup_GetVector_Owner) XXX_MessageName() string {
	return "payload.v1.Backup.GetVector.Owner"
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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 1}
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

func (*Backup_Locations) XXX_MessageName() string {
	return "payload.v1.Backup.Locations"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 1, 0}
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

func (*Backup_Locations_Request) XXX_MessageName() string {
	return "payload.v1.Backup.Locations.Request"
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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 2}
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

func (*Backup_Remove) XXX_MessageName() string {
	return "payload.v1.Backup.Remove"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 2, 0}
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

func (*Backup_Remove_Request) XXX_MessageName() string {
	return "payload.v1.Backup.Remove.Request"
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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 2, 1}
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

func (*Backup_Remove_RequestMulti) XXX_MessageName() string {
	return "payload.v1.Backup.Remove.RequestMulti"
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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 3}
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

func (*Backup_IP) XXX_MessageName() string {
	return "payload.v1.Backup.IP"
}

type Backup_IP_Register struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_IP_Register) Reset()         { *m = Backup_IP_Register{} }
func (m *Backup_IP_Register) String() string { return proto.CompactTextString(m) }
func (*Backup_IP_Register) ProtoMessage()    {}
func (*Backup_IP_Register) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{11, 3, 0}
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

func (*Backup_IP_Register) XXX_MessageName() string {
	return "payload.v1.Backup.IP.Register"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 3, 0, 0}
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

func (*Backup_IP_Register_Request) XXX_MessageName() string {
	return "payload.v1.Backup.IP.Register.Request"
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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 3, 1}
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

func (*Backup_IP_Remove) XXX_MessageName() string {
	return "payload.v1.Backup.IP.Remove"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 3, 1, 0}
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

func (*Backup_IP_Remove_Request) XXX_MessageName() string {
	return "payload.v1.Backup.IP.Remove.Request"
}

type Backup_Vector struct {
	Uuid                 string    `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Vector               []float32 `protobuf:"fixed32,3,rep,packed,name=vector,proto3" json:"vector,omitempty"`
	Ips                  []string  `protobuf:"bytes,4,rep,name=ips,proto3" json:"ips,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Backup_Vector) Reset()         { *m = Backup_Vector{} }
func (m *Backup_Vector) String() string { return proto.CompactTextString(m) }
func (*Backup_Vector) ProtoMessage()    {}
func (*Backup_Vector) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{11, 4}
}
func (m *Backup_Vector) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Vector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Vector.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Vector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Vector.Merge(m, src)
}
func (m *Backup_Vector) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Vector) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Vector.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Vector proto.InternalMessageInfo

func (m *Backup_Vector) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Backup_Vector) GetVector() []float32 {
	if m != nil {
		return m.Vector
	}
	return nil
}

func (m *Backup_Vector) GetIps() []string {
	if m != nil {
		return m.Ips
	}
	return nil
}

func (*Backup_Vector) XXX_MessageName() string {
	return "payload.v1.Backup.Vector"
}

type Backup_Vectors struct {
	Vectors              []*Backup_Vector `protobuf:"bytes,1,rep,name=vectors,proto3" json:"vectors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Backup_Vectors) Reset()         { *m = Backup_Vectors{} }
func (m *Backup_Vectors) String() string { return proto.CompactTextString(m) }
func (*Backup_Vectors) ProtoMessage()    {}
func (*Backup_Vectors) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{11, 5}
}
func (m *Backup_Vectors) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Vectors) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Vectors.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Vectors) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Vectors.Merge(m, src)
}
func (m *Backup_Vectors) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Vectors) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Vectors.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Vectors proto.InternalMessageInfo

func (m *Backup_Vectors) GetVectors() []*Backup_Vector {
	if m != nil {
		return m.Vectors
	}
	return nil
}

func (*Backup_Vectors) XXX_MessageName() string {
	return "payload.v1.Backup.Vectors"
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
	return fileDescriptor_f0518b37b4e7594b, []int{11, 6}
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

func (*Backup_Compressed) XXX_MessageName() string {
	return "payload.v1.Backup.Compressed"
}

type Backup_Compressed_Vector struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Vector               []byte   `protobuf:"bytes,3,opt,name=vector,proto3" json:"vector,omitempty"`
	Ips                  []string `protobuf:"bytes,4,rep,name=ips,proto3" json:"ips,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Backup_Compressed_Vector) Reset()         { *m = Backup_Compressed_Vector{} }
func (m *Backup_Compressed_Vector) String() string { return proto.CompactTextString(m) }
func (*Backup_Compressed_Vector) ProtoMessage()    {}
func (*Backup_Compressed_Vector) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{11, 6, 0}
}
func (m *Backup_Compressed_Vector) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Compressed_Vector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Compressed_Vector.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Compressed_Vector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Compressed_Vector.Merge(m, src)
}
func (m *Backup_Compressed_Vector) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Compressed_Vector) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Compressed_Vector.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Compressed_Vector proto.InternalMessageInfo

func (m *Backup_Compressed_Vector) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *Backup_Compressed_Vector) GetVector() []byte {
	if m != nil {
		return m.Vector
	}
	return nil
}

func (m *Backup_Compressed_Vector) GetIps() []string {
	if m != nil {
		return m.Ips
	}
	return nil
}

func (*Backup_Compressed_Vector) XXX_MessageName() string {
	return "payload.v1.Backup.Compressed.Vector"
}

type Backup_Compressed_Vectors struct {
	Vectors              []*Backup_Compressed_Vector `protobuf:"bytes,1,rep,name=vectors,proto3" json:"vectors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *Backup_Compressed_Vectors) Reset()         { *m = Backup_Compressed_Vectors{} }
func (m *Backup_Compressed_Vectors) String() string { return proto.CompactTextString(m) }
func (*Backup_Compressed_Vectors) ProtoMessage()    {}
func (*Backup_Compressed_Vectors) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{11, 6, 1}
}
func (m *Backup_Compressed_Vectors) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Backup_Compressed_Vectors) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Backup_Compressed_Vectors.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Backup_Compressed_Vectors) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Backup_Compressed_Vectors.Merge(m, src)
}
func (m *Backup_Compressed_Vectors) XXX_Size() int {
	return m.Size()
}
func (m *Backup_Compressed_Vectors) XXX_DiscardUnknown() {
	xxx_messageInfo_Backup_Compressed_Vectors.DiscardUnknown(m)
}

var xxx_messageInfo_Backup_Compressed_Vectors proto.InternalMessageInfo

func (m *Backup_Compressed_Vectors) GetVectors() []*Backup_Compressed_Vector {
	if m != nil {
		return m.Vectors
	}
	return nil
}

func (*Backup_Compressed_Vectors) XXX_MessageName() string {
	return "payload.v1.Backup.Compressed.Vectors"
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
	return fileDescriptor_f0518b37b4e7594b, []int{12}
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

func (*Info) XXX_MessageName() string {
	return "payload.v1.Info"
}

type Info_Index struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Info_Index) Reset()         { *m = Info_Index{} }
func (m *Info_Index) String() string { return proto.CompactTextString(m) }
func (*Info_Index) ProtoMessage()    {}
func (*Info_Index) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0518b37b4e7594b, []int{12, 0}
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

func (*Info_Index) XXX_MessageName() string {
	return "payload.v1.Info.Index"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{12, 0, 0}
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

func (*Info_Index_Count) XXX_MessageName() string {
	return "payload.v1.Info.Index.Count"
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
	return fileDescriptor_f0518b37b4e7594b, []int{12, 0, 1}
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

func (*Info_Index_UUID) XXX_MessageName() string {
	return "payload.v1.Info.Index.UUID"
}

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
	return fileDescriptor_f0518b37b4e7594b, []int{12, 0, 1, 0}
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

func (*Info_Index_UUID_Committed) XXX_MessageName() string {
	return "payload.v1.Info.Index.UUID.Committed"
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
	return fileDescriptor_f0518b37b4e7594b, []int{12, 0, 1, 1}
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

func (*Info_Index_UUID_Uncommitted) XXX_MessageName() string {
	return "payload.v1.Info.Index.UUID.Uncommitted"
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
	return fileDescriptor_f0518b37b4e7594b, []int{12, 1}
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

func (*Info_Pod) XXX_MessageName() string {
	return "payload.v1.Info.Pod"
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
	return fileDescriptor_f0518b37b4e7594b, []int{12, 2}
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

func (*Info_Node) XXX_MessageName() string {
	return "payload.v1.Info.Node"
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
	return fileDescriptor_f0518b37b4e7594b, []int{12, 3}
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

func (*Info_CPU) XXX_MessageName() string {
	return "payload.v1.Info.CPU"
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
	return fileDescriptor_f0518b37b4e7594b, []int{12, 4}
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

func (*Info_Memory) XXX_MessageName() string {
	return "payload.v1.Info.Memory"
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
	return fileDescriptor_f0518b37b4e7594b, []int{12, 5}
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

func (*Info_Pods) XXX_MessageName() string {
	return "payload.v1.Info.Pods"
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
	return fileDescriptor_f0518b37b4e7594b, []int{12, 6}
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

func (*Info_Nodes) XXX_MessageName() string {
	return "payload.v1.Info.Nodes"
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
	return fileDescriptor_f0518b37b4e7594b, []int{12, 7}
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

func (*Info_IPs) XXX_MessageName() string {
	return "payload.v1.Info.IPs"
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
	return fileDescriptor_f0518b37b4e7594b, []int{13}
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

func (*Empty) XXX_MessageName() string {
	return "payload.v1.Empty"
}
func init() {
	proto.RegisterType((*Search)(nil), "payload.v1.Search")
	golang_proto.RegisterType((*Search)(nil), "payload.v1.Search")
	proto.RegisterType((*Search_Request)(nil), "payload.v1.Search.Request")
	golang_proto.RegisterType((*Search_Request)(nil), "payload.v1.Search.Request")
	proto.RegisterType((*Search_MultiRequest)(nil), "payload.v1.Search.MultiRequest")
	golang_proto.RegisterType((*Search_MultiRequest)(nil), "payload.v1.Search.MultiRequest")
	proto.RegisterType((*Search_IDRequest)(nil), "payload.v1.Search.IDRequest")
	golang_proto.RegisterType((*Search_IDRequest)(nil), "payload.v1.Search.IDRequest")
	proto.RegisterType((*Search_MultiIDRequest)(nil), "payload.v1.Search.MultiIDRequest")
	golang_proto.RegisterType((*Search_MultiIDRequest)(nil), "payload.v1.Search.MultiIDRequest")
	proto.RegisterType((*Search_ObjectRequest)(nil), "payload.v1.Search.ObjectRequest")
	golang_proto.RegisterType((*Search_ObjectRequest)(nil), "payload.v1.Search.ObjectRequest")
	proto.RegisterType((*Search_Config)(nil), "payload.v1.Search.Config")
	golang_proto.RegisterType((*Search_Config)(nil), "payload.v1.Search.Config")
	proto.RegisterType((*Search_Response)(nil), "payload.v1.Search.Response")
	golang_proto.RegisterType((*Search_Response)(nil), "payload.v1.Search.Response")
	proto.RegisterType((*Search_Responses)(nil), "payload.v1.Search.Responses")
	golang_proto.RegisterType((*Search_Responses)(nil), "payload.v1.Search.Responses")
	proto.RegisterType((*Search_StreamResponse)(nil), "payload.v1.Search.StreamResponse")
	golang_proto.RegisterType((*Search_StreamResponse)(nil), "payload.v1.Search.StreamResponse")
	proto.RegisterType((*Filter)(nil), "payload.v1.Filter")
	golang_proto.RegisterType((*Filter)(nil), "payload.v1.Filter")
	proto.RegisterType((*Filter_Target)(nil), "payload.v1.Filter.Target")
	golang_proto.RegisterType((*Filter_Target)(nil), "payload.v1.Filter.Target")
	proto.RegisterType((*Filter_Config)(nil), "payload.v1.Filter.Config")
	golang_proto.RegisterType((*Filter_Config)(nil), "payload.v1.Filter.Config")
	proto.RegisterType((*Insert)(nil), "payload.v1.Insert")
	golang_proto.RegisterType((*Insert)(nil), "payload.v1.Insert")
	proto.RegisterType((*Insert_Request)(nil), "payload.v1.Insert.Request")
	golang_proto.RegisterType((*Insert_Request)(nil), "payload.v1.Insert.Request")
	proto.RegisterType((*Insert_MultiRequest)(nil), "payload.v1.Insert.MultiRequest")
	golang_proto.RegisterType((*Insert_MultiRequest)(nil), "payload.v1.Insert.MultiRequest")
	proto.RegisterType((*Insert_Config)(nil), "payload.v1.Insert.Config")
	golang_proto.RegisterType((*Insert_Config)(nil), "payload.v1.Insert.Config")
	proto.RegisterType((*Update)(nil), "payload.v1.Update")
	golang_proto.RegisterType((*Update)(nil), "payload.v1.Update")
	proto.RegisterType((*Update_Request)(nil), "payload.v1.Update.Request")
	golang_proto.RegisterType((*Update_Request)(nil), "payload.v1.Update.Request")
	proto.RegisterType((*Update_MultiRequest)(nil), "payload.v1.Update.MultiRequest")
	golang_proto.RegisterType((*Update_MultiRequest)(nil), "payload.v1.Update.MultiRequest")
	proto.RegisterType((*Update_Config)(nil), "payload.v1.Update.Config")
	golang_proto.RegisterType((*Update_Config)(nil), "payload.v1.Update.Config")
	proto.RegisterType((*Upsert)(nil), "payload.v1.Upsert")
	golang_proto.RegisterType((*Upsert)(nil), "payload.v1.Upsert")
	proto.RegisterType((*Upsert_Request)(nil), "payload.v1.Upsert.Request")
	golang_proto.RegisterType((*Upsert_Request)(nil), "payload.v1.Upsert.Request")
	proto.RegisterType((*Upsert_MultiRequest)(nil), "payload.v1.Upsert.MultiRequest")
	golang_proto.RegisterType((*Upsert_MultiRequest)(nil), "payload.v1.Upsert.MultiRequest")
	proto.RegisterType((*Upsert_Config)(nil), "payload.v1.Upsert.Config")
	golang_proto.RegisterType((*Upsert_Config)(nil), "payload.v1.Upsert.Config")
	proto.RegisterType((*Remove)(nil), "payload.v1.Remove")
	golang_proto.RegisterType((*Remove)(nil), "payload.v1.Remove")
	proto.RegisterType((*Remove_Request)(nil), "payload.v1.Remove.Request")
	golang_proto.RegisterType((*Remove_Request)(nil), "payload.v1.Remove.Request")
	proto.RegisterType((*Remove_MultiRequest)(nil), "payload.v1.Remove.MultiRequest")
	golang_proto.RegisterType((*Remove_MultiRequest)(nil), "payload.v1.Remove.MultiRequest")
	proto.RegisterType((*Remove_Config)(nil), "payload.v1.Remove.Config")
	golang_proto.RegisterType((*Remove_Config)(nil), "payload.v1.Remove.Config")
	proto.RegisterType((*Meta)(nil), "payload.v1.Meta")
	golang_proto.RegisterType((*Meta)(nil), "payload.v1.Meta")
	proto.RegisterType((*Meta_Key)(nil), "payload.v1.Meta.Key")
	golang_proto.RegisterType((*Meta_Key)(nil), "payload.v1.Meta.Key")
	proto.RegisterType((*Meta_Keys)(nil), "payload.v1.Meta.Keys")
	golang_proto.RegisterType((*Meta_Keys)(nil), "payload.v1.Meta.Keys")
	proto.RegisterType((*Meta_Val)(nil), "payload.v1.Meta.Val")
	golang_proto.RegisterType((*Meta_Val)(nil), "payload.v1.Meta.Val")
	proto.RegisterType((*Meta_Vals)(nil), "payload.v1.Meta.Vals")
	golang_proto.RegisterType((*Meta_Vals)(nil), "payload.v1.Meta.Vals")
	proto.RegisterType((*Meta_KeyVal)(nil), "payload.v1.Meta.KeyVal")
	golang_proto.RegisterType((*Meta_KeyVal)(nil), "payload.v1.Meta.KeyVal")
	proto.RegisterType((*Meta_KeyVals)(nil), "payload.v1.Meta.KeyVals")
	golang_proto.RegisterType((*Meta_KeyVals)(nil), "payload.v1.Meta.KeyVals")
	proto.RegisterType((*Object)(nil), "payload.v1.Object")
	golang_proto.RegisterType((*Object)(nil), "payload.v1.Object")
	proto.RegisterType((*Object_Distance)(nil), "payload.v1.Object.Distance")
	golang_proto.RegisterType((*Object_Distance)(nil), "payload.v1.Object.Distance")
	proto.RegisterType((*Object_StreamDistance)(nil), "payload.v1.Object.StreamDistance")
	golang_proto.RegisterType((*Object_StreamDistance)(nil), "payload.v1.Object.StreamDistance")
	proto.RegisterType((*Object_ID)(nil), "payload.v1.Object.ID")
	golang_proto.RegisterType((*Object_ID)(nil), "payload.v1.Object.ID")
	proto.RegisterType((*Object_IDs)(nil), "payload.v1.Object.IDs")
	golang_proto.RegisterType((*Object_IDs)(nil), "payload.v1.Object.IDs")
	proto.RegisterType((*Object_Vector)(nil), "payload.v1.Object.Vector")
	golang_proto.RegisterType((*Object_Vector)(nil), "payload.v1.Object.Vector")
	proto.RegisterType((*Object_Vectors)(nil), "payload.v1.Object.Vectors")
	golang_proto.RegisterType((*Object_Vectors)(nil), "payload.v1.Object.Vectors")
	proto.RegisterType((*Object_StreamVector)(nil), "payload.v1.Object.StreamVector")
	golang_proto.RegisterType((*Object_StreamVector)(nil), "payload.v1.Object.StreamVector")
	proto.RegisterType((*Object_Blob)(nil), "payload.v1.Object.Blob")
	golang_proto.RegisterType((*Object_Blob)(nil), "payload.v1.Object.Blob")
	proto.RegisterType((*Object_StreamBlob)(nil), "payload.v1.Object.StreamBlob")
	golang_proto.RegisterType((*Object_StreamBlob)(nil), "payload.v1.Object.StreamBlob")
	proto.RegisterType((*Object_Location)(nil), "payload.v1.Object.Location")
	golang_proto.RegisterType((*Object_Location)(nil), "payload.v1.Object.Location")
	proto.RegisterType((*Object_StreamLocation)(nil), "payload.v1.Object.StreamLocation")
	golang_proto.RegisterType((*Object_StreamLocation)(nil), "payload.v1.Object.StreamLocation")
	proto.RegisterType((*Object_Locations)(nil), "payload.v1.Object.Locations")
	golang_proto.RegisterType((*Object_Locations)(nil), "payload.v1.Object.Locations")
	proto.RegisterType((*Control)(nil), "payload.v1.Control")
	golang_proto.RegisterType((*Control)(nil), "payload.v1.Control")
	proto.RegisterType((*Control_CreateIndexRequest)(nil), "payload.v1.Control.CreateIndexRequest")
	golang_proto.RegisterType((*Control_CreateIndexRequest)(nil), "payload.v1.Control.CreateIndexRequest")
	proto.RegisterType((*Replication)(nil), "payload.v1.Replication")
	golang_proto.RegisterType((*Replication)(nil), "payload.v1.Replication")
	proto.RegisterType((*Replication_Recovery)(nil), "payload.v1.Replication.Recovery")
	golang_proto.RegisterType((*Replication_Recovery)(nil), "payload.v1.Replication.Recovery")
	proto.RegisterType((*Replication_Rebalance)(nil), "payload.v1.Replication.Rebalance")
	golang_proto.RegisterType((*Replication_Rebalance)(nil), "payload.v1.Replication.Rebalance")
	proto.RegisterType((*Replication_Agents)(nil), "payload.v1.Replication.Agents")
	golang_proto.RegisterType((*Replication_Agents)(nil), "payload.v1.Replication.Agents")
	proto.RegisterType((*Discoverer)(nil), "payload.v1.Discoverer")
	golang_proto.RegisterType((*Discoverer)(nil), "payload.v1.Discoverer")
	proto.RegisterType((*Discoverer_Request)(nil), "payload.v1.Discoverer.Request")
	golang_proto.RegisterType((*Discoverer_Request)(nil), "payload.v1.Discoverer.Request")
	proto.RegisterType((*Backup)(nil), "payload.v1.Backup")
	golang_proto.RegisterType((*Backup)(nil), "payload.v1.Backup")
	proto.RegisterType((*Backup_GetVector)(nil), "payload.v1.Backup.GetVector")
	golang_proto.RegisterType((*Backup_GetVector)(nil), "payload.v1.Backup.GetVector")
	proto.RegisterType((*Backup_GetVector_Request)(nil), "payload.v1.Backup.GetVector.Request")
	golang_proto.RegisterType((*Backup_GetVector_Request)(nil), "payload.v1.Backup.GetVector.Request")
	proto.RegisterType((*Backup_GetVector_Owner)(nil), "payload.v1.Backup.GetVector.Owner")
	golang_proto.RegisterType((*Backup_GetVector_Owner)(nil), "payload.v1.Backup.GetVector.Owner")
	proto.RegisterType((*Backup_Locations)(nil), "payload.v1.Backup.Locations")
	golang_proto.RegisterType((*Backup_Locations)(nil), "payload.v1.Backup.Locations")
	proto.RegisterType((*Backup_Locations_Request)(nil), "payload.v1.Backup.Locations.Request")
	golang_proto.RegisterType((*Backup_Locations_Request)(nil), "payload.v1.Backup.Locations.Request")
	proto.RegisterType((*Backup_Remove)(nil), "payload.v1.Backup.Remove")
	golang_proto.RegisterType((*Backup_Remove)(nil), "payload.v1.Backup.Remove")
	proto.RegisterType((*Backup_Remove_Request)(nil), "payload.v1.Backup.Remove.Request")
	golang_proto.RegisterType((*Backup_Remove_Request)(nil), "payload.v1.Backup.Remove.Request")
	proto.RegisterType((*Backup_Remove_RequestMulti)(nil), "payload.v1.Backup.Remove.RequestMulti")
	golang_proto.RegisterType((*Backup_Remove_RequestMulti)(nil), "payload.v1.Backup.Remove.RequestMulti")
	proto.RegisterType((*Backup_IP)(nil), "payload.v1.Backup.IP")
	golang_proto.RegisterType((*Backup_IP)(nil), "payload.v1.Backup.IP")
	proto.RegisterType((*Backup_IP_Register)(nil), "payload.v1.Backup.IP.Register")
	golang_proto.RegisterType((*Backup_IP_Register)(nil), "payload.v1.Backup.IP.Register")
	proto.RegisterType((*Backup_IP_Register_Request)(nil), "payload.v1.Backup.IP.Register.Request")
	golang_proto.RegisterType((*Backup_IP_Register_Request)(nil), "payload.v1.Backup.IP.Register.Request")
	proto.RegisterType((*Backup_IP_Remove)(nil), "payload.v1.Backup.IP.Remove")
	golang_proto.RegisterType((*Backup_IP_Remove)(nil), "payload.v1.Backup.IP.Remove")
	proto.RegisterType((*Backup_IP_Remove_Request)(nil), "payload.v1.Backup.IP.Remove.Request")
	golang_proto.RegisterType((*Backup_IP_Remove_Request)(nil), "payload.v1.Backup.IP.Remove.Request")
	proto.RegisterType((*Backup_Vector)(nil), "payload.v1.Backup.Vector")
	golang_proto.RegisterType((*Backup_Vector)(nil), "payload.v1.Backup.Vector")
	proto.RegisterType((*Backup_Vectors)(nil), "payload.v1.Backup.Vectors")
	golang_proto.RegisterType((*Backup_Vectors)(nil), "payload.v1.Backup.Vectors")
	proto.RegisterType((*Backup_Compressed)(nil), "payload.v1.Backup.Compressed")
	golang_proto.RegisterType((*Backup_Compressed)(nil), "payload.v1.Backup.Compressed")
	proto.RegisterType((*Backup_Compressed_Vector)(nil), "payload.v1.Backup.Compressed.Vector")
	golang_proto.RegisterType((*Backup_Compressed_Vector)(nil), "payload.v1.Backup.Compressed.Vector")
	proto.RegisterType((*Backup_Compressed_Vectors)(nil), "payload.v1.Backup.Compressed.Vectors")
	golang_proto.RegisterType((*Backup_Compressed_Vectors)(nil), "payload.v1.Backup.Compressed.Vectors")
	proto.RegisterType((*Info)(nil), "payload.v1.Info")
	golang_proto.RegisterType((*Info)(nil), "payload.v1.Info")
	proto.RegisterType((*Info_Index)(nil), "payload.v1.Info.Index")
	golang_proto.RegisterType((*Info_Index)(nil), "payload.v1.Info.Index")
	proto.RegisterType((*Info_Index_Count)(nil), "payload.v1.Info.Index.Count")
	golang_proto.RegisterType((*Info_Index_Count)(nil), "payload.v1.Info.Index.Count")
	proto.RegisterType((*Info_Index_UUID)(nil), "payload.v1.Info.Index.UUID")
	golang_proto.RegisterType((*Info_Index_UUID)(nil), "payload.v1.Info.Index.UUID")
	proto.RegisterType((*Info_Index_UUID_Committed)(nil), "payload.v1.Info.Index.UUID.Committed")
	golang_proto.RegisterType((*Info_Index_UUID_Committed)(nil), "payload.v1.Info.Index.UUID.Committed")
	proto.RegisterType((*Info_Index_UUID_Uncommitted)(nil), "payload.v1.Info.Index.UUID.Uncommitted")
	golang_proto.RegisterType((*Info_Index_UUID_Uncommitted)(nil), "payload.v1.Info.Index.UUID.Uncommitted")
	proto.RegisterType((*Info_Pod)(nil), "payload.v1.Info.Pod")
	golang_proto.RegisterType((*Info_Pod)(nil), "payload.v1.Info.Pod")
	proto.RegisterType((*Info_Node)(nil), "payload.v1.Info.Node")
	golang_proto.RegisterType((*Info_Node)(nil), "payload.v1.Info.Node")
	proto.RegisterType((*Info_CPU)(nil), "payload.v1.Info.CPU")
	golang_proto.RegisterType((*Info_CPU)(nil), "payload.v1.Info.CPU")
	proto.RegisterType((*Info_Memory)(nil), "payload.v1.Info.Memory")
	golang_proto.RegisterType((*Info_Memory)(nil), "payload.v1.Info.Memory")
	proto.RegisterType((*Info_Pods)(nil), "payload.v1.Info.Pods")
	golang_proto.RegisterType((*Info_Pods)(nil), "payload.v1.Info.Pods")
	proto.RegisterType((*Info_Nodes)(nil), "payload.v1.Info.Nodes")
	golang_proto.RegisterType((*Info_Nodes)(nil), "payload.v1.Info.Nodes")
	proto.RegisterType((*Info_IPs)(nil), "payload.v1.Info.IPs")
	golang_proto.RegisterType((*Info_IPs)(nil), "payload.v1.Info.IPs")
	proto.RegisterType((*Empty)(nil), "payload.v1.Empty")
	golang_proto.RegisterType((*Empty)(nil), "payload.v1.Empty")
}

func init() {
	proto.RegisterFile("apis/proto/v1/payload/payload.proto", fileDescriptor_f0518b37b4e7594b)
}
func init() {
	golang_proto.RegisterFile("apis/proto/v1/payload/payload.proto", fileDescriptor_f0518b37b4e7594b)
}

var fileDescriptor_f0518b37b4e7594b = []byte{
	// 1896 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x58, 0x4f, 0x6f, 0x1b, 0xc7,
	0x15, 0xf7, 0x2e, 0xc9, 0x25, 0xf9, 0xf4, 0xa7, 0xce, 0xc2, 0xb1, 0xe8, 0x71, 0xac, 0xaa, 0x4c,
	0x1c, 0x28, 0x6e, 0x4c, 0xd6, 0x72, 0x9d, 0xd6, 0x29, 0xe0, 0xc2, 0xa4, 0xec, 0x9a, 0x71, 0xad,
	0x10, 0xab, 0x4a, 0x28, 0x02, 0x34, 0xcc, 0x68, 0x77, 0x44, 0x6d, 0xb5, 0xdc, 0xd9, 0xee, 0x0c,
	0x69, 0x31, 0xe8, 0xb1, 0xfd, 0x02, 0xbd, 0xf5, 0x94, 0xef, 0x90, 0x4f, 0xd0, 0xa2, 0x40, 0xe1,
	0x63, 0x0f, 0xfd, 0x00, 0x85, 0x73, 0xe9, 0xad, 0xb7, 0x1e, 0x72, 0x2a, 0xe6, 0xdf, 0xee, 0x4a,
	0xa4, 0x1c, 0x55, 0x0a, 0x02, 0x9f, 0x38, 0xf3, 0xe6, 0xf7, 0xde, 0xbc, 0xf7, 0xe6, 0xfd, 0xe3,
	0xc2, 0xdb, 0x38, 0x09, 0x59, 0x3b, 0x49, 0x29, 0xa7, 0xed, 0xc9, 0x9d, 0x76, 0x82, 0xa7, 0x11,
	0xc5, 0x81, 0xf9, 0x6d, 0xc9, 0x03, 0x17, 0xcc, 0x76, 0x72, 0x07, 0xdd, 0x1e, 0x86, 0xfc, 0x60,
	0xbc, 0xd7, 0xf2, 0xe9, 0xa8, 0x3d, 0xa4, 0x43, 0xaa, 0x78, 0xf7, 0xc6, 0xfb, 0x72, 0xa7, 0x04,
	0x89, 0x95, 0x62, 0x45, 0x1b, 0x27, 0xe1, 0x43, 0x4a, 0x87, 0x11, 0x91, 0xb7, 0xaa, 0x65, 0x3b,
	0x4d, 0xfc, 0x36, 0xe3, 0x98, 0x8f, 0x99, 0xe6, 0x79, 0x58, 0xe0, 0x21, 0xf1, 0x84, 0x4e, 0x93,
	0x94, 0x1e, 0x4d, 0xd5, 0x45, 0xfe, 0xed, 0x21, 0x89, 0x6f, 0x4f, 0x70, 0x14, 0x06, 0x98, 0x93,
	0xf6, 0xcc, 0x42, 0x89, 0x68, 0x7e, 0x51, 0x05, 0x67, 0x9b, 0xe0, 0xd4, 0x3f, 0x40, 0x9f, 0x42,
	0xd5, 0x23, 0xbf, 0x1b, 0x13, 0xc6, 0xdd, 0x35, 0x70, 0x26, 0xc4, 0xe7, 0x34, 0x6d, 0x58, 0x6b,
	0xa5, 0x75, 0xbb, 0x53, 0xfb, 0xba, 0x53, 0xf9, 0x93, 0x65, 0xd7, 0x6c, 0x4f, 0xd3, 0xdd, 0x3b,
	0xe0, 0xf8, 0x34, 0xde, 0x0f, 0x87, 0x0d, 0x7b, 0xcd, 0x5a, 0x5f, 0xd8, 0xb8, 0xd6, 0xca, 0x4d,
	0x6f, 0x29, 0x81, 0xad, 0xae, 0x04, 0x78, 0x1a, 0x88, 0x1e, 0xc3, 0xe2, 0xb3, 0x71, 0xc4, 0x43,
	0x73, 0xc9, 0x07, 0x50, 0x4b, 0xd5, 0x92, 0xc9, 0x6b, 0x16, 0x36, 0xd0, 0x1c, 0x21, 0x1a, 0xed,
	0x65, 0x58, 0xb4, 0x05, 0xf5, 0xde, 0xa6, 0x11, 0xb2, 0x0c, 0x76, 0x18, 0x34, 0xac, 0x35, 0x6b,
	0xbd, 0xee, 0xd9, 0x61, 0x70, 0x1e, 0xbd, 0x3e, 0x82, 0x65, 0xa9, 0x57, 0x2e, 0xf4, 0xa7, 0x33,
	0x9a, 0xbd, 0x35, 0x47, 0x4c, 0x86, 0x2f, 0xe8, 0xf6, 0x09, 0x2c, 0x7d, 0xbc, 0xf7, 0x5b, 0xe2,
	0x73, 0x23, 0xea, 0x2a, 0x38, 0x54, 0x12, 0xa4, 0x8e, 0x8b, 0x9e, 0xde, 0x9d, 0x47, 0xcf, 0xbf,
	0x59, 0xe0, 0x28, 0x92, 0x7b, 0x03, 0x40, 0x5f, 0x39, 0xc8, 0xac, 0xaf, 0x6b, 0x4a, 0x2f, 0x70,
	0xaf, 0x41, 0x29, 0x1e, 0x8f, 0xa4, 0xe4, 0xa5, 0x4e, 0xf5, 0xeb, 0x4e, 0xf9, 0x96, 0xbd, 0x6e,
	0x79, 0x82, 0x26, 0xf4, 0x49, 0x71, 0x10, 0x8e, 0x59, 0xa3, 0xb4, 0x66, 0xad, 0xdb, 0x9e, 0xde,
	0xb9, 0x0d, 0xa8, 0x92, 0x84, 0x85, 0x11, 0x8d, 0x1b, 0x65, 0x79, 0x60, 0xb6, 0xe2, 0x84, 0x87,
	0x23, 0x42, 0xc7, 0xbc, 0x51, 0x59, 0xb3, 0xd6, 0x4b, 0x9e, 0xd9, 0xba, 0x77, 0xa1, 0xba, 0x1f,
	0x46, 0x9c, 0xa4, 0xac, 0xe1, 0xcc, 0x1a, 0xf1, 0x58, 0x1e, 0x19, 0x23, 0x0c, 0x12, 0x7d, 0x06,
	0x35, 0x8f, 0xb0, 0x84, 0xc6, 0x8c, 0x7c, 0x93, 0x19, 0xf7, 0xa0, 0x9a, 0x12, 0x36, 0x8e, 0x38,
	0x6b, 0xd8, 0xf2, 0x15, 0xae, 0x17, 0xe5, 0x2b, 0x3f, 0xb7, 0x36, 0x43, 0xc6, 0x71, 0xec, 0x13,
	0xcf, 0x60, 0xd1, 0x63, 0xa8, 0x9b, 0x1b, 0x98, 0x7b, 0x1f, 0xea, 0xa9, 0xd9, 0xe8, 0xb7, 0xbc,
	0x3e, 0x37, 0xca, 0x14, 0xc6, 0xcb, 0xd1, 0xe8, 0x0f, 0x16, 0x2c, 0x6f, 0xf3, 0x94, 0xe0, 0x51,
	0xa6, 0xf0, 0x7d, 0x11, 0x18, 0x6a, 0x2d, 0xd5, 0x7d, 0xb5, 0xb0, 0x27, 0x97, 0xbc, 0x0c, 0xee,
	0xbe, 0x0f, 0x8e, 0xca, 0x5d, 0xfd, 0xe0, 0x6e, 0x4b, 0x65, 0x75, 0x2b, 0x4d, 0xfc, 0xd6, 0xb6,
	0x3c, 0x79, 0x72, 0xc9, 0xd3, 0x98, 0x4e, 0x1d, 0xaa, 0x5a, 0x6e, 0xf3, 0x53, 0x70, 0x94, 0x2b,
	0xd1, 0x8f, 0xc0, 0xf9, 0x15, 0x4e, 0x87, 0x84, 0xbb, 0x2e, 0x94, 0x0f, 0x28, 0xe3, 0xda, 0x65,
	0x72, 0x2d, 0x68, 0x09, 0x4d, 0xb9, 0x7a, 0x75, 0x4f, 0xae, 0x51, 0x33, 0x8b, 0x18, 0xf1, 0x8a,
	0x92, 0x57, 0x79, 0xa1, 0xee, 0x99, 0x6d, 0xf3, 0xaf, 0x36, 0x38, 0xbd, 0x98, 0x91, 0x94, 0xa3,
	0x69, 0x5e, 0x01, 0x7e, 0x56, 0xa8, 0x00, 0x33, 0x4f, 0xab, 0x5d, 0xbf, 0x2b, 0x01, 0xff, 0x6f,
	0x71, 0x50, 0x77, 0x5d, 0xb0, 0x38, 0x68, 0x21, 0xb3, 0x09, 0xc8, 0x33, 0x8b, 0xef, 0xc1, 0x0a,
	0x3b, 0x0c, 0x93, 0x01, 0xe3, 0x69, 0xe8, 0xf3, 0x01, 0x39, 0x0a, 0x19, 0x1f, 0xf8, 0x07, 0xc4,
	0x3f, 0x94, 0x26, 0xd5, 0xbc, 0x2b, 0xe2, 0x78, 0x5b, 0x9e, 0x3e, 0x12, 0x87, 0x5d, 0x71, 0x56,
	0x0c, 0x6a, 0xfb, 0xac, 0x41, 0x2d, 0x7d, 0xb8, 0x93, 0x88, 0xb2, 0xfa, 0xdd, 0xf8, 0x50, 0xdd,
	0x75, 0x41, 0x1f, 0x6a, 0x21, 0xaf, 0x97, 0x0f, 0xbf, 0xbb, 0x38, 0x54, 0x77, 0x5d, 0xd8, 0x87,
	0xaf, 0x55, 0x1c, 0xfe, 0xd7, 0x02, 0xc7, 0x23, 0x23, 0x3a, 0x21, 0xc8, 0xcf, 0x7d, 0x78, 0x33,
	0xeb, 0x91, 0x0b, 0x1b, 0x6f, 0xce, 0xf1, 0x5f, 0x6f, 0xf3, 0x9b, 0x5b, 0xa7, 0x92, 0x7a, 0x41,
	0x6f, 0x69, 0x21, 0xb3, 0xde, 0xfa, 0xf9, 0x05, 0xbd, 0xd5, 0xfc, 0xa7, 0x05, 0xe5, 0x67, 0x84,
	0x63, 0xb4, 0x02, 0xa5, 0xa7, 0x64, 0xea, 0x5e, 0x86, 0xd2, 0x21, 0x99, 0xea, 0xfa, 0x28, 0x96,
	0x08, 0x41, 0xf9, 0x29, 0x99, 0x32, 0x51, 0x26, 0x0f, 0xc9, 0xd4, 0x54, 0x41, 0xb9, 0x16, 0x4c,
	0xbb, 0x38, 0x12, 0x4c, 0x13, 0x1c, 0x19, 0xa6, 0x09, 0x8e, 0x04, 0xd3, 0x2e, 0x8e, 0x24, 0xd3,
	0x04, 0x47, 0x19, 0x93, 0x58, 0xa3, 0xf7, 0xc1, 0x79, 0x4a, 0xa6, 0x9a, 0xef, 0xf8, 0x65, 0x46,
	0x92, 0x9d, 0x4b, 0xfa, 0x31, 0x54, 0x15, 0x9a, 0xb9, 0xef, 0x41, 0xe9, 0x70, 0x62, 0xfc, 0xb3,
	0x52, 0xf4, 0x8f, 0xb0, 0xa0, 0xa5, 0x60, 0x9e, 0xc0, 0x34, 0xff, 0x5e, 0x05, 0x47, 0x3d, 0x12,
	0xfa, 0x00, 0x6a, 0xa6, 0xd5, 0xcd, 0x0c, 0x3d, 0x08, 0x6a, 0x81, 0x3e, 0x93, 0x77, 0xda, 0x5e,
	0xb6, 0x2f, 0x74, 0xb1, 0x8c, 0xfd, 0x7e, 0x01, 0x3e, 0xa7, 0x8b, 0x9d, 0x68, 0xac, 0xa2, 0x8b,
	0x19, 0xf8, 0xb9, 0xbb, 0x18, 0xba, 0x01, 0x76, 0x6f, 0xd3, 0x5d, 0xc9, 0x15, 0x97, 0x73, 0x49,
	0x6a, 0x5f, 0xb6, 0x84, 0x05, 0xe2, 0x05, 0x7a, 0x9b, 0x4c, 0xf8, 0x2d, 0x0c, 0x8c, 0x9b, 0xc5,
	0x12, 0x75, 0xc1, 0x51, 0xe9, 0x7d, 0x2a, 0x6f, 0x61, 0x58, 0xb5, 0xe7, 0x0f, 0xab, 0xe8, 0x01,
	0x54, 0x95, 0x10, 0x26, 0xd2, 0x4a, 0x11, 0xcd, 0x03, 0x9c, 0x5e, 0x50, 0x3c, 0x83, 0x44, 0xbf,
	0x87, 0x45, 0xe5, 0x42, 0xad, 0xca, 0xdd, 0x33, 0x17, 0x25, 0xe1, 0x0c, 0x5d, 0x8c, 0xce, 0xed,
	0xba, 0x9f, 0x40, 0xb9, 0x13, 0xd1, 0xbd, 0xd3, 0x1d, 0x90, 0xcf, 0x98, 0x76, 0x71, 0xc6, 0x44,
	0x1c, 0x40, 0xa9, 0x2d, 0xd9, 0x6f, 0x43, 0x79, 0x2f, 0xa2, 0x7b, 0x5a, 0xe5, 0x95, 0x39, 0x2a,
	0x0b, 0xd8, 0x93, 0x4b, 0x9e, 0x84, 0x9d, 0x5f, 0xdd, 0x4d, 0xa8, 0xfd, 0x92, 0xfa, 0x98, 0x87,
	0x34, 0x16, 0x79, 0x13, 0xe3, 0x11, 0x31, 0x73, 0x8a, 0x58, 0x0b, 0xda, 0x78, 0x1c, 0x06, 0x3a,
	0x39, 0xe4, 0x5a, 0xbe, 0x7b, 0x22, 0x46, 0x52, 0xf5, 0xee, 0x49, 0x71, 0xf8, 0xca, 0x84, 0xdd,
	0x87, 0x5a, 0xa4, 0xd7, 0xaf, 0x08, 0x5b, 0x03, 0x17, 0x61, 0x6b, 0xe0, 0xe7, 0x37, 0xe6, 0x31,
	0xd4, 0x8d, 0x40, 0x39, 0x4b, 0x1a, 0x89, 0x73, 0x67, 0xc9, 0x13, 0x1a, 0x78, 0x39, 0xba, 0xf9,
	0x08, 0xaa, 0x5d, 0x1a, 0xf3, 0x94, 0x46, 0xe8, 0x43, 0x70, 0xbb, 0x29, 0xc1, 0x9c, 0xf4, 0xe2,
	0x80, 0x1c, 0x99, 0xca, 0xf9, 0x0e, 0xd4, 0x13, 0x4a, 0xa3, 0x01, 0x0b, 0x3f, 0x57, 0xee, 0xca,
	0x06, 0xf7, 0x4b, 0x5e, 0x4d, 0x9c, 0x6c, 0x87, 0x9f, 0x93, 0xe6, 0x9f, 0x6d, 0x58, 0xf0, 0x48,
	0x12, 0x85, 0x4a, 0x2e, 0xba, 0x23, 0x86, 0x69, 0x9f, 0x4e, 0x48, 0x3a, 0x75, 0x6f, 0xc2, 0x72,
	0x40, 0x22, 0xc2, 0x49, 0x30, 0xc0, 0x43, 0x12, 0x67, 0x83, 0xde, 0x92, 0xa6, 0x3e, 0x94, 0x44,
	0x84, 0xc5, 0x74, 0xbc, 0x87, 0x23, 0x99, 0xce, 0xb7, 0xe0, 0x8d, 0x83, 0x70, 0x78, 0x30, 0x18,
	0x33, 0x3c, 0x24, 0xc7, 0xd9, 0xbe, 0x27, 0x0e, 0x76, 0x04, 0x5d, 0x31, 0xba, 0xeb, 0x70, 0x39,
	0xa2, 0xcf, 0x8f, 0x43, 0x6d, 0x09, 0x5d, 0x8e, 0xe8, 0xf3, 0x02, 0x52, 0xf4, 0x3e, 0xcd, 0x73,
	0x15, 0x9c, 0x63, 0x42, 0xf5, 0x4e, 0xe8, 0x9a, 0xca, 0x5e, 0x10, 0x1c, 0x97, 0xb4, 0xa4, 0xa9,
	0x9a, 0xfd, 0x87, 0xf0, 0x46, 0x6a, 0xac, 0x8d, 0x87, 0x0a, 0xaa, 0x83, 0xe4, 0x72, 0xe1, 0x40,
	0xa2, 0x9b, 0xfb, 0x00, 0x9b, 0x21, 0x93, 0xce, 0x20, 0x29, 0xfa, 0x75, 0xde, 0xfe, 0xae, 0x17,
	0x83, 0x30, 0xcf, 0x1c, 0x15, 0x8d, 0x6f, 0x41, 0x5d, 0xfc, 0xb2, 0x04, 0xeb, 0xda, 0x59, 0xf7,
	0x72, 0x82, 0x8c, 0x5f, 0x1a, 0x10, 0xf9, 0x5f, 0x49, 0xc4, 0x2f, 0x0d, 0x48, 0xf3, 0xcb, 0x0a,
	0x38, 0x1d, 0xec, 0x1f, 0x8e, 0x13, 0xb4, 0x03, 0xf5, 0x5f, 0x10, 0xae, 0x32, 0x1d, 0xbd, 0x7b,
	0xec, 0x46, 0x19, 0xe2, 0x27, 0x6f, 0x14, 0x44, 0xb4, 0x06, 0x95, 0x8f, 0x9f, 0xc7, 0x44, 0x15,
	0xb4, 0x64, 0x36, 0x9f, 0x13, 0x74, 0xb7, 0x10, 0x74, 0x67, 0x16, 0xfb, 0x59, 0xd6, 0xf9, 0xcf,
	0xca, 0xd1, 0x82, 0x45, 0x8d, 0x93, 0x3d, 0xdc, 0x5d, 0x85, 0x8a, 0xa0, 0xeb, 0xb7, 0xca, 0xca,
	0xa8, 0xe5, 0x29, 0x32, 0xfa, 0xa3, 0x05, 0x76, 0xaf, 0x8f, 0xb6, 0x44, 0xcc, 0x0d, 0x43, 0x26,
	0xfe, 0x91, 0x74, 0xce, 0x76, 0x95, 0x8b, 0x54, 0x7e, 0xdb, 0x27, 0x04, 0xcb, 0x4c, 0x6f, 0x67,
	0x8a, 0xdf, 0xcc, 0xa5, 0x69, 0x06, 0x6b, 0x1e, 0x43, 0x3f, 0x6b, 0x09, 0x6e, 0xf1, 0x4e, 0x7d,
	0x55, 0xde, 0x0d, 0x4a, 0xa7, 0x7c, 0xba, 0xd0, 0xc5, 0xa6, 0x9c, 0x17, 0x9b, 0x33, 0xf7, 0x07,
	0xf5, 0xee, 0x33, 0xfd, 0xe1, 0x0b, 0x0b, 0xa0, 0x4b, 0x47, 0x49, 0x4a, 0x18, 0x23, 0xa2, 0x68,
	0xbc, 0x4a, 0xc1, 0xab, 0x05, 0x05, 0x65, 0xb5, 0x3e, 0x55, 0xad, 0x5e, 0xae, 0xd6, 0x83, 0x93,
	0x6a, 0xbd, 0x33, 0x47, 0xad, 0x5c, 0x85, 0x93, 0x1a, 0x36, 0x5f, 0x54, 0xa1, 0xdc, 0x8b, 0xf7,
	0x29, 0xfa, 0xd2, 0x82, 0x8a, 0x2c, 0x3c, 0xe8, 0x37, 0x50, 0xe9, 0xd2, 0x71, 0x2c, 0x3f, 0x51,
	0x30, 0x4e, 0x53, 0xa2, 0xd4, 0x5c, 0xf2, 0xf4, 0xce, 0x5d, 0x83, 0x85, 0x71, 0xec, 0xd3, 0xd1,
	0x28, 0xe4, 0x9c, 0x04, 0xfa, 0x7f, 0x65, 0x91, 0x24, 0xe6, 0x8e, 0x50, 0xc8, 0x0a, 0xe3, 0xa1,
	0x34, 0xa6, 0xe6, 0x65, 0x7b, 0xf4, 0x11, 0x94, 0x77, 0x76, 0x7a, 0x9b, 0xe8, 0xfb, 0x50, 0xef,
	0x66, 0x0c, 0x73, 0xfc, 0x81, 0x7e, 0x00, 0x0b, 0x3b, 0x05, 0x99, 0xf3, 0x20, 0xff, 0xb1, 0xa0,
	0xd4, 0xa7, 0x81, 0x7b, 0x0d, 0x6a, 0x38, 0x49, 0x06, 0x85, 0x96, 0x52, 0xc5, 0x49, 0xb2, 0xa5,
	0xbb, 0x8a, 0x24, 0xdb, 0x85, 0x4e, 0x73, 0x2c, 0xb7, 0x4b, 0x27, 0x73, 0x5b, 0xa5, 0x5f, 0xb9,
	0x10, 0xae, 0x47, 0x32, 0xfd, 0xdc, 0x77, 0xa1, 0xe4, 0x27, 0x63, 0xf9, 0xb1, 0x63, 0x61, 0xe3,
	0xca, 0xf1, 0x7f, 0x9d, 0xfb, 0xb4, 0xd5, 0xed, 0xef, 0x78, 0x02, 0xe0, 0xb6, 0xc1, 0x19, 0x91,
	0x11, 0x4d, 0xa7, 0xfa, 0xeb, 0xc7, 0xca, 0x0c, 0xf4, 0x99, 0x3c, 0xf6, 0x34, 0xcc, 0x7d, 0x4f,
	0x57, 0x93, 0xea, 0xec, 0x24, 0x2e, 0xe1, 0x5b, 0x34, 0x20, 0xaa, 0xc8, 0xa0, 0x7f, 0x5b, 0x50,
	0x16, 0xdb, 0xb9, 0x1d, 0xf4, 0x6d, 0x58, 0x0a, 0x63, 0x4e, 0xd2, 0x18, 0x47, 0x03, 0x1c, 0x04,
	0xa9, 0x36, 0x7a, 0xd1, 0x10, 0x1f, 0x06, 0x41, 0x2a, 0x40, 0xe4, 0xa8, 0x08, 0x52, 0x0e, 0x58,
	0x34, 0x44, 0x09, 0xd2, 0xa6, 0x96, 0xcf, 0x6e, 0x6a, 0xe5, 0xcc, 0xa6, 0xf6, 0x69, 0x60, 0xbe,
	0x0b, 0xcd, 0x9a, 0x2a, 0x0e, 0x3d, 0x09, 0x41, 0x4f, 0xa1, 0xd4, 0xed, 0xef, 0xb8, 0x57, 0xa0,
	0x12, 0x85, 0xa3, 0x50, 0x7d, 0xd3, 0xb0, 0x3c, 0xb5, 0x71, 0x1b, 0x50, 0xd5, 0x7f, 0x12, 0xa4,
	0x91, 0x96, 0x67, 0xb6, 0x02, 0x2f, 0x5b, 0x91, 0xb4, 0x4b, 0xd4, 0x28, 0xb1, 0x41, 0x5b, 0xe0,
	0x28, 0x4d, 0xbe, 0x25, 0x79, 0x1f, 0x2a, 0x3b, 0xdc, 0x0d, 0x28, 0x27, 0x34, 0x30, 0xc9, 0x77,
	0x65, 0x9e, 0x3d, 0x85, 0x32, 0x25, 0xb1, 0xe8, 0x01, 0x54, 0xc4, 0x13, 0x32, 0xf7, 0x1e, 0x54,
	0xc4, 0xa3, 0x1a, 0xee, 0xf9, 0x0f, 0x5f, 0xac, 0xb7, 0x12, 0x8d, 0xde, 0x84, 0x52, 0xaf, 0xcf,
	0xe4, 0xb0, 0x9f, 0xe8, 0xfe, 0x69, 0x87, 0x49, 0xb3, 0x0a, 0x95, 0x47, 0xa3, 0x84, 0x4f, 0x3b,
	0xf1, 0x8b, 0x97, 0xab, 0xd6, 0x3f, 0x5e, 0xae, 0x5a, 0xff, 0x7a, 0xb9, 0x6a, 0xfd, 0xe5, 0xab,
	0x55, 0xeb, 0xc5, 0x57, 0xab, 0x16, 0xdc, 0xa0, 0xe9, 0xb0, 0x35, 0x09, 0x30, 0x66, 0xad, 0x09,
	0x8e, 0x82, 0x16, 0x4e, 0x42, 0x71, 0x8f, 0xbe, 0xb2, 0xb3, 0xb0, 0x8b, 0xa3, 0xa0, 0xaf, 0x36,
	0x7d, 0xeb, 0x93, 0x5b, 0x85, 0x2f, 0xc8, 0x92, 0xa9, 0x2d, 0x98, 0xda, 0xea, 0xab, 0x73, 0x9a,
	0xf8, 0x85, 0x4f, 0xdd, 0x7b, 0x8e, 0xfc, 0x62, 0x7c, 0xf7, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff,
	0xf6, 0x9e, 0x06, 0x3c, 0x0a, 0x17, 0x00, 0x00,
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

func (m *Search_MultiRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Search_MultiRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search_MultiRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Requests) > 0 {
		for iNdEx := len(m.Requests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Requests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *Search_MultiIDRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Search_MultiIDRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search_MultiIDRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Requests) > 0 {
		for iNdEx := len(m.Requests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Requests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if m.Filters != nil {
		{
			size, err := m.Filters.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.Timeout != 0 {
		i = encodeVarintPayload(dAtA, i, uint64(m.Timeout))
		i--
		dAtA[i] = 0x28
	}
	if m.Epsilon != 0 {
		i -= 4
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.Epsilon))))
		i--
		dAtA[i] = 0x25
	}
	if m.Radius != 0 {
		i -= 4
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.Radius))))
		i--
		dAtA[i] = 0x1d
	}
	if m.Num != 0 {
		i = encodeVarintPayload(dAtA, i, uint64(m.Num))
		i--
		dAtA[i] = 0x10
	}
	if len(m.RequestId) > 0 {
		i -= len(m.RequestId)
		copy(dAtA[i:], m.RequestId)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.RequestId)))
		i--
		dAtA[i] = 0xa
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
			dAtA[i] = 0x12
		}
	}
	if len(m.RequestId) > 0 {
		i -= len(m.RequestId)
		copy(dAtA[i:], m.RequestId)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.RequestId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Search_Responses) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Search_Responses) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search_Responses) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Responses) > 0 {
		for iNdEx := len(m.Responses) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Responses[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *Search_StreamResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Search_StreamResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search_StreamResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Payload != nil {
		{
			size := m.Payload.Size()
			i -= size
			if _, err := m.Payload.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *Search_StreamResponse_Response) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search_StreamResponse_Response) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Response != nil {
		{
			size, err := m.Response.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *Search_StreamResponse_Status) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Search_StreamResponse_Status) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Status != nil {
		{
			size, err := m.Status.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func (m *Filter) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Filter) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Filter) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *Filter_Target) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Filter_Target) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Filter_Target) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Port != 0 {
		i = encodeVarintPayload(dAtA, i, uint64(m.Port))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Host) > 0 {
		i -= len(m.Host)
		copy(dAtA[i:], m.Host)
		i = encodeVarintPayload(dAtA, i, uint64(len(m.Host)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Filter_Config) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Filter_Config) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Filter_Config) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Targets) > 0 {
		for iNdEx := len(m.Targets) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Targets[iNdEx])
			copy(dAtA[i:], m.Targets[iNdEx])
			i = encodeVarintPayload(dAtA, i, uint64(len(m.Targets[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Insert) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Insert) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Insert) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *Insert_Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Insert_Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Insert_Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if m.Vector != nil {
		{
			size, err := m.Vector.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Insert_MultiRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Insert_MultiRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Insert_MultiRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Requests) > 0 {
		for iNdEx := len(m.Requests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Requests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *Insert_Config) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Insert_Config) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Insert_Config) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Filters != nil {
		{
			size, err := m.Filters.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.SkipStrictExistCheck {
		i--
		if m.SkipStrictExistCheck {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Update) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Update) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Update) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *Update_Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Update_Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Update_Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if m.Vector != nil {
		{
			size, err := m.Vector.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Update_MultiRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Update_MultiRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Update_MultiRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Requests) > 0 {
		for iNdEx := len(m.Requests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Requests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *Update_Config) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Update_Config) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Update_Config) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Filters != nil {
		{
			size, err := m.Filters.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.SkipStrictExistCheck {
		i--
		if m.SkipStrictExistCheck {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Upsert) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Upsert) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Upsert) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *Upsert_Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Upsert_Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Upsert_Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if m.Vector != nil {
		{
			size, err := m.Vector.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Upsert_MultiRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Upsert_MultiRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Upsert_MultiRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Requests) > 0 {
		for iNdEx := len(m.Requests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Requests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *Upsert_Config) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Upsert_Config) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Upsert_Config) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Filters != nil {
		{
			size, err := m.Filters.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.SkipStrictExistCheck {
		i--
		if m.SkipStrictExistCheck {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Remove) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Remove) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Remove) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *Remove_Request) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Remove_Request) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Remove_Request) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if m.Id != nil {
		{
			size, err := m.Id.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Remove_MultiRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Remove_MultiRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Remove_MultiRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Requests) > 0 {
		for iNdEx := len(m.Requests) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Requests[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *Remove_Config) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Remove_Config) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Remove_Config) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.SkipStrictExistCheck {
		i--
		if m.SkipStrictExistCheck {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
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

func (m *Object_StreamDistance) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_StreamDistance) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamDistance) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Payload != nil {
		{
			size := m.Payload.Size()
			i -= size
			if _, err := m.Payload.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *Object_StreamDistance_Distance) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamDistance_Distance) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Distance != nil {
		{
			size, err := m.Distance.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *Object_StreamDistance_Status) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamDistance_Status) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Status != nil {
		{
			size, err := m.Status.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
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
			f21 := math.Float32bits(float32(m.Vector[iNdEx]))
			i -= 4
			encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(f21))
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

func (m *Object_StreamVector) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_StreamVector) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamVector) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Payload != nil {
		{
			size := m.Payload.Size()
			i -= size
			if _, err := m.Payload.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *Object_StreamVector_Vector) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamVector_Vector) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Vector != nil {
		{
			size, err := m.Vector.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *Object_StreamVector_Status) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamVector_Status) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Status != nil {
		{
			size, err := m.Status.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
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

func (m *Object_StreamBlob) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_StreamBlob) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamBlob) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Payload != nil {
		{
			size := m.Payload.Size()
			i -= size
			if _, err := m.Payload.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *Object_StreamBlob_Blob) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamBlob_Blob) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Blob != nil {
		{
			size, err := m.Blob.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *Object_StreamBlob_Status) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamBlob_Status) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Status != nil {
		{
			size, err := m.Status.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
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

func (m *Object_StreamLocation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Object_StreamLocation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamLocation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Payload != nil {
		{
			size := m.Payload.Size()
			i -= size
			if _, err := m.Payload.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *Object_StreamLocation_Location) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamLocation_Location) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Location != nil {
		{
			size, err := m.Location.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *Object_StreamLocation_Status) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Object_StreamLocation_Status) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Status != nil {
		{
			size, err := m.Status.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintPayload(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
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

func (m *Backup_Vector) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Vector) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Vector) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
			f28 := math.Float32bits(float32(m.Vector[iNdEx]))
			i -= 4
			encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(f28))
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

func (m *Backup_Vectors) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Vectors) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Vectors) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *Backup_Compressed_Vector) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Compressed_Vector) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Compressed_Vector) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *Backup_Compressed_Vectors) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Backup_Compressed_Vectors) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Backup_Compressed_Vectors) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *Search_MultiRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Requests) > 0 {
		for _, e := range m.Requests {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
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

func (m *Search_MultiIDRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Requests) > 0 {
		for _, e := range m.Requests {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
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
	l = len(m.RequestId)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
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
	if m.Filters != nil {
		l = m.Filters.Size()
		n += 1 + l + sovPayload(uint64(l))
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
	l = len(m.RequestId)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
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

func (m *Search_Responses) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Responses) > 0 {
		for _, e := range m.Responses {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Search_StreamResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Payload != nil {
		n += m.Payload.Size()
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Search_StreamResponse_Response) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Response != nil {
		l = m.Response.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	return n
}
func (m *Search_StreamResponse_Status) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != nil {
		l = m.Status.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	return n
}
func (m *Filter) Size() (n int) {
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

func (m *Filter_Target) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Host)
	if l > 0 {
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.Port != 0 {
		n += 1 + sovPayload(uint64(m.Port))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Filter_Config) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Targets) > 0 {
		for _, s := range m.Targets {
			l = len(s)
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Insert) Size() (n int) {
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

func (m *Insert_Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Vector != nil {
		l = m.Vector.Size()
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

func (m *Insert_MultiRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Requests) > 0 {
		for _, e := range m.Requests {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Insert_Config) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SkipStrictExistCheck {
		n += 2
	}
	if m.Filters != nil {
		l = m.Filters.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Update) Size() (n int) {
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

func (m *Update_Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Vector != nil {
		l = m.Vector.Size()
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

func (m *Update_MultiRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Requests) > 0 {
		for _, e := range m.Requests {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Update_Config) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SkipStrictExistCheck {
		n += 2
	}
	if m.Filters != nil {
		l = m.Filters.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Upsert) Size() (n int) {
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

func (m *Upsert_Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Vector != nil {
		l = m.Vector.Size()
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

func (m *Upsert_MultiRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Requests) > 0 {
		for _, e := range m.Requests {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Upsert_Config) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SkipStrictExistCheck {
		n += 2
	}
	if m.Filters != nil {
		l = m.Filters.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Remove) Size() (n int) {
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

func (m *Remove_Request) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != nil {
		l = m.Id.Size()
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

func (m *Remove_MultiRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Requests) > 0 {
		for _, e := range m.Requests {
			l = e.Size()
			n += 1 + l + sovPayload(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Remove_Config) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SkipStrictExistCheck {
		n += 2
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

func (m *Object_StreamDistance) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Payload != nil {
		n += m.Payload.Size()
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_StreamDistance_Distance) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Distance != nil {
		l = m.Distance.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	return n
}
func (m *Object_StreamDistance_Status) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != nil {
		l = m.Status.Size()
		n += 1 + l + sovPayload(uint64(l))
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

func (m *Object_StreamVector) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Payload != nil {
		n += m.Payload.Size()
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_StreamVector_Vector) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Vector != nil {
		l = m.Vector.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	return n
}
func (m *Object_StreamVector_Status) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != nil {
		l = m.Status.Size()
		n += 1 + l + sovPayload(uint64(l))
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

func (m *Object_StreamBlob) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Payload != nil {
		n += m.Payload.Size()
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_StreamBlob_Blob) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Blob != nil {
		l = m.Blob.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	return n
}
func (m *Object_StreamBlob_Status) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != nil {
		l = m.Status.Size()
		n += 1 + l + sovPayload(uint64(l))
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

func (m *Object_StreamLocation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Payload != nil {
		n += m.Payload.Size()
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Object_StreamLocation_Location) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Location != nil {
		l = m.Location.Size()
		n += 1 + l + sovPayload(uint64(l))
	}
	return n
}
func (m *Object_StreamLocation_Status) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != nil {
		l = m.Status.Size()
		n += 1 + l + sovPayload(uint64(l))
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

func (m *Backup_Vector) Size() (n int) {
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

func (m *Backup_Vectors) Size() (n int) {
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

func (m *Backup_Compressed_Vector) Size() (n int) {
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

func (m *Backup_Compressed_Vectors) Size() (n int) {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Search_MultiRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MultiRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requests", wireType)
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
			m.Requests = append(m.Requests, &Search_Request{})
			if err := m.Requests[len(m.Requests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Search_MultiIDRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MultiIDRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiIDRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requests", wireType)
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
			m.Requests = append(m.Requests, &Search_IDRequest{})
			if err := m.Requests[len(m.Requests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestId", wireType)
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
			m.RequestId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
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
		case 3:
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
		case 4:
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
		case 5:
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
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Filters", wireType)
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
			if m.Filters == nil {
				m.Filters = &Filter_Config{}
			}
			if err := m.Filters.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
				return fmt.Errorf("proto: wrong wireType = %d for field RequestId", wireType)
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
			m.RequestId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Search_Responses) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Responses: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Responses: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Responses", wireType)
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
			m.Responses = append(m.Responses, &Search_Response{})
			if err := m.Responses[len(m.Responses)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Search_StreamResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: StreamResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StreamResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Response", wireType)
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
			v := &Search_Response{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Payload = &Search_StreamResponse_Response{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
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
			v := &rpc.Status{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Payload = &Search_StreamResponse_Status{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Filter) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Filter: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Filter: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Filter_Target) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Target: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Target: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Host", wireType)
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
			m.Host = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Port", wireType)
			}
			m.Port = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPayload
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Port |= uint32(b&0x7F) << shift
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Filter_Config) Unmarshal(dAtA []byte) error {
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Targets", wireType)
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
			m.Targets = append(m.Targets, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Insert) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Insert: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Insert: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Insert_Request) Unmarshal(dAtA []byte) error {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Vector", wireType)
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
			if m.Vector == nil {
				m.Vector = &Object_Vector{}
			}
			if err := m.Vector.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
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
				m.Config = &Insert_Config{}
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Insert_MultiRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MultiRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requests", wireType)
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
			m.Requests = append(m.Requests, &Insert_Request{})
			if err := m.Requests[len(m.Requests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Insert_Config) Unmarshal(dAtA []byte) error {
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
				return fmt.Errorf("proto: wrong wireType = %d for field SkipStrictExistCheck", wireType)
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
			m.SkipStrictExistCheck = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Filters", wireType)
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
			if m.Filters == nil {
				m.Filters = &Filter_Config{}
			}
			if err := m.Filters.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Update) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Update: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Update: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Update_Request) Unmarshal(dAtA []byte) error {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Vector", wireType)
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
			if m.Vector == nil {
				m.Vector = &Object_Vector{}
			}
			if err := m.Vector.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
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
				m.Config = &Update_Config{}
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Update_MultiRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MultiRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requests", wireType)
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
			m.Requests = append(m.Requests, &Update_Request{})
			if err := m.Requests[len(m.Requests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Update_Config) Unmarshal(dAtA []byte) error {
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
				return fmt.Errorf("proto: wrong wireType = %d for field SkipStrictExistCheck", wireType)
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
			m.SkipStrictExistCheck = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Filters", wireType)
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
			if m.Filters == nil {
				m.Filters = &Filter_Config{}
			}
			if err := m.Filters.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Upsert) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Upsert: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Upsert: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Upsert_Request) Unmarshal(dAtA []byte) error {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Vector", wireType)
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
			if m.Vector == nil {
				m.Vector = &Object_Vector{}
			}
			if err := m.Vector.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
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
				m.Config = &Upsert_Config{}
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Upsert_MultiRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MultiRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requests", wireType)
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
			m.Requests = append(m.Requests, &Upsert_Request{})
			if err := m.Requests[len(m.Requests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Upsert_Config) Unmarshal(dAtA []byte) error {
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
				return fmt.Errorf("proto: wrong wireType = %d for field SkipStrictExistCheck", wireType)
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
			m.SkipStrictExistCheck = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Filters", wireType)
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
			if m.Filters == nil {
				m.Filters = &Filter_Config{}
			}
			if err := m.Filters.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Remove) Unmarshal(dAtA []byte) error {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Remove_Request) Unmarshal(dAtA []byte) error {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
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
			if m.Id == nil {
				m.Id = &Object_ID{}
			}
			if err := m.Id.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
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
				m.Config = &Remove_Config{}
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Remove_MultiRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MultiRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requests", wireType)
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
			m.Requests = append(m.Requests, &Remove_Request{})
			if err := m.Requests[len(m.Requests)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Remove_Config) Unmarshal(dAtA []byte) error {
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
				return fmt.Errorf("proto: wrong wireType = %d for field SkipStrictExistCheck", wireType)
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
			m.SkipStrictExistCheck = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Object_StreamDistance) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: StreamDistance: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StreamDistance: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Distance", wireType)
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
			v := &Object_Distance{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Payload = &Object_StreamDistance_Distance{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
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
			v := &rpc.Status{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Payload = &Object_StreamDistance_Status{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Object_StreamVector) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: StreamVector: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StreamVector: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vector", wireType)
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
			v := &Object_Vector{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Payload = &Object_StreamVector_Vector{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
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
			v := &rpc.Status{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Payload = &Object_StreamVector_Status{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Object_StreamBlob) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: StreamBlob: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StreamBlob: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Blob", wireType)
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
			v := &Object_Blob{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Payload = &Object_StreamBlob_Blob{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
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
			v := &rpc.Status{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Payload = &Object_StreamBlob_Status{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Object_StreamLocation) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: StreamLocation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StreamLocation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Location", wireType)
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
			v := &Object_Location{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Payload = &Object_StreamLocation_Location{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
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
			v := &rpc.Status{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Payload = &Object_StreamLocation_Status{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPayload(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Backup_Vector) Unmarshal(dAtA []byte) error {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Backup_Vectors) Unmarshal(dAtA []byte) error {
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
			m.Vectors = append(m.Vectors, &Backup_Vector{})
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Backup_Compressed_Vector) Unmarshal(dAtA []byte) error {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
func (m *Backup_Compressed_Vectors) Unmarshal(dAtA []byte) error {
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
			m.Vectors = append(m.Vectors, &Backup_Compressed_Vector{})
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
			if (skippy < 0) || (iNdEx+skippy) < 0 {
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
