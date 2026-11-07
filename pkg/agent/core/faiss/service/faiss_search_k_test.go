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
	"math"
	"slices"
	"strconv"
	"testing"

	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
)

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
