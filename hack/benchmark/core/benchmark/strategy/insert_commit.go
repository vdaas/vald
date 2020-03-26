package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

func NewInsertCommit(poolSize uint32, opts ...StrategyOption) benchmark.Strategy {
	return newStrategy(append([]StrategyOption{
		WithPropName("BulkInsertCommit"),
		WithPreProp32(func(context.Context, *testing.B, core.Core32, assets.Dataset) ([]uint, error) {
			return nil, nil
		}),
		WithProp32(
			func(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				train := dataset.Train()
				return c.InsertCommit(train[int(atomic.LoadUint64(cnt))%len(train)], poolSize)
			},
		),
		WithPreProp64(func(context.Context, *testing.B, core.Core64, assets.Dataset) ([]uint, error) {
			return nil, nil
		}),
		WithProp64(
			func(ctx context.Context, b *testing.B, c core.Core64, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				train := dataset.TrainAsFloat64()
				return c.InsertCommit(train[int(atomic.LoadUint64(cnt))%len(train)], poolSize)
			},
		),
	}, opts...)...)
}
