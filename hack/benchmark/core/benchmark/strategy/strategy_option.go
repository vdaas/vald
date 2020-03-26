package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type StrategyOption func(*strategy) error

var (
	defaultStrategyOptions = []StrategyOption{}
)

func WithPreProp32(
	fn func(context.Context, *testing.B, core.Core32, assets.Dataset) ([]uint, error),
) StrategyOption {
	return func(s *strategy) error {
		if fn != nil {
			s.preProp32 = fn
		}
		return nil
	}
}

func WithProp32(
	fn func(context.Context, *testing.B, core.Core32, assets.Dataset, []uint, *uint64) (interface{}, error),
) StrategyOption {
	return func(s *strategy) error {
		if fn != nil && s.mode == core.Float32 {
			s.prop32 = fn
		}
		return nil
	}
}

func WithPreProp64(
	fn func(context.Context, *testing.B, core.Core64, assets.Dataset) ([]uint, error),
) StrategyOption {
	return func(s *strategy) error {
		if fn != nil {
			s.preProp64 = fn
		}
		return nil
	}
}

func WithProp64(
	fn func(context.Context, *testing.B, core.Core64, assets.Dataset, []uint, *uint64) (interface{}, error),
) StrategyOption {
	return func(s *strategy) error {
		if fn != nil && s.mode == core.Float64 {
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

func WithCore32(
	fn func(context.Context, *testing.B, assets.Dataset) (core.Core32, core.Closer, error),
) StrategyOption {
	return func(s *strategy) (err error) {
		if fn != nil {
			s.mode = core.Float32
			s.initCore32 = fn
		}
		return nil
	}
}

func WithCore64(
	fn func(context.Context, *testing.B, assets.Dataset) (core.Core64, core.Closer, error),
) StrategyOption {
	return func(s *strategy) error {
		if fn != nil {
			s.mode = core.Float64
			s.initCore64 = fn
		}
		return nil
	}
}
