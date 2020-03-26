package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

func NewDelete(opts ...StrategyOption) benchmark.Strategy {
	return newStrategy(append([]StrategyOption{
		WithPropName("Delete"),
		WithPreProp32(
			func(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset) (ids []uint, err error) {
				train := dataset.Train()

				n := 1000
				ids = make([]uint, 0, len(train)*n)

				for i := 0; i < n; i++ {
					inserted, errs := c.BulkInsert(train)
					err = wrapErrors(errs)
					if err != nil {
						return nil, err
					}
					ids = append(ids, inserted...)
				}

				err = c.CreateIndex(10)
				if err != nil {
					return nil, err
				}

				return
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
				train := dataset.TrainAsFloat64()

				n := 1000
				ids = make([]uint, 0, len(train)*n)

				for i := 0; i < n; i++ {
					inserted, errs := c.BulkInsert(train)
					err = wrapErrors(errs)
					if err != nil {
						return nil, err
					}
					ids = append(ids, inserted...)
				}

				err = c.CreateIndex(10)
				if err != nil {
					return nil, err
				}

				return
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
