//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"github.com/vdaas/vald/internal/errors"
)

type strategy struct {
	// 32-bit algorithm core.
	core32 algorithm.Bit32
	// 64-bit algorithm core.
	core64 algorithm.Bit64
	// Closer for the algorithm core.
	closer algorithm.Closer
	// Initialization function for the 32-bit algorithm core.
	initBit32 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit32, algorithm.Closer, error)
	// Initialization function for the 64-bit algorithm core.
	initBit64 func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit64, algorithm.Closer, error)
	// Pre-property function for the 32-bit algorithm core.
	preProp32 func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset) ([]uint, error)
	// Pre-property function for the 64-bit algorithm core.
	preProp64 func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset) ([]uint, error)
	// Property function for the 32-bit algorithm core.
	prop32 func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset, []uint, *uint64) (any, error)
	// Property function for the 64-bit algorithm core.
	prop64 func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset, []uint, *uint64) (any, error)
	// Name of the property.
	propName string
	// Algorithm mode.
	mode algorithm.Mode
	// Whether to run in parallel.
	parallel bool
}

func newStrategy(opts ...StrategyOption) benchmark.Strategy {
	s := &strategy{
		// invalid mode.
		mode: algorithm.Mode(100),
	}
	for _, opt := range append(defaultStrategyOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *strategy) Init(ctx context.Context, b *testing.B, dataset assets.Dataset) error {
	b.Helper()
	switch s.mode {
	case algorithm.Float32:
		core32, closer, err := s.initBit32(ctx, b, dataset)
		if err != nil {
			b.Error(err)
			return err
		}
		s.core32, s.closer = core32, closer
	case algorithm.Float64:
		core64, closer, err := s.initBit64(ctx, b, dataset)
		if err != nil {
			b.Error(err)
			return err
		}
		s.core64, s.closer = core64, closer
	default:
		b.Error(errors.ErrInvalidCoreMode)
		return errors.ErrInvalidCoreMode
	}
	return nil
}

func (s *strategy) PreProp(
	ctx context.Context, b *testing.B, dataset assets.Dataset,
) ([]uint, error) {
	b.Helper()

	switch s.mode {
	case algorithm.Float32:
		return s.preProp32(ctx, b, s.core32, dataset)
	case algorithm.Float64:
		return s.preProp64(ctx, b, s.core64, dataset)
	default:
		return nil, errors.ErrInvalidCoreMode
	}
}

func (s *strategy) Run(ctx context.Context, b *testing.B, dataset assets.Dataset, ids []uint) {
	b.Helper()

	var cnt uint64

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	defer b.StopTimer()

	switch s.mode {
	case algorithm.Float32:
		b.Run(s.propName, func(bb *testing.B) {
			s.float32(ctx, bb, dataset, ids, &cnt)
		})
	case algorithm.Float64:
		b.Run(s.propName, func(bb *testing.B) {
			s.float64(ctx, bb, dataset, ids, &cnt)
		})
	default:
		b.Fatal(errors.ErrInvalidCoreMode)
	}
}

func (s *strategy) Close() {
	s.closer.Close()
}

func (s *strategy) run(b *testing.B, f func()) {
	b.Helper()

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()

	if s.parallel {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				f()
			}
		})
	} else {
		for i := 0; i < b.N; i++ {
			f()
		}
	}

	b.StopTimer()
}

func (s *strategy) float32(
	ctx context.Context, b *testing.B, dataset assets.Dataset, ids []uint, cnt *uint64,
) {
	s.run(b, func() {
		_, err := s.prop32(ctx, b, s.core32, dataset, ids, cnt)
		if err != nil {
			b.Error(err)
		}
		atomic.AddUint64(cnt, 1)
	})
}

func (s *strategy) float64(
	ctx context.Context, b *testing.B, dataset assets.Dataset, ids []uint, cnt *uint64,
) {
	s.run(b, func() {
		_, err := s.prop64(ctx, b, s.core64, dataset, ids, cnt)
		if err != nil {
			b.Error(err)
		}
		atomic.AddUint64(cnt, 1)
	})
}
