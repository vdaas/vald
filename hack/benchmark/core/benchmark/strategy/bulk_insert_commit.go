package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

func NewBulkInsertCommit(poolSize uint32, opts ...StrategyOption) benchmark.Strategy {
	return newStrategy(append([]StrategyOption{
		WithPropName("BulkInsertCommit"),
		WithProp32(
			func(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				ids, errs := c.BulkInsertCommit(dataset.Train(), poolSize)
				return ids, wrapErrors(errs)
			},
		),
		WithProp64(
			func(ctx context.Context, b *testing.B, c core.Core64, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				ids, errs := c.BulkInsertCommit(dataset.TrainAsFloat64(), poolSize)
				return ids, wrapErrors(errs)
			},
		),
	}, opts...)...)
}
