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

package ngt

import (
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/singleflight"
)

// newClosedNGT returns an ngt in the state Close leaves it in: mutexes set,
// index/ospace NULL. No call in these tests may hand the NULL handles to NGT.
func newClosedNGT(t *testing.T, ot objectType) *ngt {
	t.Helper()
	return &ngt{
		idxPath:    t.TempDir(),
		dimension:  9,
		objectType: ot,
		group:      singleflight.New[*GraphStatistics](),
		mu:         new(sync.RWMutex),
		cmu:        new(sync.RWMutex),
		smu:        new(sync.Mutex),
	}
}

// ngt_save_index, ngt_get_property_info and ngt_get_object_as_* / ngt_get_allocated_object_as_float
// dereference their handle without a NULL check, so the Go side must refuse a closed index.
func TestNGT_closedIndexGuards(t *testing.T) {
	t.Parallel()
	tests := []struct {
		call func(n *ngt) error
		name string
		ot   objectType
	}{
		{
			name: "SaveIndex",
			ot:   Float,
			call: func(n *ngt) error { return n.SaveIndex() },
		},
		{
			name: "SaveIndexWithPath",
			ot:   Float,
			call: func(n *ngt) error { return n.SaveIndexWithPath(n.idxPath) },
		},
		{
			name: "GetVector Float",
			ot:   Float,
			call: func(n *ngt) error { _, err := n.GetVector(1); return err },
		},
		{
			name: "GetVector HalfFloat",
			ot:   HalfFloat,
			call: func(n *ngt) error { _, err := n.GetVector(1); return err },
		},
		{
			name: "GetVector Uint8",
			ot:   Uint8,
			call: func(n *ngt) error { _, err := n.GetVector(1); return err },
		},
		{
			name: "GetProperty",
			ot:   Float,
			call: func(n *ngt) error { _, err := n.GetProperty(); return err },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			n := newClosedNGT(t, test.ot)
			if err := test.call(n); !errors.Is(err, errIndexClosed) {
				t.Errorf("error = %v, want %v", err, errIndexClosed)
			}
		})
	}
}

func TestNGT_CloseIdempotent(t *testing.T) {
	t.Parallel()
	// A hand-built ngt without mutexes was never opened; Close must be a no-op.
	new(ngt).Close()

	n := newClosedNGT(t, Float)
	n.Close()
	n.Close()
	n.CloseWithoutSaveIndex()
	if n.index != nil || n.ospace != nil {
		t.Error("index and ospace must be nil after Close")
	}
}

const (
	// blockedFor is how long a call must stay blocked to count as excluded by the held lock.
	blockedFor = 100 * time.Millisecond
	// releaseTimeout bounds the wait for the call to finish once the lock is released.
	releaseTimeout = 5 * time.Second
)

// Close must not free the index while a reader (Search/LinearSearch) holds rLock(true).
func TestNGT_CloseWaitsForInFlightReader(t *testing.T) {
	t.Parallel()
	n := newClosedNGT(t, Float)

	n.rLock(true)
	done := make(chan struct{})
	go func() {
		defer close(done)
		n.Close()
	}()

	select {
	case <-done:
		n.rUnlock(true)
		t.Fatal("Close returned while a reader held the index lock")
	case <-time.After(blockedFor):
	}
	n.rUnlock(true)

	select {
	case <-done:
	case <-time.After(releaseTimeout):
		t.Fatal("Close did not return after the reader released the lock")
	}
}

// The statistics walk must not run while a writer (Insert/Remove) holds lock(true).
func TestNGT_GraphStatisticsWaitsForWriter(t *testing.T) {
	t.Parallel()
	n := newClosedNGT(t, Float)

	n.lock(true)
	errCh := make(chan error, 1)
	go func() {
		_, err := n.GetGraphStatistics(t.Context(), NormalStatistics)
		errCh <- err
	}()

	select {
	case <-errCh:
		n.unlock(true)
		t.Fatal("GetGraphStatistics ran while a writer held the index lock")
	case <-time.After(blockedFor):
	}
	n.unlock(true)

	select {
	case err := <-errCh:
		// NGT reports the NULL index as an error.
		if err == nil {
			t.Error("GetGraphStatistics on a closed index must return an error")
		}
	case <-time.After(releaseTimeout):
		t.Fatal("GetGraphStatistics did not return after the writer released the lock")
	}
}
