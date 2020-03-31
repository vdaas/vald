//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package strategy provides strategy for e2e testing functions
package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/internal/client"
)

type search struct {
	parallel bool
	cfg      *client.SearchConfig
}

func NewSearch(opts ...SearchOption) e2e.Strategy {
	s := new(search)
	for _, opt := range append(defaultSearchOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *search) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	if s.parallel {
		s.runParallel(ctx, b, c, dataset)
		return
	}
	s.run(ctx, b, c, dataset)
}

func (s *search) run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("Search", func(bb *testing.B) {
		queries := dataset.Query()

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			s.do(ctx, bb, c, queries[i%len(queries)])
		}
		bb.StopTimer()
	})
}

func (s *search) runParallel(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	var cnt int64
	b.Run("ParallelSearch", func(bb *testing.B) {
		queries := dataset.Query()

		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				n := int(atomic.AddInt64(&cnt, 1)) - 1
				s.do(ctx, b, c, queries[n%len(queries)])
			}
		})
		bb.StopTimer()
	})
}

func (s *search) do(ctx context.Context, b *testing.B, c client.Client, query []float32) {
	if _, err := c.Search(ctx, &client.SearchRequest{
		Vector: query,
		Config: s.cfg,
	}); err != nil {
		b.Error(err)
	}
}
