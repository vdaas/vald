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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/proto"
)

// Codec represents a gRPC codec.
type Codec struct{}

// Name represents the codec name.
const Name = "proto"

type vtprotoMessage interface {
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

// Marshal returns byte slice representing the proto message marshalling result.
func (Codec) Marshal(obj interface{}) ([]byte, error) {
	switch v := obj.(type) {
	case vtprotoMessage:
		return v.MarshalVT()
	case proto.Message:
		return proto.Marshal(v)
	default:
		return nil, errors.ErrInvalidProtoMessageType(v)
	}
}

// Unmarshal parses the byte stream data into v.
func (Codec) Unmarshal(data []byte, obj interface{}) error {
	switch v := obj.(type) {
	case vtprotoMessage:
		return v.UnmarshalVT(data)
	case proto.Message:
		return proto.Unmarshal(data, v)
	default:
		return errors.ErrInvalidProtoMessageType(v)
	}
}

func (Codec) Name() string {
	return Name
}
