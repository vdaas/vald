package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

func NewBulkInsert(opts ...StrategyOption) benchmark.Strategy {
	return newStrategy(append([]StrategyOption{
		WithPropName("BulkInsert"),
		WithPreProp32(func(context.Context, *testing.B, core.Core32, assets.Dataset) ([]uint, error) {
			return nil, nil
		}),
		WithProp32(
			func(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				ids, errs := c.BulkInsert(dataset.Train())
				return ids, wrapErrors(errs)
			},
		),
		WithPreProp64(func(context.Context, *testing.B, core.Core64, assets.Dataset) ([]uint, error) {
			return nil, nil
		}),
		WithProp64(
			func(ctx context.Context, b *testing.B, c core.Core64, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				ids, errs := c.BulkInsert(dataset.TrainAsFloat64())
				return ids, wrapErrors(errs)
			},
		),
	}, opts...)...)
}
