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
	"strconv"
	"testing"

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/test/data/vector"
)

// BenchmarkSearch measures one in-memory NGT Search (20k x 128 float, L2,
// edge size 50) for k = 10 and k = 100, covering the per-call cgo work around
// ngt_search_index_as_float: query pinning, the C results object and the
// conversion of every result back to Go.
func BenchmarkSearch(b *testing.B) {
	log.Init(log.WithLoggerType("nop"))
	const (
		dim = 128
		num = 20000
	)
	n, err := New(
		WithDimension(dim),
		WithObjectType(Float),
		WithDistanceType(L2),
		WithCreationEdgeSize(50),
		WithSearchEdgeSize(50),
		WithDefaultPoolSize(16),
	)
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(n.Close)
	vecs := vector.GaussianDistributedFloat32VectorGenerator(num, dim)
	for _, v := range vecs {
		if _, err := n.Insert(v); err != nil {
			b.Fatal(err)
		}
	}
	if err := n.CreateIndex(16); err != nil {
		b.Fatal(err)
	}
	queries := vector.GaussianDistributedFloat32VectorGenerator(1000, dim)
	for _, k := range []int{10, 100} {
		b.Run("k="+strconv.Itoa(k), func(b *testing.B) {
			b.ResetTimer()
			for i := range b.N {
				res, err := n.Search(b.Context(), queries[i%len(queries)], k, 0, 0, 50)
				if err != nil {
					b.Fatal(err)
				}
				if len(res) != k {
					b.Fatalf("Search() returned %d results, want %d", len(res), k)
				}
			}
		})
	}
}
