//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package proto provides proto file logic
package proto

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type (
	// Message is a type alias of proto.Message.
	Message = proto.Message

	// MessageV1 is a type alias of protoiface.MessageV1.
	MessageV1 = protoiface.MessageV1

	// Name is a type alias of protoreflect.Name.
	Name = protoreflect.Name
)

// Marshal returns the wire-format encoding of m.
func Marshal(m Message) ([]byte, error) {
	return proto.Marshal(m)
}

// Unmarshal parses the wire-format message in b and places the result in m.
func Unmarshal(data []byte, v Message) error {
	return proto.Unmarshal(data, v)
}

// Clone returns a deep copy of m.
func Clone(m Message) Message {
	return proto.Clone(m)
}

// ToMessageV1 downcasts Messages to V1 protobuf MessageV1
func ToMessageV1(m Message) MessageV1 {
	return protoimpl.X.ProtoMessageV1Of(m)
}
