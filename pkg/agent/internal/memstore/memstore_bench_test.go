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
package memstore

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

const benchObjects = 10000

// benchStores returns a kvs holding benchObjects objects, a vqueue where every
// fourth of them is also queued as a pending insert, and their uuids.
func benchStores(b *testing.B) (kvs.BidiMap, vqueue.Queue, []string) {
	b.Helper()
	kv := kvs.New()
	vq, err := vqueue.New()
	if err != nil {
		b.Fatal(err)
	}
	vec := make([]float32, 128)
	uuids := make([]string, benchObjects)
	for i := range uuids {
		uuids[i] = fmt.Sprintf("uuid-%08d", i)
		kv.Set(uuids[i], uint32(i+1), int64(i+1))
		if i%4 == 0 {
			if err := vq.PushInsert(uuids[i], vec, int64(i+2)); err != nil {
				b.Fatal(err)
			}
		}
	}
	return kv, vq, uuids
}

func BenchmarkExists(b *testing.B) {
	kv, vq, uuids := benchStores(b)
	var n atomic.Uint64
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, ok := Exists(kv, vq, uuids[n.Add(1)%benchObjects]); !ok {
				b.Fatal("object not found")
			}
		}
	})
}

func BenchmarkGetObject(b *testing.B) {
	kv, vq, uuids := benchStores(b)
	vec := make([]float32, 128)
	//nolint:unparam // the signature is GetObject's getVectorFn, whose error this benchmark does not exercise
	getVector := func(uint32) ([]float32, error) { return vec, nil }
	var n atomic.Uint64
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, _, err := GetObject(kv, vq, uuids[n.Add(1)%benchObjects], getVector); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkUpdateTimestamp(b *testing.B) {
	kv, vq, uuids := benchStores(b)
	per := max(benchObjects/runtime.GOMAXPROCS(0), 1)
	var w atomic.Uint64
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		// Each worker owns a disjoint block of uuids and a monotonic timestamp, so the
		// benchmark measures the update path itself and not races between workers on
		// one uuid (which legitimately fail with NotFound / newer-timestamp errors).
		base := int(w.Add(1)-1) * per % benchObjects
		ts := int64(1 << 40)
		for i := 0; pb.Next(); i++ {
			ts++
			if err := UpdateTimestamp(kv, vq, uuids[base+i%per], ts, true); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkUUIDs(b *testing.B) {
	kv, vq, _ := benchStores(b)
	b.ReportAllocs()
	for b.Loop() {
		if got := UUIDs(b.Context(), kv, vq); len(got) != benchObjects {
			b.Fatalf("got %d uuids, want %d", len(got), benchObjects)
		}
	}
}
