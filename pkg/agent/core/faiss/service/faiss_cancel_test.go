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
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

// doneDuringDeletePhaseCtx is done from the start (so vqueue ranges stop after their first
// entry), but its Err reports nil on the first call, i.e. it looks live when CreateIndex
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

// newBufferingFaiss returns a faiss on stubCore that needs 78 buffered vectors to train
// (max(nlist, 2^m)*minPointsPerCentroid), so small tests never reach the C library.
func newBufferingFaiss(t *testing.T) *faiss {
	t.Helper()
	vq, err := vqueue.New()
	if err != nil {
		t.Fatal(err)
	}
	f := &faiss{
		dim:              2,
		nlist:            1,
		m:                1,
		fmap:             make(map[string]int64),
		kvs:              kvs.New(kvs.WithConcurrency(1)),
		vq:               vq,
		core:             new(stubCore),
		kvsdbConcurrency: 1,
	}
	f.indexing.Store(false)
	f.saving.Store(false)
	return f
}

// Test_faiss_CreateIndex_doesNotApplyInsertsAfterPartialDeletes pins that a run whose ctx is
// done never applies an insert while the delete of the same Update may still be queued.
func Test_faiss_CreateIndex_doesNotApplyInsertsAfterPartialDeletes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ctx       func(t *testing.T) context.Context
		name      string
		wantIVQ   int
		wantDVQ   int
		wantNocie uint64 // increase of nocie by the cancelled run
	}{
		{
			name: "ctx done before the run starts leaves the vqueue untouched",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				return ctx
			},
			wantIVQ:   2,
			wantDVQ:   2,
			wantNocie: 0,
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
			wantNocie: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := newBufferingFaiss(t)
			defer closeKVS(t, f)

			base := time.Now().UnixNano()
			uuids := []string{"uuid-a", "uuid-b"}
			for _, uuid := range uuids {
				if err := f.InsertWithTime(uuid, []float32{1, 2}, base); err != nil {
					t.Fatal(err)
				}
			}
			if err := f.CreateIndex(t.Context()); err != nil {
				t.Fatal(err)
			}
			oids := make(map[string]uint32, len(uuids))
			for _, uuid := range uuids {
				oid, _, ok := f.kvs.Get(uuid)
				if !ok {
					t.Fatalf("kvs entry of %s is missing", uuid)
				}
				oids[uuid] = oid
			}
			// Update pairs: delete at t, insert at t+1.
			newVec := []float32{4, 5}
			for _, uuid := range uuids {
				if err := f.vq.PushDelete(uuid, base+10); err != nil {
					t.Fatal(err)
				}
				if err := f.vq.PushInsert(uuid, newVec, base+11); err != nil {
					t.Fatal(err)
				}
			}
			nocie := f.NumberOfCreateIndexExecution()

			if err := f.CreateIndex(tt.ctx(t)); !errors.Is(err, context.Canceled) {
				t.Fatalf("CreateIndex error = %v, want context.Canceled", err)
			}
			if got := f.vq.IVQLen(); got != tt.wantIVQ {
				t.Errorf("IVQLen = %d, want %d: no insert may be applied", got, tt.wantIVQ)
			}
			if got := f.vq.DVQLen(); got != tt.wantDVQ {
				t.Errorf("DVQLen = %d, want %d", got, tt.wantDVQ)
			}
			if got := f.NumberOfCreateIndexExecution(); got != nocie+tt.wantNocie {
				t.Errorf("nocie = %d, want %d", got, nocie+tt.wantNocie)
			}
			// The add buffer must hold exactly the uuids still mapped in kvs.
			if uint64(len(f.addIds)) != f.kvs.Len() || len(f.addVecs) != len(f.addIds)*f.dim {
				t.Errorf("buffer: %d ids, %d floats, kvs.Len() = %d", len(f.addIds), len(f.addVecs), f.kvs.Len())
			}
			for _, uuid := range uuids {
				// a uuid whose delete is still queued must keep its old mapping
				if _, _, dts, _ := f.vq.GetVectorWithTimestamp(uuid); dts != 0 {
					oid, _, ok := f.kvs.Get(uuid)
					if !ok || oid != oids[uuid] {
						t.Errorf("%s: oid = %d (ok %v), want %d", uuid, oid, ok, oids[uuid])
					}
				}
			}

			// The next run with a live ctx applies the rest.
			if err := f.CreateIndex(t.Context()); err != nil {
				t.Fatal(err)
			}
			if f.vq.IVQLen() != 0 || f.vq.DVQLen() != 0 {
				t.Fatalf("vqueue not drained: IVQLen = %d, DVQLen = %d", f.vq.IVQLen(), f.vq.DVQLen())
			}
			if got := f.Len(); got != uint64(len(uuids)) {
				t.Fatalf("Len = %d, want %d", got, len(uuids))
			}
			if len(f.addIds) != len(uuids) || len(f.addVecs) != len(uuids)*f.dim {
				t.Fatalf("buffer: %d ids, %d floats, want %d ids", len(f.addIds), len(f.addVecs), len(uuids))
			}
			for _, uuid := range uuids {
				oid, ts, ok := f.kvs.Get(uuid)
				if !ok {
					t.Fatalf("kvs entry of %s is missing", uuid)
				}
				if ts != base+11 {
					t.Errorf("%s: ts = %d, want %d", uuid, ts, base+11)
				}
				i := slices.Index(f.addIds, int64(oid))
				if i < 0 {
					t.Errorf("%s (oid %d) is not buffered: %v", uuid, oid, f.addIds)
					continue
				}
				if got := f.addVecs[i*f.dim : (i+1)*f.dim]; !slices.Equal(got, newVec) {
					t.Errorf("buffered vector of %s = %v, want %v", uuid, got, newVec)
				}
			}
		})
	}
}

// Test_faiss_Start_doesNotCreateIndexOnDoneCtx pins that Start's shutdown branch leaves the
// vqueue to Close: a CreateIndex there runs with a done ctx and could apply part of the queue.
func Test_faiss_Start_doesNotCreateIndexOnDoneCtx(t *testing.T) {
	t.Parallel()
	f := newBufferingFaiss(t)
	defer closeKVS(t, f)
	for _, uuid := range []string{"a", "b"} {
		if err := f.Insert(uuid, []float32{1, 2}); err != nil {
			t.Fatal(err)
		}
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	f.eg, _ = errgroup.New(t.Context())
	cancel()
	ech := f.Start(ctx)

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	for done := false; !done; {
		select {
		case err, ok := <-ech:
			if !ok {
				done = true
				break
			}
			t.Errorf("unexpected error from Start: %v", err)
		case <-timer.C:
			t.Fatal("Start did not stop on a done ctx")
		}
	}
	if got := f.vq.IVQLen(); got != 2 {
		t.Errorf("IVQLen = %d, want 2: Start must not index on shutdown", got)
	}
	if len(f.addIds) != 0 {
		t.Errorf("buffered %d vectors on shutdown, want 0", len(f.addIds))
	}
}

// Test_faiss_Close_waitsForInflightIndexing pins that Close does not run its final
// CreateIndex/SaveIndex while an earlier run is still in flight, and that the wait is bounded
// by ctx.
func Test_faiss_Close_waitsForInflightIndexing(t *testing.T) {
	t.Parallel()
	tests := []struct {
		set  func(f *faiss)
		name string
	}{
		{name: "indexing", set: func(f *faiss) { f.indexing.Store(true) }},
		{name: "saving", set: func(f *faiss) { f.saving.Store(true) }},
		{name: "waiting CreateIndex", set: func(f *faiss) { f.wfci.Store(1) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := newBufferingFaiss(t)
			f.path = t.TempDir()
			tt.set(f)

			// The in-flight flag is never cleared, so Close can only stop waiting on ctx.
			ctx, cancel := context.WithTimeout(t.Context(), 300*time.Millisecond)
			defer cancel()
			if err := f.Close(ctx); err != nil {
				t.Fatalf("Close error = %v, want nil (ctx errors are dropped)", err)
			}
			if ctx.Err() == nil {
				t.Error("Close returned before the in-flight operation finished or ctx was done")
			}
		})
	}

	t.Run("idle returns immediately", func(t *testing.T) {
		t.Parallel()
		f := newBufferingFaiss(t)
		defer closeKVS(t, f)
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		if err := f.waitForInflightIndexing(ctx); err != nil {
			t.Fatalf("waitForInflightIndexing error = %v, want nil", err)
		}
	})
}
