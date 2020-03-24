package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
)

type PreStart func(context.Context, *testing.B, interface{}, assets.Dataset) (interface{}, error)

type defaultPreStart struct{}

func (d *defaultPreStart) PreStart(ctx context.Context, b *testing.B, c interface{}, dataset assets.Dataset) (interface{}, error) {
	switch core := c.(type) {
	case core.Core32:
		return d.float32(ctx, b, dataset, core)
	case core.Core64:
		return d.float64(ctx, b, dataset, core)
	default:
		b.Fatal("not implementated")
		return nil, nil
	}
}

func (d *defaultPreStart) float32(ctx context.Context, b *testing.B, dataset assets.Dataset, core core.Core32) (interface{}, error) {
	train := dataset.Train()

	ids := make([]uint, 0, len(train)*10)
	for i := 0; i < 10; i++ {
		for _, vec := range train {
			id, err := core.Insert(vec)
			if err != nil {
				return ids, err
			}
			ids = append(ids, id)
		}
	}

	if err := core.CreateIndex(10000); err != nil {
		return ids, err
	}

	return ids, nil
}

func (d *defaultPreStart) float64(ctx context.Context, b *testing.B, dataset assets.Dataset, core core.Core64) (interface{}, error) {
	train := dataset.TrainAsFloat64()

	ids := make([]uint, 0, len(train)*10)
	for i := 0; i < 10; i++ {
		for _, vec := range train {
			id, err := core.Insert(vec)
			if err != nil {
				return ids, err
			}
			ids = append(ids, id)
		}
	}

	if err := core.CreateIndex(10000); err != nil {
		return ids, err
	}

	return ids, nil
}
