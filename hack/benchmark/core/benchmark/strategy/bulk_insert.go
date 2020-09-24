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

// Package strategy provides benchmark strategy
package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

const (
	maxBulkSize = 100000
)

func NewBulkInsert(opts ...StrategyOption) benchmark.Strategy {
	return newStrategy(append([]StrategyOption{
		WithPropName("BulkInsert"),
		WithProp32(
			func(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				size := func() int {
					if maxBulkSize < dataset.TrainSize() {
						return maxBulkSize
					} else {
						return dataset.TrainSize()
					}
				}()
				v := make([][]float32, 0, size)
				for i := 0; i < size; i++ {
					arr, err := dataset.Train(i)
					if err != nil {
						break
					}
					v = append(v, arr.([]float32))
				}
				b.StartTimer()
				defer b.StopTimer()
				ids, errs := c.BulkInsert(v)
				return ids, wrapErrors(errs)
			},
		),
		WithProp64(
			func(ctx context.Context, b *testing.B, c core.Core64, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				size := func() int {
					if maxBulkSize < dataset.TrainSize() {
						return maxBulkSize
					} else {
						return dataset.TrainSize()
					}
				}()
				v := make([][]float64, 0, size)
				for i := 0; i < size; i++ {
					arr, err := dataset.Train(i)
					if err != nil {
						break
					}
					v = append(v, float32To64(arr.([]float32)))
				}
				b.StartTimer()
				defer b.StopTimer()
				ids, errs := c.BulkInsert(v)
				return ids, wrapErrors(errs)
			},
		),
	}, opts...)...)
}
