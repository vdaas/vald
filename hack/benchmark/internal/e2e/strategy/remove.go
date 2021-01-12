//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/internal/client/v1/client"
)

type remove struct {
	parallel bool
}

func NewRemove(opts ...RemoveOption) e2e.Strategy {
	r := new(remove)
	for _, opt := range append(defaultRemoveOptions, opts...) {
		opt(r)
	}
	return r
}

func (r *remove) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	if r.parallel {
		r.runParallel(ctx, b, c, dataset)
		return
	}
	r.run(ctx, b, c, dataset)
}

func (r *remove) run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	cnt := 0
	b.Run("Remove", func(bb *testing.B) {
		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			r.do(ctx, bb, c, fmt.Sprint(cnt))
			cnt++
		}
		bb.StopTimer()
	})
}

func (r *remove) runParallel(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	var cnt int64
	b.Run("ParallelRemove", func(bb *testing.B) {
		bb.StartTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				n := int(atomic.AddInt64(&cnt, 1)) - 1
				r.do(ctx, bb, c, fmt.Sprint(n))
			}
		})
		bb.StopTimer()
	})
}

func (r *remove) do(ctx context.Context, b *testing.B, c client.Client, id string) {
	if _, err := c.Remove(ctx, &client.RemoveRequest{
		Id: &client.ObjectID{
			Id: id,
		},
	}); err != nil {
		b.Error(err)
	}
}
