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
	"testing"

	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

// afterGetKV runs afterGet once, right after the first Get, to interleave a concurrent
// kvs writer (CreateIndex or another UpdateTimestamp) between UpdateTimestamp's read and write.
type afterGetKV struct {
	kvs.BidiMap
	afterGet func(kv kvs.BidiMap)
}

func (k *afterGetKV) Get(key string) (uint32, int64, bool) {
	oid, ts, ok := k.BidiMap.Get(key)
	if f := k.afterGet; f != nil {
		k.afterGet = nil
		f(k.BidiMap)
	}
	return oid, ts, ok
}

// TestUpdateTimestamp_kvsWriteRacingCreateIndex pins that UpdateTimestamp never re-creates a
// kvs mapping that CreateIndex removed after UpdateTimestamp read it: kv.Set used to map the
// uuid back to an oid already removed from the index. A concurrent timestamp change alone
// must still succeed instead of reporting NotFound.
func TestUpdateTimestamp_kvsWriteRacingCreateIndex(t *testing.T) {
	t.Parallel()

	const (
		uuid = "uuid"
		oid  = uint32(7)
		kts  = int64(1)
		ts   = int64(10)
	)
	tests := []struct {
		queue    func(vq vqueue.Queue) error
		afterGet func(kv kvs.BidiMap)
		name     string
		wantErr  bool
		wantKV   bool
	}{
		{
			name: "kvs only, deleted by CreateIndex after the read",
			queue: func(vqueue.Queue) error {
				return nil
			},
			afterGet: func(kv kvs.BidiMap) { kv.Delete(uuid) },
			wantErr:  true,
		},
		{
			name: "kvs with a queued delete, applied by CreateIndex after the read",
			queue: func(vq vqueue.Queue) error {
				return vq.PushDelete(uuid, 2)
			},
			afterGet: func(kv kvs.BidiMap) { kv.Delete(uuid) },
			wantErr:  true,
		},
		{
			name: "kvs with a cancelled insert, deleted by CreateIndex after the read",
			queue: func(vq vqueue.Queue) error {
				if err := vq.PushInsert(uuid, []float32{1}, 2); err != nil {
					return err
				}
				return vq.PushDelete(uuid, 3)
			},
			afterGet: func(kv kvs.BidiMap) { kv.Delete(uuid) },
			wantErr:  true,
		},
		{
			name: "concurrent timestamp change still succeeds",
			queue: func(vqueue.Queue) error {
				return nil
			},
			afterGet: func(kv kvs.BidiMap) { kv.Set(uuid, oid, kts+1) },
			wantKV:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			base := kvs.New()
			t.Cleanup(func() {
				if err := base.Close(); err != nil {
					t.Error(err)
				}
			})
			base.Set(uuid, oid, kts)
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			if err = tt.queue(vq); err != nil {
				t.Fatal(err)
			}

			err = UpdateTimestamp(&afterGetKV{BidiMap: base, afterGet: tt.afterGet}, vq, uuid, ts, false)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateTimestamp() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotOid, gotTs, ok := base.Get(uuid)
			if ok != tt.wantKV {
				t.Fatalf("kvs entry exists = %v (oid %d, ts %d), want %v", ok, gotOid, gotTs, tt.wantKV)
			}
			_, invTs, invOk := base.GetInverse(oid)
			if invOk != tt.wantKV {
				t.Fatalf("inverse kvs entry exists = %v, want %v", invOk, tt.wantKV)
			}
			if !tt.wantKV {
				if base.Len() != 0 {
					t.Errorf("kvs Len = %d, want 0", base.Len())
				}
				return
			}
			if gotOid != oid || gotTs != ts || invTs != ts {
				t.Errorf("kvs = (%d, %d), inverse ts %d, want (%d, %d)", gotOid, gotTs, invTs, oid, ts)
			}
		})
	}
}
