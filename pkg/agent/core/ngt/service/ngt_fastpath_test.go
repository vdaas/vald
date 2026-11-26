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

package service_test

import (
	"context"
	"math"
	"runtime/debug"
	"slices"
	"testing"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
)

func Test_ngt_GetTimestamp_And_InsertWithConfig(t *testing.T) {
	t.Parallel()

	c := &config.NGT{
		Dimension:              4,
		DistanceType:           ngt.L2.String(),
		ObjectType:             ngt.Float.String(),
		CreationEdgeSize:       20,
		SearchEdgeSize:         10,
		ErrorBufferLimit:       100,
		EnableInMemoryMode:     true,
		AutoIndexDurationLimit: benchNoAuto,
		AutoIndexCheckDuration: benchNoAuto,
		AutoSaveIndexDuration:  benchNoAuto,
		AutoIndexLength:        10000000000,
		KVSDB:                  &config.KVSDB{Concurrency: 1},
	}
	n, err := service.New(c.Bind())
	if err != nil {
		t.Fatalf("service.New() error = %v", err)
	}
	t.Cleanup(func() {
		if cerr := n.Close(context.WithoutCancel(t.Context())); cerr != nil {
			t.Errorf("Close() error = %v", cerr)
		}
	})

	const (
		uuid1 = "id-1"
		uuid2 = "id-2"
		ts1   = int64(1000)
		ts2   = int64(2000)
		ts3   = int64(3000)
	)
	v1 := []float32{1, 2, 3, 4}
	v2 := []float32{5, 6, 7, 8}

	// Missing uuid returns ErrObjectIDNotFound.
	if _, gerr := n.GetTimestamp("missing"); !errors.Is(gerr, errors.ErrObjectIDNotFound("missing")) {
		t.Fatalf("GetTimestamp(missing) err = %v, want ErrObjectIDNotFound", gerr)
	}

	// Strict insert creates entry in vqueue.
	if ierr := n.InsertWithConfig(uuid1, v1, ts1, false); ierr != nil {
		t.Fatalf("InsertWithConfig(false) error = %v", ierr)
	}
	gotTs, err := n.GetTimestamp(uuid1)
	if err != nil || gotTs != ts1 {
		t.Fatalf("GetTimestamp(vqueue) = (%d, %v), want (%d, nil)", gotTs, err, ts1)
	}

	// Strict insert rejects duplicate id.
	if ierr := n.InsertWithConfig(uuid1, v2, ts2, false); !errors.Is(ierr, errors.ErrUUIDAlreadyExists(uuid1)) {
		t.Fatalf("duplicate InsertWithConfig(false) err = %v, want ErrUUIDAlreadyExists", ierr)
	}

	// SkipStrictExistCheck=true overwrites in vqueue cleanly.
	if ierr := n.InsertWithConfig(uuid1, v2, ts2, true); ierr != nil {
		t.Fatalf("InsertWithConfig(true) error = %v", ierr)
	}
	gotTs, err = n.GetTimestamp(uuid1)
	if err != nil || gotTs != ts2 {
		t.Fatalf("GetTimestamp(after overwrite) = (%d, %v), want (%d, nil)", gotTs, err, ts2)
	}

	// Index and verify GetTimestamp from kvs (vectorless).
	if cerr := n.CreateIndex(t.Context(), 1); cerr != nil {
		t.Fatalf("CreateIndex() error = %v", cerr)
	}
	gotTs, err = n.GetTimestamp(uuid1)
	if err != nil || gotTs != ts2 {
		t.Fatalf("GetTimestamp(indexed) = (%d, %v), want (%d, nil)", gotTs, err, ts2)
	}

	// InsertMultipleWithConfig(true) succeeds on existing id and updates after CreateIndex.
	if merr := n.InsertMultipleWithConfig(map[string][]float32{
		uuid1: v1,
		uuid2: v2,
	}, ts3, true); merr != nil {
		t.Fatalf("InsertMultipleWithConfig(true) error = %v", merr)
	}
	if cerr := n.CreateIndex(t.Context(), 2); cerr != nil {
		t.Fatalf("CreateIndex(2) error = %v", cerr)
	}
	gotTs, err = n.GetTimestamp(uuid1)
	if err != nil || gotTs != ts3 {
		t.Fatalf("GetTimestamp(after multi overwrite) = (%d, %v), want (%d, nil)", gotTs, err, ts3)
	}
	if n.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", n.Len())
	}
}

func Test_ngt_ValidatedAndPooledSearchByID(t *testing.T) {
	t.Parallel()

	c := &config.NGT{
		Dimension:              4,
		DistanceType:           ngt.L2.String(),
		ObjectType:             ngt.Float.String(),
		CreationEdgeSize:       20,
		SearchEdgeSize:         10,
		ErrorBufferLimit:       100,
		EnableInMemoryMode:     true,
		AutoIndexDurationLimit: benchNoAuto,
		AutoIndexCheckDuration: benchNoAuto,
		AutoSaveIndexDuration:  benchNoAuto,
		AutoIndexLength:        10000000000,
		KVSDB:                  &config.KVSDB{Concurrency: 1},
	}
	n, err := service.New(c.Bind())
	if err != nil {
		t.Fatalf("service.New() error = %v", err)
	}
	t.Cleanup(func() {
		if cerr := n.Close(context.WithoutCancel(t.Context())); cerr != nil {
			t.Errorf("Close() error = %v", cerr)
		}
	})

	v1 := []float32{1, 2, 3, 4}
	v2 := []float32{5, 6, 7, 8}
	vqVec := []float32{9, 10, 11, 12}
	if err = n.Insert("id-1", v1); err != nil {
		t.Fatalf("Insert(id-1) error = %v", err)
	}
	if err = n.Insert("id-2", v2); err != nil {
		t.Fatalf("Insert(id-2) error = %v", err)
	}
	if err = n.CreateIndex(t.Context(), 2); err != nil {
		t.Fatalf("CreateIndex() error = %v", err)
	}
	// Insert a third vector that remains in vqueue to verify vqueue slice is never pooled or mutated.
	if err = n.Insert("id-vq", vqVec); err != nil {
		t.Fatalf("Insert(id-vq) error = %v", err)
	}

	// SearchValidated and LinearSearchValidated match Search and LinearSearch.
	s1, err := n.Search(t.Context(), v1, 2, 0.1, -1, 0)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	s2, err := n.SearchValidated(t.Context(), v1, 2, 0.1, -1, 0)
	if err != nil {
		t.Fatalf("SearchValidated() error = %v", err)
	}
	if len(s1.GetResults()) != len(s2.GetResults()) {
		t.Fatalf("SearchValidated() len = %d, want %d", len(s2.GetResults()), len(s1.GetResults()))
	}

	l1, err := n.LinearSearch(t.Context(), v1, 2)
	if err != nil {
		t.Fatalf("LinearSearch() error = %v", err)
	}
	l2, err := n.LinearSearchValidated(t.Context(), v1, 2)
	if err != nil {
		t.Fatalf("LinearSearchValidated() error = %v", err)
	}
	if len(l1.GetResults()) != len(l2.GetResults()) || l2.GetResults()[0].GetId() != "id-1" {
		t.Fatalf("LinearSearchValidated() = %+v, want id-1 first matching %+v", l2.GetResults(), l1.GetResults())
	}

	// Missing UUID returns ErrObjectIDNotFound cleanly.
	if _, _, err = n.SearchByID(t.Context(), "missing", 2, 0.1, -1, 0); !errors.Is(err, errors.ErrObjectIDNotFound("missing")) {
		t.Errorf("SearchByID(missing) err = %v, want ErrObjectIDNotFound", err)
	}
	if _, _, err = n.LinearSearchByID(t.Context(), "missing", 2); !errors.Is(err, errors.ErrObjectIDNotFound("missing")) {
		t.Errorf("LinearSearchByID(missing) err = %v, want ErrObjectIDNotFound", err)
	}

	// Stored NaN vector in vqueue must be rejected by SearchByID / LinearSearchByID with ErrNaNOrInfVector.
	nanVec := []float32{float32(math.NaN()), 1, 2, 3}
	if err = n.Insert("id-nan", nanVec); err != nil {
		t.Fatalf("Insert(id-nan) error = %v", err)
	}
	if _, _, err = n.SearchByID(t.Context(), "id-nan", 2, 0.1, -1, 0); !errors.Is(err, errors.ErrNaNOrInfVector) {
		t.Errorf("SearchByID(id-nan) err = %v, want ErrNaNOrInfVector", err)
	}
	if _, _, err = n.LinearSearchByID(t.Context(), "id-nan", 2); !errors.Is(err, errors.ErrNaNOrInfVector) {
		t.Errorf("LinearSearchByID(id-nan) err = %v, want ErrNaNOrInfVector", err)
	}
	if err = n.Delete("id-nan"); err != nil {
		t.Fatalf("Delete(id-nan) error = %v", err)
	}

	// vqueue hit must never alias or mutate vqVec even after subsequent indexed SearchByID calls.
	if gotVec, res, serr := n.SearchByID(t.Context(), "id-vq", 2, 0.1, -1, 0); serr != nil || gotVec != nil || len(res.GetResults()) == 0 {
		t.Fatalf("SearchByID(id-vq) = (%v, %+v, %v)", gotVec, res, serr)
	}
	if gotVec, res, serr := n.LinearSearchByID(t.Context(), "id-vq", 2); serr != nil || gotVec != nil || len(res.GetResults()) == 0 {
		t.Fatalf("LinearSearchByID(id-vq) = (%v, %+v, %v)", gotVec, res, serr)
	}

	// Concurrent indexed + vqueue SearchByID & LinearSearchByID must be race-free and keep vqVec intact.
	var wg sync.WaitGroup
	for range 16 {
		wg.Go(func() {
			for range 10 {
				if _, res, gerr := n.SearchByID(t.Context(), "id-1", 2, 0.1, -1, 0); gerr != nil || len(res.GetResults()) == 0 {
					t.Errorf("concurrent SearchByID(id-1) = (%+v, %v)", res, gerr)
				}
				if _, res, gerr := n.LinearSearchByID(t.Context(), "id-2", 2); gerr != nil || len(res.GetResults()) == 0 {
					t.Errorf("concurrent LinearSearchByID(id-2) = (%+v, %v)", res, gerr)
				}
				if _, res, gerr := n.LinearSearchByID(t.Context(), "id-vq", 2); gerr != nil || len(res.GetResults()) == 0 {
					t.Errorf("concurrent LinearSearchByID(id-vq) = (%+v, %v)", res, gerr)
				}
			}
		})
	}
	wg.Wait()

	storedVQ, _, err := n.GetObject("id-vq")
	if err != nil || !slices.Equal(storedVQ, vqVec) {
		t.Fatalf("GetObject(id-vq) = (%v, %v), want %v (must not be mutated by pool)", storedVQ, err, vqVec)
	}
}

func Test_ngt_ValidatedAndPooledSearchByID_Alloc(t *testing.T) {
	// Not parallel: testing.AllocsPerRun panics if any t.Parallel() test is active.
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			if s.Key == "-race" && s.Value == "true" {
				t.Skip("skipping AllocsPerRun under -race")
			}
		}
	}

	c := &config.NGT{
		Dimension:              4,
		DistanceType:           ngt.L2.String(),
		ObjectType:             ngt.Float.String(),
		CreationEdgeSize:       20,
		SearchEdgeSize:         10,
		ErrorBufferLimit:       100,
		EnableInMemoryMode:     true,
		AutoIndexDurationLimit: benchNoAuto,
		AutoIndexCheckDuration: benchNoAuto,
		AutoSaveIndexDuration:  benchNoAuto,
		AutoIndexLength:        10000000000,
		KVSDB:                  &config.KVSDB{Concurrency: 1},
	}
	n, err := service.New(c.Bind())
	if err != nil {
		t.Fatalf("service.New() error = %v", err)
	}
	t.Cleanup(func() {
		if cerr := n.Close(context.WithoutCancel(t.Context())); cerr != nil {
			t.Errorf("Close() error = %v", cerr)
		}
	})

	if err = n.Insert("id-1", []float32{1, 2, 3, 4}); err != nil {
		t.Fatalf("Insert(id-1) error = %v", err)
	}
	if err = n.Insert("id-2", []float32{5, 6, 7, 8}); err != nil {
		t.Fatalf("Insert(id-2) error = %v", err)
	}
	ctx := t.Context()
	if err = n.CreateIndex(ctx, 2); err != nil {
		t.Fatalf("CreateIndex() error = %v", err)
	}

	// Warm pool and verify steady-state SearchByID / LinearSearchByID allocate no extra vector buffer
	// beyond the Search and response slice allocations (<= 4 allocs).
	if _, _, err = n.SearchByID(ctx, "id-1", 2, 0.1, -1, 0); err != nil {
		t.Fatalf("warm SearchByID() error = %v", err)
	}
	if allocs := testing.AllocsPerRun(50, func() {
		_, res, aerr := n.SearchByID(ctx, "id-1", 2, 0.1, -1, 0)
		if aerr != nil || len(res.GetResults()) == 0 {
			t.Fatalf("SearchByID() = (%+v, %v)", res, aerr)
		}
	}); allocs > 4 {
		t.Errorf("SearchByID() allocs = %v, want <= 4", allocs)
	}
	if allocs := testing.AllocsPerRun(50, func() {
		_, res, aerr := n.LinearSearchByID(ctx, "id-1", 2)
		if aerr != nil || len(res.GetResults()) == 0 {
			t.Fatalf("LinearSearchByID() = (%+v, %v)", res, aerr)
		}
	}); allocs > 4 {
		t.Errorf("LinearSearchByID() allocs = %v, want <= 4", allocs)
	}
}
