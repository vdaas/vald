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
	"context"
	"slices"
	"testing"

	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

type lenKV struct {
	kvs.BidiMap
	uuids []string
}

func (k lenKV) Len() uint64 { return uint64(len(k.uuids)) }

func (k lenKV) Range(_ context.Context, f func(string, uint32, int64) bool) {
	for i, uuid := range k.uuids {
		if !f(uuid, uint32(i+1), 1) {
			return
		}
	}
}

type lenVQ struct {
	vqueue.Queue
	ivq, dvq int
}

func (q lenVQ) IVQLen() int { return q.ivq }

func (q lenVQ) DVQLen() int { return q.dvq }

func (lenVQ) Range(context.Context, func(string, []float32, int64) bool) {}

func (lenVQ) DVExists(string) (int64, bool) { return 0, false }

// TestUUIDs_capacityHint pins that a delete queue longer than kvs+insert queue
// (kvs is decremented before the delete entry is popped) does not wrap the
// uint64 capacity hint and panic in make.
func TestUUIDs_capacityHint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		kv   []string
		ivq  int
		dvq  int
	}{
		{name: "more pending deletes than objects", kv: []string{"a"}, ivq: 0, dvq: 2},
		{name: "pending deletes and empty kvs", kv: nil, ivq: 0, dvq: 1},
		{name: "fewer pending deletes than objects", kv: []string{"a", "b"}, ivq: 1, dvq: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got []string
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("UUIDs panicked: %v", r)
					}
				}()
				got = UUIDs(t.Context(), lenKV{uuids: tt.kv}, lenVQ{ivq: tt.ivq, dvq: tt.dvq})
			}()
			if !slices.Equal(got, tt.kv) {
				t.Errorf("UUIDs() = %v, want %v", got, tt.kv)
			}
		})
	}
}
