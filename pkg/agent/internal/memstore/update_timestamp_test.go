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
	"slices"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

// TestUpdateTimestamp_queuedInsertAndDelete pins UpdateTimestamp for an indexed object with
// both a queued insert and a queued delete that is not older than it. Raising the kvs
// timestamp to ts used to make CreateIndex skip the queued delete (newest wins), so a newer
// delete was lost and the superseded indexed vector was served again.
func TestUpdateTimestamp_queuedInsertAndDelete(t *testing.T) {
	t.Parallel()

	const (
		uuid = "uuid"
		oid  = uint32(7)
		kts  = int64(50)
	)
	qvec := []float32{1, 2}
	tests := []struct {
		// wantErr is the expected error; nil means success.
		wantErr error
		name    string
		its     int64
		dts     int64
		ts      int64
		// wantIts is the queued insert timestamp afterwards.
		wantIts int64
		force   bool
	}{
		{
			name:    "newer delete cancels the insert: pending deletion is not found",
			its:     101,
			dts:     200,
			ts:      300,
			wantErr: errors.ErrObjectNotFound(nil, uuid),
			wantIts: 101,
		},
		{
			name:    "force on a pending deletion is not found",
			its:     101,
			dts:     200,
			ts:      150,
			force:   true,
			wantErr: errors.ErrObjectNotFound(nil, uuid),
			wantIts: 101,
		},
		{
			name:    "delete tying the insert: the queued vector is re-queued at ts",
			its:     100,
			dts:     100,
			ts:      300,
			wantIts: 300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			kv := kvs.New()
			t.Cleanup(func() {
				if err := kv.Close(); err != nil {
					t.Error(err)
				}
			})
			kv.Set(uuid, oid, kts)
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			if err = vq.PushInsert(uuid, qvec, tt.its); err != nil {
				t.Fatal(err)
			}
			if err = vq.PushDelete(uuid, tt.dts); err != nil {
				t.Fatal(err)
			}

			err = UpdateTimestamp(kv, vq, uuid, tt.ts, tt.force)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("UpdateTimestamp() error = %v, want nil", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("UpdateTimestamp() error = %v, want %v", err, tt.wantErr)
			}

			// The kvs timestamp must not pass the queued delete, or CreateIndex skips it.
			gotOid, gotKts, ok := kv.Get(uuid)
			if !ok || gotOid != oid || gotKts != kts {
				t.Errorf("kvs = (%d, %d, %v), want (%d, %d, true)", gotOid, gotKts, ok, oid, kts)
			}
			vec, its, dts, _ := vq.GetVectorWithTimestamp(uuid)
			if !slices.Equal(vec, qvec) || its != tt.wantIts {
				t.Errorf("queued insert = (%v, %d), want (%v, %d)", vec, its, qvec, tt.wantIts)
			}
			if dts != tt.dts {
				t.Errorf("queued delete ts = %d, want %d", dts, tt.dts)
			}
		})
	}
}

// When the indexed object is newer than both queued entries, CreateIndex skips
// them (newest wins), so UpdateTimestamp must update the indexed object and drop
// the stale insert; re-queuing that insert at ts would let it replace the newer
// indexed object at the next CreateIndex.
func TestUpdateTimestamp_indexedNewerThanQueued(t *testing.T) {
	t.Parallel()

	const (
		uuid = "uuid"
		oid  = uint32(7)
		kts  = int64(200)
		ts   = int64(300)
	)
	tests := []struct {
		name string
		its  int64
		dts  int64
	}{
		{name: "delete tying the stale insert", its: 100, dts: 100},
		{name: "delete newer than the stale insert", its: 100, dts: 150},
		{name: "stale insert without a delete", its: 100},
		{name: "stale insert tying the indexed object", its: kts},
		{name: "stale insert newer than its delete", its: 150, dts: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			kv := kvs.New()
			t.Cleanup(func() {
				if err := kv.Close(); err != nil {
					t.Error(err)
				}
			})
			kv.Set(uuid, oid, kts)
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			if err = vq.PushInsert(uuid, []float32{1, 2}, tt.its); err != nil {
				t.Fatal(err)
			}
			if tt.dts != 0 {
				if err = vq.PushDelete(uuid, tt.dts); err != nil {
					t.Fatal(err)
				}
			}

			if err = UpdateTimestamp(kv, vq, uuid, ts, false); err != nil {
				t.Fatalf("UpdateTimestamp() error = %v, want nil", err)
			}
			gotOid, gotKts, ok := kv.Get(uuid)
			if !ok || gotOid != oid || gotKts != ts {
				t.Errorf("kvs = (%d, %d, %v), want (%d, %d, true)", gotOid, gotKts, ok, oid, ts)
			}
			if vec, its, dts, _ := vq.GetVectorWithTimestamp(uuid); vec != nil || its != 0 || dts != 0 {
				t.Errorf("stale entries still queued: insert (%v, %d), delete %d", vec, its, dts)
			}
		})
	}
}
