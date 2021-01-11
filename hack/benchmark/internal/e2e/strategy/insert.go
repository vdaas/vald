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

type insert struct {
	parallel bool
}

func NewInsert(opts ...InsertOption) e2e.Strategy {
	i := new(insert)
	for _, opt := range append(defaultInsertOption, opts...) {
		opt(i)
	}
	return i
}

func (isrt *insert) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	if isrt.parallel {
		isrt.runParallel(ctx, b, c, dataset)
		return
	}
	isrt.run(ctx, b, c, dataset)
}

func (isrt *insert) run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	cnt := 0
	b.Run("Insert", func(bb *testing.B) {
		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		for i := 0; i < bb.N; i++ {
			v, err := dataset.Train(cnt % dataset.TrainSize())
			if err != nil {
				b.Fatal(err)
			}
			isrt.do(ctx, bb, c, fmt.Sprint(cnt), v.([]float32))
			cnt++
		}
		bb.StopTimer()
	})
}

func (isrt *insert) runParallel(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	var cnt int64
	b.Run("ParallelInsert", func(bb *testing.B) {
		bb.StopTimer()
		bb.ReportAllocs()
		bb.ResetTimer()
		bb.StartTimer()
		bb.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				n := int(atomic.AddInt64(&cnt, 1)) - 1
				v, err := dataset.Train(n % dataset.TrainSize())
				if err != nil {
					b.Fatal(err)
				}

				isrt.do(ctx, bb, c, fmt.Sprint(cnt), v.([]float32))
			}
		})
		bb.StopTimer()
	})
}

func (isrt *insert) do(ctx context.Context, b *testing.B, c client.Client, id string, vector []float32) {
	if _, err := c.Insert(ctx, &client.InsertRequest{
		Vector: &client.ObjectVector{
			Id:     id,
			Vector: vector,
		},
	}); err != nil {
		b.Error(err)
	}
}
