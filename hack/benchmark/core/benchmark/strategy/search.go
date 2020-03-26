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
				train := dataset.Train()

				ids = make([]uint, 0, len(train))
				var id uint
				for _, v := range train {
					id, err = c.Insert(v)
					if err != nil {
						return nil, err
					}
					ids = append(ids, id)
				}

				err = c.CreateIndex(10)
				if err != nil {
					return nil, err
				}

				return
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
				train := dataset.TrainAsFloat64()

				ids = make([]uint, 0, len(train))
				var id uint
				for _, v := range train {
					id, err = c.Insert(v)
					if err != nil {
						return nil, err
					}
					ids = append(ids, id)
				}

				err = c.CreateIndex(10)
				if err != nil {
					return nil, err
				}

				return
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
