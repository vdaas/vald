package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

func NewSearch(size int, epsilon, radius float32, opts ...StrategyOption) benchmark.Strategy {
	return newStrategy(append([]StrategyOption{
		WithPropName("Search"),
		WithPreProp32(
			func(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset) (ids []uint, err error) {
				return insertAndCreateIndex32(ctx, c, dataset)
			},
		),
		WithProp32(
			func(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				query := dataset.Query()
				return c.Search(query[int(atomic.LoadUint64(cnt))%len(query)], size, epsilon, radius)
			},
		),
		WithPreProp64(
			func(ctx context.Context, b *testing.B, c core.Core64, dataset assets.Dataset) (ids []uint, err error) {
				return insertAndCreateIndex64(ctx, c, dataset)
			},
		),
		WithProp64(
			func(ctx context.Context, b *testing.B, c core.Core64, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
				query := dataset.QueryAsFloat64()
				return c.Search(query[int(atomic.LoadUint64(cnt))%len(query)], size, epsilon, radius)
			},
		),
	}, opts...)...)
}
