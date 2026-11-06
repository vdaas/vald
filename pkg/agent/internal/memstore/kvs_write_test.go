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

// TestReadPaths_doNotWriteKVS pins that Exists and GetObject are pure reads: an insert that is
// superseded by a pending delete (dts == its) must not bump the kvs timestamp, and their
// observable results stay the same.
func TestReadPaths_doNotWriteKVS(t *testing.T) {
	t.Parallel()

	const (
		uuid = "uuid"
		oid  = uint32(7)
		kts  = int64(1)
		its  = int64(2)
	)
	kvVec := []float32{1, 2, 3}
	tests := []struct {
		call func(kv kvs.BidiMap, vq vqueue.Queue) error
		name string
	}{
		{
			name: "Exists",
			call: func(kv kvs.BidiMap, vq vqueue.Queue) error {
				got, ok := Exists(kv, vq, uuid)
				if !ok || got != oid {
					return errors.Errorf("Exists() = (%d, %v), want (%d, true)", got, ok, oid)
				}
				return nil
			},
		},
		{
			name: "GetObject",
			call: func(kv kvs.BidiMap, vq vqueue.Queue) error {
				vec, ts, err := GetObject(kv, vq, uuid, func(got uint32) ([]float32, error) {
					if got != oid {
						return nil, errors.Errorf("getVectorFn oid = %d, want %d", got, oid)
					}
					return kvVec, nil
				})
				if err != nil || ts != kts || !slices.Equal(vec, kvVec) {
					return errors.Errorf("GetObject() = (%v, %d, %v), want (%v, %d, nil)", vec, ts, err, kvVec, kts)
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			kv := kvs.New()
			kv.Set(uuid, oid, kts)
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			// update pattern with equal timestamps: the insert is not visible (vqok=false, its != 0).
			if err := vq.PushInsert(uuid, []float32{4, 5, 6}, its); err != nil {
				t.Fatal(err)
			}
			if err := vq.PushDelete(uuid, its); err != nil {
				t.Fatal(err)
			}

			if err := tt.call(kv, vq); err != nil {
				t.Error(err)
			}
			gotOid, gotTs, ok := kv.Get(uuid)
			if !ok || gotOid != oid || gotTs != kts {
				t.Errorf("kvs after %s = (%d, %d, %v), want (%d, %d, true)", tt.name, gotOid, gotTs, ok, oid, kts)
			}
			if _, gotTs, ok := kv.GetInverse(oid); !ok || gotTs != kts {
				t.Errorf("inverse kvs after %s = (%d, %v), want (%d, true)", tt.name, gotTs, ok, kts)
			}
		})
	}
}
