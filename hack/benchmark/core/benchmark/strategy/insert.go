package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

func NewInsert(opts ...StrategyOption) benchmark.Strategy {
	return newStrategy(append([]StrategyOption{
		WithPropName("Insert"),
		WithProp32(
			func(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				train := dataset.Train()
				return c.Insert(train[int(atomic.LoadUint64(cnt))%len(train)])
			},
		),
		WithProp64(
			func(ctx context.Context, b *testing.B, c core.Core64, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				train := dataset.TrainAsFloat64()
				return c.Insert(train[int(atomic.LoadUint64(cnt))%len(train)])
			},
		),
	}, opts...)...)
}
