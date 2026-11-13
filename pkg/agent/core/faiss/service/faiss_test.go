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
	"encoding/gob"
	"math"
	"slices"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/metadata"
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

// Test_faiss_toSearchResponse_dropsNaNDistance verifies that results whose
// distance is NaN are dropped while the order of the finite results is kept,
// and that an all-NaN result set on a non-empty index is reported as
// ErrEmptySearchResult.
func Test_faiss_toSearchResponse_dropsNaNDistance(t *testing.T) {
	t.Parallel()
	nan := float32(math.NaN())
	// kvs oid -> uuid: 1 -> "a", 2 -> "b", 3 -> "c"
	stored := []string{"a", "b", "c"}
	tests := []struct {
		wantErr   error
		name      string
		stored    []string
		sr        []algorithm.SearchResult
		wantIDs   []string
		wantDists []float32
		wantNil   bool
	}{
		{
			name:   "mixed NaN and finite keeps only finite in the original order",
			stored: stored,
			sr: []algorithm.SearchResult{
				{ID: 3, Distance: 0.3},
				{ID: 1, Distance: nan},
				{ID: 2, Distance: 0.5},
				{ID: 4, Distance: nan},
			},
			wantIDs:   []string{"c", "b"},
			wantDists: []float32{0.3, 0.5},
		},
		{
			name:   "all NaN on a non-empty index returns ErrEmptySearchResult",
			stored: stored,
			sr: []algorithm.SearchResult{
				{ID: 1, Distance: nan},
				{ID: 2, Distance: nan},
			},
			wantErr: errors.ErrEmptySearchResult,
			wantNil: true,
		},
		{
			name:    "all NaN on an empty index returns nil, nil",
			sr:      []algorithm.SearchResult{{ID: 1, Distance: nan}},
			wantNil: true,
		},
		{
			name:    "empty input on an empty index returns nil, nil",
			wantNil: true,
		},
		{
			name:    "empty input on a non-empty index returns ErrEmptySearchResult",
			stored:  stored,
			wantErr: errors.ErrEmptySearchResult,
			wantNil: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			kv := kvs.New(kvs.WithConcurrency(1))
			now := time.Now().UnixNano()
			for i, uuid := range tc.stored {
				kv.Set(uuid, uint32(i+1), now)
			}
			f := &faiss{kvs: kv}
			res, err := f.toSearchResponse(tc.sr)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("toSearchResponse() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantNil {
				if res != nil {
					t.Fatalf("toSearchResponse() = %v, want nil", res)
				}
				return
			}
			var gotIDs []string
			var gotDists []float32
			for _, r := range res.GetResults() {
				gotIDs = append(gotIDs, r.GetId())
				gotDists = append(gotDists, r.GetDistance())
			}
			if !slices.Equal(gotIDs, tc.wantIDs) || !slices.Equal(gotDists, tc.wantDists) {
				t.Errorf("toSearchResponse() = ids %v dists %v, want ids %v dists %v", gotIDs, gotDists, tc.wantIDs, tc.wantDists)
			}
		})
	}
}

// addRecordCore records the vectors passed to Add.
type addRecordCore struct {
	stubCore
	vecs []float32
	ids  []int64
}

func (c *addRecordCore) Add(n int, xb []float32, xids []int64) (int, error) {
	c.vecs = slices.Clone(xb)
	c.ids = slices.Clone(xids)
	return n, nil
}

// Test_faiss_CreateIndex_skipsNaNOrInfVector verifies that a queued NaN/Inf
// vector is dropped by CreateIndex instead of reaching Add, so it cannot make
// every later indexing run fail.
func Test_faiss_CreateIndex_skipsNaNOrInfVector(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		bad  float32
	}{
		{name: "NaN", bad: float32(math.NaN())},
		{name: "+Inf", bad: float32(math.Inf(1))},
		{name: "-Inf", bad: float32(math.Inf(-1))},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			c := new(addRecordCore)
			f := &faiss{
				dim:              2,
				nlist:            1,
				m:                1,
				path:             t.TempDir(),
				fmap:             make(map[string]int64),
				kvs:              kvs.New(kvs.WithConcurrency(1)),
				vq:               vq,
				core:             c,
				kvsdbConcurrency: 1,
			}
			defer closeKVS(t, f)
			f.indexing.Store(false)
			f.saving.Store(false)
			f.isTrained.Store(true)
			if err = f.Insert("good", []float32{1, 2}); err != nil {
				t.Fatal(err)
			}
			if err = f.Insert("bad", []float32{1, tc.bad}); err != nil {
				t.Fatal(err)
			}
			if err = f.CreateIndex(t.Context()); err != nil {
				t.Fatalf("CreateIndex() err = %v", err)
			}
			if want := []float32{1, 2}; !slices.Equal(c.vecs, want) {
				t.Errorf("Add() vectors = %v, want %v", c.vecs, want)
			}
			if len(c.ids) != 1 {
				t.Errorf("Add() ids = %v, want exactly one id", c.ids)
			}
			if _, _, ok := f.kvs.Get("bad"); ok {
				t.Error("kvs contains the skipped NaN/Inf vector")
			}
		})
	}
}

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

func writeGob(t *testing.T, path string, v any) {
	t.Helper()
	fi, err := file.Open(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	if err := gob.NewEncoder(fi).Encode(v); err != nil {
		t.Fatal(errors.Join(err, fi.Close()))
	}
	if err := fi.Close(); err != nil {
		t.Fatal(err)
	}
}

func closeKVS(t *testing.T, f *faiss) {
	t.Helper()
	if f.kvs != nil {
		if err := f.kvs.Close(); err != nil {
			t.Error(err)
		}
	}
}

// Test_faiss_loadKVS_restoresIcnt verifies that the id counter resumes after
// the largest loaded id, so a restarted agent does not reuse ids that already
// label vectors of the loaded index.
func Test_faiss_loadKVS_restoresIcnt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		ids  map[string]uint32
		name string
		icnt uint64
		want uint64
	}{
		{
			name: "next id follows the largest loaded id",
			ids:  map[string]uint32{"a": 0, "b": 5, "c": 2},
			want: 6,
		},
		{
			name: "stale counter is replaced by the loaded state",
			ids:  map[string]uint32{"a": 3},
			icnt: 100,
			want: 4,
		},
		{
			name: "empty kvsdb starts from zero",
			ids:  map[string]uint32{},
			want: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			ts := make(map[string]int64, len(test.ids))
			for k := range test.ids {
				ts[k] = 1
			}
			writeGob(t, file.Join(dir, kvsFileName), test.ids)
			writeGob(t, file.Join(dir, kvsTimestampFileName), ts)

			f := &faiss{fmap: make(map[string]int64), kvsdbConcurrency: 1, icnt: test.icnt}
			defer closeKVS(t, f)
			if err := f.loadKVS(t.Context(), dir, time.Minute); err != nil {
				t.Fatal(err)
			}
			if f.icnt != test.want {
				t.Errorf("icnt = %d, want %d", f.icnt, test.want)
			}
			if got := f.kvs.Len(); got != uint64(len(test.ids)) {
				t.Errorf("kvs.Len() = %d, want %d", got, len(test.ids))
			}
		})
	}
}

// Test_faiss_saveIndex_rollbackOnError verifies that a failed save does not
// leave the generation marked as saved, which would make every later
// SaveIndex a no-op until the next CreateIndex.
func Test_faiss_saveIndex_rollbackOnError(t *testing.T) {
	t.Parallel()
	f := new(faiss)
	f.isTrained.Store(true)
	f.nocie.Store(1)
	f.indexing.Store(true) // keeps saveIndex in its wait loop until ctx is done
	f.saving.Store(false)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	for i := range 2 {
		if err := f.saveIndex(ctx); !errors.Is(err, context.Canceled) {
			t.Fatalf("attempt %d: err = %v, want context.Canceled", i, err)
		}
		if got := f.lastNocie.Load(); got != 0 {
			t.Fatalf("attempt %d: lastNocie = %d, want 0 (rolled back)", i, got)
		}
	}
}

// stubCore is a core.Faiss that records Remove calls and fails everything else
// the tests do not expect to reach the C library.
type stubCore struct {
	removed [][]int64
}

func (*stubCore) SaveIndex() error               { return errors.New("unexpected SaveIndex") }
func (*stubCore) SaveIndexWithPath(string) error { return errors.New("unexpected SaveIndexWithPath") }
func (*stubCore) Train(int, []float32) error     { return errors.New("unexpected Train") }
func (*stubCore) Add(int, []float32, []int64) (int, error) {
	return -1, errors.New("unexpected Add")
}

func (*stubCore) Search(int, int, int, []float32) ([]algorithm.SearchResult, error) {
	return nil, errors.New("unexpected Search")
}

func (s *stubCore) Remove(size int, ids []int64) (int, error) {
	s.removed = append(s.removed, slices.Clone(ids[:size]))
	return 0, nil
}
func (*stubCore) Close() {}

// Test_faiss_CreateIndex_deleteBuffered verifies that deleting a vector that
// is still in the untrained add buffer removes it from that buffer, keeping
// addIds and addVecs aligned, instead of later adding it to the index as an
// orphan label.
func Test_faiss_CreateIndex_deleteBuffered(t *testing.T) {
	t.Parallel()
	vq, err := vqueue.New()
	if err != nil {
		t.Fatal(err)
	}
	sc := new(stubCore)
	// max(nlist, 2^m)*minPointsPerCentroid = 78 buffered vectors are needed to
	// train, so the three vectors below stay buffered.
	f := &faiss{
		dim:              2,
		nlist:            1,
		m:                1,
		fmap:             make(map[string]int64),
		kvs:              kvs.New(kvs.WithConcurrency(1)),
		vq:               vq,
		core:             sc,
		kvsdbConcurrency: 1,
	}
	defer closeKVS(t, f)
	f.indexing.Store(false)
	f.saving.Store(false)

	vecs := map[string][]float32{
		"a": {1, 1},
		"b": {2, 2},
		"c": {3, 3},
	}
	for _, uuid := range []string{"a", "b", "c"} {
		if err := f.Insert(uuid, vecs[uuid]); err != nil {
			t.Fatal(err)
		}
	}
	ctx := t.Context()
	if err := f.CreateIndex(ctx); err != nil {
		t.Fatal(err)
	}
	if f.isTrained.Load() || len(f.addIds) != 3 {
		t.Fatalf("setup: isTrained = %v, buffered = %d, want false and 3", f.isTrained.Load(), len(f.addIds))
	}

	if err := f.Delete("a"); err != nil {
		t.Fatal(err)
	}
	if err := f.CreateIndex(ctx); err != nil {
		t.Fatal(err)
	}

	if len(sc.removed) != 0 {
		t.Errorf("core.Remove called with %v for a buffered vector", sc.removed)
	}
	if len(f.addIds) != 2 || len(f.addVecs) != 2*f.dim {
		t.Fatalf("buffer: %d ids, %d floats, want 2 ids, %d floats", len(f.addIds), len(f.addVecs), 2*f.dim)
	}
	if got := f.addLen.Load(); got != 2 {
		t.Errorf("addLen = %d, want 2", got)
	}
	// Both remaining slots must belong to the live uuids b and c, each paired
	// with its own vector; with only two slots this also rules out an orphan.
	for _, uuid := range []string{"b", "c"} {
		oid, _, ok := f.kvs.Get(uuid)
		if !ok {
			t.Fatalf("kvs entry of %s is missing", uuid)
		}
		i := slices.Index(f.addIds, int64(oid))
		if i < 0 {
			t.Errorf("%s (oid %d) is not buffered: %v", uuid, oid, f.addIds)
			continue
		}
		if got := f.addVecs[i*f.dim : (i+1)*f.dim]; !slices.Equal(got, vecs[uuid]) {
			t.Errorf("buffered vector of %s = %v, want %v", uuid, got, vecs[uuid])
		}
	}
}

// Test_faiss_initFaiss_loadTimeout verifies that initFaiss returns on a load
// timeout instead of falling back while the abandoned load is still running.
// The timestamp kvsdb file is a FIFO with no writer, so loadKVS blocks in open
// and load deterministically times out.
func Test_faiss_initFaiss_loadTimeout(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := metadata.Store(file.Join(dir, metadata.AgentMetadataFileName), &metadata.Metadata{
		Faiss: &metadata.Faiss{IndexCount: 1},
	}); err != nil {
		t.Fatal(err)
	}
	writeGob(t, file.Join(dir, kvsFileName), map[string]uint32{"a": 0})
	fifo := file.Join(dir, kvsTimestampFileName)
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	// load() requires a non-empty faiss.index before it starts loading. A
	// corrupt one makes the index load fail fast, so only the kvsdb load
	// (blocked on the FIFO) keeps the concurrent load running past its timeout.
	if err := os.WriteFile(file.Join(dir, faissIndexFileName), []byte("corrupt"), 0o600); err != nil {
		t.Fatal(err)
	}

	eg, _ := errgroup.New(t.Context())
	f := &faiss{
		path:             dir,
		fmap:             make(map[string]int64),
		kvsdbConcurrency: 1,
		eg:               eg,
		minLit:           100 * time.Millisecond,
		maxLit:           100 * time.Millisecond,
		litFactor:        time.Millisecond,
	}

	err := f.initFaiss()
	if !errors.Is(err, errors.ErrIndexLoadTimeout) {
		t.Errorf("initFaiss() err = %v, want ErrIndexLoadTimeout", err)
	}

	// Release the abandoned loadKVS: an O_RDWR open of a FIFO does not block on
	// Linux and serves the reader whether it is already blocked in open or
	// opens later; the fd stays open until the abandoned load has finished.
	w, err := file.Open(fifo, os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := gob.NewEncoder(w).Encode(map[string]int64{}); err != nil {
		t.Error(err)
	}
	done := make(chan error, 1)
	go func() { done <- eg.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			t.Error(err)
		}
		closeKVS(t, f)
	case <-time.After(30 * time.Second):
		t.Error("abandoned load did not finish")
	}
	if err := w.Close(); err != nil {
		t.Error(err)
	}
}

// Test_faiss_resetTrainState verifies that the state restored by a successful
// load is cleared before initFaiss falls back to a fresh core; otherwise
// CreateIndex would skip Train and Add vectors to an untrained index.
func Test_faiss_resetTrainState(t *testing.T) {
	t.Parallel()
	f := &faiss{icnt: 7}
	f.isTrained.Store(true)
	f.resetTrainState()
	if f.isTrained.Load() {
		t.Error("isTrained = true after reset, want false")
	}
	if f.icnt != 0 {
		t.Errorf("icnt = %d after reset, want 0", f.icnt)
	}
}

// Test_faiss_waitLoad verifies that a successful load always marks the index
// as trained (the load goroutine must not cancel ctx before reporting, which
// let ctx.Done() win the select and skip the success path), and that a timed
// out load keeps ErrIndexLoadTimeout even when storing the invalid metadata
// fails, so initFaiss does not start another load beside the abandoned one.
func Test_faiss_waitLoad(t *testing.T) {
	t.Parallel()
	errLoad := errors.New("load failed")
	tests := []struct {
		wantErr     error
		wait        func(release <-chan struct{}) error
		name        string
		timeout     time.Duration
		badMetaPath bool
		wantTrained bool
		wantInvalid bool
	}{
		{
			name:        "success marks the index as trained",
			timeout:     time.Minute,
			wait:        func(<-chan struct{}) error { return nil },
			wantTrained: true,
		},
		{
			name:    "load error is returned and the index stays untrained",
			timeout: time.Minute,
			wait:    func(<-chan struct{}) error { return errLoad },
			wantErr: errLoad,
		},
		{
			name:    "timeout stores invalid metadata and returns ErrIndexLoadTimeout",
			timeout: 10 * time.Millisecond,
			wait: func(release <-chan struct{}) error {
				<-release
				return nil
			},
			wantErr:     errors.ErrIndexLoadTimeout,
			wantInvalid: true,
		},
		{
			name:    "timeout keeps ErrIndexLoadTimeout when the metadata store fails",
			timeout: 10 * time.Millisecond,
			wait: func(release <-chan struct{}) error {
				<-release
				return nil
			},
			badMetaPath: true,
			wantErr:     errors.ErrIndexLoadTimeout,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			metaPath := file.Join(dir, metadata.AgentMetadataFileName)
			if tc.badMetaPath {
				// a path below a regular file fails with ENOTDIR, also as root.
				notDir := file.Join(dir, "not-a-dir")
				writeGob(t, notDir, map[string]uint32{})
				metaPath = file.Join(notDir, metadata.AgentMetadataFileName)
			}
			// the success path is repeated so that a select woken by a
			// concurrent ctx cancellation would be observed.
			runs := 1
			if tc.wantTrained {
				runs = 100
			}
			for range runs {
				eg, _ := errgroup.New(t.Context())
				f := &faiss{eg: eg}
				release := make(chan struct{})
				ctx, cancel := context.WithTimeout(t.Context(), tc.timeout)
				err := f.waitLoad(ctx, func() error { return tc.wait(release) }, dir, metaPath, tc.timeout)
				cancel()
				close(release)
				if tc.wantErr == nil && err != nil || tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
					t.Fatalf("waitLoad() err = %v, want %v", err, tc.wantErr)
				}
				if got := f.isTrained.Load(); got != tc.wantTrained {
					t.Fatalf("isTrained = %v, want %v", got, tc.wantTrained)
				}
				done := make(chan error, 1)
				go func() { done <- eg.Wait() }()
				select {
				case err := <-done:
					if err != nil {
						t.Error(err)
					}
				case <-time.After(30 * time.Second):
					t.Fatal("load goroutine did not finish")
				}
			}
			if tc.wantInvalid {
				meta, err := metadata.Load(metaPath)
				if err != nil {
					t.Fatal(err)
				}
				if !meta.IsInvalid {
					t.Error("metadata IsInvalid = false, want true")
				}
			}
		})
	}
}

// saveOKCore is a stubCore whose SaveIndexWithPath succeeds, so saveIndex
// reaches the metadata write.
type saveOKCore struct {
	stubCore
}

func (*saveOKCore) SaveIndexWithPath(string) error { return nil }

// Test_faiss_saveIndex_cancelledCtxFullSnapshot verifies that a save started
// with a cancelled context still persists the whole kvs. kvs.Range stops
// silently on cancellation, so the snapshot used to be partial (here: empty)
// while kvsdb was truncated and the metadata IndexCount matched the truncated
// count, hiding the loss from the load-time consistency check.
func Test_faiss_saveIndex_cancelledCtxFullSnapshot(t *testing.T) {
	t.Parallel()
	const n = 100
	dir := t.TempDir()
	f := &faiss{
		path: dir,
		fmap: make(map[string]int64),
		kvs:  kvs.New(kvs.WithConcurrency(1)),
		core: new(saveOKCore),
	}
	defer closeKVS(t, f)
	f.isTrained.Store(true)
	f.indexing.Store(false)
	f.saving.Store(false)
	f.nocie.Store(1)
	for i := range n {
		f.kvs.Set("uuid-"+strconv.Itoa(i), uint32(i), int64(i+1))
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := f.saveIndex(ctx); err != nil {
		t.Fatalf("saveIndex() err = %v, want nil", err)
	}

	readGob := func(name string, v any) {
		t.Helper()
		fi, err := file.Open(file.Join(dir, name), os.O_RDONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		if err := gob.NewDecoder(fi).Decode(v); err != nil {
			t.Fatal(errors.Join(err, fi.Close()))
		}
		if err := fi.Close(); err != nil {
			t.Fatal(err)
		}
	}
	m := make(map[string]uint32)
	readGob(kvsFileName, &m)
	if len(m) != n {
		t.Errorf("kvsdb entries = %d, want %d", len(m), n)
	}
	mt := make(map[string]int64)
	readGob(kvsTimestampFileName, &mt)
	if len(mt) != n {
		t.Errorf("kvs timestamp entries = %d, want %d", len(mt), n)
	}
	meta, err := metadata.Load(file.Join(dir, metadata.AgentMetadataFileName))
	if err != nil {
		t.Fatal(err)
	}
	if meta.Faiss == nil || meta.Faiss.IndexCount != n {
		t.Errorf("metadata = %+v, want Faiss.IndexCount %d", meta.Faiss, n)
	}
}

// pathRecordingCore is a stubCore whose SaveIndexWithPath writes a stand-in
// index file and records its path, failing like the faiss writer when the path
// is a directory or its parent directory is missing.
type pathRecordingCore struct {
	stubCore
	saved []string
}

func (c *pathRecordingCore) SaveIndexWithPath(path string) error {
	c.saved = append(c.saved, path)
	// os.OpenFile (unlike file.Open) does not create missing parent
	// directories, matching faiss write_index.
	fi, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = fi.WriteString("index")
	return errors.Join(err, fi.Close())
}

// Test_faiss_saveIndex_indexFile verifies that the index is saved as a file
// inside the index directory and that the directory exists before the index
// writer runs. The directory itself used to be passed to faiss write_index,
// which cannot open a directory, so no index was ever persisted.
func Test_faiss_saveIndex_indexFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		kvs  map[string]uint32
		name string
	}{
		{
			name: "empty kvs creates the missing index directory",
		},
		{
			name: "kvsdb and index file are written to the same directory",
			kvs:  map[string]uint32{"a": 0, "b": 1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			dir := file.Join(t.TempDir(), "index")
			c := new(pathRecordingCore)
			f := &faiss{
				path: dir,
				fmap: make(map[string]int64),
				kvs:  kvs.New(kvs.WithConcurrency(1)),
				core: c,
			}
			defer closeKVS(t, f)
			f.isTrained.Store(true)
			f.indexing.Store(false)
			f.saving.Store(false)
			f.nocie.Store(1)
			for k, id := range test.kvs {
				f.kvs.Set(k, id, 1)
			}

			if err := f.saveIndex(t.Context()); err != nil {
				t.Fatalf("saveIndex() err = %v, want nil", err)
			}
			want := file.Join(dir, faissIndexFileName)
			if !slices.Equal(c.saved, []string{want}) {
				t.Errorf("SaveIndexWithPath paths = %v, want [%s]", c.saved, want)
			}
			if exist, fi, err := file.ExistsWithDetail(want); !exist || err != nil || fi == nil || fi.IsDir() {
				t.Errorf("index file %s: exist = %v, info = %v, err = %v, want a regular file", want, exist, fi, err)
			}
			if _, err := metadata.Load(file.Join(dir, metadata.AgentMetadataFileName)); err != nil {
				t.Errorf("metadata.Load() err = %v, want nil", err)
			}
		})
	}
}

// Test_faiss_load_missingIndexFile verifies that a directory with valid
// metadata and kvsdb but without the index file is reported as
// ErrIndexFileNotFound, so initFaiss takes its usual missing-index fallback.
// The directory used to be handed to faiss read_index instead.
func Test_faiss_load_missingIndexFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if err := metadata.Store(file.Join(dir, metadata.AgentMetadataFileName), &metadata.Metadata{
		Faiss: &metadata.Faiss{IndexCount: 1},
	}); err != nil {
		t.Fatal(err)
	}
	writeGob(t, file.Join(dir, kvsFileName), map[string]uint32{"a": 0})
	writeGob(t, file.Join(dir, kvsTimestampFileName), map[string]int64{"a": 1})

	eg, _ := errgroup.New(t.Context())
	f := &faiss{
		path:             dir,
		fmap:             make(map[string]int64),
		kvsdbConcurrency: 1,
		eg:               eg,
		minLit:           time.Minute,
		maxLit:           time.Minute,
		litFactor:        time.Millisecond,
	}
	defer closeKVS(t, f)

	if err := f.load(t.Context(), dir); !errors.Is(err, errors.ErrIndexFileNotFound) {
		t.Errorf("load() err = %v, want ErrIndexFileNotFound", err)
	}
}

// Test_moveAsideLegacyIndexFile verifies that a regular file left at the index
// path by an older agent, which wrote the faiss index to the path itself, is
// renamed next to the path so the index directory can be created and saved
// again; the file must be kept. Before, MkdirAll failed on it on every save.
func Test_moveAsideLegacyIndexFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		setup      func(t *testing.T, path string)
		name       string
		wantBackup bool
	}{
		{
			name:  "missing path is left alone",
			setup: func(*testing.T, string) {},
		},
		{
			name: "existing directory is left alone",
			setup: func(t *testing.T, path string) {
				t.Helper()
				if err := file.MkdirAll(path, os.ModePerm); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "legacy index file is moved aside and saving succeeds",
			setup: func(t *testing.T, path string) {
				t.Helper()
				if err := os.WriteFile(path, []byte("legacy"), os.ModePerm); err != nil {
					t.Fatal(err)
				}
			},
			wantBackup: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			parent := t.TempDir()
			path := file.Join(parent, "index")
			test.setup(t, path)

			if err := moveAsideLegacyIndexFile(path); err != nil {
				t.Fatalf("moveAsideLegacyIndexFile() err = %v, want nil", err)
			}

			entries, err := file.ReadDir(parent)
			if err != nil {
				t.Fatal(err)
			}
			var backups []string
			for _, e := range entries {
				if strings.HasPrefix(e.Name(), "index.legacy-") {
					backups = append(backups, e.Name())
				}
			}
			if !test.wantBackup {
				if len(backups) != 0 {
					t.Errorf("backups = %v, want none", backups)
				}
				return
			}
			if len(backups) != 1 {
				t.Fatalf("backups = %v, want exactly one", backups)
			}
			b, err := file.ReadFile(file.Join(parent, backups[0]))
			if err != nil || string(b) != "legacy" {
				t.Errorf("backup content = %q, err = %v, want %q", b, err, "legacy")
			}

			c := new(pathRecordingCore)
			f := &faiss{
				path: path,
				fmap: make(map[string]int64),
				kvs:  kvs.New(kvs.WithConcurrency(1)),
				core: c,
			}
			defer closeKVS(t, f)
			f.isTrained.Store(true)
			f.indexing.Store(false)
			f.saving.Store(false)
			f.nocie.Store(1)
			if err := f.saveIndex(t.Context()); err != nil {
				t.Errorf("saveIndex() after moving the legacy file aside err = %v, want nil", err)
			}
		})
	}
}

func Test_searchK(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		k    uint32
		n    uint64
		want int
	}{
		{name: "k below index size is kept", k: 5, n: 10, want: 5},
		{name: "k equal to index size is kept", k: 10, n: 10, want: 10},
		{name: "k above index size is capped", k: math.MaxUint32, n: 10, want: 10},
		{name: "empty index still searches with k=1", k: 100, n: 0, want: 1},
		{name: "k=0 stays 0 (rejected by core)", k: 0, n: 10, want: 0},
		{name: "huge index caps at MaxInt32", k: math.MaxUint32, n: math.MaxUint64, want: math.MaxInt32},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := searchK(tc.k, tc.n); got != tc.want {
				t.Errorf("searchK(%d, %d) = %d, want %d", tc.k, tc.n, got, tc.want)
			}
		})
	}
}

// searchKCore is a stubCore that records the k passed to Search and returns
// every stored id (FAISS never returns more neighbours than it holds).
type searchKCore struct {
	stubCore
	ks []int
	n  int
}

func (s *searchKCore) Search(k, _, _ int, _ []float32) ([]algorithm.SearchResult, error) {
	s.ks = append(s.ks, k)
	res := make([]algorithm.SearchResult, 0, min(k, s.n))
	for i := range min(k, s.n) {
		res = append(res, algorithm.SearchResult{ID: uint32(i), Distance: float32(i)})
	}
	return res, nil
}

// Test_faiss_Search_capsK verifies that a client-supplied k far above the
// index size reaches the core capped at the index size (the core allocates k
// slots up front), while the results are unchanged.
func Test_faiss_Search_capsK(t *testing.T) {
	t.Parallel()
	const n = 3
	sc := &searchKCore{n: n}
	f := &faiss{
		dim:  2,
		fmap: make(map[string]int64),
		kvs:  kvs.New(kvs.WithConcurrency(1)),
		core: sc,
		icnt: n,
	}
	defer closeKVS(t, f)
	f.indexing.Store(false)
	f.saving.Store(false)
	want := make([]string, 0, n)
	for i := range n {
		id := "uuid-" + strconv.Itoa(i)
		f.kvs.Set(id, uint32(i), int64(i+1))
		want = append(want, id)
	}

	res, err := f.Search(math.MaxUint32, 1, 1, []float32{0, 0})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if !slices.Equal(sc.ks, []int{n}) {
		t.Errorf("core k = %v, want [%d]", sc.ks, n)
	}
	got := make([]string, 0, len(res.GetResults()))
	for _, r := range res.GetResults() {
		got = append(got, r.GetId())
	}
	if !slices.Equal(got, want) {
		t.Errorf("ids = %v, want %v", got, want)
	}
}

// blockingTrainCore is a saveOKCore whose Train blocks until released and
// whose Add succeeds, so a CreateIndex run can be held right before it updates
// isTrained, trainSize and the add buffer.
type blockingTrainCore struct {
	entered chan struct{}
	release chan struct{}
	saveOKCore
}

func (c *blockingTrainCore) Train(int, []float32) error {
	close(c.entered)
	select {
	case <-c.release:
		return nil
	case <-time.After(30 * time.Second):
		return errors.New("Train was not released")
	}
}

func (*blockingTrainCore) Add(n int, _ []float32, _ []int64) (int, error) { return n, nil }

// Test_faiss_trainState_concurrentReaders verifies that the training state and
// the add buffer size, which CreateIndex updates under cimu, can be read by
// GetTrainSize, SaveIndex and the lock-free pre-check of another CreateIndex
// while that run is in progress. Run with -race (make test): the readers are
// not ordered with the writes made after Train returns.
func Test_faiss_trainState_concurrentReaders(t *testing.T) {
	t.Parallel()
	vq, err := vqueue.New()
	if err != nil {
		t.Fatal(err)
	}
	c := &blockingTrainCore{entered: make(chan struct{}), release: make(chan struct{})}
	f := &faiss{
		dim:              2,
		nlist:            1,
		m:                1,
		path:             t.TempDir(),
		fmap:             make(map[string]int64),
		kvs:              kvs.New(kvs.WithConcurrency(1)),
		vq:               vq,
		core:             c,
		kvsdbConcurrency: 1,
	}
	defer closeKVS(t, f)
	f.indexing.Store(false)
	f.saving.Store(false)
	// max(nlist, 2^m)*minPointsPerCentroid vectors are needed to train.
	n := 2 * minPointsPerCentroid
	for i := range n {
		if err := f.Insert("uuid-"+strconv.Itoa(i), []float32{float32(i), float32(i)}); err != nil {
			t.Fatal(err)
		}
	}
	// lastNocie != nocie, so saveIndex reads isTrained instead of returning early.
	f.nocie.Store(1)

	done := make(chan error, 1)
	go func() { done <- f.CreateIndex(t.Context()) }()
	select {
	case <-c.entered:
	case err := <-done:
		t.Fatalf("CreateIndex returned before Train: %v", err)
	case <-time.After(30 * time.Second):
		t.Fatal("CreateIndex did not reach Train")
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	close(c.release)
	if got := f.GetTrainSize(); got != 0 && got != n {
		t.Errorf("GetTrainSize() = %d during CreateIndex, want 0 or %d", got, n)
	}
	if err := f.CreateIndex(ctx); err != nil &&
		!errors.Is(err, context.Canceled) && !errors.Is(err, errors.ErrUncommittedIndexNotFound) {
		t.Errorf("concurrent CreateIndex() err = %v", err)
	}
	if err := f.saveIndex(ctx); err != nil && !errors.Is(err, context.Canceled) {
		t.Errorf("concurrent saveIndex() err = %v", err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("CreateIndex() err = %v", err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("CreateIndex did not finish")
	}
	if !f.isTrained.Load() {
		t.Error("isTrained = false after training, want true")
	}
	if got := f.GetTrainSize(); got != n {
		t.Errorf("GetTrainSize() = %d, want %d", got, n)
	}
	if got := f.addLen.Load(); got != 0 || len(f.addVecs) != 0 {
		t.Errorf("add buffer: addLen = %d, %d floats, want empty", got, len(f.addVecs))
	}
}

// update deletes at t and inserts at t+1. t == MaxInt64 used to wrap the
// insert to MinInt64, which the queue drops behind the MaxInt64 delete: the
// update reported success and the object disappeared. It must be rejected
// before the delete; ordinary timestamps keep the delete-at-t / insert-at-t+1.
func Test_faiss_update_maxTimestamp(t *testing.T) {
	t.Parallel()
	const (
		uuid = "update-ts-uuid"
		oldT = int64(100)
	)
	oldVec, newVec := []float32{9}, []float32{1}
	tests := []struct {
		name    string
		ts      int64
		wantErr bool
	}{
		{name: "rejects MaxInt64", ts: math.MaxInt64, wantErr: true},
		{name: "keeps a normal timestamp", ts: 200},
		{name: "keeps MaxInt64-1", ts: math.MaxInt64 - 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			vq, err := vqueue.New()
			if err != nil {
				t.Fatal(err)
			}
			if err = vq.PushInsert(uuid, oldVec, oldT); err != nil {
				t.Fatal(err)
			}
			f := &faiss{dim: 1, kvs: kvs.New(kvs.WithConcurrency(1)), vq: vq}
			defer closeKVS(t, f)

			err = f.update(uuid, newVec, test.ts)
			vec, its, dts, exists := vq.GetVectorWithTimestamp(uuid)
			if test.wantErr {
				if !errors.Is(err, errors.ErrInvalidTimestamp(test.ts)) {
					t.Fatalf("got error %v, want %v", err, errors.ErrInvalidTimestamp(test.ts))
				}
				if !exists || !slices.Equal(vec, oldVec) || its != oldT || dts != 0 {
					t.Errorf("rejected update changed the queue: vec=%v its=%d dts=%d exists=%t", vec, its, dts, exists)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !exists || !slices.Equal(vec, newVec) || dts != test.ts || its != test.ts+1 {
				t.Errorf("got vec=%v its=%d dts=%d exists=%t, want vec=%v its=%d dts=%d", vec, its, dts, exists, newVec, test.ts+1, test.ts)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func TestNew(t *testing.T) {
// 	type args struct {
// 		cfg  *config.Faiss
// 		opts []Option
// 	}
// 	type want struct {
// 		want Faiss
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, Faiss, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got Faiss, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           cfg:nil,
// 		           opts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           cfg:nil,
// 		           opts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := New(test.args.cfg, test.args.opts...)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_initFaiss(t *testing.T) {
// 	type args struct {
// 		opts []core.Option
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.initFaiss(test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_load(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		path string
// 		opts []core.Option
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           path:"",
// 		           opts:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           path:"",
// 		           opts:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.load(test.args.ctx, test.args.path, test.args.opts...)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_loadKVS(t *testing.T) {
// 	type args struct {
// 		ctx     context.Context
// 		path    string
// 		timeout time.Duration
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           path:"",
// 		           timeout:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           path:"",
// 		           timeout:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.loadKVS(test.args.ctx, test.args.path, test.args.timeout)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_mktmp(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.mktmp()
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Start(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want <-chan error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, <-chan error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got <-chan error) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.Start(test.args.ctx)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Train(t *testing.T) {
// 	type args struct {
// 		nb int
// 		xb []float32
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           nb:0,
// 		           xb:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           nb:0,
// 		           xb:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.Train(test.args.nb, test.args.xb)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Insert(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.Insert(test.args.uuid, test.args.vec)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_InsertWithTime(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 		t    int64
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.InsertWithTime(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_insert(t *testing.T) {
// 	type args struct {
// 		uuid       string
// 		xb         []float32
// 		t          int64
// 		validation bool
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           xb:nil,
// 		           t:0,
// 		           validation:false,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           xb:nil,
// 		           t:0,
// 		           validation:false,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.insert(test.args.uuid, test.args.xb, test.args.t, test.args.validation)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Update(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.Update(test.args.uuid, test.args.vec)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_UpdateWithTime(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 		t    int64
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.UpdateWithTime(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_update(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 		t    int64
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.update(test.args.uuid, test.args.vec, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_UpdateTimestamp(t *testing.T) {
// 	type args struct {
// 		uuid  string
// 		ts    int64
// 		force bool
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           ts:0,
// 		           force:false,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           ts:0,
// 		           force:false,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.UpdateTimestamp(test.args.uuid, test.args.ts, test.args.force)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_readyForUpdate(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		vec  []float32
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           vec:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.readyForUpdate(test.args.uuid, test.args.vec)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_CreateIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.CreateIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_SaveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.SaveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_saveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.saveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_moveAndSwitchSavedData(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.moveAndSwitchSavedData(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_CreateAndSaveIndex(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.CreateAndSaveIndex(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Search(t *testing.T) {
// 	type args struct {
// 		k      uint32
// 		nprobe uint32
// 		nq     uint32
// 		xq     []float32
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           k:0,
// 		           nprobe:0,
// 		           nq:0,
// 		           xq:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           k:0,
// 		           nprobe:0,
// 		           nq:0,
// 		           xq:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			gotRes, err := f.Search(test.args.k, test.args.nprobe, test.args.nq, test.args.xq)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Delete(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.Delete(test.args.uuid)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_DeleteWithTime(t *testing.T) {
// 	type args struct {
// 		uuid string
// 		t    int64
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           t:0,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           t:0,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.DeleteWithTime(test.args.uuid, test.args.t)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_delete(t *testing.T) {
// 	type args struct {
// 		uuid       string
// 		t          int64
// 		validation bool
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		           t:0,
// 		           validation:false,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           t:0,
// 		           validation:false,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.delete(test.args.uuid, test.args.t, test.args.validation)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Exists(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		wantOid uint32
// 		wantOk  bool
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint32, bool) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotOid uint32, gotOk bool) error {
// 		if !reflect.DeepEqual(gotOid, w.wantOid) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOid, w.wantOid)
// 		}
// 		if !reflect.DeepEqual(gotOk, w.wantOk) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotOk, w.wantOk)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			gotOid, gotOk := f.Exists(test.args.uuid)
// 			if err := checkFunc(test.want, gotOid, gotOk); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_GetObject(t *testing.T) {
// 	type args struct {
// 		uuid string
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		wantVec       []float32
// 		wantTimestamp int64
// 		err           error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []float32, int64, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotVec []float32, gotTimestamp int64, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotVec, w.wantVec) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotVec, w.wantVec)
// 		}
// 		if !reflect.DeepEqual(gotTimestamp, w.wantTimestamp) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotTimestamp, w.wantTimestamp)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           uuid:"",
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           uuid:"",
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			gotVec, gotTimestamp, err := f.GetObject(test.args.uuid)
// 			if err := checkFunc(test.want, gotVec, gotTimestamp, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_IsIndexing(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.IsIndexing()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_IsSaving(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want bool
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, bool) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got bool) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.IsSaving()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_UUIDs(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		wantUuids []string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, []string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotUuids []string) error {
// 		if !reflect.DeepEqual(gotUuids, w.wantUuids) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotUuids, w.wantUuids)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			gotUuids := f.UUIDs(test.args.ctx)
// 			if err := checkFunc(test.want, gotUuids); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_NumberOfCreateIndexExecution(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.NumberOfCreateIndexExecution()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_NumberOfProactiveGCExecution(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.NumberOfProactiveGCExecution()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_gc(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			f.gc()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Len(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.Len()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_InsertVQueueBufferLen(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.InsertVQueueBufferLen()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_DeleteVQueueBufferLen(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want uint64
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, uint64) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got uint64) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.DeleteVQueueBufferLen()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_GetDimensionSize(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.GetDimensionSize()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_GetTrainSize(t *testing.T) {
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		want int
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, int) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got int) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			got := f.GetTrainSize()
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_Close(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			err := f.Close(test.args.ctx)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_ListObjectFunc(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		fn  func(uuid string, oid uint32, ts int64) bool
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           fn:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           fn:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			f.ListObjectFunc(test.args.ctx, test.args.fn)
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_faiss_toSearchResponse(t *testing.T) {
// 	type args struct {
// 		sr []algorithm.SearchResult
// 	}
// 	type fields struct {
// 		tmpPath           atomic.Value
// 		eg                errgroup.Group
// 		kvs               kvs.BidiMap
// 		core              core.Faiss
// 		vq                vqueue.Queue
// 		saving            atomic.Value
// 		indexing          atomic.Value
// 		fmap              map[string]int64
// 		basePath          string
// 		oldPath           string
// 		path              string
// 		addIds            []int64
// 		addVecs           []float32
// 		nlist             int
// 		alen              int
// 		nocie             uint64
// 		nogce             uint64
// 		wfci              uint64
// 		kvsdbConcurrency  int
// 		dim               int
// 		lastNocie         uint64
// 		m                 int
// 		icnt              uint64
// 		dur               time.Duration
// 		sdur              time.Duration
// 		lim               time.Duration
// 		minLit            time.Duration
// 		maxLit            time.Duration
// 		litFactor         time.Duration
// 		idelay            time.Duration
// 		trainSize         int
// 		enableCopyOnWrite bool
// 		isTrained         bool
// 		dcd               bool
// 		enableProactiveGC bool
// 		inMem             bool
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           sr:nil,
// 		       },
// 		       fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           sr:nil,
// 		           },
// 		           fields: fields {
// 		           tmpPath:nil,
// 		           eg:nil,
// 		           kvs:nil,
// 		           core:nil,
// 		           vq:nil,
// 		           saving:nil,
// 		           indexing:nil,
// 		           fmap:nil,
// 		           basePath:"",
// 		           oldPath:"",
// 		           path:"",
// 		           addIds:nil,
// 		           addVecs:nil,
// 		           nlist:0,
// 		           alen:0,
// 		           nocie:0,
// 		           nogce:0,
// 		           wfci:0,
// 		           kvsdbConcurrency:0,
// 		           dim:0,
// 		           lastNocie:0,
// 		           m:0,
// 		           icnt:0,
// 		           dur:nil,
// 		           sdur:nil,
// 		           lim:nil,
// 		           minLit:nil,
// 		           maxLit:nil,
// 		           litFactor:nil,
// 		           idelay:nil,
// 		           trainSize:0,
// 		           enableCopyOnWrite:false,
// 		           isTrained:false,
// 		           dcd:false,
// 		           enableProactiveGC:false,
// 		           inMem:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			f := &faiss{
// 				tmpPath:           test.fields.tmpPath,
// 				eg:                test.fields.eg,
// 				kvs:               test.fields.kvs,
// 				core:              test.fields.core,
// 				vq:                test.fields.vq,
// 				saving:            test.fields.saving,
// 				indexing:          test.fields.indexing,
// 				fmap:              test.fields.fmap,
// 				basePath:          test.fields.basePath,
// 				oldPath:           test.fields.oldPath,
// 				path:              test.fields.path,
// 				addIds:            test.fields.addIds,
// 				addVecs:           test.fields.addVecs,
// 				nlist:             test.fields.nlist,
// 				alen:              test.fields.alen,
// 				nocie:             test.fields.nocie,
// 				nogce:             test.fields.nogce,
// 				wfci:              test.fields.wfci,
// 				kvsdbConcurrency:  test.fields.kvsdbConcurrency,
// 				dim:               test.fields.dim,
// 				lastNocie:         test.fields.lastNocie,
// 				m:                 test.fields.m,
// 				icnt:              test.fields.icnt,
// 				dur:               test.fields.dur,
// 				sdur:              test.fields.sdur,
// 				lim:               test.fields.lim,
// 				minLit:            test.fields.minLit,
// 				maxLit:            test.fields.maxLit,
// 				litFactor:         test.fields.litFactor,
// 				idelay:            test.fields.idelay,
// 				trainSize:         test.fields.trainSize,
// 				enableCopyOnWrite: test.fields.enableCopyOnWrite,
// 				isTrained:         test.fields.isTrained,
// 				dcd:               test.fields.dcd,
// 				enableProactiveGC: test.fields.enableProactiveGC,
// 				inMem:             test.fields.inMem,
// 			}
//
// 			gotRes, err := f.toSearchResponse(test.args.sr)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
