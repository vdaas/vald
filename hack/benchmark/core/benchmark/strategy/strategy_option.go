package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type StrategyOption func(*strategy)

var (
	defaultStrategyOptions = []StrategyOption{}
)

func WithPreProp32(
	fn func(context.Context, *testing.B, core.Core32, assets.Dataset) (interface{}, error),
) StrategyOption {
	return func(s *strategy) {
		if fn != nil {
			s.preProp32 = fn
		}
	}
}

func WithProp32(
	fn func(context.Context, *testing.B, core.Core32, assets.Dataset, []uint, *uint64) (interface{}, error),
) StrategyOption {
	return func(s *strategy) {
		if fn != nil {
			s.prop32 = fn
		}
	}
}

func WithPreProp64(
	fn func(context.Context, *testing.B, core.Core64, assets.Dataset) (interface{}, error),
) StrategyOption {
	return func(s *strategy) {
		if fn != nil {
			s.preProp64 = fn
		}
	}
}

func WithProp64(
	fn func(context.Context, *testing.B, core.Core64, assets.Dataset, []uint, *uint64) (interface{}, error),
) StrategyOption {
	return func(s *strategy) {
		if fn != nil {
			s.prop64 = fn
		}
	}
}

func WithPropName(str string) StrategyOption {
	return func(s *strategy) {
		if len(str) != 0 {
			s.propName = str
		}
	}
}

func WithCore32Initializer(
	fn func(context.Context, assets.Dataset) (core.Core32, func(), error),
) StrategyOption {
	return func(s *strategy) {
		if fn != nil {
			s.initCore32 = fn
			s.mode = core.Float32
		}
	}
}

func WithCore64Initializer(
	fn func(context.Context, assets.Dataset) (core.Core64, func(), error),
) StrategyOption {
	return func(s *strategy) {
		if fn != nil {
			s.initCore64 = fn
			s.mode = core.Float64
		}
	}
}
