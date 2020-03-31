package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

func NewRemove(opts ...StrategyOption) benchmark.Strategy {
	return newStrategy(append([]StrategyOption{
		WithPropName("Delete"),
		WithPreProp32(
			func(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset) (ids []uint, err error) {
				return insertAndCreateIndex32(ctx, c, dataset)
			},
		),
		WithProp32(
			func(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset, ids []uint, cnt *uint64) (obj interface{}, err error) {
				err = c.Remove(ids[int(atomic.LoadUint64(cnt))%len(ids)])
				return
			},
		),
		WithPreProp64(
			func(ctx context.Context, b *testing.B, c core.Core64, dataset assets.Dataset) (ids []uint, err error) {
				return insertAndCreateIndex64(ctx, c, dataset)
			},
		),
		WithProp64(
			func(ctx context.Context, b *testing.B, c core.Core64, dataset assets.Dataset, ids []uint, cnt *uint64) (obj interface{}, err error) {
				err = c.Remove(ids[int(atomic.LoadUint64(cnt))%len(ids)])
				return
			},
		),
	}, opts...)...)
}
