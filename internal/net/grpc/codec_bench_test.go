// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package grpc

import (
	"fmt"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/mem"
)

func BenchmarkCodec_Marshal_Unmarshal(b *testing.B) {
	pool := mem.NewTieredBufferPool(
		2<<10,
		4<<10,
		8<<10,
		16<<10,
		32<<10,
		64<<10,
		128<<10,
		256<<10,
	)
	pooledCodec := Codec{
		buffer:   pool,
		fallback: encoding.GetCodecV2(Name),
	}
	unpooledCodec := Codec{
		fallback: encoding.GetCodecV2(Name),
	}

	dims := []int{0, 16, 128, 512, 1536, 4096}
	for _, dim := range dims {
		var msg any
		if dim == 0 {
			msg = &payload.Empty{}
		} else {
			vec := make([]float32, dim)
			for i := range vec {
				vec[i] = float32(i) * 0.1
			}
			msg = &payload.Object_Vector{
				Id:     fmt.Sprintf("vec-%d", dim),
				Vector: vec,
			}
		}

		b.Run(fmt.Sprintf("Pooled_Marshal_%dd", dim), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				bs, err := pooledCodec.Marshal(msg)
				if err != nil {
					b.Fatal(err)
				}
				bs.Free()
			}
		})

		b.Run(fmt.Sprintf("Unpooled_Marshal_%dd", dim), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				bs, err := unpooledCodec.Marshal(msg)
				if err != nil {
					b.Fatal(err)
				}
				bs.Free()
			}
		})

		// Benchmark unmarshal
		bs, err := pooledCodec.Marshal(msg)
		if err != nil {
			b.Fatal(err)
		}
		data := bs.Materialize()
		bs.Free()

		b.Run(fmt.Sprintf("Pooled_Unmarshal_%dd", dim), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				var target any
				if dim == 0 {
					target = &payload.Empty{}
				} else {
					target = &payload.Object_Vector{}
				}
				bufSlice := mem.BufferSlice{mem.SliceBuffer(data)}
				err := pooledCodec.Unmarshal(bufSlice, target)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(fmt.Sprintf("Unpooled_Unmarshal_%dd", dim), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				var target any
				if dim == 0 {
					target = &payload.Empty{}
				} else {
					target = &payload.Object_Vector{}
				}
				bufSlice := mem.BufferSlice{mem.SliceBuffer(data)}
				err := unpooledCodec.Unmarshal(bufSlice, target)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(fmt.Sprintf("Pooled_RoundTrip_%dd", dim), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				bs, err := pooledCodec.Marshal(msg)
				if err != nil {
					b.Fatal(err)
				}
				var target any
				if dim == 0 {
					target = &payload.Empty{}
				} else {
					target = &payload.Object_Vector{}
				}
				err = pooledCodec.Unmarshal(bs, target)
				if err != nil {
					b.Fatal(err)
				}
				bs.Free()
			}
		})

		b.Run(fmt.Sprintf("Unpooled_RoundTrip_%dd", dim), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				bs, err := unpooledCodec.Marshal(msg)
				if err != nil {
					b.Fatal(err)
				}
				var target any
				if dim == 0 {
					target = &payload.Empty{}
				} else {
					target = &payload.Object_Vector{}
				}
				err = unpooledCodec.Unmarshal(bs, target)
				if err != nil {
					b.Fatal(err)
				}
				bs.Free()
			}
		})
	}
}
