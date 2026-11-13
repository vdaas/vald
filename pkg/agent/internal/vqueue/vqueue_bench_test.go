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
package vqueue

import (
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/sync/atomic"
)

const benchEntries = 10000

func benchQueue(b *testing.B) (Queue, []string, []float32) {
	b.Helper()
	vq, err := New()
	if err != nil {
		b.Fatal(err)
	}
	vec := make([]float32, 128)
	uuids := make([]string, benchEntries)
	for i := range uuids {
		uuids[i] = fmt.Sprintf("uuid-%08d", i)
	}
	return vq, uuids, vec
}

func BenchmarkVQueue_PushInsert(b *testing.B) {
	vq, uuids, vec := benchQueue(b)
	var (
		n  atomic.Uint64
		ts atomic.Int64
	)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := vq.PushInsert(uuids[n.Add(1)%benchEntries], vec, ts.Add(1)); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkVQueue_PushDelete(b *testing.B) {
	vq, uuids, _ := benchQueue(b)
	var (
		n  atomic.Uint64
		ts atomic.Int64
	)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := vq.PushDelete(uuids[n.Add(1)%benchEntries], ts.Add(1)); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkVQueue_GetVectorWithTimestamp(b *testing.B) {
	vq, uuids, vec := benchQueue(b)
	for i, u := range uuids {
		if err := vq.PushInsert(u, vec, int64(i+1)); err != nil {
			b.Fatal(err)
		}
	}
	var n atomic.Uint64
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, _, _, ok := vq.GetVectorWithTimestamp(uuids[n.Add(1)%benchEntries]); !ok {
				b.Fatal("not found")
			}
		}
	})
}

func BenchmarkVQueue_IVExists(b *testing.B) {
	vq, uuids, vec := benchQueue(b)
	for i, u := range uuids {
		if err := vq.PushInsert(u, vec, int64(i+1)); err != nil {
			b.Fatal(err)
		}
	}
	var n atomic.Uint64
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, ok := vq.IVExists(uuids[n.Add(1)%benchEntries]); !ok {
				b.Fatal("not found")
			}
		}
	})
}
