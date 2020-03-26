package strategy

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type insert struct{}

func NewInsert(opts ...StrategyOption) benchmark.Strategy {
	i := new(insert)
	return newStrategy(append([]StrategyOption{
		WithPropName("Insert"),
		WithProp64(i.prop64),
		WithProp32(i.prop32),
	}, opts...)...)
}

func (i *insert) prop32(ctx context.Context, b *testing.B, c core.Core32, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
	n := int(atomic.LoadUint64(cnt))
	return c.Insert(dataset.Train()[n%len(dataset.Train())])
}

func (i *insert) prop64(ctx context.Context, b *testing.B, c core.Core64, dataset assets.Dataset, ids []uint, cnt *uint64) (interface{}, error) {
	n := int(atomic.LoadUint64(cnt))
	return c.Insert(dataset.TrainAsFloat64()[n%len(dataset.TrainAsFloat64())])
}
