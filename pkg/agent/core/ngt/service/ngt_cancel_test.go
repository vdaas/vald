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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vdaas/vald/internal/sync/atomic"
)

// doneDuringDeletePhaseCtx is done from the start (so vqueue ranges stop after their first
// entry), but its Err reports nil on the first call, i.e. it looks live when createIndexLocked
// starts and is cancelled while the delete phase runs.
//
//nolint:containedctx // test fake overriding Done/Err of a parent ctx
type doneDuringDeletePhaseCtx struct {
	context.Context
	done     chan struct{}
	errCalls atomic.Int32
}

func (c *doneDuringDeletePhaseCtx) Done() <-chan struct{} { return c.done }

func (c *doneDuringDeletePhaseCtx) Err() error {
	if c.errCalls.Add(1) == 1 {
		return nil
	}
	return context.Canceled
}

// Test_ngt_createIndexLocked_doesNotApplyInsertsAfterPartialDeletes pins that a run whose ctx
// is done never applies an insert while the delete of the same Update may still be queued.
func Test_ngt_createIndexLocked_doesNotApplyInsertsAfterPartialDeletes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ctx         func(t *testing.T) context.Context
		name        string
		wantIVQ     int
		wantDVQ     int
		wantCreated bool
	}{
		{
			name: "ctx done before the run starts leaves the vqueue untouched",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				return ctx
			},
			wantIVQ:     2,
			wantDVQ:     2,
			wantCreated: false,
		},
		{
			name: "ctx done during the delete phase keeps every insert queued",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				done := make(chan struct{})
				close(done)
				return &doneDuringDeletePhaseCtx{Context: t.Context(), done: done}
			},
			wantIVQ: 2,
			wantDVQ: 1,
			// the delete already applied is still finished and counted so that it gets saved
			wantCreated: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := defaultConfig
			cfg.Dimension = 3

			nn, err := New(&cfg)
			require.NoError(t, err)
			n, ok := nn.(*ngt)
			require.True(t, ok)

			base := time.Now().UnixNano()
			uuids := []string{"uuid-a", "uuid-b"}
			for _, uuid := range uuids {
				require.NoError(t, n.InsertWithTime(uuid, []float32{1, 2, 3}, base))
			}
			require.NoError(t, n.CreateIndex(t.Context(), 10))
			oids := make(map[string]uint32, len(uuids))
			for _, uuid := range uuids {
				oid, _, ok := n.kvs.Get(uuid)
				require.True(t, ok)
				oids[uuid] = oid
			}
			// Update pairs: delete at t, insert at t+1.
			newVec := []float32{4, 5, 6}
			for _, uuid := range uuids {
				require.NoError(t, n.vq.PushDelete(uuid, base+10))
				require.NoError(t, n.vq.PushInsert(uuid, newVec, base+11))
			}
			nocie := n.NumberOfCreateIndexExecution()

			created, err := n.createIndexLocked(tt.ctx(t), 10)
			require.ErrorIs(t, err, context.Canceled)
			require.Equal(t, tt.wantCreated, created)
			require.Equal(t, tt.wantIVQ, n.vq.IVQLen(), "no insert may be applied")
			require.Equal(t, tt.wantDVQ, n.vq.DVQLen())
			if tt.wantCreated {
				require.Equal(t, nocie+1, n.NumberOfCreateIndexExecution())
			} else {
				require.Equal(t, nocie, n.NumberOfCreateIndexExecution())
			}
			for _, uuid := range uuids {
				// a uuid whose delete is still queued must keep its old mapping
				if _, _, dts, _ := n.vq.GetVectorWithTimestamp(uuid); dts != 0 {
					oid, _, ok := n.kvs.Get(uuid)
					require.True(t, ok)
					require.Equal(t, oids[uuid], oid)
				}
			}

			// The next run with a live ctx applies the rest.
			require.NoError(t, n.CreateIndex(t.Context(), 10))
			require.Zero(t, n.vq.IVQLen())
			require.Zero(t, n.vq.DVQLen())
			require.Equal(t, uint64(len(uuids)), n.Len())
			for _, uuid := range uuids {
				vec, ts, err := n.GetObject(uuid)
				require.NoError(t, err)
				require.Equal(t, newVec, vec)
				require.Equal(t, base+11, ts)
			}
		})
	}
}

// Test_ngt_copyNGT_keepsCreateIndexWaiters pins that RegenerateIndexes' copyNGT does not reset
// the waiter count of CreateIndex callers blocked on this instance; resetting it made the
// waiter's decrement wrap wfci, disabling CreateIndex dedup and hanging Close's wait.
func Test_ngt_copyNGT_keepsCreateIndexWaiters(t *testing.T) {
	t.Parallel()

	n := new(ngt)
	atomic.StoreUint64(&n.wfci, 1)
	n.copyNGT(new(ngt))
	require.Equal(t, uint64(1), atomic.LoadUint64(&n.wfci))

	// the waiter leaves its wait loop
	atomic.AddUint64(&n.wfci, ^uint64(0))
	require.Zero(t, atomic.LoadUint64(&n.wfci))

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	require.NoError(t, n.waitForInflightIndexing(ctx), "Close must not wait for a waiter that already left")
}
