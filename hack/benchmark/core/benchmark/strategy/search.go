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

// Package strategy provides benchmark strategy
package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/algorithm"
)

func NewSearch(size int, epsilon, radius float32, opts ...StrategyOption) benchmark.Strategy {
	return newStrategy(append([]StrategyOption{
		WithPropName("Search"),
		WithPreProp32(
			func(ctx context.Context, b *testing.B, c algorithm.Bit32, dataset assets.Dataset) (ids []uint, err error) {
				return insertAndCreateIndex32(ctx, c, dataset)
			},
		),
		WithProp32(
			func(ctx context.Context, b *testing.B, c algorithm.Bit32, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				v, err := dataset.Query(int(atomic.LoadUint64(cnt)) % dataset.TrainSize())
				if err != nil {
					return nil, err
				}
				b.StartTimer()
				defer b.StopTimer()
				return c.Search(v.([]float32), size, epsilon, radius)
			},
		),
		WithPreProp64(
			func(ctx context.Context, b *testing.B, c algorithm.Bit64, dataset assets.Dataset) (ids []uint, err error) {
				return insertAndCreateIndex64(ctx, c, dataset)
			},
		),
		WithProp64(
			func(ctx context.Context, b *testing.B, c algorithm.Bit64, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				v, err := dataset.Train(int(atomic.LoadUint64(cnt)) % dataset.TrainSize())
				if err != nil {
					return nil, err
				}
				b.StartTimer()
				defer b.StopTimer()
				return c.Search(float32To64(v.([]float32)), size, epsilon, radius)
			},
		),
	}, opts...)...)
}
