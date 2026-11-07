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
package service

import (
	"slices"
	"testing"
	"time"
)

// Test_faiss_CreateIndex_newestWins pins that CreateIndex applies a queued delete or insert
// only when it is newer than the indexed object. A delete older than the indexed object used
// to remove it, and a stale Update (delete at t, insert at t+1 <= indexed ts) used to repoint
// the uuid to the stale vector. An insert applied over an indexed object must remove the
// replaced object from the add buffer instead of orphaning it there.
//
//nolint:maintidx // table volume is inherent to this table-driven test
func Test_faiss_CreateIndex_newestWins(t *testing.T) {
	t.Parallel()

	const uuid = "uuid"
	oldVec := []float32{1, 2}
	newVec := []float32{3, 4}
	tests := []struct {
		op         func(t *testing.T, f *faiss, base int64)
		wantTS     func(base int64) int64
		name       string
		wantVec    []float32
		wantExists bool
		// wantOldRemoved reports whether the initially buffered object is gone from the buffer.
		wantOldRemoved bool
	}{
		{
			name: "delete older than the indexed object is skipped",
			op: func(t *testing.T, f *faiss, base int64) {
				t.Helper()
				if err := f.DeleteWithTime(uuid, base-10); err != nil {
					t.Fatal(err)
				}
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
		{
			name: "stale update is skipped",
			op: func(t *testing.T, f *faiss, base int64) {
				t.Helper()
				if err := f.UpdateWithTime(uuid, newVec, base-5); err != nil {
					t.Fatal(err)
				}
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
		{
			name: "update whose insert ties the indexed timestamp is skipped",
			op: func(t *testing.T, f *faiss, base int64) {
				t.Helper()
				if err := f.UpdateWithTime(uuid, newVec, base-1); err != nil {
					t.Fatal(err)
				}
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
		{
			name: "delete at the indexed timestamp is applied",
			op: func(t *testing.T, f *faiss, base int64) {
				t.Helper()
				if err := f.DeleteWithTime(uuid, base); err != nil {
					t.Fatal(err)
				}
			},
			wantOldRemoved: true,
		},
		{
			name: "newer update is applied",
			op: func(t *testing.T, f *faiss, base int64) {
				t.Helper()
				if err := f.UpdateWithTime(uuid, newVec, base+10); err != nil {
					t.Fatal(err)
				}
			},
			wantExists:     true,
			wantVec:        newVec,
			wantTS:         func(base int64) int64 { return base + 11 },
			wantOldRemoved: true,
		},
		{
			name: "UpdateTimestamp of a queued update keeps the new vector",
			op: func(t *testing.T, f *faiss, base int64) {
				t.Helper()
				if err := f.UpdateWithTime(uuid, newVec, base+10); err != nil {
					t.Fatal(err)
				}
				if err := f.UpdateTimestamp(uuid, base+20, false); err != nil {
					t.Fatal(err)
				}
			},
			wantExists:     true,
			wantVec:        newVec,
			wantTS:         func(base int64) int64 { return base + 20 },
			wantOldRemoved: true,
		},
		{
			// UpdateTimestamp used to raise the kvs timestamp past the queued delete, so
			// CreateIndex skipped it and kept the superseded buffered vector.
			name: "UpdateTimestamp after a queued update and a newer delete does not resurrect",
			op: func(t *testing.T, f *faiss, base int64) {
				t.Helper()
				if err := f.UpdateWithTime(uuid, newVec, base+10); err != nil {
					t.Fatal(err)
				}
				if err := f.DeleteWithTime(uuid, base+20); err != nil {
					t.Fatal(err)
				}
				if err := f.UpdateTimestamp(uuid, base+30, false); err == nil {
					t.Fatal("UpdateTimestamp of an object pending deletion succeeded")
				}
			},
			wantOldRemoved: true,
		},
		{
			name: "newer insert over the indexed object replaces it",
			op: func(t *testing.T, f *faiss, base int64) {
				t.Helper()
				if err := f.vq.PushInsert(uuid, newVec, base+5); err != nil {
					t.Fatal(err)
				}
			},
			wantExists:     true,
			wantVec:        newVec,
			wantTS:         func(base int64) int64 { return base + 5 },
			wantOldRemoved: true,
		},
		{
			name: "older insert over the indexed object is dropped",
			op: func(t *testing.T, f *faiss, base int64) {
				t.Helper()
				if err := f.vq.PushInsert(uuid, newVec, base-5); err != nil {
					t.Fatal(err)
				}
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := newBufferingFaiss(t)
			defer closeKVS(t, f)

			base := time.Now().UnixNano()
			if err := f.InsertWithTime(uuid, oldVec, base); err != nil {
				t.Fatal(err)
			}
			if err := f.CreateIndex(t.Context()); err != nil {
				t.Fatal(err)
			}
			oldOid, _, ok := f.kvs.Get(uuid)
			if !ok {
				t.Fatal("kvs entry is missing after the first CreateIndex")
			}

			tt.op(t, f, base)
			if err := f.CreateIndex(t.Context()); err != nil {
				t.Fatal(err)
			}

			if f.vq.IVQLen() != 0 || f.vq.DVQLen() != 0 {
				t.Fatalf("vqueue not drained: IVQLen = %d, DVQLen = %d", f.vq.IVQLen(), f.vq.DVQLen())
			}
			// Every object is still buffered, so the buffer must hold exactly the kvs entries.
			if uint64(len(f.addIds)) != f.kvs.Len() || len(f.addVecs) != len(f.addIds)*f.dim {
				t.Fatalf("buffer: %d ids, %d floats, kvs.Len() = %d", len(f.addIds), len(f.addVecs), f.kvs.Len())
			}
			if got := slices.Contains(f.addIds, int64(oldOid)); got == tt.wantOldRemoved {
				t.Fatalf("old oid %d buffered = %v, want %v: %v", oldOid, got, !tt.wantOldRemoved, f.addIds)
			}

			oid, kts, ok := f.kvs.Get(uuid)
			if ok != tt.wantExists {
				t.Fatalf("kvs entry exists = %v, want %v", ok, tt.wantExists)
			}
			if !tt.wantExists {
				if _, _, err := f.GetObject(uuid); err == nil {
					t.Error("GetObject of a deleted uuid succeeded")
				}
				return
			}
			if !tt.wantOldRemoved && oid != oldOid {
				t.Errorf("oid = %d, want %d: a skipped operation must not repoint the uuid", oid, oldOid)
			}
			if kts != tt.wantTS(base) {
				t.Errorf("kvs ts = %d, want %d", kts, tt.wantTS(base))
			}
			i := slices.Index(f.addIds, int64(oid))
			if i < 0 {
				t.Fatalf("oid %d is not buffered: %v", oid, f.addIds)
			}
			if got := f.addVecs[i*f.dim : (i+1)*f.dim]; !slices.Equal(got, tt.wantVec) {
				t.Errorf("buffered vector = %v, want %v", got, tt.wantVec)
			}
		})
	}
}
