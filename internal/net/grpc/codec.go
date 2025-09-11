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

package grpc

import (
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/mem"
)

// Codec represents a gRPC codec.
type Codec struct {
	fallback encoding.CodecV2
}

// Name represents the codec name.
const Name = "proto"

var bufferPool = mem.NewTieredBufferPool(
	2<<10,   // 2KB
	4<<10,   // 4KB (go page size)
	8<<10,   // 8KB
	16<<10,  // 16KB (max HTTP/2 frame size used by gRPC)
	32<<10,  // 32KB (default buffer size for io.Copy)
	64<<10,  // 64KB
	128<<10, // 128KB
	256<<10, // 256KB
	512<<10, // 512KB
	1<<20,   // 1MB
	2<<20,   // 2MB
	4<<20,   // 4MB
	8<<20,   // 8MB
	16<<20,  // 16MB
)

type vtprotoMessage interface {
	SizeVT() int
	MarshalToSizedBufferVT([]byte) (int, error)
	UnmarshalVT([]byte) error
}

// Marshal returns byte slice representing the proto message marshalling result.
func (c Codec) Marshal(obj any) (data mem.BufferSlice, err error) {
	switch m := obj.(type) {
	case vtprotoMessage:
		size := m.SizeVT()
		if mem.IsBelowBufferPoolingThreshold(size) { // less than 1KB
			buf := make([]byte, size)
			n, err := m.MarshalToSizedBufferVT(buf[:size])
			if err != nil {
				return nil, err
			}
			return mem.BufferSlice{mem.SliceBuffer(buf[:n])}, nil
		}
		buf := bufferPool.Get(size)
		n, err := m.MarshalToSizedBufferVT((*buf)[:size])
		if err != nil {
			bufferPool.Put(buf)
			return nil, err
		}
		*buf = (*buf)[:n]
		return mem.BufferSlice{mem.NewBuffer(buf, bufferPool)}, nil
	case proto.Message:
		buf, err := proto.Marshal(m)
		if err != nil {
			return nil, err
		}
		return mem.BufferSlice{mem.SliceBuffer(buf)}, nil
	default:
		return c.fallback.Marshal(obj)
	}
}

// Unmarshal parses the byte stream data into v.
func (c Codec) Unmarshal(data mem.BufferSlice, obj any) (err error) {
	switch m := obj.(type) {
	case vtprotoMessage:
		buf := data.MaterializeToBuffer(bufferPool)
		defer buf.Free()
		return m.UnmarshalVT(buf.ReadOnlyData())
	case proto.Message:
		buf := data.MaterializeToBuffer(bufferPool)
		defer buf.Free()
		return proto.Unmarshal(buf.ReadOnlyData(), m)
	default:
		return c.fallback.Unmarshal(data, m)
	}
}

func (Codec) Name() string {
	return Name
}

func init() {
	encoding.RegisterCodecV2(&Codec{
		fallback: encoding.GetCodecV2(Name),
	})
	log.Debug("Vald Customized gRPC CodecV2 Registered")
}
