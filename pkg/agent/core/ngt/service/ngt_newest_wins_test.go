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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Test_ngt_createIndexLocked_newestWins pins that CreateIndex applies a queued delete or
// insert only when it is newer than the indexed object. A delete older than the indexed
// object used to remove it, and a stale Update (delete at t, insert at t+1 <= indexed ts)
// used to repoint the uuid to the stale vector. An insert applied over an indexed object
// must remove the replaced object instead of orphaning it in the core.
func Test_ngt_createIndexLocked_newestWins(t *testing.T) {
	t.Parallel()

	const uuid = "uuid"
	oldVec := []float32{1, 2, 3}
	newVec := []float32{4, 5, 6}
	tests := []struct {
		op         func(t *testing.T, n *ngt, base int64)
		wantTS     func(base int64) int64
		name       string
		wantVec    []float32
		wantExists bool
		// wantOldRemoved reports whether the initially indexed object is gone from the core.
		wantOldRemoved bool
	}{
		{
			name: "delete older than the indexed object is skipped",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.DeleteWithTime(uuid, base-10))
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
		{
			name: "stale update is skipped",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.UpdateWithTime(uuid, newVec, base-5))
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
		{
			name: "update whose insert ties the indexed timestamp is skipped",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.UpdateWithTime(uuid, newVec, base-1))
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
		{
			name: "delete at the indexed timestamp is applied",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.DeleteWithTime(uuid, base))
			},
			wantOldRemoved: true,
		},
		{
			name: "newer delete is applied",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.DeleteWithTime(uuid, base+10))
			},
			wantOldRemoved: true,
		},
		{
			name: "newer update is applied",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.UpdateWithTime(uuid, newVec, base+10))
			},
			wantExists:     true,
			wantVec:        newVec,
			wantTS:         func(base int64) int64 { return base + 11 },
			wantOldRemoved: true,
		},
		{
			name: "UpdateTimestamp of a queued update keeps the new vector",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.UpdateWithTime(uuid, newVec, base+10))
				require.NoError(t, n.UpdateTimestamp(uuid, base+20, false))
			},
			wantExists:     true,
			wantVec:        newVec,
			wantTS:         func(base int64) int64 { return base + 20 },
			wantOldRemoved: true,
		},
		{
			// UpdateTimestamp used to raise the kvs timestamp past the queued delete, so
			// CreateIndex skipped it and served the superseded indexed vector.
			name: "UpdateTimestamp after a queued update and a newer delete does not resurrect",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.UpdateWithTime(uuid, newVec, base+10))
				require.NoError(t, n.DeleteWithTime(uuid, base+20))
				require.Error(t, n.UpdateTimestamp(uuid, base+30, false))
			},
			wantOldRemoved: true,
		},
		{
			name: "newer insert over the indexed object replaces it",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.vq.PushInsert(uuid, newVec, base+5))
			},
			wantExists:     true,
			wantVec:        newVec,
			wantTS:         func(base int64) int64 { return base + 5 },
			wantOldRemoved: true,
		},
		{
			name: "older insert over the indexed object is dropped",
			op: func(t *testing.T, n *ngt, base int64) {
				t.Helper()
				require.NoError(t, n.vq.PushInsert(uuid, newVec, base-5))
			},
			wantExists: true,
			wantVec:    oldVec,
			wantTS:     func(base int64) int64 { return base },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := defaultConfig
			cfg.Dimension = len(oldVec)

			nn, err := New(&cfg)
			require.NoError(t, err)
			n, ok := nn.(*ngt)
			require.True(t, ok)

			base := time.Now().UnixNano()
			require.NoError(t, n.InsertWithTime(uuid, oldVec, base))
			require.NoError(t, n.CreateIndex(t.Context(), 10))
			oldOid, _, ok := n.kvs.Get(uuid)
			require.True(t, ok)

			tt.op(t, n, base)
			require.NoError(t, n.CreateIndex(t.Context(), 10))

			require.Zero(t, n.vq.IVQLen())
			require.Zero(t, n.vq.DVQLen())
			coreVec, err := n.core.GetVector(uint(oldOid))
			if tt.wantOldRemoved {
				// NGT reuses removed object IDs, so an insert applied after the old object was
				// removed may get oldOid back. Only require that oldOid no longer holds oldVec.
				if err == nil {
					require.NotEqual(t, oldVec, coreVec, "the replaced or deleted object must be removed from the core")
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, oldVec, coreVec)
			}

			oid, kts, ok := n.kvs.Get(uuid)
			require.Equal(t, tt.wantExists, ok)
			if !tt.wantExists {
				require.Zero(t, n.Len())
				_, _, err = n.GetObject(uuid)
				require.Error(t, err)
				return
			}
			require.Equal(t, uint64(1), n.Len())
			if !tt.wantOldRemoved {
				require.Equal(t, oldOid, oid, "a skipped operation must not repoint the uuid")
			}
			require.Equal(t, tt.wantTS(base), kts)
			vec, ts, err := n.GetObject(uuid)
			require.NoError(t, err)
			require.Equal(t, tt.wantVec, vec)
			require.Equal(t, tt.wantTS(base), ts)
		})
	}
}
