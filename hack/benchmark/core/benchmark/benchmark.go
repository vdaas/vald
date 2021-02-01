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

// Package benchmark provides benchmark frame
package benchmark

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
)

// Benchmark is an interface for NGT benchmark.
type Benchmark interface {
	Run(context.Context, *testing.B)
}

type benchmark struct {
	name       string
	dataset    assets.Dataset
	strategies []Strategy
}

func New(b *testing.B, opts ...Option) Benchmark {
	bm := new(benchmark)
	for _, opt := range append(defaultOptions, opts...) {
		opt(bm)
	}

	fn := assets.Data(bm.name)
	if fn == nil {
		b.Fatalf("dataset provider is nil: %v", bm.name)
	}

	bm.dataset = fn(b)
	if bm.dataset == nil {
		b.Fatalf("dataset is nil: %v", bm.name)
	}

	return bm
}

func (bm *benchmark) Run(ctx context.Context, b *testing.B) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	b.Run(bm.name, func(bb *testing.B) {
		for _, strategy := range bm.strategies {
			err := func() error {
				bb.Helper()

				err := strategy.Init(ctx, bb, bm.dataset)
				if err != nil {
					return err
				}
				defer strategy.Close()

				obj, err := strategy.PreProp(ctx, bb, bm.dataset)
				if err != nil {
					return err
				}

				strategy.Run(ctx, bb, bm.dataset, obj)
				return nil
			}()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.StopTimer()
}
