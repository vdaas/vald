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
	"slices"

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/sync"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/mem"
)

// Codec represents a gRPC codec.
type Codec struct {
	buffer   mem.BufferPool
	fallback encoding.CodecV2
}

// Name represents the codec name.
const Name = "proto"

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
		buf := c.buffer.Get(size)
		n, err := m.MarshalToSizedBufferVT((*buf)[:size])
		if err != nil {
			c.buffer.Put(buf)
			return nil, err
		}
		*buf = (*buf)[:n]
		return mem.BufferSlice{mem.NewBuffer(buf, c.buffer)}, nil
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
		buf := data.MaterializeToBuffer(c.buffer)
		defer buf.Free()
		return m.UnmarshalVT(buf.ReadOnlyData())
	case proto.Message:
		buf := data.MaterializeToBuffer(c.buffer)
		defer buf.Free()
		return proto.Unmarshal(buf.ReadOnlyData(), m)
	default:
		return c.fallback.Unmarshal(data, m)
	}
}

func (Codec) Name() string {
	return Name
}

var codecOnce sync.Once

const maxTieredBufferPoolSize = 256 << 20 // 256MB

// InitCodec initializes and registers a customized gRPC CodecV2 that uses a tiered buffer pool sized according to `size`.
// If called multiple times, initialization runs only once. `size` is capped at 256MB; when `size` is less than or equal to zero the pool uses the default tier
// sizes. The pool's tier list is trimmed to tiers strictly less than `size`, or extended by doubling the largest default tier until reaching `size`. The
// initialized CodecV2 uses this pool for buffer allocations and falls back to the existing codec for types not handled by the custom implementation.
func InitCodec(size int) {
	codecOnce.Do(func() {
		if size > maxTieredBufferPoolSize {
			size = maxTieredBufferPoolSize
		}
		pool := mem.NewTieredBufferPool(func(defaultTiers []int) []int {
			if size <= 0 {
				return defaultTiers
			}

			if idx := slices.IndexFunc(defaultTiers, func(tsize int) bool {
				return tsize > size
			}); 0 <= idx {
				if idx == 0 {
					return []int{min(defaultTiers[0], size)}
				}
				return defaultTiers[:idx]
			}

			// ---- if all defaultTiers are <= size, extend: reuse last, double with shift until reaching size
			last := defaultTiers[len(defaultTiers)-1]
			for last < size {
				last <<= 1
				if last >= size {
					last = size
				}
				defaultTiers = append(defaultTiers, last)
			}
			return defaultTiers
		}([]int{
			256,       // 256B
			512,       // 512B
			1 << 10,   // 1KB
			2 << 10,   // 2KB
			4 << 10,   // 4KB (go page size)
			8 << 10,   // 8KB
			12 << 10,  // 12KB
			16 << 10,  // 16KB (max HTTP/2 frame size used by gRPC)
			24 << 10,  // 24KB
			32 << 10,  // 32KB (default buffer size for io.Copy)
			64 << 10,  // 64KB
			128 << 10, // 128KB
			256 << 10, // 256KB
			512 << 10, // 512KB
			1 << 20,   // 1MB
		})...)

		encoding.RegisterCodecV2(&Codec{
			buffer:   pool,
			fallback: encoding.GetCodecV2(Name),
		})
		log.Debugf("Vald Customized gRPC CodecV2 Registered for Size %d", size)
	})
}
