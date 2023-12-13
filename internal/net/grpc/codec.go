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
func (Codec) Marshal(obj interface{}) (data []byte, err error) {
	switch v := obj.(type) {
	case vtprotoMessage:
		data, err = v.MarshalVT()
	case proto.Message:
		data, err = proto.Marshal(v)
	default:
		err = errors.ErrInvalidProtoMessageType(v)
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Unmarshal parses the byte stream data into v.
func (Codec) Unmarshal(data []byte, obj interface{}) (err error) {
	switch v := obj.(type) {
	case vtprotoMessage:
		err = v.UnmarshalVT(data)
	case proto.Message:
		err = proto.Unmarshal(data, v)
	default:
		err = errors.ErrInvalidProtoMessageType(v)
	}
	return err
}

func (Codec) Name() string {
	return Name
}
