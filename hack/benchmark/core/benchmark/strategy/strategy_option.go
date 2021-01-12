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
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/algorithm"
)

type StrategyOption func(*strategy) error

var defaultStrategyOptions = []StrategyOption{
	WithPreProp32(func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset) ([]uint, error) {
		return nil, nil
	}),
	WithPreProp64(func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset) ([]uint, error) {
		return nil, nil
	}),
}

func WithPreProp32(
	fn func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset) ([]uint, error),
) StrategyOption {
	return func(s *strategy) error {
		if fn != nil {
			s.preProp32 = fn
		}
		return nil
	}
}

func WithProp32(
	fn func(context.Context, *testing.B, algorithm.Bit32, assets.Dataset, []uint, *uint64) (interface{}, error),
) StrategyOption {
	return func(s *strategy) error {
		if fn != nil {
			s.prop32 = fn
		}
		return nil
	}
}

func WithPreProp64(
	fn func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset) ([]uint, error),
) StrategyOption {
	return func(s *strategy) error {
		if fn != nil {
			s.preProp64 = fn
		}
		return nil
	}
}

func WithProp64(
	fn func(context.Context, *testing.B, algorithm.Bit64, assets.Dataset, []uint, *uint64) (interface{}, error),
) StrategyOption {
	return func(s *strategy) error {
		if fn != nil {
			s.prop64 = fn
		}
		return nil
	}
}

func WithPropName(str string) StrategyOption {
	return func(s *strategy) error {
		if len(str) != 0 {
			s.propName = str
		}
		return nil
	}
}

func WithBit32(
	fn func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit32, algorithm.Closer, error),
) StrategyOption {
	return func(s *strategy) (err error) {
		if fn != nil {
			s.mode = algorithm.Float32
			s.initBit32 = fn
		}
		return nil
	}
}

func WithBit64(
	fn func(context.Context, *testing.B, assets.Dataset) (algorithm.Bit64, algorithm.Closer, error),
) StrategyOption {
	return func(s *strategy) error {
		if fn != nil {
			s.mode = algorithm.Float64
			s.initBit64 = fn
		}
		return nil
	}
}

func WithParallel() StrategyOption {
	return func(s *strategy) error {
		s.parallel = true
		return nil
	}
}
